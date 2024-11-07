# Instance Metadata Service

Instance metadata is data about your instance that you can use to configure or manage the running instance which you can access from a running instance using one of the following methods:

- Instance Metadata Service Version 1 (IMDSv1) – a request/response method
- Instance Metadata Service Version 2 (IMDSv2) – a session-oriented method

CAPA defaults to use IMDSv2 as optional property when creating instances.

CAPA expose options to configure IMDSv2 as required when creating instances, as it provides a [better level of security](https://aws.amazon.com/blogs/security/defense-in-depth-open-firewalls-reverse-proxies-ssrf-vulnerabilities-ec2-instance-metadata-service/).

It is possible to configure the instance metadata options using the field called `instanceMetadataOptions` in the `AWSMachineTemplate`.

Example:

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSMachineTemplate
metadata:
  name: "test"
spec:
  template:
    spec:
      instanceMetadataOptions:
        httpEndpoint: enabled
        httpPutResponseHopLimit: 1
        httpTokens: optional
        instanceMetadataTags: disabled
```

To use IMDSv2, simply set `httpTokens` value to `required` (in other words, set the use of IMDSv2 to required).
To use IMDSv2, please also set `httpPutResponseHopLimit` value to `2`, as it is recommended in container environment according to [AWS document](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html#imds-considerations).

Similarly, this can be done with `AWSManagedMachinePool` for use with EKS Managed Nodegroups. One slight difference here is that you [must use Launch Templates to configure IMDSv2 with Autoscaling Groups](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-metadata-transition-to-version-2.html). In order to configure the LaunchTemplate, you must use a custom AMI type according to the AWS API. This can be done by setting `AWSManagedMachinePool.spec.amiType` to `CUSTOM`. This change means that you must also specify a bootstrapping script to the worker node, which allows it to be joined to the EKS cluster. The default AWS Managed Node Group bootstrap script can be found [here on Github](https://github.com/awslabs/amazon-eks-ami/blob/master/files/bootstrap.sh).

The following example will use the default Amazon EKS Worker Node AMI which includes the default EKS Bootstrapping script. This must be installed on the management cluster as a Secret, under the key `value`. The secret's name must then be included in your `MachinePool` manifest at `MachinePool.spec.template.spec.bootstrap.dataSecretName`. Some assumptions are made for this example:

- Your cluster name is `capi-imds`, which CAPA renames to `default_capi-imds-control-plane` automatically
- Your cluster is Kubernetes Version `v1.25.9`
- Your `AWSManagedCluster` is deployed in the `default` namespace along with the bootstrap secret `eks-bootstrap`

```yaml
kind: Secret
apiVersion: v1
type: Opaque
data:
  value: IyEvYmluL2Jhc2ggLXhlCi9ldGMvZWtzL2Jvb3RzdHJhcC5zaCBkZWZhdWx0X2NhcGktaW1kcy1jb250cm9sLXBsYW5l
metadata:
  name: eks-bootstrap
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSManagedMachinePool
metadata:
  name: "capi-imds-pool-launchtemplate"
spec:
  amiType: CUSTOM
  awsLaunchTemplate:
    name: my-aws-launch-template
    instanceType: t3.nano
    metadataOptions:
      httpTokens: required
      httpPutResponseHopLimit: 2
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachinePool
metadata:
  name: "capi-imds-pool-1"
spec:
  clusterName: "capi-imds"
  replicas: 1
  template:
    spec:
      version: v1.25.9
      clusterName: "capi-imds"
      bootstrap:
        dataSecretName: "eks-bootstrap"
      infrastructureRef:
        name: "capi-imds-pool-launchtemplate"
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
        kind: AWSManagedMachinePool
```

`IyEvYmluL2Jhc2ggLXhlCi9ldGMvZWtzL2Jvb3RzdHJhcC5zaCBkZWZhdWx0X2NhcGktaW1kcy1jb250cm9sLXBsYW5l` in the above secret is a Base64 encoded version of the following script:

```bash
#!/bin/bash -xe
/etc/eks/bootstrap.sh default_capi-imds-control-plane
```

If your cluster is not named `default_capi-imds-control-plane` in the AWS EKS console, you must update the name and store it as a Secret again.

See [the CLI command reference](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/modify-instance-metadata-options.html) for more information.

Before you decide to use IMDSv2 for the cluster instances, please make sure all your applications are compatible with IMDSv2.

See the [transition guide](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-metadata-transition-to-version-2.html#recommended-path-for-requiring-imdsv2) for more information.

# Instance Metadata Service

Instance metadata is data about your instance that you can use to configure or manage the running instance which you can access from a running instance using one of the following methods:

* Instance Metadata Service Version 1 (IMDSv1) – a request/response method
* Instance Metadata Service Version 2 (IMDSv2) – a session-oriented method

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

See [the CLI command reference](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/modify-instance-metadata-options.html) for more information.

Before you decide to use IMDSv2 for the cluster instances, please make sure all your applications are compatible to IMDSv2. 

See the [transition guide](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-metadata-transition-to-version-2.html#recommended-path-for-requiring-imdsv2) for more information.

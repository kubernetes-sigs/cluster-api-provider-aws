# Instance Metadata Service

Instance metadata is data about your instance that you can use to configure or manage the running instance which you can access from a running instance using one of the following methods:

* Instance Metadata Service Version 1 (IMDSv1) – a request/response method
* Instance Metadata Service Version 2 (IMDSv2) – a session-oriented method

CAPA defaults to IMDSv2 when creating instances, as it provides a [better level of security](https://aws.amazon.com/blogs/security/defense-in-depth-open-firewalls-reverse-proxies-ssrf-vulnerabilities-ec2-instance-metadata-service/).
CAPA defaults to 2 hot limit when creating instances with IMDSv2, as it is recommended in container environment according to [AWS document](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html#imds-considerations).

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
        httpPutResponseHopLimit: 2
        httpTokens: required
        instanceMetadataTags: disabled
```

To use IMDSv1, simply set `httpTokens` value to `optional` (in other words, set the use of IMDSv2 to optional).

See [the CLI command reference](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/modify-instance-metadata-options.html) for more information.

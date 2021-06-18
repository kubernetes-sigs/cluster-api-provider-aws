# Enabling Encryption

To enable encryption when creating a cluster you need to create a new KMS key that has an alias name starting with `cluster-api-provider-aws-`.

For example, `arn:aws:kms:eu-north-1:12345678901:alias/cluster-api-provider-aws-key1`.

You then need to specify this alias in the `encryptionConfig` of the `AWSManagedControlPlane`:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1alpha4
metadata:
  name: "capi-managed-test-control-plane"
spec:
  ...
  encryptionConfig:
    provider: "arn:aws:kms:eu-north-1:12345678901:alias/cluster-api-provider-aws-key1"
    resources:
    - "secrets"
```

## Custom KMS Alias Prefix

If you would like to use a different alias prefix then you can use the `kmsAliasPrefix` in the optional configuration file for **clusterawsadm**:

```bash
clusterawsadm bootstrap iam create-stack --config custom-prefix.yaml
```

And the contents of the configuration file:

```yaml
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1
kind: AWSIAMConfiguration
spec:
  eks:
    enable: true
    kmsAliasPrefix: "my-prefix-*
```
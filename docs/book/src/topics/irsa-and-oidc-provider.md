# Service Account Issuer Discovery and AWS IAM Roles for Service Accounts (IRSA)

## Overview

The Cluster API for AWS provider can optionally publish the Kubernetes service account token issuer to a public location
and configure AWS to trust it to authenticate Kubernetes service accounts. 

This enables the IAM Roles for Service Accounts (IRSA) capability on the cluster, allowing workloads running within the
cluster to authenticate to the AWS API and/or to 3rd party APIs that are able to trust OIDC providers.

This document applies only to unmanaged Kubernetes clusters created by this provider. EKS does these actions automatically for
Managed clusters.

## How to enable

### 1 - Enable Experiment feature flag

* The creation of S3 buckets and association of OIDC provider is not allowed unless the following experimental feature flags is enabled:
    * `OIDCProviderSupport` (`EXP_OIDC_PROVIDER_SUPPORT`)

### 2 - Set S3 Bucket configuration

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: capa-quickstart
spec:
  ...
  s3Bucket:
    controlPlaneIAMInstanceProfile: control-plane.cluster-api-provider-aws.sigs.k8s.io
    nodesIAMInstanceProfiles:
      - nodes.cluster-api-provider-aws.sigs.k8s.io
    name: my-cluster-bucket
```

### 3 - Enable `spec.associateOIDCProvider`

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
kind: AWSCluster
metadata:
  name: capa-quickstart
spec:
  ...
  associateOIDCProvider: true
```


### 4 - Set Service Account Issuer URL in KubeadmControlPlane Configuration

The KubeadmControlPlane configuration should be updated to set the apiServer `service-account-issuer` argument with the S3 buckets URL:

```yaml
spec:
  kubeadmConfigSpec:
    clusterConfiguration:
      apiServer:
        extraArgs:
          cloud-provider: external
          service-account-issuer: https://s3.{{ REGION }}.amazonaws.com/cluster-api-aws-provider-{{ CLUSTER_NAME }}/{{ CLUSTER_NAME }}
```

### 5 - Install AWS Pod Identity Webhook addon on workload cluster

AWS provides an open source cluster addon that operates as a mutating webhook for pods that require AWS IAM access.
The GitHub repository contains [example manifests](https://github.com/aws/amazon-eks-pod-identity-webhook/tree/master/deploy) for installing the addon.
These manifests can be applied via your preferred addon mechanism, such as Cluster API [ClusterResourceSets](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-resource-set).

>**IMPORTANT NOTE:**
> The Amazon Pod Identity Webhook addon requires cert-manager be available on the workload cluster. Please ensure this is also installed.

## How it works

The provider will perform the following actions if this feature is enabled:

1. Create an S3 bucket.
2. Upload an OpenID Connect discovery document to `<cluster_name>/.well-known/openid-configuration` and make the object world readable.

    *Example contents of Discovery Document*
    ```json
    {
      "issuer": "https://s3.ap-southeast-2.amazonaws.com/cluster-api-aws-provider-capa-quickstart-example-com/capa-quickstart",
      "jwks_uri": "https://s3.ap-southeast-2.amazonaws.com/cluster-api-aws-provider-capa-quickstart-example-com/capa-quickstart/openid/v1/jwks",
      "authorization_endpoint": "urn:kubernetes:programmatic_authorization",
      "response_types_supported": [
        "id_token"
      ],
      "subject_types_supported": [
        "public"
      ],
      "id_token_signing_alg_values_supported": [
        "RS256"
      ],
      "claims_supported": [
        "sub",
        "iss"
      ]
    }
    ```

3. Upload the Kubernetes service account public signing key to the bucket as `<cluster_name>/openid/v1/jwks` and make the object world readable.
4. Create an OpenID Connect identity provider in AWS IAM, configured with the S3 URL.

## Reference

* https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc.html

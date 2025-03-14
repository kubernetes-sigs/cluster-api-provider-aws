# Creating a ROSA HCP cluster

## Permissions
### Authentication using service account credentials
CAPA controller requires service account credentials to be able to provision ROSA HCP clusters:
1. Visit [https://console.redhat.com/iam/service-accounts](https://console.redhat.com/iam/service-accounts) and create a service account. If you already have a service account, you can skip this step.

   For every newly created service account, make sure to activate the account using the [ROSA command line tool](https://github.com/openshift/rosa). First, log in using your newly created service account
   ```shell
   rosa login --client-id ... --client-secret ...
   ```
   Then activate your service account
   ```shell
   rosa whoami
   ```

1. Create a new kubernetes secret with the service account credentials to be referenced later by `ROSAControlPlane`
    ```shell
    kubectl create secret generic rosa-creds-secret \
      --from-literal=ocmClientID='....' \
      --from-literal=ocmClientSecret='eyJhbGciOiJIUzI1NiIsI....' \
      --from-literal=ocmApiUrl='https://api.openshift.com'
    ```
    Note: to consume the secret without the need to reference it from your `ROSAControlPlane`, name your secret as `rosa-creds-secret` and create it in the CAPA manager namespace (usually `capa-system`)
    ```shell
    kubectl -n capa-system create secret generic rosa-creds-secret \
      --from-literal=ocmClientID='....' \
      --from-literal=ocmClientSecret='eyJhbGciOiJIUzI1NiIsI....' \
      --from-literal=ocmApiUrl='https://api.openshift.com'
    ```


### Authentication using SSO offline token (DEPRECATED)
The SSO offline token is being deprecated and it is recommended to use service account credentials instead, as described above.

1. Visit https://console.redhat.com/openshift/token to retrieve your SSO offline authentication token

1. Create a credentials secret within the target namespace with the token to be referenced later by `ROSAControlePlane`
    ```shell
        kubectl create secret generic rosa-creds-secret \
            --from-literal=ocmToken='eyJhbGciOiJIUzI1NiIsI....' \
        --from-literal=ocmApiUrl='https://api.openshift.com'
    ```
    Alternatively, you can edit the CAPA controller deployment to provide the credentials
    ```shell
        kubectl edit deployment -n capa-system capa-controller-manager
    ```
    and add the following environment variables to the manager container
    ```yaml
        env:
          - name: OCM_TOKEN
            value: "<token>"
          - name: OCM_API_URL
            value: "https://api.openshift.com" # or https://api.stage.openshift.com
    ```

### Migration from offline token to service account authentication

1. Visit [https://console.redhat.com/iam/service-accounts](https://console.redhat.com/iam/service-accounts) and create a new service account.

1. If you previously used kubernetes secret to specify the OCM credentials secret, edit the secret:
    ```shell
        kubectl edit secret rosa-creds-secret
    ```
    where you will remove the `ocmToken` credentials and add base64 encoded `ocmClientID` and `ocmClientSecret` credentials like so:
    ```yaml
    apiVersion: v1
     data:
       ocmApiUrl: aHR0cHM6Ly9hcGkub3BlbnNoaWZ0LmNvbQ==
       ocmClientID: Y2xpZW50X2lk...
       ocmClientSecret: Y2xpZW50X3NlY3JldA==...
     kind: Secret
     type: Opaque
    ```

1. If you previously used capa manager deployment to specify the OCM offline token as environment variable, edit the manager deployment
    ```shell
        kubectl -n capa-system edit deployment capa-controller-manager
    ```
    and remove the `OCM_TOKEN` and `OCM_API_URL` variables, followed by `kubectl -n capa-system rollout restart deploy capa-controller-manager`. Then create the new default secret in the `capa-system` namespace with
    ```shell
        kubectl -n capa-system create secret generic rosa-creds-secret \
          --from-literal=ocmClientID='....' \
          --from-literal=ocmClientSecret='eyJhbGciOiJIUzI1NiIsI....' \
          --from-literal=ocmApiUrl='https://api.openshift.com'
    ```

## Prerequisites

Follow the guide [here](https://docs.aws.amazon.com/ROSA/latest/userguide/getting-started-hcp.html) up until ["Create a ROSA with HCP Cluster"](https://docs.aws.amazon.com/ROSA/latest/userguide/getting-started-hcp.html#create-hcp-cluster-cli) to install the required tools and setup the prerequisite infrastructure. Once Step 3 is done, you will be ready to proceed with creating a ROSA HCP cluster using cluster-api.

## Creating the cluster

1. Prepare the environment:
    ```bash
    export OPENSHIFT_VERSION="4.14.5"
    export AWS_REGION="us-west-2"
    export AWS_AVAILABILITY_ZONE="us-west-2a"
    export AWS_ACCOUNT_ID="<account_id>"
    export AWS_CREATOR_ARN="<user_arn>" # can be retrieved e.g. using `aws sts get-caller-identity`

    export OIDC_CONFIG_ID="<oidc_id>" # OIDC config id creating previously with `rosa create oidc-config`
    export ACCOUNT_ROLES_PREFIX="ManagedOpenShift-HCP" # prefix used to create account IAM roles with `rosa create account-roles`
    export OPERATOR_ROLES_PREFIX="capi-rosa-quickstart"  # prefix used to create operator roles with `rosa create operator-roles --prefix <PREFIX_NAME>`

    # subnet IDs created earlier
    export PUBLIC_SUBNET_ID="subnet-0b54a1111111111111"
    export PRIVATE_SUBNET_ID="subnet-05e72222222222222"
    ```

1. Render the cluster manifest using the ROSA HCP cluster template:
    ```shell
    clusterctl generate cluster <cluster-name> --from templates/cluster-template-rosa.yaml > rosa-capi-cluster.yaml
    ```
    Note: The AWS role name must be no more than 64 characters in length. Otherwise an error will be returned. Truncate values exceeding 64 characters.

1. If a credentials secret was created earlier, edit `ROSAControlPlane` to reference it:
    ```yaml
    apiVersion: controlplane.cluster.x-k8s.io/v1beta2
    kind: ROSAControlPlane
    metadata:
      name: "capi-rosa-quickstart-control-plane"
    spec:
      credentialsSecretRef:
        name: rosa-creds-secret
    ...
    ```

1. Provide an AWS identity reference
    ```yaml
    apiVersion: controlplane.cluster.x-k8s.io/v1beta2
    kind: ROSAControlPlane
    metadata:
      name: "capi-rosa-quickstart-control-plane"
    spec:
      identityRef:
        kind: <IdentityType>
        name: <IdentityName>
    ...
    ```

    Otherwise, make sure the following `AWSClusterControllerIdentity` singleton exists in your management cluster:
    ```yaml
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
    kind: AWSClusterControllerIdentity
    metadata:
      name: "default"
    spec:
      allowedNamespaces: {}  # matches all namespaces
    ```

    see [Multi-tenancy](../multitenancy.md) for more details

1. Finally apply the manifest to create your Rosa cluster:
    ```shell
    kubectl apply -f rosa-capi-cluster.yaml
    ```

see [ROSAControlPlane CRD Reference](https://cluster-api-aws.sigs.k8s.io/crd/#controlplane.cluster.x-k8s.io/v1beta2.ROSAControlPlane) for all possible configurations.

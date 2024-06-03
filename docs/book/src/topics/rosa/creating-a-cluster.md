# Creating a ROSA cluster

## Permissions
CAPA controller requires an API token in order to be able to provision ROSA clusters:

1. Visit [https://console.redhat.com/openshift/token](https://console.redhat.com/openshift/token) to retrieve your API authentication token

1. Create a credentials secret within the target namespace with the token to be referenced later by `ROSAControlePlane`
    ```shell
    kubectl create secret generic rosa-creds-secret \
      --from-literal=ocmToken='eyJhbGciOiJIUzI1NiIsI....' \
      --from-literal=ocmApiUrl='https://api.openshift.com' 
    ```

    Alternatively, you can edit CAPA controller deployment to provide the credentials:
    ```shell
    kubectl edit deployment -n capa-system capa-controller-manager
    ```

    and add the following environment variables to the manager container:
    ```yaml
      env:
      - name: OCM_TOKEN
        value: "<token>"
      - name: OCM_API_URL
        value: "https://api.openshift.com" # or https://api.stage.openshift.com
    ```

## Prerequisites

Follow the guide [here](https://docs.aws.amazon.com/ROSA/latest/userguide/getting-started-hcp.html) up until [Step 3](https://docs.aws.amazon.com/ROSA/latest/userguide/getting-started-hcp.html#getting-started-hcp-step-3) 
to install the required tools and setup the prerequisite infrastructure.
Once Step 3 is done, you will be ready to proceed with creating a ROSA cluster using cluster-api.

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

1. Render the cluster manifest using the ROSA cluster template:
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

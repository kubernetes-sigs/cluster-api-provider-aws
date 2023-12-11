# Creating a ROSA cluster

## Permissions
CAPA controller requires an API token in order to be able to provision ROSA clusters:

1. Visit [https://console.redhat.com/openshift/token](https://console.redhat.com/openshift/token) to retrieve your API authentication token

2. Edit CAPA controller deployment:
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
    export OPENSHIFT_VERSION="openshift-v4.12.15"
    export CLUSTER_NAME="capi-rosa-quickstart"
    export AWS_REGION="us-west-2"
    export AWS_AVAILABILITY_ZONE="us-west-2a"
    export AWS_ACCOUNT_ID="<account_id"
    export AWS_CREATOR_ARN="<user_arn>" # can be retrieved e.g. using `aws sts get-caller-identity`

    export OIDC_CONFIG_ID="<oidc_id>" # OIDC config id creating previously with `rosa create oidc-config`
    export ACCOUNT_ROLES_PREFIX="ManagedOpenShift-HCP" # prefix used to create account IAM roles with `rosa create account-roles`
    export OPERATER_ROLES_PREFIX="capi-rosa-quickstart"  # prefix used to create operator roles with `rosa create operator-roles --prefix <PREFIX_NAME>`

    # subnet IDs created earlier
    export PUBLIC_SUBNET_ID="subnet-0b54a1111111111111"   
    export PRIVATE_SUBNET_ID="subnet-05e72222222222222"
    ```

1. Create a cluster using the ROSA cluster template:
    ```bash
    cat templates/cluster-template-rosa.yaml | envsubst > rosa-capi-cluster.yaml

    kubectl apply -f rosa-capi-cluster.yaml
    ```

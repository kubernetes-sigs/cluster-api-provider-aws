# Creating a ROSA HCP cluster

## Prerequisites

1. Create a management cluster using the [Quick Start Guide.](https://cluster-api-aws.sigs.k8s.io/quick-start) 


2. Install the required tools and set up the prerequisite infrastructure using the [ROSA Setup guide](https://docs.aws.amazon.com/rosa/latest/userguide/set-up.html).
   
Once these steps are complete, you are ready to create a ROSA HCP cluster.


## Authentication
The CAPA controller requires service account credentials to provision ROSA HCP clusters.  
If you already have a service account, you can skip these steps.
1. Create a service account by visiting [https://console.redhat.com/iam/service-accounts](https://console.redhat.com/iam/service-accounts).


2.   For every newly created service account, make sure to activate the account using the [ROSA command line tool](https://github.com/openshift/rosa). 
     First, log in using your newly created service account:
      ```shell
      rosa login --client-id ... --client-secret ...
      ```
3.   Then activate your service account:
      ```shell
      rosa whoami
      ```
## Permissions


1. Create a new kubernetes secret with the service account credentials to be referenced later by the `ROSAControlPlane`
    ```shell
    kubectl create secret generic rosa-creds-secret \
      --from-literal=ocmClientID='....' \
      --from-literal=ocmClientSecret='eyJhbGciOiJIUzI1NiIsI....' \
      --from-literal=ocmApiUrl='https://api.openshift.com'
    ```
    Note: to consume the secret without the need to reference it from your `ROSAControlPlane`, name your secret `rosa-creds-secret` and create it in the CAPA manager namespace (usually `capa-system`)
    ```shell
    kubectl -n capa-system create secret generic rosa-creds-secret \
      --from-literal=ocmClientID='....' \
      --from-literal=ocmClientSecret='eyJhbGciOiJIUzI1NiIsI....' \
      --from-literal=ocmApiUrl='https://api.openshift.com'
    ```


## Creating the cluster

1. Save the following to a file `rosa-role-network.yaml`:

    ```yaml
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
    kind: ROSARoleConfig
    metadata:
      name: "role-config"
    spec:
      accountRoleConfig:
        prefix: "rosa"
        version: "4.19.0"
      operatorRoleConfig:
        prefix: "rosa"
      credentialsSecretRef:
        name: rosa-creds-secret
      oidcProviderType: Managed
    ---
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
    kind: ROSANetwork
    metadata:
      name: "rosa-vpc"
    spec:
      region: "us-west-2"
      stackName: "rosa-hcp-net"
      availabilityZones:
      - "us-west-2a"
      - "us-west-2b"
      - "us-west-2c"
      cidrBlock: 10.0.0.0/16
    ```

    Apply the manifest:

    ```shell
    kubectl apply -f rosa-role-network.yaml
    ```

    Verify the `ROSARoleConfig` was successfully created. The status should contain the `accountRolesRef`, `oidcID`, `oidcProviderARN` and `operatorRolesRef`:

    ```shell
    kubectl get rosaroleconfig role-config -o yaml
    ```

    Verify the `ROSANetwork` was successfully created. The status should contain the created subnets:

    ```shell
    kubectl get rosanetwork rosa-vpc -o yaml
    ```

1. Save the following to a file `rosa-cluster.yaml`:

    ```yaml
    apiVersion: cluster.x-k8s.io/v1beta1
    kind: Cluster
    metadata:
      name: "rosa-hcp-1"
    spec:
      clusterNetwork:
        pods:
          cidrBlocks: ["192.168.0.0/16"]
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
        kind: ROSACluster
        name: "rosa-hcp-1"
      controlPlaneRef:
        apiVersion: controlplane.cluster.x-k8s.io/v1beta2
        kind: ROSAControlPlane
        name: "rosa-hcp-1-control-plane"
    ---
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
    kind: ROSACluster
    metadata:
      name: "rosa-hcp-1"
    spec: {}
    ---
    apiVersion: controlplane.cluster.x-k8s.io/v1beta2
    kind: ROSAControlPlane
    metadata:
      name: "rosa-hcp-1-control-plane"
    spec:
      credentialsSecretRef:
        name: rosa-creds-secret
      rosaClusterName: rosa-hcp-1
      domainPrefix: rosa-hcp
      rosaRoleConfigRef:
        name: role-config
      version: "4.19.0"
      region: "us-west-2"
      rosaNetworkRef:
        name: "rosa-vpc"
      network:
        machineCIDR: "10.0.0.0/16"
        podCIDR: "10.128.0.0/14"
        serviceCIDR: "172.30.0.0/16"
      defaultMachinePoolSpec:
        instanceType: "m5.xlarge"
        autoscaling:
          maxReplicas: 6
          minReplicas: 3
      additionalTags:
        env: "demo"
    ```

    Apply the manifest:

    ```shell
    kubectl apply -f rosa-cluster.yaml
    ```

1. Provide an AWS identity reference by adding an `identityRef` to the `ROSAControlPlane` spec:

    ```yaml
    apiVersion: controlplane.cluster.x-k8s.io/v1beta2
    kind: ROSAControlPlane
    metadata:
      name: "rosa-hcp-1-control-plane"
    spec:
      identityRef:
        kind: <IdentityType>
        name: <IdentityName>
    ...
    ```

    Otherwise, make sure the following `AWSClusterControllerIdentity` singleton exists in your management cluster. Save it to a file and apply it:

    ```yaml
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
    kind: AWSClusterControllerIdentity
    metadata:
      name: "default"
    spec:
      allowedNamespaces: {}  # matches all namespaces
    ```

    ```shell
    kubectl apply -f <filename>.yaml
    ```

    see [Multi-tenancy](../multitenancy.md) for more details

1. Check the `ROSAControlPlane` status:

    ```shell
    kubectl get ROSAControlPlane rosa-hcp-1-control-plane

    NAME                       CLUSTER      READY
    rosa-hcp-1-control-plane   rosa-hcp-1   true
    ```

    The ROSA HCP cluster can take around 40 minutes to be fully provisioned.

1. After provisioning has completed, verify the `ROSAMachinePool` resources were successfully created:

    ```shell
    kubectl get ROSAMachinePool

    NAME        READY   REPLICAS
    workers-0   true    1
    workers-1   true    1
    workers-2   true    1
    ```

    **Note:** The number of default `ROSAMachinePool` resources corresponds to the number of availability zones configured.

see [ROSAControlPlane CRD Reference](https://cluster-api-aws.sigs.k8s.io/crd/#controlplane.cluster.x-k8s.io/v1beta2.ROSAControlPlane) for all possible configurations.

# External Auth Providers (BYOI)

ROSA HCP allows you to Bring Your Own Identity (BYOI) to manage and authenticate cluster users.

## Enabling

To enable this feature, `enableExternalAuthProviders` field should be set to `true` on cluster creation. Changing this field afterwards will have no effect:
```yaml
---
apiVersion: controlplane.cluster.x-k8s.io/v1beta2
kind: ROSAControlPlane
metadata:
  name: "capi-rosa-quickstart-control-plane"
spec:
  enableExternalAuthProviders: true
  ....
```

Note: This feature requires OpenShift version `4.15.5` or newer.

## Usage

After creating and configuring your OIDC provider of choice, the next step is to configure ROSAControlPlane `externalAuthProviders` as follows:
```yaml
---
apiVersion: controlplane.cluster.x-k8s.io/v1beta2
kind: ROSAControlPlane
metadata:
  name: "capi-rosa-quickstart-control-plane"
spec:
  enableExternalAuthProviders: true
  externalAuthProviders:
  - name: my-oidc-provider
    issuer:
      issuerURL: https://login.microsoftonline.com/<tenant-id>/v2.0 # e.g. if using Microsoft Entra ID
      audiences:  # audiences that will be trusted by the kube-apiserver
      - "audience1" # usually the client ID
    claimMappings:
      username:
        claim: email
        prefixPolicy: ""
      groups:
        claim: groups
  ....
```

Note: `oidcProviders` only accepts one entry at the moment.

## Accessing the cluster

### Setting up RBAC

When `enableExternalAuthProviders` is set to `true`, ROSA provider will generate a temporary admin kubeconfig secret in the same namespace named `<cluster-name>-bootstrap-kubeconfig`. This kubeconfig can be used to access the cluster to setup RBAC for OIDC users/groups.

The following example binds the `cluster-admin` role to an OIDC group, giving all users in that group admin permissions.
```shell
kubectl get secret <cluster-name>-bootstrap-kubeconfig -o jsonpath='{.data.value}' | base64 -d > /tmp/capi-admin-kubeconfig
export KUBECONFIG=/tmp/capi-admin-kubeconfig

kubectl create clusterrolebinding oidc-cluster-admins --clusterrole cluster-admin --group <group-id>
```

Note: The generated bootstrap kubeconfig is only valid for 24h, and will not be usable afterwards. However, users can opt to manually delete the secret object to trigger the generation of a new one which will be valid for another 24h.

### Login using the cli

The [kubelogin kubectl plugin](https://github.com/int128/kubelogin/tree/master) can be used to login with OIDC credentials using the cli. 

### Configuring OpenShift Console

The OpenShift Console needs to be configured before it can be used to authenticate and login to the cluster. 
1. Setup a new client in your OIDC provider with the following Redirect URL: `<console-url>/auth/callback`. You can find the console URL in the status field of the `ROSAControlPlane` once the cluster is ready:
    ```shell
    kubectl get rosacontrolplane <control-plane-name> -o jsonpath='{.status.consoleURL}'
    ```

2. Create a new client secret in your OIDC provider and store the value in a kubernetes secret in the same namespace as your cluster:
    ```shell
    kubectl create secret generic console-client-secret --from-literal=clientSecret='<client-secret-value>' 
    ```

3. Configure `ROSAControlPlane` external auth provider with the created client:
    ```yaml
    ---
    apiVersion: controlplane.cluster.x-k8s.io/v1beta2
    kind: ROSAControlPlane
    metadata:
      name: "capi-rosa-quickstart-control-plane"
    spec:
      enableExternalAuthProviders: true
      externalAuthProviders:
      - name: my-oidc-provider
        issuer:
          issuerURL: https://login.microsoftonline.com/<tenant-id>/v2.0 # e.g. if using Microsoft Entra ID
          audiences:  # audiences that will be trusted by the kube-apiserver
          - "audience1"
          - <console-client-id> # <----New
        claimMappings:
          username:
            claim: email
            prefixPolicy: ""
          groups:
            claim: groups
        oidcClients:  # <----New
          - componentName: console
            componentNamespace: openshift-console
            clientID: <console-client-id>
            clientSecret:
              name: console-client-secret # secret name created in step 2
      ....
    ```

see [ROSAControlPlane CRD Reference](https://cluster-api-aws.sigs.k8s.io/crd/#controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider) for all possible configurations.

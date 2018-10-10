# CRDs workflow

0. Build your crd controller manager `make controller-manager-image`
1. Create a cluster with minikube (no longer has to be 1.9.4)
2. apply the config files:

    ```
    kubectl apply -f config/cluster_v1alpha1_cluster.yaml -f config/cluster_v1alpha1_machine.yaml -f rbac.yaml -f manager.yaml
    ```

3. Watch the logs with `kubectl logs -f controller-manager-0 -n system`
4. in a separate tab run `kubectl apply -f config/cluster.yaml`


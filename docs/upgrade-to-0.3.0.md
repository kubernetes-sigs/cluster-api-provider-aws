# Upgrading to 0.3.0

In 0.3.0, the tagging scheme changed for identifying AWS resources. In order not to lose track, there is a partial migration tool included in 0.3.0's `clusterawsadm`.

The migration path is as follows:

1. `kubectl scale statefulset -n aws-provider-system aws-provider-controller-manager --replicas=0`
2. `clusterawsadm migrate -n CLUSTER_NAME 0.3.0`
3. Update the image for the aws-provider-controller-manager
4. `kubectl scale statefulset -n aws-provider-system aws-provider-controller-manager --replicas=1`
5. Wait ~2 minutes for the security group changes to all settle
   - All of the nodes and control plane machines should have exactly one security group tagged with `kubernetes.io/cluster/<CLUSTER_NAME>=owned`: the new `CLUSTER_NAME-lb` group.
6. Find the names of your controller-manager pods, and run `kubectl exec -n kube-system -it CONTROLLER_MANAGER_POD_NAME -- sh -c 'kill 1' as a workaround for kubernetes/kubernetes#77019`

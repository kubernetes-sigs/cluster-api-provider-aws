# Create issue for ROSA

When creating issue for ROSA HCP cluster, include the logs for the capa-controller-manager and capi-controller-manager deployment pods. The logs can be saved to text file using the commands below. Also include the yaml files for all the resources used to create the ROSA HCP cluster:
- `Cluster`
- `ROSAControlPlane`
- `MachinePool`
- `ROSAMachinePool`

```shell
$ kubectl get pod -n capa-system 
NAME                                      READY   STATUS    RESTARTS   AGE
capa-controller-manager-77f5b946b-sddcg   1/1     Running   1          3d3h

$ kubectl logs -n capa-system capa-controller-manager-77f5b946b-sddcg > capa-controller-manager-logs.txt

$ kubectl get pod -n capi-system 
NAME                                       READY   STATUS    RESTARTS   AGE
capi-controller-manager-78dc897784-f8gpn   1/1     Running   18         26d

$ kubectl logs -n capi-system capi-controller-manager-78dc897784-f8gpn > capi-controller-manager-logs.txt
```

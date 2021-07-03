# EKS Addons

[EKS Addons](https://aws.amazon.com/blogs/containers/introducing-amazon-eks-add-ons/) can be used with EKS clusters created using Cluster API Provider AWS.

Addons are supported in EKS clusters using Kubernetes v1.18 or greater. 

## Installing addons

To install an addon you need to declare them by specifying the name, version and optionally how conflicts should be resolved in the `AWSManagedControlPlane`. For example:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
metadata:
  name: "capi-managed-test-control-plane"
spec:
  region: "eu-west-2"
  sshKeyName: "capi-management"
  version: "v1.18.0"
  addons:
    - name: "vpc-cni"
      version: "v1.6.3-eksbuild.1"
      conflictResolution: "overwrite"
```

Additionally, there is a cluster [flavor](https://cluster-api.sigs.k8s.io/clusterctl/commands/config-cluster.html#flavors)
called [eks-managedmachinepool-vpccni](../../templates/cluster-template-eks-managedmachinepool-vpccni.yaml) that you can use with **clusterctl**:

```shell
clusterctl generate cluster my-cluster --kubernetes-version v1.18.0 --flavor eks-managedmachinepool-vpccni > my-cluster.yaml
```

## Updating Addons

To update the version of an addon you need to edit the `AWSManagedControlPlane` instance and update the version of the addon you want to update. Using the example from the previous section we would do:

```yaml
...
  addons:
    - name: "vpc-cni"
      version: "v1.7.5-eksbuild.1"
      conflictResolution: "overwrite"
...
```

## Deleting Addons

To delete an addon from a cluster you need to edit the `AWSManagedControlPlane` instance and remove the entry for the addon you want to delete.

## Viewing installed addons

You can see what addons are installed on your EKS cluster by looking in the `Status`  of the `AWSManagedControlPlane` instance. 

Additionally you can run the following command:

```bash
clusterawsadm eks addons list-installed -n <<eksclustername>>
```

## Viewing available addons

You can see what addons are available to your EKS cluster by running the following command:

```bash
clusterawsadm eks addons list-available -n <<eksclustername>>
```

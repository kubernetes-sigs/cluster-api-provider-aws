# clusterawsadm eks
The `clusterawsadm eks` command provides EKS related utilities.


## clusterawsadm eks addons
The `clusterawsadm eks addons` command provides visibility to EKS addons.

```bash
// Lists available addons for the specified cluster
$ clusterawsadm eks addons list-available --cluster-name=<name> --region=us-west-1

// Lists installed addons to the specified cluster
$ clusterawsadm eks addons list-installed --cluster-name=<name> --region=us-west-1
```
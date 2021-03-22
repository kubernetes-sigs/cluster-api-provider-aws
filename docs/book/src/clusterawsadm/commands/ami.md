# clusterawsadm ami
The `clusterawsadm ami` command provides ami related utilities.

## clusterawsadm ami list
The `clusterawsadm ami list` command lists pre-built AMIs by Kubernetes version, OS, or AWS region.

```bash
// If no specific flag is provided, list all AMIs from the latest three Kubernetes releases (starting from the last stable one) for all supported regions and OS types.
$ clusterawsadm ami list

// List AMIs that have Kubernetes version v1.19.5 for all supported regions and OS types.
$ clusterawsadm ami list --kubernetes-version=v1.19.5

// List only ubuntu-18.04 AMI that has Kubernetes version v1.19.5 and is in us-west-1
$ clusterawsadm ami list --kubernetes-version=v1.19.5 --os=ubuntu-18.04 --region=us-west-1

```
To see the supported OS distributions and regions, see [amis](../../amis.md).

## clusterawsadm ami copy
The `clusterawsadm ami copy` command copy an AMI from an AWS account to the user's AWS account.
It first performs a lookup, after finding the AMI, copies the image to the user account.

```bash
// Finds the specified AMI in the default AWS account, then copies it to the AWS account that is retrieved from the local credentials.
$ clusterawsadm ami copy --os centos-7 --kubernetes-version=1.19.4
```

All AMIs are located in an AWS account that is hardcoded in the code. If the AMI to be copied is in the default account, `--source-account` flag should not be provided.
```bash
// Finds the specified AMI in the AWS source account# 111111111111, then copies it to the AWS account that is retrieved from the local credentials.
$ clusterawsadm ami copy --os centos-7 --kubernetes-version=1.19.4 --source-account=111111111111
```

With `--dry-run` option, you can dry-run the move action by only printing logs without taking any actual actions. Use log level verbosity `-v` to see different levels of information.


## clusterawsadm ami encrypted-copy
The `clusterawsadm ami encrypted-copy` command encrypts and copies an AMI from an AWS account to the user's AWS account.
If first encrypts and copies the AMI snapshot, then create an AMI with that snapshot.

This is especially useful for air-gapped environments, to avoid pulling the images from the Internet during cluster setup.
```bash
$ clusterawsadm ami encrypted-copy --kubernetes-version=v1.20.3 --os=ubuntu-20.04    
```
With `--dry-run` option, you can dry-run the move action by only printing logs without taking any actual actions. Use log level verbosity `-v` to see different levels of information.

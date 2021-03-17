# Overview of clusterawsadm

The `clusterawsadm` command line utility simplifies "day 0 operations" by managing identity and access management (IAM) objects necessary for this project 
and provide various AMI and EKS related commands. 

`clusterawsadm` binaries for both Darwin and Linux are published with each release, or it can be built from the source, see [Setting up Development Environment](./../development/development.md).

* use [`clusterawsadm bootstrap`](commands/bootstrap.md) to bootstrap IAM roles to be used by clusters and Cluster API Provider AWS
* use [`clusterawsadm ami`](commands/ami.md) for AMI related commands
* use [`clusterawsadm eks`](commands/eks.md) for commands related to EKS

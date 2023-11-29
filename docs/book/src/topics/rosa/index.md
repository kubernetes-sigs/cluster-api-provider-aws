# ROSA Support in the AWS Provider

- **Feature status:** Experimental
- **Feature gate (required):** ROSA=true

## Overview

The AWS provider supports creating Red Hat OpenShift Service on AWS ([ROSA](https://www.redhat.com/en/technologies/cloud-computing/openshift/aws)) based cluster. Currently the following features are supported:

- Provisioning/Deleting a ROSA cluster with hosted control planes ([HCP](https://docs.openshift.com/rosa/rosa_hcp/rosa-hcp-sts-creating-a-cluster-quickly.html))

The implementation introduces the following CRD kinds:

- ROSAControlPlane - specifies the ROSA Cluster in AWS
- ROSACluster - needed only to statisfy cluster-api contract

A new template is available in the templates folder for creating a managed ROSA workload cluster.

## SEE ALSO

* [Enabling ROSA Support](enabling.md)
* [Creating a cluster](creating-a-cluster.md)
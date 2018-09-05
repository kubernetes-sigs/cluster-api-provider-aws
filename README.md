# Kubernetes cluster-api-provider-aws Project

This repository hosts an implementation of a provider for AWS for the [cluster-api project](https://sigs.k8s.io/cluster-api).

Note: This repository is currently a skeleton implementation of a cluster-api provider, implementation will begin once there is agreement on the [Design Spec](https://docs.google.com/document/d/1G7DRQccoTY5YBrinQb6sz_fRLB9zFbCnI1O984XFk7Q).

## Development

### Requirements

* a google cloud project

### Set up

1. [Install the gcloud cli](https://cloud.google.com/sdk/install)
0. Set project: `gcloud config set project YOUR_PROJECT_NAME`
0. Pushing dev images: `make dev_push`

### clusterctl

1. `export CLUSTER_CONTROLLER_IMAGE=gcr.io/YOUR_PROJECT_NAME/aws-cluster-controller:0.0.1-dev`
0. `export MACHINE_CONTROLLER_IMAGE=gcr.io/YOUR_PROJECT_NAME/aws-machine-controller:0.0.1-dev`
0. Generate the input files with [clusterctl/examples/aws/generate-yaml.sh](/clusterctl/examples/aws/generate-yaml.sh)
0. `clusterctl create cluster -c clusterctl/examples/aws/out/cluster.yaml -m clusterctl/examples/aws/out/machines.yaml -p clusterctl/examples/aws/out/provider-components.yaml --provider aws`

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [#cluster-api on Kubernetes Slack](http://slack.k8s.io/messages/cluster-api)
- [SIG-Cluster-Lifecycle Mailing List](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).


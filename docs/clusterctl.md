# clusterctl

## Overview

`clusterctl` simplifies getting started with cluster-api-provider-aws. It wraps
the complicated steps into one command that will bring you from nothing to a
Kubernetes cluster in AWS that is running the cluster-API components.

## Install

### From release

Go to the
[releases](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases)
page and pick a version. Each release contains several binaries one of which is
`clusterctl`. Download the binary that was built for your architecture (either
darwin (os x) or linux). Windows binaries are not supported at this time. Add
the binary to your `$PATH` environment variable.

### From source

Clone the repo and run `make clusterctl`. Ensure the generated binary can be
found in the `$PATH` environment variable. The binary is built with bazel and
its location is `bazel-bin/cmd/clusterctl/darwin_amd64_pure_stripped/clusterctl`
on a darwin system.

## In depth

### `clusterctl create`

#### Bootstrap cluster

The result of a successful `create` command is a cluster running in AWS with the
cluster-api-provider-aws controllers watching Machine and Cluster objects ready
to create as many clusters as you need.

Cluster API controllers require kubernetes as a platform to
run<sup>[1](#footnote1)</sup>. The first step to `clusterctl` is to create a
bootstrap cluster or use an existing one. [KIND](https://sigs.k8s.io/kind) is
the suggested way to run an ephemeral cluster.
However, you can provide a kubeconfig to `clusterctl` if you already have a
cluster you'd like to use as a bootstrap cluster. Please see `clusterctl`'s help
flag `--help`. Note: The bootstrap cluster *must* be >= v1.12.0 of Kubernetes.

<a name="footnote1">1</a>: The controllers don't absolutely require a Kubernetes
cluster to run on, but as they were designed to run on Kubernetes it is simply
easier to create a cluster. There are plans to remove this dependency. Please
see [this Request For
Enhancement](https://github.com/kubernetes-sigs/cluster-api/issues/557).

#### Cluster API provider components

After the bootstrap cluster becomes active, `clusterctl` will apply the provider
components. This will get the cluster-api-provider-aws' controllers running on
the bootstrap cluster. These controllers watch the Machine and Cluster
[CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
and react by interacting with the AWS API to create a Kubernetes cluster.

#### Target cluster creation

Once `clusterctl` sees that the provider components are running, it will apply
the generated cluster YAML.

Once the cluster YAML is sent to the bootstrap API server, the
cluster-api-provider-aws controllers will talk to the AWS API and create
infrastructure necessary to run a Kubernetes cluster in AWS. At this point,
`clusterctl` will apply just the control plane machine node and wait for that to
become ready.

#### Addons

If specified, the addons will now be applied to the target cluster.

#### Pivot

At this point, there are two functional clusters: the bootstrap cluster (often
running locally) and a target cluster running in AWS that contains only the
control plane node(s). `clusterctl` will now move the provider components and
the Cluster and Machine objects from the bootstrap cluster to the target
cluster. This action is called the pivot.

#### The worker nodes

After the pivot has successfully completed, `clusterctl` will apply the worker
machines to the target cluster. These worker machines know how to join the
existing control plane to bring up a fully working kubernetes cluster.

#### Cleanup

`clusterctl` will cleanup the bootstrap cluster as it is no longer needed and
additional cluster can be created on the target cluster using the kubeconfig
that was written to local disk.

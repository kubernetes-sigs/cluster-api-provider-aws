# Kubernetes Cluster API Provider AWS

<p align="center">
<img src="https://github.com/kubernetes/kubernetes/raw/master/logo/logo.png"  width="100x"><a href="https://aws.amazon.com/opensource/"><img width="192x" src="https://d0.awsstatic.com/logos/powered-by-aws.png" alt="Powered by AWS Cloud Computing"></a>
</p>
<p align="center">
<!-- go doc / reference card -->
<a href="https://godoc.org/sigs.k8s.io/cluster-api-provider-aws">
 <img src="https://godoc.org/sigs.k8s.io/cluster-api-provider-aws?status.svg"></a>
</p>

------

Kubernetes-native declarative infrastructure for AWS.

## What is the Cluster API Provider AWS

The [Cluster API][cluster_api] brings
declarative, Kubernetes-style APIs to cluster creation, configuration and
management.

The API itself is shared across multiple cloud providers allowing for true AWS
hybrid deployments of Kubernetes. It is built atop the lessons learned from
previous cluster managers such as [kops][kops] and
[kubicorn][kubicorn].

## Documentation

Please see our [book](https://cluster-api-aws.sigs.k8s.io) for in-depth documentation.

## Launching a Kubernetes cluster on AWS

Check out the [Cluster API Quick Start](https://cluster-api.sigs.k8s.io/user/quick-start.html) for launching a
cluster on AWS.

## Features

- Native Kubernetes manifests and API
- Manages the bootstrapping of VPCs, gateways, security groups and instances.
- Choice of Linux distribution between Amazon Linux 2, CentOS 7 and Ubuntu 18.04,
  using [pre-baked AMIs](/docs/amis.md).
- Deploys Kubernetes control planes into private subnets with a separate
  bastion server.
- Doesn't use SSH for bootstrapping nodes.
- Installs only the minimal components to bootstrap a control plane and workers.
- Supports control planes on EC2 instances.
- Experimental [EKS support](docs/book/src/topics/eks/index.md)

------

## Compatibility with Cluster API and Kubernetes Versions

This provider's versions are compatible with the following versions of Cluster API:

|                              | Cluster API v1alpha3 (v0.3) | Cluster API v1alpha4 (v0.4) |
| ---------------------------- | --------------------------- | --------------------------- |
| AWS Provider v1alpha3 (v0.5) | ✓                           |                             |
| AWS Provider v1alpha3 (v0.6) | ✓                           |                             |
| AWS Provider v1alpha4 (v0.7) |                             | ✓                           |


This provider's versions are able to install and manage the following versions of Kubernetes:

|                              | Kubernetes 1.16 | Kubernetes 1.17 | Kubernetes 1.18 | Kubernetes 1.19 | Kubernetes 1.20 | Kubernetes 1.21 |
| ---------------------------- | --------------- | --------------- | --------------- | --------------- | --------------- | --------------- |
| AWS Provider v1alpha3 (v0.5) | ✓               | ✓               | ✓               | ✓               | ✓               |                 |
| AWS Provider v1alpha3 (v0.6) | ✓               | ✓               | ✓               | ✓               | ✓               | ✓               |
| AWS Provider v1alpha4 (v0.7) |                 |                 |                 | ✓               | ✓               | ✓               |


Each version of Cluster API for AWS will attempt to support two Kubernetes versions; e.g., Cluster API for AWS `v0.3` supports Kubernetes 1.16, 1.17, 1.18 etc.

**NOTE:** As the versioning for this project is tied to the versioning of Cluster API, future modifications to this
policy may be made to more closely align with other providers in the Cluster API ecosystem.

------

## Kubernetes versions with published AMIs

See [amis] for the list of most recently published AMIs.

------

## clusterawsadm

`clusterawsadm` CLI tool provides bootstrapping, AMI, EKS, and controller related helpers.

`clusterawsadm` binaries are released with each release, can be found under [assets](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/latest) section.

------
## Getting involved and contributing

Are you interested in contributing to cluster-api-provider-aws? We, the
maintainers and community, would love your suggestions, contributions, and help!
Also, the maintainers can be contacted at any time to learn more about how to get
involved.

In the interest of getting more new people involved we tag issues with
[`good first issue`][good_first_issue].
These are typically issues that have smaller scope but are good ways to start
to get acquainted with the codebase.

We also encourage ALL active community participants to act as if they are
maintainers, even if you don't have "official" write permissions. This is a
community effort, we are here to serve the Kubernetes community. If you have an
active interest and you want to get involved, you have real power! Don't assume
that the only people who can get things done around here are the "maintainers".

We also would love to add more "official" maintainers, so show us what you can
do!

This repository uses the Kubernetes bots.  See a full list of the commands [here][prow].

### Build the images locally

If you want to just build the CAPA containers locally, run

```
	REGISTRY=docker.io/my-reg make docker-build
```

### Tilt-based development environment

See [development][development] section for details

[development]: https://cluster-api-aws.sigs.k8s.io/development/development.html


### Implementer office hours

Maintainers hold office hours every two weeks, with sessions open to all
developers working on this project.

Office hours are hosted on a zoom video chat every other Monday
at 10:00 (Pacific) / 13:00 (Eastern) / 18:00 (Europe/London),
and are published on the [Kubernetes community meetings calendar][gcal].

### Other ways to communicate with the contributors

Please check in with us in the [#cluster-api-aws][slack] channel on Slack.

## Github issues

### Bugs

If you think you have found a bug please follow the instructions below.

- Please spend a small amount of time giving due diligence to the issue tracker. Your issue might be a duplicate.
- Get the logs from the cluster controllers. Please paste this into your issue.
- Open a [new issue][new_issue].
- Remember that users might be searching for your issue in the future, so please give it a meaningful title to help others.
- Feel free to reach out to the cluster-api community on the [kubernetes slack][slack].

### Tracking new features

We also use the issue tracker to track features. If you have an idea for a feature, or think you can help kops become even more awesome follow the steps below.

- Open a [new issue][new_issue].
- Remember that users might be searching for your issue in the future, so please
  give it a meaningful title to help others.
- Clearly define the use case, using concrete examples. EG: I type `this` and
  cluster-api-provider-aws does `that`.
- Some of our larger features will require some design. If you would like to
  include a technical design for your feature please include it in the issue.
- After the new feature is well understood, and the design agreed upon, we can
  start coding the feature. We would love for you to code it. So please open
  up a **WIP** *(work in progress)* pull request, and happy coding.

>“Amazon Web Services, AWS, and the “Powered by AWS” logo materials are
trademarks of Amazon.com, Inc. or its affiliates in the United States
and/or other countries."

<!-- References -->
[slack]: https://kubernetes.slack.com/messages/CD6U2V71N
[good_first_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22
[gcal]: https://calendar.google.com/calendar/embed?src=cgnt364vd8s86hr2phapfjc6uk%40group.calendar.google.com
[prow]: https://go.k8s.io/bot-commands
[new_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/new
[cluster_api]: https://github.com/kubernetes-sigs/cluster-api
[kops]: https://github.com/kubernetes/kops
[kubicorn]: http://kubicorn.io/
[tilt]: https://tilt.dev
[cluster_api_tilt]: https://master.cluster-api.sigs.k8s.io/developer/tilt.html
[amis]: https://cluster-api-aws.sigs.k8s.io/amis.html

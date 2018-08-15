# Cluster API
## What is the Cluster API?

The Cluster API is a Kubernetes project to bring declarative, Kubernetes-style
APIs to cluster creation, configuration, and management. It provides optional,
additive functionality on top of core Kubernetes.

Note that Cluster API effort is still in the prototype stage while we get
feedback on the API types themselves. All of the code here is to experiment with
the API and demo its abilities, in order to drive more technical feedback to the
API design. Because of this, all of the prototype code is rapidly changing.

![Cluster API Architecture](architecture.png "Cluster API Architecture")

To learn more, see the [Cluster API KEP][cluster-api-kep].

## Get involved!

* Join our Cluster API working group sessions
  * Weekly on Wednesdays @ 10:00 PT on [Zoom][zoomMeeting]
  * Previous meetings: \[ [notes][notes] | [recordings][recordings] \]

* Chat with us on [Slack](http://slack.k8s.io/): #cluster-api

* Pointers to repositories and PRs where some Cluster API provisioners are being
  developed.
  * AWS/Openshift, https://github.com/openshift/cluster-operator
  * Azure, https://github.com/platform9/azure-provider
  * GCE, https://github.com/kubernetes-sigs/cluster-api/tree/master/cloud/google
  * VSphere, https://github.com/kubernetes-sigs/cluster-api/tree/master/cloud/vsphere
  * OpenStack, https://github.com/kubernetes-sigs/cluster-api-provider-openstack

## Getting Started
### Prerequisites
* `kubectl` is required, see [here](http://kubernetes.io/docs/user-guide/prereqs/).
* `clusterctl` is a SIG-cluster-lifecycle sponsored tool to manage Cluster API clusters. See [here](clusterctl)
 
## How to use the API

To see how to build tooling on top of the Cluster API, please check out a few examples below:

* [upgrader](tools/upgrader/README.md): a cluster upgrade tool.
* [repair](tools/repair/README.md): detect problematic nodes and fix them.

[cluster-api-kep]: https://github.com/kubernetes/community/blob/master/keps/sig-cluster-lifecycle/0003-cluster-api.md
[notes]: https://docs.google.com/document/d/16ils69KImmE94RlmzjWDrkmFZysgB2J4lGnYMRN89WM/edit
[recordings]: https://www.youtube.com/playlist?list=PL69nYSiGNLP29D0nYgAGWt1ZFqS9Z7lw4
[zoomMeeting]: https://zoom.us/j/166836624

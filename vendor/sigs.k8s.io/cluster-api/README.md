# Cluster API
## What is the Cluster API?

The Cluster API is a Kubernetes project to bring declarative, Kubernetes-style
APIs to cluster creation, configuration, and management. It provides optional,
additive functionality on top of core Kubernetes.

Note that Cluster API effort is still in the prototype stage while we get
feedback on the API types themselves. All of the code here is to experiment with
the API and demo its abilities, in order to drive more technical feedback to the
API design. Because of this, all of the prototype code is rapidly changing.

## Getting Started

### Resources

* GitBook: [cluster-api.sigs.k8s.io](https://cluster-api.sigs.k8s.io)

### Prerequisites
* `kubectl` is required, see [here](http://kubernetes.io/docs/user-guide/prereqs/).
* `clusterctl` is a SIG-cluster-lifecycle sponsored tool to manage Cluster API clusters. See [here](cmd/clusterctl)

### Using `clusterctl` to create a cluster
* Doc [here](./docs/how-to-use-clusterctl.md)

![Cluster API Architecture](./docs/book/common_code/architecture.svg "Cluster API Architecture")

Learn more about the project's [scope, objectives, goals and requirements](./docs/scope-and-objectives.md), [feature proposals](./docs/proposals/) and [reference use cases](./docs/staging-use-cases.md).

### How does Cluster API compare to [Kubernetes Cloud Providers](https://kubernetes.io/docs/concepts/cluster-administration/cloud-providers/)?

Cloud Providers and the Cluster API work in concert to provide a rich Kubernetes experience in cloud environments.
The Cluster API initializes new nodes and clusters using available [providers](#Provider-Implementations).
Running clusters can then use Cloud Providers to provision support infrastructure like
[load balancers](https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/)
and [persistent volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/).

## Get involved!

* Join the [Cluster API discuss forum](https://discuss.kubernetes.io/c/contributors/cluster-api).

* Join the [sig-cluster-lifecycle](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle)
Google Group for access to documents and calendars.

* Join our Cluster API working group sessions
  * Weekly on Wednesdays @ 10:00 PT on [Zoom][zoomMeeting]
  * Previous meetings: \[ [notes][notes] | [recordings][recordings] \]

* Provider implementer office hours
  * Weekly on Tuesdays @ 12:00 PT ([Zoom][providerZoomMeetingTues]) and Wednesdays @ 15:00 CET ([Zoom][providerZoomMeetingWed])
  * Previous meetings: \[ [notes][implementerNotes] \]

* Chat with us on [Slack](http://slack.k8s.io/): #cluster-api

## Provider Implementations

The code in this repository is independent of any specific deployment environment.
Provider specific code is being developed in separate repositories, some of which
are also sponsored by SIG-cluster-lifecycle:

  * AWS, https://github.com/kubernetes-sigs/cluster-api-provider-aws
  * Azure, https://github.com/kubernetes-sigs/cluster-api-provider-azure
  * Baidu Cloud, https://github.com/baidu/cluster-api-provider-baiducloud
  * Bare Metal, https://github.com/metal3-io/cluster-api-provider-baremetal
  * DigitalOcean, https://github.com/kubernetes-sigs/cluster-api-provider-digitalocean
  * Exoscale, https://github.com/exoscale/cluster-api-provider-exoscale
  * GCP, https://github.com/kubernetes-sigs/cluster-api-provider-gcp
  * IBM Cloud, https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud
  * OpenStack, https://github.com/kubernetes-sigs/cluster-api-provider-openstack
  * Talos, https://github.com/talos-systems/cluster-api-provider-talos
  * Tencent Cloud, https://github.com/TencentCloud/cluster-api-provider-tencent
  * vSphere, https://github.com/kubernetes-sigs/cluster-api-provider-vsphere

## API Adoption

Following are the implementations managed by third-parties adopting the standard cluster-api and/or machine-api being developed here.

  * Kubermatic machine-controller, https://github.com/kubermatic/machine-controller/tree/master
  * Machine API Operator, https://github.com/openshift/machine-api-operator/tree/master
  * Machine-controller-manager, https://github.com/gardener/machine-controller-manager/tree/cluster-api

## Versioning, Maintenance, and Compatibility

- We follow [Semantic Versioning (semver)](https://semver.org/).
- Cluster API release cadence is Kubernetes Release + 6 weeks.
- The cadence is subject to change if necessary, refer to the [Milestones](https://github.com/kubernetes-sigs/cluster-api/milestones) page for up-to-date information.
- The _master_ branch is where development happens, this might include breaking changes.
- The _release-X_ branches contain stable, backward compatible code. A new _release-X_ branch is created at every major (X) release.



[notes]: https://docs.google.com/document/d/1Ys-DOR5UsgbMEeciuG0HOgDQc8kZsaWIWJeKJ1-UfbY/edit
[recordings]: https://www.youtube.com/playlist?list=PL69nYSiGNLP29D0nYgAGWt1ZFqS9Z7lw4
[zoomMeeting]: https://zoom.us/j/861487554
[implementerNotes]: https://docs.google.com/document/d/1IZ2-AZhe4r3CYiJuttyciS7bGZTTx4iMppcA8_Pr3xE/edit
[providerZoomMeetingTues]: https://zoom.us/j/140808484
[providerZoomMeetingWed]: https://zoom.us/j/424743530

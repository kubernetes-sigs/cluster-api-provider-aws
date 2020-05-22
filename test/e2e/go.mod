module sigs.k8s.io/cluster-api-provider-aws/test/e2e

go 1.13

require (
	github.com/aws/aws-sdk-go v1.31.3
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/vmware-tanzu/sonobuoy v0.17.2
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.1-0.20190704100234-640d9f240853+incompatible
	k8s.io/utils v0.0.0-20200229041039-0a110f9eb7ab
	sigs.k8s.io/cluster-api v0.3.6
	sigs.k8s.io/cluster-api-provider-aws v0.4.8
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/cluster-api-provider-aws => ../../
)

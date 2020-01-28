module sigs.k8s.io/cluster-api-provider-aws/test/e2e

go 1.13

require (
	github.com/aws/aws-sdk-go v1.25.16
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.8.1
	github.com/vmware-tanzu/sonobuoy v0.16.4
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	k8s.io/api v0.0.0-20191121015604-11707872ac1c
	k8s.io/apimachinery v0.0.0-20191121015412-41065c7a8c2a
	k8s.io/client-go v11.0.1-0.20190704100234-640d9f240853+incompatible
	sigs.k8s.io/cluster-api v0.2.6-0.20200204220036-b2ab4c203c74
	sigs.k8s.io/cluster-api-provider-aws v0.4.8
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/cluster-api-provider-aws => ../../
)

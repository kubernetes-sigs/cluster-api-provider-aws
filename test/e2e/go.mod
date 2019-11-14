module bhargavmad/cluster-api-provider-aws/test/e2e

go 1.13

require (
	github.com/aws/aws-sdk-go v1.20.19
	github.com/google/uuid v1.1.1 // indirect
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.1
	github.com/pkg/errors v0.8.1
	github.com/vmware-tanzu/sonobuoy v0.16.2
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v11.0.1-0.20190704100234-640d9f240853+incompatible
	sigs.k8s.io/cluster-api v0.2.7
	sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm v0.1.5
	sigs.k8s.io/cluster-api-provider-aws v0.4.4
	sigs.k8s.io/controller-runtime v0.3.0
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	sigs.k8s.io/cluster-api-provider-aws => ../../
)

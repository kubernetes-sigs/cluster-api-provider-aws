module sigs.k8s.io/cluster-api-provider-aws

go 1.12

require (
	github.com/aws/aws-sdk-go v1.20.19
	github.com/awslabs/goformation v0.0.0-20190310235947-776555df5a6d
	github.com/go-logr/logr v0.1.0
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/golang/mock v1.2.0
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f // indirect
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/cluster-bootstrap v0.0.0-20190711112844-b7409fb13d1b // indirect
	k8s.io/klog v0.4.0
	k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf // indirect
	k8s.io/kubernetes v1.14.2
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a
	sigs.k8s.io/cluster-api v0.2.0
	sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm v0.1.0
	sigs.k8s.io/controller-runtime v0.2.2
	sigs.k8s.io/yaml v1.1.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190704095032-f4ca3d3bdf1d
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.2.0
)

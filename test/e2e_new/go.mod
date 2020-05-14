module sigs.k8s.io/cluster-api-provider-aws/test/e2e_new

go 1.13

require (
	github.com/aws/aws-sdk-go v1.28.2
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66
	sigs.k8s.io/cluster-api v0.3.5
	sigs.k8s.io/cluster-api-provider-aws v0.5.3
	sigs.k8s.io/controller-runtime v0.5.2
)

replace sigs.k8s.io/cluster-api-provider-aws => ../../
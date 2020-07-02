module sigs.k8s.io/cluster-api-provider-aws/test/e2e_new

go 1.13

require (
	github.com/aws/aws-sdk-go v1.31.6
	github.com/onsi/ginkgo v1.12.2
	github.com/onsi/gomega v1.10.1
	k8s.io/api v0.17.7
	k8s.io/apimachinery v0.17.7
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451
	sigs.k8s.io/cluster-api v0.3.7-alpha.0.0.20200629143729-ef2b61f7d491
	sigs.k8s.io/cluster-api-provider-aws v0.5.3
	sigs.k8s.io/controller-runtime v0.5.7
	sigs.k8s.io/yaml v1.2.0
)

replace sigs.k8s.io/cluster-api-provider-aws => ../../

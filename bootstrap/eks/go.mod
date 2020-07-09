module sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/cluster-api v0.3.6
	sigs.k8s.io/cluster-api-provider-aws v0.5.4
	sigs.k8s.io/controller-runtime v0.5.2
)

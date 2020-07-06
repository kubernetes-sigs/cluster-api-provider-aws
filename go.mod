module sigs.k8s.io/cluster-api-provider-aws

go 1.13

require (
	github.com/aws/aws-sdk-go v1.33.3
	github.com/awslabs/goformation/v4 v4.11.0
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.4.3
	github.com/onsi/ginkgo v1.12.2
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.0.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
	k8s.io/api v0.17.7
	k8s.io/apimachinery v0.17.7
	k8s.io/client-go v0.17.7
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451
	sigs.k8s.io/cluster-api v0.3.7-alpha.0.0.20200629143729-ef2b61f7d491
	sigs.k8s.io/controller-runtime v0.5.7
	sigs.k8s.io/yaml v1.2.0
)

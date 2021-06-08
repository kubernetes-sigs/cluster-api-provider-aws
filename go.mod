module sigs.k8s.io/cluster-api-provider-aws

go 1.13

require (
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/amazon-vpc-cni-k8s v1.7.10
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.57
	github.com/awslabs/goformation/v4 v4.19.4
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.4.4
	github.com/google/goexpect v0.0.0-20210430020637-ab937bf7fd6f
	github.com/google/goterm v0.0.0-20200907032337-555d40f16ae2 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/sergi/go-diff v1.2.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.17.17
	k8s.io/apiextensions-apiserver v0.17.17
	k8s.io/apimachinery v0.17.17
	k8s.io/cli-runtime v0.17.17
	k8s.io/client-go v0.17.17
	k8s.io/component-base v0.17.17
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.1.0
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b
	sigs.k8s.io/aws-iam-authenticator v0.5.3
	sigs.k8s.io/cluster-api v0.3.18
	sigs.k8s.io/controller-runtime v0.5.14
	sigs.k8s.io/yaml v1.2.0
)

module sigs.k8s.io/cluster-api-provider-aws

go 1.13

replace github.com/golang/mock => github.com/golang/mock v1.4.4

require (
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/amazon-vpc-cni-k8s v1.7.10
	github.com/aws/aws-lambda-go v1.25.0
	github.com/aws/aws-sdk-go v1.40.2
	github.com/awslabs/goformation/v4 v4.19.5
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.5.0
	github.com/google/goexpect v0.0.0-20210430020637-ab937bf7fd6f
	github.com/google/goterm v0.0.0-20200907032337-555d40f16ae2 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/sergi/go-diff v1.2.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/net v0.0.0-20210716203947-853a461950ff
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.17.17
	k8s.io/apiextensions-apiserver v0.17.17
	k8s.io/apimachinery v0.17.17
	k8s.io/cli-runtime v0.17.17
	k8s.io/client-go v0.17.17
	k8s.io/component-base v0.17.17
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.1.0
	k8s.io/utils v0.0.0-20210709001253-0e1f9d693477
	sigs.k8s.io/aws-iam-authenticator v0.5.3
	sigs.k8s.io/cluster-api v0.3.22
	sigs.k8s.io/controller-runtime v0.5.14
	sigs.k8s.io/yaml v1.2.0
)

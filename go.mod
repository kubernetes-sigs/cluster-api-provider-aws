module sigs.k8s.io/cluster-api-provider-aws

go 1.16

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.0.1-0.20211028151834-d72fd59c8483

// Versions to be held for v1beta1
// sigs.k8s.io/controller-runtime on v0.10.x
// k8s.io/* on v0.22.x
// github.com/go-logr/logr on v0.4.x
// k8s.io/klog/v2 on v2.10.x
require (
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/amazon-vpc-cni-k8s v1.9.3
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.40.56
	github.com/awslabs/goformation/v4 v4.19.5
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.4.0
	github.com/gofrs/flock v0.8.1
	github.com/golang/mock v1.6.0
	github.com/google/goexpect v0.0.0-20210430020637-ab937bf7fd6f
	github.com/google/gofuzz v1.2.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/sergi/go-diff v1.2.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20211005215030-d2e5035098b3
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/cli-runtime v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/component-base v0.22.2
	k8s.io/klog/v2 v2.10.0
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
	sigs.k8s.io/aws-iam-authenticator v0.5.3
	sigs.k8s.io/cluster-api v1.0.1-0.20211028151834-d72fd59c8483
	sigs.k8s.io/cluster-api/test v1.0.1-0.20211028151834-d72fd59c8483
	sigs.k8s.io/controller-runtime v0.10.3-0.20211011182302-43ea648ec318
	sigs.k8s.io/yaml v1.3.0
)

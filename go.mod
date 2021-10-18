module sigs.k8s.io/cluster-api-provider-aws

go 1.16

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.0.0

// Versions to be held for v1beta1
// sigs.k8s.io/controller-runtime on v0.10.x
// k8s.io/* on v0.22.x
// github.com/go-logr/logr on v0.4.x
// k8s.io/klog/v2 on v2.10.x
require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/amazon-vpc-cni-k8s v1.9.3
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.40.56
	github.com/awslabs/goformation/v4 v4.19.5
	github.com/blang/semver v3.5.1+incompatible
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/containerd v1.5.7 // indirect
	github.com/coredns/caddy v1.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/docker/docker v20.10.9+incompatible // indirect
	github.com/drone/envsubst/v2 v2.0.0-20210730161058-179042472c46 // indirect
	github.com/evanphx/json-patch/v5 v5.5.0 // indirect
	github.com/go-logr/logr v0.4.0
	github.com/gofrs/flock v0.8.1
	github.com/golang/mock v1.6.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/goexpect v0.0.0-20210430020637-ab937bf7fd6f
	github.com/google/gofuzz v1.2.0
	github.com/google/goterm v0.0.0-20200907032337-555d40f16ae2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.31.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/sergi/go-diff v1.2.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20211005215030-d2e5035098b3
	golang.org/x/sys v0.0.0-20211004093028-2c5d950f24ef // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211005153810-c76a74d43a8e // indirect
	google.golang.org/grpc v1.41.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/cli-runtime v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/component-base v0.22.2
	k8s.io/klog/v2 v2.10.0
	k8s.io/kube-openapi v0.0.0-20210929172449-94abcedd1aa4 // indirect
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
	sigs.k8s.io/aws-iam-authenticator v0.5.3
	sigs.k8s.io/cluster-api v1.0.0
	sigs.k8s.io/cluster-api/test v1.0.0
	sigs.k8s.io/controller-runtime v0.10.2
	sigs.k8s.io/yaml v1.3.0
)

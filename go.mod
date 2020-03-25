module sigs.k8s.io/cluster-api-provider-aws

go 1.12

require (
	github.com/aws/aws-sdk-go v1.15.66
	github.com/blang/semver v3.5.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.2.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.7.0

	// kube 1.16
	github.com/openshift/machine-api-operator v0.2.1-0.20200320004730-01120a0337a6
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20191217005138-9e5e9d854fcc
	sigs.k8s.io/controller-runtime v0.3.1-0.20191016212439-2df793d02076
	sigs.k8s.io/controller-tools v0.2.2-0.20190930215132-4752ed2de7d2

)

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.4.0

// kubernetes-1.16.1
replace k8s.io/api => k8s.io/api v0.0.0-20191003000013-35e20aa79eb8

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191003000419-f68efa97b39e

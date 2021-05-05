module sigs.k8s.io/cluster-api-provider-aws

go 1.15

require (
	github.com/aws/aws-sdk-go v1.38.23
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.4.0
	github.com/golang/mock v1.4.4
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/openshift/api v0.0.0-20210416115537-a60c0dc032fd
	github.com/openshift/machine-api-operator v0.2.1-0.20210504014029-a132ec00f7dd

	// kube 1.18
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	k8s.io/klog/v2 v2.8.0
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009
	sigs.k8s.io/controller-runtime v0.9.0-alpha.1.0.20210413130450-7ef2da0bc161
	sigs.k8s.io/controller-tools v0.3.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20201216171336-0b00fb8d96ac
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20201209184807-075372e2ed03
)

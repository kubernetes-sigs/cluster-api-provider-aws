module sigs.k8s.io/cluster-api-provider-aws/hack/tools

go 1.16

// Dependencies which should not be upgraded during lifetime of v1alpha4
// k8s.io packages should stay on v0.21.x series
// k8s.io/klog/v2 should stay on v2.10.x series
// sigs.k8s.io/controller-tools should stay on v0.6.x series
// sigs.k8s.io/cluster-api/hack/tools needs to use v1alpha4 commits

require (
	github.com/a8m/envsubst v1.2.0
	github.com/ahmetb/gen-crd-api-reference-docs v0.3.0
	github.com/docker/docker v20.10.8+incompatible // indirect
	github.com/golang/mock v1.6.0
	github.com/golangci/golangci-lint v1.42.1
	github.com/itchyny/gojq v0.12.5
	github.com/onsi/ginkgo v1.16.4
	k8s.io/apimachinery v0.21.4
	k8s.io/code-generator v0.21.4
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.10.0 // indirect
	k8s.io/release v0.11.0
	sigs.k8s.io/cluster-api/hack/tools v0.0.0-20210812230458-f6fd5ed7dc0f
	sigs.k8s.io/controller-tools v0.6.2
	sigs.k8s.io/kind v0.11.1
	sigs.k8s.io/kustomize/kustomize/v4 v4.4.0
	sigs.k8s.io/testing_frameworks v0.1.2
)

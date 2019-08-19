module sigs.k8s.io/cluster-api-provider-aws

go 1.12

require (
	github.com/aws/aws-sdk-go v1.20.19
	github.com/awslabs/goformation v0.0.0-20190310235947-776555df5a6d
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.2.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	k8s.io/api v0.0.0-20190711103429-37c3b8b1ca65
	k8s.io/apimachinery v0.0.0-20190711103026-7bf792636534
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/cluster-bootstrap v0.0.0-20190711112844-b7409fb13d1b
	k8s.io/code-generator v0.0.0-20190311093542-50b561225d70
	k8s.io/gengo v0.0.0-20190813173942-955ffa8fcfc9 // indirect
	k8s.io/klog v0.3.3
	k8s.io/kubernetes v1.14.2
	k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5
	sigs.k8s.io/cluster-api v0.0.0-20190819183132-38edf8b4497e
	sigs.k8s.io/controller-runtime v0.2.0-rc.0
	sigs.k8s.io/controller-tools v0.2.0-rc.0
	sigs.k8s.io/testing_frameworks v0.1.2-0.20190130140139-57f07443c2d4
	sigs.k8s.io/yaml v1.1.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190704095032-f4ca3d3bdf1d
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190704094733-8f6ac2502e51
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.0.0-20190819183132-38edf8b4497e
)

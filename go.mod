module sigs.k8s.io/cluster-api-provider-aws

go 1.13

require (
	github.com/aws/aws-sdk-go v1.28.2
	github.com/awslabs/goformation/v4 v4.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.3.1
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200229041039-0a110f9eb7ab
	sigs.k8s.io/cluster-api v0.3.3
	sigs.k8s.io/controller-runtime v0.5.2
)

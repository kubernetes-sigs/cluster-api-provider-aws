module sigs.k8s.io/cluster-api-provider-aws

go 1.13

require (
	github.com/aws/aws-sdk-go v1.33.3
	github.com/awslabs/goformation/v4 v4.11.0
	github.com/frictionlessdata/tableschema-go v1.1.3
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.4.3
	github.com/google/goexpect v0.0.0-20200703111054-623d5ca06f56
	github.com/google/goterm v0.0.0-20190703233501-fc88cf888a3f // indirect
	github.com/matryer/is v1.2.0 // indirect
	github.com/onsi/ginkgo v1.12.2
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sergi/go-diff v1.0.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	k8s.io/api v0.17.8
	k8s.io/apimachinery v0.17.8
	k8s.io/client-go v0.17.8
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19
	sigs.k8s.io/cluster-api v0.3.7
	sigs.k8s.io/controller-runtime v0.5.8
	sigs.k8s.io/yaml v1.2.0
)

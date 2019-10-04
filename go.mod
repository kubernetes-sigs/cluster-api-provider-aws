module sigs.k8s.io/cluster-api-provider-aws

go 1.12

require (
	github.com/aws/aws-sdk-go v1.15.66
	github.com/blang/semver v3.5.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-log/log v0.0.0-20181211034820-a514cf01a3eb // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20181024230925-c65c006176ff // indirect
	github.com/golang/mock v1.2.0

	// kube 1.16
	github.com/openshift/cluster-api v0.0.0-20191004145247-2f02e328bd96
	github.com/openshift/cluster-api-actuator-pkg v0.0.0-20190527090340-7628df78fb4c
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/prometheus/procfs v0.0.0-20190209105433-f8d8b3f739bd // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/klog v0.4.0
	sigs.k8s.io/controller-runtime v0.2.0-beta.1.0.20190918234459-801e12a50160
	sigs.k8s.io/controller-tools v0.2.2-0.20190930215132-4752ed2de7d2

)

replace sigs.k8s.io/controller-runtime => github.com/enxebre/controller-runtime v0.2.0-beta.1.0.20190930160522-58015f7fc885

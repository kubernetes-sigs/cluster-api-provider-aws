load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "actuator.go",
        "interfaces.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/cluster",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/awsproviderconfig/v1alpha1:go_default_library",
        "//pkg/cloud/aws/services:go_default_library",
        "//pkg/cloud/aws/services/certificates:go_default_library",
        "//pkg/cloud/aws/services/ec2:go_default_library",
        "//pkg/cloud/aws/services/elb:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/elb:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/controller/error:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["actuator_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//pkg/apis/awsproviderconfig/v1alpha1:go_default_library",
        "//pkg/cloud/aws/actuators/cluster/mock_clusteriface:go_default_library",
        "//pkg/cloud/aws/services:go_default_library",
        "//pkg/cloud/aws/services/mocks:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1:go_default_library",
    ],
)

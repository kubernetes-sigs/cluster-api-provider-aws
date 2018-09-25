load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "eips.go",
        "errors.go",
        "gateways.go",
        "instances.go",
        "natgateways.go",
        "network.go",
        "routetables.go",
        "service.go",
        "subnets.go",
        "utils.go",
        "vpc.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2",
    visibility = ["//visibility:public"],
    deps = [
        "//cloud/aws/providerconfig/v1alpha1:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "gateways_test.go",
        "instances_test.go",
        "natgateways_test.go",
        "routetables_test.go",
        "subnets_test.go",
        "vpc_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//cloud/aws/providerconfig/v1alpha1:go_default_library",
        "//cloud/aws/services/ec2/mock_ec2iface:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

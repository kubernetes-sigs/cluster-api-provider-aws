load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["mock.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
    ],
)

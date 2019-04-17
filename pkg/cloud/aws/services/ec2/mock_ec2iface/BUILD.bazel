load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@bazel_gomock//:gomock.bzl", "gomock")

gomock(
    name = "mocks",
    out = "ec2api_mock.go",
    interfaces = ["EC2API"],
    library = "//vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface:go_default_library",
    package = "mock_ec2iface",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["ec2api_mock.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
    ],
)

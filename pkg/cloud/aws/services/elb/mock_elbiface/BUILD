load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@bazel_gomock//:gomock.bzl", "gomock")

gomock(
    name = "mocks",
    out = "elbapi_mock.go",
    interfaces = ["ELBAPI"],
    library = "//vendor/github.com/aws/aws-sdk-go/service/elb/elbiface:go_default_library",
    package = "mock_elbiface",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["elbapi_mock.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/elb:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
    ],
)

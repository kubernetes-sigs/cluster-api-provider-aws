load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["interface.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/model/api/codegentest/service/restxmlservice/restxmlserviceiface",
    importpath = "github.com/aws/aws-sdk-go/private/model/api/codegentest/service/restxmlservice/restxmlserviceiface",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/model/api/codegentest/service/restxmlservice:go_default_library",
    ],
)

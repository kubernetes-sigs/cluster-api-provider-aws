load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "sort_keys.go",
        "util.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/util",
    importpath = "github.com/aws/aws-sdk-go/private/util",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil:go_default_library"],
)

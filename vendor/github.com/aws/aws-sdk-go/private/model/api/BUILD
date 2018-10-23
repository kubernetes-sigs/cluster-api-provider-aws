load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "list_of_shame.go",
        "shape_alias.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/model/api",
    importpath = "github.com/aws/aws-sdk-go/private/model/api",
    visibility = ["//visibility:public"],
)

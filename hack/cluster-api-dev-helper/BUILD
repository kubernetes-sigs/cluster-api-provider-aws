load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/hack/cluster-api-dev-helper",
    visibility = ["//visibility:private"],
    deps = ["//hack/cluster-api-dev-helper/cmd:go_default_library"],
)

go_binary(
    name = "cluster-api-dev-helper",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

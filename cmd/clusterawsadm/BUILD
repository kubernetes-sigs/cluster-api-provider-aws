load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm",
    visibility = ["//visibility:private"],
    deps = ["//cmd/clusterawsadm/cmd:go_default_library"],
)

go_binary(
    name = "clusterawsadm",
    embed = [":go_default_library"],
    pure = "on",
    visibility = ["//visibility:public"],
)

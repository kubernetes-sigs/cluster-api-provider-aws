load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/cluster-api/cmd/clusterctl",
    importpath = "sigs.k8s.io/cluster-api/cmd/clusterctl",
    visibility = ["//visibility:private"],
    deps = ["//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/cmd:go_default_library"],
)

go_binary(
    name = "clusterctl",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

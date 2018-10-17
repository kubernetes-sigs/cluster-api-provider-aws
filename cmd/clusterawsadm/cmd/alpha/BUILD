load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["alpha.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha",
    visibility = ["//visibility:public"],
    deps = [
        "//cmd/clusterawsadm/cmd/alpha/bootstrap:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
    ],
)

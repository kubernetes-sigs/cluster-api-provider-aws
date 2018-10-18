load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "clusters.go",
        "deployer.go",
        "machines.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/cluster-api/pkg/errors",
    importpath = "sigs.k8s.io/cluster-api/pkg/errors",
    visibility = ["//visibility:public"],
    deps = ["//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/common:go_default_library"],
)

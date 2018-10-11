load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["deployer.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/deployer",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/common:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

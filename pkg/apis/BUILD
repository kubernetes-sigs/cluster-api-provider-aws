load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["apis.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/apis",
    visibility = ["//visibility:public"],
    deps = ["//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library"],
)

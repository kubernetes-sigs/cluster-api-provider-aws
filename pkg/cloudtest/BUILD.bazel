load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["cloudtest.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloudtest",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/go-logr/logr:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
    ],
)

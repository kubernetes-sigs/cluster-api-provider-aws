load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["proxy.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/godoc/proxy",
    importpath = "golang.org/x/tools/godoc/proxy",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/tools/godoc/env:go_default_library"],
)

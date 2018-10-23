load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["tomb.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/gopkg.in/tomb.v1",
    importpath = "gopkg.in/tomb.v1",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["tomb_test.go"],
    embed = [":go_default_library"],
)

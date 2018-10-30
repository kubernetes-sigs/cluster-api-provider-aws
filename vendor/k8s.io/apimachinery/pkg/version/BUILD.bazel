load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "helpers.go",
        "types.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/apimachinery/pkg/version",
    importpath = "k8s.io/apimachinery/pkg/version",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["helpers_test.go"],
    embed = [":go_default_library"],
)

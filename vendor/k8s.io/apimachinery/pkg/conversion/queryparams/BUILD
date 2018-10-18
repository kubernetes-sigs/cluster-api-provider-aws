load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "convert.go",
        "doc.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/apimachinery/pkg/conversion/queryparams",
    importpath = "k8s.io/apimachinery/pkg/conversion/queryparams",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["convert_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
    ],
)

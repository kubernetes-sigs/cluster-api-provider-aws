load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["clock.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/apimachinery/pkg/util/clock",
    importpath = "k8s.io/apimachinery/pkg/util/clock",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["clock_test.go"],
    embed = [":go_default_library"],
)

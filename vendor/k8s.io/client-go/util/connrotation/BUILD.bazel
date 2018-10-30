load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["connrotation.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/client-go/util/connrotation",
    importpath = "k8s.io/client-go/util/connrotation",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["connrotation_test.go"],
    embed = [":go_default_library"],
)

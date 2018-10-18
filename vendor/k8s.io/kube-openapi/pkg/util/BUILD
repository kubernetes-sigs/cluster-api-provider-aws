load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "trie.go",
        "util.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/kube-openapi/pkg/util",
    importpath = "k8s.io/kube-openapi/pkg/util",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["util_test.go"],
    embed = [":go_default_library"],
)

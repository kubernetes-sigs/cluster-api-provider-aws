load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "fuzz.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/google/gofuzz",
    importpath = "github.com/google/gofuzz",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "example_test.go",
        "fuzz_test.go",
    ],
    embed = [":go_default_library"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["a.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/tests/testdata/src/a",
    importpath = "golang.org/x/tools/go/analysis/passes/tests/testdata/src/a",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "a_test.go",
        "ax_test.go",
    ],
    embed = [":go_default_library"],
)

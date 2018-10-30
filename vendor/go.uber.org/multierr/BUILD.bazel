load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["error.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/multierr",
    importpath = "go.uber.org/multierr",
    visibility = ["//visibility:public"],
    deps = ["//vendor/go.uber.org/atomic:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "benchmarks_test.go",
        "error_test.go",
        "example_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//vendor/github.com/stretchr/testify/require:go_default_library",
    ],
)

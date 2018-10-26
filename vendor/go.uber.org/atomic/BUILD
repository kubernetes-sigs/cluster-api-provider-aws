load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "atomic.go",
        "string.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/atomic",
    importpath = "go.uber.org/atomic",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "atomic_test.go",
        "example_test.go",
        "stress_test.go",
        "string_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//vendor/github.com/stretchr/testify/require:go_default_library",
    ],
)

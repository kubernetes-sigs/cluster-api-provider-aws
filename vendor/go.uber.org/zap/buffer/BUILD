load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "buffer.go",
        "pool.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/zap/buffer",
    importpath = "go.uber.org/zap/buffer",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "buffer_test.go",
        "pool_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["@com_github_stretchr_testify//assert:go_default_library"],
)

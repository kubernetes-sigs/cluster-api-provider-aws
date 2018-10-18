load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["zapgrpc.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/zap/zapgrpc",
    importpath = "go.uber.org/zap/zapgrpc",
    visibility = ["//visibility:public"],
    deps = ["//vendor/go.uber.org/zap:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["zapgrpc_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/go.uber.org/zap:go_default_library",
        "//vendor/go.uber.org/zap/zapcore:go_default_library",
        "//vendor/go.uber.org/zap/zaptest/observer:go_default_library",
        "@com_github_stretchr_testify//require:go_default_library",
    ],
)

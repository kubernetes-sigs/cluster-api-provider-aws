load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "logger.go",
        "testingt.go",
        "timeout.go",
        "writer.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/zap/zaptest",
    importpath = "go.uber.org/zap/zaptest",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/go.uber.org/zap:go_default_library",
        "//vendor/go.uber.org/zap/internal/ztest:go_default_library",
        "//vendor/go.uber.org/zap/zapcore:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "logger_test.go",
        "testingt_test.go",
        "timeout_test.go",
        "writer_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/go.uber.org/zap:go_default_library",
        "//vendor/go.uber.org/zap/internal/ztest:go_default_library",
        "//vendor/go.uber.org/zap/zapcore:go_default_library",
        "@com_github_stretchr_testify//assert:go_default_library",
    ],
)

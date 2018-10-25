load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "logged_entry.go",
        "observer.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/zap/zaptest/observer",
    importpath = "go.uber.org/zap/zaptest/observer",
    visibility = ["//visibility:public"],
    deps = ["//vendor/go.uber.org/zap/zapcore:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "logged_entry_test.go",
        "observer_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//vendor/github.com/stretchr/testify/require:go_default_library",
        "//vendor/go.uber.org/zap:go_default_library",
        "//vendor/go.uber.org/zap/zapcore:go_default_library",
    ],
)

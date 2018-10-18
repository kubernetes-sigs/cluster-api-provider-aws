load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "deleg.go",
        "log.go",
        "null.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log",
    importpath = "sigs.k8s.io/controller-runtime/pkg/runtime/log",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/go-logr/logr:go_default_library",
        "//vendor/github.com/go-logr/zapr:go_default_library",
        "//vendor/go.uber.org/zap:go_default_library",
        "//vendor/go.uber.org/zap/zapcore:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "log_suite_test.go",
        "log_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/go-logr/logr:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest/printer:go_default_library",
    ],
)

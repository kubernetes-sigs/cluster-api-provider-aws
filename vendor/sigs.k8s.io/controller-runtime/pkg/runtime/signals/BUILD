load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "signal.go",
        "signal_posix.go",
        "signal_windows.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/runtime/signals",
    importpath = "sigs.k8s.io/controller-runtime/pkg/runtime/signals",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "signal_test.go",
        "signals_suite_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
    ],
)

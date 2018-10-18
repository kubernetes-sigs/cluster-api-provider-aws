load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "apis.go",
        "cmd.go",
        "config.go",
        "controller.go",
        "dockerfile.go",
        "webhook.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-tools/pkg/scaffold/manager",
    importpath = "sigs.k8s.io/controller-tools/pkg/scaffold/manager",
    visibility = ["//visibility:public"],
    deps = ["//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "manager_suite_test.go",
        "manager_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/scaffoldtest:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["doc.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/testing_frameworks/integration/internal/integration_tests",
    importpath = "sigs.k8s.io/testing_frameworks/integration/internal/integration_tests",
    visibility = ["//vendor/sigs.k8s.io/testing_frameworks/integration:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "apiserver_integration_test.go",
        "etcd_integration_test.go",
        "integration_suite_test.go",
        "integration_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/sigs.k8s.io/testing_frameworks/integration:go_default_library",
    ],
)

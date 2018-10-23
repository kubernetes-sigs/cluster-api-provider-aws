load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "apiserver.go",
        "control_plane.go",
        "doc.go",
        "etcd.go",
        "kubectl.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/testing_frameworks/integration",
    importpath = "sigs.k8s.io/testing_frameworks/integration",
    visibility = ["//visibility:public"],
    deps = ["//vendor/sigs.k8s.io/testing_frameworks/integration/internal:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "integration_suite_test.go",
        "kubectl_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
    ],
)

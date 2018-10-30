load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "address_manager.go",
        "apiserver.go",
        "arguments.go",
        "bin_path_finder.go",
        "etcd.go",
        "process.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/testing_frameworks/integration/internal",
    importpath = "sigs.k8s.io/testing_frameworks/integration/internal",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/onsi/gomega/gbytes:go_default_library",
        "//vendor/github.com/onsi/gomega/gexec:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "address_manager_test.go",
        "apiserver_test.go",
        "arguments_test.go",
        "bin_path_finder_test.go",
        "etcd_test.go",
        "internal_suite_test.go",
        "process_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/github.com/onsi/gomega/gexec:go_default_library",
        "//vendor/github.com/onsi/gomega/ghttp:go_default_library",
        "//vendor/sigs.k8s.io/testing_frameworks/integration:go_default_library",
    ],
)

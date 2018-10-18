load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "add.go",
        "controller.go",
        "controllersuitetest.go",
        "controllertest.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-tools/pkg/scaffold/controller",
    importpath = "sigs.k8s.io/controller-tools/pkg/scaffold/controller",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/markbates/inflect:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/resource:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "controller_suite_test.go",
        "controller_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/resource:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/scaffoldtest:go_default_library",
    ],
)

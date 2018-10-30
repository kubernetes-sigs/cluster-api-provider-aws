load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "addtoscheme.go",
        "crd.go",
        "crd_sample.go",
        "doc.go",
        "group.go",
        "register.go",
        "resource.go",
        "role.go",
        "rolebinding.go",
        "types.go",
        "typestest.go",
        "version_suitetest.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-tools/pkg/scaffold/resource",
    importpath = "sigs.k8s.io/controller-tools/pkg/scaffold/resource",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/markbates/inflect:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "resource_suite_test.go",
        "resource_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/scaffoldtest:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "boilerplate.go",
        "gitignore.go",
        "gopkg.go",
        "kustomize.go",
        "kustomize_image_patch.go",
        "makefile.go",
        "project.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-tools/pkg/scaffold/project",
    importpath = "sigs.k8s.io/controller-tools/pkg/scaffold/project",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/gopkg.in/yaml.v2:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "project_suite_test.go",
        "project_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/input:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/scaffold/scaffoldtest:go_default_library",
    ],
)

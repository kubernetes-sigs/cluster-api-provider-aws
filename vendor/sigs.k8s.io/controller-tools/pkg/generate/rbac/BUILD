load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "manifests.go",
        "parser.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-tools/pkg/generate/rbac",
    importpath = "sigs.k8s.io/controller-tools/pkg/generate/rbac",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/ghodss/yaml:go_default_library",
        "//vendor/k8s.io/api/rbac/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["parser_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/k8s.io/api/rbac/v1:go_default_library"],
)

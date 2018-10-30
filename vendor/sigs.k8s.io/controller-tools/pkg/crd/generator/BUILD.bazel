load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["generator.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-tools/pkg/crd/generator",
    importpath = "sigs.k8s.io/controller-tools/pkg/crd/generator",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/ghodss/yaml:go_default_library",
        "//vendor/github.com/spf13/afero:go_default_library",
        "//vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1:go_default_library",
        "//vendor/k8s.io/gengo/args:go_default_library",
        "//vendor/k8s.io/gengo/types:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/crd/util:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/internal/codegen:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/internal/codegen/parse:go_default_library",
        "//vendor/sigs.k8s.io/controller-tools/pkg/util:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["generator_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/spf13/afero:go_default_library"],
)

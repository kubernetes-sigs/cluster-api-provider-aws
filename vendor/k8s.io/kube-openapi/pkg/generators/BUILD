load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "api_linter.go",
        "extension.go",
        "openapi.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/kube-openapi/pkg/generators",
    importpath = "k8s.io/kube-openapi/pkg/generators",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/gengo/args:go_default_library",
        "//vendor/k8s.io/gengo/examples/set-gen/sets:go_default_library",
        "//vendor/k8s.io/gengo/generator:go_default_library",
        "//vendor/k8s.io/gengo/namer:go_default_library",
        "//vendor/k8s.io/gengo/types:go_default_library",
        "//vendor/k8s.io/kube-openapi/cmd/openapi-gen/args:go_default_library",
        "//vendor/k8s.io/kube-openapi/pkg/common:go_default_library",
        "//vendor/k8s.io/kube-openapi/pkg/generators/rules:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "extension_test.go",
        "openapi_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//vendor/k8s.io/gengo/examples/set-gen/sets:go_default_library",
        "//vendor/k8s.io/gengo/generator:go_default_library",
        "//vendor/k8s.io/gengo/namer:go_default_library",
        "//vendor/k8s.io/gengo/parser:go_default_library",
        "//vendor/k8s.io/gengo/types:go_default_library",
    ],
)

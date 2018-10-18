load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/apiextensions-apiserver",
    importpath = "k8s.io/apiextensions-apiserver",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/apiextensions-apiserver/pkg/cmd/server:go_default_library",
        "@io_k8s_apiserver//pkg/server:go_default_library",
        "@io_k8s_apiserver//pkg/util/logs:go_default_library",
    ],
)

go_binary(
    name = "apiextensions-apiserver",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

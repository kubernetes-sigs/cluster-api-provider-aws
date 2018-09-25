load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/cluster-controller",
    visibility = ["//visibility:private"],
    deps = [
        "//cloud/aws/controllers/cluster:go_default_library",
        "//cloud/aws/controllers/cluster/options:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/util/logs:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/controller/config:go_default_library",
    ],
)

go_binary(
    name = "cluster-controller",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("//build:rules.bzl", "cluster_api_binary")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/clusterctl",
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/apis/awsprovider/v1alpha1:go_default_library",
        "//pkg/cloud/aws/actuators/cluster:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/cmd:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/common:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client/config:go_default_library",
    ],
)

go_binary(
    name = "clusterctl",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

cluster_api_binary(name = "clusterctl")

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("//build:cluster_api_binary_image.bzl", "cluster_api_binary_image")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/manager",
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/apis:go_default_library",
        "//pkg/cloud/aws/actuators/cluster:go_default_library",
        "//pkg/cloud/aws/actuators/machine:go_default_library",
        "//pkg/record:go_default_library",
        "//vendor/k8s.io/klog:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/common:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/controller/cluster:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/controller/machine:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client/config:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/manager:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/signals:go_default_library",
    ],
)

go_binary(
    name = "manager",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
    pure = "on",
)

cluster_api_binary_image(name = "manager")

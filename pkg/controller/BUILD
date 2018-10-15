load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "add_cluster_controller.go",
        "add_machine_controller.go",
        "controller.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/controller",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/cloud/aws/actuators/cluster:go_default_library",
        "//pkg/cloud/aws/actuators/machine:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/controller/cluster:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/controller/machine:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/manager:go_default_library",
    ],
)

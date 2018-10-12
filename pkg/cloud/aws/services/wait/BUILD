load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["wait.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/wait",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/cloud/aws/services/awserrors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
    ],
)

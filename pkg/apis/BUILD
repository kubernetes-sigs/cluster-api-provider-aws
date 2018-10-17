load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "addtoscheme_awsprovider_v1alpha1.go",
        "apis.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/apis",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/awsprovider/v1alpha1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
    ],
)

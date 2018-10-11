load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["mock.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/cluster/mock_machine_interface",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

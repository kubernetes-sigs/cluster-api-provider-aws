load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "ec2.go",
        "elb.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/mocks",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/awsprovider/v1alpha1:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

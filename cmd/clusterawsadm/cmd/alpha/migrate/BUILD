load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["migrate.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha/migrate",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/infrastructure/v1alpha2:go_default_library",
        "@com_github_aws_aws_sdk_go//aws:go_default_library",
        "@com_github_aws_aws_sdk_go//aws/session:go_default_library",
        "@com_github_aws_aws_sdk_go//service/resourcegroupstaggingapi:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

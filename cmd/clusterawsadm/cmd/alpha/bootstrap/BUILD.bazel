load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["bootstrap.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha/bootstrap",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/cloud/aws/services/cloudformation:go_default_library",
        "//pkg/cloud/aws/services/iam:go_default_library",
        "//pkg/cloud/aws/services/sts:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/defaults:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/cloudformation:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/sts:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
    ],
)

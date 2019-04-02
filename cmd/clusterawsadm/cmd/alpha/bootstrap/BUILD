load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["bootstrap.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha/bootstrap",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/cloud/aws/actuators/cluster:go_default_library",
        "//pkg/cloud/aws/services/cloudformation:go_default_library",
        "//pkg/cloud/aws/services/iam:go_default_library",
        "//pkg/cloud/aws/services/sts:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/defaults:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/cloudformation:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/sts:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
        "//vendor/github.com/spf13/viper:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/bootstrap:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/provider:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/cmd:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/phases:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/common:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/util:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["v2.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/signer/v2",
    importpath = "github.com/aws/aws-sdk-go/private/signer/v2",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["v2_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/awstesting:go_default_library",
    ],
)

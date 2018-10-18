load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "builder.go",
        "connect.go",
        "doc.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/service/rds/rdsutils",
    importpath = "github.com/aws/aws-sdk-go/service/rds/rdsutils",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/signer/v4:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "builder_test.go",
        "connect_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library"],
)

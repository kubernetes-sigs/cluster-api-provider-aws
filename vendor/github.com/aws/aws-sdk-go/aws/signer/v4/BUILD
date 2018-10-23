load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "header_rules.go",
        "options.go",
        "uri_path.go",
        "v4.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/aws/signer/v4",
    importpath = "github.com/aws/aws-sdk-go/aws/signer/v4",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/internal/sdkio:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/rest:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "functional_1_5_test.go",
        "functional_test.go",
        "header_rules_test.go",
        "v4_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/awstesting:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/awstesting/unit:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/s3:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "decode.go",
        "defaults.go",
        "doc.go",
        "endpoints.go",
        "v3model.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/aws/endpoints",
    importpath = "github.com/aws/aws-sdk-go/aws/endpoints",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "decode_test.go",
        "endpoints_test.go",
        "example_test.go",
        "v3model_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/s3:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/sqs:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "doc.go",
        "errors.go",
        "service.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/service/lambda",
    importpath = "github.com/aws/aws-sdk-go/service/lambda",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/awsutil:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/client:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/client/metadata:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/signer/v4:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/restjson:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["examples_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
    ],
)

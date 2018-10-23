load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "debug.go",
        "decode.go",
        "encode.go",
        "error.go",
        "header.go",
        "header_value.go",
        "message.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/protocol/eventstream",
    importpath = "github.com/aws/aws-sdk-go/private/protocol/eventstream",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/aws/aws-sdk-go/aws:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "decode_test.go",
        "encode_test.go",
        "header_test.go",
        "header_value_test.go",
        "shared_test.go",
    ],
    data = glob(["testdata/**"]),
    embed = [":go_default_library"],
)

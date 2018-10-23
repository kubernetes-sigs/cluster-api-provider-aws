load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "copy.go",
        "equal.go",
        "path_value.go",
        "prettify.go",
        "string_value.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/aws/awsutil",
    importpath = "github.com/aws/aws-sdk-go/aws/awsutil",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/jmespath/go-jmespath:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "copy_test.go",
        "equal_test.go",
        "path_value_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/aws/aws-sdk-go/aws:go_default_library"],
)

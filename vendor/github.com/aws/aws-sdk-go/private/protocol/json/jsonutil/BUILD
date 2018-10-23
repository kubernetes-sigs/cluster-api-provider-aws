load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "build.go",
        "unmarshal.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/protocol/json/jsonutil",
    importpath = "github.com/aws/aws-sdk-go/private/protocol/json/jsonutil",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["build_test.go"],
    embed = [":go_default_library"],
    deps = ["@com_github_stretchr_testify//assert:go_default_library"],
)

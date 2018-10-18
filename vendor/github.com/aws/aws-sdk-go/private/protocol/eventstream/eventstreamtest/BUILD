load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["testing.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/protocol/eventstream/eventstreamtest",
    importpath = "github.com/aws/aws-sdk-go/private/protocol/eventstream/eventstreamtest",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/awstesting/unit:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/eventstream:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/eventstream/eventstreamapi:go_default_library",
        "//vendor/golang.org/x/net/http2:go_default_library",
    ],
)

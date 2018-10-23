load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["path.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/internal/sdkuri",
    importpath = "github.com/aws/aws-sdk-go/internal/sdkuri",
    visibility = ["//vendor/github.com/aws/aws-sdk-go:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = ["path_test.go"],
    embed = [":go_default_library"],
)

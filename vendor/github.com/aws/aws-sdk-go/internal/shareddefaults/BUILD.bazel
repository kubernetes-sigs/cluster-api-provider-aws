load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "ecs_container.go",
        "shared_config.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/internal/shareddefaults",
    importpath = "github.com/aws/aws-sdk-go/internal/shareddefaults",
    visibility = ["//vendor/github.com/aws/aws-sdk-go:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "shared_config_other_test.go",
        "shared_config_windows_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/aws/aws-sdk-go/awstesting:go_default_library"],
)

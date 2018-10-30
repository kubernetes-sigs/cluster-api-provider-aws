load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/awslabs/goformation/example/go-to-yaml",
    importpath = "github.com/awslabs/goformation/example/go-to-yaml",
    visibility = ["//visibility:private"],
    deps = ["//vendor/github.com/awslabs/goformation/cloudformation:go_default_library"],
)

go_binary(
    name = "go-to-yaml",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

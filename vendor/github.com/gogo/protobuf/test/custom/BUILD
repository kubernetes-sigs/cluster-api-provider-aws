load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["custom.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/gogo/protobuf/test/custom",
    importpath = "github.com/gogo/protobuf/test/custom",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["custom_test.go"],
    embed = [":go_default_library"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["pkg.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/cmd/callgraph/testdata/src/pkg",
    importpath = "golang.org/x/tools/cmd/callgraph/testdata/src/pkg",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "pkg",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["pkg_test.go"],
    embed = [":go_default_library"],
)

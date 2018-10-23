load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["a.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/findcall/testdata/src/a",
    importpath = "golang.org/x/tools/go/analysis/passes/findcall/testdata/src/a",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "a",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

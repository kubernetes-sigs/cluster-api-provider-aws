load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["versioninfo.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/versioninfo",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/spf13/cobra:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["versioninfo_test.go"],
    embed = [":go_default_library"],
)

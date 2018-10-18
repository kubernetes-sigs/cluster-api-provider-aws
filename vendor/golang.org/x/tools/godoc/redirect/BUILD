load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "hash.go",
        "redirect.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/godoc/redirect",
    importpath = "golang.org/x/tools/godoc/redirect",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["redirect_test.go"],
    embed = [":go_default_library"],
)

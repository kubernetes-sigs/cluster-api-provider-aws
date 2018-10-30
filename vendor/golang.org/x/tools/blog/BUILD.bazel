load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["blog.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/blog",
    importpath = "golang.org/x/tools/blog",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/blog/atom:go_default_library",
        "//vendor/golang.org/x/tools/present:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["blog_test.go"],
    embed = [":go_default_library"],
)

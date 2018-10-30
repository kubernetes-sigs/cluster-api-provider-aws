load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["zipfs.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/godoc/vfs/zipfs",
    importpath = "golang.org/x/tools/godoc/vfs/zipfs",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/tools/godoc/vfs:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["zipfs_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/godoc/vfs:go_default_library"],
)

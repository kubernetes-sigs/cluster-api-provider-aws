load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "emptyvfs.go",
        "namespace.go",
        "os.go",
        "vfs.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/godoc/vfs",
    importpath = "golang.org/x/tools/godoc/vfs",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "emptyvfs_test.go",
        "os_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/godoc/vfs/mapfs:go_default_library"],
)

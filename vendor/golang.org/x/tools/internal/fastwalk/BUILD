load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "fastwalk.go",
        "fastwalk_dirent_fileno.go",
        "fastwalk_dirent_ino.go",
        "fastwalk_portable.go",
        "fastwalk_unix.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/internal/fastwalk",
    importpath = "golang.org/x/tools/internal/fastwalk",
    visibility = ["//vendor/golang.org/x/tools:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = ["fastwalk_test.go"],
    embed = [":go_default_library"],
)

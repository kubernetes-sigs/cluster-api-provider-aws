load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["walk.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/internal/gopathwalk",
    importpath = "golang.org/x/tools/internal/gopathwalk",
    visibility = ["//vendor/golang.org/x/tools:__subpackages__"],
    deps = ["//vendor/golang.org/x/tools/internal/fastwalk:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["walk_test.go"],
    embed = [":go_default_library"],
)

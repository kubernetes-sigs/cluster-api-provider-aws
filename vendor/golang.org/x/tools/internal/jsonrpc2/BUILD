load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "jsonrpc2.go",
        "log.go",
        "stream.go",
        "wire.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/internal/jsonrpc2",
    importpath = "golang.org/x/tools/internal/jsonrpc2",
    visibility = ["//vendor/golang.org/x/tools:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = ["jsonrpc2_test.go"],
    embed = [":go_default_library"],
)

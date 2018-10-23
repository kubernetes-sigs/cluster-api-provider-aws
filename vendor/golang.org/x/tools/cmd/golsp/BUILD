load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/cmd/golsp",
    importpath = "golang.org/x/tools/cmd/golsp",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/golang.org/x/tools/internal/jsonrpc2:go_default_library",
        "//vendor/golang.org/x/tools/internal/lsp:go_default_library",
    ],
)

go_binary(
    name = "golsp",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "diagnostics.go",
        "format.go",
        "server.go",
        "view.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/internal/lsp",
    importpath = "golang.org/x/tools/internal/lsp",
    visibility = ["//vendor/golang.org/x/tools:__subpackages__"],
    deps = [
        "//vendor/golang.org/x/tools/go/packages:go_default_library",
        "//vendor/golang.org/x/tools/internal/jsonrpc2:go_default_library",
        "//vendor/golang.org/x/tools/internal/lsp/protocol:go_default_library",
    ],
)

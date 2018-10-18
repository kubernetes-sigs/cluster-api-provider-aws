load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "basic.go",
        "client.go",
        "diagnostics.go",
        "doc.go",
        "general.go",
        "language.go",
        "protocol.go",
        "registration.go",
        "server.go",
        "text.go",
        "window.go",
        "workspace.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/internal/lsp/protocol",
    importpath = "golang.org/x/tools/internal/lsp/protocol",
    visibility = ["//vendor/golang.org/x/tools:__subpackages__"],
    deps = ["//vendor/golang.org/x/tools/internal/jsonrpc2:go_default_library"],
)

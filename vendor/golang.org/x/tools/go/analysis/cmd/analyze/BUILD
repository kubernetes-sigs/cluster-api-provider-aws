load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["analyze.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/cmd/analyze",
    importpath = "golang.org/x/tools/go/analysis/cmd/analyze",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/golang.org/x/tools/go/analysis/multichecker:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/buildtag:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/findcall:go_default_library",
    ],
)

go_binary(
    name = "analyze",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

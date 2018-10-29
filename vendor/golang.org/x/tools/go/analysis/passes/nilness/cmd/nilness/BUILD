load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness",
    importpath = "golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/golang.org/x/tools/go/analysis/passes/nilness:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/singlechecker:go_default_library",
    ],
)

go_binary(
    name = "nilness",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

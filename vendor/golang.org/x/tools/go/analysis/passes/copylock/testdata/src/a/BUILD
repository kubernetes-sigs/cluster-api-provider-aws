load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "copylock.go",
        "copylock_func.go",
        "copylock_range.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/copylock/testdata/src/a",
    importpath = "golang.org/x/tools/go/analysis/passes/copylock/testdata/src/a",
    visibility = ["//visibility:public"],
)

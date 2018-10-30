load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["util.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/internal/analysisutil",
    importpath = "golang.org/x/tools/go/analysis/passes/internal/analysisutil",
    visibility = ["//vendor/golang.org/x/tools/go/analysis/passes:__subpackages__"],
)

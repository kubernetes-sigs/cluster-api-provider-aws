load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "analysis.go",
        "doc.go",
        "validate.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis",
    importpath = "golang.org/x/tools/go/analysis",
    visibility = ["//visibility:public"],
)

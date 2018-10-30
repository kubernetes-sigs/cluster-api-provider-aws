load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["checker.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/internal/checker",
    importpath = "golang.org/x/tools/go/analysis/internal/checker",
    visibility = ["//vendor/golang.org/x/tools/go/analysis:__subpackages__"],
    deps = [
        "//vendor/golang.org/x/tools/go/analysis:go_default_library",
        "//vendor/golang.org/x/tools/go/packages:go_default_library",
    ],
)

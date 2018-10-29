load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["multichecker.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/multichecker",
    importpath = "golang.org/x/tools/go/analysis/multichecker",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/analysis:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/internal/analysisflags:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/internal/checker:go_default_library",
    ],
)

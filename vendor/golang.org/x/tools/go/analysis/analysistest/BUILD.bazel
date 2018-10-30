load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["analysistest.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/analysistest",
    importpath = "golang.org/x/tools/go/analysis/analysistest",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/analysis:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/internal/checker:go_default_library",
        "//vendor/golang.org/x/tools/go/packages:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["analysistest_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/go/analysis/passes/findcall:go_default_library"],
)

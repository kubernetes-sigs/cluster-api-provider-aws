load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["lostcancel.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/lostcancel",
    importpath = "golang.org/x/tools/go/analysis/passes/lostcancel",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/analysis:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/ctrlflow:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/inspect:go_default_library",
        "//vendor/golang.org/x/tools/go/ast/inspector:go_default_library",
        "//vendor/golang.org/x/tools/go/cfg:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["lostcancel_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/go/analysis/analysistest:go_default_library"],
)

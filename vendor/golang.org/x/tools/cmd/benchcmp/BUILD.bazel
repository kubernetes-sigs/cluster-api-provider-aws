load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "benchcmp.go",
        "compare.go",
        "doc.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/cmd/benchcmp",
    importpath = "golang.org/x/tools/cmd/benchcmp",
    visibility = ["//visibility:private"],
    deps = ["//vendor/golang.org/x/tools/benchmark/parse:go_default_library"],
)

go_binary(
    name = "benchcmp",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "benchcmp_test.go",
        "compare_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/benchmark/parse:go_default_library"],
)

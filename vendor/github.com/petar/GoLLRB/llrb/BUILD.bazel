load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "avgvar.go",
        "iterator.go",
        "llrb.go",
        "llrb-stats.go",
        "util.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/petar/GoLLRB/llrb",
    importpath = "github.com/petar/GoLLRB/llrb",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "iterator_test.go",
        "llrb_test.go",
    ],
    embed = [":go_default_library"],
)

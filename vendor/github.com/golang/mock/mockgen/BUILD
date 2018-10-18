load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "mockgen.go",
        "parse.go",
        "reflect.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen",
    importpath = "github.com/golang/mock/mockgen",
    visibility = ["//visibility:private"],
    deps = ["//vendor/github.com/golang/mock/mockgen/model:go_default_library"],
)

go_binary(
    name = "mockgen",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "mockgen_test.go",
        "parse_test.go",
    ],
    embed = [":go_default_library"],
)

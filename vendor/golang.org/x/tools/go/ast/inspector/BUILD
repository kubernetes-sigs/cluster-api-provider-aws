load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "inspector.go",
        "typeof.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/ast/inspector",
    importpath = "golang.org/x/tools/go/ast/inspector",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["inspector_test.go"],
    embed = [":go_default_library"],
)

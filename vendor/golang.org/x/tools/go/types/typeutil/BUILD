load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "callee.go",
        "imports.go",
        "map.go",
        "methodsetcache.go",
        "ui.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/types/typeutil",
    importpath = "golang.org/x/tools/go/types/typeutil",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/tools/go/ast/astutil:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "callee_test.go",
        "example_test.go",
        "imports_test.go",
        "map_test.go",
        "ui_test.go",
    ],
    embed = [":go_default_library"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "check.go",
        "mvpkg.go",
        "rename.go",
        "spec.go",
        "util.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/refactor/rename",
    importpath = "golang.org/x/tools/refactor/rename",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/ast/astutil:go_default_library",
        "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        "//vendor/golang.org/x/tools/go/loader:go_default_library",
        "//vendor/golang.org/x/tools/go/types/typeutil:go_default_library",
        "//vendor/golang.org/x/tools/refactor/importgraph:go_default_library",
        "//vendor/golang.org/x/tools/refactor/satisfy:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "mvpkg_test.go",
        "rename_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/go/buildutil:go_default_library"],
)

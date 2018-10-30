load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "fix.go",
        "imports.go",
        "sortimports.go",
        "zstdlib.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/imports",
    importpath = "golang.org/x/tools/imports",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/ast/astutil:go_default_library",
        "//vendor/golang.org/x/tools/internal/gopathwalk:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["fix_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/go/packages/packagestest:go_default_library"],
)

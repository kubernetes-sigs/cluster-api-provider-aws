load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "export.go",
        "gopath.go",
        "modules.go",
        "modules_111.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/packages/packagestest",
    importpath = "golang.org/x/tools/go/packages/packagestest",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/tools/go/packages:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "export_test.go",
        "gopath_test.go",
        "modules_test.go",
    ],
    embed = [":go_default_library"],
)

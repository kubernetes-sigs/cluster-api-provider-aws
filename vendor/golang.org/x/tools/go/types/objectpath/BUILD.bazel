load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["objectpath.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/types/objectpath",
    importpath = "golang.org/x/tools/go/types/objectpath",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["objectpath_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        "//vendor/golang.org/x/tools/go/gcexportdata:go_default_library",
        "//vendor/golang.org/x/tools/go/loader:go_default_library",
    ],
)

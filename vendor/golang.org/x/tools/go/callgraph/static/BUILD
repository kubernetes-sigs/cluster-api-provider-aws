load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["static.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/callgraph/static",
    importpath = "golang.org/x/tools/go/callgraph/static",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
        "//vendor/golang.org/x/tools/go/ssa:go_default_library",
        "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["static_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
        "//vendor/golang.org/x/tools/go/loader:go_default_library",
        "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
    ],
)

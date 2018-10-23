load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "external.go",
        "external_darwin.go",
        "external_unix.go",
        "interp.go",
        "map.go",
        "ops.go",
        "reflect.go",
        "value.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/ssa/interp",
    importpath = "golang.org/x/tools/go/ssa/interp",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/ssa:go_default_library",
        "//vendor/golang.org/x/tools/go/types/typeutil:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["interp_test.go"],
    embed = [":go_default_library"],
    deps = select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

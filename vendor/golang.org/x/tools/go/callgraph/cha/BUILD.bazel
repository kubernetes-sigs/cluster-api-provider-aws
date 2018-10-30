load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["cha.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/callgraph/cha",
    importpath = "golang.org/x/tools/go/callgraph/cha",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
        "//vendor/golang.org/x/tools/go/ssa:go_default_library",
        "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        "//vendor/golang.org/x/tools/go/types/typeutil:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["cha_test.go"],
    embed = [":go_default_library"],
    deps = select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:nacl": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "//vendor/golang.org/x/tools/go/callgraph:go_default_library",
            "//vendor/golang.org/x/tools/go/loader:go_default_library",
            "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

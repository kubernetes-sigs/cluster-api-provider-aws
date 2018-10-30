load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "blockopt.go",
        "builder.go",
        "const.go",
        "create.go",
        "doc.go",
        "dom.go",
        "emit.go",
        "func.go",
        "identical.go",
        "identical_17.go",
        "lift.go",
        "lvalue.go",
        "methods.go",
        "mode.go",
        "print.go",
        "sanity.go",
        "source.go",
        "ssa.go",
        "testmain.go",
        "util.go",
        "wrappers.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/ssa",
    importpath = "golang.org/x/tools/go/ssa",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/ast/astutil:go_default_library",
        "//vendor/golang.org/x/tools/go/types/typeutil:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "builder_test.go",
        "example_test.go",
        "identical_test.go",
        "source_test.go",
        "stdlib_test.go",
        "testmain_test.go",
    ],
    data = glob(["testdata/**"]),
    embed = [":go_default_library"],
    deps = [
        "//vendor/golang.org/x/tools/go/ast/astutil:go_default_library",
        "//vendor/golang.org/x/tools/go/loader:go_default_library",
        "//vendor/golang.org/x/tools/go/packages:go_default_library",
        "//vendor/golang.org/x/tools/go/ssa/ssautil:go_default_library",
    ] + select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:nacl": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

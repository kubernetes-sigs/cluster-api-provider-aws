load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["vet.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/cmd/vet",
    importpath = "golang.org/x/tools/go/analysis/cmd/vet",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/golang.org/x/tools/go/analysis/multichecker:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/asmdecl:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/assign:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/atomic:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/bools:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/buildtag:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/cgocall:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/composite:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/copylock:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/findcall:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/httpresponse:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/loopclosure:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/lostcancel:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/nilfunc:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/nilness:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/pkgfact:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/printf:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/shift:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/stdmethods:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/structtag:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/tests:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/unreachable:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/unsafeptr:go_default_library",
        "//vendor/golang.org/x/tools/go/analysis/passes/unusedresult:go_default_library",
    ],
)

go_binary(
    name = "vet",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

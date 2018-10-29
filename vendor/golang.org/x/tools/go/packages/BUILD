load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "external.go",
        "golist.go",
        "golist_fallback.go",
        "golist_fallback_testmain.go",
        "packages.go",
        "visit.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/packages",
    importpath = "golang.org/x/tools/go/packages",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/go/gcexportdata:go_default_library",
        "//vendor/golang.org/x/tools/go/internal/cgo:go_default_library",
        "//vendor/golang.org/x/tools/internal/gopathwalk:go_default_library",
        "//vendor/golang.org/x/tools/internal/semver:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "example_test.go",
        "packages110_test.go",
        "packages_test.go",
        "packagescgo_test.go",
        "stdlib_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/go/packages/packagestest:go_default_library"],
)

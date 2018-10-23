load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "cgo.go",
        "cgo2.go",
        "cgo3.go",
        "cgo4.go",
    ],
    cgo = True,
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/cgocall/testdata/src/a",
    importpath = "golang.org/x/tools/go/analysis/passes/cgocall/testdata/src/a",
    visibility = ["//visibility:public"],
)

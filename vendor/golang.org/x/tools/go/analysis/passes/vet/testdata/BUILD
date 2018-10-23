load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "httpresponse.go",
        "print.go",
        "shadow.go",
        "shift.go",
        "structtag.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/passes/vet/testdata",
    importpath = "golang.org/x/tools/go/analysis/passes/vet/testdata",
    visibility = ["//visibility:public"],
)

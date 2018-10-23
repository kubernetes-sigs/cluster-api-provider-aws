load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "a.go",
        "b.go",
        "c.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/cmd/bundle/testdata/src/initial",
    importpath = "golang.org/x/tools/cmd/bundle/testdata/src/initial",
    visibility = ["//visibility:public"],
)

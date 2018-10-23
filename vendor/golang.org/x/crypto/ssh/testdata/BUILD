load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "keys.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/crypto/ssh/testdata",
    importpath = "golang.org/x/crypto/ssh/testdata",
    visibility = ["//visibility:public"],
)

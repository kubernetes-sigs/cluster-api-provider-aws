load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["generate.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen/tests/internal_pkg",
    importpath = "github.com/golang/mock/mockgen/tests/internal_pkg",
    visibility = ["//visibility:public"],
)

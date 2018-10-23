load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["input.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen/tests/internal_pkg/subdir/internal/pkg",
    importpath = "github.com/golang/mock/mockgen/tests/internal_pkg/subdir/internal/pkg",
    visibility = ["//vendor/github.com/golang/mock/mockgen/tests/internal_pkg/subdir:__subpackages__"],
)

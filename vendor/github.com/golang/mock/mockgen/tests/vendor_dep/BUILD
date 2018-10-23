load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "mock.go",
        "vendor_dep.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen/tests/vendor_dep",
    importpath = "github.com/golang/mock/mockgen/tests/vendor_dep",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/golang/mock/gomock:go_default_library"],
)

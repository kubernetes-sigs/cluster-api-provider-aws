load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["source_mock.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen/tests/import_source",
    importpath = "github.com/golang/mock/mockgen/tests/import_source",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/golang/mock/mockgen/tests/import_source/definition:go_default_library",
    ],
)

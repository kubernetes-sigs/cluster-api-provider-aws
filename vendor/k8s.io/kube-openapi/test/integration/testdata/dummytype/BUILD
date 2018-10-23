load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "alpha.go",
        "beta.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/kube-openapi/test/integration/testdata/dummytype",
    importpath = "k8s.io/kube-openapi/test/integration/testdata/dummytype",
    visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["jmespath.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/jmespath/go-jmespath/fuzz",
    importpath = "github.com/jmespath/go-jmespath/fuzz",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/jmespath/go-jmespath:go_default_library"],
)

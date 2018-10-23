load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["ocsp.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/crypto/ocsp",
    importpath = "golang.org/x/crypto/ocsp",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["ocsp_test.go"],
    embed = [":go_default_library"],
)

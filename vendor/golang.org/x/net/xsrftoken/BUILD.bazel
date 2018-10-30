load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["xsrf.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/net/xsrftoken",
    importpath = "golang.org/x/net/xsrftoken",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["xsrf_test.go"],
    embed = [":go_default_library"],
)

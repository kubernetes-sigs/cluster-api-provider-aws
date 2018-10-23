load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["diskcache.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/gregjones/httpcache/diskcache",
    importpath = "github.com/gregjones/httpcache/diskcache",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/peterbourgon/diskv:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["diskcache_test.go"],
    embed = [":go_default_library"],
)

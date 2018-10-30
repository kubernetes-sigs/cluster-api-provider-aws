load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["leveldbcache.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/gregjones/httpcache/leveldbcache",
    importpath = "github.com/gregjones/httpcache/leveldbcache",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/syndtr/goleveldb/leveldb:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["leveldbcache_test.go"],
    embed = [":go_default_library"],
)

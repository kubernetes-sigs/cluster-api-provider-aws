load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "parse.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/gengo/parser",
    importpath = "k8s.io/gengo/parser",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/gengo/types:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "local_parse_test.go",
        "parse_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/k8s.io/gengo/args:go_default_library",
        "//vendor/k8s.io/gengo/namer:go_default_library",
        "//vendor/k8s.io/gengo/types:go_default_library",
    ],
)

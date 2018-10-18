load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["nopackage.pb.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/gogo/protobuf/test/nopackage",
    importpath = "github.com/gogo/protobuf/test/nopackage",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/gogo/protobuf/proto:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["nopackage_test.go"],
    embed = [":go_default_library"],
)

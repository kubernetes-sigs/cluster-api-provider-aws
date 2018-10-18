load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["descriptor.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/protobuf/descriptor",
    importpath = "github.com/golang/protobuf/descriptor",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/protobuf/proto:go_default_library",
        "//vendor/github.com/golang/protobuf/protoc-gen-go/descriptor:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["descriptor_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/golang/protobuf/proto/test_proto:go_default_library",
        "//vendor/github.com/golang/protobuf/protoc-gen-go/descriptor:go_default_library",
    ],
)

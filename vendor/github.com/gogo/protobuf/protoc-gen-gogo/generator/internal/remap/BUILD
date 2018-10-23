load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["remap.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator/internal/remap",
    importpath = "github.com/gogo/protobuf/protoc-gen-gogo/generator/internal/remap",
    visibility = ["//vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = ["remap_test.go"],
    embed = [":go_default_library"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["conformance.pb.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/protobuf/conformance/internal/conformance_proto",
    importpath = "github.com/golang/protobuf/conformance/internal/conformance_proto",
    visibility = ["//vendor/github.com/golang/protobuf/conformance:__subpackages__"],
    deps = [
        "//vendor/github.com/golang/protobuf/proto:go_default_library",
        "//vendor/github.com/golang/protobuf/ptypes/any:go_default_library",
        "//vendor/github.com/golang/protobuf/ptypes/duration:go_default_library",
        "//vendor/github.com/golang/protobuf/ptypes/struct:go_default_library",
        "//vendor/github.com/golang/protobuf/ptypes/timestamp:go_default_library",
        "//vendor/github.com/golang/protobuf/ptypes/wrappers:go_default_library",
        "@org_golang_google_genproto//protobuf/field_mask:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "environment.go",
        "plugin.pb.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/googleapis/gnostic/plugins",
    importpath = "github.com/googleapis/gnostic/plugins",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/protobuf/proto:go_default_library",
        "//vendor/github.com/golang/protobuf/ptypes/any:go_default_library",
        "//vendor/github.com/googleapis/gnostic/OpenAPIv2:go_default_library",
        "//vendor/github.com/googleapis/gnostic/OpenAPIv3:go_default_library",
        "//vendor/github.com/googleapis/gnostic/discovery:go_default_library",
        "//vendor/github.com/googleapis/gnostic/surface:go_default_library",
    ],
)

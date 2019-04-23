package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bitmap.go",
        "interfaces.go",
        "utils.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/kubernetes/pkg/registry/core/service/allocator",
    importpath = "k8s.io/kubernetes/pkg/registry/core/service/allocator",
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//pkg/registry/core/service/allocator/storage:all-srcs",
    ],
    tags = ["automanaged"],
)

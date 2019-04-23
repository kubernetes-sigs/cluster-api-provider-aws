package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["constants.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/kubernetes/cmd/kubeadm/app/constants",
    importpath = "k8s.io/kubernetes/cmd/kubeadm/app/constants",
    deps = [
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/version:go_default_library",
        "//vendor/k8s.io/cluster-bootstrap/token/api:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/registry/core/service/ipallocator:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
)

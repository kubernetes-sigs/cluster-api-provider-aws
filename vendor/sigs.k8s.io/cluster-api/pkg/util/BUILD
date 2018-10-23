load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "retry.go",
        "util.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/cluster-api/pkg/util",
    importpath = "sigs.k8s.io/cluster-api/pkg/util",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
    ],
)

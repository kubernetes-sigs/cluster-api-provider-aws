load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "exec.go",
        "minikube.go",
        "root.go",
        "testing.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/hack/cluster-api-dev-helper/cmd",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/spf13/cobra:go_default_library"],
)

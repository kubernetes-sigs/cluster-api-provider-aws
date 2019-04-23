load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["client.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/client",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/awsprovider/v1alpha1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset:go_default_library",
    ],
)

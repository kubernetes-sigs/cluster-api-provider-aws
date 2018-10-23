load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["openstack.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/client-go/plugin/pkg/client/auth/openstack",
    importpath = "k8s.io/client-go/plugin/pkg/client/auth/openstack",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "@com_github_gophercloud_gophercloud//:go_default_library",
        "@com_github_gophercloud_gophercloud//openstack:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["openstack_test.go"],
    embed = [":go_default_library"],
)

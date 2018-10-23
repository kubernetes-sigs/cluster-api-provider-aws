load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["gcp.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/client-go/plugin/pkg/client/auth/gcp",
    importpath = "k8s.io/client-go/plugin/pkg/client/auth/gcp",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/yaml:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/util/jsonpath:go_default_library",
        "@org_golang_x_oauth2//:go_default_library",
        "@org_golang_x_oauth2//google:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["gcp_test.go"],
    embed = [":go_default_library"],
    deps = ["@org_golang_x_oauth2//:go_default_library"],
)

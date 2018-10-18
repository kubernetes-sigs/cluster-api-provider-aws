load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["azure.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/client-go/plugin/pkg/client/auth/azure",
    importpath = "k8s.io/client-go/plugin/pkg/client/auth/azure",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "@com_github_azure_go_autorest//autorest:go_default_library",
        "@com_github_azure_go_autorest//autorest/adal:go_default_library",
        "@com_github_azure_go_autorest//autorest/azure:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["azure_test.go"],
    embed = [":go_default_library"],
    deps = ["@com_github_azure_go_autorest//autorest/adal:go_default_library"],
)

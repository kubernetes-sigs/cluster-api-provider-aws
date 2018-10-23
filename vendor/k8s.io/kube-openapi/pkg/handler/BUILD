load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["handler.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/kube-openapi/pkg/handler",
    importpath = "k8s.io/kube-openapi/pkg/handler",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/emicklei/go-restful:go_default_library",
        "//vendor/github.com/golang/protobuf/proto:go_default_library",
        "//vendor/github.com/googleapis/gnostic/OpenAPIv2:go_default_library",
        "//vendor/github.com/googleapis/gnostic/compiler:go_default_library",
        "//vendor/gopkg.in/yaml.v2:go_default_library",
        "//vendor/k8s.io/kube-openapi/pkg/builder:go_default_library",
        "//vendor/k8s.io/kube-openapi/pkg/common:go_default_library",
        "@com_github_go_openapi_spec//:go_default_library",
        "@com_github_nytimes_gziphandler//:go_default_library",
        "@org_bitbucket_ww_goautoneg//:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["handler_test.go"],
    embed = [":go_default_library"],
    deps = ["@com_github_go_openapi_spec//:go_default_library"],
)

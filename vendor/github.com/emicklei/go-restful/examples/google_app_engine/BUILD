load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "restful-appstats-integration.go",
        "restful-user-service.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/emicklei/go-restful/examples/google_app_engine",
    importpath = "github.com/emicklei/go-restful/examples/google_app_engine",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/github.com/emicklei/go-restful:go_default_library",
        "@com_github_emicklei_go_restful_swagger12//:go_default_library",
        "@com_github_mjibson_appstats//:go_default_library",
        "@org_golang_google_appengine//:go_default_library",
        "@org_golang_google_appengine//memcache:go_default_library",
    ],
)

go_binary(
    name = "google_app_engine",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

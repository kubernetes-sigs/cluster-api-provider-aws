load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "restful-CORS-filter.go",
        "restful-NCSA-logging.go",
        "restful-basic-authentication.go",
        "restful-cpuprofiler-service.go",
        "restful-curly-router.go",
        "restful-encoding-filter.go",
        "restful-filters.go",
        "restful-form-handling.go",
        "restful-hello-world.go",
        "restful-html-template.go",
        "restful-multi-containers.go",
        "restful-no-cache-filter.go",
        "restful-openapi.go",
        "restful-options-filter.go",
        "restful-path-tail.go",
        "restful-pre-post-filters.go",
        "restful-resource-functions.go",
        "restful-serve-static.go",
        "restful-swagger.go",
        "restful-user-resource.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/emicklei/go-restful/examples",
    importpath = "github.com/emicklei/go-restful/examples",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/github.com/emicklei/go-restful:go_default_library",
        "@com_github_emicklei_go_restful_openapi//:go_default_library",
        "@com_github_emicklei_go_restful_swagger12//:go_default_library",
        "@com_github_go_openapi_spec//:go_default_library",
        "@com_github_gorilla_schema//:go_default_library",
    ],
)

go_binary(
    name = "examples",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "restful-curly-router_test.go",
        "restful-route_test.go",
        "restful-routefunction_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/emicklei/go-restful:go_default_library"],
)

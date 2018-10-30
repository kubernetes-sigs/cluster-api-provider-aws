load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "goformation.go",
        "test.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/awslabs/goformation",
    importpath = "github.com/awslabs/goformation",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/awslabs/goformation/cloudformation:go_default_library",
        "//vendor/github.com/awslabs/goformation/intrinsics:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "goformation_schema_test.go",
        "goformation_suite_test.go",
        "goformation_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/awslabs/goformation/cloudformation:go_default_library",
        "//vendor/github.com/awslabs/goformation/intrinsics:go_default_library",
        "//vendor/github.com/johandorland/gojsonschema:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/github.com/onsi/gomega/gstruct:go_default_library",
        "//vendor/github.com/sanathkr/yaml:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "generate.go",
        "main.go",
        "naming.go",
        "property.go",
        "resource.go",
        "spec.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/awslabs/goformation/generate",
    importpath = "github.com/awslabs/goformation/generate",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "generate",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "generate_suite_test.go",
        "property_test.go",
        "resource_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/awslabs/goformation/cloudformation:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
    ],
)

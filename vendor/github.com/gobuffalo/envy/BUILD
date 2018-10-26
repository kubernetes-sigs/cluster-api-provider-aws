load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "envy.go",
        "version.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/gobuffalo/envy",
    importpath = "github.com/gobuffalo/envy",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/joho/godotenv:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["envy_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/stretchr/testify/require:go_default_library"],
)

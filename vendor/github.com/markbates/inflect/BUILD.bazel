load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "helpers.go",
        "inflect.go",
        "name.go",
        "version.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/markbates/inflect",
    importpath = "github.com/markbates/inflect",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/gobuffalo/envy:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "inflect_test.go",
        "name_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/stretchr/testify/require:go_default_library"],
)

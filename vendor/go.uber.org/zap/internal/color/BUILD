load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["color.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/zap/internal/color",
    importpath = "go.uber.org/zap/internal/color",
    visibility = ["//vendor/go.uber.org/zap:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = ["color_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/stretchr/testify/assert:go_default_library"],
)

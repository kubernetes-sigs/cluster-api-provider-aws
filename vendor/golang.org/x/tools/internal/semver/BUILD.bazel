load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["semver.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/internal/semver",
    importpath = "golang.org/x/tools/internal/semver",
    visibility = ["//vendor/golang.org/x/tools:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = ["semver_test.go"],
    embed = [":go_default_library"],
)

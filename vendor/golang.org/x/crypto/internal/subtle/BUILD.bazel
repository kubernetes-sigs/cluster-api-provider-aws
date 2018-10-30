load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["aliasing.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/crypto/internal/subtle",
    importpath = "golang.org/x/crypto/internal/subtle",
    visibility = ["//vendor/golang.org/x/crypto:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = ["aliasing_test.go"],
    embed = [":go_default_library"],
)

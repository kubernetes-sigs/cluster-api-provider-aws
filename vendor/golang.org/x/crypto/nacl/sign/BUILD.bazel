load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["sign.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/crypto/nacl/sign",
    importpath = "golang.org/x/crypto/nacl/sign",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/crypto/ed25519:go_default_library",
        "//vendor/golang.org/x/crypto/internal/subtle:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["sign_test.go"],
    embed = [":go_default_library"],
)

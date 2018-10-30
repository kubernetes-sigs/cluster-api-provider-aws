load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/cmd/bundle",
    importpath = "golang.org/x/tools/cmd/bundle",
    visibility = ["//visibility:private"],
    deps = ["//vendor/golang.org/x/tools/go/loader:go_default_library"],
)

go_binary(
    name = "bundle",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["main_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/go/buildutil:go_default_library"],
)

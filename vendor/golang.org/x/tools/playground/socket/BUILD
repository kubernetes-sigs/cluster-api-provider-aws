load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["socket.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/playground/socket",
    importpath = "golang.org/x/tools/playground/socket",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/net/websocket:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["socket_test.go"],
    embed = [":go_default_library"],
)

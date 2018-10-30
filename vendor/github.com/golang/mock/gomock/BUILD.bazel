load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "call.go",
        "callset.go",
        "controller.go",
        "matchers.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/gomock",
    importpath = "github.com/golang/mock/gomock",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/net/context:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "call_test.go",
        "callset_test.go",
        "controller_test.go",
        "matchers_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/golang/mock/gomock/mock_matcher:go_default_library"],
)

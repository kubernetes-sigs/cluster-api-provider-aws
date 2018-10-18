load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "go_default_test",
    srcs = ["sample_test.go"],
    deps = ["//vendor/github.com/googleapis/gnostic/plugins/gnostic-go-generator/examples/v2.0/sample/sample:go_default_library"],
)

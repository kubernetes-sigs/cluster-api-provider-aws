load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "interface_fuzzer.go",
        "reflect_deepcopy.go",
        "value_fuzzer.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests",
    importpath = "k8s.io/gengo/examples/deepcopy-gen/output_tests",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/google/gofuzz:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/aliases:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/interfaces:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["output_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/google/gofuzz:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/aliases:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/builtins:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/interfaces:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/maps:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/pointer:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/slices:go_default_library",
        "//vendor/k8s.io/gengo/examples/deepcopy-gen/output_tests/structs:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "apic.go",
        "decode.go",
        "emitterc.go",
        "encode.go",
        "parserc.go",
        "readerc.go",
        "resolve.go",
        "scannerc.go",
        "sorter.go",
        "writerc.go",
        "yaml.go",
        "yamlh.go",
        "yamlprivateh.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/sanathkr/go-yaml",
    importpath = "github.com/sanathkr/go-yaml",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "decode_test.go",
        "encode_test.go",
        "example_embedded_test.go",
        "suite_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/gopkg.in/yaml.v2:go_default_library",
        "@in_gopkg_check_v1//:go_default_library",
    ],
)

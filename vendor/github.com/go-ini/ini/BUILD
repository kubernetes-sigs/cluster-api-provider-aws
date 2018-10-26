load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "error.go",
        "file.go",
        "ini.go",
        "key.go",
        "parser.go",
        "section.go",
        "struct.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/go-ini/ini",
    importpath = "github.com/go-ini/ini",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "bench_test.go",
        "file_test.go",
        "ini_internal_test.go",
        "ini_test.go",
        "key_test.go",
        "parser_test.go",
        "section_test.go",
        "struct_test.go",
    ],
    data = glob(["testdata/**"]),
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/ini.v1:go_default_library",
    ],
)

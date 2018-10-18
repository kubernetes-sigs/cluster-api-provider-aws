load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "gcexportdata.go",
        "importer.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/gcexportdata",
    importpath = "golang.org/x/tools/go/gcexportdata",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/tools/go/internal/gcimporter:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "example_test.go",
        "gcexportdata_test.go",
    ],
    data = glob(["testdata/**"]),
    embed = [":go_default_library"],
)

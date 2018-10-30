load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["gccgoexportdata.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/gccgoexportdata",
    importpath = "golang.org/x/tools/go/gccgoexportdata",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/tools/go/internal/gccgoimporter:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["gccgoexportdata_test.go"],
    data = glob(["testdata/**"]),
    embed = [":go_default_library"],
)

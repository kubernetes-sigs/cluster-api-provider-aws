load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "bugreport.go",
        "bugreport_mock.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen/tests/aux_imports_embedded_interface",
    importpath = "github.com/golang/mock/mockgen/tests/aux_imports_embedded_interface",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/golang/mock/mockgen/tests/aux_imports_embedded_interface/faux:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["bugreport_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/golang/mock/gomock:go_default_library"],
)

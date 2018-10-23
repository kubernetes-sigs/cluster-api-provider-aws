load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["greeter.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/golang/mock/mockgen/tests/custom_package_name/greeter",
    importpath = "github.com/golang/mock/mockgen/tests/custom_package_name/greeter",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/mock/mockgen/tests/custom_package_name/client/v1:go_default_library",
        "//vendor/github.com/golang/mock/mockgen/tests/custom_package_name/validator:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "greeter_mock_test.go",
        "greeter_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/golang/mock/mockgen/tests/custom_package_name/client/v1:go_default_library",
    ],
)

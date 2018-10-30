load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "deepcopy_generated.go",
        "doc.go",
        "pullrequest_types.go",
        "register.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-tools/examples/godocbot/pkg/apis/code/v1alpha1",
    importpath = "sigs.k8s.io/controller-tools/examples/godocbot/pkg/apis/code/v1alpha1",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["pullrequest_types_test.go"],
    embed = [":go_default_library"],
)

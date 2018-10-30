load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["install.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install",
    importpath = "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions:go_default_library",
        "//vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["roundtrip_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/fuzzer:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/testing/roundtrip:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "decode.go",
        "doc.go",
        "handler.go",
        "http.go",
        "response.go",
        "tls.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/internal/admission",
    importpath = "sigs.k8s.io/controller-runtime/pkg/internal/admission",
    visibility = ["//vendor/sigs.k8s.io/controller-runtime/pkg:__subpackages__"],
    deps = [
        "//vendor/k8s.io/api/admission/v1beta1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "example_admissionfunc_test.go",
        "example_decode_test.go",
        "example_handlefunc_test.go",
        "example_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/k8s.io/api/admission/v1beta1:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "decode.go",
        "doc.go",
        "http.go",
        "response.go",
        "webhook.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/webhook/admission",
    importpath = "sigs.k8s.io/controller-runtime/pkg/webhook/admission",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/mattbaird/jsonpatch:go_default_library",
        "//vendor/k8s.io/api/admission/v1beta1:go_default_library",
        "//vendor/k8s.io/api/admissionregistration/v1beta1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/patch:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/inject:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/admission/types:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/types:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "admission_suite_test.go",
        "decode_test.go",
        "http_test.go",
        "response_test.go",
        "webhook_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/mattbaird/jsonpatch:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/k8s.io/api/admission/v1beta1:go_default_library",
        "//vendor/k8s.io/api/admissionregistration/v1beta1:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/admission/types:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/types:go_default_library",
    ],
)

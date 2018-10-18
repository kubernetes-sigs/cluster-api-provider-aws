load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "provisioner.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert",
    importpath = "sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert",
    visibility = ["//vendor/sigs.k8s.io/controller-runtime/pkg/webhook:__subpackages__"],
    deps = [
        "//vendor/k8s.io/api/admissionregistration/v1beta1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/generator:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/writer:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "provisioner_test.go",
        "suite_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/k8s.io/api/admissionregistration/v1beta1:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/generator:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/writer:go_default_library",
    ],
)

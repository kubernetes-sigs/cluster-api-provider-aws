load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "certwriter.go",
        "doc.go",
        "error.go",
        "fs.go",
        "secret.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/writer",
    importpath = "sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/writer",
    visibility = ["//vendor/sigs.k8s.io/controller-runtime/pkg/webhook:__subpackages__"],
    deps = [
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/generator:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/writer/atomic:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "certwriter_test.go",
        "fs_test.go",
        "secret_test.go",
        "suite_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client/fake:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/generator:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/generator/fake:go_default_library",
    ],
)

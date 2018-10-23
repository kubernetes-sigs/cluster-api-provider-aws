load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "inject.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/runtime/inject",
    importpath = "sigs.k8s.io/controller-runtime/pkg/runtime/inject",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/cache:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/webhook/admission/types:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "inject_suite_test.go",
        "inject_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/cache:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/cache/informertest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
    ],
)

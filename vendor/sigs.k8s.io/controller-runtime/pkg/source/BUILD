load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "source.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/source",
    importpath = "sigs.k8s.io/controller-runtime/pkg/source",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/client-go/util/workqueue:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/cache:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/event:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/handler:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/predicate:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/inject:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/source/internal:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "example_test.go",
        "source_integration_test.go",
        "source_suite_test.go",
        "source_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/k8s.io/api/apps/v1:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/util/workqueue:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/cache:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/cache/informertest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/controller:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/event:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/handler:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/predicate:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/inject:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
    ],
)

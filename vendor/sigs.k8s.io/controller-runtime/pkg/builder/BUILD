load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "build.go",
        "doc.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/builder",
    importpath = "sigs.k8s.io/controller-runtime/pkg/builder",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client/apiutil:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client/config:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/controller:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/handler:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/manager:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/predicate:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/reconcile:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/source:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "build_test.go",
        "builder_suite_test.go",
        "example_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/k8s.io/api/apps/v1:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client/apiutil:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/controller:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/manager:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/reconcile:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/signals:go_default_library",
    ],
)

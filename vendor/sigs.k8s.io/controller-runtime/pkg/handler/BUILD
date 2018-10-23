load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "enqueue.go",
        "enqueue_mapped.go",
        "enqueue_owner.go",
        "eventhandler.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/handler",
    importpath = "sigs.k8s.io/controller-runtime/pkg/handler",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/client-go/util/workqueue:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/event:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/reconcile:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/inject:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "eventhandler_suite_test.go",
        "eventhandler_test.go",
        "example_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/k8s.io/api/apps/v1:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//vendor/k8s.io/client-go/util/workqueue:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/controller:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/controller/controllertest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/event:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/reconcile:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/log:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/source:go_default_library",
    ],
)

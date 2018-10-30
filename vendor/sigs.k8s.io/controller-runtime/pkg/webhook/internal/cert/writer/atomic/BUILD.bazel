load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["atomic_writer.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/writer/atomic",
    importpath = "sigs.k8s.io/controller-runtime/pkg/webhook/internal/cert/writer/atomic",
    visibility = ["//vendor/sigs.k8s.io/controller-runtime/pkg/webhook:__subpackages__"],
    deps = [
        "//vendor/github.com/go-logr/logr:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["atomic_writer_test.go"],
    embed = [":go_default_library"],
    deps = select({
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/github.com/go-logr/logr/testing:go_default_library",
            "//vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
            "//vendor/k8s.io/client-go/util/testing:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

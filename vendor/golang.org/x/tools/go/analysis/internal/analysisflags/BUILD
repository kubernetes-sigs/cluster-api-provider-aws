load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["flags.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/analysis/internal/analysisflags",
    importpath = "golang.org/x/tools/go/analysis/internal/analysisflags",
    visibility = ["//vendor/golang.org/x/tools/go/analysis:__subpackages__"],
    deps = ["//vendor/golang.org/x/tools/go/analysis:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["flags_test.go"],
    embed = [":go_default_library"],
    deps = ["//vendor/golang.org/x/tools/go/analysis:go_default_library"],
)

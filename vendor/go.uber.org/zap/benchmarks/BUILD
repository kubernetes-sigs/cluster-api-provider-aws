load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["doc.go"],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/go.uber.org/zap/benchmarks",
    importpath = "go.uber.org/zap/benchmarks",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "apex_test.go",
        "kit_test.go",
        "lion_test.go",
        "log15_test.go",
        "logrus_test.go",
        "scenario_bench_test.go",
        "zap_test.go",
        "zerolog_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/go.uber.org/multierr:go_default_library",
        "//vendor/go.uber.org/zap:go_default_library",
        "//vendor/go.uber.org/zap/internal/ztest:go_default_library",
        "//vendor/go.uber.org/zap/zapcore:go_default_library",
        "@com_github_apex_log//:go_default_library",
        "@com_github_apex_log//handlers/json:go_default_library",
        "@com_github_go_kit_kit//log:go_default_library",
        "@com_github_rs_zerolog//:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
        "@in_gopkg_inconshreveable_log15_v2//:go_default_library",
        "@io_pedge_go_lion//:go_default_library",
    ],
)

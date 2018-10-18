load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "conditions.go",
        "fnand.go",
        "fnbase64.go",
        "fnequals.go",
        "fnfindinmap.go",
        "fngetatt.go",
        "fngetazs.go",
        "fnif.go",
        "fnjoin.go",
        "fnnot.go",
        "fnor.go",
        "fnselect.go",
        "fnsplit.go",
        "fnsub.go",
        "intrinsics.go",
        "ref.go",
        "tags.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/awslabs/goformation/intrinsics",
    importpath = "github.com/awslabs/goformation/intrinsics",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/imdario/mergo:go_default_library",
        "//vendor/github.com/sanathkr/go-yaml:go_default_library",
        "//vendor/github.com/sanathkr/yaml:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "intrinsics_suite_test.go",
        "intrinsics_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
    ],
)

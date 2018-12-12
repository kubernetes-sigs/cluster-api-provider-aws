load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_test(
    name = "go_default_test",
    size = "large",
    srcs = [
        "e2e_suite_test.go",
        "metacluster_test.go",
    ],
    args = [
        "-kindBinary=$(location @io_k8s_sigs_kind//:kind)",
        "-kubectlBinary=$(location @io_k8s_kubernetes//cmd/kubectl:kubectl)",
        "-awsProviderYAML=$(location //config:aws-provider-yaml)",
        "-clusterAPIYAML=$(location //vendor/sigs.k8s.io/cluster-api/config:cluster-api-yaml)",
        "-managerImageTar=$(location //cmd/manager:manager-amd64.tar)",
    ],
    data = [
        "//config:aws-provider-yaml",
        "//vendor/sigs.k8s.io/cluster-api/config:cluster-api-yaml",
        "//cmd/manager:manager-amd64.tar",
        "@io_k8s_kubernetes//cmd/kubectl:kubectl",
        "@io_k8s_sigs_kind//:kind",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    deps = [
        "//test/e2e/util/kind:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/github.com/onsi/gomega/gstruct:go_default_library",
        "//vendor/github.com/onsi/gomega/types:go_default_library",
        "//vendor/k8s.io/api/apps/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["e2e.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/test/e2e",
    visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "go_default_test",
    size = "large",
    srcs = [
        "aws_test.go",
        "e2e_suite_test.go",
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
    rundir = ".",
    deps = [
        "//pkg/apis/awsprovider/v1alpha1:go_default_library",
        "//test/e2e/util/kind:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/github.com/onsi/gomega/gstruct:go_default_library",
        "//vendor/github.com/onsi/gomega/types:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset:go_default_library",
    ],
)

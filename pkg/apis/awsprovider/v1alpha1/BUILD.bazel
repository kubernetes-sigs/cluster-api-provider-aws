load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "awsclusterproviderconfig_types.go",
        "awsclusterproviderstatus_types.go",
        "awsmachineproviderconfig_types.go",
        "awsmachineproviderstatus_types.go",
        "doc.go",
        "register.go",
        "types.go",
        "zz_generated.deepcopy.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/cloud/aws/tags:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/json:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/runtime/scheme:go_default_library",
        "//vendor/sigs.k8s.io/yaml:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "awsclusterproviderconfig_types_test.go",
        "awsclusterproviderstatus_types_test.go",
        "awsmachineproviderconfig_types_test.go",
        "awsmachineproviderstatus_types_test.go",
        "register_test.go",
        "types_test.go",
        "v1alpha1_suite_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["manual"],
    deps = [
        "//pkg/cloudtest:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/golang.org/x/net/context:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/envtest:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["actuator.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/cluster",
    visibility = ["//visibility:public"],
    deps = [
        "//cloud/aws/providerconfig/v1alpha1:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["actuator_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//cloud/aws/providerconfig/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

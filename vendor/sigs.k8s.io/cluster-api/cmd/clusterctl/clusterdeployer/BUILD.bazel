load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "clusterdeployer.go",
        "providercomponentsstorefactory.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer",
    importpath = "sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/providercomponents:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/util:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["clusterdeployer_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

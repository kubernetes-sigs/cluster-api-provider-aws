load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "create.go",
        "create_cluster.go",
        "delete.go",
        "delete_cluster.go",
        "logutil.go",
        "root.go",
        "validate.go",
        "validate_cluster.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/cmd",
    importpath = "sigs.k8s.io/cluster-api/cmd/clusterctl/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/ghodss/yaml:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
        "//vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clientcmd:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/bootstrap/existing:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/bootstrap/minikube:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/providercomponents:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/cmd/clusterctl/validation:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/common:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/util:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/client/config:go_default_library",
        "//vendor/sigs.k8s.io/controller-runtime/pkg/manager:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["create_cluster_test.go"],
    embed = [":go_default_library"],
)

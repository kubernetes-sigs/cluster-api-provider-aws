package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "arguments.go",
        "cgroupdriver.go",
        "chroot_unix.go",
        "chroot_windows.go",
        "copy.go",
        "endpoint.go",
        "error.go",
        "marshal.go",
        "template.go",
        "version.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/kubernetes/cmd/kubeadm/app/util",
    importpath = "k8s.io/kubernetes/cmd/kubeadm/app/util",
    deps = [
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/validation:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/version:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/yaml:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//vendor/k8s.io/klog:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubeadm/app/constants:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/version:go_default_library",
        "//vendor/k8s.io/utils/exec:go_default_library",
        "//vendor/sigs.k8s.io/yaml:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//cmd/kubeadm/app/util/apiclient:all-srcs",
        "//cmd/kubeadm/app/util/audit:all-srcs",
        "//cmd/kubeadm/app/util/config:all-srcs",
        "//cmd/kubeadm/app/util/dryrun:all-srcs",
        "//cmd/kubeadm/app/util/etcd:all-srcs",
        "//cmd/kubeadm/app/util/kubeconfig:all-srcs",
        "//cmd/kubeadm/app/util/pkiutil:all-srcs",
        "//cmd/kubeadm/app/util/pubkeypin:all-srcs",
        "//cmd/kubeadm/app/util/runtime:all-srcs",
        "//cmd/kubeadm/app/util/staticpod:all-srcs",
        "//cmd/kubeadm/app/util/system:all-srcs",
    ],
    tags = ["automanaged"],
)

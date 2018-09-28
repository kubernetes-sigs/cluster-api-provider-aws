load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:push.bzl", "container_push")
load("@io_bazel_rules_docker//contrib:push-all.bzl", "docker_push")
load("@io_bazel_rules_docker//container:container.bzl", "container_bundle")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cmd/cluster-controller",
    visibility = ["//visibility:private"],
    deps = [
        "//cloud/aws/controllers/cluster:go_default_library",
        "//cloud/aws/controllers/cluster/options:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/util/logs:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/controller/config:go_default_library",
    ],
)

go_binary(
    name = "cluster-controller",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_image(
    name = "cluster-controller-amd64",
    base = "@golang-image//image",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "linux",
    # "//conditions:default": fail("only building on Linux supported. Try again with --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 ")
    pure = "on",
)

container_bundle(
    name = "cluster-controller-image-dev",
    images = {
        "$(dev_registry)/$(dev_repository)/aws-cluster-controller:0.0.1": ":cluster-controller-amd64",
    },
)

docker_push(
    name = "cluster-controller-push-dev",
    bundle = ":cluster-controller-image-dev",
)

container_bundle(
    name = "cluster-controller-image",
    images = {
        "gcr.io/cluster-api/aws-cluster-controller:0.0.1": ":cluster-controller-amd64",
    },
)

docker_push(
    name = "cluster-controller-push",
    bundle = ":cluster-controller-image",
)

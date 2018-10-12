load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "bootstrap.go",
        "cloudformation.go",
        "service.go",
    ],
    importpath = "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/cloudformation",
    visibility = ["//visibility:public"],
    deps = [
        "//cloud/aws/services/awserrors:go_default_library",
        "//cloud/aws/services/iam:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/cloudformation:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@com_github_awslabs_goformation//cloudformation:go_default_library",
    ],
)

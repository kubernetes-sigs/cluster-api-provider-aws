load("//build:go_mock.bzl", "go_mock")

go_mock(
    name = "go_default_library",
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface",
    visibility = ["//visibility:public"],
    mocks = [
        {
            "prefix": "github.com/aws/aws-sdk-go",
            "package": "service/ec2/ec2iface",
            "interface": "EC2API",
            "vendored": True,
        },
    ],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2:go_default_library",
    ],
)

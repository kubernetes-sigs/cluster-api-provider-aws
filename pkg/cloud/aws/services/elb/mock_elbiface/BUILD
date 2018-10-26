load("//build:go_mock.bzl", "go_mock")

go_mock(
    name = "go_default_library",
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface",
    visibility = ["//visibility:public"],
    mocks = [
        {
            "prefix": "github.com/aws/aws-sdk-go",
            "package": "service/elb/elbiface",
            "interface": "ELBAPI",
            "vendored": True,
        },
    ],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/elb:go_default_library",
    ],
)

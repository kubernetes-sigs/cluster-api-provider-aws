load("//build:go_mock.bzl", "go_mock")

go_mock(
    name = "go_default_library",
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/mocks",
    visibility = ["//visibility:public"],
    mocks = [
        {
            "prefix": "sigs.k8s.io/cluster-api-provider-aws",
            "package": "pkg/cloud/aws/services",
            "interface": "EC2Interface",
            "vendored": False,
        },
        {
            "prefix": "sigs.k8s.io/cluster-api-provider-aws",
            "package": "pkg/cloud/aws/services",
            "interface": "ELBInterface",
            "vendored": False,
        },
    ],
    deps = [
        "//pkg/apis/awsprovider/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

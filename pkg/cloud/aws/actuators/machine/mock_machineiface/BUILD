load("//build:go_mock.bzl", "go_mock")

go_mock(
    name = "go_default_library",
    importpath = "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine/mock_machineiface",
    visibility = ["//visibility:public"],
    mocks = [
        {
            "prefix": "sigs.k8s.io/cluster-api",
            "package": "pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1",
            "interface": "MachineInterface",
            "vendored": True,
        },
    ],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
    ],
)

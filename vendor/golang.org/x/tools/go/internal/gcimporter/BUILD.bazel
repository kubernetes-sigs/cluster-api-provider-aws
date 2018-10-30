load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "bexport.go",
        "bimport.go",
        "exportdata.go",
        "gcimporter.go",
        "iimport.go",
        "newInterface10.go",
        "newInterface11.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/internal/gcimporter",
    importpath = "golang.org/x/tools/go/internal/gcimporter",
    visibility = ["//vendor/golang.org/x/tools/go:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "bexport_test.go",
        "gcimporter11_test.go",
        "gcimporter_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/golang.org/x/tools/go/buildutil:go_default_library",
        "//vendor/golang.org/x/tools/go/loader:go_default_library",
    ],
)

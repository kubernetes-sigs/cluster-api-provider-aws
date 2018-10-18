load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "dir.go",
        "doc.go",
        "main.go",
        "play.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/cmd/present",
    importpath = "golang.org/x/tools/cmd/present",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/golang.org/x/tools/godoc/static:go_default_library",
        "//vendor/golang.org/x/tools/playground:go_default_library",
        "//vendor/golang.org/x/tools/playground/socket:go_default_library",
        "//vendor/golang.org/x/tools/present:go_default_library",
    ],
)

go_binary(
    name = "present",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

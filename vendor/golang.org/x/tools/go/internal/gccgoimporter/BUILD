load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "backdoor.go",
        "gccgoinstallation.go",
        "importer.go",
        "newInterface10.go",
        "newInterface11.go",
        "parser.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/go/internal/gccgoimporter",
    importpath = "golang.org/x/tools/go/internal/gccgoimporter",
    visibility = ["//vendor/golang.org/x/tools/go:__subpackages__"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "gccgoinstallation_test.go",
        "importer_test.go",
        "parser_test.go",
        "testenv_test.go",
    ],
    embed = [":go_default_library"],
)

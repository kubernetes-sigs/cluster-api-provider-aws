load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "corpus.go",
        "dirtrees.go",
        "format.go",
        "godoc.go",
        "index.go",
        "linkify.go",
        "meta.go",
        "page.go",
        "parser.go",
        "pres.go",
        "search.go",
        "server.go",
        "snippet.go",
        "spec.go",
        "spot.go",
        "tab.go",
        "template.go",
        "versions.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/golang.org/x/tools/godoc",
    importpath = "golang.org/x/tools/godoc",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/golang.org/x/tools/godoc/analysis:go_default_library",
        "//vendor/golang.org/x/tools/godoc/env:go_default_library",
        "//vendor/golang.org/x/tools/godoc/util:go_default_library",
        "//vendor/golang.org/x/tools/godoc/vfs:go_default_library",
        "//vendor/golang.org/x/tools/godoc/vfs/httpfs:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "dirtrees_test.go",
        "godoc17_test.go",
        "godoc_test.go",
        "index_test.go",
        "server_test.go",
        "spec_test.go",
        "versions_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/golang.org/x/tools/godoc/vfs:go_default_library",
        "//vendor/golang.org/x/tools/godoc/vfs/gatefs:go_default_library",
        "//vendor/golang.org/x/tools/godoc/vfs/mapfs:go_default_library",
    ],
)

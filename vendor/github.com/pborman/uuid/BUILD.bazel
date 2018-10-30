load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "dce.go",
        "doc.go",
        "hash.go",
        "marshal.go",
        "node.go",
        "sql.go",
        "time.go",
        "util.go",
        "uuid.go",
        "version1.go",
        "version4.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/pborman/uuid",
    importpath = "github.com/pborman/uuid",
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/google/uuid:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "marshal_test.go",
        "seq_test.go",
        "sql_test.go",
        "uuid_test.go",
    ],
    embed = [":go_default_library"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "condition.go",
        "doc.go",
        "error.go",
        "expression.go",
        "key_condition.go",
        "operand.go",
        "projection.go",
        "update.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/service/dynamodb/expression",
    importpath = "github.com/aws/aws-sdk-go/service/dynamodb/expression",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/dynamodb:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "condition_test.go",
        "error_test.go",
        "examples_test.go",
        "expression_test.go",
        "key_condition_test.go",
        "operand_test.go",
        "projection_test.go",
        "update_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/dynamodb:go_default_library",
    ],
)

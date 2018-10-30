load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "host.go",
        "idempotency.go",
        "jsonvalue.go",
        "payload.go",
        "timestamp.go",
        "unmarshal.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/private/protocol",
    importpath = "github.com/aws/aws-sdk-go/private/protocol",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/client/metadata:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "host_test.go",
        "idempotency_test.go",
        "jsonvalue_test.go",
        "protocol_test.go",
        "unmarshal_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/client/metadata:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/awstesting:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/ec2query:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/jsonrpc:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/query:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/rest:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/restjson:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/private/protocol/restxml:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)

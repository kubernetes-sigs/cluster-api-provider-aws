load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "policy.go",
        "policy_json_1_6.go",
        "policy_json_1_7.go",
        "privkey.go",
        "randomreader.go",
        "sign_cookie.go",
        "sign_url.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/github.com/aws/aws-sdk-go/service/cloudfront/sign",
    importpath = "github.com/aws/aws-sdk-go/service/cloudfront/sign",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "policy_test.go",
        "privkey_test.go",
        "sign_cookie_example_test.go",
        "sign_cookie_test.go",
        "sign_url_test.go",
    ],
    embed = [":go_default_library"],
    deps = ["//vendor/github.com/aws/aws-sdk-go/awstesting/mock:go_default_library"],
)

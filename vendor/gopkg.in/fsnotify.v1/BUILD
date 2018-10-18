load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "fen.go",
        "fsnotify.go",
        "inotify.go",
        "inotify_poller.go",
        "kqueue.go",
        "open_mode_bsd.go",
        "open_mode_darwin.go",
        "windows.go",
    ],
    importmap = "sigs.k8s.io/cluster-api-provider-aws/vendor/gopkg.in/fsnotify.v1",
    importpath = "gopkg.in/fsnotify.v1",
    visibility = ["//visibility:public"],
    deps = select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

go_test(
    name = "go_default_test",
    srcs = [
        "example_test.go",
        "fsnotify_test.go",
        "inotify_poller_test.go",
        "inotify_test.go",
        "integration_darwin_test.go",
        "integration_test.go",
    ],
    embed = [":go_default_library"],
    deps = select({
        "@io_bazel_rules_go//go/platform:android": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:nacl": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "@com_github_fsnotify_fsnotify//:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

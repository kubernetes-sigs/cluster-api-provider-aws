# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@io_bazel_rules_go//go:def.bzl", "go_test")

alias(
    name = "integration",
    actual = ":go_default_test",
)

go_test(
    name = "go_default_test",
    size = "large",
    srcs = [
        "integration_suite_test.go",
        "metacluster_test.go",
    ],
    args = [
        "-kindBinary=$(location @io_k8s_sigs_kind//:kind)",
        "-kubectlBinary=$(location @io_k8s_kubernetes//cmd/kubectl:kubectl)",
        "-providerComponentsYAML=$(location manifests/provider-components-integration.yaml)",
        "-managerImageTar=$(location //cmd/manager:manager-amd64.tar)",
    ],
    data = [
        "manifests/provider-components-integration.yaml",
        "//cmd/manager:manager-amd64.tar",
        "@io_k8s_kubernetes//cmd/kubectl:kubectl",
        "@io_k8s_sigs_kind//:kind",
    ],
    rundir = ".",
    tags = ["external"],
    deps = [
        "//test/e2e/util/kind:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/github.com/onsi/gomega/gstruct:go_default_library",
        "//vendor/github.com/onsi/gomega/types:go_default_library",
        "//vendor/k8s.io/api/apps/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
    ],
)

genrule(
    name = "integration-provider-components",
    srcs = [
        "kustomization.yaml",
        "//test:manager-version-patch.yaml",
        "//config:provider-non-stamped-base.yaml",
    ],
    outs = ["manifests/provider-components-integration.yaml"],
    cmd = " && ".join([
        "cp -f $(SRCS) $(@D)",
        "$(location @io_k8s_sigs_kustomize//:kustomize) build $(@D) > $@",
    ]),
    tools = ["@io_k8s_sigs_kustomize//:kustomize"],
    visibility = ["//visibility:public"],
)

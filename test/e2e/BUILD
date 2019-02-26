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

KUSTOMIZE = "@io_k8s_sigs_kustomize//:kustomize"

KUSTOMIZE_CMD = "$(location %s) build $(@D)" % KUSTOMIZE

alias(
    name = "e2e",
    actual = ":go_default_test",
)

go_test(
    name = "go_default_test",
    size = "large",
    srcs = [
        "aws_test.go",
        "e2e_suite_test.go",
    ],
    args = [
        "-kindBinary=$(location @io_k8s_sigs_kind//:kind)",
        "-kubectlBinary=$(location @io_k8s_kubernetes//cmd/kubectl:kubectl)",
        "-providerComponentsYAML=$(location provider-components-e2e.yaml)",
        "-managerImageTar=$(location //cmd/manager:manager-amd64.tar)",
        "-credFile=$(location manifests/provider-credentials.profile)",
        "-clusterYAML=$(location manifests/cluster.yaml)",
        "-machineYAML=$(location manifests/machines.yaml)",
        "-regionFile=$(location region.txt)",
    ],
    data = [
        "manifests/addons.yaml",
        "manifests/cluster.yaml",
        "manifests/machines.yaml",
        "manifests/provider-credentials.profile",
        "provider-components-e2e.yaml",
        "region.txt",
        "//cmd/manager:manager-amd64.tar",
        "@io_k8s_kubernetes//cmd/kubectl:kubectl",
        "@io_k8s_sigs_kind//:kind",
    ],
    rundir = ".",
    deps = [
        "//pkg/apis/awsprovider/v1alpha1:go_default_library",
        "//pkg/cloud/aws/actuators:go_default_library",
        "//pkg/cloud/aws/actuators/machine:go_default_library",
        "//pkg/cloud/aws/services/awserrors:go_default_library",
        "//pkg/cloud/aws/services/cloudformation:go_default_library",
        "//pkg/cloud/aws/services/sts:go_default_library",
        "//test/e2e/util/kind:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/client:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/aws/session:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/cloudformation:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/ec2:go_default_library",
        "//vendor/github.com/aws/aws-sdk-go/service/sts:go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/github.com/onsi/gomega/gstruct:go_default_library",
        "//vendor/github.com/onsi/gomega/types:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset:go_default_library",
    ],
)

genrule(
    name = "manifests",
    srcs = [
        "//cmd/clusterctl/examples/aws:addons.yaml",
        "//cmd/clusterctl/examples/aws:cluster.yaml.template",
        "//cmd/clusterctl/examples/aws:cluster-network-spec.yaml.template",
        "//cmd/clusterctl/examples/aws:machines.yaml.template",
        "//cmd/clusterctl/examples/aws:provider-components-base.yaml",
        "//cmd/clusterctl/examples/aws:generate-yaml.sh",
        "manifests/provider-credentials.sh",
    ],
    outs = [
        "manifests/addons.yaml",
        "manifests/cluster.yaml",
        "manifests/machines.yaml",
        "manifests/provider-components.yaml",
    ],
    cmd = " && ".join([
        "source $(location manifests/provider-credentials.sh)",
        "cp -f $(SRCS) $(@D)",
        " ".join([
            "ENVSUBST=$(location @com_github_a8m_envsubst//cmd/envsubst:envsubst)",
            "CLUSTERAWSADM=$(location //cmd/clusterawsadm)",
            "OUTPUT_DIR=$(@D)/manifests",
            "DIR=$(@D)",
            "$(location //cmd/clusterctl/examples/aws:generate-yaml.sh) -f",
        ]),
    ]),
    tools = [
        "//cmd/clusterawsadm",
        "@com_github_a8m_envsubst//cmd/envsubst:envsubst",
    ],
    visibility = ["//visibility:public"],
)

genrule(
    name = "provider_credentials_profile",
    outs = ["manifests/provider-credentials.profile"],
    cmd = " && ".join([
        "touch $@",
        "export AWS_ACCESS_KEY_ID=$$(grep ^STABLE_AWS_ACCESS_KEY_ID bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "export AWS_SECRET_ACCESS_KEY=$$(grep ^STABLE_AWS_SECRET_ACCESS_KEY bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "export AWS_SESSION_TOKEN=$$(grep ^STABLE_AWS_SESSION_TOKEN bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "echo '[default]' >> $@",
        "echo aws_access_key_id = $$AWS_ACCESS_KEY_ID >> $@",
        "echo aws_secret_access_key = $$AWS_SECRET_ACCESS_KEY >> $@",
        "echo aws_session_token = $$AWS_SESSION_TOKEN >> $@",
    ]),
    stamp = 1,
    visibility = ["//visibility:private"],
)

genrule(
    name = "provider_credentials_env",
    outs = ["manifests/provider-credentials.sh"],
    cmd = " && ".join([
        "touch $@",
        "export AWS_ACCESS_KEY_ID=$$(grep ^STABLE_AWS_ACCESS_KEY_ID bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "export AWS_SECRET_ACCESS_KEY=$$(grep ^STABLE_AWS_SECRET_ACCESS_KEY bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "export AWS_SESSION_TOKEN=$$(grep ^STABLE_AWS_SESSION_TOKEN bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "export AWS_DEFAULT_REGION=$$(grep ^STABLE_AWS_DEFAULT_REGION bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "export AWS_REGION=$$(grep ^STABLE_AWS_REGION bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "echo export AWS_ACCESS_KEY_ID=$$AWS_ACCESS_KEY_ID >> $@",
        "echo export AWS_SECRET_ACCESS_KEY=$$AWS_SECRET_ACCESS_KEY >> $@",
        "echo export AWS_SESSION_TOKEN=$$AWS_SESSION_TOKEN >> $@",
        "echo export AWS_DEFAULT_REGION=$$AWS_DEFAULT_REGION >> $@",
        "echo export AWS_REGION=$$AWS_REGION >> $@",
    ]),
    stamp = 1,
    visibility = ["//visibility:private"],
)

genrule(
    name = "e2e-provider-components",
    srcs = [
        "kustomization.yaml",
        "//test:manager-version-patch.yaml",
        "manifests/provider-components.yaml",
    ],
    outs = ["provider-components-e2e.yaml"],
    cmd = " && ".join([
        "cp -f $(SRCS) $(@D)",
        "$(location @io_k8s_sigs_kustomize//:kustomize) build $(@D) > $@",
    ]),
    tools = ["@io_k8s_sigs_kustomize//:kustomize"],
    visibility = ["//visibility:public"],
)

genrule(
    name = "e2e-region",
    outs = ["region.txt"],
    cmd = " && ".join([
        "touch $@",
        "export AWS_DEFAULT_REGION=$$(grep ^STABLE_AWS_DEFAULT_REGION bazel-out/stable-status.txt | cut -f2 -d\" \")",
        "echo $$AWS_DEFAULT_REGION >> $@",
    ]),
    stamp = 1,
    visibility = ["//visibility:private"],
)

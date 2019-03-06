# Copyright 2018 The Kubernetes Authors.
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

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "ade51a315fa17347e5c31201fdc55aa5ffb913377aa315dceb56ee9725e620ee",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.16.6/rules_go-0.16.6.tar.gz",
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "aed1c249d4ec8f703edddf35cbe9dfaca0b5f5ea6e4cd9e83e99f3b0d1136c3d",
    strip_prefix = "rules_docker-0.7.0",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.7.0.tar.gz"],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "7949fc6cc17b5b191103e97481cf8889217263acf52e00b560683413af204fcb",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.16.0/bazel-gazelle-0.16.0.tar.gz"],
)

http_archive(
    name = "io_kubernetes_build",
    sha256 = "66a44fd5f6357268340d66fbd8a502065445a7c022732fe5f6fd84d9a20f75a3",
    strip_prefix = "repo-infra-e8f2f7c3decf03e1fde9f30d249e39b8328aa8b0",
    urls = ["https://github.com/kubernetes/repo-infra/archive/e8f2f7c3decf03e1fde9f30d249e39b8328aa8b0.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.11.5",
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

load(
    "@io_bazel_rules_docker//toolchains/docker:toolchain.bzl",
    docker_toolchain_configure = "toolchain_configure",
)
load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

container_pull(
    name = "golang-image",
    registry = "registry.hub.docker.com",
    repository = "library/golang",
    tag = "1.10-alpine",
)

go_repository(
    name = "com_github_golang_dep",
    build_file_generation = "on",
    importpath = "github.com/golang/dep",
    strip_prefix = "dep-22125cfaa6ddc71e145b1535d4b7ee9744fefff2",
    urls = ["https://github.com/golang/dep/archive/22125cfaa6ddc71e145b1535d4b7ee9744fefff2.zip"],
)

go_repository(
    name = "com_github_golangci_golangci-lint",
    build_file_generation = "on",
    importpath = "github.com/golangci/golangci-lint",
    tag = "v1.15.0",
)

go_repository(
    name = "com_github_golang_mock",
    build_file_generation = "on",
    importpath = "github.com/golang/mock",
    strip_prefix = "mock-8a44ef6e8be577e050008c7886f24fc705d709fb",
    urls = ["https://github.com/golang/mock/archive/8a44ef6e8be577e050008c7886f24fc705d709fb.zip"],
)

go_repository(
    name = "io_k8s_sigs_kind",
    commit = "1284e993a84a56c994ea541dbcb97486a7b86b50",
    importpath = "sigs.k8s.io/kind",
)

go_repository(
    name = "io_k8s_sigs_kustomize",
    commit = "58492e2d83c59ed63881311f46ad6251f77dabc3",
    importpath = "sigs.k8s.io/kustomize",
)

go_repository(
    name = "io_k8s_kubernetes",
    commit = "4ed3216f3ec431b140b1d899130a69fc671678f4",  # v1.12.1
    importpath = "k8s.io/kubernetes",
)

go_repository(
    name = "com_github_a8m_envsubst",
    commit = "41dec2456c86b2a9fa51a22a808b7084b8d52c64",  # v1.1.0
    importpath = "github.com/a8m/envsubst",
)

# for @io_k8s_kubernetes
http_archive(
    name = "io_kubernetes_build",
    sha256 = "1188feb932cefad328b0a3dd75b3ebd1d79dd26dbdd723f019ceb760e27ba6d8",
    strip_prefix = "repo-infra-84d52408a061e87d45aebf5a0867246bdf66d180",
    urls = ["https://github.com/kubernetes/repo-infra/archive/84d52408a061e87d45aebf5a0867246bdf66d180.tar.gz"],
)

git_repository(
    name = "io_kubernetes_repo_infra",
    commit = "b4bc4f1552c7fc1d4654753ca9b0e5e13883429f",
    remote = "https://github.com/kubernetes/repo-infra.git",
)

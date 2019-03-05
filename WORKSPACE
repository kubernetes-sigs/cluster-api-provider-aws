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
load("//build:workspace_mirror.bzl", "mirror")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6776d68ebb897625dead17ae510eac3d5f6342367327875210df44dbe2aeeb19",
    urls = mirror("https://github.com/bazelbuild/rules_go/releases/download/0.17.1/rules_go-0.17.1.tar.gz"),
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "aed1c249d4ec8f703edddf35cbe9dfaca0b5f5ea6e4cd9e83e99f3b0d1136c3d",
    strip_prefix = "rules_docker-0.7.0",
    urls = mirror("https://github.com/bazelbuild/rules_docker/archive/v0.7.0.tar.gz"),
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
    urls = mirror("https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"),
)

http_archive(
    name = "io_k8s_repo_infra",
    sha256 = "66a44fd5f6357268340d66fbd8a502065445a7c022732fe5f6fd84d9a20f75a3",
    strip_prefix = "repo-infra-e8f2f7c3decf03e1fde9f30d249e39b8328aa8b0",
    urls = mirror("https://github.com/kubernetes/repo-infra/archive/e8f2f7c3decf03e1fde9f30d249e39b8328aa8b0.tar.gz"),
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.12",
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
    tag = "1.12-alpine",
)

go_repository(
    name = "com_github_a8m_envsubst",
    importpath = "github.com/a8m/envsubst",
    sha256 = "2e52f2b94a52728077aabd48a3586414b8ff46563f31f597b899403eb8a25b58",
    strip_prefix = "envsubst-1.1.0",
    urls = mirror("https://github.com/a8m/envsubst/archive/v1.1.0.zip"),
)

go_repository(
    name = "com_github_golang_dep",
    build_file_generation = "on",
    importpath = "github.com/golang/dep",
    sha256 = "f063992cf3546713adca414c4044892bca59f18e630bf476dfd851971d1ffb3a",
    strip_prefix = "dep-8af3a37fb20df1b93a82cc4091eeaee18a5a9a63",
    urls = mirror("https://github.com/golang/dep/archive/8af3a37fb20df1b93a82cc4091eeaee18a5a9a63.zip"),
)

go_repository(
    name = "com_github_golangci_golangci-lint",
    build_file_generation = "on",
    importpath = "github.com/golangci/golangci-lint",
    sha256 = "c094917975d3e9c6c180bd1cf275c10495ad4b20ee30f739947d5b870affd8a4",
    strip_prefix = "golangci-lint-c55a62a8de8145b86a2ff0744da080e9124b71c8",
    urls = mirror("https://github.com/golangci/golangci-lint/archive/c55a62a8de8145b86a2ff0744da080e9124b71c8.zip"),
)

go_repository(
    name = "com_github_golang_mock",
    build_file_generation = "on",
    importpath = "github.com/golang/mock",
    sha256 = "129eb7bb48cacfa984e1bf1d66e4831aab27b7037d5fa402ac6ee3311909f984",
    strip_prefix = "mock-1.2.0",
    urls = mirror("https://github.com/golang/mock/archive/v1.2.0.zip"),
)

go_repository(
    name = "io_k8s_kubernetes",
    importpath = "k8s.io/kubernetes",
    sha256 = "ca6810419cf3ff8fe1907b23c783f52778ad9a007ea8663ba5510b086d9ee929",
    strip_prefix = "kubernetes-1.13.4",
    urls = mirror("https://github.com/kubernetes/kubernetes/archive/v1.13.4.zip"),
)

http_archive(
    name = "io_k8s_repo_infra",
    sha256 = "a03eb0ec374dfb206850f9ea6509b458e1676e98c739d202463601a1481f4db9",
    strip_prefix = "repo-infra-52d76ba3344e755f5a9ab595b94b21df49448ffb",
    urls = mirror("https://github.com/kubernetes/repo-infra/archive/52d76ba3344e755f5a9ab595b94b21df49448ffb.zip"),
)

# Needed by @io_k8s_kubernetes at 1.13.4. Can remove after 1.14.0
http_archive(
    name = "io_kubernetes_build",
    sha256 = "a03eb0ec374dfb206850f9ea6509b458e1676e98c739d202463601a1481f4db9",
    strip_prefix = "repo-infra-52d76ba3344e755f5a9ab595b94b21df49448ffb",
    urls = mirror("https://github.com/kubernetes/repo-infra/archive/52d76ba3344e755f5a9ab595b94b21df49448ffb.zip"),
)

go_repository(
    name = "io_k8s_sigs_kind",
    importpath = "sigs.k8s.io/kind",
    sha256 = "33e58b7b26a0cfe3cf03ae42fea773f96bd21092a829f1b01b1c473c33288d7d",
    strip_prefix = "kind-0.1.0",
    urls = mirror("https://github.com/kubernetes-sigs/kind/archive/0.1.0.zip"),
)

go_repository(
    name = "io_k8s_sigs_kustomize",
    importpath = "sigs.k8s.io/kustomize",
    sha256 = "acdd3482bc49e27bb23e519ba3b38f4043dba58a361c48ff44e67df0905ec0a6",
    strip_prefix = "kustomize-1.0.11",
    urls = mirror("https://github.com/kubernetes-sigs/kustomize/archive/v1.0.11.zip"),
)

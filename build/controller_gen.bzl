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

# TODO: Move this to Kubebuilder repository

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_kubernetes_build//defs:go.bzl", "go_genrule")

CONTROLLER_GEN = "//vendor/sigs.k8s.io/controller-tools/cmd/controller-gen"

def _qualified_genfile(label):
  return "$$GO_GENRULE_EXECROOT/$(location %s)" % label

# controller_gen generates CRD and RBAC manifests for Kubernetes
# controllers based on controller runtime
# (https://github.com/kubernetes-sigs/controller-runtime) given
# Kubebuilder annotations in the Go source code.
# controller_gen will output manifests to config/crds and config/rbac
# within the projects genfiles.
def controller_gen(name, importpath, api, visibility, deps = []):
    outs = [
      "rbac/rbac_role.yaml",
      "rbac/rbac_role_binding.yaml",
    ]

    real_deps = [
        "//pkg/apis:go_default_library",
        "//pkg/cloud/aws/actuators/cluster:go_default_library",
        "//pkg/cloud/aws/actuators/machine:go_default_library",
    ] + deps

    for g in api:
      group = g["group"]
      version = g["version"].lower()
      types = g["types"]
      prefix = group.split(".")[0].lower()
      real_deps += [ "//pkg/apis/%s:go_default_library" % prefix]
      for t in types:
        basename = t.lower()
        out = "crds/%s_%s_%s.yaml" % (prefix, version, basename)
        outs += [out]

    cmd = """mkdir -p {source_package} && \\
             cd {source_package} && \\
             cp -f {project} {source_package} && \\
             GENDIR=$$(dirname {gendir})/../.. && \\
            {controller_gen} all && \\
            cp -fR config $$GENDIR
          """.format(
      controller_gen = _qualified_genfile(CONTROLLER_GEN),
      project = _qualified_genfile("//:PROJECT"),
      gendir = _qualified_genfile(outs[0]),
      source_package =  "$$GOPATH/src/%s" % importpath
    )

    go_genrule(
        name = name,
        outs = outs,
        srcs = ["//:PROJECT"],
        cmd = cmd,
        go_deps = real_deps,
        visibility = visibility,
        tools = [CONTROLLER_GEN],
        tags = [ "generated" ],
    )

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

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_kubernetes_build//defs:go.bzl", "go_genrule")

MOCKGEN = "@com_github_golang_mock//mockgen"
MOCKGEN_LIBS = [
  "@com_github_golang_mock//mockgen/model:go_default_library",
  "//vendor/github.com/golang/mock/gomock:go_default_library",
]
ASM_SHIM = "//build/asm_shim"
ASM_SHIM_LIB = "%s:go_default_library" % ASM_SHIM
TEXTFLAG_SHIM = "%s:textflag.h" % ASM_SHIM
SDK_INCLUDE_DIR = "$$GOROOT/pkg/include"
BOILERPLATE = "//hack:boilerplate/boilerplate.go.txt"
GO_FLAGS = "CGO_ENABLED=0"

def _qualified_genfile(label):
  return "$$GO_GENRULE_EXECROOT/$(location %s)" % label

def go_mock(name, importpath, visibility, mocks, deps):
    targets = [name + ".go"]
    srcs = [m["interface"].lower() + ".go" for m in mocks]

    go_library(
        name = name,
        srcs = srcs,
        importpath = importpath,
        deps = deps + MOCKGEN_LIBS ,
        visibility = visibility,
    )

    for m in mocks:
      package = m["package"]
      interface = m["interface"]
      prefix = m["prefix"]

      out_basename = m["interface"].lower()
      out = "%s.go" % out_basename

      if m["vendored"]:
        extra_dep = ["//vendor/%s/%s:go_default_library" % (prefix, package)]
      else:
        extra_dep = [ "//%s:go_default_library" % package ]

      full_deps = [ASM_SHIM_LIB] + MOCKGEN_LIBS + deps + extra_dep

      cmd = """mkdir -p {source_package} && \\
mkdir -p {generated_package} && \\
mkdir -p {sdk_include_dir} && \\
cp {textflag_shim} {sdk_include_dir} && \\
cd {source_package} && \\
cat {boilerplate} | sed "s/YEAR/$$(date +%Y)/g" > {qualified_out} && \\
echo "\n\n" >> {qualified_out} && \\
{go_flags} {mockgen} -package={code_package} {qualified_package} {interface} \\
    >> {qualified_out}
""".format(
            qualified_package = prefix + "/" + package,
            code_package = importpath.split("/")[-1],
            importpath = importpath,
            textflag_shim = "$(location %s)" % TEXTFLAG_SHIM,
            generated_package = "$$GO_GENRULE_EXECROOT/%s" % importpath,
            source_package = "$$GOPATH/src/%s" % importpath,
            qualified_out = _qualified_genfile(":" + out),
            mockgen = _qualified_genfile(MOCKGEN),
            sdk_include_dir = SDK_INCLUDE_DIR,
            go_flags = GO_FLAGS,
            interface = interface,
            boilerplate = _qualified_genfile(BOILERPLATE),
          )

      go_genrule(
          name = out_basename,
          srcs = [BOILERPLATE, TEXTFLAG_SHIM],
          outs = [ out ],
          cmd = cmd,
          go_deps = full_deps,
          tools = [MOCKGEN],
      )

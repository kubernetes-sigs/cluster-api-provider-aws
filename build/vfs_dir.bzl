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

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_kubernetes_build//defs:go.bzl", "go_genrule")

VFSGEN = "@com_github_shurcool_vfsgen//:go_default_library"

GO = "@go_sdk//:bin/go"

def _asset_template(importpath):
  package = importpath.split("/")[-1]
  return """
package main

import (
  "github.com/shurcool/vfsgen"
  "log"
  "net/http"
  "path/filepath"   
)

func main() {{
  dir := filepath.Dir("./unpacked/")
  fs := http.Dir(dir)
  err := vfsgen.Generate(fs, vfsgen.Options{{
    PackageName: "{package}",
  }})
  if err != nil {{
	  log.Fatalln(err)
  }}
}}

""".format(package = package)

def vfs_dir(name, importpath, src, visibility = ["//visibility:public"]):

  go_library(
      name = name,
      srcs = ["assets_vfsdata.go"],
      importpath = importpath,
      visibility = visibility,
  )

  go_genrule(
    name = "%s_bindata" % name,
    srcs = [src],
    outs = ["assets_vfsdata.go"],
    cmd = """mkdir -p unpacked && tar xfv {tar_file} -C unpacked && cat << EOF > tmp.go && {go} run tmp.go && cp assets_vfsdata.go $@
{template}
EOF
    """.format(
        template = _asset_template(importpath),
        tar_file = "$(location %s)" % src,
        go = "$(location %s)" % GO,
    ),
    go_deps = [VFSGEN],
    tools = [GO],
    visibility = visibility,
    tags = [ "generated" ]    
)

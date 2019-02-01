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

"""
Create a manager stateful set patch file for use with Kustomize
"""

load("@bazel_skylib//lib:paths.bzl", "paths")

def _install_impl(ctx):

    cmd = ""

    for s in ctx.attr.srcs:
      files = s.files.to_list()
      file = files[0].path
      root_file = paths.relativize(file,"bazel-out/k8-fastbuild/genfiles/")
      if cmd == "":
        sep = ""
      else:
        sep = " &&"
      cmd = cmd + "%s cp -f %s %s" % (sep, file, root_file)

    cmd = cmd + " && touch %s" % ctx.outputs.record.path

    ctx.actions.run_shell(
      inputs = ctx.files.srcs,
      outputs = [ctx.outputs.record],
      command = cmd,
      use_default_shell_env = True,
      execution_requirements = {
        "no-sandbox": "1",
        "no-cache": "1",
        "no-remote": "1",
        "local": "1",
      },
    )

install = rule(
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_files = True,
        ),
    },
    output_to_genfiles = True,
    outputs = {
        # "script": "%{name}.install.sh",
        "record": ".install.%{name}.record",
    },
    implementation = _install_impl,
)

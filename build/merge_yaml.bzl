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

def _merge_yaml_impl(ctx):

    out = ctx.outputs.manifest.path

    cmd = "echo --- > %s" % out

    for s in ctx.attr.srcs:
      files = s.files.to_list()
      file = files[0].path

      cmd = cmd + " && cat {file} >> {out} && echo --- >> {out}".format(file = file, out = out)

    ctx.actions.run_shell(
      inputs = ctx.files.srcs,
      outputs = [ctx.outputs.manifest],
      progress_message = "Merging to YAML file: %s" % ctx.outputs.manifest.path,
      command = cmd
    )

merge_yaml = rule(
    attrs = {
        "srcs": attr.label_list(
            mandatory = True,
            allow_files = [
                ".yml",
                ".yaml",
                ".json",
            ],
        ),
        "out": attr.string(mandatory = True),
    },
    output_to_genfiles = True,
    outputs = {"manifest": "%{out}"},
    implementation = _merge_yaml_impl,
)

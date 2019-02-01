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

# Label of the template file to use.
_TEMPLATE = "//build:stateful_set_patch.yaml"

def _genfile_impl(ctx):
    ctx.actions.expand_template(
        template = ctx.file._template,
        output = ctx.outputs.source_file,
        substitutions = {
            "<docker_repo>": ctx.var[ctx.attr.docker_repo_var],
            "<image_name>": ctx.var["MANAGER_IMAGE_NAME"],
            "<tag>": ctx.var["MANAGER_IMAGE_TAG"]
        },
    )

    return [
      OutputGroupInfo(
          compilation_outputs = [ctx.outputs.source_file],
      ),      
      DefaultInfo(
        files = depset([ctx.outputs.source_file]),
        runfiles = ctx.runfiles(files = [ctx.outputs.source_file]),
    )]    

genfile = rule(
    attrs = {
        "srcs": attr.label(mandatory = True),
    },
    output_to_genfiles = True,
    outputs = {"source_file": "%{name}.yaml"},
    implementation = _genfile,
)

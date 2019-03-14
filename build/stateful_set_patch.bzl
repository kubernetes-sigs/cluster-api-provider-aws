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

def _stateful_set_patch_impl(ctx):
    ctx.actions.expand_template(
        template = ctx.file._template,
        output = ctx.outputs.source_file,
        substitutions = {
            "<registry>": ctx.expand_make_variables("registry", ctx.attr.registry, {}),
            "<image_name>": ctx.expand_make_variables("image_name", ctx.attr.image_name, {}),
            "<tag>": ctx.expand_make_variables("tag", ctx.attr.tag, {}),
            "<pull_policy>": ctx.expand_make_variables("pull_policy", ctx.attr.pull_policy, {}),
        },
    )

    return [
        OutputGroupInfo(
            compilation_outputs = [ctx.outputs.source_file],
        ),
        DefaultInfo(
            files = depset([ctx.outputs.source_file]),
            runfiles = ctx.runfiles(files = [ctx.outputs.source_file]),
        ),
    ]

stateful_set_patch = rule(
    attrs = {
        "registry": attr.string(mandatory = True),
        "pull_policy": attr.string(default = "IfNotPresent"),
        "image_name": attr.string(default = "$(MANAGER_IMAGE_NAME)"),
        "tag": attr.string(default = "$(MANAGER_IMAGE_TAG)"),
        "_template": attr.label(
            default = Label(_TEMPLATE),
            allow_single_file = True,
        ),
    },
    output_to_genfiles = True,
    outputs = {"source_file": "%{name}.yaml"},
    implementation = _stateful_set_patch_impl,
)

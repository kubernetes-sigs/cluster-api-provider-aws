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

# This is a modified version of the same rule from kubernetes/repo-infra
# modified to add the GO SDK to the PATH environment variable.

# Writes out a script which saves the runfiles directory,
# changes to the workspace root, and then runs a command.

def _workspace_binary_script_impl(ctx):
    content = """#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail
GOBINDIR=$(dirname $(readlink {go_bin}))
export PATH=$GOBINDIR:$PATH
BASE=$(pwd)
cd $(dirname $(readlink {root_file}))
"$BASE/{cmd}" $@
""".format(
        cmd = ctx.file.cmd.short_path,
        root_file = ctx.file.root_file.short_path,
        go_bin = ctx.file.go_bin.short_path,
    )
    ctx.actions.write(
        output = ctx.outputs.executable,
        content = content,
        is_executable = True,
    )
    runfiles = ctx.runfiles(
        files = [
            ctx.file.cmd,
            ctx.file.go_bin,
            ctx.file.root_file,
        ],
    )
    return [DefaultInfo(runfiles = runfiles)]

_workspace_binary_script = rule(
    attrs = {
        "cmd": attr.label(
            mandatory = True,
            allow_files = True,
            single_file = True,
        ),
        "root_file": attr.label(
            mandatory = True,
            allow_files = True,
            single_file = True,
        ),
        "go_bin": attr.label(
            mandatory = True,
            allow_files = True,
            single_file = True,
        ),
    },
    executable = True,
    implementation = _workspace_binary_script_impl,
)

# Wraps a binary to be run in the workspace root via bazel run.
#
# For example, one might do something like
#
# workspace_binary(
#     name = "dep",
#     cmd = "//vendor/github.com/golang/dep/cmd/dep",
# )
#
# which would allow running dep with bazel run.
def workspace_binary(
        name,
        cmd,
        args = None,
        visibility = None,
        go_bin = "@go_sdk//:bin/go",
        root_file = "//:WORKSPACE"):
    script_name = name + "_script"
    _workspace_binary_script(
        name = script_name,
        cmd = cmd,
        root_file = root_file,
        go_bin = go_bin,
        tags = ["manual"],
    )
    native.sh_binary(
        name = name,
        srcs = [":" + script_name],
        args = args,
        visibility = visibility,
        tags = ["manual"],
    )

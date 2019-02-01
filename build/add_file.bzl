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

def add_file(in_file, output, path = None):
  output_path = output
  input_path = in_file.path

  if path and in_file.short_path.startswith(path):
    output_path += in_file.short_path[len(path):]

  return [
      "mkdir -p $(dirname %s)" % output_path,
      "test -L %s || ln -s $(pwd)/%s %s" % (output_path, input_path, output_path),
      ]

def add_files(in_files, path = None):
  cmds = [add_file(in_file, in_file.path.build_output, path) for in_file in in_files]
  return cmds.join(" && ")

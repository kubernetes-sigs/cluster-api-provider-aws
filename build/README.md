# Build

Contains custom Bazel rules for this repository

## cluster_api_binary

Helpers to generate Docker images and binaries for Cluster API

## go_mock

A mock generator

See an [existing definition for invocation][mock_example]

## run_in_workspace

Runs a command in the workspace, with access to the Golang toolchain

[mock_example]: ../pkg/cloud/aws/actuators/cluster/mock_clusteriface/BUILD

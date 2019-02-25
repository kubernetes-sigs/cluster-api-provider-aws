#!/usr/bin/env bash

set -euo pipefail

base=$( dirname "${BASH_SOURCE[0]}")

go run "${base}/test-prerequisites.go"
TEST_INTEGRATION=1 go test ./... -test.run=^TestIntegration -args -alsologtostderr -v=5
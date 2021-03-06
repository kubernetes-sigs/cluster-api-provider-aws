#!/bin/sh

for TARGET in "${@}"; do
  find "${TARGET}" -name '*.go' ! -path '*/vendor/*' ! -path '*/.build/*' -exec gofmt -s -w {} \+
done
git diff --exit-code


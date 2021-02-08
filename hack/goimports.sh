#!/bin/sh

for TARGET in "${@}"; do
  find "${TARGET}" -name '*.go' ! -path '*/vendor/*' ! -path '*/.build/*' -exec goimports -w {} \+
done


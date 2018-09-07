#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  yamllint -d relaxed ./examples/
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/workdir:z" \
    --entrypoint sh \
    quay.io/coreos/yamllint \
    ./hack/yaml-lint.sh
fi;

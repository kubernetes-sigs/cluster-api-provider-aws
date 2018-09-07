#!/bin/sh
# Example:  ./hack/go-lint.sh cmd/... pkg/...

if [ "$IS_CONTAINER" != "" ]; then
  golint -set_exit_status "${@}"
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/sigs.k8s.io/cluster-api-provider-aws:z" \
    --workdir /go/src/sigs.k8s.io/cluster-api-provider-aws \
    --entrypoint sh \
    quay.io/coreos/golang-testing \
    ./hack/go-lint.sh "${@}"
fi

#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  go vet "${@}"
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/cluster-api-provider-aws:z" \
    --workdir /go/src/github.com/openshift/cluster-api-provider-aws \
    quay.io/coreos/golang-testing \
    ./hack/go-vet.sh "${@}"
fi;

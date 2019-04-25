#!/bin/sh

go test -timeout 90m \
  -v github.com/openshift/cluster-api-actuator-pkg/pkg/e2e \
  -kubeconfig ${KUBECONFIG:-~/.kube/config} \
  -machine-api-namespace ${NAMESPACE:-openshift-machine-api} \
  -args -v 5 -logtostderr \
  $@

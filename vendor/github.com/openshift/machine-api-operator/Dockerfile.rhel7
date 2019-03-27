FROM registry.svc.ci.openshift.org/ocp/builder:golang-1.10 AS builder
WORKDIR /go/src/github.com/openshift/machine-api-operator
COPY . .
RUN NO_DOCKER=1 make build

FROM registry.svc.ci.openshift.org/ocp/4.0:base
COPY --from=builder /go/src/github.com/openshift/machine-api-operator/owned-manifests owned-manifests
COPY --from=builder /go/src/github.com/openshift/machine-api-operator/install manifests
COPY --from=builder /go/src/github.com/openshift/machine-api-operator/bin/machine-api-operator .
COPY --from=builder /go/src/github.com/openshift/machine-api-operator/bin/nodelink-controller .
COPY --from=builder /go/src/github.com/openshift/machine-api-operator/bin/machine-healthcheck .
LABEL io.openshift.release.operator true

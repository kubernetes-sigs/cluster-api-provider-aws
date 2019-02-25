FROM registry.svc.ci.openshift.org/ocp/builder:golang-1.10 AS builder
WORKDIR /go/src/github.com/openshift/cluster-autoscaler-operator
COPY . .
ENV NO_DOCKER=1
ENV BUILD_DEST=/go/bin/cluster-autoscaler-operator
RUN unset VERSION && make build

FROM registry.svc.ci.openshift.org/ocp/4.0:base
COPY --from=builder /go/bin/cluster-autoscaler-operator /usr/bin/
COPY --from=builder /go/src/github.com/openshift/cluster-autoscaler-operator/install /manifests
CMD ["/usr/bin/cluster-autoscaler-operator"]
LABEL io.openshift.release.operator true

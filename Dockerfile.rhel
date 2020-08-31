FROM registry.svc.ci.openshift.org/ocp/builder:rhel-8-golang-openshift-4.6 AS builder
WORKDIR /go/src/sigs.k8s.io/cluster-api-provider-aws
COPY . .
# VERSION env gets set in the openshift/release image and refers to the golang version, which interfers with our own
RUN unset VERSION \
 && GOPROXY=off NO_DOCKER=1 make build

FROM registry.svc.ci.openshift.org/ocp/4.6:base
COPY --from=builder /go/src/sigs.k8s.io/cluster-api-provider-aws/bin/machine-controller-manager /
COPY --from=builder /go/src/sigs.k8s.io/cluster-api-provider-aws/bin/termination-handler /

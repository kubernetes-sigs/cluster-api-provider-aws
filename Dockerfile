FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder
WORKDIR /go/src/sigs.k8s.io/cluster-api-provider-aws
COPY . .
# VERSION env gets set in the openshift/release image and refers to the golang version, which interfers with our own
RUN unset VERSION \
 && GOPROXY=off NO_DOCKER=1 make build

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
RUN INSTALL_PKGS=" \
      openssh \
      " && \
    yum install -y $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    yum clean all
COPY --from=builder /go/src/sigs.k8s.io/cluster-api-provider-aws/bin/manager /
COPY --from=builder /go/src/sigs.k8s.io/cluster-api-provider-aws/bin/machine-controller-manager /
COPY --from=builder /go/src/sigs.k8s.io/cluster-api-provider-aws/bin/termination-handler /

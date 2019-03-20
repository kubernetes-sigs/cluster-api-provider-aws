## Overview

This e2e tests provide a mechanism to validate OpenShift/machine-API conformance user stories regardless of the implementation details.

Different testing suites with different primary purposes can be run here, e.g UX expectations, performance expectations, etc. 

This e2e suites **must** run against any repo which **might** break the e2e holistic expected behaviour of the system related in any manner to the machine API domain:

https://github.com/openshift/machine-api-operator

https://github.com/openshift/cluster-autoscaler-operator

https://github.com/openshift/cluster-api-provider-aws

https://github.com/openshift/cluster-api-provider-libvirt

https://github.com/openshift/autoscaler

## Goals

- Validate Openshift/machine-API conformance user stories regardless of the implementation details

- Ensure consistent and reliable behavior of the system related in any manner to the machine API domain

- Provide a reusable solution across any repo which might break the holistic expected behaviour of the system

- Debuggability. Test suites should provide as detailed as possible reasons for the failure in its output

## Non Goals

- Validate behaviour of any component of the system completely unrelated to the machine API domain

- Pre create Openshift/kubernetes clusters to run the expectations against it

## Implementation

- We use [ginkgo](https://onsi.github.io/ginkgo/)

- These tests assume the existence of the system to be validated, i.e an OpenShift cluster and a Kubeconfig file

- Product expectations should be as implementation agnostic as possible e.g "When the workload increases I expect my cluster to grow".
The API Group used for this to happen and the implementation details might vary. The user expectation remains the same

- We are developing a reusable library for manipulating machine API resources that should be leveraged for writing e2e expectations

- Product expectations should be simple, easy to read and avoid details. To this end:
  - They must leverage the common library for manipulating machine API resources so this is done in a consistent manner across all e2e, e.g https://github.com/openshift/cluster-api-actuator-pkg/blob/a086b6d5db86e7e48a6ebce24343e461be233440/pkg/e2e/infra/utils.go#L113
  - Polling logic should be expressed at the Ginkgo level leveraging the common library. E.g https://github.com/openshift/cluster-api-actuator-pkg/blob/a086b6d5db86e7e48a6ebce24343e461be233440/pkg/e2e/infra/infra.go#L184 

## Run it

From your repo:

copy the e2e_test.go file, then:

`dep ensure -v`

And force dep to vendor the required deps as it's just used at dev/CI time, e.g:

```
required = [
  "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/actuators",
  "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/autoscaler",
  "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/infra",
  "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/operators",
  "github.com/openshift/cluster-autoscaler-operator/pkg/apis",
  "github.com/onsi/ginkgo",
  "github.com/onsi/gomega",
  "github.com/golang/glog",
  "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1",
  "k8s.io/client-go/kubernetes/scheme",
  "github.com/openshift/api/config/v1",
]
```

Run it:

```
go test -timeout 30m \
    -v github.com/openshift/cluster-api-actuator-pkg/pkg/e2e \
    -kubeconfig $${KUBECONFIG:-~/.kube/config} \
    -machine-api-namespace $${NAMESPACE:-kube-system} \
    -ginkgo.v \
    -args -v 5 -logtostderr
```
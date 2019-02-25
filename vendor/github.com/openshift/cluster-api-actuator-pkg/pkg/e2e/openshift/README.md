## Overview

This e2e tests provide a mechanism to test end-to-end the behavior of the components of the system (i.e Openshift) which relate in any manner to the machine API domain.

Different testing suites with different primary purposes can be run here, e.g UX expectations, performance expectations, etc. 

This e2e suites **must** run against any repo which **might** break the e2e holistic expected behaviour of the system related in any manner to the machine API domain:

https://github.com/openshift/machine-api-operator

https://github.com/openshift/cluster-autoscaler-operator

https://github.com/openshift/cluster-api-provider-aws

https://github.com/openshift/cluster-api-provider-libvirt

https://github.com/openshift/autoscaler

## Goals

- Ensure consistent and reliable behavior of the system related in any manner to the machine API domain

- Provide a reusable solution across any repo which might break the holistic expected behaviour of the system

- Debuggability. Test suites should provide as detailed as possible reasons for the failure in its output

## Non Goals

- Validate behaviour of any component of the system completely unrelated to the machine API domain

- Pre create Openshift/kubernetes clusters to run the expectations against it

## Implementation

Currently pure golang is used for running suite expectations but this can be moved to a framework eg. ginkgo as requirement rises

This tests assume the existence of the system to be validated, i.e an Openshift cluster and a Kubeconfig file

Product expectations should be as implementation agnostic as possible e.g "When the workload increases I expect my cluster to grow".
The API Group used for this to happen and the implementation details might vary. The user expectation remains the same

## Run it

From your repo:

`dep ensure -add github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/openshift`

And force dep to vendor it as it's just used at dev/CI time:

```
required = [
  "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/openshift"
]
```

Run it:
`go run ./vendor/github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/openshift/*.go -alsologtostderr`
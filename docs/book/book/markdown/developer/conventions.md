# Coding Conventions

Below is a collection of conventions, guidlines and general tips for writing code for this project.

## API Definitions

### Don't Expose 3rd Party Package Types

When adding new or modifying API types don't expose 3rd party package types/enums via the CAPA API definitions. Instead create our own versions and where provide mapping functions.

For example:
* AWS SDK [InstanceState](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/)
* CAPA [InstanceState](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/api/v1beta1/types.go#L560:L581)

### Don't use struct pointer slices

When adding new fields to an API type don't use a slice of struct pointers. This can cause issues with the code generator for the conversion functions. Instead use struct slices.

For example:

Instead of this

```go
	// Configuration options for the non root storage volumes.
	// +optional
	NonRootVolumes []*Volume `json:"nonRootVolumes,omitempty"`
```

use

```go
	// Configuration options for the non root storage volumes.
	// +optional
	NonRootVolumes []Volume `json:"nonRootVolumes,omitempty"`
```

And then within the code you can check the length or range over the slice.

## Tests

There are three types of tests written for CAPA controllers in this repo:
* Unit tests
* Integration tests
* E2E tests

In these tests, we use [fakeclient](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client/fake), [envtest](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest) and [gomock](https://pkg.go.dev/github.com/golang/mock/gomock) libraries based on the requirements of individual test types.

If any new unit, integration or E2E tests has to be added in this repo,we should follow the below conventions.

### Unit tests
These tests are meant to verify the functions inside the same controller file where we perform sanity checks, functionality checks etc.
These tests go into the file with suffix *_unit_test.go.

### Integration tests
These tests are meant to verify the overall flow of the reconcile calls in the controllers to test the flows for all the services/subcomponents of controllers as a whole.
These tests go into the file with suffix *_test.go.

### E2E tests
These tests are meant to verify the proper functioning of a CAPA cluster in an environment that resembles a real production environment. For details, refer [here](https://cluster-api-aws.sigs.k8s.io/development/e2e.html).

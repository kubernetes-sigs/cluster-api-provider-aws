# Coding Conventions

Below is a collection of conventions, guidlines and general tips for writing code for this project. 

## API Definitions

### Don't Expose 3rd Party Package Types

When adding new or modifying API types don't expose 3rd party package types/enums via the CAPA API definitions. Instead create our own versions and where provide mapping functions. 

For example:

    - AWS SDK [InstaneState](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/)
    - CAPA [InstanceState](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/api/v1alpha4/types.go#L560:L581)

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
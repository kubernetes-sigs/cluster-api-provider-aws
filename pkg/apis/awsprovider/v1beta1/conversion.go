package v1beta1

import "sigs.k8s.io/cluster-api-provider-aws/pkg/actuators/machine"

var (
	// RawExtensionFromProviderSpec marshals AWSMachineProviderConfig into a raw extension type.
	// Deprecated: Use machine.RawExtensionFromProviderSpec instead.
	RawExtensionFromProviderSpec = machine.RawExtensionFromProviderSpec
	// RawExtensionFromProviderStatus marshals AWSMachineProviderStatus into a raw extension type.
	// Deprecated: Use machine.RawExtensionFromProviderStatus instead.
	RawExtensionFromProviderStatus = machine.RawExtensionFromProviderStatus
	// ProviderSpecFromRawExtension unmarshals a raw extension into an AWSMachineProviderConfig type.
	// Deprecated: Use machine.ProviderSpecFromRawExtension instead.
	ProviderSpecFromRawExtension = machine.ProviderSpecFromRawExtension
	// ProviderStatusFromRawExtension unmarshals a raw extension into an AWSMachineProviderStatus type.
	// Deprecated: Use machine.ProviderStatusFromRawExtension instead.
	ProviderStatusFromRawExtension = machine.ProviderStatusFromRawExtension
)

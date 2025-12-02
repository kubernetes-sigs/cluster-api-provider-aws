package cloudassemblyschema

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.AmiContextQuery",
		reflect.TypeOf((*AmiContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.ArtifactManifest",
		reflect.TypeOf((*ArtifactManifest)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.cloud_assembly_schema.ArtifactMetadataEntryType",
		reflect.TypeOf((*ArtifactMetadataEntryType)(nil)).Elem(),
		map[string]interface{}{
			"ASSET": ArtifactMetadataEntryType_ASSET,
			"INFO": ArtifactMetadataEntryType_INFO,
			"WARN": ArtifactMetadataEntryType_WARN,
			"ERROR": ArtifactMetadataEntryType_ERROR,
			"LOGICAL_ID": ArtifactMetadataEntryType_LOGICAL_ID,
			"STACK_TAGS": ArtifactMetadataEntryType_STACK_TAGS,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.cloud_assembly_schema.ArtifactType",
		reflect.TypeOf((*ArtifactType)(nil)).Elem(),
		map[string]interface{}{
			"NONE": ArtifactType_NONE,
			"AWS_CLOUDFORMATION_STACK": ArtifactType_AWS_CLOUDFORMATION_STACK,
			"CDK_TREE": ArtifactType_CDK_TREE,
			"ASSET_MANIFEST": ArtifactType_ASSET_MANIFEST,
			"NESTED_CLOUD_ASSEMBLY": ArtifactType_NESTED_CLOUD_ASSEMBLY,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.AssemblyManifest",
		reflect.TypeOf((*AssemblyManifest)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.AssetManifest",
		reflect.TypeOf((*AssetManifest)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.AssetManifestProperties",
		reflect.TypeOf((*AssetManifestProperties)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.AvailabilityZonesContextQuery",
		reflect.TypeOf((*AvailabilityZonesContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.AwsCloudFormationStackProperties",
		reflect.TypeOf((*AwsCloudFormationStackProperties)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.AwsDestination",
		reflect.TypeOf((*AwsDestination)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.ContainerImageAssetMetadataEntry",
		reflect.TypeOf((*ContainerImageAssetMetadataEntry)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.cloud_assembly_schema.ContextProvider",
		reflect.TypeOf((*ContextProvider)(nil)).Elem(),
		map[string]interface{}{
			"AMI_PROVIDER": ContextProvider_AMI_PROVIDER,
			"AVAILABILITY_ZONE_PROVIDER": ContextProvider_AVAILABILITY_ZONE_PROVIDER,
			"HOSTED_ZONE_PROVIDER": ContextProvider_HOSTED_ZONE_PROVIDER,
			"SSM_PARAMETER_PROVIDER": ContextProvider_SSM_PARAMETER_PROVIDER,
			"VPC_PROVIDER": ContextProvider_VPC_PROVIDER,
			"ENDPOINT_SERVICE_AVAILABILITY_ZONE_PROVIDER": ContextProvider_ENDPOINT_SERVICE_AVAILABILITY_ZONE_PROVIDER,
			"LOAD_BALANCER_PROVIDER": ContextProvider_LOAD_BALANCER_PROVIDER,
			"LOAD_BALANCER_LISTENER_PROVIDER": ContextProvider_LOAD_BALANCER_LISTENER_PROVIDER,
			"SECURITY_GROUP_PROVIDER": ContextProvider_SECURITY_GROUP_PROVIDER,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.DockerImageAsset",
		reflect.TypeOf((*DockerImageAsset)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.DockerImageDestination",
		reflect.TypeOf((*DockerImageDestination)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.DockerImageSource",
		reflect.TypeOf((*DockerImageSource)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.EndpointServiceAvailabilityZonesContextQuery",
		reflect.TypeOf((*EndpointServiceAvailabilityZonesContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.FileAsset",
		reflect.TypeOf((*FileAsset)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.FileAssetMetadataEntry",
		reflect.TypeOf((*FileAssetMetadataEntry)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.cloud_assembly_schema.FileAssetPackaging",
		reflect.TypeOf((*FileAssetPackaging)(nil)).Elem(),
		map[string]interface{}{
			"FILE": FileAssetPackaging_FILE,
			"ZIP_DIRECTORY": FileAssetPackaging_ZIP_DIRECTORY,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.FileDestination",
		reflect.TypeOf((*FileDestination)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.FileSource",
		reflect.TypeOf((*FileSource)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.HostedZoneContextQuery",
		reflect.TypeOf((*HostedZoneContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.LoadBalancerContextQuery",
		reflect.TypeOf((*LoadBalancerContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.LoadBalancerFilter",
		reflect.TypeOf((*LoadBalancerFilter)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.LoadBalancerListenerContextQuery",
		reflect.TypeOf((*LoadBalancerListenerContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.cloud_assembly_schema.LoadBalancerListenerProtocol",
		reflect.TypeOf((*LoadBalancerListenerProtocol)(nil)).Elem(),
		map[string]interface{}{
			"HTTP": LoadBalancerListenerProtocol_HTTP,
			"HTTPS": LoadBalancerListenerProtocol_HTTPS,
			"TCP": LoadBalancerListenerProtocol_TCP,
			"TLS": LoadBalancerListenerProtocol_TLS,
			"UDP": LoadBalancerListenerProtocol_UDP,
			"TCP_UDP": LoadBalancerListenerProtocol_TCP_UDP,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.cloud_assembly_schema.LoadBalancerType",
		reflect.TypeOf((*LoadBalancerType)(nil)).Elem(),
		map[string]interface{}{
			"NETWORK": LoadBalancerType_NETWORK,
			"APPLICATION": LoadBalancerType_APPLICATION,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.cloud_assembly_schema.Manifest",
		reflect.TypeOf((*Manifest)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Manifest{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.MetadataEntry",
		reflect.TypeOf((*MetadataEntry)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.MissingContext",
		reflect.TypeOf((*MissingContext)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.NestedCloudAssemblyProperties",
		reflect.TypeOf((*NestedCloudAssemblyProperties)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.RuntimeInfo",
		reflect.TypeOf((*RuntimeInfo)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.SSMParameterContextQuery",
		reflect.TypeOf((*SSMParameterContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.SecurityGroupContextQuery",
		reflect.TypeOf((*SecurityGroupContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.Tag",
		reflect.TypeOf((*Tag)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.TreeArtifactProperties",
		reflect.TypeOf((*TreeArtifactProperties)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cloud_assembly_schema.VpcContextQuery",
		reflect.TypeOf((*VpcContextQuery)(nil)).Elem(),
	)
}

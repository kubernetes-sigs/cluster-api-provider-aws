package cxapi

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"monocdk.cx_api.AssemblyBuildOptions",
		reflect.TypeOf((*AssemblyBuildOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.AssetManifestArtifact",
		reflect.TypeOf((*AssetManifestArtifact)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "assembly", GoGetter: "Assembly"},
			_jsii_.MemberProperty{JsiiProperty: "bootstrapStackVersionSsmParameter", GoGetter: "BootstrapStackVersionSsmParameter"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberProperty{JsiiProperty: "file", GoGetter: "File"},
			_jsii_.MemberMethod{JsiiMethod: "findMetadataByType", GoMethod: "FindMetadataByType"},
			_jsii_.MemberProperty{JsiiProperty: "hierarchicalId", GoGetter: "HierarchicalId"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "manifest", GoGetter: "Manifest"},
			_jsii_.MemberProperty{JsiiProperty: "messages", GoGetter: "Messages"},
			_jsii_.MemberProperty{JsiiProperty: "requiresBootstrapStackVersion", GoGetter: "RequiresBootstrapStackVersion"},
		},
		func() interface{} {
			j := jsiiProxy_AssetManifestArtifact{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CloudArtifact)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.AwsCloudFormationStackProperties",
		reflect.TypeOf((*AwsCloudFormationStackProperties)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.CloudArtifact",
		reflect.TypeOf((*CloudArtifact)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "assembly", GoGetter: "Assembly"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberMethod{JsiiMethod: "findMetadataByType", GoMethod: "FindMetadataByType"},
			_jsii_.MemberProperty{JsiiProperty: "hierarchicalId", GoGetter: "HierarchicalId"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "manifest", GoGetter: "Manifest"},
			_jsii_.MemberProperty{JsiiProperty: "messages", GoGetter: "Messages"},
		},
		func() interface{} {
			return &jsiiProxy_CloudArtifact{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.CloudAssembly",
		reflect.TypeOf((*CloudAssembly)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "artifacts", GoGetter: "Artifacts"},
			_jsii_.MemberProperty{JsiiProperty: "directory", GoGetter: "Directory"},
			_jsii_.MemberMethod{JsiiMethod: "getNestedAssembly", GoMethod: "GetNestedAssembly"},
			_jsii_.MemberMethod{JsiiMethod: "getNestedAssemblyArtifact", GoMethod: "GetNestedAssemblyArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "getStack", GoMethod: "GetStack"},
			_jsii_.MemberMethod{JsiiMethod: "getStackArtifact", GoMethod: "GetStackArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "getStackByName", GoMethod: "GetStackByName"},
			_jsii_.MemberProperty{JsiiProperty: "manifest", GoGetter: "Manifest"},
			_jsii_.MemberProperty{JsiiProperty: "nestedAssemblies", GoGetter: "NestedAssemblies"},
			_jsii_.MemberProperty{JsiiProperty: "runtime", GoGetter: "Runtime"},
			_jsii_.MemberProperty{JsiiProperty: "stacks", GoGetter: "Stacks"},
			_jsii_.MemberProperty{JsiiProperty: "stacksRecursively", GoGetter: "StacksRecursively"},
			_jsii_.MemberMethod{JsiiMethod: "tree", GoMethod: "Tree"},
			_jsii_.MemberMethod{JsiiMethod: "tryGetArtifact", GoMethod: "TryGetArtifact"},
			_jsii_.MemberProperty{JsiiProperty: "version", GoGetter: "Version"},
		},
		func() interface{} {
			return &jsiiProxy_CloudAssembly{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.CloudAssemblyBuilder",
		reflect.TypeOf((*CloudAssemblyBuilder)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addArtifact", GoMethod: "AddArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "addMissing", GoMethod: "AddMissing"},
			_jsii_.MemberProperty{JsiiProperty: "assetOutdir", GoGetter: "AssetOutdir"},
			_jsii_.MemberMethod{JsiiMethod: "buildAssembly", GoMethod: "BuildAssembly"},
			_jsii_.MemberMethod{JsiiMethod: "createNestedAssembly", GoMethod: "CreateNestedAssembly"},
			_jsii_.MemberProperty{JsiiProperty: "outdir", GoGetter: "Outdir"},
		},
		func() interface{} {
			return &jsiiProxy_CloudAssemblyBuilder{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.CloudAssemblyBuilderProps",
		reflect.TypeOf((*CloudAssemblyBuilderProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.CloudFormationStackArtifact",
		reflect.TypeOf((*CloudFormationStackArtifact)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "assembly", GoGetter: "Assembly"},
			_jsii_.MemberProperty{JsiiProperty: "assets", GoGetter: "Assets"},
			_jsii_.MemberProperty{JsiiProperty: "assumeRoleArn", GoGetter: "AssumeRoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "bootstrapStackVersionSsmParameter", GoGetter: "BootstrapStackVersionSsmParameter"},
			_jsii_.MemberProperty{JsiiProperty: "cloudFormationExecutionRoleArn", GoGetter: "CloudFormationExecutionRoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberProperty{JsiiProperty: "displayName", GoGetter: "DisplayName"},
			_jsii_.MemberProperty{JsiiProperty: "environment", GoGetter: "Environment"},
			_jsii_.MemberMethod{JsiiMethod: "findMetadataByType", GoMethod: "FindMetadataByType"},
			_jsii_.MemberProperty{JsiiProperty: "hierarchicalId", GoGetter: "HierarchicalId"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "manifest", GoGetter: "Manifest"},
			_jsii_.MemberProperty{JsiiProperty: "messages", GoGetter: "Messages"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "originalName", GoGetter: "OriginalName"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberProperty{JsiiProperty: "requiresBootstrapStackVersion", GoGetter: "RequiresBootstrapStackVersion"},
			_jsii_.MemberProperty{JsiiProperty: "stackName", GoGetter: "StackName"},
			_jsii_.MemberProperty{JsiiProperty: "stackTemplateAssetObjectUrl", GoGetter: "StackTemplateAssetObjectUrl"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "template", GoGetter: "Template"},
			_jsii_.MemberProperty{JsiiProperty: "templateFile", GoGetter: "TemplateFile"},
			_jsii_.MemberProperty{JsiiProperty: "templateFullPath", GoGetter: "TemplateFullPath"},
			_jsii_.MemberProperty{JsiiProperty: "terminationProtection", GoGetter: "TerminationProtection"},
			_jsii_.MemberProperty{JsiiProperty: "validateOnSynth", GoGetter: "ValidateOnSynth"},
		},
		func() interface{} {
			j := jsiiProxy_CloudFormationStackArtifact{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CloudArtifact)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.EndpointServiceAvailabilityZonesContextQuery",
		reflect.TypeOf((*EndpointServiceAvailabilityZonesContextQuery)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.Environment",
		reflect.TypeOf((*Environment)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.EnvironmentPlaceholderValues",
		reflect.TypeOf((*EnvironmentPlaceholderValues)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.EnvironmentPlaceholders",
		reflect.TypeOf((*EnvironmentPlaceholders)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_EnvironmentPlaceholders{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.EnvironmentUtils",
		reflect.TypeOf((*EnvironmentUtils)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_EnvironmentUtils{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.cx_api.IEnvironmentPlaceholderProvider",
		reflect.TypeOf((*IEnvironmentPlaceholderProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "accountId", GoMethod: "AccountId"},
			_jsii_.MemberMethod{JsiiMethod: "partition", GoMethod: "Partition"},
			_jsii_.MemberMethod{JsiiMethod: "region", GoMethod: "Region"},
		},
		func() interface{} {
			return &jsiiProxy_IEnvironmentPlaceholderProvider{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.LoadBalancerContextResponse",
		reflect.TypeOf((*LoadBalancerContextResponse)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.cx_api.LoadBalancerIpAddressType",
		reflect.TypeOf((*LoadBalancerIpAddressType)(nil)).Elem(),
		map[string]interface{}{
			"IPV4": LoadBalancerIpAddressType_IPV4,
			"DUAL_STACK": LoadBalancerIpAddressType_DUAL_STACK,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.LoadBalancerListenerContextResponse",
		reflect.TypeOf((*LoadBalancerListenerContextResponse)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.MetadataEntry",
		reflect.TypeOf((*MetadataEntry)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.MetadataEntryResult",
		reflect.TypeOf((*MetadataEntryResult)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.MissingContext",
		reflect.TypeOf((*MissingContext)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.NestedCloudAssemblyArtifact",
		reflect.TypeOf((*NestedCloudAssemblyArtifact)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "assembly", GoGetter: "Assembly"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberProperty{JsiiProperty: "directoryName", GoGetter: "DirectoryName"},
			_jsii_.MemberProperty{JsiiProperty: "displayName", GoGetter: "DisplayName"},
			_jsii_.MemberMethod{JsiiMethod: "findMetadataByType", GoMethod: "FindMetadataByType"},
			_jsii_.MemberProperty{JsiiProperty: "fullPath", GoGetter: "FullPath"},
			_jsii_.MemberProperty{JsiiProperty: "hierarchicalId", GoGetter: "HierarchicalId"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "manifest", GoGetter: "Manifest"},
			_jsii_.MemberProperty{JsiiProperty: "messages", GoGetter: "Messages"},
			_jsii_.MemberProperty{JsiiProperty: "nestedAssembly", GoGetter: "NestedAssembly"},
		},
		func() interface{} {
			j := jsiiProxy_NestedCloudAssemblyArtifact{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CloudArtifact)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.RuntimeInfo",
		reflect.TypeOf((*RuntimeInfo)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.SecurityGroupContextResponse",
		reflect.TypeOf((*SecurityGroupContextResponse)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.SynthesisMessage",
		reflect.TypeOf((*SynthesisMessage)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.cx_api.SynthesisMessageLevel",
		reflect.TypeOf((*SynthesisMessageLevel)(nil)).Elem(),
		map[string]interface{}{
			"INFO": SynthesisMessageLevel_INFO,
			"WARNING": SynthesisMessageLevel_WARNING,
			"ERROR": SynthesisMessageLevel_ERROR,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.cx_api.TreeCloudArtifact",
		reflect.TypeOf((*TreeCloudArtifact)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "assembly", GoGetter: "Assembly"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberProperty{JsiiProperty: "file", GoGetter: "File"},
			_jsii_.MemberMethod{JsiiMethod: "findMetadataByType", GoMethod: "FindMetadataByType"},
			_jsii_.MemberProperty{JsiiProperty: "hierarchicalId", GoGetter: "HierarchicalId"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "manifest", GoGetter: "Manifest"},
			_jsii_.MemberProperty{JsiiProperty: "messages", GoGetter: "Messages"},
		},
		func() interface{} {
			j := jsiiProxy_TreeCloudArtifact{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CloudArtifact)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.VpcContextResponse",
		reflect.TypeOf((*VpcContextResponse)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.VpcSubnet",
		reflect.TypeOf((*VpcSubnet)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.cx_api.VpcSubnetGroup",
		reflect.TypeOf((*VpcSubnetGroup)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.cx_api.VpcSubnetGroupType",
		reflect.TypeOf((*VpcSubnetGroupType)(nil)).Elem(),
		map[string]interface{}{
			"PUBLIC": VpcSubnetGroupType_PUBLIC,
			"PRIVATE": VpcSubnetGroupType_PRIVATE,
			"ISOLATED": VpcSubnetGroupType_ISOLATED,
		},
	)
}

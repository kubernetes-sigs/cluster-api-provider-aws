package awscdk

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.Annotations",
		reflect.TypeOf((*Annotations)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeprecation", GoMethod: "AddDeprecation"},
			_jsii_.MemberMethod{JsiiMethod: "addError", GoMethod: "AddError"},
			_jsii_.MemberMethod{JsiiMethod: "addInfo", GoMethod: "AddInfo"},
			_jsii_.MemberMethod{JsiiMethod: "addWarning", GoMethod: "AddWarning"},
		},
		func() interface{} {
			return &jsiiProxy_Annotations{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.App",
		reflect.TypeOf((*App)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "account", GoGetter: "Account"},
			_jsii_.MemberProperty{JsiiProperty: "artifactId", GoGetter: "ArtifactId"},
			_jsii_.MemberProperty{JsiiProperty: "assetOutdir", GoGetter: "AssetOutdir"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outdir", GoGetter: "Outdir"},
			_jsii_.MemberProperty{JsiiProperty: "parentStage", GoGetter: "ParentStage"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "region", GoGetter: "Region"},
			_jsii_.MemberProperty{JsiiProperty: "stageName", GoGetter: "StageName"},
			_jsii_.MemberMethod{JsiiMethod: "synth", GoMethod: "Synth"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_App{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Stage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.AppProps",
		reflect.TypeOf((*AppProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Arn",
		reflect.TypeOf((*Arn)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Arn{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.ArnComponents",
		reflect.TypeOf((*ArnComponents)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.ArnFormat",
		reflect.TypeOf((*ArnFormat)(nil)).Elem(),
		map[string]interface{}{
			"NO_RESOURCE_NAME": ArnFormat_NO_RESOURCE_NAME,
			"COLON_RESOURCE_NAME": ArnFormat_COLON_RESOURCE_NAME,
			"SLASH_RESOURCE_NAME": ArnFormat_SLASH_RESOURCE_NAME,
			"SLASH_RESOURCE_SLASH_RESOURCE_NAME": ArnFormat_SLASH_RESOURCE_SLASH_RESOURCE_NAME,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Aspects",
		reflect.TypeOf((*Aspects)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
			_jsii_.MemberProperty{JsiiProperty: "aspects", GoGetter: "Aspects"},
		},
		func() interface{} {
			return &jsiiProxy_Aspects{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.AssetHashType",
		reflect.TypeOf((*AssetHashType)(nil)).Elem(),
		map[string]interface{}{
			"SOURCE": AssetHashType_SOURCE,
			"BUNDLE": AssetHashType_BUNDLE,
			"OUTPUT": AssetHashType_OUTPUT,
			"CUSTOM": AssetHashType_CUSTOM,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.AssetOptions",
		reflect.TypeOf((*AssetOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.AssetStaging",
		reflect.TypeOf((*AssetStaging)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "absoluteStagedPath", GoGetter: "AbsoluteStagedPath"},
			_jsii_.MemberProperty{JsiiProperty: "assetHash", GoGetter: "AssetHash"},
			_jsii_.MemberProperty{JsiiProperty: "isArchive", GoGetter: "IsArchive"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "packaging", GoGetter: "Packaging"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "relativeStagedPath", GoMethod: "RelativeStagedPath"},
			_jsii_.MemberProperty{JsiiProperty: "sourceHash", GoGetter: "SourceHash"},
			_jsii_.MemberProperty{JsiiProperty: "sourcePath", GoGetter: "SourcePath"},
			_jsii_.MemberProperty{JsiiProperty: "stagedPath", GoGetter: "StagedPath"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_AssetStaging{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Construct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.AssetStagingProps",
		reflect.TypeOf((*AssetStagingProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Aws",
		reflect.TypeOf((*Aws)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Aws{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.BootstraplessSynthesizer",
		reflect.TypeOf((*BootstraplessSynthesizer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "cloudFormationExecutionRoleArn", GoGetter: "CloudFormationExecutionRoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "deployRoleArn", GoGetter: "DeployRoleArn"},
			_jsii_.MemberMethod{JsiiMethod: "emitStackArtifact", GoMethod: "EmitStackArtifact"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeStackTemplate", GoMethod: "SynthesizeStackTemplate"},
		},
		func() interface{} {
			j := jsiiProxy_BootstraplessSynthesizer{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_DefaultStackSynthesizer)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.BootstraplessSynthesizerProps",
		reflect.TypeOf((*BootstraplessSynthesizerProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.BundlingDockerImage",
		reflect.TypeOf((*BundlingDockerImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "cp", GoMethod: "Cp"},
			_jsii_.MemberProperty{JsiiProperty: "image", GoGetter: "Image"},
			_jsii_.MemberMethod{JsiiMethod: "run", GoMethod: "Run"},
			_jsii_.MemberMethod{JsiiMethod: "toJSON", GoMethod: "ToJSON"},
		},
		func() interface{} {
			return &jsiiProxy_BundlingDockerImage{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.BundlingOptions",
		reflect.TypeOf((*BundlingOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.BundlingOutput",
		reflect.TypeOf((*BundlingOutput)(nil)).Elem(),
		map[string]interface{}{
			"ARCHIVED": BundlingOutput_ARCHIVED,
			"NOT_ARCHIVED": BundlingOutput_NOT_ARCHIVED,
			"AUTO_DISCOVER": BundlingOutput_AUTO_DISCOVER,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnAutoScalingReplacingUpdate",
		reflect.TypeOf((*CfnAutoScalingReplacingUpdate)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnAutoScalingRollingUpdate",
		reflect.TypeOf((*CfnAutoScalingRollingUpdate)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnAutoScalingScheduledAction",
		reflect.TypeOf((*CfnAutoScalingScheduledAction)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.CfnCapabilities",
		reflect.TypeOf((*CfnCapabilities)(nil)).Elem(),
		map[string]interface{}{
			"NONE": CfnCapabilities_NONE,
			"ANONYMOUS_IAM": CfnCapabilities_ANONYMOUS_IAM,
			"NAMED_IAM": CfnCapabilities_NAMED_IAM,
			"AUTO_EXPAND": CfnCapabilities_AUTO_EXPAND,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCodeDeployBlueGreenAdditionalOptions",
		reflect.TypeOf((*CfnCodeDeployBlueGreenAdditionalOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCodeDeployBlueGreenApplication",
		reflect.TypeOf((*CfnCodeDeployBlueGreenApplication)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCodeDeployBlueGreenApplicationTarget",
		reflect.TypeOf((*CfnCodeDeployBlueGreenApplicationTarget)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCodeDeployBlueGreenEcsAttributes",
		reflect.TypeOf((*CfnCodeDeployBlueGreenEcsAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnCodeDeployBlueGreenHook",
		reflect.TypeOf((*CfnCodeDeployBlueGreenHook)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "additionalOptions", GoGetter: "AdditionalOptions"},
			_jsii_.MemberProperty{JsiiProperty: "applications", GoGetter: "Applications"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "lifecycleEventHooks", GoGetter: "LifecycleEventHooks"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "serviceRole", GoGetter: "ServiceRole"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "trafficRoutingConfig", GoGetter: "TrafficRoutingConfig"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnCodeDeployBlueGreenHook{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnHook)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCodeDeployBlueGreenHookProps",
		reflect.TypeOf((*CfnCodeDeployBlueGreenHookProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCodeDeployBlueGreenLifecycleEventHooks",
		reflect.TypeOf((*CfnCodeDeployBlueGreenLifecycleEventHooks)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCodeDeployLambdaAliasUpdate",
		reflect.TypeOf((*CfnCodeDeployLambdaAliasUpdate)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnCondition",
		reflect.TypeOf((*CfnCondition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "expression", GoGetter: "Expression"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnCondition{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnElement)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ICfnConditionExpression)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IResolvable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnConditionProps",
		reflect.TypeOf((*CfnConditionProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCreationPolicy",
		reflect.TypeOf((*CfnCreationPolicy)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnCustomResource",
		reflect.TypeOf((*CfnCustomResource)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "serviceToken", GoGetter: "ServiceToken"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnCustomResource{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnCustomResourceProps",
		reflect.TypeOf((*CfnCustomResourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.CfnDeletionPolicy",
		reflect.TypeOf((*CfnDeletionPolicy)(nil)).Elem(),
		map[string]interface{}{
			"DELETE": CfnDeletionPolicy_DELETE,
			"RETAIN": CfnDeletionPolicy_RETAIN,
			"SNAPSHOT": CfnDeletionPolicy_SNAPSHOT,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.CfnDynamicReference",
		reflect.TypeOf((*CfnDynamicReference)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "newError", GoMethod: "NewError"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "toJSON", GoMethod: "ToJSON"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_CfnDynamicReference{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Intrinsic)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnDynamicReferenceProps",
		reflect.TypeOf((*CfnDynamicReferenceProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.CfnDynamicReferenceService",
		reflect.TypeOf((*CfnDynamicReferenceService)(nil)).Elem(),
		map[string]interface{}{
			"SSM": CfnDynamicReferenceService_SSM,
			"SSM_SECURE": CfnDynamicReferenceService_SSM_SECURE,
			"SECRETS_MANAGER": CfnDynamicReferenceService_SECRETS_MANAGER,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.CfnElement",
		reflect.TypeOf((*CfnElement)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnElement{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Construct)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.CfnHook",
		reflect.TypeOf((*CfnHook)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnHook{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnElement)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnHookProps",
		reflect.TypeOf((*CfnHookProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnInclude",
		reflect.TypeOf((*CfnInclude)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "template", GoGetter: "Template"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnInclude{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnElement)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnIncludeProps",
		reflect.TypeOf((*CfnIncludeProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnJson",
		reflect.TypeOf((*CfnJson)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toJSON", GoMethod: "ToJSON"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
		},
		func() interface{} {
			j := jsiiProxy_CfnJson{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Construct)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IResolvable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnJsonProps",
		reflect.TypeOf((*CfnJsonProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnMacro",
		reflect.TypeOf((*CfnMacro)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "functionName", GoGetter: "FunctionName"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logGroupName", GoGetter: "LogGroupName"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "logRoleArn", GoGetter: "LogRoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnMacro{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnMacroProps",
		reflect.TypeOf((*CfnMacroProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnMapping",
		reflect.TypeOf((*CfnMapping)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "findInMap", GoMethod: "FindInMap"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "setValue", GoMethod: "SetValue"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnMapping{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnRefElement)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnMappingProps",
		reflect.TypeOf((*CfnMappingProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnModuleDefaultVersion",
		reflect.TypeOf((*CfnModuleDefaultVersion)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "arn", GoGetter: "Arn"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "moduleName", GoGetter: "ModuleName"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "versionId", GoGetter: "VersionId"},
		},
		func() interface{} {
			j := jsiiProxy_CfnModuleDefaultVersion{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnModuleDefaultVersionProps",
		reflect.TypeOf((*CfnModuleDefaultVersionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnModuleVersion",
		reflect.TypeOf((*CfnModuleVersion)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrDescription", GoGetter: "AttrDescription"},
			_jsii_.MemberProperty{JsiiProperty: "attrDocumentationUrl", GoGetter: "AttrDocumentationUrl"},
			_jsii_.MemberProperty{JsiiProperty: "attrIsDefaultVersion", GoGetter: "AttrIsDefaultVersion"},
			_jsii_.MemberProperty{JsiiProperty: "attrSchema", GoGetter: "AttrSchema"},
			_jsii_.MemberProperty{JsiiProperty: "attrTimeCreated", GoGetter: "AttrTimeCreated"},
			_jsii_.MemberProperty{JsiiProperty: "attrVersionId", GoGetter: "AttrVersionId"},
			_jsii_.MemberProperty{JsiiProperty: "attrVisibility", GoGetter: "AttrVisibility"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "moduleName", GoGetter: "ModuleName"},
			_jsii_.MemberProperty{JsiiProperty: "modulePackage", GoGetter: "ModulePackage"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnModuleVersion{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnModuleVersionProps",
		reflect.TypeOf((*CfnModuleVersionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnOutput",
		reflect.TypeOf((*CfnOutput)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "condition", GoGetter: "Condition"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "exportName", GoGetter: "ExportName"},
			_jsii_.MemberProperty{JsiiProperty: "importValue", GoGetter: "ImportValue"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
		},
		func() interface{} {
			j := jsiiProxy_CfnOutput{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnElement)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnOutputProps",
		reflect.TypeOf((*CfnOutputProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnParameter",
		reflect.TypeOf((*CfnParameter)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "allowedPattern", GoGetter: "AllowedPattern"},
			_jsii_.MemberProperty{JsiiProperty: "allowedValues", GoGetter: "AllowedValues"},
			_jsii_.MemberProperty{JsiiProperty: "constraintDescription", GoGetter: "ConstraintDescription"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "default", GoGetter: "Default"},
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "maxLength", GoGetter: "MaxLength"},
			_jsii_.MemberProperty{JsiiProperty: "maxValue", GoGetter: "MaxValue"},
			_jsii_.MemberProperty{JsiiProperty: "minLength", GoGetter: "MinLength"},
			_jsii_.MemberProperty{JsiiProperty: "minValue", GoGetter: "MinValue"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "noEcho", GoGetter: "NoEcho"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
			_jsii_.MemberProperty{JsiiProperty: "valueAsList", GoGetter: "ValueAsList"},
			_jsii_.MemberProperty{JsiiProperty: "valueAsNumber", GoGetter: "ValueAsNumber"},
			_jsii_.MemberProperty{JsiiProperty: "valueAsString", GoGetter: "ValueAsString"},
		},
		func() interface{} {
			j := jsiiProxy_CfnParameter{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnElement)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnParameterProps",
		reflect.TypeOf((*CfnParameterProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnRefElement",
		reflect.TypeOf((*CfnRefElement)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnRefElement{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnElement)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.CfnResource",
		reflect.TypeOf((*CfnResource)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnResource{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnRefElement)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnResourceAutoScalingCreationPolicy",
		reflect.TypeOf((*CfnResourceAutoScalingCreationPolicy)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnResourceDefaultVersion",
		reflect.TypeOf((*CfnResourceDefaultVersion)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "typeName", GoGetter: "TypeName"},
			_jsii_.MemberProperty{JsiiProperty: "typeVersionArn", GoGetter: "TypeVersionArn"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "versionId", GoGetter: "VersionId"},
		},
		func() interface{} {
			j := jsiiProxy_CfnResourceDefaultVersion{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnResourceDefaultVersionProps",
		reflect.TypeOf((*CfnResourceDefaultVersionProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnResourceProps",
		reflect.TypeOf((*CfnResourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnResourceSignal",
		reflect.TypeOf((*CfnResourceSignal)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnResourceVersion",
		reflect.TypeOf((*CfnResourceVersion)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrIsDefaultVersion", GoGetter: "AttrIsDefaultVersion"},
			_jsii_.MemberProperty{JsiiProperty: "attrProvisioningType", GoGetter: "AttrProvisioningType"},
			_jsii_.MemberProperty{JsiiProperty: "attrTypeArn", GoGetter: "AttrTypeArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrVersionId", GoGetter: "AttrVersionId"},
			_jsii_.MemberProperty{JsiiProperty: "attrVisibility", GoGetter: "AttrVisibility"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "executionRoleArn", GoGetter: "ExecutionRoleArn"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "loggingConfig", GoGetter: "LoggingConfig"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "schemaHandlerPackage", GoGetter: "SchemaHandlerPackage"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "typeName", GoGetter: "TypeName"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnResourceVersion{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnResourceVersion.LoggingConfigProperty",
		reflect.TypeOf((*CfnResourceVersion_LoggingConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnResourceVersionProps",
		reflect.TypeOf((*CfnResourceVersionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnRule",
		reflect.TypeOf((*CfnRule)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addAssertion", GoMethod: "AddAssertion"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnRule{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnRefElement)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnRuleAssertion",
		reflect.TypeOf((*CfnRuleAssertion)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnRuleProps",
		reflect.TypeOf((*CfnRuleProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnStack",
		reflect.TypeOf((*CfnStack)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "notificationArns", GoGetter: "NotificationArns"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "templateUrl", GoGetter: "TemplateUrl"},
			_jsii_.MemberProperty{JsiiProperty: "timeoutInMinutes", GoGetter: "TimeoutInMinutes"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnStack{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnStackProps",
		reflect.TypeOf((*CfnStackProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnStackSet",
		reflect.TypeOf((*CfnStackSet)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberProperty{JsiiProperty: "administrationRoleArn", GoGetter: "AdministrationRoleArn"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrStackSetId", GoGetter: "AttrStackSetId"},
			_jsii_.MemberProperty{JsiiProperty: "autoDeployment", GoGetter: "AutoDeployment"},
			_jsii_.MemberProperty{JsiiProperty: "callAs", GoGetter: "CallAs"},
			_jsii_.MemberProperty{JsiiProperty: "capabilities", GoGetter: "Capabilities"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "executionRoleName", GoGetter: "ExecutionRoleName"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "operationPreferences", GoGetter: "OperationPreferences"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberProperty{JsiiProperty: "permissionModel", GoGetter: "PermissionModel"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "stackInstancesGroup", GoGetter: "StackInstancesGroup"},
			_jsii_.MemberProperty{JsiiProperty: "stackSetName", GoGetter: "StackSetName"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "templateBody", GoGetter: "TemplateBody"},
			_jsii_.MemberProperty{JsiiProperty: "templateUrl", GoGetter: "TemplateUrl"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnStackSet{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnStackSet.AutoDeploymentProperty",
		reflect.TypeOf((*CfnStackSet_AutoDeploymentProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnStackSet.DeploymentTargetsProperty",
		reflect.TypeOf((*CfnStackSet_DeploymentTargetsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnStackSet.OperationPreferencesProperty",
		reflect.TypeOf((*CfnStackSet_OperationPreferencesProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnStackSet.ParameterProperty",
		reflect.TypeOf((*CfnStackSet_ParameterProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnStackSet.StackInstancesProperty",
		reflect.TypeOf((*CfnStackSet_StackInstancesProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnStackSetProps",
		reflect.TypeOf((*CfnStackSetProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnTag",
		reflect.TypeOf((*CfnTag)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnTrafficRoute",
		reflect.TypeOf((*CfnTrafficRoute)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnTrafficRouting",
		reflect.TypeOf((*CfnTrafficRouting)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnTrafficRoutingConfig",
		reflect.TypeOf((*CfnTrafficRoutingConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnTrafficRoutingTimeBasedCanary",
		reflect.TypeOf((*CfnTrafficRoutingTimeBasedCanary)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnTrafficRoutingTimeBasedLinear",
		reflect.TypeOf((*CfnTrafficRoutingTimeBasedLinear)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.CfnTrafficRoutingType",
		reflect.TypeOf((*CfnTrafficRoutingType)(nil)).Elem(),
		map[string]interface{}{
			"ALL_AT_ONCE": CfnTrafficRoutingType_ALL_AT_ONCE,
			"TIME_BASED_CANARY": CfnTrafficRoutingType_TIME_BASED_CANARY,
			"TIME_BASED_LINEAR": CfnTrafficRoutingType_TIME_BASED_LINEAR,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnUpdatePolicy",
		reflect.TypeOf((*CfnUpdatePolicy)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CfnWaitCondition",
		reflect.TypeOf((*CfnWaitCondition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrData", GoGetter: "AttrData"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "count", GoGetter: "Count"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "handle", GoGetter: "Handle"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "timeout", GoGetter: "Timeout"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnWaitCondition{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.CfnWaitConditionHandle",
		reflect.TypeOf((*CfnWaitConditionHandle)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnWaitConditionHandle{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_CfnResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CfnWaitConditionProps",
		reflect.TypeOf((*CfnWaitConditionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.ConcreteDependable",
		reflect.TypeOf((*ConcreteDependable)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
		},
		func() interface{} {
			j := jsiiProxy_ConcreteDependable{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IDependable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Construct",
		reflect.TypeOf((*Construct)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Construct{}
			_jsii_.InitJsiiProxy(&j.Type__constructsConstruct)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IConstruct)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.ConstructNode",
		reflect.TypeOf((*ConstructNode)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "addError", GoMethod: "AddError"},
			_jsii_.MemberMethod{JsiiMethod: "addInfo", GoMethod: "AddInfo"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "addr", GoGetter: "Addr"},
			_jsii_.MemberMethod{JsiiMethod: "addValidation", GoMethod: "AddValidation"},
			_jsii_.MemberMethod{JsiiMethod: "addWarning", GoMethod: "AddWarning"},
			_jsii_.MemberMethod{JsiiMethod: "applyAspect", GoMethod: "ApplyAspect"},
			_jsii_.MemberProperty{JsiiProperty: "children", GoGetter: "Children"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChild", GoGetter: "DefaultChild"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberMethod{JsiiMethod: "findAll", GoMethod: "FindAll"},
			_jsii_.MemberMethod{JsiiMethod: "findChild", GoMethod: "FindChild"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "locked", GoGetter: "Locked"},
			_jsii_.MemberProperty{JsiiProperty: "metadata", GoGetter: "Metadata"},
			_jsii_.MemberProperty{JsiiProperty: "metadataEntry", GoGetter: "MetadataEntry"},
			_jsii_.MemberProperty{JsiiProperty: "path", GoGetter: "Path"},
			_jsii_.MemberProperty{JsiiProperty: "root", GoGetter: "Root"},
			_jsii_.MemberProperty{JsiiProperty: "scope", GoGetter: "Scope"},
			_jsii_.MemberProperty{JsiiProperty: "scopes", GoGetter: "Scopes"},
			_jsii_.MemberMethod{JsiiMethod: "setContext", GoMethod: "SetContext"},
			_jsii_.MemberMethod{JsiiMethod: "tryFindChild", GoMethod: "TryFindChild"},
			_jsii_.MemberMethod{JsiiMethod: "tryGetContext", GoMethod: "TryGetContext"},
			_jsii_.MemberMethod{JsiiMethod: "tryRemoveChild", GoMethod: "TryRemoveChild"},
			_jsii_.MemberProperty{JsiiProperty: "uniqueId", GoGetter: "UniqueId"},
		},
		func() interface{} {
			return &jsiiProxy_ConstructNode{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.ConstructOrder",
		reflect.TypeOf((*ConstructOrder)(nil)).Elem(),
		map[string]interface{}{
			"PREORDER": ConstructOrder_PREORDER,
			"POSTORDER": ConstructOrder_POSTORDER,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.ContextProvider",
		reflect.TypeOf((*ContextProvider)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_ContextProvider{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CopyOptions",
		reflect.TypeOf((*CopyOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CustomResource",
		reflect.TypeOf((*CustomResource)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getAttString", GoMethod: "GetAttString"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CustomResource{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Resource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CustomResourceProps",
		reflect.TypeOf((*CustomResourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.CustomResourceProvider",
		reflect.TypeOf((*CustomResourceProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "roleArn", GoGetter: "RoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceToken", GoGetter: "ServiceToken"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_CustomResourceProvider{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Construct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.CustomResourceProviderProps",
		reflect.TypeOf((*CustomResourceProviderProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.CustomResourceProviderRuntime",
		reflect.TypeOf((*CustomResourceProviderRuntime)(nil)).Elem(),
		map[string]interface{}{
			"NODEJS_12": CustomResourceProviderRuntime_NODEJS_12,
			"NODEJS_14_X": CustomResourceProviderRuntime_NODEJS_14_X,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.DefaultStackSynthesizer",
		reflect.TypeOf((*DefaultStackSynthesizer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "cloudFormationExecutionRoleArn", GoGetter: "CloudFormationExecutionRoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "deployRoleArn", GoGetter: "DeployRoleArn"},
			_jsii_.MemberMethod{JsiiMethod: "emitStackArtifact", GoMethod: "EmitStackArtifact"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeStackTemplate", GoMethod: "SynthesizeStackTemplate"},
		},
		func() interface{} {
			j := jsiiProxy_DefaultStackSynthesizer{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_StackSynthesizer)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.DefaultStackSynthesizerProps",
		reflect.TypeOf((*DefaultStackSynthesizerProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.DefaultTokenResolver",
		reflect.TypeOf((*DefaultTokenResolver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "resolveList", GoMethod: "ResolveList"},
			_jsii_.MemberMethod{JsiiMethod: "resolveString", GoMethod: "ResolveString"},
			_jsii_.MemberMethod{JsiiMethod: "resolveToken", GoMethod: "ResolveToken"},
		},
		func() interface{} {
			j := jsiiProxy_DefaultTokenResolver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ITokenResolver)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.DependableTrait",
		reflect.TypeOf((*DependableTrait)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "dependencyRoots", GoGetter: "DependencyRoots"},
		},
		func() interface{} {
			return &jsiiProxy_DependableTrait{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.Dependency",
		reflect.TypeOf((*Dependency)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.DockerBuildOptions",
		reflect.TypeOf((*DockerBuildOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.DockerIgnoreStrategy",
		reflect.TypeOf((*DockerIgnoreStrategy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
			_jsii_.MemberMethod{JsiiMethod: "ignores", GoMethod: "Ignores"},
		},
		func() interface{} {
			j := jsiiProxy_DockerIgnoreStrategy{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IgnoreStrategy)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.DockerImage",
		reflect.TypeOf((*DockerImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "cp", GoMethod: "Cp"},
			_jsii_.MemberProperty{JsiiProperty: "image", GoGetter: "Image"},
			_jsii_.MemberMethod{JsiiMethod: "run", GoMethod: "Run"},
			_jsii_.MemberMethod{JsiiMethod: "toJSON", GoMethod: "ToJSON"},
		},
		func() interface{} {
			j := jsiiProxy_DockerImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_BundlingDockerImage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.DockerImageAssetLocation",
		reflect.TypeOf((*DockerImageAssetLocation)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.DockerImageAssetSource",
		reflect.TypeOf((*DockerImageAssetSource)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.DockerRunOptions",
		reflect.TypeOf((*DockerRunOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.DockerVolume",
		reflect.TypeOf((*DockerVolume)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.DockerVolumeConsistency",
		reflect.TypeOf((*DockerVolumeConsistency)(nil)).Elem(),
		map[string]interface{}{
			"CONSISTENT": DockerVolumeConsistency_CONSISTENT,
			"DELEGATED": DockerVolumeConsistency_DELEGATED,
			"CACHED": DockerVolumeConsistency_CACHED,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Duration",
		reflect.TypeOf((*Duration)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "formatTokenToNumber", GoMethod: "FormatTokenToNumber"},
			_jsii_.MemberMethod{JsiiMethod: "isUnresolved", GoMethod: "IsUnresolved"},
			_jsii_.MemberMethod{JsiiMethod: "plus", GoMethod: "Plus"},
			_jsii_.MemberMethod{JsiiMethod: "toDays", GoMethod: "ToDays"},
			_jsii_.MemberMethod{JsiiMethod: "toHours", GoMethod: "ToHours"},
			_jsii_.MemberMethod{JsiiMethod: "toHumanString", GoMethod: "ToHumanString"},
			_jsii_.MemberMethod{JsiiMethod: "toIsoString", GoMethod: "ToIsoString"},
			_jsii_.MemberMethod{JsiiMethod: "toISOString", GoMethod: "ToISOString"},
			_jsii_.MemberMethod{JsiiMethod: "toMilliseconds", GoMethod: "ToMilliseconds"},
			_jsii_.MemberMethod{JsiiMethod: "toMinutes", GoMethod: "ToMinutes"},
			_jsii_.MemberMethod{JsiiMethod: "toSeconds", GoMethod: "ToSeconds"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "unitLabel", GoMethod: "UnitLabel"},
		},
		func() interface{} {
			return &jsiiProxy_Duration{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.EncodingOptions",
		reflect.TypeOf((*EncodingOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.Environment",
		reflect.TypeOf((*Environment)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Expiration",
		reflect.TypeOf((*Expiration)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "date", GoGetter: "Date"},
			_jsii_.MemberMethod{JsiiMethod: "isAfter", GoMethod: "IsAfter"},
			_jsii_.MemberMethod{JsiiMethod: "isBefore", GoMethod: "IsBefore"},
			_jsii_.MemberMethod{JsiiMethod: "toEpoch", GoMethod: "ToEpoch"},
		},
		func() interface{} {
			return &jsiiProxy_Expiration{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.ExportValueOptions",
		reflect.TypeOf((*ExportValueOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.FeatureFlags",
		reflect.TypeOf((*FeatureFlags)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "isEnabled", GoMethod: "IsEnabled"},
		},
		func() interface{} {
			return &jsiiProxy_FeatureFlags{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.FileAssetLocation",
		reflect.TypeOf((*FileAssetLocation)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.FileAssetPackaging",
		reflect.TypeOf((*FileAssetPackaging)(nil)).Elem(),
		map[string]interface{}{
			"ZIP_DIRECTORY": FileAssetPackaging_ZIP_DIRECTORY,
			"FILE": FileAssetPackaging_FILE,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.FileAssetSource",
		reflect.TypeOf((*FileAssetSource)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.FileCopyOptions",
		reflect.TypeOf((*FileCopyOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.FileFingerprintOptions",
		reflect.TypeOf((*FileFingerprintOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.FileSystem",
		reflect.TypeOf((*FileSystem)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_FileSystem{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.FingerprintOptions",
		reflect.TypeOf((*FingerprintOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Fn",
		reflect.TypeOf((*Fn)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Fn{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.GetContextKeyOptions",
		reflect.TypeOf((*GetContextKeyOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.GetContextKeyResult",
		reflect.TypeOf((*GetContextKeyResult)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.GetContextValueOptions",
		reflect.TypeOf((*GetContextValueOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.GetContextValueResult",
		reflect.TypeOf((*GetContextValueResult)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.GitIgnoreStrategy",
		reflect.TypeOf((*GitIgnoreStrategy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
			_jsii_.MemberMethod{JsiiMethod: "ignores", GoMethod: "Ignores"},
		},
		func() interface{} {
			j := jsiiProxy_GitIgnoreStrategy{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IgnoreStrategy)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.GlobIgnoreStrategy",
		reflect.TypeOf((*GlobIgnoreStrategy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
			_jsii_.MemberMethod{JsiiMethod: "ignores", GoMethod: "Ignores"},
		},
		func() interface{} {
			j := jsiiProxy_GlobIgnoreStrategy{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IgnoreStrategy)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IAnyProducer",
		reflect.TypeOf((*IAnyProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IAnyProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IAspect",
		reflect.TypeOf((*IAspect)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "visit", GoMethod: "Visit"},
		},
		func() interface{} {
			return &jsiiProxy_IAspect{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IAsset",
		reflect.TypeOf((*IAsset)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "assetHash", GoGetter: "AssetHash"},
		},
		func() interface{} {
			return &jsiiProxy_IAsset{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ICfnConditionExpression",
		reflect.TypeOf((*ICfnConditionExpression)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_ICfnConditionExpression{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IResolvable)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ICfnResourceOptions",
		reflect.TypeOf((*ICfnResourceOptions)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "condition", GoGetter: "Condition"},
			_jsii_.MemberProperty{JsiiProperty: "creationPolicy", GoGetter: "CreationPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "deletionPolicy", GoGetter: "DeletionPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "metadata", GoGetter: "Metadata"},
			_jsii_.MemberProperty{JsiiProperty: "updatePolicy", GoGetter: "UpdatePolicy"},
			_jsii_.MemberProperty{JsiiProperty: "updateReplacePolicy", GoGetter: "UpdateReplacePolicy"},
			_jsii_.MemberProperty{JsiiProperty: "version", GoGetter: "Version"},
		},
		func() interface{} {
			return &jsiiProxy_ICfnResourceOptions{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IConstruct",
		reflect.TypeOf((*IConstruct)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
		},
		func() interface{} {
			j := jsiiProxy_IConstruct{}
			_jsii_.InitJsiiProxy(&j.Type__constructsIConstruct)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IDependable)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IDependable",
		reflect.TypeOf((*IDependable)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_IDependable{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IFragmentConcatenator",
		reflect.TypeOf((*IFragmentConcatenator)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "join", GoMethod: "Join"},
		},
		func() interface{} {
			return &jsiiProxy_IFragmentConcatenator{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IInspectable",
		reflect.TypeOf((*IInspectable)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
		},
		func() interface{} {
			return &jsiiProxy_IInspectable{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IListProducer",
		reflect.TypeOf((*IListProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IListProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ILocalBundling",
		reflect.TypeOf((*ILocalBundling)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "tryBundle", GoMethod: "TryBundle"},
		},
		func() interface{} {
			return &jsiiProxy_ILocalBundling{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.INumberProducer",
		reflect.TypeOf((*INumberProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_INumberProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IPostProcessor",
		reflect.TypeOf((*IPostProcessor)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "postProcess", GoMethod: "PostProcess"},
		},
		func() interface{} {
			return &jsiiProxy_IPostProcessor{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IResolvable",
		reflect.TypeOf((*IResolvable)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			return &jsiiProxy_IResolvable{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IResolveContext",
		reflect.TypeOf((*IResolveContext)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "preparing", GoGetter: "Preparing"},
			_jsii_.MemberMethod{JsiiMethod: "registerPostProcessor", GoMethod: "RegisterPostProcessor"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberProperty{JsiiProperty: "scope", GoGetter: "Scope"},
		},
		func() interface{} {
			return &jsiiProxy_IResolveContext{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IResource",
		reflect.TypeOf((*IResource)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IResource{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IConstruct)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IStableAnyProducer",
		reflect.TypeOf((*IStableAnyProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IStableAnyProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IStableListProducer",
		reflect.TypeOf((*IStableListProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IStableListProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IStableNumberProducer",
		reflect.TypeOf((*IStableNumberProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IStableNumberProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IStableStringProducer",
		reflect.TypeOf((*IStableStringProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IStableStringProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IStackSynthesizer",
		reflect.TypeOf((*IStackSynthesizer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
		},
		func() interface{} {
			return &jsiiProxy_IStackSynthesizer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.IStringProducer",
		reflect.TypeOf((*IStringProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IStringProducer{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ISynthesisSession",
		reflect.TypeOf((*ISynthesisSession)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "assembly", GoGetter: "Assembly"},
			_jsii_.MemberProperty{JsiiProperty: "outdir", GoGetter: "Outdir"},
			_jsii_.MemberProperty{JsiiProperty: "validateOnSynth", GoGetter: "ValidateOnSynth"},
		},
		func() interface{} {
			return &jsiiProxy_ISynthesisSession{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ITaggable",
		reflect.TypeOf((*ITaggable)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
		},
		func() interface{} {
			return &jsiiProxy_ITaggable{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ITemplateOptions",
		reflect.TypeOf((*ITemplateOptions)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "metadata", GoGetter: "Metadata"},
			_jsii_.MemberProperty{JsiiProperty: "templateFormatVersion", GoGetter: "TemplateFormatVersion"},
			_jsii_.MemberProperty{JsiiProperty: "transform", GoGetter: "Transform"},
			_jsii_.MemberProperty{JsiiProperty: "transforms", GoGetter: "Transforms"},
		},
		func() interface{} {
			return &jsiiProxy_ITemplateOptions{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ITokenMapper",
		reflect.TypeOf((*ITokenMapper)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "mapToken", GoMethod: "MapToken"},
		},
		func() interface{} {
			return &jsiiProxy_ITokenMapper{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.ITokenResolver",
		reflect.TypeOf((*ITokenResolver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "resolveList", GoMethod: "ResolveList"},
			_jsii_.MemberMethod{JsiiMethod: "resolveString", GoMethod: "ResolveString"},
			_jsii_.MemberMethod{JsiiMethod: "resolveToken", GoMethod: "ResolveToken"},
		},
		func() interface{} {
			return &jsiiProxy_ITokenResolver{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.IgnoreMode",
		reflect.TypeOf((*IgnoreMode)(nil)).Elem(),
		map[string]interface{}{
			"GLOB": IgnoreMode_GLOB,
			"GIT": IgnoreMode_GIT,
			"DOCKER": IgnoreMode_DOCKER,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.IgnoreStrategy",
		reflect.TypeOf((*IgnoreStrategy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
			_jsii_.MemberMethod{JsiiMethod: "ignores", GoMethod: "Ignores"},
		},
		func() interface{} {
			return &jsiiProxy_IgnoreStrategy{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Intrinsic",
		reflect.TypeOf((*Intrinsic)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "newError", GoMethod: "NewError"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "toJSON", GoMethod: "ToJSON"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_Intrinsic{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IResolvable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.IntrinsicProps",
		reflect.TypeOf((*IntrinsicProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Lazy",
		reflect.TypeOf((*Lazy)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Lazy{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.LazyAnyValueOptions",
		reflect.TypeOf((*LazyAnyValueOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.LazyListValueOptions",
		reflect.TypeOf((*LazyListValueOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.LazyStringValueOptions",
		reflect.TypeOf((*LazyStringValueOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.LegacyStackSynthesizer",
		reflect.TypeOf((*LegacyStackSynthesizer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberMethod{JsiiMethod: "emitStackArtifact", GoMethod: "EmitStackArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeStackTemplate", GoMethod: "SynthesizeStackTemplate"},
		},
		func() interface{} {
			j := jsiiProxy_LegacyStackSynthesizer{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_StackSynthesizer)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Names",
		reflect.TypeOf((*Names)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Names{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.NestedStack",
		reflect.TypeOf((*NestedStack)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "account", GoGetter: "Account"},
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addTransform", GoMethod: "AddTransform"},
			_jsii_.MemberMethod{JsiiMethod: "allocateLogicalId", GoMethod: "AllocateLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "artifactId", GoGetter: "ArtifactId"},
			_jsii_.MemberProperty{JsiiProperty: "availabilityZones", GoGetter: "AvailabilityZones"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberProperty{JsiiProperty: "environment", GoGetter: "Environment"},
			_jsii_.MemberMethod{JsiiMethod: "exportValue", GoMethod: "ExportValue"},
			_jsii_.MemberMethod{JsiiMethod: "formatArn", GoMethod: "FormatArn"},
			_jsii_.MemberMethod{JsiiMethod: "getLogicalId", GoMethod: "GetLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "nested", GoGetter: "Nested"},
			_jsii_.MemberProperty{JsiiProperty: "nestedStackParent", GoGetter: "NestedStackParent"},
			_jsii_.MemberProperty{JsiiProperty: "nestedStackResource", GoGetter: "NestedStackResource"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "notificationArns", GoGetter: "NotificationArns"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "parentStack", GoGetter: "ParentStack"},
			_jsii_.MemberMethod{JsiiMethod: "parseArn", GoMethod: "ParseArn"},
			_jsii_.MemberProperty{JsiiProperty: "partition", GoGetter: "Partition"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "prepareCrossReference", GoMethod: "PrepareCrossReference"},
			_jsii_.MemberProperty{JsiiProperty: "region", GoGetter: "Region"},
			_jsii_.MemberMethod{JsiiMethod: "renameLogicalId", GoMethod: "RenameLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "reportMissingContext", GoMethod: "ReportMissingContext"},
			_jsii_.MemberMethod{JsiiMethod: "reportMissingContextKey", GoMethod: "ReportMissingContextKey"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "setParameter", GoMethod: "SetParameter"},
			_jsii_.MemberMethod{JsiiMethod: "splitArn", GoMethod: "SplitArn"},
			_jsii_.MemberProperty{JsiiProperty: "stackId", GoGetter: "StackId"},
			_jsii_.MemberProperty{JsiiProperty: "stackName", GoGetter: "StackName"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "synthesizer", GoGetter: "Synthesizer"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "templateFile", GoGetter: "TemplateFile"},
			_jsii_.MemberProperty{JsiiProperty: "templateOptions", GoGetter: "TemplateOptions"},
			_jsii_.MemberProperty{JsiiProperty: "terminationProtection", GoGetter: "TerminationProtection"},
			_jsii_.MemberMethod{JsiiMethod: "toJsonString", GoMethod: "ToJsonString"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "urlSuffix", GoGetter: "UrlSuffix"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_NestedStack{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Stack)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.NestedStackProps",
		reflect.TypeOf((*NestedStackProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.NestedStackSynthesizer",
		reflect.TypeOf((*NestedStackSynthesizer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberMethod{JsiiMethod: "emitStackArtifact", GoMethod: "EmitStackArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeStackTemplate", GoMethod: "SynthesizeStackTemplate"},
		},
		func() interface{} {
			j := jsiiProxy_NestedStackSynthesizer{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_StackSynthesizer)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.PhysicalName",
		reflect.TypeOf((*PhysicalName)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_PhysicalName{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Reference",
		reflect.TypeOf((*Reference)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "displayName", GoGetter: "DisplayName"},
			_jsii_.MemberMethod{JsiiMethod: "newError", GoMethod: "NewError"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberProperty{JsiiProperty: "target", GoGetter: "Target"},
			_jsii_.MemberMethod{JsiiMethod: "toJSON", GoMethod: "ToJSON"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_Reference{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Intrinsic)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.RemovalPolicy",
		reflect.TypeOf((*RemovalPolicy)(nil)).Elem(),
		map[string]interface{}{
			"DESTROY": RemovalPolicy_DESTROY,
			"RETAIN": RemovalPolicy_RETAIN,
			"SNAPSHOT": RemovalPolicy_SNAPSHOT,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.RemovalPolicyOptions",
		reflect.TypeOf((*RemovalPolicyOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.RemoveTag",
		reflect.TypeOf((*RemoveTag)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyTag", GoMethod: "ApplyTag"},
			_jsii_.MemberProperty{JsiiProperty: "key", GoGetter: "Key"},
			_jsii_.MemberProperty{JsiiProperty: "props", GoGetter: "Props"},
			_jsii_.MemberMethod{JsiiMethod: "visit", GoMethod: "Visit"},
		},
		func() interface{} {
			j := jsiiProxy_RemoveTag{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IAspect)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.ResolveChangeContextOptions",
		reflect.TypeOf((*ResolveChangeContextOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.ResolveOptions",
		reflect.TypeOf((*ResolveOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Resource",
		reflect.TypeOf((*Resource)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Resource{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Construct)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.ResourceEnvironment",
		reflect.TypeOf((*ResourceEnvironment)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.ResourceProps",
		reflect.TypeOf((*ResourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.ReverseOptions",
		reflect.TypeOf((*ReverseOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.ScopedAws",
		reflect.TypeOf((*ScopedAws)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "accountId", GoGetter: "AccountId"},
			_jsii_.MemberProperty{JsiiProperty: "notificationArns", GoGetter: "NotificationArns"},
			_jsii_.MemberProperty{JsiiProperty: "partition", GoGetter: "Partition"},
			_jsii_.MemberProperty{JsiiProperty: "region", GoGetter: "Region"},
			_jsii_.MemberProperty{JsiiProperty: "stackId", GoGetter: "StackId"},
			_jsii_.MemberProperty{JsiiProperty: "stackName", GoGetter: "StackName"},
			_jsii_.MemberProperty{JsiiProperty: "urlSuffix", GoGetter: "UrlSuffix"},
		},
		func() interface{} {
			return &jsiiProxy_ScopedAws{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.SecretValue",
		reflect.TypeOf((*SecretValue)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "newError", GoMethod: "NewError"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "toJSON", GoMethod: "ToJSON"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_SecretValue{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Intrinsic)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.SecretsManagerSecretOptions",
		reflect.TypeOf((*SecretsManagerSecretOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Size",
		reflect.TypeOf((*Size)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "toGibibytes", GoMethod: "ToGibibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toKibibytes", GoMethod: "ToKibibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toMebibytes", GoMethod: "ToMebibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toPebibytes", GoMethod: "ToPebibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toTebibytes", GoMethod: "ToTebibytes"},
		},
		func() interface{} {
			return &jsiiProxy_Size{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.SizeConversionOptions",
		reflect.TypeOf((*SizeConversionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.SizeRoundingBehavior",
		reflect.TypeOf((*SizeRoundingBehavior)(nil)).Elem(),
		map[string]interface{}{
			"FAIL": SizeRoundingBehavior_FAIL,
			"FLOOR": SizeRoundingBehavior_FLOOR,
			"NONE": SizeRoundingBehavior_NONE,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Stack",
		reflect.TypeOf((*Stack)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "account", GoGetter: "Account"},
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addTransform", GoMethod: "AddTransform"},
			_jsii_.MemberMethod{JsiiMethod: "allocateLogicalId", GoMethod: "AllocateLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "artifactId", GoGetter: "ArtifactId"},
			_jsii_.MemberProperty{JsiiProperty: "availabilityZones", GoGetter: "AvailabilityZones"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberProperty{JsiiProperty: "environment", GoGetter: "Environment"},
			_jsii_.MemberMethod{JsiiMethod: "exportValue", GoMethod: "ExportValue"},
			_jsii_.MemberMethod{JsiiMethod: "formatArn", GoMethod: "FormatArn"},
			_jsii_.MemberMethod{JsiiMethod: "getLogicalId", GoMethod: "GetLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "nested", GoGetter: "Nested"},
			_jsii_.MemberProperty{JsiiProperty: "nestedStackParent", GoGetter: "NestedStackParent"},
			_jsii_.MemberProperty{JsiiProperty: "nestedStackResource", GoGetter: "NestedStackResource"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "notificationArns", GoGetter: "NotificationArns"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "parentStack", GoGetter: "ParentStack"},
			_jsii_.MemberMethod{JsiiMethod: "parseArn", GoMethod: "ParseArn"},
			_jsii_.MemberProperty{JsiiProperty: "partition", GoGetter: "Partition"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "prepareCrossReference", GoMethod: "PrepareCrossReference"},
			_jsii_.MemberProperty{JsiiProperty: "region", GoGetter: "Region"},
			_jsii_.MemberMethod{JsiiMethod: "renameLogicalId", GoMethod: "RenameLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "reportMissingContext", GoMethod: "ReportMissingContext"},
			_jsii_.MemberMethod{JsiiMethod: "reportMissingContextKey", GoMethod: "ReportMissingContextKey"},
			_jsii_.MemberMethod{JsiiMethod: "resolve", GoMethod: "Resolve"},
			_jsii_.MemberMethod{JsiiMethod: "splitArn", GoMethod: "SplitArn"},
			_jsii_.MemberProperty{JsiiProperty: "stackId", GoGetter: "StackId"},
			_jsii_.MemberProperty{JsiiProperty: "stackName", GoGetter: "StackName"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "synthesizer", GoGetter: "Synthesizer"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "templateFile", GoGetter: "TemplateFile"},
			_jsii_.MemberProperty{JsiiProperty: "templateOptions", GoGetter: "TemplateOptions"},
			_jsii_.MemberProperty{JsiiProperty: "terminationProtection", GoGetter: "TerminationProtection"},
			_jsii_.MemberMethod{JsiiMethod: "toJsonString", GoMethod: "ToJsonString"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "urlSuffix", GoGetter: "UrlSuffix"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Stack{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Construct)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ITaggable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.StackProps",
		reflect.TypeOf((*StackProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.StackSynthesizer",
		reflect.TypeOf((*StackSynthesizer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDockerImageAsset", GoMethod: "AddDockerImageAsset"},
			_jsii_.MemberMethod{JsiiMethod: "addFileAsset", GoMethod: "AddFileAsset"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberMethod{JsiiMethod: "emitStackArtifact", GoMethod: "EmitStackArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeStackTemplate", GoMethod: "SynthesizeStackTemplate"},
		},
		func() interface{} {
			j := jsiiProxy_StackSynthesizer{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IStackSynthesizer)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Stage",
		reflect.TypeOf((*Stage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "account", GoGetter: "Account"},
			_jsii_.MemberProperty{JsiiProperty: "artifactId", GoGetter: "ArtifactId"},
			_jsii_.MemberProperty{JsiiProperty: "assetOutdir", GoGetter: "AssetOutdir"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outdir", GoGetter: "Outdir"},
			_jsii_.MemberProperty{JsiiProperty: "parentStage", GoGetter: "ParentStage"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "region", GoGetter: "Region"},
			_jsii_.MemberProperty{JsiiProperty: "stageName", GoGetter: "StageName"},
			_jsii_.MemberMethod{JsiiMethod: "synth", GoMethod: "Synth"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Stage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Construct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.StageProps",
		reflect.TypeOf((*StageProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.StageSynthesisOptions",
		reflect.TypeOf((*StageSynthesisOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.StringConcat",
		reflect.TypeOf((*StringConcat)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "join", GoMethod: "Join"},
		},
		func() interface{} {
			j := jsiiProxy_StringConcat{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IFragmentConcatenator)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.SymlinkFollowMode",
		reflect.TypeOf((*SymlinkFollowMode)(nil)).Elem(),
		map[string]interface{}{
			"NEVER": SymlinkFollowMode_NEVER,
			"ALWAYS": SymlinkFollowMode_ALWAYS,
			"EXTERNAL": SymlinkFollowMode_EXTERNAL,
			"BLOCK_EXTERNAL": SymlinkFollowMode_BLOCK_EXTERNAL,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.SynthesisOptions",
		reflect.TypeOf((*SynthesisOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.SynthesizeStackArtifactOptions",
		reflect.TypeOf((*SynthesizeStackArtifactOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Tag",
		reflect.TypeOf((*Tag)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyTag", GoMethod: "ApplyTag"},
			_jsii_.MemberProperty{JsiiProperty: "key", GoGetter: "Key"},
			_jsii_.MemberProperty{JsiiProperty: "props", GoGetter: "Props"},
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
			_jsii_.MemberMethod{JsiiMethod: "visit", GoMethod: "Visit"},
		},
		func() interface{} {
			j := jsiiProxy_Tag{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IAspect)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.TagManager",
		reflect.TypeOf((*TagManager)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyTagAspectHere", GoMethod: "ApplyTagAspectHere"},
			_jsii_.MemberMethod{JsiiMethod: "hasTags", GoMethod: "HasTags"},
			_jsii_.MemberMethod{JsiiMethod: "removeTag", GoMethod: "RemoveTag"},
			_jsii_.MemberMethod{JsiiMethod: "renderTags", GoMethod: "RenderTags"},
			_jsii_.MemberMethod{JsiiMethod: "setTag", GoMethod: "SetTag"},
			_jsii_.MemberProperty{JsiiProperty: "tagPropertyName", GoGetter: "TagPropertyName"},
			_jsii_.MemberMethod{JsiiMethod: "tagValues", GoMethod: "TagValues"},
		},
		func() interface{} {
			return &jsiiProxy_TagManager{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.TagManagerOptions",
		reflect.TypeOf((*TagManagerOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.TagProps",
		reflect.TypeOf((*TagProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.TagType",
		reflect.TypeOf((*TagType)(nil)).Elem(),
		map[string]interface{}{
			"STANDARD": TagType_STANDARD,
			"AUTOSCALING_GROUP": TagType_AUTOSCALING_GROUP,
			"MAP": TagType_MAP,
			"KEY_VALUE": TagType_KEY_VALUE,
			"NOT_TAGGABLE": TagType_NOT_TAGGABLE,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Tags",
		reflect.TypeOf((*Tags)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
			_jsii_.MemberMethod{JsiiMethod: "remove", GoMethod: "Remove"},
		},
		func() interface{} {
			return &jsiiProxy_Tags{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.TimeConversionOptions",
		reflect.TypeOf((*TimeConversionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.Token",
		reflect.TypeOf((*Token)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Token{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.TokenComparison",
		reflect.TypeOf((*TokenComparison)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_TokenComparison{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.Tokenization",
		reflect.TypeOf((*Tokenization)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Tokenization{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.TokenizedStringFragments",
		reflect.TypeOf((*TokenizedStringFragments)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addIntrinsic", GoMethod: "AddIntrinsic"},
			_jsii_.MemberMethod{JsiiMethod: "addLiteral", GoMethod: "AddLiteral"},
			_jsii_.MemberMethod{JsiiMethod: "addToken", GoMethod: "AddToken"},
			_jsii_.MemberProperty{JsiiProperty: "firstToken", GoGetter: "FirstToken"},
			_jsii_.MemberProperty{JsiiProperty: "firstValue", GoGetter: "FirstValue"},
			_jsii_.MemberMethod{JsiiMethod: "join", GoMethod: "Join"},
			_jsii_.MemberProperty{JsiiProperty: "length", GoGetter: "Length"},
			_jsii_.MemberMethod{JsiiMethod: "mapTokens", GoMethod: "MapTokens"},
			_jsii_.MemberProperty{JsiiProperty: "tokens", GoGetter: "Tokens"},
		},
		func() interface{} {
			return &jsiiProxy_TokenizedStringFragments{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.TreeInspector",
		reflect.TypeOf((*TreeInspector)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addAttribute", GoMethod: "AddAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "attributes", GoGetter: "Attributes"},
		},
		func() interface{} {
			return &jsiiProxy_TreeInspector{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.ValidationError",
		reflect.TypeOf((*ValidationError)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.ValidationResult",
		reflect.TypeOf((*ValidationResult)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "assertSuccess", GoMethod: "AssertSuccess"},
			_jsii_.MemberProperty{JsiiProperty: "errorMessage", GoGetter: "ErrorMessage"},
			_jsii_.MemberMethod{JsiiMethod: "errorTree", GoMethod: "ErrorTree"},
			_jsii_.MemberProperty{JsiiProperty: "isSuccess", GoGetter: "IsSuccess"},
			_jsii_.MemberMethod{JsiiMethod: "prefix", GoMethod: "Prefix"},
			_jsii_.MemberProperty{JsiiProperty: "results", GoGetter: "Results"},
		},
		func() interface{} {
			return &jsiiProxy_ValidationResult{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.ValidationResults",
		reflect.TypeOf((*ValidationResults)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "collect", GoMethod: "Collect"},
			_jsii_.MemberMethod{JsiiMethod: "errorTreeList", GoMethod: "ErrorTreeList"},
			_jsii_.MemberProperty{JsiiProperty: "isSuccess", GoGetter: "IsSuccess"},
			_jsii_.MemberProperty{JsiiProperty: "results", GoGetter: "Results"},
			_jsii_.MemberMethod{JsiiMethod: "wrap", GoMethod: "Wrap"},
		},
		func() interface{} {
			return &jsiiProxy_ValidationResults{}
		},
	)
}

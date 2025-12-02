package awscodebuild

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.Artifacts",
		reflect.TypeOf((*Artifacts)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "identifier", GoGetter: "Identifier"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
		},
		func() interface{} {
			j := jsiiProxy_Artifacts{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IArtifacts)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.ArtifactsConfig",
		reflect.TypeOf((*ArtifactsConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.ArtifactsProps",
		reflect.TypeOf((*ArtifactsProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BatchBuildConfig",
		reflect.TypeOf((*BatchBuildConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BindToCodePipelineOptions",
		reflect.TypeOf((*BindToCodePipelineOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.BitBucketSourceCredentials",
		reflect.TypeOf((*BitBucketSourceCredentials)(nil)).Elem(),
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
			j := jsiiProxy_BitBucketSourceCredentials{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BitBucketSourceCredentialsProps",
		reflect.TypeOf((*BitBucketSourceCredentialsProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BitBucketSourceProps",
		reflect.TypeOf((*BitBucketSourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BucketCacheOptions",
		reflect.TypeOf((*BucketCacheOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BuildEnvironment",
		reflect.TypeOf((*BuildEnvironment)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BuildEnvironmentVariable",
		reflect.TypeOf((*BuildEnvironmentVariable)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codebuild.BuildEnvironmentVariableType",
		reflect.TypeOf((*BuildEnvironmentVariableType)(nil)).Elem(),
		map[string]interface{}{
			"PLAINTEXT": BuildEnvironmentVariableType_PLAINTEXT,
			"PARAMETER_STORE": BuildEnvironmentVariableType_PARAMETER_STORE,
			"SECRETS_MANAGER": BuildEnvironmentVariableType_SECRETS_MANAGER,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BuildImageBindOptions",
		reflect.TypeOf((*BuildImageBindOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.BuildImageConfig",
		reflect.TypeOf((*BuildImageConfig)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.BuildSpec",
		reflect.TypeOf((*BuildSpec)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "isImmediate", GoGetter: "IsImmediate"},
			_jsii_.MemberMethod{JsiiMethod: "toBuildSpec", GoMethod: "ToBuildSpec"},
		},
		func() interface{} {
			return &jsiiProxy_BuildSpec{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.Cache",
		reflect.TypeOf((*Cache)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Cache{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.CfnProject",
		reflect.TypeOf((*CfnProject)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "artifacts", GoGetter: "Artifacts"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "badgeEnabled", GoGetter: "BadgeEnabled"},
			_jsii_.MemberProperty{JsiiProperty: "buildBatchConfig", GoGetter: "BuildBatchConfig"},
			_jsii_.MemberProperty{JsiiProperty: "cache", GoGetter: "Cache"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "concurrentBuildLimit", GoGetter: "ConcurrentBuildLimit"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "environment", GoGetter: "Environment"},
			_jsii_.MemberProperty{JsiiProperty: "fileSystemLocations", GoGetter: "FileSystemLocations"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "logsConfig", GoGetter: "LogsConfig"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "queuedTimeoutInMinutes", GoGetter: "QueuedTimeoutInMinutes"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "secondaryArtifacts", GoGetter: "SecondaryArtifacts"},
			_jsii_.MemberProperty{JsiiProperty: "secondarySources", GoGetter: "SecondarySources"},
			_jsii_.MemberProperty{JsiiProperty: "secondarySourceVersions", GoGetter: "SecondarySourceVersions"},
			_jsii_.MemberProperty{JsiiProperty: "serviceRole", GoGetter: "ServiceRole"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "source", GoGetter: "Source"},
			_jsii_.MemberProperty{JsiiProperty: "sourceVersion", GoGetter: "SourceVersion"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "timeoutInMinutes", GoGetter: "TimeoutInMinutes"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "triggers", GoGetter: "Triggers"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "vpcConfig", GoGetter: "VpcConfig"},
		},
		func() interface{} {
			j := jsiiProxy_CfnProject{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.ArtifactsProperty",
		reflect.TypeOf((*CfnProject_ArtifactsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.BatchRestrictionsProperty",
		reflect.TypeOf((*CfnProject_BatchRestrictionsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.BuildStatusConfigProperty",
		reflect.TypeOf((*CfnProject_BuildStatusConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.CloudWatchLogsConfigProperty",
		reflect.TypeOf((*CfnProject_CloudWatchLogsConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.EnvironmentProperty",
		reflect.TypeOf((*CfnProject_EnvironmentProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.EnvironmentVariableProperty",
		reflect.TypeOf((*CfnProject_EnvironmentVariableProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.GitSubmodulesConfigProperty",
		reflect.TypeOf((*CfnProject_GitSubmodulesConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.LogsConfigProperty",
		reflect.TypeOf((*CfnProject_LogsConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.ProjectBuildBatchConfigProperty",
		reflect.TypeOf((*CfnProject_ProjectBuildBatchConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.ProjectCacheProperty",
		reflect.TypeOf((*CfnProject_ProjectCacheProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.ProjectFileSystemLocationProperty",
		reflect.TypeOf((*CfnProject_ProjectFileSystemLocationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.ProjectSourceVersionProperty",
		reflect.TypeOf((*CfnProject_ProjectSourceVersionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.ProjectTriggersProperty",
		reflect.TypeOf((*CfnProject_ProjectTriggersProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.RegistryCredentialProperty",
		reflect.TypeOf((*CfnProject_RegistryCredentialProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.S3LogsConfigProperty",
		reflect.TypeOf((*CfnProject_S3LogsConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.SourceAuthProperty",
		reflect.TypeOf((*CfnProject_SourceAuthProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.SourceProperty",
		reflect.TypeOf((*CfnProject_SourceProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.VpcConfigProperty",
		reflect.TypeOf((*CfnProject_VpcConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProject.WebhookFilterProperty",
		reflect.TypeOf((*CfnProject_WebhookFilterProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnProjectProps",
		reflect.TypeOf((*CfnProjectProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.CfnReportGroup",
		reflect.TypeOf((*CfnReportGroup)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "deleteReports", GoGetter: "DeleteReports"},
			_jsii_.MemberProperty{JsiiProperty: "exportConfig", GoGetter: "ExportConfig"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
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
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnReportGroup{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnReportGroup.ReportExportConfigProperty",
		reflect.TypeOf((*CfnReportGroup_ReportExportConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnReportGroup.S3ReportExportConfigProperty",
		reflect.TypeOf((*CfnReportGroup_S3ReportExportConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnReportGroupProps",
		reflect.TypeOf((*CfnReportGroupProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.CfnSourceCredential",
		reflect.TypeOf((*CfnSourceCredential)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "authType", GoGetter: "AuthType"},
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
			_jsii_.MemberProperty{JsiiProperty: "serverType", GoGetter: "ServerType"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "token", GoGetter: "Token"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberProperty{JsiiProperty: "username", GoGetter: "Username"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnSourceCredential{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CfnSourceCredentialProps",
		reflect.TypeOf((*CfnSourceCredentialProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CloudWatchLoggingOptions",
		reflect.TypeOf((*CloudWatchLoggingOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CodeCommitSourceProps",
		reflect.TypeOf((*CodeCommitSourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.CommonProjectProps",
		reflect.TypeOf((*CommonProjectProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codebuild.ComputeType",
		reflect.TypeOf((*ComputeType)(nil)).Elem(),
		map[string]interface{}{
			"SMALL": ComputeType_SMALL,
			"MEDIUM": ComputeType_MEDIUM,
			"LARGE": ComputeType_LARGE,
			"X2_LARGE": ComputeType_X2_LARGE,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.DockerImageOptions",
		reflect.TypeOf((*DockerImageOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.EfsFileSystemLocationProps",
		reflect.TypeOf((*EfsFileSystemLocationProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codebuild.EventAction",
		reflect.TypeOf((*EventAction)(nil)).Elem(),
		map[string]interface{}{
			"PUSH": EventAction_PUSH,
			"PULL_REQUEST_CREATED": EventAction_PULL_REQUEST_CREATED,
			"PULL_REQUEST_UPDATED": EventAction_PULL_REQUEST_UPDATED,
			"PULL_REQUEST_MERGED": EventAction_PULL_REQUEST_MERGED,
			"PULL_REQUEST_REOPENED": EventAction_PULL_REQUEST_REOPENED,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.FileSystemConfig",
		reflect.TypeOf((*FileSystemConfig)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.FileSystemLocation",
		reflect.TypeOf((*FileSystemLocation)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_FileSystemLocation{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.FilterGroup",
		reflect.TypeOf((*FilterGroup)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "andActorAccountIs", GoMethod: "AndActorAccountIs"},
			_jsii_.MemberMethod{JsiiMethod: "andActorAccountIsNot", GoMethod: "AndActorAccountIsNot"},
			_jsii_.MemberMethod{JsiiMethod: "andBaseBranchIs", GoMethod: "AndBaseBranchIs"},
			_jsii_.MemberMethod{JsiiMethod: "andBaseBranchIsNot", GoMethod: "AndBaseBranchIsNot"},
			_jsii_.MemberMethod{JsiiMethod: "andBaseRefIs", GoMethod: "AndBaseRefIs"},
			_jsii_.MemberMethod{JsiiMethod: "andBaseRefIsNot", GoMethod: "AndBaseRefIsNot"},
			_jsii_.MemberMethod{JsiiMethod: "andBranchIs", GoMethod: "AndBranchIs"},
			_jsii_.MemberMethod{JsiiMethod: "andBranchIsNot", GoMethod: "AndBranchIsNot"},
			_jsii_.MemberMethod{JsiiMethod: "andCommitMessageIs", GoMethod: "AndCommitMessageIs"},
			_jsii_.MemberMethod{JsiiMethod: "andCommitMessageIsNot", GoMethod: "AndCommitMessageIsNot"},
			_jsii_.MemberMethod{JsiiMethod: "andFilePathIs", GoMethod: "AndFilePathIs"},
			_jsii_.MemberMethod{JsiiMethod: "andFilePathIsNot", GoMethod: "AndFilePathIsNot"},
			_jsii_.MemberMethod{JsiiMethod: "andHeadRefIs", GoMethod: "AndHeadRefIs"},
			_jsii_.MemberMethod{JsiiMethod: "andHeadRefIsNot", GoMethod: "AndHeadRefIsNot"},
			_jsii_.MemberMethod{JsiiMethod: "andTagIs", GoMethod: "AndTagIs"},
			_jsii_.MemberMethod{JsiiMethod: "andTagIsNot", GoMethod: "AndTagIsNot"},
		},
		func() interface{} {
			return &jsiiProxy_FilterGroup{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.GitHubEnterpriseSourceCredentials",
		reflect.TypeOf((*GitHubEnterpriseSourceCredentials)(nil)).Elem(),
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
			j := jsiiProxy_GitHubEnterpriseSourceCredentials{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.GitHubEnterpriseSourceCredentialsProps",
		reflect.TypeOf((*GitHubEnterpriseSourceCredentialsProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.GitHubEnterpriseSourceProps",
		reflect.TypeOf((*GitHubEnterpriseSourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.GitHubSourceCredentials",
		reflect.TypeOf((*GitHubSourceCredentials)(nil)).Elem(),
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
			j := jsiiProxy_GitHubSourceCredentials{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.GitHubSourceCredentialsProps",
		reflect.TypeOf((*GitHubSourceCredentialsProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.GitHubSourceProps",
		reflect.TypeOf((*GitHubSourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codebuild.IArtifacts",
		reflect.TypeOf((*IArtifacts)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "identifier", GoGetter: "Identifier"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
		},
		func() interface{} {
			return &jsiiProxy_IArtifacts{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codebuild.IBindableBuildImage",
		reflect.TypeOf((*IBindableBuildImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "defaultComputeType", GoGetter: "DefaultComputeType"},
			_jsii_.MemberProperty{JsiiProperty: "imageId", GoGetter: "ImageId"},
			_jsii_.MemberProperty{JsiiProperty: "imagePullPrincipalType", GoGetter: "ImagePullPrincipalType"},
			_jsii_.MemberProperty{JsiiProperty: "repository", GoGetter: "Repository"},
			_jsii_.MemberMethod{JsiiMethod: "runScriptBuildspec", GoMethod: "RunScriptBuildspec"},
			_jsii_.MemberProperty{JsiiProperty: "secretsManagerCredentials", GoGetter: "SecretsManagerCredentials"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_IBindableBuildImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IBuildImage)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codebuild.IBuildImage",
		reflect.TypeOf((*IBuildImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "defaultComputeType", GoGetter: "DefaultComputeType"},
			_jsii_.MemberProperty{JsiiProperty: "imageId", GoGetter: "ImageId"},
			_jsii_.MemberProperty{JsiiProperty: "imagePullPrincipalType", GoGetter: "ImagePullPrincipalType"},
			_jsii_.MemberProperty{JsiiProperty: "repository", GoGetter: "Repository"},
			_jsii_.MemberMethod{JsiiMethod: "runScriptBuildspec", GoMethod: "RunScriptBuildspec"},
			_jsii_.MemberProperty{JsiiProperty: "secretsManagerCredentials", GoGetter: "SecretsManagerCredentials"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			return &jsiiProxy_IBuildImage{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codebuild.IFileSystemLocation",
		reflect.TypeOf((*IFileSystemLocation)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_IFileSystemLocation{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codebuild.IProject",
		reflect.TypeOf((*IProject)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addToRolePolicy", GoMethod: "AddToRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "bindAsNotificationRuleSource", GoMethod: "BindAsNotificationRuleSource"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberMethod{JsiiMethod: "enableBatchBuilds", GoMethod: "EnableBatchBuilds"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "grantPrincipal", GoGetter: "GrantPrincipal"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricBuilds", GoMethod: "MetricBuilds"},
			_jsii_.MemberMethod{JsiiMethod: "metricDuration", GoMethod: "MetricDuration"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailedBuilds", GoMethod: "MetricFailedBuilds"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceededBuilds", GoMethod: "MetricSucceededBuilds"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOn", GoMethod: "NotifyOn"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOnBuildFailed", GoMethod: "NotifyOnBuildFailed"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOnBuildSucceeded", GoMethod: "NotifyOnBuildSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildFailed", GoMethod: "OnBuildFailed"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildStarted", GoMethod: "OnBuildStarted"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildSucceeded", GoMethod: "OnBuildSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onPhaseChange", GoMethod: "OnPhaseChange"},
			_jsii_.MemberMethod{JsiiMethod: "onStateChange", GoMethod: "OnStateChange"},
			_jsii_.MemberProperty{JsiiProperty: "projectArn", GoGetter: "ProjectArn"},
			_jsii_.MemberProperty{JsiiProperty: "projectName", GoGetter: "ProjectName"},
			_jsii_.MemberProperty{JsiiProperty: "role", GoGetter: "Role"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IProject{}
			_jsii_.InitJsiiProxy(&j.Type__awsec2IConnectable)
			_jsii_.InitJsiiProxy(&j.Type__awsiamIGrantable)
			_jsii_.InitJsiiProxy(&j.Type__awscodestarnotificationsINotificationRuleSource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codebuild.IReportGroup",
		reflect.TypeOf((*IReportGroup)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "reportGroupArn", GoGetter: "ReportGroupArn"},
			_jsii_.MemberProperty{JsiiProperty: "reportGroupName", GoGetter: "ReportGroupName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IReportGroup{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codebuild.ISource",
		reflect.TypeOf((*ISource)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "badgeSupported", GoGetter: "BadgeSupported"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "identifier", GoGetter: "Identifier"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
		},
		func() interface{} {
			return &jsiiProxy_ISource{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codebuild.ImagePullPrincipalType",
		reflect.TypeOf((*ImagePullPrincipalType)(nil)).Elem(),
		map[string]interface{}{
			"CODEBUILD": ImagePullPrincipalType_CODEBUILD,
			"SERVICE_ROLE": ImagePullPrincipalType_SERVICE_ROLE,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.LinuxBuildImage",
		reflect.TypeOf((*LinuxBuildImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "defaultComputeType", GoGetter: "DefaultComputeType"},
			_jsii_.MemberProperty{JsiiProperty: "imageId", GoGetter: "ImageId"},
			_jsii_.MemberProperty{JsiiProperty: "imagePullPrincipalType", GoGetter: "ImagePullPrincipalType"},
			_jsii_.MemberProperty{JsiiProperty: "repository", GoGetter: "Repository"},
			_jsii_.MemberMethod{JsiiMethod: "runScriptBuildspec", GoMethod: "RunScriptBuildspec"},
			_jsii_.MemberProperty{JsiiProperty: "secretsManagerCredentials", GoGetter: "SecretsManagerCredentials"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_LinuxBuildImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IBuildImage)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		reflect.TypeOf((*LinuxGpuBuildImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "defaultComputeType", GoGetter: "DefaultComputeType"},
			_jsii_.MemberProperty{JsiiProperty: "imageId", GoGetter: "ImageId"},
			_jsii_.MemberProperty{JsiiProperty: "imagePullPrincipalType", GoGetter: "ImagePullPrincipalType"},
			_jsii_.MemberMethod{JsiiMethod: "runScriptBuildspec", GoMethod: "RunScriptBuildspec"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_LinuxGpuBuildImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IBindableBuildImage)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codebuild.LocalCacheMode",
		reflect.TypeOf((*LocalCacheMode)(nil)).Elem(),
		map[string]interface{}{
			"SOURCE": LocalCacheMode_SOURCE,
			"DOCKER_LAYER": LocalCacheMode_DOCKER_LAYER,
			"CUSTOM": LocalCacheMode_CUSTOM,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.LoggingOptions",
		reflect.TypeOf((*LoggingOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.PhaseChangeEvent",
		reflect.TypeOf((*PhaseChangeEvent)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_PhaseChangeEvent{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.PipelineProject",
		reflect.TypeOf((*PipelineProject)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addFileSystemLocation", GoMethod: "AddFileSystemLocation"},
			_jsii_.MemberMethod{JsiiMethod: "addSecondaryArtifact", GoMethod: "AddSecondaryArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "addSecondarySource", GoMethod: "AddSecondarySource"},
			_jsii_.MemberMethod{JsiiMethod: "addToRolePolicy", GoMethod: "AddToRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "bindAsNotificationRuleSource", GoMethod: "BindAsNotificationRuleSource"},
			_jsii_.MemberMethod{JsiiMethod: "bindToCodePipeline", GoMethod: "BindToCodePipeline"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberMethod{JsiiMethod: "enableBatchBuilds", GoMethod: "EnableBatchBuilds"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "grantPrincipal", GoGetter: "GrantPrincipal"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricBuilds", GoMethod: "MetricBuilds"},
			_jsii_.MemberMethod{JsiiMethod: "metricDuration", GoMethod: "MetricDuration"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailedBuilds", GoMethod: "MetricFailedBuilds"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceededBuilds", GoMethod: "MetricSucceededBuilds"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOn", GoMethod: "NotifyOn"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOnBuildFailed", GoMethod: "NotifyOnBuildFailed"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOnBuildSucceeded", GoMethod: "NotifyOnBuildSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildFailed", GoMethod: "OnBuildFailed"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildStarted", GoMethod: "OnBuildStarted"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildSucceeded", GoMethod: "OnBuildSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onPhaseChange", GoMethod: "OnPhaseChange"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onStateChange", GoMethod: "OnStateChange"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "projectArn", GoGetter: "ProjectArn"},
			_jsii_.MemberProperty{JsiiProperty: "projectName", GoGetter: "ProjectName"},
			_jsii_.MemberProperty{JsiiProperty: "role", GoGetter: "Role"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_PipelineProject{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Project)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.PipelineProjectProps",
		reflect.TypeOf((*PipelineProjectProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.Project",
		reflect.TypeOf((*Project)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addFileSystemLocation", GoMethod: "AddFileSystemLocation"},
			_jsii_.MemberMethod{JsiiMethod: "addSecondaryArtifact", GoMethod: "AddSecondaryArtifact"},
			_jsii_.MemberMethod{JsiiMethod: "addSecondarySource", GoMethod: "AddSecondarySource"},
			_jsii_.MemberMethod{JsiiMethod: "addToRolePolicy", GoMethod: "AddToRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "bindAsNotificationRuleSource", GoMethod: "BindAsNotificationRuleSource"},
			_jsii_.MemberMethod{JsiiMethod: "bindToCodePipeline", GoMethod: "BindToCodePipeline"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberMethod{JsiiMethod: "enableBatchBuilds", GoMethod: "EnableBatchBuilds"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "grantPrincipal", GoGetter: "GrantPrincipal"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricBuilds", GoMethod: "MetricBuilds"},
			_jsii_.MemberMethod{JsiiMethod: "metricDuration", GoMethod: "MetricDuration"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailedBuilds", GoMethod: "MetricFailedBuilds"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceededBuilds", GoMethod: "MetricSucceededBuilds"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOn", GoMethod: "NotifyOn"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOnBuildFailed", GoMethod: "NotifyOnBuildFailed"},
			_jsii_.MemberMethod{JsiiMethod: "notifyOnBuildSucceeded", GoMethod: "NotifyOnBuildSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildFailed", GoMethod: "OnBuildFailed"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildStarted", GoMethod: "OnBuildStarted"},
			_jsii_.MemberMethod{JsiiMethod: "onBuildSucceeded", GoMethod: "OnBuildSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onPhaseChange", GoMethod: "OnPhaseChange"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onStateChange", GoMethod: "OnStateChange"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "projectArn", GoGetter: "ProjectArn"},
			_jsii_.MemberProperty{JsiiProperty: "projectName", GoGetter: "ProjectName"},
			_jsii_.MemberProperty{JsiiProperty: "role", GoGetter: "Role"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Project{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IProject)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codebuild.ProjectNotificationEvents",
		reflect.TypeOf((*ProjectNotificationEvents)(nil)).Elem(),
		map[string]interface{}{
			"BUILD_FAILED": ProjectNotificationEvents_BUILD_FAILED,
			"BUILD_SUCCEEDED": ProjectNotificationEvents_BUILD_SUCCEEDED,
			"BUILD_IN_PROGRESS": ProjectNotificationEvents_BUILD_IN_PROGRESS,
			"BUILD_STOPPED": ProjectNotificationEvents_BUILD_STOPPED,
			"BUILD_PHASE_FAILED": ProjectNotificationEvents_BUILD_PHASE_FAILED,
			"BUILD_PHASE_SUCCEEDED": ProjectNotificationEvents_BUILD_PHASE_SUCCEEDED,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.ProjectNotifyOnOptions",
		reflect.TypeOf((*ProjectNotifyOnOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.ProjectProps",
		reflect.TypeOf((*ProjectProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.ReportGroup",
		reflect.TypeOf((*ReportGroup)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "exportBucket", GoGetter: "ExportBucket"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "reportGroupArn", GoGetter: "ReportGroupArn"},
			_jsii_.MemberProperty{JsiiProperty: "reportGroupName", GoGetter: "ReportGroupName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ReportGroup{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IReportGroup)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.ReportGroupProps",
		reflect.TypeOf((*ReportGroupProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.S3ArtifactsProps",
		reflect.TypeOf((*S3ArtifactsProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.S3LoggingOptions",
		reflect.TypeOf((*S3LoggingOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.S3SourceProps",
		reflect.TypeOf((*S3SourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.Source",
		reflect.TypeOf((*Source)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "badgeSupported", GoGetter: "BadgeSupported"},
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "identifier", GoGetter: "Identifier"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
		},
		func() interface{} {
			j := jsiiProxy_Source{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ISource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.SourceConfig",
		reflect.TypeOf((*SourceConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.SourceProps",
		reflect.TypeOf((*SourceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.StateChangeEvent",
		reflect.TypeOf((*StateChangeEvent)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_StateChangeEvent{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		reflect.TypeOf((*UntrustedCodeBoundaryPolicy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addStatements", GoMethod: "AddStatements"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "attachToGroup", GoMethod: "AttachToGroup"},
			_jsii_.MemberMethod{JsiiMethod: "attachToRole", GoMethod: "AttachToRole"},
			_jsii_.MemberMethod{JsiiMethod: "attachToUser", GoMethod: "AttachToUser"},
			_jsii_.MemberProperty{JsiiProperty: "description", GoGetter: "Description"},
			_jsii_.MemberProperty{JsiiProperty: "document", GoGetter: "Document"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "managedPolicyArn", GoGetter: "ManagedPolicyArn"},
			_jsii_.MemberProperty{JsiiProperty: "managedPolicyName", GoGetter: "ManagedPolicyName"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "path", GoGetter: "Path"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UntrustedCodeBoundaryPolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awsiamManagedPolicy)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicyProps",
		reflect.TypeOf((*UntrustedCodeBoundaryPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codebuild.WindowsBuildImage",
		reflect.TypeOf((*WindowsBuildImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "defaultComputeType", GoGetter: "DefaultComputeType"},
			_jsii_.MemberProperty{JsiiProperty: "imageId", GoGetter: "ImageId"},
			_jsii_.MemberProperty{JsiiProperty: "imagePullPrincipalType", GoGetter: "ImagePullPrincipalType"},
			_jsii_.MemberProperty{JsiiProperty: "repository", GoGetter: "Repository"},
			_jsii_.MemberMethod{JsiiMethod: "runScriptBuildspec", GoMethod: "RunScriptBuildspec"},
			_jsii_.MemberProperty{JsiiProperty: "secretsManagerCredentials", GoGetter: "SecretsManagerCredentials"},
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_WindowsBuildImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IBuildImage)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codebuild.WindowsImageType",
		reflect.TypeOf((*WindowsImageType)(nil)).Elem(),
		map[string]interface{}{
			"STANDARD": WindowsImageType_STANDARD,
			"SERVER_2019": WindowsImageType_SERVER_2019,
		},
	)
}

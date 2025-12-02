package awsecr

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.AuthorizationToken",
		reflect.TypeOf((*AuthorizationToken)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_AuthorizationToken{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.CfnPublicRepository",
		reflect.TypeOf((*CfnPublicRepository)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "repositoryCatalogData", GoGetter: "RepositoryCatalogData"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryPolicyText", GoGetter: "RepositoryPolicyText"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnPublicRepository{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnPublicRepositoryProps",
		reflect.TypeOf((*CfnPublicRepositoryProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.CfnRegistryPolicy",
		reflect.TypeOf((*CfnRegistryPolicy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrRegistryId", GoGetter: "AttrRegistryId"},
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
			_jsii_.MemberProperty{JsiiProperty: "policyText", GoGetter: "PolicyText"},
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
			j := jsiiProxy_CfnRegistryPolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnRegistryPolicyProps",
		reflect.TypeOf((*CfnRegistryPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.CfnReplicationConfiguration",
		reflect.TypeOf((*CfnReplicationConfiguration)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrRegistryId", GoGetter: "AttrRegistryId"},
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
			_jsii_.MemberProperty{JsiiProperty: "replicationConfiguration", GoGetter: "ReplicationConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnReplicationConfiguration{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnReplicationConfiguration.ReplicationConfigurationProperty",
		reflect.TypeOf((*CfnReplicationConfiguration_ReplicationConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnReplicationConfiguration.ReplicationDestinationProperty",
		reflect.TypeOf((*CfnReplicationConfiguration_ReplicationDestinationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnReplicationConfiguration.ReplicationRuleProperty",
		reflect.TypeOf((*CfnReplicationConfiguration_ReplicationRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnReplicationConfigurationProps",
		reflect.TypeOf((*CfnReplicationConfigurationProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.CfnRepository",
		reflect.TypeOf((*CfnRepository)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrRepositoryUri", GoGetter: "AttrRepositoryUri"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionConfiguration", GoGetter: "EncryptionConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "imageScanningConfiguration", GoGetter: "ImageScanningConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "imageTagMutability", GoGetter: "ImageTagMutability"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "lifecyclePolicy", GoGetter: "LifecyclePolicy"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryPolicyText", GoGetter: "RepositoryPolicyText"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnRepository{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnRepository.EncryptionConfigurationProperty",
		reflect.TypeOf((*CfnRepository_EncryptionConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnRepository.ImageScanningConfigurationProperty",
		reflect.TypeOf((*CfnRepository_ImageScanningConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnRepository.LifecyclePolicyProperty",
		reflect.TypeOf((*CfnRepository_LifecyclePolicyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.CfnRepositoryProps",
		reflect.TypeOf((*CfnRepositoryProps)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecr.IRepository",
		reflect.TypeOf((*IRepository)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantPull", GoMethod: "GrantPull"},
			_jsii_.MemberMethod{JsiiMethod: "grantPullPush", GoMethod: "GrantPullPush"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailEvent", GoMethod: "OnCloudTrailEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailImagePushed", GoMethod: "OnCloudTrailImagePushed"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onImageScanCompleted", GoMethod: "OnImageScanCompleted"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryArn", GoGetter: "RepositoryArn"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryUri", GoGetter: "RepositoryUri"},
			_jsii_.MemberMethod{JsiiMethod: "repositoryUriForDigest", GoMethod: "RepositoryUriForDigest"},
			_jsii_.MemberMethod{JsiiMethod: "repositoryUriForTag", GoMethod: "RepositoryUriForTag"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IRepository{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.LifecycleRule",
		reflect.TypeOf((*LifecycleRule)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.OnCloudTrailImagePushedOptions",
		reflect.TypeOf((*OnCloudTrailImagePushedOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.OnImageScanCompletedOptions",
		reflect.TypeOf((*OnImageScanCompletedOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.PublicGalleryAuthorizationToken",
		reflect.TypeOf((*PublicGalleryAuthorizationToken)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_PublicGalleryAuthorizationToken{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.Repository",
		reflect.TypeOf((*Repository)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addLifecycleRule", GoMethod: "AddLifecycleRule"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantPull", GoMethod: "GrantPull"},
			_jsii_.MemberMethod{JsiiMethod: "grantPullPush", GoMethod: "GrantPullPush"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailEvent", GoMethod: "OnCloudTrailEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailImagePushed", GoMethod: "OnCloudTrailImagePushed"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onImageScanCompleted", GoMethod: "OnImageScanCompleted"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryArn", GoGetter: "RepositoryArn"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryUri", GoGetter: "RepositoryUri"},
			_jsii_.MemberMethod{JsiiMethod: "repositoryUriForDigest", GoMethod: "RepositoryUriForDigest"},
			_jsii_.MemberMethod{JsiiMethod: "repositoryUriForTag", GoMethod: "RepositoryUriForTag"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Repository{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_RepositoryBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.RepositoryAttributes",
		reflect.TypeOf((*RepositoryAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecr.RepositoryBase",
		reflect.TypeOf((*RepositoryBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantPull", GoMethod: "GrantPull"},
			_jsii_.MemberMethod{JsiiMethod: "grantPullPush", GoMethod: "GrantPullPush"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailEvent", GoMethod: "OnCloudTrailEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailImagePushed", GoMethod: "OnCloudTrailImagePushed"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onImageScanCompleted", GoMethod: "OnImageScanCompleted"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryArn", GoGetter: "RepositoryArn"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryUri", GoGetter: "RepositoryUri"},
			_jsii_.MemberMethod{JsiiMethod: "repositoryUriForDigest", GoMethod: "RepositoryUriForDigest"},
			_jsii_.MemberMethod{JsiiMethod: "repositoryUriForTag", GoMethod: "RepositoryUriForTag"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_RepositoryBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IRepository)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecr.RepositoryProps",
		reflect.TypeOf((*RepositoryProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecr.TagMutability",
		reflect.TypeOf((*TagMutability)(nil)).Elem(),
		map[string]interface{}{
			"MUTABLE": TagMutability_MUTABLE,
			"IMMUTABLE": TagMutability_IMMUTABLE,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecr.TagStatus",
		reflect.TypeOf((*TagStatus)(nil)).Elem(),
		map[string]interface{}{
			"ANY": TagStatus_ANY,
			"TAGGED": TagStatus_TAGGED,
			"UNTAGGED": TagStatus_UNTAGGED,
		},
	)
}

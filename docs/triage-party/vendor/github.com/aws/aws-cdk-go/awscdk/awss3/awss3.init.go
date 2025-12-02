package awss3

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_s3.BlockPublicAccess",
		reflect.TypeOf((*BlockPublicAccess)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "blockPublicAcls", GoGetter: "BlockPublicAcls"},
			_jsii_.MemberProperty{JsiiProperty: "blockPublicPolicy", GoGetter: "BlockPublicPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "ignorePublicAcls", GoGetter: "IgnorePublicAcls"},
			_jsii_.MemberProperty{JsiiProperty: "restrictPublicBuckets", GoGetter: "RestrictPublicBuckets"},
		},
		func() interface{} {
			return &jsiiProxy_BlockPublicAccess{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.BlockPublicAccessOptions",
		reflect.TypeOf((*BlockPublicAccessOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.Bucket",
		reflect.TypeOf((*Bucket)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addCorsRule", GoMethod: "AddCorsRule"},
			_jsii_.MemberMethod{JsiiMethod: "addEventNotification", GoMethod: "AddEventNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addInventory", GoMethod: "AddInventory"},
			_jsii_.MemberMethod{JsiiMethod: "addLifecycleRule", GoMethod: "AddLifecycleRule"},
			_jsii_.MemberMethod{JsiiMethod: "addMetric", GoMethod: "AddMetric"},
			_jsii_.MemberMethod{JsiiMethod: "addObjectCreatedNotification", GoMethod: "AddObjectCreatedNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addObjectRemovedNotification", GoMethod: "AddObjectRemovedNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "arnForObjects", GoMethod: "ArnForObjects"},
			_jsii_.MemberProperty{JsiiProperty: "autoCreatePolicy", GoGetter: "AutoCreatePolicy"},
			_jsii_.MemberProperty{JsiiProperty: "bucketArn", GoGetter: "BucketArn"},
			_jsii_.MemberProperty{JsiiProperty: "bucketDomainName", GoGetter: "BucketDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketDualStackDomainName", GoGetter: "BucketDualStackDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketName", GoGetter: "BucketName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketRegionalDomainName", GoGetter: "BucketRegionalDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketWebsiteDomainName", GoGetter: "BucketWebsiteDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketWebsiteUrl", GoGetter: "BucketWebsiteUrl"},
			_jsii_.MemberProperty{JsiiProperty: "disallowPublicAccess", GoGetter: "DisallowPublicAccess"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grantDelete", GoMethod: "GrantDelete"},
			_jsii_.MemberMethod{JsiiMethod: "grantPublicAccess", GoMethod: "GrantPublicAccess"},
			_jsii_.MemberMethod{JsiiMethod: "grantPut", GoMethod: "GrantPut"},
			_jsii_.MemberMethod{JsiiMethod: "grantPutAcl", GoMethod: "GrantPutAcl"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantReadWrite", GoMethod: "GrantReadWrite"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "isWebsite", GoGetter: "IsWebsite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailEvent", GoMethod: "OnCloudTrailEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailPutObject", GoMethod: "OnCloudTrailPutObject"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailWriteObject", GoMethod: "OnCloudTrailWriteObject"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberProperty{JsiiProperty: "policy", GoGetter: "Policy"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "s3UrlForObject", GoMethod: "S3UrlForObject"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "urlForObject", GoMethod: "UrlForObject"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "virtualHostedUrlForObject", GoMethod: "VirtualHostedUrlForObject"},
		},
		func() interface{} {
			j := jsiiProxy_Bucket{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_BucketBase)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.BucketAccessControl",
		reflect.TypeOf((*BucketAccessControl)(nil)).Elem(),
		map[string]interface{}{
			"PRIVATE": BucketAccessControl_PRIVATE,
			"PUBLIC_READ": BucketAccessControl_PUBLIC_READ,
			"PUBLIC_READ_WRITE": BucketAccessControl_PUBLIC_READ_WRITE,
			"AUTHENTICATED_READ": BucketAccessControl_AUTHENTICATED_READ,
			"LOG_DELIVERY_WRITE": BucketAccessControl_LOG_DELIVERY_WRITE,
			"BUCKET_OWNER_READ": BucketAccessControl_BUCKET_OWNER_READ,
			"BUCKET_OWNER_FULL_CONTROL": BucketAccessControl_BUCKET_OWNER_FULL_CONTROL,
			"AWS_EXEC_READ": BucketAccessControl_AWS_EXEC_READ,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.BucketAttributes",
		reflect.TypeOf((*BucketAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.BucketBase",
		reflect.TypeOf((*BucketBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addEventNotification", GoMethod: "AddEventNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addObjectCreatedNotification", GoMethod: "AddObjectCreatedNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addObjectRemovedNotification", GoMethod: "AddObjectRemovedNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "arnForObjects", GoMethod: "ArnForObjects"},
			_jsii_.MemberProperty{JsiiProperty: "autoCreatePolicy", GoGetter: "AutoCreatePolicy"},
			_jsii_.MemberProperty{JsiiProperty: "bucketArn", GoGetter: "BucketArn"},
			_jsii_.MemberProperty{JsiiProperty: "bucketDomainName", GoGetter: "BucketDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketDualStackDomainName", GoGetter: "BucketDualStackDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketName", GoGetter: "BucketName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketRegionalDomainName", GoGetter: "BucketRegionalDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketWebsiteDomainName", GoGetter: "BucketWebsiteDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketWebsiteUrl", GoGetter: "BucketWebsiteUrl"},
			_jsii_.MemberProperty{JsiiProperty: "disallowPublicAccess", GoGetter: "DisallowPublicAccess"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grantDelete", GoMethod: "GrantDelete"},
			_jsii_.MemberMethod{JsiiMethod: "grantPublicAccess", GoMethod: "GrantPublicAccess"},
			_jsii_.MemberMethod{JsiiMethod: "grantPut", GoMethod: "GrantPut"},
			_jsii_.MemberMethod{JsiiMethod: "grantPutAcl", GoMethod: "GrantPutAcl"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantReadWrite", GoMethod: "GrantReadWrite"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "isWebsite", GoGetter: "IsWebsite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailEvent", GoMethod: "OnCloudTrailEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailPutObject", GoMethod: "OnCloudTrailPutObject"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailWriteObject", GoMethod: "OnCloudTrailWriteObject"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberProperty{JsiiProperty: "policy", GoGetter: "Policy"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "s3UrlForObject", GoMethod: "S3UrlForObject"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "urlForObject", GoMethod: "UrlForObject"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "virtualHostedUrlForObject", GoMethod: "VirtualHostedUrlForObject"},
		},
		func() interface{} {
			j := jsiiProxy_BucketBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IBucket)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.BucketEncryption",
		reflect.TypeOf((*BucketEncryption)(nil)).Elem(),
		map[string]interface{}{
			"UNENCRYPTED": BucketEncryption_UNENCRYPTED,
			"KMS_MANAGED": BucketEncryption_KMS_MANAGED,
			"S3_MANAGED": BucketEncryption_S3_MANAGED,
			"KMS": BucketEncryption_KMS,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.BucketMetrics",
		reflect.TypeOf((*BucketMetrics)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.BucketNotificationDestinationConfig",
		reflect.TypeOf((*BucketNotificationDestinationConfig)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.BucketNotificationDestinationType",
		reflect.TypeOf((*BucketNotificationDestinationType)(nil)).Elem(),
		map[string]interface{}{
			"LAMBDA": BucketNotificationDestinationType_LAMBDA,
			"QUEUE": BucketNotificationDestinationType_QUEUE,
			"TOPIC": BucketNotificationDestinationType_TOPIC,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.BucketPolicy",
		reflect.TypeOf((*BucketPolicy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "document", GoGetter: "Document"},
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
			j := jsiiProxy_BucketPolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.BucketPolicyProps",
		reflect.TypeOf((*BucketPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.BucketProps",
		reflect.TypeOf((*BucketProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.CfnAccessPoint",
		reflect.TypeOf((*CfnAccessPoint)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrNetworkOrigin", GoGetter: "AttrNetworkOrigin"},
			_jsii_.MemberProperty{JsiiProperty: "bucket", GoGetter: "Bucket"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
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
			_jsii_.MemberProperty{JsiiProperty: "policy", GoGetter: "Policy"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "publicAccessBlockConfiguration", GoGetter: "PublicAccessBlockConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "vpcConfiguration", GoGetter: "VpcConfiguration"},
		},
		func() interface{} {
			j := jsiiProxy_CfnAccessPoint{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnAccessPoint.PublicAccessBlockConfigurationProperty",
		reflect.TypeOf((*CfnAccessPoint_PublicAccessBlockConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnAccessPoint.VpcConfigurationProperty",
		reflect.TypeOf((*CfnAccessPoint_VpcConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnAccessPointProps",
		reflect.TypeOf((*CfnAccessPointProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.CfnBucket",
		reflect.TypeOf((*CfnBucket)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "accelerateConfiguration", GoGetter: "AccelerateConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "accessControl", GoGetter: "AccessControl"},
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberProperty{JsiiProperty: "analyticsConfigurations", GoGetter: "AnalyticsConfigurations"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrDomainName", GoGetter: "AttrDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "attrDualStackDomainName", GoGetter: "AttrDualStackDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "attrRegionalDomainName", GoGetter: "AttrRegionalDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "attrWebsiteUrl", GoGetter: "AttrWebsiteUrl"},
			_jsii_.MemberProperty{JsiiProperty: "bucketEncryption", GoGetter: "BucketEncryption"},
			_jsii_.MemberProperty{JsiiProperty: "bucketName", GoGetter: "BucketName"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "corsConfiguration", GoGetter: "CorsConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "intelligentTieringConfigurations", GoGetter: "IntelligentTieringConfigurations"},
			_jsii_.MemberProperty{JsiiProperty: "inventoryConfigurations", GoGetter: "InventoryConfigurations"},
			_jsii_.MemberProperty{JsiiProperty: "lifecycleConfiguration", GoGetter: "LifecycleConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "loggingConfiguration", GoGetter: "LoggingConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "metricsConfigurations", GoGetter: "MetricsConfigurations"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "notificationConfiguration", GoGetter: "NotificationConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "objectLockConfiguration", GoGetter: "ObjectLockConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "objectLockEnabled", GoGetter: "ObjectLockEnabled"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "ownershipControls", GoGetter: "OwnershipControls"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "publicAccessBlockConfiguration", GoGetter: "PublicAccessBlockConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "replicationConfiguration", GoGetter: "ReplicationConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "versioningConfiguration", GoGetter: "VersioningConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "websiteConfiguration", GoGetter: "WebsiteConfiguration"},
		},
		func() interface{} {
			j := jsiiProxy_CfnBucket{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.AbortIncompleteMultipartUploadProperty",
		reflect.TypeOf((*CfnBucket_AbortIncompleteMultipartUploadProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.AccelerateConfigurationProperty",
		reflect.TypeOf((*CfnBucket_AccelerateConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.AccessControlTranslationProperty",
		reflect.TypeOf((*CfnBucket_AccessControlTranslationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.AnalyticsConfigurationProperty",
		reflect.TypeOf((*CfnBucket_AnalyticsConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.BucketEncryptionProperty",
		reflect.TypeOf((*CfnBucket_BucketEncryptionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.CorsConfigurationProperty",
		reflect.TypeOf((*CfnBucket_CorsConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.CorsRuleProperty",
		reflect.TypeOf((*CfnBucket_CorsRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.DataExportProperty",
		reflect.TypeOf((*CfnBucket_DataExportProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.DefaultRetentionProperty",
		reflect.TypeOf((*CfnBucket_DefaultRetentionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.DeleteMarkerReplicationProperty",
		reflect.TypeOf((*CfnBucket_DeleteMarkerReplicationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.DestinationProperty",
		reflect.TypeOf((*CfnBucket_DestinationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.EncryptionConfigurationProperty",
		reflect.TypeOf((*CfnBucket_EncryptionConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.FilterRuleProperty",
		reflect.TypeOf((*CfnBucket_FilterRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.IntelligentTieringConfigurationProperty",
		reflect.TypeOf((*CfnBucket_IntelligentTieringConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.InventoryConfigurationProperty",
		reflect.TypeOf((*CfnBucket_InventoryConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.LambdaConfigurationProperty",
		reflect.TypeOf((*CfnBucket_LambdaConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.LifecycleConfigurationProperty",
		reflect.TypeOf((*CfnBucket_LifecycleConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.LoggingConfigurationProperty",
		reflect.TypeOf((*CfnBucket_LoggingConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.MetricsConfigurationProperty",
		reflect.TypeOf((*CfnBucket_MetricsConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.MetricsProperty",
		reflect.TypeOf((*CfnBucket_MetricsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.NoncurrentVersionTransitionProperty",
		reflect.TypeOf((*CfnBucket_NoncurrentVersionTransitionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.NotificationConfigurationProperty",
		reflect.TypeOf((*CfnBucket_NotificationConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.NotificationFilterProperty",
		reflect.TypeOf((*CfnBucket_NotificationFilterProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ObjectLockConfigurationProperty",
		reflect.TypeOf((*CfnBucket_ObjectLockConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ObjectLockRuleProperty",
		reflect.TypeOf((*CfnBucket_ObjectLockRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.OwnershipControlsProperty",
		reflect.TypeOf((*CfnBucket_OwnershipControlsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.OwnershipControlsRuleProperty",
		reflect.TypeOf((*CfnBucket_OwnershipControlsRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.PublicAccessBlockConfigurationProperty",
		reflect.TypeOf((*CfnBucket_PublicAccessBlockConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.QueueConfigurationProperty",
		reflect.TypeOf((*CfnBucket_QueueConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.RedirectAllRequestsToProperty",
		reflect.TypeOf((*CfnBucket_RedirectAllRequestsToProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.RedirectRuleProperty",
		reflect.TypeOf((*CfnBucket_RedirectRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicaModificationsProperty",
		reflect.TypeOf((*CfnBucket_ReplicaModificationsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicationConfigurationProperty",
		reflect.TypeOf((*CfnBucket_ReplicationConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicationDestinationProperty",
		reflect.TypeOf((*CfnBucket_ReplicationDestinationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicationRuleAndOperatorProperty",
		reflect.TypeOf((*CfnBucket_ReplicationRuleAndOperatorProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicationRuleFilterProperty",
		reflect.TypeOf((*CfnBucket_ReplicationRuleFilterProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicationRuleProperty",
		reflect.TypeOf((*CfnBucket_ReplicationRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicationTimeProperty",
		reflect.TypeOf((*CfnBucket_ReplicationTimeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ReplicationTimeValueProperty",
		reflect.TypeOf((*CfnBucket_ReplicationTimeValueProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.RoutingRuleConditionProperty",
		reflect.TypeOf((*CfnBucket_RoutingRuleConditionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.RoutingRuleProperty",
		reflect.TypeOf((*CfnBucket_RoutingRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.RuleProperty",
		reflect.TypeOf((*CfnBucket_RuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.S3KeyFilterProperty",
		reflect.TypeOf((*CfnBucket_S3KeyFilterProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ServerSideEncryptionByDefaultProperty",
		reflect.TypeOf((*CfnBucket_ServerSideEncryptionByDefaultProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.ServerSideEncryptionRuleProperty",
		reflect.TypeOf((*CfnBucket_ServerSideEncryptionRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.SourceSelectionCriteriaProperty",
		reflect.TypeOf((*CfnBucket_SourceSelectionCriteriaProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.SseKmsEncryptedObjectsProperty",
		reflect.TypeOf((*CfnBucket_SseKmsEncryptedObjectsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.StorageClassAnalysisProperty",
		reflect.TypeOf((*CfnBucket_StorageClassAnalysisProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.TagFilterProperty",
		reflect.TypeOf((*CfnBucket_TagFilterProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.TieringProperty",
		reflect.TypeOf((*CfnBucket_TieringProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.TopicConfigurationProperty",
		reflect.TypeOf((*CfnBucket_TopicConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.TransitionProperty",
		reflect.TypeOf((*CfnBucket_TransitionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.VersioningConfigurationProperty",
		reflect.TypeOf((*CfnBucket_VersioningConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucket.WebsiteConfigurationProperty",
		reflect.TypeOf((*CfnBucket_WebsiteConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.CfnBucketPolicy",
		reflect.TypeOf((*CfnBucketPolicy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "bucket", GoGetter: "Bucket"},
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
			_jsii_.MemberProperty{JsiiProperty: "policyDocument", GoGetter: "PolicyDocument"},
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
			j := jsiiProxy_CfnBucketPolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucketPolicyProps",
		reflect.TypeOf((*CfnBucketPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnBucketProps",
		reflect.TypeOf((*CfnBucketProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.CfnStorageLens",
		reflect.TypeOf((*CfnStorageLens)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrStorageLensConfigurationStorageLensArn", GoGetter: "AttrStorageLensConfigurationStorageLensArn"},
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
			_jsii_.MemberProperty{JsiiProperty: "storageLensConfiguration", GoGetter: "StorageLensConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnStorageLens{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.AccountLevelProperty",
		reflect.TypeOf((*CfnStorageLens_AccountLevelProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.ActivityMetricsProperty",
		reflect.TypeOf((*CfnStorageLens_ActivityMetricsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.AwsOrgProperty",
		reflect.TypeOf((*CfnStorageLens_AwsOrgProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.BucketLevelProperty",
		reflect.TypeOf((*CfnStorageLens_BucketLevelProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.BucketsAndRegionsProperty",
		reflect.TypeOf((*CfnStorageLens_BucketsAndRegionsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.DataExportProperty",
		reflect.TypeOf((*CfnStorageLens_DataExportProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.EncryptionProperty",
		reflect.TypeOf((*CfnStorageLens_EncryptionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.PrefixLevelProperty",
		reflect.TypeOf((*CfnStorageLens_PrefixLevelProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.PrefixLevelStorageMetricsProperty",
		reflect.TypeOf((*CfnStorageLens_PrefixLevelStorageMetricsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.S3BucketDestinationProperty",
		reflect.TypeOf((*CfnStorageLens_S3BucketDestinationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.SelectionCriteriaProperty",
		reflect.TypeOf((*CfnStorageLens_SelectionCriteriaProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLens.StorageLensConfigurationProperty",
		reflect.TypeOf((*CfnStorageLens_StorageLensConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CfnStorageLensProps",
		reflect.TypeOf((*CfnStorageLensProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.CorsRule",
		reflect.TypeOf((*CorsRule)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.EventType",
		reflect.TypeOf((*EventType)(nil)).Elem(),
		map[string]interface{}{
			"OBJECT_CREATED": EventType_OBJECT_CREATED,
			"OBJECT_CREATED_PUT": EventType_OBJECT_CREATED_PUT,
			"OBJECT_CREATED_POST": EventType_OBJECT_CREATED_POST,
			"OBJECT_CREATED_COPY": EventType_OBJECT_CREATED_COPY,
			"OBJECT_CREATED_COMPLETE_MULTIPART_UPLOAD": EventType_OBJECT_CREATED_COMPLETE_MULTIPART_UPLOAD,
			"OBJECT_REMOVED": EventType_OBJECT_REMOVED,
			"OBJECT_REMOVED_DELETE": EventType_OBJECT_REMOVED_DELETE,
			"OBJECT_REMOVED_DELETE_MARKER_CREATED": EventType_OBJECT_REMOVED_DELETE_MARKER_CREATED,
			"OBJECT_RESTORE_POST": EventType_OBJECT_RESTORE_POST,
			"OBJECT_RESTORE_COMPLETED": EventType_OBJECT_RESTORE_COMPLETED,
			"REDUCED_REDUNDANCY_LOST_OBJECT": EventType_REDUCED_REDUNDANCY_LOST_OBJECT,
			"REPLICATION_OPERATION_FAILED_REPLICATION": EventType_REPLICATION_OPERATION_FAILED_REPLICATION,
			"REPLICATION_OPERATION_MISSED_THRESHOLD": EventType_REPLICATION_OPERATION_MISSED_THRESHOLD,
			"REPLICATION_OPERATION_REPLICATED_AFTER_THRESHOLD": EventType_REPLICATION_OPERATION_REPLICATED_AFTER_THRESHOLD,
			"REPLICATION_OPERATION_NOT_TRACKED": EventType_REPLICATION_OPERATION_NOT_TRACKED,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.HttpMethods",
		reflect.TypeOf((*HttpMethods)(nil)).Elem(),
		map[string]interface{}{
			"GET": HttpMethods_GET,
			"PUT": HttpMethods_PUT,
			"HEAD": HttpMethods_HEAD,
			"POST": HttpMethods_POST,
			"DELETE": HttpMethods_DELETE,
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_s3.IBucket",
		reflect.TypeOf((*IBucket)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addEventNotification", GoMethod: "AddEventNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addObjectCreatedNotification", GoMethod: "AddObjectCreatedNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addObjectRemovedNotification", GoMethod: "AddObjectRemovedNotification"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "arnForObjects", GoMethod: "ArnForObjects"},
			_jsii_.MemberProperty{JsiiProperty: "bucketArn", GoGetter: "BucketArn"},
			_jsii_.MemberProperty{JsiiProperty: "bucketDomainName", GoGetter: "BucketDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketDualStackDomainName", GoGetter: "BucketDualStackDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketName", GoGetter: "BucketName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketRegionalDomainName", GoGetter: "BucketRegionalDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketWebsiteDomainName", GoGetter: "BucketWebsiteDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "bucketWebsiteUrl", GoGetter: "BucketWebsiteUrl"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grantDelete", GoMethod: "GrantDelete"},
			_jsii_.MemberMethod{JsiiMethod: "grantPublicAccess", GoMethod: "GrantPublicAccess"},
			_jsii_.MemberMethod{JsiiMethod: "grantPut", GoMethod: "GrantPut"},
			_jsii_.MemberMethod{JsiiMethod: "grantPutAcl", GoMethod: "GrantPutAcl"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantReadWrite", GoMethod: "GrantReadWrite"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "isWebsite", GoGetter: "IsWebsite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailEvent", GoMethod: "OnCloudTrailEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailPutObject", GoMethod: "OnCloudTrailPutObject"},
			_jsii_.MemberMethod{JsiiMethod: "onCloudTrailWriteObject", GoMethod: "OnCloudTrailWriteObject"},
			_jsii_.MemberProperty{JsiiProperty: "policy", GoGetter: "Policy"},
			_jsii_.MemberMethod{JsiiMethod: "s3UrlForObject", GoMethod: "S3UrlForObject"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "urlForObject", GoMethod: "UrlForObject"},
			_jsii_.MemberMethod{JsiiMethod: "virtualHostedUrlForObject", GoMethod: "VirtualHostedUrlForObject"},
		},
		func() interface{} {
			j := jsiiProxy_IBucket{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_s3.IBucketNotificationDestination",
		reflect.TypeOf((*IBucketNotificationDestination)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_IBucketNotificationDestination{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.Inventory",
		reflect.TypeOf((*Inventory)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.InventoryDestination",
		reflect.TypeOf((*InventoryDestination)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.InventoryFormat",
		reflect.TypeOf((*InventoryFormat)(nil)).Elem(),
		map[string]interface{}{
			"CSV": InventoryFormat_CSV,
			"PARQUET": InventoryFormat_PARQUET,
			"ORC": InventoryFormat_ORC,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.InventoryFrequency",
		reflect.TypeOf((*InventoryFrequency)(nil)).Elem(),
		map[string]interface{}{
			"DAILY": InventoryFrequency_DAILY,
			"WEEKLY": InventoryFrequency_WEEKLY,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.InventoryObjectVersion",
		reflect.TypeOf((*InventoryObjectVersion)(nil)).Elem(),
		map[string]interface{}{
			"ALL": InventoryObjectVersion_ALL,
			"CURRENT": InventoryObjectVersion_CURRENT,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.LifecycleRule",
		reflect.TypeOf((*LifecycleRule)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.Location",
		reflect.TypeOf((*Location)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.NoncurrentVersionTransition",
		reflect.TypeOf((*NoncurrentVersionTransition)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.NotificationKeyFilter",
		reflect.TypeOf((*NotificationKeyFilter)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.ObjectOwnership",
		reflect.TypeOf((*ObjectOwnership)(nil)).Elem(),
		map[string]interface{}{
			"BUCKET_OWNER_PREFERRED": ObjectOwnership_BUCKET_OWNER_PREFERRED,
			"OBJECT_WRITER": ObjectOwnership_OBJECT_WRITER,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.OnCloudTrailBucketEventOptions",
		reflect.TypeOf((*OnCloudTrailBucketEventOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_s3.RedirectProtocol",
		reflect.TypeOf((*RedirectProtocol)(nil)).Elem(),
		map[string]interface{}{
			"HTTP": RedirectProtocol_HTTP,
			"HTTPS": RedirectProtocol_HTTPS,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.RedirectTarget",
		reflect.TypeOf((*RedirectTarget)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.ReplaceKey",
		reflect.TypeOf((*ReplaceKey)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "prefixWithKey", GoGetter: "PrefixWithKey"},
			_jsii_.MemberProperty{JsiiProperty: "withKey", GoGetter: "WithKey"},
		},
		func() interface{} {
			return &jsiiProxy_ReplaceKey{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.RoutingRule",
		reflect.TypeOf((*RoutingRule)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.RoutingRuleCondition",
		reflect.TypeOf((*RoutingRuleCondition)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_s3.StorageClass",
		reflect.TypeOf((*StorageClass)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
		},
		func() interface{} {
			return &jsiiProxy_StorageClass{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.Transition",
		reflect.TypeOf((*Transition)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_s3.VirtualHostedStyleUrlOptions",
		reflect.TypeOf((*VirtualHostedStyleUrlOptions)(nil)).Elem(),
	)
}

package awssecretsmanager

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.AttachedSecretOptions",
		reflect.TypeOf((*AttachedSecretOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_secretsmanager.AttachmentTargetType",
		reflect.TypeOf((*AttachmentTargetType)(nil)).Elem(),
		map[string]interface{}{
			"INSTANCE": AttachmentTargetType_INSTANCE,
			"CLUSTER": AttachmentTargetType_CLUSTER,
			"RDS_DB_PROXY": AttachmentTargetType_RDS_DB_PROXY,
			"REDSHIFT_CLUSTER": AttachmentTargetType_REDSHIFT_CLUSTER,
			"DOCDB_DB_INSTANCE": AttachmentTargetType_DOCDB_DB_INSTANCE,
			"DOCDB_DB_CLUSTER": AttachmentTargetType_DOCDB_DB_CLUSTER,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.CfnResourcePolicy",
		reflect.TypeOf((*CfnResourcePolicy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "blockPublicPolicy", GoGetter: "BlockPublicPolicy"},
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
			_jsii_.MemberProperty{JsiiProperty: "resourcePolicy", GoGetter: "ResourcePolicy"},
			_jsii_.MemberProperty{JsiiProperty: "secretId", GoGetter: "SecretId"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnResourcePolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnResourcePolicyProps",
		reflect.TypeOf((*CfnResourcePolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.CfnRotationSchedule",
		reflect.TypeOf((*CfnRotationSchedule)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "hostedRotationLambda", GoGetter: "HostedRotationLambda"},
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
			_jsii_.MemberProperty{JsiiProperty: "rotationLambdaArn", GoGetter: "RotationLambdaArn"},
			_jsii_.MemberProperty{JsiiProperty: "rotationRules", GoGetter: "RotationRules"},
			_jsii_.MemberProperty{JsiiProperty: "secretId", GoGetter: "SecretId"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnRotationSchedule{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnRotationSchedule.HostedRotationLambdaProperty",
		reflect.TypeOf((*CfnRotationSchedule_HostedRotationLambdaProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnRotationSchedule.RotationRulesProperty",
		reflect.TypeOf((*CfnRotationSchedule_RotationRulesProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnRotationScheduleProps",
		reflect.TypeOf((*CfnRotationScheduleProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.CfnSecret",
		reflect.TypeOf((*CfnSecret)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "generateSecretString", GoGetter: "GenerateSecretString"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "kmsKeyId", GoGetter: "KmsKeyId"},
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
			_jsii_.MemberProperty{JsiiProperty: "replicaRegions", GoGetter: "ReplicaRegions"},
			_jsii_.MemberProperty{JsiiProperty: "secretString", GoGetter: "SecretString"},
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
			j := jsiiProxy_CfnSecret{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnSecret.GenerateSecretStringProperty",
		reflect.TypeOf((*CfnSecret_GenerateSecretStringProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnSecret.ReplicaRegionProperty",
		reflect.TypeOf((*CfnSecret_ReplicaRegionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnSecretProps",
		reflect.TypeOf((*CfnSecretProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachment",
		reflect.TypeOf((*CfnSecretTargetAttachment)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "secretId", GoGetter: "SecretId"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetId", GoGetter: "TargetId"},
			_jsii_.MemberProperty{JsiiProperty: "targetType", GoGetter: "TargetType"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnSecretTargetAttachment{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachmentProps",
		reflect.TypeOf((*CfnSecretTargetAttachmentProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.HostedRotation",
		reflect.TypeOf((*HostedRotation)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
		},
		func() interface{} {
			j := jsiiProxy_HostedRotation{}
			_jsii_.InitJsiiProxy(&j.Type__awsec2IConnectable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.HostedRotationType",
		reflect.TypeOf((*HostedRotationType)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "isMultiUser", GoGetter: "IsMultiUser"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
		},
		func() interface{} {
			return &jsiiProxy_HostedRotationType{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_secretsmanager.ISecret",
		reflect.TypeOf((*ISecret)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addRotationSchedule", GoMethod: "AddRotationSchedule"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "attach", GoMethod: "Attach"},
			_jsii_.MemberMethod{JsiiMethod: "denyAccountRootDelete", GoMethod: "DenyAccountRootDelete"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "secretArn", GoGetter: "SecretArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretFullArn", GoGetter: "SecretFullArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretName", GoGetter: "SecretName"},
			_jsii_.MemberProperty{JsiiProperty: "secretValue", GoGetter: "SecretValue"},
			_jsii_.MemberMethod{JsiiMethod: "secretValueFromJson", GoMethod: "SecretValueFromJson"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_ISecret{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_secretsmanager.ISecretAttachmentTarget",
		reflect.TypeOf((*ISecretAttachmentTarget)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "asSecretAttachmentTarget", GoMethod: "AsSecretAttachmentTarget"},
		},
		func() interface{} {
			return &jsiiProxy_ISecretAttachmentTarget{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_secretsmanager.ISecretTargetAttachment",
		reflect.TypeOf((*ISecretTargetAttachment)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addRotationSchedule", GoMethod: "AddRotationSchedule"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "attach", GoMethod: "Attach"},
			_jsii_.MemberMethod{JsiiMethod: "denyAccountRootDelete", GoMethod: "DenyAccountRootDelete"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "secretArn", GoGetter: "SecretArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretFullArn", GoGetter: "SecretFullArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretName", GoGetter: "SecretName"},
			_jsii_.MemberProperty{JsiiProperty: "secretTargetAttachmentSecretArn", GoGetter: "SecretTargetAttachmentSecretArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretValue", GoGetter: "SecretValue"},
			_jsii_.MemberMethod{JsiiMethod: "secretValueFromJson", GoMethod: "SecretValueFromJson"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_ISecretTargetAttachment{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ISecret)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.MultiUserHostedRotationOptions",
		reflect.TypeOf((*MultiUserHostedRotationOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.ReplicaRegion",
		reflect.TypeOf((*ReplicaRegion)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.ResourcePolicy",
		reflect.TypeOf((*ResourcePolicy)(nil)).Elem(),
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
			j := jsiiProxy_ResourcePolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.ResourcePolicyProps",
		reflect.TypeOf((*ResourcePolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.RotationSchedule",
		reflect.TypeOf((*RotationSchedule)(nil)).Elem(),
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
			j := jsiiProxy_RotationSchedule{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.RotationScheduleOptions",
		reflect.TypeOf((*RotationScheduleOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.RotationScheduleProps",
		reflect.TypeOf((*RotationScheduleProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.Secret",
		reflect.TypeOf((*Secret)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addReplicaRegion", GoMethod: "AddReplicaRegion"},
			_jsii_.MemberMethod{JsiiMethod: "addRotationSchedule", GoMethod: "AddRotationSchedule"},
			_jsii_.MemberMethod{JsiiMethod: "addTargetAttachment", GoMethod: "AddTargetAttachment"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "arnForPolicies", GoGetter: "ArnForPolicies"},
			_jsii_.MemberMethod{JsiiMethod: "attach", GoMethod: "Attach"},
			_jsii_.MemberProperty{JsiiProperty: "autoCreatePolicy", GoGetter: "AutoCreatePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "denyAccountRootDelete", GoMethod: "DenyAccountRootDelete"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "secretArn", GoGetter: "SecretArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretFullArn", GoGetter: "SecretFullArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretName", GoGetter: "SecretName"},
			_jsii_.MemberProperty{JsiiProperty: "secretValue", GoGetter: "SecretValue"},
			_jsii_.MemberMethod{JsiiMethod: "secretValueFromJson", GoMethod: "SecretValueFromJson"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Secret{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ISecret)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SecretAttachmentTargetProps",
		reflect.TypeOf((*SecretAttachmentTargetProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SecretAttributes",
		reflect.TypeOf((*SecretAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SecretProps",
		reflect.TypeOf((*SecretProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.SecretRotation",
		reflect.TypeOf((*SecretRotation)(nil)).Elem(),
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
			j := jsiiProxy_SecretRotation{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		reflect.TypeOf((*SecretRotationApplication)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applicationArnForPartition", GoMethod: "ApplicationArnForPartition"},
			_jsii_.MemberProperty{JsiiProperty: "applicationId", GoGetter: "ApplicationId"},
			_jsii_.MemberProperty{JsiiProperty: "isMultiUser", GoGetter: "IsMultiUser"},
			_jsii_.MemberProperty{JsiiProperty: "semanticVersion", GoGetter: "SemanticVersion"},
			_jsii_.MemberMethod{JsiiMethod: "semanticVersionForPartition", GoMethod: "SemanticVersionForPartition"},
		},
		func() interface{} {
			return &jsiiProxy_SecretRotationApplication{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SecretRotationApplicationOptions",
		reflect.TypeOf((*SecretRotationApplicationOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SecretRotationProps",
		reflect.TypeOf((*SecretRotationProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SecretStringGenerator",
		reflect.TypeOf((*SecretStringGenerator)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_secretsmanager.SecretTargetAttachment",
		reflect.TypeOf((*SecretTargetAttachment)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addRotationSchedule", GoMethod: "AddRotationSchedule"},
			_jsii_.MemberMethod{JsiiMethod: "addToResourcePolicy", GoMethod: "AddToResourcePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "arnForPolicies", GoGetter: "ArnForPolicies"},
			_jsii_.MemberMethod{JsiiMethod: "attach", GoMethod: "Attach"},
			_jsii_.MemberProperty{JsiiProperty: "autoCreatePolicy", GoGetter: "AutoCreatePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "denyAccountRootDelete", GoMethod: "DenyAccountRootDelete"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "secretArn", GoGetter: "SecretArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretFullArn", GoGetter: "SecretFullArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretName", GoGetter: "SecretName"},
			_jsii_.MemberProperty{JsiiProperty: "secretTargetAttachmentSecretArn", GoGetter: "SecretTargetAttachmentSecretArn"},
			_jsii_.MemberProperty{JsiiProperty: "secretValue", GoGetter: "SecretValue"},
			_jsii_.MemberMethod{JsiiMethod: "secretValueFromJson", GoMethod: "SecretValueFromJson"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_SecretTargetAttachment{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ISecret)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ISecretTargetAttachment)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SecretTargetAttachmentProps",
		reflect.TypeOf((*SecretTargetAttachmentProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_secretsmanager.SingleUserHostedRotationOptions",
		reflect.TypeOf((*SingleUserHostedRotationOptions)(nil)).Elem(),
	)
}

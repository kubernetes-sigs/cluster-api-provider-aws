package awscognito

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterEnum(
		"monocdk.aws_cognito.AccountRecovery",
		reflect.TypeOf((*AccountRecovery)(nil)).Elem(),
		map[string]interface{}{
			"EMAIL_AND_PHONE_WITHOUT_MFA": AccountRecovery_EMAIL_AND_PHONE_WITHOUT_MFA,
			"PHONE_WITHOUT_MFA_AND_EMAIL": AccountRecovery_PHONE_WITHOUT_MFA_AND_EMAIL,
			"EMAIL_ONLY": AccountRecovery_EMAIL_ONLY,
			"PHONE_ONLY_WITHOUT_MFA": AccountRecovery_PHONE_ONLY_WITHOUT_MFA,
			"PHONE_AND_EMAIL": AccountRecovery_PHONE_AND_EMAIL,
			"NONE": AccountRecovery_NONE,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.AttributeMapping",
		reflect.TypeOf((*AttributeMapping)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.AuthFlow",
		reflect.TypeOf((*AuthFlow)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.AutoVerifiedAttrs",
		reflect.TypeOf((*AutoVerifiedAttrs)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.BooleanAttribute",
		reflect.TypeOf((*BooleanAttribute)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_BooleanAttribute{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ICustomAttribute)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnIdentityPool",
		reflect.TypeOf((*CfnIdentityPool)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberProperty{JsiiProperty: "allowClassicFlow", GoGetter: "AllowClassicFlow"},
			_jsii_.MemberProperty{JsiiProperty: "allowUnauthenticatedIdentities", GoGetter: "AllowUnauthenticatedIdentities"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrName", GoGetter: "AttrName"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "cognitoEvents", GoGetter: "CognitoEvents"},
			_jsii_.MemberProperty{JsiiProperty: "cognitoIdentityProviders", GoGetter: "CognitoIdentityProviders"},
			_jsii_.MemberProperty{JsiiProperty: "cognitoStreams", GoGetter: "CognitoStreams"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "developerProviderName", GoGetter: "DeveloperProviderName"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "identityPoolName", GoGetter: "IdentityPoolName"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "openIdConnectProviderArns", GoGetter: "OpenIdConnectProviderArns"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "pushSync", GoGetter: "PushSync"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "samlProviderArns", GoGetter: "SamlProviderArns"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "supportedLoginProviders", GoGetter: "SupportedLoginProviders"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnIdentityPool{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPool.CognitoIdentityProviderProperty",
		reflect.TypeOf((*CfnIdentityPool_CognitoIdentityProviderProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPool.CognitoStreamsProperty",
		reflect.TypeOf((*CfnIdentityPool_CognitoStreamsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPool.PushSyncProperty",
		reflect.TypeOf((*CfnIdentityPool_PushSyncProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPoolProps",
		reflect.TypeOf((*CfnIdentityPoolProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment",
		reflect.TypeOf((*CfnIdentityPoolRoleAttachment)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "identityPoolId", GoGetter: "IdentityPoolId"},
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
			_jsii_.MemberProperty{JsiiProperty: "roleMappings", GoGetter: "RoleMappings"},
			_jsii_.MemberProperty{JsiiProperty: "roles", GoGetter: "Roles"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnIdentityPoolRoleAttachment{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment.MappingRuleProperty",
		reflect.TypeOf((*CfnIdentityPoolRoleAttachment_MappingRuleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment.RoleMappingProperty",
		reflect.TypeOf((*CfnIdentityPoolRoleAttachment_RoleMappingProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment.RulesConfigurationTypeProperty",
		reflect.TypeOf((*CfnIdentityPoolRoleAttachment_RulesConfigurationTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachmentProps",
		reflect.TypeOf((*CfnIdentityPoolRoleAttachmentProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPool",
		reflect.TypeOf((*CfnUserPool)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "accountRecoverySetting", GoGetter: "AccountRecoverySetting"},
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberProperty{JsiiProperty: "adminCreateUserConfig", GoGetter: "AdminCreateUserConfig"},
			_jsii_.MemberProperty{JsiiProperty: "aliasAttributes", GoGetter: "AliasAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrProviderName", GoGetter: "AttrProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "attrProviderUrl", GoGetter: "AttrProviderUrl"},
			_jsii_.MemberProperty{JsiiProperty: "autoVerifiedAttributes", GoGetter: "AutoVerifiedAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "deviceConfiguration", GoGetter: "DeviceConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "emailConfiguration", GoGetter: "EmailConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "emailVerificationMessage", GoGetter: "EmailVerificationMessage"},
			_jsii_.MemberProperty{JsiiProperty: "emailVerificationSubject", GoGetter: "EmailVerificationSubject"},
			_jsii_.MemberProperty{JsiiProperty: "enabledMfas", GoGetter: "EnabledMfas"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "lambdaConfig", GoGetter: "LambdaConfig"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "mfaConfiguration", GoGetter: "MfaConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "policies", GoGetter: "Policies"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "schema", GoGetter: "Schema"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "smsAuthenticationMessage", GoGetter: "SmsAuthenticationMessage"},
			_jsii_.MemberProperty{JsiiProperty: "smsConfiguration", GoGetter: "SmsConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "smsVerificationMessage", GoGetter: "SmsVerificationMessage"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberProperty{JsiiProperty: "usernameAttributes", GoGetter: "UsernameAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "usernameConfiguration", GoGetter: "UsernameConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolAddOns", GoGetter: "UserPoolAddOns"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolName", GoGetter: "UserPoolName"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "verificationMessageTemplate", GoGetter: "VerificationMessageTemplate"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPool{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.AccountRecoverySettingProperty",
		reflect.TypeOf((*CfnUserPool_AccountRecoverySettingProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.AdminCreateUserConfigProperty",
		reflect.TypeOf((*CfnUserPool_AdminCreateUserConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.CustomEmailSenderProperty",
		reflect.TypeOf((*CfnUserPool_CustomEmailSenderProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.CustomSMSSenderProperty",
		reflect.TypeOf((*CfnUserPool_CustomSMSSenderProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.DeviceConfigurationProperty",
		reflect.TypeOf((*CfnUserPool_DeviceConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.EmailConfigurationProperty",
		reflect.TypeOf((*CfnUserPool_EmailConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.InviteMessageTemplateProperty",
		reflect.TypeOf((*CfnUserPool_InviteMessageTemplateProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.LambdaConfigProperty",
		reflect.TypeOf((*CfnUserPool_LambdaConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.NumberAttributeConstraintsProperty",
		reflect.TypeOf((*CfnUserPool_NumberAttributeConstraintsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.PasswordPolicyProperty",
		reflect.TypeOf((*CfnUserPool_PasswordPolicyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.PoliciesProperty",
		reflect.TypeOf((*CfnUserPool_PoliciesProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.RecoveryOptionProperty",
		reflect.TypeOf((*CfnUserPool_RecoveryOptionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.SchemaAttributeProperty",
		reflect.TypeOf((*CfnUserPool_SchemaAttributeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.SmsConfigurationProperty",
		reflect.TypeOf((*CfnUserPool_SmsConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.StringAttributeConstraintsProperty",
		reflect.TypeOf((*CfnUserPool_StringAttributeConstraintsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.UserPoolAddOnsProperty",
		reflect.TypeOf((*CfnUserPool_UserPoolAddOnsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.UsernameConfigurationProperty",
		reflect.TypeOf((*CfnUserPool_UsernameConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPool.VerificationMessageTemplateProperty",
		reflect.TypeOf((*CfnUserPool_VerificationMessageTemplateProperty)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolClient",
		reflect.TypeOf((*CfnUserPoolClient)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "accessTokenValidity", GoGetter: "AccessTokenValidity"},
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberProperty{JsiiProperty: "allowedOAuthFlows", GoGetter: "AllowedOAuthFlows"},
			_jsii_.MemberProperty{JsiiProperty: "allowedOAuthFlowsUserPoolClient", GoGetter: "AllowedOAuthFlowsUserPoolClient"},
			_jsii_.MemberProperty{JsiiProperty: "allowedOAuthScopes", GoGetter: "AllowedOAuthScopes"},
			_jsii_.MemberProperty{JsiiProperty: "analyticsConfiguration", GoGetter: "AnalyticsConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrClientSecret", GoGetter: "AttrClientSecret"},
			_jsii_.MemberProperty{JsiiProperty: "attrName", GoGetter: "AttrName"},
			_jsii_.MemberProperty{JsiiProperty: "callbackUrLs", GoGetter: "CallbackUrLs"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "clientName", GoGetter: "ClientName"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "defaultRedirectUri", GoGetter: "DefaultRedirectUri"},
			_jsii_.MemberProperty{JsiiProperty: "explicitAuthFlows", GoGetter: "ExplicitAuthFlows"},
			_jsii_.MemberProperty{JsiiProperty: "generateSecret", GoGetter: "GenerateSecret"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "idTokenValidity", GoGetter: "IdTokenValidity"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "logoutUrLs", GoGetter: "LogoutUrLs"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "preventUserExistenceErrors", GoGetter: "PreventUserExistenceErrors"},
			_jsii_.MemberProperty{JsiiProperty: "readAttributes", GoGetter: "ReadAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberProperty{JsiiProperty: "refreshTokenValidity", GoGetter: "RefreshTokenValidity"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "supportedIdentityProviders", GoGetter: "SupportedIdentityProviders"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tokenValidityUnits", GoGetter: "TokenValidityUnits"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "writeAttributes", GoGetter: "WriteAttributes"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolClient{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolClient.AnalyticsConfigurationProperty",
		reflect.TypeOf((*CfnUserPoolClient_AnalyticsConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolClient.TokenValidityUnitsProperty",
		reflect.TypeOf((*CfnUserPoolClient_TokenValidityUnitsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolClientProps",
		reflect.TypeOf((*CfnUserPoolClientProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolDomain",
		reflect.TypeOf((*CfnUserPoolDomain)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "customDomainConfig", GoGetter: "CustomDomainConfig"},
			_jsii_.MemberProperty{JsiiProperty: "domain", GoGetter: "Domain"},
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
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolDomain{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolDomain.CustomDomainConfigTypeProperty",
		reflect.TypeOf((*CfnUserPoolDomain_CustomDomainConfigTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolDomainProps",
		reflect.TypeOf((*CfnUserPoolDomainProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolGroup",
		reflect.TypeOf((*CfnUserPoolGroup)(nil)).Elem(),
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
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "groupName", GoGetter: "GroupName"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "precedence", GoGetter: "Precedence"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "roleArn", GoGetter: "RoleArn"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolGroup{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolGroupProps",
		reflect.TypeOf((*CfnUserPoolGroupProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolIdentityProvider",
		reflect.TypeOf((*CfnUserPoolIdentityProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attributeMapping", GoGetter: "AttributeMapping"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "idpIdentifiers", GoGetter: "IdpIdentifiers"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "providerDetails", GoGetter: "ProviderDetails"},
			_jsii_.MemberProperty{JsiiProperty: "providerName", GoGetter: "ProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "providerType", GoGetter: "ProviderType"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolIdentityProvider{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolIdentityProviderProps",
		reflect.TypeOf((*CfnUserPoolIdentityProviderProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolProps",
		reflect.TypeOf((*CfnUserPoolProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolResourceServer",
		reflect.TypeOf((*CfnUserPoolResourceServer)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "identifier", GoGetter: "Identifier"},
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
			_jsii_.MemberProperty{JsiiProperty: "scopes", GoGetter: "Scopes"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolResourceServer{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolResourceServer.ResourceServerScopeTypeProperty",
		reflect.TypeOf((*CfnUserPoolResourceServer_ResourceServerScopeTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolResourceServerProps",
		reflect.TypeOf((*CfnUserPoolResourceServerProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "accountTakeoverRiskConfiguration", GoGetter: "AccountTakeoverRiskConfiguration"},
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
			_jsii_.MemberProperty{JsiiProperty: "clientId", GoGetter: "ClientId"},
			_jsii_.MemberProperty{JsiiProperty: "compromisedCredentialsRiskConfiguration", GoGetter: "CompromisedCredentialsRiskConfiguration"},
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
			_jsii_.MemberProperty{JsiiProperty: "riskExceptionConfiguration", GoGetter: "RiskExceptionConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolRiskConfigurationAttachment{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.AccountTakeoverActionTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_AccountTakeoverActionTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.AccountTakeoverActionsTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_AccountTakeoverActionsTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.AccountTakeoverRiskConfigurationTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_AccountTakeoverRiskConfigurationTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.CompromisedCredentialsActionsTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_CompromisedCredentialsActionsTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.CompromisedCredentialsRiskConfigurationTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_CompromisedCredentialsRiskConfigurationTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.NotifyConfigurationTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_NotifyConfigurationTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.NotifyEmailTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_NotifyEmailTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment.RiskExceptionConfigurationTypeProperty",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachment_RiskExceptionConfigurationTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachmentProps",
		reflect.TypeOf((*CfnUserPoolRiskConfigurationAttachmentProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachment",
		reflect.TypeOf((*CfnUserPoolUICustomizationAttachment)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "clientId", GoGetter: "ClientId"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "css", GoGetter: "Css"},
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
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolUICustomizationAttachment{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachmentProps",
		reflect.TypeOf((*CfnUserPoolUICustomizationAttachmentProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolUser",
		reflect.TypeOf((*CfnUserPoolUser)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "clientMetadata", GoGetter: "ClientMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "desiredDeliveryMediums", GoGetter: "DesiredDeliveryMediums"},
			_jsii_.MemberProperty{JsiiProperty: "forceAliasCreation", GoGetter: "ForceAliasCreation"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "messageAction", GoGetter: "MessageAction"},
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
			_jsii_.MemberProperty{JsiiProperty: "userAttributes", GoGetter: "UserAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "username", GoGetter: "Username"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "validationData", GoGetter: "ValidationData"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolUser{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolUser.AttributeTypeProperty",
		reflect.TypeOf((*CfnUserPoolUser_AttributeTypeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolUserProps",
		reflect.TypeOf((*CfnUserPoolUserProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachment",
		reflect.TypeOf((*CfnUserPoolUserToGroupAttachment)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "groupName", GoGetter: "GroupName"},
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
			_jsii_.MemberProperty{JsiiProperty: "username", GoGetter: "Username"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnUserPoolUserToGroupAttachment{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachmentProps",
		reflect.TypeOf((*CfnUserPoolUserToGroupAttachmentProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.ClientAttributes",
		reflect.TypeOf((*ClientAttributes)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "attributes", GoMethod: "Attributes"},
			_jsii_.MemberMethod{JsiiMethod: "withCustomAttributes", GoMethod: "WithCustomAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "withStandardAttributes", GoMethod: "WithStandardAttributes"},
		},
		func() interface{} {
			return &jsiiProxy_ClientAttributes{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CognitoDomainOptions",
		reflect.TypeOf((*CognitoDomainOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CustomAttributeConfig",
		reflect.TypeOf((*CustomAttributeConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CustomAttributeProps",
		reflect.TypeOf((*CustomAttributeProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.CustomDomainOptions",
		reflect.TypeOf((*CustomDomainOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.DateTimeAttribute",
		reflect.TypeOf((*DateTimeAttribute)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_DateTimeAttribute{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ICustomAttribute)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.EmailSettings",
		reflect.TypeOf((*EmailSettings)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_cognito.ICustomAttribute",
		reflect.TypeOf((*ICustomAttribute)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_ICustomAttribute{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_cognito.IUserPool",
		reflect.TypeOf((*IUserPool)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addClient", GoMethod: "AddClient"},
			_jsii_.MemberMethod{JsiiMethod: "addDomain", GoMethod: "AddDomain"},
			_jsii_.MemberMethod{JsiiMethod: "addResourceServer", GoMethod: "AddResourceServer"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "identityProviders", GoGetter: "IdentityProviders"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "registerIdentityProvider", GoMethod: "RegisterIdentityProvider"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolArn", GoGetter: "UserPoolArn"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
		},
		func() interface{} {
			j := jsiiProxy_IUserPool{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_cognito.IUserPoolClient",
		reflect.TypeOf((*IUserPoolClient)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolClientId", GoGetter: "UserPoolClientId"},
		},
		func() interface{} {
			j := jsiiProxy_IUserPoolClient{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_cognito.IUserPoolDomain",
		reflect.TypeOf((*IUserPoolDomain)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "domainName", GoGetter: "DomainName"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IUserPoolDomain{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_cognito.IUserPoolIdentityProvider",
		reflect.TypeOf((*IUserPoolIdentityProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "providerName", GoGetter: "ProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IUserPoolIdentityProvider{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_cognito.IUserPoolResourceServer",
		reflect.TypeOf((*IUserPoolResourceServer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolResourceServerId", GoGetter: "UserPoolResourceServerId"},
		},
		func() interface{} {
			j := jsiiProxy_IUserPoolResourceServer{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_cognito.Mfa",
		reflect.TypeOf((*Mfa)(nil)).Elem(),
		map[string]interface{}{
			"OFF": Mfa_OFF,
			"OPTIONAL": Mfa_OPTIONAL,
			"REQUIRED": Mfa_REQUIRED,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.MfaSecondFactor",
		reflect.TypeOf((*MfaSecondFactor)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.NumberAttribute",
		reflect.TypeOf((*NumberAttribute)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_NumberAttribute{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ICustomAttribute)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.NumberAttributeConstraints",
		reflect.TypeOf((*NumberAttributeConstraints)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.NumberAttributeProps",
		reflect.TypeOf((*NumberAttributeProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.OAuthFlows",
		reflect.TypeOf((*OAuthFlows)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.OAuthScope",
		reflect.TypeOf((*OAuthScope)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "scopeName", GoGetter: "ScopeName"},
		},
		func() interface{} {
			return &jsiiProxy_OAuthScope{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.OAuthSettings",
		reflect.TypeOf((*OAuthSettings)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.PasswordPolicy",
		reflect.TypeOf((*PasswordPolicy)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.ProviderAttribute",
		reflect.TypeOf((*ProviderAttribute)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "attributeName", GoGetter: "AttributeName"},
		},
		func() interface{} {
			return &jsiiProxy_ProviderAttribute{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.ResourceServerScope",
		reflect.TypeOf((*ResourceServerScope)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "scopeDescription", GoGetter: "ScopeDescription"},
			_jsii_.MemberProperty{JsiiProperty: "scopeName", GoGetter: "ScopeName"},
		},
		func() interface{} {
			return &jsiiProxy_ResourceServerScope{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.ResourceServerScopeProps",
		reflect.TypeOf((*ResourceServerScopeProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.SignInAliases",
		reflect.TypeOf((*SignInAliases)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.SignInUrlOptions",
		reflect.TypeOf((*SignInUrlOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.StandardAttribute",
		reflect.TypeOf((*StandardAttribute)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.StandardAttributes",
		reflect.TypeOf((*StandardAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.StandardAttributesMask",
		reflect.TypeOf((*StandardAttributesMask)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.StringAttribute",
		reflect.TypeOf((*StringAttribute)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_StringAttribute{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ICustomAttribute)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.StringAttributeConstraints",
		reflect.TypeOf((*StringAttributeConstraints)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.StringAttributeProps",
		reflect.TypeOf((*StringAttributeProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserInvitationConfig",
		reflect.TypeOf((*UserInvitationConfig)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPool",
		reflect.TypeOf((*UserPool)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addClient", GoMethod: "AddClient"},
			_jsii_.MemberMethod{JsiiMethod: "addDomain", GoMethod: "AddDomain"},
			_jsii_.MemberMethod{JsiiMethod: "addResourceServer", GoMethod: "AddResourceServer"},
			_jsii_.MemberMethod{JsiiMethod: "addTrigger", GoMethod: "AddTrigger"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "identityProviders", GoGetter: "IdentityProviders"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerIdentityProvider", GoMethod: "RegisterIdentityProvider"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolArn", GoGetter: "UserPoolArn"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolId", GoGetter: "UserPoolId"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolProviderName", GoGetter: "UserPoolProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolProviderUrl", GoGetter: "UserPoolProviderUrl"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPool{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPool)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolClient",
		reflect.TypeOf((*UserPoolClient)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "oAuthFlows", GoGetter: "OAuthFlows"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolClientId", GoGetter: "UserPoolClientId"},
			_jsii_.MemberProperty{JsiiProperty: "userPoolClientName", GoGetter: "UserPoolClientName"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPoolClient{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPoolClient)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolClientIdentityProvider",
		reflect.TypeOf((*UserPoolClientIdentityProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
		},
		func() interface{} {
			return &jsiiProxy_UserPoolClientIdentityProvider{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolClientOptions",
		reflect.TypeOf((*UserPoolClientOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolClientProps",
		reflect.TypeOf((*UserPoolClientProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolDomain",
		reflect.TypeOf((*UserPoolDomain)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "baseUrl", GoMethod: "BaseUrl"},
			_jsii_.MemberProperty{JsiiProperty: "cloudFrontDomainName", GoGetter: "CloudFrontDomainName"},
			_jsii_.MemberProperty{JsiiProperty: "domainName", GoGetter: "DomainName"},
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
			_jsii_.MemberMethod{JsiiMethod: "signInUrl", GoMethod: "SignInUrl"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPoolDomain{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPoolDomain)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolDomainOptions",
		reflect.TypeOf((*UserPoolDomainOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolDomainProps",
		reflect.TypeOf((*UserPoolDomainProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolIdentityProvider",
		reflect.TypeOf((*UserPoolIdentityProvider)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_UserPoolIdentityProvider{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolIdentityProviderAmazon",
		reflect.TypeOf((*UserPoolIdentityProviderAmazon)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "configureAttributeMapping", GoMethod: "ConfigureAttributeMapping"},
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
			_jsii_.MemberProperty{JsiiProperty: "providerName", GoGetter: "ProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPoolIdentityProviderAmazon{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPoolIdentityProvider)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolIdentityProviderAmazonProps",
		reflect.TypeOf((*UserPoolIdentityProviderAmazonProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolIdentityProviderApple",
		reflect.TypeOf((*UserPoolIdentityProviderApple)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "configureAttributeMapping", GoMethod: "ConfigureAttributeMapping"},
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
			_jsii_.MemberProperty{JsiiProperty: "providerName", GoGetter: "ProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPoolIdentityProviderApple{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPoolIdentityProvider)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolIdentityProviderAppleProps",
		reflect.TypeOf((*UserPoolIdentityProviderAppleProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolIdentityProviderFacebook",
		reflect.TypeOf((*UserPoolIdentityProviderFacebook)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "configureAttributeMapping", GoMethod: "ConfigureAttributeMapping"},
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
			_jsii_.MemberProperty{JsiiProperty: "providerName", GoGetter: "ProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPoolIdentityProviderFacebook{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPoolIdentityProvider)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolIdentityProviderFacebookProps",
		reflect.TypeOf((*UserPoolIdentityProviderFacebookProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolIdentityProviderGoogle",
		reflect.TypeOf((*UserPoolIdentityProviderGoogle)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "configureAttributeMapping", GoMethod: "ConfigureAttributeMapping"},
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
			_jsii_.MemberProperty{JsiiProperty: "providerName", GoGetter: "ProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPoolIdentityProviderGoogle{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPoolIdentityProvider)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolIdentityProviderGoogleProps",
		reflect.TypeOf((*UserPoolIdentityProviderGoogleProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolIdentityProviderProps",
		reflect.TypeOf((*UserPoolIdentityProviderProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolOperation",
		reflect.TypeOf((*UserPoolOperation)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "operationName", GoGetter: "OperationName"},
		},
		func() interface{} {
			return &jsiiProxy_UserPoolOperation{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolProps",
		reflect.TypeOf((*UserPoolProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_cognito.UserPoolResourceServer",
		reflect.TypeOf((*UserPoolResourceServer)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "userPoolResourceServerId", GoGetter: "UserPoolResourceServerId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_UserPoolResourceServer{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IUserPoolResourceServer)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolResourceServerOptions",
		reflect.TypeOf((*UserPoolResourceServerOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolResourceServerProps",
		reflect.TypeOf((*UserPoolResourceServerProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserPoolTriggers",
		reflect.TypeOf((*UserPoolTriggers)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_cognito.UserVerificationConfig",
		reflect.TypeOf((*UserVerificationConfig)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_cognito.VerificationEmailStyle",
		reflect.TypeOf((*VerificationEmailStyle)(nil)).Elem(),
		map[string]interface{}{
			"CODE": VerificationEmailStyle_CODE,
			"LINK": VerificationEmailStyle_LINK,
		},
	)
}

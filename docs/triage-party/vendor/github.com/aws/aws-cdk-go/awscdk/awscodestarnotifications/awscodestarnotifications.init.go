package awscodestarnotifications

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_codestarnotifications.CfnNotificationRule",
		reflect.TypeOf((*CfnNotificationRule)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "detailType", GoGetter: "DetailType"},
			_jsii_.MemberProperty{JsiiProperty: "eventTypeIds", GoGetter: "EventTypeIds"},
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
			_jsii_.MemberProperty{JsiiProperty: "resource", GoGetter: "Resource"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "status", GoGetter: "Status"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "targets", GoGetter: "Targets"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnNotificationRule{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codestarnotifications.CfnNotificationRule.TargetProperty",
		reflect.TypeOf((*CfnNotificationRule_TargetProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codestarnotifications.CfnNotificationRuleProps",
		reflect.TypeOf((*CfnNotificationRuleProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codestarnotifications.DetailType",
		reflect.TypeOf((*DetailType)(nil)).Elem(),
		map[string]interface{}{
			"BASIC": DetailType_BASIC,
			"FULL": DetailType_FULL,
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codestarnotifications.INotificationRule",
		reflect.TypeOf((*INotificationRule)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addTarget", GoMethod: "AddTarget"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "notificationRuleArn", GoGetter: "NotificationRuleArn"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_INotificationRule{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codestarnotifications.INotificationRuleSource",
		reflect.TypeOf((*INotificationRuleSource)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bindAsNotificationRuleSource", GoMethod: "BindAsNotificationRuleSource"},
		},
		func() interface{} {
			return &jsiiProxy_INotificationRuleSource{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codestarnotifications.INotificationRuleTarget",
		reflect.TypeOf((*INotificationRuleTarget)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bindAsNotificationRuleTarget", GoMethod: "BindAsNotificationRuleTarget"},
		},
		func() interface{} {
			return &jsiiProxy_INotificationRuleTarget{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codestarnotifications.NotificationRule",
		reflect.TypeOf((*NotificationRule)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addTarget", GoMethod: "AddTarget"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "notificationRuleArn", GoGetter: "NotificationRuleArn"},
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
			j := jsiiProxy_NotificationRule{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INotificationRule)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codestarnotifications.NotificationRuleOptions",
		reflect.TypeOf((*NotificationRuleOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codestarnotifications.NotificationRuleProps",
		reflect.TypeOf((*NotificationRuleProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codestarnotifications.NotificationRuleSourceConfig",
		reflect.TypeOf((*NotificationRuleSourceConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codestarnotifications.NotificationRuleTargetConfig",
		reflect.TypeOf((*NotificationRuleTargetConfig)(nil)).Elem(),
	)
}

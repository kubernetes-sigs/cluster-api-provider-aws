package awselasticloadbalancing

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer",
		reflect.TypeOf((*CfnLoadBalancer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "accessLoggingPolicy", GoGetter: "AccessLoggingPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberProperty{JsiiProperty: "appCookieStickinessPolicy", GoGetter: "AppCookieStickinessPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrCanonicalHostedZoneName", GoGetter: "AttrCanonicalHostedZoneName"},
			_jsii_.MemberProperty{JsiiProperty: "attrCanonicalHostedZoneNameId", GoGetter: "AttrCanonicalHostedZoneNameId"},
			_jsii_.MemberProperty{JsiiProperty: "attrDnsName", GoGetter: "AttrDnsName"},
			_jsii_.MemberProperty{JsiiProperty: "attrSourceSecurityGroupGroupName", GoGetter: "AttrSourceSecurityGroupGroupName"},
			_jsii_.MemberProperty{JsiiProperty: "attrSourceSecurityGroupOwnerAlias", GoGetter: "AttrSourceSecurityGroupOwnerAlias"},
			_jsii_.MemberProperty{JsiiProperty: "availabilityZones", GoGetter: "AvailabilityZones"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "connectionDrainingPolicy", GoGetter: "ConnectionDrainingPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "connectionSettings", GoGetter: "ConnectionSettings"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "crossZone", GoGetter: "CrossZone"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "healthCheck", GoGetter: "HealthCheck"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "instances", GoGetter: "Instances"},
			_jsii_.MemberProperty{JsiiProperty: "lbCookieStickinessPolicy", GoGetter: "LbCookieStickinessPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "listeners", GoGetter: "Listeners"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancerName", GoGetter: "LoadBalancerName"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "policies", GoGetter: "Policies"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "scheme", GoGetter: "Scheme"},
			_jsii_.MemberProperty{JsiiProperty: "securityGroups", GoGetter: "SecurityGroups"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "subnets", GoGetter: "Subnets"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnLoadBalancer{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.AccessLoggingPolicyProperty",
		reflect.TypeOf((*CfnLoadBalancer_AccessLoggingPolicyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.AppCookieStickinessPolicyProperty",
		reflect.TypeOf((*CfnLoadBalancer_AppCookieStickinessPolicyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.ConnectionDrainingPolicyProperty",
		reflect.TypeOf((*CfnLoadBalancer_ConnectionDrainingPolicyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.ConnectionSettingsProperty",
		reflect.TypeOf((*CfnLoadBalancer_ConnectionSettingsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.HealthCheckProperty",
		reflect.TypeOf((*CfnLoadBalancer_HealthCheckProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.LBCookieStickinessPolicyProperty",
		reflect.TypeOf((*CfnLoadBalancer_LBCookieStickinessPolicyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.ListenersProperty",
		reflect.TypeOf((*CfnLoadBalancer_ListenersProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer.PoliciesProperty",
		reflect.TypeOf((*CfnLoadBalancer_PoliciesProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancerProps",
		reflect.TypeOf((*CfnLoadBalancerProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.HealthCheck",
		reflect.TypeOf((*HealthCheck)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_elasticloadbalancing.ILoadBalancerTarget",
		reflect.TypeOf((*ILoadBalancerTarget)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "attachToClassicLB", GoMethod: "AttachToClassicLB"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
		},
		func() interface{} {
			j := jsiiProxy_ILoadBalancerTarget{}
			_jsii_.InitJsiiProxy(&j.Type__awsec2IConnectable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_elasticloadbalancing.ListenerPort",
		reflect.TypeOf((*ListenerPort)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
		},
		func() interface{} {
			j := jsiiProxy_ListenerPort{}
			_jsii_.InitJsiiProxy(&j.Type__awsec2IConnectable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_elasticloadbalancing.LoadBalancer",
		reflect.TypeOf((*LoadBalancer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addListener", GoMethod: "AddListener"},
			_jsii_.MemberMethod{JsiiMethod: "addTarget", GoMethod: "AddTarget"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "listenerPorts", GoGetter: "ListenerPorts"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancerCanonicalHostedZoneName", GoGetter: "LoadBalancerCanonicalHostedZoneName"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancerCanonicalHostedZoneNameId", GoGetter: "LoadBalancerCanonicalHostedZoneNameId"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancerDnsName", GoGetter: "LoadBalancerDnsName"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancerName", GoGetter: "LoadBalancerName"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancerSourceSecurityGroupGroupName", GoGetter: "LoadBalancerSourceSecurityGroupGroupName"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancerSourceSecurityGroupOwnerAlias", GoGetter: "LoadBalancerSourceSecurityGroupOwnerAlias"},
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
			j := jsiiProxy_LoadBalancer{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.Type__awsec2IConnectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.LoadBalancerListener",
		reflect.TypeOf((*LoadBalancerListener)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_elasticloadbalancing.LoadBalancerProps",
		reflect.TypeOf((*LoadBalancerProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_elasticloadbalancing.LoadBalancingProtocol",
		reflect.TypeOf((*LoadBalancingProtocol)(nil)).Elem(),
		map[string]interface{}{
			"TCP": LoadBalancingProtocol_TCP,
			"SSL": LoadBalancingProtocol_SSL,
			"HTTP": LoadBalancingProtocol_HTTP,
			"HTTPS": LoadBalancingProtocol_HTTPS,
		},
	)
}

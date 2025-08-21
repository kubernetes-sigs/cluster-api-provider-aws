package awsecspatterns

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationListenerProps",
		reflect.TypeOf((*ApplicationListenerProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedEc2Service",
		reflect.TypeOf((*ApplicationLoadBalancedEc2Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addServiceAsTarget", GoMethod: "AddServiceAsTarget"},
			_jsii_.MemberProperty{JsiiProperty: "certificate", GoGetter: "Certificate"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "redirectListener", GoGetter: "RedirectListener"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ApplicationLoadBalancedEc2Service{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ApplicationLoadBalancedServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedEc2ServiceProps",
		reflect.TypeOf((*ApplicationLoadBalancedEc2ServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedFargateService",
		reflect.TypeOf((*ApplicationLoadBalancedFargateService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addServiceAsTarget", GoMethod: "AddServiceAsTarget"},
			_jsii_.MemberProperty{JsiiProperty: "assignPublicIp", GoGetter: "AssignPublicIp"},
			_jsii_.MemberProperty{JsiiProperty: "certificate", GoGetter: "Certificate"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "redirectListener", GoGetter: "RedirectListener"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ApplicationLoadBalancedFargateService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ApplicationLoadBalancedServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedFargateServiceProps",
		reflect.TypeOf((*ApplicationLoadBalancedFargateServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedServiceBase",
		reflect.TypeOf((*ApplicationLoadBalancedServiceBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addServiceAsTarget", GoMethod: "AddServiceAsTarget"},
			_jsii_.MemberProperty{JsiiProperty: "certificate", GoGetter: "Certificate"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "redirectListener", GoGetter: "RedirectListener"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ApplicationLoadBalancedServiceBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedServiceBaseProps",
		reflect.TypeOf((*ApplicationLoadBalancedServiceBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedServiceRecordType",
		reflect.TypeOf((*ApplicationLoadBalancedServiceRecordType)(nil)).Elem(),
		map[string]interface{}{
			"ALIAS": ApplicationLoadBalancedServiceRecordType_ALIAS,
			"CNAME": ApplicationLoadBalancedServiceRecordType_CNAME,
			"NONE": ApplicationLoadBalancedServiceRecordType_NONE,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedTaskImageOptions",
		reflect.TypeOf((*ApplicationLoadBalancedTaskImageOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancedTaskImageProps",
		reflect.TypeOf((*ApplicationLoadBalancedTaskImageProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationLoadBalancerProps",
		reflect.TypeOf((*ApplicationLoadBalancerProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ApplicationMultipleTargetGroupsEc2Service",
		reflect.TypeOf((*ApplicationMultipleTargetGroupsEc2Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPortMappingForTargets", GoMethod: "AddPortMappingForTargets"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "findListener", GoMethod: "FindListener"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "listeners", GoGetter: "Listeners"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerECSTargets", GoMethod: "RegisterECSTargets"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroups", GoGetter: "TargetGroups"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ApplicationMultipleTargetGroupsEc2Service{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ApplicationMultipleTargetGroupsServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationMultipleTargetGroupsEc2ServiceProps",
		reflect.TypeOf((*ApplicationMultipleTargetGroupsEc2ServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ApplicationMultipleTargetGroupsFargateService",
		reflect.TypeOf((*ApplicationMultipleTargetGroupsFargateService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPortMappingForTargets", GoMethod: "AddPortMappingForTargets"},
			_jsii_.MemberProperty{JsiiProperty: "assignPublicIp", GoGetter: "AssignPublicIp"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "findListener", GoMethod: "FindListener"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "listeners", GoGetter: "Listeners"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerECSTargets", GoMethod: "RegisterECSTargets"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroups", GoGetter: "TargetGroups"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ApplicationMultipleTargetGroupsFargateService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ApplicationMultipleTargetGroupsServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationMultipleTargetGroupsFargateServiceProps",
		reflect.TypeOf((*ApplicationMultipleTargetGroupsFargateServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ApplicationMultipleTargetGroupsServiceBase",
		reflect.TypeOf((*ApplicationMultipleTargetGroupsServiceBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPortMappingForTargets", GoMethod: "AddPortMappingForTargets"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "findListener", GoMethod: "FindListener"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "listeners", GoGetter: "Listeners"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerECSTargets", GoMethod: "RegisterECSTargets"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroups", GoGetter: "TargetGroups"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ApplicationMultipleTargetGroupsServiceBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationMultipleTargetGroupsServiceBaseProps",
		reflect.TypeOf((*ApplicationMultipleTargetGroupsServiceBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ApplicationTargetProps",
		reflect.TypeOf((*ApplicationTargetProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkListenerProps",
		reflect.TypeOf((*NetworkListenerProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedEc2Service",
		reflect.TypeOf((*NetworkLoadBalancedEc2Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addServiceAsTarget", GoMethod: "AddServiceAsTarget"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_NetworkLoadBalancedEc2Service{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_NetworkLoadBalancedServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedEc2ServiceProps",
		reflect.TypeOf((*NetworkLoadBalancedEc2ServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedFargateService",
		reflect.TypeOf((*NetworkLoadBalancedFargateService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addServiceAsTarget", GoMethod: "AddServiceAsTarget"},
			_jsii_.MemberProperty{JsiiProperty: "assignPublicIp", GoGetter: "AssignPublicIp"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_NetworkLoadBalancedFargateService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_NetworkLoadBalancedServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedFargateServiceProps",
		reflect.TypeOf((*NetworkLoadBalancedFargateServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedServiceBase",
		reflect.TypeOf((*NetworkLoadBalancedServiceBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addServiceAsTarget", GoMethod: "AddServiceAsTarget"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_NetworkLoadBalancedServiceBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedServiceBaseProps",
		reflect.TypeOf((*NetworkLoadBalancedServiceBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedServiceRecordType",
		reflect.TypeOf((*NetworkLoadBalancedServiceRecordType)(nil)).Elem(),
		map[string]interface{}{
			"ALIAS": NetworkLoadBalancedServiceRecordType_ALIAS,
			"CNAME": NetworkLoadBalancedServiceRecordType_CNAME,
			"NONE": NetworkLoadBalancedServiceRecordType_NONE,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedTaskImageOptions",
		reflect.TypeOf((*NetworkLoadBalancedTaskImageOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancedTaskImageProps",
		reflect.TypeOf((*NetworkLoadBalancedTaskImageProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkLoadBalancerProps",
		reflect.TypeOf((*NetworkLoadBalancerProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.NetworkMultipleTargetGroupsEc2Service",
		reflect.TypeOf((*NetworkMultipleTargetGroupsEc2Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPortMappingForTargets", GoMethod: "AddPortMappingForTargets"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "findListener", GoMethod: "FindListener"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "listeners", GoGetter: "Listeners"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerECSTargets", GoMethod: "RegisterECSTargets"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroups", GoGetter: "TargetGroups"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_NetworkMultipleTargetGroupsEc2Service{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_NetworkMultipleTargetGroupsServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkMultipleTargetGroupsEc2ServiceProps",
		reflect.TypeOf((*NetworkMultipleTargetGroupsEc2ServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.NetworkMultipleTargetGroupsFargateService",
		reflect.TypeOf((*NetworkMultipleTargetGroupsFargateService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPortMappingForTargets", GoMethod: "AddPortMappingForTargets"},
			_jsii_.MemberProperty{JsiiProperty: "assignPublicIp", GoGetter: "AssignPublicIp"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "findListener", GoMethod: "FindListener"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "listeners", GoGetter: "Listeners"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerECSTargets", GoMethod: "RegisterECSTargets"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroup", GoGetter: "TargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroups", GoGetter: "TargetGroups"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_NetworkMultipleTargetGroupsFargateService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_NetworkMultipleTargetGroupsServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkMultipleTargetGroupsFargateServiceProps",
		reflect.TypeOf((*NetworkMultipleTargetGroupsFargateServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.NetworkMultipleTargetGroupsServiceBase",
		reflect.TypeOf((*NetworkMultipleTargetGroupsServiceBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPortMappingForTargets", GoMethod: "AddPortMappingForTargets"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberMethod{JsiiMethod: "findListener", GoMethod: "FindListener"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "internalDesiredCount", GoGetter: "InternalDesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "listener", GoGetter: "Listener"},
			_jsii_.MemberProperty{JsiiProperty: "listeners", GoGetter: "Listeners"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancer", GoGetter: "LoadBalancer"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerECSTargets", GoMethod: "RegisterECSTargets"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetGroups", GoGetter: "TargetGroups"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_NetworkMultipleTargetGroupsServiceBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkMultipleTargetGroupsServiceBaseProps",
		reflect.TypeOf((*NetworkMultipleTargetGroupsServiceBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.NetworkTargetProps",
		reflect.TypeOf((*NetworkTargetProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.QueueProcessingEc2Service",
		reflect.TypeOf((*QueueProcessingEc2Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "configureAutoscalingForService", GoMethod: "ConfigureAutoscalingForService"},
			_jsii_.MemberProperty{JsiiProperty: "deadLetterQueue", GoGetter: "DeadLetterQueue"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "environment", GoGetter: "Environment"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberMethod{JsiiMethod: "grantPermissionsToService", GoMethod: "GrantPermissionsToService"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "maxCapacity", GoGetter: "MaxCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "minCapacity", GoGetter: "MinCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "scalingSteps", GoGetter: "ScalingSteps"},
			_jsii_.MemberProperty{JsiiProperty: "secrets", GoGetter: "Secrets"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberProperty{JsiiProperty: "sqsQueue", GoGetter: "SqsQueue"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_QueueProcessingEc2Service{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_QueueProcessingServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.QueueProcessingEc2ServiceProps",
		reflect.TypeOf((*QueueProcessingEc2ServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.QueueProcessingFargateService",
		reflect.TypeOf((*QueueProcessingFargateService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "configureAutoscalingForService", GoMethod: "ConfigureAutoscalingForService"},
			_jsii_.MemberProperty{JsiiProperty: "deadLetterQueue", GoGetter: "DeadLetterQueue"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "environment", GoGetter: "Environment"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberMethod{JsiiMethod: "grantPermissionsToService", GoMethod: "GrantPermissionsToService"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "maxCapacity", GoGetter: "MaxCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "minCapacity", GoGetter: "MinCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "scalingSteps", GoGetter: "ScalingSteps"},
			_jsii_.MemberProperty{JsiiProperty: "secrets", GoGetter: "Secrets"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberProperty{JsiiProperty: "sqsQueue", GoGetter: "SqsQueue"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_QueueProcessingFargateService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_QueueProcessingServiceBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.QueueProcessingFargateServiceProps",
		reflect.TypeOf((*QueueProcessingFargateServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.QueueProcessingServiceBase",
		reflect.TypeOf((*QueueProcessingServiceBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "configureAutoscalingForService", GoMethod: "ConfigureAutoscalingForService"},
			_jsii_.MemberProperty{JsiiProperty: "deadLetterQueue", GoGetter: "DeadLetterQueue"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "environment", GoGetter: "Environment"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberMethod{JsiiMethod: "grantPermissionsToService", GoMethod: "GrantPermissionsToService"},
			_jsii_.MemberProperty{JsiiProperty: "logDriver", GoGetter: "LogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "maxCapacity", GoGetter: "MaxCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "minCapacity", GoGetter: "MinCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "scalingSteps", GoGetter: "ScalingSteps"},
			_jsii_.MemberProperty{JsiiProperty: "secrets", GoGetter: "Secrets"},
			_jsii_.MemberProperty{JsiiProperty: "sqsQueue", GoGetter: "SqsQueue"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_QueueProcessingServiceBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.QueueProcessingServiceBaseProps",
		reflect.TypeOf((*QueueProcessingServiceBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ScheduledEc2Task",
		reflect.TypeOf((*ScheduledEc2Task)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addTaskAsTarget", GoMethod: "AddTaskAsTarget"},
			_jsii_.MemberMethod{JsiiMethod: "addTaskDefinitionToEventTarget", GoMethod: "AddTaskDefinitionToEventTarget"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredTaskCount", GoGetter: "DesiredTaskCount"},
			_jsii_.MemberProperty{JsiiProperty: "eventRule", GoGetter: "EventRule"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "subnetSelection", GoGetter: "SubnetSelection"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "task", GoGetter: "Task"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ScheduledEc2Task{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ScheduledTaskBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledEc2TaskDefinitionOptions",
		reflect.TypeOf((*ScheduledEc2TaskDefinitionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledEc2TaskImageOptions",
		reflect.TypeOf((*ScheduledEc2TaskImageOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledEc2TaskProps",
		reflect.TypeOf((*ScheduledEc2TaskProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ScheduledFargateTask",
		reflect.TypeOf((*ScheduledFargateTask)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addTaskAsTarget", GoMethod: "AddTaskAsTarget"},
			_jsii_.MemberMethod{JsiiMethod: "addTaskDefinitionToEventTarget", GoMethod: "AddTaskDefinitionToEventTarget"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredTaskCount", GoGetter: "DesiredTaskCount"},
			_jsii_.MemberProperty{JsiiProperty: "eventRule", GoGetter: "EventRule"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "subnetSelection", GoGetter: "SubnetSelection"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "task", GoGetter: "Task"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ScheduledFargateTask{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ScheduledTaskBase)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledFargateTaskDefinitionOptions",
		reflect.TypeOf((*ScheduledFargateTaskDefinitionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledFargateTaskImageOptions",
		reflect.TypeOf((*ScheduledFargateTaskImageOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledFargateTaskProps",
		reflect.TypeOf((*ScheduledFargateTaskProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs_patterns.ScheduledTaskBase",
		reflect.TypeOf((*ScheduledTaskBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addTaskAsTarget", GoMethod: "AddTaskAsTarget"},
			_jsii_.MemberMethod{JsiiMethod: "addTaskDefinitionToEventTarget", GoMethod: "AddTaskDefinitionToEventTarget"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "createAWSLogDriver", GoMethod: "CreateAWSLogDriver"},
			_jsii_.MemberProperty{JsiiProperty: "desiredTaskCount", GoGetter: "DesiredTaskCount"},
			_jsii_.MemberProperty{JsiiProperty: "eventRule", GoGetter: "EventRule"},
			_jsii_.MemberMethod{JsiiMethod: "getDefaultCluster", GoMethod: "GetDefaultCluster"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "subnetSelection", GoGetter: "SubnetSelection"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ScheduledTaskBase{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledTaskBaseProps",
		reflect.TypeOf((*ScheduledTaskBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs_patterns.ScheduledTaskImageProps",
		reflect.TypeOf((*ScheduledTaskImageProps)(nil)).Elem(),
	)
}

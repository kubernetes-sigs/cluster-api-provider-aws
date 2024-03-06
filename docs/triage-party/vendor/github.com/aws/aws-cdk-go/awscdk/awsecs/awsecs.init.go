package awsecs

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AddAutoScalingGroupCapacityOptions",
		reflect.TypeOf((*AddAutoScalingGroupCapacityOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AddCapacityOptions",
		reflect.TypeOf((*AddCapacityOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.AmiHardwareType",
		reflect.TypeOf((*AmiHardwareType)(nil)).Elem(),
		map[string]interface{}{
			"STANDARD": AmiHardwareType_STANDARD,
			"GPU": AmiHardwareType_GPU,
			"ARM": AmiHardwareType_ARM,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.AppMeshProxyConfiguration",
		reflect.TypeOf((*AppMeshProxyConfiguration)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_AppMeshProxyConfiguration{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ProxyConfiguration)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AppMeshProxyConfigurationConfigProps",
		reflect.TypeOf((*AppMeshProxyConfigurationConfigProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AppMeshProxyConfigurationProps",
		reflect.TypeOf((*AppMeshProxyConfigurationProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.AsgCapacityProvider",
		reflect.TypeOf((*AsgCapacityProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "autoScalingGroup", GoGetter: "AutoScalingGroup"},
			_jsii_.MemberProperty{JsiiProperty: "capacityProviderName", GoGetter: "CapacityProviderName"},
			_jsii_.MemberProperty{JsiiProperty: "enableManagedTerminationProtection", GoGetter: "EnableManagedTerminationProtection"},
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
			j := jsiiProxy_AsgCapacityProvider{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AsgCapacityProviderProps",
		reflect.TypeOf((*AsgCapacityProviderProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.AssetEnvironmentFile",
		reflect.TypeOf((*AssetEnvironmentFile)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "path", GoGetter: "Path"},
		},
		func() interface{} {
			j := jsiiProxy_AssetEnvironmentFile{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_EnvironmentFile)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.AssetImage",
		reflect.TypeOf((*AssetImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_AssetImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ContainerImage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AssetImageProps",
		reflect.TypeOf((*AssetImageProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AssociateCloudMapServiceOptions",
		reflect.TypeOf((*AssociateCloudMapServiceOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AuthorizationConfig",
		reflect.TypeOf((*AuthorizationConfig)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.AwsLogDriver",
		reflect.TypeOf((*AwsLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "logGroup", GoGetter: "LogGroup"},
		},
		func() interface{} {
			j := jsiiProxy_AwsLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.AwsLogDriverMode",
		reflect.TypeOf((*AwsLogDriverMode)(nil)).Elem(),
		map[string]interface{}{
			"BLOCKING": AwsLogDriverMode_BLOCKING,
			"NON_BLOCKING": AwsLogDriverMode_NON_BLOCKING,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.AwsLogDriverProps",
		reflect.TypeOf((*AwsLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.BaseLogDriverProps",
		reflect.TypeOf((*BaseLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.BaseService",
		reflect.TypeOf((*BaseService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "associateCloudMapService", GoMethod: "AssociateCloudMapService"},
			_jsii_.MemberMethod{JsiiMethod: "attachToApplicationTargetGroup", GoMethod: "AttachToApplicationTargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "attachToClassicLB", GoMethod: "AttachToClassicLB"},
			_jsii_.MemberMethod{JsiiMethod: "attachToNetworkTargetGroup", GoMethod: "AttachToNetworkTargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "autoScaleTaskCount", GoMethod: "AutoScaleTaskCount"},
			_jsii_.MemberProperty{JsiiProperty: "cloudmapService", GoGetter: "CloudmapService"},
			_jsii_.MemberProperty{JsiiProperty: "cloudMapService", GoGetter: "CloudMapService"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "configureAwsVpcNetworking", GoMethod: "ConfigureAwsVpcNetworking"},
			_jsii_.MemberMethod{JsiiMethod: "configureAwsVpcNetworkingWithSecurityGroups", GoMethod: "ConfigureAwsVpcNetworkingWithSecurityGroups"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberMethod{JsiiMethod: "enableCloudMap", GoMethod: "EnableCloudMap"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancers", GoGetter: "LoadBalancers"},
			_jsii_.MemberMethod{JsiiMethod: "loadBalancerTarget", GoMethod: "LoadBalancerTarget"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricCpuUtilization", GoMethod: "MetricCpuUtilization"},
			_jsii_.MemberMethod{JsiiMethod: "metricMemoryUtilization", GoMethod: "MetricMemoryUtilization"},
			_jsii_.MemberProperty{JsiiProperty: "networkConfiguration", GoGetter: "NetworkConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerLoadBalancerTargets", GoMethod: "RegisterLoadBalancerTargets"},
			_jsii_.MemberProperty{JsiiProperty: "serviceArn", GoGetter: "ServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "serviceRegistries", GoGetter: "ServiceRegistries"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_BaseService{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IBaseService)
			_jsii_.InitJsiiProxy(&j.Type__awselasticloadbalancingILoadBalancerTarget)
			_jsii_.InitJsiiProxy(&j.Type__awselasticloadbalancingv2IApplicationLoadBalancerTarget)
			_jsii_.InitJsiiProxy(&j.Type__awselasticloadbalancingv2INetworkLoadBalancerTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.BaseServiceOptions",
		reflect.TypeOf((*BaseServiceOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.BaseServiceProps",
		reflect.TypeOf((*BaseServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.BinPackResource",
		reflect.TypeOf((*BinPackResource)(nil)).Elem(),
		map[string]interface{}{
			"CPU": BinPackResource_CPU,
			"MEMORY": BinPackResource_MEMORY,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.BottleRocketImage",
		reflect.TypeOf((*BottleRocketImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "getImage", GoMethod: "GetImage"},
		},
		func() interface{} {
			j := jsiiProxy_BottleRocketImage{}
			_jsii_.InitJsiiProxy(&j.Type__awsec2IMachineImage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.BottleRocketImageProps",
		reflect.TypeOf((*BottleRocketImageProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.BottlerocketEcsVariant",
		reflect.TypeOf((*BottlerocketEcsVariant)(nil)).Elem(),
		map[string]interface{}{
			"AWS_ECS_1": BottlerocketEcsVariant_AWS_ECS_1,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.BuiltInAttributes",
		reflect.TypeOf((*BuiltInAttributes)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_BuiltInAttributes{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.Capability",
		reflect.TypeOf((*Capability)(nil)).Elem(),
		map[string]interface{}{
			"ALL": Capability_ALL,
			"AUDIT_CONTROL": Capability_AUDIT_CONTROL,
			"AUDIT_WRITE": Capability_AUDIT_WRITE,
			"BLOCK_SUSPEND": Capability_BLOCK_SUSPEND,
			"CHOWN": Capability_CHOWN,
			"DAC_OVERRIDE": Capability_DAC_OVERRIDE,
			"DAC_READ_SEARCH": Capability_DAC_READ_SEARCH,
			"FOWNER": Capability_FOWNER,
			"FSETID": Capability_FSETID,
			"IPC_LOCK": Capability_IPC_LOCK,
			"IPC_OWNER": Capability_IPC_OWNER,
			"KILL": Capability_KILL,
			"LEASE": Capability_LEASE,
			"LINUX_IMMUTABLE": Capability_LINUX_IMMUTABLE,
			"MAC_ADMIN": Capability_MAC_ADMIN,
			"MAC_OVERRIDE": Capability_MAC_OVERRIDE,
			"MKNOD": Capability_MKNOD,
			"NET_ADMIN": Capability_NET_ADMIN,
			"NET_BIND_SERVICE": Capability_NET_BIND_SERVICE,
			"NET_BROADCAST": Capability_NET_BROADCAST,
			"NET_RAW": Capability_NET_RAW,
			"SETFCAP": Capability_SETFCAP,
			"SETGID": Capability_SETGID,
			"SETPCAP": Capability_SETPCAP,
			"SETUID": Capability_SETUID,
			"SYS_ADMIN": Capability_SYS_ADMIN,
			"SYS_BOOT": Capability_SYS_BOOT,
			"SYS_CHROOT": Capability_SYS_CHROOT,
			"SYS_MODULE": Capability_SYS_MODULE,
			"SYS_NICE": Capability_SYS_NICE,
			"SYS_PACCT": Capability_SYS_PACCT,
			"SYS_PTRACE": Capability_SYS_PTRACE,
			"SYS_RAWIO": Capability_SYS_RAWIO,
			"SYS_RESOURCE": Capability_SYS_RESOURCE,
			"SYS_TIME": Capability_SYS_TIME,
			"SYS_TTY_CONFIG": Capability_SYS_TTY_CONFIG,
			"SYSLOG": Capability_SYSLOG,
			"WAKE_ALARM": Capability_WAKE_ALARM,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CapacityProviderStrategy",
		reflect.TypeOf((*CapacityProviderStrategy)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.CfnCapacityProvider",
		reflect.TypeOf((*CfnCapacityProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "autoScalingGroupProvider", GoGetter: "AutoScalingGroupProvider"},
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
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
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
			j := jsiiProxy_CfnCapacityProvider{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCapacityProvider.AutoScalingGroupProviderProperty",
		reflect.TypeOf((*CfnCapacityProvider_AutoScalingGroupProviderProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCapacityProvider.ManagedScalingProperty",
		reflect.TypeOf((*CfnCapacityProvider_ManagedScalingProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCapacityProviderProps",
		reflect.TypeOf((*CfnCapacityProviderProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.CfnCluster",
		reflect.TypeOf((*CfnCluster)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "capacityProviders", GoGetter: "CapacityProviders"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "clusterName", GoGetter: "ClusterName"},
			_jsii_.MemberProperty{JsiiProperty: "clusterSettings", GoGetter: "ClusterSettings"},
			_jsii_.MemberProperty{JsiiProperty: "configuration", GoGetter: "Configuration"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "defaultCapacityProviderStrategy", GoGetter: "DefaultCapacityProviderStrategy"},
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
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnCluster{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCluster.CapacityProviderStrategyItemProperty",
		reflect.TypeOf((*CfnCluster_CapacityProviderStrategyItemProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCluster.ClusterConfigurationProperty",
		reflect.TypeOf((*CfnCluster_ClusterConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCluster.ClusterSettingsProperty",
		reflect.TypeOf((*CfnCluster_ClusterSettingsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCluster.ExecuteCommandConfigurationProperty",
		reflect.TypeOf((*CfnCluster_ExecuteCommandConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnCluster.ExecuteCommandLogConfigurationProperty",
		reflect.TypeOf((*CfnCluster_ExecuteCommandLogConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations",
		reflect.TypeOf((*CfnClusterCapacityProviderAssociations)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "capacityProviders", GoGetter: "CapacityProviders"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "defaultCapacityProviderStrategy", GoGetter: "DefaultCapacityProviderStrategy"},
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
			j := jsiiProxy_CfnClusterCapacityProviderAssociations{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations.CapacityProviderStrategyProperty",
		reflect.TypeOf((*CfnClusterCapacityProviderAssociations_CapacityProviderStrategyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociationsProps",
		reflect.TypeOf((*CfnClusterCapacityProviderAssociationsProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnClusterProps",
		reflect.TypeOf((*CfnClusterProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.CfnPrimaryTaskSet",
		reflect.TypeOf((*CfnPrimaryTaskSet)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
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
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskSetId", GoGetter: "TaskSetId"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnPrimaryTaskSet{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnPrimaryTaskSetProps",
		reflect.TypeOf((*CfnPrimaryTaskSetProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.CfnService",
		reflect.TypeOf((*CfnService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrName", GoGetter: "AttrName"},
			_jsii_.MemberProperty{JsiiProperty: "attrServiceArn", GoGetter: "AttrServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "capacityProviderStrategy", GoGetter: "CapacityProviderStrategy"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "deploymentConfiguration", GoGetter: "DeploymentConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "deploymentController", GoGetter: "DeploymentController"},
			_jsii_.MemberProperty{JsiiProperty: "desiredCount", GoGetter: "DesiredCount"},
			_jsii_.MemberProperty{JsiiProperty: "enableEcsManagedTags", GoGetter: "EnableEcsManagedTags"},
			_jsii_.MemberProperty{JsiiProperty: "enableExecuteCommand", GoGetter: "EnableExecuteCommand"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "healthCheckGracePeriodSeconds", GoGetter: "HealthCheckGracePeriodSeconds"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "launchType", GoGetter: "LaunchType"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancers", GoGetter: "LoadBalancers"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "networkConfiguration", GoGetter: "NetworkConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "placementConstraints", GoGetter: "PlacementConstraints"},
			_jsii_.MemberProperty{JsiiProperty: "placementStrategies", GoGetter: "PlacementStrategies"},
			_jsii_.MemberProperty{JsiiProperty: "platformVersion", GoGetter: "PlatformVersion"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "propagateTags", GoGetter: "PropagateTags"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "role", GoGetter: "Role"},
			_jsii_.MemberProperty{JsiiProperty: "schedulingStrategy", GoGetter: "SchedulingStrategy"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "serviceRegistries", GoGetter: "ServiceRegistries"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnService{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.AwsVpcConfigurationProperty",
		reflect.TypeOf((*CfnService_AwsVpcConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.CapacityProviderStrategyItemProperty",
		reflect.TypeOf((*CfnService_CapacityProviderStrategyItemProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.DeploymentCircuitBreakerProperty",
		reflect.TypeOf((*CfnService_DeploymentCircuitBreakerProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.DeploymentConfigurationProperty",
		reflect.TypeOf((*CfnService_DeploymentConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.DeploymentControllerProperty",
		reflect.TypeOf((*CfnService_DeploymentControllerProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.LoadBalancerProperty",
		reflect.TypeOf((*CfnService_LoadBalancerProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.NetworkConfigurationProperty",
		reflect.TypeOf((*CfnService_NetworkConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.PlacementConstraintProperty",
		reflect.TypeOf((*CfnService_PlacementConstraintProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.PlacementStrategyProperty",
		reflect.TypeOf((*CfnService_PlacementStrategyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnService.ServiceRegistryProperty",
		reflect.TypeOf((*CfnService_ServiceRegistryProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnServiceProps",
		reflect.TypeOf((*CfnServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.CfnTaskDefinition",
		reflect.TypeOf((*CfnTaskDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrTaskDefinitionArn", GoGetter: "AttrTaskDefinitionArn"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "containerDefinitions", GoGetter: "ContainerDefinitions"},
			_jsii_.MemberProperty{JsiiProperty: "cpu", GoGetter: "Cpu"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "ephemeralStorage", GoGetter: "EphemeralStorage"},
			_jsii_.MemberProperty{JsiiProperty: "executionRoleArn", GoGetter: "ExecutionRoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "family", GoGetter: "Family"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "inferenceAccelerators", GoGetter: "InferenceAccelerators"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "ipcMode", GoGetter: "IpcMode"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "memory", GoGetter: "Memory"},
			_jsii_.MemberProperty{JsiiProperty: "networkMode", GoGetter: "NetworkMode"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "pidMode", GoGetter: "PidMode"},
			_jsii_.MemberProperty{JsiiProperty: "placementConstraints", GoGetter: "PlacementConstraints"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "proxyConfiguration", GoGetter: "ProxyConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "requiresCompatibilities", GoGetter: "RequiresCompatibilities"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "taskRoleArn", GoGetter: "TaskRoleArn"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
			_jsii_.MemberProperty{JsiiProperty: "volumes", GoGetter: "Volumes"},
		},
		func() interface{} {
			j := jsiiProxy_CfnTaskDefinition{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.AuthorizationConfigProperty",
		reflect.TypeOf((*CfnTaskDefinition_AuthorizationConfigProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.ContainerDefinitionProperty",
		reflect.TypeOf((*CfnTaskDefinition_ContainerDefinitionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.ContainerDependencyProperty",
		reflect.TypeOf((*CfnTaskDefinition_ContainerDependencyProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.DeviceProperty",
		reflect.TypeOf((*CfnTaskDefinition_DeviceProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.DockerVolumeConfigurationProperty",
		reflect.TypeOf((*CfnTaskDefinition_DockerVolumeConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.EfsVolumeConfigurationProperty",
		reflect.TypeOf((*CfnTaskDefinition_EfsVolumeConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.EnvironmentFileProperty",
		reflect.TypeOf((*CfnTaskDefinition_EnvironmentFileProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.EphemeralStorageProperty",
		reflect.TypeOf((*CfnTaskDefinition_EphemeralStorageProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.FirelensConfigurationProperty",
		reflect.TypeOf((*CfnTaskDefinition_FirelensConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.HealthCheckProperty",
		reflect.TypeOf((*CfnTaskDefinition_HealthCheckProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.HostEntryProperty",
		reflect.TypeOf((*CfnTaskDefinition_HostEntryProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.HostVolumePropertiesProperty",
		reflect.TypeOf((*CfnTaskDefinition_HostVolumePropertiesProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.InferenceAcceleratorProperty",
		reflect.TypeOf((*CfnTaskDefinition_InferenceAcceleratorProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.KernelCapabilitiesProperty",
		reflect.TypeOf((*CfnTaskDefinition_KernelCapabilitiesProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.KeyValuePairProperty",
		reflect.TypeOf((*CfnTaskDefinition_KeyValuePairProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.LinuxParametersProperty",
		reflect.TypeOf((*CfnTaskDefinition_LinuxParametersProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.LogConfigurationProperty",
		reflect.TypeOf((*CfnTaskDefinition_LogConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.MountPointProperty",
		reflect.TypeOf((*CfnTaskDefinition_MountPointProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.PortMappingProperty",
		reflect.TypeOf((*CfnTaskDefinition_PortMappingProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.ProxyConfigurationProperty",
		reflect.TypeOf((*CfnTaskDefinition_ProxyConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.RepositoryCredentialsProperty",
		reflect.TypeOf((*CfnTaskDefinition_RepositoryCredentialsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.ResourceRequirementProperty",
		reflect.TypeOf((*CfnTaskDefinition_ResourceRequirementProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.SecretProperty",
		reflect.TypeOf((*CfnTaskDefinition_SecretProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.SystemControlProperty",
		reflect.TypeOf((*CfnTaskDefinition_SystemControlProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.TaskDefinitionPlacementConstraintProperty",
		reflect.TypeOf((*CfnTaskDefinition_TaskDefinitionPlacementConstraintProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.TmpfsProperty",
		reflect.TypeOf((*CfnTaskDefinition_TmpfsProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.UlimitProperty",
		reflect.TypeOf((*CfnTaskDefinition_UlimitProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.VolumeFromProperty",
		reflect.TypeOf((*CfnTaskDefinition_VolumeFromProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinition.VolumeProperty",
		reflect.TypeOf((*CfnTaskDefinition_VolumeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskDefinitionProps",
		reflect.TypeOf((*CfnTaskDefinitionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.CfnTaskSet",
		reflect.TypeOf((*CfnTaskSet)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrId", GoGetter: "AttrId"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "externalId", GoGetter: "ExternalId"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "launchType", GoGetter: "LaunchType"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancers", GoGetter: "LoadBalancers"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "networkConfiguration", GoGetter: "NetworkConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "platformVersion", GoGetter: "PlatformVersion"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "scale", GoGetter: "Scale"},
			_jsii_.MemberProperty{JsiiProperty: "service", GoGetter: "Service"},
			_jsii_.MemberProperty{JsiiProperty: "serviceRegistries", GoGetter: "ServiceRegistries"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnTaskSet{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskSet.AwsVpcConfigurationProperty",
		reflect.TypeOf((*CfnTaskSet_AwsVpcConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskSet.LoadBalancerProperty",
		reflect.TypeOf((*CfnTaskSet_LoadBalancerProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskSet.NetworkConfigurationProperty",
		reflect.TypeOf((*CfnTaskSet_NetworkConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskSet.ScaleProperty",
		reflect.TypeOf((*CfnTaskSet_ScaleProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskSet.ServiceRegistryProperty",
		reflect.TypeOf((*CfnTaskSet_ServiceRegistryProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CfnTaskSetProps",
		reflect.TypeOf((*CfnTaskSetProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CloudMapNamespaceOptions",
		reflect.TypeOf((*CloudMapNamespaceOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CloudMapOptions",
		reflect.TypeOf((*CloudMapOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.Cluster",
		reflect.TypeOf((*Cluster)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addAsgCapacityProvider", GoMethod: "AddAsgCapacityProvider"},
			_jsii_.MemberMethod{JsiiMethod: "addAutoScalingGroup", GoMethod: "AddAutoScalingGroup"},
			_jsii_.MemberMethod{JsiiMethod: "addCapacity", GoMethod: "AddCapacity"},
			_jsii_.MemberMethod{JsiiMethod: "addCapacityProvider", GoMethod: "AddCapacityProvider"},
			_jsii_.MemberMethod{JsiiMethod: "addDefaultCloudMapNamespace", GoMethod: "AddDefaultCloudMapNamespace"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "autoscalingGroup", GoGetter: "AutoscalingGroup"},
			_jsii_.MemberProperty{JsiiProperty: "clusterArn", GoGetter: "ClusterArn"},
			_jsii_.MemberProperty{JsiiProperty: "clusterName", GoGetter: "ClusterName"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberProperty{JsiiProperty: "defaultCloudMapNamespace", GoGetter: "DefaultCloudMapNamespace"},
			_jsii_.MemberMethod{JsiiMethod: "enableFargateCapacityProviders", GoMethod: "EnableFargateCapacityProviders"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executeCommandConfiguration", GoGetter: "ExecuteCommandConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "hasEc2Capacity", GoGetter: "HasEc2Capacity"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricCpuReservation", GoMethod: "MetricCpuReservation"},
			_jsii_.MemberMethod{JsiiMethod: "metricCpuUtilization", GoMethod: "MetricCpuUtilization"},
			_jsii_.MemberMethod{JsiiMethod: "metricMemoryReservation", GoMethod: "MetricMemoryReservation"},
			_jsii_.MemberMethod{JsiiMethod: "metricMemoryUtilization", GoMethod: "MetricMemoryUtilization"},
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
			_jsii_.MemberProperty{JsiiProperty: "vpc", GoGetter: "Vpc"},
		},
		func() interface{} {
			j := jsiiProxy_Cluster{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ICluster)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ClusterAttributes",
		reflect.TypeOf((*ClusterAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ClusterProps",
		reflect.TypeOf((*ClusterProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CommonTaskDefinitionAttributes",
		reflect.TypeOf((*CommonTaskDefinitionAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CommonTaskDefinitionProps",
		reflect.TypeOf((*CommonTaskDefinitionProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.Compatibility",
		reflect.TypeOf((*Compatibility)(nil)).Elem(),
		map[string]interface{}{
			"EC2": Compatibility_EC2,
			"FARGATE": Compatibility_FARGATE,
			"EC2_AND_FARGATE": Compatibility_EC2_AND_FARGATE,
			"EXTERNAL": Compatibility_EXTERNAL,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.ContainerDefinition",
		reflect.TypeOf((*ContainerDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addContainerDependencies", GoMethod: "AddContainerDependencies"},
			_jsii_.MemberMethod{JsiiMethod: "addInferenceAcceleratorResource", GoMethod: "AddInferenceAcceleratorResource"},
			_jsii_.MemberMethod{JsiiMethod: "addLink", GoMethod: "AddLink"},
			_jsii_.MemberMethod{JsiiMethod: "addMountPoints", GoMethod: "AddMountPoints"},
			_jsii_.MemberMethod{JsiiMethod: "addPortMappings", GoMethod: "AddPortMappings"},
			_jsii_.MemberMethod{JsiiMethod: "addScratch", GoMethod: "AddScratch"},
			_jsii_.MemberMethod{JsiiMethod: "addToExecutionPolicy", GoMethod: "AddToExecutionPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addUlimits", GoMethod: "AddUlimits"},
			_jsii_.MemberMethod{JsiiMethod: "addVolumesFrom", GoMethod: "AddVolumesFrom"},
			_jsii_.MemberProperty{JsiiProperty: "containerDependencies", GoGetter: "ContainerDependencies"},
			_jsii_.MemberProperty{JsiiProperty: "containerName", GoGetter: "ContainerName"},
			_jsii_.MemberProperty{JsiiProperty: "containerPort", GoGetter: "ContainerPort"},
			_jsii_.MemberProperty{JsiiProperty: "environmentFiles", GoGetter: "EnvironmentFiles"},
			_jsii_.MemberProperty{JsiiProperty: "essential", GoGetter: "Essential"},
			_jsii_.MemberMethod{JsiiMethod: "findPortMapping", GoMethod: "FindPortMapping"},
			_jsii_.MemberProperty{JsiiProperty: "ingressPort", GoGetter: "IngressPort"},
			_jsii_.MemberProperty{JsiiProperty: "linuxParameters", GoGetter: "LinuxParameters"},
			_jsii_.MemberProperty{JsiiProperty: "logDriverConfig", GoGetter: "LogDriverConfig"},
			_jsii_.MemberProperty{JsiiProperty: "memoryLimitSpecified", GoGetter: "MemoryLimitSpecified"},
			_jsii_.MemberProperty{JsiiProperty: "mountPoints", GoGetter: "MountPoints"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "portMappings", GoGetter: "PortMappings"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "referencesSecretJsonField", GoGetter: "ReferencesSecretJsonField"},
			_jsii_.MemberMethod{JsiiMethod: "renderContainerDefinition", GoMethod: "RenderContainerDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "ulimits", GoGetter: "Ulimits"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberProperty{JsiiProperty: "volumesFrom", GoGetter: "VolumesFrom"},
		},
		func() interface{} {
			j := jsiiProxy_ContainerDefinition{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ContainerDefinitionOptions",
		reflect.TypeOf((*ContainerDefinitionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ContainerDefinitionProps",
		reflect.TypeOf((*ContainerDefinitionProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ContainerDependency",
		reflect.TypeOf((*ContainerDependency)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.ContainerDependencyCondition",
		reflect.TypeOf((*ContainerDependencyCondition)(nil)).Elem(),
		map[string]interface{}{
			"START": ContainerDependencyCondition_START,
			"COMPLETE": ContainerDependencyCondition_COMPLETE,
			"SUCCESS": ContainerDependencyCondition_SUCCESS,
			"HEALTHY": ContainerDependencyCondition_HEALTHY,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.ContainerImage",
		reflect.TypeOf((*ContainerImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_ContainerImage{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ContainerImageConfig",
		reflect.TypeOf((*ContainerImageConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.CpuUtilizationScalingProps",
		reflect.TypeOf((*CpuUtilizationScalingProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.DeploymentCircuitBreaker",
		reflect.TypeOf((*DeploymentCircuitBreaker)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.DeploymentController",
		reflect.TypeOf((*DeploymentController)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.DeploymentControllerType",
		reflect.TypeOf((*DeploymentControllerType)(nil)).Elem(),
		map[string]interface{}{
			"ECS": DeploymentControllerType_ECS,
			"CODE_DEPLOY": DeploymentControllerType_CODE_DEPLOY,
			"EXTERNAL": DeploymentControllerType_EXTERNAL,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Device",
		reflect.TypeOf((*Device)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.DevicePermission",
		reflect.TypeOf((*DevicePermission)(nil)).Elem(),
		map[string]interface{}{
			"READ": DevicePermission_READ,
			"WRITE": DevicePermission_WRITE,
			"MKNOD": DevicePermission_MKNOD,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.DockerVolumeConfiguration",
		reflect.TypeOf((*DockerVolumeConfiguration)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.Ec2Service",
		reflect.TypeOf((*Ec2Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addPlacementConstraints", GoMethod: "AddPlacementConstraints"},
			_jsii_.MemberMethod{JsiiMethod: "addPlacementStrategies", GoMethod: "AddPlacementStrategies"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "associateCloudMapService", GoMethod: "AssociateCloudMapService"},
			_jsii_.MemberMethod{JsiiMethod: "attachToApplicationTargetGroup", GoMethod: "AttachToApplicationTargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "attachToClassicLB", GoMethod: "AttachToClassicLB"},
			_jsii_.MemberMethod{JsiiMethod: "attachToNetworkTargetGroup", GoMethod: "AttachToNetworkTargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "autoScaleTaskCount", GoMethod: "AutoScaleTaskCount"},
			_jsii_.MemberProperty{JsiiProperty: "cloudmapService", GoGetter: "CloudmapService"},
			_jsii_.MemberProperty{JsiiProperty: "cloudMapService", GoGetter: "CloudMapService"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "configureAwsVpcNetworking", GoMethod: "ConfigureAwsVpcNetworking"},
			_jsii_.MemberMethod{JsiiMethod: "configureAwsVpcNetworkingWithSecurityGroups", GoMethod: "ConfigureAwsVpcNetworkingWithSecurityGroups"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberMethod{JsiiMethod: "enableCloudMap", GoMethod: "EnableCloudMap"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancers", GoGetter: "LoadBalancers"},
			_jsii_.MemberMethod{JsiiMethod: "loadBalancerTarget", GoMethod: "LoadBalancerTarget"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricCpuUtilization", GoMethod: "MetricCpuUtilization"},
			_jsii_.MemberMethod{JsiiMethod: "metricMemoryUtilization", GoMethod: "MetricMemoryUtilization"},
			_jsii_.MemberProperty{JsiiProperty: "networkConfiguration", GoGetter: "NetworkConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerLoadBalancerTargets", GoMethod: "RegisterLoadBalancerTargets"},
			_jsii_.MemberProperty{JsiiProperty: "serviceArn", GoGetter: "ServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "serviceRegistries", GoGetter: "ServiceRegistries"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Ec2Service{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_BaseService)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IEc2Service)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Ec2ServiceAttributes",
		reflect.TypeOf((*Ec2ServiceAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Ec2ServiceProps",
		reflect.TypeOf((*Ec2ServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		reflect.TypeOf((*Ec2TaskDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addContainer", GoMethod: "AddContainer"},
			_jsii_.MemberMethod{JsiiMethod: "addExtension", GoMethod: "AddExtension"},
			_jsii_.MemberMethod{JsiiMethod: "addFirelensLogRouter", GoMethod: "AddFirelensLogRouter"},
			_jsii_.MemberMethod{JsiiMethod: "addInferenceAccelerator", GoMethod: "AddInferenceAccelerator"},
			_jsii_.MemberMethod{JsiiMethod: "addPlacementConstraint", GoMethod: "AddPlacementConstraint"},
			_jsii_.MemberMethod{JsiiMethod: "addToExecutionRolePolicy", GoMethod: "AddToExecutionRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addToTaskRolePolicy", GoMethod: "AddToTaskRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addVolume", GoMethod: "AddVolume"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "compatibility", GoGetter: "Compatibility"},
			_jsii_.MemberProperty{JsiiProperty: "containers", GoGetter: "Containers"},
			_jsii_.MemberProperty{JsiiProperty: "defaultContainer", GoGetter: "DefaultContainer"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executionRole", GoGetter: "ExecutionRole"},
			_jsii_.MemberProperty{JsiiProperty: "family", GoGetter: "Family"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "inferenceAccelerators", GoGetter: "InferenceAccelerators"},
			_jsii_.MemberProperty{JsiiProperty: "isEc2Compatible", GoGetter: "IsEc2Compatible"},
			_jsii_.MemberProperty{JsiiProperty: "isExternalCompatible", GoGetter: "IsExternalCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "isFargateCompatible", GoGetter: "IsFargateCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "networkMode", GoGetter: "NetworkMode"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "obtainExecutionRole", GoMethod: "ObtainExecutionRole"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "referencesSecretJsonField", GoGetter: "ReferencesSecretJsonField"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinitionArn", GoGetter: "TaskDefinitionArn"},
			_jsii_.MemberProperty{JsiiProperty: "taskRole", GoGetter: "TaskRole"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Ec2TaskDefinition{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_TaskDefinition)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IEc2TaskDefinition)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Ec2TaskDefinitionAttributes",
		reflect.TypeOf((*Ec2TaskDefinitionAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Ec2TaskDefinitionProps",
		reflect.TypeOf((*Ec2TaskDefinitionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.EcrImage",
		reflect.TypeOf((*EcrImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "imageName", GoGetter: "ImageName"},
		},
		func() interface{} {
			j := jsiiProxy_EcrImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ContainerImage)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.EcsOptimizedAmi",
		reflect.TypeOf((*EcsOptimizedAmi)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "getImage", GoMethod: "GetImage"},
		},
		func() interface{} {
			j := jsiiProxy_EcsOptimizedAmi{}
			_jsii_.InitJsiiProxy(&j.Type__awsec2IMachineImage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.EcsOptimizedAmiProps",
		reflect.TypeOf((*EcsOptimizedAmiProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.EcsOptimizedImage",
		reflect.TypeOf((*EcsOptimizedImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "getImage", GoMethod: "GetImage"},
		},
		func() interface{} {
			j := jsiiProxy_EcsOptimizedImage{}
			_jsii_.InitJsiiProxy(&j.Type__awsec2IMachineImage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.EcsTarget",
		reflect.TypeOf((*EcsTarget)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.EfsVolumeConfiguration",
		reflect.TypeOf((*EfsVolumeConfiguration)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.EnvironmentFile",
		reflect.TypeOf((*EnvironmentFile)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_EnvironmentFile{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.EnvironmentFileConfig",
		reflect.TypeOf((*EnvironmentFileConfig)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.EnvironmentFileType",
		reflect.TypeOf((*EnvironmentFileType)(nil)).Elem(),
		map[string]interface{}{
			"S3": EnvironmentFileType_S3,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ExecuteCommandConfiguration",
		reflect.TypeOf((*ExecuteCommandConfiguration)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ExecuteCommandLogConfiguration",
		reflect.TypeOf((*ExecuteCommandLogConfiguration)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.ExecuteCommandLogging",
		reflect.TypeOf((*ExecuteCommandLogging)(nil)).Elem(),
		map[string]interface{}{
			"NONE": ExecuteCommandLogging_NONE,
			"DEFAULT": ExecuteCommandLogging_DEFAULT,
			"OVERRIDE": ExecuteCommandLogging_OVERRIDE,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.FargatePlatformVersion",
		reflect.TypeOf((*FargatePlatformVersion)(nil)).Elem(),
		map[string]interface{}{
			"LATEST": FargatePlatformVersion_LATEST,
			"VERSION1_4": FargatePlatformVersion_VERSION1_4,
			"VERSION1_3": FargatePlatformVersion_VERSION1_3,
			"VERSION1_2": FargatePlatformVersion_VERSION1_2,
			"VERSION1_1": FargatePlatformVersion_VERSION1_1,
			"VERSION1_0": FargatePlatformVersion_VERSION1_0,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.FargateService",
		reflect.TypeOf((*FargateService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "associateCloudMapService", GoMethod: "AssociateCloudMapService"},
			_jsii_.MemberMethod{JsiiMethod: "attachToApplicationTargetGroup", GoMethod: "AttachToApplicationTargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "attachToClassicLB", GoMethod: "AttachToClassicLB"},
			_jsii_.MemberMethod{JsiiMethod: "attachToNetworkTargetGroup", GoMethod: "AttachToNetworkTargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "autoScaleTaskCount", GoMethod: "AutoScaleTaskCount"},
			_jsii_.MemberProperty{JsiiProperty: "cloudmapService", GoGetter: "CloudmapService"},
			_jsii_.MemberProperty{JsiiProperty: "cloudMapService", GoGetter: "CloudMapService"},
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberMethod{JsiiMethod: "configureAwsVpcNetworking", GoMethod: "ConfigureAwsVpcNetworking"},
			_jsii_.MemberMethod{JsiiMethod: "configureAwsVpcNetworkingWithSecurityGroups", GoMethod: "ConfigureAwsVpcNetworkingWithSecurityGroups"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberMethod{JsiiMethod: "enableCloudMap", GoMethod: "EnableCloudMap"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "loadBalancers", GoGetter: "LoadBalancers"},
			_jsii_.MemberMethod{JsiiMethod: "loadBalancerTarget", GoMethod: "LoadBalancerTarget"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricCpuUtilization", GoMethod: "MetricCpuUtilization"},
			_jsii_.MemberMethod{JsiiMethod: "metricMemoryUtilization", GoMethod: "MetricMemoryUtilization"},
			_jsii_.MemberProperty{JsiiProperty: "networkConfiguration", GoGetter: "NetworkConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "registerLoadBalancerTargets", GoMethod: "RegisterLoadBalancerTargets"},
			_jsii_.MemberProperty{JsiiProperty: "serviceArn", GoGetter: "ServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "serviceRegistries", GoGetter: "ServiceRegistries"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_FargateService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_BaseService)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IFargateService)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FargateServiceAttributes",
		reflect.TypeOf((*FargateServiceAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FargateServiceProps",
		reflect.TypeOf((*FargateServiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.FargateTaskDefinition",
		reflect.TypeOf((*FargateTaskDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addContainer", GoMethod: "AddContainer"},
			_jsii_.MemberMethod{JsiiMethod: "addExtension", GoMethod: "AddExtension"},
			_jsii_.MemberMethod{JsiiMethod: "addFirelensLogRouter", GoMethod: "AddFirelensLogRouter"},
			_jsii_.MemberMethod{JsiiMethod: "addInferenceAccelerator", GoMethod: "AddInferenceAccelerator"},
			_jsii_.MemberMethod{JsiiMethod: "addPlacementConstraint", GoMethod: "AddPlacementConstraint"},
			_jsii_.MemberMethod{JsiiMethod: "addToExecutionRolePolicy", GoMethod: "AddToExecutionRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addToTaskRolePolicy", GoMethod: "AddToTaskRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addVolume", GoMethod: "AddVolume"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "compatibility", GoGetter: "Compatibility"},
			_jsii_.MemberProperty{JsiiProperty: "containers", GoGetter: "Containers"},
			_jsii_.MemberProperty{JsiiProperty: "defaultContainer", GoGetter: "DefaultContainer"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executionRole", GoGetter: "ExecutionRole"},
			_jsii_.MemberProperty{JsiiProperty: "family", GoGetter: "Family"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "inferenceAccelerators", GoGetter: "InferenceAccelerators"},
			_jsii_.MemberProperty{JsiiProperty: "isEc2Compatible", GoGetter: "IsEc2Compatible"},
			_jsii_.MemberProperty{JsiiProperty: "isExternalCompatible", GoGetter: "IsExternalCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "isFargateCompatible", GoGetter: "IsFargateCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "networkMode", GoGetter: "NetworkMode"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "obtainExecutionRole", GoMethod: "ObtainExecutionRole"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "referencesSecretJsonField", GoGetter: "ReferencesSecretJsonField"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinitionArn", GoGetter: "TaskDefinitionArn"},
			_jsii_.MemberProperty{JsiiProperty: "taskRole", GoGetter: "TaskRole"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_FargateTaskDefinition{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_TaskDefinition)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IFargateTaskDefinition)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FargateTaskDefinitionAttributes",
		reflect.TypeOf((*FargateTaskDefinitionAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FargateTaskDefinitionProps",
		reflect.TypeOf((*FargateTaskDefinitionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.FireLensLogDriver",
		reflect.TypeOf((*FireLensLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_FireLensLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FireLensLogDriverProps",
		reflect.TypeOf((*FireLensLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FirelensConfig",
		reflect.TypeOf((*FirelensConfig)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.FirelensConfigFileType",
		reflect.TypeOf((*FirelensConfigFileType)(nil)).Elem(),
		map[string]interface{}{
			"S3": FirelensConfigFileType_S3,
			"FILE": FirelensConfigFileType_FILE,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.FirelensLogRouter",
		reflect.TypeOf((*FirelensLogRouter)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addContainerDependencies", GoMethod: "AddContainerDependencies"},
			_jsii_.MemberMethod{JsiiMethod: "addInferenceAcceleratorResource", GoMethod: "AddInferenceAcceleratorResource"},
			_jsii_.MemberMethod{JsiiMethod: "addLink", GoMethod: "AddLink"},
			_jsii_.MemberMethod{JsiiMethod: "addMountPoints", GoMethod: "AddMountPoints"},
			_jsii_.MemberMethod{JsiiMethod: "addPortMappings", GoMethod: "AddPortMappings"},
			_jsii_.MemberMethod{JsiiMethod: "addScratch", GoMethod: "AddScratch"},
			_jsii_.MemberMethod{JsiiMethod: "addToExecutionPolicy", GoMethod: "AddToExecutionPolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addUlimits", GoMethod: "AddUlimits"},
			_jsii_.MemberMethod{JsiiMethod: "addVolumesFrom", GoMethod: "AddVolumesFrom"},
			_jsii_.MemberProperty{JsiiProperty: "containerDependencies", GoGetter: "ContainerDependencies"},
			_jsii_.MemberProperty{JsiiProperty: "containerName", GoGetter: "ContainerName"},
			_jsii_.MemberProperty{JsiiProperty: "containerPort", GoGetter: "ContainerPort"},
			_jsii_.MemberProperty{JsiiProperty: "environmentFiles", GoGetter: "EnvironmentFiles"},
			_jsii_.MemberProperty{JsiiProperty: "essential", GoGetter: "Essential"},
			_jsii_.MemberMethod{JsiiMethod: "findPortMapping", GoMethod: "FindPortMapping"},
			_jsii_.MemberProperty{JsiiProperty: "firelensConfig", GoGetter: "FirelensConfig"},
			_jsii_.MemberProperty{JsiiProperty: "ingressPort", GoGetter: "IngressPort"},
			_jsii_.MemberProperty{JsiiProperty: "linuxParameters", GoGetter: "LinuxParameters"},
			_jsii_.MemberProperty{JsiiProperty: "logDriverConfig", GoGetter: "LogDriverConfig"},
			_jsii_.MemberProperty{JsiiProperty: "memoryLimitSpecified", GoGetter: "MemoryLimitSpecified"},
			_jsii_.MemberProperty{JsiiProperty: "mountPoints", GoGetter: "MountPoints"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "portMappings", GoGetter: "PortMappings"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "referencesSecretJsonField", GoGetter: "ReferencesSecretJsonField"},
			_jsii_.MemberMethod{JsiiMethod: "renderContainerDefinition", GoMethod: "RenderContainerDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinition", GoGetter: "TaskDefinition"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "ulimits", GoGetter: "Ulimits"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberProperty{JsiiProperty: "volumesFrom", GoGetter: "VolumesFrom"},
		},
		func() interface{} {
			j := jsiiProxy_FirelensLogRouter{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ContainerDefinition)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FirelensLogRouterDefinitionOptions",
		reflect.TypeOf((*FirelensLogRouterDefinitionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FirelensLogRouterProps",
		reflect.TypeOf((*FirelensLogRouterProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.FirelensLogRouterType",
		reflect.TypeOf((*FirelensLogRouterType)(nil)).Elem(),
		map[string]interface{}{
			"FLUENTBIT": FirelensLogRouterType_FLUENTBIT,
			"FLUENTD": FirelensLogRouterType_FLUENTD,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FirelensOptions",
		reflect.TypeOf((*FirelensOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.FluentdLogDriver",
		reflect.TypeOf((*FluentdLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_FluentdLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.FluentdLogDriverProps",
		reflect.TypeOf((*FluentdLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.GelfCompressionType",
		reflect.TypeOf((*GelfCompressionType)(nil)).Elem(),
		map[string]interface{}{
			"GZIP": GelfCompressionType_GZIP,
			"ZLIB": GelfCompressionType_ZLIB,
			"NONE": GelfCompressionType_NONE,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.GelfLogDriver",
		reflect.TypeOf((*GelfLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_GelfLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.GelfLogDriverProps",
		reflect.TypeOf((*GelfLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.HealthCheck",
		reflect.TypeOf((*HealthCheck)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Host",
		reflect.TypeOf((*Host)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.IBaseService",
		reflect.TypeOf((*IBaseService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "cluster", GoGetter: "Cluster"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "serviceArn", GoGetter: "ServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IBaseService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IService)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.ICluster",
		reflect.TypeOf((*ICluster)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "autoscalingGroup", GoGetter: "AutoscalingGroup"},
			_jsii_.MemberProperty{JsiiProperty: "clusterArn", GoGetter: "ClusterArn"},
			_jsii_.MemberProperty{JsiiProperty: "clusterName", GoGetter: "ClusterName"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
			_jsii_.MemberProperty{JsiiProperty: "defaultCloudMapNamespace", GoGetter: "DefaultCloudMapNamespace"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executeCommandConfiguration", GoGetter: "ExecuteCommandConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "hasEc2Capacity", GoGetter: "HasEc2Capacity"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "vpc", GoGetter: "Vpc"},
		},
		func() interface{} {
			j := jsiiProxy_ICluster{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.IEc2Service",
		reflect.TypeOf((*IEc2Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "serviceArn", GoGetter: "ServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IEc2Service{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IService)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.IEc2TaskDefinition",
		reflect.TypeOf((*IEc2TaskDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "compatibility", GoGetter: "Compatibility"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executionRole", GoGetter: "ExecutionRole"},
			_jsii_.MemberProperty{JsiiProperty: "isEc2Compatible", GoGetter: "IsEc2Compatible"},
			_jsii_.MemberProperty{JsiiProperty: "isExternalCompatible", GoGetter: "IsExternalCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "isFargateCompatible", GoGetter: "IsFargateCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "networkMode", GoGetter: "NetworkMode"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinitionArn", GoGetter: "TaskDefinitionArn"},
			_jsii_.MemberProperty{JsiiProperty: "taskRole", GoGetter: "TaskRole"},
		},
		func() interface{} {
			j := jsiiProxy_IEc2TaskDefinition{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ITaskDefinition)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.IEcsLoadBalancerTarget",
		reflect.TypeOf((*IEcsLoadBalancerTarget)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "attachToApplicationTargetGroup", GoMethod: "AttachToApplicationTargetGroup"},
			_jsii_.MemberMethod{JsiiMethod: "attachToClassicLB", GoMethod: "AttachToClassicLB"},
			_jsii_.MemberMethod{JsiiMethod: "attachToNetworkTargetGroup", GoMethod: "AttachToNetworkTargetGroup"},
			_jsii_.MemberProperty{JsiiProperty: "connections", GoGetter: "Connections"},
		},
		func() interface{} {
			j := jsiiProxy_IEcsLoadBalancerTarget{}
			_jsii_.InitJsiiProxy(&j.Type__awselasticloadbalancingv2IApplicationLoadBalancerTarget)
			_jsii_.InitJsiiProxy(&j.Type__awselasticloadbalancingILoadBalancerTarget)
			_jsii_.InitJsiiProxy(&j.Type__awselasticloadbalancingv2INetworkLoadBalancerTarget)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.IFargateService",
		reflect.TypeOf((*IFargateService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "serviceArn", GoGetter: "ServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IFargateService{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IService)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.IFargateTaskDefinition",
		reflect.TypeOf((*IFargateTaskDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "compatibility", GoGetter: "Compatibility"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executionRole", GoGetter: "ExecutionRole"},
			_jsii_.MemberProperty{JsiiProperty: "isEc2Compatible", GoGetter: "IsEc2Compatible"},
			_jsii_.MemberProperty{JsiiProperty: "isExternalCompatible", GoGetter: "IsExternalCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "isFargateCompatible", GoGetter: "IsFargateCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "networkMode", GoGetter: "NetworkMode"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinitionArn", GoGetter: "TaskDefinitionArn"},
			_jsii_.MemberProperty{JsiiProperty: "taskRole", GoGetter: "TaskRole"},
		},
		func() interface{} {
			j := jsiiProxy_IFargateTaskDefinition{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ITaskDefinition)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.IService",
		reflect.TypeOf((*IService)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "serviceArn", GoGetter: "ServiceArn"},
			_jsii_.MemberProperty{JsiiProperty: "serviceName", GoGetter: "ServiceName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IService{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.ITaskDefinition",
		reflect.TypeOf((*ITaskDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "compatibility", GoGetter: "Compatibility"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executionRole", GoGetter: "ExecutionRole"},
			_jsii_.MemberProperty{JsiiProperty: "isEc2Compatible", GoGetter: "IsEc2Compatible"},
			_jsii_.MemberProperty{JsiiProperty: "isExternalCompatible", GoGetter: "IsExternalCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "isFargateCompatible", GoGetter: "IsFargateCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "networkMode", GoGetter: "NetworkMode"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinitionArn", GoGetter: "TaskDefinitionArn"},
			_jsii_.MemberProperty{JsiiProperty: "taskRole", GoGetter: "TaskRole"},
		},
		func() interface{} {
			j := jsiiProxy_ITaskDefinition{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_ecs.ITaskDefinitionExtension",
		reflect.TypeOf((*ITaskDefinitionExtension)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "extend", GoMethod: "Extend"},
		},
		func() interface{} {
			return &jsiiProxy_ITaskDefinitionExtension{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.InferenceAccelerator",
		reflect.TypeOf((*InferenceAccelerator)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.IpcMode",
		reflect.TypeOf((*IpcMode)(nil)).Elem(),
		map[string]interface{}{
			"NONE": IpcMode_NONE,
			"HOST": IpcMode_HOST,
			"TASK": IpcMode_TASK,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.JournaldLogDriver",
		reflect.TypeOf((*JournaldLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_JournaldLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.JournaldLogDriverProps",
		reflect.TypeOf((*JournaldLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.JsonFileLogDriver",
		reflect.TypeOf((*JsonFileLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_JsonFileLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.JsonFileLogDriverProps",
		reflect.TypeOf((*JsonFileLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.LaunchType",
		reflect.TypeOf((*LaunchType)(nil)).Elem(),
		map[string]interface{}{
			"EC2": LaunchType_EC2,
			"FARGATE": LaunchType_FARGATE,
			"EXTERNAL": LaunchType_EXTERNAL,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.LinuxParameters",
		reflect.TypeOf((*LinuxParameters)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addCapabilities", GoMethod: "AddCapabilities"},
			_jsii_.MemberMethod{JsiiMethod: "addDevices", GoMethod: "AddDevices"},
			_jsii_.MemberMethod{JsiiMethod: "addTmpfs", GoMethod: "AddTmpfs"},
			_jsii_.MemberMethod{JsiiMethod: "dropCapabilities", GoMethod: "DropCapabilities"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderLinuxParameters", GoMethod: "RenderLinuxParameters"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_LinuxParameters{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.LinuxParametersProps",
		reflect.TypeOf((*LinuxParametersProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.ListenerConfig",
		reflect.TypeOf((*ListenerConfig)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addTargets", GoMethod: "AddTargets"},
		},
		func() interface{} {
			return &jsiiProxy_ListenerConfig{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.LoadBalancerTargetOptions",
		reflect.TypeOf((*LoadBalancerTargetOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.LogDriver",
		reflect.TypeOf((*LogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_LogDriver{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.LogDriverConfig",
		reflect.TypeOf((*LogDriverConfig)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.LogDrivers",
		reflect.TypeOf((*LogDrivers)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_LogDrivers{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.MachineImageType",
		reflect.TypeOf((*MachineImageType)(nil)).Elem(),
		map[string]interface{}{
			"AMAZON_LINUX_2": MachineImageType_AMAZON_LINUX_2,
			"BOTTLEROCKET": MachineImageType_BOTTLEROCKET,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.MemoryUtilizationScalingProps",
		reflect.TypeOf((*MemoryUtilizationScalingProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.MountPoint",
		reflect.TypeOf((*MountPoint)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.NetworkMode",
		reflect.TypeOf((*NetworkMode)(nil)).Elem(),
		map[string]interface{}{
			"NONE": NetworkMode_NONE,
			"BRIDGE": NetworkMode_BRIDGE,
			"AWS_VPC": NetworkMode_AWS_VPC,
			"HOST": NetworkMode_HOST,
			"NAT": NetworkMode_NAT,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.PidMode",
		reflect.TypeOf((*PidMode)(nil)).Elem(),
		map[string]interface{}{
			"HOST": PidMode_HOST,
			"TASK": PidMode_TASK,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.PlacementConstraint",
		reflect.TypeOf((*PlacementConstraint)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "toJson", GoMethod: "ToJson"},
		},
		func() interface{} {
			return &jsiiProxy_PlacementConstraint{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.PlacementStrategy",
		reflect.TypeOf((*PlacementStrategy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "toJson", GoMethod: "ToJson"},
		},
		func() interface{} {
			return &jsiiProxy_PlacementStrategy{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.PortMapping",
		reflect.TypeOf((*PortMapping)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.PropagatedTagSource",
		reflect.TypeOf((*PropagatedTagSource)(nil)).Elem(),
		map[string]interface{}{
			"SERVICE": PropagatedTagSource_SERVICE,
			"TASK_DEFINITION": PropagatedTagSource_TASK_DEFINITION,
			"NONE": PropagatedTagSource_NONE,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.Protocol",
		reflect.TypeOf((*Protocol)(nil)).Elem(),
		map[string]interface{}{
			"TCP": Protocol_TCP,
			"UDP": Protocol_UDP,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.ProxyConfiguration",
		reflect.TypeOf((*ProxyConfiguration)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_ProxyConfiguration{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.ProxyConfigurations",
		reflect.TypeOf((*ProxyConfigurations)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_ProxyConfigurations{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.RepositoryImage",
		reflect.TypeOf((*RepositoryImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_RepositoryImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ContainerImage)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.RepositoryImageProps",
		reflect.TypeOf((*RepositoryImageProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.RequestCountScalingProps",
		reflect.TypeOf((*RequestCountScalingProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.S3EnvironmentFile",
		reflect.TypeOf((*S3EnvironmentFile)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_S3EnvironmentFile{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_EnvironmentFile)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.ScalableTaskCount",
		reflect.TypeOf((*ScalableTaskCount)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "doScaleOnMetric", GoMethod: "DoScaleOnMetric"},
			_jsii_.MemberMethod{JsiiMethod: "doScaleOnSchedule", GoMethod: "DoScaleOnSchedule"},
			_jsii_.MemberMethod{JsiiMethod: "doScaleToTrackMetric", GoMethod: "DoScaleToTrackMetric"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "props", GoGetter: "Props"},
			_jsii_.MemberMethod{JsiiMethod: "scaleOnCpuUtilization", GoMethod: "ScaleOnCpuUtilization"},
			_jsii_.MemberMethod{JsiiMethod: "scaleOnMemoryUtilization", GoMethod: "ScaleOnMemoryUtilization"},
			_jsii_.MemberMethod{JsiiMethod: "scaleOnMetric", GoMethod: "ScaleOnMetric"},
			_jsii_.MemberMethod{JsiiMethod: "scaleOnRequestCount", GoMethod: "ScaleOnRequestCount"},
			_jsii_.MemberMethod{JsiiMethod: "scaleOnSchedule", GoMethod: "ScaleOnSchedule"},
			_jsii_.MemberMethod{JsiiMethod: "scaleToTrackCustomMetric", GoMethod: "ScaleToTrackCustomMetric"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ScalableTaskCount{}
			_jsii_.InitJsiiProxy(&j.Type__awsapplicationautoscalingBaseScalableAttribute)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ScalableTaskCountProps",
		reflect.TypeOf((*ScalableTaskCountProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.Scope",
		reflect.TypeOf((*Scope)(nil)).Elem(),
		map[string]interface{}{
			"TASK": Scope_TASK,
			"SHARED": Scope_SHARED,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.ScratchSpace",
		reflect.TypeOf((*ScratchSpace)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.Secret",
		reflect.TypeOf((*Secret)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "arn", GoGetter: "Arn"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberProperty{JsiiProperty: "hasField", GoGetter: "HasField"},
		},
		func() interface{} {
			return &jsiiProxy_Secret{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.SplunkLogDriver",
		reflect.TypeOf((*SplunkLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_SplunkLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.SplunkLogDriverProps",
		reflect.TypeOf((*SplunkLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.SplunkLogFormat",
		reflect.TypeOf((*SplunkLogFormat)(nil)).Elem(),
		map[string]interface{}{
			"INLINE": SplunkLogFormat_INLINE,
			"JSON": SplunkLogFormat_JSON,
			"RAW": SplunkLogFormat_RAW,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.SyslogLogDriver",
		reflect.TypeOf((*SyslogLogDriver)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_SyslogLogDriver{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_LogDriver)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.SyslogLogDriverProps",
		reflect.TypeOf((*SyslogLogDriverProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.TagParameterContainerImage",
		reflect.TypeOf((*TagParameterContainerImage)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "tagParameterName", GoGetter: "TagParameterName"},
			_jsii_.MemberProperty{JsiiProperty: "tagParameterValue", GoGetter: "TagParameterValue"},
		},
		func() interface{} {
			j := jsiiProxy_TagParameterContainerImage{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ContainerImage)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_ecs.TaskDefinition",
		reflect.TypeOf((*TaskDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addContainer", GoMethod: "AddContainer"},
			_jsii_.MemberMethod{JsiiMethod: "addExtension", GoMethod: "AddExtension"},
			_jsii_.MemberMethod{JsiiMethod: "addFirelensLogRouter", GoMethod: "AddFirelensLogRouter"},
			_jsii_.MemberMethod{JsiiMethod: "addInferenceAccelerator", GoMethod: "AddInferenceAccelerator"},
			_jsii_.MemberMethod{JsiiMethod: "addPlacementConstraint", GoMethod: "AddPlacementConstraint"},
			_jsii_.MemberMethod{JsiiMethod: "addToExecutionRolePolicy", GoMethod: "AddToExecutionRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addToTaskRolePolicy", GoMethod: "AddToTaskRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "addVolume", GoMethod: "AddVolume"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "compatibility", GoGetter: "Compatibility"},
			_jsii_.MemberProperty{JsiiProperty: "containers", GoGetter: "Containers"},
			_jsii_.MemberProperty{JsiiProperty: "defaultContainer", GoGetter: "DefaultContainer"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "executionRole", GoGetter: "ExecutionRole"},
			_jsii_.MemberProperty{JsiiProperty: "family", GoGetter: "Family"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "inferenceAccelerators", GoGetter: "InferenceAccelerators"},
			_jsii_.MemberProperty{JsiiProperty: "isEc2Compatible", GoGetter: "IsEc2Compatible"},
			_jsii_.MemberProperty{JsiiProperty: "isExternalCompatible", GoGetter: "IsExternalCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "isFargateCompatible", GoGetter: "IsFargateCompatible"},
			_jsii_.MemberProperty{JsiiProperty: "networkMode", GoGetter: "NetworkMode"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "obtainExecutionRole", GoMethod: "ObtainExecutionRole"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "referencesSecretJsonField", GoGetter: "ReferencesSecretJsonField"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskDefinitionArn", GoGetter: "TaskDefinitionArn"},
			_jsii_.MemberProperty{JsiiProperty: "taskRole", GoGetter: "TaskRole"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_TaskDefinition{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_ITaskDefinition)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.TaskDefinitionAttributes",
		reflect.TypeOf((*TaskDefinitionAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.TaskDefinitionProps",
		reflect.TypeOf((*TaskDefinitionProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Tmpfs",
		reflect.TypeOf((*Tmpfs)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.TmpfsMountOption",
		reflect.TypeOf((*TmpfsMountOption)(nil)).Elem(),
		map[string]interface{}{
			"DEFAULTS": TmpfsMountOption_DEFAULTS,
			"RO": TmpfsMountOption_RO,
			"RW": TmpfsMountOption_RW,
			"SUID": TmpfsMountOption_SUID,
			"NOSUID": TmpfsMountOption_NOSUID,
			"DEV": TmpfsMountOption_DEV,
			"NODEV": TmpfsMountOption_NODEV,
			"EXEC": TmpfsMountOption_EXEC,
			"NOEXEC": TmpfsMountOption_NOEXEC,
			"SYNC": TmpfsMountOption_SYNC,
			"ASYNC": TmpfsMountOption_ASYNC,
			"DIRSYNC": TmpfsMountOption_DIRSYNC,
			"REMOUNT": TmpfsMountOption_REMOUNT,
			"MAND": TmpfsMountOption_MAND,
			"NOMAND": TmpfsMountOption_NOMAND,
			"ATIME": TmpfsMountOption_ATIME,
			"NOATIME": TmpfsMountOption_NOATIME,
			"DIRATIME": TmpfsMountOption_DIRATIME,
			"NODIRATIME": TmpfsMountOption_NODIRATIME,
			"BIND": TmpfsMountOption_BIND,
			"RBIND": TmpfsMountOption_RBIND,
			"UNBINDABLE": TmpfsMountOption_UNBINDABLE,
			"RUNBINDABLE": TmpfsMountOption_RUNBINDABLE,
			"PRIVATE": TmpfsMountOption_PRIVATE,
			"RPRIVATE": TmpfsMountOption_RPRIVATE,
			"SHARED": TmpfsMountOption_SHARED,
			"RSHARED": TmpfsMountOption_RSHARED,
			"SLAVE": TmpfsMountOption_SLAVE,
			"RSLAVE": TmpfsMountOption_RSLAVE,
			"RELATIME": TmpfsMountOption_RELATIME,
			"NORELATIME": TmpfsMountOption_NORELATIME,
			"STRICTATIME": TmpfsMountOption_STRICTATIME,
			"NOSTRICTATIME": TmpfsMountOption_NOSTRICTATIME,
			"MODE": TmpfsMountOption_MODE,
			"UID": TmpfsMountOption_UID,
			"GID": TmpfsMountOption_GID,
			"NR_INODES": TmpfsMountOption_NR_INODES,
			"NR_BLOCKS": TmpfsMountOption_NR_BLOCKS,
			"MPOL": TmpfsMountOption_MPOL,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.TrackCustomMetricProps",
		reflect.TypeOf((*TrackCustomMetricProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Ulimit",
		reflect.TypeOf((*Ulimit)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.UlimitName",
		reflect.TypeOf((*UlimitName)(nil)).Elem(),
		map[string]interface{}{
			"CORE": UlimitName_CORE,
			"CPU": UlimitName_CPU,
			"DATA": UlimitName_DATA,
			"FSIZE": UlimitName_FSIZE,
			"LOCKS": UlimitName_LOCKS,
			"MEMLOCK": UlimitName_MEMLOCK,
			"MSGQUEUE": UlimitName_MSGQUEUE,
			"NICE": UlimitName_NICE,
			"NOFILE": UlimitName_NOFILE,
			"NPROC": UlimitName_NPROC,
			"RSS": UlimitName_RSS,
			"RTPRIO": UlimitName_RTPRIO,
			"RTTIME": UlimitName_RTTIME,
			"SIGPENDING": UlimitName_SIGPENDING,
			"STACK": UlimitName_STACK,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.Volume",
		reflect.TypeOf((*Volume)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_ecs.VolumeFrom",
		reflect.TypeOf((*VolumeFrom)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_ecs.WindowsOptimizedVersion",
		reflect.TypeOf((*WindowsOptimizedVersion)(nil)).Elem(),
		map[string]interface{}{
			"SERVER_2019": WindowsOptimizedVersion_SERVER_2019,
			"SERVER_2016": WindowsOptimizedVersion_SERVER_2016,
		},
	)
}

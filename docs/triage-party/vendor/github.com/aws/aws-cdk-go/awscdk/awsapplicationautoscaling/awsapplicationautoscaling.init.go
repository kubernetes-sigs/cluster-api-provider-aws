package awsapplicationautoscaling

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.AdjustmentTier",
		reflect.TypeOf((*AdjustmentTier)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_applicationautoscaling.AdjustmentType",
		reflect.TypeOf((*AdjustmentType)(nil)).Elem(),
		map[string]interface{}{
			"CHANGE_IN_CAPACITY": AdjustmentType_CHANGE_IN_CAPACITY,
			"PERCENT_CHANGE_IN_CAPACITY": AdjustmentType_PERCENT_CHANGE_IN_CAPACITY,
			"EXACT_CAPACITY": AdjustmentType_EXACT_CAPACITY,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.BaseScalableAttribute",
		reflect.TypeOf((*BaseScalableAttribute)(nil)).Elem(),
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
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_BaseScalableAttribute{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.BaseScalableAttributeProps",
		reflect.TypeOf((*BaseScalableAttributeProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.BaseTargetTrackingProps",
		reflect.TypeOf((*BaseTargetTrackingProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.BasicStepScalingPolicyProps",
		reflect.TypeOf((*BasicStepScalingPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.BasicTargetTrackingScalingPolicyProps",
		reflect.TypeOf((*BasicTargetTrackingScalingPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget",
		reflect.TypeOf((*CfnScalableTarget)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "maxCapacity", GoGetter: "MaxCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "minCapacity", GoGetter: "MinCapacity"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "resourceId", GoGetter: "ResourceId"},
			_jsii_.MemberProperty{JsiiProperty: "roleArn", GoGetter: "RoleArn"},
			_jsii_.MemberProperty{JsiiProperty: "scalableDimension", GoGetter: "ScalableDimension"},
			_jsii_.MemberProperty{JsiiProperty: "scheduledActions", GoGetter: "ScheduledActions"},
			_jsii_.MemberProperty{JsiiProperty: "serviceNamespace", GoGetter: "ServiceNamespace"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "suspendedState", GoGetter: "SuspendedState"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnScalableTarget{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget.ScalableTargetActionProperty",
		reflect.TypeOf((*CfnScalableTarget_ScalableTargetActionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget.ScheduledActionProperty",
		reflect.TypeOf((*CfnScalableTarget_ScheduledActionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget.SuspendedStateProperty",
		reflect.TypeOf((*CfnScalableTarget_SuspendedStateProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalableTargetProps",
		reflect.TypeOf((*CfnScalableTargetProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy",
		reflect.TypeOf((*CfnScalingPolicy)(nil)).Elem(),
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
			_jsii_.MemberProperty{JsiiProperty: "policyName", GoGetter: "PolicyName"},
			_jsii_.MemberProperty{JsiiProperty: "policyType", GoGetter: "PolicyType"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "resourceId", GoGetter: "ResourceId"},
			_jsii_.MemberProperty{JsiiProperty: "scalableDimension", GoGetter: "ScalableDimension"},
			_jsii_.MemberProperty{JsiiProperty: "scalingTargetId", GoGetter: "ScalingTargetId"},
			_jsii_.MemberProperty{JsiiProperty: "serviceNamespace", GoGetter: "ServiceNamespace"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "stepScalingPolicyConfiguration", GoGetter: "StepScalingPolicyConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "targetTrackingScalingPolicyConfiguration", GoGetter: "TargetTrackingScalingPolicyConfiguration"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnScalingPolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy.CustomizedMetricSpecificationProperty",
		reflect.TypeOf((*CfnScalingPolicy_CustomizedMetricSpecificationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy.MetricDimensionProperty",
		reflect.TypeOf((*CfnScalingPolicy_MetricDimensionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy.PredefinedMetricSpecificationProperty",
		reflect.TypeOf((*CfnScalingPolicy_PredefinedMetricSpecificationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy.StepAdjustmentProperty",
		reflect.TypeOf((*CfnScalingPolicy_StepAdjustmentProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy.StepScalingPolicyConfigurationProperty",
		reflect.TypeOf((*CfnScalingPolicy_StepScalingPolicyConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy.TargetTrackingScalingPolicyConfigurationProperty",
		reflect.TypeOf((*CfnScalingPolicy_TargetTrackingScalingPolicyConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicyProps",
		reflect.TypeOf((*CfnScalingPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.CronOptions",
		reflect.TypeOf((*CronOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.EnableScalingProps",
		reflect.TypeOf((*EnableScalingProps)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_applicationautoscaling.IScalableTarget",
		reflect.TypeOf((*IScalableTarget)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "scalableTargetId", GoGetter: "ScalableTargetId"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IScalableTarget{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_applicationautoscaling.MetricAggregationType",
		reflect.TypeOf((*MetricAggregationType)(nil)).Elem(),
		map[string]interface{}{
			"AVERAGE": MetricAggregationType_AVERAGE,
			"MINIMUM": MetricAggregationType_MINIMUM,
			"MAXIMUM": MetricAggregationType_MAXIMUM,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_applicationautoscaling.PredefinedMetric",
		reflect.TypeOf((*PredefinedMetric)(nil)).Elem(),
		map[string]interface{}{
			"DYNAMODB_READ_CAPACITY_UTILIZATION": PredefinedMetric_DYNAMODB_READ_CAPACITY_UTILIZATION,
			"DYANMODB_WRITE_CAPACITY_UTILIZATION": PredefinedMetric_DYANMODB_WRITE_CAPACITY_UTILIZATION,
			"ALB_REQUEST_COUNT_PER_TARGET": PredefinedMetric_ALB_REQUEST_COUNT_PER_TARGET,
			"RDS_READER_AVERAGE_CPU_UTILIZATION": PredefinedMetric_RDS_READER_AVERAGE_CPU_UTILIZATION,
			"RDS_READER_AVERAGE_DATABASE_CONNECTIONS": PredefinedMetric_RDS_READER_AVERAGE_DATABASE_CONNECTIONS,
			"EC2_SPOT_FLEET_REQUEST_AVERAGE_CPU_UTILIZATION": PredefinedMetric_EC2_SPOT_FLEET_REQUEST_AVERAGE_CPU_UTILIZATION,
			"EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_IN": PredefinedMetric_EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_IN,
			"EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_OUT": PredefinedMetric_EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_OUT,
			"SAGEMAKER_VARIANT_INVOCATIONS_PER_INSTANCE": PredefinedMetric_SAGEMAKER_VARIANT_INVOCATIONS_PER_INSTANCE,
			"ECS_SERVICE_AVERAGE_CPU_UTILIZATION": PredefinedMetric_ECS_SERVICE_AVERAGE_CPU_UTILIZATION,
			"ECS_SERVICE_AVERAGE_MEMORY_UTILIZATION": PredefinedMetric_ECS_SERVICE_AVERAGE_MEMORY_UTILIZATION,
			"LAMBDA_PROVISIONED_CONCURRENCY_UTILIZATION": PredefinedMetric_LAMBDA_PROVISIONED_CONCURRENCY_UTILIZATION,
			"KAFKA_BROKER_STORAGE_UTILIZATION": PredefinedMetric_KAFKA_BROKER_STORAGE_UTILIZATION,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.ScalableTarget",
		reflect.TypeOf((*ScalableTarget)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addToRolePolicy", GoMethod: "AddToRolePolicy"},
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
			_jsii_.MemberProperty{JsiiProperty: "role", GoGetter: "Role"},
			_jsii_.MemberProperty{JsiiProperty: "scalableTargetId", GoGetter: "ScalableTargetId"},
			_jsii_.MemberMethod{JsiiMethod: "scaleOnMetric", GoMethod: "ScaleOnMetric"},
			_jsii_.MemberMethod{JsiiMethod: "scaleOnSchedule", GoMethod: "ScaleOnSchedule"},
			_jsii_.MemberMethod{JsiiMethod: "scaleToTrackMetric", GoMethod: "ScaleToTrackMetric"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_ScalableTarget{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IScalableTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.ScalableTargetProps",
		reflect.TypeOf((*ScalableTargetProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.ScalingInterval",
		reflect.TypeOf((*ScalingInterval)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.ScalingSchedule",
		reflect.TypeOf((*ScalingSchedule)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.Schedule",
		reflect.TypeOf((*Schedule)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "expressionString", GoGetter: "ExpressionString"},
		},
		func() interface{} {
			return &jsiiProxy_Schedule{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_applicationautoscaling.ServiceNamespace",
		reflect.TypeOf((*ServiceNamespace)(nil)).Elem(),
		map[string]interface{}{
			"ECS": ServiceNamespace_ECS,
			"ELASTIC_MAP_REDUCE": ServiceNamespace_ELASTIC_MAP_REDUCE,
			"EC2": ServiceNamespace_EC2,
			"APPSTREAM": ServiceNamespace_APPSTREAM,
			"DYNAMODB": ServiceNamespace_DYNAMODB,
			"RDS": ServiceNamespace_RDS,
			"SAGEMAKER": ServiceNamespace_SAGEMAKER,
			"CUSTOM_RESOURCE": ServiceNamespace_CUSTOM_RESOURCE,
			"LAMBDA": ServiceNamespace_LAMBDA,
			"COMPREHEND": ServiceNamespace_COMPREHEND,
			"KAFKA": ServiceNamespace_KAFKA,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.StepScalingAction",
		reflect.TypeOf((*StepScalingAction)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addAdjustment", GoMethod: "AddAdjustment"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "scalingPolicyArn", GoGetter: "ScalingPolicyArn"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_StepScalingAction{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.StepScalingActionProps",
		reflect.TypeOf((*StepScalingActionProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.StepScalingPolicy",
		reflect.TypeOf((*StepScalingPolicy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "lowerAction", GoGetter: "LowerAction"},
			_jsii_.MemberProperty{JsiiProperty: "lowerAlarm", GoGetter: "LowerAlarm"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "upperAction", GoGetter: "UpperAction"},
			_jsii_.MemberProperty{JsiiProperty: "upperAlarm", GoGetter: "UpperAlarm"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_StepScalingPolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.StepScalingPolicyProps",
		reflect.TypeOf((*StepScalingPolicyProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_applicationautoscaling.TargetTrackingScalingPolicy",
		reflect.TypeOf((*TargetTrackingScalingPolicy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "scalingPolicyArn", GoGetter: "ScalingPolicyArn"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_TargetTrackingScalingPolicy{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_applicationautoscaling.TargetTrackingScalingPolicyProps",
		reflect.TypeOf((*TargetTrackingScalingPolicyProps)(nil)).Elem(),
	)
}

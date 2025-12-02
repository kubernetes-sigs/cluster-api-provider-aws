package awsapplicationautoscaling

import (
	"time"

	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapplicationautoscaling/internal"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
)

// An adjustment.
// Experimental.
type AdjustmentTier struct {
	// What number to adjust the capacity with.
	//
	// The number is interpeted as an added capacity, a new fixed capacity or an
	// added percentage depending on the AdjustmentType value of the
	// StepScalingPolicy.
	//
	// Can be positive or negative.
	// Experimental.
	Adjustment *float64 `json:"adjustment"`
	// Lower bound where this scaling tier applies.
	//
	// The scaling tier applies if the difference between the metric
	// value and its alarm threshold is higher than this value.
	// Experimental.
	LowerBound *float64 `json:"lowerBound"`
	// Upper bound where this scaling tier applies.
	//
	// The scaling tier applies if the difference between the metric
	// value and its alarm threshold is lower than this value.
	// Experimental.
	UpperBound *float64 `json:"upperBound"`
}

// How adjustment numbers are interpreted.
// Experimental.
type AdjustmentType string

const (
	AdjustmentType_CHANGE_IN_CAPACITY AdjustmentType = "CHANGE_IN_CAPACITY"
	AdjustmentType_PERCENT_CHANGE_IN_CAPACITY AdjustmentType = "PERCENT_CHANGE_IN_CAPACITY"
	AdjustmentType_EXACT_CAPACITY AdjustmentType = "EXACT_CAPACITY"
)

// Represent an attribute for which autoscaling can be configured.
//
// This class is basically a light wrapper around ScalableTarget, but with
// all methods protected instead of public so they can be selectively
// exposed and/or more specific versions of them can be exposed by derived
// classes for individual services support autoscaling.
//
// Typical use cases:
//
// - Hide away the PredefinedMetric enum for target tracking policies.
// - Don't expose all scaling methods (for example Dynamo tables don't support
//    Step Scaling, so the Dynamo subclass won't expose this method).
// Experimental.
type BaseScalableAttribute interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	Props() *BaseScalableAttributeProps
	DoScaleOnMetric(id *string, props *BasicStepScalingPolicyProps)
	DoScaleOnSchedule(id *string, props *ScalingSchedule)
	DoScaleToTrackMetric(id *string, props *BasicTargetTrackingScalingPolicyProps)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for BaseScalableAttribute
type jsiiProxy_BaseScalableAttribute struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_BaseScalableAttribute) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseScalableAttribute) Props() *BaseScalableAttributeProps {
	var returns *BaseScalableAttributeProps
	_jsii_.Get(
		j,
		"props",
		&returns,
	)
	return returns
}


// Experimental.
func NewBaseScalableAttribute_Override(b BaseScalableAttribute, scope constructs.Construct, id *string, props *BaseScalableAttributeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.BaseScalableAttribute",
		[]interface{}{scope, id, props},
		b,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BaseScalableAttribute_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.BaseScalableAttribute",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Scale out or in based on a metric value.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) DoScaleOnMetric(id *string, props *BasicStepScalingPolicyProps) {
	_jsii_.InvokeVoid(
		b,
		"doScaleOnMetric",
		[]interface{}{id, props},
	)
}

// Scale out or in based on time.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) DoScaleOnSchedule(id *string, props *ScalingSchedule) {
	_jsii_.InvokeVoid(
		b,
		"doScaleOnSchedule",
		[]interface{}{id, props},
	)
}

// Scale out or in in order to keep a metric around a target value.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) DoScaleToTrackMetric(id *string, props *BasicTargetTrackingScalingPolicyProps) {
	_jsii_.InvokeVoid(
		b,
		"doScaleToTrackMetric",
		[]interface{}{id, props},
	)
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) OnPrepare() {
	_jsii_.InvokeVoid(
		b,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"onSynthesize",
		[]interface{}{session},
	)
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) Prepare() {
	_jsii_.InvokeVoid(
		b,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (b *jsiiProxy_BaseScalableAttribute) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a ScalableTableAttribute.
// Experimental.
type BaseScalableAttributeProps struct {
	// Maximum capacity to scale to.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// Minimum capacity to scale to.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// Scalable dimension of the attribute.
	// Experimental.
	Dimension *string `json:"dimension"`
	// Resource ID of the attribute.
	// Experimental.
	ResourceId *string `json:"resourceId"`
	// Role to use for scaling.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// Service namespace of the scalable attribute.
	// Experimental.
	ServiceNamespace ServiceNamespace `json:"serviceNamespace"`
}

// Base interface for target tracking props.
//
// Contains the attributes that are common to target tracking policies,
// except the ones relating to the metric and to the scalable target.
//
// This interface is reused by more specific target tracking props objects
// in other services.
// Experimental.
type BaseTargetTrackingProps struct {
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the scalable resource. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// scalable resource.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Period after a scale in activity completes before another scale in activity can start.
	// Experimental.
	ScaleInCooldown awscdk.Duration `json:"scaleInCooldown"`
	// Period after a scale out activity completes before another scale out activity can start.
	// Experimental.
	ScaleOutCooldown awscdk.Duration `json:"scaleOutCooldown"`
}

// Experimental.
type BasicStepScalingPolicyProps struct {
	// Metric to scale on.
	// Experimental.
	Metric awscloudwatch.IMetric `json:"metric"`
	// The intervals for scaling.
	//
	// Maps a range of metric values to a particular scaling behavior.
	// Experimental.
	ScalingSteps *[]*ScalingInterval `json:"scalingSteps"`
	// How the adjustment numbers inside 'intervals' are interpreted.
	// Experimental.
	AdjustmentType AdjustmentType `json:"adjustmentType"`
	// Grace period after scaling activity.
	//
	// Subsequent scale outs during the cooldown period are squashed so that only
	// the biggest scale out happens.
	//
	// Subsequent scale ins during the cooldown period are ignored.
	// See: https://docs.aws.amazon.com/autoscaling/application/APIReference/API_StepScalingPolicyConfiguration.html
	//
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// How many evaluation periods of the metric to wait before triggering a scaling action.
	//
	// Raising this value can be used to smooth out the metric, at the expense
	// of slower response times.
	// Experimental.
	EvaluationPeriods *float64 `json:"evaluationPeriods"`
	// Aggregation to apply to all data points over the evaluation periods.
	//
	// Only has meaning if `evaluationPeriods != 1`.
	// Experimental.
	MetricAggregationType MetricAggregationType `json:"metricAggregationType"`
	// Minimum absolute number to adjust capacity with as result of percentage scaling.
	//
	// Only when using AdjustmentType = PercentChangeInCapacity, this number controls
	// the minimum absolute effect size.
	// Experimental.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
}

// Properties for a Target Tracking policy that include the metric but exclude the target.
// Experimental.
type BasicTargetTrackingScalingPolicyProps struct {
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the scalable resource. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// scalable resource.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Period after a scale in activity completes before another scale in activity can start.
	// Experimental.
	ScaleInCooldown awscdk.Duration `json:"scaleInCooldown"`
	// Period after a scale out activity completes before another scale out activity can start.
	// Experimental.
	ScaleOutCooldown awscdk.Duration `json:"scaleOutCooldown"`
	// The target value for the metric.
	// Experimental.
	TargetValue *float64 `json:"targetValue"`
	// A custom metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	CustomMetric awscloudwatch.IMetric `json:"customMetric"`
	// A predefined metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	PredefinedMetric PredefinedMetric `json:"predefinedMetric"`
	// Identify the resource associated with the metric type.
	//
	// Only used for predefined metric ALBRequestCountPerTarget.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	ResourceLabel *string `json:"resourceLabel"`
}

// A CloudFormation `AWS::ApplicationAutoScaling::ScalableTarget`.
type CfnScalableTarget interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	MaxCapacity() *float64
	SetMaxCapacity(val *float64)
	MinCapacity() *float64
	SetMinCapacity(val *float64)
	Node() awscdk.ConstructNode
	Ref() *string
	ResourceId() *string
	SetResourceId(val *string)
	RoleArn() *string
	SetRoleArn(val *string)
	ScalableDimension() *string
	SetScalableDimension(val *string)
	ScheduledActions() interface{}
	SetScheduledActions(val interface{})
	ServiceNamespace() *string
	SetServiceNamespace(val *string)
	Stack() awscdk.Stack
	SuspendedState() interface{}
	SetSuspendedState(val interface{})
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target awscdk.CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions)
	GetAtt(attributeName *string) awscdk.Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector awscdk.TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnScalableTarget
type jsiiProxy_CfnScalableTarget struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnScalableTarget) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) MaxCapacity() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxCapacity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) MinCapacity() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"minCapacity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) ResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) ScalableDimension() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalableDimension",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) ScheduledActions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"scheduledActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) ServiceNamespace() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceNamespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) SuspendedState() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"suspendedState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalableTarget) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApplicationAutoScaling::ScalableTarget`.
func NewCfnScalableTarget(scope awscdk.Construct, id *string, props *CfnScalableTargetProps) CfnScalableTarget {
	_init_.Initialize()

	j := jsiiProxy_CfnScalableTarget{}

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApplicationAutoScaling::ScalableTarget`.
func NewCfnScalableTarget_Override(c CfnScalableTarget, scope awscdk.Construct, id *string, props *CfnScalableTargetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetMaxCapacity(val *float64) {
	_jsii_.Set(
		j,
		"maxCapacity",
		val,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetMinCapacity(val *float64) {
	_jsii_.Set(
		j,
		"minCapacity",
		val,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetResourceId(val *string) {
	_jsii_.Set(
		j,
		"resourceId",
		val,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetRoleArn(val *string) {
	_jsii_.Set(
		j,
		"roleArn",
		val,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetScalableDimension(val *string) {
	_jsii_.Set(
		j,
		"scalableDimension",
		val,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetScheduledActions(val interface{}) {
	_jsii_.Set(
		j,
		"scheduledActions",
		val,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetServiceNamespace(val *string) {
	_jsii_.Set(
		j,
		"serviceNamespace",
		val,
	)
}

func (j *jsiiProxy_CfnScalableTarget) SetSuspendedState(val interface{}) {
	_jsii_.Set(
		j,
		"suspendedState",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnScalableTarget_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnScalableTarget_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnScalableTarget_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnScalableTarget_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_applicationautoscaling.CfnScalableTarget",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) AddDeletionOverride(path *string) {
	_jsii_.InvokeVoid(
		c,
		"addDeletionOverride",
		[]interface{}{path},
	)
}

// Indicates that this resource depends on another resource and cannot be provisioned unless the other resource has been successfully provisioned.
//
// This can be used for resources across stacks (or nested stack) boundaries
// and the dependency will automatically be transferred to the relevant scope.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) AddDependsOn(target awscdk.CfnResource) {
	_jsii_.InvokeVoid(
		c,
		"addDependsOn",
		[]interface{}{target},
	)
}

// Add a value to the CloudFormation Resource Metadata.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/metadata-section-structure.html
//
// Note that this is a different set of metadata from CDK node metadata; this
// metadata ends up in the stack template under the resource, whereas CDK
// node metadata ends up in the Cloud Assembly.
//
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) AddMetadata(key *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addMetadata",
		[]interface{}{key, value},
	)
}

// Adds an override to the synthesized CloudFormation resource.
//
// To add a
// property override, either use `addPropertyOverride` or prefix `path` with
// "Properties." (i.e. `Properties.TopicName`).
//
// If the override is nested, separate each nested level using a dot (.) in the path parameter.
// If there is an array as part of the nesting, specify the index in the path.
//
// To include a literal `.` in the property name, prefix with a `\`. In most
// programming languages you will need to write this as `"\\."` because the
// `\` itself will need to be escaped.
//
// For example,
// ```typescript
// cfnResource.addOverride('Properties.GlobalSecondaryIndexes.0.Projection.NonKeyAttributes', ['myattribute']);
// cfnResource.addOverride('Properties.GlobalSecondaryIndexes.1.ProjectionType', 'INCLUDE');
// ```
// would add the overrides
// ```json
// "Properties": {
//    "GlobalSecondaryIndexes": [
//      {
//        "Projection": {
//          "NonKeyAttributes": [ "myattribute" ]
//          ...
//        }
//        ...
//      },
//      {
//        "ProjectionType": "INCLUDE"
//        ...
//      },
//    ]
//    ...
// }
// ```
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
	_jsii_.InvokeVoid(
		c,
		"applyRemovalPolicy",
		[]interface{}{policy, options},
	)
}

// Returns a token for an runtime attribute of this resource.
//
// Ideally, use generated attribute accessors (e.g. `resource.arn`), but this can be used for future compatibility
// in case there is no generated attribute.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) GetAtt(attributeName *string) awscdk.Reference {
	var returns awscdk.Reference

	_jsii_.Invoke(
		c,
		"getAtt",
		[]interface{}{attributeName},
		&returns,
	)

	return returns
}

// Retrieve a value value from the CloudFormation Resource Metadata.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/metadata-section-structure.html
//
// Note that this is a different set of metadata from CDK node metadata; this
// metadata ends up in the stack template under the resource, whereas CDK
// node metadata ends up in the Cloud Assembly.
//
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) GetMetadata(key *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"getMetadata",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Examines the CloudFormation resource and discloses attributes.
func (c *jsiiProxy_CfnScalableTarget) Inspect(inspector awscdk.TreeInspector) {
	_jsii_.InvokeVoid(
		c,
		"inspect",
		[]interface{}{inspector},
	)
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) OnPrepare() {
	_jsii_.InvokeVoid(
		c,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"onSynthesize",
		[]interface{}{session},
	)
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Overrides the auto-generated logical ID with a specific ID.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) OverrideLogicalId(newLogicalId *string) {
	_jsii_.InvokeVoid(
		c,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnScalableTarget) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"renderProperties",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Can be overridden by subclasses to determine if this resource will be rendered into the cloudformation template.
//
// Returns: `true` if the resource should be included or `false` is the resource
// should be omitted.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) ShouldSynthesize() *bool {
	var returns *bool

	_jsii_.Invoke(
		c,
		"shouldSynthesize",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (c *jsiiProxy_CfnScalableTarget) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (c *jsiiProxy_CfnScalableTarget) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnScalableTarget_ScalableTargetActionProperty struct {
	// `CfnScalableTarget.ScalableTargetActionProperty.MaxCapacity`.
	MaxCapacity *float64 `json:"maxCapacity"`
	// `CfnScalableTarget.ScalableTargetActionProperty.MinCapacity`.
	MinCapacity *float64 `json:"minCapacity"`
}

type CfnScalableTarget_ScheduledActionProperty struct {
	// `CfnScalableTarget.ScheduledActionProperty.Schedule`.
	Schedule *string `json:"schedule"`
	// `CfnScalableTarget.ScheduledActionProperty.ScheduledActionName`.
	ScheduledActionName *string `json:"scheduledActionName"`
	// `CfnScalableTarget.ScheduledActionProperty.EndTime`.
	EndTime interface{} `json:"endTime"`
	// `CfnScalableTarget.ScheduledActionProperty.ScalableTargetAction`.
	ScalableTargetAction interface{} `json:"scalableTargetAction"`
	// `CfnScalableTarget.ScheduledActionProperty.StartTime`.
	StartTime interface{} `json:"startTime"`
}

type CfnScalableTarget_SuspendedStateProperty struct {
	// `CfnScalableTarget.SuspendedStateProperty.DynamicScalingInSuspended`.
	DynamicScalingInSuspended interface{} `json:"dynamicScalingInSuspended"`
	// `CfnScalableTarget.SuspendedStateProperty.DynamicScalingOutSuspended`.
	DynamicScalingOutSuspended interface{} `json:"dynamicScalingOutSuspended"`
	// `CfnScalableTarget.SuspendedStateProperty.ScheduledScalingSuspended`.
	ScheduledScalingSuspended interface{} `json:"scheduledScalingSuspended"`
}

// Properties for defining a `AWS::ApplicationAutoScaling::ScalableTarget`.
type CfnScalableTargetProps struct {
	// `AWS::ApplicationAutoScaling::ScalableTarget.MaxCapacity`.
	MaxCapacity *float64 `json:"maxCapacity"`
	// `AWS::ApplicationAutoScaling::ScalableTarget.MinCapacity`.
	MinCapacity *float64 `json:"minCapacity"`
	// `AWS::ApplicationAutoScaling::ScalableTarget.ResourceId`.
	ResourceId *string `json:"resourceId"`
	// `AWS::ApplicationAutoScaling::ScalableTarget.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `AWS::ApplicationAutoScaling::ScalableTarget.ScalableDimension`.
	ScalableDimension *string `json:"scalableDimension"`
	// `AWS::ApplicationAutoScaling::ScalableTarget.ServiceNamespace`.
	ServiceNamespace *string `json:"serviceNamespace"`
	// `AWS::ApplicationAutoScaling::ScalableTarget.ScheduledActions`.
	ScheduledActions interface{} `json:"scheduledActions"`
	// `AWS::ApplicationAutoScaling::ScalableTarget.SuspendedState`.
	SuspendedState interface{} `json:"suspendedState"`
}

// A CloudFormation `AWS::ApplicationAutoScaling::ScalingPolicy`.
type CfnScalingPolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	PolicyName() *string
	SetPolicyName(val *string)
	PolicyType() *string
	SetPolicyType(val *string)
	Ref() *string
	ResourceId() *string
	SetResourceId(val *string)
	ScalableDimension() *string
	SetScalableDimension(val *string)
	ScalingTargetId() *string
	SetScalingTargetId(val *string)
	ServiceNamespace() *string
	SetServiceNamespace(val *string)
	Stack() awscdk.Stack
	StepScalingPolicyConfiguration() interface{}
	SetStepScalingPolicyConfiguration(val interface{})
	TargetTrackingScalingPolicyConfiguration() interface{}
	SetTargetTrackingScalingPolicyConfiguration(val interface{})
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target awscdk.CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions)
	GetAtt(attributeName *string) awscdk.Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector awscdk.TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnScalingPolicy
type jsiiProxy_CfnScalingPolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnScalingPolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) PolicyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"policyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) PolicyType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"policyType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) ResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) ScalableDimension() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalableDimension",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) ScalingTargetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalingTargetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) ServiceNamespace() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceNamespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) StepScalingPolicyConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stepScalingPolicyConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) TargetTrackingScalingPolicyConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"targetTrackingScalingPolicyConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApplicationAutoScaling::ScalingPolicy`.
func NewCfnScalingPolicy(scope awscdk.Construct, id *string, props *CfnScalingPolicyProps) CfnScalingPolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnScalingPolicy{}

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApplicationAutoScaling::ScalingPolicy`.
func NewCfnScalingPolicy_Override(c CfnScalingPolicy, scope awscdk.Construct, id *string, props *CfnScalingPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetPolicyName(val *string) {
	_jsii_.Set(
		j,
		"policyName",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetPolicyType(val *string) {
	_jsii_.Set(
		j,
		"policyType",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetResourceId(val *string) {
	_jsii_.Set(
		j,
		"resourceId",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetScalableDimension(val *string) {
	_jsii_.Set(
		j,
		"scalableDimension",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetScalingTargetId(val *string) {
	_jsii_.Set(
		j,
		"scalingTargetId",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetServiceNamespace(val *string) {
	_jsii_.Set(
		j,
		"serviceNamespace",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetStepScalingPolicyConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"stepScalingPolicyConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetTargetTrackingScalingPolicyConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"targetTrackingScalingPolicyConfiguration",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnScalingPolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnScalingPolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnScalingPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnScalingPolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_applicationautoscaling.CfnScalingPolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddDeletionOverride(path *string) {
	_jsii_.InvokeVoid(
		c,
		"addDeletionOverride",
		[]interface{}{path},
	)
}

// Indicates that this resource depends on another resource and cannot be provisioned unless the other resource has been successfully provisioned.
//
// This can be used for resources across stacks (or nested stack) boundaries
// and the dependency will automatically be transferred to the relevant scope.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddDependsOn(target awscdk.CfnResource) {
	_jsii_.InvokeVoid(
		c,
		"addDependsOn",
		[]interface{}{target},
	)
}

// Add a value to the CloudFormation Resource Metadata.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/metadata-section-structure.html
//
// Note that this is a different set of metadata from CDK node metadata; this
// metadata ends up in the stack template under the resource, whereas CDK
// node metadata ends up in the Cloud Assembly.
//
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddMetadata(key *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addMetadata",
		[]interface{}{key, value},
	)
}

// Adds an override to the synthesized CloudFormation resource.
//
// To add a
// property override, either use `addPropertyOverride` or prefix `path` with
// "Properties." (i.e. `Properties.TopicName`).
//
// If the override is nested, separate each nested level using a dot (.) in the path parameter.
// If there is an array as part of the nesting, specify the index in the path.
//
// To include a literal `.` in the property name, prefix with a `\`. In most
// programming languages you will need to write this as `"\\."` because the
// `\` itself will need to be escaped.
//
// For example,
// ```typescript
// cfnResource.addOverride('Properties.GlobalSecondaryIndexes.0.Projection.NonKeyAttributes', ['myattribute']);
// cfnResource.addOverride('Properties.GlobalSecondaryIndexes.1.ProjectionType', 'INCLUDE');
// ```
// would add the overrides
// ```json
// "Properties": {
//    "GlobalSecondaryIndexes": [
//      {
//        "Projection": {
//          "NonKeyAttributes": [ "myattribute" ]
//          ...
//        }
//        ...
//      },
//      {
//        "ProjectionType": "INCLUDE"
//        ...
//      },
//    ]
//    ...
// }
// ```
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
	_jsii_.InvokeVoid(
		c,
		"applyRemovalPolicy",
		[]interface{}{policy, options},
	)
}

// Returns a token for an runtime attribute of this resource.
//
// Ideally, use generated attribute accessors (e.g. `resource.arn`), but this can be used for future compatibility
// in case there is no generated attribute.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) GetAtt(attributeName *string) awscdk.Reference {
	var returns awscdk.Reference

	_jsii_.Invoke(
		c,
		"getAtt",
		[]interface{}{attributeName},
		&returns,
	)

	return returns
}

// Retrieve a value value from the CloudFormation Resource Metadata.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/metadata-section-structure.html
//
// Note that this is a different set of metadata from CDK node metadata; this
// metadata ends up in the stack template under the resource, whereas CDK
// node metadata ends up in the Cloud Assembly.
//
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) GetMetadata(key *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"getMetadata",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Examines the CloudFormation resource and discloses attributes.
func (c *jsiiProxy_CfnScalingPolicy) Inspect(inspector awscdk.TreeInspector) {
	_jsii_.InvokeVoid(
		c,
		"inspect",
		[]interface{}{inspector},
	)
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) OnPrepare() {
	_jsii_.InvokeVoid(
		c,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"onSynthesize",
		[]interface{}{session},
	)
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Overrides the auto-generated logical ID with a specific ID.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) OverrideLogicalId(newLogicalId *string) {
	_jsii_.InvokeVoid(
		c,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnScalingPolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"renderProperties",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Can be overridden by subclasses to determine if this resource will be rendered into the cloudformation template.
//
// Returns: `true` if the resource should be included or `false` is the resource
// should be omitted.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) ShouldSynthesize() *bool {
	var returns *bool

	_jsii_.Invoke(
		c,
		"shouldSynthesize",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnScalingPolicy_CustomizedMetricSpecificationProperty struct {
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.MetricName`.
	MetricName *string `json:"metricName"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Namespace`.
	Namespace *string `json:"namespace"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Statistic`.
	Statistic *string `json:"statistic"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Dimensions`.
	Dimensions interface{} `json:"dimensions"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Unit`.
	Unit *string `json:"unit"`
}

type CfnScalingPolicy_MetricDimensionProperty struct {
	// `CfnScalingPolicy.MetricDimensionProperty.Name`.
	Name *string `json:"name"`
	// `CfnScalingPolicy.MetricDimensionProperty.Value`.
	Value *string `json:"value"`
}

type CfnScalingPolicy_PredefinedMetricSpecificationProperty struct {
	// `CfnScalingPolicy.PredefinedMetricSpecificationProperty.PredefinedMetricType`.
	PredefinedMetricType *string `json:"predefinedMetricType"`
	// `CfnScalingPolicy.PredefinedMetricSpecificationProperty.ResourceLabel`.
	ResourceLabel *string `json:"resourceLabel"`
}

type CfnScalingPolicy_StepAdjustmentProperty struct {
	// `CfnScalingPolicy.StepAdjustmentProperty.ScalingAdjustment`.
	ScalingAdjustment *float64 `json:"scalingAdjustment"`
	// `CfnScalingPolicy.StepAdjustmentProperty.MetricIntervalLowerBound`.
	MetricIntervalLowerBound *float64 `json:"metricIntervalLowerBound"`
	// `CfnScalingPolicy.StepAdjustmentProperty.MetricIntervalUpperBound`.
	MetricIntervalUpperBound *float64 `json:"metricIntervalUpperBound"`
}

type CfnScalingPolicy_StepScalingPolicyConfigurationProperty struct {
	// `CfnScalingPolicy.StepScalingPolicyConfigurationProperty.AdjustmentType`.
	AdjustmentType *string `json:"adjustmentType"`
	// `CfnScalingPolicy.StepScalingPolicyConfigurationProperty.Cooldown`.
	Cooldown *float64 `json:"cooldown"`
	// `CfnScalingPolicy.StepScalingPolicyConfigurationProperty.MetricAggregationType`.
	MetricAggregationType *string `json:"metricAggregationType"`
	// `CfnScalingPolicy.StepScalingPolicyConfigurationProperty.MinAdjustmentMagnitude`.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
	// `CfnScalingPolicy.StepScalingPolicyConfigurationProperty.StepAdjustments`.
	StepAdjustments interface{} `json:"stepAdjustments"`
}

type CfnScalingPolicy_TargetTrackingScalingPolicyConfigurationProperty struct {
	// `CfnScalingPolicy.TargetTrackingScalingPolicyConfigurationProperty.TargetValue`.
	TargetValue *float64 `json:"targetValue"`
	// `CfnScalingPolicy.TargetTrackingScalingPolicyConfigurationProperty.CustomizedMetricSpecification`.
	CustomizedMetricSpecification interface{} `json:"customizedMetricSpecification"`
	// `CfnScalingPolicy.TargetTrackingScalingPolicyConfigurationProperty.DisableScaleIn`.
	DisableScaleIn interface{} `json:"disableScaleIn"`
	// `CfnScalingPolicy.TargetTrackingScalingPolicyConfigurationProperty.PredefinedMetricSpecification`.
	PredefinedMetricSpecification interface{} `json:"predefinedMetricSpecification"`
	// `CfnScalingPolicy.TargetTrackingScalingPolicyConfigurationProperty.ScaleInCooldown`.
	ScaleInCooldown *float64 `json:"scaleInCooldown"`
	// `CfnScalingPolicy.TargetTrackingScalingPolicyConfigurationProperty.ScaleOutCooldown`.
	ScaleOutCooldown *float64 `json:"scaleOutCooldown"`
}

// Properties for defining a `AWS::ApplicationAutoScaling::ScalingPolicy`.
type CfnScalingPolicyProps struct {
	// `AWS::ApplicationAutoScaling::ScalingPolicy.PolicyName`.
	PolicyName *string `json:"policyName"`
	// `AWS::ApplicationAutoScaling::ScalingPolicy.PolicyType`.
	PolicyType *string `json:"policyType"`
	// `AWS::ApplicationAutoScaling::ScalingPolicy.ResourceId`.
	ResourceId *string `json:"resourceId"`
	// `AWS::ApplicationAutoScaling::ScalingPolicy.ScalableDimension`.
	ScalableDimension *string `json:"scalableDimension"`
	// `AWS::ApplicationAutoScaling::ScalingPolicy.ScalingTargetId`.
	ScalingTargetId *string `json:"scalingTargetId"`
	// `AWS::ApplicationAutoScaling::ScalingPolicy.ServiceNamespace`.
	ServiceNamespace *string `json:"serviceNamespace"`
	// `AWS::ApplicationAutoScaling::ScalingPolicy.StepScalingPolicyConfiguration`.
	StepScalingPolicyConfiguration interface{} `json:"stepScalingPolicyConfiguration"`
	// `AWS::ApplicationAutoScaling::ScalingPolicy.TargetTrackingScalingPolicyConfiguration`.
	TargetTrackingScalingPolicyConfiguration interface{} `json:"targetTrackingScalingPolicyConfiguration"`
}

// Options to configure a cron expression.
//
// All fields are strings so you can use complex expressions. Absence of
// a field implies '*' or '?', whichever one is appropriate.
// See: https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html#CronExpressions
//
// Experimental.
type CronOptions struct {
	// The day of the month to run this rule at.
	// Experimental.
	Day *string `json:"day"`
	// The hour to run this rule at.
	// Experimental.
	Hour *string `json:"hour"`
	// The minute to run this rule at.
	// Experimental.
	Minute *string `json:"minute"`
	// The month to run this rule at.
	// Experimental.
	Month *string `json:"month"`
	// The day of the week to run this rule at.
	// Experimental.
	WeekDay *string `json:"weekDay"`
	// The year to run this rule at.
	// Experimental.
	Year *string `json:"year"`
}

// Properties for enabling DynamoDB capacity scaling.
// Experimental.
type EnableScalingProps struct {
	// Maximum capacity to scale to.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// Minimum capacity to scale to.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
}

// Experimental.
type IScalableTarget interface {
	awscdk.IResource
	// Experimental.
	ScalableTargetId() *string
}

// The jsii proxy for IScalableTarget
type jsiiProxy_IScalableTarget struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IScalableTarget) ScalableTargetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalableTargetId",
		&returns,
	)
	return returns
}

// How the scaling metric is going to be aggregated.
// Experimental.
type MetricAggregationType string

const (
	MetricAggregationType_AVERAGE MetricAggregationType = "AVERAGE"
	MetricAggregationType_MINIMUM MetricAggregationType = "MINIMUM"
	MetricAggregationType_MAXIMUM MetricAggregationType = "MAXIMUM"
)

// One of the predefined autoscaling metrics.
// Experimental.
type PredefinedMetric string

const (
	PredefinedMetric_DYNAMODB_READ_CAPACITY_UTILIZATION PredefinedMetric = "DYNAMODB_READ_CAPACITY_UTILIZATION"
	PredefinedMetric_DYANMODB_WRITE_CAPACITY_UTILIZATION PredefinedMetric = "DYANMODB_WRITE_CAPACITY_UTILIZATION"
	PredefinedMetric_ALB_REQUEST_COUNT_PER_TARGET PredefinedMetric = "ALB_REQUEST_COUNT_PER_TARGET"
	PredefinedMetric_RDS_READER_AVERAGE_CPU_UTILIZATION PredefinedMetric = "RDS_READER_AVERAGE_CPU_UTILIZATION"
	PredefinedMetric_RDS_READER_AVERAGE_DATABASE_CONNECTIONS PredefinedMetric = "RDS_READER_AVERAGE_DATABASE_CONNECTIONS"
	PredefinedMetric_EC2_SPOT_FLEET_REQUEST_AVERAGE_CPU_UTILIZATION PredefinedMetric = "EC2_SPOT_FLEET_REQUEST_AVERAGE_CPU_UTILIZATION"
	PredefinedMetric_EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_IN PredefinedMetric = "EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_IN"
	PredefinedMetric_EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_OUT PredefinedMetric = "EC2_SPOT_FLEET_REQUEST_AVERAGE_NETWORK_OUT"
	PredefinedMetric_SAGEMAKER_VARIANT_INVOCATIONS_PER_INSTANCE PredefinedMetric = "SAGEMAKER_VARIANT_INVOCATIONS_PER_INSTANCE"
	PredefinedMetric_ECS_SERVICE_AVERAGE_CPU_UTILIZATION PredefinedMetric = "ECS_SERVICE_AVERAGE_CPU_UTILIZATION"
	PredefinedMetric_ECS_SERVICE_AVERAGE_MEMORY_UTILIZATION PredefinedMetric = "ECS_SERVICE_AVERAGE_MEMORY_UTILIZATION"
	PredefinedMetric_LAMBDA_PROVISIONED_CONCURRENCY_UTILIZATION PredefinedMetric = "LAMBDA_PROVISIONED_CONCURRENCY_UTILIZATION"
	PredefinedMetric_KAFKA_BROKER_STORAGE_UTILIZATION PredefinedMetric = "KAFKA_BROKER_STORAGE_UTILIZATION"
)

// Define a scalable target.
// Experimental.
type ScalableTarget interface {
	awscdk.Resource
	IScalableTarget
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Role() awsiam.IRole
	ScalableTargetId() *string
	Stack() awscdk.Stack
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	ScaleOnMetric(id *string, props *BasicStepScalingPolicyProps) StepScalingPolicy
	ScaleOnSchedule(id *string, action *ScalingSchedule)
	ScaleToTrackMetric(id *string, props *BasicTargetTrackingScalingPolicyProps) TargetTrackingScalingPolicy
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ScalableTarget
type jsiiProxy_ScalableTarget struct {
	internal.Type__awscdkResource
	jsiiProxy_IScalableTarget
}

func (j *jsiiProxy_ScalableTarget) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScalableTarget) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScalableTarget) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScalableTarget) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScalableTarget) ScalableTargetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalableTargetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScalableTarget) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewScalableTarget(scope constructs.Construct, id *string, props *ScalableTargetProps) ScalableTarget {
	_init_.Initialize()

	j := jsiiProxy_ScalableTarget{}

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.ScalableTarget",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewScalableTarget_Override(s ScalableTarget, scope constructs.Construct, id *string, props *ScalableTargetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.ScalableTarget",
		[]interface{}{scope, id, props},
		s,
	)
}

// Experimental.
func ScalableTarget_FromScalableTargetId(scope constructs.Construct, id *string, scalableTargetId *string) IScalableTarget {
	_init_.Initialize()

	var returns IScalableTarget

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.ScalableTarget",
		"fromScalableTargetId",
		[]interface{}{scope, id, scalableTargetId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ScalableTarget_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.ScalableTarget",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ScalableTarget_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.ScalableTarget",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a policy statement to the role's policy.
// Experimental.
func (s *jsiiProxy_ScalableTarget) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		s,
		"addToRolePolicy",
		[]interface{}{statement},
	)
}

// Apply the given removal policy to this resource.
//
// The Removal Policy controls what happens to this resource when it stops
// being managed by CloudFormation, either because you've removed it from the
// CDK application or because you've made a change that requires the resource
// to be replaced.
//
// The resource can be deleted (`RemovalPolicy.DELETE`), or left in your AWS
// account for data recovery and cleanup later (`RemovalPolicy.RETAIN`).
// Experimental.
func (s *jsiiProxy_ScalableTarget) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (s *jsiiProxy_ScalableTarget) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns an environment-sensitive token that should be used for the resource's "ARN" attribute (e.g. `bucket.bucketArn`).
//
// Normally, this token will resolve to `arnAttr`, but if the resource is
// referenced across environments, `arnComponents` will be used to synthesize
// a concrete ARN with the resource's physical name. Make sure to reference
// `this.physicalName` in `arnComponents`.
// Experimental.
func (s *jsiiProxy_ScalableTarget) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceArnAttribute",
		[]interface{}{arnAttr, arnComponents},
		&returns,
	)

	return returns
}

// Returns an environment-sensitive token that should be used for the resource's "name" attribute (e.g. `bucket.bucketName`).
//
// Normally, this token will resolve to `nameAttr`, but if the resource is
// referenced across environments, it will be resolved to `this.physicalName`,
// which will be a concrete name.
// Experimental.
func (s *jsiiProxy_ScalableTarget) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (s *jsiiProxy_ScalableTarget) OnPrepare() {
	_jsii_.InvokeVoid(
		s,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_ScalableTarget) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"onSynthesize",
		[]interface{}{session},
	)
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (s *jsiiProxy_ScalableTarget) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (s *jsiiProxy_ScalableTarget) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Scale out or in, in response to a metric.
// Experimental.
func (s *jsiiProxy_ScalableTarget) ScaleOnMetric(id *string, props *BasicStepScalingPolicyProps) StepScalingPolicy {
	var returns StepScalingPolicy

	_jsii_.Invoke(
		s,
		"scaleOnMetric",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Scale out or in based on time.
// Experimental.
func (s *jsiiProxy_ScalableTarget) ScaleOnSchedule(id *string, action *ScalingSchedule) {
	_jsii_.InvokeVoid(
		s,
		"scaleOnSchedule",
		[]interface{}{id, action},
	)
}

// Scale out or in in order to keep a metric around a target value.
// Experimental.
func (s *jsiiProxy_ScalableTarget) ScaleToTrackMetric(id *string, props *BasicTargetTrackingScalingPolicyProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		s,
		"scaleToTrackMetric",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_ScalableTarget) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_ScalableTarget) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (s *jsiiProxy_ScalableTarget) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a scalable target.
// Experimental.
type ScalableTargetProps struct {
	// The maximum value that Application Auto Scaling can use to scale a target during a scaling activity.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// The minimum value that Application Auto Scaling can use to scale a target during a scaling activity.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// The resource identifier to associate with this scalable target.
	//
	// This string consists of the resource type and unique identifier.
	//
	// TODO: EXAMPLE
	//
	// See: https://docs.aws.amazon.com/autoscaling/application/APIReference/API_RegisterScalableTarget.html
	//
	// Experimental.
	ResourceId *string `json:"resourceId"`
	// The scalable dimension that's associated with the scalable target.
	//
	// Specify the service namespace, resource type, and scaling property.
	//
	// TODO: EXAMPLE
	//
	// See: https://docs.aws.amazon.com/autoscaling/application/APIReference/API_ScalingPolicy.html
	//
	// Experimental.
	ScalableDimension *string `json:"scalableDimension"`
	// The namespace of the AWS service that provides the resource or custom-resource for a resource provided by your own application or service.
	//
	// For valid AWS service namespace values, see the RegisterScalableTarget
	// action in the Application Auto Scaling API Reference.
	// See: https://docs.aws.amazon.com/autoscaling/application/APIReference/API_RegisterScalableTarget.html
	//
	// Experimental.
	ServiceNamespace ServiceNamespace `json:"serviceNamespace"`
	// Role that allows Application Auto Scaling to modify your scalable target.
	// Experimental.
	Role awsiam.IRole `json:"role"`
}

// A range of metric values in which to apply a certain scaling operation.
// Experimental.
type ScalingInterval struct {
	// The capacity adjustment to apply in this interval.
	//
	// The number is interpreted differently based on AdjustmentType:
	//
	// - ChangeInCapacity: add the adjustment to the current capacity.
	//   The number can be positive or negative.
	// - PercentChangeInCapacity: add or remove the given percentage of the current
	//    capacity to itself. The number can be in the range [-100..100].
	// - ExactCapacity: set the capacity to this number. The number must
	//    be positive.
	// Experimental.
	Change *float64 `json:"change"`
	// The lower bound of the interval.
	//
	// The scaling adjustment will be applied if the metric is higher than this value.
	// Experimental.
	Lower *float64 `json:"lower"`
	// The upper bound of the interval.
	//
	// The scaling adjustment will be applied if the metric is lower than this value.
	// Experimental.
	Upper *float64 `json:"upper"`
}

// A scheduled scaling action.
// Experimental.
type ScalingSchedule struct {
	// When to perform this action.
	// Experimental.
	Schedule Schedule `json:"schedule"`
	// When this scheduled action expires.
	// Experimental.
	EndTime *time.Time `json:"endTime"`
	// The new maximum capacity.
	//
	// During the scheduled time, the current capacity is above the maximum
	// capacity, Application Auto Scaling scales in to the maximum capacity.
	//
	// At least one of maxCapacity and minCapacity must be supplied.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// The new minimum capacity.
	//
	// During the scheduled time, if the current capacity is below the minimum
	// capacity, Application Auto Scaling scales out to the minimum capacity.
	//
	// At least one of maxCapacity and minCapacity must be supplied.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// When this scheduled action becomes active.
	// Experimental.
	StartTime *time.Time `json:"startTime"`
}

// Schedule for scheduled scaling actions.
// Experimental.
type Schedule interface {
	ExpressionString() *string
}

// The jsii proxy struct for Schedule
type jsiiProxy_Schedule struct {
	_ byte // padding
}

func (j *jsiiProxy_Schedule) ExpressionString() *string {
	var returns *string
	_jsii_.Get(
		j,
		"expressionString",
		&returns,
	)
	return returns
}


// Experimental.
func NewSchedule_Override(s Schedule) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.Schedule",
		nil, // no parameters
		s,
	)
}

// Construct a Schedule from a moment in time.
// Experimental.
func Schedule_At(moment *time.Time) Schedule {
	_init_.Initialize()

	var returns Schedule

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.Schedule",
		"at",
		[]interface{}{moment},
		&returns,
	)

	return returns
}

// Create a schedule from a set of cron fields.
// Experimental.
func Schedule_Cron(options *CronOptions) Schedule {
	_init_.Initialize()

	var returns Schedule

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.Schedule",
		"cron",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Construct a schedule from a literal schedule expression.
// Experimental.
func Schedule_Expression(expression *string) Schedule {
	_init_.Initialize()

	var returns Schedule

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.Schedule",
		"expression",
		[]interface{}{expression},
		&returns,
	)

	return returns
}

// Construct a schedule from an interval and a time unit.
// Experimental.
func Schedule_Rate(duration awscdk.Duration) Schedule {
	_init_.Initialize()

	var returns Schedule

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.Schedule",
		"rate",
		[]interface{}{duration},
		&returns,
	)

	return returns
}

// The service that supports Application AutoScaling.
// Experimental.
type ServiceNamespace string

const (
	ServiceNamespace_ECS ServiceNamespace = "ECS"
	ServiceNamespace_ELASTIC_MAP_REDUCE ServiceNamespace = "ELASTIC_MAP_REDUCE"
	ServiceNamespace_EC2 ServiceNamespace = "EC2"
	ServiceNamespace_APPSTREAM ServiceNamespace = "APPSTREAM"
	ServiceNamespace_DYNAMODB ServiceNamespace = "DYNAMODB"
	ServiceNamespace_RDS ServiceNamespace = "RDS"
	ServiceNamespace_SAGEMAKER ServiceNamespace = "SAGEMAKER"
	ServiceNamespace_CUSTOM_RESOURCE ServiceNamespace = "CUSTOM_RESOURCE"
	ServiceNamespace_LAMBDA ServiceNamespace = "LAMBDA"
	ServiceNamespace_COMPREHEND ServiceNamespace = "COMPREHEND"
	ServiceNamespace_KAFKA ServiceNamespace = "KAFKA"
)

// Define a step scaling action.
//
// This kind of scaling policy adjusts the target capacity in configurable
// steps. The size of the step is configurable based on the metric's distance
// to its alarm threshold.
//
// This Action must be used as the target of a CloudWatch alarm to take effect.
// Experimental.
type StepScalingAction interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	ScalingPolicyArn() *string
	AddAdjustment(adjustment *AdjustmentTier)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for StepScalingAction
type jsiiProxy_StepScalingAction struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_StepScalingAction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingAction) ScalingPolicyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalingPolicyArn",
		&returns,
	)
	return returns
}


// Experimental.
func NewStepScalingAction(scope constructs.Construct, id *string, props *StepScalingActionProps) StepScalingAction {
	_init_.Initialize()

	j := jsiiProxy_StepScalingAction{}

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.StepScalingAction",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStepScalingAction_Override(s StepScalingAction, scope constructs.Construct, id *string, props *StepScalingActionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.StepScalingAction",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func StepScalingAction_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.StepScalingAction",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add an adjusment interval to the ScalingAction.
// Experimental.
func (s *jsiiProxy_StepScalingAction) AddAdjustment(adjustment *AdjustmentTier) {
	_jsii_.InvokeVoid(
		s,
		"addAdjustment",
		[]interface{}{adjustment},
	)
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (s *jsiiProxy_StepScalingAction) OnPrepare() {
	_jsii_.InvokeVoid(
		s,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_StepScalingAction) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"onSynthesize",
		[]interface{}{session},
	)
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (s *jsiiProxy_StepScalingAction) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (s *jsiiProxy_StepScalingAction) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_StepScalingAction) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_StepScalingAction) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (s *jsiiProxy_StepScalingAction) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a scaling policy.
// Experimental.
type StepScalingActionProps struct {
	// The scalable target.
	// Experimental.
	ScalingTarget IScalableTarget `json:"scalingTarget"`
	// How the adjustment numbers are interpreted.
	// Experimental.
	AdjustmentType AdjustmentType `json:"adjustmentType"`
	// Grace period after scaling activity.
	//
	// For scale out policies, multiple scale outs during the cooldown period are
	// squashed so that only the biggest scale out happens.
	//
	// For scale in policies, subsequent scale ins during the cooldown period are
	// ignored.
	// See: https://docs.aws.amazon.com/autoscaling/application/APIReference/API_StepScalingPolicyConfiguration.html
	//
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// The aggregation type for the CloudWatch metrics.
	// Experimental.
	MetricAggregationType MetricAggregationType `json:"metricAggregationType"`
	// Minimum absolute number to adjust capacity with as result of percentage scaling.
	//
	// Only when using AdjustmentType = PercentChangeInCapacity, this number controls
	// the minimum absolute effect size.
	// Experimental.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
}

// Define a scaling strategy which scales depending on absolute values of some metric.
//
// You can specify the scaling behavior for various values of the metric.
//
// Implemented using one or more CloudWatch alarms and Step Scaling Policies.
// Experimental.
type StepScalingPolicy interface {
	awscdk.Construct
	LowerAction() StepScalingAction
	LowerAlarm() awscloudwatch.Alarm
	Node() awscdk.ConstructNode
	UpperAction() StepScalingAction
	UpperAlarm() awscloudwatch.Alarm
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for StepScalingPolicy
type jsiiProxy_StepScalingPolicy struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_StepScalingPolicy) LowerAction() StepScalingAction {
	var returns StepScalingAction
	_jsii_.Get(
		j,
		"lowerAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) LowerAlarm() awscloudwatch.Alarm {
	var returns awscloudwatch.Alarm
	_jsii_.Get(
		j,
		"lowerAlarm",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) UpperAction() StepScalingAction {
	var returns StepScalingAction
	_jsii_.Get(
		j,
		"upperAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) UpperAlarm() awscloudwatch.Alarm {
	var returns awscloudwatch.Alarm
	_jsii_.Get(
		j,
		"upperAlarm",
		&returns,
	)
	return returns
}


// Experimental.
func NewStepScalingPolicy(scope constructs.Construct, id *string, props *StepScalingPolicyProps) StepScalingPolicy {
	_init_.Initialize()

	j := jsiiProxy_StepScalingPolicy{}

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.StepScalingPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStepScalingPolicy_Override(s StepScalingPolicy, scope constructs.Construct, id *string, props *StepScalingPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.StepScalingPolicy",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func StepScalingPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.StepScalingPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) OnPrepare() {
	_jsii_.InvokeVoid(
		s,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"onSynthesize",
		[]interface{}{session},
	)
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type StepScalingPolicyProps struct {
	// Metric to scale on.
	// Experimental.
	Metric awscloudwatch.IMetric `json:"metric"`
	// The intervals for scaling.
	//
	// Maps a range of metric values to a particular scaling behavior.
	// Experimental.
	ScalingSteps *[]*ScalingInterval `json:"scalingSteps"`
	// How the adjustment numbers inside 'intervals' are interpreted.
	// Experimental.
	AdjustmentType AdjustmentType `json:"adjustmentType"`
	// Grace period after scaling activity.
	//
	// Subsequent scale outs during the cooldown period are squashed so that only
	// the biggest scale out happens.
	//
	// Subsequent scale ins during the cooldown period are ignored.
	// See: https://docs.aws.amazon.com/autoscaling/application/APIReference/API_StepScalingPolicyConfiguration.html
	//
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// How many evaluation periods of the metric to wait before triggering a scaling action.
	//
	// Raising this value can be used to smooth out the metric, at the expense
	// of slower response times.
	// Experimental.
	EvaluationPeriods *float64 `json:"evaluationPeriods"`
	// Aggregation to apply to all data points over the evaluation periods.
	//
	// Only has meaning if `evaluationPeriods != 1`.
	// Experimental.
	MetricAggregationType MetricAggregationType `json:"metricAggregationType"`
	// Minimum absolute number to adjust capacity with as result of percentage scaling.
	//
	// Only when using AdjustmentType = PercentChangeInCapacity, this number controls
	// the minimum absolute effect size.
	// Experimental.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
	// The scaling target.
	// Experimental.
	ScalingTarget IScalableTarget `json:"scalingTarget"`
}

// Experimental.
type TargetTrackingScalingPolicy interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	ScalingPolicyArn() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for TargetTrackingScalingPolicy
type jsiiProxy_TargetTrackingScalingPolicy struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_TargetTrackingScalingPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetTrackingScalingPolicy) ScalingPolicyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalingPolicyArn",
		&returns,
	)
	return returns
}


// Experimental.
func NewTargetTrackingScalingPolicy(scope constructs.Construct, id *string, props *TargetTrackingScalingPolicyProps) TargetTrackingScalingPolicy {
	_init_.Initialize()

	j := jsiiProxy_TargetTrackingScalingPolicy{}

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.TargetTrackingScalingPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewTargetTrackingScalingPolicy_Override(t TargetTrackingScalingPolicy, scope constructs.Construct, id *string, props *TargetTrackingScalingPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_applicationautoscaling.TargetTrackingScalingPolicy",
		[]interface{}{scope, id, props},
		t,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func TargetTrackingScalingPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_applicationautoscaling.TargetTrackingScalingPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) OnPrepare() {
	_jsii_.InvokeVoid(
		t,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"onSynthesize",
		[]interface{}{session},
	)
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Perform final modifications before synthesis.
//
// This method can be implemented by derived constructs in order to perform
// final changes before synthesis. prepare() will be called after child
// constructs have been prepared.
//
// This is an advanced framework feature. Only use this if you
// understand the implications.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) Prepare() {
	_jsii_.InvokeVoid(
		t,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a concrete TargetTrackingPolicy.
//
// Adds the scalingTarget.
// Experimental.
type TargetTrackingScalingPolicyProps struct {
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the scalable resource. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// scalable resource.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Period after a scale in activity completes before another scale in activity can start.
	// Experimental.
	ScaleInCooldown awscdk.Duration `json:"scaleInCooldown"`
	// Period after a scale out activity completes before another scale out activity can start.
	// Experimental.
	ScaleOutCooldown awscdk.Duration `json:"scaleOutCooldown"`
	// The target value for the metric.
	// Experimental.
	TargetValue *float64 `json:"targetValue"`
	// A custom metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	CustomMetric awscloudwatch.IMetric `json:"customMetric"`
	// A predefined metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	PredefinedMetric PredefinedMetric `json:"predefinedMetric"`
	// Identify the resource associated with the metric type.
	//
	// Only used for predefined metric ALBRequestCountPerTarget.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	ResourceLabel *string `json:"resourceLabel"`
	// Experimental.
	ScalingTarget IScalableTarget `json:"scalingTarget"`
}


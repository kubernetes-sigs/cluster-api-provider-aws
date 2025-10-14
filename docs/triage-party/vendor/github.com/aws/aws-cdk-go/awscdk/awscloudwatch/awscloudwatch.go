package awscloudwatch

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch/internal"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
)

// An alarm on a CloudWatch metric.
// Experimental.
type Alarm interface {
	AlarmBase
	AlarmActionArns() *[]*string
	SetAlarmActionArns(val *[]*string)
	AlarmArn() *string
	AlarmName() *string
	Env() *awscdk.ResourceEnvironment
	InsufficientDataActionArns() *[]*string
	SetInsufficientDataActionArns(val *[]*string)
	Metric() IMetric
	Node() awscdk.ConstructNode
	OkActionArns() *[]*string
	SetOkActionArns(val *[]*string)
	PhysicalName() *string
	Stack() awscdk.Stack
	AddAlarmAction(actions ...IAlarmAction)
	AddInsufficientDataAction(actions ...IAlarmAction)
	AddOkAction(actions ...IAlarmAction)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderAlarmRule() *string
	Synthesize(session awscdk.ISynthesisSession)
	ToAnnotation() *HorizontalAnnotation
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Alarm
type jsiiProxy_Alarm struct {
	jsiiProxy_AlarmBase
}

func (j *jsiiProxy_Alarm) AlarmActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"alarmActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) AlarmArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) AlarmName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) InsufficientDataActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"insufficientDataActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) Metric() IMetric {
	var returns IMetric
	_jsii_.Get(
		j,
		"metric",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) OkActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"okActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alarm) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewAlarm(scope constructs.Construct, id *string, props *AlarmProps) Alarm {
	_init_.Initialize()

	j := jsiiProxy_Alarm{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Alarm",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewAlarm_Override(a Alarm, scope constructs.Construct, id *string, props *AlarmProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Alarm",
		[]interface{}{scope, id, props},
		a,
	)
}

func (j *jsiiProxy_Alarm) SetAlarmActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"alarmActionArns",
		val,
	)
}

func (j *jsiiProxy_Alarm) SetInsufficientDataActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"insufficientDataActionArns",
		val,
	)
}

func (j *jsiiProxy_Alarm) SetOkActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"okActionArns",
		val,
	)
}

// Import an existing CloudWatch alarm provided an ARN.
// Experimental.
func Alarm_FromAlarmArn(scope constructs.Construct, id *string, alarmArn *string) IAlarm {
	_init_.Initialize()

	var returns IAlarm

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.Alarm",
		"fromAlarmArn",
		[]interface{}{scope, id, alarmArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Alarm_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.Alarm",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Alarm_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.Alarm",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Trigger this action if the alarm fires.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (a *jsiiProxy_Alarm) AddAlarmAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addAlarmAction",
		args,
	)
}

// Trigger this action if there is insufficient data to evaluate the alarm.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (a *jsiiProxy_Alarm) AddInsufficientDataAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addInsufficientDataAction",
		args,
	)
}

// Trigger this action if the alarm returns from breaching state into ok state.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (a *jsiiProxy_Alarm) AddOkAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addOkAction",
		args,
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
func (a *jsiiProxy_Alarm) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_Alarm) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Alarm) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Alarm) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Alarm) OnPrepare() {
	_jsii_.InvokeVoid(
		a,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_Alarm) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
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
func (a *jsiiProxy_Alarm) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Alarm) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// AlarmRule indicating ALARM state for Alarm.
// Experimental.
func (a *jsiiProxy_Alarm) RenderAlarmRule() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"renderAlarmRule",
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
func (a *jsiiProxy_Alarm) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Turn this alarm into a horizontal annotation.
//
// This is useful if you want to represent an Alarm in a non-AlarmWidget.
// An `AlarmWidget` can directly show an alarm, but it can only show a
// single alarm and no other metrics. Instead, you can convert the alarm to
// a HorizontalAnnotation and add it as an annotation to another graph.
//
// This might be useful if:
//
// - You want to show multiple alarms inside a single graph, for example if
//    you have both a "small margin/long period" alarm as well as a
//    "large margin/short period" alarm.
//
// - You want to show an Alarm line in a graph with multiple metrics in it.
// Experimental.
func (a *jsiiProxy_Alarm) ToAnnotation() *HorizontalAnnotation {
	var returns *HorizontalAnnotation

	_jsii_.Invoke(
		a,
		"toAnnotation",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_Alarm) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Alarm) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for an alarm action.
// Experimental.
type AlarmActionConfig struct {
	// Return the ARN that should be used for a CloudWatch Alarm action.
	// Experimental.
	AlarmActionArn *string `json:"alarmActionArn"`
}

// The base class for Alarm and CompositeAlarm resources.
// Experimental.
type AlarmBase interface {
	awscdk.Resource
	IAlarm
	AlarmActionArns() *[]*string
	SetAlarmActionArns(val *[]*string)
	AlarmArn() *string
	AlarmName() *string
	Env() *awscdk.ResourceEnvironment
	InsufficientDataActionArns() *[]*string
	SetInsufficientDataActionArns(val *[]*string)
	Node() awscdk.ConstructNode
	OkActionArns() *[]*string
	SetOkActionArns(val *[]*string)
	PhysicalName() *string
	Stack() awscdk.Stack
	AddAlarmAction(actions ...IAlarmAction)
	AddInsufficientDataAction(actions ...IAlarmAction)
	AddOkAction(actions ...IAlarmAction)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderAlarmRule() *string
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for AlarmBase
type jsiiProxy_AlarmBase struct {
	internal.Type__awscdkResource
	jsiiProxy_IAlarm
}

func (j *jsiiProxy_AlarmBase) AlarmActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"alarmActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) AlarmArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) AlarmName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) InsufficientDataActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"insufficientDataActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) OkActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"okActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewAlarmBase_Override(a AlarmBase, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.AlarmBase",
		[]interface{}{scope, id, props},
		a,
	)
}

func (j *jsiiProxy_AlarmBase) SetAlarmActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"alarmActionArns",
		val,
	)
}

func (j *jsiiProxy_AlarmBase) SetInsufficientDataActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"insufficientDataActionArns",
		val,
	)
}

func (j *jsiiProxy_AlarmBase) SetOkActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"okActionArns",
		val,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func AlarmBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func AlarmBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Trigger this action if the alarm fires.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (a *jsiiProxy_AlarmBase) AddAlarmAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addAlarmAction",
		args,
	)
}

// Trigger this action if there is insufficient data to evaluate the alarm.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (a *jsiiProxy_AlarmBase) AddInsufficientDataAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addInsufficientDataAction",
		args,
	)
}

// Trigger this action if the alarm returns from breaching state into ok state.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (a *jsiiProxy_AlarmBase) AddOkAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addOkAction",
		args,
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
func (a *jsiiProxy_AlarmBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_AlarmBase) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_AlarmBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_AlarmBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_AlarmBase) OnPrepare() {
	_jsii_.InvokeVoid(
		a,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_AlarmBase) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
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
func (a *jsiiProxy_AlarmBase) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_AlarmBase) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// AlarmRule indicating ALARM state for Alarm.
// Experimental.
func (a *jsiiProxy_AlarmBase) RenderAlarmRule() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"renderAlarmRule",
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
func (a *jsiiProxy_AlarmBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_AlarmBase) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_AlarmBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for Alarms.
// Experimental.
type AlarmProps struct {
	// The number of periods over which data is compared to the specified threshold.
	// Experimental.
	EvaluationPeriods *float64 `json:"evaluationPeriods"`
	// The value against which the specified statistic is compared.
	// Experimental.
	Threshold *float64 `json:"threshold"`
	// Whether the actions for this alarm are enabled.
	// Experimental.
	ActionsEnabled *bool `json:"actionsEnabled"`
	// Description for the alarm.
	// Experimental.
	AlarmDescription *string `json:"alarmDescription"`
	// Name of the alarm.
	// Experimental.
	AlarmName *string `json:"alarmName"`
	// Comparison to use to check if metric is breaching.
	// Experimental.
	ComparisonOperator ComparisonOperator `json:"comparisonOperator"`
	// The number of datapoints that must be breaching to trigger the alarm.
	//
	// This is used only if you are setting an "M
	// out of N" alarm. In that case, this value is the M. For more information, see Evaluating an Alarm in the Amazon
	// CloudWatch User Guide.
	// See: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/AlarmThatSendsEmail.html#alarm-evaluation
	//
	// Experimental.
	DatapointsToAlarm *float64 `json:"datapointsToAlarm"`
	// Specifies whether to evaluate the data and potentially change the alarm state if there are too few data points to be statistically significant.
	//
	// Used only for alarms that are based on percentiles.
	// Experimental.
	EvaluateLowSampleCountPercentile *string `json:"evaluateLowSampleCountPercentile"`
	// The period over which the specified statistic is applied.
	//
	// Cannot be used with `MathExpression` objects.
	// Deprecated: Use `metric.with({ period: ... })` to encode the period into the Metric object
	Period awscdk.Duration `json:"period"`
	// What function to use for aggregating.
	//
	// Can be one of the following:
	//
	// - "Minimum" | "min"
	// - "Maximum" | "max"
	// - "Average" | "avg"
	// - "Sum" | "sum"
	// - "SampleCount | "n"
	// - "pNN.NN"
	//
	// Cannot be used with `MathExpression` objects.
	// Deprecated: Use `metric.with({ statistic: ... })` to encode the period into the Metric object
	Statistic *string `json:"statistic"`
	// Sets how this alarm is to handle missing data points.
	// Experimental.
	TreatMissingData TreatMissingData `json:"treatMissingData"`
	// The metric to add the alarm on.
	//
	// Metric objects can be obtained from most resources, or you can construct
	// custom Metric objects by instantiating one.
	// Experimental.
	Metric IMetric `json:"metric"`
}

// Class with static functions to build AlarmRule for Composite Alarms.
// Experimental.
type AlarmRule interface {
}

// The jsii proxy struct for AlarmRule
type jsiiProxy_AlarmRule struct {
	_ byte // padding
}

// Experimental.
func NewAlarmRule() AlarmRule {
	_init_.Initialize()

	j := jsiiProxy_AlarmRule{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.AlarmRule",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewAlarmRule_Override(a AlarmRule) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.AlarmRule",
		nil, // no parameters
		a,
	)
}

// function to join all provided AlarmRules with AND operator.
// Experimental.
func AlarmRule_AllOf(operands ...IAlarmRule) IAlarmRule {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range operands {
		args = append(args, a)
	}

	var returns IAlarmRule

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmRule",
		"allOf",
		args,
		&returns,
	)

	return returns
}

// function to join all provided AlarmRules with OR operator.
// Experimental.
func AlarmRule_AnyOf(operands ...IAlarmRule) IAlarmRule {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range operands {
		args = append(args, a)
	}

	var returns IAlarmRule

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmRule",
		"anyOf",
		args,
		&returns,
	)

	return returns
}

// function to build Rule Expression for given IAlarm and AlarmState.
// Experimental.
func AlarmRule_FromAlarm(alarm IAlarm, alarmState AlarmState) IAlarmRule {
	_init_.Initialize()

	var returns IAlarmRule

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmRule",
		"fromAlarm",
		[]interface{}{alarm, alarmState},
		&returns,
	)

	return returns
}

// function to build TRUE/FALSE intent for Rule Expression.
// Experimental.
func AlarmRule_FromBoolean(value *bool) IAlarmRule {
	_init_.Initialize()

	var returns IAlarmRule

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmRule",
		"fromBoolean",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// function to build Rule Expression for given Alarm Rule string.
// Experimental.
func AlarmRule_FromString(alarmRule *string) IAlarmRule {
	_init_.Initialize()

	var returns IAlarmRule

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmRule",
		"fromString",
		[]interface{}{alarmRule},
		&returns,
	)

	return returns
}

// function to wrap provided AlarmRule in NOT operator.
// Experimental.
func AlarmRule_Not(operand IAlarmRule) IAlarmRule {
	_init_.Initialize()

	var returns IAlarmRule

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.AlarmRule",
		"not",
		[]interface{}{operand},
		&returns,
	)

	return returns
}

// Enumeration indicates state of Alarm used in building Alarm Rule.
// Experimental.
type AlarmState string

const (
	AlarmState_ALARM AlarmState = "ALARM"
	AlarmState_OK AlarmState = "OK"
	AlarmState_INSUFFICIENT_DATA AlarmState = "INSUFFICIENT_DATA"
)

// A dashboard widget that displays alarms in a grid view.
// Experimental.
type AlarmStatusWidget interface {
	ConcreteWidget
	Height() *float64
	Width() *float64
	X() *float64
	SetX(val *float64)
	Y() *float64
	SetY(val *float64)
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for AlarmStatusWidget
type jsiiProxy_AlarmStatusWidget struct {
	jsiiProxy_ConcreteWidget
}

func (j *jsiiProxy_AlarmStatusWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmStatusWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmStatusWidget) X() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"x",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmStatusWidget) Y() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"y",
		&returns,
	)
	return returns
}


// Experimental.
func NewAlarmStatusWidget(props *AlarmStatusWidgetProps) AlarmStatusWidget {
	_init_.Initialize()

	j := jsiiProxy_AlarmStatusWidget{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.AlarmStatusWidget",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewAlarmStatusWidget_Override(a AlarmStatusWidget, props *AlarmStatusWidgetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.AlarmStatusWidget",
		[]interface{}{props},
		a,
	)
}

func (j *jsiiProxy_AlarmStatusWidget) SetX(val *float64) {
	_jsii_.Set(
		j,
		"x",
		val,
	)
}

func (j *jsiiProxy_AlarmStatusWidget) SetY(val *float64) {
	_jsii_.Set(
		j,
		"y",
		val,
	)
}

// Place the widget at a given position.
// Experimental.
func (a *jsiiProxy_AlarmStatusWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		a,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (a *jsiiProxy_AlarmStatusWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		a,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for an Alarm Status Widget.
// Experimental.
type AlarmStatusWidgetProps struct {
	// CloudWatch Alarms to show in widget.
	// Experimental.
	Alarms *[]IAlarm `json:"alarms"`
	// Height of the widget.
	// Experimental.
	Height *float64 `json:"height"`
	// The title of the widget.
	// Experimental.
	Title *string `json:"title"`
	// Width of the widget, in a grid of 24 units wide.
	// Experimental.
	Width *float64 `json:"width"`
}

// Display the metric associated with an alarm, including the alarm line.
// Experimental.
type AlarmWidget interface {
	ConcreteWidget
	Height() *float64
	Width() *float64
	X() *float64
	SetX(val *float64)
	Y() *float64
	SetY(val *float64)
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for AlarmWidget
type jsiiProxy_AlarmWidget struct {
	jsiiProxy_ConcreteWidget
}

func (j *jsiiProxy_AlarmWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmWidget) X() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"x",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AlarmWidget) Y() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"y",
		&returns,
	)
	return returns
}


// Experimental.
func NewAlarmWidget(props *AlarmWidgetProps) AlarmWidget {
	_init_.Initialize()

	j := jsiiProxy_AlarmWidget{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.AlarmWidget",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewAlarmWidget_Override(a AlarmWidget, props *AlarmWidgetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.AlarmWidget",
		[]interface{}{props},
		a,
	)
}

func (j *jsiiProxy_AlarmWidget) SetX(val *float64) {
	_jsii_.Set(
		j,
		"x",
		val,
	)
}

func (j *jsiiProxy_AlarmWidget) SetY(val *float64) {
	_jsii_.Set(
		j,
		"y",
		val,
	)
}

// Place the widget at a given position.
// Experimental.
func (a *jsiiProxy_AlarmWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		a,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (a *jsiiProxy_AlarmWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		a,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for an AlarmWidget.
// Experimental.
type AlarmWidgetProps struct {
	// Height of the widget.
	// Experimental.
	Height *float64 `json:"height"`
	// The region the metrics of this graph should be taken from.
	// Experimental.
	Region *string `json:"region"`
	// Title for the graph.
	// Experimental.
	Title *string `json:"title"`
	// Width of the widget, in a grid of 24 units wide.
	// Experimental.
	Width *float64 `json:"width"`
	// The alarm to show.
	// Experimental.
	Alarm IAlarm `json:"alarm"`
	// Left Y axis.
	// Experimental.
	LeftYAxis *YAxisProps `json:"leftYAxis"`
}

// A CloudFormation `AWS::CloudWatch::Alarm`.
type CfnAlarm interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ActionsEnabled() interface{}
	SetActionsEnabled(val interface{})
	AlarmActions() *[]*string
	SetAlarmActions(val *[]*string)
	AlarmDescription() *string
	SetAlarmDescription(val *string)
	AlarmName() *string
	SetAlarmName(val *string)
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ComparisonOperator() *string
	SetComparisonOperator(val *string)
	CreationStack() *[]*string
	DatapointsToAlarm() *float64
	SetDatapointsToAlarm(val *float64)
	Dimensions() interface{}
	SetDimensions(val interface{})
	EvaluateLowSampleCountPercentile() *string
	SetEvaluateLowSampleCountPercentile(val *string)
	EvaluationPeriods() *float64
	SetEvaluationPeriods(val *float64)
	ExtendedStatistic() *string
	SetExtendedStatistic(val *string)
	InsufficientDataActions() *[]*string
	SetInsufficientDataActions(val *[]*string)
	LogicalId() *string
	MetricName() *string
	SetMetricName(val *string)
	Metrics() interface{}
	SetMetrics(val interface{})
	Namespace() *string
	SetNamespace(val *string)
	Node() awscdk.ConstructNode
	OkActions() *[]*string
	SetOkActions(val *[]*string)
	Period() *float64
	SetPeriod(val *float64)
	Ref() *string
	Stack() awscdk.Stack
	Statistic() *string
	SetStatistic(val *string)
	Threshold() *float64
	SetThreshold(val *float64)
	ThresholdMetricId() *string
	SetThresholdMetricId(val *string)
	TreatMissingData() *string
	SetTreatMissingData(val *string)
	Unit() *string
	SetUnit(val *string)
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

// The jsii proxy struct for CfnAlarm
type jsiiProxy_CfnAlarm struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAlarm) ActionsEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"actionsEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) AlarmActions() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"alarmActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) AlarmDescription() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmDescription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) AlarmName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) ComparisonOperator() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comparisonOperator",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) DatapointsToAlarm() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"datapointsToAlarm",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Dimensions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"dimensions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) EvaluateLowSampleCountPercentile() *string {
	var returns *string
	_jsii_.Get(
		j,
		"evaluateLowSampleCountPercentile",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) EvaluationPeriods() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"evaluationPeriods",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) ExtendedStatistic() *string {
	var returns *string
	_jsii_.Get(
		j,
		"extendedStatistic",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) InsufficientDataActions() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"insufficientDataActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) MetricName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"metricName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Metrics() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"metrics",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Namespace() *string {
	var returns *string
	_jsii_.Get(
		j,
		"namespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) OkActions() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"okActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Period() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"period",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Statistic() *string {
	var returns *string
	_jsii_.Get(
		j,
		"statistic",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Threshold() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"threshold",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) ThresholdMetricId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"thresholdMetricId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) TreatMissingData() *string {
	var returns *string
	_jsii_.Get(
		j,
		"treatMissingData",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) Unit() *string {
	var returns *string
	_jsii_.Get(
		j,
		"unit",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlarm) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudWatch::Alarm`.
func NewCfnAlarm(scope awscdk.Construct, id *string, props *CfnAlarmProps) CfnAlarm {
	_init_.Initialize()

	j := jsiiProxy_CfnAlarm{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnAlarm",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudWatch::Alarm`.
func NewCfnAlarm_Override(c CfnAlarm, scope awscdk.Construct, id *string, props *CfnAlarmProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnAlarm",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAlarm) SetActionsEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"actionsEnabled",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetAlarmActions(val *[]*string) {
	_jsii_.Set(
		j,
		"alarmActions",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetAlarmDescription(val *string) {
	_jsii_.Set(
		j,
		"alarmDescription",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetAlarmName(val *string) {
	_jsii_.Set(
		j,
		"alarmName",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetComparisonOperator(val *string) {
	_jsii_.Set(
		j,
		"comparisonOperator",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetDatapointsToAlarm(val *float64) {
	_jsii_.Set(
		j,
		"datapointsToAlarm",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetDimensions(val interface{}) {
	_jsii_.Set(
		j,
		"dimensions",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetEvaluateLowSampleCountPercentile(val *string) {
	_jsii_.Set(
		j,
		"evaluateLowSampleCountPercentile",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetEvaluationPeriods(val *float64) {
	_jsii_.Set(
		j,
		"evaluationPeriods",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetExtendedStatistic(val *string) {
	_jsii_.Set(
		j,
		"extendedStatistic",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetInsufficientDataActions(val *[]*string) {
	_jsii_.Set(
		j,
		"insufficientDataActions",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetMetricName(val *string) {
	_jsii_.Set(
		j,
		"metricName",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetMetrics(val interface{}) {
	_jsii_.Set(
		j,
		"metrics",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetNamespace(val *string) {
	_jsii_.Set(
		j,
		"namespace",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetOkActions(val *[]*string) {
	_jsii_.Set(
		j,
		"okActions",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetPeriod(val *float64) {
	_jsii_.Set(
		j,
		"period",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetStatistic(val *string) {
	_jsii_.Set(
		j,
		"statistic",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetThreshold(val *float64) {
	_jsii_.Set(
		j,
		"threshold",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetThresholdMetricId(val *string) {
	_jsii_.Set(
		j,
		"thresholdMetricId",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetTreatMissingData(val *string) {
	_jsii_.Set(
		j,
		"treatMissingData",
		val,
	)
}

func (j *jsiiProxy_CfnAlarm) SetUnit(val *string) {
	_jsii_.Set(
		j,
		"unit",
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
func CfnAlarm_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnAlarm",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAlarm_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnAlarm",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAlarm_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnAlarm",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAlarm_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.CfnAlarm",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAlarm) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAlarm) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAlarm) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAlarm) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAlarm) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAlarm) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAlarm) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAlarm) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAlarm) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAlarm) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAlarm) OnPrepare() {
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
func (c *jsiiProxy_CfnAlarm) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAlarm) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAlarm) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAlarm) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAlarm) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAlarm) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAlarm) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAlarm) ToString() *string {
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
func (c *jsiiProxy_CfnAlarm) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAlarm) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnAlarm_DimensionProperty struct {
	// `CfnAlarm.DimensionProperty.Name`.
	Name *string `json:"name"`
	// `CfnAlarm.DimensionProperty.Value`.
	Value *string `json:"value"`
}

type CfnAlarm_MetricDataQueryProperty struct {
	// `CfnAlarm.MetricDataQueryProperty.Id`.
	Id *string `json:"id"`
	// `CfnAlarm.MetricDataQueryProperty.Expression`.
	Expression *string `json:"expression"`
	// `CfnAlarm.MetricDataQueryProperty.Label`.
	Label *string `json:"label"`
	// `CfnAlarm.MetricDataQueryProperty.MetricStat`.
	MetricStat interface{} `json:"metricStat"`
	// `CfnAlarm.MetricDataQueryProperty.Period`.
	Period *float64 `json:"period"`
	// `CfnAlarm.MetricDataQueryProperty.ReturnData`.
	ReturnData interface{} `json:"returnData"`
}

type CfnAlarm_MetricProperty struct {
	// `CfnAlarm.MetricProperty.Dimensions`.
	Dimensions interface{} `json:"dimensions"`
	// `CfnAlarm.MetricProperty.MetricName`.
	MetricName *string `json:"metricName"`
	// `CfnAlarm.MetricProperty.Namespace`.
	Namespace *string `json:"namespace"`
}

type CfnAlarm_MetricStatProperty struct {
	// `CfnAlarm.MetricStatProperty.Metric`.
	Metric interface{} `json:"metric"`
	// `CfnAlarm.MetricStatProperty.Period`.
	Period *float64 `json:"period"`
	// `CfnAlarm.MetricStatProperty.Stat`.
	Stat *string `json:"stat"`
	// `CfnAlarm.MetricStatProperty.Unit`.
	Unit *string `json:"unit"`
}

// Properties for defining a `AWS::CloudWatch::Alarm`.
type CfnAlarmProps struct {
	// `AWS::CloudWatch::Alarm.ComparisonOperator`.
	ComparisonOperator *string `json:"comparisonOperator"`
	// `AWS::CloudWatch::Alarm.EvaluationPeriods`.
	EvaluationPeriods *float64 `json:"evaluationPeriods"`
	// `AWS::CloudWatch::Alarm.ActionsEnabled`.
	ActionsEnabled interface{} `json:"actionsEnabled"`
	// `AWS::CloudWatch::Alarm.AlarmActions`.
	AlarmActions *[]*string `json:"alarmActions"`
	// `AWS::CloudWatch::Alarm.AlarmDescription`.
	AlarmDescription *string `json:"alarmDescription"`
	// `AWS::CloudWatch::Alarm.AlarmName`.
	AlarmName *string `json:"alarmName"`
	// `AWS::CloudWatch::Alarm.DatapointsToAlarm`.
	DatapointsToAlarm *float64 `json:"datapointsToAlarm"`
	// `AWS::CloudWatch::Alarm.Dimensions`.
	Dimensions interface{} `json:"dimensions"`
	// `AWS::CloudWatch::Alarm.EvaluateLowSampleCountPercentile`.
	EvaluateLowSampleCountPercentile *string `json:"evaluateLowSampleCountPercentile"`
	// `AWS::CloudWatch::Alarm.ExtendedStatistic`.
	ExtendedStatistic *string `json:"extendedStatistic"`
	// `AWS::CloudWatch::Alarm.InsufficientDataActions`.
	InsufficientDataActions *[]*string `json:"insufficientDataActions"`
	// `AWS::CloudWatch::Alarm.MetricName`.
	MetricName *string `json:"metricName"`
	// `AWS::CloudWatch::Alarm.Metrics`.
	Metrics interface{} `json:"metrics"`
	// `AWS::CloudWatch::Alarm.Namespace`.
	Namespace *string `json:"namespace"`
	// `AWS::CloudWatch::Alarm.OKActions`.
	OkActions *[]*string `json:"okActions"`
	// `AWS::CloudWatch::Alarm.Period`.
	Period *float64 `json:"period"`
	// `AWS::CloudWatch::Alarm.Statistic`.
	Statistic *string `json:"statistic"`
	// `AWS::CloudWatch::Alarm.Threshold`.
	Threshold *float64 `json:"threshold"`
	// `AWS::CloudWatch::Alarm.ThresholdMetricId`.
	ThresholdMetricId *string `json:"thresholdMetricId"`
	// `AWS::CloudWatch::Alarm.TreatMissingData`.
	TreatMissingData *string `json:"treatMissingData"`
	// `AWS::CloudWatch::Alarm.Unit`.
	Unit *string `json:"unit"`
}

// A CloudFormation `AWS::CloudWatch::AnomalyDetector`.
type CfnAnomalyDetector interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Configuration() interface{}
	SetConfiguration(val interface{})
	CreationStack() *[]*string
	Dimensions() interface{}
	SetDimensions(val interface{})
	LogicalId() *string
	MetricName() *string
	SetMetricName(val *string)
	Namespace() *string
	SetNamespace(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	Stat() *string
	SetStat(val *string)
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

// The jsii proxy struct for CfnAnomalyDetector
type jsiiProxy_CfnAnomalyDetector struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAnomalyDetector) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) Configuration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"configuration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) Dimensions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"dimensions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) MetricName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"metricName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) Namespace() *string {
	var returns *string
	_jsii_.Get(
		j,
		"namespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) Stat() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stat",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAnomalyDetector) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudWatch::AnomalyDetector`.
func NewCfnAnomalyDetector(scope awscdk.Construct, id *string, props *CfnAnomalyDetectorProps) CfnAnomalyDetector {
	_init_.Initialize()

	j := jsiiProxy_CfnAnomalyDetector{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnAnomalyDetector",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudWatch::AnomalyDetector`.
func NewCfnAnomalyDetector_Override(c CfnAnomalyDetector, scope awscdk.Construct, id *string, props *CfnAnomalyDetectorProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnAnomalyDetector",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAnomalyDetector) SetConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"configuration",
		val,
	)
}

func (j *jsiiProxy_CfnAnomalyDetector) SetDimensions(val interface{}) {
	_jsii_.Set(
		j,
		"dimensions",
		val,
	)
}

func (j *jsiiProxy_CfnAnomalyDetector) SetMetricName(val *string) {
	_jsii_.Set(
		j,
		"metricName",
		val,
	)
}

func (j *jsiiProxy_CfnAnomalyDetector) SetNamespace(val *string) {
	_jsii_.Set(
		j,
		"namespace",
		val,
	)
}

func (j *jsiiProxy_CfnAnomalyDetector) SetStat(val *string) {
	_jsii_.Set(
		j,
		"stat",
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
func CfnAnomalyDetector_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnAnomalyDetector",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAnomalyDetector_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnAnomalyDetector",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAnomalyDetector_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnAnomalyDetector",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAnomalyDetector_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.CfnAnomalyDetector",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAnomalyDetector) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAnomalyDetector) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAnomalyDetector) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAnomalyDetector) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAnomalyDetector) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAnomalyDetector) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAnomalyDetector) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAnomalyDetector) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAnomalyDetector) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAnomalyDetector) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAnomalyDetector) OnPrepare() {
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
func (c *jsiiProxy_CfnAnomalyDetector) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAnomalyDetector) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAnomalyDetector) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAnomalyDetector) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAnomalyDetector) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAnomalyDetector) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAnomalyDetector) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAnomalyDetector) ToString() *string {
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
func (c *jsiiProxy_CfnAnomalyDetector) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAnomalyDetector) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnAnomalyDetector_ConfigurationProperty struct {
	// `CfnAnomalyDetector.ConfigurationProperty.ExcludedTimeRanges`.
	ExcludedTimeRanges interface{} `json:"excludedTimeRanges"`
	// `CfnAnomalyDetector.ConfigurationProperty.MetricTimeZone`.
	MetricTimeZone *string `json:"metricTimeZone"`
}

type CfnAnomalyDetector_DimensionProperty struct {
	// `CfnAnomalyDetector.DimensionProperty.Name`.
	Name *string `json:"name"`
	// `CfnAnomalyDetector.DimensionProperty.Value`.
	Value *string `json:"value"`
}

type CfnAnomalyDetector_RangeProperty struct {
	// `CfnAnomalyDetector.RangeProperty.EndTime`.
	EndTime *string `json:"endTime"`
	// `CfnAnomalyDetector.RangeProperty.StartTime`.
	StartTime *string `json:"startTime"`
}

// Properties for defining a `AWS::CloudWatch::AnomalyDetector`.
type CfnAnomalyDetectorProps struct {
	// `AWS::CloudWatch::AnomalyDetector.MetricName`.
	MetricName *string `json:"metricName"`
	// `AWS::CloudWatch::AnomalyDetector.Namespace`.
	Namespace *string `json:"namespace"`
	// `AWS::CloudWatch::AnomalyDetector.Stat`.
	Stat *string `json:"stat"`
	// `AWS::CloudWatch::AnomalyDetector.Configuration`.
	Configuration interface{} `json:"configuration"`
	// `AWS::CloudWatch::AnomalyDetector.Dimensions`.
	Dimensions interface{} `json:"dimensions"`
}

// A CloudFormation `AWS::CloudWatch::CompositeAlarm`.
type CfnCompositeAlarm interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ActionsEnabled() interface{}
	SetActionsEnabled(val interface{})
	AlarmActions() *[]*string
	SetAlarmActions(val *[]*string)
	AlarmDescription() *string
	SetAlarmDescription(val *string)
	AlarmName() *string
	SetAlarmName(val *string)
	AlarmRule() *string
	SetAlarmRule(val *string)
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	InsufficientDataActions() *[]*string
	SetInsufficientDataActions(val *[]*string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	OkActions() *[]*string
	SetOkActions(val *[]*string)
	Ref() *string
	Stack() awscdk.Stack
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

// The jsii proxy struct for CfnCompositeAlarm
type jsiiProxy_CfnCompositeAlarm struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnCompositeAlarm) ActionsEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"actionsEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) AlarmActions() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"alarmActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) AlarmDescription() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmDescription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) AlarmName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) AlarmRule() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmRule",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) InsufficientDataActions() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"insufficientDataActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) OkActions() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"okActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCompositeAlarm) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudWatch::CompositeAlarm`.
func NewCfnCompositeAlarm(scope awscdk.Construct, id *string, props *CfnCompositeAlarmProps) CfnCompositeAlarm {
	_init_.Initialize()

	j := jsiiProxy_CfnCompositeAlarm{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnCompositeAlarm",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudWatch::CompositeAlarm`.
func NewCfnCompositeAlarm_Override(c CfnCompositeAlarm, scope awscdk.Construct, id *string, props *CfnCompositeAlarmProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnCompositeAlarm",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCompositeAlarm) SetActionsEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"actionsEnabled",
		val,
	)
}

func (j *jsiiProxy_CfnCompositeAlarm) SetAlarmActions(val *[]*string) {
	_jsii_.Set(
		j,
		"alarmActions",
		val,
	)
}

func (j *jsiiProxy_CfnCompositeAlarm) SetAlarmDescription(val *string) {
	_jsii_.Set(
		j,
		"alarmDescription",
		val,
	)
}

func (j *jsiiProxy_CfnCompositeAlarm) SetAlarmName(val *string) {
	_jsii_.Set(
		j,
		"alarmName",
		val,
	)
}

func (j *jsiiProxy_CfnCompositeAlarm) SetAlarmRule(val *string) {
	_jsii_.Set(
		j,
		"alarmRule",
		val,
	)
}

func (j *jsiiProxy_CfnCompositeAlarm) SetInsufficientDataActions(val *[]*string) {
	_jsii_.Set(
		j,
		"insufficientDataActions",
		val,
	)
}

func (j *jsiiProxy_CfnCompositeAlarm) SetOkActions(val *[]*string) {
	_jsii_.Set(
		j,
		"okActions",
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
func CfnCompositeAlarm_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnCompositeAlarm",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnCompositeAlarm_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnCompositeAlarm",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCompositeAlarm_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnCompositeAlarm",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnCompositeAlarm_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.CfnCompositeAlarm",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnCompositeAlarm) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnCompositeAlarm) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnCompositeAlarm) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnCompositeAlarm) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnCompositeAlarm) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnCompositeAlarm) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnCompositeAlarm) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnCompositeAlarm) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnCompositeAlarm) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnCompositeAlarm) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnCompositeAlarm) OnPrepare() {
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
func (c *jsiiProxy_CfnCompositeAlarm) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCompositeAlarm) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCompositeAlarm) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCompositeAlarm) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnCompositeAlarm) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnCompositeAlarm) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnCompositeAlarm) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCompositeAlarm) ToString() *string {
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
func (c *jsiiProxy_CfnCompositeAlarm) Validate() *[]*string {
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
func (c *jsiiProxy_CfnCompositeAlarm) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudWatch::CompositeAlarm`.
type CfnCompositeAlarmProps struct {
	// `AWS::CloudWatch::CompositeAlarm.AlarmName`.
	AlarmName *string `json:"alarmName"`
	// `AWS::CloudWatch::CompositeAlarm.AlarmRule`.
	AlarmRule *string `json:"alarmRule"`
	// `AWS::CloudWatch::CompositeAlarm.ActionsEnabled`.
	ActionsEnabled interface{} `json:"actionsEnabled"`
	// `AWS::CloudWatch::CompositeAlarm.AlarmActions`.
	AlarmActions *[]*string `json:"alarmActions"`
	// `AWS::CloudWatch::CompositeAlarm.AlarmDescription`.
	AlarmDescription *string `json:"alarmDescription"`
	// `AWS::CloudWatch::CompositeAlarm.InsufficientDataActions`.
	InsufficientDataActions *[]*string `json:"insufficientDataActions"`
	// `AWS::CloudWatch::CompositeAlarm.OKActions`.
	OkActions *[]*string `json:"okActions"`
}

// A CloudFormation `AWS::CloudWatch::Dashboard`.
type CfnDashboard interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DashboardBody() *string
	SetDashboardBody(val *string)
	DashboardName() *string
	SetDashboardName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
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

// The jsii proxy struct for CfnDashboard
type jsiiProxy_CfnDashboard struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDashboard) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) DashboardBody() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dashboardBody",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) DashboardName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dashboardName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDashboard) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudWatch::Dashboard`.
func NewCfnDashboard(scope awscdk.Construct, id *string, props *CfnDashboardProps) CfnDashboard {
	_init_.Initialize()

	j := jsiiProxy_CfnDashboard{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnDashboard",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudWatch::Dashboard`.
func NewCfnDashboard_Override(c CfnDashboard, scope awscdk.Construct, id *string, props *CfnDashboardProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnDashboard",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDashboard) SetDashboardBody(val *string) {
	_jsii_.Set(
		j,
		"dashboardBody",
		val,
	)
}

func (j *jsiiProxy_CfnDashboard) SetDashboardName(val *string) {
	_jsii_.Set(
		j,
		"dashboardName",
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
func CfnDashboard_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnDashboard",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnDashboard_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnDashboard",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnDashboard_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnDashboard",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDashboard_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.CfnDashboard",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnDashboard) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnDashboard) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnDashboard) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnDashboard) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnDashboard) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnDashboard) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnDashboard) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnDashboard) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnDashboard) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnDashboard) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnDashboard) OnPrepare() {
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
func (c *jsiiProxy_CfnDashboard) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDashboard) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnDashboard) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnDashboard) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnDashboard) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnDashboard) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnDashboard) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDashboard) ToString() *string {
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
func (c *jsiiProxy_CfnDashboard) Validate() *[]*string {
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
func (c *jsiiProxy_CfnDashboard) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudWatch::Dashboard`.
type CfnDashboardProps struct {
	// `AWS::CloudWatch::Dashboard.DashboardBody`.
	DashboardBody *string `json:"dashboardBody"`
	// `AWS::CloudWatch::Dashboard.DashboardName`.
	DashboardName *string `json:"dashboardName"`
}

// A CloudFormation `AWS::CloudWatch::InsightRule`.
type CfnInsightRule interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrRuleName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RuleBody() *string
	SetRuleBody(val *string)
	RuleName() *string
	SetRuleName(val *string)
	RuleState() *string
	SetRuleState(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
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

// The jsii proxy struct for CfnInsightRule
type jsiiProxy_CfnInsightRule struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnInsightRule) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) AttrRuleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRuleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) RuleBody() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ruleBody",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) RuleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ruleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) RuleState() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ruleState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInsightRule) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudWatch::InsightRule`.
func NewCfnInsightRule(scope awscdk.Construct, id *string, props *CfnInsightRuleProps) CfnInsightRule {
	_init_.Initialize()

	j := jsiiProxy_CfnInsightRule{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnInsightRule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudWatch::InsightRule`.
func NewCfnInsightRule_Override(c CfnInsightRule, scope awscdk.Construct, id *string, props *CfnInsightRuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnInsightRule",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnInsightRule) SetRuleBody(val *string) {
	_jsii_.Set(
		j,
		"ruleBody",
		val,
	)
}

func (j *jsiiProxy_CfnInsightRule) SetRuleName(val *string) {
	_jsii_.Set(
		j,
		"ruleName",
		val,
	)
}

func (j *jsiiProxy_CfnInsightRule) SetRuleState(val *string) {
	_jsii_.Set(
		j,
		"ruleState",
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
func CfnInsightRule_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnInsightRule",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnInsightRule_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnInsightRule",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnInsightRule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnInsightRule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnInsightRule_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.CfnInsightRule",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnInsightRule) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnInsightRule) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnInsightRule) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnInsightRule) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnInsightRule) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnInsightRule) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnInsightRule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnInsightRule) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnInsightRule) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnInsightRule) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnInsightRule) OnPrepare() {
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
func (c *jsiiProxy_CfnInsightRule) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnInsightRule) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnInsightRule) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnInsightRule) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnInsightRule) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnInsightRule) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnInsightRule) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnInsightRule) ToString() *string {
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
func (c *jsiiProxy_CfnInsightRule) Validate() *[]*string {
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
func (c *jsiiProxy_CfnInsightRule) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudWatch::InsightRule`.
type CfnInsightRuleProps struct {
	// `AWS::CloudWatch::InsightRule.RuleBody`.
	RuleBody *string `json:"ruleBody"`
	// `AWS::CloudWatch::InsightRule.RuleName`.
	RuleName *string `json:"ruleName"`
	// `AWS::CloudWatch::InsightRule.RuleState`.
	RuleState *string `json:"ruleState"`
	// `AWS::CloudWatch::InsightRule.Tags`.
	Tags interface{} `json:"tags"`
}

// A CloudFormation `AWS::CloudWatch::MetricStream`.
type CfnMetricStream interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrCreationDate() *string
	AttrLastUpdateDate() *string
	AttrState() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	ExcludeFilters() interface{}
	SetExcludeFilters(val interface{})
	FirehoseArn() *string
	SetFirehoseArn(val *string)
	IncludeFilters() interface{}
	SetIncludeFilters(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	OutputFormat() *string
	SetOutputFormat(val *string)
	Ref() *string
	RoleArn() *string
	SetRoleArn(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
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

// The jsii proxy struct for CfnMetricStream
type jsiiProxy_CfnMetricStream struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnMetricStream) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) AttrCreationDate() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCreationDate",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) AttrLastUpdateDate() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrLastUpdateDate",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) AttrState() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) ExcludeFilters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"excludeFilters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) FirehoseArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"firehoseArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) IncludeFilters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"includeFilters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) OutputFormat() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputFormat",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMetricStream) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudWatch::MetricStream`.
func NewCfnMetricStream(scope awscdk.Construct, id *string, props *CfnMetricStreamProps) CfnMetricStream {
	_init_.Initialize()

	j := jsiiProxy_CfnMetricStream{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnMetricStream",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudWatch::MetricStream`.
func NewCfnMetricStream_Override(c CfnMetricStream, scope awscdk.Construct, id *string, props *CfnMetricStreamProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CfnMetricStream",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnMetricStream) SetExcludeFilters(val interface{}) {
	_jsii_.Set(
		j,
		"excludeFilters",
		val,
	)
}

func (j *jsiiProxy_CfnMetricStream) SetFirehoseArn(val *string) {
	_jsii_.Set(
		j,
		"firehoseArn",
		val,
	)
}

func (j *jsiiProxy_CfnMetricStream) SetIncludeFilters(val interface{}) {
	_jsii_.Set(
		j,
		"includeFilters",
		val,
	)
}

func (j *jsiiProxy_CfnMetricStream) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnMetricStream) SetOutputFormat(val *string) {
	_jsii_.Set(
		j,
		"outputFormat",
		val,
	)
}

func (j *jsiiProxy_CfnMetricStream) SetRoleArn(val *string) {
	_jsii_.Set(
		j,
		"roleArn",
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
func CfnMetricStream_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnMetricStream",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnMetricStream_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnMetricStream",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnMetricStream_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CfnMetricStream",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnMetricStream_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.CfnMetricStream",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnMetricStream) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnMetricStream) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnMetricStream) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnMetricStream) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnMetricStream) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnMetricStream) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnMetricStream) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnMetricStream) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnMetricStream) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnMetricStream) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnMetricStream) OnPrepare() {
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
func (c *jsiiProxy_CfnMetricStream) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnMetricStream) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnMetricStream) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnMetricStream) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnMetricStream) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnMetricStream) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnMetricStream) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnMetricStream) ToString() *string {
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
func (c *jsiiProxy_CfnMetricStream) Validate() *[]*string {
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
func (c *jsiiProxy_CfnMetricStream) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnMetricStream_MetricStreamFilterProperty struct {
	// `CfnMetricStream.MetricStreamFilterProperty.Namespace`.
	Namespace *string `json:"namespace"`
}

// Properties for defining a `AWS::CloudWatch::MetricStream`.
type CfnMetricStreamProps struct {
	// `AWS::CloudWatch::MetricStream.FirehoseArn`.
	FirehoseArn *string `json:"firehoseArn"`
	// `AWS::CloudWatch::MetricStream.OutputFormat`.
	OutputFormat *string `json:"outputFormat"`
	// `AWS::CloudWatch::MetricStream.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `AWS::CloudWatch::MetricStream.ExcludeFilters`.
	ExcludeFilters interface{} `json:"excludeFilters"`
	// `AWS::CloudWatch::MetricStream.IncludeFilters`.
	IncludeFilters interface{} `json:"includeFilters"`
	// `AWS::CloudWatch::MetricStream.Name`.
	Name *string `json:"name"`
	// `AWS::CloudWatch::MetricStream.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A set of standard colours that can be used in annotations in a GraphWidget.
// Experimental.
type Color interface {
}

// The jsii proxy struct for Color
type jsiiProxy_Color struct {
	_ byte // padding
}

// Experimental.
func NewColor() Color {
	_init_.Initialize()

	j := jsiiProxy_Color{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Color",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewColor_Override(c Color) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Color",
		nil, // no parameters
		c,
	)
}

func Color_BLUE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"BLUE",
		&returns,
	)
	return returns
}

func Color_BROWN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"BROWN",
		&returns,
	)
	return returns
}

func Color_GREEN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"GREEN",
		&returns,
	)
	return returns
}

func Color_GREY() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"GREY",
		&returns,
	)
	return returns
}

func Color_ORANGE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"ORANGE",
		&returns,
	)
	return returns
}

func Color_PINK() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"PINK",
		&returns,
	)
	return returns
}

func Color_PURPLE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"PURPLE",
		&returns,
	)
	return returns
}

func Color_RED() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cloudwatch.Color",
		"RED",
		&returns,
	)
	return returns
}

// A widget that contains other widgets in a vertical column.
//
// Widgets will be laid out next to each other
// Experimental.
type Column interface {
	IWidget
	Height() *float64
	Width() *float64
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for Column
type jsiiProxy_Column struct {
	jsiiProxy_IWidget
}

func (j *jsiiProxy_Column) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Column) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}


// Experimental.
func NewColumn(widgets ...IWidget) Column {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range widgets {
		args = append(args, a)
	}

	j := jsiiProxy_Column{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Column",
		args,
		&j,
	)

	return &j
}

// Experimental.
func NewColumn_Override(c Column, widgets ...IWidget) {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range widgets {
		args = append(args, a)
	}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Column",
		args,
		c,
	)
}

// Place the widget at a given position.
// Experimental.
func (c *jsiiProxy_Column) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		c,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (c *jsiiProxy_Column) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		c,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options shared by most methods accepting metric options.
// Experimental.
type CommonMetricOptions struct {
	// Account which this metric comes from.
	// Experimental.
	Account *string `json:"account"`
	// The hex color code, prefixed with '#' (e.g. '#00ff00'), to use when this metric is rendered on a graph. The `Color` class has a set of standard colors that can be used here.
	// Experimental.
	Color *string `json:"color"`
	// Dimensions of the metric.
	// Deprecated: Use 'dimensionsMap' instead.
	Dimensions *map[string]interface{} `json:"dimensions"`
	// Dimensions of the metric.
	// Experimental.
	DimensionsMap *map[string]*string `json:"dimensionsMap"`
	// Label for this metric when added to a Graph in a Dashboard.
	// Experimental.
	Label *string `json:"label"`
	// The period over which the specified statistic is applied.
	// Experimental.
	Period awscdk.Duration `json:"period"`
	// Region which this metric comes from.
	// Experimental.
	Region *string `json:"region"`
	// What function to use for aggregating.
	//
	// Can be one of the following:
	//
	// - "Minimum" | "min"
	// - "Maximum" | "max"
	// - "Average" | "avg"
	// - "Sum" | "sum"
	// - "SampleCount | "n"
	// - "pNN.NN"
	// Experimental.
	Statistic *string `json:"statistic"`
	// Unit used to filter the metric stream.
	//
	// Only refer to datums emitted to the metric stream with the given unit and
	// ignore all others. Only useful when datums are being emitted to the same
	// metric stream under different units.
	//
	// The default is to use all matric datums in the stream, regardless of unit,
	// which is recommended in nearly all cases.
	//
	// CloudWatch does not honor this property for graphs.
	// Experimental.
	Unit Unit `json:"unit"`
}

// Comparison operator for evaluating alarms.
// Experimental.
type ComparisonOperator string

const (
	ComparisonOperator_GREATER_THAN_OR_EQUAL_TO_THRESHOLD ComparisonOperator = "GREATER_THAN_OR_EQUAL_TO_THRESHOLD"
	ComparisonOperator_GREATER_THAN_THRESHOLD ComparisonOperator = "GREATER_THAN_THRESHOLD"
	ComparisonOperator_LESS_THAN_THRESHOLD ComparisonOperator = "LESS_THAN_THRESHOLD"
	ComparisonOperator_LESS_THAN_OR_EQUAL_TO_THRESHOLD ComparisonOperator = "LESS_THAN_OR_EQUAL_TO_THRESHOLD"
	ComparisonOperator_LESS_THAN_LOWER_OR_GREATER_THAN_UPPER_THRESHOLD ComparisonOperator = "LESS_THAN_LOWER_OR_GREATER_THAN_UPPER_THRESHOLD"
	ComparisonOperator_GREATER_THAN_UPPER_THRESHOLD ComparisonOperator = "GREATER_THAN_UPPER_THRESHOLD"
	ComparisonOperator_LESS_THAN_LOWER_THRESHOLD ComparisonOperator = "LESS_THAN_LOWER_THRESHOLD"
)

// A Composite Alarm based on Alarm Rule.
// Experimental.
type CompositeAlarm interface {
	AlarmBase
	AlarmActionArns() *[]*string
	SetAlarmActionArns(val *[]*string)
	AlarmArn() *string
	AlarmName() *string
	Env() *awscdk.ResourceEnvironment
	InsufficientDataActionArns() *[]*string
	SetInsufficientDataActionArns(val *[]*string)
	Node() awscdk.ConstructNode
	OkActionArns() *[]*string
	SetOkActionArns(val *[]*string)
	PhysicalName() *string
	Stack() awscdk.Stack
	AddAlarmAction(actions ...IAlarmAction)
	AddInsufficientDataAction(actions ...IAlarmAction)
	AddOkAction(actions ...IAlarmAction)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderAlarmRule() *string
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CompositeAlarm
type jsiiProxy_CompositeAlarm struct {
	jsiiProxy_AlarmBase
}

func (j *jsiiProxy_CompositeAlarm) AlarmActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"alarmActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) AlarmArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) AlarmName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) InsufficientDataActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"insufficientDataActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) OkActionArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"okActionArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositeAlarm) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewCompositeAlarm(scope constructs.Construct, id *string, props *CompositeAlarmProps) CompositeAlarm {
	_init_.Initialize()

	j := jsiiProxy_CompositeAlarm{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CompositeAlarm",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCompositeAlarm_Override(c CompositeAlarm, scope constructs.Construct, id *string, props *CompositeAlarmProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.CompositeAlarm",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CompositeAlarm) SetAlarmActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"alarmActionArns",
		val,
	)
}

func (j *jsiiProxy_CompositeAlarm) SetInsufficientDataActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"insufficientDataActionArns",
		val,
	)
}

func (j *jsiiProxy_CompositeAlarm) SetOkActionArns(val *[]*string) {
	_jsii_.Set(
		j,
		"okActionArns",
		val,
	)
}

// Import an existing CloudWatch composite alarm provided an ARN.
// Experimental.
func CompositeAlarm_FromCompositeAlarmArn(scope constructs.Construct, id *string, compositeAlarmArn *string) IAlarm {
	_init_.Initialize()

	var returns IAlarm

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CompositeAlarm",
		"fromCompositeAlarmArn",
		[]interface{}{scope, id, compositeAlarmArn},
		&returns,
	)

	return returns
}

// Import an existing CloudWatch composite alarm provided an Name.
// Experimental.
func CompositeAlarm_FromCompositeAlarmName(scope constructs.Construct, id *string, compositeAlarmName *string) IAlarm {
	_init_.Initialize()

	var returns IAlarm

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CompositeAlarm",
		"fromCompositeAlarmName",
		[]interface{}{scope, id, compositeAlarmName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CompositeAlarm_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CompositeAlarm",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func CompositeAlarm_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.CompositeAlarm",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Trigger this action if the alarm fires.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (c *jsiiProxy_CompositeAlarm) AddAlarmAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addAlarmAction",
		args,
	)
}

// Trigger this action if there is insufficient data to evaluate the alarm.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (c *jsiiProxy_CompositeAlarm) AddInsufficientDataAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addInsufficientDataAction",
		args,
	)
}

// Trigger this action if the alarm returns from breaching state into ok state.
//
// Typically the ARN of an SNS topic or ARN of an AutoScaling policy.
// Experimental.
func (c *jsiiProxy_CompositeAlarm) AddOkAction(actions ...IAlarmAction) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addOkAction",
		args,
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
func (c *jsiiProxy_CompositeAlarm) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		c,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (c *jsiiProxy_CompositeAlarm) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CompositeAlarm) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CompositeAlarm) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CompositeAlarm) OnPrepare() {
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
func (c *jsiiProxy_CompositeAlarm) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CompositeAlarm) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CompositeAlarm) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// AlarmRule indicating ALARM state for Alarm.
// Experimental.
func (c *jsiiProxy_CompositeAlarm) RenderAlarmRule() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"renderAlarmRule",
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
func (c *jsiiProxy_CompositeAlarm) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CompositeAlarm) ToString() *string {
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
func (c *jsiiProxy_CompositeAlarm) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for creating a Composite Alarm.
// Experimental.
type CompositeAlarmProps struct {
	// Expression that specifies which other alarms are to be evaluated to determine this composite alarm's state.
	// Experimental.
	AlarmRule IAlarmRule `json:"alarmRule"`
	// Whether the actions for this alarm are enabled.
	// Experimental.
	ActionsEnabled *bool `json:"actionsEnabled"`
	// Description for the alarm.
	// Experimental.
	AlarmDescription *string `json:"alarmDescription"`
	// Name of the alarm.
	// Experimental.
	CompositeAlarmName *string `json:"compositeAlarmName"`
}

// A real CloudWatch widget that has its own fixed size and remembers its position.
//
// This is in contrast to other widgets which exist for layout purposes.
// Experimental.
type ConcreteWidget interface {
	IWidget
	Height() *float64
	Width() *float64
	X() *float64
	SetX(val *float64)
	Y() *float64
	SetY(val *float64)
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for ConcreteWidget
type jsiiProxy_ConcreteWidget struct {
	jsiiProxy_IWidget
}

func (j *jsiiProxy_ConcreteWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConcreteWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConcreteWidget) X() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"x",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConcreteWidget) Y() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"y",
		&returns,
	)
	return returns
}


// Experimental.
func NewConcreteWidget_Override(c ConcreteWidget, width *float64, height *float64) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.ConcreteWidget",
		[]interface{}{width, height},
		c,
	)
}

func (j *jsiiProxy_ConcreteWidget) SetX(val *float64) {
	_jsii_.Set(
		j,
		"x",
		val,
	)
}

func (j *jsiiProxy_ConcreteWidget) SetY(val *float64) {
	_jsii_.Set(
		j,
		"y",
		val,
	)
}

// Place the widget at a given position.
// Experimental.
func (c *jsiiProxy_ConcreteWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		c,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (c *jsiiProxy_ConcreteWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		c,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties needed to make an alarm from a metric.
// Experimental.
type CreateAlarmOptions struct {
	// The number of periods over which data is compared to the specified threshold.
	// Experimental.
	EvaluationPeriods *float64 `json:"evaluationPeriods"`
	// The value against which the specified statistic is compared.
	// Experimental.
	Threshold *float64 `json:"threshold"`
	// Whether the actions for this alarm are enabled.
	// Experimental.
	ActionsEnabled *bool `json:"actionsEnabled"`
	// Description for the alarm.
	// Experimental.
	AlarmDescription *string `json:"alarmDescription"`
	// Name of the alarm.
	// Experimental.
	AlarmName *string `json:"alarmName"`
	// Comparison to use to check if metric is breaching.
	// Experimental.
	ComparisonOperator ComparisonOperator `json:"comparisonOperator"`
	// The number of datapoints that must be breaching to trigger the alarm.
	//
	// This is used only if you are setting an "M
	// out of N" alarm. In that case, this value is the M. For more information, see Evaluating an Alarm in the Amazon
	// CloudWatch User Guide.
	// See: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/AlarmThatSendsEmail.html#alarm-evaluation
	//
	// Experimental.
	DatapointsToAlarm *float64 `json:"datapointsToAlarm"`
	// Specifies whether to evaluate the data and potentially change the alarm state if there are too few data points to be statistically significant.
	//
	// Used only for alarms that are based on percentiles.
	// Experimental.
	EvaluateLowSampleCountPercentile *string `json:"evaluateLowSampleCountPercentile"`
	// The period over which the specified statistic is applied.
	//
	// Cannot be used with `MathExpression` objects.
	// Deprecated: Use `metric.with({ period: ... })` to encode the period into the Metric object
	Period awscdk.Duration `json:"period"`
	// What function to use for aggregating.
	//
	// Can be one of the following:
	//
	// - "Minimum" | "min"
	// - "Maximum" | "max"
	// - "Average" | "avg"
	// - "Sum" | "sum"
	// - "SampleCount | "n"
	// - "pNN.NN"
	//
	// Cannot be used with `MathExpression` objects.
	// Deprecated: Use `metric.with({ statistic: ... })` to encode the period into the Metric object
	Statistic *string `json:"statistic"`
	// Sets how this alarm is to handle missing data points.
	// Experimental.
	TreatMissingData TreatMissingData `json:"treatMissingData"`
}

// A CloudWatch dashboard.
// Experimental.
type Dashboard interface {
	awscdk.Resource
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddWidgets(widgets ...IWidget)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Dashboard
type jsiiProxy_Dashboard struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_Dashboard) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Dashboard) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Dashboard) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Dashboard) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewDashboard(scope constructs.Construct, id *string, props *DashboardProps) Dashboard {
	_init_.Initialize()

	j := jsiiProxy_Dashboard{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Dashboard",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewDashboard_Override(d Dashboard, scope constructs.Construct, id *string, props *DashboardProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Dashboard",
		[]interface{}{scope, id, props},
		d,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Dashboard_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.Dashboard",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Dashboard_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.Dashboard",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a widget to the dashboard.
//
// Widgets given in multiple calls to add() will be laid out stacked on
// top of each other.
//
// Multiple widgets added in the same call to add() will be laid out next
// to each other.
// Experimental.
func (d *jsiiProxy_Dashboard) AddWidgets(widgets ...IWidget) {
	args := []interface{}{}
	for _, a := range widgets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		d,
		"addWidgets",
		args,
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
func (d *jsiiProxy_Dashboard) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		d,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (d *jsiiProxy_Dashboard) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Dashboard) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Dashboard) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Dashboard) OnPrepare() {
	_jsii_.InvokeVoid(
		d,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (d *jsiiProxy_Dashboard) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
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
func (d *jsiiProxy_Dashboard) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Dashboard) Prepare() {
	_jsii_.InvokeVoid(
		d,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (d *jsiiProxy_Dashboard) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (d *jsiiProxy_Dashboard) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Dashboard) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		d,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining a CloudWatch Dashboard.
// Experimental.
type DashboardProps struct {
	// Name of the dashboard.
	//
	// If set, must only contain alphanumerics, dash (-) and underscore (_)
	// Experimental.
	DashboardName *string `json:"dashboardName"`
	// The end of the time range to use for each widget on the dashboard when the dashboard loads.
	//
	// If you specify a value for end, you must also specify a value for start.
	// Specify an absolute time in the ISO 8601 format. For example, 2018-12-17T06:00:00.000Z.
	// Experimental.
	End *string `json:"end"`
	// Use this field to specify the period for the graphs when the dashboard loads.
	//
	// Specifying `Auto` causes the period of all graphs on the dashboard to automatically adapt to the time range of the dashboard.
	// Specifying `Inherit` ensures that the period set for each graph is always obeyed.
	// Experimental.
	PeriodOverride PeriodOverride `json:"periodOverride"`
	// The start of the time range to use for each widget on the dashboard.
	//
	// You can specify start without specifying end to specify a relative time range that ends with the current time.
	// In this case, the value of start must begin with -P, and you can use M, H, D, W and M as abbreviations for
	// minutes, hours, days, weeks and months. For example, -PT8H shows the last 8 hours and -P3M shows the last three months.
	// You can also use start along with an end field, to specify an absolute time range.
	// When specifying an absolute time range, use the ISO 8601 format. For example, 2018-12-17T06:00:00.000Z.
	// Experimental.
	Start *string `json:"start"`
	// Initial set of widgets on the dashboard.
	//
	// One array represents a row of widgets.
	// Experimental.
	Widgets *[]*[]IWidget `json:"widgets"`
}

// Metric dimension.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cw-dimension.html
//
// Experimental.
type Dimension struct {
	// Name of the dimension.
	// Experimental.
	Name *string `json:"name"`
	// Value of the dimension.
	// Experimental.
	Value interface{} `json:"value"`
}

// A dashboard widget that displays metrics.
// Experimental.
type GraphWidget interface {
	ConcreteWidget
	Height() *float64
	Width() *float64
	X() *float64
	SetX(val *float64)
	Y() *float64
	SetY(val *float64)
	AddLeftMetric(metric IMetric)
	AddRightMetric(metric IMetric)
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for GraphWidget
type jsiiProxy_GraphWidget struct {
	jsiiProxy_ConcreteWidget
}

func (j *jsiiProxy_GraphWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GraphWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GraphWidget) X() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"x",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GraphWidget) Y() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"y",
		&returns,
	)
	return returns
}


// Experimental.
func NewGraphWidget(props *GraphWidgetProps) GraphWidget {
	_init_.Initialize()

	j := jsiiProxy_GraphWidget{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.GraphWidget",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewGraphWidget_Override(g GraphWidget, props *GraphWidgetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.GraphWidget",
		[]interface{}{props},
		g,
	)
}

func (j *jsiiProxy_GraphWidget) SetX(val *float64) {
	_jsii_.Set(
		j,
		"x",
		val,
	)
}

func (j *jsiiProxy_GraphWidget) SetY(val *float64) {
	_jsii_.Set(
		j,
		"y",
		val,
	)
}

// Add another metric to the left Y axis of the GraphWidget.
// Experimental.
func (g *jsiiProxy_GraphWidget) AddLeftMetric(metric IMetric) {
	_jsii_.InvokeVoid(
		g,
		"addLeftMetric",
		[]interface{}{metric},
	)
}

// Add another metric to the right Y axis of the GraphWidget.
// Experimental.
func (g *jsiiProxy_GraphWidget) AddRightMetric(metric IMetric) {
	_jsii_.InvokeVoid(
		g,
		"addRightMetric",
		[]interface{}{metric},
	)
}

// Place the widget at a given position.
// Experimental.
func (g *jsiiProxy_GraphWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		g,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (g *jsiiProxy_GraphWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		g,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a GraphWidget.
// Experimental.
type GraphWidgetProps struct {
	// Height of the widget.
	// Experimental.
	Height *float64 `json:"height"`
	// The region the metrics of this graph should be taken from.
	// Experimental.
	Region *string `json:"region"`
	// Title for the graph.
	// Experimental.
	Title *string `json:"title"`
	// Width of the widget, in a grid of 24 units wide.
	// Experimental.
	Width *float64 `json:"width"`
	// Metrics to display on left Y axis.
	// Experimental.
	Left *[]IMetric `json:"left"`
	// Annotations for the left Y axis.
	// Experimental.
	LeftAnnotations *[]*HorizontalAnnotation `json:"leftAnnotations"`
	// Left Y axis.
	// Experimental.
	LeftYAxis *YAxisProps `json:"leftYAxis"`
	// Position of the legend.
	// Experimental.
	LegendPosition LegendPosition `json:"legendPosition"`
	// Whether the graph should show live data.
	// Experimental.
	LiveData *bool `json:"liveData"`
	// The default period for all metrics in this widget.
	//
	// The period is the length of time represented by one data point on the graph.
	// This default can be overridden within each metric definition.
	// Experimental.
	Period awscdk.Duration `json:"period"`
	// Metrics to display on right Y axis.
	// Experimental.
	Right *[]IMetric `json:"right"`
	// Annotations for the right Y axis.
	// Experimental.
	RightAnnotations *[]*HorizontalAnnotation `json:"rightAnnotations"`
	// Right Y axis.
	// Experimental.
	RightYAxis *YAxisProps `json:"rightYAxis"`
	// Whether to show the value from the entire time range. Only applicable for Bar and Pie charts.
	//
	// If false, values will be from the most recent period of your chosen time range;
	// if true, shows the value from the entire time range.
	// Experimental.
	SetPeriodToTimeRange *bool `json:"setPeriodToTimeRange"`
	// Whether the graph should be shown as stacked lines.
	// Experimental.
	Stacked *bool `json:"stacked"`
	// The default statistic to be displayed for each metric.
	//
	// This default can be overridden within the definition of each individual metric
	// Experimental.
	Statistic *string `json:"statistic"`
	// Display this metric.
	// Experimental.
	View GraphWidgetView `json:"view"`
}

// Types of view.
// Experimental.
type GraphWidgetView string

const (
	GraphWidgetView_TIME_SERIES GraphWidgetView = "TIME_SERIES"
	GraphWidgetView_BAR GraphWidgetView = "BAR"
	GraphWidgetView_PIE GraphWidgetView = "PIE"
)

// Horizontal annotation to be added to a graph.
// Experimental.
type HorizontalAnnotation struct {
	// The value of the annotation.
	// Experimental.
	Value *float64 `json:"value"`
	// The hex color code, prefixed with '#' (e.g. '#00ff00'), to be used for the annotation. The `Color` class has a set of standard colors that can be used here.
	// Experimental.
	Color *string `json:"color"`
	// Add shading above or below the annotation.
	// Experimental.
	Fill Shading `json:"fill"`
	// Label for the annotation.
	// Experimental.
	Label *string `json:"label"`
	// Whether the annotation is visible.
	// Experimental.
	Visible *bool `json:"visible"`
}

// Represents a CloudWatch Alarm.
// Experimental.
type IAlarm interface {
	IAlarmRule
	awscdk.IResource
	// Alarm ARN (i.e. arn:aws:cloudwatch:<region>:<account-id>:alarm:Foo).
	// Experimental.
	AlarmArn() *string
	// Name of the alarm.
	// Experimental.
	AlarmName() *string
}

// The jsii proxy for IAlarm
type jsiiProxy_IAlarm struct {
	jsiiProxy_IAlarmRule
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IAlarm) RenderAlarmRule() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"renderAlarmRule",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IAlarm) AlarmArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAlarm) AlarmName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alarmName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAlarm) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAlarm) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAlarm) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// Interface for objects that can be the targets of CloudWatch alarm actions.
// Experimental.
type IAlarmAction interface {
	// Return the properties required to send alarm actions to this CloudWatch alarm.
	// Experimental.
	Bind(scope awscdk.Construct, alarm IAlarm) *AlarmActionConfig
}

// The jsii proxy for IAlarmAction
type jsiiProxy_IAlarmAction struct {
	_ byte // padding
}

func (i *jsiiProxy_IAlarmAction) Bind(scope awscdk.Construct, alarm IAlarm) *AlarmActionConfig {
	var returns *AlarmActionConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, alarm},
		&returns,
	)

	return returns
}

// Interface for Alarm Rule.
// Experimental.
type IAlarmRule interface {
	// serialized representation of Alarm Rule to be used when building the Composite Alarm resource.
	// Experimental.
	RenderAlarmRule() *string
}

// The jsii proxy for IAlarmRule
type jsiiProxy_IAlarmRule struct {
	_ byte // padding
}

func (i *jsiiProxy_IAlarmRule) RenderAlarmRule() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"renderAlarmRule",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Interface for metrics.
// Experimental.
type IMetric interface {
	// Turn this metric object into an alarm configuration.
	// Deprecated: Use `toMetricConfig()` instead.
	ToAlarmConfig() *MetricAlarmConfig
	// Turn this metric object into a graph configuration.
	// Deprecated: Use `toMetricConfig()` instead.
	ToGraphConfig() *MetricGraphConfig
	// Inspect the details of the metric object.
	// Experimental.
	ToMetricConfig() *MetricConfig
}

// The jsii proxy for IMetric
type jsiiProxy_IMetric struct {
	_ byte // padding
}

func (i *jsiiProxy_IMetric) ToAlarmConfig() *MetricAlarmConfig {
	var returns *MetricAlarmConfig

	_jsii_.Invoke(
		i,
		"toAlarmConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IMetric) ToGraphConfig() *MetricGraphConfig {
	var returns *MetricGraphConfig

	_jsii_.Invoke(
		i,
		"toGraphConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IMetric) ToMetricConfig() *MetricConfig {
	var returns *MetricConfig

	_jsii_.Invoke(
		i,
		"toMetricConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A single dashboard widget.
// Experimental.
type IWidget interface {
	// Place the widget at a given position.
	// Experimental.
	Position(x *float64, y *float64)
	// Return the widget JSON for use in the dashboard.
	// Experimental.
	ToJson() *[]interface{}
	// The amount of vertical grid units the widget will take up.
	// Experimental.
	Height() *float64
	// The amount of horizontal grid units the widget will take up.
	// Experimental.
	Width() *float64
}

// The jsii proxy for IWidget
type jsiiProxy_IWidget struct {
	_ byte // padding
}

func (i *jsiiProxy_IWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		i,
		"position",
		[]interface{}{x, y},
	)
}

func (i *jsiiProxy_IWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		i,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

// The position of the legend on a GraphWidget.
// Experimental.
type LegendPosition string

const (
	LegendPosition_BOTTOM LegendPosition = "BOTTOM"
	LegendPosition_RIGHT LegendPosition = "RIGHT"
	LegendPosition_HIDDEN LegendPosition = "HIDDEN"
)

// Types of view.
// Experimental.
type LogQueryVisualizationType string

const (
	LogQueryVisualizationType_TABLE LogQueryVisualizationType = "TABLE"
	LogQueryVisualizationType_LINE LogQueryVisualizationType = "LINE"
	LogQueryVisualizationType_STACKEDAREA LogQueryVisualizationType = "STACKEDAREA"
	LogQueryVisualizationType_BAR LogQueryVisualizationType = "BAR"
	LogQueryVisualizationType_PIE LogQueryVisualizationType = "PIE"
)

// Display query results from Logs Insights.
// Experimental.
type LogQueryWidget interface {
	ConcreteWidget
	Height() *float64
	Width() *float64
	X() *float64
	SetX(val *float64)
	Y() *float64
	SetY(val *float64)
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for LogQueryWidget
type jsiiProxy_LogQueryWidget struct {
	jsiiProxy_ConcreteWidget
}

func (j *jsiiProxy_LogQueryWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LogQueryWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LogQueryWidget) X() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"x",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LogQueryWidget) Y() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"y",
		&returns,
	)
	return returns
}


// Experimental.
func NewLogQueryWidget(props *LogQueryWidgetProps) LogQueryWidget {
	_init_.Initialize()

	j := jsiiProxy_LogQueryWidget{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.LogQueryWidget",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewLogQueryWidget_Override(l LogQueryWidget, props *LogQueryWidgetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.LogQueryWidget",
		[]interface{}{props},
		l,
	)
}

func (j *jsiiProxy_LogQueryWidget) SetX(val *float64) {
	_jsii_.Set(
		j,
		"x",
		val,
	)
}

func (j *jsiiProxy_LogQueryWidget) SetY(val *float64) {
	_jsii_.Set(
		j,
		"y",
		val,
	)
}

// Place the widget at a given position.
// Experimental.
func (l *jsiiProxy_LogQueryWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		l,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (l *jsiiProxy_LogQueryWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		l,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a Query widget.
// Experimental.
type LogQueryWidgetProps struct {
	// Names of log groups to query.
	// Experimental.
	LogGroupNames *[]*string `json:"logGroupNames"`
	// Height of the widget.
	// Experimental.
	Height *float64 `json:"height"`
	// A sequence of lines to use to build the query.
	//
	// The query will be built by joining the lines together using `\n|`.
	// Experimental.
	QueryLines *[]*string `json:"queryLines"`
	// Full query string for log insights.
	//
	// Be sure to prepend every new line with a newline and pipe character
	// (`\n|`).
	// Experimental.
	QueryString *string `json:"queryString"`
	// The region the metrics of this widget should be taken from.
	// Experimental.
	Region *string `json:"region"`
	// Title for the widget.
	// Experimental.
	Title *string `json:"title"`
	// The type of view to use.
	// Experimental.
	View LogQueryVisualizationType `json:"view"`
	// Width of the widget, in a grid of 24 units wide.
	// Experimental.
	Width *float64 `json:"width"`
}

// A math expression built with metric(s) emitted by a service.
//
// The math expression is a combination of an expression (x+y) and metrics to apply expression on.
// It also contains metadata which is used only in graphs, such as color and label.
// It makes sense to embed this in here, so that compound constructs can attach
// that metadata to metrics they expose.
//
// This class does not represent a resource, so hence is not a construct. Instead,
// MathExpression is an abstraction that makes it easy to specify metrics for use in both
// alarms and graphs.
// Experimental.
type MathExpression interface {
	IMetric
	Color() *string
	Expression() *string
	Label() *string
	Period() awscdk.Duration
	UsingMetrics() *map[string]IMetric
	CreateAlarm(scope awscdk.Construct, id *string, props *CreateAlarmOptions) Alarm
	ToAlarmConfig() *MetricAlarmConfig
	ToGraphConfig() *MetricGraphConfig
	ToMetricConfig() *MetricConfig
	ToString() *string
	With(props *MathExpressionOptions) MathExpression
}

// The jsii proxy struct for MathExpression
type jsiiProxy_MathExpression struct {
	jsiiProxy_IMetric
}

func (j *jsiiProxy_MathExpression) Color() *string {
	var returns *string
	_jsii_.Get(
		j,
		"color",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_MathExpression) Expression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"expression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_MathExpression) Label() *string {
	var returns *string
	_jsii_.Get(
		j,
		"label",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_MathExpression) Period() awscdk.Duration {
	var returns awscdk.Duration
	_jsii_.Get(
		j,
		"period",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_MathExpression) UsingMetrics() *map[string]IMetric {
	var returns *map[string]IMetric
	_jsii_.Get(
		j,
		"usingMetrics",
		&returns,
	)
	return returns
}


// Experimental.
func NewMathExpression(props *MathExpressionProps) MathExpression {
	_init_.Initialize()

	j := jsiiProxy_MathExpression{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.MathExpression",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewMathExpression_Override(m MathExpression, props *MathExpressionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.MathExpression",
		[]interface{}{props},
		m,
	)
}

// Make a new Alarm for this metric.
//
// Combines both properties that may adjust the metric (aggregation) as well
// as alarm properties.
// Experimental.
func (m *jsiiProxy_MathExpression) CreateAlarm(scope awscdk.Construct, id *string, props *CreateAlarmOptions) Alarm {
	var returns Alarm

	_jsii_.Invoke(
		m,
		"createAlarm",
		[]interface{}{scope, id, props},
		&returns,
	)

	return returns
}

// Turn this metric object into an alarm configuration.
// Deprecated: use toMetricConfig()
func (m *jsiiProxy_MathExpression) ToAlarmConfig() *MetricAlarmConfig {
	var returns *MetricAlarmConfig

	_jsii_.Invoke(
		m,
		"toAlarmConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Turn this metric object into a graph configuration.
// Deprecated: use toMetricConfig()
func (m *jsiiProxy_MathExpression) ToGraphConfig() *MetricGraphConfig {
	var returns *MetricGraphConfig

	_jsii_.Invoke(
		m,
		"toGraphConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Inspect the details of the metric object.
// Experimental.
func (m *jsiiProxy_MathExpression) ToMetricConfig() *MetricConfig {
	var returns *MetricConfig

	_jsii_.Invoke(
		m,
		"toMetricConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (m *jsiiProxy_MathExpression) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		m,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return a copy of Metric with properties changed.
//
// All properties except namespace and metricName can be changed.
// Experimental.
func (m *jsiiProxy_MathExpression) With(props *MathExpressionOptions) MathExpression {
	var returns MathExpression

	_jsii_.Invoke(
		m,
		"with",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Configurable options for MathExpressions.
// Experimental.
type MathExpressionOptions struct {
	// Color for this metric when added to a Graph in a Dashboard.
	// Experimental.
	Color *string `json:"color"`
	// Label for this metric when added to a Graph in a Dashboard.
	// Experimental.
	Label *string `json:"label"`
	// The period over which the expression's statistics are applied.
	//
	// This period overrides all periods in the metrics used in this
	// math expression.
	// Experimental.
	Period awscdk.Duration `json:"period"`
}

// Properties for a MathExpression.
// Experimental.
type MathExpressionProps struct {
	// Color for this metric when added to a Graph in a Dashboard.
	// Experimental.
	Color *string `json:"color"`
	// Label for this metric when added to a Graph in a Dashboard.
	// Experimental.
	Label *string `json:"label"`
	// The period over which the expression's statistics are applied.
	//
	// This period overrides all periods in the metrics used in this
	// math expression.
	// Experimental.
	Period awscdk.Duration `json:"period"`
	// The expression defining the metric.
	// Experimental.
	Expression *string `json:"expression"`
	// The metrics used in the expression, in a map.
	//
	// The key is the identifier that represents the given metric in the
	// expression, and the value is the actual Metric object.
	// Experimental.
	UsingMetrics *map[string]IMetric `json:"usingMetrics"`
}

// A metric emitted by a service.
//
// The metric is a combination of a metric identifier (namespace, name and dimensions)
// and an aggregation function (statistic, period and unit).
//
// It also contains metadata which is used only in graphs, such as color and label.
// It makes sense to embed this in here, so that compound constructs can attach
// that metadata to metrics they expose.
//
// This class does not represent a resource, so hence is not a construct. Instead,
// Metric is an abstraction that makes it easy to specify metrics for use in both
// alarms and graphs.
// Experimental.
type Metric interface {
	IMetric
	Account() *string
	Color() *string
	Dimensions() *map[string]interface{}
	Label() *string
	MetricName() *string
	Namespace() *string
	Period() awscdk.Duration
	Region() *string
	Statistic() *string
	Unit() Unit
	AttachTo(scope constructs.IConstruct) Metric
	CreateAlarm(scope awscdk.Construct, id *string, props *CreateAlarmOptions) Alarm
	ToAlarmConfig() *MetricAlarmConfig
	ToGraphConfig() *MetricGraphConfig
	ToMetricConfig() *MetricConfig
	ToString() *string
	With(props *MetricOptions) Metric
}

// The jsii proxy struct for Metric
type jsiiProxy_Metric struct {
	jsiiProxy_IMetric
}

func (j *jsiiProxy_Metric) Account() *string {
	var returns *string
	_jsii_.Get(
		j,
		"account",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Color() *string {
	var returns *string
	_jsii_.Get(
		j,
		"color",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Dimensions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"dimensions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Label() *string {
	var returns *string
	_jsii_.Get(
		j,
		"label",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) MetricName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"metricName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Namespace() *string {
	var returns *string
	_jsii_.Get(
		j,
		"namespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Period() awscdk.Duration {
	var returns awscdk.Duration
	_jsii_.Get(
		j,
		"period",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Statistic() *string {
	var returns *string
	_jsii_.Get(
		j,
		"statistic",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Metric) Unit() Unit {
	var returns Unit
	_jsii_.Get(
		j,
		"unit",
		&returns,
	)
	return returns
}


// Experimental.
func NewMetric(props *MetricProps) Metric {
	_init_.Initialize()

	j := jsiiProxy_Metric{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Metric",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewMetric_Override(m Metric, props *MetricProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Metric",
		[]interface{}{props},
		m,
	)
}

// Grant permissions to the given identity to write metrics.
// Experimental.
func Metric_GrantPutMetricData(grantee awsiam.IGrantable) awsiam.Grant {
	_init_.Initialize()

	var returns awsiam.Grant

	_jsii_.StaticInvoke(
		"monocdk.aws_cloudwatch.Metric",
		"grantPutMetricData",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Attach the metric object to the given construct scope.
//
// Returns a Metric object that uses the account and region from the Stack
// the given construct is defined in. If the metric is subsequently used
// in a Dashboard or Alarm in a different Stack defined in a different
// account or region, the appropriate 'region' and 'account' fields
// will be added to it.
//
// If the scope we attach to is in an environment-agnostic stack,
// nothing is done and the same Metric object is returned.
// Experimental.
func (m *jsiiProxy_Metric) AttachTo(scope constructs.IConstruct) Metric {
	var returns Metric

	_jsii_.Invoke(
		m,
		"attachTo",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Make a new Alarm for this metric.
//
// Combines both properties that may adjust the metric (aggregation) as well
// as alarm properties.
// Experimental.
func (m *jsiiProxy_Metric) CreateAlarm(scope awscdk.Construct, id *string, props *CreateAlarmOptions) Alarm {
	var returns Alarm

	_jsii_.Invoke(
		m,
		"createAlarm",
		[]interface{}{scope, id, props},
		&returns,
	)

	return returns
}

// Turn this metric object into an alarm configuration.
// Deprecated: use toMetricConfig()
func (m *jsiiProxy_Metric) ToAlarmConfig() *MetricAlarmConfig {
	var returns *MetricAlarmConfig

	_jsii_.Invoke(
		m,
		"toAlarmConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Turn this metric object into a graph configuration.
// Deprecated: use toMetricConfig()
func (m *jsiiProxy_Metric) ToGraphConfig() *MetricGraphConfig {
	var returns *MetricGraphConfig

	_jsii_.Invoke(
		m,
		"toGraphConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Inspect the details of the metric object.
// Experimental.
func (m *jsiiProxy_Metric) ToMetricConfig() *MetricConfig {
	var returns *MetricConfig

	_jsii_.Invoke(
		m,
		"toMetricConfig",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (m *jsiiProxy_Metric) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		m,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return a copy of Metric `with` properties changed.
//
// All properties except namespace and metricName can be changed.
// Experimental.
func (m *jsiiProxy_Metric) With(props *MetricOptions) Metric {
	var returns Metric

	_jsii_.Invoke(
		m,
		"with",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Properties used to construct the Metric identifying part of an Alarm.
// Deprecated: Replaced by MetricConfig
type MetricAlarmConfig struct {
	// Name of the metric.
	// Deprecated: Replaced by MetricConfig
	MetricName *string `json:"metricName"`
	// Namespace of the metric.
	// Deprecated: Replaced by MetricConfig
	Namespace *string `json:"namespace"`
	// How many seconds to aggregate over.
	// Deprecated: Replaced by MetricConfig
	Period *float64 `json:"period"`
	// The dimensions to apply to the alarm.
	// Deprecated: Replaced by MetricConfig
	Dimensions *[]*Dimension `json:"dimensions"`
	// Percentile aggregation function to use.
	// Deprecated: Replaced by MetricConfig
	ExtendedStatistic *string `json:"extendedStatistic"`
	// Simple aggregation function to use.
	// Deprecated: Replaced by MetricConfig
	Statistic Statistic `json:"statistic"`
	// The unit of the alarm.
	// Deprecated: Replaced by MetricConfig
	Unit Unit `json:"unit"`
}

// Properties of a rendered metric.
// Experimental.
type MetricConfig struct {
	// In case the metric is a math expression, the details of the math expression.
	// Experimental.
	MathExpression *MetricExpressionConfig `json:"mathExpression"`
	// In case the metric represents a query, the details of the query.
	// Experimental.
	MetricStat *MetricStatConfig `json:"metricStat"`
	// Additional properties which will be rendered if the metric is used in a dashboard.
	//
	// Examples are 'label' and 'color', but any key in here will be
	// added to dashboard graphs.
	// Experimental.
	RenderingProperties *map[string]interface{} `json:"renderingProperties"`
}

// Properties for a concrete metric.
// Experimental.
type MetricExpressionConfig struct {
	// Math expression for the metric.
	// Experimental.
	Expression *string `json:"expression"`
	// How many seconds to aggregate over.
	// Experimental.
	Period *float64 `json:"period"`
	// Metrics used in the math expression.
	// Experimental.
	UsingMetrics *map[string]IMetric `json:"usingMetrics"`
}

// Properties used to construct the Metric identifying part of a Graph.
// Deprecated: Replaced by MetricConfig
type MetricGraphConfig struct {
	// Name of the metric.
	// Deprecated: Replaced by MetricConfig
	MetricName *string `json:"metricName"`
	// Namespace of the metric.
	// Deprecated: Replaced by MetricConfig
	Namespace *string `json:"namespace"`
	// How many seconds to aggregate over.
	// Deprecated: Use `period` in `renderingProperties`
	Period *float64 `json:"period"`
	// Rendering properties override yAxis parameter of the widget object.
	// Deprecated: Replaced by MetricConfig
	RenderingProperties *MetricRenderingProperties `json:"renderingProperties"`
	// Color for the graph line.
	// Deprecated: Use `color` in `renderingProperties`
	Color *string `json:"color"`
	// The dimensions to apply to the alarm.
	// Deprecated: Replaced by MetricConfig
	Dimensions *[]*Dimension `json:"dimensions"`
	// Label for the metric.
	// Deprecated: Use `label` in `renderingProperties`
	Label *string `json:"label"`
	// Aggregation function to use (can be either simple or a percentile).
	// Deprecated: Use `stat` in `renderingProperties`
	Statistic *string `json:"statistic"`
	// The unit of the alarm.
	// Deprecated: not used in dashboard widgets
	Unit Unit `json:"unit"`
}

// Properties of a metric that can be changed.
// Experimental.
type MetricOptions struct {
	// Account which this metric comes from.
	// Experimental.
	Account *string `json:"account"`
	// The hex color code, prefixed with '#' (e.g. '#00ff00'), to use when this metric is rendered on a graph. The `Color` class has a set of standard colors that can be used here.
	// Experimental.
	Color *string `json:"color"`
	// Dimensions of the metric.
	// Deprecated: Use 'dimensionsMap' instead.
	Dimensions *map[string]interface{} `json:"dimensions"`
	// Dimensions of the metric.
	// Experimental.
	DimensionsMap *map[string]*string `json:"dimensionsMap"`
	// Label for this metric when added to a Graph in a Dashboard.
	// Experimental.
	Label *string `json:"label"`
	// The period over which the specified statistic is applied.
	// Experimental.
	Period awscdk.Duration `json:"period"`
	// Region which this metric comes from.
	// Experimental.
	Region *string `json:"region"`
	// What function to use for aggregating.
	//
	// Can be one of the following:
	//
	// - "Minimum" | "min"
	// - "Maximum" | "max"
	// - "Average" | "avg"
	// - "Sum" | "sum"
	// - "SampleCount | "n"
	// - "pNN.NN"
	// Experimental.
	Statistic *string `json:"statistic"`
	// Unit used to filter the metric stream.
	//
	// Only refer to datums emitted to the metric stream with the given unit and
	// ignore all others. Only useful when datums are being emitted to the same
	// metric stream under different units.
	//
	// The default is to use all matric datums in the stream, regardless of unit,
	// which is recommended in nearly all cases.
	//
	// CloudWatch does not honor this property for graphs.
	// Experimental.
	Unit Unit `json:"unit"`
}

// Properties for a metric.
// Experimental.
type MetricProps struct {
	// Account which this metric comes from.
	// Experimental.
	Account *string `json:"account"`
	// The hex color code, prefixed with '#' (e.g. '#00ff00'), to use when this metric is rendered on a graph. The `Color` class has a set of standard colors that can be used here.
	// Experimental.
	Color *string `json:"color"`
	// Dimensions of the metric.
	// Deprecated: Use 'dimensionsMap' instead.
	Dimensions *map[string]interface{} `json:"dimensions"`
	// Dimensions of the metric.
	// Experimental.
	DimensionsMap *map[string]*string `json:"dimensionsMap"`
	// Label for this metric when added to a Graph in a Dashboard.
	// Experimental.
	Label *string `json:"label"`
	// The period over which the specified statistic is applied.
	// Experimental.
	Period awscdk.Duration `json:"period"`
	// Region which this metric comes from.
	// Experimental.
	Region *string `json:"region"`
	// What function to use for aggregating.
	//
	// Can be one of the following:
	//
	// - "Minimum" | "min"
	// - "Maximum" | "max"
	// - "Average" | "avg"
	// - "Sum" | "sum"
	// - "SampleCount | "n"
	// - "pNN.NN"
	// Experimental.
	Statistic *string `json:"statistic"`
	// Unit used to filter the metric stream.
	//
	// Only refer to datums emitted to the metric stream with the given unit and
	// ignore all others. Only useful when datums are being emitted to the same
	// metric stream under different units.
	//
	// The default is to use all matric datums in the stream, regardless of unit,
	// which is recommended in nearly all cases.
	//
	// CloudWatch does not honor this property for graphs.
	// Experimental.
	Unit Unit `json:"unit"`
	// Name of the metric.
	// Experimental.
	MetricName *string `json:"metricName"`
	// Namespace of the metric.
	// Experimental.
	Namespace *string `json:"namespace"`
}

// Custom rendering properties that override the default rendering properties specified in the yAxis parameter of the widget object.
// Deprecated: Replaced by MetricConfig.
type MetricRenderingProperties struct {
	// How many seconds to aggregate over.
	// Deprecated: Replaced by MetricConfig.
	Period *float64 `json:"period"`
	// The hex color code, prefixed with '#' (e.g. '#00ff00'), to use when this metric is rendered on a graph. The `Color` class has a set of standard colors that can be used here.
	// Deprecated: Replaced by MetricConfig.
	Color *string `json:"color"`
	// Label for the metric.
	// Deprecated: Replaced by MetricConfig.
	Label *string `json:"label"`
	// Aggregation function to use (can be either simple or a percentile).
	// Deprecated: Replaced by MetricConfig.
	Stat *string `json:"stat"`
}

// Properties for a concrete metric.
//
// NOTE: `unit` is no longer on this object since it is only used for `Alarms`, and doesn't mean what one
// would expect it to mean there anyway. It is most likely to be misused.
// Experimental.
type MetricStatConfig struct {
	// Name of the metric.
	// Experimental.
	MetricName *string `json:"metricName"`
	// Namespace of the metric.
	// Experimental.
	Namespace *string `json:"namespace"`
	// How many seconds to aggregate over.
	// Experimental.
	Period awscdk.Duration `json:"period"`
	// Aggregation function to use (can be either simple or a percentile).
	// Experimental.
	Statistic *string `json:"statistic"`
	// Account which this metric comes from.
	// Experimental.
	Account *string `json:"account"`
	// The dimensions to apply to the alarm.
	// Experimental.
	Dimensions *[]*Dimension `json:"dimensions"`
	// Region which this metric comes from.
	// Experimental.
	Region *string `json:"region"`
	// Unit used to filter the metric stream.
	//
	// Only refer to datums emitted to the metric stream with the given unit and
	// ignore all others. Only useful when datums are being emitted to the same
	// metric stream under different units.
	//
	// This field has been renamed from plain `unit` to clearly communicate
	// its purpose.
	// Experimental.
	UnitFilter Unit `json:"unitFilter"`
}

// Basic properties for widgets that display metrics.
// Experimental.
type MetricWidgetProps struct {
	// Height of the widget.
	// Experimental.
	Height *float64 `json:"height"`
	// The region the metrics of this graph should be taken from.
	// Experimental.
	Region *string `json:"region"`
	// Title for the graph.
	// Experimental.
	Title *string `json:"title"`
	// Width of the widget, in a grid of 24 units wide.
	// Experimental.
	Width *float64 `json:"width"`
}

// Specify the period for graphs when the CloudWatch dashboard loads.
// Experimental.
type PeriodOverride string

const (
	PeriodOverride_AUTO PeriodOverride = "AUTO"
	PeriodOverride_INHERIT PeriodOverride = "INHERIT"
)

// A widget that contains other widgets in a horizontal row.
//
// Widgets will be laid out next to each other
// Experimental.
type Row interface {
	IWidget
	Height() *float64
	Width() *float64
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for Row
type jsiiProxy_Row struct {
	jsiiProxy_IWidget
}

func (j *jsiiProxy_Row) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Row) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}


// Experimental.
func NewRow(widgets ...IWidget) Row {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range widgets {
		args = append(args, a)
	}

	j := jsiiProxy_Row{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Row",
		args,
		&j,
	)

	return &j
}

// Experimental.
func NewRow_Override(r Row, widgets ...IWidget) {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range widgets {
		args = append(args, a)
	}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Row",
		args,
		r,
	)
}

// Place the widget at a given position.
// Experimental.
func (r *jsiiProxy_Row) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		r,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (r *jsiiProxy_Row) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		r,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Fill shading options that will be used with an annotation.
// Experimental.
type Shading string

const (
	Shading_NONE Shading = "NONE"
	Shading_ABOVE Shading = "ABOVE"
	Shading_BELOW Shading = "BELOW"
)

// A dashboard widget that displays the most recent value for every metric.
// Experimental.
type SingleValueWidget interface {
	ConcreteWidget
	Height() *float64
	Width() *float64
	X() *float64
	SetX(val *float64)
	Y() *float64
	SetY(val *float64)
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for SingleValueWidget
type jsiiProxy_SingleValueWidget struct {
	jsiiProxy_ConcreteWidget
}

func (j *jsiiProxy_SingleValueWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingleValueWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingleValueWidget) X() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"x",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingleValueWidget) Y() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"y",
		&returns,
	)
	return returns
}


// Experimental.
func NewSingleValueWidget(props *SingleValueWidgetProps) SingleValueWidget {
	_init_.Initialize()

	j := jsiiProxy_SingleValueWidget{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.SingleValueWidget",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewSingleValueWidget_Override(s SingleValueWidget, props *SingleValueWidgetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.SingleValueWidget",
		[]interface{}{props},
		s,
	)
}

func (j *jsiiProxy_SingleValueWidget) SetX(val *float64) {
	_jsii_.Set(
		j,
		"x",
		val,
	)
}

func (j *jsiiProxy_SingleValueWidget) SetY(val *float64) {
	_jsii_.Set(
		j,
		"y",
		val,
	)
}

// Place the widget at a given position.
// Experimental.
func (s *jsiiProxy_SingleValueWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		s,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (s *jsiiProxy_SingleValueWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		s,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a SingleValueWidget.
// Experimental.
type SingleValueWidgetProps struct {
	// Height of the widget.
	// Experimental.
	Height *float64 `json:"height"`
	// The region the metrics of this graph should be taken from.
	// Experimental.
	Region *string `json:"region"`
	// Title for the graph.
	// Experimental.
	Title *string `json:"title"`
	// Width of the widget, in a grid of 24 units wide.
	// Experimental.
	Width *float64 `json:"width"`
	// Metrics to display.
	// Experimental.
	Metrics *[]IMetric `json:"metrics"`
	// Whether to show as many digits as can fit, before rounding.
	// Experimental.
	FullPrecision *bool `json:"fullPrecision"`
	// Whether to show the value from the entire time range.
	// Experimental.
	SetPeriodToTimeRange *bool `json:"setPeriodToTimeRange"`
}

// A widget that doesn't display anything but takes up space.
// Experimental.
type Spacer interface {
	IWidget
	Height() *float64
	Width() *float64
	Position(_x *float64, _y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for Spacer
type jsiiProxy_Spacer struct {
	jsiiProxy_IWidget
}

func (j *jsiiProxy_Spacer) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Spacer) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}


// Experimental.
func NewSpacer(props *SpacerProps) Spacer {
	_init_.Initialize()

	j := jsiiProxy_Spacer{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Spacer",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewSpacer_Override(s Spacer, props *SpacerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.Spacer",
		[]interface{}{props},
		s,
	)
}

// Place the widget at a given position.
// Experimental.
func (s *jsiiProxy_Spacer) Position(_x *float64, _y *float64) {
	_jsii_.InvokeVoid(
		s,
		"position",
		[]interface{}{_x, _y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (s *jsiiProxy_Spacer) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		s,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Props of the spacer.
// Experimental.
type SpacerProps struct {
	// Height of the spacer.
	// Experimental.
	Height *float64 `json:"height"`
	// Width of the spacer.
	// Experimental.
	Width *float64 `json:"width"`
}

// Statistic to use over the aggregation period.
// Experimental.
type Statistic string

const (
	Statistic_SAMPLE_COUNT Statistic = "SAMPLE_COUNT"
	Statistic_AVERAGE Statistic = "AVERAGE"
	Statistic_SUM Statistic = "SUM"
	Statistic_MINIMUM Statistic = "MINIMUM"
	Statistic_MAXIMUM Statistic = "MAXIMUM"
)

// A dashboard widget that displays MarkDown.
// Experimental.
type TextWidget interface {
	ConcreteWidget
	Height() *float64
	Width() *float64
	X() *float64
	SetX(val *float64)
	Y() *float64
	SetY(val *float64)
	Position(x *float64, y *float64)
	ToJson() *[]interface{}
}

// The jsii proxy struct for TextWidget
type jsiiProxy_TextWidget struct {
	jsiiProxy_ConcreteWidget
}

func (j *jsiiProxy_TextWidget) Height() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"height",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TextWidget) Width() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"width",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TextWidget) X() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"x",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TextWidget) Y() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"y",
		&returns,
	)
	return returns
}


// Experimental.
func NewTextWidget(props *TextWidgetProps) TextWidget {
	_init_.Initialize()

	j := jsiiProxy_TextWidget{}

	_jsii_.Create(
		"monocdk.aws_cloudwatch.TextWidget",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewTextWidget_Override(t TextWidget, props *TextWidgetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cloudwatch.TextWidget",
		[]interface{}{props},
		t,
	)
}

func (j *jsiiProxy_TextWidget) SetX(val *float64) {
	_jsii_.Set(
		j,
		"x",
		val,
	)
}

func (j *jsiiProxy_TextWidget) SetY(val *float64) {
	_jsii_.Set(
		j,
		"y",
		val,
	)
}

// Place the widget at a given position.
// Experimental.
func (t *jsiiProxy_TextWidget) Position(x *float64, y *float64) {
	_jsii_.InvokeVoid(
		t,
		"position",
		[]interface{}{x, y},
	)
}

// Return the widget JSON for use in the dashboard.
// Experimental.
func (t *jsiiProxy_TextWidget) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		t,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a Text widget.
// Experimental.
type TextWidgetProps struct {
	// The text to display, in MarkDown format.
	// Experimental.
	Markdown *string `json:"markdown"`
	// Height of the widget.
	// Experimental.
	Height *float64 `json:"height"`
	// Width of the widget, in a grid of 24 units wide.
	// Experimental.
	Width *float64 `json:"width"`
}

// Specify how missing data points are treated during alarm evaluation.
// Experimental.
type TreatMissingData string

const (
	TreatMissingData_BREACHING TreatMissingData = "BREACHING"
	TreatMissingData_NOT_BREACHING TreatMissingData = "NOT_BREACHING"
	TreatMissingData_IGNORE TreatMissingData = "IGNORE"
	TreatMissingData_MISSING TreatMissingData = "MISSING"
)

// Unit for metric.
// Experimental.
type Unit string

const (
	Unit_SECONDS Unit = "SECONDS"
	Unit_MICROSECONDS Unit = "MICROSECONDS"
	Unit_MILLISECONDS Unit = "MILLISECONDS"
	Unit_BYTES Unit = "BYTES"
	Unit_KILOBYTES Unit = "KILOBYTES"
	Unit_MEGABYTES Unit = "MEGABYTES"
	Unit_GIGABYTES Unit = "GIGABYTES"
	Unit_TERABYTES Unit = "TERABYTES"
	Unit_BITS Unit = "BITS"
	Unit_KILOBITS Unit = "KILOBITS"
	Unit_MEGABITS Unit = "MEGABITS"
	Unit_GIGABITS Unit = "GIGABITS"
	Unit_TERABITS Unit = "TERABITS"
	Unit_PERCENT Unit = "PERCENT"
	Unit_COUNT Unit = "COUNT"
	Unit_BYTES_PER_SECOND Unit = "BYTES_PER_SECOND"
	Unit_KILOBYTES_PER_SECOND Unit = "KILOBYTES_PER_SECOND"
	Unit_MEGABYTES_PER_SECOND Unit = "MEGABYTES_PER_SECOND"
	Unit_GIGABYTES_PER_SECOND Unit = "GIGABYTES_PER_SECOND"
	Unit_TERABYTES_PER_SECOND Unit = "TERABYTES_PER_SECOND"
	Unit_BITS_PER_SECOND Unit = "BITS_PER_SECOND"
	Unit_KILOBITS_PER_SECOND Unit = "KILOBITS_PER_SECOND"
	Unit_MEGABITS_PER_SECOND Unit = "MEGABITS_PER_SECOND"
	Unit_GIGABITS_PER_SECOND Unit = "GIGABITS_PER_SECOND"
	Unit_TERABITS_PER_SECOND Unit = "TERABITS_PER_SECOND"
	Unit_COUNT_PER_SECOND Unit = "COUNT_PER_SECOND"
	Unit_NONE Unit = "NONE"
)

// Properties for a Y-Axis.
// Experimental.
type YAxisProps struct {
	// The label.
	// Experimental.
	Label *string `json:"label"`
	// The max value.
	// Experimental.
	Max *float64 `json:"max"`
	// The min value.
	// Experimental.
	Min *float64 `json:"min"`
	// Whether to show units.
	// Experimental.
	ShowUnits *bool `json:"showUnits"`
}


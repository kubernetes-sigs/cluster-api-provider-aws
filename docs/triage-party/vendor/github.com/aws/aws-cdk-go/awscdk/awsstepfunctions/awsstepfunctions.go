package awsstepfunctions

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/awsstepfunctions/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// Define a new Step Functions Activity.
// Experimental.
type Activity interface {
	awscdk.Resource
	IActivity
	ActivityArn() *string
	ActivityName() *string
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(identity awsiam.IGrantable, actions ...*string) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHeartbeatTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRunTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricScheduled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricScheduleTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Activity
type jsiiProxy_Activity struct {
	internal.Type__awscdkResource
	jsiiProxy_IActivity
}

func (j *jsiiProxy_Activity) ActivityArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"activityArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Activity) ActivityName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"activityName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Activity) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Activity) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Activity) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Activity) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewActivity(scope constructs.Construct, id *string, props *ActivityProps) Activity {
	_init_.Initialize()

	j := jsiiProxy_Activity{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Activity",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewActivity_Override(a Activity, scope constructs.Construct, id *string, props *ActivityProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Activity",
		[]interface{}{scope, id, props},
		a,
	)
}

// Construct an Activity from an existing Activity ARN.
// Experimental.
func Activity_FromActivityArn(scope constructs.Construct, id *string, activityArn *string) IActivity {
	_init_.Initialize()

	var returns IActivity

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Activity",
		"fromActivityArn",
		[]interface{}{scope, id, activityArn},
		&returns,
	)

	return returns
}

// Construct an Activity from an existing Activity Name.
// Experimental.
func Activity_FromActivityName(scope constructs.Construct, id *string, activityName *string) IActivity {
	_init_.Initialize()

	var returns IActivity

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Activity",
		"fromActivityName",
		[]interface{}{scope, id, activityName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Activity_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Activity",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Activity_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Activity",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
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
func (a *jsiiProxy_Activity) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_Activity) GeneratePhysicalName() *string {
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
func (a *jsiiProxy_Activity) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (a *jsiiProxy_Activity) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity permissions on this Activity.
// Experimental.
func (a *jsiiProxy_Activity) Grant(identity awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{identity}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Return the given named metric for this Activity.
// Experimental.
func (a *jsiiProxy_Activity) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity fails.
// Experimental.
func (a *jsiiProxy_Activity) MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times the heartbeat times out for this activity.
// Experimental.
func (a *jsiiProxy_Activity) MetricHeartbeatTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHeartbeatTimedOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, between the time the activity starts and the time it closes.
// Experimental.
func (a *jsiiProxy_Activity) MetricRunTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricRunTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity is scheduled.
// Experimental.
func (a *jsiiProxy_Activity) MetricScheduled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricScheduled",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, for which the activity stays in the schedule state.
// Experimental.
func (a *jsiiProxy_Activity) MetricScheduleTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricScheduleTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity is started.
// Experimental.
func (a *jsiiProxy_Activity) MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricStarted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity succeeds.
// Experimental.
func (a *jsiiProxy_Activity) MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricSucceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, between the time the activity is scheduled and the time it closes.
// Experimental.
func (a *jsiiProxy_Activity) MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity times out.
// Experimental.
func (a *jsiiProxy_Activity) MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTimedOut",
		[]interface{}{props},
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
func (a *jsiiProxy_Activity) OnPrepare() {
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
func (a *jsiiProxy_Activity) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_Activity) OnValidate() *[]*string {
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
func (a *jsiiProxy_Activity) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_Activity) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_Activity) ToString() *string {
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
func (a *jsiiProxy_Activity) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining a new Step Functions Activity.
// Experimental.
type ActivityProps struct {
	// The name for this activity.
	// Experimental.
	ActivityName *string `json:"activityName"`
}

// Options for selecting the choice paths.
// Experimental.
type AfterwardsOptions struct {
	// Whether to include error handling states.
	//
	// If this is true, all states which are error handlers (added through 'onError')
	// and states reachable via error handlers will be included as well.
	// Experimental.
	IncludeErrorHandlers *bool `json:"includeErrorHandlers"`
	// Whether to include the default/otherwise transition for the current Choice state.
	//
	// If this is true and the current Choice does not have a default outgoing
	// transition, one will be added included when .next() is called on the chain.
	// Experimental.
	IncludeOtherwise *bool `json:"includeOtherwise"`
}

// Error handler details.
// Experimental.
type CatchProps struct {
	// Errors to recover from by going to the given state.
	//
	// A list of error strings to retry, which can be either predefined errors
	// (for example Errors.NoChoiceMatched) or a self-defined error.
	// Experimental.
	Errors *[]*string `json:"errors"`
	// JSONPath expression to indicate where to inject the error data.
	//
	// May also be the special value DISCARD, which will cause the error
	// data to be discarded.
	// Experimental.
	ResultPath *string `json:"resultPath"`
}

// A CloudFormation `AWS::StepFunctions::Activity`.
type CfnActivity interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
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

// The jsii proxy struct for CfnActivity
type jsiiProxy_CfnActivity struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnActivity) AttrName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnActivity) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::StepFunctions::Activity`.
func NewCfnActivity(scope awscdk.Construct, id *string, props *CfnActivityProps) CfnActivity {
	_init_.Initialize()

	j := jsiiProxy_CfnActivity{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.CfnActivity",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::StepFunctions::Activity`.
func NewCfnActivity_Override(c CfnActivity, scope awscdk.Construct, id *string, props *CfnActivityProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.CfnActivity",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnActivity) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
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
func CfnActivity_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CfnActivity",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnActivity_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CfnActivity",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnActivity_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CfnActivity",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnActivity_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.CfnActivity",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnActivity) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnActivity) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnActivity) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnActivity) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnActivity) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnActivity) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnActivity) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnActivity) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnActivity) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnActivity) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnActivity) OnPrepare() {
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
func (c *jsiiProxy_CfnActivity) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnActivity) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnActivity) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnActivity) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnActivity) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnActivity) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnActivity) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnActivity) ToString() *string {
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
func (c *jsiiProxy_CfnActivity) Validate() *[]*string {
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
func (c *jsiiProxy_CfnActivity) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnActivity_TagsEntryProperty struct {
	// `CfnActivity.TagsEntryProperty.Key`.
	Key *string `json:"key"`
	// `CfnActivity.TagsEntryProperty.Value`.
	Value *string `json:"value"`
}

// Properties for defining a `AWS::StepFunctions::Activity`.
type CfnActivityProps struct {
	// `AWS::StepFunctions::Activity.Name`.
	Name *string `json:"name"`
	// `AWS::StepFunctions::Activity.Tags`.
	Tags *[]*CfnActivity_TagsEntryProperty `json:"tags"`
}

// A CloudFormation `AWS::StepFunctions::StateMachine`.
type CfnStateMachine interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Definition() interface{}
	SetDefinition(val interface{})
	DefinitionS3Location() interface{}
	SetDefinitionS3Location(val interface{})
	DefinitionString() *string
	SetDefinitionString(val *string)
	DefinitionSubstitutions() interface{}
	SetDefinitionSubstitutions(val interface{})
	LoggingConfiguration() interface{}
	SetLoggingConfiguration(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RoleArn() *string
	SetRoleArn(val *string)
	Stack() awscdk.Stack
	StateMachineName() *string
	SetStateMachineName(val *string)
	StateMachineType() *string
	SetStateMachineType(val *string)
	Tags() awscdk.TagManager
	TracingConfiguration() interface{}
	SetTracingConfiguration(val interface{})
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

// The jsii proxy struct for CfnStateMachine
type jsiiProxy_CfnStateMachine struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnStateMachine) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) AttrName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) Definition() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"definition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) DefinitionS3Location() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"definitionS3Location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) DefinitionString() *string {
	var returns *string
	_jsii_.Get(
		j,
		"definitionString",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) DefinitionSubstitutions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"definitionSubstitutions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) LoggingConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"loggingConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) StateMachineName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateMachineName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) StateMachineType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateMachineType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) TracingConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"tracingConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStateMachine) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::StepFunctions::StateMachine`.
func NewCfnStateMachine(scope awscdk.Construct, id *string, props *CfnStateMachineProps) CfnStateMachine {
	_init_.Initialize()

	j := jsiiProxy_CfnStateMachine{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.CfnStateMachine",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::StepFunctions::StateMachine`.
func NewCfnStateMachine_Override(c CfnStateMachine, scope awscdk.Construct, id *string, props *CfnStateMachineProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.CfnStateMachine",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetDefinition(val interface{}) {
	_jsii_.Set(
		j,
		"definition",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetDefinitionS3Location(val interface{}) {
	_jsii_.Set(
		j,
		"definitionS3Location",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetDefinitionString(val *string) {
	_jsii_.Set(
		j,
		"definitionString",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetDefinitionSubstitutions(val interface{}) {
	_jsii_.Set(
		j,
		"definitionSubstitutions",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetLoggingConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"loggingConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetRoleArn(val *string) {
	_jsii_.Set(
		j,
		"roleArn",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetStateMachineName(val *string) {
	_jsii_.Set(
		j,
		"stateMachineName",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetStateMachineType(val *string) {
	_jsii_.Set(
		j,
		"stateMachineType",
		val,
	)
}

func (j *jsiiProxy_CfnStateMachine) SetTracingConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"tracingConfiguration",
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
func CfnStateMachine_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CfnStateMachine",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnStateMachine_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CfnStateMachine",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnStateMachine_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CfnStateMachine",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStateMachine_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.CfnStateMachine",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnStateMachine) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnStateMachine) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnStateMachine) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnStateMachine) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnStateMachine) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnStateMachine) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnStateMachine) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnStateMachine) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnStateMachine) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnStateMachine) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnStateMachine) OnPrepare() {
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
func (c *jsiiProxy_CfnStateMachine) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStateMachine) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnStateMachine) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnStateMachine) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnStateMachine) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnStateMachine) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnStateMachine) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStateMachine) ToString() *string {
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
func (c *jsiiProxy_CfnStateMachine) Validate() *[]*string {
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
func (c *jsiiProxy_CfnStateMachine) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnStateMachine_CloudWatchLogsLogGroupProperty struct {
	// `CfnStateMachine.CloudWatchLogsLogGroupProperty.LogGroupArn`.
	LogGroupArn *string `json:"logGroupArn"`
}

type CfnStateMachine_DefinitionProperty struct {
}

type CfnStateMachine_LogDestinationProperty struct {
	// `CfnStateMachine.LogDestinationProperty.CloudWatchLogsLogGroup`.
	CloudWatchLogsLogGroup interface{} `json:"cloudWatchLogsLogGroup"`
}

type CfnStateMachine_LoggingConfigurationProperty struct {
	// `CfnStateMachine.LoggingConfigurationProperty.Destinations`.
	Destinations interface{} `json:"destinations"`
	// `CfnStateMachine.LoggingConfigurationProperty.IncludeExecutionData`.
	IncludeExecutionData interface{} `json:"includeExecutionData"`
	// `CfnStateMachine.LoggingConfigurationProperty.Level`.
	Level *string `json:"level"`
}

type CfnStateMachine_S3LocationProperty struct {
	// `CfnStateMachine.S3LocationProperty.Bucket`.
	Bucket *string `json:"bucket"`
	// `CfnStateMachine.S3LocationProperty.Key`.
	Key *string `json:"key"`
	// `CfnStateMachine.S3LocationProperty.Version`.
	Version *string `json:"version"`
}

type CfnStateMachine_TagsEntryProperty struct {
	// `CfnStateMachine.TagsEntryProperty.Key`.
	Key *string `json:"key"`
	// `CfnStateMachine.TagsEntryProperty.Value`.
	Value *string `json:"value"`
}

type CfnStateMachine_TracingConfigurationProperty struct {
	// `CfnStateMachine.TracingConfigurationProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
}

// Properties for defining a `AWS::StepFunctions::StateMachine`.
type CfnStateMachineProps struct {
	// `AWS::StepFunctions::StateMachine.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `AWS::StepFunctions::StateMachine.Definition`.
	Definition interface{} `json:"definition"`
	// `AWS::StepFunctions::StateMachine.DefinitionS3Location`.
	DefinitionS3Location interface{} `json:"definitionS3Location"`
	// `AWS::StepFunctions::StateMachine.DefinitionString`.
	DefinitionString *string `json:"definitionString"`
	// `AWS::StepFunctions::StateMachine.DefinitionSubstitutions`.
	DefinitionSubstitutions interface{} `json:"definitionSubstitutions"`
	// `AWS::StepFunctions::StateMachine.LoggingConfiguration`.
	LoggingConfiguration interface{} `json:"loggingConfiguration"`
	// `AWS::StepFunctions::StateMachine.StateMachineName`.
	StateMachineName *string `json:"stateMachineName"`
	// `AWS::StepFunctions::StateMachine.StateMachineType`.
	StateMachineType *string `json:"stateMachineType"`
	// `AWS::StepFunctions::StateMachine.Tags`.
	Tags *[]*CfnStateMachine_TagsEntryProperty `json:"tags"`
	// `AWS::StepFunctions::StateMachine.TracingConfiguration`.
	TracingConfiguration interface{} `json:"tracingConfiguration"`
}

// A collection of states to chain onto.
//
// A Chain has a start and zero or more chainable ends. If there are
// zero ends, calling next() on the Chain will fail.
// Experimental.
type Chain interface {
	IChainable
	EndStates() *[]INextable
	Id() *string
	StartState() State
	Next(next IChainable) Chain
	ToSingleState(id *string, props *ParallelProps) Parallel
}

// The jsii proxy struct for Chain
type jsiiProxy_Chain struct {
	jsiiProxy_IChainable
}

func (j *jsiiProxy_Chain) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Chain) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Chain) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}


// Make a Chain with specific start and end states, and a last-added Chainable.
// Experimental.
func Chain_Custom(startState State, endStates *[]INextable, lastAdded IChainable) Chain {
	_init_.Initialize()

	var returns Chain

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Chain",
		"custom",
		[]interface{}{startState, endStates, lastAdded},
		&returns,
	)

	return returns
}

// Make a Chain with the start from one chain and the ends from another.
// Experimental.
func Chain_Sequence(start IChainable, next IChainable) Chain {
	_init_.Initialize()

	var returns Chain

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Chain",
		"sequence",
		[]interface{}{start, next},
		&returns,
	)

	return returns
}

// Begin a new Chain from one chainable.
// Experimental.
func Chain_Start(state IChainable) Chain {
	_init_.Initialize()

	var returns Chain

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Chain",
		"start",
		[]interface{}{state},
		&returns,
	)

	return returns
}

// Continue normal execution with the given state.
// Experimental.
func (c *jsiiProxy_Chain) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		c,
		"next",
		[]interface{}{next},
		&returns,
	)

	return returns
}

// Return a single state that encompasses all states in the chain.
//
// This can be used to add error handling to a sequence of states.
//
// Be aware that this changes the result of the inner state machine
// to be an array with the result of the state machine in it. Adjust
// your paths accordingly. For example, change 'outputPath' to
// '$[0]'.
// Experimental.
func (c *jsiiProxy_Chain) ToSingleState(id *string, props *ParallelProps) Parallel {
	var returns Parallel

	_jsii_.Invoke(
		c,
		"toSingleState",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Define a Choice in the state machine.
//
// A choice state can be used to make decisions based on the execution
// state.
// Experimental.
type Choice interface {
	State
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	Afterwards(options *AfterwardsOptions) Chain
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Otherwise(def IChainable) Choice
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	When(condition Condition, next IChainable) Choice
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Choice
type jsiiProxy_Choice struct {
	jsiiProxy_State
}

func (j *jsiiProxy_Choice) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Choice) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewChoice(scope constructs.Construct, id *string, props *ChoiceProps) Choice {
	_init_.Initialize()

	j := jsiiProxy_Choice{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Choice",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewChoice_Override(c Choice, scope constructs.Construct, id *string, props *ChoiceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Choice",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_Choice) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Choice) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func Choice_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Choice",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func Choice_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Choice",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func Choice_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Choice",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Choice_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Choice",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func Choice_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Choice",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (c *jsiiProxy_Choice) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a choice branch to this state.
// Experimental.
func (c *jsiiProxy_Choice) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		c,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (c *jsiiProxy_Choice) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (c *jsiiProxy_Choice) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		c,
		"addPrefix",
		[]interface{}{x},
	)
}

// Return a Chain that contains all reachable end states from this Choice.
//
// Use this to combine all possible choice paths back.
// Experimental.
func (c *jsiiProxy_Choice) Afterwards(options *AfterwardsOptions) Chain {
	var returns Chain

	_jsii_.Invoke(
		c,
		"afterwards",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (c *jsiiProxy_Choice) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (c *jsiiProxy_Choice) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		c,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (c *jsiiProxy_Choice) MakeNext(next State) {
	_jsii_.InvokeVoid(
		c,
		"makeNext",
		[]interface{}{next},
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
func (c *jsiiProxy_Choice) OnPrepare() {
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
func (c *jsiiProxy_Choice) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_Choice) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// If none of the given conditions match, continue execution with the given state.
//
// If no conditions match and no otherwise() has been given, an execution
// error will be raised.
// Experimental.
func (c *jsiiProxy_Choice) Otherwise(def IChainable) Choice {
	var returns Choice

	_jsii_.Invoke(
		c,
		"otherwise",
		[]interface{}{def},
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
func (c *jsiiProxy_Choice) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (c *jsiiProxy_Choice) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (c *jsiiProxy_Choice) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (c *jsiiProxy_Choice) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (c *jsiiProxy_Choice) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (c *jsiiProxy_Choice) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (c *jsiiProxy_Choice) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (c *jsiiProxy_Choice) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderRetryCatch",
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
func (c *jsiiProxy_Choice) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (c *jsiiProxy_Choice) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_Choice) ToString() *string {
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
func (c *jsiiProxy_Choice) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// If the given condition matches, continue execution with the given state.
// Experimental.
func (c *jsiiProxy_Choice) When(condition Condition, next IChainable) Choice {
	var returns Choice

	_jsii_.Invoke(
		c,
		"when",
		[]interface{}{condition, next},
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (c *jsiiProxy_Choice) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a Choice state.
// Experimental.
type ChoiceProps struct {
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
}

// A Condition for use in a Choice state branch.
// Experimental.
type Condition interface {
	RenderCondition() interface{}
}

// The jsii proxy struct for Condition
type jsiiProxy_Condition struct {
	_ byte // padding
}

// Experimental.
func NewCondition_Override(c Condition) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Condition",
		nil, // no parameters
		c,
	)
}

// Combine two or more conditions with a logical AND.
// Experimental.
func Condition_And(conditions ...Condition) Condition {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range conditions {
		args = append(args, a)
	}

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"and",
		args,
		&returns,
	)

	return returns
}

// Matches if a boolean field has the given value.
// Experimental.
func Condition_BooleanEquals(variable *string, value *bool) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"booleanEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a boolean field equals to a value at a given mapping path.
// Experimental.
func Condition_BooleanEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"booleanEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if variable is boolean.
// Experimental.
func Condition_IsBoolean(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isBoolean",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is not boolean.
// Experimental.
func Condition_IsNotBoolean(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNotBoolean",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is not null.
// Experimental.
func Condition_IsNotNull(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNotNull",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is not numeric.
// Experimental.
func Condition_IsNotNumeric(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNotNumeric",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is not present.
// Experimental.
func Condition_IsNotPresent(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNotPresent",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is not a string.
// Experimental.
func Condition_IsNotString(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNotString",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is not a timestamp.
// Experimental.
func Condition_IsNotTimestamp(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNotTimestamp",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is Null.
// Experimental.
func Condition_IsNull(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNull",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is numeric.
// Experimental.
func Condition_IsNumeric(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isNumeric",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is present.
// Experimental.
func Condition_IsPresent(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isPresent",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is a string.
// Experimental.
func Condition_IsString(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isString",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Matches if variable is a timestamp.
// Experimental.
func Condition_IsTimestamp(variable *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"isTimestamp",
		[]interface{}{variable},
		&returns,
	)

	return returns
}

// Negate a condition.
// Experimental.
func Condition_Not(condition Condition) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"not",
		[]interface{}{condition},
		&returns,
	)

	return returns
}

// Matches if a numeric field has the given value.
// Experimental.
func Condition_NumberEquals(variable *string, value *float64) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field has the value in a given mapping path.
// Experimental.
func Condition_NumberEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is greater than the given value.
// Experimental.
func Condition_NumberGreaterThan(variable *string, value *float64) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberGreaterThan",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is greater than or equal to the given value.
// Experimental.
func Condition_NumberGreaterThanEquals(variable *string, value *float64) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberGreaterThanEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is greater than or equal to the value at a given mapping path.
// Experimental.
func Condition_NumberGreaterThanEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberGreaterThanEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is greater than the value at a given mapping path.
// Experimental.
func Condition_NumberGreaterThanJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberGreaterThanJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is less than the given value.
// Experimental.
func Condition_NumberLessThan(variable *string, value *float64) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberLessThan",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is less than or equal to the given value.
// Experimental.
func Condition_NumberLessThanEquals(variable *string, value *float64) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberLessThanEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is less than or equal to the numeric value at given mapping path.
// Experimental.
func Condition_NumberLessThanEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberLessThanEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a numeric field is less than the value at the given mapping path.
// Experimental.
func Condition_NumberLessThanJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"numberLessThanJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Combine two or more conditions with a logical OR.
// Experimental.
func Condition_Or(conditions ...Condition) Condition {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range conditions {
		args = append(args, a)
	}

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"or",
		args,
		&returns,
	)

	return returns
}

// Matches if a string field has the given value.
// Experimental.
func Condition_StringEquals(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field equals to a value at a given mapping path.
// Experimental.
func Condition_StringEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts after a given value.
// Experimental.
func Condition_StringGreaterThan(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringGreaterThan",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts after or equal to a given value.
// Experimental.
func Condition_StringGreaterThanEquals(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringGreaterThanEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts after or equal to value at a given mapping path.
// Experimental.
func Condition_StringGreaterThanEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringGreaterThanEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts after a value at a given mapping path.
// Experimental.
func Condition_StringGreaterThanJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringGreaterThanJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts before a given value.
// Experimental.
func Condition_StringLessThan(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringLessThan",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts equal to or before a given value.
// Experimental.
func Condition_StringLessThanEquals(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringLessThanEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts equal to or before a given mapping.
// Experimental.
func Condition_StringLessThanEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringLessThanEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a string field sorts before a given value at a particular mapping.
// Experimental.
func Condition_StringLessThanJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringLessThanJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a field matches a string pattern that can contain a wild card (*) e.g: log-*.txt or *LATEST*. No other characters other than "*" have any special meaning - * can be escaped: \\*.
// Experimental.
func Condition_StringMatches(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"stringMatches",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is the same time as the given timestamp.
// Experimental.
func Condition_TimestampEquals(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is the same time as the timestamp at a given mapping path.
// Experimental.
func Condition_TimestampEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is after the given timestamp.
// Experimental.
func Condition_TimestampGreaterThan(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampGreaterThan",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is after or equal to the given timestamp.
// Experimental.
func Condition_TimestampGreaterThanEquals(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampGreaterThanEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is after or equal to the timestamp at a given mapping path.
// Experimental.
func Condition_TimestampGreaterThanEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampGreaterThanEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is after the timestamp at a given mapping path.
// Experimental.
func Condition_TimestampGreaterThanJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampGreaterThanJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is before the given timestamp.
// Experimental.
func Condition_TimestampLessThan(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampLessThan",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is before or equal to the given timestamp.
// Experimental.
func Condition_TimestampLessThanEquals(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampLessThanEquals",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is before or equal to the timestamp at a given mapping path.
// Experimental.
func Condition_TimestampLessThanEqualsJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampLessThanEqualsJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Matches if a timestamp field is before the timestamp at a given mapping path.
// Experimental.
func Condition_TimestampLessThanJsonPath(variable *string, value *string) Condition {
	_init_.Initialize()

	var returns Condition

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Condition",
		"timestampLessThanJsonPath",
		[]interface{}{variable, value},
		&returns,
	)

	return returns
}

// Render Amazon States Language JSON for the condition.
// Experimental.
func (c *jsiiProxy_Condition) RenderCondition() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderCondition",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Extract a field from the State Machine Context data.
// See: https://docs.aws.amazon.com/step-functions/latest/dg/connect-to-resource.html#wait-token-contextobject
//
// Deprecated: replaced by `JsonPath`
type Context interface {
}

// The jsii proxy struct for Context
type jsiiProxy_Context struct {
	_ byte // padding
}

// Instead of using a literal number, get the value from a JSON path.
// Deprecated: replaced by `JsonPath`
func Context_NumberAt(path *string) *float64 {
	_init_.Initialize()

	var returns *float64

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Context",
		"numberAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Instead of using a literal string, get the value from a JSON path.
// Deprecated: replaced by `JsonPath`
func Context_StringAt(path *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Context",
		"stringAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

func Context_EntireContext() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Context",
		"entireContext",
		&returns,
	)
	return returns
}

func Context_TaskToken() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Context",
		"taskToken",
		&returns,
	)
	return returns
}

// State defined by supplying Amazon States Language (ASL) in the state machine.
// Experimental.
type CustomState interface {
	State
	IChainable
	INextable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for CustomState
type jsiiProxy_CustomState struct {
	jsiiProxy_State
	jsiiProxy_IChainable
	jsiiProxy_INextable
}

func (j *jsiiProxy_CustomState) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomState) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewCustomState(scope constructs.Construct, id *string, props *CustomStateProps) CustomState {
	_init_.Initialize()

	j := jsiiProxy_CustomState{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.CustomState",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCustomState_Override(c CustomState, scope constructs.Construct, id *string, props *CustomStateProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.CustomState",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CustomState) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_CustomState) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func CustomState_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CustomState",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func CustomState_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CustomState",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func CustomState_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CustomState",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CustomState_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.CustomState",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func CustomState_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.CustomState",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (c *jsiiProxy_CustomState) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a choice branch to this state.
// Experimental.
func (c *jsiiProxy_CustomState) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		c,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (c *jsiiProxy_CustomState) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (c *jsiiProxy_CustomState) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		c,
		"addPrefix",
		[]interface{}{x},
	)
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (c *jsiiProxy_CustomState) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (c *jsiiProxy_CustomState) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		c,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (c *jsiiProxy_CustomState) MakeNext(next State) {
	_jsii_.InvokeVoid(
		c,
		"makeNext",
		[]interface{}{next},
	)
}

// Continue normal execution with the given state.
// Experimental.
func (c *jsiiProxy_CustomState) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		c,
		"next",
		[]interface{}{next},
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
func (c *jsiiProxy_CustomState) OnPrepare() {
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
func (c *jsiiProxy_CustomState) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CustomState) OnValidate() *[]*string {
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
func (c *jsiiProxy_CustomState) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (c *jsiiProxy_CustomState) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (c *jsiiProxy_CustomState) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (c *jsiiProxy_CustomState) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (c *jsiiProxy_CustomState) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (c *jsiiProxy_CustomState) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (c *jsiiProxy_CustomState) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (c *jsiiProxy_CustomState) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"renderRetryCatch",
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
func (c *jsiiProxy_CustomState) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns the Amazon States Language object for this state.
// Experimental.
func (c *jsiiProxy_CustomState) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CustomState) ToString() *string {
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
func (c *jsiiProxy_CustomState) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (c *jsiiProxy_CustomState) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		c,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a custom state definition.
// Experimental.
type CustomStateProps struct {
	// Amazon States Language (JSON-based) definition of the state.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/concepts-amazon-states-language.html
	//
	// Experimental.
	StateJson *map[string]interface{} `json:"stateJson"`
}

// Extract a field from the State Machine data that gets passed around between states.
// Deprecated: replaced by `JsonPath`
type Data interface {
}

// The jsii proxy struct for Data
type jsiiProxy_Data struct {
	_ byte // padding
}

// Determines if the indicated string is an encoded JSON path.
// Deprecated: replaced by `JsonPath`
func Data_IsJsonPathString(value *string) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Data",
		"isJsonPathString",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// Instead of using a literal string list, get the value from a JSON path.
// Deprecated: replaced by `JsonPath`
func Data_ListAt(path *string) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Data",
		"listAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Instead of using a literal number, get the value from a JSON path.
// Deprecated: replaced by `JsonPath`
func Data_NumberAt(path *string) *float64 {
	_init_.Initialize()

	var returns *float64

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Data",
		"numberAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Instead of using a literal string, get the value from a JSON path.
// Deprecated: replaced by `JsonPath`
func Data_StringAt(path *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Data",
		"stringAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

func Data_EntirePayload() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Data",
		"entirePayload",
		&returns,
	)
	return returns
}

// Predefined error strings Error names in Amazon States Language - https://states-language.net/spec.html#appendix-a Error handling in Step Functions - https://docs.aws.amazon.com/step-functions/latest/dg/concepts-error-handling.html.
// Experimental.
type Errors interface {
}

// The jsii proxy struct for Errors
type jsiiProxy_Errors struct {
	_ byte // padding
}

// Experimental.
func NewErrors() Errors {
	_init_.Initialize()

	j := jsiiProxy_Errors{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Errors",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewErrors_Override(e Errors) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Errors",
		nil, // no parameters
		e,
	)
}

func Errors_ALL() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"ALL",
		&returns,
	)
	return returns
}

func Errors_BRANCH_FAILED() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"BRANCH_FAILED",
		&returns,
	)
	return returns
}

func Errors_NO_CHOICE_MATCHED() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"NO_CHOICE_MATCHED",
		&returns,
	)
	return returns
}

func Errors_PARAMETER_PATH_FAILURE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"PARAMETER_PATH_FAILURE",
		&returns,
	)
	return returns
}

func Errors_PERMISSIONS() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"PERMISSIONS",
		&returns,
	)
	return returns
}

func Errors_RESULT_PATH_MATCH_FAILURE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"RESULT_PATH_MATCH_FAILURE",
		&returns,
	)
	return returns
}

func Errors_TASKS_FAILED() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"TASKS_FAILED",
		&returns,
	)
	return returns
}

func Errors_TIMEOUT() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.Errors",
		"TIMEOUT",
		&returns,
	)
	return returns
}

// Define a Fail state in the state machine.
//
// Reaching a Fail state terminates the state execution in failure.
// Experimental.
type Fail interface {
	State
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Fail
type jsiiProxy_Fail struct {
	jsiiProxy_State
}

func (j *jsiiProxy_Fail) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Fail) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewFail(scope constructs.Construct, id *string, props *FailProps) Fail {
	_init_.Initialize()

	j := jsiiProxy_Fail{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Fail",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewFail_Override(f Fail, scope constructs.Construct, id *string, props *FailProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Fail",
		[]interface{}{scope, id, props},
		f,
	)
}

func (j *jsiiProxy_Fail) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Fail) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func Fail_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Fail",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func Fail_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Fail",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func Fail_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Fail",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Fail_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Fail",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func Fail_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Fail",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (f *jsiiProxy_Fail) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		f,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a choice branch to this state.
// Experimental.
func (f *jsiiProxy_Fail) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		f,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (f *jsiiProxy_Fail) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		f,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (f *jsiiProxy_Fail) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		f,
		"addPrefix",
		[]interface{}{x},
	)
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (f *jsiiProxy_Fail) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		f,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (f *jsiiProxy_Fail) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		f,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (f *jsiiProxy_Fail) MakeNext(next State) {
	_jsii_.InvokeVoid(
		f,
		"makeNext",
		[]interface{}{next},
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
func (f *jsiiProxy_Fail) OnPrepare() {
	_jsii_.InvokeVoid(
		f,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (f *jsiiProxy_Fail) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
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
func (f *jsiiProxy_Fail) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_Fail) Prepare() {
	_jsii_.InvokeVoid(
		f,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (f *jsiiProxy_Fail) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		f,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (f *jsiiProxy_Fail) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		f,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (f *jsiiProxy_Fail) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		f,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (f *jsiiProxy_Fail) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		f,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (f *jsiiProxy_Fail) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		f,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (f *jsiiProxy_Fail) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		f,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (f *jsiiProxy_Fail) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		f,
		"renderRetryCatch",
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
func (f *jsiiProxy_Fail) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (f *jsiiProxy_Fail) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		f,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (f *jsiiProxy_Fail) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_Fail) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (f *jsiiProxy_Fail) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		f,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a Fail state.
// Experimental.
type FailProps struct {
	// A description for the cause of the failure.
	// Experimental.
	Cause *string `json:"cause"`
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// Error code used to represent this failure.
	// Experimental.
	Error *string `json:"error"`
}

// Helper functions to work with structures containing fields.
// Experimental.
type FieldUtils interface {
}

// The jsii proxy struct for FieldUtils
type jsiiProxy_FieldUtils struct {
	_ byte // padding
}

// Returns whether the given task structure contains the TaskToken field anywhere.
//
// The field is considered included if the field itself or one of its containing
// fields occurs anywhere in the payload.
// Experimental.
func FieldUtils_ContainsTaskToken(obj *map[string]interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.FieldUtils",
		"containsTaskToken",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Return all JSON paths used in the given structure.
// Experimental.
func FieldUtils_FindReferencedPaths(obj *map[string]interface{}) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.FieldUtils",
		"findReferencedPaths",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Render a JSON structure containing fields to the right StepFunctions structure.
// Experimental.
func FieldUtils_RenderObject(obj *map[string]interface{}) *map[string]interface{} {
	_init_.Initialize()

	var returns *map[string]interface{}

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.FieldUtils",
		"renderObject",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Options for finding reachable states.
// Experimental.
type FindStateOptions struct {
	// Whether or not to follow error-handling transitions.
	// Experimental.
	IncludeErrorHandlers *bool `json:"includeErrorHandlers"`
}

// Represents a Step Functions Activity https://docs.aws.amazon.com/step-functions/latest/dg/concepts-activities.html.
// Experimental.
type IActivity interface {
	awscdk.IResource
	// The ARN of the activity.
	// Experimental.
	ActivityArn() *string
	// The name of the activity.
	// Experimental.
	ActivityName() *string
}

// The jsii proxy for IActivity
type jsiiProxy_IActivity struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IActivity) ActivityArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"activityArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IActivity) ActivityName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"activityName",
		&returns,
	)
	return returns
}

// Interface for objects that can be used in a Chain.
// Experimental.
type IChainable interface {
	// The chainable end state(s) of this chainable.
	// Experimental.
	EndStates() *[]INextable
	// Descriptive identifier for this chainable.
	// Experimental.
	Id() *string
	// The start state of this chainable.
	// Experimental.
	StartState() State
}

// The jsii proxy for IChainable
type jsiiProxy_IChainable struct {
	_ byte // padding
}

func (j *jsiiProxy_IChainable) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IChainable) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IChainable) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

// Interface for states that can have 'next' states.
// Experimental.
type INextable interface {
	// Go to the indicated state after this state.
	//
	// Returns: The chain of states built up
	// Experimental.
	Next(state IChainable) Chain
}

// The jsii proxy for INextable
type jsiiProxy_INextable struct {
	_ byte // padding
}

func (i *jsiiProxy_INextable) Next(state IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		i,
		"next",
		[]interface{}{state},
		&returns,
	)

	return returns
}

// A State Machine.
// Experimental.
type IStateMachine interface {
	awsiam.IGrantable
	awscdk.IResource
	// Grant the given identity custom permissions.
	// Experimental.
	Grant(identity awsiam.IGrantable, actions ...*string) awsiam.Grant
	// Grant the given identity permissions for all executions of a state machine.
	// Experimental.
	GrantExecution(identity awsiam.IGrantable, actions ...*string) awsiam.Grant
	// Grant the given identity read permissions for this state machine.
	// Experimental.
	GrantRead(identity awsiam.IGrantable) awsiam.Grant
	// Grant the given identity permissions to start an execution of this state machine.
	// Experimental.
	GrantStartExecution(identity awsiam.IGrantable) awsiam.Grant
	// Grant the given identity read permissions for this state machine.
	// Experimental.
	GrantTaskResponse(identity awsiam.IGrantable) awsiam.Grant
	// Return the given named metric for this State Machine's executions.
	// Experimental.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of executions that were aborted.
	// Experimental.
	MetricAborted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of executions that failed.
	// Experimental.
	MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of executions that were started.
	// Experimental.
	MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of executions that succeeded.
	// Experimental.
	MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of executions that were throttled.
	// Experimental.
	MetricThrottled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the interval, in milliseconds, between the time the execution starts and the time it closes.
	// Experimental.
	MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of executions that timed out.
	// Experimental.
	MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The ARN of the state machine.
	// Experimental.
	StateMachineArn() *string
}

// The jsii proxy for IStateMachine
type jsiiProxy_IStateMachine struct {
	internal.Type__awsiamIGrantable
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IStateMachine) Grant(identity awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{identity}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grant",
		args,
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) GrantExecution(identity awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{identity}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantExecution",
		args,
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) GrantRead(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantRead",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) GrantStartExecution(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantStartExecution",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) GrantTaskResponse(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantTaskResponse",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) MetricAborted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricAborted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricStarted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricSucceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) MetricThrottled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricThrottled",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStateMachine) MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricTimedOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IStateMachine) StateMachineArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateMachineArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStateMachine) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStateMachine) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStateMachine) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStateMachine) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// Interface for resources that can be used as tasks.
// Deprecated: replaced by `TaskStateBase`.
type IStepFunctionsTask interface {
	// Called when the task object is used in a workflow.
	// Deprecated: replaced by `TaskStateBase`.
	Bind(task Task) *StepFunctionsTaskConfig
}

// The jsii proxy for IStepFunctionsTask
type jsiiProxy_IStepFunctionsTask struct {
	_ byte // padding
}

func (i *jsiiProxy_IStepFunctionsTask) Bind(task Task) *StepFunctionsTaskConfig {
	var returns *StepFunctionsTaskConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{task},
		&returns,
	)

	return returns
}

// The type of task input.
// Experimental.
type InputType string

const (
	InputType_TEXT InputType = "TEXT"
	InputType_OBJECT InputType = "OBJECT"
)

// AWS Step Functions integrates with services directly in the Amazon States Language.
//
// You can control these AWS services using service integration patterns:
// See: https://docs.aws.amazon.com/step-functions/latest/dg/connect-to-resource.html
//
// Experimental.
type IntegrationPattern string

const (
	IntegrationPattern_REQUEST_RESPONSE IntegrationPattern = "REQUEST_RESPONSE"
	IntegrationPattern_RUN_JOB IntegrationPattern = "RUN_JOB"
	IntegrationPattern_WAIT_FOR_TASK_TOKEN IntegrationPattern = "WAIT_FOR_TASK_TOKEN"
)

// Extract a field from the State Machine data or context that gets passed around between states.
// See: https://docs.aws.amazon.com/step-functions/latest/dg/amazon-states-language-paths.html
//
// Experimental.
type JsonPath interface {
}

// The jsii proxy struct for JsonPath
type jsiiProxy_JsonPath struct {
	_ byte // padding
}

// Determines if the indicated string is an encoded JSON path.
// Experimental.
func JsonPath_IsEncodedJsonPath(value *string) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.JsonPath",
		"isEncodedJsonPath",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// Instead of using a literal string list, get the value from a JSON path.
// Experimental.
func JsonPath_ListAt(path *string) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.JsonPath",
		"listAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Instead of using a literal number, get the value from a JSON path.
// Experimental.
func JsonPath_NumberAt(path *string) *float64 {
	_init_.Initialize()

	var returns *float64

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.JsonPath",
		"numberAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Instead of using a literal string, get the value from a JSON path.
// Experimental.
func JsonPath_StringAt(path *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.JsonPath",
		"stringAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

func JsonPath_DISCARD() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.JsonPath",
		"DISCARD",
		&returns,
	)
	return returns
}

func JsonPath_EntireContext() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.JsonPath",
		"entireContext",
		&returns,
	)
	return returns
}

func JsonPath_EntirePayload() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.JsonPath",
		"entirePayload",
		&returns,
	)
	return returns
}

func JsonPath_TaskToken() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_stepfunctions.JsonPath",
		"taskToken",
		&returns,
	)
	return returns
}

// Defines which category of execution history events are logged.
// See: https://docs.aws.amazon.com/step-functions/latest/dg/cloudwatch-log-level.html
//
// Experimental.
type LogLevel string

const (
	LogLevel_OFF LogLevel = "OFF"
	LogLevel_ALL LogLevel = "ALL"
	LogLevel_ERROR LogLevel = "ERROR"
	LogLevel_FATAL LogLevel = "FATAL"
)

// Defines what execution history events are logged and where they are logged.
// Experimental.
type LogOptions struct {
	// The log group where the execution history events will be logged.
	// Experimental.
	Destination awslogs.ILogGroup `json:"destination"`
	// Determines whether execution data is included in your log.
	// Experimental.
	IncludeExecutionData *bool `json:"includeExecutionData"`
	// Defines which category of execution history events are logged.
	// Experimental.
	Level LogLevel `json:"level"`
}

// Define a Map state in the state machine.
//
// A `Map` state can be used to run a set of steps for each element of an input array.
// A Map state will execute the same steps for multiple entries of an array in the state input.
//
// While the Parallel state executes multiple branches of steps using the same input, a Map state
// will execute the same steps for multiple entries of an array in the state input.
// See: https://docs.aws.amazon.com/step-functions/latest/dg/amazon-states-language-map-state.html
//
// Experimental.
type Map interface {
	State
	INextable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddCatch(handler IChainable, props *CatchProps) Map
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	AddRetry(props *RetryProps) Map
	BindToGraph(graph StateGraph)
	Iterator(iterator IChainable) Map
	MakeDefault(def State)
	MakeNext(next State)
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Map
type jsiiProxy_Map struct {
	jsiiProxy_State
	jsiiProxy_INextable
}

func (j *jsiiProxy_Map) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Map) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewMap(scope constructs.Construct, id *string, props *MapProps) Map {
	_init_.Initialize()

	j := jsiiProxy_Map{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Map",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewMap_Override(m Map, scope constructs.Construct, id *string, props *MapProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Map",
		[]interface{}{scope, id, props},
		m,
	)
}

func (j *jsiiProxy_Map) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Map) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func Map_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Map",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func Map_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Map",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func Map_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Map",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Map_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Map",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func Map_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Map",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (m *jsiiProxy_Map) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		m,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a recovery handler for this state.
//
// When a particular error occurs, execution will continue at the error
// handler instead of failing the state machine execution.
// Experimental.
func (m *jsiiProxy_Map) AddCatch(handler IChainable, props *CatchProps) Map {
	var returns Map

	_jsii_.Invoke(
		m,
		"addCatch",
		[]interface{}{handler, props},
		&returns,
	)

	return returns
}

// Add a choice branch to this state.
// Experimental.
func (m *jsiiProxy_Map) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		m,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (m *jsiiProxy_Map) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		m,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (m *jsiiProxy_Map) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		m,
		"addPrefix",
		[]interface{}{x},
	)
}

// Add retry configuration for this state.
//
// This controls if and how the execution will be retried if a particular
// error occurs.
// Experimental.
func (m *jsiiProxy_Map) AddRetry(props *RetryProps) Map {
	var returns Map

	_jsii_.Invoke(
		m,
		"addRetry",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (m *jsiiProxy_Map) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		m,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Define iterator state machine in Map.
// Experimental.
func (m *jsiiProxy_Map) Iterator(iterator IChainable) Map {
	var returns Map

	_jsii_.Invoke(
		m,
		"iterator",
		[]interface{}{iterator},
		&returns,
	)

	return returns
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (m *jsiiProxy_Map) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		m,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (m *jsiiProxy_Map) MakeNext(next State) {
	_jsii_.InvokeVoid(
		m,
		"makeNext",
		[]interface{}{next},
	)
}

// Continue normal execution with the given state.
// Experimental.
func (m *jsiiProxy_Map) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		m,
		"next",
		[]interface{}{next},
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
func (m *jsiiProxy_Map) OnPrepare() {
	_jsii_.InvokeVoid(
		m,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (m *jsiiProxy_Map) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		m,
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
func (m *jsiiProxy_Map) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		m,
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
func (m *jsiiProxy_Map) Prepare() {
	_jsii_.InvokeVoid(
		m,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (m *jsiiProxy_Map) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		m,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (m *jsiiProxy_Map) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		m,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (m *jsiiProxy_Map) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		m,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (m *jsiiProxy_Map) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		m,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (m *jsiiProxy_Map) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		m,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (m *jsiiProxy_Map) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		m,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (m *jsiiProxy_Map) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		m,
		"renderRetryCatch",
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
func (m *jsiiProxy_Map) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		m,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (m *jsiiProxy_Map) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		m,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (m *jsiiProxy_Map) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		m,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate this state.
// Experimental.
func (m *jsiiProxy_Map) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		m,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (m *jsiiProxy_Map) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		m,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a Map state.
// Experimental.
type MapProps struct {
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select the array to iterate over.
	// Experimental.
	ItemsPath *string `json:"itemsPath"`
	// MaxConcurrency.
	//
	// An upper bound on the number of iterations you want running at once.
	// Experimental.
	MaxConcurrency *float64 `json:"maxConcurrency"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
	// The JSON that you want to override your default iteration input.
	// Experimental.
	Parameters *map[string]interface{} `json:"parameters"`
	// JSONPath expression to indicate where to inject the state's output.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the state's
	// input to become its output.
	// Experimental.
	ResultPath *string `json:"resultPath"`
	// The JSON that will replace the state's raw result and become the effective result before ResultPath is applied.
	//
	// You can use ResultSelector to create a payload with values that are static
	// or selected from the state's raw result.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-resultselector
	//
	// Experimental.
	ResultSelector *map[string]interface{} `json:"resultSelector"`
}

// Define a Parallel state in the state machine.
//
// A Parallel state can be used to run one or more state machines at the same
// time.
//
// The Result of a Parallel state is an array of the results of its substatemachines.
// Experimental.
type Parallel interface {
	State
	INextable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddCatch(handler IChainable, props *CatchProps) Parallel
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	AddRetry(props *RetryProps) Parallel
	BindToGraph(graph StateGraph)
	Branch(branches ...IChainable) Parallel
	MakeDefault(def State)
	MakeNext(next State)
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Parallel
type jsiiProxy_Parallel struct {
	jsiiProxy_State
	jsiiProxy_INextable
}

func (j *jsiiProxy_Parallel) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Parallel) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewParallel(scope constructs.Construct, id *string, props *ParallelProps) Parallel {
	_init_.Initialize()

	j := jsiiProxy_Parallel{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Parallel",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewParallel_Override(p Parallel, scope constructs.Construct, id *string, props *ParallelProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Parallel",
		[]interface{}{scope, id, props},
		p,
	)
}

func (j *jsiiProxy_Parallel) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Parallel) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func Parallel_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Parallel",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func Parallel_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Parallel",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func Parallel_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Parallel",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Parallel_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Parallel",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func Parallel_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Parallel",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (p *jsiiProxy_Parallel) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a recovery handler for this state.
//
// When a particular error occurs, execution will continue at the error
// handler instead of failing the state machine execution.
// Experimental.
func (p *jsiiProxy_Parallel) AddCatch(handler IChainable, props *CatchProps) Parallel {
	var returns Parallel

	_jsii_.Invoke(
		p,
		"addCatch",
		[]interface{}{handler, props},
		&returns,
	)

	return returns
}

// Add a choice branch to this state.
// Experimental.
func (p *jsiiProxy_Parallel) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		p,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (p *jsiiProxy_Parallel) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (p *jsiiProxy_Parallel) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		p,
		"addPrefix",
		[]interface{}{x},
	)
}

// Add retry configuration for this state.
//
// This controls if and how the execution will be retried if a particular
// error occurs.
// Experimental.
func (p *jsiiProxy_Parallel) AddRetry(props *RetryProps) Parallel {
	var returns Parallel

	_jsii_.Invoke(
		p,
		"addRetry",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (p *jsiiProxy_Parallel) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Define one or more branches to run in parallel.
// Experimental.
func (p *jsiiProxy_Parallel) Branch(branches ...IChainable) Parallel {
	args := []interface{}{}
	for _, a := range branches {
		args = append(args, a)
	}

	var returns Parallel

	_jsii_.Invoke(
		p,
		"branch",
		args,
		&returns,
	)

	return returns
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (p *jsiiProxy_Parallel) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		p,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (p *jsiiProxy_Parallel) MakeNext(next State) {
	_jsii_.InvokeVoid(
		p,
		"makeNext",
		[]interface{}{next},
	)
}

// Continue normal execution with the given state.
// Experimental.
func (p *jsiiProxy_Parallel) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		p,
		"next",
		[]interface{}{next},
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
func (p *jsiiProxy_Parallel) OnPrepare() {
	_jsii_.InvokeVoid(
		p,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_Parallel) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
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
func (p *jsiiProxy_Parallel) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
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
func (p *jsiiProxy_Parallel) Prepare() {
	_jsii_.InvokeVoid(
		p,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Parallel) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Parallel) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Parallel) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Parallel) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Parallel) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Parallel) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Parallel) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderRetryCatch",
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
func (p *jsiiProxy_Parallel) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (p *jsiiProxy_Parallel) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		p,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_Parallel) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate this state.
// Experimental.
func (p *jsiiProxy_Parallel) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (p *jsiiProxy_Parallel) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a Parallel state.
// Experimental.
type ParallelProps struct {
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
	// JSONPath expression to indicate where to inject the state's output.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the state's
	// input to become its output.
	// Experimental.
	ResultPath *string `json:"resultPath"`
	// The JSON that will replace the state's raw result and become the effective result before ResultPath is applied.
	//
	// You can use ResultSelector to create a payload with values that are static
	// or selected from the state's raw result.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-resultselector
	//
	// Experimental.
	ResultSelector *map[string]interface{} `json:"resultSelector"`
}

// Define a Pass in the state machine.
//
// A Pass state can be used to transform the current exeuction's state.
// Experimental.
type Pass interface {
	State
	INextable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Pass
type jsiiProxy_Pass struct {
	jsiiProxy_State
	jsiiProxy_INextable
}

func (j *jsiiProxy_Pass) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pass) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewPass(scope constructs.Construct, id *string, props *PassProps) Pass {
	_init_.Initialize()

	j := jsiiProxy_Pass{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Pass",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewPass_Override(p Pass, scope constructs.Construct, id *string, props *PassProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Pass",
		[]interface{}{scope, id, props},
		p,
	)
}

func (j *jsiiProxy_Pass) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Pass) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func Pass_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Pass",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func Pass_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Pass",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func Pass_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Pass",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Pass_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Pass",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func Pass_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Pass",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (p *jsiiProxy_Pass) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a choice branch to this state.
// Experimental.
func (p *jsiiProxy_Pass) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		p,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (p *jsiiProxy_Pass) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (p *jsiiProxy_Pass) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		p,
		"addPrefix",
		[]interface{}{x},
	)
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (p *jsiiProxy_Pass) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (p *jsiiProxy_Pass) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		p,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (p *jsiiProxy_Pass) MakeNext(next State) {
	_jsii_.InvokeVoid(
		p,
		"makeNext",
		[]interface{}{next},
	)
}

// Continue normal execution with the given state.
// Experimental.
func (p *jsiiProxy_Pass) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		p,
		"next",
		[]interface{}{next},
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
func (p *jsiiProxy_Pass) OnPrepare() {
	_jsii_.InvokeVoid(
		p,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_Pass) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
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
func (p *jsiiProxy_Pass) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
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
func (p *jsiiProxy_Pass) Prepare() {
	_jsii_.InvokeVoid(
		p,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Pass) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Pass) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Pass) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Pass) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Pass) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Pass) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (p *jsiiProxy_Pass) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"renderRetryCatch",
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
func (p *jsiiProxy_Pass) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (p *jsiiProxy_Pass) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		p,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_Pass) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		p,
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
func (p *jsiiProxy_Pass) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (p *jsiiProxy_Pass) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		p,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a Pass state.
// Experimental.
type PassProps struct {
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
	// Parameters pass a collection of key-value pairs, either static values or JSONPath expressions that select from the input.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-parameters
	//
	// Experimental.
	Parameters *map[string]interface{} `json:"parameters"`
	// If given, treat as the result of this operation.
	//
	// Can be used to inject or replace the current execution state.
	// Experimental.
	Result Result `json:"result"`
	// JSONPath expression to indicate where to inject the state's output.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the state's
	// input to become its output.
	// Experimental.
	ResultPath *string `json:"resultPath"`
}

// The result of a Pass operation.
// Experimental.
type Result interface {
	Value() interface{}
}

// The jsii proxy struct for Result
type jsiiProxy_Result struct {
	_ byte // padding
}

func (j *jsiiProxy_Result) Value() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


// Experimental.
func NewResult(value interface{}) Result {
	_init_.Initialize()

	j := jsiiProxy_Result{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Result",
		[]interface{}{value},
		&j,
	)

	return &j
}

// Experimental.
func NewResult_Override(r Result, value interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Result",
		[]interface{}{value},
		r,
	)
}

// The result of the operation is an array.
// Experimental.
func Result_FromArray(value *[]interface{}) Result {
	_init_.Initialize()

	var returns Result

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Result",
		"fromArray",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// The result of the operation is a boolean.
// Experimental.
func Result_FromBoolean(value *bool) Result {
	_init_.Initialize()

	var returns Result

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Result",
		"fromBoolean",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// The result of the operation is a number.
// Experimental.
func Result_FromNumber(value *float64) Result {
	_init_.Initialize()

	var returns Result

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Result",
		"fromNumber",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// The result of the operation is an object.
// Experimental.
func Result_FromObject(value *map[string]interface{}) Result {
	_init_.Initialize()

	var returns Result

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Result",
		"fromObject",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// The result of the operation is a string.
// Experimental.
func Result_FromString(value *string) Result {
	_init_.Initialize()

	var returns Result

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Result",
		"fromString",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// Retry details.
// Experimental.
type RetryProps struct {
	// Multiplication for how much longer the wait interval gets on every retry.
	// Experimental.
	BackoffRate *float64 `json:"backoffRate"`
	// Errors to retry.
	//
	// A list of error strings to retry, which can be either predefined errors
	// (for example Errors.NoChoiceMatched) or a self-defined error.
	// Experimental.
	Errors *[]*string `json:"errors"`
	// How many seconds to wait initially before retrying.
	// Experimental.
	Interval awscdk.Duration `json:"interval"`
	// How many times to retry this particular error.
	//
	// May be 0 to disable retry for specific errors (in case you have
	// a catch-all retry policy).
	// Experimental.
	MaxAttempts *float64 `json:"maxAttempts"`
}

// Three ways to call an integrated service: Request Response, Run a Job and Wait for a Callback with Task Token.
// See: https://docs.aws.amazon.com/step-functions/latest/dg/connect-to-resource.html
//
// Here, they are named as FIRE_AND_FORGET, SYNC and WAIT_FOR_TASK_TOKEN respectfully.
//
// Experimental.
type ServiceIntegrationPattern string

const (
	ServiceIntegrationPattern_FIRE_AND_FORGET ServiceIntegrationPattern = "FIRE_AND_FORGET"
	ServiceIntegrationPattern_SYNC ServiceIntegrationPattern = "SYNC"
	ServiceIntegrationPattern_WAIT_FOR_TASK_TOKEN ServiceIntegrationPattern = "WAIT_FOR_TASK_TOKEN"
)

// Options for creating a single state.
// Experimental.
type SingleStateOptions struct {
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
	// JSONPath expression to indicate where to inject the state's output.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the state's
	// input to become its output.
	// Experimental.
	ResultPath *string `json:"resultPath"`
	// The JSON that will replace the state's raw result and become the effective result before ResultPath is applied.
	//
	// You can use ResultSelector to create a payload with values that are static
	// or selected from the state's raw result.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-resultselector
	//
	// Experimental.
	ResultSelector *map[string]interface{} `json:"resultSelector"`
	// String to prefix all stateIds in the state machine with.
	// Experimental.
	PrefixStates *string `json:"prefixStates"`
	// ID of newly created containing state.
	// Experimental.
	StateId *string `json:"stateId"`
}

// Base class for all other state classes.
// Experimental.
type State interface {
	awscdk.Construct
	IChainable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for State
type jsiiProxy_State struct {
	internal.Type__awscdkConstruct
	jsiiProxy_IChainable
}

func (j *jsiiProxy_State) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_State) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewState_Override(s State, scope constructs.Construct, id *string, props *StateProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.State",
		[]interface{}{scope, id, props},
		s,
	)
}

func (j *jsiiProxy_State) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_State) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func State_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.State",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func State_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.State",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func State_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.State",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func State_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.State",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func State_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.State",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (s *jsiiProxy_State) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a choice branch to this state.
// Experimental.
func (s *jsiiProxy_State) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		s,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (s *jsiiProxy_State) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (s *jsiiProxy_State) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		s,
		"addPrefix",
		[]interface{}{x},
	)
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (s *jsiiProxy_State) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (s *jsiiProxy_State) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		s,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (s *jsiiProxy_State) MakeNext(next State) {
	_jsii_.InvokeVoid(
		s,
		"makeNext",
		[]interface{}{next},
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
func (s *jsiiProxy_State) OnPrepare() {
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
func (s *jsiiProxy_State) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_State) OnValidate() *[]*string {
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
func (s *jsiiProxy_State) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (s *jsiiProxy_State) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (s *jsiiProxy_State) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (s *jsiiProxy_State) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (s *jsiiProxy_State) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (s *jsiiProxy_State) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (s *jsiiProxy_State) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (s *jsiiProxy_State) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderRetryCatch",
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
func (s *jsiiProxy_State) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Render the state as JSON.
// Experimental.
func (s *jsiiProxy_State) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_State) ToString() *string {
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
func (s *jsiiProxy_State) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (s *jsiiProxy_State) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// A collection of connected states.
//
// A StateGraph is used to keep track of all states that are connected (have
// transitions between them). It does not include the substatemachines in
// a Parallel's branches: those are their own StateGraphs, but the graphs
// themselves have a hierarchical relationship as well.
//
// By assigning states to a definitive StateGraph, we verify that no state
// machines are constructed. In particular:
//
// - Every state object can only ever be in 1 StateGraph, and not inadvertently
//    be used in two graphs.
// - Every stateId must be unique across all states in the entire state
//    machine.
//
// All policy statements in all states in all substatemachines are bubbled so
// that the top-level StateMachine instantiation can read them all and add
// them to the IAM Role.
//
// You do not need to instantiate this class; it is used internally.
// Experimental.
type StateGraph interface {
	PolicyStatements() *[]awsiam.PolicyStatement
	StartState() State
	Timeout() awscdk.Duration
	SetTimeout(val awscdk.Duration)
	RegisterPolicyStatement(statement awsiam.PolicyStatement)
	RegisterState(state State)
	RegisterSuperGraph(graph StateGraph)
	ToGraphJson() *map[string]interface{}
	ToString() *string
}

// The jsii proxy struct for StateGraph
type jsiiProxy_StateGraph struct {
	_ byte // padding
}

func (j *jsiiProxy_StateGraph) PolicyStatements() *[]awsiam.PolicyStatement {
	var returns *[]awsiam.PolicyStatement
	_jsii_.Get(
		j,
		"policyStatements",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateGraph) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateGraph) Timeout() awscdk.Duration {
	var returns awscdk.Duration
	_jsii_.Get(
		j,
		"timeout",
		&returns,
	)
	return returns
}


// Experimental.
func NewStateGraph(startState State, graphDescription *string) StateGraph {
	_init_.Initialize()

	j := jsiiProxy_StateGraph{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.StateGraph",
		[]interface{}{startState, graphDescription},
		&j,
	)

	return &j
}

// Experimental.
func NewStateGraph_Override(s StateGraph, startState State, graphDescription *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.StateGraph",
		[]interface{}{startState, graphDescription},
		s,
	)
}

func (j *jsiiProxy_StateGraph) SetTimeout(val awscdk.Duration) {
	_jsii_.Set(
		j,
		"timeout",
		val,
	)
}

// Register a Policy Statement used by states in this graph.
// Experimental.
func (s *jsiiProxy_StateGraph) RegisterPolicyStatement(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		s,
		"registerPolicyStatement",
		[]interface{}{statement},
	)
}

// Register a state as part of this graph.
//
// Called by State.bindToGraph().
// Experimental.
func (s *jsiiProxy_StateGraph) RegisterState(state State) {
	_jsii_.InvokeVoid(
		s,
		"registerState",
		[]interface{}{state},
	)
}

// Register this graph as a child of the given graph.
//
// Resource changes will be bubbled up to the given graph.
// Experimental.
func (s *jsiiProxy_StateGraph) RegisterSuperGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"registerSuperGraph",
		[]interface{}{graph},
	)
}

// Return the Amazon States Language JSON for this graph.
// Experimental.
func (s *jsiiProxy_StateGraph) ToGraphJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"toGraphJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return a string description of this graph.
// Experimental.
func (s *jsiiProxy_StateGraph) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Define a StepFunctions State Machine.
// Experimental.
type StateMachine interface {
	awscdk.Resource
	IStateMachine
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() awsiam.IPrincipal
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	StateMachineArn() *string
	StateMachineName() *string
	StateMachineType() StateMachineType
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(identity awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantExecution(identity awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantRead(identity awsiam.IGrantable) awsiam.Grant
	GrantStartExecution(identity awsiam.IGrantable) awsiam.Grant
	GrantTaskResponse(identity awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricAborted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for StateMachine
type jsiiProxy_StateMachine struct {
	internal.Type__awscdkResource
	jsiiProxy_IStateMachine
}

func (j *jsiiProxy_StateMachine) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) StateMachineArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateMachineArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) StateMachineName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateMachineName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachine) StateMachineType() StateMachineType {
	var returns StateMachineType
	_jsii_.Get(
		j,
		"stateMachineType",
		&returns,
	)
	return returns
}


// Experimental.
func NewStateMachine(scope constructs.Construct, id *string, props *StateMachineProps) StateMachine {
	_init_.Initialize()

	j := jsiiProxy_StateMachine{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.StateMachine",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStateMachine_Override(s StateMachine, scope constructs.Construct, id *string, props *StateMachineProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.StateMachine",
		[]interface{}{scope, id, props},
		s,
	)
}

// Import a state machine.
// Experimental.
func StateMachine_FromStateMachineArn(scope constructs.Construct, id *string, stateMachineArn *string) IStateMachine {
	_init_.Initialize()

	var returns IStateMachine

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateMachine",
		"fromStateMachineArn",
		[]interface{}{scope, id, stateMachineArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func StateMachine_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateMachine",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func StateMachine_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateMachine",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add the given statement to the role's policy.
// Experimental.
func (s *jsiiProxy_StateMachine) AddToRolePolicy(statement awsiam.PolicyStatement) {
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
func (s *jsiiProxy_StateMachine) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (s *jsiiProxy_StateMachine) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_StateMachine) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_StateMachine) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity custom permissions.
// Experimental.
func (s *jsiiProxy_StateMachine) Grant(identity awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{identity}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Grant the given identity permissions on all executions of the state machine.
// Experimental.
func (s *jsiiProxy_StateMachine) GrantExecution(identity awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{identity}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantExecution",
		args,
		&returns,
	)

	return returns
}

// Grant the given identity permissions to read results from state machine.
// Experimental.
func (s *jsiiProxy_StateMachine) GrantRead(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantRead",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to start an execution of this state machine.
// Experimental.
func (s *jsiiProxy_StateMachine) GrantStartExecution(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantStartExecution",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

// Grant the given identity task response permissions on a state machine.
// Experimental.
func (s *jsiiProxy_StateMachine) GrantTaskResponse(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantTaskResponse",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

// Return the given named metric for this State Machine's executions.
// Experimental.
func (s *jsiiProxy_StateMachine) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of executions that were aborted.
// Experimental.
func (s *jsiiProxy_StateMachine) MetricAborted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricAborted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of executions that failed.
// Experimental.
func (s *jsiiProxy_StateMachine) MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of executions that were started.
// Experimental.
func (s *jsiiProxy_StateMachine) MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricStarted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of executions that succeeded.
// Experimental.
func (s *jsiiProxy_StateMachine) MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricSucceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of executions that were throttled.
// Experimental.
func (s *jsiiProxy_StateMachine) MetricThrottled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricThrottled",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the interval, in milliseconds, between the time the execution starts and the time it closes.
// Experimental.
func (s *jsiiProxy_StateMachine) MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of executions that timed out.
// Experimental.
func (s *jsiiProxy_StateMachine) MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricTimedOut",
		[]interface{}{props},
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
func (s *jsiiProxy_StateMachine) OnPrepare() {
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
func (s *jsiiProxy_StateMachine) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_StateMachine) OnValidate() *[]*string {
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
func (s *jsiiProxy_StateMachine) Prepare() {
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
func (s *jsiiProxy_StateMachine) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_StateMachine) ToString() *string {
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
func (s *jsiiProxy_StateMachine) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Base class for reusable state machine fragments.
// Experimental.
type StateMachineFragment interface {
	awscdk.Construct
	IChainable
	EndStates() *[]INextable
	Id() *string
	Node() awscdk.ConstructNode
	StartState() State
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	PrefixStates(prefix *string) StateMachineFragment
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToSingleState(options *SingleStateOptions) Parallel
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for StateMachineFragment
type jsiiProxy_StateMachineFragment struct {
	internal.Type__awscdkConstruct
	jsiiProxy_IChainable
}

func (j *jsiiProxy_StateMachineFragment) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachineFragment) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachineFragment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StateMachineFragment) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}


// Experimental.
func NewStateMachineFragment_Override(s StateMachineFragment, scope constructs.Construct, id *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.StateMachineFragment",
		[]interface{}{scope, id},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func StateMachineFragment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateMachineFragment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Continue normal execution with the given state.
// Experimental.
func (s *jsiiProxy_StateMachineFragment) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		s,
		"next",
		[]interface{}{next},
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
func (s *jsiiProxy_StateMachineFragment) OnPrepare() {
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
func (s *jsiiProxy_StateMachineFragment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_StateMachineFragment) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Prefix the IDs of all states in this state machine fragment.
//
// Use this to avoid multiple copies of the state machine all having the
// same state IDs.
// Experimental.
func (s *jsiiProxy_StateMachineFragment) PrefixStates(prefix *string) StateMachineFragment {
	var returns StateMachineFragment

	_jsii_.Invoke(
		s,
		"prefixStates",
		[]interface{}{prefix},
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
func (s *jsiiProxy_StateMachineFragment) Prepare() {
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
func (s *jsiiProxy_StateMachineFragment) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Wrap all states in this state machine fragment up into a single state.
//
// This can be used to add retry or error handling onto this state
// machine fragment.
//
// Be aware that this changes the result of the inner state machine
// to be an array with the result of the state machine in it. Adjust
// your paths accordingly. For example, change 'outputPath' to
// '$[0]'.
// Experimental.
func (s *jsiiProxy_StateMachineFragment) ToSingleState(options *SingleStateOptions) Parallel {
	var returns Parallel

	_jsii_.Invoke(
		s,
		"toSingleState",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_StateMachineFragment) ToString() *string {
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
func (s *jsiiProxy_StateMachineFragment) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining a State Machine.
// Experimental.
type StateMachineProps struct {
	// Definition for this state machine.
	// Experimental.
	Definition IChainable `json:"definition"`
	// Defines what execution history events are logged and where they are logged.
	// Experimental.
	Logs *LogOptions `json:"logs"`
	// The execution role for the state machine service.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// A name for the state machine.
	// Experimental.
	StateMachineName *string `json:"stateMachineName"`
	// Type of the state machine.
	// Experimental.
	StateMachineType StateMachineType `json:"stateMachineType"`
	// Maximum run time for this state machine.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// Specifies whether Amazon X-Ray tracing is enabled for this state machine.
	// Experimental.
	TracingEnabled *bool `json:"tracingEnabled"`
}

// Two types of state machines are available in AWS Step Functions: EXPRESS AND STANDARD.
// See: https://docs.aws.amazon.com/step-functions/latest/dg/concepts-standard-vs-express.html
//
// Experimental.
type StateMachineType string

const (
	StateMachineType_EXPRESS StateMachineType = "EXPRESS"
	StateMachineType_STANDARD StateMachineType = "STANDARD"
)

// Properties shared by all states.
// Experimental.
type StateProps struct {
	// A comment describing this state.
	// Experimental.
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
	// Parameters pass a collection of key-value pairs, either static values or JSONPath expressions that select from the input.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-parameters
	//
	// Experimental.
	Parameters *map[string]interface{} `json:"parameters"`
	// JSONPath expression to indicate where to inject the state's output.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the state's
	// input to become its output.
	// Experimental.
	ResultPath *string `json:"resultPath"`
	// The JSON that will replace the state's raw result and become the effective result before ResultPath is applied.
	//
	// You can use ResultSelector to create a payload with values that are static
	// or selected from the state's raw result.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-resultselector
	//
	// Experimental.
	ResultSelector *map[string]interface{} `json:"resultSelector"`
}

// Metrics on the rate limiting performed on state machine execution.
//
// These rate limits are shared across all state machines.
// Experimental.
type StateTransitionMetric interface {
}

// The jsii proxy struct for StateTransitionMetric
type jsiiProxy_StateTransitionMetric struct {
	_ byte // padding
}

// Experimental.
func NewStateTransitionMetric() StateTransitionMetric {
	_init_.Initialize()

	j := jsiiProxy_StateTransitionMetric{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewStateTransitionMetric_Override(s StateTransitionMetric) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		nil, // no parameters
		s,
	)
}

// Return the given named metric for the service's state transition metrics.
// Experimental.
func StateTransitionMetric_Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of available state transitions per second.
// Experimental.
func StateTransitionMetric_MetricConsumedCapacity(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		"metricConsumedCapacity",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of available state transitions.
// Experimental.
func StateTransitionMetric_MetricProvisionedBucketSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		"metricProvisionedBucketSize",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the provisioned steady-state execution rate.
// Experimental.
func StateTransitionMetric_MetricProvisionedRefillRate(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		"metricProvisionedRefillRate",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of throttled state transitions.
// Experimental.
func StateTransitionMetric_MetricThrottledEvents(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		"metricThrottledEvents",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Properties that define what kind of task should be created.
// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
type StepFunctionsTaskConfig struct {
	// The resource that represents the work to be executed.
	//
	// Either the ARN of a Lambda Function or Activity, or a special
	// ARN.
	// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
	ResourceArn *string `json:"resourceArn"`
	// Maximum time between heart beats.
	//
	// If the time between heart beats takes longer than this, a 'Timeout' error is raised.
	//
	// This is only relevant when using an Activity type as resource.
	// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
	Heartbeat awscdk.Duration `json:"heartbeat"`
	// The dimensions to attach to metrics.
	// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
	MetricDimensions *map[string]interface{} `json:"metricDimensions"`
	// Prefix for plural metric names of activity actions.
	// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
	MetricPrefixPlural *string `json:"metricPrefixPlural"`
	// Prefix for singular metric names of activity actions.
	// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
	MetricPrefixSingular *string `json:"metricPrefixSingular"`
	// Parameters pass a collection of key-value pairs, either static values or JSONPath expressions that select from the input.
	//
	// The meaning of these parameters is task-dependent.
	//
	// Its values will be merged with the `parameters` property which is configured directly
	// on the Task state.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-parameters
	//
	// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
	Parameters *map[string]interface{} `json:"parameters"`
	// Additional policy statements to add to the execution role.
	// Deprecated: used by `IStepFunctionsTask`. `IStepFunctionsTask` is deprecated and replaced by `TaskStateBase`.
	PolicyStatements *[]awsiam.PolicyStatement `json:"policyStatements"`
}

// Define a Succeed state in the state machine.
//
// Reaching a Succeed state terminates the state execution in success.
// Experimental.
type Succeed interface {
	State
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Succeed
type jsiiProxy_Succeed struct {
	jsiiProxy_State
}

func (j *jsiiProxy_Succeed) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Succeed) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewSucceed(scope constructs.Construct, id *string, props *SucceedProps) Succeed {
	_init_.Initialize()

	j := jsiiProxy_Succeed{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Succeed",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSucceed_Override(s Succeed, scope constructs.Construct, id *string, props *SucceedProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Succeed",
		[]interface{}{scope, id, props},
		s,
	)
}

func (j *jsiiProxy_Succeed) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Succeed) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func Succeed_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Succeed",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func Succeed_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Succeed",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func Succeed_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Succeed",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Succeed_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Succeed",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func Succeed_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Succeed",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (s *jsiiProxy_Succeed) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a choice branch to this state.
// Experimental.
func (s *jsiiProxy_Succeed) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		s,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (s *jsiiProxy_Succeed) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (s *jsiiProxy_Succeed) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		s,
		"addPrefix",
		[]interface{}{x},
	)
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (s *jsiiProxy_Succeed) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (s *jsiiProxy_Succeed) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		s,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (s *jsiiProxy_Succeed) MakeNext(next State) {
	_jsii_.InvokeVoid(
		s,
		"makeNext",
		[]interface{}{next},
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
func (s *jsiiProxy_Succeed) OnPrepare() {
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
func (s *jsiiProxy_Succeed) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_Succeed) OnValidate() *[]*string {
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
func (s *jsiiProxy_Succeed) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (s *jsiiProxy_Succeed) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (s *jsiiProxy_Succeed) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (s *jsiiProxy_Succeed) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (s *jsiiProxy_Succeed) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (s *jsiiProxy_Succeed) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (s *jsiiProxy_Succeed) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (s *jsiiProxy_Succeed) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"renderRetryCatch",
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
func (s *jsiiProxy_Succeed) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (s *jsiiProxy_Succeed) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_Succeed) ToString() *string {
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
func (s *jsiiProxy_Succeed) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (s *jsiiProxy_Succeed) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		s,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a Succeed state.
// Experimental.
type SucceedProps struct {
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
}

// Define a Task state in the state machine.
//
// Reaching a Task state causes some work to be executed, represented by the
// Task's resource property. Task constructs represent a generic Amazon
// States Language Task.
//
// For some resource types, more specific subclasses of Task may be available
// which are more convenient to use.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
type Task interface {
	State
	INextable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddCatch(handler IChainable, props *CatchProps) Task
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	AddRetry(props *RetryProps) Task
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHeartbeatTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRunTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricScheduled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricScheduleTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Task
type jsiiProxy_Task struct {
	jsiiProxy_State
	jsiiProxy_INextable
}

func (j *jsiiProxy_Task) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Task) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func NewTask(scope constructs.Construct, id *string, props *TaskProps) Task {
	_init_.Initialize()

	j := jsiiProxy_Task{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Task",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func NewTask_Override(t Task, scope constructs.Construct, id *string, props *TaskProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Task",
		[]interface{}{scope, id, props},
		t,
	)
}

func (j *jsiiProxy_Task) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Task) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func Task_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Task",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func Task_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Task",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func Task_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Task",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func Task_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Task",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func Task_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Task",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a recovery handler for this state.
//
// When a particular error occurs, execution will continue at the error
// handler instead of failing the state machine execution.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) AddCatch(handler IChainable, props *CatchProps) Task {
	var returns Task

	_jsii_.Invoke(
		t,
		"addCatch",
		[]interface{}{handler, props},
		&returns,
	)

	return returns
}

// Add a choice branch to this state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		t,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		t,
		"addPrefix",
		[]interface{}{x},
	)
}

// Add retry configuration for this state.
//
// This controls if and how the execution will be retried if a particular
// error occurs.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) AddRetry(props *RetryProps) Task {
	var returns Task

	_jsii_.Invoke(
		t,
		"addRetry",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		t,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MakeNext(next State) {
	_jsii_.InvokeVoid(
		t,
		"makeNext",
		[]interface{}{next},
	)
}

// Return the given named metric for this Task.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity fails.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times the heartbeat times out for this activity.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricHeartbeatTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricHeartbeatTimedOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, between the time the Task starts and the time it closes.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricRunTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricRunTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity is scheduled.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricScheduled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricScheduled",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, for which the activity stays in the schedule state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricScheduleTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricScheduleTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity is started.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricStarted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity succeeds.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSucceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, between the time the activity is scheduled and the time it closes.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity times out.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricTimedOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Continue normal execution with the given state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		t,
		"next",
		[]interface{}{next},
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
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) OnPrepare() {
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
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) OnValidate() *[]*string {
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
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) Prepare() {
	_jsii_.InvokeVoid(
		t,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderRetryCatch",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		t,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) ToString() *string {
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
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
func (t *jsiiProxy_Task) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Type union for task classes that accept multiple types of payload.
// Experimental.
type TaskInput interface {
	Type() InputType
	Value() interface{}
}

// The jsii proxy struct for TaskInput
type jsiiProxy_TaskInput struct {
	_ byte // padding
}

func (j *jsiiProxy_TaskInput) Type() InputType {
	var returns InputType
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskInput) Value() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


// Use a part of the task context as task input.
//
// Use this when you want to use a subobject or string from
// the current task context as complete payload
// to a task.
// Deprecated: Use `fromJsonPathAt`.
func TaskInput_FromContextAt(path *string) TaskInput {
	_init_.Initialize()

	var returns TaskInput

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskInput",
		"fromContextAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Use a part of the execution data as task input.
//
// Use this when you want to use a subobject or string from
// the current state machine execution as complete payload
// to a task.
// Deprecated: Use `fromJsonPathAt`.
func TaskInput_FromDataAt(path *string) TaskInput {
	_init_.Initialize()

	var returns TaskInput

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskInput",
		"fromDataAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Use a part of the execution data or task context as task input.
//
// Use this when you want to use a subobject or string from
// the current state machine execution or the current task context
// as complete payload to a task.
// Experimental.
func TaskInput_FromJsonPathAt(path *string) TaskInput {
	_init_.Initialize()

	var returns TaskInput

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskInput",
		"fromJsonPathAt",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Use an object as task input.
//
// This object may contain JSON path fields as object values, if desired.
// Experimental.
func TaskInput_FromObject(obj *map[string]interface{}) TaskInput {
	_init_.Initialize()

	var returns TaskInput

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskInput",
		"fromObject",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Use a literal string as task input.
//
// This might be a JSON-encoded object, or just a text.
// Experimental.
func TaskInput_FromText(text *string) TaskInput {
	_init_.Initialize()

	var returns TaskInput

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskInput",
		"fromText",
		[]interface{}{text},
		&returns,
	)

	return returns
}

// Task Metrics.
// Experimental.
type TaskMetricsConfig struct {
	// The dimensions to attach to metrics.
	// Experimental.
	MetricDimensions *map[string]interface{} `json:"metricDimensions"`
	// Prefix for plural metric names of activity actions.
	// Experimental.
	MetricPrefixPlural *string `json:"metricPrefixPlural"`
	// Prefix for singular metric names of activity actions.
	// Experimental.
	MetricPrefixSingular *string `json:"metricPrefixSingular"`
}

// Props that are common to all tasks.
// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
type TaskProps struct {
	// Actual task to be invoked in this workflow.
	// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
	Task IStepFunctionsTask `json:"task"`
	// An optional description for this state.
	// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
	Comment *string `json:"comment"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
	InputPath *string `json:"inputPath"`
	// JSONPath expression to select part of the state to be the output to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
	OutputPath *string `json:"outputPath"`
	// Parameters to invoke the task with.
	//
	// It is not recommended to use this field. The object that is passed in
	// the `task` property will take care of returning the right values for the
	// `Parameters` field in the Step Functions definition.
	//
	// The various classes that implement `IStepFunctionsTask` will take a
	// properties which make sense for the task type. For example, for
	// `InvokeFunction` the field that populates the `parameters` field will be
	// called `payload`, and for the `PublishToTopic` the `parameters` field
	// will be populated via a combination of the referenced topic, subject and
	// message.
	//
	// If passed anyway, the keys in this map will override the parameters
	// returned by the task object.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-parameters
	//
	// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
	Parameters *map[string]interface{} `json:"parameters"`
	// JSONPath expression to indicate where to inject the state's output.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the state's
	// input to become its output.
	// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
	ResultPath *string `json:"resultPath"`
	// Maximum run time of this state.
	//
	// If the state takes longer than this amount of time to complete, a 'Timeout' error is raised.
	// Deprecated: - replaced by service integration specific classes (i.e. LambdaInvoke, SnsPublish)
	Timeout awscdk.Duration `json:"timeout"`
}

// Define a Task state in the state machine.
//
// Reaching a Task state causes some work to be executed, represented by the
// Task's resource property. Task constructs represent a generic Amazon
// States Language Task.
//
// For some resource types, more specific subclasses of Task may be available
// which are more convenient to use.
// Experimental.
type TaskStateBase interface {
	State
	INextable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	TaskMetrics() *TaskMetricsConfig
	TaskPolicies() *[]awsiam.PolicyStatement
	AddBranch(branch StateGraph)
	AddCatch(handler IChainable, props *CatchProps) TaskStateBase
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	AddRetry(props *RetryProps) TaskStateBase
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHeartbeatTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRunTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricScheduled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricScheduleTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for TaskStateBase
type jsiiProxy_TaskStateBase struct {
	jsiiProxy_State
	jsiiProxy_INextable
}

func (j *jsiiProxy_TaskStateBase) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) TaskMetrics() *TaskMetricsConfig {
	var returns *TaskMetricsConfig
	_jsii_.Get(
		j,
		"taskMetrics",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskStateBase) TaskPolicies() *[]awsiam.PolicyStatement {
	var returns *[]awsiam.PolicyStatement
	_jsii_.Get(
		j,
		"taskPolicies",
		&returns,
	)
	return returns
}


// Experimental.
func NewTaskStateBase_Override(t TaskStateBase, scope constructs.Construct, id *string, props *TaskStateBaseProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.TaskStateBase",
		[]interface{}{scope, id, props},
		t,
	)
}

func (j *jsiiProxy_TaskStateBase) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_TaskStateBase) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func TaskStateBase_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskStateBase",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func TaskStateBase_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskStateBase",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func TaskStateBase_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskStateBase",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func TaskStateBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.TaskStateBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func TaskStateBase_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.TaskStateBase",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a recovery handler for this state.
//
// When a particular error occurs, execution will continue at the error
// handler instead of failing the state machine execution.
// Experimental.
func (t *jsiiProxy_TaskStateBase) AddCatch(handler IChainable, props *CatchProps) TaskStateBase {
	var returns TaskStateBase

	_jsii_.Invoke(
		t,
		"addCatch",
		[]interface{}{handler, props},
		&returns,
	)

	return returns
}

// Add a choice branch to this state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		t,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		t,
		"addPrefix",
		[]interface{}{x},
	)
}

// Add retry configuration for this state.
//
// This controls if and how the execution will be retried if a particular
// error occurs.
// Experimental.
func (t *jsiiProxy_TaskStateBase) AddRetry(props *RetryProps) TaskStateBase {
	var returns TaskStateBase

	_jsii_.Invoke(
		t,
		"addRetry",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (t *jsiiProxy_TaskStateBase) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		t,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MakeNext(next State) {
	_jsii_.InvokeVoid(
		t,
		"makeNext",
		[]interface{}{next},
	)
}

// Return the given named metric for this Task.
// Experimental.
func (t *jsiiProxy_TaskStateBase) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity fails.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times the heartbeat times out for this activity.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricHeartbeatTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricHeartbeatTimedOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, between the time the Task starts and the time it closes.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricRunTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricRunTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity is scheduled.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricScheduled(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricScheduled",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, for which the activity stays in the schedule state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricScheduleTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricScheduleTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity is started.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricStarted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricStarted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity succeeds.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricSucceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSucceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The interval, in milliseconds, between the time the activity is scheduled and the time it closes.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of times this activity times out.
// Experimental.
func (t *jsiiProxy_TaskStateBase) MetricTimedOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricTimedOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Continue normal execution with the given state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		t,
		"next",
		[]interface{}{next},
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
func (t *jsiiProxy_TaskStateBase) OnPrepare() {
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
func (t *jsiiProxy_TaskStateBase) OnSynthesize(session constructs.ISynthesisSession) {
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
func (t *jsiiProxy_TaskStateBase) OnValidate() *[]*string {
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
func (t *jsiiProxy_TaskStateBase) Prepare() {
	_jsii_.InvokeVoid(
		t,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (t *jsiiProxy_TaskStateBase) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (t *jsiiProxy_TaskStateBase) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (t *jsiiProxy_TaskStateBase) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (t *jsiiProxy_TaskStateBase) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (t *jsiiProxy_TaskStateBase) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (t *jsiiProxy_TaskStateBase) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (t *jsiiProxy_TaskStateBase) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderRetryCatch",
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
func (t *jsiiProxy_TaskStateBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (t *jsiiProxy_TaskStateBase) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		t,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TaskStateBase) ToString() *string {
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
func (t *jsiiProxy_TaskStateBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (t *jsiiProxy_TaskStateBase) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		t,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Props that are common to all tasks.
// Experimental.
type TaskStateBaseProps struct {
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
	// Timeout for the heartbeat.
	// Experimental.
	Heartbeat awscdk.Duration `json:"heartbeat"`
	// JSONPath expression to select part of the state to be the input to this state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// input to be the empty object {}.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// AWS Step Functions integrates with services directly in the Amazon States Language.
	//
	// You can control these AWS services using service integration patterns
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/connect-to-resource.html#connect-wait-token
	//
	// Experimental.
	IntegrationPattern IntegrationPattern `json:"integrationPattern"`
	// JSONPath expression to select select a portion of the state output to pass to the next state.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the effective
	// output to be the empty object {}.
	// Experimental.
	OutputPath *string `json:"outputPath"`
	// JSONPath expression to indicate where to inject the state's output.
	//
	// May also be the special value JsonPath.DISCARD, which will cause the state's
	// input to become its output.
	// Experimental.
	ResultPath *string `json:"resultPath"`
	// The JSON that will replace the state's raw result and become the effective result before ResultPath is applied.
	//
	// You can use ResultSelector to create a payload with values that are static
	// or selected from the state's raw result.
	// See: https://docs.aws.amazon.com/step-functions/latest/dg/input-output-inputpath-params.html#input-output-resultselector
	//
	// Experimental.
	ResultSelector *map[string]interface{} `json:"resultSelector"`
	// Timeout for the state machine.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
}

// Define a Wait state in the state machine.
//
// A Wait state can be used to delay execution of the state machine for a while.
// Experimental.
type Wait interface {
	State
	INextable
	Branches() *[]StateGraph
	Comment() *string
	DefaultChoice() State
	SetDefaultChoice(val State)
	EndStates() *[]INextable
	Id() *string
	InputPath() *string
	Iteration() StateGraph
	SetIteration(val StateGraph)
	Node() awscdk.ConstructNode
	OutputPath() *string
	Parameters() *map[string]interface{}
	ResultPath() *string
	ResultSelector() *map[string]interface{}
	StartState() State
	StateId() *string
	AddBranch(branch StateGraph)
	AddChoice(condition Condition, next State)
	AddIterator(iteration StateGraph)
	AddPrefix(x *string)
	BindToGraph(graph StateGraph)
	MakeDefault(def State)
	MakeNext(next State)
	Next(next IChainable) Chain
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderBranches() interface{}
	RenderChoices() interface{}
	RenderInputOutput() interface{}
	RenderIterator() interface{}
	RenderNextEnd() interface{}
	RenderResultSelector() interface{}
	RenderRetryCatch() interface{}
	Synthesize(session awscdk.ISynthesisSession)
	ToStateJson() *map[string]interface{}
	ToString() *string
	Validate() *[]*string
	WhenBoundToGraph(graph StateGraph)
}

// The jsii proxy struct for Wait
type jsiiProxy_Wait struct {
	jsiiProxy_State
	jsiiProxy_INextable
}

func (j *jsiiProxy_Wait) Branches() *[]StateGraph {
	var returns *[]StateGraph
	_jsii_.Get(
		j,
		"branches",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) Comment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"comment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) DefaultChoice() State {
	var returns State
	_jsii_.Get(
		j,
		"defaultChoice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) EndStates() *[]INextable {
	var returns *[]INextable
	_jsii_.Get(
		j,
		"endStates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) InputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"inputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) Iteration() StateGraph {
	var returns StateGraph
	_jsii_.Get(
		j,
		"iteration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) OutputPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outputPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) Parameters() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) ResultPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resultPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) ResultSelector() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"resultSelector",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) StartState() State {
	var returns State
	_jsii_.Get(
		j,
		"startState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Wait) StateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stateId",
		&returns,
	)
	return returns
}


// Experimental.
func NewWait(scope constructs.Construct, id *string, props *WaitProps) Wait {
	_init_.Initialize()

	j := jsiiProxy_Wait{}

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Wait",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewWait_Override(w Wait, scope constructs.Construct, id *string, props *WaitProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_stepfunctions.Wait",
		[]interface{}{scope, id, props},
		w,
	)
}

func (j *jsiiProxy_Wait) SetDefaultChoice(val State) {
	_jsii_.Set(
		j,
		"defaultChoice",
		val,
	)
}

func (j *jsiiProxy_Wait) SetIteration(val StateGraph) {
	_jsii_.Set(
		j,
		"iteration",
		val,
	)
}

// Return only the states that allow chaining from an array of states.
// Experimental.
func Wait_FilterNextables(states *[]State) *[]INextable {
	_init_.Initialize()

	var returns *[]INextable

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Wait",
		"filterNextables",
		[]interface{}{states},
		&returns,
	)

	return returns
}

// Find the set of end states states reachable through transitions from the given start state.
// Experimental.
func Wait_FindReachableEndStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Wait",
		"findReachableEndStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Find the set of states reachable through transitions from the given start state.
//
// This does not retrieve states from within sub-graphs, such as states within a Parallel state's branch.
// Experimental.
func Wait_FindReachableStates(start State, options *FindStateOptions) *[]State {
	_init_.Initialize()

	var returns *[]State

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Wait",
		"findReachableStates",
		[]interface{}{start, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Wait_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.Wait",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a prefix to the stateId of all States found in a construct tree.
// Experimental.
func Wait_PrefixStates(root constructs.IConstruct, prefix *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_stepfunctions.Wait",
		"prefixStates",
		[]interface{}{root, prefix},
	)
}

// Add a paralle branch to this state.
// Experimental.
func (w *jsiiProxy_Wait) AddBranch(branch StateGraph) {
	_jsii_.InvokeVoid(
		w,
		"addBranch",
		[]interface{}{branch},
	)
}

// Add a choice branch to this state.
// Experimental.
func (w *jsiiProxy_Wait) AddChoice(condition Condition, next State) {
	_jsii_.InvokeVoid(
		w,
		"addChoice",
		[]interface{}{condition, next},
	)
}

// Add a map iterator to this state.
// Experimental.
func (w *jsiiProxy_Wait) AddIterator(iteration StateGraph) {
	_jsii_.InvokeVoid(
		w,
		"addIterator",
		[]interface{}{iteration},
	)
}

// Add a prefix to the stateId of this state.
// Experimental.
func (w *jsiiProxy_Wait) AddPrefix(x *string) {
	_jsii_.InvokeVoid(
		w,
		"addPrefix",
		[]interface{}{x},
	)
}

// Register this state as part of the given graph.
//
// Don't call this. It will be called automatically when you work
// with states normally.
// Experimental.
func (w *jsiiProxy_Wait) BindToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		w,
		"bindToGraph",
		[]interface{}{graph},
	)
}

// Make the indicated state the default choice transition of this state.
// Experimental.
func (w *jsiiProxy_Wait) MakeDefault(def State) {
	_jsii_.InvokeVoid(
		w,
		"makeDefault",
		[]interface{}{def},
	)
}

// Make the indicated state the default transition of this state.
// Experimental.
func (w *jsiiProxy_Wait) MakeNext(next State) {
	_jsii_.InvokeVoid(
		w,
		"makeNext",
		[]interface{}{next},
	)
}

// Continue normal execution with the given state.
// Experimental.
func (w *jsiiProxy_Wait) Next(next IChainable) Chain {
	var returns Chain

	_jsii_.Invoke(
		w,
		"next",
		[]interface{}{next},
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
func (w *jsiiProxy_Wait) OnPrepare() {
	_jsii_.InvokeVoid(
		w,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (w *jsiiProxy_Wait) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		w,
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
func (w *jsiiProxy_Wait) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		w,
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
func (w *jsiiProxy_Wait) Prepare() {
	_jsii_.InvokeVoid(
		w,
		"prepare",
		nil, // no parameters
	)
}

// Render parallel branches in ASL JSON format.
// Experimental.
func (w *jsiiProxy_Wait) RenderBranches() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"renderBranches",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the choices in ASL JSON format.
// Experimental.
func (w *jsiiProxy_Wait) RenderChoices() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"renderChoices",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render InputPath/Parameters/OutputPath in ASL JSON format.
// Experimental.
func (w *jsiiProxy_Wait) RenderInputOutput() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"renderInputOutput",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render map iterator in ASL JSON format.
// Experimental.
func (w *jsiiProxy_Wait) RenderIterator() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"renderIterator",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render the default next state in ASL JSON format.
// Experimental.
func (w *jsiiProxy_Wait) RenderNextEnd() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"renderNextEnd",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render ResultSelector in ASL JSON format.
// Experimental.
func (w *jsiiProxy_Wait) RenderResultSelector() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"renderResultSelector",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Render error recovery options in ASL JSON format.
// Experimental.
func (w *jsiiProxy_Wait) RenderRetryCatch() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"renderRetryCatch",
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
func (w *jsiiProxy_Wait) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		w,
		"synthesize",
		[]interface{}{session},
	)
}

// Return the Amazon States Language object for this state.
// Experimental.
func (w *jsiiProxy_Wait) ToStateJson() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		w,
		"toStateJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (w *jsiiProxy_Wait) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		w,
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
func (w *jsiiProxy_Wait) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		w,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Called whenever this state is bound to a graph.
//
// Can be overridden by subclasses.
// Experimental.
func (w *jsiiProxy_Wait) WhenBoundToGraph(graph StateGraph) {
	_jsii_.InvokeVoid(
		w,
		"whenBoundToGraph",
		[]interface{}{graph},
	)
}

// Properties for defining a Wait state.
// Experimental.
type WaitProps struct {
	// Wait duration.
	// Experimental.
	Time WaitTime `json:"time"`
	// An optional description for this state.
	// Experimental.
	Comment *string `json:"comment"`
}

// Represents the Wait state which delays a state machine from continuing for a specified time.
// See: https://docs.aws.amazon.com/step-functions/latest/dg/amazon-states-language-wait-state.html
//
// Experimental.
type WaitTime interface {
}

// The jsii proxy struct for WaitTime
type jsiiProxy_WaitTime struct {
	_ byte // padding
}

// Wait a fixed amount of time.
// Experimental.
func WaitTime_Duration(duration awscdk.Duration) WaitTime {
	_init_.Initialize()

	var returns WaitTime

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.WaitTime",
		"duration",
		[]interface{}{duration},
		&returns,
	)

	return returns
}

// Wait for a number of seconds stored in the state object.
//
// TODO: EXAMPLE
//
// Experimental.
func WaitTime_SecondsPath(path *string) WaitTime {
	_init_.Initialize()

	var returns WaitTime

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.WaitTime",
		"secondsPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Wait until the given ISO8601 timestamp.
//
// TODO: EXAMPLE
//
// Experimental.
func WaitTime_Timestamp(timestamp *string) WaitTime {
	_init_.Initialize()

	var returns WaitTime

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.WaitTime",
		"timestamp",
		[]interface{}{timestamp},
		&returns,
	)

	return returns
}

// Wait until a timestamp found in the state object.
//
// TODO: EXAMPLE
//
// Experimental.
func WaitTime_TimestampPath(path *string) WaitTime {
	_init_.Initialize()

	var returns WaitTime

	_jsii_.StaticInvoke(
		"monocdk.aws_stepfunctions.WaitTime",
		"timestampPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}


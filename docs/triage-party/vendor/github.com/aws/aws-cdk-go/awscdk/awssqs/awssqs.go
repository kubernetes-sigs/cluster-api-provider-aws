package awssqs

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awssqs/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// A CloudFormation `AWS::SQS::Queue`.
type CfnQueue interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrQueueName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ContentBasedDeduplication() interface{}
	SetContentBasedDeduplication(val interface{})
	CreationStack() *[]*string
	DeduplicationScope() *string
	SetDeduplicationScope(val *string)
	DelaySeconds() *float64
	SetDelaySeconds(val *float64)
	FifoQueue() interface{}
	SetFifoQueue(val interface{})
	FifoThroughputLimit() *string
	SetFifoThroughputLimit(val *string)
	KmsDataKeyReusePeriodSeconds() *float64
	SetKmsDataKeyReusePeriodSeconds(val *float64)
	KmsMasterKeyId() *string
	SetKmsMasterKeyId(val *string)
	LogicalId() *string
	MaximumMessageSize() *float64
	SetMaximumMessageSize(val *float64)
	MessageRetentionPeriod() *float64
	SetMessageRetentionPeriod(val *float64)
	Node() awscdk.ConstructNode
	QueueName() *string
	SetQueueName(val *string)
	ReceiveMessageWaitTimeSeconds() *float64
	SetReceiveMessageWaitTimeSeconds(val *float64)
	RedrivePolicy() interface{}
	SetRedrivePolicy(val interface{})
	Ref() *string
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	UpdatedProperites() *map[string]interface{}
	VisibilityTimeout() *float64
	SetVisibilityTimeout(val *float64)
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

// The jsii proxy struct for CfnQueue
type jsiiProxy_CfnQueue struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnQueue) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) AttrQueueName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrQueueName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) ContentBasedDeduplication() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"contentBasedDeduplication",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) DeduplicationScope() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deduplicationScope",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) DelaySeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"delaySeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) FifoQueue() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"fifoQueue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) FifoThroughputLimit() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fifoThroughputLimit",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) KmsDataKeyReusePeriodSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"kmsDataKeyReusePeriodSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) KmsMasterKeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kmsMasterKeyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) MaximumMessageSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maximumMessageSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) MessageRetentionPeriod() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"messageRetentionPeriod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) QueueName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) ReceiveMessageWaitTimeSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"receiveMessageWaitTimeSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) RedrivePolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"redrivePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueue) VisibilityTimeout() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"visibilityTimeout",
		&returns,
	)
	return returns
}


// Create a new `AWS::SQS::Queue`.
func NewCfnQueue(scope awscdk.Construct, id *string, props *CfnQueueProps) CfnQueue {
	_init_.Initialize()

	j := jsiiProxy_CfnQueue{}

	_jsii_.Create(
		"monocdk.aws_sqs.CfnQueue",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SQS::Queue`.
func NewCfnQueue_Override(c CfnQueue, scope awscdk.Construct, id *string, props *CfnQueueProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sqs.CfnQueue",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnQueue) SetContentBasedDeduplication(val interface{}) {
	_jsii_.Set(
		j,
		"contentBasedDeduplication",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetDeduplicationScope(val *string) {
	_jsii_.Set(
		j,
		"deduplicationScope",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetDelaySeconds(val *float64) {
	_jsii_.Set(
		j,
		"delaySeconds",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetFifoQueue(val interface{}) {
	_jsii_.Set(
		j,
		"fifoQueue",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetFifoThroughputLimit(val *string) {
	_jsii_.Set(
		j,
		"fifoThroughputLimit",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetKmsDataKeyReusePeriodSeconds(val *float64) {
	_jsii_.Set(
		j,
		"kmsDataKeyReusePeriodSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetKmsMasterKeyId(val *string) {
	_jsii_.Set(
		j,
		"kmsMasterKeyId",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetMaximumMessageSize(val *float64) {
	_jsii_.Set(
		j,
		"maximumMessageSize",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetMessageRetentionPeriod(val *float64) {
	_jsii_.Set(
		j,
		"messageRetentionPeriod",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetQueueName(val *string) {
	_jsii_.Set(
		j,
		"queueName",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetReceiveMessageWaitTimeSeconds(val *float64) {
	_jsii_.Set(
		j,
		"receiveMessageWaitTimeSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetRedrivePolicy(val interface{}) {
	_jsii_.Set(
		j,
		"redrivePolicy",
		val,
	)
}

func (j *jsiiProxy_CfnQueue) SetVisibilityTimeout(val *float64) {
	_jsii_.Set(
		j,
		"visibilityTimeout",
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
func CfnQueue_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.CfnQueue",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnQueue_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.CfnQueue",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnQueue_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.CfnQueue",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnQueue_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_sqs.CfnQueue",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnQueue) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnQueue) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnQueue) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnQueue) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnQueue) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnQueue) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnQueue) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnQueue) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnQueue) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnQueue) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnQueue) OnPrepare() {
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
func (c *jsiiProxy_CfnQueue) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnQueue) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnQueue) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnQueue) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnQueue) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnQueue) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnQueue) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnQueue) ToString() *string {
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
func (c *jsiiProxy_CfnQueue) Validate() *[]*string {
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
func (c *jsiiProxy_CfnQueue) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// A CloudFormation `AWS::SQS::QueuePolicy`.
type CfnQueuePolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	PolicyDocument() interface{}
	SetPolicyDocument(val interface{})
	Queues() *[]*string
	SetQueues(val *[]*string)
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

// The jsii proxy struct for CfnQueuePolicy
type jsiiProxy_CfnQueuePolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnQueuePolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) PolicyDocument() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policyDocument",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) Queues() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"queues",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnQueuePolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SQS::QueuePolicy`.
func NewCfnQueuePolicy(scope awscdk.Construct, id *string, props *CfnQueuePolicyProps) CfnQueuePolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnQueuePolicy{}

	_jsii_.Create(
		"monocdk.aws_sqs.CfnQueuePolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SQS::QueuePolicy`.
func NewCfnQueuePolicy_Override(c CfnQueuePolicy, scope awscdk.Construct, id *string, props *CfnQueuePolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sqs.CfnQueuePolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnQueuePolicy) SetPolicyDocument(val interface{}) {
	_jsii_.Set(
		j,
		"policyDocument",
		val,
	)
}

func (j *jsiiProxy_CfnQueuePolicy) SetQueues(val *[]*string) {
	_jsii_.Set(
		j,
		"queues",
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
func CfnQueuePolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.CfnQueuePolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnQueuePolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.CfnQueuePolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnQueuePolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.CfnQueuePolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnQueuePolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_sqs.CfnQueuePolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnQueuePolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnQueuePolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnQueuePolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnQueuePolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnQueuePolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnQueuePolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnQueuePolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnQueuePolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnQueuePolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnQueuePolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnQueuePolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnQueuePolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnQueuePolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnQueuePolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnQueuePolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnQueuePolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnQueuePolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnQueuePolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnQueuePolicy) ToString() *string {
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
func (c *jsiiProxy_CfnQueuePolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnQueuePolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::SQS::QueuePolicy`.
type CfnQueuePolicyProps struct {
	// `AWS::SQS::QueuePolicy.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
	// `AWS::SQS::QueuePolicy.Queues`.
	Queues *[]*string `json:"queues"`
}

// Properties for defining a `AWS::SQS::Queue`.
type CfnQueueProps struct {
	// `AWS::SQS::Queue.ContentBasedDeduplication`.
	ContentBasedDeduplication interface{} `json:"contentBasedDeduplication"`
	// `AWS::SQS::Queue.DeduplicationScope`.
	DeduplicationScope *string `json:"deduplicationScope"`
	// `AWS::SQS::Queue.DelaySeconds`.
	DelaySeconds *float64 `json:"delaySeconds"`
	// `AWS::SQS::Queue.FifoQueue`.
	FifoQueue interface{} `json:"fifoQueue"`
	// `AWS::SQS::Queue.FifoThroughputLimit`.
	FifoThroughputLimit *string `json:"fifoThroughputLimit"`
	// `AWS::SQS::Queue.KmsDataKeyReusePeriodSeconds`.
	KmsDataKeyReusePeriodSeconds *float64 `json:"kmsDataKeyReusePeriodSeconds"`
	// `AWS::SQS::Queue.KmsMasterKeyId`.
	KmsMasterKeyId *string `json:"kmsMasterKeyId"`
	// `AWS::SQS::Queue.MaximumMessageSize`.
	MaximumMessageSize *float64 `json:"maximumMessageSize"`
	// `AWS::SQS::Queue.MessageRetentionPeriod`.
	MessageRetentionPeriod *float64 `json:"messageRetentionPeriod"`
	// `AWS::SQS::Queue.QueueName`.
	QueueName *string `json:"queueName"`
	// `AWS::SQS::Queue.ReceiveMessageWaitTimeSeconds`.
	ReceiveMessageWaitTimeSeconds *float64 `json:"receiveMessageWaitTimeSeconds"`
	// `AWS::SQS::Queue.RedrivePolicy`.
	RedrivePolicy interface{} `json:"redrivePolicy"`
	// `AWS::SQS::Queue.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::SQS::Queue.VisibilityTimeout`.
	VisibilityTimeout *float64 `json:"visibilityTimeout"`
}

// Dead letter queue settings.
// Experimental.
type DeadLetterQueue struct {
	// The number of times a message can be unsuccesfully dequeued before being moved to the dead-letter queue.
	// Experimental.
	MaxReceiveCount *float64 `json:"maxReceiveCount"`
	// The dead-letter queue to which Amazon SQS moves messages after the value of maxReceiveCount is exceeded.
	// Experimental.
	Queue IQueue `json:"queue"`
}

// Represents an SQS queue.
// Experimental.
type IQueue interface {
	awscdk.IResource
	// Adds a statement to the IAM resource policy associated with this queue.
	//
	// If this queue was created in this stack (`new Queue`), a queue policy
	// will be automatically created upon the first call to `addToPolicy`. If
	// the queue is imported (`Queue.import`), then this is a no-op.
	// Experimental.
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	// Grant the actions defined in queueActions to the identity Principal given on this SQS queue resource.
	// Experimental.
	Grant(grantee awsiam.IGrantable, queueActions ...*string) awsiam.Grant
	// Grant permissions to consume messages from a queue.
	//
	// This will grant the following permissions:
	//
	//    - sqs:ChangeMessageVisibility
	//    - sqs:DeleteMessage
	//    - sqs:ReceiveMessage
	//    - sqs:GetQueueAttributes
	//    - sqs:GetQueueUrl
	// Experimental.
	GrantConsumeMessages(grantee awsiam.IGrantable) awsiam.Grant
	// Grant an IAM principal permissions to purge all messages from the queue.
	//
	// This will grant the following permissions:
	//
	//   - sqs:PurgeQueue
	//   - sqs:GetQueueAttributes
	//   - sqs:GetQueueUrl
	// Experimental.
	GrantPurge(grantee awsiam.IGrantable) awsiam.Grant
	// Grant access to send messages to a queue to the given identity.
	//
	// This will grant the following permissions:
	//
	//   - sqs:SendMessage
	//   - sqs:GetQueueAttributes
	//   - sqs:GetQueueUrl
	// Experimental.
	GrantSendMessages(grantee awsiam.IGrantable) awsiam.Grant
	// Return the given named metric for this Queue.
	// Experimental.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The approximate age of the oldest non-deleted message in the queue.
	//
	// Maximum over 5 minutes
	// Experimental.
	MetricApproximateAgeOfOldestMessage(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages in the queue that are delayed and not available for reading immediately.
	//
	// Maximum over 5 minutes
	// Experimental.
	MetricApproximateNumberOfMessagesDelayed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages that are in flight.
	//
	// Maximum over 5 minutes
	// Experimental.
	MetricApproximateNumberOfMessagesNotVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages available for retrieval from the queue.
	//
	// Maximum over 5 minutes
	// Experimental.
	MetricApproximateNumberOfMessagesVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of ReceiveMessage API calls that did not return a message.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfEmptyReceives(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages deleted from the queue.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfMessagesDeleted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages returned by calls to the ReceiveMessage action.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfMessagesReceived(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages added to a queue.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfMessagesSent(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The size of messages added to a queue.
	//
	// Average over 5 minutes
	// Experimental.
	MetricSentMessageSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// If this queue is server-side encrypted, this is the KMS encryption key.
	// Experimental.
	EncryptionMasterKey() awskms.IKey
	// Whether this queue is an Amazon SQS FIFO queue.
	//
	// If false, this is a standard queue.
	// Experimental.
	Fifo() *bool
	// The ARN of this queue.
	// Experimental.
	QueueArn() *string
	// The name of this queue.
	// Experimental.
	QueueName() *string
	// The URL of this queue.
	// Experimental.
	QueueUrl() *string
}

// The jsii proxy for IQueue
type jsiiProxy_IQueue struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IQueue) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		i,
		"addToResourcePolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) Grant(grantee awsiam.IGrantable, queueActions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
	for _, a := range queueActions {
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

func (i *jsiiProxy_IQueue) GrantConsumeMessages(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantConsumeMessages",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) GrantPurge(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPurge",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) GrantSendMessages(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantSendMessages",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricApproximateAgeOfOldestMessage(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricApproximateAgeOfOldestMessage",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricApproximateNumberOfMessagesDelayed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricApproximateNumberOfMessagesDelayed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricApproximateNumberOfMessagesNotVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricApproximateNumberOfMessagesNotVisible",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricApproximateNumberOfMessagesVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricApproximateNumberOfMessagesVisible",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricNumberOfEmptyReceives(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfEmptyReceives",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricNumberOfMessagesDeleted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfMessagesDeleted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricNumberOfMessagesReceived(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfMessagesReceived",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricNumberOfMessagesSent(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfMessagesSent",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IQueue) MetricSentMessageSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricSentMessageSize",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IQueue) EncryptionMasterKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionMasterKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IQueue) Fifo() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"fifo",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IQueue) QueueArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IQueue) QueueName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IQueue) QueueUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueUrl",
		&returns,
	)
	return returns
}

// A new Amazon SQS queue.
// Experimental.
type Queue interface {
	QueueBase
	AutoCreatePolicy() *bool
	EncryptionMasterKey() awskms.IKey
	Env() *awscdk.ResourceEnvironment
	Fifo() *bool
	Node() awscdk.ConstructNode
	PhysicalName() *string
	QueueArn() *string
	QueueName() *string
	QueueUrl() *string
	Stack() awscdk.Stack
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantConsumeMessages(grantee awsiam.IGrantable) awsiam.Grant
	GrantPurge(grantee awsiam.IGrantable) awsiam.Grant
	GrantSendMessages(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateAgeOfOldestMessage(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateNumberOfMessagesDelayed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateNumberOfMessagesNotVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateNumberOfMessagesVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfEmptyReceives(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesDeleted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesReceived(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesSent(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSentMessageSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Queue
type jsiiProxy_Queue struct {
	jsiiProxy_QueueBase
}

func (j *jsiiProxy_Queue) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) EncryptionMasterKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionMasterKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) Fifo() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"fifo",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) QueueArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) QueueName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) QueueUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Queue) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewQueue(scope constructs.Construct, id *string, props *QueueProps) Queue {
	_init_.Initialize()

	j := jsiiProxy_Queue{}

	_jsii_.Create(
		"monocdk.aws_sqs.Queue",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewQueue_Override(q Queue, scope constructs.Construct, id *string, props *QueueProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sqs.Queue",
		[]interface{}{scope, id, props},
		q,
	)
}

// Import an existing SQS queue provided an ARN.
// Experimental.
func Queue_FromQueueArn(scope constructs.Construct, id *string, queueArn *string) IQueue {
	_init_.Initialize()

	var returns IQueue

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.Queue",
		"fromQueueArn",
		[]interface{}{scope, id, queueArn},
		&returns,
	)

	return returns
}

// Import an existing queue.
// Experimental.
func Queue_FromQueueAttributes(scope constructs.Construct, id *string, attrs *QueueAttributes) IQueue {
	_init_.Initialize()

	var returns IQueue

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.Queue",
		"fromQueueAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Queue_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.Queue",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Queue_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.Queue",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a statement to the IAM resource policy associated with this queue.
//
// If this queue was created in this stack (`new Queue`), a queue policy
// will be automatically created upon the first call to `addToPolicy`. If
// the queue is imported (`Queue.import`), then this is a no-op.
// Experimental.
func (q *jsiiProxy_Queue) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		q,
		"addToResourcePolicy",
		[]interface{}{statement},
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
func (q *jsiiProxy_Queue) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		q,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (q *jsiiProxy_Queue) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_Queue) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_Queue) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		q,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the actions defined in queueActions to the identity Principal given on this SQS queue resource.
// Experimental.
func (q *jsiiProxy_Queue) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Grant permissions to consume messages from a queue.
//
// This will grant the following permissions:
//
//    - sqs:ChangeMessageVisibility
//    - sqs:DeleteMessage
//    - sqs:ReceiveMessage
//    - sqs:GetQueueAttributes
//    - sqs:GetQueueUrl
// Experimental.
func (q *jsiiProxy_Queue) GrantConsumeMessages(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grantConsumeMessages",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant an IAM principal permissions to purge all messages from the queue.
//
// This will grant the following permissions:
//
//   - sqs:PurgeQueue
//   - sqs:GetQueueAttributes
//   - sqs:GetQueueUrl
// Experimental.
func (q *jsiiProxy_Queue) GrantPurge(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grantPurge",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant access to send messages to a queue to the given identity.
//
// This will grant the following permissions:
//
//   - sqs:SendMessage
//   - sqs:GetQueueAttributes
//   - sqs:GetQueueUrl
// Experimental.
func (q *jsiiProxy_Queue) GrantSendMessages(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grantSendMessages",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Queue.
// Experimental.
func (q *jsiiProxy_Queue) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The approximate age of the oldest non-deleted message in the queue.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricApproximateAgeOfOldestMessage(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateAgeOfOldestMessage",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages in the queue that are delayed and not available for reading immediately.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricApproximateNumberOfMessagesDelayed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateNumberOfMessagesDelayed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that are in flight.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricApproximateNumberOfMessagesNotVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateNumberOfMessagesNotVisible",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages available for retrieval from the queue.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricApproximateNumberOfMessagesVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateNumberOfMessagesVisible",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of ReceiveMessage API calls that did not return a message.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricNumberOfEmptyReceives(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfEmptyReceives",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages deleted from the queue.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricNumberOfMessagesDeleted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfMessagesDeleted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages returned by calls to the ReceiveMessage action.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricNumberOfMessagesReceived(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfMessagesReceived",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages added to a queue.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricNumberOfMessagesSent(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfMessagesSent",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The size of messages added to a queue.
//
// Average over 5 minutes
// Experimental.
func (q *jsiiProxy_Queue) MetricSentMessageSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricSentMessageSize",
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
func (q *jsiiProxy_Queue) OnPrepare() {
	_jsii_.InvokeVoid(
		q,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_Queue) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
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
func (q *jsiiProxy_Queue) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_Queue) Prepare() {
	_jsii_.InvokeVoid(
		q,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_Queue) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (q *jsiiProxy_Queue) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
// Experimental.
func (q *jsiiProxy_Queue) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Reference to a queue.
// Experimental.
type QueueAttributes struct {
	// The ARN of the queue.
	// Experimental.
	QueueArn *string `json:"queueArn"`
	// KMS encryption key, if this queue is server-side encrypted by a KMS key.
	// Experimental.
	KeyArn *string `json:"keyArn"`
	// The name of the queue.
	// Experimental.
	QueueName *string `json:"queueName"`
	// The URL of the queue.
	// See: https://docs.aws.amazon.com/sdk-for-net/v2/developer-guide/QueueURL.html
	//
	// Experimental.
	QueueUrl *string `json:"queueUrl"`
}

// Reference to a new or existing Amazon SQS queue.
// Experimental.
type QueueBase interface {
	awscdk.Resource
	IQueue
	AutoCreatePolicy() *bool
	EncryptionMasterKey() awskms.IKey
	Env() *awscdk.ResourceEnvironment
	Fifo() *bool
	Node() awscdk.ConstructNode
	PhysicalName() *string
	QueueArn() *string
	QueueName() *string
	QueueUrl() *string
	Stack() awscdk.Stack
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantConsumeMessages(grantee awsiam.IGrantable) awsiam.Grant
	GrantPurge(grantee awsiam.IGrantable) awsiam.Grant
	GrantSendMessages(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateAgeOfOldestMessage(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateNumberOfMessagesDelayed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateNumberOfMessagesNotVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricApproximateNumberOfMessagesVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfEmptyReceives(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesDeleted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesReceived(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesSent(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSentMessageSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for QueueBase
type jsiiProxy_QueueBase struct {
	internal.Type__awscdkResource
	jsiiProxy_IQueue
}

func (j *jsiiProxy_QueueBase) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) EncryptionMasterKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionMasterKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) Fifo() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"fifo",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) QueueArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) QueueName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) QueueUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueueBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewQueueBase_Override(q QueueBase, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sqs.QueueBase",
		[]interface{}{scope, id, props},
		q,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func QueueBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.QueueBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func QueueBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.QueueBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a statement to the IAM resource policy associated with this queue.
//
// If this queue was created in this stack (`new Queue`), a queue policy
// will be automatically created upon the first call to `addToPolicy`. If
// the queue is imported (`Queue.import`), then this is a no-op.
// Experimental.
func (q *jsiiProxy_QueueBase) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		q,
		"addToResourcePolicy",
		[]interface{}{statement},
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
func (q *jsiiProxy_QueueBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		q,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (q *jsiiProxy_QueueBase) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueueBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueueBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		q,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the actions defined in queueActions to the identity Principal given on this SQS queue resource.
// Experimental.
func (q *jsiiProxy_QueueBase) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Grant permissions to consume messages from a queue.
//
// This will grant the following permissions:
//
//    - sqs:ChangeMessageVisibility
//    - sqs:DeleteMessage
//    - sqs:ReceiveMessage
//    - sqs:GetQueueAttributes
//    - sqs:GetQueueUrl
// Experimental.
func (q *jsiiProxy_QueueBase) GrantConsumeMessages(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grantConsumeMessages",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant an IAM principal permissions to purge all messages from the queue.
//
// This will grant the following permissions:
//
//   - sqs:PurgeQueue
//   - sqs:GetQueueAttributes
//   - sqs:GetQueueUrl
// Experimental.
func (q *jsiiProxy_QueueBase) GrantPurge(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grantPurge",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant access to send messages to a queue to the given identity.
//
// This will grant the following permissions:
//
//   - sqs:SendMessage
//   - sqs:GetQueueAttributes
//   - sqs:GetQueueUrl
// Experimental.
func (q *jsiiProxy_QueueBase) GrantSendMessages(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grantSendMessages",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Queue.
// Experimental.
func (q *jsiiProxy_QueueBase) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The approximate age of the oldest non-deleted message in the queue.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricApproximateAgeOfOldestMessage(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateAgeOfOldestMessage",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages in the queue that are delayed and not available for reading immediately.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricApproximateNumberOfMessagesDelayed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateNumberOfMessagesDelayed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that are in flight.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricApproximateNumberOfMessagesNotVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateNumberOfMessagesNotVisible",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages available for retrieval from the queue.
//
// Maximum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricApproximateNumberOfMessagesVisible(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricApproximateNumberOfMessagesVisible",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of ReceiveMessage API calls that did not return a message.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricNumberOfEmptyReceives(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfEmptyReceives",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages deleted from the queue.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricNumberOfMessagesDeleted(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfMessagesDeleted",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages returned by calls to the ReceiveMessage action.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricNumberOfMessagesReceived(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfMessagesReceived",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages added to a queue.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricNumberOfMessagesSent(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricNumberOfMessagesSent",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The size of messages added to a queue.
//
// Average over 5 minutes
// Experimental.
func (q *jsiiProxy_QueueBase) MetricSentMessageSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricSentMessageSize",
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
func (q *jsiiProxy_QueueBase) OnPrepare() {
	_jsii_.InvokeVoid(
		q,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_QueueBase) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
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
func (q *jsiiProxy_QueueBase) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueueBase) Prepare() {
	_jsii_.InvokeVoid(
		q,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_QueueBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (q *jsiiProxy_QueueBase) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
// Experimental.
func (q *jsiiProxy_QueueBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// What kind of encryption to apply to this queue.
// Experimental.
type QueueEncryption string

const (
	QueueEncryption_UNENCRYPTED QueueEncryption = "UNENCRYPTED"
	QueueEncryption_KMS_MANAGED QueueEncryption = "KMS_MANAGED"
	QueueEncryption_KMS QueueEncryption = "KMS"
)

// Applies a policy to SQS queues.
// Experimental.
type QueuePolicy interface {
	awscdk.Resource
	Document() awsiam.PolicyDocument
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
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

// The jsii proxy struct for QueuePolicy
type jsiiProxy_QueuePolicy struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_QueuePolicy) Document() awsiam.PolicyDocument {
	var returns awsiam.PolicyDocument
	_jsii_.Get(
		j,
		"document",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueuePolicy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueuePolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueuePolicy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QueuePolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewQueuePolicy(scope constructs.Construct, id *string, props *QueuePolicyProps) QueuePolicy {
	_init_.Initialize()

	j := jsiiProxy_QueuePolicy{}

	_jsii_.Create(
		"monocdk.aws_sqs.QueuePolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewQueuePolicy_Override(q QueuePolicy, scope constructs.Construct, id *string, props *QueuePolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sqs.QueuePolicy",
		[]interface{}{scope, id, props},
		q,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func QueuePolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.QueuePolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func QueuePolicy_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sqs.QueuePolicy",
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
func (q *jsiiProxy_QueuePolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		q,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (q *jsiiProxy_QueuePolicy) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueuePolicy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueuePolicy) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueuePolicy) OnPrepare() {
	_jsii_.InvokeVoid(
		q,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_QueuePolicy) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
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
func (q *jsiiProxy_QueuePolicy) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueuePolicy) Prepare() {
	_jsii_.InvokeVoid(
		q,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_QueuePolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (q *jsiiProxy_QueuePolicy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QueuePolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to associate SQS queues with a policy.
// Experimental.
type QueuePolicyProps struct {
	// The set of queues this policy applies to.
	// Experimental.
	Queues *[]IQueue `json:"queues"`
}

// Properties for creating a new Queue.
// Experimental.
type QueueProps struct {
	// Specifies whether to enable content-based deduplication.
	//
	// During the deduplication interval (5 minutes), Amazon SQS treats
	// messages that are sent with identical content (excluding attributes) as
	// duplicates and delivers only one copy of the message.
	//
	// If you don't enable content-based deduplication and you want to deduplicate
	// messages, provide an explicit deduplication ID in your SendMessage() call.
	//
	// (Only applies to FIFO queues.)
	// Experimental.
	ContentBasedDeduplication *bool `json:"contentBasedDeduplication"`
	// The length of time that Amazon SQS reuses a data key before calling KMS again.
	//
	// The value must be an integer between 60 (1 minute) and 86,400 (24
	// hours). The default is 300 (5 minutes).
	// Experimental.
	DataKeyReuse awscdk.Duration `json:"dataKeyReuse"`
	// Send messages to this queue if they were unsuccessfully dequeued a number of times.
	// Experimental.
	DeadLetterQueue *DeadLetterQueue `json:"deadLetterQueue"`
	// The time in seconds that the delivery of all messages in the queue is delayed.
	//
	// You can specify an integer value of 0 to 900 (15 minutes). The default
	// value is 0.
	// Experimental.
	DeliveryDelay awscdk.Duration `json:"deliveryDelay"`
	// Whether the contents of the queue are encrypted, and by what type of key.
	//
	// Be aware that encryption is not available in all regions, please see the docs
	// for current availability details.
	// Experimental.
	Encryption QueueEncryption `json:"encryption"`
	// External KMS master key to use for queue encryption.
	//
	// Individual messages will be encrypted using data keys. The data keys in
	// turn will be encrypted using this key, and reused for a maximum of
	// `dataKeyReuseSecs` seconds.
	//
	// If the 'encryptionMasterKey' property is set, 'encryption' type will be
	// implicitly set to "KMS".
	// Experimental.
	EncryptionMasterKey awskms.IKey `json:"encryptionMasterKey"`
	// Whether this a first-in-first-out (FIFO) queue.
	// Experimental.
	Fifo *bool `json:"fifo"`
	// The limit of how many bytes that a message can contain before Amazon SQS rejects it.
	//
	// You can specify an integer value from 1024 bytes (1 KiB) to 262144 bytes
	// (256 KiB). The default value is 262144 (256 KiB).
	// Experimental.
	MaxMessageSizeBytes *float64 `json:"maxMessageSizeBytes"`
	// A name for the queue.
	//
	// If specified and this is a FIFO queue, must end in the string '.fifo'.
	// Experimental.
	QueueName *string `json:"queueName"`
	// Default wait time for ReceiveMessage calls.
	//
	// Does not wait if set to 0, otherwise waits this amount of seconds
	// by default for messages to arrive.
	//
	// For more information, see Amazon SQS Long Poll.
	// Experimental.
	ReceiveMessageWaitTime awscdk.Duration `json:"receiveMessageWaitTime"`
	// Policy to apply when the user pool is removed from the stack.
	//
	// Even though queues are technically stateful, their contents are transient and it
	// is common to add and remove Queues while rearchitecting your application. The
	// default is therefore `DESTROY`. Change it to `RETAIN` if the messages are so
	// valuable that accidentally losing them would be unacceptable.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// The number of seconds that Amazon SQS retains a message.
	//
	// You can specify an integer value from 60 seconds (1 minute) to 1209600
	// seconds (14 days). The default value is 345600 seconds (4 days).
	// Experimental.
	RetentionPeriod awscdk.Duration `json:"retentionPeriod"`
	// Timeout of processing a single message.
	//
	// After dequeuing, the processor has this much time to handle the message
	// and delete it from the queue before it becomes visible again for dequeueing
	// by another processor.
	//
	// Values must be from 0 to 43200 seconds (12 hours). If you don't specify
	// a value, AWS CloudFormation uses the default value of 30 seconds.
	// Experimental.
	VisibilityTimeout awscdk.Duration `json:"visibilityTimeout"`
}


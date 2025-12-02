package awssns

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awscodestarnotifications"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awssns/internal"
	"github.com/aws/aws-cdk-go/awscdk/awssqs"
	"github.com/aws/constructs-go/constructs/v3"
)

// Between condition for a numeric attribute.
// Experimental.
type BetweenCondition struct {
	// The start value.
	// Experimental.
	Start *float64 `json:"start"`
	// The stop value.
	// Experimental.
	Stop *float64 `json:"stop"`
}

// A CloudFormation `AWS::SNS::Subscription`.
type CfnSubscription interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DeliveryPolicy() interface{}
	SetDeliveryPolicy(val interface{})
	Endpoint() *string
	SetEndpoint(val *string)
	FilterPolicy() interface{}
	SetFilterPolicy(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	Protocol() *string
	SetProtocol(val *string)
	RawMessageDelivery() interface{}
	SetRawMessageDelivery(val interface{})
	RedrivePolicy() interface{}
	SetRedrivePolicy(val interface{})
	Ref() *string
	Region() *string
	SetRegion(val *string)
	Stack() awscdk.Stack
	SubscriptionRoleArn() *string
	SetSubscriptionRoleArn(val *string)
	TopicArn() *string
	SetTopicArn(val *string)
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

// The jsii proxy struct for CfnSubscription
type jsiiProxy_CfnSubscription struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnSubscription) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) DeliveryPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deliveryPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) Endpoint() *string {
	var returns *string
	_jsii_.Get(
		j,
		"endpoint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) FilterPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"filterPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) Protocol() *string {
	var returns *string
	_jsii_.Get(
		j,
		"protocol",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) RawMessageDelivery() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawMessageDelivery",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) RedrivePolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"redrivePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) SubscriptionRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"subscriptionRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) TopicArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSubscription) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SNS::Subscription`.
func NewCfnSubscription(scope awscdk.Construct, id *string, props *CfnSubscriptionProps) CfnSubscription {
	_init_.Initialize()

	j := jsiiProxy_CfnSubscription{}

	_jsii_.Create(
		"monocdk.aws_sns.CfnSubscription",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SNS::Subscription`.
func NewCfnSubscription_Override(c CfnSubscription, scope awscdk.Construct, id *string, props *CfnSubscriptionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.CfnSubscription",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnSubscription) SetDeliveryPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"deliveryPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetEndpoint(val *string) {
	_jsii_.Set(
		j,
		"endpoint",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetFilterPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"filterPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetProtocol(val *string) {
	_jsii_.Set(
		j,
		"protocol",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetRawMessageDelivery(val interface{}) {
	_jsii_.Set(
		j,
		"rawMessageDelivery",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetRedrivePolicy(val interface{}) {
	_jsii_.Set(
		j,
		"redrivePolicy",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetRegion(val *string) {
	_jsii_.Set(
		j,
		"region",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetSubscriptionRoleArn(val *string) {
	_jsii_.Set(
		j,
		"subscriptionRoleArn",
		val,
	)
}

func (j *jsiiProxy_CfnSubscription) SetTopicArn(val *string) {
	_jsii_.Set(
		j,
		"topicArn",
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
func CfnSubscription_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnSubscription",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnSubscription_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnSubscription",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnSubscription_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnSubscription",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnSubscription_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_sns.CfnSubscription",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnSubscription) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnSubscription) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnSubscription) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnSubscription) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnSubscription) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnSubscription) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnSubscription) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnSubscription) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnSubscription) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnSubscription) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnSubscription) OnPrepare() {
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
func (c *jsiiProxy_CfnSubscription) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSubscription) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnSubscription) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnSubscription) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnSubscription) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnSubscription) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnSubscription) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSubscription) ToString() *string {
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
func (c *jsiiProxy_CfnSubscription) Validate() *[]*string {
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
func (c *jsiiProxy_CfnSubscription) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::SNS::Subscription`.
type CfnSubscriptionProps struct {
	// `AWS::SNS::Subscription.Protocol`.
	Protocol *string `json:"protocol"`
	// `AWS::SNS::Subscription.TopicArn`.
	TopicArn *string `json:"topicArn"`
	// `AWS::SNS::Subscription.DeliveryPolicy`.
	DeliveryPolicy interface{} `json:"deliveryPolicy"`
	// `AWS::SNS::Subscription.Endpoint`.
	Endpoint *string `json:"endpoint"`
	// `AWS::SNS::Subscription.FilterPolicy`.
	FilterPolicy interface{} `json:"filterPolicy"`
	// `AWS::SNS::Subscription.RawMessageDelivery`.
	RawMessageDelivery interface{} `json:"rawMessageDelivery"`
	// `AWS::SNS::Subscription.RedrivePolicy`.
	RedrivePolicy interface{} `json:"redrivePolicy"`
	// `AWS::SNS::Subscription.Region`.
	Region *string `json:"region"`
	// `AWS::SNS::Subscription.SubscriptionRoleArn`.
	SubscriptionRoleArn *string `json:"subscriptionRoleArn"`
}

// A CloudFormation `AWS::SNS::Topic`.
type CfnTopic interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrTopicName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ContentBasedDeduplication() interface{}
	SetContentBasedDeduplication(val interface{})
	CreationStack() *[]*string
	DisplayName() *string
	SetDisplayName(val *string)
	FifoTopic() interface{}
	SetFifoTopic(val interface{})
	KmsMasterKeyId() *string
	SetKmsMasterKeyId(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	Subscription() interface{}
	SetSubscription(val interface{})
	Tags() awscdk.TagManager
	TopicName() *string
	SetTopicName(val *string)
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

// The jsii proxy struct for CfnTopic
type jsiiProxy_CfnTopic struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnTopic) AttrTopicName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrTopicName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) ContentBasedDeduplication() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"contentBasedDeduplication",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) DisplayName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"displayName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) FifoTopic() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"fifoTopic",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) KmsMasterKeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kmsMasterKeyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) Subscription() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"subscription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) TopicName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopic) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SNS::Topic`.
func NewCfnTopic(scope awscdk.Construct, id *string, props *CfnTopicProps) CfnTopic {
	_init_.Initialize()

	j := jsiiProxy_CfnTopic{}

	_jsii_.Create(
		"monocdk.aws_sns.CfnTopic",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SNS::Topic`.
func NewCfnTopic_Override(c CfnTopic, scope awscdk.Construct, id *string, props *CfnTopicProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.CfnTopic",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnTopic) SetContentBasedDeduplication(val interface{}) {
	_jsii_.Set(
		j,
		"contentBasedDeduplication",
		val,
	)
}

func (j *jsiiProxy_CfnTopic) SetDisplayName(val *string) {
	_jsii_.Set(
		j,
		"displayName",
		val,
	)
}

func (j *jsiiProxy_CfnTopic) SetFifoTopic(val interface{}) {
	_jsii_.Set(
		j,
		"fifoTopic",
		val,
	)
}

func (j *jsiiProxy_CfnTopic) SetKmsMasterKeyId(val *string) {
	_jsii_.Set(
		j,
		"kmsMasterKeyId",
		val,
	)
}

func (j *jsiiProxy_CfnTopic) SetSubscription(val interface{}) {
	_jsii_.Set(
		j,
		"subscription",
		val,
	)
}

func (j *jsiiProxy_CfnTopic) SetTopicName(val *string) {
	_jsii_.Set(
		j,
		"topicName",
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
func CfnTopic_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnTopic",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnTopic_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnTopic",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnTopic_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnTopic",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnTopic_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_sns.CfnTopic",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnTopic) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnTopic) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnTopic) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnTopic) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnTopic) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnTopic) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnTopic) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnTopic) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnTopic) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnTopic) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnTopic) OnPrepare() {
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
func (c *jsiiProxy_CfnTopic) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTopic) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnTopic) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnTopic) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnTopic) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnTopic) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnTopic) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTopic) ToString() *string {
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
func (c *jsiiProxy_CfnTopic) Validate() *[]*string {
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
func (c *jsiiProxy_CfnTopic) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnTopic_SubscriptionProperty struct {
	// `CfnTopic.SubscriptionProperty.Endpoint`.
	Endpoint *string `json:"endpoint"`
	// `CfnTopic.SubscriptionProperty.Protocol`.
	Protocol *string `json:"protocol"`
}

// A CloudFormation `AWS::SNS::TopicPolicy`.
type CfnTopicPolicy interface {
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
	Ref() *string
	Stack() awscdk.Stack
	Topics() *[]*string
	SetTopics(val *[]*string)
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

// The jsii proxy struct for CfnTopicPolicy
type jsiiProxy_CfnTopicPolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnTopicPolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) PolicyDocument() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policyDocument",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) Topics() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"topics",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTopicPolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SNS::TopicPolicy`.
func NewCfnTopicPolicy(scope awscdk.Construct, id *string, props *CfnTopicPolicyProps) CfnTopicPolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnTopicPolicy{}

	_jsii_.Create(
		"monocdk.aws_sns.CfnTopicPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SNS::TopicPolicy`.
func NewCfnTopicPolicy_Override(c CfnTopicPolicy, scope awscdk.Construct, id *string, props *CfnTopicPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.CfnTopicPolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnTopicPolicy) SetPolicyDocument(val interface{}) {
	_jsii_.Set(
		j,
		"policyDocument",
		val,
	)
}

func (j *jsiiProxy_CfnTopicPolicy) SetTopics(val *[]*string) {
	_jsii_.Set(
		j,
		"topics",
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
func CfnTopicPolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnTopicPolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnTopicPolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnTopicPolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnTopicPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.CfnTopicPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnTopicPolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_sns.CfnTopicPolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnTopicPolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnTopicPolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnTopicPolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnTopicPolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnTopicPolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnTopicPolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnTopicPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnTopicPolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnTopicPolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnTopicPolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnTopicPolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnTopicPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTopicPolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnTopicPolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnTopicPolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnTopicPolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnTopicPolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnTopicPolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTopicPolicy) ToString() *string {
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
func (c *jsiiProxy_CfnTopicPolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnTopicPolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::SNS::TopicPolicy`.
type CfnTopicPolicyProps struct {
	// `AWS::SNS::TopicPolicy.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
	// `AWS::SNS::TopicPolicy.Topics`.
	Topics *[]*string `json:"topics"`
}

// Properties for defining a `AWS::SNS::Topic`.
type CfnTopicProps struct {
	// `AWS::SNS::Topic.ContentBasedDeduplication`.
	ContentBasedDeduplication interface{} `json:"contentBasedDeduplication"`
	// `AWS::SNS::Topic.DisplayName`.
	DisplayName *string `json:"displayName"`
	// `AWS::SNS::Topic.FifoTopic`.
	FifoTopic interface{} `json:"fifoTopic"`
	// `AWS::SNS::Topic.KmsMasterKeyId`.
	KmsMasterKeyId *string `json:"kmsMasterKeyId"`
	// `AWS::SNS::Topic.Subscription`.
	Subscription interface{} `json:"subscription"`
	// `AWS::SNS::Topic.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::SNS::Topic.TopicName`.
	TopicName *string `json:"topicName"`
}

// Represents an SNS topic.
// Experimental.
type ITopic interface {
	awscodestarnotifications.INotificationRuleTarget
	awscdk.IResource
	// Subscribe some endpoint to this topic.
	// Experimental.
	AddSubscription(subscription ITopicSubscription)
	// Adds a statement to the IAM resource policy associated with this topic.
	//
	// If this topic was created in this stack (`new Topic`), a topic policy
	// will be automatically created upon the first call to `addToPolicy`. If
	// the topic is imported (`Topic.import`), then this is a no-op.
	// Experimental.
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	// Grant topic publishing permissions to the given identity.
	// Experimental.
	GrantPublish(identity awsiam.IGrantable) awsiam.Grant
	// Return the given named metric for this Topic.
	// Experimental.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages published to your Amazon SNS topics.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfMessagesPublished(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages successfully delivered from your Amazon SNS topics to subscribing endpoints.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfNotificationsDelivered(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages that Amazon SNS failed to deliver.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfNotificationsFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages that were rejected by subscription filter policies.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfNotificationsFilteredOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages that were rejected by subscription filter policies because the messages' attributes are invalid.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfNotificationsFilteredOutInvalidAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of messages that were rejected by subscription filter policies because the messages have no attributes.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricNumberOfNotificationsFilteredOutNoMessageAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the size of messages published through this topic.
	//
	// Average over 5 minutes
	// Experimental.
	MetricPublishSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The charges you have accrued since the start of the current calendar month for sending SMS messages.
	//
	// Maximum over 5 minutes
	// Experimental.
	MetricSMSMonthToDateSpentUSD(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The rate of successful SMS message deliveries.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricSMSSuccessRate(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The ARN of the topic.
	// Experimental.
	TopicArn() *string
	// The name of the topic.
	// Experimental.
	TopicName() *string
}

// The jsii proxy for ITopic
type jsiiProxy_ITopic struct {
	internal.Type__awscodestarnotificationsINotificationRuleTarget
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_ITopic) AddSubscription(subscription ITopicSubscription) {
	_jsii_.InvokeVoid(
		i,
		"addSubscription",
		[]interface{}{subscription},
	)
}

func (i *jsiiProxy_ITopic) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		i,
		"addToResourcePolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) GrantPublish(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPublish",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricNumberOfMessagesPublished(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfMessagesPublished",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricNumberOfNotificationsDelivered(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfNotificationsDelivered",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricNumberOfNotificationsFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfNotificationsFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricNumberOfNotificationsFilteredOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfNotificationsFilteredOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricNumberOfNotificationsFilteredOutInvalidAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfNotificationsFilteredOutInvalidAttributes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricNumberOfNotificationsFilteredOutNoMessageAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricNumberOfNotificationsFilteredOutNoMessageAttributes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricPublishSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPublishSize",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricSMSMonthToDateSpentUSD(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricSMSMonthToDateSpentUSD",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) MetricSMSSuccessRate(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricSMSSuccessRate",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITopic) BindAsNotificationRuleTarget(scope constructs.Construct) *awscodestarnotifications.NotificationRuleTargetConfig {
	var returns *awscodestarnotifications.NotificationRuleTargetConfig

	_jsii_.Invoke(
		i,
		"bindAsNotificationRuleTarget",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_ITopic) TopicArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITopic) TopicName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITopic) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITopic) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITopic) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// Topic subscription.
// Experimental.
type ITopicSubscription interface {
	// Returns a configuration used to subscribe to an SNS topic.
	// Experimental.
	Bind(topic ITopic) *TopicSubscriptionConfig
}

// The jsii proxy for ITopicSubscription
type jsiiProxy_ITopicSubscription struct {
	_ byte // padding
}

func (i *jsiiProxy_ITopicSubscription) Bind(topic ITopic) *TopicSubscriptionConfig {
	var returns *TopicSubscriptionConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{topic},
		&returns,
	)

	return returns
}

// Conditions that can be applied to numeric attributes.
// Experimental.
type NumericConditions struct {
	// Match one or more values.
	// Experimental.
	Allowlist *[]*float64 `json:"allowlist"`
	// Match values that are between the specified values.
	// Experimental.
	Between *BetweenCondition `json:"between"`
	// Match values that are strictly between the specified values.
	// Experimental.
	BetweenStrict *BetweenCondition `json:"betweenStrict"`
	// Match values that are greater than the specified value.
	// Experimental.
	GreaterThan *float64 `json:"greaterThan"`
	// Match values that are greater than or equal to the specified value.
	// Experimental.
	GreaterThanOrEqualTo *float64 `json:"greaterThanOrEqualTo"`
	// Match values that are less than the specified value.
	// Experimental.
	LessThan *float64 `json:"lessThan"`
	// Match values that are less than or equal to the specified value.
	// Experimental.
	LessThanOrEqualTo *float64 `json:"lessThanOrEqualTo"`
	// Match one or more values.
	// Deprecated: use `allowlist`
	Whitelist *[]*float64 `json:"whitelist"`
}

// Conditions that can be applied to string attributes.
// Experimental.
type StringConditions struct {
	// Match one or more values.
	// Experimental.
	Allowlist *[]*string `json:"allowlist"`
	// Match any value that doesn't include any of the specified values.
	// Deprecated: use `denylist`
	Blacklist *[]*string `json:"blacklist"`
	// Match any value that doesn't include any of the specified values.
	// Experimental.
	Denylist *[]*string `json:"denylist"`
	// Matches values that begins with the specified prefixes.
	// Experimental.
	MatchPrefixes *[]*string `json:"matchPrefixes"`
	// Match one or more values.
	// Deprecated: use `allowlist`
	Whitelist *[]*string `json:"whitelist"`
}

// A new subscription.
//
// Prefer to use the `ITopic.addSubscription()` methods to create instances of
// this class.
// Experimental.
type Subscription interface {
	awscdk.Resource
	DeadLetterQueue() awssqs.IQueue
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

// The jsii proxy struct for Subscription
type jsiiProxy_Subscription struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_Subscription) DeadLetterQueue() awssqs.IQueue {
	var returns awssqs.IQueue
	_jsii_.Get(
		j,
		"deadLetterQueue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Subscription) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Subscription) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Subscription) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Subscription) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewSubscription(scope constructs.Construct, id *string, props *SubscriptionProps) Subscription {
	_init_.Initialize()

	j := jsiiProxy_Subscription{}

	_jsii_.Create(
		"monocdk.aws_sns.Subscription",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSubscription_Override(s Subscription, scope constructs.Construct, id *string, props *SubscriptionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.Subscription",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Subscription_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.Subscription",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Subscription_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.Subscription",
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
func (s *jsiiProxy_Subscription) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (s *jsiiProxy_Subscription) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_Subscription) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_Subscription) GetResourceNameAttribute(nameAttr *string) *string {
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
func (s *jsiiProxy_Subscription) OnPrepare() {
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
func (s *jsiiProxy_Subscription) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_Subscription) OnValidate() *[]*string {
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
func (s *jsiiProxy_Subscription) Prepare() {
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
func (s *jsiiProxy_Subscription) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_Subscription) ToString() *string {
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
func (s *jsiiProxy_Subscription) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A subscription filter for an attribute.
// Experimental.
type SubscriptionFilter interface {
	Conditions() *[]interface{}
}

// The jsii proxy struct for SubscriptionFilter
type jsiiProxy_SubscriptionFilter struct {
	_ byte // padding
}

func (j *jsiiProxy_SubscriptionFilter) Conditions() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}


// Experimental.
func NewSubscriptionFilter(conditions *[]interface{}) SubscriptionFilter {
	_init_.Initialize()

	j := jsiiProxy_SubscriptionFilter{}

	_jsii_.Create(
		"monocdk.aws_sns.SubscriptionFilter",
		[]interface{}{conditions},
		&j,
	)

	return &j
}

// Experimental.
func NewSubscriptionFilter_Override(s SubscriptionFilter, conditions *[]interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.SubscriptionFilter",
		[]interface{}{conditions},
		s,
	)
}

// Returns a subscription filter for attribute key matching.
// Experimental.
func SubscriptionFilter_ExistsFilter() SubscriptionFilter {
	_init_.Initialize()

	var returns SubscriptionFilter

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.SubscriptionFilter",
		"existsFilter",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a subscription filter for a numeric attribute.
// Experimental.
func SubscriptionFilter_NumericFilter(numericConditions *NumericConditions) SubscriptionFilter {
	_init_.Initialize()

	var returns SubscriptionFilter

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.SubscriptionFilter",
		"numericFilter",
		[]interface{}{numericConditions},
		&returns,
	)

	return returns
}

// Returns a subscription filter for a string attribute.
// Experimental.
func SubscriptionFilter_StringFilter(stringConditions *StringConditions) SubscriptionFilter {
	_init_.Initialize()

	var returns SubscriptionFilter

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.SubscriptionFilter",
		"stringFilter",
		[]interface{}{stringConditions},
		&returns,
	)

	return returns
}

// Options for creating a new subscription.
// Experimental.
type SubscriptionOptions struct {
	// The subscription endpoint.
	//
	// The meaning of this value depends on the value for 'protocol'.
	// Experimental.
	Endpoint *string `json:"endpoint"`
	// What type of subscription to add.
	// Experimental.
	Protocol SubscriptionProtocol `json:"protocol"`
	// Queue to be used as dead letter queue.
	//
	// If not passed no dead letter queue is enabled.
	// Experimental.
	DeadLetterQueue awssqs.IQueue `json:"deadLetterQueue"`
	// The filter policy.
	// Experimental.
	FilterPolicy *map[string]SubscriptionFilter `json:"filterPolicy"`
	// true if raw message delivery is enabled for the subscription.
	//
	// Raw messages are free of JSON formatting and can be
	// sent to HTTP/S and Amazon SQS endpoints. For more information, see GetSubscriptionAttributes in the Amazon Simple
	// Notification Service API Reference.
	// Experimental.
	RawMessageDelivery *bool `json:"rawMessageDelivery"`
	// The region where the topic resides, in the case of cross-region subscriptions.
	// Experimental.
	Region *string `json:"region"`
}

// Properties for creating a new subscription.
// Experimental.
type SubscriptionProps struct {
	// The subscription endpoint.
	//
	// The meaning of this value depends on the value for 'protocol'.
	// Experimental.
	Endpoint *string `json:"endpoint"`
	// What type of subscription to add.
	// Experimental.
	Protocol SubscriptionProtocol `json:"protocol"`
	// Queue to be used as dead letter queue.
	//
	// If not passed no dead letter queue is enabled.
	// Experimental.
	DeadLetterQueue awssqs.IQueue `json:"deadLetterQueue"`
	// The filter policy.
	// Experimental.
	FilterPolicy *map[string]SubscriptionFilter `json:"filterPolicy"`
	// true if raw message delivery is enabled for the subscription.
	//
	// Raw messages are free of JSON formatting and can be
	// sent to HTTP/S and Amazon SQS endpoints. For more information, see GetSubscriptionAttributes in the Amazon Simple
	// Notification Service API Reference.
	// Experimental.
	RawMessageDelivery *bool `json:"rawMessageDelivery"`
	// The region where the topic resides, in the case of cross-region subscriptions.
	// Experimental.
	Region *string `json:"region"`
	// The topic to subscribe to.
	// Experimental.
	Topic ITopic `json:"topic"`
}

// The type of subscription, controlling the type of the endpoint parameter.
// Experimental.
type SubscriptionProtocol string

const (
	SubscriptionProtocol_HTTP SubscriptionProtocol = "HTTP"
	SubscriptionProtocol_HTTPS SubscriptionProtocol = "HTTPS"
	SubscriptionProtocol_EMAIL SubscriptionProtocol = "EMAIL"
	SubscriptionProtocol_EMAIL_JSON SubscriptionProtocol = "EMAIL_JSON"
	SubscriptionProtocol_SMS SubscriptionProtocol = "SMS"
	SubscriptionProtocol_SQS SubscriptionProtocol = "SQS"
	SubscriptionProtocol_APPLICATION SubscriptionProtocol = "APPLICATION"
	SubscriptionProtocol_LAMBDA SubscriptionProtocol = "LAMBDA"
)

// A new SNS topic.
// Experimental.
type Topic interface {
	TopicBase
	AutoCreatePolicy() *bool
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	TopicArn() *string
	TopicName() *string
	AddSubscription(subscription ITopicSubscription)
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	BindAsNotificationRuleTarget(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleTargetConfig
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantPublish(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesPublished(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsDelivered(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFilteredOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFilteredOutInvalidAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFilteredOutNoMessageAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPublishSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSMSMonthToDateSpentUSD(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSMSSuccessRate(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Topic
type jsiiProxy_Topic struct {
	jsiiProxy_TopicBase
}

func (j *jsiiProxy_Topic) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Topic) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Topic) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Topic) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Topic) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Topic) TopicArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Topic) TopicName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicName",
		&returns,
	)
	return returns
}


// Experimental.
func NewTopic(scope constructs.Construct, id *string, props *TopicProps) Topic {
	_init_.Initialize()

	j := jsiiProxy_Topic{}

	_jsii_.Create(
		"monocdk.aws_sns.Topic",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewTopic_Override(t Topic, scope constructs.Construct, id *string, props *TopicProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.Topic",
		[]interface{}{scope, id, props},
		t,
	)
}

// Import an existing SNS topic provided an ARN.
// Experimental.
func Topic_FromTopicArn(scope constructs.Construct, id *string, topicArn *string) ITopic {
	_init_.Initialize()

	var returns ITopic

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.Topic",
		"fromTopicArn",
		[]interface{}{scope, id, topicArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Topic_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.Topic",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Topic_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.Topic",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Subscribe some endpoint to this topic.
// Experimental.
func (t *jsiiProxy_Topic) AddSubscription(subscription ITopicSubscription) {
	_jsii_.InvokeVoid(
		t,
		"addSubscription",
		[]interface{}{subscription},
	)
}

// Adds a statement to the IAM resource policy associated with this topic.
//
// If this topic was created in this stack (`new Topic`), a topic policy
// will be automatically created upon the first call to `addToPolicy`. If
// the topic is imported (`Topic.import`), then this is a no-op.
// Experimental.
func (t *jsiiProxy_Topic) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_Topic) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		t,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Represents a notification target That allows SNS topic to associate with this rule target.
// Experimental.
func (t *jsiiProxy_Topic) BindAsNotificationRuleTarget(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleTargetConfig {
	var returns *awscodestarnotifications.NotificationRuleTargetConfig

	_jsii_.Invoke(
		t,
		"bindAsNotificationRuleTarget",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// Experimental.
func (t *jsiiProxy_Topic) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_Topic) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_Topic) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant topic publishing permissions to the given identity.
// Experimental.
func (t *jsiiProxy_Topic) GrantPublish(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantPublish",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Topic.
// Experimental.
func (t *jsiiProxy_Topic) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The number of messages published to your Amazon SNS topics.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricNumberOfMessagesPublished(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfMessagesPublished",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages successfully delivered from your Amazon SNS topics to subscribing endpoints.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricNumberOfNotificationsDelivered(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsDelivered",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that Amazon SNS failed to deliver.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricNumberOfNotificationsFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that were rejected by subscription filter policies.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricNumberOfNotificationsFilteredOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFilteredOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that were rejected by subscription filter policies because the messages' attributes are invalid.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricNumberOfNotificationsFilteredOutInvalidAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFilteredOutInvalidAttributes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that were rejected by subscription filter policies because the messages have no attributes.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricNumberOfNotificationsFilteredOutNoMessageAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFilteredOutNoMessageAttributes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the size of messages published through this topic.
//
// Average over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricPublishSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricPublishSize",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The charges you have accrued since the start of the current calendar month for sending SMS messages.
//
// Maximum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricSMSMonthToDateSpentUSD(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSMSMonthToDateSpentUSD",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The rate of successful SMS message deliveries.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_Topic) MetricSMSSuccessRate(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSMSSuccessRate",
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
func (t *jsiiProxy_Topic) OnPrepare() {
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
func (t *jsiiProxy_Topic) OnSynthesize(session constructs.ISynthesisSession) {
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
func (t *jsiiProxy_Topic) OnValidate() *[]*string {
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
func (t *jsiiProxy_Topic) Prepare() {
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
func (t *jsiiProxy_Topic) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_Topic) ToString() *string {
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
// Experimental.
func (t *jsiiProxy_Topic) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Either a new or imported Topic.
// Experimental.
type TopicBase interface {
	awscdk.Resource
	ITopic
	AutoCreatePolicy() *bool
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	TopicArn() *string
	TopicName() *string
	AddSubscription(subscription ITopicSubscription)
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	BindAsNotificationRuleTarget(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleTargetConfig
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantPublish(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfMessagesPublished(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsDelivered(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFilteredOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFilteredOutInvalidAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNumberOfNotificationsFilteredOutNoMessageAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPublishSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSMSMonthToDateSpentUSD(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSMSSuccessRate(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for TopicBase
type jsiiProxy_TopicBase struct {
	internal.Type__awscdkResource
	jsiiProxy_ITopic
}

func (j *jsiiProxy_TopicBase) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicBase) TopicArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicBase) TopicName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicName",
		&returns,
	)
	return returns
}


// Experimental.
func NewTopicBase_Override(t TopicBase, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.TopicBase",
		[]interface{}{scope, id, props},
		t,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func TopicBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.TopicBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func TopicBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.TopicBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Subscribe some endpoint to this topic.
// Experimental.
func (t *jsiiProxy_TopicBase) AddSubscription(subscription ITopicSubscription) {
	_jsii_.InvokeVoid(
		t,
		"addSubscription",
		[]interface{}{subscription},
	)
}

// Adds a statement to the IAM resource policy associated with this topic.
//
// If this topic was created in this stack (`new Topic`), a topic policy
// will be automatically created upon the first call to `addToPolicy`. If
// the topic is imported (`Topic.import`), then this is a no-op.
// Experimental.
func (t *jsiiProxy_TopicBase) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TopicBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		t,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Represents a notification target That allows SNS topic to associate with this rule target.
// Experimental.
func (t *jsiiProxy_TopicBase) BindAsNotificationRuleTarget(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleTargetConfig {
	var returns *awscodestarnotifications.NotificationRuleTargetConfig

	_jsii_.Invoke(
		t,
		"bindAsNotificationRuleTarget",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// Experimental.
func (t *jsiiProxy_TopicBase) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TopicBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TopicBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant topic publishing permissions to the given identity.
// Experimental.
func (t *jsiiProxy_TopicBase) GrantPublish(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		t,
		"grantPublish",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Topic.
// Experimental.
func (t *jsiiProxy_TopicBase) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The number of messages published to your Amazon SNS topics.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricNumberOfMessagesPublished(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfMessagesPublished",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages successfully delivered from your Amazon SNS topics to subscribing endpoints.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricNumberOfNotificationsDelivered(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsDelivered",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that Amazon SNS failed to deliver.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricNumberOfNotificationsFailed(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFailed",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that were rejected by subscription filter policies.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricNumberOfNotificationsFilteredOut(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFilteredOut",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that were rejected by subscription filter policies because the messages' attributes are invalid.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricNumberOfNotificationsFilteredOutInvalidAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFilteredOutInvalidAttributes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of messages that were rejected by subscription filter policies because the messages have no attributes.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricNumberOfNotificationsFilteredOutNoMessageAttributes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricNumberOfNotificationsFilteredOutNoMessageAttributes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the size of messages published through this topic.
//
// Average over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricPublishSize(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricPublishSize",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The charges you have accrued since the start of the current calendar month for sending SMS messages.
//
// Maximum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricSMSMonthToDateSpentUSD(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSMSMonthToDateSpentUSD",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The rate of successful SMS message deliveries.
//
// Sum over 5 minutes
// Experimental.
func (t *jsiiProxy_TopicBase) MetricSMSSuccessRate(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		t,
		"metricSMSSuccessRate",
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
func (t *jsiiProxy_TopicBase) OnPrepare() {
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
func (t *jsiiProxy_TopicBase) OnSynthesize(session constructs.ISynthesisSession) {
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
func (t *jsiiProxy_TopicBase) OnValidate() *[]*string {
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
func (t *jsiiProxy_TopicBase) Prepare() {
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
func (t *jsiiProxy_TopicBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TopicBase) ToString() *string {
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
// Experimental.
func (t *jsiiProxy_TopicBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Applies a policy to SNS topics.
// Experimental.
type TopicPolicy interface {
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

// The jsii proxy struct for TopicPolicy
type jsiiProxy_TopicPolicy struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_TopicPolicy) Document() awsiam.PolicyDocument {
	var returns awsiam.PolicyDocument
	_jsii_.Get(
		j,
		"document",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicPolicy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicPolicy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TopicPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewTopicPolicy(scope constructs.Construct, id *string, props *TopicPolicyProps) TopicPolicy {
	_init_.Initialize()

	j := jsiiProxy_TopicPolicy{}

	_jsii_.Create(
		"monocdk.aws_sns.TopicPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewTopicPolicy_Override(t TopicPolicy, scope constructs.Construct, id *string, props *TopicPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_sns.TopicPolicy",
		[]interface{}{scope, id, props},
		t,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func TopicPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.TopicPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func TopicPolicy_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_sns.TopicPolicy",
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
func (t *jsiiProxy_TopicPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		t,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (t *jsiiProxy_TopicPolicy) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TopicPolicy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TopicPolicy) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TopicPolicy) OnPrepare() {
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
func (t *jsiiProxy_TopicPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (t *jsiiProxy_TopicPolicy) OnValidate() *[]*string {
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
func (t *jsiiProxy_TopicPolicy) Prepare() {
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
func (t *jsiiProxy_TopicPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TopicPolicy) ToString() *string {
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
func (t *jsiiProxy_TopicPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to associate SNS topics with a policy.
// Experimental.
type TopicPolicyProps struct {
	// The set of topics this policy applies to.
	// Experimental.
	Topics *[]ITopic `json:"topics"`
	// IAM policy document to apply to topic(s).
	// Experimental.
	PolicyDocument awsiam.PolicyDocument `json:"policyDocument"`
}

// Properties for a new SNS topic.
// Experimental.
type TopicProps struct {
	// Enables content-based deduplication for FIFO topics.
	// Experimental.
	ContentBasedDeduplication *bool `json:"contentBasedDeduplication"`
	// A developer-defined string that can be used to identify this SNS topic.
	// Experimental.
	DisplayName *string `json:"displayName"`
	// Set to true to create a FIFO topic.
	// Experimental.
	Fifo *bool `json:"fifo"`
	// A KMS Key, either managed by this CDK app, or imported.
	// Experimental.
	MasterKey awskms.IKey `json:"masterKey"`
	// A name for the topic.
	//
	// If you don't specify a name, AWS CloudFormation generates a unique
	// physical ID and uses that ID for the topic name. For more information,
	// see Name Type.
	// Experimental.
	TopicName *string `json:"topicName"`
}

// Subscription configuration.
// Experimental.
type TopicSubscriptionConfig struct {
	// The subscription endpoint.
	//
	// The meaning of this value depends on the value for 'protocol'.
	// Experimental.
	Endpoint *string `json:"endpoint"`
	// What type of subscription to add.
	// Experimental.
	Protocol SubscriptionProtocol `json:"protocol"`
	// Queue to be used as dead letter queue.
	//
	// If not passed no dead letter queue is enabled.
	// Experimental.
	DeadLetterQueue awssqs.IQueue `json:"deadLetterQueue"`
	// The filter policy.
	// Experimental.
	FilterPolicy *map[string]SubscriptionFilter `json:"filterPolicy"`
	// true if raw message delivery is enabled for the subscription.
	//
	// Raw messages are free of JSON formatting and can be
	// sent to HTTP/S and Amazon SQS endpoints. For more information, see GetSubscriptionAttributes in the Amazon Simple
	// Notification Service API Reference.
	// Experimental.
	RawMessageDelivery *bool `json:"rawMessageDelivery"`
	// The region where the topic resides, in the case of cross-region subscriptions.
	// Experimental.
	Region *string `json:"region"`
	// The id of the SNS subscription resource created under `scope`.
	//
	// In most
	// cases, it is recommended to use the `uniqueId` of the topic you are
	// subscribing to.
	// Experimental.
	SubscriberId *string `json:"subscriberId"`
	// The scope in which to create the SNS subscription resource.
	//
	// Normally you'd
	// want the subscription to be created on the consuming stack because the
	// topic is usually referenced by the consumer's resource policy (e.g. SQS
	// queue policy). Otherwise, it will cause a cyclic reference.
	//
	// If this is undefined, the subscription will be created on the topic's stack.
	// Experimental.
	SubscriberScope awscdk.Construct `json:"subscriberScope"`
}


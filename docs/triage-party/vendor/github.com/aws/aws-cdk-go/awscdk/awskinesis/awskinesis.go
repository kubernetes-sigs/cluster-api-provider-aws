package awskinesis

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskinesis/internal"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/constructs-go/constructs/v3"
)

// A CloudFormation `AWS::Kinesis::Stream`.
type CfnStream interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	RetentionPeriodHours() *float64
	SetRetentionPeriodHours(val *float64)
	ShardCount() *float64
	SetShardCount(val *float64)
	Stack() awscdk.Stack
	StreamEncryption() interface{}
	SetStreamEncryption(val interface{})
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

// The jsii proxy struct for CfnStream
type jsiiProxy_CfnStream struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnStream) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) RetentionPeriodHours() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"retentionPeriodHours",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) ShardCount() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"shardCount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) StreamEncryption() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"streamEncryption",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStream) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Kinesis::Stream`.
func NewCfnStream(scope awscdk.Construct, id *string, props *CfnStreamProps) CfnStream {
	_init_.Initialize()

	j := jsiiProxy_CfnStream{}

	_jsii_.Create(
		"monocdk.aws_kinesis.CfnStream",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Kinesis::Stream`.
func NewCfnStream_Override(c CfnStream, scope awscdk.Construct, id *string, props *CfnStreamProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kinesis.CfnStream",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStream) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnStream) SetRetentionPeriodHours(val *float64) {
	_jsii_.Set(
		j,
		"retentionPeriodHours",
		val,
	)
}

func (j *jsiiProxy_CfnStream) SetShardCount(val *float64) {
	_jsii_.Set(
		j,
		"shardCount",
		val,
	)
}

func (j *jsiiProxy_CfnStream) SetStreamEncryption(val interface{}) {
	_jsii_.Set(
		j,
		"streamEncryption",
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
func CfnStream_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.CfnStream",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnStream_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.CfnStream",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnStream_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.CfnStream",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStream_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_kinesis.CfnStream",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnStream) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnStream) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnStream) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnStream) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnStream) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnStream) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnStream) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnStream) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnStream) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnStream) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnStream) OnPrepare() {
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
func (c *jsiiProxy_CfnStream) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStream) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnStream) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnStream) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnStream) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnStream) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnStream) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStream) ToString() *string {
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
func (c *jsiiProxy_CfnStream) Validate() *[]*string {
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
func (c *jsiiProxy_CfnStream) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnStream_StreamEncryptionProperty struct {
	// `CfnStream.StreamEncryptionProperty.EncryptionType`.
	EncryptionType *string `json:"encryptionType"`
	// `CfnStream.StreamEncryptionProperty.KeyId`.
	KeyId *string `json:"keyId"`
}

// A CloudFormation `AWS::Kinesis::StreamConsumer`.
type CfnStreamConsumer interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrConsumerArn() *string
	AttrConsumerCreationTimestamp() *string
	AttrConsumerName() *string
	AttrConsumerStatus() *string
	AttrStreamArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ConsumerName() *string
	SetConsumerName(val *string)
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	StreamArn() *string
	SetStreamArn(val *string)
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

// The jsii proxy struct for CfnStreamConsumer
type jsiiProxy_CfnStreamConsumer struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnStreamConsumer) AttrConsumerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrConsumerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) AttrConsumerCreationTimestamp() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrConsumerCreationTimestamp",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) AttrConsumerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrConsumerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) AttrConsumerStatus() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrConsumerStatus",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) AttrStreamArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrStreamArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) ConsumerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"consumerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) StreamArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"streamArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStreamConsumer) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Kinesis::StreamConsumer`.
func NewCfnStreamConsumer(scope awscdk.Construct, id *string, props *CfnStreamConsumerProps) CfnStreamConsumer {
	_init_.Initialize()

	j := jsiiProxy_CfnStreamConsumer{}

	_jsii_.Create(
		"monocdk.aws_kinesis.CfnStreamConsumer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Kinesis::StreamConsumer`.
func NewCfnStreamConsumer_Override(c CfnStreamConsumer, scope awscdk.Construct, id *string, props *CfnStreamConsumerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kinesis.CfnStreamConsumer",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStreamConsumer) SetConsumerName(val *string) {
	_jsii_.Set(
		j,
		"consumerName",
		val,
	)
}

func (j *jsiiProxy_CfnStreamConsumer) SetStreamArn(val *string) {
	_jsii_.Set(
		j,
		"streamArn",
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
func CfnStreamConsumer_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.CfnStreamConsumer",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnStreamConsumer_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.CfnStreamConsumer",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnStreamConsumer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.CfnStreamConsumer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStreamConsumer_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_kinesis.CfnStreamConsumer",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnStreamConsumer) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnStreamConsumer) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnStreamConsumer) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnStreamConsumer) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnStreamConsumer) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnStreamConsumer) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnStreamConsumer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnStreamConsumer) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnStreamConsumer) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnStreamConsumer) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnStreamConsumer) OnPrepare() {
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
func (c *jsiiProxy_CfnStreamConsumer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStreamConsumer) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnStreamConsumer) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnStreamConsumer) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnStreamConsumer) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnStreamConsumer) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnStreamConsumer) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStreamConsumer) ToString() *string {
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
func (c *jsiiProxy_CfnStreamConsumer) Validate() *[]*string {
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
func (c *jsiiProxy_CfnStreamConsumer) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Kinesis::StreamConsumer`.
type CfnStreamConsumerProps struct {
	// `AWS::Kinesis::StreamConsumer.ConsumerName`.
	ConsumerName *string `json:"consumerName"`
	// `AWS::Kinesis::StreamConsumer.StreamARN`.
	StreamArn *string `json:"streamArn"`
}

// Properties for defining a `AWS::Kinesis::Stream`.
type CfnStreamProps struct {
	// `AWS::Kinesis::Stream.ShardCount`.
	ShardCount *float64 `json:"shardCount"`
	// `AWS::Kinesis::Stream.Name`.
	Name *string `json:"name"`
	// `AWS::Kinesis::Stream.RetentionPeriodHours`.
	RetentionPeriodHours *float64 `json:"retentionPeriodHours"`
	// `AWS::Kinesis::Stream.StreamEncryption`.
	StreamEncryption interface{} `json:"streamEncryption"`
	// `AWS::Kinesis::Stream.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A Kinesis Stream.
// Experimental.
type IStream interface {
	awscdk.IResource
	// Grant the indicated permissions on this stream to the provided IAM principal.
	// Experimental.
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	// Grant read permissions for this stream and its contents to an IAM principal (Role/Group/User).
	//
	// If an encryption key is used, permission to ues the key to decrypt the
	// contents of the stream will also be granted.
	// Experimental.
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	// Grants read/write permissions for this stream and its contents to an IAM principal (Role/Group/User).
	//
	// If an encryption key is used, permission to use the key for
	// encrypt/decrypt will also be granted.
	// Experimental.
	GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant
	// Grant write permissions for this stream and its contents to an IAM principal (Role/Group/User).
	//
	// If an encryption key is used, permission to ues the key to encrypt the
	// contents of the stream will also be granted.
	// Experimental.
	GrantWrite(grantee awsiam.IGrantable) awsiam.Grant
	// Return stream metric based from its metric name.
	// Experimental.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of records retrieved from the shard, measured over the specified time period.
	//
	// Minimum, Maximum, and
	// Average statistics represent the records in a single GetRecords operation for the stream in the specified time
	// period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricGetRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of bytes retrieved from the Kinesis stream, measured over the specified time period.
	//
	// Minimum, Maximum,
	// and Average statistics represent the bytes in a single GetRecords operation for the stream in the specified time
	// period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricGetRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The age of the last record in all GetRecords calls made against a Kinesis stream, measured over the specified time period.
	//
	// Age is the difference between the current time and when the last record of the GetRecords call was written
	// to the stream. The Minimum and Maximum statistics can be used to track the progress of Kinesis consumer
	// applications. A value of zero indicates that the records being read are completely caught up with the stream.
	//
	// The metric defaults to maximum over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricGetRecordsIteratorAgeMilliseconds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The time taken per GetRecords operation, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricGetRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of successful GetRecords operations per stream, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricGetRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of bytes successfully put to the Kinesis stream over the specified time period.
	//
	// This metric includes
	// bytes from PutRecord and PutRecords operations. Minimum, Maximum, and Average statistics represent the bytes in a
	// single put operation for the stream in the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricIncomingBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of records successfully put to the Kinesis stream over the specified time period.
	//
	// This metric includes
	// record counts from PutRecord and PutRecords operations. Minimum, Maximum, and Average statistics represent the
	// records in a single put operation for the stream in the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricIncomingRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of bytes put to the Kinesis stream using the PutRecord operation over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The time taken per PutRecord operation, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of bytes put to the Kinesis stream using the PutRecords operation over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of records rejected due to internal failures in a PutRecords operation per Kinesis data stream, measured over the specified time period.
	//
	// Occasional internal failures are to be expected and should be retried.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordsFailedRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The time taken per PutRecords operation, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of PutRecords operations where at least one record succeeded, per Kinesis stream, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of successful records in a PutRecords operation per Kinesis data stream, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordsSuccessfulRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of records rejected due to throttling in a PutRecords operation per Kinesis data stream, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordsThrottledRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The total number of records sent in a PutRecords operation per Kinesis data stream, measured over the specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordsTotalRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of successful PutRecord operations per Kinesis stream, measured over the specified time period.
	//
	// Average
	// reflects the percentage of successful writes to a stream.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricPutRecordSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of GetRecords calls throttled for the stream over the specified time period.
	//
	// The most commonly used
	// statistic for this metric is Average.
	//
	// When the Minimum statistic has a value of 1, all records were throttled for the stream during the specified time
	// period.
	//
	// When the Maximum statistic has a value of 0 (zero), no records were throttled for the stream during the specified
	// time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties
	// Experimental.
	MetricReadProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The number of records rejected due to throttling for the stream over the specified time period.
	//
	// This metric
	// includes throttling from PutRecord and PutRecords operations.
	//
	// When the Minimum statistic has a non-zero value, records were being throttled for the stream during the specified
	// time period.
	//
	// When the Maximum statistic has a value of 0 (zero), no records were being throttled for the stream during the
	// specified time period.
	//
	// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
	// Experimental.
	MetricWriteProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Optional KMS encryption key associated with this stream.
	// Experimental.
	EncryptionKey() awskms.IKey
	// The ARN of the stream.
	// Experimental.
	StreamArn() *string
	// The name of the stream.
	// Experimental.
	StreamName() *string
}

// The jsii proxy for IStream
type jsiiProxy_IStream struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IStream) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
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

func (i *jsiiProxy_IStream) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantReadWrite",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) GrantWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantWrite",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricGetRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricGetRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricGetRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricGetRecordsBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricGetRecordsIteratorAgeMilliseconds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricGetRecordsIteratorAgeMilliseconds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricGetRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricGetRecordsLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricGetRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricGetRecordsSuccess",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricIncomingBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricIncomingBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricIncomingRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricIncomingRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordsBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordsFailedRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordsFailedRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordsLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordsSuccess",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordsSuccessfulRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordsSuccessfulRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordsThrottledRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordsThrottledRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordsTotalRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordsTotalRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricPutRecordSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricPutRecordSuccess",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricReadProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricReadProvisionedThroughputExceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStream) MetricWriteProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricWriteProvisionedThroughputExceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IStream) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStream) StreamArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"streamArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStream) StreamName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"streamName",
		&returns,
	)
	return returns
}

// A Kinesis stream.
//
// Can be encrypted with a KMS key.
// Experimental.
type Stream interface {
	awscdk.Resource
	IStream
	EncryptionKey() awskms.IKey
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	StreamArn() *string
	StreamName() *string
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant
	GrantWrite(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricGetRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricGetRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricGetRecordsIteratorAgeMilliseconds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricGetRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricGetRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIncomingBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIncomingRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordsFailedRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordsSuccessfulRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordsThrottledRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordsTotalRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricPutRecordSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricReadProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricWriteProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Stream
type jsiiProxy_Stream struct {
	internal.Type__awscdkResource
	jsiiProxy_IStream
}

func (j *jsiiProxy_Stream) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stream) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stream) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stream) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stream) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stream) StreamArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"streamArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stream) StreamName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"streamName",
		&returns,
	)
	return returns
}


// Experimental.
func NewStream(scope constructs.Construct, id *string, props *StreamProps) Stream {
	_init_.Initialize()

	j := jsiiProxy_Stream{}

	_jsii_.Create(
		"monocdk.aws_kinesis.Stream",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStream_Override(s Stream, scope constructs.Construct, id *string, props *StreamProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kinesis.Stream",
		[]interface{}{scope, id, props},
		s,
	)
}

// Import an existing Kinesis Stream provided an ARN.
// Experimental.
func Stream_FromStreamArn(scope constructs.Construct, id *string, streamArn *string) IStream {
	_init_.Initialize()

	var returns IStream

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.Stream",
		"fromStreamArn",
		[]interface{}{scope, id, streamArn},
		&returns,
	)

	return returns
}

// Creates a Stream construct that represents an external stream.
// Experimental.
func Stream_FromStreamAttributes(scope constructs.Construct, id *string, attrs *StreamAttributes) IStream {
	_init_.Initialize()

	var returns IStream

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.Stream",
		"fromStreamAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Stream_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.Stream",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Stream_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesis.Stream",
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
func (s *jsiiProxy_Stream) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (s *jsiiProxy_Stream) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_Stream) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_Stream) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the indicated permissions on this stream to the given IAM principal (Role/Group/User).
// Experimental.
func (s *jsiiProxy_Stream) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
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

// Grant read permissions for this stream and its contents to an IAM principal (Role/Group/User).
//
// If an encryption key is used, permission to ues the key to decrypt the
// contents of the stream will also be granted.
// Experimental.
func (s *jsiiProxy_Stream) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grants read/write permissions for this stream and its contents to an IAM principal (Role/Group/User).
//
// If an encryption key is used, permission to use the key for
// encrypt/decrypt will also be granted.
// Experimental.
func (s *jsiiProxy_Stream) GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantReadWrite",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant write permissions for this stream and its contents to an IAM principal (Role/Group/User).
//
// If an encryption key is used, permission to ues the key to encrypt the
// contents of the stream will also be granted.
// Experimental.
func (s *jsiiProxy_Stream) GrantWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantWrite",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return stream metric based from its metric name.
// Experimental.
func (s *jsiiProxy_Stream) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The number of records retrieved from the shard, measured over the specified time period.
//
// Minimum, Maximum, and
// Average statistics represent the records in a single GetRecords operation for the stream in the specified time
// period.
//
// average
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricGetRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricGetRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of bytes retrieved from the Kinesis stream, measured over the specified time period.
//
// Minimum, Maximum,
// and Average statistics represent the bytes in a single GetRecords operation for the stream in the specified time
// period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricGetRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricGetRecordsBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The age of the last record in all GetRecords calls made against a Kinesis stream, measured over the specified time period.
//
// Age is the difference between the current time and when the last record of the GetRecords call was written
// to the stream. The Minimum and Maximum statistics can be used to track the progress of Kinesis consumer
// applications. A value of zero indicates that the records being read are completely caught up with the stream.
//
// The metric defaults to maximum over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricGetRecordsIteratorAgeMilliseconds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricGetRecordsIteratorAgeMilliseconds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of successful GetRecords operations per stream, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricGetRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricGetRecordsLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of successful GetRecords operations per stream, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricGetRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricGetRecordsSuccess",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of bytes successfully put to the Kinesis stream over the specified time period.
//
// This metric includes
// bytes from PutRecord and PutRecords operations. Minimum, Maximum, and Average statistics represent the bytes in a
// single put operation for the stream in the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricIncomingBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricIncomingBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of records successfully put to the Kinesis stream over the specified time period.
//
// This metric includes
// record counts from PutRecord and PutRecords operations. Minimum, Maximum, and Average statistics represent the
// records in a single put operation for the stream in the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricIncomingRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricIncomingRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of bytes put to the Kinesis stream using the PutRecord operation over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time taken per PutRecord operation, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of bytes put to the Kinesis stream using the PutRecords operation over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordsBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordsBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of records rejected due to internal failures in a PutRecords operation per Kinesis data stream, measured over the specified time period.
//
// Occasional internal failures are to be expected and should be retried.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordsFailedRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordsFailedRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time taken per PutRecords operation, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordsLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordsLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of PutRecords operations where at least one record succeeded, per Kinesis stream, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordsSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordsSuccess",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of successful records in a PutRecords operation per Kinesis data stream, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordsSuccessfulRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordsSuccessfulRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of records rejected due to throttling in a PutRecords operation per Kinesis data stream, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordsThrottledRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordsThrottledRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of records sent in a PutRecords operation per Kinesis data stream, measured over the specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordsTotalRecords(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordsTotalRecords",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of successful PutRecord operations per Kinesis stream, measured over the specified time period.
//
// Average
// reflects the percentage of successful writes to a stream.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricPutRecordSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricPutRecordSuccess",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of GetRecords calls throttled for the stream over the specified time period.
//
// The most commonly used
// statistic for this metric is Average.
//
// When the Minimum statistic has a value of 1, all records were throttled for the stream during the specified time
// period.
//
// When the Maximum statistic has a value of 0 (zero), no records were throttled for the stream during the specified
// time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties
// Experimental.
func (s *jsiiProxy_Stream) MetricReadProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricReadProvisionedThroughputExceeded",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of records rejected due to throttling for the stream over the specified time period.
//
// This metric
// includes throttling from PutRecord and PutRecords operations.
//
// When the Minimum statistic has a non-zero value, records were being throttled for the stream during the specified
// time period.
//
// When the Maximum statistic has a value of 0 (zero), no records were being throttled for the stream during the
// specified time period.
//
// The metric defaults to average over 5 minutes, it can be changed by passing `statistic` and `period` properties.
// Experimental.
func (s *jsiiProxy_Stream) MetricWriteProvisionedThroughputExceeded(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricWriteProvisionedThroughputExceeded",
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
func (s *jsiiProxy_Stream) OnPrepare() {
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
func (s *jsiiProxy_Stream) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_Stream) OnValidate() *[]*string {
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
func (s *jsiiProxy_Stream) Prepare() {
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
func (s *jsiiProxy_Stream) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_Stream) ToString() *string {
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
func (s *jsiiProxy_Stream) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A reference to a stream.
//
// The easiest way to instantiate is to call
// `stream.export()`. Then, the consumer can use `Stream.import(this, ref)` and
// get a `Stream`.
// Experimental.
type StreamAttributes struct {
	// The ARN of the stream.
	// Experimental.
	StreamArn *string `json:"streamArn"`
	// The KMS key securing the contents of the stream if encryption is enabled.
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
}

// What kind of server-side encryption to apply to this stream.
// Experimental.
type StreamEncryption string

const (
	StreamEncryption_UNENCRYPTED StreamEncryption = "UNENCRYPTED"
	StreamEncryption_KMS StreamEncryption = "KMS"
	StreamEncryption_MANAGED StreamEncryption = "MANAGED"
)

// Properties for a Kinesis Stream.
// Experimental.
type StreamProps struct {
	// The kind of server-side encryption to apply to this stream.
	//
	// If you choose KMS, you can specify a KMS key via `encryptionKey`. If
	// encryption key is not specified, a key will automatically be created.
	// Experimental.
	Encryption StreamEncryption `json:"encryption"`
	// External KMS key to use for stream encryption.
	//
	// The 'encryption' property must be set to "Kms".
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// The number of hours for the data records that are stored in shards to remain accessible.
	// Experimental.
	RetentionPeriod awscdk.Duration `json:"retentionPeriod"`
	// The number of shards for the stream.
	// Experimental.
	ShardCount *float64 `json:"shardCount"`
	// Enforces a particular physical stream name.
	// Experimental.
	StreamName *string `json:"streamName"`
}


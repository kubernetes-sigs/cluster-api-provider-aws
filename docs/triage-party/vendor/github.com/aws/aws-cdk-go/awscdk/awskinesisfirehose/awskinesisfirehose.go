package awskinesisfirehose

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awskinesisfirehose/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// A CloudFormation `AWS::KinesisFirehose::DeliveryStream`.
type CfnDeliveryStream interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DeliveryStreamEncryptionConfigurationInput() interface{}
	SetDeliveryStreamEncryptionConfigurationInput(val interface{})
	DeliveryStreamName() *string
	SetDeliveryStreamName(val *string)
	DeliveryStreamType() *string
	SetDeliveryStreamType(val *string)
	ElasticsearchDestinationConfiguration() interface{}
	SetElasticsearchDestinationConfiguration(val interface{})
	ExtendedS3DestinationConfiguration() interface{}
	SetExtendedS3DestinationConfiguration(val interface{})
	HttpEndpointDestinationConfiguration() interface{}
	SetHttpEndpointDestinationConfiguration(val interface{})
	KinesisStreamSourceConfiguration() interface{}
	SetKinesisStreamSourceConfiguration(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	RedshiftDestinationConfiguration() interface{}
	SetRedshiftDestinationConfiguration(val interface{})
	Ref() *string
	S3DestinationConfiguration() interface{}
	SetS3DestinationConfiguration(val interface{})
	SplunkDestinationConfiguration() interface{}
	SetSplunkDestinationConfiguration(val interface{})
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

// The jsii proxy struct for CfnDeliveryStream
type jsiiProxy_CfnDeliveryStream struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDeliveryStream) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) DeliveryStreamEncryptionConfigurationInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deliveryStreamEncryptionConfigurationInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) DeliveryStreamName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deliveryStreamName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) DeliveryStreamType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deliveryStreamType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) ElasticsearchDestinationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"elasticsearchDestinationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) ExtendedS3DestinationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"extendedS3DestinationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) HttpEndpointDestinationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"httpEndpointDestinationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) KinesisStreamSourceConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"kinesisStreamSourceConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) RedshiftDestinationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"redshiftDestinationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) S3DestinationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"s3DestinationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) SplunkDestinationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"splunkDestinationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeliveryStream) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::KinesisFirehose::DeliveryStream`.
func NewCfnDeliveryStream(scope awscdk.Construct, id *string, props *CfnDeliveryStreamProps) CfnDeliveryStream {
	_init_.Initialize()

	j := jsiiProxy_CfnDeliveryStream{}

	_jsii_.Create(
		"monocdk.aws_kinesisfirehose.CfnDeliveryStream",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::KinesisFirehose::DeliveryStream`.
func NewCfnDeliveryStream_Override(c CfnDeliveryStream, scope awscdk.Construct, id *string, props *CfnDeliveryStreamProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kinesisfirehose.CfnDeliveryStream",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetDeliveryStreamEncryptionConfigurationInput(val interface{}) {
	_jsii_.Set(
		j,
		"deliveryStreamEncryptionConfigurationInput",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetDeliveryStreamName(val *string) {
	_jsii_.Set(
		j,
		"deliveryStreamName",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetDeliveryStreamType(val *string) {
	_jsii_.Set(
		j,
		"deliveryStreamType",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetElasticsearchDestinationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"elasticsearchDestinationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetExtendedS3DestinationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"extendedS3DestinationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetHttpEndpointDestinationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"httpEndpointDestinationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetKinesisStreamSourceConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"kinesisStreamSourceConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetRedshiftDestinationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"redshiftDestinationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetS3DestinationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"s3DestinationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnDeliveryStream) SetSplunkDestinationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"splunkDestinationConfiguration",
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
func CfnDeliveryStream_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesisfirehose.CfnDeliveryStream",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnDeliveryStream_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesisfirehose.CfnDeliveryStream",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnDeliveryStream_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kinesisfirehose.CfnDeliveryStream",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDeliveryStream_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_kinesisfirehose.CfnDeliveryStream",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnDeliveryStream) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnDeliveryStream) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnDeliveryStream) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnDeliveryStream) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnDeliveryStream) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnDeliveryStream) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnDeliveryStream) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnDeliveryStream) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnDeliveryStream) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnDeliveryStream) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnDeliveryStream) OnPrepare() {
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
func (c *jsiiProxy_CfnDeliveryStream) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDeliveryStream) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnDeliveryStream) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnDeliveryStream) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnDeliveryStream) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnDeliveryStream) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnDeliveryStream) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDeliveryStream) ToString() *string {
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
func (c *jsiiProxy_CfnDeliveryStream) Validate() *[]*string {
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
func (c *jsiiProxy_CfnDeliveryStream) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnDeliveryStream_BufferingHintsProperty struct {
	// `CfnDeliveryStream.BufferingHintsProperty.IntervalInSeconds`.
	IntervalInSeconds *float64 `json:"intervalInSeconds"`
	// `CfnDeliveryStream.BufferingHintsProperty.SizeInMBs`.
	SizeInMBs *float64 `json:"sizeInMBs"`
}

type CfnDeliveryStream_CloudWatchLoggingOptionsProperty struct {
	// `CfnDeliveryStream.CloudWatchLoggingOptionsProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `CfnDeliveryStream.CloudWatchLoggingOptionsProperty.LogGroupName`.
	LogGroupName *string `json:"logGroupName"`
	// `CfnDeliveryStream.CloudWatchLoggingOptionsProperty.LogStreamName`.
	LogStreamName *string `json:"logStreamName"`
}

type CfnDeliveryStream_CopyCommandProperty struct {
	// `CfnDeliveryStream.CopyCommandProperty.DataTableName`.
	DataTableName *string `json:"dataTableName"`
	// `CfnDeliveryStream.CopyCommandProperty.CopyOptions`.
	CopyOptions *string `json:"copyOptions"`
	// `CfnDeliveryStream.CopyCommandProperty.DataTableColumns`.
	DataTableColumns *string `json:"dataTableColumns"`
}

type CfnDeliveryStream_DataFormatConversionConfigurationProperty struct {
	// `CfnDeliveryStream.DataFormatConversionConfigurationProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `CfnDeliveryStream.DataFormatConversionConfigurationProperty.InputFormatConfiguration`.
	InputFormatConfiguration interface{} `json:"inputFormatConfiguration"`
	// `CfnDeliveryStream.DataFormatConversionConfigurationProperty.OutputFormatConfiguration`.
	OutputFormatConfiguration interface{} `json:"outputFormatConfiguration"`
	// `CfnDeliveryStream.DataFormatConversionConfigurationProperty.SchemaConfiguration`.
	SchemaConfiguration interface{} `json:"schemaConfiguration"`
}

type CfnDeliveryStream_DeliveryStreamEncryptionConfigurationInputProperty struct {
	// `CfnDeliveryStream.DeliveryStreamEncryptionConfigurationInputProperty.KeyType`.
	KeyType *string `json:"keyType"`
	// `CfnDeliveryStream.DeliveryStreamEncryptionConfigurationInputProperty.KeyARN`.
	KeyArn *string `json:"keyArn"`
}

type CfnDeliveryStream_DeserializerProperty struct {
	// `CfnDeliveryStream.DeserializerProperty.HiveJsonSerDe`.
	HiveJsonSerDe interface{} `json:"hiveJsonSerDe"`
	// `CfnDeliveryStream.DeserializerProperty.OpenXJsonSerDe`.
	OpenXJsonSerDe interface{} `json:"openXJsonSerDe"`
}

type CfnDeliveryStream_ElasticsearchBufferingHintsProperty struct {
	// `CfnDeliveryStream.ElasticsearchBufferingHintsProperty.IntervalInSeconds`.
	IntervalInSeconds *float64 `json:"intervalInSeconds"`
	// `CfnDeliveryStream.ElasticsearchBufferingHintsProperty.SizeInMBs`.
	SizeInMBs *float64 `json:"sizeInMBs"`
}

type CfnDeliveryStream_ElasticsearchDestinationConfigurationProperty struct {
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.IndexName`.
	IndexName *string `json:"indexName"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.S3Configuration`.
	S3Configuration interface{} `json:"s3Configuration"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.BufferingHints`.
	BufferingHints interface{} `json:"bufferingHints"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.CloudWatchLoggingOptions`.
	CloudWatchLoggingOptions interface{} `json:"cloudWatchLoggingOptions"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.ClusterEndpoint`.
	ClusterEndpoint *string `json:"clusterEndpoint"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.DomainARN`.
	DomainArn *string `json:"domainArn"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.IndexRotationPeriod`.
	IndexRotationPeriod *string `json:"indexRotationPeriod"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.ProcessingConfiguration`.
	ProcessingConfiguration interface{} `json:"processingConfiguration"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.RetryOptions`.
	RetryOptions interface{} `json:"retryOptions"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.S3BackupMode`.
	S3BackupMode *string `json:"s3BackupMode"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.TypeName`.
	TypeName *string `json:"typeName"`
	// `CfnDeliveryStream.ElasticsearchDestinationConfigurationProperty.VpcConfiguration`.
	VpcConfiguration interface{} `json:"vpcConfiguration"`
}

type CfnDeliveryStream_ElasticsearchRetryOptionsProperty struct {
	// `CfnDeliveryStream.ElasticsearchRetryOptionsProperty.DurationInSeconds`.
	DurationInSeconds *float64 `json:"durationInSeconds"`
}

type CfnDeliveryStream_EncryptionConfigurationProperty struct {
	// `CfnDeliveryStream.EncryptionConfigurationProperty.KMSEncryptionConfig`.
	KmsEncryptionConfig interface{} `json:"kmsEncryptionConfig"`
	// `CfnDeliveryStream.EncryptionConfigurationProperty.NoEncryptionConfig`.
	NoEncryptionConfig *string `json:"noEncryptionConfig"`
}

type CfnDeliveryStream_ExtendedS3DestinationConfigurationProperty struct {
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.BucketARN`.
	BucketArn *string `json:"bucketArn"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.BufferingHints`.
	BufferingHints interface{} `json:"bufferingHints"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.CloudWatchLoggingOptions`.
	CloudWatchLoggingOptions interface{} `json:"cloudWatchLoggingOptions"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.CompressionFormat`.
	CompressionFormat *string `json:"compressionFormat"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.DataFormatConversionConfiguration`.
	DataFormatConversionConfiguration interface{} `json:"dataFormatConversionConfiguration"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.EncryptionConfiguration`.
	EncryptionConfiguration interface{} `json:"encryptionConfiguration"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.ErrorOutputPrefix`.
	ErrorOutputPrefix *string `json:"errorOutputPrefix"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.ProcessingConfiguration`.
	ProcessingConfiguration interface{} `json:"processingConfiguration"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.S3BackupConfiguration`.
	S3BackupConfiguration interface{} `json:"s3BackupConfiguration"`
	// `CfnDeliveryStream.ExtendedS3DestinationConfigurationProperty.S3BackupMode`.
	S3BackupMode *string `json:"s3BackupMode"`
}

type CfnDeliveryStream_HiveJsonSerDeProperty struct {
	// `CfnDeliveryStream.HiveJsonSerDeProperty.TimestampFormats`.
	TimestampFormats *[]*string `json:"timestampFormats"`
}

type CfnDeliveryStream_HttpEndpointCommonAttributeProperty struct {
	// `CfnDeliveryStream.HttpEndpointCommonAttributeProperty.AttributeName`.
	AttributeName *string `json:"attributeName"`
	// `CfnDeliveryStream.HttpEndpointCommonAttributeProperty.AttributeValue`.
	AttributeValue *string `json:"attributeValue"`
}

type CfnDeliveryStream_HttpEndpointConfigurationProperty struct {
	// `CfnDeliveryStream.HttpEndpointConfigurationProperty.Url`.
	Url *string `json:"url"`
	// `CfnDeliveryStream.HttpEndpointConfigurationProperty.AccessKey`.
	AccessKey *string `json:"accessKey"`
	// `CfnDeliveryStream.HttpEndpointConfigurationProperty.Name`.
	Name *string `json:"name"`
}

type CfnDeliveryStream_HttpEndpointDestinationConfigurationProperty struct {
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.EndpointConfiguration`.
	EndpointConfiguration interface{} `json:"endpointConfiguration"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.S3Configuration`.
	S3Configuration interface{} `json:"s3Configuration"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.BufferingHints`.
	BufferingHints interface{} `json:"bufferingHints"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.CloudWatchLoggingOptions`.
	CloudWatchLoggingOptions interface{} `json:"cloudWatchLoggingOptions"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.ProcessingConfiguration`.
	ProcessingConfiguration interface{} `json:"processingConfiguration"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.RequestConfiguration`.
	RequestConfiguration interface{} `json:"requestConfiguration"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.RetryOptions`.
	RetryOptions interface{} `json:"retryOptions"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnDeliveryStream.HttpEndpointDestinationConfigurationProperty.S3BackupMode`.
	S3BackupMode *string `json:"s3BackupMode"`
}

type CfnDeliveryStream_HttpEndpointRequestConfigurationProperty struct {
	// `CfnDeliveryStream.HttpEndpointRequestConfigurationProperty.CommonAttributes`.
	CommonAttributes interface{} `json:"commonAttributes"`
	// `CfnDeliveryStream.HttpEndpointRequestConfigurationProperty.ContentEncoding`.
	ContentEncoding *string `json:"contentEncoding"`
}

type CfnDeliveryStream_InputFormatConfigurationProperty struct {
	// `CfnDeliveryStream.InputFormatConfigurationProperty.Deserializer`.
	Deserializer interface{} `json:"deserializer"`
}

type CfnDeliveryStream_KMSEncryptionConfigProperty struct {
	// `CfnDeliveryStream.KMSEncryptionConfigProperty.AWSKMSKeyARN`.
	AwskmsKeyArn *string `json:"awskmsKeyArn"`
}

type CfnDeliveryStream_KinesisStreamSourceConfigurationProperty struct {
	// `CfnDeliveryStream.KinesisStreamSourceConfigurationProperty.KinesisStreamARN`.
	KinesisStreamArn *string `json:"kinesisStreamArn"`
	// `CfnDeliveryStream.KinesisStreamSourceConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
}

type CfnDeliveryStream_OpenXJsonSerDeProperty struct {
	// `CfnDeliveryStream.OpenXJsonSerDeProperty.CaseInsensitive`.
	CaseInsensitive interface{} `json:"caseInsensitive"`
	// `CfnDeliveryStream.OpenXJsonSerDeProperty.ColumnToJsonKeyMappings`.
	ColumnToJsonKeyMappings interface{} `json:"columnToJsonKeyMappings"`
	// `CfnDeliveryStream.OpenXJsonSerDeProperty.ConvertDotsInJsonKeysToUnderscores`.
	ConvertDotsInJsonKeysToUnderscores interface{} `json:"convertDotsInJsonKeysToUnderscores"`
}

type CfnDeliveryStream_OrcSerDeProperty struct {
	// `CfnDeliveryStream.OrcSerDeProperty.BlockSizeBytes`.
	BlockSizeBytes *float64 `json:"blockSizeBytes"`
	// `CfnDeliveryStream.OrcSerDeProperty.BloomFilterColumns`.
	BloomFilterColumns *[]*string `json:"bloomFilterColumns"`
	// `CfnDeliveryStream.OrcSerDeProperty.BloomFilterFalsePositiveProbability`.
	BloomFilterFalsePositiveProbability *float64 `json:"bloomFilterFalsePositiveProbability"`
	// `CfnDeliveryStream.OrcSerDeProperty.Compression`.
	Compression *string `json:"compression"`
	// `CfnDeliveryStream.OrcSerDeProperty.DictionaryKeyThreshold`.
	DictionaryKeyThreshold *float64 `json:"dictionaryKeyThreshold"`
	// `CfnDeliveryStream.OrcSerDeProperty.EnablePadding`.
	EnablePadding interface{} `json:"enablePadding"`
	// `CfnDeliveryStream.OrcSerDeProperty.FormatVersion`.
	FormatVersion *string `json:"formatVersion"`
	// `CfnDeliveryStream.OrcSerDeProperty.PaddingTolerance`.
	PaddingTolerance *float64 `json:"paddingTolerance"`
	// `CfnDeliveryStream.OrcSerDeProperty.RowIndexStride`.
	RowIndexStride *float64 `json:"rowIndexStride"`
	// `CfnDeliveryStream.OrcSerDeProperty.StripeSizeBytes`.
	StripeSizeBytes *float64 `json:"stripeSizeBytes"`
}

type CfnDeliveryStream_OutputFormatConfigurationProperty struct {
	// `CfnDeliveryStream.OutputFormatConfigurationProperty.Serializer`.
	Serializer interface{} `json:"serializer"`
}

type CfnDeliveryStream_ParquetSerDeProperty struct {
	// `CfnDeliveryStream.ParquetSerDeProperty.BlockSizeBytes`.
	BlockSizeBytes *float64 `json:"blockSizeBytes"`
	// `CfnDeliveryStream.ParquetSerDeProperty.Compression`.
	Compression *string `json:"compression"`
	// `CfnDeliveryStream.ParquetSerDeProperty.EnableDictionaryCompression`.
	EnableDictionaryCompression interface{} `json:"enableDictionaryCompression"`
	// `CfnDeliveryStream.ParquetSerDeProperty.MaxPaddingBytes`.
	MaxPaddingBytes *float64 `json:"maxPaddingBytes"`
	// `CfnDeliveryStream.ParquetSerDeProperty.PageSizeBytes`.
	PageSizeBytes *float64 `json:"pageSizeBytes"`
	// `CfnDeliveryStream.ParquetSerDeProperty.WriterVersion`.
	WriterVersion *string `json:"writerVersion"`
}

type CfnDeliveryStream_ProcessingConfigurationProperty struct {
	// `CfnDeliveryStream.ProcessingConfigurationProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `CfnDeliveryStream.ProcessingConfigurationProperty.Processors`.
	Processors interface{} `json:"processors"`
}

type CfnDeliveryStream_ProcessorParameterProperty struct {
	// `CfnDeliveryStream.ProcessorParameterProperty.ParameterName`.
	ParameterName *string `json:"parameterName"`
	// `CfnDeliveryStream.ProcessorParameterProperty.ParameterValue`.
	ParameterValue *string `json:"parameterValue"`
}

type CfnDeliveryStream_ProcessorProperty struct {
	// `CfnDeliveryStream.ProcessorProperty.Type`.
	Type *string `json:"type"`
	// `CfnDeliveryStream.ProcessorProperty.Parameters`.
	Parameters interface{} `json:"parameters"`
}

type CfnDeliveryStream_RedshiftDestinationConfigurationProperty struct {
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.ClusterJDBCURL`.
	ClusterJdbcurl *string `json:"clusterJdbcurl"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.CopyCommand`.
	CopyCommand interface{} `json:"copyCommand"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.Password`.
	Password *string `json:"password"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.S3Configuration`.
	S3Configuration interface{} `json:"s3Configuration"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.Username`.
	Username *string `json:"username"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.CloudWatchLoggingOptions`.
	CloudWatchLoggingOptions interface{} `json:"cloudWatchLoggingOptions"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.ProcessingConfiguration`.
	ProcessingConfiguration interface{} `json:"processingConfiguration"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.RetryOptions`.
	RetryOptions interface{} `json:"retryOptions"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.S3BackupConfiguration`.
	S3BackupConfiguration interface{} `json:"s3BackupConfiguration"`
	// `CfnDeliveryStream.RedshiftDestinationConfigurationProperty.S3BackupMode`.
	S3BackupMode *string `json:"s3BackupMode"`
}

type CfnDeliveryStream_RedshiftRetryOptionsProperty struct {
	// `CfnDeliveryStream.RedshiftRetryOptionsProperty.DurationInSeconds`.
	DurationInSeconds *float64 `json:"durationInSeconds"`
}

type CfnDeliveryStream_RetryOptionsProperty struct {
	// `CfnDeliveryStream.RetryOptionsProperty.DurationInSeconds`.
	DurationInSeconds *float64 `json:"durationInSeconds"`
}

type CfnDeliveryStream_S3DestinationConfigurationProperty struct {
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.BucketARN`.
	BucketArn *string `json:"bucketArn"`
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.BufferingHints`.
	BufferingHints interface{} `json:"bufferingHints"`
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.CloudWatchLoggingOptions`.
	CloudWatchLoggingOptions interface{} `json:"cloudWatchLoggingOptions"`
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.CompressionFormat`.
	CompressionFormat *string `json:"compressionFormat"`
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.EncryptionConfiguration`.
	EncryptionConfiguration interface{} `json:"encryptionConfiguration"`
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.ErrorOutputPrefix`.
	ErrorOutputPrefix *string `json:"errorOutputPrefix"`
	// `CfnDeliveryStream.S3DestinationConfigurationProperty.Prefix`.
	Prefix *string `json:"prefix"`
}

type CfnDeliveryStream_SchemaConfigurationProperty struct {
	// `CfnDeliveryStream.SchemaConfigurationProperty.CatalogId`.
	CatalogId *string `json:"catalogId"`
	// `CfnDeliveryStream.SchemaConfigurationProperty.DatabaseName`.
	DatabaseName *string `json:"databaseName"`
	// `CfnDeliveryStream.SchemaConfigurationProperty.Region`.
	Region *string `json:"region"`
	// `CfnDeliveryStream.SchemaConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnDeliveryStream.SchemaConfigurationProperty.TableName`.
	TableName *string `json:"tableName"`
	// `CfnDeliveryStream.SchemaConfigurationProperty.VersionId`.
	VersionId *string `json:"versionId"`
}

type CfnDeliveryStream_SerializerProperty struct {
	// `CfnDeliveryStream.SerializerProperty.OrcSerDe`.
	OrcSerDe interface{} `json:"orcSerDe"`
	// `CfnDeliveryStream.SerializerProperty.ParquetSerDe`.
	ParquetSerDe interface{} `json:"parquetSerDe"`
}

type CfnDeliveryStream_SplunkDestinationConfigurationProperty struct {
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.HECEndpoint`.
	HecEndpoint *string `json:"hecEndpoint"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.HECEndpointType`.
	HecEndpointType *string `json:"hecEndpointType"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.HECToken`.
	HecToken *string `json:"hecToken"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.S3Configuration`.
	S3Configuration interface{} `json:"s3Configuration"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.CloudWatchLoggingOptions`.
	CloudWatchLoggingOptions interface{} `json:"cloudWatchLoggingOptions"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.HECAcknowledgmentTimeoutInSeconds`.
	HecAcknowledgmentTimeoutInSeconds *float64 `json:"hecAcknowledgmentTimeoutInSeconds"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.ProcessingConfiguration`.
	ProcessingConfiguration interface{} `json:"processingConfiguration"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.RetryOptions`.
	RetryOptions interface{} `json:"retryOptions"`
	// `CfnDeliveryStream.SplunkDestinationConfigurationProperty.S3BackupMode`.
	S3BackupMode *string `json:"s3BackupMode"`
}

type CfnDeliveryStream_SplunkRetryOptionsProperty struct {
	// `CfnDeliveryStream.SplunkRetryOptionsProperty.DurationInSeconds`.
	DurationInSeconds *float64 `json:"durationInSeconds"`
}

type CfnDeliveryStream_VpcConfigurationProperty struct {
	// `CfnDeliveryStream.VpcConfigurationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnDeliveryStream.VpcConfigurationProperty.SecurityGroupIds`.
	SecurityGroupIds *[]*string `json:"securityGroupIds"`
	// `CfnDeliveryStream.VpcConfigurationProperty.SubnetIds`.
	SubnetIds *[]*string `json:"subnetIds"`
}

// Properties for defining a `AWS::KinesisFirehose::DeliveryStream`.
type CfnDeliveryStreamProps struct {
	// `AWS::KinesisFirehose::DeliveryStream.DeliveryStreamEncryptionConfigurationInput`.
	DeliveryStreamEncryptionConfigurationInput interface{} `json:"deliveryStreamEncryptionConfigurationInput"`
	// `AWS::KinesisFirehose::DeliveryStream.DeliveryStreamName`.
	DeliveryStreamName *string `json:"deliveryStreamName"`
	// `AWS::KinesisFirehose::DeliveryStream.DeliveryStreamType`.
	DeliveryStreamType *string `json:"deliveryStreamType"`
	// `AWS::KinesisFirehose::DeliveryStream.ElasticsearchDestinationConfiguration`.
	ElasticsearchDestinationConfiguration interface{} `json:"elasticsearchDestinationConfiguration"`
	// `AWS::KinesisFirehose::DeliveryStream.ExtendedS3DestinationConfiguration`.
	ExtendedS3DestinationConfiguration interface{} `json:"extendedS3DestinationConfiguration"`
	// `AWS::KinesisFirehose::DeliveryStream.HttpEndpointDestinationConfiguration`.
	HttpEndpointDestinationConfiguration interface{} `json:"httpEndpointDestinationConfiguration"`
	// `AWS::KinesisFirehose::DeliveryStream.KinesisStreamSourceConfiguration`.
	KinesisStreamSourceConfiguration interface{} `json:"kinesisStreamSourceConfiguration"`
	// `AWS::KinesisFirehose::DeliveryStream.RedshiftDestinationConfiguration`.
	RedshiftDestinationConfiguration interface{} `json:"redshiftDestinationConfiguration"`
	// `AWS::KinesisFirehose::DeliveryStream.S3DestinationConfiguration`.
	S3DestinationConfiguration interface{} `json:"s3DestinationConfiguration"`
	// `AWS::KinesisFirehose::DeliveryStream.SplunkDestinationConfiguration`.
	SplunkDestinationConfiguration interface{} `json:"splunkDestinationConfiguration"`
	// `AWS::KinesisFirehose::DeliveryStream.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}


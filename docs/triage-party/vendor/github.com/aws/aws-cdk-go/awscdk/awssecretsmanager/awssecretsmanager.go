package awssecretsmanager

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awssecretsmanager/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// Options to add a secret attachment to a secret.
// Deprecated: use `secret.attach()` instead
type AttachedSecretOptions struct {
	// The target to attach the secret to.
	// Deprecated: use `secret.attach()` instead
	Target ISecretAttachmentTarget `json:"target"`
}

// The type of service or database that's being associated with the secret.
// Experimental.
type AttachmentTargetType string

const (
	AttachmentTargetType_INSTANCE AttachmentTargetType = "INSTANCE"
	AttachmentTargetType_CLUSTER AttachmentTargetType = "CLUSTER"
	AttachmentTargetType_RDS_DB_PROXY AttachmentTargetType = "RDS_DB_PROXY"
	AttachmentTargetType_REDSHIFT_CLUSTER AttachmentTargetType = "REDSHIFT_CLUSTER"
	AttachmentTargetType_DOCDB_DB_INSTANCE AttachmentTargetType = "DOCDB_DB_INSTANCE"
	AttachmentTargetType_DOCDB_DB_CLUSTER AttachmentTargetType = "DOCDB_DB_CLUSTER"
)

// A CloudFormation `AWS::SecretsManager::ResourcePolicy`.
type CfnResourcePolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	BlockPublicPolicy() interface{}
	SetBlockPublicPolicy(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	ResourcePolicy() interface{}
	SetResourcePolicy(val interface{})
	SecretId() *string
	SetSecretId(val *string)
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

// The jsii proxy struct for CfnResourcePolicy
type jsiiProxy_CfnResourcePolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnResourcePolicy) BlockPublicPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"blockPublicPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) ResourcePolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"resourcePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) SecretId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourcePolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SecretsManager::ResourcePolicy`.
func NewCfnResourcePolicy(scope awscdk.Construct, id *string, props *CfnResourcePolicyProps) CfnResourcePolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnResourcePolicy{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnResourcePolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SecretsManager::ResourcePolicy`.
func NewCfnResourcePolicy_Override(c CfnResourcePolicy, scope awscdk.Construct, id *string, props *CfnResourcePolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnResourcePolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnResourcePolicy) SetBlockPublicPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"blockPublicPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnResourcePolicy) SetResourcePolicy(val interface{}) {
	_jsii_.Set(
		j,
		"resourcePolicy",
		val,
	)
}

func (j *jsiiProxy_CfnResourcePolicy) SetSecretId(val *string) {
	_jsii_.Set(
		j,
		"secretId",
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
func CfnResourcePolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnResourcePolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnResourcePolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnResourcePolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnResourcePolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnResourcePolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnResourcePolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.CfnResourcePolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnResourcePolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnResourcePolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnResourcePolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnResourcePolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnResourcePolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnResourcePolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnResourcePolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnResourcePolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnResourcePolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnResourcePolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnResourcePolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnResourcePolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnResourcePolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnResourcePolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnResourcePolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnResourcePolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnResourcePolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnResourcePolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnResourcePolicy) ToString() *string {
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
func (c *jsiiProxy_CfnResourcePolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnResourcePolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::SecretsManager::ResourcePolicy`.
type CfnResourcePolicyProps struct {
	// `AWS::SecretsManager::ResourcePolicy.ResourcePolicy`.
	ResourcePolicy interface{} `json:"resourcePolicy"`
	// `AWS::SecretsManager::ResourcePolicy.SecretId`.
	SecretId *string `json:"secretId"`
	// `AWS::SecretsManager::ResourcePolicy.BlockPublicPolicy`.
	BlockPublicPolicy interface{} `json:"blockPublicPolicy"`
}

// A CloudFormation `AWS::SecretsManager::RotationSchedule`.
type CfnRotationSchedule interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	HostedRotationLambda() interface{}
	SetHostedRotationLambda(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RotationLambdaArn() *string
	SetRotationLambdaArn(val *string)
	RotationRules() interface{}
	SetRotationRules(val interface{})
	SecretId() *string
	SetSecretId(val *string)
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

// The jsii proxy struct for CfnRotationSchedule
type jsiiProxy_CfnRotationSchedule struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRotationSchedule) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) HostedRotationLambda() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"hostedRotationLambda",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) RotationLambdaArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"rotationLambdaArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) RotationRules() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rotationRules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) SecretId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRotationSchedule) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SecretsManager::RotationSchedule`.
func NewCfnRotationSchedule(scope awscdk.Construct, id *string, props *CfnRotationScheduleProps) CfnRotationSchedule {
	_init_.Initialize()

	j := jsiiProxy_CfnRotationSchedule{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnRotationSchedule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SecretsManager::RotationSchedule`.
func NewCfnRotationSchedule_Override(c CfnRotationSchedule, scope awscdk.Construct, id *string, props *CfnRotationScheduleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnRotationSchedule",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRotationSchedule) SetHostedRotationLambda(val interface{}) {
	_jsii_.Set(
		j,
		"hostedRotationLambda",
		val,
	)
}

func (j *jsiiProxy_CfnRotationSchedule) SetRotationLambdaArn(val *string) {
	_jsii_.Set(
		j,
		"rotationLambdaArn",
		val,
	)
}

func (j *jsiiProxy_CfnRotationSchedule) SetRotationRules(val interface{}) {
	_jsii_.Set(
		j,
		"rotationRules",
		val,
	)
}

func (j *jsiiProxy_CfnRotationSchedule) SetSecretId(val *string) {
	_jsii_.Set(
		j,
		"secretId",
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
func CfnRotationSchedule_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnRotationSchedule",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnRotationSchedule_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnRotationSchedule",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnRotationSchedule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnRotationSchedule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRotationSchedule_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.CfnRotationSchedule",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnRotationSchedule) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnRotationSchedule) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnRotationSchedule) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnRotationSchedule) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnRotationSchedule) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnRotationSchedule) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnRotationSchedule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnRotationSchedule) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnRotationSchedule) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnRotationSchedule) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnRotationSchedule) OnPrepare() {
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
func (c *jsiiProxy_CfnRotationSchedule) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRotationSchedule) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnRotationSchedule) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnRotationSchedule) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnRotationSchedule) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnRotationSchedule) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnRotationSchedule) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRotationSchedule) ToString() *string {
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
func (c *jsiiProxy_CfnRotationSchedule) Validate() *[]*string {
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
func (c *jsiiProxy_CfnRotationSchedule) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnRotationSchedule_HostedRotationLambdaProperty struct {
	// `CfnRotationSchedule.HostedRotationLambdaProperty.RotationType`.
	RotationType *string `json:"rotationType"`
	// `CfnRotationSchedule.HostedRotationLambdaProperty.KmsKeyArn`.
	KmsKeyArn *string `json:"kmsKeyArn"`
	// `CfnRotationSchedule.HostedRotationLambdaProperty.MasterSecretArn`.
	MasterSecretArn *string `json:"masterSecretArn"`
	// `CfnRotationSchedule.HostedRotationLambdaProperty.MasterSecretKmsKeyArn`.
	MasterSecretKmsKeyArn *string `json:"masterSecretKmsKeyArn"`
	// `CfnRotationSchedule.HostedRotationLambdaProperty.RotationLambdaName`.
	RotationLambdaName *string `json:"rotationLambdaName"`
	// `CfnRotationSchedule.HostedRotationLambdaProperty.VpcSecurityGroupIds`.
	VpcSecurityGroupIds *string `json:"vpcSecurityGroupIds"`
	// `CfnRotationSchedule.HostedRotationLambdaProperty.VpcSubnetIds`.
	VpcSubnetIds *string `json:"vpcSubnetIds"`
}

type CfnRotationSchedule_RotationRulesProperty struct {
	// `CfnRotationSchedule.RotationRulesProperty.AutomaticallyAfterDays`.
	AutomaticallyAfterDays *float64 `json:"automaticallyAfterDays"`
}

// Properties for defining a `AWS::SecretsManager::RotationSchedule`.
type CfnRotationScheduleProps struct {
	// `AWS::SecretsManager::RotationSchedule.SecretId`.
	SecretId *string `json:"secretId"`
	// `AWS::SecretsManager::RotationSchedule.HostedRotationLambda`.
	HostedRotationLambda interface{} `json:"hostedRotationLambda"`
	// `AWS::SecretsManager::RotationSchedule.RotationLambdaARN`.
	RotationLambdaArn *string `json:"rotationLambdaArn"`
	// `AWS::SecretsManager::RotationSchedule.RotationRules`.
	RotationRules interface{} `json:"rotationRules"`
}

// A CloudFormation `AWS::SecretsManager::Secret`.
type CfnSecret interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	GenerateSecretString() interface{}
	SetGenerateSecretString(val interface{})
	KmsKeyId() *string
	SetKmsKeyId(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	ReplicaRegions() interface{}
	SetReplicaRegions(val interface{})
	SecretString() *string
	SetSecretString(val *string)
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

// The jsii proxy struct for CfnSecret
type jsiiProxy_CfnSecret struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnSecret) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) GenerateSecretString() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"generateSecretString",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) KmsKeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kmsKeyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) ReplicaRegions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"replicaRegions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) SecretString() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretString",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecret) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SecretsManager::Secret`.
func NewCfnSecret(scope awscdk.Construct, id *string, props *CfnSecretProps) CfnSecret {
	_init_.Initialize()

	j := jsiiProxy_CfnSecret{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnSecret",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SecretsManager::Secret`.
func NewCfnSecret_Override(c CfnSecret, scope awscdk.Construct, id *string, props *CfnSecretProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnSecret",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnSecret) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnSecret) SetGenerateSecretString(val interface{}) {
	_jsii_.Set(
		j,
		"generateSecretString",
		val,
	)
}

func (j *jsiiProxy_CfnSecret) SetKmsKeyId(val *string) {
	_jsii_.Set(
		j,
		"kmsKeyId",
		val,
	)
}

func (j *jsiiProxy_CfnSecret) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnSecret) SetReplicaRegions(val interface{}) {
	_jsii_.Set(
		j,
		"replicaRegions",
		val,
	)
}

func (j *jsiiProxy_CfnSecret) SetSecretString(val *string) {
	_jsii_.Set(
		j,
		"secretString",
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
func CfnSecret_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnSecret",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnSecret_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnSecret",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnSecret_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnSecret",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnSecret_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.CfnSecret",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnSecret) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnSecret) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnSecret) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnSecret) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnSecret) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnSecret) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnSecret) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnSecret) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnSecret) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnSecret) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnSecret) OnPrepare() {
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
func (c *jsiiProxy_CfnSecret) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSecret) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnSecret) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnSecret) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnSecret) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnSecret) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnSecret) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSecret) ToString() *string {
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
func (c *jsiiProxy_CfnSecret) Validate() *[]*string {
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
func (c *jsiiProxy_CfnSecret) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnSecret_GenerateSecretStringProperty struct {
	// `CfnSecret.GenerateSecretStringProperty.ExcludeCharacters`.
	ExcludeCharacters *string `json:"excludeCharacters"`
	// `CfnSecret.GenerateSecretStringProperty.ExcludeLowercase`.
	ExcludeLowercase interface{} `json:"excludeLowercase"`
	// `CfnSecret.GenerateSecretStringProperty.ExcludeNumbers`.
	ExcludeNumbers interface{} `json:"excludeNumbers"`
	// `CfnSecret.GenerateSecretStringProperty.ExcludePunctuation`.
	ExcludePunctuation interface{} `json:"excludePunctuation"`
	// `CfnSecret.GenerateSecretStringProperty.ExcludeUppercase`.
	ExcludeUppercase interface{} `json:"excludeUppercase"`
	// `CfnSecret.GenerateSecretStringProperty.GenerateStringKey`.
	GenerateStringKey *string `json:"generateStringKey"`
	// `CfnSecret.GenerateSecretStringProperty.IncludeSpace`.
	IncludeSpace interface{} `json:"includeSpace"`
	// `CfnSecret.GenerateSecretStringProperty.PasswordLength`.
	PasswordLength *float64 `json:"passwordLength"`
	// `CfnSecret.GenerateSecretStringProperty.RequireEachIncludedType`.
	RequireEachIncludedType interface{} `json:"requireEachIncludedType"`
	// `CfnSecret.GenerateSecretStringProperty.SecretStringTemplate`.
	SecretStringTemplate *string `json:"secretStringTemplate"`
}

type CfnSecret_ReplicaRegionProperty struct {
	// `CfnSecret.ReplicaRegionProperty.Region`.
	Region *string `json:"region"`
	// `CfnSecret.ReplicaRegionProperty.KmsKeyId`.
	KmsKeyId *string `json:"kmsKeyId"`
}

// Properties for defining a `AWS::SecretsManager::Secret`.
type CfnSecretProps struct {
	// `AWS::SecretsManager::Secret.Description`.
	Description *string `json:"description"`
	// `AWS::SecretsManager::Secret.GenerateSecretString`.
	GenerateSecretString interface{} `json:"generateSecretString"`
	// `AWS::SecretsManager::Secret.KmsKeyId`.
	KmsKeyId *string `json:"kmsKeyId"`
	// `AWS::SecretsManager::Secret.Name`.
	Name *string `json:"name"`
	// `AWS::SecretsManager::Secret.ReplicaRegions`.
	ReplicaRegions interface{} `json:"replicaRegions"`
	// `AWS::SecretsManager::Secret.SecretString`.
	SecretString *string `json:"secretString"`
	// `AWS::SecretsManager::Secret.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::SecretsManager::SecretTargetAttachment`.
type CfnSecretTargetAttachment interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	SecretId() *string
	SetSecretId(val *string)
	Stack() awscdk.Stack
	TargetId() *string
	SetTargetId(val *string)
	TargetType() *string
	SetTargetType(val *string)
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

// The jsii proxy struct for CfnSecretTargetAttachment
type jsiiProxy_CfnSecretTargetAttachment struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnSecretTargetAttachment) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) SecretId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) TargetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) TargetType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSecretTargetAttachment) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::SecretsManager::SecretTargetAttachment`.
func NewCfnSecretTargetAttachment(scope awscdk.Construct, id *string, props *CfnSecretTargetAttachmentProps) CfnSecretTargetAttachment {
	_init_.Initialize()

	j := jsiiProxy_CfnSecretTargetAttachment{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::SecretsManager::SecretTargetAttachment`.
func NewCfnSecretTargetAttachment_Override(c CfnSecretTargetAttachment, scope awscdk.Construct, id *string, props *CfnSecretTargetAttachmentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachment",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnSecretTargetAttachment) SetSecretId(val *string) {
	_jsii_.Set(
		j,
		"secretId",
		val,
	)
}

func (j *jsiiProxy_CfnSecretTargetAttachment) SetTargetId(val *string) {
	_jsii_.Set(
		j,
		"targetId",
		val,
	)
}

func (j *jsiiProxy_CfnSecretTargetAttachment) SetTargetType(val *string) {
	_jsii_.Set(
		j,
		"targetType",
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
func CfnSecretTargetAttachment_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachment",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnSecretTargetAttachment_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachment",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnSecretTargetAttachment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnSecretTargetAttachment_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.CfnSecretTargetAttachment",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnSecretTargetAttachment) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnSecretTargetAttachment) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnSecretTargetAttachment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) OnPrepare() {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnSecretTargetAttachment) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) ToString() *string {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) Validate() *[]*string {
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
func (c *jsiiProxy_CfnSecretTargetAttachment) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::SecretsManager::SecretTargetAttachment`.
type CfnSecretTargetAttachmentProps struct {
	// `AWS::SecretsManager::SecretTargetAttachment.SecretId`.
	SecretId *string `json:"secretId"`
	// `AWS::SecretsManager::SecretTargetAttachment.TargetId`.
	TargetId *string `json:"targetId"`
	// `AWS::SecretsManager::SecretTargetAttachment.TargetType`.
	TargetType *string `json:"targetType"`
}

// A hosted rotation.
// Experimental.
type HostedRotation interface {
	awsec2.IConnectable
	Connections() awsec2.Connections
	Bind(secret ISecret, scope constructs.Construct) *CfnRotationSchedule_HostedRotationLambdaProperty
}

// The jsii proxy struct for HostedRotation
type jsiiProxy_HostedRotation struct {
	internal.Type__awsec2IConnectable
}

func (j *jsiiProxy_HostedRotation) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}


// MariaDB Multi User.
// Experimental.
func HostedRotation_MariaDbMultiUser(options *MultiUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"mariaDbMultiUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// MariaDB Single User.
// Experimental.
func HostedRotation_MariaDbSingleUser(options *SingleUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"mariaDbSingleUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// MongoDB Multi User.
// Experimental.
func HostedRotation_MongoDbMultiUser(options *MultiUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"mongoDbMultiUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// MongoDB Single User.
// Experimental.
func HostedRotation_MongoDbSingleUser(options *SingleUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"mongoDbSingleUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// MySQL Multi User.
// Experimental.
func HostedRotation_MysqlMultiUser(options *MultiUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"mysqlMultiUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// MySQL Single User.
// Experimental.
func HostedRotation_MysqlSingleUser(options *SingleUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"mysqlSingleUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Oracle Multi User.
// Experimental.
func HostedRotation_OracleMultiUser(options *MultiUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"oracleMultiUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Oracle Single User.
// Experimental.
func HostedRotation_OracleSingleUser(options *SingleUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"oracleSingleUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// PostgreSQL Multi User.
// Experimental.
func HostedRotation_PostgreSqlMultiUser(options *MultiUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"postgreSqlMultiUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// PostgreSQL Single User.
// Experimental.
func HostedRotation_PostgreSqlSingleUser(options *SingleUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"postgreSqlSingleUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Redshift Multi User.
// Experimental.
func HostedRotation_RedshiftMultiUser(options *MultiUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"redshiftMultiUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Redshift Single User.
// Experimental.
func HostedRotation_RedshiftSingleUser(options *SingleUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"redshiftSingleUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// SQL Server Multi User.
// Experimental.
func HostedRotation_SqlServerMultiUser(options *MultiUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"sqlServerMultiUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// SQL Server Single User.
// Experimental.
func HostedRotation_SqlServerSingleUser(options *SingleUserHostedRotationOptions) HostedRotation {
	_init_.Initialize()

	var returns HostedRotation

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.HostedRotation",
		"sqlServerSingleUser",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Binds this hosted rotation to a secret.
// Experimental.
func (h *jsiiProxy_HostedRotation) Bind(secret ISecret, scope constructs.Construct) *CfnRotationSchedule_HostedRotationLambdaProperty {
	var returns *CfnRotationSchedule_HostedRotationLambdaProperty

	_jsii_.Invoke(
		h,
		"bind",
		[]interface{}{secret, scope},
		&returns,
	)

	return returns
}

// Hosted rotation type.
// Experimental.
type HostedRotationType interface {
	IsMultiUser() *bool
	Name() *string
}

// The jsii proxy struct for HostedRotationType
type jsiiProxy_HostedRotationType struct {
	_ byte // padding
}

func (j *jsiiProxy_HostedRotationType) IsMultiUser() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isMultiUser",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HostedRotationType) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}


func HostedRotationType_MARIADB_MULTI_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"MARIADB_MULTI_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_MARIADB_SINGLE_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"MARIADB_SINGLE_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_MONGODB_MULTI_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"MONGODB_MULTI_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_MONGODB_SINGLE_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"MONGODB_SINGLE_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_MYSQL_MULTI_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"MYSQL_MULTI_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_MYSQL_SINGLE_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"MYSQL_SINGLE_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_ORACLE_MULTI_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"ORACLE_MULTI_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_ORACLE_SINGLE_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"ORACLE_SINGLE_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_POSTGRESQL_MULTI_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"POSTGRESQL_MULTI_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_POSTGRESQL_SINGLE_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"POSTGRESQL_SINGLE_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_REDSHIFT_MULTI_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"REDSHIFT_MULTI_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_REDSHIFT_SINGLE_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"REDSHIFT_SINGLE_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_SQLSERVER_MULTI_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"SQLSERVER_MULTI_USER",
		&returns,
	)
	return returns
}

func HostedRotationType_SQLSERVER_SINGLE_USER() HostedRotationType {
	_init_.Initialize()
	var returns HostedRotationType
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.HostedRotationType",
		"SQLSERVER_SINGLE_USER",
		&returns,
	)
	return returns
}

// A secret in AWS Secrets Manager.
// Experimental.
type ISecret interface {
	awscdk.IResource
	// Adds a rotation schedule to the secret.
	// Experimental.
	AddRotationSchedule(id *string, options *RotationScheduleOptions) RotationSchedule
	// Adds a statement to the IAM resource policy associated with this secret.
	//
	// If this secret was created in this stack, a resource policy will be
	// automatically created upon the first call to `addToResourcePolicy`. If
	// the secret is imported, then this is a no-op.
	// Experimental.
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	// Attach a target to this secret.
	//
	// Returns: An attached secret
	// Experimental.
	Attach(target ISecretAttachmentTarget) ISecret
	// Denies the `DeleteSecret` action to all principals within the current account.
	// Experimental.
	DenyAccountRootDelete()
	// Grants reading the secret value to some role.
	// Experimental.
	GrantRead(grantee awsiam.IGrantable, versionStages *[]*string) awsiam.Grant
	// Grants writing and updating the secret value to some role.
	// Experimental.
	GrantWrite(grantee awsiam.IGrantable) awsiam.Grant
	// Interpret the secret as a JSON object and return a field's value from it as a `SecretValue`.
	// Experimental.
	SecretValueFromJson(key *string) awscdk.SecretValue
	// The customer-managed encryption key that is used to encrypt this secret, if any.
	//
	// When not specified, the default
	// KMS key for the account and region is being used.
	// Experimental.
	EncryptionKey() awskms.IKey
	// The ARN of the secret in AWS Secrets Manager.
	//
	// Will return the full ARN if available, otherwise a partial arn.
	// For secrets imported by the deprecated `fromSecretName`, it will return the `secretName`.
	// Experimental.
	SecretArn() *string
	// The full ARN of the secret in AWS Secrets Manager, which is the ARN including the Secrets Manager-supplied 6-character suffix.
	//
	// This is equal to `secretArn` in most cases, but is undefined when a full ARN is not available (e.g., secrets imported by name).
	// Experimental.
	SecretFullArn() *string
	// The name of the secret.
	//
	// For "owned" secrets, this will be the full resource name (secret name + suffix), unless the
	// '@aws-cdk/aws-secretsmanager:parseOwnedSecretName' feature flag is set.
	// Experimental.
	SecretName() *string
	// Retrieve the value of the stored secret as a `SecretValue`.
	// Experimental.
	SecretValue() awscdk.SecretValue
}

// The jsii proxy for ISecret
type jsiiProxy_ISecret struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_ISecret) AddRotationSchedule(id *string, options *RotationScheduleOptions) RotationSchedule {
	var returns RotationSchedule

	_jsii_.Invoke(
		i,
		"addRotationSchedule",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ISecret) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		i,
		"addToResourcePolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ISecret) Attach(target ISecretAttachmentTarget) ISecret {
	var returns ISecret

	_jsii_.Invoke(
		i,
		"attach",
		[]interface{}{target},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ISecret) DenyAccountRootDelete() {
	_jsii_.InvokeVoid(
		i,
		"denyAccountRootDelete",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ISecret) GrantRead(grantee awsiam.IGrantable, versionStages *[]*string) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantRead",
		[]interface{}{grantee, versionStages},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ISecret) GrantWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantWrite",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ISecret) SecretValueFromJson(key *string) awscdk.SecretValue {
	var returns awscdk.SecretValue

	_jsii_.Invoke(
		i,
		"secretValueFromJson",
		[]interface{}{key},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_ISecret) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISecret) SecretArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISecret) SecretFullArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretFullArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISecret) SecretName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISecret) SecretValue() awscdk.SecretValue {
	var returns awscdk.SecretValue
	_jsii_.Get(
		j,
		"secretValue",
		&returns,
	)
	return returns
}

// A secret attachment target.
// Experimental.
type ISecretAttachmentTarget interface {
	// Renders the target specifications.
	// Experimental.
	AsSecretAttachmentTarget() *SecretAttachmentTargetProps
}

// The jsii proxy for ISecretAttachmentTarget
type jsiiProxy_ISecretAttachmentTarget struct {
	_ byte // padding
}

func (i *jsiiProxy_ISecretAttachmentTarget) AsSecretAttachmentTarget() *SecretAttachmentTargetProps {
	var returns *SecretAttachmentTargetProps

	_jsii_.Invoke(
		i,
		"asSecretAttachmentTarget",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type ISecretTargetAttachment interface {
	ISecret
	// Same as `secretArn`.
	// Experimental.
	SecretTargetAttachmentSecretArn() *string
}

// The jsii proxy for ISecretTargetAttachment
type jsiiProxy_ISecretTargetAttachment struct {
	jsiiProxy_ISecret
}

func (j *jsiiProxy_ISecretTargetAttachment) SecretTargetAttachmentSecretArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretTargetAttachmentSecretArn",
		&returns,
	)
	return returns
}

// Multi user hosted rotation options.
// Experimental.
type MultiUserHostedRotationOptions struct {
	// A name for the Lambda created to rotate the secret.
	// Experimental.
	FunctionName *string `json:"functionName"`
	// A list of security groups for the Lambda created to rotate the secret.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The VPC where the Lambda rotation function will run.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// The type of subnets in the VPC where the Lambda rotation function will run.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// The master secret for a multi user rotation scheme.
	// Experimental.
	MasterSecret ISecret `json:"masterSecret"`
}

// Secret replica region.
// Experimental.
type ReplicaRegion struct {
	// The name of the region.
	// Experimental.
	Region *string `json:"region"`
	// The customer-managed encryption key to use for encrypting the secret value.
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
}

// Secret Resource Policy.
// Experimental.
type ResourcePolicy interface {
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

// The jsii proxy struct for ResourcePolicy
type jsiiProxy_ResourcePolicy struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_ResourcePolicy) Document() awsiam.PolicyDocument {
	var returns awsiam.PolicyDocument
	_jsii_.Get(
		j,
		"document",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourcePolicy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourcePolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourcePolicy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourcePolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewResourcePolicy(scope constructs.Construct, id *string, props *ResourcePolicyProps) ResourcePolicy {
	_init_.Initialize()

	j := jsiiProxy_ResourcePolicy{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.ResourcePolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewResourcePolicy_Override(r ResourcePolicy, scope constructs.Construct, id *string, props *ResourcePolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.ResourcePolicy",
		[]interface{}{scope, id, props},
		r,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func ResourcePolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.ResourcePolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ResourcePolicy_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.ResourcePolicy",
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
func (r *jsiiProxy_ResourcePolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_ResourcePolicy) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_ResourcePolicy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_ResourcePolicy) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_ResourcePolicy) OnPrepare() {
	_jsii_.InvokeVoid(
		r,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (r *jsiiProxy_ResourcePolicy) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
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
func (r *jsiiProxy_ResourcePolicy) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_ResourcePolicy) Prepare() {
	_jsii_.InvokeVoid(
		r,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (r *jsiiProxy_ResourcePolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_ResourcePolicy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_ResourcePolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for a ResourcePolicy.
// Experimental.
type ResourcePolicyProps struct {
	// The secret to attach a resource-based permissions policy.
	// Experimental.
	Secret ISecret `json:"secret"`
}

// A rotation schedule.
// Experimental.
type RotationSchedule interface {
	awscdk.Resource
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

// The jsii proxy struct for RotationSchedule
type jsiiProxy_RotationSchedule struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_RotationSchedule) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RotationSchedule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RotationSchedule) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RotationSchedule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRotationSchedule(scope constructs.Construct, id *string, props *RotationScheduleProps) RotationSchedule {
	_init_.Initialize()

	j := jsiiProxy_RotationSchedule{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.RotationSchedule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRotationSchedule_Override(r RotationSchedule, scope constructs.Construct, id *string, props *RotationScheduleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.RotationSchedule",
		[]interface{}{scope, id, props},
		r,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func RotationSchedule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.RotationSchedule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func RotationSchedule_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.RotationSchedule",
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
func (r *jsiiProxy_RotationSchedule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_RotationSchedule) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_RotationSchedule) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_RotationSchedule) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_RotationSchedule) OnPrepare() {
	_jsii_.InvokeVoid(
		r,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (r *jsiiProxy_RotationSchedule) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
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
func (r *jsiiProxy_RotationSchedule) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_RotationSchedule) Prepare() {
	_jsii_.InvokeVoid(
		r,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (r *jsiiProxy_RotationSchedule) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_RotationSchedule) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_RotationSchedule) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options to add a rotation schedule to a secret.
// Experimental.
type RotationScheduleOptions struct {
	// Specifies the number of days after the previous rotation before Secrets Manager triggers the next automatic rotation.
	// Experimental.
	AutomaticallyAfter awscdk.Duration `json:"automaticallyAfter"`
	// Hosted rotation.
	// Experimental.
	HostedRotation HostedRotation `json:"hostedRotation"`
	// A Lambda function that can rotate the secret.
	// Experimental.
	RotationLambda awslambda.IFunction `json:"rotationLambda"`
}

// Construction properties for a RotationSchedule.
// Experimental.
type RotationScheduleProps struct {
	// Specifies the number of days after the previous rotation before Secrets Manager triggers the next automatic rotation.
	// Experimental.
	AutomaticallyAfter awscdk.Duration `json:"automaticallyAfter"`
	// Hosted rotation.
	// Experimental.
	HostedRotation HostedRotation `json:"hostedRotation"`
	// A Lambda function that can rotate the secret.
	// Experimental.
	RotationLambda awslambda.IFunction `json:"rotationLambda"`
	// The secret to rotate.
	//
	// If hosted rotation is used, this must be a JSON string with the following format:
	//
	// ```
	// {
	//    "engine": <required: database engine>,
	//    "host": <required: instance host name>,
	//    "username": <required: username>,
	//    "password": <required: password>,
	//    "dbname": <optional: database name>,
	//    "port": <optional: if not specified, default port will be used>,
	//    "masterarn": <required for multi user rotation: the arn of the master secret which will be used to create users/change passwords>
	// }
	// ```
	//
	// This is typically the case for a secret referenced from an `AWS::SecretsManager::SecretTargetAttachment`
	// or an `ISecret` returned by the `attach()` method of `Secret`.
	// Experimental.
	Secret ISecret `json:"secret"`
}

// Creates a new secret in AWS SecretsManager.
// Experimental.
type Secret interface {
	awscdk.Resource
	ISecret
	ArnForPolicies() *string
	AutoCreatePolicy() *bool
	EncryptionKey() awskms.IKey
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	SecretArn() *string
	SecretFullArn() *string
	SecretName() *string
	SecretValue() awscdk.SecretValue
	Stack() awscdk.Stack
	AddReplicaRegion(region *string, encryptionKey awskms.IKey)
	AddRotationSchedule(id *string, options *RotationScheduleOptions) RotationSchedule
	AddTargetAttachment(id *string, options *AttachedSecretOptions) SecretTargetAttachment
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	Attach(target ISecretAttachmentTarget) ISecret
	DenyAccountRootDelete()
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantRead(grantee awsiam.IGrantable, versionStages *[]*string) awsiam.Grant
	GrantWrite(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	SecretValueFromJson(jsonField *string) awscdk.SecretValue
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Secret
type jsiiProxy_Secret struct {
	internal.Type__awscdkResource
	jsiiProxy_ISecret
}

func (j *jsiiProxy_Secret) ArnForPolicies() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arnForPolicies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) SecretArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) SecretFullArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretFullArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) SecretName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) SecretValue() awscdk.SecretValue {
	var returns awscdk.SecretValue
	_jsii_.Get(
		j,
		"secretValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewSecret(scope constructs.Construct, id *string, props *SecretProps) Secret {
	_init_.Initialize()

	j := jsiiProxy_Secret{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.Secret",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSecret_Override(s Secret, scope constructs.Construct, id *string, props *SecretProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.Secret",
		[]interface{}{scope, id, props},
		s,
	)
}

// Deprecated: use `fromSecretCompleteArn` or `fromSecretPartialArn`
func Secret_FromSecretArn(scope constructs.Construct, id *string, secretArn *string) ISecret {
	_init_.Initialize()

	var returns ISecret

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"fromSecretArn",
		[]interface{}{scope, id, secretArn},
		&returns,
	)

	return returns
}

// Import an existing secret into the Stack.
// Experimental.
func Secret_FromSecretAttributes(scope constructs.Construct, id *string, attrs *SecretAttributes) ISecret {
	_init_.Initialize()

	var returns ISecret

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"fromSecretAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Imports a secret by complete ARN.
//
// The complete ARN is the ARN with the Secrets Manager-supplied suffix.
// Experimental.
func Secret_FromSecretCompleteArn(scope constructs.Construct, id *string, secretCompleteArn *string) ISecret {
	_init_.Initialize()

	var returns ISecret

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"fromSecretCompleteArn",
		[]interface{}{scope, id, secretCompleteArn},
		&returns,
	)

	return returns
}

// Imports a secret by secret name;
//
// the ARN of the Secret will be set to the secret name.
// A secret with this name must exist in the same account & region.
// Deprecated: use `fromSecretNameV2`
func Secret_FromSecretName(scope constructs.Construct, id *string, secretName *string) ISecret {
	_init_.Initialize()

	var returns ISecret

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"fromSecretName",
		[]interface{}{scope, id, secretName},
		&returns,
	)

	return returns
}

// Imports a secret by secret name.
//
// A secret with this name must exist in the same account & region.
// Replaces the deprecated `fromSecretName`.
// Experimental.
func Secret_FromSecretNameV2(scope constructs.Construct, id *string, secretName *string) ISecret {
	_init_.Initialize()

	var returns ISecret

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"fromSecretNameV2",
		[]interface{}{scope, id, secretName},
		&returns,
	)

	return returns
}

// Imports a secret by partial ARN.
//
// The partial ARN is the ARN without the Secrets Manager-supplied suffix.
// Experimental.
func Secret_FromSecretPartialArn(scope constructs.Construct, id *string, secretPartialArn *string) ISecret {
	_init_.Initialize()

	var returns ISecret

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"fromSecretPartialArn",
		[]interface{}{scope, id, secretPartialArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Secret_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Secret_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.Secret",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a replica region for the secret.
// Experimental.
func (s *jsiiProxy_Secret) AddReplicaRegion(region *string, encryptionKey awskms.IKey) {
	_jsii_.InvokeVoid(
		s,
		"addReplicaRegion",
		[]interface{}{region, encryptionKey},
	)
}

// Adds a rotation schedule to the secret.
// Experimental.
func (s *jsiiProxy_Secret) AddRotationSchedule(id *string, options *RotationScheduleOptions) RotationSchedule {
	var returns RotationSchedule

	_jsii_.Invoke(
		s,
		"addRotationSchedule",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a target attachment to the secret.
//
// Returns: an AttachedSecret
// Deprecated: use `attach()` instead
func (s *jsiiProxy_Secret) AddTargetAttachment(id *string, options *AttachedSecretOptions) SecretTargetAttachment {
	var returns SecretTargetAttachment

	_jsii_.Invoke(
		s,
		"addTargetAttachment",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a statement to the IAM resource policy associated with this secret.
//
// If this secret was created in this stack, a resource policy will be
// automatically created upon the first call to `addToResourcePolicy`. If
// the secret is imported, then this is a no-op.
// Experimental.
func (s *jsiiProxy_Secret) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		s,
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
func (s *jsiiProxy_Secret) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attach a target to this secret.
//
// Returns: An attached secret
// Experimental.
func (s *jsiiProxy_Secret) Attach(target ISecretAttachmentTarget) ISecret {
	var returns ISecret

	_jsii_.Invoke(
		s,
		"attach",
		[]interface{}{target},
		&returns,
	)

	return returns
}

// Denies the `DeleteSecret` action to all principals within the current account.
// Experimental.
func (s *jsiiProxy_Secret) DenyAccountRootDelete() {
	_jsii_.InvokeVoid(
		s,
		"denyAccountRootDelete",
		nil, // no parameters
	)
}

// Experimental.
func (s *jsiiProxy_Secret) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_Secret) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_Secret) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grants reading the secret value to some role.
// Experimental.
func (s *jsiiProxy_Secret) GrantRead(grantee awsiam.IGrantable, versionStages *[]*string) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantRead",
		[]interface{}{grantee, versionStages},
		&returns,
	)

	return returns
}

// Grants writing and updating the secret value to some role.
// Experimental.
func (s *jsiiProxy_Secret) GrantWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantWrite",
		[]interface{}{grantee},
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
func (s *jsiiProxy_Secret) OnPrepare() {
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
func (s *jsiiProxy_Secret) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_Secret) OnValidate() *[]*string {
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
func (s *jsiiProxy_Secret) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Interpret the secret as a JSON object and return a field's value from it as a `SecretValue`.
// Experimental.
func (s *jsiiProxy_Secret) SecretValueFromJson(jsonField *string) awscdk.SecretValue {
	var returns awscdk.SecretValue

	_jsii_.Invoke(
		s,
		"secretValueFromJson",
		[]interface{}{jsonField},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_Secret) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_Secret) ToString() *string {
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
// Experimental.
func (s *jsiiProxy_Secret) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Attachment target specifications.
// Experimental.
type SecretAttachmentTargetProps struct {
	// The id of the target to attach the secret to.
	// Experimental.
	TargetId *string `json:"targetId"`
	// The type of the target to attach the secret to.
	// Experimental.
	TargetType AttachmentTargetType `json:"targetType"`
}

// Attributes required to import an existing secret into the Stack.
//
// One ARN format (`secretArn`, `secretCompleteArn`, `secretPartialArn`) must be provided.
// Experimental.
type SecretAttributes struct {
	// The encryption key that is used to encrypt the secret, unless the default SecretsManager key is used.
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// The ARN of the secret in SecretsManager.
	//
	// Cannot be used with `secretCompleteArn` or `secretPartialArn`.
	// Deprecated: use `secretCompleteArn` or `secretPartialArn` instead.
	SecretArn *string `json:"secretArn"`
	// The complete ARN of the secret in SecretsManager.
	//
	// This is the ARN including the Secrets Manager 6-character suffix.
	// Cannot be used with `secretArn` or `secretPartialArn`.
	// Experimental.
	SecretCompleteArn *string `json:"secretCompleteArn"`
	// The partial ARN of the secret in SecretsManager.
	//
	// This is the ARN without the Secrets Manager 6-character suffix.
	// Cannot be used with `secretArn` or `secretCompleteArn`.
	// Experimental.
	SecretPartialArn *string `json:"secretPartialArn"`
}

// The properties required to create a new secret in AWS Secrets Manager.
// Experimental.
type SecretProps struct {
	// An optional, human-friendly description of the secret.
	// Experimental.
	Description *string `json:"description"`
	// The customer-managed encryption key to use for encrypting the secret value.
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// Configuration for how to generate a secret value.
	// Experimental.
	GenerateSecretString *SecretStringGenerator `json:"generateSecretString"`
	// Policy to apply when the secret is removed from this stack.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// A list of regions where to replicate this secret.
	// Experimental.
	ReplicaRegions *[]*ReplicaRegion `json:"replicaRegions"`
	// A name for the secret.
	//
	// Note that deleting secrets from SecretsManager does not happen immediately, but after a 7 to
	// 30 days blackout period. During that period, it is not possible to create another secret that shares the same name.
	// Experimental.
	SecretName *string `json:"secretName"`
}

// Secret rotation for a service or database.
// Experimental.
type SecretRotation interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for SecretRotation
type jsiiProxy_SecretRotation struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_SecretRotation) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


// Experimental.
func NewSecretRotation(scope constructs.Construct, id *string, props *SecretRotationProps) SecretRotation {
	_init_.Initialize()

	j := jsiiProxy_SecretRotation{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.SecretRotation",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSecretRotation_Override(s SecretRotation, scope constructs.Construct, id *string, props *SecretRotationProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.SecretRotation",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func SecretRotation_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.SecretRotation",
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
func (s *jsiiProxy_SecretRotation) OnPrepare() {
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
func (s *jsiiProxy_SecretRotation) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_SecretRotation) OnValidate() *[]*string {
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
func (s *jsiiProxy_SecretRotation) Prepare() {
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
func (s *jsiiProxy_SecretRotation) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_SecretRotation) ToString() *string {
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
func (s *jsiiProxy_SecretRotation) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A secret rotation serverless application.
// Experimental.
type SecretRotationApplication interface {
	ApplicationId() *string
	IsMultiUser() *bool
	SemanticVersion() *string
	ApplicationArnForPartition(partition *string) *string
	SemanticVersionForPartition(partition *string) *string
}

// The jsii proxy struct for SecretRotationApplication
type jsiiProxy_SecretRotationApplication struct {
	_ byte // padding
}

func (j *jsiiProxy_SecretRotationApplication) ApplicationId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"applicationId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretRotationApplication) IsMultiUser() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isMultiUser",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretRotationApplication) SemanticVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"semanticVersion",
		&returns,
	)
	return returns
}


// Experimental.
func NewSecretRotationApplication(applicationId *string, semanticVersion *string, options *SecretRotationApplicationOptions) SecretRotationApplication {
	_init_.Initialize()

	j := jsiiProxy_SecretRotationApplication{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		[]interface{}{applicationId, semanticVersion, options},
		&j,
	)

	return &j
}

// Experimental.
func NewSecretRotationApplication_Override(s SecretRotationApplication, applicationId *string, semanticVersion *string, options *SecretRotationApplicationOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		[]interface{}{applicationId, semanticVersion, options},
		s,
	)
}

func SecretRotationApplication_MARIADB_ROTATION_MULTI_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"MARIADB_ROTATION_MULTI_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_MARIADB_ROTATION_SINGLE_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"MARIADB_ROTATION_SINGLE_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_MONGODB_ROTATION_MULTI_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"MONGODB_ROTATION_MULTI_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_MONGODB_ROTATION_SINGLE_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"MONGODB_ROTATION_SINGLE_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_MYSQL_ROTATION_MULTI_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"MYSQL_ROTATION_MULTI_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_MYSQL_ROTATION_SINGLE_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"MYSQL_ROTATION_SINGLE_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_ORACLE_ROTATION_MULTI_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"ORACLE_ROTATION_MULTI_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_ORACLE_ROTATION_SINGLE_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"ORACLE_ROTATION_SINGLE_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_POSTGRES_ROTATION_MULTI_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"POSTGRES_ROTATION_MULTI_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_POSTGRES_ROTATION_SINGLE_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"POSTGRES_ROTATION_SINGLE_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_REDSHIFT_ROTATION_MULTI_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"REDSHIFT_ROTATION_MULTI_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_REDSHIFT_ROTATION_SINGLE_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"REDSHIFT_ROTATION_SINGLE_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_SQLSERVER_ROTATION_MULTI_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"SQLSERVER_ROTATION_MULTI_USER",
		&returns,
	)
	return returns
}

func SecretRotationApplication_SQLSERVER_ROTATION_SINGLE_USER() SecretRotationApplication {
	_init_.Initialize()
	var returns SecretRotationApplication
	_jsii_.StaticGet(
		"monocdk.aws_secretsmanager.SecretRotationApplication",
		"SQLSERVER_ROTATION_SINGLE_USER",
		&returns,
	)
	return returns
}

// Returns the application ARN for the current partition.
//
// Can be used in combination with a `CfnMapping` to automatically select the correct ARN based on the current partition.
// Experimental.
func (s *jsiiProxy_SecretRotationApplication) ApplicationArnForPartition(partition *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"applicationArnForPartition",
		[]interface{}{partition},
		&returns,
	)

	return returns
}

// The semantic version of the app for the current partition.
//
// Can be used in combination with a `CfnMapping` to automatically select the correct version based on the current partition.
// Experimental.
func (s *jsiiProxy_SecretRotationApplication) SemanticVersionForPartition(partition *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"semanticVersionForPartition",
		[]interface{}{partition},
		&returns,
	)

	return returns
}

// Options for a SecretRotationApplication.
// Experimental.
type SecretRotationApplicationOptions struct {
	// Whether the rotation application uses the mutli user scheme.
	// Experimental.
	IsMultiUser *bool `json:"isMultiUser"`
}

// Construction properties for a SecretRotation.
// Experimental.
type SecretRotationProps struct {
	// The serverless application for the rotation.
	// Experimental.
	Application SecretRotationApplication `json:"application"`
	// The secret to rotate. It must be a JSON string with the following format:.
	//
	// ```
	// {
	//    "engine": <required: database engine>,
	//    "host": <required: instance host name>,
	//    "username": <required: username>,
	//    "password": <required: password>,
	//    "dbname": <optional: database name>,
	//    "port": <optional: if not specified, default port will be used>,
	//    "masterarn": <required for multi user rotation: the arn of the master secret which will be used to create users/change passwords>
	// }
	// ```
	//
	// This is typically the case for a secret referenced from an `AWS::SecretsManager::SecretTargetAttachment`
	// or an `ISecret` returned by the `attach()` method of `Secret`.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-secrettargetattachment.html
	//
	// Experimental.
	Secret ISecret `json:"secret"`
	// The target service or database.
	// Experimental.
	Target awsec2.IConnectable `json:"target"`
	// The VPC where the Lambda rotation function will run.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Specifies the number of days after the previous rotation before Secrets Manager triggers the next automatic rotation.
	// Experimental.
	AutomaticallyAfter awscdk.Duration `json:"automaticallyAfter"`
	// Characters which should not appear in the generated password.
	// Experimental.
	ExcludeCharacters *string `json:"excludeCharacters"`
	// The master secret for a multi user rotation scheme.
	// Experimental.
	MasterSecret ISecret `json:"masterSecret"`
	// The security group for the Lambda rotation function.
	// Experimental.
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// The type of subnets in the VPC where the Lambda rotation function will run.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
}

// Configuration to generate secrets such as passwords automatically.
// Experimental.
type SecretStringGenerator struct {
	// A string that includes characters that shouldn't be included in the generated password.
	//
	// The string can be a minimum
	// of ``0`` and a maximum of ``4096`` characters long.
	// Experimental.
	ExcludeCharacters *string `json:"excludeCharacters"`
	// Specifies that the generated password shouldn't include lowercase letters.
	// Experimental.
	ExcludeLowercase *bool `json:"excludeLowercase"`
	// Specifies that the generated password shouldn't include digits.
	// Experimental.
	ExcludeNumbers *bool `json:"excludeNumbers"`
	// Specifies that the generated password shouldn't include punctuation characters.
	// Experimental.
	ExcludePunctuation *bool `json:"excludePunctuation"`
	// Specifies that the generated password shouldn't include uppercase letters.
	// Experimental.
	ExcludeUppercase *bool `json:"excludeUppercase"`
	// The JSON key name that's used to add the generated password to the JSON structure specified by the ``secretStringTemplate`` parameter.
	//
	// If you specify ``generateStringKey`` then ``secretStringTemplate``
	// must be also be specified.
	// Experimental.
	GenerateStringKey *string `json:"generateStringKey"`
	// Specifies that the generated password can include the space character.
	// Experimental.
	IncludeSpace *bool `json:"includeSpace"`
	// The desired length of the generated password.
	// Experimental.
	PasswordLength *float64 `json:"passwordLength"`
	// Specifies whether the generated password must include at least one of every allowed character type.
	// Experimental.
	RequireEachIncludedType *bool `json:"requireEachIncludedType"`
	// A properly structured JSON string that the generated password can be added to.
	//
	// The ``generateStringKey`` is
	// combined with the generated random string and inserted into the JSON structure that's specified by this parameter.
	// The merged JSON string is returned as the completed SecretString of the secret. If you specify ``secretStringTemplate``
	// then ``generateStringKey`` must be also be specified.
	// Experimental.
	SecretStringTemplate *string `json:"secretStringTemplate"`
}

// An attached secret.
// Experimental.
type SecretTargetAttachment interface {
	awscdk.Resource
	ISecret
	ISecretTargetAttachment
	ArnForPolicies() *string
	AutoCreatePolicy() *bool
	EncryptionKey() awskms.IKey
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	SecretArn() *string
	SecretFullArn() *string
	SecretName() *string
	SecretTargetAttachmentSecretArn() *string
	SecretValue() awscdk.SecretValue
	Stack() awscdk.Stack
	AddRotationSchedule(id *string, options *RotationScheduleOptions) RotationSchedule
	AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	Attach(target ISecretAttachmentTarget) ISecret
	DenyAccountRootDelete()
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantRead(grantee awsiam.IGrantable, versionStages *[]*string) awsiam.Grant
	GrantWrite(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	SecretValueFromJson(jsonField *string) awscdk.SecretValue
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for SecretTargetAttachment
type jsiiProxy_SecretTargetAttachment struct {
	internal.Type__awscdkResource
	jsiiProxy_ISecret
	jsiiProxy_ISecretTargetAttachment
}

func (j *jsiiProxy_SecretTargetAttachment) ArnForPolicies() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arnForPolicies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) SecretArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) SecretFullArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretFullArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) SecretName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) SecretTargetAttachmentSecretArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretTargetAttachmentSecretArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) SecretValue() awscdk.SecretValue {
	var returns awscdk.SecretValue
	_jsii_.Get(
		j,
		"secretValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SecretTargetAttachment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewSecretTargetAttachment(scope constructs.Construct, id *string, props *SecretTargetAttachmentProps) SecretTargetAttachment {
	_init_.Initialize()

	j := jsiiProxy_SecretTargetAttachment{}

	_jsii_.Create(
		"monocdk.aws_secretsmanager.SecretTargetAttachment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSecretTargetAttachment_Override(s SecretTargetAttachment, scope constructs.Construct, id *string, props *SecretTargetAttachmentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_secretsmanager.SecretTargetAttachment",
		[]interface{}{scope, id, props},
		s,
	)
}

// Experimental.
func SecretTargetAttachment_FromSecretTargetAttachmentSecretArn(scope constructs.Construct, id *string, secretTargetAttachmentSecretArn *string) ISecretTargetAttachment {
	_init_.Initialize()

	var returns ISecretTargetAttachment

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.SecretTargetAttachment",
		"fromSecretTargetAttachmentSecretArn",
		[]interface{}{scope, id, secretTargetAttachmentSecretArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func SecretTargetAttachment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.SecretTargetAttachment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func SecretTargetAttachment_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_secretsmanager.SecretTargetAttachment",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a rotation schedule to the secret.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) AddRotationSchedule(id *string, options *RotationScheduleOptions) RotationSchedule {
	var returns RotationSchedule

	_jsii_.Invoke(
		s,
		"addRotationSchedule",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a statement to the IAM resource policy associated with this secret.
//
// If this secret was created in this stack, a resource policy will be
// automatically created upon the first call to `addToResourcePolicy`. If
// the secret is imported, then this is a no-op.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) AddToResourcePolicy(statement awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		s,
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
func (s *jsiiProxy_SecretTargetAttachment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attach a target to this secret.
//
// Returns: An attached secret
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) Attach(target ISecretAttachmentTarget) ISecret {
	var returns ISecret

	_jsii_.Invoke(
		s,
		"attach",
		[]interface{}{target},
		&returns,
	)

	return returns
}

// Denies the `DeleteSecret` action to all principals within the current account.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) DenyAccountRootDelete() {
	_jsii_.InvokeVoid(
		s,
		"denyAccountRootDelete",
		nil, // no parameters
	)
}

// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_SecretTargetAttachment) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_SecretTargetAttachment) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grants reading the secret value to some role.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) GrantRead(grantee awsiam.IGrantable, versionStages *[]*string) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantRead",
		[]interface{}{grantee, versionStages},
		&returns,
	)

	return returns
}

// Grants writing and updating the secret value to some role.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) GrantWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantWrite",
		[]interface{}{grantee},
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
func (s *jsiiProxy_SecretTargetAttachment) OnPrepare() {
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
func (s *jsiiProxy_SecretTargetAttachment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_SecretTargetAttachment) OnValidate() *[]*string {
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
func (s *jsiiProxy_SecretTargetAttachment) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Interpret the secret as a JSON object and return a field's value from it as a `SecretValue`.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) SecretValueFromJson(jsonField *string) awscdk.SecretValue {
	var returns awscdk.SecretValue

	_jsii_.Invoke(
		s,
		"secretValueFromJson",
		[]interface{}{jsonField},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) ToString() *string {
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
// Experimental.
func (s *jsiiProxy_SecretTargetAttachment) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for an AttachedSecret.
// Experimental.
type SecretTargetAttachmentProps struct {
	// The target to attach the secret to.
	// Deprecated.
	Target ISecretAttachmentTarget `json:"target"`
	// The secret to attach to the target.
	// Experimental.
	Secret ISecret `json:"secret"`
}

// Single user hosted rotation options.
// Experimental.
type SingleUserHostedRotationOptions struct {
	// A name for the Lambda created to rotate the secret.
	// Experimental.
	FunctionName *string `json:"functionName"`
	// A list of security groups for the Lambda created to rotate the secret.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The VPC where the Lambda rotation function will run.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// The type of subnets in the VPC where the Lambda rotation function will run.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
}


package awskms

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// Defines a display name for a customer master key (CMK) in AWS Key Management Service (AWS KMS).
//
// Using an alias to refer to a key can help you simplify key
// management. For example, when rotating keys, you can just update the alias
// mapping instead of tracking and changing key IDs. For more information, see
// Working with Aliases in the AWS Key Management Service Developer Guide.
//
// You can also add an alias for a key by calling `key.addAlias(alias)`.
// Experimental.
type Alias interface {
	awscdk.Resource
	IAlias
	AliasName() *string
	AliasTargetKey() IKey
	Env() *awscdk.ResourceEnvironment
	KeyArn() *string
	KeyId() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddAlias(alias *string) Alias
	AddToResourcePolicy(statement awsiam.PolicyStatement, allowNoOp *bool) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantDecrypt(grantee awsiam.IGrantable) awsiam.Grant
	GrantEncrypt(grantee awsiam.IGrantable) awsiam.Grant
	GrantEncryptDecrypt(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Alias
type jsiiProxy_Alias struct {
	internal.Type__awscdkResource
	jsiiProxy_IAlias
}

func (j *jsiiProxy_Alias) AliasName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"aliasName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) AliasTargetKey() IKey {
	var returns IKey
	_jsii_.Get(
		j,
		"aliasTargetKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) KeyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) KeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewAlias(scope constructs.Construct, id *string, props *AliasProps) Alias {
	_init_.Initialize()

	j := jsiiProxy_Alias{}

	_jsii_.Create(
		"monocdk.aws_kms.Alias",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewAlias_Override(a Alias, scope constructs.Construct, id *string, props *AliasProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kms.Alias",
		[]interface{}{scope, id, props},
		a,
	)
}

// Import an existing KMS Alias defined outside the CDK app.
// Experimental.
func Alias_FromAliasAttributes(scope constructs.Construct, id *string, attrs *AliasAttributes) IAlias {
	_init_.Initialize()

	var returns IAlias

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Alias",
		"fromAliasAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Import an existing KMS Alias defined outside the CDK app, by the alias name.
//
// This method should be used
// instead of 'fromAliasAttributes' when the underlying KMS Key ARN is not available.
// This Alias will not have a direct reference to the KMS Key, so addAlias and grant* methods are not supported.
// Experimental.
func Alias_FromAliasName(scope constructs.Construct, id *string, aliasName *string) IAlias {
	_init_.Initialize()

	var returns IAlias

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Alias",
		"fromAliasName",
		[]interface{}{scope, id, aliasName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Alias_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Alias",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Alias_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Alias",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Defines a new alias for the key.
// Experimental.
func (a *jsiiProxy_Alias) AddAlias(alias *string) Alias {
	var returns Alias

	_jsii_.Invoke(
		a,
		"addAlias",
		[]interface{}{alias},
		&returns,
	)

	return returns
}

// Adds a statement to the KMS key resource policy.
// Experimental.
func (a *jsiiProxy_Alias) AddToResourcePolicy(statement awsiam.PolicyStatement, allowNoOp *bool) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		a,
		"addToResourcePolicy",
		[]interface{}{statement, allowNoOp},
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
func (a *jsiiProxy_Alias) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_Alias) GeneratePhysicalName() *string {
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
func (a *jsiiProxy_Alias) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (a *jsiiProxy_Alias) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the indicated permissions on this key to the given principal.
// Experimental.
func (a *jsiiProxy_Alias) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
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

// Grant decryption permissions using this key to the given principal.
// Experimental.
func (a *jsiiProxy_Alias) GrantDecrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grantDecrypt",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant encryption permissions using this key to the given principal.
// Experimental.
func (a *jsiiProxy_Alias) GrantEncrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grantEncrypt",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant encryption and decryption permissions using this key to the given principal.
// Experimental.
func (a *jsiiProxy_Alias) GrantEncryptDecrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grantEncryptDecrypt",
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
func (a *jsiiProxy_Alias) OnPrepare() {
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
func (a *jsiiProxy_Alias) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_Alias) OnValidate() *[]*string {
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
func (a *jsiiProxy_Alias) Prepare() {
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
func (a *jsiiProxy_Alias) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_Alias) ToString() *string {
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
func (a *jsiiProxy_Alias) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties of a reference to an existing KMS Alias.
// Experimental.
type AliasAttributes struct {
	// Specifies the alias name.
	//
	// This value must begin with alias/ followed by a name (i.e. alias/ExampleAlias)
	// Experimental.
	AliasName *string `json:"aliasName"`
	// The customer master key (CMK) to which the Alias refers.
	// Experimental.
	AliasTargetKey IKey `json:"aliasTargetKey"`
}

// Construction properties for a KMS Key Alias object.
// Experimental.
type AliasProps struct {
	// The name of the alias.
	//
	// The name must start with alias followed by a
	// forward slash, such as alias/. You can't specify aliases that begin with
	// alias/AWS. These aliases are reserved.
	// Experimental.
	AliasName *string `json:"aliasName"`
	// The ID of the key for which you are creating the alias.
	//
	// Specify the key's
	// globally unique identifier or Amazon Resource Name (ARN). You can't
	// specify another alias.
	// Experimental.
	TargetKey IKey `json:"targetKey"`
	// Policy to apply when the alias is removed from this stack.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
}

// A CloudFormation `AWS::KMS::Alias`.
type CfnAlias interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AliasName() *string
	SetAliasName(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	TargetKeyId() *string
	SetTargetKeyId(val *string)
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

// The jsii proxy struct for CfnAlias
type jsiiProxy_CfnAlias struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAlias) AliasName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"aliasName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) TargetKeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetKeyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::KMS::Alias`.
func NewCfnAlias(scope awscdk.Construct, id *string, props *CfnAliasProps) CfnAlias {
	_init_.Initialize()

	j := jsiiProxy_CfnAlias{}

	_jsii_.Create(
		"monocdk.aws_kms.CfnAlias",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::KMS::Alias`.
func NewCfnAlias_Override(c CfnAlias, scope awscdk.Construct, id *string, props *CfnAliasProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kms.CfnAlias",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAlias) SetAliasName(val *string) {
	_jsii_.Set(
		j,
		"aliasName",
		val,
	)
}

func (j *jsiiProxy_CfnAlias) SetTargetKeyId(val *string) {
	_jsii_.Set(
		j,
		"targetKeyId",
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
func CfnAlias_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.CfnAlias",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAlias_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.CfnAlias",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAlias_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.CfnAlias",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAlias_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_kms.CfnAlias",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAlias) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAlias) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAlias) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAlias) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAlias) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAlias) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAlias) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAlias) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAlias) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAlias) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAlias) OnPrepare() {
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
func (c *jsiiProxy_CfnAlias) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAlias) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAlias) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAlias) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAlias) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAlias) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAlias) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAlias) ToString() *string {
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
func (c *jsiiProxy_CfnAlias) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAlias) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::KMS::Alias`.
type CfnAliasProps struct {
	// `AWS::KMS::Alias.AliasName`.
	AliasName *string `json:"aliasName"`
	// `AWS::KMS::Alias.TargetKeyId`.
	TargetKeyId *string `json:"targetKeyId"`
}

// A CloudFormation `AWS::KMS::Key`.
type CfnKey interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrKeyId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	Enabled() interface{}
	SetEnabled(val interface{})
	EnableKeyRotation() interface{}
	SetEnableKeyRotation(val interface{})
	KeyPolicy() interface{}
	SetKeyPolicy(val interface{})
	KeySpec() *string
	SetKeySpec(val *string)
	KeyUsage() *string
	SetKeyUsage(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	PendingWindowInDays() *float64
	SetPendingWindowInDays(val *float64)
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

// The jsii proxy struct for CfnKey
type jsiiProxy_CfnKey struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnKey) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) AttrKeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrKeyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) Enabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"enabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) EnableKeyRotation() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"enableKeyRotation",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) KeyPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"keyPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) KeySpec() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keySpec",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) KeyUsage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyUsage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) PendingWindowInDays() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"pendingWindowInDays",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnKey) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::KMS::Key`.
func NewCfnKey(scope awscdk.Construct, id *string, props *CfnKeyProps) CfnKey {
	_init_.Initialize()

	j := jsiiProxy_CfnKey{}

	_jsii_.Create(
		"monocdk.aws_kms.CfnKey",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::KMS::Key`.
func NewCfnKey_Override(c CfnKey, scope awscdk.Construct, id *string, props *CfnKeyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kms.CfnKey",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnKey) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnKey) SetEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"enabled",
		val,
	)
}

func (j *jsiiProxy_CfnKey) SetEnableKeyRotation(val interface{}) {
	_jsii_.Set(
		j,
		"enableKeyRotation",
		val,
	)
}

func (j *jsiiProxy_CfnKey) SetKeyPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"keyPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnKey) SetKeySpec(val *string) {
	_jsii_.Set(
		j,
		"keySpec",
		val,
	)
}

func (j *jsiiProxy_CfnKey) SetKeyUsage(val *string) {
	_jsii_.Set(
		j,
		"keyUsage",
		val,
	)
}

func (j *jsiiProxy_CfnKey) SetPendingWindowInDays(val *float64) {
	_jsii_.Set(
		j,
		"pendingWindowInDays",
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
func CfnKey_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.CfnKey",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnKey_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.CfnKey",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnKey_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.CfnKey",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnKey_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_kms.CfnKey",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnKey) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnKey) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnKey) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnKey) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnKey) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnKey) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnKey) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnKey) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnKey) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnKey) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnKey) OnPrepare() {
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
func (c *jsiiProxy_CfnKey) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnKey) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnKey) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnKey) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnKey) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnKey) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnKey) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnKey) ToString() *string {
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
func (c *jsiiProxy_CfnKey) Validate() *[]*string {
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
func (c *jsiiProxy_CfnKey) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::KMS::Key`.
type CfnKeyProps struct {
	// `AWS::KMS::Key.KeyPolicy`.
	KeyPolicy interface{} `json:"keyPolicy"`
	// `AWS::KMS::Key.Description`.
	Description *string `json:"description"`
	// `AWS::KMS::Key.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `AWS::KMS::Key.EnableKeyRotation`.
	EnableKeyRotation interface{} `json:"enableKeyRotation"`
	// `AWS::KMS::Key.KeySpec`.
	KeySpec *string `json:"keySpec"`
	// `AWS::KMS::Key.KeyUsage`.
	KeyUsage *string `json:"keyUsage"`
	// `AWS::KMS::Key.PendingWindowInDays`.
	PendingWindowInDays *float64 `json:"pendingWindowInDays"`
	// `AWS::KMS::Key.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A KMS Key alias.
//
// An alias can be used in all places that expect a key.
// Experimental.
type IAlias interface {
	IKey
	// The name of the alias.
	// Experimental.
	AliasName() *string
	// The Key to which the Alias refers.
	// Experimental.
	AliasTargetKey() IKey
}

// The jsii proxy for IAlias
type jsiiProxy_IAlias struct {
	jsiiProxy_IKey
}

func (j *jsiiProxy_IAlias) AliasName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"aliasName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAlias) AliasTargetKey() IKey {
	var returns IKey
	_jsii_.Get(
		j,
		"aliasTargetKey",
		&returns,
	)
	return returns
}

// A KMS Key, either managed by this CDK app, or imported.
// Experimental.
type IKey interface {
	awscdk.IResource
	// Defines a new alias for the key.
	// Experimental.
	AddAlias(alias *string) Alias
	// Adds a statement to the KMS key resource policy.
	// Experimental.
	AddToResourcePolicy(statement awsiam.PolicyStatement, allowNoOp *bool) *awsiam.AddToResourcePolicyResult
	// Grant the indicated permissions on this key to the given principal.
	// Experimental.
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	// Grant decryption permissions using this key to the given principal.
	// Experimental.
	GrantDecrypt(grantee awsiam.IGrantable) awsiam.Grant
	// Grant encryption permissions using this key to the given principal.
	// Experimental.
	GrantEncrypt(grantee awsiam.IGrantable) awsiam.Grant
	// Grant encryption and decryption permissions using this key to the given principal.
	// Experimental.
	GrantEncryptDecrypt(grantee awsiam.IGrantable) awsiam.Grant
	// The ARN of the key.
	// Experimental.
	KeyArn() *string
	// The ID of the key (the part that looks something like: 1234abcd-12ab-34cd-56ef-1234567890ab).
	// Experimental.
	KeyId() *string
}

// The jsii proxy for IKey
type jsiiProxy_IKey struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IKey) AddAlias(alias *string) Alias {
	var returns Alias

	_jsii_.Invoke(
		i,
		"addAlias",
		[]interface{}{alias},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IKey) AddToResourcePolicy(statement awsiam.PolicyStatement, allowNoOp *bool) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		i,
		"addToResourcePolicy",
		[]interface{}{statement, allowNoOp},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IKey) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
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

func (i *jsiiProxy_IKey) GrantDecrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantDecrypt",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IKey) GrantEncrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantEncrypt",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IKey) GrantEncryptDecrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantEncryptDecrypt",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IKey) KeyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IKey) KeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyId",
		&returns,
	)
	return returns
}

// Defines a KMS key.
// Experimental.
type Key interface {
	awscdk.Resource
	IKey
	Env() *awscdk.ResourceEnvironment
	KeyArn() *string
	KeyId() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Policy() awsiam.PolicyDocument
	Stack() awscdk.Stack
	TrustAccountIdentities() *bool
	AddAlias(aliasName *string) Alias
	AddToResourcePolicy(statement awsiam.PolicyStatement, allowNoOp *bool) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantAdmin(grantee awsiam.IGrantable) awsiam.Grant
	GrantDecrypt(grantee awsiam.IGrantable) awsiam.Grant
	GrantEncrypt(grantee awsiam.IGrantable) awsiam.Grant
	GrantEncryptDecrypt(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Key
type jsiiProxy_Key struct {
	internal.Type__awscdkResource
	jsiiProxy_IKey
}

func (j *jsiiProxy_Key) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Key) KeyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Key) KeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Key) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Key) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Key) Policy() awsiam.PolicyDocument {
	var returns awsiam.PolicyDocument
	_jsii_.Get(
		j,
		"policy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Key) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Key) TrustAccountIdentities() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"trustAccountIdentities",
		&returns,
	)
	return returns
}


// Experimental.
func NewKey(scope constructs.Construct, id *string, props *KeyProps) Key {
	_init_.Initialize()

	j := jsiiProxy_Key{}

	_jsii_.Create(
		"monocdk.aws_kms.Key",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewKey_Override(k Key, scope constructs.Construct, id *string, props *KeyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kms.Key",
		[]interface{}{scope, id, props},
		k,
	)
}

// Create a mutable {@link IKey} based on a low-level {@link CfnKey}.
//
// This is most useful when combined with the cloudformation-include module.
// This method is different than {@link fromKeyArn()} because the {@link IKey}
// returned from this method is mutable;
// meaning, calling any mutating methods on it,
// like {@link IKey.addToResourcePolicy()},
// will actually be reflected in the resulting template,
// as opposed to the object returned from {@link fromKeyArn()},
// on which calling those methods would have no effect.
// Experimental.
func Key_FromCfnKey(cfnKey CfnKey) IKey {
	_init_.Initialize()

	var returns IKey

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Key",
		"fromCfnKey",
		[]interface{}{cfnKey},
		&returns,
	)

	return returns
}

// Import an externally defined KMS Key using its ARN.
// Experimental.
func Key_FromKeyArn(scope constructs.Construct, id *string, keyArn *string) IKey {
	_init_.Initialize()

	var returns IKey

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Key",
		"fromKeyArn",
		[]interface{}{scope, id, keyArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Key_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Key",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Key_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_kms.Key",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Defines a new alias for the key.
// Experimental.
func (k *jsiiProxy_Key) AddAlias(aliasName *string) Alias {
	var returns Alias

	_jsii_.Invoke(
		k,
		"addAlias",
		[]interface{}{aliasName},
		&returns,
	)

	return returns
}

// Adds a statement to the KMS key resource policy.
// Experimental.
func (k *jsiiProxy_Key) AddToResourcePolicy(statement awsiam.PolicyStatement, allowNoOp *bool) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		k,
		"addToResourcePolicy",
		[]interface{}{statement, allowNoOp},
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
func (k *jsiiProxy_Key) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		k,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (k *jsiiProxy_Key) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		k,
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
func (k *jsiiProxy_Key) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		k,
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
func (k *jsiiProxy_Key) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		k,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the indicated permissions on this key to the given principal.
//
// This modifies both the principal's policy as well as the resource policy,
// since the default CloudFormation setup for KMS keys is that the policy
// must not be empty and so default grants won't work.
// Experimental.
func (k *jsiiProxy_Key) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		k,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Grant admins permissions using this key to the given principal.
//
// Key administrators have permissions to manage the key (e.g., change permissions, revoke), but do not have permissions
// to use the key in cryptographic operations (e.g., encrypt, decrypt).
// Experimental.
func (k *jsiiProxy_Key) GrantAdmin(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		k,
		"grantAdmin",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant decryption permissions using this key to the given principal.
// Experimental.
func (k *jsiiProxy_Key) GrantDecrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		k,
		"grantDecrypt",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant encryption permissions using this key to the given principal.
// Experimental.
func (k *jsiiProxy_Key) GrantEncrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		k,
		"grantEncrypt",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant encryption and decryption permissions using this key to the given principal.
// Experimental.
func (k *jsiiProxy_Key) GrantEncryptDecrypt(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		k,
		"grantEncryptDecrypt",
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
func (k *jsiiProxy_Key) OnPrepare() {
	_jsii_.InvokeVoid(
		k,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (k *jsiiProxy_Key) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		k,
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
func (k *jsiiProxy_Key) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		k,
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
func (k *jsiiProxy_Key) Prepare() {
	_jsii_.InvokeVoid(
		k,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (k *jsiiProxy_Key) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		k,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (k *jsiiProxy_Key) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		k,
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
func (k *jsiiProxy_Key) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		k,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for a KMS Key object.
// Experimental.
type KeyProps struct {
	// A list of principals to add as key administrators to the key policy.
	//
	// Key administrators have permissions to manage the key (e.g., change permissions, revoke), but do not have permissions
	// to use the key in cryptographic operations (e.g., encrypt, decrypt).
	//
	// These principals will be added to the default key policy (if none specified), or to the specified policy (if provided).
	// Experimental.
	Admins *[]awsiam.IPrincipal `json:"admins"`
	// Initial alias to add to the key.
	//
	// More aliases can be added later by calling `addAlias`.
	// Experimental.
	Alias *string `json:"alias"`
	// A description of the key.
	//
	// Use a description that helps your users decide
	// whether the key is appropriate for a particular task.
	// Experimental.
	Description *string `json:"description"`
	// Indicates whether the key is available for use.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// Indicates whether AWS KMS rotates the key.
	// Experimental.
	EnableKeyRotation *bool `json:"enableKeyRotation"`
	// The cryptographic configuration of the key. The valid value depends on usage of the key.
	//
	// IMPORTANT: If you change this property of an existing key, the existing key is scheduled for deletion
	// and a new key is created with the specified value.
	// Experimental.
	KeySpec KeySpec `json:"keySpec"`
	// The cryptographic operations for which the key can be used.
	//
	// IMPORTANT: If you change this property of an existing key, the existing key is scheduled for deletion
	// and a new key is created with the specified value.
	// Experimental.
	KeyUsage KeyUsage `json:"keyUsage"`
	// Specifies the number of days in the waiting period before AWS KMS deletes a CMK that has been removed from a CloudFormation stack.
	//
	// When you remove a customer master key (CMK) from a CloudFormation stack, AWS KMS schedules the CMK for deletion
	// and starts the mandatory waiting period. The PendingWindowInDays property determines the length of waiting period.
	// During the waiting period, the key state of CMK is Pending Deletion, which prevents the CMK from being used in
	// cryptographic operations. When the waiting period expires, AWS KMS permanently deletes the CMK.
	//
	// Enter a value between 7 and 30 days.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kms-key.html#cfn-kms-key-pendingwindowindays
	//
	// Experimental.
	PendingWindow awscdk.Duration `json:"pendingWindow"`
	// Custom policy document to attach to the KMS key.
	//
	// NOTE - If the `@aws-cdk/aws-kms:defaultKeyPolicies` feature flag is set (the default for new projects),
	// this policy will *override* the default key policy and become the only key policy for the key. If the
	// feature flag is not set, this policy will be appended to the default key policy.
	// Experimental.
	Policy awsiam.PolicyDocument `json:"policy"`
	// Whether the encryption key should be retained when it is removed from the Stack.
	//
	// This is useful when one wants to
	// retain access to data that was encrypted with a key that is being retired.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// Whether the key usage can be granted by IAM policies.
	//
	// Setting this to true adds a default statement which delegates key
	// access control completely to the identity's IAM policy (similar
	// to how it works for other AWS resources). This matches the default behavior
	// when creating KMS keys via the API or console.
	//
	// If the `@aws-cdk/aws-kms:defaultKeyPolicies` feature flag is set (the default for new projects),
	// this flag will always be treated as 'true' and does not need to be explicitly set.
	// See: https://docs.aws.amazon.com/kms/latest/developerguide/key-policies.html#key-policy-default-allow-root-enable-iam
	//
	// Deprecated: redundant with the `@aws-cdk/aws-kms:defaultKeyPolicies` feature flag
	TrustAccountIdentities *bool `json:"trustAccountIdentities"`
}

// The key spec, represents the cryptographic configuration of keys.
// Experimental.
type KeySpec string

const (
	KeySpec_SYMMETRIC_DEFAULT KeySpec = "SYMMETRIC_DEFAULT"
	KeySpec_RSA_2048 KeySpec = "RSA_2048"
	KeySpec_RSA_3072 KeySpec = "RSA_3072"
	KeySpec_RSA_4096 KeySpec = "RSA_4096"
	KeySpec_ECC_NIST_P256 KeySpec = "ECC_NIST_P256"
	KeySpec_ECC_NIST_P384 KeySpec = "ECC_NIST_P384"
	KeySpec_ECC_NIST_P521 KeySpec = "ECC_NIST_P521"
	KeySpec_ECC_SECG_P256K1 KeySpec = "ECC_SECG_P256K1"
)

// The key usage, represents the cryptographic operations of keys.
// Experimental.
type KeyUsage string

const (
	KeyUsage_ENCRYPT_DECRYPT KeyUsage = "ENCRYPT_DECRYPT"
	KeyUsage_SIGN_VERIFY KeyUsage = "SIGN_VERIFY"
)

// A principal to allow access to a key if it's being used through another AWS service.
// Experimental.
type ViaServicePrincipal interface {
	awsiam.PrincipalBase
	AssumeRoleAction() *string
	GrantPrincipal() awsiam.IPrincipal
	PolicyFragment() awsiam.PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement awsiam.PolicyStatement) *bool
	AddToPrincipalPolicy(_statement awsiam.PolicyStatement) *awsiam.AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) awsiam.IPrincipal
}

// The jsii proxy struct for ViaServicePrincipal
type jsiiProxy_ViaServicePrincipal struct {
	internal.Type__awsiamPrincipalBase
}

func (j *jsiiProxy_ViaServicePrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ViaServicePrincipal) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ViaServicePrincipal) PolicyFragment() awsiam.PrincipalPolicyFragment {
	var returns awsiam.PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ViaServicePrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewViaServicePrincipal(serviceName *string, basePrincipal awsiam.IPrincipal) ViaServicePrincipal {
	_init_.Initialize()

	j := jsiiProxy_ViaServicePrincipal{}

	_jsii_.Create(
		"monocdk.aws_kms.ViaServicePrincipal",
		[]interface{}{serviceName, basePrincipal},
		&j,
	)

	return &j
}

// Experimental.
func NewViaServicePrincipal_Override(v ViaServicePrincipal, serviceName *string, basePrincipal awsiam.IPrincipal) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_kms.ViaServicePrincipal",
		[]interface{}{serviceName, basePrincipal},
		v,
	)
}

// Add to the policy of this principal.
// Experimental.
func (v *jsiiProxy_ViaServicePrincipal) AddToPolicy(statement awsiam.PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		v,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (v *jsiiProxy_ViaServicePrincipal) AddToPrincipalPolicy(_statement awsiam.PolicyStatement) *awsiam.AddToPrincipalPolicyResult {
	var returns *awsiam.AddToPrincipalPolicyResult

	_jsii_.Invoke(
		v,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (v *jsiiProxy_ViaServicePrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		v,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (v *jsiiProxy_ViaServicePrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		v,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (v *jsiiProxy_ViaServicePrincipal) WithConditions(conditions *map[string]interface{}) awsiam.IPrincipal {
	var returns awsiam.IPrincipal

	_jsii_.Invoke(
		v,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}


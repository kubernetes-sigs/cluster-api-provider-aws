package awsevents

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsevents/internal"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
)

// Define an EventBridge Archive.
// Experimental.
type Archive interface {
	awscdk.Resource
	ArchiveArn() *string
	ArchiveName() *string
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

// The jsii proxy struct for Archive
type jsiiProxy_Archive struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_Archive) ArchiveArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"archiveArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Archive) ArchiveName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"archiveName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Archive) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Archive) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Archive) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Archive) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewArchive(scope constructs.Construct, id *string, props *ArchiveProps) Archive {
	_init_.Initialize()

	j := jsiiProxy_Archive{}

	_jsii_.Create(
		"monocdk.aws_events.Archive",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewArchive_Override(a Archive, scope constructs.Construct, id *string, props *ArchiveProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.Archive",
		[]interface{}{scope, id, props},
		a,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Archive_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.Archive",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Archive_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.Archive",
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
func (a *jsiiProxy_Archive) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_Archive) GeneratePhysicalName() *string {
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
func (a *jsiiProxy_Archive) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (a *jsiiProxy_Archive) GetResourceNameAttribute(nameAttr *string) *string {
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
func (a *jsiiProxy_Archive) OnPrepare() {
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
func (a *jsiiProxy_Archive) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_Archive) OnValidate() *[]*string {
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
func (a *jsiiProxy_Archive) Prepare() {
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
func (a *jsiiProxy_Archive) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_Archive) ToString() *string {
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
func (a *jsiiProxy_Archive) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The event archive properties.
// Experimental.
type ArchiveProps struct {
	// An event pattern to use to filter events sent to the archive.
	// Experimental.
	EventPattern *EventPattern `json:"eventPattern"`
	// The name of the archive.
	// Experimental.
	ArchiveName *string `json:"archiveName"`
	// A description for the archive.
	// Experimental.
	Description *string `json:"description"`
	// The number of days to retain events for.
	//
	// Default value is 0. If set to 0, events are retained indefinitely.
	// Experimental.
	Retention awscdk.Duration `json:"retention"`
	// The event source associated with the archive.
	// Experimental.
	SourceEventBus IEventBus `json:"sourceEventBus"`
}

// The event archive base properties.
// Experimental.
type BaseArchiveProps struct {
	// An event pattern to use to filter events sent to the archive.
	// Experimental.
	EventPattern *EventPattern `json:"eventPattern"`
	// The name of the archive.
	// Experimental.
	ArchiveName *string `json:"archiveName"`
	// A description for the archive.
	// Experimental.
	Description *string `json:"description"`
	// The number of days to retain events for.
	//
	// Default value is 0. If set to 0, events are retained indefinitely.
	// Experimental.
	Retention awscdk.Duration `json:"retention"`
}

// A CloudFormation `AWS::Events::ApiDestination`.
type CfnApiDestination interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ConnectionArn() *string
	SetConnectionArn(val *string)
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	HttpMethod() *string
	SetHttpMethod(val *string)
	InvocationEndpoint() *string
	SetInvocationEndpoint(val *string)
	InvocationRateLimitPerSecond() *float64
	SetInvocationRateLimitPerSecond(val *float64)
	LogicalId() *string
	Name() *string
	SetName(val *string)
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

// The jsii proxy struct for CfnApiDestination
type jsiiProxy_CfnApiDestination struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnApiDestination) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) ConnectionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"connectionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) HttpMethod() *string {
	var returns *string
	_jsii_.Get(
		j,
		"httpMethod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) InvocationEndpoint() *string {
	var returns *string
	_jsii_.Get(
		j,
		"invocationEndpoint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) InvocationRateLimitPerSecond() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"invocationRateLimitPerSecond",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiDestination) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Events::ApiDestination`.
func NewCfnApiDestination(scope awscdk.Construct, id *string, props *CfnApiDestinationProps) CfnApiDestination {
	_init_.Initialize()

	j := jsiiProxy_CfnApiDestination{}

	_jsii_.Create(
		"monocdk.aws_events.CfnApiDestination",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Events::ApiDestination`.
func NewCfnApiDestination_Override(c CfnApiDestination, scope awscdk.Construct, id *string, props *CfnApiDestinationProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.CfnApiDestination",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnApiDestination) SetConnectionArn(val *string) {
	_jsii_.Set(
		j,
		"connectionArn",
		val,
	)
}

func (j *jsiiProxy_CfnApiDestination) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnApiDestination) SetHttpMethod(val *string) {
	_jsii_.Set(
		j,
		"httpMethod",
		val,
	)
}

func (j *jsiiProxy_CfnApiDestination) SetInvocationEndpoint(val *string) {
	_jsii_.Set(
		j,
		"invocationEndpoint",
		val,
	)
}

func (j *jsiiProxy_CfnApiDestination) SetInvocationRateLimitPerSecond(val *float64) {
	_jsii_.Set(
		j,
		"invocationRateLimitPerSecond",
		val,
	)
}

func (j *jsiiProxy_CfnApiDestination) SetName(val *string) {
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
func CfnApiDestination_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnApiDestination",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnApiDestination_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnApiDestination",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnApiDestination_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnApiDestination",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnApiDestination_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.CfnApiDestination",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnApiDestination) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnApiDestination) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnApiDestination) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnApiDestination) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnApiDestination) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnApiDestination) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnApiDestination) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnApiDestination) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnApiDestination) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnApiDestination) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnApiDestination) OnPrepare() {
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
func (c *jsiiProxy_CfnApiDestination) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnApiDestination) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnApiDestination) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnApiDestination) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnApiDestination) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnApiDestination) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnApiDestination) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnApiDestination) ToString() *string {
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
func (c *jsiiProxy_CfnApiDestination) Validate() *[]*string {
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
func (c *jsiiProxy_CfnApiDestination) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Events::ApiDestination`.
type CfnApiDestinationProps struct {
	// `AWS::Events::ApiDestination.ConnectionArn`.
	ConnectionArn *string `json:"connectionArn"`
	// `AWS::Events::ApiDestination.HttpMethod`.
	HttpMethod *string `json:"httpMethod"`
	// `AWS::Events::ApiDestination.InvocationEndpoint`.
	InvocationEndpoint *string `json:"invocationEndpoint"`
	// `AWS::Events::ApiDestination.Description`.
	Description *string `json:"description"`
	// `AWS::Events::ApiDestination.InvocationRateLimitPerSecond`.
	InvocationRateLimitPerSecond *float64 `json:"invocationRateLimitPerSecond"`
	// `AWS::Events::ApiDestination.Name`.
	Name *string `json:"name"`
}

// A CloudFormation `AWS::Events::Archive`.
type CfnArchive interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ArchiveName() *string
	SetArchiveName(val *string)
	AttrArchiveName() *string
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	EventPattern() interface{}
	SetEventPattern(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RetentionDays() *float64
	SetRetentionDays(val *float64)
	SourceArn() *string
	SetSourceArn(val *string)
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

// The jsii proxy struct for CfnArchive
type jsiiProxy_CfnArchive struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnArchive) ArchiveName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"archiveName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) AttrArchiveName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArchiveName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) EventPattern() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"eventPattern",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) RetentionDays() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"retentionDays",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) SourceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnArchive) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Events::Archive`.
func NewCfnArchive(scope awscdk.Construct, id *string, props *CfnArchiveProps) CfnArchive {
	_init_.Initialize()

	j := jsiiProxy_CfnArchive{}

	_jsii_.Create(
		"monocdk.aws_events.CfnArchive",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Events::Archive`.
func NewCfnArchive_Override(c CfnArchive, scope awscdk.Construct, id *string, props *CfnArchiveProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.CfnArchive",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnArchive) SetArchiveName(val *string) {
	_jsii_.Set(
		j,
		"archiveName",
		val,
	)
}

func (j *jsiiProxy_CfnArchive) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnArchive) SetEventPattern(val interface{}) {
	_jsii_.Set(
		j,
		"eventPattern",
		val,
	)
}

func (j *jsiiProxy_CfnArchive) SetRetentionDays(val *float64) {
	_jsii_.Set(
		j,
		"retentionDays",
		val,
	)
}

func (j *jsiiProxy_CfnArchive) SetSourceArn(val *string) {
	_jsii_.Set(
		j,
		"sourceArn",
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
func CfnArchive_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnArchive",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnArchive_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnArchive",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnArchive_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnArchive",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnArchive_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.CfnArchive",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnArchive) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnArchive) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnArchive) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnArchive) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnArchive) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnArchive) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnArchive) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnArchive) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnArchive) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnArchive) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnArchive) OnPrepare() {
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
func (c *jsiiProxy_CfnArchive) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnArchive) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnArchive) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnArchive) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnArchive) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnArchive) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnArchive) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnArchive) ToString() *string {
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
func (c *jsiiProxy_CfnArchive) Validate() *[]*string {
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
func (c *jsiiProxy_CfnArchive) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Events::Archive`.
type CfnArchiveProps struct {
	// `AWS::Events::Archive.SourceArn`.
	SourceArn *string `json:"sourceArn"`
	// `AWS::Events::Archive.ArchiveName`.
	ArchiveName *string `json:"archiveName"`
	// `AWS::Events::Archive.Description`.
	Description *string `json:"description"`
	// `AWS::Events::Archive.EventPattern`.
	EventPattern interface{} `json:"eventPattern"`
	// `AWS::Events::Archive.RetentionDays`.
	RetentionDays *float64 `json:"retentionDays"`
}

// A CloudFormation `AWS::Events::Connection`.
type CfnConnection interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrSecretArn() *string
	AuthorizationType() *string
	SetAuthorizationType(val *string)
	AuthParameters() interface{}
	SetAuthParameters(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
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

// The jsii proxy struct for CfnConnection
type jsiiProxy_CfnConnection struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnConnection) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) AttrSecretArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrSecretArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) AuthorizationType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) AuthParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"authParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnConnection) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Events::Connection`.
func NewCfnConnection(scope awscdk.Construct, id *string, props *CfnConnectionProps) CfnConnection {
	_init_.Initialize()

	j := jsiiProxy_CfnConnection{}

	_jsii_.Create(
		"monocdk.aws_events.CfnConnection",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Events::Connection`.
func NewCfnConnection_Override(c CfnConnection, scope awscdk.Construct, id *string, props *CfnConnectionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.CfnConnection",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnConnection) SetAuthorizationType(val *string) {
	_jsii_.Set(
		j,
		"authorizationType",
		val,
	)
}

func (j *jsiiProxy_CfnConnection) SetAuthParameters(val interface{}) {
	_jsii_.Set(
		j,
		"authParameters",
		val,
	)
}

func (j *jsiiProxy_CfnConnection) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnConnection) SetName(val *string) {
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
func CfnConnection_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnConnection",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnConnection_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnConnection",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnConnection_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnConnection",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnConnection_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.CfnConnection",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnConnection) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnConnection) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnConnection) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnConnection) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnConnection) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnConnection) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnConnection) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnConnection) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnConnection) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnConnection) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnConnection) OnPrepare() {
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
func (c *jsiiProxy_CfnConnection) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnConnection) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnConnection) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnConnection) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnConnection) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnConnection) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnConnection) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnConnection) ToString() *string {
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
func (c *jsiiProxy_CfnConnection) Validate() *[]*string {
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
func (c *jsiiProxy_CfnConnection) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Events::Connection`.
type CfnConnectionProps struct {
	// `AWS::Events::Connection.AuthorizationType`.
	AuthorizationType *string `json:"authorizationType"`
	// `AWS::Events::Connection.AuthParameters`.
	AuthParameters interface{} `json:"authParameters"`
	// `AWS::Events::Connection.Description`.
	Description *string `json:"description"`
	// `AWS::Events::Connection.Name`.
	Name *string `json:"name"`
}

// A CloudFormation `AWS::Events::EventBus`.
type CfnEventBus interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrName() *string
	AttrPolicy() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	EventSourceName() *string
	SetEventSourceName(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
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

// The jsii proxy struct for CfnEventBus
type jsiiProxy_CfnEventBus struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnEventBus) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) AttrName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) AttrPolicy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) EventSourceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventSourceName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBus) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Events::EventBus`.
func NewCfnEventBus(scope awscdk.Construct, id *string, props *CfnEventBusProps) CfnEventBus {
	_init_.Initialize()

	j := jsiiProxy_CfnEventBus{}

	_jsii_.Create(
		"monocdk.aws_events.CfnEventBus",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Events::EventBus`.
func NewCfnEventBus_Override(c CfnEventBus, scope awscdk.Construct, id *string, props *CfnEventBusProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.CfnEventBus",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnEventBus) SetEventSourceName(val *string) {
	_jsii_.Set(
		j,
		"eventSourceName",
		val,
	)
}

func (j *jsiiProxy_CfnEventBus) SetName(val *string) {
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
func CfnEventBus_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnEventBus",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnEventBus_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnEventBus",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnEventBus_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnEventBus",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnEventBus_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.CfnEventBus",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnEventBus) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnEventBus) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnEventBus) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnEventBus) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnEventBus) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnEventBus) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnEventBus) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnEventBus) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnEventBus) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnEventBus) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnEventBus) OnPrepare() {
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
func (c *jsiiProxy_CfnEventBus) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventBus) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnEventBus) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnEventBus) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnEventBus) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnEventBus) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnEventBus) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventBus) ToString() *string {
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
func (c *jsiiProxy_CfnEventBus) Validate() *[]*string {
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
func (c *jsiiProxy_CfnEventBus) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// A CloudFormation `AWS::Events::EventBusPolicy`.
type CfnEventBusPolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Action() *string
	SetAction(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Condition() interface{}
	SetCondition(val interface{})
	CreationStack() *[]*string
	EventBusName() *string
	SetEventBusName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Principal() *string
	SetPrincipal(val *string)
	Ref() *string
	Stack() awscdk.Stack
	Statement() interface{}
	SetStatement(val interface{})
	StatementId() *string
	SetStatementId(val *string)
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

// The jsii proxy struct for CfnEventBusPolicy
type jsiiProxy_CfnEventBusPolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnEventBusPolicy) Action() *string {
	var returns *string
	_jsii_.Get(
		j,
		"action",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) Condition() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"condition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) EventBusName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) Principal() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) Statement() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"statement",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) StatementId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"statementId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventBusPolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Events::EventBusPolicy`.
func NewCfnEventBusPolicy(scope awscdk.Construct, id *string, props *CfnEventBusPolicyProps) CfnEventBusPolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnEventBusPolicy{}

	_jsii_.Create(
		"monocdk.aws_events.CfnEventBusPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Events::EventBusPolicy`.
func NewCfnEventBusPolicy_Override(c CfnEventBusPolicy, scope awscdk.Construct, id *string, props *CfnEventBusPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.CfnEventBusPolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnEventBusPolicy) SetAction(val *string) {
	_jsii_.Set(
		j,
		"action",
		val,
	)
}

func (j *jsiiProxy_CfnEventBusPolicy) SetCondition(val interface{}) {
	_jsii_.Set(
		j,
		"condition",
		val,
	)
}

func (j *jsiiProxy_CfnEventBusPolicy) SetEventBusName(val *string) {
	_jsii_.Set(
		j,
		"eventBusName",
		val,
	)
}

func (j *jsiiProxy_CfnEventBusPolicy) SetPrincipal(val *string) {
	_jsii_.Set(
		j,
		"principal",
		val,
	)
}

func (j *jsiiProxy_CfnEventBusPolicy) SetStatement(val interface{}) {
	_jsii_.Set(
		j,
		"statement",
		val,
	)
}

func (j *jsiiProxy_CfnEventBusPolicy) SetStatementId(val *string) {
	_jsii_.Set(
		j,
		"statementId",
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
func CfnEventBusPolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnEventBusPolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnEventBusPolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnEventBusPolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnEventBusPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnEventBusPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnEventBusPolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.CfnEventBusPolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnEventBusPolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnEventBusPolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnEventBusPolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnEventBusPolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnEventBusPolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnEventBusPolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnEventBusPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnEventBusPolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnEventBusPolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnEventBusPolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnEventBusPolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnEventBusPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventBusPolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnEventBusPolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnEventBusPolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnEventBusPolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnEventBusPolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnEventBusPolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventBusPolicy) ToString() *string {
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
func (c *jsiiProxy_CfnEventBusPolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnEventBusPolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnEventBusPolicy_ConditionProperty struct {
	// `CfnEventBusPolicy.ConditionProperty.Key`.
	Key *string `json:"key"`
	// `CfnEventBusPolicy.ConditionProperty.Type`.
	Type *string `json:"type"`
	// `CfnEventBusPolicy.ConditionProperty.Value`.
	Value *string `json:"value"`
}

// Properties for defining a `AWS::Events::EventBusPolicy`.
type CfnEventBusPolicyProps struct {
	// `AWS::Events::EventBusPolicy.StatementId`.
	StatementId *string `json:"statementId"`
	// `AWS::Events::EventBusPolicy.Action`.
	Action *string `json:"action"`
	// `AWS::Events::EventBusPolicy.Condition`.
	Condition interface{} `json:"condition"`
	// `AWS::Events::EventBusPolicy.EventBusName`.
	EventBusName *string `json:"eventBusName"`
	// `AWS::Events::EventBusPolicy.Principal`.
	Principal *string `json:"principal"`
	// `AWS::Events::EventBusPolicy.Statement`.
	Statement interface{} `json:"statement"`
}

// Properties for defining a `AWS::Events::EventBus`.
type CfnEventBusProps struct {
	// `AWS::Events::EventBus.Name`.
	Name *string `json:"name"`
	// `AWS::Events::EventBus.EventSourceName`.
	EventSourceName *string `json:"eventSourceName"`
}

// A CloudFormation `AWS::Events::Rule`.
type CfnRule interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	EventBusName() *string
	SetEventBusName(val *string)
	EventPattern() interface{}
	SetEventPattern(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	RoleArn() *string
	SetRoleArn(val *string)
	ScheduleExpression() *string
	SetScheduleExpression(val *string)
	Stack() awscdk.Stack
	State() *string
	SetState(val *string)
	Targets() interface{}
	SetTargets(val interface{})
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

// The jsii proxy struct for CfnRule
type jsiiProxy_CfnRule struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRule) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) EventBusName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) EventPattern() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"eventPattern",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) ScheduleExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scheduleExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) State() *string {
	var returns *string
	_jsii_.Get(
		j,
		"state",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) Targets() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"targets",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Events::Rule`.
func NewCfnRule(scope awscdk.Construct, id *string, props *CfnRuleProps) CfnRule {
	_init_.Initialize()

	j := jsiiProxy_CfnRule{}

	_jsii_.Create(
		"monocdk.aws_events.CfnRule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Events::Rule`.
func NewCfnRule_Override(c CfnRule, scope awscdk.Construct, id *string, props *CfnRuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.CfnRule",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRule) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnRule) SetEventBusName(val *string) {
	_jsii_.Set(
		j,
		"eventBusName",
		val,
	)
}

func (j *jsiiProxy_CfnRule) SetEventPattern(val interface{}) {
	_jsii_.Set(
		j,
		"eventPattern",
		val,
	)
}

func (j *jsiiProxy_CfnRule) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnRule) SetRoleArn(val *string) {
	_jsii_.Set(
		j,
		"roleArn",
		val,
	)
}

func (j *jsiiProxy_CfnRule) SetScheduleExpression(val *string) {
	_jsii_.Set(
		j,
		"scheduleExpression",
		val,
	)
}

func (j *jsiiProxy_CfnRule) SetState(val *string) {
	_jsii_.Set(
		j,
		"state",
		val,
	)
}

func (j *jsiiProxy_CfnRule) SetTargets(val interface{}) {
	_jsii_.Set(
		j,
		"targets",
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
func CfnRule_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnRule",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnRule_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnRule",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnRule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.CfnRule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRule_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.CfnRule",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnRule) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnRule) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnRule) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnRule) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnRule) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnRule) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnRule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnRule) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnRule) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnRule) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnRule) OnPrepare() {
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
func (c *jsiiProxy_CfnRule) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRule) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnRule) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnRule) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnRule) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnRule) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnRule) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRule) ToString() *string {
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
func (c *jsiiProxy_CfnRule) Validate() *[]*string {
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
func (c *jsiiProxy_CfnRule) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnRule_AwsVpcConfigurationProperty struct {
	// `CfnRule.AwsVpcConfigurationProperty.Subnets`.
	Subnets *[]*string `json:"subnets"`
	// `CfnRule.AwsVpcConfigurationProperty.AssignPublicIp`.
	AssignPublicIp *string `json:"assignPublicIp"`
	// `CfnRule.AwsVpcConfigurationProperty.SecurityGroups`.
	SecurityGroups *[]*string `json:"securityGroups"`
}

type CfnRule_BatchArrayPropertiesProperty struct {
	// `CfnRule.BatchArrayPropertiesProperty.Size`.
	Size *float64 `json:"size"`
}

type CfnRule_BatchParametersProperty struct {
	// `CfnRule.BatchParametersProperty.JobDefinition`.
	JobDefinition *string `json:"jobDefinition"`
	// `CfnRule.BatchParametersProperty.JobName`.
	JobName *string `json:"jobName"`
	// `CfnRule.BatchParametersProperty.ArrayProperties`.
	ArrayProperties interface{} `json:"arrayProperties"`
	// `CfnRule.BatchParametersProperty.RetryStrategy`.
	RetryStrategy interface{} `json:"retryStrategy"`
}

type CfnRule_BatchRetryStrategyProperty struct {
	// `CfnRule.BatchRetryStrategyProperty.Attempts`.
	Attempts *float64 `json:"attempts"`
}

type CfnRule_DeadLetterConfigProperty struct {
	// `CfnRule.DeadLetterConfigProperty.Arn`.
	Arn *string `json:"arn"`
}

type CfnRule_EcsParametersProperty struct {
	// `CfnRule.EcsParametersProperty.TaskDefinitionArn`.
	TaskDefinitionArn *string `json:"taskDefinitionArn"`
	// `CfnRule.EcsParametersProperty.Group`.
	Group *string `json:"group"`
	// `CfnRule.EcsParametersProperty.LaunchType`.
	LaunchType *string `json:"launchType"`
	// `CfnRule.EcsParametersProperty.NetworkConfiguration`.
	NetworkConfiguration interface{} `json:"networkConfiguration"`
	// `CfnRule.EcsParametersProperty.PlatformVersion`.
	PlatformVersion *string `json:"platformVersion"`
	// `CfnRule.EcsParametersProperty.TaskCount`.
	TaskCount *float64 `json:"taskCount"`
}

type CfnRule_HttpParametersProperty struct {
	// `CfnRule.HttpParametersProperty.HeaderParameters`.
	HeaderParameters interface{} `json:"headerParameters"`
	// `CfnRule.HttpParametersProperty.PathParameterValues`.
	PathParameterValues *[]*string `json:"pathParameterValues"`
	// `CfnRule.HttpParametersProperty.QueryStringParameters`.
	QueryStringParameters interface{} `json:"queryStringParameters"`
}

type CfnRule_InputTransformerProperty struct {
	// `CfnRule.InputTransformerProperty.InputTemplate`.
	InputTemplate *string `json:"inputTemplate"`
	// `CfnRule.InputTransformerProperty.InputPathsMap`.
	InputPathsMap interface{} `json:"inputPathsMap"`
}

type CfnRule_KinesisParametersProperty struct {
	// `CfnRule.KinesisParametersProperty.PartitionKeyPath`.
	PartitionKeyPath *string `json:"partitionKeyPath"`
}

type CfnRule_NetworkConfigurationProperty struct {
	// `CfnRule.NetworkConfigurationProperty.AwsVpcConfiguration`.
	AwsVpcConfiguration interface{} `json:"awsVpcConfiguration"`
}

type CfnRule_RedshiftDataParametersProperty struct {
	// `CfnRule.RedshiftDataParametersProperty.Database`.
	Database *string `json:"database"`
	// `CfnRule.RedshiftDataParametersProperty.Sql`.
	Sql *string `json:"sql"`
	// `CfnRule.RedshiftDataParametersProperty.DbUser`.
	DbUser *string `json:"dbUser"`
	// `CfnRule.RedshiftDataParametersProperty.SecretManagerArn`.
	SecretManagerArn *string `json:"secretManagerArn"`
	// `CfnRule.RedshiftDataParametersProperty.StatementName`.
	StatementName *string `json:"statementName"`
	// `CfnRule.RedshiftDataParametersProperty.WithEvent`.
	WithEvent interface{} `json:"withEvent"`
}

type CfnRule_RetryPolicyProperty struct {
	// `CfnRule.RetryPolicyProperty.MaximumEventAgeInSeconds`.
	MaximumEventAgeInSeconds *float64 `json:"maximumEventAgeInSeconds"`
	// `CfnRule.RetryPolicyProperty.MaximumRetryAttempts`.
	MaximumRetryAttempts *float64 `json:"maximumRetryAttempts"`
}

type CfnRule_RunCommandParametersProperty struct {
	// `CfnRule.RunCommandParametersProperty.RunCommandTargets`.
	RunCommandTargets interface{} `json:"runCommandTargets"`
}

type CfnRule_RunCommandTargetProperty struct {
	// `CfnRule.RunCommandTargetProperty.Key`.
	Key *string `json:"key"`
	// `CfnRule.RunCommandTargetProperty.Values`.
	Values *[]*string `json:"values"`
}

type CfnRule_SqsParametersProperty struct {
	// `CfnRule.SqsParametersProperty.MessageGroupId`.
	MessageGroupId *string `json:"messageGroupId"`
}

type CfnRule_TargetProperty struct {
	// `CfnRule.TargetProperty.Arn`.
	Arn *string `json:"arn"`
	// `CfnRule.TargetProperty.Id`.
	Id *string `json:"id"`
	// `CfnRule.TargetProperty.BatchParameters`.
	BatchParameters interface{} `json:"batchParameters"`
	// `CfnRule.TargetProperty.DeadLetterConfig`.
	DeadLetterConfig interface{} `json:"deadLetterConfig"`
	// `CfnRule.TargetProperty.EcsParameters`.
	EcsParameters interface{} `json:"ecsParameters"`
	// `CfnRule.TargetProperty.HttpParameters`.
	HttpParameters interface{} `json:"httpParameters"`
	// `CfnRule.TargetProperty.Input`.
	Input *string `json:"input"`
	// `CfnRule.TargetProperty.InputPath`.
	InputPath *string `json:"inputPath"`
	// `CfnRule.TargetProperty.InputTransformer`.
	InputTransformer interface{} `json:"inputTransformer"`
	// `CfnRule.TargetProperty.KinesisParameters`.
	KinesisParameters interface{} `json:"kinesisParameters"`
	// `CfnRule.TargetProperty.RedshiftDataParameters`.
	RedshiftDataParameters interface{} `json:"redshiftDataParameters"`
	// `CfnRule.TargetProperty.RetryPolicy`.
	RetryPolicy interface{} `json:"retryPolicy"`
	// `CfnRule.TargetProperty.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `CfnRule.TargetProperty.RunCommandParameters`.
	RunCommandParameters interface{} `json:"runCommandParameters"`
	// `CfnRule.TargetProperty.SqsParameters`.
	SqsParameters interface{} `json:"sqsParameters"`
}

// Properties for defining a `AWS::Events::Rule`.
type CfnRuleProps struct {
	// `AWS::Events::Rule.Description`.
	Description *string `json:"description"`
	// `AWS::Events::Rule.EventBusName`.
	EventBusName *string `json:"eventBusName"`
	// `AWS::Events::Rule.EventPattern`.
	EventPattern interface{} `json:"eventPattern"`
	// `AWS::Events::Rule.Name`.
	Name *string `json:"name"`
	// `AWS::Events::Rule.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `AWS::Events::Rule.ScheduleExpression`.
	ScheduleExpression *string `json:"scheduleExpression"`
	// `AWS::Events::Rule.State`.
	State *string `json:"state"`
	// `AWS::Events::Rule.Targets`.
	Targets interface{} `json:"targets"`
}

// Options to configure a cron expression.
//
// All fields are strings so you can use complex expressions. Absence of
// a field implies '*' or '?', whichever one is appropriate.
// See: https://docs.aws.amazon.com/eventbridge/latest/userguide/scheduled-events.html#cron-expressions
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

// Define an EventBridge EventBus.
// Experimental.
type EventBus interface {
	awscdk.Resource
	IEventBus
	Env() *awscdk.ResourceEnvironment
	EventBusArn() *string
	EventBusName() *string
	EventBusPolicy() *string
	EventSourceName() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	Archive(id *string, props *BaseArchiveProps) Archive
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantPutEventsTo(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for EventBus
type jsiiProxy_EventBus struct {
	internal.Type__awscdkResource
	jsiiProxy_IEventBus
}

func (j *jsiiProxy_EventBus) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventBus) EventBusArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventBus) EventBusName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventBus) EventBusPolicy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventBus) EventSourceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventSourceName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventBus) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventBus) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventBus) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewEventBus(scope constructs.Construct, id *string, props *EventBusProps) EventBus {
	_init_.Initialize()

	j := jsiiProxy_EventBus{}

	_jsii_.Create(
		"monocdk.aws_events.EventBus",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewEventBus_Override(e EventBus, scope constructs.Construct, id *string, props *EventBusProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.EventBus",
		[]interface{}{scope, id, props},
		e,
	)
}

// Import an existing event bus resource.
// Experimental.
func EventBus_FromEventBusArn(scope constructs.Construct, id *string, eventBusArn *string) IEventBus {
	_init_.Initialize()

	var returns IEventBus

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventBus",
		"fromEventBusArn",
		[]interface{}{scope, id, eventBusArn},
		&returns,
	)

	return returns
}

// Import an existing event bus resource.
// Experimental.
func EventBus_FromEventBusAttributes(scope constructs.Construct, id *string, attrs *EventBusAttributes) IEventBus {
	_init_.Initialize()

	var returns IEventBus

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventBus",
		"fromEventBusAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Import an existing event bus resource.
// Experimental.
func EventBus_FromEventBusName(scope constructs.Construct, id *string, eventBusName *string) IEventBus {
	_init_.Initialize()

	var returns IEventBus

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventBus",
		"fromEventBusName",
		[]interface{}{scope, id, eventBusName},
		&returns,
	)

	return returns
}

// Permits an IAM Principal to send custom events to EventBridge so that they can be matched to rules.
// Experimental.
func EventBus_GrantAllPutEvents(grantee awsiam.IGrantable) awsiam.Grant {
	_init_.Initialize()

	var returns awsiam.Grant

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventBus",
		"grantAllPutEvents",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Permits an IAM Principal to send custom events to EventBridge so that they can be matched to rules.
// Deprecated: use grantAllPutEvents instead
func EventBus_GrantPutEvents(grantee awsiam.IGrantable) awsiam.Grant {
	_init_.Initialize()

	var returns awsiam.Grant

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventBus",
		"grantPutEvents",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func EventBus_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventBus",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func EventBus_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventBus",
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
func (e *jsiiProxy_EventBus) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		e,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Create an EventBridge archive to send events to.
//
// When you create an archive, incoming events might not immediately start being sent to the archive.
// Allow a short period of time for changes to take effect.
// Experimental.
func (e *jsiiProxy_EventBus) Archive(id *string, props *BaseArchiveProps) Archive {
	var returns Archive

	_jsii_.Invoke(
		e,
		"archive",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Experimental.
func (e *jsiiProxy_EventBus) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		e,
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
func (e *jsiiProxy_EventBus) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		e,
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
func (e *jsiiProxy_EventBus) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		e,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grants an IAM Principal to send custom events to the eventBus so that they can be matched to rules.
// Experimental.
func (e *jsiiProxy_EventBus) GrantPutEventsTo(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		e,
		"grantPutEventsTo",
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
func (e *jsiiProxy_EventBus) OnPrepare() {
	_jsii_.InvokeVoid(
		e,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (e *jsiiProxy_EventBus) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		e,
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
func (e *jsiiProxy_EventBus) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		e,
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
func (e *jsiiProxy_EventBus) Prepare() {
	_jsii_.InvokeVoid(
		e,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (e *jsiiProxy_EventBus) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		e,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (e *jsiiProxy_EventBus) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		e,
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
func (e *jsiiProxy_EventBus) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		e,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Interface with properties necessary to import a reusable EventBus.
// Experimental.
type EventBusAttributes struct {
	// The ARN of this event bus resource.
	// Experimental.
	EventBusArn *string `json:"eventBusArn"`
	// The physical ID of this event bus resource.
	// Experimental.
	EventBusName *string `json:"eventBusName"`
	// The JSON policy of this event bus resource.
	// Experimental.
	EventBusPolicy *string `json:"eventBusPolicy"`
	// The partner event source to associate with this event bus resource.
	// Experimental.
	EventSourceName *string `json:"eventSourceName"`
}

// Properties to define an event bus.
// Experimental.
type EventBusProps struct {
	// The name of the event bus you are creating Note: If 'eventSourceName' is passed in, you cannot set this.
	// Experimental.
	EventBusName *string `json:"eventBusName"`
	// The partner event source to associate with this event bus resource Note: If 'eventBusName' is passed in, you cannot set this.
	// Experimental.
	EventSourceName *string `json:"eventSourceName"`
}

// Represents a field in the event pattern.
// Experimental.
type EventField interface {
	awscdk.IResolvable
	CreationStack() *[]*string
	DisplayHint() *string
	Path() *string
	Resolve(_ctx awscdk.IResolveContext) interface{}
	ToJSON() *string
	ToString() *string
}

// The jsii proxy struct for EventField
type jsiiProxy_EventField struct {
	internal.Type__awscdkIResolvable
}

func (j *jsiiProxy_EventField) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventField) DisplayHint() *string {
	var returns *string
	_jsii_.Get(
		j,
		"displayHint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventField) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}


// Extract a custom JSON path from the event.
// Experimental.
func EventField_FromPath(path *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_events.EventField",
		"fromPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

func EventField_Account() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.EventField",
		"account",
		&returns,
	)
	return returns
}

func EventField_DetailType() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.EventField",
		"detailType",
		&returns,
	)
	return returns
}

func EventField_EventId() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.EventField",
		"eventId",
		&returns,
	)
	return returns
}

func EventField_Region() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.EventField",
		"region",
		&returns,
	)
	return returns
}

func EventField_Source() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.EventField",
		"source",
		&returns,
	)
	return returns
}

func EventField_Time() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_events.EventField",
		"time",
		&returns,
	)
	return returns
}

// Produce the Token's value at resolution time.
// Experimental.
func (e *jsiiProxy_EventField) Resolve(_ctx awscdk.IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		e,
		"resolve",
		[]interface{}{_ctx},
		&returns,
	)

	return returns
}

// Convert the path to the field in the event pattern to JSON.
// Experimental.
func (e *jsiiProxy_EventField) ToJSON() *string {
	var returns *string

	_jsii_.Invoke(
		e,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return a string representation of this resolvable object.
//
// Returns a reversible string representation.
// Experimental.
func (e *jsiiProxy_EventField) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		e,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Events in Amazon CloudWatch Events are represented as JSON objects. For more information about JSON objects, see RFC 7159.
//
// Rules use event patterns to select events and route them to targets. A
// pattern either matches an event or it doesn't. Event patterns are represented
// as JSON objects with a structure that is similar to that of events, for
// example:
//
// It is important to remember the following about event pattern matching:
//
// - For a pattern to match an event, the event must contain all the field names
//    listed in the pattern. The field names must appear in the event with the
//    same nesting structure.
//
// - Other fields of the event not mentioned in the pattern are ignored;
//    effectively, there is a ``"*": "*"`` wildcard for fields not mentioned.
//
// - The matching is exact (character-by-character), without case-folding or any
//    other string normalization.
//
// - The values being matched follow JSON rules: Strings enclosed in quotes,
//    numbers, and the unquoted keywords true, false, and null.
//
// - Number matching is at the string representation level. For example, 300,
//    300.0, and 3.0e2 are not considered equal.
// See: https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/CloudWatchEventsandEventPatterns.html
//
// Experimental.
type EventPattern struct {
	// The 12-digit number identifying an AWS account.
	// Experimental.
	Account *[]*string `json:"account"`
	// A JSON object, whose content is at the discretion of the service originating the event.
	// Experimental.
	Detail *map[string]interface{} `json:"detail"`
	// Identifies, in combination with the source field, the fields and values that appear in the detail field.
	//
	// Represents the "detail-type" event field.
	// Experimental.
	DetailType *[]*string `json:"detailType"`
	// A unique value is generated for every event.
	//
	// This can be helpful in
	// tracing events as they move through rules to targets, and are processed.
	// Experimental.
	Id *[]*string `json:"id"`
	// Identifies the AWS region where the event originated.
	// Experimental.
	Region *[]*string `json:"region"`
	// This JSON array contains ARNs that identify resources that are involved in the event.
	//
	// Inclusion of these ARNs is at the discretion of the
	// service.
	//
	// For example, Amazon EC2 instance state-changes include Amazon EC2
	// instance ARNs, Auto Scaling events include ARNs for both instances and
	// Auto Scaling groups, but API calls with AWS CloudTrail do not include
	// resource ARNs.
	// Experimental.
	Resources *[]*string `json:"resources"`
	// Identifies the service that sourced the event.
	//
	// All events sourced from
	// within AWS begin with "aws." Customer-generated events can have any value
	// here, as long as it doesn't begin with "aws." We recommend the use of
	// Java package-name style reverse domain-name strings.
	//
	// To find the correct value for source for an AWS service, see the table in
	// AWS Service Namespaces. For example, the source value for Amazon
	// CloudFront is aws.cloudfront.
	// See: http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html#genref-aws-service-namespaces
	//
	// Experimental.
	Source *[]*string `json:"source"`
	// The event timestamp, which can be specified by the service originating the event.
	//
	// If the event spans a time interval, the service might choose
	// to report the start time, so this value can be noticeably before the time
	// the event is actually received.
	// Experimental.
	Time *[]*string `json:"time"`
	// By default, this is set to 0 (zero) in all events.
	// Experimental.
	Version *[]*string `json:"version"`
}

// Interface which all EventBus based classes MUST implement.
// Experimental.
type IEventBus interface {
	awscdk.IResource
	// Create an EventBridge archive to send events to.
	//
	// When you create an archive, incoming events might not immediately start being sent to the archive.
	// Allow a short period of time for changes to take effect.
	// Experimental.
	Archive(id *string, props *BaseArchiveProps) Archive
	// Grants an IAM Principal to send custom events to the eventBus so that they can be matched to rules.
	// Experimental.
	GrantPutEventsTo(grantee awsiam.IGrantable) awsiam.Grant
	// The ARN of this event bus resource.
	// Experimental.
	EventBusArn() *string
	// The physical ID of this event bus resource.
	// Experimental.
	EventBusName() *string
	// The JSON policy of this event bus resource.
	// Experimental.
	EventBusPolicy() *string
	// The partner event source to associate with this event bus resource.
	// Experimental.
	EventSourceName() *string
}

// The jsii proxy for IEventBus
type jsiiProxy_IEventBus struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IEventBus) Archive(id *string, props *BaseArchiveProps) Archive {
	var returns Archive

	_jsii_.Invoke(
		i,
		"archive",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IEventBus) GrantPutEventsTo(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPutEventsTo",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IEventBus) EventBusArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IEventBus) EventBusName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IEventBus) EventBusPolicy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventBusPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IEventBus) EventSourceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventSourceName",
		&returns,
	)
	return returns
}

// Represents an EventBridge Rule.
// Experimental.
type IRule interface {
	awscdk.IResource
	// The value of the event rule Amazon Resource Name (ARN), such as arn:aws:events:us-east-2:123456789012:rule/example.
	// Experimental.
	RuleArn() *string
	// The name event rule.
	// Experimental.
	RuleName() *string
}

// The jsii proxy for IRule
type jsiiProxy_IRule struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IRule) RuleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ruleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRule) RuleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ruleName",
		&returns,
	)
	return returns
}

// An abstract target for EventRules.
// Experimental.
type IRuleTarget interface {
	// Returns the rule target specification.
	//
	// NOTE: Do not use the various `inputXxx` options. They can be set in a call to `addTarget`.
	// Experimental.
	Bind(rule IRule, id *string) *RuleTargetConfig
}

// The jsii proxy for IRuleTarget
type jsiiProxy_IRuleTarget struct {
	_ byte // padding
}

func (i *jsiiProxy_IRuleTarget) Bind(rule IRule, id *string) *RuleTargetConfig {
	var returns *RuleTargetConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{rule, id},
		&returns,
	)

	return returns
}

// Standard set of options for `onXxx` event handlers on construct.
// Experimental.
type OnEventOptions struct {
	// A description of the rule's purpose.
	// Experimental.
	Description *string `json:"description"`
	// Additional restrictions for the event to route to the specified target.
	//
	// The method that generates the rule probably imposes some type of event
	// filtering. The filtering implied by what you pass here is added
	// on top of that filtering.
	// See: https://docs.aws.amazon.com/eventbridge/latest/userguide/eventbridge-and-event-patterns.html
	//
	// Experimental.
	EventPattern *EventPattern `json:"eventPattern"`
	// A name for the rule.
	// Experimental.
	RuleName *string `json:"ruleName"`
	// The target to register for the event.
	// Experimental.
	Target IRuleTarget `json:"target"`
}

// Defines an EventBridge Rule in this stack.
// Experimental.
type Rule interface {
	awscdk.Resource
	IRule
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RuleArn() *string
	RuleName() *string
	Stack() awscdk.Stack
	AddEventPattern(eventPattern *EventPattern)
	AddTarget(target IRuleTarget)
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

// The jsii proxy struct for Rule
type jsiiProxy_Rule struct {
	internal.Type__awscdkResource
	jsiiProxy_IRule
}

func (j *jsiiProxy_Rule) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Rule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Rule) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Rule) RuleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ruleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Rule) RuleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ruleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Rule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRule(scope constructs.Construct, id *string, props *RuleProps) Rule {
	_init_.Initialize()

	j := jsiiProxy_Rule{}

	_jsii_.Create(
		"monocdk.aws_events.Rule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRule_Override(r Rule, scope constructs.Construct, id *string, props *RuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.Rule",
		[]interface{}{scope, id, props},
		r,
	)
}

// Import an existing EventBridge Rule provided an ARN.
// Experimental.
func Rule_FromEventRuleArn(scope constructs.Construct, id *string, eventRuleArn *string) IRule {
	_init_.Initialize()

	var returns IRule

	_jsii_.StaticInvoke(
		"monocdk.aws_events.Rule",
		"fromEventRuleArn",
		[]interface{}{scope, id, eventRuleArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Rule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.Rule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Rule_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_events.Rule",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds an event pattern filter to this rule.
//
// If a pattern was already specified,
// these values are merged into the existing pattern.
//
// For example, if the rule already contains the pattern:
//
//     {
//       "resources": [ "r1" ],
//       "detail": {
//         "hello": [ 1 ]
//       }
//     }
//
// And `addEventPattern` is called with the pattern:
//
//     {
//       "resources": [ "r2" ],
//       "detail": {
//         "foo": [ "bar" ]
//       }
//     }
//
// The resulting event pattern will be:
//
//     {
//       "resources": [ "r1", "r2" ],
//       "detail": {
//         "hello": [ 1 ],
//         "foo": [ "bar" ]
//       }
//     }
// Experimental.
func (r *jsiiProxy_Rule) AddEventPattern(eventPattern *EventPattern) {
	_jsii_.InvokeVoid(
		r,
		"addEventPattern",
		[]interface{}{eventPattern},
	)
}

// Adds a target to the rule. The abstract class RuleTarget can be extended to define new targets.
//
// No-op if target is undefined.
// Experimental.
func (r *jsiiProxy_Rule) AddTarget(target IRuleTarget) {
	_jsii_.InvokeVoid(
		r,
		"addTarget",
		[]interface{}{target},
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
func (r *jsiiProxy_Rule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_Rule) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_Rule) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_Rule) GetResourceNameAttribute(nameAttr *string) *string {
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
func (r *jsiiProxy_Rule) OnPrepare() {
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
func (r *jsiiProxy_Rule) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_Rule) OnValidate() *[]*string {
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
func (r *jsiiProxy_Rule) Prepare() {
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
func (r *jsiiProxy_Rule) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_Rule) ToString() *string {
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
// Experimental.
func (r *jsiiProxy_Rule) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining an EventBridge Rule.
// Experimental.
type RuleProps struct {
	// A description of the rule's purpose.
	// Experimental.
	Description *string `json:"description"`
	// Indicates whether the rule is enabled.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// The event bus to associate with this rule.
	// Experimental.
	EventBus IEventBus `json:"eventBus"`
	// Describes which events EventBridge routes to the specified target.
	//
	// These routed events are matched events. For more information, see Events
	// and Event Patterns in the Amazon EventBridge User Guide.
	// See: https://docs.aws.amazon.com/eventbridge/latest/userguide/eventbridge-and-event-patterns.html
	//
	// You must specify this property (either via props or via
	// `addEventPattern`), the `scheduleExpression` property, or both. The
	// method `addEventPattern` can be used to add filter values to the event
	// pattern.
	//
	// Experimental.
	EventPattern *EventPattern `json:"eventPattern"`
	// A name for the rule.
	// Experimental.
	RuleName *string `json:"ruleName"`
	// The schedule or rate (frequency) that determines when EventBridge runs the rule.
	//
	// For more information, see Schedule Expression Syntax for
	// Rules in the Amazon EventBridge User Guide.
	// See: https://docs.aws.amazon.com/eventbridge/latest/userguide/scheduled-events.html
	//
	// You must specify this property, the `eventPattern` property, or both.
	//
	// Experimental.
	Schedule Schedule `json:"schedule"`
	// Targets to invoke when this rule matches an event.
	//
	// Input will be the full matched event. If you wish to specify custom
	// target input, use `addTarget(target[, inputOptions])`.
	// Experimental.
	Targets *[]IRuleTarget `json:"targets"`
}

// Properties for an event rule target.
// Experimental.
type RuleTargetConfig struct {
	// The Amazon Resource Name (ARN) of the target.
	// Experimental.
	Arn *string `json:"arn"`
	// Parameters used when the rule invokes Amazon AWS Batch Job/Queue.
	// Experimental.
	BatchParameters *CfnRule_BatchParametersProperty `json:"batchParameters"`
	// Contains information about a dead-letter queue configuration.
	// Experimental.
	DeadLetterConfig *CfnRule_DeadLetterConfigProperty `json:"deadLetterConfig"`
	// The Amazon ECS task definition and task count to use, if the event target is an Amazon ECS task.
	// Experimental.
	EcsParameters *CfnRule_EcsParametersProperty `json:"ecsParameters"`
	// Parameters used when the rule invoke api gateway.
	// Experimental.
	HttpParameters *CfnRule_HttpParametersProperty `json:"httpParameters"`
	// A unique, user-defined identifier for the target.
	//
	// Acceptable values
	// include alphanumeric characters, periods (.), hyphens (-), and
	// underscores (_).
	// Deprecated: no replacement. we will always use an autogenerated id.
	Id *string `json:"id"`
	// What input to send to the event target.
	// Experimental.
	Input RuleTargetInput `json:"input"`
	// Settings that control shard assignment, when the target is a Kinesis stream.
	//
	// If you don't include this parameter, eventId is used as the
	// partition key.
	// Experimental.
	KinesisParameters *CfnRule_KinesisParametersProperty `json:"kinesisParameters"`
	// A RetryPolicy object that includes information about the retry policy settings.
	// Experimental.
	RetryPolicy *CfnRule_RetryPolicyProperty `json:"retryPolicy"`
	// Role to use to invoke this event target.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// Parameters used when the rule invokes Amazon EC2 Systems Manager Run Command.
	// Experimental.
	RunCommandParameters *CfnRule_RunCommandParametersProperty `json:"runCommandParameters"`
	// Parameters used when the FIFO sqs queue is used an event target by the rule.
	// Experimental.
	SqsParameters *CfnRule_SqsParametersProperty `json:"sqsParameters"`
	// The resource that is backing this target.
	//
	// This is the resource that will actually have some action performed on it when used as a target
	// (for example, start a build for a CodeBuild project).
	// We need it to determine whether the rule belongs to a different account than the target -
	// if so, we generate a more complex setup,
	// including an additional stack containing the EventBusPolicy.
	// See: https://docs.aws.amazon.com/eventbridge/latest/userguide/eventbridge-cross-account-event-delivery.html
	//
	// Experimental.
	TargetResource awscdk.IConstruct `json:"targetResource"`
}

// The input to send to the event target.
// Experimental.
type RuleTargetInput interface {
	Bind(rule IRule) *RuleTargetInputProperties
}

// The jsii proxy struct for RuleTargetInput
type jsiiProxy_RuleTargetInput struct {
	_ byte // padding
}

// Experimental.
func NewRuleTargetInput_Override(r RuleTargetInput) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_events.RuleTargetInput",
		nil, // no parameters
		r,
	)
}

// Take the event target input from a path in the event JSON.
// Experimental.
func RuleTargetInput_FromEventPath(path *string) RuleTargetInput {
	_init_.Initialize()

	var returns RuleTargetInput

	_jsii_.StaticInvoke(
		"monocdk.aws_events.RuleTargetInput",
		"fromEventPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Pass text to the event target, splitting on newlines.
//
// This is only useful when passing to a target that does not
// take a single argument.
//
// May contain strings returned by EventField.from() to substitute in parts
// of the matched event.
// Experimental.
func RuleTargetInput_FromMultilineText(text *string) RuleTargetInput {
	_init_.Initialize()

	var returns RuleTargetInput

	_jsii_.StaticInvoke(
		"monocdk.aws_events.RuleTargetInput",
		"fromMultilineText",
		[]interface{}{text},
		&returns,
	)

	return returns
}

// Pass a JSON object to the event target.
//
// May contain strings returned by EventField.from() to substitute in parts of the
// matched event.
// Experimental.
func RuleTargetInput_FromObject(obj interface{}) RuleTargetInput {
	_init_.Initialize()

	var returns RuleTargetInput

	_jsii_.StaticInvoke(
		"monocdk.aws_events.RuleTargetInput",
		"fromObject",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Pass text to the event target.
//
// May contain strings returned by EventField.from() to substitute in parts of the
// matched event.
// Experimental.
func RuleTargetInput_FromText(text *string) RuleTargetInput {
	_init_.Initialize()

	var returns RuleTargetInput

	_jsii_.StaticInvoke(
		"monocdk.aws_events.RuleTargetInput",
		"fromText",
		[]interface{}{text},
		&returns,
	)

	return returns
}

// Return the input properties for this input object.
// Experimental.
func (r *jsiiProxy_RuleTargetInput) Bind(rule IRule) *RuleTargetInputProperties {
	var returns *RuleTargetInputProperties

	_jsii_.Invoke(
		r,
		"bind",
		[]interface{}{rule},
		&returns,
	)

	return returns
}

// The input properties for an event target.
// Experimental.
type RuleTargetInputProperties struct {
	// Literal input to the target service (must be valid JSON).
	// Experimental.
	Input *string `json:"input"`
	// JsonPath to take input from the input event.
	// Experimental.
	InputPath *string `json:"inputPath"`
	// Paths map to extract values from event and insert into `inputTemplate`.
	// Experimental.
	InputPathsMap *map[string]*string `json:"inputPathsMap"`
	// Input template to insert paths map into.
	// Experimental.
	InputTemplate *string `json:"inputTemplate"`
}

// Schedule for scheduled event rules.
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
		"monocdk.aws_events.Schedule",
		nil, // no parameters
		s,
	)
}

// Create a schedule from a set of cron fields.
// Experimental.
func Schedule_Cron(options *CronOptions) Schedule {
	_init_.Initialize()

	var returns Schedule

	_jsii_.StaticInvoke(
		"monocdk.aws_events.Schedule",
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
		"monocdk.aws_events.Schedule",
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
		"monocdk.aws_events.Schedule",
		"rate",
		[]interface{}{duration},
		&returns,
	)

	return returns
}


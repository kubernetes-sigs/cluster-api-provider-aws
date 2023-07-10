package awscodeguruprofiler

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscodeguruprofiler/internal"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
)

// A CloudFormation `AWS::CodeGuruProfiler::ProfilingGroup`.
type CfnProfilingGroup interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AgentPermissions() interface{}
	SetAgentPermissions(val interface{})
	AnomalyDetectionNotificationConfiguration() interface{}
	SetAnomalyDetectionNotificationConfiguration(val interface{})
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ComputePlatform() *string
	SetComputePlatform(val *string)
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	ProfilingGroupName() *string
	SetProfilingGroupName(val *string)
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

// The jsii proxy struct for CfnProfilingGroup
type jsiiProxy_CfnProfilingGroup struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnProfilingGroup) AgentPermissions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"agentPermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) AnomalyDetectionNotificationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"anomalyDetectionNotificationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) ComputePlatform() *string {
	var returns *string
	_jsii_.Get(
		j,
		"computePlatform",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) ProfilingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"profilingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProfilingGroup) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodeGuruProfiler::ProfilingGroup`.
func NewCfnProfilingGroup(scope awscdk.Construct, id *string, props *CfnProfilingGroupProps) CfnProfilingGroup {
	_init_.Initialize()

	j := jsiiProxy_CfnProfilingGroup{}

	_jsii_.Create(
		"monocdk.aws_codeguruprofiler.CfnProfilingGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodeGuruProfiler::ProfilingGroup`.
func NewCfnProfilingGroup_Override(c CfnProfilingGroup, scope awscdk.Construct, id *string, props *CfnProfilingGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codeguruprofiler.CfnProfilingGroup",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnProfilingGroup) SetAgentPermissions(val interface{}) {
	_jsii_.Set(
		j,
		"agentPermissions",
		val,
	)
}

func (j *jsiiProxy_CfnProfilingGroup) SetAnomalyDetectionNotificationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"anomalyDetectionNotificationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnProfilingGroup) SetComputePlatform(val *string) {
	_jsii_.Set(
		j,
		"computePlatform",
		val,
	)
}

func (j *jsiiProxy_CfnProfilingGroup) SetProfilingGroupName(val *string) {
	_jsii_.Set(
		j,
		"profilingGroupName",
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
func CfnProfilingGroup_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codeguruprofiler.CfnProfilingGroup",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnProfilingGroup_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codeguruprofiler.CfnProfilingGroup",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnProfilingGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codeguruprofiler.CfnProfilingGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnProfilingGroup_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codeguruprofiler.CfnProfilingGroup",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnProfilingGroup) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnProfilingGroup) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnProfilingGroup) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnProfilingGroup) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnProfilingGroup) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnProfilingGroup) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnProfilingGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnProfilingGroup) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnProfilingGroup) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnProfilingGroup) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnProfilingGroup) OnPrepare() {
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
func (c *jsiiProxy_CfnProfilingGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnProfilingGroup) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnProfilingGroup) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnProfilingGroup) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnProfilingGroup) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnProfilingGroup) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnProfilingGroup) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnProfilingGroup) ToString() *string {
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
func (c *jsiiProxy_CfnProfilingGroup) Validate() *[]*string {
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
func (c *jsiiProxy_CfnProfilingGroup) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnProfilingGroup_ChannelProperty struct {
	// `CfnProfilingGroup.ChannelProperty.channelUri`.
	ChannelUri *string `json:"channelUri"`
	// `CfnProfilingGroup.ChannelProperty.channelId`.
	ChannelId *string `json:"channelId"`
}

// Properties for defining a `AWS::CodeGuruProfiler::ProfilingGroup`.
type CfnProfilingGroupProps struct {
	// `AWS::CodeGuruProfiler::ProfilingGroup.ProfilingGroupName`.
	ProfilingGroupName *string `json:"profilingGroupName"`
	// `AWS::CodeGuruProfiler::ProfilingGroup.AgentPermissions`.
	AgentPermissions interface{} `json:"agentPermissions"`
	// `AWS::CodeGuruProfiler::ProfilingGroup.AnomalyDetectionNotificationConfiguration`.
	AnomalyDetectionNotificationConfiguration interface{} `json:"anomalyDetectionNotificationConfiguration"`
	// `AWS::CodeGuruProfiler::ProfilingGroup.ComputePlatform`.
	ComputePlatform *string `json:"computePlatform"`
	// `AWS::CodeGuruProfiler::ProfilingGroup.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// The compute platform of the profiling group.
// Experimental.
type ComputePlatform string

const (
	ComputePlatform_AWS_LAMBDA ComputePlatform = "AWS_LAMBDA"
	ComputePlatform_DEFAULT ComputePlatform = "DEFAULT"
)

// IResource represents a Profiling Group.
// Experimental.
type IProfilingGroup interface {
	awscdk.IResource
	// Grant access to publish profiling information to the Profiling Group to the given identity.
	//
	// This will grant the following permissions:
	//
	//   - codeguru-profiler:ConfigureAgent
	//   - codeguru-profiler:PostAgentProfile
	// Experimental.
	GrantPublish(grantee awsiam.IGrantable) awsiam.Grant
	// Grant access to read profiling information from the Profiling Group to the given identity.
	//
	// This will grant the following permissions:
	//
	//   - codeguru-profiler:GetProfile
	//   - codeguru-profiler:DescribeProfilingGroup
	// Experimental.
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	// A name for the profiling group.
	// Experimental.
	ProfilingGroupName() *string
}

// The jsii proxy for IProfilingGroup
type jsiiProxy_IProfilingGroup struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IProfilingGroup) GrantPublish(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPublish",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProfilingGroup) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IProfilingGroup) ProfilingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"profilingGroupName",
		&returns,
	)
	return returns
}

// A new Profiling Group.
// Experimental.
type ProfilingGroup interface {
	awscdk.Resource
	IProfilingGroup
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ProfilingGroupArn() *string
	ProfilingGroupName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantPublish(grantee awsiam.IGrantable) awsiam.Grant
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ProfilingGroup
type jsiiProxy_ProfilingGroup struct {
	internal.Type__awscdkResource
	jsiiProxy_IProfilingGroup
}

func (j *jsiiProxy_ProfilingGroup) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProfilingGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProfilingGroup) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProfilingGroup) ProfilingGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"profilingGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProfilingGroup) ProfilingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"profilingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProfilingGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewProfilingGroup(scope constructs.Construct, id *string, props *ProfilingGroupProps) ProfilingGroup {
	_init_.Initialize()

	j := jsiiProxy_ProfilingGroup{}

	_jsii_.Create(
		"monocdk.aws_codeguruprofiler.ProfilingGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewProfilingGroup_Override(p ProfilingGroup, scope constructs.Construct, id *string, props *ProfilingGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codeguruprofiler.ProfilingGroup",
		[]interface{}{scope, id, props},
		p,
	)
}

// Import an existing Profiling Group provided an ARN.
// Experimental.
func ProfilingGroup_FromProfilingGroupArn(scope constructs.Construct, id *string, profilingGroupArn *string) IProfilingGroup {
	_init_.Initialize()

	var returns IProfilingGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_codeguruprofiler.ProfilingGroup",
		"fromProfilingGroupArn",
		[]interface{}{scope, id, profilingGroupArn},
		&returns,
	)

	return returns
}

// Import an existing Profiling Group provided a Profiling Group Name.
// Experimental.
func ProfilingGroup_FromProfilingGroupName(scope constructs.Construct, id *string, profilingGroupName *string) IProfilingGroup {
	_init_.Initialize()

	var returns IProfilingGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_codeguruprofiler.ProfilingGroup",
		"fromProfilingGroupName",
		[]interface{}{scope, id, profilingGroupName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ProfilingGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codeguruprofiler.ProfilingGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ProfilingGroup_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codeguruprofiler.ProfilingGroup",
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
func (p *jsiiProxy_ProfilingGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		p,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (p *jsiiProxy_ProfilingGroup) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		p,
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
func (p *jsiiProxy_ProfilingGroup) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		p,
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
func (p *jsiiProxy_ProfilingGroup) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant access to publish profiling information to the Profiling Group to the given identity.
//
// This will grant the following permissions:
//
//   - codeguru-profiler:ConfigureAgent
//   - codeguru-profiler:PostAgentProfile
// Experimental.
func (p *jsiiProxy_ProfilingGroup) GrantPublish(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		p,
		"grantPublish",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant access to read profiling information from the Profiling Group to the given identity.
//
// This will grant the following permissions:
//
//   - codeguru-profiler:GetProfile
//   - codeguru-profiler:DescribeProfilingGroup
// Experimental.
func (p *jsiiProxy_ProfilingGroup) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		p,
		"grantRead",
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
func (p *jsiiProxy_ProfilingGroup) OnPrepare() {
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
func (p *jsiiProxy_ProfilingGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (p *jsiiProxy_ProfilingGroup) OnValidate() *[]*string {
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
func (p *jsiiProxy_ProfilingGroup) Prepare() {
	_jsii_.InvokeVoid(
		p,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_ProfilingGroup) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_ProfilingGroup) ToString() *string {
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
func (p *jsiiProxy_ProfilingGroup) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for creating a new Profiling Group.
// Experimental.
type ProfilingGroupProps struct {
	// The compute platform of the profiling group.
	// Experimental.
	ComputePlatform ComputePlatform `json:"computePlatform"`
	// A name for the profiling group.
	// Experimental.
	ProfilingGroupName *string `json:"profilingGroupName"`
}


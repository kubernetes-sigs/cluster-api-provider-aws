package awscodestarnotifications

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscodestarnotifications/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// A CloudFormation `AWS::CodeStarNotifications::NotificationRule`.
type CfnNotificationRule interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DetailType() *string
	SetDetailType(val *string)
	EventTypeIds() *[]*string
	SetEventTypeIds(val *[]*string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Resource() *string
	SetResource(val *string)
	Stack() awscdk.Stack
	Status() *string
	SetStatus(val *string)
	Tags() awscdk.TagManager
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

// The jsii proxy struct for CfnNotificationRule
type jsiiProxy_CfnNotificationRule struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnNotificationRule) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) DetailType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"detailType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) EventTypeIds() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"eventTypeIds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Resource() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Status() *string {
	var returns *string
	_jsii_.Get(
		j,
		"status",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) Targets() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"targets",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnNotificationRule) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodeStarNotifications::NotificationRule`.
func NewCfnNotificationRule(scope awscdk.Construct, id *string, props *CfnNotificationRuleProps) CfnNotificationRule {
	_init_.Initialize()

	j := jsiiProxy_CfnNotificationRule{}

	_jsii_.Create(
		"monocdk.aws_codestarnotifications.CfnNotificationRule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodeStarNotifications::NotificationRule`.
func NewCfnNotificationRule_Override(c CfnNotificationRule, scope awscdk.Construct, id *string, props *CfnNotificationRuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codestarnotifications.CfnNotificationRule",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnNotificationRule) SetDetailType(val *string) {
	_jsii_.Set(
		j,
		"detailType",
		val,
	)
}

func (j *jsiiProxy_CfnNotificationRule) SetEventTypeIds(val *[]*string) {
	_jsii_.Set(
		j,
		"eventTypeIds",
		val,
	)
}

func (j *jsiiProxy_CfnNotificationRule) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnNotificationRule) SetResource(val *string) {
	_jsii_.Set(
		j,
		"resource",
		val,
	)
}

func (j *jsiiProxy_CfnNotificationRule) SetStatus(val *string) {
	_jsii_.Set(
		j,
		"status",
		val,
	)
}

func (j *jsiiProxy_CfnNotificationRule) SetTargets(val interface{}) {
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
func CfnNotificationRule_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codestarnotifications.CfnNotificationRule",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnNotificationRule_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codestarnotifications.CfnNotificationRule",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnNotificationRule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codestarnotifications.CfnNotificationRule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnNotificationRule_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codestarnotifications.CfnNotificationRule",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnNotificationRule) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnNotificationRule) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnNotificationRule) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnNotificationRule) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnNotificationRule) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnNotificationRule) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnNotificationRule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnNotificationRule) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnNotificationRule) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnNotificationRule) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnNotificationRule) OnPrepare() {
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
func (c *jsiiProxy_CfnNotificationRule) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnNotificationRule) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnNotificationRule) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnNotificationRule) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnNotificationRule) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnNotificationRule) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnNotificationRule) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnNotificationRule) ToString() *string {
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
func (c *jsiiProxy_CfnNotificationRule) Validate() *[]*string {
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
func (c *jsiiProxy_CfnNotificationRule) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnNotificationRule_TargetProperty struct {
	// `CfnNotificationRule.TargetProperty.TargetAddress`.
	TargetAddress *string `json:"targetAddress"`
	// `CfnNotificationRule.TargetProperty.TargetType`.
	TargetType *string `json:"targetType"`
}

// Properties for defining a `AWS::CodeStarNotifications::NotificationRule`.
type CfnNotificationRuleProps struct {
	// `AWS::CodeStarNotifications::NotificationRule.DetailType`.
	DetailType *string `json:"detailType"`
	// `AWS::CodeStarNotifications::NotificationRule.EventTypeIds`.
	EventTypeIds *[]*string `json:"eventTypeIds"`
	// `AWS::CodeStarNotifications::NotificationRule.Name`.
	Name *string `json:"name"`
	// `AWS::CodeStarNotifications::NotificationRule.Resource`.
	Resource *string `json:"resource"`
	// `AWS::CodeStarNotifications::NotificationRule.Targets`.
	Targets interface{} `json:"targets"`
	// `AWS::CodeStarNotifications::NotificationRule.Status`.
	Status *string `json:"status"`
	// `AWS::CodeStarNotifications::NotificationRule.Tags`.
	Tags interface{} `json:"tags"`
}

// The level of detail to include in the notifications for this resource.
// Experimental.
type DetailType string

const (
	DetailType_BASIC DetailType = "BASIC"
	DetailType_FULL DetailType = "FULL"
)

// Represents a notification rule.
// Experimental.
type INotificationRule interface {
	awscdk.IResource
	// Adds target to notification rule.
	//
	// Returns: boolean - return true if it had any effect
	// Experimental.
	AddTarget(target INotificationRuleTarget) *bool
	// The ARN of the notification rule (i.e. arn:aws:codestar-notifications:::notificationrule/01234abcde).
	// Experimental.
	NotificationRuleArn() *string
}

// The jsii proxy for INotificationRule
type jsiiProxy_INotificationRule struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_INotificationRule) AddTarget(target INotificationRuleTarget) *bool {
	var returns *bool

	_jsii_.Invoke(
		i,
		"addTarget",
		[]interface{}{target},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_INotificationRule) NotificationRuleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"notificationRuleArn",
		&returns,
	)
	return returns
}

// Represents a notification source The source that allows CodeBuild and CodePipeline to associate with this rule.
// Experimental.
type INotificationRuleSource interface {
	// Returns a source configuration for notification rule.
	// Experimental.
	BindAsNotificationRuleSource(scope constructs.Construct) *NotificationRuleSourceConfig
}

// The jsii proxy for INotificationRuleSource
type jsiiProxy_INotificationRuleSource struct {
	_ byte // padding
}

func (i *jsiiProxy_INotificationRuleSource) BindAsNotificationRuleSource(scope constructs.Construct) *NotificationRuleSourceConfig {
	var returns *NotificationRuleSourceConfig

	_jsii_.Invoke(
		i,
		"bindAsNotificationRuleSource",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Represents a notification target That allows AWS Chatbot and SNS topic to associate with this rule target.
// Experimental.
type INotificationRuleTarget interface {
	// Returns a target configuration for notification rule.
	// Experimental.
	BindAsNotificationRuleTarget(scope constructs.Construct) *NotificationRuleTargetConfig
}

// The jsii proxy for INotificationRuleTarget
type jsiiProxy_INotificationRuleTarget struct {
	_ byte // padding
}

func (i *jsiiProxy_INotificationRuleTarget) BindAsNotificationRuleTarget(scope constructs.Construct) *NotificationRuleTargetConfig {
	var returns *NotificationRuleTargetConfig

	_jsii_.Invoke(
		i,
		"bindAsNotificationRuleTarget",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// A new notification rule.
// Experimental.
type NotificationRule interface {
	awscdk.Resource
	INotificationRule
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	NotificationRuleArn() *string
	PhysicalName() *string
	Stack() awscdk.Stack
	AddTarget(target INotificationRuleTarget) *bool
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

// The jsii proxy struct for NotificationRule
type jsiiProxy_NotificationRule struct {
	internal.Type__awscdkResource
	jsiiProxy_INotificationRule
}

func (j *jsiiProxy_NotificationRule) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NotificationRule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NotificationRule) NotificationRuleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"notificationRuleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NotificationRule) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NotificationRule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewNotificationRule(scope constructs.Construct, id *string, props *NotificationRuleProps) NotificationRule {
	_init_.Initialize()

	j := jsiiProxy_NotificationRule{}

	_jsii_.Create(
		"monocdk.aws_codestarnotifications.NotificationRule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewNotificationRule_Override(n NotificationRule, scope constructs.Construct, id *string, props *NotificationRuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codestarnotifications.NotificationRule",
		[]interface{}{scope, id, props},
		n,
	)
}

// Import an existing notification rule provided an ARN.
// Experimental.
func NotificationRule_FromNotificationRuleArn(scope constructs.Construct, id *string, notificationRuleArn *string) INotificationRule {
	_init_.Initialize()

	var returns INotificationRule

	_jsii_.StaticInvoke(
		"monocdk.aws_codestarnotifications.NotificationRule",
		"fromNotificationRuleArn",
		[]interface{}{scope, id, notificationRuleArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func NotificationRule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codestarnotifications.NotificationRule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func NotificationRule_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codestarnotifications.NotificationRule",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds target to notification rule.
// Experimental.
func (n *jsiiProxy_NotificationRule) AddTarget(target INotificationRuleTarget) *bool {
	var returns *bool

	_jsii_.Invoke(
		n,
		"addTarget",
		[]interface{}{target},
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
func (n *jsiiProxy_NotificationRule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		n,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (n *jsiiProxy_NotificationRule) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NotificationRule) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NotificationRule) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NotificationRule) OnPrepare() {
	_jsii_.InvokeVoid(
		n,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NotificationRule) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
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
func (n *jsiiProxy_NotificationRule) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NotificationRule) Prepare() {
	_jsii_.InvokeVoid(
		n,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NotificationRule) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (n *jsiiProxy_NotificationRule) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NotificationRule) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Standard set of options for `notifyOnXxx` codestar notification handler on construct.
// Experimental.
type NotificationRuleOptions struct {
	// The level of detail to include in the notifications for this resource.
	//
	// BASIC will include only the contents of the event as it would appear in AWS CloudWatch.
	// FULL will include any supplemental information provided by AWS CodeStar Notifications and/or the service for the resource for which the notification is created.
	// Experimental.
	DetailType DetailType `json:"detailType"`
	// The status of the notification rule.
	//
	// If the enabled is set to DISABLED, notifications aren't sent for the notification rule.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// The name for the notification rule.
	//
	// Notification rule names must be unique in your AWS account.
	// Experimental.
	NotificationRuleName *string `json:"notificationRuleName"`
}

// Properties for a new notification rule.
// Experimental.
type NotificationRuleProps struct {
	// The level of detail to include in the notifications for this resource.
	//
	// BASIC will include only the contents of the event as it would appear in AWS CloudWatch.
	// FULL will include any supplemental information provided by AWS CodeStar Notifications and/or the service for the resource for which the notification is created.
	// Experimental.
	DetailType DetailType `json:"detailType"`
	// The status of the notification rule.
	//
	// If the enabled is set to DISABLED, notifications aren't sent for the notification rule.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// The name for the notification rule.
	//
	// Notification rule names must be unique in your AWS account.
	// Experimental.
	NotificationRuleName *string `json:"notificationRuleName"`
	// A list of event types associated with this notification rule.
	//
	// For a complete list of event types and IDs, see Notification concepts in the Developer Tools Console User Guide.
	// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#concepts-api
	//
	// Experimental.
	Events *[]*string `json:"events"`
	// The Amazon Resource Name (ARN) of the resource to associate with the notification rule.
	//
	// Currently, Supported sources include pipelines in AWS CodePipeline and build projects in AWS CodeBuild in this L2 constructor.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-codestarnotifications-notificationrule.html#cfn-codestarnotifications-notificationrule-resource
	//
	// Experimental.
	Source INotificationRuleSource `json:"source"`
	// The targets to register for the notification destination.
	// Experimental.
	Targets *[]INotificationRuleTarget `json:"targets"`
}

// Information about the Codebuild or CodePipeline associated with a notification source.
// Experimental.
type NotificationRuleSourceConfig struct {
	// The Amazon Resource Name (ARN) of the notification source.
	// Experimental.
	SourceArn *string `json:"sourceArn"`
}

// Information about the SNS topic or AWS Chatbot client associated with a notification target.
// Experimental.
type NotificationRuleTargetConfig struct {
	// The Amazon Resource Name (ARN) of the Amazon SNS topic or AWS Chatbot client.
	// Experimental.
	TargetAddress *string `json:"targetAddress"`
	// The target type.
	//
	// Can be an Amazon SNS topic or AWS Chatbot client.
	// Experimental.
	TargetType *string `json:"targetType"`
}


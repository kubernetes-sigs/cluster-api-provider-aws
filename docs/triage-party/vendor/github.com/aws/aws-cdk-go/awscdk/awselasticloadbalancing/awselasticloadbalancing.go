package awselasticloadbalancing

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancing/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// A CloudFormation `AWS::ElasticLoadBalancing::LoadBalancer`.
type CfnLoadBalancer interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AccessLoggingPolicy() interface{}
	SetAccessLoggingPolicy(val interface{})
	AppCookieStickinessPolicy() interface{}
	SetAppCookieStickinessPolicy(val interface{})
	AttrCanonicalHostedZoneName() *string
	AttrCanonicalHostedZoneNameId() *string
	AttrDnsName() *string
	AttrSourceSecurityGroupGroupName() *string
	AttrSourceSecurityGroupOwnerAlias() *string
	AvailabilityZones() *[]*string
	SetAvailabilityZones(val *[]*string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ConnectionDrainingPolicy() interface{}
	SetConnectionDrainingPolicy(val interface{})
	ConnectionSettings() interface{}
	SetConnectionSettings(val interface{})
	CreationStack() *[]*string
	CrossZone() interface{}
	SetCrossZone(val interface{})
	HealthCheck() interface{}
	SetHealthCheck(val interface{})
	Instances() *[]*string
	SetInstances(val *[]*string)
	LbCookieStickinessPolicy() interface{}
	SetLbCookieStickinessPolicy(val interface{})
	Listeners() interface{}
	SetListeners(val interface{})
	LoadBalancerName() *string
	SetLoadBalancerName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Policies() interface{}
	SetPolicies(val interface{})
	Ref() *string
	Scheme() *string
	SetScheme(val *string)
	SecurityGroups() *[]*string
	SetSecurityGroups(val *[]*string)
	Stack() awscdk.Stack
	Subnets() *[]*string
	SetSubnets(val *[]*string)
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

// The jsii proxy struct for CfnLoadBalancer
type jsiiProxy_CfnLoadBalancer struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnLoadBalancer) AccessLoggingPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accessLoggingPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AppCookieStickinessPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"appCookieStickinessPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrCanonicalHostedZoneName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCanonicalHostedZoneName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrCanonicalHostedZoneNameId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCanonicalHostedZoneNameId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrSourceSecurityGroupGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrSourceSecurityGroupGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrSourceSecurityGroupOwnerAlias() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrSourceSecurityGroupOwnerAlias",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AvailabilityZones() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"availabilityZones",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) ConnectionDrainingPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"connectionDrainingPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) ConnectionSettings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"connectionSettings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CrossZone() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"crossZone",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) HealthCheck() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"healthCheck",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Instances() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"instances",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) LbCookieStickinessPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"lbCookieStickinessPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Listeners() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"listeners",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) LoadBalancerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Policies() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Scheme() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scheme",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) SecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Subnets() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"subnets",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ElasticLoadBalancing::LoadBalancer`.
func NewCfnLoadBalancer(scope awscdk.Construct, id *string, props *CfnLoadBalancerProps) CfnLoadBalancer {
	_init_.Initialize()

	j := jsiiProxy_CfnLoadBalancer{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ElasticLoadBalancing::LoadBalancer`.
func NewCfnLoadBalancer_Override(c CfnLoadBalancer, scope awscdk.Construct, id *string, props *CfnLoadBalancerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetAccessLoggingPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"accessLoggingPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetAppCookieStickinessPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"appCookieStickinessPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetAvailabilityZones(val *[]*string) {
	_jsii_.Set(
		j,
		"availabilityZones",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetConnectionDrainingPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"connectionDrainingPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetConnectionSettings(val interface{}) {
	_jsii_.Set(
		j,
		"connectionSettings",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetCrossZone(val interface{}) {
	_jsii_.Set(
		j,
		"crossZone",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetHealthCheck(val interface{}) {
	_jsii_.Set(
		j,
		"healthCheck",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetInstances(val *[]*string) {
	_jsii_.Set(
		j,
		"instances",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetLbCookieStickinessPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"lbCookieStickinessPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetListeners(val interface{}) {
	_jsii_.Set(
		j,
		"listeners",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetLoadBalancerName(val *string) {
	_jsii_.Set(
		j,
		"loadBalancerName",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetPolicies(val interface{}) {
	_jsii_.Set(
		j,
		"policies",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetScheme(val *string) {
	_jsii_.Set(
		j,
		"scheme",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetSecurityGroups(val *[]*string) {
	_jsii_.Set(
		j,
		"securityGroups",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetSubnets(val *[]*string) {
	_jsii_.Set(
		j,
		"subnets",
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
func CfnLoadBalancer_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnLoadBalancer_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnLoadBalancer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnLoadBalancer_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_elasticloadbalancing.CfnLoadBalancer",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnLoadBalancer) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnLoadBalancer) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnLoadBalancer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnLoadBalancer) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnLoadBalancer) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnLoadBalancer) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnLoadBalancer) OnPrepare() {
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
func (c *jsiiProxy_CfnLoadBalancer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLoadBalancer) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnLoadBalancer) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnLoadBalancer) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnLoadBalancer) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnLoadBalancer) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnLoadBalancer) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLoadBalancer) ToString() *string {
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
func (c *jsiiProxy_CfnLoadBalancer) Validate() *[]*string {
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
func (c *jsiiProxy_CfnLoadBalancer) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnLoadBalancer_AccessLoggingPolicyProperty struct {
	// `CfnLoadBalancer.AccessLoggingPolicyProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `CfnLoadBalancer.AccessLoggingPolicyProperty.S3BucketName`.
	S3BucketName *string `json:"s3BucketName"`
	// `CfnLoadBalancer.AccessLoggingPolicyProperty.EmitInterval`.
	EmitInterval *float64 `json:"emitInterval"`
	// `CfnLoadBalancer.AccessLoggingPolicyProperty.S3BucketPrefix`.
	S3BucketPrefix *string `json:"s3BucketPrefix"`
}

type CfnLoadBalancer_AppCookieStickinessPolicyProperty struct {
	// `CfnLoadBalancer.AppCookieStickinessPolicyProperty.CookieName`.
	CookieName *string `json:"cookieName"`
	// `CfnLoadBalancer.AppCookieStickinessPolicyProperty.PolicyName`.
	PolicyName *string `json:"policyName"`
}

type CfnLoadBalancer_ConnectionDrainingPolicyProperty struct {
	// `CfnLoadBalancer.ConnectionDrainingPolicyProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `CfnLoadBalancer.ConnectionDrainingPolicyProperty.Timeout`.
	Timeout *float64 `json:"timeout"`
}

type CfnLoadBalancer_ConnectionSettingsProperty struct {
	// `CfnLoadBalancer.ConnectionSettingsProperty.IdleTimeout`.
	IdleTimeout *float64 `json:"idleTimeout"`
}

type CfnLoadBalancer_HealthCheckProperty struct {
	// `CfnLoadBalancer.HealthCheckProperty.HealthyThreshold`.
	HealthyThreshold *string `json:"healthyThreshold"`
	// `CfnLoadBalancer.HealthCheckProperty.Interval`.
	Interval *string `json:"interval"`
	// `CfnLoadBalancer.HealthCheckProperty.Target`.
	Target *string `json:"target"`
	// `CfnLoadBalancer.HealthCheckProperty.Timeout`.
	Timeout *string `json:"timeout"`
	// `CfnLoadBalancer.HealthCheckProperty.UnhealthyThreshold`.
	UnhealthyThreshold *string `json:"unhealthyThreshold"`
}

type CfnLoadBalancer_LBCookieStickinessPolicyProperty struct {
	// `CfnLoadBalancer.LBCookieStickinessPolicyProperty.CookieExpirationPeriod`.
	CookieExpirationPeriod *string `json:"cookieExpirationPeriod"`
	// `CfnLoadBalancer.LBCookieStickinessPolicyProperty.PolicyName`.
	PolicyName *string `json:"policyName"`
}

type CfnLoadBalancer_ListenersProperty struct {
	// `CfnLoadBalancer.ListenersProperty.InstancePort`.
	InstancePort *string `json:"instancePort"`
	// `CfnLoadBalancer.ListenersProperty.LoadBalancerPort`.
	LoadBalancerPort *string `json:"loadBalancerPort"`
	// `CfnLoadBalancer.ListenersProperty.Protocol`.
	Protocol *string `json:"protocol"`
	// `CfnLoadBalancer.ListenersProperty.InstanceProtocol`.
	InstanceProtocol *string `json:"instanceProtocol"`
	// `CfnLoadBalancer.ListenersProperty.PolicyNames`.
	PolicyNames *[]*string `json:"policyNames"`
	// `CfnLoadBalancer.ListenersProperty.SSLCertificateId`.
	SslCertificateId *string `json:"sslCertificateId"`
}

type CfnLoadBalancer_PoliciesProperty struct {
	// `CfnLoadBalancer.PoliciesProperty.Attributes`.
	Attributes interface{} `json:"attributes"`
	// `CfnLoadBalancer.PoliciesProperty.PolicyName`.
	PolicyName *string `json:"policyName"`
	// `CfnLoadBalancer.PoliciesProperty.PolicyType`.
	PolicyType *string `json:"policyType"`
	// `CfnLoadBalancer.PoliciesProperty.InstancePorts`.
	InstancePorts *[]*string `json:"instancePorts"`
	// `CfnLoadBalancer.PoliciesProperty.LoadBalancerPorts`.
	LoadBalancerPorts *[]*string `json:"loadBalancerPorts"`
}

// Properties for defining a `AWS::ElasticLoadBalancing::LoadBalancer`.
type CfnLoadBalancerProps struct {
	// `AWS::ElasticLoadBalancing::LoadBalancer.Listeners`.
	Listeners interface{} `json:"listeners"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.AccessLoggingPolicy`.
	AccessLoggingPolicy interface{} `json:"accessLoggingPolicy"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.AppCookieStickinessPolicy`.
	AppCookieStickinessPolicy interface{} `json:"appCookieStickinessPolicy"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.AvailabilityZones`.
	AvailabilityZones *[]*string `json:"availabilityZones"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.ConnectionDrainingPolicy`.
	ConnectionDrainingPolicy interface{} `json:"connectionDrainingPolicy"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.ConnectionSettings`.
	ConnectionSettings interface{} `json:"connectionSettings"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.CrossZone`.
	CrossZone interface{} `json:"crossZone"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.HealthCheck`.
	HealthCheck interface{} `json:"healthCheck"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.Instances`.
	Instances *[]*string `json:"instances"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.LBCookieStickinessPolicy`.
	LbCookieStickinessPolicy interface{} `json:"lbCookieStickinessPolicy"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.LoadBalancerName`.
	LoadBalancerName *string `json:"loadBalancerName"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.Policies`.
	Policies interface{} `json:"policies"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.Scheme`.
	Scheme *string `json:"scheme"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.SecurityGroups`.
	SecurityGroups *[]*string `json:"securityGroups"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.Subnets`.
	Subnets *[]*string `json:"subnets"`
	// `AWS::ElasticLoadBalancing::LoadBalancer.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// Describe the health check to a load balancer.
// Experimental.
type HealthCheck struct {
	// What port number to health check on.
	// Experimental.
	Port *float64 `json:"port"`
	// After how many successful checks is an instance considered healthy.
	// Experimental.
	HealthyThreshold *float64 `json:"healthyThreshold"`
	// Number of seconds between health checks.
	// Experimental.
	Interval awscdk.Duration `json:"interval"`
	// What path to use for HTTP or HTTPS health check (must return 200).
	//
	// For SSL and TCP health checks, accepting connections is enough to be considered
	// healthy.
	// Experimental.
	Path *string `json:"path"`
	// What protocol to use for health checking.
	//
	// The protocol is automatically determined from the port if it's not supplied.
	// Experimental.
	Protocol LoadBalancingProtocol `json:"protocol"`
	// Health check timeout.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// After how many unsuccessful checks is an instance considered unhealthy.
	// Experimental.
	UnhealthyThreshold *float64 `json:"unhealthyThreshold"`
}

// Interface that is going to be implemented by constructs that you can load balance to.
// Experimental.
type ILoadBalancerTarget interface {
	awsec2.IConnectable
	// Attach load-balanced target to a classic ELB.
	// Experimental.
	AttachToClassicLB(loadBalancer LoadBalancer)
}

// The jsii proxy for ILoadBalancerTarget
type jsiiProxy_ILoadBalancerTarget struct {
	internal.Type__awsec2IConnectable
}

func (i *jsiiProxy_ILoadBalancerTarget) AttachToClassicLB(loadBalancer LoadBalancer) {
	_jsii_.InvokeVoid(
		i,
		"attachToClassicLB",
		[]interface{}{loadBalancer},
	)
}

// Reference to a listener's port just created.
//
// This implements IConnectable with a default port (the port that an ELB
// listener was just created on) for a given security group so that it can be
// conveniently used just like any Connectable. E.g:
//
//     const listener = elb.addListener(...);
//
//     listener.connections.allowDefaultPortFromAnyIPv4();
//     // or
//     instance.connections.allowToDefaultPort(listener);
// Experimental.
type ListenerPort interface {
	awsec2.IConnectable
	Connections() awsec2.Connections
}

// The jsii proxy struct for ListenerPort
type jsiiProxy_ListenerPort struct {
	internal.Type__awsec2IConnectable
}

func (j *jsiiProxy_ListenerPort) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}


// Experimental.
func NewListenerPort(securityGroup awsec2.ISecurityGroup, defaultPort awsec2.Port) ListenerPort {
	_init_.Initialize()

	j := jsiiProxy_ListenerPort{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancing.ListenerPort",
		[]interface{}{securityGroup, defaultPort},
		&j,
	)

	return &j
}

// Experimental.
func NewListenerPort_Override(l ListenerPort, securityGroup awsec2.ISecurityGroup, defaultPort awsec2.Port) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancing.ListenerPort",
		[]interface{}{securityGroup, defaultPort},
		l,
	)
}

// A load balancer with a single listener.
//
// Routes to a fleet of of instances in a VPC.
// Experimental.
type LoadBalancer interface {
	awscdk.Resource
	awsec2.IConnectable
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	ListenerPorts() *[]ListenerPort
	LoadBalancerCanonicalHostedZoneName() *string
	LoadBalancerCanonicalHostedZoneNameId() *string
	LoadBalancerDnsName() *string
	LoadBalancerName() *string
	LoadBalancerSourceSecurityGroupGroupName() *string
	LoadBalancerSourceSecurityGroupOwnerAlias() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddListener(listener *LoadBalancerListener) ListenerPort
	AddTarget(target ILoadBalancerTarget)
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

// The jsii proxy struct for LoadBalancer
type jsiiProxy_LoadBalancer struct {
	internal.Type__awscdkResource
	internal.Type__awsec2IConnectable
}

func (j *jsiiProxy_LoadBalancer) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) ListenerPorts() *[]ListenerPort {
	var returns *[]ListenerPort
	_jsii_.Get(
		j,
		"listenerPorts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) LoadBalancerCanonicalHostedZoneName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) LoadBalancerCanonicalHostedZoneNameId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneNameId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) LoadBalancerDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) LoadBalancerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) LoadBalancerSourceSecurityGroupGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerSourceSecurityGroupGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) LoadBalancerSourceSecurityGroupOwnerAlias() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerSourceSecurityGroupOwnerAlias",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewLoadBalancer(scope constructs.Construct, id *string, props *LoadBalancerProps) LoadBalancer {
	_init_.Initialize()

	j := jsiiProxy_LoadBalancer{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancing.LoadBalancer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewLoadBalancer_Override(l LoadBalancer, scope constructs.Construct, id *string, props *LoadBalancerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancing.LoadBalancer",
		[]interface{}{scope, id, props},
		l,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func LoadBalancer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancing.LoadBalancer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func LoadBalancer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancing.LoadBalancer",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a backend to the load balancer.
//
// Returns: A ListenerPort object that controls connections to the listener port
// Experimental.
func (l *jsiiProxy_LoadBalancer) AddListener(listener *LoadBalancerListener) ListenerPort {
	var returns ListenerPort

	_jsii_.Invoke(
		l,
		"addListener",
		[]interface{}{listener},
		&returns,
	)

	return returns
}

// Experimental.
func (l *jsiiProxy_LoadBalancer) AddTarget(target ILoadBalancerTarget) {
	_jsii_.InvokeVoid(
		l,
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
func (l *jsiiProxy_LoadBalancer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		l,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (l *jsiiProxy_LoadBalancer) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		l,
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
func (l *jsiiProxy_LoadBalancer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		l,
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
func (l *jsiiProxy_LoadBalancer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		l,
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
func (l *jsiiProxy_LoadBalancer) OnPrepare() {
	_jsii_.InvokeVoid(
		l,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (l *jsiiProxy_LoadBalancer) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
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
func (l *jsiiProxy_LoadBalancer) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
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
func (l *jsiiProxy_LoadBalancer) Prepare() {
	_jsii_.InvokeVoid(
		l,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (l *jsiiProxy_LoadBalancer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (l *jsiiProxy_LoadBalancer) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		l,
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
func (l *jsiiProxy_LoadBalancer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Add a backend to the load balancer.
// Experimental.
type LoadBalancerListener struct {
	// External listening port.
	// Experimental.
	ExternalPort *float64 `json:"externalPort"`
	// Allow connections to the load balancer from the given set of connection peers.
	//
	// By default, connections will be allowed from anywhere. Set this to an empty list
	// to deny connections, or supply a custom list of peers to allow connections from
	// (IP ranges or security groups).
	// Experimental.
	AllowConnectionsFrom *[]awsec2.IConnectable `json:"allowConnectionsFrom"`
	// What public protocol to use for load balancing.
	//
	// Either 'tcp', 'ssl', 'http' or 'https'.
	//
	// May be omitted if the external port is either 80 or 443.
	// Experimental.
	ExternalProtocol LoadBalancingProtocol `json:"externalProtocol"`
	// Instance listening port.
	//
	// Same as the externalPort if not specified.
	// Experimental.
	InternalPort *float64 `json:"internalPort"`
	// What public protocol to use for load balancing.
	//
	// Either 'tcp', 'ssl', 'http' or 'https'.
	//
	// May be omitted if the internal port is either 80 or 443.
	//
	// The instance protocol is 'tcp' if the front-end protocol
	// is 'tcp' or 'ssl', the instance protocol is 'http' if the
	// front-end protocol is 'https'.
	// Experimental.
	InternalProtocol LoadBalancingProtocol `json:"internalProtocol"`
	// SSL policy names.
	// Experimental.
	PolicyNames *[]*string `json:"policyNames"`
	// the ARN of the SSL certificate.
	// Experimental.
	SslCertificateArn *string `json:"sslCertificateArn"`
	// the ARN of the SSL certificate.
	// Deprecated: - use sslCertificateArn instead
	SslCertificateId *string `json:"sslCertificateId"`
}

// Construction properties for a LoadBalancer.
// Experimental.
type LoadBalancerProps struct {
	// VPC network of the fleet instances.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Enable Loadbalancer access logs Can be used to avoid manual work as aws console Required S3 bucket name , enabled flag Can add interval for pushing log Can set bucket prefix in order to provide folder name inside bucket.
	// Experimental.
	AccessLoggingPolicy *CfnLoadBalancer_AccessLoggingPolicyProperty `json:"accessLoggingPolicy"`
	// Whether cross zone load balancing is enabled.
	//
	// This controls whether the load balancer evenly distributes requests
	// across each availability zone
	// Experimental.
	CrossZone *bool `json:"crossZone"`
	// Health check settings for the load balancing targets.
	//
	// Not required but recommended.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// Whether this is an internet-facing Load Balancer.
	//
	// This controls whether the LB has a public IP address assigned. It does
	// not open up the Load Balancer's security groups to public internet access.
	// Experimental.
	InternetFacing *bool `json:"internetFacing"`
	// What listeners to set up for the load balancer.
	//
	// Can also be added by .addListener()
	// Experimental.
	Listeners *[]*LoadBalancerListener `json:"listeners"`
	// Which subnets to deploy the load balancer.
	//
	// Can be used to define a specific set of subnets to deploy the load balancer to.
	// Useful multiple public or private subnets are covering the same availability zone.
	// Experimental.
	SubnetSelection *awsec2.SubnetSelection `json:"subnetSelection"`
	// What targets to load balance to.
	//
	// Can also be added by .addTarget()
	// Experimental.
	Targets *[]ILoadBalancerTarget `json:"targets"`
}

// Experimental.
type LoadBalancingProtocol string

const (
	LoadBalancingProtocol_TCP LoadBalancingProtocol = "TCP"
	LoadBalancingProtocol_SSL LoadBalancingProtocol = "SSL"
	LoadBalancingProtocol_HTTP LoadBalancingProtocol = "HTTP"
	LoadBalancingProtocol_HTTPS LoadBalancingProtocol = "HTTPS"
)


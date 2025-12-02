package awsautoscaling

import (
	"time"

	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsautoscaling/internal"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancing"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awssns"
	"github.com/aws/constructs-go/constructs/v3"
)

// An adjustment.
// Experimental.
type AdjustmentTier struct {
	// What number to adjust the capacity with.
	//
	// The number is interpeted as an added capacity, a new fixed capacity or an
	// added percentage depending on the AdjustmentType value of the
	// StepScalingPolicy.
	//
	// Can be positive or negative.
	// Experimental.
	Adjustment *float64 `json:"adjustment"`
	// Lower bound where this scaling tier applies.
	//
	// The scaling tier applies if the difference between the metric
	// value and its alarm threshold is higher than this value.
	// Experimental.
	LowerBound *float64 `json:"lowerBound"`
	// Upper bound where this scaling tier applies.
	//
	// The scaling tier applies if the difference between the metric
	// value and its alarm threshold is lower than this value.
	// Experimental.
	UpperBound *float64 `json:"upperBound"`
}

// How adjustment numbers are interpreted.
// Experimental.
type AdjustmentType string

const (
	AdjustmentType_CHANGE_IN_CAPACITY AdjustmentType = "CHANGE_IN_CAPACITY"
	AdjustmentType_PERCENT_CHANGE_IN_CAPACITY AdjustmentType = "PERCENT_CHANGE_IN_CAPACITY"
	AdjustmentType_EXACT_CAPACITY AdjustmentType = "EXACT_CAPACITY"
)

// Options for applying CloudFormation init to an instance or instance group.
// Experimental.
type ApplyCloudFormationInitOptions struct {
	// ConfigSet to activate.
	// Experimental.
	ConfigSets *[]*string `json:"configSets"`
	// Force instance replacement by embedding a config fingerprint.
	//
	// If `true` (the default), a hash of the config will be embedded into the
	// UserData, so that if the config changes, the UserData changes and
	// instances will be replaced (given an UpdatePolicy has been configured on
	// the AutoScalingGroup).
	//
	// If `false`, no such hash will be embedded, and if the CloudFormation Init
	// config changes nothing will happen to the running instances. If a
	// config update introduces errors, you will not notice until after the
	// CloudFormation deployment successfully finishes and the next instance
	// fails to launch.
	// Experimental.
	EmbedFingerprint *bool `json:"embedFingerprint"`
	// Don't fail the instance creation when cfn-init fails.
	//
	// You can use this to prevent CloudFormation from rolling back when
	// instances fail to start up, to help in debugging.
	// Experimental.
	IgnoreFailures *bool `json:"ignoreFailures"`
	// Print the results of running cfn-init to the Instance System Log.
	//
	// By default, the output of running cfn-init is written to a log file
	// on the instance. Set this to `true` to print it to the System Log
	// (visible from the EC2 Console), `false` to not print it.
	//
	// (Be aware that the system log is refreshed at certain points in
	// time of the instance life cycle, and successful execution may
	// not always show up).
	// Experimental.
	PrintLog *bool `json:"printLog"`
}

// A Fleet represents a managed set of EC2 instances.
//
// The Fleet models a number of AutoScalingGroups, a launch configuration, a
// security group and an instance role.
//
// It allows adding arbitrary commands to the startup scripts of the instances
// in the fleet.
//
// The ASG spans the availability zones specified by vpcSubnets, falling back to
// the Vpc default strategy if not specified.
// Experimental.
type AutoScalingGroup interface {
	awscdk.Resource
	IAutoScalingGroup
	awsec2.IConnectable
	awselasticloadbalancing.ILoadBalancerTarget
	awselasticloadbalancingv2.IApplicationLoadBalancerTarget
	awselasticloadbalancingv2.INetworkLoadBalancerTarget
	AlbTargetGroup() awselasticloadbalancingv2.ApplicationTargetGroup
	SetAlbTargetGroup(val awselasticloadbalancingv2.ApplicationTargetGroup)
	AutoScalingGroupArn() *string
	AutoScalingGroupName() *string
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() awsiam.IPrincipal
	MaxInstanceLifetime() awscdk.Duration
	NewInstancesProtectedFromScaleIn() *bool
	SetNewInstancesProtectedFromScaleIn(val *bool)
	Node() awscdk.ConstructNode
	OsType() awsec2.OperatingSystemType
	PhysicalName() *string
	Role() awsiam.IRole
	SpotPrice() *string
	Stack() awscdk.Stack
	UserData() awsec2.UserData
	AddLifecycleHook(id *string, props *BasicLifecycleHookProps) LifecycleHook
	AddSecurityGroup(securityGroup awsec2.ISecurityGroup)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	AddUserData(commands ...*string)
	ApplyCloudFormationInit(init awsec2.CloudFormationInit, options *ApplyCloudFormationInitOptions)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AreNewInstancesProtectedFromScaleIn() *bool
	AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer)
	AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	ProtectNewInstancesFromScaleIn()
	ScaleOnCpuUtilization(id *string, props *CpuUtilizationScalingProps) TargetTrackingScalingPolicy
	ScaleOnIncomingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy
	ScaleOnMetric(id *string, props *BasicStepScalingPolicyProps) StepScalingPolicy
	ScaleOnOutgoingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy
	ScaleOnRequestCount(id *string, props *RequestCountScalingProps) TargetTrackingScalingPolicy
	ScaleOnSchedule(id *string, props *BasicScheduledActionProps) ScheduledAction
	ScaleToTrackMetric(id *string, props *MetricTargetTrackingProps) TargetTrackingScalingPolicy
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for AutoScalingGroup
type jsiiProxy_AutoScalingGroup struct {
	internal.Type__awscdkResource
	jsiiProxy_IAutoScalingGroup
	internal.Type__awsec2IConnectable
	internal.Type__awselasticloadbalancingILoadBalancerTarget
	internal.Type__awselasticloadbalancingv2IApplicationLoadBalancerTarget
	internal.Type__awselasticloadbalancingv2INetworkLoadBalancerTarget
}

func (j *jsiiProxy_AutoScalingGroup) AlbTargetGroup() awselasticloadbalancingv2.ApplicationTargetGroup {
	var returns awselasticloadbalancingv2.ApplicationTargetGroup
	_jsii_.Get(
		j,
		"albTargetGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) AutoScalingGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) AutoScalingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) MaxInstanceLifetime() awscdk.Duration {
	var returns awscdk.Duration
	_jsii_.Get(
		j,
		"maxInstanceLifetime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) NewInstancesProtectedFromScaleIn() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"newInstancesProtectedFromScaleIn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) OsType() awsec2.OperatingSystemType {
	var returns awsec2.OperatingSystemType
	_jsii_.Get(
		j,
		"osType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) SpotPrice() *string {
	var returns *string
	_jsii_.Get(
		j,
		"spotPrice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AutoScalingGroup) UserData() awsec2.UserData {
	var returns awsec2.UserData
	_jsii_.Get(
		j,
		"userData",
		&returns,
	)
	return returns
}


// Experimental.
func NewAutoScalingGroup(scope constructs.Construct, id *string, props *AutoScalingGroupProps) AutoScalingGroup {
	_init_.Initialize()

	j := jsiiProxy_AutoScalingGroup{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.AutoScalingGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewAutoScalingGroup_Override(a AutoScalingGroup, scope constructs.Construct, id *string, props *AutoScalingGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.AutoScalingGroup",
		[]interface{}{scope, id, props},
		a,
	)
}

func (j *jsiiProxy_AutoScalingGroup) SetAlbTargetGroup(val awselasticloadbalancingv2.ApplicationTargetGroup) {
	_jsii_.Set(
		j,
		"albTargetGroup",
		val,
	)
}

func (j *jsiiProxy_AutoScalingGroup) SetNewInstancesProtectedFromScaleIn(val *bool) {
	_jsii_.Set(
		j,
		"newInstancesProtectedFromScaleIn",
		val,
	)
}

// Experimental.
func AutoScalingGroup_FromAutoScalingGroupName(scope constructs.Construct, id *string, autoScalingGroupName *string) IAutoScalingGroup {
	_init_.Initialize()

	var returns IAutoScalingGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.AutoScalingGroup",
		"fromAutoScalingGroupName",
		[]interface{}{scope, id, autoScalingGroupName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func AutoScalingGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.AutoScalingGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func AutoScalingGroup_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.AutoScalingGroup",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Send a message to either an SQS queue or SNS topic when instances launch or terminate.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AddLifecycleHook(id *string, props *BasicLifecycleHookProps) LifecycleHook {
	var returns LifecycleHook

	_jsii_.Invoke(
		a,
		"addLifecycleHook",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Add the security group to all instances via the launch configuration security groups array.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AddSecurityGroup(securityGroup awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		a,
		"addSecurityGroup",
		[]interface{}{securityGroup},
	)
}

// Adds a statement to the IAM role assumed by instances of this fleet.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		a,
		"addToRolePolicy",
		[]interface{}{statement},
	)
}

// Add command to the startup script of fleet instances.
//
// The command must be in the scripting language supported by the fleet's OS (i.e. Linux/Windows).
// Does nothing for imported ASGs.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AddUserData(commands ...*string) {
	args := []interface{}{}
	for _, a := range commands {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addUserData",
		args,
	)
}

// Use a CloudFormation Init configuration at instance startup.
//
// This does the following:
//
// - Attaches the CloudFormation Init metadata to the AutoScalingGroup resource.
// - Add commands to the UserData to run `cfn-init` and `cfn-signal`.
// - Update the instance's CreationPolicy to wait for `cfn-init` to finish
//    before reporting success.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ApplyCloudFormationInit(init awsec2.CloudFormationInit, options *ApplyCloudFormationInitOptions) {
	_jsii_.InvokeVoid(
		a,
		"applyCloudFormationInit",
		[]interface{}{init, options},
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
func (a *jsiiProxy_AutoScalingGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Returns `true` if newly-launched instances are protected from scale-in.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AreNewInstancesProtectedFromScaleIn() *bool {
	var returns *bool

	_jsii_.Invoke(
		a,
		"areNewInstancesProtectedFromScaleIn",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Attach to ELBv2 Application Target Group.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		a,
		"attachToApplicationTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// Attach to a classic load balancer.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer) {
	_jsii_.InvokeVoid(
		a,
		"attachToClassicLB",
		[]interface{}{loadBalancer},
	)
}

// Attach to ELBv2 Application Target Group.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		a,
		"attachToNetworkTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// Experimental.
func (a *jsiiProxy_AutoScalingGroup) GeneratePhysicalName() *string {
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
func (a *jsiiProxy_AutoScalingGroup) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (a *jsiiProxy_AutoScalingGroup) GetResourceNameAttribute(nameAttr *string) *string {
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
func (a *jsiiProxy_AutoScalingGroup) OnPrepare() {
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
func (a *jsiiProxy_AutoScalingGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_AutoScalingGroup) OnValidate() *[]*string {
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
func (a *jsiiProxy_AutoScalingGroup) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Ensures newly-launched instances are protected from scale-in.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ProtectNewInstancesFromScaleIn() {
	_jsii_.InvokeVoid(
		a,
		"protectNewInstancesFromScaleIn",
		nil, // no parameters
	)
}

// Scale out or in to achieve a target CPU utilization.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ScaleOnCpuUtilization(id *string, props *CpuUtilizationScalingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		a,
		"scaleOnCpuUtilization",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Scale out or in to achieve a target network ingress rate.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ScaleOnIncomingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		a,
		"scaleOnIncomingBytes",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Scale out or in, in response to a metric.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ScaleOnMetric(id *string, props *BasicStepScalingPolicyProps) StepScalingPolicy {
	var returns StepScalingPolicy

	_jsii_.Invoke(
		a,
		"scaleOnMetric",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Scale out or in to achieve a target network egress rate.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ScaleOnOutgoingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		a,
		"scaleOnOutgoingBytes",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Scale out or in to achieve a target request handling rate.
//
// The AutoScalingGroup must have been attached to an Application Load Balancer
// in order to be able to call this.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ScaleOnRequestCount(id *string, props *RequestCountScalingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		a,
		"scaleOnRequestCount",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Scale out or in based on time.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ScaleOnSchedule(id *string, props *BasicScheduledActionProps) ScheduledAction {
	var returns ScheduledAction

	_jsii_.Invoke(
		a,
		"scaleOnSchedule",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Scale out or in in order to keep a metric around a target value.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ScaleToTrackMetric(id *string, props *MetricTargetTrackingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		a,
		"scaleToTrackMetric",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_AutoScalingGroup) ToString() *string {
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
func (a *jsiiProxy_AutoScalingGroup) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties of a Fleet.
// Experimental.
type AutoScalingGroupProps struct {
	// Whether the instances can initiate connections to anywhere by default.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Whether instances in the Auto Scaling Group should have public IP addresses associated with them.
	// Experimental.
	AssociatePublicIpAddress *bool `json:"associatePublicIpAddress"`
	// The name of the Auto Scaling group.
	//
	// This name must be unique per Region per account.
	// Experimental.
	AutoScalingGroupName *string `json:"autoScalingGroupName"`
	// Specifies how block devices are exposed to the instance. You can specify virtual devices and EBS volumes.
	//
	// Each instance that is launched has an associated root device volume,
	// either an Amazon EBS volume or an instance store volume.
	// You can use block device mappings to specify additional EBS volumes or
	// instance store volumes to attach to an instance when it is launched.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/block-device-mapping-concepts.html
	//
	// Experimental.
	BlockDevices *[]*BlockDevice `json:"blockDevices"`
	// Default scaling cooldown for this AutoScalingGroup.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Initial amount of instances in the fleet.
	//
	// If this is set to a number, every deployment will reset the amount of
	// instances to this number. It is recommended to leave this value blank.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-group.html#cfn-as-group-desiredcapacity
	//
	// Experimental.
	DesiredCapacity *float64 `json:"desiredCapacity"`
	// Enable monitoring for group metrics, these metrics describe the group rather than any of its instances.
	//
	// To report all group metrics use `GroupMetrics.all()`
	// Group metrics are reported in a granularity of 1 minute at no additional charge.
	// Experimental.
	GroupMetrics *[]GroupMetrics `json:"groupMetrics"`
	// Configuration for health checks.
	// Experimental.
	HealthCheck HealthCheck `json:"healthCheck"`
	// If the ASG has scheduled actions, don't reset unchanged group sizes.
	//
	// Only used if the ASG has scheduled actions (which may scale your ASG up
	// or down regardless of cdk deployments). If true, the size of the group
	// will only be reset if it has been changed in the CDK app. If false, the
	// sizes will always be changed back to what they were in the CDK app
	// on deployment.
	// Experimental.
	IgnoreUnmodifiedSizeProperties *bool `json:"ignoreUnmodifiedSizeProperties"`
	// Controls whether instances in this group are launched with detailed or basic monitoring.
	//
	// When detailed monitoring is enabled, Amazon CloudWatch generates metrics every minute and your account
	// is charged a fee. When you disable detailed monitoring, CloudWatch generates metrics every 5 minutes.
	// See: https://docs.aws.amazon.com/autoscaling/latest/userguide/as-instance-monitoring.html#enable-as-instance-metrics
	//
	// Experimental.
	InstanceMonitoring Monitoring `json:"instanceMonitoring"`
	// Name of SSH keypair to grant access to instances.
	// Experimental.
	KeyName *string `json:"keyName"`
	// Maximum number of instances in the fleet.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// The maximum amount of time that an instance can be in service.
	//
	// The maximum duration applies
	// to all current and future instances in the group. As an instance approaches its maximum duration,
	// it is terminated and replaced, and cannot be used again.
	//
	// You must specify a value of at least 604,800 seconds (7 days). To clear a previously set value,
	// leave this property undefined.
	// See: https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-max-instance-lifetime.html
	//
	// Experimental.
	MaxInstanceLifetime awscdk.Duration `json:"maxInstanceLifetime"`
	// Minimum number of instances in the fleet.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// Whether newly-launched instances are protected from termination by Amazon EC2 Auto Scaling when scaling in.
	//
	// By default, Auto Scaling can terminate an instance at any time after launch
	// when scaling in an Auto Scaling Group, subject to the group's termination
	// policy. However, you may wish to protect newly-launched instances from
	// being scaled in if they are going to run critical applications that should
	// not be prematurely terminated.
	//
	// This flag must be enabled if the Auto Scaling Group will be associated with
	// an ECS Capacity Provider with managed termination protection.
	// Experimental.
	NewInstancesProtectedFromScaleIn *bool `json:"newInstancesProtectedFromScaleIn"`
	// Configure autoscaling group to send notifications about fleet changes to an SNS topic(s).
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-group.html#cfn-as-group-notificationconfigurations
	//
	// Experimental.
	Notifications *[]*NotificationConfiguration `json:"notifications"`
	// SNS topic to send notifications about fleet changes.
	// Deprecated: use `notifications`
	NotificationsTopic awssns.ITopic `json:"notificationsTopic"`
	// Configuration for replacing updates.
	//
	// Only used if updateType == UpdateType.ReplacingUpdate. Specifies how
	// many instances must signal success for the update to succeed.
	// Deprecated: Use `signals` instead
	ReplacingUpdateMinSuccessfulInstancesPercent *float64 `json:"replacingUpdateMinSuccessfulInstancesPercent"`
	// How many ResourceSignal calls CloudFormation expects before the resource is considered created.
	// Deprecated: Use `signals` instead.
	ResourceSignalCount *float64 `json:"resourceSignalCount"`
	// The length of time to wait for the resourceSignalCount.
	//
	// The maximum value is 43200 (12 hours).
	// Deprecated: Use `signals` instead.
	ResourceSignalTimeout awscdk.Duration `json:"resourceSignalTimeout"`
	// Configuration for rolling updates.
	//
	// Only used if updateType == UpdateType.RollingUpdate.
	// Deprecated: Use `updatePolicy` instead
	RollingUpdateConfiguration *RollingUpdateConfiguration `json:"rollingUpdateConfiguration"`
	// Configure waiting for signals during deployment.
	//
	// Use this to pause the CloudFormation deployment to wait for the instances
	// in the AutoScalingGroup to report successful startup during
	// creation and updates. The UserData script needs to invoke `cfn-signal`
	// with a success or failure code after it is done setting up the instance.
	//
	// Without waiting for signals, the CloudFormation deployment will proceed as
	// soon as the AutoScalingGroup has been created or updated but before the
	// instances in the group have been started.
	//
	// For example, to have instances wait for an Elastic Load Balancing health check before
	// they signal success, add a health-check verification by using the
	// cfn-init helper script. For an example, see the verify_instance_health
	// command in the Auto Scaling rolling updates sample template:
	//
	// https://github.com/awslabs/aws-cloudformation-templates/blob/master/aws/services/AutoScaling/AutoScalingRollingUpdates.yaml
	// Experimental.
	Signals Signals `json:"signals"`
	// The maximum hourly price (in USD) to be paid for any Spot Instance launched to fulfill the request.
	//
	// Spot Instances are
	// launched when the price you specify exceeds the current Spot market price.
	// Experimental.
	SpotPrice *string `json:"spotPrice"`
	// What to do when an AutoScalingGroup's instance configuration is changed.
	//
	// This is applied when any of the settings on the ASG are changed that
	// affect how the instances should be created (VPC, instance type, startup
	// scripts, etc.). It indicates how the existing instances should be
	// replaced with new instances matching the new config. By default, nothing
	// is done and only new instances are launched with the new config.
	// Experimental.
	UpdatePolicy UpdatePolicy `json:"updatePolicy"`
	// What to do when an AutoScalingGroup's instance configuration is changed.
	//
	// This is applied when any of the settings on the ASG are changed that
	// affect how the instances should be created (VPC, instance type, startup
	// scripts, etc.). It indicates how the existing instances should be
	// replaced with new instances matching the new config. By default, nothing
	// is done and only new instances are launched with the new config.
	// Deprecated: Use `updatePolicy` instead
	UpdateType UpdateType `json:"updateType"`
	// Where to place instances within the VPC.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// Type of instance to launch.
	// Experimental.
	InstanceType awsec2.InstanceType `json:"instanceType"`
	// AMI to launch.
	// Experimental.
	MachineImage awsec2.IMachineImage `json:"machineImage"`
	// VPC to launch these instances in.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Apply the given CloudFormation Init configuration to the instances in the AutoScalingGroup at startup.
	//
	// If you specify `init`, you must also specify `signals` to configure
	// the number of instances to wait for and the timeout for waiting for the
	// init process.
	// Experimental.
	Init awsec2.CloudFormationInit `json:"init"`
	// Use the given options for applying CloudFormation Init.
	//
	// Describes the configsets to use and the timeout to wait
	// Experimental.
	InitOptions *ApplyCloudFormationInitOptions `json:"initOptions"`
	// An IAM role to associate with the instance profile assigned to this Auto Scaling Group.
	//
	// The role must be assumable by the service principal `ec2.amazonaws.com`:
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// Security group to launch the instances in.
	// Experimental.
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// Specific UserData to use.
	//
	// The UserData may still be mutated after creation.
	// Experimental.
	UserData awsec2.UserData `json:"userData"`
}

// Base interface for target tracking props.
//
// Contains the attributes that are common to target tracking policies,
// except the ones relating to the metric and to the scalable target.
//
// This interface is reused by more specific target tracking props objects.
// Experimental.
type BaseTargetTrackingProps struct {
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the autoscaling group. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// group.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
}

// Basic properties for a lifecycle hook.
// Experimental.
type BasicLifecycleHookProps struct {
	// The state of the Amazon EC2 instance to which you want to attach the lifecycle hook.
	// Experimental.
	LifecycleTransition LifecycleTransition `json:"lifecycleTransition"`
	// The target of the lifecycle hook.
	// Experimental.
	NotificationTarget ILifecycleHookTarget `json:"notificationTarget"`
	// The action the Auto Scaling group takes when the lifecycle hook timeout elapses or if an unexpected failure occurs.
	// Experimental.
	DefaultResult DefaultResult `json:"defaultResult"`
	// Maximum time between calls to RecordLifecycleActionHeartbeat for the hook.
	//
	// If the lifecycle hook times out, perform the action in DefaultResult.
	// Experimental.
	HeartbeatTimeout awscdk.Duration `json:"heartbeatTimeout"`
	// Name of the lifecycle hook.
	// Experimental.
	LifecycleHookName *string `json:"lifecycleHookName"`
	// Additional data to pass to the lifecycle hook target.
	// Experimental.
	NotificationMetadata *string `json:"notificationMetadata"`
	// The role that allows publishing to the notification target.
	// Experimental.
	Role awsiam.IRole `json:"role"`
}

// Properties for a scheduled scaling action.
// Experimental.
type BasicScheduledActionProps struct {
	// When to perform this action.
	//
	// Supports cron expressions.
	//
	// For more information about cron expressions, see https://en.wikipedia.org/wiki/Cron.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Schedule Schedule `json:"schedule"`
	// The new desired capacity.
	//
	// At the scheduled time, set the desired capacity to the given capacity.
	//
	// At least one of maxCapacity, minCapacity, or desiredCapacity must be supplied.
	// Experimental.
	DesiredCapacity *float64 `json:"desiredCapacity"`
	// When this scheduled action expires.
	// Experimental.
	EndTime *time.Time `json:"endTime"`
	// The new maximum capacity.
	//
	// At the scheduled time, set the maximum capacity to the given capacity.
	//
	// At least one of maxCapacity, minCapacity, or desiredCapacity must be supplied.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// The new minimum capacity.
	//
	// At the scheduled time, set the minimum capacity to the given capacity.
	//
	// At least one of maxCapacity, minCapacity, or desiredCapacity must be supplied.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// When this scheduled action becomes active.
	// Experimental.
	StartTime *time.Time `json:"startTime"`
}

// Experimental.
type BasicStepScalingPolicyProps struct {
	// Metric to scale on.
	// Experimental.
	Metric awscloudwatch.IMetric `json:"metric"`
	// The intervals for scaling.
	//
	// Maps a range of metric values to a particular scaling behavior.
	// Experimental.
	ScalingSteps *[]*ScalingInterval `json:"scalingSteps"`
	// How the adjustment numbers inside 'intervals' are interpreted.
	// Experimental.
	AdjustmentType AdjustmentType `json:"adjustmentType"`
	// Grace period after scaling activity.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// How many evaluation periods of the metric to wait before triggering a scaling action.
	//
	// Raising this value can be used to smooth out the metric, at the expense
	// of slower response times.
	// Experimental.
	EvaluationPeriods *float64 `json:"evaluationPeriods"`
	// Aggregation to apply to all data points over the evaluation periods.
	//
	// Only has meaning if `evaluationPeriods != 1`.
	// Experimental.
	MetricAggregationType MetricAggregationType `json:"metricAggregationType"`
	// Minimum absolute number to adjust capacity with as result of percentage scaling.
	//
	// Only when using AdjustmentType = PercentChangeInCapacity, this number controls
	// the minimum absolute effect size.
	// Experimental.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
}

// Properties for a Target Tracking policy that include the metric but exclude the target.
// Experimental.
type BasicTargetTrackingScalingPolicyProps struct {
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the autoscaling group. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// group.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// The target value for the metric.
	// Experimental.
	TargetValue *float64 `json:"targetValue"`
	// A custom metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	CustomMetric awscloudwatch.IMetric `json:"customMetric"`
	// A predefined metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	PredefinedMetric PredefinedMetric `json:"predefinedMetric"`
	// The resource label associated with the predefined metric.
	//
	// Should be supplied if the predefined metric is ALBRequestCountPerTarget, and the
	// format should be:
	//
	// app/<load-balancer-name>/<load-balancer-id>/targetgroup/<target-group-name>/<target-group-id>
	// Experimental.
	ResourceLabel *string `json:"resourceLabel"`
}

// Block device.
// Experimental.
type BlockDevice struct {
	// The device name exposed to the EC2 instance.
	//
	// TODO: EXAMPLE
	//
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/device_naming.html
	//
	// Experimental.
	DeviceName *string `json:"deviceName"`
	// Defines the block device volume, to be either an Amazon EBS volume or an ephemeral instance store volume.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Volume BlockDeviceVolume `json:"volume"`
	// If false, the device mapping will be suppressed.
	//
	// If set to false for the root device, the instance might fail the Amazon EC2 health check.
	// Amazon EC2 Auto Scaling launches a replacement instance if the instance fails the health check.
	// Deprecated: use `BlockDeviceVolume.noDevice()` as the volume to supress a mapping.
	MappingEnabled *bool `json:"mappingEnabled"`
}

// Describes a block device mapping for an EC2 instance or Auto Scaling group.
// Experimental.
type BlockDeviceVolume interface {
	EbsDevice() *EbsDeviceProps
	VirtualName() *string
}

// The jsii proxy struct for BlockDeviceVolume
type jsiiProxy_BlockDeviceVolume struct {
	_ byte // padding
}

func (j *jsiiProxy_BlockDeviceVolume) EbsDevice() *EbsDeviceProps {
	var returns *EbsDeviceProps
	_jsii_.Get(
		j,
		"ebsDevice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BlockDeviceVolume) VirtualName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"virtualName",
		&returns,
	)
	return returns
}


// Experimental.
func NewBlockDeviceVolume(ebsDevice *EbsDeviceProps, virtualName *string) BlockDeviceVolume {
	_init_.Initialize()

	j := jsiiProxy_BlockDeviceVolume{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.BlockDeviceVolume",
		[]interface{}{ebsDevice, virtualName},
		&j,
	)

	return &j
}

// Experimental.
func NewBlockDeviceVolume_Override(b BlockDeviceVolume, ebsDevice *EbsDeviceProps, virtualName *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.BlockDeviceVolume",
		[]interface{}{ebsDevice, virtualName},
		b,
	)
}

// Creates a new Elastic Block Storage device.
// Experimental.
func BlockDeviceVolume_Ebs(volumeSize *float64, options *EbsDeviceOptions) BlockDeviceVolume {
	_init_.Initialize()

	var returns BlockDeviceVolume

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.BlockDeviceVolume",
		"ebs",
		[]interface{}{volumeSize, options},
		&returns,
	)

	return returns
}

// Creates a new Elastic Block Storage device from an existing snapshot.
// Experimental.
func BlockDeviceVolume_EbsFromSnapshot(snapshotId *string, options *EbsDeviceSnapshotOptions) BlockDeviceVolume {
	_init_.Initialize()

	var returns BlockDeviceVolume

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.BlockDeviceVolume",
		"ebsFromSnapshot",
		[]interface{}{snapshotId, options},
		&returns,
	)

	return returns
}

// Creates a virtual, ephemeral device.
//
// The name will be in the form ephemeral{volumeIndex}.
// Experimental.
func BlockDeviceVolume_Ephemeral(volumeIndex *float64) BlockDeviceVolume {
	_init_.Initialize()

	var returns BlockDeviceVolume

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.BlockDeviceVolume",
		"ephemeral",
		[]interface{}{volumeIndex},
		&returns,
	)

	return returns
}

// Supresses a volume mapping.
// Experimental.
func BlockDeviceVolume_NoDevice() BlockDeviceVolume {
	_init_.Initialize()

	var returns BlockDeviceVolume

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.BlockDeviceVolume",
		"noDevice",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A CloudFormation `AWS::AutoScaling::AutoScalingGroup`.
type CfnAutoScalingGroup interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AutoScalingGroupName() *string
	SetAutoScalingGroupName(val *string)
	AvailabilityZones() *[]*string
	SetAvailabilityZones(val *[]*string)
	CapacityRebalance() interface{}
	SetCapacityRebalance(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Context() *string
	SetContext(val *string)
	Cooldown() *string
	SetCooldown(val *string)
	CreationStack() *[]*string
	DesiredCapacity() *string
	SetDesiredCapacity(val *string)
	HealthCheckGracePeriod() *float64
	SetHealthCheckGracePeriod(val *float64)
	HealthCheckType() *string
	SetHealthCheckType(val *string)
	InstanceId() *string
	SetInstanceId(val *string)
	LaunchConfigurationName() *string
	SetLaunchConfigurationName(val *string)
	LaunchTemplate() interface{}
	SetLaunchTemplate(val interface{})
	LifecycleHookSpecificationList() interface{}
	SetLifecycleHookSpecificationList(val interface{})
	LoadBalancerNames() *[]*string
	SetLoadBalancerNames(val *[]*string)
	LogicalId() *string
	MaxInstanceLifetime() *float64
	SetMaxInstanceLifetime(val *float64)
	MaxSize() *string
	SetMaxSize(val *string)
	MetricsCollection() interface{}
	SetMetricsCollection(val interface{})
	MinSize() *string
	SetMinSize(val *string)
	MixedInstancesPolicy() interface{}
	SetMixedInstancesPolicy(val interface{})
	NewInstancesProtectedFromScaleIn() interface{}
	SetNewInstancesProtectedFromScaleIn(val interface{})
	Node() awscdk.ConstructNode
	NotificationConfigurations() interface{}
	SetNotificationConfigurations(val interface{})
	PlacementGroup() *string
	SetPlacementGroup(val *string)
	Ref() *string
	ServiceLinkedRoleArn() *string
	SetServiceLinkedRoleArn(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	TargetGroupArns() *[]*string
	SetTargetGroupArns(val *[]*string)
	TerminationPolicies() *[]*string
	SetTerminationPolicies(val *[]*string)
	UpdatedProperites() *map[string]interface{}
	VpcZoneIdentifier() *[]*string
	SetVpcZoneIdentifier(val *[]*string)
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

// The jsii proxy struct for CfnAutoScalingGroup
type jsiiProxy_CfnAutoScalingGroup struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAutoScalingGroup) AutoScalingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) AvailabilityZones() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"availabilityZones",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) CapacityRebalance() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"capacityRebalance",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) Context() *string {
	var returns *string
	_jsii_.Get(
		j,
		"context",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) Cooldown() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cooldown",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) DesiredCapacity() *string {
	var returns *string
	_jsii_.Get(
		j,
		"desiredCapacity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) HealthCheckGracePeriod() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"healthCheckGracePeriod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) HealthCheckType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"healthCheckType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) InstanceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"instanceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) LaunchConfigurationName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"launchConfigurationName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) LaunchTemplate() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"launchTemplate",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) LifecycleHookSpecificationList() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"lifecycleHookSpecificationList",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) LoadBalancerNames() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"loadBalancerNames",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) MaxInstanceLifetime() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxInstanceLifetime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) MaxSize() *string {
	var returns *string
	_jsii_.Get(
		j,
		"maxSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) MetricsCollection() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"metricsCollection",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) MinSize() *string {
	var returns *string
	_jsii_.Get(
		j,
		"minSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) MixedInstancesPolicy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"mixedInstancesPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) NewInstancesProtectedFromScaleIn() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"newInstancesProtectedFromScaleIn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) NotificationConfigurations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"notificationConfigurations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) PlacementGroup() *string {
	var returns *string
	_jsii_.Get(
		j,
		"placementGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) ServiceLinkedRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceLinkedRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) TargetGroupArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"targetGroupArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) TerminationPolicies() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"terminationPolicies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAutoScalingGroup) VpcZoneIdentifier() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"vpcZoneIdentifier",
		&returns,
	)
	return returns
}


// Create a new `AWS::AutoScaling::AutoScalingGroup`.
func NewCfnAutoScalingGroup(scope awscdk.Construct, id *string, props *CfnAutoScalingGroupProps) CfnAutoScalingGroup {
	_init_.Initialize()

	j := jsiiProxy_CfnAutoScalingGroup{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnAutoScalingGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::AutoScaling::AutoScalingGroup`.
func NewCfnAutoScalingGroup_Override(c CfnAutoScalingGroup, scope awscdk.Construct, id *string, props *CfnAutoScalingGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnAutoScalingGroup",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetAutoScalingGroupName(val *string) {
	_jsii_.Set(
		j,
		"autoScalingGroupName",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetAvailabilityZones(val *[]*string) {
	_jsii_.Set(
		j,
		"availabilityZones",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetCapacityRebalance(val interface{}) {
	_jsii_.Set(
		j,
		"capacityRebalance",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetContext(val *string) {
	_jsii_.Set(
		j,
		"context",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetCooldown(val *string) {
	_jsii_.Set(
		j,
		"cooldown",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetDesiredCapacity(val *string) {
	_jsii_.Set(
		j,
		"desiredCapacity",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetHealthCheckGracePeriod(val *float64) {
	_jsii_.Set(
		j,
		"healthCheckGracePeriod",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetHealthCheckType(val *string) {
	_jsii_.Set(
		j,
		"healthCheckType",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetInstanceId(val *string) {
	_jsii_.Set(
		j,
		"instanceId",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetLaunchConfigurationName(val *string) {
	_jsii_.Set(
		j,
		"launchConfigurationName",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetLaunchTemplate(val interface{}) {
	_jsii_.Set(
		j,
		"launchTemplate",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetLifecycleHookSpecificationList(val interface{}) {
	_jsii_.Set(
		j,
		"lifecycleHookSpecificationList",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetLoadBalancerNames(val *[]*string) {
	_jsii_.Set(
		j,
		"loadBalancerNames",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetMaxInstanceLifetime(val *float64) {
	_jsii_.Set(
		j,
		"maxInstanceLifetime",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetMaxSize(val *string) {
	_jsii_.Set(
		j,
		"maxSize",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetMetricsCollection(val interface{}) {
	_jsii_.Set(
		j,
		"metricsCollection",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetMinSize(val *string) {
	_jsii_.Set(
		j,
		"minSize",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetMixedInstancesPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"mixedInstancesPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetNewInstancesProtectedFromScaleIn(val interface{}) {
	_jsii_.Set(
		j,
		"newInstancesProtectedFromScaleIn",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetNotificationConfigurations(val interface{}) {
	_jsii_.Set(
		j,
		"notificationConfigurations",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetPlacementGroup(val *string) {
	_jsii_.Set(
		j,
		"placementGroup",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetServiceLinkedRoleArn(val *string) {
	_jsii_.Set(
		j,
		"serviceLinkedRoleArn",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetTargetGroupArns(val *[]*string) {
	_jsii_.Set(
		j,
		"targetGroupArns",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetTerminationPolicies(val *[]*string) {
	_jsii_.Set(
		j,
		"terminationPolicies",
		val,
	)
}

func (j *jsiiProxy_CfnAutoScalingGroup) SetVpcZoneIdentifier(val *[]*string) {
	_jsii_.Set(
		j,
		"vpcZoneIdentifier",
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
func CfnAutoScalingGroup_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnAutoScalingGroup",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAutoScalingGroup_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnAutoScalingGroup",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAutoScalingGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnAutoScalingGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAutoScalingGroup_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.CfnAutoScalingGroup",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAutoScalingGroup) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAutoScalingGroup) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAutoScalingGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAutoScalingGroup) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAutoScalingGroup) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) OnPrepare() {
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
func (c *jsiiProxy_CfnAutoScalingGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAutoScalingGroup) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAutoScalingGroup) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAutoScalingGroup) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAutoScalingGroup) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAutoScalingGroup) ToString() *string {
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
func (c *jsiiProxy_CfnAutoScalingGroup) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAutoScalingGroup) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnAutoScalingGroup_InstancesDistributionProperty struct {
	// `CfnAutoScalingGroup.InstancesDistributionProperty.OnDemandAllocationStrategy`.
	OnDemandAllocationStrategy *string `json:"onDemandAllocationStrategy"`
	// `CfnAutoScalingGroup.InstancesDistributionProperty.OnDemandBaseCapacity`.
	OnDemandBaseCapacity *float64 `json:"onDemandBaseCapacity"`
	// `CfnAutoScalingGroup.InstancesDistributionProperty.OnDemandPercentageAboveBaseCapacity`.
	OnDemandPercentageAboveBaseCapacity *float64 `json:"onDemandPercentageAboveBaseCapacity"`
	// `CfnAutoScalingGroup.InstancesDistributionProperty.SpotAllocationStrategy`.
	SpotAllocationStrategy *string `json:"spotAllocationStrategy"`
	// `CfnAutoScalingGroup.InstancesDistributionProperty.SpotInstancePools`.
	SpotInstancePools *float64 `json:"spotInstancePools"`
	// `CfnAutoScalingGroup.InstancesDistributionProperty.SpotMaxPrice`.
	SpotMaxPrice *string `json:"spotMaxPrice"`
}

type CfnAutoScalingGroup_LaunchTemplateOverridesProperty struct {
	// `CfnAutoScalingGroup.LaunchTemplateOverridesProperty.InstanceType`.
	InstanceType *string `json:"instanceType"`
	// `CfnAutoScalingGroup.LaunchTemplateOverridesProperty.LaunchTemplateSpecification`.
	LaunchTemplateSpecification interface{} `json:"launchTemplateSpecification"`
	// `CfnAutoScalingGroup.LaunchTemplateOverridesProperty.WeightedCapacity`.
	WeightedCapacity *string `json:"weightedCapacity"`
}

type CfnAutoScalingGroup_LaunchTemplateProperty struct {
	// `CfnAutoScalingGroup.LaunchTemplateProperty.LaunchTemplateSpecification`.
	LaunchTemplateSpecification interface{} `json:"launchTemplateSpecification"`
	// `CfnAutoScalingGroup.LaunchTemplateProperty.Overrides`.
	Overrides interface{} `json:"overrides"`
}

type CfnAutoScalingGroup_LaunchTemplateSpecificationProperty struct {
	// `CfnAutoScalingGroup.LaunchTemplateSpecificationProperty.Version`.
	Version *string `json:"version"`
	// `CfnAutoScalingGroup.LaunchTemplateSpecificationProperty.LaunchTemplateId`.
	LaunchTemplateId *string `json:"launchTemplateId"`
	// `CfnAutoScalingGroup.LaunchTemplateSpecificationProperty.LaunchTemplateName`.
	LaunchTemplateName *string `json:"launchTemplateName"`
}

type CfnAutoScalingGroup_LifecycleHookSpecificationProperty struct {
	// `CfnAutoScalingGroup.LifecycleHookSpecificationProperty.LifecycleHookName`.
	LifecycleHookName *string `json:"lifecycleHookName"`
	// `CfnAutoScalingGroup.LifecycleHookSpecificationProperty.LifecycleTransition`.
	LifecycleTransition *string `json:"lifecycleTransition"`
	// `CfnAutoScalingGroup.LifecycleHookSpecificationProperty.DefaultResult`.
	DefaultResult *string `json:"defaultResult"`
	// `CfnAutoScalingGroup.LifecycleHookSpecificationProperty.HeartbeatTimeout`.
	HeartbeatTimeout *float64 `json:"heartbeatTimeout"`
	// `CfnAutoScalingGroup.LifecycleHookSpecificationProperty.NotificationMetadata`.
	NotificationMetadata *string `json:"notificationMetadata"`
	// `CfnAutoScalingGroup.LifecycleHookSpecificationProperty.NotificationTargetARN`.
	NotificationTargetArn *string `json:"notificationTargetArn"`
	// `CfnAutoScalingGroup.LifecycleHookSpecificationProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
}

type CfnAutoScalingGroup_MetricsCollectionProperty struct {
	// `CfnAutoScalingGroup.MetricsCollectionProperty.Granularity`.
	Granularity *string `json:"granularity"`
	// `CfnAutoScalingGroup.MetricsCollectionProperty.Metrics`.
	Metrics *[]*string `json:"metrics"`
}

type CfnAutoScalingGroup_MixedInstancesPolicyProperty struct {
	// `CfnAutoScalingGroup.MixedInstancesPolicyProperty.LaunchTemplate`.
	LaunchTemplate interface{} `json:"launchTemplate"`
	// `CfnAutoScalingGroup.MixedInstancesPolicyProperty.InstancesDistribution`.
	InstancesDistribution interface{} `json:"instancesDistribution"`
}

type CfnAutoScalingGroup_NotificationConfigurationProperty struct {
	// `CfnAutoScalingGroup.NotificationConfigurationProperty.TopicARN`.
	TopicArn *string `json:"topicArn"`
	// `CfnAutoScalingGroup.NotificationConfigurationProperty.NotificationTypes`.
	NotificationTypes *[]*string `json:"notificationTypes"`
}

type CfnAutoScalingGroup_TagPropertyProperty struct {
	// `CfnAutoScalingGroup.TagPropertyProperty.Key`.
	Key *string `json:"key"`
	// `CfnAutoScalingGroup.TagPropertyProperty.PropagateAtLaunch`.
	PropagateAtLaunch interface{} `json:"propagateAtLaunch"`
	// `CfnAutoScalingGroup.TagPropertyProperty.Value`.
	Value *string `json:"value"`
}

// Properties for defining a `AWS::AutoScaling::AutoScalingGroup`.
type CfnAutoScalingGroupProps struct {
	// `AWS::AutoScaling::AutoScalingGroup.MaxSize`.
	MaxSize *string `json:"maxSize"`
	// `AWS::AutoScaling::AutoScalingGroup.MinSize`.
	MinSize *string `json:"minSize"`
	// `AWS::AutoScaling::AutoScalingGroup.AutoScalingGroupName`.
	AutoScalingGroupName *string `json:"autoScalingGroupName"`
	// `AWS::AutoScaling::AutoScalingGroup.AvailabilityZones`.
	AvailabilityZones *[]*string `json:"availabilityZones"`
	// `AWS::AutoScaling::AutoScalingGroup.CapacityRebalance`.
	CapacityRebalance interface{} `json:"capacityRebalance"`
	// `AWS::AutoScaling::AutoScalingGroup.Context`.
	Context *string `json:"context"`
	// `AWS::AutoScaling::AutoScalingGroup.Cooldown`.
	Cooldown *string `json:"cooldown"`
	// `AWS::AutoScaling::AutoScalingGroup.DesiredCapacity`.
	DesiredCapacity *string `json:"desiredCapacity"`
	// `AWS::AutoScaling::AutoScalingGroup.HealthCheckGracePeriod`.
	HealthCheckGracePeriod *float64 `json:"healthCheckGracePeriod"`
	// `AWS::AutoScaling::AutoScalingGroup.HealthCheckType`.
	HealthCheckType *string `json:"healthCheckType"`
	// `AWS::AutoScaling::AutoScalingGroup.InstanceId`.
	InstanceId *string `json:"instanceId"`
	// `AWS::AutoScaling::AutoScalingGroup.LaunchConfigurationName`.
	LaunchConfigurationName *string `json:"launchConfigurationName"`
	// `AWS::AutoScaling::AutoScalingGroup.LaunchTemplate`.
	LaunchTemplate interface{} `json:"launchTemplate"`
	// `AWS::AutoScaling::AutoScalingGroup.LifecycleHookSpecificationList`.
	LifecycleHookSpecificationList interface{} `json:"lifecycleHookSpecificationList"`
	// `AWS::AutoScaling::AutoScalingGroup.LoadBalancerNames`.
	LoadBalancerNames *[]*string `json:"loadBalancerNames"`
	// `AWS::AutoScaling::AutoScalingGroup.MaxInstanceLifetime`.
	MaxInstanceLifetime *float64 `json:"maxInstanceLifetime"`
	// `AWS::AutoScaling::AutoScalingGroup.MetricsCollection`.
	MetricsCollection interface{} `json:"metricsCollection"`
	// `AWS::AutoScaling::AutoScalingGroup.MixedInstancesPolicy`.
	MixedInstancesPolicy interface{} `json:"mixedInstancesPolicy"`
	// `AWS::AutoScaling::AutoScalingGroup.NewInstancesProtectedFromScaleIn`.
	NewInstancesProtectedFromScaleIn interface{} `json:"newInstancesProtectedFromScaleIn"`
	// `AWS::AutoScaling::AutoScalingGroup.NotificationConfigurations`.
	NotificationConfigurations interface{} `json:"notificationConfigurations"`
	// `AWS::AutoScaling::AutoScalingGroup.PlacementGroup`.
	PlacementGroup *string `json:"placementGroup"`
	// `AWS::AutoScaling::AutoScalingGroup.ServiceLinkedRoleARN`.
	ServiceLinkedRoleArn *string `json:"serviceLinkedRoleArn"`
	// `AWS::AutoScaling::AutoScalingGroup.Tags`.
	Tags *[]*CfnAutoScalingGroup_TagPropertyProperty `json:"tags"`
	// `AWS::AutoScaling::AutoScalingGroup.TargetGroupARNs`.
	TargetGroupArns *[]*string `json:"targetGroupArns"`
	// `AWS::AutoScaling::AutoScalingGroup.TerminationPolicies`.
	TerminationPolicies *[]*string `json:"terminationPolicies"`
	// `AWS::AutoScaling::AutoScalingGroup.VPCZoneIdentifier`.
	VpcZoneIdentifier *[]*string `json:"vpcZoneIdentifier"`
}

// A CloudFormation `AWS::AutoScaling::LaunchConfiguration`.
type CfnLaunchConfiguration interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AssociatePublicIpAddress() interface{}
	SetAssociatePublicIpAddress(val interface{})
	BlockDeviceMappings() interface{}
	SetBlockDeviceMappings(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClassicLinkVpcId() *string
	SetClassicLinkVpcId(val *string)
	ClassicLinkVpcSecurityGroups() *[]*string
	SetClassicLinkVpcSecurityGroups(val *[]*string)
	CreationStack() *[]*string
	EbsOptimized() interface{}
	SetEbsOptimized(val interface{})
	IamInstanceProfile() *string
	SetIamInstanceProfile(val *string)
	ImageId() *string
	SetImageId(val *string)
	InstanceId() *string
	SetInstanceId(val *string)
	InstanceMonitoring() interface{}
	SetInstanceMonitoring(val interface{})
	InstanceType() *string
	SetInstanceType(val *string)
	KernelId() *string
	SetKernelId(val *string)
	KeyName() *string
	SetKeyName(val *string)
	LaunchConfigurationName() *string
	SetLaunchConfigurationName(val *string)
	LogicalId() *string
	MetadataOptions() interface{}
	SetMetadataOptions(val interface{})
	Node() awscdk.ConstructNode
	PlacementTenancy() *string
	SetPlacementTenancy(val *string)
	RamDiskId() *string
	SetRamDiskId(val *string)
	Ref() *string
	SecurityGroups() *[]*string
	SetSecurityGroups(val *[]*string)
	SpotPrice() *string
	SetSpotPrice(val *string)
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserData() *string
	SetUserData(val *string)
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

// The jsii proxy struct for CfnLaunchConfiguration
type jsiiProxy_CfnLaunchConfiguration struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnLaunchConfiguration) AssociatePublicIpAddress() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"associatePublicIpAddress",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) BlockDeviceMappings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"blockDeviceMappings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) ClassicLinkVpcId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"classicLinkVpcId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) ClassicLinkVpcSecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"classicLinkVpcSecurityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) EbsOptimized() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ebsOptimized",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) IamInstanceProfile() *string {
	var returns *string
	_jsii_.Get(
		j,
		"iamInstanceProfile",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) ImageId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) InstanceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"instanceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) InstanceMonitoring() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"instanceMonitoring",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) InstanceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"instanceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) KernelId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kernelId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) KeyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) LaunchConfigurationName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"launchConfigurationName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) MetadataOptions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"metadataOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) PlacementTenancy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"placementTenancy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) RamDiskId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ramDiskId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) SecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) SpotPrice() *string {
	var returns *string
	_jsii_.Get(
		j,
		"spotPrice",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLaunchConfiguration) UserData() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userData",
		&returns,
	)
	return returns
}


// Create a new `AWS::AutoScaling::LaunchConfiguration`.
func NewCfnLaunchConfiguration(scope awscdk.Construct, id *string, props *CfnLaunchConfigurationProps) CfnLaunchConfiguration {
	_init_.Initialize()

	j := jsiiProxy_CfnLaunchConfiguration{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnLaunchConfiguration",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::AutoScaling::LaunchConfiguration`.
func NewCfnLaunchConfiguration_Override(c CfnLaunchConfiguration, scope awscdk.Construct, id *string, props *CfnLaunchConfigurationProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnLaunchConfiguration",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetAssociatePublicIpAddress(val interface{}) {
	_jsii_.Set(
		j,
		"associatePublicIpAddress",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetBlockDeviceMappings(val interface{}) {
	_jsii_.Set(
		j,
		"blockDeviceMappings",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetClassicLinkVpcId(val *string) {
	_jsii_.Set(
		j,
		"classicLinkVpcId",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetClassicLinkVpcSecurityGroups(val *[]*string) {
	_jsii_.Set(
		j,
		"classicLinkVpcSecurityGroups",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetEbsOptimized(val interface{}) {
	_jsii_.Set(
		j,
		"ebsOptimized",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetIamInstanceProfile(val *string) {
	_jsii_.Set(
		j,
		"iamInstanceProfile",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetImageId(val *string) {
	_jsii_.Set(
		j,
		"imageId",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetInstanceId(val *string) {
	_jsii_.Set(
		j,
		"instanceId",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetInstanceMonitoring(val interface{}) {
	_jsii_.Set(
		j,
		"instanceMonitoring",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetInstanceType(val *string) {
	_jsii_.Set(
		j,
		"instanceType",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetKernelId(val *string) {
	_jsii_.Set(
		j,
		"kernelId",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetKeyName(val *string) {
	_jsii_.Set(
		j,
		"keyName",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetLaunchConfigurationName(val *string) {
	_jsii_.Set(
		j,
		"launchConfigurationName",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetMetadataOptions(val interface{}) {
	_jsii_.Set(
		j,
		"metadataOptions",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetPlacementTenancy(val *string) {
	_jsii_.Set(
		j,
		"placementTenancy",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetRamDiskId(val *string) {
	_jsii_.Set(
		j,
		"ramDiskId",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetSecurityGroups(val *[]*string) {
	_jsii_.Set(
		j,
		"securityGroups",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetSpotPrice(val *string) {
	_jsii_.Set(
		j,
		"spotPrice",
		val,
	)
}

func (j *jsiiProxy_CfnLaunchConfiguration) SetUserData(val *string) {
	_jsii_.Set(
		j,
		"userData",
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
func CfnLaunchConfiguration_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnLaunchConfiguration",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnLaunchConfiguration_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnLaunchConfiguration",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnLaunchConfiguration_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnLaunchConfiguration",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnLaunchConfiguration_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.CfnLaunchConfiguration",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnLaunchConfiguration) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnLaunchConfiguration) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnLaunchConfiguration) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnLaunchConfiguration) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnLaunchConfiguration) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) OnPrepare() {
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
func (c *jsiiProxy_CfnLaunchConfiguration) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnLaunchConfiguration) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnLaunchConfiguration) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnLaunchConfiguration) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnLaunchConfiguration) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLaunchConfiguration) ToString() *string {
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
func (c *jsiiProxy_CfnLaunchConfiguration) Validate() *[]*string {
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
func (c *jsiiProxy_CfnLaunchConfiguration) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnLaunchConfiguration_BlockDeviceMappingProperty struct {
	// `CfnLaunchConfiguration.BlockDeviceMappingProperty.DeviceName`.
	DeviceName *string `json:"deviceName"`
	// `CfnLaunchConfiguration.BlockDeviceMappingProperty.Ebs`.
	Ebs interface{} `json:"ebs"`
	// `CfnLaunchConfiguration.BlockDeviceMappingProperty.NoDevice`.
	NoDevice interface{} `json:"noDevice"`
	// `CfnLaunchConfiguration.BlockDeviceMappingProperty.VirtualName`.
	VirtualName *string `json:"virtualName"`
}

type CfnLaunchConfiguration_BlockDeviceProperty struct {
	// `CfnLaunchConfiguration.BlockDeviceProperty.DeleteOnTermination`.
	DeleteOnTermination interface{} `json:"deleteOnTermination"`
	// `CfnLaunchConfiguration.BlockDeviceProperty.Encrypted`.
	Encrypted interface{} `json:"encrypted"`
	// `CfnLaunchConfiguration.BlockDeviceProperty.Iops`.
	Iops *float64 `json:"iops"`
	// `CfnLaunchConfiguration.BlockDeviceProperty.SnapshotId`.
	SnapshotId *string `json:"snapshotId"`
	// `CfnLaunchConfiguration.BlockDeviceProperty.VolumeSize`.
	VolumeSize *float64 `json:"volumeSize"`
	// `CfnLaunchConfiguration.BlockDeviceProperty.VolumeType`.
	VolumeType *string `json:"volumeType"`
}

type CfnLaunchConfiguration_MetadataOptionsProperty struct {
	// `CfnLaunchConfiguration.MetadataOptionsProperty.HttpEndpoint`.
	HttpEndpoint *string `json:"httpEndpoint"`
	// `CfnLaunchConfiguration.MetadataOptionsProperty.HttpPutResponseHopLimit`.
	HttpPutResponseHopLimit *float64 `json:"httpPutResponseHopLimit"`
	// `CfnLaunchConfiguration.MetadataOptionsProperty.HttpTokens`.
	HttpTokens *string `json:"httpTokens"`
}

// Properties for defining a `AWS::AutoScaling::LaunchConfiguration`.
type CfnLaunchConfigurationProps struct {
	// `AWS::AutoScaling::LaunchConfiguration.ImageId`.
	ImageId *string `json:"imageId"`
	// `AWS::AutoScaling::LaunchConfiguration.InstanceType`.
	InstanceType *string `json:"instanceType"`
	// `AWS::AutoScaling::LaunchConfiguration.AssociatePublicIpAddress`.
	AssociatePublicIpAddress interface{} `json:"associatePublicIpAddress"`
	// `AWS::AutoScaling::LaunchConfiguration.BlockDeviceMappings`.
	BlockDeviceMappings interface{} `json:"blockDeviceMappings"`
	// `AWS::AutoScaling::LaunchConfiguration.ClassicLinkVPCId`.
	ClassicLinkVpcId *string `json:"classicLinkVpcId"`
	// `AWS::AutoScaling::LaunchConfiguration.ClassicLinkVPCSecurityGroups`.
	ClassicLinkVpcSecurityGroups *[]*string `json:"classicLinkVpcSecurityGroups"`
	// `AWS::AutoScaling::LaunchConfiguration.EbsOptimized`.
	EbsOptimized interface{} `json:"ebsOptimized"`
	// `AWS::AutoScaling::LaunchConfiguration.IamInstanceProfile`.
	IamInstanceProfile *string `json:"iamInstanceProfile"`
	// `AWS::AutoScaling::LaunchConfiguration.InstanceId`.
	InstanceId *string `json:"instanceId"`
	// `AWS::AutoScaling::LaunchConfiguration.InstanceMonitoring`.
	InstanceMonitoring interface{} `json:"instanceMonitoring"`
	// `AWS::AutoScaling::LaunchConfiguration.KernelId`.
	KernelId *string `json:"kernelId"`
	// `AWS::AutoScaling::LaunchConfiguration.KeyName`.
	KeyName *string `json:"keyName"`
	// `AWS::AutoScaling::LaunchConfiguration.LaunchConfigurationName`.
	LaunchConfigurationName *string `json:"launchConfigurationName"`
	// `AWS::AutoScaling::LaunchConfiguration.MetadataOptions`.
	MetadataOptions interface{} `json:"metadataOptions"`
	// `AWS::AutoScaling::LaunchConfiguration.PlacementTenancy`.
	PlacementTenancy *string `json:"placementTenancy"`
	// `AWS::AutoScaling::LaunchConfiguration.RamDiskId`.
	RamDiskId *string `json:"ramDiskId"`
	// `AWS::AutoScaling::LaunchConfiguration.SecurityGroups`.
	SecurityGroups *[]*string `json:"securityGroups"`
	// `AWS::AutoScaling::LaunchConfiguration.SpotPrice`.
	SpotPrice *string `json:"spotPrice"`
	// `AWS::AutoScaling::LaunchConfiguration.UserData`.
	UserData *string `json:"userData"`
}

// A CloudFormation `AWS::AutoScaling::LifecycleHook`.
type CfnLifecycleHook interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AutoScalingGroupName() *string
	SetAutoScalingGroupName(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DefaultResult() *string
	SetDefaultResult(val *string)
	HeartbeatTimeout() *float64
	SetHeartbeatTimeout(val *float64)
	LifecycleHookName() *string
	SetLifecycleHookName(val *string)
	LifecycleTransition() *string
	SetLifecycleTransition(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	NotificationMetadata() *string
	SetNotificationMetadata(val *string)
	NotificationTargetArn() *string
	SetNotificationTargetArn(val *string)
	Ref() *string
	RoleArn() *string
	SetRoleArn(val *string)
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

// The jsii proxy struct for CfnLifecycleHook
type jsiiProxy_CfnLifecycleHook struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnLifecycleHook) AutoScalingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) DefaultResult() *string {
	var returns *string
	_jsii_.Get(
		j,
		"defaultResult",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) HeartbeatTimeout() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"heartbeatTimeout",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) LifecycleHookName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"lifecycleHookName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) LifecycleTransition() *string {
	var returns *string
	_jsii_.Get(
		j,
		"lifecycleTransition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) NotificationMetadata() *string {
	var returns *string
	_jsii_.Get(
		j,
		"notificationMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) NotificationTargetArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"notificationTargetArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLifecycleHook) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::AutoScaling::LifecycleHook`.
func NewCfnLifecycleHook(scope awscdk.Construct, id *string, props *CfnLifecycleHookProps) CfnLifecycleHook {
	_init_.Initialize()

	j := jsiiProxy_CfnLifecycleHook{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnLifecycleHook",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::AutoScaling::LifecycleHook`.
func NewCfnLifecycleHook_Override(c CfnLifecycleHook, scope awscdk.Construct, id *string, props *CfnLifecycleHookProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnLifecycleHook",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetAutoScalingGroupName(val *string) {
	_jsii_.Set(
		j,
		"autoScalingGroupName",
		val,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetDefaultResult(val *string) {
	_jsii_.Set(
		j,
		"defaultResult",
		val,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetHeartbeatTimeout(val *float64) {
	_jsii_.Set(
		j,
		"heartbeatTimeout",
		val,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetLifecycleHookName(val *string) {
	_jsii_.Set(
		j,
		"lifecycleHookName",
		val,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetLifecycleTransition(val *string) {
	_jsii_.Set(
		j,
		"lifecycleTransition",
		val,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetNotificationMetadata(val *string) {
	_jsii_.Set(
		j,
		"notificationMetadata",
		val,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetNotificationTargetArn(val *string) {
	_jsii_.Set(
		j,
		"notificationTargetArn",
		val,
	)
}

func (j *jsiiProxy_CfnLifecycleHook) SetRoleArn(val *string) {
	_jsii_.Set(
		j,
		"roleArn",
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
func CfnLifecycleHook_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnLifecycleHook",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnLifecycleHook_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnLifecycleHook",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnLifecycleHook_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnLifecycleHook",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnLifecycleHook_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.CfnLifecycleHook",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnLifecycleHook) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnLifecycleHook) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnLifecycleHook) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnLifecycleHook) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnLifecycleHook) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnLifecycleHook) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnLifecycleHook) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnLifecycleHook) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnLifecycleHook) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnLifecycleHook) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnLifecycleHook) OnPrepare() {
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
func (c *jsiiProxy_CfnLifecycleHook) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLifecycleHook) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnLifecycleHook) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnLifecycleHook) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnLifecycleHook) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnLifecycleHook) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnLifecycleHook) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLifecycleHook) ToString() *string {
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
func (c *jsiiProxy_CfnLifecycleHook) Validate() *[]*string {
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
func (c *jsiiProxy_CfnLifecycleHook) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::AutoScaling::LifecycleHook`.
type CfnLifecycleHookProps struct {
	// `AWS::AutoScaling::LifecycleHook.AutoScalingGroupName`.
	AutoScalingGroupName *string `json:"autoScalingGroupName"`
	// `AWS::AutoScaling::LifecycleHook.LifecycleTransition`.
	LifecycleTransition *string `json:"lifecycleTransition"`
	// `AWS::AutoScaling::LifecycleHook.DefaultResult`.
	DefaultResult *string `json:"defaultResult"`
	// `AWS::AutoScaling::LifecycleHook.HeartbeatTimeout`.
	HeartbeatTimeout *float64 `json:"heartbeatTimeout"`
	// `AWS::AutoScaling::LifecycleHook.LifecycleHookName`.
	LifecycleHookName *string `json:"lifecycleHookName"`
	// `AWS::AutoScaling::LifecycleHook.NotificationMetadata`.
	NotificationMetadata *string `json:"notificationMetadata"`
	// `AWS::AutoScaling::LifecycleHook.NotificationTargetARN`.
	NotificationTargetArn *string `json:"notificationTargetArn"`
	// `AWS::AutoScaling::LifecycleHook.RoleARN`.
	RoleArn *string `json:"roleArn"`
}

// A CloudFormation `AWS::AutoScaling::ScalingPolicy`.
type CfnScalingPolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AdjustmentType() *string
	SetAdjustmentType(val *string)
	AutoScalingGroupName() *string
	SetAutoScalingGroupName(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Cooldown() *string
	SetCooldown(val *string)
	CreationStack() *[]*string
	EstimatedInstanceWarmup() *float64
	SetEstimatedInstanceWarmup(val *float64)
	LogicalId() *string
	MetricAggregationType() *string
	SetMetricAggregationType(val *string)
	MinAdjustmentMagnitude() *float64
	SetMinAdjustmentMagnitude(val *float64)
	Node() awscdk.ConstructNode
	PolicyType() *string
	SetPolicyType(val *string)
	Ref() *string
	ScalingAdjustment() *float64
	SetScalingAdjustment(val *float64)
	Stack() awscdk.Stack
	StepAdjustments() interface{}
	SetStepAdjustments(val interface{})
	TargetTrackingConfiguration() interface{}
	SetTargetTrackingConfiguration(val interface{})
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

// The jsii proxy struct for CfnScalingPolicy
type jsiiProxy_CfnScalingPolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnScalingPolicy) AdjustmentType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"adjustmentType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) AutoScalingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) Cooldown() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cooldown",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) EstimatedInstanceWarmup() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"estimatedInstanceWarmup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) MetricAggregationType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"metricAggregationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) MinAdjustmentMagnitude() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"minAdjustmentMagnitude",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) PolicyType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"policyType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) ScalingAdjustment() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"scalingAdjustment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) StepAdjustments() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stepAdjustments",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) TargetTrackingConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"targetTrackingConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScalingPolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::AutoScaling::ScalingPolicy`.
func NewCfnScalingPolicy(scope awscdk.Construct, id *string, props *CfnScalingPolicyProps) CfnScalingPolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnScalingPolicy{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnScalingPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::AutoScaling::ScalingPolicy`.
func NewCfnScalingPolicy_Override(c CfnScalingPolicy, scope awscdk.Construct, id *string, props *CfnScalingPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnScalingPolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetAdjustmentType(val *string) {
	_jsii_.Set(
		j,
		"adjustmentType",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetAutoScalingGroupName(val *string) {
	_jsii_.Set(
		j,
		"autoScalingGroupName",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetCooldown(val *string) {
	_jsii_.Set(
		j,
		"cooldown",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetEstimatedInstanceWarmup(val *float64) {
	_jsii_.Set(
		j,
		"estimatedInstanceWarmup",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetMetricAggregationType(val *string) {
	_jsii_.Set(
		j,
		"metricAggregationType",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetMinAdjustmentMagnitude(val *float64) {
	_jsii_.Set(
		j,
		"minAdjustmentMagnitude",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetPolicyType(val *string) {
	_jsii_.Set(
		j,
		"policyType",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetScalingAdjustment(val *float64) {
	_jsii_.Set(
		j,
		"scalingAdjustment",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetStepAdjustments(val interface{}) {
	_jsii_.Set(
		j,
		"stepAdjustments",
		val,
	)
}

func (j *jsiiProxy_CfnScalingPolicy) SetTargetTrackingConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"targetTrackingConfiguration",
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
func CfnScalingPolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnScalingPolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnScalingPolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnScalingPolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnScalingPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnScalingPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnScalingPolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.CfnScalingPolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnScalingPolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnScalingPolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnScalingPolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnScalingPolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnScalingPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnScalingPolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnScalingPolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnScalingPolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnScalingPolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnScalingPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnScalingPolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnScalingPolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnScalingPolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnScalingPolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnScalingPolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnScalingPolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnScalingPolicy) ToString() *string {
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
func (c *jsiiProxy_CfnScalingPolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnScalingPolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnScalingPolicy_CustomizedMetricSpecificationProperty struct {
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.MetricName`.
	MetricName *string `json:"metricName"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Namespace`.
	Namespace *string `json:"namespace"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Statistic`.
	Statistic *string `json:"statistic"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Dimensions`.
	Dimensions interface{} `json:"dimensions"`
	// `CfnScalingPolicy.CustomizedMetricSpecificationProperty.Unit`.
	Unit *string `json:"unit"`
}

type CfnScalingPolicy_MetricDimensionProperty struct {
	// `CfnScalingPolicy.MetricDimensionProperty.Name`.
	Name *string `json:"name"`
	// `CfnScalingPolicy.MetricDimensionProperty.Value`.
	Value *string `json:"value"`
}

type CfnScalingPolicy_PredefinedMetricSpecificationProperty struct {
	// `CfnScalingPolicy.PredefinedMetricSpecificationProperty.PredefinedMetricType`.
	PredefinedMetricType *string `json:"predefinedMetricType"`
	// `CfnScalingPolicy.PredefinedMetricSpecificationProperty.ResourceLabel`.
	ResourceLabel *string `json:"resourceLabel"`
}

type CfnScalingPolicy_StepAdjustmentProperty struct {
	// `CfnScalingPolicy.StepAdjustmentProperty.ScalingAdjustment`.
	ScalingAdjustment *float64 `json:"scalingAdjustment"`
	// `CfnScalingPolicy.StepAdjustmentProperty.MetricIntervalLowerBound`.
	MetricIntervalLowerBound *float64 `json:"metricIntervalLowerBound"`
	// `CfnScalingPolicy.StepAdjustmentProperty.MetricIntervalUpperBound`.
	MetricIntervalUpperBound *float64 `json:"metricIntervalUpperBound"`
}

type CfnScalingPolicy_TargetTrackingConfigurationProperty struct {
	// `CfnScalingPolicy.TargetTrackingConfigurationProperty.TargetValue`.
	TargetValue *float64 `json:"targetValue"`
	// `CfnScalingPolicy.TargetTrackingConfigurationProperty.CustomizedMetricSpecification`.
	CustomizedMetricSpecification interface{} `json:"customizedMetricSpecification"`
	// `CfnScalingPolicy.TargetTrackingConfigurationProperty.DisableScaleIn`.
	DisableScaleIn interface{} `json:"disableScaleIn"`
	// `CfnScalingPolicy.TargetTrackingConfigurationProperty.PredefinedMetricSpecification`.
	PredefinedMetricSpecification interface{} `json:"predefinedMetricSpecification"`
}

// Properties for defining a `AWS::AutoScaling::ScalingPolicy`.
type CfnScalingPolicyProps struct {
	// `AWS::AutoScaling::ScalingPolicy.AutoScalingGroupName`.
	AutoScalingGroupName *string `json:"autoScalingGroupName"`
	// `AWS::AutoScaling::ScalingPolicy.AdjustmentType`.
	AdjustmentType *string `json:"adjustmentType"`
	// `AWS::AutoScaling::ScalingPolicy.Cooldown`.
	Cooldown *string `json:"cooldown"`
	// `AWS::AutoScaling::ScalingPolicy.EstimatedInstanceWarmup`.
	EstimatedInstanceWarmup *float64 `json:"estimatedInstanceWarmup"`
	// `AWS::AutoScaling::ScalingPolicy.MetricAggregationType`.
	MetricAggregationType *string `json:"metricAggregationType"`
	// `AWS::AutoScaling::ScalingPolicy.MinAdjustmentMagnitude`.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
	// `AWS::AutoScaling::ScalingPolicy.PolicyType`.
	PolicyType *string `json:"policyType"`
	// `AWS::AutoScaling::ScalingPolicy.ScalingAdjustment`.
	ScalingAdjustment *float64 `json:"scalingAdjustment"`
	// `AWS::AutoScaling::ScalingPolicy.StepAdjustments`.
	StepAdjustments interface{} `json:"stepAdjustments"`
	// `AWS::AutoScaling::ScalingPolicy.TargetTrackingConfiguration`.
	TargetTrackingConfiguration interface{} `json:"targetTrackingConfiguration"`
}

// A CloudFormation `AWS::AutoScaling::ScheduledAction`.
type CfnScheduledAction interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AutoScalingGroupName() *string
	SetAutoScalingGroupName(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DesiredCapacity() *float64
	SetDesiredCapacity(val *float64)
	EndTime() *string
	SetEndTime(val *string)
	LogicalId() *string
	MaxSize() *float64
	SetMaxSize(val *float64)
	MinSize() *float64
	SetMinSize(val *float64)
	Node() awscdk.ConstructNode
	Recurrence() *string
	SetRecurrence(val *string)
	Ref() *string
	Stack() awscdk.Stack
	StartTime() *string
	SetStartTime(val *string)
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

// The jsii proxy struct for CfnScheduledAction
type jsiiProxy_CfnScheduledAction struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnScheduledAction) AutoScalingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) DesiredCapacity() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"desiredCapacity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) EndTime() *string {
	var returns *string
	_jsii_.Get(
		j,
		"endTime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) MaxSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) MinSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"minSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) Recurrence() *string {
	var returns *string
	_jsii_.Get(
		j,
		"recurrence",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) StartTime() *string {
	var returns *string
	_jsii_.Get(
		j,
		"startTime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnScheduledAction) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::AutoScaling::ScheduledAction`.
func NewCfnScheduledAction(scope awscdk.Construct, id *string, props *CfnScheduledActionProps) CfnScheduledAction {
	_init_.Initialize()

	j := jsiiProxy_CfnScheduledAction{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnScheduledAction",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::AutoScaling::ScheduledAction`.
func NewCfnScheduledAction_Override(c CfnScheduledAction, scope awscdk.Construct, id *string, props *CfnScheduledActionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnScheduledAction",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnScheduledAction) SetAutoScalingGroupName(val *string) {
	_jsii_.Set(
		j,
		"autoScalingGroupName",
		val,
	)
}

func (j *jsiiProxy_CfnScheduledAction) SetDesiredCapacity(val *float64) {
	_jsii_.Set(
		j,
		"desiredCapacity",
		val,
	)
}

func (j *jsiiProxy_CfnScheduledAction) SetEndTime(val *string) {
	_jsii_.Set(
		j,
		"endTime",
		val,
	)
}

func (j *jsiiProxy_CfnScheduledAction) SetMaxSize(val *float64) {
	_jsii_.Set(
		j,
		"maxSize",
		val,
	)
}

func (j *jsiiProxy_CfnScheduledAction) SetMinSize(val *float64) {
	_jsii_.Set(
		j,
		"minSize",
		val,
	)
}

func (j *jsiiProxy_CfnScheduledAction) SetRecurrence(val *string) {
	_jsii_.Set(
		j,
		"recurrence",
		val,
	)
}

func (j *jsiiProxy_CfnScheduledAction) SetStartTime(val *string) {
	_jsii_.Set(
		j,
		"startTime",
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
func CfnScheduledAction_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnScheduledAction",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnScheduledAction_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnScheduledAction",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnScheduledAction_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnScheduledAction",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnScheduledAction_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.CfnScheduledAction",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnScheduledAction) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnScheduledAction) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnScheduledAction) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnScheduledAction) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnScheduledAction) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnScheduledAction) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnScheduledAction) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnScheduledAction) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnScheduledAction) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnScheduledAction) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnScheduledAction) OnPrepare() {
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
func (c *jsiiProxy_CfnScheduledAction) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnScheduledAction) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnScheduledAction) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnScheduledAction) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnScheduledAction) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnScheduledAction) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnScheduledAction) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnScheduledAction) ToString() *string {
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
func (c *jsiiProxy_CfnScheduledAction) Validate() *[]*string {
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
func (c *jsiiProxy_CfnScheduledAction) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::AutoScaling::ScheduledAction`.
type CfnScheduledActionProps struct {
	// `AWS::AutoScaling::ScheduledAction.AutoScalingGroupName`.
	AutoScalingGroupName *string `json:"autoScalingGroupName"`
	// `AWS::AutoScaling::ScheduledAction.DesiredCapacity`.
	DesiredCapacity *float64 `json:"desiredCapacity"`
	// `AWS::AutoScaling::ScheduledAction.EndTime`.
	EndTime *string `json:"endTime"`
	// `AWS::AutoScaling::ScheduledAction.MaxSize`.
	MaxSize *float64 `json:"maxSize"`
	// `AWS::AutoScaling::ScheduledAction.MinSize`.
	MinSize *float64 `json:"minSize"`
	// `AWS::AutoScaling::ScheduledAction.Recurrence`.
	Recurrence *string `json:"recurrence"`
	// `AWS::AutoScaling::ScheduledAction.StartTime`.
	StartTime *string `json:"startTime"`
}

// A CloudFormation `AWS::AutoScaling::WarmPool`.
type CfnWarmPool interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AutoScalingGroupName() *string
	SetAutoScalingGroupName(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	MaxGroupPreparedCapacity() *float64
	SetMaxGroupPreparedCapacity(val *float64)
	MinSize() *float64
	SetMinSize(val *float64)
	Node() awscdk.ConstructNode
	PoolState() *string
	SetPoolState(val *string)
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

// The jsii proxy struct for CfnWarmPool
type jsiiProxy_CfnWarmPool struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnWarmPool) AutoScalingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) MaxGroupPreparedCapacity() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxGroupPreparedCapacity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) MinSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"minSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) PoolState() *string {
	var returns *string
	_jsii_.Get(
		j,
		"poolState",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWarmPool) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::AutoScaling::WarmPool`.
func NewCfnWarmPool(scope awscdk.Construct, id *string, props *CfnWarmPoolProps) CfnWarmPool {
	_init_.Initialize()

	j := jsiiProxy_CfnWarmPool{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnWarmPool",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::AutoScaling::WarmPool`.
func NewCfnWarmPool_Override(c CfnWarmPool, scope awscdk.Construct, id *string, props *CfnWarmPoolProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.CfnWarmPool",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnWarmPool) SetAutoScalingGroupName(val *string) {
	_jsii_.Set(
		j,
		"autoScalingGroupName",
		val,
	)
}

func (j *jsiiProxy_CfnWarmPool) SetMaxGroupPreparedCapacity(val *float64) {
	_jsii_.Set(
		j,
		"maxGroupPreparedCapacity",
		val,
	)
}

func (j *jsiiProxy_CfnWarmPool) SetMinSize(val *float64) {
	_jsii_.Set(
		j,
		"minSize",
		val,
	)
}

func (j *jsiiProxy_CfnWarmPool) SetPoolState(val *string) {
	_jsii_.Set(
		j,
		"poolState",
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
func CfnWarmPool_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnWarmPool",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnWarmPool_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnWarmPool",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnWarmPool_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.CfnWarmPool",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnWarmPool_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.CfnWarmPool",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnWarmPool) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnWarmPool) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnWarmPool) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnWarmPool) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnWarmPool) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnWarmPool) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnWarmPool) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnWarmPool) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnWarmPool) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnWarmPool) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnWarmPool) OnPrepare() {
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
func (c *jsiiProxy_CfnWarmPool) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnWarmPool) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnWarmPool) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnWarmPool) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnWarmPool) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnWarmPool) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnWarmPool) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnWarmPool) ToString() *string {
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
func (c *jsiiProxy_CfnWarmPool) Validate() *[]*string {
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
func (c *jsiiProxy_CfnWarmPool) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::AutoScaling::WarmPool`.
type CfnWarmPoolProps struct {
	// `AWS::AutoScaling::WarmPool.AutoScalingGroupName`.
	AutoScalingGroupName *string `json:"autoScalingGroupName"`
	// `AWS::AutoScaling::WarmPool.MaxGroupPreparedCapacity`.
	MaxGroupPreparedCapacity *float64 `json:"maxGroupPreparedCapacity"`
	// `AWS::AutoScaling::WarmPool.MinSize`.
	MinSize *float64 `json:"minSize"`
	// `AWS::AutoScaling::WarmPool.PoolState`.
	PoolState *string `json:"poolState"`
}

// Basic properties of an AutoScalingGroup, except the exact machines to run and where they should run.
//
// Constructs that want to create AutoScalingGroups can inherit
// this interface and specialize the essential parts in various ways.
// Experimental.
type CommonAutoScalingGroupProps struct {
	// Whether the instances can initiate connections to anywhere by default.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Whether instances in the Auto Scaling Group should have public IP addresses associated with them.
	// Experimental.
	AssociatePublicIpAddress *bool `json:"associatePublicIpAddress"`
	// The name of the Auto Scaling group.
	//
	// This name must be unique per Region per account.
	// Experimental.
	AutoScalingGroupName *string `json:"autoScalingGroupName"`
	// Specifies how block devices are exposed to the instance. You can specify virtual devices and EBS volumes.
	//
	// Each instance that is launched has an associated root device volume,
	// either an Amazon EBS volume or an instance store volume.
	// You can use block device mappings to specify additional EBS volumes or
	// instance store volumes to attach to an instance when it is launched.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/block-device-mapping-concepts.html
	//
	// Experimental.
	BlockDevices *[]*BlockDevice `json:"blockDevices"`
	// Default scaling cooldown for this AutoScalingGroup.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Initial amount of instances in the fleet.
	//
	// If this is set to a number, every deployment will reset the amount of
	// instances to this number. It is recommended to leave this value blank.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-group.html#cfn-as-group-desiredcapacity
	//
	// Experimental.
	DesiredCapacity *float64 `json:"desiredCapacity"`
	// Enable monitoring for group metrics, these metrics describe the group rather than any of its instances.
	//
	// To report all group metrics use `GroupMetrics.all()`
	// Group metrics are reported in a granularity of 1 minute at no additional charge.
	// Experimental.
	GroupMetrics *[]GroupMetrics `json:"groupMetrics"`
	// Configuration for health checks.
	// Experimental.
	HealthCheck HealthCheck `json:"healthCheck"`
	// If the ASG has scheduled actions, don't reset unchanged group sizes.
	//
	// Only used if the ASG has scheduled actions (which may scale your ASG up
	// or down regardless of cdk deployments). If true, the size of the group
	// will only be reset if it has been changed in the CDK app. If false, the
	// sizes will always be changed back to what they were in the CDK app
	// on deployment.
	// Experimental.
	IgnoreUnmodifiedSizeProperties *bool `json:"ignoreUnmodifiedSizeProperties"`
	// Controls whether instances in this group are launched with detailed or basic monitoring.
	//
	// When detailed monitoring is enabled, Amazon CloudWatch generates metrics every minute and your account
	// is charged a fee. When you disable detailed monitoring, CloudWatch generates metrics every 5 minutes.
	// See: https://docs.aws.amazon.com/autoscaling/latest/userguide/as-instance-monitoring.html#enable-as-instance-metrics
	//
	// Experimental.
	InstanceMonitoring Monitoring `json:"instanceMonitoring"`
	// Name of SSH keypair to grant access to instances.
	// Experimental.
	KeyName *string `json:"keyName"`
	// Maximum number of instances in the fleet.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// The maximum amount of time that an instance can be in service.
	//
	// The maximum duration applies
	// to all current and future instances in the group. As an instance approaches its maximum duration,
	// it is terminated and replaced, and cannot be used again.
	//
	// You must specify a value of at least 604,800 seconds (7 days). To clear a previously set value,
	// leave this property undefined.
	// See: https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-max-instance-lifetime.html
	//
	// Experimental.
	MaxInstanceLifetime awscdk.Duration `json:"maxInstanceLifetime"`
	// Minimum number of instances in the fleet.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// Whether newly-launched instances are protected from termination by Amazon EC2 Auto Scaling when scaling in.
	//
	// By default, Auto Scaling can terminate an instance at any time after launch
	// when scaling in an Auto Scaling Group, subject to the group's termination
	// policy. However, you may wish to protect newly-launched instances from
	// being scaled in if they are going to run critical applications that should
	// not be prematurely terminated.
	//
	// This flag must be enabled if the Auto Scaling Group will be associated with
	// an ECS Capacity Provider with managed termination protection.
	// Experimental.
	NewInstancesProtectedFromScaleIn *bool `json:"newInstancesProtectedFromScaleIn"`
	// Configure autoscaling group to send notifications about fleet changes to an SNS topic(s).
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-group.html#cfn-as-group-notificationconfigurations
	//
	// Experimental.
	Notifications *[]*NotificationConfiguration `json:"notifications"`
	// SNS topic to send notifications about fleet changes.
	// Deprecated: use `notifications`
	NotificationsTopic awssns.ITopic `json:"notificationsTopic"`
	// Configuration for replacing updates.
	//
	// Only used if updateType == UpdateType.ReplacingUpdate. Specifies how
	// many instances must signal success for the update to succeed.
	// Deprecated: Use `signals` instead
	ReplacingUpdateMinSuccessfulInstancesPercent *float64 `json:"replacingUpdateMinSuccessfulInstancesPercent"`
	// How many ResourceSignal calls CloudFormation expects before the resource is considered created.
	// Deprecated: Use `signals` instead.
	ResourceSignalCount *float64 `json:"resourceSignalCount"`
	// The length of time to wait for the resourceSignalCount.
	//
	// The maximum value is 43200 (12 hours).
	// Deprecated: Use `signals` instead.
	ResourceSignalTimeout awscdk.Duration `json:"resourceSignalTimeout"`
	// Configuration for rolling updates.
	//
	// Only used if updateType == UpdateType.RollingUpdate.
	// Deprecated: Use `updatePolicy` instead
	RollingUpdateConfiguration *RollingUpdateConfiguration `json:"rollingUpdateConfiguration"`
	// Configure waiting for signals during deployment.
	//
	// Use this to pause the CloudFormation deployment to wait for the instances
	// in the AutoScalingGroup to report successful startup during
	// creation and updates. The UserData script needs to invoke `cfn-signal`
	// with a success or failure code after it is done setting up the instance.
	//
	// Without waiting for signals, the CloudFormation deployment will proceed as
	// soon as the AutoScalingGroup has been created or updated but before the
	// instances in the group have been started.
	//
	// For example, to have instances wait for an Elastic Load Balancing health check before
	// they signal success, add a health-check verification by using the
	// cfn-init helper script. For an example, see the verify_instance_health
	// command in the Auto Scaling rolling updates sample template:
	//
	// https://github.com/awslabs/aws-cloudformation-templates/blob/master/aws/services/AutoScaling/AutoScalingRollingUpdates.yaml
	// Experimental.
	Signals Signals `json:"signals"`
	// The maximum hourly price (in USD) to be paid for any Spot Instance launched to fulfill the request.
	//
	// Spot Instances are
	// launched when the price you specify exceeds the current Spot market price.
	// Experimental.
	SpotPrice *string `json:"spotPrice"`
	// What to do when an AutoScalingGroup's instance configuration is changed.
	//
	// This is applied when any of the settings on the ASG are changed that
	// affect how the instances should be created (VPC, instance type, startup
	// scripts, etc.). It indicates how the existing instances should be
	// replaced with new instances matching the new config. By default, nothing
	// is done and only new instances are launched with the new config.
	// Experimental.
	UpdatePolicy UpdatePolicy `json:"updatePolicy"`
	// What to do when an AutoScalingGroup's instance configuration is changed.
	//
	// This is applied when any of the settings on the ASG are changed that
	// affect how the instances should be created (VPC, instance type, startup
	// scripts, etc.). It indicates how the existing instances should be
	// replaced with new instances matching the new config. By default, nothing
	// is done and only new instances are launched with the new config.
	// Deprecated: Use `updatePolicy` instead
	UpdateType UpdateType `json:"updateType"`
	// Where to place instances within the VPC.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
}

// Properties for enabling scaling based on CPU utilization.
// Experimental.
type CpuUtilizationScalingProps struct {
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the autoscaling group. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// group.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// Target average CPU utilization across the task.
	// Experimental.
	TargetUtilizationPercent *float64 `json:"targetUtilizationPercent"`
}

// Options to configure a cron expression.
//
// All fields are strings so you can use complex expressions. Absence of
// a field implies '*' or '?', whichever one is appropriate.
// See: http://crontab.org/
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
}

// Experimental.
type DefaultResult string

const (
	DefaultResult_CONTINUE DefaultResult = "CONTINUE"
	DefaultResult_ABANDON DefaultResult = "ABANDON"
)

// Block device options for an EBS volume.
// Experimental.
type EbsDeviceOptions struct {
	// Indicates whether to delete the volume when the instance is terminated.
	// Experimental.
	DeleteOnTermination *bool `json:"deleteOnTermination"`
	// The number of I/O operations per second (IOPS) to provision for the volume.
	//
	// Must only be set for {@link volumeType}: {@link EbsDeviceVolumeType.IO1}
	//
	// The maximum ratio of IOPS to volume size (in GiB) is 50:1, so for 5,000 provisioned IOPS,
	// you need at least 100 GiB storage on the volume.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	Iops *float64 `json:"iops"`
	// The EBS volume type.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	VolumeType EbsDeviceVolumeType `json:"volumeType"`
	// Specifies whether the EBS volume is encrypted.
	//
	// Encrypted EBS volumes can only be attached to instances that support Amazon EBS encryption
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html#EBSEncryption_supported_instances
	//
	// Experimental.
	Encrypted *bool `json:"encrypted"`
}

// Base block device options for an EBS volume.
// Experimental.
type EbsDeviceOptionsBase struct {
	// Indicates whether to delete the volume when the instance is terminated.
	// Experimental.
	DeleteOnTermination *bool `json:"deleteOnTermination"`
	// The number of I/O operations per second (IOPS) to provision for the volume.
	//
	// Must only be set for {@link volumeType}: {@link EbsDeviceVolumeType.IO1}
	//
	// The maximum ratio of IOPS to volume size (in GiB) is 50:1, so for 5,000 provisioned IOPS,
	// you need at least 100 GiB storage on the volume.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	Iops *float64 `json:"iops"`
	// The EBS volume type.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	VolumeType EbsDeviceVolumeType `json:"volumeType"`
}

// Properties of an EBS block device.
// Experimental.
type EbsDeviceProps struct {
	// Indicates whether to delete the volume when the instance is terminated.
	// Experimental.
	DeleteOnTermination *bool `json:"deleteOnTermination"`
	// The number of I/O operations per second (IOPS) to provision for the volume.
	//
	// Must only be set for {@link volumeType}: {@link EbsDeviceVolumeType.IO1}
	//
	// The maximum ratio of IOPS to volume size (in GiB) is 50:1, so for 5,000 provisioned IOPS,
	// you need at least 100 GiB storage on the volume.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	Iops *float64 `json:"iops"`
	// The EBS volume type.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	VolumeType EbsDeviceVolumeType `json:"volumeType"`
	// The volume size, in Gibibytes (GiB).
	//
	// If you specify volumeSize, it must be equal or greater than the size of the snapshot.
	// Experimental.
	VolumeSize *float64 `json:"volumeSize"`
	// The snapshot ID of the volume to use.
	// Experimental.
	SnapshotId *string `json:"snapshotId"`
}

// Block device options for an EBS volume created from a snapshot.
// Experimental.
type EbsDeviceSnapshotOptions struct {
	// Indicates whether to delete the volume when the instance is terminated.
	// Experimental.
	DeleteOnTermination *bool `json:"deleteOnTermination"`
	// The number of I/O operations per second (IOPS) to provision for the volume.
	//
	// Must only be set for {@link volumeType}: {@link EbsDeviceVolumeType.IO1}
	//
	// The maximum ratio of IOPS to volume size (in GiB) is 50:1, so for 5,000 provisioned IOPS,
	// you need at least 100 GiB storage on the volume.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	Iops *float64 `json:"iops"`
	// The EBS volume type.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	//
	// Experimental.
	VolumeType EbsDeviceVolumeType `json:"volumeType"`
	// The volume size, in Gibibytes (GiB).
	//
	// If you specify volumeSize, it must be equal or greater than the size of the snapshot.
	// Experimental.
	VolumeSize *float64 `json:"volumeSize"`
}

// Supported EBS volume types for blockDevices.
// Experimental.
type EbsDeviceVolumeType string

const (
	EbsDeviceVolumeType_STANDARD EbsDeviceVolumeType = "STANDARD"
	EbsDeviceVolumeType_IO1 EbsDeviceVolumeType = "IO1"
	EbsDeviceVolumeType_IO2 EbsDeviceVolumeType = "IO2"
	EbsDeviceVolumeType_GP2 EbsDeviceVolumeType = "GP2"
	EbsDeviceVolumeType_GP3 EbsDeviceVolumeType = "GP3"
	EbsDeviceVolumeType_ST1 EbsDeviceVolumeType = "ST1"
	EbsDeviceVolumeType_SC1 EbsDeviceVolumeType = "SC1"
)

// EC2 Heath check options.
// Experimental.
type Ec2HealthCheckOptions struct {
	// Specified the time Auto Scaling waits before checking the health status of an EC2 instance that has come into service.
	// Experimental.
	Grace awscdk.Duration `json:"grace"`
}

// ELB Heath check options.
// Experimental.
type ElbHealthCheckOptions struct {
	// Specified the time Auto Scaling waits before checking the health status of an EC2 instance that has come into service.
	//
	// This option is required for ELB health checks.
	// Experimental.
	Grace awscdk.Duration `json:"grace"`
}

// Group metrics that an Auto Scaling group sends to Amazon CloudWatch.
// Experimental.
type GroupMetric interface {
	Name() *string
}

// The jsii proxy struct for GroupMetric
type jsiiProxy_GroupMetric struct {
	_ byte // padding
}

func (j *jsiiProxy_GroupMetric) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}


// Experimental.
func NewGroupMetric(name *string) GroupMetric {
	_init_.Initialize()

	j := jsiiProxy_GroupMetric{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.GroupMetric",
		[]interface{}{name},
		&j,
	)

	return &j
}

// Experimental.
func NewGroupMetric_Override(g GroupMetric, name *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.GroupMetric",
		[]interface{}{name},
		g,
	)
}

func GroupMetric_DESIRED_CAPACITY() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"DESIRED_CAPACITY",
		&returns,
	)
	return returns
}

func GroupMetric_IN_SERVICE_INSTANCES() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"IN_SERVICE_INSTANCES",
		&returns,
	)
	return returns
}

func GroupMetric_MAX_SIZE() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"MAX_SIZE",
		&returns,
	)
	return returns
}

func GroupMetric_MIN_SIZE() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"MIN_SIZE",
		&returns,
	)
	return returns
}

func GroupMetric_PENDING_INSTANCES() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"PENDING_INSTANCES",
		&returns,
	)
	return returns
}

func GroupMetric_STANDBY_INSTANCES() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"STANDBY_INSTANCES",
		&returns,
	)
	return returns
}

func GroupMetric_TERMINATING_INSTANCES() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"TERMINATING_INSTANCES",
		&returns,
	)
	return returns
}

func GroupMetric_TOTAL_INSTANCES() GroupMetric {
	_init_.Initialize()
	var returns GroupMetric
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.GroupMetric",
		"TOTAL_INSTANCES",
		&returns,
	)
	return returns
}

// A set of group metrics.
// Experimental.
type GroupMetrics interface {
}

// The jsii proxy struct for GroupMetrics
type jsiiProxy_GroupMetrics struct {
	_ byte // padding
}

// Experimental.
func NewGroupMetrics(metrics ...GroupMetric) GroupMetrics {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range metrics {
		args = append(args, a)
	}

	j := jsiiProxy_GroupMetrics{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.GroupMetrics",
		args,
		&j,
	)

	return &j
}

// Experimental.
func NewGroupMetrics_Override(g GroupMetrics, metrics ...GroupMetric) {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range metrics {
		args = append(args, a)
	}

	_jsii_.Create(
		"monocdk.aws_autoscaling.GroupMetrics",
		args,
		g,
	)
}

// Report all group metrics.
// Experimental.
func GroupMetrics_All() GroupMetrics {
	_init_.Initialize()

	var returns GroupMetrics

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.GroupMetrics",
		"all",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Health check settings.
// Experimental.
type HealthCheck interface {
	GracePeriod() awscdk.Duration
	Type() *string
}

// The jsii proxy struct for HealthCheck
type jsiiProxy_HealthCheck struct {
	_ byte // padding
}

func (j *jsiiProxy_HealthCheck) GracePeriod() awscdk.Duration {
	var returns awscdk.Duration
	_jsii_.Get(
		j,
		"gracePeriod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_HealthCheck) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Use EC2 for health checks.
// Experimental.
func HealthCheck_Ec2(options *Ec2HealthCheckOptions) HealthCheck {
	_init_.Initialize()

	var returns HealthCheck

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.HealthCheck",
		"ec2",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Use ELB for health checks.
//
// It considers the instance unhealthy if it fails either the EC2 status checks or the load balancer health checks.
// Experimental.
func HealthCheck_Elb(options *ElbHealthCheckOptions) HealthCheck {
	_init_.Initialize()

	var returns HealthCheck

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.HealthCheck",
		"elb",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// An AutoScalingGroup.
// Experimental.
type IAutoScalingGroup interface {
	awsiam.IGrantable
	awscdk.IResource
	// Send a message to either an SQS queue or SNS topic when instances launch or terminate.
	// Experimental.
	AddLifecycleHook(id *string, props *BasicLifecycleHookProps) LifecycleHook
	// Add command to the startup script of fleet instances.
	//
	// The command must be in the scripting language supported by the fleet's OS (i.e. Linux/Windows).
	// Does nothing for imported ASGs.
	// Experimental.
	AddUserData(commands ...*string)
	// Scale out or in to achieve a target CPU utilization.
	// Experimental.
	ScaleOnCpuUtilization(id *string, props *CpuUtilizationScalingProps) TargetTrackingScalingPolicy
	// Scale out or in to achieve a target network ingress rate.
	// Experimental.
	ScaleOnIncomingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy
	// Scale out or in, in response to a metric.
	// Experimental.
	ScaleOnMetric(id *string, props *BasicStepScalingPolicyProps) StepScalingPolicy
	// Scale out or in to achieve a target network egress rate.
	// Experimental.
	ScaleOnOutgoingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy
	// Scale out or in based on time.
	// Experimental.
	ScaleOnSchedule(id *string, props *BasicScheduledActionProps) ScheduledAction
	// Scale out or in in order to keep a metric around a target value.
	// Experimental.
	ScaleToTrackMetric(id *string, props *MetricTargetTrackingProps) TargetTrackingScalingPolicy
	// The arn of the AutoScalingGroup.
	// Experimental.
	AutoScalingGroupArn() *string
	// The name of the AutoScalingGroup.
	// Experimental.
	AutoScalingGroupName() *string
	// The operating system family that the instances in this auto-scaling group belong to.
	//
	// Is 'UNKNOWN' for imported ASGs.
	// Experimental.
	OsType() awsec2.OperatingSystemType
}

// The jsii proxy for IAutoScalingGroup
type jsiiProxy_IAutoScalingGroup struct {
	internal.Type__awsiamIGrantable
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IAutoScalingGroup) AddLifecycleHook(id *string, props *BasicLifecycleHookProps) LifecycleHook {
	var returns LifecycleHook

	_jsii_.Invoke(
		i,
		"addLifecycleHook",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IAutoScalingGroup) AddUserData(commands ...*string) {
	args := []interface{}{}
	for _, a := range commands {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		i,
		"addUserData",
		args,
	)
}

func (i *jsiiProxy_IAutoScalingGroup) ScaleOnCpuUtilization(id *string, props *CpuUtilizationScalingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		i,
		"scaleOnCpuUtilization",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IAutoScalingGroup) ScaleOnIncomingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		i,
		"scaleOnIncomingBytes",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IAutoScalingGroup) ScaleOnMetric(id *string, props *BasicStepScalingPolicyProps) StepScalingPolicy {
	var returns StepScalingPolicy

	_jsii_.Invoke(
		i,
		"scaleOnMetric",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IAutoScalingGroup) ScaleOnOutgoingBytes(id *string, props *NetworkUtilizationScalingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		i,
		"scaleOnOutgoingBytes",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IAutoScalingGroup) ScaleOnSchedule(id *string, props *BasicScheduledActionProps) ScheduledAction {
	var returns ScheduledAction

	_jsii_.Invoke(
		i,
		"scaleOnSchedule",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IAutoScalingGroup) ScaleToTrackMetric(id *string, props *MetricTargetTrackingProps) TargetTrackingScalingPolicy {
	var returns TargetTrackingScalingPolicy

	_jsii_.Invoke(
		i,
		"scaleToTrackMetric",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IAutoScalingGroup) AutoScalingGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAutoScalingGroup) AutoScalingGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"autoScalingGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAutoScalingGroup) OsType() awsec2.OperatingSystemType {
	var returns awsec2.OperatingSystemType
	_jsii_.Get(
		j,
		"osType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAutoScalingGroup) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAutoScalingGroup) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAutoScalingGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAutoScalingGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// A basic lifecycle hook object.
// Experimental.
type ILifecycleHook interface {
	awscdk.IResource
	// The role for the lifecycle hook to execute.
	// Experimental.
	Role() awsiam.IRole
}

// The jsii proxy for ILifecycleHook
type jsiiProxy_ILifecycleHook struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_ILifecycleHook) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

// Interface for autoscaling lifecycle hook targets.
// Experimental.
type ILifecycleHookTarget interface {
	// Called when this object is used as the target of a lifecycle hook.
	// Experimental.
	Bind(scope awscdk.Construct, lifecycleHook ILifecycleHook) *LifecycleHookTargetConfig
}

// The jsii proxy for ILifecycleHookTarget
type jsiiProxy_ILifecycleHookTarget struct {
	_ byte // padding
}

func (i *jsiiProxy_ILifecycleHookTarget) Bind(scope awscdk.Construct, lifecycleHook ILifecycleHook) *LifecycleHookTargetConfig {
	var returns *LifecycleHookTargetConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, lifecycleHook},
		&returns,
	)

	return returns
}

// Define a life cycle hook.
// Experimental.
type LifecycleHook interface {
	awscdk.Resource
	ILifecycleHook
	Env() *awscdk.ResourceEnvironment
	LifecycleHookName() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Role() awsiam.IRole
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

// The jsii proxy struct for LifecycleHook
type jsiiProxy_LifecycleHook struct {
	internal.Type__awscdkResource
	jsiiProxy_ILifecycleHook
}

func (j *jsiiProxy_LifecycleHook) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LifecycleHook) LifecycleHookName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"lifecycleHookName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LifecycleHook) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LifecycleHook) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LifecycleHook) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LifecycleHook) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewLifecycleHook(scope constructs.Construct, id *string, props *LifecycleHookProps) LifecycleHook {
	_init_.Initialize()

	j := jsiiProxy_LifecycleHook{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.LifecycleHook",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewLifecycleHook_Override(l LifecycleHook, scope constructs.Construct, id *string, props *LifecycleHookProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.LifecycleHook",
		[]interface{}{scope, id, props},
		l,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func LifecycleHook_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.LifecycleHook",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func LifecycleHook_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.LifecycleHook",
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
func (l *jsiiProxy_LifecycleHook) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		l,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (l *jsiiProxy_LifecycleHook) GeneratePhysicalName() *string {
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
func (l *jsiiProxy_LifecycleHook) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (l *jsiiProxy_LifecycleHook) GetResourceNameAttribute(nameAttr *string) *string {
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
func (l *jsiiProxy_LifecycleHook) OnPrepare() {
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
func (l *jsiiProxy_LifecycleHook) OnSynthesize(session constructs.ISynthesisSession) {
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
func (l *jsiiProxy_LifecycleHook) OnValidate() *[]*string {
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
func (l *jsiiProxy_LifecycleHook) Prepare() {
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
func (l *jsiiProxy_LifecycleHook) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (l *jsiiProxy_LifecycleHook) ToString() *string {
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
func (l *jsiiProxy_LifecycleHook) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a Lifecycle hook.
// Experimental.
type LifecycleHookProps struct {
	// The state of the Amazon EC2 instance to which you want to attach the lifecycle hook.
	// Experimental.
	LifecycleTransition LifecycleTransition `json:"lifecycleTransition"`
	// The target of the lifecycle hook.
	// Experimental.
	NotificationTarget ILifecycleHookTarget `json:"notificationTarget"`
	// The action the Auto Scaling group takes when the lifecycle hook timeout elapses or if an unexpected failure occurs.
	// Experimental.
	DefaultResult DefaultResult `json:"defaultResult"`
	// Maximum time between calls to RecordLifecycleActionHeartbeat for the hook.
	//
	// If the lifecycle hook times out, perform the action in DefaultResult.
	// Experimental.
	HeartbeatTimeout awscdk.Duration `json:"heartbeatTimeout"`
	// Name of the lifecycle hook.
	// Experimental.
	LifecycleHookName *string `json:"lifecycleHookName"`
	// Additional data to pass to the lifecycle hook target.
	// Experimental.
	NotificationMetadata *string `json:"notificationMetadata"`
	// The role that allows publishing to the notification target.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// The AutoScalingGroup to add the lifecycle hook to.
	// Experimental.
	AutoScalingGroup IAutoScalingGroup `json:"autoScalingGroup"`
}

// Properties to add the target to a lifecycle hook.
// Experimental.
type LifecycleHookTargetConfig struct {
	// The ARN to use as the notification target.
	// Experimental.
	NotificationTargetArn *string `json:"notificationTargetArn"`
}

// What instance transition to attach the hook to.
// Experimental.
type LifecycleTransition string

const (
	LifecycleTransition_INSTANCE_LAUNCHING LifecycleTransition = "INSTANCE_LAUNCHING"
	LifecycleTransition_INSTANCE_TERMINATING LifecycleTransition = "INSTANCE_TERMINATING"
)

// How the scaling metric is going to be aggregated.
// Experimental.
type MetricAggregationType string

const (
	MetricAggregationType_AVERAGE MetricAggregationType = "AVERAGE"
	MetricAggregationType_MINIMUM MetricAggregationType = "MINIMUM"
	MetricAggregationType_MAXIMUM MetricAggregationType = "MAXIMUM"
)

// Properties for enabling tracking of an arbitrary metric.
// Experimental.
type MetricTargetTrackingProps struct {
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the autoscaling group. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// group.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// Metric to track.
	//
	// The metric must represent a utilization, so that if it's higher than the
	// target value, your ASG should scale out, and if it's lower it should
	// scale in.
	// Experimental.
	Metric awscloudwatch.IMetric `json:"metric"`
	// Value to keep the metric around.
	// Experimental.
	TargetValue *float64 `json:"targetValue"`
}

// The monitoring mode for instances launched in an autoscaling group.
// Experimental.
type Monitoring string

const (
	Monitoring_BASIC Monitoring = "BASIC"
	Monitoring_DETAILED Monitoring = "DETAILED"
)

// Properties for enabling scaling based on network utilization.
// Experimental.
type NetworkUtilizationScalingProps struct {
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the autoscaling group. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// group.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// Target average bytes/seconds on each instance.
	// Experimental.
	TargetBytesPerSecond *float64 `json:"targetBytesPerSecond"`
}

// AutoScalingGroup fleet change notifications configurations.
//
// You can configure AutoScaling to send an SNS notification whenever your Auto Scaling group scales.
// Experimental.
type NotificationConfiguration struct {
	// SNS topic to send notifications about fleet scaling events.
	// Experimental.
	Topic awssns.ITopic `json:"topic"`
	// Which fleet scaling events triggers a notification.
	// Experimental.
	ScalingEvents ScalingEvents `json:"scalingEvents"`
}

// One of the predefined autoscaling metrics.
// Experimental.
type PredefinedMetric string

const (
	PredefinedMetric_ASG_AVERAGE_CPU_UTILIZATION PredefinedMetric = "ASG_AVERAGE_CPU_UTILIZATION"
	PredefinedMetric_ASG_AVERAGE_NETWORK_IN PredefinedMetric = "ASG_AVERAGE_NETWORK_IN"
	PredefinedMetric_ASG_AVERAGE_NETWORK_OUT PredefinedMetric = "ASG_AVERAGE_NETWORK_OUT"
	PredefinedMetric_ALB_REQUEST_COUNT_PER_TARGET PredefinedMetric = "ALB_REQUEST_COUNT_PER_TARGET"
)

// Input for Signals.renderCreationPolicy.
// Experimental.
type RenderSignalsOptions struct {
	// The desiredCapacity of the ASG.
	// Experimental.
	DesiredCapacity *float64 `json:"desiredCapacity"`
	// The minSize of the ASG.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
}

// Properties for enabling scaling based on request/second.
// Experimental.
type RequestCountScalingProps struct {
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the autoscaling group. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// group.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// Target average requests/minute on each instance.
	// Experimental.
	TargetRequestsPerMinute *float64 `json:"targetRequestsPerMinute"`
	// Target average requests/seconds on each instance.
	// Deprecated: Use 'targetRequestsPerMinute' instead
	TargetRequestsPerSecond *float64 `json:"targetRequestsPerSecond"`
}

// Additional settings when a rolling update is selected.
// Deprecated: use `UpdatePolicy.rollingUpdate()`
type RollingUpdateConfiguration struct {
	// The maximum number of instances that AWS CloudFormation updates at once.
	// Deprecated: use `UpdatePolicy.rollingUpdate()`
	MaxBatchSize *float64 `json:"maxBatchSize"`
	// The minimum number of instances that must be in service before more instances are replaced.
	//
	// This number affects the speed of the replacement.
	// Deprecated: use `UpdatePolicy.rollingUpdate()`
	MinInstancesInService *float64 `json:"minInstancesInService"`
	// The percentage of instances that must signal success for an update to succeed.
	//
	// If an instance doesn't send a signal within the time specified in the
	// pauseTime property, AWS CloudFormation assumes that the instance wasn't
	// updated.
	//
	// This number affects the success of the replacement.
	//
	// If you specify this property, you must also enable the
	// waitOnResourceSignals and pauseTime properties.
	// Deprecated: use `UpdatePolicy.rollingUpdate()`
	MinSuccessfulInstancesPercent *float64 `json:"minSuccessfulInstancesPercent"`
	// The pause time after making a change to a batch of instances.
	//
	// This is intended to give those instances time to start software applications.
	//
	// Specify PauseTime in the ISO8601 duration format (in the format
	// PT#H#M#S, where each # is the number of hours, minutes, and seconds,
	// respectively). The maximum PauseTime is one hour (PT1H).
	// Deprecated: use `UpdatePolicy.rollingUpdate()`
	PauseTime awscdk.Duration `json:"pauseTime"`
	// Specifies the Auto Scaling processes to suspend during a stack update.
	//
	// Suspending processes prevents Auto Scaling from interfering with a stack
	// update.
	// Deprecated: use `UpdatePolicy.rollingUpdate()`
	SuspendProcesses *[]ScalingProcess `json:"suspendProcesses"`
	// Specifies whether the Auto Scaling group waits on signals from new instances during an update.
	//
	// AWS CloudFormation must receive a signal from each new instance within
	// the specified PauseTime before continuing the update.
	//
	// To have instances wait for an Elastic Load Balancing health check before
	// they signal success, add a health-check verification by using the
	// cfn-init helper script. For an example, see the verify_instance_health
	// command in the Auto Scaling rolling updates sample template.
	// Deprecated: use `UpdatePolicy.rollingUpdate()`
	WaitOnResourceSignals *bool `json:"waitOnResourceSignals"`
}

// Options for customizing the rolling update.
// Experimental.
type RollingUpdateOptions struct {
	// The maximum number of instances that AWS CloudFormation updates at once.
	//
	// This number affects the speed of the replacement.
	// Experimental.
	MaxBatchSize *float64 `json:"maxBatchSize"`
	// The minimum number of instances that must be in service before more instances are replaced.
	//
	// This number affects the speed of the replacement.
	// Experimental.
	MinInstancesInService *float64 `json:"minInstancesInService"`
	// The percentage of instances that must signal success for the update to succeed.
	// Experimental.
	MinSuccessPercentage *float64 `json:"minSuccessPercentage"`
	// The pause time after making a change to a batch of instances.
	// Experimental.
	PauseTime awscdk.Duration `json:"pauseTime"`
	// Specifies the Auto Scaling processes to suspend during a stack update.
	//
	// Suspending processes prevents Auto Scaling from interfering with a stack
	// update.
	// Experimental.
	SuspendProcesses *[]ScalingProcess `json:"suspendProcesses"`
	// Specifies whether the Auto Scaling group waits on signals from new instances during an update.
	// Experimental.
	WaitOnResourceSignals *bool `json:"waitOnResourceSignals"`
}

// Fleet scaling events.
// Experimental.
type ScalingEvent string

const (
	ScalingEvent_INSTANCE_LAUNCH ScalingEvent = "INSTANCE_LAUNCH"
	ScalingEvent_INSTANCE_TERMINATE ScalingEvent = "INSTANCE_TERMINATE"
	ScalingEvent_INSTANCE_TERMINATE_ERROR ScalingEvent = "INSTANCE_TERMINATE_ERROR"
	ScalingEvent_INSTANCE_LAUNCH_ERROR ScalingEvent = "INSTANCE_LAUNCH_ERROR"
	ScalingEvent_TEST_NOTIFICATION ScalingEvent = "TEST_NOTIFICATION"
)

// A list of ScalingEvents, you can use one of the predefined lists, such as ScalingEvents.ERRORS or create a custom group by instantiating a `NotificationTypes` object, e.g: `new NotificationTypes(`NotificationType.INSTANCE_LAUNCH`)`.
// Experimental.
type ScalingEvents interface {
}

// The jsii proxy struct for ScalingEvents
type jsiiProxy_ScalingEvents struct {
	_ byte // padding
}

// Experimental.
func NewScalingEvents(types ...ScalingEvent) ScalingEvents {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range types {
		args = append(args, a)
	}

	j := jsiiProxy_ScalingEvents{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.ScalingEvents",
		args,
		&j,
	)

	return &j
}

// Experimental.
func NewScalingEvents_Override(s ScalingEvents, types ...ScalingEvent) {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range types {
		args = append(args, a)
	}

	_jsii_.Create(
		"monocdk.aws_autoscaling.ScalingEvents",
		args,
		s,
	)
}

func ScalingEvents_ALL() ScalingEvents {
	_init_.Initialize()
	var returns ScalingEvents
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.ScalingEvents",
		"ALL",
		&returns,
	)
	return returns
}

func ScalingEvents_ERRORS() ScalingEvents {
	_init_.Initialize()
	var returns ScalingEvents
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.ScalingEvents",
		"ERRORS",
		&returns,
	)
	return returns
}

func ScalingEvents_LAUNCH_EVENTS() ScalingEvents {
	_init_.Initialize()
	var returns ScalingEvents
	_jsii_.StaticGet(
		"monocdk.aws_autoscaling.ScalingEvents",
		"LAUNCH_EVENTS",
		&returns,
	)
	return returns
}

// A range of metric values in which to apply a certain scaling operation.
// Experimental.
type ScalingInterval struct {
	// The capacity adjustment to apply in this interval.
	//
	// The number is interpreted differently based on AdjustmentType:
	//
	// - ChangeInCapacity: add the adjustment to the current capacity.
	//   The number can be positive or negative.
	// - PercentChangeInCapacity: add or remove the given percentage of the current
	//    capacity to itself. The number can be in the range [-100..100].
	// - ExactCapacity: set the capacity to this number. The number must
	//    be positive.
	// Experimental.
	Change *float64 `json:"change"`
	// The lower bound of the interval.
	//
	// The scaling adjustment will be applied if the metric is higher than this value.
	// Experimental.
	Lower *float64 `json:"lower"`
	// The upper bound of the interval.
	//
	// The scaling adjustment will be applied if the metric is lower than this value.
	// Experimental.
	Upper *float64 `json:"upper"`
}

// Experimental.
type ScalingProcess string

const (
	ScalingProcess_LAUNCH ScalingProcess = "LAUNCH"
	ScalingProcess_TERMINATE ScalingProcess = "TERMINATE"
	ScalingProcess_HEALTH_CHECK ScalingProcess = "HEALTH_CHECK"
	ScalingProcess_REPLACE_UNHEALTHY ScalingProcess = "REPLACE_UNHEALTHY"
	ScalingProcess_AZ_REBALANCE ScalingProcess = "AZ_REBALANCE"
	ScalingProcess_ALARM_NOTIFICATION ScalingProcess = "ALARM_NOTIFICATION"
	ScalingProcess_SCHEDULED_ACTIONS ScalingProcess = "SCHEDULED_ACTIONS"
	ScalingProcess_ADD_TO_LOAD_BALANCER ScalingProcess = "ADD_TO_LOAD_BALANCER"
)

// Schedule for scheduled scaling actions.
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
		"monocdk.aws_autoscaling.Schedule",
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
		"monocdk.aws_autoscaling.Schedule",
		"cron",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Construct a schedule from a literal schedule expression.
// See: http://crontab.org/
//
// Experimental.
func Schedule_Expression(expression *string) Schedule {
	_init_.Initialize()

	var returns Schedule

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.Schedule",
		"expression",
		[]interface{}{expression},
		&returns,
	)

	return returns
}

// Define a scheduled scaling action.
// Experimental.
type ScheduledAction interface {
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

// The jsii proxy struct for ScheduledAction
type jsiiProxy_ScheduledAction struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_ScheduledAction) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScheduledAction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScheduledAction) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScheduledAction) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewScheduledAction(scope constructs.Construct, id *string, props *ScheduledActionProps) ScheduledAction {
	_init_.Initialize()

	j := jsiiProxy_ScheduledAction{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.ScheduledAction",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewScheduledAction_Override(s ScheduledAction, scope constructs.Construct, id *string, props *ScheduledActionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.ScheduledAction",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func ScheduledAction_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.ScheduledAction",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ScheduledAction_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.ScheduledAction",
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
func (s *jsiiProxy_ScheduledAction) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (s *jsiiProxy_ScheduledAction) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_ScheduledAction) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_ScheduledAction) GetResourceNameAttribute(nameAttr *string) *string {
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
func (s *jsiiProxy_ScheduledAction) OnPrepare() {
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
func (s *jsiiProxy_ScheduledAction) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_ScheduledAction) OnValidate() *[]*string {
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
func (s *jsiiProxy_ScheduledAction) Prepare() {
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
func (s *jsiiProxy_ScheduledAction) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_ScheduledAction) ToString() *string {
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
func (s *jsiiProxy_ScheduledAction) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a scheduled action on an AutoScalingGroup.
// Experimental.
type ScheduledActionProps struct {
	// When to perform this action.
	//
	// Supports cron expressions.
	//
	// For more information about cron expressions, see https://en.wikipedia.org/wiki/Cron.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Schedule Schedule `json:"schedule"`
	// The new desired capacity.
	//
	// At the scheduled time, set the desired capacity to the given capacity.
	//
	// At least one of maxCapacity, minCapacity, or desiredCapacity must be supplied.
	// Experimental.
	DesiredCapacity *float64 `json:"desiredCapacity"`
	// When this scheduled action expires.
	// Experimental.
	EndTime *time.Time `json:"endTime"`
	// The new maximum capacity.
	//
	// At the scheduled time, set the maximum capacity to the given capacity.
	//
	// At least one of maxCapacity, minCapacity, or desiredCapacity must be supplied.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// The new minimum capacity.
	//
	// At the scheduled time, set the minimum capacity to the given capacity.
	//
	// At least one of maxCapacity, minCapacity, or desiredCapacity must be supplied.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// When this scheduled action becomes active.
	// Experimental.
	StartTime *time.Time `json:"startTime"`
	// The AutoScalingGroup to apply the scheduled actions to.
	// Experimental.
	AutoScalingGroup IAutoScalingGroup `json:"autoScalingGroup"`
}

// Configure whether the AutoScalingGroup waits for signals.
//
// If you do configure waiting for signals, you should make sure the instances
// invoke `cfn-signal` somewhere in their UserData to signal that they have
// started up (either successfully or unsuccessfully).
//
// Signals are used both during intial creation and subsequent updates.
// Experimental.
type Signals interface {
	DoRender(options *SignalsOptions, count *float64) *awscdk.CfnCreationPolicy
	RenderCreationPolicy(renderOptions *RenderSignalsOptions) *awscdk.CfnCreationPolicy
}

// The jsii proxy struct for Signals
type jsiiProxy_Signals struct {
	_ byte // padding
}

// Experimental.
func NewSignals_Override(s Signals) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.Signals",
		nil, // no parameters
		s,
	)
}

// Wait for the desiredCapacity of the AutoScalingGroup amount of signals to have been received.
//
// If no desiredCapacity has been configured, wait for minCapacity signals intead.
//
// This number is used during initial creation and during replacing updates.
// During rolling updates, all updated instances must send a signal.
// Experimental.
func Signals_WaitForAll(options *SignalsOptions) Signals {
	_init_.Initialize()

	var returns Signals

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.Signals",
		"waitForAll",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Wait for a specific amount of signals to have been received.
//
// You should send one signal per instance, so this represents the number of
// instances to wait for.
//
// This number is used during initial creation and during replacing updates.
// During rolling updates, all updated instances must send a signal.
// Experimental.
func Signals_WaitForCount(count *float64, options *SignalsOptions) Signals {
	_init_.Initialize()

	var returns Signals

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.Signals",
		"waitForCount",
		[]interface{}{count, options},
		&returns,
	)

	return returns
}

// Wait for the minCapacity of the AutoScalingGroup amount of signals to have been received.
//
// This number is used during initial creation and during replacing updates.
// During rolling updates, all updated instances must send a signal.
// Experimental.
func Signals_WaitForMinCapacity(options *SignalsOptions) Signals {
	_init_.Initialize()

	var returns Signals

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.Signals",
		"waitForMinCapacity",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Helper to render the actual creation policy, as the logic between them is quite similar.
// Experimental.
func (s *jsiiProxy_Signals) DoRender(options *SignalsOptions, count *float64) *awscdk.CfnCreationPolicy {
	var returns *awscdk.CfnCreationPolicy

	_jsii_.Invoke(
		s,
		"doRender",
		[]interface{}{options, count},
		&returns,
	)

	return returns
}

// Render the ASG's CreationPolicy.
// Experimental.
func (s *jsiiProxy_Signals) RenderCreationPolicy(renderOptions *RenderSignalsOptions) *awscdk.CfnCreationPolicy {
	var returns *awscdk.CfnCreationPolicy

	_jsii_.Invoke(
		s,
		"renderCreationPolicy",
		[]interface{}{renderOptions},
		&returns,
	)

	return returns
}

// Customization options for Signal handling.
// Experimental.
type SignalsOptions struct {
	// The percentage of signals that need to be successful.
	//
	// If this number is less than 100, a percentage of signals may be failure
	// signals while still succeeding the creation or update in CloudFormation.
	// Experimental.
	MinSuccessPercentage *float64 `json:"minSuccessPercentage"`
	// How long to wait for the signals to be sent.
	//
	// This should reflect how long it takes your instances to start up
	// (including instance start time and instance initialization time).
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
}

// Define a step scaling action.
//
// This kind of scaling policy adjusts the target capacity in configurable
// steps. The size of the step is configurable based on the metric's distance
// to its alarm threshold.
//
// This Action must be used as the target of a CloudWatch alarm to take effect.
// Experimental.
type StepScalingAction interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	ScalingPolicyArn() *string
	AddAdjustment(adjustment *AdjustmentTier)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for StepScalingAction
type jsiiProxy_StepScalingAction struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_StepScalingAction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingAction) ScalingPolicyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalingPolicyArn",
		&returns,
	)
	return returns
}


// Experimental.
func NewStepScalingAction(scope constructs.Construct, id *string, props *StepScalingActionProps) StepScalingAction {
	_init_.Initialize()

	j := jsiiProxy_StepScalingAction{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.StepScalingAction",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStepScalingAction_Override(s StepScalingAction, scope constructs.Construct, id *string, props *StepScalingActionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.StepScalingAction",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func StepScalingAction_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.StepScalingAction",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add an adjusment interval to the ScalingAction.
// Experimental.
func (s *jsiiProxy_StepScalingAction) AddAdjustment(adjustment *AdjustmentTier) {
	_jsii_.InvokeVoid(
		s,
		"addAdjustment",
		[]interface{}{adjustment},
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
func (s *jsiiProxy_StepScalingAction) OnPrepare() {
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
func (s *jsiiProxy_StepScalingAction) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_StepScalingAction) OnValidate() *[]*string {
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
func (s *jsiiProxy_StepScalingAction) Prepare() {
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
func (s *jsiiProxy_StepScalingAction) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_StepScalingAction) ToString() *string {
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
func (s *jsiiProxy_StepScalingAction) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a scaling policy.
// Experimental.
type StepScalingActionProps struct {
	// The auto scaling group.
	// Experimental.
	AutoScalingGroup IAutoScalingGroup `json:"autoScalingGroup"`
	// How the adjustment numbers are interpreted.
	// Experimental.
	AdjustmentType AdjustmentType `json:"adjustmentType"`
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// The aggregation type for the CloudWatch metrics.
	// Experimental.
	MetricAggregationType MetricAggregationType `json:"metricAggregationType"`
	// Minimum absolute number to adjust capacity with as result of percentage scaling.
	//
	// Only when using AdjustmentType = PercentChangeInCapacity, this number controls
	// the minimum absolute effect size.
	// Experimental.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
}

// Define a acaling strategy which scales depending on absolute values of some metric.
//
// You can specify the scaling behavior for various values of the metric.
//
// Implemented using one or more CloudWatch alarms and Step Scaling Policies.
// Experimental.
type StepScalingPolicy interface {
	awscdk.Construct
	LowerAction() StepScalingAction
	LowerAlarm() awscloudwatch.Alarm
	Node() awscdk.ConstructNode
	UpperAction() StepScalingAction
	UpperAlarm() awscloudwatch.Alarm
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for StepScalingPolicy
type jsiiProxy_StepScalingPolicy struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_StepScalingPolicy) LowerAction() StepScalingAction {
	var returns StepScalingAction
	_jsii_.Get(
		j,
		"lowerAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) LowerAlarm() awscloudwatch.Alarm {
	var returns awscloudwatch.Alarm
	_jsii_.Get(
		j,
		"lowerAlarm",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) UpperAction() StepScalingAction {
	var returns StepScalingAction
	_jsii_.Get(
		j,
		"upperAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_StepScalingPolicy) UpperAlarm() awscloudwatch.Alarm {
	var returns awscloudwatch.Alarm
	_jsii_.Get(
		j,
		"upperAlarm",
		&returns,
	)
	return returns
}


// Experimental.
func NewStepScalingPolicy(scope constructs.Construct, id *string, props *StepScalingPolicyProps) StepScalingPolicy {
	_init_.Initialize()

	j := jsiiProxy_StepScalingPolicy{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.StepScalingPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStepScalingPolicy_Override(s StepScalingPolicy, scope constructs.Construct, id *string, props *StepScalingPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.StepScalingPolicy",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func StepScalingPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.StepScalingPolicy",
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
func (s *jsiiProxy_StepScalingPolicy) OnPrepare() {
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
func (s *jsiiProxy_StepScalingPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_StepScalingPolicy) OnValidate() *[]*string {
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
func (s *jsiiProxy_StepScalingPolicy) Prepare() {
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
func (s *jsiiProxy_StepScalingPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_StepScalingPolicy) ToString() *string {
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
func (s *jsiiProxy_StepScalingPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type StepScalingPolicyProps struct {
	// Metric to scale on.
	// Experimental.
	Metric awscloudwatch.IMetric `json:"metric"`
	// The intervals for scaling.
	//
	// Maps a range of metric values to a particular scaling behavior.
	// Experimental.
	ScalingSteps *[]*ScalingInterval `json:"scalingSteps"`
	// How the adjustment numbers inside 'intervals' are interpreted.
	// Experimental.
	AdjustmentType AdjustmentType `json:"adjustmentType"`
	// Grace period after scaling activity.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// How many evaluation periods of the metric to wait before triggering a scaling action.
	//
	// Raising this value can be used to smooth out the metric, at the expense
	// of slower response times.
	// Experimental.
	EvaluationPeriods *float64 `json:"evaluationPeriods"`
	// Aggregation to apply to all data points over the evaluation periods.
	//
	// Only has meaning if `evaluationPeriods != 1`.
	// Experimental.
	MetricAggregationType MetricAggregationType `json:"metricAggregationType"`
	// Minimum absolute number to adjust capacity with as result of percentage scaling.
	//
	// Only when using AdjustmentType = PercentChangeInCapacity, this number controls
	// the minimum absolute effect size.
	// Experimental.
	MinAdjustmentMagnitude *float64 `json:"minAdjustmentMagnitude"`
	// The auto scaling group.
	// Experimental.
	AutoScalingGroup IAutoScalingGroup `json:"autoScalingGroup"`
}

// Experimental.
type TargetTrackingScalingPolicy interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	ScalingPolicyArn() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for TargetTrackingScalingPolicy
type jsiiProxy_TargetTrackingScalingPolicy struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_TargetTrackingScalingPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetTrackingScalingPolicy) ScalingPolicyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scalingPolicyArn",
		&returns,
	)
	return returns
}


// Experimental.
func NewTargetTrackingScalingPolicy(scope constructs.Construct, id *string, props *TargetTrackingScalingPolicyProps) TargetTrackingScalingPolicy {
	_init_.Initialize()

	j := jsiiProxy_TargetTrackingScalingPolicy{}

	_jsii_.Create(
		"monocdk.aws_autoscaling.TargetTrackingScalingPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewTargetTrackingScalingPolicy_Override(t TargetTrackingScalingPolicy, scope constructs.Construct, id *string, props *TargetTrackingScalingPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.TargetTrackingScalingPolicy",
		[]interface{}{scope, id, props},
		t,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func TargetTrackingScalingPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.TargetTrackingScalingPolicy",
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
func (t *jsiiProxy_TargetTrackingScalingPolicy) OnPrepare() {
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
func (t *jsiiProxy_TargetTrackingScalingPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (t *jsiiProxy_TargetTrackingScalingPolicy) OnValidate() *[]*string {
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
func (t *jsiiProxy_TargetTrackingScalingPolicy) Prepare() {
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
func (t *jsiiProxy_TargetTrackingScalingPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TargetTrackingScalingPolicy) ToString() *string {
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
func (t *jsiiProxy_TargetTrackingScalingPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a concrete TargetTrackingPolicy.
//
// Adds the scalingTarget.
// Experimental.
type TargetTrackingScalingPolicyProps struct {
	// Period after a scaling completes before another scaling activity can start.
	// Experimental.
	Cooldown awscdk.Duration `json:"cooldown"`
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the autoscaling group. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// group.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// Estimated time until a newly launched instance can send metrics to CloudWatch.
	// Experimental.
	EstimatedInstanceWarmup awscdk.Duration `json:"estimatedInstanceWarmup"`
	// The target value for the metric.
	// Experimental.
	TargetValue *float64 `json:"targetValue"`
	// A custom metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	CustomMetric awscloudwatch.IMetric `json:"customMetric"`
	// A predefined metric for application autoscaling.
	//
	// The metric must track utilization. Scaling out will happen if the metric is higher than
	// the target value, scaling in will happen in the metric is lower than the target value.
	//
	// Exactly one of customMetric or predefinedMetric must be specified.
	// Experimental.
	PredefinedMetric PredefinedMetric `json:"predefinedMetric"`
	// The resource label associated with the predefined metric.
	//
	// Should be supplied if the predefined metric is ALBRequestCountPerTarget, and the
	// format should be:
	//
	// app/<load-balancer-name>/<load-balancer-id>/targetgroup/<target-group-name>/<target-group-id>
	// Experimental.
	ResourceLabel *string `json:"resourceLabel"`
	// Experimental.
	AutoScalingGroup IAutoScalingGroup `json:"autoScalingGroup"`
}

// How existing instances should be updated.
// Experimental.
type UpdatePolicy interface {
}

// The jsii proxy struct for UpdatePolicy
type jsiiProxy_UpdatePolicy struct {
	_ byte // padding
}

// Experimental.
func NewUpdatePolicy_Override(u UpdatePolicy) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_autoscaling.UpdatePolicy",
		nil, // no parameters
		u,
	)
}

// Create a new AutoScalingGroup and switch over to it.
// Experimental.
func UpdatePolicy_ReplacingUpdate() UpdatePolicy {
	_init_.Initialize()

	var returns UpdatePolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.UpdatePolicy",
		"replacingUpdate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Replace the instances in the AutoScalingGroup one by one, or in batches.
// Experimental.
func UpdatePolicy_RollingUpdate(options *RollingUpdateOptions) UpdatePolicy {
	_init_.Initialize()

	var returns UpdatePolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_autoscaling.UpdatePolicy",
		"rollingUpdate",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// The type of update to perform on instances in this AutoScalingGroup.
// Deprecated: Use UpdatePolicy instead
type UpdateType string

const (
	UpdateType_NONE UpdateType = "NONE"
	UpdateType_REPLACING_UPDATE UpdateType = "REPLACING_UPDATE"
	UpdateType_ROLLING_UPDATE UpdateType = "ROLLING_UPDATE"
)


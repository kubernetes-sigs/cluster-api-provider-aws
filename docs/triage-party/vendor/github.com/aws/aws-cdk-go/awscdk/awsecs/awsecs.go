package awsecs

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/assets"
	"github.com/aws/aws-cdk-go/awscdk/awsapplicationautoscaling"
	"github.com/aws/aws-cdk-go/awscdk/awsautoscaling"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/awsecs/internal"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancing"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/awssecretsmanager"
	"github.com/aws/aws-cdk-go/awscdk/awsservicediscovery"
	"github.com/aws/aws-cdk-go/awscdk/awssns"
	"github.com/aws/aws-cdk-go/awscdk/awsssm"
	"github.com/aws/constructs-go/constructs/v3"
)

// The properties for adding an AutoScalingGroup.
// Experimental.
type AddAutoScalingGroupCapacityOptions struct {
	// Specifies whether the containers can access the container instance role.
	// Experimental.
	CanContainersAccessInstanceRole *bool `json:"canContainersAccessInstanceRole"`
	// Specify the machine image type.
	// Experimental.
	MachineImageType MachineImageType `json:"machineImageType"`
	// Specify whether to enable Automated Draining for Spot Instances running Amazon ECS Services.
	//
	// For more information, see [Using Spot Instances](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/container-instance-spot.html).
	// Experimental.
	SpotInstanceDraining *bool `json:"spotInstanceDraining"`
	// The time period to wait before force terminating an instance that is draining.
	//
	// This creates a Lambda function that is used by a lifecycle hook for the
	// AutoScalingGroup that will delay instance termination until all ECS tasks
	// have drained from the instance. Set to 0 to disable task draining.
	//
	// Set to 0 to disable task draining.
	// Deprecated: The lifecycle draining hook is not configured if using the EC2 Capacity Provider. Enable managed termination protection instead.
	TaskDrainTime awscdk.Duration `json:"taskDrainTime"`
	// If {@link AddAutoScalingGroupCapacityOptions.taskDrainTime} is non-zero, then the ECS cluster creates an SNS Topic to as part of a system to drain instances of tasks when the instance is being shut down. If this property is provided, then this key will be used to encrypt the contents of that SNS Topic. See [SNS Data Encryption](https://docs.aws.amazon.com/sns/latest/dg/sns-data-encryption.html) for more information.
	// Experimental.
	TopicEncryptionKey awskms.IKey `json:"topicEncryptionKey"`
}

// The properties for adding instance capacity to an AutoScalingGroup.
// Experimental.
type AddCapacityOptions struct {
	// Specifies whether the containers can access the container instance role.
	// Experimental.
	CanContainersAccessInstanceRole *bool `json:"canContainersAccessInstanceRole"`
	// Specify the machine image type.
	// Experimental.
	MachineImageType MachineImageType `json:"machineImageType"`
	// Specify whether to enable Automated Draining for Spot Instances running Amazon ECS Services.
	//
	// For more information, see [Using Spot Instances](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/container-instance-spot.html).
	// Experimental.
	SpotInstanceDraining *bool `json:"spotInstanceDraining"`
	// The time period to wait before force terminating an instance that is draining.
	//
	// This creates a Lambda function that is used by a lifecycle hook for the
	// AutoScalingGroup that will delay instance termination until all ECS tasks
	// have drained from the instance. Set to 0 to disable task draining.
	//
	// Set to 0 to disable task draining.
	// Deprecated: The lifecycle draining hook is not configured if using the EC2 Capacity Provider. Enable managed termination protection instead.
	TaskDrainTime awscdk.Duration `json:"taskDrainTime"`
	// If {@link AddAutoScalingGroupCapacityOptions.taskDrainTime} is non-zero, then the ECS cluster creates an SNS Topic to as part of a system to drain instances of tasks when the instance is being shut down. If this property is provided, then this key will be used to encrypt the contents of that SNS Topic. See [SNS Data Encryption](https://docs.aws.amazon.com/sns/latest/dg/sns-data-encryption.html) for more information.
	// Experimental.
	TopicEncryptionKey awskms.IKey `json:"topicEncryptionKey"`
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
	BlockDevices *[]*awsautoscaling.BlockDevice `json:"blockDevices"`
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
	GroupMetrics *[]awsautoscaling.GroupMetrics `json:"groupMetrics"`
	// Configuration for health checks.
	// Experimental.
	HealthCheck awsautoscaling.HealthCheck `json:"healthCheck"`
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
	InstanceMonitoring awsautoscaling.Monitoring `json:"instanceMonitoring"`
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
	Notifications *[]*awsautoscaling.NotificationConfiguration `json:"notifications"`
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
	RollingUpdateConfiguration *awsautoscaling.RollingUpdateConfiguration `json:"rollingUpdateConfiguration"`
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
	Signals awsautoscaling.Signals `json:"signals"`
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
	UpdatePolicy awsautoscaling.UpdatePolicy `json:"updatePolicy"`
	// What to do when an AutoScalingGroup's instance configuration is changed.
	//
	// This is applied when any of the settings on the ASG are changed that
	// affect how the instances should be created (VPC, instance type, startup
	// scripts, etc.). It indicates how the existing instances should be
	// replaced with new instances matching the new config. By default, nothing
	// is done and only new instances are launched with the new config.
	// Deprecated: Use `updatePolicy` instead
	UpdateType awsautoscaling.UpdateType `json:"updateType"`
	// Where to place instances within the VPC.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// The EC2 instance type to use when launching instances into the AutoScalingGroup.
	// Experimental.
	InstanceType awsec2.InstanceType `json:"instanceType"`
	// The ECS-optimized AMI variant to use.
	//
	// For more information, see
	// [Amazon ECS-optimized AMIs](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-optimized_AMI.html).
	// You must define either `machineImage` or `machineImageType`, not both.
	// Experimental.
	MachineImage awsec2.IMachineImage `json:"machineImage"`
}

// The ECS-optimized AMI variant to use.
//
// For more information, see
// [Amazon ECS-optimized AMIs](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-optimized_AMI.html).
// Experimental.
type AmiHardwareType string

const (
	AmiHardwareType_STANDARD AmiHardwareType = "STANDARD"
	AmiHardwareType_GPU AmiHardwareType = "GPU"
	AmiHardwareType_ARM AmiHardwareType = "ARM"
)

// The class for App Mesh proxy configurations.
//
// For tasks using the EC2 launch type, the container instances require at least version 1.26.0 of the container agent and at least version
// 1.26.0-1 of the ecs-init package to enable a proxy configuration. If your container instances are launched from the Amazon ECS-optimized
// AMI version 20190301 or later, then they contain the required versions of the container agent and ecs-init.
// For more information, see [Amazon ECS-optimized AMIs](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-optimized_AMI.html).
//
// For tasks using the Fargate launch type, the task or service requires platform version 1.3.0 or later.
// Experimental.
type AppMeshProxyConfiguration interface {
	ProxyConfiguration
	Bind(_scope awscdk.Construct, _taskDefinition TaskDefinition) *CfnTaskDefinition_ProxyConfigurationProperty
}

// The jsii proxy struct for AppMeshProxyConfiguration
type jsiiProxy_AppMeshProxyConfiguration struct {
	jsiiProxy_ProxyConfiguration
}

// Constructs a new instance of the AppMeshProxyConfiguration class.
// Experimental.
func NewAppMeshProxyConfiguration(props *AppMeshProxyConfigurationConfigProps) AppMeshProxyConfiguration {
	_init_.Initialize()

	j := jsiiProxy_AppMeshProxyConfiguration{}

	_jsii_.Create(
		"monocdk.aws_ecs.AppMeshProxyConfiguration",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the AppMeshProxyConfiguration class.
// Experimental.
func NewAppMeshProxyConfiguration_Override(a AppMeshProxyConfiguration, props *AppMeshProxyConfigurationConfigProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.AppMeshProxyConfiguration",
		[]interface{}{props},
		a,
	)
}

// Called when the proxy configuration is configured on a task definition.
// Experimental.
func (a *jsiiProxy_AppMeshProxyConfiguration) Bind(_scope awscdk.Construct, _taskDefinition TaskDefinition) *CfnTaskDefinition_ProxyConfigurationProperty {
	var returns *CfnTaskDefinition_ProxyConfigurationProperty

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{_scope, _taskDefinition},
		&returns,
	)

	return returns
}

// The configuration to use when setting an App Mesh proxy configuration.
// Experimental.
type AppMeshProxyConfigurationConfigProps struct {
	// The name of the container that will serve as the App Mesh proxy.
	// Experimental.
	ContainerName *string `json:"containerName"`
	// The set of network configuration parameters to provide the Container Network Interface (CNI) plugin.
	// Experimental.
	Properties *AppMeshProxyConfigurationProps `json:"properties"`
}

// Interface for setting the properties of proxy configuration.
// Experimental.
type AppMeshProxyConfigurationProps struct {
	// The list of ports that the application uses.
	//
	// Network traffic to these ports is forwarded to the ProxyIngressPort and ProxyEgressPort.
	// Experimental.
	AppPorts *[]*float64 `json:"appPorts"`
	// Specifies the port that outgoing traffic from the AppPorts is directed to.
	// Experimental.
	ProxyEgressPort *float64 `json:"proxyEgressPort"`
	// Specifies the port that incoming traffic to the AppPorts is directed to.
	// Experimental.
	ProxyIngressPort *float64 `json:"proxyIngressPort"`
	// The egress traffic going to these specified IP addresses is ignored and not redirected to the ProxyEgressPort.
	//
	// It can be an empty list.
	// Experimental.
	EgressIgnoredIPs *[]*string `json:"egressIgnoredIPs"`
	// The egress traffic going to these specified ports is ignored and not redirected to the ProxyEgressPort.
	//
	// It can be an empty list.
	// Experimental.
	EgressIgnoredPorts *[]*float64 `json:"egressIgnoredPorts"`
	// The group ID (GID) of the proxy container as defined by the user parameter in a container definition.
	//
	// This is used to ensure the proxy ignores its own traffic. If IgnoredUID is specified, this field can be empty.
	// Experimental.
	IgnoredGID *float64 `json:"ignoredGID"`
	// The user ID (UID) of the proxy container as defined by the user parameter in a container definition.
	//
	// This is used to ensure the proxy ignores its own traffic. If IgnoredGID is specified, this field can be empty.
	// Experimental.
	IgnoredUID *float64 `json:"ignoredUID"`
}

// An Auto Scaling Group Capacity Provider.
//
// This allows an ECS cluster to target
// a specific EC2 Auto Scaling Group for the placement of tasks. Optionally (and
// recommended), ECS can manage the number of instances in the ASG to fit the
// tasks, and can ensure that instances are not prematurely terminated while
// there are still tasks running on them.
// Experimental.
type AsgCapacityProvider interface {
	awscdk.Construct
	AutoScalingGroup() awsautoscaling.AutoScalingGroup
	CapacityProviderName() *string
	EnableManagedTerminationProtection() *bool
	Node() awscdk.ConstructNode
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for AsgCapacityProvider
type jsiiProxy_AsgCapacityProvider struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_AsgCapacityProvider) AutoScalingGroup() awsautoscaling.AutoScalingGroup {
	var returns awsautoscaling.AutoScalingGroup
	_jsii_.Get(
		j,
		"autoScalingGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AsgCapacityProvider) CapacityProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"capacityProviderName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AsgCapacityProvider) EnableManagedTerminationProtection() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"enableManagedTerminationProtection",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AsgCapacityProvider) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


// Experimental.
func NewAsgCapacityProvider(scope constructs.Construct, id *string, props *AsgCapacityProviderProps) AsgCapacityProvider {
	_init_.Initialize()

	j := jsiiProxy_AsgCapacityProvider{}

	_jsii_.Create(
		"monocdk.aws_ecs.AsgCapacityProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewAsgCapacityProvider_Override(a AsgCapacityProvider, scope constructs.Construct, id *string, props *AsgCapacityProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.AsgCapacityProvider",
		[]interface{}{scope, id, props},
		a,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func AsgCapacityProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AsgCapacityProvider",
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
func (a *jsiiProxy_AsgCapacityProvider) OnPrepare() {
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
func (a *jsiiProxy_AsgCapacityProvider) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_AsgCapacityProvider) OnValidate() *[]*string {
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
func (a *jsiiProxy_AsgCapacityProvider) Prepare() {
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
func (a *jsiiProxy_AsgCapacityProvider) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_AsgCapacityProvider) ToString() *string {
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
func (a *jsiiProxy_AsgCapacityProvider) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The options for creating an Auto Scaling Group Capacity Provider.
// Experimental.
type AsgCapacityProviderProps struct {
	// Specifies whether the containers can access the container instance role.
	// Experimental.
	CanContainersAccessInstanceRole *bool `json:"canContainersAccessInstanceRole"`
	// Specify the machine image type.
	// Experimental.
	MachineImageType MachineImageType `json:"machineImageType"`
	// Specify whether to enable Automated Draining for Spot Instances running Amazon ECS Services.
	//
	// For more information, see [Using Spot Instances](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/container-instance-spot.html).
	// Experimental.
	SpotInstanceDraining *bool `json:"spotInstanceDraining"`
	// The time period to wait before force terminating an instance that is draining.
	//
	// This creates a Lambda function that is used by a lifecycle hook for the
	// AutoScalingGroup that will delay instance termination until all ECS tasks
	// have drained from the instance. Set to 0 to disable task draining.
	//
	// Set to 0 to disable task draining.
	// Deprecated: The lifecycle draining hook is not configured if using the EC2 Capacity Provider. Enable managed termination protection instead.
	TaskDrainTime awscdk.Duration `json:"taskDrainTime"`
	// If {@link AddAutoScalingGroupCapacityOptions.taskDrainTime} is non-zero, then the ECS cluster creates an SNS Topic to as part of a system to drain instances of tasks when the instance is being shut down. If this property is provided, then this key will be used to encrypt the contents of that SNS Topic. See [SNS Data Encryption](https://docs.aws.amazon.com/sns/latest/dg/sns-data-encryption.html) for more information.
	// Experimental.
	TopicEncryptionKey awskms.IKey `json:"topicEncryptionKey"`
	// The autoscaling group to add as a Capacity Provider.
	// Experimental.
	AutoScalingGroup awsautoscaling.IAutoScalingGroup `json:"autoScalingGroup"`
	// The name for the capacity provider.
	// Experimental.
	CapacityProviderName *string `json:"capacityProviderName"`
	// Whether to enable managed scaling.
	// Experimental.
	EnableManagedScaling *bool `json:"enableManagedScaling"`
	// Whether to enable managed termination protection.
	// Experimental.
	EnableManagedTerminationProtection *bool `json:"enableManagedTerminationProtection"`
	// Maximum scaling step size.
	//
	// In most cases this should be left alone.
	// Experimental.
	MaximumScalingStepSize *float64 `json:"maximumScalingStepSize"`
	// Minimum scaling step size.
	//
	// In most cases this should be left alone.
	// Experimental.
	MinimumScalingStepSize *float64 `json:"minimumScalingStepSize"`
	// Target capacity percent.
	//
	// In most cases this should be left alone.
	// Experimental.
	TargetCapacityPercent *float64 `json:"targetCapacityPercent"`
}

// Environment file from a local directory.
// Experimental.
type AssetEnvironmentFile interface {
	EnvironmentFile
	Path() *string
	Bind(scope awscdk.Construct) *EnvironmentFileConfig
}

// The jsii proxy struct for AssetEnvironmentFile
type jsiiProxy_AssetEnvironmentFile struct {
	jsiiProxy_EnvironmentFile
}

func (j *jsiiProxy_AssetEnvironmentFile) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}


// Experimental.
func NewAssetEnvironmentFile(path *string, options *awss3assets.AssetOptions) AssetEnvironmentFile {
	_init_.Initialize()

	j := jsiiProxy_AssetEnvironmentFile{}

	_jsii_.Create(
		"monocdk.aws_ecs.AssetEnvironmentFile",
		[]interface{}{path, options},
		&j,
	)

	return &j
}

// Experimental.
func NewAssetEnvironmentFile_Override(a AssetEnvironmentFile, path *string, options *awss3assets.AssetOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.AssetEnvironmentFile",
		[]interface{}{path, options},
		a,
	)
}

// Loads the environment file from a local disk path.
// Experimental.
func AssetEnvironmentFile_FromAsset(path *string, options *awss3assets.AssetOptions) AssetEnvironmentFile {
	_init_.Initialize()

	var returns AssetEnvironmentFile

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AssetEnvironmentFile",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Loads the environment file from an S3 bucket.
//
// Returns: `S3EnvironmentFile` associated with the specified S3 object.
// Experimental.
func AssetEnvironmentFile_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3EnvironmentFile {
	_init_.Initialize()

	var returns S3EnvironmentFile

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AssetEnvironmentFile",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Called when the container is initialized to allow this object to bind to the stack.
// Experimental.
func (a *jsiiProxy_AssetEnvironmentFile) Bind(scope awscdk.Construct) *EnvironmentFileConfig {
	var returns *EnvironmentFileConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// An image that will be built from a local directory with a Dockerfile.
// Experimental.
type AssetImage interface {
	ContainerImage
	Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig
}

// The jsii proxy struct for AssetImage
type jsiiProxy_AssetImage struct {
	jsiiProxy_ContainerImage
}

// Constructs a new instance of the AssetImage class.
// Experimental.
func NewAssetImage(directory *string, props *AssetImageProps) AssetImage {
	_init_.Initialize()

	j := jsiiProxy_AssetImage{}

	_jsii_.Create(
		"monocdk.aws_ecs.AssetImage",
		[]interface{}{directory, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the AssetImage class.
// Experimental.
func NewAssetImage_Override(a AssetImage, directory *string, props *AssetImageProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.AssetImage",
		[]interface{}{directory, props},
		a,
	)
}

// Reference an image that's constructed directly from sources on disk.
//
// If you already have a `DockerImageAsset` instance, you can use the
// `ContainerImage.fromDockerImageAsset` method instead.
// Experimental.
func AssetImage_FromAsset(directory *string, props *AssetImageProps) AssetImage {
	_init_.Initialize()

	var returns AssetImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AssetImage",
		"fromAsset",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Use an existing `DockerImageAsset` for this container image.
// Experimental.
func AssetImage_FromDockerImageAsset(asset awsecrassets.DockerImageAsset) ContainerImage {
	_init_.Initialize()

	var returns ContainerImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AssetImage",
		"fromDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Reference an image in an ECR repository.
// Experimental.
func AssetImage_FromEcrRepository(repository awsecr.IRepository, tag *string) EcrImage {
	_init_.Initialize()

	var returns EcrImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AssetImage",
		"fromEcrRepository",
		[]interface{}{repository, tag},
		&returns,
	)

	return returns
}

// Reference an image on DockerHub or another online registry.
// Experimental.
func AssetImage_FromRegistry(name *string, props *RepositoryImageProps) RepositoryImage {
	_init_.Initialize()

	var returns RepositoryImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AssetImage",
		"fromRegistry",
		[]interface{}{name, props},
		&returns,
	)

	return returns
}

// Called when the image is used by a ContainerDefinition.
// Experimental.
func (a *jsiiProxy_AssetImage) Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig {
	var returns *ContainerImageConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope, containerDefinition},
		&returns,
	)

	return returns
}

// The properties for building an AssetImage.
// Experimental.
type AssetImageProps struct {
	// Glob patterns to exclude from the copy.
	// Experimental.
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Deprecated: use `followSymlinks` instead
	Follow assets.FollowMode `json:"follow"`
	// The ignore behavior to use for exclude patterns.
	// Experimental.
	IgnoreMode awscdk.IgnoreMode `json:"ignoreMode"`
	// Extra information to encode into the fingerprint (e.g. build instructions and other inputs).
	// Experimental.
	ExtraHash *string `json:"extraHash"`
	// A strategy for how to handle symlinks.
	// Experimental.
	FollowSymlinks awscdk.SymlinkFollowMode `json:"followSymlinks"`
	// Build args to pass to the `docker build` command.
	//
	// Since Docker build arguments are resolved before deployment, keys and
	// values cannot refer to unresolved tokens (such as `lambda.functionArn` or
	// `queue.queueUrl`).
	// Experimental.
	BuildArgs *map[string]*string `json:"buildArgs"`
	// Path to the Dockerfile (relative to the directory).
	// Experimental.
	File *string `json:"file"`
	// ECR repository name.
	//
	// Specify this property if you need to statically address the image, e.g.
	// from a Kubernetes Pod. Note, this is only the repository name, without the
	// registry and the tag parts.
	// Deprecated: to control the location of docker image assets, please override
	// `Stack.addDockerImageAsset`. this feature will be removed in future
	// releases.
	RepositoryName *string `json:"repositoryName"`
	// Docker target to build to.
	// Experimental.
	Target *string `json:"target"`
}

// The options for using a cloudmap service.
// Experimental.
type AssociateCloudMapServiceOptions struct {
	// The cloudmap service to register with.
	// Experimental.
	Service awsservicediscovery.IService `json:"service"`
	// The container to point to for a SRV record.
	// Experimental.
	Container ContainerDefinition `json:"container"`
	// The port to point to for a SRV record.
	// Experimental.
	ContainerPort *float64 `json:"containerPort"`
}

// The authorization configuration details for the Amazon EFS file system.
// Experimental.
type AuthorizationConfig struct {
	// The access point ID to use.
	//
	// If an access point is specified, the root directory value will be
	// relative to the directory set for the access point.
	// If specified, transit encryption must be enabled in the EFSVolumeConfiguration.
	// Experimental.
	AccessPointId *string `json:"accessPointId"`
	// Whether or not to use the Amazon ECS task IAM role defined in a task definition when mounting the Amazon EFS file system.
	//
	// If enabled, transit encryption must be enabled in the EFSVolumeConfiguration.
	//
	// Valid values: ENABLED | DISABLED
	// Experimental.
	Iam *string `json:"iam"`
}

// A log driver that sends log information to CloudWatch Logs.
// Experimental.
type AwsLogDriver interface {
	LogDriver
	LogGroup() awslogs.ILogGroup
	SetLogGroup(val awslogs.ILogGroup)
	Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for AwsLogDriver
type jsiiProxy_AwsLogDriver struct {
	jsiiProxy_LogDriver
}

func (j *jsiiProxy_AwsLogDriver) LogGroup() awslogs.ILogGroup {
	var returns awslogs.ILogGroup
	_jsii_.Get(
		j,
		"logGroup",
		&returns,
	)
	return returns
}


// Constructs a new instance of the AwsLogDriver class.
// Experimental.
func NewAwsLogDriver(props *AwsLogDriverProps) AwsLogDriver {
	_init_.Initialize()

	j := jsiiProxy_AwsLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.AwsLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the AwsLogDriver class.
// Experimental.
func NewAwsLogDriver_Override(a AwsLogDriver, props *AwsLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.AwsLogDriver",
		[]interface{}{props},
		a,
	)
}

func (j *jsiiProxy_AwsLogDriver) SetLogGroup(val awslogs.ILogGroup) {
	_jsii_.Set(
		j,
		"logGroup",
		val,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func AwsLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.AwsLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (a *jsiiProxy_AwsLogDriver) Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope, containerDefinition},
		&returns,
	)

	return returns
}

// awslogs provides two modes for delivering messages from the container to the log driver.
// Experimental.
type AwsLogDriverMode string

const (
	AwsLogDriverMode_BLOCKING AwsLogDriverMode = "BLOCKING"
	AwsLogDriverMode_NON_BLOCKING AwsLogDriverMode = "NON_BLOCKING"
)

// Specifies the awslogs log driver configuration options.
// Experimental.
type AwsLogDriverProps struct {
	// Prefix for the log streams.
	//
	// The awslogs-stream-prefix option allows you to associate a log stream
	// with the specified prefix, the container name, and the ID of the Amazon
	// ECS task to which the container belongs. If you specify a prefix with
	// this option, then the log stream takes the following format:
	//
	//      prefix-name/container-name/ecs-task-id
	// Experimental.
	StreamPrefix *string `json:"streamPrefix"`
	// This option defines a multiline start pattern in Python strftime format.
	//
	// A log message consists of a line that matches the pattern and any
	// following lines that don’t match the pattern. Thus the matched line is
	// the delimiter between log messages.
	// Experimental.
	DatetimeFormat *string `json:"datetimeFormat"`
	// The log group to log to.
	// Experimental.
	LogGroup awslogs.ILogGroup `json:"logGroup"`
	// The number of days log events are kept in CloudWatch Logs when the log group is automatically created by this construct.
	// Experimental.
	LogRetention awslogs.RetentionDays `json:"logRetention"`
	// The delivery mode of log messages from the container to awslogs.
	// Experimental.
	Mode AwsLogDriverMode `json:"mode"`
	// This option defines a multiline start pattern using a regular expression.
	//
	// A log message consists of a line that matches the pattern and any
	// following lines that don’t match the pattern. Thus the matched line is
	// the delimiter between log messages.
	//
	// This option is ignored if datetimeFormat is also configured.
	// Experimental.
	MultilinePattern *string `json:"multilinePattern"`
}

// Experimental.
type BaseLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
}

// The base class for Ec2Service and FargateService services.
// Experimental.
type BaseService interface {
	awscdk.Resource
	IBaseService
	awselasticloadbalancing.ILoadBalancerTarget
	awselasticloadbalancingv2.IApplicationLoadBalancerTarget
	awselasticloadbalancingv2.INetworkLoadBalancerTarget
	CloudmapService() awsservicediscovery.Service
	SetCloudmapService(val awsservicediscovery.Service)
	CloudMapService() awsservicediscovery.IService
	Cluster() ICluster
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	LoadBalancers() *[]*CfnService_LoadBalancerProperty
	SetLoadBalancers(val *[]*CfnService_LoadBalancerProperty)
	NetworkConfiguration() *CfnService_NetworkConfigurationProperty
	SetNetworkConfiguration(val *CfnService_NetworkConfigurationProperty)
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ServiceArn() *string
	ServiceName() *string
	ServiceRegistries() *[]*CfnService_ServiceRegistryProperty
	SetServiceRegistries(val *[]*CfnService_ServiceRegistryProperty)
	Stack() awscdk.Stack
	TaskDefinition() TaskDefinition
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AssociateCloudMapService(options *AssociateCloudMapServiceOptions)
	AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer)
	AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	AutoScaleTaskCount(props *awsapplicationautoscaling.EnableScalingProps) ScalableTaskCount
	ConfigureAwsVpcNetworking(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroup awsec2.ISecurityGroup)
	ConfigureAwsVpcNetworkingWithSecurityGroups(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroups *[]awsec2.ISecurityGroup)
	EnableCloudMap(options *CloudMapOptions) awsservicediscovery.Service
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LoadBalancerTarget(options *LoadBalancerTargetOptions) IEcsLoadBalancerTarget
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RegisterLoadBalancerTargets(targets ...*EcsTarget)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for BaseService
type jsiiProxy_BaseService struct {
	internal.Type__awscdkResource
	jsiiProxy_IBaseService
	internal.Type__awselasticloadbalancingILoadBalancerTarget
	internal.Type__awselasticloadbalancingv2IApplicationLoadBalancerTarget
	internal.Type__awselasticloadbalancingv2INetworkLoadBalancerTarget
}

func (j *jsiiProxy_BaseService) CloudmapService() awsservicediscovery.Service {
	var returns awsservicediscovery.Service
	_jsii_.Get(
		j,
		"cloudmapService",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) CloudMapService() awsservicediscovery.IService {
	var returns awsservicediscovery.IService
	_jsii_.Get(
		j,
		"cloudMapService",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) Cluster() ICluster {
	var returns ICluster
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) LoadBalancers() *[]*CfnService_LoadBalancerProperty {
	var returns *[]*CfnService_LoadBalancerProperty
	_jsii_.Get(
		j,
		"loadBalancers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) NetworkConfiguration() *CfnService_NetworkConfigurationProperty {
	var returns *CfnService_NetworkConfigurationProperty
	_jsii_.Get(
		j,
		"networkConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) ServiceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) ServiceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) ServiceRegistries() *[]*CfnService_ServiceRegistryProperty {
	var returns *[]*CfnService_ServiceRegistryProperty
	_jsii_.Get(
		j,
		"serviceRegistries",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseService) TaskDefinition() TaskDefinition {
	var returns TaskDefinition
	_jsii_.Get(
		j,
		"taskDefinition",
		&returns,
	)
	return returns
}


// Constructs a new instance of the BaseService class.
// Experimental.
func NewBaseService_Override(b BaseService, scope constructs.Construct, id *string, props *BaseServiceProps, additionalProps interface{}, taskDefinition TaskDefinition) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.BaseService",
		[]interface{}{scope, id, props, additionalProps, taskDefinition},
		b,
	)
}

func (j *jsiiProxy_BaseService) SetCloudmapService(val awsservicediscovery.Service) {
	_jsii_.Set(
		j,
		"cloudmapService",
		val,
	)
}

func (j *jsiiProxy_BaseService) SetLoadBalancers(val *[]*CfnService_LoadBalancerProperty) {
	_jsii_.Set(
		j,
		"loadBalancers",
		val,
	)
}

func (j *jsiiProxy_BaseService) SetNetworkConfiguration(val *CfnService_NetworkConfigurationProperty) {
	_jsii_.Set(
		j,
		"networkConfiguration",
		val,
	)
}

func (j *jsiiProxy_BaseService) SetServiceRegistries(val *[]*CfnService_ServiceRegistryProperty) {
	_jsii_.Set(
		j,
		"serviceRegistries",
		val,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BaseService_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.BaseService",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func BaseService_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.BaseService",
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
func (b *jsiiProxy_BaseService) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Associates this service with a CloudMap service.
// Experimental.
func (b *jsiiProxy_BaseService) AssociateCloudMapService(options *AssociateCloudMapServiceOptions) {
	_jsii_.InvokeVoid(
		b,
		"associateCloudMapService",
		[]interface{}{options},
	)
}

// This method is called to attach this service to an Application Load Balancer.
//
// Don't call this function directly. Instead, call `listener.addTargets()`
// to add this service to a load balancer.
// Experimental.
func (b *jsiiProxy_BaseService) AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		b,
		"attachToApplicationTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// Registers the service as a target of a Classic Load Balancer (CLB).
//
// Don't call this. Call `loadBalancer.addTarget()` instead.
// Experimental.
func (b *jsiiProxy_BaseService) AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer) {
	_jsii_.InvokeVoid(
		b,
		"attachToClassicLB",
		[]interface{}{loadBalancer},
	)
}

// This method is called to attach this service to a Network Load Balancer.
//
// Don't call this function directly. Instead, call `listener.addTargets()`
// to add this service to a load balancer.
// Experimental.
func (b *jsiiProxy_BaseService) AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		b,
		"attachToNetworkTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// An attribute representing the minimum and maximum task count for an AutoScalingGroup.
// Experimental.
func (b *jsiiProxy_BaseService) AutoScaleTaskCount(props *awsapplicationautoscaling.EnableScalingProps) ScalableTaskCount {
	var returns ScalableTaskCount

	_jsii_.Invoke(
		b,
		"autoScaleTaskCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method is called to create a networkConfiguration.
// Deprecated: use configureAwsVpcNetworkingWithSecurityGroups instead.
func (b *jsiiProxy_BaseService) ConfigureAwsVpcNetworking(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroup awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		b,
		"configureAwsVpcNetworking",
		[]interface{}{vpc, assignPublicIp, vpcSubnets, securityGroup},
	)
}

// This method is called to create a networkConfiguration.
// Experimental.
func (b *jsiiProxy_BaseService) ConfigureAwsVpcNetworkingWithSecurityGroups(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroups *[]awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		b,
		"configureAwsVpcNetworkingWithSecurityGroups",
		[]interface{}{vpc, assignPublicIp, vpcSubnets, securityGroups},
	)
}

// Enable CloudMap service discovery for the service.
//
// Returns: The created CloudMap service
// Experimental.
func (b *jsiiProxy_BaseService) EnableCloudMap(options *CloudMapOptions) awsservicediscovery.Service {
	var returns awsservicediscovery.Service

	_jsii_.Invoke(
		b,
		"enableCloudMap",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Experimental.
func (b *jsiiProxy_BaseService) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		b,
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
func (b *jsiiProxy_BaseService) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		b,
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
func (b *jsiiProxy_BaseService) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Return a load balancing target for a specific container and port.
//
// Use this function to create a load balancer target if you want to load balance to
// another container than the first essential container or the first mapped port on
// the container.
//
// Use the return value of this function where you would normally use a load balancer
// target, instead of the `Service` object itself.
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_BaseService) LoadBalancerTarget(options *LoadBalancerTargetOptions) IEcsLoadBalancerTarget {
	var returns IEcsLoadBalancerTarget

	_jsii_.Invoke(
		b,
		"loadBalancerTarget",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// This method returns the specified CloudWatch metric name for this service.
// Experimental.
func (b *jsiiProxy_BaseService) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		b,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this service's CPU utilization.
// Experimental.
func (b *jsiiProxy_BaseService) MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		b,
		"metricCpuUtilization",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this service's memory utilization.
// Experimental.
func (b *jsiiProxy_BaseService) MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		b,
		"metricMemoryUtilization",
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
func (b *jsiiProxy_BaseService) OnPrepare() {
	_jsii_.InvokeVoid(
		b,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_BaseService) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
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
func (b *jsiiProxy_BaseService) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
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
func (b *jsiiProxy_BaseService) Prepare() {
	_jsii_.InvokeVoid(
		b,
		"prepare",
		nil, // no parameters
	)
}

// Use this function to create all load balancer targets to be registered in this service, add them to target groups, and attach target groups to listeners accordingly.
//
// Alternatively, you can use `listener.addTargets()` to create targets and add them to target groups.
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_BaseService) RegisterLoadBalancerTargets(targets ...*EcsTarget) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		b,
		"registerLoadBalancerTargets",
		args,
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_BaseService) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BaseService) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		b,
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
func (b *jsiiProxy_BaseService) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The properties for the base Ec2Service or FargateService service.
// Experimental.
type BaseServiceOptions struct {
	// The name of the cluster that hosts the service.
	// Experimental.
	Cluster ICluster `json:"cluster"`
	// A list of Capacity Provider strategies used to place a service.
	// Experimental.
	CapacityProviderStrategies *[]*CapacityProviderStrategy `json:"capacityProviderStrategies"`
	// Whether to enable the deployment circuit breaker.
	//
	// If this property is defined, circuit breaker will be implicitly
	// enabled.
	// Experimental.
	CircuitBreaker *DeploymentCircuitBreaker `json:"circuitBreaker"`
	// The options for configuring an Amazon ECS service to use service discovery.
	// Experimental.
	CloudMapOptions *CloudMapOptions `json:"cloudMapOptions"`
	// Specifies which deployment controller to use for the service.
	//
	// For more information, see
	// [Amazon ECS Deployment Types](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-types.html)
	// Experimental.
	DeploymentController *DeploymentController `json:"deploymentController"`
	// The desired number of instantiations of the task definition to keep running on the service.
	// Experimental.
	DesiredCount *float64 `json:"desiredCount"`
	// Specifies whether to enable Amazon ECS managed tags for the tasks within the service.
	//
	// For more information, see
	// [Tagging Your Amazon ECS Resources](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-using-tags.html)
	// Experimental.
	EnableECSManagedTags *bool `json:"enableECSManagedTags"`
	// Whether to enable the ability to execute into a container.
	// Experimental.
	EnableExecuteCommand *bool `json:"enableExecuteCommand"`
	// The period of time, in seconds, that the Amazon ECS service scheduler ignores unhealthy Elastic Load Balancing target health checks after a task has first started.
	// Experimental.
	HealthCheckGracePeriod awscdk.Duration `json:"healthCheckGracePeriod"`
	// The maximum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that can run in a service during a deployment.
	// Experimental.
	MaxHealthyPercent *float64 `json:"maxHealthyPercent"`
	// The minimum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that must continue to run and remain healthy during a deployment.
	// Experimental.
	MinHealthyPercent *float64 `json:"minHealthyPercent"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Valid values are: PropagatedTagSource.SERVICE, PropagatedTagSource.TASK_DEFINITION or PropagatedTagSource.NONE
	// Experimental.
	PropagateTags PropagatedTagSource `json:"propagateTags"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Tags can only be propagated to the tasks within the service during service creation.
	// Deprecated: Use `propagateTags` instead.
	PropagateTaskTagsFrom PropagatedTagSource `json:"propagateTaskTagsFrom"`
	// The name of the service.
	// Experimental.
	ServiceName *string `json:"serviceName"`
}

// Complete base service properties that are required to be supplied by the implementation of the BaseService class.
// Experimental.
type BaseServiceProps struct {
	// The name of the cluster that hosts the service.
	// Experimental.
	Cluster ICluster `json:"cluster"`
	// A list of Capacity Provider strategies used to place a service.
	// Experimental.
	CapacityProviderStrategies *[]*CapacityProviderStrategy `json:"capacityProviderStrategies"`
	// Whether to enable the deployment circuit breaker.
	//
	// If this property is defined, circuit breaker will be implicitly
	// enabled.
	// Experimental.
	CircuitBreaker *DeploymentCircuitBreaker `json:"circuitBreaker"`
	// The options for configuring an Amazon ECS service to use service discovery.
	// Experimental.
	CloudMapOptions *CloudMapOptions `json:"cloudMapOptions"`
	// Specifies which deployment controller to use for the service.
	//
	// For more information, see
	// [Amazon ECS Deployment Types](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-types.html)
	// Experimental.
	DeploymentController *DeploymentController `json:"deploymentController"`
	// The desired number of instantiations of the task definition to keep running on the service.
	// Experimental.
	DesiredCount *float64 `json:"desiredCount"`
	// Specifies whether to enable Amazon ECS managed tags for the tasks within the service.
	//
	// For more information, see
	// [Tagging Your Amazon ECS Resources](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-using-tags.html)
	// Experimental.
	EnableECSManagedTags *bool `json:"enableECSManagedTags"`
	// Whether to enable the ability to execute into a container.
	// Experimental.
	EnableExecuteCommand *bool `json:"enableExecuteCommand"`
	// The period of time, in seconds, that the Amazon ECS service scheduler ignores unhealthy Elastic Load Balancing target health checks after a task has first started.
	// Experimental.
	HealthCheckGracePeriod awscdk.Duration `json:"healthCheckGracePeriod"`
	// The maximum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that can run in a service during a deployment.
	// Experimental.
	MaxHealthyPercent *float64 `json:"maxHealthyPercent"`
	// The minimum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that must continue to run and remain healthy during a deployment.
	// Experimental.
	MinHealthyPercent *float64 `json:"minHealthyPercent"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Valid values are: PropagatedTagSource.SERVICE, PropagatedTagSource.TASK_DEFINITION or PropagatedTagSource.NONE
	// Experimental.
	PropagateTags PropagatedTagSource `json:"propagateTags"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Tags can only be propagated to the tasks within the service during service creation.
	// Deprecated: Use `propagateTags` instead.
	PropagateTaskTagsFrom PropagatedTagSource `json:"propagateTaskTagsFrom"`
	// The name of the service.
	// Experimental.
	ServiceName *string `json:"serviceName"`
	// The launch type on which to run your service.
	//
	// LaunchType will be omitted if capacity provider strategies are specified on the service.
	// See: - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ecs-service.html#cfn-ecs-service-capacityproviderstrategy
	//
	// Valid values are: LaunchType.ECS or LaunchType.FARGATE or LaunchType.EXTERNAL
	//
	// Experimental.
	LaunchType LaunchType `json:"launchType"`
}

// Instance resource used for bin packing.
// Experimental.
type BinPackResource string

const (
	BinPackResource_CPU BinPackResource = "CPU"
	BinPackResource_MEMORY BinPackResource = "MEMORY"
)

// Construct an Bottlerocket image from the latest AMI published in SSM.
// Experimental.
type BottleRocketImage interface {
	awsec2.IMachineImage
	GetImage(scope awscdk.Construct) *awsec2.MachineImageConfig
}

// The jsii proxy struct for BottleRocketImage
type jsiiProxy_BottleRocketImage struct {
	internal.Type__awsec2IMachineImage
}

// Constructs a new instance of the BottleRocketImage class.
// Experimental.
func NewBottleRocketImage(props *BottleRocketImageProps) BottleRocketImage {
	_init_.Initialize()

	j := jsiiProxy_BottleRocketImage{}

	_jsii_.Create(
		"monocdk.aws_ecs.BottleRocketImage",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the BottleRocketImage class.
// Experimental.
func NewBottleRocketImage_Override(b BottleRocketImage, props *BottleRocketImageProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.BottleRocketImage",
		[]interface{}{props},
		b,
	)
}

// Return the correct image.
// Experimental.
func (b *jsiiProxy_BottleRocketImage) GetImage(scope awscdk.Construct) *awsec2.MachineImageConfig {
	var returns *awsec2.MachineImageConfig

	_jsii_.Invoke(
		b,
		"getImage",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Properties for BottleRocketImage.
// Experimental.
type BottleRocketImageProps struct {
	// The Amazon ECS variant to use.
	//
	// Only `aws-ecs-1` is currently available
	// Experimental.
	Variant BottlerocketEcsVariant `json:"variant"`
}

// Amazon ECS variant.
// Experimental.
type BottlerocketEcsVariant string

const (
	BottlerocketEcsVariant_AWS_ECS_1 BottlerocketEcsVariant = "AWS_ECS_1"
)

// The built-in container instance attributes.
// Experimental.
type BuiltInAttributes interface {
}

// The jsii proxy struct for BuiltInAttributes
type jsiiProxy_BuiltInAttributes struct {
	_ byte // padding
}

// Experimental.
func NewBuiltInAttributes() BuiltInAttributes {
	_init_.Initialize()

	j := jsiiProxy_BuiltInAttributes{}

	_jsii_.Create(
		"monocdk.aws_ecs.BuiltInAttributes",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewBuiltInAttributes_Override(b BuiltInAttributes) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.BuiltInAttributes",
		nil, // no parameters
		b,
	)
}

func BuiltInAttributes_AMI_ID() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.BuiltInAttributes",
		"AMI_ID",
		&returns,
	)
	return returns
}

func BuiltInAttributes_AVAILABILITY_ZONE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.BuiltInAttributes",
		"AVAILABILITY_ZONE",
		&returns,
	)
	return returns
}

func BuiltInAttributes_INSTANCE_ID() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.BuiltInAttributes",
		"INSTANCE_ID",
		&returns,
	)
	return returns
}

func BuiltInAttributes_INSTANCE_TYPE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.BuiltInAttributes",
		"INSTANCE_TYPE",
		&returns,
	)
	return returns
}

func BuiltInAttributes_OS_TYPE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.BuiltInAttributes",
		"OS_TYPE",
		&returns,
	)
	return returns
}

// A Linux capability.
// Experimental.
type Capability string

const (
	Capability_ALL Capability = "ALL"
	Capability_AUDIT_CONTROL Capability = "AUDIT_CONTROL"
	Capability_AUDIT_WRITE Capability = "AUDIT_WRITE"
	Capability_BLOCK_SUSPEND Capability = "BLOCK_SUSPEND"
	Capability_CHOWN Capability = "CHOWN"
	Capability_DAC_OVERRIDE Capability = "DAC_OVERRIDE"
	Capability_DAC_READ_SEARCH Capability = "DAC_READ_SEARCH"
	Capability_FOWNER Capability = "FOWNER"
	Capability_FSETID Capability = "FSETID"
	Capability_IPC_LOCK Capability = "IPC_LOCK"
	Capability_IPC_OWNER Capability = "IPC_OWNER"
	Capability_KILL Capability = "KILL"
	Capability_LEASE Capability = "LEASE"
	Capability_LINUX_IMMUTABLE Capability = "LINUX_IMMUTABLE"
	Capability_MAC_ADMIN Capability = "MAC_ADMIN"
	Capability_MAC_OVERRIDE Capability = "MAC_OVERRIDE"
	Capability_MKNOD Capability = "MKNOD"
	Capability_NET_ADMIN Capability = "NET_ADMIN"
	Capability_NET_BIND_SERVICE Capability = "NET_BIND_SERVICE"
	Capability_NET_BROADCAST Capability = "NET_BROADCAST"
	Capability_NET_RAW Capability = "NET_RAW"
	Capability_SETFCAP Capability = "SETFCAP"
	Capability_SETGID Capability = "SETGID"
	Capability_SETPCAP Capability = "SETPCAP"
	Capability_SETUID Capability = "SETUID"
	Capability_SYS_ADMIN Capability = "SYS_ADMIN"
	Capability_SYS_BOOT Capability = "SYS_BOOT"
	Capability_SYS_CHROOT Capability = "SYS_CHROOT"
	Capability_SYS_MODULE Capability = "SYS_MODULE"
	Capability_SYS_NICE Capability = "SYS_NICE"
	Capability_SYS_PACCT Capability = "SYS_PACCT"
	Capability_SYS_PTRACE Capability = "SYS_PTRACE"
	Capability_SYS_RAWIO Capability = "SYS_RAWIO"
	Capability_SYS_RESOURCE Capability = "SYS_RESOURCE"
	Capability_SYS_TIME Capability = "SYS_TIME"
	Capability_SYS_TTY_CONFIG Capability = "SYS_TTY_CONFIG"
	Capability_SYSLOG Capability = "SYSLOG"
	Capability_WAKE_ALARM Capability = "WAKE_ALARM"
)

// A Capacity Provider strategy to use for the service.
//
// NOTE: defaultCapacityProviderStrategy on cluster not currently supported.
// Experimental.
type CapacityProviderStrategy struct {
	// The name of the capacity provider.
	// Experimental.
	CapacityProvider *string `json:"capacityProvider"`
	// The base value designates how many tasks, at a minimum, to run on the specified capacity provider.
	//
	// Only one
	// capacity provider in a capacity provider strategy can have a base defined. If no value is specified, the default
	// value of 0 is used.
	// Experimental.
	Base *float64 `json:"base"`
	// The weight value designates the relative percentage of the total number of tasks launched that should use the specified capacity provider.
	//
	// The weight value is taken into consideration after the base value, if defined, is satisfied.
	// Experimental.
	Weight *float64 `json:"weight"`
}

// A CloudFormation `AWS::ECS::CapacityProvider`.
type CfnCapacityProvider interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AutoScalingGroupProvider() interface{}
	SetAutoScalingGroupProvider(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
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

// The jsii proxy struct for CfnCapacityProvider
type jsiiProxy_CfnCapacityProvider struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnCapacityProvider) AutoScalingGroupProvider() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"autoScalingGroupProvider",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCapacityProvider) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ECS::CapacityProvider`.
func NewCfnCapacityProvider(scope awscdk.Construct, id *string, props *CfnCapacityProviderProps) CfnCapacityProvider {
	_init_.Initialize()

	j := jsiiProxy_CfnCapacityProvider{}

	_jsii_.Create(
		"monocdk.aws_ecs.CfnCapacityProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ECS::CapacityProvider`.
func NewCfnCapacityProvider_Override(c CfnCapacityProvider, scope awscdk.Construct, id *string, props *CfnCapacityProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.CfnCapacityProvider",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCapacityProvider) SetAutoScalingGroupProvider(val interface{}) {
	_jsii_.Set(
		j,
		"autoScalingGroupProvider",
		val,
	)
}

func (j *jsiiProxy_CfnCapacityProvider) SetName(val *string) {
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
func CfnCapacityProvider_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnCapacityProvider",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnCapacityProvider_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnCapacityProvider",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCapacityProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnCapacityProvider",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnCapacityProvider_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.CfnCapacityProvider",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnCapacityProvider) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnCapacityProvider) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnCapacityProvider) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnCapacityProvider) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnCapacityProvider) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnCapacityProvider) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnCapacityProvider) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnCapacityProvider) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnCapacityProvider) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnCapacityProvider) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnCapacityProvider) OnPrepare() {
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
func (c *jsiiProxy_CfnCapacityProvider) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCapacityProvider) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCapacityProvider) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCapacityProvider) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnCapacityProvider) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnCapacityProvider) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnCapacityProvider) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCapacityProvider) ToString() *string {
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
func (c *jsiiProxy_CfnCapacityProvider) Validate() *[]*string {
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
func (c *jsiiProxy_CfnCapacityProvider) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnCapacityProvider_AutoScalingGroupProviderProperty struct {
	// `CfnCapacityProvider.AutoScalingGroupProviderProperty.AutoScalingGroupArn`.
	AutoScalingGroupArn *string `json:"autoScalingGroupArn"`
	// `CfnCapacityProvider.AutoScalingGroupProviderProperty.ManagedScaling`.
	ManagedScaling interface{} `json:"managedScaling"`
	// `CfnCapacityProvider.AutoScalingGroupProviderProperty.ManagedTerminationProtection`.
	ManagedTerminationProtection *string `json:"managedTerminationProtection"`
}

type CfnCapacityProvider_ManagedScalingProperty struct {
	// `CfnCapacityProvider.ManagedScalingProperty.InstanceWarmupPeriod`.
	InstanceWarmupPeriod *float64 `json:"instanceWarmupPeriod"`
	// `CfnCapacityProvider.ManagedScalingProperty.MaximumScalingStepSize`.
	MaximumScalingStepSize *float64 `json:"maximumScalingStepSize"`
	// `CfnCapacityProvider.ManagedScalingProperty.MinimumScalingStepSize`.
	MinimumScalingStepSize *float64 `json:"minimumScalingStepSize"`
	// `CfnCapacityProvider.ManagedScalingProperty.Status`.
	Status *string `json:"status"`
	// `CfnCapacityProvider.ManagedScalingProperty.TargetCapacity`.
	TargetCapacity *float64 `json:"targetCapacity"`
}

// Properties for defining a `AWS::ECS::CapacityProvider`.
type CfnCapacityProviderProps struct {
	// `AWS::ECS::CapacityProvider.AutoScalingGroupProvider`.
	AutoScalingGroupProvider interface{} `json:"autoScalingGroupProvider"`
	// `AWS::ECS::CapacityProvider.Name`.
	Name *string `json:"name"`
	// `AWS::ECS::CapacityProvider.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::ECS::Cluster`.
type CfnCluster interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CapacityProviders() *[]*string
	SetCapacityProviders(val *[]*string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClusterName() *string
	SetClusterName(val *string)
	ClusterSettings() interface{}
	SetClusterSettings(val interface{})
	Configuration() interface{}
	SetConfiguration(val interface{})
	CreationStack() *[]*string
	DefaultCapacityProviderStrategy() interface{}
	SetDefaultCapacityProviderStrategy(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
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

// The jsii proxy struct for CfnCluster
type jsiiProxy_CfnCluster struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnCluster) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) CapacityProviders() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"capacityProviders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) ClusterName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clusterName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) ClusterSettings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"clusterSettings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) Configuration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"configuration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) DefaultCapacityProviderStrategy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"defaultCapacityProviderStrategy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCluster) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ECS::Cluster`.
func NewCfnCluster(scope awscdk.Construct, id *string, props *CfnClusterProps) CfnCluster {
	_init_.Initialize()

	j := jsiiProxy_CfnCluster{}

	_jsii_.Create(
		"monocdk.aws_ecs.CfnCluster",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ECS::Cluster`.
func NewCfnCluster_Override(c CfnCluster, scope awscdk.Construct, id *string, props *CfnClusterProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.CfnCluster",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCluster) SetCapacityProviders(val *[]*string) {
	_jsii_.Set(
		j,
		"capacityProviders",
		val,
	)
}

func (j *jsiiProxy_CfnCluster) SetClusterName(val *string) {
	_jsii_.Set(
		j,
		"clusterName",
		val,
	)
}

func (j *jsiiProxy_CfnCluster) SetClusterSettings(val interface{}) {
	_jsii_.Set(
		j,
		"clusterSettings",
		val,
	)
}

func (j *jsiiProxy_CfnCluster) SetConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"configuration",
		val,
	)
}

func (j *jsiiProxy_CfnCluster) SetDefaultCapacityProviderStrategy(val interface{}) {
	_jsii_.Set(
		j,
		"defaultCapacityProviderStrategy",
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
func CfnCluster_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnCluster",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnCluster_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnCluster",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCluster_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnCluster",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnCluster_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.CfnCluster",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnCluster) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnCluster) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnCluster) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnCluster) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnCluster) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnCluster) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnCluster) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnCluster) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnCluster) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnCluster) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnCluster) OnPrepare() {
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
func (c *jsiiProxy_CfnCluster) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCluster) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCluster) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCluster) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnCluster) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnCluster) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnCluster) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCluster) ToString() *string {
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
func (c *jsiiProxy_CfnCluster) Validate() *[]*string {
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
func (c *jsiiProxy_CfnCluster) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnCluster_CapacityProviderStrategyItemProperty struct {
	// `CfnCluster.CapacityProviderStrategyItemProperty.Base`.
	Base *float64 `json:"base"`
	// `CfnCluster.CapacityProviderStrategyItemProperty.CapacityProvider`.
	CapacityProvider *string `json:"capacityProvider"`
	// `CfnCluster.CapacityProviderStrategyItemProperty.Weight`.
	Weight *float64 `json:"weight"`
}

type CfnCluster_ClusterConfigurationProperty struct {
	// `CfnCluster.ClusterConfigurationProperty.ExecuteCommandConfiguration`.
	ExecuteCommandConfiguration interface{} `json:"executeCommandConfiguration"`
}

type CfnCluster_ClusterSettingsProperty struct {
	// `CfnCluster.ClusterSettingsProperty.Name`.
	Name *string `json:"name"`
	// `CfnCluster.ClusterSettingsProperty.Value`.
	Value *string `json:"value"`
}

type CfnCluster_ExecuteCommandConfigurationProperty struct {
	// `CfnCluster.ExecuteCommandConfigurationProperty.KmsKeyId`.
	KmsKeyId *string `json:"kmsKeyId"`
	// `CfnCluster.ExecuteCommandConfigurationProperty.LogConfiguration`.
	LogConfiguration interface{} `json:"logConfiguration"`
	// `CfnCluster.ExecuteCommandConfigurationProperty.Logging`.
	Logging *string `json:"logging"`
}

type CfnCluster_ExecuteCommandLogConfigurationProperty struct {
	// `CfnCluster.ExecuteCommandLogConfigurationProperty.CloudWatchEncryptionEnabled`.
	CloudWatchEncryptionEnabled interface{} `json:"cloudWatchEncryptionEnabled"`
	// `CfnCluster.ExecuteCommandLogConfigurationProperty.CloudWatchLogGroupName`.
	CloudWatchLogGroupName *string `json:"cloudWatchLogGroupName"`
	// `CfnCluster.ExecuteCommandLogConfigurationProperty.S3BucketName`.
	S3BucketName *string `json:"s3BucketName"`
	// `CfnCluster.ExecuteCommandLogConfigurationProperty.S3EncryptionEnabled`.
	S3EncryptionEnabled interface{} `json:"s3EncryptionEnabled"`
	// `CfnCluster.ExecuteCommandLogConfigurationProperty.S3KeyPrefix`.
	S3KeyPrefix *string `json:"s3KeyPrefix"`
}

// A CloudFormation `AWS::ECS::ClusterCapacityProviderAssociations`.
type CfnClusterCapacityProviderAssociations interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CapacityProviders() *[]*string
	SetCapacityProviders(val *[]*string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Cluster() *string
	SetCluster(val *string)
	CreationStack() *[]*string
	DefaultCapacityProviderStrategy() interface{}
	SetDefaultCapacityProviderStrategy(val interface{})
	LogicalId() *string
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

// The jsii proxy struct for CfnClusterCapacityProviderAssociations
type jsiiProxy_CfnClusterCapacityProviderAssociations struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) CapacityProviders() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"capacityProviders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) Cluster() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) DefaultCapacityProviderStrategy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"defaultCapacityProviderStrategy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ECS::ClusterCapacityProviderAssociations`.
func NewCfnClusterCapacityProviderAssociations(scope awscdk.Construct, id *string, props *CfnClusterCapacityProviderAssociationsProps) CfnClusterCapacityProviderAssociations {
	_init_.Initialize()

	j := jsiiProxy_CfnClusterCapacityProviderAssociations{}

	_jsii_.Create(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ECS::ClusterCapacityProviderAssociations`.
func NewCfnClusterCapacityProviderAssociations_Override(c CfnClusterCapacityProviderAssociations, scope awscdk.Construct, id *string, props *CfnClusterCapacityProviderAssociationsProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) SetCapacityProviders(val *[]*string) {
	_jsii_.Set(
		j,
		"capacityProviders",
		val,
	)
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) SetCluster(val *string) {
	_jsii_.Set(
		j,
		"cluster",
		val,
	)
}

func (j *jsiiProxy_CfnClusterCapacityProviderAssociations) SetDefaultCapacityProviderStrategy(val interface{}) {
	_jsii_.Set(
		j,
		"defaultCapacityProviderStrategy",
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
func CfnClusterCapacityProviderAssociations_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnClusterCapacityProviderAssociations_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnClusterCapacityProviderAssociations_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnClusterCapacityProviderAssociations_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.CfnClusterCapacityProviderAssociations",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) OnPrepare() {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) ToString() *string {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) Validate() *[]*string {
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
func (c *jsiiProxy_CfnClusterCapacityProviderAssociations) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnClusterCapacityProviderAssociations_CapacityProviderStrategyProperty struct {
	// `CfnClusterCapacityProviderAssociations.CapacityProviderStrategyProperty.CapacityProvider`.
	CapacityProvider *string `json:"capacityProvider"`
	// `CfnClusterCapacityProviderAssociations.CapacityProviderStrategyProperty.Base`.
	Base *float64 `json:"base"`
	// `CfnClusterCapacityProviderAssociations.CapacityProviderStrategyProperty.Weight`.
	Weight *float64 `json:"weight"`
}

// Properties for defining a `AWS::ECS::ClusterCapacityProviderAssociations`.
type CfnClusterCapacityProviderAssociationsProps struct {
	// `AWS::ECS::ClusterCapacityProviderAssociations.CapacityProviders`.
	CapacityProviders *[]*string `json:"capacityProviders"`
	// `AWS::ECS::ClusterCapacityProviderAssociations.Cluster`.
	Cluster *string `json:"cluster"`
	// `AWS::ECS::ClusterCapacityProviderAssociations.DefaultCapacityProviderStrategy`.
	DefaultCapacityProviderStrategy interface{} `json:"defaultCapacityProviderStrategy"`
}

// Properties for defining a `AWS::ECS::Cluster`.
type CfnClusterProps struct {
	// `AWS::ECS::Cluster.CapacityProviders`.
	CapacityProviders *[]*string `json:"capacityProviders"`
	// `AWS::ECS::Cluster.ClusterName`.
	ClusterName *string `json:"clusterName"`
	// `AWS::ECS::Cluster.ClusterSettings`.
	ClusterSettings interface{} `json:"clusterSettings"`
	// `AWS::ECS::Cluster.Configuration`.
	Configuration interface{} `json:"configuration"`
	// `AWS::ECS::Cluster.DefaultCapacityProviderStrategy`.
	DefaultCapacityProviderStrategy interface{} `json:"defaultCapacityProviderStrategy"`
	// `AWS::ECS::Cluster.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::ECS::PrimaryTaskSet`.
type CfnPrimaryTaskSet interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Cluster() *string
	SetCluster(val *string)
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Service() *string
	SetService(val *string)
	Stack() awscdk.Stack
	TaskSetId() *string
	SetTaskSetId(val *string)
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

// The jsii proxy struct for CfnPrimaryTaskSet
type jsiiProxy_CfnPrimaryTaskSet struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnPrimaryTaskSet) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) Cluster() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) Service() *string {
	var returns *string
	_jsii_.Get(
		j,
		"service",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) TaskSetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskSetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPrimaryTaskSet) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ECS::PrimaryTaskSet`.
func NewCfnPrimaryTaskSet(scope awscdk.Construct, id *string, props *CfnPrimaryTaskSetProps) CfnPrimaryTaskSet {
	_init_.Initialize()

	j := jsiiProxy_CfnPrimaryTaskSet{}

	_jsii_.Create(
		"monocdk.aws_ecs.CfnPrimaryTaskSet",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ECS::PrimaryTaskSet`.
func NewCfnPrimaryTaskSet_Override(c CfnPrimaryTaskSet, scope awscdk.Construct, id *string, props *CfnPrimaryTaskSetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.CfnPrimaryTaskSet",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnPrimaryTaskSet) SetCluster(val *string) {
	_jsii_.Set(
		j,
		"cluster",
		val,
	)
}

func (j *jsiiProxy_CfnPrimaryTaskSet) SetService(val *string) {
	_jsii_.Set(
		j,
		"service",
		val,
	)
}

func (j *jsiiProxy_CfnPrimaryTaskSet) SetTaskSetId(val *string) {
	_jsii_.Set(
		j,
		"taskSetId",
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
func CfnPrimaryTaskSet_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnPrimaryTaskSet",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnPrimaryTaskSet_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnPrimaryTaskSet",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnPrimaryTaskSet_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnPrimaryTaskSet",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnPrimaryTaskSet_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.CfnPrimaryTaskSet",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnPrimaryTaskSet) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnPrimaryTaskSet) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnPrimaryTaskSet) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) OnPrepare() {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnPrimaryTaskSet) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) ToString() *string {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) Validate() *[]*string {
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
func (c *jsiiProxy_CfnPrimaryTaskSet) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ECS::PrimaryTaskSet`.
type CfnPrimaryTaskSetProps struct {
	// `AWS::ECS::PrimaryTaskSet.Cluster`.
	Cluster *string `json:"cluster"`
	// `AWS::ECS::PrimaryTaskSet.Service`.
	Service *string `json:"service"`
	// `AWS::ECS::PrimaryTaskSet.TaskSetId`.
	TaskSetId *string `json:"taskSetId"`
}

// A CloudFormation `AWS::ECS::Service`.
type CfnService interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrName() *string
	AttrServiceArn() *string
	CapacityProviderStrategy() interface{}
	SetCapacityProviderStrategy(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Cluster() *string
	SetCluster(val *string)
	CreationStack() *[]*string
	DeploymentConfiguration() interface{}
	SetDeploymentConfiguration(val interface{})
	DeploymentController() interface{}
	SetDeploymentController(val interface{})
	DesiredCount() *float64
	SetDesiredCount(val *float64)
	EnableEcsManagedTags() interface{}
	SetEnableEcsManagedTags(val interface{})
	EnableExecuteCommand() interface{}
	SetEnableExecuteCommand(val interface{})
	HealthCheckGracePeriodSeconds() *float64
	SetHealthCheckGracePeriodSeconds(val *float64)
	LaunchType() *string
	SetLaunchType(val *string)
	LoadBalancers() interface{}
	SetLoadBalancers(val interface{})
	LogicalId() *string
	NetworkConfiguration() interface{}
	SetNetworkConfiguration(val interface{})
	Node() awscdk.ConstructNode
	PlacementConstraints() interface{}
	SetPlacementConstraints(val interface{})
	PlacementStrategies() interface{}
	SetPlacementStrategies(val interface{})
	PlatformVersion() *string
	SetPlatformVersion(val *string)
	PropagateTags() *string
	SetPropagateTags(val *string)
	Ref() *string
	Role() *string
	SetRole(val *string)
	SchedulingStrategy() *string
	SetSchedulingStrategy(val *string)
	ServiceName() *string
	SetServiceName(val *string)
	ServiceRegistries() interface{}
	SetServiceRegistries(val interface{})
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	TaskDefinition() *string
	SetTaskDefinition(val *string)
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

// The jsii proxy struct for CfnService
type jsiiProxy_CfnService struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnService) AttrName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) AttrServiceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrServiceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) CapacityProviderStrategy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"capacityProviderStrategy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) Cluster() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) DeploymentConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deploymentConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) DeploymentController() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deploymentController",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) DesiredCount() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"desiredCount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) EnableEcsManagedTags() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"enableEcsManagedTags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) EnableExecuteCommand() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"enableExecuteCommand",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) HealthCheckGracePeriodSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"healthCheckGracePeriodSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) LaunchType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"launchType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) LoadBalancers() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"loadBalancers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) NetworkConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"networkConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) PlacementConstraints() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"placementConstraints",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) PlacementStrategies() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"placementStrategies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) PlatformVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"platformVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) PropagateTags() *string {
	var returns *string
	_jsii_.Get(
		j,
		"propagateTags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) Role() *string {
	var returns *string
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) SchedulingStrategy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"schedulingStrategy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) ServiceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) ServiceRegistries() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"serviceRegistries",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) TaskDefinition() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskDefinition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnService) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ECS::Service`.
func NewCfnService(scope awscdk.Construct, id *string, props *CfnServiceProps) CfnService {
	_init_.Initialize()

	j := jsiiProxy_CfnService{}

	_jsii_.Create(
		"monocdk.aws_ecs.CfnService",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ECS::Service`.
func NewCfnService_Override(c CfnService, scope awscdk.Construct, id *string, props *CfnServiceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.CfnService",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnService) SetCapacityProviderStrategy(val interface{}) {
	_jsii_.Set(
		j,
		"capacityProviderStrategy",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetCluster(val *string) {
	_jsii_.Set(
		j,
		"cluster",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetDeploymentConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"deploymentConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetDeploymentController(val interface{}) {
	_jsii_.Set(
		j,
		"deploymentController",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetDesiredCount(val *float64) {
	_jsii_.Set(
		j,
		"desiredCount",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetEnableEcsManagedTags(val interface{}) {
	_jsii_.Set(
		j,
		"enableEcsManagedTags",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetEnableExecuteCommand(val interface{}) {
	_jsii_.Set(
		j,
		"enableExecuteCommand",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetHealthCheckGracePeriodSeconds(val *float64) {
	_jsii_.Set(
		j,
		"healthCheckGracePeriodSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetLaunchType(val *string) {
	_jsii_.Set(
		j,
		"launchType",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetLoadBalancers(val interface{}) {
	_jsii_.Set(
		j,
		"loadBalancers",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetNetworkConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"networkConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetPlacementConstraints(val interface{}) {
	_jsii_.Set(
		j,
		"placementConstraints",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetPlacementStrategies(val interface{}) {
	_jsii_.Set(
		j,
		"placementStrategies",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetPlatformVersion(val *string) {
	_jsii_.Set(
		j,
		"platformVersion",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetPropagateTags(val *string) {
	_jsii_.Set(
		j,
		"propagateTags",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetRole(val *string) {
	_jsii_.Set(
		j,
		"role",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetSchedulingStrategy(val *string) {
	_jsii_.Set(
		j,
		"schedulingStrategy",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetServiceName(val *string) {
	_jsii_.Set(
		j,
		"serviceName",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetServiceRegistries(val interface{}) {
	_jsii_.Set(
		j,
		"serviceRegistries",
		val,
	)
}

func (j *jsiiProxy_CfnService) SetTaskDefinition(val *string) {
	_jsii_.Set(
		j,
		"taskDefinition",
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
func CfnService_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnService",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnService_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnService",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnService_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnService",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnService_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.CfnService",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnService) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnService) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnService) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnService) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnService) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnService) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnService) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnService) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnService) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnService) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnService) OnPrepare() {
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
func (c *jsiiProxy_CfnService) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnService) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnService) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnService) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnService) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnService) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnService) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnService) ToString() *string {
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
func (c *jsiiProxy_CfnService) Validate() *[]*string {
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
func (c *jsiiProxy_CfnService) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnService_AwsVpcConfigurationProperty struct {
	// `CfnService.AwsVpcConfigurationProperty.Subnets`.
	Subnets *[]*string `json:"subnets"`
	// `CfnService.AwsVpcConfigurationProperty.AssignPublicIp`.
	AssignPublicIp *string `json:"assignPublicIp"`
	// `CfnService.AwsVpcConfigurationProperty.SecurityGroups`.
	SecurityGroups *[]*string `json:"securityGroups"`
}

type CfnService_CapacityProviderStrategyItemProperty struct {
	// `CfnService.CapacityProviderStrategyItemProperty.Base`.
	Base *float64 `json:"base"`
	// `CfnService.CapacityProviderStrategyItemProperty.CapacityProvider`.
	CapacityProvider *string `json:"capacityProvider"`
	// `CfnService.CapacityProviderStrategyItemProperty.Weight`.
	Weight *float64 `json:"weight"`
}

type CfnService_DeploymentCircuitBreakerProperty struct {
	// `CfnService.DeploymentCircuitBreakerProperty.Enable`.
	Enable interface{} `json:"enable"`
	// `CfnService.DeploymentCircuitBreakerProperty.Rollback`.
	Rollback interface{} `json:"rollback"`
}

type CfnService_DeploymentConfigurationProperty struct {
	// `CfnService.DeploymentConfigurationProperty.DeploymentCircuitBreaker`.
	DeploymentCircuitBreaker interface{} `json:"deploymentCircuitBreaker"`
	// `CfnService.DeploymentConfigurationProperty.MaximumPercent`.
	MaximumPercent *float64 `json:"maximumPercent"`
	// `CfnService.DeploymentConfigurationProperty.MinimumHealthyPercent`.
	MinimumHealthyPercent *float64 `json:"minimumHealthyPercent"`
}

type CfnService_DeploymentControllerProperty struct {
	// `CfnService.DeploymentControllerProperty.Type`.
	Type *string `json:"type"`
}

type CfnService_LoadBalancerProperty struct {
	// `CfnService.LoadBalancerProperty.ContainerPort`.
	ContainerPort *float64 `json:"containerPort"`
	// `CfnService.LoadBalancerProperty.ContainerName`.
	ContainerName *string `json:"containerName"`
	// `CfnService.LoadBalancerProperty.LoadBalancerName`.
	LoadBalancerName *string `json:"loadBalancerName"`
	// `CfnService.LoadBalancerProperty.TargetGroupArn`.
	TargetGroupArn *string `json:"targetGroupArn"`
}

type CfnService_NetworkConfigurationProperty struct {
	// `CfnService.NetworkConfigurationProperty.AwsvpcConfiguration`.
	AwsvpcConfiguration interface{} `json:"awsvpcConfiguration"`
}

type CfnService_PlacementConstraintProperty struct {
	// `CfnService.PlacementConstraintProperty.Type`.
	Type *string `json:"type"`
	// `CfnService.PlacementConstraintProperty.Expression`.
	Expression *string `json:"expression"`
}

type CfnService_PlacementStrategyProperty struct {
	// `CfnService.PlacementStrategyProperty.Type`.
	Type *string `json:"type"`
	// `CfnService.PlacementStrategyProperty.Field`.
	Field *string `json:"field"`
}

type CfnService_ServiceRegistryProperty struct {
	// `CfnService.ServiceRegistryProperty.ContainerName`.
	ContainerName *string `json:"containerName"`
	// `CfnService.ServiceRegistryProperty.ContainerPort`.
	ContainerPort *float64 `json:"containerPort"`
	// `CfnService.ServiceRegistryProperty.Port`.
	Port *float64 `json:"port"`
	// `CfnService.ServiceRegistryProperty.RegistryArn`.
	RegistryArn *string `json:"registryArn"`
}

// Properties for defining a `AWS::ECS::Service`.
type CfnServiceProps struct {
	// `AWS::ECS::Service.CapacityProviderStrategy`.
	CapacityProviderStrategy interface{} `json:"capacityProviderStrategy"`
	// `AWS::ECS::Service.Cluster`.
	Cluster *string `json:"cluster"`
	// `AWS::ECS::Service.DeploymentConfiguration`.
	DeploymentConfiguration interface{} `json:"deploymentConfiguration"`
	// `AWS::ECS::Service.DeploymentController`.
	DeploymentController interface{} `json:"deploymentController"`
	// `AWS::ECS::Service.DesiredCount`.
	DesiredCount *float64 `json:"desiredCount"`
	// `AWS::ECS::Service.EnableECSManagedTags`.
	EnableEcsManagedTags interface{} `json:"enableEcsManagedTags"`
	// `AWS::ECS::Service.EnableExecuteCommand`.
	EnableExecuteCommand interface{} `json:"enableExecuteCommand"`
	// `AWS::ECS::Service.HealthCheckGracePeriodSeconds`.
	HealthCheckGracePeriodSeconds *float64 `json:"healthCheckGracePeriodSeconds"`
	// `AWS::ECS::Service.LaunchType`.
	LaunchType *string `json:"launchType"`
	// `AWS::ECS::Service.LoadBalancers`.
	LoadBalancers interface{} `json:"loadBalancers"`
	// `AWS::ECS::Service.NetworkConfiguration`.
	NetworkConfiguration interface{} `json:"networkConfiguration"`
	// `AWS::ECS::Service.PlacementConstraints`.
	PlacementConstraints interface{} `json:"placementConstraints"`
	// `AWS::ECS::Service.PlacementStrategies`.
	PlacementStrategies interface{} `json:"placementStrategies"`
	// `AWS::ECS::Service.PlatformVersion`.
	PlatformVersion *string `json:"platformVersion"`
	// `AWS::ECS::Service.PropagateTags`.
	PropagateTags *string `json:"propagateTags"`
	// `AWS::ECS::Service.Role`.
	Role *string `json:"role"`
	// `AWS::ECS::Service.SchedulingStrategy`.
	SchedulingStrategy *string `json:"schedulingStrategy"`
	// `AWS::ECS::Service.ServiceName`.
	ServiceName *string `json:"serviceName"`
	// `AWS::ECS::Service.ServiceRegistries`.
	ServiceRegistries interface{} `json:"serviceRegistries"`
	// `AWS::ECS::Service.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::ECS::Service.TaskDefinition`.
	TaskDefinition *string `json:"taskDefinition"`
}

// A CloudFormation `AWS::ECS::TaskDefinition`.
type CfnTaskDefinition interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrTaskDefinitionArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ContainerDefinitions() interface{}
	SetContainerDefinitions(val interface{})
	Cpu() *string
	SetCpu(val *string)
	CreationStack() *[]*string
	EphemeralStorage() interface{}
	SetEphemeralStorage(val interface{})
	ExecutionRoleArn() *string
	SetExecutionRoleArn(val *string)
	Family() *string
	SetFamily(val *string)
	InferenceAccelerators() interface{}
	SetInferenceAccelerators(val interface{})
	IpcMode() *string
	SetIpcMode(val *string)
	LogicalId() *string
	Memory() *string
	SetMemory(val *string)
	NetworkMode() *string
	SetNetworkMode(val *string)
	Node() awscdk.ConstructNode
	PidMode() *string
	SetPidMode(val *string)
	PlacementConstraints() interface{}
	SetPlacementConstraints(val interface{})
	ProxyConfiguration() interface{}
	SetProxyConfiguration(val interface{})
	Ref() *string
	RequiresCompatibilities() *[]*string
	SetRequiresCompatibilities(val *[]*string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	TaskRoleArn() *string
	SetTaskRoleArn(val *string)
	UpdatedProperites() *map[string]interface{}
	Volumes() interface{}
	SetVolumes(val interface{})
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

// The jsii proxy struct for CfnTaskDefinition
type jsiiProxy_CfnTaskDefinition struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnTaskDefinition) AttrTaskDefinitionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrTaskDefinitionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) ContainerDefinitions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"containerDefinitions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Cpu() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cpu",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) EphemeralStorage() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ephemeralStorage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) ExecutionRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"executionRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Family() *string {
	var returns *string
	_jsii_.Get(
		j,
		"family",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) InferenceAccelerators() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"inferenceAccelerators",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) IpcMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ipcMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Memory() *string {
	var returns *string
	_jsii_.Get(
		j,
		"memory",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) NetworkMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) PidMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pidMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) PlacementConstraints() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"placementConstraints",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) ProxyConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"proxyConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) RequiresCompatibilities() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"requiresCompatibilities",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) TaskRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskDefinition) Volumes() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"volumes",
		&returns,
	)
	return returns
}


// Create a new `AWS::ECS::TaskDefinition`.
func NewCfnTaskDefinition(scope awscdk.Construct, id *string, props *CfnTaskDefinitionProps) CfnTaskDefinition {
	_init_.Initialize()

	j := jsiiProxy_CfnTaskDefinition{}

	_jsii_.Create(
		"monocdk.aws_ecs.CfnTaskDefinition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ECS::TaskDefinition`.
func NewCfnTaskDefinition_Override(c CfnTaskDefinition, scope awscdk.Construct, id *string, props *CfnTaskDefinitionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.CfnTaskDefinition",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetContainerDefinitions(val interface{}) {
	_jsii_.Set(
		j,
		"containerDefinitions",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetCpu(val *string) {
	_jsii_.Set(
		j,
		"cpu",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetEphemeralStorage(val interface{}) {
	_jsii_.Set(
		j,
		"ephemeralStorage",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetExecutionRoleArn(val *string) {
	_jsii_.Set(
		j,
		"executionRoleArn",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetFamily(val *string) {
	_jsii_.Set(
		j,
		"family",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetInferenceAccelerators(val interface{}) {
	_jsii_.Set(
		j,
		"inferenceAccelerators",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetIpcMode(val *string) {
	_jsii_.Set(
		j,
		"ipcMode",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetMemory(val *string) {
	_jsii_.Set(
		j,
		"memory",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetNetworkMode(val *string) {
	_jsii_.Set(
		j,
		"networkMode",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetPidMode(val *string) {
	_jsii_.Set(
		j,
		"pidMode",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetPlacementConstraints(val interface{}) {
	_jsii_.Set(
		j,
		"placementConstraints",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetProxyConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"proxyConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetRequiresCompatibilities(val *[]*string) {
	_jsii_.Set(
		j,
		"requiresCompatibilities",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetTaskRoleArn(val *string) {
	_jsii_.Set(
		j,
		"taskRoleArn",
		val,
	)
}

func (j *jsiiProxy_CfnTaskDefinition) SetVolumes(val interface{}) {
	_jsii_.Set(
		j,
		"volumes",
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
func CfnTaskDefinition_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnTaskDefinition",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnTaskDefinition_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnTaskDefinition",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnTaskDefinition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnTaskDefinition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnTaskDefinition_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.CfnTaskDefinition",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnTaskDefinition) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnTaskDefinition) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnTaskDefinition) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnTaskDefinition) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnTaskDefinition) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnTaskDefinition) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnTaskDefinition) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnTaskDefinition) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnTaskDefinition) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnTaskDefinition) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnTaskDefinition) OnPrepare() {
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
func (c *jsiiProxy_CfnTaskDefinition) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTaskDefinition) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnTaskDefinition) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnTaskDefinition) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnTaskDefinition) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnTaskDefinition) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnTaskDefinition) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTaskDefinition) ToString() *string {
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
func (c *jsiiProxy_CfnTaskDefinition) Validate() *[]*string {
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
func (c *jsiiProxy_CfnTaskDefinition) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnTaskDefinition_AuthorizationConfigProperty struct {
	// `CfnTaskDefinition.AuthorizationConfigProperty.AccessPointId`.
	AccessPointId *string `json:"accessPointId"`
	// `CfnTaskDefinition.AuthorizationConfigProperty.IAM`.
	Iam *string `json:"iam"`
}

type CfnTaskDefinition_ContainerDefinitionProperty struct {
	// `CfnTaskDefinition.ContainerDefinitionProperty.Command`.
	Command *[]*string `json:"command"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Cpu`.
	Cpu *float64 `json:"cpu"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.DependsOn`.
	DependsOn interface{} `json:"dependsOn"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.DisableNetworking`.
	DisableNetworking interface{} `json:"disableNetworking"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.DnsSearchDomains`.
	DnsSearchDomains *[]*string `json:"dnsSearchDomains"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.DnsServers`.
	DnsServers *[]*string `json:"dnsServers"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.DockerLabels`.
	DockerLabels interface{} `json:"dockerLabels"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.DockerSecurityOptions`.
	DockerSecurityOptions *[]*string `json:"dockerSecurityOptions"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.EntryPoint`.
	EntryPoint *[]*string `json:"entryPoint"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Environment`.
	Environment interface{} `json:"environment"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.EnvironmentFiles`.
	EnvironmentFiles interface{} `json:"environmentFiles"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Essential`.
	Essential interface{} `json:"essential"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.ExtraHosts`.
	ExtraHosts interface{} `json:"extraHosts"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.FirelensConfiguration`.
	FirelensConfiguration interface{} `json:"firelensConfiguration"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.HealthCheck`.
	HealthCheck interface{} `json:"healthCheck"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Hostname`.
	Hostname *string `json:"hostname"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Image`.
	Image *string `json:"image"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Interactive`.
	Interactive interface{} `json:"interactive"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Links`.
	Links *[]*string `json:"links"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.LinuxParameters`.
	LinuxParameters interface{} `json:"linuxParameters"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.LogConfiguration`.
	LogConfiguration interface{} `json:"logConfiguration"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Memory`.
	Memory *float64 `json:"memory"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.MemoryReservation`.
	MemoryReservation *float64 `json:"memoryReservation"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.MountPoints`.
	MountPoints interface{} `json:"mountPoints"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Name`.
	Name *string `json:"name"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.PortMappings`.
	PortMappings interface{} `json:"portMappings"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Privileged`.
	Privileged interface{} `json:"privileged"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.PseudoTerminal`.
	PseudoTerminal interface{} `json:"pseudoTerminal"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.ReadonlyRootFilesystem`.
	ReadonlyRootFilesystem interface{} `json:"readonlyRootFilesystem"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.RepositoryCredentials`.
	RepositoryCredentials interface{} `json:"repositoryCredentials"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.ResourceRequirements`.
	ResourceRequirements interface{} `json:"resourceRequirements"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Secrets`.
	Secrets interface{} `json:"secrets"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.StartTimeout`.
	StartTimeout *float64 `json:"startTimeout"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.StopTimeout`.
	StopTimeout *float64 `json:"stopTimeout"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.SystemControls`.
	SystemControls interface{} `json:"systemControls"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.Ulimits`.
	Ulimits interface{} `json:"ulimits"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.User`.
	User *string `json:"user"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.VolumesFrom`.
	VolumesFrom interface{} `json:"volumesFrom"`
	// `CfnTaskDefinition.ContainerDefinitionProperty.WorkingDirectory`.
	WorkingDirectory *string `json:"workingDirectory"`
}

type CfnTaskDefinition_ContainerDependencyProperty struct {
	// `CfnTaskDefinition.ContainerDependencyProperty.Condition`.
	Condition *string `json:"condition"`
	// `CfnTaskDefinition.ContainerDependencyProperty.ContainerName`.
	ContainerName *string `json:"containerName"`
}

type CfnTaskDefinition_DeviceProperty struct {
	// `CfnTaskDefinition.DeviceProperty.ContainerPath`.
	ContainerPath *string `json:"containerPath"`
	// `CfnTaskDefinition.DeviceProperty.HostPath`.
	HostPath *string `json:"hostPath"`
	// `CfnTaskDefinition.DeviceProperty.Permissions`.
	Permissions *[]*string `json:"permissions"`
}

type CfnTaskDefinition_DockerVolumeConfigurationProperty struct {
	// `CfnTaskDefinition.DockerVolumeConfigurationProperty.Autoprovision`.
	Autoprovision interface{} `json:"autoprovision"`
	// `CfnTaskDefinition.DockerVolumeConfigurationProperty.Driver`.
	Driver *string `json:"driver"`
	// `CfnTaskDefinition.DockerVolumeConfigurationProperty.DriverOpts`.
	DriverOpts interface{} `json:"driverOpts"`
	// `CfnTaskDefinition.DockerVolumeConfigurationProperty.Labels`.
	Labels interface{} `json:"labels"`
	// `CfnTaskDefinition.DockerVolumeConfigurationProperty.Scope`.
	Scope *string `json:"scope"`
}

type CfnTaskDefinition_EfsVolumeConfigurationProperty struct {
	// `CfnTaskDefinition.EfsVolumeConfigurationProperty.FileSystemId`.
	FileSystemId *string `json:"fileSystemId"`
	// `CfnTaskDefinition.EfsVolumeConfigurationProperty.AuthorizationConfig`.
	AuthorizationConfig interface{} `json:"authorizationConfig"`
	// `CfnTaskDefinition.EfsVolumeConfigurationProperty.RootDirectory`.
	RootDirectory *string `json:"rootDirectory"`
	// `CfnTaskDefinition.EfsVolumeConfigurationProperty.TransitEncryption`.
	TransitEncryption *string `json:"transitEncryption"`
	// `CfnTaskDefinition.EfsVolumeConfigurationProperty.TransitEncryptionPort`.
	TransitEncryptionPort *float64 `json:"transitEncryptionPort"`
}

type CfnTaskDefinition_EnvironmentFileProperty struct {
	// `CfnTaskDefinition.EnvironmentFileProperty.Type`.
	Type *string `json:"type"`
	// `CfnTaskDefinition.EnvironmentFileProperty.Value`.
	Value *string `json:"value"`
}

type CfnTaskDefinition_EphemeralStorageProperty struct {
	// `CfnTaskDefinition.EphemeralStorageProperty.SizeInGiB`.
	SizeInGiB *float64 `json:"sizeInGiB"`
}

type CfnTaskDefinition_FirelensConfigurationProperty struct {
	// `CfnTaskDefinition.FirelensConfigurationProperty.Options`.
	Options interface{} `json:"options"`
	// `CfnTaskDefinition.FirelensConfigurationProperty.Type`.
	Type *string `json:"type"`
}

type CfnTaskDefinition_HealthCheckProperty struct {
	// `CfnTaskDefinition.HealthCheckProperty.Command`.
	Command *[]*string `json:"command"`
	// `CfnTaskDefinition.HealthCheckProperty.Interval`.
	Interval *float64 `json:"interval"`
	// `CfnTaskDefinition.HealthCheckProperty.Retries`.
	Retries *float64 `json:"retries"`
	// `CfnTaskDefinition.HealthCheckProperty.StartPeriod`.
	StartPeriod *float64 `json:"startPeriod"`
	// `CfnTaskDefinition.HealthCheckProperty.Timeout`.
	Timeout *float64 `json:"timeout"`
}

type CfnTaskDefinition_HostEntryProperty struct {
	// `CfnTaskDefinition.HostEntryProperty.Hostname`.
	Hostname *string `json:"hostname"`
	// `CfnTaskDefinition.HostEntryProperty.IpAddress`.
	IpAddress *string `json:"ipAddress"`
}

type CfnTaskDefinition_HostVolumePropertiesProperty struct {
	// `CfnTaskDefinition.HostVolumePropertiesProperty.SourcePath`.
	SourcePath *string `json:"sourcePath"`
}

type CfnTaskDefinition_InferenceAcceleratorProperty struct {
	// `CfnTaskDefinition.InferenceAcceleratorProperty.DeviceName`.
	DeviceName *string `json:"deviceName"`
	// `CfnTaskDefinition.InferenceAcceleratorProperty.DeviceType`.
	DeviceType *string `json:"deviceType"`
}

type CfnTaskDefinition_KernelCapabilitiesProperty struct {
	// `CfnTaskDefinition.KernelCapabilitiesProperty.Add`.
	Add *[]*string `json:"add"`
	// `CfnTaskDefinition.KernelCapabilitiesProperty.Drop`.
	Drop *[]*string `json:"drop"`
}

type CfnTaskDefinition_KeyValuePairProperty struct {
	// `CfnTaskDefinition.KeyValuePairProperty.Name`.
	Name *string `json:"name"`
	// `CfnTaskDefinition.KeyValuePairProperty.Value`.
	Value *string `json:"value"`
}

type CfnTaskDefinition_LinuxParametersProperty struct {
	// `CfnTaskDefinition.LinuxParametersProperty.Capabilities`.
	Capabilities interface{} `json:"capabilities"`
	// `CfnTaskDefinition.LinuxParametersProperty.Devices`.
	Devices interface{} `json:"devices"`
	// `CfnTaskDefinition.LinuxParametersProperty.InitProcessEnabled`.
	InitProcessEnabled interface{} `json:"initProcessEnabled"`
	// `CfnTaskDefinition.LinuxParametersProperty.MaxSwap`.
	MaxSwap *float64 `json:"maxSwap"`
	// `CfnTaskDefinition.LinuxParametersProperty.SharedMemorySize`.
	SharedMemorySize *float64 `json:"sharedMemorySize"`
	// `CfnTaskDefinition.LinuxParametersProperty.Swappiness`.
	Swappiness *float64 `json:"swappiness"`
	// `CfnTaskDefinition.LinuxParametersProperty.Tmpfs`.
	Tmpfs interface{} `json:"tmpfs"`
}

type CfnTaskDefinition_LogConfigurationProperty struct {
	// `CfnTaskDefinition.LogConfigurationProperty.LogDriver`.
	LogDriver *string `json:"logDriver"`
	// `CfnTaskDefinition.LogConfigurationProperty.Options`.
	Options interface{} `json:"options"`
	// `CfnTaskDefinition.LogConfigurationProperty.SecretOptions`.
	SecretOptions interface{} `json:"secretOptions"`
}

type CfnTaskDefinition_MountPointProperty struct {
	// `CfnTaskDefinition.MountPointProperty.ContainerPath`.
	ContainerPath *string `json:"containerPath"`
	// `CfnTaskDefinition.MountPointProperty.ReadOnly`.
	ReadOnly interface{} `json:"readOnly"`
	// `CfnTaskDefinition.MountPointProperty.SourceVolume`.
	SourceVolume *string `json:"sourceVolume"`
}

type CfnTaskDefinition_PortMappingProperty struct {
	// `CfnTaskDefinition.PortMappingProperty.ContainerPort`.
	ContainerPort *float64 `json:"containerPort"`
	// `CfnTaskDefinition.PortMappingProperty.HostPort`.
	HostPort *float64 `json:"hostPort"`
	// `CfnTaskDefinition.PortMappingProperty.Protocol`.
	Protocol *string `json:"protocol"`
}

type CfnTaskDefinition_ProxyConfigurationProperty struct {
	// `CfnTaskDefinition.ProxyConfigurationProperty.ContainerName`.
	ContainerName *string `json:"containerName"`
	// `CfnTaskDefinition.ProxyConfigurationProperty.ProxyConfigurationProperties`.
	ProxyConfigurationProperties interface{} `json:"proxyConfigurationProperties"`
	// `CfnTaskDefinition.ProxyConfigurationProperty.Type`.
	Type *string `json:"type"`
}

type CfnTaskDefinition_RepositoryCredentialsProperty struct {
	// `CfnTaskDefinition.RepositoryCredentialsProperty.CredentialsParameter`.
	CredentialsParameter *string `json:"credentialsParameter"`
}

type CfnTaskDefinition_ResourceRequirementProperty struct {
	// `CfnTaskDefinition.ResourceRequirementProperty.Type`.
	Type *string `json:"type"`
	// `CfnTaskDefinition.ResourceRequirementProperty.Value`.
	Value *string `json:"value"`
}

type CfnTaskDefinition_SecretProperty struct {
	// `CfnTaskDefinition.SecretProperty.Name`.
	Name *string `json:"name"`
	// `CfnTaskDefinition.SecretProperty.ValueFrom`.
	ValueFrom *string `json:"valueFrom"`
}

type CfnTaskDefinition_SystemControlProperty struct {
	// `CfnTaskDefinition.SystemControlProperty.Namespace`.
	Namespace *string `json:"namespace"`
	// `CfnTaskDefinition.SystemControlProperty.Value`.
	Value *string `json:"value"`
}

type CfnTaskDefinition_TaskDefinitionPlacementConstraintProperty struct {
	// `CfnTaskDefinition.TaskDefinitionPlacementConstraintProperty.Type`.
	Type *string `json:"type"`
	// `CfnTaskDefinition.TaskDefinitionPlacementConstraintProperty.Expression`.
	Expression *string `json:"expression"`
}

type CfnTaskDefinition_TmpfsProperty struct {
	// `CfnTaskDefinition.TmpfsProperty.Size`.
	Size *float64 `json:"size"`
	// `CfnTaskDefinition.TmpfsProperty.ContainerPath`.
	ContainerPath *string `json:"containerPath"`
	// `CfnTaskDefinition.TmpfsProperty.MountOptions`.
	MountOptions *[]*string `json:"mountOptions"`
}

type CfnTaskDefinition_UlimitProperty struct {
	// `CfnTaskDefinition.UlimitProperty.HardLimit`.
	HardLimit *float64 `json:"hardLimit"`
	// `CfnTaskDefinition.UlimitProperty.Name`.
	Name *string `json:"name"`
	// `CfnTaskDefinition.UlimitProperty.SoftLimit`.
	SoftLimit *float64 `json:"softLimit"`
}

type CfnTaskDefinition_VolumeFromProperty struct {
	// `CfnTaskDefinition.VolumeFromProperty.ReadOnly`.
	ReadOnly interface{} `json:"readOnly"`
	// `CfnTaskDefinition.VolumeFromProperty.SourceContainer`.
	SourceContainer *string `json:"sourceContainer"`
}

type CfnTaskDefinition_VolumeProperty struct {
	// `CfnTaskDefinition.VolumeProperty.DockerVolumeConfiguration`.
	DockerVolumeConfiguration interface{} `json:"dockerVolumeConfiguration"`
	// `CfnTaskDefinition.VolumeProperty.EfsVolumeConfiguration`.
	EfsVolumeConfiguration interface{} `json:"efsVolumeConfiguration"`
	// `CfnTaskDefinition.VolumeProperty.Host`.
	Host interface{} `json:"host"`
	// `CfnTaskDefinition.VolumeProperty.Name`.
	Name *string `json:"name"`
}

// Properties for defining a `AWS::ECS::TaskDefinition`.
type CfnTaskDefinitionProps struct {
	// `AWS::ECS::TaskDefinition.ContainerDefinitions`.
	ContainerDefinitions interface{} `json:"containerDefinitions"`
	// `AWS::ECS::TaskDefinition.Cpu`.
	Cpu *string `json:"cpu"`
	// `AWS::ECS::TaskDefinition.EphemeralStorage`.
	EphemeralStorage interface{} `json:"ephemeralStorage"`
	// `AWS::ECS::TaskDefinition.ExecutionRoleArn`.
	ExecutionRoleArn *string `json:"executionRoleArn"`
	// `AWS::ECS::TaskDefinition.Family`.
	Family *string `json:"family"`
	// `AWS::ECS::TaskDefinition.InferenceAccelerators`.
	InferenceAccelerators interface{} `json:"inferenceAccelerators"`
	// `AWS::ECS::TaskDefinition.IpcMode`.
	IpcMode *string `json:"ipcMode"`
	// `AWS::ECS::TaskDefinition.Memory`.
	Memory *string `json:"memory"`
	// `AWS::ECS::TaskDefinition.NetworkMode`.
	NetworkMode *string `json:"networkMode"`
	// `AWS::ECS::TaskDefinition.PidMode`.
	PidMode *string `json:"pidMode"`
	// `AWS::ECS::TaskDefinition.PlacementConstraints`.
	PlacementConstraints interface{} `json:"placementConstraints"`
	// `AWS::ECS::TaskDefinition.ProxyConfiguration`.
	ProxyConfiguration interface{} `json:"proxyConfiguration"`
	// `AWS::ECS::TaskDefinition.RequiresCompatibilities`.
	RequiresCompatibilities *[]*string `json:"requiresCompatibilities"`
	// `AWS::ECS::TaskDefinition.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::ECS::TaskDefinition.TaskRoleArn`.
	TaskRoleArn *string `json:"taskRoleArn"`
	// `AWS::ECS::TaskDefinition.Volumes`.
	Volumes interface{} `json:"volumes"`
}

// A CloudFormation `AWS::ECS::TaskSet`.
type CfnTaskSet interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Cluster() *string
	SetCluster(val *string)
	CreationStack() *[]*string
	ExternalId() *string
	SetExternalId(val *string)
	LaunchType() *string
	SetLaunchType(val *string)
	LoadBalancers() interface{}
	SetLoadBalancers(val interface{})
	LogicalId() *string
	NetworkConfiguration() interface{}
	SetNetworkConfiguration(val interface{})
	Node() awscdk.ConstructNode
	PlatformVersion() *string
	SetPlatformVersion(val *string)
	Ref() *string
	Scale() interface{}
	SetScale(val interface{})
	Service() *string
	SetService(val *string)
	ServiceRegistries() interface{}
	SetServiceRegistries(val interface{})
	Stack() awscdk.Stack
	TaskDefinition() *string
	SetTaskDefinition(val *string)
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

// The jsii proxy struct for CfnTaskSet
type jsiiProxy_CfnTaskSet struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnTaskSet) AttrId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) Cluster() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) ExternalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"externalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) LaunchType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"launchType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) LoadBalancers() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"loadBalancers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) NetworkConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"networkConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) PlatformVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"platformVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) Scale() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"scale",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) Service() *string {
	var returns *string
	_jsii_.Get(
		j,
		"service",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) ServiceRegistries() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"serviceRegistries",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) TaskDefinition() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskDefinition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTaskSet) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ECS::TaskSet`.
func NewCfnTaskSet(scope awscdk.Construct, id *string, props *CfnTaskSetProps) CfnTaskSet {
	_init_.Initialize()

	j := jsiiProxy_CfnTaskSet{}

	_jsii_.Create(
		"monocdk.aws_ecs.CfnTaskSet",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ECS::TaskSet`.
func NewCfnTaskSet_Override(c CfnTaskSet, scope awscdk.Construct, id *string, props *CfnTaskSetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.CfnTaskSet",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetCluster(val *string) {
	_jsii_.Set(
		j,
		"cluster",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetExternalId(val *string) {
	_jsii_.Set(
		j,
		"externalId",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetLaunchType(val *string) {
	_jsii_.Set(
		j,
		"launchType",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetLoadBalancers(val interface{}) {
	_jsii_.Set(
		j,
		"loadBalancers",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetNetworkConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"networkConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetPlatformVersion(val *string) {
	_jsii_.Set(
		j,
		"platformVersion",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetScale(val interface{}) {
	_jsii_.Set(
		j,
		"scale",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetService(val *string) {
	_jsii_.Set(
		j,
		"service",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetServiceRegistries(val interface{}) {
	_jsii_.Set(
		j,
		"serviceRegistries",
		val,
	)
}

func (j *jsiiProxy_CfnTaskSet) SetTaskDefinition(val *string) {
	_jsii_.Set(
		j,
		"taskDefinition",
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
func CfnTaskSet_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnTaskSet",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnTaskSet_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnTaskSet",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnTaskSet_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.CfnTaskSet",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnTaskSet_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_ecs.CfnTaskSet",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnTaskSet) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnTaskSet) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnTaskSet) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnTaskSet) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnTaskSet) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnTaskSet) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnTaskSet) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnTaskSet) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnTaskSet) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnTaskSet) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnTaskSet) OnPrepare() {
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
func (c *jsiiProxy_CfnTaskSet) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTaskSet) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnTaskSet) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnTaskSet) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnTaskSet) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnTaskSet) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnTaskSet) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTaskSet) ToString() *string {
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
func (c *jsiiProxy_CfnTaskSet) Validate() *[]*string {
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
func (c *jsiiProxy_CfnTaskSet) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnTaskSet_AwsVpcConfigurationProperty struct {
	// `CfnTaskSet.AwsVpcConfigurationProperty.Subnets`.
	Subnets *[]*string `json:"subnets"`
	// `CfnTaskSet.AwsVpcConfigurationProperty.AssignPublicIp`.
	AssignPublicIp *string `json:"assignPublicIp"`
	// `CfnTaskSet.AwsVpcConfigurationProperty.SecurityGroups`.
	SecurityGroups *[]*string `json:"securityGroups"`
}

type CfnTaskSet_LoadBalancerProperty struct {
	// `CfnTaskSet.LoadBalancerProperty.ContainerName`.
	ContainerName *string `json:"containerName"`
	// `CfnTaskSet.LoadBalancerProperty.ContainerPort`.
	ContainerPort *float64 `json:"containerPort"`
	// `CfnTaskSet.LoadBalancerProperty.LoadBalancerName`.
	LoadBalancerName *string `json:"loadBalancerName"`
	// `CfnTaskSet.LoadBalancerProperty.TargetGroupArn`.
	TargetGroupArn *string `json:"targetGroupArn"`
}

type CfnTaskSet_NetworkConfigurationProperty struct {
	// `CfnTaskSet.NetworkConfigurationProperty.AwsVpcConfiguration`.
	AwsVpcConfiguration interface{} `json:"awsVpcConfiguration"`
}

type CfnTaskSet_ScaleProperty struct {
	// `CfnTaskSet.ScaleProperty.Unit`.
	Unit *string `json:"unit"`
	// `CfnTaskSet.ScaleProperty.Value`.
	Value *float64 `json:"value"`
}

type CfnTaskSet_ServiceRegistryProperty struct {
	// `CfnTaskSet.ServiceRegistryProperty.ContainerName`.
	ContainerName *string `json:"containerName"`
	// `CfnTaskSet.ServiceRegistryProperty.ContainerPort`.
	ContainerPort *float64 `json:"containerPort"`
	// `CfnTaskSet.ServiceRegistryProperty.Port`.
	Port *float64 `json:"port"`
	// `CfnTaskSet.ServiceRegistryProperty.RegistryArn`.
	RegistryArn *string `json:"registryArn"`
}

// Properties for defining a `AWS::ECS::TaskSet`.
type CfnTaskSetProps struct {
	// `AWS::ECS::TaskSet.Cluster`.
	Cluster *string `json:"cluster"`
	// `AWS::ECS::TaskSet.Service`.
	Service *string `json:"service"`
	// `AWS::ECS::TaskSet.TaskDefinition`.
	TaskDefinition *string `json:"taskDefinition"`
	// `AWS::ECS::TaskSet.ExternalId`.
	ExternalId *string `json:"externalId"`
	// `AWS::ECS::TaskSet.LaunchType`.
	LaunchType *string `json:"launchType"`
	// `AWS::ECS::TaskSet.LoadBalancers`.
	LoadBalancers interface{} `json:"loadBalancers"`
	// `AWS::ECS::TaskSet.NetworkConfiguration`.
	NetworkConfiguration interface{} `json:"networkConfiguration"`
	// `AWS::ECS::TaskSet.PlatformVersion`.
	PlatformVersion *string `json:"platformVersion"`
	// `AWS::ECS::TaskSet.Scale`.
	Scale interface{} `json:"scale"`
	// `AWS::ECS::TaskSet.ServiceRegistries`.
	ServiceRegistries interface{} `json:"serviceRegistries"`
}

// The options for creating an AWS Cloud Map namespace.
// Experimental.
type CloudMapNamespaceOptions struct {
	// The name of the namespace, such as example.com.
	// Experimental.
	Name *string `json:"name"`
	// The type of CloudMap Namespace to create.
	// Experimental.
	Type awsservicediscovery.NamespaceType `json:"type"`
	// The VPC to associate the namespace with.
	//
	// This property is required for private DNS namespaces.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
}

// The options to enabling AWS Cloud Map for an Amazon ECS service.
// Experimental.
type CloudMapOptions struct {
	// The service discovery namespace for the Cloud Map service to attach to the ECS service.
	// Experimental.
	CloudMapNamespace awsservicediscovery.INamespace `json:"cloudMapNamespace"`
	// The container to point to for a SRV record.
	// Experimental.
	Container ContainerDefinition `json:"container"`
	// The port to point to for a SRV record.
	// Experimental.
	ContainerPort *float64 `json:"containerPort"`
	// The DNS record type that you want AWS Cloud Map to create.
	//
	// The supported record types are A or SRV.
	// Experimental.
	DnsRecordType awsservicediscovery.DnsRecordType `json:"dnsRecordType"`
	// The amount of time that you want DNS resolvers to cache the settings for this record.
	// Experimental.
	DnsTtl awscdk.Duration `json:"dnsTtl"`
	// The number of 30-second intervals that you want Cloud Map to wait after receiving an UpdateInstanceCustomHealthStatus request before it changes the health status of a service instance.
	//
	// NOTE: This is used for HealthCheckCustomConfig
	// Experimental.
	FailureThreshold *float64 `json:"failureThreshold"`
	// The name of the Cloud Map service to attach to the ECS service.
	// Experimental.
	Name *string `json:"name"`
}

// A regional grouping of one or more container instances on which you can run tasks and services.
// Experimental.
type Cluster interface {
	awscdk.Resource
	ICluster
	AutoscalingGroup() awsautoscaling.IAutoScalingGroup
	ClusterArn() *string
	ClusterName() *string
	Connections() awsec2.Connections
	DefaultCloudMapNamespace() awsservicediscovery.INamespace
	Env() *awscdk.ResourceEnvironment
	ExecuteCommandConfiguration() *ExecuteCommandConfiguration
	HasEc2Capacity() *bool
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	Vpc() awsec2.IVpc
	AddAsgCapacityProvider(provider AsgCapacityProvider, options *AddAutoScalingGroupCapacityOptions)
	AddAutoScalingGroup(autoScalingGroup awsautoscaling.AutoScalingGroup, options *AddAutoScalingGroupCapacityOptions)
	AddCapacity(id *string, options *AddCapacityOptions) awsautoscaling.AutoScalingGroup
	AddCapacityProvider(provider *string)
	AddDefaultCloudMapNamespace(options *CloudMapNamespaceOptions) awsservicediscovery.INamespace
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	EnableFargateCapacityProviders()
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCpuReservation(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricMemoryReservation(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Cluster
type jsiiProxy_Cluster struct {
	internal.Type__awscdkResource
	jsiiProxy_ICluster
}

func (j *jsiiProxy_Cluster) AutoscalingGroup() awsautoscaling.IAutoScalingGroup {
	var returns awsautoscaling.IAutoScalingGroup
	_jsii_.Get(
		j,
		"autoscalingGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) ClusterArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clusterArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) ClusterName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clusterName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) DefaultCloudMapNamespace() awsservicediscovery.INamespace {
	var returns awsservicediscovery.INamespace
	_jsii_.Get(
		j,
		"defaultCloudMapNamespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) ExecuteCommandConfiguration() *ExecuteCommandConfiguration {
	var returns *ExecuteCommandConfiguration
	_jsii_.Get(
		j,
		"executeCommandConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) HasEc2Capacity() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"hasEc2Capacity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Cluster) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}


// Constructs a new instance of the Cluster class.
// Experimental.
func NewCluster(scope constructs.Construct, id *string, props *ClusterProps) Cluster {
	_init_.Initialize()

	j := jsiiProxy_Cluster{}

	_jsii_.Create(
		"monocdk.aws_ecs.Cluster",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the Cluster class.
// Experimental.
func NewCluster_Override(c Cluster, scope constructs.Construct, id *string, props *ClusterProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.Cluster",
		[]interface{}{scope, id, props},
		c,
	)
}

// Import an existing cluster to the stack from its attributes.
// Experimental.
func Cluster_FromClusterAttributes(scope constructs.Construct, id *string, attrs *ClusterAttributes) ICluster {
	_init_.Initialize()

	var returns ICluster

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Cluster",
		"fromClusterAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Cluster_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Cluster",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Cluster_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Cluster",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// This method adds an Auto Scaling Group Capacity Provider to a cluster.
// Experimental.
func (c *jsiiProxy_Cluster) AddAsgCapacityProvider(provider AsgCapacityProvider, options *AddAutoScalingGroupCapacityOptions) {
	_jsii_.InvokeVoid(
		c,
		"addAsgCapacityProvider",
		[]interface{}{provider, options},
	)
}

// This method adds compute capacity to a cluster using the specified AutoScalingGroup.
// Deprecated: Use {@link Cluster.addAsgCapacityProvider} instead.
func (c *jsiiProxy_Cluster) AddAutoScalingGroup(autoScalingGroup awsautoscaling.AutoScalingGroup, options *AddAutoScalingGroupCapacityOptions) {
	_jsii_.InvokeVoid(
		c,
		"addAutoScalingGroup",
		[]interface{}{autoScalingGroup, options},
	)
}

// This method adds compute capacity to a cluster by creating an AutoScalingGroup with the specified options.
//
// Returns the AutoScalingGroup so you can add autoscaling settings to it.
// Deprecated: Use {@link Cluster.addAsgCapacityProvider} instead.
func (c *jsiiProxy_Cluster) AddCapacity(id *string, options *AddCapacityOptions) awsautoscaling.AutoScalingGroup {
	var returns awsautoscaling.AutoScalingGroup

	_jsii_.Invoke(
		c,
		"addCapacity",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// This method enables the Fargate or Fargate Spot capacity providers on the cluster.
// See: {@link addAsgCapacityProvider} to add an Auto Scaling Group capacity provider to the cluster.
//
// Deprecated: Use {@link enableFargateCapacityProviders} instead.
func (c *jsiiProxy_Cluster) AddCapacityProvider(provider *string) {
	_jsii_.InvokeVoid(
		c,
		"addCapacityProvider",
		[]interface{}{provider},
	)
}

// Add an AWS Cloud Map DNS namespace for this cluster.
//
// NOTE: HttpNamespaces are not supported, as ECS always requires a DNSConfig when registering an instance to a Cloud
// Map service.
// Experimental.
func (c *jsiiProxy_Cluster) AddDefaultCloudMapNamespace(options *CloudMapNamespaceOptions) awsservicediscovery.INamespace {
	var returns awsservicediscovery.INamespace

	_jsii_.Invoke(
		c,
		"addDefaultCloudMapNamespace",
		[]interface{}{options},
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
func (c *jsiiProxy_Cluster) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		c,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Enable the Fargate capacity providers for this cluster.
// Experimental.
func (c *jsiiProxy_Cluster) EnableFargateCapacityProviders() {
	_jsii_.InvokeVoid(
		c,
		"enableFargateCapacityProviders",
		nil, // no parameters
	)
}

// Experimental.
func (c *jsiiProxy_Cluster) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_Cluster) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_Cluster) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// This method returns the specifed CloudWatch metric for this cluster.
// Experimental.
func (c *jsiiProxy_Cluster) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		c,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this clusters CPU reservation.
// Experimental.
func (c *jsiiProxy_Cluster) MetricCpuReservation(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		c,
		"metricCpuReservation",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this clusters CPU utilization.
// Experimental.
func (c *jsiiProxy_Cluster) MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		c,
		"metricCpuUtilization",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this clusters memory reservation.
// Experimental.
func (c *jsiiProxy_Cluster) MetricMemoryReservation(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		c,
		"metricMemoryReservation",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this clusters memory utilization.
// Experimental.
func (c *jsiiProxy_Cluster) MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		c,
		"metricMemoryUtilization",
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
func (c *jsiiProxy_Cluster) OnPrepare() {
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
func (c *jsiiProxy_Cluster) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_Cluster) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_Cluster) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_Cluster) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_Cluster) ToString() *string {
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
func (c *jsiiProxy_Cluster) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The properties to import from the ECS cluster.
// Experimental.
type ClusterAttributes struct {
	// The name of the cluster.
	// Experimental.
	ClusterName *string `json:"clusterName"`
	// The security groups associated with the container instances registered to the cluster.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The VPC associated with the cluster.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Autoscaling group added to the cluster if capacity is added.
	// Experimental.
	AutoscalingGroup awsautoscaling.IAutoScalingGroup `json:"autoscalingGroup"`
	// The Amazon Resource Name (ARN) that identifies the cluster.
	// Experimental.
	ClusterArn *string `json:"clusterArn"`
	// The AWS Cloud Map namespace to associate with the cluster.
	// Experimental.
	DefaultCloudMapNamespace awsservicediscovery.INamespace `json:"defaultCloudMapNamespace"`
	// The execute command configuration for the cluster.
	// Experimental.
	ExecuteCommandConfiguration *ExecuteCommandConfiguration `json:"executeCommandConfiguration"`
	// Specifies whether the cluster has EC2 instance capacity.
	// Experimental.
	HasEc2Capacity *bool `json:"hasEc2Capacity"`
}

// The properties used to define an ECS cluster.
// Experimental.
type ClusterProps struct {
	// The ec2 capacity to add to the cluster.
	// Experimental.
	Capacity *AddCapacityOptions `json:"capacity"`
	// The capacity providers to add to the cluster.
	// Deprecated: Use {@link ClusterProps.enableFargateCapacityProviders} instead.
	CapacityProviders *[]*string `json:"capacityProviders"`
	// The name for the cluster.
	// Experimental.
	ClusterName *string `json:"clusterName"`
	// If true CloudWatch Container Insights will be enabled for the cluster.
	// Experimental.
	ContainerInsights *bool `json:"containerInsights"`
	// The service discovery namespace created in this cluster.
	// Experimental.
	DefaultCloudMapNamespace *CloudMapNamespaceOptions `json:"defaultCloudMapNamespace"`
	// Whether to enable Fargate Capacity Providers.
	// Experimental.
	EnableFargateCapacityProviders *bool `json:"enableFargateCapacityProviders"`
	// The execute command configuration for the cluster.
	// Experimental.
	ExecuteCommandConfiguration *ExecuteCommandConfiguration `json:"executeCommandConfiguration"`
	// The VPC where your ECS instances will be running or your ENIs will be deployed.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
}

// The common task definition attributes used across all types of task definitions.
// Experimental.
type CommonTaskDefinitionAttributes struct {
	// The arn of the task definition.
	// Experimental.
	TaskDefinitionArn *string `json:"taskDefinitionArn"`
	// The networking mode to use for the containers in the task.
	// Experimental.
	NetworkMode NetworkMode `json:"networkMode"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
}

// The common properties for all task definitions.
//
// For more information, see
// [Task Definition Parameters](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html).
// Experimental.
type CommonTaskDefinitionProps struct {
	// The name of the IAM task execution role that grants the ECS agent to call AWS APIs on your behalf.
	//
	// The role will be used to retrieve container images from ECR and create CloudWatch log groups.
	// Experimental.
	ExecutionRole awsiam.IRole `json:"executionRole"`
	// The name of a family that this task definition is registered to.
	//
	// A family groups multiple versions of a task definition.
	// Experimental.
	Family *string `json:"family"`
	// The configuration details for the App Mesh proxy.
	// Experimental.
	ProxyConfiguration ProxyConfiguration `json:"proxyConfiguration"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
	// The list of volume definitions for the task.
	//
	// For more information, see
	// [Task Definition Parameter Volumes](https://docs.aws.amazon.com/AmazonECS/latest/developerguide//task_definition_parameters.html#volumes).
	// Experimental.
	Volumes *[]*Volume `json:"volumes"`
}

// The task launch type compatibility requirement.
// Experimental.
type Compatibility string

const (
	Compatibility_EC2 Compatibility = "EC2"
	Compatibility_FARGATE Compatibility = "FARGATE"
	Compatibility_EC2_AND_FARGATE Compatibility = "EC2_AND_FARGATE"
	Compatibility_EXTERNAL Compatibility = "EXTERNAL"
)

// A container definition is used in a task definition to describe the containers that are launched as part of a task.
// Experimental.
type ContainerDefinition interface {
	awscdk.Construct
	ContainerDependencies() *[]*ContainerDependency
	ContainerName() *string
	ContainerPort() *float64
	EnvironmentFiles() *[]*EnvironmentFileConfig
	Essential() *bool
	IngressPort() *float64
	LinuxParameters() LinuxParameters
	LogDriverConfig() *LogDriverConfig
	MemoryLimitSpecified() *bool
	MountPoints() *[]*MountPoint
	Node() awscdk.ConstructNode
	PortMappings() *[]*PortMapping
	ReferencesSecretJsonField() *bool
	TaskDefinition() TaskDefinition
	Ulimits() *[]*Ulimit
	VolumesFrom() *[]*VolumeFrom
	AddContainerDependencies(containerDependencies ...*ContainerDependency)
	AddInferenceAcceleratorResource(inferenceAcceleratorResources ...*string)
	AddLink(container ContainerDefinition, alias *string)
	AddMountPoints(mountPoints ...*MountPoint)
	AddPortMappings(portMappings ...*PortMapping)
	AddScratch(scratch *ScratchSpace)
	AddToExecutionPolicy(statement awsiam.PolicyStatement)
	AddUlimits(ulimits ...*Ulimit)
	AddVolumesFrom(volumesFrom ...*VolumeFrom)
	FindPortMapping(containerPort *float64, protocol Protocol) *PortMapping
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderContainerDefinition(_taskDefinition TaskDefinition) *CfnTaskDefinition_ContainerDefinitionProperty
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ContainerDefinition
type jsiiProxy_ContainerDefinition struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_ContainerDefinition) ContainerDependencies() *[]*ContainerDependency {
	var returns *[]*ContainerDependency
	_jsii_.Get(
		j,
		"containerDependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) ContainerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"containerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) ContainerPort() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"containerPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) EnvironmentFiles() *[]*EnvironmentFileConfig {
	var returns *[]*EnvironmentFileConfig
	_jsii_.Get(
		j,
		"environmentFiles",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) Essential() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"essential",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) IngressPort() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"ingressPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) LinuxParameters() LinuxParameters {
	var returns LinuxParameters
	_jsii_.Get(
		j,
		"linuxParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) LogDriverConfig() *LogDriverConfig {
	var returns *LogDriverConfig
	_jsii_.Get(
		j,
		"logDriverConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) MemoryLimitSpecified() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"memoryLimitSpecified",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) MountPoints() *[]*MountPoint {
	var returns *[]*MountPoint
	_jsii_.Get(
		j,
		"mountPoints",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) PortMappings() *[]*PortMapping {
	var returns *[]*PortMapping
	_jsii_.Get(
		j,
		"portMappings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) ReferencesSecretJsonField() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"referencesSecretJsonField",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) TaskDefinition() TaskDefinition {
	var returns TaskDefinition
	_jsii_.Get(
		j,
		"taskDefinition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) Ulimits() *[]*Ulimit {
	var returns *[]*Ulimit
	_jsii_.Get(
		j,
		"ulimits",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ContainerDefinition) VolumesFrom() *[]*VolumeFrom {
	var returns *[]*VolumeFrom
	_jsii_.Get(
		j,
		"volumesFrom",
		&returns,
	)
	return returns
}


// Constructs a new instance of the ContainerDefinition class.
// Experimental.
func NewContainerDefinition(scope constructs.Construct, id *string, props *ContainerDefinitionProps) ContainerDefinition {
	_init_.Initialize()

	j := jsiiProxy_ContainerDefinition{}

	_jsii_.Create(
		"monocdk.aws_ecs.ContainerDefinition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the ContainerDefinition class.
// Experimental.
func NewContainerDefinition_Override(c ContainerDefinition, scope constructs.Construct, id *string, props *ContainerDefinitionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.ContainerDefinition",
		[]interface{}{scope, id, props},
		c,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func ContainerDefinition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ContainerDefinition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// This method adds one or more container dependencies to the container.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddContainerDependencies(containerDependencies ...*ContainerDependency) {
	args := []interface{}{}
	for _, a := range containerDependencies {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addContainerDependencies",
		args,
	)
}

// This method adds one or more resources to the container.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddInferenceAcceleratorResource(inferenceAcceleratorResources ...*string) {
	args := []interface{}{}
	for _, a := range inferenceAcceleratorResources {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addInferenceAcceleratorResource",
		args,
	)
}

// This method adds a link which allows containers to communicate with each other without the need for port mappings.
//
// This parameter is only supported if the task definition is using the bridge network mode.
// Warning: The --link flag is a legacy feature of Docker. It may eventually be removed.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddLink(container ContainerDefinition, alias *string) {
	_jsii_.InvokeVoid(
		c,
		"addLink",
		[]interface{}{container, alias},
	)
}

// This method adds one or more mount points for data volumes to the container.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddMountPoints(mountPoints ...*MountPoint) {
	args := []interface{}{}
	for _, a := range mountPoints {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addMountPoints",
		args,
	)
}

// This method adds one or more port mappings to the container.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddPortMappings(portMappings ...*PortMapping) {
	args := []interface{}{}
	for _, a := range portMappings {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addPortMappings",
		args,
	)
}

// This method mounts temporary disk space to the container.
//
// This adds the correct container mountPoint and task definition volume.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddScratch(scratch *ScratchSpace) {
	_jsii_.InvokeVoid(
		c,
		"addScratch",
		[]interface{}{scratch},
	)
}

// This method adds the specified statement to the IAM task execution policy in the task definition.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddToExecutionPolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		c,
		"addToExecutionPolicy",
		[]interface{}{statement},
	)
}

// This method adds one or more ulimits to the container.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddUlimits(ulimits ...*Ulimit) {
	args := []interface{}{}
	for _, a := range ulimits {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addUlimits",
		args,
	)
}

// This method adds one or more volumes to the container.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) AddVolumesFrom(volumesFrom ...*VolumeFrom) {
	args := []interface{}{}
	for _, a := range volumesFrom {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addVolumesFrom",
		args,
	)
}

// Returns the host port for the requested container port if it exists.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) FindPortMapping(containerPort *float64, protocol Protocol) *PortMapping {
	var returns *PortMapping

	_jsii_.Invoke(
		c,
		"findPortMapping",
		[]interface{}{containerPort, protocol},
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
func (c *jsiiProxy_ContainerDefinition) OnPrepare() {
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
func (c *jsiiProxy_ContainerDefinition) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_ContainerDefinition) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_ContainerDefinition) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Render this container definition to a CloudFormation object.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) RenderContainerDefinition(_taskDefinition TaskDefinition) *CfnTaskDefinition_ContainerDefinitionProperty {
	var returns *CfnTaskDefinition_ContainerDefinitionProperty

	_jsii_.Invoke(
		c,
		"renderContainerDefinition",
		[]interface{}{_taskDefinition},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_ContainerDefinition) ToString() *string {
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
func (c *jsiiProxy_ContainerDefinition) Validate() *[]*string {
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
type ContainerDefinitionOptions struct {
	// The image used to start a container.
	//
	// This string is passed directly to the Docker daemon.
	// Images in the Docker Hub registry are available by default.
	// Other repositories are specified with either repository-url/image:tag or repository-url/image@digest.
	// TODO: Update these to specify using classes of IContainerImage
	// Experimental.
	Image ContainerImage `json:"image"`
	// The command that is passed to the container.
	//
	// If you provide a shell command as a single string, you have to quote command-line arguments.
	// Experimental.
	Command *[]*string `json:"command"`
	// The name of the container.
	// Experimental.
	ContainerName *string `json:"containerName"`
	// The minimum number of CPU units to reserve for the container.
	// Experimental.
	Cpu *float64 `json:"cpu"`
	// Specifies whether networking is disabled within the container.
	//
	// When this parameter is true, networking is disabled within the container.
	// Experimental.
	DisableNetworking *bool `json:"disableNetworking"`
	// A list of DNS search domains that are presented to the container.
	// Experimental.
	DnsSearchDomains *[]*string `json:"dnsSearchDomains"`
	// A list of DNS servers that are presented to the container.
	// Experimental.
	DnsServers *[]*string `json:"dnsServers"`
	// A key/value map of labels to add to the container.
	// Experimental.
	DockerLabels *map[string]*string `json:"dockerLabels"`
	// A list of strings to provide custom labels for SELinux and AppArmor multi-level security systems.
	// Experimental.
	DockerSecurityOptions *[]*string `json:"dockerSecurityOptions"`
	// The ENTRYPOINT value to pass to the container.
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	EntryPoint *[]*string `json:"entryPoint"`
	// The environment variables to pass to the container.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The environment files to pass to the container.
	// See: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/taskdef-envfiles.html
	//
	// Experimental.
	EnvironmentFiles *[]EnvironmentFile `json:"environmentFiles"`
	// Specifies whether the container is marked essential.
	//
	// If the essential parameter of a container is marked as true, and that container fails
	// or stops for any reason, all other containers that are part of the task are stopped.
	// If the essential parameter of a container is marked as false, then its failure does not
	// affect the rest of the containers in a task. All tasks must have at least one essential container.
	//
	// If this parameter is omitted, a container is assumed to be essential.
	// Experimental.
	Essential *bool `json:"essential"`
	// A list of hostnames and IP address mappings to append to the /etc/hosts file on the container.
	// Experimental.
	ExtraHosts *map[string]*string `json:"extraHosts"`
	// The number of GPUs assigned to the container.
	// Experimental.
	GpuCount *float64 `json:"gpuCount"`
	// The health check command and associated configuration parameters for the container.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The hostname to use for your container.
	// Experimental.
	Hostname *string `json:"hostname"`
	// The inference accelerators referenced by the container.
	// Experimental.
	InferenceAcceleratorResources *[]*string `json:"inferenceAcceleratorResources"`
	// Linux-specific modifications that are applied to the container, such as Linux kernel capabilities.
	//
	// For more information see [KernelCapabilities](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_KernelCapabilities.html).
	// Experimental.
	LinuxParameters LinuxParameters `json:"linuxParameters"`
	// The log configuration specification for the container.
	// Experimental.
	Logging LogDriver `json:"logging"`
	// The amount (in MiB) of memory to present to the container.
	//
	// If your container attempts to exceed the allocated memory, the container
	// is terminated.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryLimitMiB *float64 `json:"memoryLimitMiB"`
	// The soft limit (in MiB) of memory to reserve for the container.
	//
	// When system memory is under heavy contention, Docker attempts to keep the
	// container memory to this soft limit. However, your container can consume more
	// memory when it needs to, up to either the hard limit specified with the memory
	// parameter (if applicable), or all of the available memory on the container
	// instance, whichever comes first.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryReservationMiB *float64 `json:"memoryReservationMiB"`
	// The port mappings to add to the container definition.
	// Experimental.
	PortMappings *[]*PortMapping `json:"portMappings"`
	// Specifies whether the container is marked as privileged.
	//
	// When this parameter is true, the container is given elevated privileges on the host container instance (similar to the root user).
	// Experimental.
	Privileged *bool `json:"privileged"`
	// When this parameter is true, the container is given read-only access to its root file system.
	// Experimental.
	ReadonlyRootFilesystem *bool `json:"readonlyRootFilesystem"`
	// The secret environment variables to pass to the container.
	// Experimental.
	Secrets *map[string]Secret `json:"secrets"`
	// Time duration (in seconds) to wait before giving up on resolving dependencies for a container.
	// Experimental.
	StartTimeout awscdk.Duration `json:"startTimeout"`
	// Time duration (in seconds) to wait before the container is forcefully killed if it doesn't exit normally on its own.
	// Experimental.
	StopTimeout awscdk.Duration `json:"stopTimeout"`
	// The user name to use inside the container.
	// Experimental.
	User *string `json:"user"`
	// The working directory in which to run commands inside the container.
	// Experimental.
	WorkingDirectory *string `json:"workingDirectory"`
}

// The properties in a container definition.
// Experimental.
type ContainerDefinitionProps struct {
	// The image used to start a container.
	//
	// This string is passed directly to the Docker daemon.
	// Images in the Docker Hub registry are available by default.
	// Other repositories are specified with either repository-url/image:tag or repository-url/image@digest.
	// TODO: Update these to specify using classes of IContainerImage
	// Experimental.
	Image ContainerImage `json:"image"`
	// The command that is passed to the container.
	//
	// If you provide a shell command as a single string, you have to quote command-line arguments.
	// Experimental.
	Command *[]*string `json:"command"`
	// The name of the container.
	// Experimental.
	ContainerName *string `json:"containerName"`
	// The minimum number of CPU units to reserve for the container.
	// Experimental.
	Cpu *float64 `json:"cpu"`
	// Specifies whether networking is disabled within the container.
	//
	// When this parameter is true, networking is disabled within the container.
	// Experimental.
	DisableNetworking *bool `json:"disableNetworking"`
	// A list of DNS search domains that are presented to the container.
	// Experimental.
	DnsSearchDomains *[]*string `json:"dnsSearchDomains"`
	// A list of DNS servers that are presented to the container.
	// Experimental.
	DnsServers *[]*string `json:"dnsServers"`
	// A key/value map of labels to add to the container.
	// Experimental.
	DockerLabels *map[string]*string `json:"dockerLabels"`
	// A list of strings to provide custom labels for SELinux and AppArmor multi-level security systems.
	// Experimental.
	DockerSecurityOptions *[]*string `json:"dockerSecurityOptions"`
	// The ENTRYPOINT value to pass to the container.
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	EntryPoint *[]*string `json:"entryPoint"`
	// The environment variables to pass to the container.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The environment files to pass to the container.
	// See: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/taskdef-envfiles.html
	//
	// Experimental.
	EnvironmentFiles *[]EnvironmentFile `json:"environmentFiles"`
	// Specifies whether the container is marked essential.
	//
	// If the essential parameter of a container is marked as true, and that container fails
	// or stops for any reason, all other containers that are part of the task are stopped.
	// If the essential parameter of a container is marked as false, then its failure does not
	// affect the rest of the containers in a task. All tasks must have at least one essential container.
	//
	// If this parameter is omitted, a container is assumed to be essential.
	// Experimental.
	Essential *bool `json:"essential"`
	// A list of hostnames and IP address mappings to append to the /etc/hosts file on the container.
	// Experimental.
	ExtraHosts *map[string]*string `json:"extraHosts"`
	// The number of GPUs assigned to the container.
	// Experimental.
	GpuCount *float64 `json:"gpuCount"`
	// The health check command and associated configuration parameters for the container.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The hostname to use for your container.
	// Experimental.
	Hostname *string `json:"hostname"`
	// The inference accelerators referenced by the container.
	// Experimental.
	InferenceAcceleratorResources *[]*string `json:"inferenceAcceleratorResources"`
	// Linux-specific modifications that are applied to the container, such as Linux kernel capabilities.
	//
	// For more information see [KernelCapabilities](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_KernelCapabilities.html).
	// Experimental.
	LinuxParameters LinuxParameters `json:"linuxParameters"`
	// The log configuration specification for the container.
	// Experimental.
	Logging LogDriver `json:"logging"`
	// The amount (in MiB) of memory to present to the container.
	//
	// If your container attempts to exceed the allocated memory, the container
	// is terminated.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryLimitMiB *float64 `json:"memoryLimitMiB"`
	// The soft limit (in MiB) of memory to reserve for the container.
	//
	// When system memory is under heavy contention, Docker attempts to keep the
	// container memory to this soft limit. However, your container can consume more
	// memory when it needs to, up to either the hard limit specified with the memory
	// parameter (if applicable), or all of the available memory on the container
	// instance, whichever comes first.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryReservationMiB *float64 `json:"memoryReservationMiB"`
	// The port mappings to add to the container definition.
	// Experimental.
	PortMappings *[]*PortMapping `json:"portMappings"`
	// Specifies whether the container is marked as privileged.
	//
	// When this parameter is true, the container is given elevated privileges on the host container instance (similar to the root user).
	// Experimental.
	Privileged *bool `json:"privileged"`
	// When this parameter is true, the container is given read-only access to its root file system.
	// Experimental.
	ReadonlyRootFilesystem *bool `json:"readonlyRootFilesystem"`
	// The secret environment variables to pass to the container.
	// Experimental.
	Secrets *map[string]Secret `json:"secrets"`
	// Time duration (in seconds) to wait before giving up on resolving dependencies for a container.
	// Experimental.
	StartTimeout awscdk.Duration `json:"startTimeout"`
	// Time duration (in seconds) to wait before the container is forcefully killed if it doesn't exit normally on its own.
	// Experimental.
	StopTimeout awscdk.Duration `json:"stopTimeout"`
	// The user name to use inside the container.
	// Experimental.
	User *string `json:"user"`
	// The working directory in which to run commands inside the container.
	// Experimental.
	WorkingDirectory *string `json:"workingDirectory"`
	// The name of the task definition that includes this container definition.
	//
	// [disable-awslint:ref-via-interface]
	// Experimental.
	TaskDefinition TaskDefinition `json:"taskDefinition"`
}

// The details of a dependency on another container in the task definition.
// See: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ContainerDependency.html
//
// Experimental.
type ContainerDependency struct {
	// The container to depend on.
	// Experimental.
	Container ContainerDefinition `json:"container"`
	// The state the container needs to be in to satisfy the dependency and proceed with startup.
	//
	// Valid values are ContainerDependencyCondition.START, ContainerDependencyCondition.COMPLETE,
	// ContainerDependencyCondition.SUCCESS and ContainerDependencyCondition.HEALTHY.
	// Experimental.
	Condition ContainerDependencyCondition `json:"condition"`
}

// Experimental.
type ContainerDependencyCondition string

const (
	ContainerDependencyCondition_START ContainerDependencyCondition = "START"
	ContainerDependencyCondition_COMPLETE ContainerDependencyCondition = "COMPLETE"
	ContainerDependencyCondition_SUCCESS ContainerDependencyCondition = "SUCCESS"
	ContainerDependencyCondition_HEALTHY ContainerDependencyCondition = "HEALTHY"
)

// Constructs for types of container images.
// Experimental.
type ContainerImage interface {
	Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig
}

// The jsii proxy struct for ContainerImage
type jsiiProxy_ContainerImage struct {
	_ byte // padding
}

// Experimental.
func NewContainerImage_Override(c ContainerImage) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.ContainerImage",
		nil, // no parameters
		c,
	)
}

// Reference an image that's constructed directly from sources on disk.
//
// If you already have a `DockerImageAsset` instance, you can use the
// `ContainerImage.fromDockerImageAsset` method instead.
// Experimental.
func ContainerImage_FromAsset(directory *string, props *AssetImageProps) AssetImage {
	_init_.Initialize()

	var returns AssetImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ContainerImage",
		"fromAsset",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Use an existing `DockerImageAsset` for this container image.
// Experimental.
func ContainerImage_FromDockerImageAsset(asset awsecrassets.DockerImageAsset) ContainerImage {
	_init_.Initialize()

	var returns ContainerImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ContainerImage",
		"fromDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Reference an image in an ECR repository.
// Experimental.
func ContainerImage_FromEcrRepository(repository awsecr.IRepository, tag *string) EcrImage {
	_init_.Initialize()

	var returns EcrImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ContainerImage",
		"fromEcrRepository",
		[]interface{}{repository, tag},
		&returns,
	)

	return returns
}

// Reference an image on DockerHub or another online registry.
// Experimental.
func ContainerImage_FromRegistry(name *string, props *RepositoryImageProps) RepositoryImage {
	_init_.Initialize()

	var returns RepositoryImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ContainerImage",
		"fromRegistry",
		[]interface{}{name, props},
		&returns,
	)

	return returns
}

// Called when the image is used by a ContainerDefinition.
// Experimental.
func (c *jsiiProxy_ContainerImage) Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig {
	var returns *ContainerImageConfig

	_jsii_.Invoke(
		c,
		"bind",
		[]interface{}{scope, containerDefinition},
		&returns,
	)

	return returns
}

// The configuration for creating a container image.
// Experimental.
type ContainerImageConfig struct {
	// Specifies the name of the container image.
	// Experimental.
	ImageName *string `json:"imageName"`
	// Specifies the credentials used to access the image repository.
	// Experimental.
	RepositoryCredentials *CfnTaskDefinition_RepositoryCredentialsProperty `json:"repositoryCredentials"`
}

// The properties for enabling scaling based on CPU utilization.
// Experimental.
type CpuUtilizationScalingProps struct {
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the scalable resource. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// scalable resource.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Period after a scale in activity completes before another scale in activity can start.
	// Experimental.
	ScaleInCooldown awscdk.Duration `json:"scaleInCooldown"`
	// Period after a scale out activity completes before another scale out activity can start.
	// Experimental.
	ScaleOutCooldown awscdk.Duration `json:"scaleOutCooldown"`
	// The target value for CPU utilization across all tasks in the service.
	// Experimental.
	TargetUtilizationPercent *float64 `json:"targetUtilizationPercent"`
}

// The deployment circuit breaker to use for the service.
// Experimental.
type DeploymentCircuitBreaker struct {
	// Whether to enable rollback on deployment failure.
	// Experimental.
	Rollback *bool `json:"rollback"`
}

// The deployment controller to use for the service.
// Experimental.
type DeploymentController struct {
	// The deployment controller type to use.
	// Experimental.
	Type DeploymentControllerType `json:"type"`
}

// The deployment controller type to use for the service.
// Experimental.
type DeploymentControllerType string

const (
	DeploymentControllerType_ECS DeploymentControllerType = "ECS"
	DeploymentControllerType_CODE_DEPLOY DeploymentControllerType = "CODE_DEPLOY"
	DeploymentControllerType_EXTERNAL DeploymentControllerType = "EXTERNAL"
)

// A container instance host device.
// Experimental.
type Device struct {
	// The path for the device on the host container instance.
	// Experimental.
	HostPath *string `json:"hostPath"`
	// The path inside the container at which to expose the host device.
	// Experimental.
	ContainerPath *string `json:"containerPath"`
	// The explicit permissions to provide to the container for the device.
	//
	// By default, the container has permissions for read, write, and mknod for the device.
	// Experimental.
	Permissions *[]DevicePermission `json:"permissions"`
}

// Permissions for device access.
// Experimental.
type DevicePermission string

const (
	DevicePermission_READ DevicePermission = "READ"
	DevicePermission_WRITE DevicePermission = "WRITE"
	DevicePermission_MKNOD DevicePermission = "MKNOD"
)

// The configuration for a Docker volume.
//
// Docker volumes are only supported when you are using the EC2 launch type.
// Experimental.
type DockerVolumeConfiguration struct {
	// The Docker volume driver to use.
	// Experimental.
	Driver *string `json:"driver"`
	// The scope for the Docker volume that determines its lifecycle.
	// Experimental.
	Scope Scope `json:"scope"`
	// Specifies whether the Docker volume should be created if it does not already exist.
	//
	// If true is specified, the Docker volume will be created for you.
	// Experimental.
	Autoprovision *bool `json:"autoprovision"`
	// A map of Docker driver-specific options passed through.
	// Experimental.
	DriverOpts *map[string]*string `json:"driverOpts"`
	// Custom metadata to add to your Docker volume.
	// Experimental.
	Labels *map[string]*string `json:"labels"`
}

// This creates a service using the EC2 launch type on an ECS cluster.
// Experimental.
type Ec2Service interface {
	BaseService
	IEc2Service
	CloudmapService() awsservicediscovery.Service
	SetCloudmapService(val awsservicediscovery.Service)
	CloudMapService() awsservicediscovery.IService
	Cluster() ICluster
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	LoadBalancers() *[]*CfnService_LoadBalancerProperty
	SetLoadBalancers(val *[]*CfnService_LoadBalancerProperty)
	NetworkConfiguration() *CfnService_NetworkConfigurationProperty
	SetNetworkConfiguration(val *CfnService_NetworkConfigurationProperty)
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ServiceArn() *string
	ServiceName() *string
	ServiceRegistries() *[]*CfnService_ServiceRegistryProperty
	SetServiceRegistries(val *[]*CfnService_ServiceRegistryProperty)
	Stack() awscdk.Stack
	TaskDefinition() TaskDefinition
	AddPlacementConstraints(constraints ...PlacementConstraint)
	AddPlacementStrategies(strategies ...PlacementStrategy)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AssociateCloudMapService(options *AssociateCloudMapServiceOptions)
	AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer)
	AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	AutoScaleTaskCount(props *awsapplicationautoscaling.EnableScalingProps) ScalableTaskCount
	ConfigureAwsVpcNetworking(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroup awsec2.ISecurityGroup)
	ConfigureAwsVpcNetworkingWithSecurityGroups(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroups *[]awsec2.ISecurityGroup)
	EnableCloudMap(options *CloudMapOptions) awsservicediscovery.Service
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LoadBalancerTarget(options *LoadBalancerTargetOptions) IEcsLoadBalancerTarget
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RegisterLoadBalancerTargets(targets ...*EcsTarget)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Ec2Service
type jsiiProxy_Ec2Service struct {
	jsiiProxy_BaseService
	jsiiProxy_IEc2Service
}

func (j *jsiiProxy_Ec2Service) CloudmapService() awsservicediscovery.Service {
	var returns awsservicediscovery.Service
	_jsii_.Get(
		j,
		"cloudmapService",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) CloudMapService() awsservicediscovery.IService {
	var returns awsservicediscovery.IService
	_jsii_.Get(
		j,
		"cloudMapService",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) Cluster() ICluster {
	var returns ICluster
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) LoadBalancers() *[]*CfnService_LoadBalancerProperty {
	var returns *[]*CfnService_LoadBalancerProperty
	_jsii_.Get(
		j,
		"loadBalancers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) NetworkConfiguration() *CfnService_NetworkConfigurationProperty {
	var returns *CfnService_NetworkConfigurationProperty
	_jsii_.Get(
		j,
		"networkConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) ServiceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) ServiceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) ServiceRegistries() *[]*CfnService_ServiceRegistryProperty {
	var returns *[]*CfnService_ServiceRegistryProperty
	_jsii_.Get(
		j,
		"serviceRegistries",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2Service) TaskDefinition() TaskDefinition {
	var returns TaskDefinition
	_jsii_.Get(
		j,
		"taskDefinition",
		&returns,
	)
	return returns
}


// Constructs a new instance of the Ec2Service class.
// Experimental.
func NewEc2Service(scope constructs.Construct, id *string, props *Ec2ServiceProps) Ec2Service {
	_init_.Initialize()

	j := jsiiProxy_Ec2Service{}

	_jsii_.Create(
		"monocdk.aws_ecs.Ec2Service",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the Ec2Service class.
// Experimental.
func NewEc2Service_Override(e Ec2Service, scope constructs.Construct, id *string, props *Ec2ServiceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.Ec2Service",
		[]interface{}{scope, id, props},
		e,
	)
}

func (j *jsiiProxy_Ec2Service) SetCloudmapService(val awsservicediscovery.Service) {
	_jsii_.Set(
		j,
		"cloudmapService",
		val,
	)
}

func (j *jsiiProxy_Ec2Service) SetLoadBalancers(val *[]*CfnService_LoadBalancerProperty) {
	_jsii_.Set(
		j,
		"loadBalancers",
		val,
	)
}

func (j *jsiiProxy_Ec2Service) SetNetworkConfiguration(val *CfnService_NetworkConfigurationProperty) {
	_jsii_.Set(
		j,
		"networkConfiguration",
		val,
	)
}

func (j *jsiiProxy_Ec2Service) SetServiceRegistries(val *[]*CfnService_ServiceRegistryProperty) {
	_jsii_.Set(
		j,
		"serviceRegistries",
		val,
	)
}

// Imports from the specified service ARN.
// Experimental.
func Ec2Service_FromEc2ServiceArn(scope constructs.Construct, id *string, ec2ServiceArn *string) IEc2Service {
	_init_.Initialize()

	var returns IEc2Service

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2Service",
		"fromEc2ServiceArn",
		[]interface{}{scope, id, ec2ServiceArn},
		&returns,
	)

	return returns
}

// Imports from the specified service attrributes.
// Experimental.
func Ec2Service_FromEc2ServiceAttributes(scope constructs.Construct, id *string, attrs *Ec2ServiceAttributes) IBaseService {
	_init_.Initialize()

	var returns IBaseService

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2Service",
		"fromEc2ServiceAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Ec2Service_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2Service",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Ec2Service_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2Service",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds one or more placement constraints to use for tasks in the service.
//
// For more information, see
// [Amazon ECS Task Placement Constraints](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-placement-constraints.html).
// Experimental.
func (e *jsiiProxy_Ec2Service) AddPlacementConstraints(constraints ...PlacementConstraint) {
	args := []interface{}{}
	for _, a := range constraints {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		e,
		"addPlacementConstraints",
		args,
	)
}

// Adds one or more placement strategies to use for tasks in the service.
//
// For more information, see
// [Amazon ECS Task Placement Strategies](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-placement-strategies.html).
// Experimental.
func (e *jsiiProxy_Ec2Service) AddPlacementStrategies(strategies ...PlacementStrategy) {
	args := []interface{}{}
	for _, a := range strategies {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		e,
		"addPlacementStrategies",
		args,
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
func (e *jsiiProxy_Ec2Service) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		e,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Associates this service with a CloudMap service.
// Experimental.
func (e *jsiiProxy_Ec2Service) AssociateCloudMapService(options *AssociateCloudMapServiceOptions) {
	_jsii_.InvokeVoid(
		e,
		"associateCloudMapService",
		[]interface{}{options},
	)
}

// This method is called to attach this service to an Application Load Balancer.
//
// Don't call this function directly. Instead, call `listener.addTargets()`
// to add this service to a load balancer.
// Experimental.
func (e *jsiiProxy_Ec2Service) AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		e,
		"attachToApplicationTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// Registers the service as a target of a Classic Load Balancer (CLB).
//
// Don't call this. Call `loadBalancer.addTarget()` instead.
// Experimental.
func (e *jsiiProxy_Ec2Service) AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer) {
	_jsii_.InvokeVoid(
		e,
		"attachToClassicLB",
		[]interface{}{loadBalancer},
	)
}

// This method is called to attach this service to a Network Load Balancer.
//
// Don't call this function directly. Instead, call `listener.addTargets()`
// to add this service to a load balancer.
// Experimental.
func (e *jsiiProxy_Ec2Service) AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		e,
		"attachToNetworkTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// An attribute representing the minimum and maximum task count for an AutoScalingGroup.
// Experimental.
func (e *jsiiProxy_Ec2Service) AutoScaleTaskCount(props *awsapplicationautoscaling.EnableScalingProps) ScalableTaskCount {
	var returns ScalableTaskCount

	_jsii_.Invoke(
		e,
		"autoScaleTaskCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method is called to create a networkConfiguration.
// Deprecated: use configureAwsVpcNetworkingWithSecurityGroups instead.
func (e *jsiiProxy_Ec2Service) ConfigureAwsVpcNetworking(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroup awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		e,
		"configureAwsVpcNetworking",
		[]interface{}{vpc, assignPublicIp, vpcSubnets, securityGroup},
	)
}

// This method is called to create a networkConfiguration.
// Experimental.
func (e *jsiiProxy_Ec2Service) ConfigureAwsVpcNetworkingWithSecurityGroups(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroups *[]awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		e,
		"configureAwsVpcNetworkingWithSecurityGroups",
		[]interface{}{vpc, assignPublicIp, vpcSubnets, securityGroups},
	)
}

// Enable CloudMap service discovery for the service.
//
// Returns: The created CloudMap service
// Experimental.
func (e *jsiiProxy_Ec2Service) EnableCloudMap(options *CloudMapOptions) awsservicediscovery.Service {
	var returns awsservicediscovery.Service

	_jsii_.Invoke(
		e,
		"enableCloudMap",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Experimental.
func (e *jsiiProxy_Ec2Service) GeneratePhysicalName() *string {
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
func (e *jsiiProxy_Ec2Service) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (e *jsiiProxy_Ec2Service) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		e,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Return a load balancing target for a specific container and port.
//
// Use this function to create a load balancer target if you want to load balance to
// another container than the first essential container or the first mapped port on
// the container.
//
// Use the return value of this function where you would normally use a load balancer
// target, instead of the `Service` object itself.
//
// TODO: EXAMPLE
//
// Experimental.
func (e *jsiiProxy_Ec2Service) LoadBalancerTarget(options *LoadBalancerTargetOptions) IEcsLoadBalancerTarget {
	var returns IEcsLoadBalancerTarget

	_jsii_.Invoke(
		e,
		"loadBalancerTarget",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// This method returns the specified CloudWatch metric name for this service.
// Experimental.
func (e *jsiiProxy_Ec2Service) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		e,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this service's CPU utilization.
// Experimental.
func (e *jsiiProxy_Ec2Service) MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		e,
		"metricCpuUtilization",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this service's memory utilization.
// Experimental.
func (e *jsiiProxy_Ec2Service) MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		e,
		"metricMemoryUtilization",
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
func (e *jsiiProxy_Ec2Service) OnPrepare() {
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
func (e *jsiiProxy_Ec2Service) OnSynthesize(session constructs.ISynthesisSession) {
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
func (e *jsiiProxy_Ec2Service) OnValidate() *[]*string {
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
func (e *jsiiProxy_Ec2Service) Prepare() {
	_jsii_.InvokeVoid(
		e,
		"prepare",
		nil, // no parameters
	)
}

// Use this function to create all load balancer targets to be registered in this service, add them to target groups, and attach target groups to listeners accordingly.
//
// Alternatively, you can use `listener.addTargets()` to create targets and add them to target groups.
//
// TODO: EXAMPLE
//
// Experimental.
func (e *jsiiProxy_Ec2Service) RegisterLoadBalancerTargets(targets ...*EcsTarget) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		e,
		"registerLoadBalancerTargets",
		args,
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (e *jsiiProxy_Ec2Service) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		e,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (e *jsiiProxy_Ec2Service) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		e,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validates this Ec2Service.
// Experimental.
func (e *jsiiProxy_Ec2Service) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		e,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The properties to import from the service using the EC2 launch type.
// Experimental.
type Ec2ServiceAttributes struct {
	// The cluster that hosts the service.
	// Experimental.
	Cluster ICluster `json:"cluster"`
	// The service ARN.
	// Experimental.
	ServiceArn *string `json:"serviceArn"`
	// The name of the service.
	// Experimental.
	ServiceName *string `json:"serviceName"`
}

// The properties for defining a service using the EC2 launch type.
// Experimental.
type Ec2ServiceProps struct {
	// The name of the cluster that hosts the service.
	// Experimental.
	Cluster ICluster `json:"cluster"`
	// A list of Capacity Provider strategies used to place a service.
	// Experimental.
	CapacityProviderStrategies *[]*CapacityProviderStrategy `json:"capacityProviderStrategies"`
	// Whether to enable the deployment circuit breaker.
	//
	// If this property is defined, circuit breaker will be implicitly
	// enabled.
	// Experimental.
	CircuitBreaker *DeploymentCircuitBreaker `json:"circuitBreaker"`
	// The options for configuring an Amazon ECS service to use service discovery.
	// Experimental.
	CloudMapOptions *CloudMapOptions `json:"cloudMapOptions"`
	// Specifies which deployment controller to use for the service.
	//
	// For more information, see
	// [Amazon ECS Deployment Types](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-types.html)
	// Experimental.
	DeploymentController *DeploymentController `json:"deploymentController"`
	// The desired number of instantiations of the task definition to keep running on the service.
	// Experimental.
	DesiredCount *float64 `json:"desiredCount"`
	// Specifies whether to enable Amazon ECS managed tags for the tasks within the service.
	//
	// For more information, see
	// [Tagging Your Amazon ECS Resources](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-using-tags.html)
	// Experimental.
	EnableECSManagedTags *bool `json:"enableECSManagedTags"`
	// Whether to enable the ability to execute into a container.
	// Experimental.
	EnableExecuteCommand *bool `json:"enableExecuteCommand"`
	// The period of time, in seconds, that the Amazon ECS service scheduler ignores unhealthy Elastic Load Balancing target health checks after a task has first started.
	// Experimental.
	HealthCheckGracePeriod awscdk.Duration `json:"healthCheckGracePeriod"`
	// The maximum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that can run in a service during a deployment.
	// Experimental.
	MaxHealthyPercent *float64 `json:"maxHealthyPercent"`
	// The minimum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that must continue to run and remain healthy during a deployment.
	// Experimental.
	MinHealthyPercent *float64 `json:"minHealthyPercent"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Valid values are: PropagatedTagSource.SERVICE, PropagatedTagSource.TASK_DEFINITION or PropagatedTagSource.NONE
	// Experimental.
	PropagateTags PropagatedTagSource `json:"propagateTags"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Tags can only be propagated to the tasks within the service during service creation.
	// Deprecated: Use `propagateTags` instead.
	PropagateTaskTagsFrom PropagatedTagSource `json:"propagateTaskTagsFrom"`
	// The name of the service.
	// Experimental.
	ServiceName *string `json:"serviceName"`
	// The task definition to use for tasks in the service.
	//
	// [disable-awslint:ref-via-interface]
	// Experimental.
	TaskDefinition TaskDefinition `json:"taskDefinition"`
	// Specifies whether the task's elastic network interface receives a public IP address.
	//
	// If true, each task will receive a public IP address.
	//
	// This property is only used for tasks that use the awsvpc network mode.
	// Experimental.
	AssignPublicIp *bool `json:"assignPublicIp"`
	// Specifies whether the service will use the daemon scheduling strategy.
	//
	// If true, the service scheduler deploys exactly one task on each container instance in your cluster.
	//
	// When you are using this strategy, do not specify a desired number of tasks orany task placement strategies.
	// Experimental.
	Daemon *bool `json:"daemon"`
	// The placement constraints to use for tasks in the service.
	//
	// For more information, see
	// [Amazon ECS Task Placement Constraints](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-placement-constraints.html).
	// Experimental.
	PlacementConstraints *[]PlacementConstraint `json:"placementConstraints"`
	// The placement strategies to use for tasks in the service.
	//
	// For more information, see
	// [Amazon ECS Task Placement Strategies](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-placement-strategies.html).
	// Experimental.
	PlacementStrategies *[]PlacementStrategy `json:"placementStrategies"`
	// The security groups to associate with the service.
	//
	// If you do not specify a security group, the default security group for the VPC is used.
	//
	// This property is only used for tasks that use the awsvpc network mode.
	// Deprecated: use securityGroups instead.
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// The security groups to associate with the service.
	//
	// If you do not specify a security group, the default security group for the VPC is used.
	//
	// This property is only used for tasks that use the awsvpc network mode.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The subnets to associate with the service.
	//
	// This property is only used for tasks that use the awsvpc network mode.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
}

// The details of a task definition run on an EC2 cluster.
// Experimental.
type Ec2TaskDefinition interface {
	TaskDefinition
	IEc2TaskDefinition
	Compatibility() Compatibility
	Containers() *[]ContainerDefinition
	DefaultContainer() ContainerDefinition
	SetDefaultContainer(val ContainerDefinition)
	Env() *awscdk.ResourceEnvironment
	ExecutionRole() awsiam.IRole
	Family() *string
	InferenceAccelerators() *[]*InferenceAccelerator
	IsEc2Compatible() *bool
	IsExternalCompatible() *bool
	IsFargateCompatible() *bool
	NetworkMode() NetworkMode
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ReferencesSecretJsonField() *bool
	Stack() awscdk.Stack
	TaskDefinitionArn() *string
	TaskRole() awsiam.IRole
	AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition
	AddExtension(extension ITaskDefinitionExtension)
	AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter
	AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator)
	AddPlacementConstraint(constraint PlacementConstraint)
	AddToExecutionRolePolicy(statement awsiam.PolicyStatement)
	AddToTaskRolePolicy(statement awsiam.PolicyStatement)
	AddVolume(volume *Volume)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	ObtainExecutionRole() awsiam.IRole
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Ec2TaskDefinition
type jsiiProxy_Ec2TaskDefinition struct {
	jsiiProxy_TaskDefinition
	jsiiProxy_IEc2TaskDefinition
}

func (j *jsiiProxy_Ec2TaskDefinition) Compatibility() Compatibility {
	var returns Compatibility
	_jsii_.Get(
		j,
		"compatibility",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) Containers() *[]ContainerDefinition {
	var returns *[]ContainerDefinition
	_jsii_.Get(
		j,
		"containers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) DefaultContainer() ContainerDefinition {
	var returns ContainerDefinition
	_jsii_.Get(
		j,
		"defaultContainer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) ExecutionRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"executionRole",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) Family() *string {
	var returns *string
	_jsii_.Get(
		j,
		"family",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) InferenceAccelerators() *[]*InferenceAccelerator {
	var returns *[]*InferenceAccelerator
	_jsii_.Get(
		j,
		"inferenceAccelerators",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) IsEc2Compatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isEc2Compatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) IsExternalCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isExternalCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) IsFargateCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isFargateCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) NetworkMode() NetworkMode {
	var returns NetworkMode
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) ReferencesSecretJsonField() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"referencesSecretJsonField",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) TaskDefinitionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskDefinitionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Ec2TaskDefinition) TaskRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"taskRole",
		&returns,
	)
	return returns
}


// Constructs a new instance of the Ec2TaskDefinition class.
// Experimental.
func NewEc2TaskDefinition(scope constructs.Construct, id *string, props *Ec2TaskDefinitionProps) Ec2TaskDefinition {
	_init_.Initialize()

	j := jsiiProxy_Ec2TaskDefinition{}

	_jsii_.Create(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the Ec2TaskDefinition class.
// Experimental.
func NewEc2TaskDefinition_Override(e Ec2TaskDefinition, scope constructs.Construct, id *string, props *Ec2TaskDefinitionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		[]interface{}{scope, id, props},
		e,
	)
}

func (j *jsiiProxy_Ec2TaskDefinition) SetDefaultContainer(val ContainerDefinition) {
	_jsii_.Set(
		j,
		"defaultContainer",
		val,
	)
}

// Imports a task definition from the specified task definition ARN.
// Experimental.
func Ec2TaskDefinition_FromEc2TaskDefinitionArn(scope constructs.Construct, id *string, ec2TaskDefinitionArn *string) IEc2TaskDefinition {
	_init_.Initialize()

	var returns IEc2TaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		"fromEc2TaskDefinitionArn",
		[]interface{}{scope, id, ec2TaskDefinitionArn},
		&returns,
	)

	return returns
}

// Imports an existing Ec2 task definition from its attributes.
// Experimental.
func Ec2TaskDefinition_FromEc2TaskDefinitionAttributes(scope constructs.Construct, id *string, attrs *Ec2TaskDefinitionAttributes) IEc2TaskDefinition {
	_init_.Initialize()

	var returns IEc2TaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		"fromEc2TaskDefinitionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Imports a task definition from the specified task definition ARN.
//
// The task will have a compatibility of EC2+Fargate.
// Experimental.
func Ec2TaskDefinition_FromTaskDefinitionArn(scope constructs.Construct, id *string, taskDefinitionArn *string) ITaskDefinition {
	_init_.Initialize()

	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		"fromTaskDefinitionArn",
		[]interface{}{scope, id, taskDefinitionArn},
		&returns,
	)

	return returns
}

// Create a task definition from a task definition reference.
// Experimental.
func Ec2TaskDefinition_FromTaskDefinitionAttributes(scope constructs.Construct, id *string, attrs *TaskDefinitionAttributes) ITaskDefinition {
	_init_.Initialize()

	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		"fromTaskDefinitionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Ec2TaskDefinition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Ec2TaskDefinition_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Ec2TaskDefinition",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a new container to the task definition.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition {
	var returns ContainerDefinition

	_jsii_.Invoke(
		e,
		"addContainer",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds the specified extension to the task definition.
//
// Extension can be used to apply a packaged modification to
// a task definition.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddExtension(extension ITaskDefinitionExtension) {
	_jsii_.InvokeVoid(
		e,
		"addExtension",
		[]interface{}{extension},
	)
}

// Adds a firelens log router to the task definition.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter {
	var returns FirelensLogRouter

	_jsii_.Invoke(
		e,
		"addFirelensLogRouter",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds an inference accelerator to the task definition.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator) {
	_jsii_.InvokeVoid(
		e,
		"addInferenceAccelerator",
		[]interface{}{inferenceAccelerator},
	)
}

// Adds the specified placement constraint to the task definition.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddPlacementConstraint(constraint PlacementConstraint) {
	_jsii_.InvokeVoid(
		e,
		"addPlacementConstraint",
		[]interface{}{constraint},
	)
}

// Adds a policy statement to the task execution IAM role.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddToExecutionRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		e,
		"addToExecutionRolePolicy",
		[]interface{}{statement},
	)
}

// Adds a policy statement to the task IAM role.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddToTaskRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		e,
		"addToTaskRolePolicy",
		[]interface{}{statement},
	)
}

// Adds a volume to the task definition.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) AddVolume(volume *Volume) {
	_jsii_.InvokeVoid(
		e,
		"addVolume",
		[]interface{}{volume},
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
func (e *jsiiProxy_Ec2TaskDefinition) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		e,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) GeneratePhysicalName() *string {
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
func (e *jsiiProxy_Ec2TaskDefinition) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (e *jsiiProxy_Ec2TaskDefinition) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		e,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Creates the task execution IAM role if it doesn't already exist.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) ObtainExecutionRole() awsiam.IRole {
	var returns awsiam.IRole

	_jsii_.Invoke(
		e,
		"obtainExecutionRole",
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
func (e *jsiiProxy_Ec2TaskDefinition) OnPrepare() {
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
func (e *jsiiProxy_Ec2TaskDefinition) OnSynthesize(session constructs.ISynthesisSession) {
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
func (e *jsiiProxy_Ec2TaskDefinition) OnValidate() *[]*string {
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
func (e *jsiiProxy_Ec2TaskDefinition) Prepare() {
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
func (e *jsiiProxy_Ec2TaskDefinition) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		e,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		e,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validates the task definition.
// Experimental.
func (e *jsiiProxy_Ec2TaskDefinition) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		e,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Attributes used to import an existing EC2 task definition.
// Experimental.
type Ec2TaskDefinitionAttributes struct {
	// The arn of the task definition.
	// Experimental.
	TaskDefinitionArn *string `json:"taskDefinitionArn"`
	// The networking mode to use for the containers in the task.
	// Experimental.
	NetworkMode NetworkMode `json:"networkMode"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
}

// The properties for a task definition run on an EC2 cluster.
// Experimental.
type Ec2TaskDefinitionProps struct {
	// The name of the IAM task execution role that grants the ECS agent to call AWS APIs on your behalf.
	//
	// The role will be used to retrieve container images from ECR and create CloudWatch log groups.
	// Experimental.
	ExecutionRole awsiam.IRole `json:"executionRole"`
	// The name of a family that this task definition is registered to.
	//
	// A family groups multiple versions of a task definition.
	// Experimental.
	Family *string `json:"family"`
	// The configuration details for the App Mesh proxy.
	// Experimental.
	ProxyConfiguration ProxyConfiguration `json:"proxyConfiguration"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
	// The list of volume definitions for the task.
	//
	// For more information, see
	// [Task Definition Parameter Volumes](https://docs.aws.amazon.com/AmazonECS/latest/developerguide//task_definition_parameters.html#volumes).
	// Experimental.
	Volumes *[]*Volume `json:"volumes"`
	// The inference accelerators to use for the containers in the task.
	//
	// Not supported in Fargate.
	// Experimental.
	InferenceAccelerators *[]*InferenceAccelerator `json:"inferenceAccelerators"`
	// The IPC resource namespace to use for the containers in the task.
	//
	// Not supported in Fargate and Windows containers.
	// Experimental.
	IpcMode IpcMode `json:"ipcMode"`
	// The Docker networking mode to use for the containers in the task.
	//
	// The valid values are none, bridge, awsvpc, and host.
	// Experimental.
	NetworkMode NetworkMode `json:"networkMode"`
	// The process namespace to use for the containers in the task.
	//
	// Not supported in Fargate and Windows containers.
	// Experimental.
	PidMode PidMode `json:"pidMode"`
	// An array of placement constraint objects to use for the task.
	//
	// You can
	// specify a maximum of 10 constraints per task (this limit includes
	// constraints in the task definition and those specified at run time).
	// Experimental.
	PlacementConstraints *[]PlacementConstraint `json:"placementConstraints"`
}

// An image from an Amazon ECR repository.
// Experimental.
type EcrImage interface {
	ContainerImage
	ImageName() *string
	Bind(_scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig
}

// The jsii proxy struct for EcrImage
type jsiiProxy_EcrImage struct {
	jsiiProxy_ContainerImage
}

func (j *jsiiProxy_EcrImage) ImageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageName",
		&returns,
	)
	return returns
}


// Constructs a new instance of the EcrImage class.
// Experimental.
func NewEcrImage(repository awsecr.IRepository, tagOrDigest *string) EcrImage {
	_init_.Initialize()

	j := jsiiProxy_EcrImage{}

	_jsii_.Create(
		"monocdk.aws_ecs.EcrImage",
		[]interface{}{repository, tagOrDigest},
		&j,
	)

	return &j
}

// Constructs a new instance of the EcrImage class.
// Experimental.
func NewEcrImage_Override(e EcrImage, repository awsecr.IRepository, tagOrDigest *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.EcrImage",
		[]interface{}{repository, tagOrDigest},
		e,
	)
}

// Reference an image that's constructed directly from sources on disk.
//
// If you already have a `DockerImageAsset` instance, you can use the
// `ContainerImage.fromDockerImageAsset` method instead.
// Experimental.
func EcrImage_FromAsset(directory *string, props *AssetImageProps) AssetImage {
	_init_.Initialize()

	var returns AssetImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EcrImage",
		"fromAsset",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Use an existing `DockerImageAsset` for this container image.
// Experimental.
func EcrImage_FromDockerImageAsset(asset awsecrassets.DockerImageAsset) ContainerImage {
	_init_.Initialize()

	var returns ContainerImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EcrImage",
		"fromDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Reference an image in an ECR repository.
// Experimental.
func EcrImage_FromEcrRepository(repository awsecr.IRepository, tag *string) EcrImage {
	_init_.Initialize()

	var returns EcrImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EcrImage",
		"fromEcrRepository",
		[]interface{}{repository, tag},
		&returns,
	)

	return returns
}

// Reference an image on DockerHub or another online registry.
// Experimental.
func EcrImage_FromRegistry(name *string, props *RepositoryImageProps) RepositoryImage {
	_init_.Initialize()

	var returns RepositoryImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EcrImage",
		"fromRegistry",
		[]interface{}{name, props},
		&returns,
	)

	return returns
}

// Called when the image is used by a ContainerDefinition.
// Experimental.
func (e *jsiiProxy_EcrImage) Bind(_scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig {
	var returns *ContainerImageConfig

	_jsii_.Invoke(
		e,
		"bind",
		[]interface{}{_scope, containerDefinition},
		&returns,
	)

	return returns
}

// Construct a Linux or Windows machine image from the latest ECS Optimized AMI published in SSM.
// Deprecated: see {@link EcsOptimizedImage#amazonLinux}, {@link EcsOptimizedImage#amazonLinux} and {@link EcsOptimizedImage#windows}
type EcsOptimizedAmi interface {
	awsec2.IMachineImage
	GetImage(scope awscdk.Construct) *awsec2.MachineImageConfig
}

// The jsii proxy struct for EcsOptimizedAmi
type jsiiProxy_EcsOptimizedAmi struct {
	internal.Type__awsec2IMachineImage
}

// Constructs a new instance of the EcsOptimizedAmi class.
// Deprecated: see {@link EcsOptimizedImage#amazonLinux}, {@link EcsOptimizedImage#amazonLinux} and {@link EcsOptimizedImage#windows}
func NewEcsOptimizedAmi(props *EcsOptimizedAmiProps) EcsOptimizedAmi {
	_init_.Initialize()

	j := jsiiProxy_EcsOptimizedAmi{}

	_jsii_.Create(
		"monocdk.aws_ecs.EcsOptimizedAmi",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the EcsOptimizedAmi class.
// Deprecated: see {@link EcsOptimizedImage#amazonLinux}, {@link EcsOptimizedImage#amazonLinux} and {@link EcsOptimizedImage#windows}
func NewEcsOptimizedAmi_Override(e EcsOptimizedAmi, props *EcsOptimizedAmiProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.EcsOptimizedAmi",
		[]interface{}{props},
		e,
	)
}

// Return the correct image.
// Deprecated: see {@link EcsOptimizedImage#amazonLinux}, {@link EcsOptimizedImage#amazonLinux} and {@link EcsOptimizedImage#windows}
func (e *jsiiProxy_EcsOptimizedAmi) GetImage(scope awscdk.Construct) *awsec2.MachineImageConfig {
	var returns *awsec2.MachineImageConfig

	_jsii_.Invoke(
		e,
		"getImage",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// The properties that define which ECS-optimized AMI is used.
// Deprecated: see {@link EcsOptimizedImage}
type EcsOptimizedAmiProps struct {
	// The Amazon Linux generation to use.
	// Deprecated: see {@link EcsOptimizedImage}
	Generation awsec2.AmazonLinuxGeneration `json:"generation"`
	// The ECS-optimized AMI variant to use.
	// Deprecated: see {@link EcsOptimizedImage}
	HardwareType AmiHardwareType `json:"hardwareType"`
	// The Windows Server version to use.
	// Deprecated: see {@link EcsOptimizedImage}
	WindowsVersion WindowsOptimizedVersion `json:"windowsVersion"`
}

// Construct a Linux or Windows machine image from the latest ECS Optimized AMI published in SSM.
// Experimental.
type EcsOptimizedImage interface {
	awsec2.IMachineImage
	GetImage(scope awscdk.Construct) *awsec2.MachineImageConfig
}

// The jsii proxy struct for EcsOptimizedImage
type jsiiProxy_EcsOptimizedImage struct {
	internal.Type__awsec2IMachineImage
}

// Construct an Amazon Linux AMI image from the latest ECS Optimized AMI published in SSM.
// Experimental.
func EcsOptimizedImage_AmazonLinux() EcsOptimizedImage {
	_init_.Initialize()

	var returns EcsOptimizedImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EcsOptimizedImage",
		"amazonLinux",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construct an Amazon Linux 2 image from the latest ECS Optimized AMI published in SSM.
// Experimental.
func EcsOptimizedImage_AmazonLinux2(hardwareType AmiHardwareType) EcsOptimizedImage {
	_init_.Initialize()

	var returns EcsOptimizedImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EcsOptimizedImage",
		"amazonLinux2",
		[]interface{}{hardwareType},
		&returns,
	)

	return returns
}

// Construct a Windows image from the latest ECS Optimized AMI published in SSM.
// Experimental.
func EcsOptimizedImage_Windows(windowsVersion WindowsOptimizedVersion) EcsOptimizedImage {
	_init_.Initialize()

	var returns EcsOptimizedImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EcsOptimizedImage",
		"windows",
		[]interface{}{windowsVersion},
		&returns,
	)

	return returns
}

// Return the correct image.
// Experimental.
func (e *jsiiProxy_EcsOptimizedImage) GetImage(scope awscdk.Construct) *awsec2.MachineImageConfig {
	var returns *awsec2.MachineImageConfig

	_jsii_.Invoke(
		e,
		"getImage",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Experimental.
type EcsTarget struct {
	// The name of the container.
	// Experimental.
	ContainerName *string `json:"containerName"`
	// Listener and properties for adding target group to the listener.
	// Experimental.
	Listener ListenerConfig `json:"listener"`
	// ID for a target group to be created.
	// Experimental.
	NewTargetGroupId *string `json:"newTargetGroupId"`
	// The port number of the container.
	//
	// Only applicable when using application/network load balancers.
	// Experimental.
	ContainerPort *float64 `json:"containerPort"`
	// The protocol used for the port mapping.
	//
	// Only applicable when using application load balancers.
	// Experimental.
	Protocol Protocol `json:"protocol"`
}

// The configuration for an Elastic FileSystem volume.
// Experimental.
type EfsVolumeConfiguration struct {
	// The Amazon EFS file system ID to use.
	// Experimental.
	FileSystemId *string `json:"fileSystemId"`
	// The authorization configuration details for the Amazon EFS file system.
	// Experimental.
	AuthorizationConfig *AuthorizationConfig `json:"authorizationConfig"`
	// The directory within the Amazon EFS file system to mount as the root directory inside the host.
	//
	// Specifying / will have the same effect as omitting this parameter.
	// Experimental.
	RootDirectory *string `json:"rootDirectory"`
	// Whether or not to enable encryption for Amazon EFS data in transit between the Amazon ECS host and the Amazon EFS server.
	//
	// Transit encryption must be enabled if Amazon EFS IAM authorization is used.
	//
	// Valid values: ENABLED | DISABLED
	// Experimental.
	TransitEncryption *string `json:"transitEncryption"`
	// The port to use when sending encrypted data between the Amazon ECS host and the Amazon EFS server.
	//
	// EFS mount helper uses.
	// Experimental.
	TransitEncryptionPort *float64 `json:"transitEncryptionPort"`
}

// Constructs for types of environment files.
// Experimental.
type EnvironmentFile interface {
	Bind(scope awscdk.Construct) *EnvironmentFileConfig
}

// The jsii proxy struct for EnvironmentFile
type jsiiProxy_EnvironmentFile struct {
	_ byte // padding
}

// Experimental.
func NewEnvironmentFile_Override(e EnvironmentFile) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.EnvironmentFile",
		nil, // no parameters
		e,
	)
}

// Loads the environment file from a local disk path.
// Experimental.
func EnvironmentFile_FromAsset(path *string, options *awss3assets.AssetOptions) AssetEnvironmentFile {
	_init_.Initialize()

	var returns AssetEnvironmentFile

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EnvironmentFile",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Loads the environment file from an S3 bucket.
//
// Returns: `S3EnvironmentFile` associated with the specified S3 object.
// Experimental.
func EnvironmentFile_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3EnvironmentFile {
	_init_.Initialize()

	var returns S3EnvironmentFile

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.EnvironmentFile",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Called when the container is initialized to allow this object to bind to the stack.
// Experimental.
func (e *jsiiProxy_EnvironmentFile) Bind(scope awscdk.Construct) *EnvironmentFileConfig {
	var returns *EnvironmentFileConfig

	_jsii_.Invoke(
		e,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Configuration for the environment file.
// Experimental.
type EnvironmentFileConfig struct {
	// The type of environment file.
	// Experimental.
	FileType EnvironmentFileType `json:"fileType"`
	// The location of the environment file in S3.
	// Experimental.
	S3Location *awss3.Location `json:"s3Location"`
}

// Type of environment file to be included in the container definition.
// Experimental.
type EnvironmentFileType string

const (
	EnvironmentFileType_S3 EnvironmentFileType = "S3"
)

// The details of the execute command configuration.
//
// For more information, see
// [ExecuteCommandConfiguration] https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-executecommandconfiguration.html
// Experimental.
type ExecuteCommandConfiguration struct {
	// The AWS Key Management Service key ID to encrypt the data between the local client and the container.
	// Experimental.
	KmsKey awskms.IKey `json:"kmsKey"`
	// The log configuration for the results of the execute command actions.
	//
	// The logs can be sent to CloudWatch Logs or an Amazon S3 bucket.
	// Experimental.
	LogConfiguration *ExecuteCommandLogConfiguration `json:"logConfiguration"`
	// The log settings to use for logging the execute command session.
	// Experimental.
	Logging ExecuteCommandLogging `json:"logging"`
}

// The log configuration for the results of the execute command actions.
//
// The logs can be sent to CloudWatch Logs and/ or an Amazon S3 bucket.
// For more information, see [ExecuteCommandLogConfiguration] https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-executecommandlogconfiguration.html
// Experimental.
type ExecuteCommandLogConfiguration struct {
	// Whether or not to enable encryption on the CloudWatch logs.
	// Experimental.
	CloudWatchEncryptionEnabled *bool `json:"cloudWatchEncryptionEnabled"`
	// The name of the CloudWatch log group to send logs to.
	//
	// The CloudWatch log group must already be created.
	// Experimental.
	CloudWatchLogGroup awslogs.ILogGroup `json:"cloudWatchLogGroup"`
	// The name of the S3 bucket to send logs to.
	//
	// The S3 bucket must already be created.
	// Experimental.
	S3Bucket awss3.IBucket `json:"s3Bucket"`
	// Whether or not to enable encryption on the CloudWatch logs.
	// Experimental.
	S3EncryptionEnabled *bool `json:"s3EncryptionEnabled"`
	// An optional folder in the S3 bucket to place logs in.
	// Experimental.
	S3KeyPrefix *string `json:"s3KeyPrefix"`
}

// The log settings to use to for logging the execute command session.
//
// For more information, see
// [Logging] https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-executecommandconfiguration.html#cfn-ecs-cluster-executecommandconfiguration-logging
// Experimental.
type ExecuteCommandLogging string

const (
	ExecuteCommandLogging_NONE ExecuteCommandLogging = "NONE"
	ExecuteCommandLogging_DEFAULT ExecuteCommandLogging = "DEFAULT"
	ExecuteCommandLogging_OVERRIDE ExecuteCommandLogging = "OVERRIDE"
)

// The platform version on which to run your service.
// See: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html
//
// Experimental.
type FargatePlatformVersion string

const (
	FargatePlatformVersion_LATEST FargatePlatformVersion = "LATEST"
	FargatePlatformVersion_VERSION1_4 FargatePlatformVersion = "VERSION1_4"
	FargatePlatformVersion_VERSION1_3 FargatePlatformVersion = "VERSION1_3"
	FargatePlatformVersion_VERSION1_2 FargatePlatformVersion = "VERSION1_2"
	FargatePlatformVersion_VERSION1_1 FargatePlatformVersion = "VERSION1_1"
	FargatePlatformVersion_VERSION1_0 FargatePlatformVersion = "VERSION1_0"
)

// This creates a service using the Fargate launch type on an ECS cluster.
// Experimental.
type FargateService interface {
	BaseService
	IFargateService
	CloudmapService() awsservicediscovery.Service
	SetCloudmapService(val awsservicediscovery.Service)
	CloudMapService() awsservicediscovery.IService
	Cluster() ICluster
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	LoadBalancers() *[]*CfnService_LoadBalancerProperty
	SetLoadBalancers(val *[]*CfnService_LoadBalancerProperty)
	NetworkConfiguration() *CfnService_NetworkConfigurationProperty
	SetNetworkConfiguration(val *CfnService_NetworkConfigurationProperty)
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ServiceArn() *string
	ServiceName() *string
	ServiceRegistries() *[]*CfnService_ServiceRegistryProperty
	SetServiceRegistries(val *[]*CfnService_ServiceRegistryProperty)
	Stack() awscdk.Stack
	TaskDefinition() TaskDefinition
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AssociateCloudMapService(options *AssociateCloudMapServiceOptions)
	AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer)
	AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps
	AutoScaleTaskCount(props *awsapplicationautoscaling.EnableScalingProps) ScalableTaskCount
	ConfigureAwsVpcNetworking(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroup awsec2.ISecurityGroup)
	ConfigureAwsVpcNetworkingWithSecurityGroups(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroups *[]awsec2.ISecurityGroup)
	EnableCloudMap(options *CloudMapOptions) awsservicediscovery.Service
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LoadBalancerTarget(options *LoadBalancerTargetOptions) IEcsLoadBalancerTarget
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RegisterLoadBalancerTargets(targets ...*EcsTarget)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for FargateService
type jsiiProxy_FargateService struct {
	jsiiProxy_BaseService
	jsiiProxy_IFargateService
}

func (j *jsiiProxy_FargateService) CloudmapService() awsservicediscovery.Service {
	var returns awsservicediscovery.Service
	_jsii_.Get(
		j,
		"cloudmapService",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) CloudMapService() awsservicediscovery.IService {
	var returns awsservicediscovery.IService
	_jsii_.Get(
		j,
		"cloudMapService",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) Cluster() ICluster {
	var returns ICluster
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) LoadBalancers() *[]*CfnService_LoadBalancerProperty {
	var returns *[]*CfnService_LoadBalancerProperty
	_jsii_.Get(
		j,
		"loadBalancers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) NetworkConfiguration() *CfnService_NetworkConfigurationProperty {
	var returns *CfnService_NetworkConfigurationProperty
	_jsii_.Get(
		j,
		"networkConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) ServiceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) ServiceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) ServiceRegistries() *[]*CfnService_ServiceRegistryProperty {
	var returns *[]*CfnService_ServiceRegistryProperty
	_jsii_.Get(
		j,
		"serviceRegistries",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateService) TaskDefinition() TaskDefinition {
	var returns TaskDefinition
	_jsii_.Get(
		j,
		"taskDefinition",
		&returns,
	)
	return returns
}


// Constructs a new instance of the FargateService class.
// Experimental.
func NewFargateService(scope constructs.Construct, id *string, props *FargateServiceProps) FargateService {
	_init_.Initialize()

	j := jsiiProxy_FargateService{}

	_jsii_.Create(
		"monocdk.aws_ecs.FargateService",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the FargateService class.
// Experimental.
func NewFargateService_Override(f FargateService, scope constructs.Construct, id *string, props *FargateServiceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.FargateService",
		[]interface{}{scope, id, props},
		f,
	)
}

func (j *jsiiProxy_FargateService) SetCloudmapService(val awsservicediscovery.Service) {
	_jsii_.Set(
		j,
		"cloudmapService",
		val,
	)
}

func (j *jsiiProxy_FargateService) SetLoadBalancers(val *[]*CfnService_LoadBalancerProperty) {
	_jsii_.Set(
		j,
		"loadBalancers",
		val,
	)
}

func (j *jsiiProxy_FargateService) SetNetworkConfiguration(val *CfnService_NetworkConfigurationProperty) {
	_jsii_.Set(
		j,
		"networkConfiguration",
		val,
	)
}

func (j *jsiiProxy_FargateService) SetServiceRegistries(val *[]*CfnService_ServiceRegistryProperty) {
	_jsii_.Set(
		j,
		"serviceRegistries",
		val,
	)
}

// Imports from the specified service ARN.
// Experimental.
func FargateService_FromFargateServiceArn(scope constructs.Construct, id *string, fargateServiceArn *string) IFargateService {
	_init_.Initialize()

	var returns IFargateService

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateService",
		"fromFargateServiceArn",
		[]interface{}{scope, id, fargateServiceArn},
		&returns,
	)

	return returns
}

// Imports from the specified service attrributes.
// Experimental.
func FargateService_FromFargateServiceAttributes(scope constructs.Construct, id *string, attrs *FargateServiceAttributes) IBaseService {
	_init_.Initialize()

	var returns IBaseService

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateService",
		"fromFargateServiceAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func FargateService_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateService",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func FargateService_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateService",
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
func (f *jsiiProxy_FargateService) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		f,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Associates this service with a CloudMap service.
// Experimental.
func (f *jsiiProxy_FargateService) AssociateCloudMapService(options *AssociateCloudMapServiceOptions) {
	_jsii_.InvokeVoid(
		f,
		"associateCloudMapService",
		[]interface{}{options},
	)
}

// This method is called to attach this service to an Application Load Balancer.
//
// Don't call this function directly. Instead, call `listener.addTargets()`
// to add this service to a load balancer.
// Experimental.
func (f *jsiiProxy_FargateService) AttachToApplicationTargetGroup(targetGroup awselasticloadbalancingv2.IApplicationTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		f,
		"attachToApplicationTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// Registers the service as a target of a Classic Load Balancer (CLB).
//
// Don't call this. Call `loadBalancer.addTarget()` instead.
// Experimental.
func (f *jsiiProxy_FargateService) AttachToClassicLB(loadBalancer awselasticloadbalancing.LoadBalancer) {
	_jsii_.InvokeVoid(
		f,
		"attachToClassicLB",
		[]interface{}{loadBalancer},
	)
}

// This method is called to attach this service to a Network Load Balancer.
//
// Don't call this function directly. Instead, call `listener.addTargets()`
// to add this service to a load balancer.
// Experimental.
func (f *jsiiProxy_FargateService) AttachToNetworkTargetGroup(targetGroup awselasticloadbalancingv2.INetworkTargetGroup) *awselasticloadbalancingv2.LoadBalancerTargetProps {
	var returns *awselasticloadbalancingv2.LoadBalancerTargetProps

	_jsii_.Invoke(
		f,
		"attachToNetworkTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// An attribute representing the minimum and maximum task count for an AutoScalingGroup.
// Experimental.
func (f *jsiiProxy_FargateService) AutoScaleTaskCount(props *awsapplicationautoscaling.EnableScalingProps) ScalableTaskCount {
	var returns ScalableTaskCount

	_jsii_.Invoke(
		f,
		"autoScaleTaskCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method is called to create a networkConfiguration.
// Deprecated: use configureAwsVpcNetworkingWithSecurityGroups instead.
func (f *jsiiProxy_FargateService) ConfigureAwsVpcNetworking(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroup awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		f,
		"configureAwsVpcNetworking",
		[]interface{}{vpc, assignPublicIp, vpcSubnets, securityGroup},
	)
}

// This method is called to create a networkConfiguration.
// Experimental.
func (f *jsiiProxy_FargateService) ConfigureAwsVpcNetworkingWithSecurityGroups(vpc awsec2.IVpc, assignPublicIp *bool, vpcSubnets *awsec2.SubnetSelection, securityGroups *[]awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		f,
		"configureAwsVpcNetworkingWithSecurityGroups",
		[]interface{}{vpc, assignPublicIp, vpcSubnets, securityGroups},
	)
}

// Enable CloudMap service discovery for the service.
//
// Returns: The created CloudMap service
// Experimental.
func (f *jsiiProxy_FargateService) EnableCloudMap(options *CloudMapOptions) awsservicediscovery.Service {
	var returns awsservicediscovery.Service

	_jsii_.Invoke(
		f,
		"enableCloudMap",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Experimental.
func (f *jsiiProxy_FargateService) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FargateService) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FargateService) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Return a load balancing target for a specific container and port.
//
// Use this function to create a load balancer target if you want to load balance to
// another container than the first essential container or the first mapped port on
// the container.
//
// Use the return value of this function where you would normally use a load balancer
// target, instead of the `Service` object itself.
//
// TODO: EXAMPLE
//
// Experimental.
func (f *jsiiProxy_FargateService) LoadBalancerTarget(options *LoadBalancerTargetOptions) IEcsLoadBalancerTarget {
	var returns IEcsLoadBalancerTarget

	_jsii_.Invoke(
		f,
		"loadBalancerTarget",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// This method returns the specified CloudWatch metric name for this service.
// Experimental.
func (f *jsiiProxy_FargateService) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this service's CPU utilization.
// Experimental.
func (f *jsiiProxy_FargateService) MetricCpuUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricCpuUtilization",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// This method returns the CloudWatch metric for this service's memory utilization.
// Experimental.
func (f *jsiiProxy_FargateService) MetricMemoryUtilization(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricMemoryUtilization",
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
func (f *jsiiProxy_FargateService) OnPrepare() {
	_jsii_.InvokeVoid(
		f,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (f *jsiiProxy_FargateService) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
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
func (f *jsiiProxy_FargateService) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FargateService) Prepare() {
	_jsii_.InvokeVoid(
		f,
		"prepare",
		nil, // no parameters
	)
}

// Use this function to create all load balancer targets to be registered in this service, add them to target groups, and attach target groups to listeners accordingly.
//
// Alternatively, you can use `listener.addTargets()` to create targets and add them to target groups.
//
// TODO: EXAMPLE
//
// Experimental.
func (f *jsiiProxy_FargateService) RegisterLoadBalancerTargets(targets ...*EcsTarget) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"registerLoadBalancerTargets",
		args,
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (f *jsiiProxy_FargateService) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (f *jsiiProxy_FargateService) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FargateService) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The properties to import from the service using the Fargate launch type.
// Experimental.
type FargateServiceAttributes struct {
	// The cluster that hosts the service.
	// Experimental.
	Cluster ICluster `json:"cluster"`
	// The service ARN.
	// Experimental.
	ServiceArn *string `json:"serviceArn"`
	// The name of the service.
	// Experimental.
	ServiceName *string `json:"serviceName"`
}

// The properties for defining a service using the Fargate launch type.
// Experimental.
type FargateServiceProps struct {
	// The name of the cluster that hosts the service.
	// Experimental.
	Cluster ICluster `json:"cluster"`
	// A list of Capacity Provider strategies used to place a service.
	// Experimental.
	CapacityProviderStrategies *[]*CapacityProviderStrategy `json:"capacityProviderStrategies"`
	// Whether to enable the deployment circuit breaker.
	//
	// If this property is defined, circuit breaker will be implicitly
	// enabled.
	// Experimental.
	CircuitBreaker *DeploymentCircuitBreaker `json:"circuitBreaker"`
	// The options for configuring an Amazon ECS service to use service discovery.
	// Experimental.
	CloudMapOptions *CloudMapOptions `json:"cloudMapOptions"`
	// Specifies which deployment controller to use for the service.
	//
	// For more information, see
	// [Amazon ECS Deployment Types](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-types.html)
	// Experimental.
	DeploymentController *DeploymentController `json:"deploymentController"`
	// The desired number of instantiations of the task definition to keep running on the service.
	// Experimental.
	DesiredCount *float64 `json:"desiredCount"`
	// Specifies whether to enable Amazon ECS managed tags for the tasks within the service.
	//
	// For more information, see
	// [Tagging Your Amazon ECS Resources](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-using-tags.html)
	// Experimental.
	EnableECSManagedTags *bool `json:"enableECSManagedTags"`
	// Whether to enable the ability to execute into a container.
	// Experimental.
	EnableExecuteCommand *bool `json:"enableExecuteCommand"`
	// The period of time, in seconds, that the Amazon ECS service scheduler ignores unhealthy Elastic Load Balancing target health checks after a task has first started.
	// Experimental.
	HealthCheckGracePeriod awscdk.Duration `json:"healthCheckGracePeriod"`
	// The maximum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that can run in a service during a deployment.
	// Experimental.
	MaxHealthyPercent *float64 `json:"maxHealthyPercent"`
	// The minimum number of tasks, specified as a percentage of the Amazon ECS service's DesiredCount value, that must continue to run and remain healthy during a deployment.
	// Experimental.
	MinHealthyPercent *float64 `json:"minHealthyPercent"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Valid values are: PropagatedTagSource.SERVICE, PropagatedTagSource.TASK_DEFINITION or PropagatedTagSource.NONE
	// Experimental.
	PropagateTags PropagatedTagSource `json:"propagateTags"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks in the service.
	//
	// Tags can only be propagated to the tasks within the service during service creation.
	// Deprecated: Use `propagateTags` instead.
	PropagateTaskTagsFrom PropagatedTagSource `json:"propagateTaskTagsFrom"`
	// The name of the service.
	// Experimental.
	ServiceName *string `json:"serviceName"`
	// The task definition to use for tasks in the service.
	//
	// [disable-awslint:ref-via-interface]
	// Experimental.
	TaskDefinition TaskDefinition `json:"taskDefinition"`
	// Specifies whether the task's elastic network interface receives a public IP address.
	//
	// If true, each task will receive a public IP address.
	// Experimental.
	AssignPublicIp *bool `json:"assignPublicIp"`
	// The platform version on which to run your service.
	//
	// If one is not specified, the LATEST platform version is used by default. For more information, see
	// [AWS Fargate Platform Versions](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html)
	// in the Amazon Elastic Container Service Developer Guide.
	// Experimental.
	PlatformVersion FargatePlatformVersion `json:"platformVersion"`
	// The security groups to associate with the service.
	//
	// If you do not specify a security group, the default security group for the VPC is used.
	// Deprecated: use securityGroups instead.
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// The security groups to associate with the service.
	//
	// If you do not specify a security group, the default security group for the VPC is used.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The subnets to associate with the service.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
}

// The details of a task definition run on a Fargate cluster.
// Experimental.
type FargateTaskDefinition interface {
	TaskDefinition
	IFargateTaskDefinition
	Compatibility() Compatibility
	Containers() *[]ContainerDefinition
	DefaultContainer() ContainerDefinition
	SetDefaultContainer(val ContainerDefinition)
	Env() *awscdk.ResourceEnvironment
	ExecutionRole() awsiam.IRole
	Family() *string
	InferenceAccelerators() *[]*InferenceAccelerator
	IsEc2Compatible() *bool
	IsExternalCompatible() *bool
	IsFargateCompatible() *bool
	NetworkMode() NetworkMode
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ReferencesSecretJsonField() *bool
	Stack() awscdk.Stack
	TaskDefinitionArn() *string
	TaskRole() awsiam.IRole
	AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition
	AddExtension(extension ITaskDefinitionExtension)
	AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter
	AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator)
	AddPlacementConstraint(constraint PlacementConstraint)
	AddToExecutionRolePolicy(statement awsiam.PolicyStatement)
	AddToTaskRolePolicy(statement awsiam.PolicyStatement)
	AddVolume(volume *Volume)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	ObtainExecutionRole() awsiam.IRole
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for FargateTaskDefinition
type jsiiProxy_FargateTaskDefinition struct {
	jsiiProxy_TaskDefinition
	jsiiProxy_IFargateTaskDefinition
}

func (j *jsiiProxy_FargateTaskDefinition) Compatibility() Compatibility {
	var returns Compatibility
	_jsii_.Get(
		j,
		"compatibility",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Containers() *[]ContainerDefinition {
	var returns *[]ContainerDefinition
	_jsii_.Get(
		j,
		"containers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) DefaultContainer() ContainerDefinition {
	var returns ContainerDefinition
	_jsii_.Get(
		j,
		"defaultContainer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) ExecutionRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"executionRole",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Family() *string {
	var returns *string
	_jsii_.Get(
		j,
		"family",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) InferenceAccelerators() *[]*InferenceAccelerator {
	var returns *[]*InferenceAccelerator
	_jsii_.Get(
		j,
		"inferenceAccelerators",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) IsEc2Compatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isEc2Compatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) IsExternalCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isExternalCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) IsFargateCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isFargateCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) NetworkMode() NetworkMode {
	var returns NetworkMode
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) ReferencesSecretJsonField() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"referencesSecretJsonField",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) TaskDefinitionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskDefinitionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FargateTaskDefinition) TaskRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"taskRole",
		&returns,
	)
	return returns
}


// Constructs a new instance of the FargateTaskDefinition class.
// Experimental.
func NewFargateTaskDefinition(scope constructs.Construct, id *string, props *FargateTaskDefinitionProps) FargateTaskDefinition {
	_init_.Initialize()

	j := jsiiProxy_FargateTaskDefinition{}

	_jsii_.Create(
		"monocdk.aws_ecs.FargateTaskDefinition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the FargateTaskDefinition class.
// Experimental.
func NewFargateTaskDefinition_Override(f FargateTaskDefinition, scope constructs.Construct, id *string, props *FargateTaskDefinitionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.FargateTaskDefinition",
		[]interface{}{scope, id, props},
		f,
	)
}

func (j *jsiiProxy_FargateTaskDefinition) SetDefaultContainer(val ContainerDefinition) {
	_jsii_.Set(
		j,
		"defaultContainer",
		val,
	)
}

// Imports a task definition from the specified task definition ARN.
// Experimental.
func FargateTaskDefinition_FromFargateTaskDefinitionArn(scope constructs.Construct, id *string, fargateTaskDefinitionArn *string) IFargateTaskDefinition {
	_init_.Initialize()

	var returns IFargateTaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateTaskDefinition",
		"fromFargateTaskDefinitionArn",
		[]interface{}{scope, id, fargateTaskDefinitionArn},
		&returns,
	)

	return returns
}

// Import an existing Fargate task definition from its attributes.
// Experimental.
func FargateTaskDefinition_FromFargateTaskDefinitionAttributes(scope constructs.Construct, id *string, attrs *FargateTaskDefinitionAttributes) IFargateTaskDefinition {
	_init_.Initialize()

	var returns IFargateTaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateTaskDefinition",
		"fromFargateTaskDefinitionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Imports a task definition from the specified task definition ARN.
//
// The task will have a compatibility of EC2+Fargate.
// Experimental.
func FargateTaskDefinition_FromTaskDefinitionArn(scope constructs.Construct, id *string, taskDefinitionArn *string) ITaskDefinition {
	_init_.Initialize()

	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateTaskDefinition",
		"fromTaskDefinitionArn",
		[]interface{}{scope, id, taskDefinitionArn},
		&returns,
	)

	return returns
}

// Create a task definition from a task definition reference.
// Experimental.
func FargateTaskDefinition_FromTaskDefinitionAttributes(scope constructs.Construct, id *string, attrs *TaskDefinitionAttributes) ITaskDefinition {
	_init_.Initialize()

	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateTaskDefinition",
		"fromTaskDefinitionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func FargateTaskDefinition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateTaskDefinition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func FargateTaskDefinition_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FargateTaskDefinition",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a new container to the task definition.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition {
	var returns ContainerDefinition

	_jsii_.Invoke(
		f,
		"addContainer",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds the specified extension to the task definition.
//
// Extension can be used to apply a packaged modification to
// a task definition.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddExtension(extension ITaskDefinitionExtension) {
	_jsii_.InvokeVoid(
		f,
		"addExtension",
		[]interface{}{extension},
	)
}

// Adds a firelens log router to the task definition.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter {
	var returns FirelensLogRouter

	_jsii_.Invoke(
		f,
		"addFirelensLogRouter",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds an inference accelerator to the task definition.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator) {
	_jsii_.InvokeVoid(
		f,
		"addInferenceAccelerator",
		[]interface{}{inferenceAccelerator},
	)
}

// Adds the specified placement constraint to the task definition.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddPlacementConstraint(constraint PlacementConstraint) {
	_jsii_.InvokeVoid(
		f,
		"addPlacementConstraint",
		[]interface{}{constraint},
	)
}

// Adds a policy statement to the task execution IAM role.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddToExecutionRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		f,
		"addToExecutionRolePolicy",
		[]interface{}{statement},
	)
}

// Adds a policy statement to the task IAM role.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddToTaskRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		f,
		"addToTaskRolePolicy",
		[]interface{}{statement},
	)
}

// Adds a volume to the task definition.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) AddVolume(volume *Volume) {
	_jsii_.InvokeVoid(
		f,
		"addVolume",
		[]interface{}{volume},
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
func (f *jsiiProxy_FargateTaskDefinition) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		f,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FargateTaskDefinition) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FargateTaskDefinition) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Creates the task execution IAM role if it doesn't already exist.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) ObtainExecutionRole() awsiam.IRole {
	var returns awsiam.IRole

	_jsii_.Invoke(
		f,
		"obtainExecutionRole",
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
func (f *jsiiProxy_FargateTaskDefinition) OnPrepare() {
	_jsii_.InvokeVoid(
		f,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
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
func (f *jsiiProxy_FargateTaskDefinition) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FargateTaskDefinition) Prepare() {
	_jsii_.InvokeVoid(
		f,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validates the task definition.
// Experimental.
func (f *jsiiProxy_FargateTaskDefinition) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Attributes used to import an existing Fargate task definition.
// Experimental.
type FargateTaskDefinitionAttributes struct {
	// The arn of the task definition.
	// Experimental.
	TaskDefinitionArn *string `json:"taskDefinitionArn"`
	// The networking mode to use for the containers in the task.
	// Experimental.
	NetworkMode NetworkMode `json:"networkMode"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
}

// The properties for a task definition.
// Experimental.
type FargateTaskDefinitionProps struct {
	// The name of the IAM task execution role that grants the ECS agent to call AWS APIs on your behalf.
	//
	// The role will be used to retrieve container images from ECR and create CloudWatch log groups.
	// Experimental.
	ExecutionRole awsiam.IRole `json:"executionRole"`
	// The name of a family that this task definition is registered to.
	//
	// A family groups multiple versions of a task definition.
	// Experimental.
	Family *string `json:"family"`
	// The configuration details for the App Mesh proxy.
	// Experimental.
	ProxyConfiguration ProxyConfiguration `json:"proxyConfiguration"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
	// The list of volume definitions for the task.
	//
	// For more information, see
	// [Task Definition Parameter Volumes](https://docs.aws.amazon.com/AmazonECS/latest/developerguide//task_definition_parameters.html#volumes).
	// Experimental.
	Volumes *[]*Volume `json:"volumes"`
	// The number of cpu units used by the task.
	//
	// For tasks using the Fargate launch type,
	// this field is required and you must use one of the following values,
	// which determines your range of valid values for the memory parameter:
	//
	// 256 (.25 vCPU) - Available memory values: 512 (0.5 GB), 1024 (1 GB), 2048 (2 GB)
	//
	// 512 (.5 vCPU) - Available memory values: 1024 (1 GB), 2048 (2 GB), 3072 (3 GB), 4096 (4 GB)
	//
	// 1024 (1 vCPU) - Available memory values: 2048 (2 GB), 3072 (3 GB), 4096 (4 GB), 5120 (5 GB), 6144 (6 GB), 7168 (7 GB), 8192 (8 GB)
	//
	// 2048 (2 vCPU) - Available memory values: Between 4096 (4 GB) and 16384 (16 GB) in increments of 1024 (1 GB)
	//
	// 4096 (4 vCPU) - Available memory values: Between 8192 (8 GB) and 30720 (30 GB) in increments of 1024 (1 GB)
	// Experimental.
	Cpu *float64 `json:"cpu"`
	// The amount (in MiB) of memory used by the task.
	//
	// For tasks using the Fargate launch type,
	// this field is required and you must use one of the following values, which determines your range of valid values for the cpu parameter:
	//
	// 512 (0.5 GB), 1024 (1 GB), 2048 (2 GB) - Available cpu values: 256 (.25 vCPU)
	//
	// 1024 (1 GB), 2048 (2 GB), 3072 (3 GB), 4096 (4 GB) - Available cpu values: 512 (.5 vCPU)
	//
	// 2048 (2 GB), 3072 (3 GB), 4096 (4 GB), 5120 (5 GB), 6144 (6 GB), 7168 (7 GB), 8192 (8 GB) - Available cpu values: 1024 (1 vCPU)
	//
	// Between 4096 (4 GB) and 16384 (16 GB) in increments of 1024 (1 GB) - Available cpu values: 2048 (2 vCPU)
	//
	// Between 8192 (8 GB) and 30720 (30 GB) in increments of 1024 (1 GB) - Available cpu values: 4096 (4 vCPU)
	// Experimental.
	MemoryLimitMiB *float64 `json:"memoryLimitMiB"`
}

// FireLens enables you to use task definition parameters to route logs to an AWS service   or AWS Partner Network (APN) destination for log storage and analytics.
// Experimental.
type FireLensLogDriver interface {
	LogDriver
	Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for FireLensLogDriver
type jsiiProxy_FireLensLogDriver struct {
	jsiiProxy_LogDriver
}

// Constructs a new instance of the FireLensLogDriver class.
// Experimental.
func NewFireLensLogDriver(props *FireLensLogDriverProps) FireLensLogDriver {
	_init_.Initialize()

	j := jsiiProxy_FireLensLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.FireLensLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the FireLensLogDriver class.
// Experimental.
func NewFireLensLogDriver_Override(f FireLensLogDriver, props *FireLensLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.FireLensLogDriver",
		[]interface{}{props},
		f,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func FireLensLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FireLensLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (f *jsiiProxy_FireLensLogDriver) Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		f,
		"bind",
		[]interface{}{_scope, _containerDefinition},
		&returns,
	)

	return returns
}

// Specifies the firelens log driver configuration options.
// Experimental.
type FireLensLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
	// The configuration options to send to the log driver.
	// Experimental.
	Options *map[string]*string `json:"options"`
}

// Firelens Configuration https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html#firelens-taskdef.
// Experimental.
type FirelensConfig struct {
	// The log router to use.
	// Experimental.
	Type FirelensLogRouterType `json:"type"`
	// Firelens options.
	// Experimental.
	Options *FirelensOptions `json:"options"`
}

// Firelens configuration file type, s3 or file path.
//
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html#firelens-taskdef-customconfig
// Experimental.
type FirelensConfigFileType string

const (
	FirelensConfigFileType_S3 FirelensConfigFileType = "S3"
	FirelensConfigFileType_FILE FirelensConfigFileType = "FILE"
)

// Firelens log router.
// Experimental.
type FirelensLogRouter interface {
	ContainerDefinition
	ContainerDependencies() *[]*ContainerDependency
	ContainerName() *string
	ContainerPort() *float64
	EnvironmentFiles() *[]*EnvironmentFileConfig
	Essential() *bool
	FirelensConfig() *FirelensConfig
	IngressPort() *float64
	LinuxParameters() LinuxParameters
	LogDriverConfig() *LogDriverConfig
	MemoryLimitSpecified() *bool
	MountPoints() *[]*MountPoint
	Node() awscdk.ConstructNode
	PortMappings() *[]*PortMapping
	ReferencesSecretJsonField() *bool
	TaskDefinition() TaskDefinition
	Ulimits() *[]*Ulimit
	VolumesFrom() *[]*VolumeFrom
	AddContainerDependencies(containerDependencies ...*ContainerDependency)
	AddInferenceAcceleratorResource(inferenceAcceleratorResources ...*string)
	AddLink(container ContainerDefinition, alias *string)
	AddMountPoints(mountPoints ...*MountPoint)
	AddPortMappings(portMappings ...*PortMapping)
	AddScratch(scratch *ScratchSpace)
	AddToExecutionPolicy(statement awsiam.PolicyStatement)
	AddUlimits(ulimits ...*Ulimit)
	AddVolumesFrom(volumesFrom ...*VolumeFrom)
	FindPortMapping(containerPort *float64, protocol Protocol) *PortMapping
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderContainerDefinition(_taskDefinition TaskDefinition) *CfnTaskDefinition_ContainerDefinitionProperty
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for FirelensLogRouter
type jsiiProxy_FirelensLogRouter struct {
	jsiiProxy_ContainerDefinition
}

func (j *jsiiProxy_FirelensLogRouter) ContainerDependencies() *[]*ContainerDependency {
	var returns *[]*ContainerDependency
	_jsii_.Get(
		j,
		"containerDependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) ContainerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"containerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) ContainerPort() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"containerPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) EnvironmentFiles() *[]*EnvironmentFileConfig {
	var returns *[]*EnvironmentFileConfig
	_jsii_.Get(
		j,
		"environmentFiles",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) Essential() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"essential",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) FirelensConfig() *FirelensConfig {
	var returns *FirelensConfig
	_jsii_.Get(
		j,
		"firelensConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) IngressPort() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"ingressPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) LinuxParameters() LinuxParameters {
	var returns LinuxParameters
	_jsii_.Get(
		j,
		"linuxParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) LogDriverConfig() *LogDriverConfig {
	var returns *LogDriverConfig
	_jsii_.Get(
		j,
		"logDriverConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) MemoryLimitSpecified() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"memoryLimitSpecified",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) MountPoints() *[]*MountPoint {
	var returns *[]*MountPoint
	_jsii_.Get(
		j,
		"mountPoints",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) PortMappings() *[]*PortMapping {
	var returns *[]*PortMapping
	_jsii_.Get(
		j,
		"portMappings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) ReferencesSecretJsonField() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"referencesSecretJsonField",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) TaskDefinition() TaskDefinition {
	var returns TaskDefinition
	_jsii_.Get(
		j,
		"taskDefinition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) Ulimits() *[]*Ulimit {
	var returns *[]*Ulimit
	_jsii_.Get(
		j,
		"ulimits",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FirelensLogRouter) VolumesFrom() *[]*VolumeFrom {
	var returns *[]*VolumeFrom
	_jsii_.Get(
		j,
		"volumesFrom",
		&returns,
	)
	return returns
}


// Constructs a new instance of the FirelensLogRouter class.
// Experimental.
func NewFirelensLogRouter(scope constructs.Construct, id *string, props *FirelensLogRouterProps) FirelensLogRouter {
	_init_.Initialize()

	j := jsiiProxy_FirelensLogRouter{}

	_jsii_.Create(
		"monocdk.aws_ecs.FirelensLogRouter",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the FirelensLogRouter class.
// Experimental.
func NewFirelensLogRouter_Override(f FirelensLogRouter, scope constructs.Construct, id *string, props *FirelensLogRouterProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.FirelensLogRouter",
		[]interface{}{scope, id, props},
		f,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func FirelensLogRouter_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FirelensLogRouter",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// This method adds one or more container dependencies to the container.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddContainerDependencies(containerDependencies ...*ContainerDependency) {
	args := []interface{}{}
	for _, a := range containerDependencies {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"addContainerDependencies",
		args,
	)
}

// This method adds one or more resources to the container.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddInferenceAcceleratorResource(inferenceAcceleratorResources ...*string) {
	args := []interface{}{}
	for _, a := range inferenceAcceleratorResources {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"addInferenceAcceleratorResource",
		args,
	)
}

// This method adds a link which allows containers to communicate with each other without the need for port mappings.
//
// This parameter is only supported if the task definition is using the bridge network mode.
// Warning: The --link flag is a legacy feature of Docker. It may eventually be removed.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddLink(container ContainerDefinition, alias *string) {
	_jsii_.InvokeVoid(
		f,
		"addLink",
		[]interface{}{container, alias},
	)
}

// This method adds one or more mount points for data volumes to the container.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddMountPoints(mountPoints ...*MountPoint) {
	args := []interface{}{}
	for _, a := range mountPoints {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"addMountPoints",
		args,
	)
}

// This method adds one or more port mappings to the container.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddPortMappings(portMappings ...*PortMapping) {
	args := []interface{}{}
	for _, a := range portMappings {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"addPortMappings",
		args,
	)
}

// This method mounts temporary disk space to the container.
//
// This adds the correct container mountPoint and task definition volume.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddScratch(scratch *ScratchSpace) {
	_jsii_.InvokeVoid(
		f,
		"addScratch",
		[]interface{}{scratch},
	)
}

// This method adds the specified statement to the IAM task execution policy in the task definition.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddToExecutionPolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		f,
		"addToExecutionPolicy",
		[]interface{}{statement},
	)
}

// This method adds one or more ulimits to the container.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddUlimits(ulimits ...*Ulimit) {
	args := []interface{}{}
	for _, a := range ulimits {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"addUlimits",
		args,
	)
}

// This method adds one or more volumes to the container.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) AddVolumesFrom(volumesFrom ...*VolumeFrom) {
	args := []interface{}{}
	for _, a := range volumesFrom {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"addVolumesFrom",
		args,
	)
}

// Returns the host port for the requested container port if it exists.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) FindPortMapping(containerPort *float64, protocol Protocol) *PortMapping {
	var returns *PortMapping

	_jsii_.Invoke(
		f,
		"findPortMapping",
		[]interface{}{containerPort, protocol},
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
func (f *jsiiProxy_FirelensLogRouter) OnPrepare() {
	_jsii_.InvokeVoid(
		f,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
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
func (f *jsiiProxy_FirelensLogRouter) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FirelensLogRouter) Prepare() {
	_jsii_.InvokeVoid(
		f,
		"prepare",
		nil, // no parameters
	)
}

// Render this container definition to a CloudFormation object.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) RenderContainerDefinition(_taskDefinition TaskDefinition) *CfnTaskDefinition_ContainerDefinitionProperty {
	var returns *CfnTaskDefinition_ContainerDefinitionProperty

	_jsii_.Invoke(
		f,
		"renderContainerDefinition",
		[]interface{}{_taskDefinition},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (f *jsiiProxy_FirelensLogRouter) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		f,
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
func (f *jsiiProxy_FirelensLogRouter) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The options for creating a firelens log router.
// Experimental.
type FirelensLogRouterDefinitionOptions struct {
	// The image used to start a container.
	//
	// This string is passed directly to the Docker daemon.
	// Images in the Docker Hub registry are available by default.
	// Other repositories are specified with either repository-url/image:tag or repository-url/image@digest.
	// TODO: Update these to specify using classes of IContainerImage
	// Experimental.
	Image ContainerImage `json:"image"`
	// The command that is passed to the container.
	//
	// If you provide a shell command as a single string, you have to quote command-line arguments.
	// Experimental.
	Command *[]*string `json:"command"`
	// The name of the container.
	// Experimental.
	ContainerName *string `json:"containerName"`
	// The minimum number of CPU units to reserve for the container.
	// Experimental.
	Cpu *float64 `json:"cpu"`
	// Specifies whether networking is disabled within the container.
	//
	// When this parameter is true, networking is disabled within the container.
	// Experimental.
	DisableNetworking *bool `json:"disableNetworking"`
	// A list of DNS search domains that are presented to the container.
	// Experimental.
	DnsSearchDomains *[]*string `json:"dnsSearchDomains"`
	// A list of DNS servers that are presented to the container.
	// Experimental.
	DnsServers *[]*string `json:"dnsServers"`
	// A key/value map of labels to add to the container.
	// Experimental.
	DockerLabels *map[string]*string `json:"dockerLabels"`
	// A list of strings to provide custom labels for SELinux and AppArmor multi-level security systems.
	// Experimental.
	DockerSecurityOptions *[]*string `json:"dockerSecurityOptions"`
	// The ENTRYPOINT value to pass to the container.
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	EntryPoint *[]*string `json:"entryPoint"`
	// The environment variables to pass to the container.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The environment files to pass to the container.
	// See: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/taskdef-envfiles.html
	//
	// Experimental.
	EnvironmentFiles *[]EnvironmentFile `json:"environmentFiles"`
	// Specifies whether the container is marked essential.
	//
	// If the essential parameter of a container is marked as true, and that container fails
	// or stops for any reason, all other containers that are part of the task are stopped.
	// If the essential parameter of a container is marked as false, then its failure does not
	// affect the rest of the containers in a task. All tasks must have at least one essential container.
	//
	// If this parameter is omitted, a container is assumed to be essential.
	// Experimental.
	Essential *bool `json:"essential"`
	// A list of hostnames and IP address mappings to append to the /etc/hosts file on the container.
	// Experimental.
	ExtraHosts *map[string]*string `json:"extraHosts"`
	// The number of GPUs assigned to the container.
	// Experimental.
	GpuCount *float64 `json:"gpuCount"`
	// The health check command and associated configuration parameters for the container.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The hostname to use for your container.
	// Experimental.
	Hostname *string `json:"hostname"`
	// The inference accelerators referenced by the container.
	// Experimental.
	InferenceAcceleratorResources *[]*string `json:"inferenceAcceleratorResources"`
	// Linux-specific modifications that are applied to the container, such as Linux kernel capabilities.
	//
	// For more information see [KernelCapabilities](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_KernelCapabilities.html).
	// Experimental.
	LinuxParameters LinuxParameters `json:"linuxParameters"`
	// The log configuration specification for the container.
	// Experimental.
	Logging LogDriver `json:"logging"`
	// The amount (in MiB) of memory to present to the container.
	//
	// If your container attempts to exceed the allocated memory, the container
	// is terminated.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryLimitMiB *float64 `json:"memoryLimitMiB"`
	// The soft limit (in MiB) of memory to reserve for the container.
	//
	// When system memory is under heavy contention, Docker attempts to keep the
	// container memory to this soft limit. However, your container can consume more
	// memory when it needs to, up to either the hard limit specified with the memory
	// parameter (if applicable), or all of the available memory on the container
	// instance, whichever comes first.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryReservationMiB *float64 `json:"memoryReservationMiB"`
	// The port mappings to add to the container definition.
	// Experimental.
	PortMappings *[]*PortMapping `json:"portMappings"`
	// Specifies whether the container is marked as privileged.
	//
	// When this parameter is true, the container is given elevated privileges on the host container instance (similar to the root user).
	// Experimental.
	Privileged *bool `json:"privileged"`
	// When this parameter is true, the container is given read-only access to its root file system.
	// Experimental.
	ReadonlyRootFilesystem *bool `json:"readonlyRootFilesystem"`
	// The secret environment variables to pass to the container.
	// Experimental.
	Secrets *map[string]Secret `json:"secrets"`
	// Time duration (in seconds) to wait before giving up on resolving dependencies for a container.
	// Experimental.
	StartTimeout awscdk.Duration `json:"startTimeout"`
	// Time duration (in seconds) to wait before the container is forcefully killed if it doesn't exit normally on its own.
	// Experimental.
	StopTimeout awscdk.Duration `json:"stopTimeout"`
	// The user name to use inside the container.
	// Experimental.
	User *string `json:"user"`
	// The working directory in which to run commands inside the container.
	// Experimental.
	WorkingDirectory *string `json:"workingDirectory"`
	// Firelens configuration.
	// Experimental.
	FirelensConfig *FirelensConfig `json:"firelensConfig"`
}

// The properties in a firelens log router.
// Experimental.
type FirelensLogRouterProps struct {
	// The image used to start a container.
	//
	// This string is passed directly to the Docker daemon.
	// Images in the Docker Hub registry are available by default.
	// Other repositories are specified with either repository-url/image:tag or repository-url/image@digest.
	// TODO: Update these to specify using classes of IContainerImage
	// Experimental.
	Image ContainerImage `json:"image"`
	// The command that is passed to the container.
	//
	// If you provide a shell command as a single string, you have to quote command-line arguments.
	// Experimental.
	Command *[]*string `json:"command"`
	// The name of the container.
	// Experimental.
	ContainerName *string `json:"containerName"`
	// The minimum number of CPU units to reserve for the container.
	// Experimental.
	Cpu *float64 `json:"cpu"`
	// Specifies whether networking is disabled within the container.
	//
	// When this parameter is true, networking is disabled within the container.
	// Experimental.
	DisableNetworking *bool `json:"disableNetworking"`
	// A list of DNS search domains that are presented to the container.
	// Experimental.
	DnsSearchDomains *[]*string `json:"dnsSearchDomains"`
	// A list of DNS servers that are presented to the container.
	// Experimental.
	DnsServers *[]*string `json:"dnsServers"`
	// A key/value map of labels to add to the container.
	// Experimental.
	DockerLabels *map[string]*string `json:"dockerLabels"`
	// A list of strings to provide custom labels for SELinux and AppArmor multi-level security systems.
	// Experimental.
	DockerSecurityOptions *[]*string `json:"dockerSecurityOptions"`
	// The ENTRYPOINT value to pass to the container.
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	EntryPoint *[]*string `json:"entryPoint"`
	// The environment variables to pass to the container.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The environment files to pass to the container.
	// See: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/taskdef-envfiles.html
	//
	// Experimental.
	EnvironmentFiles *[]EnvironmentFile `json:"environmentFiles"`
	// Specifies whether the container is marked essential.
	//
	// If the essential parameter of a container is marked as true, and that container fails
	// or stops for any reason, all other containers that are part of the task are stopped.
	// If the essential parameter of a container is marked as false, then its failure does not
	// affect the rest of the containers in a task. All tasks must have at least one essential container.
	//
	// If this parameter is omitted, a container is assumed to be essential.
	// Experimental.
	Essential *bool `json:"essential"`
	// A list of hostnames and IP address mappings to append to the /etc/hosts file on the container.
	// Experimental.
	ExtraHosts *map[string]*string `json:"extraHosts"`
	// The number of GPUs assigned to the container.
	// Experimental.
	GpuCount *float64 `json:"gpuCount"`
	// The health check command and associated configuration parameters for the container.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The hostname to use for your container.
	// Experimental.
	Hostname *string `json:"hostname"`
	// The inference accelerators referenced by the container.
	// Experimental.
	InferenceAcceleratorResources *[]*string `json:"inferenceAcceleratorResources"`
	// Linux-specific modifications that are applied to the container, such as Linux kernel capabilities.
	//
	// For more information see [KernelCapabilities](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_KernelCapabilities.html).
	// Experimental.
	LinuxParameters LinuxParameters `json:"linuxParameters"`
	// The log configuration specification for the container.
	// Experimental.
	Logging LogDriver `json:"logging"`
	// The amount (in MiB) of memory to present to the container.
	//
	// If your container attempts to exceed the allocated memory, the container
	// is terminated.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryLimitMiB *float64 `json:"memoryLimitMiB"`
	// The soft limit (in MiB) of memory to reserve for the container.
	//
	// When system memory is under heavy contention, Docker attempts to keep the
	// container memory to this soft limit. However, your container can consume more
	// memory when it needs to, up to either the hard limit specified with the memory
	// parameter (if applicable), or all of the available memory on the container
	// instance, whichever comes first.
	//
	// At least one of memoryLimitMiB and memoryReservationMiB is required for non-Fargate services.
	// Experimental.
	MemoryReservationMiB *float64 `json:"memoryReservationMiB"`
	// The port mappings to add to the container definition.
	// Experimental.
	PortMappings *[]*PortMapping `json:"portMappings"`
	// Specifies whether the container is marked as privileged.
	//
	// When this parameter is true, the container is given elevated privileges on the host container instance (similar to the root user).
	// Experimental.
	Privileged *bool `json:"privileged"`
	// When this parameter is true, the container is given read-only access to its root file system.
	// Experimental.
	ReadonlyRootFilesystem *bool `json:"readonlyRootFilesystem"`
	// The secret environment variables to pass to the container.
	// Experimental.
	Secrets *map[string]Secret `json:"secrets"`
	// Time duration (in seconds) to wait before giving up on resolving dependencies for a container.
	// Experimental.
	StartTimeout awscdk.Duration `json:"startTimeout"`
	// Time duration (in seconds) to wait before the container is forcefully killed if it doesn't exit normally on its own.
	// Experimental.
	StopTimeout awscdk.Duration `json:"stopTimeout"`
	// The user name to use inside the container.
	// Experimental.
	User *string `json:"user"`
	// The working directory in which to run commands inside the container.
	// Experimental.
	WorkingDirectory *string `json:"workingDirectory"`
	// The name of the task definition that includes this container definition.
	//
	// [disable-awslint:ref-via-interface]
	// Experimental.
	TaskDefinition TaskDefinition `json:"taskDefinition"`
	// Firelens configuration.
	// Experimental.
	FirelensConfig *FirelensConfig `json:"firelensConfig"`
}

// Firelens log router type, fluentbit or fluentd.
//
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html
// Experimental.
type FirelensLogRouterType string

const (
	FirelensLogRouterType_FLUENTBIT FirelensLogRouterType = "FLUENTBIT"
	FirelensLogRouterType_FLUENTD FirelensLogRouterType = "FLUENTD"
)

// The options for firelens log router https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html#firelens-taskdef-customconfig.
// Experimental.
type FirelensOptions struct {
	// Custom configuration file, S3 ARN or a file path.
	// Experimental.
	ConfigFileValue *string `json:"configFileValue"`
	// Custom configuration file, s3 or file.
	// Experimental.
	ConfigFileType FirelensConfigFileType `json:"configFileType"`
	// By default, Amazon ECS adds additional fields in your log entries that help identify the source of the logs.
	//
	// You can disable this action by setting enable-ecs-log-metadata to false.
	// Experimental.
	EnableECSLogMetadata *bool `json:"enableECSLogMetadata"`
}

// A log driver that sends log information to journald Logs.
// Experimental.
type FluentdLogDriver interface {
	LogDriver
	Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for FluentdLogDriver
type jsiiProxy_FluentdLogDriver struct {
	jsiiProxy_LogDriver
}

// Constructs a new instance of the FluentdLogDriver class.
// Experimental.
func NewFluentdLogDriver(props *FluentdLogDriverProps) FluentdLogDriver {
	_init_.Initialize()

	j := jsiiProxy_FluentdLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.FluentdLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the FluentdLogDriver class.
// Experimental.
func NewFluentdLogDriver_Override(f FluentdLogDriver, props *FluentdLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.FluentdLogDriver",
		[]interface{}{props},
		f,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func FluentdLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.FluentdLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (f *jsiiProxy_FluentdLogDriver) Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		f,
		"bind",
		[]interface{}{_scope, _containerDefinition},
		&returns,
	)

	return returns
}

// Specifies the fluentd log driver configuration options.
//
// [Source](https://docs.docker.com/config/containers/logging/fluentd/)
// Experimental.
type FluentdLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
	// By default, the logging driver connects to localhost:24224.
	//
	// Supply the
	// address option to connect to a different address. tcp(default) and unix
	// sockets are supported.
	// Experimental.
	Address *string `json:"address"`
	// Docker connects to Fluentd in the background.
	//
	// Messages are buffered until
	// the connection is established.
	// Experimental.
	AsyncConnect *bool `json:"asyncConnect"`
	// The amount of data to buffer before flushing to disk.
	// Experimental.
	BufferLimit *float64 `json:"bufferLimit"`
	// The maximum number of retries.
	// Experimental.
	MaxRetries *float64 `json:"maxRetries"`
	// How long to wait between retries.
	// Experimental.
	RetryWait awscdk.Duration `json:"retryWait"`
	// Generates event logs in nanosecond resolution.
	// Experimental.
	SubSecondPrecision *bool `json:"subSecondPrecision"`
}

// The type of compression the GELF driver uses to compress each log message.
// Experimental.
type GelfCompressionType string

const (
	GelfCompressionType_GZIP GelfCompressionType = "GZIP"
	GelfCompressionType_ZLIB GelfCompressionType = "ZLIB"
	GelfCompressionType_NONE GelfCompressionType = "NONE"
)

// A log driver that sends log information to journald Logs.
// Experimental.
type GelfLogDriver interface {
	LogDriver
	Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for GelfLogDriver
type jsiiProxy_GelfLogDriver struct {
	jsiiProxy_LogDriver
}

// Constructs a new instance of the GelfLogDriver class.
// Experimental.
func NewGelfLogDriver(props *GelfLogDriverProps) GelfLogDriver {
	_init_.Initialize()

	j := jsiiProxy_GelfLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.GelfLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the GelfLogDriver class.
// Experimental.
func NewGelfLogDriver_Override(g GelfLogDriver, props *GelfLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.GelfLogDriver",
		[]interface{}{props},
		g,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func GelfLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.GelfLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (g *jsiiProxy_GelfLogDriver) Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		g,
		"bind",
		[]interface{}{_scope, _containerDefinition},
		&returns,
	)

	return returns
}

// Specifies the journald log driver configuration options.
//
// [Source](https://docs.docker.com/config/containers/logging/gelf/)
// Experimental.
type GelfLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
	// The address of the GELF server.
	//
	// tcp and udp are the only supported URI
	// specifier and you must specify the port.
	// Experimental.
	Address *string `json:"address"`
	// UDP Only The level of compression when gzip or zlib is the gelf-compression-type.
	//
	// An integer in the range of -1 to 9 (BestCompression). Higher levels provide more
	// compression at lower speed. Either -1 or 0 disables compression.
	// Experimental.
	CompressionLevel *float64 `json:"compressionLevel"`
	// UDP Only The type of compression the GELF driver uses to compress each log message.
	//
	// Allowed values are gzip, zlib and none.
	// Experimental.
	CompressionType GelfCompressionType `json:"compressionType"`
	// TCP Only The maximum number of reconnection attempts when the connection drop.
	//
	// A positive integer.
	// Experimental.
	TcpMaxReconnect *float64 `json:"tcpMaxReconnect"`
	// TCP Only The number of seconds to wait between reconnection attempts.
	//
	// A positive integer.
	// Experimental.
	TcpReconnectDelay awscdk.Duration `json:"tcpReconnectDelay"`
}

// The health check command and associated configuration parameters for the container.
// Experimental.
type HealthCheck struct {
	// A string array representing the command that the container runs to determine if it is healthy.
	//
	// The string array must start with CMD to execute the command arguments directly, or
	// CMD-SHELL to run the command with the container's default shell.
	//
	// For example: [ "CMD-SHELL", "curl -f http://localhost/ || exit 1" ]
	// Experimental.
	Command *[]*string `json:"command"`
	// The time period in seconds between each health check execution.
	//
	// You may specify between 5 and 300 seconds.
	// Experimental.
	Interval awscdk.Duration `json:"interval"`
	// The number of times to retry a failed health check before the container is considered unhealthy.
	//
	// You may specify between 1 and 10 retries.
	// Experimental.
	Retries *float64 `json:"retries"`
	// The optional grace period within which to provide containers time to bootstrap before failed health checks count towards the maximum number of retries.
	//
	// You may specify between 0 and 300 seconds.
	// Experimental.
	StartPeriod awscdk.Duration `json:"startPeriod"`
	// The time period in seconds to wait for a health check to succeed before it is considered a failure.
	//
	// You may specify between 2 and 60 seconds.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
}

// The details on a container instance bind mount host volume.
// Experimental.
type Host struct {
	// Specifies the path on the host container instance that is presented to the container.
	//
	// If the sourcePath value does not exist on the host container instance, the Docker daemon creates it.
	// If the location does exist, the contents of the source path folder are exported.
	//
	// This property is not supported for tasks that use the Fargate launch type.
	// Experimental.
	SourcePath *string `json:"sourcePath"`
}

// The interface for BaseService.
// Experimental.
type IBaseService interface {
	IService
	// The cluster that hosts the service.
	// Experimental.
	Cluster() ICluster
}

// The jsii proxy for IBaseService
type jsiiProxy_IBaseService struct {
	jsiiProxy_IService
}

func (j *jsiiProxy_IBaseService) Cluster() ICluster {
	var returns ICluster
	_jsii_.Get(
		j,
		"cluster",
		&returns,
	)
	return returns
}

// A regional grouping of one or more container instances on which you can run tasks and services.
// Experimental.
type ICluster interface {
	awscdk.IResource
	// The autoscaling group added to the cluster if capacity is associated to the cluster.
	// Experimental.
	AutoscalingGroup() awsautoscaling.IAutoScalingGroup
	// The Amazon Resource Name (ARN) that identifies the cluster.
	// Experimental.
	ClusterArn() *string
	// The name of the cluster.
	// Experimental.
	ClusterName() *string
	// Manage the allowed network connections for the cluster with Security Groups.
	// Experimental.
	Connections() awsec2.Connections
	// The AWS Cloud Map namespace to associate with the cluster.
	// Experimental.
	DefaultCloudMapNamespace() awsservicediscovery.INamespace
	// The execute command configuration for the cluster.
	// Experimental.
	ExecuteCommandConfiguration() *ExecuteCommandConfiguration
	// Specifies whether the cluster has EC2 instance capacity.
	// Experimental.
	HasEc2Capacity() *bool
	// The VPC associated with the cluster.
	// Experimental.
	Vpc() awsec2.IVpc
}

// The jsii proxy for ICluster
type jsiiProxy_ICluster struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_ICluster) AutoscalingGroup() awsautoscaling.IAutoScalingGroup {
	var returns awsautoscaling.IAutoScalingGroup
	_jsii_.Get(
		j,
		"autoscalingGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICluster) ClusterArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clusterArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICluster) ClusterName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clusterName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICluster) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICluster) DefaultCloudMapNamespace() awsservicediscovery.INamespace {
	var returns awsservicediscovery.INamespace
	_jsii_.Get(
		j,
		"defaultCloudMapNamespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICluster) ExecuteCommandConfiguration() *ExecuteCommandConfiguration {
	var returns *ExecuteCommandConfiguration
	_jsii_.Get(
		j,
		"executeCommandConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICluster) HasEc2Capacity() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"hasEc2Capacity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICluster) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}

// The interface for a service using the EC2 launch type on an ECS cluster.
// Experimental.
type IEc2Service interface {
	IService
}

// The jsii proxy for IEc2Service
type jsiiProxy_IEc2Service struct {
	jsiiProxy_IService
}

// The interface of a task definition run on an EC2 cluster.
// Experimental.
type IEc2TaskDefinition interface {
	ITaskDefinition
}

// The jsii proxy for IEc2TaskDefinition
type jsiiProxy_IEc2TaskDefinition struct {
	jsiiProxy_ITaskDefinition
}

// Interface for ECS load balancer target.
// Experimental.
type IEcsLoadBalancerTarget interface {
	awselasticloadbalancingv2.IApplicationLoadBalancerTarget
	awselasticloadbalancing.ILoadBalancerTarget
	awselasticloadbalancingv2.INetworkLoadBalancerTarget
}

// The jsii proxy for IEcsLoadBalancerTarget
type jsiiProxy_IEcsLoadBalancerTarget struct {
	internal.Type__awselasticloadbalancingv2IApplicationLoadBalancerTarget
	internal.Type__awselasticloadbalancingILoadBalancerTarget
	internal.Type__awselasticloadbalancingv2INetworkLoadBalancerTarget
}

// The interface for a service using the Fargate launch type on an ECS cluster.
// Experimental.
type IFargateService interface {
	IService
}

// The jsii proxy for IFargateService
type jsiiProxy_IFargateService struct {
	jsiiProxy_IService
}

// The interface of a task definition run on a Fargate cluster.
// Experimental.
type IFargateTaskDefinition interface {
	ITaskDefinition
}

// The jsii proxy for IFargateTaskDefinition
type jsiiProxy_IFargateTaskDefinition struct {
	jsiiProxy_ITaskDefinition
}

// The interface for a service.
// Experimental.
type IService interface {
	awscdk.IResource
	// The Amazon Resource Name (ARN) of the service.
	// Experimental.
	ServiceArn() *string
	// The name of the service.
	// Experimental.
	ServiceName() *string
}

// The jsii proxy for IService
type jsiiProxy_IService struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IService) ServiceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IService) ServiceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceName",
		&returns,
	)
	return returns
}

// The interface for all task definitions.
// Experimental.
type ITaskDefinition interface {
	awscdk.IResource
	// What launch types this task definition should be compatible with.
	// Experimental.
	Compatibility() Compatibility
	// Execution role for this task definition.
	// Experimental.
	ExecutionRole() awsiam.IRole
	// Return true if the task definition can be run on an EC2 cluster.
	// Experimental.
	IsEc2Compatible() *bool
	// Return true if the task definition can be run on a ECS Anywhere cluster.
	// Experimental.
	IsExternalCompatible() *bool
	// Return true if the task definition can be run on a Fargate cluster.
	// Experimental.
	IsFargateCompatible() *bool
	// The networking mode to use for the containers in the task.
	// Experimental.
	NetworkMode() NetworkMode
	// ARN of this task definition.
	// Experimental.
	TaskDefinitionArn() *string
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole() awsiam.IRole
}

// The jsii proxy for ITaskDefinition
type jsiiProxy_ITaskDefinition struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_ITaskDefinition) Compatibility() Compatibility {
	var returns Compatibility
	_jsii_.Get(
		j,
		"compatibility",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITaskDefinition) ExecutionRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"executionRole",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITaskDefinition) IsEc2Compatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isEc2Compatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITaskDefinition) IsExternalCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isExternalCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITaskDefinition) IsFargateCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isFargateCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITaskDefinition) NetworkMode() NetworkMode {
	var returns NetworkMode
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITaskDefinition) TaskDefinitionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskDefinitionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITaskDefinition) TaskRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"taskRole",
		&returns,
	)
	return returns
}

// An extension for Task Definitions.
//
// Classes that want to make changes to a TaskDefinition (such as
// adding helper containers) can implement this interface, and can
// then be "added" to a TaskDefinition like so:
//
//     taskDefinition.addExtension(new MyExtension("some_parameter"));
// Experimental.
type ITaskDefinitionExtension interface {
	// Apply the extension to the given TaskDefinition.
	// Experimental.
	Extend(taskDefinition TaskDefinition)
}

// The jsii proxy for ITaskDefinitionExtension
type jsiiProxy_ITaskDefinitionExtension struct {
	_ byte // padding
}

func (i *jsiiProxy_ITaskDefinitionExtension) Extend(taskDefinition TaskDefinition) {
	_jsii_.InvokeVoid(
		i,
		"extend",
		[]interface{}{taskDefinition},
	)
}

// Elastic Inference Accelerator.
//
// For more information, see [Elastic Inference Basics](https://docs.aws.amazon.com/elastic-inference/latest/developerguide/basics.html)
// Experimental.
type InferenceAccelerator struct {
	// The Elastic Inference accelerator device name.
	// Experimental.
	DeviceName *string `json:"deviceName"`
	// The Elastic Inference accelerator type to use.
	//
	// The allowed values are: eia2.medium, eia2.large and eia2.xlarge.
	// Experimental.
	DeviceType *string `json:"deviceType"`
}

// The IPC resource namespace to use for the containers in the task.
// Experimental.
type IpcMode string

const (
	IpcMode_NONE IpcMode = "NONE"
	IpcMode_HOST IpcMode = "HOST"
	IpcMode_TASK IpcMode = "TASK"
)

// A log driver that sends log information to journald Logs.
// Experimental.
type JournaldLogDriver interface {
	LogDriver
	Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for JournaldLogDriver
type jsiiProxy_JournaldLogDriver struct {
	jsiiProxy_LogDriver
}

// Constructs a new instance of the JournaldLogDriver class.
// Experimental.
func NewJournaldLogDriver(props *JournaldLogDriverProps) JournaldLogDriver {
	_init_.Initialize()

	j := jsiiProxy_JournaldLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.JournaldLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the JournaldLogDriver class.
// Experimental.
func NewJournaldLogDriver_Override(j JournaldLogDriver, props *JournaldLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.JournaldLogDriver",
		[]interface{}{props},
		j,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func JournaldLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.JournaldLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (j *jsiiProxy_JournaldLogDriver) Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		j,
		"bind",
		[]interface{}{_scope, _containerDefinition},
		&returns,
	)

	return returns
}

// Specifies the journald log driver configuration options.
//
// [Source](https://docs.docker.com/config/containers/logging/journald/)
// Experimental.
type JournaldLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
}

// A log driver that sends log information to json-file Logs.
// Experimental.
type JsonFileLogDriver interface {
	LogDriver
	Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for JsonFileLogDriver
type jsiiProxy_JsonFileLogDriver struct {
	jsiiProxy_LogDriver
}

// Constructs a new instance of the JsonFileLogDriver class.
// Experimental.
func NewJsonFileLogDriver(props *JsonFileLogDriverProps) JsonFileLogDriver {
	_init_.Initialize()

	j := jsiiProxy_JsonFileLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.JsonFileLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the JsonFileLogDriver class.
// Experimental.
func NewJsonFileLogDriver_Override(j JsonFileLogDriver, props *JsonFileLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.JsonFileLogDriver",
		[]interface{}{props},
		j,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func JsonFileLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.JsonFileLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (j *jsiiProxy_JsonFileLogDriver) Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		j,
		"bind",
		[]interface{}{_scope, _containerDefinition},
		&returns,
	)

	return returns
}

// Specifies the json-file log driver configuration options.
//
// [Source](https://docs.docker.com/config/containers/logging/json-file/)
// Experimental.
type JsonFileLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
	// Toggles compression for rotated logs.
	// Experimental.
	Compress *bool `json:"compress"`
	// The maximum number of log files that can be present.
	//
	// If rolling the logs creates
	// excess files, the oldest file is removed. Only effective when max-size is also set.
	// A positive integer.
	// Experimental.
	MaxFile *float64 `json:"maxFile"`
	// The maximum size of the log before it is rolled.
	//
	// A positive integer plus a modifier
	// representing the unit of measure (k, m, or g).
	// Experimental.
	MaxSize *string `json:"maxSize"`
}

// The launch type of an ECS service.
// Experimental.
type LaunchType string

const (
	LaunchType_EC2 LaunchType = "EC2"
	LaunchType_FARGATE LaunchType = "FARGATE"
	LaunchType_EXTERNAL LaunchType = "EXTERNAL"
)

// Linux-specific options that are applied to the container.
// Experimental.
type LinuxParameters interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	AddCapabilities(cap ...Capability)
	AddDevices(device ...*Device)
	AddTmpfs(tmpfs ...*Tmpfs)
	DropCapabilities(cap ...Capability)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RenderLinuxParameters() *CfnTaskDefinition_LinuxParametersProperty
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for LinuxParameters
type jsiiProxy_LinuxParameters struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_LinuxParameters) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


// Constructs a new instance of the LinuxParameters class.
// Experimental.
func NewLinuxParameters(scope constructs.Construct, id *string, props *LinuxParametersProps) LinuxParameters {
	_init_.Initialize()

	j := jsiiProxy_LinuxParameters{}

	_jsii_.Create(
		"monocdk.aws_ecs.LinuxParameters",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the LinuxParameters class.
// Experimental.
func NewLinuxParameters_Override(l LinuxParameters, scope constructs.Construct, id *string, props *LinuxParametersProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.LinuxParameters",
		[]interface{}{scope, id, props},
		l,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func LinuxParameters_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LinuxParameters",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Adds one or more Linux capabilities to the Docker configuration of a container.
//
// Only works with EC2 launch type.
// Experimental.
func (l *jsiiProxy_LinuxParameters) AddCapabilities(cap ...Capability) {
	args := []interface{}{}
	for _, a := range cap {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		l,
		"addCapabilities",
		args,
	)
}

// Adds one or more host devices to a container.
// Experimental.
func (l *jsiiProxy_LinuxParameters) AddDevices(device ...*Device) {
	args := []interface{}{}
	for _, a := range device {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		l,
		"addDevices",
		args,
	)
}

// Specifies the container path, mount options, and size (in MiB) of the tmpfs mount for a container.
//
// Only works with EC2 launch type.
// Experimental.
func (l *jsiiProxy_LinuxParameters) AddTmpfs(tmpfs ...*Tmpfs) {
	args := []interface{}{}
	for _, a := range tmpfs {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		l,
		"addTmpfs",
		args,
	)
}

// Removes one or more Linux capabilities to the Docker configuration of a container.
//
// Only works with EC2 launch type.
// Experimental.
func (l *jsiiProxy_LinuxParameters) DropCapabilities(cap ...Capability) {
	args := []interface{}{}
	for _, a := range cap {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		l,
		"dropCapabilities",
		args,
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
func (l *jsiiProxy_LinuxParameters) OnPrepare() {
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
func (l *jsiiProxy_LinuxParameters) OnSynthesize(session constructs.ISynthesisSession) {
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
func (l *jsiiProxy_LinuxParameters) OnValidate() *[]*string {
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
func (l *jsiiProxy_LinuxParameters) Prepare() {
	_jsii_.InvokeVoid(
		l,
		"prepare",
		nil, // no parameters
	)
}

// Renders the Linux parameters to a CloudFormation object.
// Experimental.
func (l *jsiiProxy_LinuxParameters) RenderLinuxParameters() *CfnTaskDefinition_LinuxParametersProperty {
	var returns *CfnTaskDefinition_LinuxParametersProperty

	_jsii_.Invoke(
		l,
		"renderLinuxParameters",
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
func (l *jsiiProxy_LinuxParameters) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (l *jsiiProxy_LinuxParameters) ToString() *string {
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
func (l *jsiiProxy_LinuxParameters) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The properties for defining Linux-specific options that are applied to the container.
// Experimental.
type LinuxParametersProps struct {
	// Specifies whether to run an init process inside the container that forwards signals and reaps processes.
	// Experimental.
	InitProcessEnabled *bool `json:"initProcessEnabled"`
	// The value for the size (in MiB) of the /dev/shm volume.
	// Experimental.
	SharedMemorySize *float64 `json:"sharedMemorySize"`
}

// Base class for configuring listener when registering targets.
// Experimental.
type ListenerConfig interface {
	AddTargets(id *string, target *LoadBalancerTargetOptions, service BaseService)
}

// The jsii proxy struct for ListenerConfig
type jsiiProxy_ListenerConfig struct {
	_ byte // padding
}

// Experimental.
func NewListenerConfig_Override(l ListenerConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.ListenerConfig",
		nil, // no parameters
		l,
	)
}

// Create a config for adding target group to ALB listener.
// Experimental.
func ListenerConfig_ApplicationListener(listener awselasticloadbalancingv2.ApplicationListener, props *awselasticloadbalancingv2.AddApplicationTargetsProps) ListenerConfig {
	_init_.Initialize()

	var returns ListenerConfig

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ListenerConfig",
		"applicationListener",
		[]interface{}{listener, props},
		&returns,
	)

	return returns
}

// Create a config for adding target group to NLB listener.
// Experimental.
func ListenerConfig_NetworkListener(listener awselasticloadbalancingv2.NetworkListener, props *awselasticloadbalancingv2.AddNetworkTargetsProps) ListenerConfig {
	_init_.Initialize()

	var returns ListenerConfig

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ListenerConfig",
		"networkListener",
		[]interface{}{listener, props},
		&returns,
	)

	return returns
}

// Create and attach a target group to listener.
// Experimental.
func (l *jsiiProxy_ListenerConfig) AddTargets(id *string, target *LoadBalancerTargetOptions, service BaseService) {
	_jsii_.InvokeVoid(
		l,
		"addTargets",
		[]interface{}{id, target, service},
	)
}

// Properties for defining an ECS target.
//
// The port mapping for it must already have been created through addPortMapping().
// Experimental.
type LoadBalancerTargetOptions struct {
	// The name of the container.
	// Experimental.
	ContainerName *string `json:"containerName"`
	// The port number of the container.
	//
	// Only applicable when using application/network load balancers.
	// Experimental.
	ContainerPort *float64 `json:"containerPort"`
	// The protocol used for the port mapping.
	//
	// Only applicable when using application load balancers.
	// Experimental.
	Protocol Protocol `json:"protocol"`
}

// The base class for log drivers.
// Experimental.
type LogDriver interface {
	Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for LogDriver
type jsiiProxy_LogDriver struct {
	_ byte // padding
}

// Experimental.
func NewLogDriver_Override(l LogDriver) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.LogDriver",
		nil, // no parameters
		l,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func LogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (l *jsiiProxy_LogDriver) Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		l,
		"bind",
		[]interface{}{scope, containerDefinition},
		&returns,
	)

	return returns
}

// The configuration to use when creating a log driver.
// Experimental.
type LogDriverConfig struct {
	// The log driver to use for the container.
	//
	// The valid values listed for this parameter are log drivers
	// that the Amazon ECS container agent can communicate with by default.
	//
	// For tasks using the Fargate launch type, the supported log drivers are awslogs, splunk, and awsfirelens.
	// For tasks using the EC2 launch type, the supported log drivers are awslogs, fluentd, gelf, json-file, journald,
	// logentries,syslog, splunk, and awsfirelens.
	//
	// For more information about using the awslogs log driver, see
	// [Using the awslogs Log Driver](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html)
	// in the Amazon Elastic Container Service Developer Guide.
	//
	// For more information about using the awsfirelens log driver, see
	// [Custom Log Routing](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html)
	// in the Amazon Elastic Container Service Developer Guide.
	// Experimental.
	LogDriver *string `json:"logDriver"`
	// The configuration options to send to the log driver.
	// Experimental.
	Options *map[string]*string `json:"options"`
}

// The base class for log drivers.
// Experimental.
type LogDrivers interface {
}

// The jsii proxy struct for LogDrivers
type jsiiProxy_LogDrivers struct {
	_ byte // padding
}

// Experimental.
func NewLogDrivers() LogDrivers {
	_init_.Initialize()

	j := jsiiProxy_LogDrivers{}

	_jsii_.Create(
		"monocdk.aws_ecs.LogDrivers",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewLogDrivers_Override(l LogDrivers) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.LogDrivers",
		nil, // no parameters
		l,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func LogDrivers_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Creates a log driver configuration that sends log information to firelens log router.
//
// For detail configurations, please refer to Amazon ECS FireLens Examples:
// https://github.com/aws-samples/amazon-ecs-firelens-examples
// Experimental.
func LogDrivers_Firelens(props *FireLensLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"firelens",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Creates a log driver configuration that sends log information to fluentd Logs.
// Experimental.
func LogDrivers_Fluentd(props *FluentdLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"fluentd",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Creates a log driver configuration that sends log information to gelf Logs.
// Experimental.
func LogDrivers_Gelf(props *GelfLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"gelf",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Creates a log driver configuration that sends log information to journald Logs.
// Experimental.
func LogDrivers_Journald(props *JournaldLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"journald",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Creates a log driver configuration that sends log information to json-file Logs.
// Experimental.
func LogDrivers_JsonFile(props *JsonFileLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"jsonFile",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Creates a log driver configuration that sends log information to splunk Logs.
// Experimental.
func LogDrivers_Splunk(props *SplunkLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"splunk",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Creates a log driver configuration that sends log information to syslog Logs.
// Experimental.
func LogDrivers_Syslog(props *SyslogLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.LogDrivers",
		"syslog",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The machine image type.
// Experimental.
type MachineImageType string

const (
	MachineImageType_AMAZON_LINUX_2 MachineImageType = "AMAZON_LINUX_2"
	MachineImageType_BOTTLEROCKET MachineImageType = "BOTTLEROCKET"
)

// The properties for enabling scaling based on memory utilization.
// Experimental.
type MemoryUtilizationScalingProps struct {
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the scalable resource. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// scalable resource.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Period after a scale in activity completes before another scale in activity can start.
	// Experimental.
	ScaleInCooldown awscdk.Duration `json:"scaleInCooldown"`
	// Period after a scale out activity completes before another scale out activity can start.
	// Experimental.
	ScaleOutCooldown awscdk.Duration `json:"scaleOutCooldown"`
	// The target value for memory utilization across all tasks in the service.
	// Experimental.
	TargetUtilizationPercent *float64 `json:"targetUtilizationPercent"`
}

// The details of data volume mount points for a container.
// Experimental.
type MountPoint struct {
	// The path on the container to mount the host volume at.
	// Experimental.
	ContainerPath *string `json:"containerPath"`
	// Specifies whether to give the container read-only access to the volume.
	//
	// If this value is true, the container has read-only access to the volume.
	// If this value is false, then the container can write to the volume.
	// Experimental.
	ReadOnly *bool `json:"readOnly"`
	// The name of the volume to mount.
	//
	// Must be a volume name referenced in the name parameter of task definition volume.
	// Experimental.
	SourceVolume *string `json:"sourceVolume"`
}

// The networking mode to use for the containers in the task.
// Experimental.
type NetworkMode string

const (
	NetworkMode_NONE NetworkMode = "NONE"
	NetworkMode_BRIDGE NetworkMode = "BRIDGE"
	NetworkMode_AWS_VPC NetworkMode = "AWS_VPC"
	NetworkMode_HOST NetworkMode = "HOST"
	NetworkMode_NAT NetworkMode = "NAT"
)

// The process namespace to use for the containers in the task.
// Experimental.
type PidMode string

const (
	PidMode_HOST PidMode = "HOST"
	PidMode_TASK PidMode = "TASK"
)

// The placement constraints to use for tasks in the service. For more information, see [Amazon ECS Task Placement Constraints](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-placement-constraints.html).
//
// Tasks will only be placed on instances that match these rules.
// Experimental.
type PlacementConstraint interface {
	ToJson() *[]*CfnService_PlacementConstraintProperty
}

// The jsii proxy struct for PlacementConstraint
type jsiiProxy_PlacementConstraint struct {
	_ byte // padding
}

// Use distinctInstance to ensure that each task in a particular group is running on a different container instance.
// Experimental.
func PlacementConstraint_DistinctInstances() PlacementConstraint {
	_init_.Initialize()

	var returns PlacementConstraint

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementConstraint",
		"distinctInstances",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Use memberOf to restrict the selection to a group of valid candidates specified by a query expression.
//
// Multiple expressions can be specified. For more information, see
// [Cluster Query Language](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cluster-query-language.html).
//
// You can specify multiple expressions in one call. The tasks will only be placed on instances matching all expressions.
// See: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cluster-query-language.html
//
// Experimental.
func PlacementConstraint_MemberOf(expressions ...*string) PlacementConstraint {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range expressions {
		args = append(args, a)
	}

	var returns PlacementConstraint

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementConstraint",
		"memberOf",
		args,
		&returns,
	)

	return returns
}

// Return the placement JSON.
// Experimental.
func (p *jsiiProxy_PlacementConstraint) ToJson() *[]*CfnService_PlacementConstraintProperty {
	var returns *[]*CfnService_PlacementConstraintProperty

	_jsii_.Invoke(
		p,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The placement strategies to use for tasks in the service. For more information, see [Amazon ECS Task Placement Strategies](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-placement-strategies.html).
//
// Tasks will preferentially be placed on instances that match these rules.
// Experimental.
type PlacementStrategy interface {
	ToJson() *[]*CfnService_PlacementStrategyProperty
}

// The jsii proxy struct for PlacementStrategy
type jsiiProxy_PlacementStrategy struct {
	_ byte // padding
}

// Places tasks on the container instances with the least available capacity of the specified resource.
// Experimental.
func PlacementStrategy_PackedBy(resource BinPackResource) PlacementStrategy {
	_init_.Initialize()

	var returns PlacementStrategy

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementStrategy",
		"packedBy",
		[]interface{}{resource},
		&returns,
	)

	return returns
}

// Places tasks on container instances with the least available amount of CPU capacity.
//
// This minimizes the number of instances in use.
// Experimental.
func PlacementStrategy_PackedByCpu() PlacementStrategy {
	_init_.Initialize()

	var returns PlacementStrategy

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementStrategy",
		"packedByCpu",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Places tasks on container instances with the least available amount of memory capacity.
//
// This minimizes the number of instances in use.
// Experimental.
func PlacementStrategy_PackedByMemory() PlacementStrategy {
	_init_.Initialize()

	var returns PlacementStrategy

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementStrategy",
		"packedByMemory",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Places tasks randomly.
// Experimental.
func PlacementStrategy_Randomly() PlacementStrategy {
	_init_.Initialize()

	var returns PlacementStrategy

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementStrategy",
		"randomly",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Places tasks evenly based on the specified value.
//
// You can use one of the built-in attributes found on `BuiltInAttributes`
// or supply your own custom instance attributes. If more than one attribute
// is supplied, spreading is done in order.
// Experimental.
func PlacementStrategy_SpreadAcross(fields ...*string) PlacementStrategy {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range fields {
		args = append(args, a)
	}

	var returns PlacementStrategy

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementStrategy",
		"spreadAcross",
		args,
		&returns,
	)

	return returns
}

// Places tasks evenly across all container instances in the cluster.
// Experimental.
func PlacementStrategy_SpreadAcrossInstances() PlacementStrategy {
	_init_.Initialize()

	var returns PlacementStrategy

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.PlacementStrategy",
		"spreadAcrossInstances",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return the placement JSON.
// Experimental.
func (p *jsiiProxy_PlacementStrategy) ToJson() *[]*CfnService_PlacementStrategyProperty {
	var returns *[]*CfnService_PlacementStrategyProperty

	_jsii_.Invoke(
		p,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Port mappings allow containers to access ports on the host container instance to send or receive traffic.
// Experimental.
type PortMapping struct {
	// The port number on the container that is bound to the user-specified or automatically assigned host port.
	//
	// If you are using containers in a task with the awsvpc or host network mode, exposed ports should be specified using containerPort.
	// If you are using containers in a task with the bridge network mode and you specify a container port and not a host port,
	// your container automatically receives a host port in the ephemeral port range.
	//
	// For more information, see hostPort.
	// Port mappings that are automatically assigned in this way do not count toward the 100 reserved ports limit of a container instance.
	// Experimental.
	ContainerPort *float64 `json:"containerPort"`
	// The port number on the container instance to reserve for your container.
	//
	// If you are using containers in a task with the awsvpc or host network mode,
	// the hostPort can either be left blank or set to the same value as the containerPort.
	//
	// If you are using containers in a task with the bridge network mode,
	// you can specify a non-reserved host port for your container port mapping, or
	// you can omit the hostPort (or set it to 0) while specifying a containerPort and
	// your container automatically receives a port in the ephemeral port range for
	// your container instance operating system and Docker version.
	// Experimental.
	HostPort *float64 `json:"hostPort"`
	// The protocol used for the port mapping.
	//
	// Valid values are Protocol.TCP and Protocol.UDP.
	// Experimental.
	Protocol Protocol `json:"protocol"`
}

// Propagate tags from either service or task definition.
// Experimental.
type PropagatedTagSource string

const (
	PropagatedTagSource_SERVICE PropagatedTagSource = "SERVICE"
	PropagatedTagSource_TASK_DEFINITION PropagatedTagSource = "TASK_DEFINITION"
	PropagatedTagSource_NONE PropagatedTagSource = "NONE"
)

// Network protocol.
// Experimental.
type Protocol string

const (
	Protocol_TCP Protocol = "TCP"
	Protocol_UDP Protocol = "UDP"
)

// The base class for proxy configurations.
// Experimental.
type ProxyConfiguration interface {
	Bind(_scope awscdk.Construct, _taskDefinition TaskDefinition) *CfnTaskDefinition_ProxyConfigurationProperty
}

// The jsii proxy struct for ProxyConfiguration
type jsiiProxy_ProxyConfiguration struct {
	_ byte // padding
}

// Experimental.
func NewProxyConfiguration_Override(p ProxyConfiguration) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.ProxyConfiguration",
		nil, // no parameters
		p,
	)
}

// Called when the proxy configuration is configured on a task definition.
// Experimental.
func (p *jsiiProxy_ProxyConfiguration) Bind(_scope awscdk.Construct, _taskDefinition TaskDefinition) *CfnTaskDefinition_ProxyConfigurationProperty {
	var returns *CfnTaskDefinition_ProxyConfigurationProperty

	_jsii_.Invoke(
		p,
		"bind",
		[]interface{}{_scope, _taskDefinition},
		&returns,
	)

	return returns
}

// The base class for proxy configurations.
// Experimental.
type ProxyConfigurations interface {
}

// The jsii proxy struct for ProxyConfigurations
type jsiiProxy_ProxyConfigurations struct {
	_ byte // padding
}

// Experimental.
func NewProxyConfigurations() ProxyConfigurations {
	_init_.Initialize()

	j := jsiiProxy_ProxyConfigurations{}

	_jsii_.Create(
		"monocdk.aws_ecs.ProxyConfigurations",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewProxyConfigurations_Override(p ProxyConfigurations) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.ProxyConfigurations",
		nil, // no parameters
		p,
	)
}

// Constructs a new instance of the ProxyConfiguration class.
// Experimental.
func ProxyConfigurations_AppMeshProxyConfiguration(props *AppMeshProxyConfigurationConfigProps) ProxyConfiguration {
	_init_.Initialize()

	var returns ProxyConfiguration

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ProxyConfigurations",
		"appMeshProxyConfiguration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// An image hosted in a public or private repository.
//
// For images hosted in Amazon ECR, see
// [EcrImage](https://docs.aws.amazon.com/AmazonECR/latest/userguide/images.html).
// Experimental.
type RepositoryImage interface {
	ContainerImage
	Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig
}

// The jsii proxy struct for RepositoryImage
type jsiiProxy_RepositoryImage struct {
	jsiiProxy_ContainerImage
}

// Constructs a new instance of the RepositoryImage class.
// Experimental.
func NewRepositoryImage(imageName *string, props *RepositoryImageProps) RepositoryImage {
	_init_.Initialize()

	j := jsiiProxy_RepositoryImage{}

	_jsii_.Create(
		"monocdk.aws_ecs.RepositoryImage",
		[]interface{}{imageName, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the RepositoryImage class.
// Experimental.
func NewRepositoryImage_Override(r RepositoryImage, imageName *string, props *RepositoryImageProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.RepositoryImage",
		[]interface{}{imageName, props},
		r,
	)
}

// Reference an image that's constructed directly from sources on disk.
//
// If you already have a `DockerImageAsset` instance, you can use the
// `ContainerImage.fromDockerImageAsset` method instead.
// Experimental.
func RepositoryImage_FromAsset(directory *string, props *AssetImageProps) AssetImage {
	_init_.Initialize()

	var returns AssetImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.RepositoryImage",
		"fromAsset",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Use an existing `DockerImageAsset` for this container image.
// Experimental.
func RepositoryImage_FromDockerImageAsset(asset awsecrassets.DockerImageAsset) ContainerImage {
	_init_.Initialize()

	var returns ContainerImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.RepositoryImage",
		"fromDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Reference an image in an ECR repository.
// Experimental.
func RepositoryImage_FromEcrRepository(repository awsecr.IRepository, tag *string) EcrImage {
	_init_.Initialize()

	var returns EcrImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.RepositoryImage",
		"fromEcrRepository",
		[]interface{}{repository, tag},
		&returns,
	)

	return returns
}

// Reference an image on DockerHub or another online registry.
// Experimental.
func RepositoryImage_FromRegistry(name *string, props *RepositoryImageProps) RepositoryImage {
	_init_.Initialize()

	var returns RepositoryImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.RepositoryImage",
		"fromRegistry",
		[]interface{}{name, props},
		&returns,
	)

	return returns
}

// Called when the image is used by a ContainerDefinition.
// Experimental.
func (r *jsiiProxy_RepositoryImage) Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig {
	var returns *ContainerImageConfig

	_jsii_.Invoke(
		r,
		"bind",
		[]interface{}{scope, containerDefinition},
		&returns,
	)

	return returns
}

// The properties for an image hosted in a public or private repository.
// Experimental.
type RepositoryImageProps struct {
	// The secret to expose to the container that contains the credentials for the image repository.
	//
	// The supported value is the full ARN of an AWS Secrets Manager secret.
	// Experimental.
	Credentials awssecretsmanager.ISecret `json:"credentials"`
}

// The properties for enabling scaling based on Application Load Balancer (ALB) request counts.
// Experimental.
type RequestCountScalingProps struct {
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the scalable resource. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// scalable resource.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Period after a scale in activity completes before another scale in activity can start.
	// Experimental.
	ScaleInCooldown awscdk.Duration `json:"scaleInCooldown"`
	// Period after a scale out activity completes before another scale out activity can start.
	// Experimental.
	ScaleOutCooldown awscdk.Duration `json:"scaleOutCooldown"`
	// The number of ALB requests per target.
	// Experimental.
	RequestsPerTarget *float64 `json:"requestsPerTarget"`
	// The ALB target group name.
	// Experimental.
	TargetGroup awselasticloadbalancingv2.ApplicationTargetGroup `json:"targetGroup"`
}

// Environment file from S3.
// Experimental.
type S3EnvironmentFile interface {
	EnvironmentFile
	Bind(_scope awscdk.Construct) *EnvironmentFileConfig
}

// The jsii proxy struct for S3EnvironmentFile
type jsiiProxy_S3EnvironmentFile struct {
	jsiiProxy_EnvironmentFile
}

// Experimental.
func NewS3EnvironmentFile(bucket awss3.IBucket, key *string, objectVersion *string) S3EnvironmentFile {
	_init_.Initialize()

	j := jsiiProxy_S3EnvironmentFile{}

	_jsii_.Create(
		"monocdk.aws_ecs.S3EnvironmentFile",
		[]interface{}{bucket, key, objectVersion},
		&j,
	)

	return &j
}

// Experimental.
func NewS3EnvironmentFile_Override(s S3EnvironmentFile, bucket awss3.IBucket, key *string, objectVersion *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.S3EnvironmentFile",
		[]interface{}{bucket, key, objectVersion},
		s,
	)
}

// Loads the environment file from a local disk path.
// Experimental.
func S3EnvironmentFile_FromAsset(path *string, options *awss3assets.AssetOptions) AssetEnvironmentFile {
	_init_.Initialize()

	var returns AssetEnvironmentFile

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.S3EnvironmentFile",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Loads the environment file from an S3 bucket.
//
// Returns: `S3EnvironmentFile` associated with the specified S3 object.
// Experimental.
func S3EnvironmentFile_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3EnvironmentFile {
	_init_.Initialize()

	var returns S3EnvironmentFile

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.S3EnvironmentFile",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Called when the container is initialized to allow this object to bind to the stack.
// Experimental.
func (s *jsiiProxy_S3EnvironmentFile) Bind(_scope awscdk.Construct) *EnvironmentFileConfig {
	var returns *EnvironmentFileConfig

	_jsii_.Invoke(
		s,
		"bind",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// The scalable attribute representing task count.
// Experimental.
type ScalableTaskCount interface {
	awsapplicationautoscaling.BaseScalableAttribute
	Node() awscdk.ConstructNode
	Props() *awsapplicationautoscaling.BaseScalableAttributeProps
	DoScaleOnMetric(id *string, props *awsapplicationautoscaling.BasicStepScalingPolicyProps)
	DoScaleOnSchedule(id *string, props *awsapplicationautoscaling.ScalingSchedule)
	DoScaleToTrackMetric(id *string, props *awsapplicationautoscaling.BasicTargetTrackingScalingPolicyProps)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	ScaleOnCpuUtilization(id *string, props *CpuUtilizationScalingProps)
	ScaleOnMemoryUtilization(id *string, props *MemoryUtilizationScalingProps)
	ScaleOnMetric(id *string, props *awsapplicationautoscaling.BasicStepScalingPolicyProps)
	ScaleOnRequestCount(id *string, props *RequestCountScalingProps)
	ScaleOnSchedule(id *string, props *awsapplicationautoscaling.ScalingSchedule)
	ScaleToTrackCustomMetric(id *string, props *TrackCustomMetricProps)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ScalableTaskCount
type jsiiProxy_ScalableTaskCount struct {
	internal.Type__awsapplicationautoscalingBaseScalableAttribute
}

func (j *jsiiProxy_ScalableTaskCount) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScalableTaskCount) Props() *awsapplicationautoscaling.BaseScalableAttributeProps {
	var returns *awsapplicationautoscaling.BaseScalableAttributeProps
	_jsii_.Get(
		j,
		"props",
		&returns,
	)
	return returns
}


// Constructs a new instance of the ScalableTaskCount class.
// Experimental.
func NewScalableTaskCount(scope constructs.Construct, id *string, props *ScalableTaskCountProps) ScalableTaskCount {
	_init_.Initialize()

	j := jsiiProxy_ScalableTaskCount{}

	_jsii_.Create(
		"monocdk.aws_ecs.ScalableTaskCount",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the ScalableTaskCount class.
// Experimental.
func NewScalableTaskCount_Override(s ScalableTaskCount, scope constructs.Construct, id *string, props *ScalableTaskCountProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.ScalableTaskCount",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func ScalableTaskCount_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.ScalableTaskCount",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Scale out or in based on a metric value.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) DoScaleOnMetric(id *string, props *awsapplicationautoscaling.BasicStepScalingPolicyProps) {
	_jsii_.InvokeVoid(
		s,
		"doScaleOnMetric",
		[]interface{}{id, props},
	)
}

// Scale out or in based on time.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) DoScaleOnSchedule(id *string, props *awsapplicationautoscaling.ScalingSchedule) {
	_jsii_.InvokeVoid(
		s,
		"doScaleOnSchedule",
		[]interface{}{id, props},
	)
}

// Scale out or in in order to keep a metric around a target value.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) DoScaleToTrackMetric(id *string, props *awsapplicationautoscaling.BasicTargetTrackingScalingPolicyProps) {
	_jsii_.InvokeVoid(
		s,
		"doScaleToTrackMetric",
		[]interface{}{id, props},
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
func (s *jsiiProxy_ScalableTaskCount) OnPrepare() {
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
func (s *jsiiProxy_ScalableTaskCount) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_ScalableTaskCount) OnValidate() *[]*string {
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
func (s *jsiiProxy_ScalableTaskCount) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Scales in or out to achieve a target CPU utilization.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) ScaleOnCpuUtilization(id *string, props *CpuUtilizationScalingProps) {
	_jsii_.InvokeVoid(
		s,
		"scaleOnCpuUtilization",
		[]interface{}{id, props},
	)
}

// Scales in or out to achieve a target memory utilization.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) ScaleOnMemoryUtilization(id *string, props *MemoryUtilizationScalingProps) {
	_jsii_.InvokeVoid(
		s,
		"scaleOnMemoryUtilization",
		[]interface{}{id, props},
	)
}

// Scales in or out based on a specified metric value.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) ScaleOnMetric(id *string, props *awsapplicationautoscaling.BasicStepScalingPolicyProps) {
	_jsii_.InvokeVoid(
		s,
		"scaleOnMetric",
		[]interface{}{id, props},
	)
}

// Scales in or out to achieve a target Application Load Balancer request count per target.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) ScaleOnRequestCount(id *string, props *RequestCountScalingProps) {
	_jsii_.InvokeVoid(
		s,
		"scaleOnRequestCount",
		[]interface{}{id, props},
	)
}

// Scales in or out based on a specified scheduled time.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) ScaleOnSchedule(id *string, props *awsapplicationautoscaling.ScalingSchedule) {
	_jsii_.InvokeVoid(
		s,
		"scaleOnSchedule",
		[]interface{}{id, props},
	)
}

// Scales in or out to achieve a target on a custom metric.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) ScaleToTrackCustomMetric(id *string, props *TrackCustomMetricProps) {
	_jsii_.InvokeVoid(
		s,
		"scaleToTrackCustomMetric",
		[]interface{}{id, props},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_ScalableTaskCount) ToString() *string {
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
func (s *jsiiProxy_ScalableTaskCount) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The properties of a scalable attribute representing task count.
// Experimental.
type ScalableTaskCountProps struct {
	// Maximum capacity to scale to.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// Minimum capacity to scale to.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
	// Scalable dimension of the attribute.
	// Experimental.
	Dimension *string `json:"dimension"`
	// Resource ID of the attribute.
	// Experimental.
	ResourceId *string `json:"resourceId"`
	// Role to use for scaling.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// Service namespace of the scalable attribute.
	// Experimental.
	ServiceNamespace awsapplicationautoscaling.ServiceNamespace `json:"serviceNamespace"`
}

// The scope for the Docker volume that determines its lifecycle.
//
// Docker volumes that are scoped to a task are automatically provisioned when the task starts and destroyed when the task stops.
// Docker volumes that are scoped as shared persist after the task stops.
// Experimental.
type Scope string

const (
	Scope_TASK Scope = "TASK"
	Scope_SHARED Scope = "SHARED"
)

// The temporary disk space mounted to the container.
// Experimental.
type ScratchSpace struct {
	// The path on the container to mount the scratch volume at.
	// Experimental.
	ContainerPath *string `json:"containerPath"`
	// The name of the scratch volume to mount.
	//
	// Must be a volume name referenced in the name parameter of task definition volume.
	// Experimental.
	Name *string `json:"name"`
	// Specifies whether to give the container read-only access to the scratch volume.
	//
	// If this value is true, the container has read-only access to the scratch volume.
	// If this value is false, then the container can write to the scratch volume.
	// Experimental.
	ReadOnly *bool `json:"readOnly"`
	// Experimental.
	SourcePath *string `json:"sourcePath"`
}

// A secret environment variable.
// Experimental.
type Secret interface {
	Arn() *string
	HasField() *bool
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
}

// The jsii proxy struct for Secret
type jsiiProxy_Secret struct {
	_ byte // padding
}

func (j *jsiiProxy_Secret) Arn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Secret) HasField() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"hasField",
		&returns,
	)
	return returns
}


// Experimental.
func NewSecret_Override(s Secret) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.Secret",
		nil, // no parameters
		s,
	)
}

// Creates a environment variable value from a secret stored in AWS Secrets Manager.
// Experimental.
func Secret_FromSecretsManager(secret awssecretsmanager.ISecret, field *string) Secret {
	_init_.Initialize()

	var returns Secret

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Secret",
		"fromSecretsManager",
		[]interface{}{secret, field},
		&returns,
	)

	return returns
}

// Creates an environment variable value from a parameter stored in AWS Systems Manager Parameter Store.
// Experimental.
func Secret_FromSsmParameter(parameter awsssm.IParameter) Secret {
	_init_.Initialize()

	var returns Secret

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.Secret",
		"fromSsmParameter",
		[]interface{}{parameter},
		&returns,
	)

	return returns
}

// Grants reading the secret to a principal.
// Experimental.
func (s *jsiiProxy_Secret) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// A log driver that sends log information to splunk Logs.
// Experimental.
type SplunkLogDriver interface {
	LogDriver
	Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for SplunkLogDriver
type jsiiProxy_SplunkLogDriver struct {
	jsiiProxy_LogDriver
}

// Constructs a new instance of the SplunkLogDriver class.
// Experimental.
func NewSplunkLogDriver(props *SplunkLogDriverProps) SplunkLogDriver {
	_init_.Initialize()

	j := jsiiProxy_SplunkLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.SplunkLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the SplunkLogDriver class.
// Experimental.
func NewSplunkLogDriver_Override(s SplunkLogDriver, props *SplunkLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.SplunkLogDriver",
		[]interface{}{props},
		s,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func SplunkLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.SplunkLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (s *jsiiProxy_SplunkLogDriver) Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		s,
		"bind",
		[]interface{}{_scope, _containerDefinition},
		&returns,
	)

	return returns
}

// Specifies the splunk log driver configuration options.
//
// [Source](https://docs.docker.com/config/containers/logging/splunk/)
// Experimental.
type SplunkLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
	// Splunk HTTP Event Collector token.
	// Experimental.
	Token awscdk.SecretValue `json:"token"`
	// Path to your Splunk Enterprise, self-service Splunk Cloud instance, or Splunk Cloud managed cluster (including port and scheme used by HTTP Event Collector) in one of the following formats: https://your_splunk_instance:8088 or https://input-prd-p-XXXXXXX.cloud.splunk.com:8088 or https://http-inputs-XXXXXXXX.splunkcloud.com.
	// Experimental.
	Url *string `json:"url"`
	// Name to use for validating server certificate.
	// Experimental.
	CaName *string `json:"caName"`
	// Path to root certificate.
	// Experimental.
	CaPath *string `json:"caPath"`
	// Message format.
	//
	// Can be inline, json or raw.
	// Experimental.
	Format SplunkLogFormat `json:"format"`
	// Enable/disable gzip compression to send events to Splunk Enterprise or Splunk Cloud instance.
	// Experimental.
	Gzip *bool `json:"gzip"`
	// Set compression level for gzip.
	//
	// Valid values are -1 (default), 0 (no compression),
	// 1 (best speed) ... 9 (best compression).
	// Experimental.
	GzipLevel *float64 `json:"gzipLevel"`
	// Event index.
	// Experimental.
	Index *string `json:"index"`
	// Ignore server certificate validation.
	// Experimental.
	InsecureSkipVerify *string `json:"insecureSkipVerify"`
	// Event source.
	// Experimental.
	Source *string `json:"source"`
	// Event source type.
	// Experimental.
	SourceType *string `json:"sourceType"`
	// Verify on start, that docker can connect to Splunk server.
	// Experimental.
	VerifyConnection *bool `json:"verifyConnection"`
}

// Log Message Format.
// Experimental.
type SplunkLogFormat string

const (
	SplunkLogFormat_INLINE SplunkLogFormat = "INLINE"
	SplunkLogFormat_JSON SplunkLogFormat = "JSON"
	SplunkLogFormat_RAW SplunkLogFormat = "RAW"
)

// A log driver that sends log information to syslog Logs.
// Experimental.
type SyslogLogDriver interface {
	LogDriver
	Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig
}

// The jsii proxy struct for SyslogLogDriver
type jsiiProxy_SyslogLogDriver struct {
	jsiiProxy_LogDriver
}

// Constructs a new instance of the SyslogLogDriver class.
// Experimental.
func NewSyslogLogDriver(props *SyslogLogDriverProps) SyslogLogDriver {
	_init_.Initialize()

	j := jsiiProxy_SyslogLogDriver{}

	_jsii_.Create(
		"monocdk.aws_ecs.SyslogLogDriver",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Constructs a new instance of the SyslogLogDriver class.
// Experimental.
func NewSyslogLogDriver_Override(s SyslogLogDriver, props *SyslogLogDriverProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.SyslogLogDriver",
		[]interface{}{props},
		s,
	)
}

// Creates a log driver configuration that sends log information to CloudWatch Logs.
// Experimental.
func SyslogLogDriver_AwsLogs(props *AwsLogDriverProps) LogDriver {
	_init_.Initialize()

	var returns LogDriver

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.SyslogLogDriver",
		"awsLogs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called when the log driver is configured on a container.
// Experimental.
func (s *jsiiProxy_SyslogLogDriver) Bind(_scope awscdk.Construct, _containerDefinition ContainerDefinition) *LogDriverConfig {
	var returns *LogDriverConfig

	_jsii_.Invoke(
		s,
		"bind",
		[]interface{}{_scope, _containerDefinition},
		&returns,
	)

	return returns
}

// Specifies the syslog log driver configuration options.
//
// [Source](https://docs.docker.com/config/containers/logging/syslog/)
// Experimental.
type SyslogLogDriverProps struct {
	// The env option takes an array of keys.
	//
	// If there is collision between
	// label and env keys, the value of the env takes precedence. Adds additional fields
	// to the extra attributes of a logging message.
	// Experimental.
	Env *[]*string `json:"env"`
	// The env-regex option is similar to and compatible with env.
	//
	// Its value is a regular
	// expression to match logging-related environment variables. It is used for advanced
	// log tag options.
	// Experimental.
	EnvRegex *string `json:"envRegex"`
	// The labels option takes an array of keys.
	//
	// If there is collision
	// between label and env keys, the value of the env takes precedence. Adds additional
	// fields to the extra attributes of a logging message.
	// Experimental.
	Labels *[]*string `json:"labels"`
	// By default, Docker uses the first 12 characters of the container ID to tag log messages.
	//
	// Refer to the log tag option documentation for customizing the
	// log tag format.
	// Experimental.
	Tag *string `json:"tag"`
	// The address of an external syslog server.
	//
	// The URI specifier may be
	// [tcp|udp|tcp+tls]://host:port, unix://path, or unixgram://path.
	// Experimental.
	Address *string `json:"address"`
	// The syslog facility to use.
	//
	// Can be the number or name for any valid
	// syslog facility. See the syslog documentation:
	// https://tools.ietf.org/html/rfc5424#section-6.2.1.
	// Experimental.
	Facility *string `json:"facility"`
	// The syslog message format to use.
	//
	// If not specified the local UNIX syslog
	// format is used, without a specified hostname. Specify rfc3164 for the RFC-3164
	// compatible format, rfc5424 for RFC-5424 compatible format, or rfc5424micro
	// for RFC-5424 compatible format with microsecond timestamp resolution.
	// Experimental.
	Format *string `json:"format"`
	// The absolute path to the trust certificates signed by the CA.
	//
	// Ignored
	// if the address protocol is not tcp+tls.
	// Experimental.
	TlsCaCert *string `json:"tlsCaCert"`
	// The absolute path to the TLS certificate file.
	//
	// Ignored if the address
	// protocol is not tcp+tls.
	// Experimental.
	TlsCert *string `json:"tlsCert"`
	// The absolute path to the TLS key file.
	//
	// Ignored if the address protocol
	// is not tcp+tls.
	// Experimental.
	TlsKey *string `json:"tlsKey"`
	// If set to true, TLS verification is skipped when connecting to the syslog daemon.
	//
	// Ignored if the address protocol is not tcp+tls.
	// Experimental.
	TlsSkipVerify *bool `json:"tlsSkipVerify"`
}

// A special type of {@link ContainerImage} that uses an ECR repository for the image, but a CloudFormation Parameter for the tag of the image in that repository.
//
// This allows providing this tag through the Parameter at deploy time,
// for example in a CodePipeline that pushes a new tag of the image to the repository during a build step,
// and then provides that new tag through the CloudFormation Parameter in the deploy step.
// See: #tagParameterName
//
// Experimental.
type TagParameterContainerImage interface {
	ContainerImage
	TagParameterName() *string
	TagParameterValue() *string
	Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig
}

// The jsii proxy struct for TagParameterContainerImage
type jsiiProxy_TagParameterContainerImage struct {
	jsiiProxy_ContainerImage
}

func (j *jsiiProxy_TagParameterContainerImage) TagParameterName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tagParameterName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TagParameterContainerImage) TagParameterValue() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tagParameterValue",
		&returns,
	)
	return returns
}


// Experimental.
func NewTagParameterContainerImage(repository awsecr.IRepository) TagParameterContainerImage {
	_init_.Initialize()

	j := jsiiProxy_TagParameterContainerImage{}

	_jsii_.Create(
		"monocdk.aws_ecs.TagParameterContainerImage",
		[]interface{}{repository},
		&j,
	)

	return &j
}

// Experimental.
func NewTagParameterContainerImage_Override(t TagParameterContainerImage, repository awsecr.IRepository) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.TagParameterContainerImage",
		[]interface{}{repository},
		t,
	)
}

// Reference an image that's constructed directly from sources on disk.
//
// If you already have a `DockerImageAsset` instance, you can use the
// `ContainerImage.fromDockerImageAsset` method instead.
// Experimental.
func TagParameterContainerImage_FromAsset(directory *string, props *AssetImageProps) AssetImage {
	_init_.Initialize()

	var returns AssetImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TagParameterContainerImage",
		"fromAsset",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Use an existing `DockerImageAsset` for this container image.
// Experimental.
func TagParameterContainerImage_FromDockerImageAsset(asset awsecrassets.DockerImageAsset) ContainerImage {
	_init_.Initialize()

	var returns ContainerImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TagParameterContainerImage",
		"fromDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Reference an image in an ECR repository.
// Experimental.
func TagParameterContainerImage_FromEcrRepository(repository awsecr.IRepository, tag *string) EcrImage {
	_init_.Initialize()

	var returns EcrImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TagParameterContainerImage",
		"fromEcrRepository",
		[]interface{}{repository, tag},
		&returns,
	)

	return returns
}

// Reference an image on DockerHub or another online registry.
// Experimental.
func TagParameterContainerImage_FromRegistry(name *string, props *RepositoryImageProps) RepositoryImage {
	_init_.Initialize()

	var returns RepositoryImage

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TagParameterContainerImage",
		"fromRegistry",
		[]interface{}{name, props},
		&returns,
	)

	return returns
}

// Called when the image is used by a ContainerDefinition.
// Experimental.
func (t *jsiiProxy_TagParameterContainerImage) Bind(scope awscdk.Construct, containerDefinition ContainerDefinition) *ContainerImageConfig {
	var returns *ContainerImageConfig

	_jsii_.Invoke(
		t,
		"bind",
		[]interface{}{scope, containerDefinition},
		&returns,
	)

	return returns
}

// The base class for all task definitions.
// Experimental.
type TaskDefinition interface {
	awscdk.Resource
	ITaskDefinition
	Compatibility() Compatibility
	Containers() *[]ContainerDefinition
	DefaultContainer() ContainerDefinition
	SetDefaultContainer(val ContainerDefinition)
	Env() *awscdk.ResourceEnvironment
	ExecutionRole() awsiam.IRole
	Family() *string
	InferenceAccelerators() *[]*InferenceAccelerator
	IsEc2Compatible() *bool
	IsExternalCompatible() *bool
	IsFargateCompatible() *bool
	NetworkMode() NetworkMode
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ReferencesSecretJsonField() *bool
	Stack() awscdk.Stack
	TaskDefinitionArn() *string
	TaskRole() awsiam.IRole
	AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition
	AddExtension(extension ITaskDefinitionExtension)
	AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter
	AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator)
	AddPlacementConstraint(constraint PlacementConstraint)
	AddToExecutionRolePolicy(statement awsiam.PolicyStatement)
	AddToTaskRolePolicy(statement awsiam.PolicyStatement)
	AddVolume(volume *Volume)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	ObtainExecutionRole() awsiam.IRole
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for TaskDefinition
type jsiiProxy_TaskDefinition struct {
	internal.Type__awscdkResource
	jsiiProxy_ITaskDefinition
}

func (j *jsiiProxy_TaskDefinition) Compatibility() Compatibility {
	var returns Compatibility
	_jsii_.Get(
		j,
		"compatibility",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) Containers() *[]ContainerDefinition {
	var returns *[]ContainerDefinition
	_jsii_.Get(
		j,
		"containers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) DefaultContainer() ContainerDefinition {
	var returns ContainerDefinition
	_jsii_.Get(
		j,
		"defaultContainer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) ExecutionRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"executionRole",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) Family() *string {
	var returns *string
	_jsii_.Get(
		j,
		"family",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) InferenceAccelerators() *[]*InferenceAccelerator {
	var returns *[]*InferenceAccelerator
	_jsii_.Get(
		j,
		"inferenceAccelerators",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) IsEc2Compatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isEc2Compatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) IsExternalCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isExternalCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) IsFargateCompatible() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isFargateCompatible",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) NetworkMode() NetworkMode {
	var returns NetworkMode
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) ReferencesSecretJsonField() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"referencesSecretJsonField",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) TaskDefinitionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"taskDefinitionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TaskDefinition) TaskRole() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"taskRole",
		&returns,
	)
	return returns
}


// Constructs a new instance of the TaskDefinition class.
// Experimental.
func NewTaskDefinition(scope constructs.Construct, id *string, props *TaskDefinitionProps) TaskDefinition {
	_init_.Initialize()

	j := jsiiProxy_TaskDefinition{}

	_jsii_.Create(
		"monocdk.aws_ecs.TaskDefinition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Constructs a new instance of the TaskDefinition class.
// Experimental.
func NewTaskDefinition_Override(t TaskDefinition, scope constructs.Construct, id *string, props *TaskDefinitionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecs.TaskDefinition",
		[]interface{}{scope, id, props},
		t,
	)
}

func (j *jsiiProxy_TaskDefinition) SetDefaultContainer(val ContainerDefinition) {
	_jsii_.Set(
		j,
		"defaultContainer",
		val,
	)
}

// Imports a task definition from the specified task definition ARN.
//
// The task will have a compatibility of EC2+Fargate.
// Experimental.
func TaskDefinition_FromTaskDefinitionArn(scope constructs.Construct, id *string, taskDefinitionArn *string) ITaskDefinition {
	_init_.Initialize()

	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TaskDefinition",
		"fromTaskDefinitionArn",
		[]interface{}{scope, id, taskDefinitionArn},
		&returns,
	)

	return returns
}

// Create a task definition from a task definition reference.
// Experimental.
func TaskDefinition_FromTaskDefinitionAttributes(scope constructs.Construct, id *string, attrs *TaskDefinitionAttributes) ITaskDefinition {
	_init_.Initialize()

	var returns ITaskDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TaskDefinition",
		"fromTaskDefinitionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func TaskDefinition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TaskDefinition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func TaskDefinition_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecs.TaskDefinition",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a new container to the task definition.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddContainer(id *string, props *ContainerDefinitionOptions) ContainerDefinition {
	var returns ContainerDefinition

	_jsii_.Invoke(
		t,
		"addContainer",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds the specified extension to the task definition.
//
// Extension can be used to apply a packaged modification to
// a task definition.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddExtension(extension ITaskDefinitionExtension) {
	_jsii_.InvokeVoid(
		t,
		"addExtension",
		[]interface{}{extension},
	)
}

// Adds a firelens log router to the task definition.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddFirelensLogRouter(id *string, props *FirelensLogRouterDefinitionOptions) FirelensLogRouter {
	var returns FirelensLogRouter

	_jsii_.Invoke(
		t,
		"addFirelensLogRouter",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds an inference accelerator to the task definition.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddInferenceAccelerator(inferenceAccelerator *InferenceAccelerator) {
	_jsii_.InvokeVoid(
		t,
		"addInferenceAccelerator",
		[]interface{}{inferenceAccelerator},
	)
}

// Adds the specified placement constraint to the task definition.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddPlacementConstraint(constraint PlacementConstraint) {
	_jsii_.InvokeVoid(
		t,
		"addPlacementConstraint",
		[]interface{}{constraint},
	)
}

// Adds a policy statement to the task execution IAM role.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddToExecutionRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		t,
		"addToExecutionRolePolicy",
		[]interface{}{statement},
	)
}

// Adds a policy statement to the task IAM role.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddToTaskRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		t,
		"addToTaskRolePolicy",
		[]interface{}{statement},
	)
}

// Adds a volume to the task definition.
// Experimental.
func (t *jsiiProxy_TaskDefinition) AddVolume(volume *Volume) {
	_jsii_.InvokeVoid(
		t,
		"addVolume",
		[]interface{}{volume},
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
func (t *jsiiProxy_TaskDefinition) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		t,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (t *jsiiProxy_TaskDefinition) GeneratePhysicalName() *string {
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
func (t *jsiiProxy_TaskDefinition) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (t *jsiiProxy_TaskDefinition) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Creates the task execution IAM role if it doesn't already exist.
// Experimental.
func (t *jsiiProxy_TaskDefinition) ObtainExecutionRole() awsiam.IRole {
	var returns awsiam.IRole

	_jsii_.Invoke(
		t,
		"obtainExecutionRole",
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
func (t *jsiiProxy_TaskDefinition) OnPrepare() {
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
func (t *jsiiProxy_TaskDefinition) OnSynthesize(session constructs.ISynthesisSession) {
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
func (t *jsiiProxy_TaskDefinition) OnValidate() *[]*string {
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
func (t *jsiiProxy_TaskDefinition) Prepare() {
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
func (t *jsiiProxy_TaskDefinition) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TaskDefinition) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validates the task definition.
// Experimental.
func (t *jsiiProxy_TaskDefinition) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A reference to an existing task definition.
// Experimental.
type TaskDefinitionAttributes struct {
	// The arn of the task definition.
	// Experimental.
	TaskDefinitionArn *string `json:"taskDefinitionArn"`
	// The networking mode to use for the containers in the task.
	// Experimental.
	NetworkMode NetworkMode `json:"networkMode"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
	// What launch types this task definition should be compatible with.
	// Experimental.
	Compatibility Compatibility `json:"compatibility"`
}

// The properties for task definitions.
// Experimental.
type TaskDefinitionProps struct {
	// The name of the IAM task execution role that grants the ECS agent to call AWS APIs on your behalf.
	//
	// The role will be used to retrieve container images from ECR and create CloudWatch log groups.
	// Experimental.
	ExecutionRole awsiam.IRole `json:"executionRole"`
	// The name of a family that this task definition is registered to.
	//
	// A family groups multiple versions of a task definition.
	// Experimental.
	Family *string `json:"family"`
	// The configuration details for the App Mesh proxy.
	// Experimental.
	ProxyConfiguration ProxyConfiguration `json:"proxyConfiguration"`
	// The name of the IAM role that grants containers in the task permission to call AWS APIs on your behalf.
	// Experimental.
	TaskRole awsiam.IRole `json:"taskRole"`
	// The list of volume definitions for the task.
	//
	// For more information, see
	// [Task Definition Parameter Volumes](https://docs.aws.amazon.com/AmazonECS/latest/developerguide//task_definition_parameters.html#volumes).
	// Experimental.
	Volumes *[]*Volume `json:"volumes"`
	// The task launch type compatiblity requirement.
	// Experimental.
	Compatibility Compatibility `json:"compatibility"`
	// The number of cpu units used by the task.
	//
	// If you are using the EC2 launch type, this field is optional and any value can be used.
	// If you are using the Fargate launch type, this field is required and you must use one of the following values,
	// which determines your range of valid values for the memory parameter:
	//
	// 256 (.25 vCPU) - Available memory values: 512 (0.5 GB), 1024 (1 GB), 2048 (2 GB)
	//
	// 512 (.5 vCPU) - Available memory values: 1024 (1 GB), 2048 (2 GB), 3072 (3 GB), 4096 (4 GB)
	//
	// 1024 (1 vCPU) - Available memory values: 2048 (2 GB), 3072 (3 GB), 4096 (4 GB), 5120 (5 GB), 6144 (6 GB), 7168 (7 GB), 8192 (8 GB)
	//
	// 2048 (2 vCPU) - Available memory values: Between 4096 (4 GB) and 16384 (16 GB) in increments of 1024 (1 GB)
	//
	// 4096 (4 vCPU) - Available memory values: Between 8192 (8 GB) and 30720 (30 GB) in increments of 1024 (1 GB)
	// Experimental.
	Cpu *string `json:"cpu"`
	// The inference accelerators to use for the containers in the task.
	//
	// Not supported in Fargate.
	// Experimental.
	InferenceAccelerators *[]*InferenceAccelerator `json:"inferenceAccelerators"`
	// The IPC resource namespace to use for the containers in the task.
	//
	// Not supported in Fargate and Windows containers.
	// Experimental.
	IpcMode IpcMode `json:"ipcMode"`
	// The amount (in MiB) of memory used by the task.
	//
	// If using the EC2 launch type, this field is optional and any value can be used.
	// If using the Fargate launch type, this field is required and you must use one of the following values,
	// which determines your range of valid values for the cpu parameter:
	//
	// 512 (0.5 GB), 1024 (1 GB), 2048 (2 GB) - Available cpu values: 256 (.25 vCPU)
	//
	// 1024 (1 GB), 2048 (2 GB), 3072 (3 GB), 4096 (4 GB) - Available cpu values: 512 (.5 vCPU)
	//
	// 2048 (2 GB), 3072 (3 GB), 4096 (4 GB), 5120 (5 GB), 6144 (6 GB), 7168 (7 GB), 8192 (8 GB) - Available cpu values: 1024 (1 vCPU)
	//
	// Between 4096 (4 GB) and 16384 (16 GB) in increments of 1024 (1 GB) - Available cpu values: 2048 (2 vCPU)
	//
	// Between 8192 (8 GB) and 30720 (30 GB) in increments of 1024 (1 GB) - Available cpu values: 4096 (4 vCPU)
	// Experimental.
	MemoryMiB *string `json:"memoryMiB"`
	// The networking mode to use for the containers in the task.
	//
	// On Fargate, the only supported networking mode is AwsVpc.
	// Experimental.
	NetworkMode NetworkMode `json:"networkMode"`
	// The process namespace to use for the containers in the task.
	//
	// Not supported in Fargate and Windows containers.
	// Experimental.
	PidMode PidMode `json:"pidMode"`
	// The placement constraints to use for tasks in the service.
	//
	// You can specify a maximum of 10 constraints per task (this limit includes
	// constraints in the task definition and those specified at run time).
	//
	// Not supported in Fargate.
	// Experimental.
	PlacementConstraints *[]PlacementConstraint `json:"placementConstraints"`
}

// The details of a tmpfs mount for a container.
// Experimental.
type Tmpfs struct {
	// The absolute file path where the tmpfs volume is to be mounted.
	// Experimental.
	ContainerPath *string `json:"containerPath"`
	// The size (in MiB) of the tmpfs volume.
	// Experimental.
	Size *float64 `json:"size"`
	// The list of tmpfs volume mount options.
	//
	// For more information, see
	// [TmpfsMountOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Tmpfs.html).
	// Experimental.
	MountOptions *[]TmpfsMountOption `json:"mountOptions"`
}

// The supported options for a tmpfs mount for a container.
// Experimental.
type TmpfsMountOption string

const (
	TmpfsMountOption_DEFAULTS TmpfsMountOption = "DEFAULTS"
	TmpfsMountOption_RO TmpfsMountOption = "RO"
	TmpfsMountOption_RW TmpfsMountOption = "RW"
	TmpfsMountOption_SUID TmpfsMountOption = "SUID"
	TmpfsMountOption_NOSUID TmpfsMountOption = "NOSUID"
	TmpfsMountOption_DEV TmpfsMountOption = "DEV"
	TmpfsMountOption_NODEV TmpfsMountOption = "NODEV"
	TmpfsMountOption_EXEC TmpfsMountOption = "EXEC"
	TmpfsMountOption_NOEXEC TmpfsMountOption = "NOEXEC"
	TmpfsMountOption_SYNC TmpfsMountOption = "SYNC"
	TmpfsMountOption_ASYNC TmpfsMountOption = "ASYNC"
	TmpfsMountOption_DIRSYNC TmpfsMountOption = "DIRSYNC"
	TmpfsMountOption_REMOUNT TmpfsMountOption = "REMOUNT"
	TmpfsMountOption_MAND TmpfsMountOption = "MAND"
	TmpfsMountOption_NOMAND TmpfsMountOption = "NOMAND"
	TmpfsMountOption_ATIME TmpfsMountOption = "ATIME"
	TmpfsMountOption_NOATIME TmpfsMountOption = "NOATIME"
	TmpfsMountOption_DIRATIME TmpfsMountOption = "DIRATIME"
	TmpfsMountOption_NODIRATIME TmpfsMountOption = "NODIRATIME"
	TmpfsMountOption_BIND TmpfsMountOption = "BIND"
	TmpfsMountOption_RBIND TmpfsMountOption = "RBIND"
	TmpfsMountOption_UNBINDABLE TmpfsMountOption = "UNBINDABLE"
	TmpfsMountOption_RUNBINDABLE TmpfsMountOption = "RUNBINDABLE"
	TmpfsMountOption_PRIVATE TmpfsMountOption = "PRIVATE"
	TmpfsMountOption_RPRIVATE TmpfsMountOption = "RPRIVATE"
	TmpfsMountOption_SHARED TmpfsMountOption = "SHARED"
	TmpfsMountOption_RSHARED TmpfsMountOption = "RSHARED"
	TmpfsMountOption_SLAVE TmpfsMountOption = "SLAVE"
	TmpfsMountOption_RSLAVE TmpfsMountOption = "RSLAVE"
	TmpfsMountOption_RELATIME TmpfsMountOption = "RELATIME"
	TmpfsMountOption_NORELATIME TmpfsMountOption = "NORELATIME"
	TmpfsMountOption_STRICTATIME TmpfsMountOption = "STRICTATIME"
	TmpfsMountOption_NOSTRICTATIME TmpfsMountOption = "NOSTRICTATIME"
	TmpfsMountOption_MODE TmpfsMountOption = "MODE"
	TmpfsMountOption_UID TmpfsMountOption = "UID"
	TmpfsMountOption_GID TmpfsMountOption = "GID"
	TmpfsMountOption_NR_INODES TmpfsMountOption = "NR_INODES"
	TmpfsMountOption_NR_BLOCKS TmpfsMountOption = "NR_BLOCKS"
	TmpfsMountOption_MPOL TmpfsMountOption = "MPOL"
)

// The properties for enabling target tracking scaling based on a custom CloudWatch metric.
// Experimental.
type TrackCustomMetricProps struct {
	// Indicates whether scale in by the target tracking policy is disabled.
	//
	// If the value is true, scale in is disabled and the target tracking policy
	// won't remove capacity from the scalable resource. Otherwise, scale in is
	// enabled and the target tracking policy can remove capacity from the
	// scalable resource.
	// Experimental.
	DisableScaleIn *bool `json:"disableScaleIn"`
	// A name for the scaling policy.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Period after a scale in activity completes before another scale in activity can start.
	// Experimental.
	ScaleInCooldown awscdk.Duration `json:"scaleInCooldown"`
	// Period after a scale out activity completes before another scale out activity can start.
	// Experimental.
	ScaleOutCooldown awscdk.Duration `json:"scaleOutCooldown"`
	// The custom CloudWatch metric to track.
	//
	// The metric must represent utilization; that is, you will always get the following behavior:
	//
	// - metric > targetValue => scale out
	// - metric < targetValue => scale in
	// Experimental.
	Metric awscloudwatch.IMetric `json:"metric"`
	// The target value for the custom CloudWatch metric.
	// Experimental.
	TargetValue *float64 `json:"targetValue"`
}

// The ulimit settings to pass to the container.
//
// NOTE: Does not work for Windows containers.
// Experimental.
type Ulimit struct {
	// The hard limit for the ulimit type.
	// Experimental.
	HardLimit *float64 `json:"hardLimit"`
	// The type of the ulimit.
	//
	// For more information, see [UlimitName](https://docs.aws.amazon.com/cdk/api/latest/typescript/api/aws-ecs/ulimitname.html#aws_ecs_UlimitName).
	// Experimental.
	Name UlimitName `json:"name"`
	// The soft limit for the ulimit type.
	// Experimental.
	SoftLimit *float64 `json:"softLimit"`
}

// Type of resource to set a limit on.
// Experimental.
type UlimitName string

const (
	UlimitName_CORE UlimitName = "CORE"
	UlimitName_CPU UlimitName = "CPU"
	UlimitName_DATA UlimitName = "DATA"
	UlimitName_FSIZE UlimitName = "FSIZE"
	UlimitName_LOCKS UlimitName = "LOCKS"
	UlimitName_MEMLOCK UlimitName = "MEMLOCK"
	UlimitName_MSGQUEUE UlimitName = "MSGQUEUE"
	UlimitName_NICE UlimitName = "NICE"
	UlimitName_NOFILE UlimitName = "NOFILE"
	UlimitName_NPROC UlimitName = "NPROC"
	UlimitName_RSS UlimitName = "RSS"
	UlimitName_RTPRIO UlimitName = "RTPRIO"
	UlimitName_RTTIME UlimitName = "RTTIME"
	UlimitName_SIGPENDING UlimitName = "SIGPENDING"
	UlimitName_STACK UlimitName = "STACK"
)

// A data volume used in a task definition.
//
// For tasks that use a Docker volume, specify a DockerVolumeConfiguration.
// For tasks that use a bind mount host volume, specify a host and optional sourcePath.
//
// For more information, see [Using Data Volumes in Tasks](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_data_volumes.html).
// Experimental.
type Volume struct {
	// The name of the volume.
	//
	// Up to 255 letters (uppercase and lowercase), numbers, and hyphens are allowed.
	// This name is referenced in the sourceVolume parameter of container definition mountPoints.
	// Experimental.
	Name *string `json:"name"`
	// This property is specified when you are using Docker volumes.
	//
	// Docker volumes are only supported when you are using the EC2 launch type.
	// Windows containers only support the use of the local driver.
	// To use bind mounts, specify a host instead.
	// Experimental.
	DockerVolumeConfiguration *DockerVolumeConfiguration `json:"dockerVolumeConfiguration"`
	// This property is specified when you are using Amazon EFS.
	//
	// When specifying Amazon EFS volumes in tasks using the Fargate launch type,
	// Fargate creates a supervisor container that is responsible for managing the Amazon EFS volume.
	// The supervisor container uses a small amount of the task's memory.
	// The supervisor container is visible when querying the task metadata version 4 endpoint,
	// but is not visible in CloudWatch Container Insights.
	// Experimental.
	EfsVolumeConfiguration *EfsVolumeConfiguration `json:"efsVolumeConfiguration"`
	// This property is specified when you are using bind mount host volumes.
	//
	// Bind mount host volumes are supported when you are using either the EC2 or Fargate launch types.
	// The contents of the host parameter determine whether your bind mount host volume persists on the
	// host container instance and where it is stored. If the host parameter is empty, then the Docker
	// daemon assigns a host path for your data volume. However, the data is not guaranteed to persist
	// after the containers associated with it stop running.
	// Experimental.
	Host *Host `json:"host"`
}

// The details on a data volume from another container in the same task definition.
// Experimental.
type VolumeFrom struct {
	// Specifies whether the container has read-only access to the volume.
	//
	// If this value is true, the container has read-only access to the volume.
	// If this value is false, then the container can write to the volume.
	// Experimental.
	ReadOnly *bool `json:"readOnly"`
	// The name of another container within the same task definition from which to mount volumes.
	// Experimental.
	SourceContainer *string `json:"sourceContainer"`
}

// ECS-optimized Windows version list.
// Experimental.
type WindowsOptimizedVersion string

const (
	WindowsOptimizedVersion_SERVER_2019 WindowsOptimizedVersion = "SERVER_2019"
	WindowsOptimizedVersion_SERVER_2016 WindowsOptimizedVersion = "SERVER_2016"
)


package awslambda

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/assets"
	"github.com/aws/aws-cdk-go/awscdk/awsapplicationautoscaling"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awscodeguruprofiler"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/awsefs"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awslambda/internal"
	"github.com/aws/aws-cdk-go/awscdk/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/awssigner"
	"github.com/aws/aws-cdk-go/awscdk/awssqs"
	"github.com/aws/constructs-go/constructs/v3"
)

// A new alias to a particular version of a Lambda function.
// Experimental.
type Alias interface {
	QualifiedFunctionBase
	IAlias
	AliasName() *string
	CanCreatePermissions() *bool
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	FunctionArn() *string
	FunctionName() *string
	GrantPrincipal() awsiam.IPrincipal
	IsBoundToVpc() *bool
	Lambda() IFunction
	LatestVersion() IVersion
	Node() awscdk.ConstructNode
	PermissionsNode() awscdk.ConstructNode
	PhysicalName() *string
	Qualifier() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	Version() IVersion
	AddAutoScaling(options *AutoScalingOptions) IScalableFunctionAttribute
	AddEventSource(source IEventSource)
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	AddPermission(id *string, permission *Permission)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
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
	jsiiProxy_QualifiedFunctionBase
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

func (j *jsiiProxy_Alias) CanCreatePermissions() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"canCreatePermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
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

func (j *jsiiProxy_Alias) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) Lambda() IFunction {
	var returns IFunction
	_jsii_.Get(
		j,
		"lambda",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
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

func (j *jsiiProxy_Alias) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
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

func (j *jsiiProxy_Alias) Qualifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"qualifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Alias) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
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

func (j *jsiiProxy_Alias) Version() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


// Experimental.
func NewAlias(scope constructs.Construct, id *string, props *AliasProps) Alias {
	_init_.Initialize()

	j := jsiiProxy_Alias{}

	_jsii_.Create(
		"monocdk.aws_lambda.Alias",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewAlias_Override(a Alias, scope constructs.Construct, id *string, props *AliasProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.Alias",
		[]interface{}{scope, id, props},
		a,
	)
}

// Experimental.
func Alias_FromAliasAttributes(scope constructs.Construct, id *string, attrs *AliasAttributes) IAlias {
	_init_.Initialize()

	var returns IAlias

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Alias",
		"fromAliasAttributes",
		[]interface{}{scope, id, attrs},
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
		"monocdk.aws_lambda.Alias",
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
		"monocdk.aws_lambda.Alias",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Configure provisioned concurrency autoscaling on a function alias.
//
// Returns a scalable attribute that can call
// `scaleOnUtilization()` and `scaleOnSchedule()`.
// Experimental.
func (a *jsiiProxy_Alias) AddAutoScaling(options *AutoScalingOptions) IScalableFunctionAttribute {
	var returns IScalableFunctionAttribute

	_jsii_.Invoke(
		a,
		"addAutoScaling",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Adds an event source to this function.
//
// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
//
// The following example adds an SQS Queue as an event source:
// ```
// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
// myFunction.addEventSource(new SqsEventSource(myQueue));
// ```
// Experimental.
func (a *jsiiProxy_Alias) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		a,
		"addEventSource",
		[]interface{}{source},
	)
}

// Adds an event source that maps to this AWS Lambda function.
// Experimental.
func (a *jsiiProxy_Alias) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		a,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a permission to the Lambda resource policy.
// See: Permission for details.
//
// Experimental.
func (a *jsiiProxy_Alias) AddPermission(id *string, permission *Permission) {
	_jsii_.InvokeVoid(
		a,
		"addPermission",
		[]interface{}{id, permission},
	)
}

// Adds a statement to the IAM role assumed by the instance.
// Experimental.
func (a *jsiiProxy_Alias) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		a,
		"addToRolePolicy",
		[]interface{}{statement},
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
func (a *jsiiProxy_Alias) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Configures options for asynchronous invocation.
// Experimental.
func (a *jsiiProxy_Alias) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		a,
		"configureAsyncInvoke",
		[]interface{}{options},
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

// Grant the given identity permissions to invoke this Lambda.
// Experimental.
func (a *jsiiProxy_Alias) GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grantInvoke",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Function.
// Experimental.
func (a *jsiiProxy_Alias) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// How long execution of this Lambda takes.
//
// Average over 5 minutes
// Experimental.
func (a *jsiiProxy_Alias) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How many invocations of this Lambda fail.
//
// Sum over 5 minutes
// Experimental.
func (a *jsiiProxy_Alias) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is invoked.
//
// Sum over 5 minutes
// Experimental.
func (a *jsiiProxy_Alias) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is throttled.
//
// Sum over 5 minutes
// Experimental.
func (a *jsiiProxy_Alias) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricThrottles",
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

// Experimental.
type AliasAttributes struct {
	// Experimental.
	AliasName *string `json:"aliasName"`
	// Experimental.
	AliasVersion IVersion `json:"aliasVersion"`
}

// Options for `lambda.Alias`.
// Experimental.
type AliasOptions struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Additional versions with individual weights this alias points to.
	//
	// Individual additional version weights specified here should add up to
	// (less than) one. All remaining weight is routed to the default
	// version.
	//
	// For example, the config is
	//
	//     version: "1"
	//     additionalVersions: [{ version: "2", weight: 0.05 }]
	//
	// Then 5% of traffic will be routed to function version 2, while
	// the remaining 95% of traffic will be routed to function version 1.
	// Experimental.
	AdditionalVersions *[]*VersionWeight `json:"additionalVersions"`
	// Description for the alias.
	// Experimental.
	Description *string `json:"description"`
	// Specifies a provisioned concurrency configuration for a function's alias.
	// Experimental.
	ProvisionedConcurrentExecutions *float64 `json:"provisionedConcurrentExecutions"`
}

// Properties for a new Lambda alias.
// Experimental.
type AliasProps struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Additional versions with individual weights this alias points to.
	//
	// Individual additional version weights specified here should add up to
	// (less than) one. All remaining weight is routed to the default
	// version.
	//
	// For example, the config is
	//
	//     version: "1"
	//     additionalVersions: [{ version: "2", weight: 0.05 }]
	//
	// Then 5% of traffic will be routed to function version 2, while
	// the remaining 95% of traffic will be routed to function version 1.
	// Experimental.
	AdditionalVersions *[]*VersionWeight `json:"additionalVersions"`
	// Description for the alias.
	// Experimental.
	Description *string `json:"description"`
	// Specifies a provisioned concurrency configuration for a function's alias.
	// Experimental.
	ProvisionedConcurrentExecutions *float64 `json:"provisionedConcurrentExecutions"`
	// Name of this alias.
	// Experimental.
	AliasName *string `json:"aliasName"`
	// Function version this alias refers to.
	//
	// Use lambda.addVersion() to obtain a new lambda version to refer to.
	// Experimental.
	Version IVersion `json:"version"`
}

// Lambda code from a local directory.
// Experimental.
type AssetCode interface {
	Code
	IsInline() *bool
	Path() *string
	Bind(scope awscdk.Construct) *CodeConfig
	BindToResource(resource awscdk.CfnResource, options *ResourceBindOptions)
}

// The jsii proxy struct for AssetCode
type jsiiProxy_AssetCode struct {
	jsiiProxy_Code
}

func (j *jsiiProxy_AssetCode) IsInline() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isInline",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetCode) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}


// Experimental.
func NewAssetCode(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	j := jsiiProxy_AssetCode{}

	_jsii_.Create(
		"monocdk.aws_lambda.AssetCode",
		[]interface{}{path, options},
		&j,
	)

	return &j
}

// Experimental.
func NewAssetCode_Override(a AssetCode, path *string, options *awss3assets.AssetOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.AssetCode",
		[]interface{}{path, options},
		a,
	)
}

// DEPRECATED.
// Deprecated: use `fromAsset`
func AssetCode_Asset(path *string) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"asset",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromBucket`
func AssetCode_Bucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"bucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromCfnParameters`
func AssetCode_CfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"cfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from a local disk path.
// Experimental.
func AssetCode_FromAsset(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func AssetCode_FromAssetImage(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	var returns AssetImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"fromAssetImage",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Lambda handler code as an S3 object.
// Experimental.
func AssetCode_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Creates a new Lambda source defined using CloudFormation parameters.
//
// Returns: a new instance of `CfnParametersCode`
// Experimental.
func AssetCode_FromCfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"fromCfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from an asset created by a Docker build.
//
// By default, the asset is expected to be located at `/asset` in the
// image.
// Experimental.
func AssetCode_FromDockerBuild(path *string, options *DockerBuildAssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"fromDockerBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func AssetCode_FromEcrImage(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	var returns EcrImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"fromEcrImage",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Inline code for Lambda handler.
//
// Returns: `LambdaInlineCode` with inline code.
// Experimental.
func AssetCode_FromInline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"fromInline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromInline`
func AssetCode_Inline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetCode",
		"inline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// Called when the lambda or layer is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (a *jsiiProxy_AssetCode) Bind(scope awscdk.Construct) *CodeConfig {
	var returns *CodeConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Called after the CFN function resource has been created to allow the code class to bind to it.
//
// Specifically it's required to allow assets to add
// metadata for tooling like SAM CLI to be able to find their origins.
// Experimental.
func (a *jsiiProxy_AssetCode) BindToResource(resource awscdk.CfnResource, options *ResourceBindOptions) {
	_jsii_.InvokeVoid(
		a,
		"bindToResource",
		[]interface{}{resource, options},
	)
}

// Represents an ECR image that will be constructed from the specified asset and can be bound as Lambda code.
// Experimental.
type AssetImageCode interface {
	Code
	IsInline() *bool
	Bind(scope awscdk.Construct) *CodeConfig
	BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions)
}

// The jsii proxy struct for AssetImageCode
type jsiiProxy_AssetImageCode struct {
	jsiiProxy_Code
}

func (j *jsiiProxy_AssetImageCode) IsInline() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isInline",
		&returns,
	)
	return returns
}


// Experimental.
func NewAssetImageCode(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	j := jsiiProxy_AssetImageCode{}

	_jsii_.Create(
		"monocdk.aws_lambda.AssetImageCode",
		[]interface{}{directory, props},
		&j,
	)

	return &j
}

// Experimental.
func NewAssetImageCode_Override(a AssetImageCode, directory *string, props *AssetImageCodeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.AssetImageCode",
		[]interface{}{directory, props},
		a,
	)
}

// DEPRECATED.
// Deprecated: use `fromAsset`
func AssetImageCode_Asset(path *string) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"asset",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromBucket`
func AssetImageCode_Bucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"bucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromCfnParameters`
func AssetImageCode_CfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"cfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from a local disk path.
// Experimental.
func AssetImageCode_FromAsset(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func AssetImageCode_FromAssetImage(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	var returns AssetImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"fromAssetImage",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Lambda handler code as an S3 object.
// Experimental.
func AssetImageCode_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Creates a new Lambda source defined using CloudFormation parameters.
//
// Returns: a new instance of `CfnParametersCode`
// Experimental.
func AssetImageCode_FromCfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"fromCfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from an asset created by a Docker build.
//
// By default, the asset is expected to be located at `/asset` in the
// image.
// Experimental.
func AssetImageCode_FromDockerBuild(path *string, options *DockerBuildAssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"fromDockerBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func AssetImageCode_FromEcrImage(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	var returns EcrImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"fromEcrImage",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Inline code for Lambda handler.
//
// Returns: `LambdaInlineCode` with inline code.
// Experimental.
func AssetImageCode_FromInline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"fromInline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromInline`
func AssetImageCode_Inline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.AssetImageCode",
		"inline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// Called when the lambda or layer is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (a *jsiiProxy_AssetImageCode) Bind(scope awscdk.Construct) *CodeConfig {
	var returns *CodeConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Called after the CFN function resource has been created to allow the code class to bind to it.
//
// Specifically it's required to allow assets to add
// metadata for tooling like SAM CLI to be able to find their origins.
// Experimental.
func (a *jsiiProxy_AssetImageCode) BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions) {
	_jsii_.InvokeVoid(
		a,
		"bindToResource",
		[]interface{}{_resource, _options},
	)
}

// Properties to initialize a new AssetImage.
// Experimental.
type AssetImageCodeProps struct {
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
	// Specify or override the CMD on the specified Docker image or Dockerfile.
	//
	// This needs to be in the 'exec form', viz., `[ 'executable', 'param1', 'param2' ]`.
	// See: https://docs.docker.com/engine/reference/builder/#cmd
	//
	// Experimental.
	Cmd *[]*string `json:"cmd"`
	// Specify or override the ENTRYPOINT on the specified Docker image or Dockerfile.
	//
	// An ENTRYPOINT allows you to configure a container that will run as an executable.
	// This needs to be in the 'exec form', viz., `[ 'executable', 'param1', 'param2' ]`.
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	Entrypoint *[]*string `json:"entrypoint"`
}

// Properties for enabling Lambda autoscaling.
// Experimental.
type AutoScalingOptions struct {
	// Maximum capacity to scale to.
	// Experimental.
	MaxCapacity *float64 `json:"maxCapacity"`
	// Minimum capacity to scale to.
	// Experimental.
	MinCapacity *float64 `json:"minCapacity"`
}

// A CloudFormation `AWS::Lambda::Alias`.
type CfnAlias interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	FunctionName() *string
	SetFunctionName(val *string)
	FunctionVersion() *string
	SetFunctionVersion(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	ProvisionedConcurrencyConfig() interface{}
	SetProvisionedConcurrencyConfig(val interface{})
	Ref() *string
	RoutingConfig() interface{}
	SetRoutingConfig(val interface{})
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

// The jsii proxy struct for CfnAlias
type jsiiProxy_CfnAlias struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
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

func (j *jsiiProxy_CfnAlias) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAlias) FunctionVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionVersion",
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

func (j *jsiiProxy_CfnAlias) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
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

func (j *jsiiProxy_CfnAlias) ProvisionedConcurrencyConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"provisionedConcurrencyConfig",
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

func (j *jsiiProxy_CfnAlias) RoutingConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"routingConfig",
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

func (j *jsiiProxy_CfnAlias) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::Alias`.
func NewCfnAlias(scope awscdk.Construct, id *string, props *CfnAliasProps) CfnAlias {
	_init_.Initialize()

	j := jsiiProxy_CfnAlias{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnAlias",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::Alias`.
func NewCfnAlias_Override(c CfnAlias, scope awscdk.Construct, id *string, props *CfnAliasProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnAlias",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAlias) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnAlias) SetFunctionName(val *string) {
	_jsii_.Set(
		j,
		"functionName",
		val,
	)
}

func (j *jsiiProxy_CfnAlias) SetFunctionVersion(val *string) {
	_jsii_.Set(
		j,
		"functionVersion",
		val,
	)
}

func (j *jsiiProxy_CfnAlias) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnAlias) SetProvisionedConcurrencyConfig(val interface{}) {
	_jsii_.Set(
		j,
		"provisionedConcurrencyConfig",
		val,
	)
}

func (j *jsiiProxy_CfnAlias) SetRoutingConfig(val interface{}) {
	_jsii_.Set(
		j,
		"routingConfig",
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
		"monocdk.aws_lambda.CfnAlias",
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
		"monocdk.aws_lambda.CfnAlias",
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
		"monocdk.aws_lambda.CfnAlias",
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
		"monocdk.aws_lambda.CfnAlias",
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

type CfnAlias_AliasRoutingConfigurationProperty struct {
	// `CfnAlias.AliasRoutingConfigurationProperty.AdditionalVersionWeights`.
	AdditionalVersionWeights interface{} `json:"additionalVersionWeights"`
}

type CfnAlias_ProvisionedConcurrencyConfigurationProperty struct {
	// `CfnAlias.ProvisionedConcurrencyConfigurationProperty.ProvisionedConcurrentExecutions`.
	ProvisionedConcurrentExecutions *float64 `json:"provisionedConcurrentExecutions"`
}

type CfnAlias_VersionWeightProperty struct {
	// `CfnAlias.VersionWeightProperty.FunctionVersion`.
	FunctionVersion *string `json:"functionVersion"`
	// `CfnAlias.VersionWeightProperty.FunctionWeight`.
	FunctionWeight *float64 `json:"functionWeight"`
}

// Properties for defining a `AWS::Lambda::Alias`.
type CfnAliasProps struct {
	// `AWS::Lambda::Alias.FunctionName`.
	FunctionName *string `json:"functionName"`
	// `AWS::Lambda::Alias.FunctionVersion`.
	FunctionVersion *string `json:"functionVersion"`
	// `AWS::Lambda::Alias.Name`.
	Name *string `json:"name"`
	// `AWS::Lambda::Alias.Description`.
	Description *string `json:"description"`
	// `AWS::Lambda::Alias.ProvisionedConcurrencyConfig`.
	ProvisionedConcurrencyConfig interface{} `json:"provisionedConcurrencyConfig"`
	// `AWS::Lambda::Alias.RoutingConfig`.
	RoutingConfig interface{} `json:"routingConfig"`
}

// A CloudFormation `AWS::Lambda::CodeSigningConfig`.
type CfnCodeSigningConfig interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AllowedPublishers() interface{}
	SetAllowedPublishers(val interface{})
	AttrCodeSigningConfigArn() *string
	AttrCodeSigningConfigId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CodeSigningPolicies() interface{}
	SetCodeSigningPolicies(val interface{})
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
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

// The jsii proxy struct for CfnCodeSigningConfig
type jsiiProxy_CfnCodeSigningConfig struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnCodeSigningConfig) AllowedPublishers() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"allowedPublishers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) AttrCodeSigningConfigArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCodeSigningConfigArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) AttrCodeSigningConfigId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCodeSigningConfigId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) CodeSigningPolicies() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"codeSigningPolicies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeSigningConfig) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::CodeSigningConfig`.
func NewCfnCodeSigningConfig(scope awscdk.Construct, id *string, props *CfnCodeSigningConfigProps) CfnCodeSigningConfig {
	_init_.Initialize()

	j := jsiiProxy_CfnCodeSigningConfig{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnCodeSigningConfig",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::CodeSigningConfig`.
func NewCfnCodeSigningConfig_Override(c CfnCodeSigningConfig, scope awscdk.Construct, id *string, props *CfnCodeSigningConfigProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnCodeSigningConfig",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCodeSigningConfig) SetAllowedPublishers(val interface{}) {
	_jsii_.Set(
		j,
		"allowedPublishers",
		val,
	)
}

func (j *jsiiProxy_CfnCodeSigningConfig) SetCodeSigningPolicies(val interface{}) {
	_jsii_.Set(
		j,
		"codeSigningPolicies",
		val,
	)
}

func (j *jsiiProxy_CfnCodeSigningConfig) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
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
func CfnCodeSigningConfig_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnCodeSigningConfig",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnCodeSigningConfig_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnCodeSigningConfig",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCodeSigningConfig_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnCodeSigningConfig",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnCodeSigningConfig_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnCodeSigningConfig",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnCodeSigningConfig) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnCodeSigningConfig) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnCodeSigningConfig) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnCodeSigningConfig) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnCodeSigningConfig) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) OnPrepare() {
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
func (c *jsiiProxy_CfnCodeSigningConfig) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCodeSigningConfig) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnCodeSigningConfig) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnCodeSigningConfig) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnCodeSigningConfig) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCodeSigningConfig) ToString() *string {
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
func (c *jsiiProxy_CfnCodeSigningConfig) Validate() *[]*string {
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
func (c *jsiiProxy_CfnCodeSigningConfig) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnCodeSigningConfig_AllowedPublishersProperty struct {
	// `CfnCodeSigningConfig.AllowedPublishersProperty.SigningProfileVersionArns`.
	SigningProfileVersionArns *[]*string `json:"signingProfileVersionArns"`
}

type CfnCodeSigningConfig_CodeSigningPoliciesProperty struct {
	// `CfnCodeSigningConfig.CodeSigningPoliciesProperty.UntrustedArtifactOnDeployment`.
	UntrustedArtifactOnDeployment *string `json:"untrustedArtifactOnDeployment"`
}

// Properties for defining a `AWS::Lambda::CodeSigningConfig`.
type CfnCodeSigningConfigProps struct {
	// `AWS::Lambda::CodeSigningConfig.AllowedPublishers`.
	AllowedPublishers interface{} `json:"allowedPublishers"`
	// `AWS::Lambda::CodeSigningConfig.CodeSigningPolicies`.
	CodeSigningPolicies interface{} `json:"codeSigningPolicies"`
	// `AWS::Lambda::CodeSigningConfig.Description`.
	Description *string `json:"description"`
}

// A CloudFormation `AWS::Lambda::EventInvokeConfig`.
type CfnEventInvokeConfig interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DestinationConfig() interface{}
	SetDestinationConfig(val interface{})
	FunctionName() *string
	SetFunctionName(val *string)
	LogicalId() *string
	MaximumEventAgeInSeconds() *float64
	SetMaximumEventAgeInSeconds(val *float64)
	MaximumRetryAttempts() *float64
	SetMaximumRetryAttempts(val *float64)
	Node() awscdk.ConstructNode
	Qualifier() *string
	SetQualifier(val *string)
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

// The jsii proxy struct for CfnEventInvokeConfig
type jsiiProxy_CfnEventInvokeConfig struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnEventInvokeConfig) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) DestinationConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"destinationConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) MaximumEventAgeInSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maximumEventAgeInSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) MaximumRetryAttempts() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maximumRetryAttempts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) Qualifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"qualifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventInvokeConfig) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::EventInvokeConfig`.
func NewCfnEventInvokeConfig(scope awscdk.Construct, id *string, props *CfnEventInvokeConfigProps) CfnEventInvokeConfig {
	_init_.Initialize()

	j := jsiiProxy_CfnEventInvokeConfig{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnEventInvokeConfig",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::EventInvokeConfig`.
func NewCfnEventInvokeConfig_Override(c CfnEventInvokeConfig, scope awscdk.Construct, id *string, props *CfnEventInvokeConfigProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnEventInvokeConfig",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnEventInvokeConfig) SetDestinationConfig(val interface{}) {
	_jsii_.Set(
		j,
		"destinationConfig",
		val,
	)
}

func (j *jsiiProxy_CfnEventInvokeConfig) SetFunctionName(val *string) {
	_jsii_.Set(
		j,
		"functionName",
		val,
	)
}

func (j *jsiiProxy_CfnEventInvokeConfig) SetMaximumEventAgeInSeconds(val *float64) {
	_jsii_.Set(
		j,
		"maximumEventAgeInSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnEventInvokeConfig) SetMaximumRetryAttempts(val *float64) {
	_jsii_.Set(
		j,
		"maximumRetryAttempts",
		val,
	)
}

func (j *jsiiProxy_CfnEventInvokeConfig) SetQualifier(val *string) {
	_jsii_.Set(
		j,
		"qualifier",
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
func CfnEventInvokeConfig_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnEventInvokeConfig",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnEventInvokeConfig_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnEventInvokeConfig",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnEventInvokeConfig_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnEventInvokeConfig",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnEventInvokeConfig_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnEventInvokeConfig",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnEventInvokeConfig) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnEventInvokeConfig) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnEventInvokeConfig) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnEventInvokeConfig) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnEventInvokeConfig) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) OnPrepare() {
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
func (c *jsiiProxy_CfnEventInvokeConfig) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnEventInvokeConfig) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnEventInvokeConfig) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnEventInvokeConfig) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnEventInvokeConfig) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventInvokeConfig) ToString() *string {
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
func (c *jsiiProxy_CfnEventInvokeConfig) Validate() *[]*string {
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
func (c *jsiiProxy_CfnEventInvokeConfig) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnEventInvokeConfig_DestinationConfigProperty struct {
	// `CfnEventInvokeConfig.DestinationConfigProperty.OnFailure`.
	OnFailure interface{} `json:"onFailure"`
	// `CfnEventInvokeConfig.DestinationConfigProperty.OnSuccess`.
	OnSuccess interface{} `json:"onSuccess"`
}

type CfnEventInvokeConfig_OnFailureProperty struct {
	// `CfnEventInvokeConfig.OnFailureProperty.Destination`.
	Destination *string `json:"destination"`
}

type CfnEventInvokeConfig_OnSuccessProperty struct {
	// `CfnEventInvokeConfig.OnSuccessProperty.Destination`.
	Destination *string `json:"destination"`
}

// Properties for defining a `AWS::Lambda::EventInvokeConfig`.
type CfnEventInvokeConfigProps struct {
	// `AWS::Lambda::EventInvokeConfig.FunctionName`.
	FunctionName *string `json:"functionName"`
	// `AWS::Lambda::EventInvokeConfig.Qualifier`.
	Qualifier *string `json:"qualifier"`
	// `AWS::Lambda::EventInvokeConfig.DestinationConfig`.
	DestinationConfig interface{} `json:"destinationConfig"`
	// `AWS::Lambda::EventInvokeConfig.MaximumEventAgeInSeconds`.
	MaximumEventAgeInSeconds *float64 `json:"maximumEventAgeInSeconds"`
	// `AWS::Lambda::EventInvokeConfig.MaximumRetryAttempts`.
	MaximumRetryAttempts *float64 `json:"maximumRetryAttempts"`
}

// A CloudFormation `AWS::Lambda::EventSourceMapping`.
type CfnEventSourceMapping interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrId() *string
	BatchSize() *float64
	SetBatchSize(val *float64)
	BisectBatchOnFunctionError() interface{}
	SetBisectBatchOnFunctionError(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DestinationConfig() interface{}
	SetDestinationConfig(val interface{})
	Enabled() interface{}
	SetEnabled(val interface{})
	EventSourceArn() *string
	SetEventSourceArn(val *string)
	FunctionName() *string
	SetFunctionName(val *string)
	FunctionResponseTypes() *[]*string
	SetFunctionResponseTypes(val *[]*string)
	LogicalId() *string
	MaximumBatchingWindowInSeconds() *float64
	SetMaximumBatchingWindowInSeconds(val *float64)
	MaximumRecordAgeInSeconds() *float64
	SetMaximumRecordAgeInSeconds(val *float64)
	MaximumRetryAttempts() *float64
	SetMaximumRetryAttempts(val *float64)
	Node() awscdk.ConstructNode
	ParallelizationFactor() *float64
	SetParallelizationFactor(val *float64)
	Queues() *[]*string
	SetQueues(val *[]*string)
	Ref() *string
	SelfManagedEventSource() interface{}
	SetSelfManagedEventSource(val interface{})
	SourceAccessConfigurations() interface{}
	SetSourceAccessConfigurations(val interface{})
	Stack() awscdk.Stack
	StartingPosition() *string
	SetStartingPosition(val *string)
	Topics() *[]*string
	SetTopics(val *[]*string)
	TumblingWindowInSeconds() *float64
	SetTumblingWindowInSeconds(val *float64)
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

// The jsii proxy struct for CfnEventSourceMapping
type jsiiProxy_CfnEventSourceMapping struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnEventSourceMapping) AttrId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) BatchSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"batchSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) BisectBatchOnFunctionError() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"bisectBatchOnFunctionError",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) DestinationConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"destinationConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) Enabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"enabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) EventSourceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventSourceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) FunctionResponseTypes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"functionResponseTypes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) MaximumBatchingWindowInSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maximumBatchingWindowInSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) MaximumRecordAgeInSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maximumRecordAgeInSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) MaximumRetryAttempts() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maximumRetryAttempts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) ParallelizationFactor() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"parallelizationFactor",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) Queues() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"queues",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) SelfManagedEventSource() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"selfManagedEventSource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) SourceAccessConfigurations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"sourceAccessConfigurations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) StartingPosition() *string {
	var returns *string
	_jsii_.Get(
		j,
		"startingPosition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) Topics() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"topics",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) TumblingWindowInSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"tumblingWindowInSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnEventSourceMapping) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::EventSourceMapping`.
func NewCfnEventSourceMapping(scope awscdk.Construct, id *string, props *CfnEventSourceMappingProps) CfnEventSourceMapping {
	_init_.Initialize()

	j := jsiiProxy_CfnEventSourceMapping{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnEventSourceMapping",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::EventSourceMapping`.
func NewCfnEventSourceMapping_Override(c CfnEventSourceMapping, scope awscdk.Construct, id *string, props *CfnEventSourceMappingProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnEventSourceMapping",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetBatchSize(val *float64) {
	_jsii_.Set(
		j,
		"batchSize",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetBisectBatchOnFunctionError(val interface{}) {
	_jsii_.Set(
		j,
		"bisectBatchOnFunctionError",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetDestinationConfig(val interface{}) {
	_jsii_.Set(
		j,
		"destinationConfig",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"enabled",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetEventSourceArn(val *string) {
	_jsii_.Set(
		j,
		"eventSourceArn",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetFunctionName(val *string) {
	_jsii_.Set(
		j,
		"functionName",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetFunctionResponseTypes(val *[]*string) {
	_jsii_.Set(
		j,
		"functionResponseTypes",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetMaximumBatchingWindowInSeconds(val *float64) {
	_jsii_.Set(
		j,
		"maximumBatchingWindowInSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetMaximumRecordAgeInSeconds(val *float64) {
	_jsii_.Set(
		j,
		"maximumRecordAgeInSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetMaximumRetryAttempts(val *float64) {
	_jsii_.Set(
		j,
		"maximumRetryAttempts",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetParallelizationFactor(val *float64) {
	_jsii_.Set(
		j,
		"parallelizationFactor",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetQueues(val *[]*string) {
	_jsii_.Set(
		j,
		"queues",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetSelfManagedEventSource(val interface{}) {
	_jsii_.Set(
		j,
		"selfManagedEventSource",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetSourceAccessConfigurations(val interface{}) {
	_jsii_.Set(
		j,
		"sourceAccessConfigurations",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetStartingPosition(val *string) {
	_jsii_.Set(
		j,
		"startingPosition",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetTopics(val *[]*string) {
	_jsii_.Set(
		j,
		"topics",
		val,
	)
}

func (j *jsiiProxy_CfnEventSourceMapping) SetTumblingWindowInSeconds(val *float64) {
	_jsii_.Set(
		j,
		"tumblingWindowInSeconds",
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
func CfnEventSourceMapping_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnEventSourceMapping",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnEventSourceMapping_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnEventSourceMapping",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnEventSourceMapping_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnEventSourceMapping",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnEventSourceMapping_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnEventSourceMapping",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnEventSourceMapping) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnEventSourceMapping) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnEventSourceMapping) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnEventSourceMapping) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnEventSourceMapping) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnEventSourceMapping) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnEventSourceMapping) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnEventSourceMapping) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnEventSourceMapping) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnEventSourceMapping) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnEventSourceMapping) OnPrepare() {
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
func (c *jsiiProxy_CfnEventSourceMapping) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventSourceMapping) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnEventSourceMapping) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnEventSourceMapping) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnEventSourceMapping) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnEventSourceMapping) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnEventSourceMapping) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnEventSourceMapping) ToString() *string {
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
func (c *jsiiProxy_CfnEventSourceMapping) Validate() *[]*string {
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
func (c *jsiiProxy_CfnEventSourceMapping) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnEventSourceMapping_DestinationConfigProperty struct {
	// `CfnEventSourceMapping.DestinationConfigProperty.OnFailure`.
	OnFailure interface{} `json:"onFailure"`
}

type CfnEventSourceMapping_EndpointsProperty struct {
	// `CfnEventSourceMapping.EndpointsProperty.KafkaBootstrapServers`.
	KafkaBootstrapServers *[]*string `json:"kafkaBootstrapServers"`
}

type CfnEventSourceMapping_OnFailureProperty struct {
	// `CfnEventSourceMapping.OnFailureProperty.Destination`.
	Destination *string `json:"destination"`
}

type CfnEventSourceMapping_SelfManagedEventSourceProperty struct {
	// `CfnEventSourceMapping.SelfManagedEventSourceProperty.Endpoints`.
	Endpoints interface{} `json:"endpoints"`
}

type CfnEventSourceMapping_SourceAccessConfigurationProperty struct {
	// `CfnEventSourceMapping.SourceAccessConfigurationProperty.Type`.
	Type *string `json:"type"`
	// `CfnEventSourceMapping.SourceAccessConfigurationProperty.URI`.
	Uri *string `json:"uri"`
}

// Properties for defining a `AWS::Lambda::EventSourceMapping`.
type CfnEventSourceMappingProps struct {
	// `AWS::Lambda::EventSourceMapping.FunctionName`.
	FunctionName *string `json:"functionName"`
	// `AWS::Lambda::EventSourceMapping.BatchSize`.
	BatchSize *float64 `json:"batchSize"`
	// `AWS::Lambda::EventSourceMapping.BisectBatchOnFunctionError`.
	BisectBatchOnFunctionError interface{} `json:"bisectBatchOnFunctionError"`
	// `AWS::Lambda::EventSourceMapping.DestinationConfig`.
	DestinationConfig interface{} `json:"destinationConfig"`
	// `AWS::Lambda::EventSourceMapping.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `AWS::Lambda::EventSourceMapping.EventSourceArn`.
	EventSourceArn *string `json:"eventSourceArn"`
	// `AWS::Lambda::EventSourceMapping.FunctionResponseTypes`.
	FunctionResponseTypes *[]*string `json:"functionResponseTypes"`
	// `AWS::Lambda::EventSourceMapping.MaximumBatchingWindowInSeconds`.
	MaximumBatchingWindowInSeconds *float64 `json:"maximumBatchingWindowInSeconds"`
	// `AWS::Lambda::EventSourceMapping.MaximumRecordAgeInSeconds`.
	MaximumRecordAgeInSeconds *float64 `json:"maximumRecordAgeInSeconds"`
	// `AWS::Lambda::EventSourceMapping.MaximumRetryAttempts`.
	MaximumRetryAttempts *float64 `json:"maximumRetryAttempts"`
	// `AWS::Lambda::EventSourceMapping.ParallelizationFactor`.
	ParallelizationFactor *float64 `json:"parallelizationFactor"`
	// `AWS::Lambda::EventSourceMapping.Queues`.
	Queues *[]*string `json:"queues"`
	// `AWS::Lambda::EventSourceMapping.SelfManagedEventSource`.
	SelfManagedEventSource interface{} `json:"selfManagedEventSource"`
	// `AWS::Lambda::EventSourceMapping.SourceAccessConfigurations`.
	SourceAccessConfigurations interface{} `json:"sourceAccessConfigurations"`
	// `AWS::Lambda::EventSourceMapping.StartingPosition`.
	StartingPosition *string `json:"startingPosition"`
	// `AWS::Lambda::EventSourceMapping.Topics`.
	Topics *[]*string `json:"topics"`
	// `AWS::Lambda::EventSourceMapping.TumblingWindowInSeconds`.
	TumblingWindowInSeconds *float64 `json:"tumblingWindowInSeconds"`
}

// A CloudFormation `AWS::Lambda::Function`.
type CfnFunction interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Code() interface{}
	SetCode(val interface{})
	CodeSigningConfigArn() *string
	SetCodeSigningConfigArn(val *string)
	CreationStack() *[]*string
	DeadLetterConfig() interface{}
	SetDeadLetterConfig(val interface{})
	Description() *string
	SetDescription(val *string)
	Environment() interface{}
	SetEnvironment(val interface{})
	FileSystemConfigs() interface{}
	SetFileSystemConfigs(val interface{})
	FunctionName() *string
	SetFunctionName(val *string)
	Handler() *string
	SetHandler(val *string)
	ImageConfig() interface{}
	SetImageConfig(val interface{})
	KmsKeyArn() *string
	SetKmsKeyArn(val *string)
	Layers() *[]*string
	SetLayers(val *[]*string)
	LogicalId() *string
	MemorySize() *float64
	SetMemorySize(val *float64)
	Node() awscdk.ConstructNode
	PackageType() *string
	SetPackageType(val *string)
	Ref() *string
	ReservedConcurrentExecutions() *float64
	SetReservedConcurrentExecutions(val *float64)
	Role() *string
	SetRole(val *string)
	Runtime() *string
	SetRuntime(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	Timeout() *float64
	SetTimeout(val *float64)
	TracingConfig() interface{}
	SetTracingConfig(val interface{})
	UpdatedProperites() *map[string]interface{}
	VpcConfig() interface{}
	SetVpcConfig(val interface{})
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

// The jsii proxy struct for CfnFunction
type jsiiProxy_CfnFunction struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnFunction) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Code() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"code",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) CodeSigningConfigArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codeSigningConfigArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) DeadLetterConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deadLetterConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Environment() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"environment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) FileSystemConfigs() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"fileSystemConfigs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Handler() *string {
	var returns *string
	_jsii_.Get(
		j,
		"handler",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) ImageConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"imageConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) KmsKeyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kmsKeyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Layers() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"layers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) MemorySize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memorySize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) PackageType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"packageType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) ReservedConcurrentExecutions() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"reservedConcurrentExecutions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Role() *string {
	var returns *string
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Runtime() *string {
	var returns *string
	_jsii_.Get(
		j,
		"runtime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) Timeout() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"timeout",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) TracingConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"tracingConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnFunction) VpcConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"vpcConfig",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::Function`.
func NewCfnFunction(scope awscdk.Construct, id *string, props *CfnFunctionProps) CfnFunction {
	_init_.Initialize()

	j := jsiiProxy_CfnFunction{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnFunction",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::Function`.
func NewCfnFunction_Override(c CfnFunction, scope awscdk.Construct, id *string, props *CfnFunctionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnFunction",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnFunction) SetCode(val interface{}) {
	_jsii_.Set(
		j,
		"code",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetCodeSigningConfigArn(val *string) {
	_jsii_.Set(
		j,
		"codeSigningConfigArn",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetDeadLetterConfig(val interface{}) {
	_jsii_.Set(
		j,
		"deadLetterConfig",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetEnvironment(val interface{}) {
	_jsii_.Set(
		j,
		"environment",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetFileSystemConfigs(val interface{}) {
	_jsii_.Set(
		j,
		"fileSystemConfigs",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetFunctionName(val *string) {
	_jsii_.Set(
		j,
		"functionName",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetHandler(val *string) {
	_jsii_.Set(
		j,
		"handler",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetImageConfig(val interface{}) {
	_jsii_.Set(
		j,
		"imageConfig",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetKmsKeyArn(val *string) {
	_jsii_.Set(
		j,
		"kmsKeyArn",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetLayers(val *[]*string) {
	_jsii_.Set(
		j,
		"layers",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetMemorySize(val *float64) {
	_jsii_.Set(
		j,
		"memorySize",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetPackageType(val *string) {
	_jsii_.Set(
		j,
		"packageType",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetReservedConcurrentExecutions(val *float64) {
	_jsii_.Set(
		j,
		"reservedConcurrentExecutions",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetRole(val *string) {
	_jsii_.Set(
		j,
		"role",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetRuntime(val *string) {
	_jsii_.Set(
		j,
		"runtime",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetTimeout(val *float64) {
	_jsii_.Set(
		j,
		"timeout",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetTracingConfig(val interface{}) {
	_jsii_.Set(
		j,
		"tracingConfig",
		val,
	)
}

func (j *jsiiProxy_CfnFunction) SetVpcConfig(val interface{}) {
	_jsii_.Set(
		j,
		"vpcConfig",
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
func CfnFunction_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnFunction",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnFunction_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnFunction",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnFunction_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnFunction",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnFunction_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnFunction",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnFunction) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnFunction) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnFunction) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnFunction) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnFunction) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnFunction) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnFunction) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnFunction) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnFunction) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnFunction) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnFunction) OnPrepare() {
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
func (c *jsiiProxy_CfnFunction) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnFunction) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnFunction) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnFunction) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnFunction) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnFunction) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnFunction) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnFunction) ToString() *string {
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
func (c *jsiiProxy_CfnFunction) Validate() *[]*string {
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
func (c *jsiiProxy_CfnFunction) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnFunction_CodeProperty struct {
	// `CfnFunction.CodeProperty.ImageUri`.
	ImageUri *string `json:"imageUri"`
	// `CfnFunction.CodeProperty.S3Bucket`.
	S3Bucket *string `json:"s3Bucket"`
	// `CfnFunction.CodeProperty.S3Key`.
	S3Key *string `json:"s3Key"`
	// `CfnFunction.CodeProperty.S3ObjectVersion`.
	S3ObjectVersion *string `json:"s3ObjectVersion"`
	// `CfnFunction.CodeProperty.ZipFile`.
	ZipFile *string `json:"zipFile"`
}

type CfnFunction_DeadLetterConfigProperty struct {
	// `CfnFunction.DeadLetterConfigProperty.TargetArn`.
	TargetArn *string `json:"targetArn"`
}

type CfnFunction_EnvironmentProperty struct {
	// `CfnFunction.EnvironmentProperty.Variables`.
	Variables interface{} `json:"variables"`
}

type CfnFunction_FileSystemConfigProperty struct {
	// `CfnFunction.FileSystemConfigProperty.Arn`.
	Arn *string `json:"arn"`
	// `CfnFunction.FileSystemConfigProperty.LocalMountPath`.
	LocalMountPath *string `json:"localMountPath"`
}

type CfnFunction_ImageConfigProperty struct {
	// `CfnFunction.ImageConfigProperty.Command`.
	Command *[]*string `json:"command"`
	// `CfnFunction.ImageConfigProperty.EntryPoint`.
	EntryPoint *[]*string `json:"entryPoint"`
	// `CfnFunction.ImageConfigProperty.WorkingDirectory`.
	WorkingDirectory *string `json:"workingDirectory"`
}

type CfnFunction_TracingConfigProperty struct {
	// `CfnFunction.TracingConfigProperty.Mode`.
	Mode *string `json:"mode"`
}

type CfnFunction_VpcConfigProperty struct {
	// `CfnFunction.VpcConfigProperty.SecurityGroupIds`.
	SecurityGroupIds *[]*string `json:"securityGroupIds"`
	// `CfnFunction.VpcConfigProperty.SubnetIds`.
	SubnetIds *[]*string `json:"subnetIds"`
}

// Properties for defining a `AWS::Lambda::Function`.
type CfnFunctionProps struct {
	// `AWS::Lambda::Function.Code`.
	Code interface{} `json:"code"`
	// `AWS::Lambda::Function.Role`.
	Role *string `json:"role"`
	// `AWS::Lambda::Function.CodeSigningConfigArn`.
	CodeSigningConfigArn *string `json:"codeSigningConfigArn"`
	// `AWS::Lambda::Function.DeadLetterConfig`.
	DeadLetterConfig interface{} `json:"deadLetterConfig"`
	// `AWS::Lambda::Function.Description`.
	Description *string `json:"description"`
	// `AWS::Lambda::Function.Environment`.
	Environment interface{} `json:"environment"`
	// `AWS::Lambda::Function.FileSystemConfigs`.
	FileSystemConfigs interface{} `json:"fileSystemConfigs"`
	// `AWS::Lambda::Function.FunctionName`.
	FunctionName *string `json:"functionName"`
	// `AWS::Lambda::Function.Handler`.
	Handler *string `json:"handler"`
	// `AWS::Lambda::Function.ImageConfig`.
	ImageConfig interface{} `json:"imageConfig"`
	// `AWS::Lambda::Function.KmsKeyArn`.
	KmsKeyArn *string `json:"kmsKeyArn"`
	// `AWS::Lambda::Function.Layers`.
	Layers *[]*string `json:"layers"`
	// `AWS::Lambda::Function.MemorySize`.
	MemorySize *float64 `json:"memorySize"`
	// `AWS::Lambda::Function.PackageType`.
	PackageType *string `json:"packageType"`
	// `AWS::Lambda::Function.ReservedConcurrentExecutions`.
	ReservedConcurrentExecutions *float64 `json:"reservedConcurrentExecutions"`
	// `AWS::Lambda::Function.Runtime`.
	Runtime *string `json:"runtime"`
	// `AWS::Lambda::Function.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::Lambda::Function.Timeout`.
	Timeout *float64 `json:"timeout"`
	// `AWS::Lambda::Function.TracingConfig`.
	TracingConfig interface{} `json:"tracingConfig"`
	// `AWS::Lambda::Function.VpcConfig`.
	VpcConfig interface{} `json:"vpcConfig"`
}

// A CloudFormation `AWS::Lambda::LayerVersion`.
type CfnLayerVersion interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CompatibleRuntimes() *[]*string
	SetCompatibleRuntimes(val *[]*string)
	Content() interface{}
	SetContent(val interface{})
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LayerName() *string
	SetLayerName(val *string)
	LicenseInfo() *string
	SetLicenseInfo(val *string)
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

// The jsii proxy struct for CfnLayerVersion
type jsiiProxy_CfnLayerVersion struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnLayerVersion) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) CompatibleRuntimes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"compatibleRuntimes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) Content() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"content",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) LayerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"layerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) LicenseInfo() *string {
	var returns *string
	_jsii_.Get(
		j,
		"licenseInfo",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersion) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::LayerVersion`.
func NewCfnLayerVersion(scope awscdk.Construct, id *string, props *CfnLayerVersionProps) CfnLayerVersion {
	_init_.Initialize()

	j := jsiiProxy_CfnLayerVersion{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnLayerVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::LayerVersion`.
func NewCfnLayerVersion_Override(c CfnLayerVersion, scope awscdk.Construct, id *string, props *CfnLayerVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnLayerVersion",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnLayerVersion) SetCompatibleRuntimes(val *[]*string) {
	_jsii_.Set(
		j,
		"compatibleRuntimes",
		val,
	)
}

func (j *jsiiProxy_CfnLayerVersion) SetContent(val interface{}) {
	_jsii_.Set(
		j,
		"content",
		val,
	)
}

func (j *jsiiProxy_CfnLayerVersion) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnLayerVersion) SetLayerName(val *string) {
	_jsii_.Set(
		j,
		"layerName",
		val,
	)
}

func (j *jsiiProxy_CfnLayerVersion) SetLicenseInfo(val *string) {
	_jsii_.Set(
		j,
		"licenseInfo",
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
func CfnLayerVersion_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnLayerVersion",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnLayerVersion_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnLayerVersion",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnLayerVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnLayerVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnLayerVersion_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnLayerVersion",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnLayerVersion) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnLayerVersion) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnLayerVersion) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnLayerVersion) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnLayerVersion) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnLayerVersion) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnLayerVersion) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnLayerVersion) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnLayerVersion) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnLayerVersion) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnLayerVersion) OnPrepare() {
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
func (c *jsiiProxy_CfnLayerVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLayerVersion) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnLayerVersion) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnLayerVersion) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnLayerVersion) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnLayerVersion) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnLayerVersion) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLayerVersion) ToString() *string {
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
func (c *jsiiProxy_CfnLayerVersion) Validate() *[]*string {
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
func (c *jsiiProxy_CfnLayerVersion) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnLayerVersion_ContentProperty struct {
	// `CfnLayerVersion.ContentProperty.S3Bucket`.
	S3Bucket *string `json:"s3Bucket"`
	// `CfnLayerVersion.ContentProperty.S3Key`.
	S3Key *string `json:"s3Key"`
	// `CfnLayerVersion.ContentProperty.S3ObjectVersion`.
	S3ObjectVersion *string `json:"s3ObjectVersion"`
}

// A CloudFormation `AWS::Lambda::LayerVersionPermission`.
type CfnLayerVersionPermission interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Action() *string
	SetAction(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LayerVersionArn() *string
	SetLayerVersionArn(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	OrganizationId() *string
	SetOrganizationId(val *string)
	Principal() *string
	SetPrincipal(val *string)
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

// The jsii proxy struct for CfnLayerVersionPermission
type jsiiProxy_CfnLayerVersionPermission struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnLayerVersionPermission) Action() *string {
	var returns *string
	_jsii_.Get(
		j,
		"action",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) LayerVersionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"layerVersionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) OrganizationId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"organizationId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) Principal() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLayerVersionPermission) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::LayerVersionPermission`.
func NewCfnLayerVersionPermission(scope awscdk.Construct, id *string, props *CfnLayerVersionPermissionProps) CfnLayerVersionPermission {
	_init_.Initialize()

	j := jsiiProxy_CfnLayerVersionPermission{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnLayerVersionPermission",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::LayerVersionPermission`.
func NewCfnLayerVersionPermission_Override(c CfnLayerVersionPermission, scope awscdk.Construct, id *string, props *CfnLayerVersionPermissionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnLayerVersionPermission",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnLayerVersionPermission) SetAction(val *string) {
	_jsii_.Set(
		j,
		"action",
		val,
	)
}

func (j *jsiiProxy_CfnLayerVersionPermission) SetLayerVersionArn(val *string) {
	_jsii_.Set(
		j,
		"layerVersionArn",
		val,
	)
}

func (j *jsiiProxy_CfnLayerVersionPermission) SetOrganizationId(val *string) {
	_jsii_.Set(
		j,
		"organizationId",
		val,
	)
}

func (j *jsiiProxy_CfnLayerVersionPermission) SetPrincipal(val *string) {
	_jsii_.Set(
		j,
		"principal",
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
func CfnLayerVersionPermission_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnLayerVersionPermission",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnLayerVersionPermission_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnLayerVersionPermission",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnLayerVersionPermission_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnLayerVersionPermission",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnLayerVersionPermission_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnLayerVersionPermission",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnLayerVersionPermission) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnLayerVersionPermission) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnLayerVersionPermission) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnLayerVersionPermission) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnLayerVersionPermission) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) OnPrepare() {
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
func (c *jsiiProxy_CfnLayerVersionPermission) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnLayerVersionPermission) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnLayerVersionPermission) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnLayerVersionPermission) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnLayerVersionPermission) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLayerVersionPermission) ToString() *string {
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
func (c *jsiiProxy_CfnLayerVersionPermission) Validate() *[]*string {
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
func (c *jsiiProxy_CfnLayerVersionPermission) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Lambda::LayerVersionPermission`.
type CfnLayerVersionPermissionProps struct {
	// `AWS::Lambda::LayerVersionPermission.Action`.
	Action *string `json:"action"`
	// `AWS::Lambda::LayerVersionPermission.LayerVersionArn`.
	LayerVersionArn *string `json:"layerVersionArn"`
	// `AWS::Lambda::LayerVersionPermission.Principal`.
	Principal *string `json:"principal"`
	// `AWS::Lambda::LayerVersionPermission.OrganizationId`.
	OrganizationId *string `json:"organizationId"`
}

// Properties for defining a `AWS::Lambda::LayerVersion`.
type CfnLayerVersionProps struct {
	// `AWS::Lambda::LayerVersion.Content`.
	Content interface{} `json:"content"`
	// `AWS::Lambda::LayerVersion.CompatibleRuntimes`.
	CompatibleRuntimes *[]*string `json:"compatibleRuntimes"`
	// `AWS::Lambda::LayerVersion.Description`.
	Description *string `json:"description"`
	// `AWS::Lambda::LayerVersion.LayerName`.
	LayerName *string `json:"layerName"`
	// `AWS::Lambda::LayerVersion.LicenseInfo`.
	LicenseInfo *string `json:"licenseInfo"`
}

// Lambda code defined using 2 CloudFormation parameters.
//
// Useful when you don't have access to the code of your Lambda from your CDK code, so you can't use Assets,
// and you want to deploy the Lambda in a CodePipeline, using CloudFormation Actions -
// you can fill the parameters using the {@link #assign} method.
// Experimental.
type CfnParametersCode interface {
	Code
	BucketNameParam() *string
	IsInline() *bool
	ObjectKeyParam() *string
	Assign(location *awss3.Location) *map[string]interface{}
	Bind(scope awscdk.Construct) *CodeConfig
	BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions)
}

// The jsii proxy struct for CfnParametersCode
type jsiiProxy_CfnParametersCode struct {
	jsiiProxy_Code
}

func (j *jsiiProxy_CfnParametersCode) BucketNameParam() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketNameParam",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParametersCode) IsInline() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isInline",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParametersCode) ObjectKeyParam() *string {
	var returns *string
	_jsii_.Get(
		j,
		"objectKeyParam",
		&returns,
	)
	return returns
}


// Experimental.
func NewCfnParametersCode(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	j := jsiiProxy_CfnParametersCode{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnParametersCode",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewCfnParametersCode_Override(c CfnParametersCode, props *CfnParametersCodeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnParametersCode",
		[]interface{}{props},
		c,
	)
}

// DEPRECATED.
// Deprecated: use `fromAsset`
func CfnParametersCode_Asset(path *string) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"asset",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromBucket`
func CfnParametersCode_Bucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"bucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromCfnParameters`
func CfnParametersCode_CfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"cfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from a local disk path.
// Experimental.
func CfnParametersCode_FromAsset(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func CfnParametersCode_FromAssetImage(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	var returns AssetImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"fromAssetImage",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Lambda handler code as an S3 object.
// Experimental.
func CfnParametersCode_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Creates a new Lambda source defined using CloudFormation parameters.
//
// Returns: a new instance of `CfnParametersCode`
// Experimental.
func CfnParametersCode_FromCfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"fromCfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from an asset created by a Docker build.
//
// By default, the asset is expected to be located at `/asset` in the
// image.
// Experimental.
func CfnParametersCode_FromDockerBuild(path *string, options *DockerBuildAssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"fromDockerBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func CfnParametersCode_FromEcrImage(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	var returns EcrImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"fromEcrImage",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Inline code for Lambda handler.
//
// Returns: `LambdaInlineCode` with inline code.
// Experimental.
func CfnParametersCode_FromInline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"fromInline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromInline`
func CfnParametersCode_Inline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnParametersCode",
		"inline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// Create a parameters map from this instance's CloudFormation parameters.
//
// It returns a map with 2 keys that correspond to the names of the parameters defined in this Lambda code,
// and as values it contains the appropriate expressions pointing at the provided S3 location
// (most likely, obtained from a CodePipeline Artifact by calling the `artifact.s3Location` method).
// The result should be provided to the CloudFormation Action
// that is deploying the Stack that the Lambda with this code is part of,
// in the `parameterOverrides` property.
// Experimental.
func (c *jsiiProxy_CfnParametersCode) Assign(location *awss3.Location) *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"assign",
		[]interface{}{location},
		&returns,
	)

	return returns
}

// Called when the lambda or layer is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (c *jsiiProxy_CfnParametersCode) Bind(scope awscdk.Construct) *CodeConfig {
	var returns *CodeConfig

	_jsii_.Invoke(
		c,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Called after the CFN function resource has been created to allow the code class to bind to it.
//
// Specifically it's required to allow assets to add
// metadata for tooling like SAM CLI to be able to find their origins.
// Experimental.
func (c *jsiiProxy_CfnParametersCode) BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions) {
	_jsii_.InvokeVoid(
		c,
		"bindToResource",
		[]interface{}{_resource, _options},
	)
}

// Construction properties for {@link CfnParametersCode}.
// Experimental.
type CfnParametersCodeProps struct {
	// The CloudFormation parameter that represents the name of the S3 Bucket where the Lambda code will be located in.
	//
	// Must be of type 'String'.
	// Experimental.
	BucketNameParam awscdk.CfnParameter `json:"bucketNameParam"`
	// The CloudFormation parameter that represents the path inside the S3 Bucket where the Lambda code will be located at.
	//
	// Must be of type 'String'.
	// Experimental.
	ObjectKeyParam awscdk.CfnParameter `json:"objectKeyParam"`
}

// A CloudFormation `AWS::Lambda::Permission`.
type CfnPermission interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Action() *string
	SetAction(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	EventSourceToken() *string
	SetEventSourceToken(val *string)
	FunctionName() *string
	SetFunctionName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Principal() *string
	SetPrincipal(val *string)
	Ref() *string
	SourceAccount() *string
	SetSourceAccount(val *string)
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

// The jsii proxy struct for CfnPermission
type jsiiProxy_CfnPermission struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnPermission) Action() *string {
	var returns *string
	_jsii_.Get(
		j,
		"action",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) EventSourceToken() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventSourceToken",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) Principal() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) SourceAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceAccount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) SourceArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPermission) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::Permission`.
func NewCfnPermission(scope awscdk.Construct, id *string, props *CfnPermissionProps) CfnPermission {
	_init_.Initialize()

	j := jsiiProxy_CfnPermission{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnPermission",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::Permission`.
func NewCfnPermission_Override(c CfnPermission, scope awscdk.Construct, id *string, props *CfnPermissionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnPermission",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnPermission) SetAction(val *string) {
	_jsii_.Set(
		j,
		"action",
		val,
	)
}

func (j *jsiiProxy_CfnPermission) SetEventSourceToken(val *string) {
	_jsii_.Set(
		j,
		"eventSourceToken",
		val,
	)
}

func (j *jsiiProxy_CfnPermission) SetFunctionName(val *string) {
	_jsii_.Set(
		j,
		"functionName",
		val,
	)
}

func (j *jsiiProxy_CfnPermission) SetPrincipal(val *string) {
	_jsii_.Set(
		j,
		"principal",
		val,
	)
}

func (j *jsiiProxy_CfnPermission) SetSourceAccount(val *string) {
	_jsii_.Set(
		j,
		"sourceAccount",
		val,
	)
}

func (j *jsiiProxy_CfnPermission) SetSourceArn(val *string) {
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
func CfnPermission_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnPermission",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnPermission_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnPermission",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnPermission_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnPermission",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnPermission_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnPermission",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnPermission) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnPermission) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnPermission) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnPermission) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnPermission) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnPermission) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnPermission) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnPermission) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnPermission) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnPermission) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnPermission) OnPrepare() {
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
func (c *jsiiProxy_CfnPermission) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPermission) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnPermission) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnPermission) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnPermission) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnPermission) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnPermission) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPermission) ToString() *string {
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
func (c *jsiiProxy_CfnPermission) Validate() *[]*string {
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
func (c *jsiiProxy_CfnPermission) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Lambda::Permission`.
type CfnPermissionProps struct {
	// `AWS::Lambda::Permission.Action`.
	Action *string `json:"action"`
	// `AWS::Lambda::Permission.FunctionName`.
	FunctionName *string `json:"functionName"`
	// `AWS::Lambda::Permission.Principal`.
	Principal *string `json:"principal"`
	// `AWS::Lambda::Permission.EventSourceToken`.
	EventSourceToken *string `json:"eventSourceToken"`
	// `AWS::Lambda::Permission.SourceAccount`.
	SourceAccount *string `json:"sourceAccount"`
	// `AWS::Lambda::Permission.SourceArn`.
	SourceArn *string `json:"sourceArn"`
}

// A CloudFormation `AWS::Lambda::Version`.
type CfnVersion interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrVersion() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CodeSha256() *string
	SetCodeSha256(val *string)
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	FunctionName() *string
	SetFunctionName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	ProvisionedConcurrencyConfig() interface{}
	SetProvisionedConcurrencyConfig(val interface{})
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

// The jsii proxy struct for CfnVersion
type jsiiProxy_CfnVersion struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnVersion) AttrVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) CodeSha256() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codeSha256",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) ProvisionedConcurrencyConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"provisionedConcurrencyConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVersion) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Lambda::Version`.
func NewCfnVersion(scope awscdk.Construct, id *string, props *CfnVersionProps) CfnVersion {
	_init_.Initialize()

	j := jsiiProxy_CfnVersion{}

	_jsii_.Create(
		"monocdk.aws_lambda.CfnVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Lambda::Version`.
func NewCfnVersion_Override(c CfnVersion, scope awscdk.Construct, id *string, props *CfnVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CfnVersion",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnVersion) SetCodeSha256(val *string) {
	_jsii_.Set(
		j,
		"codeSha256",
		val,
	)
}

func (j *jsiiProxy_CfnVersion) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnVersion) SetFunctionName(val *string) {
	_jsii_.Set(
		j,
		"functionName",
		val,
	)
}

func (j *jsiiProxy_CfnVersion) SetProvisionedConcurrencyConfig(val interface{}) {
	_jsii_.Set(
		j,
		"provisionedConcurrencyConfig",
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
func CfnVersion_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnVersion",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnVersion_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnVersion",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CfnVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnVersion_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.CfnVersion",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnVersion) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnVersion) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnVersion) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnVersion) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnVersion) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnVersion) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnVersion) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnVersion) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnVersion) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnVersion) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnVersion) OnPrepare() {
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
func (c *jsiiProxy_CfnVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnVersion) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnVersion) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnVersion) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnVersion) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnVersion) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnVersion) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnVersion) ToString() *string {
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
func (c *jsiiProxy_CfnVersion) Validate() *[]*string {
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
func (c *jsiiProxy_CfnVersion) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnVersion_ProvisionedConcurrencyConfigurationProperty struct {
	// `CfnVersion.ProvisionedConcurrencyConfigurationProperty.ProvisionedConcurrentExecutions`.
	ProvisionedConcurrentExecutions *float64 `json:"provisionedConcurrentExecutions"`
}

// Properties for defining a `AWS::Lambda::Version`.
type CfnVersionProps struct {
	// `AWS::Lambda::Version.FunctionName`.
	FunctionName *string `json:"functionName"`
	// `AWS::Lambda::Version.CodeSha256`.
	CodeSha256 *string `json:"codeSha256"`
	// `AWS::Lambda::Version.Description`.
	Description *string `json:"description"`
	// `AWS::Lambda::Version.ProvisionedConcurrencyConfig`.
	ProvisionedConcurrencyConfig interface{} `json:"provisionedConcurrencyConfig"`
}

// Represents the Lambda Handler Code.
// Experimental.
type Code interface {
	IsInline() *bool
	Bind(scope awscdk.Construct) *CodeConfig
	BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions)
}

// The jsii proxy struct for Code
type jsiiProxy_Code struct {
	_ byte // padding
}

func (j *jsiiProxy_Code) IsInline() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isInline",
		&returns,
	)
	return returns
}


// Experimental.
func NewCode_Override(c Code) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.Code",
		nil, // no parameters
		c,
	)
}

// DEPRECATED.
// Deprecated: use `fromAsset`
func Code_Asset(path *string) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"asset",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromBucket`
func Code_Bucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"bucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromCfnParameters`
func Code_CfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"cfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from a local disk path.
// Experimental.
func Code_FromAsset(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func Code_FromAssetImage(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	var returns AssetImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"fromAssetImage",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Lambda handler code as an S3 object.
// Experimental.
func Code_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Creates a new Lambda source defined using CloudFormation parameters.
//
// Returns: a new instance of `CfnParametersCode`
// Experimental.
func Code_FromCfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"fromCfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from an asset created by a Docker build.
//
// By default, the asset is expected to be located at `/asset` in the
// image.
// Experimental.
func Code_FromDockerBuild(path *string, options *DockerBuildAssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"fromDockerBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func Code_FromEcrImage(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	var returns EcrImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"fromEcrImage",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Inline code for Lambda handler.
//
// Returns: `LambdaInlineCode` with inline code.
// Experimental.
func Code_FromInline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"fromInline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromInline`
func Code_Inline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Code",
		"inline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// Called when the lambda or layer is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (c *jsiiProxy_Code) Bind(scope awscdk.Construct) *CodeConfig {
	var returns *CodeConfig

	_jsii_.Invoke(
		c,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Called after the CFN function resource has been created to allow the code class to bind to it.
//
// Specifically it's required to allow assets to add
// metadata for tooling like SAM CLI to be able to find their origins.
// Experimental.
func (c *jsiiProxy_Code) BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions) {
	_jsii_.InvokeVoid(
		c,
		"bindToResource",
		[]interface{}{_resource, _options},
	)
}

// Result of binding `Code` into a `Function`.
// Experimental.
type CodeConfig struct {
	// Docker image configuration (mutually exclusive with `s3Location` and `inlineCode`).
	// Experimental.
	Image *CodeImageConfig `json:"image"`
	// Inline code (mutually exclusive with `s3Location` and `image`).
	// Experimental.
	InlineCode *string `json:"inlineCode"`
	// The location of the code in S3 (mutually exclusive with `inlineCode` and `image`).
	// Experimental.
	S3Location *awss3.Location `json:"s3Location"`
}

// Result of the bind when an ECR image is used.
// Experimental.
type CodeImageConfig struct {
	// URI to the Docker image.
	// Experimental.
	ImageUri *string `json:"imageUri"`
	// Specify or override the CMD on the specified Docker image or Dockerfile.
	//
	// This needs to be in the 'exec form', viz., `[ 'executable', 'param1', 'param2' ]`.
	// See: https://docs.docker.com/engine/reference/builder/#cmd
	//
	// Experimental.
	Cmd *[]*string `json:"cmd"`
	// Specify or override the ENTRYPOINT on the specified Docker image or Dockerfile.
	//
	// An ENTRYPOINT allows you to configure a container that will run as an executable.
	// This needs to be in the 'exec form', viz., `[ 'executable', 'param1', 'param2' ]`.
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	Entrypoint *[]*string `json:"entrypoint"`
}

// Defines a Code Signing Config.
// Experimental.
type CodeSigningConfig interface {
	awscdk.Resource
	ICodeSigningConfig
	CodeSigningConfigArn() *string
	CodeSigningConfigId() *string
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

// The jsii proxy struct for CodeSigningConfig
type jsiiProxy_CodeSigningConfig struct {
	internal.Type__awscdkResource
	jsiiProxy_ICodeSigningConfig
}

func (j *jsiiProxy_CodeSigningConfig) CodeSigningConfigArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codeSigningConfigArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CodeSigningConfig) CodeSigningConfigId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codeSigningConfigId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CodeSigningConfig) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CodeSigningConfig) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CodeSigningConfig) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CodeSigningConfig) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewCodeSigningConfig(scope constructs.Construct, id *string, props *CodeSigningConfigProps) CodeSigningConfig {
	_init_.Initialize()

	j := jsiiProxy_CodeSigningConfig{}

	_jsii_.Create(
		"monocdk.aws_lambda.CodeSigningConfig",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCodeSigningConfig_Override(c CodeSigningConfig, scope constructs.Construct, id *string, props *CodeSigningConfigProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.CodeSigningConfig",
		[]interface{}{scope, id, props},
		c,
	)
}

// Creates a Signing Profile construct that represents an external Signing Profile.
// Experimental.
func CodeSigningConfig_FromCodeSigningConfigArn(scope constructs.Construct, id *string, codeSigningConfigArn *string) ICodeSigningConfig {
	_init_.Initialize()

	var returns ICodeSigningConfig

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CodeSigningConfig",
		"fromCodeSigningConfigArn",
		[]interface{}{scope, id, codeSigningConfigArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CodeSigningConfig_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CodeSigningConfig",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func CodeSigningConfig_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.CodeSigningConfig",
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
func (c *jsiiProxy_CodeSigningConfig) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		c,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (c *jsiiProxy_CodeSigningConfig) GeneratePhysicalName() *string {
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
func (c *jsiiProxy_CodeSigningConfig) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (c *jsiiProxy_CodeSigningConfig) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CodeSigningConfig) OnPrepare() {
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
func (c *jsiiProxy_CodeSigningConfig) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CodeSigningConfig) OnValidate() *[]*string {
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
func (c *jsiiProxy_CodeSigningConfig) Prepare() {
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
func (c *jsiiProxy_CodeSigningConfig) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CodeSigningConfig) ToString() *string {
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
func (c *jsiiProxy_CodeSigningConfig) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for a Code Signing Config object.
// Experimental.
type CodeSigningConfigProps struct {
	// List of signing profiles that defines a trusted user who can sign a code package.
	// Experimental.
	SigningProfiles *[]awssigner.ISigningProfile `json:"signingProfiles"`
	// Code signing configuration description.
	// Experimental.
	Description *string `json:"description"`
	// Code signing configuration policy for deployment validation failure.
	//
	// If you set the policy to Enforce, Lambda blocks the deployment request
	// if signature validation checks fail.
	// If you set the policy to Warn, Lambda allows the deployment and
	// creates a CloudWatch log.
	// Experimental.
	UntrustedArtifactOnDeployment UntrustedArtifactOnDeployment `json:"untrustedArtifactOnDeployment"`
}

// A destination configuration.
// Experimental.
type DestinationConfig struct {
	// The Amazon Resource Name (ARN) of the destination resource.
	// Experimental.
	Destination *string `json:"destination"`
}

// Options when binding a destination to a function.
// Experimental.
type DestinationOptions struct {
	// The destination type.
	// Experimental.
	Type DestinationType `json:"type"`
}

// The type of destination.
// Experimental.
type DestinationType string

const (
	DestinationType_FAILURE DestinationType = "FAILURE"
	DestinationType_SUCCESS DestinationType = "SUCCESS"
)

// A destination configuration.
// Experimental.
type DlqDestinationConfig struct {
	// The Amazon Resource Name (ARN) of the destination resource.
	// Experimental.
	Destination *string `json:"destination"`
}

// Options when creating an asset from a Docker build.
// Experimental.
type DockerBuildAssetOptions struct {
	// Build args.
	// Experimental.
	BuildArgs *map[string]*string `json:"buildArgs"`
	// Name of the Dockerfile, must relative to the docker build path.
	// Experimental.
	File *string `json:"file"`
	// Set platform if server is multi-platform capable.
	//
	// _Requires Docker Engine API v1.38+_.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Platform *string `json:"platform"`
	// The path in the Docker image where the asset is located after the build operation.
	// Experimental.
	ImagePath *string `json:"imagePath"`
	// The path on the local filesystem where the asset will be copied using `docker cp`.
	// Experimental.
	OutputPath *string `json:"outputPath"`
}

// Code property for the DockerImageFunction construct.
// Experimental.
type DockerImageCode interface {
}

// The jsii proxy struct for DockerImageCode
type jsiiProxy_DockerImageCode struct {
	_ byte // padding
}

// Experimental.
func NewDockerImageCode_Override(d DockerImageCode) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.DockerImageCode",
		nil, // no parameters
		d,
	)
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func DockerImageCode_FromEcr(repository awsecr.IRepository, props *EcrImageCodeProps) DockerImageCode {
	_init_.Initialize()

	var returns DockerImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageCode",
		"fromEcr",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func DockerImageCode_FromImageAsset(directory *string, props *AssetImageCodeProps) DockerImageCode {
	_init_.Initialize()

	var returns DockerImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageCode",
		"fromImageAsset",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Create a lambda function where the handler is a docker image.
// Experimental.
type DockerImageFunction interface {
	Function
	CanCreatePermissions() *bool
	Connections() awsec2.Connections
	CurrentVersion() Version
	DeadLetterQueue() awssqs.IQueue
	Env() *awscdk.ResourceEnvironment
	FunctionArn() *string
	FunctionName() *string
	GrantPrincipal() awsiam.IPrincipal
	IsBoundToVpc() *bool
	LatestVersion() IVersion
	LogGroup() awslogs.ILogGroup
	Node() awscdk.ConstructNode
	PermissionsNode() awscdk.ConstructNode
	PhysicalName() *string
	Role() awsiam.IRole
	Runtime() Runtime
	Stack() awscdk.Stack
	AddEnvironment(key *string, value *string, options *EnvironmentOptions) Function
	AddEventSource(source IEventSource)
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	AddLayers(layers ...ILayerVersion)
	AddPermission(id *string, permission *Permission)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	AddVersion(name *string, codeSha256 *string, description *string, provisionedExecutions *float64, asyncInvokeConfig *EventInvokeConfigOptions) Version
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for DockerImageFunction
type jsiiProxy_DockerImageFunction struct {
	jsiiProxy_Function
}

func (j *jsiiProxy_DockerImageFunction) CanCreatePermissions() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"canCreatePermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) CurrentVersion() Version {
	var returns Version
	_jsii_.Get(
		j,
		"currentVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) DeadLetterQueue() awssqs.IQueue {
	var returns awssqs.IQueue
	_jsii_.Get(
		j,
		"deadLetterQueue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) LogGroup() awslogs.ILogGroup {
	var returns awslogs.ILogGroup
	_jsii_.Get(
		j,
		"logGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) Runtime() Runtime {
	var returns Runtime
	_jsii_.Get(
		j,
		"runtime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageFunction) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewDockerImageFunction(scope constructs.Construct, id *string, props *DockerImageFunctionProps) DockerImageFunction {
	_init_.Initialize()

	j := jsiiProxy_DockerImageFunction{}

	_jsii_.Create(
		"monocdk.aws_lambda.DockerImageFunction",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewDockerImageFunction_Override(d DockerImageFunction, scope constructs.Construct, id *string, props *DockerImageFunctionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.DockerImageFunction",
		[]interface{}{scope, id, props},
		d,
	)
}

// Record whether specific properties in the `AWS::Lambda::Function` resource should also be associated to the Version resource.
//
// See 'currentVersion' section in the module README for more details.
// Experimental.
func DockerImageFunction_ClassifyVersionProperty(propertyName *string, locked *bool) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_lambda.DockerImageFunction",
		"classifyVersionProperty",
		[]interface{}{propertyName, locked},
	)
}

// Import a lambda function into the CDK using its ARN.
// Experimental.
func DockerImageFunction_FromFunctionArn(scope constructs.Construct, id *string, functionArn *string) IFunction {
	_init_.Initialize()

	var returns IFunction

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"fromFunctionArn",
		[]interface{}{scope, id, functionArn},
		&returns,
	)

	return returns
}

// Creates a Lambda function object which represents a function not defined within this stack.
// Experimental.
func DockerImageFunction_FromFunctionAttributes(scope constructs.Construct, id *string, attrs *FunctionAttributes) IFunction {
	_init_.Initialize()

	var returns IFunction

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"fromFunctionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func DockerImageFunction_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func DockerImageFunction_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return the given named metric for this Lambda.
// Experimental.
func DockerImageFunction_MetricAll(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"metricAll",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of concurrent executions across all Lambdas.
// Experimental.
func DockerImageFunction_MetricAllConcurrentExecutions(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"metricAllConcurrentExecutions",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the Duration executing all Lambdas.
// Experimental.
func DockerImageFunction_MetricAllDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"metricAllDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of Errors executing all Lambdas.
// Experimental.
func DockerImageFunction_MetricAllErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"metricAllErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of invocations of all Lambdas.
// Experimental.
func DockerImageFunction_MetricAllInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"metricAllInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of throttled invocations of all Lambdas.
// Experimental.
func DockerImageFunction_MetricAllThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"metricAllThrottles",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of unreserved concurrent executions across all Lambdas.
// Experimental.
func DockerImageFunction_MetricAllUnreservedConcurrentExecutions(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.DockerImageFunction",
		"metricAllUnreservedConcurrentExecutions",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Adds an environment variable to this Lambda function.
//
// If this is a ref to a Lambda function, this operation results in a no-op.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) AddEnvironment(key *string, value *string, options *EnvironmentOptions) Function {
	var returns Function

	_jsii_.Invoke(
		d,
		"addEnvironment",
		[]interface{}{key, value, options},
		&returns,
	)

	return returns
}

// Adds an event source to this function.
//
// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
//
// The following example adds an SQS Queue as an event source:
// ```
// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
// myFunction.addEventSource(new SqsEventSource(myQueue));
// ```
// Experimental.
func (d *jsiiProxy_DockerImageFunction) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		d,
		"addEventSource",
		[]interface{}{source},
	)
}

// Adds an event source that maps to this AWS Lambda function.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		d,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds one or more Lambda Layers to this Lambda function.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) AddLayers(layers ...ILayerVersion) {
	args := []interface{}{}
	for _, a := range layers {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		d,
		"addLayers",
		args,
	)
}

// Adds a permission to the Lambda resource policy.
// See: Permission for details.
//
// Experimental.
func (d *jsiiProxy_DockerImageFunction) AddPermission(id *string, permission *Permission) {
	_jsii_.InvokeVoid(
		d,
		"addPermission",
		[]interface{}{id, permission},
	)
}

// Adds a statement to the IAM role assumed by the instance.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		d,
		"addToRolePolicy",
		[]interface{}{statement},
	)
}

// Add a new version for this Lambda.
//
// If you want to deploy through CloudFormation and use aliases, you need to
// add a new version (with a new name) to your Lambda every time you want to
// deploy an update. An alias can then refer to the newly created Version.
//
// All versions should have distinct names, and you should not delete versions
// as long as your Alias needs to refer to them.
//
// Returns: A new Version object.
// Deprecated: This method will create an AWS::Lambda::Version resource which
// snapshots the AWS Lambda function *at the time of its creation* and it
// won't get updated when the function changes. Instead, use
// `this.currentVersion` to obtain a reference to a version resource that gets
// automatically recreated when the function configuration (or code) changes.
func (d *jsiiProxy_DockerImageFunction) AddVersion(name *string, codeSha256 *string, description *string, provisionedExecutions *float64, asyncInvokeConfig *EventInvokeConfigOptions) Version {
	var returns Version

	_jsii_.Invoke(
		d,
		"addVersion",
		[]interface{}{name, codeSha256, description, provisionedExecutions, asyncInvokeConfig},
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
func (d *jsiiProxy_DockerImageFunction) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		d,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Configures options for asynchronous invocation.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		d,
		"configureAsyncInvoke",
		[]interface{}{options},
	)
}

// Experimental.
func (d *jsiiProxy_DockerImageFunction) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_DockerImageFunction) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_DockerImageFunction) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to invoke this Lambda.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		d,
		"grantInvoke",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Function.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// How long execution of this Lambda takes.
//
// Average over 5 minutes
// Experimental.
func (d *jsiiProxy_DockerImageFunction) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How many invocations of this Lambda fail.
//
// Sum over 5 minutes
// Experimental.
func (d *jsiiProxy_DockerImageFunction) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is invoked.
//
// Sum over 5 minutes
// Experimental.
func (d *jsiiProxy_DockerImageFunction) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is throttled.
//
// Sum over 5 minutes
// Experimental.
func (d *jsiiProxy_DockerImageFunction) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		d,
		"metricThrottles",
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
func (d *jsiiProxy_DockerImageFunction) OnPrepare() {
	_jsii_.InvokeVoid(
		d,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
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
func (d *jsiiProxy_DockerImageFunction) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_DockerImageFunction) Prepare() {
	_jsii_.InvokeVoid(
		d,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (d *jsiiProxy_DockerImageFunction) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_DockerImageFunction) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		d,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to configure a new DockerImageFunction construct.
// Experimental.
type DockerImageFunctionProps struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Whether to allow the Lambda to send all network traffic.
	//
	// If set to false, you must individually add traffic rules to allow the
	// Lambda to connect to network targets.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Lambda Functions in a public subnet can NOT access the internet.
	//
	// Use this property to acknowledge this limitation and still place the function in a public subnet.
	// See: https://stackoverflow.com/questions/52992085/why-cant-an-aws-lambda-function-inside-a-public-subnet-in-a-vpc-connect-to-the/52994841#52994841
	//
	// Experimental.
	AllowPublicSubnet *bool `json:"allowPublicSubnet"`
	// Code signing config associated with this function.
	// Experimental.
	CodeSigningConfig ICodeSigningConfig `json:"codeSigningConfig"`
	// Options for the `lambda.Version` resource automatically created by the `fn.currentVersion` method.
	// Experimental.
	CurrentVersionOptions *VersionOptions `json:"currentVersionOptions"`
	// The SQS queue to use if DLQ is enabled.
	// Experimental.
	DeadLetterQueue awssqs.IQueue `json:"deadLetterQueue"`
	// Enabled DLQ.
	//
	// If `deadLetterQueue` is undefined,
	// an SQS queue with default options will be defined for your Function.
	// Experimental.
	DeadLetterQueueEnabled *bool `json:"deadLetterQueueEnabled"`
	// A description of the function.
	// Experimental.
	Description *string `json:"description"`
	// Key-value pairs that Lambda caches and makes available for your Lambda functions.
	//
	// Use environment variables to apply configuration changes, such
	// as test and production environment configurations, without changing your
	// Lambda function source code.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The AWS KMS key that's used to encrypt your function's environment variables.
	// Experimental.
	EnvironmentEncryption awskms.IKey `json:"environmentEncryption"`
	// Event sources for this function.
	//
	// You can also add event sources using `addEventSource`.
	// Experimental.
	Events *[]IEventSource `json:"events"`
	// The filesystem configuration for the lambda function.
	// Experimental.
	Filesystem FileSystem `json:"filesystem"`
	// A name for the function.
	// Experimental.
	FunctionName *string `json:"functionName"`
	// Initial policy statements to add to the created Lambda Role.
	//
	// You can call `addToRolePolicy` to the created lambda to add statements post creation.
	// Experimental.
	InitialPolicy *[]awsiam.PolicyStatement `json:"initialPolicy"`
	// A list of layers to add to the function's execution environment.
	//
	// You can configure your Lambda function to pull in
	// additional code during initialization in the form of layers. Layers are packages of libraries or other dependencies
	// that can be used by multiple functions.
	// Experimental.
	Layers *[]ILayerVersion `json:"layers"`
	// The number of days log events are kept in CloudWatch Logs.
	//
	// When updating
	// this property, unsetting it doesn't remove the log retention policy. To
	// remove the retention policy, set the value to `INFINITE`.
	// Experimental.
	LogRetention awslogs.RetentionDays `json:"logRetention"`
	// When log retention is specified, a custom resource attempts to create the CloudWatch log group.
	//
	// These options control the retry policy when interacting with CloudWatch APIs.
	// Experimental.
	LogRetentionRetryOptions *LogRetentionRetryOptions `json:"logRetentionRetryOptions"`
	// The IAM role for the Lambda function associated with the custom resource that sets the retention policy.
	// Experimental.
	LogRetentionRole awsiam.IRole `json:"logRetentionRole"`
	// The amount of memory, in MB, that is allocated to your Lambda function.
	//
	// Lambda uses this value to proportionally allocate the amount of CPU
	// power. For more information, see Resource Model in the AWS Lambda
	// Developer Guide.
	// Experimental.
	MemorySize *float64 `json:"memorySize"`
	// Enable profiling.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	Profiling *bool `json:"profiling"`
	// Profiling Group.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	ProfilingGroup awscodeguruprofiler.IProfilingGroup `json:"profilingGroup"`
	// The maximum of concurrent executions you want to reserve for the function.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/concurrent-executions.html
	//
	// Experimental.
	ReservedConcurrentExecutions *float64 `json:"reservedConcurrentExecutions"`
	// Lambda execution role.
	//
	// This is the role that will be assumed by the function upon execution.
	// It controls the permissions that the function will have. The Role must
	// be assumable by the 'lambda.amazonaws.com' service principal.
	//
	// The default Role automatically has permissions granted for Lambda execution. If you
	// provide a Role, you must add the relevant AWS managed policies yourself.
	//
	// The relevant managed policies are "service-role/AWSLambdaBasicExecutionRole" and
	// "service-role/AWSLambdaVPCAccessExecutionRole".
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// What security group to associate with the Lambda's network interfaces. This property is being deprecated, consider using securityGroups instead.
	//
	// Only used if 'vpc' is supplied.
	//
	// Use securityGroups property instead.
	// Function constructor will throw an error if both are specified.
	// Deprecated: - This property is deprecated, use securityGroups instead
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// The list of security groups to associate with the Lambda's network interfaces.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The function execution time (in seconds) after which Lambda terminates the function.
	//
	// Because the execution time affects cost, set this value
	// based on the function's expected execution time.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// Enable AWS X-Ray Tracing for Lambda Function.
	// Experimental.
	Tracing Tracing `json:"tracing"`
	// VPC network to place Lambda network interfaces.
	//
	// Specify this if the Lambda function needs to access resources in a VPC.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Where to place the network interfaces within the VPC.
	//
	// Only used if 'vpc' is supplied. Note: internet access for Lambdas
	// requires a NAT gateway, so picking Public subnets is not allowed.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// The source code of your Lambda function.
	//
	// You can point to a file in an
	// Amazon Simple Storage Service (Amazon S3) bucket or specify your source
	// code as inline text.
	// Experimental.
	Code DockerImageCode `json:"code"`
}

// Represents a Docker image in ECR that can be bound as Lambda Code.
// Experimental.
type EcrImageCode interface {
	Code
	IsInline() *bool
	Bind(_arg awscdk.Construct) *CodeConfig
	BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions)
}

// The jsii proxy struct for EcrImageCode
type jsiiProxy_EcrImageCode struct {
	jsiiProxy_Code
}

func (j *jsiiProxy_EcrImageCode) IsInline() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isInline",
		&returns,
	)
	return returns
}


// Experimental.
func NewEcrImageCode(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	j := jsiiProxy_EcrImageCode{}

	_jsii_.Create(
		"monocdk.aws_lambda.EcrImageCode",
		[]interface{}{repository, props},
		&j,
	)

	return &j
}

// Experimental.
func NewEcrImageCode_Override(e EcrImageCode, repository awsecr.IRepository, props *EcrImageCodeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.EcrImageCode",
		[]interface{}{repository, props},
		e,
	)
}

// DEPRECATED.
// Deprecated: use `fromAsset`
func EcrImageCode_Asset(path *string) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"asset",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromBucket`
func EcrImageCode_Bucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"bucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromCfnParameters`
func EcrImageCode_CfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"cfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from a local disk path.
// Experimental.
func EcrImageCode_FromAsset(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func EcrImageCode_FromAssetImage(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	var returns AssetImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"fromAssetImage",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Lambda handler code as an S3 object.
// Experimental.
func EcrImageCode_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Creates a new Lambda source defined using CloudFormation parameters.
//
// Returns: a new instance of `CfnParametersCode`
// Experimental.
func EcrImageCode_FromCfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"fromCfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from an asset created by a Docker build.
//
// By default, the asset is expected to be located at `/asset` in the
// image.
// Experimental.
func EcrImageCode_FromDockerBuild(path *string, options *DockerBuildAssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"fromDockerBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func EcrImageCode_FromEcrImage(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	var returns EcrImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"fromEcrImage",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Inline code for Lambda handler.
//
// Returns: `LambdaInlineCode` with inline code.
// Experimental.
func EcrImageCode_FromInline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"fromInline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromInline`
func EcrImageCode_Inline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EcrImageCode",
		"inline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// Called when the lambda or layer is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (e *jsiiProxy_EcrImageCode) Bind(_arg awscdk.Construct) *CodeConfig {
	var returns *CodeConfig

	_jsii_.Invoke(
		e,
		"bind",
		[]interface{}{_arg},
		&returns,
	)

	return returns
}

// Called after the CFN function resource has been created to allow the code class to bind to it.
//
// Specifically it's required to allow assets to add
// metadata for tooling like SAM CLI to be able to find their origins.
// Experimental.
func (e *jsiiProxy_EcrImageCode) BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions) {
	_jsii_.InvokeVoid(
		e,
		"bindToResource",
		[]interface{}{_resource, _options},
	)
}

// Properties to initialize a new EcrImageCode.
// Experimental.
type EcrImageCodeProps struct {
	// Specify or override the CMD on the specified Docker image or Dockerfile.
	//
	// This needs to be in the 'exec form', viz., `[ 'executable', 'param1', 'param2' ]`.
	// See: https://docs.docker.com/engine/reference/builder/#cmd
	//
	// Experimental.
	Cmd *[]*string `json:"cmd"`
	// Specify or override the ENTRYPOINT on the specified Docker image or Dockerfile.
	//
	// An ENTRYPOINT allows you to configure a container that will run as an executable.
	// This needs to be in the 'exec form', viz., `[ 'executable', 'param1', 'param2' ]`.
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	Entrypoint *[]*string `json:"entrypoint"`
	// The image tag to use when pulling the image from ECR.
	// Experimental.
	Tag *string `json:"tag"`
}

// Environment variables options.
// Experimental.
type EnvironmentOptions struct {
	// When used in Lambda@Edge via edgeArn() API, these environment variables will be removed.
	//
	// If not set, an error will be thrown.
	// See: https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/lambda-requirements-limits.html#lambda-requirements-lambda-function-configuration
	//
	// Experimental.
	RemoveInEdge *bool `json:"removeInEdge"`
}

// Configure options for asynchronous invocation on a version or an alias.
//
// By default, Lambda retries an asynchronous invocation twice if the function
// returns an error. It retains events in a queue for up to six hours. When an
// event fails all processing attempts or stays in the asynchronous invocation
// queue for too long, Lambda discards it.
// Experimental.
type EventInvokeConfig interface {
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

// The jsii proxy struct for EventInvokeConfig
type jsiiProxy_EventInvokeConfig struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_EventInvokeConfig) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventInvokeConfig) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventInvokeConfig) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventInvokeConfig) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewEventInvokeConfig(scope constructs.Construct, id *string, props *EventInvokeConfigProps) EventInvokeConfig {
	_init_.Initialize()

	j := jsiiProxy_EventInvokeConfig{}

	_jsii_.Create(
		"monocdk.aws_lambda.EventInvokeConfig",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewEventInvokeConfig_Override(e EventInvokeConfig, scope constructs.Construct, id *string, props *EventInvokeConfigProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.EventInvokeConfig",
		[]interface{}{scope, id, props},
		e,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func EventInvokeConfig_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EventInvokeConfig",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func EventInvokeConfig_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EventInvokeConfig",
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
func (e *jsiiProxy_EventInvokeConfig) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		e,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (e *jsiiProxy_EventInvokeConfig) GeneratePhysicalName() *string {
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
func (e *jsiiProxy_EventInvokeConfig) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (e *jsiiProxy_EventInvokeConfig) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		e,
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
func (e *jsiiProxy_EventInvokeConfig) OnPrepare() {
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
func (e *jsiiProxy_EventInvokeConfig) OnSynthesize(session constructs.ISynthesisSession) {
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
func (e *jsiiProxy_EventInvokeConfig) OnValidate() *[]*string {
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
func (e *jsiiProxy_EventInvokeConfig) Prepare() {
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
func (e *jsiiProxy_EventInvokeConfig) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		e,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (e *jsiiProxy_EventInvokeConfig) ToString() *string {
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
func (e *jsiiProxy_EventInvokeConfig) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		e,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options to add an EventInvokeConfig to a function.
// Experimental.
type EventInvokeConfigOptions struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
}

// Properties for an EventInvokeConfig.
// Experimental.
type EventInvokeConfigProps struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// The Lambda function.
	// Experimental.
	Function IFunction `json:"function"`
	// The qualifier.
	// Experimental.
	Qualifier *string `json:"qualifier"`
}

// Defines a Lambda EventSourceMapping resource.
//
// Usually, you won't need to define the mapping yourself. This will usually be done by
// event sources. For example, to add an SQS event source to a function:
//
//     import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
//     lambda.addEventSource(new SqsEventSource(sqs));
//
// The `SqsEventSource` class will automatically create the mapping, and will also
// modify the Lambda's execution role so it can consume messages from the queue.
// Experimental.
type EventSourceMapping interface {
	awscdk.Resource
	IEventSourceMapping
	Env() *awscdk.ResourceEnvironment
	EventSourceMappingId() *string
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

// The jsii proxy struct for EventSourceMapping
type jsiiProxy_EventSourceMapping struct {
	internal.Type__awscdkResource
	jsiiProxy_IEventSourceMapping
}

func (j *jsiiProxy_EventSourceMapping) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventSourceMapping) EventSourceMappingId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventSourceMappingId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventSourceMapping) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventSourceMapping) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_EventSourceMapping) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewEventSourceMapping(scope constructs.Construct, id *string, props *EventSourceMappingProps) EventSourceMapping {
	_init_.Initialize()

	j := jsiiProxy_EventSourceMapping{}

	_jsii_.Create(
		"monocdk.aws_lambda.EventSourceMapping",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewEventSourceMapping_Override(e EventSourceMapping, scope constructs.Construct, id *string, props *EventSourceMappingProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.EventSourceMapping",
		[]interface{}{scope, id, props},
		e,
	)
}

// Import an event source into this stack from its event source id.
// Experimental.
func EventSourceMapping_FromEventSourceMappingId(scope constructs.Construct, id *string, eventSourceMappingId *string) IEventSourceMapping {
	_init_.Initialize()

	var returns IEventSourceMapping

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EventSourceMapping",
		"fromEventSourceMappingId",
		[]interface{}{scope, id, eventSourceMappingId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func EventSourceMapping_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EventSourceMapping",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func EventSourceMapping_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.EventSourceMapping",
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
func (e *jsiiProxy_EventSourceMapping) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		e,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (e *jsiiProxy_EventSourceMapping) GeneratePhysicalName() *string {
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
func (e *jsiiProxy_EventSourceMapping) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (e *jsiiProxy_EventSourceMapping) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		e,
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
func (e *jsiiProxy_EventSourceMapping) OnPrepare() {
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
func (e *jsiiProxy_EventSourceMapping) OnSynthesize(session constructs.ISynthesisSession) {
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
func (e *jsiiProxy_EventSourceMapping) OnValidate() *[]*string {
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
func (e *jsiiProxy_EventSourceMapping) Prepare() {
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
func (e *jsiiProxy_EventSourceMapping) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		e,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (e *jsiiProxy_EventSourceMapping) ToString() *string {
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
func (e *jsiiProxy_EventSourceMapping) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		e,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type EventSourceMappingOptions struct {
	// The largest number of records that AWS Lambda will retrieve from your event source at the time of invoking your function.
	//
	// Your function receives an
	// event with all the retrieved records.
	//
	// Valid Range: Minimum value of 1. Maximum value of 10000.
	// Experimental.
	BatchSize *float64 `json:"batchSize"`
	// If the function returns an error, split the batch in two and retry.
	// Experimental.
	BisectBatchOnError *bool `json:"bisectBatchOnError"`
	// Set to false to disable the event source upon creation.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// The Amazon Resource Name (ARN) of the event source.
	//
	// Any record added to
	// this stream can invoke the Lambda function.
	// Experimental.
	EventSourceArn *string `json:"eventSourceArn"`
	// A list of host and port pairs that are the addresses of the Kafka brokers in a self managed "bootstrap" Kafka cluster that a Kafka client connects to initially to bootstrap itself.
	//
	// They are in the format `abc.example.com:9096`.
	// Experimental.
	KafkaBootstrapServers *[]*string `json:"kafkaBootstrapServers"`
	// The name of the Kafka topic.
	// Experimental.
	KafkaTopic *string `json:"kafkaTopic"`
	// The maximum amount of time to gather records before invoking the function.
	//
	// Maximum of Duration.minutes(5)
	// Experimental.
	MaxBatchingWindow awscdk.Duration `json:"maxBatchingWindow"`
	// The maximum age of a record that Lambda sends to a function for processing.
	//
	// Valid Range:
	// * Minimum value of 60 seconds
	// * Maximum value of 7 days
	// Experimental.
	MaxRecordAge awscdk.Duration `json:"maxRecordAge"`
	// An Amazon SQS queue or Amazon SNS topic destination for discarded records.
	// Experimental.
	OnFailure IEventSourceDlq `json:"onFailure"`
	// The number of batches to process from each shard concurrently.
	//
	// Valid Range:
	// * Minimum value of 1
	// * Maximum value of 10
	// Experimental.
	ParallelizationFactor *float64 `json:"parallelizationFactor"`
	// Allow functions to return partially successful responses for a batch of records.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-ddb-batchfailurereporting
	//
	// Experimental.
	ReportBatchItemFailures *bool `json:"reportBatchItemFailures"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Set to `undefined` if you want lambda to keep retrying infinitely or until
	// the record expires.
	//
	// Valid Range:
	// * Minimum value of 0
	// * Maximum value of 10000
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Specific settings like the authentication protocol or the VPC components to secure access to your event source.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-eventsourcemapping-sourceaccessconfiguration.html
	//
	// Experimental.
	SourceAccessConfigurations *[]*SourceAccessConfiguration `json:"sourceAccessConfigurations"`
	// The position in the DynamoDB, Kinesis or MSK stream where AWS Lambda should start reading.
	// See: https://docs.aws.amazon.com/kinesis/latest/APIReference/API_GetShardIterator.html#Kinesis-GetShardIterator-request-ShardIteratorType
	//
	// Experimental.
	StartingPosition StartingPosition `json:"startingPosition"`
	// The size of the tumbling windows to group records sent to DynamoDB or Kinesis.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-ddb-windows
	//
	// Valid Range: 0 - 15 minutes
	//
	// Experimental.
	TumblingWindow awscdk.Duration `json:"tumblingWindow"`
}

// Properties for declaring a new event source mapping.
// Experimental.
type EventSourceMappingProps struct {
	// The largest number of records that AWS Lambda will retrieve from your event source at the time of invoking your function.
	//
	// Your function receives an
	// event with all the retrieved records.
	//
	// Valid Range: Minimum value of 1. Maximum value of 10000.
	// Experimental.
	BatchSize *float64 `json:"batchSize"`
	// If the function returns an error, split the batch in two and retry.
	// Experimental.
	BisectBatchOnError *bool `json:"bisectBatchOnError"`
	// Set to false to disable the event source upon creation.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// The Amazon Resource Name (ARN) of the event source.
	//
	// Any record added to
	// this stream can invoke the Lambda function.
	// Experimental.
	EventSourceArn *string `json:"eventSourceArn"`
	// A list of host and port pairs that are the addresses of the Kafka brokers in a self managed "bootstrap" Kafka cluster that a Kafka client connects to initially to bootstrap itself.
	//
	// They are in the format `abc.example.com:9096`.
	// Experimental.
	KafkaBootstrapServers *[]*string `json:"kafkaBootstrapServers"`
	// The name of the Kafka topic.
	// Experimental.
	KafkaTopic *string `json:"kafkaTopic"`
	// The maximum amount of time to gather records before invoking the function.
	//
	// Maximum of Duration.minutes(5)
	// Experimental.
	MaxBatchingWindow awscdk.Duration `json:"maxBatchingWindow"`
	// The maximum age of a record that Lambda sends to a function for processing.
	//
	// Valid Range:
	// * Minimum value of 60 seconds
	// * Maximum value of 7 days
	// Experimental.
	MaxRecordAge awscdk.Duration `json:"maxRecordAge"`
	// An Amazon SQS queue or Amazon SNS topic destination for discarded records.
	// Experimental.
	OnFailure IEventSourceDlq `json:"onFailure"`
	// The number of batches to process from each shard concurrently.
	//
	// Valid Range:
	// * Minimum value of 1
	// * Maximum value of 10
	// Experimental.
	ParallelizationFactor *float64 `json:"parallelizationFactor"`
	// Allow functions to return partially successful responses for a batch of records.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-ddb-batchfailurereporting
	//
	// Experimental.
	ReportBatchItemFailures *bool `json:"reportBatchItemFailures"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Set to `undefined` if you want lambda to keep retrying infinitely or until
	// the record expires.
	//
	// Valid Range:
	// * Minimum value of 0
	// * Maximum value of 10000
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Specific settings like the authentication protocol or the VPC components to secure access to your event source.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-eventsourcemapping-sourceaccessconfiguration.html
	//
	// Experimental.
	SourceAccessConfigurations *[]*SourceAccessConfiguration `json:"sourceAccessConfigurations"`
	// The position in the DynamoDB, Kinesis or MSK stream where AWS Lambda should start reading.
	// See: https://docs.aws.amazon.com/kinesis/latest/APIReference/API_GetShardIterator.html#Kinesis-GetShardIterator-request-ShardIteratorType
	//
	// Experimental.
	StartingPosition StartingPosition `json:"startingPosition"`
	// The size of the tumbling windows to group records sent to DynamoDB or Kinesis.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-ddb-windows
	//
	// Valid Range: 0 - 15 minutes
	//
	// Experimental.
	TumblingWindow awscdk.Duration `json:"tumblingWindow"`
	// The target AWS Lambda function.
	// Experimental.
	Target IFunction `json:"target"`
}

// Represents the filesystem for the Lambda function.
// Experimental.
type FileSystem interface {
	Config() *FileSystemConfig
}

// The jsii proxy struct for FileSystem
type jsiiProxy_FileSystem struct {
	_ byte // padding
}

func (j *jsiiProxy_FileSystem) Config() *FileSystemConfig {
	var returns *FileSystemConfig
	_jsii_.Get(
		j,
		"config",
		&returns,
	)
	return returns
}


// Experimental.
func NewFileSystem(config *FileSystemConfig) FileSystem {
	_init_.Initialize()

	j := jsiiProxy_FileSystem{}

	_jsii_.Create(
		"monocdk.aws_lambda.FileSystem",
		[]interface{}{config},
		&j,
	)

	return &j
}

// Experimental.
func NewFileSystem_Override(f FileSystem, config *FileSystemConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.FileSystem",
		[]interface{}{config},
		f,
	)
}

// mount the filesystem from Amazon EFS.
// Experimental.
func FileSystem_FromEfsAccessPoint(ap awsefs.IAccessPoint, mountPath *string) FileSystem {
	_init_.Initialize()

	var returns FileSystem

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.FileSystem",
		"fromEfsAccessPoint",
		[]interface{}{ap, mountPath},
		&returns,
	)

	return returns
}

// FileSystem configurations for the Lambda function.
// Experimental.
type FileSystemConfig struct {
	// ARN of the access point.
	// Experimental.
	Arn *string `json:"arn"`
	// mount path in the lambda runtime environment.
	// Experimental.
	LocalMountPath *string `json:"localMountPath"`
	// connections object used to allow ingress traffic from lambda function.
	// Experimental.
	Connections awsec2.Connections `json:"connections"`
	// array of IDependable that lambda function depends on.
	// Experimental.
	Dependency *[]awscdk.IDependable `json:"dependency"`
	// additional IAM policies required for the lambda function.
	// Experimental.
	Policies *[]awsiam.PolicyStatement `json:"policies"`
}

// Deploys a file from inside the construct library as a function.
//
// The supplied file is subject to the 4096 bytes limit of being embedded in a
// CloudFormation template.
//
// The construct includes an associated role with the lambda.
//
// This construct does not yet reproduce all features from the underlying resource
// library.
// Experimental.
type Function interface {
	FunctionBase
	CanCreatePermissions() *bool
	Connections() awsec2.Connections
	CurrentVersion() Version
	DeadLetterQueue() awssqs.IQueue
	Env() *awscdk.ResourceEnvironment
	FunctionArn() *string
	FunctionName() *string
	GrantPrincipal() awsiam.IPrincipal
	IsBoundToVpc() *bool
	LatestVersion() IVersion
	LogGroup() awslogs.ILogGroup
	Node() awscdk.ConstructNode
	PermissionsNode() awscdk.ConstructNode
	PhysicalName() *string
	Role() awsiam.IRole
	Runtime() Runtime
	Stack() awscdk.Stack
	AddEnvironment(key *string, value *string, options *EnvironmentOptions) Function
	AddEventSource(source IEventSource)
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	AddLayers(layers ...ILayerVersion)
	AddPermission(id *string, permission *Permission)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	AddVersion(name *string, codeSha256 *string, description *string, provisionedExecutions *float64, asyncInvokeConfig *EventInvokeConfigOptions) Version
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Function
type jsiiProxy_Function struct {
	jsiiProxy_FunctionBase
}

func (j *jsiiProxy_Function) CanCreatePermissions() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"canCreatePermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) CurrentVersion() Version {
	var returns Version
	_jsii_.Get(
		j,
		"currentVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) DeadLetterQueue() awssqs.IQueue {
	var returns awssqs.IQueue
	_jsii_.Get(
		j,
		"deadLetterQueue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) LogGroup() awslogs.ILogGroup {
	var returns awslogs.ILogGroup
	_jsii_.Get(
		j,
		"logGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) Runtime() Runtime {
	var returns Runtime
	_jsii_.Get(
		j,
		"runtime",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Function) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewFunction(scope constructs.Construct, id *string, props *FunctionProps) Function {
	_init_.Initialize()

	j := jsiiProxy_Function{}

	_jsii_.Create(
		"monocdk.aws_lambda.Function",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewFunction_Override(f Function, scope constructs.Construct, id *string, props *FunctionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.Function",
		[]interface{}{scope, id, props},
		f,
	)
}

// Record whether specific properties in the `AWS::Lambda::Function` resource should also be associated to the Version resource.
//
// See 'currentVersion' section in the module README for more details.
// Experimental.
func Function_ClassifyVersionProperty(propertyName *string, locked *bool) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.aws_lambda.Function",
		"classifyVersionProperty",
		[]interface{}{propertyName, locked},
	)
}

// Import a lambda function into the CDK using its ARN.
// Experimental.
func Function_FromFunctionArn(scope constructs.Construct, id *string, functionArn *string) IFunction {
	_init_.Initialize()

	var returns IFunction

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"fromFunctionArn",
		[]interface{}{scope, id, functionArn},
		&returns,
	)

	return returns
}

// Creates a Lambda function object which represents a function not defined within this stack.
// Experimental.
func Function_FromFunctionAttributes(scope constructs.Construct, id *string, attrs *FunctionAttributes) IFunction {
	_init_.Initialize()

	var returns IFunction

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"fromFunctionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Function_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Function_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return the given named metric for this Lambda.
// Experimental.
func Function_MetricAll(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"metricAll",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of concurrent executions across all Lambdas.
// Experimental.
func Function_MetricAllConcurrentExecutions(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"metricAllConcurrentExecutions",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the Duration executing all Lambdas.
// Experimental.
func Function_MetricAllDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"metricAllDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of Errors executing all Lambdas.
// Experimental.
func Function_MetricAllErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"metricAllErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of invocations of all Lambdas.
// Experimental.
func Function_MetricAllInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"metricAllInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of throttled invocations of all Lambdas.
// Experimental.
func Function_MetricAllThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"metricAllThrottles",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of unreserved concurrent executions across all Lambdas.
// Experimental.
func Function_MetricAllUnreservedConcurrentExecutions(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	_init_.Initialize()

	var returns awscloudwatch.Metric

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Function",
		"metricAllUnreservedConcurrentExecutions",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Adds an environment variable to this Lambda function.
//
// If this is a ref to a Lambda function, this operation results in a no-op.
// Experimental.
func (f *jsiiProxy_Function) AddEnvironment(key *string, value *string, options *EnvironmentOptions) Function {
	var returns Function

	_jsii_.Invoke(
		f,
		"addEnvironment",
		[]interface{}{key, value, options},
		&returns,
	)

	return returns
}

// Adds an event source to this function.
//
// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
//
// The following example adds an SQS Queue as an event source:
// ```
// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
// myFunction.addEventSource(new SqsEventSource(myQueue));
// ```
// Experimental.
func (f *jsiiProxy_Function) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		f,
		"addEventSource",
		[]interface{}{source},
	)
}

// Adds an event source that maps to this AWS Lambda function.
// Experimental.
func (f *jsiiProxy_Function) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		f,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds one or more Lambda Layers to this Lambda function.
// Experimental.
func (f *jsiiProxy_Function) AddLayers(layers ...ILayerVersion) {
	args := []interface{}{}
	for _, a := range layers {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		f,
		"addLayers",
		args,
	)
}

// Adds a permission to the Lambda resource policy.
// See: Permission for details.
//
// Experimental.
func (f *jsiiProxy_Function) AddPermission(id *string, permission *Permission) {
	_jsii_.InvokeVoid(
		f,
		"addPermission",
		[]interface{}{id, permission},
	)
}

// Adds a statement to the IAM role assumed by the instance.
// Experimental.
func (f *jsiiProxy_Function) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		f,
		"addToRolePolicy",
		[]interface{}{statement},
	)
}

// Add a new version for this Lambda.
//
// If you want to deploy through CloudFormation and use aliases, you need to
// add a new version (with a new name) to your Lambda every time you want to
// deploy an update. An alias can then refer to the newly created Version.
//
// All versions should have distinct names, and you should not delete versions
// as long as your Alias needs to refer to them.
//
// Returns: A new Version object.
// Deprecated: This method will create an AWS::Lambda::Version resource which
// snapshots the AWS Lambda function *at the time of its creation* and it
// won't get updated when the function changes. Instead, use
// `this.currentVersion` to obtain a reference to a version resource that gets
// automatically recreated when the function configuration (or code) changes.
func (f *jsiiProxy_Function) AddVersion(name *string, codeSha256 *string, description *string, provisionedExecutions *float64, asyncInvokeConfig *EventInvokeConfigOptions) Version {
	var returns Version

	_jsii_.Invoke(
		f,
		"addVersion",
		[]interface{}{name, codeSha256, description, provisionedExecutions, asyncInvokeConfig},
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
func (f *jsiiProxy_Function) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		f,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Configures options for asynchronous invocation.
// Experimental.
func (f *jsiiProxy_Function) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		f,
		"configureAsyncInvoke",
		[]interface{}{options},
	)
}

// Experimental.
func (f *jsiiProxy_Function) GeneratePhysicalName() *string {
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
func (f *jsiiProxy_Function) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (f *jsiiProxy_Function) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to invoke this Lambda.
// Experimental.
func (f *jsiiProxy_Function) GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		f,
		"grantInvoke",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Function.
// Experimental.
func (f *jsiiProxy_Function) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// How long execution of this Lambda takes.
//
// Average over 5 minutes
// Experimental.
func (f *jsiiProxy_Function) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How many invocations of this Lambda fail.
//
// Sum over 5 minutes
// Experimental.
func (f *jsiiProxy_Function) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is invoked.
//
// Sum over 5 minutes
// Experimental.
func (f *jsiiProxy_Function) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is throttled.
//
// Sum over 5 minutes
// Experimental.
func (f *jsiiProxy_Function) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricThrottles",
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
func (f *jsiiProxy_Function) OnPrepare() {
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
func (f *jsiiProxy_Function) OnSynthesize(session constructs.ISynthesisSession) {
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
func (f *jsiiProxy_Function) OnValidate() *[]*string {
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
func (f *jsiiProxy_Function) Prepare() {
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
func (f *jsiiProxy_Function) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (f *jsiiProxy_Function) ToString() *string {
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
func (f *jsiiProxy_Function) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents a Lambda function defined outside of this stack.
// Experimental.
type FunctionAttributes struct {
	// The ARN of the Lambda function.
	//
	// Format: arn:<partition>:lambda:<region>:<account-id>:function:<function-name>
	// Experimental.
	FunctionArn *string `json:"functionArn"`
	// The IAM execution role associated with this function.
	//
	// If the role is not specified, any role-related operations will no-op.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// Setting this property informs the CDK that the imported function is in the same environment as the stack.
	//
	// This affects certain behaviours such as, whether this function's permission can be modified.
	// When not configured, the CDK attempts to auto-determine this. For environment agnostic stacks, i.e., stacks
	// where the account is not specified with the `env` property, this is determined to be false.
	//
	// Set this to property *ONLY IF* the imported function is in the same account as the stack
	// it's imported in.
	// Experimental.
	SameEnvironment *bool `json:"sameEnvironment"`
	// The security group of this Lambda, if in a VPC.
	//
	// This needs to be given in order to support allowing connections
	// to this Lambda.
	// Experimental.
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// Id of the security group of this Lambda, if in a VPC.
	//
	// This needs to be given in order to support allowing connections
	// to this Lambda.
	// Deprecated: use `securityGroup` instead
	SecurityGroupId *string `json:"securityGroupId"`
}

// Experimental.
type FunctionBase interface {
	awscdk.Resource
	awsec2.IClientVpnConnectionHandler
	IFunction
	CanCreatePermissions() *bool
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	FunctionArn() *string
	FunctionName() *string
	GrantPrincipal() awsiam.IPrincipal
	IsBoundToVpc() *bool
	LatestVersion() IVersion
	Node() awscdk.ConstructNode
	PermissionsNode() awscdk.ConstructNode
	PhysicalName() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	AddEventSource(source IEventSource)
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	AddPermission(id *string, permission *Permission)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for FunctionBase
type jsiiProxy_FunctionBase struct {
	internal.Type__awscdkResource
	internal.Type__awsec2IClientVpnConnectionHandler
	jsiiProxy_IFunction
}

func (j *jsiiProxy_FunctionBase) CanCreatePermissions() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"canCreatePermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FunctionBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewFunctionBase_Override(f FunctionBase, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.FunctionBase",
		[]interface{}{scope, id, props},
		f,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func FunctionBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.FunctionBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func FunctionBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.FunctionBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds an event source to this function.
//
// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
//
// The following example adds an SQS Queue as an event source:
// ```
// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
// myFunction.addEventSource(new SqsEventSource(myQueue));
// ```
// Experimental.
func (f *jsiiProxy_FunctionBase) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		f,
		"addEventSource",
		[]interface{}{source},
	)
}

// Adds an event source that maps to this AWS Lambda function.
// Experimental.
func (f *jsiiProxy_FunctionBase) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		f,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a permission to the Lambda resource policy.
// See: Permission for details.
//
// Experimental.
func (f *jsiiProxy_FunctionBase) AddPermission(id *string, permission *Permission) {
	_jsii_.InvokeVoid(
		f,
		"addPermission",
		[]interface{}{id, permission},
	)
}

// Adds a statement to the IAM role assumed by the instance.
// Experimental.
func (f *jsiiProxy_FunctionBase) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		f,
		"addToRolePolicy",
		[]interface{}{statement},
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
func (f *jsiiProxy_FunctionBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		f,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Configures options for asynchronous invocation.
// Experimental.
func (f *jsiiProxy_FunctionBase) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		f,
		"configureAsyncInvoke",
		[]interface{}{options},
	)
}

// Experimental.
func (f *jsiiProxy_FunctionBase) GeneratePhysicalName() *string {
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
func (f *jsiiProxy_FunctionBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (f *jsiiProxy_FunctionBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to invoke this Lambda.
// Experimental.
func (f *jsiiProxy_FunctionBase) GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		f,
		"grantInvoke",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Function.
// Experimental.
func (f *jsiiProxy_FunctionBase) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// How long execution of this Lambda takes.
//
// Average over 5 minutes
// Experimental.
func (f *jsiiProxy_FunctionBase) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How many invocations of this Lambda fail.
//
// Sum over 5 minutes
// Experimental.
func (f *jsiiProxy_FunctionBase) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is invoked.
//
// Sum over 5 minutes
// Experimental.
func (f *jsiiProxy_FunctionBase) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is throttled.
//
// Sum over 5 minutes
// Experimental.
func (f *jsiiProxy_FunctionBase) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		f,
		"metricThrottles",
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
func (f *jsiiProxy_FunctionBase) OnPrepare() {
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
func (f *jsiiProxy_FunctionBase) OnSynthesize(session constructs.ISynthesisSession) {
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
func (f *jsiiProxy_FunctionBase) OnValidate() *[]*string {
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
func (f *jsiiProxy_FunctionBase) Prepare() {
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
func (f *jsiiProxy_FunctionBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		f,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (f *jsiiProxy_FunctionBase) ToString() *string {
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
func (f *jsiiProxy_FunctionBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		f,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Non runtime options.
// Experimental.
type FunctionOptions struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Whether to allow the Lambda to send all network traffic.
	//
	// If set to false, you must individually add traffic rules to allow the
	// Lambda to connect to network targets.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Lambda Functions in a public subnet can NOT access the internet.
	//
	// Use this property to acknowledge this limitation and still place the function in a public subnet.
	// See: https://stackoverflow.com/questions/52992085/why-cant-an-aws-lambda-function-inside-a-public-subnet-in-a-vpc-connect-to-the/52994841#52994841
	//
	// Experimental.
	AllowPublicSubnet *bool `json:"allowPublicSubnet"`
	// Code signing config associated with this function.
	// Experimental.
	CodeSigningConfig ICodeSigningConfig `json:"codeSigningConfig"`
	// Options for the `lambda.Version` resource automatically created by the `fn.currentVersion` method.
	// Experimental.
	CurrentVersionOptions *VersionOptions `json:"currentVersionOptions"`
	// The SQS queue to use if DLQ is enabled.
	// Experimental.
	DeadLetterQueue awssqs.IQueue `json:"deadLetterQueue"`
	// Enabled DLQ.
	//
	// If `deadLetterQueue` is undefined,
	// an SQS queue with default options will be defined for your Function.
	// Experimental.
	DeadLetterQueueEnabled *bool `json:"deadLetterQueueEnabled"`
	// A description of the function.
	// Experimental.
	Description *string `json:"description"`
	// Key-value pairs that Lambda caches and makes available for your Lambda functions.
	//
	// Use environment variables to apply configuration changes, such
	// as test and production environment configurations, without changing your
	// Lambda function source code.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The AWS KMS key that's used to encrypt your function's environment variables.
	// Experimental.
	EnvironmentEncryption awskms.IKey `json:"environmentEncryption"`
	// Event sources for this function.
	//
	// You can also add event sources using `addEventSource`.
	// Experimental.
	Events *[]IEventSource `json:"events"`
	// The filesystem configuration for the lambda function.
	// Experimental.
	Filesystem FileSystem `json:"filesystem"`
	// A name for the function.
	// Experimental.
	FunctionName *string `json:"functionName"`
	// Initial policy statements to add to the created Lambda Role.
	//
	// You can call `addToRolePolicy` to the created lambda to add statements post creation.
	// Experimental.
	InitialPolicy *[]awsiam.PolicyStatement `json:"initialPolicy"`
	// A list of layers to add to the function's execution environment.
	//
	// You can configure your Lambda function to pull in
	// additional code during initialization in the form of layers. Layers are packages of libraries or other dependencies
	// that can be used by multiple functions.
	// Experimental.
	Layers *[]ILayerVersion `json:"layers"`
	// The number of days log events are kept in CloudWatch Logs.
	//
	// When updating
	// this property, unsetting it doesn't remove the log retention policy. To
	// remove the retention policy, set the value to `INFINITE`.
	// Experimental.
	LogRetention awslogs.RetentionDays `json:"logRetention"`
	// When log retention is specified, a custom resource attempts to create the CloudWatch log group.
	//
	// These options control the retry policy when interacting with CloudWatch APIs.
	// Experimental.
	LogRetentionRetryOptions *LogRetentionRetryOptions `json:"logRetentionRetryOptions"`
	// The IAM role for the Lambda function associated with the custom resource that sets the retention policy.
	// Experimental.
	LogRetentionRole awsiam.IRole `json:"logRetentionRole"`
	// The amount of memory, in MB, that is allocated to your Lambda function.
	//
	// Lambda uses this value to proportionally allocate the amount of CPU
	// power. For more information, see Resource Model in the AWS Lambda
	// Developer Guide.
	// Experimental.
	MemorySize *float64 `json:"memorySize"`
	// Enable profiling.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	Profiling *bool `json:"profiling"`
	// Profiling Group.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	ProfilingGroup awscodeguruprofiler.IProfilingGroup `json:"profilingGroup"`
	// The maximum of concurrent executions you want to reserve for the function.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/concurrent-executions.html
	//
	// Experimental.
	ReservedConcurrentExecutions *float64 `json:"reservedConcurrentExecutions"`
	// Lambda execution role.
	//
	// This is the role that will be assumed by the function upon execution.
	// It controls the permissions that the function will have. The Role must
	// be assumable by the 'lambda.amazonaws.com' service principal.
	//
	// The default Role automatically has permissions granted for Lambda execution. If you
	// provide a Role, you must add the relevant AWS managed policies yourself.
	//
	// The relevant managed policies are "service-role/AWSLambdaBasicExecutionRole" and
	// "service-role/AWSLambdaVPCAccessExecutionRole".
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// What security group to associate with the Lambda's network interfaces. This property is being deprecated, consider using securityGroups instead.
	//
	// Only used if 'vpc' is supplied.
	//
	// Use securityGroups property instead.
	// Function constructor will throw an error if both are specified.
	// Deprecated: - This property is deprecated, use securityGroups instead
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// The list of security groups to associate with the Lambda's network interfaces.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The function execution time (in seconds) after which Lambda terminates the function.
	//
	// Because the execution time affects cost, set this value
	// based on the function's expected execution time.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// Enable AWS X-Ray Tracing for Lambda Function.
	// Experimental.
	Tracing Tracing `json:"tracing"`
	// VPC network to place Lambda network interfaces.
	//
	// Specify this if the Lambda function needs to access resources in a VPC.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Where to place the network interfaces within the VPC.
	//
	// Only used if 'vpc' is supplied. Note: internet access for Lambdas
	// requires a NAT gateway, so picking Public subnets is not allowed.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
}

// Experimental.
type FunctionProps struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Whether to allow the Lambda to send all network traffic.
	//
	// If set to false, you must individually add traffic rules to allow the
	// Lambda to connect to network targets.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Lambda Functions in a public subnet can NOT access the internet.
	//
	// Use this property to acknowledge this limitation and still place the function in a public subnet.
	// See: https://stackoverflow.com/questions/52992085/why-cant-an-aws-lambda-function-inside-a-public-subnet-in-a-vpc-connect-to-the/52994841#52994841
	//
	// Experimental.
	AllowPublicSubnet *bool `json:"allowPublicSubnet"`
	// Code signing config associated with this function.
	// Experimental.
	CodeSigningConfig ICodeSigningConfig `json:"codeSigningConfig"`
	// Options for the `lambda.Version` resource automatically created by the `fn.currentVersion` method.
	// Experimental.
	CurrentVersionOptions *VersionOptions `json:"currentVersionOptions"`
	// The SQS queue to use if DLQ is enabled.
	// Experimental.
	DeadLetterQueue awssqs.IQueue `json:"deadLetterQueue"`
	// Enabled DLQ.
	//
	// If `deadLetterQueue` is undefined,
	// an SQS queue with default options will be defined for your Function.
	// Experimental.
	DeadLetterQueueEnabled *bool `json:"deadLetterQueueEnabled"`
	// A description of the function.
	// Experimental.
	Description *string `json:"description"`
	// Key-value pairs that Lambda caches and makes available for your Lambda functions.
	//
	// Use environment variables to apply configuration changes, such
	// as test and production environment configurations, without changing your
	// Lambda function source code.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The AWS KMS key that's used to encrypt your function's environment variables.
	// Experimental.
	EnvironmentEncryption awskms.IKey `json:"environmentEncryption"`
	// Event sources for this function.
	//
	// You can also add event sources using `addEventSource`.
	// Experimental.
	Events *[]IEventSource `json:"events"`
	// The filesystem configuration for the lambda function.
	// Experimental.
	Filesystem FileSystem `json:"filesystem"`
	// A name for the function.
	// Experimental.
	FunctionName *string `json:"functionName"`
	// Initial policy statements to add to the created Lambda Role.
	//
	// You can call `addToRolePolicy` to the created lambda to add statements post creation.
	// Experimental.
	InitialPolicy *[]awsiam.PolicyStatement `json:"initialPolicy"`
	// A list of layers to add to the function's execution environment.
	//
	// You can configure your Lambda function to pull in
	// additional code during initialization in the form of layers. Layers are packages of libraries or other dependencies
	// that can be used by multiple functions.
	// Experimental.
	Layers *[]ILayerVersion `json:"layers"`
	// The number of days log events are kept in CloudWatch Logs.
	//
	// When updating
	// this property, unsetting it doesn't remove the log retention policy. To
	// remove the retention policy, set the value to `INFINITE`.
	// Experimental.
	LogRetention awslogs.RetentionDays `json:"logRetention"`
	// When log retention is specified, a custom resource attempts to create the CloudWatch log group.
	//
	// These options control the retry policy when interacting with CloudWatch APIs.
	// Experimental.
	LogRetentionRetryOptions *LogRetentionRetryOptions `json:"logRetentionRetryOptions"`
	// The IAM role for the Lambda function associated with the custom resource that sets the retention policy.
	// Experimental.
	LogRetentionRole awsiam.IRole `json:"logRetentionRole"`
	// The amount of memory, in MB, that is allocated to your Lambda function.
	//
	// Lambda uses this value to proportionally allocate the amount of CPU
	// power. For more information, see Resource Model in the AWS Lambda
	// Developer Guide.
	// Experimental.
	MemorySize *float64 `json:"memorySize"`
	// Enable profiling.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	Profiling *bool `json:"profiling"`
	// Profiling Group.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	ProfilingGroup awscodeguruprofiler.IProfilingGroup `json:"profilingGroup"`
	// The maximum of concurrent executions you want to reserve for the function.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/concurrent-executions.html
	//
	// Experimental.
	ReservedConcurrentExecutions *float64 `json:"reservedConcurrentExecutions"`
	// Lambda execution role.
	//
	// This is the role that will be assumed by the function upon execution.
	// It controls the permissions that the function will have. The Role must
	// be assumable by the 'lambda.amazonaws.com' service principal.
	//
	// The default Role automatically has permissions granted for Lambda execution. If you
	// provide a Role, you must add the relevant AWS managed policies yourself.
	//
	// The relevant managed policies are "service-role/AWSLambdaBasicExecutionRole" and
	// "service-role/AWSLambdaVPCAccessExecutionRole".
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// What security group to associate with the Lambda's network interfaces. This property is being deprecated, consider using securityGroups instead.
	//
	// Only used if 'vpc' is supplied.
	//
	// Use securityGroups property instead.
	// Function constructor will throw an error if both are specified.
	// Deprecated: - This property is deprecated, use securityGroups instead
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// The list of security groups to associate with the Lambda's network interfaces.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The function execution time (in seconds) after which Lambda terminates the function.
	//
	// Because the execution time affects cost, set this value
	// based on the function's expected execution time.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// Enable AWS X-Ray Tracing for Lambda Function.
	// Experimental.
	Tracing Tracing `json:"tracing"`
	// VPC network to place Lambda network interfaces.
	//
	// Specify this if the Lambda function needs to access resources in a VPC.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Where to place the network interfaces within the VPC.
	//
	// Only used if 'vpc' is supplied. Note: internet access for Lambdas
	// requires a NAT gateway, so picking Public subnets is not allowed.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// The source code of your Lambda function.
	//
	// You can point to a file in an
	// Amazon Simple Storage Service (Amazon S3) bucket or specify your source
	// code as inline text.
	// Experimental.
	Code Code `json:"code"`
	// The name of the method within your code that Lambda calls to execute your function.
	//
	// The format includes the file name. It can also include
	// namespaces and other qualifiers, depending on the runtime.
	// For more information, see https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-features.html#gettingstarted-features-programmingmodel.
	//
	// Use `Handler.FROM_IMAGE` when defining a function from a Docker image.
	//
	// NOTE: If you specify your source code as inline text by specifying the
	// ZipFile property within the Code property, specify index.function_name as
	// the handler.
	// Experimental.
	Handler *string `json:"handler"`
	// The runtime environment for the Lambda function that you are uploading.
	//
	// For valid values, see the Runtime property in the AWS Lambda Developer
	// Guide.
	//
	// Use `Runtime.FROM_IMAGE` when when defining a function from a Docker image.
	// Experimental.
	Runtime Runtime `json:"runtime"`
}

// Lambda function handler.
// Experimental.
type Handler interface {
}

// The jsii proxy struct for Handler
type jsiiProxy_Handler struct {
	_ byte // padding
}

func Handler_FROM_IMAGE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Handler",
		"FROM_IMAGE",
		&returns,
	)
	return returns
}

// Experimental.
type IAlias interface {
	IFunction
	// Name of this alias.
	// Experimental.
	AliasName() *string
	// The underlying Lambda function version.
	// Experimental.
	Version() IVersion
}

// The jsii proxy for IAlias
type jsiiProxy_IAlias struct {
	jsiiProxy_IFunction
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

func (j *jsiiProxy_IAlias) Version() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}

// A Code Signing Config.
// Experimental.
type ICodeSigningConfig interface {
	awscdk.IResource
	// The ARN of Code Signing Config.
	// Experimental.
	CodeSigningConfigArn() *string
	// The id of Code Signing Config.
	// Experimental.
	CodeSigningConfigId() *string
}

// The jsii proxy for ICodeSigningConfig
type jsiiProxy_ICodeSigningConfig struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_ICodeSigningConfig) CodeSigningConfigArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codeSigningConfigArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICodeSigningConfig) CodeSigningConfigId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"codeSigningConfigId",
		&returns,
	)
	return returns
}

// A Lambda destination.
// Experimental.
type IDestination interface {
	// Binds this destination to the Lambda function.
	// Experimental.
	Bind(scope awscdk.Construct, fn IFunction, options *DestinationOptions) *DestinationConfig
}

// The jsii proxy for IDestination
type jsiiProxy_IDestination struct {
	_ byte // padding
}

func (i *jsiiProxy_IDestination) Bind(scope awscdk.Construct, fn IFunction, options *DestinationOptions) *DestinationConfig {
	var returns *DestinationConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, fn, options},
		&returns,
	)

	return returns
}

// An abstract class which represents an AWS Lambda event source.
// Experimental.
type IEventSource interface {
	// Called by `lambda.addEventSource` to allow the event source to bind to this function.
	// Experimental.
	Bind(target IFunction)
}

// The jsii proxy for IEventSource
type jsiiProxy_IEventSource struct {
	_ byte // padding
}

func (i *jsiiProxy_IEventSource) Bind(target IFunction) {
	_jsii_.InvokeVoid(
		i,
		"bind",
		[]interface{}{target},
	)
}

// A DLQ for an event source.
// Experimental.
type IEventSourceDlq interface {
	// Returns the DLQ destination config of the DLQ.
	// Experimental.
	Bind(target IEventSourceMapping, targetHandler IFunction) *DlqDestinationConfig
}

// The jsii proxy for IEventSourceDlq
type jsiiProxy_IEventSourceDlq struct {
	_ byte // padding
}

func (i *jsiiProxy_IEventSourceDlq) Bind(target IEventSourceMapping, targetHandler IFunction) *DlqDestinationConfig {
	var returns *DlqDestinationConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{target, targetHandler},
		&returns,
	)

	return returns
}

// Represents an event source mapping for a lambda function.
// See: https://docs.aws.amazon.com/lambda/latest/dg/invocation-eventsourcemapping.html
//
// Experimental.
type IEventSourceMapping interface {
	awscdk.IResource
	// The identifier for this EventSourceMapping.
	// Experimental.
	EventSourceMappingId() *string
}

// The jsii proxy for IEventSourceMapping
type jsiiProxy_IEventSourceMapping struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IEventSourceMapping) EventSourceMappingId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"eventSourceMappingId",
		&returns,
	)
	return returns
}

// Experimental.
type IFunction interface {
	awsec2.IConnectable
	awsiam.IGrantable
	awscdk.IResource
	// Adds an event source to this function.
	//
	// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
	//
	// The following example adds an SQS Queue as an event source:
	// ```
	// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
	// myFunction.addEventSource(new SqsEventSource(myQueue));
	// ```
	// Experimental.
	AddEventSource(source IEventSource)
	// Adds an event source that maps to this AWS Lambda function.
	// Experimental.
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	// Adds a permission to the Lambda resource policy.
	// See: Permission for details.
	//
	// Experimental.
	AddPermission(id *string, permission *Permission)
	// Adds a statement to the IAM role assumed by the instance.
	// Experimental.
	AddToRolePolicy(statement awsiam.PolicyStatement)
	// Configures options for asynchronous invocation.
	// Experimental.
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	// Grant the given identity permissions to invoke this Lambda.
	// Experimental.
	GrantInvoke(identity awsiam.IGrantable) awsiam.Grant
	// Return the given named metric for this Lambda Return the given named metric for this Function.
	// Experimental.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the Duration of this Lambda How long execution of this Lambda takes.
	//
	// Average over 5 minutes
	// Experimental.
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// How many invocations of this Lambda fail.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of invocations of this Lambda How often this Lambda is invoked.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Metric for the number of throttled invocations of this Lambda How often this Lambda is throttled.
	//
	// Sum over 5 minutes
	// Experimental.
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// The ARN fo the function.
	// Experimental.
	FunctionArn() *string
	// The name of the function.
	// Experimental.
	FunctionName() *string
	// Whether or not this Lambda function was bound to a VPC.
	//
	// If this is is `false`, trying to access the `connections` object will fail.
	// Experimental.
	IsBoundToVpc() *bool
	// The `$LATEST` version of this function.
	//
	// Note that this is reference to a non-specific AWS Lambda version, which
	// means the function this version refers to can return different results in
	// different invocations.
	//
	// To obtain a reference to an explicit version which references the current
	// function configuration, use `lambdaFunction.currentVersion` instead.
	// Experimental.
	LatestVersion() IVersion
	// The construct node where permissions are attached.
	// Experimental.
	PermissionsNode() awscdk.ConstructNode
	// The IAM role associated with this function.
	// Experimental.
	Role() awsiam.IRole
}

// The jsii proxy for IFunction
type jsiiProxy_IFunction struct {
	internal.Type__awsec2IConnectable
	internal.Type__awsiamIGrantable
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IFunction) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		i,
		"addEventSource",
		[]interface{}{source},
	)
}

func (i *jsiiProxy_IFunction) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		i,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IFunction) AddPermission(id *string, permission *Permission) {
	_jsii_.InvokeVoid(
		i,
		"addPermission",
		[]interface{}{id, permission},
	)
}

func (i *jsiiProxy_IFunction) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		i,
		"addToRolePolicy",
		[]interface{}{statement},
	)
}

func (i *jsiiProxy_IFunction) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		i,
		"configureAsyncInvoke",
		[]interface{}{options},
	)
}

func (i *jsiiProxy_IFunction) GrantInvoke(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantInvoke",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IFunction) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IFunction) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IFunction) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IFunction) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IFunction) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricThrottles",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IFunction) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IFunction) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// Experimental.
type ILayerVersion interface {
	awscdk.IResource
	// Add permission for this layer version to specific entities.
	//
	// Usage within
	// the same account where the layer is defined is always allowed and does not
	// require calling this method. Note that the principal that creates the
	// Lambda function using the layer (for example, a CloudFormation changeset
	// execution role) also needs to have the ``lambda:GetLayerVersion``
	// permission on the layer version.
	// Experimental.
	AddPermission(id *string, permission *LayerVersionPermission)
	// The runtimes compatible with this Layer.
	// Experimental.
	CompatibleRuntimes() *[]Runtime
	// The ARN of the Lambda Layer version that this Layer defines.
	// Experimental.
	LayerVersionArn() *string
}

// The jsii proxy for ILayerVersion
type jsiiProxy_ILayerVersion struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_ILayerVersion) AddPermission(id *string, permission *LayerVersionPermission) {
	_jsii_.InvokeVoid(
		i,
		"addPermission",
		[]interface{}{id, permission},
	)
}

func (j *jsiiProxy_ILayerVersion) CompatibleRuntimes() *[]Runtime {
	var returns *[]Runtime
	_jsii_.Get(
		j,
		"compatibleRuntimes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ILayerVersion) LayerVersionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"layerVersionArn",
		&returns,
	)
	return returns
}

// Interface for scalable attributes.
// Experimental.
type IScalableFunctionAttribute interface {
	awscdk.IConstruct
	// Scale out or in based on schedule.
	// Experimental.
	ScaleOnSchedule(id *string, actions *awsapplicationautoscaling.ScalingSchedule)
	// Scale out or in to keep utilization at a given level.
	//
	// The utilization is tracked by the
	// LambdaProvisionedConcurrencyUtilization metric, emitted by lambda. See:
	// https://docs.aws.amazon.com/lambda/latest/dg/monitoring-metrics.html#monitoring-metrics-concurrency
	// Experimental.
	ScaleOnUtilization(options *UtilizationScalingOptions)
}

// The jsii proxy for IScalableFunctionAttribute
type jsiiProxy_IScalableFunctionAttribute struct {
	internal.Type__awscdkIConstruct
}

func (i *jsiiProxy_IScalableFunctionAttribute) ScaleOnSchedule(id *string, actions *awsapplicationautoscaling.ScalingSchedule) {
	_jsii_.InvokeVoid(
		i,
		"scaleOnSchedule",
		[]interface{}{id, actions},
	)
}

func (i *jsiiProxy_IScalableFunctionAttribute) ScaleOnUtilization(options *UtilizationScalingOptions) {
	_jsii_.InvokeVoid(
		i,
		"scaleOnUtilization",
		[]interface{}{options},
	)
}

// Experimental.
type IVersion interface {
	IFunction
	// Defines an alias for this version.
	// Experimental.
	AddAlias(aliasName *string, options *AliasOptions) Alias
	// The ARN of the version for Lambda@Edge.
	// Experimental.
	EdgeArn() *string
	// The underlying AWS Lambda function.
	// Experimental.
	Lambda() IFunction
	// The most recently deployed version of this function.
	// Experimental.
	Version() *string
}

// The jsii proxy for IVersion
type jsiiProxy_IVersion struct {
	jsiiProxy_IFunction
}

func (i *jsiiProxy_IVersion) AddAlias(aliasName *string, options *AliasOptions) Alias {
	var returns Alias

	_jsii_.Invoke(
		i,
		"addAlias",
		[]interface{}{aliasName, options},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IVersion) EdgeArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"edgeArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IVersion) Lambda() IFunction {
	var returns IFunction
	_jsii_.Get(
		j,
		"lambda",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IVersion) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}

// Lambda code from an inline string (limited to 4KiB).
// Experimental.
type InlineCode interface {
	Code
	IsInline() *bool
	Bind(_scope awscdk.Construct) *CodeConfig
	BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions)
}

// The jsii proxy struct for InlineCode
type jsiiProxy_InlineCode struct {
	jsiiProxy_Code
}

func (j *jsiiProxy_InlineCode) IsInline() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isInline",
		&returns,
	)
	return returns
}


// Experimental.
func NewInlineCode(code *string) InlineCode {
	_init_.Initialize()

	j := jsiiProxy_InlineCode{}

	_jsii_.Create(
		"monocdk.aws_lambda.InlineCode",
		[]interface{}{code},
		&j,
	)

	return &j
}

// Experimental.
func NewInlineCode_Override(i InlineCode, code *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.InlineCode",
		[]interface{}{code},
		i,
	)
}

// DEPRECATED.
// Deprecated: use `fromAsset`
func InlineCode_Asset(path *string) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"asset",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromBucket`
func InlineCode_Bucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"bucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromCfnParameters`
func InlineCode_CfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"cfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from a local disk path.
// Experimental.
func InlineCode_FromAsset(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func InlineCode_FromAssetImage(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	var returns AssetImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"fromAssetImage",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Lambda handler code as an S3 object.
// Experimental.
func InlineCode_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Creates a new Lambda source defined using CloudFormation parameters.
//
// Returns: a new instance of `CfnParametersCode`
// Experimental.
func InlineCode_FromCfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"fromCfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from an asset created by a Docker build.
//
// By default, the asset is expected to be located at `/asset` in the
// image.
// Experimental.
func InlineCode_FromDockerBuild(path *string, options *DockerBuildAssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"fromDockerBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func InlineCode_FromEcrImage(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	var returns EcrImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"fromEcrImage",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Inline code for Lambda handler.
//
// Returns: `LambdaInlineCode` with inline code.
// Experimental.
func InlineCode_FromInline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"fromInline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromInline`
func InlineCode_Inline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.InlineCode",
		"inline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// Called when the lambda or layer is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (i *jsiiProxy_InlineCode) Bind(_scope awscdk.Construct) *CodeConfig {
	var returns *CodeConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// Called after the CFN function resource has been created to allow the code class to bind to it.
//
// Specifically it's required to allow assets to add
// metadata for tooling like SAM CLI to be able to find their origins.
// Experimental.
func (i *jsiiProxy_InlineCode) BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions) {
	_jsii_.InvokeVoid(
		i,
		"bindToResource",
		[]interface{}{_resource, _options},
	)
}

// Experimental.
type LambdaRuntimeProps struct {
	// The Docker image name to be used for bundling in this runtime.
	// Experimental.
	BundlingDockerImage *string `json:"bundlingDockerImage"`
	// Whether this runtime is integrated with and supported for profiling using Amazon CodeGuru Profiler.
	// Experimental.
	SupportsCodeGuruProfiling *bool `json:"supportsCodeGuruProfiling"`
	// Whether the ``ZipFile`` (aka inline code) property can be used with this runtime.
	// Experimental.
	SupportsInlineCode *bool `json:"supportsInlineCode"`
}

// Defines a new Lambda Layer version.
// Experimental.
type LayerVersion interface {
	awscdk.Resource
	ILayerVersion
	CompatibleRuntimes() *[]Runtime
	Env() *awscdk.ResourceEnvironment
	LayerVersionArn() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddPermission(id *string, permission *LayerVersionPermission)
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

// The jsii proxy struct for LayerVersion
type jsiiProxy_LayerVersion struct {
	internal.Type__awscdkResource
	jsiiProxy_ILayerVersion
}

func (j *jsiiProxy_LayerVersion) CompatibleRuntimes() *[]Runtime {
	var returns *[]Runtime
	_jsii_.Get(
		j,
		"compatibleRuntimes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LayerVersion) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LayerVersion) LayerVersionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"layerVersionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LayerVersion) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LayerVersion) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LayerVersion) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewLayerVersion(scope constructs.Construct, id *string, props *LayerVersionProps) LayerVersion {
	_init_.Initialize()

	j := jsiiProxy_LayerVersion{}

	_jsii_.Create(
		"monocdk.aws_lambda.LayerVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewLayerVersion_Override(l LayerVersion, scope constructs.Construct, id *string, props *LayerVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.LayerVersion",
		[]interface{}{scope, id, props},
		l,
	)
}

// Imports a layer version by ARN.
//
// Assumes it is compatible with all Lambda runtimes.
// Experimental.
func LayerVersion_FromLayerVersionArn(scope constructs.Construct, id *string, layerVersionArn *string) ILayerVersion {
	_init_.Initialize()

	var returns ILayerVersion

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.LayerVersion",
		"fromLayerVersionArn",
		[]interface{}{scope, id, layerVersionArn},
		&returns,
	)

	return returns
}

// Imports a Layer that has been defined externally.
// Experimental.
func LayerVersion_FromLayerVersionAttributes(scope constructs.Construct, id *string, attrs *LayerVersionAttributes) ILayerVersion {
	_init_.Initialize()

	var returns ILayerVersion

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.LayerVersion",
		"fromLayerVersionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func LayerVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.LayerVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func LayerVersion_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.LayerVersion",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add permission for this layer version to specific entities.
//
// Usage within
// the same account where the layer is defined is always allowed and does not
// require calling this method. Note that the principal that creates the
// Lambda function using the layer (for example, a CloudFormation changeset
// execution role) also needs to have the ``lambda:GetLayerVersion``
// permission on the layer version.
// Experimental.
func (l *jsiiProxy_LayerVersion) AddPermission(id *string, permission *LayerVersionPermission) {
	_jsii_.InvokeVoid(
		l,
		"addPermission",
		[]interface{}{id, permission},
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
func (l *jsiiProxy_LayerVersion) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		l,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (l *jsiiProxy_LayerVersion) GeneratePhysicalName() *string {
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
func (l *jsiiProxy_LayerVersion) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (l *jsiiProxy_LayerVersion) GetResourceNameAttribute(nameAttr *string) *string {
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
func (l *jsiiProxy_LayerVersion) OnPrepare() {
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
func (l *jsiiProxy_LayerVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (l *jsiiProxy_LayerVersion) OnValidate() *[]*string {
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
func (l *jsiiProxy_LayerVersion) Prepare() {
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
func (l *jsiiProxy_LayerVersion) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (l *jsiiProxy_LayerVersion) ToString() *string {
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
func (l *jsiiProxy_LayerVersion) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties necessary to import a LayerVersion.
// Experimental.
type LayerVersionAttributes struct {
	// The ARN of the LayerVersion.
	// Experimental.
	LayerVersionArn *string `json:"layerVersionArn"`
	// The list of compatible runtimes with this Layer.
	// Experimental.
	CompatibleRuntimes *[]Runtime `json:"compatibleRuntimes"`
}

// Non runtime options.
// Experimental.
type LayerVersionOptions struct {
	// The description the this Lambda Layer.
	// Experimental.
	Description *string `json:"description"`
	// The name of the layer.
	// Experimental.
	LayerVersionName *string `json:"layerVersionName"`
	// The SPDX licence identifier or URL to the license file for this layer.
	// Experimental.
	License *string `json:"license"`
	// Whether to retain this version of the layer when a new version is added or when the stack is deleted.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
}

// Identification of an account (or organization) that is allowed to access a Lambda Layer Version.
// Experimental.
type LayerVersionPermission struct {
	// The AWS Account id of the account that is authorized to use a Lambda Layer Version.
	//
	// The wild-card ``'*'`` can be
	// used to grant access to "any" account (or any account in an organization when ``organizationId`` is specified).
	// Experimental.
	AccountId *string `json:"accountId"`
	// The ID of the AWS Organization to which the grant is restricted.
	//
	// Can only be specified if ``accountId`` is ``'*'``
	// Experimental.
	OrganizationId *string `json:"organizationId"`
}

// Experimental.
type LayerVersionProps struct {
	// The description the this Lambda Layer.
	// Experimental.
	Description *string `json:"description"`
	// The name of the layer.
	// Experimental.
	LayerVersionName *string `json:"layerVersionName"`
	// The SPDX licence identifier or URL to the license file for this layer.
	// Experimental.
	License *string `json:"license"`
	// Whether to retain this version of the layer when a new version is added or when the stack is deleted.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// The content of this Layer.
	//
	// Using `Code.fromInline` is not supported.
	// Experimental.
	Code Code `json:"code"`
	// The runtimes compatible with this Layer.
	// Experimental.
	CompatibleRuntimes *[]Runtime `json:"compatibleRuntimes"`
}

// Creates a custom resource to control the retention policy of a CloudWatch Logs log group.
//
// The log group is created if it doesn't already exist. The policy
// is removed when `retentionDays` is `undefined` or equal to `Infinity`.
// Deprecated: use `LogRetention` from '
type LogRetention interface {
	awslogs.LogRetention
	LogGroupArn() *string
	Node() awscdk.ConstructNode
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for LogRetention
type jsiiProxy_LogRetention struct {
	internal.Type__awslogsLogRetention
}

func (j *jsiiProxy_LogRetention) LogGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LogRetention) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


// Deprecated: use `LogRetention` from '
func NewLogRetention(scope constructs.Construct, id *string, props *LogRetentionProps) LogRetention {
	_init_.Initialize()

	j := jsiiProxy_LogRetention{}

	_jsii_.Create(
		"monocdk.aws_lambda.LogRetention",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Deprecated: use `LogRetention` from '
func NewLogRetention_Override(l LogRetention, scope constructs.Construct, id *string, props *LogRetentionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.LogRetention",
		[]interface{}{scope, id, props},
		l,
	)
}

// Return whether the given object is a Construct.
// Deprecated: use `LogRetention` from '
func LogRetention_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.LogRetention",
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
// Deprecated: use `LogRetention` from '
func (l *jsiiProxy_LogRetention) OnPrepare() {
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
// Deprecated: use `LogRetention` from '
func (l *jsiiProxy_LogRetention) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: use `LogRetention` from '
func (l *jsiiProxy_LogRetention) OnValidate() *[]*string {
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
// Deprecated: use `LogRetention` from '
func (l *jsiiProxy_LogRetention) Prepare() {
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
// Deprecated: use `LogRetention` from '
func (l *jsiiProxy_LogRetention) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Deprecated: use `LogRetention` from '
func (l *jsiiProxy_LogRetention) ToString() *string {
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
// Deprecated: use `LogRetention` from '
func (l *jsiiProxy_LogRetention) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for a LogRetention.
// Deprecated: use `LogRetentionProps` from '
type LogRetentionProps struct {
	// The log group name.
	// Deprecated: use `LogRetentionProps` from '
	LogGroupName *string `json:"logGroupName"`
	// The number of days log events are kept in CloudWatch Logs.
	// Deprecated: use `LogRetentionProps` from '
	Retention awslogs.RetentionDays `json:"retention"`
	// The region where the log group should be created.
	// Deprecated: use `LogRetentionProps` from '
	LogGroupRegion *string `json:"logGroupRegion"`
	// Retry options for all AWS API calls.
	// Deprecated: use `LogRetentionProps` from '
	LogRetentionRetryOptions *awslogs.LogRetentionRetryOptions `json:"logRetentionRetryOptions"`
	// The IAM role for the Lambda function associated with the custom resource.
	// Deprecated: use `LogRetentionProps` from '
	Role awsiam.IRole `json:"role"`
}

// Retry options for all AWS API calls.
// Experimental.
type LogRetentionRetryOptions struct {
	// The base duration to use in the exponential backoff for operation retries.
	// Experimental.
	Base awscdk.Duration `json:"base"`
	// The maximum amount of retries.
	// Experimental.
	MaxRetries *float64 `json:"maxRetries"`
}

// Represents a permission statement that can be added to a Lambda's resource policy via the `addToResourcePolicy` method.
// Experimental.
type Permission struct {
	// The entity for which you are granting permission to invoke the Lambda function.
	//
	// This entity can be any valid AWS service principal, such as
	// s3.amazonaws.com or sns.amazonaws.com, or, if you are granting
	// cross-account permission, an AWS account ID. For example, you might want
	// to allow a custom application in another AWS account to push events to
	// Lambda by invoking your function.
	//
	// The principal can be either an AccountPrincipal or a ServicePrincipal.
	// Experimental.
	Principal awsiam.IPrincipal `json:"principal"`
	// The Lambda actions that you want to allow in this statement.
	//
	// For example,
	// you can specify lambda:CreateFunction to specify a certain action, or use
	// a wildcard (``lambda:*``) to grant permission to all Lambda actions. For a
	// list of actions, see Actions and Condition Context Keys for AWS Lambda in
	// the IAM User Guide.
	// Experimental.
	Action *string `json:"action"`
	// A unique token that must be supplied by the principal invoking the function.
	// Experimental.
	EventSourceToken *string `json:"eventSourceToken"`
	// The scope to which the permission constructs be attached.
	//
	// The default is
	// the Lambda function construct itself, but this would need to be different
	// in cases such as cross-stack references where the Permissions would need
	// to sit closer to the consumer of this permission (i.e., the caller).
	// Experimental.
	Scope awscdk.Construct `json:"scope"`
	// The AWS account ID (without hyphens) of the source owner.
	//
	// For example, if
	// you specify an S3 bucket in the SourceArn property, this value is the
	// bucket owner's account ID. You can use this property to ensure that all
	// source principals are owned by a specific account.
	// Experimental.
	SourceAccount *string `json:"sourceAccount"`
	// The ARN of a resource that is invoking your function.
	//
	// When granting
	// Amazon Simple Storage Service (Amazon S3) permission to invoke your
	// function, specify this property with the bucket ARN as its value. This
	// ensures that events generated only from the specified bucket, not just
	// any bucket from any AWS account that creates a mapping to your function,
	// can invoke the function.
	// Experimental.
	SourceArn *string `json:"sourceArn"`
}

// Experimental.
type QualifiedFunctionBase interface {
	FunctionBase
	CanCreatePermissions() *bool
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	FunctionArn() *string
	FunctionName() *string
	GrantPrincipal() awsiam.IPrincipal
	IsBoundToVpc() *bool
	Lambda() IFunction
	LatestVersion() IVersion
	Node() awscdk.ConstructNode
	PermissionsNode() awscdk.ConstructNode
	PhysicalName() *string
	Qualifier() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	AddEventSource(source IEventSource)
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	AddPermission(id *string, permission *Permission)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for QualifiedFunctionBase
type jsiiProxy_QualifiedFunctionBase struct {
	jsiiProxy_FunctionBase
}

func (j *jsiiProxy_QualifiedFunctionBase) CanCreatePermissions() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"canCreatePermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) Lambda() IFunction {
	var returns IFunction
	_jsii_.Get(
		j,
		"lambda",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) Qualifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"qualifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_QualifiedFunctionBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewQualifiedFunctionBase_Override(q QualifiedFunctionBase, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.QualifiedFunctionBase",
		[]interface{}{scope, id, props},
		q,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func QualifiedFunctionBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.QualifiedFunctionBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func QualifiedFunctionBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.QualifiedFunctionBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds an event source to this function.
//
// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
//
// The following example adds an SQS Queue as an event source:
// ```
// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
// myFunction.addEventSource(new SqsEventSource(myQueue));
// ```
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		q,
		"addEventSource",
		[]interface{}{source},
	)
}

// Adds an event source that maps to this AWS Lambda function.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		q,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a permission to the Lambda resource policy.
// See: Permission for details.
//
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) AddPermission(id *string, permission *Permission) {
	_jsii_.InvokeVoid(
		q,
		"addPermission",
		[]interface{}{id, permission},
	)
}

// Adds a statement to the IAM role assumed by the instance.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		q,
		"addToRolePolicy",
		[]interface{}{statement},
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
func (q *jsiiProxy_QualifiedFunctionBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		q,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Configures options for asynchronous invocation.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		q,
		"configureAsyncInvoke",
		[]interface{}{options},
	)
}

// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QualifiedFunctionBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QualifiedFunctionBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		q,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to invoke this Lambda.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		q,
		"grantInvoke",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Function.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// How long execution of this Lambda takes.
//
// Average over 5 minutes
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How many invocations of this Lambda fail.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is invoked.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is throttled.
//
// Sum over 5 minutes
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		q,
		"metricThrottles",
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
func (q *jsiiProxy_QualifiedFunctionBase) OnPrepare() {
	_jsii_.InvokeVoid(
		q,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
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
func (q *jsiiProxy_QualifiedFunctionBase) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QualifiedFunctionBase) Prepare() {
	_jsii_.InvokeVoid(
		q,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		q,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (q *jsiiProxy_QualifiedFunctionBase) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		q,
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
func (q *jsiiProxy_QualifiedFunctionBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		q,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type ResourceBindOptions struct {
	// The name of the CloudFormation property to annotate with asset metadata.
	// See: https://github.com/aws/aws-cdk/issues/1432
	//
	// Experimental.
	ResourceProperty *string `json:"resourceProperty"`
}

// Lambda function runtime environment.
//
// If you need to use a runtime name that doesn't exist as a static member, you
// can instantiate a `Runtime` object, e.g: `new Runtime('nodejs99.99')`.
// Experimental.
type Runtime interface {
	BundlingDockerImage() awscdk.BundlingDockerImage
	BundlingImage() awscdk.DockerImage
	Family() RuntimeFamily
	Name() *string
	SupportsCodeGuruProfiling() *bool
	SupportsInlineCode() *bool
	RuntimeEquals(other Runtime) *bool
	ToString() *string
}

// The jsii proxy struct for Runtime
type jsiiProxy_Runtime struct {
	_ byte // padding
}

func (j *jsiiProxy_Runtime) BundlingDockerImage() awscdk.BundlingDockerImage {
	var returns awscdk.BundlingDockerImage
	_jsii_.Get(
		j,
		"bundlingDockerImage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Runtime) BundlingImage() awscdk.DockerImage {
	var returns awscdk.DockerImage
	_jsii_.Get(
		j,
		"bundlingImage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Runtime) Family() RuntimeFamily {
	var returns RuntimeFamily
	_jsii_.Get(
		j,
		"family",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Runtime) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Runtime) SupportsCodeGuruProfiling() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"supportsCodeGuruProfiling",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Runtime) SupportsInlineCode() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"supportsInlineCode",
		&returns,
	)
	return returns
}


// Experimental.
func NewRuntime(name *string, family RuntimeFamily, props *LambdaRuntimeProps) Runtime {
	_init_.Initialize()

	j := jsiiProxy_Runtime{}

	_jsii_.Create(
		"monocdk.aws_lambda.Runtime",
		[]interface{}{name, family, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRuntime_Override(r Runtime, name *string, family RuntimeFamily, props *LambdaRuntimeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.Runtime",
		[]interface{}{name, family, props},
		r,
	)
}

func Runtime_ALL() *[]Runtime {
	_init_.Initialize()
	var returns *[]Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"ALL",
		&returns,
	)
	return returns
}

func Runtime_DOTNET_CORE_1() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"DOTNET_CORE_1",
		&returns,
	)
	return returns
}

func Runtime_DOTNET_CORE_2() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"DOTNET_CORE_2",
		&returns,
	)
	return returns
}

func Runtime_DOTNET_CORE_2_1() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"DOTNET_CORE_2_1",
		&returns,
	)
	return returns
}

func Runtime_DOTNET_CORE_3_1() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"DOTNET_CORE_3_1",
		&returns,
	)
	return returns
}

func Runtime_FROM_IMAGE() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"FROM_IMAGE",
		&returns,
	)
	return returns
}

func Runtime_GO_1_X() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"GO_1_X",
		&returns,
	)
	return returns
}

func Runtime_JAVA_11() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"JAVA_11",
		&returns,
	)
	return returns
}

func Runtime_JAVA_8() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"JAVA_8",
		&returns,
	)
	return returns
}

func Runtime_JAVA_8_CORRETTO() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"JAVA_8_CORRETTO",
		&returns,
	)
	return returns
}

func Runtime_NODEJS() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"NODEJS",
		&returns,
	)
	return returns
}

func Runtime_NODEJS_10_X() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"NODEJS_10_X",
		&returns,
	)
	return returns
}

func Runtime_NODEJS_12_X() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"NODEJS_12_X",
		&returns,
	)
	return returns
}

func Runtime_NODEJS_14_X() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"NODEJS_14_X",
		&returns,
	)
	return returns
}

func Runtime_NODEJS_4_3() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"NODEJS_4_3",
		&returns,
	)
	return returns
}

func Runtime_NODEJS_6_10() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"NODEJS_6_10",
		&returns,
	)
	return returns
}

func Runtime_NODEJS_8_10() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"NODEJS_8_10",
		&returns,
	)
	return returns
}

func Runtime_PROVIDED() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"PROVIDED",
		&returns,
	)
	return returns
}

func Runtime_PROVIDED_AL2() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"PROVIDED_AL2",
		&returns,
	)
	return returns
}

func Runtime_PYTHON_2_7() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"PYTHON_2_7",
		&returns,
	)
	return returns
}

func Runtime_PYTHON_3_6() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"PYTHON_3_6",
		&returns,
	)
	return returns
}

func Runtime_PYTHON_3_7() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"PYTHON_3_7",
		&returns,
	)
	return returns
}

func Runtime_PYTHON_3_8() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"PYTHON_3_8",
		&returns,
	)
	return returns
}

func Runtime_RUBY_2_5() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"RUBY_2_5",
		&returns,
	)
	return returns
}

func Runtime_RUBY_2_7() Runtime {
	_init_.Initialize()
	var returns Runtime
	_jsii_.StaticGet(
		"monocdk.aws_lambda.Runtime",
		"RUBY_2_7",
		&returns,
	)
	return returns
}

// Experimental.
func (r *jsiiProxy_Runtime) RuntimeEquals(other Runtime) *bool {
	var returns *bool

	_jsii_.Invoke(
		r,
		"runtimeEquals",
		[]interface{}{other},
		&returns,
	)

	return returns
}

// Experimental.
func (r *jsiiProxy_Runtime) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type RuntimeFamily string

const (
	RuntimeFamily_NODEJS RuntimeFamily = "NODEJS"
	RuntimeFamily_JAVA RuntimeFamily = "JAVA"
	RuntimeFamily_PYTHON RuntimeFamily = "PYTHON"
	RuntimeFamily_DOTNET_CORE RuntimeFamily = "DOTNET_CORE"
	RuntimeFamily_GO RuntimeFamily = "GO"
	RuntimeFamily_RUBY RuntimeFamily = "RUBY"
	RuntimeFamily_OTHER RuntimeFamily = "OTHER"
)

// Lambda code from an S3 archive.
// Experimental.
type S3Code interface {
	Code
	IsInline() *bool
	Bind(_scope awscdk.Construct) *CodeConfig
	BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions)
}

// The jsii proxy struct for S3Code
type jsiiProxy_S3Code struct {
	jsiiProxy_Code
}

func (j *jsiiProxy_S3Code) IsInline() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isInline",
		&returns,
	)
	return returns
}


// Experimental.
func NewS3Code(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	j := jsiiProxy_S3Code{}

	_jsii_.Create(
		"monocdk.aws_lambda.S3Code",
		[]interface{}{bucket, key, objectVersion},
		&j,
	)

	return &j
}

// Experimental.
func NewS3Code_Override(s S3Code, bucket awss3.IBucket, key *string, objectVersion *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.S3Code",
		[]interface{}{bucket, key, objectVersion},
		s,
	)
}

// DEPRECATED.
// Deprecated: use `fromAsset`
func S3Code_Asset(path *string) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"asset",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromBucket`
func S3Code_Bucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"bucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromCfnParameters`
func S3Code_CfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"cfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from a local disk path.
// Experimental.
func S3Code_FromAsset(path *string, options *awss3assets.AssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Create an ECR image from the specified asset and bind it as the Lambda code.
// Experimental.
func S3Code_FromAssetImage(directory *string, props *AssetImageCodeProps) AssetImageCode {
	_init_.Initialize()

	var returns AssetImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"fromAssetImage",
		[]interface{}{directory, props},
		&returns,
	)

	return returns
}

// Lambda handler code as an S3 object.
// Experimental.
func S3Code_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3Code {
	_init_.Initialize()

	var returns S3Code

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Creates a new Lambda source defined using CloudFormation parameters.
//
// Returns: a new instance of `CfnParametersCode`
// Experimental.
func S3Code_FromCfnParameters(props *CfnParametersCodeProps) CfnParametersCode {
	_init_.Initialize()

	var returns CfnParametersCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"fromCfnParameters",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Loads the function code from an asset created by a Docker build.
//
// By default, the asset is expected to be located at `/asset` in the
// image.
// Experimental.
func S3Code_FromDockerBuild(path *string, options *DockerBuildAssetOptions) AssetCode {
	_init_.Initialize()

	var returns AssetCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"fromDockerBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Use an existing ECR image as the Lambda code.
// Experimental.
func S3Code_FromEcrImage(repository awsecr.IRepository, props *EcrImageCodeProps) EcrImageCode {
	_init_.Initialize()

	var returns EcrImageCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"fromEcrImage",
		[]interface{}{repository, props},
		&returns,
	)

	return returns
}

// Inline code for Lambda handler.
//
// Returns: `LambdaInlineCode` with inline code.
// Experimental.
func S3Code_FromInline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"fromInline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// DEPRECATED.
// Deprecated: use `fromInline`
func S3Code_Inline(code *string) InlineCode {
	_init_.Initialize()

	var returns InlineCode

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.S3Code",
		"inline",
		[]interface{}{code},
		&returns,
	)

	return returns
}

// Called when the lambda or layer is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (s *jsiiProxy_S3Code) Bind(_scope awscdk.Construct) *CodeConfig {
	var returns *CodeConfig

	_jsii_.Invoke(
		s,
		"bind",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// Called after the CFN function resource has been created to allow the code class to bind to it.
//
// Specifically it's required to allow assets to add
// metadata for tooling like SAM CLI to be able to find their origins.
// Experimental.
func (s *jsiiProxy_S3Code) BindToResource(_resource awscdk.CfnResource, _options *ResourceBindOptions) {
	_jsii_.InvokeVoid(
		s,
		"bindToResource",
		[]interface{}{_resource, _options},
	)
}

// A Lambda that will only ever be added to a stack once.
//
// This construct is a way to guarantee that the lambda function will be guaranteed to be part of the stack,
// once and only once, irrespective of how many times the construct is declared to be part of the stack.
// This is guaranteed as long as the `uuid` property and the optional `lambdaPurpose` property stay the same
// whenever they're declared into the stack.
// Experimental.
type SingletonFunction interface {
	FunctionBase
	CanCreatePermissions() *bool
	Connections() awsec2.Connections
	CurrentVersion() Version
	Env() *awscdk.ResourceEnvironment
	FunctionArn() *string
	FunctionName() *string
	GrantPrincipal() awsiam.IPrincipal
	IsBoundToVpc() *bool
	LatestVersion() IVersion
	Node() awscdk.ConstructNode
	PermissionsNode() awscdk.ConstructNode
	PhysicalName() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	AddDependency(up ...awscdk.IDependable)
	AddEventSource(source IEventSource)
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	AddPermission(name *string, permission *Permission)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	DependOn(down awscdk.IConstruct)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for SingletonFunction
type jsiiProxy_SingletonFunction struct {
	jsiiProxy_FunctionBase
}

func (j *jsiiProxy_SingletonFunction) CanCreatePermissions() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"canCreatePermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) CurrentVersion() Version {
	var returns Version
	_jsii_.Get(
		j,
		"currentVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SingletonFunction) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewSingletonFunction(scope constructs.Construct, id *string, props *SingletonFunctionProps) SingletonFunction {
	_init_.Initialize()

	j := jsiiProxy_SingletonFunction{}

	_jsii_.Create(
		"monocdk.aws_lambda.SingletonFunction",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSingletonFunction_Override(s SingletonFunction, scope constructs.Construct, id *string, props *SingletonFunctionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.SingletonFunction",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func SingletonFunction_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.SingletonFunction",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func SingletonFunction_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.SingletonFunction",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Using node.addDependency() does not work on this method as the underlying lambda function is modeled as a singleton across the stack. Use this method instead to declare dependencies.
// Experimental.
func (s *jsiiProxy_SingletonFunction) AddDependency(up ...awscdk.IDependable) {
	args := []interface{}{}
	for _, a := range up {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		s,
		"addDependency",
		args,
	)
}

// Adds an event source to this function.
//
// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
//
// The following example adds an SQS Queue as an event source:
// ```
// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
// myFunction.addEventSource(new SqsEventSource(myQueue));
// ```
// Experimental.
func (s *jsiiProxy_SingletonFunction) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		s,
		"addEventSource",
		[]interface{}{source},
	)
}

// Adds an event source that maps to this AWS Lambda function.
// Experimental.
func (s *jsiiProxy_SingletonFunction) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		s,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a permission to the Lambda resource policy.
// Experimental.
func (s *jsiiProxy_SingletonFunction) AddPermission(name *string, permission *Permission) {
	_jsii_.InvokeVoid(
		s,
		"addPermission",
		[]interface{}{name, permission},
	)
}

// Adds a statement to the IAM role assumed by the instance.
// Experimental.
func (s *jsiiProxy_SingletonFunction) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		s,
		"addToRolePolicy",
		[]interface{}{statement},
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
func (s *jsiiProxy_SingletonFunction) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Configures options for asynchronous invocation.
// Experimental.
func (s *jsiiProxy_SingletonFunction) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		s,
		"configureAsyncInvoke",
		[]interface{}{options},
	)
}

// The SingletonFunction construct cannot be added as a dependency of another construct using node.addDependency(). Use this method instead to declare this as a dependency of another construct.
// Experimental.
func (s *jsiiProxy_SingletonFunction) DependOn(down awscdk.IConstruct) {
	_jsii_.InvokeVoid(
		s,
		"dependOn",
		[]interface{}{down},
	)
}

// Experimental.
func (s *jsiiProxy_SingletonFunction) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_SingletonFunction) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_SingletonFunction) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to invoke this Lambda.
// Experimental.
func (s *jsiiProxy_SingletonFunction) GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		s,
		"grantInvoke",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Function.
// Experimental.
func (s *jsiiProxy_SingletonFunction) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// How long execution of this Lambda takes.
//
// Average over 5 minutes
// Experimental.
func (s *jsiiProxy_SingletonFunction) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How many invocations of this Lambda fail.
//
// Sum over 5 minutes
// Experimental.
func (s *jsiiProxy_SingletonFunction) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is invoked.
//
// Sum over 5 minutes
// Experimental.
func (s *jsiiProxy_SingletonFunction) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is throttled.
//
// Sum over 5 minutes
// Experimental.
func (s *jsiiProxy_SingletonFunction) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricThrottles",
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
func (s *jsiiProxy_SingletonFunction) OnPrepare() {
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
func (s *jsiiProxy_SingletonFunction) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_SingletonFunction) OnValidate() *[]*string {
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
func (s *jsiiProxy_SingletonFunction) Prepare() {
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
func (s *jsiiProxy_SingletonFunction) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_SingletonFunction) ToString() *string {
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
func (s *jsiiProxy_SingletonFunction) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a newly created singleton Lambda.
// Experimental.
type SingletonFunctionProps struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// Whether to allow the Lambda to send all network traffic.
	//
	// If set to false, you must individually add traffic rules to allow the
	// Lambda to connect to network targets.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Lambda Functions in a public subnet can NOT access the internet.
	//
	// Use this property to acknowledge this limitation and still place the function in a public subnet.
	// See: https://stackoverflow.com/questions/52992085/why-cant-an-aws-lambda-function-inside-a-public-subnet-in-a-vpc-connect-to-the/52994841#52994841
	//
	// Experimental.
	AllowPublicSubnet *bool `json:"allowPublicSubnet"`
	// Code signing config associated with this function.
	// Experimental.
	CodeSigningConfig ICodeSigningConfig `json:"codeSigningConfig"`
	// Options for the `lambda.Version` resource automatically created by the `fn.currentVersion` method.
	// Experimental.
	CurrentVersionOptions *VersionOptions `json:"currentVersionOptions"`
	// The SQS queue to use if DLQ is enabled.
	// Experimental.
	DeadLetterQueue awssqs.IQueue `json:"deadLetterQueue"`
	// Enabled DLQ.
	//
	// If `deadLetterQueue` is undefined,
	// an SQS queue with default options will be defined for your Function.
	// Experimental.
	DeadLetterQueueEnabled *bool `json:"deadLetterQueueEnabled"`
	// A description of the function.
	// Experimental.
	Description *string `json:"description"`
	// Key-value pairs that Lambda caches and makes available for your Lambda functions.
	//
	// Use environment variables to apply configuration changes, such
	// as test and production environment configurations, without changing your
	// Lambda function source code.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The AWS KMS key that's used to encrypt your function's environment variables.
	// Experimental.
	EnvironmentEncryption awskms.IKey `json:"environmentEncryption"`
	// Event sources for this function.
	//
	// You can also add event sources using `addEventSource`.
	// Experimental.
	Events *[]IEventSource `json:"events"`
	// The filesystem configuration for the lambda function.
	// Experimental.
	Filesystem FileSystem `json:"filesystem"`
	// A name for the function.
	// Experimental.
	FunctionName *string `json:"functionName"`
	// Initial policy statements to add to the created Lambda Role.
	//
	// You can call `addToRolePolicy` to the created lambda to add statements post creation.
	// Experimental.
	InitialPolicy *[]awsiam.PolicyStatement `json:"initialPolicy"`
	// A list of layers to add to the function's execution environment.
	//
	// You can configure your Lambda function to pull in
	// additional code during initialization in the form of layers. Layers are packages of libraries or other dependencies
	// that can be used by multiple functions.
	// Experimental.
	Layers *[]ILayerVersion `json:"layers"`
	// The number of days log events are kept in CloudWatch Logs.
	//
	// When updating
	// this property, unsetting it doesn't remove the log retention policy. To
	// remove the retention policy, set the value to `INFINITE`.
	// Experimental.
	LogRetention awslogs.RetentionDays `json:"logRetention"`
	// When log retention is specified, a custom resource attempts to create the CloudWatch log group.
	//
	// These options control the retry policy when interacting with CloudWatch APIs.
	// Experimental.
	LogRetentionRetryOptions *LogRetentionRetryOptions `json:"logRetentionRetryOptions"`
	// The IAM role for the Lambda function associated with the custom resource that sets the retention policy.
	// Experimental.
	LogRetentionRole awsiam.IRole `json:"logRetentionRole"`
	// The amount of memory, in MB, that is allocated to your Lambda function.
	//
	// Lambda uses this value to proportionally allocate the amount of CPU
	// power. For more information, see Resource Model in the AWS Lambda
	// Developer Guide.
	// Experimental.
	MemorySize *float64 `json:"memorySize"`
	// Enable profiling.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	Profiling *bool `json:"profiling"`
	// Profiling Group.
	// See: https://docs.aws.amazon.com/codeguru/latest/profiler-ug/setting-up-lambda.html
	//
	// Experimental.
	ProfilingGroup awscodeguruprofiler.IProfilingGroup `json:"profilingGroup"`
	// The maximum of concurrent executions you want to reserve for the function.
	// See: https://docs.aws.amazon.com/lambda/latest/dg/concurrent-executions.html
	//
	// Experimental.
	ReservedConcurrentExecutions *float64 `json:"reservedConcurrentExecutions"`
	// Lambda execution role.
	//
	// This is the role that will be assumed by the function upon execution.
	// It controls the permissions that the function will have. The Role must
	// be assumable by the 'lambda.amazonaws.com' service principal.
	//
	// The default Role automatically has permissions granted for Lambda execution. If you
	// provide a Role, you must add the relevant AWS managed policies yourself.
	//
	// The relevant managed policies are "service-role/AWSLambdaBasicExecutionRole" and
	// "service-role/AWSLambdaVPCAccessExecutionRole".
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// What security group to associate with the Lambda's network interfaces. This property is being deprecated, consider using securityGroups instead.
	//
	// Only used if 'vpc' is supplied.
	//
	// Use securityGroups property instead.
	// Function constructor will throw an error if both are specified.
	// Deprecated: - This property is deprecated, use securityGroups instead
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// The list of security groups to associate with the Lambda's network interfaces.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// The function execution time (in seconds) after which Lambda terminates the function.
	//
	// Because the execution time affects cost, set this value
	// based on the function's expected execution time.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// Enable AWS X-Ray Tracing for Lambda Function.
	// Experimental.
	Tracing Tracing `json:"tracing"`
	// VPC network to place Lambda network interfaces.
	//
	// Specify this if the Lambda function needs to access resources in a VPC.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Where to place the network interfaces within the VPC.
	//
	// Only used if 'vpc' is supplied. Note: internet access for Lambdas
	// requires a NAT gateway, so picking Public subnets is not allowed.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// The source code of your Lambda function.
	//
	// You can point to a file in an
	// Amazon Simple Storage Service (Amazon S3) bucket or specify your source
	// code as inline text.
	// Experimental.
	Code Code `json:"code"`
	// The name of the method within your code that Lambda calls to execute your function.
	//
	// The format includes the file name. It can also include
	// namespaces and other qualifiers, depending on the runtime.
	// For more information, see https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-features.html#gettingstarted-features-programmingmodel.
	//
	// Use `Handler.FROM_IMAGE` when defining a function from a Docker image.
	//
	// NOTE: If you specify your source code as inline text by specifying the
	// ZipFile property within the Code property, specify index.function_name as
	// the handler.
	// Experimental.
	Handler *string `json:"handler"`
	// The runtime environment for the Lambda function that you are uploading.
	//
	// For valid values, see the Runtime property in the AWS Lambda Developer
	// Guide.
	//
	// Use `Runtime.FROM_IMAGE` when when defining a function from a Docker image.
	// Experimental.
	Runtime Runtime `json:"runtime"`
	// A unique identifier to identify this lambda.
	//
	// The identifier should be unique across all custom resource providers.
	// We recommend generating a UUID per provider.
	// Experimental.
	Uuid *string `json:"uuid"`
	// A descriptive name for the purpose of this Lambda.
	//
	// If the Lambda does not have a physical name, this string will be
	// reflected its generated name. The combination of lambdaPurpose
	// and uuid must be unique.
	// Experimental.
	LambdaPurpose *string `json:"lambdaPurpose"`
}

// Specific settings like the authentication protocol or the VPC components to secure access to your event source.
// Experimental.
type SourceAccessConfiguration struct {
	// The type of authentication protocol or the VPC components for your event source.
	//
	// For example: "SASL_SCRAM_512_AUTH".
	// Experimental.
	Type SourceAccessConfigurationType `json:"type"`
	// The value for your chosen configuration in type.
	//
	// For example: "URI": "arn:aws:secretsmanager:us-east-1:01234567890:secret:MyBrokerSecretName".
	// The exact string depends on the type.
	// See: SourceAccessConfigurationType
	//
	// Experimental.
	Uri *string `json:"uri"`
}

// The type of authentication protocol or the VPC components for your event source's SourceAccessConfiguration.
// See: https://docs.aws.amazon.com/lambda/latest/dg/API_SourceAccessConfiguration.html#SSS-Type-SourceAccessConfiguration-Type
//
// Experimental.
type SourceAccessConfigurationType interface {
	Type() *string
}

// The jsii proxy struct for SourceAccessConfigurationType
type jsiiProxy_SourceAccessConfigurationType struct {
	_ byte // padding
}

func (j *jsiiProxy_SourceAccessConfigurationType) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// A custom source access configuration property.
// Experimental.
func SourceAccessConfigurationType_Of(name *string) SourceAccessConfigurationType {
	_init_.Initialize()

	var returns SourceAccessConfigurationType

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.SourceAccessConfigurationType",
		"of",
		[]interface{}{name},
		&returns,
	)

	return returns
}

func SourceAccessConfigurationType_BASIC_AUTH() SourceAccessConfigurationType {
	_init_.Initialize()
	var returns SourceAccessConfigurationType
	_jsii_.StaticGet(
		"monocdk.aws_lambda.SourceAccessConfigurationType",
		"BASIC_AUTH",
		&returns,
	)
	return returns
}

func SourceAccessConfigurationType_SASL_SCRAM_256_AUTH() SourceAccessConfigurationType {
	_init_.Initialize()
	var returns SourceAccessConfigurationType
	_jsii_.StaticGet(
		"monocdk.aws_lambda.SourceAccessConfigurationType",
		"SASL_SCRAM_256_AUTH",
		&returns,
	)
	return returns
}

func SourceAccessConfigurationType_SASL_SCRAM_512_AUTH() SourceAccessConfigurationType {
	_init_.Initialize()
	var returns SourceAccessConfigurationType
	_jsii_.StaticGet(
		"monocdk.aws_lambda.SourceAccessConfigurationType",
		"SASL_SCRAM_512_AUTH",
		&returns,
	)
	return returns
}

func SourceAccessConfigurationType_VPC_SECURITY_GROUP() SourceAccessConfigurationType {
	_init_.Initialize()
	var returns SourceAccessConfigurationType
	_jsii_.StaticGet(
		"monocdk.aws_lambda.SourceAccessConfigurationType",
		"VPC_SECURITY_GROUP",
		&returns,
	)
	return returns
}

func SourceAccessConfigurationType_VPC_SUBNET() SourceAccessConfigurationType {
	_init_.Initialize()
	var returns SourceAccessConfigurationType
	_jsii_.StaticGet(
		"monocdk.aws_lambda.SourceAccessConfigurationType",
		"VPC_SUBNET",
		&returns,
	)
	return returns
}

// The position in the DynamoDB, Kinesis or MSK stream where AWS Lambda should start reading.
// Experimental.
type StartingPosition string

const (
	StartingPosition_TRIM_HORIZON StartingPosition = "TRIM_HORIZON"
	StartingPosition_LATEST StartingPosition = "LATEST"
)

// X-Ray Tracing Modes (https://docs.aws.amazon.com/lambda/latest/dg/API_TracingConfig.html).
// Experimental.
type Tracing string

const (
	Tracing_ACTIVE Tracing = "ACTIVE"
	Tracing_PASS_THROUGH Tracing = "PASS_THROUGH"
	Tracing_DISABLED Tracing = "DISABLED"
)

// Code signing configuration policy for deployment validation failure.
// Experimental.
type UntrustedArtifactOnDeployment string

const (
	UntrustedArtifactOnDeployment_ENFORCE UntrustedArtifactOnDeployment = "ENFORCE"
	UntrustedArtifactOnDeployment_WARN UntrustedArtifactOnDeployment = "WARN"
)

// Options for enabling Lambda utilization tracking.
// Experimental.
type UtilizationScalingOptions struct {
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
	// Utilization target for the attribute.
	//
	// For example, .5 indicates that 50 percent of allocated provisioned concurrency is in use.
	// Experimental.
	UtilizationTarget *float64 `json:"utilizationTarget"`
}

// A single newly-deployed version of a Lambda function.
//
// This object exists to--at deploy time--query the "then-current" version of
// the Lambda function that it refers to. This Version object can then be
// used in `Alias` to refer to a particular deployment of a Lambda.
//
// This means that for every new update you deploy to your Lambda (using the
// CDK and Aliases), you must always create a new Version object. In
// particular, it must have a different name, so that a new resource is
// created.
//
// If you want to ensure that you're associating the right version with
// the right deployment, specify the `codeSha256` property while
// creating the `Version.
// Experimental.
type Version interface {
	QualifiedFunctionBase
	IVersion
	CanCreatePermissions() *bool
	Connections() awsec2.Connections
	EdgeArn() *string
	Env() *awscdk.ResourceEnvironment
	FunctionArn() *string
	FunctionName() *string
	GrantPrincipal() awsiam.IPrincipal
	IsBoundToVpc() *bool
	Lambda() IFunction
	LatestVersion() IVersion
	Node() awscdk.ConstructNode
	PermissionsNode() awscdk.ConstructNode
	PhysicalName() *string
	Qualifier() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	Version() *string
	AddAlias(aliasName *string, options *AliasOptions) Alias
	AddEventSource(source IEventSource)
	AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping
	AddPermission(id *string, permission *Permission)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAsyncInvoke(options *EventInvokeConfigOptions)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Version
type jsiiProxy_Version struct {
	jsiiProxy_QualifiedFunctionBase
	jsiiProxy_IVersion
}

func (j *jsiiProxy_Version) CanCreatePermissions() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"canCreatePermissions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) EdgeArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"edgeArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) FunctionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) IsBoundToVpc() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isBoundToVpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Lambda() IFunction {
	var returns IFunction
	_jsii_.Get(
		j,
		"lambda",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) LatestVersion() IVersion {
	var returns IVersion
	_jsii_.Get(
		j,
		"latestVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) PermissionsNode() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"permissionsNode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Qualifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"qualifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Version) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


// Experimental.
func NewVersion(scope constructs.Construct, id *string, props *VersionProps) Version {
	_init_.Initialize()

	j := jsiiProxy_Version{}

	_jsii_.Create(
		"monocdk.aws_lambda.Version",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewVersion_Override(v Version, scope constructs.Construct, id *string, props *VersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_lambda.Version",
		[]interface{}{scope, id, props},
		v,
	)
}

// Construct a Version object from a Version ARN.
// Experimental.
func Version_FromVersionArn(scope constructs.Construct, id *string, versionArn *string) IVersion {
	_init_.Initialize()

	var returns IVersion

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Version",
		"fromVersionArn",
		[]interface{}{scope, id, versionArn},
		&returns,
	)

	return returns
}

// Experimental.
func Version_FromVersionAttributes(scope constructs.Construct, id *string, attrs *VersionAttributes) IVersion {
	_init_.Initialize()

	var returns IVersion

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Version",
		"fromVersionAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Version_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Version",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Version_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_lambda.Version",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Defines an alias for this version.
// Experimental.
func (v *jsiiProxy_Version) AddAlias(aliasName *string, options *AliasOptions) Alias {
	var returns Alias

	_jsii_.Invoke(
		v,
		"addAlias",
		[]interface{}{aliasName, options},
		&returns,
	)

	return returns
}

// Adds an event source to this function.
//
// Event sources are implemented in the @aws-cdk/aws-lambda-event-sources module.
//
// The following example adds an SQS Queue as an event source:
// ```
// import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
// myFunction.addEventSource(new SqsEventSource(myQueue));
// ```
// Experimental.
func (v *jsiiProxy_Version) AddEventSource(source IEventSource) {
	_jsii_.InvokeVoid(
		v,
		"addEventSource",
		[]interface{}{source},
	)
}

// Adds an event source that maps to this AWS Lambda function.
// Experimental.
func (v *jsiiProxy_Version) AddEventSourceMapping(id *string, options *EventSourceMappingOptions) EventSourceMapping {
	var returns EventSourceMapping

	_jsii_.Invoke(
		v,
		"addEventSourceMapping",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a permission to the Lambda resource policy.
// See: Permission for details.
//
// Experimental.
func (v *jsiiProxy_Version) AddPermission(id *string, permission *Permission) {
	_jsii_.InvokeVoid(
		v,
		"addPermission",
		[]interface{}{id, permission},
	)
}

// Adds a statement to the IAM role assumed by the instance.
// Experimental.
func (v *jsiiProxy_Version) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		v,
		"addToRolePolicy",
		[]interface{}{statement},
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
func (v *jsiiProxy_Version) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		v,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Configures options for asynchronous invocation.
// Experimental.
func (v *jsiiProxy_Version) ConfigureAsyncInvoke(options *EventInvokeConfigOptions) {
	_jsii_.InvokeVoid(
		v,
		"configureAsyncInvoke",
		[]interface{}{options},
	)
}

// Experimental.
func (v *jsiiProxy_Version) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_Version) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_Version) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		v,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to invoke this Lambda.
// Experimental.
func (v *jsiiProxy_Version) GrantInvoke(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		v,
		"grantInvoke",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Return the given named metric for this Function.
// Experimental.
func (v *jsiiProxy_Version) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		v,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// How long execution of this Lambda takes.
//
// Average over 5 minutes
// Experimental.
func (v *jsiiProxy_Version) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		v,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How many invocations of this Lambda fail.
//
// Sum over 5 minutes
// Experimental.
func (v *jsiiProxy_Version) MetricErrors(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		v,
		"metricErrors",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is invoked.
//
// Sum over 5 minutes
// Experimental.
func (v *jsiiProxy_Version) MetricInvocations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		v,
		"metricInvocations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// How often this Lambda is throttled.
//
// Sum over 5 minutes
// Experimental.
func (v *jsiiProxy_Version) MetricThrottles(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		v,
		"metricThrottles",
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
func (v *jsiiProxy_Version) OnPrepare() {
	_jsii_.InvokeVoid(
		v,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (v *jsiiProxy_Version) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		v,
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
func (v *jsiiProxy_Version) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_Version) Prepare() {
	_jsii_.InvokeVoid(
		v,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (v *jsiiProxy_Version) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		v,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (v *jsiiProxy_Version) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_Version) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		v,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type VersionAttributes struct {
	// The lambda function.
	// Experimental.
	Lambda IFunction `json:"lambda"`
	// The version.
	// Experimental.
	Version *string `json:"version"`
}

// Options for `lambda.Version`.
// Experimental.
type VersionOptions struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// SHA256 of the version of the Lambda source code.
	//
	// Specify to validate that you're deploying the right version.
	// Experimental.
	CodeSha256 *string `json:"codeSha256"`
	// Description of the version.
	// Experimental.
	Description *string `json:"description"`
	// Specifies a provisioned concurrency configuration for a function's version.
	// Experimental.
	ProvisionedConcurrentExecutions *float64 `json:"provisionedConcurrentExecutions"`
	// Whether to retain old versions of this function when a new version is created.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
}

// Properties for a new Lambda version.
// Experimental.
type VersionProps struct {
	// The maximum age of a request that Lambda sends to a function for processing.
	//
	// Minimum: 60 seconds
	// Maximum: 6 hours
	// Experimental.
	MaxEventAge awscdk.Duration `json:"maxEventAge"`
	// The destination for failed invocations.
	// Experimental.
	OnFailure IDestination `json:"onFailure"`
	// The destination for successful invocations.
	// Experimental.
	OnSuccess IDestination `json:"onSuccess"`
	// The maximum number of times to retry when the function returns an error.
	//
	// Minimum: 0
	// Maximum: 2
	// Experimental.
	RetryAttempts *float64 `json:"retryAttempts"`
	// SHA256 of the version of the Lambda source code.
	//
	// Specify to validate that you're deploying the right version.
	// Experimental.
	CodeSha256 *string `json:"codeSha256"`
	// Description of the version.
	// Experimental.
	Description *string `json:"description"`
	// Specifies a provisioned concurrency configuration for a function's version.
	// Experimental.
	ProvisionedConcurrentExecutions *float64 `json:"provisionedConcurrentExecutions"`
	// Whether to retain old versions of this function when a new version is created.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// Function to get the value of.
	// Experimental.
	Lambda IFunction `json:"lambda"`
}

// A version/weight pair for routing traffic to Lambda functions.
// Experimental.
type VersionWeight struct {
	// The version to route traffic to.
	// Experimental.
	Version IVersion `json:"version"`
	// How much weight to assign to this version (0..1).
	// Experimental.
	Weight *float64 `json:"weight"`
}


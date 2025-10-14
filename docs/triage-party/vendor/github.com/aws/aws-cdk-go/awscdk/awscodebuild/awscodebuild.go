package awscodebuild

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awscodebuild/internal"
	"github.com/aws/aws-cdk-go/awscdk/awscodecommit"
	"github.com/aws/aws-cdk-go/awscdk/awscodestarnotifications"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/aws-cdk-go/awscdk/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v3"
)

// Artifacts definition for a CodeBuild Project.
// Experimental.
type Artifacts interface {
	IArtifacts
	Identifier() *string
	Type() *string
	Bind(_scope awscdk.Construct, _project IProject) *ArtifactsConfig
}

// The jsii proxy struct for Artifacts
type jsiiProxy_Artifacts struct {
	jsiiProxy_IArtifacts
}

func (j *jsiiProxy_Artifacts) Identifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Artifacts) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Experimental.
func NewArtifacts_Override(a Artifacts, props *ArtifactsProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.Artifacts",
		[]interface{}{props},
		a,
	)
}

// Experimental.
func Artifacts_S3(props *S3ArtifactsProps) IArtifacts {
	_init_.Initialize()

	var returns IArtifacts

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Artifacts",
		"s3",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Callback when an Artifacts class is used in a CodeBuild Project.
// Experimental.
func (a *jsiiProxy_Artifacts) Bind(_scope awscdk.Construct, _project IProject) *ArtifactsConfig {
	var returns *ArtifactsConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{_scope, _project},
		&returns,
	)

	return returns
}

// The type returned from {@link IArtifacts#bind}.
// Experimental.
type ArtifactsConfig struct {
	// The low-level CloudFormation artifacts property.
	// Experimental.
	ArtifactsProperty *CfnProject_ArtifactsProperty `json:"artifactsProperty"`
}

// Properties common to all Artifacts classes.
// Experimental.
type ArtifactsProps struct {
	// The artifact identifier.
	//
	// This property is required on secondary artifacts.
	// Experimental.
	Identifier *string `json:"identifier"`
}

// The type returned from {@link IProject#enableBatchBuilds}.
// Experimental.
type BatchBuildConfig struct {
	// The IAM batch service Role of this Project.
	// Experimental.
	Role awsiam.IRole `json:"role"`
}

// The extra options passed to the {@link IProject.bindToCodePipeline} method.
// Experimental.
type BindToCodePipelineOptions struct {
	// The artifact bucket that will be used by the action that invokes this project.
	// Experimental.
	ArtifactBucket awss3.IBucket `json:"artifactBucket"`
}

// The source credentials used when contacting the BitBucket API.
//
// **Note**: CodeBuild only allows a single credential for BitBucket
// to be saved in a given AWS account in a given region -
// any attempt to add more than one will result in an error.
// Experimental.
type BitBucketSourceCredentials interface {
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

// The jsii proxy struct for BitBucketSourceCredentials
type jsiiProxy_BitBucketSourceCredentials struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_BitBucketSourceCredentials) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BitBucketSourceCredentials) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BitBucketSourceCredentials) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BitBucketSourceCredentials) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewBitBucketSourceCredentials(scope constructs.Construct, id *string, props *BitBucketSourceCredentialsProps) BitBucketSourceCredentials {
	_init_.Initialize()

	j := jsiiProxy_BitBucketSourceCredentials{}

	_jsii_.Create(
		"monocdk.aws_codebuild.BitBucketSourceCredentials",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewBitBucketSourceCredentials_Override(b BitBucketSourceCredentials, scope constructs.Construct, id *string, props *BitBucketSourceCredentialsProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.BitBucketSourceCredentials",
		[]interface{}{scope, id, props},
		b,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BitBucketSourceCredentials_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.BitBucketSourceCredentials",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func BitBucketSourceCredentials_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.BitBucketSourceCredentials",
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
func (b *jsiiProxy_BitBucketSourceCredentials) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (b *jsiiProxy_BitBucketSourceCredentials) GeneratePhysicalName() *string {
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
func (b *jsiiProxy_BitBucketSourceCredentials) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (b *jsiiProxy_BitBucketSourceCredentials) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
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
func (b *jsiiProxy_BitBucketSourceCredentials) OnPrepare() {
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
func (b *jsiiProxy_BitBucketSourceCredentials) OnSynthesize(session constructs.ISynthesisSession) {
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
func (b *jsiiProxy_BitBucketSourceCredentials) OnValidate() *[]*string {
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
func (b *jsiiProxy_BitBucketSourceCredentials) Prepare() {
	_jsii_.InvokeVoid(
		b,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_BitBucketSourceCredentials) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BitBucketSourceCredentials) ToString() *string {
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
func (b *jsiiProxy_BitBucketSourceCredentials) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties of {@link BitBucketSourceCredentials}.
// Experimental.
type BitBucketSourceCredentialsProps struct {
	// Your BitBucket application password.
	// Experimental.
	Password awscdk.SecretValue `json:"password"`
	// Your BitBucket username.
	// Experimental.
	Username awscdk.SecretValue `json:"username"`
}

// Construction properties for {@link BitBucketSource}.
// Experimental.
type BitBucketSourceProps struct {
	// The source identifier.
	//
	// This property is required on secondary sources.
	// Experimental.
	Identifier *string `json:"identifier"`
	// The BitBucket account/user that owns the repo.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Owner *string `json:"owner"`
	// The name of the repo (without the username).
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Repo *string `json:"repo"`
	// The commit ID, pull request ID, branch name, or tag name that corresponds to the version of the source code you want to build.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	BranchOrRef *string `json:"branchOrRef"`
	// The depth of history to download.
	//
	// Minimum value is 0.
	// If this value is 0, greater than 25, or not provided,
	// then the full history is downloaded with each build of the project.
	// Experimental.
	CloneDepth *float64 `json:"cloneDepth"`
	// Whether to fetch submodules while cloning git repo.
	// Experimental.
	FetchSubmodules *bool `json:"fetchSubmodules"`
	// Whether to send notifications on your build's start and end.
	// Experimental.
	ReportBuildStatus *bool `json:"reportBuildStatus"`
	// Whether to create a webhook that will trigger a build every time an event happens in the repository.
	// Experimental.
	Webhook *bool `json:"webhook"`
	// A list of webhook filters that can constraint what events in the repository will trigger a build.
	//
	// A build is triggered if any of the provided filter groups match.
	// Only valid if `webhook` was not provided as false.
	// Experimental.
	WebhookFilters *[]FilterGroup `json:"webhookFilters"`
	// Trigger a batch build from a webhook instead of a standard one.
	//
	// Enabling this will enable batch builds on the CodeBuild project.
	// Experimental.
	WebhookTriggersBatchBuild *bool `json:"webhookTriggersBatchBuild"`
}

// Experimental.
type BucketCacheOptions struct {
	// The prefix to use to store the cache in the bucket.
	// Experimental.
	Prefix *string `json:"prefix"`
}

// Experimental.
type BuildEnvironment struct {
	// The image used for the builds.
	// Experimental.
	BuildImage IBuildImage `json:"buildImage"`
	// The type of compute to use for this build.
	//
	// See the {@link ComputeType} enum for the possible values.
	// Experimental.
	ComputeType ComputeType `json:"computeType"`
	// The environment variables that your builds can use.
	// Experimental.
	EnvironmentVariables *map[string]*BuildEnvironmentVariable `json:"environmentVariables"`
	// Indicates how the project builds Docker images.
	//
	// Specify true to enable
	// running the Docker daemon inside a Docker container. This value must be
	// set to true only if this build project will be used to build Docker
	// images, and the specified build environment image is not one provided by
	// AWS CodeBuild with Docker support. Otherwise, all associated builds that
	// attempt to interact with the Docker daemon will fail.
	// Experimental.
	Privileged *bool `json:"privileged"`
}

// Experimental.
type BuildEnvironmentVariable struct {
	// The value of the environment variable.
	//
	// For plain-text variables (the default), this is the literal value of variable.
	// For SSM parameter variables, pass the name of the parameter here (`parameterName` property of `IParameter`).
	// For SecretsManager variables secrets, pass either the secret name (`secretName` property of `ISecret`)
	// or the secret ARN (`secretArn` property of `ISecret`) here,
	// along with optional SecretsManager qualifiers separated by ':', like the JSON key, or the version or stage
	// (see https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec.env.secrets-manager for details).
	// Experimental.
	Value interface{} `json:"value"`
	// The type of environment variable.
	// Experimental.
	Type BuildEnvironmentVariableType `json:"type"`
}

// Experimental.
type BuildEnvironmentVariableType string

const (
	BuildEnvironmentVariableType_PLAINTEXT BuildEnvironmentVariableType = "PLAINTEXT"
	BuildEnvironmentVariableType_PARAMETER_STORE BuildEnvironmentVariableType = "PARAMETER_STORE"
	BuildEnvironmentVariableType_SECRETS_MANAGER BuildEnvironmentVariableType = "SECRETS_MANAGER"
)

// Optional arguments to {@link IBuildImage.binder} - currently empty.
// Experimental.
type BuildImageBindOptions struct {
}

// The return type from {@link IBuildImage.binder} - currently empty.
// Experimental.
type BuildImageConfig struct {
}

// BuildSpec for CodeBuild projects.
// Experimental.
type BuildSpec interface {
	IsImmediate() *bool
	ToBuildSpec() *string
}

// The jsii proxy struct for BuildSpec
type jsiiProxy_BuildSpec struct {
	_ byte // padding
}

func (j *jsiiProxy_BuildSpec) IsImmediate() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isImmediate",
		&returns,
	)
	return returns
}


// Experimental.
func NewBuildSpec_Override(b BuildSpec) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.BuildSpec",
		nil, // no parameters
		b,
	)
}

// Experimental.
func BuildSpec_FromObject(value *map[string]interface{}) BuildSpec {
	_init_.Initialize()

	var returns BuildSpec

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.BuildSpec",
		"fromObject",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// Create a buildspec from an object that will be rendered as YAML in the resulting CloudFormation template.
// Experimental.
func BuildSpec_FromObjectToYaml(value *map[string]interface{}) BuildSpec {
	_init_.Initialize()

	var returns BuildSpec

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.BuildSpec",
		"fromObjectToYaml",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// Use a file from the source as buildspec.
//
// Use this if you want to use a file different from 'buildspec.yml'`
// Experimental.
func BuildSpec_FromSourceFilename(filename *string) BuildSpec {
	_init_.Initialize()

	var returns BuildSpec

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.BuildSpec",
		"fromSourceFilename",
		[]interface{}{filename},
		&returns,
	)

	return returns
}

// Render the represented BuildSpec.
// Experimental.
func (b *jsiiProxy_BuildSpec) ToBuildSpec() *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"toBuildSpec",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Cache options for CodeBuild Project.
//
// A cache can store reusable pieces of your build environment and use them across multiple builds.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-caching.html
//
// Experimental.
type Cache interface {
}

// The jsii proxy struct for Cache
type jsiiProxy_Cache struct {
	_ byte // padding
}

// Experimental.
func NewCache_Override(c Cache) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.Cache",
		nil, // no parameters
		c,
	)
}

// Create an S3 caching strategy.
// Experimental.
func Cache_Bucket(bucket awss3.IBucket, options *BucketCacheOptions) Cache {
	_init_.Initialize()

	var returns Cache

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Cache",
		"bucket",
		[]interface{}{bucket, options},
		&returns,
	)

	return returns
}

// Create a local caching strategy.
// Experimental.
func Cache_Local(modes ...LocalCacheMode) Cache {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range modes {
		args = append(args, a)
	}

	var returns Cache

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Cache",
		"local",
		args,
		&returns,
	)

	return returns
}

// Experimental.
func Cache_None() Cache {
	_init_.Initialize()

	var returns Cache

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Cache",
		"none",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A CloudFormation `AWS::CodeBuild::Project`.
type CfnProject interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Artifacts() interface{}
	SetArtifacts(val interface{})
	AttrArn() *string
	BadgeEnabled() interface{}
	SetBadgeEnabled(val interface{})
	BuildBatchConfig() interface{}
	SetBuildBatchConfig(val interface{})
	Cache() interface{}
	SetCache(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ConcurrentBuildLimit() *float64
	SetConcurrentBuildLimit(val *float64)
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	EncryptionKey() *string
	SetEncryptionKey(val *string)
	Environment() interface{}
	SetEnvironment(val interface{})
	FileSystemLocations() interface{}
	SetFileSystemLocations(val interface{})
	LogicalId() *string
	LogsConfig() interface{}
	SetLogsConfig(val interface{})
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	QueuedTimeoutInMinutes() *float64
	SetQueuedTimeoutInMinutes(val *float64)
	Ref() *string
	SecondaryArtifacts() interface{}
	SetSecondaryArtifacts(val interface{})
	SecondarySources() interface{}
	SetSecondarySources(val interface{})
	SecondarySourceVersions() interface{}
	SetSecondarySourceVersions(val interface{})
	ServiceRole() *string
	SetServiceRole(val *string)
	Source() interface{}
	SetSource(val interface{})
	SourceVersion() *string
	SetSourceVersion(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	TimeoutInMinutes() *float64
	SetTimeoutInMinutes(val *float64)
	Triggers() interface{}
	SetTriggers(val interface{})
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

// The jsii proxy struct for CfnProject
type jsiiProxy_CfnProject struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnProject) Artifacts() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"artifacts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) BadgeEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"badgeEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) BuildBatchConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"buildBatchConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Cache() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"cache",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) ConcurrentBuildLimit() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"concurrentBuildLimit",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) EncryptionKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Environment() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"environment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) FileSystemLocations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"fileSystemLocations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) LogsConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"logsConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) QueuedTimeoutInMinutes() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"queuedTimeoutInMinutes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) SecondaryArtifacts() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"secondaryArtifacts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) SecondarySources() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"secondarySources",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) SecondarySourceVersions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"secondarySourceVersions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) ServiceRole() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceRole",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Source() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) SourceVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) TimeoutInMinutes() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"timeoutInMinutes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) Triggers() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"triggers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnProject) VpcConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"vpcConfig",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodeBuild::Project`.
func NewCfnProject(scope awscdk.Construct, id *string, props *CfnProjectProps) CfnProject {
	_init_.Initialize()

	j := jsiiProxy_CfnProject{}

	_jsii_.Create(
		"monocdk.aws_codebuild.CfnProject",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodeBuild::Project`.
func NewCfnProject_Override(c CfnProject, scope awscdk.Construct, id *string, props *CfnProjectProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.CfnProject",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnProject) SetArtifacts(val interface{}) {
	_jsii_.Set(
		j,
		"artifacts",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetBadgeEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"badgeEnabled",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetBuildBatchConfig(val interface{}) {
	_jsii_.Set(
		j,
		"buildBatchConfig",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetCache(val interface{}) {
	_jsii_.Set(
		j,
		"cache",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetConcurrentBuildLimit(val *float64) {
	_jsii_.Set(
		j,
		"concurrentBuildLimit",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetEncryptionKey(val *string) {
	_jsii_.Set(
		j,
		"encryptionKey",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetEnvironment(val interface{}) {
	_jsii_.Set(
		j,
		"environment",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetFileSystemLocations(val interface{}) {
	_jsii_.Set(
		j,
		"fileSystemLocations",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetLogsConfig(val interface{}) {
	_jsii_.Set(
		j,
		"logsConfig",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetQueuedTimeoutInMinutes(val *float64) {
	_jsii_.Set(
		j,
		"queuedTimeoutInMinutes",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetSecondaryArtifacts(val interface{}) {
	_jsii_.Set(
		j,
		"secondaryArtifacts",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetSecondarySources(val interface{}) {
	_jsii_.Set(
		j,
		"secondarySources",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetSecondarySourceVersions(val interface{}) {
	_jsii_.Set(
		j,
		"secondarySourceVersions",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetServiceRole(val *string) {
	_jsii_.Set(
		j,
		"serviceRole",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetSource(val interface{}) {
	_jsii_.Set(
		j,
		"source",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetSourceVersion(val *string) {
	_jsii_.Set(
		j,
		"sourceVersion",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetTimeoutInMinutes(val *float64) {
	_jsii_.Set(
		j,
		"timeoutInMinutes",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetTriggers(val interface{}) {
	_jsii_.Set(
		j,
		"triggers",
		val,
	)
}

func (j *jsiiProxy_CfnProject) SetVpcConfig(val interface{}) {
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
func CfnProject_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnProject",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnProject_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnProject",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnProject_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnProject",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnProject_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.CfnProject",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnProject) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnProject) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnProject) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnProject) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnProject) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnProject) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnProject) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnProject) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnProject) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnProject) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnProject) OnPrepare() {
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
func (c *jsiiProxy_CfnProject) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnProject) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnProject) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnProject) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnProject) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnProject) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnProject) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnProject) ToString() *string {
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
func (c *jsiiProxy_CfnProject) Validate() *[]*string {
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
func (c *jsiiProxy_CfnProject) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnProject_ArtifactsProperty struct {
	// `CfnProject.ArtifactsProperty.Type`.
	Type *string `json:"type"`
	// `CfnProject.ArtifactsProperty.ArtifactIdentifier`.
	ArtifactIdentifier *string `json:"artifactIdentifier"`
	// `CfnProject.ArtifactsProperty.EncryptionDisabled`.
	EncryptionDisabled interface{} `json:"encryptionDisabled"`
	// `CfnProject.ArtifactsProperty.Location`.
	Location *string `json:"location"`
	// `CfnProject.ArtifactsProperty.Name`.
	Name *string `json:"name"`
	// `CfnProject.ArtifactsProperty.NamespaceType`.
	NamespaceType *string `json:"namespaceType"`
	// `CfnProject.ArtifactsProperty.OverrideArtifactName`.
	OverrideArtifactName interface{} `json:"overrideArtifactName"`
	// `CfnProject.ArtifactsProperty.Packaging`.
	Packaging *string `json:"packaging"`
	// `CfnProject.ArtifactsProperty.Path`.
	Path *string `json:"path"`
}

type CfnProject_BatchRestrictionsProperty struct {
	// `CfnProject.BatchRestrictionsProperty.ComputeTypesAllowed`.
	ComputeTypesAllowed *[]*string `json:"computeTypesAllowed"`
	// `CfnProject.BatchRestrictionsProperty.MaximumBuildsAllowed`.
	MaximumBuildsAllowed *float64 `json:"maximumBuildsAllowed"`
}

type CfnProject_BuildStatusConfigProperty struct {
	// `CfnProject.BuildStatusConfigProperty.Context`.
	Context *string `json:"context"`
	// `CfnProject.BuildStatusConfigProperty.TargetUrl`.
	TargetUrl *string `json:"targetUrl"`
}

type CfnProject_CloudWatchLogsConfigProperty struct {
	// `CfnProject.CloudWatchLogsConfigProperty.Status`.
	Status *string `json:"status"`
	// `CfnProject.CloudWatchLogsConfigProperty.GroupName`.
	GroupName *string `json:"groupName"`
	// `CfnProject.CloudWatchLogsConfigProperty.StreamName`.
	StreamName *string `json:"streamName"`
}

type CfnProject_EnvironmentProperty struct {
	// `CfnProject.EnvironmentProperty.ComputeType`.
	ComputeType *string `json:"computeType"`
	// `CfnProject.EnvironmentProperty.Image`.
	Image *string `json:"image"`
	// `CfnProject.EnvironmentProperty.Type`.
	Type *string `json:"type"`
	// `CfnProject.EnvironmentProperty.Certificate`.
	Certificate *string `json:"certificate"`
	// `CfnProject.EnvironmentProperty.EnvironmentVariables`.
	EnvironmentVariables interface{} `json:"environmentVariables"`
	// `CfnProject.EnvironmentProperty.ImagePullCredentialsType`.
	ImagePullCredentialsType *string `json:"imagePullCredentialsType"`
	// `CfnProject.EnvironmentProperty.PrivilegedMode`.
	PrivilegedMode interface{} `json:"privilegedMode"`
	// `CfnProject.EnvironmentProperty.RegistryCredential`.
	RegistryCredential interface{} `json:"registryCredential"`
}

type CfnProject_EnvironmentVariableProperty struct {
	// `CfnProject.EnvironmentVariableProperty.Name`.
	Name *string `json:"name"`
	// `CfnProject.EnvironmentVariableProperty.Value`.
	Value *string `json:"value"`
	// `CfnProject.EnvironmentVariableProperty.Type`.
	Type *string `json:"type"`
}

type CfnProject_GitSubmodulesConfigProperty struct {
	// `CfnProject.GitSubmodulesConfigProperty.FetchSubmodules`.
	FetchSubmodules interface{} `json:"fetchSubmodules"`
}

type CfnProject_LogsConfigProperty struct {
	// `CfnProject.LogsConfigProperty.CloudWatchLogs`.
	CloudWatchLogs interface{} `json:"cloudWatchLogs"`
	// `CfnProject.LogsConfigProperty.S3Logs`.
	S3Logs interface{} `json:"s3Logs"`
}

type CfnProject_ProjectBuildBatchConfigProperty struct {
	// `CfnProject.ProjectBuildBatchConfigProperty.CombineArtifacts`.
	CombineArtifacts interface{} `json:"combineArtifacts"`
	// `CfnProject.ProjectBuildBatchConfigProperty.Restrictions`.
	Restrictions interface{} `json:"restrictions"`
	// `CfnProject.ProjectBuildBatchConfigProperty.ServiceRole`.
	ServiceRole *string `json:"serviceRole"`
	// `CfnProject.ProjectBuildBatchConfigProperty.TimeoutInMins`.
	TimeoutInMins *float64 `json:"timeoutInMins"`
}

type CfnProject_ProjectCacheProperty struct {
	// `CfnProject.ProjectCacheProperty.Type`.
	Type *string `json:"type"`
	// `CfnProject.ProjectCacheProperty.Location`.
	Location *string `json:"location"`
	// `CfnProject.ProjectCacheProperty.Modes`.
	Modes *[]*string `json:"modes"`
}

type CfnProject_ProjectFileSystemLocationProperty struct {
	// `CfnProject.ProjectFileSystemLocationProperty.Identifier`.
	Identifier *string `json:"identifier"`
	// `CfnProject.ProjectFileSystemLocationProperty.Location`.
	Location *string `json:"location"`
	// `CfnProject.ProjectFileSystemLocationProperty.MountPoint`.
	MountPoint *string `json:"mountPoint"`
	// `CfnProject.ProjectFileSystemLocationProperty.Type`.
	Type *string `json:"type"`
	// `CfnProject.ProjectFileSystemLocationProperty.MountOptions`.
	MountOptions *string `json:"mountOptions"`
}

type CfnProject_ProjectSourceVersionProperty struct {
	// `CfnProject.ProjectSourceVersionProperty.SourceIdentifier`.
	SourceIdentifier *string `json:"sourceIdentifier"`
	// `CfnProject.ProjectSourceVersionProperty.SourceVersion`.
	SourceVersion *string `json:"sourceVersion"`
}

type CfnProject_ProjectTriggersProperty struct {
	// `CfnProject.ProjectTriggersProperty.BuildType`.
	BuildType *string `json:"buildType"`
	// `CfnProject.ProjectTriggersProperty.FilterGroups`.
	FilterGroups interface{} `json:"filterGroups"`
	// `CfnProject.ProjectTriggersProperty.Webhook`.
	Webhook interface{} `json:"webhook"`
}

type CfnProject_RegistryCredentialProperty struct {
	// `CfnProject.RegistryCredentialProperty.Credential`.
	Credential *string `json:"credential"`
	// `CfnProject.RegistryCredentialProperty.CredentialProvider`.
	CredentialProvider *string `json:"credentialProvider"`
}

type CfnProject_S3LogsConfigProperty struct {
	// `CfnProject.S3LogsConfigProperty.Status`.
	Status *string `json:"status"`
	// `CfnProject.S3LogsConfigProperty.EncryptionDisabled`.
	EncryptionDisabled interface{} `json:"encryptionDisabled"`
	// `CfnProject.S3LogsConfigProperty.Location`.
	Location *string `json:"location"`
}

type CfnProject_SourceAuthProperty struct {
	// `CfnProject.SourceAuthProperty.Type`.
	Type *string `json:"type"`
	// `CfnProject.SourceAuthProperty.Resource`.
	Resource *string `json:"resource"`
}

type CfnProject_SourceProperty struct {
	// `CfnProject.SourceProperty.Type`.
	Type *string `json:"type"`
	// `CfnProject.SourceProperty.Auth`.
	Auth interface{} `json:"auth"`
	// `CfnProject.SourceProperty.BuildSpec`.
	BuildSpec *string `json:"buildSpec"`
	// `CfnProject.SourceProperty.BuildStatusConfig`.
	BuildStatusConfig interface{} `json:"buildStatusConfig"`
	// `CfnProject.SourceProperty.GitCloneDepth`.
	GitCloneDepth *float64 `json:"gitCloneDepth"`
	// `CfnProject.SourceProperty.GitSubmodulesConfig`.
	GitSubmodulesConfig interface{} `json:"gitSubmodulesConfig"`
	// `CfnProject.SourceProperty.InsecureSsl`.
	InsecureSsl interface{} `json:"insecureSsl"`
	// `CfnProject.SourceProperty.Location`.
	Location *string `json:"location"`
	// `CfnProject.SourceProperty.ReportBuildStatus`.
	ReportBuildStatus interface{} `json:"reportBuildStatus"`
	// `CfnProject.SourceProperty.SourceIdentifier`.
	SourceIdentifier *string `json:"sourceIdentifier"`
}

type CfnProject_VpcConfigProperty struct {
	// `CfnProject.VpcConfigProperty.SecurityGroupIds`.
	SecurityGroupIds *[]*string `json:"securityGroupIds"`
	// `CfnProject.VpcConfigProperty.Subnets`.
	Subnets *[]*string `json:"subnets"`
	// `CfnProject.VpcConfigProperty.VpcId`.
	VpcId *string `json:"vpcId"`
}

type CfnProject_WebhookFilterProperty struct {
	// `CfnProject.WebhookFilterProperty.Pattern`.
	Pattern *string `json:"pattern"`
	// `CfnProject.WebhookFilterProperty.Type`.
	Type *string `json:"type"`
	// `CfnProject.WebhookFilterProperty.ExcludeMatchedPattern`.
	ExcludeMatchedPattern interface{} `json:"excludeMatchedPattern"`
}

// Properties for defining a `AWS::CodeBuild::Project`.
type CfnProjectProps struct {
	// `AWS::CodeBuild::Project.Artifacts`.
	Artifacts interface{} `json:"artifacts"`
	// `AWS::CodeBuild::Project.Environment`.
	Environment interface{} `json:"environment"`
	// `AWS::CodeBuild::Project.ServiceRole`.
	ServiceRole *string `json:"serviceRole"`
	// `AWS::CodeBuild::Project.Source`.
	Source interface{} `json:"source"`
	// `AWS::CodeBuild::Project.BadgeEnabled`.
	BadgeEnabled interface{} `json:"badgeEnabled"`
	// `AWS::CodeBuild::Project.BuildBatchConfig`.
	BuildBatchConfig interface{} `json:"buildBatchConfig"`
	// `AWS::CodeBuild::Project.Cache`.
	Cache interface{} `json:"cache"`
	// `AWS::CodeBuild::Project.ConcurrentBuildLimit`.
	ConcurrentBuildLimit *float64 `json:"concurrentBuildLimit"`
	// `AWS::CodeBuild::Project.Description`.
	Description *string `json:"description"`
	// `AWS::CodeBuild::Project.EncryptionKey`.
	EncryptionKey *string `json:"encryptionKey"`
	// `AWS::CodeBuild::Project.FileSystemLocations`.
	FileSystemLocations interface{} `json:"fileSystemLocations"`
	// `AWS::CodeBuild::Project.LogsConfig`.
	LogsConfig interface{} `json:"logsConfig"`
	// `AWS::CodeBuild::Project.Name`.
	Name *string `json:"name"`
	// `AWS::CodeBuild::Project.QueuedTimeoutInMinutes`.
	QueuedTimeoutInMinutes *float64 `json:"queuedTimeoutInMinutes"`
	// `AWS::CodeBuild::Project.SecondaryArtifacts`.
	SecondaryArtifacts interface{} `json:"secondaryArtifacts"`
	// `AWS::CodeBuild::Project.SecondarySources`.
	SecondarySources interface{} `json:"secondarySources"`
	// `AWS::CodeBuild::Project.SecondarySourceVersions`.
	SecondarySourceVersions interface{} `json:"secondarySourceVersions"`
	// `AWS::CodeBuild::Project.SourceVersion`.
	SourceVersion *string `json:"sourceVersion"`
	// `AWS::CodeBuild::Project.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::CodeBuild::Project.TimeoutInMinutes`.
	TimeoutInMinutes *float64 `json:"timeoutInMinutes"`
	// `AWS::CodeBuild::Project.Triggers`.
	Triggers interface{} `json:"triggers"`
	// `AWS::CodeBuild::Project.VpcConfig`.
	VpcConfig interface{} `json:"vpcConfig"`
}

// A CloudFormation `AWS::CodeBuild::ReportGroup`.
type CfnReportGroup interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DeleteReports() interface{}
	SetDeleteReports(val interface{})
	ExportConfig() interface{}
	SetExportConfig(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	Type() *string
	SetType(val *string)
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

// The jsii proxy struct for CfnReportGroup
type jsiiProxy_CfnReportGroup struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnReportGroup) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) DeleteReports() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deleteReports",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) ExportConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"exportConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnReportGroup) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodeBuild::ReportGroup`.
func NewCfnReportGroup(scope awscdk.Construct, id *string, props *CfnReportGroupProps) CfnReportGroup {
	_init_.Initialize()

	j := jsiiProxy_CfnReportGroup{}

	_jsii_.Create(
		"monocdk.aws_codebuild.CfnReportGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodeBuild::ReportGroup`.
func NewCfnReportGroup_Override(c CfnReportGroup, scope awscdk.Construct, id *string, props *CfnReportGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.CfnReportGroup",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnReportGroup) SetDeleteReports(val interface{}) {
	_jsii_.Set(
		j,
		"deleteReports",
		val,
	)
}

func (j *jsiiProxy_CfnReportGroup) SetExportConfig(val interface{}) {
	_jsii_.Set(
		j,
		"exportConfig",
		val,
	)
}

func (j *jsiiProxy_CfnReportGroup) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnReportGroup) SetType(val *string) {
	_jsii_.Set(
		j,
		"type",
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
func CfnReportGroup_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnReportGroup",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnReportGroup_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnReportGroup",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnReportGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnReportGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnReportGroup_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.CfnReportGroup",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnReportGroup) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnReportGroup) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnReportGroup) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnReportGroup) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnReportGroup) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnReportGroup) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnReportGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnReportGroup) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnReportGroup) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnReportGroup) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnReportGroup) OnPrepare() {
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
func (c *jsiiProxy_CfnReportGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnReportGroup) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnReportGroup) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnReportGroup) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnReportGroup) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnReportGroup) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnReportGroup) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnReportGroup) ToString() *string {
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
func (c *jsiiProxy_CfnReportGroup) Validate() *[]*string {
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
func (c *jsiiProxy_CfnReportGroup) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnReportGroup_ReportExportConfigProperty struct {
	// `CfnReportGroup.ReportExportConfigProperty.ExportConfigType`.
	ExportConfigType *string `json:"exportConfigType"`
	// `CfnReportGroup.ReportExportConfigProperty.S3Destination`.
	S3Destination interface{} `json:"s3Destination"`
}

type CfnReportGroup_S3ReportExportConfigProperty struct {
	// `CfnReportGroup.S3ReportExportConfigProperty.Bucket`.
	Bucket *string `json:"bucket"`
	// `CfnReportGroup.S3ReportExportConfigProperty.BucketOwner`.
	BucketOwner *string `json:"bucketOwner"`
	// `CfnReportGroup.S3ReportExportConfigProperty.EncryptionDisabled`.
	EncryptionDisabled interface{} `json:"encryptionDisabled"`
	// `CfnReportGroup.S3ReportExportConfigProperty.EncryptionKey`.
	EncryptionKey *string `json:"encryptionKey"`
	// `CfnReportGroup.S3ReportExportConfigProperty.Packaging`.
	Packaging *string `json:"packaging"`
	// `CfnReportGroup.S3ReportExportConfigProperty.Path`.
	Path *string `json:"path"`
}

// Properties for defining a `AWS::CodeBuild::ReportGroup`.
type CfnReportGroupProps struct {
	// `AWS::CodeBuild::ReportGroup.ExportConfig`.
	ExportConfig interface{} `json:"exportConfig"`
	// `AWS::CodeBuild::ReportGroup.Type`.
	Type *string `json:"type"`
	// `AWS::CodeBuild::ReportGroup.DeleteReports`.
	DeleteReports interface{} `json:"deleteReports"`
	// `AWS::CodeBuild::ReportGroup.Name`.
	Name *string `json:"name"`
	// `AWS::CodeBuild::ReportGroup.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::CodeBuild::SourceCredential`.
type CfnSourceCredential interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AuthType() *string
	SetAuthType(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	ServerType() *string
	SetServerType(val *string)
	Stack() awscdk.Stack
	Token() *string
	SetToken(val *string)
	UpdatedProperites() *map[string]interface{}
	Username() *string
	SetUsername(val *string)
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

// The jsii proxy struct for CfnSourceCredential
type jsiiProxy_CfnSourceCredential struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnSourceCredential) AuthType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) ServerType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serverType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) Token() *string {
	var returns *string
	_jsii_.Get(
		j,
		"token",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSourceCredential) Username() *string {
	var returns *string
	_jsii_.Get(
		j,
		"username",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodeBuild::SourceCredential`.
func NewCfnSourceCredential(scope awscdk.Construct, id *string, props *CfnSourceCredentialProps) CfnSourceCredential {
	_init_.Initialize()

	j := jsiiProxy_CfnSourceCredential{}

	_jsii_.Create(
		"monocdk.aws_codebuild.CfnSourceCredential",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodeBuild::SourceCredential`.
func NewCfnSourceCredential_Override(c CfnSourceCredential, scope awscdk.Construct, id *string, props *CfnSourceCredentialProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.CfnSourceCredential",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnSourceCredential) SetAuthType(val *string) {
	_jsii_.Set(
		j,
		"authType",
		val,
	)
}

func (j *jsiiProxy_CfnSourceCredential) SetServerType(val *string) {
	_jsii_.Set(
		j,
		"serverType",
		val,
	)
}

func (j *jsiiProxy_CfnSourceCredential) SetToken(val *string) {
	_jsii_.Set(
		j,
		"token",
		val,
	)
}

func (j *jsiiProxy_CfnSourceCredential) SetUsername(val *string) {
	_jsii_.Set(
		j,
		"username",
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
func CfnSourceCredential_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnSourceCredential",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnSourceCredential_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnSourceCredential",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnSourceCredential_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.CfnSourceCredential",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnSourceCredential_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.CfnSourceCredential",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnSourceCredential) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnSourceCredential) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnSourceCredential) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnSourceCredential) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnSourceCredential) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnSourceCredential) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnSourceCredential) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnSourceCredential) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnSourceCredential) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnSourceCredential) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnSourceCredential) OnPrepare() {
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
func (c *jsiiProxy_CfnSourceCredential) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSourceCredential) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnSourceCredential) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnSourceCredential) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnSourceCredential) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnSourceCredential) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnSourceCredential) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSourceCredential) ToString() *string {
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
func (c *jsiiProxy_CfnSourceCredential) Validate() *[]*string {
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
func (c *jsiiProxy_CfnSourceCredential) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CodeBuild::SourceCredential`.
type CfnSourceCredentialProps struct {
	// `AWS::CodeBuild::SourceCredential.AuthType`.
	AuthType *string `json:"authType"`
	// `AWS::CodeBuild::SourceCredential.ServerType`.
	ServerType *string `json:"serverType"`
	// `AWS::CodeBuild::SourceCredential.Token`.
	Token *string `json:"token"`
	// `AWS::CodeBuild::SourceCredential.Username`.
	Username *string `json:"username"`
}

// Information about logs built to a CloudWatch Log Group for a build project.
// Experimental.
type CloudWatchLoggingOptions struct {
	// The current status of the logs in Amazon CloudWatch Logs for a build project.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// The Log Group to send logs to.
	// Experimental.
	LogGroup awslogs.ILogGroup `json:"logGroup"`
	// The prefix of the stream name of the Amazon CloudWatch Logs.
	// Experimental.
	Prefix *string `json:"prefix"`
}

// Construction properties for {@link CodeCommitSource}.
// Experimental.
type CodeCommitSourceProps struct {
	// The source identifier.
	//
	// This property is required on secondary sources.
	// Experimental.
	Identifier *string `json:"identifier"`
	// Experimental.
	Repository awscodecommit.IRepository `json:"repository"`
	// The commit ID, pull request ID, branch name, or tag name that corresponds to the version of the source code you want to build.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	BranchOrRef *string `json:"branchOrRef"`
	// The depth of history to download.
	//
	// Minimum value is 0.
	// If this value is 0, greater than 25, or not provided,
	// then the full history is downloaded with each build of the project.
	// Experimental.
	CloneDepth *float64 `json:"cloneDepth"`
	// Whether to fetch submodules while cloning git repo.
	// Experimental.
	FetchSubmodules *bool `json:"fetchSubmodules"`
}

// Experimental.
type CommonProjectProps struct {
	// Whether to allow the CodeBuild to send all network traffic.
	//
	// If set to false, you must individually add traffic rules to allow the
	// CodeBuild project to connect to network targets.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Indicates whether AWS CodeBuild generates a publicly accessible URL for your project's build badge.
	//
	// For more information, see Build Badges Sample
	// in the AWS CodeBuild User Guide.
	// Experimental.
	Badge *bool `json:"badge"`
	// Filename or contents of buildspec in JSON format.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec-ref-example
	//
	// Experimental.
	BuildSpec BuildSpec `json:"buildSpec"`
	// Caching strategy to use.
	// Experimental.
	Cache Cache `json:"cache"`
	// Whether to check for the presence of any secrets in the environment variables of the default type, BuildEnvironmentVariableType.PLAINTEXT. Since using a secret for the value of that kind of variable would result in it being displayed in plain text in the AWS Console, the construct will throw an exception if it detects a secret was passed there. Pass this property as false if you want to skip this validation, and keep using a secret in a plain text environment variable.
	// Experimental.
	CheckSecretsInPlainTextEnvVariables *bool `json:"checkSecretsInPlainTextEnvVariables"`
	// Maximum number of concurrent builds.
	//
	// Minimum value is 1 and maximum is account build limit.
	// Experimental.
	ConcurrentBuildLimit *float64 `json:"concurrentBuildLimit"`
	// A description of the project.
	//
	// Use the description to identify the purpose
	// of the project.
	// Experimental.
	Description *string `json:"description"`
	// Encryption key to use to read and write artifacts.
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// Build environment to use for the build.
	// Experimental.
	Environment *BuildEnvironment `json:"environment"`
	// Additional environment variables to add to the build environment.
	// Experimental.
	EnvironmentVariables *map[string]*BuildEnvironmentVariable `json:"environmentVariables"`
	// An  ProjectFileSystemLocation objects for a CodeBuild build project.
	//
	// A ProjectFileSystemLocation object specifies the identifier, location, mountOptions, mountPoint,
	// and type of a file system created using Amazon Elastic File System.
	// Experimental.
	FileSystemLocations *[]IFileSystemLocation `json:"fileSystemLocations"`
	// Add permissions to this project's role to create and use test report groups with name starting with the name of this project.
	//
	// That is the standard report group that gets created when a simple name
	// (in contrast to an ARN)
	// is used in the 'reports' section of the buildspec of this project.
	// This is usually harmless, but you can turn these off if you don't plan on using test
	// reports in this project.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/test-report-group-naming.html
	//
	// Experimental.
	GrantReportGroupPermissions *bool `json:"grantReportGroupPermissions"`
	// Information about logs for the build project.
	//
	// A project can create logs in Amazon CloudWatch Logs, an S3 bucket, or both.
	// Experimental.
	Logging *LoggingOptions `json:"logging"`
	// The physical, human-readable name of the CodeBuild Project.
	// Experimental.
	ProjectName *string `json:"projectName"`
	// The number of minutes after which AWS CodeBuild stops the build if it's still in queue.
	//
	// For valid values, see the timeoutInMinutes field in the AWS
	// CodeBuild User Guide.
	// Experimental.
	QueuedTimeout awscdk.Duration `json:"queuedTimeout"`
	// Service Role to assume while running the build.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// What security group to associate with the codebuild project's network interfaces.
	//
	// If no security group is identified, one will be created automatically.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// Where to place the network interfaces within the VPC.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SubnetSelection *awsec2.SubnetSelection `json:"subnetSelection"`
	// The number of minutes after which AWS CodeBuild stops the build if it's not complete.
	//
	// For valid values, see the timeoutInMinutes field in the AWS
	// CodeBuild User Guide.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// VPC network to place codebuild network interfaces.
	//
	// Specify this if the codebuild project needs to access resources in a VPC.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
}

// Build machine compute type.
// Experimental.
type ComputeType string

const (
	ComputeType_SMALL ComputeType = "SMALL"
	ComputeType_MEDIUM ComputeType = "MEDIUM"
	ComputeType_LARGE ComputeType = "LARGE"
	ComputeType_X2_LARGE ComputeType = "X2_LARGE"
)

// The options when creating a CodeBuild Docker build image using {@link LinuxBuildImage.fromDockerRegistry} or {@link WindowsBuildImage.fromDockerRegistry}.
// Experimental.
type DockerImageOptions struct {
	// The credentials, stored in Secrets Manager, used for accessing the repository holding the image, if the repository is private.
	// Experimental.
	SecretsManagerCredentials awssecretsmanager.ISecret `json:"secretsManagerCredentials"`
}

// Construction properties for {@link EfsFileSystemLocation}.
// Experimental.
type EfsFileSystemLocationProps struct {
	// The name used to access a file system created by Amazon EFS.
	// Experimental.
	Identifier *string `json:"identifier"`
	// A string that specifies the location of the file system, like Amazon EFS.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Location *string `json:"location"`
	// The location in the container where you mount the file system.
	// Experimental.
	MountPoint *string `json:"mountPoint"`
	// The mount options for a file system such as Amazon EFS.
	// Experimental.
	MountOptions *string `json:"mountOptions"`
}

// The types of webhook event actions.
// Experimental.
type EventAction string

const (
	EventAction_PUSH EventAction = "PUSH"
	EventAction_PULL_REQUEST_CREATED EventAction = "PULL_REQUEST_CREATED"
	EventAction_PULL_REQUEST_UPDATED EventAction = "PULL_REQUEST_UPDATED"
	EventAction_PULL_REQUEST_MERGED EventAction = "PULL_REQUEST_MERGED"
	EventAction_PULL_REQUEST_REOPENED EventAction = "PULL_REQUEST_REOPENED"
)

// The type returned from {@link IFileSystemLocation#bind}.
// Experimental.
type FileSystemConfig struct {
	// File system location wrapper property.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projectfilesystemlocation.html
	//
	// Experimental.
	Location *CfnProject_ProjectFileSystemLocationProperty `json:"location"`
}

// FileSystemLocation provider definition for a CodeBuild Project.
// Experimental.
type FileSystemLocation interface {
}

// The jsii proxy struct for FileSystemLocation
type jsiiProxy_FileSystemLocation struct {
	_ byte // padding
}

// Experimental.
func NewFileSystemLocation() FileSystemLocation {
	_init_.Initialize()

	j := jsiiProxy_FileSystemLocation{}

	_jsii_.Create(
		"monocdk.aws_codebuild.FileSystemLocation",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewFileSystemLocation_Override(f FileSystemLocation) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.FileSystemLocation",
		nil, // no parameters
		f,
	)
}

// EFS file system provider.
// Experimental.
func FileSystemLocation_Efs(props *EfsFileSystemLocationProps) IFileSystemLocation {
	_init_.Initialize()

	var returns IFileSystemLocation

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.FileSystemLocation",
		"efs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// An object that represents a group of filter conditions for a webhook.
//
// Every condition in a given FilterGroup must be true in order for the whole group to be true.
// You construct instances of it by calling the {@link #inEventOf} static factory method,
// and then calling various `andXyz` instance methods to create modified instances of it
// (this class is immutable).
//
// You pass instances of this class to the `webhookFilters` property when constructing a source.
// Experimental.
type FilterGroup interface {
	AndActorAccountIs(pattern *string) FilterGroup
	AndActorAccountIsNot(pattern *string) FilterGroup
	AndBaseBranchIs(branchName *string) FilterGroup
	AndBaseBranchIsNot(branchName *string) FilterGroup
	AndBaseRefIs(pattern *string) FilterGroup
	AndBaseRefIsNot(pattern *string) FilterGroup
	AndBranchIs(branchName *string) FilterGroup
	AndBranchIsNot(branchName *string) FilterGroup
	AndCommitMessageIs(commitMessage *string) FilterGroup
	AndCommitMessageIsNot(commitMessage *string) FilterGroup
	AndFilePathIs(pattern *string) FilterGroup
	AndFilePathIsNot(pattern *string) FilterGroup
	AndHeadRefIs(pattern *string) FilterGroup
	AndHeadRefIsNot(pattern *string) FilterGroup
	AndTagIs(tagName *string) FilterGroup
	AndTagIsNot(tagName *string) FilterGroup
}

// The jsii proxy struct for FilterGroup
type jsiiProxy_FilterGroup struct {
	_ byte // padding
}

// Creates a new event FilterGroup that triggers on any of the provided actions.
// Experimental.
func FilterGroup_InEventOf(actions ...EventAction) FilterGroup {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns FilterGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.FilterGroup",
		"inEventOf",
		args,
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the account ID of the actor initiating the event must match the given pattern.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndActorAccountIs(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andActorAccountIs",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the account ID of the actor initiating the event must not match the given pattern.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndActorAccountIsNot(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andActorAccountIsNot",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the Pull Request that is the source of the event must target the given base branch.
//
// Note that you cannot use this method if this Group contains the `PUSH` event action.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndBaseBranchIs(branchName *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andBaseBranchIs",
		[]interface{}{branchName},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the Pull Request that is the source of the event must not target the given base branch.
//
// Note that you cannot use this method if this Group contains the `PUSH` event action.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndBaseBranchIsNot(branchName *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andBaseBranchIsNot",
		[]interface{}{branchName},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the Pull Request that is the source of the event must target the given Git reference.
//
// Note that you cannot use this method if this Group contains the `PUSH` event action.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndBaseRefIs(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andBaseRefIs",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the Pull Request that is the source of the event must not target the given Git reference.
//
// Note that you cannot use this method if this Group contains the `PUSH` event action.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndBaseRefIsNot(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andBaseRefIsNot",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must affect the given branch.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndBranchIs(branchName *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andBranchIs",
		[]interface{}{branchName},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must not affect the given branch.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndBranchIsNot(branchName *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andBranchIsNot",
		[]interface{}{branchName},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must affect a head commit with the given message.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndCommitMessageIs(commitMessage *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andCommitMessageIs",
		[]interface{}{commitMessage},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must not affect a head commit with the given message.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndCommitMessageIsNot(commitMessage *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andCommitMessageIsNot",
		[]interface{}{commitMessage},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the push that is the source of the event must affect a file that matches the given pattern.
//
// Note that you can only use this method if this Group contains only the `PUSH` event action,
// and only for GitHub, Bitbucket and GitHubEnterprise sources.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndFilePathIs(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andFilePathIs",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the push that is the source of the event must not affect a file that matches the given pattern.
//
// Note that you can only use this method if this Group contains only the `PUSH` event action,
// and only for GitHub, Bitbucket and GitHubEnterprise sources.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndFilePathIsNot(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andFilePathIsNot",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must affect a Git reference (ie., a branch or a tag) that matches the given pattern.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndHeadRefIs(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andHeadRefIs",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must not affect a Git reference (ie., a branch or a tag) that matches the given pattern.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndHeadRefIsNot(pattern *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andHeadRefIsNot",
		[]interface{}{pattern},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must affect the given tag.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndTagIs(tagName *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andTagIs",
		[]interface{}{tagName},
		&returns,
	)

	return returns
}

// Create a new FilterGroup with an added condition: the event must not affect the given tag.
// Experimental.
func (f *jsiiProxy_FilterGroup) AndTagIsNot(tagName *string) FilterGroup {
	var returns FilterGroup

	_jsii_.Invoke(
		f,
		"andTagIsNot",
		[]interface{}{tagName},
		&returns,
	)

	return returns
}

// The source credentials used when contacting the GitHub Enterprise API.
//
// **Note**: CodeBuild only allows a single credential for GitHub Enterprise
// to be saved in a given AWS account in a given region -
// any attempt to add more than one will result in an error.
// Experimental.
type GitHubEnterpriseSourceCredentials interface {
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

// The jsii proxy struct for GitHubEnterpriseSourceCredentials
type jsiiProxy_GitHubEnterpriseSourceCredentials struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_GitHubEnterpriseSourceCredentials) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubEnterpriseSourceCredentials) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubEnterpriseSourceCredentials) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubEnterpriseSourceCredentials) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewGitHubEnterpriseSourceCredentials(scope constructs.Construct, id *string, props *GitHubEnterpriseSourceCredentialsProps) GitHubEnterpriseSourceCredentials {
	_init_.Initialize()

	j := jsiiProxy_GitHubEnterpriseSourceCredentials{}

	_jsii_.Create(
		"monocdk.aws_codebuild.GitHubEnterpriseSourceCredentials",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGitHubEnterpriseSourceCredentials_Override(g GitHubEnterpriseSourceCredentials, scope constructs.Construct, id *string, props *GitHubEnterpriseSourceCredentialsProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.GitHubEnterpriseSourceCredentials",
		[]interface{}{scope, id, props},
		g,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func GitHubEnterpriseSourceCredentials_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.GitHubEnterpriseSourceCredentials",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func GitHubEnterpriseSourceCredentials_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.GitHubEnterpriseSourceCredentials",
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
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		g,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) OnPrepare() {
	_jsii_.InvokeVoid(
		g,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		g,
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
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) Prepare() {
	_jsii_.InvokeVoid(
		g,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		g,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubEnterpriseSourceCredentials) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		g,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Creation properties for {@link GitHubEnterpriseSourceCredentials}.
// Experimental.
type GitHubEnterpriseSourceCredentialsProps struct {
	// The personal access token to use when contacting the instance of the GitHub Enterprise API.
	// Experimental.
	AccessToken awscdk.SecretValue `json:"accessToken"`
}

// Construction properties for {@link GitHubEnterpriseSource}.
// Experimental.
type GitHubEnterpriseSourceProps struct {
	// The source identifier.
	//
	// This property is required on secondary sources.
	// Experimental.
	Identifier *string `json:"identifier"`
	// The HTTPS URL of the repository in your GitHub Enterprise installation.
	// Experimental.
	HttpsCloneUrl *string `json:"httpsCloneUrl"`
	// The commit ID, pull request ID, branch name, or tag name that corresponds to the version of the source code you want to build.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	BranchOrRef *string `json:"branchOrRef"`
	// The depth of history to download.
	//
	// Minimum value is 0.
	// If this value is 0, greater than 25, or not provided,
	// then the full history is downloaded with each build of the project.
	// Experimental.
	CloneDepth *float64 `json:"cloneDepth"`
	// Whether to fetch submodules while cloning git repo.
	// Experimental.
	FetchSubmodules *bool `json:"fetchSubmodules"`
	// Whether to ignore SSL errors when connecting to the repository.
	// Experimental.
	IgnoreSslErrors *bool `json:"ignoreSslErrors"`
	// Whether to send notifications on your build's start and end.
	// Experimental.
	ReportBuildStatus *bool `json:"reportBuildStatus"`
	// Whether to create a webhook that will trigger a build every time an event happens in the repository.
	// Experimental.
	Webhook *bool `json:"webhook"`
	// A list of webhook filters that can constraint what events in the repository will trigger a build.
	//
	// A build is triggered if any of the provided filter groups match.
	// Only valid if `webhook` was not provided as false.
	// Experimental.
	WebhookFilters *[]FilterGroup `json:"webhookFilters"`
	// Trigger a batch build from a webhook instead of a standard one.
	//
	// Enabling this will enable batch builds on the CodeBuild project.
	// Experimental.
	WebhookTriggersBatchBuild *bool `json:"webhookTriggersBatchBuild"`
}

// The source credentials used when contacting the GitHub API.
//
// **Note**: CodeBuild only allows a single credential for GitHub
// to be saved in a given AWS account in a given region -
// any attempt to add more than one will result in an error.
// Experimental.
type GitHubSourceCredentials interface {
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

// The jsii proxy struct for GitHubSourceCredentials
type jsiiProxy_GitHubSourceCredentials struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_GitHubSourceCredentials) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubSourceCredentials) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubSourceCredentials) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GitHubSourceCredentials) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewGitHubSourceCredentials(scope constructs.Construct, id *string, props *GitHubSourceCredentialsProps) GitHubSourceCredentials {
	_init_.Initialize()

	j := jsiiProxy_GitHubSourceCredentials{}

	_jsii_.Create(
		"monocdk.aws_codebuild.GitHubSourceCredentials",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGitHubSourceCredentials_Override(g GitHubSourceCredentials, scope constructs.Construct, id *string, props *GitHubSourceCredentialsProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.GitHubSourceCredentials",
		[]interface{}{scope, id, props},
		g,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func GitHubSourceCredentials_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.GitHubSourceCredentials",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func GitHubSourceCredentials_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.GitHubSourceCredentials",
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
func (g *jsiiProxy_GitHubSourceCredentials) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		g,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (g *jsiiProxy_GitHubSourceCredentials) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubSourceCredentials) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubSourceCredentials) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubSourceCredentials) OnPrepare() {
	_jsii_.InvokeVoid(
		g,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (g *jsiiProxy_GitHubSourceCredentials) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		g,
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
func (g *jsiiProxy_GitHubSourceCredentials) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubSourceCredentials) Prepare() {
	_jsii_.InvokeVoid(
		g,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (g *jsiiProxy_GitHubSourceCredentials) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		g,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (g *jsiiProxy_GitHubSourceCredentials) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		g,
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
func (g *jsiiProxy_GitHubSourceCredentials) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		g,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Creation properties for {@link GitHubSourceCredentials}.
// Experimental.
type GitHubSourceCredentialsProps struct {
	// The personal access token to use when contacting the GitHub API.
	// Experimental.
	AccessToken awscdk.SecretValue `json:"accessToken"`
}

// Construction properties for {@link GitHubSource} and {@link GitHubEnterpriseSource}.
// Experimental.
type GitHubSourceProps struct {
	// The source identifier.
	//
	// This property is required on secondary sources.
	// Experimental.
	Identifier *string `json:"identifier"`
	// The GitHub account/user that owns the repo.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Owner *string `json:"owner"`
	// The name of the repo (without the username).
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Repo *string `json:"repo"`
	// The commit ID, pull request ID, branch name, or tag name that corresponds to the version of the source code you want to build.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	BranchOrRef *string `json:"branchOrRef"`
	// The depth of history to download.
	//
	// Minimum value is 0.
	// If this value is 0, greater than 25, or not provided,
	// then the full history is downloaded with each build of the project.
	// Experimental.
	CloneDepth *float64 `json:"cloneDepth"`
	// Whether to fetch submodules while cloning git repo.
	// Experimental.
	FetchSubmodules *bool `json:"fetchSubmodules"`
	// Whether to send notifications on your build's start and end.
	// Experimental.
	ReportBuildStatus *bool `json:"reportBuildStatus"`
	// Whether to create a webhook that will trigger a build every time an event happens in the repository.
	// Experimental.
	Webhook *bool `json:"webhook"`
	// A list of webhook filters that can constraint what events in the repository will trigger a build.
	//
	// A build is triggered if any of the provided filter groups match.
	// Only valid if `webhook` was not provided as false.
	// Experimental.
	WebhookFilters *[]FilterGroup `json:"webhookFilters"`
	// Trigger a batch build from a webhook instead of a standard one.
	//
	// Enabling this will enable batch builds on the CodeBuild project.
	// Experimental.
	WebhookTriggersBatchBuild *bool `json:"webhookTriggersBatchBuild"`
}

// The abstract interface of a CodeBuild build output.
//
// Implemented by {@link Artifacts}.
// Experimental.
type IArtifacts interface {
	// Callback when an Artifacts class is used in a CodeBuild Project.
	// Experimental.
	Bind(scope awscdk.Construct, project IProject) *ArtifactsConfig
	// The artifact identifier.
	//
	// This property is required on secondary artifacts.
	// Experimental.
	Identifier() *string
	// The CodeBuild type of this artifact.
	// Experimental.
	Type() *string
}

// The jsii proxy for IArtifacts
type jsiiProxy_IArtifacts struct {
	_ byte // padding
}

func (i *jsiiProxy_IArtifacts) Bind(scope awscdk.Construct, project IProject) *ArtifactsConfig {
	var returns *ArtifactsConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, project},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IArtifacts) Identifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IArtifacts) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

// A variant of {@link IBuildImage} that allows binding to the project.
// Experimental.
type IBindableBuildImage interface {
	IBuildImage
	// Function that allows the build image access to the construct tree.
	// Experimental.
	Bind(scope awscdk.Construct, project IProject, options *BuildImageBindOptions) *BuildImageConfig
}

// The jsii proxy for IBindableBuildImage
type jsiiProxy_IBindableBuildImage struct {
	jsiiProxy_IBuildImage
}

func (i *jsiiProxy_IBindableBuildImage) Bind(scope awscdk.Construct, project IProject, options *BuildImageBindOptions) *BuildImageConfig {
	var returns *BuildImageConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, project, options},
		&returns,
	)

	return returns
}

// Represents a Docker image used for the CodeBuild Project builds.
//
// Use the concrete subclasses, either:
// {@link LinuxBuildImage} or {@link WindowsBuildImage}.
// Experimental.
type IBuildImage interface {
	// Make a buildspec to run the indicated script.
	// Experimental.
	RunScriptBuildspec(entrypoint *string) BuildSpec
	// Allows the image a chance to validate whether the passed configuration is correct.
	// Experimental.
	Validate(buildEnvironment *BuildEnvironment) *[]*string
	// The default {@link ComputeType} to use with this image, if one was not specified in {@link BuildEnvironment#computeType} explicitly.
	// Experimental.
	DefaultComputeType() ComputeType
	// The Docker image identifier that the build environment uses.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-available.html
	//
	// Experimental.
	ImageId() *string
	// The type of principal that CodeBuild will use to pull this build Docker image.
	// Experimental.
	ImagePullPrincipalType() ImagePullPrincipalType
	// An optional ECR repository that the image is hosted in.
	// Experimental.
	Repository() awsecr.IRepository
	// The secretsManagerCredentials for access to a private registry.
	// Experimental.
	SecretsManagerCredentials() awssecretsmanager.ISecret
	// The type of build environment.
	// Experimental.
	Type() *string
}

// The jsii proxy for IBuildImage
type jsiiProxy_IBuildImage struct {
	_ byte // padding
}

func (i *jsiiProxy_IBuildImage) RunScriptBuildspec(entrypoint *string) BuildSpec {
	var returns BuildSpec

	_jsii_.Invoke(
		i,
		"runScriptBuildspec",
		[]interface{}{entrypoint},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBuildImage) Validate(buildEnvironment *BuildEnvironment) *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		i,
		"validate",
		[]interface{}{buildEnvironment},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IBuildImage) DefaultComputeType() ComputeType {
	var returns ComputeType
	_jsii_.Get(
		j,
		"defaultComputeType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBuildImage) ImageId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBuildImage) ImagePullPrincipalType() ImagePullPrincipalType {
	var returns ImagePullPrincipalType
	_jsii_.Get(
		j,
		"imagePullPrincipalType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBuildImage) Repository() awsecr.IRepository {
	var returns awsecr.IRepository
	_jsii_.Get(
		j,
		"repository",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBuildImage) SecretsManagerCredentials() awssecretsmanager.ISecret {
	var returns awssecretsmanager.ISecret
	_jsii_.Get(
		j,
		"secretsManagerCredentials",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBuildImage) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

// The interface of a CodeBuild FileSystemLocation.
//
// Implemented by {@link EfsFileSystemLocation}.
// Experimental.
type IFileSystemLocation interface {
	// Called by the project when a file system is added so it can perform binding operations on this file system location.
	// Experimental.
	Bind(scope awscdk.Construct, project IProject) *FileSystemConfig
}

// The jsii proxy for IFileSystemLocation
type jsiiProxy_IFileSystemLocation struct {
	_ byte // padding
}

func (i *jsiiProxy_IFileSystemLocation) Bind(scope awscdk.Construct, project IProject) *FileSystemConfig {
	var returns *FileSystemConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, project},
		&returns,
	)

	return returns
}

// Experimental.
type IProject interface {
	awsec2.IConnectable
	awsiam.IGrantable
	awscodestarnotifications.INotificationRuleSource
	awscdk.IResource
	// Experimental.
	AddToRolePolicy(policyStatement awsiam.PolicyStatement)
	// Enable batch builds.
	//
	// Returns an object contining the batch service role if batch builds
	// could be enabled.
	// Experimental.
	EnableBatchBuilds() *BatchBuildConfig
	// Returns: a CloudWatch metric associated with this build project.
	// Experimental.
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Measures the number of builds triggered.
	//
	// Units: Count
	//
	// Valid CloudWatch statistics: Sum
	// Experimental.
	MetricBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Measures the duration of all builds over time.
	//
	// Units: Seconds
	//
	// Valid CloudWatch statistics: Average (recommended), Maximum, Minimum
	// Experimental.
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Measures the number of builds that failed because of client error or because of a timeout.
	//
	// Units: Count
	//
	// Valid CloudWatch statistics: Sum
	// Experimental.
	MetricFailedBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Measures the number of successful builds.
	//
	// Units: Count
	//
	// Valid CloudWatch statistics: Sum
	// Experimental.
	MetricSucceededBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	// Defines a CodeStar Notification rule triggered when the project events emitted by you specified, it very similar to `onEvent` API.
	//
	// You can also use the methods `notifyOnBuildSucceeded` and
	// `notifyOnBuildFailed` to define rules for these specific event emitted.
	//
	// Returns: CodeStar Notifications rule associated with this build project.
	// Experimental.
	NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *ProjectNotifyOnOptions) awscodestarnotifications.INotificationRule
	// Defines a CodeStar notification rule which triggers when a build fails.
	// Experimental.
	NotifyOnBuildFailed(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	// Defines a CodeStar notification rule which triggers when a build completes successfully.
	// Experimental.
	NotifyOnBuildSucceeded(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	// Defines an event rule which triggers when a build fails.
	// Experimental.
	OnBuildFailed(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines an event rule which triggers when a build starts.
	// Experimental.
	OnBuildStarted(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines an event rule which triggers when a build completes successfully.
	// Experimental.
	OnBuildSucceeded(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule triggered when something happens with this project.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
	//
	// Experimental.
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule that triggers upon phase change of this build project.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
	//
	// Experimental.
	OnPhaseChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule triggered when the build project state changes.
	//
	// You can filter specific build status events using an event
	// pattern filter on the `build-status` detail field:
	//
	//     const rule = project.onStateChange('OnBuildStarted', { target });
	//     rule.addEventPattern({
	//       detail: {
	//         'build-status': [
	//           "IN_PROGRESS",
	//           "SUCCEEDED",
	//           "FAILED",
	//           "STOPPED"
	//         ]
	//       }
	//     });
	//
	// You can also use the methods `onBuildFailed` and `onBuildSucceeded` to define rules for
	// these specific state changes.
	//
	// To access fields from the event in the event target input,
	// use the static fields on the `StateChangeEvent` class.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
	//
	// Experimental.
	OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// The ARN of this Project.
	// Experimental.
	ProjectArn() *string
	// The human-visible name of this Project.
	// Experimental.
	ProjectName() *string
	// The IAM service Role of this Project.
	//
	// Undefined for imported Projects.
	// Experimental.
	Role() awsiam.IRole
}

// The jsii proxy for IProject
type jsiiProxy_IProject struct {
	internal.Type__awsec2IConnectable
	internal.Type__awsiamIGrantable
	internal.Type__awscodestarnotificationsINotificationRuleSource
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IProject) AddToRolePolicy(policyStatement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		i,
		"addToRolePolicy",
		[]interface{}{policyStatement},
	)
}

func (i *jsiiProxy_IProject) EnableBatchBuilds() *BatchBuildConfig {
	var returns *BatchBuildConfig

	_jsii_.Invoke(
		i,
		"enableBatchBuilds",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) MetricBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) MetricFailedBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricFailedBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) MetricSucceededBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		i,
		"metricSucceededBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *ProjectNotifyOnOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOn",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) NotifyOnBuildFailed(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOnBuildFailed",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) NotifyOnBuildSucceeded(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOnBuildSucceeded",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) OnBuildFailed(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onBuildFailed",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) OnBuildStarted(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onBuildStarted",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) OnBuildSucceeded(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onBuildSucceeded",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) OnPhaseChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onPhaseChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onStateChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IProject) BindAsNotificationRuleSource(scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig {
	var returns *awscodestarnotifications.NotificationRuleSourceConfig

	_jsii_.Invoke(
		i,
		"bindAsNotificationRuleSource",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IProject) ProjectArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"projectArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IProject) ProjectName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"projectName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IProject) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IProject) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IProject) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IProject) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IProject) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IProject) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// The interface representing the ReportGroup resource - either an existing one, imported using the {@link ReportGroup.fromReportGroupName} method, or a new one, created with the {@link ReportGroup} class.
// Experimental.
type IReportGroup interface {
	awscdk.IResource
	// Grants the given entity permissions to write (that is, upload reports to) this report group.
	// Experimental.
	GrantWrite(identity awsiam.IGrantable) awsiam.Grant
	// The ARN of the ReportGroup.
	// Experimental.
	ReportGroupArn() *string
	// The name of the ReportGroup.
	// Experimental.
	ReportGroupName() *string
}

// The jsii proxy for IReportGroup
type jsiiProxy_IReportGroup struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IReportGroup) GrantWrite(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantWrite",
		[]interface{}{identity},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IReportGroup) ReportGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"reportGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IReportGroup) ReportGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"reportGroupName",
		&returns,
	)
	return returns
}

// The abstract interface of a CodeBuild source.
//
// Implemented by {@link Source}.
// Experimental.
type ISource interface {
	// Experimental.
	Bind(scope awscdk.Construct, project IProject) *SourceConfig
	// Experimental.
	BadgeSupported() *bool
	// Experimental.
	Identifier() *string
	// Experimental.
	Type() *string
}

// The jsii proxy for ISource
type jsiiProxy_ISource struct {
	_ byte // padding
}

func (i *jsiiProxy_ISource) Bind(scope awscdk.Construct, project IProject) *SourceConfig {
	var returns *SourceConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, project},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_ISource) BadgeSupported() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"badgeSupported",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISource) Identifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISource) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

// The type of principal CodeBuild will use to pull your build Docker image.
// Experimental.
type ImagePullPrincipalType string

const (
	ImagePullPrincipalType_CODEBUILD ImagePullPrincipalType = "CODEBUILD"
	ImagePullPrincipalType_SERVICE_ROLE ImagePullPrincipalType = "SERVICE_ROLE"
)

// A CodeBuild image running Linux.
//
// This class has a bunch of public constants that represent the most popular images.
//
// You can also specify a custom image using one of the static methods:
//
// - LinuxBuildImage.fromDockerRegistry(image[, { secretsManagerCredentials }])
// - LinuxBuildImage.fromEcrRepository(repo[, tag])
// - LinuxBuildImage.fromAsset(parent, id, props)
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-available.html
//
// Experimental.
type LinuxBuildImage interface {
	IBuildImage
	DefaultComputeType() ComputeType
	ImageId() *string
	ImagePullPrincipalType() ImagePullPrincipalType
	Repository() awsecr.IRepository
	SecretsManagerCredentials() awssecretsmanager.ISecret
	Type() *string
	RunScriptBuildspec(entrypoint *string) BuildSpec
	Validate(_arg *BuildEnvironment) *[]*string
}

// The jsii proxy struct for LinuxBuildImage
type jsiiProxy_LinuxBuildImage struct {
	jsiiProxy_IBuildImage
}

func (j *jsiiProxy_LinuxBuildImage) DefaultComputeType() ComputeType {
	var returns ComputeType
	_jsii_.Get(
		j,
		"defaultComputeType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxBuildImage) ImageId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxBuildImage) ImagePullPrincipalType() ImagePullPrincipalType {
	var returns ImagePullPrincipalType
	_jsii_.Get(
		j,
		"imagePullPrincipalType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxBuildImage) Repository() awsecr.IRepository {
	var returns awsecr.IRepository
	_jsii_.Get(
		j,
		"repository",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxBuildImage) SecretsManagerCredentials() awssecretsmanager.ISecret {
	var returns awssecretsmanager.ISecret
	_jsii_.Get(
		j,
		"secretsManagerCredentials",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxBuildImage) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Uses an Docker image asset as a Linux build image.
// Experimental.
func LinuxBuildImage_FromAsset(scope constructs.Construct, id *string, props *awsecrassets.DockerImageAssetProps) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"fromAsset",
		[]interface{}{scope, id, props},
		&returns,
	)

	return returns
}

// Uses a Docker image provided by CodeBuild.
//
// Returns: A Docker image provided by CodeBuild.
//
// TODO: EXAMPLE
//
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-available.html
//
// Experimental.
func LinuxBuildImage_FromCodeBuildImageId(id *string) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"fromCodeBuildImageId",
		[]interface{}{id},
		&returns,
	)

	return returns
}

// Returns: a Linux build image from a Docker Hub image.
// Experimental.
func LinuxBuildImage_FromDockerRegistry(name *string, options *DockerImageOptions) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"fromDockerRegistry",
		[]interface{}{name, options},
		&returns,
	)

	return returns
}

// Returns: A Linux build image from an ECR repository.
//
// NOTE: if the repository is external (i.e. imported), then we won't be able to add
// a resource policy statement for it so CodeBuild can pull the image.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-ecr.html
//
// Experimental.
func LinuxBuildImage_FromEcrRepository(repository awsecr.IRepository, tag *string) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"fromEcrRepository",
		[]interface{}{repository, tag},
		&returns,
	)

	return returns
}

func LinuxBuildImage_AMAZON_LINUX_2() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"AMAZON_LINUX_2",
		&returns,
	)
	return returns
}

func LinuxBuildImage_AMAZON_LINUX_2_2() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"AMAZON_LINUX_2_2",
		&returns,
	)
	return returns
}

func LinuxBuildImage_AMAZON_LINUX_2_3() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"AMAZON_LINUX_2_3",
		&returns,
	)
	return returns
}

func LinuxBuildImage_AMAZON_LINUX_2_ARM() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"AMAZON_LINUX_2_ARM",
		&returns,
	)
	return returns
}

func LinuxBuildImage_STANDARD_1_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"STANDARD_1_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_STANDARD_2_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"STANDARD_2_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_STANDARD_3_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"STANDARD_3_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_STANDARD_4_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"STANDARD_4_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_STANDARD_5_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"STANDARD_5_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_ANDROID_JAVA8_24_4_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_ANDROID_JAVA8_24_4_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_ANDROID_JAVA8_26_1_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_ANDROID_JAVA8_26_1_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_BASE() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_BASE",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_DOCKER_17_09_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_DOCKER_17_09_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_DOCKER_18_09_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_DOCKER_18_09_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_DOTNET_CORE_1_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_DOTNET_CORE_1_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_DOTNET_CORE_2_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_DOTNET_CORE_2_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_DOTNET_CORE_2_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_DOTNET_CORE_2_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_GOLANG_1_10() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_GOLANG_1_10",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_GOLANG_1_11() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_GOLANG_1_11",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_NODEJS_10_1_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_NODEJS_10_1_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_NODEJS_10_14_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_NODEJS_10_14_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_NODEJS_6_3_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_NODEJS_6_3_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_NODEJS_8_11_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_NODEJS_8_11_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_OPEN_JDK_11() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_OPEN_JDK_11",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_OPEN_JDK_8() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_OPEN_JDK_8",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_OPEN_JDK_9() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_OPEN_JDK_9",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PHP_5_6() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PHP_5_6",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PHP_7_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PHP_7_0",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PHP_7_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PHP_7_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PYTHON_2_7_12() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PYTHON_2_7_12",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PYTHON_3_3_6() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PYTHON_3_3_6",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PYTHON_3_4_5() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PYTHON_3_4_5",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PYTHON_3_5_2() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PYTHON_3_5_2",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PYTHON_3_6_5() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PYTHON_3_6_5",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_PYTHON_3_7_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_PYTHON_3_7_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_RUBY_2_2_5() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_RUBY_2_2_5",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_RUBY_2_3_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_RUBY_2_3_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_RUBY_2_5_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_RUBY_2_5_1",
		&returns,
	)
	return returns
}

func LinuxBuildImage_UBUNTU_14_04_RUBY_2_5_3() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxBuildImage",
		"UBUNTU_14_04_RUBY_2_5_3",
		&returns,
	)
	return returns
}

// Make a buildspec to run the indicated script.
// Experimental.
func (l *jsiiProxy_LinuxBuildImage) RunScriptBuildspec(entrypoint *string) BuildSpec {
	var returns BuildSpec

	_jsii_.Invoke(
		l,
		"runScriptBuildspec",
		[]interface{}{entrypoint},
		&returns,
	)

	return returns
}

// Allows the image a chance to validate whether the passed configuration is correct.
// Experimental.
func (l *jsiiProxy_LinuxBuildImage) Validate(_arg *BuildEnvironment) *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		[]interface{}{_arg},
		&returns,
	)

	return returns
}

// A CodeBuild GPU image running Linux.
//
// This class has public constants that represent the most popular GPU images from AWS Deep Learning Containers.
// See: https://aws.amazon.com/releasenotes/available-deep-learning-containers-images
//
// Experimental.
type LinuxGpuBuildImage interface {
	IBindableBuildImage
	DefaultComputeType() ComputeType
	ImageId() *string
	ImagePullPrincipalType() ImagePullPrincipalType
	Type() *string
	Bind(scope awscdk.Construct, project IProject, _options *BuildImageBindOptions) *BuildImageConfig
	RunScriptBuildspec(entrypoint *string) BuildSpec
	Validate(buildEnvironment *BuildEnvironment) *[]*string
}

// The jsii proxy struct for LinuxGpuBuildImage
type jsiiProxy_LinuxGpuBuildImage struct {
	jsiiProxy_IBindableBuildImage
}

func (j *jsiiProxy_LinuxGpuBuildImage) DefaultComputeType() ComputeType {
	var returns ComputeType
	_jsii_.Get(
		j,
		"defaultComputeType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxGpuBuildImage) ImageId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxGpuBuildImage) ImagePullPrincipalType() ImagePullPrincipalType {
	var returns ImagePullPrincipalType
	_jsii_.Get(
		j,
		"imagePullPrincipalType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LinuxGpuBuildImage) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Returns a Linux GPU build image from AWS Deep Learning Containers.
// See: https://aws.amazon.com/releasenotes/available-deep-learning-containers-images
//
// Experimental.
func LinuxGpuBuildImage_AwsDeepLearningContainersImage(repositoryName *string, tag *string, account *string) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"awsDeepLearningContainersImage",
		[]interface{}{repositoryName, tag, account},
		&returns,
	)

	return returns
}

func LinuxGpuBuildImage_DLC_MXNET_1_4_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_MXNET_1_4_1",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_MXNET_1_6_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_MXNET_1_6_0",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_PYTORCH_1_2_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_PYTORCH_1_2_0",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_PYTORCH_1_3_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_PYTORCH_1_3_1",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_PYTORCH_1_4_0_INFERENCE() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_PYTORCH_1_4_0_INFERENCE",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_PYTORCH_1_4_0_TRAINING() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_PYTORCH_1_4_0_TRAINING",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_PYTORCH_1_5_0_INFERENCE() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_PYTORCH_1_5_0_INFERENCE",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_PYTORCH_1_5_0_TRAINING() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_PYTORCH_1_5_0_TRAINING",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_1_14_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_1_14_0",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_1_15_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_1_15_0",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_1_15_2_INFERENCE() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_1_15_2_INFERENCE",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_1_15_2_TRAINING() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_1_15_2_TRAINING",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_2_0_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_2_0_0",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_2_0_1() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_2_0_1",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_2_1_0_INFERENCE() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_2_1_0_INFERENCE",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_2_1_0_TRAINING() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_2_1_0_TRAINING",
		&returns,
	)
	return returns
}

func LinuxGpuBuildImage_DLC_TENSORFLOW_2_2_0_TRAINING() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.LinuxGpuBuildImage",
		"DLC_TENSORFLOW_2_2_0_TRAINING",
		&returns,
	)
	return returns
}

// Function that allows the build image access to the construct tree.
// Experimental.
func (l *jsiiProxy_LinuxGpuBuildImage) Bind(scope awscdk.Construct, project IProject, _options *BuildImageBindOptions) *BuildImageConfig {
	var returns *BuildImageConfig

	_jsii_.Invoke(
		l,
		"bind",
		[]interface{}{scope, project, _options},
		&returns,
	)

	return returns
}

// Make a buildspec to run the indicated script.
// Experimental.
func (l *jsiiProxy_LinuxGpuBuildImage) RunScriptBuildspec(entrypoint *string) BuildSpec {
	var returns BuildSpec

	_jsii_.Invoke(
		l,
		"runScriptBuildspec",
		[]interface{}{entrypoint},
		&returns,
	)

	return returns
}

// Allows the image a chance to validate whether the passed configuration is correct.
// Experimental.
func (l *jsiiProxy_LinuxGpuBuildImage) Validate(buildEnvironment *BuildEnvironment) *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		[]interface{}{buildEnvironment},
		&returns,
	)

	return returns
}

// Local cache modes to enable for the CodeBuild Project.
// Experimental.
type LocalCacheMode string

const (
	LocalCacheMode_SOURCE LocalCacheMode = "SOURCE"
	LocalCacheMode_DOCKER_LAYER LocalCacheMode = "DOCKER_LAYER"
	LocalCacheMode_CUSTOM LocalCacheMode = "CUSTOM"
)

// Information about logs for the build project.
//
// A project can create logs in Amazon CloudWatch Logs, an S3 bucket, or both.
// Experimental.
type LoggingOptions struct {
	// Information about Amazon CloudWatch Logs for a build project.
	// Experimental.
	CloudWatch *CloudWatchLoggingOptions `json:"cloudWatch"`
	// Information about logs built to an S3 bucket for a build project.
	// Experimental.
	S3 *S3LoggingOptions `json:"s3"`
}

// Event fields for the CodeBuild "phase change" event.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html#sample-build-notifications-ref
//
// Experimental.
type PhaseChangeEvent interface {
}

// The jsii proxy struct for PhaseChangeEvent
type jsiiProxy_PhaseChangeEvent struct {
	_ byte // padding
}

func PhaseChangeEvent_BuildComplete() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.PhaseChangeEvent",
		"buildComplete",
		&returns,
	)
	return returns
}

func PhaseChangeEvent_BuildId() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.PhaseChangeEvent",
		"buildId",
		&returns,
	)
	return returns
}

func PhaseChangeEvent_CompletedPhase() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.PhaseChangeEvent",
		"completedPhase",
		&returns,
	)
	return returns
}

func PhaseChangeEvent_CompletedPhaseDurationSeconds() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.PhaseChangeEvent",
		"completedPhaseDurationSeconds",
		&returns,
	)
	return returns
}

func PhaseChangeEvent_CompletedPhaseStatus() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.PhaseChangeEvent",
		"completedPhaseStatus",
		&returns,
	)
	return returns
}

func PhaseChangeEvent_ProjectName() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.PhaseChangeEvent",
		"projectName",
		&returns,
	)
	return returns
}

// A convenience class for CodeBuild Projects that are used in CodePipeline.
// Experimental.
type PipelineProject interface {
	Project
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() awsiam.IPrincipal
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ProjectArn() *string
	ProjectName() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	AddFileSystemLocation(fileSystemLocation IFileSystemLocation)
	AddSecondaryArtifact(secondaryArtifact IArtifacts)
	AddSecondarySource(secondarySource ISource)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	BindAsNotificationRuleSource(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig
	BindToCodePipeline(_scope awscdk.Construct, options *BindToCodePipelineOptions)
	EnableBatchBuilds() *BatchBuildConfig
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricFailedBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSucceededBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *ProjectNotifyOnOptions) awscodestarnotifications.INotificationRule
	NotifyOnBuildFailed(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	NotifyOnBuildSucceeded(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	OnBuildFailed(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnBuildStarted(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnBuildSucceeded(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnPhaseChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnPrepare()
	OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for PipelineProject
type jsiiProxy_PipelineProject struct {
	jsiiProxy_Project
}

func (j *jsiiProxy_PipelineProject) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) ProjectArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"projectArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) ProjectName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"projectName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PipelineProject) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewPipelineProject(scope constructs.Construct, id *string, props *PipelineProjectProps) PipelineProject {
	_init_.Initialize()

	j := jsiiProxy_PipelineProject{}

	_jsii_.Create(
		"monocdk.aws_codebuild.PipelineProject",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewPipelineProject_Override(p PipelineProject, scope constructs.Construct, id *string, props *PipelineProjectProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.PipelineProject",
		[]interface{}{scope, id, props},
		p,
	)
}

// Experimental.
func PipelineProject_FromProjectArn(scope constructs.Construct, id *string, projectArn *string) IProject {
	_init_.Initialize()

	var returns IProject

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.PipelineProject",
		"fromProjectArn",
		[]interface{}{scope, id, projectArn},
		&returns,
	)

	return returns
}

// Import a Project defined either outside the CDK, or in a different CDK Stack (and exported using the {@link export} method).
//
// Returns: a reference to the existing Project
// Experimental.
func PipelineProject_FromProjectName(scope constructs.Construct, id *string, projectName *string) IProject {
	_init_.Initialize()

	var returns IProject

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.PipelineProject",
		"fromProjectName",
		[]interface{}{scope, id, projectName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func PipelineProject_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.PipelineProject",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func PipelineProject_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.PipelineProject",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Convert the environment variables map of string to {@link BuildEnvironmentVariable}, which is the customer-facing type, to a list of {@link CfnProject.EnvironmentVariableProperty}, which is the representation of environment variables in CloudFormation.
//
// Returns: an array of {@link CfnProject.EnvironmentVariableProperty} instances
// Experimental.
func PipelineProject_SerializeEnvVariables(environmentVariables *map[string]*BuildEnvironmentVariable, validateNoPlainTextSecrets *bool, principal awsiam.IGrantable) *[]*CfnProject_EnvironmentVariableProperty {
	_init_.Initialize()

	var returns *[]*CfnProject_EnvironmentVariableProperty

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.PipelineProject",
		"serializeEnvVariables",
		[]interface{}{environmentVariables, validateNoPlainTextSecrets, principal},
		&returns,
	)

	return returns
}

// Adds a fileSystemLocation to the Project.
// Experimental.
func (p *jsiiProxy_PipelineProject) AddFileSystemLocation(fileSystemLocation IFileSystemLocation) {
	_jsii_.InvokeVoid(
		p,
		"addFileSystemLocation",
		[]interface{}{fileSystemLocation},
	)
}

// Adds a secondary artifact to the Project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-multi-in-out.html
//
// Experimental.
func (p *jsiiProxy_PipelineProject) AddSecondaryArtifact(secondaryArtifact IArtifacts) {
	_jsii_.InvokeVoid(
		p,
		"addSecondaryArtifact",
		[]interface{}{secondaryArtifact},
	)
}

// Adds a secondary source to the Project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-multi-in-out.html
//
// Experimental.
func (p *jsiiProxy_PipelineProject) AddSecondarySource(secondarySource ISource) {
	_jsii_.InvokeVoid(
		p,
		"addSecondarySource",
		[]interface{}{secondarySource},
	)
}

// Add a permission only if there's a policy attached.
// Experimental.
func (p *jsiiProxy_PipelineProject) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		p,
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
func (p *jsiiProxy_PipelineProject) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		p,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Returns a source configuration for notification rule.
// Experimental.
func (p *jsiiProxy_PipelineProject) BindAsNotificationRuleSource(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig {
	var returns *awscodestarnotifications.NotificationRuleSourceConfig

	_jsii_.Invoke(
		p,
		"bindAsNotificationRuleSource",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// A callback invoked when the given project is added to a CodePipeline.
// Experimental.
func (p *jsiiProxy_PipelineProject) BindToCodePipeline(_scope awscdk.Construct, options *BindToCodePipelineOptions) {
	_jsii_.InvokeVoid(
		p,
		"bindToCodePipeline",
		[]interface{}{_scope, options},
	)
}

// Enable batch builds.
//
// Returns an object contining the batch service role if batch builds
// could be enabled.
// Experimental.
func (p *jsiiProxy_PipelineProject) EnableBatchBuilds() *BatchBuildConfig {
	var returns *BatchBuildConfig

	_jsii_.Invoke(
		p,
		"enableBatchBuilds",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (p *jsiiProxy_PipelineProject) GeneratePhysicalName() *string {
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
func (p *jsiiProxy_PipelineProject) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (p *jsiiProxy_PipelineProject) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns: a CloudWatch metric associated with this build project.
// Experimental.
func (p *jsiiProxy_PipelineProject) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Measures the number of builds triggered.
//
// Units: Count
//
// Valid CloudWatch statistics: Sum
// Experimental.
func (p *jsiiProxy_PipelineProject) MetricBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Measures the duration of all builds over time.
//
// Units: Seconds
//
// Valid CloudWatch statistics: Average (recommended), Maximum, Minimum
// Experimental.
func (p *jsiiProxy_PipelineProject) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Measures the number of builds that failed because of client error or because of a timeout.
//
// Units: Count
//
// Valid CloudWatch statistics: Sum
// Experimental.
func (p *jsiiProxy_PipelineProject) MetricFailedBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricFailedBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Measures the number of successful builds.
//
// Units: Count
//
// Valid CloudWatch statistics: Sum
// Experimental.
func (p *jsiiProxy_PipelineProject) MetricSucceededBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricSucceededBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Defines a CodeStar Notification rule triggered when the project events emitted by you specified, it very similar to `onEvent` API.
//
// You can also use the methods `notifyOnBuildSucceeded` and
// `notifyOnBuildFailed` to define rules for these specific event emitted.
// Experimental.
func (p *jsiiProxy_PipelineProject) NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *ProjectNotifyOnOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOn",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Defines a CodeStar notification rule which triggers when a build fails.
// Experimental.
func (p *jsiiProxy_PipelineProject) NotifyOnBuildFailed(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnBuildFailed",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Defines a CodeStar notification rule which triggers when a build completes successfully.
// Experimental.
func (p *jsiiProxy_PipelineProject) NotifyOnBuildSucceeded(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnBuildSucceeded",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Defines an event rule which triggers when a build fails.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// Experimental.
func (p *jsiiProxy_PipelineProject) OnBuildFailed(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onBuildFailed",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an event rule which triggers when a build starts.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// Experimental.
func (p *jsiiProxy_PipelineProject) OnBuildStarted(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onBuildStarted",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an event rule which triggers when a build completes successfully.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// Experimental.
func (p *jsiiProxy_PipelineProject) OnBuildSucceeded(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onBuildSucceeded",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule triggered when something happens with this project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
//
// Experimental.
func (p *jsiiProxy_PipelineProject) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule that triggers upon phase change of this build project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
//
// Experimental.
func (p *jsiiProxy_PipelineProject) OnPhaseChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onPhaseChange",
		[]interface{}{id, options},
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
func (p *jsiiProxy_PipelineProject) OnPrepare() {
	_jsii_.InvokeVoid(
		p,
		"onPrepare",
		nil, // no parameters
	)
}

// Defines a CloudWatch event rule triggered when the build project state changes.
//
// You can filter specific build status events using an event
// pattern filter on the `build-status` detail field:
//
//     const rule = project.onStateChange('OnBuildStarted', { target });
//     rule.addEventPattern({
//       detail: {
//         'build-status': [
//           "IN_PROGRESS",
//           "SUCCEEDED",
//           "FAILED",
//           "STOPPED"
//         ]
//       }
//     });
//
// You can also use the methods `onBuildFailed` and `onBuildSucceeded` to define rules for
// these specific state changes.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
//
// Experimental.
func (p *jsiiProxy_PipelineProject) OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onStateChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_PipelineProject) OnSynthesize(session constructs.ISynthesisSession) {
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
func (p *jsiiProxy_PipelineProject) OnValidate() *[]*string {
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
func (p *jsiiProxy_PipelineProject) Prepare() {
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
func (p *jsiiProxy_PipelineProject) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_PipelineProject) ToString() *string {
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
// Experimental.
func (p *jsiiProxy_PipelineProject) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type PipelineProjectProps struct {
	// Whether to allow the CodeBuild to send all network traffic.
	//
	// If set to false, you must individually add traffic rules to allow the
	// CodeBuild project to connect to network targets.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Indicates whether AWS CodeBuild generates a publicly accessible URL for your project's build badge.
	//
	// For more information, see Build Badges Sample
	// in the AWS CodeBuild User Guide.
	// Experimental.
	Badge *bool `json:"badge"`
	// Filename or contents of buildspec in JSON format.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec-ref-example
	//
	// Experimental.
	BuildSpec BuildSpec `json:"buildSpec"`
	// Caching strategy to use.
	// Experimental.
	Cache Cache `json:"cache"`
	// Whether to check for the presence of any secrets in the environment variables of the default type, BuildEnvironmentVariableType.PLAINTEXT. Since using a secret for the value of that kind of variable would result in it being displayed in plain text in the AWS Console, the construct will throw an exception if it detects a secret was passed there. Pass this property as false if you want to skip this validation, and keep using a secret in a plain text environment variable.
	// Experimental.
	CheckSecretsInPlainTextEnvVariables *bool `json:"checkSecretsInPlainTextEnvVariables"`
	// Maximum number of concurrent builds.
	//
	// Minimum value is 1 and maximum is account build limit.
	// Experimental.
	ConcurrentBuildLimit *float64 `json:"concurrentBuildLimit"`
	// A description of the project.
	//
	// Use the description to identify the purpose
	// of the project.
	// Experimental.
	Description *string `json:"description"`
	// Encryption key to use to read and write artifacts.
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// Build environment to use for the build.
	// Experimental.
	Environment *BuildEnvironment `json:"environment"`
	// Additional environment variables to add to the build environment.
	// Experimental.
	EnvironmentVariables *map[string]*BuildEnvironmentVariable `json:"environmentVariables"`
	// An  ProjectFileSystemLocation objects for a CodeBuild build project.
	//
	// A ProjectFileSystemLocation object specifies the identifier, location, mountOptions, mountPoint,
	// and type of a file system created using Amazon Elastic File System.
	// Experimental.
	FileSystemLocations *[]IFileSystemLocation `json:"fileSystemLocations"`
	// Add permissions to this project's role to create and use test report groups with name starting with the name of this project.
	//
	// That is the standard report group that gets created when a simple name
	// (in contrast to an ARN)
	// is used in the 'reports' section of the buildspec of this project.
	// This is usually harmless, but you can turn these off if you don't plan on using test
	// reports in this project.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/test-report-group-naming.html
	//
	// Experimental.
	GrantReportGroupPermissions *bool `json:"grantReportGroupPermissions"`
	// Information about logs for the build project.
	//
	// A project can create logs in Amazon CloudWatch Logs, an S3 bucket, or both.
	// Experimental.
	Logging *LoggingOptions `json:"logging"`
	// The physical, human-readable name of the CodeBuild Project.
	// Experimental.
	ProjectName *string `json:"projectName"`
	// The number of minutes after which AWS CodeBuild stops the build if it's still in queue.
	//
	// For valid values, see the timeoutInMinutes field in the AWS
	// CodeBuild User Guide.
	// Experimental.
	QueuedTimeout awscdk.Duration `json:"queuedTimeout"`
	// Service Role to assume while running the build.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// What security group to associate with the codebuild project's network interfaces.
	//
	// If no security group is identified, one will be created automatically.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// Where to place the network interfaces within the VPC.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SubnetSelection *awsec2.SubnetSelection `json:"subnetSelection"`
	// The number of minutes after which AWS CodeBuild stops the build if it's not complete.
	//
	// For valid values, see the timeoutInMinutes field in the AWS
	// CodeBuild User Guide.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// VPC network to place codebuild network interfaces.
	//
	// Specify this if the codebuild project needs to access resources in a VPC.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
}

// A representation of a CodeBuild Project.
// Experimental.
type Project interface {
	awscdk.Resource
	IProject
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() awsiam.IPrincipal
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ProjectArn() *string
	ProjectName() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	AddFileSystemLocation(fileSystemLocation IFileSystemLocation)
	AddSecondaryArtifact(secondaryArtifact IArtifacts)
	AddSecondarySource(secondarySource ISource)
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	BindAsNotificationRuleSource(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig
	BindToCodePipeline(_scope awscdk.Construct, options *BindToCodePipelineOptions)
	EnableBatchBuilds() *BatchBuildConfig
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricFailedBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricSucceededBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *ProjectNotifyOnOptions) awscodestarnotifications.INotificationRule
	NotifyOnBuildFailed(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	NotifyOnBuildSucceeded(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	OnBuildFailed(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnBuildStarted(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnBuildSucceeded(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnPhaseChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnPrepare()
	OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Project
type jsiiProxy_Project struct {
	internal.Type__awscdkResource
	jsiiProxy_IProject
}

func (j *jsiiProxy_Project) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) GrantPrincipal() awsiam.IPrincipal {
	var returns awsiam.IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) ProjectArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"projectArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) ProjectName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"projectName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Project) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewProject(scope constructs.Construct, id *string, props *ProjectProps) Project {
	_init_.Initialize()

	j := jsiiProxy_Project{}

	_jsii_.Create(
		"monocdk.aws_codebuild.Project",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewProject_Override(p Project, scope constructs.Construct, id *string, props *ProjectProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.Project",
		[]interface{}{scope, id, props},
		p,
	)
}

// Experimental.
func Project_FromProjectArn(scope constructs.Construct, id *string, projectArn *string) IProject {
	_init_.Initialize()

	var returns IProject

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Project",
		"fromProjectArn",
		[]interface{}{scope, id, projectArn},
		&returns,
	)

	return returns
}

// Import a Project defined either outside the CDK, or in a different CDK Stack (and exported using the {@link export} method).
//
// Returns: a reference to the existing Project
// Experimental.
func Project_FromProjectName(scope constructs.Construct, id *string, projectName *string) IProject {
	_init_.Initialize()

	var returns IProject

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Project",
		"fromProjectName",
		[]interface{}{scope, id, projectName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Project_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Project",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Project_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Project",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Convert the environment variables map of string to {@link BuildEnvironmentVariable}, which is the customer-facing type, to a list of {@link CfnProject.EnvironmentVariableProperty}, which is the representation of environment variables in CloudFormation.
//
// Returns: an array of {@link CfnProject.EnvironmentVariableProperty} instances
// Experimental.
func Project_SerializeEnvVariables(environmentVariables *map[string]*BuildEnvironmentVariable, validateNoPlainTextSecrets *bool, principal awsiam.IGrantable) *[]*CfnProject_EnvironmentVariableProperty {
	_init_.Initialize()

	var returns *[]*CfnProject_EnvironmentVariableProperty

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Project",
		"serializeEnvVariables",
		[]interface{}{environmentVariables, validateNoPlainTextSecrets, principal},
		&returns,
	)

	return returns
}

// Adds a fileSystemLocation to the Project.
// Experimental.
func (p *jsiiProxy_Project) AddFileSystemLocation(fileSystemLocation IFileSystemLocation) {
	_jsii_.InvokeVoid(
		p,
		"addFileSystemLocation",
		[]interface{}{fileSystemLocation},
	)
}

// Adds a secondary artifact to the Project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-multi-in-out.html
//
// Experimental.
func (p *jsiiProxy_Project) AddSecondaryArtifact(secondaryArtifact IArtifacts) {
	_jsii_.InvokeVoid(
		p,
		"addSecondaryArtifact",
		[]interface{}{secondaryArtifact},
	)
}

// Adds a secondary source to the Project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-multi-in-out.html
//
// Experimental.
func (p *jsiiProxy_Project) AddSecondarySource(secondarySource ISource) {
	_jsii_.InvokeVoid(
		p,
		"addSecondarySource",
		[]interface{}{secondarySource},
	)
}

// Add a permission only if there's a policy attached.
// Experimental.
func (p *jsiiProxy_Project) AddToRolePolicy(statement awsiam.PolicyStatement) {
	_jsii_.InvokeVoid(
		p,
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
func (p *jsiiProxy_Project) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		p,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Returns a source configuration for notification rule.
// Experimental.
func (p *jsiiProxy_Project) BindAsNotificationRuleSource(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig {
	var returns *awscodestarnotifications.NotificationRuleSourceConfig

	_jsii_.Invoke(
		p,
		"bindAsNotificationRuleSource",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// A callback invoked when the given project is added to a CodePipeline.
// Experimental.
func (p *jsiiProxy_Project) BindToCodePipeline(_scope awscdk.Construct, options *BindToCodePipelineOptions) {
	_jsii_.InvokeVoid(
		p,
		"bindToCodePipeline",
		[]interface{}{_scope, options},
	)
}

// Enable batch builds.
//
// Returns an object contining the batch service role if batch builds
// could be enabled.
// Experimental.
func (p *jsiiProxy_Project) EnableBatchBuilds() *BatchBuildConfig {
	var returns *BatchBuildConfig

	_jsii_.Invoke(
		p,
		"enableBatchBuilds",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (p *jsiiProxy_Project) GeneratePhysicalName() *string {
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
func (p *jsiiProxy_Project) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (p *jsiiProxy_Project) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns: a CloudWatch metric associated with this build project.
// Experimental.
func (p *jsiiProxy_Project) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Measures the number of builds triggered.
//
// Units: Count
//
// Valid CloudWatch statistics: Sum
// Experimental.
func (p *jsiiProxy_Project) MetricBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Measures the duration of all builds over time.
//
// Units: Seconds
//
// Valid CloudWatch statistics: Average (recommended), Maximum, Minimum
// Experimental.
func (p *jsiiProxy_Project) MetricDuration(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricDuration",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Measures the number of builds that failed because of client error or because of a timeout.
//
// Units: Count
//
// Valid CloudWatch statistics: Sum
// Experimental.
func (p *jsiiProxy_Project) MetricFailedBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricFailedBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Measures the number of successful builds.
//
// Units: Count
//
// Valid CloudWatch statistics: Sum
// Experimental.
func (p *jsiiProxy_Project) MetricSucceededBuilds(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		p,
		"metricSucceededBuilds",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Defines a CodeStar Notification rule triggered when the project events emitted by you specified, it very similar to `onEvent` API.
//
// You can also use the methods `notifyOnBuildSucceeded` and
// `notifyOnBuildFailed` to define rules for these specific event emitted.
// Experimental.
func (p *jsiiProxy_Project) NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *ProjectNotifyOnOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOn",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Defines a CodeStar notification rule which triggers when a build fails.
// Experimental.
func (p *jsiiProxy_Project) NotifyOnBuildFailed(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnBuildFailed",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Defines a CodeStar notification rule which triggers when a build completes successfully.
// Experimental.
func (p *jsiiProxy_Project) NotifyOnBuildSucceeded(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnBuildSucceeded",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Defines an event rule which triggers when a build fails.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// Experimental.
func (p *jsiiProxy_Project) OnBuildFailed(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onBuildFailed",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an event rule which triggers when a build starts.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// Experimental.
func (p *jsiiProxy_Project) OnBuildStarted(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onBuildStarted",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an event rule which triggers when a build completes successfully.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// Experimental.
func (p *jsiiProxy_Project) OnBuildSucceeded(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onBuildSucceeded",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule triggered when something happens with this project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
//
// Experimental.
func (p *jsiiProxy_Project) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule that triggers upon phase change of this build project.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
//
// Experimental.
func (p *jsiiProxy_Project) OnPhaseChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onPhaseChange",
		[]interface{}{id, options},
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
func (p *jsiiProxy_Project) OnPrepare() {
	_jsii_.InvokeVoid(
		p,
		"onPrepare",
		nil, // no parameters
	)
}

// Defines a CloudWatch event rule triggered when the build project state changes.
//
// You can filter specific build status events using an event
// pattern filter on the `build-status` detail field:
//
//     const rule = project.onStateChange('OnBuildStarted', { target });
//     rule.addEventPattern({
//       detail: {
//         'build-status': [
//           "IN_PROGRESS",
//           "SUCCEEDED",
//           "FAILED",
//           "STOPPED"
//         ]
//       }
//     });
//
// You can also use the methods `onBuildFailed` and `onBuildSucceeded` to define rules for
// these specific state changes.
//
// To access fields from the event in the event target input,
// use the static fields on the `StateChangeEvent` class.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html
//
// Experimental.
func (p *jsiiProxy_Project) OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onStateChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_Project) OnSynthesize(session constructs.ISynthesisSession) {
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
func (p *jsiiProxy_Project) OnValidate() *[]*string {
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
func (p *jsiiProxy_Project) Prepare() {
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
func (p *jsiiProxy_Project) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_Project) ToString() *string {
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
// Experimental.
func (p *jsiiProxy_Project) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The list of event types for AWS Codebuild.
// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#events-ref-buildproject
//
// Experimental.
type ProjectNotificationEvents string

const (
	ProjectNotificationEvents_BUILD_FAILED ProjectNotificationEvents = "BUILD_FAILED"
	ProjectNotificationEvents_BUILD_SUCCEEDED ProjectNotificationEvents = "BUILD_SUCCEEDED"
	ProjectNotificationEvents_BUILD_IN_PROGRESS ProjectNotificationEvents = "BUILD_IN_PROGRESS"
	ProjectNotificationEvents_BUILD_STOPPED ProjectNotificationEvents = "BUILD_STOPPED"
	ProjectNotificationEvents_BUILD_PHASE_FAILED ProjectNotificationEvents = "BUILD_PHASE_FAILED"
	ProjectNotificationEvents_BUILD_PHASE_SUCCEEDED ProjectNotificationEvents = "BUILD_PHASE_SUCCEEDED"
)

// Additional options to pass to the notification rule.
// Experimental.
type ProjectNotifyOnOptions struct {
	// The level of detail to include in the notifications for this resource.
	//
	// BASIC will include only the contents of the event as it would appear in AWS CloudWatch.
	// FULL will include any supplemental information provided by AWS CodeStar Notifications and/or the service for the resource for which the notification is created.
	// Experimental.
	DetailType awscodestarnotifications.DetailType `json:"detailType"`
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
	// A list of event types associated with this notification rule for CodeBuild Project.
	//
	// For a complete list of event types and IDs, see Notification concepts in the Developer Tools Console User Guide.
	// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#concepts-api
	//
	// Experimental.
	Events *[]ProjectNotificationEvents `json:"events"`
}

// Experimental.
type ProjectProps struct {
	// Whether to allow the CodeBuild to send all network traffic.
	//
	// If set to false, you must individually add traffic rules to allow the
	// CodeBuild project to connect to network targets.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	AllowAllOutbound *bool `json:"allowAllOutbound"`
	// Indicates whether AWS CodeBuild generates a publicly accessible URL for your project's build badge.
	//
	// For more information, see Build Badges Sample
	// in the AWS CodeBuild User Guide.
	// Experimental.
	Badge *bool `json:"badge"`
	// Filename or contents of buildspec in JSON format.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html#build-spec-ref-example
	//
	// Experimental.
	BuildSpec BuildSpec `json:"buildSpec"`
	// Caching strategy to use.
	// Experimental.
	Cache Cache `json:"cache"`
	// Whether to check for the presence of any secrets in the environment variables of the default type, BuildEnvironmentVariableType.PLAINTEXT. Since using a secret for the value of that kind of variable would result in it being displayed in plain text in the AWS Console, the construct will throw an exception if it detects a secret was passed there. Pass this property as false if you want to skip this validation, and keep using a secret in a plain text environment variable.
	// Experimental.
	CheckSecretsInPlainTextEnvVariables *bool `json:"checkSecretsInPlainTextEnvVariables"`
	// Maximum number of concurrent builds.
	//
	// Minimum value is 1 and maximum is account build limit.
	// Experimental.
	ConcurrentBuildLimit *float64 `json:"concurrentBuildLimit"`
	// A description of the project.
	//
	// Use the description to identify the purpose
	// of the project.
	// Experimental.
	Description *string `json:"description"`
	// Encryption key to use to read and write artifacts.
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// Build environment to use for the build.
	// Experimental.
	Environment *BuildEnvironment `json:"environment"`
	// Additional environment variables to add to the build environment.
	// Experimental.
	EnvironmentVariables *map[string]*BuildEnvironmentVariable `json:"environmentVariables"`
	// An  ProjectFileSystemLocation objects for a CodeBuild build project.
	//
	// A ProjectFileSystemLocation object specifies the identifier, location, mountOptions, mountPoint,
	// and type of a file system created using Amazon Elastic File System.
	// Experimental.
	FileSystemLocations *[]IFileSystemLocation `json:"fileSystemLocations"`
	// Add permissions to this project's role to create and use test report groups with name starting with the name of this project.
	//
	// That is the standard report group that gets created when a simple name
	// (in contrast to an ARN)
	// is used in the 'reports' section of the buildspec of this project.
	// This is usually harmless, but you can turn these off if you don't plan on using test
	// reports in this project.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/test-report-group-naming.html
	//
	// Experimental.
	GrantReportGroupPermissions *bool `json:"grantReportGroupPermissions"`
	// Information about logs for the build project.
	//
	// A project can create logs in Amazon CloudWatch Logs, an S3 bucket, or both.
	// Experimental.
	Logging *LoggingOptions `json:"logging"`
	// The physical, human-readable name of the CodeBuild Project.
	// Experimental.
	ProjectName *string `json:"projectName"`
	// The number of minutes after which AWS CodeBuild stops the build if it's still in queue.
	//
	// For valid values, see the timeoutInMinutes field in the AWS
	// CodeBuild User Guide.
	// Experimental.
	QueuedTimeout awscdk.Duration `json:"queuedTimeout"`
	// Service Role to assume while running the build.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// What security group to associate with the codebuild project's network interfaces.
	//
	// If no security group is identified, one will be created automatically.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SecurityGroups *[]awsec2.ISecurityGroup `json:"securityGroups"`
	// Where to place the network interfaces within the VPC.
	//
	// Only used if 'vpc' is supplied.
	// Experimental.
	SubnetSelection *awsec2.SubnetSelection `json:"subnetSelection"`
	// The number of minutes after which AWS CodeBuild stops the build if it's not complete.
	//
	// For valid values, see the timeoutInMinutes field in the AWS
	// CodeBuild User Guide.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// VPC network to place codebuild network interfaces.
	//
	// Specify this if the codebuild project needs to access resources in a VPC.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Defines where build artifacts will be stored.
	//
	// Could be: PipelineBuildArtifacts, NoArtifacts and S3Artifacts.
	// Experimental.
	Artifacts IArtifacts `json:"artifacts"`
	// The secondary artifacts for the Project.
	//
	// Can also be added after the Project has been created by using the {@link Project#addSecondaryArtifact} method.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-multi-in-out.html
	//
	// Experimental.
	SecondaryArtifacts *[]IArtifacts `json:"secondaryArtifacts"`
	// The secondary sources for the Project.
	//
	// Can be also added after the Project has been created by using the {@link Project#addSecondarySource} method.
	// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-multi-in-out.html
	//
	// Experimental.
	SecondarySources *[]ISource `json:"secondarySources"`
	// The source of the build.
	//
	// *Note*: if {@link NoSource} is given as the source,
	// then you need to provide an explicit `buildSpec`.
	// Experimental.
	Source ISource `json:"source"`
}

// The ReportGroup resource class.
// Experimental.
type ReportGroup interface {
	awscdk.Resource
	IReportGroup
	Env() *awscdk.ResourceEnvironment
	ExportBucket() awss3.IBucket
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ReportGroupArn() *string
	ReportGroupName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantWrite(identity awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ReportGroup
type jsiiProxy_ReportGroup struct {
	internal.Type__awscdkResource
	jsiiProxy_IReportGroup
}

func (j *jsiiProxy_ReportGroup) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ReportGroup) ExportBucket() awss3.IBucket {
	var returns awss3.IBucket
	_jsii_.Get(
		j,
		"exportBucket",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ReportGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ReportGroup) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ReportGroup) ReportGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"reportGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ReportGroup) ReportGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"reportGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ReportGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewReportGroup(scope constructs.Construct, id *string, props *ReportGroupProps) ReportGroup {
	_init_.Initialize()

	j := jsiiProxy_ReportGroup{}

	_jsii_.Create(
		"monocdk.aws_codebuild.ReportGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewReportGroup_Override(r ReportGroup, scope constructs.Construct, id *string, props *ReportGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.ReportGroup",
		[]interface{}{scope, id, props},
		r,
	)
}

// Reference an existing ReportGroup, defined outside of the CDK code, by name.
// Experimental.
func ReportGroup_FromReportGroupName(scope constructs.Construct, id *string, reportGroupName *string) IReportGroup {
	_init_.Initialize()

	var returns IReportGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.ReportGroup",
		"fromReportGroupName",
		[]interface{}{scope, id, reportGroupName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ReportGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.ReportGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ReportGroup_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.ReportGroup",
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
func (r *jsiiProxy_ReportGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_ReportGroup) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_ReportGroup) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_ReportGroup) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grants the given entity permissions to write (that is, upload reports to) this report group.
// Experimental.
func (r *jsiiProxy_ReportGroup) GrantWrite(identity awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grantWrite",
		[]interface{}{identity},
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
func (r *jsiiProxy_ReportGroup) OnPrepare() {
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
func (r *jsiiProxy_ReportGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_ReportGroup) OnValidate() *[]*string {
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
func (r *jsiiProxy_ReportGroup) Prepare() {
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
func (r *jsiiProxy_ReportGroup) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_ReportGroup) ToString() *string {
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
func (r *jsiiProxy_ReportGroup) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for {@link ReportGroup}.
// Experimental.
type ReportGroupProps struct {
	// An optional S3 bucket to export the reports to.
	// Experimental.
	ExportBucket awss3.IBucket `json:"exportBucket"`
	// What to do when this resource is deleted from a stack.
	//
	// As CodeBuild does not allow deleting a ResourceGroup that has reports inside of it,
	// this is set to retain the resource by default.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// The physical name of the report group.
	// Experimental.
	ReportGroupName *string `json:"reportGroupName"`
	// Whether to output the report files into the export bucket as-is, or create a ZIP from them before doing the export.
	//
	// Ignored if {@link exportBucket} has not been provided.
	// Experimental.
	ZipExport *bool `json:"zipExport"`
}

// Construction properties for {@link S3Artifacts}.
// Experimental.
type S3ArtifactsProps struct {
	// The artifact identifier.
	//
	// This property is required on secondary artifacts.
	// Experimental.
	Identifier *string `json:"identifier"`
	// The name of the output bucket.
	// Experimental.
	Bucket awss3.IBucket `json:"bucket"`
	// If this is false, build output will not be encrypted.
	//
	// This is useful if the artifact to publish a static website or sharing content with others
	// Experimental.
	Encryption *bool `json:"encryption"`
	// Indicates if the build ID should be included in the path.
	//
	// If this is set to true,
	// then the build artifact will be stored in "<path>/<build-id>/<name>".
	// Experimental.
	IncludeBuildId *bool `json:"includeBuildId"`
	// The name of the build output ZIP file or folder inside the bucket.
	//
	// The full S3 object key will be "<path>/<build-id>/<name>" or
	// "<path>/<name>" depending on whether `includeBuildId` is set to true.
	//
	// If not set, `overrideArtifactName` will be set and the name from the
	// buildspec will be used instead.
	// Experimental.
	Name *string `json:"name"`
	// If this is true, all build output will be packaged into a single .zip file. Otherwise, all files will be uploaded to <path>/<name>.
	// Experimental.
	PackageZip *bool `json:"packageZip"`
	// The path inside of the bucket for the build output .zip file or folder. If a value is not specified, then build output will be stored at the root of the bucket (or under the <build-id> directory if `includeBuildId` is set to true).
	// Experimental.
	Path *string `json:"path"`
}

// Information about logs built to an S3 bucket for a build project.
// Experimental.
type S3LoggingOptions struct {
	// The S3 Bucket to send logs to.
	// Experimental.
	Bucket awss3.IBucket `json:"bucket"`
	// The current status of the logs in Amazon CloudWatch Logs for a build project.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// Encrypt the S3 build log output.
	// Experimental.
	Encrypted *bool `json:"encrypted"`
	// The path prefix for S3 logs.
	// Experimental.
	Prefix *string `json:"prefix"`
}

// Construction properties for {@link S3Source}.
// Experimental.
type S3SourceProps struct {
	// The source identifier.
	//
	// This property is required on secondary sources.
	// Experimental.
	Identifier *string `json:"identifier"`
	// Experimental.
	Bucket awss3.IBucket `json:"bucket"`
	// Experimental.
	Path *string `json:"path"`
	// The version ID of the object that represents the build input ZIP file to use.
	// Experimental.
	Version *string `json:"version"`
}

// Source provider definition for a CodeBuild Project.
// Experimental.
type Source interface {
	ISource
	BadgeSupported() *bool
	Identifier() *string
	Type() *string
	Bind(_scope awscdk.Construct, _project IProject) *SourceConfig
}

// The jsii proxy struct for Source
type jsiiProxy_Source struct {
	jsiiProxy_ISource
}

func (j *jsiiProxy_Source) BadgeSupported() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"badgeSupported",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Source) Identifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Source) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Experimental.
func NewSource_Override(s Source, props *SourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.Source",
		[]interface{}{props},
		s,
	)
}

// Experimental.
func Source_BitBucket(props *BitBucketSourceProps) ISource {
	_init_.Initialize()

	var returns ISource

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Source",
		"bitBucket",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Experimental.
func Source_CodeCommit(props *CodeCommitSourceProps) ISource {
	_init_.Initialize()

	var returns ISource

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Source",
		"codeCommit",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Experimental.
func Source_GitHub(props *GitHubSourceProps) ISource {
	_init_.Initialize()

	var returns ISource

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Source",
		"gitHub",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Experimental.
func Source_GitHubEnterprise(props *GitHubEnterpriseSourceProps) ISource {
	_init_.Initialize()

	var returns ISource

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Source",
		"gitHubEnterprise",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Experimental.
func Source_S3(props *S3SourceProps) ISource {
	_init_.Initialize()

	var returns ISource

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.Source",
		"s3",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Called by the project when the source is added so that the source can perform binding operations on the source.
//
// For example, it can grant permissions to the
// code build project to read from the S3 bucket.
// Experimental.
func (s *jsiiProxy_Source) Bind(_scope awscdk.Construct, _project IProject) *SourceConfig {
	var returns *SourceConfig

	_jsii_.Invoke(
		s,
		"bind",
		[]interface{}{_scope, _project},
		&returns,
	)

	return returns
}

// The type returned from {@link ISource#bind}.
// Experimental.
type SourceConfig struct {
	// Experimental.
	SourceProperty *CfnProject_SourceProperty `json:"sourceProperty"`
	// Experimental.
	BuildTriggers *CfnProject_ProjectTriggersProperty `json:"buildTriggers"`
	// `AWS::CodeBuild::Project.SourceVersion`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-codebuild-project.html#cfn-codebuild-project-sourceversion
	//
	// Experimental.
	SourceVersion *string `json:"sourceVersion"`
}

// Properties common to all Source classes.
// Experimental.
type SourceProps struct {
	// The source identifier.
	//
	// This property is required on secondary sources.
	// Experimental.
	Identifier *string `json:"identifier"`
}

// Event fields for the CodeBuild "state change" event.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-build-notifications.html#sample-build-notifications-ref
//
// Experimental.
type StateChangeEvent interface {
}

// The jsii proxy struct for StateChangeEvent
type jsiiProxy_StateChangeEvent struct {
	_ byte // padding
}

func StateChangeEvent_BuildId() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.StateChangeEvent",
		"buildId",
		&returns,
	)
	return returns
}

func StateChangeEvent_BuildStatus() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.StateChangeEvent",
		"buildStatus",
		&returns,
	)
	return returns
}

func StateChangeEvent_CurrentPhase() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.StateChangeEvent",
		"currentPhase",
		&returns,
	)
	return returns
}

func StateChangeEvent_ProjectName() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.StateChangeEvent",
		"projectName",
		&returns,
	)
	return returns
}

// Permissions Boundary for a CodeBuild Project running untrusted code.
//
// This class is a Policy, intended to be used as a Permissions Boundary
// for a CodeBuild project. It allows most of the actions necessary to run
// the CodeBuild project, but disallows reading from Parameter Store
// and Secrets Manager.
//
// Use this when your CodeBuild project is running untrusted code (for
// example, if you are using one to automatically build Pull Requests
// that anyone can submit), and you want to prevent your future self
// from accidentally exposing Secrets to this build.
//
// (The reason you might want to do this is because otherwise anyone
// who can submit a Pull Request to your project can write a script
// to email those secrets to themselves).
//
// TODO: EXAMPLE
//
// Experimental.
type UntrustedCodeBoundaryPolicy interface {
	awsiam.ManagedPolicy
	Description() *string
	Document() awsiam.PolicyDocument
	Env() *awscdk.ResourceEnvironment
	ManagedPolicyArn() *string
	ManagedPolicyName() *string
	Node() awscdk.ConstructNode
	Path() *string
	PhysicalName() *string
	Stack() awscdk.Stack
	AddStatements(statement ...awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AttachToGroup(group awsiam.IGroup)
	AttachToRole(role awsiam.IRole)
	AttachToUser(user awsiam.IUser)
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

// The jsii proxy struct for UntrustedCodeBoundaryPolicy
type jsiiProxy_UntrustedCodeBoundaryPolicy struct {
	internal.Type__awsiamManagedPolicy
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) Document() awsiam.PolicyDocument {
	var returns awsiam.PolicyDocument
	_jsii_.Get(
		j,
		"document",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) ManagedPolicyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"managedPolicyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) ManagedPolicyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"managedPolicyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UntrustedCodeBoundaryPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewUntrustedCodeBoundaryPolicy(scope constructs.Construct, id *string, props *UntrustedCodeBoundaryPolicyProps) UntrustedCodeBoundaryPolicy {
	_init_.Initialize()

	j := jsiiProxy_UntrustedCodeBoundaryPolicy{}

	_jsii_.Create(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUntrustedCodeBoundaryPolicy_Override(u UntrustedCodeBoundaryPolicy, scope constructs.Construct, id *string, props *UntrustedCodeBoundaryPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		[]interface{}{scope, id, props},
		u,
	)
}

// Import a managed policy from one of the policies that AWS manages.
//
// For this managed policy, you only need to know the name to be able to use it.
//
// Some managed policy names start with "service-role/", some start with
// "job-function/", and some don't start with anything. Do include the
// prefix when constructing this object.
// Experimental.
func UntrustedCodeBoundaryPolicy_FromAwsManagedPolicyName(managedPolicyName *string) awsiam.IManagedPolicy {
	_init_.Initialize()

	var returns awsiam.IManagedPolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		"fromAwsManagedPolicyName",
		[]interface{}{managedPolicyName},
		&returns,
	)

	return returns
}

// Import an external managed policy by ARN.
//
// For this managed policy, you only need to know the ARN to be able to use it.
// This can be useful if you got the ARN from a CloudFormation Export.
//
// If the imported Managed Policy ARN is a Token (such as a
// `CfnParameter.valueAsString` or a `Fn.importValue()`) *and* the referenced
// managed policy has a `path` (like `arn:...:policy/AdminPolicy/AdminAllow`), the
// `managedPolicyName` property will not resolve to the correct value. Instead it
// will resolve to the first path component. We unfortunately cannot express
// the correct calculation of the full path name as a CloudFormation
// expression. In this scenario the Managed Policy ARN should be supplied without the
// `path` in order to resolve the correct managed policy resource.
// Experimental.
func UntrustedCodeBoundaryPolicy_FromManagedPolicyArn(scope constructs.Construct, id *string, managedPolicyArn *string) awsiam.IManagedPolicy {
	_init_.Initialize()

	var returns awsiam.IManagedPolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		"fromManagedPolicyArn",
		[]interface{}{scope, id, managedPolicyArn},
		&returns,
	)

	return returns
}

// Import a customer managed policy from the managedPolicyName.
//
// For this managed policy, you only need to know the name to be able to use it.
// Experimental.
func UntrustedCodeBoundaryPolicy_FromManagedPolicyName(scope constructs.Construct, id *string, managedPolicyName *string) awsiam.IManagedPolicy {
	_init_.Initialize()

	var returns awsiam.IManagedPolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		"fromManagedPolicyName",
		[]interface{}{scope, id, managedPolicyName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func UntrustedCodeBoundaryPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UntrustedCodeBoundaryPolicy_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.UntrustedCodeBoundaryPolicy",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a statement to the policy document.
// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) AddStatements(statement ...awsiam.PolicyStatement) {
	args := []interface{}{}
	for _, a := range statement {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		u,
		"addStatements",
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
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attaches this policy to a group.
// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) AttachToGroup(group awsiam.IGroup) {
	_jsii_.InvokeVoid(
		u,
		"attachToGroup",
		[]interface{}{group},
	)
}

// Attaches this policy to a role.
// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) AttachToRole(role awsiam.IRole) {
	_jsii_.InvokeVoid(
		u,
		"attachToRole",
		[]interface{}{role},
	)
}

// Attaches this policy to a user.
// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) AttachToUser(user awsiam.IUser) {
	_jsii_.InvokeVoid(
		u,
		"attachToUser",
		[]interface{}{user},
	)
}

// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		u,
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
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		u,
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
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		u,
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
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) OnPrepare() {
	_jsii_.InvokeVoid(
		u,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
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
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
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
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) Prepare() {
	_jsii_.InvokeVoid(
		u,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		u,
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
func (u *jsiiProxy_UntrustedCodeBoundaryPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for UntrustedCodeBoundaryPolicy.
// Experimental.
type UntrustedCodeBoundaryPolicyProps struct {
	// Additional statements to add to the default set of statements.
	// Experimental.
	AdditionalStatements *[]awsiam.PolicyStatement `json:"additionalStatements"`
	// The name of the managed policy.
	// Experimental.
	ManagedPolicyName *string `json:"managedPolicyName"`
}

// A CodeBuild image running Windows.
//
// This class has a bunch of public constants that represent the most popular images.
//
// You can also specify a custom image using one of the static methods:
//
// - WindowsBuildImage.fromDockerRegistry(image[, { secretsManagerCredentials }, imageType])
// - WindowsBuildImage.fromEcrRepository(repo[, tag, imageType])
// - WindowsBuildImage.fromAsset(parent, id, props, [, imageType])
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-available.html
//
// Experimental.
type WindowsBuildImage interface {
	IBuildImage
	DefaultComputeType() ComputeType
	ImageId() *string
	ImagePullPrincipalType() ImagePullPrincipalType
	Repository() awsecr.IRepository
	SecretsManagerCredentials() awssecretsmanager.ISecret
	Type() *string
	RunScriptBuildspec(entrypoint *string) BuildSpec
	Validate(buildEnvironment *BuildEnvironment) *[]*string
}

// The jsii proxy struct for WindowsBuildImage
type jsiiProxy_WindowsBuildImage struct {
	jsiiProxy_IBuildImage
}

func (j *jsiiProxy_WindowsBuildImage) DefaultComputeType() ComputeType {
	var returns ComputeType
	_jsii_.Get(
		j,
		"defaultComputeType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WindowsBuildImage) ImageId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WindowsBuildImage) ImagePullPrincipalType() ImagePullPrincipalType {
	var returns ImagePullPrincipalType
	_jsii_.Get(
		j,
		"imagePullPrincipalType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WindowsBuildImage) Repository() awsecr.IRepository {
	var returns awsecr.IRepository
	_jsii_.Get(
		j,
		"repository",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WindowsBuildImage) SecretsManagerCredentials() awssecretsmanager.ISecret {
	var returns awssecretsmanager.ISecret
	_jsii_.Get(
		j,
		"secretsManagerCredentials",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WindowsBuildImage) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Uses an Docker image asset as a Windows build image.
// Experimental.
func WindowsBuildImage_FromAsset(scope constructs.Construct, id *string, props *awsecrassets.DockerImageAssetProps, imageType WindowsImageType) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.WindowsBuildImage",
		"fromAsset",
		[]interface{}{scope, id, props, imageType},
		&returns,
	)

	return returns
}

// Returns: a Windows build image from a Docker Hub image.
// Experimental.
func WindowsBuildImage_FromDockerRegistry(name *string, options *DockerImageOptions, imageType WindowsImageType) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.WindowsBuildImage",
		"fromDockerRegistry",
		[]interface{}{name, options, imageType},
		&returns,
	)

	return returns
}

// Returns: A Linux build image from an ECR repository.
//
// NOTE: if the repository is external (i.e. imported), then we won't be able to add
// a resource policy statement for it so CodeBuild can pull the image.
// See: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-ecr.html
//
// Experimental.
func WindowsBuildImage_FromEcrRepository(repository awsecr.IRepository, tag *string, imageType WindowsImageType) IBuildImage {
	_init_.Initialize()

	var returns IBuildImage

	_jsii_.StaticInvoke(
		"monocdk.aws_codebuild.WindowsBuildImage",
		"fromEcrRepository",
		[]interface{}{repository, tag, imageType},
		&returns,
	)

	return returns
}

func WindowsBuildImage_WIN_SERVER_CORE_2016_BASE() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.WindowsBuildImage",
		"WIN_SERVER_CORE_2016_BASE",
		&returns,
	)
	return returns
}

func WindowsBuildImage_WIN_SERVER_CORE_2019_BASE() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.WindowsBuildImage",
		"WIN_SERVER_CORE_2019_BASE",
		&returns,
	)
	return returns
}

func WindowsBuildImage_WINDOWS_BASE_2_0() IBuildImage {
	_init_.Initialize()
	var returns IBuildImage
	_jsii_.StaticGet(
		"monocdk.aws_codebuild.WindowsBuildImage",
		"WINDOWS_BASE_2_0",
		&returns,
	)
	return returns
}

// Make a buildspec to run the indicated script.
// Experimental.
func (w *jsiiProxy_WindowsBuildImage) RunScriptBuildspec(entrypoint *string) BuildSpec {
	var returns BuildSpec

	_jsii_.Invoke(
		w,
		"runScriptBuildspec",
		[]interface{}{entrypoint},
		&returns,
	)

	return returns
}

// Allows the image a chance to validate whether the passed configuration is correct.
// Experimental.
func (w *jsiiProxy_WindowsBuildImage) Validate(buildEnvironment *BuildEnvironment) *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		w,
		"validate",
		[]interface{}{buildEnvironment},
		&returns,
	)

	return returns
}

// Environment type for Windows Docker images.
// Experimental.
type WindowsImageType string

const (
	WindowsImageType_STANDARD WindowsImageType = "STANDARD"
	WindowsImageType_SERVER_2019 WindowsImageType = "SERVER_2019"
)


// An experiment to bundle the entire CDK into a single module
package awscdk

import (
	"time"

	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk/cloudassemblyschema"
	"github.com/aws/aws-cdk-go/awscdk/cxapi"
	"github.com/aws/aws-cdk-go/awscdk/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// Includes API for attaching annotations such as warning messages to constructs.
// Experimental.
type Annotations interface {
	AddDeprecation(api *string, message *string)
	AddError(message *string)
	AddInfo(message *string)
	AddWarning(message *string)
}

// The jsii proxy struct for Annotations
type jsiiProxy_Annotations struct {
	_ byte // padding
}

// Returns the annotations API for a construct scope.
// Experimental.
func Annotations_Of(scope constructs.IConstruct) Annotations {
	_init_.Initialize()

	var returns Annotations

	_jsii_.StaticInvoke(
		"monocdk.Annotations",
		"of",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Adds a deprecation warning for a specific API.
//
// Deprecations will be added only once per construct as a warning and will be
// deduplicated based on the `api`.
//
// If the environment variable `CDK_BLOCK_DEPRECATIONS` is set, this method
// will throw an error instead with the deprecation message.
// Experimental.
func (a *jsiiProxy_Annotations) AddDeprecation(api *string, message *string) {
	_jsii_.InvokeVoid(
		a,
		"addDeprecation",
		[]interface{}{api, message},
	)
}

// Adds an { "error": <message> } metadata entry to this construct.
//
// The toolkit will fail synthesis when errors are reported.
// Experimental.
func (a *jsiiProxy_Annotations) AddError(message *string) {
	_jsii_.InvokeVoid(
		a,
		"addError",
		[]interface{}{message},
	)
}

// Adds an info metadata entry to this construct.
//
// The CLI will display the info message when apps are synthesized.
// Experimental.
func (a *jsiiProxy_Annotations) AddInfo(message *string) {
	_jsii_.InvokeVoid(
		a,
		"addInfo",
		[]interface{}{message},
	)
}

// Adds a warning metadata entry to this construct.
//
// The CLI will display the warning when an app is synthesized, or fail if run
// in --strict mode.
// Experimental.
func (a *jsiiProxy_Annotations) AddWarning(message *string) {
	_jsii_.InvokeVoid(
		a,
		"addWarning",
		[]interface{}{message},
	)
}

// A construct which represents an entire CDK app. This construct is normally the root of the construct tree.
//
// You would normally define an `App` instance in your program's entrypoint,
// then define constructs where the app is used as the parent scope.
//
// After all the child constructs are defined within the app, you should call
// `app.synth()` which will emit a "cloud assembly" from this app into the
// directory specified by `outdir`. Cloud assemblies includes artifacts such as
// CloudFormation templates and assets that are needed to deploy this app into
// the AWS cloud.
// See: https://docs.aws.amazon.com/cdk/latest/guide/apps.html
//
// Experimental.
type App interface {
	Stage
	Account() *string
	ArtifactId() *string
	AssetOutdir() *string
	Node() ConstructNode
	Outdir() *string
	ParentStage() Stage
	Region() *string
	StageName() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synth(options *StageSynthesisOptions) cxapi.CloudAssembly
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for App
type jsiiProxy_App struct {
	jsiiProxy_Stage
}

func (j *jsiiProxy_App) Account() *string {
	var returns *string
	_jsii_.Get(
		j,
		"account",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_App) ArtifactId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"artifactId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_App) AssetOutdir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assetOutdir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_App) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_App) Outdir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outdir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_App) ParentStage() Stage {
	var returns Stage
	_jsii_.Get(
		j,
		"parentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_App) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_App) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}


// Initializes a CDK application.
// Experimental.
func NewApp(props *AppProps) App {
	_init_.Initialize()

	j := jsiiProxy_App{}

	_jsii_.Create(
		"monocdk.App",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Initializes a CDK application.
// Experimental.
func NewApp_Override(a App, props *AppProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.App",
		[]interface{}{props},
		a,
	)
}

// Checks if an object is an instance of the `App` class.
//
// Returns: `true` if `obj` is an `App`.
// Experimental.
func App_IsApp(obj interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.App",
		"isApp",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func App_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.App",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Test whether the given construct is a stage.
// Experimental.
func App_IsStage(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.App",
		"isStage",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return the stage this construct is contained with, if available.
//
// If called
// on a nested stage, returns its parent.
// Experimental.
func App_Of(construct constructs.IConstruct) Stage {
	_init_.Initialize()

	var returns Stage

	_jsii_.StaticInvoke(
		"monocdk.App",
		"of",
		[]interface{}{construct},
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
func (a *jsiiProxy_App) OnPrepare() {
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
func (a *jsiiProxy_App) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_App) OnValidate() *[]*string {
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
func (a *jsiiProxy_App) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Synthesize this stage into a cloud assembly.
//
// Once an assembly has been synthesized, it cannot be modified. Subsequent
// calls will return the same assembly.
// Experimental.
func (a *jsiiProxy_App) Synth(options *StageSynthesisOptions) cxapi.CloudAssembly {
	var returns cxapi.CloudAssembly

	_jsii_.Invoke(
		a,
		"synth",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_App) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_App) ToString() *string {
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
func (a *jsiiProxy_App) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Initialization props for apps.
// Experimental.
type AppProps struct {
	// Include runtime versioning information in the Stacks of this app.
	// Experimental.
	AnalyticsReporting *bool `json:"analyticsReporting"`
	// Automatically call `synth()` before the program exits.
	//
	// If you set this, you don't have to call `synth()` explicitly. Note that
	// this feature is only available for certain programming languages, and
	// calling `synth()` is still recommended.
	// Experimental.
	AutoSynth *bool `json:"autoSynth"`
	// Additional context values for the application.
	//
	// Context set by the CLI or the `context` key in `cdk.json` has precedence.
	//
	// Context can be read from any construct using `node.getContext(key)`.
	// Experimental.
	Context *map[string]interface{} `json:"context"`
	// The output directory into which to emit synthesized artifacts.
	// Experimental.
	Outdir *string `json:"outdir"`
	// Include runtime versioning information in the Stacks of this app.
	// Deprecated: use `versionReporting` instead
	RuntimeInfo *bool `json:"runtimeInfo"`
	// Include construct creation stack trace in the `aws:cdk:trace` metadata key of all constructs.
	// Experimental.
	StackTraces *bool `json:"stackTraces"`
	// Include construct tree metadata as part of the Cloud Assembly.
	// Experimental.
	TreeMetadata *bool `json:"treeMetadata"`
}

// Experimental.
type Arn interface {
}

// The jsii proxy struct for Arn
type jsiiProxy_Arn struct {
	_ byte // padding
}

// Extract the full resource name from an ARN.
//
// Necessary for resource names (paths) that may contain the separator, like
// `arn:aws:iam::111111111111:role/path/to/role/name`.
//
// Only works if we statically know the expected `resourceType` beforehand, since we're going
// to use that to split the string on ':<resourceType>/' (and take the right-hand side).
//
// We can't extract the 'resourceType' from the ARN at hand, because CloudFormation Expressions
// only allow literals in the 'separator' argument to `{ Fn::Split }`, and so it can't be
// `{ Fn::Select: [5, { Fn::Split: [':', ARN] }}`.
//
// Only necessary for ARN formats for which the type-name separator is `/`.
// Experimental.
func Arn_ExtractResourceName(arn *string, resourceType *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Arn",
		"extractResourceName",
		[]interface{}{arn, resourceType},
		&returns,
	)

	return returns
}

// Creates an ARN from components.
//
// If `partition`, `region` or `account` are not specified, the stack's
// partition, region and account will be used.
//
// If any component is the empty string, an empty string will be inserted
// into the generated ARN at the location that component corresponds to.
//
// The ARN will be formatted as follows:
//
//    arn:{partition}:{service}:{region}:{account}:{resource}{sep}{resource-name}
//
// The required ARN pieces that are omitted will be taken from the stack that
// the 'scope' is attached to. If all ARN pieces are supplied, the supplied scope
// can be 'undefined'.
// Experimental.
func Arn_Format(components *ArnComponents, stack Stack) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Arn",
		"format",
		[]interface{}{components, stack},
		&returns,
	)

	return returns
}

// Given an ARN, parses it and returns components.
//
// IF THE ARN IS A CONCRETE STRING...
//
// ...it will be parsed and validated. The separator (`sep`) will be set to '/'
// if the 6th component includes a '/', in which case, `resource` will be set
// to the value before the '/' and `resourceName` will be the rest. In case
// there is no '/', `resource` will be set to the 6th components and
// `resourceName` will be set to the rest of the string.
//
// IF THE ARN IS A TOKEN...
//
// ...it cannot be validated, since we don't have the actual value yet at the
// time of this function call. You will have to supply `sepIfToken` and
// whether or not ARNs of the expected format usually have resource names
// in order to parse it properly. The resulting `ArnComponents` object will
// contain tokens for the subexpressions of the ARN, not string literals.
//
// If the resource name could possibly contain the separator char, the actual
// resource name cannot be properly parsed. This only occurs if the separator
// char is '/', and happens for example for S3 object ARNs, IAM Role ARNs,
// IAM OIDC Provider ARNs, etc. To properly extract the resource name from a
// Tokenized ARN, you must know the resource type and call
// `Arn.extractResourceName`.
//
// Returns: an ArnComponents object which allows access to the various
// components of the ARN.
// Deprecated: use split instead
func Arn_Parse(arn *string, sepIfToken *string, hasName *bool) *ArnComponents {
	_init_.Initialize()

	var returns *ArnComponents

	_jsii_.StaticInvoke(
		"monocdk.Arn",
		"parse",
		[]interface{}{arn, sepIfToken, hasName},
		&returns,
	)

	return returns
}

// Splits the provided ARN into its components.
//
// Works both if 'arn' is a string like 'arn:aws:s3:::bucket',
// and a Token representing a dynamic CloudFormation expression
// (in which case the returned components will also be dynamic CloudFormation expressions,
// encoded as Tokens).
// Experimental.
func Arn_Split(arn *string, arnFormat ArnFormat) *ArnComponents {
	_init_.Initialize()

	var returns *ArnComponents

	_jsii_.StaticInvoke(
		"monocdk.Arn",
		"split",
		[]interface{}{arn, arnFormat},
		&returns,
	)

	return returns
}

// Experimental.
type ArnComponents struct {
	// Resource type (e.g. "table", "autoScalingGroup", "certificate"). For some resource types, e.g. S3 buckets, this field defines the bucket name.
	// Experimental.
	Resource *string `json:"resource"`
	// The service namespace that identifies the AWS product (for example, 's3', 'iam', 'codepipline').
	// Experimental.
	Service *string `json:"service"`
	// The ID of the AWS account that owns the resource, without the hyphens.
	//
	// For example, 123456789012. Note that the ARNs for some resources don't
	// require an account number, so this component might be omitted.
	// Experimental.
	Account *string `json:"account"`
	// The specific ARN format to use for this ARN value.
	// Experimental.
	ArnFormat ArnFormat `json:"arnFormat"`
	// The partition that the resource is in.
	//
	// For standard AWS regions, the
	// partition is aws. If you have resources in other partitions, the
	// partition is aws-partitionname. For example, the partition for resources
	// in the China (Beijing) region is aws-cn.
	// Experimental.
	Partition *string `json:"partition"`
	// The region the resource resides in.
	//
	// Note that the ARNs for some resources
	// do not require a region, so this component might be omitted.
	// Experimental.
	Region *string `json:"region"`
	// Resource name or path within the resource (i.e. S3 bucket object key) or a wildcard such as ``"*"``. This is service-dependent.
	// Experimental.
	ResourceName *string `json:"resourceName"`
	// Separator between resource type and the resource.
	//
	// Can be either '/', ':' or an empty string. Will only be used if resourceName is defined.
	// Deprecated: use arnFormat instead
	Sep *string `json:"sep"`
}

// An enum representing the various ARN formats that different services use.
// Experimental.
type ArnFormat string

const (
	ArnFormat_NO_RESOURCE_NAME ArnFormat = "NO_RESOURCE_NAME"
	ArnFormat_COLON_RESOURCE_NAME ArnFormat = "COLON_RESOURCE_NAME"
	ArnFormat_SLASH_RESOURCE_NAME ArnFormat = "SLASH_RESOURCE_NAME"
	ArnFormat_SLASH_RESOURCE_SLASH_RESOURCE_NAME ArnFormat = "SLASH_RESOURCE_SLASH_RESOURCE_NAME"
)

// Aspects can be applied to CDK tree scopes and can operate on the tree before synthesis.
// Experimental.
type Aspects interface {
	Aspects() *[]IAspect
	Add(aspect IAspect)
}

// The jsii proxy struct for Aspects
type jsiiProxy_Aspects struct {
	_ byte // padding
}

func (j *jsiiProxy_Aspects) Aspects() *[]IAspect {
	var returns *[]IAspect
	_jsii_.Get(
		j,
		"aspects",
		&returns,
	)
	return returns
}


// Returns the `Aspects` object associated with a construct scope.
// Experimental.
func Aspects_Of(scope IConstruct) Aspects {
	_init_.Initialize()

	var returns Aspects

	_jsii_.StaticInvoke(
		"monocdk.Aspects",
		"of",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Adds an aspect to apply this scope before synthesis.
// Experimental.
func (a *jsiiProxy_Aspects) Add(aspect IAspect) {
	_jsii_.InvokeVoid(
		a,
		"add",
		[]interface{}{aspect},
	)
}

// The type of asset hash.
//
// NOTE: the hash is used in order to identify a specific revision of the asset, and
// used for optimizing and caching deployment activities related to this asset such as
// packaging, uploading to Amazon S3, etc.
// Experimental.
type AssetHashType string

const (
	AssetHashType_SOURCE AssetHashType = "SOURCE"
	AssetHashType_BUNDLE AssetHashType = "BUNDLE"
	AssetHashType_OUTPUT AssetHashType = "OUTPUT"
	AssetHashType_CUSTOM AssetHashType = "CUSTOM"
)

// Asset hash options.
// Experimental.
type AssetOptions struct {
	// Specify a custom hash for this asset.
	//
	// If `assetHashType` is set it must
	// be set to `AssetHashType.CUSTOM`. For consistency, this custom hash will
	// be SHA256 hashed and encoded as hex. The resulting hash will be the asset
	// hash.
	//
	// NOTE: the hash is used in order to identify a specific revision of the asset, and
	// used for optimizing and caching deployment activities related to this asset such as
	// packaging, uploading to Amazon S3, etc. If you chose to customize the hash, you will
	// need to make sure it is updated every time the asset changes, or otherwise it is
	// possible that some deployments will not be invalidated.
	// Experimental.
	AssetHash *string `json:"assetHash"`
	// Specifies the type of hash to calculate for this asset.
	//
	// If `assetHash` is configured, this option must be `undefined` or
	// `AssetHashType.CUSTOM`.
	// Experimental.
	AssetHashType AssetHashType `json:"assetHashType"`
	// Bundle the asset by executing a command in a Docker container or a custom bundling provider.
	//
	// The asset path will be mounted at `/asset-input`. The Docker
	// container is responsible for putting content at `/asset-output`.
	// The content at `/asset-output` will be zipped and used as the
	// final asset.
	// Experimental.
	Bundling *BundlingOptions `json:"bundling"`
}

// Stages a file or directory from a location on the file system into a staging directory.
//
// This is controlled by the context key 'aws:cdk:asset-staging' and enabled
// by the CLI by default in order to ensure that when the CDK app exists, all
// assets are available for deployment. Otherwise, if an app references assets
// in temporary locations, those will not be available when it exists (see
// https://github.com/aws/aws-cdk/issues/1716).
//
// The `stagedPath` property is a stringified token that represents the location
// of the file or directory after staging. It will be resolved only during the
// "prepare" stage and may be either the original path or the staged path
// depending on the context setting.
//
// The file/directory are staged based on their content hash (fingerprint). This
// means that only if content was changed, copy will happen.
// Experimental.
type AssetStaging interface {
	Construct
	AbsoluteStagedPath() *string
	AssetHash() *string
	IsArchive() *bool
	Node() ConstructNode
	Packaging() FileAssetPackaging
	SourceHash() *string
	SourcePath() *string
	StagedPath() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RelativeStagedPath(stack Stack) *string
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for AssetStaging
type jsiiProxy_AssetStaging struct {
	jsiiProxy_Construct
}

func (j *jsiiProxy_AssetStaging) AbsoluteStagedPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"absoluteStagedPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetStaging) AssetHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assetHash",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetStaging) IsArchive() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isArchive",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetStaging) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetStaging) Packaging() FileAssetPackaging {
	var returns FileAssetPackaging
	_jsii_.Get(
		j,
		"packaging",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetStaging) SourceHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceHash",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetStaging) SourcePath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourcePath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AssetStaging) StagedPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stagedPath",
		&returns,
	)
	return returns
}


// Experimental.
func NewAssetStaging(scope constructs.Construct, id *string, props *AssetStagingProps) AssetStaging {
	_init_.Initialize()

	j := jsiiProxy_AssetStaging{}

	_jsii_.Create(
		"monocdk.AssetStaging",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewAssetStaging_Override(a AssetStaging, scope constructs.Construct, id *string, props *AssetStagingProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.AssetStaging",
		[]interface{}{scope, id, props},
		a,
	)
}

// Clears the asset hash cache.
// Experimental.
func AssetStaging_ClearAssetHashCache() {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.AssetStaging",
		"clearAssetHashCache",
		nil, // no parameters
	)
}

// Return whether the given object is a Construct.
// Experimental.
func AssetStaging_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.AssetStaging",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func AssetStaging_BUNDLING_INPUT_DIR() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.AssetStaging",
		"BUNDLING_INPUT_DIR",
		&returns,
	)
	return returns
}

func AssetStaging_BUNDLING_OUTPUT_DIR() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.AssetStaging",
		"BUNDLING_OUTPUT_DIR",
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
func (a *jsiiProxy_AssetStaging) OnPrepare() {
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
func (a *jsiiProxy_AssetStaging) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_AssetStaging) OnValidate() *[]*string {
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
func (a *jsiiProxy_AssetStaging) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Return the path to the staged asset, relative to the Cloud Assembly (manifest) directory of the given stack.
//
// Only returns a relative path if the asset was staged, returns an absolute path if
// it was not staged.
//
// A bundled asset might end up in the outDir and still not count as
// "staged"; if asset staging is disabled we're technically expected to
// reference source directories, but we don't have a source directory for the
// bundled outputs (as the bundle output is written to a temporary
// directory). Nevertheless, we will still return an absolute path.
//
// A non-obvious directory layout may look like this:
//
// ```
//    CLOUD ASSEMBLY ROOT
//      +-- asset.12345abcdef/
//      +-- assembly-Stage
//            +-- MyStack.template.json
//            +-- MyStack.assets.json <- will contain { "path": "../asset.12345abcdef" }
// ```
// Experimental.
func (a *jsiiProxy_AssetStaging) RelativeStagedPath(stack Stack) *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"relativeStagedPath",
		[]interface{}{stack},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_AssetStaging) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_AssetStaging) ToString() *string {
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
func (a *jsiiProxy_AssetStaging) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Initialization properties for `AssetStaging`.
// Experimental.
type AssetStagingProps struct {
	// Glob patterns to exclude from the copy.
	// Experimental.
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Experimental.
	Follow SymlinkFollowMode `json:"follow"`
	// The ignore behavior to use for exclude patterns.
	// Experimental.
	IgnoreMode IgnoreMode `json:"ignoreMode"`
	// Extra information to encode into the fingerprint (e.g. build instructions and other inputs).
	// Experimental.
	ExtraHash *string `json:"extraHash"`
	// Specify a custom hash for this asset.
	//
	// If `assetHashType` is set it must
	// be set to `AssetHashType.CUSTOM`. For consistency, this custom hash will
	// be SHA256 hashed and encoded as hex. The resulting hash will be the asset
	// hash.
	//
	// NOTE: the hash is used in order to identify a specific revision of the asset, and
	// used for optimizing and caching deployment activities related to this asset such as
	// packaging, uploading to Amazon S3, etc. If you chose to customize the hash, you will
	// need to make sure it is updated every time the asset changes, or otherwise it is
	// possible that some deployments will not be invalidated.
	// Experimental.
	AssetHash *string `json:"assetHash"`
	// Specifies the type of hash to calculate for this asset.
	//
	// If `assetHash` is configured, this option must be `undefined` or
	// `AssetHashType.CUSTOM`.
	// Experimental.
	AssetHashType AssetHashType `json:"assetHashType"`
	// Bundle the asset by executing a command in a Docker container or a custom bundling provider.
	//
	// The asset path will be mounted at `/asset-input`. The Docker
	// container is responsible for putting content at `/asset-output`.
	// The content at `/asset-output` will be zipped and used as the
	// final asset.
	// Experimental.
	Bundling *BundlingOptions `json:"bundling"`
	// The source file or directory to copy from.
	// Experimental.
	SourcePath *string `json:"sourcePath"`
}

// Accessor for pseudo parameters.
//
// Since pseudo parameters need to be anchored to a stack somewhere in the
// construct tree, this class takes an scope parameter; the pseudo parameter
// values can be obtained as properties from an scoped object.
// Experimental.
type Aws interface {
}

// The jsii proxy struct for Aws
type jsiiProxy_Aws struct {
	_ byte // padding
}

func Aws_ACCOUNT_ID() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"ACCOUNT_ID",
		&returns,
	)
	return returns
}

func Aws_NO_VALUE() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"NO_VALUE",
		&returns,
	)
	return returns
}

func Aws_NOTIFICATION_ARNS() *[]*string {
	_init_.Initialize()
	var returns *[]*string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"NOTIFICATION_ARNS",
		&returns,
	)
	return returns
}

func Aws_PARTITION() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"PARTITION",
		&returns,
	)
	return returns
}

func Aws_REGION() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"REGION",
		&returns,
	)
	return returns
}

func Aws_STACK_ID() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"STACK_ID",
		&returns,
	)
	return returns
}

func Aws_STACK_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"STACK_NAME",
		&returns,
	)
	return returns
}

func Aws_URL_SUFFIX() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.Aws",
		"URL_SUFFIX",
		&returns,
	)
	return returns
}

// A special synthesizer that behaves similarly to DefaultStackSynthesizer, but doesn't require bootstrapping the environment it operates in.
//
// Because of that, stacks using it cannot have assets inside of them.
// Used by the CodePipeline construct for the support stacks needed for
// cross-region replication S3 buckets.
// Experimental.
type BootstraplessSynthesizer interface {
	DefaultStackSynthesizer
	CloudFormationExecutionRoleArn() *string
	DeployRoleArn() *string
	Stack() Stack
	AddDockerImageAsset(_asset *DockerImageAssetSource) *DockerImageAssetLocation
	AddFileAsset(_asset *FileAssetSource) *FileAssetLocation
	Bind(stack Stack)
	EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions)
	Synthesize(session ISynthesisSession)
	SynthesizeStackTemplate(stack Stack, session ISynthesisSession)
}

// The jsii proxy struct for BootstraplessSynthesizer
type jsiiProxy_BootstraplessSynthesizer struct {
	jsiiProxy_DefaultStackSynthesizer
}

func (j *jsiiProxy_BootstraplessSynthesizer) CloudFormationExecutionRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cloudFormationExecutionRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BootstraplessSynthesizer) DeployRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deployRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BootstraplessSynthesizer) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewBootstraplessSynthesizer(props *BootstraplessSynthesizerProps) BootstraplessSynthesizer {
	_init_.Initialize()

	j := jsiiProxy_BootstraplessSynthesizer{}

	_jsii_.Create(
		"monocdk.BootstraplessSynthesizer",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewBootstraplessSynthesizer_Override(b BootstraplessSynthesizer, props *BootstraplessSynthesizerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.BootstraplessSynthesizer",
		[]interface{}{props},
		b,
	)
}

func BootstraplessSynthesizer_DEFAULT_BOOTSTRAP_STACK_VERSION_SSM_PARAMETER() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_BOOTSTRAP_STACK_VERSION_SSM_PARAMETER",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_CLOUDFORMATION_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_CLOUDFORMATION_ROLE_ARN",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_DEPLOY_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_DEPLOY_ROLE_ARN",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_FILE_ASSET_KEY_ARN_EXPORT_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_FILE_ASSET_KEY_ARN_EXPORT_NAME",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_FILE_ASSET_PREFIX() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_FILE_ASSET_PREFIX",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_FILE_ASSET_PUBLISHING_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_FILE_ASSET_PUBLISHING_ROLE_ARN",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_FILE_ASSETS_BUCKET_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_FILE_ASSETS_BUCKET_NAME",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_IMAGE_ASSET_PUBLISHING_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_IMAGE_ASSET_PUBLISHING_ROLE_ARN",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_IMAGE_ASSETS_REPOSITORY_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_IMAGE_ASSETS_REPOSITORY_NAME",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_LOOKUP_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_LOOKUP_ROLE_ARN",
		&returns,
	)
	return returns
}

func BootstraplessSynthesizer_DEFAULT_QUALIFIER() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.BootstraplessSynthesizer",
		"DEFAULT_QUALIFIER",
		&returns,
	)
	return returns
}

// Register a Docker Image Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (b *jsiiProxy_BootstraplessSynthesizer) AddDockerImageAsset(_asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		b,
		"addDockerImageAsset",
		[]interface{}{_asset},
		&returns,
	)

	return returns
}

// Register a File Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (b *jsiiProxy_BootstraplessSynthesizer) AddFileAsset(_asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		b,
		"addFileAsset",
		[]interface{}{_asset},
		&returns,
	)

	return returns
}

// Bind to the stack this environment is going to be used on.
//
// Must be called before any of the other methods are called.
// Experimental.
func (b *jsiiProxy_BootstraplessSynthesizer) Bind(stack Stack) {
	_jsii_.InvokeVoid(
		b,
		"bind",
		[]interface{}{stack},
	)
}

// Write the stack artifact to the session.
//
// Use default settings to add a CloudFormationStackArtifact artifact to
// the given synthesis session.
// Experimental.
func (b *jsiiProxy_BootstraplessSynthesizer) EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions) {
	_jsii_.InvokeVoid(
		b,
		"emitStackArtifact",
		[]interface{}{stack, session, options},
	)
}

// Synthesize the associated stack to the session.
// Experimental.
func (b *jsiiProxy_BootstraplessSynthesizer) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Have the stack write out its template.
// Experimental.
func (b *jsiiProxy_BootstraplessSynthesizer) SynthesizeStackTemplate(stack Stack, session ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesizeStackTemplate",
		[]interface{}{stack, session},
	)
}

// Construction properties of {@link BootstraplessSynthesizer}.
// Experimental.
type BootstraplessSynthesizerProps struct {
	// The CFN execution Role ARN to use.
	// Experimental.
	CloudFormationExecutionRoleArn *string `json:"cloudFormationExecutionRoleArn"`
	// The deploy Role ARN to use.
	// Experimental.
	DeployRoleArn *string `json:"deployRoleArn"`
}

// A Docker image used for asset bundling.
// Deprecated: use DockerImage
type BundlingDockerImage interface {
	Image() *string
	Cp(imagePath *string, outputPath *string) *string
	Run(options *DockerRunOptions)
	ToJSON() *string
}

// The jsii proxy struct for BundlingDockerImage
type jsiiProxy_BundlingDockerImage struct {
	_ byte // padding
}

func (j *jsiiProxy_BundlingDockerImage) Image() *string {
	var returns *string
	_jsii_.Get(
		j,
		"image",
		&returns,
	)
	return returns
}


// Deprecated: use DockerImage
func NewBundlingDockerImage(image *string, _imageHash *string) BundlingDockerImage {
	_init_.Initialize()

	j := jsiiProxy_BundlingDockerImage{}

	_jsii_.Create(
		"monocdk.BundlingDockerImage",
		[]interface{}{image, _imageHash},
		&j,
	)

	return &j
}

// Deprecated: use DockerImage
func NewBundlingDockerImage_Override(b BundlingDockerImage, image *string, _imageHash *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.BundlingDockerImage",
		[]interface{}{image, _imageHash},
		b,
	)
}

// Reference an image that's built directly from sources on disk.
// Deprecated: use DockerImage.fromBuild()
func BundlingDockerImage_FromAsset(path *string, options *DockerBuildOptions) BundlingDockerImage {
	_init_.Initialize()

	var returns BundlingDockerImage

	_jsii_.StaticInvoke(
		"monocdk.BundlingDockerImage",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Reference an image on DockerHub or another online registry.
// Deprecated: use DockerImage
func BundlingDockerImage_FromRegistry(image *string) DockerImage {
	_init_.Initialize()

	var returns DockerImage

	_jsii_.StaticInvoke(
		"monocdk.BundlingDockerImage",
		"fromRegistry",
		[]interface{}{image},
		&returns,
	)

	return returns
}

// Copies a file or directory out of the Docker image to the local filesystem.
//
// If `outputPath` is omitted the destination path is a temporary directory.
//
// Returns: the destination path
// Deprecated: use DockerImage
func (b *jsiiProxy_BundlingDockerImage) Cp(imagePath *string, outputPath *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"cp",
		[]interface{}{imagePath, outputPath},
		&returns,
	)

	return returns
}

// Runs a Docker image.
// Deprecated: use DockerImage
func (b *jsiiProxy_BundlingDockerImage) Run(options *DockerRunOptions) {
	_jsii_.InvokeVoid(
		b,
		"run",
		[]interface{}{options},
	)
}

// Provides a stable representation of this image for JSON serialization.
//
// Returns: The overridden image name if set or image hash name in that order
// Deprecated: use DockerImage
func (b *jsiiProxy_BundlingDockerImage) ToJSON() *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Bundling options.
// Experimental.
type BundlingOptions struct {
	// The Docker image where the command will run.
	// Experimental.
	Image DockerImage `json:"image"`
	// The command to run in the Docker container.
	//
	// TODO: EXAMPLE
	//
	// See: https://docs.docker.com/engine/reference/run/
	//
	// Experimental.
	Command *[]*string `json:"command"`
	// The entrypoint to run in the Docker container.
	//
	// TODO: EXAMPLE
	//
	// See: https://docs.docker.com/engine/reference/builder/#entrypoint
	//
	// Experimental.
	Entrypoint *[]*string `json:"entrypoint"`
	// The environment variables to pass to the Docker container.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// Local bundling provider.
	//
	// The provider implements a method `tryBundle()` which should return `true`
	// if local bundling was performed. If `false` is returned, docker bundling
	// will be done.
	// Experimental.
	Local ILocalBundling `json:"local"`
	// The type of output that this bundling operation is producing.
	// Experimental.
	OutputType BundlingOutput `json:"outputType"`
	// The user to use when running the Docker container.
	//
	// user | user:group | uid | uid:gid | user:gid | uid:group
	// See: https://docs.docker.com/engine/reference/run/#user
	//
	// Experimental.
	User *string `json:"user"`
	// Additional Docker volumes to mount.
	// Experimental.
	Volumes *[]*DockerVolume `json:"volumes"`
	// Working directory inside the Docker container.
	// Experimental.
	WorkingDirectory *string `json:"workingDirectory"`
}

// The type of output that a bundling operation is producing.
// Experimental.
type BundlingOutput string

const (
	BundlingOutput_ARCHIVED BundlingOutput = "ARCHIVED"
	BundlingOutput_NOT_ARCHIVED BundlingOutput = "NOT_ARCHIVED"
	BundlingOutput_AUTO_DISCOVER BundlingOutput = "AUTO_DISCOVER"
)

// Specifies whether an Auto Scaling group and the instances it contains are replaced during an update.
//
// During replacement,
// AWS CloudFormation retains the old group until it finishes creating the new one. If the update fails, AWS CloudFormation
// can roll back to the old Auto Scaling group and delete the new Auto Scaling group.
//
// While AWS CloudFormation creates the new group, it doesn't detach or attach any instances. After successfully creating
// the new Auto Scaling group, AWS CloudFormation deletes the old Auto Scaling group during the cleanup process.
//
// When you set the WillReplace parameter, remember to specify a matching CreationPolicy. If the minimum number of
// instances (specified by the MinSuccessfulInstancesPercent property) don't signal success within the Timeout period
// (specified in the CreationPolicy policy), the replacement update fails and AWS CloudFormation rolls back to the old
// Auto Scaling group.
// Experimental.
type CfnAutoScalingReplacingUpdate struct {
	// Experimental.
	WillReplace *bool `json:"willReplace"`
}

// To specify how AWS CloudFormation handles rolling updates for an Auto Scaling group, use the AutoScalingRollingUpdate policy.
//
// Rolling updates enable you to specify whether AWS CloudFormation updates instances that are in an Auto Scaling
// group in batches or all at once.
// Experimental.
type CfnAutoScalingRollingUpdate struct {
	// Specifies the maximum number of instances that AWS CloudFormation updates.
	// Experimental.
	MaxBatchSize *float64 `json:"maxBatchSize"`
	// Specifies the minimum number of instances that must be in service within the Auto Scaling group while AWS CloudFormation updates old instances.
	// Experimental.
	MinInstancesInService *float64 `json:"minInstancesInService"`
	// Specifies the percentage of instances in an Auto Scaling rolling update that must signal success for an update to succeed.
	//
	// You can specify a value from 0 to 100. AWS CloudFormation rounds to the nearest tenth of a percent. For example, if you
	// update five instances with a minimum successful percentage of 50, three instances must signal success.
	//
	// If an instance doesn't send a signal within the time specified in the PauseTime property, AWS CloudFormation assumes
	// that the instance wasn't updated.
	//
	// If you specify this property, you must also enable the WaitOnResourceSignals and PauseTime properties.
	// Experimental.
	MinSuccessfulInstancesPercent *float64 `json:"minSuccessfulInstancesPercent"`
	// The amount of time that AWS CloudFormation pauses after making a change to a batch of instances to give those instances time to start software applications.
	//
	// For example, you might need to specify PauseTime when scaling up the number of
	// instances in an Auto Scaling group.
	//
	// If you enable the WaitOnResourceSignals property, PauseTime is the amount of time that AWS CloudFormation should wait
	// for the Auto Scaling group to receive the required number of valid signals from added or replaced instances. If the
	// PauseTime is exceeded before the Auto Scaling group receives the required number of signals, the update fails. For best
	// results, specify a time period that gives your applications sufficient time to get started. If the update needs to be
	// rolled back, a short PauseTime can cause the rollback to fail.
	//
	// Specify PauseTime in the ISO8601 duration format (in the format PT#H#M#S, where each # is the number of hours, minutes,
	// and seconds, respectively). The maximum PauseTime is one hour (PT1H).
	// Experimental.
	PauseTime *string `json:"pauseTime"`
	// Specifies the Auto Scaling processes to suspend during a stack update.
	//
	// Suspending processes prevents Auto Scaling from
	// interfering with a stack update. For example, you can suspend alarming so that Auto Scaling doesn't execute scaling
	// policies associated with an alarm. For valid values, see the ScalingProcesses.member.N parameter for the SuspendProcesses
	// action in the Auto Scaling API Reference.
	// Experimental.
	SuspendProcesses *[]*string `json:"suspendProcesses"`
	// Specifies whether the Auto Scaling group waits on signals from new instances during an update.
	//
	// Use this property to
	// ensure that instances have completed installing and configuring applications before the Auto Scaling group update proceeds.
	// AWS CloudFormation suspends the update of an Auto Scaling group after new EC2 instances are launched into the group.
	// AWS CloudFormation must receive a signal from each new instance within the specified PauseTime before continuing the update.
	// To signal the Auto Scaling group, use the cfn-signal helper script or SignalResource API.
	//
	// To have instances wait for an Elastic Load Balancing health check before they signal success, add a health-check
	// verification by using the cfn-init helper script. For an example, see the verify_instance_health command in the Auto Scaling
	// rolling updates sample template.
	// Experimental.
	WaitOnResourceSignals *bool `json:"waitOnResourceSignals"`
}

// With scheduled actions, the group size properties of an Auto Scaling group can change at any time.
//
// When you update a
// stack with an Auto Scaling group and scheduled action, AWS CloudFormation always sets the group size property values of
// your Auto Scaling group to the values that are defined in the AWS::AutoScaling::AutoScalingGroup resource of your template,
// even if a scheduled action is in effect.
//
// If you do not want AWS CloudFormation to change any of the group size property values when you have a scheduled action in
// effect, use the AutoScalingScheduledAction update policy to prevent AWS CloudFormation from changing the MinSize, MaxSize,
// or DesiredCapacity properties unless you have modified these values in your template.\
// Experimental.
type CfnAutoScalingScheduledAction struct {
	// Experimental.
	IgnoreUnmodifiedGroupSizeProperties *bool `json:"ignoreUnmodifiedGroupSizeProperties"`
}

// Capabilities that affect whether CloudFormation is allowed to change IAM resources.
// Experimental.
type CfnCapabilities string

const (
	CfnCapabilities_NONE CfnCapabilities = "NONE"
	CfnCapabilities_ANONYMOUS_IAM CfnCapabilities = "ANONYMOUS_IAM"
	CfnCapabilities_NAMED_IAM CfnCapabilities = "NAMED_IAM"
	CfnCapabilities_AUTO_EXPAND CfnCapabilities = "AUTO_EXPAND"
)

// Additional options for the blue/green deployment.
//
// The type of the {@link CfnCodeDeployBlueGreenHookProps.additionalOptions} property.
// Experimental.
type CfnCodeDeployBlueGreenAdditionalOptions struct {
	// Specifies time to wait, in minutes, before terminating the blue resources.
	// Experimental.
	TerminationWaitTimeInMinutes *float64 `json:"terminationWaitTimeInMinutes"`
}

// The application actually being deployed.
//
// Type of the {@link CfnCodeDeployBlueGreenHookProps.applications} property.
// Experimental.
type CfnCodeDeployBlueGreenApplication struct {
	// The detailed attributes of the deployed target.
	// Experimental.
	EcsAttributes *CfnCodeDeployBlueGreenEcsAttributes `json:"ecsAttributes"`
	// The target that is being deployed.
	// Experimental.
	Target *CfnCodeDeployBlueGreenApplicationTarget `json:"target"`
}

// Type of the {@link CfnCodeDeployBlueGreenApplication.target} property.
// Experimental.
type CfnCodeDeployBlueGreenApplicationTarget struct {
	// The logical id of the target resource.
	// Experimental.
	LogicalId *string `json:"logicalId"`
	// The resource type of the target being deployed.
	//
	// Right now, the only allowed value is 'AWS::ECS::Service'.
	// Experimental.
	Type *string `json:"type"`
}

// The attributes of the ECS Service being deployed.
//
// Type of the {@link CfnCodeDeployBlueGreenApplication.ecsAttributes} property.
// Experimental.
type CfnCodeDeployBlueGreenEcsAttributes struct {
	// The logical IDs of the blue and green, respectively, AWS::ECS::TaskDefinition task definitions.
	// Experimental.
	TaskDefinitions *[]*string `json:"taskDefinitions"`
	// The logical IDs of the blue and green, respectively, AWS::ECS::TaskSet task sets.
	// Experimental.
	TaskSets *[]*string `json:"taskSets"`
	// The traffic routing configuration.
	// Experimental.
	TrafficRouting *CfnTrafficRouting `json:"trafficRouting"`
}

// A CloudFormation Hook for CodeDeploy blue-green ECS deployments.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/blue-green.html#blue-green-template-reference
//
// Experimental.
type CfnCodeDeployBlueGreenHook interface {
	CfnHook
	AdditionalOptions() *CfnCodeDeployBlueGreenAdditionalOptions
	SetAdditionalOptions(val *CfnCodeDeployBlueGreenAdditionalOptions)
	Applications() *[]*CfnCodeDeployBlueGreenApplication
	SetApplications(val *[]*CfnCodeDeployBlueGreenApplication)
	CreationStack() *[]*string
	LifecycleEventHooks() *CfnCodeDeployBlueGreenLifecycleEventHooks
	SetLifecycleEventHooks(val *CfnCodeDeployBlueGreenLifecycleEventHooks)
	LogicalId() *string
	Node() ConstructNode
	ServiceRole() *string
	SetServiceRole(val *string)
	Stack() Stack
	TrafficRoutingConfig() *CfnTrafficRoutingConfig
	SetTrafficRoutingConfig(val *CfnTrafficRoutingConfig)
	Type() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(_props *map[string]interface{}) *map[string]interface{}
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnCodeDeployBlueGreenHook
type jsiiProxy_CfnCodeDeployBlueGreenHook struct {
	jsiiProxy_CfnHook
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) AdditionalOptions() *CfnCodeDeployBlueGreenAdditionalOptions {
	var returns *CfnCodeDeployBlueGreenAdditionalOptions
	_jsii_.Get(
		j,
		"additionalOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) Applications() *[]*CfnCodeDeployBlueGreenApplication {
	var returns *[]*CfnCodeDeployBlueGreenApplication
	_jsii_.Get(
		j,
		"applications",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) LifecycleEventHooks() *CfnCodeDeployBlueGreenLifecycleEventHooks {
	var returns *CfnCodeDeployBlueGreenLifecycleEventHooks
	_jsii_.Get(
		j,
		"lifecycleEventHooks",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) ServiceRole() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceRole",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) TrafficRoutingConfig() *CfnTrafficRoutingConfig {
	var returns *CfnTrafficRoutingConfig
	_jsii_.Get(
		j,
		"trafficRoutingConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Creates a new CodeDeploy blue-green ECS Hook.
// Experimental.
func NewCfnCodeDeployBlueGreenHook(scope constructs.Construct, id *string, props *CfnCodeDeployBlueGreenHookProps) CfnCodeDeployBlueGreenHook {
	_init_.Initialize()

	j := jsiiProxy_CfnCodeDeployBlueGreenHook{}

	_jsii_.Create(
		"monocdk.CfnCodeDeployBlueGreenHook",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates a new CodeDeploy blue-green ECS Hook.
// Experimental.
func NewCfnCodeDeployBlueGreenHook_Override(c CfnCodeDeployBlueGreenHook, scope constructs.Construct, id *string, props *CfnCodeDeployBlueGreenHookProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnCodeDeployBlueGreenHook",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) SetAdditionalOptions(val *CfnCodeDeployBlueGreenAdditionalOptions) {
	_jsii_.Set(
		j,
		"additionalOptions",
		val,
	)
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) SetApplications(val *[]*CfnCodeDeployBlueGreenApplication) {
	_jsii_.Set(
		j,
		"applications",
		val,
	)
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) SetLifecycleEventHooks(val *CfnCodeDeployBlueGreenLifecycleEventHooks) {
	_jsii_.Set(
		j,
		"lifecycleEventHooks",
		val,
	)
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) SetServiceRole(val *string) {
	_jsii_.Set(
		j,
		"serviceRole",
		val,
	)
}

func (j *jsiiProxy_CfnCodeDeployBlueGreenHook) SetTrafficRoutingConfig(val *CfnTrafficRoutingConfig) {
	_jsii_.Set(
		j,
		"trafficRoutingConfig",
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
func CfnCodeDeployBlueGreenHook_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnCodeDeployBlueGreenHook",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCodeDeployBlueGreenHook_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnCodeDeployBlueGreenHook",
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
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) OnPrepare() {
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
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Experimental.
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) RenderProperties(_props *map[string]interface{}) *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"renderProperties",
		[]interface{}{_props},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) ToString() *string {
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
func (c *jsiiProxy_CfnCodeDeployBlueGreenHook) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties of {@link CfnCodeDeployBlueGreenHook}.
// Experimental.
type CfnCodeDeployBlueGreenHookProps struct {
	// Properties of the Amazon ECS applications being deployed.
	// Experimental.
	Applications *[]*CfnCodeDeployBlueGreenApplication `json:"applications"`
	// The IAM Role for CloudFormation to use to perform blue-green deployments.
	// Experimental.
	ServiceRole *string `json:"serviceRole"`
	// Additional options for the blue/green deployment.
	// Experimental.
	AdditionalOptions *CfnCodeDeployBlueGreenAdditionalOptions `json:"additionalOptions"`
	// Use lifecycle event hooks to specify a Lambda function that CodeDeploy can call to validate a deployment.
	//
	// You can use the same function or a different one for deployment lifecycle events.
	// Following completion of the validation tests,
	// the Lambda {@link CfnCodeDeployBlueGreenLifecycleEventHooks.afterAllowTraffic}
	// function calls back CodeDeploy and delivers a result of 'Succeeded' or 'Failed'.
	// Experimental.
	LifecycleEventHooks *CfnCodeDeployBlueGreenLifecycleEventHooks `json:"lifecycleEventHooks"`
	// Traffic routing configuration settings.
	// Experimental.
	TrafficRoutingConfig *CfnTrafficRoutingConfig `json:"trafficRoutingConfig"`
}

// Lifecycle events for blue-green deployments.
//
// The type of the {@link CfnCodeDeployBlueGreenHookProps.lifecycleEventHooks} property.
// Experimental.
type CfnCodeDeployBlueGreenLifecycleEventHooks struct {
	// Function to use to run tasks after the test listener serves traffic to the replacement task set.
	// Experimental.
	AfterAllowTestTraffic *string `json:"afterAllowTestTraffic"`
	// Function to use to run tasks after the second target group serves traffic to the replacement task set.
	// Experimental.
	AfterAllowTraffic *string `json:"afterAllowTraffic"`
	// Function to use to run tasks after the replacement task set is created and one of the target groups is associated with it.
	// Experimental.
	AfterInstall *string `json:"afterInstall"`
	// Function to use to run tasks after the second target group is associated with the replacement task set, but before traffic is shifted to the replacement task set.
	// Experimental.
	BeforeAllowTraffic *string `json:"beforeAllowTraffic"`
	// Function to use to run tasks before the replacement task set is created.
	// Experimental.
	BeforeInstall *string `json:"beforeInstall"`
}

// To perform an AWS CodeDeploy deployment when the version changes on an AWS::Lambda::Alias resource, use the CodeDeployLambdaAliasUpdate update policy.
// Experimental.
type CfnCodeDeployLambdaAliasUpdate struct {
	// The name of the AWS CodeDeploy application.
	// Experimental.
	ApplicationName *string `json:"applicationName"`
	// The name of the AWS CodeDeploy deployment group.
	//
	// This is where the traffic-shifting policy is set.
	// Experimental.
	DeploymentGroupName *string `json:"deploymentGroupName"`
	// The name of the Lambda function to run after traffic routing completes.
	// Experimental.
	AfterAllowTrafficHook *string `json:"afterAllowTrafficHook"`
	// The name of the Lambda function to run before traffic routing starts.
	// Experimental.
	BeforeAllowTrafficHook *string `json:"beforeAllowTrafficHook"`
}

// Represents a CloudFormation condition, for resources which must be conditionally created and the determination must be made at deploy time.
// Experimental.
type CfnCondition interface {
	CfnElement
	ICfnConditionExpression
	IResolvable
	CreationStack() *[]*string
	Expression() ICfnConditionExpression
	SetExpression(val ICfnConditionExpression)
	LogicalId() *string
	Node() ConstructNode
	Stack() Stack
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	Resolve(_context IResolveContext) interface{}
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnCondition
type jsiiProxy_CfnCondition struct {
	jsiiProxy_CfnElement
	jsiiProxy_ICfnConditionExpression
	jsiiProxy_IResolvable
}

func (j *jsiiProxy_CfnCondition) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCondition) Expression() ICfnConditionExpression {
	var returns ICfnConditionExpression
	_jsii_.Get(
		j,
		"expression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCondition) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCondition) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCondition) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Build a new condition.
//
// The condition must be constructed with a condition token,
// that the condition is based on.
// Experimental.
func NewCfnCondition(scope constructs.Construct, id *string, props *CfnConditionProps) CfnCondition {
	_init_.Initialize()

	j := jsiiProxy_CfnCondition{}

	_jsii_.Create(
		"monocdk.CfnCondition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Build a new condition.
//
// The condition must be constructed with a condition token,
// that the condition is based on.
// Experimental.
func NewCfnCondition_Override(c CfnCondition, scope constructs.Construct, id *string, props *CfnConditionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnCondition",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCondition) SetExpression(val ICfnConditionExpression) {
	_jsii_.Set(
		j,
		"expression",
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
func CfnCondition_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnCondition",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCondition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnCondition",
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
func (c *jsiiProxy_CfnCondition) OnPrepare() {
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
func (c *jsiiProxy_CfnCondition) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCondition) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCondition) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCondition) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Synthesizes the condition.
// Experimental.
func (c *jsiiProxy_CfnCondition) Resolve(_context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnCondition) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnCondition) ToString() *string {
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
func (c *jsiiProxy_CfnCondition) Validate() *[]*string {
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
type CfnConditionProps struct {
	// The expression that the condition will evaluate.
	// Experimental.
	Expression ICfnConditionExpression `json:"expression"`
}

// Associate the CreationPolicy attribute with a resource to prevent its status from reaching create complete until AWS CloudFormation receives a specified number of success signals or the timeout period is exceeded.
//
// To signal a
// resource, you can use the cfn-signal helper script or SignalResource API. AWS CloudFormation publishes valid signals
// to the stack events so that you track the number of signals sent.
//
// The creation policy is invoked only when AWS CloudFormation creates the associated resource. Currently, the only
// AWS CloudFormation resources that support creation policies are AWS::AutoScaling::AutoScalingGroup, AWS::EC2::Instance,
// and AWS::CloudFormation::WaitCondition.
//
// Use the CreationPolicy attribute when you want to wait on resource configuration actions before stack creation proceeds.
// For example, if you install and configure software applications on an EC2 instance, you might want those applications to
// be running before proceeding. In such cases, you can add a CreationPolicy attribute to the instance, and then send a success
// signal to the instance after the applications are installed and configured. For a detailed example, see Deploying Applications
// on Amazon EC2 with AWS CloudFormation.
// Experimental.
type CfnCreationPolicy struct {
	// For an Auto Scaling group replacement update, specifies how many instances must signal success for the update to succeed.
	// Experimental.
	AutoScalingCreationPolicy *CfnResourceAutoScalingCreationPolicy `json:"autoScalingCreationPolicy"`
	// When AWS CloudFormation creates the associated resource, configures the number of required success signals and the length of time that AWS CloudFormation waits for those signals.
	// Experimental.
	ResourceSignal *CfnResourceSignal `json:"resourceSignal"`
}

// A CloudFormation `AWS::CloudFormation::CustomResource`.
type CfnCustomResource interface {
	CfnResource
	IInspectable
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	ServiceToken() *string
	SetServiceToken(val *string)
	Stack() Stack
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnCustomResource
type jsiiProxy_CfnCustomResource struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnCustomResource) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) ServiceToken() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceToken",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomResource) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::CustomResource`.
func NewCfnCustomResource(scope Construct, id *string, props *CfnCustomResourceProps) CfnCustomResource {
	_init_.Initialize()

	j := jsiiProxy_CfnCustomResource{}

	_jsii_.Create(
		"monocdk.CfnCustomResource",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::CustomResource`.
func NewCfnCustomResource_Override(c CfnCustomResource, scope Construct, id *string, props *CfnCustomResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnCustomResource",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCustomResource) SetServiceToken(val *string) {
	_jsii_.Set(
		j,
		"serviceToken",
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
func CfnCustomResource_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnCustomResource",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnCustomResource_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnCustomResource",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCustomResource_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnCustomResource",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnCustomResource_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnCustomResource",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnCustomResource) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnCustomResource) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnCustomResource) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnCustomResource) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnCustomResource) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnCustomResource) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnCustomResource) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnCustomResource) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnCustomResource) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnCustomResource) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnCustomResource) OnPrepare() {
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
func (c *jsiiProxy_CfnCustomResource) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCustomResource) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCustomResource) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCustomResource) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnCustomResource) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnCustomResource) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnCustomResource) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnCustomResource) ToString() *string {
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
func (c *jsiiProxy_CfnCustomResource) Validate() *[]*string {
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
func (c *jsiiProxy_CfnCustomResource) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudFormation::CustomResource`.
type CfnCustomResourceProps struct {
	// `AWS::CloudFormation::CustomResource.ServiceToken`.
	ServiceToken *string `json:"serviceToken"`
}

// With the DeletionPolicy attribute you can preserve or (in some cases) backup a resource when its stack is deleted.
//
// You specify a DeletionPolicy attribute for each resource that you want to control. If a resource has no DeletionPolicy
// attribute, AWS CloudFormation deletes the resource by default. Note that this capability also applies to update operations
// that lead to resources being removed.
// Experimental.
type CfnDeletionPolicy string

const (
	CfnDeletionPolicy_DELETE CfnDeletionPolicy = "DELETE"
	CfnDeletionPolicy_RETAIN CfnDeletionPolicy = "RETAIN"
	CfnDeletionPolicy_SNAPSHOT CfnDeletionPolicy = "SNAPSHOT"
)

// References a dynamically retrieved value.
//
// This is a Construct so that subclasses will (eventually) be able to attach
// metadata to themselves without having to change call signatures.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/dynamic-references.html
//
// Experimental.
type CfnDynamicReference interface {
	Intrinsic
	CreationStack() *[]*string
	NewError(message *string) interface{}
	Resolve(_context IResolveContext) interface{}
	ToJSON() interface{}
	ToString() *string
}

// The jsii proxy struct for CfnDynamicReference
type jsiiProxy_CfnDynamicReference struct {
	jsiiProxy_Intrinsic
}

func (j *jsiiProxy_CfnDynamicReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}


// Experimental.
func NewCfnDynamicReference(service CfnDynamicReferenceService, key *string) CfnDynamicReference {
	_init_.Initialize()

	j := jsiiProxy_CfnDynamicReference{}

	_jsii_.Create(
		"monocdk.CfnDynamicReference",
		[]interface{}{service, key},
		&j,
	)

	return &j
}

// Experimental.
func NewCfnDynamicReference_Override(c CfnDynamicReference, service CfnDynamicReferenceService, key *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnDynamicReference",
		[]interface{}{service, key},
		c,
	)
}

// Creates a throwable Error object that contains the token creation stack trace.
// Experimental.
func (c *jsiiProxy_CfnDynamicReference) NewError(message *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"newError",
		[]interface{}{message},
		&returns,
	)

	return returns
}

// Produce the Token's value at resolution time.
// Experimental.
func (c *jsiiProxy_CfnDynamicReference) Resolve(_context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

// Turn this Token into JSON.
//
// Called automatically when JSON.stringify() is called on a Token.
// Experimental.
func (c *jsiiProxy_CfnDynamicReference) ToJSON() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Convert an instance of this Token to a string.
//
// This method will be called implicitly by language runtimes if the object
// is embedded into a string. We treat it the same as an explicit
// stringification.
// Experimental.
func (c *jsiiProxy_CfnDynamicReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a Dynamic Reference.
// Experimental.
type CfnDynamicReferenceProps struct {
	// The reference key of the dynamic reference.
	// Experimental.
	ReferenceKey *string `json:"referenceKey"`
	// The service to retrieve the dynamic reference from.
	// Experimental.
	Service CfnDynamicReferenceService `json:"service"`
}

// The service to retrieve the dynamic reference from.
// Experimental.
type CfnDynamicReferenceService string

const (
	CfnDynamicReferenceService_SSM CfnDynamicReferenceService = "SSM"
	CfnDynamicReferenceService_SSM_SECURE CfnDynamicReferenceService = "SSM_SECURE"
	CfnDynamicReferenceService_SECRETS_MANAGER CfnDynamicReferenceService = "SECRETS_MANAGER"
)

// An element of a CloudFormation stack.
// Experimental.
type CfnElement interface {
	Construct
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Stack() Stack
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnElement
type jsiiProxy_CfnElement struct {
	jsiiProxy_Construct
}

func (j *jsiiProxy_CfnElement) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnElement) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnElement) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnElement) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Creates an entity and binds it to a tree.
//
// Note that the root of the tree must be a Stack object (not just any Root).
// Experimental.
func NewCfnElement_Override(c CfnElement, scope constructs.Construct, id *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnElement",
		[]interface{}{scope, id},
		c,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnElement_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnElement",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnElement_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnElement",
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
func (c *jsiiProxy_CfnElement) OnPrepare() {
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
func (c *jsiiProxy_CfnElement) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnElement) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnElement) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnElement) Prepare() {
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
func (c *jsiiProxy_CfnElement) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnElement) ToString() *string {
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
func (c *jsiiProxy_CfnElement) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents a CloudFormation resource.
// Experimental.
type CfnHook interface {
	CfnElement
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Stack() Stack
	Type() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnHook
type jsiiProxy_CfnHook struct {
	jsiiProxy_CfnElement
}

func (j *jsiiProxy_CfnHook) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnHook) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnHook) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnHook) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnHook) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}


// Creates a new Hook object.
// Experimental.
func NewCfnHook(scope constructs.Construct, id *string, props *CfnHookProps) CfnHook {
	_init_.Initialize()

	j := jsiiProxy_CfnHook{}

	_jsii_.Create(
		"monocdk.CfnHook",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates a new Hook object.
// Experimental.
func NewCfnHook_Override(c CfnHook, scope constructs.Construct, id *string, props *CfnHookProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnHook",
		[]interface{}{scope, id, props},
		c,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnHook_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnHook",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnHook_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnHook",
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
func (c *jsiiProxy_CfnHook) OnPrepare() {
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
func (c *jsiiProxy_CfnHook) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnHook) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnHook) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnHook) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Experimental.
func (c *jsiiProxy_CfnHook) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		c,
		"renderProperties",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnHook) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnHook) ToString() *string {
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
func (c *jsiiProxy_CfnHook) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties of {@link CfnHook}.
// Experimental.
type CfnHookProps struct {
	// The type of the hook (for example, "AWS::CodeDeploy::BlueGreen").
	// Experimental.
	Type *string `json:"type"`
	// The properties of the hook.
	// Experimental.
	Properties *map[string]interface{} `json:"properties"`
}

// Includes a CloudFormation template into a stack.
//
// All elements of the template will be merged into
// the current stack, together with any elements created programmatically.
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
type CfnInclude interface {
	CfnElement
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Stack() Stack
	Template() *map[string]interface{}
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnInclude
type jsiiProxy_CfnInclude struct {
	jsiiProxy_CfnElement
}

func (j *jsiiProxy_CfnInclude) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInclude) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInclude) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInclude) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInclude) Template() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"template",
		&returns,
	)
	return returns
}


// Creates an adopted template construct.
//
// The template will be incorporated into the stack as-is with no changes at all.
// This means that logical IDs of entities within this template may conflict with logical IDs of entities that are part of the
// stack.
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func NewCfnInclude(scope constructs.Construct, id *string, props *CfnIncludeProps) CfnInclude {
	_init_.Initialize()

	j := jsiiProxy_CfnInclude{}

	_jsii_.Create(
		"monocdk.CfnInclude",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates an adopted template construct.
//
// The template will be incorporated into the stack as-is with no changes at all.
// This means that logical IDs of entities within this template may conflict with logical IDs of entities that are part of the
// stack.
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func NewCfnInclude_Override(c CfnInclude, scope constructs.Construct, id *string, props *CfnIncludeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnInclude",
		[]interface{}{scope, id, props},
		c,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func CfnInclude_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnInclude",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func CfnInclude_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnInclude",
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
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) OnPrepare() {
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
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) OnValidate() *[]*string {
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
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) Prepare() {
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
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) ToString() *string {
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
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
func (c *jsiiProxy_CfnInclude) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Construction properties for {@link CfnInclude}.
// Deprecated: use the CfnInclude class from the cloudformation-include module instead
type CfnIncludeProps struct {
	// The CloudFormation template to include in the stack (as is).
	// Deprecated: use the CfnInclude class from the cloudformation-include module instead
	Template *map[string]interface{} `json:"template"`
}

// Captures a synthesis-time JSON object a CloudFormation reference which resolves during deployment to the resolved values of the JSON object.
//
// The main use case for this is to overcome a limitation in CloudFormation that
// does not allow using intrinsic functions as dictionary keys (because
// dictionary keys in JSON must be strings). Specifically this is common in IAM
// conditions such as `StringEquals: { lhs: "rhs" }` where you want "lhs" to be
// a reference.
//
// This object is resolvable, so it can be used as a value.
//
// This construct is backed by a custom resource.
// Experimental.
type CfnJson interface {
	Construct
	IResolvable
	CreationStack() *[]*string
	Node() ConstructNode
	Value() Reference
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Resolve(_arg IResolveContext) interface{}
	Synthesize(session ISynthesisSession)
	ToJSON() *string
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnJson
type jsiiProxy_CfnJson struct {
	jsiiProxy_Construct
	jsiiProxy_IResolvable
}

func (j *jsiiProxy_CfnJson) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnJson) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnJson) Value() Reference {
	var returns Reference
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


// Experimental.
func NewCfnJson(scope constructs.Construct, id *string, props *CfnJsonProps) CfnJson {
	_init_.Initialize()

	j := jsiiProxy_CfnJson{}

	_jsii_.Create(
		"monocdk.CfnJson",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCfnJson_Override(c CfnJson, scope constructs.Construct, id *string, props *CfnJsonProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnJson",
		[]interface{}{scope, id, props},
		c,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func CfnJson_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnJson",
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
func (c *jsiiProxy_CfnJson) OnPrepare() {
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
func (c *jsiiProxy_CfnJson) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnJson) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnJson) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Produce the Token's value at resolution time.
// Experimental.
func (c *jsiiProxy_CfnJson) Resolve(_arg IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"resolve",
		[]interface{}{_arg},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnJson) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// This is required in case someone JSON.stringifys an object which refrences this object. Otherwise, we'll get a cyclic JSON reference.
// Experimental.
func (c *jsiiProxy_CfnJson) ToJSON() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnJson) ToString() *string {
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
func (c *jsiiProxy_CfnJson) Validate() *[]*string {
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
type CfnJsonProps struct {
	// The value to resolve.
	//
	// Can be any JavaScript object, including tokens and
	// references in keys or values.
	// Experimental.
	Value interface{} `json:"value"`
}

// A CloudFormation `AWS::CloudFormation::Macro`.
type CfnMacro interface {
	CfnResource
	IInspectable
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	FunctionName() *string
	SetFunctionName(val *string)
	LogGroupName() *string
	SetLogGroupName(val *string)
	LogicalId() *string
	LogRoleArn() *string
	SetLogRoleArn(val *string)
	Name() *string
	SetName(val *string)
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnMacro
type jsiiProxy_CfnMacro struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnMacro) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) FunctionName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"functionName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) LogGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) LogRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMacro) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::Macro`.
func NewCfnMacro(scope Construct, id *string, props *CfnMacroProps) CfnMacro {
	_init_.Initialize()

	j := jsiiProxy_CfnMacro{}

	_jsii_.Create(
		"monocdk.CfnMacro",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::Macro`.
func NewCfnMacro_Override(c CfnMacro, scope Construct, id *string, props *CfnMacroProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnMacro",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnMacro) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnMacro) SetFunctionName(val *string) {
	_jsii_.Set(
		j,
		"functionName",
		val,
	)
}

func (j *jsiiProxy_CfnMacro) SetLogGroupName(val *string) {
	_jsii_.Set(
		j,
		"logGroupName",
		val,
	)
}

func (j *jsiiProxy_CfnMacro) SetLogRoleArn(val *string) {
	_jsii_.Set(
		j,
		"logRoleArn",
		val,
	)
}

func (j *jsiiProxy_CfnMacro) SetName(val *string) {
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
func CfnMacro_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnMacro",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnMacro_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnMacro",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnMacro_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnMacro",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnMacro_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnMacro",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnMacro) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnMacro) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnMacro) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnMacro) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnMacro) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnMacro) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnMacro) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnMacro) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnMacro) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnMacro) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnMacro) OnPrepare() {
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
func (c *jsiiProxy_CfnMacro) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnMacro) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnMacro) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnMacro) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnMacro) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnMacro) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnMacro) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnMacro) ToString() *string {
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
func (c *jsiiProxy_CfnMacro) Validate() *[]*string {
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
func (c *jsiiProxy_CfnMacro) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudFormation::Macro`.
type CfnMacroProps struct {
	// `AWS::CloudFormation::Macro.FunctionName`.
	FunctionName *string `json:"functionName"`
	// `AWS::CloudFormation::Macro.Name`.
	Name *string `json:"name"`
	// `AWS::CloudFormation::Macro.Description`.
	Description *string `json:"description"`
	// `AWS::CloudFormation::Macro.LogGroupName`.
	LogGroupName *string `json:"logGroupName"`
	// `AWS::CloudFormation::Macro.LogRoleARN`.
	LogRoleArn *string `json:"logRoleArn"`
}

// Represents a CloudFormation mapping.
// Experimental.
type CfnMapping interface {
	CfnRefElement
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	FindInMap(key1 *string, key2 *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	SetValue(key1 *string, key2 *string, value interface{})
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnMapping
type jsiiProxy_CfnMapping struct {
	jsiiProxy_CfnRefElement
}

func (j *jsiiProxy_CfnMapping) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMapping) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMapping) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMapping) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMapping) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewCfnMapping(scope constructs.Construct, id *string, props *CfnMappingProps) CfnMapping {
	_init_.Initialize()

	j := jsiiProxy_CfnMapping{}

	_jsii_.Create(
		"monocdk.CfnMapping",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCfnMapping_Override(c CfnMapping, scope constructs.Construct, id *string, props *CfnMappingProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnMapping",
		[]interface{}{scope, id, props},
		c,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnMapping_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnMapping",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnMapping_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnMapping",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Returns: A reference to a value in the map based on the two keys.
// Experimental.
func (c *jsiiProxy_CfnMapping) FindInMap(key1 *string, key2 *string) *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"findInMap",
		[]interface{}{key1, key2},
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
func (c *jsiiProxy_CfnMapping) OnPrepare() {
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
func (c *jsiiProxy_CfnMapping) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnMapping) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnMapping) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnMapping) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Sets a value in the map based on the two keys.
// Experimental.
func (c *jsiiProxy_CfnMapping) SetValue(key1 *string, key2 *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"setValue",
		[]interface{}{key1, key2, value},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnMapping) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnMapping) ToString() *string {
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
func (c *jsiiProxy_CfnMapping) Validate() *[]*string {
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
type CfnMappingProps struct {
	// Mapping of key to a set of corresponding set of named values.
	//
	// The key identifies a map of name-value pairs and must be unique within the mapping.
	//
	// For example, if you want to set values based on a region, you can create a mapping
	// that uses the region name as a key and contains the values you want to specify for
	// each specific region.
	// Experimental.
	Mapping *map[string]*map[string]interface{} `json:"mapping"`
}

// A CloudFormation `AWS::CloudFormation::ModuleDefaultVersion`.
type CfnModuleDefaultVersion interface {
	CfnResource
	IInspectable
	Arn() *string
	SetArn(val *string)
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	ModuleName() *string
	SetModuleName(val *string)
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	UpdatedProperites() *map[string]interface{}
	VersionId() *string
	SetVersionId(val *string)
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnModuleDefaultVersion
type jsiiProxy_CfnModuleDefaultVersion struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnModuleDefaultVersion) Arn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) ModuleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"moduleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleDefaultVersion) VersionId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"versionId",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::ModuleDefaultVersion`.
func NewCfnModuleDefaultVersion(scope Construct, id *string, props *CfnModuleDefaultVersionProps) CfnModuleDefaultVersion {
	_init_.Initialize()

	j := jsiiProxy_CfnModuleDefaultVersion{}

	_jsii_.Create(
		"monocdk.CfnModuleDefaultVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::ModuleDefaultVersion`.
func NewCfnModuleDefaultVersion_Override(c CfnModuleDefaultVersion, scope Construct, id *string, props *CfnModuleDefaultVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnModuleDefaultVersion",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnModuleDefaultVersion) SetArn(val *string) {
	_jsii_.Set(
		j,
		"arn",
		val,
	)
}

func (j *jsiiProxy_CfnModuleDefaultVersion) SetModuleName(val *string) {
	_jsii_.Set(
		j,
		"moduleName",
		val,
	)
}

func (j *jsiiProxy_CfnModuleDefaultVersion) SetVersionId(val *string) {
	_jsii_.Set(
		j,
		"versionId",
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
func CfnModuleDefaultVersion_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnModuleDefaultVersion",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnModuleDefaultVersion_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnModuleDefaultVersion",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnModuleDefaultVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnModuleDefaultVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnModuleDefaultVersion_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnModuleDefaultVersion",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnModuleDefaultVersion) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnModuleDefaultVersion) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnModuleDefaultVersion) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnModuleDefaultVersion) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) OnPrepare() {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnModuleDefaultVersion) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) ToString() *string {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) Validate() *[]*string {
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
func (c *jsiiProxy_CfnModuleDefaultVersion) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudFormation::ModuleDefaultVersion`.
type CfnModuleDefaultVersionProps struct {
	// `AWS::CloudFormation::ModuleDefaultVersion.Arn`.
	Arn *string `json:"arn"`
	// `AWS::CloudFormation::ModuleDefaultVersion.ModuleName`.
	ModuleName *string `json:"moduleName"`
	// `AWS::CloudFormation::ModuleDefaultVersion.VersionId`.
	VersionId *string `json:"versionId"`
}

// A CloudFormation `AWS::CloudFormation::ModuleVersion`.
type CfnModuleVersion interface {
	CfnResource
	IInspectable
	AttrArn() *string
	AttrDescription() *string
	AttrDocumentationUrl() *string
	AttrIsDefaultVersion() IResolvable
	AttrSchema() *string
	AttrTimeCreated() *string
	AttrVersionId() *string
	AttrVisibility() *string
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	ModuleName() *string
	SetModuleName(val *string)
	ModulePackage() *string
	SetModulePackage(val *string)
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnModuleVersion
type jsiiProxy_CfnModuleVersion struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnModuleVersion) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) AttrDescription() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDescription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) AttrDocumentationUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDocumentationUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) AttrIsDefaultVersion() IResolvable {
	var returns IResolvable
	_jsii_.Get(
		j,
		"attrIsDefaultVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) AttrSchema() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrSchema",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) AttrTimeCreated() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrTimeCreated",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) AttrVersionId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrVersionId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) AttrVisibility() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrVisibility",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) ModuleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"moduleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) ModulePackage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"modulePackage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModuleVersion) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::ModuleVersion`.
func NewCfnModuleVersion(scope Construct, id *string, props *CfnModuleVersionProps) CfnModuleVersion {
	_init_.Initialize()

	j := jsiiProxy_CfnModuleVersion{}

	_jsii_.Create(
		"monocdk.CfnModuleVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::ModuleVersion`.
func NewCfnModuleVersion_Override(c CfnModuleVersion, scope Construct, id *string, props *CfnModuleVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnModuleVersion",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnModuleVersion) SetModuleName(val *string) {
	_jsii_.Set(
		j,
		"moduleName",
		val,
	)
}

func (j *jsiiProxy_CfnModuleVersion) SetModulePackage(val *string) {
	_jsii_.Set(
		j,
		"modulePackage",
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
func CfnModuleVersion_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnModuleVersion",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnModuleVersion_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnModuleVersion",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnModuleVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnModuleVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnModuleVersion_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnModuleVersion",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnModuleVersion) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnModuleVersion) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnModuleVersion) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnModuleVersion) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnModuleVersion) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnModuleVersion) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnModuleVersion) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnModuleVersion) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnModuleVersion) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnModuleVersion) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnModuleVersion) OnPrepare() {
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
func (c *jsiiProxy_CfnModuleVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnModuleVersion) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnModuleVersion) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnModuleVersion) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnModuleVersion) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnModuleVersion) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnModuleVersion) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnModuleVersion) ToString() *string {
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
func (c *jsiiProxy_CfnModuleVersion) Validate() *[]*string {
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
func (c *jsiiProxy_CfnModuleVersion) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudFormation::ModuleVersion`.
type CfnModuleVersionProps struct {
	// `AWS::CloudFormation::ModuleVersion.ModuleName`.
	ModuleName *string `json:"moduleName"`
	// `AWS::CloudFormation::ModuleVersion.ModulePackage`.
	ModulePackage *string `json:"modulePackage"`
}

// Experimental.
type CfnOutput interface {
	CfnElement
	Condition() CfnCondition
	SetCondition(val CfnCondition)
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	ExportName() *string
	SetExportName(val *string)
	ImportValue() *string
	LogicalId() *string
	Node() ConstructNode
	Stack() Stack
	Value() interface{}
	SetValue(val interface{})
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnOutput
type jsiiProxy_CfnOutput struct {
	jsiiProxy_CfnElement
}

func (j *jsiiProxy_CfnOutput) Condition() CfnCondition {
	var returns CfnCondition
	_jsii_.Get(
		j,
		"condition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) ExportName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"exportName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) ImportValue() *string {
	var returns *string
	_jsii_.Get(
		j,
		"importValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOutput) Value() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


// Creates an CfnOutput value for this stack.
// Experimental.
func NewCfnOutput(scope constructs.Construct, id *string, props *CfnOutputProps) CfnOutput {
	_init_.Initialize()

	j := jsiiProxy_CfnOutput{}

	_jsii_.Create(
		"monocdk.CfnOutput",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates an CfnOutput value for this stack.
// Experimental.
func NewCfnOutput_Override(c CfnOutput, scope constructs.Construct, id *string, props *CfnOutputProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnOutput",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnOutput) SetCondition(val CfnCondition) {
	_jsii_.Set(
		j,
		"condition",
		val,
	)
}

func (j *jsiiProxy_CfnOutput) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnOutput) SetExportName(val *string) {
	_jsii_.Set(
		j,
		"exportName",
		val,
	)
}

func (j *jsiiProxy_CfnOutput) SetValue(val interface{}) {
	_jsii_.Set(
		j,
		"value",
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
func CfnOutput_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnOutput",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnOutput_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnOutput",
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
func (c *jsiiProxy_CfnOutput) OnPrepare() {
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
func (c *jsiiProxy_CfnOutput) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnOutput) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnOutput) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnOutput) Prepare() {
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
func (c *jsiiProxy_CfnOutput) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnOutput) ToString() *string {
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
// Experimental.
func (c *jsiiProxy_CfnOutput) Validate() *[]*string {
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
type CfnOutputProps struct {
	// The value of the property returned by the aws cloudformation describe-stacks command.
	//
	// The value of an output can include literals, parameter references, pseudo-parameters,
	// a mapping value, or intrinsic functions.
	// Experimental.
	Value *string `json:"value"`
	// A condition to associate with this output value.
	//
	// If the condition evaluates
	// to `false`, this output value will not be included in the stack.
	// Experimental.
	Condition CfnCondition `json:"condition"`
	// A String type that describes the output value.
	//
	// The description can be a maximum of 4 K in length.
	// Experimental.
	Description *string `json:"description"`
	// The name used to export the value of this output across stacks.
	//
	// To import the value from another stack, use `Fn.importValue(exportName)`.
	// Experimental.
	ExportName *string `json:"exportName"`
}

// A CloudFormation parameter.
//
// Use the optional Parameters section to customize your templates.
// Parameters enable you to input custom values to your template each time you create or
// update a stack.
// Experimental.
type CfnParameter interface {
	CfnElement
	AllowedPattern() *string
	SetAllowedPattern(val *string)
	AllowedValues() *[]*string
	SetAllowedValues(val *[]*string)
	ConstraintDescription() *string
	SetConstraintDescription(val *string)
	CreationStack() *[]*string
	Default() interface{}
	SetDefault(val interface{})
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	MaxLength() *float64
	SetMaxLength(val *float64)
	MaxValue() *float64
	SetMaxValue(val *float64)
	MinLength() *float64
	SetMinLength(val *float64)
	MinValue() *float64
	SetMinValue(val *float64)
	Node() ConstructNode
	NoEcho() *bool
	SetNoEcho(val *bool)
	Stack() Stack
	Type() *string
	SetType(val *string)
	Value() IResolvable
	ValueAsList() *[]*string
	ValueAsNumber() *float64
	ValueAsString() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	Resolve(_context IResolveContext) interface{}
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnParameter
type jsiiProxy_CfnParameter struct {
	jsiiProxy_CfnElement
}

func (j *jsiiProxy_CfnParameter) AllowedPattern() *string {
	var returns *string
	_jsii_.Get(
		j,
		"allowedPattern",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) AllowedValues() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"allowedValues",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) ConstraintDescription() *string {
	var returns *string
	_jsii_.Get(
		j,
		"constraintDescription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) Default() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"default",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) MaxLength() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxLength",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) MaxValue() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) MinLength() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"minLength",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) MinValue() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"minValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) NoEcho() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"noEcho",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) Value() IResolvable {
	var returns IResolvable
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) ValueAsList() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"valueAsList",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) ValueAsNumber() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"valueAsNumber",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnParameter) ValueAsString() *string {
	var returns *string
	_jsii_.Get(
		j,
		"valueAsString",
		&returns,
	)
	return returns
}


// Creates a parameter construct.
//
// Note that the name (logical ID) of the parameter will derive from it's `coname` and location
// within the stack. Therefore, it is recommended that parameters are defined at the stack level.
// Experimental.
func NewCfnParameter(scope constructs.Construct, id *string, props *CfnParameterProps) CfnParameter {
	_init_.Initialize()

	j := jsiiProxy_CfnParameter{}

	_jsii_.Create(
		"monocdk.CfnParameter",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates a parameter construct.
//
// Note that the name (logical ID) of the parameter will derive from it's `coname` and location
// within the stack. Therefore, it is recommended that parameters are defined at the stack level.
// Experimental.
func NewCfnParameter_Override(c CfnParameter, scope constructs.Construct, id *string, props *CfnParameterProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnParameter",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnParameter) SetAllowedPattern(val *string) {
	_jsii_.Set(
		j,
		"allowedPattern",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetAllowedValues(val *[]*string) {
	_jsii_.Set(
		j,
		"allowedValues",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetConstraintDescription(val *string) {
	_jsii_.Set(
		j,
		"constraintDescription",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetDefault(val interface{}) {
	_jsii_.Set(
		j,
		"default",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetMaxLength(val *float64) {
	_jsii_.Set(
		j,
		"maxLength",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetMaxValue(val *float64) {
	_jsii_.Set(
		j,
		"maxValue",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetMinLength(val *float64) {
	_jsii_.Set(
		j,
		"minLength",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetMinValue(val *float64) {
	_jsii_.Set(
		j,
		"minValue",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetNoEcho(val *bool) {
	_jsii_.Set(
		j,
		"noEcho",
		val,
	)
}

func (j *jsiiProxy_CfnParameter) SetType(val *string) {
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
func CfnParameter_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnParameter",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnParameter_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnParameter",
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
func (c *jsiiProxy_CfnParameter) OnPrepare() {
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
func (c *jsiiProxy_CfnParameter) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnParameter) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnParameter) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnParameter) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Experimental.
func (c *jsiiProxy_CfnParameter) Resolve(_context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnParameter) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnParameter) ToString() *string {
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
func (c *jsiiProxy_CfnParameter) Validate() *[]*string {
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
type CfnParameterProps struct {
	// A regular expression that represents the patterns to allow for String types.
	// Experimental.
	AllowedPattern *string `json:"allowedPattern"`
	// An array containing the list of values allowed for the parameter.
	// Experimental.
	AllowedValues *[]*string `json:"allowedValues"`
	// A string that explains a constraint when the constraint is violated.
	//
	// For example, without a constraint description, a parameter that has an allowed
	// pattern of [A-Za-z0-9]+ displays the following error message when the user specifies
	// an invalid value:
	// Experimental.
	ConstraintDescription *string `json:"constraintDescription"`
	// A value of the appropriate type for the template to use if no value is specified when a stack is created.
	//
	// If you define constraints for the parameter, you must specify
	// a value that adheres to those constraints.
	// Experimental.
	Default interface{} `json:"default"`
	// A string of up to 4000 characters that describes the parameter.
	// Experimental.
	Description *string `json:"description"`
	// An integer value that determines the largest number of characters you want to allow for String types.
	// Experimental.
	MaxLength *float64 `json:"maxLength"`
	// A numeric value that determines the largest numeric value you want to allow for Number types.
	// Experimental.
	MaxValue *float64 `json:"maxValue"`
	// An integer value that determines the smallest number of characters you want to allow for String types.
	// Experimental.
	MinLength *float64 `json:"minLength"`
	// A numeric value that determines the smallest numeric value you want to allow for Number types.
	// Experimental.
	MinValue *float64 `json:"minValue"`
	// Whether to mask the parameter value when anyone makes a call that describes the stack.
	//
	// If you set the value to ``true``, the parameter value is masked with asterisks (``*****``).
	// Experimental.
	NoEcho *bool `json:"noEcho"`
	// The data type for the parameter (DataType).
	// Experimental.
	Type *string `json:"type"`
}

// Base class for referenceable CloudFormation constructs which are not Resources.
//
// These constructs are things like Conditions and Parameters, can be
// referenced by taking the `.ref` attribute.
//
// Resource constructs do not inherit from CfnRefElement because they have their
// own, more specific types returned from the .ref attribute. Also, some
// resources aren't referenceable at all (such as BucketPolicies or GatewayAttachments).
// Experimental.
type CfnRefElement interface {
	CfnElement
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnRefElement
type jsiiProxy_CfnRefElement struct {
	jsiiProxy_CfnElement
}

func (j *jsiiProxy_CfnRefElement) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRefElement) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRefElement) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRefElement) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRefElement) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Creates an entity and binds it to a tree.
//
// Note that the root of the tree must be a Stack object (not just any Root).
// Experimental.
func NewCfnRefElement_Override(c CfnRefElement, scope constructs.Construct, id *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnRefElement",
		[]interface{}{scope, id},
		c,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnRefElement_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnRefElement",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnRefElement_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnRefElement",
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
func (c *jsiiProxy_CfnRefElement) OnPrepare() {
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
func (c *jsiiProxy_CfnRefElement) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRefElement) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnRefElement) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnRefElement) Prepare() {
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
func (c *jsiiProxy_CfnRefElement) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CfnRefElement) ToString() *string {
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
func (c *jsiiProxy_CfnRefElement) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents a CloudFormation resource.
// Experimental.
type CfnResource interface {
	CfnRefElement
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnResource
type jsiiProxy_CfnResource struct {
	jsiiProxy_CfnRefElement
}

func (j *jsiiProxy_CfnResource) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Creates a resource construct.
// Experimental.
func NewCfnResource(scope constructs.Construct, id *string, props *CfnResourceProps) CfnResource {
	_init_.Initialize()

	j := jsiiProxy_CfnResource{}

	_jsii_.Create(
		"monocdk.CfnResource",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates a resource construct.
// Experimental.
func NewCfnResource_Override(c CfnResource, scope constructs.Construct, id *string, props *CfnResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnResource",
		[]interface{}{scope, id, props},
		c,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnResource_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResource",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnResource_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResource",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnResource_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResource",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnResource) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnResource) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnResource) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnResource) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnResource) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnResource) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnResource) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnResource) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnResource) GetMetadata(key *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"getMetadata",
		[]interface{}{key},
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
func (c *jsiiProxy_CfnResource) OnPrepare() {
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
func (c *jsiiProxy_CfnResource) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnResource) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnResource) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnResource) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Experimental.
func (c *jsiiProxy_CfnResource) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnResource) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnResource) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnResource) ToString() *string {
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
func (c *jsiiProxy_CfnResource) Validate() *[]*string {
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
func (c *jsiiProxy_CfnResource) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// For an Auto Scaling group replacement update, specifies how many instances must signal success for the update to succeed.
// Experimental.
type CfnResourceAutoScalingCreationPolicy struct {
	// Specifies the percentage of instances in an Auto Scaling replacement update that must signal success for the update to succeed.
	//
	// You can specify a value from 0 to 100. AWS CloudFormation rounds to the nearest tenth of a percent.
	// For example, if you update five instances with a minimum successful percentage of 50, three instances must signal success.
	// If an instance doesn't send a signal within the time specified by the Timeout property, AWS CloudFormation assumes that the
	// instance wasn't created.
	// Experimental.
	MinSuccessfulInstancesPercent *float64 `json:"minSuccessfulInstancesPercent"`
}

// A CloudFormation `AWS::CloudFormation::ResourceDefaultVersion`.
type CfnResourceDefaultVersion interface {
	CfnResource
	IInspectable
	AttrArn() *string
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	TypeName() *string
	SetTypeName(val *string)
	TypeVersionArn() *string
	SetTypeVersionArn(val *string)
	UpdatedProperites() *map[string]interface{}
	VersionId() *string
	SetVersionId(val *string)
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnResourceDefaultVersion
type jsiiProxy_CfnResourceDefaultVersion struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnResourceDefaultVersion) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) TypeName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"typeName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) TypeVersionArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"typeVersionArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceDefaultVersion) VersionId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"versionId",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::ResourceDefaultVersion`.
func NewCfnResourceDefaultVersion(scope Construct, id *string, props *CfnResourceDefaultVersionProps) CfnResourceDefaultVersion {
	_init_.Initialize()

	j := jsiiProxy_CfnResourceDefaultVersion{}

	_jsii_.Create(
		"monocdk.CfnResourceDefaultVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::ResourceDefaultVersion`.
func NewCfnResourceDefaultVersion_Override(c CfnResourceDefaultVersion, scope Construct, id *string, props *CfnResourceDefaultVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnResourceDefaultVersion",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnResourceDefaultVersion) SetTypeName(val *string) {
	_jsii_.Set(
		j,
		"typeName",
		val,
	)
}

func (j *jsiiProxy_CfnResourceDefaultVersion) SetTypeVersionArn(val *string) {
	_jsii_.Set(
		j,
		"typeVersionArn",
		val,
	)
}

func (j *jsiiProxy_CfnResourceDefaultVersion) SetVersionId(val *string) {
	_jsii_.Set(
		j,
		"versionId",
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
func CfnResourceDefaultVersion_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResourceDefaultVersion",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnResourceDefaultVersion_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResourceDefaultVersion",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnResourceDefaultVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResourceDefaultVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnResourceDefaultVersion_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnResourceDefaultVersion",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnResourceDefaultVersion) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnResourceDefaultVersion) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnResourceDefaultVersion) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnResourceDefaultVersion) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) OnPrepare() {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnResourceDefaultVersion) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) ToString() *string {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) Validate() *[]*string {
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
func (c *jsiiProxy_CfnResourceDefaultVersion) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudFormation::ResourceDefaultVersion`.
type CfnResourceDefaultVersionProps struct {
	// `AWS::CloudFormation::ResourceDefaultVersion.TypeName`.
	TypeName *string `json:"typeName"`
	// `AWS::CloudFormation::ResourceDefaultVersion.TypeVersionArn`.
	TypeVersionArn *string `json:"typeVersionArn"`
	// `AWS::CloudFormation::ResourceDefaultVersion.VersionId`.
	VersionId *string `json:"versionId"`
}

// Experimental.
type CfnResourceProps struct {
	// CloudFormation resource type (e.g. `AWS::S3::Bucket`).
	// Experimental.
	Type *string `json:"type"`
	// Resource properties.
	// Experimental.
	Properties *map[string]interface{} `json:"properties"`
}

// When AWS CloudFormation creates the associated resource, configures the number of required success signals and the length of time that AWS CloudFormation waits for those signals.
// Experimental.
type CfnResourceSignal struct {
	// The number of success signals AWS CloudFormation must receive before it sets the resource status as CREATE_COMPLETE.
	//
	// If the resource receives a failure signal or doesn't receive the specified number of signals before the timeout period
	// expires, the resource creation fails and AWS CloudFormation rolls the stack back.
	// Experimental.
	Count *float64 `json:"count"`
	// The length of time that AWS CloudFormation waits for the number of signals that was specified in the Count property.
	//
	// The timeout period starts after AWS CloudFormation starts creating the resource, and the timeout expires no sooner
	// than the time you specify but can occur shortly thereafter. The maximum time that you can specify is 12 hours.
	// Experimental.
	Timeout *string `json:"timeout"`
}

// A CloudFormation `AWS::CloudFormation::ResourceVersion`.
type CfnResourceVersion interface {
	CfnResource
	IInspectable
	AttrArn() *string
	AttrIsDefaultVersion() IResolvable
	AttrProvisioningType() *string
	AttrTypeArn() *string
	AttrVersionId() *string
	AttrVisibility() *string
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	ExecutionRoleArn() *string
	SetExecutionRoleArn(val *string)
	LoggingConfig() interface{}
	SetLoggingConfig(val interface{})
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	SchemaHandlerPackage() *string
	SetSchemaHandlerPackage(val *string)
	Stack() Stack
	TypeName() *string
	SetTypeName(val *string)
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnResourceVersion
type jsiiProxy_CfnResourceVersion struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnResourceVersion) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) AttrIsDefaultVersion() IResolvable {
	var returns IResolvable
	_jsii_.Get(
		j,
		"attrIsDefaultVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) AttrProvisioningType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrProvisioningType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) AttrTypeArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrTypeArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) AttrVersionId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrVersionId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) AttrVisibility() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrVisibility",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) ExecutionRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"executionRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) LoggingConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"loggingConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) SchemaHandlerPackage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"schemaHandlerPackage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) TypeName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"typeName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResourceVersion) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::ResourceVersion`.
func NewCfnResourceVersion(scope Construct, id *string, props *CfnResourceVersionProps) CfnResourceVersion {
	_init_.Initialize()

	j := jsiiProxy_CfnResourceVersion{}

	_jsii_.Create(
		"monocdk.CfnResourceVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::ResourceVersion`.
func NewCfnResourceVersion_Override(c CfnResourceVersion, scope Construct, id *string, props *CfnResourceVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnResourceVersion",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnResourceVersion) SetExecutionRoleArn(val *string) {
	_jsii_.Set(
		j,
		"executionRoleArn",
		val,
	)
}

func (j *jsiiProxy_CfnResourceVersion) SetLoggingConfig(val interface{}) {
	_jsii_.Set(
		j,
		"loggingConfig",
		val,
	)
}

func (j *jsiiProxy_CfnResourceVersion) SetSchemaHandlerPackage(val *string) {
	_jsii_.Set(
		j,
		"schemaHandlerPackage",
		val,
	)
}

func (j *jsiiProxy_CfnResourceVersion) SetTypeName(val *string) {
	_jsii_.Set(
		j,
		"typeName",
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
func CfnResourceVersion_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResourceVersion",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnResourceVersion_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResourceVersion",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnResourceVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnResourceVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnResourceVersion_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnResourceVersion",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnResourceVersion) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnResourceVersion) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnResourceVersion) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnResourceVersion) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnResourceVersion) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnResourceVersion) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnResourceVersion) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnResourceVersion) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnResourceVersion) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnResourceVersion) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnResourceVersion) OnPrepare() {
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
func (c *jsiiProxy_CfnResourceVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnResourceVersion) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnResourceVersion) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnResourceVersion) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnResourceVersion) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnResourceVersion) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnResourceVersion) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnResourceVersion) ToString() *string {
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
func (c *jsiiProxy_CfnResourceVersion) Validate() *[]*string {
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
func (c *jsiiProxy_CfnResourceVersion) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnResourceVersion_LoggingConfigProperty struct {
	// `CfnResourceVersion.LoggingConfigProperty.LogGroupName`.
	LogGroupName *string `json:"logGroupName"`
	// `CfnResourceVersion.LoggingConfigProperty.LogRoleArn`.
	LogRoleArn *string `json:"logRoleArn"`
}

// Properties for defining a `AWS::CloudFormation::ResourceVersion`.
type CfnResourceVersionProps struct {
	// `AWS::CloudFormation::ResourceVersion.SchemaHandlerPackage`.
	SchemaHandlerPackage *string `json:"schemaHandlerPackage"`
	// `AWS::CloudFormation::ResourceVersion.TypeName`.
	TypeName *string `json:"typeName"`
	// `AWS::CloudFormation::ResourceVersion.ExecutionRoleArn`.
	ExecutionRoleArn *string `json:"executionRoleArn"`
	// `AWS::CloudFormation::ResourceVersion.LoggingConfig`.
	LoggingConfig interface{} `json:"loggingConfig"`
}

// The Rules that define template constraints in an AWS Service Catalog portfolio describe when end users can use the template and which values they can specify for parameters that are declared in the AWS CloudFormation template used to create the product they are attempting to use.
//
// Rules
// are useful for preventing end users from inadvertently specifying an incorrect value.
// For example, you can add a rule to verify whether end users specified a valid subnet in a
// given VPC or used m1.small instance types for test environments. AWS CloudFormation uses
// rules to validate parameter values before it creates the resources for the product.
//
// A rule can include a RuleCondition property and must include an Assertions property.
// For each rule, you can define only one rule condition; you can define one or more asserts within the Assertions property.
// You define a rule condition and assertions by using rule-specific intrinsic functions.
// Experimental.
type CfnRule interface {
	CfnRefElement
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	AddAssertion(condition ICfnConditionExpression, description *string)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CfnRule
type jsiiProxy_CfnRule struct {
	jsiiProxy_CfnRefElement
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

func (j *jsiiProxy_CfnRule) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRule) Node() ConstructNode {
	var returns ConstructNode
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

func (j *jsiiProxy_CfnRule) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Creates and adds a rule.
// Experimental.
func NewCfnRule(scope constructs.Construct, id *string, props *CfnRuleProps) CfnRule {
	_init_.Initialize()

	j := jsiiProxy_CfnRule{}

	_jsii_.Create(
		"monocdk.CfnRule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates and adds a rule.
// Experimental.
func NewCfnRule_Override(c CfnRule, scope constructs.Construct, id *string, props *CfnRuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnRule",
		[]interface{}{scope, id, props},
		c,
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
		"monocdk.CfnRule",
		"isCfnElement",
		[]interface{}{x},
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
		"monocdk.CfnRule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Adds an assertion to the rule.
// Experimental.
func (c *jsiiProxy_CfnRule) AddAssertion(condition ICfnConditionExpression, description *string) {
	_jsii_.InvokeVoid(
		c,
		"addAssertion",
		[]interface{}{condition, description},
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

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CfnRule) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
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

// A rule assertion.
// Experimental.
type CfnRuleAssertion struct {
	// The assertion.
	// Experimental.
	Assert ICfnConditionExpression `json:"assert"`
	// The assertion description.
	// Experimental.
	AssertDescription *string `json:"assertDescription"`
}

// A rule can include a RuleCondition property and must include an Assertions property.
//
// For each rule, you can define only one rule condition; you can define one or more asserts within the Assertions property.
// You define a rule condition and assertions by using rule-specific intrinsic functions.
//
// You can use the following rule-specific intrinsic functions to define rule conditions and assertions:
//
//   Fn::And
//   Fn::Contains
//   Fn::EachMemberEquals
//   Fn::EachMemberIn
//   Fn::Equals
//   Fn::If
//   Fn::Not
//   Fn::Or
//   Fn::RefAll
//   Fn::ValueOf
//   Fn::ValueOfAll
//
// https://docs.aws.amazon.com/servicecatalog/latest/adminguide/reference-template_constraint_rules.html
// Experimental.
type CfnRuleProps struct {
	// Assertions which define the rule.
	// Experimental.
	Assertions *[]*CfnRuleAssertion `json:"assertions"`
	// If the rule condition evaluates to false, the rule doesn't take effect.
	//
	// If the function in the rule condition evaluates to true, expressions in each assert are evaluated and applied.
	// Experimental.
	RuleCondition ICfnConditionExpression `json:"ruleCondition"`
}

// A CloudFormation `AWS::CloudFormation::Stack`.
type CfnStack interface {
	CfnResource
	IInspectable
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	NotificationArns() *[]*string
	SetNotificationArns(val *[]*string)
	Parameters() interface{}
	SetParameters(val interface{})
	Ref() *string
	Stack() Stack
	Tags() TagManager
	TemplateUrl() *string
	SetTemplateUrl(val *string)
	TimeoutInMinutes() *float64
	SetTimeoutInMinutes(val *float64)
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnStack
type jsiiProxy_CfnStack struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnStack) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) NotificationArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"notificationArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) Parameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) Tags() TagManager {
	var returns TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) TemplateUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) TimeoutInMinutes() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"timeoutInMinutes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStack) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::Stack`.
func NewCfnStack(scope Construct, id *string, props *CfnStackProps) CfnStack {
	_init_.Initialize()

	j := jsiiProxy_CfnStack{}

	_jsii_.Create(
		"monocdk.CfnStack",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::Stack`.
func NewCfnStack_Override(c CfnStack, scope Construct, id *string, props *CfnStackProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnStack",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStack) SetNotificationArns(val *[]*string) {
	_jsii_.Set(
		j,
		"notificationArns",
		val,
	)
}

func (j *jsiiProxy_CfnStack) SetParameters(val interface{}) {
	_jsii_.Set(
		j,
		"parameters",
		val,
	)
}

func (j *jsiiProxy_CfnStack) SetTemplateUrl(val *string) {
	_jsii_.Set(
		j,
		"templateUrl",
		val,
	)
}

func (j *jsiiProxy_CfnStack) SetTimeoutInMinutes(val *float64) {
	_jsii_.Set(
		j,
		"timeoutInMinutes",
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
func CfnStack_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnStack",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnStack_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnStack",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnStack_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnStack",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStack_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnStack",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnStack) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnStack) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnStack) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnStack) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnStack) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnStack) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnStack) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnStack) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnStack) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnStack) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnStack) OnPrepare() {
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
func (c *jsiiProxy_CfnStack) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStack) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnStack) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnStack) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnStack) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnStack) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnStack) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnStack) ToString() *string {
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
func (c *jsiiProxy_CfnStack) Validate() *[]*string {
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
func (c *jsiiProxy_CfnStack) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudFormation::Stack`.
type CfnStackProps struct {
	// `AWS::CloudFormation::Stack.TemplateURL`.
	TemplateUrl *string `json:"templateUrl"`
	// `AWS::CloudFormation::Stack.NotificationARNs`.
	NotificationArns *[]*string `json:"notificationArns"`
	// `AWS::CloudFormation::Stack.Parameters`.
	Parameters interface{} `json:"parameters"`
	// `AWS::CloudFormation::Stack.Tags`.
	Tags *[]*CfnTag `json:"tags"`
	// `AWS::CloudFormation::Stack.TimeoutInMinutes`.
	TimeoutInMinutes *float64 `json:"timeoutInMinutes"`
}

// A CloudFormation `AWS::CloudFormation::StackSet`.
type CfnStackSet interface {
	CfnResource
	IInspectable
	AdministrationRoleArn() *string
	SetAdministrationRoleArn(val *string)
	AttrStackSetId() *string
	AutoDeployment() interface{}
	SetAutoDeployment(val interface{})
	CallAs() *string
	SetCallAs(val *string)
	Capabilities() *[]*string
	SetCapabilities(val *[]*string)
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	ExecutionRoleName() *string
	SetExecutionRoleName(val *string)
	LogicalId() *string
	Node() ConstructNode
	OperationPreferences() interface{}
	SetOperationPreferences(val interface{})
	Parameters() interface{}
	SetParameters(val interface{})
	PermissionModel() *string
	SetPermissionModel(val *string)
	Ref() *string
	Stack() Stack
	StackInstancesGroup() interface{}
	SetStackInstancesGroup(val interface{})
	StackSetName() *string
	SetStackSetName(val *string)
	Tags() TagManager
	TemplateBody() *string
	SetTemplateBody(val *string)
	TemplateUrl() *string
	SetTemplateUrl(val *string)
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnStackSet
type jsiiProxy_CfnStackSet struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnStackSet) AdministrationRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"administrationRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) AttrStackSetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrStackSetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) AutoDeployment() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"autoDeployment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) CallAs() *string {
	var returns *string
	_jsii_.Get(
		j,
		"callAs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) Capabilities() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"capabilities",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) ExecutionRoleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"executionRoleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) OperationPreferences() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"operationPreferences",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) Parameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) PermissionModel() *string {
	var returns *string
	_jsii_.Get(
		j,
		"permissionModel",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) StackInstancesGroup() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stackInstancesGroup",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) StackSetName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackSetName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) Tags() TagManager {
	var returns TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) TemplateBody() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateBody",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) TemplateUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStackSet) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::StackSet`.
func NewCfnStackSet(scope Construct, id *string, props *CfnStackSetProps) CfnStackSet {
	_init_.Initialize()

	j := jsiiProxy_CfnStackSet{}

	_jsii_.Create(
		"monocdk.CfnStackSet",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::StackSet`.
func NewCfnStackSet_Override(c CfnStackSet, scope Construct, id *string, props *CfnStackSetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnStackSet",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStackSet) SetAdministrationRoleArn(val *string) {
	_jsii_.Set(
		j,
		"administrationRoleArn",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetAutoDeployment(val interface{}) {
	_jsii_.Set(
		j,
		"autoDeployment",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetCallAs(val *string) {
	_jsii_.Set(
		j,
		"callAs",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetCapabilities(val *[]*string) {
	_jsii_.Set(
		j,
		"capabilities",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetExecutionRoleName(val *string) {
	_jsii_.Set(
		j,
		"executionRoleName",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetOperationPreferences(val interface{}) {
	_jsii_.Set(
		j,
		"operationPreferences",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetParameters(val interface{}) {
	_jsii_.Set(
		j,
		"parameters",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetPermissionModel(val *string) {
	_jsii_.Set(
		j,
		"permissionModel",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetStackInstancesGroup(val interface{}) {
	_jsii_.Set(
		j,
		"stackInstancesGroup",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetStackSetName(val *string) {
	_jsii_.Set(
		j,
		"stackSetName",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetTemplateBody(val *string) {
	_jsii_.Set(
		j,
		"templateBody",
		val,
	)
}

func (j *jsiiProxy_CfnStackSet) SetTemplateUrl(val *string) {
	_jsii_.Set(
		j,
		"templateUrl",
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
func CfnStackSet_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnStackSet",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnStackSet_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnStackSet",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnStackSet_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnStackSet",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStackSet_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnStackSet",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnStackSet) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnStackSet) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnStackSet) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnStackSet) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnStackSet) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnStackSet) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnStackSet) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnStackSet) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnStackSet) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnStackSet) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnStackSet) OnPrepare() {
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
func (c *jsiiProxy_CfnStackSet) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStackSet) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnStackSet) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnStackSet) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnStackSet) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnStackSet) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnStackSet) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnStackSet) ToString() *string {
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
func (c *jsiiProxy_CfnStackSet) Validate() *[]*string {
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
func (c *jsiiProxy_CfnStackSet) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnStackSet_AutoDeploymentProperty struct {
	// `CfnStackSet.AutoDeploymentProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `CfnStackSet.AutoDeploymentProperty.RetainStacksOnAccountRemoval`.
	RetainStacksOnAccountRemoval interface{} `json:"retainStacksOnAccountRemoval"`
}

type CfnStackSet_DeploymentTargetsProperty struct {
	// `CfnStackSet.DeploymentTargetsProperty.Accounts`.
	Accounts *[]*string `json:"accounts"`
	// `CfnStackSet.DeploymentTargetsProperty.OrganizationalUnitIds`.
	OrganizationalUnitIds *[]*string `json:"organizationalUnitIds"`
}

type CfnStackSet_OperationPreferencesProperty struct {
	// `CfnStackSet.OperationPreferencesProperty.FailureToleranceCount`.
	FailureToleranceCount *float64 `json:"failureToleranceCount"`
	// `CfnStackSet.OperationPreferencesProperty.FailureTolerancePercentage`.
	FailureTolerancePercentage *float64 `json:"failureTolerancePercentage"`
	// `CfnStackSet.OperationPreferencesProperty.MaxConcurrentCount`.
	MaxConcurrentCount *float64 `json:"maxConcurrentCount"`
	// `CfnStackSet.OperationPreferencesProperty.MaxConcurrentPercentage`.
	MaxConcurrentPercentage *float64 `json:"maxConcurrentPercentage"`
	// `CfnStackSet.OperationPreferencesProperty.RegionConcurrencyType`.
	RegionConcurrencyType *string `json:"regionConcurrencyType"`
	// `CfnStackSet.OperationPreferencesProperty.RegionOrder`.
	RegionOrder *[]*string `json:"regionOrder"`
}

type CfnStackSet_ParameterProperty struct {
	// `CfnStackSet.ParameterProperty.ParameterKey`.
	ParameterKey *string `json:"parameterKey"`
	// `CfnStackSet.ParameterProperty.ParameterValue`.
	ParameterValue *string `json:"parameterValue"`
}

type CfnStackSet_StackInstancesProperty struct {
	// `CfnStackSet.StackInstancesProperty.DeploymentTargets`.
	DeploymentTargets interface{} `json:"deploymentTargets"`
	// `CfnStackSet.StackInstancesProperty.Regions`.
	Regions *[]*string `json:"regions"`
	// `CfnStackSet.StackInstancesProperty.ParameterOverrides`.
	ParameterOverrides interface{} `json:"parameterOverrides"`
}

// Properties for defining a `AWS::CloudFormation::StackSet`.
type CfnStackSetProps struct {
	// `AWS::CloudFormation::StackSet.PermissionModel`.
	PermissionModel *string `json:"permissionModel"`
	// `AWS::CloudFormation::StackSet.StackSetName`.
	StackSetName *string `json:"stackSetName"`
	// `AWS::CloudFormation::StackSet.AdministrationRoleARN`.
	AdministrationRoleArn *string `json:"administrationRoleArn"`
	// `AWS::CloudFormation::StackSet.AutoDeployment`.
	AutoDeployment interface{} `json:"autoDeployment"`
	// `AWS::CloudFormation::StackSet.CallAs`.
	CallAs *string `json:"callAs"`
	// `AWS::CloudFormation::StackSet.Capabilities`.
	Capabilities *[]*string `json:"capabilities"`
	// `AWS::CloudFormation::StackSet.Description`.
	Description *string `json:"description"`
	// `AWS::CloudFormation::StackSet.ExecutionRoleName`.
	ExecutionRoleName *string `json:"executionRoleName"`
	// `AWS::CloudFormation::StackSet.OperationPreferences`.
	OperationPreferences interface{} `json:"operationPreferences"`
	// `AWS::CloudFormation::StackSet.Parameters`.
	Parameters interface{} `json:"parameters"`
	// `AWS::CloudFormation::StackSet.StackInstancesGroup`.
	StackInstancesGroup interface{} `json:"stackInstancesGroup"`
	// `AWS::CloudFormation::StackSet.Tags`.
	Tags *[]*CfnTag `json:"tags"`
	// `AWS::CloudFormation::StackSet.TemplateBody`.
	TemplateBody *string `json:"templateBody"`
	// `AWS::CloudFormation::StackSet.TemplateURL`.
	TemplateUrl *string `json:"templateUrl"`
}

// Experimental.
type CfnTag struct {
	// Experimental.
	Key *string `json:"key"`
	// Experimental.
	Value *string `json:"value"`
}

// A traffic route, representing where the traffic is being directed to.
// Experimental.
type CfnTrafficRoute struct {
	// The logical id of the target resource.
	// Experimental.
	LogicalId *string `json:"logicalId"`
	// The resource type of the route.
	//
	// Today, the only allowed value is 'AWS::ElasticLoadBalancingV2::Listener'.
	// Experimental.
	Type *string `json:"type"`
}

// Type of the {@link CfnCodeDeployBlueGreenEcsAttributes.trafficRouting} property.
// Experimental.
type CfnTrafficRouting struct {
	// The listener to be used by your load balancer to direct traffic to your target groups.
	// Experimental.
	ProdTrafficRoute *CfnTrafficRoute `json:"prodTrafficRoute"`
	// The logical IDs of the blue and green, respectively, AWS::ElasticLoadBalancingV2::TargetGroup target groups.
	// Experimental.
	TargetGroups *[]*string `json:"targetGroups"`
	// The listener to be used by your load balancer to direct traffic to your target groups.
	// Experimental.
	TestTrafficRoute *CfnTrafficRoute `json:"testTrafficRoute"`
}

// Traffic routing configuration settings.
//
// The type of the {@link CfnCodeDeployBlueGreenHookProps.trafficRoutingConfig} property.
// Experimental.
type CfnTrafficRoutingConfig struct {
	// The type of traffic shifting used by the blue-green deployment configuration.
	// Experimental.
	Type CfnTrafficRoutingType `json:"type"`
	// The configuration for traffic routing when {@link type} is {@link CfnTrafficRoutingType.TIME_BASED_CANARY}.
	// Experimental.
	TimeBasedCanary *CfnTrafficRoutingTimeBasedCanary `json:"timeBasedCanary"`
	// The configuration for traffic routing when {@link type} is {@link CfnTrafficRoutingType.TIME_BASED_LINEAR}.
	// Experimental.
	TimeBasedLinear *CfnTrafficRoutingTimeBasedLinear `json:"timeBasedLinear"`
}

// The traffic routing configuration if {@link CfnTrafficRoutingConfig.type} is {@link CfnTrafficRoutingType.TIME_BASED_CANARY}.
// Experimental.
type CfnTrafficRoutingTimeBasedCanary struct {
	// The number of minutes between the first and second traffic shifts of a time-based canary deployment.
	// Experimental.
	BakeTimeMins *float64 `json:"bakeTimeMins"`
	// The percentage of traffic to shift in the first increment of a time-based canary deployment.
	//
	// The step percentage must be 14% or greater.
	// Experimental.
	StepPercentage *float64 `json:"stepPercentage"`
}

// The traffic routing configuration if {@link CfnTrafficRoutingConfig.type} is {@link CfnTrafficRoutingType.TIME_BASED_LINEAR}.
// Experimental.
type CfnTrafficRoutingTimeBasedLinear struct {
	// The number of minutes between the first and second traffic shifts of a time-based linear deployment.
	// Experimental.
	BakeTimeMins *float64 `json:"bakeTimeMins"`
	// The percentage of traffic that is shifted at the start of each increment of a time-based linear deployment.
	//
	// The step percentage must be 14% or greater.
	// Experimental.
	StepPercentage *float64 `json:"stepPercentage"`
}

// The possible types of traffic shifting for the blue-green deployment configuration.
//
// The type of the {@link CfnTrafficRoutingConfig.type} property.
// Experimental.
type CfnTrafficRoutingType string

const (
	CfnTrafficRoutingType_ALL_AT_ONCE CfnTrafficRoutingType = "ALL_AT_ONCE"
	CfnTrafficRoutingType_TIME_BASED_CANARY CfnTrafficRoutingType = "TIME_BASED_CANARY"
	CfnTrafficRoutingType_TIME_BASED_LINEAR CfnTrafficRoutingType = "TIME_BASED_LINEAR"
)

// Use the UpdatePolicy attribute to specify how AWS CloudFormation handles updates to the AWS::AutoScaling::AutoScalingGroup resource.
//
// AWS CloudFormation invokes one of three update policies depending on the type of change you make or whether a
// scheduled action is associated with the Auto Scaling group.
// Experimental.
type CfnUpdatePolicy struct {
	// Specifies whether an Auto Scaling group and the instances it contains are replaced during an update.
	//
	// During replacement,
	// AWS CloudFormation retains the old group until it finishes creating the new one. If the update fails, AWS CloudFormation
	// can roll back to the old Auto Scaling group and delete the new Auto Scaling group.
	// Experimental.
	AutoScalingReplacingUpdate *CfnAutoScalingReplacingUpdate `json:"autoScalingReplacingUpdate"`
	// To specify how AWS CloudFormation handles rolling updates for an Auto Scaling group, use the AutoScalingRollingUpdate policy.
	//
	// Rolling updates enable you to specify whether AWS CloudFormation updates instances that are in an Auto Scaling
	// group in batches or all at once.
	// Experimental.
	AutoScalingRollingUpdate *CfnAutoScalingRollingUpdate `json:"autoScalingRollingUpdate"`
	// To specify how AWS CloudFormation handles updates for the MinSize, MaxSize, and DesiredCapacity properties when the AWS::AutoScaling::AutoScalingGroup resource has an associated scheduled action, use the AutoScalingScheduledAction policy.
	// Experimental.
	AutoScalingScheduledAction *CfnAutoScalingScheduledAction `json:"autoScalingScheduledAction"`
	// To perform an AWS CodeDeploy deployment when the version changes on an AWS::Lambda::Alias resource, use the CodeDeployLambdaAliasUpdate update policy.
	// Experimental.
	CodeDeployLambdaAliasUpdate *CfnCodeDeployLambdaAliasUpdate `json:"codeDeployLambdaAliasUpdate"`
	// To upgrade an Amazon ES domain to a new version of Elasticsearch rather than replacing the entire AWS::Elasticsearch::Domain resource, use the EnableVersionUpgrade update policy.
	// Experimental.
	EnableVersionUpgrade *bool `json:"enableVersionUpgrade"`
	// To modify a replication group's shards by adding or removing shards, rather than replacing the entire AWS::ElastiCache::ReplicationGroup resource, use the UseOnlineResharding update policy.
	// Experimental.
	UseOnlineResharding *bool `json:"useOnlineResharding"`
}

// A CloudFormation `AWS::CloudFormation::WaitCondition`.
type CfnWaitCondition interface {
	CfnResource
	IInspectable
	AttrData() IResolvable
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Count() *float64
	SetCount(val *float64)
	CreationStack() *[]*string
	Handle() *string
	SetHandle(val *string)
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	Timeout() *string
	SetTimeout(val *string)
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnWaitCondition
type jsiiProxy_CfnWaitCondition struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnWaitCondition) AttrData() IResolvable {
	var returns IResolvable
	_jsii_.Get(
		j,
		"attrData",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) Count() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"count",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) Handle() *string {
	var returns *string
	_jsii_.Get(
		j,
		"handle",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) Timeout() *string {
	var returns *string
	_jsii_.Get(
		j,
		"timeout",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitCondition) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::WaitCondition`.
func NewCfnWaitCondition(scope Construct, id *string, props *CfnWaitConditionProps) CfnWaitCondition {
	_init_.Initialize()

	j := jsiiProxy_CfnWaitCondition{}

	_jsii_.Create(
		"monocdk.CfnWaitCondition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::WaitCondition`.
func NewCfnWaitCondition_Override(c CfnWaitCondition, scope Construct, id *string, props *CfnWaitConditionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnWaitCondition",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnWaitCondition) SetCount(val *float64) {
	_jsii_.Set(
		j,
		"count",
		val,
	)
}

func (j *jsiiProxy_CfnWaitCondition) SetHandle(val *string) {
	_jsii_.Set(
		j,
		"handle",
		val,
	)
}

func (j *jsiiProxy_CfnWaitCondition) SetTimeout(val *string) {
	_jsii_.Set(
		j,
		"timeout",
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
func CfnWaitCondition_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnWaitCondition",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnWaitCondition_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnWaitCondition",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnWaitCondition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnWaitCondition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnWaitCondition_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnWaitCondition",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnWaitCondition) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnWaitCondition) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnWaitCondition) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnWaitCondition) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnWaitCondition) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnWaitCondition) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnWaitCondition) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnWaitCondition) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnWaitCondition) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnWaitCondition) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnWaitCondition) OnPrepare() {
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
func (c *jsiiProxy_CfnWaitCondition) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnWaitCondition) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnWaitCondition) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnWaitCondition) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnWaitCondition) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnWaitCondition) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnWaitCondition) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnWaitCondition) ToString() *string {
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
func (c *jsiiProxy_CfnWaitCondition) Validate() *[]*string {
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
func (c *jsiiProxy_CfnWaitCondition) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// A CloudFormation `AWS::CloudFormation::WaitConditionHandle`.
type CfnWaitConditionHandle interface {
	CfnResource
	IInspectable
	CfnOptions() ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() ConstructNode
	Ref() *string
	Stack() Stack
	UpdatedProperites() *map[string]interface{}
	AddDeletionOverride(path *string)
	AddDependsOn(target CfnResource)
	AddMetadata(key *string, value interface{})
	AddOverride(path *string, value interface{})
	AddPropertyDeletionOverride(propertyPath *string)
	AddPropertyOverride(propertyPath *string, value interface{})
	ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions)
	GetAtt(attributeName *string) Reference
	GetMetadata(key *string) interface{}
	Inspect(inspector TreeInspector)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	OverrideLogicalId(newLogicalId *string)
	Prepare()
	RenderProperties(props *map[string]interface{}) *map[string]interface{}
	ShouldSynthesize() *bool
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	ValidateProperties(_properties interface{})
}

// The jsii proxy struct for CfnWaitConditionHandle
type jsiiProxy_CfnWaitConditionHandle struct {
	jsiiProxy_CfnResource
	jsiiProxy_IInspectable
}

func (j *jsiiProxy_CfnWaitConditionHandle) CfnOptions() ICfnResourceOptions {
	var returns ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWaitConditionHandle) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CloudFormation::WaitConditionHandle`.
func NewCfnWaitConditionHandle(scope Construct, id *string) CfnWaitConditionHandle {
	_init_.Initialize()

	j := jsiiProxy_CfnWaitConditionHandle{}

	_jsii_.Create(
		"monocdk.CfnWaitConditionHandle",
		[]interface{}{scope, id},
		&j,
	)

	return &j
}

// Create a new `AWS::CloudFormation::WaitConditionHandle`.
func NewCfnWaitConditionHandle_Override(c CfnWaitConditionHandle, scope Construct, id *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CfnWaitConditionHandle",
		[]interface{}{scope, id},
		c,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Experimental.
func CfnWaitConditionHandle_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnWaitConditionHandle",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnWaitConditionHandle_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnWaitConditionHandle",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnWaitConditionHandle_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CfnWaitConditionHandle",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnWaitConditionHandle_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.CfnWaitConditionHandle",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnWaitConditionHandle) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) AddDependsOn(target CfnResource) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnWaitConditionHandle) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnWaitConditionHandle) ApplyRemovalPolicy(policy RemovalPolicy, options *RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) GetAtt(attributeName *string) Reference {
	var returns Reference

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
func (c *jsiiProxy_CfnWaitConditionHandle) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnWaitConditionHandle) Inspect(inspector TreeInspector) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) OnPrepare() {
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
func (c *jsiiProxy_CfnWaitConditionHandle) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnWaitConditionHandle) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Experimental.
func (c *jsiiProxy_CfnWaitConditionHandle) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnWaitConditionHandle) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnWaitConditionHandle) Synthesize(session ISynthesisSession) {
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
func (c *jsiiProxy_CfnWaitConditionHandle) ToString() *string {
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
func (c *jsiiProxy_CfnWaitConditionHandle) Validate() *[]*string {
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
func (c *jsiiProxy_CfnWaitConditionHandle) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::CloudFormation::WaitCondition`.
type CfnWaitConditionProps struct {
	// `AWS::CloudFormation::WaitCondition.Count`.
	Count *float64 `json:"count"`
	// `AWS::CloudFormation::WaitCondition.Handle`.
	Handle *string `json:"handle"`
	// `AWS::CloudFormation::WaitCondition.Timeout`.
	Timeout *string `json:"timeout"`
}

// A set of constructs to be used as a dependable.
//
// This class can be used when a set of constructs which are disjoint in the
// construct tree needs to be combined to be used as a single dependable.
// Experimental.
type ConcreteDependable interface {
	IDependable
	Add(construct IConstruct)
}

// The jsii proxy struct for ConcreteDependable
type jsiiProxy_ConcreteDependable struct {
	jsiiProxy_IDependable
}

// Experimental.
func NewConcreteDependable() ConcreteDependable {
	_init_.Initialize()

	j := jsiiProxy_ConcreteDependable{}

	_jsii_.Create(
		"monocdk.ConcreteDependable",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewConcreteDependable_Override(c ConcreteDependable) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.ConcreteDependable",
		nil, // no parameters
		c,
	)
}

// Add a construct to the dependency roots.
// Experimental.
func (c *jsiiProxy_ConcreteDependable) Add(construct IConstruct) {
	_jsii_.InvokeVoid(
		c,
		"add",
		[]interface{}{construct},
	)
}

// Represents the building block of the construct graph.
//
// All constructs besides the root construct must be created within the scope of
// another construct.
// Experimental.
type Construct interface {
	constructs.Construct
	IConstruct
	Node() ConstructNode
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Construct
type jsiiProxy_Construct struct {
	internal.Type__constructsConstruct
	jsiiProxy_IConstruct
}

func (j *jsiiProxy_Construct) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


// Experimental.
func NewConstruct(scope constructs.Construct, id *string) Construct {
	_init_.Initialize()

	j := jsiiProxy_Construct{}

	_jsii_.Create(
		"monocdk.Construct",
		[]interface{}{scope, id},
		&j,
	)

	return &j
}

// Experimental.
func NewConstruct_Override(c Construct, scope constructs.Construct, id *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.Construct",
		[]interface{}{scope, id},
		c,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Construct_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Construct",
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
func (c *jsiiProxy_Construct) OnPrepare() {
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
func (c *jsiiProxy_Construct) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_Construct) OnValidate() *[]*string {
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
func (c *jsiiProxy_Construct) Prepare() {
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
func (c *jsiiProxy_Construct) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_Construct) ToString() *string {
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
func (c *jsiiProxy_Construct) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents the construct node in the scope tree.
// Experimental.
type ConstructNode interface {
	Addr() *string
	Children() *[]IConstruct
	DefaultChild() IConstruct
	SetDefaultChild(val IConstruct)
	Dependencies() *[]*Dependency
	Id() *string
	Locked() *bool
	Metadata() *[]*cxapi.MetadataEntry
	MetadataEntry() *[]*constructs.MetadataEntry
	Path() *string
	Root() IConstruct
	Scope() IConstruct
	Scopes() *[]IConstruct
	UniqueId() *string
	AddDependency(dependencies ...IDependable)
	AddError(message *string)
	AddInfo(message *string)
	AddMetadata(type_ *string, data interface{}, fromFunction interface{})
	AddValidation(validation constructs.IValidation)
	AddWarning(message *string)
	ApplyAspect(aspect IAspect)
	FindAll(order ConstructOrder) *[]IConstruct
	FindChild(id *string) IConstruct
	SetContext(key *string, value interface{})
	TryFindChild(id *string) IConstruct
	TryGetContext(key *string) interface{}
	TryRemoveChild(childName *string) *bool
}

// The jsii proxy struct for ConstructNode
type jsiiProxy_ConstructNode struct {
	_ byte // padding
}

func (j *jsiiProxy_ConstructNode) Addr() *string {
	var returns *string
	_jsii_.Get(
		j,
		"addr",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Children() *[]IConstruct {
	var returns *[]IConstruct
	_jsii_.Get(
		j,
		"children",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) DefaultChild() IConstruct {
	var returns IConstruct
	_jsii_.Get(
		j,
		"defaultChild",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Dependencies() *[]*Dependency {
	var returns *[]*Dependency
	_jsii_.Get(
		j,
		"dependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Locked() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"locked",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Metadata() *[]*cxapi.MetadataEntry {
	var returns *[]*cxapi.MetadataEntry
	_jsii_.Get(
		j,
		"metadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) MetadataEntry() *[]*constructs.MetadataEntry {
	var returns *[]*constructs.MetadataEntry
	_jsii_.Get(
		j,
		"metadataEntry",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Root() IConstruct {
	var returns IConstruct
	_jsii_.Get(
		j,
		"root",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Scope() IConstruct {
	var returns IConstruct
	_jsii_.Get(
		j,
		"scope",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) Scopes() *[]IConstruct {
	var returns *[]IConstruct
	_jsii_.Get(
		j,
		"scopes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ConstructNode) UniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"uniqueId",
		&returns,
	)
	return returns
}


// Experimental.
func NewConstructNode(host Construct, scope IConstruct, id *string) ConstructNode {
	_init_.Initialize()

	j := jsiiProxy_ConstructNode{}

	_jsii_.Create(
		"monocdk.ConstructNode",
		[]interface{}{host, scope, id},
		&j,
	)

	return &j
}

// Experimental.
func NewConstructNode_Override(c ConstructNode, host Construct, scope IConstruct, id *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.ConstructNode",
		[]interface{}{host, scope, id},
		c,
	)
}

func (j *jsiiProxy_ConstructNode) SetDefaultChild(val IConstruct) {
	_jsii_.Set(
		j,
		"defaultChild",
		val,
	)
}

// Invokes "prepare" on all constructs (depth-first, post-order) in the tree under `node`.
// Deprecated: Use `app.synth()` instead
func ConstructNode_Prepare(node ConstructNode) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.ConstructNode",
		"prepare",
		[]interface{}{node},
	)
}

// Synthesizes a CloudAssembly from a construct tree.
// Deprecated: Use `app.synth()` or `stage.synth()` instead
func ConstructNode_Synth(node ConstructNode, options *SynthesisOptions) cxapi.CloudAssembly {
	_init_.Initialize()

	var returns cxapi.CloudAssembly

	_jsii_.StaticInvoke(
		"monocdk.ConstructNode",
		"synth",
		[]interface{}{node, options},
		&returns,
	)

	return returns
}

// Invokes "validate" on all constructs in the tree (depth-first, pre-order) and returns the list of all errors.
//
// An empty list indicates that there are no errors.
// Experimental.
func ConstructNode_Validate(node ConstructNode) *[]*ValidationError {
	_init_.Initialize()

	var returns *[]*ValidationError

	_jsii_.StaticInvoke(
		"monocdk.ConstructNode",
		"validate",
		[]interface{}{node},
		&returns,
	)

	return returns
}

func ConstructNode_PATH_SEP() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.ConstructNode",
		"PATH_SEP",
		&returns,
	)
	return returns
}

// Add an ordering dependency on another Construct.
//
// All constructs in the dependency's scope will be deployed before any
// construct in this construct's scope.
// Experimental.
func (c *jsiiProxy_ConstructNode) AddDependency(dependencies ...IDependable) {
	args := []interface{}{}
	for _, a := range dependencies {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addDependency",
		args,
	)
}

// DEPRECATED: Adds an { "error": <message> } metadata entry to this construct.
//
// The toolkit will fail synthesis when errors are reported.
// Deprecated: use `Annotations.of(construct).addError()`
func (c *jsiiProxy_ConstructNode) AddError(message *string) {
	_jsii_.InvokeVoid(
		c,
		"addError",
		[]interface{}{message},
	)
}

// DEPRECATED: Adds a { "info": <message> } metadata entry to this construct.
//
// The toolkit will display the info message when apps are synthesized.
// Deprecated: use `Annotations.of(construct).addInfo()`
func (c *jsiiProxy_ConstructNode) AddInfo(message *string) {
	_jsii_.InvokeVoid(
		c,
		"addInfo",
		[]interface{}{message},
	)
}

// Adds a metadata entry to this construct.
//
// Entries are arbitrary values and will also include a stack trace to allow tracing back to
// the code location for when the entry was added. It can be used, for example, to include source
// mapping in CloudFormation templates to improve diagnostics.
// Experimental.
func (c *jsiiProxy_ConstructNode) AddMetadata(type_ *string, data interface{}, fromFunction interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addMetadata",
		[]interface{}{type_, data, fromFunction},
	)
}

// Add a validator to this construct Node.
// Experimental.
func (c *jsiiProxy_ConstructNode) AddValidation(validation constructs.IValidation) {
	_jsii_.InvokeVoid(
		c,
		"addValidation",
		[]interface{}{validation},
	)
}

// DEPRECATED: Adds a { "warning": <message> } metadata entry to this construct.
//
// The toolkit will display the warning when an app is synthesized, or fail
// if run in --strict mode.
// Deprecated: use `Annotations.of(construct).addWarning()`
func (c *jsiiProxy_ConstructNode) AddWarning(message *string) {
	_jsii_.InvokeVoid(
		c,
		"addWarning",
		[]interface{}{message},
	)
}

// DEPRECATED: Applies the aspect to this Constructs node.
// Deprecated: This API is going to be removed in the next major version of
// the AWS CDK. Please use `Aspects.of(scope).add()` instead.
func (c *jsiiProxy_ConstructNode) ApplyAspect(aspect IAspect) {
	_jsii_.InvokeVoid(
		c,
		"applyAspect",
		[]interface{}{aspect},
	)
}

// Return this construct and all of its children in the given order.
// Experimental.
func (c *jsiiProxy_ConstructNode) FindAll(order ConstructOrder) *[]IConstruct {
	var returns *[]IConstruct

	_jsii_.Invoke(
		c,
		"findAll",
		[]interface{}{order},
		&returns,
	)

	return returns
}

// Return a direct child by id.
//
// Throws an error if the child is not found.
//
// Returns: Child with the given id.
// Experimental.
func (c *jsiiProxy_ConstructNode) FindChild(id *string) IConstruct {
	var returns IConstruct

	_jsii_.Invoke(
		c,
		"findChild",
		[]interface{}{id},
		&returns,
	)

	return returns
}

// This can be used to set contextual values.
//
// Context must be set before any children are added, since children may consult context info during construction.
// If the key already exists, it will be overridden.
// Experimental.
func (c *jsiiProxy_ConstructNode) SetContext(key *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"setContext",
		[]interface{}{key, value},
	)
}

// Return a direct child by id, or undefined.
//
// Returns: the child if found, or undefined
// Experimental.
func (c *jsiiProxy_ConstructNode) TryFindChild(id *string) IConstruct {
	var returns IConstruct

	_jsii_.Invoke(
		c,
		"tryFindChild",
		[]interface{}{id},
		&returns,
	)

	return returns
}

// Retrieves a value from tree context.
//
// Context is usually initialized at the root, but can be overridden at any point in the tree.
//
// Returns: The context value or `undefined` if there is no context value for the key.
// Experimental.
func (c *jsiiProxy_ConstructNode) TryGetContext(key *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		c,
		"tryGetContext",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Remove the child with the given name, if present.
//
// Returns: Whether a child with the given name was deleted.
// Experimental.
func (c *jsiiProxy_ConstructNode) TryRemoveChild(childName *string) *bool {
	var returns *bool

	_jsii_.Invoke(
		c,
		"tryRemoveChild",
		[]interface{}{childName},
		&returns,
	)

	return returns
}

// In what order to return constructs.
// Experimental.
type ConstructOrder string

const (
	ConstructOrder_PREORDER ConstructOrder = "PREORDER"
	ConstructOrder_POSTORDER ConstructOrder = "POSTORDER"
)

// Base class for the model side of context providers.
//
// Instances of this class communicate with context provider plugins in the 'cdk
// toolkit' via context variables (input), outputting specialized queries for
// more context variables (output).
//
// ContextProvider needs access to a Construct to hook into the context mechanism.
// Experimental.
type ContextProvider interface {
}

// The jsii proxy struct for ContextProvider
type jsiiProxy_ContextProvider struct {
	_ byte // padding
}

// Returns: the context key or undefined if a key cannot be rendered (due to tokens used in any of the props)
// Experimental.
func ContextProvider_GetKey(scope constructs.Construct, options *GetContextKeyOptions) *GetContextKeyResult {
	_init_.Initialize()

	var returns *GetContextKeyResult

	_jsii_.StaticInvoke(
		"monocdk.ContextProvider",
		"getKey",
		[]interface{}{scope, options},
		&returns,
	)

	return returns
}

// Experimental.
func ContextProvider_GetValue(scope constructs.Construct, options *GetContextValueOptions) *GetContextValueResult {
	_init_.Initialize()

	var returns *GetContextValueResult

	_jsii_.StaticInvoke(
		"monocdk.ContextProvider",
		"getValue",
		[]interface{}{scope, options},
		&returns,
	)

	return returns
}

// Options applied when copying directories.
// Experimental.
type CopyOptions struct {
	// Glob patterns to exclude from the copy.
	// Experimental.
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Experimental.
	Follow SymlinkFollowMode `json:"follow"`
	// The ignore behavior to use for exclude patterns.
	// Experimental.
	IgnoreMode IgnoreMode `json:"ignoreMode"`
}

// Custom resource that is implemented using a Lambda.
//
// As a custom resource author, you should be publishing a subclass of this class
// that hides the choice of provider, and accepts a strongly-typed properties
// object with the properties your provider accepts.
// Experimental.
type CustomResource interface {
	Resource
	Env() *ResourceEnvironment
	Node() ConstructNode
	PhysicalName() *string
	Ref() *string
	Stack() Stack
	ApplyRemovalPolicy(policy RemovalPolicy)
	GeneratePhysicalName() *string
	GetAtt(attributeName *string) Reference
	GetAttString(attributeName *string) *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CustomResource
type jsiiProxy_CustomResource struct {
	jsiiProxy_Resource
}

func (j *jsiiProxy_CustomResource) Env() *ResourceEnvironment {
	var returns *ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomResource) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomResource) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomResource) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomResource) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewCustomResource(scope constructs.Construct, id *string, props *CustomResourceProps) CustomResource {
	_init_.Initialize()

	j := jsiiProxy_CustomResource{}

	_jsii_.Create(
		"monocdk.CustomResource",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCustomResource_Override(c CustomResource, scope constructs.Construct, id *string, props *CustomResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CustomResource",
		[]interface{}{scope, id, props},
		c,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func CustomResource_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CustomResource",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func CustomResource_IsResource(construct IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CustomResource",
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
func (c *jsiiProxy_CustomResource) ApplyRemovalPolicy(policy RemovalPolicy) {
	_jsii_.InvokeVoid(
		c,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (c *jsiiProxy_CustomResource) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns the value of an attribute of the custom resource of an arbitrary type.
//
// Attributes are returned from the custom resource provider through the
// `Data` map where the key is the attribute name.
//
// Returns: a token for `Fn::GetAtt`. Use `Token.asXxx` to encode the returned `Reference` as a specific type or
// use the convenience `getAttString` for string attributes.
// Experimental.
func (c *jsiiProxy_CustomResource) GetAtt(attributeName *string) Reference {
	var returns Reference

	_jsii_.Invoke(
		c,
		"getAtt",
		[]interface{}{attributeName},
		&returns,
	)

	return returns
}

// Returns the value of an attribute of the custom resource of type string.
//
// Attributes are returned from the custom resource provider through the
// `Data` map where the key is the attribute name.
//
// Returns: a token for `Fn::GetAtt` encoded as a string.
// Experimental.
func (c *jsiiProxy_CustomResource) GetAttString(attributeName *string) *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"getAttString",
		[]interface{}{attributeName},
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
func (c *jsiiProxy_CustomResource) GetResourceArnAttribute(arnAttr *string, arnComponents *ArnComponents) *string {
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
func (c *jsiiProxy_CustomResource) GetResourceNameAttribute(nameAttr *string) *string {
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
func (c *jsiiProxy_CustomResource) OnPrepare() {
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
func (c *jsiiProxy_CustomResource) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CustomResource) OnValidate() *[]*string {
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
func (c *jsiiProxy_CustomResource) Prepare() {
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
func (c *jsiiProxy_CustomResource) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CustomResource) ToString() *string {
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
func (c *jsiiProxy_CustomResource) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to provide a Lambda-backed custom resource.
// Experimental.
type CustomResourceProps struct {
	// The ARN of the provider which implements this custom resource type.
	//
	// You can implement a provider by listening to raw AWS CloudFormation events
	// and specify the ARN of an SNS topic (`topic.topicArn`) or the ARN of an AWS
	// Lambda function (`lambda.functionArn`) or use the CDK's custom [resource
	// provider framework] which makes it easier to implement robust providers.
	//
	// [resource provider framework]:
	// https://docs.aws.amazon.com/cdk/api/latest/docs/custom-resources-readme.html
	//
	// Provider framework:
	//
	// ```ts
	// // use the provider framework from aws-cdk/custom-resources:
	// const provider = new customresources.Provider(this, 'ResourceProvider', {
	//    onEventHandler,
	//    isCompleteHandler, // optional
	// });
	//
	// new CustomResource(this, 'MyResource', {
	//    serviceToken: provider.serviceToken,
	// });
	// ```
	//
	// AWS Lambda function:
	//
	// ```ts
	// // invoke an AWS Lambda function when a lifecycle event occurs:
	// new CustomResource(this, 'MyResource', {
	//    serviceToken: myFunction.functionArn,
	// });
	// ```
	//
	// SNS topic:
	//
	// ```ts
	// // publish lifecycle events to an SNS topic:
	// new CustomResource(this, 'MyResource', {
	//    serviceToken: myTopic.topicArn,
	// });
	// ```
	// Experimental.
	ServiceToken *string `json:"serviceToken"`
	// Convert all property keys to pascal case.
	// Experimental.
	PascalCaseProperties *bool `json:"pascalCaseProperties"`
	// Properties to pass to the Lambda.
	// Experimental.
	Properties *map[string]interface{} `json:"properties"`
	// The policy to apply when this resource is removed from the application.
	// Experimental.
	RemovalPolicy RemovalPolicy `json:"removalPolicy"`
	// For custom resources, you can specify AWS::CloudFormation::CustomResource (the default) as the resource type, or you can specify your own resource type name.
	//
	// For example, you can use "Custom::MyCustomResourceTypeName".
	//
	// Custom resource type names must begin with "Custom::" and can include
	// alphanumeric characters and the following characters: _@-. You can specify
	// a custom resource type name up to a maximum length of 60 characters. You
	// cannot change the type during an update.
	//
	// Using your own resource type names helps you quickly differentiate the
	// types of custom resources in your stack. For example, if you had two custom
	// resources that conduct two different ping tests, you could name their type
	// as Custom::PingTester to make them easily identifiable as ping testers
	// (instead of using AWS::CloudFormation::CustomResource).
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cfn-customresource.html#aws-cfn-resource-type-name
	//
	// Experimental.
	ResourceType *string `json:"resourceType"`
}

// An AWS-Lambda backed custom resource provider.
// Experimental.
type CustomResourceProvider interface {
	Construct
	Node() ConstructNode
	RoleArn() *string
	ServiceToken() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for CustomResourceProvider
type jsiiProxy_CustomResourceProvider struct {
	jsiiProxy_Construct
}

func (j *jsiiProxy_CustomResourceProvider) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomResourceProvider) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CustomResourceProvider) ServiceToken() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serviceToken",
		&returns,
	)
	return returns
}


// Experimental.
func NewCustomResourceProvider(scope constructs.Construct, id *string, props *CustomResourceProviderProps) CustomResourceProvider {
	_init_.Initialize()

	j := jsiiProxy_CustomResourceProvider{}

	_jsii_.Create(
		"monocdk.CustomResourceProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCustomResourceProvider_Override(c CustomResourceProvider, scope constructs.Construct, id *string, props *CustomResourceProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.CustomResourceProvider",
		[]interface{}{scope, id, props},
		c,
	)
}

// Returns a stack-level singleton ARN (service token) for the custom resource provider.
//
// Returns: the service token of the custom resource provider, which should be
// used when defining a `CustomResource`.
// Experimental.
func CustomResourceProvider_GetOrCreate(scope constructs.Construct, uniqueid *string, props *CustomResourceProviderProps) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.CustomResourceProvider",
		"getOrCreate",
		[]interface{}{scope, uniqueid, props},
		&returns,
	)

	return returns
}

// Returns a stack-level singleton for the custom resource provider.
//
// Returns: the service token of the custom resource provider, which should be
// used when defining a `CustomResource`.
// Experimental.
func CustomResourceProvider_GetOrCreateProvider(scope constructs.Construct, uniqueid *string, props *CustomResourceProviderProps) CustomResourceProvider {
	_init_.Initialize()

	var returns CustomResourceProvider

	_jsii_.StaticInvoke(
		"monocdk.CustomResourceProvider",
		"getOrCreateProvider",
		[]interface{}{scope, uniqueid, props},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CustomResourceProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.CustomResourceProvider",
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
func (c *jsiiProxy_CustomResourceProvider) OnPrepare() {
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
func (c *jsiiProxy_CustomResourceProvider) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CustomResourceProvider) OnValidate() *[]*string {
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
func (c *jsiiProxy_CustomResourceProvider) Prepare() {
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
func (c *jsiiProxy_CustomResourceProvider) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CustomResourceProvider) ToString() *string {
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
func (c *jsiiProxy_CustomResourceProvider) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Initialization properties for `CustomResourceProvider`.
// Experimental.
type CustomResourceProviderProps struct {
	// A local file system directory with the provider's code.
	//
	// The code will be
	// bundled into a zip asset and wired to the provider's AWS Lambda function.
	// Experimental.
	CodeDirectory *string `json:"codeDirectory"`
	// The AWS Lambda runtime and version to use for the provider.
	// Experimental.
	Runtime CustomResourceProviderRuntime `json:"runtime"`
	// A description of the function.
	// Experimental.
	Description *string `json:"description"`
	// Key-value pairs that are passed to Lambda as Environment.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// The amount of memory that your function has access to.
	//
	// Increasing the
	// function's memory also increases its CPU allocation.
	// Experimental.
	MemorySize Size `json:"memorySize"`
	// A set of IAM policy statements to include in the inline policy of the provider's lambda function.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	PolicyStatements *[]interface{} `json:"policyStatements"`
	// AWS Lambda timeout for the provider.
	// Experimental.
	Timeout Duration `json:"timeout"`
}

// The lambda runtime to use for the resource provider.
//
// This also indicates
// which language is used for the handler.
// Experimental.
type CustomResourceProviderRuntime string

const (
	CustomResourceProviderRuntime_NODEJS_12 CustomResourceProviderRuntime = "NODEJS_12"
	CustomResourceProviderRuntime_NODEJS_14_X CustomResourceProviderRuntime = "NODEJS_14_X"
)

// Uses conventionally named roles and reify asset storage locations.
//
// This synthesizer is the only StackSynthesizer that generates
// an asset manifest, and is required to deploy CDK applications using the
// `@aws-cdk/app-delivery` CI/CD library.
//
// Requires the environment to have been bootstrapped with Bootstrap Stack V2.
// Experimental.
type DefaultStackSynthesizer interface {
	StackSynthesizer
	CloudFormationExecutionRoleArn() *string
	DeployRoleArn() *string
	Stack() Stack
	AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation
	AddFileAsset(asset *FileAssetSource) *FileAssetLocation
	Bind(stack Stack)
	EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions)
	Synthesize(session ISynthesisSession)
	SynthesizeStackTemplate(stack Stack, session ISynthesisSession)
}

// The jsii proxy struct for DefaultStackSynthesizer
type jsiiProxy_DefaultStackSynthesizer struct {
	jsiiProxy_StackSynthesizer
}

func (j *jsiiProxy_DefaultStackSynthesizer) CloudFormationExecutionRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cloudFormationExecutionRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DefaultStackSynthesizer) DeployRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deployRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DefaultStackSynthesizer) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewDefaultStackSynthesizer(props *DefaultStackSynthesizerProps) DefaultStackSynthesizer {
	_init_.Initialize()

	j := jsiiProxy_DefaultStackSynthesizer{}

	_jsii_.Create(
		"monocdk.DefaultStackSynthesizer",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewDefaultStackSynthesizer_Override(d DefaultStackSynthesizer, props *DefaultStackSynthesizerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.DefaultStackSynthesizer",
		[]interface{}{props},
		d,
	)
}

func DefaultStackSynthesizer_DEFAULT_BOOTSTRAP_STACK_VERSION_SSM_PARAMETER() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_BOOTSTRAP_STACK_VERSION_SSM_PARAMETER",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_CLOUDFORMATION_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_CLOUDFORMATION_ROLE_ARN",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_DEPLOY_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_DEPLOY_ROLE_ARN",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_FILE_ASSET_KEY_ARN_EXPORT_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_FILE_ASSET_KEY_ARN_EXPORT_NAME",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_FILE_ASSET_PREFIX() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_FILE_ASSET_PREFIX",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_FILE_ASSET_PUBLISHING_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_FILE_ASSET_PUBLISHING_ROLE_ARN",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_FILE_ASSETS_BUCKET_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_FILE_ASSETS_BUCKET_NAME",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_IMAGE_ASSET_PUBLISHING_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_IMAGE_ASSET_PUBLISHING_ROLE_ARN",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_IMAGE_ASSETS_REPOSITORY_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_IMAGE_ASSETS_REPOSITORY_NAME",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_LOOKUP_ROLE_ARN() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_LOOKUP_ROLE_ARN",
		&returns,
	)
	return returns
}

func DefaultStackSynthesizer_DEFAULT_QUALIFIER() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.DefaultStackSynthesizer",
		"DEFAULT_QUALIFIER",
		&returns,
	)
	return returns
}

// Register a Docker Image Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (d *jsiiProxy_DefaultStackSynthesizer) AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		d,
		"addDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Register a File Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (d *jsiiProxy_DefaultStackSynthesizer) AddFileAsset(asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		d,
		"addFileAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Bind to the stack this environment is going to be used on.
//
// Must be called before any of the other methods are called.
// Experimental.
func (d *jsiiProxy_DefaultStackSynthesizer) Bind(stack Stack) {
	_jsii_.InvokeVoid(
		d,
		"bind",
		[]interface{}{stack},
	)
}

// Write the stack artifact to the session.
//
// Use default settings to add a CloudFormationStackArtifact artifact to
// the given synthesis session.
// Experimental.
func (d *jsiiProxy_DefaultStackSynthesizer) EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions) {
	_jsii_.InvokeVoid(
		d,
		"emitStackArtifact",
		[]interface{}{stack, session, options},
	)
}

// Synthesize the associated stack to the session.
// Experimental.
func (d *jsiiProxy_DefaultStackSynthesizer) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
		"synthesize",
		[]interface{}{session},
	)
}

// Have the stack write out its template.
// Experimental.
func (d *jsiiProxy_DefaultStackSynthesizer) SynthesizeStackTemplate(stack Stack, session ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
		"synthesizeStackTemplate",
		[]interface{}{stack, session},
	)
}

// Configuration properties for DefaultStackSynthesizer.
// Experimental.
type DefaultStackSynthesizerProps struct {
	// Bootstrap stack version SSM parameter.
	//
	// The placeholder `${Qualifier}` will be replaced with the value of qualifier.
	// Experimental.
	BootstrapStackVersionSsmParameter *string `json:"bootstrapStackVersionSsmParameter"`
	// bucketPrefix to use while storing S3 Assets.
	// Experimental.
	BucketPrefix *string `json:"bucketPrefix"`
	// The role CloudFormation will assume when deploying the Stack.
	//
	// You must supply this if you have given a non-standard name to the execution role.
	//
	// The placeholders `${Qualifier}`, `${AWS::AccountId}` and `${AWS::Region}` will
	// be replaced with the values of qualifier and the stack's account and region,
	// respectively.
	// Experimental.
	CloudFormationExecutionRole *string `json:"cloudFormationExecutionRole"`
	// The role to assume to initiate a deployment in this environment.
	//
	// You must supply this if you have given a non-standard name to the publishing role.
	//
	// The placeholders `${Qualifier}`, `${AWS::AccountId}` and `${AWS::Region}` will
	// be replaced with the values of qualifier and the stack's account and region,
	// respectively.
	// Experimental.
	DeployRoleArn *string `json:"deployRoleArn"`
	// Name of the CloudFormation Export with the asset key name.
	//
	// You must supply this if you have given a non-standard name to the KMS key export
	//
	// The placeholders `${Qualifier}`, `${AWS::AccountId}` and `${AWS::Region}` will
	// be replaced with the values of qualifier and the stack's account and region,
	// respectively.
	// Deprecated: This property is not used anymore
	FileAssetKeyArnExportName *string `json:"fileAssetKeyArnExportName"`
	// External ID to use when assuming role for file asset publishing.
	// Experimental.
	FileAssetPublishingExternalId *string `json:"fileAssetPublishingExternalId"`
	// The role to use to publish file assets to the S3 bucket in this environment.
	//
	// You must supply this if you have given a non-standard name to the publishing role.
	//
	// The placeholders `${Qualifier}`, `${AWS::AccountId}` and `${AWS::Region}` will
	// be replaced with the values of qualifier and the stack's account and region,
	// respectively.
	// Experimental.
	FileAssetPublishingRoleArn *string `json:"fileAssetPublishingRoleArn"`
	// Name of the S3 bucket to hold file assets.
	//
	// You must supply this if you have given a non-standard name to the staging bucket.
	//
	// The placeholders `${Qualifier}`, `${AWS::AccountId}` and `${AWS::Region}` will
	// be replaced with the values of qualifier and the stack's account and region,
	// respectively.
	// Experimental.
	FileAssetsBucketName *string `json:"fileAssetsBucketName"`
	// Whether to add a Rule to the stack template verifying the bootstrap stack version.
	//
	// This generally should be left set to `true`, unless you explicitly
	// want to be able to deploy to an unbootstrapped environment.
	// Experimental.
	GenerateBootstrapVersionRule *bool `json:"generateBootstrapVersionRule"`
	// External ID to use when assuming role for image asset publishing.
	// Experimental.
	ImageAssetPublishingExternalId *string `json:"imageAssetPublishingExternalId"`
	// The role to use to publish image assets to the ECR repository in this environment.
	//
	// You must supply this if you have given a non-standard name to the publishing role.
	//
	// The placeholders `${Qualifier}`, `${AWS::AccountId}` and `${AWS::Region}` will
	// be replaced with the values of qualifier and the stack's account and region,
	// respectively.
	// Experimental.
	ImageAssetPublishingRoleArn *string `json:"imageAssetPublishingRoleArn"`
	// Name of the ECR repository to hold Docker Image assets.
	//
	// You must supply this if you have given a non-standard name to the ECR repository.
	//
	// The placeholders `${Qualifier}`, `${AWS::AccountId}` and `${AWS::Region}` will
	// be replaced with the values of qualifier and the stack's account and region,
	// respectively.
	// Experimental.
	ImageAssetsRepositoryName *string `json:"imageAssetsRepositoryName"`
	// The role to use to look up values from the target AWS account during synthesis.
	// Experimental.
	LookupRoleArn *string `json:"lookupRoleArn"`
	// Qualifier to disambiguate multiple environments in the same account.
	//
	// You can use this and leave the other naming properties empty if you have deployed
	// the bootstrap environment with standard names but only differnet qualifiers.
	// Experimental.
	Qualifier *string `json:"qualifier"`
}

// Default resolver implementation.
// Experimental.
type DefaultTokenResolver interface {
	ITokenResolver
	ResolveList(xs *[]*string, context IResolveContext) interface{}
	ResolveString(fragments TokenizedStringFragments, context IResolveContext) interface{}
	ResolveToken(t IResolvable, context IResolveContext, postProcessor IPostProcessor) interface{}
}

// The jsii proxy struct for DefaultTokenResolver
type jsiiProxy_DefaultTokenResolver struct {
	jsiiProxy_ITokenResolver
}

// Experimental.
func NewDefaultTokenResolver(concat IFragmentConcatenator) DefaultTokenResolver {
	_init_.Initialize()

	j := jsiiProxy_DefaultTokenResolver{}

	_jsii_.Create(
		"monocdk.DefaultTokenResolver",
		[]interface{}{concat},
		&j,
	)

	return &j
}

// Experimental.
func NewDefaultTokenResolver_Override(d DefaultTokenResolver, concat IFragmentConcatenator) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.DefaultTokenResolver",
		[]interface{}{concat},
		d,
	)
}

// Resolve a tokenized list.
// Experimental.
func (d *jsiiProxy_DefaultTokenResolver) ResolveList(xs *[]*string, context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"resolveList",
		[]interface{}{xs, context},
		&returns,
	)

	return returns
}

// Resolve string fragments to Tokens.
// Experimental.
func (d *jsiiProxy_DefaultTokenResolver) ResolveString(fragments TokenizedStringFragments, context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"resolveString",
		[]interface{}{fragments, context},
		&returns,
	)

	return returns
}

// Default Token resolution.
//
// Resolve the Token, recurse into whatever it returns,
// then finally post-process it.
// Experimental.
func (d *jsiiProxy_DefaultTokenResolver) ResolveToken(t IResolvable, context IResolveContext, postProcessor IPostProcessor) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"resolveToken",
		[]interface{}{t, context, postProcessor},
		&returns,
	)

	return returns
}

// Trait for IDependable.
//
// Traits are interfaces that are privately implemented by objects. Instead of
// showing up in the public interface of a class, they need to be queried
// explicitly. This is used to implement certain framework features that are
// not intended to be used by Construct consumers, and so should be hidden
// from accidental use.
//
// TODO: EXAMPLE
//
// Experimental.
type DependableTrait interface {
	DependencyRoots() *[]IConstruct
}

// The jsii proxy struct for DependableTrait
type jsiiProxy_DependableTrait struct {
	_ byte // padding
}

func (j *jsiiProxy_DependableTrait) DependencyRoots() *[]IConstruct {
	var returns *[]IConstruct
	_jsii_.Get(
		j,
		"dependencyRoots",
		&returns,
	)
	return returns
}


// Experimental.
func NewDependableTrait_Override(d DependableTrait) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.DependableTrait",
		nil, // no parameters
		d,
	)
}

// Return the matching DependableTrait for the given class instance.
// Experimental.
func DependableTrait_Get(instance IDependable) DependableTrait {
	_init_.Initialize()

	var returns DependableTrait

	_jsii_.StaticInvoke(
		"monocdk.DependableTrait",
		"get",
		[]interface{}{instance},
		&returns,
	)

	return returns
}

// Register `instance` to have the given DependableTrait.
//
// Should be called in the class constructor.
// Experimental.
func DependableTrait_Implement(instance IDependable, trait DependableTrait) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.DependableTrait",
		"implement",
		[]interface{}{instance, trait},
	)
}

// A single dependency.
// Experimental.
type Dependency struct {
	// Source the dependency.
	// Experimental.
	Source IConstruct `json:"source"`
	// Target of the dependency.
	// Experimental.
	Target IConstruct `json:"target"`
}

// Docker build options.
// Experimental.
type DockerBuildOptions struct {
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
}

// Ignores file paths based on the [`.dockerignore specification`](https://docs.docker.com/engine/reference/builder/#dockerignore-file).
// Experimental.
type DockerIgnoreStrategy interface {
	IgnoreStrategy
	Add(pattern *string)
	Ignores(absoluteFilePath *string) *bool
}

// The jsii proxy struct for DockerIgnoreStrategy
type jsiiProxy_DockerIgnoreStrategy struct {
	jsiiProxy_IgnoreStrategy
}

// Experimental.
func NewDockerIgnoreStrategy(absoluteRootPath *string, patterns *[]*string) DockerIgnoreStrategy {
	_init_.Initialize()

	j := jsiiProxy_DockerIgnoreStrategy{}

	_jsii_.Create(
		"monocdk.DockerIgnoreStrategy",
		[]interface{}{absoluteRootPath, patterns},
		&j,
	)

	return &j
}

// Experimental.
func NewDockerIgnoreStrategy_Override(d DockerIgnoreStrategy, absoluteRootPath *string, patterns *[]*string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.DockerIgnoreStrategy",
		[]interface{}{absoluteRootPath, patterns},
		d,
	)
}

// Ignores file paths based on the [`.dockerignore specification`](https://docs.docker.com/engine/reference/builder/#dockerignore-file).
//
// Returns: `DockerIgnorePattern` associated with the given patterns.
// Experimental.
func DockerIgnoreStrategy_Docker(absoluteRootPath *string, patterns *[]*string) DockerIgnoreStrategy {
	_init_.Initialize()

	var returns DockerIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.DockerIgnoreStrategy",
		"docker",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Creates an IgnoreStrategy based on the `ignoreMode` and `exclude` in a `CopyOptions`.
//
// Returns: `IgnoreStrategy` based on the `CopyOptions`
// Experimental.
func DockerIgnoreStrategy_FromCopyOptions(options *CopyOptions, absoluteRootPath *string) IgnoreStrategy {
	_init_.Initialize()

	var returns IgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.DockerIgnoreStrategy",
		"fromCopyOptions",
		[]interface{}{options, absoluteRootPath},
		&returns,
	)

	return returns
}

// Ignores file paths based on the [`.gitignore specification`](https://git-scm.com/docs/gitignore).
//
// Returns: `GitIgnorePattern` associated with the given patterns.
// Experimental.
func DockerIgnoreStrategy_Git(absoluteRootPath *string, patterns *[]*string) GitIgnoreStrategy {
	_init_.Initialize()

	var returns GitIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.DockerIgnoreStrategy",
		"git",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Ignores file paths based on simple glob patterns.
//
// Returns: `GlobIgnorePattern` associated with the given patterns.
// Experimental.
func DockerIgnoreStrategy_Glob(absoluteRootPath *string, patterns *[]*string) GlobIgnoreStrategy {
	_init_.Initialize()

	var returns GlobIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.DockerIgnoreStrategy",
		"glob",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Adds another pattern.
// Experimental.
func (d *jsiiProxy_DockerIgnoreStrategy) Add(pattern *string) {
	_jsii_.InvokeVoid(
		d,
		"add",
		[]interface{}{pattern},
	)
}

// Determines whether a given file path should be ignored or not.
//
// Returns: `true` if the file should be ignored
// Experimental.
func (d *jsiiProxy_DockerIgnoreStrategy) Ignores(absoluteFilePath *string) *bool {
	var returns *bool

	_jsii_.Invoke(
		d,
		"ignores",
		[]interface{}{absoluteFilePath},
		&returns,
	)

	return returns
}

// A Docker image.
// Experimental.
type DockerImage interface {
	BundlingDockerImage
	Image() *string
	Cp(imagePath *string, outputPath *string) *string
	Run(options *DockerRunOptions)
	ToJSON() *string
}

// The jsii proxy struct for DockerImage
type jsiiProxy_DockerImage struct {
	jsiiProxy_BundlingDockerImage
}

func (j *jsiiProxy_DockerImage) Image() *string {
	var returns *string
	_jsii_.Get(
		j,
		"image",
		&returns,
	)
	return returns
}


// Experimental.
func NewDockerImage(image *string, _imageHash *string) DockerImage {
	_init_.Initialize()

	j := jsiiProxy_DockerImage{}

	_jsii_.Create(
		"monocdk.DockerImage",
		[]interface{}{image, _imageHash},
		&j,
	)

	return &j
}

// Experimental.
func NewDockerImage_Override(d DockerImage, image *string, _imageHash *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.DockerImage",
		[]interface{}{image, _imageHash},
		d,
	)
}

// Reference an image that's built directly from sources on disk.
// Deprecated: use DockerImage.fromBuild()
func DockerImage_FromAsset(path *string, options *DockerBuildOptions) BundlingDockerImage {
	_init_.Initialize()

	var returns BundlingDockerImage

	_jsii_.StaticInvoke(
		"monocdk.DockerImage",
		"fromAsset",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Builds a Docker image.
// Experimental.
func DockerImage_FromBuild(path *string, options *DockerBuildOptions) DockerImage {
	_init_.Initialize()

	var returns DockerImage

	_jsii_.StaticInvoke(
		"monocdk.DockerImage",
		"fromBuild",
		[]interface{}{path, options},
		&returns,
	)

	return returns
}

// Reference an image on DockerHub or another online registry.
// Experimental.
func DockerImage_FromRegistry(image *string) DockerImage {
	_init_.Initialize()

	var returns DockerImage

	_jsii_.StaticInvoke(
		"monocdk.DockerImage",
		"fromRegistry",
		[]interface{}{image},
		&returns,
	)

	return returns
}

// Copies a file or directory out of the Docker image to the local filesystem.
//
// If `outputPath` is omitted the destination path is a temporary directory.
//
// Returns: the destination path
// Experimental.
func (d *jsiiProxy_DockerImage) Cp(imagePath *string, outputPath *string) *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"cp",
		[]interface{}{imagePath, outputPath},
		&returns,
	)

	return returns
}

// Runs a Docker image.
// Experimental.
func (d *jsiiProxy_DockerImage) Run(options *DockerRunOptions) {
	_jsii_.InvokeVoid(
		d,
		"run",
		[]interface{}{options},
	)
}

// Provides a stable representation of this image for JSON serialization.
//
// Returns: The overridden image name if set or image hash name in that order
// Experimental.
func (d *jsiiProxy_DockerImage) ToJSON() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The location of the published docker image.
//
// This is where the image can be
// consumed at runtime.
// Experimental.
type DockerImageAssetLocation struct {
	// The URI of the image in Amazon ECR.
	// Experimental.
	ImageUri *string `json:"imageUri"`
	// The name of the ECR repository.
	// Experimental.
	RepositoryName *string `json:"repositoryName"`
}

// Experimental.
type DockerImageAssetSource struct {
	// The hash of the contents of the docker build context.
	//
	// This hash is used
	// throughout the system to identify this image and avoid duplicate work
	// in case the source did not change.
	//
	// NOTE: this means that if you wish to update your docker image, you
	// must make a modification to the source (e.g. add some metadata to your Dockerfile).
	// Experimental.
	SourceHash *string `json:"sourceHash"`
	// The directory where the Dockerfile is stored, must be relative to the cloud assembly root.
	// Experimental.
	DirectoryName *string `json:"directoryName"`
	// Build args to pass to the `docker build` command.
	//
	// Since Docker build arguments are resolved before deployment, keys and
	// values cannot refer to unresolved tokens (such as `lambda.functionArn` or
	// `queue.queueUrl`).
	//
	// Only allowed when `directoryName` is specified.
	// Experimental.
	DockerBuildArgs *map[string]*string `json:"dockerBuildArgs"`
	// Docker target to build to.
	//
	// Only allowed when `directoryName` is specified.
	// Experimental.
	DockerBuildTarget *string `json:"dockerBuildTarget"`
	// Path to the Dockerfile (relative to the directory).
	//
	// Only allowed when `directoryName` is specified.
	// Experimental.
	DockerFile *string `json:"dockerFile"`
	// An external command that will produce the packaged asset.
	//
	// The command should produce the name of a local Docker image on `stdout`.
	// Experimental.
	Executable *[]*string `json:"executable"`
	// ECR repository name.
	//
	// Specify this property if you need to statically address the image, e.g.
	// from a Kubernetes Pod. Note, this is only the repository name, without the
	// registry and the tag parts.
	// Deprecated: repository name should be specified at the environment-level and not at the image level
	RepositoryName *string `json:"repositoryName"`
}

// Docker run options.
// Experimental.
type DockerRunOptions struct {
	// The command to run in the container.
	// Experimental.
	Command *[]*string `json:"command"`
	// The entrypoint to run in the container.
	// Experimental.
	Entrypoint *[]*string `json:"entrypoint"`
	// The environment variables to pass to the container.
	// Experimental.
	Environment *map[string]*string `json:"environment"`
	// [Security configuration](https://docs.docker.com/engine/reference/run/#security-configuration) when running the docker container.
	// Experimental.
	SecurityOpt *string `json:"securityOpt"`
	// The user to use when running the container.
	// Experimental.
	User *string `json:"user"`
	// Docker volumes to mount.
	// Experimental.
	Volumes *[]*DockerVolume `json:"volumes"`
	// Working directory inside the container.
	// Experimental.
	WorkingDirectory *string `json:"workingDirectory"`
}

// A Docker volume.
// Experimental.
type DockerVolume struct {
	// The path where the file or directory is mounted in the container.
	// Experimental.
	ContainerPath *string `json:"containerPath"`
	// The path to the file or directory on the host machine.
	// Experimental.
	HostPath *string `json:"hostPath"`
	// Mount consistency.
	//
	// Only applicable for macOS
	// See: https://docs.docker.com/storage/bind-mounts/#configure-mount-consistency-for-macos
	//
	// Experimental.
	Consistency DockerVolumeConsistency `json:"consistency"`
}

// Supported Docker volume consistency types.
//
// Only valid on macOS due to the way file storage works on Mac
// Experimental.
type DockerVolumeConsistency string

const (
	DockerVolumeConsistency_CONSISTENT DockerVolumeConsistency = "CONSISTENT"
	DockerVolumeConsistency_DELEGATED DockerVolumeConsistency = "DELEGATED"
	DockerVolumeConsistency_CACHED DockerVolumeConsistency = "CACHED"
)

// Represents a length of time.
//
// The amount can be specified either as a literal value (e.g: `10`) which
// cannot be negative, or as an unresolved number token.
//
// When the amount is passed as a token, unit conversion is not possible.
// Experimental.
type Duration interface {
	FormatTokenToNumber() *string
	IsUnresolved() *bool
	Plus(rhs Duration) Duration
	ToDays(opts *TimeConversionOptions) *float64
	ToHours(opts *TimeConversionOptions) *float64
	ToHumanString() *string
	ToIsoString() *string
	ToISOString() *string
	ToMilliseconds(opts *TimeConversionOptions) *float64
	ToMinutes(opts *TimeConversionOptions) *float64
	ToSeconds(opts *TimeConversionOptions) *float64
	ToString() *string
	UnitLabel() *string
}

// The jsii proxy struct for Duration
type jsiiProxy_Duration struct {
	_ byte // padding
}

// Create a Duration representing an amount of days.
//
// Returns: a new `Duration` representing `amount` Days.
// Experimental.
func Duration_Days(amount *float64) Duration {
	_init_.Initialize()

	var returns Duration

	_jsii_.StaticInvoke(
		"monocdk.Duration",
		"days",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of hours.
//
// Returns: a new `Duration` representing `amount` Hours.
// Experimental.
func Duration_Hours(amount *float64) Duration {
	_init_.Initialize()

	var returns Duration

	_jsii_.StaticInvoke(
		"monocdk.Duration",
		"hours",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of milliseconds.
//
// Returns: a new `Duration` representing `amount` ms.
// Experimental.
func Duration_Millis(amount *float64) Duration {
	_init_.Initialize()

	var returns Duration

	_jsii_.StaticInvoke(
		"monocdk.Duration",
		"millis",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of minutes.
//
// Returns: a new `Duration` representing `amount` Minutes.
// Experimental.
func Duration_Minutes(amount *float64) Duration {
	_init_.Initialize()

	var returns Duration

	_jsii_.StaticInvoke(
		"monocdk.Duration",
		"minutes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Parse a period formatted according to the ISO 8601 standard.
//
// Returns: the parsed `Duration`.
// See: https://www.iso.org/fr/standard/70907.html
//
// Experimental.
func Duration_Parse(duration *string) Duration {
	_init_.Initialize()

	var returns Duration

	_jsii_.StaticInvoke(
		"monocdk.Duration",
		"parse",
		[]interface{}{duration},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of seconds.
//
// Returns: a new `Duration` representing `amount` Seconds.
// Experimental.
func Duration_Seconds(amount *float64) Duration {
	_init_.Initialize()

	var returns Duration

	_jsii_.StaticInvoke(
		"monocdk.Duration",
		"seconds",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Returns stringified number of duration.
// Experimental.
func (d *jsiiProxy_Duration) FormatTokenToNumber() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"formatTokenToNumber",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Checks if duration is a token or a resolvable object.
// Experimental.
func (d *jsiiProxy_Duration) IsUnresolved() *bool {
	var returns *bool

	_jsii_.Invoke(
		d,
		"isUnresolved",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Add two Durations together.
// Experimental.
func (d *jsiiProxy_Duration) Plus(rhs Duration) Duration {
	var returns Duration

	_jsii_.Invoke(
		d,
		"plus",
		[]interface{}{rhs},
		&returns,
	)

	return returns
}

// Return the total number of days in this Duration.
//
// Returns: the value of this `Duration` expressed in Days.
// Experimental.
func (d *jsiiProxy_Duration) ToDays(opts *TimeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		d,
		"toDays",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Return the total number of hours in this Duration.
//
// Returns: the value of this `Duration` expressed in Hours.
// Experimental.
func (d *jsiiProxy_Duration) ToHours(opts *TimeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		d,
		"toHours",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Turn this duration into a human-readable string.
// Experimental.
func (d *jsiiProxy_Duration) ToHumanString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toHumanString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return an ISO 8601 representation of this period.
//
// Returns: a string starting with 'P' describing the period
// See: https://www.iso.org/fr/standard/70907.html
//
// Experimental.
func (d *jsiiProxy_Duration) ToIsoString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toIsoString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return an ISO 8601 representation of this period.
//
// Returns: a string starting with 'P' describing the period
// See: https://www.iso.org/fr/standard/70907.html
//
// Deprecated: Use `toIsoString()` instead.
func (d *jsiiProxy_Duration) ToISOString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toISOString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return the total number of milliseconds in this Duration.
//
// Returns: the value of this `Duration` expressed in Milliseconds.
// Experimental.
func (d *jsiiProxy_Duration) ToMilliseconds(opts *TimeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		d,
		"toMilliseconds",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Return the total number of minutes in this Duration.
//
// Returns: the value of this `Duration` expressed in Minutes.
// Experimental.
func (d *jsiiProxy_Duration) ToMinutes(opts *TimeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		d,
		"toMinutes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Return the total number of seconds in this Duration.
//
// Returns: the value of this `Duration` expressed in Seconds.
// Experimental.
func (d *jsiiProxy_Duration) ToSeconds(opts *TimeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		d,
		"toSeconds",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Returns a string representation of this `Duration` that is also a Token that cannot be successfully resolved.
//
// This
// protects users against inadvertently stringifying a `Duration` object, when they should have called one of the
// `to*` methods instead.
// Experimental.
func (d *jsiiProxy_Duration) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns unit of the duration.
// Experimental.
func (d *jsiiProxy_Duration) UnitLabel() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"unitLabel",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to string encodings.
// Experimental.
type EncodingOptions struct {
	// A hint for the Token's purpose when stringifying it.
	// Experimental.
	DisplayHint *string `json:"displayHint"`
}

// The deployment environment for a stack.
// Experimental.
type Environment struct {
	// The AWS account ID for this environment.
	//
	// This can be either a concrete value such as `585191031104` or `Aws.accountId` which
	// indicates that account ID will only be determined during deployment (it
	// will resolve to the CloudFormation intrinsic `{"Ref":"AWS::AccountId"}`).
	// Note that certain features, such as cross-stack references and
	// environmental context providers require concerete region information and
	// will cause this stack to emit synthesis errors.
	// Experimental.
	Account *string `json:"account"`
	// The AWS region for this environment.
	//
	// This can be either a concrete value such as `eu-west-2` or `Aws.region`
	// which indicates that account ID will only be determined during deployment
	// (it will resolve to the CloudFormation intrinsic `{"Ref":"AWS::Region"}`).
	// Note that certain features, such as cross-stack references and
	// environmental context providers require concerete region information and
	// will cause this stack to emit synthesis errors.
	// Experimental.
	Region *string `json:"region"`
}

// Represents a date of expiration.
//
// The amount can be specified either as a Date object, timestamp, Duration or string.
// Experimental.
type Expiration interface {
	Date() *time.Time
	IsAfter(t Duration) *bool
	IsBefore(t Duration) *bool
	ToEpoch() *float64
}

// The jsii proxy struct for Expiration
type jsiiProxy_Expiration struct {
	_ byte // padding
}

func (j *jsiiProxy_Expiration) Date() *time.Time {
	var returns *time.Time
	_jsii_.Get(
		j,
		"date",
		&returns,
	)
	return returns
}


// Expire once the specified duration has passed since deployment time.
// Experimental.
func Expiration_After(t Duration) Expiration {
	_init_.Initialize()

	var returns Expiration

	_jsii_.StaticInvoke(
		"monocdk.Expiration",
		"after",
		[]interface{}{t},
		&returns,
	)

	return returns
}

// Expire at the specified date.
// Experimental.
func Expiration_AtDate(d *time.Time) Expiration {
	_init_.Initialize()

	var returns Expiration

	_jsii_.StaticInvoke(
		"monocdk.Expiration",
		"atDate",
		[]interface{}{d},
		&returns,
	)

	return returns
}

// Expire at the specified timestamp.
// Experimental.
func Expiration_AtTimestamp(t *float64) Expiration {
	_init_.Initialize()

	var returns Expiration

	_jsii_.StaticInvoke(
		"monocdk.Expiration",
		"atTimestamp",
		[]interface{}{t},
		&returns,
	)

	return returns
}

// Expire at specified date, represented as a string.
// Experimental.
func Expiration_FromString(s *string) Expiration {
	_init_.Initialize()

	var returns Expiration

	_jsii_.StaticInvoke(
		"monocdk.Expiration",
		"fromString",
		[]interface{}{s},
		&returns,
	)

	return returns
}

// Check if Exipiration expires after input.
// Experimental.
func (e *jsiiProxy_Expiration) IsAfter(t Duration) *bool {
	var returns *bool

	_jsii_.Invoke(
		e,
		"isAfter",
		[]interface{}{t},
		&returns,
	)

	return returns
}

// Check if Exipiration expires before input.
// Experimental.
func (e *jsiiProxy_Expiration) IsBefore(t Duration) *bool {
	var returns *bool

	_jsii_.Invoke(
		e,
		"isBefore",
		[]interface{}{t},
		&returns,
	)

	return returns
}

// Exipration Value in a formatted Unix Epoch Time in seconds.
// Experimental.
func (e *jsiiProxy_Expiration) ToEpoch() *float64 {
	var returns *float64

	_jsii_.Invoke(
		e,
		"toEpoch",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options for the `stack.exportValue()` method.
// Experimental.
type ExportValueOptions struct {
	// The name of the export to create.
	// Experimental.
	Name *string `json:"name"`
}

// Features that are implemented behind a flag in order to preserve backwards compatibility for existing apps.
//
// The list of flags are available in the
// `@aws-cdk/cx-api` module.
//
// The state of the flag for this application is stored as a CDK context variable.
// Experimental.
type FeatureFlags interface {
	IsEnabled(featureFlag *string) *bool
}

// The jsii proxy struct for FeatureFlags
type jsiiProxy_FeatureFlags struct {
	_ byte // padding
}

// Inspect feature flags on the construct node's context.
// Experimental.
func FeatureFlags_Of(scope Construct) FeatureFlags {
	_init_.Initialize()

	var returns FeatureFlags

	_jsii_.StaticInvoke(
		"monocdk.FeatureFlags",
		"of",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Check whether a feature flag is enabled.
//
// If configured, the flag is present in
// the construct node context. Falls back to the defaults defined in the `cx-api`
// module.
// Experimental.
func (f *jsiiProxy_FeatureFlags) IsEnabled(featureFlag *string) *bool {
	var returns *bool

	_jsii_.Invoke(
		f,
		"isEnabled",
		[]interface{}{featureFlag},
		&returns,
	)

	return returns
}

// The location of the published file asset.
//
// This is where the asset
// can be consumed at runtime.
// Experimental.
type FileAssetLocation struct {
	// The name of the Amazon S3 bucket.
	// Experimental.
	BucketName *string `json:"bucketName"`
	// The HTTP URL of this asset on Amazon S3.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	HttpUrl *string `json:"httpUrl"`
	// The Amazon S3 object key.
	// Experimental.
	ObjectKey *string `json:"objectKey"`
	// The S3 URL of this asset on Amazon S3.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	S3ObjectUrl *string `json:"s3ObjectUrl"`
	// The ARN of the KMS key used to encrypt the file asset bucket, if any.
	//
	// If so, the consuming role should be given "kms:Decrypt" permissions in its
	// identity policy.
	//
	// It's the responsibility of they key's creator to make sure that all
	// consumers that the key's key policy is configured such that the key can be used
	// by all consumers that need it.
	//
	// The default bootstrap stack provisioned by the CDK CLI ensures this, and
	// can be used as an example for how to configure the key properly.
	// Deprecated: Since bootstrap bucket v4, the key policy properly allows use of the
	// key via the bucket and no additional parameters have to be granted anymore.
	KmsKeyArn *string `json:"kmsKeyArn"`
	// The HTTP URL of this asset on Amazon S3.
	// Deprecated: use `httpUrl`
	S3Url *string `json:"s3Url"`
}

// Packaging modes for file assets.
// Experimental.
type FileAssetPackaging string

const (
	FileAssetPackaging_ZIP_DIRECTORY FileAssetPackaging = "ZIP_DIRECTORY"
	FileAssetPackaging_FILE FileAssetPackaging = "FILE"
)

// Represents the source for a file asset.
// Experimental.
type FileAssetSource struct {
	// A hash on the content source.
	//
	// This hash is used to uniquely identify this
	// asset throughout the system. If this value doesn't change, the asset will
	// not be rebuilt or republished.
	// Experimental.
	SourceHash *string `json:"sourceHash"`
	// An external command that will produce the packaged asset.
	//
	// The command should produce the location of a ZIP file on `stdout`.
	// Experimental.
	Executable *[]*string `json:"executable"`
	// The path, relative to the root of the cloud assembly, in which this asset source resides.
	//
	// This can be a path to a file or a directory, depending on the
	// packaging type.
	// Experimental.
	FileName *string `json:"fileName"`
	// Which type of packaging to perform.
	// Experimental.
	Packaging FileAssetPackaging `json:"packaging"`
}

// Options applied when copying directories into the staging location.
// Experimental.
type FileCopyOptions struct {
	// Glob patterns to exclude from the copy.
	// Experimental.
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Experimental.
	FollowSymlinks SymlinkFollowMode `json:"followSymlinks"`
	// The ignore behavior to use for exclude patterns.
	// Experimental.
	IgnoreMode IgnoreMode `json:"ignoreMode"`
}

// Options related to calculating source hash.
// Experimental.
type FileFingerprintOptions struct {
	// Glob patterns to exclude from the copy.
	// Experimental.
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Experimental.
	FollowSymlinks SymlinkFollowMode `json:"followSymlinks"`
	// The ignore behavior to use for exclude patterns.
	// Experimental.
	IgnoreMode IgnoreMode `json:"ignoreMode"`
	// Extra information to encode into the fingerprint (e.g. build instructions and other inputs).
	// Experimental.
	ExtraHash *string `json:"extraHash"`
}

// File system utilities.
// Experimental.
type FileSystem interface {
}

// The jsii proxy struct for FileSystem
type jsiiProxy_FileSystem struct {
	_ byte // padding
}

// Experimental.
func NewFileSystem() FileSystem {
	_init_.Initialize()

	j := jsiiProxy_FileSystem{}

	_jsii_.Create(
		"monocdk.FileSystem",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewFileSystem_Override(f FileSystem) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.FileSystem",
		nil, // no parameters
		f,
	)
}

// Copies an entire directory structure.
// Experimental.
func FileSystem_CopyDirectory(srcDir *string, destDir *string, options *CopyOptions, rootDir *string) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.FileSystem",
		"copyDirectory",
		[]interface{}{srcDir, destDir, options, rootDir},
	)
}

// Produces fingerprint based on the contents of a single file or an entire directory tree.
//
// The fingerprint will also include:
// 1. An extra string if defined in `options.extra`.
// 2. The set of exclude patterns, if defined in `options.exclude`
// 3. The symlink follow mode value.
// Experimental.
func FileSystem_Fingerprint(fileOrDirectory *string, options *FingerprintOptions) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.FileSystem",
		"fingerprint",
		[]interface{}{fileOrDirectory, options},
		&returns,
	)

	return returns
}

// Checks whether a directory is empty.
// Experimental.
func FileSystem_IsEmpty(dir *string) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.FileSystem",
		"isEmpty",
		[]interface{}{dir},
		&returns,
	)

	return returns
}

// Creates a unique temporary directory in the **system temp directory**.
// Experimental.
func FileSystem_Mkdtemp(prefix *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.FileSystem",
		"mkdtemp",
		[]interface{}{prefix},
		&returns,
	)

	return returns
}

func FileSystem_Tmpdir() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.FileSystem",
		"tmpdir",
		&returns,
	)
	return returns
}

// Options related to calculating source hash.
// Experimental.
type FingerprintOptions struct {
	// Glob patterns to exclude from the copy.
	// Experimental.
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Experimental.
	Follow SymlinkFollowMode `json:"follow"`
	// The ignore behavior to use for exclude patterns.
	// Experimental.
	IgnoreMode IgnoreMode `json:"ignoreMode"`
	// Extra information to encode into the fingerprint (e.g. build instructions and other inputs).
	// Experimental.
	ExtraHash *string `json:"extraHash"`
}

// CloudFormation intrinsic functions.
//
// http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference.html
// Experimental.
type Fn interface {
}

// The jsii proxy struct for Fn
type jsiiProxy_Fn struct {
	_ byte // padding
}

// The intrinsic function ``Fn::Base64`` returns the Base64 representation of the input string.
//
// This function is typically used to pass encoded data to
// Amazon EC2 instances by way of the UserData property.
//
// Returns: a token represented as a string
// Experimental.
func Fn_Base64(data *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"base64",
		[]interface{}{data},
		&returns,
	)

	return returns
}

// The intrinsic function ``Fn::Cidr`` returns the specified Cidr address block.
//
// Returns: a token represented as a string
// Experimental.
func Fn_Cidr(ipBlock *string, count *float64, sizeMask *string) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"cidr",
		[]interface{}{ipBlock, count, sizeMask},
		&returns,
	)

	return returns
}

// Returns true if all the specified conditions evaluate to true, or returns false if any one of the conditions evaluates to false.
//
// ``Fn::And`` acts as
// an AND operator. The minimum number of conditions that you can include is
// 1.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionAnd(conditions ...ICfnConditionExpression) ICfnConditionExpression {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range conditions {
		args = append(args, a)
	}

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionAnd",
		args,
		&returns,
	)

	return returns
}

// Returns true if a specified string matches at least one value in a list of strings.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionContains(listOfStrings *[]*string, value *string) ICfnConditionExpression {
	_init_.Initialize()

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionContains",
		[]interface{}{listOfStrings, value},
		&returns,
	)

	return returns
}

// Returns true if a specified string matches all values in a list.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionEachMemberEquals(listOfStrings *[]*string, value *string) ICfnConditionExpression {
	_init_.Initialize()

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionEachMemberEquals",
		[]interface{}{listOfStrings, value},
		&returns,
	)

	return returns
}

// Returns true if each member in a list of strings matches at least one value in a second list of strings.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionEachMemberIn(stringsToCheck *[]*string, stringsToMatch *[]*string) ICfnConditionExpression {
	_init_.Initialize()

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionEachMemberIn",
		[]interface{}{stringsToCheck, stringsToMatch},
		&returns,
	)

	return returns
}

// Compares if two values are equal.
//
// Returns true if the two values are equal
// or false if they aren't.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionEquals(lhs interface{}, rhs interface{}) ICfnConditionExpression {
	_init_.Initialize()

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionEquals",
		[]interface{}{lhs, rhs},
		&returns,
	)

	return returns
}

// Returns one value if the specified condition evaluates to true and another value if the specified condition evaluates to false.
//
// Currently, AWS
// CloudFormation supports the ``Fn::If`` intrinsic function in the metadata
// attribute, update policy attribute, and property values in the Resources
// section and Outputs sections of a template. You can use the AWS::NoValue
// pseudo parameter as a return value to remove the corresponding property.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionIf(conditionId *string, valueIfTrue interface{}, valueIfFalse interface{}) ICfnConditionExpression {
	_init_.Initialize()

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionIf",
		[]interface{}{conditionId, valueIfTrue, valueIfFalse},
		&returns,
	)

	return returns
}

// Returns true for a condition that evaluates to false or returns false for a condition that evaluates to true.
//
// ``Fn::Not`` acts as a NOT operator.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionNot(condition ICfnConditionExpression) ICfnConditionExpression {
	_init_.Initialize()

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionNot",
		[]interface{}{condition},
		&returns,
	)

	return returns
}

// Returns true if any one of the specified conditions evaluate to true, or returns false if all of the conditions evaluates to false.
//
// ``Fn::Or`` acts
// as an OR operator. The minimum number of conditions that you can include is
// 1.
//
// Returns: an FnCondition token
// Experimental.
func Fn_ConditionOr(conditions ...ICfnConditionExpression) ICfnConditionExpression {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range conditions {
		args = append(args, a)
	}

	var returns ICfnConditionExpression

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"conditionOr",
		args,
		&returns,
	)

	return returns
}

// The intrinsic function ``Fn::FindInMap`` returns the value corresponding to keys in a two-level map that is declared in the Mappings section.
//
// Returns: a token represented as a string
// Experimental.
func Fn_FindInMap(mapName *string, topLevelKey *string, secondLevelKey *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"findInMap",
		[]interface{}{mapName, topLevelKey, secondLevelKey},
		&returns,
	)

	return returns
}

// The ``Fn::GetAtt`` intrinsic function returns the value of an attribute from a resource in the template.
//
// Returns: an IResolvable object
// Experimental.
func Fn_GetAtt(logicalNameOfResource *string, attributeName *string) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"getAtt",
		[]interface{}{logicalNameOfResource, attributeName},
		&returns,
	)

	return returns
}

// The intrinsic function ``Fn::GetAZs`` returns an array that lists Availability Zones for a specified region.
//
// Because customers have access to
// different Availability Zones, the intrinsic function ``Fn::GetAZs`` enables
// template authors to write templates that adapt to the calling user's
// access. That way you don't have to hard-code a full list of Availability
// Zones for a specified region.
//
// Returns: a token represented as a string array
// Experimental.
func Fn_GetAzs(region *string) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"getAzs",
		[]interface{}{region},
		&returns,
	)

	return returns
}

// Like `Fn.importValue`, but import a list with a known length.
//
// If you explicitly want a list with an unknown length, call `Fn.split(',',
// Fn.importValue(exportName))`. See the documentation of `Fn.split` to read
// more about the limitations of using lists of unknown length.
//
// `Fn.importListValue(exportName, assumedLength)` is the same as
// `Fn.split(',', Fn.importValue(exportName), assumedLength)`,
// but easier to read and impossible to forget to pass `assumedLength`.
// Experimental.
func Fn_ImportListValue(sharedValueToImport *string, assumedLength *float64, delimiter *string) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"importListValue",
		[]interface{}{sharedValueToImport, assumedLength, delimiter},
		&returns,
	)

	return returns
}

// The intrinsic function ``Fn::ImportValue`` returns the value of an output exported by another stack.
//
// You typically use this function to create
// cross-stack references. In the following example template snippets, Stack A
// exports VPC security group values and Stack B imports them.
//
// Returns: a token represented as a string
// Experimental.
func Fn_ImportValue(sharedValueToImport *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"importValue",
		[]interface{}{sharedValueToImport},
		&returns,
	)

	return returns
}

// The intrinsic function ``Fn::Join`` appends a set of values into a single value, separated by the specified delimiter.
//
// If a delimiter is the empty
// string, the set of values are concatenated with no delimiter.
//
// Returns: a token represented as a string
// Experimental.
func Fn_Join(delimiter *string, listOfValues *[]*string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"join",
		[]interface{}{delimiter, listOfValues},
		&returns,
	)

	return returns
}

// Given an url, parse the domain name.
// Experimental.
func Fn_ParseDomainName(url *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"parseDomainName",
		[]interface{}{url},
		&returns,
	)

	return returns
}

// The ``Ref`` intrinsic function returns the value of the specified parameter or resource.
//
// Note that it doesn't validate the logicalName, it mainly serves paremeter/resource reference defined in a ``CfnInclude`` template.
// Experimental.
func Fn_Ref(logicalName *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"ref",
		[]interface{}{logicalName},
		&returns,
	)

	return returns
}

// Returns all values for a specified parameter type.
//
// Returns: a token represented as a string array
// Experimental.
func Fn_RefAll(parameterType *string) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"refAll",
		[]interface{}{parameterType},
		&returns,
	)

	return returns
}

// The intrinsic function ``Fn::Select`` returns a single object from a list of objects by index.
//
// Returns: a token represented as a string
// Experimental.
func Fn_Select(index *float64, array *[]*string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"select",
		[]interface{}{index, array},
		&returns,
	)

	return returns
}

// Split a string token into a token list of string values.
//
// Specify the location of splits with a delimiter such as ',' (a comma).
// Renders to the `Fn::Split` intrinsic function.
//
// Lists with unknown lengths (default)
// -------------------------------------
//
// Since this function is used to work with deploy-time values, if `assumedLength`
// is not given the CDK cannot know the length of the resulting list at synthesis time.
// This brings the following restrictions:
//
// - You must use `Fn.select(i, list)` to pick elements out of the list (you must not use
//    `list[i]`).
// - You cannot add elements to the list, remove elements from the list,
//    combine two such lists together, or take a slice of the list.
// - You cannot pass the list to constructs that do any of the above.
//
// The only valid operation with such a tokenized list is to pass it unmodified to a
// CloudFormation Resource construct.
//
// Lists with assumed lengths
// --------------------------
//
// Pass `assumedLength` if you know the length of the list that will be
// produced by splitting. The actual list length at deploy time may be
// *longer* than the number you pass, but not *shorter*.
//
// The returned list will look like:
//
// ```
// [Fn.select(0, split), Fn.select(1, split), Fn.select(2, split), ...]
// ```
//
// The restrictions from the section "Lists with unknown lengths" will now be lifted,
// at the expense of having to know and fix the length of the list.
//
// Returns: a token represented as a string array
// Experimental.
func Fn_Split(delimiter *string, source *string, assumedLength *float64) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"split",
		[]interface{}{delimiter, source, assumedLength},
		&returns,
	)

	return returns
}

// The intrinsic function ``Fn::Sub`` substitutes variables in an input string with values that you specify.
//
// In your templates, you can use this function
// to construct commands or outputs that include values that aren't available
// until you create or update a stack.
//
// Returns: a token represented as a string
// Experimental.
func Fn_Sub(body *string, variables *map[string]*string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"sub",
		[]interface{}{body, variables},
		&returns,
	)

	return returns
}

// Creates a token representing the ``Fn::Transform`` expression.
//
// Returns: a token representing the transform expression
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-transform.html
//
// Experimental.
func Fn_Transform(macroName *string, parameters *map[string]interface{}) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"transform",
		[]interface{}{macroName, parameters},
		&returns,
	)

	return returns
}

// Returns an attribute value or list of values for a specific parameter and attribute.
//
// Returns: a token represented as a string
// Experimental.
func Fn_ValueOf(parameterOrLogicalId *string, attribute *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"valueOf",
		[]interface{}{parameterOrLogicalId, attribute},
		&returns,
	)

	return returns
}

// Returns a list of all attribute values for a given parameter type and attribute.
//
// Returns: a token represented as a string array
// Experimental.
func Fn_ValueOfAll(parameterType *string, attribute *string) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Fn",
		"valueOfAll",
		[]interface{}{parameterType, attribute},
		&returns,
	)

	return returns
}

// Experimental.
type GetContextKeyOptions struct {
	// The context provider to query.
	// Experimental.
	Provider *string `json:"provider"`
	// Provider-specific properties.
	// Experimental.
	Props *map[string]interface{} `json:"props"`
}

// Experimental.
type GetContextKeyResult struct {
	// Experimental.
	Key *string `json:"key"`
	// Experimental.
	Props *map[string]interface{} `json:"props"`
}

// Experimental.
type GetContextValueOptions struct {
	// The context provider to query.
	// Experimental.
	Provider *string `json:"provider"`
	// Provider-specific properties.
	// Experimental.
	Props *map[string]interface{} `json:"props"`
	// The value to return if the context value was not found and a missing context is reported.
	//
	// This should be a dummy value that should preferably
	// fail during deployment since it represents an invalid state.
	// Experimental.
	DummyValue interface{} `json:"dummyValue"`
}

// Experimental.
type GetContextValueResult struct {
	// Experimental.
	Value interface{} `json:"value"`
}

// Ignores file paths based on the [`.gitignore specification`](https://git-scm.com/docs/gitignore).
// Experimental.
type GitIgnoreStrategy interface {
	IgnoreStrategy
	Add(pattern *string)
	Ignores(absoluteFilePath *string) *bool
}

// The jsii proxy struct for GitIgnoreStrategy
type jsiiProxy_GitIgnoreStrategy struct {
	jsiiProxy_IgnoreStrategy
}

// Experimental.
func NewGitIgnoreStrategy(absoluteRootPath *string, patterns *[]*string) GitIgnoreStrategy {
	_init_.Initialize()

	j := jsiiProxy_GitIgnoreStrategy{}

	_jsii_.Create(
		"monocdk.GitIgnoreStrategy",
		[]interface{}{absoluteRootPath, patterns},
		&j,
	)

	return &j
}

// Experimental.
func NewGitIgnoreStrategy_Override(g GitIgnoreStrategy, absoluteRootPath *string, patterns *[]*string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.GitIgnoreStrategy",
		[]interface{}{absoluteRootPath, patterns},
		g,
	)
}

// Ignores file paths based on the [`.dockerignore specification`](https://docs.docker.com/engine/reference/builder/#dockerignore-file).
//
// Returns: `DockerIgnorePattern` associated with the given patterns.
// Experimental.
func GitIgnoreStrategy_Docker(absoluteRootPath *string, patterns *[]*string) DockerIgnoreStrategy {
	_init_.Initialize()

	var returns DockerIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GitIgnoreStrategy",
		"docker",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Creates an IgnoreStrategy based on the `ignoreMode` and `exclude` in a `CopyOptions`.
//
// Returns: `IgnoreStrategy` based on the `CopyOptions`
// Experimental.
func GitIgnoreStrategy_FromCopyOptions(options *CopyOptions, absoluteRootPath *string) IgnoreStrategy {
	_init_.Initialize()

	var returns IgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GitIgnoreStrategy",
		"fromCopyOptions",
		[]interface{}{options, absoluteRootPath},
		&returns,
	)

	return returns
}

// Ignores file paths based on the [`.gitignore specification`](https://git-scm.com/docs/gitignore).
//
// Returns: `GitIgnorePattern` associated with the given patterns.
// Experimental.
func GitIgnoreStrategy_Git(absoluteRootPath *string, patterns *[]*string) GitIgnoreStrategy {
	_init_.Initialize()

	var returns GitIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GitIgnoreStrategy",
		"git",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Ignores file paths based on simple glob patterns.
//
// Returns: `GlobIgnorePattern` associated with the given patterns.
// Experimental.
func GitIgnoreStrategy_Glob(absoluteRootPath *string, patterns *[]*string) GlobIgnoreStrategy {
	_init_.Initialize()

	var returns GlobIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GitIgnoreStrategy",
		"glob",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Adds another pattern.
// Experimental.
func (g *jsiiProxy_GitIgnoreStrategy) Add(pattern *string) {
	_jsii_.InvokeVoid(
		g,
		"add",
		[]interface{}{pattern},
	)
}

// Determines whether a given file path should be ignored or not.
//
// Returns: `true` if the file should be ignored
// Experimental.
func (g *jsiiProxy_GitIgnoreStrategy) Ignores(absoluteFilePath *string) *bool {
	var returns *bool

	_jsii_.Invoke(
		g,
		"ignores",
		[]interface{}{absoluteFilePath},
		&returns,
	)

	return returns
}

// Ignores file paths based on simple glob patterns.
// Experimental.
type GlobIgnoreStrategy interface {
	IgnoreStrategy
	Add(pattern *string)
	Ignores(absoluteFilePath *string) *bool
}

// The jsii proxy struct for GlobIgnoreStrategy
type jsiiProxy_GlobIgnoreStrategy struct {
	jsiiProxy_IgnoreStrategy
}

// Experimental.
func NewGlobIgnoreStrategy(absoluteRootPath *string, patterns *[]*string) GlobIgnoreStrategy {
	_init_.Initialize()

	j := jsiiProxy_GlobIgnoreStrategy{}

	_jsii_.Create(
		"monocdk.GlobIgnoreStrategy",
		[]interface{}{absoluteRootPath, patterns},
		&j,
	)

	return &j
}

// Experimental.
func NewGlobIgnoreStrategy_Override(g GlobIgnoreStrategy, absoluteRootPath *string, patterns *[]*string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.GlobIgnoreStrategy",
		[]interface{}{absoluteRootPath, patterns},
		g,
	)
}

// Ignores file paths based on the [`.dockerignore specification`](https://docs.docker.com/engine/reference/builder/#dockerignore-file).
//
// Returns: `DockerIgnorePattern` associated with the given patterns.
// Experimental.
func GlobIgnoreStrategy_Docker(absoluteRootPath *string, patterns *[]*string) DockerIgnoreStrategy {
	_init_.Initialize()

	var returns DockerIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GlobIgnoreStrategy",
		"docker",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Creates an IgnoreStrategy based on the `ignoreMode` and `exclude` in a `CopyOptions`.
//
// Returns: `IgnoreStrategy` based on the `CopyOptions`
// Experimental.
func GlobIgnoreStrategy_FromCopyOptions(options *CopyOptions, absoluteRootPath *string) IgnoreStrategy {
	_init_.Initialize()

	var returns IgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GlobIgnoreStrategy",
		"fromCopyOptions",
		[]interface{}{options, absoluteRootPath},
		&returns,
	)

	return returns
}

// Ignores file paths based on the [`.gitignore specification`](https://git-scm.com/docs/gitignore).
//
// Returns: `GitIgnorePattern` associated with the given patterns.
// Experimental.
func GlobIgnoreStrategy_Git(absoluteRootPath *string, patterns *[]*string) GitIgnoreStrategy {
	_init_.Initialize()

	var returns GitIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GlobIgnoreStrategy",
		"git",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Ignores file paths based on simple glob patterns.
//
// Returns: `GlobIgnorePattern` associated with the given patterns.
// Experimental.
func GlobIgnoreStrategy_Glob(absoluteRootPath *string, patterns *[]*string) GlobIgnoreStrategy {
	_init_.Initialize()

	var returns GlobIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.GlobIgnoreStrategy",
		"glob",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Adds another pattern.
// Experimental.
func (g *jsiiProxy_GlobIgnoreStrategy) Add(pattern *string) {
	_jsii_.InvokeVoid(
		g,
		"add",
		[]interface{}{pattern},
	)
}

// Determines whether a given file path should be ignored or not.
//
// Returns: `true` if the file should be ignored
// Experimental.
func (g *jsiiProxy_GlobIgnoreStrategy) Ignores(absoluteFilePath *string) *bool {
	var returns *bool

	_jsii_.Invoke(
		g,
		"ignores",
		[]interface{}{absoluteFilePath},
		&returns,
	)

	return returns
}

// Interface for lazy untyped value producers.
// Experimental.
type IAnyProducer interface {
	// Produce the value.
	// Experimental.
	Produce(context IResolveContext) interface{}
}

// The jsii proxy for IAnyProducer
type jsiiProxy_IAnyProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_IAnyProducer) Produce(context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"produce",
		[]interface{}{context},
		&returns,
	)

	return returns
}

// Represents an Aspect.
// Experimental.
type IAspect interface {
	// All aspects can visit an IConstruct.
	// Experimental.
	Visit(node IConstruct)
}

// The jsii proxy for IAspect
type jsiiProxy_IAspect struct {
	_ byte // padding
}

func (i *jsiiProxy_IAspect) Visit(node IConstruct) {
	_jsii_.InvokeVoid(
		i,
		"visit",
		[]interface{}{node},
	)
}

// Common interface for all assets.
// Experimental.
type IAsset interface {
	// A hash of this asset, which is available at construction time.
	//
	// As this is a plain string, it
	// can be used in construct IDs in order to enforce creation of a new resource when the content
	// hash has changed.
	// Experimental.
	AssetHash() *string
}

// The jsii proxy for IAsset
type jsiiProxy_IAsset struct {
	_ byte // padding
}

func (j *jsiiProxy_IAsset) AssetHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assetHash",
		&returns,
	)
	return returns
}

// Represents a CloudFormation element that can be used within a Condition.
//
// You can use intrinsic functions, such as ``Fn.conditionIf``,
// ``Fn.conditionEquals``, and ``Fn.conditionNot``, to conditionally create
// stack resources. These conditions are evaluated based on input parameters
// that you declare when you create or update a stack. After you define all your
// conditions, you can associate them with resources or resource properties in
// the Resources and Outputs sections of a template.
//
// You define all conditions in the Conditions section of a template except for
// ``Fn.conditionIf`` conditions. You can use the ``Fn.conditionIf`` condition
// in the metadata attribute, update policy attribute, and property values in
// the Resources section and Outputs sections of a template.
//
// You might use conditions when you want to reuse a template that can create
// resources in different contexts, such as a test environment versus a
// production environment. In your template, you can add an EnvironmentType
// input parameter, which accepts either prod or test as inputs. For the
// production environment, you might include Amazon EC2 instances with certain
// capabilities; however, for the test environment, you want to use less
// capabilities to save costs. With conditions, you can define which resources
// are created and how they're configured for each environment type.
//
// You can use `toString` when you wish to embed a condition expression
// in a property value that accepts a `string`. For example:
//
// ```ts
// new sqs.Queue(this, 'MyQueue', {
//    queueName: Fn.conditionIf('Condition', 'Hello', 'World').toString()
// });
// ```
// Experimental.
type ICfnConditionExpression interface {
	IResolvable
}

// The jsii proxy for ICfnConditionExpression
type jsiiProxy_ICfnConditionExpression struct {
	jsiiProxy_IResolvable
}

// Experimental.
type ICfnResourceOptions interface {
	// A condition to associate with this resource.
	//
	// This means that only if the condition evaluates to 'true' when the stack
	// is deployed, the resource will be included. This is provided to allow CDK projects to produce legacy templates, but noramlly
	// there is no need to use it in CDK projects.
	// Experimental.
	Condition() CfnCondition
	// A condition to associate with this resource.
	//
	// This means that only if the condition evaluates to 'true' when the stack
	// is deployed, the resource will be included. This is provided to allow CDK projects to produce legacy templates, but noramlly
	// there is no need to use it in CDK projects.
	// Experimental.
	SetCondition(c CfnCondition)
	// Associate the CreationPolicy attribute with a resource to prevent its status from reaching create complete until AWS CloudFormation receives a specified number of success signals or the timeout period is exceeded.
	//
	// To signal a
	// resource, you can use the cfn-signal helper script or SignalResource API. AWS CloudFormation publishes valid signals
	// to the stack events so that you track the number of signals sent.
	// Experimental.
	CreationPolicy() *CfnCreationPolicy
	// Associate the CreationPolicy attribute with a resource to prevent its status from reaching create complete until AWS CloudFormation receives a specified number of success signals or the timeout period is exceeded.
	//
	// To signal a
	// resource, you can use the cfn-signal helper script or SignalResource API. AWS CloudFormation publishes valid signals
	// to the stack events so that you track the number of signals sent.
	// Experimental.
	SetCreationPolicy(c *CfnCreationPolicy)
	// With the DeletionPolicy attribute you can preserve or (in some cases) backup a resource when its stack is deleted.
	//
	// You specify a DeletionPolicy attribute for each resource that you want to control. If a resource has no DeletionPolicy
	// attribute, AWS CloudFormation deletes the resource by default. Note that this capability also applies to update operations
	// that lead to resources being removed.
	// Experimental.
	DeletionPolicy() CfnDeletionPolicy
	// With the DeletionPolicy attribute you can preserve or (in some cases) backup a resource when its stack is deleted.
	//
	// You specify a DeletionPolicy attribute for each resource that you want to control. If a resource has no DeletionPolicy
	// attribute, AWS CloudFormation deletes the resource by default. Note that this capability also applies to update operations
	// that lead to resources being removed.
	// Experimental.
	SetDeletionPolicy(d CfnDeletionPolicy)
	// The description of this resource.
	//
	// Used for informational purposes only, is not processed in any way
	// (and stays with the CloudFormation template, is not passed to the underlying resource,
	// even if it does have a 'description' property).
	// Experimental.
	Description() *string
	// The description of this resource.
	//
	// Used for informational purposes only, is not processed in any way
	// (and stays with the CloudFormation template, is not passed to the underlying resource,
	// even if it does have a 'description' property).
	// Experimental.
	SetDescription(d *string)
	// Metadata associated with the CloudFormation resource.
	//
	// This is not the same as the construct metadata which can be added
	// using construct.addMetadata(), but would not appear in the CloudFormation template automatically.
	// Experimental.
	Metadata() *map[string]interface{}
	// Metadata associated with the CloudFormation resource.
	//
	// This is not the same as the construct metadata which can be added
	// using construct.addMetadata(), but would not appear in the CloudFormation template automatically.
	// Experimental.
	SetMetadata(m *map[string]interface{})
	// Use the UpdatePolicy attribute to specify how AWS CloudFormation handles updates to the AWS::AutoScaling::AutoScalingGroup resource.
	//
	// AWS CloudFormation invokes one of three update policies depending on the type of change you make or whether a
	// scheduled action is associated with the Auto Scaling group.
	// Experimental.
	UpdatePolicy() *CfnUpdatePolicy
	// Use the UpdatePolicy attribute to specify how AWS CloudFormation handles updates to the AWS::AutoScaling::AutoScalingGroup resource.
	//
	// AWS CloudFormation invokes one of three update policies depending on the type of change you make or whether a
	// scheduled action is associated with the Auto Scaling group.
	// Experimental.
	SetUpdatePolicy(u *CfnUpdatePolicy)
	// Use the UpdateReplacePolicy attribute to retain or (in some cases) backup the existing physical instance of a resource when it is replaced during a stack update operation.
	// Experimental.
	UpdateReplacePolicy() CfnDeletionPolicy
	// Use the UpdateReplacePolicy attribute to retain or (in some cases) backup the existing physical instance of a resource when it is replaced during a stack update operation.
	// Experimental.
	SetUpdateReplacePolicy(u CfnDeletionPolicy)
	// The version of this resource.
	//
	// Used only for custom CloudFormation resources.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cfn-customresource.html
	//
	// Experimental.
	Version() *string
	// The version of this resource.
	//
	// Used only for custom CloudFormation resources.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cfn-customresource.html
	//
	// Experimental.
	SetVersion(v *string)
}

// The jsii proxy for ICfnResourceOptions
type jsiiProxy_ICfnResourceOptions struct {
	_ byte // padding
}

func (j *jsiiProxy_ICfnResourceOptions) Condition() CfnCondition {
	var returns CfnCondition
	_jsii_.Get(
		j,
		"condition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetCondition(val CfnCondition) {
	_jsii_.Set(
		j,
		"condition",
		val,
	)
}

func (j *jsiiProxy_ICfnResourceOptions) CreationPolicy() *CfnCreationPolicy {
	var returns *CfnCreationPolicy
	_jsii_.Get(
		j,
		"creationPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetCreationPolicy(val *CfnCreationPolicy) {
	_jsii_.Set(
		j,
		"creationPolicy",
		val,
	)
}

func (j *jsiiProxy_ICfnResourceOptions) DeletionPolicy() CfnDeletionPolicy {
	var returns CfnDeletionPolicy
	_jsii_.Get(
		j,
		"deletionPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetDeletionPolicy(val CfnDeletionPolicy) {
	_jsii_.Set(
		j,
		"deletionPolicy",
		val,
	)
}

func (j *jsiiProxy_ICfnResourceOptions) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_ICfnResourceOptions) Metadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"metadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetMetadata(val *map[string]interface{}) {
	_jsii_.Set(
		j,
		"metadata",
		val,
	)
}

func (j *jsiiProxy_ICfnResourceOptions) UpdatePolicy() *CfnUpdatePolicy {
	var returns *CfnUpdatePolicy
	_jsii_.Get(
		j,
		"updatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetUpdatePolicy(val *CfnUpdatePolicy) {
	_jsii_.Set(
		j,
		"updatePolicy",
		val,
	)
}

func (j *jsiiProxy_ICfnResourceOptions) UpdateReplacePolicy() CfnDeletionPolicy {
	var returns CfnDeletionPolicy
	_jsii_.Get(
		j,
		"updateReplacePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetUpdateReplacePolicy(val CfnDeletionPolicy) {
	_jsii_.Set(
		j,
		"updateReplacePolicy",
		val,
	)
}

func (j *jsiiProxy_ICfnResourceOptions) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ICfnResourceOptions) SetVersion(val *string) {
	_jsii_.Set(
		j,
		"version",
		val,
	)
}

// Represents a construct.
// Experimental.
type IConstruct interface {
	constructs.IConstruct
	IDependable
	// The construct tree node for this construct.
	// Experimental.
	Node() ConstructNode
}

// The jsii proxy for IConstruct
type jsiiProxy_IConstruct struct {
	internal.Type__constructsIConstruct
	jsiiProxy_IDependable
}

func (j *jsiiProxy_IConstruct) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

// Trait marker for classes that can be depended upon.
//
// The presence of this interface indicates that an object has
// an `IDependableTrait` implementation.
//
// This interface can be used to take an (ordering) dependency on a set of
// constructs. An ordering dependency implies that the resources represented by
// those constructs are deployed before the resources depending ON them are
// deployed.
// Experimental.
type IDependable interface {
}

// The jsii proxy for IDependable
type jsiiProxy_IDependable struct {
	_ byte // padding
}

// Function used to concatenate symbols in the target document language.
//
// Interface so it could potentially be exposed over jsii.
// Experimental.
type IFragmentConcatenator interface {
	// Join the fragment on the left and on the right.
	// Experimental.
	Join(left interface{}, right interface{}) interface{}
}

// The jsii proxy for IFragmentConcatenator
type jsiiProxy_IFragmentConcatenator struct {
	_ byte // padding
}

func (i *jsiiProxy_IFragmentConcatenator) Join(left interface{}, right interface{}) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"join",
		[]interface{}{left, right},
		&returns,
	)

	return returns
}

// Interface for examining a construct and exposing metadata.
// Experimental.
type IInspectable interface {
	// Examines construct.
	// Experimental.
	Inspect(inspector TreeInspector)
}

// The jsii proxy for IInspectable
type jsiiProxy_IInspectable struct {
	_ byte // padding
}

func (i *jsiiProxy_IInspectable) Inspect(inspector TreeInspector) {
	_jsii_.InvokeVoid(
		i,
		"inspect",
		[]interface{}{inspector},
	)
}

// Interface for lazy list producers.
// Experimental.
type IListProducer interface {
	// Produce the list value.
	// Experimental.
	Produce(context IResolveContext) *[]*string
}

// The jsii proxy for IListProducer
type jsiiProxy_IListProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_IListProducer) Produce(context IResolveContext) *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		i,
		"produce",
		[]interface{}{context},
		&returns,
	)

	return returns
}

// Local bundling.
// Experimental.
type ILocalBundling interface {
	// This method is called before attempting docker bundling to allow the bundler to be executed locally.
	//
	// If the local bundler exists, and bundling
	// was performed locally, return `true`. Otherwise, return `false`.
	// Experimental.
	TryBundle(outputDir *string, options *BundlingOptions) *bool
}

// The jsii proxy for ILocalBundling
type jsiiProxy_ILocalBundling struct {
	_ byte // padding
}

func (i *jsiiProxy_ILocalBundling) TryBundle(outputDir *string, options *BundlingOptions) *bool {
	var returns *bool

	_jsii_.Invoke(
		i,
		"tryBundle",
		[]interface{}{outputDir, options},
		&returns,
	)

	return returns
}

// Interface for lazy number producers.
// Experimental.
type INumberProducer interface {
	// Produce the number value.
	// Experimental.
	Produce(context IResolveContext) *float64
}

// The jsii proxy for INumberProducer
type jsiiProxy_INumberProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_INumberProducer) Produce(context IResolveContext) *float64 {
	var returns *float64

	_jsii_.Invoke(
		i,
		"produce",
		[]interface{}{context},
		&returns,
	)

	return returns
}

// A Token that can post-process the complete resolved value, after resolve() has recursed over it.
// Experimental.
type IPostProcessor interface {
	// Process the completely resolved value, after full recursion/resolution has happened.
	// Experimental.
	PostProcess(input interface{}, context IResolveContext) interface{}
}

// The jsii proxy for IPostProcessor
type jsiiProxy_IPostProcessor struct {
	_ byte // padding
}

func (i *jsiiProxy_IPostProcessor) PostProcess(input interface{}, context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"postProcess",
		[]interface{}{input, context},
		&returns,
	)

	return returns
}

// Interface for values that can be resolvable later.
//
// Tokens are special objects that participate in synthesis.
// Experimental.
type IResolvable interface {
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(context IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// This may return an array with a single informational element indicating how
	// to get this property populated, if it was skipped for performance reasons.
	// Experimental.
	CreationStack() *[]*string
}

// The jsii proxy for IResolvable
type jsiiProxy_IResolvable struct {
	_ byte // padding
}

func (i *jsiiProxy_IResolvable) Resolve(context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"resolve",
		[]interface{}{context},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IResolvable) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IResolvable) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

// Current resolution context for tokens.
// Experimental.
type IResolveContext interface {
	// Use this postprocessor after the entire token structure has been resolved.
	// Experimental.
	RegisterPostProcessor(postProcessor IPostProcessor)
	// Resolve an inner object.
	// Experimental.
	Resolve(x interface{}, options *ResolveChangeContextOptions) interface{}
	// True when we are still preparing, false if we're rendering the final output.
	// Experimental.
	Preparing() *bool
	// The scope from which resolution has been initiated.
	// Experimental.
	Scope() IConstruct
}

// The jsii proxy for IResolveContext
type jsiiProxy_IResolveContext struct {
	_ byte // padding
}

func (i *jsiiProxy_IResolveContext) RegisterPostProcessor(postProcessor IPostProcessor) {
	_jsii_.InvokeVoid(
		i,
		"registerPostProcessor",
		[]interface{}{postProcessor},
	)
}

func (i *jsiiProxy_IResolveContext) Resolve(x interface{}, options *ResolveChangeContextOptions) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"resolve",
		[]interface{}{x, options},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IResolveContext) Preparing() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"preparing",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResolveContext) Scope() IConstruct {
	var returns IConstruct
	_jsii_.Get(
		j,
		"scope",
		&returns,
	)
	return returns
}

// Interface for the Resource construct.
// Experimental.
type IResource interface {
	IConstruct
	// The environment this resource belongs to.
	//
	// For resources that are created and managed by the CDK
	// (generally, those created by creating new class instances like Role, Bucket, etc.),
	// this is always the same as the environment of the stack they belong to;
	// however, for imported resources
	// (those obtained from static methods like fromRoleArn, fromBucketName, etc.),
	// that might be different than the stack they were imported into.
	// Experimental.
	Env() *ResourceEnvironment
	// The stack in which this resource is defined.
	// Experimental.
	Stack() Stack
}

// The jsii proxy for IResource
type jsiiProxy_IResource struct {
	jsiiProxy_IConstruct
}

func (j *jsiiProxy_IResource) Env() *ResourceEnvironment {
	var returns *ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// Interface for (stable) lazy untyped value producers.
// Experimental.
type IStableAnyProducer interface {
	// Produce the value.
	// Experimental.
	Produce() interface{}
}

// The jsii proxy for IStableAnyProducer
type jsiiProxy_IStableAnyProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_IStableAnyProducer) Produce() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"produce",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Interface for (stable) lazy list producers.
// Experimental.
type IStableListProducer interface {
	// Produce the list value.
	// Experimental.
	Produce() *[]*string
}

// The jsii proxy for IStableListProducer
type jsiiProxy_IStableListProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_IStableListProducer) Produce() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		i,
		"produce",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Interface for (stable) lazy number producers.
// Experimental.
type IStableNumberProducer interface {
	// Produce the number value.
	// Experimental.
	Produce() *float64
}

// The jsii proxy for IStableNumberProducer
type jsiiProxy_IStableNumberProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_IStableNumberProducer) Produce() *float64 {
	var returns *float64

	_jsii_.Invoke(
		i,
		"produce",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Interface for (stable) lazy string producers.
// Experimental.
type IStableStringProducer interface {
	// Produce the string value.
	// Experimental.
	Produce() *string
}

// The jsii proxy for IStableStringProducer
type jsiiProxy_IStableStringProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_IStableStringProducer) Produce() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"produce",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Encodes information how a certain Stack should be deployed.
// Experimental.
type IStackSynthesizer interface {
	// Register a Docker Image Asset.
	//
	// Returns the parameters that can be used to refer to the asset inside the template.
	// Experimental.
	AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation
	// Register a File Asset.
	//
	// Returns the parameters that can be used to refer to the asset inside the template.
	// Experimental.
	AddFileAsset(asset *FileAssetSource) *FileAssetLocation
	// Bind to the stack this environment is going to be used on.
	//
	// Must be called before any of the other methods are called.
	// Experimental.
	Bind(stack Stack)
	// Synthesize the associated stack to the session.
	// Experimental.
	Synthesize(session ISynthesisSession)
}

// The jsii proxy for IStackSynthesizer
type jsiiProxy_IStackSynthesizer struct {
	_ byte // padding
}

func (i *jsiiProxy_IStackSynthesizer) AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		i,
		"addDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStackSynthesizer) AddFileAsset(asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		i,
		"addFileAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IStackSynthesizer) Bind(stack Stack) {
	_jsii_.InvokeVoid(
		i,
		"bind",
		[]interface{}{stack},
	)
}

func (i *jsiiProxy_IStackSynthesizer) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		i,
		"synthesize",
		[]interface{}{session},
	)
}

// Interface for lazy string producers.
// Experimental.
type IStringProducer interface {
	// Produce the string value.
	// Experimental.
	Produce(context IResolveContext) *string
}

// The jsii proxy for IStringProducer
type jsiiProxy_IStringProducer struct {
	_ byte // padding
}

func (i *jsiiProxy_IStringProducer) Produce(context IResolveContext) *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"produce",
		[]interface{}{context},
		&returns,
	)

	return returns
}

// Represents a single session of synthesis.
//
// Passed into `Construct.synthesize()` methods.
// Experimental.
type ISynthesisSession interface {
	// Cloud assembly builder.
	// Experimental.
	Assembly() cxapi.CloudAssemblyBuilder
	// Cloud assembly builder.
	// Experimental.
	SetAssembly(a cxapi.CloudAssemblyBuilder)
	// The output directory for this synthesis session.
	// Experimental.
	Outdir() *string
	// The output directory for this synthesis session.
	// Experimental.
	SetOutdir(o *string)
	// Whether the stack should be validated after synthesis to check for error metadata.
	// Experimental.
	ValidateOnSynth() *bool
	// Whether the stack should be validated after synthesis to check for error metadata.
	// Experimental.
	SetValidateOnSynth(v *bool)
}

// The jsii proxy for ISynthesisSession
type jsiiProxy_ISynthesisSession struct {
	_ byte // padding
}

func (j *jsiiProxy_ISynthesisSession) Assembly() cxapi.CloudAssemblyBuilder {
	var returns cxapi.CloudAssemblyBuilder
	_jsii_.Get(
		j,
		"assembly",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISynthesisSession) SetAssembly(val cxapi.CloudAssemblyBuilder) {
	_jsii_.Set(
		j,
		"assembly",
		val,
	)
}

func (j *jsiiProxy_ISynthesisSession) Outdir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outdir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISynthesisSession) SetOutdir(val *string) {
	_jsii_.Set(
		j,
		"outdir",
		val,
	)
}

func (j *jsiiProxy_ISynthesisSession) ValidateOnSynth() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"validateOnSynth",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ISynthesisSession) SetValidateOnSynth(val *bool) {
	_jsii_.Set(
		j,
		"validateOnSynth",
		val,
	)
}

// Interface to implement tags.
// Experimental.
type ITaggable interface {
	// TagManager to set, remove and format tags.
	// Experimental.
	Tags() TagManager
}

// The jsii proxy for ITaggable
type jsiiProxy_ITaggable struct {
	_ byte // padding
}

func (j *jsiiProxy_ITaggable) Tags() TagManager {
	var returns TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

// CloudFormation template options for a stack.
// Experimental.
type ITemplateOptions interface {
	// Gets or sets the description of this stack.
	//
	// If provided, it will be included in the CloudFormation template's "Description" attribute.
	// Experimental.
	Description() *string
	// Gets or sets the description of this stack.
	//
	// If provided, it will be included in the CloudFormation template's "Description" attribute.
	// Experimental.
	SetDescription(d *string)
	// Metadata associated with the CloudFormation template.
	// Experimental.
	Metadata() *map[string]interface{}
	// Metadata associated with the CloudFormation template.
	// Experimental.
	SetMetadata(m *map[string]interface{})
	// Gets or sets the AWSTemplateFormatVersion field of the CloudFormation template.
	// Experimental.
	TemplateFormatVersion() *string
	// Gets or sets the AWSTemplateFormatVersion field of the CloudFormation template.
	// Experimental.
	SetTemplateFormatVersion(t *string)
	// Gets or sets the top-level template transform for this stack (e.g. "AWS::Serverless-2016-10-31").
	// Deprecated: use `transforms` instead.
	Transform() *string
	// Gets or sets the top-level template transform for this stack (e.g. "AWS::Serverless-2016-10-31").
	// Deprecated: use `transforms` instead.
	SetTransform(t *string)
	// Gets or sets the top-level template transform(s) for this stack (e.g. `["AWS::Serverless-2016-10-31"]`).
	// Experimental.
	Transforms() *[]*string
	// Gets or sets the top-level template transform(s) for this stack (e.g. `["AWS::Serverless-2016-10-31"]`).
	// Experimental.
	SetTransforms(t *[]*string)
}

// The jsii proxy for ITemplateOptions
type jsiiProxy_ITemplateOptions struct {
	_ byte // padding
}

func (j *jsiiProxy_ITemplateOptions) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITemplateOptions) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_ITemplateOptions) Metadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"metadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITemplateOptions) SetMetadata(val *map[string]interface{}) {
	_jsii_.Set(
		j,
		"metadata",
		val,
	)
}

func (j *jsiiProxy_ITemplateOptions) TemplateFormatVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateFormatVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITemplateOptions) SetTemplateFormatVersion(val *string) {
	_jsii_.Set(
		j,
		"templateFormatVersion",
		val,
	)
}

func (j *jsiiProxy_ITemplateOptions) Transform() *string {
	var returns *string
	_jsii_.Get(
		j,
		"transform",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITemplateOptions) SetTransform(val *string) {
	_jsii_.Set(
		j,
		"transform",
		val,
	)
}

func (j *jsiiProxy_ITemplateOptions) Transforms() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"transforms",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITemplateOptions) SetTransforms(val *[]*string) {
	_jsii_.Set(
		j,
		"transforms",
		val,
	)
}

// Interface to apply operation to tokens in a string.
//
// Interface so it can be exported via jsii.
// Experimental.
type ITokenMapper interface {
	// Replace a single token.
	// Experimental.
	MapToken(t IResolvable) interface{}
}

// The jsii proxy for ITokenMapper
type jsiiProxy_ITokenMapper struct {
	_ byte // padding
}

func (i *jsiiProxy_ITokenMapper) MapToken(t IResolvable) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"mapToken",
		[]interface{}{t},
		&returns,
	)

	return returns
}

// How to resolve tokens.
// Experimental.
type ITokenResolver interface {
	// Resolve a tokenized list.
	// Experimental.
	ResolveList(l *[]*string, context IResolveContext) interface{}
	// Resolve a string with at least one stringified token in it.
	//
	// (May use concatenation)
	// Experimental.
	ResolveString(s TokenizedStringFragments, context IResolveContext) interface{}
	// Resolve a single token.
	// Experimental.
	ResolveToken(t IResolvable, context IResolveContext, postProcessor IPostProcessor) interface{}
}

// The jsii proxy for ITokenResolver
type jsiiProxy_ITokenResolver struct {
	_ byte // padding
}

func (i *jsiiProxy_ITokenResolver) ResolveList(l *[]*string, context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"resolveList",
		[]interface{}{l, context},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITokenResolver) ResolveString(s TokenizedStringFragments, context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"resolveString",
		[]interface{}{s, context},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ITokenResolver) ResolveToken(t IResolvable, context IResolveContext, postProcessor IPostProcessor) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"resolveToken",
		[]interface{}{t, context, postProcessor},
		&returns,
	)

	return returns
}

// Determines the ignore behavior to use.
// Experimental.
type IgnoreMode string

const (
	IgnoreMode_GLOB IgnoreMode = "GLOB"
	IgnoreMode_GIT IgnoreMode = "GIT"
	IgnoreMode_DOCKER IgnoreMode = "DOCKER"
)

// Represents file path ignoring behavior.
// Experimental.
type IgnoreStrategy interface {
	Add(pattern *string)
	Ignores(absoluteFilePath *string) *bool
}

// The jsii proxy struct for IgnoreStrategy
type jsiiProxy_IgnoreStrategy struct {
	_ byte // padding
}

// Experimental.
func NewIgnoreStrategy_Override(i IgnoreStrategy) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.IgnoreStrategy",
		nil, // no parameters
		i,
	)
}

// Ignores file paths based on the [`.dockerignore specification`](https://docs.docker.com/engine/reference/builder/#dockerignore-file).
//
// Returns: `DockerIgnorePattern` associated with the given patterns.
// Experimental.
func IgnoreStrategy_Docker(absoluteRootPath *string, patterns *[]*string) DockerIgnoreStrategy {
	_init_.Initialize()

	var returns DockerIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.IgnoreStrategy",
		"docker",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Creates an IgnoreStrategy based on the `ignoreMode` and `exclude` in a `CopyOptions`.
//
// Returns: `IgnoreStrategy` based on the `CopyOptions`
// Experimental.
func IgnoreStrategy_FromCopyOptions(options *CopyOptions, absoluteRootPath *string) IgnoreStrategy {
	_init_.Initialize()

	var returns IgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.IgnoreStrategy",
		"fromCopyOptions",
		[]interface{}{options, absoluteRootPath},
		&returns,
	)

	return returns
}

// Ignores file paths based on the [`.gitignore specification`](https://git-scm.com/docs/gitignore).
//
// Returns: `GitIgnorePattern` associated with the given patterns.
// Experimental.
func IgnoreStrategy_Git(absoluteRootPath *string, patterns *[]*string) GitIgnoreStrategy {
	_init_.Initialize()

	var returns GitIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.IgnoreStrategy",
		"git",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Ignores file paths based on simple glob patterns.
//
// Returns: `GlobIgnorePattern` associated with the given patterns.
// Experimental.
func IgnoreStrategy_Glob(absoluteRootPath *string, patterns *[]*string) GlobIgnoreStrategy {
	_init_.Initialize()

	var returns GlobIgnoreStrategy

	_jsii_.StaticInvoke(
		"monocdk.IgnoreStrategy",
		"glob",
		[]interface{}{absoluteRootPath, patterns},
		&returns,
	)

	return returns
}

// Adds another pattern.
// Experimental.
func (i *jsiiProxy_IgnoreStrategy) Add(pattern *string) {
	_jsii_.InvokeVoid(
		i,
		"add",
		[]interface{}{pattern},
	)
}

// Determines whether a given file path should be ignored or not.
//
// Returns: `true` if the file should be ignored
// Experimental.
func (i *jsiiProxy_IgnoreStrategy) Ignores(absoluteFilePath *string) *bool {
	var returns *bool

	_jsii_.Invoke(
		i,
		"ignores",
		[]interface{}{absoluteFilePath},
		&returns,
	)

	return returns
}

// Token subclass that represents values intrinsic to the target document language.
//
// WARNING: this class should not be externally exposed, but is currently visible
// because of a limitation of jsii (https://github.com/aws/jsii/issues/524).
//
// This class will disappear in a future release and should not be used.
// Experimental.
type Intrinsic interface {
	IResolvable
	CreationStack() *[]*string
	NewError(message *string) interface{}
	Resolve(_context IResolveContext) interface{}
	ToJSON() interface{}
	ToString() *string
}

// The jsii proxy struct for Intrinsic
type jsiiProxy_Intrinsic struct {
	jsiiProxy_IResolvable
}

func (j *jsiiProxy_Intrinsic) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}


// Experimental.
func NewIntrinsic(value interface{}, options *IntrinsicProps) Intrinsic {
	_init_.Initialize()

	j := jsiiProxy_Intrinsic{}

	_jsii_.Create(
		"monocdk.Intrinsic",
		[]interface{}{value, options},
		&j,
	)

	return &j
}

// Experimental.
func NewIntrinsic_Override(i Intrinsic, value interface{}, options *IntrinsicProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.Intrinsic",
		[]interface{}{value, options},
		i,
	)
}

// Creates a throwable Error object that contains the token creation stack trace.
// Experimental.
func (i *jsiiProxy_Intrinsic) NewError(message *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"newError",
		[]interface{}{message},
		&returns,
	)

	return returns
}

// Produce the Token's value at resolution time.
// Experimental.
func (i *jsiiProxy_Intrinsic) Resolve(_context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

// Turn this Token into JSON.
//
// Called automatically when JSON.stringify() is called on a Token.
// Experimental.
func (i *jsiiProxy_Intrinsic) ToJSON() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		i,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Convert an instance of this Token to a string.
//
// This method will be called implicitly by language runtimes if the object
// is embedded into a string. We treat it the same as an explicit
// stringification.
// Experimental.
func (i *jsiiProxy_Intrinsic) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Customization properties for an Intrinsic token.
// Experimental.
type IntrinsicProps struct {
	// Capture the stack trace of where this token is created.
	// Experimental.
	StackTrace *bool `json:"stackTrace"`
}

// Lazily produce a value.
//
// Can be used to return a string, list or numeric value whose actual value
// will only be calculated later, during synthesis.
// Experimental.
type Lazy interface {
}

// The jsii proxy struct for Lazy
type jsiiProxy_Lazy struct {
	_ byte // padding
}

// Defer the one-time calculation of an arbitrarily typed value to synthesis time.
//
// Use this if you want to render an object to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
//
// The inner function will only be invoked one time and cannot depend on
// resolution context.
// Experimental.
func Lazy_Any(producer IStableAnyProducer, options *LazyAnyValueOptions) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"any",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the one-time calculation of an arbitrarily typed value to synthesis time.
//
// Use this if you want to render an object to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
// Deprecated: Use `Lazy.any()` or `Lazy.uncachedAny()` instead.
func Lazy_AnyValue(producer IAnyProducer, options *LazyAnyValueOptions) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"anyValue",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the one-time calculation of a list value to synthesis time.
//
// Use this if you want to render a list to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
//
// If you are simply looking to force a value to a `string[]` type and don't need
// the calculation to be deferred, use `Token.asList()` instead.
//
// The inner function will only be invoked once, and the resolved value
// cannot depend on the Stack the Token is used in.
// Experimental.
func Lazy_List(producer IStableListProducer, options *LazyListValueOptions) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"list",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the one-time calculation of a list value to synthesis time.
//
// Use this if you want to render a list to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
//
// If you are simply looking to force a value to a `string[]` type and don't need
// the calculation to be deferred, use `Token.asList()` instead.
// Deprecated: Use `Lazy.list()` or `Lazy.uncachedList()` instead.
func Lazy_ListValue(producer IListProducer, options *LazyListValueOptions) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"listValue",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the one-time calculation of a number value to synthesis time.
//
// Use this if you want to render a number to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
//
// If you are simply looking to force a value to a `number` type and don't need
// the calculation to be deferred, use `Token.asNumber()` instead.
//
// The inner function will only be invoked once, and the resolved value
// cannot depend on the Stack the Token is used in.
// Experimental.
func Lazy_Number(producer IStableNumberProducer) *float64 {
	_init_.Initialize()

	var returns *float64

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"number",
		[]interface{}{producer},
		&returns,
	)

	return returns
}

// Defer the one-time calculation of a number value to synthesis time.
//
// Use this if you want to render a number to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
//
// If you are simply looking to force a value to a `number` type and don't need
// the calculation to be deferred, use `Token.asNumber()` instead.
// Deprecated: Use `Lazy.number()` or `Lazy.uncachedNumber()` instead.
func Lazy_NumberValue(producer INumberProducer) *float64 {
	_init_.Initialize()

	var returns *float64

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"numberValue",
		[]interface{}{producer},
		&returns,
	)

	return returns
}

// Defer the one-time calculation of a string value to synthesis time.
//
// Use this if you want to render a string to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
//
// If you are simply looking to force a value to a `string` type and don't need
// the calculation to be deferred, use `Token.asString()` instead.
//
// The inner function will only be invoked once, and the resolved value
// cannot depend on the Stack the Token is used in.
// Experimental.
func Lazy_String(producer IStableStringProducer, options *LazyStringValueOptions) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"string",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the calculation of a string value to synthesis time.
//
// Use this if you want to render a string to a template whose actual value depends on
// some state mutation that may happen after the construct has been created.
//
// If you are simply looking to force a value to a `string` type and don't need
// the calculation to be deferred, use `Token.asString()` instead.
// Deprecated: Use `Lazy.string()` or `Lazy.uncachedString()` instead.
func Lazy_StringValue(producer IStringProducer, options *LazyStringValueOptions) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"stringValue",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the calculation of an untyped value to synthesis time.
//
// Use of this function is not recommended; unless you know you need it for sure, you
// probably don't. Use `Lazy.any()` instead.
//
// The inner function may be invoked multiple times during synthesis. You
// should only use this method if the returned value depends on variables
// that may change during the Aspect application phase of synthesis, or if
// the value depends on the Stack the value is being used in. Both of these
// cases are rare, and only ever occur for AWS Construct Library authors.
// Experimental.
func Lazy_UncachedAny(producer IAnyProducer, options *LazyAnyValueOptions) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"uncachedAny",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the calculation of a list value to synthesis time.
//
// Use of this function is not recommended; unless you know you need it for sure, you
// probably don't. Use `Lazy.list()` instead.
//
// The inner function may be invoked multiple times during synthesis. You
// should only use this method if the returned value depends on variables
// that may change during the Aspect application phase of synthesis, or if
// the value depends on the Stack the value is being used in. Both of these
// cases are rare, and only ever occur for AWS Construct Library authors.
// Experimental.
func Lazy_UncachedList(producer IListProducer, options *LazyListValueOptions) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"uncachedList",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Defer the calculation of a number value to synthesis time.
//
// Use of this function is not recommended; unless you know you need it for sure, you
// probably don't. Use `Lazy.number()` instead.
//
// The inner function may be invoked multiple times during synthesis. You
// should only use this method if the returned value depends on variables
// that may change during the Aspect application phase of synthesis, or if
// the value depends on the Stack the value is being used in. Both of these
// cases are rare, and only ever occur for AWS Construct Library authors.
// Experimental.
func Lazy_UncachedNumber(producer INumberProducer) *float64 {
	_init_.Initialize()

	var returns *float64

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"uncachedNumber",
		[]interface{}{producer},
		&returns,
	)

	return returns
}

// Defer the calculation of a string value to synthesis time.
//
// Use of this function is not recommended; unless you know you need it for sure, you
// probably don't. Use `Lazy.string()` instead.
//
// The inner function may be invoked multiple times during synthesis. You
// should only use this method if the returned value depends on variables
// that may change during the Aspect application phase of synthesis, or if
// the value depends on the Stack the value is being used in. Both of these
// cases are rare, and only ever occur for AWS Construct Library authors.
// Experimental.
func Lazy_UncachedString(producer IStringProducer, options *LazyStringValueOptions) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Lazy",
		"uncachedString",
		[]interface{}{producer, options},
		&returns,
	)

	return returns
}

// Options for creating lazy untyped tokens.
// Experimental.
type LazyAnyValueOptions struct {
	// Use the given name as a display hint.
	// Experimental.
	DisplayHint *string `json:"displayHint"`
	// If the produced value is an array and it is empty, return 'undefined' instead.
	// Experimental.
	OmitEmptyArray *bool `json:"omitEmptyArray"`
}

// Options for creating a lazy list token.
// Experimental.
type LazyListValueOptions struct {
	// Use the given name as a display hint.
	// Experimental.
	DisplayHint *string `json:"displayHint"`
	// If the produced list is empty, return 'undefined' instead.
	// Experimental.
	OmitEmpty *bool `json:"omitEmpty"`
}

// Options for creating a lazy string token.
// Experimental.
type LazyStringValueOptions struct {
	// Use the given name as a display hint.
	// Experimental.
	DisplayHint *string `json:"displayHint"`
}

// Use the original deployment environment.
//
// This deployment environment is restricted in cross-environment deployments,
// CI/CD deployments, and will use up CloudFormation parameters in your template.
//
// This is the only StackSynthesizer that supports customizing asset behavior
// by overriding `Stack.addFileAsset()` and `Stack.addDockerImageAsset()`.
// Experimental.
type LegacyStackSynthesizer interface {
	StackSynthesizer
	AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation
	AddFileAsset(asset *FileAssetSource) *FileAssetLocation
	Bind(stack Stack)
	EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions)
	Synthesize(session ISynthesisSession)
	SynthesizeStackTemplate(stack Stack, session ISynthesisSession)
}

// The jsii proxy struct for LegacyStackSynthesizer
type jsiiProxy_LegacyStackSynthesizer struct {
	jsiiProxy_StackSynthesizer
}

// Experimental.
func NewLegacyStackSynthesizer() LegacyStackSynthesizer {
	_init_.Initialize()

	j := jsiiProxy_LegacyStackSynthesizer{}

	_jsii_.Create(
		"monocdk.LegacyStackSynthesizer",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewLegacyStackSynthesizer_Override(l LegacyStackSynthesizer) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.LegacyStackSynthesizer",
		nil, // no parameters
		l,
	)
}

// Register a Docker Image Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (l *jsiiProxy_LegacyStackSynthesizer) AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		l,
		"addDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Register a File Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (l *jsiiProxy_LegacyStackSynthesizer) AddFileAsset(asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		l,
		"addFileAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Bind to the stack this environment is going to be used on.
//
// Must be called before any of the other methods are called.
// Experimental.
func (l *jsiiProxy_LegacyStackSynthesizer) Bind(stack Stack) {
	_jsii_.InvokeVoid(
		l,
		"bind",
		[]interface{}{stack},
	)
}

// Write the stack artifact to the session.
//
// Use default settings to add a CloudFormationStackArtifact artifact to
// the given synthesis session.
// Experimental.
func (l *jsiiProxy_LegacyStackSynthesizer) EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions) {
	_jsii_.InvokeVoid(
		l,
		"emitStackArtifact",
		[]interface{}{stack, session, options},
	)
}

// Synthesize the associated stack to the session.
// Experimental.
func (l *jsiiProxy_LegacyStackSynthesizer) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Have the stack write out its template.
// Experimental.
func (l *jsiiProxy_LegacyStackSynthesizer) SynthesizeStackTemplate(stack Stack, session ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesizeStackTemplate",
		[]interface{}{stack, session},
	)
}

// Functions for devising unique names for constructs.
//
// For example, those can be
// used to allocate unique physical names for resources.
// Experimental.
type Names interface {
}

// The jsii proxy struct for Names
type jsiiProxy_Names struct {
	_ byte // padding
}

// Returns a CloudFormation-compatible unique identifier for a construct based on its path.
//
// The identifier includes a human readable portion rendered
// from the path components and a hash suffix.
//
// TODO (v2): replace with API to use `constructs.Node`.
//
// Returns: a unique id based on the construct path
// Experimental.
func Names_NodeUniqueId(node ConstructNode) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Names",
		"nodeUniqueId",
		[]interface{}{node},
		&returns,
	)

	return returns
}

// Returns a CloudFormation-compatible unique identifier for a construct based on its path.
//
// The identifier includes a human readable portion rendered
// from the path components and a hash suffix.
//
// Returns: a unique id based on the construct path
// Experimental.
func Names_UniqueId(construct constructs.Construct) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Names",
		"uniqueId",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// A CloudFormation nested stack.
//
// When you apply template changes to update a top-level stack, CloudFormation
// updates the top-level stack and initiates an update to its nested stacks.
// CloudFormation updates the resources of modified nested stacks, but does not
// update the resources of unmodified nested stacks.
//
// Furthermore, this stack will not be treated as an independent deployment
// artifact (won't be listed in "cdk list" or deployable through "cdk deploy"),
// but rather only synthesized as a template and uploaded as an asset to S3.
//
// Cross references of resource attributes between the parent stack and the
// nested stack will automatically be translated to stack parameters and
// outputs.
// Experimental.
type NestedStack interface {
	Stack
	Account() *string
	ArtifactId() *string
	AvailabilityZones() *[]*string
	Dependencies() *[]Stack
	Environment() *string
	Nested() *bool
	NestedStackParent() Stack
	NestedStackResource() CfnResource
	Node() ConstructNode
	NotificationArns() *[]*string
	ParentStack() Stack
	Partition() *string
	Region() *string
	StackId() *string
	StackName() *string
	Synthesizer() IStackSynthesizer
	Tags() TagManager
	TemplateFile() *string
	TemplateOptions() ITemplateOptions
	TerminationProtection() *bool
	UrlSuffix() *string
	AddDependency(target Stack, reason *string)
	AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation
	AddFileAsset(asset *FileAssetSource) *FileAssetLocation
	AddTransform(transform *string)
	AllocateLogicalId(cfnElement CfnElement) *string
	ExportValue(exportedValue interface{}, options *ExportValueOptions) *string
	FormatArn(components *ArnComponents) *string
	GetLogicalId(element CfnElement) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	ParseArn(arn *string, sepIfToken *string, hasName *bool) *ArnComponents
	Prepare()
	PrepareCrossReference(_sourceStack Stack, reference Reference) IResolvable
	RenameLogicalId(oldId *string, newId *string)
	ReportMissingContext(report *cxapi.MissingContext)
	ReportMissingContextKey(report *cloudassemblyschema.MissingContext)
	Resolve(obj interface{}) interface{}
	SetParameter(name *string, value *string)
	SplitArn(arn *string, arnFormat ArnFormat) *ArnComponents
	Synthesize(session ISynthesisSession)
	ToJsonString(obj interface{}, space *float64) *string
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for NestedStack
type jsiiProxy_NestedStack struct {
	jsiiProxy_Stack
}

func (j *jsiiProxy_NestedStack) Account() *string {
	var returns *string
	_jsii_.Get(
		j,
		"account",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) ArtifactId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"artifactId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) AvailabilityZones() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"availabilityZones",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Dependencies() *[]Stack {
	var returns *[]Stack
	_jsii_.Get(
		j,
		"dependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Environment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"environment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Nested() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"nested",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) NestedStackParent() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"nestedStackParent",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) NestedStackResource() CfnResource {
	var returns CfnResource
	_jsii_.Get(
		j,
		"nestedStackResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) NotificationArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"notificationArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) ParentStack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"parentStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Partition() *string {
	var returns *string
	_jsii_.Get(
		j,
		"partition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) StackId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) StackName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Synthesizer() IStackSynthesizer {
	var returns IStackSynthesizer
	_jsii_.Get(
		j,
		"synthesizer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) Tags() TagManager {
	var returns TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) TemplateFile() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateFile",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) TemplateOptions() ITemplateOptions {
	var returns ITemplateOptions
	_jsii_.Get(
		j,
		"templateOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) TerminationProtection() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"terminationProtection",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NestedStack) UrlSuffix() *string {
	var returns *string
	_jsii_.Get(
		j,
		"urlSuffix",
		&returns,
	)
	return returns
}


// Experimental.
func NewNestedStack(scope constructs.Construct, id *string, props *NestedStackProps) NestedStack {
	_init_.Initialize()

	j := jsiiProxy_NestedStack{}

	_jsii_.Create(
		"monocdk.NestedStack",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewNestedStack_Override(n NestedStack, scope constructs.Construct, id *string, props *NestedStackProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.NestedStack",
		[]interface{}{scope, id, props},
		n,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func NestedStack_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.NestedStack",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Checks if `x` is an object of type `NestedStack`.
// Experimental.
func NestedStack_IsNestedStack(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.NestedStack",
		"isNestedStack",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Stack.
//
// We do attribute detection since we can't reliably use 'instanceof'.
// Experimental.
func NestedStack_IsStack(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.NestedStack",
		"isStack",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Looks up the first stack scope in which `construct` is defined.
//
// Fails if there is no stack up the tree.
// Experimental.
func NestedStack_Of(construct constructs.IConstruct) Stack {
	_init_.Initialize()

	var returns Stack

	_jsii_.StaticInvoke(
		"monocdk.NestedStack",
		"of",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a dependency between this stack and another stack.
//
// This can be used to define dependencies between any two stacks within an
// app, and also supports nested stacks.
// Experimental.
func (n *jsiiProxy_NestedStack) AddDependency(target Stack, reason *string) {
	_jsii_.InvokeVoid(
		n,
		"addDependency",
		[]interface{}{target, reason},
	)
}

// Register a docker image asset on this Stack.
// Deprecated: Use `stack.synthesizer.addDockerImageAsset()` if you are calling,
// and a different `IStackSynthesizer` class if you are implementing.
func (n *jsiiProxy_NestedStack) AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		n,
		"addDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Register a file asset on this Stack.
// Deprecated: Use `stack.synthesizer.addFileAsset()` if you are calling,
// and a different IStackSynthesizer class if you are implementing.
func (n *jsiiProxy_NestedStack) AddFileAsset(asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		n,
		"addFileAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Add a Transform to this stack. A Transform is a macro that AWS CloudFormation uses to process your template.
//
// Duplicate values are removed when stack is synthesized.
//
// TODO: EXAMPLE
//
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/transform-section-structure.html
//
// Experimental.
func (n *jsiiProxy_NestedStack) AddTransform(transform *string) {
	_jsii_.InvokeVoid(
		n,
		"addTransform",
		[]interface{}{transform},
	)
}

// Returns the naming scheme used to allocate logical IDs.
//
// By default, uses
// the `HashedAddressingScheme` but this method can be overridden to customize
// this behavior.
//
// In order to make sure logical IDs are unique and stable, we hash the resource
// construct tree path (i.e. toplevel/secondlevel/.../myresource) and add it as
// a suffix to the path components joined without a separator (CloudFormation
// IDs only allow alphanumeric characters).
//
// The result will be:
//
//    <path.join('')><md5(path.join('/')>
//      "human"      "hash"
//
// If the "human" part of the ID exceeds 240 characters, we simply trim it so
// the total ID doesn't exceed CloudFormation's 255 character limit.
//
// We only take 8 characters from the md5 hash (0.000005 chance of collision).
//
// Special cases:
//
// - If the path only contains a single component (i.e. it's a top-level
//    resource), we won't add the hash to it. The hash is not needed for
//    disamiguation and also, it allows for a more straightforward migration an
//    existing CloudFormation template to a CDK stack without logical ID changes
//    (or renames).
// - For aesthetic reasons, if the last components of the path are the same
//    (i.e. `L1/L2/Pipeline/Pipeline`), they will be de-duplicated to make the
//    resulting human portion of the ID more pleasing: `L1L2Pipeline<HASH>`
//    instead of `L1L2PipelinePipeline<HASH>`
// - If a component is named "Default" it will be omitted from the path. This
//    allows refactoring higher level abstractions around constructs without affecting
//    the IDs of already deployed resources.
// - If a component is named "Resource" it will be omitted from the user-visible
//    path, but included in the hash. This reduces visual noise in the human readable
//    part of the identifier.
// Experimental.
func (n *jsiiProxy_NestedStack) AllocateLogicalId(cfnElement CfnElement) *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"allocateLogicalId",
		[]interface{}{cfnElement},
		&returns,
	)

	return returns
}

// Create a CloudFormation Export for a value.
//
// Returns a string representing the corresponding `Fn.importValue()`
// expression for this Export. You can control the name for the export by
// passing the `name` option.
//
// If you don't supply a value for `name`, the value you're exporting must be
// a Resource attribute (for example: `bucket.bucketName`) and it will be
// given the same name as the automatic cross-stack reference that would be created
// if you used the attribute in another Stack.
//
// One of the uses for this method is to *remove* the relationship between
// two Stacks established by automatic cross-stack references. It will
// temporarily ensure that the CloudFormation Export still exists while you
// remove the reference from the consuming stack. After that, you can remove
// the resource and the manual export.
//
// ## Example
//
// Here is how the process works. Let's say there are two stacks,
// `producerStack` and `consumerStack`, and `producerStack` has a bucket
// called `bucket`, which is referenced by `consumerStack` (perhaps because
// an AWS Lambda Function writes into it, or something like that).
//
// It is not safe to remove `producerStack.bucket` because as the bucket is being
// deleted, `consumerStack` might still be using it.
//
// Instead, the process takes two deployments:
//
// ### Deployment 1: break the relationship
//
// - Make sure `consumerStack` no longer references `bucket.bucketName` (maybe the consumer
//    stack now uses its own bucket, or it writes to an AWS DynamoDB table, or maybe you just
//    remove the Lambda Function altogether).
// - In the `ProducerStack` class, call `this.exportValue(this.bucket.bucketName)`. This
//    will make sure the CloudFormation Export continues to exist while the relationship
//    between the two stacks is being broken.
// - Deploy (this will effectively only change the `consumerStack`, but it's safe to deploy both).
//
// ### Deployment 2: remove the bucket resource
//
// - You are now free to remove the `bucket` resource from `producerStack`.
// - Don't forget to remove the `exportValue()` call as well.
// - Deploy again (this time only the `producerStack` will be changed -- the bucket will be deleted).
// Experimental.
func (n *jsiiProxy_NestedStack) ExportValue(exportedValue interface{}, options *ExportValueOptions) *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"exportValue",
		[]interface{}{exportedValue, options},
		&returns,
	)

	return returns
}

// Creates an ARN from components.
//
// If `partition`, `region` or `account` are not specified, the stack's
// partition, region and account will be used.
//
// If any component is the empty string, an empty string will be inserted
// into the generated ARN at the location that component corresponds to.
//
// The ARN will be formatted as follows:
//
//    arn:{partition}:{service}:{region}:{account}:{resource}{sep}}{resource-name}
//
// The required ARN pieces that are omitted will be taken from the stack that
// the 'scope' is attached to. If all ARN pieces are supplied, the supplied scope
// can be 'undefined'.
// Experimental.
func (n *jsiiProxy_NestedStack) FormatArn(components *ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"formatArn",
		[]interface{}{components},
		&returns,
	)

	return returns
}

// Allocates a stack-unique CloudFormation-compatible logical identity for a specific resource.
//
// This method is called when a `CfnElement` is created and used to render the
// initial logical identity of resources. Logical ID renames are applied at
// this stage.
//
// This method uses the protected method `allocateLogicalId` to render the
// logical ID for an element. To modify the naming scheme, extend the `Stack`
// class and override this method.
// Experimental.
func (n *jsiiProxy_NestedStack) GetLogicalId(element CfnElement) *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"getLogicalId",
		[]interface{}{element},
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
func (n *jsiiProxy_NestedStack) OnPrepare() {
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
func (n *jsiiProxy_NestedStack) OnSynthesize(session constructs.ISynthesisSession) {
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
func (n *jsiiProxy_NestedStack) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Given an ARN, parses it and returns components.
//
// IF THE ARN IS A CONCRETE STRING...
//
// ...it will be parsed and validated. The separator (`sep`) will be set to '/'
// if the 6th component includes a '/', in which case, `resource` will be set
// to the value before the '/' and `resourceName` will be the rest. In case
// there is no '/', `resource` will be set to the 6th components and
// `resourceName` will be set to the rest of the string.
//
// IF THE ARN IS A TOKEN...
//
// ...it cannot be validated, since we don't have the actual value yet at the
// time of this function call. You will have to supply `sepIfToken` and
// whether or not ARNs of the expected format usually have resource names
// in order to parse it properly. The resulting `ArnComponents` object will
// contain tokens for the subexpressions of the ARN, not string literals.
//
// If the resource name could possibly contain the separator char, the actual
// resource name cannot be properly parsed. This only occurs if the separator
// char is '/', and happens for example for S3 object ARNs, IAM Role ARNs,
// IAM OIDC Provider ARNs, etc. To properly extract the resource name from a
// Tokenized ARN, you must know the resource type and call
// `Arn.extractResourceName`.
//
// Returns: an ArnComponents object which allows access to the various
// components of the ARN.
// Deprecated: use splitArn instead
func (n *jsiiProxy_NestedStack) ParseArn(arn *string, sepIfToken *string, hasName *bool) *ArnComponents {
	var returns *ArnComponents

	_jsii_.Invoke(
		n,
		"parseArn",
		[]interface{}{arn, sepIfToken, hasName},
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
func (n *jsiiProxy_NestedStack) Prepare() {
	_jsii_.InvokeVoid(
		n,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated.
//
// Returns: reference itself without any change
// See: https://github.com/aws/aws-cdk/pull/7187
//
// Deprecated: cross reference handling has been moved to `App.prepare()`.
func (n *jsiiProxy_NestedStack) PrepareCrossReference(_sourceStack Stack, reference Reference) IResolvable {
	var returns IResolvable

	_jsii_.Invoke(
		n,
		"prepareCrossReference",
		[]interface{}{_sourceStack, reference},
		&returns,
	)

	return returns
}

// Rename a generated logical identities.
//
// To modify the naming scheme strategy, extend the `Stack` class and
// override the `allocateLogicalId` method.
// Experimental.
func (n *jsiiProxy_NestedStack) RenameLogicalId(oldId *string, newId *string) {
	_jsii_.InvokeVoid(
		n,
		"renameLogicalId",
		[]interface{}{oldId, newId},
	)
}

// DEPRECATED.
// Deprecated: use `reportMissingContextKey()`
func (n *jsiiProxy_NestedStack) ReportMissingContext(report *cxapi.MissingContext) {
	_jsii_.InvokeVoid(
		n,
		"reportMissingContext",
		[]interface{}{report},
	)
}

// Indicate that a context key was expected.
//
// Contains instructions which will be emitted into the cloud assembly on how
// the key should be supplied.
// Experimental.
func (n *jsiiProxy_NestedStack) ReportMissingContextKey(report *cloudassemblyschema.MissingContext) {
	_jsii_.InvokeVoid(
		n,
		"reportMissingContextKey",
		[]interface{}{report},
	)
}

// Resolve a tokenized value in the context of the current stack.
// Experimental.
func (n *jsiiProxy_NestedStack) Resolve(obj interface{}) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		n,
		"resolve",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Assign a value to one of the nested stack parameters.
// Experimental.
func (n *jsiiProxy_NestedStack) SetParameter(name *string, value *string) {
	_jsii_.InvokeVoid(
		n,
		"setParameter",
		[]interface{}{name, value},
	)
}

// Splits the provided ARN into its components.
//
// Works both if 'arn' is a string like 'arn:aws:s3:::bucket',
// and a Token representing a dynamic CloudFormation expression
// (in which case the returned components will also be dynamic CloudFormation expressions,
// encoded as Tokens).
// Experimental.
func (n *jsiiProxy_NestedStack) SplitArn(arn *string, arnFormat ArnFormat) *ArnComponents {
	var returns *ArnComponents

	_jsii_.Invoke(
		n,
		"splitArn",
		[]interface{}{arn, arnFormat},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NestedStack) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
		"synthesize",
		[]interface{}{session},
	)
}

// Convert an object, potentially containing tokens, to a JSON string.
// Experimental.
func (n *jsiiProxy_NestedStack) ToJsonString(obj interface{}, space *float64) *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"toJsonString",
		[]interface{}{obj, space},
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (n *jsiiProxy_NestedStack) ToString() *string {
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
func (n *jsiiProxy_NestedStack) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Initialization props for the `NestedStack` construct.
// Experimental.
type NestedStackProps struct {
	// The Simple Notification Service (SNS) topics to publish stack related events.
	// Experimental.
	NotificationArns *[]*string `json:"notificationArns"`
	// The set value pairs that represent the parameters passed to CloudFormation when this nested stack is created.
	//
	// Each parameter has a name corresponding
	// to a parameter defined in the embedded template and a value representing
	// the value that you want to set for the parameter.
	//
	// The nested stack construct will automatically synthesize parameters in order
	// to bind references from the parent stack(s) into the nested stack.
	// Experimental.
	Parameters *map[string]*string `json:"parameters"`
	// Policy to apply when the nested stack is removed.
	//
	// The default is `Destroy`, because all Removal Policies of resources inside the
	// Nested Stack should already have been set correctly. You normally should
	// not need to set this value.
	// Experimental.
	RemovalPolicy RemovalPolicy `json:"removalPolicy"`
	// The length of time that CloudFormation waits for the nested stack to reach the CREATE_COMPLETE state.
	//
	// When CloudFormation detects that the nested stack has reached the
	// CREATE_COMPLETE state, it marks the nested stack resource as
	// CREATE_COMPLETE in the parent stack and resumes creating the parent stack.
	// If the timeout period expires before the nested stack reaches
	// CREATE_COMPLETE, CloudFormation marks the nested stack as failed and rolls
	// back both the nested stack and parent stack.
	// Experimental.
	Timeout Duration `json:"timeout"`
}

// Deployment environment for a nested stack.
//
// Interoperates with the StackSynthesizer of the parent stack.
// Experimental.
type NestedStackSynthesizer interface {
	StackSynthesizer
	AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation
	AddFileAsset(asset *FileAssetSource) *FileAssetLocation
	Bind(stack Stack)
	EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions)
	Synthesize(session ISynthesisSession)
	SynthesizeStackTemplate(stack Stack, session ISynthesisSession)
}

// The jsii proxy struct for NestedStackSynthesizer
type jsiiProxy_NestedStackSynthesizer struct {
	jsiiProxy_StackSynthesizer
}

// Experimental.
func NewNestedStackSynthesizer(parentDeployment IStackSynthesizer) NestedStackSynthesizer {
	_init_.Initialize()

	j := jsiiProxy_NestedStackSynthesizer{}

	_jsii_.Create(
		"monocdk.NestedStackSynthesizer",
		[]interface{}{parentDeployment},
		&j,
	)

	return &j
}

// Experimental.
func NewNestedStackSynthesizer_Override(n NestedStackSynthesizer, parentDeployment IStackSynthesizer) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.NestedStackSynthesizer",
		[]interface{}{parentDeployment},
		n,
	)
}

// Register a Docker Image Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (n *jsiiProxy_NestedStackSynthesizer) AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		n,
		"addDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Register a File Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (n *jsiiProxy_NestedStackSynthesizer) AddFileAsset(asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		n,
		"addFileAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Bind to the stack this environment is going to be used on.
//
// Must be called before any of the other methods are called.
// Experimental.
func (n *jsiiProxy_NestedStackSynthesizer) Bind(stack Stack) {
	_jsii_.InvokeVoid(
		n,
		"bind",
		[]interface{}{stack},
	)
}

// Write the stack artifact to the session.
//
// Use default settings to add a CloudFormationStackArtifact artifact to
// the given synthesis session.
// Experimental.
func (n *jsiiProxy_NestedStackSynthesizer) EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions) {
	_jsii_.InvokeVoid(
		n,
		"emitStackArtifact",
		[]interface{}{stack, session, options},
	)
}

// Synthesize the associated stack to the session.
// Experimental.
func (n *jsiiProxy_NestedStackSynthesizer) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
		"synthesize",
		[]interface{}{session},
	)
}

// Have the stack write out its template.
// Experimental.
func (n *jsiiProxy_NestedStackSynthesizer) SynthesizeStackTemplate(stack Stack, session ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
		"synthesizeStackTemplate",
		[]interface{}{stack, session},
	)
}

// Includes special markers for automatic generation of physical names.
// Experimental.
type PhysicalName interface {
}

// The jsii proxy struct for PhysicalName
type jsiiProxy_PhysicalName struct {
	_ byte // padding
}

func PhysicalName_GENERATE_IF_NEEDED() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.PhysicalName",
		"GENERATE_IF_NEEDED",
		&returns,
	)
	return returns
}

// An intrinsic Token that represents a reference to a construct.
//
// References are recorded.
// Experimental.
type Reference interface {
	Intrinsic
	CreationStack() *[]*string
	DisplayName() *string
	Target() IConstruct
	NewError(message *string) interface{}
	Resolve(_context IResolveContext) interface{}
	ToJSON() interface{}
	ToString() *string
}

// The jsii proxy struct for Reference
type jsiiProxy_Reference struct {
	jsiiProxy_Intrinsic
}

func (j *jsiiProxy_Reference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Reference) DisplayName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"displayName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Reference) Target() IConstruct {
	var returns IConstruct
	_jsii_.Get(
		j,
		"target",
		&returns,
	)
	return returns
}


// Experimental.
func NewReference_Override(r Reference, value interface{}, target IConstruct, displayName *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.Reference",
		[]interface{}{value, target, displayName},
		r,
	)
}

// Check whether this is actually a Reference.
// Experimental.
func Reference_IsReference(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Reference",
		"isReference",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Creates a throwable Error object that contains the token creation stack trace.
// Experimental.
func (r *jsiiProxy_Reference) NewError(message *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		r,
		"newError",
		[]interface{}{message},
		&returns,
	)

	return returns
}

// Produce the Token's value at resolution time.
// Experimental.
func (r *jsiiProxy_Reference) Resolve(_context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		r,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

// Turn this Token into JSON.
//
// Called automatically when JSON.stringify() is called on a Token.
// Experimental.
func (r *jsiiProxy_Reference) ToJSON() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		r,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Convert an instance of this Token to a string.
//
// This method will be called implicitly by language runtimes if the object
// is embedded into a string. We treat it the same as an explicit
// stringification.
// Experimental.
func (r *jsiiProxy_Reference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Possible values for a resource's Removal Policy.
//
// The removal policy controls what happens to the resource if it stops being
// managed by CloudFormation. This can happen in one of three situations:
//
// - The resource is removed from the template, so CloudFormation stops managing it;
// - A change to the resource is made that requires it to be replaced, so CloudFormation stops
//    managing it;
// - The stack is deleted, so CloudFormation stops managing all resources in it.
//
// The Removal Policy applies to all above cases.
//
// Many stateful resources in the AWS Construct Library will accept a
// `removalPolicy` as a property, typically defaulting it to `RETAIN`.
//
// If the AWS Construct Library resource does not accept a `removalPolicy`
// argument, you can always configure it by using the escape hatch mechanism,
// as shown in the following example:
//
// ```ts
// const cfnBucket = bucket.node.findChild('Resource') as cdk.CfnResource;
// cfnBucket.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY);
// ```
// Experimental.
type RemovalPolicy string

const (
	RemovalPolicy_DESTROY RemovalPolicy = "DESTROY"
	RemovalPolicy_RETAIN RemovalPolicy = "RETAIN"
	RemovalPolicy_SNAPSHOT RemovalPolicy = "SNAPSHOT"
)

// Experimental.
type RemovalPolicyOptions struct {
	// Apply the same deletion policy to the resource's "UpdateReplacePolicy".
	// Experimental.
	ApplyToUpdateReplacePolicy *bool `json:"applyToUpdateReplacePolicy"`
	// The default policy to apply in case the removal policy is not defined.
	// Experimental.
	Default RemovalPolicy `json:"default"`
}

// The RemoveTag Aspect will handle removing tags from this node and children.
// Experimental.
type RemoveTag interface {
	IAspect
	Key() *string
	Props() *TagProps
	ApplyTag(resource ITaggable)
	Visit(construct IConstruct)
}

// The jsii proxy struct for RemoveTag
type jsiiProxy_RemoveTag struct {
	jsiiProxy_IAspect
}

func (j *jsiiProxy_RemoveTag) Key() *string {
	var returns *string
	_jsii_.Get(
		j,
		"key",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RemoveTag) Props() *TagProps {
	var returns *TagProps
	_jsii_.Get(
		j,
		"props",
		&returns,
	)
	return returns
}


// Experimental.
func NewRemoveTag(key *string, props *TagProps) RemoveTag {
	_init_.Initialize()

	j := jsiiProxy_RemoveTag{}

	_jsii_.Create(
		"monocdk.RemoveTag",
		[]interface{}{key, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRemoveTag_Override(r RemoveTag, key *string, props *TagProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.RemoveTag",
		[]interface{}{key, props},
		r,
	)
}

// Experimental.
func (r *jsiiProxy_RemoveTag) ApplyTag(resource ITaggable) {
	_jsii_.InvokeVoid(
		r,
		"applyTag",
		[]interface{}{resource},
	)
}

// All aspects can visit an IConstruct.
// Experimental.
func (r *jsiiProxy_RemoveTag) Visit(construct IConstruct) {
	_jsii_.InvokeVoid(
		r,
		"visit",
		[]interface{}{construct},
	)
}

// Options that can be changed while doing a recursive resolve.
// Experimental.
type ResolveChangeContextOptions struct {
	// Change the 'allowIntrinsicKeys' option.
	// Experimental.
	AllowIntrinsicKeys *bool `json:"allowIntrinsicKeys"`
}

// Options to the resolve() operation.
//
// NOT the same as the ResolveContext; ResolveContext is exposed to Token
// implementors and resolution hooks, whereas this struct is just to bundle
// a number of things that would otherwise be arguments to resolve() in a
// readable way.
// Experimental.
type ResolveOptions struct {
	// The resolver to apply to any resolvable tokens found.
	// Experimental.
	Resolver ITokenResolver `json:"resolver"`
	// The scope from which resolution is performed.
	// Experimental.
	Scope constructs.IConstruct `json:"scope"`
	// Whether the resolution is being executed during the prepare phase or not.
	// Experimental.
	Preparing *bool `json:"preparing"`
	// Whether to remove undefined elements from arrays and objects when resolving.
	// Experimental.
	RemoveEmpty *bool `json:"removeEmpty"`
}

// A construct which represents an AWS resource.
// Experimental.
type Resource interface {
	Construct
	IResource
	Env() *ResourceEnvironment
	Node() ConstructNode
	PhysicalName() *string
	Stack() Stack
	ApplyRemovalPolicy(policy RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Resource
type jsiiProxy_Resource struct {
	jsiiProxy_Construct
	jsiiProxy_IResource
}

func (j *jsiiProxy_Resource) Env() *ResourceEnvironment {
	var returns *ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) Stack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewResource_Override(r Resource, scope constructs.Construct, id *string, props *ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.Resource",
		[]interface{}{scope, id, props},
		r,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Resource_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Resource",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Resource_IsResource(construct IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Resource",
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
func (r *jsiiProxy_Resource) ApplyRemovalPolicy(policy RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_Resource) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_Resource) GetResourceArnAttribute(arnAttr *string, arnComponents *ArnComponents) *string {
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
func (r *jsiiProxy_Resource) GetResourceNameAttribute(nameAttr *string) *string {
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
func (r *jsiiProxy_Resource) OnPrepare() {
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
func (r *jsiiProxy_Resource) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_Resource) OnValidate() *[]*string {
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
func (r *jsiiProxy_Resource) Prepare() {
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
func (r *jsiiProxy_Resource) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_Resource) ToString() *string {
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
func (r *jsiiProxy_Resource) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents the environment a given resource lives in.
//
// Used as the return value for the {@link IResource.env} property.
// Experimental.
type ResourceEnvironment struct {
	// The AWS account ID that this resource belongs to.
	//
	// Since this can be a Token
	// (for example, when the account is CloudFormation's AWS::AccountId intrinsic),
	// make sure to use Token.compareStrings()
	// instead of just comparing the values for equality.
	// Experimental.
	Account *string `json:"account"`
	// The AWS region that this resource belongs to.
	//
	// Since this can be a Token
	// (for example, when the region is CloudFormation's AWS::Region intrinsic),
	// make sure to use Token.compareStrings()
	// instead of just comparing the values for equality.
	// Experimental.
	Region *string `json:"region"`
}

// Construction properties for {@link Resource}.
// Experimental.
type ResourceProps struct {
	// The AWS account ID this resource belongs to.
	// Experimental.
	Account *string `json:"account"`
	// ARN to deduce region and account from.
	//
	// The ARN is parsed and the account and region are taken from the ARN.
	// This should be used for imported resources.
	//
	// Cannot be supplied together with either `account` or `region`.
	// Experimental.
	EnvironmentFromArn *string `json:"environmentFromArn"`
	// The value passed in by users to the physical name prop of the resource.
	//
	// - `undefined` implies that a physical name will be allocated by
	//    CloudFormation during deployment.
	// - a concrete value implies a specific physical name
	// - `PhysicalName.GENERATE_IF_NEEDED` is a marker that indicates that a physical will only be generated
	//    by the CDK if it is needed for cross-environment references. Otherwise, it will be allocated by CloudFormation.
	// Experimental.
	PhysicalName *string `json:"physicalName"`
	// The AWS region this resource belongs to.
	// Experimental.
	Region *string `json:"region"`
}

// Options for the 'reverse()' operation.
// Experimental.
type ReverseOptions struct {
	// Fail if the given string is a concatenation.
	//
	// If `false`, just return `undefined`.
	// Experimental.
	FailConcat *bool `json:"failConcat"`
}

// Accessor for scoped pseudo parameters.
//
// These pseudo parameters are anchored to a stack somewhere in the construct
// tree, and their values will be exported automatically.
// Experimental.
type ScopedAws interface {
	AccountId() *string
	NotificationArns() *[]*string
	Partition() *string
	Region() *string
	StackId() *string
	StackName() *string
	UrlSuffix() *string
}

// The jsii proxy struct for ScopedAws
type jsiiProxy_ScopedAws struct {
	_ byte // padding
}

func (j *jsiiProxy_ScopedAws) AccountId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"accountId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScopedAws) NotificationArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"notificationArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScopedAws) Partition() *string {
	var returns *string
	_jsii_.Get(
		j,
		"partition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScopedAws) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScopedAws) StackId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScopedAws) StackName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ScopedAws) UrlSuffix() *string {
	var returns *string
	_jsii_.Get(
		j,
		"urlSuffix",
		&returns,
	)
	return returns
}


// Experimental.
func NewScopedAws(scope Construct) ScopedAws {
	_init_.Initialize()

	j := jsiiProxy_ScopedAws{}

	_jsii_.Create(
		"monocdk.ScopedAws",
		[]interface{}{scope},
		&j,
	)

	return &j
}

// Experimental.
func NewScopedAws_Override(s ScopedAws, scope Construct) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.ScopedAws",
		[]interface{}{scope},
		s,
	)
}

// Work with secret values in the CDK.
//
// Secret values in the CDK (such as those retrieved from SecretsManager) are
// represented as regular strings, just like other values that are only
// available at deployment time.
//
// To help you avoid accidental mistakes which would lead to you putting your
// secret values directly into a CloudFormation template, constructs that take
// secret values will not allow you to pass in a literal secret value. They do
// so by calling `Secret.assertSafeSecret()`.
//
// You can escape the check by calling `Secret.plainText()`, but doing
// so is highly discouraged.
// Experimental.
type SecretValue interface {
	Intrinsic
	CreationStack() *[]*string
	NewError(message *string) interface{}
	Resolve(_context IResolveContext) interface{}
	ToJSON() interface{}
	ToString() *string
}

// The jsii proxy struct for SecretValue
type jsiiProxy_SecretValue struct {
	jsiiProxy_Intrinsic
}

func (j *jsiiProxy_SecretValue) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}


// Experimental.
func NewSecretValue(value interface{}, options *IntrinsicProps) SecretValue {
	_init_.Initialize()

	j := jsiiProxy_SecretValue{}

	_jsii_.Create(
		"monocdk.SecretValue",
		[]interface{}{value, options},
		&j,
	)

	return &j
}

// Experimental.
func NewSecretValue_Override(s SecretValue, value interface{}, options *IntrinsicProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.SecretValue",
		[]interface{}{value, options},
		s,
	)
}

// Obtain the secret value through a CloudFormation dynamic reference.
//
// If possible, use `SecretValue.ssmSecure` or `SecretValue.secretsManager` directly.
// Experimental.
func SecretValue_CfnDynamicReference(ref CfnDynamicReference) SecretValue {
	_init_.Initialize()

	var returns SecretValue

	_jsii_.StaticInvoke(
		"monocdk.SecretValue",
		"cfnDynamicReference",
		[]interface{}{ref},
		&returns,
	)

	return returns
}

// Obtain the secret value through a CloudFormation parameter.
//
// Generally, this is not a recommended approach. AWS Secrets Manager is the
// recommended way to reference secrets.
// Experimental.
func SecretValue_CfnParameter(param CfnParameter) SecretValue {
	_init_.Initialize()

	var returns SecretValue

	_jsii_.StaticInvoke(
		"monocdk.SecretValue",
		"cfnParameter",
		[]interface{}{param},
		&returns,
	)

	return returns
}

// Construct a literal secret value for use with secret-aware constructs.
//
// *Do not use this method for any secrets that you care about.*
//
// The only reasonable use case for using this method is when you are testing.
// Experimental.
func SecretValue_PlainText(secret *string) SecretValue {
	_init_.Initialize()

	var returns SecretValue

	_jsii_.StaticInvoke(
		"monocdk.SecretValue",
		"plainText",
		[]interface{}{secret},
		&returns,
	)

	return returns
}

// Creates a `SecretValue` with a value which is dynamically loaded from AWS Secrets Manager.
// Experimental.
func SecretValue_SecretsManager(secretId *string, options *SecretsManagerSecretOptions) SecretValue {
	_init_.Initialize()

	var returns SecretValue

	_jsii_.StaticInvoke(
		"monocdk.SecretValue",
		"secretsManager",
		[]interface{}{secretId, options},
		&returns,
	)

	return returns
}

// Use a secret value stored from a Systems Manager (SSM) parameter.
// Experimental.
func SecretValue_SsmSecure(parameterName *string, version *string) SecretValue {
	_init_.Initialize()

	var returns SecretValue

	_jsii_.StaticInvoke(
		"monocdk.SecretValue",
		"ssmSecure",
		[]interface{}{parameterName, version},
		&returns,
	)

	return returns
}

// Creates a throwable Error object that contains the token creation stack trace.
// Experimental.
func (s *jsiiProxy_SecretValue) NewError(message *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"newError",
		[]interface{}{message},
		&returns,
	)

	return returns
}

// Produce the Token's value at resolution time.
// Experimental.
func (s *jsiiProxy_SecretValue) Resolve(_context IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

// Turn this Token into JSON.
//
// Called automatically when JSON.stringify() is called on a Token.
// Experimental.
func (s *jsiiProxy_SecretValue) ToJSON() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Convert an instance of this Token to a string.
//
// This method will be called implicitly by language runtimes if the object
// is embedded into a string. We treat it the same as an explicit
// stringification.
// Experimental.
func (s *jsiiProxy_SecretValue) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options for referencing a secret value from Secrets Manager.
// Experimental.
type SecretsManagerSecretOptions struct {
	// The key of a JSON field to retrieve.
	//
	// This can only be used if the secret
	// stores a JSON object.
	// Experimental.
	JsonField *string `json:"jsonField"`
	// Specifies the unique identifier of the version of the secret you want to use.
	//
	// Can specify at most one of `versionId` and `versionStage`.
	// Experimental.
	VersionId *string `json:"versionId"`
	// Specifies the secret version that you want to retrieve by the staging label attached to the version.
	//
	// Can specify at most one of `versionId` and `versionStage`.
	// Experimental.
	VersionStage *string `json:"versionStage"`
}

// Represents the amount of digital storage.
//
// The amount can be specified either as a literal value (e.g: `10`) which
// cannot be negative, or as an unresolved number token.
//
// When the amount is passed as a token, unit conversion is not possible.
// Experimental.
type Size interface {
	ToGibibytes(opts *SizeConversionOptions) *float64
	ToKibibytes(opts *SizeConversionOptions) *float64
	ToMebibytes(opts *SizeConversionOptions) *float64
	ToPebibytes(opts *SizeConversionOptions) *float64
	ToTebibytes(opts *SizeConversionOptions) *float64
}

// The jsii proxy struct for Size
type jsiiProxy_Size struct {
	_ byte // padding
}

// Create a Storage representing an amount gibibytes.
//
// 1 GiB = 1024 MiB
//
// Returns: a new `Size` instance
// Experimental.
func Size_Gibibytes(amount *float64) Size {
	_init_.Initialize()

	var returns Size

	_jsii_.StaticInvoke(
		"monocdk.Size",
		"gibibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount kibibytes.
//
// 1 KiB = 1024 bytes
//
// Returns: a new `Size` instance
// Experimental.
func Size_Kibibytes(amount *float64) Size {
	_init_.Initialize()

	var returns Size

	_jsii_.StaticInvoke(
		"monocdk.Size",
		"kibibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount mebibytes.
//
// 1 MiB = 1024 KiB
//
// Returns: a new `Size` instance
// Experimental.
func Size_Mebibytes(amount *float64) Size {
	_init_.Initialize()

	var returns Size

	_jsii_.StaticInvoke(
		"monocdk.Size",
		"mebibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount pebibytes.
//
// 1 PiB = 1024 TiB
// Deprecated: use `pebibytes` instead
func Size_Pebibyte(amount *float64) Size {
	_init_.Initialize()

	var returns Size

	_jsii_.StaticInvoke(
		"monocdk.Size",
		"pebibyte",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount pebibytes.
//
// 1 PiB = 1024 TiB
//
// Returns: a new `Size` instance
// Experimental.
func Size_Pebibytes(amount *float64) Size {
	_init_.Initialize()

	var returns Size

	_jsii_.StaticInvoke(
		"monocdk.Size",
		"pebibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount tebibytes.
//
// 1 TiB = 1024 GiB
//
// Returns: a new `Size` instance
// Experimental.
func Size_Tebibytes(amount *float64) Size {
	_init_.Initialize()

	var returns Size

	_jsii_.StaticInvoke(
		"monocdk.Size",
		"tebibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Return this storage as a total number of gibibytes.
//
// Returns: the quantity of bytes expressed in gibibytes
// Experimental.
func (s *jsiiProxy_Size) ToGibibytes(opts *SizeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		s,
		"toGibibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Return this storage as a total number of kibibytes.
//
// Returns: the quantity of bytes expressed in kibibytes
// Experimental.
func (s *jsiiProxy_Size) ToKibibytes(opts *SizeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		s,
		"toKibibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Return this storage as a total number of mebibytes.
//
// Returns: the quantity of bytes expressed in mebibytes
// Experimental.
func (s *jsiiProxy_Size) ToMebibytes(opts *SizeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		s,
		"toMebibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Return this storage as a total number of pebibytes.
//
// Returns: the quantity of bytes expressed in pebibytes
// Experimental.
func (s *jsiiProxy_Size) ToPebibytes(opts *SizeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		s,
		"toPebibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Return this storage as a total number of tebibytes.
//
// Returns: the quantity of bytes expressed in tebibytes
// Experimental.
func (s *jsiiProxy_Size) ToTebibytes(opts *SizeConversionOptions) *float64 {
	var returns *float64

	_jsii_.Invoke(
		s,
		"toTebibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

// Options for how to convert time to a different unit.
// Experimental.
type SizeConversionOptions struct {
	// How conversions should behave when it encounters a non-integer result.
	// Experimental.
	Rounding SizeRoundingBehavior `json:"rounding"`
}

// Rounding behaviour when converting between units of `Size`.
// Experimental.
type SizeRoundingBehavior string

const (
	SizeRoundingBehavior_FAIL SizeRoundingBehavior = "FAIL"
	SizeRoundingBehavior_FLOOR SizeRoundingBehavior = "FLOOR"
	SizeRoundingBehavior_NONE SizeRoundingBehavior = "NONE"
)

// A root construct which represents a single CloudFormation stack.
// Experimental.
type Stack interface {
	Construct
	ITaggable
	Account() *string
	ArtifactId() *string
	AvailabilityZones() *[]*string
	Dependencies() *[]Stack
	Environment() *string
	Nested() *bool
	NestedStackParent() Stack
	NestedStackResource() CfnResource
	Node() ConstructNode
	NotificationArns() *[]*string
	ParentStack() Stack
	Partition() *string
	Region() *string
	StackId() *string
	StackName() *string
	Synthesizer() IStackSynthesizer
	Tags() TagManager
	TemplateFile() *string
	TemplateOptions() ITemplateOptions
	TerminationProtection() *bool
	UrlSuffix() *string
	AddDependency(target Stack, reason *string)
	AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation
	AddFileAsset(asset *FileAssetSource) *FileAssetLocation
	AddTransform(transform *string)
	AllocateLogicalId(cfnElement CfnElement) *string
	ExportValue(exportedValue interface{}, options *ExportValueOptions) *string
	FormatArn(components *ArnComponents) *string
	GetLogicalId(element CfnElement) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	ParseArn(arn *string, sepIfToken *string, hasName *bool) *ArnComponents
	Prepare()
	PrepareCrossReference(_sourceStack Stack, reference Reference) IResolvable
	RenameLogicalId(oldId *string, newId *string)
	ReportMissingContext(report *cxapi.MissingContext)
	ReportMissingContextKey(report *cloudassemblyschema.MissingContext)
	Resolve(obj interface{}) interface{}
	SplitArn(arn *string, arnFormat ArnFormat) *ArnComponents
	Synthesize(session ISynthesisSession)
	ToJsonString(obj interface{}, space *float64) *string
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Stack
type jsiiProxy_Stack struct {
	jsiiProxy_Construct
	jsiiProxy_ITaggable
}

func (j *jsiiProxy_Stack) Account() *string {
	var returns *string
	_jsii_.Get(
		j,
		"account",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ArtifactId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"artifactId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) AvailabilityZones() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"availabilityZones",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Dependencies() *[]Stack {
	var returns *[]Stack
	_jsii_.Get(
		j,
		"dependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Environment() *string {
	var returns *string
	_jsii_.Get(
		j,
		"environment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Nested() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"nested",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) NestedStackParent() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"nestedStackParent",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) NestedStackResource() CfnResource {
	var returns CfnResource
	_jsii_.Get(
		j,
		"nestedStackResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) NotificationArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"notificationArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ParentStack() Stack {
	var returns Stack
	_jsii_.Get(
		j,
		"parentStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Partition() *string {
	var returns *string
	_jsii_.Get(
		j,
		"partition",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StackId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StackName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Synthesizer() IStackSynthesizer {
	var returns IStackSynthesizer
	_jsii_.Get(
		j,
		"synthesizer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Tags() TagManager {
	var returns TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) TemplateFile() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateFile",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) TemplateOptions() ITemplateOptions {
	var returns ITemplateOptions
	_jsii_.Get(
		j,
		"templateOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) TerminationProtection() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"terminationProtection",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) UrlSuffix() *string {
	var returns *string
	_jsii_.Get(
		j,
		"urlSuffix",
		&returns,
	)
	return returns
}


// Creates a new stack.
// Experimental.
func NewStack(scope constructs.Construct, id *string, props *StackProps) Stack {
	_init_.Initialize()

	j := jsiiProxy_Stack{}

	_jsii_.Create(
		"monocdk.Stack",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Creates a new stack.
// Experimental.
func NewStack_Override(s Stack, scope constructs.Construct, id *string, props *StackProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.Stack",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Stack_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Stack",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Stack.
//
// We do attribute detection since we can't reliably use 'instanceof'.
// Experimental.
func Stack_IsStack(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Stack",
		"isStack",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Looks up the first stack scope in which `construct` is defined.
//
// Fails if there is no stack up the tree.
// Experimental.
func Stack_Of(construct constructs.IConstruct) Stack {
	_init_.Initialize()

	var returns Stack

	_jsii_.StaticInvoke(
		"monocdk.Stack",
		"of",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a dependency between this stack and another stack.
//
// This can be used to define dependencies between any two stacks within an
// app, and also supports nested stacks.
// Experimental.
func (s *jsiiProxy_Stack) AddDependency(target Stack, reason *string) {
	_jsii_.InvokeVoid(
		s,
		"addDependency",
		[]interface{}{target, reason},
	)
}

// Register a docker image asset on this Stack.
// Deprecated: Use `stack.synthesizer.addDockerImageAsset()` if you are calling,
// and a different `IStackSynthesizer` class if you are implementing.
func (s *jsiiProxy_Stack) AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		s,
		"addDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Register a file asset on this Stack.
// Deprecated: Use `stack.synthesizer.addFileAsset()` if you are calling,
// and a different IStackSynthesizer class if you are implementing.
func (s *jsiiProxy_Stack) AddFileAsset(asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		s,
		"addFileAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Add a Transform to this stack. A Transform is a macro that AWS CloudFormation uses to process your template.
//
// Duplicate values are removed when stack is synthesized.
//
// TODO: EXAMPLE
//
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/transform-section-structure.html
//
// Experimental.
func (s *jsiiProxy_Stack) AddTransform(transform *string) {
	_jsii_.InvokeVoid(
		s,
		"addTransform",
		[]interface{}{transform},
	)
}

// Returns the naming scheme used to allocate logical IDs.
//
// By default, uses
// the `HashedAddressingScheme` but this method can be overridden to customize
// this behavior.
//
// In order to make sure logical IDs are unique and stable, we hash the resource
// construct tree path (i.e. toplevel/secondlevel/.../myresource) and add it as
// a suffix to the path components joined without a separator (CloudFormation
// IDs only allow alphanumeric characters).
//
// The result will be:
//
//    <path.join('')><md5(path.join('/')>
//      "human"      "hash"
//
// If the "human" part of the ID exceeds 240 characters, we simply trim it so
// the total ID doesn't exceed CloudFormation's 255 character limit.
//
// We only take 8 characters from the md5 hash (0.000005 chance of collision).
//
// Special cases:
//
// - If the path only contains a single component (i.e. it's a top-level
//    resource), we won't add the hash to it. The hash is not needed for
//    disamiguation and also, it allows for a more straightforward migration an
//    existing CloudFormation template to a CDK stack without logical ID changes
//    (or renames).
// - For aesthetic reasons, if the last components of the path are the same
//    (i.e. `L1/L2/Pipeline/Pipeline`), they will be de-duplicated to make the
//    resulting human portion of the ID more pleasing: `L1L2Pipeline<HASH>`
//    instead of `L1L2PipelinePipeline<HASH>`
// - If a component is named "Default" it will be omitted from the path. This
//    allows refactoring higher level abstractions around constructs without affecting
//    the IDs of already deployed resources.
// - If a component is named "Resource" it will be omitted from the user-visible
//    path, but included in the hash. This reduces visual noise in the human readable
//    part of the identifier.
// Experimental.
func (s *jsiiProxy_Stack) AllocateLogicalId(cfnElement CfnElement) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"allocateLogicalId",
		[]interface{}{cfnElement},
		&returns,
	)

	return returns
}

// Create a CloudFormation Export for a value.
//
// Returns a string representing the corresponding `Fn.importValue()`
// expression for this Export. You can control the name for the export by
// passing the `name` option.
//
// If you don't supply a value for `name`, the value you're exporting must be
// a Resource attribute (for example: `bucket.bucketName`) and it will be
// given the same name as the automatic cross-stack reference that would be created
// if you used the attribute in another Stack.
//
// One of the uses for this method is to *remove* the relationship between
// two Stacks established by automatic cross-stack references. It will
// temporarily ensure that the CloudFormation Export still exists while you
// remove the reference from the consuming stack. After that, you can remove
// the resource and the manual export.
//
// ## Example
//
// Here is how the process works. Let's say there are two stacks,
// `producerStack` and `consumerStack`, and `producerStack` has a bucket
// called `bucket`, which is referenced by `consumerStack` (perhaps because
// an AWS Lambda Function writes into it, or something like that).
//
// It is not safe to remove `producerStack.bucket` because as the bucket is being
// deleted, `consumerStack` might still be using it.
//
// Instead, the process takes two deployments:
//
// ### Deployment 1: break the relationship
//
// - Make sure `consumerStack` no longer references `bucket.bucketName` (maybe the consumer
//    stack now uses its own bucket, or it writes to an AWS DynamoDB table, or maybe you just
//    remove the Lambda Function altogether).
// - In the `ProducerStack` class, call `this.exportValue(this.bucket.bucketName)`. This
//    will make sure the CloudFormation Export continues to exist while the relationship
//    between the two stacks is being broken.
// - Deploy (this will effectively only change the `consumerStack`, but it's safe to deploy both).
//
// ### Deployment 2: remove the bucket resource
//
// - You are now free to remove the `bucket` resource from `producerStack`.
// - Don't forget to remove the `exportValue()` call as well.
// - Deploy again (this time only the `producerStack` will be changed -- the bucket will be deleted).
// Experimental.
func (s *jsiiProxy_Stack) ExportValue(exportedValue interface{}, options *ExportValueOptions) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"exportValue",
		[]interface{}{exportedValue, options},
		&returns,
	)

	return returns
}

// Creates an ARN from components.
//
// If `partition`, `region` or `account` are not specified, the stack's
// partition, region and account will be used.
//
// If any component is the empty string, an empty string will be inserted
// into the generated ARN at the location that component corresponds to.
//
// The ARN will be formatted as follows:
//
//    arn:{partition}:{service}:{region}:{account}:{resource}{sep}}{resource-name}
//
// The required ARN pieces that are omitted will be taken from the stack that
// the 'scope' is attached to. If all ARN pieces are supplied, the supplied scope
// can be 'undefined'.
// Experimental.
func (s *jsiiProxy_Stack) FormatArn(components *ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"formatArn",
		[]interface{}{components},
		&returns,
	)

	return returns
}

// Allocates a stack-unique CloudFormation-compatible logical identity for a specific resource.
//
// This method is called when a `CfnElement` is created and used to render the
// initial logical identity of resources. Logical ID renames are applied at
// this stage.
//
// This method uses the protected method `allocateLogicalId` to render the
// logical ID for an element. To modify the naming scheme, extend the `Stack`
// class and override this method.
// Experimental.
func (s *jsiiProxy_Stack) GetLogicalId(element CfnElement) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getLogicalId",
		[]interface{}{element},
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
func (s *jsiiProxy_Stack) OnPrepare() {
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
func (s *jsiiProxy_Stack) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_Stack) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"onValidate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Given an ARN, parses it and returns components.
//
// IF THE ARN IS A CONCRETE STRING...
//
// ...it will be parsed and validated. The separator (`sep`) will be set to '/'
// if the 6th component includes a '/', in which case, `resource` will be set
// to the value before the '/' and `resourceName` will be the rest. In case
// there is no '/', `resource` will be set to the 6th components and
// `resourceName` will be set to the rest of the string.
//
// IF THE ARN IS A TOKEN...
//
// ...it cannot be validated, since we don't have the actual value yet at the
// time of this function call. You will have to supply `sepIfToken` and
// whether or not ARNs of the expected format usually have resource names
// in order to parse it properly. The resulting `ArnComponents` object will
// contain tokens for the subexpressions of the ARN, not string literals.
//
// If the resource name could possibly contain the separator char, the actual
// resource name cannot be properly parsed. This only occurs if the separator
// char is '/', and happens for example for S3 object ARNs, IAM Role ARNs,
// IAM OIDC Provider ARNs, etc. To properly extract the resource name from a
// Tokenized ARN, you must know the resource type and call
// `Arn.extractResourceName`.
//
// Returns: an ArnComponents object which allows access to the various
// components of the ARN.
// Deprecated: use splitArn instead
func (s *jsiiProxy_Stack) ParseArn(arn *string, sepIfToken *string, hasName *bool) *ArnComponents {
	var returns *ArnComponents

	_jsii_.Invoke(
		s,
		"parseArn",
		[]interface{}{arn, sepIfToken, hasName},
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
func (s *jsiiProxy_Stack) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated.
//
// Returns: reference itself without any change
// See: https://github.com/aws/aws-cdk/pull/7187
//
// Deprecated: cross reference handling has been moved to `App.prepare()`.
func (s *jsiiProxy_Stack) PrepareCrossReference(_sourceStack Stack, reference Reference) IResolvable {
	var returns IResolvable

	_jsii_.Invoke(
		s,
		"prepareCrossReference",
		[]interface{}{_sourceStack, reference},
		&returns,
	)

	return returns
}

// Rename a generated logical identities.
//
// To modify the naming scheme strategy, extend the `Stack` class and
// override the `allocateLogicalId` method.
// Experimental.
func (s *jsiiProxy_Stack) RenameLogicalId(oldId *string, newId *string) {
	_jsii_.InvokeVoid(
		s,
		"renameLogicalId",
		[]interface{}{oldId, newId},
	)
}

// DEPRECATED.
// Deprecated: use `reportMissingContextKey()`
func (s *jsiiProxy_Stack) ReportMissingContext(report *cxapi.MissingContext) {
	_jsii_.InvokeVoid(
		s,
		"reportMissingContext",
		[]interface{}{report},
	)
}

// Indicate that a context key was expected.
//
// Contains instructions which will be emitted into the cloud assembly on how
// the key should be supplied.
// Experimental.
func (s *jsiiProxy_Stack) ReportMissingContextKey(report *cloudassemblyschema.MissingContext) {
	_jsii_.InvokeVoid(
		s,
		"reportMissingContextKey",
		[]interface{}{report},
	)
}

// Resolve a tokenized value in the context of the current stack.
// Experimental.
func (s *jsiiProxy_Stack) Resolve(obj interface{}) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"resolve",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Splits the provided ARN into its components.
//
// Works both if 'arn' is a string like 'arn:aws:s3:::bucket',
// and a Token representing a dynamic CloudFormation expression
// (in which case the returned components will also be dynamic CloudFormation expressions,
// encoded as Tokens).
// Experimental.
func (s *jsiiProxy_Stack) SplitArn(arn *string, arnFormat ArnFormat) *ArnComponents {
	var returns *ArnComponents

	_jsii_.Invoke(
		s,
		"splitArn",
		[]interface{}{arn, arnFormat},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_Stack) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Convert an object, potentially containing tokens, to a JSON string.
// Experimental.
func (s *jsiiProxy_Stack) ToJsonString(obj interface{}, space *float64) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toJsonString",
		[]interface{}{obj, space},
		&returns,
	)

	return returns
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_Stack) ToString() *string {
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
func (s *jsiiProxy_Stack) Validate() *[]*string {
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
type StackProps struct {
	// Include runtime versioning information in this Stack.
	// Experimental.
	AnalyticsReporting *bool `json:"analyticsReporting"`
	// A description of the stack.
	// Experimental.
	Description *string `json:"description"`
	// The AWS environment (account/region) where this stack will be deployed.
	//
	// Set the `region`/`account` fields of `env` to either a concrete value to
	// select the indicated environment (recommended for production stacks), or to
	// the values of environment variables
	// `CDK_DEFAULT_REGION`/`CDK_DEFAULT_ACCOUNT` to let the target environment
	// depend on the AWS credentials/configuration that the CDK CLI is executed
	// under (recommended for development stacks).
	//
	// If the `Stack` is instantiated inside a `Stage`, any undefined
	// `region`/`account` fields from `env` will default to the same field on the
	// encompassing `Stage`, if configured there.
	//
	// If either `region` or `account` are not set nor inherited from `Stage`, the
	// Stack will be considered "*environment-agnostic*"". Environment-agnostic
	// stacks can be deployed to any environment but may not be able to take
	// advantage of all features of the CDK. For example, they will not be able to
	// use environmental context lookups such as `ec2.Vpc.fromLookup` and will not
	// automatically translate Service Principals to the right format based on the
	// environment's AWS partition, and other such enhancements.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Env *Environment `json:"env"`
	// Name to deploy the stack with.
	// Experimental.
	StackName *string `json:"stackName"`
	// Synthesis method to use while deploying this stack.
	// Experimental.
	Synthesizer IStackSynthesizer `json:"synthesizer"`
	// Stack tags that will be applied to all the taggable resources and the stack itself.
	// Experimental.
	Tags *map[string]*string `json:"tags"`
	// Whether to enable termination protection for this stack.
	// Experimental.
	TerminationProtection *bool `json:"terminationProtection"`
}

// Base class for implementing an IStackSynthesizer.
//
// This class needs to exist to provide public surface area for external
// implementations of stack synthesizers. The protected methods give
// access to functions that are otherwise @_internal to the framework
// and could not be accessed by external implementors.
// Experimental.
type StackSynthesizer interface {
	IStackSynthesizer
	AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation
	AddFileAsset(asset *FileAssetSource) *FileAssetLocation
	Bind(stack Stack)
	EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions)
	Synthesize(session ISynthesisSession)
	SynthesizeStackTemplate(stack Stack, session ISynthesisSession)
}

// The jsii proxy struct for StackSynthesizer
type jsiiProxy_StackSynthesizer struct {
	jsiiProxy_IStackSynthesizer
}

// Experimental.
func NewStackSynthesizer_Override(s StackSynthesizer) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.StackSynthesizer",
		nil, // no parameters
		s,
	)
}

// Register a Docker Image Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (s *jsiiProxy_StackSynthesizer) AddDockerImageAsset(asset *DockerImageAssetSource) *DockerImageAssetLocation {
	var returns *DockerImageAssetLocation

	_jsii_.Invoke(
		s,
		"addDockerImageAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Register a File Asset.
//
// Returns the parameters that can be used to refer to the asset inside the template.
// Experimental.
func (s *jsiiProxy_StackSynthesizer) AddFileAsset(asset *FileAssetSource) *FileAssetLocation {
	var returns *FileAssetLocation

	_jsii_.Invoke(
		s,
		"addFileAsset",
		[]interface{}{asset},
		&returns,
	)

	return returns
}

// Bind to the stack this environment is going to be used on.
//
// Must be called before any of the other methods are called.
// Experimental.
func (s *jsiiProxy_StackSynthesizer) Bind(stack Stack) {
	_jsii_.InvokeVoid(
		s,
		"bind",
		[]interface{}{stack},
	)
}

// Write the stack artifact to the session.
//
// Use default settings to add a CloudFormationStackArtifact artifact to
// the given synthesis session.
// Experimental.
func (s *jsiiProxy_StackSynthesizer) EmitStackArtifact(stack Stack, session ISynthesisSession, options *SynthesizeStackArtifactOptions) {
	_jsii_.InvokeVoid(
		s,
		"emitStackArtifact",
		[]interface{}{stack, session, options},
	)
}

// Synthesize the associated stack to the session.
// Experimental.
func (s *jsiiProxy_StackSynthesizer) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Have the stack write out its template.
// Experimental.
func (s *jsiiProxy_StackSynthesizer) SynthesizeStackTemplate(stack Stack, session ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesizeStackTemplate",
		[]interface{}{stack, session},
	)
}

// An abstract application modeling unit consisting of Stacks that should be deployed together.
//
// Derive a subclass of `Stage` and use it to model a single instance of your
// application.
//
// You can then instantiate your subclass multiple times to model multiple
// copies of your application which should be be deployed to different
// environments.
// Experimental.
type Stage interface {
	Construct
	Account() *string
	ArtifactId() *string
	AssetOutdir() *string
	Node() ConstructNode
	Outdir() *string
	ParentStage() Stage
	Region() *string
	StageName() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synth(options *StageSynthesisOptions) cxapi.CloudAssembly
	Synthesize(session ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Stage
type jsiiProxy_Stage struct {
	jsiiProxy_Construct
}

func (j *jsiiProxy_Stage) Account() *string {
	var returns *string
	_jsii_.Get(
		j,
		"account",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) ArtifactId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"artifactId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) AssetOutdir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assetOutdir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) Node() ConstructNode {
	var returns ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) Outdir() *string {
	var returns *string
	_jsii_.Get(
		j,
		"outdir",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) ParentStage() Stage {
	var returns Stage
	_jsii_.Get(
		j,
		"parentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) Region() *string {
	var returns *string
	_jsii_.Get(
		j,
		"region",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}


// Experimental.
func NewStage(scope constructs.Construct, id *string, props *StageProps) Stage {
	_init_.Initialize()

	j := jsiiProxy_Stage{}

	_jsii_.Create(
		"monocdk.Stage",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStage_Override(s Stage, scope constructs.Construct, id *string, props *StageProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.Stage",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Stage_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Stage",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Test whether the given construct is a stage.
// Experimental.
func Stage_IsStage(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Stage",
		"isStage",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return the stage this construct is contained with, if available.
//
// If called
// on a nested stage, returns its parent.
// Experimental.
func Stage_Of(construct constructs.IConstruct) Stage {
	_init_.Initialize()

	var returns Stage

	_jsii_.StaticInvoke(
		"monocdk.Stage",
		"of",
		[]interface{}{construct},
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
func (s *jsiiProxy_Stage) OnPrepare() {
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
func (s *jsiiProxy_Stage) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_Stage) OnValidate() *[]*string {
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
func (s *jsiiProxy_Stage) Prepare() {
	_jsii_.InvokeVoid(
		s,
		"prepare",
		nil, // no parameters
	)
}

// Synthesize this stage into a cloud assembly.
//
// Once an assembly has been synthesized, it cannot be modified. Subsequent
// calls will return the same assembly.
// Experimental.
func (s *jsiiProxy_Stage) Synth(options *StageSynthesisOptions) cxapi.CloudAssembly {
	var returns cxapi.CloudAssembly

	_jsii_.Invoke(
		s,
		"synth",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (s *jsiiProxy_Stage) Synthesize(session ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_Stage) ToString() *string {
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
func (s *jsiiProxy_Stage) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Initialization props for a stage.
// Experimental.
type StageProps struct {
	// Default AWS environment (account/region) for `Stack`s in this `Stage`.
	//
	// Stacks defined inside this `Stage` with either `region` or `account` missing
	// from its env will use the corresponding field given here.
	//
	// If either `region` or `account`is is not configured for `Stack` (either on
	// the `Stack` itself or on the containing `Stage`), the Stack will be
	// *environment-agnostic*.
	//
	// Environment-agnostic stacks can be deployed to any environment, may not be
	// able to take advantage of all features of the CDK. For example, they will
	// not be able to use environmental context lookups, will not automatically
	// translate Service Principals to the right format based on the environment's
	// AWS partition, and other such enhancements.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	Env *Environment `json:"env"`
	// The output directory into which to emit synthesized artifacts.
	//
	// Can only be specified if this stage is the root stage (the app). If this is
	// specified and this stage is nested within another stage, an error will be
	// thrown.
	// Experimental.
	Outdir *string `json:"outdir"`
}

// Options for assembly synthesis.
// Experimental.
type StageSynthesisOptions struct {
	// Force a re-synth, even if the stage has already been synthesized.
	//
	// This is used by tests to allow for incremental verification of the output.
	// Do not use in production.
	// Experimental.
	Force *bool `json:"force"`
	// Should we skip construct validation.
	// Experimental.
	SkipValidation *bool `json:"skipValidation"`
	// Whether the stack should be validated after synthesis to check for error metadata.
	// Experimental.
	ValidateOnSynthesis *bool `json:"validateOnSynthesis"`
}

// Converts all fragments to strings and concats those.
//
// Drops 'undefined's.
// Experimental.
type StringConcat interface {
	IFragmentConcatenator
	Join(left interface{}, right interface{}) interface{}
}

// The jsii proxy struct for StringConcat
type jsiiProxy_StringConcat struct {
	jsiiProxy_IFragmentConcatenator
}

// Experimental.
func NewStringConcat() StringConcat {
	_init_.Initialize()

	j := jsiiProxy_StringConcat{}

	_jsii_.Create(
		"monocdk.StringConcat",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewStringConcat_Override(s StringConcat) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.StringConcat",
		nil, // no parameters
		s,
	)
}

// Join the fragment on the left and on the right.
// Experimental.
func (s *jsiiProxy_StringConcat) Join(left interface{}, right interface{}) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"join",
		[]interface{}{left, right},
		&returns,
	)

	return returns
}

// Determines how symlinks are followed.
// Experimental.
type SymlinkFollowMode string

const (
	SymlinkFollowMode_NEVER SymlinkFollowMode = "NEVER"
	SymlinkFollowMode_ALWAYS SymlinkFollowMode = "ALWAYS"
	SymlinkFollowMode_EXTERNAL SymlinkFollowMode = "EXTERNAL"
	SymlinkFollowMode_BLOCK_EXTERNAL SymlinkFollowMode = "BLOCK_EXTERNAL"
)

// Options for synthesis.
// Deprecated: use `app.synth()` or `stage.synth()` instead
type SynthesisOptions struct {
	// Include the specified runtime information (module versions) in manifest.
	// Deprecated: All template modifications that should result from this should
	// have already been inserted into the template.
	RuntimeInfo *cxapi.RuntimeInfo `json:"runtimeInfo"`
	// The output directory into which to synthesize the cloud assembly.
	// Deprecated: use `app.synth()` or `stage.synth()` instead
	Outdir *string `json:"outdir"`
	// Whether synthesis should skip the validation phase.
	// Deprecated: use `app.synth()` or `stage.synth()` instead
	SkipValidation *bool `json:"skipValidation"`
	// Whether the stack should be validated after synthesis to check for error metadata.
	// Deprecated: use `app.synth()` or `stage.synth()` instead
	ValidateOnSynthesis *bool `json:"validateOnSynthesis"`
}

// Stack artifact options.
//
// A subset of `cxschema.AwsCloudFormationStackProperties` of optional settings that need to be
// configurable by synthesizers, plus `additionalDependencies`.
// Experimental.
type SynthesizeStackArtifactOptions struct {
	// Identifiers of additional dependencies.
	// Experimental.
	AdditionalDependencies *[]*string `json:"additionalDependencies"`
	// The role that needs to be assumed to deploy the stack.
	// Experimental.
	AssumeRoleArn *string `json:"assumeRoleArn"`
	// SSM parameter where the bootstrap stack version number can be found.
	//
	// Only used if `requiresBootstrapStackVersion` is set.
	//
	// - If this value is not set, the bootstrap stack name must be known at
	//    deployment time so the stack version can be looked up from the stack
	//    outputs.
	// - If this value is set, the bootstrap stack can have any name because
	//    we won't need to look it up.
	// Experimental.
	BootstrapStackVersionSsmParameter *string `json:"bootstrapStackVersionSsmParameter"`
	// The role that is passed to CloudFormation to execute the change set.
	// Experimental.
	CloudFormationExecutionRoleArn *string `json:"cloudFormationExecutionRoleArn"`
	// Values for CloudFormation stack parameters that should be passed when the stack is deployed.
	// Experimental.
	Parameters *map[string]*string `json:"parameters"`
	// Version of bootstrap stack required to deploy this stack.
	// Experimental.
	RequiresBootstrapStackVersion *float64 `json:"requiresBootstrapStackVersion"`
	// If the stack template has already been included in the asset manifest, its asset URL.
	// Experimental.
	StackTemplateAssetObjectUrl *string `json:"stackTemplateAssetObjectUrl"`
}

// The Tag Aspect will handle adding a tag to this node and cascading tags to children.
// Experimental.
type Tag interface {
	IAspect
	Key() *string
	Props() *TagProps
	Value() *string
	ApplyTag(resource ITaggable)
	Visit(construct IConstruct)
}

// The jsii proxy struct for Tag
type jsiiProxy_Tag struct {
	jsiiProxy_IAspect
}

func (j *jsiiProxy_Tag) Key() *string {
	var returns *string
	_jsii_.Get(
		j,
		"key",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Tag) Props() *TagProps {
	var returns *TagProps
	_jsii_.Get(
		j,
		"props",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Tag) Value() *string {
	var returns *string
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


// Experimental.
func NewTag(key *string, value *string, props *TagProps) Tag {
	_init_.Initialize()

	j := jsiiProxy_Tag{}

	_jsii_.Create(
		"monocdk.Tag",
		[]interface{}{key, value, props},
		&j,
	)

	return &j
}

// Experimental.
func NewTag_Override(t Tag, key *string, value *string, props *TagProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.Tag",
		[]interface{}{key, value, props},
		t,
	)
}

// DEPRECATED: add tags to the node of a construct and all its the taggable children.
// Deprecated: use `Tags.of(scope).add()`
func Tag_Add(scope Construct, key *string, value *string, props *TagProps) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.Tag",
		"add",
		[]interface{}{scope, key, value, props},
	)
}

// DEPRECATED: remove tags to the node of a construct and all its the taggable children.
// Deprecated: use `Tags.of(scope).remove()`
func Tag_Remove(scope Construct, key *string, props *TagProps) {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.Tag",
		"remove",
		[]interface{}{scope, key, props},
	)
}

// Experimental.
func (t *jsiiProxy_Tag) ApplyTag(resource ITaggable) {
	_jsii_.InvokeVoid(
		t,
		"applyTag",
		[]interface{}{resource},
	)
}

// All aspects can visit an IConstruct.
// Experimental.
func (t *jsiiProxy_Tag) Visit(construct IConstruct) {
	_jsii_.InvokeVoid(
		t,
		"visit",
		[]interface{}{construct},
	)
}

// TagManager facilitates a common implementation of tagging for Constructs.
// Experimental.
type TagManager interface {
	TagPropertyName() *string
	ApplyTagAspectHere(include *[]*string, exclude *[]*string) *bool
	HasTags() *bool
	RemoveTag(key *string, priority *float64)
	RenderTags() interface{}
	SetTag(key *string, value *string, priority *float64, applyToLaunchedInstances *bool)
	TagValues() *map[string]*string
}

// The jsii proxy struct for TagManager
type jsiiProxy_TagManager struct {
	_ byte // padding
}

func (j *jsiiProxy_TagManager) TagPropertyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tagPropertyName",
		&returns,
	)
	return returns
}


// Experimental.
func NewTagManager(tagType TagType, resourceTypeName *string, tagStructure interface{}, options *TagManagerOptions) TagManager {
	_init_.Initialize()

	j := jsiiProxy_TagManager{}

	_jsii_.Create(
		"monocdk.TagManager",
		[]interface{}{tagType, resourceTypeName, tagStructure, options},
		&j,
	)

	return &j
}

// Experimental.
func NewTagManager_Override(t TagManager, tagType TagType, resourceTypeName *string, tagStructure interface{}, options *TagManagerOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.TagManager",
		[]interface{}{tagType, resourceTypeName, tagStructure, options},
		t,
	)
}

// Check whether the given construct is Taggable.
// Experimental.
func TagManager_IsTaggable(construct interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.TagManager",
		"isTaggable",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Determine if the aspect applies here.
//
// Looks at the include and exclude resourceTypeName arrays to determine if
// the aspect applies here
// Experimental.
func (t *jsiiProxy_TagManager) ApplyTagAspectHere(include *[]*string, exclude *[]*string) *bool {
	var returns *bool

	_jsii_.Invoke(
		t,
		"applyTagAspectHere",
		[]interface{}{include, exclude},
		&returns,
	)

	return returns
}

// Returns true if there are any tags defined.
// Experimental.
func (t *jsiiProxy_TagManager) HasTags() *bool {
	var returns *bool

	_jsii_.Invoke(
		t,
		"hasTags",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Removes the specified tag from the array if it exists.
// Experimental.
func (t *jsiiProxy_TagManager) RemoveTag(key *string, priority *float64) {
	_jsii_.InvokeVoid(
		t,
		"removeTag",
		[]interface{}{key, priority},
	)
}

// Renders tags into the proper format based on TagType.
// Experimental.
func (t *jsiiProxy_TagManager) RenderTags() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"renderTags",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Adds the specified tag to the array of tags.
// Experimental.
func (t *jsiiProxy_TagManager) SetTag(key *string, value *string, priority *float64, applyToLaunchedInstances *bool) {
	_jsii_.InvokeVoid(
		t,
		"setTag",
		[]interface{}{key, value, priority, applyToLaunchedInstances},
	)
}

// Render the tags in a readable format.
// Experimental.
func (t *jsiiProxy_TagManager) TagValues() *map[string]*string {
	var returns *map[string]*string

	_jsii_.Invoke(
		t,
		"tagValues",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options to configure TagManager behavior.
// Experimental.
type TagManagerOptions struct {
	// The name of the property in CloudFormation for these tags.
	//
	// Normally this is `tags`, but Cognito UserPool uses UserPoolTags
	// Experimental.
	TagPropertyName *string `json:"tagPropertyName"`
}

// Properties for a tag.
// Experimental.
type TagProps struct {
	// Whether the tag should be applied to instances in an AutoScalingGroup.
	// Experimental.
	ApplyToLaunchedInstances *bool `json:"applyToLaunchedInstances"`
	// An array of Resource Types that will not receive this tag.
	//
	// An empty array will allow this tag to be applied to all resources. A
	// non-empty array will apply this tag only if the Resource type is not in
	// this array.
	// Experimental.
	ExcludeResourceTypes *[]*string `json:"excludeResourceTypes"`
	// An array of Resource Types that will receive this tag.
	//
	// An empty array will match any Resource. A non-empty array will apply this
	// tag only to Resource types that are included in this array.
	// Experimental.
	IncludeResourceTypes *[]*string `json:"includeResourceTypes"`
	// Priority of the tag operation.
	//
	// Higher or equal priority tags will take precedence.
	//
	// Setting priority will enable the user to control tags when they need to not
	// follow the default precedence pattern of last applied and closest to the
	// construct in the tree.
	// Experimental.
	Priority *float64 `json:"priority"`
}

// Experimental.
type TagType string

const (
	TagType_STANDARD TagType = "STANDARD"
	TagType_AUTOSCALING_GROUP TagType = "AUTOSCALING_GROUP"
	TagType_MAP TagType = "MAP"
	TagType_KEY_VALUE TagType = "KEY_VALUE"
	TagType_NOT_TAGGABLE TagType = "NOT_TAGGABLE"
)

// Manages AWS tags for all resources within a construct scope.
// Experimental.
type Tags interface {
	Add(key *string, value *string, props *TagProps)
	Remove(key *string, props *TagProps)
}

// The jsii proxy struct for Tags
type jsiiProxy_Tags struct {
	_ byte // padding
}

// Returns the tags API for this scope.
// Experimental.
func Tags_Of(scope IConstruct) Tags {
	_init_.Initialize()

	var returns Tags

	_jsii_.StaticInvoke(
		"monocdk.Tags",
		"of",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// add tags to the node of a construct and all its the taggable children.
// Experimental.
func (t *jsiiProxy_Tags) Add(key *string, value *string, props *TagProps) {
	_jsii_.InvokeVoid(
		t,
		"add",
		[]interface{}{key, value, props},
	)
}

// remove tags to the node of a construct and all its the taggable children.
// Experimental.
func (t *jsiiProxy_Tags) Remove(key *string, props *TagProps) {
	_jsii_.InvokeVoid(
		t,
		"remove",
		[]interface{}{key, props},
	)
}

// Options for how to convert time to a different unit.
// Experimental.
type TimeConversionOptions struct {
	// If `true`, conversions into a larger time unit (e.g. `Seconds` to `Minutes`) will fail if the result is not an integer.
	// Experimental.
	Integral *bool `json:"integral"`
}

// Represents a special or lazily-evaluated value.
//
// Can be used to delay evaluation of a certain value in case, for example,
// that it requires some context or late-bound data. Can also be used to
// mark values that need special processing at document rendering time.
//
// Tokens can be embedded into strings while retaining their original
// semantics.
// Experimental.
type Token interface {
}

// The jsii proxy struct for Token
type jsiiProxy_Token struct {
	_ byte // padding
}

// Return a resolvable representation of the given value.
// Experimental.
func Token_AsAny(value interface{}) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Token",
		"asAny",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// Return a reversible list representation of this token.
// Experimental.
func Token_AsList(value interface{}, options *EncodingOptions) *[]*string {
	_init_.Initialize()

	var returns *[]*string

	_jsii_.StaticInvoke(
		"monocdk.Token",
		"asList",
		[]interface{}{value, options},
		&returns,
	)

	return returns
}

// Return a reversible number representation of this token.
// Experimental.
func Token_AsNumber(value interface{}) *float64 {
	_init_.Initialize()

	var returns *float64

	_jsii_.StaticInvoke(
		"monocdk.Token",
		"asNumber",
		[]interface{}{value},
		&returns,
	)

	return returns
}

// Return a reversible string representation of this token.
//
// If the Token is initialized with a literal, the stringified value of the
// literal is returned. Otherwise, a special quoted string representation
// of the Token is returned that can be embedded into other strings.
//
// Strings with quoted Tokens in them can be restored back into
// complex values with the Tokens restored by calling `resolve()`
// on the string.
// Experimental.
func Token_AsString(value interface{}, options *EncodingOptions) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Token",
		"asString",
		[]interface{}{value, options},
		&returns,
	)

	return returns
}

// Compare two strings that might contain Tokens with each other.
// Experimental.
func Token_CompareStrings(possibleToken1 *string, possibleToken2 *string) TokenComparison {
	_init_.Initialize()

	var returns TokenComparison

	_jsii_.StaticInvoke(
		"monocdk.Token",
		"compareStrings",
		[]interface{}{possibleToken1, possibleToken2},
		&returns,
	)

	return returns
}

// Returns true if obj represents an unresolved value.
//
// One of these must be true:
//
// - `obj` is an IResolvable
// - `obj` is a string containing at least one encoded `IResolvable`
// - `obj` is either an encoded number or list
//
// This does NOT recurse into lists or objects to see if they
// containing resolvables.
// Experimental.
func Token_IsUnresolved(obj interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Token",
		"isUnresolved",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// An enum-like class that represents the result of comparing two Tokens.
//
// The return type of {@link Token.compareStrings}.
// Experimental.
type TokenComparison interface {
}

// The jsii proxy struct for TokenComparison
type jsiiProxy_TokenComparison struct {
	_ byte // padding
}

func TokenComparison_BOTH_UNRESOLVED() TokenComparison {
	_init_.Initialize()
	var returns TokenComparison
	_jsii_.StaticGet(
		"monocdk.TokenComparison",
		"BOTH_UNRESOLVED",
		&returns,
	)
	return returns
}

func TokenComparison_DIFFERENT() TokenComparison {
	_init_.Initialize()
	var returns TokenComparison
	_jsii_.StaticGet(
		"monocdk.TokenComparison",
		"DIFFERENT",
		&returns,
	)
	return returns
}

func TokenComparison_ONE_UNRESOLVED() TokenComparison {
	_init_.Initialize()
	var returns TokenComparison
	_jsii_.StaticGet(
		"monocdk.TokenComparison",
		"ONE_UNRESOLVED",
		&returns,
	)
	return returns
}

func TokenComparison_SAME() TokenComparison {
	_init_.Initialize()
	var returns TokenComparison
	_jsii_.StaticGet(
		"monocdk.TokenComparison",
		"SAME",
		&returns,
	)
	return returns
}

// Less oft-needed functions to manipulate Tokens.
// Experimental.
type Tokenization interface {
}

// The jsii proxy struct for Tokenization
type jsiiProxy_Tokenization struct {
	_ byte // padding
}

// Return whether the given object is an IResolvable object.
//
// This is different from Token.isUnresolved() which will also check for
// encoded Tokens, whereas this method will only do a type check on the given
// object.
// Experimental.
func Tokenization_IsResolvable(obj interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"isResolvable",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Resolves an object by evaluating all tokens and removing any undefined or empty objects or arrays.
//
// Values can only be primitives, arrays or tokens. Other objects (i.e. with methods) will be rejected.
// Experimental.
func Tokenization_Resolve(obj interface{}, options *ResolveOptions) interface{} {
	_init_.Initialize()

	var returns interface{}

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"resolve",
		[]interface{}{obj, options},
		&returns,
	)

	return returns
}

// Reverse any value into a Resolvable, if possible.
//
// In case of a string, the string must not be a concatenation.
// Experimental.
func Tokenization_Reverse(x interface{}, options *ReverseOptions) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"reverse",
		[]interface{}{x, options},
		&returns,
	)

	return returns
}

// Un-encode a string which is either a complete encoded token, or doesn't contain tokens at all.
//
// It's illegal for the string to be a concatenation of an encoded token and something else.
// Experimental.
func Tokenization_ReverseCompleteString(s *string) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"reverseCompleteString",
		[]interface{}{s},
		&returns,
	)

	return returns
}

// Un-encode a Tokenized value from a list.
// Experimental.
func Tokenization_ReverseList(l *[]*string) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"reverseList",
		[]interface{}{l},
		&returns,
	)

	return returns
}

// Un-encode a Tokenized value from a number.
// Experimental.
func Tokenization_ReverseNumber(n *float64) IResolvable {
	_init_.Initialize()

	var returns IResolvable

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"reverseNumber",
		[]interface{}{n},
		&returns,
	)

	return returns
}

// Un-encode a string potentially containing encoded tokens.
// Experimental.
func Tokenization_ReverseString(s *string) TokenizedStringFragments {
	_init_.Initialize()

	var returns TokenizedStringFragments

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"reverseString",
		[]interface{}{s},
		&returns,
	)

	return returns
}

// Stringify a number directly or lazily if it's a Token.
//
// If it is an object (i.e., { Ref: 'SomeLogicalId' }), return it as-is.
// Experimental.
func Tokenization_StringifyNumber(x *float64) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.Tokenization",
		"stringifyNumber",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Fragments of a concatenated string containing stringified Tokens.
// Experimental.
type TokenizedStringFragments interface {
	FirstToken() IResolvable
	FirstValue() interface{}
	Length() *float64
	Tokens() *[]IResolvable
	AddIntrinsic(value interface{})
	AddLiteral(lit interface{})
	AddToken(token IResolvable)
	Join(concat IFragmentConcatenator) interface{}
	MapTokens(mapper ITokenMapper) TokenizedStringFragments
}

// The jsii proxy struct for TokenizedStringFragments
type jsiiProxy_TokenizedStringFragments struct {
	_ byte // padding
}

func (j *jsiiProxy_TokenizedStringFragments) FirstToken() IResolvable {
	var returns IResolvable
	_jsii_.Get(
		j,
		"firstToken",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenizedStringFragments) FirstValue() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"firstValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenizedStringFragments) Length() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"length",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenizedStringFragments) Tokens() *[]IResolvable {
	var returns *[]IResolvable
	_jsii_.Get(
		j,
		"tokens",
		&returns,
	)
	return returns
}


// Experimental.
func NewTokenizedStringFragments() TokenizedStringFragments {
	_init_.Initialize()

	j := jsiiProxy_TokenizedStringFragments{}

	_jsii_.Create(
		"monocdk.TokenizedStringFragments",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewTokenizedStringFragments_Override(t TokenizedStringFragments) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.TokenizedStringFragments",
		nil, // no parameters
		t,
	)
}

// Experimental.
func (t *jsiiProxy_TokenizedStringFragments) AddIntrinsic(value interface{}) {
	_jsii_.InvokeVoid(
		t,
		"addIntrinsic",
		[]interface{}{value},
	)
}

// Experimental.
func (t *jsiiProxy_TokenizedStringFragments) AddLiteral(lit interface{}) {
	_jsii_.InvokeVoid(
		t,
		"addLiteral",
		[]interface{}{lit},
	)
}

// Experimental.
func (t *jsiiProxy_TokenizedStringFragments) AddToken(token IResolvable) {
	_jsii_.InvokeVoid(
		t,
		"addToken",
		[]interface{}{token},
	)
}

// Combine the string fragments using the given joiner.
//
// If there are any
// Experimental.
func (t *jsiiProxy_TokenizedStringFragments) Join(concat IFragmentConcatenator) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		t,
		"join",
		[]interface{}{concat},
		&returns,
	)

	return returns
}

// Apply a transformation function to all tokens in the string.
// Experimental.
func (t *jsiiProxy_TokenizedStringFragments) MapTokens(mapper ITokenMapper) TokenizedStringFragments {
	var returns TokenizedStringFragments

	_jsii_.Invoke(
		t,
		"mapTokens",
		[]interface{}{mapper},
		&returns,
	)

	return returns
}

// Inspector that maintains an attribute bag.
// Experimental.
type TreeInspector interface {
	Attributes() *map[string]interface{}
	AddAttribute(key *string, value interface{})
}

// The jsii proxy struct for TreeInspector
type jsiiProxy_TreeInspector struct {
	_ byte // padding
}

func (j *jsiiProxy_TreeInspector) Attributes() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"attributes",
		&returns,
	)
	return returns
}


// Experimental.
func NewTreeInspector() TreeInspector {
	_init_.Initialize()

	j := jsiiProxy_TreeInspector{}

	_jsii_.Create(
		"monocdk.TreeInspector",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewTreeInspector_Override(t TreeInspector) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.TreeInspector",
		nil, // no parameters
		t,
	)
}

// Adds attribute to bag.
//
// Keys should be added by convention to prevent conflicts
// i.e. L1 constructs will contain attributes with keys prefixed with aws:cdk:cloudformation
// Experimental.
func (t *jsiiProxy_TreeInspector) AddAttribute(key *string, value interface{}) {
	_jsii_.InvokeVoid(
		t,
		"addAttribute",
		[]interface{}{key, value},
	)
}

// An error returned during the validation phase.
// Experimental.
type ValidationError struct {
	// The error message.
	// Experimental.
	Message *string `json:"message"`
	// The construct which emitted the error.
	// Experimental.
	Source Construct `json:"source"`
}

// Representation of validation results.
//
// Models a tree of validation errors so that we have as much information as possible
// about the failure that occurred.
// Experimental.
type ValidationResult interface {
	ErrorMessage() *string
	IsSuccess() *bool
	Results() ValidationResults
	AssertSuccess()
	ErrorTree() *string
	Prefix(message *string) ValidationResult
}

// The jsii proxy struct for ValidationResult
type jsiiProxy_ValidationResult struct {
	_ byte // padding
}

func (j *jsiiProxy_ValidationResult) ErrorMessage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"errorMessage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ValidationResult) IsSuccess() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isSuccess",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ValidationResult) Results() ValidationResults {
	var returns ValidationResults
	_jsii_.Get(
		j,
		"results",
		&returns,
	)
	return returns
}


// Experimental.
func NewValidationResult(errorMessage *string, results ValidationResults) ValidationResult {
	_init_.Initialize()

	j := jsiiProxy_ValidationResult{}

	_jsii_.Create(
		"monocdk.ValidationResult",
		[]interface{}{errorMessage, results},
		&j,
	)

	return &j
}

// Experimental.
func NewValidationResult_Override(v ValidationResult, errorMessage *string, results ValidationResults) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.ValidationResult",
		[]interface{}{errorMessage, results},
		v,
	)
}

// Turn a failed validation into an exception.
// Experimental.
func (v *jsiiProxy_ValidationResult) AssertSuccess() {
	_jsii_.InvokeVoid(
		v,
		"assertSuccess",
		nil, // no parameters
	)
}

// Return a string rendering of the tree of validation failures.
// Experimental.
func (v *jsiiProxy_ValidationResult) ErrorTree() *string {
	var returns *string

	_jsii_.Invoke(
		v,
		"errorTree",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Wrap this result with an error message, if it concerns an error.
// Experimental.
func (v *jsiiProxy_ValidationResult) Prefix(message *string) ValidationResult {
	var returns ValidationResult

	_jsii_.Invoke(
		v,
		"prefix",
		[]interface{}{message},
		&returns,
	)

	return returns
}

// A collection of validation results.
// Experimental.
type ValidationResults interface {
	IsSuccess() *bool
	Results() *[]ValidationResult
	SetResults(val *[]ValidationResult)
	Collect(result ValidationResult)
	ErrorTreeList() *string
	Wrap(message *string) ValidationResult
}

// The jsii proxy struct for ValidationResults
type jsiiProxy_ValidationResults struct {
	_ byte // padding
}

func (j *jsiiProxy_ValidationResults) IsSuccess() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isSuccess",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ValidationResults) Results() *[]ValidationResult {
	var returns *[]ValidationResult
	_jsii_.Get(
		j,
		"results",
		&returns,
	)
	return returns
}


// Experimental.
func NewValidationResults(results *[]ValidationResult) ValidationResults {
	_init_.Initialize()

	j := jsiiProxy_ValidationResults{}

	_jsii_.Create(
		"monocdk.ValidationResults",
		[]interface{}{results},
		&j,
	)

	return &j
}

// Experimental.
func NewValidationResults_Override(v ValidationResults, results *[]ValidationResult) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.ValidationResults",
		[]interface{}{results},
		v,
	)
}

func (j *jsiiProxy_ValidationResults) SetResults(val *[]ValidationResult) {
	_jsii_.Set(
		j,
		"results",
		val,
	)
}

// Experimental.
func (v *jsiiProxy_ValidationResults) Collect(result ValidationResult) {
	_jsii_.InvokeVoid(
		v,
		"collect",
		[]interface{}{result},
	)
}

// Experimental.
func (v *jsiiProxy_ValidationResults) ErrorTreeList() *string {
	var returns *string

	_jsii_.Invoke(
		v,
		"errorTreeList",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Wrap up all validation results into a single tree node.
//
// If there are failures in the collection, add a message, otherwise
// return a success.
// Experimental.
func (v *jsiiProxy_ValidationResults) Wrap(message *string) ValidationResult {
	var returns ValidationResult

	_jsii_.Invoke(
		v,
		"wrap",
		[]interface{}{message},
		&returns,
	)

	return returns
}


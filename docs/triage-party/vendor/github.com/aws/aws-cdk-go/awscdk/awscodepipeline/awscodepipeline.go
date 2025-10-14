package awscodepipeline

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscodepipeline/internal"
	"github.com/aws/aws-cdk-go/awscdk/awscodestarnotifications"
	"github.com/aws/aws-cdk-go/awscdk/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/constructs-go/constructs/v3"
)

// Low-level class for generic CodePipeline Actions implementing the {@link IAction} interface.
//
// Contains some common logic that can be re-used by all {@link IAction} implementations.
// If you're writing your own Action class,
// feel free to extend this class.
// Experimental.
type Action interface {
	IAction
	ActionProperties() *ActionProperties
	ProvidedActionProperties() *ActionProperties
	Bind(scope awscdk.Construct, stage IStage, options *ActionBindOptions) *ActionConfig
	Bound(scope awscdk.Construct, stage IStage, options *ActionBindOptions) *ActionConfig
	OnStateChange(name *string, target awsevents.IRuleTarget, options *awsevents.RuleProps) awsevents.Rule
	VariableExpression(variableName *string) *string
}

// The jsii proxy struct for Action
type jsiiProxy_Action struct {
	jsiiProxy_IAction
}

func (j *jsiiProxy_Action) ActionProperties() *ActionProperties {
	var returns *ActionProperties
	_jsii_.Get(
		j,
		"actionProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Action) ProvidedActionProperties() *ActionProperties {
	var returns *ActionProperties
	_jsii_.Get(
		j,
		"providedActionProperties",
		&returns,
	)
	return returns
}


// Experimental.
func NewAction_Override(a Action) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.Action",
		nil, // no parameters
		a,
	)
}

// The callback invoked when this Action is added to a Pipeline.
// Experimental.
func (a *jsiiProxy_Action) Bind(scope awscdk.Construct, stage IStage, options *ActionBindOptions) *ActionConfig {
	var returns *ActionConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope, stage, options},
		&returns,
	)

	return returns
}

// This is a renamed version of the {@link IAction.bind} method.
// Experimental.
func (a *jsiiProxy_Action) Bound(scope awscdk.Construct, stage IStage, options *ActionBindOptions) *ActionConfig {
	var returns *ActionConfig

	_jsii_.Invoke(
		a,
		"bound",
		[]interface{}{scope, stage, options},
		&returns,
	)

	return returns
}

// Creates an Event that will be triggered whenever the state of this Action changes.
// Experimental.
func (a *jsiiProxy_Action) OnStateChange(name *string, target awsevents.IRuleTarget, options *awsevents.RuleProps) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		a,
		"onStateChange",
		[]interface{}{name, target, options},
		&returns,
	)

	return returns
}

// Experimental.
func (a *jsiiProxy_Action) VariableExpression(variableName *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"variableExpression",
		[]interface{}{variableName},
		&returns,
	)

	return returns
}

// Specifies the constraints on the number of input and output artifacts an action can have.
//
// The constraints for each action type are documented on the
// {@link https://docs.aws.amazon.com/codepipeline/latest/userguide/reference-pipeline-structure.html Pipeline Structure Reference} page.
// Experimental.
type ActionArtifactBounds struct {
	// Experimental.
	MaxInputs *float64 `json:"maxInputs"`
	// Experimental.
	MaxOutputs *float64 `json:"maxOutputs"`
	// Experimental.
	MinInputs *float64 `json:"minInputs"`
	// Experimental.
	MinOutputs *float64 `json:"minOutputs"`
}

// Experimental.
type ActionBindOptions struct {
	// Experimental.
	Bucket awss3.IBucket `json:"bucket"`
	// Experimental.
	Role awsiam.IRole `json:"role"`
}

// Experimental.
type ActionCategory string

const (
	ActionCategory_SOURCE ActionCategory = "SOURCE"
	ActionCategory_BUILD ActionCategory = "BUILD"
	ActionCategory_TEST ActionCategory = "TEST"
	ActionCategory_APPROVAL ActionCategory = "APPROVAL"
	ActionCategory_DEPLOY ActionCategory = "DEPLOY"
	ActionCategory_INVOKE ActionCategory = "INVOKE"
)

// Experimental.
type ActionConfig struct {
	// Experimental.
	Configuration interface{} `json:"configuration"`
}

// Experimental.
type ActionProperties struct {
	// Experimental.
	ActionName *string `json:"actionName"`
	// Experimental.
	ArtifactBounds *ActionArtifactBounds `json:"artifactBounds"`
	// The category of the action.
	//
	// The category defines which action type the owner
	// (the entity that performs the action) performs.
	// Experimental.
	Category ActionCategory `json:"category"`
	// The service provider that the action calls.
	// Experimental.
	Provider *string `json:"provider"`
	// The account the Action is supposed to live in.
	//
	// For Actions backed by resources,
	// this is inferred from the Stack {@link resource} is part of.
	// However, some Actions, like the CloudFormation ones,
	// are not backed by any resource, and they still might want to be cross-account.
	// In general, a concrete Action class should specify either {@link resource},
	// or {@link account} - but not both.
	// Experimental.
	Account *string `json:"account"`
	// Experimental.
	Inputs *[]Artifact `json:"inputs"`
	// Experimental.
	Outputs *[]Artifact `json:"outputs"`
	// Experimental.
	Owner *string `json:"owner"`
	// The AWS region the given Action resides in.
	//
	// Note that a cross-region Pipeline requires replication buckets to function correctly.
	// You can provide their names with the {@link PipelineProps#crossRegionReplicationBuckets} property.
	// If you don't, the CodePipeline Construct will create new Stacks in your CDK app containing those buckets,
	// that you will need to `cdk deploy` before deploying the main, Pipeline-containing Stack.
	// Experimental.
	Region *string `json:"region"`
	// The optional resource that is backing this Action.
	//
	// This is used for automatically handling Actions backed by
	// resources from a different account and/or region.
	// Experimental.
	Resource awscdk.IResource `json:"resource"`
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// The order in which AWS CodePipeline runs this action. For more information, see the AWS CodePipeline User Guide.
	//
	// https://docs.aws.amazon.com/codepipeline/latest/userguide/reference-pipeline-structure.html#action-requirements
	// Experimental.
	RunOrder *float64 `json:"runOrder"`
	// The name of the namespace to use for variables emitted by this action.
	// Experimental.
	VariablesNamespace *string `json:"variablesNamespace"`
	// Experimental.
	Version *string `json:"version"`
}

// An output artifact of an action.
//
// Artifacts can be used as input by some actions.
// Experimental.
type Artifact interface {
	ArtifactName() *string
	BucketName() *string
	ObjectKey() *string
	S3Location() *awss3.Location
	Url() *string
	AtPath(fileName *string) ArtifactPath
	GetMetadata(key *string) interface{}
	GetParam(jsonFile *string, keyName *string) *string
	SetMetadata(key *string, value interface{})
	ToString() *string
}

// The jsii proxy struct for Artifact
type jsiiProxy_Artifact struct {
	_ byte // padding
}

func (j *jsiiProxy_Artifact) ArtifactName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"artifactName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Artifact) BucketName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Artifact) ObjectKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"objectKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Artifact) S3Location() *awss3.Location {
	var returns *awss3.Location
	_jsii_.Get(
		j,
		"s3Location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Artifact) Url() *string {
	var returns *string
	_jsii_.Get(
		j,
		"url",
		&returns,
	)
	return returns
}


// Experimental.
func NewArtifact(artifactName *string) Artifact {
	_init_.Initialize()

	j := jsiiProxy_Artifact{}

	_jsii_.Create(
		"monocdk.aws_codepipeline.Artifact",
		[]interface{}{artifactName},
		&j,
	)

	return &j
}

// Experimental.
func NewArtifact_Override(a Artifact, artifactName *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.Artifact",
		[]interface{}{artifactName},
		a,
	)
}

// A static factory method used to create instances of the Artifact class.
//
// Mainly meant to be used from `decdk`.
// Experimental.
func Artifact_Artifact(name *string) Artifact {
	_init_.Initialize()

	var returns Artifact

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.Artifact",
		"artifact",
		[]interface{}{name},
		&returns,
	)

	return returns
}

// Returns an ArtifactPath for a file within this artifact.
//
// CfnOutput is in the form "<artifact-name>::<file-name>"
// Experimental.
func (a *jsiiProxy_Artifact) AtPath(fileName *string) ArtifactPath {
	var returns ArtifactPath

	_jsii_.Invoke(
		a,
		"atPath",
		[]interface{}{fileName},
		&returns,
	)

	return returns
}

// Retrieve the metadata stored in this artifact under the given key.
//
// If there is no metadata stored under the given key,
// null will be returned.
// Experimental.
func (a *jsiiProxy_Artifact) GetMetadata(key *string) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		a,
		"getMetadata",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Returns a token for a value inside a JSON file within this artifact.
// Experimental.
func (a *jsiiProxy_Artifact) GetParam(jsonFile *string, keyName *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"getParam",
		[]interface{}{jsonFile, keyName},
		&returns,
	)

	return returns
}

// Add arbitrary extra payload to the artifact under a given key.
//
// This can be used by CodePipeline actions to communicate data between themselves.
// If metadata was already present under the given key,
// it will be overwritten with the new value.
// Experimental.
func (a *jsiiProxy_Artifact) SetMetadata(key *string, value interface{}) {
	_jsii_.InvokeVoid(
		a,
		"setMetadata",
		[]interface{}{key, value},
	)
}

// Experimental.
func (a *jsiiProxy_Artifact) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A specific file within an output artifact.
//
// The most common use case for this is specifying the template file
// for a CloudFormation action.
// Experimental.
type ArtifactPath interface {
	Artifact() Artifact
	FileName() *string
	Location() *string
}

// The jsii proxy struct for ArtifactPath
type jsiiProxy_ArtifactPath struct {
	_ byte // padding
}

func (j *jsiiProxy_ArtifactPath) Artifact() Artifact {
	var returns Artifact
	_jsii_.Get(
		j,
		"artifact",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ArtifactPath) FileName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fileName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ArtifactPath) Location() *string {
	var returns *string
	_jsii_.Get(
		j,
		"location",
		&returns,
	)
	return returns
}


// Experimental.
func NewArtifactPath(artifact Artifact, fileName *string) ArtifactPath {
	_init_.Initialize()

	j := jsiiProxy_ArtifactPath{}

	_jsii_.Create(
		"monocdk.aws_codepipeline.ArtifactPath",
		[]interface{}{artifact, fileName},
		&j,
	)

	return &j
}

// Experimental.
func NewArtifactPath_Override(a ArtifactPath, artifact Artifact, fileName *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.ArtifactPath",
		[]interface{}{artifact, fileName},
		a,
	)
}

// Experimental.
func ArtifactPath_ArtifactPath(artifactName *string, fileName *string) ArtifactPath {
	_init_.Initialize()

	var returns ArtifactPath

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.ArtifactPath",
		"artifactPath",
		[]interface{}{artifactName, fileName},
		&returns,
	)

	return returns
}

// A CloudFormation `AWS::CodePipeline::CustomActionType`.
type CfnCustomActionType interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Category() *string
	SetCategory(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ConfigurationProperties() interface{}
	SetConfigurationProperties(val interface{})
	CreationStack() *[]*string
	InputArtifactDetails() interface{}
	SetInputArtifactDetails(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	OutputArtifactDetails() interface{}
	SetOutputArtifactDetails(val interface{})
	Provider() *string
	SetProvider(val *string)
	Ref() *string
	Settings() interface{}
	SetSettings(val interface{})
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	UpdatedProperites() *map[string]interface{}
	Version() *string
	SetVersion(val *string)
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

// The jsii proxy struct for CfnCustomActionType
type jsiiProxy_CfnCustomActionType struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnCustomActionType) Category() *string {
	var returns *string
	_jsii_.Get(
		j,
		"category",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) ConfigurationProperties() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"configurationProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) InputArtifactDetails() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"inputArtifactDetails",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) OutputArtifactDetails() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"outputArtifactDetails",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) Provider() *string {
	var returns *string
	_jsii_.Get(
		j,
		"provider",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) Settings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"settings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnCustomActionType) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodePipeline::CustomActionType`.
func NewCfnCustomActionType(scope awscdk.Construct, id *string, props *CfnCustomActionTypeProps) CfnCustomActionType {
	_init_.Initialize()

	j := jsiiProxy_CfnCustomActionType{}

	_jsii_.Create(
		"monocdk.aws_codepipeline.CfnCustomActionType",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodePipeline::CustomActionType`.
func NewCfnCustomActionType_Override(c CfnCustomActionType, scope awscdk.Construct, id *string, props *CfnCustomActionTypeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.CfnCustomActionType",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnCustomActionType) SetCategory(val *string) {
	_jsii_.Set(
		j,
		"category",
		val,
	)
}

func (j *jsiiProxy_CfnCustomActionType) SetConfigurationProperties(val interface{}) {
	_jsii_.Set(
		j,
		"configurationProperties",
		val,
	)
}

func (j *jsiiProxy_CfnCustomActionType) SetInputArtifactDetails(val interface{}) {
	_jsii_.Set(
		j,
		"inputArtifactDetails",
		val,
	)
}

func (j *jsiiProxy_CfnCustomActionType) SetOutputArtifactDetails(val interface{}) {
	_jsii_.Set(
		j,
		"outputArtifactDetails",
		val,
	)
}

func (j *jsiiProxy_CfnCustomActionType) SetProvider(val *string) {
	_jsii_.Set(
		j,
		"provider",
		val,
	)
}

func (j *jsiiProxy_CfnCustomActionType) SetSettings(val interface{}) {
	_jsii_.Set(
		j,
		"settings",
		val,
	)
}

func (j *jsiiProxy_CfnCustomActionType) SetVersion(val *string) {
	_jsii_.Set(
		j,
		"version",
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
func CfnCustomActionType_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnCustomActionType",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnCustomActionType_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnCustomActionType",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnCustomActionType_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnCustomActionType",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnCustomActionType_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codepipeline.CfnCustomActionType",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnCustomActionType) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnCustomActionType) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnCustomActionType) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnCustomActionType) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnCustomActionType) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnCustomActionType) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnCustomActionType) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnCustomActionType) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnCustomActionType) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnCustomActionType) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnCustomActionType) OnPrepare() {
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
func (c *jsiiProxy_CfnCustomActionType) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCustomActionType) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnCustomActionType) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnCustomActionType) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnCustomActionType) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnCustomActionType) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnCustomActionType) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnCustomActionType) ToString() *string {
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
func (c *jsiiProxy_CfnCustomActionType) Validate() *[]*string {
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
func (c *jsiiProxy_CfnCustomActionType) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnCustomActionType_ArtifactDetailsProperty struct {
	// `CfnCustomActionType.ArtifactDetailsProperty.MaximumCount`.
	MaximumCount *float64 `json:"maximumCount"`
	// `CfnCustomActionType.ArtifactDetailsProperty.MinimumCount`.
	MinimumCount *float64 `json:"minimumCount"`
}

type CfnCustomActionType_ConfigurationPropertiesProperty struct {
	// `CfnCustomActionType.ConfigurationPropertiesProperty.Key`.
	Key interface{} `json:"key"`
	// `CfnCustomActionType.ConfigurationPropertiesProperty.Name`.
	Name *string `json:"name"`
	// `CfnCustomActionType.ConfigurationPropertiesProperty.Required`.
	Required interface{} `json:"required"`
	// `CfnCustomActionType.ConfigurationPropertiesProperty.Secret`.
	Secret interface{} `json:"secret"`
	// `CfnCustomActionType.ConfigurationPropertiesProperty.Description`.
	Description *string `json:"description"`
	// `CfnCustomActionType.ConfigurationPropertiesProperty.Queryable`.
	Queryable interface{} `json:"queryable"`
	// `CfnCustomActionType.ConfigurationPropertiesProperty.Type`.
	Type *string `json:"type"`
}

type CfnCustomActionType_SettingsProperty struct {
	// `CfnCustomActionType.SettingsProperty.EntityUrlTemplate`.
	EntityUrlTemplate *string `json:"entityUrlTemplate"`
	// `CfnCustomActionType.SettingsProperty.ExecutionUrlTemplate`.
	ExecutionUrlTemplate *string `json:"executionUrlTemplate"`
	// `CfnCustomActionType.SettingsProperty.RevisionUrlTemplate`.
	RevisionUrlTemplate *string `json:"revisionUrlTemplate"`
	// `CfnCustomActionType.SettingsProperty.ThirdPartyConfigurationUrl`.
	ThirdPartyConfigurationUrl *string `json:"thirdPartyConfigurationUrl"`
}

// Properties for defining a `AWS::CodePipeline::CustomActionType`.
type CfnCustomActionTypeProps struct {
	// `AWS::CodePipeline::CustomActionType.Category`.
	Category *string `json:"category"`
	// `AWS::CodePipeline::CustomActionType.InputArtifactDetails`.
	InputArtifactDetails interface{} `json:"inputArtifactDetails"`
	// `AWS::CodePipeline::CustomActionType.OutputArtifactDetails`.
	OutputArtifactDetails interface{} `json:"outputArtifactDetails"`
	// `AWS::CodePipeline::CustomActionType.Provider`.
	Provider *string `json:"provider"`
	// `AWS::CodePipeline::CustomActionType.Version`.
	Version *string `json:"version"`
	// `AWS::CodePipeline::CustomActionType.ConfigurationProperties`.
	ConfigurationProperties interface{} `json:"configurationProperties"`
	// `AWS::CodePipeline::CustomActionType.Settings`.
	Settings interface{} `json:"settings"`
	// `AWS::CodePipeline::CustomActionType.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::CodePipeline::Pipeline`.
type CfnPipeline interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ArtifactStore() interface{}
	SetArtifactStore(val interface{})
	ArtifactStores() interface{}
	SetArtifactStores(val interface{})
	AttrVersion() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DisableInboundStageTransitions() interface{}
	SetDisableInboundStageTransitions(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	RestartExecutionOnUpdate() interface{}
	SetRestartExecutionOnUpdate(val interface{})
	RoleArn() *string
	SetRoleArn(val *string)
	Stack() awscdk.Stack
	Stages() interface{}
	SetStages(val interface{})
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

// The jsii proxy struct for CfnPipeline
type jsiiProxy_CfnPipeline struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnPipeline) ArtifactStore() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"artifactStore",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) ArtifactStores() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"artifactStores",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) AttrVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) DisableInboundStageTransitions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"disableInboundStageTransitions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) RestartExecutionOnUpdate() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"restartExecutionOnUpdate",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) Stages() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stages",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPipeline) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodePipeline::Pipeline`.
func NewCfnPipeline(scope awscdk.Construct, id *string, props *CfnPipelineProps) CfnPipeline {
	_init_.Initialize()

	j := jsiiProxy_CfnPipeline{}

	_jsii_.Create(
		"monocdk.aws_codepipeline.CfnPipeline",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodePipeline::Pipeline`.
func NewCfnPipeline_Override(c CfnPipeline, scope awscdk.Construct, id *string, props *CfnPipelineProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.CfnPipeline",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnPipeline) SetArtifactStore(val interface{}) {
	_jsii_.Set(
		j,
		"artifactStore",
		val,
	)
}

func (j *jsiiProxy_CfnPipeline) SetArtifactStores(val interface{}) {
	_jsii_.Set(
		j,
		"artifactStores",
		val,
	)
}

func (j *jsiiProxy_CfnPipeline) SetDisableInboundStageTransitions(val interface{}) {
	_jsii_.Set(
		j,
		"disableInboundStageTransitions",
		val,
	)
}

func (j *jsiiProxy_CfnPipeline) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnPipeline) SetRestartExecutionOnUpdate(val interface{}) {
	_jsii_.Set(
		j,
		"restartExecutionOnUpdate",
		val,
	)
}

func (j *jsiiProxy_CfnPipeline) SetRoleArn(val *string) {
	_jsii_.Set(
		j,
		"roleArn",
		val,
	)
}

func (j *jsiiProxy_CfnPipeline) SetStages(val interface{}) {
	_jsii_.Set(
		j,
		"stages",
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
func CfnPipeline_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnPipeline",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnPipeline_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnPipeline",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnPipeline_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnPipeline",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnPipeline_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codepipeline.CfnPipeline",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnPipeline) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnPipeline) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnPipeline) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnPipeline) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnPipeline) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnPipeline) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnPipeline) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnPipeline) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnPipeline) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnPipeline) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnPipeline) OnPrepare() {
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
func (c *jsiiProxy_CfnPipeline) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPipeline) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnPipeline) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnPipeline) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnPipeline) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnPipeline) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnPipeline) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPipeline) ToString() *string {
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
func (c *jsiiProxy_CfnPipeline) Validate() *[]*string {
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
func (c *jsiiProxy_CfnPipeline) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnPipeline_ActionDeclarationProperty struct {
	// `CfnPipeline.ActionDeclarationProperty.ActionTypeId`.
	ActionTypeId interface{} `json:"actionTypeId"`
	// `CfnPipeline.ActionDeclarationProperty.Name`.
	Name *string `json:"name"`
	// `CfnPipeline.ActionDeclarationProperty.Configuration`.
	Configuration interface{} `json:"configuration"`
	// `CfnPipeline.ActionDeclarationProperty.InputArtifacts`.
	InputArtifacts interface{} `json:"inputArtifacts"`
	// `CfnPipeline.ActionDeclarationProperty.Namespace`.
	Namespace *string `json:"namespace"`
	// `CfnPipeline.ActionDeclarationProperty.OutputArtifacts`.
	OutputArtifacts interface{} `json:"outputArtifacts"`
	// `CfnPipeline.ActionDeclarationProperty.Region`.
	Region *string `json:"region"`
	// `CfnPipeline.ActionDeclarationProperty.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `CfnPipeline.ActionDeclarationProperty.RunOrder`.
	RunOrder *float64 `json:"runOrder"`
}

type CfnPipeline_ActionTypeIdProperty struct {
	// `CfnPipeline.ActionTypeIdProperty.Category`.
	Category *string `json:"category"`
	// `CfnPipeline.ActionTypeIdProperty.Owner`.
	Owner *string `json:"owner"`
	// `CfnPipeline.ActionTypeIdProperty.Provider`.
	Provider *string `json:"provider"`
	// `CfnPipeline.ActionTypeIdProperty.Version`.
	Version *string `json:"version"`
}

type CfnPipeline_ArtifactStoreMapProperty struct {
	// `CfnPipeline.ArtifactStoreMapProperty.ArtifactStore`.
	ArtifactStore interface{} `json:"artifactStore"`
	// `CfnPipeline.ArtifactStoreMapProperty.Region`.
	Region *string `json:"region"`
}

type CfnPipeline_ArtifactStoreProperty struct {
	// `CfnPipeline.ArtifactStoreProperty.Location`.
	Location *string `json:"location"`
	// `CfnPipeline.ArtifactStoreProperty.Type`.
	Type *string `json:"type"`
	// `CfnPipeline.ArtifactStoreProperty.EncryptionKey`.
	EncryptionKey interface{} `json:"encryptionKey"`
}

type CfnPipeline_BlockerDeclarationProperty struct {
	// `CfnPipeline.BlockerDeclarationProperty.Name`.
	Name *string `json:"name"`
	// `CfnPipeline.BlockerDeclarationProperty.Type`.
	Type *string `json:"type"`
}

type CfnPipeline_EncryptionKeyProperty struct {
	// `CfnPipeline.EncryptionKeyProperty.Id`.
	Id *string `json:"id"`
	// `CfnPipeline.EncryptionKeyProperty.Type`.
	Type *string `json:"type"`
}

type CfnPipeline_InputArtifactProperty struct {
	// `CfnPipeline.InputArtifactProperty.Name`.
	Name *string `json:"name"`
}

type CfnPipeline_OutputArtifactProperty struct {
	// `CfnPipeline.OutputArtifactProperty.Name`.
	Name *string `json:"name"`
}

type CfnPipeline_StageDeclarationProperty struct {
	// `CfnPipeline.StageDeclarationProperty.Actions`.
	Actions interface{} `json:"actions"`
	// `CfnPipeline.StageDeclarationProperty.Name`.
	Name *string `json:"name"`
	// `CfnPipeline.StageDeclarationProperty.Blockers`.
	Blockers interface{} `json:"blockers"`
}

type CfnPipeline_StageTransitionProperty struct {
	// `CfnPipeline.StageTransitionProperty.Reason`.
	Reason *string `json:"reason"`
	// `CfnPipeline.StageTransitionProperty.StageName`.
	StageName *string `json:"stageName"`
}

// Properties for defining a `AWS::CodePipeline::Pipeline`.
type CfnPipelineProps struct {
	// `AWS::CodePipeline::Pipeline.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `AWS::CodePipeline::Pipeline.Stages`.
	Stages interface{} `json:"stages"`
	// `AWS::CodePipeline::Pipeline.ArtifactStore`.
	ArtifactStore interface{} `json:"artifactStore"`
	// `AWS::CodePipeline::Pipeline.ArtifactStores`.
	ArtifactStores interface{} `json:"artifactStores"`
	// `AWS::CodePipeline::Pipeline.DisableInboundStageTransitions`.
	DisableInboundStageTransitions interface{} `json:"disableInboundStageTransitions"`
	// `AWS::CodePipeline::Pipeline.Name`.
	Name *string `json:"name"`
	// `AWS::CodePipeline::Pipeline.RestartExecutionOnUpdate`.
	RestartExecutionOnUpdate interface{} `json:"restartExecutionOnUpdate"`
	// `AWS::CodePipeline::Pipeline.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::CodePipeline::Webhook`.
type CfnWebhook interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrUrl() *string
	Authentication() *string
	SetAuthentication(val *string)
	AuthenticationConfiguration() interface{}
	SetAuthenticationConfiguration(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Filters() interface{}
	SetFilters(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	RegisterWithThirdParty() interface{}
	SetRegisterWithThirdParty(val interface{})
	Stack() awscdk.Stack
	TargetAction() *string
	SetTargetAction(val *string)
	TargetPipeline() *string
	SetTargetPipeline(val *string)
	TargetPipelineVersion() *float64
	SetTargetPipelineVersion(val *float64)
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

// The jsii proxy struct for CfnWebhook
type jsiiProxy_CfnWebhook struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnWebhook) AttrUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) Authentication() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authentication",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) AuthenticationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"authenticationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) Filters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"filters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) RegisterWithThirdParty() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"registerWithThirdParty",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) TargetAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) TargetPipeline() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetPipeline",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) TargetPipelineVersion() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"targetPipelineVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnWebhook) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodePipeline::Webhook`.
func NewCfnWebhook(scope awscdk.Construct, id *string, props *CfnWebhookProps) CfnWebhook {
	_init_.Initialize()

	j := jsiiProxy_CfnWebhook{}

	_jsii_.Create(
		"monocdk.aws_codepipeline.CfnWebhook",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodePipeline::Webhook`.
func NewCfnWebhook_Override(c CfnWebhook, scope awscdk.Construct, id *string, props *CfnWebhookProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.CfnWebhook",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnWebhook) SetAuthentication(val *string) {
	_jsii_.Set(
		j,
		"authentication",
		val,
	)
}

func (j *jsiiProxy_CfnWebhook) SetAuthenticationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"authenticationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnWebhook) SetFilters(val interface{}) {
	_jsii_.Set(
		j,
		"filters",
		val,
	)
}

func (j *jsiiProxy_CfnWebhook) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnWebhook) SetRegisterWithThirdParty(val interface{}) {
	_jsii_.Set(
		j,
		"registerWithThirdParty",
		val,
	)
}

func (j *jsiiProxy_CfnWebhook) SetTargetAction(val *string) {
	_jsii_.Set(
		j,
		"targetAction",
		val,
	)
}

func (j *jsiiProxy_CfnWebhook) SetTargetPipeline(val *string) {
	_jsii_.Set(
		j,
		"targetPipeline",
		val,
	)
}

func (j *jsiiProxy_CfnWebhook) SetTargetPipelineVersion(val *float64) {
	_jsii_.Set(
		j,
		"targetPipelineVersion",
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
func CfnWebhook_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnWebhook",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnWebhook_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnWebhook",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnWebhook_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.CfnWebhook",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnWebhook_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codepipeline.CfnWebhook",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnWebhook) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnWebhook) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnWebhook) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnWebhook) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnWebhook) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnWebhook) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnWebhook) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnWebhook) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnWebhook) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnWebhook) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnWebhook) OnPrepare() {
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
func (c *jsiiProxy_CfnWebhook) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnWebhook) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnWebhook) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnWebhook) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnWebhook) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnWebhook) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnWebhook) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnWebhook) ToString() *string {
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
func (c *jsiiProxy_CfnWebhook) Validate() *[]*string {
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
func (c *jsiiProxy_CfnWebhook) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnWebhook_WebhookAuthConfigurationProperty struct {
	// `CfnWebhook.WebhookAuthConfigurationProperty.AllowedIPRange`.
	AllowedIpRange *string `json:"allowedIpRange"`
	// `CfnWebhook.WebhookAuthConfigurationProperty.SecretToken`.
	SecretToken *string `json:"secretToken"`
}

type CfnWebhook_WebhookFilterRuleProperty struct {
	// `CfnWebhook.WebhookFilterRuleProperty.JsonPath`.
	JsonPath *string `json:"jsonPath"`
	// `CfnWebhook.WebhookFilterRuleProperty.MatchEquals`.
	MatchEquals *string `json:"matchEquals"`
}

// Properties for defining a `AWS::CodePipeline::Webhook`.
type CfnWebhookProps struct {
	// `AWS::CodePipeline::Webhook.Authentication`.
	Authentication *string `json:"authentication"`
	// `AWS::CodePipeline::Webhook.AuthenticationConfiguration`.
	AuthenticationConfiguration interface{} `json:"authenticationConfiguration"`
	// `AWS::CodePipeline::Webhook.Filters`.
	Filters interface{} `json:"filters"`
	// `AWS::CodePipeline::Webhook.TargetAction`.
	TargetAction *string `json:"targetAction"`
	// `AWS::CodePipeline::Webhook.TargetPipeline`.
	TargetPipeline *string `json:"targetPipeline"`
	// `AWS::CodePipeline::Webhook.TargetPipelineVersion`.
	TargetPipelineVersion *float64 `json:"targetPipelineVersion"`
	// `AWS::CodePipeline::Webhook.Name`.
	Name *string `json:"name"`
	// `AWS::CodePipeline::Webhook.RegisterWithThirdParty`.
	RegisterWithThirdParty interface{} `json:"registerWithThirdParty"`
}

// Common properties shared by all Actions.
// Experimental.
type CommonActionProps struct {
	// The physical, human-readable name of the Action.
	//
	// Note that Action names must be unique within a single Stage.
	// Experimental.
	ActionName *string `json:"actionName"`
	// The runOrder property for this Action.
	//
	// RunOrder determines the relative order in which multiple Actions in the same Stage execute.
	// See: https://docs.aws.amazon.com/codepipeline/latest/userguide/reference-pipeline-structure.html
	//
	// Experimental.
	RunOrder *float64 `json:"runOrder"`
	// The name of the namespace to use for variables emitted by this action.
	// Experimental.
	VariablesNamespace *string `json:"variablesNamespace"`
}

// Common properties shared by all Actions whose {@link ActionProperties.owner} field is 'AWS' (or unset, as 'AWS' is the default).
// Experimental.
type CommonAwsActionProps struct {
	// The physical, human-readable name of the Action.
	//
	// Note that Action names must be unique within a single Stage.
	// Experimental.
	ActionName *string `json:"actionName"`
	// The runOrder property for this Action.
	//
	// RunOrder determines the relative order in which multiple Actions in the same Stage execute.
	// See: https://docs.aws.amazon.com/codepipeline/latest/userguide/reference-pipeline-structure.html
	//
	// Experimental.
	RunOrder *float64 `json:"runOrder"`
	// The name of the namespace to use for variables emitted by this action.
	// Experimental.
	VariablesNamespace *string `json:"variablesNamespace"`
	// The Role in which context's this Action will be executing in.
	//
	// The Pipeline's Role will assume this Role
	// (the required permissions for that will be granted automatically)
	// right before executing this Action.
	// This Action will be passed into your {@link IAction.bind}
	// method in the {@link ActionBindOptions.role} property.
	// Experimental.
	Role awsiam.IRole `json:"role"`
}

// An interface representing resources generated in order to support the cross-region capabilities of CodePipeline.
//
// You get instances of this interface from the {@link Pipeline#crossRegionSupport} property.
// Experimental.
type CrossRegionSupport struct {
	// The replication Bucket used by CodePipeline to operate in this region.
	//
	// Belongs to {@link stack}.
	// Experimental.
	ReplicationBucket awss3.IBucket `json:"replicationBucket"`
	// The Stack that has been created to house the replication Bucket required for this  region.
	// Experimental.
	Stack awscdk.Stack `json:"stack"`
}

// The CodePipeline variables that are global, not bound to a specific action.
//
// This class defines a bunch of static fields that represent the different variables.
// These can be used can be used in any action configuration.
// Experimental.
type GlobalVariables interface {
}

// The jsii proxy struct for GlobalVariables
type jsiiProxy_GlobalVariables struct {
	_ byte // padding
}

// Experimental.
func NewGlobalVariables() GlobalVariables {
	_init_.Initialize()

	j := jsiiProxy_GlobalVariables{}

	_jsii_.Create(
		"monocdk.aws_codepipeline.GlobalVariables",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewGlobalVariables_Override(g GlobalVariables) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.GlobalVariables",
		nil, // no parameters
		g,
	)
}

func GlobalVariables_ExecutionId() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codepipeline.GlobalVariables",
		"executionId",
		&returns,
	)
	return returns
}

// A Pipeline Action.
//
// If you want to implement this interface,
// consider extending the {@link Action} class,
// which contains some common logic.
// Experimental.
type IAction interface {
	// The callback invoked when this Action is added to a Pipeline.
	// Experimental.
	Bind(scope awscdk.Construct, stage IStage, options *ActionBindOptions) *ActionConfig
	// Creates an Event that will be triggered whenever the state of this Action changes.
	// Experimental.
	OnStateChange(name *string, target awsevents.IRuleTarget, options *awsevents.RuleProps) awsevents.Rule
	// The simple properties of the Action, like its Owner, name, etc.
	//
	// Note that this accessor will be called before the {@link bind} callback.
	// Experimental.
	ActionProperties() *ActionProperties
}

// The jsii proxy for IAction
type jsiiProxy_IAction struct {
	_ byte // padding
}

func (i *jsiiProxy_IAction) Bind(scope awscdk.Construct, stage IStage, options *ActionBindOptions) *ActionConfig {
	var returns *ActionConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, stage, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IAction) OnStateChange(name *string, target awsevents.IRuleTarget, options *awsevents.RuleProps) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onStateChange",
		[]interface{}{name, target, options},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IAction) ActionProperties() *ActionProperties {
	var returns *ActionProperties
	_jsii_.Get(
		j,
		"actionProperties",
		&returns,
	)
	return returns
}

// The abstract view of an AWS CodePipeline as required and used by Actions.
//
// It extends {@link events.IRuleTarget},
// so this interface can be used as a Target for CloudWatch Events.
// Experimental.
type IPipeline interface {
	awscodestarnotifications.INotificationRuleSource
	awscdk.IResource
	// Defines a CodeStar notification rule triggered when the pipeline events emitted by you specified, it very similar to `onEvent` API.
	//
	// You can also use the methods `notifyOnExecutionStateChange`, `notifyOnAnyStageStateChange`,
	// `notifyOnAnyActionStateChange` and `notifyOnAnyManualApprovalStateChange`
	// to define rules for these specific event emitted.
	//
	// Returns: CodeStar notification rule associated with this build project.
	// Experimental.
	NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *PipelineNotifyOnOptions) awscodestarnotifications.INotificationRule
	// Define an notification rule triggered by the set of the "Action execution" events emitted from this pipeline.
	// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#events-ref-pipeline
	//
	// Experimental.
	NotifyOnAnyActionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	// Define an notification rule triggered by the set of the "Manual approval" events emitted from this pipeline.
	// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#events-ref-pipeline
	//
	// Experimental.
	NotifyOnAnyManualApprovalStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	// Define an notification rule triggered by the set of the "Stage execution" events emitted from this pipeline.
	// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#events-ref-pipeline
	//
	// Experimental.
	NotifyOnAnyStageStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	// Define an notification rule triggered by the set of the "Pipeline execution" events emitted from this pipeline.
	// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#events-ref-pipeline
	//
	// Experimental.
	NotifyOnExecutionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	// Define an event rule triggered by this CodePipeline.
	// Experimental.
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Define an event rule triggered by the "CodePipeline Pipeline Execution State Change" event emitted from this pipeline.
	// Experimental.
	OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// The ARN of the Pipeline.
	// Experimental.
	PipelineArn() *string
	// The name of the Pipeline.
	// Experimental.
	PipelineName() *string
}

// The jsii proxy for IPipeline
type jsiiProxy_IPipeline struct {
	internal.Type__awscodestarnotificationsINotificationRuleSource
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IPipeline) NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *PipelineNotifyOnOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOn",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPipeline) NotifyOnAnyActionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOnAnyActionStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPipeline) NotifyOnAnyManualApprovalStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOnAnyManualApprovalStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPipeline) NotifyOnAnyStageStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOnAnyStageStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPipeline) NotifyOnExecutionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		i,
		"notifyOnExecutionStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPipeline) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPipeline) OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onStateChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPipeline) BindAsNotificationRuleSource(scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig {
	var returns *awscodestarnotifications.NotificationRuleSourceConfig

	_jsii_.Invoke(
		i,
		"bindAsNotificationRuleSource",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IPipeline) PipelineArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pipelineArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IPipeline) PipelineName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pipelineName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IPipeline) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IPipeline) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IPipeline) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// The abstract interface of a Pipeline Stage that is used by Actions.
// Experimental.
type IStage interface {
	// Experimental.
	AddAction(action IAction)
	// Experimental.
	OnStateChange(name *string, target awsevents.IRuleTarget, options *awsevents.RuleProps) awsevents.Rule
	// The actions belonging to this stage.
	// Experimental.
	Actions() *[]IAction
	// Experimental.
	Pipeline() IPipeline
	// The physical, human-readable name of this Pipeline Stage.
	// Experimental.
	StageName() *string
}

// The jsii proxy for IStage
type jsiiProxy_IStage struct {
	_ byte // padding
}

func (i *jsiiProxy_IStage) AddAction(action IAction) {
	_jsii_.InvokeVoid(
		i,
		"addAction",
		[]interface{}{action},
	)
}

func (i *jsiiProxy_IStage) OnStateChange(name *string, target awsevents.IRuleTarget, options *awsevents.RuleProps) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onStateChange",
		[]interface{}{name, target, options},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IStage) Actions() *[]IAction {
	var returns *[]IAction
	_jsii_.Get(
		j,
		"actions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStage) Pipeline() IPipeline {
	var returns IPipeline
	_jsii_.Get(
		j,
		"pipeline",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStage) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}

// An AWS CodePipeline pipeline with its associated IAM role and S3 bucket.
//
// TODO: EXAMPLE
//
// Experimental.
type Pipeline interface {
	awscdk.Resource
	IPipeline
	ArtifactBucket() awss3.IBucket
	CrossRegionSupport() *map[string]*CrossRegionSupport
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	PipelineArn() *string
	PipelineName() *string
	PipelineVersion() *string
	Role() awsiam.IRole
	Stack() awscdk.Stack
	StageCount() *float64
	Stages() *[]IStage
	AddStage(props *StageOptions) IStage
	AddToRolePolicy(statement awsiam.PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	BindAsNotificationRuleSource(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *PipelineNotifyOnOptions) awscodestarnotifications.INotificationRule
	NotifyOnAnyActionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	NotifyOnAnyManualApprovalStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	NotifyOnAnyStageStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	NotifyOnExecutionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnPrepare()
	OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Stage(stageName *string) IStage
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Pipeline
type jsiiProxy_Pipeline struct {
	internal.Type__awscdkResource
	jsiiProxy_IPipeline
}

func (j *jsiiProxy_Pipeline) ArtifactBucket() awss3.IBucket {
	var returns awss3.IBucket
	_jsii_.Get(
		j,
		"artifactBucket",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) CrossRegionSupport() *map[string]*CrossRegionSupport {
	var returns *map[string]*CrossRegionSupport
	_jsii_.Get(
		j,
		"crossRegionSupport",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) PipelineArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pipelineArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) PipelineName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pipelineName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) PipelineVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pipelineVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) StageCount() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"stageCount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Pipeline) Stages() *[]IStage {
	var returns *[]IStage
	_jsii_.Get(
		j,
		"stages",
		&returns,
	)
	return returns
}


// Experimental.
func NewPipeline(scope constructs.Construct, id *string, props *PipelineProps) Pipeline {
	_init_.Initialize()

	j := jsiiProxy_Pipeline{}

	_jsii_.Create(
		"monocdk.aws_codepipeline.Pipeline",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewPipeline_Override(p Pipeline, scope constructs.Construct, id *string, props *PipelineProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codepipeline.Pipeline",
		[]interface{}{scope, id, props},
		p,
	)
}

// Import a pipeline into this app.
// Experimental.
func Pipeline_FromPipelineArn(scope constructs.Construct, id *string, pipelineArn *string) IPipeline {
	_init_.Initialize()

	var returns IPipeline

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.Pipeline",
		"fromPipelineArn",
		[]interface{}{scope, id, pipelineArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Pipeline_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.Pipeline",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Pipeline_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codepipeline.Pipeline",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Creates a new Stage, and adds it to this Pipeline.
//
// Returns: the newly created Stage
// Experimental.
func (p *jsiiProxy_Pipeline) AddStage(props *StageOptions) IStage {
	var returns IStage

	_jsii_.Invoke(
		p,
		"addStage",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Adds a statement to the pipeline role.
// Experimental.
func (p *jsiiProxy_Pipeline) AddToRolePolicy(statement awsiam.PolicyStatement) {
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
func (p *jsiiProxy_Pipeline) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		p,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Returns a source configuration for notification rule.
// Experimental.
func (p *jsiiProxy_Pipeline) BindAsNotificationRuleSource(_scope constructs.Construct) *awscodestarnotifications.NotificationRuleSourceConfig {
	var returns *awscodestarnotifications.NotificationRuleSourceConfig

	_jsii_.Invoke(
		p,
		"bindAsNotificationRuleSource",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// Experimental.
func (p *jsiiProxy_Pipeline) GeneratePhysicalName() *string {
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
func (p *jsiiProxy_Pipeline) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (p *jsiiProxy_Pipeline) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Defines a CodeStar notification rule triggered when the pipeline events emitted by you specified, it very similar to `onEvent` API.
//
// You can also use the methods `notifyOnExecutionStateChange`, `notifyOnAnyStageStateChange`,
// `notifyOnAnyActionStateChange` and `notifyOnAnyManualApprovalStateChange`
// to define rules for these specific event emitted.
// Experimental.
func (p *jsiiProxy_Pipeline) NotifyOn(id *string, target awscodestarnotifications.INotificationRuleTarget, options *PipelineNotifyOnOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOn",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Define an notification rule triggered by the set of the "Action execution" events emitted from this pipeline.
// Experimental.
func (p *jsiiProxy_Pipeline) NotifyOnAnyActionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnAnyActionStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Define an notification rule triggered by the set of the "Manual approval" events emitted from this pipeline.
// Experimental.
func (p *jsiiProxy_Pipeline) NotifyOnAnyManualApprovalStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnAnyManualApprovalStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Define an notification rule triggered by the set of the "Stage execution" events emitted from this pipeline.
// Experimental.
func (p *jsiiProxy_Pipeline) NotifyOnAnyStageStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnAnyStageStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Define an notification rule triggered by the set of the "Pipeline execution" events emitted from this pipeline.
// Experimental.
func (p *jsiiProxy_Pipeline) NotifyOnExecutionStateChange(id *string, target awscodestarnotifications.INotificationRuleTarget, options *awscodestarnotifications.NotificationRuleOptions) awscodestarnotifications.INotificationRule {
	var returns awscodestarnotifications.INotificationRule

	_jsii_.Invoke(
		p,
		"notifyOnExecutionStateChange",
		[]interface{}{id, target, options},
		&returns,
	)

	return returns
}

// Defines an event rule triggered by this CodePipeline.
// Experimental.
func (p *jsiiProxy_Pipeline) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		p,
		"onEvent",
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
func (p *jsiiProxy_Pipeline) OnPrepare() {
	_jsii_.InvokeVoid(
		p,
		"onPrepare",
		nil, // no parameters
	)
}

// Defines an event rule triggered by the "CodePipeline Pipeline Execution State Change" event emitted from this pipeline.
// Experimental.
func (p *jsiiProxy_Pipeline) OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
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
func (p *jsiiProxy_Pipeline) OnSynthesize(session constructs.ISynthesisSession) {
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
func (p *jsiiProxy_Pipeline) OnValidate() *[]*string {
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
func (p *jsiiProxy_Pipeline) Prepare() {
	_jsii_.InvokeVoid(
		p,
		"prepare",
		nil, // no parameters
	)
}

// Access one of the pipeline's stages by stage name.
// Experimental.
func (p *jsiiProxy_Pipeline) Stage(stageName *string) IStage {
	var returns IStage

	_jsii_.Invoke(
		p,
		"stage",
		[]interface{}{stageName},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_Pipeline) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_Pipeline) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the pipeline structure.
//
// Validation happens according to the rules documented at
//
// https://docs.aws.amazon.com/codepipeline/latest/userguide/reference-pipeline-structure.html#pipeline-requirements
// Experimental.
func (p *jsiiProxy_Pipeline) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The list of event types for AWS Codepipeline Pipeline.
// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#events-ref-pipeline
//
// Experimental.
type PipelineNotificationEvents string

const (
	PipelineNotificationEvents_PIPELINE_EXECUTION_FAILED PipelineNotificationEvents = "PIPELINE_EXECUTION_FAILED"
	PipelineNotificationEvents_PIPELINE_EXECUTION_CANCELED PipelineNotificationEvents = "PIPELINE_EXECUTION_CANCELED"
	PipelineNotificationEvents_PIPELINE_EXECUTION_STARTED PipelineNotificationEvents = "PIPELINE_EXECUTION_STARTED"
	PipelineNotificationEvents_PIPELINE_EXECUTION_RESUMED PipelineNotificationEvents = "PIPELINE_EXECUTION_RESUMED"
	PipelineNotificationEvents_PIPELINE_EXECUTION_SUCCEEDED PipelineNotificationEvents = "PIPELINE_EXECUTION_SUCCEEDED"
	PipelineNotificationEvents_PIPELINE_EXECUTION_SUPERSEDED PipelineNotificationEvents = "PIPELINE_EXECUTION_SUPERSEDED"
	PipelineNotificationEvents_STAGE_EXECUTION_STARTED PipelineNotificationEvents = "STAGE_EXECUTION_STARTED"
	PipelineNotificationEvents_STAGE_EXECUTION_SUCCEEDED PipelineNotificationEvents = "STAGE_EXECUTION_SUCCEEDED"
	PipelineNotificationEvents_STAGE_EXECUTION_RESUMED PipelineNotificationEvents = "STAGE_EXECUTION_RESUMED"
	PipelineNotificationEvents_STAGE_EXECUTION_CANCELED PipelineNotificationEvents = "STAGE_EXECUTION_CANCELED"
	PipelineNotificationEvents_STAGE_EXECUTION_FAILED PipelineNotificationEvents = "STAGE_EXECUTION_FAILED"
	PipelineNotificationEvents_ACTION_EXECUTION_SUCCEEDED PipelineNotificationEvents = "ACTION_EXECUTION_SUCCEEDED"
	PipelineNotificationEvents_ACTION_EXECUTION_FAILED PipelineNotificationEvents = "ACTION_EXECUTION_FAILED"
	PipelineNotificationEvents_ACTION_EXECUTION_CANCELED PipelineNotificationEvents = "ACTION_EXECUTION_CANCELED"
	PipelineNotificationEvents_ACTION_EXECUTION_STARTED PipelineNotificationEvents = "ACTION_EXECUTION_STARTED"
	PipelineNotificationEvents_MANUAL_APPROVAL_FAILED PipelineNotificationEvents = "MANUAL_APPROVAL_FAILED"
	PipelineNotificationEvents_MANUAL_APPROVAL_NEEDED PipelineNotificationEvents = "MANUAL_APPROVAL_NEEDED"
	PipelineNotificationEvents_MANUAL_APPROVAL_SUCCEEDED PipelineNotificationEvents = "MANUAL_APPROVAL_SUCCEEDED"
)

// Additional options to pass to the notification rule.
// Experimental.
type PipelineNotifyOnOptions struct {
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
	// A list of event types associated with this notification rule for CodePipeline Pipeline.
	//
	// For a complete list of event types and IDs, see Notification concepts in the Developer Tools Console User Guide.
	// See: https://docs.aws.amazon.com/dtconsole/latest/userguide/concepts.html#concepts-api
	//
	// Experimental.
	Events *[]PipelineNotificationEvents `json:"events"`
}

// Experimental.
type PipelineProps struct {
	// The S3 bucket used by this Pipeline to store artifacts.
	// Experimental.
	ArtifactBucket awss3.IBucket `json:"artifactBucket"`
	// Create KMS keys for cross-account deployments.
	//
	// This controls whether the pipeline is enabled for cross-account deployments.
	//
	// By default cross-account deployments are enabled, but this feature requires
	// that KMS Customer Master Keys are created which have a cost of $1/month.
	//
	// If you do not need cross-account deployments, you can set this to `false` to
	// not create those keys and save on that cost (the artifact bucket will be
	// encrypted with an AWS-managed key). However, cross-account deployments will
	// no longer be possible.
	// Experimental.
	CrossAccountKeys *bool `json:"crossAccountKeys"`
	// A map of region to S3 bucket name used for cross-region CodePipeline.
	//
	// For every Action that you specify targeting a different region than the Pipeline itself,
	// if you don't provide an explicit Bucket for that region using this property,
	// the construct will automatically create a Stack containing an S3 Bucket in that region.
	// Experimental.
	CrossRegionReplicationBuckets *map[string]awss3.IBucket `json:"crossRegionReplicationBuckets"`
	// Name of the pipeline.
	// Experimental.
	PipelineName *string `json:"pipelineName"`
	// Indicates whether to rerun the AWS CodePipeline pipeline after you update it.
	// Experimental.
	RestartExecutionOnUpdate *bool `json:"restartExecutionOnUpdate"`
	// The IAM role to be assumed by this Pipeline.
	// Experimental.
	Role awsiam.IRole `json:"role"`
	// The list of Stages, in order, to create this Pipeline with.
	//
	// You can always add more Stages later by calling {@link Pipeline#addStage}.
	// Experimental.
	Stages *[]*StageProps `json:"stages"`
}

// Experimental.
type StageOptions struct {
	// The physical, human-readable name to assign to this Pipeline Stage.
	// Experimental.
	StageName *string `json:"stageName"`
	// The list of Actions to create this Stage with.
	//
	// You can always add more Actions later by calling {@link IStage#addAction}.
	// Experimental.
	Actions *[]IAction `json:"actions"`
	// Experimental.
	Placement *StagePlacement `json:"placement"`
}

// Allows you to control where to place a new Stage when it's added to the Pipeline.
//
// Note that you can provide only one of the below properties -
// specifying more than one will result in a validation error.
// See: #justAfter
//
// Experimental.
type StagePlacement struct {
	// Inserts the new Stage as a child of the given Stage (changing its current child Stage, if it had one).
	// Experimental.
	JustAfter IStage `json:"justAfter"`
	// Inserts the new Stage as a parent of the given Stage (changing its current parent Stage, if it had one).
	// Experimental.
	RightBefore IStage `json:"rightBefore"`
}

// Construction properties of a Pipeline Stage.
// Experimental.
type StageProps struct {
	// The physical, human-readable name to assign to this Pipeline Stage.
	// Experimental.
	StageName *string `json:"stageName"`
	// The list of Actions to create this Stage with.
	//
	// You can always add more Actions later by calling {@link IStage#addAction}.
	// Experimental.
	Actions *[]IAction `json:"actions"`
}


package awscodecommit

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscodecommit/internal"
	"github.com/aws/aws-cdk-go/awscdk/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
)

// A CloudFormation `AWS::CodeCommit::Repository`.
type CfnRepository interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	AttrCloneUrlHttp() *string
	AttrCloneUrlSsh() *string
	AttrName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Code() interface{}
	SetCode(val interface{})
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RepositoryDescription() *string
	SetRepositoryDescription(val *string)
	RepositoryName() *string
	SetRepositoryName(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	Triggers() interface{}
	SetTriggers(val interface{})
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

// The jsii proxy struct for CfnRepository
type jsiiProxy_CfnRepository struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRepository) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) AttrCloneUrlHttp() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCloneUrlHttp",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) AttrCloneUrlSsh() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCloneUrlSsh",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) AttrName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) Code() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"code",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) RepositoryDescription() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryDescription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) RepositoryName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) Triggers() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"triggers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRepository) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::CodeCommit::Repository`.
func NewCfnRepository(scope awscdk.Construct, id *string, props *CfnRepositoryProps) CfnRepository {
	_init_.Initialize()

	j := jsiiProxy_CfnRepository{}

	_jsii_.Create(
		"monocdk.aws_codecommit.CfnRepository",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::CodeCommit::Repository`.
func NewCfnRepository_Override(c CfnRepository, scope awscdk.Construct, id *string, props *CfnRepositoryProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codecommit.CfnRepository",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRepository) SetCode(val interface{}) {
	_jsii_.Set(
		j,
		"code",
		val,
	)
}

func (j *jsiiProxy_CfnRepository) SetRepositoryDescription(val *string) {
	_jsii_.Set(
		j,
		"repositoryDescription",
		val,
	)
}

func (j *jsiiProxy_CfnRepository) SetRepositoryName(val *string) {
	_jsii_.Set(
		j,
		"repositoryName",
		val,
	)
}

func (j *jsiiProxy_CfnRepository) SetTriggers(val interface{}) {
	_jsii_.Set(
		j,
		"triggers",
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
func CfnRepository_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codecommit.CfnRepository",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnRepository_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codecommit.CfnRepository",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnRepository_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codecommit.CfnRepository",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRepository_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.CfnRepository",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnRepository) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnRepository) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnRepository) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnRepository) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnRepository) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnRepository) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnRepository) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnRepository) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnRepository) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnRepository) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnRepository) OnPrepare() {
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
func (c *jsiiProxy_CfnRepository) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRepository) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnRepository) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnRepository) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnRepository) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnRepository) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnRepository) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRepository) ToString() *string {
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
func (c *jsiiProxy_CfnRepository) Validate() *[]*string {
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
func (c *jsiiProxy_CfnRepository) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnRepository_CodeProperty struct {
	// `CfnRepository.CodeProperty.S3`.
	S3 interface{} `json:"s3"`
	// `CfnRepository.CodeProperty.BranchName`.
	BranchName *string `json:"branchName"`
}

type CfnRepository_RepositoryTriggerProperty struct {
	// `CfnRepository.RepositoryTriggerProperty.DestinationArn`.
	DestinationArn *string `json:"destinationArn"`
	// `CfnRepository.RepositoryTriggerProperty.Events`.
	Events *[]*string `json:"events"`
	// `CfnRepository.RepositoryTriggerProperty.Name`.
	Name *string `json:"name"`
	// `CfnRepository.RepositoryTriggerProperty.Branches`.
	Branches *[]*string `json:"branches"`
	// `CfnRepository.RepositoryTriggerProperty.CustomData`.
	CustomData *string `json:"customData"`
}

type CfnRepository_S3Property struct {
	// `CfnRepository.S3Property.Bucket`.
	Bucket *string `json:"bucket"`
	// `CfnRepository.S3Property.Key`.
	Key *string `json:"key"`
	// `CfnRepository.S3Property.ObjectVersion`.
	ObjectVersion *string `json:"objectVersion"`
}

// Properties for defining a `AWS::CodeCommit::Repository`.
type CfnRepositoryProps struct {
	// `AWS::CodeCommit::Repository.RepositoryName`.
	RepositoryName *string `json:"repositoryName"`
	// `AWS::CodeCommit::Repository.Code`.
	Code interface{} `json:"code"`
	// `AWS::CodeCommit::Repository.RepositoryDescription`.
	RepositoryDescription *string `json:"repositoryDescription"`
	// `AWS::CodeCommit::Repository.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::CodeCommit::Repository.Triggers`.
	Triggers interface{} `json:"triggers"`
}

// Experimental.
type IRepository interface {
	awscdk.IResource
	// Grant the given principal identity permissions to perform the actions on this repository.
	// Experimental.
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	// Grant the given identity permissions to pull this repository.
	// Experimental.
	GrantPull(grantee awsiam.IGrantable) awsiam.Grant
	// Grant the given identity permissions to pull and push this repository.
	// Experimental.
	GrantPullPush(grantee awsiam.IGrantable) awsiam.Grant
	// Grant the given identity permissions to read this repository.
	// Experimental.
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	// Defines a CloudWatch event rule which triggers when a comment is made on a commit.
	// Experimental.
	OnCommentOnCommit(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers when a comment is made on a pull request.
	// Experimental.
	OnCommentOnPullRequest(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers when a commit is pushed to a branch.
	// Experimental.
	OnCommit(id *string, options *OnCommitOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers for repository events.
	//
	// Use
	// `rule.addEventPattern(pattern)` to specify a filter.
	// Experimental.
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers when a pull request state is changed.
	// Experimental.
	OnPullRequestStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers when a reference is created (i.e. a new branch/tag is created) to the repository.
	// Experimental.
	OnReferenceCreated(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers when a reference is delete (i.e. a branch/tag is deleted) from the repository.
	// Experimental.
	OnReferenceDeleted(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers when a reference is updated (i.e. a commit is pushed to an existing or new branch) from the repository.
	// Experimental.
	OnReferenceUpdated(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// Defines a CloudWatch event rule which triggers when a "CodeCommit Repository State Change" event occurs.
	// Experimental.
	OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	// The ARN of this Repository.
	// Experimental.
	RepositoryArn() *string
	// The HTTPS (GRC) clone URL.
	//
	// HTTPS (GRC) is the protocol to use with git-remote-codecommit (GRC).
	//
	// It is the recommended method for supporting connections made with federated
	// access, identity providers, and temporary credentials.
	// See: https://docs.aws.amazon.com/codecommit/latest/userguide/setting-up-git-remote-codecommit.html
	//
	// Experimental.
	RepositoryCloneUrlGrc() *string
	// The HTTP clone URL.
	// Experimental.
	RepositoryCloneUrlHttp() *string
	// The SSH clone URL.
	// Experimental.
	RepositoryCloneUrlSsh() *string
	// The human-visible name of this Repository.
	// Experimental.
	RepositoryName() *string
}

// The jsii proxy for IRepository
type jsiiProxy_IRepository struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IRepository) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grant",
		args,
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) GrantPull(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPull",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) GrantPullPush(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPullPush",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnCommentOnCommit(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onCommentOnCommit",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnCommentOnPullRequest(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onCommentOnPullRequest",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnCommit(id *string, options *OnCommitOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onCommit",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnPullRequestStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onPullRequestStateChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnReferenceCreated(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onReferenceCreated",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnReferenceDeleted(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onReferenceDeleted",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnReferenceUpdated(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onReferenceUpdated",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRepository) OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onStateChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IRepository) RepositoryArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRepository) RepositoryCloneUrlGrc() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryCloneUrlGrc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRepository) RepositoryCloneUrlHttp() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryCloneUrlHttp",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRepository) RepositoryCloneUrlSsh() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryCloneUrlSsh",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRepository) RepositoryName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryName",
		&returns,
	)
	return returns
}

// Options for the onCommit() method.
// Experimental.
type OnCommitOptions struct {
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
	EventPattern *awsevents.EventPattern `json:"eventPattern"`
	// A name for the rule.
	// Experimental.
	RuleName *string `json:"ruleName"`
	// The target to register for the event.
	// Experimental.
	Target awsevents.IRuleTarget `json:"target"`
	// The branch to monitor.
	// Experimental.
	Branches *[]*string `json:"branches"`
}

// Fields of CloudWatch Events that change references.
// See: https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/EventTypes.html#codebuild_event_type
//
// Experimental.
type ReferenceEvent interface {
}

// The jsii proxy struct for ReferenceEvent
type jsiiProxy_ReferenceEvent struct {
	_ byte // padding
}

func ReferenceEvent_CommitId() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.ReferenceEvent",
		"commitId",
		&returns,
	)
	return returns
}

func ReferenceEvent_EventType() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.ReferenceEvent",
		"eventType",
		&returns,
	)
	return returns
}

func ReferenceEvent_ReferenceFullName() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.ReferenceEvent",
		"referenceFullName",
		&returns,
	)
	return returns
}

func ReferenceEvent_ReferenceName() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.ReferenceEvent",
		"referenceName",
		&returns,
	)
	return returns
}

func ReferenceEvent_ReferenceType() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.ReferenceEvent",
		"referenceType",
		&returns,
	)
	return returns
}

func ReferenceEvent_RepositoryId() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.ReferenceEvent",
		"repositoryId",
		&returns,
	)
	return returns
}

func ReferenceEvent_RepositoryName() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_codecommit.ReferenceEvent",
		"repositoryName",
		&returns,
	)
	return returns
}

// Provides a CodeCommit Repository.
// Experimental.
type Repository interface {
	awscdk.Resource
	IRepository
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RepositoryArn() *string
	RepositoryCloneUrlGrc() *string
	RepositoryCloneUrlHttp() *string
	RepositoryCloneUrlSsh() *string
	RepositoryName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant
	GrantPull(grantee awsiam.IGrantable) awsiam.Grant
	GrantPullPush(grantee awsiam.IGrantable) awsiam.Grant
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	Notify(arn *string, options *RepositoryTriggerOptions) Repository
	OnCommentOnCommit(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnCommentOnPullRequest(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnCommit(id *string, options *OnCommitOptions) awsevents.Rule
	OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnPrepare()
	OnPullRequestStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnReferenceCreated(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnReferenceDeleted(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnReferenceUpdated(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Repository
type jsiiProxy_Repository struct {
	internal.Type__awscdkResource
	jsiiProxy_IRepository
}

func (j *jsiiProxy_Repository) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) RepositoryArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) RepositoryCloneUrlGrc() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryCloneUrlGrc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) RepositoryCloneUrlHttp() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryCloneUrlHttp",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) RepositoryCloneUrlSsh() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryCloneUrlSsh",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) RepositoryName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"repositoryName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Repository) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRepository(scope constructs.Construct, id *string, props *RepositoryProps) Repository {
	_init_.Initialize()

	j := jsiiProxy_Repository{}

	_jsii_.Create(
		"monocdk.aws_codecommit.Repository",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRepository_Override(r Repository, scope constructs.Construct, id *string, props *RepositoryProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_codecommit.Repository",
		[]interface{}{scope, id, props},
		r,
	)
}

// Imports a codecommit repository.
// Experimental.
func Repository_FromRepositoryArn(scope constructs.Construct, id *string, repositoryArn *string) IRepository {
	_init_.Initialize()

	var returns IRepository

	_jsii_.StaticInvoke(
		"monocdk.aws_codecommit.Repository",
		"fromRepositoryArn",
		[]interface{}{scope, id, repositoryArn},
		&returns,
	)

	return returns
}

// Experimental.
func Repository_FromRepositoryName(scope constructs.Construct, id *string, repositoryName *string) IRepository {
	_init_.Initialize()

	var returns IRepository

	_jsii_.StaticInvoke(
		"monocdk.aws_codecommit.Repository",
		"fromRepositoryName",
		[]interface{}{scope, id, repositoryName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Repository_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codecommit.Repository",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Repository_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_codecommit.Repository",
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
func (r *jsiiProxy_Repository) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_Repository) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_Repository) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_Repository) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the given principal identity permissions to perform the actions on this repository.
// Experimental.
func (r *jsiiProxy_Repository) Grant(grantee awsiam.IGrantable, actions ...*string) awsiam.Grant {
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Grant the given identity permissions to pull this repository.
// Experimental.
func (r *jsiiProxy_Repository) GrantPull(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grantPull",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to pull and push this repository.
// Experimental.
func (r *jsiiProxy_Repository) GrantPullPush(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grantPullPush",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Grant the given identity permissions to read this repository.
// Experimental.
func (r *jsiiProxy_Repository) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Create a trigger to notify another service to run actions on repository events.
// Experimental.
func (r *jsiiProxy_Repository) Notify(arn *string, options *RepositoryTriggerOptions) Repository {
	var returns Repository

	_jsii_.Invoke(
		r,
		"notify",
		[]interface{}{arn, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers when a comment is made on a commit.
// Experimental.
func (r *jsiiProxy_Repository) OnCommentOnCommit(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
		"onCommentOnCommit",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers when a comment is made on a pull request.
// Experimental.
func (r *jsiiProxy_Repository) OnCommentOnPullRequest(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
		"onCommentOnPullRequest",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers when a commit is pushed to a branch.
// Experimental.
func (r *jsiiProxy_Repository) OnCommit(id *string, options *OnCommitOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
		"onCommit",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers for repository events.
//
// Use
// `rule.addEventPattern(pattern)` to specify a filter.
// Experimental.
func (r *jsiiProxy_Repository) OnEvent(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_Repository) OnPrepare() {
	_jsii_.InvokeVoid(
		r,
		"onPrepare",
		nil, // no parameters
	)
}

// Defines a CloudWatch event rule which triggers when a pull request state is changed.
// Experimental.
func (r *jsiiProxy_Repository) OnPullRequestStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
		"onPullRequestStateChange",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers when a reference is created (i.e. a new branch/tag is created) to the repository.
// Experimental.
func (r *jsiiProxy_Repository) OnReferenceCreated(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
		"onReferenceCreated",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers when a reference is delete (i.e. a branch/tag is deleted) from the repository.
// Experimental.
func (r *jsiiProxy_Repository) OnReferenceDeleted(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
		"onReferenceDeleted",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers when a reference is updated (i.e. a commit is pushed to an existing or new branch) from the repository.
// Experimental.
func (r *jsiiProxy_Repository) OnReferenceUpdated(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
		"onReferenceUpdated",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines a CloudWatch event rule which triggers when a "CodeCommit Repository State Change" event occurs.
// Experimental.
func (r *jsiiProxy_Repository) OnStateChange(id *string, options *awsevents.OnEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_Repository) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_Repository) OnValidate() *[]*string {
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
func (r *jsiiProxy_Repository) Prepare() {
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
func (r *jsiiProxy_Repository) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_Repository) ToString() *string {
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
func (r *jsiiProxy_Repository) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Repository events that will cause the trigger to run actions in another service.
// Experimental.
type RepositoryEventTrigger string

const (
	RepositoryEventTrigger_ALL RepositoryEventTrigger = "ALL"
	RepositoryEventTrigger_UPDATE_REF RepositoryEventTrigger = "UPDATE_REF"
	RepositoryEventTrigger_CREATE_REF RepositoryEventTrigger = "CREATE_REF"
	RepositoryEventTrigger_DELETE_REF RepositoryEventTrigger = "DELETE_REF"
)

// Experimental.
type RepositoryProps struct {
	// Name of the repository.
	//
	// This property is required for all CodeCommit repositories.
	// Experimental.
	RepositoryName *string `json:"repositoryName"`
	// A description of the repository.
	//
	// Use the description to identify the
	// purpose of the repository.
	// Experimental.
	Description *string `json:"description"`
}

// Creates for a repository trigger to an SNS topic or Lambda function.
// Experimental.
type RepositoryTriggerOptions struct {
	// The names of the branches in the AWS CodeCommit repository that contain events that you want to include in the trigger.
	//
	// If you don't specify at
	// least one branch, the trigger applies to all branches.
	// Experimental.
	Branches *[]*string `json:"branches"`
	// When an event is triggered, additional information that AWS CodeCommit includes when it sends information to the target.
	// Experimental.
	CustomData *string `json:"customData"`
	// The repository events for which AWS CodeCommit sends information to the target, which you specified in the DestinationArn property.If you don't specify events, the trigger runs for all repository events.
	// Experimental.
	Events *[]RepositoryEventTrigger `json:"events"`
	// A name for the trigger.Triggers on a repository must have unique names.
	// Experimental.
	Name *string `json:"name"`
}


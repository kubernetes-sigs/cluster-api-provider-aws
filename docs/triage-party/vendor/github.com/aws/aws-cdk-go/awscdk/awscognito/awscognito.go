package awscognito

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/awscognito/internal"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/constructs-go/constructs/v3"
)

// How will a user be able to recover their account?
//
// When a user forgets their password, they can have a code sent to their verified email or verified phone to recover their account.
// You can choose the preferred way to send codes below.
// We recommend not allowing phone to be used for both password resets and multi-factor authentication (MFA).
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/how-to-recover-a-user-account.html
//
// Experimental.
type AccountRecovery string

const (
	AccountRecovery_EMAIL_AND_PHONE_WITHOUT_MFA AccountRecovery = "EMAIL_AND_PHONE_WITHOUT_MFA"
	AccountRecovery_PHONE_WITHOUT_MFA_AND_EMAIL AccountRecovery = "PHONE_WITHOUT_MFA_AND_EMAIL"
	AccountRecovery_EMAIL_ONLY AccountRecovery = "EMAIL_ONLY"
	AccountRecovery_PHONE_ONLY_WITHOUT_MFA AccountRecovery = "PHONE_ONLY_WITHOUT_MFA"
	AccountRecovery_PHONE_AND_EMAIL AccountRecovery = "PHONE_AND_EMAIL"
	AccountRecovery_NONE AccountRecovery = "NONE"
)

// The mapping of user pool attributes to the attributes provided by the identity providers.
// Experimental.
type AttributeMapping struct {
	// The user's postal address is a required attribute.
	// Experimental.
	Address ProviderAttribute `json:"address"`
	// The user's birthday.
	// Experimental.
	Birthdate ProviderAttribute `json:"birthdate"`
	// Specify custom attribute mapping here and mapping for any standard attributes not supported yet.
	// Experimental.
	Custom *map[string]ProviderAttribute `json:"custom"`
	// The user's e-mail address.
	// Experimental.
	Email ProviderAttribute `json:"email"`
	// The surname or last name of user.
	// Experimental.
	FamilyName ProviderAttribute `json:"familyName"`
	// The user's full name in displayable form.
	// Experimental.
	Fullname ProviderAttribute `json:"fullname"`
	// The user's gender.
	// Experimental.
	Gender ProviderAttribute `json:"gender"`
	// The user's first name or give name.
	// Experimental.
	GivenName ProviderAttribute `json:"givenName"`
	// Time, the user's information was last updated.
	// Experimental.
	LastUpdateTime ProviderAttribute `json:"lastUpdateTime"`
	// The user's locale.
	// Experimental.
	Locale ProviderAttribute `json:"locale"`
	// The user's middle name.
	// Experimental.
	MiddleName ProviderAttribute `json:"middleName"`
	// The user's nickname or casual name.
	// Experimental.
	Nickname ProviderAttribute `json:"nickname"`
	// The user's telephone number.
	// Experimental.
	PhoneNumber ProviderAttribute `json:"phoneNumber"`
	// The user's preferred username.
	// Experimental.
	PreferredUsername ProviderAttribute `json:"preferredUsername"`
	// The URL to the user's profile page.
	// Experimental.
	ProfilePage ProviderAttribute `json:"profilePage"`
	// The URL to the user's profile picture.
	// Experimental.
	ProfilePicture ProviderAttribute `json:"profilePicture"`
	// The user's time zone.
	// Experimental.
	Timezone ProviderAttribute `json:"timezone"`
	// The URL to the user's web page or blog.
	// Experimental.
	Website ProviderAttribute `json:"website"`
}

// Types of authentication flow.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow.html
//
// Experimental.
type AuthFlow struct {
	// Enable admin based user password authentication flow.
	// Experimental.
	AdminUserPassword *bool `json:"adminUserPassword"`
	// Enable custom authentication flow.
	// Experimental.
	Custom *bool `json:"custom"`
	// Enable auth using username & password.
	// Experimental.
	UserPassword *bool `json:"userPassword"`
	// Enable SRP based authentication.
	// Experimental.
	UserSrp *bool `json:"userSrp"`
}

// Attributes that can be automatically verified for users in a user pool.
// Experimental.
type AutoVerifiedAttrs struct {
	// Whether the email address of the user should be auto verified at sign up.
	//
	// Note: If both `email` and `phone` is set, Cognito only verifies the phone number. To also verify email, see here -
	// https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-email-phone-verification.html
	// Experimental.
	Email *bool `json:"email"`
	// Whether the phone number of the user should be auto verified at sign up.
	// Experimental.
	Phone *bool `json:"phone"`
}

// The Boolean custom attribute type.
// Experimental.
type BooleanAttribute interface {
	ICustomAttribute
	Bind() *CustomAttributeConfig
}

// The jsii proxy struct for BooleanAttribute
type jsiiProxy_BooleanAttribute struct {
	jsiiProxy_ICustomAttribute
}

// Experimental.
func NewBooleanAttribute(props *CustomAttributeProps) BooleanAttribute {
	_init_.Initialize()

	j := jsiiProxy_BooleanAttribute{}

	_jsii_.Create(
		"monocdk.aws_cognito.BooleanAttribute",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewBooleanAttribute_Override(b BooleanAttribute, props *CustomAttributeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.BooleanAttribute",
		[]interface{}{props},
		b,
	)
}

// Bind this custom attribute type to the values as expected by CloudFormation.
// Experimental.
func (b *jsiiProxy_BooleanAttribute) Bind() *CustomAttributeConfig {
	var returns *CustomAttributeConfig

	_jsii_.Invoke(
		b,
		"bind",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A CloudFormation `AWS::Cognito::IdentityPool`.
type CfnIdentityPool interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AllowClassicFlow() interface{}
	SetAllowClassicFlow(val interface{})
	AllowUnauthenticatedIdentities() interface{}
	SetAllowUnauthenticatedIdentities(val interface{})
	AttrName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CognitoEvents() interface{}
	SetCognitoEvents(val interface{})
	CognitoIdentityProviders() interface{}
	SetCognitoIdentityProviders(val interface{})
	CognitoStreams() interface{}
	SetCognitoStreams(val interface{})
	CreationStack() *[]*string
	DeveloperProviderName() *string
	SetDeveloperProviderName(val *string)
	IdentityPoolName() *string
	SetIdentityPoolName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	OpenIdConnectProviderArns() *[]*string
	SetOpenIdConnectProviderArns(val *[]*string)
	PushSync() interface{}
	SetPushSync(val interface{})
	Ref() *string
	SamlProviderArns() *[]*string
	SetSamlProviderArns(val *[]*string)
	Stack() awscdk.Stack
	SupportedLoginProviders() interface{}
	SetSupportedLoginProviders(val interface{})
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

// The jsii proxy struct for CfnIdentityPool
type jsiiProxy_CfnIdentityPool struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnIdentityPool) AllowClassicFlow() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"allowClassicFlow",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) AllowUnauthenticatedIdentities() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"allowUnauthenticatedIdentities",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) AttrName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) CognitoEvents() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"cognitoEvents",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) CognitoIdentityProviders() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"cognitoIdentityProviders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) CognitoStreams() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"cognitoStreams",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) DeveloperProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"developerProviderName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) IdentityPoolName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identityPoolName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) OpenIdConnectProviderArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"openIdConnectProviderArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) PushSync() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"pushSync",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) SamlProviderArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"samlProviderArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) SupportedLoginProviders() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"supportedLoginProviders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPool) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::IdentityPool`.
func NewCfnIdentityPool(scope awscdk.Construct, id *string, props *CfnIdentityPoolProps) CfnIdentityPool {
	_init_.Initialize()

	j := jsiiProxy_CfnIdentityPool{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnIdentityPool",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::IdentityPool`.
func NewCfnIdentityPool_Override(c CfnIdentityPool, scope awscdk.Construct, id *string, props *CfnIdentityPoolProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnIdentityPool",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetAllowClassicFlow(val interface{}) {
	_jsii_.Set(
		j,
		"allowClassicFlow",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetAllowUnauthenticatedIdentities(val interface{}) {
	_jsii_.Set(
		j,
		"allowUnauthenticatedIdentities",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetCognitoEvents(val interface{}) {
	_jsii_.Set(
		j,
		"cognitoEvents",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetCognitoIdentityProviders(val interface{}) {
	_jsii_.Set(
		j,
		"cognitoIdentityProviders",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetCognitoStreams(val interface{}) {
	_jsii_.Set(
		j,
		"cognitoStreams",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetDeveloperProviderName(val *string) {
	_jsii_.Set(
		j,
		"developerProviderName",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetIdentityPoolName(val *string) {
	_jsii_.Set(
		j,
		"identityPoolName",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetOpenIdConnectProviderArns(val *[]*string) {
	_jsii_.Set(
		j,
		"openIdConnectProviderArns",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetPushSync(val interface{}) {
	_jsii_.Set(
		j,
		"pushSync",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetSamlProviderArns(val *[]*string) {
	_jsii_.Set(
		j,
		"samlProviderArns",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPool) SetSupportedLoginProviders(val interface{}) {
	_jsii_.Set(
		j,
		"supportedLoginProviders",
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
func CfnIdentityPool_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnIdentityPool",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnIdentityPool_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnIdentityPool",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnIdentityPool_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnIdentityPool",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnIdentityPool_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnIdentityPool",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnIdentityPool) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnIdentityPool) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnIdentityPool) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnIdentityPool) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnIdentityPool) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnIdentityPool) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnIdentityPool) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnIdentityPool) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnIdentityPool) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnIdentityPool) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnIdentityPool) OnPrepare() {
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
func (c *jsiiProxy_CfnIdentityPool) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnIdentityPool) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnIdentityPool) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnIdentityPool) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnIdentityPool) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnIdentityPool) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnIdentityPool) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnIdentityPool) ToString() *string {
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
func (c *jsiiProxy_CfnIdentityPool) Validate() *[]*string {
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
func (c *jsiiProxy_CfnIdentityPool) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnIdentityPool_CognitoIdentityProviderProperty struct {
	// `CfnIdentityPool.CognitoIdentityProviderProperty.ClientId`.
	ClientId *string `json:"clientId"`
	// `CfnIdentityPool.CognitoIdentityProviderProperty.ProviderName`.
	ProviderName *string `json:"providerName"`
	// `CfnIdentityPool.CognitoIdentityProviderProperty.ServerSideTokenCheck`.
	ServerSideTokenCheck interface{} `json:"serverSideTokenCheck"`
}

type CfnIdentityPool_CognitoStreamsProperty struct {
	// `CfnIdentityPool.CognitoStreamsProperty.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `CfnIdentityPool.CognitoStreamsProperty.StreamingStatus`.
	StreamingStatus *string `json:"streamingStatus"`
	// `CfnIdentityPool.CognitoStreamsProperty.StreamName`.
	StreamName *string `json:"streamName"`
}

type CfnIdentityPool_PushSyncProperty struct {
	// `CfnIdentityPool.PushSyncProperty.ApplicationArns`.
	ApplicationArns *[]*string `json:"applicationArns"`
	// `CfnIdentityPool.PushSyncProperty.RoleArn`.
	RoleArn *string `json:"roleArn"`
}

// Properties for defining a `AWS::Cognito::IdentityPool`.
type CfnIdentityPoolProps struct {
	// `AWS::Cognito::IdentityPool.AllowUnauthenticatedIdentities`.
	AllowUnauthenticatedIdentities interface{} `json:"allowUnauthenticatedIdentities"`
	// `AWS::Cognito::IdentityPool.AllowClassicFlow`.
	AllowClassicFlow interface{} `json:"allowClassicFlow"`
	// `AWS::Cognito::IdentityPool.CognitoEvents`.
	CognitoEvents interface{} `json:"cognitoEvents"`
	// `AWS::Cognito::IdentityPool.CognitoIdentityProviders`.
	CognitoIdentityProviders interface{} `json:"cognitoIdentityProviders"`
	// `AWS::Cognito::IdentityPool.CognitoStreams`.
	CognitoStreams interface{} `json:"cognitoStreams"`
	// `AWS::Cognito::IdentityPool.DeveloperProviderName`.
	DeveloperProviderName *string `json:"developerProviderName"`
	// `AWS::Cognito::IdentityPool.IdentityPoolName`.
	IdentityPoolName *string `json:"identityPoolName"`
	// `AWS::Cognito::IdentityPool.OpenIdConnectProviderARNs`.
	OpenIdConnectProviderArns *[]*string `json:"openIdConnectProviderArns"`
	// `AWS::Cognito::IdentityPool.PushSync`.
	PushSync interface{} `json:"pushSync"`
	// `AWS::Cognito::IdentityPool.SamlProviderARNs`.
	SamlProviderArns *[]*string `json:"samlProviderArns"`
	// `AWS::Cognito::IdentityPool.SupportedLoginProviders`.
	SupportedLoginProviders interface{} `json:"supportedLoginProviders"`
}

// A CloudFormation `AWS::Cognito::IdentityPoolRoleAttachment`.
type CfnIdentityPoolRoleAttachment interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	IdentityPoolId() *string
	SetIdentityPoolId(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RoleMappings() interface{}
	SetRoleMappings(val interface{})
	Roles() interface{}
	SetRoles(val interface{})
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

// The jsii proxy struct for CfnIdentityPoolRoleAttachment
type jsiiProxy_CfnIdentityPoolRoleAttachment struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) IdentityPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identityPoolId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) RoleMappings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"roleMappings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) Roles() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"roles",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::IdentityPoolRoleAttachment`.
func NewCfnIdentityPoolRoleAttachment(scope awscdk.Construct, id *string, props *CfnIdentityPoolRoleAttachmentProps) CfnIdentityPoolRoleAttachment {
	_init_.Initialize()

	j := jsiiProxy_CfnIdentityPoolRoleAttachment{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::IdentityPoolRoleAttachment`.
func NewCfnIdentityPoolRoleAttachment_Override(c CfnIdentityPoolRoleAttachment, scope awscdk.Construct, id *string, props *CfnIdentityPoolRoleAttachmentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) SetIdentityPoolId(val *string) {
	_jsii_.Set(
		j,
		"identityPoolId",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) SetRoleMappings(val interface{}) {
	_jsii_.Set(
		j,
		"roleMappings",
		val,
	)
}

func (j *jsiiProxy_CfnIdentityPoolRoleAttachment) SetRoles(val interface{}) {
	_jsii_.Set(
		j,
		"roles",
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
func CfnIdentityPoolRoleAttachment_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnIdentityPoolRoleAttachment_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnIdentityPoolRoleAttachment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnIdentityPoolRoleAttachment_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnIdentityPoolRoleAttachment",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) OnPrepare() {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) ToString() *string {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) Validate() *[]*string {
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
func (c *jsiiProxy_CfnIdentityPoolRoleAttachment) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnIdentityPoolRoleAttachment_MappingRuleProperty struct {
	// `CfnIdentityPoolRoleAttachment.MappingRuleProperty.Claim`.
	Claim *string `json:"claim"`
	// `CfnIdentityPoolRoleAttachment.MappingRuleProperty.MatchType`.
	MatchType *string `json:"matchType"`
	// `CfnIdentityPoolRoleAttachment.MappingRuleProperty.RoleARN`.
	RoleArn *string `json:"roleArn"`
	// `CfnIdentityPoolRoleAttachment.MappingRuleProperty.Value`.
	Value *string `json:"value"`
}

type CfnIdentityPoolRoleAttachment_RoleMappingProperty struct {
	// `CfnIdentityPoolRoleAttachment.RoleMappingProperty.Type`.
	Type *string `json:"type"`
	// `CfnIdentityPoolRoleAttachment.RoleMappingProperty.AmbiguousRoleResolution`.
	AmbiguousRoleResolution *string `json:"ambiguousRoleResolution"`
	// `CfnIdentityPoolRoleAttachment.RoleMappingProperty.IdentityProvider`.
	IdentityProvider *string `json:"identityProvider"`
	// `CfnIdentityPoolRoleAttachment.RoleMappingProperty.RulesConfiguration`.
	RulesConfiguration interface{} `json:"rulesConfiguration"`
}

type CfnIdentityPoolRoleAttachment_RulesConfigurationTypeProperty struct {
	// `CfnIdentityPoolRoleAttachment.RulesConfigurationTypeProperty.Rules`.
	Rules interface{} `json:"rules"`
}

// Properties for defining a `AWS::Cognito::IdentityPoolRoleAttachment`.
type CfnIdentityPoolRoleAttachmentProps struct {
	// `AWS::Cognito::IdentityPoolRoleAttachment.IdentityPoolId`.
	IdentityPoolId *string `json:"identityPoolId"`
	// `AWS::Cognito::IdentityPoolRoleAttachment.RoleMappings`.
	RoleMappings interface{} `json:"roleMappings"`
	// `AWS::Cognito::IdentityPoolRoleAttachment.Roles`.
	Roles interface{} `json:"roles"`
}

// A CloudFormation `AWS::Cognito::UserPool`.
type CfnUserPool interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AccountRecoverySetting() interface{}
	SetAccountRecoverySetting(val interface{})
	AdminCreateUserConfig() interface{}
	SetAdminCreateUserConfig(val interface{})
	AliasAttributes() *[]*string
	SetAliasAttributes(val *[]*string)
	AttrArn() *string
	AttrProviderName() *string
	AttrProviderUrl() *string
	AutoVerifiedAttributes() *[]*string
	SetAutoVerifiedAttributes(val *[]*string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DeviceConfiguration() interface{}
	SetDeviceConfiguration(val interface{})
	EmailConfiguration() interface{}
	SetEmailConfiguration(val interface{})
	EmailVerificationMessage() *string
	SetEmailVerificationMessage(val *string)
	EmailVerificationSubject() *string
	SetEmailVerificationSubject(val *string)
	EnabledMfas() *[]*string
	SetEnabledMfas(val *[]*string)
	LambdaConfig() interface{}
	SetLambdaConfig(val interface{})
	LogicalId() *string
	MfaConfiguration() *string
	SetMfaConfiguration(val *string)
	Node() awscdk.ConstructNode
	Policies() interface{}
	SetPolicies(val interface{})
	Ref() *string
	Schema() interface{}
	SetSchema(val interface{})
	SmsAuthenticationMessage() *string
	SetSmsAuthenticationMessage(val *string)
	SmsConfiguration() interface{}
	SetSmsConfiguration(val interface{})
	SmsVerificationMessage() *string
	SetSmsVerificationMessage(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	UpdatedProperites() *map[string]interface{}
	UsernameAttributes() *[]*string
	SetUsernameAttributes(val *[]*string)
	UsernameConfiguration() interface{}
	SetUsernameConfiguration(val interface{})
	UserPoolAddOns() interface{}
	SetUserPoolAddOns(val interface{})
	UserPoolName() *string
	SetUserPoolName(val *string)
	VerificationMessageTemplate() interface{}
	SetVerificationMessageTemplate(val interface{})
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

// The jsii proxy struct for CfnUserPool
type jsiiProxy_CfnUserPool struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPool) AccountRecoverySetting() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accountRecoverySetting",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) AdminCreateUserConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"adminCreateUserConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) AliasAttributes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"aliasAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) AttrProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrProviderName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) AttrProviderUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrProviderUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) AutoVerifiedAttributes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"autoVerifiedAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) DeviceConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deviceConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) EmailConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"emailConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) EmailVerificationMessage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"emailVerificationMessage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) EmailVerificationSubject() *string {
	var returns *string
	_jsii_.Get(
		j,
		"emailVerificationSubject",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) EnabledMfas() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"enabledMfas",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) LambdaConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"lambdaConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) MfaConfiguration() *string {
	var returns *string
	_jsii_.Get(
		j,
		"mfaConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) Policies() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) Schema() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"schema",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) SmsAuthenticationMessage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"smsAuthenticationMessage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) SmsConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"smsConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) SmsVerificationMessage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"smsVerificationMessage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) UsernameAttributes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"usernameAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) UsernameConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"usernameConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) UserPoolAddOns() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"userPoolAddOns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) UserPoolName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPool) VerificationMessageTemplate() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"verificationMessageTemplate",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPool`.
func NewCfnUserPool(scope awscdk.Construct, id *string, props *CfnUserPoolProps) CfnUserPool {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPool{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPool",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPool`.
func NewCfnUserPool_Override(c CfnUserPool, scope awscdk.Construct, id *string, props *CfnUserPoolProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPool",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPool) SetAccountRecoverySetting(val interface{}) {
	_jsii_.Set(
		j,
		"accountRecoverySetting",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetAdminCreateUserConfig(val interface{}) {
	_jsii_.Set(
		j,
		"adminCreateUserConfig",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetAliasAttributes(val *[]*string) {
	_jsii_.Set(
		j,
		"aliasAttributes",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetAutoVerifiedAttributes(val *[]*string) {
	_jsii_.Set(
		j,
		"autoVerifiedAttributes",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetDeviceConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"deviceConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetEmailConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"emailConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetEmailVerificationMessage(val *string) {
	_jsii_.Set(
		j,
		"emailVerificationMessage",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetEmailVerificationSubject(val *string) {
	_jsii_.Set(
		j,
		"emailVerificationSubject",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetEnabledMfas(val *[]*string) {
	_jsii_.Set(
		j,
		"enabledMfas",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetLambdaConfig(val interface{}) {
	_jsii_.Set(
		j,
		"lambdaConfig",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetMfaConfiguration(val *string) {
	_jsii_.Set(
		j,
		"mfaConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetPolicies(val interface{}) {
	_jsii_.Set(
		j,
		"policies",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetSchema(val interface{}) {
	_jsii_.Set(
		j,
		"schema",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetSmsAuthenticationMessage(val *string) {
	_jsii_.Set(
		j,
		"smsAuthenticationMessage",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetSmsConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"smsConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetSmsVerificationMessage(val *string) {
	_jsii_.Set(
		j,
		"smsVerificationMessage",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetUsernameAttributes(val *[]*string) {
	_jsii_.Set(
		j,
		"usernameAttributes",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetUsernameConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"usernameConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetUserPoolAddOns(val interface{}) {
	_jsii_.Set(
		j,
		"userPoolAddOns",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetUserPoolName(val *string) {
	_jsii_.Set(
		j,
		"userPoolName",
		val,
	)
}

func (j *jsiiProxy_CfnUserPool) SetVerificationMessageTemplate(val interface{}) {
	_jsii_.Set(
		j,
		"verificationMessageTemplate",
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
func CfnUserPool_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPool",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPool_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPool",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPool_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPool",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPool_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPool",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPool) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPool) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPool) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPool) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPool) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPool) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPool) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPool) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPool) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPool) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPool) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPool) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPool) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPool) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPool) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPool) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPool) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPool) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPool) ToString() *string {
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
func (c *jsiiProxy_CfnUserPool) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPool) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUserPool_AccountRecoverySettingProperty struct {
	// `CfnUserPool.AccountRecoverySettingProperty.RecoveryMechanisms`.
	RecoveryMechanisms interface{} `json:"recoveryMechanisms"`
}

type CfnUserPool_AdminCreateUserConfigProperty struct {
	// `CfnUserPool.AdminCreateUserConfigProperty.AllowAdminCreateUserOnly`.
	AllowAdminCreateUserOnly interface{} `json:"allowAdminCreateUserOnly"`
	// `CfnUserPool.AdminCreateUserConfigProperty.InviteMessageTemplate`.
	InviteMessageTemplate interface{} `json:"inviteMessageTemplate"`
	// `CfnUserPool.AdminCreateUserConfigProperty.UnusedAccountValidityDays`.
	UnusedAccountValidityDays *float64 `json:"unusedAccountValidityDays"`
}

type CfnUserPool_CustomEmailSenderProperty struct {
	// `CfnUserPool.CustomEmailSenderProperty.LambdaArn`.
	LambdaArn *string `json:"lambdaArn"`
	// `CfnUserPool.CustomEmailSenderProperty.LambdaVersion`.
	LambdaVersion *string `json:"lambdaVersion"`
}

type CfnUserPool_CustomSMSSenderProperty struct {
	// `CfnUserPool.CustomSMSSenderProperty.LambdaArn`.
	LambdaArn *string `json:"lambdaArn"`
	// `CfnUserPool.CustomSMSSenderProperty.LambdaVersion`.
	LambdaVersion *string `json:"lambdaVersion"`
}

type CfnUserPool_DeviceConfigurationProperty struct {
	// `CfnUserPool.DeviceConfigurationProperty.ChallengeRequiredOnNewDevice`.
	ChallengeRequiredOnNewDevice interface{} `json:"challengeRequiredOnNewDevice"`
	// `CfnUserPool.DeviceConfigurationProperty.DeviceOnlyRememberedOnUserPrompt`.
	DeviceOnlyRememberedOnUserPrompt interface{} `json:"deviceOnlyRememberedOnUserPrompt"`
}

type CfnUserPool_EmailConfigurationProperty struct {
	// `CfnUserPool.EmailConfigurationProperty.ConfigurationSet`.
	ConfigurationSet *string `json:"configurationSet"`
	// `CfnUserPool.EmailConfigurationProperty.EmailSendingAccount`.
	EmailSendingAccount *string `json:"emailSendingAccount"`
	// `CfnUserPool.EmailConfigurationProperty.From`.
	From *string `json:"from"`
	// `CfnUserPool.EmailConfigurationProperty.ReplyToEmailAddress`.
	ReplyToEmailAddress *string `json:"replyToEmailAddress"`
	// `CfnUserPool.EmailConfigurationProperty.SourceArn`.
	SourceArn *string `json:"sourceArn"`
}

type CfnUserPool_InviteMessageTemplateProperty struct {
	// `CfnUserPool.InviteMessageTemplateProperty.EmailMessage`.
	EmailMessage *string `json:"emailMessage"`
	// `CfnUserPool.InviteMessageTemplateProperty.EmailSubject`.
	EmailSubject *string `json:"emailSubject"`
	// `CfnUserPool.InviteMessageTemplateProperty.SMSMessage`.
	SmsMessage *string `json:"smsMessage"`
}

type CfnUserPool_LambdaConfigProperty struct {
	// `CfnUserPool.LambdaConfigProperty.CreateAuthChallenge`.
	CreateAuthChallenge *string `json:"createAuthChallenge"`
	// `CfnUserPool.LambdaConfigProperty.CustomEmailSender`.
	CustomEmailSender interface{} `json:"customEmailSender"`
	// `CfnUserPool.LambdaConfigProperty.CustomMessage`.
	CustomMessage *string `json:"customMessage"`
	// `CfnUserPool.LambdaConfigProperty.CustomSMSSender`.
	CustomSmsSender interface{} `json:"customSmsSender"`
	// `CfnUserPool.LambdaConfigProperty.DefineAuthChallenge`.
	DefineAuthChallenge *string `json:"defineAuthChallenge"`
	// `CfnUserPool.LambdaConfigProperty.KMSKeyID`.
	KmsKeyId *string `json:"kmsKeyId"`
	// `CfnUserPool.LambdaConfigProperty.PostAuthentication`.
	PostAuthentication *string `json:"postAuthentication"`
	// `CfnUserPool.LambdaConfigProperty.PostConfirmation`.
	PostConfirmation *string `json:"postConfirmation"`
	// `CfnUserPool.LambdaConfigProperty.PreAuthentication`.
	PreAuthentication *string `json:"preAuthentication"`
	// `CfnUserPool.LambdaConfigProperty.PreSignUp`.
	PreSignUp *string `json:"preSignUp"`
	// `CfnUserPool.LambdaConfigProperty.PreTokenGeneration`.
	PreTokenGeneration *string `json:"preTokenGeneration"`
	// `CfnUserPool.LambdaConfigProperty.UserMigration`.
	UserMigration *string `json:"userMigration"`
	// `CfnUserPool.LambdaConfigProperty.VerifyAuthChallengeResponse`.
	VerifyAuthChallengeResponse *string `json:"verifyAuthChallengeResponse"`
}

type CfnUserPool_NumberAttributeConstraintsProperty struct {
	// `CfnUserPool.NumberAttributeConstraintsProperty.MaxValue`.
	MaxValue *string `json:"maxValue"`
	// `CfnUserPool.NumberAttributeConstraintsProperty.MinValue`.
	MinValue *string `json:"minValue"`
}

type CfnUserPool_PasswordPolicyProperty struct {
	// `CfnUserPool.PasswordPolicyProperty.MinimumLength`.
	MinimumLength *float64 `json:"minimumLength"`
	// `CfnUserPool.PasswordPolicyProperty.RequireLowercase`.
	RequireLowercase interface{} `json:"requireLowercase"`
	// `CfnUserPool.PasswordPolicyProperty.RequireNumbers`.
	RequireNumbers interface{} `json:"requireNumbers"`
	// `CfnUserPool.PasswordPolicyProperty.RequireSymbols`.
	RequireSymbols interface{} `json:"requireSymbols"`
	// `CfnUserPool.PasswordPolicyProperty.RequireUppercase`.
	RequireUppercase interface{} `json:"requireUppercase"`
	// `CfnUserPool.PasswordPolicyProperty.TemporaryPasswordValidityDays`.
	TemporaryPasswordValidityDays *float64 `json:"temporaryPasswordValidityDays"`
}

type CfnUserPool_PoliciesProperty struct {
	// `CfnUserPool.PoliciesProperty.PasswordPolicy`.
	PasswordPolicy interface{} `json:"passwordPolicy"`
}

type CfnUserPool_RecoveryOptionProperty struct {
	// `CfnUserPool.RecoveryOptionProperty.Name`.
	Name *string `json:"name"`
	// `CfnUserPool.RecoveryOptionProperty.Priority`.
	Priority *float64 `json:"priority"`
}

type CfnUserPool_SchemaAttributeProperty struct {
	// `CfnUserPool.SchemaAttributeProperty.AttributeDataType`.
	AttributeDataType *string `json:"attributeDataType"`
	// `CfnUserPool.SchemaAttributeProperty.DeveloperOnlyAttribute`.
	DeveloperOnlyAttribute interface{} `json:"developerOnlyAttribute"`
	// `CfnUserPool.SchemaAttributeProperty.Mutable`.
	Mutable interface{} `json:"mutable"`
	// `CfnUserPool.SchemaAttributeProperty.Name`.
	Name *string `json:"name"`
	// `CfnUserPool.SchemaAttributeProperty.NumberAttributeConstraints`.
	NumberAttributeConstraints interface{} `json:"numberAttributeConstraints"`
	// `CfnUserPool.SchemaAttributeProperty.Required`.
	Required interface{} `json:"required"`
	// `CfnUserPool.SchemaAttributeProperty.StringAttributeConstraints`.
	StringAttributeConstraints interface{} `json:"stringAttributeConstraints"`
}

type CfnUserPool_SmsConfigurationProperty struct {
	// `CfnUserPool.SmsConfigurationProperty.ExternalId`.
	ExternalId *string `json:"externalId"`
	// `CfnUserPool.SmsConfigurationProperty.SnsCallerArn`.
	SnsCallerArn *string `json:"snsCallerArn"`
}

type CfnUserPool_StringAttributeConstraintsProperty struct {
	// `CfnUserPool.StringAttributeConstraintsProperty.MaxLength`.
	MaxLength *string `json:"maxLength"`
	// `CfnUserPool.StringAttributeConstraintsProperty.MinLength`.
	MinLength *string `json:"minLength"`
}

type CfnUserPool_UserPoolAddOnsProperty struct {
	// `CfnUserPool.UserPoolAddOnsProperty.AdvancedSecurityMode`.
	AdvancedSecurityMode *string `json:"advancedSecurityMode"`
}

type CfnUserPool_UsernameConfigurationProperty struct {
	// `CfnUserPool.UsernameConfigurationProperty.CaseSensitive`.
	CaseSensitive interface{} `json:"caseSensitive"`
}

type CfnUserPool_VerificationMessageTemplateProperty struct {
	// `CfnUserPool.VerificationMessageTemplateProperty.DefaultEmailOption`.
	DefaultEmailOption *string `json:"defaultEmailOption"`
	// `CfnUserPool.VerificationMessageTemplateProperty.EmailMessage`.
	EmailMessage *string `json:"emailMessage"`
	// `CfnUserPool.VerificationMessageTemplateProperty.EmailMessageByLink`.
	EmailMessageByLink *string `json:"emailMessageByLink"`
	// `CfnUserPool.VerificationMessageTemplateProperty.EmailSubject`.
	EmailSubject *string `json:"emailSubject"`
	// `CfnUserPool.VerificationMessageTemplateProperty.EmailSubjectByLink`.
	EmailSubjectByLink *string `json:"emailSubjectByLink"`
	// `CfnUserPool.VerificationMessageTemplateProperty.SmsMessage`.
	SmsMessage *string `json:"smsMessage"`
}

// A CloudFormation `AWS::Cognito::UserPoolClient`.
type CfnUserPoolClient interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AccessTokenValidity() *float64
	SetAccessTokenValidity(val *float64)
	AllowedOAuthFlows() *[]*string
	SetAllowedOAuthFlows(val *[]*string)
	AllowedOAuthFlowsUserPoolClient() interface{}
	SetAllowedOAuthFlowsUserPoolClient(val interface{})
	AllowedOAuthScopes() *[]*string
	SetAllowedOAuthScopes(val *[]*string)
	AnalyticsConfiguration() interface{}
	SetAnalyticsConfiguration(val interface{})
	AttrClientSecret() *string
	AttrName() *string
	CallbackUrLs() *[]*string
	SetCallbackUrLs(val *[]*string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClientName() *string
	SetClientName(val *string)
	CreationStack() *[]*string
	DefaultRedirectUri() *string
	SetDefaultRedirectUri(val *string)
	ExplicitAuthFlows() *[]*string
	SetExplicitAuthFlows(val *[]*string)
	GenerateSecret() interface{}
	SetGenerateSecret(val interface{})
	IdTokenValidity() *float64
	SetIdTokenValidity(val *float64)
	LogicalId() *string
	LogoutUrLs() *[]*string
	SetLogoutUrLs(val *[]*string)
	Node() awscdk.ConstructNode
	PreventUserExistenceErrors() *string
	SetPreventUserExistenceErrors(val *string)
	ReadAttributes() *[]*string
	SetReadAttributes(val *[]*string)
	Ref() *string
	RefreshTokenValidity() *float64
	SetRefreshTokenValidity(val *float64)
	Stack() awscdk.Stack
	SupportedIdentityProviders() *[]*string
	SetSupportedIdentityProviders(val *[]*string)
	TokenValidityUnits() interface{}
	SetTokenValidityUnits(val interface{})
	UpdatedProperites() *map[string]interface{}
	UserPoolId() *string
	SetUserPoolId(val *string)
	WriteAttributes() *[]*string
	SetWriteAttributes(val *[]*string)
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

// The jsii proxy struct for CfnUserPoolClient
type jsiiProxy_CfnUserPoolClient struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolClient) AccessTokenValidity() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"accessTokenValidity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) AllowedOAuthFlows() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"allowedOAuthFlows",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) AllowedOAuthFlowsUserPoolClient() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"allowedOAuthFlowsUserPoolClient",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) AllowedOAuthScopes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"allowedOAuthScopes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) AnalyticsConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"analyticsConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) AttrClientSecret() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrClientSecret",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) AttrName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) CallbackUrLs() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"callbackUrLs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) ClientName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clientName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) DefaultRedirectUri() *string {
	var returns *string
	_jsii_.Get(
		j,
		"defaultRedirectUri",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) ExplicitAuthFlows() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"explicitAuthFlows",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) GenerateSecret() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"generateSecret",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) IdTokenValidity() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"idTokenValidity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) LogoutUrLs() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"logoutUrLs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) PreventUserExistenceErrors() *string {
	var returns *string
	_jsii_.Get(
		j,
		"preventUserExistenceErrors",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) ReadAttributes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"readAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) RefreshTokenValidity() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"refreshTokenValidity",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) SupportedIdentityProviders() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"supportedIdentityProviders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) TokenValidityUnits() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"tokenValidityUnits",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolClient) WriteAttributes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"writeAttributes",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolClient`.
func NewCfnUserPoolClient(scope awscdk.Construct, id *string, props *CfnUserPoolClientProps) CfnUserPoolClient {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolClient{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolClient",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolClient`.
func NewCfnUserPoolClient_Override(c CfnUserPoolClient, scope awscdk.Construct, id *string, props *CfnUserPoolClientProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolClient",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetAccessTokenValidity(val *float64) {
	_jsii_.Set(
		j,
		"accessTokenValidity",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetAllowedOAuthFlows(val *[]*string) {
	_jsii_.Set(
		j,
		"allowedOAuthFlows",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetAllowedOAuthFlowsUserPoolClient(val interface{}) {
	_jsii_.Set(
		j,
		"allowedOAuthFlowsUserPoolClient",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetAllowedOAuthScopes(val *[]*string) {
	_jsii_.Set(
		j,
		"allowedOAuthScopes",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetAnalyticsConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"analyticsConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetCallbackUrLs(val *[]*string) {
	_jsii_.Set(
		j,
		"callbackUrLs",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetClientName(val *string) {
	_jsii_.Set(
		j,
		"clientName",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetDefaultRedirectUri(val *string) {
	_jsii_.Set(
		j,
		"defaultRedirectUri",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetExplicitAuthFlows(val *[]*string) {
	_jsii_.Set(
		j,
		"explicitAuthFlows",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetGenerateSecret(val interface{}) {
	_jsii_.Set(
		j,
		"generateSecret",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetIdTokenValidity(val *float64) {
	_jsii_.Set(
		j,
		"idTokenValidity",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetLogoutUrLs(val *[]*string) {
	_jsii_.Set(
		j,
		"logoutUrLs",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetPreventUserExistenceErrors(val *string) {
	_jsii_.Set(
		j,
		"preventUserExistenceErrors",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetReadAttributes(val *[]*string) {
	_jsii_.Set(
		j,
		"readAttributes",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetRefreshTokenValidity(val *float64) {
	_jsii_.Set(
		j,
		"refreshTokenValidity",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetSupportedIdentityProviders(val *[]*string) {
	_jsii_.Set(
		j,
		"supportedIdentityProviders",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetTokenValidityUnits(val interface{}) {
	_jsii_.Set(
		j,
		"tokenValidityUnits",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolClient) SetWriteAttributes(val *[]*string) {
	_jsii_.Set(
		j,
		"writeAttributes",
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
func CfnUserPoolClient_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolClient",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolClient_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolClient",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolClient_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolClient",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolClient_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolClient",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolClient) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolClient) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolClient) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolClient) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolClient) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolClient) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolClient) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolClient) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolClient) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolClient) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolClient) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolClient) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolClient) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolClient) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolClient) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolClient) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolClient) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolClient) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolClient) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolClient) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolClient) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUserPoolClient_AnalyticsConfigurationProperty struct {
	// `CfnUserPoolClient.AnalyticsConfigurationProperty.ApplicationArn`.
	ApplicationArn *string `json:"applicationArn"`
	// `CfnUserPoolClient.AnalyticsConfigurationProperty.ApplicationId`.
	ApplicationId *string `json:"applicationId"`
	// `CfnUserPoolClient.AnalyticsConfigurationProperty.ExternalId`.
	ExternalId *string `json:"externalId"`
	// `CfnUserPoolClient.AnalyticsConfigurationProperty.RoleArn`.
	RoleArn *string `json:"roleArn"`
	// `CfnUserPoolClient.AnalyticsConfigurationProperty.UserDataShared`.
	UserDataShared interface{} `json:"userDataShared"`
}

type CfnUserPoolClient_TokenValidityUnitsProperty struct {
	// `CfnUserPoolClient.TokenValidityUnitsProperty.AccessToken`.
	AccessToken *string `json:"accessToken"`
	// `CfnUserPoolClient.TokenValidityUnitsProperty.IdToken`.
	IdToken *string `json:"idToken"`
	// `CfnUserPoolClient.TokenValidityUnitsProperty.RefreshToken`.
	RefreshToken *string `json:"refreshToken"`
}

// Properties for defining a `AWS::Cognito::UserPoolClient`.
type CfnUserPoolClientProps struct {
	// `AWS::Cognito::UserPoolClient.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolClient.AccessTokenValidity`.
	AccessTokenValidity *float64 `json:"accessTokenValidity"`
	// `AWS::Cognito::UserPoolClient.AllowedOAuthFlows`.
	AllowedOAuthFlows *[]*string `json:"allowedOAuthFlows"`
	// `AWS::Cognito::UserPoolClient.AllowedOAuthFlowsUserPoolClient`.
	AllowedOAuthFlowsUserPoolClient interface{} `json:"allowedOAuthFlowsUserPoolClient"`
	// `AWS::Cognito::UserPoolClient.AllowedOAuthScopes`.
	AllowedOAuthScopes *[]*string `json:"allowedOAuthScopes"`
	// `AWS::Cognito::UserPoolClient.AnalyticsConfiguration`.
	AnalyticsConfiguration interface{} `json:"analyticsConfiguration"`
	// `AWS::Cognito::UserPoolClient.CallbackURLs`.
	CallbackUrLs *[]*string `json:"callbackUrLs"`
	// `AWS::Cognito::UserPoolClient.ClientName`.
	ClientName *string `json:"clientName"`
	// `AWS::Cognito::UserPoolClient.DefaultRedirectURI`.
	DefaultRedirectUri *string `json:"defaultRedirectUri"`
	// `AWS::Cognito::UserPoolClient.ExplicitAuthFlows`.
	ExplicitAuthFlows *[]*string `json:"explicitAuthFlows"`
	// `AWS::Cognito::UserPoolClient.GenerateSecret`.
	GenerateSecret interface{} `json:"generateSecret"`
	// `AWS::Cognito::UserPoolClient.IdTokenValidity`.
	IdTokenValidity *float64 `json:"idTokenValidity"`
	// `AWS::Cognito::UserPoolClient.LogoutURLs`.
	LogoutUrLs *[]*string `json:"logoutUrLs"`
	// `AWS::Cognito::UserPoolClient.PreventUserExistenceErrors`.
	PreventUserExistenceErrors *string `json:"preventUserExistenceErrors"`
	// `AWS::Cognito::UserPoolClient.ReadAttributes`.
	ReadAttributes *[]*string `json:"readAttributes"`
	// `AWS::Cognito::UserPoolClient.RefreshTokenValidity`.
	RefreshTokenValidity *float64 `json:"refreshTokenValidity"`
	// `AWS::Cognito::UserPoolClient.SupportedIdentityProviders`.
	SupportedIdentityProviders *[]*string `json:"supportedIdentityProviders"`
	// `AWS::Cognito::UserPoolClient.TokenValidityUnits`.
	TokenValidityUnits interface{} `json:"tokenValidityUnits"`
	// `AWS::Cognito::UserPoolClient.WriteAttributes`.
	WriteAttributes *[]*string `json:"writeAttributes"`
}

// A CloudFormation `AWS::Cognito::UserPoolDomain`.
type CfnUserPoolDomain interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	CustomDomainConfig() interface{}
	SetCustomDomainConfig(val interface{})
	Domain() *string
	SetDomain(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserPoolId() *string
	SetUserPoolId(val *string)
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

// The jsii proxy struct for CfnUserPoolDomain
type jsiiProxy_CfnUserPoolDomain struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolDomain) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) CustomDomainConfig() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"customDomainConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) Domain() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domain",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolDomain) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolDomain`.
func NewCfnUserPoolDomain(scope awscdk.Construct, id *string, props *CfnUserPoolDomainProps) CfnUserPoolDomain {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolDomain{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolDomain",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolDomain`.
func NewCfnUserPoolDomain_Override(c CfnUserPoolDomain, scope awscdk.Construct, id *string, props *CfnUserPoolDomainProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolDomain",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolDomain) SetCustomDomainConfig(val interface{}) {
	_jsii_.Set(
		j,
		"customDomainConfig",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolDomain) SetDomain(val *string) {
	_jsii_.Set(
		j,
		"domain",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolDomain) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
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
func CfnUserPoolDomain_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolDomain",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolDomain_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolDomain",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolDomain_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolDomain",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolDomain_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolDomain",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolDomain) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolDomain) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolDomain) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolDomain) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolDomain) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolDomain) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolDomain) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolDomain) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolDomain) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolDomain) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolDomain) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolDomain) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolDomain) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolDomain) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolDomain) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolDomain) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolDomain) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolDomain) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolDomain) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolDomain) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolDomain) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUserPoolDomain_CustomDomainConfigTypeProperty struct {
	// `CfnUserPoolDomain.CustomDomainConfigTypeProperty.CertificateArn`.
	CertificateArn *string `json:"certificateArn"`
}

// Properties for defining a `AWS::Cognito::UserPoolDomain`.
type CfnUserPoolDomainProps struct {
	// `AWS::Cognito::UserPoolDomain.Domain`.
	Domain *string `json:"domain"`
	// `AWS::Cognito::UserPoolDomain.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolDomain.CustomDomainConfig`.
	CustomDomainConfig interface{} `json:"customDomainConfig"`
}

// A CloudFormation `AWS::Cognito::UserPoolGroup`.
type CfnUserPoolGroup interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	GroupName() *string
	SetGroupName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Precedence() *float64
	SetPrecedence(val *float64)
	Ref() *string
	RoleArn() *string
	SetRoleArn(val *string)
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserPoolId() *string
	SetUserPoolId(val *string)
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

// The jsii proxy struct for CfnUserPoolGroup
type jsiiProxy_CfnUserPoolGroup struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolGroup) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) GroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) Precedence() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"precedence",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolGroup) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolGroup`.
func NewCfnUserPoolGroup(scope awscdk.Construct, id *string, props *CfnUserPoolGroupProps) CfnUserPoolGroup {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolGroup{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolGroup`.
func NewCfnUserPoolGroup_Override(c CfnUserPoolGroup, scope awscdk.Construct, id *string, props *CfnUserPoolGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolGroup",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolGroup) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolGroup) SetGroupName(val *string) {
	_jsii_.Set(
		j,
		"groupName",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolGroup) SetPrecedence(val *float64) {
	_jsii_.Set(
		j,
		"precedence",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolGroup) SetRoleArn(val *string) {
	_jsii_.Set(
		j,
		"roleArn",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolGroup) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
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
func CfnUserPoolGroup_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolGroup",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolGroup_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolGroup",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolGroup_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolGroup",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolGroup) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolGroup) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolGroup) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolGroup) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolGroup) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolGroup) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolGroup) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolGroup) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolGroup) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolGroup) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolGroup) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolGroup) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolGroup) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolGroup) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolGroup) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolGroup) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolGroup) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolGroup) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolGroup) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Cognito::UserPoolGroup`.
type CfnUserPoolGroupProps struct {
	// `AWS::Cognito::UserPoolGroup.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolGroup.Description`.
	Description *string `json:"description"`
	// `AWS::Cognito::UserPoolGroup.GroupName`.
	GroupName *string `json:"groupName"`
	// `AWS::Cognito::UserPoolGroup.Precedence`.
	Precedence *float64 `json:"precedence"`
	// `AWS::Cognito::UserPoolGroup.RoleArn`.
	RoleArn *string `json:"roleArn"`
}

// A CloudFormation `AWS::Cognito::UserPoolIdentityProvider`.
type CfnUserPoolIdentityProvider interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttributeMapping() interface{}
	SetAttributeMapping(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	IdpIdentifiers() *[]*string
	SetIdpIdentifiers(val *[]*string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	ProviderDetails() interface{}
	SetProviderDetails(val interface{})
	ProviderName() *string
	SetProviderName(val *string)
	ProviderType() *string
	SetProviderType(val *string)
	Ref() *string
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserPoolId() *string
	SetUserPoolId(val *string)
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

// The jsii proxy struct for CfnUserPoolIdentityProvider
type jsiiProxy_CfnUserPoolIdentityProvider struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) AttributeMapping() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"attributeMapping",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) IdpIdentifiers() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"idpIdentifiers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) ProviderDetails() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"providerDetails",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) ProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"providerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) ProviderType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"providerType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolIdentityProvider`.
func NewCfnUserPoolIdentityProvider(scope awscdk.Construct, id *string, props *CfnUserPoolIdentityProviderProps) CfnUserPoolIdentityProvider {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolIdentityProvider{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolIdentityProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolIdentityProvider`.
func NewCfnUserPoolIdentityProvider_Override(c CfnUserPoolIdentityProvider, scope awscdk.Construct, id *string, props *CfnUserPoolIdentityProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolIdentityProvider",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) SetAttributeMapping(val interface{}) {
	_jsii_.Set(
		j,
		"attributeMapping",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) SetIdpIdentifiers(val *[]*string) {
	_jsii_.Set(
		j,
		"idpIdentifiers",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) SetProviderDetails(val interface{}) {
	_jsii_.Set(
		j,
		"providerDetails",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) SetProviderName(val *string) {
	_jsii_.Set(
		j,
		"providerName",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) SetProviderType(val *string) {
	_jsii_.Set(
		j,
		"providerType",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolIdentityProvider) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
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
func CfnUserPoolIdentityProvider_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolIdentityProvider",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolIdentityProvider_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolIdentityProvider",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolIdentityProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolIdentityProvider",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolIdentityProvider_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolIdentityProvider",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolIdentityProvider) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolIdentityProvider) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolIdentityProvider) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolIdentityProvider) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolIdentityProvider) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Cognito::UserPoolIdentityProvider`.
type CfnUserPoolIdentityProviderProps struct {
	// `AWS::Cognito::UserPoolIdentityProvider.ProviderName`.
	ProviderName *string `json:"providerName"`
	// `AWS::Cognito::UserPoolIdentityProvider.ProviderType`.
	ProviderType *string `json:"providerType"`
	// `AWS::Cognito::UserPoolIdentityProvider.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolIdentityProvider.AttributeMapping`.
	AttributeMapping interface{} `json:"attributeMapping"`
	// `AWS::Cognito::UserPoolIdentityProvider.IdpIdentifiers`.
	IdpIdentifiers *[]*string `json:"idpIdentifiers"`
	// `AWS::Cognito::UserPoolIdentityProvider.ProviderDetails`.
	ProviderDetails interface{} `json:"providerDetails"`
}

// Properties for defining a `AWS::Cognito::UserPool`.
type CfnUserPoolProps struct {
	// `AWS::Cognito::UserPool.AccountRecoverySetting`.
	AccountRecoverySetting interface{} `json:"accountRecoverySetting"`
	// `AWS::Cognito::UserPool.AdminCreateUserConfig`.
	AdminCreateUserConfig interface{} `json:"adminCreateUserConfig"`
	// `AWS::Cognito::UserPool.AliasAttributes`.
	AliasAttributes *[]*string `json:"aliasAttributes"`
	// `AWS::Cognito::UserPool.AutoVerifiedAttributes`.
	AutoVerifiedAttributes *[]*string `json:"autoVerifiedAttributes"`
	// `AWS::Cognito::UserPool.DeviceConfiguration`.
	DeviceConfiguration interface{} `json:"deviceConfiguration"`
	// `AWS::Cognito::UserPool.EmailConfiguration`.
	EmailConfiguration interface{} `json:"emailConfiguration"`
	// `AWS::Cognito::UserPool.EmailVerificationMessage`.
	EmailVerificationMessage *string `json:"emailVerificationMessage"`
	// `AWS::Cognito::UserPool.EmailVerificationSubject`.
	EmailVerificationSubject *string `json:"emailVerificationSubject"`
	// `AWS::Cognito::UserPool.EnabledMfas`.
	EnabledMfas *[]*string `json:"enabledMfas"`
	// `AWS::Cognito::UserPool.LambdaConfig`.
	LambdaConfig interface{} `json:"lambdaConfig"`
	// `AWS::Cognito::UserPool.MfaConfiguration`.
	MfaConfiguration *string `json:"mfaConfiguration"`
	// `AWS::Cognito::UserPool.Policies`.
	Policies interface{} `json:"policies"`
	// `AWS::Cognito::UserPool.Schema`.
	Schema interface{} `json:"schema"`
	// `AWS::Cognito::UserPool.SmsAuthenticationMessage`.
	SmsAuthenticationMessage *string `json:"smsAuthenticationMessage"`
	// `AWS::Cognito::UserPool.SmsConfiguration`.
	SmsConfiguration interface{} `json:"smsConfiguration"`
	// `AWS::Cognito::UserPool.SmsVerificationMessage`.
	SmsVerificationMessage *string `json:"smsVerificationMessage"`
	// `AWS::Cognito::UserPool.UsernameAttributes`.
	UsernameAttributes *[]*string `json:"usernameAttributes"`
	// `AWS::Cognito::UserPool.UsernameConfiguration`.
	UsernameConfiguration interface{} `json:"usernameConfiguration"`
	// `AWS::Cognito::UserPool.UserPoolAddOns`.
	UserPoolAddOns interface{} `json:"userPoolAddOns"`
	// `AWS::Cognito::UserPool.UserPoolName`.
	UserPoolName *string `json:"userPoolName"`
	// `AWS::Cognito::UserPool.UserPoolTags`.
	UserPoolTags interface{} `json:"userPoolTags"`
	// `AWS::Cognito::UserPool.VerificationMessageTemplate`.
	VerificationMessageTemplate interface{} `json:"verificationMessageTemplate"`
}

// A CloudFormation `AWS::Cognito::UserPoolResourceServer`.
type CfnUserPoolResourceServer interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Identifier() *string
	SetIdentifier(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Scopes() interface{}
	SetScopes(val interface{})
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserPoolId() *string
	SetUserPoolId(val *string)
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

// The jsii proxy struct for CfnUserPoolResourceServer
type jsiiProxy_CfnUserPoolResourceServer struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolResourceServer) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) Identifier() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identifier",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) Scopes() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"scopes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolResourceServer) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolResourceServer`.
func NewCfnUserPoolResourceServer(scope awscdk.Construct, id *string, props *CfnUserPoolResourceServerProps) CfnUserPoolResourceServer {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolResourceServer{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolResourceServer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolResourceServer`.
func NewCfnUserPoolResourceServer_Override(c CfnUserPoolResourceServer, scope awscdk.Construct, id *string, props *CfnUserPoolResourceServerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolResourceServer",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolResourceServer) SetIdentifier(val *string) {
	_jsii_.Set(
		j,
		"identifier",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolResourceServer) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolResourceServer) SetScopes(val interface{}) {
	_jsii_.Set(
		j,
		"scopes",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolResourceServer) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
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
func CfnUserPoolResourceServer_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolResourceServer",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolResourceServer_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolResourceServer",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolResourceServer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolResourceServer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolResourceServer_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolResourceServer",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolResourceServer) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolResourceServer) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolResourceServer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolResourceServer) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolResourceServer) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUserPoolResourceServer_ResourceServerScopeTypeProperty struct {
	// `CfnUserPoolResourceServer.ResourceServerScopeTypeProperty.ScopeDescription`.
	ScopeDescription *string `json:"scopeDescription"`
	// `CfnUserPoolResourceServer.ResourceServerScopeTypeProperty.ScopeName`.
	ScopeName *string `json:"scopeName"`
}

// Properties for defining a `AWS::Cognito::UserPoolResourceServer`.
type CfnUserPoolResourceServerProps struct {
	// `AWS::Cognito::UserPoolResourceServer.Identifier`.
	Identifier *string `json:"identifier"`
	// `AWS::Cognito::UserPoolResourceServer.Name`.
	Name *string `json:"name"`
	// `AWS::Cognito::UserPoolResourceServer.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolResourceServer.Scopes`.
	Scopes interface{} `json:"scopes"`
}

// A CloudFormation `AWS::Cognito::UserPoolRiskConfigurationAttachment`.
type CfnUserPoolRiskConfigurationAttachment interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AccountTakeoverRiskConfiguration() interface{}
	SetAccountTakeoverRiskConfiguration(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClientId() *string
	SetClientId(val *string)
	CompromisedCredentialsRiskConfiguration() interface{}
	SetCompromisedCredentialsRiskConfiguration(val interface{})
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RiskExceptionConfiguration() interface{}
	SetRiskExceptionConfiguration(val interface{})
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserPoolId() *string
	SetUserPoolId(val *string)
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

// The jsii proxy struct for CfnUserPoolRiskConfigurationAttachment
type jsiiProxy_CfnUserPoolRiskConfigurationAttachment struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) AccountTakeoverRiskConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accountTakeoverRiskConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) ClientId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clientId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) CompromisedCredentialsRiskConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"compromisedCredentialsRiskConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) RiskExceptionConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"riskExceptionConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolRiskConfigurationAttachment`.
func NewCfnUserPoolRiskConfigurationAttachment(scope awscdk.Construct, id *string, props *CfnUserPoolRiskConfigurationAttachmentProps) CfnUserPoolRiskConfigurationAttachment {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolRiskConfigurationAttachment{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolRiskConfigurationAttachment`.
func NewCfnUserPoolRiskConfigurationAttachment_Override(c CfnUserPoolRiskConfigurationAttachment, scope awscdk.Construct, id *string, props *CfnUserPoolRiskConfigurationAttachmentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) SetAccountTakeoverRiskConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"accountTakeoverRiskConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) SetClientId(val *string) {
	_jsii_.Set(
		j,
		"clientId",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) SetCompromisedCredentialsRiskConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"compromisedCredentialsRiskConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) SetRiskExceptionConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"riskExceptionConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
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
func CfnUserPoolRiskConfigurationAttachment_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolRiskConfigurationAttachment_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolRiskConfigurationAttachment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolRiskConfigurationAttachment_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolRiskConfigurationAttachment",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolRiskConfigurationAttachment) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUserPoolRiskConfigurationAttachment_AccountTakeoverActionTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.AccountTakeoverActionTypeProperty.EventAction`.
	EventAction *string `json:"eventAction"`
	// `CfnUserPoolRiskConfigurationAttachment.AccountTakeoverActionTypeProperty.Notify`.
	Notify interface{} `json:"notify"`
}

type CfnUserPoolRiskConfigurationAttachment_AccountTakeoverActionsTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.AccountTakeoverActionsTypeProperty.HighAction`.
	HighAction interface{} `json:"highAction"`
	// `CfnUserPoolRiskConfigurationAttachment.AccountTakeoverActionsTypeProperty.LowAction`.
	LowAction interface{} `json:"lowAction"`
	// `CfnUserPoolRiskConfigurationAttachment.AccountTakeoverActionsTypeProperty.MediumAction`.
	MediumAction interface{} `json:"mediumAction"`
}

type CfnUserPoolRiskConfigurationAttachment_AccountTakeoverRiskConfigurationTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.AccountTakeoverRiskConfigurationTypeProperty.Actions`.
	Actions interface{} `json:"actions"`
	// `CfnUserPoolRiskConfigurationAttachment.AccountTakeoverRiskConfigurationTypeProperty.NotifyConfiguration`.
	NotifyConfiguration interface{} `json:"notifyConfiguration"`
}

type CfnUserPoolRiskConfigurationAttachment_CompromisedCredentialsActionsTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.CompromisedCredentialsActionsTypeProperty.EventAction`.
	EventAction *string `json:"eventAction"`
}

type CfnUserPoolRiskConfigurationAttachment_CompromisedCredentialsRiskConfigurationTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.CompromisedCredentialsRiskConfigurationTypeProperty.Actions`.
	Actions interface{} `json:"actions"`
	// `CfnUserPoolRiskConfigurationAttachment.CompromisedCredentialsRiskConfigurationTypeProperty.EventFilter`.
	EventFilter *[]*string `json:"eventFilter"`
}

type CfnUserPoolRiskConfigurationAttachment_NotifyConfigurationTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.NotifyConfigurationTypeProperty.SourceArn`.
	SourceArn *string `json:"sourceArn"`
	// `CfnUserPoolRiskConfigurationAttachment.NotifyConfigurationTypeProperty.BlockEmail`.
	BlockEmail interface{} `json:"blockEmail"`
	// `CfnUserPoolRiskConfigurationAttachment.NotifyConfigurationTypeProperty.From`.
	From *string `json:"from"`
	// `CfnUserPoolRiskConfigurationAttachment.NotifyConfigurationTypeProperty.MfaEmail`.
	MfaEmail interface{} `json:"mfaEmail"`
	// `CfnUserPoolRiskConfigurationAttachment.NotifyConfigurationTypeProperty.NoActionEmail`.
	NoActionEmail interface{} `json:"noActionEmail"`
	// `CfnUserPoolRiskConfigurationAttachment.NotifyConfigurationTypeProperty.ReplyTo`.
	ReplyTo *string `json:"replyTo"`
}

type CfnUserPoolRiskConfigurationAttachment_NotifyEmailTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.NotifyEmailTypeProperty.Subject`.
	Subject *string `json:"subject"`
	// `CfnUserPoolRiskConfigurationAttachment.NotifyEmailTypeProperty.HtmlBody`.
	HtmlBody *string `json:"htmlBody"`
	// `CfnUserPoolRiskConfigurationAttachment.NotifyEmailTypeProperty.TextBody`.
	TextBody *string `json:"textBody"`
}

type CfnUserPoolRiskConfigurationAttachment_RiskExceptionConfigurationTypeProperty struct {
	// `CfnUserPoolRiskConfigurationAttachment.RiskExceptionConfigurationTypeProperty.BlockedIPRangeList`.
	BlockedIpRangeList *[]*string `json:"blockedIpRangeList"`
	// `CfnUserPoolRiskConfigurationAttachment.RiskExceptionConfigurationTypeProperty.SkippedIPRangeList`.
	SkippedIpRangeList *[]*string `json:"skippedIpRangeList"`
}

// Properties for defining a `AWS::Cognito::UserPoolRiskConfigurationAttachment`.
type CfnUserPoolRiskConfigurationAttachmentProps struct {
	// `AWS::Cognito::UserPoolRiskConfigurationAttachment.ClientId`.
	ClientId *string `json:"clientId"`
	// `AWS::Cognito::UserPoolRiskConfigurationAttachment.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolRiskConfigurationAttachment.AccountTakeoverRiskConfiguration`.
	AccountTakeoverRiskConfiguration interface{} `json:"accountTakeoverRiskConfiguration"`
	// `AWS::Cognito::UserPoolRiskConfigurationAttachment.CompromisedCredentialsRiskConfiguration`.
	CompromisedCredentialsRiskConfiguration interface{} `json:"compromisedCredentialsRiskConfiguration"`
	// `AWS::Cognito::UserPoolRiskConfigurationAttachment.RiskExceptionConfiguration`.
	RiskExceptionConfiguration interface{} `json:"riskExceptionConfiguration"`
}

// A CloudFormation `AWS::Cognito::UserPoolUICustomizationAttachment`.
type CfnUserPoolUICustomizationAttachment interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClientId() *string
	SetClientId(val *string)
	CreationStack() *[]*string
	Css() *string
	SetCss(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserPoolId() *string
	SetUserPoolId(val *string)
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

// The jsii proxy struct for CfnUserPoolUICustomizationAttachment
type jsiiProxy_CfnUserPoolUICustomizationAttachment struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) ClientId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clientId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) Css() *string {
	var returns *string
	_jsii_.Get(
		j,
		"css",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolUICustomizationAttachment`.
func NewCfnUserPoolUICustomizationAttachment(scope awscdk.Construct, id *string, props *CfnUserPoolUICustomizationAttachmentProps) CfnUserPoolUICustomizationAttachment {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolUICustomizationAttachment{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolUICustomizationAttachment`.
func NewCfnUserPoolUICustomizationAttachment_Override(c CfnUserPoolUICustomizationAttachment, scope awscdk.Construct, id *string, props *CfnUserPoolUICustomizationAttachmentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachment",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) SetClientId(val *string) {
	_jsii_.Set(
		j,
		"clientId",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) SetCss(val *string) {
	_jsii_.Set(
		j,
		"css",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUICustomizationAttachment) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
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
func CfnUserPoolUICustomizationAttachment_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachment",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolUICustomizationAttachment_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachment",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolUICustomizationAttachment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolUICustomizationAttachment_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolUICustomizationAttachment",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolUICustomizationAttachment) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Cognito::UserPoolUICustomizationAttachment`.
type CfnUserPoolUICustomizationAttachmentProps struct {
	// `AWS::Cognito::UserPoolUICustomizationAttachment.ClientId`.
	ClientId *string `json:"clientId"`
	// `AWS::Cognito::UserPoolUICustomizationAttachment.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolUICustomizationAttachment.CSS`.
	Css *string `json:"css"`
}

// A CloudFormation `AWS::Cognito::UserPoolUser`.
type CfnUserPoolUser interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClientMetadata() interface{}
	SetClientMetadata(val interface{})
	CreationStack() *[]*string
	DesiredDeliveryMediums() *[]*string
	SetDesiredDeliveryMediums(val *[]*string)
	ForceAliasCreation() interface{}
	SetForceAliasCreation(val interface{})
	LogicalId() *string
	MessageAction() *string
	SetMessageAction(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UserAttributes() interface{}
	SetUserAttributes(val interface{})
	Username() *string
	SetUsername(val *string)
	UserPoolId() *string
	SetUserPoolId(val *string)
	ValidationData() interface{}
	SetValidationData(val interface{})
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

// The jsii proxy struct for CfnUserPoolUser
type jsiiProxy_CfnUserPoolUser struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolUser) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) ClientMetadata() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"clientMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) DesiredDeliveryMediums() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"desiredDeliveryMediums",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) ForceAliasCreation() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"forceAliasCreation",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) MessageAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"messageAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) UserAttributes() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"userAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) Username() *string {
	var returns *string
	_jsii_.Get(
		j,
		"username",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUser) ValidationData() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"validationData",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolUser`.
func NewCfnUserPoolUser(scope awscdk.Construct, id *string, props *CfnUserPoolUserProps) CfnUserPoolUser {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolUser{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolUser",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolUser`.
func NewCfnUserPoolUser_Override(c CfnUserPoolUser, scope awscdk.Construct, id *string, props *CfnUserPoolUserProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolUser",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetClientMetadata(val interface{}) {
	_jsii_.Set(
		j,
		"clientMetadata",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetDesiredDeliveryMediums(val *[]*string) {
	_jsii_.Set(
		j,
		"desiredDeliveryMediums",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetForceAliasCreation(val interface{}) {
	_jsii_.Set(
		j,
		"forceAliasCreation",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetMessageAction(val *string) {
	_jsii_.Set(
		j,
		"messageAction",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetUserAttributes(val interface{}) {
	_jsii_.Set(
		j,
		"userAttributes",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetUsername(val *string) {
	_jsii_.Set(
		j,
		"username",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUser) SetValidationData(val interface{}) {
	_jsii_.Set(
		j,
		"validationData",
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
func CfnUserPoolUser_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUser",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolUser_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUser",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolUser_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUser",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolUser_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolUser",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUser) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolUser) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolUser) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolUser) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUser) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolUser) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUser) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolUser) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolUser) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolUser) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolUser) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolUser) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolUser) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolUser) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolUser) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolUser) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolUser) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolUser) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolUser) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolUser) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolUser) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUserPoolUser_AttributeTypeProperty struct {
	// `CfnUserPoolUser.AttributeTypeProperty.Name`.
	Name *string `json:"name"`
	// `CfnUserPoolUser.AttributeTypeProperty.Value`.
	Value *string `json:"value"`
}

// Properties for defining a `AWS::Cognito::UserPoolUser`.
type CfnUserPoolUserProps struct {
	// `AWS::Cognito::UserPoolUser.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
	// `AWS::Cognito::UserPoolUser.ClientMetadata`.
	ClientMetadata interface{} `json:"clientMetadata"`
	// `AWS::Cognito::UserPoolUser.DesiredDeliveryMediums`.
	DesiredDeliveryMediums *[]*string `json:"desiredDeliveryMediums"`
	// `AWS::Cognito::UserPoolUser.ForceAliasCreation`.
	ForceAliasCreation interface{} `json:"forceAliasCreation"`
	// `AWS::Cognito::UserPoolUser.MessageAction`.
	MessageAction *string `json:"messageAction"`
	// `AWS::Cognito::UserPoolUser.UserAttributes`.
	UserAttributes interface{} `json:"userAttributes"`
	// `AWS::Cognito::UserPoolUser.Username`.
	Username *string `json:"username"`
	// `AWS::Cognito::UserPoolUser.ValidationData`.
	ValidationData interface{} `json:"validationData"`
}

// A CloudFormation `AWS::Cognito::UserPoolUserToGroupAttachment`.
type CfnUserPoolUserToGroupAttachment interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	GroupName() *string
	SetGroupName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	Username() *string
	SetUsername(val *string)
	UserPoolId() *string
	SetUserPoolId(val *string)
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

// The jsii proxy struct for CfnUserPoolUserToGroupAttachment
type jsiiProxy_CfnUserPoolUserToGroupAttachment struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) GroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) Username() *string {
	var returns *string
	_jsii_.Get(
		j,
		"username",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}


// Create a new `AWS::Cognito::UserPoolUserToGroupAttachment`.
func NewCfnUserPoolUserToGroupAttachment(scope awscdk.Construct, id *string, props *CfnUserPoolUserToGroupAttachmentProps) CfnUserPoolUserToGroupAttachment {
	_init_.Initialize()

	j := jsiiProxy_CfnUserPoolUserToGroupAttachment{}

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::Cognito::UserPoolUserToGroupAttachment`.
func NewCfnUserPoolUserToGroupAttachment_Override(c CfnUserPoolUserToGroupAttachment, scope awscdk.Construct, id *string, props *CfnUserPoolUserToGroupAttachmentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachment",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) SetGroupName(val *string) {
	_jsii_.Set(
		j,
		"groupName",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) SetUsername(val *string) {
	_jsii_.Set(
		j,
		"username",
		val,
	)
}

func (j *jsiiProxy_CfnUserPoolUserToGroupAttachment) SetUserPoolId(val *string) {
	_jsii_.Set(
		j,
		"userPoolId",
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
func CfnUserPoolUserToGroupAttachment_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachment",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserPoolUserToGroupAttachment_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachment",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserPoolUserToGroupAttachment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserPoolUserToGroupAttachment_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_cognito.CfnUserPoolUserToGroupAttachment",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) OnPrepare() {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) ToString() *string {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserPoolUserToGroupAttachment) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::Cognito::UserPoolUserToGroupAttachment`.
type CfnUserPoolUserToGroupAttachmentProps struct {
	// `AWS::Cognito::UserPoolUserToGroupAttachment.GroupName`.
	GroupName *string `json:"groupName"`
	// `AWS::Cognito::UserPoolUserToGroupAttachment.Username`.
	Username *string `json:"username"`
	// `AWS::Cognito::UserPoolUserToGroupAttachment.UserPoolId`.
	UserPoolId *string `json:"userPoolId"`
}

// A set of attributes, useful to set Read and Write attributes.
// Experimental.
type ClientAttributes interface {
	Attributes() *[]*string
	WithCustomAttributes(attributes ...*string) ClientAttributes
	WithStandardAttributes(attributes *StandardAttributesMask) ClientAttributes
}

// The jsii proxy struct for ClientAttributes
type jsiiProxy_ClientAttributes struct {
	_ byte // padding
}

// Creates a ClientAttributes with the specified attributes.
// Experimental.
func NewClientAttributes() ClientAttributes {
	_init_.Initialize()

	j := jsiiProxy_ClientAttributes{}

	_jsii_.Create(
		"monocdk.aws_cognito.ClientAttributes",
		nil, // no parameters
		&j,
	)

	return &j
}

// Creates a ClientAttributes with the specified attributes.
// Experimental.
func NewClientAttributes_Override(c ClientAttributes) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.ClientAttributes",
		nil, // no parameters
		c,
	)
}

// The list of attributes represented by this ClientAttributes.
// Experimental.
func (c *jsiiProxy_ClientAttributes) Attributes() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"attributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Creates a custom ClientAttributes with the specified attributes.
// Experimental.
func (c *jsiiProxy_ClientAttributes) WithCustomAttributes(attributes ...*string) ClientAttributes {
	args := []interface{}{}
	for _, a := range attributes {
		args = append(args, a)
	}

	var returns ClientAttributes

	_jsii_.Invoke(
		c,
		"withCustomAttributes",
		args,
		&returns,
	)

	return returns
}

// Creates a custom ClientAttributes with the specified attributes.
// Experimental.
func (c *jsiiProxy_ClientAttributes) WithStandardAttributes(attributes *StandardAttributesMask) ClientAttributes {
	var returns ClientAttributes

	_jsii_.Invoke(
		c,
		"withStandardAttributes",
		[]interface{}{attributes},
		&returns,
	)

	return returns
}

// Options while specifying a cognito prefix domain.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-assign-domain-prefix.html
//
// Experimental.
type CognitoDomainOptions struct {
	// The prefix to the Cognito hosted domain name that will be associated with the user pool.
	// Experimental.
	DomainPrefix *string `json:"domainPrefix"`
}

// Configuration that will be fed into CloudFormation for any custom attribute type.
// Experimental.
type CustomAttributeConfig struct {
	// The data type of the custom attribute.
	// See: https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_SchemaAttributeType.html#CognitoUserPools-Type-SchemaAttributeType-AttributeDataType
	//
	// Experimental.
	DataType *string `json:"dataType"`
	// Specifies whether the value of the attribute can be changed.
	//
	// For any user pool attribute that's mapped to an identity provider attribute, you must set this parameter to true.
	// Amazon Cognito updates mapped attributes when users sign in to your application through an identity provider.
	// If an attribute is immutable, Amazon Cognito throws an error when it attempts to update the attribute.
	// Experimental.
	Mutable *bool `json:"mutable"`
	// The constraints for a custom attribute of the 'Number' data type.
	// Experimental.
	NumberConstraints *NumberAttributeConstraints `json:"numberConstraints"`
	// The constraints for a custom attribute of 'String' data type.
	// Experimental.
	StringConstraints *StringAttributeConstraints `json:"stringConstraints"`
}

// Constraints that can be applied to a custom attribute of any type.
// Experimental.
type CustomAttributeProps struct {
	// Specifies whether the value of the attribute can be changed.
	//
	// For any user pool attribute that's mapped to an identity provider attribute, you must set this parameter to true.
	// Amazon Cognito updates mapped attributes when users sign in to your application through an identity provider.
	// If an attribute is immutable, Amazon Cognito throws an error when it attempts to update the attribute.
	// Experimental.
	Mutable *bool `json:"mutable"`
}

// Options while specifying custom domain.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-add-custom-domain.html
//
// Experimental.
type CustomDomainOptions struct {
	// The certificate to associate with this domain.
	// Experimental.
	Certificate awscertificatemanager.ICertificate `json:"certificate"`
	// The custom domain name that you would like to associate with this User Pool.
	// Experimental.
	DomainName *string `json:"domainName"`
}

// The DateTime custom attribute type.
// Experimental.
type DateTimeAttribute interface {
	ICustomAttribute
	Bind() *CustomAttributeConfig
}

// The jsii proxy struct for DateTimeAttribute
type jsiiProxy_DateTimeAttribute struct {
	jsiiProxy_ICustomAttribute
}

// Experimental.
func NewDateTimeAttribute(props *CustomAttributeProps) DateTimeAttribute {
	_init_.Initialize()

	j := jsiiProxy_DateTimeAttribute{}

	_jsii_.Create(
		"monocdk.aws_cognito.DateTimeAttribute",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewDateTimeAttribute_Override(d DateTimeAttribute, props *CustomAttributeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.DateTimeAttribute",
		[]interface{}{props},
		d,
	)
}

// Bind this custom attribute type to the values as expected by CloudFormation.
// Experimental.
func (d *jsiiProxy_DateTimeAttribute) Bind() *CustomAttributeConfig {
	var returns *CustomAttributeConfig

	_jsii_.Invoke(
		d,
		"bind",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Email settings for the user pool.
// Experimental.
type EmailSettings struct {
	// The 'from' address on the emails received by the user.
	// Experimental.
	From *string `json:"from"`
	// The 'replyTo' address on the emails received by the user as defined by IETF RFC-5322.
	//
	// When set, most email clients recognize to change 'to' line to this address when a reply is drafted.
	// Experimental.
	ReplyTo *string `json:"replyTo"`
}

// Represents a custom attribute type.
// Experimental.
type ICustomAttribute interface {
	// Bind this custom attribute type to the values as expected by CloudFormation.
	// Experimental.
	Bind() *CustomAttributeConfig
}

// The jsii proxy for ICustomAttribute
type jsiiProxy_ICustomAttribute struct {
	_ byte // padding
}

func (i *jsiiProxy_ICustomAttribute) Bind() *CustomAttributeConfig {
	var returns *CustomAttributeConfig

	_jsii_.Invoke(
		i,
		"bind",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents a Cognito UserPool.
// Experimental.
type IUserPool interface {
	awscdk.IResource
	// Add a new app client to this user pool.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-client-apps.html
	//
	// Experimental.
	AddClient(id *string, options *UserPoolClientOptions) UserPoolClient
	// Associate a domain to this user pool.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-assign-domain.html
	//
	// Experimental.
	AddDomain(id *string, options *UserPoolDomainOptions) UserPoolDomain
	// Add a new resource server to this user pool.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-resource-servers.html
	//
	// Experimental.
	AddResourceServer(id *string, options *UserPoolResourceServerOptions) UserPoolResourceServer
	// Register an identity provider with this user pool.
	// Experimental.
	RegisterIdentityProvider(provider IUserPoolIdentityProvider)
	// Get all identity providers registered with this user pool.
	// Experimental.
	IdentityProviders() *[]IUserPoolIdentityProvider
	// The ARN of this user pool resource.
	// Experimental.
	UserPoolArn() *string
	// The physical ID of this user pool resource.
	// Experimental.
	UserPoolId() *string
}

// The jsii proxy for IUserPool
type jsiiProxy_IUserPool struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IUserPool) AddClient(id *string, options *UserPoolClientOptions) UserPoolClient {
	var returns UserPoolClient

	_jsii_.Invoke(
		i,
		"addClient",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IUserPool) AddDomain(id *string, options *UserPoolDomainOptions) UserPoolDomain {
	var returns UserPoolDomain

	_jsii_.Invoke(
		i,
		"addDomain",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IUserPool) AddResourceServer(id *string, options *UserPoolResourceServerOptions) UserPoolResourceServer {
	var returns UserPoolResourceServer

	_jsii_.Invoke(
		i,
		"addResourceServer",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IUserPool) RegisterIdentityProvider(provider IUserPoolIdentityProvider) {
	_jsii_.InvokeVoid(
		i,
		"registerIdentityProvider",
		[]interface{}{provider},
	)
}

func (j *jsiiProxy_IUserPool) IdentityProviders() *[]IUserPoolIdentityProvider {
	var returns *[]IUserPoolIdentityProvider
	_jsii_.Get(
		j,
		"identityProviders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IUserPool) UserPoolArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IUserPool) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}

// Represents a Cognito user pool client.
// Experimental.
type IUserPoolClient interface {
	awscdk.IResource
	// Name of the application client.
	// Experimental.
	UserPoolClientId() *string
}

// The jsii proxy for IUserPoolClient
type jsiiProxy_IUserPoolClient struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IUserPoolClient) UserPoolClientId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolClientId",
		&returns,
	)
	return returns
}

// Represents a user pool domain.
// Experimental.
type IUserPoolDomain interface {
	awscdk.IResource
	// The domain that was specified to be created.
	//
	// If `customDomain` was selected, this holds the full domain name that was specified.
	// If the `cognitoDomain` was used, it contains the prefix to the Cognito hosted domain.
	// Experimental.
	DomainName() *string
}

// The jsii proxy for IUserPoolDomain
type jsiiProxy_IUserPoolDomain struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IUserPoolDomain) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

// Represents a UserPoolIdentityProvider.
// Experimental.
type IUserPoolIdentityProvider interface {
	awscdk.IResource
	// The primary identifier of this identity provider.
	// Experimental.
	ProviderName() *string
}

// The jsii proxy for IUserPoolIdentityProvider
type jsiiProxy_IUserPoolIdentityProvider struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IUserPoolIdentityProvider) ProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"providerName",
		&returns,
	)
	return returns
}

// Represents a Cognito user pool resource server.
// Experimental.
type IUserPoolResourceServer interface {
	awscdk.IResource
	// Resource server id.
	// Experimental.
	UserPoolResourceServerId() *string
}

// The jsii proxy for IUserPoolResourceServer
type jsiiProxy_IUserPoolResourceServer struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IUserPoolResourceServer) UserPoolResourceServerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolResourceServerId",
		&returns,
	)
	return returns
}

// The different ways in which a user pool's MFA enforcement can be configured.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-mfa.html
//
// Experimental.
type Mfa string

const (
	Mfa_OFF Mfa = "OFF"
	Mfa_OPTIONAL Mfa = "OPTIONAL"
	Mfa_REQUIRED Mfa = "REQUIRED"
)

// The different ways in which a user pool can obtain their MFA token for sign in.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-mfa.html
//
// Experimental.
type MfaSecondFactor struct {
	// The MFA token is a time-based one time password that is generated by a hardware or software token.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-mfa-totp.html
	//
	// Experimental.
	Otp *bool `json:"otp"`
	// The MFA token is sent to the user via SMS to their verified phone numbers.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-mfa-sms-text-message.html
	//
	// Experimental.
	Sms *bool `json:"sms"`
}

// The Number custom attribute type.
// Experimental.
type NumberAttribute interface {
	ICustomAttribute
	Bind() *CustomAttributeConfig
}

// The jsii proxy struct for NumberAttribute
type jsiiProxy_NumberAttribute struct {
	jsiiProxy_ICustomAttribute
}

// Experimental.
func NewNumberAttribute(props *NumberAttributeProps) NumberAttribute {
	_init_.Initialize()

	j := jsiiProxy_NumberAttribute{}

	_jsii_.Create(
		"monocdk.aws_cognito.NumberAttribute",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewNumberAttribute_Override(n NumberAttribute, props *NumberAttributeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.NumberAttribute",
		[]interface{}{props},
		n,
	)
}

// Bind this custom attribute type to the values as expected by CloudFormation.
// Experimental.
func (n *jsiiProxy_NumberAttribute) Bind() *CustomAttributeConfig {
	var returns *CustomAttributeConfig

	_jsii_.Invoke(
		n,
		"bind",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Constraints that can be applied to a custom attribute of number type.
// Experimental.
type NumberAttributeConstraints struct {
	// Maximum value of this attribute.
	// Experimental.
	Max *float64 `json:"max"`
	// Minimum value of this attribute.
	// Experimental.
	Min *float64 `json:"min"`
}

// Props for NumberAttr.
// Experimental.
type NumberAttributeProps struct {
	// Maximum value of this attribute.
	// Experimental.
	Max *float64 `json:"max"`
	// Minimum value of this attribute.
	// Experimental.
	Min *float64 `json:"min"`
	// Specifies whether the value of the attribute can be changed.
	//
	// For any user pool attribute that's mapped to an identity provider attribute, you must set this parameter to true.
	// Amazon Cognito updates mapped attributes when users sign in to your application through an identity provider.
	// If an attribute is immutable, Amazon Cognito throws an error when it attempts to update the attribute.
	// Experimental.
	Mutable *bool `json:"mutable"`
}

// Types of OAuth grant flows.
// See: - the 'Allowed OAuth Flows' section at https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-idp-settings.html
//
// Experimental.
type OAuthFlows struct {
	// Initiate an authorization code grant flow, which provides an authorization code as the response.
	// Experimental.
	AuthorizationCodeGrant *bool `json:"authorizationCodeGrant"`
	// Client should get the access token and ID token from the token endpoint using a combination of client and client_secret.
	// Experimental.
	ClientCredentials *bool `json:"clientCredentials"`
	// The client should get the access token and ID token directly.
	// Experimental.
	ImplicitCodeGrant *bool `json:"implicitCodeGrant"`
}

// OAuth scopes that are allowed with this client.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-idp-settings.html
//
// Experimental.
type OAuthScope interface {
	ScopeName() *string
}

// The jsii proxy struct for OAuthScope
type jsiiProxy_OAuthScope struct {
	_ byte // padding
}

func (j *jsiiProxy_OAuthScope) ScopeName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scopeName",
		&returns,
	)
	return returns
}


// Custom scope is one that you define for your own resource server in the Resource Servers.
//
// The format is 'resource-server-identifier/scope'.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-define-resource-servers.html
//
// Experimental.
func OAuthScope_Custom(name *string) OAuthScope {
	_init_.Initialize()

	var returns OAuthScope

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.OAuthScope",
		"custom",
		[]interface{}{name},
		&returns,
	)

	return returns
}

// Adds a custom scope that's tied to a resource server in your stack.
// Experimental.
func OAuthScope_ResourceServer(server IUserPoolResourceServer, scope ResourceServerScope) OAuthScope {
	_init_.Initialize()

	var returns OAuthScope

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.OAuthScope",
		"resourceServer",
		[]interface{}{server, scope},
		&returns,
	)

	return returns
}

func OAuthScope_COGNITO_ADMIN() OAuthScope {
	_init_.Initialize()
	var returns OAuthScope
	_jsii_.StaticGet(
		"monocdk.aws_cognito.OAuthScope",
		"COGNITO_ADMIN",
		&returns,
	)
	return returns
}

func OAuthScope_EMAIL() OAuthScope {
	_init_.Initialize()
	var returns OAuthScope
	_jsii_.StaticGet(
		"monocdk.aws_cognito.OAuthScope",
		"EMAIL",
		&returns,
	)
	return returns
}

func OAuthScope_OPENID() OAuthScope {
	_init_.Initialize()
	var returns OAuthScope
	_jsii_.StaticGet(
		"monocdk.aws_cognito.OAuthScope",
		"OPENID",
		&returns,
	)
	return returns
}

func OAuthScope_PHONE() OAuthScope {
	_init_.Initialize()
	var returns OAuthScope
	_jsii_.StaticGet(
		"monocdk.aws_cognito.OAuthScope",
		"PHONE",
		&returns,
	)
	return returns
}

func OAuthScope_PROFILE() OAuthScope {
	_init_.Initialize()
	var returns OAuthScope
	_jsii_.StaticGet(
		"monocdk.aws_cognito.OAuthScope",
		"PROFILE",
		&returns,
	)
	return returns
}

// OAuth settings to configure the interaction between the app and this client.
// Experimental.
type OAuthSettings struct {
	// List of allowed redirect URLs for the identity providers.
	// Experimental.
	CallbackUrls *[]*string `json:"callbackUrls"`
	// OAuth flows that are allowed with this client.
	// See: - the 'Allowed OAuth Flows' section at https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-idp-settings.html
	//
	// Experimental.
	Flows *OAuthFlows `json:"flows"`
	// List of allowed logout URLs for the identity providers.
	// Experimental.
	LogoutUrls *[]*string `json:"logoutUrls"`
	// OAuth scopes that are allowed with this client.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-app-idp-settings.html
	//
	// Experimental.
	Scopes *[]OAuthScope `json:"scopes"`
}

// Password policy for User Pools.
// Experimental.
type PasswordPolicy struct {
	// Minimum length required for a user's password.
	// Experimental.
	MinLength *float64 `json:"minLength"`
	// Whether the user is required to have digits in their password.
	// Experimental.
	RequireDigits *bool `json:"requireDigits"`
	// Whether the user is required to have lowercase characters in their password.
	// Experimental.
	RequireLowercase *bool `json:"requireLowercase"`
	// Whether the user is required to have symbols in their password.
	// Experimental.
	RequireSymbols *bool `json:"requireSymbols"`
	// Whether the user is required to have uppercase characters in their password.
	// Experimental.
	RequireUppercase *bool `json:"requireUppercase"`
	// The length of time the temporary password generated by an admin is valid.
	//
	// This must be provided as whole days, like Duration.days(3) or Duration.hours(48).
	// Fractional days, such as Duration.hours(20), will generate an error.
	// Experimental.
	TempPasswordValidity awscdk.Duration `json:"tempPasswordValidity"`
}

// An attribute available from a third party identity provider.
// Experimental.
type ProviderAttribute interface {
	AttributeName() *string
}

// The jsii proxy struct for ProviderAttribute
type jsiiProxy_ProviderAttribute struct {
	_ byte // padding
}

func (j *jsiiProxy_ProviderAttribute) AttributeName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attributeName",
		&returns,
	)
	return returns
}


// Use this to specify an attribute from the identity provider that is not pre-defined in the CDK.
// Experimental.
func ProviderAttribute_Other(attributeName *string) ProviderAttribute {
	_init_.Initialize()

	var returns ProviderAttribute

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.ProviderAttribute",
		"other",
		[]interface{}{attributeName},
		&returns,
	)

	return returns
}

func ProviderAttribute_AMAZON_EMAIL() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"AMAZON_EMAIL",
		&returns,
	)
	return returns
}

func ProviderAttribute_AMAZON_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"AMAZON_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_AMAZON_POSTAL_CODE() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"AMAZON_POSTAL_CODE",
		&returns,
	)
	return returns
}

func ProviderAttribute_AMAZON_USER_ID() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"AMAZON_USER_ID",
		&returns,
	)
	return returns
}

func ProviderAttribute_APPLE_EMAIL() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"APPLE_EMAIL",
		&returns,
	)
	return returns
}

func ProviderAttribute_APPLE_FIRST_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"APPLE_FIRST_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_APPLE_LAST_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"APPLE_LAST_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_APPLE_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"APPLE_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_BIRTHDAY() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_BIRTHDAY",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_EMAIL() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_EMAIL",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_FIRST_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_FIRST_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_GENDER() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_GENDER",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_ID() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_ID",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_LAST_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_LAST_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_LOCALE() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_LOCALE",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_MIDDLE_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_MIDDLE_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_FACEBOOK_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"FACEBOOK_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_BIRTHDAYS() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_BIRTHDAYS",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_EMAIL() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_EMAIL",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_FAMILY_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_FAMILY_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_GENDER() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_GENDER",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_GIVEN_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_GIVEN_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_NAME() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_NAME",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_NAMES() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_NAMES",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_PHONE_NUMBERS() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_PHONE_NUMBERS",
		&returns,
	)
	return returns
}

func ProviderAttribute_GOOGLE_PICTURE() ProviderAttribute {
	_init_.Initialize()
	var returns ProviderAttribute
	_jsii_.StaticGet(
		"monocdk.aws_cognito.ProviderAttribute",
		"GOOGLE_PICTURE",
		&returns,
	)
	return returns
}

// A scope for ResourceServer.
// Experimental.
type ResourceServerScope interface {
	ScopeDescription() *string
	ScopeName() *string
}

// The jsii proxy struct for ResourceServerScope
type jsiiProxy_ResourceServerScope struct {
	_ byte // padding
}

func (j *jsiiProxy_ResourceServerScope) ScopeDescription() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scopeDescription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceServerScope) ScopeName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scopeName",
		&returns,
	)
	return returns
}


// Experimental.
func NewResourceServerScope(props *ResourceServerScopeProps) ResourceServerScope {
	_init_.Initialize()

	j := jsiiProxy_ResourceServerScope{}

	_jsii_.Create(
		"monocdk.aws_cognito.ResourceServerScope",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewResourceServerScope_Override(r ResourceServerScope, props *ResourceServerScopeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.ResourceServerScope",
		[]interface{}{props},
		r,
	)
}

// Props to initialize ResourceServerScope.
// Experimental.
type ResourceServerScopeProps struct {
	// A description of the scope.
	// Experimental.
	ScopeDescription *string `json:"scopeDescription"`
	// The name of the scope.
	// Experimental.
	ScopeName *string `json:"scopeName"`
}

// The different ways in which users of this pool can sign up or sign in.
// Experimental.
type SignInAliases struct {
	// Whether a user is allowed to sign up or sign in with an email address.
	// Experimental.
	Email *bool `json:"email"`
	// Whether a user is allowed to sign up or sign in with a phone number.
	// Experimental.
	Phone *bool `json:"phone"`
	// Whether a user is allowed to sign in with a secondary username, that can be set and modified after sign up.
	//
	// Can only be used in conjunction with `USERNAME`.
	// Experimental.
	PreferredUsername *bool `json:"preferredUsername"`
	// Whether user is allowed to sign up or sign in with a username.
	// Experimental.
	Username *bool `json:"username"`
}

// Options to customize the behaviour of `signInUrl()`.
// Experimental.
type SignInUrlOptions struct {
	// Where to redirect to after sign in.
	// Experimental.
	RedirectUri *string `json:"redirectUri"`
	// The path in the URI where the sign-in page is located.
	// Experimental.
	SignInPath *string `json:"signInPath"`
}

// Standard attribute that can be marked as required or mutable.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html#cognito-user-pools-standard-attributes
//
// Experimental.
type StandardAttribute struct {
	// Specifies whether the value of the attribute can be changed.
	//
	// For any user pool attribute that's mapped to an identity provider attribute, this must be set to `true`.
	// Amazon Cognito updates mapped attributes when users sign in to your application through an identity provider.
	// If an attribute is immutable, Amazon Cognito throws an error when it attempts to update the attribute.
	// Experimental.
	Mutable *bool `json:"mutable"`
	// Specifies whether the attribute is required upon user registration.
	//
	// If the attribute is required and the user does not provide a value, registration or sign-in will fail.
	// Experimental.
	Required *bool `json:"required"`
}

// The set of standard attributes that can be marked as required or mutable.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html#cognito-user-pools-standard-attributes
//
// Experimental.
type StandardAttributes struct {
	// The user's postal address.
	// Experimental.
	Address *StandardAttribute `json:"address"`
	// The user's birthday, represented as an ISO 8601:2004 format.
	// Experimental.
	Birthdate *StandardAttribute `json:"birthdate"`
	// The user's e-mail address, represented as an RFC 5322 [RFC5322] addr-spec.
	// Experimental.
	Email *StandardAttribute `json:"email"`
	// DEPRECATED.
	// Deprecated: this is not a standard attribute and was incorrectly added to the CDK.
	// It is a Cognito built-in attribute and cannot be controlled as part of user pool creation.
	EmailVerified *StandardAttribute `json:"emailVerified"`
	// The surname or last name of the user.
	// Experimental.
	FamilyName *StandardAttribute `json:"familyName"`
	// The user's full name in displayable form, including all name parts, titles and suffixes.
	// Experimental.
	Fullname *StandardAttribute `json:"fullname"`
	// The user's gender.
	// Experimental.
	Gender *StandardAttribute `json:"gender"`
	// The user's first name or give name.
	// Experimental.
	GivenName *StandardAttribute `json:"givenName"`
	// The time, the user's information was last updated.
	// Experimental.
	LastUpdateTime *StandardAttribute `json:"lastUpdateTime"`
	// The user's locale, represented as a BCP47 [RFC5646] language tag.
	// Experimental.
	Locale *StandardAttribute `json:"locale"`
	// The user's middle name.
	// Experimental.
	MiddleName *StandardAttribute `json:"middleName"`
	// The user's nickname or casual name.
	// Experimental.
	Nickname *StandardAttribute `json:"nickname"`
	// The user's telephone number.
	// Experimental.
	PhoneNumber *StandardAttribute `json:"phoneNumber"`
	// DEPRECATED.
	// Deprecated: this is not a standard attribute and was incorrectly added to the CDK.
	// It is a Cognito built-in attribute and cannot be controlled as part of user pool creation.
	PhoneNumberVerified *StandardAttribute `json:"phoneNumberVerified"`
	// The user's preffered username, different from the immutable user name.
	// Experimental.
	PreferredUsername *StandardAttribute `json:"preferredUsername"`
	// The URL to the user's profile page.
	// Experimental.
	ProfilePage *StandardAttribute `json:"profilePage"`
	// The URL to the user's profile picture.
	// Experimental.
	ProfilePicture *StandardAttribute `json:"profilePicture"`
	// The user's time zone.
	// Experimental.
	Timezone *StandardAttribute `json:"timezone"`
	// The URL to the user's web page or blog.
	// Experimental.
	Website *StandardAttribute `json:"website"`
}

// This interface contains standard attributes recognized by Cognito from https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html including built-in attributes `email_verified` and `phone_number_verified`.
// Experimental.
type StandardAttributesMask struct {
	// The user's postal address.
	// Experimental.
	Address *bool `json:"address"`
	// The user's birthday, represented as an ISO 8601:2004 format.
	// Experimental.
	Birthdate *bool `json:"birthdate"`
	// The user's e-mail address, represented as an RFC 5322 [RFC5322] addr-spec.
	// Experimental.
	Email *bool `json:"email"`
	// Whether the email address has been verified.
	// Experimental.
	EmailVerified *bool `json:"emailVerified"`
	// The surname or last name of the user.
	// Experimental.
	FamilyName *bool `json:"familyName"`
	// The user's full name in displayable form, including all name parts, titles and suffixes.
	// Experimental.
	Fullname *bool `json:"fullname"`
	// The user's gender.
	// Experimental.
	Gender *bool `json:"gender"`
	// The user's first name or give name.
	// Experimental.
	GivenName *bool `json:"givenName"`
	// The time, the user's information was last updated.
	// Experimental.
	LastUpdateTime *bool `json:"lastUpdateTime"`
	// The user's locale, represented as a BCP47 [RFC5646] language tag.
	// Experimental.
	Locale *bool `json:"locale"`
	// The user's middle name.
	// Experimental.
	MiddleName *bool `json:"middleName"`
	// The user's nickname or casual name.
	// Experimental.
	Nickname *bool `json:"nickname"`
	// The user's telephone number.
	// Experimental.
	PhoneNumber *bool `json:"phoneNumber"`
	// Whether the phone number has been verified.
	// Experimental.
	PhoneNumberVerified *bool `json:"phoneNumberVerified"`
	// The user's preffered username, different from the immutable user name.
	// Experimental.
	PreferredUsername *bool `json:"preferredUsername"`
	// The URL to the user's profile page.
	// Experimental.
	ProfilePage *bool `json:"profilePage"`
	// The URL to the user's profile picture.
	// Experimental.
	ProfilePicture *bool `json:"profilePicture"`
	// The user's time zone.
	// Experimental.
	Timezone *bool `json:"timezone"`
	// The URL to the user's web page or blog.
	// Experimental.
	Website *bool `json:"website"`
}

// The String custom attribute type.
// Experimental.
type StringAttribute interface {
	ICustomAttribute
	Bind() *CustomAttributeConfig
}

// The jsii proxy struct for StringAttribute
type jsiiProxy_StringAttribute struct {
	jsiiProxy_ICustomAttribute
}

// Experimental.
func NewStringAttribute(props *StringAttributeProps) StringAttribute {
	_init_.Initialize()

	j := jsiiProxy_StringAttribute{}

	_jsii_.Create(
		"monocdk.aws_cognito.StringAttribute",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewStringAttribute_Override(s StringAttribute, props *StringAttributeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.StringAttribute",
		[]interface{}{props},
		s,
	)
}

// Bind this custom attribute type to the values as expected by CloudFormation.
// Experimental.
func (s *jsiiProxy_StringAttribute) Bind() *CustomAttributeConfig {
	var returns *CustomAttributeConfig

	_jsii_.Invoke(
		s,
		"bind",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Constraints that can be applied to a custom attribute of string type.
// Experimental.
type StringAttributeConstraints struct {
	// Maximum length of this attribute.
	// Experimental.
	MaxLen *float64 `json:"maxLen"`
	// Minimum length of this attribute.
	// Experimental.
	MinLen *float64 `json:"minLen"`
}

// Props for constructing a StringAttr.
// Experimental.
type StringAttributeProps struct {
	// Maximum length of this attribute.
	// Experimental.
	MaxLen *float64 `json:"maxLen"`
	// Minimum length of this attribute.
	// Experimental.
	MinLen *float64 `json:"minLen"`
	// Specifies whether the value of the attribute can be changed.
	//
	// For any user pool attribute that's mapped to an identity provider attribute, you must set this parameter to true.
	// Amazon Cognito updates mapped attributes when users sign in to your application through an identity provider.
	// If an attribute is immutable, Amazon Cognito throws an error when it attempts to update the attribute.
	// Experimental.
	Mutable *bool `json:"mutable"`
}

// User pool configuration when administrators sign users up.
// Experimental.
type UserInvitationConfig struct {
	// The template to the email body that is sent to the user when an administrator signs them up to the user pool.
	// Experimental.
	EmailBody *string `json:"emailBody"`
	// The template to the email subject that is sent to the user when an administrator signs them up to the user pool.
	// Experimental.
	EmailSubject *string `json:"emailSubject"`
	// The template to the SMS message that is sent to the user when an administrator signs them up to the user pool.
	// Experimental.
	SmsMessage *string `json:"smsMessage"`
}

// Define a Cognito User Pool.
// Experimental.
type UserPool interface {
	awscdk.Resource
	IUserPool
	Env() *awscdk.ResourceEnvironment
	IdentityProviders() *[]IUserPoolIdentityProvider
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	UserPoolArn() *string
	UserPoolId() *string
	UserPoolProviderName() *string
	UserPoolProviderUrl() *string
	AddClient(id *string, options *UserPoolClientOptions) UserPoolClient
	AddDomain(id *string, options *UserPoolDomainOptions) UserPoolDomain
	AddResourceServer(id *string, options *UserPoolResourceServerOptions) UserPoolResourceServer
	AddTrigger(operation UserPoolOperation, fn awslambda.IFunction)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RegisterIdentityProvider(provider IUserPoolIdentityProvider)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for UserPool
type jsiiProxy_UserPool struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPool
}

func (j *jsiiProxy_UserPool) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) IdentityProviders() *[]IUserPoolIdentityProvider {
	var returns *[]IUserPoolIdentityProvider
	_jsii_.Get(
		j,
		"identityProviders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) UserPoolArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) UserPoolId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) UserPoolProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolProviderName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPool) UserPoolProviderUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolProviderUrl",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPool(scope constructs.Construct, id *string, props *UserPoolProps) UserPool {
	_init_.Initialize()

	j := jsiiProxy_UserPool{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPool",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPool_Override(u UserPool, scope constructs.Construct, id *string, props *UserPoolProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPool",
		[]interface{}{scope, id, props},
		u,
	)
}

// Import an existing user pool based on its ARN.
// Experimental.
func UserPool_FromUserPoolArn(scope constructs.Construct, id *string, userPoolArn *string) IUserPool {
	_init_.Initialize()

	var returns IUserPool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPool",
		"fromUserPoolArn",
		[]interface{}{scope, id, userPoolArn},
		&returns,
	)

	return returns
}

// Import an existing user pool based on its id.
// Experimental.
func UserPool_FromUserPoolId(scope constructs.Construct, id *string, userPoolId *string) IUserPool {
	_init_.Initialize()

	var returns IUserPool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPool",
		"fromUserPoolId",
		[]interface{}{scope, id, userPoolId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func UserPool_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPool",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPool_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPool",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a new app client to this user pool.
// Experimental.
func (u *jsiiProxy_UserPool) AddClient(id *string, options *UserPoolClientOptions) UserPoolClient {
	var returns UserPoolClient

	_jsii_.Invoke(
		u,
		"addClient",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Associate a domain to this user pool.
// Experimental.
func (u *jsiiProxy_UserPool) AddDomain(id *string, options *UserPoolDomainOptions) UserPoolDomain {
	var returns UserPoolDomain

	_jsii_.Invoke(
		u,
		"addDomain",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Add a new resource server to this user pool.
// Experimental.
func (u *jsiiProxy_UserPool) AddResourceServer(id *string, options *UserPoolResourceServerOptions) UserPoolResourceServer {
	var returns UserPoolResourceServer

	_jsii_.Invoke(
		u,
		"addResourceServer",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Add a lambda trigger to a user pool operation.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html
//
// Experimental.
func (u *jsiiProxy_UserPool) AddTrigger(operation UserPoolOperation, fn awslambda.IFunction) {
	_jsii_.InvokeVoid(
		u,
		"addTrigger",
		[]interface{}{operation, fn},
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
func (u *jsiiProxy_UserPool) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UserPool) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPool) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPool) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPool) OnPrepare() {
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
func (u *jsiiProxy_UserPool) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPool) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPool) Prepare() {
	_jsii_.InvokeVoid(
		u,
		"prepare",
		nil, // no parameters
	)
}

// Register an identity provider with this user pool.
// Experimental.
func (u *jsiiProxy_UserPool) RegisterIdentityProvider(provider IUserPoolIdentityProvider) {
	_jsii_.InvokeVoid(
		u,
		"registerIdentityProvider",
		[]interface{}{provider},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (u *jsiiProxy_UserPool) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPool) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPool) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Define a UserPool App Client.
// Experimental.
type UserPoolClient interface {
	awscdk.Resource
	IUserPoolClient
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	OAuthFlows() *OAuthFlows
	PhysicalName() *string
	Stack() awscdk.Stack
	UserPoolClientId() *string
	UserPoolClientName() *string
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

// The jsii proxy struct for UserPoolClient
type jsiiProxy_UserPoolClient struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPoolClient
}

func (j *jsiiProxy_UserPoolClient) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolClient) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolClient) OAuthFlows() *OAuthFlows {
	var returns *OAuthFlows
	_jsii_.Get(
		j,
		"oAuthFlows",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolClient) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolClient) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolClient) UserPoolClientId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolClientId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolClient) UserPoolClientName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolClientName",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPoolClient(scope constructs.Construct, id *string, props *UserPoolClientProps) UserPoolClient {
	_init_.Initialize()

	j := jsiiProxy_UserPoolClient{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolClient",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPoolClient_Override(u UserPoolClient, scope constructs.Construct, id *string, props *UserPoolClientProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolClient",
		[]interface{}{scope, id, props},
		u,
	)
}

// Import a user pool client given its id.
// Experimental.
func UserPoolClient_FromUserPoolClientId(scope constructs.Construct, id *string, userPoolClientId *string) IUserPoolClient {
	_init_.Initialize()

	var returns IUserPoolClient

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolClient",
		"fromUserPoolClientId",
		[]interface{}{scope, id, userPoolClientId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func UserPoolClient_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolClient",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPoolClient_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolClient",
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
func (u *jsiiProxy_UserPoolClient) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UserPoolClient) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPoolClient) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPoolClient) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPoolClient) OnPrepare() {
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
func (u *jsiiProxy_UserPoolClient) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPoolClient) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPoolClient) Prepare() {
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
func (u *jsiiProxy_UserPoolClient) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPoolClient) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPoolClient) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Identity providers supported by the UserPoolClient.
// Experimental.
type UserPoolClientIdentityProvider interface {
	Name() *string
}

// The jsii proxy struct for UserPoolClientIdentityProvider
type jsiiProxy_UserPoolClientIdentityProvider struct {
	_ byte // padding
}

func (j *jsiiProxy_UserPoolClientIdentityProvider) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}


// Specify a provider not yet supported by the CDK.
// Experimental.
func UserPoolClientIdentityProvider_Custom(name *string) UserPoolClientIdentityProvider {
	_init_.Initialize()

	var returns UserPoolClientIdentityProvider

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolClientIdentityProvider",
		"custom",
		[]interface{}{name},
		&returns,
	)

	return returns
}

func UserPoolClientIdentityProvider_AMAZON() UserPoolClientIdentityProvider {
	_init_.Initialize()
	var returns UserPoolClientIdentityProvider
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolClientIdentityProvider",
		"AMAZON",
		&returns,
	)
	return returns
}

func UserPoolClientIdentityProvider_APPLE() UserPoolClientIdentityProvider {
	_init_.Initialize()
	var returns UserPoolClientIdentityProvider
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolClientIdentityProvider",
		"APPLE",
		&returns,
	)
	return returns
}

func UserPoolClientIdentityProvider_COGNITO() UserPoolClientIdentityProvider {
	_init_.Initialize()
	var returns UserPoolClientIdentityProvider
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolClientIdentityProvider",
		"COGNITO",
		&returns,
	)
	return returns
}

func UserPoolClientIdentityProvider_FACEBOOK() UserPoolClientIdentityProvider {
	_init_.Initialize()
	var returns UserPoolClientIdentityProvider
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolClientIdentityProvider",
		"FACEBOOK",
		&returns,
	)
	return returns
}

func UserPoolClientIdentityProvider_GOOGLE() UserPoolClientIdentityProvider {
	_init_.Initialize()
	var returns UserPoolClientIdentityProvider
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolClientIdentityProvider",
		"GOOGLE",
		&returns,
	)
	return returns
}

// Options to create a UserPoolClient.
// Experimental.
type UserPoolClientOptions struct {
	// Validity of the access token.
	//
	// Values between 5 minutes and 1 day are valid. The duration can not be longer than the refresh token validity.
	// See: https://docs.aws.amazon.com/en_us/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html#amazon-cognito-user-pools-using-the-access-token
	//
	// Experimental.
	AccessTokenValidity awscdk.Duration `json:"accessTokenValidity"`
	// The set of OAuth authentication flows to enable on the client.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow.html
	//
	// Experimental.
	AuthFlows *AuthFlow `json:"authFlows"`
	// Turns off all OAuth interactions for this client.
	// Experimental.
	DisableOAuth *bool `json:"disableOAuth"`
	// Whether to generate a client secret.
	// Experimental.
	GenerateSecret *bool `json:"generateSecret"`
	// Validity of the ID token.
	//
	// Values between 5 minutes and 1 day are valid. The duration can not be longer than the refresh token validity.
	// See: https://docs.aws.amazon.com/en_us/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html#amazon-cognito-user-pools-using-the-id-token
	//
	// Experimental.
	IdTokenValidity awscdk.Duration `json:"idTokenValidity"`
	// OAuth settings for this client to interact with the app.
	//
	// An error is thrown when this is specified and `disableOAuth` is set.
	// Experimental.
	OAuth *OAuthSettings `json:"oAuth"`
	// Whether Cognito returns a UserNotFoundException exception when the user does not exist in the user pool (false), or whether it returns another type of error that doesn't reveal the user's absence.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pool-managing-errors.html
	//
	// Experimental.
	PreventUserExistenceErrors *bool `json:"preventUserExistenceErrors"`
	// The set of attributes this client will be able to read.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html#user-pool-settings-attribute-permissions-and-scopes
	//
	// Experimental.
	ReadAttributes ClientAttributes `json:"readAttributes"`
	// Validity of the refresh token.
	//
	// Values between 60 minutes and 10 years are valid.
	// See: https://docs.aws.amazon.com/en_us/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html#amazon-cognito-user-pools-using-the-refresh-token
	//
	// Experimental.
	RefreshTokenValidity awscdk.Duration `json:"refreshTokenValidity"`
	// The list of identity providers that users should be able to use to sign in using this client.
	// Experimental.
	SupportedIdentityProviders *[]UserPoolClientIdentityProvider `json:"supportedIdentityProviders"`
	// Name of the application client.
	// Experimental.
	UserPoolClientName *string `json:"userPoolClientName"`
	// The set of attributes this client will be able to write.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html#user-pool-settings-attribute-permissions-and-scopes
	//
	// Experimental.
	WriteAttributes ClientAttributes `json:"writeAttributes"`
}

// Properties for the UserPoolClient construct.
// Experimental.
type UserPoolClientProps struct {
	// Validity of the access token.
	//
	// Values between 5 minutes and 1 day are valid. The duration can not be longer than the refresh token validity.
	// See: https://docs.aws.amazon.com/en_us/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html#amazon-cognito-user-pools-using-the-access-token
	//
	// Experimental.
	AccessTokenValidity awscdk.Duration `json:"accessTokenValidity"`
	// The set of OAuth authentication flows to enable on the client.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow.html
	//
	// Experimental.
	AuthFlows *AuthFlow `json:"authFlows"`
	// Turns off all OAuth interactions for this client.
	// Experimental.
	DisableOAuth *bool `json:"disableOAuth"`
	// Whether to generate a client secret.
	// Experimental.
	GenerateSecret *bool `json:"generateSecret"`
	// Validity of the ID token.
	//
	// Values between 5 minutes and 1 day are valid. The duration can not be longer than the refresh token validity.
	// See: https://docs.aws.amazon.com/en_us/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html#amazon-cognito-user-pools-using-the-id-token
	//
	// Experimental.
	IdTokenValidity awscdk.Duration `json:"idTokenValidity"`
	// OAuth settings for this client to interact with the app.
	//
	// An error is thrown when this is specified and `disableOAuth` is set.
	// Experimental.
	OAuth *OAuthSettings `json:"oAuth"`
	// Whether Cognito returns a UserNotFoundException exception when the user does not exist in the user pool (false), or whether it returns another type of error that doesn't reveal the user's absence.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pool-managing-errors.html
	//
	// Experimental.
	PreventUserExistenceErrors *bool `json:"preventUserExistenceErrors"`
	// The set of attributes this client will be able to read.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html#user-pool-settings-attribute-permissions-and-scopes
	//
	// Experimental.
	ReadAttributes ClientAttributes `json:"readAttributes"`
	// Validity of the refresh token.
	//
	// Values between 60 minutes and 10 years are valid.
	// See: https://docs.aws.amazon.com/en_us/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html#amazon-cognito-user-pools-using-the-refresh-token
	//
	// Experimental.
	RefreshTokenValidity awscdk.Duration `json:"refreshTokenValidity"`
	// The list of identity providers that users should be able to use to sign in using this client.
	// Experimental.
	SupportedIdentityProviders *[]UserPoolClientIdentityProvider `json:"supportedIdentityProviders"`
	// Name of the application client.
	// Experimental.
	UserPoolClientName *string `json:"userPoolClientName"`
	// The set of attributes this client will be able to write.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html#user-pool-settings-attribute-permissions-and-scopes
	//
	// Experimental.
	WriteAttributes ClientAttributes `json:"writeAttributes"`
	// The UserPool resource this client will have access to.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
}

// Define a user pool domain.
// Experimental.
type UserPoolDomain interface {
	awscdk.Resource
	IUserPoolDomain
	CloudFrontDomainName() *string
	DomainName() *string
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	BaseUrl() *string
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	SignInUrl(client UserPoolClient, options *SignInUrlOptions) *string
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for UserPoolDomain
type jsiiProxy_UserPoolDomain struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPoolDomain
}

func (j *jsiiProxy_UserPoolDomain) CloudFrontDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cloudFrontDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolDomain) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolDomain) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolDomain) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolDomain) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolDomain) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPoolDomain(scope constructs.Construct, id *string, props *UserPoolDomainProps) UserPoolDomain {
	_init_.Initialize()

	j := jsiiProxy_UserPoolDomain{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolDomain",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPoolDomain_Override(u UserPoolDomain, scope constructs.Construct, id *string, props *UserPoolDomainProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolDomain",
		[]interface{}{scope, id, props},
		u,
	)
}

// Import a UserPoolDomain given its domain name.
// Experimental.
func UserPoolDomain_FromDomainName(scope constructs.Construct, id *string, userPoolDomainName *string) IUserPoolDomain {
	_init_.Initialize()

	var returns IUserPoolDomain

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolDomain",
		"fromDomainName",
		[]interface{}{scope, id, userPoolDomainName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func UserPoolDomain_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolDomain",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPoolDomain_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolDomain",
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
func (u *jsiiProxy_UserPoolDomain) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// The URL to the hosted UI associated with this domain.
// Experimental.
func (u *jsiiProxy_UserPoolDomain) BaseUrl() *string {
	var returns *string

	_jsii_.Invoke(
		u,
		"baseUrl",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (u *jsiiProxy_UserPoolDomain) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPoolDomain) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPoolDomain) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPoolDomain) OnPrepare() {
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
func (u *jsiiProxy_UserPoolDomain) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPoolDomain) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPoolDomain) Prepare() {
	_jsii_.InvokeVoid(
		u,
		"prepare",
		nil, // no parameters
	)
}

// The URL to the sign in page in this domain using a specific UserPoolClient.
// Experimental.
func (u *jsiiProxy_UserPoolDomain) SignInUrl(client UserPoolClient, options *SignInUrlOptions) *string {
	var returns *string

	_jsii_.Invoke(
		u,
		"signInUrl",
		[]interface{}{client, options},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (u *jsiiProxy_UserPoolDomain) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPoolDomain) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPoolDomain) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options to create a UserPoolDomain.
// Experimental.
type UserPoolDomainOptions struct {
	// Associate a cognito prefix domain with your user pool Either `customDomain` or `cognitoDomain` must be specified.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-assign-domain-prefix.html
	//
	// Experimental.
	CognitoDomain *CognitoDomainOptions `json:"cognitoDomain"`
	// Associate a custom domain with your user pool Either `customDomain` or `cognitoDomain` must be specified.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-add-custom-domain.html
	//
	// Experimental.
	CustomDomain *CustomDomainOptions `json:"customDomain"`
}

// Props for UserPoolDomain construct.
// Experimental.
type UserPoolDomainProps struct {
	// Associate a cognito prefix domain with your user pool Either `customDomain` or `cognitoDomain` must be specified.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-assign-domain-prefix.html
	//
	// Experimental.
	CognitoDomain *CognitoDomainOptions `json:"cognitoDomain"`
	// Associate a custom domain with your user pool Either `customDomain` or `cognitoDomain` must be specified.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-add-custom-domain.html
	//
	// Experimental.
	CustomDomain *CustomDomainOptions `json:"customDomain"`
	// The user pool to which this domain should be associated.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
}

// User pool third-party identity providers.
// Experimental.
type UserPoolIdentityProvider interface {
}

// The jsii proxy struct for UserPoolIdentityProvider
type jsiiProxy_UserPoolIdentityProvider struct {
	_ byte // padding
}

// Import an existing UserPoolIdentityProvider.
// Experimental.
func UserPoolIdentityProvider_FromProviderName(scope constructs.Construct, id *string, providerName *string) IUserPoolIdentityProvider {
	_init_.Initialize()

	var returns IUserPoolIdentityProvider

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProvider",
		"fromProviderName",
		[]interface{}{scope, id, providerName},
		&returns,
	)

	return returns
}

// Represents a identity provider that integrates with 'Login with Amazon'.
// Experimental.
type UserPoolIdentityProviderAmazon interface {
	awscdk.Resource
	IUserPoolIdentityProvider
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ProviderName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAttributeMapping() interface{}
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

// The jsii proxy struct for UserPoolIdentityProviderAmazon
type jsiiProxy_UserPoolIdentityProviderAmazon struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPoolIdentityProvider
}

func (j *jsiiProxy_UserPoolIdentityProviderAmazon) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderAmazon) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderAmazon) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderAmazon) ProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"providerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderAmazon) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPoolIdentityProviderAmazon(scope constructs.Construct, id *string, props *UserPoolIdentityProviderAmazonProps) UserPoolIdentityProviderAmazon {
	_init_.Initialize()

	j := jsiiProxy_UserPoolIdentityProviderAmazon{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderAmazon",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPoolIdentityProviderAmazon_Override(u UserPoolIdentityProviderAmazon, scope constructs.Construct, id *string, props *UserPoolIdentityProviderAmazonProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderAmazon",
		[]interface{}{scope, id, props},
		u,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func UserPoolIdentityProviderAmazon_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderAmazon",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPoolIdentityProviderAmazon_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderAmazon",
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) ConfigureAttributeMapping() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		u,
		"configureAttributeMapping",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) OnPrepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) Prepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderAmazon) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to initialize UserPoolAmazonIdentityProvider.
// Experimental.
type UserPoolIdentityProviderAmazonProps struct {
	// The user pool to which this construct provides identities.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
	// Mapping attributes from the identity provider to standard and custom attributes of the user pool.
	// Experimental.
	AttributeMapping *AttributeMapping `json:"attributeMapping"`
	// The client id recognized by 'Login with Amazon' APIs.
	// See: https://developer.amazon.com/docs/login-with-amazon/security-profile.html#client-identifier
	//
	// Experimental.
	ClientId *string `json:"clientId"`
	// The client secret to be accompanied with clientId for 'Login with Amazon' APIs to authenticate the client.
	// See: https://developer.amazon.com/docs/login-with-amazon/security-profile.html#client-identifier
	//
	// Experimental.
	ClientSecret *string `json:"clientSecret"`
	// The types of user profile data to obtain for the Amazon profile.
	// See: https://developer.amazon.com/docs/login-with-amazon/customer-profile.html
	//
	// Experimental.
	Scopes *[]*string `json:"scopes"`
}

// Represents a identity provider that integrates with 'Apple'.
// Experimental.
type UserPoolIdentityProviderApple interface {
	awscdk.Resource
	IUserPoolIdentityProvider
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ProviderName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAttributeMapping() interface{}
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

// The jsii proxy struct for UserPoolIdentityProviderApple
type jsiiProxy_UserPoolIdentityProviderApple struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPoolIdentityProvider
}

func (j *jsiiProxy_UserPoolIdentityProviderApple) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderApple) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderApple) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderApple) ProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"providerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderApple) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPoolIdentityProviderApple(scope constructs.Construct, id *string, props *UserPoolIdentityProviderAppleProps) UserPoolIdentityProviderApple {
	_init_.Initialize()

	j := jsiiProxy_UserPoolIdentityProviderApple{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderApple",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPoolIdentityProviderApple_Override(u UserPoolIdentityProviderApple, scope constructs.Construct, id *string, props *UserPoolIdentityProviderAppleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderApple",
		[]interface{}{scope, id, props},
		u,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func UserPoolIdentityProviderApple_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderApple",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPoolIdentityProviderApple_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderApple",
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderApple) ConfigureAttributeMapping() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		u,
		"configureAttributeMapping",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderApple) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) OnPrepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) Prepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderApple) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderApple) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderApple) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to initialize UserPoolAppleIdentityProvider.
// Experimental.
type UserPoolIdentityProviderAppleProps struct {
	// The user pool to which this construct provides identities.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
	// Mapping attributes from the identity provider to standard and custom attributes of the user pool.
	// Experimental.
	AttributeMapping *AttributeMapping `json:"attributeMapping"`
	// The client id recognized by Apple APIs.
	// See: https://developer.apple.com/documentation/sign_in_with_apple/clientconfigi/3230948-clientid
	//
	// Experimental.
	ClientId *string `json:"clientId"`
	// The keyId (of the same key, which content has to be later supplied as `privateKey`) for Apple APIs to authenticate the client.
	// Experimental.
	KeyId *string `json:"keyId"`
	// The privateKey content for Apple APIs to authenticate the client.
	// Experimental.
	PrivateKey *string `json:"privateKey"`
	// The teamId for Apple APIs to authenticate the client.
	// Experimental.
	TeamId *string `json:"teamId"`
	// The list of apple permissions to obtain for getting access to the apple profile.
	// See: https://developer.apple.com/documentation/sign_in_with_apple/clientconfigi/3230955-scope
	//
	// Experimental.
	Scopes *[]*string `json:"scopes"`
}

// Represents a identity provider that integrates with 'Facebook Login'.
// Experimental.
type UserPoolIdentityProviderFacebook interface {
	awscdk.Resource
	IUserPoolIdentityProvider
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ProviderName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAttributeMapping() interface{}
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

// The jsii proxy struct for UserPoolIdentityProviderFacebook
type jsiiProxy_UserPoolIdentityProviderFacebook struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPoolIdentityProvider
}

func (j *jsiiProxy_UserPoolIdentityProviderFacebook) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderFacebook) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderFacebook) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderFacebook) ProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"providerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderFacebook) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPoolIdentityProviderFacebook(scope constructs.Construct, id *string, props *UserPoolIdentityProviderFacebookProps) UserPoolIdentityProviderFacebook {
	_init_.Initialize()

	j := jsiiProxy_UserPoolIdentityProviderFacebook{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderFacebook",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPoolIdentityProviderFacebook_Override(u UserPoolIdentityProviderFacebook, scope constructs.Construct, id *string, props *UserPoolIdentityProviderFacebookProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderFacebook",
		[]interface{}{scope, id, props},
		u,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func UserPoolIdentityProviderFacebook_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderFacebook",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPoolIdentityProviderFacebook_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderFacebook",
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) ConfigureAttributeMapping() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		u,
		"configureAttributeMapping",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) OnPrepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) Prepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderFacebook) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to initialize UserPoolFacebookIdentityProvider.
// Experimental.
type UserPoolIdentityProviderFacebookProps struct {
	// The user pool to which this construct provides identities.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
	// Mapping attributes from the identity provider to standard and custom attributes of the user pool.
	// Experimental.
	AttributeMapping *AttributeMapping `json:"attributeMapping"`
	// The client id recognized by Facebook APIs.
	// Experimental.
	ClientId *string `json:"clientId"`
	// The client secret to be accompanied with clientUd for Facebook to authenticate the client.
	// See: https://developers.facebook.com/docs/facebook-login/security#appsecret
	//
	// Experimental.
	ClientSecret *string `json:"clientSecret"`
	// The Facebook API version to use.
	// Experimental.
	ApiVersion *string `json:"apiVersion"`
	// The list of facebook permissions to obtain for getting access to the Facebook profile.
	// See: https://developers.facebook.com/docs/facebook-login/permissions
	//
	// Experimental.
	Scopes *[]*string `json:"scopes"`
}

// Represents a identity provider that integrates with 'Google'.
// Experimental.
type UserPoolIdentityProviderGoogle interface {
	awscdk.Resource
	IUserPoolIdentityProvider
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	ProviderName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ConfigureAttributeMapping() interface{}
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

// The jsii proxy struct for UserPoolIdentityProviderGoogle
type jsiiProxy_UserPoolIdentityProviderGoogle struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPoolIdentityProvider
}

func (j *jsiiProxy_UserPoolIdentityProviderGoogle) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderGoogle) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderGoogle) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderGoogle) ProviderName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"providerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolIdentityProviderGoogle) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPoolIdentityProviderGoogle(scope constructs.Construct, id *string, props *UserPoolIdentityProviderGoogleProps) UserPoolIdentityProviderGoogle {
	_init_.Initialize()

	j := jsiiProxy_UserPoolIdentityProviderGoogle{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderGoogle",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPoolIdentityProviderGoogle_Override(u UserPoolIdentityProviderGoogle, scope constructs.Construct, id *string, props *UserPoolIdentityProviderGoogleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolIdentityProviderGoogle",
		[]interface{}{scope, id, props},
		u,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func UserPoolIdentityProviderGoogle_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderGoogle",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPoolIdentityProviderGoogle_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolIdentityProviderGoogle",
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) ConfigureAttributeMapping() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		u,
		"configureAttributeMapping",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) OnPrepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) Prepare() {
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
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPoolIdentityProviderGoogle) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to initialize UserPoolGoogleIdentityProvider.
// Experimental.
type UserPoolIdentityProviderGoogleProps struct {
	// The user pool to which this construct provides identities.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
	// Mapping attributes from the identity provider to standard and custom attributes of the user pool.
	// Experimental.
	AttributeMapping *AttributeMapping `json:"attributeMapping"`
	// The client id recognized by Google APIs.
	// See: https://developers.google.com/identity/sign-in/web/sign-in#specify_your_apps_client_id
	//
	// Experimental.
	ClientId *string `json:"clientId"`
	// The client secret to be accompanied with clientId for Google APIs to authenticate the client.
	// See: https://developers.google.com/identity/sign-in/web/sign-in
	//
	// Experimental.
	ClientSecret *string `json:"clientSecret"`
	// The list of google permissions to obtain for getting access to the google profile.
	// See: https://developers.google.com/identity/sign-in/web/sign-in
	//
	// Experimental.
	Scopes *[]*string `json:"scopes"`
}

// Properties to create a new instance of UserPoolIdentityProvider.
// Experimental.
type UserPoolIdentityProviderProps struct {
	// The user pool to which this construct provides identities.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
	// Mapping attributes from the identity provider to standard and custom attributes of the user pool.
	// Experimental.
	AttributeMapping *AttributeMapping `json:"attributeMapping"`
}

// User pool operations to which lambda triggers can be attached.
// Experimental.
type UserPoolOperation interface {
	OperationName() *string
}

// The jsii proxy struct for UserPoolOperation
type jsiiProxy_UserPoolOperation struct {
	_ byte // padding
}

func (j *jsiiProxy_UserPoolOperation) OperationName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"operationName",
		&returns,
	)
	return returns
}


// A custom user pool operation.
// Experimental.
func UserPoolOperation_Of(name *string) UserPoolOperation {
	_init_.Initialize()

	var returns UserPoolOperation

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolOperation",
		"of",
		[]interface{}{name},
		&returns,
	)

	return returns
}

func UserPoolOperation_CREATE_AUTH_CHALLENGE() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"CREATE_AUTH_CHALLENGE",
		&returns,
	)
	return returns
}

func UserPoolOperation_CUSTOM_MESSAGE() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"CUSTOM_MESSAGE",
		&returns,
	)
	return returns
}

func UserPoolOperation_DEFINE_AUTH_CHALLENGE() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"DEFINE_AUTH_CHALLENGE",
		&returns,
	)
	return returns
}

func UserPoolOperation_POST_AUTHENTICATION() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"POST_AUTHENTICATION",
		&returns,
	)
	return returns
}

func UserPoolOperation_POST_CONFIRMATION() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"POST_CONFIRMATION",
		&returns,
	)
	return returns
}

func UserPoolOperation_PRE_AUTHENTICATION() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"PRE_AUTHENTICATION",
		&returns,
	)
	return returns
}

func UserPoolOperation_PRE_SIGN_UP() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"PRE_SIGN_UP",
		&returns,
	)
	return returns
}

func UserPoolOperation_PRE_TOKEN_GENERATION() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"PRE_TOKEN_GENERATION",
		&returns,
	)
	return returns
}

func UserPoolOperation_USER_MIGRATION() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"USER_MIGRATION",
		&returns,
	)
	return returns
}

func UserPoolOperation_VERIFY_AUTH_CHALLENGE_RESPONSE() UserPoolOperation {
	_init_.Initialize()
	var returns UserPoolOperation
	_jsii_.StaticGet(
		"monocdk.aws_cognito.UserPoolOperation",
		"VERIFY_AUTH_CHALLENGE_RESPONSE",
		&returns,
	)
	return returns
}

// Props for the UserPool construct.
// Experimental.
type UserPoolProps struct {
	// How will a user be able to recover their account?
	// Experimental.
	AccountRecovery AccountRecovery `json:"accountRecovery"`
	// Attributes which Cognito will look to verify automatically upon user sign up.
	//
	// EMAIL and PHONE are the only available options.
	// Experimental.
	AutoVerify *AutoVerifiedAttrs `json:"autoVerify"`
	// Define a set of custom attributes that can be configured for each user in the user pool.
	// Experimental.
	CustomAttributes *map[string]ICustomAttribute `json:"customAttributes"`
	// Email settings for a user pool.
	// Experimental.
	EmailSettings *EmailSettings `json:"emailSettings"`
	// Setting this would explicitly enable or disable SMS role creation.
	//
	// When left unspecified, CDK will determine based on other properties if a role is needed or not.
	// Experimental.
	EnableSmsRole *bool `json:"enableSmsRole"`
	// Lambda functions to use for supported Cognito triggers.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html
	//
	// Experimental.
	LambdaTriggers *UserPoolTriggers `json:"lambdaTriggers"`
	// Configure whether users of this user pool can or are required use MFA to sign in.
	// Experimental.
	Mfa Mfa `json:"mfa"`
	// The SMS message template sent during MFA verification.
	//
	// Use '{####}' in the template where Cognito should insert the verification code.
	// Experimental.
	MfaMessage *string `json:"mfaMessage"`
	// Configure the MFA types that users can use in this user pool.
	//
	// Ignored if `mfa` is set to `OFF`.
	// Experimental.
	MfaSecondFactor *MfaSecondFactor `json:"mfaSecondFactor"`
	// Password policy for this user pool.
	// Experimental.
	PasswordPolicy *PasswordPolicy `json:"passwordPolicy"`
	// Policy to apply when the user pool is removed from the stack.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// Whether self sign up should be enabled.
	//
	// This can be further configured via the `selfSignUp` property.
	// Experimental.
	SelfSignUpEnabled *bool `json:"selfSignUpEnabled"`
	// Methods in which a user registers or signs in to a user pool.
	//
	// Allows either username with aliases OR sign in with email, phone, or both.
	//
	// Read the sections on usernames and aliases to learn more -
	// https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html
	//
	// To match with 'Option 1' in the above link, with a verified email, this property should be set to
	// `{ username: true, email: true }`. To match with 'Option 2' in the above link with both a verified email and phone
	// number, this property should be set to `{ email: true, phone: true }`.
	// Experimental.
	SignInAliases *SignInAliases `json:"signInAliases"`
	// Whether sign-in aliases should be evaluated with case sensitivity.
	//
	// For example, when this option is set to false, users will be able to sign in using either `MyUsername` or `myusername`.
	// Experimental.
	SignInCaseSensitive *bool `json:"signInCaseSensitive"`
	// The IAM role that Cognito will assume while sending SMS messages.
	// Experimental.
	SmsRole awsiam.IRole `json:"smsRole"`
	// The 'ExternalId' that Cognito service must using when assuming the `smsRole`, if the role is restricted with an 'sts:ExternalId' conditional.
	//
	// Learn more about ExternalId here - https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html
	//
	// This property will be ignored if `smsRole` is not specified.
	// Experimental.
	SmsRoleExternalId *string `json:"smsRoleExternalId"`
	// The set of attributes that are required for every user in the user pool.
	//
	// Read more on attributes here - https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-attributes.html
	// Experimental.
	StandardAttributes *StandardAttributes `json:"standardAttributes"`
	// Configuration around admins signing up users into a user pool.
	// Experimental.
	UserInvitation *UserInvitationConfig `json:"userInvitation"`
	// Name of the user pool.
	// Experimental.
	UserPoolName *string `json:"userPoolName"`
	// Configuration around users signing themselves up to the user pool.
	//
	// Enable or disable self sign-up via the `selfSignUpEnabled` property.
	// Experimental.
	UserVerification *UserVerificationConfig `json:"userVerification"`
}

// Defines a User Pool OAuth2.0 Resource Server.
// Experimental.
type UserPoolResourceServer interface {
	awscdk.Resource
	IUserPoolResourceServer
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	UserPoolResourceServerId() *string
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

// The jsii proxy struct for UserPoolResourceServer
type jsiiProxy_UserPoolResourceServer struct {
	internal.Type__awscdkResource
	jsiiProxy_IUserPoolResourceServer
}

func (j *jsiiProxy_UserPoolResourceServer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolResourceServer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolResourceServer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolResourceServer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UserPoolResourceServer) UserPoolResourceServerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userPoolResourceServerId",
		&returns,
	)
	return returns
}


// Experimental.
func NewUserPoolResourceServer(scope constructs.Construct, id *string, props *UserPoolResourceServerProps) UserPoolResourceServer {
	_init_.Initialize()

	j := jsiiProxy_UserPoolResourceServer{}

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolResourceServer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUserPoolResourceServer_Override(u UserPoolResourceServer, scope constructs.Construct, id *string, props *UserPoolResourceServerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_cognito.UserPoolResourceServer",
		[]interface{}{scope, id, props},
		u,
	)
}

// Import a user pool resource client given its id.
// Experimental.
func UserPoolResourceServer_FromUserPoolResourceServerId(scope constructs.Construct, id *string, userPoolResourceServerId *string) IUserPoolResourceServer {
	_init_.Initialize()

	var returns IUserPoolResourceServer

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolResourceServer",
		"fromUserPoolResourceServerId",
		[]interface{}{scope, id, userPoolResourceServerId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func UserPoolResourceServer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolResourceServer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UserPoolResourceServer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_cognito.UserPoolResourceServer",
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
func (u *jsiiProxy_UserPoolResourceServer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UserPoolResourceServer) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UserPoolResourceServer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UserPoolResourceServer) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UserPoolResourceServer) OnPrepare() {
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
func (u *jsiiProxy_UserPoolResourceServer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UserPoolResourceServer) OnValidate() *[]*string {
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
func (u *jsiiProxy_UserPoolResourceServer) Prepare() {
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
func (u *jsiiProxy_UserPoolResourceServer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UserPoolResourceServer) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (u *jsiiProxy_UserPoolResourceServer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options to create a UserPoolResourceServer.
// Experimental.
type UserPoolResourceServerOptions struct {
	// A unique resource server identifier for the resource server.
	// Experimental.
	Identifier *string `json:"identifier"`
	// Oauth scopes.
	// Experimental.
	Scopes *[]ResourceServerScope `json:"scopes"`
	// A friendly name for the resource server.
	// Experimental.
	UserPoolResourceServerName *string `json:"userPoolResourceServerName"`
}

// Properties for the UserPoolResourceServer construct.
// Experimental.
type UserPoolResourceServerProps struct {
	// A unique resource server identifier for the resource server.
	// Experimental.
	Identifier *string `json:"identifier"`
	// Oauth scopes.
	// Experimental.
	Scopes *[]ResourceServerScope `json:"scopes"`
	// A friendly name for the resource server.
	// Experimental.
	UserPoolResourceServerName *string `json:"userPoolResourceServerName"`
	// The user pool to add this resource server to.
	// Experimental.
	UserPool IUserPool `json:"userPool"`
}

// Triggers for a user pool.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools-working-with-aws-lambda-triggers.html
//
// Experimental.
type UserPoolTriggers struct {
	// Creates an authentication challenge.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-create-auth-challenge.html
	//
	// Experimental.
	CreateAuthChallenge awslambda.IFunction `json:"createAuthChallenge"`
	// A custom Message AWS Lambda trigger.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-custom-message.html
	//
	// Experimental.
	CustomMessage awslambda.IFunction `json:"customMessage"`
	// Defines the authentication challenge.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-define-auth-challenge.html
	//
	// Experimental.
	DefineAuthChallenge awslambda.IFunction `json:"defineAuthChallenge"`
	// A post-authentication AWS Lambda trigger.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-post-authentication.html
	//
	// Experimental.
	PostAuthentication awslambda.IFunction `json:"postAuthentication"`
	// A post-confirmation AWS Lambda trigger.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-post-confirmation.html
	//
	// Experimental.
	PostConfirmation awslambda.IFunction `json:"postConfirmation"`
	// A pre-authentication AWS Lambda trigger.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-pre-authentication.html
	//
	// Experimental.
	PreAuthentication awslambda.IFunction `json:"preAuthentication"`
	// A pre-registration AWS Lambda trigger.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-pre-sign-up.html
	//
	// Experimental.
	PreSignUp awslambda.IFunction `json:"preSignUp"`
	// A pre-token-generation AWS Lambda trigger.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-pre-token-generation.html
	//
	// Experimental.
	PreTokenGeneration awslambda.IFunction `json:"preTokenGeneration"`
	// A user-migration AWS Lambda trigger.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-migrate-user.html
	//
	// Experimental.
	UserMigration awslambda.IFunction `json:"userMigration"`
	// Verifies the authentication challenge response.
	// See: https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-verify-auth-challenge-response.html
	//
	// Experimental.
	VerifyAuthChallengeResponse awslambda.IFunction `json:"verifyAuthChallengeResponse"`
}

// User pool configuration for user self sign up.
// Experimental.
type UserVerificationConfig struct {
	// The email body template for the verification email sent to the user upon sign up.
	//
	// See https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pool-settings-message-templates.html to
	// learn more about message templates.
	// Experimental.
	EmailBody *string `json:"emailBody"`
	// Emails can be verified either using a code or a link.
	//
	// Learn more at https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pool-settings-email-verification-message-customization.html
	// Experimental.
	EmailStyle VerificationEmailStyle `json:"emailStyle"`
	// The email subject template for the verification email sent to the user upon sign up.
	//
	// See https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pool-settings-message-templates.html to
	// learn more about message templates.
	// Experimental.
	EmailSubject *string `json:"emailSubject"`
	// The message template for the verification SMS sent to the user upon sign up.
	//
	// See https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pool-settings-message-templates.html to
	// learn more about message templates.
	// Experimental.
	SmsMessage *string `json:"smsMessage"`
}

// The email verification style.
// Experimental.
type VerificationEmailStyle string

const (
	VerificationEmailStyle_CODE VerificationEmailStyle = "CODE"
	VerificationEmailStyle_LINK VerificationEmailStyle = "LINK"
)


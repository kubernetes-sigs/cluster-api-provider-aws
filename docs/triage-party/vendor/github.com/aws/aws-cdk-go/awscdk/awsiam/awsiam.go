package awsiam

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsiam/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// Specify AWS account ID as the principal entity in a policy to delegate authority to the account.
// Experimental.
type AccountPrincipal interface {
	ArnPrincipal
	AccountId() interface{}
	Arn() *string
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for AccountPrincipal
type jsiiProxy_AccountPrincipal struct {
	jsiiProxy_ArnPrincipal
}

func (j *jsiiProxy_AccountPrincipal) AccountId() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accountId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountPrincipal) Arn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewAccountPrincipal(accountId interface{}) AccountPrincipal {
	_init_.Initialize()

	j := jsiiProxy_AccountPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.AccountPrincipal",
		[]interface{}{accountId},
		&j,
	)

	return &j
}

// Experimental.
func NewAccountPrincipal_Override(a AccountPrincipal, accountId interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.AccountPrincipal",
		[]interface{}{accountId},
		a,
	)
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_AccountPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		a,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_AccountPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		a,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (a *jsiiProxy_AccountPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		a,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (a *jsiiProxy_AccountPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (a *jsiiProxy_AccountPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		a,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// Use the AWS account into which a stack is deployed as the principal entity in a policy.
// Experimental.
type AccountRootPrincipal interface {
	AccountPrincipal
	AccountId() interface{}
	Arn() *string
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for AccountRootPrincipal
type jsiiProxy_AccountRootPrincipal struct {
	jsiiProxy_AccountPrincipal
}

func (j *jsiiProxy_AccountRootPrincipal) AccountId() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accountId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountRootPrincipal) Arn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountRootPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountRootPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountRootPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AccountRootPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewAccountRootPrincipal() AccountRootPrincipal {
	_init_.Initialize()

	j := jsiiProxy_AccountRootPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.AccountRootPrincipal",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewAccountRootPrincipal_Override(a AccountRootPrincipal) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.AccountRootPrincipal",
		nil, // no parameters
		a,
	)
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_AccountRootPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		a,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_AccountRootPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		a,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (a *jsiiProxy_AccountRootPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		a,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (a *jsiiProxy_AccountRootPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (a *jsiiProxy_AccountRootPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		a,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// Result of calling `addToPrincipalPolicy`.
// Experimental.
type AddToPrincipalPolicyResult struct {
	// Whether the statement was added to the identity's policies.
	// Experimental.
	StatementAdded *bool `json:"statementAdded"`
	// Dependable which allows depending on the policy change being applied.
	// Experimental.
	PolicyDependable awscdk.IDependable `json:"policyDependable"`
}

// Result of calling addToResourcePolicy.
// Experimental.
type AddToResourcePolicyResult struct {
	// Whether the statement was added.
	// Experimental.
	StatementAdded *bool `json:"statementAdded"`
	// Dependable which allows depending on the policy change being applied.
	// Experimental.
	PolicyDependable awscdk.IDependable `json:"policyDependable"`
}

// A principal representing all identities in all accounts.
// Experimental.
type AnyPrincipal interface {
	ArnPrincipal
	Arn() *string
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for AnyPrincipal
type jsiiProxy_AnyPrincipal struct {
	jsiiProxy_ArnPrincipal
}

func (j *jsiiProxy_AnyPrincipal) Arn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AnyPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AnyPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AnyPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_AnyPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewAnyPrincipal() AnyPrincipal {
	_init_.Initialize()

	j := jsiiProxy_AnyPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.AnyPrincipal",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewAnyPrincipal_Override(a AnyPrincipal) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.AnyPrincipal",
		nil, // no parameters
		a,
	)
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_AnyPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		a,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_AnyPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		a,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (a *jsiiProxy_AnyPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		a,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (a *jsiiProxy_AnyPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (a *jsiiProxy_AnyPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		a,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// A principal representing all identities in all accounts.
// Deprecated: use `AnyPrincipal`
type Anyone interface {
	AnyPrincipal
	Arn() *string
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for Anyone
type jsiiProxy_Anyone struct {
	jsiiProxy_AnyPrincipal
}

func (j *jsiiProxy_Anyone) Arn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Anyone) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Anyone) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Anyone) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Anyone) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Deprecated: use `AnyPrincipal`
func NewAnyone() Anyone {
	_init_.Initialize()

	j := jsiiProxy_Anyone{}

	_jsii_.Create(
		"monocdk.aws_iam.Anyone",
		nil, // no parameters
		&j,
	)

	return &j
}

// Deprecated: use `AnyPrincipal`
func NewAnyone_Override(a Anyone) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.Anyone",
		nil, // no parameters
		a,
	)
}

// Add to the policy of this principal.
// Deprecated: use `AnyPrincipal`
func (a *jsiiProxy_Anyone) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		a,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Deprecated: use `AnyPrincipal`
func (a *jsiiProxy_Anyone) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		a,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Deprecated: use `AnyPrincipal`
func (a *jsiiProxy_Anyone) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		a,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Deprecated: use `AnyPrincipal`
func (a *jsiiProxy_Anyone) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Deprecated: use `AnyPrincipal`
func (a *jsiiProxy_Anyone) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		a,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// Specify a principal by the Amazon Resource Name (ARN).
//
// You can specify AWS accounts, IAM users, Federated SAML users, IAM roles, and specific assumed-role sessions.
// You cannot specify IAM groups or instance profiles as principals
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_principal.html
//
// Experimental.
type ArnPrincipal interface {
	PrincipalBase
	Arn() *string
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for ArnPrincipal
type jsiiProxy_ArnPrincipal struct {
	jsiiProxy_PrincipalBase
}

func (j *jsiiProxy_ArnPrincipal) Arn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"arn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ArnPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ArnPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ArnPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ArnPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewArnPrincipal(arn *string) ArnPrincipal {
	_init_.Initialize()

	j := jsiiProxy_ArnPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.ArnPrincipal",
		[]interface{}{arn},
		&j,
	)

	return &j
}

// Experimental.
func NewArnPrincipal_Override(a ArnPrincipal, arn *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.ArnPrincipal",
		[]interface{}{arn},
		a,
	)
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_ArnPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		a,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (a *jsiiProxy_ArnPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		a,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (a *jsiiProxy_ArnPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		a,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (a *jsiiProxy_ArnPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (a *jsiiProxy_ArnPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		a,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// A policy principal for canonicalUserIds - useful for S3 bucket policies that use Origin Access identities.
//
// See https://docs.aws.amazon.com/general/latest/gr/acct-identifiers.html
//
// and
//
// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/private-content-restricting-access-to-s3.html
//
// for more details.
// Experimental.
type CanonicalUserPrincipal interface {
	PrincipalBase
	AssumeRoleAction() *string
	CanonicalUserId() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for CanonicalUserPrincipal
type jsiiProxy_CanonicalUserPrincipal struct {
	jsiiProxy_PrincipalBase
}

func (j *jsiiProxy_CanonicalUserPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CanonicalUserPrincipal) CanonicalUserId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"canonicalUserId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CanonicalUserPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CanonicalUserPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CanonicalUserPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewCanonicalUserPrincipal(canonicalUserId *string) CanonicalUserPrincipal {
	_init_.Initialize()

	j := jsiiProxy_CanonicalUserPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.CanonicalUserPrincipal",
		[]interface{}{canonicalUserId},
		&j,
	)

	return &j
}

// Experimental.
func NewCanonicalUserPrincipal_Override(c CanonicalUserPrincipal, canonicalUserId *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CanonicalUserPrincipal",
		[]interface{}{canonicalUserId},
		c,
	)
}

// Add to the policy of this principal.
// Experimental.
func (c *jsiiProxy_CanonicalUserPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		c,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (c *jsiiProxy_CanonicalUserPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		c,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (c *jsiiProxy_CanonicalUserPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		c,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (c *jsiiProxy_CanonicalUserPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (c *jsiiProxy_CanonicalUserPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		c,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// A CloudFormation `AWS::IAM::AccessKey`.
type CfnAccessKey interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrSecretAccessKey() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Serial() *float64
	SetSerial(val *float64)
	Stack() awscdk.Stack
	Status() *string
	SetStatus(val *string)
	UpdatedProperites() *map[string]interface{}
	UserName() *string
	SetUserName(val *string)
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

// The jsii proxy struct for CfnAccessKey
type jsiiProxy_CfnAccessKey struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAccessKey) AttrSecretAccessKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrSecretAccessKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) Serial() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"serial",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) Status() *string {
	var returns *string
	_jsii_.Get(
		j,
		"status",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessKey) UserName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userName",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::AccessKey`.
func NewCfnAccessKey(scope awscdk.Construct, id *string, props *CfnAccessKeyProps) CfnAccessKey {
	_init_.Initialize()

	j := jsiiProxy_CfnAccessKey{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnAccessKey",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::AccessKey`.
func NewCfnAccessKey_Override(c CfnAccessKey, scope awscdk.Construct, id *string, props *CfnAccessKeyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnAccessKey",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAccessKey) SetSerial(val *float64) {
	_jsii_.Set(
		j,
		"serial",
		val,
	)
}

func (j *jsiiProxy_CfnAccessKey) SetStatus(val *string) {
	_jsii_.Set(
		j,
		"status",
		val,
	)
}

func (j *jsiiProxy_CfnAccessKey) SetUserName(val *string) {
	_jsii_.Set(
		j,
		"userName",
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
func CfnAccessKey_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnAccessKey",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAccessKey_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnAccessKey",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAccessKey_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnAccessKey",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAccessKey_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnAccessKey",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAccessKey) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAccessKey) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAccessKey) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAccessKey) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAccessKey) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAccessKey) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAccessKey) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAccessKey) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAccessKey) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAccessKey) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAccessKey) OnPrepare() {
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
func (c *jsiiProxy_CfnAccessKey) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAccessKey) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAccessKey) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAccessKey) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAccessKey) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAccessKey) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAccessKey) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAccessKey) ToString() *string {
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
func (c *jsiiProxy_CfnAccessKey) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAccessKey) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::AccessKey`.
type CfnAccessKeyProps struct {
	// `AWS::IAM::AccessKey.UserName`.
	UserName *string `json:"userName"`
	// `AWS::IAM::AccessKey.Serial`.
	Serial *float64 `json:"serial"`
	// `AWS::IAM::AccessKey.Status`.
	Status *string `json:"status"`
}

// A CloudFormation `AWS::IAM::Group`.
type CfnGroup interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	GroupName() *string
	SetGroupName(val *string)
	LogicalId() *string
	ManagedPolicyArns() *[]*string
	SetManagedPolicyArns(val *[]*string)
	Node() awscdk.ConstructNode
	Path() *string
	SetPath(val *string)
	Policies() interface{}
	SetPolicies(val interface{})
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

// The jsii proxy struct for CfnGroup
type jsiiProxy_CfnGroup struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnGroup) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) GroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) ManagedPolicyArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"managedPolicyArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) Policies() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGroup) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::Group`.
func NewCfnGroup(scope awscdk.Construct, id *string, props *CfnGroupProps) CfnGroup {
	_init_.Initialize()

	j := jsiiProxy_CfnGroup{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::Group`.
func NewCfnGroup_Override(c CfnGroup, scope awscdk.Construct, id *string, props *CfnGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnGroup",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnGroup) SetGroupName(val *string) {
	_jsii_.Set(
		j,
		"groupName",
		val,
	)
}

func (j *jsiiProxy_CfnGroup) SetManagedPolicyArns(val *[]*string) {
	_jsii_.Set(
		j,
		"managedPolicyArns",
		val,
	)
}

func (j *jsiiProxy_CfnGroup) SetPath(val *string) {
	_jsii_.Set(
		j,
		"path",
		val,
	)
}

func (j *jsiiProxy_CfnGroup) SetPolicies(val interface{}) {
	_jsii_.Set(
		j,
		"policies",
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
func CfnGroup_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnGroup",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnGroup_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnGroup",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnGroup_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnGroup",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnGroup) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnGroup) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnGroup) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnGroup) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnGroup) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnGroup) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnGroup) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnGroup) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnGroup) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnGroup) OnPrepare() {
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
func (c *jsiiProxy_CfnGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnGroup) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnGroup) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnGroup) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnGroup) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnGroup) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnGroup) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnGroup) ToString() *string {
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
func (c *jsiiProxy_CfnGroup) Validate() *[]*string {
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
func (c *jsiiProxy_CfnGroup) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnGroup_PolicyProperty struct {
	// `CfnGroup.PolicyProperty.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
	// `CfnGroup.PolicyProperty.PolicyName`.
	PolicyName *string `json:"policyName"`
}

// Properties for defining a `AWS::IAM::Group`.
type CfnGroupProps struct {
	// `AWS::IAM::Group.GroupName`.
	GroupName *string `json:"groupName"`
	// `AWS::IAM::Group.ManagedPolicyArns`.
	ManagedPolicyArns *[]*string `json:"managedPolicyArns"`
	// `AWS::IAM::Group.Path`.
	Path *string `json:"path"`
	// `AWS::IAM::Group.Policies`.
	Policies interface{} `json:"policies"`
}

// A CloudFormation `AWS::IAM::InstanceProfile`.
type CfnInstanceProfile interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	InstanceProfileName() *string
	SetInstanceProfileName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Path() *string
	SetPath(val *string)
	Ref() *string
	Roles() *[]*string
	SetRoles(val *[]*string)
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

// The jsii proxy struct for CfnInstanceProfile
type jsiiProxy_CfnInstanceProfile struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnInstanceProfile) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) InstanceProfileName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"instanceProfileName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) Roles() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"roles",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnInstanceProfile) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::InstanceProfile`.
func NewCfnInstanceProfile(scope awscdk.Construct, id *string, props *CfnInstanceProfileProps) CfnInstanceProfile {
	_init_.Initialize()

	j := jsiiProxy_CfnInstanceProfile{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnInstanceProfile",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::InstanceProfile`.
func NewCfnInstanceProfile_Override(c CfnInstanceProfile, scope awscdk.Construct, id *string, props *CfnInstanceProfileProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnInstanceProfile",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnInstanceProfile) SetInstanceProfileName(val *string) {
	_jsii_.Set(
		j,
		"instanceProfileName",
		val,
	)
}

func (j *jsiiProxy_CfnInstanceProfile) SetPath(val *string) {
	_jsii_.Set(
		j,
		"path",
		val,
	)
}

func (j *jsiiProxy_CfnInstanceProfile) SetRoles(val *[]*string) {
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
func CfnInstanceProfile_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnInstanceProfile",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnInstanceProfile_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnInstanceProfile",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnInstanceProfile_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnInstanceProfile",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnInstanceProfile_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnInstanceProfile",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnInstanceProfile) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnInstanceProfile) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnInstanceProfile) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnInstanceProfile) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnInstanceProfile) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnInstanceProfile) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnInstanceProfile) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnInstanceProfile) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnInstanceProfile) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnInstanceProfile) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnInstanceProfile) OnPrepare() {
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
func (c *jsiiProxy_CfnInstanceProfile) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnInstanceProfile) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnInstanceProfile) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnInstanceProfile) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnInstanceProfile) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnInstanceProfile) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnInstanceProfile) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnInstanceProfile) ToString() *string {
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
func (c *jsiiProxy_CfnInstanceProfile) Validate() *[]*string {
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
func (c *jsiiProxy_CfnInstanceProfile) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::InstanceProfile`.
type CfnInstanceProfileProps struct {
	// `AWS::IAM::InstanceProfile.Roles`.
	Roles *[]*string `json:"roles"`
	// `AWS::IAM::InstanceProfile.InstanceProfileName`.
	InstanceProfileName *string `json:"instanceProfileName"`
	// `AWS::IAM::InstanceProfile.Path`.
	Path *string `json:"path"`
}

// A CloudFormation `AWS::IAM::ManagedPolicy`.
type CfnManagedPolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	Groups() *[]*string
	SetGroups(val *[]*string)
	LogicalId() *string
	ManagedPolicyName() *string
	SetManagedPolicyName(val *string)
	Node() awscdk.ConstructNode
	Path() *string
	SetPath(val *string)
	PolicyDocument() interface{}
	SetPolicyDocument(val interface{})
	Ref() *string
	Roles() *[]*string
	SetRoles(val *[]*string)
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	Users() *[]*string
	SetUsers(val *[]*string)
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

// The jsii proxy struct for CfnManagedPolicy
type jsiiProxy_CfnManagedPolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnManagedPolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Groups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"groups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) ManagedPolicyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"managedPolicyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) PolicyDocument() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policyDocument",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Roles() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"roles",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnManagedPolicy) Users() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"users",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::ManagedPolicy`.
func NewCfnManagedPolicy(scope awscdk.Construct, id *string, props *CfnManagedPolicyProps) CfnManagedPolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnManagedPolicy{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnManagedPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::ManagedPolicy`.
func NewCfnManagedPolicy_Override(c CfnManagedPolicy, scope awscdk.Construct, id *string, props *CfnManagedPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnManagedPolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnManagedPolicy) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnManagedPolicy) SetGroups(val *[]*string) {
	_jsii_.Set(
		j,
		"groups",
		val,
	)
}

func (j *jsiiProxy_CfnManagedPolicy) SetManagedPolicyName(val *string) {
	_jsii_.Set(
		j,
		"managedPolicyName",
		val,
	)
}

func (j *jsiiProxy_CfnManagedPolicy) SetPath(val *string) {
	_jsii_.Set(
		j,
		"path",
		val,
	)
}

func (j *jsiiProxy_CfnManagedPolicy) SetPolicyDocument(val interface{}) {
	_jsii_.Set(
		j,
		"policyDocument",
		val,
	)
}

func (j *jsiiProxy_CfnManagedPolicy) SetRoles(val *[]*string) {
	_jsii_.Set(
		j,
		"roles",
		val,
	)
}

func (j *jsiiProxy_CfnManagedPolicy) SetUsers(val *[]*string) {
	_jsii_.Set(
		j,
		"users",
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
func CfnManagedPolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnManagedPolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnManagedPolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnManagedPolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnManagedPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnManagedPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnManagedPolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnManagedPolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnManagedPolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnManagedPolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnManagedPolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnManagedPolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnManagedPolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnManagedPolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnManagedPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnManagedPolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnManagedPolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnManagedPolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnManagedPolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnManagedPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnManagedPolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnManagedPolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnManagedPolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnManagedPolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnManagedPolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnManagedPolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnManagedPolicy) ToString() *string {
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
func (c *jsiiProxy_CfnManagedPolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnManagedPolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::ManagedPolicy`.
type CfnManagedPolicyProps struct {
	// `AWS::IAM::ManagedPolicy.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
	// `AWS::IAM::ManagedPolicy.Description`.
	Description *string `json:"description"`
	// `AWS::IAM::ManagedPolicy.Groups`.
	Groups *[]*string `json:"groups"`
	// `AWS::IAM::ManagedPolicy.ManagedPolicyName`.
	ManagedPolicyName *string `json:"managedPolicyName"`
	// `AWS::IAM::ManagedPolicy.Path`.
	Path *string `json:"path"`
	// `AWS::IAM::ManagedPolicy.Roles`.
	Roles *[]*string `json:"roles"`
	// `AWS::IAM::ManagedPolicy.Users`.
	Users *[]*string `json:"users"`
}

// A CloudFormation `AWS::IAM::OIDCProvider`.
type CfnOIDCProvider interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClientIdList() *[]*string
	SetClientIdList(val *[]*string)
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	ThumbprintList() *[]*string
	SetThumbprintList(val *[]*string)
	UpdatedProperites() *map[string]interface{}
	Url() *string
	SetUrl(val *string)
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

// The jsii proxy struct for CfnOIDCProvider
type jsiiProxy_CfnOIDCProvider struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnOIDCProvider) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) ClientIdList() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"clientIdList",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) ThumbprintList() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"thumbprintList",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnOIDCProvider) Url() *string {
	var returns *string
	_jsii_.Get(
		j,
		"url",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::OIDCProvider`.
func NewCfnOIDCProvider(scope awscdk.Construct, id *string, props *CfnOIDCProviderProps) CfnOIDCProvider {
	_init_.Initialize()

	j := jsiiProxy_CfnOIDCProvider{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnOIDCProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::OIDCProvider`.
func NewCfnOIDCProvider_Override(c CfnOIDCProvider, scope awscdk.Construct, id *string, props *CfnOIDCProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnOIDCProvider",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnOIDCProvider) SetClientIdList(val *[]*string) {
	_jsii_.Set(
		j,
		"clientIdList",
		val,
	)
}

func (j *jsiiProxy_CfnOIDCProvider) SetThumbprintList(val *[]*string) {
	_jsii_.Set(
		j,
		"thumbprintList",
		val,
	)
}

func (j *jsiiProxy_CfnOIDCProvider) SetUrl(val *string) {
	_jsii_.Set(
		j,
		"url",
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
func CfnOIDCProvider_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnOIDCProvider",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnOIDCProvider_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnOIDCProvider",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnOIDCProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnOIDCProvider",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnOIDCProvider_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnOIDCProvider",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnOIDCProvider) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnOIDCProvider) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnOIDCProvider) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnOIDCProvider) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnOIDCProvider) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnOIDCProvider) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnOIDCProvider) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnOIDCProvider) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnOIDCProvider) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnOIDCProvider) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnOIDCProvider) OnPrepare() {
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
func (c *jsiiProxy_CfnOIDCProvider) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnOIDCProvider) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnOIDCProvider) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnOIDCProvider) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnOIDCProvider) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnOIDCProvider) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnOIDCProvider) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnOIDCProvider) ToString() *string {
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
func (c *jsiiProxy_CfnOIDCProvider) Validate() *[]*string {
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
func (c *jsiiProxy_CfnOIDCProvider) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::OIDCProvider`.
type CfnOIDCProviderProps struct {
	// `AWS::IAM::OIDCProvider.ThumbprintList`.
	ThumbprintList *[]*string `json:"thumbprintList"`
	// `AWS::IAM::OIDCProvider.ClientIdList`.
	ClientIdList *[]*string `json:"clientIdList"`
	// `AWS::IAM::OIDCProvider.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::IAM::OIDCProvider.Url`.
	Url *string `json:"url"`
}

// A CloudFormation `AWS::IAM::Policy`.
type CfnPolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Groups() *[]*string
	SetGroups(val *[]*string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	PolicyDocument() interface{}
	SetPolicyDocument(val interface{})
	PolicyName() *string
	SetPolicyName(val *string)
	Ref() *string
	Roles() *[]*string
	SetRoles(val *[]*string)
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	Users() *[]*string
	SetUsers(val *[]*string)
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

// The jsii proxy struct for CfnPolicy
type jsiiProxy_CfnPolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnPolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) Groups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"groups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) PolicyDocument() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policyDocument",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) PolicyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"policyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) Roles() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"roles",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnPolicy) Users() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"users",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::Policy`.
func NewCfnPolicy(scope awscdk.Construct, id *string, props *CfnPolicyProps) CfnPolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnPolicy{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::Policy`.
func NewCfnPolicy_Override(c CfnPolicy, scope awscdk.Construct, id *string, props *CfnPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnPolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnPolicy) SetGroups(val *[]*string) {
	_jsii_.Set(
		j,
		"groups",
		val,
	)
}

func (j *jsiiProxy_CfnPolicy) SetPolicyDocument(val interface{}) {
	_jsii_.Set(
		j,
		"policyDocument",
		val,
	)
}

func (j *jsiiProxy_CfnPolicy) SetPolicyName(val *string) {
	_jsii_.Set(
		j,
		"policyName",
		val,
	)
}

func (j *jsiiProxy_CfnPolicy) SetRoles(val *[]*string) {
	_jsii_.Set(
		j,
		"roles",
		val,
	)
}

func (j *jsiiProxy_CfnPolicy) SetUsers(val *[]*string) {
	_jsii_.Set(
		j,
		"users",
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
func CfnPolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnPolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnPolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnPolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnPolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnPolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnPolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnPolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnPolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnPolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnPolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnPolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnPolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnPolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnPolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnPolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnPolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnPolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnPolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnPolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnPolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnPolicy) ToString() *string {
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
func (c *jsiiProxy_CfnPolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnPolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::Policy`.
type CfnPolicyProps struct {
	// `AWS::IAM::Policy.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
	// `AWS::IAM::Policy.PolicyName`.
	PolicyName *string `json:"policyName"`
	// `AWS::IAM::Policy.Groups`.
	Groups *[]*string `json:"groups"`
	// `AWS::IAM::Policy.Roles`.
	Roles *[]*string `json:"roles"`
	// `AWS::IAM::Policy.Users`.
	Users *[]*string `json:"users"`
}

// A CloudFormation `AWS::IAM::Role`.
type CfnRole interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AssumeRolePolicyDocument() interface{}
	SetAssumeRolePolicyDocument(val interface{})
	AttrArn() *string
	AttrRoleId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	ManagedPolicyArns() *[]*string
	SetManagedPolicyArns(val *[]*string)
	MaxSessionDuration() *float64
	SetMaxSessionDuration(val *float64)
	Node() awscdk.ConstructNode
	Path() *string
	SetPath(val *string)
	PermissionsBoundary() *string
	SetPermissionsBoundary(val *string)
	Policies() interface{}
	SetPolicies(val interface{})
	Ref() *string
	RoleName() *string
	SetRoleName(val *string)
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

// The jsii proxy struct for CfnRole
type jsiiProxy_CfnRole struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRole) AssumeRolePolicyDocument() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"assumeRolePolicyDocument",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) AttrRoleId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRoleId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) ManagedPolicyArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"managedPolicyArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) MaxSessionDuration() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"maxSessionDuration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) PermissionsBoundary() *string {
	var returns *string
	_jsii_.Get(
		j,
		"permissionsBoundary",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) Policies() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) RoleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRole) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::Role`.
func NewCfnRole(scope awscdk.Construct, id *string, props *CfnRoleProps) CfnRole {
	_init_.Initialize()

	j := jsiiProxy_CfnRole{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnRole",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::Role`.
func NewCfnRole_Override(c CfnRole, scope awscdk.Construct, id *string, props *CfnRoleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnRole",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRole) SetAssumeRolePolicyDocument(val interface{}) {
	_jsii_.Set(
		j,
		"assumeRolePolicyDocument",
		val,
	)
}

func (j *jsiiProxy_CfnRole) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnRole) SetManagedPolicyArns(val *[]*string) {
	_jsii_.Set(
		j,
		"managedPolicyArns",
		val,
	)
}

func (j *jsiiProxy_CfnRole) SetMaxSessionDuration(val *float64) {
	_jsii_.Set(
		j,
		"maxSessionDuration",
		val,
	)
}

func (j *jsiiProxy_CfnRole) SetPath(val *string) {
	_jsii_.Set(
		j,
		"path",
		val,
	)
}

func (j *jsiiProxy_CfnRole) SetPermissionsBoundary(val *string) {
	_jsii_.Set(
		j,
		"permissionsBoundary",
		val,
	)
}

func (j *jsiiProxy_CfnRole) SetPolicies(val interface{}) {
	_jsii_.Set(
		j,
		"policies",
		val,
	)
}

func (j *jsiiProxy_CfnRole) SetRoleName(val *string) {
	_jsii_.Set(
		j,
		"roleName",
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
func CfnRole_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnRole",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnRole_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnRole",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnRole_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnRole",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRole_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnRole",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnRole) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnRole) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnRole) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnRole) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnRole) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnRole) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnRole) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnRole) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnRole) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnRole) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnRole) OnPrepare() {
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
func (c *jsiiProxy_CfnRole) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRole) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnRole) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnRole) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnRole) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnRole) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnRole) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRole) ToString() *string {
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
func (c *jsiiProxy_CfnRole) Validate() *[]*string {
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
func (c *jsiiProxy_CfnRole) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnRole_PolicyProperty struct {
	// `CfnRole.PolicyProperty.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
	// `CfnRole.PolicyProperty.PolicyName`.
	PolicyName *string `json:"policyName"`
}

// Properties for defining a `AWS::IAM::Role`.
type CfnRoleProps struct {
	// `AWS::IAM::Role.AssumeRolePolicyDocument`.
	AssumeRolePolicyDocument interface{} `json:"assumeRolePolicyDocument"`
	// `AWS::IAM::Role.Description`.
	Description *string `json:"description"`
	// `AWS::IAM::Role.ManagedPolicyArns`.
	ManagedPolicyArns *[]*string `json:"managedPolicyArns"`
	// `AWS::IAM::Role.MaxSessionDuration`.
	MaxSessionDuration *float64 `json:"maxSessionDuration"`
	// `AWS::IAM::Role.Path`.
	Path *string `json:"path"`
	// `AWS::IAM::Role.PermissionsBoundary`.
	PermissionsBoundary *string `json:"permissionsBoundary"`
	// `AWS::IAM::Role.Policies`.
	Policies interface{} `json:"policies"`
	// `AWS::IAM::Role.RoleName`.
	RoleName *string `json:"roleName"`
	// `AWS::IAM::Role.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::IAM::SAMLProvider`.
type CfnSAMLProvider interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	SamlMetadataDocument() *string
	SetSamlMetadataDocument(val *string)
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

// The jsii proxy struct for CfnSAMLProvider
type jsiiProxy_CfnSAMLProvider struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnSAMLProvider) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) SamlMetadataDocument() *string {
	var returns *string
	_jsii_.Get(
		j,
		"samlMetadataDocument",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnSAMLProvider) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::SAMLProvider`.
func NewCfnSAMLProvider(scope awscdk.Construct, id *string, props *CfnSAMLProviderProps) CfnSAMLProvider {
	_init_.Initialize()

	j := jsiiProxy_CfnSAMLProvider{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnSAMLProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::SAMLProvider`.
func NewCfnSAMLProvider_Override(c CfnSAMLProvider, scope awscdk.Construct, id *string, props *CfnSAMLProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnSAMLProvider",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnSAMLProvider) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnSAMLProvider) SetSamlMetadataDocument(val *string) {
	_jsii_.Set(
		j,
		"samlMetadataDocument",
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
func CfnSAMLProvider_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnSAMLProvider",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnSAMLProvider_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnSAMLProvider",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnSAMLProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnSAMLProvider",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnSAMLProvider_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnSAMLProvider",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnSAMLProvider) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnSAMLProvider) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnSAMLProvider) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnSAMLProvider) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnSAMLProvider) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnSAMLProvider) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnSAMLProvider) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnSAMLProvider) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnSAMLProvider) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnSAMLProvider) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnSAMLProvider) OnPrepare() {
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
func (c *jsiiProxy_CfnSAMLProvider) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSAMLProvider) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnSAMLProvider) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnSAMLProvider) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnSAMLProvider) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnSAMLProvider) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnSAMLProvider) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnSAMLProvider) ToString() *string {
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
func (c *jsiiProxy_CfnSAMLProvider) Validate() *[]*string {
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
func (c *jsiiProxy_CfnSAMLProvider) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::SAMLProvider`.
type CfnSAMLProviderProps struct {
	// `AWS::IAM::SAMLProvider.SamlMetadataDocument`.
	SamlMetadataDocument *string `json:"samlMetadataDocument"`
	// `AWS::IAM::SAMLProvider.Name`.
	Name *string `json:"name"`
	// `AWS::IAM::SAMLProvider.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::IAM::ServerCertificate`.
type CfnServerCertificate interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CertificateBody() *string
	SetCertificateBody(val *string)
	CertificateChain() *string
	SetCertificateChain(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Path() *string
	SetPath(val *string)
	PrivateKey() *string
	SetPrivateKey(val *string)
	Ref() *string
	ServerCertificateName() *string
	SetServerCertificateName(val *string)
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

// The jsii proxy struct for CfnServerCertificate
type jsiiProxy_CfnServerCertificate struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnServerCertificate) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) CertificateBody() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certificateBody",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) CertificateChain() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certificateChain",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) PrivateKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"privateKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) ServerCertificateName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serverCertificateName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServerCertificate) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::ServerCertificate`.
func NewCfnServerCertificate(scope awscdk.Construct, id *string, props *CfnServerCertificateProps) CfnServerCertificate {
	_init_.Initialize()

	j := jsiiProxy_CfnServerCertificate{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnServerCertificate",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::ServerCertificate`.
func NewCfnServerCertificate_Override(c CfnServerCertificate, scope awscdk.Construct, id *string, props *CfnServerCertificateProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnServerCertificate",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnServerCertificate) SetCertificateBody(val *string) {
	_jsii_.Set(
		j,
		"certificateBody",
		val,
	)
}

func (j *jsiiProxy_CfnServerCertificate) SetCertificateChain(val *string) {
	_jsii_.Set(
		j,
		"certificateChain",
		val,
	)
}

func (j *jsiiProxy_CfnServerCertificate) SetPath(val *string) {
	_jsii_.Set(
		j,
		"path",
		val,
	)
}

func (j *jsiiProxy_CfnServerCertificate) SetPrivateKey(val *string) {
	_jsii_.Set(
		j,
		"privateKey",
		val,
	)
}

func (j *jsiiProxy_CfnServerCertificate) SetServerCertificateName(val *string) {
	_jsii_.Set(
		j,
		"serverCertificateName",
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
func CfnServerCertificate_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnServerCertificate",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnServerCertificate_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnServerCertificate",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnServerCertificate_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnServerCertificate",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnServerCertificate_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnServerCertificate",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnServerCertificate) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnServerCertificate) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnServerCertificate) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnServerCertificate) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnServerCertificate) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnServerCertificate) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnServerCertificate) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnServerCertificate) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnServerCertificate) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnServerCertificate) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnServerCertificate) OnPrepare() {
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
func (c *jsiiProxy_CfnServerCertificate) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnServerCertificate) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnServerCertificate) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnServerCertificate) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnServerCertificate) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnServerCertificate) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnServerCertificate) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnServerCertificate) ToString() *string {
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
func (c *jsiiProxy_CfnServerCertificate) Validate() *[]*string {
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
func (c *jsiiProxy_CfnServerCertificate) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::ServerCertificate`.
type CfnServerCertificateProps struct {
	// `AWS::IAM::ServerCertificate.CertificateBody`.
	CertificateBody *string `json:"certificateBody"`
	// `AWS::IAM::ServerCertificate.CertificateChain`.
	CertificateChain *string `json:"certificateChain"`
	// `AWS::IAM::ServerCertificate.Path`.
	Path *string `json:"path"`
	// `AWS::IAM::ServerCertificate.PrivateKey`.
	PrivateKey *string `json:"privateKey"`
	// `AWS::IAM::ServerCertificate.ServerCertificateName`.
	ServerCertificateName *string `json:"serverCertificateName"`
	// `AWS::IAM::ServerCertificate.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::IAM::ServiceLinkedRole`.
type CfnServiceLinkedRole interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AwsServiceName() *string
	SetAwsServiceName(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	CustomSuffix() *string
	SetCustomSuffix(val *string)
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

// The jsii proxy struct for CfnServiceLinkedRole
type jsiiProxy_CfnServiceLinkedRole struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnServiceLinkedRole) AwsServiceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"awsServiceName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) CustomSuffix() *string {
	var returns *string
	_jsii_.Get(
		j,
		"customSuffix",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnServiceLinkedRole) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::ServiceLinkedRole`.
func NewCfnServiceLinkedRole(scope awscdk.Construct, id *string, props *CfnServiceLinkedRoleProps) CfnServiceLinkedRole {
	_init_.Initialize()

	j := jsiiProxy_CfnServiceLinkedRole{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnServiceLinkedRole",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::ServiceLinkedRole`.
func NewCfnServiceLinkedRole_Override(c CfnServiceLinkedRole, scope awscdk.Construct, id *string, props *CfnServiceLinkedRoleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnServiceLinkedRole",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnServiceLinkedRole) SetAwsServiceName(val *string) {
	_jsii_.Set(
		j,
		"awsServiceName",
		val,
	)
}

func (j *jsiiProxy_CfnServiceLinkedRole) SetCustomSuffix(val *string) {
	_jsii_.Set(
		j,
		"customSuffix",
		val,
	)
}

func (j *jsiiProxy_CfnServiceLinkedRole) SetDescription(val *string) {
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
func CfnServiceLinkedRole_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnServiceLinkedRole",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnServiceLinkedRole_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnServiceLinkedRole",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnServiceLinkedRole_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnServiceLinkedRole",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnServiceLinkedRole_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnServiceLinkedRole",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnServiceLinkedRole) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnServiceLinkedRole) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnServiceLinkedRole) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnServiceLinkedRole) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnServiceLinkedRole) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) OnPrepare() {
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
func (c *jsiiProxy_CfnServiceLinkedRole) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnServiceLinkedRole) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnServiceLinkedRole) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnServiceLinkedRole) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnServiceLinkedRole) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnServiceLinkedRole) ToString() *string {
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
func (c *jsiiProxy_CfnServiceLinkedRole) Validate() *[]*string {
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
func (c *jsiiProxy_CfnServiceLinkedRole) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::ServiceLinkedRole`.
type CfnServiceLinkedRoleProps struct {
	// `AWS::IAM::ServiceLinkedRole.AWSServiceName`.
	AwsServiceName *string `json:"awsServiceName"`
	// `AWS::IAM::ServiceLinkedRole.CustomSuffix`.
	CustomSuffix *string `json:"customSuffix"`
	// `AWS::IAM::ServiceLinkedRole.Description`.
	Description *string `json:"description"`
}

// A CloudFormation `AWS::IAM::User`.
type CfnUser interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Groups() *[]*string
	SetGroups(val *[]*string)
	LogicalId() *string
	LoginProfile() interface{}
	SetLoginProfile(val interface{})
	ManagedPolicyArns() *[]*string
	SetManagedPolicyArns(val *[]*string)
	Node() awscdk.ConstructNode
	Path() *string
	SetPath(val *string)
	PermissionsBoundary() *string
	SetPermissionsBoundary(val *string)
	Policies() interface{}
	SetPolicies(val interface{})
	Ref() *string
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	UpdatedProperites() *map[string]interface{}
	UserName() *string
	SetUserName(val *string)
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

// The jsii proxy struct for CfnUser
type jsiiProxy_CfnUser struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUser) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) Groups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"groups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) LoginProfile() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"loginProfile",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) ManagedPolicyArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"managedPolicyArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) PermissionsBoundary() *string {
	var returns *string
	_jsii_.Get(
		j,
		"permissionsBoundary",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) Policies() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUser) UserName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userName",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::User`.
func NewCfnUser(scope awscdk.Construct, id *string, props *CfnUserProps) CfnUser {
	_init_.Initialize()

	j := jsiiProxy_CfnUser{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnUser",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::User`.
func NewCfnUser_Override(c CfnUser, scope awscdk.Construct, id *string, props *CfnUserProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnUser",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUser) SetGroups(val *[]*string) {
	_jsii_.Set(
		j,
		"groups",
		val,
	)
}

func (j *jsiiProxy_CfnUser) SetLoginProfile(val interface{}) {
	_jsii_.Set(
		j,
		"loginProfile",
		val,
	)
}

func (j *jsiiProxy_CfnUser) SetManagedPolicyArns(val *[]*string) {
	_jsii_.Set(
		j,
		"managedPolicyArns",
		val,
	)
}

func (j *jsiiProxy_CfnUser) SetPath(val *string) {
	_jsii_.Set(
		j,
		"path",
		val,
	)
}

func (j *jsiiProxy_CfnUser) SetPermissionsBoundary(val *string) {
	_jsii_.Set(
		j,
		"permissionsBoundary",
		val,
	)
}

func (j *jsiiProxy_CfnUser) SetPolicies(val interface{}) {
	_jsii_.Set(
		j,
		"policies",
		val,
	)
}

func (j *jsiiProxy_CfnUser) SetUserName(val *string) {
	_jsii_.Set(
		j,
		"userName",
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
func CfnUser_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnUser",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUser_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnUser",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUser_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnUser",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUser_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnUser",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUser) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUser) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUser) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUser) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUser) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUser) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUser) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUser) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUser) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUser) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUser) OnPrepare() {
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
func (c *jsiiProxy_CfnUser) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUser) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUser) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUser) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUser) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUser) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUser) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUser) ToString() *string {
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
func (c *jsiiProxy_CfnUser) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUser) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUser_LoginProfileProperty struct {
	// `CfnUser.LoginProfileProperty.Password`.
	Password *string `json:"password"`
	// `CfnUser.LoginProfileProperty.PasswordResetRequired`.
	PasswordResetRequired interface{} `json:"passwordResetRequired"`
}

type CfnUser_PolicyProperty struct {
	// `CfnUser.PolicyProperty.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
	// `CfnUser.PolicyProperty.PolicyName`.
	PolicyName *string `json:"policyName"`
}

// Properties for defining a `AWS::IAM::User`.
type CfnUserProps struct {
	// `AWS::IAM::User.Groups`.
	Groups *[]*string `json:"groups"`
	// `AWS::IAM::User.LoginProfile`.
	LoginProfile interface{} `json:"loginProfile"`
	// `AWS::IAM::User.ManagedPolicyArns`.
	ManagedPolicyArns *[]*string `json:"managedPolicyArns"`
	// `AWS::IAM::User.Path`.
	Path *string `json:"path"`
	// `AWS::IAM::User.PermissionsBoundary`.
	PermissionsBoundary *string `json:"permissionsBoundary"`
	// `AWS::IAM::User.Policies`.
	Policies interface{} `json:"policies"`
	// `AWS::IAM::User.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::IAM::User.UserName`.
	UserName *string `json:"userName"`
}

// A CloudFormation `AWS::IAM::UserToGroupAddition`.
type CfnUserToGroupAddition interface {
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
	Users() *[]*string
	SetUsers(val *[]*string)
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

// The jsii proxy struct for CfnUserToGroupAddition
type jsiiProxy_CfnUserToGroupAddition struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUserToGroupAddition) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) GroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUserToGroupAddition) Users() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"users",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::UserToGroupAddition`.
func NewCfnUserToGroupAddition(scope awscdk.Construct, id *string, props *CfnUserToGroupAdditionProps) CfnUserToGroupAddition {
	_init_.Initialize()

	j := jsiiProxy_CfnUserToGroupAddition{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnUserToGroupAddition",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::UserToGroupAddition`.
func NewCfnUserToGroupAddition_Override(c CfnUserToGroupAddition, scope awscdk.Construct, id *string, props *CfnUserToGroupAdditionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnUserToGroupAddition",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUserToGroupAddition) SetGroupName(val *string) {
	_jsii_.Set(
		j,
		"groupName",
		val,
	)
}

func (j *jsiiProxy_CfnUserToGroupAddition) SetUsers(val *[]*string) {
	_jsii_.Set(
		j,
		"users",
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
func CfnUserToGroupAddition_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnUserToGroupAddition",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUserToGroupAddition_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnUserToGroupAddition",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUserToGroupAddition_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnUserToGroupAddition",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUserToGroupAddition_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnUserToGroupAddition",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUserToGroupAddition) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUserToGroupAddition) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUserToGroupAddition) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUserToGroupAddition) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUserToGroupAddition) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) OnPrepare() {
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
func (c *jsiiProxy_CfnUserToGroupAddition) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUserToGroupAddition) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUserToGroupAddition) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUserToGroupAddition) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUserToGroupAddition) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUserToGroupAddition) ToString() *string {
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
func (c *jsiiProxy_CfnUserToGroupAddition) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUserToGroupAddition) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::UserToGroupAddition`.
type CfnUserToGroupAdditionProps struct {
	// `AWS::IAM::UserToGroupAddition.GroupName`.
	GroupName *string `json:"groupName"`
	// `AWS::IAM::UserToGroupAddition.Users`.
	Users *[]*string `json:"users"`
}

// A CloudFormation `AWS::IAM::VirtualMFADevice`.
type CfnVirtualMFADevice interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrSerialNumber() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Path() *string
	SetPath(val *string)
	Ref() *string
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	UpdatedProperites() *map[string]interface{}
	Users() *[]*string
	SetUsers(val *[]*string)
	VirtualMfaDeviceName() *string
	SetVirtualMfaDeviceName(val *string)
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

// The jsii proxy struct for CfnVirtualMFADevice
type jsiiProxy_CfnVirtualMFADevice struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnVirtualMFADevice) AttrSerialNumber() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrSerialNumber",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) Users() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"users",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVirtualMFADevice) VirtualMfaDeviceName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"virtualMfaDeviceName",
		&returns,
	)
	return returns
}


// Create a new `AWS::IAM::VirtualMFADevice`.
func NewCfnVirtualMFADevice(scope awscdk.Construct, id *string, props *CfnVirtualMFADeviceProps) CfnVirtualMFADevice {
	_init_.Initialize()

	j := jsiiProxy_CfnVirtualMFADevice{}

	_jsii_.Create(
		"monocdk.aws_iam.CfnVirtualMFADevice",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::IAM::VirtualMFADevice`.
func NewCfnVirtualMFADevice_Override(c CfnVirtualMFADevice, scope awscdk.Construct, id *string, props *CfnVirtualMFADeviceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.CfnVirtualMFADevice",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnVirtualMFADevice) SetPath(val *string) {
	_jsii_.Set(
		j,
		"path",
		val,
	)
}

func (j *jsiiProxy_CfnVirtualMFADevice) SetUsers(val *[]*string) {
	_jsii_.Set(
		j,
		"users",
		val,
	)
}

func (j *jsiiProxy_CfnVirtualMFADevice) SetVirtualMfaDeviceName(val *string) {
	_jsii_.Set(
		j,
		"virtualMfaDeviceName",
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
func CfnVirtualMFADevice_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnVirtualMFADevice",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnVirtualMFADevice_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnVirtualMFADevice",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnVirtualMFADevice_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.CfnVirtualMFADevice",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnVirtualMFADevice_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_iam.CfnVirtualMFADevice",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnVirtualMFADevice) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnVirtualMFADevice) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnVirtualMFADevice) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnVirtualMFADevice) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnVirtualMFADevice) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) OnPrepare() {
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
func (c *jsiiProxy_CfnVirtualMFADevice) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnVirtualMFADevice) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnVirtualMFADevice) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnVirtualMFADevice) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnVirtualMFADevice) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnVirtualMFADevice) ToString() *string {
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
func (c *jsiiProxy_CfnVirtualMFADevice) Validate() *[]*string {
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
func (c *jsiiProxy_CfnVirtualMFADevice) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::IAM::VirtualMFADevice`.
type CfnVirtualMFADeviceProps struct {
	// `AWS::IAM::VirtualMFADevice.Users`.
	Users *[]*string `json:"users"`
	// `AWS::IAM::VirtualMFADevice.Path`.
	Path *string `json:"path"`
	// `AWS::IAM::VirtualMFADevice.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::IAM::VirtualMFADevice.VirtualMfaDeviceName`.
	VirtualMfaDeviceName *string `json:"virtualMfaDeviceName"`
}

// Basic options for a grant operation.
// Experimental.
type CommonGrantOptions struct {
	// The actions to grant.
	// Experimental.
	Actions *[]*string `json:"actions"`
	// The principal to grant to.
	// Experimental.
	Grantee IGrantable `json:"grantee"`
	// The resource ARNs to grant to.
	// Experimental.
	ResourceArns *[]*string `json:"resourceArns"`
}

// Composite dependable.
//
// Not as simple as eagerly getting the dependency roots from the
// inner dependables, as they may be mutable so we need to defer
// the query.
// Experimental.
type CompositeDependable interface {
	awscdk.IDependable
}

// The jsii proxy struct for CompositeDependable
type jsiiProxy_CompositeDependable struct {
	internal.Type__awscdkIDependable
}

// Experimental.
func NewCompositeDependable(dependables ...awscdk.IDependable) CompositeDependable {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range dependables {
		args = append(args, a)
	}

	j := jsiiProxy_CompositeDependable{}

	_jsii_.Create(
		"monocdk.aws_iam.CompositeDependable",
		args,
		&j,
	)

	return &j
}

// Experimental.
func NewCompositeDependable_Override(c CompositeDependable, dependables ...awscdk.IDependable) {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range dependables {
		args = append(args, a)
	}

	_jsii_.Create(
		"monocdk.aws_iam.CompositeDependable",
		args,
		c,
	)
}

// Represents a principal that has multiple types of principals.
//
// A composite principal cannot
// have conditions. i.e. multiple ServicePrincipals that form a composite principal
// Experimental.
type CompositePrincipal interface {
	PrincipalBase
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddPrincipals(principals ...PrincipalBase) CompositePrincipal
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for CompositePrincipal
type jsiiProxy_CompositePrincipal struct {
	jsiiProxy_PrincipalBase
}

func (j *jsiiProxy_CompositePrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositePrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositePrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CompositePrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewCompositePrincipal(principals ...PrincipalBase) CompositePrincipal {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range principals {
		args = append(args, a)
	}

	j := jsiiProxy_CompositePrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.CompositePrincipal",
		args,
		&j,
	)

	return &j
}

// Experimental.
func NewCompositePrincipal_Override(c CompositePrincipal, principals ...PrincipalBase) {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range principals {
		args = append(args, a)
	}

	_jsii_.Create(
		"monocdk.aws_iam.CompositePrincipal",
		args,
		c,
	)
}

// Adds IAM principals to the composite principal.
//
// Composite principals cannot have
// conditions.
// Experimental.
func (c *jsiiProxy_CompositePrincipal) AddPrincipals(principals ...PrincipalBase) CompositePrincipal {
	args := []interface{}{}
	for _, a := range principals {
		args = append(args, a)
	}

	var returns CompositePrincipal

	_jsii_.Invoke(
		c,
		"addPrincipals",
		args,
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (c *jsiiProxy_CompositePrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		c,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (c *jsiiProxy_CompositePrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		c,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (c *jsiiProxy_CompositePrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		c,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (c *jsiiProxy_CompositePrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (c *jsiiProxy_CompositePrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		c,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// The Effect element of an IAM policy.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_effect.html
//
// Experimental.
type Effect string

const (
	Effect_ALLOW Effect = "ALLOW"
	Effect_DENY Effect = "DENY"
)

// Principal entity that represents a federated identity provider such as Amazon Cognito, that can be used to provide temporary security credentials to users who have been authenticated.
//
// Additional condition keys are available when the temporary security credentials are used to make a request.
// You can use these keys to write policies that limit the access of federated users.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_iam-condition-keys.html#condition-keys-wif
//
// Experimental.
type FederatedPrincipal interface {
	PrincipalBase
	AssumeRoleAction() *string
	Conditions() *map[string]interface{}
	Federated() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for FederatedPrincipal
type jsiiProxy_FederatedPrincipal struct {
	jsiiProxy_PrincipalBase
}

func (j *jsiiProxy_FederatedPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FederatedPrincipal) Conditions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FederatedPrincipal) Federated() *string {
	var returns *string
	_jsii_.Get(
		j,
		"federated",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FederatedPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FederatedPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_FederatedPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewFederatedPrincipal(federated *string, conditions *map[string]interface{}, assumeRoleAction *string) FederatedPrincipal {
	_init_.Initialize()

	j := jsiiProxy_FederatedPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.FederatedPrincipal",
		[]interface{}{federated, conditions, assumeRoleAction},
		&j,
	)

	return &j
}

// Experimental.
func NewFederatedPrincipal_Override(f FederatedPrincipal, federated *string, conditions *map[string]interface{}, assumeRoleAction *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.FederatedPrincipal",
		[]interface{}{federated, conditions, assumeRoleAction},
		f,
	)
}

// Add to the policy of this principal.
// Experimental.
func (f *jsiiProxy_FederatedPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		f,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (f *jsiiProxy_FederatedPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		f,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (f *jsiiProxy_FederatedPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		f,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (f *jsiiProxy_FederatedPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		f,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (f *jsiiProxy_FederatedPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		f,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// Options allowing customizing the behavior of {@link Role.fromRoleArn}.
// Experimental.
type FromRoleArnOptions struct {
	// Whether the imported role can be modified by attaching policy resources to it.
	// Experimental.
	Mutable *bool `json:"mutable"`
}

// Result of a grant() operation.
//
// This class is not instantiable by consumers on purpose, so that they will be
// required to call the Grant factory functions.
// Experimental.
type Grant interface {
	awscdk.IDependable
	PrincipalStatement() PolicyStatement
	ResourceStatement() PolicyStatement
	Success() *bool
	ApplyBefore(constructs ...awscdk.IConstruct)
	AssertSuccess()
}

// The jsii proxy struct for Grant
type jsiiProxy_Grant struct {
	internal.Type__awscdkIDependable
}

func (j *jsiiProxy_Grant) PrincipalStatement() PolicyStatement {
	var returns PolicyStatement
	_jsii_.Get(
		j,
		"principalStatement",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Grant) ResourceStatement() PolicyStatement {
	var returns PolicyStatement
	_jsii_.Get(
		j,
		"resourceStatement",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Grant) Success() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"success",
		&returns,
	)
	return returns
}


// Try to grant the given permissions to the given principal.
//
// Absence of a principal leads to a warning, but failing to add
// the permissions to a present principal is not an error.
// Experimental.
func Grant_AddToPrincipal(options *GrantOnPrincipalOptions) Grant {
	_init_.Initialize()

	var returns Grant

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Grant",
		"addToPrincipal",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Add a grant both on the principal and on the resource.
//
// As long as any principal is given, granting on the principal may fail (in
// case of a non-identity principal), but granting on the resource will
// never fail.
//
// Statement will be the resource statement.
// Experimental.
func Grant_AddToPrincipalAndResource(options *GrantOnPrincipalAndResourceOptions) Grant {
	_init_.Initialize()

	var returns Grant

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Grant",
		"addToPrincipalAndResource",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Grant the given permissions to the principal.
//
// The permissions will be added to the principal policy primarily, falling
// back to the resource policy if necessary. The permissions must be granted
// somewhere.
//
// - Trying to grant permissions to a principal that does not admit adding to
//    the principal policy while not providing a resource with a resource policy
//    is an error.
// - Trying to grant permissions to an absent principal (possible in the
//    case of imported resources) leads to a warning being added to the
//    resource construct.
// Experimental.
func Grant_AddToPrincipalOrResource(options *GrantWithResourceOptions) Grant {
	_init_.Initialize()

	var returns Grant

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Grant",
		"addToPrincipalOrResource",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Returns a "no-op" `Grant` object which represents a "dropped grant".
//
// This can be used for e.g. imported resources where you may not be able to modify
// the resource's policy or some underlying policy which you don't know about.
// Experimental.
func Grant_Drop(grantee IGrantable, _intent *string) Grant {
	_init_.Initialize()

	var returns Grant

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Grant",
		"drop",
		[]interface{}{grantee, _intent},
		&returns,
	)

	return returns
}

// Make sure this grant is applied before the given constructs are deployed.
//
// The same as construct.node.addDependency(grant), but slightly nicer to read.
// Experimental.
func (g *jsiiProxy_Grant) ApplyBefore(constructs ...awscdk.IConstruct) {
	args := []interface{}{}
	for _, a := range constructs {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		g,
		"applyBefore",
		args,
	)
}

// Throw an error if this grant wasn't successful.
// Experimental.
func (g *jsiiProxy_Grant) AssertSuccess() {
	_jsii_.InvokeVoid(
		g,
		"assertSuccess",
		nil, // no parameters
	)
}

// Options for a grant operation to both identity and resource.
// Experimental.
type GrantOnPrincipalAndResourceOptions struct {
	// The actions to grant.
	// Experimental.
	Actions *[]*string `json:"actions"`
	// The principal to grant to.
	// Experimental.
	Grantee IGrantable `json:"grantee"`
	// The resource ARNs to grant to.
	// Experimental.
	ResourceArns *[]*string `json:"resourceArns"`
	// The resource with a resource policy.
	//
	// The statement will always be added to the resource policy.
	// Experimental.
	Resource IResourceWithPolicy `json:"resource"`
	// The principal to use in the statement for the resource policy.
	// Experimental.
	ResourcePolicyPrincipal IPrincipal `json:"resourcePolicyPrincipal"`
	// When referring to the resource in a resource policy, use this as ARN.
	//
	// (Depending on the resource type, this needs to be '*' in a resource policy).
	// Experimental.
	ResourceSelfArns *[]*string `json:"resourceSelfArns"`
}

// Options for a grant operation that only applies to principals.
// Experimental.
type GrantOnPrincipalOptions struct {
	// The actions to grant.
	// Experimental.
	Actions *[]*string `json:"actions"`
	// The principal to grant to.
	// Experimental.
	Grantee IGrantable `json:"grantee"`
	// The resource ARNs to grant to.
	// Experimental.
	ResourceArns *[]*string `json:"resourceArns"`
	// Construct to report warnings on in case grant could not be registered.
	// Experimental.
	Scope awscdk.IConstruct `json:"scope"`
}

// Options for a grant operation.
// Experimental.
type GrantWithResourceOptions struct {
	// The actions to grant.
	// Experimental.
	Actions *[]*string `json:"actions"`
	// The principal to grant to.
	// Experimental.
	Grantee IGrantable `json:"grantee"`
	// The resource ARNs to grant to.
	// Experimental.
	ResourceArns *[]*string `json:"resourceArns"`
	// The resource with a resource policy.
	//
	// The statement will be added to the resource policy if it couldn't be
	// added to the principal policy.
	// Experimental.
	Resource IResourceWithPolicy `json:"resource"`
	// When referring to the resource in a resource policy, use this as ARN.
	//
	// (Depending on the resource type, this needs to be '*' in a resource policy).
	// Experimental.
	ResourceSelfArns *[]*string `json:"resourceSelfArns"`
}

// An IAM Group (collection of IAM users) lets you specify permissions for multiple users, which can make it easier to manage permissions for those users.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_groups.html
//
// Experimental.
type Group interface {
	awscdk.Resource
	IGroup
	AssumeRoleAction() *string
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() IPrincipal
	GroupArn() *string
	GroupName() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	Stack() awscdk.Stack
	AddManagedPolicy(policy IManagedPolicy)
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult
	AddUser(user IUser)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AttachInlinePolicy(policy Policy)
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

// The jsii proxy struct for Group
type jsiiProxy_Group struct {
	internal.Type__awscdkResource
	jsiiProxy_IGroup
}

func (j *jsiiProxy_Group) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) GroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) GroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Group) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewGroup(scope constructs.Construct, id *string, props *GroupProps) Group {
	_init_.Initialize()

	j := jsiiProxy_Group{}

	_jsii_.Create(
		"monocdk.aws_iam.Group",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGroup_Override(g Group, scope constructs.Construct, id *string, props *GroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.Group",
		[]interface{}{scope, id, props},
		g,
	)
}

// Import an external group by ARN.
//
// If the imported Group ARN is a Token (such as a
// `CfnParameter.valueAsString` or a `Fn.importValue()`) *and* the referenced
// group has a `path` (like `arn:...:group/AdminGroup/NetworkAdmin`), the
// `groupName` property will not resolve to the correct value. Instead it
// will resolve to the first path component. We unfortunately cannot express
// the correct calculation of the full path name as a CloudFormation
// expression. In this scenario the Group ARN should be supplied without the
// `path` in order to resolve the correct group resource.
// Experimental.
func Group_FromGroupArn(scope constructs.Construct, id *string, groupArn *string) IGroup {
	_init_.Initialize()

	var returns IGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Group",
		"fromGroupArn",
		[]interface{}{scope, id, groupArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Group_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Group",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Group_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Group",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Attaches a managed policy to this group.
// Experimental.
func (g *jsiiProxy_Group) AddManagedPolicy(policy IManagedPolicy) {
	_jsii_.InvokeVoid(
		g,
		"addManagedPolicy",
		[]interface{}{policy},
	)
}

// Add to the policy of this principal.
// Experimental.
func (g *jsiiProxy_Group) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		g,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Adds an IAM statement to the default policy.
// Experimental.
func (g *jsiiProxy_Group) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		g,
		"addToPrincipalPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Adds a user to this group.
// Experimental.
func (g *jsiiProxy_Group) AddUser(user IUser) {
	_jsii_.InvokeVoid(
		g,
		"addUser",
		[]interface{}{user},
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
func (g *jsiiProxy_Group) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		g,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attaches a policy to this group.
// Experimental.
func (g *jsiiProxy_Group) AttachInlinePolicy(policy Policy) {
	_jsii_.InvokeVoid(
		g,
		"attachInlinePolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (g *jsiiProxy_Group) GeneratePhysicalName() *string {
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
func (g *jsiiProxy_Group) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (g *jsiiProxy_Group) GetResourceNameAttribute(nameAttr *string) *string {
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
func (g *jsiiProxy_Group) OnPrepare() {
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
func (g *jsiiProxy_Group) OnSynthesize(session constructs.ISynthesisSession) {
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
func (g *jsiiProxy_Group) OnValidate() *[]*string {
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
func (g *jsiiProxy_Group) Prepare() {
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
func (g *jsiiProxy_Group) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		g,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (g *jsiiProxy_Group) ToString() *string {
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
func (g *jsiiProxy_Group) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		g,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining an IAM group.
// Experimental.
type GroupProps struct {
	// A name for the IAM group.
	//
	// For valid values, see the GroupName parameter
	// for the CreateGroup action in the IAM API Reference. If you don't specify
	// a name, AWS CloudFormation generates a unique physical ID and uses that
	// ID for the group name.
	//
	// If you specify a name, you must specify the CAPABILITY_NAMED_IAM value to
	// acknowledge your template's capabilities. For more information, see
	// Acknowledging IAM Resources in AWS CloudFormation Templates.
	// Experimental.
	GroupName *string `json:"groupName"`
	// A list of managed policies associated with this role.
	//
	// You can add managed policies later using
	// `addManagedPolicy(ManagedPolicy.fromAwsManagedPolicyName(policyName))`.
	// Experimental.
	ManagedPolicies *[]IManagedPolicy `json:"managedPolicies"`
	// The path to the group.
	//
	// For more information about paths, see [IAM
	// Identifiers](http://docs.aws.amazon.com/IAM/latest/UserGuide/index.html?Using_Identifiers.html)
	// in the IAM User Guide.
	// Experimental.
	Path *string `json:"path"`
}

// Any object that has an associated principal that a permission can be granted to.
// Experimental.
type IGrantable interface {
	// The principal to grant permissions to.
	// Experimental.
	GrantPrincipal() IPrincipal
}

// The jsii proxy for IGrantable
type jsiiProxy_IGrantable struct {
	_ byte // padding
}

func (j *jsiiProxy_IGrantable) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

// Represents an IAM Group.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_groups.html
//
// Experimental.
type IGroup interface {
	IIdentity
	// Returns the IAM Group ARN.
	// Experimental.
	GroupArn() *string
	// Returns the IAM Group Name.
	// Experimental.
	GroupName() *string
}

// The jsii proxy for IGroup
type jsiiProxy_IGroup struct {
	jsiiProxy_IIdentity
}

func (j *jsiiProxy_IGroup) GroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IGroup) GroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"groupName",
		&returns,
	)
	return returns
}

// A construct that represents an IAM principal, such as a user, group or role.
// Experimental.
type IIdentity interface {
	IPrincipal
	awscdk.IResource
	// Attaches a managed policy to this principal.
	// Experimental.
	AddManagedPolicy(policy IManagedPolicy)
	// Attaches an inline policy to this principal.
	//
	// This is the same as calling `policy.addToXxx(principal)`.
	// Experimental.
	AttachInlinePolicy(policy Policy)
}

// The jsii proxy for IIdentity
type jsiiProxy_IIdentity struct {
	jsiiProxy_IPrincipal
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IIdentity) AddManagedPolicy(policy IManagedPolicy) {
	_jsii_.InvokeVoid(
		i,
		"addManagedPolicy",
		[]interface{}{policy},
	)
}

func (i *jsiiProxy_IIdentity) AttachInlinePolicy(policy Policy) {
	_jsii_.InvokeVoid(
		i,
		"attachInlinePolicy",
		[]interface{}{policy},
	)
}

func (i *jsiiProxy_IIdentity) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		i,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IIdentity) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		i,
		"addToPrincipalPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IIdentity) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IIdentity) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IIdentity) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IIdentity) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IIdentity) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IIdentity) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IIdentity) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// A managed policy.
// Experimental.
type IManagedPolicy interface {
	// The ARN of the managed policy.
	// Experimental.
	ManagedPolicyArn() *string
}

// The jsii proxy for IManagedPolicy
type jsiiProxy_IManagedPolicy struct {
	_ byte // padding
}

func (j *jsiiProxy_IManagedPolicy) ManagedPolicyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"managedPolicyArn",
		&returns,
	)
	return returns
}

// Represents an IAM OpenID Connect provider.
// Experimental.
type IOpenIdConnectProvider interface {
	awscdk.IResource
	// The Amazon Resource Name (ARN) of the IAM OpenID Connect provider.
	// Experimental.
	OpenIdConnectProviderArn() *string
	// The issuer for OIDC Provider.
	// Experimental.
	OpenIdConnectProviderIssuer() *string
}

// The jsii proxy for IOpenIdConnectProvider
type jsiiProxy_IOpenIdConnectProvider struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IOpenIdConnectProvider) OpenIdConnectProviderArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"openIdConnectProviderArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IOpenIdConnectProvider) OpenIdConnectProviderIssuer() *string {
	var returns *string
	_jsii_.Get(
		j,
		"openIdConnectProviderIssuer",
		&returns,
	)
	return returns
}

// Represents an IAM Policy.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_manage.html
//
// Experimental.
type IPolicy interface {
	awscdk.IResource
	// The name of this policy.
	// Experimental.
	PolicyName() *string
}

// The jsii proxy for IPolicy
type jsiiProxy_IPolicy struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IPolicy) PolicyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"policyName",
		&returns,
	)
	return returns
}

// Represents a logical IAM principal.
//
// An IPrincipal describes a logical entity that can perform AWS API calls
// against sets of resources, optionally under certain conditions.
//
// Examples of simple principals are IAM objects that you create, such
// as Users or Roles.
//
// An example of a more complex principals is a `ServicePrincipal` (such as
// `new ServicePrincipal("sns.amazonaws.com")`, which represents the Simple
// Notifications Service).
//
// A single logical Principal may also map to a set of physical principals.
// For example, `new OrganizationPrincipal('o-1234')` represents all
// identities that are part of the given AWS Organization.
// Experimental.
type IPrincipal interface {
	IGrantable
	// Add to the policy of this principal.
	//
	// Returns: true if the statement was added, false if the principal in
	// question does not have a policy document to add the statement to.
	// Deprecated: Use `addToPrincipalPolicy` instead.
	AddToPolicy(statement PolicyStatement) *bool
	// Add to the policy of this principal.
	// Experimental.
	AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult
	// When this Principal is used in an AssumeRole policy, the action to use.
	// Experimental.
	AssumeRoleAction() *string
	// Return the policy fragment that identifies this principal in a Policy.
	// Experimental.
	PolicyFragment() PrincipalPolicyFragment
	// The AWS account ID of this principal.
	//
	// Can be undefined when the account is not known
	// (for example, for service principals).
	// Can be a Token - in that case,
	// it's assumed to be AWS::AccountId.
	// Experimental.
	PrincipalAccount() *string
}

// The jsii proxy for IPrincipal
type jsiiProxy_IPrincipal struct {
	jsiiProxy_IGrantable
}

func (i *jsiiProxy_IPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		i,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IPrincipal) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		i,
		"addToPrincipalPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}

// A resource with a resource policy that can be added to.
// Experimental.
type IResourceWithPolicy interface {
	awscdk.IResource
	// Add a statement to the resource's resource policy.
	// Experimental.
	AddToResourcePolicy(statement PolicyStatement) *AddToResourcePolicyResult
}

// The jsii proxy for IResourceWithPolicy
type jsiiProxy_IResourceWithPolicy struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IResourceWithPolicy) AddToResourcePolicy(statement PolicyStatement) *AddToResourcePolicyResult {
	var returns *AddToResourcePolicyResult

	_jsii_.Invoke(
		i,
		"addToResourcePolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// A Role object.
// Experimental.
type IRole interface {
	IIdentity
	// Grant the actions defined in actions to the identity Principal on this resource.
	// Experimental.
	Grant(grantee IPrincipal, actions ...*string) Grant
	// Grant permissions to the given principal to pass this role.
	// Experimental.
	GrantPassRole(grantee IPrincipal) Grant
	// Returns the ARN of this role.
	// Experimental.
	RoleArn() *string
	// Returns the name of this role.
	// Experimental.
	RoleName() *string
}

// The jsii proxy for IRole
type jsiiProxy_IRole struct {
	jsiiProxy_IIdentity
}

func (i *jsiiProxy_IRole) Grant(grantee IPrincipal, actions ...*string) Grant {
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns Grant

	_jsii_.Invoke(
		i,
		"grant",
		args,
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IRole) GrantPassRole(grantee IPrincipal) Grant {
	var returns Grant

	_jsii_.Invoke(
		i,
		"grantPassRole",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IRole) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRole) RoleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleName",
		&returns,
	)
	return returns
}

// A SAML provider.
// Experimental.
type ISamlProvider interface {
	awscdk.IResource
	// The Amazon Resource Name (ARN) of the provider.
	// Experimental.
	SamlProviderArn() *string
}

// The jsii proxy for ISamlProvider
type jsiiProxy_ISamlProvider struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_ISamlProvider) SamlProviderArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"samlProviderArn",
		&returns,
	)
	return returns
}

// Represents an IAM user.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users.html
//
// Experimental.
type IUser interface {
	IIdentity
	// Adds this user to a group.
	// Experimental.
	AddToGroup(group IGroup)
	// The user's ARN.
	// Experimental.
	UserArn() *string
	// The user's name.
	// Experimental.
	UserName() *string
}

// The jsii proxy for IUser
type jsiiProxy_IUser struct {
	jsiiProxy_IIdentity
}

func (i *jsiiProxy_IUser) AddToGroup(group IGroup) {
	_jsii_.InvokeVoid(
		i,
		"addToGroup",
		[]interface{}{group},
	)
}

func (j *jsiiProxy_IUser) UserArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IUser) UserName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userName",
		&returns,
	)
	return returns
}

// An IAM role that only gets attached to the construct tree once it gets used, not before.
//
// This construct can be used to simplify logic in other constructs
// which need to create a role but only if certain configurations occur
// (such as when AutoScaling is configured). The role can be configured in one
// place, but if it never gets used it doesn't get instantiated and will
// not be synthesized or deployed.
// Experimental.
type LazyRole interface {
	awscdk.Resource
	IRole
	AssumeRoleAction() *string
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() IPrincipal
	Node() awscdk.ConstructNode
	PhysicalName() *string
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	RoleArn() *string
	RoleId() *string
	RoleName() *string
	Stack() awscdk.Stack
	AddManagedPolicy(policy IManagedPolicy)
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AttachInlinePolicy(policy Policy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(identity IPrincipal, actions ...*string) Grant
	GrantPassRole(identity IPrincipal) Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for LazyRole
type jsiiProxy_LazyRole struct {
	internal.Type__awscdkResource
	jsiiProxy_IRole
}

func (j *jsiiProxy_LazyRole) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) RoleId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) RoleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LazyRole) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewLazyRole(scope constructs.Construct, id *string, props *LazyRoleProps) LazyRole {
	_init_.Initialize()

	j := jsiiProxy_LazyRole{}

	_jsii_.Create(
		"monocdk.aws_iam.LazyRole",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewLazyRole_Override(l LazyRole, scope constructs.Construct, id *string, props *LazyRoleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.LazyRole",
		[]interface{}{scope, id, props},
		l,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func LazyRole_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.LazyRole",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func LazyRole_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.LazyRole",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Attaches a managed policy to this role.
// Experimental.
func (l *jsiiProxy_LazyRole) AddManagedPolicy(policy IManagedPolicy) {
	_jsii_.InvokeVoid(
		l,
		"addManagedPolicy",
		[]interface{}{policy},
	)
}

// Add to the policy of this principal.
// Experimental.
func (l *jsiiProxy_LazyRole) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		l,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Adds a permission to the role's default policy document.
//
// If there is no default policy attached to this role, it will be created.
// Experimental.
func (l *jsiiProxy_LazyRole) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		l,
		"addToPrincipalPolicy",
		[]interface{}{statement},
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
func (l *jsiiProxy_LazyRole) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		l,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attaches a policy to this role.
// Experimental.
func (l *jsiiProxy_LazyRole) AttachInlinePolicy(policy Policy) {
	_jsii_.InvokeVoid(
		l,
		"attachInlinePolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (l *jsiiProxy_LazyRole) GeneratePhysicalName() *string {
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
func (l *jsiiProxy_LazyRole) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (l *jsiiProxy_LazyRole) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		l,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the actions defined in actions to the identity Principal on this resource.
// Experimental.
func (l *jsiiProxy_LazyRole) Grant(identity IPrincipal, actions ...*string) Grant {
	args := []interface{}{identity}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns Grant

	_jsii_.Invoke(
		l,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Grant permissions to the given principal to pass this role.
// Experimental.
func (l *jsiiProxy_LazyRole) GrantPassRole(identity IPrincipal) Grant {
	var returns Grant

	_jsii_.Invoke(
		l,
		"grantPassRole",
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
func (l *jsiiProxy_LazyRole) OnPrepare() {
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
func (l *jsiiProxy_LazyRole) OnSynthesize(session constructs.ISynthesisSession) {
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
func (l *jsiiProxy_LazyRole) OnValidate() *[]*string {
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
func (l *jsiiProxy_LazyRole) Prepare() {
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
func (l *jsiiProxy_LazyRole) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (l *jsiiProxy_LazyRole) ToString() *string {
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
func (l *jsiiProxy_LazyRole) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining a LazyRole.
// Experimental.
type LazyRoleProps struct {
	// The IAM principal (i.e. `new ServicePrincipal('sns.amazonaws.com')`) which can assume this role.
	//
	// You can later modify the assume role policy document by accessing it via
	// the `assumeRolePolicy` property.
	// Experimental.
	AssumedBy IPrincipal `json:"assumedBy"`
	// A description of the role.
	//
	// It can be up to 1000 characters long.
	// Experimental.
	Description *string `json:"description"`
	// ID that the role assumer needs to provide when assuming this role.
	//
	// If the configured and provided external IDs do not match, the
	// AssumeRole operation will fail.
	// Deprecated: see {@link externalIds}
	ExternalId *string `json:"externalId"`
	// List of IDs that the role assumer needs to provide one of when assuming this role.
	//
	// If the configured and provided external IDs do not match, the
	// AssumeRole operation will fail.
	// Experimental.
	ExternalIds *[]*string `json:"externalIds"`
	// A list of named policies to inline into this role.
	//
	// These policies will be
	// created with the role, whereas those added by ``addToPolicy`` are added
	// using a separate CloudFormation resource (allowing a way around circular
	// dependencies that could otherwise be introduced).
	// Experimental.
	InlinePolicies *map[string]PolicyDocument `json:"inlinePolicies"`
	// A list of managed policies associated with this role.
	//
	// You can add managed policies later using
	// `addManagedPolicy(ManagedPolicy.fromAwsManagedPolicyName(policyName))`.
	// Experimental.
	ManagedPolicies *[]IManagedPolicy `json:"managedPolicies"`
	// The maximum session duration that you want to set for the specified role.
	//
	// This setting can have a value from 1 hour (3600sec) to 12 (43200sec) hours.
	//
	// Anyone who assumes the role from the AWS CLI or API can use the
	// DurationSeconds API parameter or the duration-seconds CLI parameter to
	// request a longer session. The MaxSessionDuration setting determines the
	// maximum duration that can be requested using the DurationSeconds
	// parameter.
	//
	// If users don't specify a value for the DurationSeconds parameter, their
	// security credentials are valid for one hour by default. This applies when
	// you use the AssumeRole* API operations or the assume-role* CLI operations
	// but does not apply when you use those operations to create a console URL.
	// Experimental.
	MaxSessionDuration awscdk.Duration `json:"maxSessionDuration"`
	// The path associated with this role.
	//
	// For information about IAM paths, see
	// Friendly Names and Paths in IAM User Guide.
	// Experimental.
	Path *string `json:"path"`
	// AWS supports permissions boundaries for IAM entities (users or roles).
	//
	// A permissions boundary is an advanced feature for using a managed policy
	// to set the maximum permissions that an identity-based policy can grant to
	// an IAM entity. An entity's permissions boundary allows it to perform only
	// the actions that are allowed by both its identity-based policies and its
	// permissions boundaries.
	// Experimental.
	PermissionsBoundary IManagedPolicy `json:"permissionsBoundary"`
	// A name for the IAM role.
	//
	// For valid values, see the RoleName parameter for
	// the CreateRole action in the IAM API Reference.
	//
	// IMPORTANT: If you specify a name, you cannot perform updates that require
	// replacement of this resource. You can perform updates that require no or
	// some interruption. If you must replace the resource, specify a new name.
	//
	// If you specify a name, you must specify the CAPABILITY_NAMED_IAM value to
	// acknowledge your template's capabilities. For more information, see
	// Acknowledging IAM Resources in AWS CloudFormation Templates.
	// Experimental.
	RoleName *string `json:"roleName"`
}

// Managed policy.
// Experimental.
type ManagedPolicy interface {
	awscdk.Resource
	IManagedPolicy
	Description() *string
	Document() PolicyDocument
	Env() *awscdk.ResourceEnvironment
	ManagedPolicyArn() *string
	ManagedPolicyName() *string
	Node() awscdk.ConstructNode
	Path() *string
	PhysicalName() *string
	Stack() awscdk.Stack
	AddStatements(statement ...PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AttachToGroup(group IGroup)
	AttachToRole(role IRole)
	AttachToUser(user IUser)
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

// The jsii proxy struct for ManagedPolicy
type jsiiProxy_ManagedPolicy struct {
	internal.Type__awscdkResource
	jsiiProxy_IManagedPolicy
}

func (j *jsiiProxy_ManagedPolicy) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) Document() PolicyDocument {
	var returns PolicyDocument
	_jsii_.Get(
		j,
		"document",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) ManagedPolicyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"managedPolicyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) ManagedPolicyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"managedPolicyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ManagedPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewManagedPolicy(scope constructs.Construct, id *string, props *ManagedPolicyProps) ManagedPolicy {
	_init_.Initialize()

	j := jsiiProxy_ManagedPolicy{}

	_jsii_.Create(
		"monocdk.aws_iam.ManagedPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewManagedPolicy_Override(m ManagedPolicy, scope constructs.Construct, id *string, props *ManagedPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.ManagedPolicy",
		[]interface{}{scope, id, props},
		m,
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
func ManagedPolicy_FromAwsManagedPolicyName(managedPolicyName *string) IManagedPolicy {
	_init_.Initialize()

	var returns IManagedPolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.ManagedPolicy",
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
func ManagedPolicy_FromManagedPolicyArn(scope constructs.Construct, id *string, managedPolicyArn *string) IManagedPolicy {
	_init_.Initialize()

	var returns IManagedPolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.ManagedPolicy",
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
func ManagedPolicy_FromManagedPolicyName(scope constructs.Construct, id *string, managedPolicyName *string) IManagedPolicy {
	_init_.Initialize()

	var returns IManagedPolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.ManagedPolicy",
		"fromManagedPolicyName",
		[]interface{}{scope, id, managedPolicyName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ManagedPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.ManagedPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ManagedPolicy_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.ManagedPolicy",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a statement to the policy document.
// Experimental.
func (m *jsiiProxy_ManagedPolicy) AddStatements(statement ...PolicyStatement) {
	args := []interface{}{}
	for _, a := range statement {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		m,
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
func (m *jsiiProxy_ManagedPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		m,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attaches this policy to a group.
// Experimental.
func (m *jsiiProxy_ManagedPolicy) AttachToGroup(group IGroup) {
	_jsii_.InvokeVoid(
		m,
		"attachToGroup",
		[]interface{}{group},
	)
}

// Attaches this policy to a role.
// Experimental.
func (m *jsiiProxy_ManagedPolicy) AttachToRole(role IRole) {
	_jsii_.InvokeVoid(
		m,
		"attachToRole",
		[]interface{}{role},
	)
}

// Attaches this policy to a user.
// Experimental.
func (m *jsiiProxy_ManagedPolicy) AttachToUser(user IUser) {
	_jsii_.InvokeVoid(
		m,
		"attachToUser",
		[]interface{}{user},
	)
}

// Experimental.
func (m *jsiiProxy_ManagedPolicy) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		m,
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
func (m *jsiiProxy_ManagedPolicy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		m,
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
func (m *jsiiProxy_ManagedPolicy) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		m,
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
func (m *jsiiProxy_ManagedPolicy) OnPrepare() {
	_jsii_.InvokeVoid(
		m,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (m *jsiiProxy_ManagedPolicy) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		m,
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
func (m *jsiiProxy_ManagedPolicy) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		m,
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
func (m *jsiiProxy_ManagedPolicy) Prepare() {
	_jsii_.InvokeVoid(
		m,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (m *jsiiProxy_ManagedPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		m,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (m *jsiiProxy_ManagedPolicy) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		m,
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
func (m *jsiiProxy_ManagedPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		m,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining an IAM managed policy.
// Experimental.
type ManagedPolicyProps struct {
	// A description of the managed policy.
	//
	// Typically used to store information about the
	// permissions defined in the policy. For example, "Grants access to production DynamoDB tables."
	// The policy description is immutable. After a value is assigned, it cannot be changed.
	// Experimental.
	Description *string `json:"description"`
	// Initial PolicyDocument to use for this ManagedPolicy.
	//
	// If omited, any
	// `PolicyStatement` provided in the `statements` property will be applied
	// against the empty default `PolicyDocument`.
	// Experimental.
	Document PolicyDocument `json:"document"`
	// Groups to attach this policy to.
	//
	// You can also use `attachToGroup(group)` to attach this policy to a group.
	// Experimental.
	Groups *[]IGroup `json:"groups"`
	// The name of the managed policy.
	//
	// If you specify multiple policies for an entity,
	// specify unique names. For example, if you specify a list of policies for
	// an IAM role, each policy must have a unique name.
	// Experimental.
	ManagedPolicyName *string `json:"managedPolicyName"`
	// The path for the policy.
	//
	// This parameter allows (through its regex pattern) a string of characters
	// consisting of either a forward slash (/) by itself or a string that must begin and end with forward slashes.
	// In addition, it can contain any ASCII character from the ! (\u0021) through the DEL character (\u007F),
	// including most punctuation characters, digits, and upper and lowercased letters.
	//
	// For more information about paths, see IAM Identifiers in the IAM User Guide.
	// Experimental.
	Path *string `json:"path"`
	// Roles to attach this policy to.
	//
	// You can also use `attachToRole(role)` to attach this policy to a role.
	// Experimental.
	Roles *[]IRole `json:"roles"`
	// Initial set of permissions to add to this policy document.
	//
	// You can also use `addPermission(statement)` to add permissions later.
	// Experimental.
	Statements *[]PolicyStatement `json:"statements"`
	// Users to attach this policy to.
	//
	// You can also use `attachToUser(user)` to attach this policy to a user.
	// Experimental.
	Users *[]IUser `json:"users"`
}

// A principal that represents a federated identity provider as from a OpenID Connect provider.
// Experimental.
type OpenIdConnectPrincipal interface {
	WebIdentityPrincipal
	AssumeRoleAction() *string
	Conditions() *map[string]interface{}
	Federated() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for OpenIdConnectPrincipal
type jsiiProxy_OpenIdConnectPrincipal struct {
	jsiiProxy_WebIdentityPrincipal
}

func (j *jsiiProxy_OpenIdConnectPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectPrincipal) Conditions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectPrincipal) Federated() *string {
	var returns *string
	_jsii_.Get(
		j,
		"federated",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewOpenIdConnectPrincipal(openIdConnectProvider IOpenIdConnectProvider, conditions *map[string]interface{}) OpenIdConnectPrincipal {
	_init_.Initialize()

	j := jsiiProxy_OpenIdConnectPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.OpenIdConnectPrincipal",
		[]interface{}{openIdConnectProvider, conditions},
		&j,
	)

	return &j
}

// Experimental.
func NewOpenIdConnectPrincipal_Override(o OpenIdConnectPrincipal, openIdConnectProvider IOpenIdConnectProvider, conditions *map[string]interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.OpenIdConnectPrincipal",
		[]interface{}{openIdConnectProvider, conditions},
		o,
	)
}

// Add to the policy of this principal.
// Experimental.
func (o *jsiiProxy_OpenIdConnectPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		o,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (o *jsiiProxy_OpenIdConnectPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		o,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (o *jsiiProxy_OpenIdConnectPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		o,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (o *jsiiProxy_OpenIdConnectPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		o,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (o *jsiiProxy_OpenIdConnectPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		o,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// IAM OIDC identity providers are entities in IAM that describe an external identity provider (IdP) service that supports the OpenID Connect (OIDC) standard, such as Google or Salesforce.
//
// You use an IAM OIDC identity provider
// when you want to establish trust between an OIDC-compatible IdP and your AWS
// account. This is useful when creating a mobile app or web application that
// requires access to AWS resources, but you don't want to create custom sign-in
// code or manage your own user identities.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_oidc.html
//
// Experimental.
type OpenIdConnectProvider interface {
	awscdk.Resource
	IOpenIdConnectProvider
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	OpenIdConnectProviderArn() *string
	OpenIdConnectProviderIssuer() *string
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

// The jsii proxy struct for OpenIdConnectProvider
type jsiiProxy_OpenIdConnectProvider struct {
	internal.Type__awscdkResource
	jsiiProxy_IOpenIdConnectProvider
}

func (j *jsiiProxy_OpenIdConnectProvider) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectProvider) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectProvider) OpenIdConnectProviderArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"openIdConnectProviderArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectProvider) OpenIdConnectProviderIssuer() *string {
	var returns *string
	_jsii_.Get(
		j,
		"openIdConnectProviderIssuer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectProvider) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OpenIdConnectProvider) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Defines an OpenID Connect provider.
// Experimental.
func NewOpenIdConnectProvider(scope constructs.Construct, id *string, props *OpenIdConnectProviderProps) OpenIdConnectProvider {
	_init_.Initialize()

	j := jsiiProxy_OpenIdConnectProvider{}

	_jsii_.Create(
		"monocdk.aws_iam.OpenIdConnectProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Defines an OpenID Connect provider.
// Experimental.
func NewOpenIdConnectProvider_Override(o OpenIdConnectProvider, scope constructs.Construct, id *string, props *OpenIdConnectProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.OpenIdConnectProvider",
		[]interface{}{scope, id, props},
		o,
	)
}

// Imports an Open ID connect provider from an ARN.
// Experimental.
func OpenIdConnectProvider_FromOpenIdConnectProviderArn(scope constructs.Construct, id *string, openIdConnectProviderArn *string) IOpenIdConnectProvider {
	_init_.Initialize()

	var returns IOpenIdConnectProvider

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.OpenIdConnectProvider",
		"fromOpenIdConnectProviderArn",
		[]interface{}{scope, id, openIdConnectProviderArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func OpenIdConnectProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.OpenIdConnectProvider",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func OpenIdConnectProvider_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.OpenIdConnectProvider",
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
func (o *jsiiProxy_OpenIdConnectProvider) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		o,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (o *jsiiProxy_OpenIdConnectProvider) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		o,
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
func (o *jsiiProxy_OpenIdConnectProvider) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		o,
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
func (o *jsiiProxy_OpenIdConnectProvider) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		o,
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
func (o *jsiiProxy_OpenIdConnectProvider) OnPrepare() {
	_jsii_.InvokeVoid(
		o,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (o *jsiiProxy_OpenIdConnectProvider) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		o,
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
func (o *jsiiProxy_OpenIdConnectProvider) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		o,
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
func (o *jsiiProxy_OpenIdConnectProvider) Prepare() {
	_jsii_.InvokeVoid(
		o,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (o *jsiiProxy_OpenIdConnectProvider) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		o,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (o *jsiiProxy_OpenIdConnectProvider) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		o,
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
func (o *jsiiProxy_OpenIdConnectProvider) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		o,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Initialization properties for `OpenIdConnectProvider`.
// Experimental.
type OpenIdConnectProviderProps struct {
	// The URL of the identity provider.
	//
	// The URL must begin with https:// and
	// should correspond to the iss claim in the provider's OpenID Connect ID
	// tokens. Per the OIDC standard, path components are allowed but query
	// parameters are not. Typically the URL consists of only a hostname, like
	// https://server.example.org or https://example.com.
	//
	// You cannot register the same provider multiple times in a single AWS
	// account. If you try to submit a URL that has already been used for an
	// OpenID Connect provider in the AWS account, you will get an error.
	// Experimental.
	Url *string `json:"url"`
	// A list of client IDs (also known as audiences).
	//
	// When a mobile or web app
	// registers with an OpenID Connect provider, they establish a value that
	// identifies the application. (This is the value that's sent as the client_id
	// parameter on OAuth requests.)
	//
	// You can register multiple client IDs with the same provider. For example,
	// you might have multiple applications that use the same OIDC provider. You
	// cannot register more than 100 client IDs with a single IAM OIDC provider.
	//
	// Client IDs are up to 255 characters long.
	// Experimental.
	ClientIds *[]*string `json:"clientIds"`
	// A list of server certificate thumbprints for the OpenID Connect (OIDC) identity provider's server certificates.
	//
	// Typically this list includes only one entry. However, IAM lets you have up
	// to five thumbprints for an OIDC provider. This lets you maintain multiple
	// thumbprints if the identity provider is rotating certificates.
	//
	// The server certificate thumbprint is the hex-encoded SHA-1 hash value of
	// the X.509 certificate used by the domain where the OpenID Connect provider
	// makes its keys available. It is always a 40-character string.
	//
	// You must provide at least one thumbprint when creating an IAM OIDC
	// provider. For example, assume that the OIDC provider is server.example.com
	// and the provider stores its keys at
	// https://keys.server.example.com/openid-connect. In that case, the
	// thumbprint string would be the hex-encoded SHA-1 hash value of the
	// certificate used by https://keys.server.example.com.
	// Experimental.
	Thumbprints *[]*string `json:"thumbprints"`
}

// A principal that represents an AWS Organization.
// Experimental.
type OrganizationPrincipal interface {
	PrincipalBase
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	OrganizationId() *string
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for OrganizationPrincipal
type jsiiProxy_OrganizationPrincipal struct {
	jsiiProxy_PrincipalBase
}

func (j *jsiiProxy_OrganizationPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OrganizationPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OrganizationPrincipal) OrganizationId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"organizationId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OrganizationPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_OrganizationPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewOrganizationPrincipal(organizationId *string) OrganizationPrincipal {
	_init_.Initialize()

	j := jsiiProxy_OrganizationPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.OrganizationPrincipal",
		[]interface{}{organizationId},
		&j,
	)

	return &j
}

// Experimental.
func NewOrganizationPrincipal_Override(o OrganizationPrincipal, organizationId *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.OrganizationPrincipal",
		[]interface{}{organizationId},
		o,
	)
}

// Add to the policy of this principal.
// Experimental.
func (o *jsiiProxy_OrganizationPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		o,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (o *jsiiProxy_OrganizationPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		o,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (o *jsiiProxy_OrganizationPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		o,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (o *jsiiProxy_OrganizationPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		o,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (o *jsiiProxy_OrganizationPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		o,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// Modify the Permissions Boundaries of Users and Roles in a construct tree.
//
// TODO: EXAMPLE
//
// Experimental.
type PermissionsBoundary interface {
	Apply(boundaryPolicy IManagedPolicy)
	Clear()
}

// The jsii proxy struct for PermissionsBoundary
type jsiiProxy_PermissionsBoundary struct {
	_ byte // padding
}

// Access the Permissions Boundaries of a construct tree.
// Experimental.
func PermissionsBoundary_Of(scope constructs.IConstruct) PermissionsBoundary {
	_init_.Initialize()

	var returns PermissionsBoundary

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.PermissionsBoundary",
		"of",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Apply the given policy as Permissions Boundary to all Roles and Users in the scope.
//
// Will override any Permissions Boundaries configured previously; in case
// a Permission Boundary is applied in multiple scopes, the Boundary applied
// closest to the Role wins.
// Experimental.
func (p *jsiiProxy_PermissionsBoundary) Apply(boundaryPolicy IManagedPolicy) {
	_jsii_.InvokeVoid(
		p,
		"apply",
		[]interface{}{boundaryPolicy},
	)
}

// Remove previously applied Permissions Boundaries.
// Experimental.
func (p *jsiiProxy_PermissionsBoundary) Clear() {
	_jsii_.InvokeVoid(
		p,
		"clear",
		nil, // no parameters
	)
}

// The AWS::IAM::Policy resource associates an IAM policy with IAM users, roles, or groups.
//
// For more information about IAM policies, see [Overview of IAM
// Policies](http://docs.aws.amazon.com/IAM/latest/UserGuide/policies_overview.html)
// in the IAM User Guide guide.
// Experimental.
type Policy interface {
	awscdk.Resource
	IPolicy
	Document() PolicyDocument
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	PolicyName() *string
	Stack() awscdk.Stack
	AddStatements(statement ...PolicyStatement)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AttachToGroup(group IGroup)
	AttachToRole(role IRole)
	AttachToUser(user IUser)
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

// The jsii proxy struct for Policy
type jsiiProxy_Policy struct {
	internal.Type__awscdkResource
	jsiiProxy_IPolicy
}

func (j *jsiiProxy_Policy) Document() PolicyDocument {
	var returns PolicyDocument
	_jsii_.Get(
		j,
		"document",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Policy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Policy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Policy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Policy) PolicyName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"policyName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Policy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewPolicy(scope constructs.Construct, id *string, props *PolicyProps) Policy {
	_init_.Initialize()

	j := jsiiProxy_Policy{}

	_jsii_.Create(
		"monocdk.aws_iam.Policy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewPolicy_Override(p Policy, scope constructs.Construct, id *string, props *PolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.Policy",
		[]interface{}{scope, id, props},
		p,
	)
}

// Import a policy in this app based on its name.
// Experimental.
func Policy_FromPolicyName(scope constructs.Construct, id *string, policyName *string) IPolicy {
	_init_.Initialize()

	var returns IPolicy

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Policy",
		"fromPolicyName",
		[]interface{}{scope, id, policyName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Policy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Policy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Policy_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Policy",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a statement to the policy document.
// Experimental.
func (p *jsiiProxy_Policy) AddStatements(statement ...PolicyStatement) {
	args := []interface{}{}
	for _, a := range statement {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
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
func (p *jsiiProxy_Policy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		p,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attaches this policy to a group.
// Experimental.
func (p *jsiiProxy_Policy) AttachToGroup(group IGroup) {
	_jsii_.InvokeVoid(
		p,
		"attachToGroup",
		[]interface{}{group},
	)
}

// Attaches this policy to a role.
// Experimental.
func (p *jsiiProxy_Policy) AttachToRole(role IRole) {
	_jsii_.InvokeVoid(
		p,
		"attachToRole",
		[]interface{}{role},
	)
}

// Attaches this policy to a user.
// Experimental.
func (p *jsiiProxy_Policy) AttachToUser(user IUser) {
	_jsii_.InvokeVoid(
		p,
		"attachToUser",
		[]interface{}{user},
	)
}

// Experimental.
func (p *jsiiProxy_Policy) GeneratePhysicalName() *string {
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
func (p *jsiiProxy_Policy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (p *jsiiProxy_Policy) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		p,
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
func (p *jsiiProxy_Policy) OnPrepare() {
	_jsii_.InvokeVoid(
		p,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_Policy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (p *jsiiProxy_Policy) OnValidate() *[]*string {
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
func (p *jsiiProxy_Policy) Prepare() {
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
func (p *jsiiProxy_Policy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_Policy) ToString() *string {
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
func (p *jsiiProxy_Policy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A PolicyDocument is a collection of statements.
// Experimental.
type PolicyDocument interface {
	awscdk.IResolvable
	CreationStack() *[]*string
	IsEmpty() *bool
	StatementCount() *float64
	AddStatements(statement ...PolicyStatement)
	Resolve(context awscdk.IResolveContext) interface{}
	ToJSON() interface{}
	ToString() *string
	ValidateForAnyPolicy() *[]*string
	ValidateForIdentityPolicy() *[]*string
	ValidateForResourcePolicy() *[]*string
}

// The jsii proxy struct for PolicyDocument
type jsiiProxy_PolicyDocument struct {
	internal.Type__awscdkIResolvable
}

func (j *jsiiProxy_PolicyDocument) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PolicyDocument) IsEmpty() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isEmpty",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PolicyDocument) StatementCount() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"statementCount",
		&returns,
	)
	return returns
}


// Experimental.
func NewPolicyDocument(props *PolicyDocumentProps) PolicyDocument {
	_init_.Initialize()

	j := jsiiProxy_PolicyDocument{}

	_jsii_.Create(
		"monocdk.aws_iam.PolicyDocument",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewPolicyDocument_Override(p PolicyDocument, props *PolicyDocumentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.PolicyDocument",
		[]interface{}{props},
		p,
	)
}

// Creates a new PolicyDocument based on the object provided.
//
// This will accept an object created from the `.toJSON()` call
// Experimental.
func PolicyDocument_FromJson(obj interface{}) PolicyDocument {
	_init_.Initialize()

	var returns PolicyDocument

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.PolicyDocument",
		"fromJson",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Adds a statement to the policy document.
// Experimental.
func (p *jsiiProxy_PolicyDocument) AddStatements(statement ...PolicyStatement) {
	args := []interface{}{}
	for _, a := range statement {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
		"addStatements",
		args,
	)
}

// Produce the Token's value at resolution time.
// Experimental.
func (p *jsiiProxy_PolicyDocument) Resolve(context awscdk.IResolveContext) interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"resolve",
		[]interface{}{context},
		&returns,
	)

	return returns
}

// JSON-ify the document.
//
// Used when JSON.stringify() is called
// Experimental.
func (p *jsiiProxy_PolicyDocument) ToJSON() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Encode the policy document as a string.
// Experimental.
func (p *jsiiProxy_PolicyDocument) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate that all policy statements in the policy document satisfies the requirements for any policy.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html#access_policies-json
//
// Experimental.
func (p *jsiiProxy_PolicyDocument) ValidateForAnyPolicy() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validateForAnyPolicy",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate that all policy statements in the policy document satisfies the requirements for an identity-based policy.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html#access_policies-json
//
// Experimental.
func (p *jsiiProxy_PolicyDocument) ValidateForIdentityPolicy() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validateForIdentityPolicy",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate that all policy statements in the policy document satisfies the requirements for a resource-based policy.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html#access_policies-json
//
// Experimental.
func (p *jsiiProxy_PolicyDocument) ValidateForResourcePolicy() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validateForResourcePolicy",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a new PolicyDocument.
// Experimental.
type PolicyDocumentProps struct {
	// Automatically assign Statement Ids to all statements.
	// Experimental.
	AssignSids *bool `json:"assignSids"`
	// Initial statements to add to the policy document.
	// Experimental.
	Statements *[]PolicyStatement `json:"statements"`
}

// Properties for defining an IAM inline policy document.
// Experimental.
type PolicyProps struct {
	// Initial PolicyDocument to use for this Policy.
	//
	// If omited, any
	// `PolicyStatement` provided in the `statements` property will be applied
	// against the empty default `PolicyDocument`.
	// Experimental.
	Document PolicyDocument `json:"document"`
	// Force creation of an `AWS::IAM::Policy`.
	//
	// Unless set to `true`, this `Policy` construct will not materialize to an
	// `AWS::IAM::Policy` CloudFormation resource in case it would have no effect
	// (for example, if it remains unattached to an IAM identity or if it has no
	// statements). This is generally desired behavior, since it prevents
	// creating invalid--and hence undeployable--CloudFormation templates.
	//
	// In cases where you know the policy must be created and it is actually
	// an error if no statements have been added to it, you can set this to `true`.
	// Experimental.
	Force *bool `json:"force"`
	// Groups to attach this policy to.
	//
	// You can also use `attachToGroup(group)` to attach this policy to a group.
	// Experimental.
	Groups *[]IGroup `json:"groups"`
	// The name of the policy.
	//
	// If you specify multiple policies for an entity,
	// specify unique names. For example, if you specify a list of policies for
	// an IAM role, each policy must have a unique name.
	// Experimental.
	PolicyName *string `json:"policyName"`
	// Roles to attach this policy to.
	//
	// You can also use `attachToRole(role)` to attach this policy to a role.
	// Experimental.
	Roles *[]IRole `json:"roles"`
	// Initial set of permissions to add to this policy document.
	//
	// You can also use `addStatements(...statement)` to add permissions later.
	// Experimental.
	Statements *[]PolicyStatement `json:"statements"`
	// Users to attach this policy to.
	//
	// You can also use `attachToUser(user)` to attach this policy to a user.
	// Experimental.
	Users *[]IUser `json:"users"`
}

// Represents a statement in an IAM policy document.
// Experimental.
type PolicyStatement interface {
	Effect() Effect
	SetEffect(val Effect)
	HasPrincipal() *bool
	HasResource() *bool
	Sid() *string
	SetSid(val *string)
	AddAccountCondition(accountId *string)
	AddAccountRootPrincipal()
	AddActions(actions ...*string)
	AddAllResources()
	AddAnyPrincipal()
	AddArnPrincipal(arn *string)
	AddAwsAccountPrincipal(accountId *string)
	AddCanonicalUserPrincipal(canonicalUserId *string)
	AddCondition(key *string, value interface{})
	AddConditions(conditions *map[string]interface{})
	AddFederatedPrincipal(federated interface{}, conditions *map[string]interface{})
	AddNotActions(notActions ...*string)
	AddNotPrincipals(notPrincipals ...IPrincipal)
	AddNotResources(arns ...*string)
	AddPrincipals(principals ...IPrincipal)
	AddResources(arns ...*string)
	AddServicePrincipal(service *string, opts *ServicePrincipalOpts)
	ToJSON() interface{}
	ToStatementJson() interface{}
	ToString() *string
	ValidateForAnyPolicy() *[]*string
	ValidateForIdentityPolicy() *[]*string
	ValidateForResourcePolicy() *[]*string
}

// The jsii proxy struct for PolicyStatement
type jsiiProxy_PolicyStatement struct {
	_ byte // padding
}

func (j *jsiiProxy_PolicyStatement) Effect() Effect {
	var returns Effect
	_jsii_.Get(
		j,
		"effect",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PolicyStatement) HasPrincipal() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"hasPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PolicyStatement) HasResource() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"hasResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PolicyStatement) Sid() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sid",
		&returns,
	)
	return returns
}


// Experimental.
func NewPolicyStatement(props *PolicyStatementProps) PolicyStatement {
	_init_.Initialize()

	j := jsiiProxy_PolicyStatement{}

	_jsii_.Create(
		"monocdk.aws_iam.PolicyStatement",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewPolicyStatement_Override(p PolicyStatement, props *PolicyStatementProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.PolicyStatement",
		[]interface{}{props},
		p,
	)
}

func (j *jsiiProxy_PolicyStatement) SetEffect(val Effect) {
	_jsii_.Set(
		j,
		"effect",
		val,
	)
}

func (j *jsiiProxy_PolicyStatement) SetSid(val *string) {
	_jsii_.Set(
		j,
		"sid",
		val,
	)
}

// Creates a new PolicyStatement based on the object provided.
//
// This will accept an object created from the `.toJSON()` call
// Experimental.
func PolicyStatement_FromJson(obj interface{}) PolicyStatement {
	_init_.Initialize()

	var returns PolicyStatement

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.PolicyStatement",
		"fromJson",
		[]interface{}{obj},
		&returns,
	)

	return returns
}

// Add a condition that limits to a given account.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddAccountCondition(accountId *string) {
	_jsii_.InvokeVoid(
		p,
		"addAccountCondition",
		[]interface{}{accountId},
	)
}

// Adds an AWS account root user principal to this policy statement.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddAccountRootPrincipal() {
	_jsii_.InvokeVoid(
		p,
		"addAccountRootPrincipal",
		nil, // no parameters
	)
}

// Specify allowed actions into the "Action" section of the policy statement.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_action.html
//
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddActions(actions ...*string) {
	args := []interface{}{}
	for _, a := range actions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
		"addActions",
		args,
	)
}

// Adds a ``"*"`` resource to this statement.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddAllResources() {
	_jsii_.InvokeVoid(
		p,
		"addAllResources",
		nil, // no parameters
	)
}

// Adds all identities in all accounts ("*") to this policy statement.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddAnyPrincipal() {
	_jsii_.InvokeVoid(
		p,
		"addAnyPrincipal",
		nil, // no parameters
	)
}

// Specify a principal using the ARN  identifier of the principal.
//
// You cannot specify IAM groups and instance profiles as principals.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddArnPrincipal(arn *string) {
	_jsii_.InvokeVoid(
		p,
		"addArnPrincipal",
		[]interface{}{arn},
	)
}

// Specify AWS account ID as the principal entity to the "Principal" section of a policy statement.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddAwsAccountPrincipal(accountId *string) {
	_jsii_.InvokeVoid(
		p,
		"addAwsAccountPrincipal",
		[]interface{}{accountId},
	)
}

// Adds a canonical user ID principal to this policy document.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddCanonicalUserPrincipal(canonicalUserId *string) {
	_jsii_.InvokeVoid(
		p,
		"addCanonicalUserPrincipal",
		[]interface{}{canonicalUserId},
	)
}

// Add a condition to the Policy.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddCondition(key *string, value interface{}) {
	_jsii_.InvokeVoid(
		p,
		"addCondition",
		[]interface{}{key, value},
	)
}

// Add multiple conditions to the Policy.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddConditions(conditions *map[string]interface{}) {
	_jsii_.InvokeVoid(
		p,
		"addConditions",
		[]interface{}{conditions},
	)
}

// Adds a federated identity provider such as Amazon Cognito to this policy statement.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddFederatedPrincipal(federated interface{}, conditions *map[string]interface{}) {
	_jsii_.InvokeVoid(
		p,
		"addFederatedPrincipal",
		[]interface{}{federated, conditions},
	)
}

// Explicitly allow all actions except the specified list of actions into the "NotAction" section of the policy document.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_notaction.html
//
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddNotActions(notActions ...*string) {
	args := []interface{}{}
	for _, a := range notActions {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
		"addNotActions",
		args,
	)
}

// Specify principals that is not allowed or denied access to the "NotPrincipal" section of a policy statement.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_notprincipal.html
//
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddNotPrincipals(notPrincipals ...IPrincipal) {
	args := []interface{}{}
	for _, a := range notPrincipals {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
		"addNotPrincipals",
		args,
	)
}

// Specify resources that this policy statement will not apply to in the "NotResource" section of this policy statement.
//
// All resources except the specified list will be matched.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_notresource.html
//
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddNotResources(arns ...*string) {
	args := []interface{}{}
	for _, a := range arns {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
		"addNotResources",
		args,
	)
}

// Adds principals to the "Principal" section of a policy statement.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_principal.html
//
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddPrincipals(principals ...IPrincipal) {
	args := []interface{}{}
	for _, a := range principals {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
		"addPrincipals",
		args,
	)
}

// Specify resources that this policy statement applies into the "Resource" section of this policy statement.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_resource.html
//
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddResources(arns ...*string) {
	args := []interface{}{}
	for _, a := range arns {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		p,
		"addResources",
		args,
	)
}

// Adds a service principal to this policy statement.
// Experimental.
func (p *jsiiProxy_PolicyStatement) AddServicePrincipal(service *string, opts *ServicePrincipalOpts) {
	_jsii_.InvokeVoid(
		p,
		"addServicePrincipal",
		[]interface{}{service, opts},
	)
}

// JSON-ify the statement.
//
// Used when JSON.stringify() is called
// Experimental.
func (p *jsiiProxy_PolicyStatement) ToJSON() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// JSON-ify the policy statement.
//
// Used when JSON.stringify() is called
// Experimental.
func (p *jsiiProxy_PolicyStatement) ToStatementJson() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		p,
		"toStatementJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

// String representation of this policy statement.
// Experimental.
func (p *jsiiProxy_PolicyStatement) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate that the policy statement satisfies base requirements for a policy.
// Experimental.
func (p *jsiiProxy_PolicyStatement) ValidateForAnyPolicy() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validateForAnyPolicy",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate that the policy statement satisfies all requirements for an identity-based policy.
// Experimental.
func (p *jsiiProxy_PolicyStatement) ValidateForIdentityPolicy() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validateForIdentityPolicy",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate that the policy statement satisfies all requirements for a resource-based policy.
// Experimental.
func (p *jsiiProxy_PolicyStatement) ValidateForResourcePolicy() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validateForResourcePolicy",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Interface for creating a policy statement.
// Experimental.
type PolicyStatementProps struct {
	// List of actions to add to the statement.
	// Experimental.
	Actions *[]*string `json:"actions"`
	// Conditions to add to the statement.
	// Experimental.
	Conditions *map[string]interface{} `json:"conditions"`
	// Whether to allow or deny the actions in this statement.
	// Experimental.
	Effect Effect `json:"effect"`
	// List of not actions to add to the statement.
	// Experimental.
	NotActions *[]*string `json:"notActions"`
	// List of not principals to add to the statement.
	// Experimental.
	NotPrincipals *[]IPrincipal `json:"notPrincipals"`
	// NotResource ARNs to add to the statement.
	// Experimental.
	NotResources *[]*string `json:"notResources"`
	// List of principals to add to the statement.
	// Experimental.
	Principals *[]IPrincipal `json:"principals"`
	// Resource ARNs to add to the statement.
	// Experimental.
	Resources *[]*string `json:"resources"`
	// The Sid (statement ID) is an optional identifier that you provide for the policy statement.
	//
	// You can assign a Sid value to each statement in a
	// statement array. In services that let you specify an ID element, such as
	// SQS and SNS, the Sid value is just a sub-ID of the policy document's ID. In
	// IAM, the Sid value must be unique within a JSON policy.
	// Experimental.
	Sid *string `json:"sid"`
}

// Base class for policy principals.
// Experimental.
type PrincipalBase interface {
	IPrincipal
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for PrincipalBase
type jsiiProxy_PrincipalBase struct {
	jsiiProxy_IPrincipal
}

func (j *jsiiProxy_PrincipalBase) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PrincipalBase) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PrincipalBase) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PrincipalBase) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewPrincipalBase_Override(p PrincipalBase) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.PrincipalBase",
		nil, // no parameters
		p,
	)
}

// Add to the policy of this principal.
// Experimental.
func (p *jsiiProxy_PrincipalBase) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		p,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (p *jsiiProxy_PrincipalBase) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		p,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (p *jsiiProxy_PrincipalBase) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		p,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (p *jsiiProxy_PrincipalBase) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (p *jsiiProxy_PrincipalBase) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		p,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// A collection of the fields in a PolicyStatement that can be used to identify a principal.
//
// This consists of the JSON used in the "Principal" field, and optionally a
// set of "Condition"s that need to be applied to the policy.
// Experimental.
type PrincipalPolicyFragment interface {
	Conditions() *map[string]interface{}
	PrincipalJson() *map[string]*[]*string
}

// The jsii proxy struct for PrincipalPolicyFragment
type jsiiProxy_PrincipalPolicyFragment struct {
	_ byte // padding
}

func (j *jsiiProxy_PrincipalPolicyFragment) Conditions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PrincipalPolicyFragment) PrincipalJson() *map[string]*[]*string {
	var returns *map[string]*[]*string
	_jsii_.Get(
		j,
		"principalJson",
		&returns,
	)
	return returns
}


// Experimental.
func NewPrincipalPolicyFragment(principalJson *map[string]*[]*string, conditions *map[string]interface{}) PrincipalPolicyFragment {
	_init_.Initialize()

	j := jsiiProxy_PrincipalPolicyFragment{}

	_jsii_.Create(
		"monocdk.aws_iam.PrincipalPolicyFragment",
		[]interface{}{principalJson, conditions},
		&j,
	)

	return &j
}

// Experimental.
func NewPrincipalPolicyFragment_Override(p PrincipalPolicyFragment, principalJson *map[string]*[]*string, conditions *map[string]interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.PrincipalPolicyFragment",
		[]interface{}{principalJson, conditions},
		p,
	)
}

// An IAM principal with additional conditions specifying when the policy is in effect.
//
// For more information about conditions, see:
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition.html
// Experimental.
type PrincipalWithConditions interface {
	IPrincipal
	AssumeRoleAction() *string
	Conditions() *map[string]interface{}
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	AddCondition(key *string, value interface{})
	AddConditions(conditions *map[string]interface{})
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
}

// The jsii proxy struct for PrincipalWithConditions
type jsiiProxy_PrincipalWithConditions struct {
	jsiiProxy_IPrincipal
}

func (j *jsiiProxy_PrincipalWithConditions) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PrincipalWithConditions) Conditions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PrincipalWithConditions) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_PrincipalWithConditions) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}


// Experimental.
func NewPrincipalWithConditions(principal IPrincipal, conditions *map[string]interface{}) PrincipalWithConditions {
	_init_.Initialize()

	j := jsiiProxy_PrincipalWithConditions{}

	_jsii_.Create(
		"monocdk.aws_iam.PrincipalWithConditions",
		[]interface{}{principal, conditions},
		&j,
	)

	return &j
}

// Experimental.
func NewPrincipalWithConditions_Override(p PrincipalWithConditions, principal IPrincipal, conditions *map[string]interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.PrincipalWithConditions",
		[]interface{}{principal, conditions},
		p,
	)
}

// Add a condition to the principal.
// Experimental.
func (p *jsiiProxy_PrincipalWithConditions) AddCondition(key *string, value interface{}) {
	_jsii_.InvokeVoid(
		p,
		"addCondition",
		[]interface{}{key, value},
	)
}

// Adds multiple conditions to the principal.
//
// Values from the conditions parameter will overwrite existing values with the same operator
// and key.
// Experimental.
func (p *jsiiProxy_PrincipalWithConditions) AddConditions(conditions *map[string]interface{}) {
	_jsii_.InvokeVoid(
		p,
		"addConditions",
		[]interface{}{conditions},
	)
}

// Add to the policy of this principal.
// Experimental.
func (p *jsiiProxy_PrincipalWithConditions) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		p,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (p *jsiiProxy_PrincipalWithConditions) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		p,
		"addToPrincipalPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (p *jsiiProxy_PrincipalWithConditions) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		p,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (p *jsiiProxy_PrincipalWithConditions) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// IAM Role.
//
// Defines an IAM role. The role is created with an assume policy document associated with
// the specified AWS service principal defined in `serviceAssumeRole`.
// Experimental.
type Role interface {
	awscdk.Resource
	IRole
	AssumeRoleAction() *string
	AssumeRolePolicy() PolicyDocument
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() IPrincipal
	Node() awscdk.ConstructNode
	PermissionsBoundary() IManagedPolicy
	PhysicalName() *string
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	RoleArn() *string
	RoleId() *string
	RoleName() *string
	Stack() awscdk.Stack
	AddManagedPolicy(policy IManagedPolicy)
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AttachInlinePolicy(policy Policy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Grant(grantee IPrincipal, actions ...*string) Grant
	GrantPassRole(identity IPrincipal) Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
	WithoutPolicyUpdates() IRole
}

// The jsii proxy struct for Role
type jsiiProxy_Role struct {
	internal.Type__awscdkResource
	jsiiProxy_IRole
}

func (j *jsiiProxy_Role) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) AssumeRolePolicy() PolicyDocument {
	var returns PolicyDocument
	_jsii_.Get(
		j,
		"assumeRolePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) PermissionsBoundary() IManagedPolicy {
	var returns IManagedPolicy
	_jsii_.Get(
		j,
		"permissionsBoundary",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) RoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) RoleId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) RoleName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"roleName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Role) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRole(scope constructs.Construct, id *string, props *RoleProps) Role {
	_init_.Initialize()

	j := jsiiProxy_Role{}

	_jsii_.Create(
		"monocdk.aws_iam.Role",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRole_Override(r Role, scope constructs.Construct, id *string, props *RoleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.Role",
		[]interface{}{scope, id, props},
		r,
	)
}

// Import an external role by ARN.
//
// If the imported Role ARN is a Token (such as a
// `CfnParameter.valueAsString` or a `Fn.importValue()`) *and* the referenced
// role has a `path` (like `arn:...:role/AdminRoles/Alice`), the
// `roleName` property will not resolve to the correct value. Instead it
// will resolve to the first path component. We unfortunately cannot express
// the correct calculation of the full path name as a CloudFormation
// expression. In this scenario the Role ARN should be supplied without the
// `path` in order to resolve the correct role resource.
// Experimental.
func Role_FromRoleArn(scope constructs.Construct, id *string, roleArn *string, options *FromRoleArnOptions) IRole {
	_init_.Initialize()

	var returns IRole

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Role",
		"fromRoleArn",
		[]interface{}{scope, id, roleArn, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Role_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Role",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Role_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.Role",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Attaches a managed policy to this role.
// Experimental.
func (r *jsiiProxy_Role) AddManagedPolicy(policy IManagedPolicy) {
	_jsii_.InvokeVoid(
		r,
		"addManagedPolicy",
		[]interface{}{policy},
	)
}

// Add to the policy of this principal.
// Experimental.
func (r *jsiiProxy_Role) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		r,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Adds a permission to the role's default policy document.
//
// If there is no default policy attached to this role, it will be created.
// Experimental.
func (r *jsiiProxy_Role) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		r,
		"addToPrincipalPolicy",
		[]interface{}{statement},
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
func (r *jsiiProxy_Role) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attaches a policy to this role.
// Experimental.
func (r *jsiiProxy_Role) AttachInlinePolicy(policy Policy) {
	_jsii_.InvokeVoid(
		r,
		"attachInlinePolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_Role) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_Role) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_Role) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grant the actions defined in actions to the identity Principal on this resource.
// Experimental.
func (r *jsiiProxy_Role) Grant(grantee IPrincipal, actions ...*string) Grant {
	args := []interface{}{grantee}
	for _, a := range actions {
		args = append(args, a)
	}

	var returns Grant

	_jsii_.Invoke(
		r,
		"grant",
		args,
		&returns,
	)

	return returns
}

// Grant permissions to the given principal to pass this role.
// Experimental.
func (r *jsiiProxy_Role) GrantPassRole(identity IPrincipal) Grant {
	var returns Grant

	_jsii_.Invoke(
		r,
		"grantPassRole",
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
func (r *jsiiProxy_Role) OnPrepare() {
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
func (r *jsiiProxy_Role) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_Role) OnValidate() *[]*string {
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
func (r *jsiiProxy_Role) Prepare() {
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
func (r *jsiiProxy_Role) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_Role) ToString() *string {
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
// Experimental.
func (r *jsiiProxy_Role) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Return a copy of this Role object whose Policies will not be updated.
//
// Use the object returned by this method if you want this Role to be used by
// a construct without it automatically updating the Role's Policies.
//
// If you do, you are responsible for adding the correct statements to the
// Role's policies yourself.
// Experimental.
func (r *jsiiProxy_Role) WithoutPolicyUpdates() IRole {
	var returns IRole

	_jsii_.Invoke(
		r,
		"withoutPolicyUpdates",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining an IAM Role.
// Experimental.
type RoleProps struct {
	// The IAM principal (i.e. `new ServicePrincipal('sns.amazonaws.com')`) which can assume this role.
	//
	// You can later modify the assume role policy document by accessing it via
	// the `assumeRolePolicy` property.
	// Experimental.
	AssumedBy IPrincipal `json:"assumedBy"`
	// A description of the role.
	//
	// It can be up to 1000 characters long.
	// Experimental.
	Description *string `json:"description"`
	// ID that the role assumer needs to provide when assuming this role.
	//
	// If the configured and provided external IDs do not match, the
	// AssumeRole operation will fail.
	// Deprecated: see {@link externalIds}
	ExternalId *string `json:"externalId"`
	// List of IDs that the role assumer needs to provide one of when assuming this role.
	//
	// If the configured and provided external IDs do not match, the
	// AssumeRole operation will fail.
	// Experimental.
	ExternalIds *[]*string `json:"externalIds"`
	// A list of named policies to inline into this role.
	//
	// These policies will be
	// created with the role, whereas those added by ``addToPolicy`` are added
	// using a separate CloudFormation resource (allowing a way around circular
	// dependencies that could otherwise be introduced).
	// Experimental.
	InlinePolicies *map[string]PolicyDocument `json:"inlinePolicies"`
	// A list of managed policies associated with this role.
	//
	// You can add managed policies later using
	// `addManagedPolicy(ManagedPolicy.fromAwsManagedPolicyName(policyName))`.
	// Experimental.
	ManagedPolicies *[]IManagedPolicy `json:"managedPolicies"`
	// The maximum session duration that you want to set for the specified role.
	//
	// This setting can have a value from 1 hour (3600sec) to 12 (43200sec) hours.
	//
	// Anyone who assumes the role from the AWS CLI or API can use the
	// DurationSeconds API parameter or the duration-seconds CLI parameter to
	// request a longer session. The MaxSessionDuration setting determines the
	// maximum duration that can be requested using the DurationSeconds
	// parameter.
	//
	// If users don't specify a value for the DurationSeconds parameter, their
	// security credentials are valid for one hour by default. This applies when
	// you use the AssumeRole* API operations or the assume-role* CLI operations
	// but does not apply when you use those operations to create a console URL.
	// Experimental.
	MaxSessionDuration awscdk.Duration `json:"maxSessionDuration"`
	// The path associated with this role.
	//
	// For information about IAM paths, see
	// Friendly Names and Paths in IAM User Guide.
	// Experimental.
	Path *string `json:"path"`
	// AWS supports permissions boundaries for IAM entities (users or roles).
	//
	// A permissions boundary is an advanced feature for using a managed policy
	// to set the maximum permissions that an identity-based policy can grant to
	// an IAM entity. An entity's permissions boundary allows it to perform only
	// the actions that are allowed by both its identity-based policies and its
	// permissions boundaries.
	// Experimental.
	PermissionsBoundary IManagedPolicy `json:"permissionsBoundary"`
	// A name for the IAM role.
	//
	// For valid values, see the RoleName parameter for
	// the CreateRole action in the IAM API Reference.
	//
	// IMPORTANT: If you specify a name, you cannot perform updates that require
	// replacement of this resource. You can perform updates that require no or
	// some interruption. If you must replace the resource, specify a new name.
	//
	// If you specify a name, you must specify the CAPABILITY_NAMED_IAM value to
	// acknowledge your template's capabilities. For more information, see
	// Acknowledging IAM Resources in AWS CloudFormation Templates.
	// Experimental.
	RoleName *string `json:"roleName"`
}

// Principal entity that represents a SAML federated identity provider for programmatic and AWS Management Console access.
// Experimental.
type SamlConsolePrincipal interface {
	SamlPrincipal
	AssumeRoleAction() *string
	Conditions() *map[string]interface{}
	Federated() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for SamlConsolePrincipal
type jsiiProxy_SamlConsolePrincipal struct {
	jsiiProxy_SamlPrincipal
}

func (j *jsiiProxy_SamlConsolePrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlConsolePrincipal) Conditions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlConsolePrincipal) Federated() *string {
	var returns *string
	_jsii_.Get(
		j,
		"federated",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlConsolePrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlConsolePrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlConsolePrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewSamlConsolePrincipal(samlProvider ISamlProvider, conditions *map[string]interface{}) SamlConsolePrincipal {
	_init_.Initialize()

	j := jsiiProxy_SamlConsolePrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.SamlConsolePrincipal",
		[]interface{}{samlProvider, conditions},
		&j,
	)

	return &j
}

// Experimental.
func NewSamlConsolePrincipal_Override(s SamlConsolePrincipal, samlProvider ISamlProvider, conditions *map[string]interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.SamlConsolePrincipal",
		[]interface{}{samlProvider, conditions},
		s,
	)
}

// Add to the policy of this principal.
// Experimental.
func (s *jsiiProxy_SamlConsolePrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		s,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (s *jsiiProxy_SamlConsolePrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		s,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (s *jsiiProxy_SamlConsolePrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		s,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (s *jsiiProxy_SamlConsolePrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (s *jsiiProxy_SamlConsolePrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		s,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// A SAML metadata document.
// Experimental.
type SamlMetadataDocument interface {
	Xml() *string
}

// The jsii proxy struct for SamlMetadataDocument
type jsiiProxy_SamlMetadataDocument struct {
	_ byte // padding
}

func (j *jsiiProxy_SamlMetadataDocument) Xml() *string {
	var returns *string
	_jsii_.Get(
		j,
		"xml",
		&returns,
	)
	return returns
}


// Experimental.
func NewSamlMetadataDocument_Override(s SamlMetadataDocument) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.SamlMetadataDocument",
		nil, // no parameters
		s,
	)
}

// Create a SAML metadata document from a XML file.
// Experimental.
func SamlMetadataDocument_FromFile(path *string) SamlMetadataDocument {
	_init_.Initialize()

	var returns SamlMetadataDocument

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.SamlMetadataDocument",
		"fromFile",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Create a SAML metadata document from a XML string.
// Experimental.
func SamlMetadataDocument_FromXml(xml *string) SamlMetadataDocument {
	_init_.Initialize()

	var returns SamlMetadataDocument

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.SamlMetadataDocument",
		"fromXml",
		[]interface{}{xml},
		&returns,
	)

	return returns
}

// Principal entity that represents a SAML federated identity provider.
// Experimental.
type SamlPrincipal interface {
	FederatedPrincipal
	AssumeRoleAction() *string
	Conditions() *map[string]interface{}
	Federated() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for SamlPrincipal
type jsiiProxy_SamlPrincipal struct {
	jsiiProxy_FederatedPrincipal
}

func (j *jsiiProxy_SamlPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlPrincipal) Conditions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlPrincipal) Federated() *string {
	var returns *string
	_jsii_.Get(
		j,
		"federated",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewSamlPrincipal(samlProvider ISamlProvider, conditions *map[string]interface{}) SamlPrincipal {
	_init_.Initialize()

	j := jsiiProxy_SamlPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.SamlPrincipal",
		[]interface{}{samlProvider, conditions},
		&j,
	)

	return &j
}

// Experimental.
func NewSamlPrincipal_Override(s SamlPrincipal, samlProvider ISamlProvider, conditions *map[string]interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.SamlPrincipal",
		[]interface{}{samlProvider, conditions},
		s,
	)
}

// Add to the policy of this principal.
// Experimental.
func (s *jsiiProxy_SamlPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		s,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (s *jsiiProxy_SamlPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		s,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (s *jsiiProxy_SamlPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		s,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (s *jsiiProxy_SamlPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (s *jsiiProxy_SamlPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		s,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// A SAML provider.
// Experimental.
type SamlProvider interface {
	awscdk.Resource
	ISamlProvider
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	SamlProviderArn() *string
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

// The jsii proxy struct for SamlProvider
type jsiiProxy_SamlProvider struct {
	internal.Type__awscdkResource
	jsiiProxy_ISamlProvider
}

func (j *jsiiProxy_SamlProvider) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlProvider) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlProvider) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlProvider) SamlProviderArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"samlProviderArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SamlProvider) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewSamlProvider(scope constructs.Construct, id *string, props *SamlProviderProps) SamlProvider {
	_init_.Initialize()

	j := jsiiProxy_SamlProvider{}

	_jsii_.Create(
		"monocdk.aws_iam.SamlProvider",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSamlProvider_Override(s SamlProvider, scope constructs.Construct, id *string, props *SamlProviderProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.SamlProvider",
		[]interface{}{scope, id, props},
		s,
	)
}

// Import an existing provider.
// Experimental.
func SamlProvider_FromSamlProviderArn(scope constructs.Construct, id *string, samlProviderArn *string) ISamlProvider {
	_init_.Initialize()

	var returns ISamlProvider

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.SamlProvider",
		"fromSamlProviderArn",
		[]interface{}{scope, id, samlProviderArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func SamlProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.SamlProvider",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func SamlProvider_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.SamlProvider",
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
func (s *jsiiProxy_SamlProvider) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (s *jsiiProxy_SamlProvider) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_SamlProvider) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_SamlProvider) GetResourceNameAttribute(nameAttr *string) *string {
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
func (s *jsiiProxy_SamlProvider) OnPrepare() {
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
func (s *jsiiProxy_SamlProvider) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_SamlProvider) OnValidate() *[]*string {
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
func (s *jsiiProxy_SamlProvider) Prepare() {
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
func (s *jsiiProxy_SamlProvider) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_SamlProvider) ToString() *string {
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
func (s *jsiiProxy_SamlProvider) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a SAML provider.
// Experimental.
type SamlProviderProps struct {
	// An XML document generated by an identity provider (IdP) that supports SAML 2.0. The document includes the issuer's name, expiration information, and keys that can be used to validate the SAML authentication response (assertions) that are received from the IdP. You must generate the metadata document using the identity management software that is used as your organization's IdP.
	// Experimental.
	MetadataDocument SamlMetadataDocument `json:"metadataDocument"`
	// The name of the provider to create.
	//
	// This parameter allows a string of characters consisting of upper and
	// lowercase alphanumeric characters with no spaces. You can also include
	// any of the following characters: _+=,.@-
	//
	// Length must be between 1 and 128 characters.
	// Experimental.
	Name *string `json:"name"`
}

// An IAM principal that represents an AWS service (i.e. sqs.amazonaws.com).
// Experimental.
type ServicePrincipal interface {
	PrincipalBase
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	Service() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for ServicePrincipal
type jsiiProxy_ServicePrincipal struct {
	jsiiProxy_PrincipalBase
}

func (j *jsiiProxy_ServicePrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServicePrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServicePrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServicePrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ServicePrincipal) Service() *string {
	var returns *string
	_jsii_.Get(
		j,
		"service",
		&returns,
	)
	return returns
}


// Experimental.
func NewServicePrincipal(service *string, opts *ServicePrincipalOpts) ServicePrincipal {
	_init_.Initialize()

	j := jsiiProxy_ServicePrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.ServicePrincipal",
		[]interface{}{service, opts},
		&j,
	)

	return &j
}

// Experimental.
func NewServicePrincipal_Override(s ServicePrincipal, service *string, opts *ServicePrincipalOpts) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.ServicePrincipal",
		[]interface{}{service, opts},
		s,
	)
}

// Add to the policy of this principal.
// Experimental.
func (s *jsiiProxy_ServicePrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		s,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (s *jsiiProxy_ServicePrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		s,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (s *jsiiProxy_ServicePrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		s,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (s *jsiiProxy_ServicePrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (s *jsiiProxy_ServicePrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		s,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}

// Options for a service principal.
// Experimental.
type ServicePrincipalOpts struct {
	// Additional conditions to add to the Service Principal.
	// Experimental.
	Conditions *map[string]interface{} `json:"conditions"`
	// The region in which the service is operating.
	// Experimental.
	Region *string `json:"region"`
}

// A principal for use in resources that need to have a role but it's unknown.
//
// Some resources have roles associated with them which they assume, such as
// Lambda Functions, CodeBuild projects, StepFunctions machines, etc.
//
// When those resources are imported, their actual roles are not always
// imported with them. When that happens, we use an instance of this class
// instead, which will add user warnings when statements are attempted to be
// added to it.
// Experimental.
type UnknownPrincipal interface {
	IPrincipal
	AssumeRoleAction() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult
}

// The jsii proxy struct for UnknownPrincipal
type jsiiProxy_UnknownPrincipal struct {
	jsiiProxy_IPrincipal
}

func (j *jsiiProxy_UnknownPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UnknownPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UnknownPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}


// Experimental.
func NewUnknownPrincipal(props *UnknownPrincipalProps) UnknownPrincipal {
	_init_.Initialize()

	j := jsiiProxy_UnknownPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.UnknownPrincipal",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewUnknownPrincipal_Override(u UnknownPrincipal, props *UnknownPrincipalProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.UnknownPrincipal",
		[]interface{}{props},
		u,
	)
}

// Add to the policy of this principal.
// Experimental.
func (u *jsiiProxy_UnknownPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		u,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (u *jsiiProxy_UnknownPrincipal) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		u,
		"addToPrincipalPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Properties for an UnknownPrincipal.
// Experimental.
type UnknownPrincipalProps struct {
	// The resource the role proxy is for.
	// Experimental.
	Resource constructs.IConstruct `json:"resource"`
}

// Define a new IAM user.
// Experimental.
type User interface {
	awscdk.Resource
	IIdentity
	IUser
	AssumeRoleAction() *string
	Env() *awscdk.ResourceEnvironment
	GrantPrincipal() IPrincipal
	Node() awscdk.ConstructNode
	PermissionsBoundary() IManagedPolicy
	PhysicalName() *string
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	Stack() awscdk.Stack
	UserArn() *string
	UserName() *string
	AddManagedPolicy(policy IManagedPolicy)
	AddToGroup(group IGroup)
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	AttachInlinePolicy(policy Policy)
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

// The jsii proxy struct for User
type jsiiProxy_User struct {
	internal.Type__awscdkResource
	jsiiProxy_IIdentity
	jsiiProxy_IUser
}

func (j *jsiiProxy_User) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) PermissionsBoundary() IManagedPolicy {
	var returns IManagedPolicy
	_jsii_.Get(
		j,
		"permissionsBoundary",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) UserArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_User) UserName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"userName",
		&returns,
	)
	return returns
}


// Experimental.
func NewUser(scope constructs.Construct, id *string, props *UserProps) User {
	_init_.Initialize()

	j := jsiiProxy_User{}

	_jsii_.Create(
		"monocdk.aws_iam.User",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUser_Override(u User, scope constructs.Construct, id *string, props *UserProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.User",
		[]interface{}{scope, id, props},
		u,
	)
}

// Import an existing user given a user ARN.
// Experimental.
func User_FromUserArn(scope constructs.Construct, id *string, userArn *string) IUser {
	_init_.Initialize()

	var returns IUser

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.User",
		"fromUserArn",
		[]interface{}{scope, id, userArn},
		&returns,
	)

	return returns
}

// Import an existing user given user attributes.
// Experimental.
func User_FromUserAttributes(scope constructs.Construct, id *string, attrs *UserAttributes) IUser {
	_init_.Initialize()

	var returns IUser

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.User",
		"fromUserAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Import an existing user given a username.
// Experimental.
func User_FromUserName(scope constructs.Construct, id *string, userName *string) IUser {
	_init_.Initialize()

	var returns IUser

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.User",
		"fromUserName",
		[]interface{}{scope, id, userName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func User_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.User",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func User_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_iam.User",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Attaches a managed policy to the user.
// Experimental.
func (u *jsiiProxy_User) AddManagedPolicy(policy IManagedPolicy) {
	_jsii_.InvokeVoid(
		u,
		"addManagedPolicy",
		[]interface{}{policy},
	)
}

// Adds this user to a group.
// Experimental.
func (u *jsiiProxy_User) AddToGroup(group IGroup) {
	_jsii_.InvokeVoid(
		u,
		"addToGroup",
		[]interface{}{group},
	)
}

// Add to the policy of this principal.
// Experimental.
func (u *jsiiProxy_User) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		u,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Adds an IAM statement to the default policy.
//
// Returns: true
// Experimental.
func (u *jsiiProxy_User) AddToPrincipalPolicy(statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		u,
		"addToPrincipalPolicy",
		[]interface{}{statement},
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
func (u *jsiiProxy_User) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Attaches a policy to this user.
// Experimental.
func (u *jsiiProxy_User) AttachInlinePolicy(policy Policy) {
	_jsii_.InvokeVoid(
		u,
		"attachInlinePolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_User) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_User) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_User) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_User) OnPrepare() {
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
func (u *jsiiProxy_User) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_User) OnValidate() *[]*string {
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
func (u *jsiiProxy_User) Prepare() {
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
func (u *jsiiProxy_User) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_User) ToString() *string {
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
func (u *jsiiProxy_User) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents a user defined outside of this stack.
// Experimental.
type UserAttributes struct {
	// The ARN of the user.
	//
	// Format: arn:<partition>:iam::<account-id>:user/<user-name-with-path>
	// Experimental.
	UserArn *string `json:"userArn"`
}

// Properties for defining an IAM user.
// Experimental.
type UserProps struct {
	// Groups to add this user to.
	//
	// You can also use `addToGroup` to add this
	// user to a group.
	// Experimental.
	Groups *[]IGroup `json:"groups"`
	// A list of managed policies associated with this role.
	//
	// You can add managed policies later using
	// `addManagedPolicy(ManagedPolicy.fromAwsManagedPolicyName(policyName))`.
	// Experimental.
	ManagedPolicies *[]IManagedPolicy `json:"managedPolicies"`
	// The password for the user. This is required so the user can access the AWS Management Console.
	//
	// You can use `SecretValue.plainText` to specify a password in plain text or
	// use `secretsmanager.Secret.fromSecretAttributes` to reference a secret in
	// Secrets Manager.
	// Experimental.
	Password awscdk.SecretValue `json:"password"`
	// Specifies whether the user is required to set a new password the next time the user logs in to the AWS Management Console.
	//
	// If this is set to 'true', you must also specify "initialPassword".
	// Experimental.
	PasswordResetRequired *bool `json:"passwordResetRequired"`
	// The path for the user name.
	//
	// For more information about paths, see IAM
	// Identifiers in the IAM User Guide.
	// Experimental.
	Path *string `json:"path"`
	// AWS supports permissions boundaries for IAM entities (users or roles).
	//
	// A permissions boundary is an advanced feature for using a managed policy
	// to set the maximum permissions that an identity-based policy can grant to
	// an IAM entity. An entity's permissions boundary allows it to perform only
	// the actions that are allowed by both its identity-based policies and its
	// permissions boundaries.
	// Experimental.
	PermissionsBoundary IManagedPolicy `json:"permissionsBoundary"`
	// A name for the IAM user.
	//
	// For valid values, see the UserName parameter for
	// the CreateUser action in the IAM API Reference. If you don't specify a
	// name, AWS CloudFormation generates a unique physical ID and uses that ID
	// for the user name.
	//
	// If you specify a name, you cannot perform updates that require
	// replacement of this resource. You can perform updates that require no or
	// some interruption. If you must replace the resource, specify a new name.
	//
	// If you specify a name, you must specify the CAPABILITY_NAMED_IAM value to
	// acknowledge your template's capabilities. For more information, see
	// Acknowledging IAM Resources in AWS CloudFormation Templates.
	// Experimental.
	UserName *string `json:"userName"`
}

// A principal that represents a federated identity provider as Web Identity such as Cognito, Amazon, Facebook, Google, etc.
// Experimental.
type WebIdentityPrincipal interface {
	FederatedPrincipal
	AssumeRoleAction() *string
	Conditions() *map[string]interface{}
	Federated() *string
	GrantPrincipal() IPrincipal
	PolicyFragment() PrincipalPolicyFragment
	PrincipalAccount() *string
	AddToPolicy(statement PolicyStatement) *bool
	AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult
	ToJSON() *map[string]*[]*string
	ToString() *string
	WithConditions(conditions *map[string]interface{}) IPrincipal
}

// The jsii proxy struct for WebIdentityPrincipal
type jsiiProxy_WebIdentityPrincipal struct {
	jsiiProxy_FederatedPrincipal
}

func (j *jsiiProxy_WebIdentityPrincipal) AssumeRoleAction() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assumeRoleAction",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WebIdentityPrincipal) Conditions() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WebIdentityPrincipal) Federated() *string {
	var returns *string
	_jsii_.Get(
		j,
		"federated",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WebIdentityPrincipal) GrantPrincipal() IPrincipal {
	var returns IPrincipal
	_jsii_.Get(
		j,
		"grantPrincipal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WebIdentityPrincipal) PolicyFragment() PrincipalPolicyFragment {
	var returns PrincipalPolicyFragment
	_jsii_.Get(
		j,
		"policyFragment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_WebIdentityPrincipal) PrincipalAccount() *string {
	var returns *string
	_jsii_.Get(
		j,
		"principalAccount",
		&returns,
	)
	return returns
}


// Experimental.
func NewWebIdentityPrincipal(identityProvider *string, conditions *map[string]interface{}) WebIdentityPrincipal {
	_init_.Initialize()

	j := jsiiProxy_WebIdentityPrincipal{}

	_jsii_.Create(
		"monocdk.aws_iam.WebIdentityPrincipal",
		[]interface{}{identityProvider, conditions},
		&j,
	)

	return &j
}

// Experimental.
func NewWebIdentityPrincipal_Override(w WebIdentityPrincipal, identityProvider *string, conditions *map[string]interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_iam.WebIdentityPrincipal",
		[]interface{}{identityProvider, conditions},
		w,
	)
}

// Add to the policy of this principal.
// Experimental.
func (w *jsiiProxy_WebIdentityPrincipal) AddToPolicy(statement PolicyStatement) *bool {
	var returns *bool

	_jsii_.Invoke(
		w,
		"addToPolicy",
		[]interface{}{statement},
		&returns,
	)

	return returns
}

// Add to the policy of this principal.
// Experimental.
func (w *jsiiProxy_WebIdentityPrincipal) AddToPrincipalPolicy(_statement PolicyStatement) *AddToPrincipalPolicyResult {
	var returns *AddToPrincipalPolicyResult

	_jsii_.Invoke(
		w,
		"addToPrincipalPolicy",
		[]interface{}{_statement},
		&returns,
	)

	return returns
}

// JSON-ify the principal.
//
// Used when JSON.stringify() is called
// Experimental.
func (w *jsiiProxy_WebIdentityPrincipal) ToJSON() *map[string]*[]*string {
	var returns *map[string]*[]*string

	_jsii_.Invoke(
		w,
		"toJSON",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a string representation of an object.
// Experimental.
func (w *jsiiProxy_WebIdentityPrincipal) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		w,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns a new PrincipalWithConditions using this principal as the base, with the passed conditions added.
//
// When there is a value for the same operator and key in both the principal and the
// conditions parameter, the value from the conditions parameter will be used.
//
// Returns: a new PrincipalWithConditions object.
// Experimental.
func (w *jsiiProxy_WebIdentityPrincipal) WithConditions(conditions *map[string]interface{}) IPrincipal {
	var returns IPrincipal

	_jsii_.Invoke(
		w,
		"withConditions",
		[]interface{}{conditions},
		&returns,
	)

	return returns
}


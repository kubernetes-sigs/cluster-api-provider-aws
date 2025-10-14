package awss3

import (
	"time"

	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awss3/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// Experimental.
type BlockPublicAccess interface {
	BlockPublicAcls() *bool
	SetBlockPublicAcls(val *bool)
	BlockPublicPolicy() *bool
	SetBlockPublicPolicy(val *bool)
	IgnorePublicAcls() *bool
	SetIgnorePublicAcls(val *bool)
	RestrictPublicBuckets() *bool
	SetRestrictPublicBuckets(val *bool)
}

// The jsii proxy struct for BlockPublicAccess
type jsiiProxy_BlockPublicAccess struct {
	_ byte // padding
}

func (j *jsiiProxy_BlockPublicAccess) BlockPublicAcls() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"blockPublicAcls",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BlockPublicAccess) BlockPublicPolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"blockPublicPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BlockPublicAccess) IgnorePublicAcls() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"ignorePublicAcls",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BlockPublicAccess) RestrictPublicBuckets() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"restrictPublicBuckets",
		&returns,
	)
	return returns
}


// Experimental.
func NewBlockPublicAccess(options *BlockPublicAccessOptions) BlockPublicAccess {
	_init_.Initialize()

	j := jsiiProxy_BlockPublicAccess{}

	_jsii_.Create(
		"monocdk.aws_s3.BlockPublicAccess",
		[]interface{}{options},
		&j,
	)

	return &j
}

// Experimental.
func NewBlockPublicAccess_Override(b BlockPublicAccess, options *BlockPublicAccessOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.BlockPublicAccess",
		[]interface{}{options},
		b,
	)
}

func (j *jsiiProxy_BlockPublicAccess) SetBlockPublicAcls(val *bool) {
	_jsii_.Set(
		j,
		"blockPublicAcls",
		val,
	)
}

func (j *jsiiProxy_BlockPublicAccess) SetBlockPublicPolicy(val *bool) {
	_jsii_.Set(
		j,
		"blockPublicPolicy",
		val,
	)
}

func (j *jsiiProxy_BlockPublicAccess) SetIgnorePublicAcls(val *bool) {
	_jsii_.Set(
		j,
		"ignorePublicAcls",
		val,
	)
}

func (j *jsiiProxy_BlockPublicAccess) SetRestrictPublicBuckets(val *bool) {
	_jsii_.Set(
		j,
		"restrictPublicBuckets",
		val,
	)
}

func BlockPublicAccess_BLOCK_ACLS() BlockPublicAccess {
	_init_.Initialize()
	var returns BlockPublicAccess
	_jsii_.StaticGet(
		"monocdk.aws_s3.BlockPublicAccess",
		"BLOCK_ACLS",
		&returns,
	)
	return returns
}

func BlockPublicAccess_BLOCK_ALL() BlockPublicAccess {
	_init_.Initialize()
	var returns BlockPublicAccess
	_jsii_.StaticGet(
		"monocdk.aws_s3.BlockPublicAccess",
		"BLOCK_ALL",
		&returns,
	)
	return returns
}

// Experimental.
type BlockPublicAccessOptions struct {
	// Whether to block public ACLs.
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/access-control-block-public-access.html#access-control-block-public-access-options
	//
	// Experimental.
	BlockPublicAcls *bool `json:"blockPublicAcls"`
	// Whether to block public policy.
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/access-control-block-public-access.html#access-control-block-public-access-options
	//
	// Experimental.
	BlockPublicPolicy *bool `json:"blockPublicPolicy"`
	// Whether to ignore public ACLs.
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/access-control-block-public-access.html#access-control-block-public-access-options
	//
	// Experimental.
	IgnorePublicAcls *bool `json:"ignorePublicAcls"`
	// Whether to restrict public access.
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/access-control-block-public-access.html#access-control-block-public-access-options
	//
	// Experimental.
	RestrictPublicBuckets *bool `json:"restrictPublicBuckets"`
}

// An S3 bucket with associated policy objects.
//
// This bucket does not yet have all features that exposed by the underlying
// BucketResource.
// Experimental.
type Bucket interface {
	BucketBase
	AutoCreatePolicy() *bool
	SetAutoCreatePolicy(val *bool)
	BucketArn() *string
	BucketDomainName() *string
	BucketDualStackDomainName() *string
	BucketName() *string
	BucketRegionalDomainName() *string
	BucketWebsiteDomainName() *string
	BucketWebsiteUrl() *string
	DisallowPublicAccess() *bool
	SetDisallowPublicAccess(val *bool)
	EncryptionKey() awskms.IKey
	Env() *awscdk.ResourceEnvironment
	IsWebsite() *bool
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Policy() BucketPolicy
	SetPolicy(val BucketPolicy)
	Stack() awscdk.Stack
	AddCorsRule(rule *CorsRule)
	AddEventNotification(event EventType, dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	AddInventory(inventory *Inventory)
	AddLifecycleRule(rule *LifecycleRule)
	AddMetric(metric *BucketMetrics)
	AddObjectCreatedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	AddObjectRemovedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	AddToResourcePolicy(permission awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ArnForObjects(keyPattern *string) *string
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantDelete(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantPublicAccess(keyPrefix *string, allowedActions ...*string) awsiam.Grant
	GrantPut(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantPutAcl(identity awsiam.IGrantable, objectsKeyPattern *string) awsiam.Grant
	GrantRead(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantReadWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	OnCloudTrailEvent(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	OnCloudTrailPutObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	OnCloudTrailWriteObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	S3UrlForObject(key *string) *string
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	UrlForObject(key *string) *string
	Validate() *[]*string
	VirtualHostedUrlForObject(key *string, options *VirtualHostedStyleUrlOptions) *string
}

// The jsii proxy struct for Bucket
type jsiiProxy_Bucket struct {
	jsiiProxy_BucketBase
}

func (j *jsiiProxy_Bucket) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) BucketArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) BucketDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) BucketDualStackDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketDualStackDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) BucketName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) BucketRegionalDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketRegionalDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) BucketWebsiteDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketWebsiteDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) BucketWebsiteUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketWebsiteUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) DisallowPublicAccess() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"disallowPublicAccess",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) IsWebsite() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isWebsite",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) Policy() BucketPolicy {
	var returns BucketPolicy
	_jsii_.Get(
		j,
		"policy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Bucket) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewBucket(scope constructs.Construct, id *string, props *BucketProps) Bucket {
	_init_.Initialize()

	j := jsiiProxy_Bucket{}

	_jsii_.Create(
		"monocdk.aws_s3.Bucket",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewBucket_Override(b Bucket, scope constructs.Construct, id *string, props *BucketProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.Bucket",
		[]interface{}{scope, id, props},
		b,
	)
}

func (j *jsiiProxy_Bucket) SetAutoCreatePolicy(val *bool) {
	_jsii_.Set(
		j,
		"autoCreatePolicy",
		val,
	)
}

func (j *jsiiProxy_Bucket) SetDisallowPublicAccess(val *bool) {
	_jsii_.Set(
		j,
		"disallowPublicAccess",
		val,
	)
}

func (j *jsiiProxy_Bucket) SetPolicy(val BucketPolicy) {
	_jsii_.Set(
		j,
		"policy",
		val,
	)
}

// Experimental.
func Bucket_FromBucketArn(scope constructs.Construct, id *string, bucketArn *string) IBucket {
	_init_.Initialize()

	var returns IBucket

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.Bucket",
		"fromBucketArn",
		[]interface{}{scope, id, bucketArn},
		&returns,
	)

	return returns
}

// Creates a Bucket construct that represents an external bucket.
// Experimental.
func Bucket_FromBucketAttributes(scope constructs.Construct, id *string, attrs *BucketAttributes) IBucket {
	_init_.Initialize()

	var returns IBucket

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.Bucket",
		"fromBucketAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Experimental.
func Bucket_FromBucketName(scope constructs.Construct, id *string, bucketName *string) IBucket {
	_init_.Initialize()

	var returns IBucket

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.Bucket",
		"fromBucketName",
		[]interface{}{scope, id, bucketName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Bucket_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.Bucket",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Bucket_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.Bucket",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a cross-origin access configuration for objects in an Amazon S3 bucket.
// Experimental.
func (b *jsiiProxy_Bucket) AddCorsRule(rule *CorsRule) {
	_jsii_.InvokeVoid(
		b,
		"addCorsRule",
		[]interface{}{rule},
	)
}

// Adds a bucket notification event destination.
//
// TODO: EXAMPLE
//
// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html
//
// Experimental.
func (b *jsiiProxy_Bucket) AddEventNotification(event EventType, dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{event, dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		b,
		"addEventNotification",
		args,
	)
}

// Add an inventory configuration.
// Experimental.
func (b *jsiiProxy_Bucket) AddInventory(inventory *Inventory) {
	_jsii_.InvokeVoid(
		b,
		"addInventory",
		[]interface{}{inventory},
	)
}

// Add a lifecycle rule to the bucket.
// Experimental.
func (b *jsiiProxy_Bucket) AddLifecycleRule(rule *LifecycleRule) {
	_jsii_.InvokeVoid(
		b,
		"addLifecycleRule",
		[]interface{}{rule},
	)
}

// Adds a metrics configuration for the CloudWatch request metrics from the bucket.
// Experimental.
func (b *jsiiProxy_Bucket) AddMetric(metric *BucketMetrics) {
	_jsii_.InvokeVoid(
		b,
		"addMetric",
		[]interface{}{metric},
	)
}

// Subscribes a destination to receive notifications when an object is created in the bucket.
//
// This is identical to calling
// `onEvent(EventType.ObjectCreated)`.
// Experimental.
func (b *jsiiProxy_Bucket) AddObjectCreatedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		b,
		"addObjectCreatedNotification",
		args,
	)
}

// Subscribes a destination to receive notifications when an object is removed from the bucket.
//
// This is identical to calling
// `onEvent(EventType.ObjectRemoved)`.
// Experimental.
func (b *jsiiProxy_Bucket) AddObjectRemovedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		b,
		"addObjectRemovedNotification",
		args,
	)
}

// Adds a statement to the resource policy for a principal (i.e. account/role/service) to perform actions on this bucket and/or it's contents. Use `bucketArn` and `arnForObjects(keys)` to obtain ARNs for this bucket or objects.
// Experimental.
func (b *jsiiProxy_Bucket) AddToResourcePolicy(permission awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		b,
		"addToResourcePolicy",
		[]interface{}{permission},
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
func (b *jsiiProxy_Bucket) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Returns an ARN that represents all objects within the bucket that match the key pattern specified.
//
// To represent all keys, specify ``"*"``.
//
// If you need to specify a keyPattern with multiple components, concatenate them into a single string, e.g.:
//
//    arnForObjects(`home/${team}/${user}/*`)
// Experimental.
func (b *jsiiProxy_Bucket) ArnForObjects(keyPattern *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"arnForObjects",
		[]interface{}{keyPattern},
		&returns,
	)

	return returns
}

// Experimental.
func (b *jsiiProxy_Bucket) GeneratePhysicalName() *string {
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
func (b *jsiiProxy_Bucket) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (b *jsiiProxy_Bucket) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grants s3:DeleteObject* permission to an IAM principal for objects in this bucket.
// Experimental.
func (b *jsiiProxy_Bucket) GrantDelete(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantDelete",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Allows unrestricted access to objects from this bucket.
//
// IMPORTANT: This permission allows anyone to perform actions on S3 objects
// in this bucket, which is useful for when you configure your bucket as a
// website and want everyone to be able to read objects in the bucket without
// needing to authenticate.
//
// Without arguments, this method will grant read ("s3:GetObject") access to
// all objects ("*") in the bucket.
//
// The method returns the `iam.Grant` object, which can then be modified
// as needed. For example, you can add a condition that will restrict access only
// to an IPv4 range like this:
//
//      const grant = bucket.grantPublicAccess();
//      grant.resourceStatement!.addCondition(‘IpAddress’, { “aws:SourceIp”: “54.240.143.0/24” });
// Experimental.
func (b *jsiiProxy_Bucket) GrantPublicAccess(keyPrefix *string, allowedActions ...*string) awsiam.Grant {
	args := []interface{}{keyPrefix}
	for _, a := range allowedActions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantPublicAccess",
		args,
		&returns,
	)

	return returns
}

// Grants s3:PutObject* and s3:Abort* permissions for this bucket to an IAM principal.
//
// If encryption is used, permission to use the key to encrypt the contents
// of written files will also be granted to the same principal.
// Experimental.
func (b *jsiiProxy_Bucket) GrantPut(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantPut",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grant the given IAM identity permissions to modify the ACLs of objects in the given Bucket.
//
// If your application has the '@aws-cdk/aws-s3:grantWriteWithoutAcl' feature flag set,
// calling {@link grantWrite} or {@link grantReadWrite} no longer grants permissions to modify the ACLs of the objects;
// in this case, if you need to modify object ACLs, call this method explicitly.
// Experimental.
func (b *jsiiProxy_Bucket) GrantPutAcl(identity awsiam.IGrantable, objectsKeyPattern *string) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantPutAcl",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grant read permissions for this bucket and it's contents to an IAM principal (Role/Group/User).
//
// If encryption is used, permission to use the key to decrypt the contents
// of the bucket will also be granted to the same principal.
// Experimental.
func (b *jsiiProxy_Bucket) GrantRead(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantRead",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grants read/write permissions for this bucket and it's contents to an IAM principal (Role/Group/User).
//
// If an encryption key is used, permission to use the key for
// encrypt/decrypt will also be granted.
//
// Before CDK version 1.85.0, this method granted the `s3:PutObject*` permission that included `s3:PutObjectAcl`,
// which could be used to grant read/write object access to IAM principals in other accounts.
// If you want to get rid of that behavior, update your CDK version to 1.85.0 or later,
// and make sure the `@aws-cdk/aws-s3:grantWriteWithoutAcl` feature flag is set to `true`
// in the `context` key of your cdk.json file.
// If you've already updated, but still need the principal to have permissions to modify the ACLs,
// use the {@link grantPutAcl} method.
// Experimental.
func (b *jsiiProxy_Bucket) GrantReadWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantReadWrite",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grant write permissions to this bucket to an IAM principal.
//
// If encryption is used, permission to use the key to encrypt the contents
// of written files will also be granted to the same principal.
//
// Before CDK version 1.85.0, this method granted the `s3:PutObject*` permission that included `s3:PutObjectAcl`,
// which could be used to grant read/write object access to IAM principals in other accounts.
// If you want to get rid of that behavior, update your CDK version to 1.85.0 or later,
// and make sure the `@aws-cdk/aws-s3:grantWriteWithoutAcl` feature flag is set to `true`
// in the `context` key of your cdk.json file.
// If you've already updated, but still need the principal to have permissions to modify the ACLs,
// use the {@link grantPutAcl} method.
// Experimental.
func (b *jsiiProxy_Bucket) GrantWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantWrite",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Define a CloudWatch event that triggers when something happens to this repository.
//
// Requires that there exists at least one CloudTrail Trail in your account
// that captures the event. This method will not create the Trail.
// Experimental.
func (b *jsiiProxy_Bucket) OnCloudTrailEvent(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		b,
		"onCloudTrailEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an AWS CloudWatch event that triggers when an object is uploaded to the specified paths (keys) in this bucket using the PutObject API call.
//
// Note that some tools like `aws s3 cp` will automatically use either
// PutObject or the multipart upload API depending on the file size,
// so using `onCloudTrailWriteObject` may be preferable.
//
// Requires that there exists at least one CloudTrail Trail in your account
// that captures the event. This method will not create the Trail.
// Experimental.
func (b *jsiiProxy_Bucket) OnCloudTrailPutObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		b,
		"onCloudTrailPutObject",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an AWS CloudWatch event that triggers when an object at the specified paths (keys) in this bucket are written to.
//
// This includes
// the events PutObject, CopyObject, and CompleteMultipartUpload.
//
// Note that some tools like `aws s3 cp` will automatically use either
// PutObject or the multipart upload API depending on the file size,
// so using this method may be preferable to `onCloudTrailPutObject`.
//
// Requires that there exists at least one CloudTrail Trail in your account
// that captures the event. This method will not create the Trail.
// Experimental.
func (b *jsiiProxy_Bucket) OnCloudTrailWriteObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		b,
		"onCloudTrailWriteObject",
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
func (b *jsiiProxy_Bucket) OnPrepare() {
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
func (b *jsiiProxy_Bucket) OnSynthesize(session constructs.ISynthesisSession) {
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
func (b *jsiiProxy_Bucket) OnValidate() *[]*string {
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
func (b *jsiiProxy_Bucket) Prepare() {
	_jsii_.InvokeVoid(
		b,
		"prepare",
		nil, // no parameters
	)
}

// The S3 URL of an S3 object.
//
// For example:
//
// Returns: an ObjectS3Url token
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_Bucket) S3UrlForObject(key *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"s3UrlForObject",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_Bucket) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_Bucket) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The https URL of an S3 object.
//
// Specify `regional: false` at the options
// for non-regional URLs. For example:
//
// Returns: an ObjectS3Url token
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_Bucket) UrlForObject(key *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"urlForObject",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
// Experimental.
func (b *jsiiProxy_Bucket) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The virtual hosted-style URL of an S3 object.
//
// Specify `regional: false` at
// the options for non-regional URL. For example:
//
// Returns: an ObjectS3Url token
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_Bucket) VirtualHostedUrlForObject(key *string, options *VirtualHostedStyleUrlOptions) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"virtualHostedUrlForObject",
		[]interface{}{key, options},
		&returns,
	)

	return returns
}

// Default bucket access control types.
// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html
//
// Experimental.
type BucketAccessControl string

const (
	BucketAccessControl_PRIVATE BucketAccessControl = "PRIVATE"
	BucketAccessControl_PUBLIC_READ BucketAccessControl = "PUBLIC_READ"
	BucketAccessControl_PUBLIC_READ_WRITE BucketAccessControl = "PUBLIC_READ_WRITE"
	BucketAccessControl_AUTHENTICATED_READ BucketAccessControl = "AUTHENTICATED_READ"
	BucketAccessControl_LOG_DELIVERY_WRITE BucketAccessControl = "LOG_DELIVERY_WRITE"
	BucketAccessControl_BUCKET_OWNER_READ BucketAccessControl = "BUCKET_OWNER_READ"
	BucketAccessControl_BUCKET_OWNER_FULL_CONTROL BucketAccessControl = "BUCKET_OWNER_FULL_CONTROL"
	BucketAccessControl_AWS_EXEC_READ BucketAccessControl = "AWS_EXEC_READ"
)

// A reference to a bucket.
//
// The easiest way to instantiate is to call
// `bucket.export()`. Then, the consumer can use `Bucket.import(this, ref)` and
// get a `Bucket`.
// Experimental.
type BucketAttributes struct {
	// The account this existing bucket belongs to.
	// Experimental.
	Account *string `json:"account"`
	// The ARN of the bucket.
	//
	// At least one of bucketArn or bucketName must be
	// defined in order to initialize a bucket ref.
	// Experimental.
	BucketArn *string `json:"bucketArn"`
	// The domain name of the bucket.
	// Experimental.
	BucketDomainName *string `json:"bucketDomainName"`
	// The IPv6 DNS name of the specified bucket.
	// Experimental.
	BucketDualStackDomainName *string `json:"bucketDualStackDomainName"`
	// The name of the bucket.
	//
	// If the underlying value of ARN is a string, the
	// name will be parsed from the ARN. Otherwise, the name is optional, but
	// some features that require the bucket name such as auto-creating a bucket
	// policy, won't work.
	// Experimental.
	BucketName *string `json:"bucketName"`
	// The regional domain name of the specified bucket.
	// Experimental.
	BucketRegionalDomainName *string `json:"bucketRegionalDomainName"`
	// The format of the website URL of the bucket.
	//
	// This should be true for
	// regions launched since 2014.
	// Experimental.
	BucketWebsiteNewUrlFormat *bool `json:"bucketWebsiteNewUrlFormat"`
	// The website URL of the bucket (if static web hosting is enabled).
	// Experimental.
	BucketWebsiteUrl *string `json:"bucketWebsiteUrl"`
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// If this bucket has been configured for static website hosting.
	// Experimental.
	IsWebsite *bool `json:"isWebsite"`
	// The region this existing bucket is in.
	// Experimental.
	Region *string `json:"region"`
}

// Represents an S3 Bucket.
//
// Buckets can be either defined within this stack:
//
//    new Bucket(this, 'MyBucket', { props });
//
// Or imported from an existing bucket:
//
//    Bucket.import(this, 'MyImportedBucket', { bucketArn: ... });
//
// You can also export a bucket and import it into another stack:
//
//    const ref = myBucket.export();
//    Bucket.import(this, 'MyImportedBucket', ref);
// Experimental.
type BucketBase interface {
	awscdk.Resource
	IBucket
	AutoCreatePolicy() *bool
	SetAutoCreatePolicy(val *bool)
	BucketArn() *string
	BucketDomainName() *string
	BucketDualStackDomainName() *string
	BucketName() *string
	BucketRegionalDomainName() *string
	BucketWebsiteDomainName() *string
	BucketWebsiteUrl() *string
	DisallowPublicAccess() *bool
	SetDisallowPublicAccess(val *bool)
	EncryptionKey() awskms.IKey
	Env() *awscdk.ResourceEnvironment
	IsWebsite() *bool
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Policy() BucketPolicy
	SetPolicy(val BucketPolicy)
	Stack() awscdk.Stack
	AddEventNotification(event EventType, dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	AddObjectCreatedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	AddObjectRemovedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	AddToResourcePolicy(permission awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ArnForObjects(keyPattern *string) *string
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantDelete(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantPublicAccess(keyPrefix *string, allowedActions ...*string) awsiam.Grant
	GrantPut(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantPutAcl(identity awsiam.IGrantable, objectsKeyPattern *string) awsiam.Grant
	GrantRead(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantReadWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	GrantWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	OnCloudTrailEvent(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	OnCloudTrailPutObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	OnCloudTrailWriteObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	S3UrlForObject(key *string) *string
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	UrlForObject(key *string) *string
	Validate() *[]*string
	VirtualHostedUrlForObject(key *string, options *VirtualHostedStyleUrlOptions) *string
}

// The jsii proxy struct for BucketBase
type jsiiProxy_BucketBase struct {
	internal.Type__awscdkResource
	jsiiProxy_IBucket
}

func (j *jsiiProxy_BucketBase) AutoCreatePolicy() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"autoCreatePolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) BucketArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) BucketDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) BucketDualStackDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketDualStackDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) BucketName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) BucketRegionalDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketRegionalDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) BucketWebsiteDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketWebsiteDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) BucketWebsiteUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketWebsiteUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) DisallowPublicAccess() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"disallowPublicAccess",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) IsWebsite() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isWebsite",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) Policy() BucketPolicy {
	var returns BucketPolicy
	_jsii_.Get(
		j,
		"policy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewBucketBase_Override(b BucketBase, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.BucketBase",
		[]interface{}{scope, id, props},
		b,
	)
}

func (j *jsiiProxy_BucketBase) SetAutoCreatePolicy(val *bool) {
	_jsii_.Set(
		j,
		"autoCreatePolicy",
		val,
	)
}

func (j *jsiiProxy_BucketBase) SetDisallowPublicAccess(val *bool) {
	_jsii_.Set(
		j,
		"disallowPublicAccess",
		val,
	)
}

func (j *jsiiProxy_BucketBase) SetPolicy(val BucketPolicy) {
	_jsii_.Set(
		j,
		"policy",
		val,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BucketBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.BucketBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func BucketBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.BucketBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a bucket notification event destination.
//
// TODO: EXAMPLE
//
// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html
//
// Experimental.
func (b *jsiiProxy_BucketBase) AddEventNotification(event EventType, dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{event, dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		b,
		"addEventNotification",
		args,
	)
}

// Subscribes a destination to receive notifications when an object is created in the bucket.
//
// This is identical to calling
// `onEvent(EventType.ObjectCreated)`.
// Experimental.
func (b *jsiiProxy_BucketBase) AddObjectCreatedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		b,
		"addObjectCreatedNotification",
		args,
	)
}

// Subscribes a destination to receive notifications when an object is removed from the bucket.
//
// This is identical to calling
// `onEvent(EventType.ObjectRemoved)`.
// Experimental.
func (b *jsiiProxy_BucketBase) AddObjectRemovedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		b,
		"addObjectRemovedNotification",
		args,
	)
}

// Adds a statement to the resource policy for a principal (i.e. account/role/service) to perform actions on this bucket and/or it's contents. Use `bucketArn` and `arnForObjects(keys)` to obtain ARNs for this bucket or objects.
// Experimental.
func (b *jsiiProxy_BucketBase) AddToResourcePolicy(permission awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		b,
		"addToResourcePolicy",
		[]interface{}{permission},
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
func (b *jsiiProxy_BucketBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Returns an ARN that represents all objects within the bucket that match the key pattern specified.
//
// To represent all keys, specify ``"*"``.
//
// If you need to specify a keyPattern with multiple components, concatenate them into a single string, e.g.:
//
//    arnForObjects(`home/${team}/${user}/*`)
// Experimental.
func (b *jsiiProxy_BucketBase) ArnForObjects(keyPattern *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"arnForObjects",
		[]interface{}{keyPattern},
		&returns,
	)

	return returns
}

// Experimental.
func (b *jsiiProxy_BucketBase) GeneratePhysicalName() *string {
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
func (b *jsiiProxy_BucketBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (b *jsiiProxy_BucketBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Grants s3:DeleteObject* permission to an IAM principal for objects in this bucket.
// Experimental.
func (b *jsiiProxy_BucketBase) GrantDelete(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantDelete",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Allows unrestricted access to objects from this bucket.
//
// IMPORTANT: This permission allows anyone to perform actions on S3 objects
// in this bucket, which is useful for when you configure your bucket as a
// website and want everyone to be able to read objects in the bucket without
// needing to authenticate.
//
// Without arguments, this method will grant read ("s3:GetObject") access to
// all objects ("*") in the bucket.
//
// The method returns the `iam.Grant` object, which can then be modified
// as needed. For example, you can add a condition that will restrict access only
// to an IPv4 range like this:
//
//      const grant = bucket.grantPublicAccess();
//      grant.resourceStatement!.addCondition(‘IpAddress’, { “aws:SourceIp”: “54.240.143.0/24” });
// Experimental.
func (b *jsiiProxy_BucketBase) GrantPublicAccess(keyPrefix *string, allowedActions ...*string) awsiam.Grant {
	args := []interface{}{keyPrefix}
	for _, a := range allowedActions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantPublicAccess",
		args,
		&returns,
	)

	return returns
}

// Grants s3:PutObject* and s3:Abort* permissions for this bucket to an IAM principal.
//
// If encryption is used, permission to use the key to encrypt the contents
// of written files will also be granted to the same principal.
// Experimental.
func (b *jsiiProxy_BucketBase) GrantPut(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantPut",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grant the given IAM identity permissions to modify the ACLs of objects in the given Bucket.
//
// If your application has the '@aws-cdk/aws-s3:grantWriteWithoutAcl' feature flag set,
// calling {@link grantWrite} or {@link grantReadWrite} no longer grants permissions to modify the ACLs of the objects;
// in this case, if you need to modify object ACLs, call this method explicitly.
// Experimental.
func (b *jsiiProxy_BucketBase) GrantPutAcl(identity awsiam.IGrantable, objectsKeyPattern *string) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantPutAcl",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grant read permissions for this bucket and it's contents to an IAM principal (Role/Group/User).
//
// If encryption is used, permission to use the key to decrypt the contents
// of the bucket will also be granted to the same principal.
// Experimental.
func (b *jsiiProxy_BucketBase) GrantRead(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantRead",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grants read/write permissions for this bucket and it's contents to an IAM principal (Role/Group/User).
//
// If an encryption key is used, permission to use the key for
// encrypt/decrypt will also be granted.
//
// Before CDK version 1.85.0, this method granted the `s3:PutObject*` permission that included `s3:PutObjectAcl`,
// which could be used to grant read/write object access to IAM principals in other accounts.
// If you want to get rid of that behavior, update your CDK version to 1.85.0 or later,
// and make sure the `@aws-cdk/aws-s3:grantWriteWithoutAcl` feature flag is set to `true`
// in the `context` key of your cdk.json file.
// If you've already updated, but still need the principal to have permissions to modify the ACLs,
// use the {@link grantPutAcl} method.
// Experimental.
func (b *jsiiProxy_BucketBase) GrantReadWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantReadWrite",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Grant write permissions to this bucket to an IAM principal.
//
// If encryption is used, permission to use the key to encrypt the contents
// of written files will also be granted to the same principal.
//
// Before CDK version 1.85.0, this method granted the `s3:PutObject*` permission that included `s3:PutObjectAcl`,
// which could be used to grant read/write object access to IAM principals in other accounts.
// If you want to get rid of that behavior, update your CDK version to 1.85.0 or later,
// and make sure the `@aws-cdk/aws-s3:grantWriteWithoutAcl` feature flag is set to `true`
// in the `context` key of your cdk.json file.
// If you've already updated, but still need the principal to have permissions to modify the ACLs,
// use the {@link grantPutAcl} method.
// Experimental.
func (b *jsiiProxy_BucketBase) GrantWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		b,
		"grantWrite",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

// Define a CloudWatch event that triggers when something happens to this repository.
//
// Requires that there exists at least one CloudTrail Trail in your account
// that captures the event. This method will not create the Trail.
// Experimental.
func (b *jsiiProxy_BucketBase) OnCloudTrailEvent(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		b,
		"onCloudTrailEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an AWS CloudWatch event that triggers when an object is uploaded to the specified paths (keys) in this bucket using the PutObject API call.
//
// Note that some tools like `aws s3 cp` will automatically use either
// PutObject or the multipart upload API depending on the file size,
// so using `onCloudTrailWriteObject` may be preferable.
//
// Requires that there exists at least one CloudTrail Trail in your account
// that captures the event. This method will not create the Trail.
// Experimental.
func (b *jsiiProxy_BucketBase) OnCloudTrailPutObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		b,
		"onCloudTrailPutObject",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an AWS CloudWatch event that triggers when an object at the specified paths (keys) in this bucket are written to.
//
// This includes
// the events PutObject, CopyObject, and CompleteMultipartUpload.
//
// Note that some tools like `aws s3 cp` will automatically use either
// PutObject or the multipart upload API depending on the file size,
// so using this method may be preferable to `onCloudTrailPutObject`.
//
// Requires that there exists at least one CloudTrail Trail in your account
// that captures the event. This method will not create the Trail.
// Experimental.
func (b *jsiiProxy_BucketBase) OnCloudTrailWriteObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		b,
		"onCloudTrailWriteObject",
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
func (b *jsiiProxy_BucketBase) OnPrepare() {
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
func (b *jsiiProxy_BucketBase) OnSynthesize(session constructs.ISynthesisSession) {
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
func (b *jsiiProxy_BucketBase) OnValidate() *[]*string {
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
func (b *jsiiProxy_BucketBase) Prepare() {
	_jsii_.InvokeVoid(
		b,
		"prepare",
		nil, // no parameters
	)
}

// The S3 URL of an S3 object.
//
// For example:
//
// Returns: an ObjectS3Url token
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_BucketBase) S3UrlForObject(key *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"s3UrlForObject",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_BucketBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BucketBase) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The https URL of an S3 object.
//
// Specify `regional: false` at the options
// for non-regional URLs. For example:
//
// Returns: an ObjectS3Url token
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_BucketBase) UrlForObject(key *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"urlForObject",
		[]interface{}{key},
		&returns,
	)

	return returns
}

// Validate the current construct.
//
// This method can be implemented by derived constructs in order to perform
// validation logic. It is called on all constructs before synthesis.
// Experimental.
func (b *jsiiProxy_BucketBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The virtual hosted-style URL of an S3 object.
//
// Specify `regional: false` at
// the options for non-regional URL. For example:
//
// Returns: an ObjectS3Url token
//
// TODO: EXAMPLE
//
// Experimental.
func (b *jsiiProxy_BucketBase) VirtualHostedUrlForObject(key *string, options *VirtualHostedStyleUrlOptions) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"virtualHostedUrlForObject",
		[]interface{}{key, options},
		&returns,
	)

	return returns
}

// What kind of server-side encryption to apply to this bucket.
// Experimental.
type BucketEncryption string

const (
	BucketEncryption_UNENCRYPTED BucketEncryption = "UNENCRYPTED"
	BucketEncryption_KMS_MANAGED BucketEncryption = "KMS_MANAGED"
	BucketEncryption_S3_MANAGED BucketEncryption = "S3_MANAGED"
	BucketEncryption_KMS BucketEncryption = "KMS"
)

// Specifies a metrics configuration for the CloudWatch request metrics from an Amazon S3 bucket.
// Experimental.
type BucketMetrics struct {
	// The ID used to identify the metrics configuration.
	// Experimental.
	Id *string `json:"id"`
	// The prefix that an object must have to be included in the metrics results.
	// Experimental.
	Prefix *string `json:"prefix"`
	// Specifies a list of tag filters to use as a metrics configuration filter.
	//
	// The metrics configuration includes only objects that meet the filter's criteria.
	// Experimental.
	TagFilters *map[string]interface{} `json:"tagFilters"`
}

// Represents the properties of a notification destination.
// Experimental.
type BucketNotificationDestinationConfig struct {
	// The ARN of the destination (i.e. Lambda, SNS, SQS).
	// Experimental.
	Arn *string `json:"arn"`
	// The notification type.
	// Experimental.
	Type BucketNotificationDestinationType `json:"type"`
	// Any additional dependencies that should be resolved before the bucket notification can be configured (for example, the SNS Topic Policy resource).
	// Experimental.
	Dependencies *[]awscdk.IDependable `json:"dependencies"`
}

// Supported types of notification destinations.
// Experimental.
type BucketNotificationDestinationType string

const (
	BucketNotificationDestinationType_LAMBDA BucketNotificationDestinationType = "LAMBDA"
	BucketNotificationDestinationType_QUEUE BucketNotificationDestinationType = "QUEUE"
	BucketNotificationDestinationType_TOPIC BucketNotificationDestinationType = "TOPIC"
)

// Applies an Amazon S3 bucket policy to an Amazon S3 bucket.
// Experimental.
type BucketPolicy interface {
	awscdk.Resource
	Document() awsiam.PolicyDocument
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(removalPolicy awscdk.RemovalPolicy)
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

// The jsii proxy struct for BucketPolicy
type jsiiProxy_BucketPolicy struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_BucketPolicy) Document() awsiam.PolicyDocument {
	var returns awsiam.PolicyDocument
	_jsii_.Get(
		j,
		"document",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketPolicy) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketPolicy) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BucketPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewBucketPolicy(scope constructs.Construct, id *string, props *BucketPolicyProps) BucketPolicy {
	_init_.Initialize()

	j := jsiiProxy_BucketPolicy{}

	_jsii_.Create(
		"monocdk.aws_s3.BucketPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewBucketPolicy_Override(b BucketPolicy, scope constructs.Construct, id *string, props *BucketPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.BucketPolicy",
		[]interface{}{scope, id, props},
		b,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BucketPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.BucketPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func BucketPolicy_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.BucketPolicy",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Sets the removal policy for the BucketPolicy.
// Experimental.
func (b *jsiiProxy_BucketPolicy) ApplyRemovalPolicy(removalPolicy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{removalPolicy},
	)
}

// Experimental.
func (b *jsiiProxy_BucketPolicy) GeneratePhysicalName() *string {
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
func (b *jsiiProxy_BucketPolicy) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (b *jsiiProxy_BucketPolicy) GetResourceNameAttribute(nameAttr *string) *string {
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
func (b *jsiiProxy_BucketPolicy) OnPrepare() {
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
func (b *jsiiProxy_BucketPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (b *jsiiProxy_BucketPolicy) OnValidate() *[]*string {
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
func (b *jsiiProxy_BucketPolicy) Prepare() {
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
func (b *jsiiProxy_BucketPolicy) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BucketPolicy) ToString() *string {
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
func (b *jsiiProxy_BucketPolicy) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type BucketPolicyProps struct {
	// The Amazon S3 bucket that the policy applies to.
	// Experimental.
	Bucket IBucket `json:"bucket"`
	// Policy to apply when the policy is removed from this stack.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
}

// Experimental.
type BucketProps struct {
	// Specifies a canned ACL that grants predefined permissions to the bucket.
	// Experimental.
	AccessControl BucketAccessControl `json:"accessControl"`
	// Whether all objects should be automatically deleted when the bucket is removed from the stack or when the stack is deleted.
	//
	// Requires the `removalPolicy` to be set to `RemovalPolicy.DESTROY`.
	// Experimental.
	AutoDeleteObjects *bool `json:"autoDeleteObjects"`
	// The block public access configuration of this bucket.
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/access-control-block-public-access.html
	//
	// Experimental.
	BlockPublicAccess BlockPublicAccess `json:"blockPublicAccess"`
	// Specifies whether Amazon S3 should use an S3 Bucket Key with server-side encryption using KMS (SSE-KMS) for new objects in the bucket.
	//
	// Only relevant, when Encryption is set to {@link BucketEncryption.KMS}
	// Experimental.
	BucketKeyEnabled *bool `json:"bucketKeyEnabled"`
	// Physical name of this bucket.
	// Experimental.
	BucketName *string `json:"bucketName"`
	// The CORS configuration of this bucket.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-cors.html
	//
	// Experimental.
	Cors *[]*CorsRule `json:"cors"`
	// The kind of server-side encryption to apply to this bucket.
	//
	// If you choose KMS, you can specify a KMS key via `encryptionKey`. If
	// encryption key is not specified, a key will automatically be created.
	// Experimental.
	Encryption BucketEncryption `json:"encryption"`
	// External KMS key to use for bucket encryption.
	//
	// The 'encryption' property must be either not specified or set to "Kms".
	// An error will be emitted if encryption is set to "Unencrypted" or
	// "Managed".
	// Experimental.
	EncryptionKey awskms.IKey `json:"encryptionKey"`
	// Enforces SSL for requests.
	//
	// S3.5 of the AWS Foundational Security Best Practices Regarding S3.
	// See: https://docs.aws.amazon.com/config/latest/developerguide/s3-bucket-ssl-requests-only.html
	//
	// Experimental.
	EnforceSSL *bool `json:"enforceSSL"`
	// The inventory configuration of the bucket.
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/storage-inventory.html
	//
	// Experimental.
	Inventories *[]*Inventory `json:"inventories"`
	// Rules that define how Amazon S3 manages objects during their lifetime.
	// Experimental.
	LifecycleRules *[]*LifecycleRule `json:"lifecycleRules"`
	// The metrics configuration of this bucket.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-metricsconfiguration.html
	//
	// Experimental.
	Metrics *[]*BucketMetrics `json:"metrics"`
	// The objectOwnership of the bucket.
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/about-object-ownership.html
	//
	// Experimental.
	ObjectOwnership ObjectOwnership `json:"objectOwnership"`
	// Grants public read access to all objects in the bucket.
	//
	// Similar to calling `bucket.grantPublicAccess()`
	// Experimental.
	PublicReadAccess *bool `json:"publicReadAccess"`
	// Policy to apply when the bucket is removed from this stack.
	// Experimental.
	RemovalPolicy awscdk.RemovalPolicy `json:"removalPolicy"`
	// Destination bucket for the server access logs.
	// Experimental.
	ServerAccessLogsBucket IBucket `json:"serverAccessLogsBucket"`
	// Optional log file prefix to use for the bucket's access logs.
	//
	// If defined without "serverAccessLogsBucket", enables access logs to current bucket with this prefix.
	// Experimental.
	ServerAccessLogsPrefix *string `json:"serverAccessLogsPrefix"`
	// Whether this bucket should have versioning turned on or not.
	// Experimental.
	Versioned *bool `json:"versioned"`
	// The name of the error document (e.g. "404.html") for the website. `websiteIndexDocument` must also be set if this is set.
	// Experimental.
	WebsiteErrorDocument *string `json:"websiteErrorDocument"`
	// The name of the index document (e.g. "index.html") for the website. Enables static website hosting for this bucket.
	// Experimental.
	WebsiteIndexDocument *string `json:"websiteIndexDocument"`
	// Specifies the redirect behavior of all requests to a website endpoint of a bucket.
	//
	// If you specify this property, you can't specify "websiteIndexDocument", "websiteErrorDocument" nor , "websiteRoutingRules".
	// Experimental.
	WebsiteRedirect *RedirectTarget `json:"websiteRedirect"`
	// Rules that define when a redirect is applied and the redirect behavior.
	// Experimental.
	WebsiteRoutingRules *[]*RoutingRule `json:"websiteRoutingRules"`
}

// A CloudFormation `AWS::S3::AccessPoint`.
type CfnAccessPoint interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrNetworkOrigin() *string
	Bucket() *string
	SetBucket(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Policy() interface{}
	SetPolicy(val interface{})
	PublicAccessBlockConfiguration() interface{}
	SetPublicAccessBlockConfiguration(val interface{})
	Ref() *string
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	VpcConfiguration() interface{}
	SetVpcConfiguration(val interface{})
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

// The jsii proxy struct for CfnAccessPoint
type jsiiProxy_CfnAccessPoint struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAccessPoint) AttrNetworkOrigin() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrNetworkOrigin",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) Bucket() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucket",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) Policy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) PublicAccessBlockConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"publicAccessBlockConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccessPoint) VpcConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"vpcConfiguration",
		&returns,
	)
	return returns
}


// Create a new `AWS::S3::AccessPoint`.
func NewCfnAccessPoint(scope awscdk.Construct, id *string, props *CfnAccessPointProps) CfnAccessPoint {
	_init_.Initialize()

	j := jsiiProxy_CfnAccessPoint{}

	_jsii_.Create(
		"monocdk.aws_s3.CfnAccessPoint",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::S3::AccessPoint`.
func NewCfnAccessPoint_Override(c CfnAccessPoint, scope awscdk.Construct, id *string, props *CfnAccessPointProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.CfnAccessPoint",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAccessPoint) SetBucket(val *string) {
	_jsii_.Set(
		j,
		"bucket",
		val,
	)
}

func (j *jsiiProxy_CfnAccessPoint) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnAccessPoint) SetPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"policy",
		val,
	)
}

func (j *jsiiProxy_CfnAccessPoint) SetPublicAccessBlockConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"publicAccessBlockConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnAccessPoint) SetVpcConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"vpcConfiguration",
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
func CfnAccessPoint_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnAccessPoint",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAccessPoint_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnAccessPoint",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAccessPoint_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnAccessPoint",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAccessPoint_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_s3.CfnAccessPoint",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAccessPoint) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAccessPoint) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAccessPoint) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAccessPoint) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAccessPoint) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAccessPoint) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAccessPoint) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAccessPoint) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAccessPoint) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAccessPoint) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAccessPoint) OnPrepare() {
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
func (c *jsiiProxy_CfnAccessPoint) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAccessPoint) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAccessPoint) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAccessPoint) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAccessPoint) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAccessPoint) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAccessPoint) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAccessPoint) ToString() *string {
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
func (c *jsiiProxy_CfnAccessPoint) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAccessPoint) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnAccessPoint_PublicAccessBlockConfigurationProperty struct {
	// `CfnAccessPoint.PublicAccessBlockConfigurationProperty.BlockPublicAcls`.
	BlockPublicAcls interface{} `json:"blockPublicAcls"`
	// `CfnAccessPoint.PublicAccessBlockConfigurationProperty.BlockPublicPolicy`.
	BlockPublicPolicy interface{} `json:"blockPublicPolicy"`
	// `CfnAccessPoint.PublicAccessBlockConfigurationProperty.IgnorePublicAcls`.
	IgnorePublicAcls interface{} `json:"ignorePublicAcls"`
	// `CfnAccessPoint.PublicAccessBlockConfigurationProperty.RestrictPublicBuckets`.
	RestrictPublicBuckets interface{} `json:"restrictPublicBuckets"`
}

type CfnAccessPoint_VpcConfigurationProperty struct {
	// `CfnAccessPoint.VpcConfigurationProperty.VpcId`.
	VpcId *string `json:"vpcId"`
}

// Properties for defining a `AWS::S3::AccessPoint`.
type CfnAccessPointProps struct {
	// `AWS::S3::AccessPoint.Bucket`.
	Bucket *string `json:"bucket"`
	// `AWS::S3::AccessPoint.Name`.
	Name *string `json:"name"`
	// `AWS::S3::AccessPoint.Policy`.
	Policy interface{} `json:"policy"`
	// `AWS::S3::AccessPoint.PublicAccessBlockConfiguration`.
	PublicAccessBlockConfiguration interface{} `json:"publicAccessBlockConfiguration"`
	// `AWS::S3::AccessPoint.VpcConfiguration`.
	VpcConfiguration interface{} `json:"vpcConfiguration"`
}

// A CloudFormation `AWS::S3::Bucket`.
type CfnBucket interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AccelerateConfiguration() interface{}
	SetAccelerateConfiguration(val interface{})
	AccessControl() *string
	SetAccessControl(val *string)
	AnalyticsConfigurations() interface{}
	SetAnalyticsConfigurations(val interface{})
	AttrArn() *string
	AttrDomainName() *string
	AttrDualStackDomainName() *string
	AttrRegionalDomainName() *string
	AttrWebsiteUrl() *string
	BucketEncryption() interface{}
	SetBucketEncryption(val interface{})
	BucketName() *string
	SetBucketName(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CorsConfiguration() interface{}
	SetCorsConfiguration(val interface{})
	CreationStack() *[]*string
	IntelligentTieringConfigurations() interface{}
	SetIntelligentTieringConfigurations(val interface{})
	InventoryConfigurations() interface{}
	SetInventoryConfigurations(val interface{})
	LifecycleConfiguration() interface{}
	SetLifecycleConfiguration(val interface{})
	LoggingConfiguration() interface{}
	SetLoggingConfiguration(val interface{})
	LogicalId() *string
	MetricsConfigurations() interface{}
	SetMetricsConfigurations(val interface{})
	Node() awscdk.ConstructNode
	NotificationConfiguration() interface{}
	SetNotificationConfiguration(val interface{})
	ObjectLockConfiguration() interface{}
	SetObjectLockConfiguration(val interface{})
	ObjectLockEnabled() interface{}
	SetObjectLockEnabled(val interface{})
	OwnershipControls() interface{}
	SetOwnershipControls(val interface{})
	PublicAccessBlockConfiguration() interface{}
	SetPublicAccessBlockConfiguration(val interface{})
	Ref() *string
	ReplicationConfiguration() interface{}
	SetReplicationConfiguration(val interface{})
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	UpdatedProperites() *map[string]interface{}
	VersioningConfiguration() interface{}
	SetVersioningConfiguration(val interface{})
	WebsiteConfiguration() interface{}
	SetWebsiteConfiguration(val interface{})
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

// The jsii proxy struct for CfnBucket
type jsiiProxy_CfnBucket struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnBucket) AccelerateConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accelerateConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) AccessControl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"accessControl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) AnalyticsConfigurations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"analyticsConfigurations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) AttrArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) AttrDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) AttrDualStackDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDualStackDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) AttrRegionalDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRegionalDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) AttrWebsiteUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrWebsiteUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) BucketEncryption() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"bucketEncryption",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) BucketName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) CorsConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"corsConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) IntelligentTieringConfigurations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"intelligentTieringConfigurations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) InventoryConfigurations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"inventoryConfigurations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) LifecycleConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"lifecycleConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) LoggingConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"loggingConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) MetricsConfigurations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"metricsConfigurations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) NotificationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"notificationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) ObjectLockConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"objectLockConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) ObjectLockEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"objectLockEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) OwnershipControls() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ownershipControls",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) PublicAccessBlockConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"publicAccessBlockConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) ReplicationConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"replicationConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) VersioningConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"versioningConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucket) WebsiteConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"websiteConfiguration",
		&returns,
	)
	return returns
}


// Create a new `AWS::S3::Bucket`.
func NewCfnBucket(scope awscdk.Construct, id *string, props *CfnBucketProps) CfnBucket {
	_init_.Initialize()

	j := jsiiProxy_CfnBucket{}

	_jsii_.Create(
		"monocdk.aws_s3.CfnBucket",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::S3::Bucket`.
func NewCfnBucket_Override(c CfnBucket, scope awscdk.Construct, id *string, props *CfnBucketProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.CfnBucket",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnBucket) SetAccelerateConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"accelerateConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetAccessControl(val *string) {
	_jsii_.Set(
		j,
		"accessControl",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetAnalyticsConfigurations(val interface{}) {
	_jsii_.Set(
		j,
		"analyticsConfigurations",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetBucketEncryption(val interface{}) {
	_jsii_.Set(
		j,
		"bucketEncryption",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetBucketName(val *string) {
	_jsii_.Set(
		j,
		"bucketName",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetCorsConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"corsConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetIntelligentTieringConfigurations(val interface{}) {
	_jsii_.Set(
		j,
		"intelligentTieringConfigurations",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetInventoryConfigurations(val interface{}) {
	_jsii_.Set(
		j,
		"inventoryConfigurations",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetLifecycleConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"lifecycleConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetLoggingConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"loggingConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetMetricsConfigurations(val interface{}) {
	_jsii_.Set(
		j,
		"metricsConfigurations",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetNotificationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"notificationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetObjectLockConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"objectLockConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetObjectLockEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"objectLockEnabled",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetOwnershipControls(val interface{}) {
	_jsii_.Set(
		j,
		"ownershipControls",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetPublicAccessBlockConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"publicAccessBlockConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetReplicationConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"replicationConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetVersioningConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"versioningConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnBucket) SetWebsiteConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"websiteConfiguration",
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
func CfnBucket_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnBucket",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnBucket_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnBucket",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnBucket_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnBucket",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnBucket_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_s3.CfnBucket",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnBucket) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnBucket) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnBucket) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnBucket) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnBucket) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnBucket) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnBucket) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnBucket) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnBucket) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnBucket) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnBucket) OnPrepare() {
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
func (c *jsiiProxy_CfnBucket) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnBucket) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnBucket) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnBucket) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnBucket) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnBucket) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnBucket) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnBucket) ToString() *string {
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
func (c *jsiiProxy_CfnBucket) Validate() *[]*string {
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
func (c *jsiiProxy_CfnBucket) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnBucket_AbortIncompleteMultipartUploadProperty struct {
	// `CfnBucket.AbortIncompleteMultipartUploadProperty.DaysAfterInitiation`.
	DaysAfterInitiation *float64 `json:"daysAfterInitiation"`
}

type CfnBucket_AccelerateConfigurationProperty struct {
	// `CfnBucket.AccelerateConfigurationProperty.AccelerationStatus`.
	AccelerationStatus *string `json:"accelerationStatus"`
}

type CfnBucket_AccessControlTranslationProperty struct {
	// `CfnBucket.AccessControlTranslationProperty.Owner`.
	Owner *string `json:"owner"`
}

type CfnBucket_AnalyticsConfigurationProperty struct {
	// `CfnBucket.AnalyticsConfigurationProperty.Id`.
	Id *string `json:"id"`
	// `CfnBucket.AnalyticsConfigurationProperty.StorageClassAnalysis`.
	StorageClassAnalysis interface{} `json:"storageClassAnalysis"`
	// `CfnBucket.AnalyticsConfigurationProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnBucket.AnalyticsConfigurationProperty.TagFilters`.
	TagFilters interface{} `json:"tagFilters"`
}

type CfnBucket_BucketEncryptionProperty struct {
	// `CfnBucket.BucketEncryptionProperty.ServerSideEncryptionConfiguration`.
	ServerSideEncryptionConfiguration interface{} `json:"serverSideEncryptionConfiguration"`
}

type CfnBucket_CorsConfigurationProperty struct {
	// `CfnBucket.CorsConfigurationProperty.CorsRules`.
	CorsRules interface{} `json:"corsRules"`
}

type CfnBucket_CorsRuleProperty struct {
	// `CfnBucket.CorsRuleProperty.AllowedMethods`.
	AllowedMethods *[]*string `json:"allowedMethods"`
	// `CfnBucket.CorsRuleProperty.AllowedOrigins`.
	AllowedOrigins *[]*string `json:"allowedOrigins"`
	// `CfnBucket.CorsRuleProperty.AllowedHeaders`.
	AllowedHeaders *[]*string `json:"allowedHeaders"`
	// `CfnBucket.CorsRuleProperty.ExposedHeaders`.
	ExposedHeaders *[]*string `json:"exposedHeaders"`
	// `CfnBucket.CorsRuleProperty.Id`.
	Id *string `json:"id"`
	// `CfnBucket.CorsRuleProperty.MaxAge`.
	MaxAge *float64 `json:"maxAge"`
}

type CfnBucket_DataExportProperty struct {
	// `CfnBucket.DataExportProperty.Destination`.
	Destination interface{} `json:"destination"`
	// `CfnBucket.DataExportProperty.OutputSchemaVersion`.
	OutputSchemaVersion *string `json:"outputSchemaVersion"`
}

type CfnBucket_DefaultRetentionProperty struct {
	// `CfnBucket.DefaultRetentionProperty.Days`.
	Days *float64 `json:"days"`
	// `CfnBucket.DefaultRetentionProperty.Mode`.
	Mode *string `json:"mode"`
	// `CfnBucket.DefaultRetentionProperty.Years`.
	Years *float64 `json:"years"`
}

type CfnBucket_DeleteMarkerReplicationProperty struct {
	// `CfnBucket.DeleteMarkerReplicationProperty.Status`.
	Status *string `json:"status"`
}

type CfnBucket_DestinationProperty struct {
	// `CfnBucket.DestinationProperty.BucketArn`.
	BucketArn *string `json:"bucketArn"`
	// `CfnBucket.DestinationProperty.Format`.
	Format *string `json:"format"`
	// `CfnBucket.DestinationProperty.BucketAccountId`.
	BucketAccountId *string `json:"bucketAccountId"`
	// `CfnBucket.DestinationProperty.Prefix`.
	Prefix *string `json:"prefix"`
}

type CfnBucket_EncryptionConfigurationProperty struct {
	// `CfnBucket.EncryptionConfigurationProperty.ReplicaKmsKeyID`.
	ReplicaKmsKeyId *string `json:"replicaKmsKeyId"`
}

type CfnBucket_FilterRuleProperty struct {
	// `CfnBucket.FilterRuleProperty.Name`.
	Name *string `json:"name"`
	// `CfnBucket.FilterRuleProperty.Value`.
	Value *string `json:"value"`
}

type CfnBucket_IntelligentTieringConfigurationProperty struct {
	// `CfnBucket.IntelligentTieringConfigurationProperty.Id`.
	Id *string `json:"id"`
	// `CfnBucket.IntelligentTieringConfigurationProperty.Status`.
	Status *string `json:"status"`
	// `CfnBucket.IntelligentTieringConfigurationProperty.Tierings`.
	Tierings interface{} `json:"tierings"`
	// `CfnBucket.IntelligentTieringConfigurationProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnBucket.IntelligentTieringConfigurationProperty.TagFilters`.
	TagFilters interface{} `json:"tagFilters"`
}

type CfnBucket_InventoryConfigurationProperty struct {
	// `CfnBucket.InventoryConfigurationProperty.Destination`.
	Destination interface{} `json:"destination"`
	// `CfnBucket.InventoryConfigurationProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `CfnBucket.InventoryConfigurationProperty.Id`.
	Id *string `json:"id"`
	// `CfnBucket.InventoryConfigurationProperty.IncludedObjectVersions`.
	IncludedObjectVersions *string `json:"includedObjectVersions"`
	// `CfnBucket.InventoryConfigurationProperty.ScheduleFrequency`.
	ScheduleFrequency *string `json:"scheduleFrequency"`
	// `CfnBucket.InventoryConfigurationProperty.OptionalFields`.
	OptionalFields *[]*string `json:"optionalFields"`
	// `CfnBucket.InventoryConfigurationProperty.Prefix`.
	Prefix *string `json:"prefix"`
}

type CfnBucket_LambdaConfigurationProperty struct {
	// `CfnBucket.LambdaConfigurationProperty.Event`.
	Event *string `json:"event"`
	// `CfnBucket.LambdaConfigurationProperty.Function`.
	Function *string `json:"function"`
	// `CfnBucket.LambdaConfigurationProperty.Filter`.
	Filter interface{} `json:"filter"`
}

type CfnBucket_LifecycleConfigurationProperty struct {
	// `CfnBucket.LifecycleConfigurationProperty.Rules`.
	Rules interface{} `json:"rules"`
}

type CfnBucket_LoggingConfigurationProperty struct {
	// `CfnBucket.LoggingConfigurationProperty.DestinationBucketName`.
	DestinationBucketName *string `json:"destinationBucketName"`
	// `CfnBucket.LoggingConfigurationProperty.LogFilePrefix`.
	LogFilePrefix *string `json:"logFilePrefix"`
}

type CfnBucket_MetricsConfigurationProperty struct {
	// `CfnBucket.MetricsConfigurationProperty.Id`.
	Id *string `json:"id"`
	// `CfnBucket.MetricsConfigurationProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnBucket.MetricsConfigurationProperty.TagFilters`.
	TagFilters interface{} `json:"tagFilters"`
}

type CfnBucket_MetricsProperty struct {
	// `CfnBucket.MetricsProperty.Status`.
	Status *string `json:"status"`
	// `CfnBucket.MetricsProperty.EventThreshold`.
	EventThreshold interface{} `json:"eventThreshold"`
}

type CfnBucket_NoncurrentVersionTransitionProperty struct {
	// `CfnBucket.NoncurrentVersionTransitionProperty.StorageClass`.
	StorageClass *string `json:"storageClass"`
	// `CfnBucket.NoncurrentVersionTransitionProperty.TransitionInDays`.
	TransitionInDays *float64 `json:"transitionInDays"`
}

type CfnBucket_NotificationConfigurationProperty struct {
	// `CfnBucket.NotificationConfigurationProperty.LambdaConfigurations`.
	LambdaConfigurations interface{} `json:"lambdaConfigurations"`
	// `CfnBucket.NotificationConfigurationProperty.QueueConfigurations`.
	QueueConfigurations interface{} `json:"queueConfigurations"`
	// `CfnBucket.NotificationConfigurationProperty.TopicConfigurations`.
	TopicConfigurations interface{} `json:"topicConfigurations"`
}

type CfnBucket_NotificationFilterProperty struct {
	// `CfnBucket.NotificationFilterProperty.S3Key`.
	S3Key interface{} `json:"s3Key"`
}

type CfnBucket_ObjectLockConfigurationProperty struct {
	// `CfnBucket.ObjectLockConfigurationProperty.ObjectLockEnabled`.
	ObjectLockEnabled *string `json:"objectLockEnabled"`
	// `CfnBucket.ObjectLockConfigurationProperty.Rule`.
	Rule interface{} `json:"rule"`
}

type CfnBucket_ObjectLockRuleProperty struct {
	// `CfnBucket.ObjectLockRuleProperty.DefaultRetention`.
	DefaultRetention interface{} `json:"defaultRetention"`
}

type CfnBucket_OwnershipControlsProperty struct {
	// `CfnBucket.OwnershipControlsProperty.Rules`.
	Rules interface{} `json:"rules"`
}

type CfnBucket_OwnershipControlsRuleProperty struct {
	// `CfnBucket.OwnershipControlsRuleProperty.ObjectOwnership`.
	ObjectOwnership *string `json:"objectOwnership"`
}

type CfnBucket_PublicAccessBlockConfigurationProperty struct {
	// `CfnBucket.PublicAccessBlockConfigurationProperty.BlockPublicAcls`.
	BlockPublicAcls interface{} `json:"blockPublicAcls"`
	// `CfnBucket.PublicAccessBlockConfigurationProperty.BlockPublicPolicy`.
	BlockPublicPolicy interface{} `json:"blockPublicPolicy"`
	// `CfnBucket.PublicAccessBlockConfigurationProperty.IgnorePublicAcls`.
	IgnorePublicAcls interface{} `json:"ignorePublicAcls"`
	// `CfnBucket.PublicAccessBlockConfigurationProperty.RestrictPublicBuckets`.
	RestrictPublicBuckets interface{} `json:"restrictPublicBuckets"`
}

type CfnBucket_QueueConfigurationProperty struct {
	// `CfnBucket.QueueConfigurationProperty.Event`.
	Event *string `json:"event"`
	// `CfnBucket.QueueConfigurationProperty.Queue`.
	Queue *string `json:"queue"`
	// `CfnBucket.QueueConfigurationProperty.Filter`.
	Filter interface{} `json:"filter"`
}

type CfnBucket_RedirectAllRequestsToProperty struct {
	// `CfnBucket.RedirectAllRequestsToProperty.HostName`.
	HostName *string `json:"hostName"`
	// `CfnBucket.RedirectAllRequestsToProperty.Protocol`.
	Protocol *string `json:"protocol"`
}

type CfnBucket_RedirectRuleProperty struct {
	// `CfnBucket.RedirectRuleProperty.HostName`.
	HostName *string `json:"hostName"`
	// `CfnBucket.RedirectRuleProperty.HttpRedirectCode`.
	HttpRedirectCode *string `json:"httpRedirectCode"`
	// `CfnBucket.RedirectRuleProperty.Protocol`.
	Protocol *string `json:"protocol"`
	// `CfnBucket.RedirectRuleProperty.ReplaceKeyPrefixWith`.
	ReplaceKeyPrefixWith *string `json:"replaceKeyPrefixWith"`
	// `CfnBucket.RedirectRuleProperty.ReplaceKeyWith`.
	ReplaceKeyWith *string `json:"replaceKeyWith"`
}

type CfnBucket_ReplicaModificationsProperty struct {
	// `CfnBucket.ReplicaModificationsProperty.Status`.
	Status *string `json:"status"`
}

type CfnBucket_ReplicationConfigurationProperty struct {
	// `CfnBucket.ReplicationConfigurationProperty.Role`.
	Role *string `json:"role"`
	// `CfnBucket.ReplicationConfigurationProperty.Rules`.
	Rules interface{} `json:"rules"`
}

type CfnBucket_ReplicationDestinationProperty struct {
	// `CfnBucket.ReplicationDestinationProperty.Bucket`.
	Bucket *string `json:"bucket"`
	// `CfnBucket.ReplicationDestinationProperty.AccessControlTranslation`.
	AccessControlTranslation interface{} `json:"accessControlTranslation"`
	// `CfnBucket.ReplicationDestinationProperty.Account`.
	Account *string `json:"account"`
	// `CfnBucket.ReplicationDestinationProperty.EncryptionConfiguration`.
	EncryptionConfiguration interface{} `json:"encryptionConfiguration"`
	// `CfnBucket.ReplicationDestinationProperty.Metrics`.
	Metrics interface{} `json:"metrics"`
	// `CfnBucket.ReplicationDestinationProperty.ReplicationTime`.
	ReplicationTime interface{} `json:"replicationTime"`
	// `CfnBucket.ReplicationDestinationProperty.StorageClass`.
	StorageClass *string `json:"storageClass"`
}

type CfnBucket_ReplicationRuleAndOperatorProperty struct {
	// `CfnBucket.ReplicationRuleAndOperatorProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnBucket.ReplicationRuleAndOperatorProperty.TagFilters`.
	TagFilters interface{} `json:"tagFilters"`
}

type CfnBucket_ReplicationRuleFilterProperty struct {
	// `CfnBucket.ReplicationRuleFilterProperty.And`.
	And interface{} `json:"and"`
	// `CfnBucket.ReplicationRuleFilterProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnBucket.ReplicationRuleFilterProperty.TagFilter`.
	TagFilter interface{} `json:"tagFilter"`
}

type CfnBucket_ReplicationRuleProperty struct {
	// `CfnBucket.ReplicationRuleProperty.Destination`.
	Destination interface{} `json:"destination"`
	// `CfnBucket.ReplicationRuleProperty.Status`.
	Status *string `json:"status"`
	// `CfnBucket.ReplicationRuleProperty.DeleteMarkerReplication`.
	DeleteMarkerReplication interface{} `json:"deleteMarkerReplication"`
	// `CfnBucket.ReplicationRuleProperty.Filter`.
	Filter interface{} `json:"filter"`
	// `CfnBucket.ReplicationRuleProperty.Id`.
	Id *string `json:"id"`
	// `CfnBucket.ReplicationRuleProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnBucket.ReplicationRuleProperty.Priority`.
	Priority *float64 `json:"priority"`
	// `CfnBucket.ReplicationRuleProperty.SourceSelectionCriteria`.
	SourceSelectionCriteria interface{} `json:"sourceSelectionCriteria"`
}

type CfnBucket_ReplicationTimeProperty struct {
	// `CfnBucket.ReplicationTimeProperty.Status`.
	Status *string `json:"status"`
	// `CfnBucket.ReplicationTimeProperty.Time`.
	Time interface{} `json:"time"`
}

type CfnBucket_ReplicationTimeValueProperty struct {
	// `CfnBucket.ReplicationTimeValueProperty.Minutes`.
	Minutes *float64 `json:"minutes"`
}

type CfnBucket_RoutingRuleConditionProperty struct {
	// `CfnBucket.RoutingRuleConditionProperty.HttpErrorCodeReturnedEquals`.
	HttpErrorCodeReturnedEquals *string `json:"httpErrorCodeReturnedEquals"`
	// `CfnBucket.RoutingRuleConditionProperty.KeyPrefixEquals`.
	KeyPrefixEquals *string `json:"keyPrefixEquals"`
}

type CfnBucket_RoutingRuleProperty struct {
	// `CfnBucket.RoutingRuleProperty.RedirectRule`.
	RedirectRule interface{} `json:"redirectRule"`
	// `CfnBucket.RoutingRuleProperty.RoutingRuleCondition`.
	RoutingRuleCondition interface{} `json:"routingRuleCondition"`
}

type CfnBucket_RuleProperty struct {
	// `CfnBucket.RuleProperty.Status`.
	Status *string `json:"status"`
	// `CfnBucket.RuleProperty.AbortIncompleteMultipartUpload`.
	AbortIncompleteMultipartUpload interface{} `json:"abortIncompleteMultipartUpload"`
	// `CfnBucket.RuleProperty.ExpirationDate`.
	ExpirationDate interface{} `json:"expirationDate"`
	// `CfnBucket.RuleProperty.ExpirationInDays`.
	ExpirationInDays *float64 `json:"expirationInDays"`
	// `CfnBucket.RuleProperty.ExpiredObjectDeleteMarker`.
	ExpiredObjectDeleteMarker interface{} `json:"expiredObjectDeleteMarker"`
	// `CfnBucket.RuleProperty.Id`.
	Id *string `json:"id"`
	// `CfnBucket.RuleProperty.NoncurrentVersionExpirationInDays`.
	NoncurrentVersionExpirationInDays *float64 `json:"noncurrentVersionExpirationInDays"`
	// `CfnBucket.RuleProperty.NoncurrentVersionTransition`.
	NoncurrentVersionTransition interface{} `json:"noncurrentVersionTransition"`
	// `CfnBucket.RuleProperty.NoncurrentVersionTransitions`.
	NoncurrentVersionTransitions interface{} `json:"noncurrentVersionTransitions"`
	// `CfnBucket.RuleProperty.Prefix`.
	Prefix *string `json:"prefix"`
	// `CfnBucket.RuleProperty.TagFilters`.
	TagFilters interface{} `json:"tagFilters"`
	// `CfnBucket.RuleProperty.Transition`.
	Transition interface{} `json:"transition"`
	// `CfnBucket.RuleProperty.Transitions`.
	Transitions interface{} `json:"transitions"`
}

type CfnBucket_S3KeyFilterProperty struct {
	// `CfnBucket.S3KeyFilterProperty.Rules`.
	Rules interface{} `json:"rules"`
}

type CfnBucket_ServerSideEncryptionByDefaultProperty struct {
	// `CfnBucket.ServerSideEncryptionByDefaultProperty.SSEAlgorithm`.
	SseAlgorithm *string `json:"sseAlgorithm"`
	// `CfnBucket.ServerSideEncryptionByDefaultProperty.KMSMasterKeyID`.
	KmsMasterKeyId *string `json:"kmsMasterKeyId"`
}

type CfnBucket_ServerSideEncryptionRuleProperty struct {
	// `CfnBucket.ServerSideEncryptionRuleProperty.BucketKeyEnabled`.
	BucketKeyEnabled interface{} `json:"bucketKeyEnabled"`
	// `CfnBucket.ServerSideEncryptionRuleProperty.ServerSideEncryptionByDefault`.
	ServerSideEncryptionByDefault interface{} `json:"serverSideEncryptionByDefault"`
}

type CfnBucket_SourceSelectionCriteriaProperty struct {
	// `CfnBucket.SourceSelectionCriteriaProperty.ReplicaModifications`.
	ReplicaModifications interface{} `json:"replicaModifications"`
	// `CfnBucket.SourceSelectionCriteriaProperty.SseKmsEncryptedObjects`.
	SseKmsEncryptedObjects interface{} `json:"sseKmsEncryptedObjects"`
}

type CfnBucket_SseKmsEncryptedObjectsProperty struct {
	// `CfnBucket.SseKmsEncryptedObjectsProperty.Status`.
	Status *string `json:"status"`
}

type CfnBucket_StorageClassAnalysisProperty struct {
	// `CfnBucket.StorageClassAnalysisProperty.DataExport`.
	DataExport interface{} `json:"dataExport"`
}

type CfnBucket_TagFilterProperty struct {
	// `CfnBucket.TagFilterProperty.Key`.
	Key *string `json:"key"`
	// `CfnBucket.TagFilterProperty.Value`.
	Value *string `json:"value"`
}

type CfnBucket_TieringProperty struct {
	// `CfnBucket.TieringProperty.AccessTier`.
	AccessTier *string `json:"accessTier"`
	// `CfnBucket.TieringProperty.Days`.
	Days *float64 `json:"days"`
}

type CfnBucket_TopicConfigurationProperty struct {
	// `CfnBucket.TopicConfigurationProperty.Event`.
	Event *string `json:"event"`
	// `CfnBucket.TopicConfigurationProperty.Topic`.
	Topic *string `json:"topic"`
	// `CfnBucket.TopicConfigurationProperty.Filter`.
	Filter interface{} `json:"filter"`
}

type CfnBucket_TransitionProperty struct {
	// `CfnBucket.TransitionProperty.StorageClass`.
	StorageClass *string `json:"storageClass"`
	// `CfnBucket.TransitionProperty.TransitionDate`.
	TransitionDate interface{} `json:"transitionDate"`
	// `CfnBucket.TransitionProperty.TransitionInDays`.
	TransitionInDays *float64 `json:"transitionInDays"`
}

type CfnBucket_VersioningConfigurationProperty struct {
	// `CfnBucket.VersioningConfigurationProperty.Status`.
	Status *string `json:"status"`
}

type CfnBucket_WebsiteConfigurationProperty struct {
	// `CfnBucket.WebsiteConfigurationProperty.ErrorDocument`.
	ErrorDocument *string `json:"errorDocument"`
	// `CfnBucket.WebsiteConfigurationProperty.IndexDocument`.
	IndexDocument *string `json:"indexDocument"`
	// `CfnBucket.WebsiteConfigurationProperty.RedirectAllRequestsTo`.
	RedirectAllRequestsTo interface{} `json:"redirectAllRequestsTo"`
	// `CfnBucket.WebsiteConfigurationProperty.RoutingRules`.
	RoutingRules interface{} `json:"routingRules"`
}

// A CloudFormation `AWS::S3::BucketPolicy`.
type CfnBucketPolicy interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Bucket() *string
	SetBucket(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	PolicyDocument() interface{}
	SetPolicyDocument(val interface{})
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

// The jsii proxy struct for CfnBucketPolicy
type jsiiProxy_CfnBucketPolicy struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnBucketPolicy) Bucket() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucket",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) PolicyDocument() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policyDocument",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBucketPolicy) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::S3::BucketPolicy`.
func NewCfnBucketPolicy(scope awscdk.Construct, id *string, props *CfnBucketPolicyProps) CfnBucketPolicy {
	_init_.Initialize()

	j := jsiiProxy_CfnBucketPolicy{}

	_jsii_.Create(
		"monocdk.aws_s3.CfnBucketPolicy",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::S3::BucketPolicy`.
func NewCfnBucketPolicy_Override(c CfnBucketPolicy, scope awscdk.Construct, id *string, props *CfnBucketPolicyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.CfnBucketPolicy",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnBucketPolicy) SetBucket(val *string) {
	_jsii_.Set(
		j,
		"bucket",
		val,
	)
}

func (j *jsiiProxy_CfnBucketPolicy) SetPolicyDocument(val interface{}) {
	_jsii_.Set(
		j,
		"policyDocument",
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
func CfnBucketPolicy_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnBucketPolicy",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnBucketPolicy_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnBucketPolicy",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnBucketPolicy_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnBucketPolicy",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnBucketPolicy_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_s3.CfnBucketPolicy",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnBucketPolicy) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnBucketPolicy) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnBucketPolicy) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnBucketPolicy) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnBucketPolicy) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnBucketPolicy) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnBucketPolicy) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnBucketPolicy) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnBucketPolicy) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnBucketPolicy) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnBucketPolicy) OnPrepare() {
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
func (c *jsiiProxy_CfnBucketPolicy) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnBucketPolicy) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnBucketPolicy) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnBucketPolicy) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnBucketPolicy) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnBucketPolicy) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnBucketPolicy) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnBucketPolicy) ToString() *string {
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
func (c *jsiiProxy_CfnBucketPolicy) Validate() *[]*string {
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
func (c *jsiiProxy_CfnBucketPolicy) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::S3::BucketPolicy`.
type CfnBucketPolicyProps struct {
	// `AWS::S3::BucketPolicy.Bucket`.
	Bucket *string `json:"bucket"`
	// `AWS::S3::BucketPolicy.PolicyDocument`.
	PolicyDocument interface{} `json:"policyDocument"`
}

// Properties for defining a `AWS::S3::Bucket`.
type CfnBucketProps struct {
	// `AWS::S3::Bucket.AccelerateConfiguration`.
	AccelerateConfiguration interface{} `json:"accelerateConfiguration"`
	// `AWS::S3::Bucket.AccessControl`.
	AccessControl *string `json:"accessControl"`
	// `AWS::S3::Bucket.AnalyticsConfigurations`.
	AnalyticsConfigurations interface{} `json:"analyticsConfigurations"`
	// `AWS::S3::Bucket.BucketEncryption`.
	BucketEncryption interface{} `json:"bucketEncryption"`
	// `AWS::S3::Bucket.BucketName`.
	BucketName *string `json:"bucketName"`
	// `AWS::S3::Bucket.CorsConfiguration`.
	CorsConfiguration interface{} `json:"corsConfiguration"`
	// `AWS::S3::Bucket.IntelligentTieringConfigurations`.
	IntelligentTieringConfigurations interface{} `json:"intelligentTieringConfigurations"`
	// `AWS::S3::Bucket.InventoryConfigurations`.
	InventoryConfigurations interface{} `json:"inventoryConfigurations"`
	// `AWS::S3::Bucket.LifecycleConfiguration`.
	LifecycleConfiguration interface{} `json:"lifecycleConfiguration"`
	// `AWS::S3::Bucket.LoggingConfiguration`.
	LoggingConfiguration interface{} `json:"loggingConfiguration"`
	// `AWS::S3::Bucket.MetricsConfigurations`.
	MetricsConfigurations interface{} `json:"metricsConfigurations"`
	// `AWS::S3::Bucket.NotificationConfiguration`.
	NotificationConfiguration interface{} `json:"notificationConfiguration"`
	// `AWS::S3::Bucket.ObjectLockConfiguration`.
	ObjectLockConfiguration interface{} `json:"objectLockConfiguration"`
	// `AWS::S3::Bucket.ObjectLockEnabled`.
	ObjectLockEnabled interface{} `json:"objectLockEnabled"`
	// `AWS::S3::Bucket.OwnershipControls`.
	OwnershipControls interface{} `json:"ownershipControls"`
	// `AWS::S3::Bucket.PublicAccessBlockConfiguration`.
	PublicAccessBlockConfiguration interface{} `json:"publicAccessBlockConfiguration"`
	// `AWS::S3::Bucket.ReplicationConfiguration`.
	ReplicationConfiguration interface{} `json:"replicationConfiguration"`
	// `AWS::S3::Bucket.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::S3::Bucket.VersioningConfiguration`.
	VersioningConfiguration interface{} `json:"versioningConfiguration"`
	// `AWS::S3::Bucket.WebsiteConfiguration`.
	WebsiteConfiguration interface{} `json:"websiteConfiguration"`
}

// A CloudFormation `AWS::S3::StorageLens`.
type CfnStorageLens interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrStorageLensConfigurationStorageLensArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	StorageLensConfiguration() interface{}
	SetStorageLensConfiguration(val interface{})
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

// The jsii proxy struct for CfnStorageLens
type jsiiProxy_CfnStorageLens struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnStorageLens) AttrStorageLensConfigurationStorageLensArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrStorageLensConfigurationStorageLensArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) StorageLensConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"storageLensConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStorageLens) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::S3::StorageLens`.
func NewCfnStorageLens(scope awscdk.Construct, id *string, props *CfnStorageLensProps) CfnStorageLens {
	_init_.Initialize()

	j := jsiiProxy_CfnStorageLens{}

	_jsii_.Create(
		"monocdk.aws_s3.CfnStorageLens",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::S3::StorageLens`.
func NewCfnStorageLens_Override(c CfnStorageLens, scope awscdk.Construct, id *string, props *CfnStorageLensProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.CfnStorageLens",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStorageLens) SetStorageLensConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"storageLensConfiguration",
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
func CfnStorageLens_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnStorageLens",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnStorageLens_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnStorageLens",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnStorageLens_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.CfnStorageLens",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStorageLens_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_s3.CfnStorageLens",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnStorageLens) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnStorageLens) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnStorageLens) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnStorageLens) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnStorageLens) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnStorageLens) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnStorageLens) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnStorageLens) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnStorageLens) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnStorageLens) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnStorageLens) OnPrepare() {
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
func (c *jsiiProxy_CfnStorageLens) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStorageLens) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnStorageLens) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnStorageLens) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnStorageLens) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnStorageLens) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnStorageLens) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStorageLens) ToString() *string {
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
func (c *jsiiProxy_CfnStorageLens) Validate() *[]*string {
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
func (c *jsiiProxy_CfnStorageLens) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnStorageLens_AccountLevelProperty struct {
	// `CfnStorageLens.AccountLevelProperty.BucketLevel`.
	BucketLevel interface{} `json:"bucketLevel"`
	// `CfnStorageLens.AccountLevelProperty.ActivityMetrics`.
	ActivityMetrics interface{} `json:"activityMetrics"`
}

type CfnStorageLens_ActivityMetricsProperty struct {
	// `CfnStorageLens.ActivityMetricsProperty.IsEnabled`.
	IsEnabled interface{} `json:"isEnabled"`
}

type CfnStorageLens_AwsOrgProperty struct {
	// `CfnStorageLens.AwsOrgProperty.Arn`.
	Arn *string `json:"arn"`
}

type CfnStorageLens_BucketLevelProperty struct {
	// `CfnStorageLens.BucketLevelProperty.ActivityMetrics`.
	ActivityMetrics interface{} `json:"activityMetrics"`
	// `CfnStorageLens.BucketLevelProperty.PrefixLevel`.
	PrefixLevel interface{} `json:"prefixLevel"`
}

type CfnStorageLens_BucketsAndRegionsProperty struct {
	// `CfnStorageLens.BucketsAndRegionsProperty.Buckets`.
	Buckets *[]*string `json:"buckets"`
	// `CfnStorageLens.BucketsAndRegionsProperty.Regions`.
	Regions *[]*string `json:"regions"`
}

type CfnStorageLens_DataExportProperty struct {
	// `CfnStorageLens.DataExportProperty.S3BucketDestination`.
	S3BucketDestination interface{} `json:"s3BucketDestination"`
}

type CfnStorageLens_EncryptionProperty struct {
}

type CfnStorageLens_PrefixLevelProperty struct {
	// `CfnStorageLens.PrefixLevelProperty.StorageMetrics`.
	StorageMetrics interface{} `json:"storageMetrics"`
}

type CfnStorageLens_PrefixLevelStorageMetricsProperty struct {
	// `CfnStorageLens.PrefixLevelStorageMetricsProperty.IsEnabled`.
	IsEnabled interface{} `json:"isEnabled"`
	// `CfnStorageLens.PrefixLevelStorageMetricsProperty.SelectionCriteria`.
	SelectionCriteria interface{} `json:"selectionCriteria"`
}

type CfnStorageLens_S3BucketDestinationProperty struct {
	// `CfnStorageLens.S3BucketDestinationProperty.AccountId`.
	AccountId *string `json:"accountId"`
	// `CfnStorageLens.S3BucketDestinationProperty.Arn`.
	Arn *string `json:"arn"`
	// `CfnStorageLens.S3BucketDestinationProperty.Format`.
	Format *string `json:"format"`
	// `CfnStorageLens.S3BucketDestinationProperty.OutputSchemaVersion`.
	OutputSchemaVersion *string `json:"outputSchemaVersion"`
	// `CfnStorageLens.S3BucketDestinationProperty.Encryption`.
	Encryption interface{} `json:"encryption"`
	// `CfnStorageLens.S3BucketDestinationProperty.Prefix`.
	Prefix *string `json:"prefix"`
}

type CfnStorageLens_SelectionCriteriaProperty struct {
	// `CfnStorageLens.SelectionCriteriaProperty.Delimiter`.
	Delimiter *string `json:"delimiter"`
	// `CfnStorageLens.SelectionCriteriaProperty.MaxDepth`.
	MaxDepth *float64 `json:"maxDepth"`
	// `CfnStorageLens.SelectionCriteriaProperty.MinStorageBytesPercentage`.
	MinStorageBytesPercentage *float64 `json:"minStorageBytesPercentage"`
}

type CfnStorageLens_StorageLensConfigurationProperty struct {
	// `CfnStorageLens.StorageLensConfigurationProperty.AccountLevel`.
	AccountLevel interface{} `json:"accountLevel"`
	// `CfnStorageLens.StorageLensConfigurationProperty.Id`.
	Id *string `json:"id"`
	// `CfnStorageLens.StorageLensConfigurationProperty.IsEnabled`.
	IsEnabled interface{} `json:"isEnabled"`
	// `CfnStorageLens.StorageLensConfigurationProperty.AwsOrg`.
	AwsOrg interface{} `json:"awsOrg"`
	// `CfnStorageLens.StorageLensConfigurationProperty.DataExport`.
	DataExport interface{} `json:"dataExport"`
	// `CfnStorageLens.StorageLensConfigurationProperty.Exclude`.
	Exclude interface{} `json:"exclude"`
	// `CfnStorageLens.StorageLensConfigurationProperty.Include`.
	Include interface{} `json:"include"`
	// `CfnStorageLens.StorageLensConfigurationProperty.StorageLensArn`.
	StorageLensArn *string `json:"storageLensArn"`
}

// Properties for defining a `AWS::S3::StorageLens`.
type CfnStorageLensProps struct {
	// `AWS::S3::StorageLens.StorageLensConfiguration`.
	StorageLensConfiguration interface{} `json:"storageLensConfiguration"`
	// `AWS::S3::StorageLens.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// Specifies a cross-origin access rule for an Amazon S3 bucket.
// Experimental.
type CorsRule struct {
	// An HTTP method that you allow the origin to execute.
	// Experimental.
	AllowedMethods *[]HttpMethods `json:"allowedMethods"`
	// One or more origins you want customers to be able to access the bucket from.
	// Experimental.
	AllowedOrigins *[]*string `json:"allowedOrigins"`
	// Headers that are specified in the Access-Control-Request-Headers header.
	// Experimental.
	AllowedHeaders *[]*string `json:"allowedHeaders"`
	// One or more headers in the response that you want customers to be able to access from their applications.
	// Experimental.
	ExposedHeaders *[]*string `json:"exposedHeaders"`
	// A unique identifier for this rule.
	// Experimental.
	Id *string `json:"id"`
	// The time in seconds that your browser is to cache the preflight response for the specified resource.
	// Experimental.
	MaxAge *float64 `json:"maxAge"`
}

// Notification event types.
// Experimental.
type EventType string

const (
	EventType_OBJECT_CREATED EventType = "OBJECT_CREATED"
	EventType_OBJECT_CREATED_PUT EventType = "OBJECT_CREATED_PUT"
	EventType_OBJECT_CREATED_POST EventType = "OBJECT_CREATED_POST"
	EventType_OBJECT_CREATED_COPY EventType = "OBJECT_CREATED_COPY"
	EventType_OBJECT_CREATED_COMPLETE_MULTIPART_UPLOAD EventType = "OBJECT_CREATED_COMPLETE_MULTIPART_UPLOAD"
	EventType_OBJECT_REMOVED EventType = "OBJECT_REMOVED"
	EventType_OBJECT_REMOVED_DELETE EventType = "OBJECT_REMOVED_DELETE"
	EventType_OBJECT_REMOVED_DELETE_MARKER_CREATED EventType = "OBJECT_REMOVED_DELETE_MARKER_CREATED"
	EventType_OBJECT_RESTORE_POST EventType = "OBJECT_RESTORE_POST"
	EventType_OBJECT_RESTORE_COMPLETED EventType = "OBJECT_RESTORE_COMPLETED"
	EventType_REDUCED_REDUNDANCY_LOST_OBJECT EventType = "REDUCED_REDUNDANCY_LOST_OBJECT"
	EventType_REPLICATION_OPERATION_FAILED_REPLICATION EventType = "REPLICATION_OPERATION_FAILED_REPLICATION"
	EventType_REPLICATION_OPERATION_MISSED_THRESHOLD EventType = "REPLICATION_OPERATION_MISSED_THRESHOLD"
	EventType_REPLICATION_OPERATION_REPLICATED_AFTER_THRESHOLD EventType = "REPLICATION_OPERATION_REPLICATED_AFTER_THRESHOLD"
	EventType_REPLICATION_OPERATION_NOT_TRACKED EventType = "REPLICATION_OPERATION_NOT_TRACKED"
)

// All http request methods.
// Experimental.
type HttpMethods string

const (
	HttpMethods_GET HttpMethods = "GET"
	HttpMethods_PUT HttpMethods = "PUT"
	HttpMethods_HEAD HttpMethods = "HEAD"
	HttpMethods_POST HttpMethods = "POST"
	HttpMethods_DELETE HttpMethods = "DELETE"
)

// Experimental.
type IBucket interface {
	awscdk.IResource
	// Adds a bucket notification event destination.
	//
	// TODO: EXAMPLE
	//
	// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html
	//
	// Experimental.
	AddEventNotification(event EventType, dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	// Subscribes a destination to receive notifications when an object is created in the bucket.
	//
	// This is identical to calling
	// `onEvent(EventType.ObjectCreated)`.
	// Experimental.
	AddObjectCreatedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	// Subscribes a destination to receive notifications when an object is removed from the bucket.
	//
	// This is identical to calling
	// `onEvent(EventType.ObjectRemoved)`.
	// Experimental.
	AddObjectRemovedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter)
	// Adds a statement to the resource policy for a principal (i.e. account/role/service) to perform actions on this bucket and/or it's contents. Use `bucketArn` and `arnForObjects(keys)` to obtain ARNs for this bucket or objects.
	// Experimental.
	AddToResourcePolicy(permission awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult
	// Returns an ARN that represents all objects within the bucket that match the key pattern specified.
	//
	// To represent all keys, specify ``"*"``.
	// Experimental.
	ArnForObjects(keyPattern *string) *string
	// Grants s3:DeleteObject* permission to an IAM principal for objects in this bucket.
	// Experimental.
	GrantDelete(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	// Allows unrestricted access to objects from this bucket.
	//
	// IMPORTANT: This permission allows anyone to perform actions on S3 objects
	// in this bucket, which is useful for when you configure your bucket as a
	// website and want everyone to be able to read objects in the bucket without
	// needing to authenticate.
	//
	// Without arguments, this method will grant read ("s3:GetObject") access to
	// all objects ("*") in the bucket.
	//
	// The method returns the `iam.Grant` object, which can then be modified
	// as needed. For example, you can add a condition that will restrict access only
	// to an IPv4 range like this:
	//
	//      const grant = bucket.grantPublicAccess();
	//      grant.resourceStatement!.addCondition(‘IpAddress’, { “aws:SourceIp”: “54.240.143.0/24” });
	//
	// Returns: The `iam.PolicyStatement` object, which can be used to apply e.g. conditions.
	// Experimental.
	GrantPublicAccess(keyPrefix *string, allowedActions ...*string) awsiam.Grant
	// Grants s3:PutObject* and s3:Abort* permissions for this bucket to an IAM principal.
	//
	// If encryption is used, permission to use the key to encrypt the contents
	// of written files will also be granted to the same principal.
	// Experimental.
	GrantPut(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	// Grant the given IAM identity permissions to modify the ACLs of objects in the given Bucket.
	//
	// If your application has the '@aws-cdk/aws-s3:grantWriteWithoutAcl' feature flag set,
	// calling {@link grantWrite} or {@link grantReadWrite} no longer grants permissions to modify the ACLs of the objects;
	// in this case, if you need to modify object ACLs, call this method explicitly.
	// Experimental.
	GrantPutAcl(identity awsiam.IGrantable, objectsKeyPattern *string) awsiam.Grant
	// Grant read permissions for this bucket and it's contents to an IAM principal (Role/Group/User).
	//
	// If encryption is used, permission to use the key to decrypt the contents
	// of the bucket will also be granted to the same principal.
	// Experimental.
	GrantRead(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	// Grants read/write permissions for this bucket and it's contents to an IAM principal (Role/Group/User).
	//
	// If an encryption key is used, permission to use the key for
	// encrypt/decrypt will also be granted.
	//
	// Before CDK version 1.85.0, this method granted the `s3:PutObject*` permission that included `s3:PutObjectAcl`,
	// which could be used to grant read/write object access to IAM principals in other accounts.
	// If you want to get rid of that behavior, update your CDK version to 1.85.0 or later,
	// and make sure the `@aws-cdk/aws-s3:grantWriteWithoutAcl` feature flag is set to `true`
	// in the `context` key of your cdk.json file.
	// If you've already updated, but still need the principal to have permissions to modify the ACLs,
	// use the {@link grantPutAcl} method.
	// Experimental.
	GrantReadWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	// Grant write permissions to this bucket to an IAM principal.
	//
	// If encryption is used, permission to use the key to encrypt the contents
	// of written files will also be granted to the same principal.
	//
	// Before CDK version 1.85.0, this method granted the `s3:PutObject*` permission that included `s3:PutObjectAcl`,
	// which could be used to grant read/write object access to IAM principals in other accounts.
	// If you want to get rid of that behavior, update your CDK version to 1.85.0 or later,
	// and make sure the `@aws-cdk/aws-s3:grantWriteWithoutAcl` feature flag is set to `true`
	// in the `context` key of your cdk.json file.
	// If you've already updated, but still need the principal to have permissions to modify the ACLs,
	// use the {@link grantPutAcl} method.
	// Experimental.
	GrantWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant
	// Defines a CloudWatch event that triggers when something happens to this bucket.
	//
	// Requires that there exists at least one CloudTrail Trail in your account
	// that captures the event. This method will not create the Trail.
	// Experimental.
	OnCloudTrailEvent(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	// Defines an AWS CloudWatch event that triggers when an object is uploaded to the specified paths (keys) in this bucket using the PutObject API call.
	//
	// Note that some tools like `aws s3 cp` will automatically use either
	// PutObject or the multipart upload API depending on the file size,
	// so using `onCloudTrailWriteObject` may be preferable.
	//
	// Requires that there exists at least one CloudTrail Trail in your account
	// that captures the event. This method will not create the Trail.
	// Experimental.
	OnCloudTrailPutObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	// Defines an AWS CloudWatch event that triggers when an object at the specified paths (keys) in this bucket are written to.
	//
	// This includes
	// the events PutObject, CopyObject, and CompleteMultipartUpload.
	//
	// Note that some tools like `aws s3 cp` will automatically use either
	// PutObject or the multipart upload API depending on the file size,
	// so using this method may be preferable to `onCloudTrailPutObject`.
	//
	// Requires that there exists at least one CloudTrail Trail in your account
	// that captures the event. This method will not create the Trail.
	// Experimental.
	OnCloudTrailWriteObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule
	// The S3 URL of an S3 object.
	//
	// For example:
	//
	// Returns: an ObjectS3Url token
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	S3UrlForObject(key *string) *string
	// The https URL of an S3 object.
	//
	// For example:
	//
	// Returns: an ObjectS3Url token
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	UrlForObject(key *string) *string
	// The virtual hosted-style URL of an S3 object.
	//
	// Specify `regional: false` at
	// the options for non-regional URL. For example:
	//
	// Returns: an ObjectS3Url token
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	VirtualHostedUrlForObject(key *string, options *VirtualHostedStyleUrlOptions) *string
	// The ARN of the bucket.
	// Experimental.
	BucketArn() *string
	// The IPv4 DNS name of the specified bucket.
	// Experimental.
	BucketDomainName() *string
	// The IPv6 DNS name of the specified bucket.
	// Experimental.
	BucketDualStackDomainName() *string
	// The name of the bucket.
	// Experimental.
	BucketName() *string
	// The regional domain name of the specified bucket.
	// Experimental.
	BucketRegionalDomainName() *string
	// The Domain name of the static website.
	// Experimental.
	BucketWebsiteDomainName() *string
	// The URL of the static website.
	// Experimental.
	BucketWebsiteUrl() *string
	// Optional KMS encryption key associated with this bucket.
	// Experimental.
	EncryptionKey() awskms.IKey
	// If this bucket has been configured for static website hosting.
	// Experimental.
	IsWebsite() *bool
	// The resource policy associated with this bucket.
	//
	// If `autoCreatePolicy` is true, a `BucketPolicy` will be created upon the
	// first call to addToResourcePolicy(s).
	// Experimental.
	Policy() BucketPolicy
	// The resource policy associated with this bucket.
	//
	// If `autoCreatePolicy` is true, a `BucketPolicy` will be created upon the
	// first call to addToResourcePolicy(s).
	// Experimental.
	SetPolicy(p BucketPolicy)
}

// The jsii proxy for IBucket
type jsiiProxy_IBucket struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IBucket) AddEventNotification(event EventType, dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{event, dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		i,
		"addEventNotification",
		args,
	)
}

func (i *jsiiProxy_IBucket) AddObjectCreatedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		i,
		"addObjectCreatedNotification",
		args,
	)
}

func (i *jsiiProxy_IBucket) AddObjectRemovedNotification(dest IBucketNotificationDestination, filters ...*NotificationKeyFilter) {
	args := []interface{}{dest}
	for _, a := range filters {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		i,
		"addObjectRemovedNotification",
		args,
	)
}

func (i *jsiiProxy_IBucket) AddToResourcePolicy(permission awsiam.PolicyStatement) *awsiam.AddToResourcePolicyResult {
	var returns *awsiam.AddToResourcePolicyResult

	_jsii_.Invoke(
		i,
		"addToResourcePolicy",
		[]interface{}{permission},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) ArnForObjects(keyPattern *string) *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"arnForObjects",
		[]interface{}{keyPattern},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) GrantDelete(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantDelete",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) GrantPublicAccess(keyPrefix *string, allowedActions ...*string) awsiam.Grant {
	args := []interface{}{keyPrefix}
	for _, a := range allowedActions {
		args = append(args, a)
	}

	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPublicAccess",
		args,
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) GrantPut(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPut",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) GrantPutAcl(identity awsiam.IGrantable, objectsKeyPattern *string) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantPutAcl",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) GrantRead(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantRead",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) GrantReadWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantReadWrite",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) GrantWrite(identity awsiam.IGrantable, objectsKeyPattern interface{}) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		i,
		"grantWrite",
		[]interface{}{identity, objectsKeyPattern},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) OnCloudTrailEvent(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onCloudTrailEvent",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) OnCloudTrailPutObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onCloudTrailPutObject",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) OnCloudTrailWriteObject(id *string, options *OnCloudTrailBucketEventOptions) awsevents.Rule {
	var returns awsevents.Rule

	_jsii_.Invoke(
		i,
		"onCloudTrailWriteObject",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) S3UrlForObject(key *string) *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"s3UrlForObject",
		[]interface{}{key},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) UrlForObject(key *string) *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"urlForObject",
		[]interface{}{key},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IBucket) VirtualHostedUrlForObject(key *string, options *VirtualHostedStyleUrlOptions) *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"virtualHostedUrlForObject",
		[]interface{}{key, options},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IBucket) BucketArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) BucketDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) BucketDualStackDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketDualStackDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) BucketName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) BucketRegionalDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketRegionalDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) BucketWebsiteDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketWebsiteDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) BucketWebsiteUrl() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketWebsiteUrl",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) EncryptionKey() awskms.IKey {
	var returns awskms.IKey
	_jsii_.Get(
		j,
		"encryptionKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) IsWebsite() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isWebsite",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) Policy() BucketPolicy {
	var returns BucketPolicy
	_jsii_.Get(
		j,
		"policy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IBucket) SetPolicy(val BucketPolicy) {
	_jsii_.Set(
		j,
		"policy",
		val,
	)
}

// Implemented by constructs that can be used as bucket notification destinations.
// Experimental.
type IBucketNotificationDestination interface {
	// Registers this resource to receive notifications for the specified bucket.
	//
	// This method will only be called once for each destination/bucket
	// pair and the result will be cached, so there is no need to implement
	// idempotency in each destination.
	// Experimental.
	Bind(scope awscdk.Construct, bucket IBucket) *BucketNotificationDestinationConfig
}

// The jsii proxy for IBucketNotificationDestination
type jsiiProxy_IBucketNotificationDestination struct {
	_ byte // padding
}

func (i *jsiiProxy_IBucketNotificationDestination) Bind(scope awscdk.Construct, bucket IBucket) *BucketNotificationDestinationConfig {
	var returns *BucketNotificationDestinationConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{scope, bucket},
		&returns,
	)

	return returns
}

// Specifies the inventory configuration of an S3 Bucket.
// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/storage-inventory.html
//
// Experimental.
type Inventory struct {
	// The destination of the inventory.
	// Experimental.
	Destination *InventoryDestination `json:"destination"`
	// Whether the inventory is enabled or disabled.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// The format of the inventory.
	// Experimental.
	Format InventoryFormat `json:"format"`
	// Frequency at which the inventory should be generated.
	// Experimental.
	Frequency InventoryFrequency `json:"frequency"`
	// If the inventory should contain all the object versions or only the current one.
	// Experimental.
	IncludeObjectVersions InventoryObjectVersion `json:"includeObjectVersions"`
	// The inventory configuration ID.
	// Experimental.
	InventoryId *string `json:"inventoryId"`
	// The inventory will only include objects that meet the prefix filter criteria.
	// Experimental.
	ObjectsPrefix *string `json:"objectsPrefix"`
	// A list of optional fields to be included in the inventory result.
	// Experimental.
	OptionalFields *[]*string `json:"optionalFields"`
}

// The destination of the inventory.
// Experimental.
type InventoryDestination struct {
	// Bucket where all inventories will be saved in.
	// Experimental.
	Bucket IBucket `json:"bucket"`
	// The account ID that owns the destination S3 bucket.
	//
	// If no account ID is provided, the owner is not validated before exporting data.
	// It's recommended to set an account ID to prevent problems if the destination bucket ownership changes.
	// Experimental.
	BucketOwner *string `json:"bucketOwner"`
	// The prefix to be used when saving the inventory.
	// Experimental.
	Prefix *string `json:"prefix"`
}

// All supported inventory list formats.
// Experimental.
type InventoryFormat string

const (
	InventoryFormat_CSV InventoryFormat = "CSV"
	InventoryFormat_PARQUET InventoryFormat = "PARQUET"
	InventoryFormat_ORC InventoryFormat = "ORC"
)

// All supported inventory frequencies.
// Experimental.
type InventoryFrequency string

const (
	InventoryFrequency_DAILY InventoryFrequency = "DAILY"
	InventoryFrequency_WEEKLY InventoryFrequency = "WEEKLY"
)

// Inventory version support.
// Experimental.
type InventoryObjectVersion string

const (
	InventoryObjectVersion_ALL InventoryObjectVersion = "ALL"
	InventoryObjectVersion_CURRENT InventoryObjectVersion = "CURRENT"
)

// Declaration of a Life cycle rule.
// Experimental.
type LifecycleRule struct {
	// Specifies a lifecycle rule that aborts incomplete multipart uploads to an Amazon S3 bucket.
	//
	// The AbortIncompleteMultipartUpload property type creates a lifecycle
	// rule that aborts incomplete multipart uploads to an Amazon S3 bucket.
	// When Amazon S3 aborts a multipart upload, it deletes all parts
	// associated with the multipart upload.
	// Experimental.
	AbortIncompleteMultipartUploadAfter awscdk.Duration `json:"abortIncompleteMultipartUploadAfter"`
	// Whether this rule is enabled.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// Indicates the number of days after creation when objects are deleted from Amazon S3 and Amazon Glacier.
	//
	// If you specify an expiration and transition time, you must use the same
	// time unit for both properties (either in days or by date). The
	// expiration time must also be later than the transition time.
	// Experimental.
	Expiration awscdk.Duration `json:"expiration"`
	// Indicates when objects are deleted from Amazon S3 and Amazon Glacier.
	//
	// The date value must be in ISO 8601 format. The time is always midnight UTC.
	//
	// If you specify an expiration and transition time, you must use the same
	// time unit for both properties (either in days or by date). The
	// expiration time must also be later than the transition time.
	// Experimental.
	ExpirationDate *time.Time `json:"expirationDate"`
	// Indicates whether Amazon S3 will remove a delete marker with no noncurrent versions.
	//
	// If set to true, the delete marker will be expired.
	// Experimental.
	ExpiredObjectDeleteMarker *bool `json:"expiredObjectDeleteMarker"`
	// A unique identifier for this rule.
	//
	// The value cannot be more than 255 characters.
	// Experimental.
	Id *string `json:"id"`
	// Time between when a new version of the object is uploaded to the bucket and when old versions of the object expire.
	//
	// For buckets with versioning enabled (or suspended), specifies the time,
	// in days, between when a new version of the object is uploaded to the
	// bucket and when old versions of the object expire. When object versions
	// expire, Amazon S3 permanently deletes them. If you specify a transition
	// and expiration time, the expiration time must be later than the
	// transition time.
	// Experimental.
	NoncurrentVersionExpiration awscdk.Duration `json:"noncurrentVersionExpiration"`
	// One or more transition rules that specify when non-current objects transition to a specified storage class.
	//
	// Only for for buckets with versioning enabled (or suspended).
	//
	// If you specify a transition and expiration time, the expiration time
	// must be later than the transition time.
	// Experimental.
	NoncurrentVersionTransitions *[]*NoncurrentVersionTransition `json:"noncurrentVersionTransitions"`
	// Object key prefix that identifies one or more objects to which this rule applies.
	// Experimental.
	Prefix *string `json:"prefix"`
	// The TagFilter property type specifies tags to use to identify a subset of objects for an Amazon S3 bucket.
	// Experimental.
	TagFilters *map[string]interface{} `json:"tagFilters"`
	// One or more transition rules that specify when an object transitions to a specified storage class.
	//
	// If you specify an expiration and transition time, you must use the same
	// time unit for both properties (either in days or by date). The
	// expiration time must also be later than the transition time.
	// Experimental.
	Transitions *[]*Transition `json:"transitions"`
}

// An interface that represents the location of a specific object in an S3 Bucket.
// Experimental.
type Location struct {
	// The name of the S3 Bucket the object is in.
	// Experimental.
	BucketName *string `json:"bucketName"`
	// The path inside the Bucket where the object is located at.
	// Experimental.
	ObjectKey *string `json:"objectKey"`
	// The S3 object version.
	// Experimental.
	ObjectVersion *string `json:"objectVersion"`
}

// Describes when noncurrent versions transition to a specified storage class.
// Experimental.
type NoncurrentVersionTransition struct {
	// The storage class to which you want the object to transition.
	// Experimental.
	StorageClass StorageClass `json:"storageClass"`
	// Indicates the number of days after creation when objects are transitioned to the specified storage class.
	// Experimental.
	TransitionAfter awscdk.Duration `json:"transitionAfter"`
}

// Experimental.
type NotificationKeyFilter struct {
	// S3 keys must have the specified prefix.
	// Experimental.
	Prefix *string `json:"prefix"`
	// S3 keys must have the specified suffix.
	// Experimental.
	Suffix *string `json:"suffix"`
}

// The ObjectOwnership of the bucket.
// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/about-object-ownership.html
//
// Experimental.
type ObjectOwnership string

const (
	ObjectOwnership_BUCKET_OWNER_PREFERRED ObjectOwnership = "BUCKET_OWNER_PREFERRED"
	ObjectOwnership_OBJECT_WRITER ObjectOwnership = "OBJECT_WRITER"
)

// Options for the onCloudTrailPutObject method.
// Experimental.
type OnCloudTrailBucketEventOptions struct {
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
	// Only watch changes to these object paths.
	// Experimental.
	Paths *[]*string `json:"paths"`
}

// All http request methods.
// Experimental.
type RedirectProtocol string

const (
	RedirectProtocol_HTTP RedirectProtocol = "HTTP"
	RedirectProtocol_HTTPS RedirectProtocol = "HTTPS"
)

// Specifies a redirect behavior of all requests to a website endpoint of a bucket.
// Experimental.
type RedirectTarget struct {
	// Name of the host where requests are redirected.
	// Experimental.
	HostName *string `json:"hostName"`
	// Protocol to use when redirecting requests.
	// Experimental.
	Protocol RedirectProtocol `json:"protocol"`
}

// Experimental.
type ReplaceKey interface {
	PrefixWithKey() *string
	WithKey() *string
}

// The jsii proxy struct for ReplaceKey
type jsiiProxy_ReplaceKey struct {
	_ byte // padding
}

func (j *jsiiProxy_ReplaceKey) PrefixWithKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"prefixWithKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ReplaceKey) WithKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"withKey",
		&returns,
	)
	return returns
}


// The object key prefix to use in the redirect request.
// Experimental.
func ReplaceKey_PrefixWith(keyReplacement *string) ReplaceKey {
	_init_.Initialize()

	var returns ReplaceKey

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.ReplaceKey",
		"prefixWith",
		[]interface{}{keyReplacement},
		&returns,
	)

	return returns
}

// The specific object key to use in the redirect request.
// Experimental.
func ReplaceKey_With(keyReplacement *string) ReplaceKey {
	_init_.Initialize()

	var returns ReplaceKey

	_jsii_.StaticInvoke(
		"monocdk.aws_s3.ReplaceKey",
		"with",
		[]interface{}{keyReplacement},
		&returns,
	)

	return returns
}

// Rule that define when a redirect is applied and the redirect behavior.
// See: https://docs.aws.amazon.com/AmazonS3/latest/dev/how-to-page-redirect.html
//
// Experimental.
type RoutingRule struct {
	// Specifies a condition that must be met for the specified redirect to apply.
	// Experimental.
	Condition *RoutingRuleCondition `json:"condition"`
	// The host name to use in the redirect request.
	// Experimental.
	HostName *string `json:"hostName"`
	// The HTTP redirect code to use on the response.
	// Experimental.
	HttpRedirectCode *string `json:"httpRedirectCode"`
	// Protocol to use when redirecting requests.
	// Experimental.
	Protocol RedirectProtocol `json:"protocol"`
	// Specifies the object key prefix to use in the redirect request.
	// Experimental.
	ReplaceKey ReplaceKey `json:"replaceKey"`
}

// Experimental.
type RoutingRuleCondition struct {
	// The HTTP error code when the redirect is applied.
	//
	// In the event of an error, if the error code equals this value, then the specified redirect is applied.
	//
	// If both condition properties are specified, both must be true for the redirect to be applied.
	// Experimental.
	HttpErrorCodeReturnedEquals *string `json:"httpErrorCodeReturnedEquals"`
	// The object key name prefix when the redirect is applied.
	//
	// If both condition properties are specified, both must be true for the redirect to be applied.
	// Experimental.
	KeyPrefixEquals *string `json:"keyPrefixEquals"`
}

// Storage class to move an object to.
// Experimental.
type StorageClass interface {
	Value() *string
	ToString() *string
}

// The jsii proxy struct for StorageClass
type jsiiProxy_StorageClass struct {
	_ byte // padding
}

func (j *jsiiProxy_StorageClass) Value() *string {
	var returns *string
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


// Experimental.
func NewStorageClass(value *string) StorageClass {
	_init_.Initialize()

	j := jsiiProxy_StorageClass{}

	_jsii_.Create(
		"monocdk.aws_s3.StorageClass",
		[]interface{}{value},
		&j,
	)

	return &j
}

// Experimental.
func NewStorageClass_Override(s StorageClass, value *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_s3.StorageClass",
		[]interface{}{value},
		s,
	)
}

func StorageClass_DEEP_ARCHIVE() StorageClass {
	_init_.Initialize()
	var returns StorageClass
	_jsii_.StaticGet(
		"monocdk.aws_s3.StorageClass",
		"DEEP_ARCHIVE",
		&returns,
	)
	return returns
}

func StorageClass_GLACIER() StorageClass {
	_init_.Initialize()
	var returns StorageClass
	_jsii_.StaticGet(
		"monocdk.aws_s3.StorageClass",
		"GLACIER",
		&returns,
	)
	return returns
}

func StorageClass_INFREQUENT_ACCESS() StorageClass {
	_init_.Initialize()
	var returns StorageClass
	_jsii_.StaticGet(
		"monocdk.aws_s3.StorageClass",
		"INFREQUENT_ACCESS",
		&returns,
	)
	return returns
}

func StorageClass_INTELLIGENT_TIERING() StorageClass {
	_init_.Initialize()
	var returns StorageClass
	_jsii_.StaticGet(
		"monocdk.aws_s3.StorageClass",
		"INTELLIGENT_TIERING",
		&returns,
	)
	return returns
}

func StorageClass_ONE_ZONE_INFREQUENT_ACCESS() StorageClass {
	_init_.Initialize()
	var returns StorageClass
	_jsii_.StaticGet(
		"monocdk.aws_s3.StorageClass",
		"ONE_ZONE_INFREQUENT_ACCESS",
		&returns,
	)
	return returns
}

// Experimental.
func (s *jsiiProxy_StorageClass) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Describes when an object transitions to a specified storage class.
// Experimental.
type Transition struct {
	// The storage class to which you want the object to transition.
	// Experimental.
	StorageClass StorageClass `json:"storageClass"`
	// Indicates the number of days after creation when objects are transitioned to the specified storage class.
	// Experimental.
	TransitionAfter awscdk.Duration `json:"transitionAfter"`
	// Indicates when objects are transitioned to the specified storage class.
	//
	// The date value must be in ISO 8601 format. The time is always midnight UTC.
	// Experimental.
	TransitionDate *time.Time `json:"transitionDate"`
}

// Options for creating Virtual-Hosted style URL.
// Experimental.
type VirtualHostedStyleUrlOptions struct {
	// Specifies the URL includes the region.
	// Experimental.
	Regional *bool `json:"regional"`
}


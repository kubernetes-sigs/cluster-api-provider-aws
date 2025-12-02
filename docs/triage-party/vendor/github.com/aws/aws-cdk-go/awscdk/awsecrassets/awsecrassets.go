package awsecrassets

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/assets"
	"github.com/aws/aws-cdk-go/awscdk/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/awsecrassets/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// An asset that represents a Docker image.
//
// The image will be created in build time and uploaded to an ECR repository.
// Experimental.
type DockerImageAsset interface {
	awscdk.Construct
	assets.IAsset
	AssetHash() *string
	ImageUri() *string
	SetImageUri(val *string)
	Node() awscdk.ConstructNode
	Repository() awsecr.IRepository
	SetRepository(val awsecr.IRepository)
	SourceHash() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for DockerImageAsset
type jsiiProxy_DockerImageAsset struct {
	internal.Type__awscdkConstruct
	internal.Type__assetsIAsset
}

func (j *jsiiProxy_DockerImageAsset) AssetHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assetHash",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageAsset) ImageUri() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageUri",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageAsset) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageAsset) Repository() awsecr.IRepository {
	var returns awsecr.IRepository
	_jsii_.Get(
		j,
		"repository",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerImageAsset) SourceHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceHash",
		&returns,
	)
	return returns
}


// Experimental.
func NewDockerImageAsset(scope constructs.Construct, id *string, props *DockerImageAssetProps) DockerImageAsset {
	_init_.Initialize()

	j := jsiiProxy_DockerImageAsset{}

	_jsii_.Create(
		"monocdk.aws_ecr_assets.DockerImageAsset",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewDockerImageAsset_Override(d DockerImageAsset, scope constructs.Construct, id *string, props *DockerImageAssetProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_ecr_assets.DockerImageAsset",
		[]interface{}{scope, id, props},
		d,
	)
}

func (j *jsiiProxy_DockerImageAsset) SetImageUri(val *string) {
	_jsii_.Set(
		j,
		"imageUri",
		val,
	)
}

func (j *jsiiProxy_DockerImageAsset) SetRepository(val awsecr.IRepository) {
	_jsii_.Set(
		j,
		"repository",
		val,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func DockerImageAsset_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_ecr_assets.DockerImageAsset",
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
func (d *jsiiProxy_DockerImageAsset) OnPrepare() {
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
func (d *jsiiProxy_DockerImageAsset) OnSynthesize(session constructs.ISynthesisSession) {
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
func (d *jsiiProxy_DockerImageAsset) OnValidate() *[]*string {
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
func (d *jsiiProxy_DockerImageAsset) Prepare() {
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
func (d *jsiiProxy_DockerImageAsset) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (d *jsiiProxy_DockerImageAsset) ToString() *string {
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
func (d *jsiiProxy_DockerImageAsset) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		d,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options for DockerImageAsset.
// Experimental.
type DockerImageAssetOptions struct {
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

// Props for DockerImageAssets.
// Experimental.
type DockerImageAssetProps struct {
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
	// The directory where the Dockerfile is stored.
	// Experimental.
	Directory *string `json:"directory"`
}


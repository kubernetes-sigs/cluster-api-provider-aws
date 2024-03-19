package assets

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/assets/internal"
	"github.com/aws/constructs-go/constructs/v3"
)

// Obtains applied when copying directories into the staging location.
// Deprecated: see `core.CopyOptions`
type CopyOptions struct {
	// Glob patterns to exclude from the copy.
	// Deprecated: see `core.CopyOptions`
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Deprecated: use `followSymlinks` instead
	Follow FollowMode `json:"follow"`
	// The ignore behavior to use for exclude patterns.
	// Deprecated: see `core.CopyOptions`
	IgnoreMode awscdk.IgnoreMode `json:"ignoreMode"`
}

// Options related to calculating source hash.
// Deprecated: see `core.FingerprintOptions`
type FingerprintOptions struct {
	// Glob patterns to exclude from the copy.
	// Deprecated: see `core.FingerprintOptions`
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Deprecated: use `followSymlinks` instead
	Follow FollowMode `json:"follow"`
	// The ignore behavior to use for exclude patterns.
	// Deprecated: see `core.FingerprintOptions`
	IgnoreMode awscdk.IgnoreMode `json:"ignoreMode"`
	// Extra information to encode into the fingerprint (e.g. build instructions and other inputs).
	// Deprecated: see `core.FingerprintOptions`
	ExtraHash *string `json:"extraHash"`
}

// Symlink follow mode.
// Deprecated: see `core.SymlinkFollowMode`
type FollowMode string

const (
	FollowMode_NEVER FollowMode = "NEVER"
	FollowMode_ALWAYS FollowMode = "ALWAYS"
	FollowMode_EXTERNAL FollowMode = "EXTERNAL"
	FollowMode_BLOCK_EXTERNAL FollowMode = "BLOCK_EXTERNAL"
)

// Common interface for all assets.
// Deprecated: use `core.IAsset`
type IAsset interface {
	// A hash of the source of this asset, which is available at construction time.
	//
	// As this is a plain
	// string, it can be used in construct IDs in order to enforce creation of a new resource when
	// the content hash has changed.
	// Deprecated: use `core.IAsset`
	SourceHash() *string
}

// The jsii proxy for IAsset
type jsiiProxy_IAsset struct {
	_ byte // padding
}

func (j *jsiiProxy_IAsset) SourceHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceHash",
		&returns,
	)
	return returns
}

// Deprecated.
// Deprecated: use `core.AssetStaging`
type Staging interface {
	awscdk.AssetStaging
	AbsoluteStagedPath() *string
	AssetHash() *string
	IsArchive() *bool
	Node() awscdk.ConstructNode
	Packaging() awscdk.FileAssetPackaging
	SourceHash() *string
	SourcePath() *string
	StagedPath() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RelativeStagedPath(stack awscdk.Stack) *string
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Staging
type jsiiProxy_Staging struct {
	internal.Type__awscdkAssetStaging
}

func (j *jsiiProxy_Staging) AbsoluteStagedPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"absoluteStagedPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Staging) AssetHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"assetHash",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Staging) IsArchive() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"isArchive",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Staging) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Staging) Packaging() awscdk.FileAssetPackaging {
	var returns awscdk.FileAssetPackaging
	_jsii_.Get(
		j,
		"packaging",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Staging) SourceHash() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourceHash",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Staging) SourcePath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sourcePath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Staging) StagedPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stagedPath",
		&returns,
	)
	return returns
}


// Deprecated: use `core.AssetStaging`
func NewStaging(scope awscdk.Construct, id *string, props *StagingProps) Staging {
	_init_.Initialize()

	j := jsiiProxy_Staging{}

	_jsii_.Create(
		"monocdk.assets.Staging",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Deprecated: use `core.AssetStaging`
func NewStaging_Override(s Staging, scope awscdk.Construct, id *string, props *StagingProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.assets.Staging",
		[]interface{}{scope, id, props},
		s,
	)
}

// Clears the asset hash cache.
// Deprecated: use `core.AssetStaging`
func Staging_ClearAssetHashCache() {
	_init_.Initialize()

	_jsii_.StaticInvokeVoid(
		"monocdk.assets.Staging",
		"clearAssetHashCache",
		nil, // no parameters
	)
}

// Return whether the given object is a Construct.
// Deprecated: use `core.AssetStaging`
func Staging_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.assets.Staging",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func Staging_BUNDLING_INPUT_DIR() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.assets.Staging",
		"BUNDLING_INPUT_DIR",
		&returns,
	)
	return returns
}

func Staging_BUNDLING_OUTPUT_DIR() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.assets.Staging",
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
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) OnPrepare() {
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
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) OnValidate() *[]*string {
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
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) Prepare() {
	_jsii_.InvokeVoid(
		s,
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
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) RelativeStagedPath(stack awscdk.Stack) *string {
	var returns *string

	_jsii_.Invoke(
		s,
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
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) ToString() *string {
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
// Deprecated: use `core.AssetStaging`
func (s *jsiiProxy_Staging) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated.
// Deprecated: use `core.AssetStagingProps`
type StagingProps struct {
	// Glob patterns to exclude from the copy.
	// Deprecated: use `core.AssetStagingProps`
	Exclude *[]*string `json:"exclude"`
	// A strategy for how to handle symlinks.
	// Deprecated: use `followSymlinks` instead
	Follow FollowMode `json:"follow"`
	// The ignore behavior to use for exclude patterns.
	// Deprecated: use `core.AssetStagingProps`
	IgnoreMode awscdk.IgnoreMode `json:"ignoreMode"`
	// Extra information to encode into the fingerprint (e.g. build instructions and other inputs).
	// Deprecated: use `core.AssetStagingProps`
	ExtraHash *string `json:"extraHash"`
	// Local file or directory to stage.
	// Deprecated: use `core.AssetStagingProps`
	SourcePath *string `json:"sourcePath"`
}


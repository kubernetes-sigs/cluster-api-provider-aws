package assets

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterStruct(
		"monocdk.assets.CopyOptions",
		reflect.TypeOf((*CopyOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.assets.FingerprintOptions",
		reflect.TypeOf((*FingerprintOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.assets.FollowMode",
		reflect.TypeOf((*FollowMode)(nil)).Elem(),
		map[string]interface{}{
			"NEVER": FollowMode_NEVER,
			"ALWAYS": FollowMode_ALWAYS,
			"EXTERNAL": FollowMode_EXTERNAL,
			"BLOCK_EXTERNAL": FollowMode_BLOCK_EXTERNAL,
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.assets.IAsset",
		reflect.TypeOf((*IAsset)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "sourceHash", GoGetter: "SourceHash"},
		},
		func() interface{} {
			return &jsiiProxy_IAsset{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.assets.Staging",
		reflect.TypeOf((*Staging)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "absoluteStagedPath", GoGetter: "AbsoluteStagedPath"},
			_jsii_.MemberProperty{JsiiProperty: "assetHash", GoGetter: "AssetHash"},
			_jsii_.MemberProperty{JsiiProperty: "isArchive", GoGetter: "IsArchive"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "packaging", GoGetter: "Packaging"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "relativeStagedPath", GoMethod: "RelativeStagedPath"},
			_jsii_.MemberProperty{JsiiProperty: "sourceHash", GoGetter: "SourceHash"},
			_jsii_.MemberProperty{JsiiProperty: "sourcePath", GoGetter: "SourcePath"},
			_jsii_.MemberProperty{JsiiProperty: "stagedPath", GoGetter: "StagedPath"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Staging{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkAssetStaging)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.assets.StagingProps",
		reflect.TypeOf((*StagingProps)(nil)).Elem(),
	)
}

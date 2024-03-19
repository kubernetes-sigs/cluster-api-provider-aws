package constructs

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"constructs.Construct",
		reflect.TypeOf((*Construct)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_Construct{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IConstruct)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"constructs.ConstructMetadata",
		reflect.TypeOf((*ConstructMetadata)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_ConstructMetadata{}
		},
	)
	_jsii_.RegisterStruct(
		"constructs.ConstructOptions",
		reflect.TypeOf((*ConstructOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"constructs.ConstructOrder",
		reflect.TypeOf((*ConstructOrder)(nil)).Elem(),
		map[string]interface{}{
			"PREORDER": ConstructOrder_PREORDER,
			"POSTORDER": ConstructOrder_POSTORDER,
		},
	)
	_jsii_.RegisterStruct(
		"constructs.Dependency",
		reflect.TypeOf((*Dependency)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"constructs.IAspect",
		reflect.TypeOf((*IAspect)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "visit", GoMethod: "Visit"},
		},
		func() interface{} {
			return &jsiiProxy_IAspect{}
		},
	)
	_jsii_.RegisterInterface(
		"constructs.IConstruct",
		reflect.TypeOf((*IConstruct)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_IConstruct{}
		},
	)
	_jsii_.RegisterInterface(
		"constructs.INodeFactory",
		reflect.TypeOf((*INodeFactory)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "createNode", GoMethod: "CreateNode"},
		},
		func() interface{} {
			return &jsiiProxy_INodeFactory{}
		},
	)
	_jsii_.RegisterInterface(
		"constructs.ISynthesisSession",
		reflect.TypeOf((*ISynthesisSession)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "outdir", GoGetter: "Outdir"},
		},
		func() interface{} {
			return &jsiiProxy_ISynthesisSession{}
		},
	)
	_jsii_.RegisterInterface(
		"constructs.IValidation",
		reflect.TypeOf((*IValidation)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			return &jsiiProxy_IValidation{}
		},
	)
	_jsii_.RegisterStruct(
		"constructs.MetadataEntry",
		reflect.TypeOf((*MetadataEntry)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"constructs.Node",
		reflect.TypeOf((*Node)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "addError", GoMethod: "AddError"},
			_jsii_.MemberMethod{JsiiMethod: "addInfo", GoMethod: "AddInfo"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "addr", GoGetter: "Addr"},
			_jsii_.MemberMethod{JsiiMethod: "addValidation", GoMethod: "AddValidation"},
			_jsii_.MemberMethod{JsiiMethod: "addWarning", GoMethod: "AddWarning"},
			_jsii_.MemberMethod{JsiiMethod: "applyAspect", GoMethod: "ApplyAspect"},
			_jsii_.MemberProperty{JsiiProperty: "children", GoGetter: "Children"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChild", GoGetter: "DefaultChild"},
			_jsii_.MemberProperty{JsiiProperty: "dependencies", GoGetter: "Dependencies"},
			_jsii_.MemberMethod{JsiiMethod: "findAll", GoMethod: "FindAll"},
			_jsii_.MemberMethod{JsiiMethod: "findChild", GoMethod: "FindChild"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "locked", GoGetter: "Locked"},
			_jsii_.MemberProperty{JsiiProperty: "metadata", GoGetter: "Metadata"},
			_jsii_.MemberProperty{JsiiProperty: "path", GoGetter: "Path"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "root", GoGetter: "Root"},
			_jsii_.MemberProperty{JsiiProperty: "scope", GoGetter: "Scope"},
			_jsii_.MemberProperty{JsiiProperty: "scopes", GoGetter: "Scopes"},
			_jsii_.MemberMethod{JsiiMethod: "setContext", GoMethod: "SetContext"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "tryFindChild", GoMethod: "TryFindChild"},
			_jsii_.MemberMethod{JsiiMethod: "tryGetContext", GoMethod: "TryGetContext"},
			_jsii_.MemberMethod{JsiiMethod: "tryRemoveChild", GoMethod: "TryRemoveChild"},
			_jsii_.MemberProperty{JsiiProperty: "uniqueId", GoGetter: "UniqueId"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			return &jsiiProxy_Node{}
		},
	)
	_jsii_.RegisterStruct(
		"constructs.SynthesisOptions",
		reflect.TypeOf((*SynthesisOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"constructs.ValidationError",
		reflect.TypeOf((*ValidationError)(nil)).Elem(),
	)
}

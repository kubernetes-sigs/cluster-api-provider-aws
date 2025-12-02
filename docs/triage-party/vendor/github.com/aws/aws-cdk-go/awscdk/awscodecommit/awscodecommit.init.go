package awscodecommit

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_codecommit.CfnRepository",
		reflect.TypeOf((*CfnRepository)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrCloneUrlHttp", GoGetter: "AttrCloneUrlHttp"},
			_jsii_.MemberProperty{JsiiProperty: "attrCloneUrlSsh", GoGetter: "AttrCloneUrlSsh"},
			_jsii_.MemberProperty{JsiiProperty: "attrName", GoGetter: "AttrName"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "code", GoGetter: "Code"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryDescription", GoGetter: "RepositoryDescription"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "triggers", GoGetter: "Triggers"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnRepository{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codecommit.CfnRepository.CodeProperty",
		reflect.TypeOf((*CfnRepository_CodeProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codecommit.CfnRepository.RepositoryTriggerProperty",
		reflect.TypeOf((*CfnRepository_RepositoryTriggerProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codecommit.CfnRepository.S3Property",
		reflect.TypeOf((*CfnRepository_S3Property)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codecommit.CfnRepositoryProps",
		reflect.TypeOf((*CfnRepositoryProps)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_codecommit.IRepository",
		reflect.TypeOf((*IRepository)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantPull", GoMethod: "GrantPull"},
			_jsii_.MemberMethod{JsiiMethod: "grantPullPush", GoMethod: "GrantPullPush"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onCommentOnCommit", GoMethod: "OnCommentOnCommit"},
			_jsii_.MemberMethod{JsiiMethod: "onCommentOnPullRequest", GoMethod: "OnCommentOnPullRequest"},
			_jsii_.MemberMethod{JsiiMethod: "onCommit", GoMethod: "OnCommit"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onPullRequestStateChange", GoMethod: "OnPullRequestStateChange"},
			_jsii_.MemberMethod{JsiiMethod: "onReferenceCreated", GoMethod: "OnReferenceCreated"},
			_jsii_.MemberMethod{JsiiMethod: "onReferenceDeleted", GoMethod: "OnReferenceDeleted"},
			_jsii_.MemberMethod{JsiiMethod: "onReferenceUpdated", GoMethod: "OnReferenceUpdated"},
			_jsii_.MemberMethod{JsiiMethod: "onStateChange", GoMethod: "OnStateChange"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryArn", GoGetter: "RepositoryArn"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryCloneUrlGrc", GoGetter: "RepositoryCloneUrlGrc"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryCloneUrlHttp", GoGetter: "RepositoryCloneUrlHttp"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryCloneUrlSsh", GoGetter: "RepositoryCloneUrlSsh"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IRepository{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codecommit.OnCommitOptions",
		reflect.TypeOf((*OnCommitOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codecommit.ReferenceEvent",
		reflect.TypeOf((*ReferenceEvent)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_ReferenceEvent{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_codecommit.Repository",
		reflect.TypeOf((*Repository)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantPull", GoMethod: "GrantPull"},
			_jsii_.MemberMethod{JsiiMethod: "grantPullPush", GoMethod: "GrantPullPush"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "notify", GoMethod: "Notify"},
			_jsii_.MemberMethod{JsiiMethod: "onCommentOnCommit", GoMethod: "OnCommentOnCommit"},
			_jsii_.MemberMethod{JsiiMethod: "onCommentOnPullRequest", GoMethod: "OnCommentOnPullRequest"},
			_jsii_.MemberMethod{JsiiMethod: "onCommit", GoMethod: "OnCommit"},
			_jsii_.MemberMethod{JsiiMethod: "onEvent", GoMethod: "OnEvent"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onPullRequestStateChange", GoMethod: "OnPullRequestStateChange"},
			_jsii_.MemberMethod{JsiiMethod: "onReferenceCreated", GoMethod: "OnReferenceCreated"},
			_jsii_.MemberMethod{JsiiMethod: "onReferenceDeleted", GoMethod: "OnReferenceDeleted"},
			_jsii_.MemberMethod{JsiiMethod: "onReferenceUpdated", GoMethod: "OnReferenceUpdated"},
			_jsii_.MemberMethod{JsiiMethod: "onStateChange", GoMethod: "OnStateChange"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryArn", GoGetter: "RepositoryArn"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryCloneUrlGrc", GoGetter: "RepositoryCloneUrlGrc"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryCloneUrlHttp", GoGetter: "RepositoryCloneUrlHttp"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryCloneUrlSsh", GoGetter: "RepositoryCloneUrlSsh"},
			_jsii_.MemberProperty{JsiiProperty: "repositoryName", GoGetter: "RepositoryName"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Repository{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IRepository)
			return &j
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_codecommit.RepositoryEventTrigger",
		reflect.TypeOf((*RepositoryEventTrigger)(nil)).Elem(),
		map[string]interface{}{
			"ALL": RepositoryEventTrigger_ALL,
			"UPDATE_REF": RepositoryEventTrigger_UPDATE_REF,
			"CREATE_REF": RepositoryEventTrigger_CREATE_REF,
			"DELETE_REF": RepositoryEventTrigger_DELETE_REF,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codecommit.RepositoryProps",
		reflect.TypeOf((*RepositoryProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_codecommit.RepositoryTriggerOptions",
		reflect.TypeOf((*RepositoryTriggerOptions)(nil)).Elem(),
	)
}

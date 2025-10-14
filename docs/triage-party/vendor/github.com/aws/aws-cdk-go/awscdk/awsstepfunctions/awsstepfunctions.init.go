package awsstepfunctions

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Activity",
		reflect.TypeOf((*Activity)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "activityArn", GoGetter: "ActivityArn"},
			_jsii_.MemberProperty{JsiiProperty: "activityName", GoGetter: "ActivityName"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailed", GoMethod: "MetricFailed"},
			_jsii_.MemberMethod{JsiiMethod: "metricHeartbeatTimedOut", GoMethod: "MetricHeartbeatTimedOut"},
			_jsii_.MemberMethod{JsiiMethod: "metricRunTime", GoMethod: "MetricRunTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricScheduled", GoMethod: "MetricScheduled"},
			_jsii_.MemberMethod{JsiiMethod: "metricScheduleTime", GoMethod: "MetricScheduleTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricStarted", GoMethod: "MetricStarted"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceeded", GoMethod: "MetricSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "metricTime", GoMethod: "MetricTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricTimedOut", GoMethod: "MetricTimedOut"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Activity{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IActivity)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.ActivityProps",
		reflect.TypeOf((*ActivityProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.AfterwardsOptions",
		reflect.TypeOf((*AfterwardsOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CatchProps",
		reflect.TypeOf((*CatchProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.CfnActivity",
		reflect.TypeOf((*CfnActivity)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrName", GoGetter: "AttrName"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnActivity{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnActivity.TagsEntryProperty",
		reflect.TypeOf((*CfnActivity_TagsEntryProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnActivityProps",
		reflect.TypeOf((*CfnActivityProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.CfnStateMachine",
		reflect.TypeOf((*CfnStateMachine)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrName", GoGetter: "AttrName"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "creationStack", GoGetter: "CreationStack"},
			_jsii_.MemberProperty{JsiiProperty: "definition", GoGetter: "Definition"},
			_jsii_.MemberProperty{JsiiProperty: "definitionS3Location", GoGetter: "DefinitionS3Location"},
			_jsii_.MemberProperty{JsiiProperty: "definitionString", GoGetter: "DefinitionString"},
			_jsii_.MemberProperty{JsiiProperty: "definitionSubstitutions", GoGetter: "DefinitionSubstitutions"},
			_jsii_.MemberMethod{JsiiMethod: "getAtt", GoMethod: "GetAtt"},
			_jsii_.MemberMethod{JsiiMethod: "getMetadata", GoMethod: "GetMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "inspect", GoMethod: "Inspect"},
			_jsii_.MemberProperty{JsiiProperty: "loggingConfiguration", GoGetter: "LoggingConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "logicalId", GoGetter: "LogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "ref", GoGetter: "Ref"},
			_jsii_.MemberMethod{JsiiMethod: "renderProperties", GoMethod: "RenderProperties"},
			_jsii_.MemberProperty{JsiiProperty: "roleArn", GoGetter: "RoleArn"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "stateMachineName", GoGetter: "StateMachineName"},
			_jsii_.MemberProperty{JsiiProperty: "stateMachineType", GoGetter: "StateMachineType"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "tracingConfiguration", GoGetter: "TracingConfiguration"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnStateMachine{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachine.CloudWatchLogsLogGroupProperty",
		reflect.TypeOf((*CfnStateMachine_CloudWatchLogsLogGroupProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachine.DefinitionProperty",
		reflect.TypeOf((*CfnStateMachine_DefinitionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachine.LogDestinationProperty",
		reflect.TypeOf((*CfnStateMachine_LogDestinationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachine.LoggingConfigurationProperty",
		reflect.TypeOf((*CfnStateMachine_LoggingConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachine.S3LocationProperty",
		reflect.TypeOf((*CfnStateMachine_S3LocationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachine.TagsEntryProperty",
		reflect.TypeOf((*CfnStateMachine_TagsEntryProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachine.TracingConfigurationProperty",
		reflect.TypeOf((*CfnStateMachine_TracingConfigurationProperty)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CfnStateMachineProps",
		reflect.TypeOf((*CfnStateMachineProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Chain",
		reflect.TypeOf((*Chain)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberMethod{JsiiMethod: "toSingleState", GoMethod: "ToSingleState"},
		},
		func() interface{} {
			j := jsiiProxy_Chain{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IChainable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Choice",
		reflect.TypeOf((*Choice)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "afterwards", GoMethod: "Afterwards"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "otherwise", GoMethod: "Otherwise"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "when", GoMethod: "When"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Choice{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.ChoiceProps",
		reflect.TypeOf((*ChoiceProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Condition",
		reflect.TypeOf((*Condition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "renderCondition", GoMethod: "RenderCondition"},
		},
		func() interface{} {
			return &jsiiProxy_Condition{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Context",
		reflect.TypeOf((*Context)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Context{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.CustomState",
		reflect.TypeOf((*CustomState)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_CustomState{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IChainable)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INextable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.CustomStateProps",
		reflect.TypeOf((*CustomStateProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Data",
		reflect.TypeOf((*Data)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Data{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Errors",
		reflect.TypeOf((*Errors)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Errors{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Fail",
		reflect.TypeOf((*Fail)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Fail{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.FailProps",
		reflect.TypeOf((*FailProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.FieldUtils",
		reflect.TypeOf((*FieldUtils)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_FieldUtils{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.FindStateOptions",
		reflect.TypeOf((*FindStateOptions)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_stepfunctions.IActivity",
		reflect.TypeOf((*IActivity)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "activityArn", GoGetter: "ActivityArn"},
			_jsii_.MemberProperty{JsiiProperty: "activityName", GoGetter: "ActivityName"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
		},
		func() interface{} {
			j := jsiiProxy_IActivity{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_stepfunctions.IChainable",
		reflect.TypeOf((*IChainable)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
		},
		func() interface{} {
			return &jsiiProxy_IChainable{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_stepfunctions.INextable",
		reflect.TypeOf((*INextable)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
		},
		func() interface{} {
			return &jsiiProxy_INextable{}
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_stepfunctions.IStateMachine",
		reflect.TypeOf((*IStateMachine)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantExecution", GoMethod: "GrantExecution"},
			_jsii_.MemberProperty{JsiiProperty: "grantPrincipal", GoGetter: "GrantPrincipal"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantStartExecution", GoMethod: "GrantStartExecution"},
			_jsii_.MemberMethod{JsiiMethod: "grantTaskResponse", GoMethod: "GrantTaskResponse"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricAborted", GoMethod: "MetricAborted"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailed", GoMethod: "MetricFailed"},
			_jsii_.MemberMethod{JsiiMethod: "metricStarted", GoMethod: "MetricStarted"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceeded", GoMethod: "MetricSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "metricThrottled", GoMethod: "MetricThrottled"},
			_jsii_.MemberMethod{JsiiMethod: "metricTime", GoMethod: "MetricTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricTimedOut", GoMethod: "MetricTimedOut"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "stateMachineArn", GoGetter: "StateMachineArn"},
		},
		func() interface{} {
			j := jsiiProxy_IStateMachine{}
			_jsii_.InitJsiiProxy(&j.Type__awsiamIGrantable)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_stepfunctions.IStepFunctionsTask",
		reflect.TypeOf((*IStepFunctionsTask)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			return &jsiiProxy_IStepFunctionsTask{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_stepfunctions.InputType",
		reflect.TypeOf((*InputType)(nil)).Elem(),
		map[string]interface{}{
			"TEXT": InputType_TEXT,
			"OBJECT": InputType_OBJECT,
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_stepfunctions.IntegrationPattern",
		reflect.TypeOf((*IntegrationPattern)(nil)).Elem(),
		map[string]interface{}{
			"REQUEST_RESPONSE": IntegrationPattern_REQUEST_RESPONSE,
			"RUN_JOB": IntegrationPattern_RUN_JOB,
			"WAIT_FOR_TASK_TOKEN": IntegrationPattern_WAIT_FOR_TASK_TOKEN,
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.JsonPath",
		reflect.TypeOf((*JsonPath)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_JsonPath{}
		},
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_stepfunctions.LogLevel",
		reflect.TypeOf((*LogLevel)(nil)).Elem(),
		map[string]interface{}{
			"OFF": LogLevel_OFF,
			"ALL": LogLevel_ALL,
			"ERROR": LogLevel_ERROR,
			"FATAL": LogLevel_FATAL,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.LogOptions",
		reflect.TypeOf((*LogOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Map",
		reflect.TypeOf((*Map)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addCatch", GoMethod: "AddCatch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "addRetry", GoMethod: "AddRetry"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "iterator", GoMethod: "Iterator"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Map{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INextable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.MapProps",
		reflect.TypeOf((*MapProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Parallel",
		reflect.TypeOf((*Parallel)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addCatch", GoMethod: "AddCatch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "addRetry", GoMethod: "AddRetry"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberMethod{JsiiMethod: "branch", GoMethod: "Branch"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Parallel{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INextable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.ParallelProps",
		reflect.TypeOf((*ParallelProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Pass",
		reflect.TypeOf((*Pass)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Pass{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INextable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.PassProps",
		reflect.TypeOf((*PassProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Result",
		reflect.TypeOf((*Result)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
		},
		func() interface{} {
			return &jsiiProxy_Result{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.RetryProps",
		reflect.TypeOf((*RetryProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_stepfunctions.ServiceIntegrationPattern",
		reflect.TypeOf((*ServiceIntegrationPattern)(nil)).Elem(),
		map[string]interface{}{
			"FIRE_AND_FORGET": ServiceIntegrationPattern_FIRE_AND_FORGET,
			"SYNC": ServiceIntegrationPattern_SYNC,
			"WAIT_FOR_TASK_TOKEN": ServiceIntegrationPattern_WAIT_FOR_TASK_TOKEN,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.SingleStateOptions",
		reflect.TypeOf((*SingleStateOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.State",
		reflect.TypeOf((*State)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_State{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IChainable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.StateGraph",
		reflect.TypeOf((*StateGraph)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "policyStatements", GoGetter: "PolicyStatements"},
			_jsii_.MemberMethod{JsiiMethod: "registerPolicyStatement", GoMethod: "RegisterPolicyStatement"},
			_jsii_.MemberMethod{JsiiMethod: "registerState", GoMethod: "RegisterState"},
			_jsii_.MemberMethod{JsiiMethod: "registerSuperGraph", GoMethod: "RegisterSuperGraph"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "timeout", GoGetter: "Timeout"},
			_jsii_.MemberMethod{JsiiMethod: "toGraphJson", GoMethod: "ToGraphJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			return &jsiiProxy_StateGraph{}
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.StateMachine",
		reflect.TypeOf((*StateMachine)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addToRolePolicy", GoMethod: "AddToRolePolicy"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantExecution", GoMethod: "GrantExecution"},
			_jsii_.MemberProperty{JsiiProperty: "grantPrincipal", GoGetter: "GrantPrincipal"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantStartExecution", GoMethod: "GrantStartExecution"},
			_jsii_.MemberMethod{JsiiMethod: "grantTaskResponse", GoMethod: "GrantTaskResponse"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricAborted", GoMethod: "MetricAborted"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailed", GoMethod: "MetricFailed"},
			_jsii_.MemberMethod{JsiiMethod: "metricStarted", GoMethod: "MetricStarted"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceeded", GoMethod: "MetricSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "metricThrottled", GoMethod: "MetricThrottled"},
			_jsii_.MemberMethod{JsiiMethod: "metricTime", GoMethod: "MetricTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricTimedOut", GoMethod: "MetricTimedOut"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "role", GoGetter: "Role"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "stateMachineArn", GoGetter: "StateMachineArn"},
			_jsii_.MemberProperty{JsiiProperty: "stateMachineName", GoGetter: "StateMachineName"},
			_jsii_.MemberProperty{JsiiProperty: "stateMachineType", GoGetter: "StateMachineType"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_StateMachine{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IStateMachine)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.StateMachineFragment",
		reflect.TypeOf((*StateMachineFragment)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberMethod{JsiiMethod: "prefixStates", GoMethod: "PrefixStates"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toSingleState", GoMethod: "ToSingleState"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_StateMachineFragment{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkConstruct)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IChainable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.StateMachineProps",
		reflect.TypeOf((*StateMachineProps)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_stepfunctions.StateMachineType",
		reflect.TypeOf((*StateMachineType)(nil)).Elem(),
		map[string]interface{}{
			"EXPRESS": StateMachineType_EXPRESS,
			"STANDARD": StateMachineType_STANDARD,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.StateProps",
		reflect.TypeOf((*StateProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.StateTransitionMetric",
		reflect.TypeOf((*StateTransitionMetric)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_StateTransitionMetric{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.StepFunctionsTaskConfig",
		reflect.TypeOf((*StepFunctionsTaskConfig)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Succeed",
		reflect.TypeOf((*Succeed)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Succeed{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.SucceedProps",
		reflect.TypeOf((*SucceedProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Task",
		reflect.TypeOf((*Task)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addCatch", GoMethod: "AddCatch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "addRetry", GoMethod: "AddRetry"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailed", GoMethod: "MetricFailed"},
			_jsii_.MemberMethod{JsiiMethod: "metricHeartbeatTimedOut", GoMethod: "MetricHeartbeatTimedOut"},
			_jsii_.MemberMethod{JsiiMethod: "metricRunTime", GoMethod: "MetricRunTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricScheduled", GoMethod: "MetricScheduled"},
			_jsii_.MemberMethod{JsiiMethod: "metricScheduleTime", GoMethod: "MetricScheduleTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricStarted", GoMethod: "MetricStarted"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceeded", GoMethod: "MetricSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "metricTime", GoMethod: "MetricTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricTimedOut", GoMethod: "MetricTimedOut"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Task{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INextable)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.TaskInput",
		reflect.TypeOf((*TaskInput)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "type", GoGetter: "Type"},
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
		},
		func() interface{} {
			return &jsiiProxy_TaskInput{}
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.TaskMetricsConfig",
		reflect.TypeOf((*TaskMetricsConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.TaskProps",
		reflect.TypeOf((*TaskProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.TaskStateBase",
		reflect.TypeOf((*TaskStateBase)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addCatch", GoMethod: "AddCatch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "addRetry", GoMethod: "AddRetry"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricFailed", GoMethod: "MetricFailed"},
			_jsii_.MemberMethod{JsiiMethod: "metricHeartbeatTimedOut", GoMethod: "MetricHeartbeatTimedOut"},
			_jsii_.MemberMethod{JsiiMethod: "metricRunTime", GoMethod: "MetricRunTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricScheduled", GoMethod: "MetricScheduled"},
			_jsii_.MemberMethod{JsiiMethod: "metricScheduleTime", GoMethod: "MetricScheduleTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricStarted", GoMethod: "MetricStarted"},
			_jsii_.MemberMethod{JsiiMethod: "metricSucceeded", GoMethod: "MetricSucceeded"},
			_jsii_.MemberMethod{JsiiMethod: "metricTime", GoMethod: "MetricTime"},
			_jsii_.MemberMethod{JsiiMethod: "metricTimedOut", GoMethod: "MetricTimedOut"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "taskMetrics", GoGetter: "TaskMetrics"},
			_jsii_.MemberProperty{JsiiProperty: "taskPolicies", GoGetter: "TaskPolicies"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_TaskStateBase{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INextable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.TaskStateBaseProps",
		reflect.TypeOf((*TaskStateBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.Wait",
		reflect.TypeOf((*Wait)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addBranch", GoMethod: "AddBranch"},
			_jsii_.MemberMethod{JsiiMethod: "addChoice", GoMethod: "AddChoice"},
			_jsii_.MemberMethod{JsiiMethod: "addIterator", GoMethod: "AddIterator"},
			_jsii_.MemberMethod{JsiiMethod: "addPrefix", GoMethod: "AddPrefix"},
			_jsii_.MemberMethod{JsiiMethod: "bindToGraph", GoMethod: "BindToGraph"},
			_jsii_.MemberProperty{JsiiProperty: "branches", GoGetter: "Branches"},
			_jsii_.MemberProperty{JsiiProperty: "comment", GoGetter: "Comment"},
			_jsii_.MemberProperty{JsiiProperty: "defaultChoice", GoGetter: "DefaultChoice"},
			_jsii_.MemberProperty{JsiiProperty: "endStates", GoGetter: "EndStates"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "inputPath", GoGetter: "InputPath"},
			_jsii_.MemberProperty{JsiiProperty: "iteration", GoGetter: "Iteration"},
			_jsii_.MemberMethod{JsiiMethod: "makeDefault", GoMethod: "MakeDefault"},
			_jsii_.MemberMethod{JsiiMethod: "makeNext", GoMethod: "MakeNext"},
			_jsii_.MemberMethod{JsiiMethod: "next", GoMethod: "Next"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "outputPath", GoGetter: "OutputPath"},
			_jsii_.MemberProperty{JsiiProperty: "parameters", GoGetter: "Parameters"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberMethod{JsiiMethod: "renderBranches", GoMethod: "RenderBranches"},
			_jsii_.MemberMethod{JsiiMethod: "renderChoices", GoMethod: "RenderChoices"},
			_jsii_.MemberMethod{JsiiMethod: "renderInputOutput", GoMethod: "RenderInputOutput"},
			_jsii_.MemberMethod{JsiiMethod: "renderIterator", GoMethod: "RenderIterator"},
			_jsii_.MemberMethod{JsiiMethod: "renderNextEnd", GoMethod: "RenderNextEnd"},
			_jsii_.MemberMethod{JsiiMethod: "renderResultSelector", GoMethod: "RenderResultSelector"},
			_jsii_.MemberMethod{JsiiMethod: "renderRetryCatch", GoMethod: "RenderRetryCatch"},
			_jsii_.MemberProperty{JsiiProperty: "resultPath", GoGetter: "ResultPath"},
			_jsii_.MemberProperty{JsiiProperty: "resultSelector", GoGetter: "ResultSelector"},
			_jsii_.MemberProperty{JsiiProperty: "startState", GoGetter: "StartState"},
			_jsii_.MemberProperty{JsiiProperty: "stateId", GoGetter: "StateId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toStateJson", GoMethod: "ToStateJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "whenBoundToGraph", GoMethod: "WhenBoundToGraph"},
		},
		func() interface{} {
			j := jsiiProxy_Wait{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_State)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_INextable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_stepfunctions.WaitProps",
		reflect.TypeOf((*WaitProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_stepfunctions.WaitTime",
		reflect.TypeOf((*WaitTime)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_WaitTime{}
		},
	)
}

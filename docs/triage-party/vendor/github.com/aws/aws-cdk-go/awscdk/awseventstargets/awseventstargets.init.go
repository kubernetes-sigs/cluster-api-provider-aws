package awseventstargets

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.ApiGateway",
		reflect.TypeOf((*ApiGateway)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "restApi", GoGetter: "RestApi"},
		},
		func() interface{} {
			j := jsiiProxy_ApiGateway{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.ApiGatewayProps",
		reflect.TypeOf((*ApiGatewayProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.AwsApi",
		reflect.TypeOf((*AwsApi)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_AwsApi{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.AwsApiInput",
		reflect.TypeOf((*AwsApiInput)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.AwsApiProps",
		reflect.TypeOf((*AwsApiProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.BatchJob",
		reflect.TypeOf((*BatchJob)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_BatchJob{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.BatchJobProps",
		reflect.TypeOf((*BatchJobProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.CloudWatchLogGroup",
		reflect.TypeOf((*CloudWatchLogGroup)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_CloudWatchLogGroup{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.CodeBuildProject",
		reflect.TypeOf((*CodeBuildProject)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_CodeBuildProject{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.CodeBuildProjectProps",
		reflect.TypeOf((*CodeBuildProjectProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.CodePipeline",
		reflect.TypeOf((*CodePipeline)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_CodePipeline{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.CodePipelineTargetOptions",
		reflect.TypeOf((*CodePipelineTargetOptions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.ContainerOverride",
		reflect.TypeOf((*ContainerOverride)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.EcsTask",
		reflect.TypeOf((*EcsTask)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "securityGroup", GoGetter: "SecurityGroup"},
			_jsii_.MemberProperty{JsiiProperty: "securityGroups", GoGetter: "SecurityGroups"},
		},
		func() interface{} {
			j := jsiiProxy_EcsTask{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.EcsTaskProps",
		reflect.TypeOf((*EcsTaskProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.EventBus",
		reflect.TypeOf((*EventBus)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_EventBus{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.EventBusProps",
		reflect.TypeOf((*EventBusProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.KinesisFirehoseStream",
		reflect.TypeOf((*KinesisFirehoseStream)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_KinesisFirehoseStream{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.KinesisFirehoseStreamProps",
		reflect.TypeOf((*KinesisFirehoseStreamProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.KinesisStream",
		reflect.TypeOf((*KinesisStream)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_KinesisStream{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.KinesisStreamProps",
		reflect.TypeOf((*KinesisStreamProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.LambdaFunction",
		reflect.TypeOf((*LambdaFunction)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
		},
		func() interface{} {
			j := jsiiProxy_LambdaFunction{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.LambdaFunctionProps",
		reflect.TypeOf((*LambdaFunctionProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.LogGroupProps",
		reflect.TypeOf((*LogGroupProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.SfnStateMachine",
		reflect.TypeOf((*SfnStateMachine)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "machine", GoGetter: "Machine"},
		},
		func() interface{} {
			j := jsiiProxy_SfnStateMachine{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.SfnStateMachineProps",
		reflect.TypeOf((*SfnStateMachineProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.SnsTopic",
		reflect.TypeOf((*SnsTopic)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "topic", GoGetter: "Topic"},
		},
		func() interface{} {
			j := jsiiProxy_SnsTopic{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.SnsTopicProps",
		reflect.TypeOf((*SnsTopicProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_events_targets.SqsQueue",
		reflect.TypeOf((*SqsQueue)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "bind", GoMethod: "Bind"},
			_jsii_.MemberProperty{JsiiProperty: "queue", GoGetter: "Queue"},
		},
		func() interface{} {
			j := jsiiProxy_SqsQueue{}
			_jsii_.InitJsiiProxy(&j.Type__awseventsIRuleTarget)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.SqsQueueProps",
		reflect.TypeOf((*SqsQueueProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.TargetBaseProps",
		reflect.TypeOf((*TargetBaseProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_events_targets.TaskEnvironmentVariable",
		reflect.TypeOf((*TaskEnvironmentVariable)(nil)).Elem(),
	)
}

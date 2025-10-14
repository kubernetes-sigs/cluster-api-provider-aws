package awskinesis

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"monocdk.aws_kinesis.CfnStream",
		reflect.TypeOf((*CfnStream)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrArn", GoGetter: "AttrArn"},
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
			_jsii_.MemberProperty{JsiiProperty: "retentionPeriodHours", GoGetter: "RetentionPeriodHours"},
			_jsii_.MemberProperty{JsiiProperty: "shardCount", GoGetter: "ShardCount"},
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "streamEncryption", GoGetter: "StreamEncryption"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnStream{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_kinesis.CfnStream.StreamEncryptionProperty",
		reflect.TypeOf((*CfnStream_StreamEncryptionProperty)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"monocdk.aws_kinesis.CfnStreamConsumer",
		reflect.TypeOf((*CfnStreamConsumer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDeletionOverride", GoMethod: "AddDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addDependsOn", GoMethod: "AddDependsOn"},
			_jsii_.MemberMethod{JsiiMethod: "addMetadata", GoMethod: "AddMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyDeletionOverride", GoMethod: "AddPropertyDeletionOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addPropertyOverride", GoMethod: "AddPropertyOverride"},
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "attrConsumerArn", GoGetter: "AttrConsumerArn"},
			_jsii_.MemberProperty{JsiiProperty: "attrConsumerCreationTimestamp", GoGetter: "AttrConsumerCreationTimestamp"},
			_jsii_.MemberProperty{JsiiProperty: "attrConsumerName", GoGetter: "AttrConsumerName"},
			_jsii_.MemberProperty{JsiiProperty: "attrConsumerStatus", GoGetter: "AttrConsumerStatus"},
			_jsii_.MemberProperty{JsiiProperty: "attrStreamArn", GoGetter: "AttrStreamArn"},
			_jsii_.MemberProperty{JsiiProperty: "cfnOptions", GoGetter: "CfnOptions"},
			_jsii_.MemberProperty{JsiiProperty: "cfnProperties", GoGetter: "CfnProperties"},
			_jsii_.MemberProperty{JsiiProperty: "cfnResourceType", GoGetter: "CfnResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "consumerName", GoGetter: "ConsumerName"},
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
			_jsii_.MemberMethod{JsiiMethod: "shouldSynthesize", GoMethod: "ShouldSynthesize"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "streamArn", GoGetter: "StreamArn"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "updatedProperites", GoGetter: "UpdatedProperites"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
			_jsii_.MemberMethod{JsiiMethod: "validateProperties", GoMethod: "ValidateProperties"},
		},
		func() interface{} {
			j := jsiiProxy_CfnStreamConsumer{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkCfnResource)
			_jsii_.InitJsiiProxy(&j.Type__awscdkIInspectable)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_kinesis.CfnStreamConsumerProps",
		reflect.TypeOf((*CfnStreamConsumerProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_kinesis.CfnStreamProps",
		reflect.TypeOf((*CfnStreamProps)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"monocdk.aws_kinesis.IStream",
		reflect.TypeOf((*IStream)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantReadWrite", GoMethod: "GrantReadWrite"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecords", GoMethod: "MetricGetRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsBytes", GoMethod: "MetricGetRecordsBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsIteratorAgeMilliseconds", GoMethod: "MetricGetRecordsIteratorAgeMilliseconds"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsLatency", GoMethod: "MetricGetRecordsLatency"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsSuccess", GoMethod: "MetricGetRecordsSuccess"},
			_jsii_.MemberMethod{JsiiMethod: "metricIncomingBytes", GoMethod: "MetricIncomingBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricIncomingRecords", GoMethod: "MetricIncomingRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordBytes", GoMethod: "MetricPutRecordBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordLatency", GoMethod: "MetricPutRecordLatency"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsBytes", GoMethod: "MetricPutRecordsBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsFailedRecords", GoMethod: "MetricPutRecordsFailedRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsLatency", GoMethod: "MetricPutRecordsLatency"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsSuccess", GoMethod: "MetricPutRecordsSuccess"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsSuccessfulRecords", GoMethod: "MetricPutRecordsSuccessfulRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsThrottledRecords", GoMethod: "MetricPutRecordsThrottledRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsTotalRecords", GoMethod: "MetricPutRecordsTotalRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordSuccess", GoMethod: "MetricPutRecordSuccess"},
			_jsii_.MemberMethod{JsiiMethod: "metricReadProvisionedThroughputExceeded", GoMethod: "MetricReadProvisionedThroughputExceeded"},
			_jsii_.MemberMethod{JsiiMethod: "metricWriteProvisionedThroughputExceeded", GoMethod: "MetricWriteProvisionedThroughputExceeded"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "streamArn", GoGetter: "StreamArn"},
			_jsii_.MemberProperty{JsiiProperty: "streamName", GoGetter: "StreamName"},
		},
		func() interface{} {
			j := jsiiProxy_IStream{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkIResource)
			return &j
		},
	)
	_jsii_.RegisterClass(
		"monocdk.aws_kinesis.Stream",
		reflect.TypeOf((*Stream)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "applyRemovalPolicy", GoMethod: "ApplyRemovalPolicy"},
			_jsii_.MemberProperty{JsiiProperty: "encryptionKey", GoGetter: "EncryptionKey"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberMethod{JsiiMethod: "generatePhysicalName", GoMethod: "GeneratePhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceArnAttribute", GoMethod: "GetResourceArnAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getResourceNameAttribute", GoMethod: "GetResourceNameAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "grant", GoMethod: "Grant"},
			_jsii_.MemberMethod{JsiiMethod: "grantRead", GoMethod: "GrantRead"},
			_jsii_.MemberMethod{JsiiMethod: "grantReadWrite", GoMethod: "GrantReadWrite"},
			_jsii_.MemberMethod{JsiiMethod: "grantWrite", GoMethod: "GrantWrite"},
			_jsii_.MemberMethod{JsiiMethod: "metric", GoMethod: "Metric"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecords", GoMethod: "MetricGetRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsBytes", GoMethod: "MetricGetRecordsBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsIteratorAgeMilliseconds", GoMethod: "MetricGetRecordsIteratorAgeMilliseconds"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsLatency", GoMethod: "MetricGetRecordsLatency"},
			_jsii_.MemberMethod{JsiiMethod: "metricGetRecordsSuccess", GoMethod: "MetricGetRecordsSuccess"},
			_jsii_.MemberMethod{JsiiMethod: "metricIncomingBytes", GoMethod: "MetricIncomingBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricIncomingRecords", GoMethod: "MetricIncomingRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordBytes", GoMethod: "MetricPutRecordBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordLatency", GoMethod: "MetricPutRecordLatency"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsBytes", GoMethod: "MetricPutRecordsBytes"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsFailedRecords", GoMethod: "MetricPutRecordsFailedRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsLatency", GoMethod: "MetricPutRecordsLatency"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsSuccess", GoMethod: "MetricPutRecordsSuccess"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsSuccessfulRecords", GoMethod: "MetricPutRecordsSuccessfulRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsThrottledRecords", GoMethod: "MetricPutRecordsThrottledRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordsTotalRecords", GoMethod: "MetricPutRecordsTotalRecords"},
			_jsii_.MemberMethod{JsiiMethod: "metricPutRecordSuccess", GoMethod: "MetricPutRecordSuccess"},
			_jsii_.MemberMethod{JsiiMethod: "metricReadProvisionedThroughputExceeded", GoMethod: "MetricReadProvisionedThroughputExceeded"},
			_jsii_.MemberMethod{JsiiMethod: "metricWriteProvisionedThroughputExceeded", GoMethod: "MetricWriteProvisionedThroughputExceeded"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "onPrepare", GoMethod: "OnPrepare"},
			_jsii_.MemberMethod{JsiiMethod: "onSynthesize", GoMethod: "OnSynthesize"},
			_jsii_.MemberMethod{JsiiMethod: "onValidate", GoMethod: "OnValidate"},
			_jsii_.MemberProperty{JsiiProperty: "physicalName", GoGetter: "PhysicalName"},
			_jsii_.MemberMethod{JsiiMethod: "prepare", GoMethod: "Prepare"},
			_jsii_.MemberProperty{JsiiProperty: "stack", GoGetter: "Stack"},
			_jsii_.MemberProperty{JsiiProperty: "streamArn", GoGetter: "StreamArn"},
			_jsii_.MemberProperty{JsiiProperty: "streamName", GoGetter: "StreamName"},
			_jsii_.MemberMethod{JsiiMethod: "synthesize", GoMethod: "Synthesize"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "validate", GoMethod: "Validate"},
		},
		func() interface{} {
			j := jsiiProxy_Stream{}
			_jsii_.InitJsiiProxy(&j.Type__awscdkResource)
			_jsii_.InitJsiiProxy(&j.jsiiProxy_IStream)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_kinesis.StreamAttributes",
		reflect.TypeOf((*StreamAttributes)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"monocdk.aws_kinesis.StreamEncryption",
		reflect.TypeOf((*StreamEncryption)(nil)).Elem(),
		map[string]interface{}{
			"UNENCRYPTED": StreamEncryption_UNENCRYPTED,
			"KMS": StreamEncryption_KMS,
			"MANAGED": StreamEncryption_MANAGED,
		},
	)
	_jsii_.RegisterStruct(
		"monocdk.aws_kinesis.StreamProps",
		reflect.TypeOf((*StreamProps)(nil)).Elem(),
	)
}

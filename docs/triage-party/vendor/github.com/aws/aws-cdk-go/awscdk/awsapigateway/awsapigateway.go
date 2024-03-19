package awsapigateway

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigateway/internal"
	"github.com/aws/aws-cdk-go/awscdk/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/constructs-go/constructs/v3"
)

// Options when binding a log destination to a RestApi Stage.
// Experimental.
type AccessLogDestinationConfig struct {
	// The Amazon Resource Name (ARN) of the destination resource.
	// Experimental.
	DestinationArn *string `json:"destinationArn"`
}

// $context variables that can be used to customize access log pattern.
// Experimental.
type AccessLogField interface {
}

// The jsii proxy struct for AccessLogField
type jsiiProxy_AccessLogField struct {
	_ byte // padding
}

// Experimental.
func NewAccessLogField() AccessLogField {
	_init_.Initialize()

	j := jsiiProxy_AccessLogField{}

	_jsii_.Create(
		"monocdk.aws_apigateway.AccessLogField",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewAccessLogField_Override(a AccessLogField) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.AccessLogField",
		nil, // no parameters
		a,
	)
}

// The API owner's AWS account ID.
// Experimental.
func AccessLogField_ContextAccountId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextAccountId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The identifier API Gateway assigns to your API.
// Experimental.
func AccessLogField_ContextApiId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextApiId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The stringified value of the specified key-value pair of the `context` map returned from an API Gateway Lambda authorizer function.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-use-lambda-authorizer.html
//
// Experimental.
func AccessLogField_ContextAuthorizer(property *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextAuthorizer",
		[]interface{}{property},
		&returns,
	)

	return returns
}

// A property of the claims returned from the Amazon Cognito user pool after the method caller is successfully authenticated.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-integrate-with-cognito.html
//
// Experimental.
func AccessLogField_ContextAuthorizerClaims(property *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextAuthorizerClaims",
		[]interface{}{property},
		&returns,
	)

	return returns
}

// The authorizer latency in ms.
// Experimental.
func AccessLogField_ContextAuthorizerIntegrationLatency() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextAuthorizerIntegrationLatency",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The principal user identification associated with the token sent by the client and returned from an API Gateway Lambda authorizer (formerly known as a custom authorizer).
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-use-lambda-authorizer.html
//
// Experimental.
func AccessLogField_ContextAuthorizerPrincipalId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextAuthorizerPrincipalId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The AWS endpoint's request ID.
// Experimental.
func AccessLogField_ContextAwsEndpointRequestId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextAwsEndpointRequestId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The full domain name used to invoke the API.
//
// This should be the same as the incoming `Host` header.
// Experimental.
func AccessLogField_ContextDomainName() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextDomainName",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The first label of the `$context.domainName`. This is often used as a caller/customer identifier.
// Experimental.
func AccessLogField_ContextDomainPrefix() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextDomainPrefix",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A string containing an API Gateway error message.
// Experimental.
func AccessLogField_ContextErrorMessage() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextErrorMessage",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The quoted value of $context.error.message, namely "$context.error.message".
// Experimental.
func AccessLogField_ContextErrorMessageString() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextErrorMessageString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A type of GatewayResponse.
//
// This variable can only be used for simple variable substitution in a GatewayResponse body-mapping template,
// which is not processed by the Velocity Template Language engine, and in access logging.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/customize-gateway-responses.html
//
// Experimental.
func AccessLogField_ContextErrorResponseType() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextErrorResponseType",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A string containing a detailed validation error message.
// Experimental.
func AccessLogField_ContextErrorValidationErrorString() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextErrorValidationErrorString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The extended ID that API Gateway assigns to the API request, which contains more useful information for debugging/troubleshooting.
// Experimental.
func AccessLogField_ContextExtendedRequestId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextExtendedRequestId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The HTTP method used.
//
// Valid values include: `DELETE`, `GET`, `HEAD`, `OPTIONS`, `PATCH`, `POST`, and `PUT`.
// Experimental.
func AccessLogField_ContextHttpMethod() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextHttpMethod",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The AWS account ID associated with the request.
// Experimental.
func AccessLogField_ContextIdentityAccountId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityAccountId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// For API methods that require an API key, this variable is the API key associated with the method request.
//
// For methods that don't require an API key, this variable is
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-api-usage-plans.html
//
// Experimental.
func AccessLogField_ContextIdentityApiKey() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityApiKey",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The API key ID associated with an API request that requires an API key.
// Experimental.
func AccessLogField_ContextIdentityApiKeyId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityApiKeyId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The principal identifier of the caller making the request.
// Experimental.
func AccessLogField_ContextIdentityCaller() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityCaller",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The Amazon Cognito authentication provider used by the caller making the request.
//
// Available only if the request was signed with Amazon Cognito credentials.
// See: https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-identity.html
//
// Experimental.
func AccessLogField_ContextIdentityCognitoAuthenticationProvider() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityCognitoAuthenticationProvider",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The Amazon Cognito authentication type of the caller making the request.
//
// Available only if the request was signed with Amazon Cognito credentials.
// Experimental.
func AccessLogField_ContextIdentityCognitoAuthenticationType() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityCognitoAuthenticationType",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The Amazon Cognito identity ID of the caller making the request.
//
// Available only if the request was signed with Amazon Cognito credentials.
// Experimental.
func AccessLogField_ContextIdentityCognitoIdentityId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityCognitoIdentityId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The Amazon Cognito identity pool ID of the caller making the request.
//
// Available only if the request was signed with Amazon Cognito credentials.
// Experimental.
func AccessLogField_ContextIdentityCognitoIdentityPoolId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityCognitoIdentityPoolId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The AWS organization ID.
// Experimental.
func AccessLogField_ContextIdentityPrincipalOrgId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityPrincipalOrgId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The source IP address of the TCP connection making the request to API Gateway.
//
// Warning: You should not trust this value if there is any chance that the `X-Forwarded-For` header could be forged.
// Experimental.
func AccessLogField_ContextIdentitySourceIp() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentitySourceIp",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The principal identifier of the user making the request.
//
// Used in Lambda authorizers.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html
//
// Experimental.
func AccessLogField_ContextIdentityUser() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityUser",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The User-Agent header of the API caller.
// Experimental.
func AccessLogField_ContextIdentityUserAgent() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityUserAgent",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The Amazon Resource Name (ARN) of the effective user identified after authentication.
// See: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users.html
//
// Experimental.
func AccessLogField_ContextIdentityUserArn() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIdentityUserArn",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The integration latency in ms.
// Experimental.
func AccessLogField_ContextIntegrationLatency() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIntegrationLatency",
		nil, // no parameters
		&returns,
	)

	return returns
}

// For Lambda proxy integration, this parameter represents the status code returned from AWS Lambda, not from the backend Lambda function.
// Experimental.
func AccessLogField_ContextIntegrationStatus() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextIntegrationStatus",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The request path.
//
// For example, for a non-proxy request URL of https://{rest-api-id.execute-api.{region}.amazonaws.com/{stage}/root/child,
// this value is /{stage}/root/child.
// Experimental.
func AccessLogField_ContextPath() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextPath",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The request protocol, for example, HTTP/1.1.
// Experimental.
func AccessLogField_ContextProtocol() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextProtocol",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The ID that API Gateway assigns to the API request.
// Experimental.
func AccessLogField_ContextRequestId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextRequestId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The request header override.
//
// If this parameter is defined, it contains the headers to be used instead of the HTTP Headers that are defined in the Integration Request pane.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-override-request-response-parameters.html
//
// Experimental.
func AccessLogField_ContextRequestOverrideHeader(headerName *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextRequestOverrideHeader",
		[]interface{}{headerName},
		&returns,
	)

	return returns
}

// The request path override.
//
// If this parameter is defined,
// it contains the request path to be used instead of the URL Path Parameters that are defined in the Integration Request pane.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-override-request-response-parameters.html
//
// Experimental.
func AccessLogField_ContextRequestOverridePath(pathName *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextRequestOverridePath",
		[]interface{}{pathName},
		&returns,
	)

	return returns
}

// The request query string override.
//
// If this parameter is defined, it contains the request query strings to be used instead
// of the URL Query String Parameters that are defined in the Integration Request pane.
// Experimental.
func AccessLogField_ContextRequestOverrideQuerystring(querystringName *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextRequestOverrideQuerystring",
		[]interface{}{querystringName},
		&returns,
	)

	return returns
}

// The CLF-formatted request time (dd/MMM/yyyy:HH:mm:ss +-hhmm).
// Experimental.
func AccessLogField_ContextRequestTime() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextRequestTime",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The Epoch-formatted request time.
// Experimental.
func AccessLogField_ContextRequestTimeEpoch() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextRequestTimeEpoch",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The identifier that API Gateway assigns to your resource.
// Experimental.
func AccessLogField_ContextResourceId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextResourceId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The path to your resource.
//
// For example, for the non-proxy request URI of `https://{rest-api-id.execute-api.{region}.amazonaws.com/{stage}/root/child`,
// The $context.resourcePath value is `/root/child`.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-create-api-step-by-step.html
//
// Experimental.
func AccessLogField_ContextResourcePath() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextResourcePath",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The response latency in ms.
// Experimental.
func AccessLogField_ContextResponseLatency() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextResponseLatency",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The response payload length.
// Experimental.
func AccessLogField_ContextResponseLength() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextResponseLength",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The response header override.
//
// If this parameter is defined, it contains the header to be returned instead of the Response header
// that is defined as the Default mapping in the Integration Response pane.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-override-request-response-parameters.html
//
// Experimental.
func AccessLogField_ContextResponseOverrideHeader(headerName *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextResponseOverrideHeader",
		[]interface{}{headerName},
		&returns,
	)

	return returns
}

// The response status code override.
//
// If this parameter is defined, it contains the status code to be returned instead of the Method response status
// that is defined as the Default mapping in the Integration Response pane.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-override-request-response-parameters.html
//
// Experimental.
func AccessLogField_ContextResponseOverrideStatus() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextResponseOverrideStatus",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The deployment stage of the API request (for example, `Beta` or `Prod`).
// Experimental.
func AccessLogField_ContextStage() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextStage",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The method response status.
// Experimental.
func AccessLogField_ContextStatus() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextStatus",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The response received from AWS WAF: `WAF_ALLOW` or `WAF_BLOCK`.
//
// Will not be set if the stage is not associated with a web ACL.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-control-access-aws-waf.html
//
// Experimental.
func AccessLogField_ContextWafResponseCode() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextWafResponseCode",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The complete ARN of the web ACL that is used to decide whether to allow or block the request.
//
// Will not be set if the stage is not associated with a web ACL.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-control-access-aws-waf.html
//
// Experimental.
func AccessLogField_ContextWebaclArn() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextWebaclArn",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The trace ID for the X-Ray trace.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-enabling-xray.html
//
// Experimental.
func AccessLogField_ContextXrayTraceId() *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogField",
		"contextXrayTraceId",
		nil, // no parameters
		&returns,
	)

	return returns
}

// factory methods for access log format.
// Experimental.
type AccessLogFormat interface {
	ToString() *string
}

// The jsii proxy struct for AccessLogFormat
type jsiiProxy_AccessLogFormat struct {
	_ byte // padding
}

// Generate Common Log Format.
// Experimental.
func AccessLogFormat_Clf() AccessLogFormat {
	_init_.Initialize()

	var returns AccessLogFormat

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogFormat",
		"clf",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Custom log format.
//
// You can create any log format string. You can easily get the $ context variable by using the methods of AccessLogField.
//
// TODO: EXAMPLE
//
// Experimental.
func AccessLogFormat_Custom(format *string) AccessLogFormat {
	_init_.Initialize()

	var returns AccessLogFormat

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogFormat",
		"custom",
		[]interface{}{format},
		&returns,
	)

	return returns
}

// Access log will be produced in the JSON format with a set of fields most useful in the access log.
//
// All fields are turned on by default with the
// option to turn off specific fields.
// Experimental.
func AccessLogFormat_JsonWithStandardFields(fields *JsonWithStandardFieldProps) AccessLogFormat {
	_init_.Initialize()

	var returns AccessLogFormat

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AccessLogFormat",
		"jsonWithStandardFields",
		[]interface{}{fields},
		&returns,
	)

	return returns
}

// Output a format string to be used with CloudFormation.
// Experimental.
func (a *jsiiProxy_AccessLogFormat) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options to the UsagePlan.addApiKey() method.
// Experimental.
type AddApiKeyOptions struct {
	// Override the CloudFormation logical id of the AWS::ApiGateway::UsagePlanKey resource.
	// Experimental.
	OverrideLogicalId *string `json:"overrideLogicalId"`
}

// Represents an OpenAPI definition asset.
// Experimental.
type ApiDefinition interface {
	Bind(scope awscdk.Construct) *ApiDefinitionConfig
}

// The jsii proxy struct for ApiDefinition
type jsiiProxy_ApiDefinition struct {
	_ byte // padding
}

// Experimental.
func NewApiDefinition_Override(a ApiDefinition) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.ApiDefinition",
		nil, // no parameters
		a,
	)
}

// Loads the API specification from a local disk asset.
// Experimental.
func ApiDefinition_FromAsset(file *string, options *awss3assets.AssetOptions) AssetApiDefinition {
	_init_.Initialize()

	var returns AssetApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ApiDefinition",
		"fromAsset",
		[]interface{}{file, options},
		&returns,
	)

	return returns
}

// Creates an API definition from a specification file in an S3 bucket.
// Experimental.
func ApiDefinition_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3ApiDefinition {
	_init_.Initialize()

	var returns S3ApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ApiDefinition",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Create an API definition from an inline object.
//
// The inline object must follow the
// schema of OpenAPI 2.0 or OpenAPI 3.0
//
// TODO: EXAMPLE
//
// Experimental.
func ApiDefinition_FromInline(definition interface{}) InlineApiDefinition {
	_init_.Initialize()

	var returns InlineApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ApiDefinition",
		"fromInline",
		[]interface{}{definition},
		&returns,
	)

	return returns
}

// Called when the specification is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (a *jsiiProxy_ApiDefinition) Bind(scope awscdk.Construct) *ApiDefinitionConfig {
	var returns *ApiDefinitionConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Post-Binding Configuration for a CDK construct.
// Experimental.
type ApiDefinitionConfig struct {
	// Inline specification (mutually exclusive with `s3Location`).
	// Experimental.
	InlineDefinition interface{} `json:"inlineDefinition"`
	// The location of the specification in S3 (mutually exclusive with `inlineDefinition`).
	// Experimental.
	S3Location *ApiDefinitionS3Location `json:"s3Location"`
}

// S3 location of the API definition file.
// Experimental.
type ApiDefinitionS3Location struct {
	// The S3 bucket.
	// Experimental.
	Bucket *string `json:"bucket"`
	// The S3 key.
	// Experimental.
	Key *string `json:"key"`
	// An optional version.
	// Experimental.
	Version *string `json:"version"`
}

// An API Gateway ApiKey.
//
// An ApiKey can be distributed to API clients that are executing requests
// for Method resources that require an Api Key.
// Experimental.
type ApiKey interface {
	awscdk.Resource
	IApiKey
	Env() *awscdk.ResourceEnvironment
	KeyArn() *string
	KeyId() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant
	GrantWrite(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ApiKey
type jsiiProxy_ApiKey struct {
	internal.Type__awscdkResource
	jsiiProxy_IApiKey
}

func (j *jsiiProxy_ApiKey) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApiKey) KeyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApiKey) KeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApiKey) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApiKey) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApiKey) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewApiKey(scope constructs.Construct, id *string, props *ApiKeyProps) ApiKey {
	_init_.Initialize()

	j := jsiiProxy_ApiKey{}

	_jsii_.Create(
		"monocdk.aws_apigateway.ApiKey",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewApiKey_Override(a ApiKey, scope constructs.Construct, id *string, props *ApiKeyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.ApiKey",
		[]interface{}{scope, id, props},
		a,
	)
}

// Import an ApiKey by its Id.
// Experimental.
func ApiKey_FromApiKeyId(scope constructs.Construct, id *string, apiKeyId *string) IApiKey {
	_init_.Initialize()

	var returns IApiKey

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ApiKey",
		"fromApiKeyId",
		[]interface{}{scope, id, apiKeyId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ApiKey_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ApiKey",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ApiKey_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ApiKey",
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
func (a *jsiiProxy_ApiKey) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_ApiKey) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_ApiKey) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_ApiKey) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Permits the IAM principal all read operations through this key.
// Experimental.
func (a *jsiiProxy_ApiKey) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Permits the IAM principal all read and write operations through this key.
// Experimental.
func (a *jsiiProxy_ApiKey) GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grantReadWrite",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Permits the IAM principal all write operations through this key.
// Experimental.
func (a *jsiiProxy_ApiKey) GrantWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		a,
		"grantWrite",
		[]interface{}{grantee},
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
func (a *jsiiProxy_ApiKey) OnPrepare() {
	_jsii_.InvokeVoid(
		a,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_ApiKey) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
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
func (a *jsiiProxy_ApiKey) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_ApiKey) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_ApiKey) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_ApiKey) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_ApiKey) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// The options for creating an API Key.
// Experimental.
type ApiKeyOptions struct {
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// A name for the API key.
	//
	// If you don't specify a name, AWS CloudFormation generates a unique physical ID and uses that ID for the API key name.
	// Experimental.
	ApiKeyName *string `json:"apiKeyName"`
	// The value of the API key.
	//
	// Must be at least 20 characters long.
	// Experimental.
	Value *string `json:"value"`
}

// ApiKey Properties.
// Experimental.
type ApiKeyProps struct {
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// A name for the API key.
	//
	// If you don't specify a name, AWS CloudFormation generates a unique physical ID and uses that ID for the API key name.
	// Experimental.
	ApiKeyName *string `json:"apiKeyName"`
	// The value of the API key.
	//
	// Must be at least 20 characters long.
	// Experimental.
	Value *string `json:"value"`
	// An AWS Marketplace customer identifier to use when integrating with the AWS SaaS Marketplace.
	// Experimental.
	CustomerId *string `json:"customerId"`
	// A description of the purpose of the API key.
	// Experimental.
	Description *string `json:"description"`
	// Indicates whether the API key can be used by clients.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// Specifies whether the key identifier is distinct from the created API key value.
	// Experimental.
	GenerateDistinctId *bool `json:"generateDistinctId"`
	// A list of resources this api key is associated with.
	// Experimental.
	Resources *[]IRestApi `json:"resources"`
}

// Experimental.
type ApiKeySourceType string

const (
	ApiKeySourceType_HEADER ApiKeySourceType = "HEADER"
	ApiKeySourceType_AUTHORIZER ApiKeySourceType = "AUTHORIZER"
)

// OpenAPI specification from a local file.
// Experimental.
type AssetApiDefinition interface {
	ApiDefinition
	Bind(scope awscdk.Construct) *ApiDefinitionConfig
}

// The jsii proxy struct for AssetApiDefinition
type jsiiProxy_AssetApiDefinition struct {
	jsiiProxy_ApiDefinition
}

// Experimental.
func NewAssetApiDefinition(path *string, options *awss3assets.AssetOptions) AssetApiDefinition {
	_init_.Initialize()

	j := jsiiProxy_AssetApiDefinition{}

	_jsii_.Create(
		"monocdk.aws_apigateway.AssetApiDefinition",
		[]interface{}{path, options},
		&j,
	)

	return &j
}

// Experimental.
func NewAssetApiDefinition_Override(a AssetApiDefinition, path *string, options *awss3assets.AssetOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.AssetApiDefinition",
		[]interface{}{path, options},
		a,
	)
}

// Loads the API specification from a local disk asset.
// Experimental.
func AssetApiDefinition_FromAsset(file *string, options *awss3assets.AssetOptions) AssetApiDefinition {
	_init_.Initialize()

	var returns AssetApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AssetApiDefinition",
		"fromAsset",
		[]interface{}{file, options},
		&returns,
	)

	return returns
}

// Creates an API definition from a specification file in an S3 bucket.
// Experimental.
func AssetApiDefinition_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3ApiDefinition {
	_init_.Initialize()

	var returns S3ApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AssetApiDefinition",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Create an API definition from an inline object.
//
// The inline object must follow the
// schema of OpenAPI 2.0 or OpenAPI 3.0
//
// TODO: EXAMPLE
//
// Experimental.
func AssetApiDefinition_FromInline(definition interface{}) InlineApiDefinition {
	_init_.Initialize()

	var returns InlineApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.AssetApiDefinition",
		"fromInline",
		[]interface{}{definition},
		&returns,
	)

	return returns
}

// Called when the specification is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (a *jsiiProxy_AssetApiDefinition) Bind(scope awscdk.Construct) *ApiDefinitionConfig {
	var returns *ApiDefinitionConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{scope},
		&returns,
	)

	return returns
}

// Experimental.
type AuthorizationType string

const (
	AuthorizationType_NONE AuthorizationType = "NONE"
	AuthorizationType_IAM AuthorizationType = "IAM"
	AuthorizationType_CUSTOM AuthorizationType = "CUSTOM"
	AuthorizationType_COGNITO AuthorizationType = "COGNITO"
)

// Base class for all custom authorizers.
// Experimental.
type Authorizer interface {
	awscdk.Resource
	IAuthorizer
	AuthorizationType() AuthorizationType
	AuthorizerId() *string
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
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

// The jsii proxy struct for Authorizer
type jsiiProxy_Authorizer struct {
	internal.Type__awscdkResource
	jsiiProxy_IAuthorizer
}

func (j *jsiiProxy_Authorizer) AuthorizationType() AuthorizationType {
	var returns AuthorizationType
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Authorizer) AuthorizerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Authorizer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Authorizer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Authorizer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Authorizer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewAuthorizer_Override(a Authorizer, scope constructs.Construct, id *string, props *awscdk.ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.Authorizer",
		[]interface{}{scope, id, props},
		a,
	)
}

// Return whether the given object is an Authorizer.
// Experimental.
func Authorizer_IsAuthorizer(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Authorizer",
		"isAuthorizer",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Authorizer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Authorizer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Authorizer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Authorizer",
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
func (a *jsiiProxy_Authorizer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_Authorizer) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Authorizer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Authorizer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Authorizer) OnPrepare() {
	_jsii_.InvokeVoid(
		a,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_Authorizer) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
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
func (a *jsiiProxy_Authorizer) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Authorizer) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_Authorizer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_Authorizer) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
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
func (a *jsiiProxy_Authorizer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// This type of integration lets an API expose AWS service actions.
//
// It is
// intended for calling all AWS service actions, but is not recommended for
// calling a Lambda function, because the Lambda custom integration is a legacy
// technology.
// Experimental.
type AwsIntegration interface {
	Integration
	Bind(method Method) *IntegrationConfig
}

// The jsii proxy struct for AwsIntegration
type jsiiProxy_AwsIntegration struct {
	jsiiProxy_Integration
}

// Experimental.
func NewAwsIntegration(props *AwsIntegrationProps) AwsIntegration {
	_init_.Initialize()

	j := jsiiProxy_AwsIntegration{}

	_jsii_.Create(
		"monocdk.aws_apigateway.AwsIntegration",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewAwsIntegration_Override(a AwsIntegration, props *AwsIntegrationProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.AwsIntegration",
		[]interface{}{props},
		a,
	)
}

// Can be overridden by subclasses to allow the integration to interact with the method being integrated, access the REST API object, method ARNs, etc.
// Experimental.
func (a *jsiiProxy_AwsIntegration) Bind(method Method) *IntegrationConfig {
	var returns *IntegrationConfig

	_jsii_.Invoke(
		a,
		"bind",
		[]interface{}{method},
		&returns,
	)

	return returns
}

// Experimental.
type AwsIntegrationProps struct {
	// The name of the integrated AWS service (e.g. `s3`).
	// Experimental.
	Service *string `json:"service"`
	// The AWS action to perform in the integration.
	//
	// Use `actionParams` to specify key-value params for the action.
	//
	// Mutually exclusive with `path`.
	// Experimental.
	Action *string `json:"action"`
	// Parameters for the action.
	//
	// `action` must be set, and `path` must be undefined.
	// The action params will be URL encoded.
	// Experimental.
	ActionParameters *map[string]*string `json:"actionParameters"`
	// The integration's HTTP method type.
	// Experimental.
	IntegrationHttpMethod *string `json:"integrationHttpMethod"`
	// Integration options, such as content handling, request/response mapping, etc.
	// Experimental.
	Options *IntegrationOptions `json:"options"`
	// The path to use for path-base APIs.
	//
	// For example, for S3 GET, you can set path to `bucket/key`.
	// For lambda, you can set path to `2015-03-31/functions/${function-arn}/invocations`
	//
	// Mutually exclusive with the `action` options.
	// Experimental.
	Path *string `json:"path"`
	// Use AWS_PROXY integration.
	// Experimental.
	Proxy *bool `json:"proxy"`
	// The region of the integrated AWS service.
	// Experimental.
	Region *string `json:"region"`
	// A designated subdomain supported by certain AWS service for fast host-name lookup.
	// Experimental.
	Subdomain *string `json:"subdomain"`
}

// This resource creates a base path that clients who call your API must use in the invocation URL.
//
// Unless you're importing a domain with `DomainName.fromDomainNameAttributes()`,
// you can use `DomainName.addBasePathMapping()` to define mappings.
// Experimental.
type BasePathMapping interface {
	awscdk.Resource
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
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

// The jsii proxy struct for BasePathMapping
type jsiiProxy_BasePathMapping struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_BasePathMapping) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BasePathMapping) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BasePathMapping) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BasePathMapping) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewBasePathMapping(scope constructs.Construct, id *string, props *BasePathMappingProps) BasePathMapping {
	_init_.Initialize()

	j := jsiiProxy_BasePathMapping{}

	_jsii_.Create(
		"monocdk.aws_apigateway.BasePathMapping",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewBasePathMapping_Override(b BasePathMapping, scope constructs.Construct, id *string, props *BasePathMappingProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.BasePathMapping",
		[]interface{}{scope, id, props},
		b,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BasePathMapping_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.BasePathMapping",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func BasePathMapping_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.BasePathMapping",
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
func (b *jsiiProxy_BasePathMapping) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (b *jsiiProxy_BasePathMapping) GeneratePhysicalName() *string {
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
func (b *jsiiProxy_BasePathMapping) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (b *jsiiProxy_BasePathMapping) GetResourceNameAttribute(nameAttr *string) *string {
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
func (b *jsiiProxy_BasePathMapping) OnPrepare() {
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
func (b *jsiiProxy_BasePathMapping) OnSynthesize(session constructs.ISynthesisSession) {
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
func (b *jsiiProxy_BasePathMapping) OnValidate() *[]*string {
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
func (b *jsiiProxy_BasePathMapping) Prepare() {
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
func (b *jsiiProxy_BasePathMapping) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BasePathMapping) ToString() *string {
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
func (b *jsiiProxy_BasePathMapping) Validate() *[]*string {
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
type BasePathMappingOptions struct {
	// The base path name that callers of the API must provide in the URL after the domain name (e.g. `example.com/base-path`). If you specify this property, it can't be an empty string.
	// Experimental.
	BasePath *string `json:"basePath"`
	// The Deployment stage of API [disable-awslint:ref-via-interface].
	// Experimental.
	Stage Stage `json:"stage"`
}

// Experimental.
type BasePathMappingProps struct {
	// The base path name that callers of the API must provide in the URL after the domain name (e.g. `example.com/base-path`). If you specify this property, it can't be an empty string.
	// Experimental.
	BasePath *string `json:"basePath"`
	// The Deployment stage of API [disable-awslint:ref-via-interface].
	// Experimental.
	Stage Stage `json:"stage"`
	// The DomainName to associate with this base path mapping.
	// Experimental.
	DomainName IDomainName `json:"domainName"`
	// The RestApi resource to target.
	// Experimental.
	RestApi IRestApi `json:"restApi"`
}

// A CloudFormation `AWS::ApiGateway::Account`.
type CfnAccount interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CloudWatchRoleArn() *string
	SetCloudWatchRoleArn(val *string)
	CreationStack() *[]*string
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

// The jsii proxy struct for CfnAccount
type jsiiProxy_CfnAccount struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAccount) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) CloudWatchRoleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cloudWatchRoleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAccount) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::Account`.
func NewCfnAccount(scope awscdk.Construct, id *string, props *CfnAccountProps) CfnAccount {
	_init_.Initialize()

	j := jsiiProxy_CfnAccount{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnAccount",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::Account`.
func NewCfnAccount_Override(c CfnAccount, scope awscdk.Construct, id *string, props *CfnAccountProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnAccount",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAccount) SetCloudWatchRoleArn(val *string) {
	_jsii_.Set(
		j,
		"cloudWatchRoleArn",
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
func CfnAccount_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAccount",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAccount_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAccount",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAccount_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAccount",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAccount_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnAccount",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAccount) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAccount) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAccount) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAccount) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAccount) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAccount) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAccount) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAccount) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAccount) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAccount) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAccount) OnPrepare() {
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
func (c *jsiiProxy_CfnAccount) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAccount) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAccount) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAccount) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAccount) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAccount) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAccount) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAccount) ToString() *string {
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
func (c *jsiiProxy_CfnAccount) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAccount) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::Account`.
type CfnAccountProps struct {
	// `AWS::ApiGateway::Account.CloudWatchRoleArn`.
	CloudWatchRoleArn *string `json:"cloudWatchRoleArn"`
}

// A CloudFormation `AWS::ApiGateway::ApiKey`.
type CfnApiKey interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrApiKeyId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	CustomerId() *string
	SetCustomerId(val *string)
	Description() *string
	SetDescription(val *string)
	Enabled() interface{}
	SetEnabled(val interface{})
	GenerateDistinctId() interface{}
	SetGenerateDistinctId(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	StageKeys() interface{}
	SetStageKeys(val interface{})
	Tags() awscdk.TagManager
	UpdatedProperites() *map[string]interface{}
	Value() *string
	SetValue(val *string)
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

// The jsii proxy struct for CfnApiKey
type jsiiProxy_CfnApiKey struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnApiKey) AttrApiKeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrApiKeyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) CustomerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"customerId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Enabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"enabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) GenerateDistinctId() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"generateDistinctId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) StageKeys() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stageKeys",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiKey) Value() *string {
	var returns *string
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::ApiKey`.
func NewCfnApiKey(scope awscdk.Construct, id *string, props *CfnApiKeyProps) CfnApiKey {
	_init_.Initialize()

	j := jsiiProxy_CfnApiKey{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnApiKey",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::ApiKey`.
func NewCfnApiKey_Override(c CfnApiKey, scope awscdk.Construct, id *string, props *CfnApiKeyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnApiKey",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnApiKey) SetCustomerId(val *string) {
	_jsii_.Set(
		j,
		"customerId",
		val,
	)
}

func (j *jsiiProxy_CfnApiKey) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnApiKey) SetEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"enabled",
		val,
	)
}

func (j *jsiiProxy_CfnApiKey) SetGenerateDistinctId(val interface{}) {
	_jsii_.Set(
		j,
		"generateDistinctId",
		val,
	)
}

func (j *jsiiProxy_CfnApiKey) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnApiKey) SetStageKeys(val interface{}) {
	_jsii_.Set(
		j,
		"stageKeys",
		val,
	)
}

func (j *jsiiProxy_CfnApiKey) SetValue(val *string) {
	_jsii_.Set(
		j,
		"value",
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
func CfnApiKey_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiKey",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnApiKey_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiKey",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnApiKey_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiKey",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnApiKey_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnApiKey",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnApiKey) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnApiKey) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnApiKey) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnApiKey) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnApiKey) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnApiKey) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnApiKey) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnApiKey) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnApiKey) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnApiKey) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnApiKey) OnPrepare() {
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
func (c *jsiiProxy_CfnApiKey) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnApiKey) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnApiKey) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnApiKey) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnApiKey) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnApiKey) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnApiKey) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnApiKey) ToString() *string {
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
func (c *jsiiProxy_CfnApiKey) Validate() *[]*string {
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
func (c *jsiiProxy_CfnApiKey) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnApiKey_StageKeyProperty struct {
	// `CfnApiKey.StageKeyProperty.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `CfnApiKey.StageKeyProperty.StageName`.
	StageName *string `json:"stageName"`
}

// Properties for defining a `AWS::ApiGateway::ApiKey`.
type CfnApiKeyProps struct {
	// `AWS::ApiGateway::ApiKey.CustomerId`.
	CustomerId *string `json:"customerId"`
	// `AWS::ApiGateway::ApiKey.Description`.
	Description *string `json:"description"`
	// `AWS::ApiGateway::ApiKey.Enabled`.
	Enabled interface{} `json:"enabled"`
	// `AWS::ApiGateway::ApiKey.GenerateDistinctId`.
	GenerateDistinctId interface{} `json:"generateDistinctId"`
	// `AWS::ApiGateway::ApiKey.Name`.
	Name *string `json:"name"`
	// `AWS::ApiGateway::ApiKey.StageKeys`.
	StageKeys interface{} `json:"stageKeys"`
	// `AWS::ApiGateway::ApiKey.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::ApiGateway::ApiKey.Value`.
	Value *string `json:"value"`
}

// A CloudFormation `AWS::ApiGatewayV2::ApiMapping`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-apimapping.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnApiMappingV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	ApiMappingKey() *string
	SetApiMappingKey(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DomainName() *string
	SetDomainName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	Stage() *string
	SetStage(val *string)
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

// The jsii proxy struct for CfnApiMappingV2
type jsiiProxy_CfnApiMappingV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnApiMappingV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) ApiMappingKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiMappingKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) Stage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiMappingV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::ApiMapping`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnApiMappingV2(scope awscdk.Construct, id *string, props *CfnApiMappingV2Props) CfnApiMappingV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnApiMappingV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnApiMappingV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::ApiMapping`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnApiMappingV2_Override(c CfnApiMappingV2, scope awscdk.Construct, id *string, props *CfnApiMappingV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnApiMappingV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnApiMappingV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnApiMappingV2) SetApiMappingKey(val *string) {
	_jsii_.Set(
		j,
		"apiMappingKey",
		val,
	)
}

func (j *jsiiProxy_CfnApiMappingV2) SetDomainName(val *string) {
	_jsii_.Set(
		j,
		"domainName",
		val,
	)
}

func (j *jsiiProxy_CfnApiMappingV2) SetStage(val *string) {
	_jsii_.Set(
		j,
		"stage",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnApiMappingV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiMappingV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnApiMappingV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiMappingV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnApiMappingV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiMappingV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnApiMappingV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnApiMappingV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiMappingV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGatewayV2::ApiMapping`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-apimapping.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnApiMappingV2Props struct {
	// `AWS::ApiGatewayV2::ApiMapping.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-apimapping.html#cfn-apigatewayv2-apimapping-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::ApiMapping.DomainName`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-apimapping.html#cfn-apigatewayv2-apimapping-domainname
	//
	// Deprecated: moved to package aws-apigatewayv2
	DomainName *string `json:"domainName"`
	// `AWS::ApiGatewayV2::ApiMapping.Stage`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-apimapping.html#cfn-apigatewayv2-apimapping-stage
	//
	// Deprecated: moved to package aws-apigatewayv2
	Stage *string `json:"stage"`
	// `AWS::ApiGatewayV2::ApiMapping.ApiMappingKey`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-apimapping.html#cfn-apigatewayv2-apimapping-apimappingkey
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiMappingKey *string `json:"apiMappingKey"`
}

// A CloudFormation `AWS::ApiGatewayV2::Api`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnApiV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiKeySelectionExpression() *string
	SetApiKeySelectionExpression(val *string)
	BasePath() *string
	SetBasePath(val *string)
	Body() interface{}
	SetBody(val interface{})
	BodyS3Location() interface{}
	SetBodyS3Location(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CorsConfiguration() interface{}
	SetCorsConfiguration(val interface{})
	CreationStack() *[]*string
	CredentialsArn() *string
	SetCredentialsArn(val *string)
	Description() *string
	SetDescription(val *string)
	DisableSchemaValidation() interface{}
	SetDisableSchemaValidation(val interface{})
	FailOnWarnings() interface{}
	SetFailOnWarnings(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	ProtocolType() *string
	SetProtocolType(val *string)
	Ref() *string
	RouteKey() *string
	SetRouteKey(val *string)
	RouteSelectionExpression() *string
	SetRouteSelectionExpression(val *string)
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	Target() *string
	SetTarget(val *string)
	UpdatedProperites() *map[string]interface{}
	Version() *string
	SetVersion(val *string)
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

// The jsii proxy struct for CfnApiV2
type jsiiProxy_CfnApiV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnApiV2) ApiKeySelectionExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiKeySelectionExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) BasePath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"basePath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Body() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"body",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) BodyS3Location() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"bodyS3Location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) CorsConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"corsConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) CredentialsArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"credentialsArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) DisableSchemaValidation() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"disableSchemaValidation",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) FailOnWarnings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"failOnWarnings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) ProtocolType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"protocolType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) RouteKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"routeKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) RouteSelectionExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"routeSelectionExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Target() *string {
	var returns *string
	_jsii_.Get(
		j,
		"target",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnApiV2) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::Api`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnApiV2(scope awscdk.Construct, id *string, props *CfnApiV2Props) CfnApiV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnApiV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnApiV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::Api`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnApiV2_Override(c CfnApiV2, scope awscdk.Construct, id *string, props *CfnApiV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnApiV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnApiV2) SetApiKeySelectionExpression(val *string) {
	_jsii_.Set(
		j,
		"apiKeySelectionExpression",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetBasePath(val *string) {
	_jsii_.Set(
		j,
		"basePath",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetBody(val interface{}) {
	_jsii_.Set(
		j,
		"body",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetBodyS3Location(val interface{}) {
	_jsii_.Set(
		j,
		"bodyS3Location",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetCorsConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"corsConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetCredentialsArn(val *string) {
	_jsii_.Set(
		j,
		"credentialsArn",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetDisableSchemaValidation(val interface{}) {
	_jsii_.Set(
		j,
		"disableSchemaValidation",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetFailOnWarnings(val interface{}) {
	_jsii_.Set(
		j,
		"failOnWarnings",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetProtocolType(val *string) {
	_jsii_.Set(
		j,
		"protocolType",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetRouteKey(val *string) {
	_jsii_.Set(
		j,
		"routeKey",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetRouteSelectionExpression(val *string) {
	_jsii_.Set(
		j,
		"routeSelectionExpression",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetTarget(val *string) {
	_jsii_.Set(
		j,
		"target",
		val,
	)
}

func (j *jsiiProxy_CfnApiV2) SetVersion(val *string) {
	_jsii_.Set(
		j,
		"version",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnApiV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnApiV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnApiV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnApiV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnApiV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnApiV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnApiV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-bodys3location.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnApiV2_BodyS3LocationProperty struct {
	// `CfnApiV2.BodyS3LocationProperty.Bucket`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-bodys3location.html#cfn-apigatewayv2-api-bodys3location-bucket
	//
	// Deprecated: moved to package aws-apigatewayv2
	Bucket *string `json:"bucket"`
	// `CfnApiV2.BodyS3LocationProperty.Etag`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-bodys3location.html#cfn-apigatewayv2-api-bodys3location-etag
	//
	// Deprecated: moved to package aws-apigatewayv2
	Etag *string `json:"etag"`
	// `CfnApiV2.BodyS3LocationProperty.Key`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-bodys3location.html#cfn-apigatewayv2-api-bodys3location-key
	//
	// Deprecated: moved to package aws-apigatewayv2
	Key *string `json:"key"`
	// `CfnApiV2.BodyS3LocationProperty.Version`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-bodys3location.html#cfn-apigatewayv2-api-bodys3location-version
	//
	// Deprecated: moved to package aws-apigatewayv2
	Version *string `json:"version"`
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-cors.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnApiV2_CorsProperty struct {
	// `CfnApiV2.CorsProperty.AllowCredentials`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-cors.html#cfn-apigatewayv2-api-cors-allowcredentials
	//
	// Deprecated: moved to package aws-apigatewayv2
	AllowCredentials interface{} `json:"allowCredentials"`
	// `CfnApiV2.CorsProperty.AllowHeaders`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-cors.html#cfn-apigatewayv2-api-cors-allowheaders
	//
	// Deprecated: moved to package aws-apigatewayv2
	AllowHeaders *[]*string `json:"allowHeaders"`
	// `CfnApiV2.CorsProperty.AllowMethods`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-cors.html#cfn-apigatewayv2-api-cors-allowmethods
	//
	// Deprecated: moved to package aws-apigatewayv2
	AllowMethods *[]*string `json:"allowMethods"`
	// `CfnApiV2.CorsProperty.AllowOrigins`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-cors.html#cfn-apigatewayv2-api-cors-alloworigins
	//
	// Deprecated: moved to package aws-apigatewayv2
	AllowOrigins *[]*string `json:"allowOrigins"`
	// `CfnApiV2.CorsProperty.ExposeHeaders`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-cors.html#cfn-apigatewayv2-api-cors-exposeheaders
	//
	// Deprecated: moved to package aws-apigatewayv2
	ExposeHeaders *[]*string `json:"exposeHeaders"`
	// `CfnApiV2.CorsProperty.MaxAge`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-api-cors.html#cfn-apigatewayv2-api-cors-maxage
	//
	// Deprecated: moved to package aws-apigatewayv2
	MaxAge *float64 `json:"maxAge"`
}

// Properties for defining a `AWS::ApiGatewayV2::Api`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnApiV2Props struct {
	// `AWS::ApiGatewayV2::Api.ApiKeySelectionExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-apikeyselectionexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiKeySelectionExpression *string `json:"apiKeySelectionExpression"`
	// `AWS::ApiGatewayV2::Api.BasePath`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-basepath
	//
	// Deprecated: moved to package aws-apigatewayv2
	BasePath *string `json:"basePath"`
	// `AWS::ApiGatewayV2::Api.Body`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-body
	//
	// Deprecated: moved to package aws-apigatewayv2
	Body interface{} `json:"body"`
	// `AWS::ApiGatewayV2::Api.BodyS3Location`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-bodys3location
	//
	// Deprecated: moved to package aws-apigatewayv2
	BodyS3Location interface{} `json:"bodyS3Location"`
	// `AWS::ApiGatewayV2::Api.CorsConfiguration`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-corsconfiguration
	//
	// Deprecated: moved to package aws-apigatewayv2
	CorsConfiguration interface{} `json:"corsConfiguration"`
	// `AWS::ApiGatewayV2::Api.CredentialsArn`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-credentialsarn
	//
	// Deprecated: moved to package aws-apigatewayv2
	CredentialsArn *string `json:"credentialsArn"`
	// `AWS::ApiGatewayV2::Api.Description`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-description
	//
	// Deprecated: moved to package aws-apigatewayv2
	Description *string `json:"description"`
	// `AWS::ApiGatewayV2::Api.DisableSchemaValidation`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-disableschemavalidation
	//
	// Deprecated: moved to package aws-apigatewayv2
	DisableSchemaValidation interface{} `json:"disableSchemaValidation"`
	// `AWS::ApiGatewayV2::Api.FailOnWarnings`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-failonwarnings
	//
	// Deprecated: moved to package aws-apigatewayv2
	FailOnWarnings interface{} `json:"failOnWarnings"`
	// `AWS::ApiGatewayV2::Api.Name`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-name
	//
	// Deprecated: moved to package aws-apigatewayv2
	Name *string `json:"name"`
	// `AWS::ApiGatewayV2::Api.ProtocolType`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-protocoltype
	//
	// Deprecated: moved to package aws-apigatewayv2
	ProtocolType *string `json:"protocolType"`
	// `AWS::ApiGatewayV2::Api.RouteKey`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-routekey
	//
	// Deprecated: moved to package aws-apigatewayv2
	RouteKey *string `json:"routeKey"`
	// `AWS::ApiGatewayV2::Api.RouteSelectionExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-routeselectionexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	RouteSelectionExpression *string `json:"routeSelectionExpression"`
	// `AWS::ApiGatewayV2::Api.Tags`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-tags
	//
	// Deprecated: moved to package aws-apigatewayv2
	Tags interface{} `json:"tags"`
	// `AWS::ApiGatewayV2::Api.Target`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-target
	//
	// Deprecated: moved to package aws-apigatewayv2
	Target *string `json:"target"`
	// `AWS::ApiGatewayV2::Api.Version`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html#cfn-apigatewayv2-api-version
	//
	// Deprecated: moved to package aws-apigatewayv2
	Version *string `json:"version"`
}

// A CloudFormation `AWS::ApiGateway::Authorizer`.
type CfnAuthorizer interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AuthorizerCredentials() *string
	SetAuthorizerCredentials(val *string)
	AuthorizerResultTtlInSeconds() *float64
	SetAuthorizerResultTtlInSeconds(val *float64)
	AuthorizerUri() *string
	SetAuthorizerUri(val *string)
	AuthType() *string
	SetAuthType(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	IdentitySource() *string
	SetIdentitySource(val *string)
	IdentityValidationExpression() *string
	SetIdentityValidationExpression(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	ProviderArns() *[]*string
	SetProviderArns(val *[]*string)
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Stack() awscdk.Stack
	Type() *string
	SetType(val *string)
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

// The jsii proxy struct for CfnAuthorizer
type jsiiProxy_CfnAuthorizer struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAuthorizer) AuthorizerCredentials() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerCredentials",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) AuthorizerResultTtlInSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"authorizerResultTtlInSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) AuthorizerUri() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerUri",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) AuthType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) IdentitySource() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identitySource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) IdentityValidationExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identityValidationExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) ProviderArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"providerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizer) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::Authorizer`.
func NewCfnAuthorizer(scope awscdk.Construct, id *string, props *CfnAuthorizerProps) CfnAuthorizer {
	_init_.Initialize()

	j := jsiiProxy_CfnAuthorizer{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnAuthorizer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::Authorizer`.
func NewCfnAuthorizer_Override(c CfnAuthorizer, scope awscdk.Construct, id *string, props *CfnAuthorizerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnAuthorizer",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetAuthorizerCredentials(val *string) {
	_jsii_.Set(
		j,
		"authorizerCredentials",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetAuthorizerResultTtlInSeconds(val *float64) {
	_jsii_.Set(
		j,
		"authorizerResultTtlInSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetAuthorizerUri(val *string) {
	_jsii_.Set(
		j,
		"authorizerUri",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetAuthType(val *string) {
	_jsii_.Set(
		j,
		"authType",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetIdentitySource(val *string) {
	_jsii_.Set(
		j,
		"identitySource",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetIdentityValidationExpression(val *string) {
	_jsii_.Set(
		j,
		"identityValidationExpression",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetProviderArns(val *[]*string) {
	_jsii_.Set(
		j,
		"providerArns",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizer) SetType(val *string) {
	_jsii_.Set(
		j,
		"type",
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
func CfnAuthorizer_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAuthorizer",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnAuthorizer_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAuthorizer",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnAuthorizer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAuthorizer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAuthorizer_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnAuthorizer",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnAuthorizer) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnAuthorizer) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnAuthorizer) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnAuthorizer) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnAuthorizer) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnAuthorizer) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnAuthorizer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnAuthorizer) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnAuthorizer) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnAuthorizer) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnAuthorizer) OnPrepare() {
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
func (c *jsiiProxy_CfnAuthorizer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAuthorizer) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnAuthorizer) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnAuthorizer) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnAuthorizer) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnAuthorizer) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnAuthorizer) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnAuthorizer) ToString() *string {
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
func (c *jsiiProxy_CfnAuthorizer) Validate() *[]*string {
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
func (c *jsiiProxy_CfnAuthorizer) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::Authorizer`.
type CfnAuthorizerProps struct {
	// `AWS::ApiGateway::Authorizer.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::Authorizer.Type`.
	Type *string `json:"type"`
	// `AWS::ApiGateway::Authorizer.AuthorizerCredentials`.
	AuthorizerCredentials *string `json:"authorizerCredentials"`
	// `AWS::ApiGateway::Authorizer.AuthorizerResultTtlInSeconds`.
	AuthorizerResultTtlInSeconds *float64 `json:"authorizerResultTtlInSeconds"`
	// `AWS::ApiGateway::Authorizer.AuthorizerUri`.
	AuthorizerUri *string `json:"authorizerUri"`
	// `AWS::ApiGateway::Authorizer.AuthType`.
	AuthType *string `json:"authType"`
	// `AWS::ApiGateway::Authorizer.IdentitySource`.
	IdentitySource *string `json:"identitySource"`
	// `AWS::ApiGateway::Authorizer.IdentityValidationExpression`.
	IdentityValidationExpression *string `json:"identityValidationExpression"`
	// `AWS::ApiGateway::Authorizer.Name`.
	Name *string `json:"name"`
	// `AWS::ApiGateway::Authorizer.ProviderARNs`.
	ProviderArns *[]*string `json:"providerArns"`
}

// A CloudFormation `AWS::ApiGatewayV2::Authorizer`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnAuthorizerV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	AuthorizerCredentialsArn() *string
	SetAuthorizerCredentialsArn(val *string)
	AuthorizerResultTtlInSeconds() *float64
	SetAuthorizerResultTtlInSeconds(val *float64)
	AuthorizerType() *string
	SetAuthorizerType(val *string)
	AuthorizerUri() *string
	SetAuthorizerUri(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	IdentitySource() *[]*string
	SetIdentitySource(val *[]*string)
	IdentityValidationExpression() *string
	SetIdentityValidationExpression(val *string)
	JwtConfiguration() interface{}
	SetJwtConfiguration(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
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

// The jsii proxy struct for CfnAuthorizerV2
type jsiiProxy_CfnAuthorizerV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnAuthorizerV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) AuthorizerCredentialsArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerCredentialsArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) AuthorizerResultTtlInSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"authorizerResultTtlInSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) AuthorizerType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) AuthorizerUri() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerUri",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) IdentitySource() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"identitySource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) IdentityValidationExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"identityValidationExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) JwtConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"jwtConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnAuthorizerV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::Authorizer`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnAuthorizerV2(scope awscdk.Construct, id *string, props *CfnAuthorizerV2Props) CfnAuthorizerV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnAuthorizerV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnAuthorizerV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::Authorizer`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnAuthorizerV2_Override(c CfnAuthorizerV2, scope awscdk.Construct, id *string, props *CfnAuthorizerV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnAuthorizerV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetAuthorizerCredentialsArn(val *string) {
	_jsii_.Set(
		j,
		"authorizerCredentialsArn",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetAuthorizerResultTtlInSeconds(val *float64) {
	_jsii_.Set(
		j,
		"authorizerResultTtlInSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetAuthorizerType(val *string) {
	_jsii_.Set(
		j,
		"authorizerType",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetAuthorizerUri(val *string) {
	_jsii_.Set(
		j,
		"authorizerUri",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetIdentitySource(val *[]*string) {
	_jsii_.Set(
		j,
		"identitySource",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetIdentityValidationExpression(val *string) {
	_jsii_.Set(
		j,
		"identityValidationExpression",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetJwtConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"jwtConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnAuthorizerV2) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnAuthorizerV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAuthorizerV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnAuthorizerV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAuthorizerV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnAuthorizerV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnAuthorizerV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnAuthorizerV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnAuthorizerV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnAuthorizerV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-authorizer-jwtconfiguration.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnAuthorizerV2_JWTConfigurationProperty struct {
	// `CfnAuthorizerV2.JWTConfigurationProperty.Audience`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-authorizer-jwtconfiguration.html#cfn-apigatewayv2-authorizer-jwtconfiguration-audience
	//
	// Deprecated: moved to package aws-apigatewayv2
	Audience *[]*string `json:"audience"`
	// `CfnAuthorizerV2.JWTConfigurationProperty.Issuer`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-authorizer-jwtconfiguration.html#cfn-apigatewayv2-authorizer-jwtconfiguration-issuer
	//
	// Deprecated: moved to package aws-apigatewayv2
	Issuer *string `json:"issuer"`
}

// Properties for defining a `AWS::ApiGatewayV2::Authorizer`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnAuthorizerV2Props struct {
	// `AWS::ApiGatewayV2::Authorizer.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::Authorizer.AuthorizerType`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-authorizertype
	//
	// Deprecated: moved to package aws-apigatewayv2
	AuthorizerType *string `json:"authorizerType"`
	// `AWS::ApiGatewayV2::Authorizer.IdentitySource`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-identitysource
	//
	// Deprecated: moved to package aws-apigatewayv2
	IdentitySource *[]*string `json:"identitySource"`
	// `AWS::ApiGatewayV2::Authorizer.Name`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-name
	//
	// Deprecated: moved to package aws-apigatewayv2
	Name *string `json:"name"`
	// `AWS::ApiGatewayV2::Authorizer.AuthorizerCredentialsArn`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-authorizercredentialsarn
	//
	// Deprecated: moved to package aws-apigatewayv2
	AuthorizerCredentialsArn *string `json:"authorizerCredentialsArn"`
	// `AWS::ApiGatewayV2::Authorizer.AuthorizerResultTtlInSeconds`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-authorizerresultttlinseconds
	//
	// Deprecated: moved to package aws-apigatewayv2
	AuthorizerResultTtlInSeconds *float64 `json:"authorizerResultTtlInSeconds"`
	// `AWS::ApiGatewayV2::Authorizer.AuthorizerUri`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-authorizeruri
	//
	// Deprecated: moved to package aws-apigatewayv2
	AuthorizerUri *string `json:"authorizerUri"`
	// `AWS::ApiGatewayV2::Authorizer.IdentityValidationExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-identityvalidationexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	IdentityValidationExpression *string `json:"identityValidationExpression"`
	// `AWS::ApiGatewayV2::Authorizer.JwtConfiguration`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html#cfn-apigatewayv2-authorizer-jwtconfiguration
	//
	// Deprecated: moved to package aws-apigatewayv2
	JwtConfiguration interface{} `json:"jwtConfiguration"`
}

// A CloudFormation `AWS::ApiGateway::BasePathMapping`.
type CfnBasePathMapping interface {
	awscdk.CfnResource
	awscdk.IInspectable
	BasePath() *string
	SetBasePath(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DomainName() *string
	SetDomainName(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Stack() awscdk.Stack
	Stage() *string
	SetStage(val *string)
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

// The jsii proxy struct for CfnBasePathMapping
type jsiiProxy_CfnBasePathMapping struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnBasePathMapping) BasePath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"basePath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) Stage() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnBasePathMapping) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::BasePathMapping`.
func NewCfnBasePathMapping(scope awscdk.Construct, id *string, props *CfnBasePathMappingProps) CfnBasePathMapping {
	_init_.Initialize()

	j := jsiiProxy_CfnBasePathMapping{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnBasePathMapping",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::BasePathMapping`.
func NewCfnBasePathMapping_Override(c CfnBasePathMapping, scope awscdk.Construct, id *string, props *CfnBasePathMappingProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnBasePathMapping",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnBasePathMapping) SetBasePath(val *string) {
	_jsii_.Set(
		j,
		"basePath",
		val,
	)
}

func (j *jsiiProxy_CfnBasePathMapping) SetDomainName(val *string) {
	_jsii_.Set(
		j,
		"domainName",
		val,
	)
}

func (j *jsiiProxy_CfnBasePathMapping) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

func (j *jsiiProxy_CfnBasePathMapping) SetStage(val *string) {
	_jsii_.Set(
		j,
		"stage",
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
func CfnBasePathMapping_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnBasePathMapping",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnBasePathMapping_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnBasePathMapping",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnBasePathMapping_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnBasePathMapping",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnBasePathMapping_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnBasePathMapping",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnBasePathMapping) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnBasePathMapping) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnBasePathMapping) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnBasePathMapping) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnBasePathMapping) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnBasePathMapping) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnBasePathMapping) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnBasePathMapping) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnBasePathMapping) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnBasePathMapping) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnBasePathMapping) OnPrepare() {
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
func (c *jsiiProxy_CfnBasePathMapping) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnBasePathMapping) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnBasePathMapping) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnBasePathMapping) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnBasePathMapping) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnBasePathMapping) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnBasePathMapping) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnBasePathMapping) ToString() *string {
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
func (c *jsiiProxy_CfnBasePathMapping) Validate() *[]*string {
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
func (c *jsiiProxy_CfnBasePathMapping) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::BasePathMapping`.
type CfnBasePathMappingProps struct {
	// `AWS::ApiGateway::BasePathMapping.DomainName`.
	DomainName *string `json:"domainName"`
	// `AWS::ApiGateway::BasePathMapping.BasePath`.
	BasePath *string `json:"basePath"`
	// `AWS::ApiGateway::BasePathMapping.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::BasePathMapping.Stage`.
	Stage *string `json:"stage"`
}

// A CloudFormation `AWS::ApiGateway::ClientCertificate`.
type CfnClientCertificate interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrClientCertificateId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
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

// The jsii proxy struct for CfnClientCertificate
type jsiiProxy_CfnClientCertificate struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnClientCertificate) AttrClientCertificateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrClientCertificateId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnClientCertificate) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::ClientCertificate`.
func NewCfnClientCertificate(scope awscdk.Construct, id *string, props *CfnClientCertificateProps) CfnClientCertificate {
	_init_.Initialize()

	j := jsiiProxy_CfnClientCertificate{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnClientCertificate",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::ClientCertificate`.
func NewCfnClientCertificate_Override(c CfnClientCertificate, scope awscdk.Construct, id *string, props *CfnClientCertificateProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnClientCertificate",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnClientCertificate) SetDescription(val *string) {
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
func CfnClientCertificate_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnClientCertificate",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnClientCertificate_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnClientCertificate",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnClientCertificate_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnClientCertificate",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnClientCertificate_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnClientCertificate",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnClientCertificate) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnClientCertificate) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnClientCertificate) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnClientCertificate) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnClientCertificate) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnClientCertificate) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnClientCertificate) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnClientCertificate) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnClientCertificate) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnClientCertificate) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnClientCertificate) OnPrepare() {
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
func (c *jsiiProxy_CfnClientCertificate) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnClientCertificate) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnClientCertificate) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnClientCertificate) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnClientCertificate) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnClientCertificate) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnClientCertificate) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnClientCertificate) ToString() *string {
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
func (c *jsiiProxy_CfnClientCertificate) Validate() *[]*string {
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
func (c *jsiiProxy_CfnClientCertificate) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::ClientCertificate`.
type CfnClientCertificateProps struct {
	// `AWS::ApiGateway::ClientCertificate.Description`.
	Description *string `json:"description"`
	// `AWS::ApiGateway::ClientCertificate.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::ApiGateway::Deployment`.
type CfnDeployment interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DeploymentCanarySettings() interface{}
	SetDeploymentCanarySettings(val interface{})
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Stack() awscdk.Stack
	StageDescription() interface{}
	SetStageDescription(val interface{})
	StageName() *string
	SetStageName(val *string)
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

// The jsii proxy struct for CfnDeployment
type jsiiProxy_CfnDeployment struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDeployment) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) DeploymentCanarySettings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"deploymentCanarySettings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) StageDescription() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stageDescription",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeployment) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::Deployment`.
func NewCfnDeployment(scope awscdk.Construct, id *string, props *CfnDeploymentProps) CfnDeployment {
	_init_.Initialize()

	j := jsiiProxy_CfnDeployment{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDeployment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::Deployment`.
func NewCfnDeployment_Override(c CfnDeployment, scope awscdk.Construct, id *string, props *CfnDeploymentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDeployment",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDeployment) SetDeploymentCanarySettings(val interface{}) {
	_jsii_.Set(
		j,
		"deploymentCanarySettings",
		val,
	)
}

func (j *jsiiProxy_CfnDeployment) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnDeployment) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

func (j *jsiiProxy_CfnDeployment) SetStageDescription(val interface{}) {
	_jsii_.Set(
		j,
		"stageDescription",
		val,
	)
}

func (j *jsiiProxy_CfnDeployment) SetStageName(val *string) {
	_jsii_.Set(
		j,
		"stageName",
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
func CfnDeployment_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDeployment",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnDeployment_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDeployment",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnDeployment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDeployment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDeployment_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnDeployment",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnDeployment) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnDeployment) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnDeployment) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnDeployment) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnDeployment) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnDeployment) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnDeployment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnDeployment) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnDeployment) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnDeployment) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnDeployment) OnPrepare() {
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
func (c *jsiiProxy_CfnDeployment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDeployment) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnDeployment) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnDeployment) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnDeployment) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnDeployment) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnDeployment) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDeployment) ToString() *string {
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
func (c *jsiiProxy_CfnDeployment) Validate() *[]*string {
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
func (c *jsiiProxy_CfnDeployment) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnDeployment_AccessLogSettingProperty struct {
	// `CfnDeployment.AccessLogSettingProperty.DestinationArn`.
	DestinationArn *string `json:"destinationArn"`
	// `CfnDeployment.AccessLogSettingProperty.Format`.
	Format *string `json:"format"`
}

type CfnDeployment_CanarySettingProperty struct {
	// `CfnDeployment.CanarySettingProperty.PercentTraffic`.
	PercentTraffic *float64 `json:"percentTraffic"`
	// `CfnDeployment.CanarySettingProperty.StageVariableOverrides`.
	StageVariableOverrides interface{} `json:"stageVariableOverrides"`
	// `CfnDeployment.CanarySettingProperty.UseStageCache`.
	UseStageCache interface{} `json:"useStageCache"`
}

type CfnDeployment_DeploymentCanarySettingsProperty struct {
	// `CfnDeployment.DeploymentCanarySettingsProperty.PercentTraffic`.
	PercentTraffic *float64 `json:"percentTraffic"`
	// `CfnDeployment.DeploymentCanarySettingsProperty.StageVariableOverrides`.
	StageVariableOverrides interface{} `json:"stageVariableOverrides"`
	// `CfnDeployment.DeploymentCanarySettingsProperty.UseStageCache`.
	UseStageCache interface{} `json:"useStageCache"`
}

type CfnDeployment_MethodSettingProperty struct {
	// `CfnDeployment.MethodSettingProperty.CacheDataEncrypted`.
	CacheDataEncrypted interface{} `json:"cacheDataEncrypted"`
	// `CfnDeployment.MethodSettingProperty.CacheTtlInSeconds`.
	CacheTtlInSeconds *float64 `json:"cacheTtlInSeconds"`
	// `CfnDeployment.MethodSettingProperty.CachingEnabled`.
	CachingEnabled interface{} `json:"cachingEnabled"`
	// `CfnDeployment.MethodSettingProperty.DataTraceEnabled`.
	DataTraceEnabled interface{} `json:"dataTraceEnabled"`
	// `CfnDeployment.MethodSettingProperty.HttpMethod`.
	HttpMethod *string `json:"httpMethod"`
	// `CfnDeployment.MethodSettingProperty.LoggingLevel`.
	LoggingLevel *string `json:"loggingLevel"`
	// `CfnDeployment.MethodSettingProperty.MetricsEnabled`.
	MetricsEnabled interface{} `json:"metricsEnabled"`
	// `CfnDeployment.MethodSettingProperty.ResourcePath`.
	ResourcePath *string `json:"resourcePath"`
	// `CfnDeployment.MethodSettingProperty.ThrottlingBurstLimit`.
	ThrottlingBurstLimit *float64 `json:"throttlingBurstLimit"`
	// `CfnDeployment.MethodSettingProperty.ThrottlingRateLimit`.
	ThrottlingRateLimit *float64 `json:"throttlingRateLimit"`
}

type CfnDeployment_StageDescriptionProperty struct {
	// `CfnDeployment.StageDescriptionProperty.AccessLogSetting`.
	AccessLogSetting interface{} `json:"accessLogSetting"`
	// `CfnDeployment.StageDescriptionProperty.CacheClusterEnabled`.
	CacheClusterEnabled interface{} `json:"cacheClusterEnabled"`
	// `CfnDeployment.StageDescriptionProperty.CacheClusterSize`.
	CacheClusterSize *string `json:"cacheClusterSize"`
	// `CfnDeployment.StageDescriptionProperty.CacheDataEncrypted`.
	CacheDataEncrypted interface{} `json:"cacheDataEncrypted"`
	// `CfnDeployment.StageDescriptionProperty.CacheTtlInSeconds`.
	CacheTtlInSeconds *float64 `json:"cacheTtlInSeconds"`
	// `CfnDeployment.StageDescriptionProperty.CachingEnabled`.
	CachingEnabled interface{} `json:"cachingEnabled"`
	// `CfnDeployment.StageDescriptionProperty.CanarySetting`.
	CanarySetting interface{} `json:"canarySetting"`
	// `CfnDeployment.StageDescriptionProperty.ClientCertificateId`.
	ClientCertificateId *string `json:"clientCertificateId"`
	// `CfnDeployment.StageDescriptionProperty.DataTraceEnabled`.
	DataTraceEnabled interface{} `json:"dataTraceEnabled"`
	// `CfnDeployment.StageDescriptionProperty.Description`.
	Description *string `json:"description"`
	// `CfnDeployment.StageDescriptionProperty.DocumentationVersion`.
	DocumentationVersion *string `json:"documentationVersion"`
	// `CfnDeployment.StageDescriptionProperty.LoggingLevel`.
	LoggingLevel *string `json:"loggingLevel"`
	// `CfnDeployment.StageDescriptionProperty.MethodSettings`.
	MethodSettings interface{} `json:"methodSettings"`
	// `CfnDeployment.StageDescriptionProperty.MetricsEnabled`.
	MetricsEnabled interface{} `json:"metricsEnabled"`
	// `CfnDeployment.StageDescriptionProperty.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `CfnDeployment.StageDescriptionProperty.ThrottlingBurstLimit`.
	ThrottlingBurstLimit *float64 `json:"throttlingBurstLimit"`
	// `CfnDeployment.StageDescriptionProperty.ThrottlingRateLimit`.
	ThrottlingRateLimit *float64 `json:"throttlingRateLimit"`
	// `CfnDeployment.StageDescriptionProperty.TracingEnabled`.
	TracingEnabled interface{} `json:"tracingEnabled"`
	// `CfnDeployment.StageDescriptionProperty.Variables`.
	Variables interface{} `json:"variables"`
}

// Properties for defining a `AWS::ApiGateway::Deployment`.
type CfnDeploymentProps struct {
	// `AWS::ApiGateway::Deployment.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::Deployment.DeploymentCanarySettings`.
	DeploymentCanarySettings interface{} `json:"deploymentCanarySettings"`
	// `AWS::ApiGateway::Deployment.Description`.
	Description *string `json:"description"`
	// `AWS::ApiGateway::Deployment.StageDescription`.
	StageDescription interface{} `json:"stageDescription"`
	// `AWS::ApiGateway::Deployment.StageName`.
	StageName *string `json:"stageName"`
}

// A CloudFormation `AWS::ApiGatewayV2::Deployment`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-deployment.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnDeploymentV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	StageName() *string
	SetStageName(val *string)
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

// The jsii proxy struct for CfnDeploymentV2
type jsiiProxy_CfnDeploymentV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDeploymentV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDeploymentV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::Deployment`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnDeploymentV2(scope awscdk.Construct, id *string, props *CfnDeploymentV2Props) CfnDeploymentV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnDeploymentV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDeploymentV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::Deployment`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnDeploymentV2_Override(c CfnDeploymentV2, scope awscdk.Construct, id *string, props *CfnDeploymentV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDeploymentV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDeploymentV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnDeploymentV2) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnDeploymentV2) SetStageName(val *string) {
	_jsii_.Set(
		j,
		"stageName",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnDeploymentV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDeploymentV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnDeploymentV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDeploymentV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnDeploymentV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDeploymentV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDeploymentV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnDeploymentV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDeploymentV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGatewayV2::Deployment`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-deployment.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnDeploymentV2Props struct {
	// `AWS::ApiGatewayV2::Deployment.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-deployment.html#cfn-apigatewayv2-deployment-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::Deployment.Description`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-deployment.html#cfn-apigatewayv2-deployment-description
	//
	// Deprecated: moved to package aws-apigatewayv2
	Description *string `json:"description"`
	// `AWS::ApiGatewayV2::Deployment.StageName`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-deployment.html#cfn-apigatewayv2-deployment-stagename
	//
	// Deprecated: moved to package aws-apigatewayv2
	StageName *string `json:"stageName"`
}

// A CloudFormation `AWS::ApiGateway::DocumentationPart`.
type CfnDocumentationPart interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Location() interface{}
	SetLocation(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	Properties() *string
	SetProperties(val *string)
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
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

// The jsii proxy struct for CfnDocumentationPart
type jsiiProxy_CfnDocumentationPart struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDocumentationPart) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) Location() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) Properties() *string {
	var returns *string
	_jsii_.Get(
		j,
		"properties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationPart) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::DocumentationPart`.
func NewCfnDocumentationPart(scope awscdk.Construct, id *string, props *CfnDocumentationPartProps) CfnDocumentationPart {
	_init_.Initialize()

	j := jsiiProxy_CfnDocumentationPart{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDocumentationPart",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::DocumentationPart`.
func NewCfnDocumentationPart_Override(c CfnDocumentationPart, scope awscdk.Construct, id *string, props *CfnDocumentationPartProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDocumentationPart",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDocumentationPart) SetLocation(val interface{}) {
	_jsii_.Set(
		j,
		"location",
		val,
	)
}

func (j *jsiiProxy_CfnDocumentationPart) SetProperties(val *string) {
	_jsii_.Set(
		j,
		"properties",
		val,
	)
}

func (j *jsiiProxy_CfnDocumentationPart) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
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
func CfnDocumentationPart_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDocumentationPart",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnDocumentationPart_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDocumentationPart",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnDocumentationPart_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDocumentationPart",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDocumentationPart_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnDocumentationPart",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnDocumentationPart) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnDocumentationPart) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnDocumentationPart) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnDocumentationPart) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnDocumentationPart) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnDocumentationPart) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnDocumentationPart) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnDocumentationPart) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnDocumentationPart) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnDocumentationPart) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnDocumentationPart) OnPrepare() {
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
func (c *jsiiProxy_CfnDocumentationPart) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDocumentationPart) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnDocumentationPart) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnDocumentationPart) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnDocumentationPart) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnDocumentationPart) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnDocumentationPart) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDocumentationPart) ToString() *string {
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
func (c *jsiiProxy_CfnDocumentationPart) Validate() *[]*string {
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
func (c *jsiiProxy_CfnDocumentationPart) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnDocumentationPart_LocationProperty struct {
	// `CfnDocumentationPart.LocationProperty.Method`.
	Method *string `json:"method"`
	// `CfnDocumentationPart.LocationProperty.Name`.
	Name *string `json:"name"`
	// `CfnDocumentationPart.LocationProperty.Path`.
	Path *string `json:"path"`
	// `CfnDocumentationPart.LocationProperty.StatusCode`.
	StatusCode *string `json:"statusCode"`
	// `CfnDocumentationPart.LocationProperty.Type`.
	Type *string `json:"type"`
}

// Properties for defining a `AWS::ApiGateway::DocumentationPart`.
type CfnDocumentationPartProps struct {
	// `AWS::ApiGateway::DocumentationPart.Location`.
	Location interface{} `json:"location"`
	// `AWS::ApiGateway::DocumentationPart.Properties`.
	Properties *string `json:"properties"`
	// `AWS::ApiGateway::DocumentationPart.RestApiId`.
	RestApiId *string `json:"restApiId"`
}

// A CloudFormation `AWS::ApiGateway::DocumentationVersion`.
type CfnDocumentationVersion interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	DocumentationVersion() *string
	SetDocumentationVersion(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
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

// The jsii proxy struct for CfnDocumentationVersion
type jsiiProxy_CfnDocumentationVersion struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDocumentationVersion) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) DocumentationVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"documentationVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDocumentationVersion) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::DocumentationVersion`.
func NewCfnDocumentationVersion(scope awscdk.Construct, id *string, props *CfnDocumentationVersionProps) CfnDocumentationVersion {
	_init_.Initialize()

	j := jsiiProxy_CfnDocumentationVersion{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDocumentationVersion",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::DocumentationVersion`.
func NewCfnDocumentationVersion_Override(c CfnDocumentationVersion, scope awscdk.Construct, id *string, props *CfnDocumentationVersionProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDocumentationVersion",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDocumentationVersion) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnDocumentationVersion) SetDocumentationVersion(val *string) {
	_jsii_.Set(
		j,
		"documentationVersion",
		val,
	)
}

func (j *jsiiProxy_CfnDocumentationVersion) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
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
func CfnDocumentationVersion_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDocumentationVersion",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnDocumentationVersion_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDocumentationVersion",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnDocumentationVersion_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDocumentationVersion",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDocumentationVersion_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnDocumentationVersion",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnDocumentationVersion) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnDocumentationVersion) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnDocumentationVersion) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnDocumentationVersion) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnDocumentationVersion) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnDocumentationVersion) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnDocumentationVersion) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnDocumentationVersion) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnDocumentationVersion) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnDocumentationVersion) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnDocumentationVersion) OnPrepare() {
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
func (c *jsiiProxy_CfnDocumentationVersion) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDocumentationVersion) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnDocumentationVersion) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnDocumentationVersion) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnDocumentationVersion) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnDocumentationVersion) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnDocumentationVersion) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDocumentationVersion) ToString() *string {
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
func (c *jsiiProxy_CfnDocumentationVersion) Validate() *[]*string {
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
func (c *jsiiProxy_CfnDocumentationVersion) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::DocumentationVersion`.
type CfnDocumentationVersionProps struct {
	// `AWS::ApiGateway::DocumentationVersion.DocumentationVersion`.
	DocumentationVersion *string `json:"documentationVersion"`
	// `AWS::ApiGateway::DocumentationVersion.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::DocumentationVersion.Description`.
	Description *string `json:"description"`
}

// A CloudFormation `AWS::ApiGateway::DomainName`.
type CfnDomainName interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrDistributionDomainName() *string
	AttrDistributionHostedZoneId() *string
	AttrRegionalDomainName() *string
	AttrRegionalHostedZoneId() *string
	CertificateArn() *string
	SetCertificateArn(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DomainName() *string
	SetDomainName(val *string)
	EndpointConfiguration() interface{}
	SetEndpointConfiguration(val interface{})
	LogicalId() *string
	MutualTlsAuthentication() interface{}
	SetMutualTlsAuthentication(val interface{})
	Node() awscdk.ConstructNode
	Ref() *string
	RegionalCertificateArn() *string
	SetRegionalCertificateArn(val *string)
	SecurityPolicy() *string
	SetSecurityPolicy(val *string)
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

// The jsii proxy struct for CfnDomainName
type jsiiProxy_CfnDomainName struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDomainName) AttrDistributionDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDistributionDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) AttrDistributionHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDistributionHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) AttrRegionalDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRegionalDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) AttrRegionalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRegionalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) CertificateArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certificateArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) EndpointConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"endpointConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) MutualTlsAuthentication() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"mutualTlsAuthentication",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) RegionalCertificateArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"regionalCertificateArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) SecurityPolicy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"securityPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainName) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::DomainName`.
func NewCfnDomainName(scope awscdk.Construct, id *string, props *CfnDomainNameProps) CfnDomainName {
	_init_.Initialize()

	j := jsiiProxy_CfnDomainName{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDomainName",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::DomainName`.
func NewCfnDomainName_Override(c CfnDomainName, scope awscdk.Construct, id *string, props *CfnDomainNameProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDomainName",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDomainName) SetCertificateArn(val *string) {
	_jsii_.Set(
		j,
		"certificateArn",
		val,
	)
}

func (j *jsiiProxy_CfnDomainName) SetDomainName(val *string) {
	_jsii_.Set(
		j,
		"domainName",
		val,
	)
}

func (j *jsiiProxy_CfnDomainName) SetEndpointConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"endpointConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnDomainName) SetMutualTlsAuthentication(val interface{}) {
	_jsii_.Set(
		j,
		"mutualTlsAuthentication",
		val,
	)
}

func (j *jsiiProxy_CfnDomainName) SetRegionalCertificateArn(val *string) {
	_jsii_.Set(
		j,
		"regionalCertificateArn",
		val,
	)
}

func (j *jsiiProxy_CfnDomainName) SetSecurityPolicy(val *string) {
	_jsii_.Set(
		j,
		"securityPolicy",
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
func CfnDomainName_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDomainName",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnDomainName_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDomainName",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnDomainName_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDomainName",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDomainName_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnDomainName",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnDomainName) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnDomainName) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnDomainName) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnDomainName) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnDomainName) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnDomainName) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnDomainName) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnDomainName) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnDomainName) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnDomainName) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnDomainName) OnPrepare() {
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
func (c *jsiiProxy_CfnDomainName) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDomainName) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnDomainName) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnDomainName) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnDomainName) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnDomainName) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnDomainName) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnDomainName) ToString() *string {
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
func (c *jsiiProxy_CfnDomainName) Validate() *[]*string {
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
func (c *jsiiProxy_CfnDomainName) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnDomainName_EndpointConfigurationProperty struct {
	// `CfnDomainName.EndpointConfigurationProperty.Types`.
	Types *[]*string `json:"types"`
}

type CfnDomainName_MutualTlsAuthenticationProperty struct {
	// `CfnDomainName.MutualTlsAuthenticationProperty.TruststoreUri`.
	TruststoreUri *string `json:"truststoreUri"`
	// `CfnDomainName.MutualTlsAuthenticationProperty.TruststoreVersion`.
	TruststoreVersion *string `json:"truststoreVersion"`
}

// Properties for defining a `AWS::ApiGateway::DomainName`.
type CfnDomainNameProps struct {
	// `AWS::ApiGateway::DomainName.CertificateArn`.
	CertificateArn *string `json:"certificateArn"`
	// `AWS::ApiGateway::DomainName.DomainName`.
	DomainName *string `json:"domainName"`
	// `AWS::ApiGateway::DomainName.EndpointConfiguration`.
	EndpointConfiguration interface{} `json:"endpointConfiguration"`
	// `AWS::ApiGateway::DomainName.MutualTlsAuthentication`.
	MutualTlsAuthentication interface{} `json:"mutualTlsAuthentication"`
	// `AWS::ApiGateway::DomainName.RegionalCertificateArn`.
	RegionalCertificateArn *string `json:"regionalCertificateArn"`
	// `AWS::ApiGateway::DomainName.SecurityPolicy`.
	SecurityPolicy *string `json:"securityPolicy"`
	// `AWS::ApiGateway::DomainName.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::ApiGatewayV2::DomainName`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-domainname.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnDomainNameV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrRegionalDomainName() *string
	AttrRegionalHostedZoneId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DomainName() *string
	SetDomainName(val *string)
	DomainNameConfigurations() interface{}
	SetDomainNameConfigurations(val interface{})
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
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

// The jsii proxy struct for CfnDomainNameV2
type jsiiProxy_CfnDomainNameV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnDomainNameV2) AttrRegionalDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRegionalDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) AttrRegionalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRegionalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) DomainNameConfigurations() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"domainNameConfigurations",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnDomainNameV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::DomainName`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnDomainNameV2(scope awscdk.Construct, id *string, props *CfnDomainNameV2Props) CfnDomainNameV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnDomainNameV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDomainNameV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::DomainName`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnDomainNameV2_Override(c CfnDomainNameV2, scope awscdk.Construct, id *string, props *CfnDomainNameV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnDomainNameV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnDomainNameV2) SetDomainName(val *string) {
	_jsii_.Set(
		j,
		"domainName",
		val,
	)
}

func (j *jsiiProxy_CfnDomainNameV2) SetDomainNameConfigurations(val interface{}) {
	_jsii_.Set(
		j,
		"domainNameConfigurations",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnDomainNameV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDomainNameV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnDomainNameV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDomainNameV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnDomainNameV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnDomainNameV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnDomainNameV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnDomainNameV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnDomainNameV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnDomainNameV2_DomainNameConfigurationProperty struct {
	// `CfnDomainNameV2.DomainNameConfigurationProperty.CertificateArn`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-certificatearn
	//
	// Deprecated: moved to package aws-apigatewayv2
	CertificateArn *string `json:"certificateArn"`
	// `CfnDomainNameV2.DomainNameConfigurationProperty.CertificateName`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-certificatename
	//
	// Deprecated: moved to package aws-apigatewayv2
	CertificateName *string `json:"certificateName"`
	// `CfnDomainNameV2.DomainNameConfigurationProperty.EndpointType`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-domainname-domainnameconfiguration.html#cfn-apigatewayv2-domainname-domainnameconfiguration-endpointtype
	//
	// Deprecated: moved to package aws-apigatewayv2
	EndpointType *string `json:"endpointType"`
}

// Properties for defining a `AWS::ApiGatewayV2::DomainName`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-domainname.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnDomainNameV2Props struct {
	// `AWS::ApiGatewayV2::DomainName.DomainName`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-domainname.html#cfn-apigatewayv2-domainname-domainname
	//
	// Deprecated: moved to package aws-apigatewayv2
	DomainName *string `json:"domainName"`
	// `AWS::ApiGatewayV2::DomainName.DomainNameConfigurations`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-domainname.html#cfn-apigatewayv2-domainname-domainnameconfigurations
	//
	// Deprecated: moved to package aws-apigatewayv2
	DomainNameConfigurations interface{} `json:"domainNameConfigurations"`
	// `AWS::ApiGatewayV2::DomainName.Tags`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-domainname.html#cfn-apigatewayv2-domainname-tags
	//
	// Deprecated: moved to package aws-apigatewayv2
	Tags interface{} `json:"tags"`
}

// A CloudFormation `AWS::ApiGateway::GatewayResponse`.
type CfnGatewayResponse interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	ResponseParameters() interface{}
	SetResponseParameters(val interface{})
	ResponseTemplates() interface{}
	SetResponseTemplates(val interface{})
	ResponseType() *string
	SetResponseType(val *string)
	RestApiId() *string
	SetRestApiId(val *string)
	Stack() awscdk.Stack
	StatusCode() *string
	SetStatusCode(val *string)
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

// The jsii proxy struct for CfnGatewayResponse
type jsiiProxy_CfnGatewayResponse struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnGatewayResponse) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) ResponseParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"responseParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) ResponseTemplates() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"responseTemplates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) ResponseType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"responseType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) StatusCode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"statusCode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnGatewayResponse) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::GatewayResponse`.
func NewCfnGatewayResponse(scope awscdk.Construct, id *string, props *CfnGatewayResponseProps) CfnGatewayResponse {
	_init_.Initialize()

	j := jsiiProxy_CfnGatewayResponse{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnGatewayResponse",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::GatewayResponse`.
func NewCfnGatewayResponse_Override(c CfnGatewayResponse, scope awscdk.Construct, id *string, props *CfnGatewayResponseProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnGatewayResponse",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnGatewayResponse) SetResponseParameters(val interface{}) {
	_jsii_.Set(
		j,
		"responseParameters",
		val,
	)
}

func (j *jsiiProxy_CfnGatewayResponse) SetResponseTemplates(val interface{}) {
	_jsii_.Set(
		j,
		"responseTemplates",
		val,
	)
}

func (j *jsiiProxy_CfnGatewayResponse) SetResponseType(val *string) {
	_jsii_.Set(
		j,
		"responseType",
		val,
	)
}

func (j *jsiiProxy_CfnGatewayResponse) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

func (j *jsiiProxy_CfnGatewayResponse) SetStatusCode(val *string) {
	_jsii_.Set(
		j,
		"statusCode",
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
func CfnGatewayResponse_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnGatewayResponse",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnGatewayResponse_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnGatewayResponse",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnGatewayResponse_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnGatewayResponse",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnGatewayResponse_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnGatewayResponse",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnGatewayResponse) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnGatewayResponse) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnGatewayResponse) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnGatewayResponse) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnGatewayResponse) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnGatewayResponse) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnGatewayResponse) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnGatewayResponse) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnGatewayResponse) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnGatewayResponse) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnGatewayResponse) OnPrepare() {
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
func (c *jsiiProxy_CfnGatewayResponse) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnGatewayResponse) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnGatewayResponse) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnGatewayResponse) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnGatewayResponse) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnGatewayResponse) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnGatewayResponse) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnGatewayResponse) ToString() *string {
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
func (c *jsiiProxy_CfnGatewayResponse) Validate() *[]*string {
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
func (c *jsiiProxy_CfnGatewayResponse) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::GatewayResponse`.
type CfnGatewayResponseProps struct {
	// `AWS::ApiGateway::GatewayResponse.ResponseType`.
	ResponseType *string `json:"responseType"`
	// `AWS::ApiGateway::GatewayResponse.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::GatewayResponse.ResponseParameters`.
	ResponseParameters interface{} `json:"responseParameters"`
	// `AWS::ApiGateway::GatewayResponse.ResponseTemplates`.
	ResponseTemplates interface{} `json:"responseTemplates"`
	// `AWS::ApiGateway::GatewayResponse.StatusCode`.
	StatusCode *string `json:"statusCode"`
}

// A CloudFormation `AWS::ApiGatewayV2::IntegrationResponse`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnIntegrationResponseV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ContentHandlingStrategy() *string
	SetContentHandlingStrategy(val *string)
	CreationStack() *[]*string
	IntegrationId() *string
	SetIntegrationId(val *string)
	IntegrationResponseKey() *string
	SetIntegrationResponseKey(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	ResponseParameters() interface{}
	SetResponseParameters(val interface{})
	ResponseTemplates() interface{}
	SetResponseTemplates(val interface{})
	Stack() awscdk.Stack
	TemplateSelectionExpression() *string
	SetTemplateSelectionExpression(val *string)
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

// The jsii proxy struct for CfnIntegrationResponseV2
type jsiiProxy_CfnIntegrationResponseV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnIntegrationResponseV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) ContentHandlingStrategy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"contentHandlingStrategy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) IntegrationId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"integrationId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) IntegrationResponseKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"integrationResponseKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) ResponseParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"responseParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) ResponseTemplates() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"responseTemplates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) TemplateSelectionExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateSelectionExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationResponseV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::IntegrationResponse`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnIntegrationResponseV2(scope awscdk.Construct, id *string, props *CfnIntegrationResponseV2Props) CfnIntegrationResponseV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnIntegrationResponseV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnIntegrationResponseV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::IntegrationResponse`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnIntegrationResponseV2_Override(c CfnIntegrationResponseV2, scope awscdk.Construct, id *string, props *CfnIntegrationResponseV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnIntegrationResponseV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnIntegrationResponseV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationResponseV2) SetContentHandlingStrategy(val *string) {
	_jsii_.Set(
		j,
		"contentHandlingStrategy",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationResponseV2) SetIntegrationId(val *string) {
	_jsii_.Set(
		j,
		"integrationId",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationResponseV2) SetIntegrationResponseKey(val *string) {
	_jsii_.Set(
		j,
		"integrationResponseKey",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationResponseV2) SetResponseParameters(val interface{}) {
	_jsii_.Set(
		j,
		"responseParameters",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationResponseV2) SetResponseTemplates(val interface{}) {
	_jsii_.Set(
		j,
		"responseTemplates",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationResponseV2) SetTemplateSelectionExpression(val *string) {
	_jsii_.Set(
		j,
		"templateSelectionExpression",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnIntegrationResponseV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnIntegrationResponseV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnIntegrationResponseV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnIntegrationResponseV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnIntegrationResponseV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnIntegrationResponseV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnIntegrationResponseV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnIntegrationResponseV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationResponseV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGatewayV2::IntegrationResponse`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnIntegrationResponseV2Props struct {
	// `AWS::ApiGatewayV2::IntegrationResponse.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::IntegrationResponse.IntegrationId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-integrationid
	//
	// Deprecated: moved to package aws-apigatewayv2
	IntegrationId *string `json:"integrationId"`
	// `AWS::ApiGatewayV2::IntegrationResponse.IntegrationResponseKey`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-integrationresponsekey
	//
	// Deprecated: moved to package aws-apigatewayv2
	IntegrationResponseKey *string `json:"integrationResponseKey"`
	// `AWS::ApiGatewayV2::IntegrationResponse.ContentHandlingStrategy`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-contenthandlingstrategy
	//
	// Deprecated: moved to package aws-apigatewayv2
	ContentHandlingStrategy *string `json:"contentHandlingStrategy"`
	// `AWS::ApiGatewayV2::IntegrationResponse.ResponseParameters`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-responseparameters
	//
	// Deprecated: moved to package aws-apigatewayv2
	ResponseParameters interface{} `json:"responseParameters"`
	// `AWS::ApiGatewayV2::IntegrationResponse.ResponseTemplates`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-responsetemplates
	//
	// Deprecated: moved to package aws-apigatewayv2
	ResponseTemplates interface{} `json:"responseTemplates"`
	// `AWS::ApiGatewayV2::IntegrationResponse.TemplateSelectionExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-templateselectionexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	TemplateSelectionExpression *string `json:"templateSelectionExpression"`
}

// A CloudFormation `AWS::ApiGatewayV2::Integration`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnIntegrationV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ConnectionType() *string
	SetConnectionType(val *string)
	ContentHandlingStrategy() *string
	SetContentHandlingStrategy(val *string)
	CreationStack() *[]*string
	CredentialsArn() *string
	SetCredentialsArn(val *string)
	Description() *string
	SetDescription(val *string)
	IntegrationMethod() *string
	SetIntegrationMethod(val *string)
	IntegrationType() *string
	SetIntegrationType(val *string)
	IntegrationUri() *string
	SetIntegrationUri(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	PassthroughBehavior() *string
	SetPassthroughBehavior(val *string)
	PayloadFormatVersion() *string
	SetPayloadFormatVersion(val *string)
	Ref() *string
	RequestParameters() interface{}
	SetRequestParameters(val interface{})
	RequestTemplates() interface{}
	SetRequestTemplates(val interface{})
	Stack() awscdk.Stack
	TemplateSelectionExpression() *string
	SetTemplateSelectionExpression(val *string)
	TimeoutInMillis() *float64
	SetTimeoutInMillis(val *float64)
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

// The jsii proxy struct for CfnIntegrationV2
type jsiiProxy_CfnIntegrationV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnIntegrationV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) ConnectionType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"connectionType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) ContentHandlingStrategy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"contentHandlingStrategy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) CredentialsArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"credentialsArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) IntegrationMethod() *string {
	var returns *string
	_jsii_.Get(
		j,
		"integrationMethod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) IntegrationType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"integrationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) IntegrationUri() *string {
	var returns *string
	_jsii_.Get(
		j,
		"integrationUri",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) PassthroughBehavior() *string {
	var returns *string
	_jsii_.Get(
		j,
		"passthroughBehavior",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) PayloadFormatVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"payloadFormatVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) RequestParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"requestParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) RequestTemplates() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"requestTemplates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) TemplateSelectionExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"templateSelectionExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) TimeoutInMillis() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"timeoutInMillis",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnIntegrationV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::Integration`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnIntegrationV2(scope awscdk.Construct, id *string, props *CfnIntegrationV2Props) CfnIntegrationV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnIntegrationV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnIntegrationV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::Integration`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnIntegrationV2_Override(c CfnIntegrationV2, scope awscdk.Construct, id *string, props *CfnIntegrationV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnIntegrationV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetConnectionType(val *string) {
	_jsii_.Set(
		j,
		"connectionType",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetContentHandlingStrategy(val *string) {
	_jsii_.Set(
		j,
		"contentHandlingStrategy",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetCredentialsArn(val *string) {
	_jsii_.Set(
		j,
		"credentialsArn",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetIntegrationMethod(val *string) {
	_jsii_.Set(
		j,
		"integrationMethod",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetIntegrationType(val *string) {
	_jsii_.Set(
		j,
		"integrationType",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetIntegrationUri(val *string) {
	_jsii_.Set(
		j,
		"integrationUri",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetPassthroughBehavior(val *string) {
	_jsii_.Set(
		j,
		"passthroughBehavior",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetPayloadFormatVersion(val *string) {
	_jsii_.Set(
		j,
		"payloadFormatVersion",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetRequestParameters(val interface{}) {
	_jsii_.Set(
		j,
		"requestParameters",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetRequestTemplates(val interface{}) {
	_jsii_.Set(
		j,
		"requestTemplates",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetTemplateSelectionExpression(val *string) {
	_jsii_.Set(
		j,
		"templateSelectionExpression",
		val,
	)
}

func (j *jsiiProxy_CfnIntegrationV2) SetTimeoutInMillis(val *float64) {
	_jsii_.Set(
		j,
		"timeoutInMillis",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnIntegrationV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnIntegrationV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnIntegrationV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnIntegrationV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnIntegrationV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnIntegrationV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnIntegrationV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnIntegrationV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnIntegrationV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGatewayV2::Integration`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnIntegrationV2Props struct {
	// `AWS::ApiGatewayV2::Integration.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::Integration.IntegrationType`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-integrationtype
	//
	// Deprecated: moved to package aws-apigatewayv2
	IntegrationType *string `json:"integrationType"`
	// `AWS::ApiGatewayV2::Integration.ConnectionType`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-connectiontype
	//
	// Deprecated: moved to package aws-apigatewayv2
	ConnectionType *string `json:"connectionType"`
	// `AWS::ApiGatewayV2::Integration.ContentHandlingStrategy`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-contenthandlingstrategy
	//
	// Deprecated: moved to package aws-apigatewayv2
	ContentHandlingStrategy *string `json:"contentHandlingStrategy"`
	// `AWS::ApiGatewayV2::Integration.CredentialsArn`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-credentialsarn
	//
	// Deprecated: moved to package aws-apigatewayv2
	CredentialsArn *string `json:"credentialsArn"`
	// `AWS::ApiGatewayV2::Integration.Description`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-description
	//
	// Deprecated: moved to package aws-apigatewayv2
	Description *string `json:"description"`
	// `AWS::ApiGatewayV2::Integration.IntegrationMethod`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-integrationmethod
	//
	// Deprecated: moved to package aws-apigatewayv2
	IntegrationMethod *string `json:"integrationMethod"`
	// `AWS::ApiGatewayV2::Integration.IntegrationUri`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-integrationuri
	//
	// Deprecated: moved to package aws-apigatewayv2
	IntegrationUri *string `json:"integrationUri"`
	// `AWS::ApiGatewayV2::Integration.PassthroughBehavior`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-passthroughbehavior
	//
	// Deprecated: moved to package aws-apigatewayv2
	PassthroughBehavior *string `json:"passthroughBehavior"`
	// `AWS::ApiGatewayV2::Integration.PayloadFormatVersion`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-payloadformatversion
	//
	// Deprecated: moved to package aws-apigatewayv2
	PayloadFormatVersion *string `json:"payloadFormatVersion"`
	// `AWS::ApiGatewayV2::Integration.RequestParameters`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-requestparameters
	//
	// Deprecated: moved to package aws-apigatewayv2
	RequestParameters interface{} `json:"requestParameters"`
	// `AWS::ApiGatewayV2::Integration.RequestTemplates`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-requesttemplates
	//
	// Deprecated: moved to package aws-apigatewayv2
	RequestTemplates interface{} `json:"requestTemplates"`
	// `AWS::ApiGatewayV2::Integration.TemplateSelectionExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-templateselectionexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	TemplateSelectionExpression *string `json:"templateSelectionExpression"`
	// `AWS::ApiGatewayV2::Integration.TimeoutInMillis`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html#cfn-apigatewayv2-integration-timeoutinmillis
	//
	// Deprecated: moved to package aws-apigatewayv2
	TimeoutInMillis *float64 `json:"timeoutInMillis"`
}

// A CloudFormation `AWS::ApiGateway::Method`.
type CfnMethod interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiKeyRequired() interface{}
	SetApiKeyRequired(val interface{})
	AuthorizationScopes() *[]*string
	SetAuthorizationScopes(val *[]*string)
	AuthorizationType() *string
	SetAuthorizationType(val *string)
	AuthorizerId() *string
	SetAuthorizerId(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	HttpMethod() *string
	SetHttpMethod(val *string)
	Integration() interface{}
	SetIntegration(val interface{})
	LogicalId() *string
	MethodResponses() interface{}
	SetMethodResponses(val interface{})
	Node() awscdk.ConstructNode
	OperationName() *string
	SetOperationName(val *string)
	Ref() *string
	RequestModels() interface{}
	SetRequestModels(val interface{})
	RequestParameters() interface{}
	SetRequestParameters(val interface{})
	RequestValidatorId() *string
	SetRequestValidatorId(val *string)
	ResourceId() *string
	SetResourceId(val *string)
	RestApiId() *string
	SetRestApiId(val *string)
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

// The jsii proxy struct for CfnMethod
type jsiiProxy_CfnMethod struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnMethod) ApiKeyRequired() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"apiKeyRequired",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) AuthorizationScopes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"authorizationScopes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) AuthorizationType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) AuthorizerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) HttpMethod() *string {
	var returns *string
	_jsii_.Get(
		j,
		"httpMethod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) Integration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"integration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) MethodResponses() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"methodResponses",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) OperationName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"operationName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) RequestModels() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"requestModels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) RequestParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"requestParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) RequestValidatorId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"requestValidatorId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) ResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnMethod) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::Method`.
func NewCfnMethod(scope awscdk.Construct, id *string, props *CfnMethodProps) CfnMethod {
	_init_.Initialize()

	j := jsiiProxy_CfnMethod{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnMethod",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::Method`.
func NewCfnMethod_Override(c CfnMethod, scope awscdk.Construct, id *string, props *CfnMethodProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnMethod",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnMethod) SetApiKeyRequired(val interface{}) {
	_jsii_.Set(
		j,
		"apiKeyRequired",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetAuthorizationScopes(val *[]*string) {
	_jsii_.Set(
		j,
		"authorizationScopes",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetAuthorizationType(val *string) {
	_jsii_.Set(
		j,
		"authorizationType",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetAuthorizerId(val *string) {
	_jsii_.Set(
		j,
		"authorizerId",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetHttpMethod(val *string) {
	_jsii_.Set(
		j,
		"httpMethod",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetIntegration(val interface{}) {
	_jsii_.Set(
		j,
		"integration",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetMethodResponses(val interface{}) {
	_jsii_.Set(
		j,
		"methodResponses",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetOperationName(val *string) {
	_jsii_.Set(
		j,
		"operationName",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetRequestModels(val interface{}) {
	_jsii_.Set(
		j,
		"requestModels",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetRequestParameters(val interface{}) {
	_jsii_.Set(
		j,
		"requestParameters",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetRequestValidatorId(val *string) {
	_jsii_.Set(
		j,
		"requestValidatorId",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetResourceId(val *string) {
	_jsii_.Set(
		j,
		"resourceId",
		val,
	)
}

func (j *jsiiProxy_CfnMethod) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
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
func CfnMethod_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnMethod",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnMethod_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnMethod",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnMethod_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnMethod",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnMethod_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnMethod",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnMethod) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnMethod) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnMethod) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnMethod) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnMethod) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnMethod) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnMethod) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnMethod) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnMethod) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnMethod) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnMethod) OnPrepare() {
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
func (c *jsiiProxy_CfnMethod) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnMethod) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnMethod) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnMethod) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnMethod) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnMethod) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnMethod) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnMethod) ToString() *string {
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
func (c *jsiiProxy_CfnMethod) Validate() *[]*string {
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
func (c *jsiiProxy_CfnMethod) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnMethod_IntegrationProperty struct {
	// `CfnMethod.IntegrationProperty.CacheKeyParameters`.
	CacheKeyParameters *[]*string `json:"cacheKeyParameters"`
	// `CfnMethod.IntegrationProperty.CacheNamespace`.
	CacheNamespace *string `json:"cacheNamespace"`
	// `CfnMethod.IntegrationProperty.ConnectionId`.
	ConnectionId *string `json:"connectionId"`
	// `CfnMethod.IntegrationProperty.ConnectionType`.
	ConnectionType *string `json:"connectionType"`
	// `CfnMethod.IntegrationProperty.ContentHandling`.
	ContentHandling *string `json:"contentHandling"`
	// `CfnMethod.IntegrationProperty.Credentials`.
	Credentials *string `json:"credentials"`
	// `CfnMethod.IntegrationProperty.IntegrationHttpMethod`.
	IntegrationHttpMethod *string `json:"integrationHttpMethod"`
	// `CfnMethod.IntegrationProperty.IntegrationResponses`.
	IntegrationResponses interface{} `json:"integrationResponses"`
	// `CfnMethod.IntegrationProperty.PassthroughBehavior`.
	PassthroughBehavior *string `json:"passthroughBehavior"`
	// `CfnMethod.IntegrationProperty.RequestParameters`.
	RequestParameters interface{} `json:"requestParameters"`
	// `CfnMethod.IntegrationProperty.RequestTemplates`.
	RequestTemplates interface{} `json:"requestTemplates"`
	// `CfnMethod.IntegrationProperty.TimeoutInMillis`.
	TimeoutInMillis *float64 `json:"timeoutInMillis"`
	// `CfnMethod.IntegrationProperty.Type`.
	Type *string `json:"type"`
	// `CfnMethod.IntegrationProperty.Uri`.
	Uri *string `json:"uri"`
}

type CfnMethod_IntegrationResponseProperty struct {
	// `CfnMethod.IntegrationResponseProperty.StatusCode`.
	StatusCode *string `json:"statusCode"`
	// `CfnMethod.IntegrationResponseProperty.ContentHandling`.
	ContentHandling *string `json:"contentHandling"`
	// `CfnMethod.IntegrationResponseProperty.ResponseParameters`.
	ResponseParameters interface{} `json:"responseParameters"`
	// `CfnMethod.IntegrationResponseProperty.ResponseTemplates`.
	ResponseTemplates interface{} `json:"responseTemplates"`
	// `CfnMethod.IntegrationResponseProperty.SelectionPattern`.
	SelectionPattern *string `json:"selectionPattern"`
}

type CfnMethod_MethodResponseProperty struct {
	// `CfnMethod.MethodResponseProperty.StatusCode`.
	StatusCode *string `json:"statusCode"`
	// `CfnMethod.MethodResponseProperty.ResponseModels`.
	ResponseModels interface{} `json:"responseModels"`
	// `CfnMethod.MethodResponseProperty.ResponseParameters`.
	ResponseParameters interface{} `json:"responseParameters"`
}

// Properties for defining a `AWS::ApiGateway::Method`.
type CfnMethodProps struct {
	// `AWS::ApiGateway::Method.HttpMethod`.
	HttpMethod *string `json:"httpMethod"`
	// `AWS::ApiGateway::Method.ResourceId`.
	ResourceId *string `json:"resourceId"`
	// `AWS::ApiGateway::Method.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::Method.ApiKeyRequired`.
	ApiKeyRequired interface{} `json:"apiKeyRequired"`
	// `AWS::ApiGateway::Method.AuthorizationScopes`.
	AuthorizationScopes *[]*string `json:"authorizationScopes"`
	// `AWS::ApiGateway::Method.AuthorizationType`.
	AuthorizationType *string `json:"authorizationType"`
	// `AWS::ApiGateway::Method.AuthorizerId`.
	AuthorizerId *string `json:"authorizerId"`
	// `AWS::ApiGateway::Method.Integration`.
	Integration interface{} `json:"integration"`
	// `AWS::ApiGateway::Method.MethodResponses`.
	MethodResponses interface{} `json:"methodResponses"`
	// `AWS::ApiGateway::Method.OperationName`.
	OperationName *string `json:"operationName"`
	// `AWS::ApiGateway::Method.RequestModels`.
	RequestModels interface{} `json:"requestModels"`
	// `AWS::ApiGateway::Method.RequestParameters`.
	RequestParameters interface{} `json:"requestParameters"`
	// `AWS::ApiGateway::Method.RequestValidatorId`.
	RequestValidatorId *string `json:"requestValidatorId"`
}

// A CloudFormation `AWS::ApiGateway::Model`.
type CfnModel interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ContentType() *string
	SetContentType(val *string)
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Schema() interface{}
	SetSchema(val interface{})
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

// The jsii proxy struct for CfnModel
type jsiiProxy_CfnModel struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnModel) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) ContentType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"contentType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) Schema() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"schema",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModel) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::Model`.
func NewCfnModel(scope awscdk.Construct, id *string, props *CfnModelProps) CfnModel {
	_init_.Initialize()

	j := jsiiProxy_CfnModel{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnModel",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::Model`.
func NewCfnModel_Override(c CfnModel, scope awscdk.Construct, id *string, props *CfnModelProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnModel",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnModel) SetContentType(val *string) {
	_jsii_.Set(
		j,
		"contentType",
		val,
	)
}

func (j *jsiiProxy_CfnModel) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnModel) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnModel) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

func (j *jsiiProxy_CfnModel) SetSchema(val interface{}) {
	_jsii_.Set(
		j,
		"schema",
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
func CfnModel_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnModel",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnModel_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnModel",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnModel_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnModel",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnModel_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnModel",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnModel) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnModel) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnModel) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnModel) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnModel) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnModel) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnModel) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnModel) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnModel) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnModel) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnModel) OnPrepare() {
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
func (c *jsiiProxy_CfnModel) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnModel) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnModel) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnModel) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnModel) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnModel) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnModel) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnModel) ToString() *string {
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
func (c *jsiiProxy_CfnModel) Validate() *[]*string {
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
func (c *jsiiProxy_CfnModel) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::Model`.
type CfnModelProps struct {
	// `AWS::ApiGateway::Model.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::Model.ContentType`.
	ContentType *string `json:"contentType"`
	// `AWS::ApiGateway::Model.Description`.
	Description *string `json:"description"`
	// `AWS::ApiGateway::Model.Name`.
	Name *string `json:"name"`
	// `AWS::ApiGateway::Model.Schema`.
	Schema interface{} `json:"schema"`
}

// A CloudFormation `AWS::ApiGatewayV2::Model`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-model.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnModelV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ContentType() *string
	SetContentType(val *string)
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Schema() interface{}
	SetSchema(val interface{})
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

// The jsii proxy struct for CfnModelV2
type jsiiProxy_CfnModelV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnModelV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) ContentType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"contentType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) Schema() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"schema",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnModelV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::Model`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnModelV2(scope awscdk.Construct, id *string, props *CfnModelV2Props) CfnModelV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnModelV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnModelV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::Model`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnModelV2_Override(c CfnModelV2, scope awscdk.Construct, id *string, props *CfnModelV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnModelV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnModelV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnModelV2) SetContentType(val *string) {
	_jsii_.Set(
		j,
		"contentType",
		val,
	)
}

func (j *jsiiProxy_CfnModelV2) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnModelV2) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnModelV2) SetSchema(val interface{}) {
	_jsii_.Set(
		j,
		"schema",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnModelV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnModelV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnModelV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnModelV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnModelV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnModelV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnModelV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnModelV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnModelV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGatewayV2::Model`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-model.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnModelV2Props struct {
	// `AWS::ApiGatewayV2::Model.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-model.html#cfn-apigatewayv2-model-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::Model.Name`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-model.html#cfn-apigatewayv2-model-name
	//
	// Deprecated: moved to package aws-apigatewayv2
	Name *string `json:"name"`
	// `AWS::ApiGatewayV2::Model.Schema`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-model.html#cfn-apigatewayv2-model-schema
	//
	// Deprecated: moved to package aws-apigatewayv2
	Schema interface{} `json:"schema"`
	// `AWS::ApiGatewayV2::Model.ContentType`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-model.html#cfn-apigatewayv2-model-contenttype
	//
	// Deprecated: moved to package aws-apigatewayv2
	ContentType *string `json:"contentType"`
	// `AWS::ApiGatewayV2::Model.Description`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-model.html#cfn-apigatewayv2-model-description
	//
	// Deprecated: moved to package aws-apigatewayv2
	Description *string `json:"description"`
}

// A CloudFormation `AWS::ApiGateway::RequestValidator`.
type CfnRequestValidator interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrRequestValidatorId() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	ValidateRequestBody() interface{}
	SetValidateRequestBody(val interface{})
	ValidateRequestParameters() interface{}
	SetValidateRequestParameters(val interface{})
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

// The jsii proxy struct for CfnRequestValidator
type jsiiProxy_CfnRequestValidator struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRequestValidator) AttrRequestValidatorId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRequestValidatorId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) ValidateRequestBody() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"validateRequestBody",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRequestValidator) ValidateRequestParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"validateRequestParameters",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::RequestValidator`.
func NewCfnRequestValidator(scope awscdk.Construct, id *string, props *CfnRequestValidatorProps) CfnRequestValidator {
	_init_.Initialize()

	j := jsiiProxy_CfnRequestValidator{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRequestValidator",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::RequestValidator`.
func NewCfnRequestValidator_Override(c CfnRequestValidator, scope awscdk.Construct, id *string, props *CfnRequestValidatorProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRequestValidator",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRequestValidator) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnRequestValidator) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

func (j *jsiiProxy_CfnRequestValidator) SetValidateRequestBody(val interface{}) {
	_jsii_.Set(
		j,
		"validateRequestBody",
		val,
	)
}

func (j *jsiiProxy_CfnRequestValidator) SetValidateRequestParameters(val interface{}) {
	_jsii_.Set(
		j,
		"validateRequestParameters",
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
func CfnRequestValidator_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRequestValidator",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnRequestValidator_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRequestValidator",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnRequestValidator_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRequestValidator",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRequestValidator_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnRequestValidator",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnRequestValidator) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnRequestValidator) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnRequestValidator) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnRequestValidator) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnRequestValidator) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnRequestValidator) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnRequestValidator) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnRequestValidator) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnRequestValidator) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnRequestValidator) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnRequestValidator) OnPrepare() {
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
func (c *jsiiProxy_CfnRequestValidator) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRequestValidator) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnRequestValidator) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnRequestValidator) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnRequestValidator) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnRequestValidator) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnRequestValidator) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRequestValidator) ToString() *string {
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
func (c *jsiiProxy_CfnRequestValidator) Validate() *[]*string {
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
func (c *jsiiProxy_CfnRequestValidator) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::RequestValidator`.
type CfnRequestValidatorProps struct {
	// `AWS::ApiGateway::RequestValidator.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::RequestValidator.Name`.
	Name *string `json:"name"`
	// `AWS::ApiGateway::RequestValidator.ValidateRequestBody`.
	ValidateRequestBody interface{} `json:"validateRequestBody"`
	// `AWS::ApiGateway::RequestValidator.ValidateRequestParameters`.
	ValidateRequestParameters interface{} `json:"validateRequestParameters"`
}

// A CloudFormation `AWS::ApiGateway::Resource`.
type CfnResource interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	Node() awscdk.ConstructNode
	ParentId() *string
	SetParentId(val *string)
	PathPart() *string
	SetPathPart(val *string)
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
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

// The jsii proxy struct for CfnResource
type jsiiProxy_CfnResource struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnResource) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) ParentId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"parentId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) PathPart() *string {
	var returns *string
	_jsii_.Get(
		j,
		"pathPart",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnResource) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::Resource`.
func NewCfnResource(scope awscdk.Construct, id *string, props *CfnResourceProps) CfnResource {
	_init_.Initialize()

	j := jsiiProxy_CfnResource{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnResource",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::Resource`.
func NewCfnResource_Override(c CfnResource, scope awscdk.Construct, id *string, props *CfnResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnResource",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnResource) SetParentId(val *string) {
	_jsii_.Set(
		j,
		"parentId",
		val,
	)
}

func (j *jsiiProxy_CfnResource) SetPathPart(val *string) {
	_jsii_.Set(
		j,
		"pathPart",
		val,
	)
}

func (j *jsiiProxy_CfnResource) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
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
func CfnResource_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnResource",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnResource_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnResource",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnResource_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnResource",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnResource_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnResource",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnResource) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnResource) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnResource) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnResource) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnResource) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnResource) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnResource) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnResource) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnResource) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnResource) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnResource) OnPrepare() {
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
func (c *jsiiProxy_CfnResource) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnResource) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnResource) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnResource) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnResource) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnResource) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnResource) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnResource) ToString() *string {
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
func (c *jsiiProxy_CfnResource) Validate() *[]*string {
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
func (c *jsiiProxy_CfnResource) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::Resource`.
type CfnResourceProps struct {
	// `AWS::ApiGateway::Resource.ParentId`.
	ParentId *string `json:"parentId"`
	// `AWS::ApiGateway::Resource.PathPart`.
	PathPart *string `json:"pathPart"`
	// `AWS::ApiGateway::Resource.RestApiId`.
	RestApiId *string `json:"restApiId"`
}

// A CloudFormation `AWS::ApiGateway::RestApi`.
type CfnRestApi interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiKeySourceType() *string
	SetApiKeySourceType(val *string)
	AttrRootResourceId() *string
	BinaryMediaTypes() *[]*string
	SetBinaryMediaTypes(val *[]*string)
	Body() interface{}
	SetBody(val interface{})
	BodyS3Location() interface{}
	SetBodyS3Location(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CloneFrom() *string
	SetCloneFrom(val *string)
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	DisableExecuteApiEndpoint() interface{}
	SetDisableExecuteApiEndpoint(val interface{})
	EndpointConfiguration() interface{}
	SetEndpointConfiguration(val interface{})
	FailOnWarnings() interface{}
	SetFailOnWarnings(val interface{})
	LogicalId() *string
	MinimumCompressionSize() *float64
	SetMinimumCompressionSize(val *float64)
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Parameters() interface{}
	SetParameters(val interface{})
	Policy() interface{}
	SetPolicy(val interface{})
	Ref() *string
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

// The jsii proxy struct for CfnRestApi
type jsiiProxy_CfnRestApi struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRestApi) ApiKeySourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiKeySourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) AttrRootResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRootResourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) BinaryMediaTypes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"binaryMediaTypes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Body() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"body",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) BodyS3Location() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"bodyS3Location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) CloneFrom() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cloneFrom",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) DisableExecuteApiEndpoint() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"disableExecuteApiEndpoint",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) EndpointConfiguration() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"endpointConfiguration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) FailOnWarnings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"failOnWarnings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) MinimumCompressionSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"minimumCompressionSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Parameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"parameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Policy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"policy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRestApi) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::RestApi`.
func NewCfnRestApi(scope awscdk.Construct, id *string, props *CfnRestApiProps) CfnRestApi {
	_init_.Initialize()

	j := jsiiProxy_CfnRestApi{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRestApi",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::RestApi`.
func NewCfnRestApi_Override(c CfnRestApi, scope awscdk.Construct, id *string, props *CfnRestApiProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRestApi",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRestApi) SetApiKeySourceType(val *string) {
	_jsii_.Set(
		j,
		"apiKeySourceType",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetBinaryMediaTypes(val *[]*string) {
	_jsii_.Set(
		j,
		"binaryMediaTypes",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetBody(val interface{}) {
	_jsii_.Set(
		j,
		"body",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetBodyS3Location(val interface{}) {
	_jsii_.Set(
		j,
		"bodyS3Location",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetCloneFrom(val *string) {
	_jsii_.Set(
		j,
		"cloneFrom",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetDisableExecuteApiEndpoint(val interface{}) {
	_jsii_.Set(
		j,
		"disableExecuteApiEndpoint",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetEndpointConfiguration(val interface{}) {
	_jsii_.Set(
		j,
		"endpointConfiguration",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetFailOnWarnings(val interface{}) {
	_jsii_.Set(
		j,
		"failOnWarnings",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetMinimumCompressionSize(val *float64) {
	_jsii_.Set(
		j,
		"minimumCompressionSize",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetParameters(val interface{}) {
	_jsii_.Set(
		j,
		"parameters",
		val,
	)
}

func (j *jsiiProxy_CfnRestApi) SetPolicy(val interface{}) {
	_jsii_.Set(
		j,
		"policy",
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
func CfnRestApi_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRestApi",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnRestApi_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRestApi",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnRestApi_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRestApi",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRestApi_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnRestApi",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnRestApi) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnRestApi) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnRestApi) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnRestApi) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnRestApi) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnRestApi) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnRestApi) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnRestApi) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnRestApi) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnRestApi) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnRestApi) OnPrepare() {
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
func (c *jsiiProxy_CfnRestApi) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRestApi) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnRestApi) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnRestApi) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnRestApi) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnRestApi) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnRestApi) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnRestApi) ToString() *string {
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
func (c *jsiiProxy_CfnRestApi) Validate() *[]*string {
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
func (c *jsiiProxy_CfnRestApi) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnRestApi_EndpointConfigurationProperty struct {
	// `CfnRestApi.EndpointConfigurationProperty.Types`.
	Types *[]*string `json:"types"`
	// `CfnRestApi.EndpointConfigurationProperty.VpcEndpointIds`.
	VpcEndpointIds *[]*string `json:"vpcEndpointIds"`
}

type CfnRestApi_S3LocationProperty struct {
	// `CfnRestApi.S3LocationProperty.Bucket`.
	Bucket *string `json:"bucket"`
	// `CfnRestApi.S3LocationProperty.ETag`.
	ETag *string `json:"eTag"`
	// `CfnRestApi.S3LocationProperty.Key`.
	Key *string `json:"key"`
	// `CfnRestApi.S3LocationProperty.Version`.
	Version *string `json:"version"`
}

// Properties for defining a `AWS::ApiGateway::RestApi`.
type CfnRestApiProps struct {
	// `AWS::ApiGateway::RestApi.ApiKeySourceType`.
	ApiKeySourceType *string `json:"apiKeySourceType"`
	// `AWS::ApiGateway::RestApi.BinaryMediaTypes`.
	BinaryMediaTypes *[]*string `json:"binaryMediaTypes"`
	// `AWS::ApiGateway::RestApi.Body`.
	Body interface{} `json:"body"`
	// `AWS::ApiGateway::RestApi.BodyS3Location`.
	BodyS3Location interface{} `json:"bodyS3Location"`
	// `AWS::ApiGateway::RestApi.CloneFrom`.
	CloneFrom *string `json:"cloneFrom"`
	// `AWS::ApiGateway::RestApi.Description`.
	Description *string `json:"description"`
	// `AWS::ApiGateway::RestApi.DisableExecuteApiEndpoint`.
	DisableExecuteApiEndpoint interface{} `json:"disableExecuteApiEndpoint"`
	// `AWS::ApiGateway::RestApi.EndpointConfiguration`.
	EndpointConfiguration interface{} `json:"endpointConfiguration"`
	// `AWS::ApiGateway::RestApi.FailOnWarnings`.
	FailOnWarnings interface{} `json:"failOnWarnings"`
	// `AWS::ApiGateway::RestApi.MinimumCompressionSize`.
	MinimumCompressionSize *float64 `json:"minimumCompressionSize"`
	// `AWS::ApiGateway::RestApi.Name`.
	Name *string `json:"name"`
	// `AWS::ApiGateway::RestApi.Parameters`.
	Parameters interface{} `json:"parameters"`
	// `AWS::ApiGateway::RestApi.Policy`.
	Policy interface{} `json:"policy"`
	// `AWS::ApiGateway::RestApi.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
}

// A CloudFormation `AWS::ApiGatewayV2::RouteResponse`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnRouteResponseV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	ModelSelectionExpression() *string
	SetModelSelectionExpression(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	ResponseModels() interface{}
	SetResponseModels(val interface{})
	ResponseParameters() interface{}
	SetResponseParameters(val interface{})
	RouteId() *string
	SetRouteId(val *string)
	RouteResponseKey() *string
	SetRouteResponseKey(val *string)
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

// The jsii proxy struct for CfnRouteResponseV2
type jsiiProxy_CfnRouteResponseV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRouteResponseV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) ModelSelectionExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"modelSelectionExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) ResponseModels() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"responseModels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) ResponseParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"responseParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) RouteId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"routeId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) RouteResponseKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"routeResponseKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteResponseV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::RouteResponse`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnRouteResponseV2(scope awscdk.Construct, id *string, props *CfnRouteResponseV2Props) CfnRouteResponseV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnRouteResponseV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRouteResponseV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::RouteResponse`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnRouteResponseV2_Override(c CfnRouteResponseV2, scope awscdk.Construct, id *string, props *CfnRouteResponseV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRouteResponseV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRouteResponseV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnRouteResponseV2) SetModelSelectionExpression(val *string) {
	_jsii_.Set(
		j,
		"modelSelectionExpression",
		val,
	)
}

func (j *jsiiProxy_CfnRouteResponseV2) SetResponseModels(val interface{}) {
	_jsii_.Set(
		j,
		"responseModels",
		val,
	)
}

func (j *jsiiProxy_CfnRouteResponseV2) SetResponseParameters(val interface{}) {
	_jsii_.Set(
		j,
		"responseParameters",
		val,
	)
}

func (j *jsiiProxy_CfnRouteResponseV2) SetRouteId(val *string) {
	_jsii_.Set(
		j,
		"routeId",
		val,
	)
}

func (j *jsiiProxy_CfnRouteResponseV2) SetRouteResponseKey(val *string) {
	_jsii_.Set(
		j,
		"routeResponseKey",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnRouteResponseV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRouteResponseV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnRouteResponseV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRouteResponseV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnRouteResponseV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRouteResponseV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRouteResponseV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnRouteResponseV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteResponseV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-routeresponse-parameterconstraints.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnRouteResponseV2_ParameterConstraintsProperty struct {
	// `CfnRouteResponseV2.ParameterConstraintsProperty.Required`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-routeresponse-parameterconstraints.html#cfn-apigatewayv2-routeresponse-parameterconstraints-required
	//
	// Deprecated: moved to package aws-apigatewayv2
	Required interface{} `json:"required"`
}

// Properties for defining a `AWS::ApiGatewayV2::RouteResponse`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnRouteResponseV2Props struct {
	// `AWS::ApiGatewayV2::RouteResponse.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html#cfn-apigatewayv2-routeresponse-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::RouteResponse.RouteId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html#cfn-apigatewayv2-routeresponse-routeid
	//
	// Deprecated: moved to package aws-apigatewayv2
	RouteId *string `json:"routeId"`
	// `AWS::ApiGatewayV2::RouteResponse.RouteResponseKey`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html#cfn-apigatewayv2-routeresponse-routeresponsekey
	//
	// Deprecated: moved to package aws-apigatewayv2
	RouteResponseKey *string `json:"routeResponseKey"`
	// `AWS::ApiGatewayV2::RouteResponse.ModelSelectionExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html#cfn-apigatewayv2-routeresponse-modelselectionexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	ModelSelectionExpression *string `json:"modelSelectionExpression"`
	// `AWS::ApiGatewayV2::RouteResponse.ResponseModels`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html#cfn-apigatewayv2-routeresponse-responsemodels
	//
	// Deprecated: moved to package aws-apigatewayv2
	ResponseModels interface{} `json:"responseModels"`
	// `AWS::ApiGatewayV2::RouteResponse.ResponseParameters`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-routeresponse.html#cfn-apigatewayv2-routeresponse-responseparameters
	//
	// Deprecated: moved to package aws-apigatewayv2
	ResponseParameters interface{} `json:"responseParameters"`
}

// A CloudFormation `AWS::ApiGatewayV2::Route`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnRouteV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiId() *string
	SetApiId(val *string)
	ApiKeyRequired() interface{}
	SetApiKeyRequired(val interface{})
	AuthorizationScopes() *[]*string
	SetAuthorizationScopes(val *[]*string)
	AuthorizationType() *string
	SetAuthorizationType(val *string)
	AuthorizerId() *string
	SetAuthorizerId(val *string)
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	LogicalId() *string
	ModelSelectionExpression() *string
	SetModelSelectionExpression(val *string)
	Node() awscdk.ConstructNode
	OperationName() *string
	SetOperationName(val *string)
	Ref() *string
	RequestModels() interface{}
	SetRequestModels(val interface{})
	RequestParameters() interface{}
	SetRequestParameters(val interface{})
	RouteKey() *string
	SetRouteKey(val *string)
	RouteResponseSelectionExpression() *string
	SetRouteResponseSelectionExpression(val *string)
	Stack() awscdk.Stack
	Target() *string
	SetTarget(val *string)
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

// The jsii proxy struct for CfnRouteV2
type jsiiProxy_CfnRouteV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnRouteV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) ApiKeyRequired() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"apiKeyRequired",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) AuthorizationScopes() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"authorizationScopes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) AuthorizationType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) AuthorizerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) ModelSelectionExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"modelSelectionExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) OperationName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"operationName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) RequestModels() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"requestModels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) RequestParameters() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"requestParameters",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) RouteKey() *string {
	var returns *string
	_jsii_.Get(
		j,
		"routeKey",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) RouteResponseSelectionExpression() *string {
	var returns *string
	_jsii_.Get(
		j,
		"routeResponseSelectionExpression",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) Target() *string {
	var returns *string
	_jsii_.Get(
		j,
		"target",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnRouteV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::Route`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnRouteV2(scope awscdk.Construct, id *string, props *CfnRouteV2Props) CfnRouteV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnRouteV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRouteV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::Route`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnRouteV2_Override(c CfnRouteV2, scope awscdk.Construct, id *string, props *CfnRouteV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnRouteV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetApiKeyRequired(val interface{}) {
	_jsii_.Set(
		j,
		"apiKeyRequired",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetAuthorizationScopes(val *[]*string) {
	_jsii_.Set(
		j,
		"authorizationScopes",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetAuthorizationType(val *string) {
	_jsii_.Set(
		j,
		"authorizationType",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetAuthorizerId(val *string) {
	_jsii_.Set(
		j,
		"authorizerId",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetModelSelectionExpression(val *string) {
	_jsii_.Set(
		j,
		"modelSelectionExpression",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetOperationName(val *string) {
	_jsii_.Set(
		j,
		"operationName",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetRequestModels(val interface{}) {
	_jsii_.Set(
		j,
		"requestModels",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetRequestParameters(val interface{}) {
	_jsii_.Set(
		j,
		"requestParameters",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetRouteKey(val *string) {
	_jsii_.Set(
		j,
		"routeKey",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetRouteResponseSelectionExpression(val *string) {
	_jsii_.Set(
		j,
		"routeResponseSelectionExpression",
		val,
	)
}

func (j *jsiiProxy_CfnRouteV2) SetTarget(val *string) {
	_jsii_.Set(
		j,
		"target",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnRouteV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRouteV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnRouteV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRouteV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnRouteV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnRouteV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnRouteV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnRouteV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnRouteV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-route-parameterconstraints.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnRouteV2_ParameterConstraintsProperty struct {
	// `CfnRouteV2.ParameterConstraintsProperty.Required`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-route-parameterconstraints.html#cfn-apigatewayv2-route-parameterconstraints-required
	//
	// Deprecated: moved to package aws-apigatewayv2
	Required interface{} `json:"required"`
}

// Properties for defining a `AWS::ApiGatewayV2::Route`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnRouteV2Props struct {
	// `AWS::ApiGatewayV2::Route.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::Route.RouteKey`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-routekey
	//
	// Deprecated: moved to package aws-apigatewayv2
	RouteKey *string `json:"routeKey"`
	// `AWS::ApiGatewayV2::Route.ApiKeyRequired`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-apikeyrequired
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiKeyRequired interface{} `json:"apiKeyRequired"`
	// `AWS::ApiGatewayV2::Route.AuthorizationScopes`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-authorizationscopes
	//
	// Deprecated: moved to package aws-apigatewayv2
	AuthorizationScopes *[]*string `json:"authorizationScopes"`
	// `AWS::ApiGatewayV2::Route.AuthorizationType`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-authorizationtype
	//
	// Deprecated: moved to package aws-apigatewayv2
	AuthorizationType *string `json:"authorizationType"`
	// `AWS::ApiGatewayV2::Route.AuthorizerId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-authorizerid
	//
	// Deprecated: moved to package aws-apigatewayv2
	AuthorizerId *string `json:"authorizerId"`
	// `AWS::ApiGatewayV2::Route.ModelSelectionExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-modelselectionexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	ModelSelectionExpression *string `json:"modelSelectionExpression"`
	// `AWS::ApiGatewayV2::Route.OperationName`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-operationname
	//
	// Deprecated: moved to package aws-apigatewayv2
	OperationName *string `json:"operationName"`
	// `AWS::ApiGatewayV2::Route.RequestModels`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-requestmodels
	//
	// Deprecated: moved to package aws-apigatewayv2
	RequestModels interface{} `json:"requestModels"`
	// `AWS::ApiGatewayV2::Route.RequestParameters`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-requestparameters
	//
	// Deprecated: moved to package aws-apigatewayv2
	RequestParameters interface{} `json:"requestParameters"`
	// `AWS::ApiGatewayV2::Route.RouteResponseSelectionExpression`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-routeresponseselectionexpression
	//
	// Deprecated: moved to package aws-apigatewayv2
	RouteResponseSelectionExpression *string `json:"routeResponseSelectionExpression"`
	// `AWS::ApiGatewayV2::Route.Target`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html#cfn-apigatewayv2-route-target
	//
	// Deprecated: moved to package aws-apigatewayv2
	Target *string `json:"target"`
}

// A CloudFormation `AWS::ApiGateway::Stage`.
type CfnStage interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AccessLogSetting() interface{}
	SetAccessLogSetting(val interface{})
	CacheClusterEnabled() interface{}
	SetCacheClusterEnabled(val interface{})
	CacheClusterSize() *string
	SetCacheClusterSize(val *string)
	CanarySetting() interface{}
	SetCanarySetting(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClientCertificateId() *string
	SetClientCertificateId(val *string)
	CreationStack() *[]*string
	DeploymentId() *string
	SetDeploymentId(val *string)
	Description() *string
	SetDescription(val *string)
	DocumentationVersion() *string
	SetDocumentationVersion(val *string)
	LogicalId() *string
	MethodSettings() interface{}
	SetMethodSettings(val interface{})
	Node() awscdk.ConstructNode
	Ref() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Stack() awscdk.Stack
	StageName() *string
	SetStageName(val *string)
	Tags() awscdk.TagManager
	TracingEnabled() interface{}
	SetTracingEnabled(val interface{})
	UpdatedProperites() *map[string]interface{}
	Variables() interface{}
	SetVariables(val interface{})
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

// The jsii proxy struct for CfnStage
type jsiiProxy_CfnStage struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnStage) AccessLogSetting() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accessLogSetting",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) CacheClusterEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"cacheClusterEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) CacheClusterSize() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cacheClusterSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) CanarySetting() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"canarySetting",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) ClientCertificateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clientCertificateId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) DeploymentId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deploymentId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) DocumentationVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"documentationVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) MethodSettings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"methodSettings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) TracingEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"tracingEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStage) Variables() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"variables",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::Stage`.
func NewCfnStage(scope awscdk.Construct, id *string, props *CfnStageProps) CfnStage {
	_init_.Initialize()

	j := jsiiProxy_CfnStage{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnStage",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::Stage`.
func NewCfnStage_Override(c CfnStage, scope awscdk.Construct, id *string, props *CfnStageProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnStage",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStage) SetAccessLogSetting(val interface{}) {
	_jsii_.Set(
		j,
		"accessLogSetting",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetCacheClusterEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"cacheClusterEnabled",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetCacheClusterSize(val *string) {
	_jsii_.Set(
		j,
		"cacheClusterSize",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetCanarySetting(val interface{}) {
	_jsii_.Set(
		j,
		"canarySetting",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetClientCertificateId(val *string) {
	_jsii_.Set(
		j,
		"clientCertificateId",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetDeploymentId(val *string) {
	_jsii_.Set(
		j,
		"deploymentId",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetDocumentationVersion(val *string) {
	_jsii_.Set(
		j,
		"documentationVersion",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetMethodSettings(val interface{}) {
	_jsii_.Set(
		j,
		"methodSettings",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetStageName(val *string) {
	_jsii_.Set(
		j,
		"stageName",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetTracingEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"tracingEnabled",
		val,
	)
}

func (j *jsiiProxy_CfnStage) SetVariables(val interface{}) {
	_jsii_.Set(
		j,
		"variables",
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
func CfnStage_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnStage",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnStage_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnStage",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnStage_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnStage",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStage_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnStage",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnStage) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnStage) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnStage) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnStage) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnStage) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnStage) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnStage) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnStage) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnStage) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnStage) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnStage) OnPrepare() {
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
func (c *jsiiProxy_CfnStage) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStage) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnStage) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnStage) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnStage) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnStage) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnStage) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnStage) ToString() *string {
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
func (c *jsiiProxy_CfnStage) Validate() *[]*string {
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
func (c *jsiiProxy_CfnStage) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnStage_AccessLogSettingProperty struct {
	// `CfnStage.AccessLogSettingProperty.DestinationArn`.
	DestinationArn *string `json:"destinationArn"`
	// `CfnStage.AccessLogSettingProperty.Format`.
	Format *string `json:"format"`
}

type CfnStage_CanarySettingProperty struct {
	// `CfnStage.CanarySettingProperty.DeploymentId`.
	DeploymentId *string `json:"deploymentId"`
	// `CfnStage.CanarySettingProperty.PercentTraffic`.
	PercentTraffic *float64 `json:"percentTraffic"`
	// `CfnStage.CanarySettingProperty.StageVariableOverrides`.
	StageVariableOverrides interface{} `json:"stageVariableOverrides"`
	// `CfnStage.CanarySettingProperty.UseStageCache`.
	UseStageCache interface{} `json:"useStageCache"`
}

type CfnStage_MethodSettingProperty struct {
	// `CfnStage.MethodSettingProperty.CacheDataEncrypted`.
	CacheDataEncrypted interface{} `json:"cacheDataEncrypted"`
	// `CfnStage.MethodSettingProperty.CacheTtlInSeconds`.
	CacheTtlInSeconds *float64 `json:"cacheTtlInSeconds"`
	// `CfnStage.MethodSettingProperty.CachingEnabled`.
	CachingEnabled interface{} `json:"cachingEnabled"`
	// `CfnStage.MethodSettingProperty.DataTraceEnabled`.
	DataTraceEnabled interface{} `json:"dataTraceEnabled"`
	// `CfnStage.MethodSettingProperty.HttpMethod`.
	HttpMethod *string `json:"httpMethod"`
	// `CfnStage.MethodSettingProperty.LoggingLevel`.
	LoggingLevel *string `json:"loggingLevel"`
	// `CfnStage.MethodSettingProperty.MetricsEnabled`.
	MetricsEnabled interface{} `json:"metricsEnabled"`
	// `CfnStage.MethodSettingProperty.ResourcePath`.
	ResourcePath *string `json:"resourcePath"`
	// `CfnStage.MethodSettingProperty.ThrottlingBurstLimit`.
	ThrottlingBurstLimit *float64 `json:"throttlingBurstLimit"`
	// `CfnStage.MethodSettingProperty.ThrottlingRateLimit`.
	ThrottlingRateLimit *float64 `json:"throttlingRateLimit"`
}

// Properties for defining a `AWS::ApiGateway::Stage`.
type CfnStageProps struct {
	// `AWS::ApiGateway::Stage.RestApiId`.
	RestApiId *string `json:"restApiId"`
	// `AWS::ApiGateway::Stage.AccessLogSetting`.
	AccessLogSetting interface{} `json:"accessLogSetting"`
	// `AWS::ApiGateway::Stage.CacheClusterEnabled`.
	CacheClusterEnabled interface{} `json:"cacheClusterEnabled"`
	// `AWS::ApiGateway::Stage.CacheClusterSize`.
	CacheClusterSize *string `json:"cacheClusterSize"`
	// `AWS::ApiGateway::Stage.CanarySetting`.
	CanarySetting interface{} `json:"canarySetting"`
	// `AWS::ApiGateway::Stage.ClientCertificateId`.
	ClientCertificateId *string `json:"clientCertificateId"`
	// `AWS::ApiGateway::Stage.DeploymentId`.
	DeploymentId *string `json:"deploymentId"`
	// `AWS::ApiGateway::Stage.Description`.
	Description *string `json:"description"`
	// `AWS::ApiGateway::Stage.DocumentationVersion`.
	DocumentationVersion *string `json:"documentationVersion"`
	// `AWS::ApiGateway::Stage.MethodSettings`.
	MethodSettings interface{} `json:"methodSettings"`
	// `AWS::ApiGateway::Stage.StageName`.
	StageName *string `json:"stageName"`
	// `AWS::ApiGateway::Stage.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::ApiGateway::Stage.TracingEnabled`.
	TracingEnabled interface{} `json:"tracingEnabled"`
	// `AWS::ApiGateway::Stage.Variables`.
	Variables interface{} `json:"variables"`
}

// A CloudFormation `AWS::ApiGatewayV2::Stage`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnStageV2 interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AccessLogSettings() interface{}
	SetAccessLogSettings(val interface{})
	ApiId() *string
	SetApiId(val *string)
	AutoDeploy() interface{}
	SetAutoDeploy(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	ClientCertificateId() *string
	SetClientCertificateId(val *string)
	CreationStack() *[]*string
	DefaultRouteSettings() interface{}
	SetDefaultRouteSettings(val interface{})
	DeploymentId() *string
	SetDeploymentId(val *string)
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	RouteSettings() interface{}
	SetRouteSettings(val interface{})
	Stack() awscdk.Stack
	StageName() *string
	SetStageName(val *string)
	StageVariables() interface{}
	SetStageVariables(val interface{})
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

// The jsii proxy struct for CfnStageV2
type jsiiProxy_CfnStageV2 struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnStageV2) AccessLogSettings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"accessLogSettings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) ApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"apiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) AutoDeploy() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"autoDeploy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) ClientCertificateId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"clientCertificateId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) DefaultRouteSettings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"defaultRouteSettings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) DeploymentId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deploymentId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) RouteSettings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"routeSettings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) StageVariables() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"stageVariables",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnStageV2) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGatewayV2::Stage`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnStageV2(scope awscdk.Construct, id *string, props *CfnStageV2Props) CfnStageV2 {
	_init_.Initialize()

	j := jsiiProxy_CfnStageV2{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnStageV2",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGatewayV2::Stage`.
// Deprecated: moved to package aws-apigatewayv2
func NewCfnStageV2_Override(c CfnStageV2, scope awscdk.Construct, id *string, props *CfnStageV2Props) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnStageV2",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnStageV2) SetAccessLogSettings(val interface{}) {
	_jsii_.Set(
		j,
		"accessLogSettings",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetApiId(val *string) {
	_jsii_.Set(
		j,
		"apiId",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetAutoDeploy(val interface{}) {
	_jsii_.Set(
		j,
		"autoDeploy",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetClientCertificateId(val *string) {
	_jsii_.Set(
		j,
		"clientCertificateId",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetDefaultRouteSettings(val interface{}) {
	_jsii_.Set(
		j,
		"defaultRouteSettings",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetDeploymentId(val *string) {
	_jsii_.Set(
		j,
		"deploymentId",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetRouteSettings(val interface{}) {
	_jsii_.Set(
		j,
		"routeSettings",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetStageName(val *string) {
	_jsii_.Set(
		j,
		"stageName",
		val,
	)
}

func (j *jsiiProxy_CfnStageV2) SetStageVariables(val interface{}) {
	_jsii_.Set(
		j,
		"stageVariables",
		val,
	)
}

// Returns `true` if a construct is a stack element (i.e. part of the synthesized cloudformation template).
//
// Uses duck-typing instead of `instanceof` to allow stack elements from different
// versions of this library to be included in the same stack.
//
// Returns: The construct as a stack element or undefined if it is not a stack element.
// Deprecated: moved to package aws-apigatewayv2
func CfnStageV2_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnStageV2",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Deprecated: moved to package aws-apigatewayv2
func CfnStageV2_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnStageV2",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Deprecated: moved to package aws-apigatewayv2
func CfnStageV2_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnStageV2",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnStageV2_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnStageV2",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) AddDeletionOverride(path *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) AddDependsOn(target awscdk.CfnResource) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) AddMetadata(key *string, value interface{}) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) AddPropertyDeletionOverride(propertyPath *string) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyDeletionOverride",
		[]interface{}{propertyPath},
	)
}

// Adds an override to a resource property.
//
// Syntactic sugar for `addOverride("Properties.<...>", value)`.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) GetAtt(attributeName *string) awscdk.Reference {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) GetMetadata(key *string) interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) Inspect(inspector awscdk.TreeInspector) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) OnPrepare() {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) OnSynthesize(session constructs.ISynthesisSession) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) OnValidate() *[]*string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) OverrideLogicalId(newLogicalId *string) {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) ShouldSynthesize() *bool {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
//
// Returns: a string representation of this resource
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) ToString() *string {
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
// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Deprecated: moved to package aws-apigatewayv2
func (c *jsiiProxy_CfnStageV2) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-accesslogsettings.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnStageV2_AccessLogSettingsProperty struct {
	// `CfnStageV2.AccessLogSettingsProperty.DestinationArn`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-accesslogsettings.html#cfn-apigatewayv2-stage-accesslogsettings-destinationarn
	//
	// Deprecated: moved to package aws-apigatewayv2
	DestinationArn *string `json:"destinationArn"`
	// `CfnStageV2.AccessLogSettingsProperty.Format`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-accesslogsettings.html#cfn-apigatewayv2-stage-accesslogsettings-format
	//
	// Deprecated: moved to package aws-apigatewayv2
	Format *string `json:"format"`
}

// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnStageV2_RouteSettingsProperty struct {
	// `CfnStageV2.RouteSettingsProperty.DataTraceEnabled`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-datatraceenabled
	//
	// Deprecated: moved to package aws-apigatewayv2
	DataTraceEnabled interface{} `json:"dataTraceEnabled"`
	// `CfnStageV2.RouteSettingsProperty.DetailedMetricsEnabled`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-detailedmetricsenabled
	//
	// Deprecated: moved to package aws-apigatewayv2
	DetailedMetricsEnabled interface{} `json:"detailedMetricsEnabled"`
	// `CfnStageV2.RouteSettingsProperty.LoggingLevel`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-logginglevel
	//
	// Deprecated: moved to package aws-apigatewayv2
	LoggingLevel *string `json:"loggingLevel"`
	// `CfnStageV2.RouteSettingsProperty.ThrottlingBurstLimit`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-throttlingburstlimit
	//
	// Deprecated: moved to package aws-apigatewayv2
	ThrottlingBurstLimit *float64 `json:"throttlingBurstLimit"`
	// `CfnStageV2.RouteSettingsProperty.ThrottlingRateLimit`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-throttlingratelimit
	//
	// Deprecated: moved to package aws-apigatewayv2
	ThrottlingRateLimit *float64 `json:"throttlingRateLimit"`
}

// Properties for defining a `AWS::ApiGatewayV2::Stage`.
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html
//
// Deprecated: moved to package aws-apigatewayv2
type CfnStageV2Props struct {
	// `AWS::ApiGatewayV2::Stage.ApiId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-apiid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ApiId *string `json:"apiId"`
	// `AWS::ApiGatewayV2::Stage.StageName`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-stagename
	//
	// Deprecated: moved to package aws-apigatewayv2
	StageName *string `json:"stageName"`
	// `AWS::ApiGatewayV2::Stage.AccessLogSettings`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-accesslogsettings
	//
	// Deprecated: moved to package aws-apigatewayv2
	AccessLogSettings interface{} `json:"accessLogSettings"`
	// `AWS::ApiGatewayV2::Stage.AutoDeploy`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-autodeploy
	//
	// Deprecated: moved to package aws-apigatewayv2
	AutoDeploy interface{} `json:"autoDeploy"`
	// `AWS::ApiGatewayV2::Stage.ClientCertificateId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-clientcertificateid
	//
	// Deprecated: moved to package aws-apigatewayv2
	ClientCertificateId *string `json:"clientCertificateId"`
	// `AWS::ApiGatewayV2::Stage.DefaultRouteSettings`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-defaultroutesettings
	//
	// Deprecated: moved to package aws-apigatewayv2
	DefaultRouteSettings interface{} `json:"defaultRouteSettings"`
	// `AWS::ApiGatewayV2::Stage.DeploymentId`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-deploymentid
	//
	// Deprecated: moved to package aws-apigatewayv2
	DeploymentId *string `json:"deploymentId"`
	// `AWS::ApiGatewayV2::Stage.Description`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-description
	//
	// Deprecated: moved to package aws-apigatewayv2
	Description *string `json:"description"`
	// `AWS::ApiGatewayV2::Stage.RouteSettings`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-routesettings
	//
	// Deprecated: moved to package aws-apigatewayv2
	RouteSettings interface{} `json:"routeSettings"`
	// `AWS::ApiGatewayV2::Stage.StageVariables`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-stagevariables
	//
	// Deprecated: moved to package aws-apigatewayv2
	StageVariables interface{} `json:"stageVariables"`
	// `AWS::ApiGatewayV2::Stage.Tags`.
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html#cfn-apigatewayv2-stage-tags
	//
	// Deprecated: moved to package aws-apigatewayv2
	Tags interface{} `json:"tags"`
}

// A CloudFormation `AWS::ApiGateway::UsagePlan`.
type CfnUsagePlan interface {
	awscdk.CfnResource
	awscdk.IInspectable
	ApiStages() interface{}
	SetApiStages(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Quota() interface{}
	SetQuota(val interface{})
	Ref() *string
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	Throttle() interface{}
	SetThrottle(val interface{})
	UpdatedProperites() *map[string]interface{}
	UsagePlanName() *string
	SetUsagePlanName(val *string)
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

// The jsii proxy struct for CfnUsagePlan
type jsiiProxy_CfnUsagePlan struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUsagePlan) ApiStages() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"apiStages",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) Quota() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"quota",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) Throttle() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"throttle",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlan) UsagePlanName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"usagePlanName",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::UsagePlan`.
func NewCfnUsagePlan(scope awscdk.Construct, id *string, props *CfnUsagePlanProps) CfnUsagePlan {
	_init_.Initialize()

	j := jsiiProxy_CfnUsagePlan{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnUsagePlan",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::UsagePlan`.
func NewCfnUsagePlan_Override(c CfnUsagePlan, scope awscdk.Construct, id *string, props *CfnUsagePlanProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnUsagePlan",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUsagePlan) SetApiStages(val interface{}) {
	_jsii_.Set(
		j,
		"apiStages",
		val,
	)
}

func (j *jsiiProxy_CfnUsagePlan) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnUsagePlan) SetQuota(val interface{}) {
	_jsii_.Set(
		j,
		"quota",
		val,
	)
}

func (j *jsiiProxy_CfnUsagePlan) SetThrottle(val interface{}) {
	_jsii_.Set(
		j,
		"throttle",
		val,
	)
}

func (j *jsiiProxy_CfnUsagePlan) SetUsagePlanName(val *string) {
	_jsii_.Set(
		j,
		"usagePlanName",
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
func CfnUsagePlan_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnUsagePlan",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUsagePlan_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnUsagePlan",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUsagePlan_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnUsagePlan",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUsagePlan_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnUsagePlan",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUsagePlan) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUsagePlan) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUsagePlan) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUsagePlan) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUsagePlan) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUsagePlan) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUsagePlan) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUsagePlan) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUsagePlan) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUsagePlan) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUsagePlan) OnPrepare() {
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
func (c *jsiiProxy_CfnUsagePlan) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUsagePlan) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUsagePlan) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUsagePlan) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUsagePlan) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUsagePlan) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUsagePlan) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUsagePlan) ToString() *string {
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
func (c *jsiiProxy_CfnUsagePlan) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUsagePlan) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnUsagePlan_ApiStageProperty struct {
	// `CfnUsagePlan.ApiStageProperty.ApiId`.
	ApiId *string `json:"apiId"`
	// `CfnUsagePlan.ApiStageProperty.Stage`.
	Stage *string `json:"stage"`
	// `CfnUsagePlan.ApiStageProperty.Throttle`.
	Throttle interface{} `json:"throttle"`
}

type CfnUsagePlan_QuotaSettingsProperty struct {
	// `CfnUsagePlan.QuotaSettingsProperty.Limit`.
	Limit *float64 `json:"limit"`
	// `CfnUsagePlan.QuotaSettingsProperty.Offset`.
	Offset *float64 `json:"offset"`
	// `CfnUsagePlan.QuotaSettingsProperty.Period`.
	Period *string `json:"period"`
}

type CfnUsagePlan_ThrottleSettingsProperty struct {
	// `CfnUsagePlan.ThrottleSettingsProperty.BurstLimit`.
	BurstLimit *float64 `json:"burstLimit"`
	// `CfnUsagePlan.ThrottleSettingsProperty.RateLimit`.
	RateLimit *float64 `json:"rateLimit"`
}

// A CloudFormation `AWS::ApiGateway::UsagePlanKey`.
type CfnUsagePlanKey interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	KeyId() *string
	SetKeyId(val *string)
	KeyType() *string
	SetKeyType(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	UpdatedProperites() *map[string]interface{}
	UsagePlanId() *string
	SetUsagePlanId(val *string)
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

// The jsii proxy struct for CfnUsagePlanKey
type jsiiProxy_CfnUsagePlanKey struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnUsagePlanKey) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) KeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) KeyType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnUsagePlanKey) UsagePlanId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"usagePlanId",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::UsagePlanKey`.
func NewCfnUsagePlanKey(scope awscdk.Construct, id *string, props *CfnUsagePlanKeyProps) CfnUsagePlanKey {
	_init_.Initialize()

	j := jsiiProxy_CfnUsagePlanKey{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnUsagePlanKey",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::UsagePlanKey`.
func NewCfnUsagePlanKey_Override(c CfnUsagePlanKey, scope awscdk.Construct, id *string, props *CfnUsagePlanKeyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnUsagePlanKey",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnUsagePlanKey) SetKeyId(val *string) {
	_jsii_.Set(
		j,
		"keyId",
		val,
	)
}

func (j *jsiiProxy_CfnUsagePlanKey) SetKeyType(val *string) {
	_jsii_.Set(
		j,
		"keyType",
		val,
	)
}

func (j *jsiiProxy_CfnUsagePlanKey) SetUsagePlanId(val *string) {
	_jsii_.Set(
		j,
		"usagePlanId",
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
func CfnUsagePlanKey_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnUsagePlanKey",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnUsagePlanKey_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnUsagePlanKey",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnUsagePlanKey_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnUsagePlanKey",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnUsagePlanKey_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnUsagePlanKey",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnUsagePlanKey) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnUsagePlanKey) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnUsagePlanKey) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnUsagePlanKey) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnUsagePlanKey) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnUsagePlanKey) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnUsagePlanKey) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnUsagePlanKey) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnUsagePlanKey) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnUsagePlanKey) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnUsagePlanKey) OnPrepare() {
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
func (c *jsiiProxy_CfnUsagePlanKey) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUsagePlanKey) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnUsagePlanKey) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnUsagePlanKey) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnUsagePlanKey) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnUsagePlanKey) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnUsagePlanKey) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnUsagePlanKey) ToString() *string {
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
func (c *jsiiProxy_CfnUsagePlanKey) Validate() *[]*string {
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
func (c *jsiiProxy_CfnUsagePlanKey) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::UsagePlanKey`.
type CfnUsagePlanKeyProps struct {
	// `AWS::ApiGateway::UsagePlanKey.KeyId`.
	KeyId *string `json:"keyId"`
	// `AWS::ApiGateway::UsagePlanKey.KeyType`.
	KeyType *string `json:"keyType"`
	// `AWS::ApiGateway::UsagePlanKey.UsagePlanId`.
	UsagePlanId *string `json:"usagePlanId"`
}

// Properties for defining a `AWS::ApiGateway::UsagePlan`.
type CfnUsagePlanProps struct {
	// `AWS::ApiGateway::UsagePlan.ApiStages`.
	ApiStages interface{} `json:"apiStages"`
	// `AWS::ApiGateway::UsagePlan.Description`.
	Description *string `json:"description"`
	// `AWS::ApiGateway::UsagePlan.Quota`.
	Quota interface{} `json:"quota"`
	// `AWS::ApiGateway::UsagePlan.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::ApiGateway::UsagePlan.Throttle`.
	Throttle interface{} `json:"throttle"`
	// `AWS::ApiGateway::UsagePlan.UsagePlanName`.
	UsagePlanName *string `json:"usagePlanName"`
}

// A CloudFormation `AWS::ApiGateway::VpcLink`.
type CfnVpcLink interface {
	awscdk.CfnResource
	awscdk.IInspectable
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	Description() *string
	SetDescription(val *string)
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Stack() awscdk.Stack
	TargetArns() *[]*string
	SetTargetArns(val *[]*string)
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

// The jsii proxy struct for CfnVpcLink
type jsiiProxy_CfnVpcLink struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnVpcLink) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) Description() *string {
	var returns *string
	_jsii_.Get(
		j,
		"description",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) TargetArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"targetArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnVpcLink) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ApiGateway::VpcLink`.
func NewCfnVpcLink(scope awscdk.Construct, id *string, props *CfnVpcLinkProps) CfnVpcLink {
	_init_.Initialize()

	j := jsiiProxy_CfnVpcLink{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnVpcLink",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ApiGateway::VpcLink`.
func NewCfnVpcLink_Override(c CfnVpcLink, scope awscdk.Construct, id *string, props *CfnVpcLinkProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CfnVpcLink",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnVpcLink) SetDescription(val *string) {
	_jsii_.Set(
		j,
		"description",
		val,
	)
}

func (j *jsiiProxy_CfnVpcLink) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnVpcLink) SetTargetArns(val *[]*string) {
	_jsii_.Set(
		j,
		"targetArns",
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
func CfnVpcLink_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnVpcLink",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnVpcLink_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnVpcLink",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnVpcLink_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CfnVpcLink",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnVpcLink_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.CfnVpcLink",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnVpcLink) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnVpcLink) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnVpcLink) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnVpcLink) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnVpcLink) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnVpcLink) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnVpcLink) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnVpcLink) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnVpcLink) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnVpcLink) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnVpcLink) OnPrepare() {
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
func (c *jsiiProxy_CfnVpcLink) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnVpcLink) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnVpcLink) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnVpcLink) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnVpcLink) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnVpcLink) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnVpcLink) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnVpcLink) ToString() *string {
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
func (c *jsiiProxy_CfnVpcLink) Validate() *[]*string {
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
func (c *jsiiProxy_CfnVpcLink) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

// Properties for defining a `AWS::ApiGateway::VpcLink`.
type CfnVpcLinkProps struct {
	// `AWS::ApiGateway::VpcLink.Name`.
	Name *string `json:"name"`
	// `AWS::ApiGateway::VpcLink.TargetArns`.
	TargetArns *[]*string `json:"targetArns"`
	// `AWS::ApiGateway::VpcLink.Description`.
	Description *string `json:"description"`
}

// Cognito user pools based custom authorizer.
// Experimental.
type CognitoUserPoolsAuthorizer interface {
	Authorizer
	IAuthorizer
	AuthorizationType() AuthorizationType
	AuthorizerArn() *string
	AuthorizerId() *string
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
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

// The jsii proxy struct for CognitoUserPoolsAuthorizer
type jsiiProxy_CognitoUserPoolsAuthorizer struct {
	jsiiProxy_Authorizer
	jsiiProxy_IAuthorizer
}

func (j *jsiiProxy_CognitoUserPoolsAuthorizer) AuthorizationType() AuthorizationType {
	var returns AuthorizationType
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CognitoUserPoolsAuthorizer) AuthorizerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CognitoUserPoolsAuthorizer) AuthorizerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CognitoUserPoolsAuthorizer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CognitoUserPoolsAuthorizer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CognitoUserPoolsAuthorizer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CognitoUserPoolsAuthorizer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewCognitoUserPoolsAuthorizer(scope constructs.Construct, id *string, props *CognitoUserPoolsAuthorizerProps) CognitoUserPoolsAuthorizer {
	_init_.Initialize()

	j := jsiiProxy_CognitoUserPoolsAuthorizer{}

	_jsii_.Create(
		"monocdk.aws_apigateway.CognitoUserPoolsAuthorizer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewCognitoUserPoolsAuthorizer_Override(c CognitoUserPoolsAuthorizer, scope constructs.Construct, id *string, props *CognitoUserPoolsAuthorizerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.CognitoUserPoolsAuthorizer",
		[]interface{}{scope, id, props},
		c,
	)
}

// Return whether the given object is an Authorizer.
// Experimental.
func CognitoUserPoolsAuthorizer_IsAuthorizer(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CognitoUserPoolsAuthorizer",
		"isAuthorizer",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CognitoUserPoolsAuthorizer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CognitoUserPoolsAuthorizer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func CognitoUserPoolsAuthorizer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.CognitoUserPoolsAuthorizer",
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		c,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) OnPrepare() {
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		c,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) ToString() *string {
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
func (c *jsiiProxy_CognitoUserPoolsAuthorizer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		c,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for CognitoUserPoolsAuthorizer.
// Experimental.
type CognitoUserPoolsAuthorizerProps struct {
	// The user pools to associate with this authorizer.
	// Experimental.
	CognitoUserPools *[]awscognito.IUserPool `json:"cognitoUserPools"`
	// An optional human friendly name for the authorizer.
	//
	// Note that, this is not the primary identifier of the authorizer.
	// Experimental.
	AuthorizerName *string `json:"authorizerName"`
	// The request header mapping expression for the bearer token.
	//
	// This is typically passed as part of the header, in which case
	// this should be `method.request.header.Authorizer` where Authorizer is the header containing the bearer token.
	// See: https://docs.aws.amazon.com/apigateway/api-reference/link-relation/authorizer-create/#identitySource
	//
	// Experimental.
	IdentitySource *string `json:"identitySource"`
	// How long APIGateway should cache the results.
	//
	// Max 1 hour.
	// Disable caching by setting this to 0.
	// Experimental.
	ResultsCacheTtl awscdk.Duration `json:"resultsCacheTtl"`
}

// Experimental.
type ConnectionType string

const (
	ConnectionType_INTERNET ConnectionType = "INTERNET"
	ConnectionType_VPC_LINK ConnectionType = "VPC_LINK"
)

// Experimental.
type ContentHandling string

const (
	ContentHandling_CONVERT_TO_BINARY ContentHandling = "CONVERT_TO_BINARY"
	ContentHandling_CONVERT_TO_TEXT ContentHandling = "CONVERT_TO_TEXT"
)

// Experimental.
type Cors interface {
}

// The jsii proxy struct for Cors
type jsiiProxy_Cors struct {
	_ byte // padding
}

func Cors_ALL_METHODS() *[]*string {
	_init_.Initialize()
	var returns *[]*string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.Cors",
		"ALL_METHODS",
		&returns,
	)
	return returns
}

func Cors_ALL_ORIGINS() *[]*string {
	_init_.Initialize()
	var returns *[]*string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.Cors",
		"ALL_ORIGINS",
		&returns,
	)
	return returns
}

func Cors_DEFAULT_HEADERS() *[]*string {
	_init_.Initialize()
	var returns *[]*string
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.Cors",
		"DEFAULT_HEADERS",
		&returns,
	)
	return returns
}

// Experimental.
type CorsOptions struct {
	// Specifies the list of origins that are allowed to make requests to this resource.
	//
	// If you wish to allow all origins, specify `Cors.ALL_ORIGINS` or
	// `[ * ]`.
	//
	// Responses will include the `Access-Control-Allow-Origin` response header.
	// If `Cors.ALL_ORIGINS` is specified, the `Vary: Origin` response header will
	// also be included.
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	//
	// Experimental.
	AllowOrigins *[]*string `json:"allowOrigins"`
	// The Access-Control-Allow-Credentials response header tells browsers whether to expose the response to frontend JavaScript code when the request's credentials mode (Request.credentials) is "include".
	//
	// When a request's credentials mode (Request.credentials) is "include",
	// browsers will only expose the response to frontend JavaScript code if the
	// Access-Control-Allow-Credentials value is true.
	//
	// Credentials are cookies, authorization headers or TLS client certificates.
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	//
	// Experimental.
	AllowCredentials *bool `json:"allowCredentials"`
	// The Access-Control-Allow-Headers response header is used in response to a preflight request which includes the Access-Control-Request-Headers to indicate which HTTP headers can be used during the actual request.
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	//
	// Experimental.
	AllowHeaders *[]*string `json:"allowHeaders"`
	// The Access-Control-Allow-Methods response header specifies the method or methods allowed when accessing the resource in response to a preflight request.
	//
	// If `ANY` is specified, it will be expanded to `Cors.ALL_METHODS`.
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	//
	// Experimental.
	AllowMethods *[]*string `json:"allowMethods"`
	// Sets Access-Control-Max-Age to -1, which means that caching is disabled.
	//
	// This option cannot be used with `maxAge`.
	// Experimental.
	DisableCache *bool `json:"disableCache"`
	// The Access-Control-Expose-Headers response header indicates which headers can be exposed as part of the response by listing their names.
	//
	// If you want clients to be able to access other headers, you have to list
	// them using the Access-Control-Expose-Headers header.
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	//
	// Experimental.
	ExposeHeaders *[]*string `json:"exposeHeaders"`
	// The Access-Control-Max-Age response header indicates how long the results of a preflight request (that is the information contained in the Access-Control-Allow-Methods and Access-Control-Allow-Headers headers) can be cached.
	//
	// To disable caching altogether use `disableCache: true`.
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	//
	// Experimental.
	MaxAge awscdk.Duration `json:"maxAge"`
	// Specifies the response status code returned from the OPTIONS method.
	// Experimental.
	StatusCode *float64 `json:"statusCode"`
}

// A Deployment of a REST API.
//
// An immutable representation of a RestApi resource that can be called by users
// using Stages. A deployment must be associated with a Stage for it to be
// callable over the Internet.
//
// Normally, you don't need to define deployments manually. The RestApi
// construct manages a Deployment resource that represents the latest model. It
// can be accessed through `restApi.latestDeployment` (unless `deploy: false` is
// set when defining the `RestApi`).
//
// If you manually define this resource, you will need to know that since
// deployments are immutable, as long as the resource's logical ID doesn't
// change, the deployment will represent the snapshot in time in which the
// resource was created. This means that if you modify the RestApi model (i.e.
// add methods or resources), these changes will not be reflected unless a new
// deployment resource is created.
//
// To achieve this behavior, the method `addToLogicalId(data)` can be used to
// augment the logical ID generated for the deployment resource such that it
// will include arbitrary data. This is done automatically for the
// `restApi.latestDeployment` deployment.
//
// Furthermore, since a deployment does not reference any of the REST API
// resources and methods, CloudFormation will likely provision it before these
// resources are created, which means that it will represent a "half-baked"
// model. Use the `node.addDependency(dep)` method to circumvent that. This is done
// automatically for the `restApi.latestDeployment` deployment.
// Experimental.
type Deployment interface {
	awscdk.Resource
	Api() IRestApi
	DeploymentId() *string
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddToLogicalId(data interface{})
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

// The jsii proxy struct for Deployment
type jsiiProxy_Deployment struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_Deployment) Api() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"api",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Deployment) DeploymentId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"deploymentId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Deployment) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Deployment) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Deployment) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Deployment) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewDeployment(scope constructs.Construct, id *string, props *DeploymentProps) Deployment {
	_init_.Initialize()

	j := jsiiProxy_Deployment{}

	_jsii_.Create(
		"monocdk.aws_apigateway.Deployment",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewDeployment_Override(d Deployment, scope constructs.Construct, id *string, props *DeploymentProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.Deployment",
		[]interface{}{scope, id, props},
		d,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Deployment_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Deployment",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Deployment_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Deployment",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds a component to the hash that determines this Deployment resource's logical ID.
//
// This should be called by constructs of the API Gateway model that want to
// invalidate the deployment when their settings change. The component will
// be resolve()ed during synthesis so tokens are welcome.
// Experimental.
func (d *jsiiProxy_Deployment) AddToLogicalId(data interface{}) {
	_jsii_.InvokeVoid(
		d,
		"addToLogicalId",
		[]interface{}{data},
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
func (d *jsiiProxy_Deployment) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		d,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (d *jsiiProxy_Deployment) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Deployment) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Deployment) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_Deployment) OnPrepare() {
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
func (d *jsiiProxy_Deployment) OnSynthesize(session constructs.ISynthesisSession) {
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
func (d *jsiiProxy_Deployment) OnValidate() *[]*string {
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
func (d *jsiiProxy_Deployment) Prepare() {
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
func (d *jsiiProxy_Deployment) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (d *jsiiProxy_Deployment) ToString() *string {
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
func (d *jsiiProxy_Deployment) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		d,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type DeploymentProps struct {
	// The Rest API to deploy.
	// Experimental.
	Api IRestApi `json:"api"`
	// A description of the purpose of the API Gateway deployment.
	// Experimental.
	Description *string `json:"description"`
	// When an API Gateway model is updated, a new deployment will automatically be created.
	//
	// If this is true, the old API Gateway Deployment resource will not be deleted.
	// This will allow manually reverting back to a previous deployment in case for example
	// Experimental.
	RetainDeployments *bool `json:"retainDeployments"`
}

// Experimental.
type DomainName interface {
	awscdk.Resource
	IDomainName
	DomainName() *string
	DomainNameAliasDomainName() *string
	DomainNameAliasHostedZoneId() *string
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddBasePathMapping(targetApi IRestApi, options *BasePathMappingOptions) BasePathMapping
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

// The jsii proxy struct for DomainName
type jsiiProxy_DomainName struct {
	internal.Type__awscdkResource
	jsiiProxy_IDomainName
}

func (j *jsiiProxy_DomainName) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DomainName) DomainNameAliasDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainNameAliasDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DomainName) DomainNameAliasHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainNameAliasHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DomainName) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DomainName) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DomainName) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DomainName) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewDomainName(scope constructs.Construct, id *string, props *DomainNameProps) DomainName {
	_init_.Initialize()

	j := jsiiProxy_DomainName{}

	_jsii_.Create(
		"monocdk.aws_apigateway.DomainName",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewDomainName_Override(d DomainName, scope constructs.Construct, id *string, props *DomainNameProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.DomainName",
		[]interface{}{scope, id, props},
		d,
	)
}

// Imports an existing domain name.
// Experimental.
func DomainName_FromDomainNameAttributes(scope constructs.Construct, id *string, attrs *DomainNameAttributes) IDomainName {
	_init_.Initialize()

	var returns IDomainName

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.DomainName",
		"fromDomainNameAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func DomainName_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.DomainName",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func DomainName_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.DomainName",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Maps this domain to an API endpoint.
// Experimental.
func (d *jsiiProxy_DomainName) AddBasePathMapping(targetApi IRestApi, options *BasePathMappingOptions) BasePathMapping {
	var returns BasePathMapping

	_jsii_.Invoke(
		d,
		"addBasePathMapping",
		[]interface{}{targetApi, options},
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
func (d *jsiiProxy_DomainName) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		d,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (d *jsiiProxy_DomainName) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_DomainName) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_DomainName) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		d,
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
func (d *jsiiProxy_DomainName) OnPrepare() {
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
func (d *jsiiProxy_DomainName) OnSynthesize(session constructs.ISynthesisSession) {
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
func (d *jsiiProxy_DomainName) OnValidate() *[]*string {
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
func (d *jsiiProxy_DomainName) Prepare() {
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
func (d *jsiiProxy_DomainName) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		d,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (d *jsiiProxy_DomainName) ToString() *string {
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
func (d *jsiiProxy_DomainName) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		d,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type DomainNameAttributes struct {
	// The domain name (e.g. `example.com`).
	// Experimental.
	DomainName *string `json:"domainName"`
	// The Route53 hosted zone ID to use in order to connect a record set to this domain through an alias.
	// Experimental.
	DomainNameAliasHostedZoneId *string `json:"domainNameAliasHostedZoneId"`
	// The Route53 alias target to use in order to connect a record set to this domain through an alias.
	// Experimental.
	DomainNameAliasTarget *string `json:"domainNameAliasTarget"`
}

// Experimental.
type DomainNameOptions struct {
	// The reference to an AWS-managed certificate for use by the edge-optimized endpoint for the domain name.
	//
	// For "EDGE" domain names, the certificate
	// needs to be in the US East (N. Virginia) region.
	// Experimental.
	Certificate awscertificatemanager.ICertificate `json:"certificate"`
	// The custom domain name for your API.
	//
	// Uppercase letters are not supported.
	// Experimental.
	DomainName *string `json:"domainName"`
	// The type of endpoint for this DomainName.
	// Experimental.
	EndpointType EndpointType `json:"endpointType"`
	// The mutual TLS authentication configuration for a custom domain name.
	// Experimental.
	Mtls *MTLSConfig `json:"mtls"`
	// The Transport Layer Security (TLS) version + cipher suite for this domain name.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-domainname.html
	//
	// Experimental.
	SecurityPolicy SecurityPolicy `json:"securityPolicy"`
}

// Experimental.
type DomainNameProps struct {
	// The reference to an AWS-managed certificate for use by the edge-optimized endpoint for the domain name.
	//
	// For "EDGE" domain names, the certificate
	// needs to be in the US East (N. Virginia) region.
	// Experimental.
	Certificate awscertificatemanager.ICertificate `json:"certificate"`
	// The custom domain name for your API.
	//
	// Uppercase letters are not supported.
	// Experimental.
	DomainName *string `json:"domainName"`
	// The type of endpoint for this DomainName.
	// Experimental.
	EndpointType EndpointType `json:"endpointType"`
	// The mutual TLS authentication configuration for a custom domain name.
	// Experimental.
	Mtls *MTLSConfig `json:"mtls"`
	// The Transport Layer Security (TLS) version + cipher suite for this domain name.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-domainname.html
	//
	// Experimental.
	SecurityPolicy SecurityPolicy `json:"securityPolicy"`
	// If specified, all requests to this domain will be mapped to the production deployment of this API.
	//
	// If you wish to map this domain to multiple APIs
	// with different base paths, don't specify this option and use
	// `addBasePathMapping`.
	// Experimental.
	Mapping IRestApi `json:"mapping"`
}

// Represents a reference to a REST API's Empty model, which is available as part of the model collection by default.
//
// This can be used for mapping
// JSON responses from an integration to what is returned to a client,
// where strong typing is not required. In the absence of any defined
// model, the Empty model will be used to return the response payload
// unmapped.
//
// Definition
// {
//    "$schema" : "http://json-schema.org/draft-04/schema#",
//    "title" : "Empty Schema",
//    "type" : "object"
// }
// See: Model.EMPTY_MODEL
//
// Deprecated: You should use
type EmptyModel interface {
	IModel
	ModelId() *string
}

// The jsii proxy struct for EmptyModel
type jsiiProxy_EmptyModel struct {
	jsiiProxy_IModel
}

func (j *jsiiProxy_EmptyModel) ModelId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"modelId",
		&returns,
	)
	return returns
}


// Deprecated: You should use
func NewEmptyModel() EmptyModel {
	_init_.Initialize()

	j := jsiiProxy_EmptyModel{}

	_jsii_.Create(
		"monocdk.aws_apigateway.EmptyModel",
		nil, // no parameters
		&j,
	)

	return &j
}

// Deprecated: You should use
func NewEmptyModel_Override(e EmptyModel) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.EmptyModel",
		nil, // no parameters
		e,
	)
}

// The endpoint configuration of a REST API, including VPCs and endpoint types.
//
// EndpointConfiguration is a property of the AWS::ApiGateway::RestApi resource.
// Experimental.
type EndpointConfiguration struct {
	// A list of endpoint types of an API or its custom domain name.
	// Experimental.
	Types *[]EndpointType `json:"types"`
	// A list of VPC Endpoints against which to create Route53 ALIASes.
	// Experimental.
	VpcEndpoints *[]awsec2.IVpcEndpoint `json:"vpcEndpoints"`
}

// Experimental.
type EndpointType string

const (
	EndpointType_EDGE EndpointType = "EDGE"
	EndpointType_REGIONAL EndpointType = "REGIONAL"
	EndpointType_PRIVATE EndpointType = "PRIVATE"
)

// Represents a reference to a REST API's Error model, which is available as part of the model collection by default.
//
// This can be used for mapping
// error JSON responses from an integration to a client, where a simple
// generic message field is sufficient to map and return an error payload.
//
// Definition
// {
//    "$schema" : "http://json-schema.org/draft-04/schema#",
//    "title" : "Error Schema",
//    "type" : "object",
//    "properties" : {
//      "message" : { "type" : "string" }
//    }
// }
// See: Model.ERROR_MODEL
//
// Deprecated: You should use
type ErrorModel interface {
	IModel
	ModelId() *string
}

// The jsii proxy struct for ErrorModel
type jsiiProxy_ErrorModel struct {
	jsiiProxy_IModel
}

func (j *jsiiProxy_ErrorModel) ModelId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"modelId",
		&returns,
	)
	return returns
}


// Deprecated: You should use
func NewErrorModel() ErrorModel {
	_init_.Initialize()

	j := jsiiProxy_ErrorModel{}

	_jsii_.Create(
		"monocdk.aws_apigateway.ErrorModel",
		nil, // no parameters
		&j,
	)

	return &j
}

// Deprecated: You should use
func NewErrorModel_Override(e ErrorModel) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.ErrorModel",
		nil, // no parameters
		e,
	)
}

// Configure the response received by clients, produced from the API Gateway backend.
// Experimental.
type GatewayResponse interface {
	awscdk.Resource
	IGatewayResponse
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
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

// The jsii proxy struct for GatewayResponse
type jsiiProxy_GatewayResponse struct {
	internal.Type__awscdkResource
	jsiiProxy_IGatewayResponse
}

func (j *jsiiProxy_GatewayResponse) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GatewayResponse) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GatewayResponse) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_GatewayResponse) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewGatewayResponse(scope constructs.Construct, id *string, props *GatewayResponseProps) GatewayResponse {
	_init_.Initialize()

	j := jsiiProxy_GatewayResponse{}

	_jsii_.Create(
		"monocdk.aws_apigateway.GatewayResponse",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewGatewayResponse_Override(g GatewayResponse, scope constructs.Construct, id *string, props *GatewayResponseProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.GatewayResponse",
		[]interface{}{scope, id, props},
		g,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func GatewayResponse_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.GatewayResponse",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func GatewayResponse_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.GatewayResponse",
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
func (g *jsiiProxy_GatewayResponse) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		g,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (g *jsiiProxy_GatewayResponse) GeneratePhysicalName() *string {
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
func (g *jsiiProxy_GatewayResponse) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (g *jsiiProxy_GatewayResponse) GetResourceNameAttribute(nameAttr *string) *string {
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
func (g *jsiiProxy_GatewayResponse) OnPrepare() {
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
func (g *jsiiProxy_GatewayResponse) OnSynthesize(session constructs.ISynthesisSession) {
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
func (g *jsiiProxy_GatewayResponse) OnValidate() *[]*string {
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
func (g *jsiiProxy_GatewayResponse) Prepare() {
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
func (g *jsiiProxy_GatewayResponse) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		g,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (g *jsiiProxy_GatewayResponse) ToString() *string {
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
func (g *jsiiProxy_GatewayResponse) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		g,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options to add gateway response.
// Experimental.
type GatewayResponseOptions struct {
	// Response type to associate with gateway response.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/supported-gateway-response-types.html
	//
	// Experimental.
	Type ResponseType `json:"type"`
	// Custom headers parameters for response.
	// Experimental.
	ResponseHeaders *map[string]*string `json:"responseHeaders"`
	// Http status code for response.
	// Experimental.
	StatusCode *string `json:"statusCode"`
	// Custom templates to get mapped as response.
	// Experimental.
	Templates *map[string]*string `json:"templates"`
}

// Properties for a new gateway response.
// Experimental.
type GatewayResponseProps struct {
	// Response type to associate with gateway response.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/supported-gateway-response-types.html
	//
	// Experimental.
	Type ResponseType `json:"type"`
	// Custom headers parameters for response.
	// Experimental.
	ResponseHeaders *map[string]*string `json:"responseHeaders"`
	// Http status code for response.
	// Experimental.
	StatusCode *string `json:"statusCode"`
	// Custom templates to get mapped as response.
	// Experimental.
	Templates *map[string]*string `json:"templates"`
	// Rest api resource to target.
	// Experimental.
	RestApi IRestApi `json:"restApi"`
}

// You can integrate an API method with an HTTP endpoint using the HTTP proxy integration or the HTTP custom integration,.
//
// With the proxy integration, the setup is simple. You only need to set the
// HTTP method and the HTTP endpoint URI, according to the backend requirements,
// if you are not concerned with content encoding or caching.
//
// With the custom integration, the setup is more involved. In addition to the
// proxy integration setup steps, you need to specify how the incoming request
// data is mapped to the integration request and how the resulting integration
// response data is mapped to the method response.
// Experimental.
type HttpIntegration interface {
	Integration
	Bind(_method Method) *IntegrationConfig
}

// The jsii proxy struct for HttpIntegration
type jsiiProxy_HttpIntegration struct {
	jsiiProxy_Integration
}

// Experimental.
func NewHttpIntegration(url *string, props *HttpIntegrationProps) HttpIntegration {
	_init_.Initialize()

	j := jsiiProxy_HttpIntegration{}

	_jsii_.Create(
		"monocdk.aws_apigateway.HttpIntegration",
		[]interface{}{url, props},
		&j,
	)

	return &j
}

// Experimental.
func NewHttpIntegration_Override(h HttpIntegration, url *string, props *HttpIntegrationProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.HttpIntegration",
		[]interface{}{url, props},
		h,
	)
}

// Can be overridden by subclasses to allow the integration to interact with the method being integrated, access the REST API object, method ARNs, etc.
// Experimental.
func (h *jsiiProxy_HttpIntegration) Bind(_method Method) *IntegrationConfig {
	var returns *IntegrationConfig

	_jsii_.Invoke(
		h,
		"bind",
		[]interface{}{_method},
		&returns,
	)

	return returns
}

// Experimental.
type HttpIntegrationProps struct {
	// HTTP method to use when invoking the backend URL.
	// Experimental.
	HttpMethod *string `json:"httpMethod"`
	// Integration options, such as request/resopnse mapping, content handling, etc.
	// Experimental.
	Options *IntegrationOptions `json:"options"`
	// Determines whether to use proxy integration or custom integration.
	// Experimental.
	Proxy *bool `json:"proxy"`
}

// Access log destination for a RestApi Stage.
// Experimental.
type IAccessLogDestination interface {
	// Binds this destination to the RestApi Stage.
	// Experimental.
	Bind(stage IStage) *AccessLogDestinationConfig
}

// The jsii proxy for IAccessLogDestination
type jsiiProxy_IAccessLogDestination struct {
	_ byte // padding
}

func (i *jsiiProxy_IAccessLogDestination) Bind(stage IStage) *AccessLogDestinationConfig {
	var returns *AccessLogDestinationConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{stage},
		&returns,
	)

	return returns
}

// API keys are alphanumeric string values that you distribute to app developer customers to grant access to your API.
// Experimental.
type IApiKey interface {
	awscdk.IResource
	// The API key ARN.
	// Experimental.
	KeyArn() *string
	// The API key ID.
	// Experimental.
	KeyId() *string
}

// The jsii proxy for IApiKey
type jsiiProxy_IApiKey struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IApiKey) KeyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApiKey) KeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyId",
		&returns,
	)
	return returns
}

// Represents an API Gateway authorizer.
// Experimental.
type IAuthorizer interface {
	// The authorization type of this authorizer.
	// Experimental.
	AuthorizationType() AuthorizationType
	// The authorizer ID.
	// Experimental.
	AuthorizerId() *string
}

// The jsii proxy for IAuthorizer
type jsiiProxy_IAuthorizer struct {
	_ byte // padding
}

func (j *jsiiProxy_IAuthorizer) AuthorizationType() AuthorizationType {
	var returns AuthorizationType
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IAuthorizer) AuthorizerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerId",
		&returns,
	)
	return returns
}

// Experimental.
type IDomainName interface {
	awscdk.IResource
	// The domain name (e.g. `example.com`).
	// Experimental.
	DomainName() *string
	// The Route53 alias target to use in order to connect a record set to this domain through an alias.
	// Experimental.
	DomainNameAliasDomainName() *string
	// The Route53 hosted zone ID to use in order to connect a record set to this domain through an alias.
	// Experimental.
	DomainNameAliasHostedZoneId() *string
}

// The jsii proxy for IDomainName
type jsiiProxy_IDomainName struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IDomainName) DomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IDomainName) DomainNameAliasDomainName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainNameAliasDomainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IDomainName) DomainNameAliasHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"domainNameAliasHostedZoneId",
		&returns,
	)
	return returns
}

// Represents gateway response resource.
// Experimental.
type IGatewayResponse interface {
	awscdk.IResource
}

// The jsii proxy for IGatewayResponse
type jsiiProxy_IGatewayResponse struct {
	internal.Type__awscdkIResource
}

// Experimental.
type IModel interface {
	// Returns the model name, such as 'myModel'.
	// Experimental.
	ModelId() *string
}

// The jsii proxy for IModel
type jsiiProxy_IModel struct {
	_ byte // padding
}

func (j *jsiiProxy_IModel) ModelId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"modelId",
		&returns,
	)
	return returns
}

// Experimental.
type IRequestValidator interface {
	awscdk.IResource
	// ID of the request validator, such as abc123.
	// Experimental.
	RequestValidatorId() *string
}

// The jsii proxy for IRequestValidator
type jsiiProxy_IRequestValidator struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IRequestValidator) RequestValidatorId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"requestValidatorId",
		&returns,
	)
	return returns
}

// Experimental.
type IResource interface {
	awscdk.IResource
	// Adds an OPTIONS method to this resource which responds to Cross-Origin Resource Sharing (CORS) preflight requests.
	//
	// Cross-Origin Resource Sharing (CORS) is a mechanism that uses additional
	// HTTP headers to tell browsers to give a web application running at one
	// origin, access to selected resources from a different origin. A web
	// application executes a cross-origin HTTP request when it requests a
	// resource that has a different origin (domain, protocol, or port) from its
	// own.
	//
	// Returns: a `Method` object
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
	//
	// Experimental.
	AddCorsPreflight(options *CorsOptions) Method
	// Defines a new method for this resource.
	//
	// Returns: The newly created `Method` object.
	// Experimental.
	AddMethod(httpMethod *string, target Integration, options *MethodOptions) Method
	// Adds a greedy proxy resource ("{proxy+}") and an ANY method to this route.
	// Experimental.
	AddProxy(options *ProxyResourceOptions) ProxyResource
	// Defines a new child resource where this resource is the parent.
	//
	// Returns: A Resource object
	// Experimental.
	AddResource(pathPart *string, options *ResourceOptions) Resource
	// Retrieves a child resource by path part.
	//
	// Returns: the child resource or undefined if not found
	// Experimental.
	GetResource(pathPart *string) IResource
	// Gets or create all resources leading up to the specified path.
	//
	// - Path may only start with "/" if this method is called on the root resource.
	// - All resources are created using default options.
	//
	// Returns: a new or existing resource.
	// Experimental.
	ResourceForPath(path *string) Resource
	// The rest API that this resource is part of.
	//
	// The reason we need the RestApi object itself and not just the ID is because the model
	// is being tracked by the top-level RestApi object for the purpose of calculating it's
	// hash to determine the ID of the deployment. This allows us to automatically update
	// the deployment when the model of the REST API changes.
	// Experimental.
	Api() IRestApi
	// Default options for CORS preflight OPTIONS method.
	// Experimental.
	DefaultCorsPreflightOptions() *CorsOptions
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration() Integration
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions() *MethodOptions
	// The parent of this resource or undefined for the root resource.
	// Experimental.
	ParentResource() IResource
	// The full path of this resource.
	// Experimental.
	Path() *string
	// The ID of the resource.
	// Experimental.
	ResourceId() *string
	// The rest API that this resource is part of.
	// Deprecated: - Throws an error if this Resource is not associated with an instance of `RestApi`. Use `api` instead.
	RestApi() RestApi
}

// The jsii proxy for IResource
type jsiiProxy_IResource struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IResource) AddCorsPreflight(options *CorsOptions) Method {
	var returns Method

	_jsii_.Invoke(
		i,
		"addCorsPreflight",
		[]interface{}{options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IResource) AddMethod(httpMethod *string, target Integration, options *MethodOptions) Method {
	var returns Method

	_jsii_.Invoke(
		i,
		"addMethod",
		[]interface{}{httpMethod, target, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IResource) AddProxy(options *ProxyResourceOptions) ProxyResource {
	var returns ProxyResource

	_jsii_.Invoke(
		i,
		"addProxy",
		[]interface{}{options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IResource) AddResource(pathPart *string, options *ResourceOptions) Resource {
	var returns Resource

	_jsii_.Invoke(
		i,
		"addResource",
		[]interface{}{pathPart, options},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IResource) GetResource(pathPart *string) IResource {
	var returns IResource

	_jsii_.Invoke(
		i,
		"getResource",
		[]interface{}{pathPart},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IResource) ResourceForPath(path *string) Resource {
	var returns Resource

	_jsii_.Invoke(
		i,
		"resourceForPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IResource) Api() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"api",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) DefaultCorsPreflightOptions() *CorsOptions {
	var returns *CorsOptions
	_jsii_.Get(
		j,
		"defaultCorsPreflightOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) DefaultIntegration() Integration {
	var returns Integration
	_jsii_.Get(
		j,
		"defaultIntegration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) DefaultMethodOptions() *MethodOptions {
	var returns *MethodOptions
	_jsii_.Get(
		j,
		"defaultMethodOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) ParentResource() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"parentResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) ResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IResource) RestApi() RestApi {
	var returns RestApi
	_jsii_.Get(
		j,
		"restApi",
		&returns,
	)
	return returns
}

// Experimental.
type IRestApi interface {
	awscdk.IResource
	// Gets the "execute-api" ARN.
	//
	// Returns: The "execute-api" ARN.
	// Experimental.
	ArnForExecuteApi(method *string, path *string, stage *string) *string
	// API Gateway stage that points to the latest deployment (if defined).
	// Experimental.
	DeploymentStage() Stage
	// API Gateway stage that points to the latest deployment (if defined).
	// Experimental.
	SetDeploymentStage(d Stage)
	// API Gateway deployment that represents the latest changes of the API.
	//
	// This resource will be automatically updated every time the REST API model changes.
	// `undefined` when no deployment is configured.
	// Experimental.
	LatestDeployment() Deployment
	// The ID of this API Gateway RestApi.
	// Experimental.
	RestApiId() *string
	// The resource ID of the root resource.
	// Experimental.
	RestApiRootResourceId() *string
	// Represents the root resource ("/") of this API. Use it to define the API model:.
	//
	// api.root.addMethod('ANY', redirectToHomePage); // "ANY /"
	//     api.root.addResource('friends').addMethod('GET', getFriendsHandler); // "GET /friends"
	// Experimental.
	Root() IResource
}

// The jsii proxy for IRestApi
type jsiiProxy_IRestApi struct {
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IRestApi) ArnForExecuteApi(method *string, path *string, stage *string) *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"arnForExecuteApi",
		[]interface{}{method, path, stage},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IRestApi) DeploymentStage() Stage {
	var returns Stage
	_jsii_.Get(
		j,
		"deploymentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRestApi) SetDeploymentStage(val Stage) {
	_jsii_.Set(
		j,
		"deploymentStage",
		val,
	)
}

func (j *jsiiProxy_IRestApi) LatestDeployment() Deployment {
	var returns Deployment
	_jsii_.Get(
		j,
		"latestDeployment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRestApi) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRestApi) RestApiRootResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiRootResourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IRestApi) Root() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"root",
		&returns,
	)
	return returns
}

// Represents an APIGateway Stage.
// Experimental.
type IStage interface {
	awscdk.IResource
	// RestApi to which this stage is associated.
	// Experimental.
	RestApi() IRestApi
	// Name of this stage.
	// Experimental.
	StageName() *string
}

// The jsii proxy for IStage
type jsiiProxy_IStage struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IStage) RestApi() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"restApi",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IStage) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}

// Represents an API Gateway VpcLink.
// Experimental.
type IVpcLink interface {
	awscdk.IResource
	// Physical ID of the VpcLink resource.
	// Experimental.
	VpcLinkId() *string
}

// The jsii proxy for IVpcLink
type jsiiProxy_IVpcLink struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_IVpcLink) VpcLinkId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"vpcLinkId",
		&returns,
	)
	return returns
}

// Represents an identity source.
//
// The source can be specified either as a literal value (e.g: `Auth`) which
// cannot be blank, or as an unresolved string token.
// Experimental.
type IdentitySource interface {
}

// The jsii proxy struct for IdentitySource
type jsiiProxy_IdentitySource struct {
	_ byte // padding
}

// Experimental.
func NewIdentitySource() IdentitySource {
	_init_.Initialize()

	j := jsiiProxy_IdentitySource{}

	_jsii_.Create(
		"monocdk.aws_apigateway.IdentitySource",
		nil, // no parameters
		&j,
	)

	return &j
}

// Experimental.
func NewIdentitySource_Override(i IdentitySource) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.IdentitySource",
		nil, // no parameters
		i,
	)
}

// Provides a properly formatted request context identity source.
//
// Returns: a request context identity source.
// Experimental.
func IdentitySource_Context(context *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.IdentitySource",
		"context",
		[]interface{}{context},
		&returns,
	)

	return returns
}

// Provides a properly formatted header identity source.
//
// Returns: a header identity source.
// Experimental.
func IdentitySource_Header(headerName *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.IdentitySource",
		"header",
		[]interface{}{headerName},
		&returns,
	)

	return returns
}

// Provides a properly formatted query string identity source.
//
// Returns: a query string identity source.
// Experimental.
func IdentitySource_QueryString(queryString *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.IdentitySource",
		"queryString",
		[]interface{}{queryString},
		&returns,
	)

	return returns
}

// Provides a properly formatted API Gateway stage variable identity source.
//
// Returns: an API Gateway stage variable identity source.
// Experimental.
func IdentitySource_StageVariable(stageVariable *string) *string {
	_init_.Initialize()

	var returns *string

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.IdentitySource",
		"stageVariable",
		[]interface{}{stageVariable},
		&returns,
	)

	return returns
}

// OpenAPI specification from an inline JSON object.
// Experimental.
type InlineApiDefinition interface {
	ApiDefinition
	Bind(_scope awscdk.Construct) *ApiDefinitionConfig
}

// The jsii proxy struct for InlineApiDefinition
type jsiiProxy_InlineApiDefinition struct {
	jsiiProxy_ApiDefinition
}

// Experimental.
func NewInlineApiDefinition(definition interface{}) InlineApiDefinition {
	_init_.Initialize()

	j := jsiiProxy_InlineApiDefinition{}

	_jsii_.Create(
		"monocdk.aws_apigateway.InlineApiDefinition",
		[]interface{}{definition},
		&j,
	)

	return &j
}

// Experimental.
func NewInlineApiDefinition_Override(i InlineApiDefinition, definition interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.InlineApiDefinition",
		[]interface{}{definition},
		i,
	)
}

// Loads the API specification from a local disk asset.
// Experimental.
func InlineApiDefinition_FromAsset(file *string, options *awss3assets.AssetOptions) AssetApiDefinition {
	_init_.Initialize()

	var returns AssetApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.InlineApiDefinition",
		"fromAsset",
		[]interface{}{file, options},
		&returns,
	)

	return returns
}

// Creates an API definition from a specification file in an S3 bucket.
// Experimental.
func InlineApiDefinition_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3ApiDefinition {
	_init_.Initialize()

	var returns S3ApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.InlineApiDefinition",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Create an API definition from an inline object.
//
// The inline object must follow the
// schema of OpenAPI 2.0 or OpenAPI 3.0
//
// TODO: EXAMPLE
//
// Experimental.
func InlineApiDefinition_FromInline(definition interface{}) InlineApiDefinition {
	_init_.Initialize()

	var returns InlineApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.InlineApiDefinition",
		"fromInline",
		[]interface{}{definition},
		&returns,
	)

	return returns
}

// Called when the specification is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (i *jsiiProxy_InlineApiDefinition) Bind(_scope awscdk.Construct) *ApiDefinitionConfig {
	var returns *ApiDefinitionConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// Base class for backend integrations for an API Gateway method.
//
// Use one of the concrete classes such as `MockIntegration`, `AwsIntegration`, `LambdaIntegration`
// or implement on your own by specifying the set of props.
// Experimental.
type Integration interface {
	Bind(_method Method) *IntegrationConfig
}

// The jsii proxy struct for Integration
type jsiiProxy_Integration struct {
	_ byte // padding
}

// Experimental.
func NewIntegration(props *IntegrationProps) Integration {
	_init_.Initialize()

	j := jsiiProxy_Integration{}

	_jsii_.Create(
		"monocdk.aws_apigateway.Integration",
		[]interface{}{props},
		&j,
	)

	return &j
}

// Experimental.
func NewIntegration_Override(i Integration, props *IntegrationProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.Integration",
		[]interface{}{props},
		i,
	)
}

// Can be overridden by subclasses to allow the integration to interact with the method being integrated, access the REST API object, method ARNs, etc.
// Experimental.
func (i *jsiiProxy_Integration) Bind(_method Method) *IntegrationConfig {
	var returns *IntegrationConfig

	_jsii_.Invoke(
		i,
		"bind",
		[]interface{}{_method},
		&returns,
	)

	return returns
}

// Result of binding an Integration to a Method.
// Experimental.
type IntegrationConfig struct {
	// Specifies an API method integration type.
	// Experimental.
	Type IntegrationType `json:"type"`
	// This value is included in computing the Deployment's fingerprint.
	//
	// When the fingerprint
	// changes, a new deployment is triggered.
	// This property should contain values associated with the Integration that upon changing
	// should trigger a fresh the Deployment needs to be refreshed.
	// Experimental.
	DeploymentToken *string `json:"deploymentToken"`
	// The integration's HTTP method type.
	// Experimental.
	IntegrationHttpMethod *string `json:"integrationHttpMethod"`
	// Integration options.
	// Experimental.
	Options *IntegrationOptions `json:"options"`
	// The Uniform Resource Identifier (URI) for the integration.
	// See: https://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
	//
	// Experimental.
	Uri *string `json:"uri"`
}

// Experimental.
type IntegrationOptions struct {
	// A list of request parameters whose values are to be cached.
	//
	// It determines
	// request parameters that will make it into the cache key.
	// Experimental.
	CacheKeyParameters *[]*string `json:"cacheKeyParameters"`
	// An API-specific tag group of related cached parameters.
	// Experimental.
	CacheNamespace *string `json:"cacheNamespace"`
	// The type of network connection to the integration endpoint.
	// Experimental.
	ConnectionType ConnectionType `json:"connectionType"`
	// Specifies how to handle request payload content type conversions.
	// Experimental.
	ContentHandling ContentHandling `json:"contentHandling"`
	// Requires that the caller's identity be passed through from the request.
	// Experimental.
	CredentialsPassthrough *bool `json:"credentialsPassthrough"`
	// An IAM role that API Gateway assumes.
	//
	// Mutually exclusive with `credentialsPassThrough`.
	// Experimental.
	CredentialsRole awsiam.IRole `json:"credentialsRole"`
	// The response that API Gateway provides after a method's backend completes processing a request.
	//
	// API Gateway intercepts the response from the
	// backend so that you can control how API Gateway surfaces backend
	// responses. For example, you can map the backend status codes to codes
	// that you define.
	// Experimental.
	IntegrationResponses *[]*IntegrationResponse `json:"integrationResponses"`
	// Specifies the pass-through behavior for incoming requests based on the Content-Type header in the request, and the available mapping templates specified as the requestTemplates property on the Integration resource.
	//
	// There are three valid values: WHEN_NO_MATCH, WHEN_NO_TEMPLATES, and
	// NEVER.
	// Experimental.
	PassthroughBehavior PassthroughBehavior `json:"passthroughBehavior"`
	// The request parameters that API Gateway sends with the backend request.
	//
	// Specify request parameters as key-value pairs (string-to-string
	// mappings), with a destination as the key and a source as the value.
	//
	// Specify the destination by using the following pattern
	// integration.request.location.name, where location is querystring, path,
	// or header, and name is a valid, unique parameter name.
	//
	// The source must be an existing method request parameter or a static
	// value. You must enclose static values in single quotation marks and
	// pre-encode these values based on their destination in the request.
	// Experimental.
	RequestParameters *map[string]*string `json:"requestParameters"`
	// A map of Apache Velocity templates that are applied on the request payload.
	//
	// The template that API Gateway uses is based on the value of the
	// Content-Type header that's sent by the client. The content type value is
	// the key, and the template is the value (specified as a string), such as
	// the following snippet:
	//
	// ```
	//    { "application/json": "{ \"statusCode\": 200 }" }
	// ```
	// See: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
	//
	// Experimental.
	RequestTemplates *map[string]*string `json:"requestTemplates"`
	// The maximum amount of time an integration will run before it returns without a response.
	//
	// Must be between 50 milliseconds and 29 seconds.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// The VpcLink used for the integration.
	//
	// Required if connectionType is VPC_LINK
	// Experimental.
	VpcLink IVpcLink `json:"vpcLink"`
}

// Experimental.
type IntegrationProps struct {
	// Specifies an API method integration type.
	// Experimental.
	Type IntegrationType `json:"type"`
	// The integration's HTTP method type.
	//
	// Required unless you use a MOCK integration.
	// Experimental.
	IntegrationHttpMethod *string `json:"integrationHttpMethod"`
	// Integration options.
	// Experimental.
	Options *IntegrationOptions `json:"options"`
	// The Uniform Resource Identifier (URI) for the integration.
	//
	// - If you specify HTTP for the `type` property, specify the API endpoint URL.
	// - If you specify MOCK for the `type` property, don't specify this property.
	// - If you specify AWS for the `type` property, specify an AWS service that
	//    follows this form: `arn:partition:apigateway:region:subdomain.service|service:path|action/service_api.`
	//    For example, a Lambda function URI follows this form:
	//    arn:partition:apigateway:region:lambda:path/path. The path is usually in the
	//    form /2015-03-31/functions/LambdaFunctionARN/invocations.
	// See: https://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
	//
	// Experimental.
	Uri interface{} `json:"uri"`
}

// Experimental.
type IntegrationResponse struct {
	// The status code that API Gateway uses to map the integration response to a MethodResponse status code.
	// Experimental.
	StatusCode *string `json:"statusCode"`
	// Specifies how to handle request payload content type conversions.
	// Experimental.
	ContentHandling ContentHandling `json:"contentHandling"`
	// The response parameters from the backend response that API Gateway sends to the method response.
	//
	// Use the destination as the key and the source as the value:
	//
	// - The destination must be an existing response parameter in the
	//    MethodResponse property.
	// - The source must be an existing method request parameter or a static
	//    value. You must enclose static values in single quotation marks and
	//    pre-encode these values based on the destination specified in the
	//    request.
	// See: http://docs.aws.amazon.com/apigateway/latest/developerguide/request-response-data-mappings.html
	//
	// Experimental.
	ResponseParameters *map[string]*string `json:"responseParameters"`
	// The templates that are used to transform the integration response body.
	//
	// Specify templates as key-value pairs, with a content type as the key and
	// a template as the value.
	// See: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
	//
	// Experimental.
	ResponseTemplates *map[string]*string `json:"responseTemplates"`
	// Specifies the regular expression (regex) pattern used to choose an integration response based on the response from the back end.
	//
	// For example, if the success response returns nothing and the error response returns some string, you
	// could use the ``.+`` regex to match error response. However, make sure that the error response does not contain any
	// newline (``\n``) character in such cases. If the back end is an AWS Lambda function, the AWS Lambda function error
	// header is matched. For all other HTTP and AWS back ends, the HTTP status code is matched.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-integration-settings-integration-response.html
	//
	// Experimental.
	SelectionPattern *string `json:"selectionPattern"`
}

// Experimental.
type IntegrationType string

const (
	IntegrationType_AWS IntegrationType = "AWS"
	IntegrationType_AWS_PROXY IntegrationType = "AWS_PROXY"
	IntegrationType_HTTP IntegrationType = "HTTP"
	IntegrationType_HTTP_PROXY IntegrationType = "HTTP_PROXY"
	IntegrationType_MOCK IntegrationType = "MOCK"
)

// Represents a JSON schema definition of the structure of a REST API model.
//
// Copied from npm module jsonschema.
// See: https://github.com/tdegrunt/jsonschema
//
// Experimental.
type JsonSchema struct {
	// Experimental.
	AdditionalItems *[]*JsonSchema `json:"additionalItems"`
	// Experimental.
	AdditionalProperties interface{} `json:"additionalProperties"`
	// Experimental.
	AllOf *[]*JsonSchema `json:"allOf"`
	// Experimental.
	AnyOf *[]*JsonSchema `json:"anyOf"`
	// Experimental.
	Contains interface{} `json:"contains"`
	// The default value if you use an enum.
	// Experimental.
	Default interface{} `json:"default"`
	// Experimental.
	Definitions *map[string]*JsonSchema `json:"definitions"`
	// Experimental.
	Dependencies *map[string]interface{} `json:"dependencies"`
	// Experimental.
	Description *string `json:"description"`
	// Experimental.
	Enum *[]interface{} `json:"enum"`
	// Experimental.
	ExclusiveMaximum *bool `json:"exclusiveMaximum"`
	// Experimental.
	ExclusiveMinimum *bool `json:"exclusiveMinimum"`
	// Experimental.
	Format *string `json:"format"`
	// Experimental.
	Id *string `json:"id"`
	// Experimental.
	Items interface{} `json:"items"`
	// Experimental.
	Maximum *float64 `json:"maximum"`
	// Experimental.
	MaxItems *float64 `json:"maxItems"`
	// Experimental.
	MaxLength *float64 `json:"maxLength"`
	// Experimental.
	MaxProperties *float64 `json:"maxProperties"`
	// Experimental.
	Minimum *float64 `json:"minimum"`
	// Experimental.
	MinItems *float64 `json:"minItems"`
	// Experimental.
	MinLength *float64 `json:"minLength"`
	// Experimental.
	MinProperties *float64 `json:"minProperties"`
	// Experimental.
	MultipleOf *float64 `json:"multipleOf"`
	// Experimental.
	Not **JsonSchema `json:"not"`
	// Experimental.
	OneOf *[]*JsonSchema `json:"oneOf"`
	// Experimental.
	Pattern *string `json:"pattern"`
	// Experimental.
	PatternProperties *map[string]*JsonSchema `json:"patternProperties"`
	// Experimental.
	Properties *map[string]*JsonSchema `json:"properties"`
	// Experimental.
	PropertyNames **JsonSchema `json:"propertyNames"`
	// Experimental.
	Ref *string `json:"ref"`
	// Experimental.
	Required *[]*string `json:"required"`
	// Experimental.
	Schema JsonSchemaVersion `json:"schema"`
	// Experimental.
	Title *string `json:"title"`
	// Experimental.
	Type interface{} `json:"type"`
	// Experimental.
	UniqueItems *bool `json:"uniqueItems"`
}

// Experimental.
type JsonSchemaType string

const (
	JsonSchemaType_NULL JsonSchemaType = "NULL"
	JsonSchemaType_BOOLEAN JsonSchemaType = "BOOLEAN"
	JsonSchemaType_OBJECT JsonSchemaType = "OBJECT"
	JsonSchemaType_ARRAY JsonSchemaType = "ARRAY"
	JsonSchemaType_NUMBER JsonSchemaType = "NUMBER"
	JsonSchemaType_INTEGER JsonSchemaType = "INTEGER"
	JsonSchemaType_STRING JsonSchemaType = "STRING"
)

// Experimental.
type JsonSchemaVersion string

const (
	JsonSchemaVersion_DRAFT4 JsonSchemaVersion = "DRAFT4"
	JsonSchemaVersion_DRAFT7 JsonSchemaVersion = "DRAFT7"
)

// Properties for controlling items output in JSON standard format.
// Experimental.
type JsonWithStandardFieldProps struct {
	// If this flag is enabled, the principal identifier of the caller will be output to the log.
	// Experimental.
	Caller *bool `json:"caller"`
	// If this flag is enabled, the http method will be output to the log.
	// Experimental.
	HttpMethod *bool `json:"httpMethod"`
	// If this flag is enabled, the source IP of request will be output to the log.
	// Experimental.
	Ip *bool `json:"ip"`
	// If this flag is enabled, the request protocol will be output to the log.
	// Experimental.
	Protocol *bool `json:"protocol"`
	// If this flag is enabled, the CLF-formatted request time((dd/MMM/yyyy:HH:mm:ss +-hhmm) will be output to the log.
	// Experimental.
	RequestTime *bool `json:"requestTime"`
	// If this flag is enabled, the path to your resource will be output to the log.
	// Experimental.
	ResourcePath *bool `json:"resourcePath"`
	// If this flag is enabled, the response payload length will be output to the log.
	// Experimental.
	ResponseLength *bool `json:"responseLength"`
	// If this flag is enabled, the method response status will be output to the log.
	// Experimental.
	Status *bool `json:"status"`
	// If this flag is enabled, the principal identifier of the user will be output to the log.
	// Experimental.
	User *bool `json:"user"`
}

// Base properties for all lambda authorizers.
// Experimental.
type LambdaAuthorizerProps struct {
	// The handler for the authorizer lambda function.
	//
	// The handler must follow a very specific protocol on the input it receives and the output it needs to produce.
	// API Gateway has documented the handler's input specification
	// {@link https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-input.html | here} and output specification
	// {@link https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html | here}.
	// Experimental.
	Handler awslambda.IFunction `json:"handler"`
	// An optional IAM role for APIGateway to assume before calling the Lambda-based authorizer.
	//
	// The IAM role must be
	// assumable by 'apigateway.amazonaws.com'.
	// Experimental.
	AssumeRole awsiam.IRole `json:"assumeRole"`
	// An optional human friendly name for the authorizer.
	//
	// Note that, this is not the primary identifier of the authorizer.
	// Experimental.
	AuthorizerName *string `json:"authorizerName"`
	// How long APIGateway should cache the results.
	//
	// Max 1 hour.
	// Disable caching by setting this to 0.
	// Experimental.
	ResultsCacheTtl awscdk.Duration `json:"resultsCacheTtl"`
}

// Integrates an AWS Lambda function to an API Gateway method.
//
// TODO: EXAMPLE
//
// Experimental.
type LambdaIntegration interface {
	AwsIntegration
	Bind(method Method) *IntegrationConfig
}

// The jsii proxy struct for LambdaIntegration
type jsiiProxy_LambdaIntegration struct {
	jsiiProxy_AwsIntegration
}

// Experimental.
func NewLambdaIntegration(handler awslambda.IFunction, options *LambdaIntegrationOptions) LambdaIntegration {
	_init_.Initialize()

	j := jsiiProxy_LambdaIntegration{}

	_jsii_.Create(
		"monocdk.aws_apigateway.LambdaIntegration",
		[]interface{}{handler, options},
		&j,
	)

	return &j
}

// Experimental.
func NewLambdaIntegration_Override(l LambdaIntegration, handler awslambda.IFunction, options *LambdaIntegrationOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.LambdaIntegration",
		[]interface{}{handler, options},
		l,
	)
}

// Can be overridden by subclasses to allow the integration to interact with the method being integrated, access the REST API object, method ARNs, etc.
// Experimental.
func (l *jsiiProxy_LambdaIntegration) Bind(method Method) *IntegrationConfig {
	var returns *IntegrationConfig

	_jsii_.Invoke(
		l,
		"bind",
		[]interface{}{method},
		&returns,
	)

	return returns
}

// Experimental.
type LambdaIntegrationOptions struct {
	// A list of request parameters whose values are to be cached.
	//
	// It determines
	// request parameters that will make it into the cache key.
	// Experimental.
	CacheKeyParameters *[]*string `json:"cacheKeyParameters"`
	// An API-specific tag group of related cached parameters.
	// Experimental.
	CacheNamespace *string `json:"cacheNamespace"`
	// The type of network connection to the integration endpoint.
	// Experimental.
	ConnectionType ConnectionType `json:"connectionType"`
	// Specifies how to handle request payload content type conversions.
	// Experimental.
	ContentHandling ContentHandling `json:"contentHandling"`
	// Requires that the caller's identity be passed through from the request.
	// Experimental.
	CredentialsPassthrough *bool `json:"credentialsPassthrough"`
	// An IAM role that API Gateway assumes.
	//
	// Mutually exclusive with `credentialsPassThrough`.
	// Experimental.
	CredentialsRole awsiam.IRole `json:"credentialsRole"`
	// The response that API Gateway provides after a method's backend completes processing a request.
	//
	// API Gateway intercepts the response from the
	// backend so that you can control how API Gateway surfaces backend
	// responses. For example, you can map the backend status codes to codes
	// that you define.
	// Experimental.
	IntegrationResponses *[]*IntegrationResponse `json:"integrationResponses"`
	// Specifies the pass-through behavior for incoming requests based on the Content-Type header in the request, and the available mapping templates specified as the requestTemplates property on the Integration resource.
	//
	// There are three valid values: WHEN_NO_MATCH, WHEN_NO_TEMPLATES, and
	// NEVER.
	// Experimental.
	PassthroughBehavior PassthroughBehavior `json:"passthroughBehavior"`
	// The request parameters that API Gateway sends with the backend request.
	//
	// Specify request parameters as key-value pairs (string-to-string
	// mappings), with a destination as the key and a source as the value.
	//
	// Specify the destination by using the following pattern
	// integration.request.location.name, where location is querystring, path,
	// or header, and name is a valid, unique parameter name.
	//
	// The source must be an existing method request parameter or a static
	// value. You must enclose static values in single quotation marks and
	// pre-encode these values based on their destination in the request.
	// Experimental.
	RequestParameters *map[string]*string `json:"requestParameters"`
	// A map of Apache Velocity templates that are applied on the request payload.
	//
	// The template that API Gateway uses is based on the value of the
	// Content-Type header that's sent by the client. The content type value is
	// the key, and the template is the value (specified as a string), such as
	// the following snippet:
	//
	// ```
	//    { "application/json": "{ \"statusCode\": 200 }" }
	// ```
	// See: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
	//
	// Experimental.
	RequestTemplates *map[string]*string `json:"requestTemplates"`
	// The maximum amount of time an integration will run before it returns without a response.
	//
	// Must be between 50 milliseconds and 29 seconds.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// The VpcLink used for the integration.
	//
	// Required if connectionType is VPC_LINK
	// Experimental.
	VpcLink IVpcLink `json:"vpcLink"`
	// Allow invoking method from AWS Console UI (for testing purposes).
	//
	// This will add another permission to the AWS Lambda resource policy which
	// will allow the `test-invoke-stage` stage to invoke this handler. If this
	// is set to `false`, the function will only be usable from the deployment
	// endpoint.
	// Experimental.
	AllowTestInvoke *bool `json:"allowTestInvoke"`
	// Use proxy integration or normal (request/response mapping) integration.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-output-format
	//
	// Experimental.
	Proxy *bool `json:"proxy"`
}

// Defines an API Gateway REST API with AWS Lambda proxy integration.
//
// Use the `proxy` property to define a greedy proxy ("{proxy+}") and "ANY"
// method from the specified path. If not defined, you will need to explicity
// add resources and methods to the API.
// Experimental.
type LambdaRestApi interface {
	RestApi
	DeploymentStage() Stage
	SetDeploymentStage(val Stage)
	DomainName() DomainName
	Env() *awscdk.ResourceEnvironment
	LatestDeployment() Deployment
	Methods() *[]Method
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RestApiId() *string
	RestApiName() *string
	RestApiRootResourceId() *string
	Root() IResource
	Stack() awscdk.Stack
	Url() *string
	AddApiKey(id *string, options *ApiKeyOptions) IApiKey
	AddDomainName(id *string, options *DomainNameOptions) DomainName
	AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse
	AddModel(id *string, props *ModelOptions) Model
	AddRequestValidator(id *string, props *RequestValidatorOptions) RequestValidator
	AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ArnForExecuteApi(method *string, path *string, stage *string) *string
	ConfigureCloudWatchRole(apiResource CfnRestApi)
	ConfigureDeployment(props *RestApiBaseProps)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	UrlForPath(path *string) *string
	Validate() *[]*string
}

// The jsii proxy struct for LambdaRestApi
type jsiiProxy_LambdaRestApi struct {
	jsiiProxy_RestApi
}

func (j *jsiiProxy_LambdaRestApi) DeploymentStage() Stage {
	var returns Stage
	_jsii_.Get(
		j,
		"deploymentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) DomainName() DomainName {
	var returns DomainName
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) LatestDeployment() Deployment {
	var returns Deployment
	_jsii_.Get(
		j,
		"latestDeployment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) Methods() *[]Method {
	var returns *[]Method
	_jsii_.Get(
		j,
		"methods",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) RestApiName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) RestApiRootResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiRootResourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) Root() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"root",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_LambdaRestApi) Url() *string {
	var returns *string
	_jsii_.Get(
		j,
		"url",
		&returns,
	)
	return returns
}


// Experimental.
func NewLambdaRestApi(scope constructs.Construct, id *string, props *LambdaRestApiProps) LambdaRestApi {
	_init_.Initialize()

	j := jsiiProxy_LambdaRestApi{}

	_jsii_.Create(
		"monocdk.aws_apigateway.LambdaRestApi",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewLambdaRestApi_Override(l LambdaRestApi, scope constructs.Construct, id *string, props *LambdaRestApiProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.LambdaRestApi",
		[]interface{}{scope, id, props},
		l,
	)
}

func (j *jsiiProxy_LambdaRestApi) SetDeploymentStage(val Stage) {
	_jsii_.Set(
		j,
		"deploymentStage",
		val,
	)
}

// Import an existing RestApi that can be configured with additional Methods and Resources.
// Experimental.
func LambdaRestApi_FromRestApiAttributes(scope constructs.Construct, id *string, attrs *RestApiAttributes) IRestApi {
	_init_.Initialize()

	var returns IRestApi

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.LambdaRestApi",
		"fromRestApiAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Import an existing RestApi.
// Experimental.
func LambdaRestApi_FromRestApiId(scope constructs.Construct, id *string, restApiId *string) IRestApi {
	_init_.Initialize()

	var returns IRestApi

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.LambdaRestApi",
		"fromRestApiId",
		[]interface{}{scope, id, restApiId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func LambdaRestApi_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.LambdaRestApi",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func LambdaRestApi_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.LambdaRestApi",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add an ApiKey.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) AddApiKey(id *string, options *ApiKeyOptions) IApiKey {
	var returns IApiKey

	_jsii_.Invoke(
		l,
		"addApiKey",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an API Gateway domain name and maps it to this API.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) AddDomainName(id *string, options *DomainNameOptions) DomainName {
	var returns DomainName

	_jsii_.Invoke(
		l,
		"addDomainName",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a new gateway response.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse {
	var returns GatewayResponse

	_jsii_.Invoke(
		l,
		"addGatewayResponse",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a new model.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) AddModel(id *string, props *ModelOptions) Model {
	var returns Model

	_jsii_.Invoke(
		l,
		"addModel",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds a new request validator.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) AddRequestValidator(id *string, props *RequestValidatorOptions) RequestValidator {
	var returns RequestValidator

	_jsii_.Invoke(
		l,
		"addRequestValidator",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds a usage plan.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan {
	var returns UsagePlan

	_jsii_.Invoke(
		l,
		"addUsagePlan",
		[]interface{}{id, props},
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
func (l *jsiiProxy_LambdaRestApi) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		l,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Gets the "execute-api" ARN.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) ArnForExecuteApi(method *string, path *string, stage *string) *string {
	var returns *string

	_jsii_.Invoke(
		l,
		"arnForExecuteApi",
		[]interface{}{method, path, stage},
		&returns,
	)

	return returns
}

// Deprecated: This method will be made internal. No replacement
func (l *jsiiProxy_LambdaRestApi) ConfigureCloudWatchRole(apiResource CfnRestApi) {
	_jsii_.InvokeVoid(
		l,
		"configureCloudWatchRole",
		[]interface{}{apiResource},
	)
}

// Deprecated: This method will be made internal. No replacement
func (l *jsiiProxy_LambdaRestApi) ConfigureDeployment(props *RestApiBaseProps) {
	_jsii_.InvokeVoid(
		l,
		"configureDeployment",
		[]interface{}{props},
	)
}

// Experimental.
func (l *jsiiProxy_LambdaRestApi) GeneratePhysicalName() *string {
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
func (l *jsiiProxy_LambdaRestApi) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (l *jsiiProxy_LambdaRestApi) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		l,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns the given named metric for this API.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the API cache in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (l *jsiiProxy_LambdaRestApi) MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metricCacheHitCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the backend in a given period, when API caching is enabled.
//
// Default: sum over 5 minutes
// Experimental.
func (l *jsiiProxy_LambdaRestApi) MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metricCacheMissCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of client-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (l *jsiiProxy_LambdaRestApi) MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metricClientError",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the total number API requests in a given period.
//
// Default: sample count over 5 minutes
// Experimental.
func (l *jsiiProxy_LambdaRestApi) MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metricCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the time between when API Gateway relays a request to the backend and when it receives a response from the backend.
//
// Default: average over 5 minutes.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metricIntegrationLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time between when API Gateway receives a request from a client and when it returns a response to the client.
//
// The latency includes the integration latency and other API Gateway overhead.
//
// Default: average over 5 minutes.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metricLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of server-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (l *jsiiProxy_LambdaRestApi) MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		l,
		"metricServerError",
		[]interface{}{props},
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
func (l *jsiiProxy_LambdaRestApi) OnPrepare() {
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
func (l *jsiiProxy_LambdaRestApi) OnSynthesize(session constructs.ISynthesisSession) {
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
func (l *jsiiProxy_LambdaRestApi) OnValidate() *[]*string {
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
func (l *jsiiProxy_LambdaRestApi) Prepare() {
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
func (l *jsiiProxy_LambdaRestApi) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		l,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		l,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns the URL for an HTTP path.
//
// Fails if `deploymentStage` is not set either by `deploy` or explicitly.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) UrlForPath(path *string) *string {
	var returns *string

	_jsii_.Invoke(
		l,
		"urlForPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Performs validation of the REST API.
// Experimental.
func (l *jsiiProxy_LambdaRestApi) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		l,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type LambdaRestApiProps struct {
	// Automatically configure an AWS CloudWatch role for API Gateway.
	// Experimental.
	CloudWatchRole *bool `json:"cloudWatchRole"`
	// Indicates if a Deployment should be automatically created for this API, and recreated when the API model (resources, methods) changes.
	//
	// Since API Gateway deployments are immutable, When this option is enabled
	// (by default), an AWS::ApiGateway::Deployment resource will automatically
	// created with a logical ID that hashes the API model (methods, resources
	// and options). This means that when the model changes, the logical ID of
	// this CloudFormation resource will change, and a new deployment will be
	// created.
	//
	// If this is set, `latestDeployment` will refer to the `Deployment` object
	// and `deploymentStage` will refer to a `Stage` that points to this
	// deployment. To customize the stage options, use the `deployOptions`
	// property.
	//
	// A CloudFormation Output will also be defined with the root URL endpoint
	// of this REST API.
	// Experimental.
	Deploy *bool `json:"deploy"`
	// Options for the API Gateway stage that will always point to the latest deployment when `deploy` is enabled.
	//
	// If `deploy` is disabled,
	// this value cannot be set.
	// Experimental.
	DeployOptions *StageOptions `json:"deployOptions"`
	// Specifies whether clients can invoke the API using the default execute-api endpoint.
	//
	// To require that clients use a custom domain name to invoke the
	// API, disable the default endpoint.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-restapi.html
	//
	// Experimental.
	DisableExecuteApiEndpoint *bool `json:"disableExecuteApiEndpoint"`
	// Configure a custom domain name and map it to this API.
	// Experimental.
	DomainName *DomainNameOptions `json:"domainName"`
	// Export name for the CfnOutput containing the API endpoint.
	// Experimental.
	EndpointExportName *string `json:"endpointExportName"`
	// A list of the endpoint types of the API.
	//
	// Use this property when creating
	// an API.
	// Experimental.
	EndpointTypes *[]EndpointType `json:"endpointTypes"`
	// Indicates whether to roll back the resource if a warning occurs while API Gateway is creating the RestApi resource.
	// Experimental.
	FailOnWarnings *bool `json:"failOnWarnings"`
	// Custom header parameters for the request.
	// See: https://docs.aws.amazon.com/cli/latest/reference/apigateway/import-rest-api.html
	//
	// Experimental.
	Parameters *map[string]*string `json:"parameters"`
	// A policy document that contains the permissions for this RestApi.
	// Experimental.
	Policy awsiam.PolicyDocument `json:"policy"`
	// A name for the API Gateway RestApi resource.
	// Experimental.
	RestApiName *string `json:"restApiName"`
	// Retains old deployment resources when the API changes.
	//
	// This allows
	// manually reverting stages to point to old deployments via the AWS
	// Console.
	// Experimental.
	RetainDeployments *bool `json:"retainDeployments"`
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// The source of the API key for metering requests according to a usage plan.
	// Experimental.
	ApiKeySourceType ApiKeySourceType `json:"apiKeySourceType"`
	// The list of binary media mime-types that are supported by the RestApi resource, such as "image/png" or "application/octet-stream".
	// Experimental.
	BinaryMediaTypes *[]*string `json:"binaryMediaTypes"`
	// The ID of the API Gateway RestApi resource that you want to clone.
	// Experimental.
	CloneFrom IRestApi `json:"cloneFrom"`
	// A description of the purpose of this API Gateway RestApi resource.
	// Experimental.
	Description *string `json:"description"`
	// The EndpointConfiguration property type specifies the endpoint types of a REST API.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-endpointconfiguration.html
	//
	// Experimental.
	EndpointConfiguration *EndpointConfiguration `json:"endpointConfiguration"`
	// A nullable integer that is used to enable compression (with non-negative between 0 and 10485760 (10M) bytes, inclusive) or disable compression (when undefined) on an API.
	//
	// When compression is enabled, compression or
	// decompression is not applied on the payload if the payload size is
	// smaller than this value. Setting it to zero allows compression for any
	// payload size.
	// Experimental.
	MinimumCompressionSize *float64 `json:"minimumCompressionSize"`
	// The default Lambda function that handles all requests from this API.
	//
	// This handler will be used as a the default integration for all methods in
	// this API, unless specified otherwise in `addMethod`.
	// Experimental.
	Handler awslambda.IFunction `json:"handler"`
	// Deprecated: the `LambdaRestApiProps` now extends `RestApiProps`, so all
	// options are just available here. Note that the options specified in
	// `options` will be overridden by any props specified at the root level.
	Options *RestApiProps `json:"options"`
	// If true, route all requests to the Lambda Function.
	//
	// If set to false, you will need to explicitly define the API model using
	// `addResource` and `addMethod` (or `addProxy`).
	// Experimental.
	Proxy *bool `json:"proxy"`
}

// Use CloudWatch Logs as a custom access log destination for API Gateway.
// Experimental.
type LogGroupLogDestination interface {
	IAccessLogDestination
	Bind(_stage IStage) *AccessLogDestinationConfig
}

// The jsii proxy struct for LogGroupLogDestination
type jsiiProxy_LogGroupLogDestination struct {
	jsiiProxy_IAccessLogDestination
}

// Experimental.
func NewLogGroupLogDestination(logGroup awslogs.ILogGroup) LogGroupLogDestination {
	_init_.Initialize()

	j := jsiiProxy_LogGroupLogDestination{}

	_jsii_.Create(
		"monocdk.aws_apigateway.LogGroupLogDestination",
		[]interface{}{logGroup},
		&j,
	)

	return &j
}

// Experimental.
func NewLogGroupLogDestination_Override(l LogGroupLogDestination, logGroup awslogs.ILogGroup) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.LogGroupLogDestination",
		[]interface{}{logGroup},
		l,
	)
}

// Binds this destination to the CloudWatch Logs.
// Experimental.
func (l *jsiiProxy_LogGroupLogDestination) Bind(_stage IStage) *AccessLogDestinationConfig {
	var returns *AccessLogDestinationConfig

	_jsii_.Invoke(
		l,
		"bind",
		[]interface{}{_stage},
		&returns,
	)

	return returns
}

// The mTLS authentication configuration for a custom domain name.
// Experimental.
type MTLSConfig struct {
	// The bucket that the trust store is hosted in.
	// Experimental.
	Bucket awss3.IBucket `json:"bucket"`
	// The key in S3 to look at for the trust store.
	// Experimental.
	Key *string `json:"key"`
	// The version of the S3 object that contains your truststore.
	//
	// To specify a version, you must have versioning enabled for the S3 bucket.
	// Experimental.
	Version *string `json:"version"`
}

// Experimental.
type Method interface {
	awscdk.Resource
	Api() IRestApi
	Env() *awscdk.ResourceEnvironment
	HttpMethod() *string
	MethodArn() *string
	MethodId() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Resource() IResource
	RestApi() RestApi
	Stack() awscdk.Stack
	TestMethodArn() *string
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

// The jsii proxy struct for Method
type jsiiProxy_Method struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_Method) Api() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"api",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) HttpMethod() *string {
	var returns *string
	_jsii_.Get(
		j,
		"httpMethod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) MethodArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"methodArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) MethodId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"methodId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) Resource() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"resource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) RestApi() RestApi {
	var returns RestApi
	_jsii_.Get(
		j,
		"restApi",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Method) TestMethodArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"testMethodArn",
		&returns,
	)
	return returns
}


// Experimental.
func NewMethod(scope constructs.Construct, id *string, props *MethodProps) Method {
	_init_.Initialize()

	j := jsiiProxy_Method{}

	_jsii_.Create(
		"monocdk.aws_apigateway.Method",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewMethod_Override(m Method, scope constructs.Construct, id *string, props *MethodProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.Method",
		[]interface{}{scope, id, props},
		m,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Method_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Method",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Method_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Method",
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
func (m *jsiiProxy_Method) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		m,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (m *jsiiProxy_Method) GeneratePhysicalName() *string {
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
func (m *jsiiProxy_Method) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (m *jsiiProxy_Method) GetResourceNameAttribute(nameAttr *string) *string {
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
func (m *jsiiProxy_Method) OnPrepare() {
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
func (m *jsiiProxy_Method) OnSynthesize(session constructs.ISynthesisSession) {
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
func (m *jsiiProxy_Method) OnValidate() *[]*string {
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
func (m *jsiiProxy_Method) Prepare() {
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
func (m *jsiiProxy_Method) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		m,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (m *jsiiProxy_Method) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (m *jsiiProxy_Method) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		m,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type MethodDeploymentOptions struct {
	// Indicates whether the cached responses are encrypted.
	// Experimental.
	CacheDataEncrypted *bool `json:"cacheDataEncrypted"`
	// Specifies the time to live (TTL), in seconds, for cached responses.
	//
	// The
	// higher the TTL, the longer the response will be cached.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-caching.html
	//
	// Experimental.
	CacheTtl awscdk.Duration `json:"cacheTtl"`
	// Specifies whether responses should be cached and returned for requests.
	//
	// A
	// cache cluster must be enabled on the stage for responses to be cached.
	// Experimental.
	CachingEnabled *bool `json:"cachingEnabled"`
	// Specifies whether data trace logging is enabled for this method, which effects the log entries pushed to Amazon CloudWatch Logs.
	// Experimental.
	DataTraceEnabled *bool `json:"dataTraceEnabled"`
	// Specifies the logging level for this method, which effects the log entries pushed to Amazon CloudWatch Logs.
	// Experimental.
	LoggingLevel MethodLoggingLevel `json:"loggingLevel"`
	// Specifies whether Amazon CloudWatch metrics are enabled for this method.
	// Experimental.
	MetricsEnabled *bool `json:"metricsEnabled"`
	// Specifies the throttling burst limit.
	//
	// The total rate of all requests in your AWS account is limited to 5,000 requests.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html
	//
	// Experimental.
	ThrottlingBurstLimit *float64 `json:"throttlingBurstLimit"`
	// Specifies the throttling rate limit.
	//
	// The total rate of all requests in your AWS account is limited to 10,000 requests per second (rps).
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html
	//
	// Experimental.
	ThrottlingRateLimit *float64 `json:"throttlingRateLimit"`
}

// Experimental.
type MethodLoggingLevel string

const (
	MethodLoggingLevel_OFF MethodLoggingLevel = "OFF"
	MethodLoggingLevel_ERROR MethodLoggingLevel = "ERROR"
	MethodLoggingLevel_INFO MethodLoggingLevel = "INFO"
)

// Experimental.
type MethodOptions struct {
	// Indicates whether the method requires clients to submit a valid API key.
	// Experimental.
	ApiKeyRequired *bool `json:"apiKeyRequired"`
	// A list of authorization scopes configured on the method.
	//
	// The scopes are used with
	// a COGNITO_USER_POOLS authorizer to authorize the method invocation.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-method.html#cfn-apigateway-method-authorizationscopes
	//
	// Experimental.
	AuthorizationScopes *[]*string `json:"authorizationScopes"`
	// Method authorization. If the value is set of `Custom`, an `authorizer` must also be specified.
	//
	// If you're using one of the authorizers that are available via the {@link Authorizer} class, such as {@link Authorizer#token()},
	// it is recommended that this option not be specified. The authorizer will take care of setting the correct authorization type.
	// However, specifying an authorization type using this property that conflicts with what is expected by the {@link Authorizer}
	// will result in an error.
	// Experimental.
	AuthorizationType AuthorizationType `json:"authorizationType"`
	// If `authorizationType` is `Custom`, this specifies the ID of the method authorizer resource.
	//
	// If specified, the value of `authorizationType` must be set to `Custom`
	// Experimental.
	Authorizer IAuthorizer `json:"authorizer"`
	// The responses that can be sent to the client who calls the method.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-method-settings-method-response.html
	//
	// Experimental.
	MethodResponses *[]*MethodResponse `json:"methodResponses"`
	// A friendly operation name for the method.
	//
	// For example, you can assign the
	// OperationName of ListPets for the GET /pets method.
	// Experimental.
	OperationName *string `json:"operationName"`
	// The models which describe data structure of request payload.
	//
	// When
	// combined with `requestValidator` or `requestValidatorOptions`, the service
	// will validate the API request payload before it reaches the API's Integration (including proxies).
	// Specify `requestModels` as key-value pairs, with a content type
	// (e.g. `'application/json'`) as the key and an API Gateway Model as the value.
	//
	// TODO: EXAMPLE
	//
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-method-settings-method-request.html#setup-method-request-model
	//
	// Experimental.
	RequestModels *map[string]IModel `json:"requestModels"`
	// The request parameters that API Gateway accepts.
	//
	// Specify request parameters
	// as key-value pairs (string-to-Boolean mapping), with a source as the key and
	// a Boolean as the value. The Boolean specifies whether a parameter is required.
	// A source must match the format method.request.location.name, where the location
	// is querystring, path, or header, and name is a valid, unique parameter name.
	// Experimental.
	RequestParameters *map[string]*bool `json:"requestParameters"`
	// The ID of the associated request validator.
	//
	// Only one of `requestValidator` or `requestValidatorOptions` must be specified.
	// Works together with `requestModels` or `requestParameters` to validate
	// the request before it reaches integration like Lambda Proxy Integration.
	// Experimental.
	RequestValidator IRequestValidator `json:"requestValidator"`
	// Request validator options to create new validator Only one of `requestValidator` or `requestValidatorOptions` must be specified.
	//
	// Works together with `requestModels` or `requestParameters` to validate
	// the request before it reaches integration like Lambda Proxy Integration.
	// Experimental.
	RequestValidatorOptions *RequestValidatorOptions `json:"requestValidatorOptions"`
}

// Experimental.
type MethodProps struct {
	// The HTTP method ("GET", "POST", "PUT", ...) that clients use to call this method.
	// Experimental.
	HttpMethod *string `json:"httpMethod"`
	// The resource this method is associated with.
	//
	// For root resource methods,
	// specify the `RestApi` object.
	// Experimental.
	Resource IResource `json:"resource"`
	// The backend system that the method calls when it receives a request.
	// Experimental.
	Integration Integration `json:"integration"`
	// Method options.
	// Experimental.
	Options *MethodOptions `json:"options"`
}

// Experimental.
type MethodResponse struct {
	// The method response's status code, which you map to an IntegrationResponse.
	//
	// Required.
	// Experimental.
	StatusCode *string `json:"statusCode"`
	// The resources used for the response's content type.
	//
	// Specify response models as
	// key-value pairs (string-to-string maps), with a content type as the key and a Model
	// resource name as the value.
	// Experimental.
	ResponseModels *map[string]IModel `json:"responseModels"`
	// Response parameters that API Gateway sends to the client that called a method.
	//
	// Specify response parameters as key-value pairs (string-to-Boolean maps), with
	// a destination as the key and a Boolean as the value. Specify the destination
	// using the following pattern: method.response.header.name, where the name is a
	// valid, unique header name. The Boolean specifies whether a parameter is required.
	// Experimental.
	ResponseParameters *map[string]*bool `json:"responseParameters"`
}

// This type of integration lets API Gateway return a response without sending the request further to the backend.
//
// This is useful for API testing because it
// can be used to test the integration set up without incurring charges for
// using the backend and to enable collaborative development of an API. In
// collaborative development, a team can isolate their development effort by
// setting up simulations of API components owned by other teams by using the
// MOCK integrations. It is also used to return CORS-related headers to ensure
// that the API method permits CORS access. In fact, the API Gateway console
// integrates the OPTIONS method to support CORS with a mock integration.
// Gateway responses are other examples of mock integrations.
// Experimental.
type MockIntegration interface {
	Integration
	Bind(_method Method) *IntegrationConfig
}

// The jsii proxy struct for MockIntegration
type jsiiProxy_MockIntegration struct {
	jsiiProxy_Integration
}

// Experimental.
func NewMockIntegration(options *IntegrationOptions) MockIntegration {
	_init_.Initialize()

	j := jsiiProxy_MockIntegration{}

	_jsii_.Create(
		"monocdk.aws_apigateway.MockIntegration",
		[]interface{}{options},
		&j,
	)

	return &j
}

// Experimental.
func NewMockIntegration_Override(m MockIntegration, options *IntegrationOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.MockIntegration",
		[]interface{}{options},
		m,
	)
}

// Can be overridden by subclasses to allow the integration to interact with the method being integrated, access the REST API object, method ARNs, etc.
// Experimental.
func (m *jsiiProxy_MockIntegration) Bind(_method Method) *IntegrationConfig {
	var returns *IntegrationConfig

	_jsii_.Invoke(
		m,
		"bind",
		[]interface{}{_method},
		&returns,
	)

	return returns
}

// Experimental.
type Model interface {
	awscdk.Resource
	IModel
	Env() *awscdk.ResourceEnvironment
	ModelId() *string
	Node() awscdk.ConstructNode
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

// The jsii proxy struct for Model
type jsiiProxy_Model struct {
	internal.Type__awscdkResource
	jsiiProxy_IModel
}

func (j *jsiiProxy_Model) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Model) ModelId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"modelId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Model) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Model) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Model) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewModel(scope constructs.Construct, id *string, props *ModelProps) Model {
	_init_.Initialize()

	j := jsiiProxy_Model{}

	_jsii_.Create(
		"monocdk.aws_apigateway.Model",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewModel_Override(m Model, scope constructs.Construct, id *string, props *ModelProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.Model",
		[]interface{}{scope, id, props},
		m,
	)
}

// Experimental.
func Model_FromModelName(scope constructs.Construct, id *string, modelName *string) IModel {
	_init_.Initialize()

	var returns IModel

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Model",
		"fromModelName",
		[]interface{}{scope, id, modelName},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Model_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Model",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Model_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Model",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

func Model_EMPTY_MODEL() IModel {
	_init_.Initialize()
	var returns IModel
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.Model",
		"EMPTY_MODEL",
		&returns,
	)
	return returns
}

func Model_ERROR_MODEL() IModel {
	_init_.Initialize()
	var returns IModel
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.Model",
		"ERROR_MODEL",
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
func (m *jsiiProxy_Model) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		m,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (m *jsiiProxy_Model) GeneratePhysicalName() *string {
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
func (m *jsiiProxy_Model) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (m *jsiiProxy_Model) GetResourceNameAttribute(nameAttr *string) *string {
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
func (m *jsiiProxy_Model) OnPrepare() {
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
func (m *jsiiProxy_Model) OnSynthesize(session constructs.ISynthesisSession) {
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
func (m *jsiiProxy_Model) OnValidate() *[]*string {
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
func (m *jsiiProxy_Model) Prepare() {
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
func (m *jsiiProxy_Model) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		m,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (m *jsiiProxy_Model) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (m *jsiiProxy_Model) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		m,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type ModelOptions struct {
	// The schema to use to transform data to one or more output formats.
	//
	// Specify null ({}) if you don't want to specify a schema.
	// Experimental.
	Schema *JsonSchema `json:"schema"`
	// The content type for the model.
	//
	// You can also force a
	// content type in the request or response model mapping.
	// Experimental.
	ContentType *string `json:"contentType"`
	// A description that identifies this model.
	// Experimental.
	Description *string `json:"description"`
	// A name for the model.
	//
	// Important
	//   If you specify a name, you cannot perform updates that
	//   require replacement of this resource. You can perform
	//   updates that require no or some interruption. If you
	//   must replace the resource, specify a new name.
	// Experimental.
	ModelName *string `json:"modelName"`
}

// Experimental.
type ModelProps struct {
	// The schema to use to transform data to one or more output formats.
	//
	// Specify null ({}) if you don't want to specify a schema.
	// Experimental.
	Schema *JsonSchema `json:"schema"`
	// The content type for the model.
	//
	// You can also force a
	// content type in the request or response model mapping.
	// Experimental.
	ContentType *string `json:"contentType"`
	// A description that identifies this model.
	// Experimental.
	Description *string `json:"description"`
	// A name for the model.
	//
	// Important
	//   If you specify a name, you cannot perform updates that
	//   require replacement of this resource. You can perform
	//   updates that require no or some interruption. If you
	//   must replace the resource, specify a new name.
	// Experimental.
	ModelName *string `json:"modelName"`
	// The rest API that this model is part of.
	//
	// The reason we need the RestApi object itself and not just the ID is because the model
	// is being tracked by the top-level RestApi object for the purpose of calculating it's
	// hash to determine the ID of the deployment. This allows us to automatically update
	// the deployment when the model of the REST API changes.
	// Experimental.
	RestApi IRestApi `json:"restApi"`
}

// Experimental.
type PassthroughBehavior string

const (
	PassthroughBehavior_WHEN_NO_MATCH PassthroughBehavior = "WHEN_NO_MATCH"
	PassthroughBehavior_NEVER PassthroughBehavior = "NEVER"
	PassthroughBehavior_WHEN_NO_TEMPLATES PassthroughBehavior = "WHEN_NO_TEMPLATES"
)

// Time period for which quota settings apply.
// Experimental.
type Period string

const (
	Period_DAY Period = "DAY"
	Period_WEEK Period = "WEEK"
	Period_MONTH Period = "MONTH"
)

// Defines a {proxy+} greedy resource and an ANY method on a route.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html
//
// Experimental.
type ProxyResource interface {
	Resource
	AnyMethod() Method
	Api() IRestApi
	DefaultCorsPreflightOptions() *CorsOptions
	DefaultIntegration() Integration
	DefaultMethodOptions() *MethodOptions
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	ParentResource() IResource
	Path() *string
	PhysicalName() *string
	ResourceId() *string
	RestApi() RestApi
	Stack() awscdk.Stack
	Url() *string
	AddCorsPreflight(options *CorsOptions) Method
	AddMethod(httpMethod *string, integration Integration, options *MethodOptions) Method
	AddProxy(options *ProxyResourceOptions) ProxyResource
	AddResource(pathPart *string, options *ResourceOptions) Resource
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResource(pathPart *string) IResource
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	ResourceForPath(path *string) Resource
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ProxyResource
type jsiiProxy_ProxyResource struct {
	jsiiProxy_Resource
}

func (j *jsiiProxy_ProxyResource) AnyMethod() Method {
	var returns Method
	_jsii_.Get(
		j,
		"anyMethod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) Api() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"api",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) DefaultCorsPreflightOptions() *CorsOptions {
	var returns *CorsOptions
	_jsii_.Get(
		j,
		"defaultCorsPreflightOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) DefaultIntegration() Integration {
	var returns Integration
	_jsii_.Get(
		j,
		"defaultIntegration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) DefaultMethodOptions() *MethodOptions {
	var returns *MethodOptions
	_jsii_.Get(
		j,
		"defaultMethodOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) ParentResource() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"parentResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) ResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) RestApi() RestApi {
	var returns RestApi
	_jsii_.Get(
		j,
		"restApi",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ProxyResource) Url() *string {
	var returns *string
	_jsii_.Get(
		j,
		"url",
		&returns,
	)
	return returns
}


// Experimental.
func NewProxyResource(scope constructs.Construct, id *string, props *ProxyResourceProps) ProxyResource {
	_init_.Initialize()

	j := jsiiProxy_ProxyResource{}

	_jsii_.Create(
		"monocdk.aws_apigateway.ProxyResource",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewProxyResource_Override(p ProxyResource, scope constructs.Construct, id *string, props *ProxyResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.ProxyResource",
		[]interface{}{scope, id, props},
		p,
	)
}

// Import an existing resource.
// Experimental.
func ProxyResource_FromResourceAttributes(scope constructs.Construct, id *string, attrs *ResourceAttributes) IResource {
	_init_.Initialize()

	var returns IResource

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ProxyResource",
		"fromResourceAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ProxyResource_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ProxyResource",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ProxyResource_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ProxyResource",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds an OPTIONS method to this resource which responds to Cross-Origin Resource Sharing (CORS) preflight requests.
//
// Cross-Origin Resource Sharing (CORS) is a mechanism that uses additional
// HTTP headers to tell browsers to give a web application running at one
// origin, access to selected resources from a different origin. A web
// application executes a cross-origin HTTP request when it requests a
// resource that has a different origin (domain, protocol, or port) from its
// own.
// Experimental.
func (p *jsiiProxy_ProxyResource) AddCorsPreflight(options *CorsOptions) Method {
	var returns Method

	_jsii_.Invoke(
		p,
		"addCorsPreflight",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Defines a new method for this resource.
// Experimental.
func (p *jsiiProxy_ProxyResource) AddMethod(httpMethod *string, integration Integration, options *MethodOptions) Method {
	var returns Method

	_jsii_.Invoke(
		p,
		"addMethod",
		[]interface{}{httpMethod, integration, options},
		&returns,
	)

	return returns
}

// Adds a greedy proxy resource ("{proxy+}") and an ANY method to this route.
// Experimental.
func (p *jsiiProxy_ProxyResource) AddProxy(options *ProxyResourceOptions) ProxyResource {
	var returns ProxyResource

	_jsii_.Invoke(
		p,
		"addProxy",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Defines a new child resource where this resource is the parent.
// Experimental.
func (p *jsiiProxy_ProxyResource) AddResource(pathPart *string, options *ResourceOptions) Resource {
	var returns Resource

	_jsii_.Invoke(
		p,
		"addResource",
		[]interface{}{pathPart, options},
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
func (p *jsiiProxy_ProxyResource) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		p,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (p *jsiiProxy_ProxyResource) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		p,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Retrieves a child resource by path part.
// Experimental.
func (p *jsiiProxy_ProxyResource) GetResource(pathPart *string) IResource {
	var returns IResource

	_jsii_.Invoke(
		p,
		"getResource",
		[]interface{}{pathPart},
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
func (p *jsiiProxy_ProxyResource) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (p *jsiiProxy_ProxyResource) GetResourceNameAttribute(nameAttr *string) *string {
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
func (p *jsiiProxy_ProxyResource) OnPrepare() {
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
func (p *jsiiProxy_ProxyResource) OnSynthesize(session constructs.ISynthesisSession) {
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
func (p *jsiiProxy_ProxyResource) OnValidate() *[]*string {
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
func (p *jsiiProxy_ProxyResource) Prepare() {
	_jsii_.InvokeVoid(
		p,
		"prepare",
		nil, // no parameters
	)
}

// Gets or create all resources leading up to the specified path.
//
// - Path may only start with "/" if this method is called on the root resource.
// - All resources are created using default options.
// Experimental.
func (p *jsiiProxy_ProxyResource) ResourceForPath(path *string) Resource {
	var returns Resource

	_jsii_.Invoke(
		p,
		"resourceForPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (p *jsiiProxy_ProxyResource) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		p,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (p *jsiiProxy_ProxyResource) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (p *jsiiProxy_ProxyResource) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		p,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type ProxyResourceOptions struct {
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// Adds an "ANY" method to this resource.
	//
	// If set to `false`, you will have to explicitly
	// add methods to this resource after it's created.
	// Experimental.
	AnyMethod *bool `json:"anyMethod"`
}

// Experimental.
type ProxyResourceProps struct {
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// Adds an "ANY" method to this resource.
	//
	// If set to `false`, you will have to explicitly
	// add methods to this resource after it's created.
	// Experimental.
	AnyMethod *bool `json:"anyMethod"`
	// The parent resource of this resource.
	//
	// You can either pass another
	// `Resource` object or a `RestApi` object here.
	// Experimental.
	Parent IResource `json:"parent"`
}

// Specifies the maximum number of requests that clients can make to API Gateway APIs.
// Experimental.
type QuotaSettings struct {
	// The maximum number of requests that users can make within the specified time period.
	// Experimental.
	Limit *float64 `json:"limit"`
	// For the initial time period, the number of requests to subtract from the specified limit.
	// Experimental.
	Offset *float64 `json:"offset"`
	// The time period for which the maximum limit of requests applies.
	// Experimental.
	Period Period `json:"period"`
}

// An API Gateway ApiKey, for which a rate limiting configuration can be specified.
// Experimental.
type RateLimitedApiKey interface {
	awscdk.Resource
	IApiKey
	Env() *awscdk.ResourceEnvironment
	KeyArn() *string
	KeyId() *string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	GrantRead(grantee awsiam.IGrantable) awsiam.Grant
	GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant
	GrantWrite(grantee awsiam.IGrantable) awsiam.Grant
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for RateLimitedApiKey
type jsiiProxy_RateLimitedApiKey struct {
	internal.Type__awscdkResource
	jsiiProxy_IApiKey
}

func (j *jsiiProxy_RateLimitedApiKey) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RateLimitedApiKey) KeyArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RateLimitedApiKey) KeyId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RateLimitedApiKey) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RateLimitedApiKey) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RateLimitedApiKey) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRateLimitedApiKey(scope constructs.Construct, id *string, props *RateLimitedApiKeyProps) RateLimitedApiKey {
	_init_.Initialize()

	j := jsiiProxy_RateLimitedApiKey{}

	_jsii_.Create(
		"monocdk.aws_apigateway.RateLimitedApiKey",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRateLimitedApiKey_Override(r RateLimitedApiKey, scope constructs.Construct, id *string, props *RateLimitedApiKeyProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.RateLimitedApiKey",
		[]interface{}{scope, id, props},
		r,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func RateLimitedApiKey_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RateLimitedApiKey",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func RateLimitedApiKey_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RateLimitedApiKey",
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
func (r *jsiiProxy_RateLimitedApiKey) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_RateLimitedApiKey) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_RateLimitedApiKey) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_RateLimitedApiKey) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Permits the IAM principal all read operations through this key.
// Experimental.
func (r *jsiiProxy_RateLimitedApiKey) GrantRead(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grantRead",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Permits the IAM principal all read and write operations through this key.
// Experimental.
func (r *jsiiProxy_RateLimitedApiKey) GrantReadWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grantReadWrite",
		[]interface{}{grantee},
		&returns,
	)

	return returns
}

// Permits the IAM principal all write operations through this key.
// Experimental.
func (r *jsiiProxy_RateLimitedApiKey) GrantWrite(grantee awsiam.IGrantable) awsiam.Grant {
	var returns awsiam.Grant

	_jsii_.Invoke(
		r,
		"grantWrite",
		[]interface{}{grantee},
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
func (r *jsiiProxy_RateLimitedApiKey) OnPrepare() {
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
func (r *jsiiProxy_RateLimitedApiKey) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_RateLimitedApiKey) OnValidate() *[]*string {
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
func (r *jsiiProxy_RateLimitedApiKey) Prepare() {
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
func (r *jsiiProxy_RateLimitedApiKey) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_RateLimitedApiKey) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (r *jsiiProxy_RateLimitedApiKey) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// RateLimitedApiKey properties.
// Experimental.
type RateLimitedApiKeyProps struct {
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// A name for the API key.
	//
	// If you don't specify a name, AWS CloudFormation generates a unique physical ID and uses that ID for the API key name.
	// Experimental.
	ApiKeyName *string `json:"apiKeyName"`
	// The value of the API key.
	//
	// Must be at least 20 characters long.
	// Experimental.
	Value *string `json:"value"`
	// An AWS Marketplace customer identifier to use when integrating with the AWS SaaS Marketplace.
	// Experimental.
	CustomerId *string `json:"customerId"`
	// A description of the purpose of the API key.
	// Experimental.
	Description *string `json:"description"`
	// Indicates whether the API key can be used by clients.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// Specifies whether the key identifier is distinct from the created API key value.
	// Experimental.
	GenerateDistinctId *bool `json:"generateDistinctId"`
	// A list of resources this api key is associated with.
	// Experimental.
	Resources *[]IRestApi `json:"resources"`
	// API Stages to be associated with the RateLimitedApiKey.
	// Experimental.
	ApiStages *[]*UsagePlanPerApiStage `json:"apiStages"`
	// Number of requests clients can make in a given time period.
	// Experimental.
	Quota *QuotaSettings `json:"quota"`
	// Overall throttle settings for the API.
	// Experimental.
	Throttle *ThrottleSettings `json:"throttle"`
}

// Request-based lambda authorizer that recognizes the caller's identity via request parameters, such as headers, paths, query strings, stage variables, or context variables.
//
// Based on the request, authorization is performed by a lambda function.
// Experimental.
type RequestAuthorizer interface {
	Authorizer
	IAuthorizer
	AuthorizationType() AuthorizationType
	AuthorizerArn() *string
	AuthorizerId() *string
	Env() *awscdk.ResourceEnvironment
	Handler() awslambda.IFunction
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Role() awsiam.IRole
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LazyRestApiId() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	SetupPermissions()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for RequestAuthorizer
type jsiiProxy_RequestAuthorizer struct {
	jsiiProxy_Authorizer
	jsiiProxy_IAuthorizer
}

func (j *jsiiProxy_RequestAuthorizer) AuthorizationType() AuthorizationType {
	var returns AuthorizationType
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) AuthorizerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) AuthorizerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) Handler() awslambda.IFunction {
	var returns awslambda.IFunction
	_jsii_.Get(
		j,
		"handler",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestAuthorizer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRequestAuthorizer(scope constructs.Construct, id *string, props *RequestAuthorizerProps) RequestAuthorizer {
	_init_.Initialize()

	j := jsiiProxy_RequestAuthorizer{}

	_jsii_.Create(
		"monocdk.aws_apigateway.RequestAuthorizer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRequestAuthorizer_Override(r RequestAuthorizer, scope constructs.Construct, id *string, props *RequestAuthorizerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.RequestAuthorizer",
		[]interface{}{scope, id, props},
		r,
	)
}

func (j *jsiiProxy_RequestAuthorizer) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

// Return whether the given object is an Authorizer.
// Experimental.
func RequestAuthorizer_IsAuthorizer(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RequestAuthorizer",
		"isAuthorizer",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func RequestAuthorizer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RequestAuthorizer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func RequestAuthorizer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RequestAuthorizer",
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
func (r *jsiiProxy_RequestAuthorizer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_RequestAuthorizer) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_RequestAuthorizer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_RequestAuthorizer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns a token that resolves to the Rest Api Id at the time of synthesis.
//
// Throws an error, during token resolution, if no RestApi is attached to this authorizer.
// Experimental.
func (r *jsiiProxy_RequestAuthorizer) LazyRestApiId() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"lazyRestApiId",
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
func (r *jsiiProxy_RequestAuthorizer) OnPrepare() {
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
func (r *jsiiProxy_RequestAuthorizer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_RequestAuthorizer) OnValidate() *[]*string {
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
func (r *jsiiProxy_RequestAuthorizer) Prepare() {
	_jsii_.InvokeVoid(
		r,
		"prepare",
		nil, // no parameters
	)
}

// Sets up the permissions necessary for the API Gateway service to invoke the Lambda function.
// Experimental.
func (r *jsiiProxy_RequestAuthorizer) SetupPermissions() {
	_jsii_.InvokeVoid(
		r,
		"setupPermissions",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (r *jsiiProxy_RequestAuthorizer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_RequestAuthorizer) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (r *jsiiProxy_RequestAuthorizer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for RequestAuthorizer.
// Experimental.
type RequestAuthorizerProps struct {
	// The handler for the authorizer lambda function.
	//
	// The handler must follow a very specific protocol on the input it receives and the output it needs to produce.
	// API Gateway has documented the handler's input specification
	// {@link https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-input.html | here} and output specification
	// {@link https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html | here}.
	// Experimental.
	Handler awslambda.IFunction `json:"handler"`
	// An optional IAM role for APIGateway to assume before calling the Lambda-based authorizer.
	//
	// The IAM role must be
	// assumable by 'apigateway.amazonaws.com'.
	// Experimental.
	AssumeRole awsiam.IRole `json:"assumeRole"`
	// An optional human friendly name for the authorizer.
	//
	// Note that, this is not the primary identifier of the authorizer.
	// Experimental.
	AuthorizerName *string `json:"authorizerName"`
	// How long APIGateway should cache the results.
	//
	// Max 1 hour.
	// Disable caching by setting this to 0.
	// Experimental.
	ResultsCacheTtl awscdk.Duration `json:"resultsCacheTtl"`
	// An array of request header mapping expressions for identities.
	//
	// Supported parameter types are
	// Header, Query String, Stage Variable, and Context. For instance, extracting an authorization
	// token from a header would use the identity source `IdentitySource.header('Authorizer')`.
	//
	// Note: API Gateway uses the specified identity sources as the request authorizer caching key. When caching is
	// enabled, API Gateway calls the authorizer's Lambda function only after successfully verifying that all the
	// specified identity sources are present at runtime. If a specified identify source is missing, null, or empty,
	// API Gateway returns a 401 Unauthorized response without calling the authorizer Lambda function.
	// See: https://docs.aws.amazon.com/apigateway/api-reference/link-relation/authorizer-create/#identitySource
	//
	// Experimental.
	IdentitySources *[]*string `json:"identitySources"`
}

// Experimental.
type RequestValidator interface {
	awscdk.Resource
	IRequestValidator
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RequestValidatorId() *string
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

// The jsii proxy struct for RequestValidator
type jsiiProxy_RequestValidator struct {
	internal.Type__awscdkResource
	jsiiProxy_IRequestValidator
}

func (j *jsiiProxy_RequestValidator) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestValidator) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestValidator) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestValidator) RequestValidatorId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"requestValidatorId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RequestValidator) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRequestValidator(scope constructs.Construct, id *string, props *RequestValidatorProps) RequestValidator {
	_init_.Initialize()

	j := jsiiProxy_RequestValidator{}

	_jsii_.Create(
		"monocdk.aws_apigateway.RequestValidator",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRequestValidator_Override(r RequestValidator, scope constructs.Construct, id *string, props *RequestValidatorProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.RequestValidator",
		[]interface{}{scope, id, props},
		r,
	)
}

// Experimental.
func RequestValidator_FromRequestValidatorId(scope constructs.Construct, id *string, requestValidatorId *string) IRequestValidator {
	_init_.Initialize()

	var returns IRequestValidator

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RequestValidator",
		"fromRequestValidatorId",
		[]interface{}{scope, id, requestValidatorId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func RequestValidator_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RequestValidator",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func RequestValidator_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RequestValidator",
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
func (r *jsiiProxy_RequestValidator) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_RequestValidator) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_RequestValidator) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_RequestValidator) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_RequestValidator) OnPrepare() {
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
func (r *jsiiProxy_RequestValidator) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_RequestValidator) OnValidate() *[]*string {
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
func (r *jsiiProxy_RequestValidator) Prepare() {
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
func (r *jsiiProxy_RequestValidator) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_RequestValidator) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (r *jsiiProxy_RequestValidator) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type RequestValidatorOptions struct {
	// The name of this request validator.
	// Experimental.
	RequestValidatorName *string `json:"requestValidatorName"`
	// Indicates whether to validate the request body according to the configured schema for the targeted API and method.
	// Experimental.
	ValidateRequestBody *bool `json:"validateRequestBody"`
	// Indicates whether to validate request parameters.
	// Experimental.
	ValidateRequestParameters *bool `json:"validateRequestParameters"`
}

// Experimental.
type RequestValidatorProps struct {
	// The name of this request validator.
	// Experimental.
	RequestValidatorName *string `json:"requestValidatorName"`
	// Indicates whether to validate the request body according to the configured schema for the targeted API and method.
	// Experimental.
	ValidateRequestBody *bool `json:"validateRequestBody"`
	// Indicates whether to validate request parameters.
	// Experimental.
	ValidateRequestParameters *bool `json:"validateRequestParameters"`
	// The rest API that this model is part of.
	//
	// The reason we need the RestApi object itself and not just the ID is because the model
	// is being tracked by the top-level RestApi object for the purpose of calculating it's
	// hash to determine the ID of the deployment. This allows us to automatically update
	// the deployment when the model of the REST API changes.
	// Experimental.
	RestApi IRestApi `json:"restApi"`
}

// Experimental.
type Resource interface {
	ResourceBase
	Api() IRestApi
	DefaultCorsPreflightOptions() *CorsOptions
	DefaultIntegration() Integration
	DefaultMethodOptions() *MethodOptions
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	ParentResource() IResource
	Path() *string
	PhysicalName() *string
	ResourceId() *string
	RestApi() RestApi
	Stack() awscdk.Stack
	Url() *string
	AddCorsPreflight(options *CorsOptions) Method
	AddMethod(httpMethod *string, integration Integration, options *MethodOptions) Method
	AddProxy(options *ProxyResourceOptions) ProxyResource
	AddResource(pathPart *string, options *ResourceOptions) Resource
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResource(pathPart *string) IResource
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	ResourceForPath(path *string) Resource
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for Resource
type jsiiProxy_Resource struct {
	jsiiProxy_ResourceBase
}

func (j *jsiiProxy_Resource) Api() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"api",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) DefaultCorsPreflightOptions() *CorsOptions {
	var returns *CorsOptions
	_jsii_.Get(
		j,
		"defaultCorsPreflightOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) DefaultIntegration() Integration {
	var returns Integration
	_jsii_.Get(
		j,
		"defaultIntegration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) DefaultMethodOptions() *MethodOptions {
	var returns *MethodOptions
	_jsii_.Get(
		j,
		"defaultMethodOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) ParentResource() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"parentResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) ResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) RestApi() RestApi {
	var returns RestApi
	_jsii_.Get(
		j,
		"restApi",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Resource) Url() *string {
	var returns *string
	_jsii_.Get(
		j,
		"url",
		&returns,
	)
	return returns
}


// Experimental.
func NewResource(scope constructs.Construct, id *string, props *ResourceProps) Resource {
	_init_.Initialize()

	j := jsiiProxy_Resource{}

	_jsii_.Create(
		"monocdk.aws_apigateway.Resource",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewResource_Override(r Resource, scope constructs.Construct, id *string, props *ResourceProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.Resource",
		[]interface{}{scope, id, props},
		r,
	)
}

// Import an existing resource.
// Experimental.
func Resource_FromResourceAttributes(scope constructs.Construct, id *string, attrs *ResourceAttributes) IResource {
	_init_.Initialize()

	var returns IResource

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Resource",
		"fromResourceAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func Resource_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Resource",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Resource_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Resource",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds an OPTIONS method to this resource which responds to Cross-Origin Resource Sharing (CORS) preflight requests.
//
// Cross-Origin Resource Sharing (CORS) is a mechanism that uses additional
// HTTP headers to tell browsers to give a web application running at one
// origin, access to selected resources from a different origin. A web
// application executes a cross-origin HTTP request when it requests a
// resource that has a different origin (domain, protocol, or port) from its
// own.
// Experimental.
func (r *jsiiProxy_Resource) AddCorsPreflight(options *CorsOptions) Method {
	var returns Method

	_jsii_.Invoke(
		r,
		"addCorsPreflight",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Defines a new method for this resource.
// Experimental.
func (r *jsiiProxy_Resource) AddMethod(httpMethod *string, integration Integration, options *MethodOptions) Method {
	var returns Method

	_jsii_.Invoke(
		r,
		"addMethod",
		[]interface{}{httpMethod, integration, options},
		&returns,
	)

	return returns
}

// Adds a greedy proxy resource ("{proxy+}") and an ANY method to this route.
// Experimental.
func (r *jsiiProxy_Resource) AddProxy(options *ProxyResourceOptions) ProxyResource {
	var returns ProxyResource

	_jsii_.Invoke(
		r,
		"addProxy",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Defines a new child resource where this resource is the parent.
// Experimental.
func (r *jsiiProxy_Resource) AddResource(pathPart *string, options *ResourceOptions) Resource {
	var returns Resource

	_jsii_.Invoke(
		r,
		"addResource",
		[]interface{}{pathPart, options},
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
func (r *jsiiProxy_Resource) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_Resource) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Retrieves a child resource by path part.
// Experimental.
func (r *jsiiProxy_Resource) GetResource(pathPart *string) IResource {
	var returns IResource

	_jsii_.Invoke(
		r,
		"getResource",
		[]interface{}{pathPart},
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
func (r *jsiiProxy_Resource) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_Resource) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_Resource) OnPrepare() {
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
func (r *jsiiProxy_Resource) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_Resource) OnValidate() *[]*string {
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
func (r *jsiiProxy_Resource) Prepare() {
	_jsii_.InvokeVoid(
		r,
		"prepare",
		nil, // no parameters
	)
}

// Gets or create all resources leading up to the specified path.
//
// - Path may only start with "/" if this method is called on the root resource.
// - All resources are created using default options.
// Experimental.
func (r *jsiiProxy_Resource) ResourceForPath(path *string) Resource {
	var returns Resource

	_jsii_.Invoke(
		r,
		"resourceForPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (r *jsiiProxy_Resource) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_Resource) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (r *jsiiProxy_Resource) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Attributes that can be specified when importing a Resource.
// Experimental.
type ResourceAttributes struct {
	// The full path of this resource.
	// Experimental.
	Path *string `json:"path"`
	// The ID of the resource.
	// Experimental.
	ResourceId *string `json:"resourceId"`
	// The rest API that this resource is part of.
	// Experimental.
	RestApi IRestApi `json:"restApi"`
}

// Experimental.
type ResourceBase interface {
	awscdk.Resource
	IResource
	Api() IRestApi
	DefaultCorsPreflightOptions() *CorsOptions
	DefaultIntegration() Integration
	DefaultMethodOptions() *MethodOptions
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	ParentResource() IResource
	Path() *string
	PhysicalName() *string
	ResourceId() *string
	RestApi() RestApi
	Stack() awscdk.Stack
	Url() *string
	AddCorsPreflight(options *CorsOptions) Method
	AddMethod(httpMethod *string, integration Integration, options *MethodOptions) Method
	AddProxy(options *ProxyResourceOptions) ProxyResource
	AddResource(pathPart *string, options *ResourceOptions) Resource
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResource(pathPart *string) IResource
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	ResourceForPath(path *string) Resource
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ResourceBase
type jsiiProxy_ResourceBase struct {
	internal.Type__awscdkResource
	jsiiProxy_IResource
}

func (j *jsiiProxy_ResourceBase) Api() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"api",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) DefaultCorsPreflightOptions() *CorsOptions {
	var returns *CorsOptions
	_jsii_.Get(
		j,
		"defaultCorsPreflightOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) DefaultIntegration() Integration {
	var returns Integration
	_jsii_.Get(
		j,
		"defaultIntegration",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) DefaultMethodOptions() *MethodOptions {
	var returns *MethodOptions
	_jsii_.Get(
		j,
		"defaultMethodOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) ParentResource() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"parentResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) Path() *string {
	var returns *string
	_jsii_.Get(
		j,
		"path",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) ResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) RestApi() RestApi {
	var returns RestApi
	_jsii_.Get(
		j,
		"restApi",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ResourceBase) Url() *string {
	var returns *string
	_jsii_.Get(
		j,
		"url",
		&returns,
	)
	return returns
}


// Experimental.
func NewResourceBase_Override(r ResourceBase, scope constructs.Construct, id *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.ResourceBase",
		[]interface{}{scope, id},
		r,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func ResourceBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ResourceBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ResourceBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ResourceBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds an OPTIONS method to this resource which responds to Cross-Origin Resource Sharing (CORS) preflight requests.
//
// Cross-Origin Resource Sharing (CORS) is a mechanism that uses additional
// HTTP headers to tell browsers to give a web application running at one
// origin, access to selected resources from a different origin. A web
// application executes a cross-origin HTTP request when it requests a
// resource that has a different origin (domain, protocol, or port) from its
// own.
// Experimental.
func (r *jsiiProxy_ResourceBase) AddCorsPreflight(options *CorsOptions) Method {
	var returns Method

	_jsii_.Invoke(
		r,
		"addCorsPreflight",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Defines a new method for this resource.
// Experimental.
func (r *jsiiProxy_ResourceBase) AddMethod(httpMethod *string, integration Integration, options *MethodOptions) Method {
	var returns Method

	_jsii_.Invoke(
		r,
		"addMethod",
		[]interface{}{httpMethod, integration, options},
		&returns,
	)

	return returns
}

// Adds a greedy proxy resource ("{proxy+}") and an ANY method to this route.
// Experimental.
func (r *jsiiProxy_ResourceBase) AddProxy(options *ProxyResourceOptions) ProxyResource {
	var returns ProxyResource

	_jsii_.Invoke(
		r,
		"addProxy",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Defines a new child resource where this resource is the parent.
// Experimental.
func (r *jsiiProxy_ResourceBase) AddResource(pathPart *string, options *ResourceOptions) Resource {
	var returns Resource

	_jsii_.Invoke(
		r,
		"addResource",
		[]interface{}{pathPart, options},
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
func (r *jsiiProxy_ResourceBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (r *jsiiProxy_ResourceBase) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"generatePhysicalName",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Retrieves a child resource by path part.
// Experimental.
func (r *jsiiProxy_ResourceBase) GetResource(pathPart *string) IResource {
	var returns IResource

	_jsii_.Invoke(
		r,
		"getResource",
		[]interface{}{pathPart},
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
func (r *jsiiProxy_ResourceBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_ResourceBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
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
func (r *jsiiProxy_ResourceBase) OnPrepare() {
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
func (r *jsiiProxy_ResourceBase) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_ResourceBase) OnValidate() *[]*string {
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
func (r *jsiiProxy_ResourceBase) Prepare() {
	_jsii_.InvokeVoid(
		r,
		"prepare",
		nil, // no parameters
	)
}

// Gets or create all resources leading up to the specified path.
//
// - Path may only start with "/" if this method is called on the root resource.
// - All resources are created using default options.
// Experimental.
func (r *jsiiProxy_ResourceBase) ResourceForPath(path *string) Resource {
	var returns Resource

	_jsii_.Invoke(
		r,
		"resourceForPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (r *jsiiProxy_ResourceBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_ResourceBase) ToString() *string {
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
//
// Returns: An array of validation error messages, or an empty array if the construct is valid.
// Experimental.
func (r *jsiiProxy_ResourceBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type ResourceOptions struct {
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
}

// Experimental.
type ResourceProps struct {
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// The parent resource of this resource.
	//
	// You can either pass another
	// `Resource` object or a `RestApi` object here.
	// Experimental.
	Parent IResource `json:"parent"`
	// A path name for the resource.
	// Experimental.
	PathPart *string `json:"pathPart"`
}

// Supported types of gateway responses.
// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/supported-gateway-response-types.html
//
// Experimental.
type ResponseType interface {
	ResponseType() *string
}

// The jsii proxy struct for ResponseType
type jsiiProxy_ResponseType struct {
	_ byte // padding
}

func (j *jsiiProxy_ResponseType) ResponseType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"responseType",
		&returns,
	)
	return returns
}


// A custom response type to support future cases.
// Experimental.
func ResponseType_Of(type_ *string) ResponseType {
	_init_.Initialize()

	var returns ResponseType

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.ResponseType",
		"of",
		[]interface{}{type_},
		&returns,
	)

	return returns
}

func ResponseType_ACCESS_DENIED() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"ACCESS_DENIED",
		&returns,
	)
	return returns
}

func ResponseType_API_CONFIGURATION_ERROR() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"API_CONFIGURATION_ERROR",
		&returns,
	)
	return returns
}

func ResponseType_AUTHORIZER_CONFIGURATION_ERROR() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"AUTHORIZER_CONFIGURATION_ERROR",
		&returns,
	)
	return returns
}

func ResponseType_AUTHORIZER_FAILURE() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"AUTHORIZER_FAILURE",
		&returns,
	)
	return returns
}

func ResponseType_BAD_REQUEST_BODY() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"BAD_REQUEST_BODY",
		&returns,
	)
	return returns
}

func ResponseType_BAD_REQUEST_PARAMETERS() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"BAD_REQUEST_PARAMETERS",
		&returns,
	)
	return returns
}

func ResponseType_DEFAULT_4XX() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"DEFAULT_4XX",
		&returns,
	)
	return returns
}

func ResponseType_DEFAULT_5XX() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"DEFAULT_5XX",
		&returns,
	)
	return returns
}

func ResponseType_EXPIRED_TOKEN() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"EXPIRED_TOKEN",
		&returns,
	)
	return returns
}

func ResponseType_INTEGRATION_FAILURE() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"INTEGRATION_FAILURE",
		&returns,
	)
	return returns
}

func ResponseType_INTEGRATION_TIMEOUT() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"INTEGRATION_TIMEOUT",
		&returns,
	)
	return returns
}

func ResponseType_INVALID_API_KEY() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"INVALID_API_KEY",
		&returns,
	)
	return returns
}

func ResponseType_INVALID_SIGNATURE() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"INVALID_SIGNATURE",
		&returns,
	)
	return returns
}

func ResponseType_MISSING_AUTHENTICATION_TOKEN() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"MISSING_AUTHENTICATION_TOKEN",
		&returns,
	)
	return returns
}

func ResponseType_QUOTA_EXCEEDED() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"QUOTA_EXCEEDED",
		&returns,
	)
	return returns
}

func ResponseType_REQUEST_TOO_LARGE() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"REQUEST_TOO_LARGE",
		&returns,
	)
	return returns
}

func ResponseType_RESOURCE_NOT_FOUND() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"RESOURCE_NOT_FOUND",
		&returns,
	)
	return returns
}

func ResponseType_THROTTLED() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"THROTTLED",
		&returns,
	)
	return returns
}

func ResponseType_UNAUTHORIZED() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"UNAUTHORIZED",
		&returns,
	)
	return returns
}

func ResponseType_UNSUPPORTED_MEDIA_TYPE() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"UNSUPPORTED_MEDIA_TYPE",
		&returns,
	)
	return returns
}

func ResponseType_WAF_FILTERED() ResponseType {
	_init_.Initialize()
	var returns ResponseType
	_jsii_.StaticGet(
		"monocdk.aws_apigateway.ResponseType",
		"WAF_FILTERED",
		&returns,
	)
	return returns
}

// Represents a REST API in Amazon API Gateway.
//
// Use `addResource` and `addMethod` to configure the API model.
//
// By default, the API will automatically be deployed and accessible from a
// public endpoint.
// Experimental.
type RestApi interface {
	RestApiBase
	DeploymentStage() Stage
	SetDeploymentStage(val Stage)
	DomainName() DomainName
	Env() *awscdk.ResourceEnvironment
	LatestDeployment() Deployment
	Methods() *[]Method
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RestApiId() *string
	RestApiName() *string
	RestApiRootResourceId() *string
	Root() IResource
	Stack() awscdk.Stack
	Url() *string
	AddApiKey(id *string, options *ApiKeyOptions) IApiKey
	AddDomainName(id *string, options *DomainNameOptions) DomainName
	AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse
	AddModel(id *string, props *ModelOptions) Model
	AddRequestValidator(id *string, props *RequestValidatorOptions) RequestValidator
	AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ArnForExecuteApi(method *string, path *string, stage *string) *string
	ConfigureCloudWatchRole(apiResource CfnRestApi)
	ConfigureDeployment(props *RestApiBaseProps)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	UrlForPath(path *string) *string
	Validate() *[]*string
}

// The jsii proxy struct for RestApi
type jsiiProxy_RestApi struct {
	jsiiProxy_RestApiBase
}

func (j *jsiiProxy_RestApi) DeploymentStage() Stage {
	var returns Stage
	_jsii_.Get(
		j,
		"deploymentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) DomainName() DomainName {
	var returns DomainName
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) LatestDeployment() Deployment {
	var returns Deployment
	_jsii_.Get(
		j,
		"latestDeployment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) Methods() *[]Method {
	var returns *[]Method
	_jsii_.Get(
		j,
		"methods",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) RestApiName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) RestApiRootResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiRootResourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) Root() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"root",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApi) Url() *string {
	var returns *string
	_jsii_.Get(
		j,
		"url",
		&returns,
	)
	return returns
}


// Experimental.
func NewRestApi(scope constructs.Construct, id *string, props *RestApiProps) RestApi {
	_init_.Initialize()

	j := jsiiProxy_RestApi{}

	_jsii_.Create(
		"monocdk.aws_apigateway.RestApi",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewRestApi_Override(r RestApi, scope constructs.Construct, id *string, props *RestApiProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.RestApi",
		[]interface{}{scope, id, props},
		r,
	)
}

func (j *jsiiProxy_RestApi) SetDeploymentStage(val Stage) {
	_jsii_.Set(
		j,
		"deploymentStage",
		val,
	)
}

// Import an existing RestApi that can be configured with additional Methods and Resources.
// Experimental.
func RestApi_FromRestApiAttributes(scope constructs.Construct, id *string, attrs *RestApiAttributes) IRestApi {
	_init_.Initialize()

	var returns IRestApi

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RestApi",
		"fromRestApiAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Import an existing RestApi.
// Experimental.
func RestApi_FromRestApiId(scope constructs.Construct, id *string, restApiId *string) IRestApi {
	_init_.Initialize()

	var returns IRestApi

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RestApi",
		"fromRestApiId",
		[]interface{}{scope, id, restApiId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func RestApi_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RestApi",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func RestApi_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RestApi",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add an ApiKey.
// Experimental.
func (r *jsiiProxy_RestApi) AddApiKey(id *string, options *ApiKeyOptions) IApiKey {
	var returns IApiKey

	_jsii_.Invoke(
		r,
		"addApiKey",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an API Gateway domain name and maps it to this API.
// Experimental.
func (r *jsiiProxy_RestApi) AddDomainName(id *string, options *DomainNameOptions) DomainName {
	var returns DomainName

	_jsii_.Invoke(
		r,
		"addDomainName",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a new gateway response.
// Experimental.
func (r *jsiiProxy_RestApi) AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse {
	var returns GatewayResponse

	_jsii_.Invoke(
		r,
		"addGatewayResponse",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a new model.
// Experimental.
func (r *jsiiProxy_RestApi) AddModel(id *string, props *ModelOptions) Model {
	var returns Model

	_jsii_.Invoke(
		r,
		"addModel",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds a new request validator.
// Experimental.
func (r *jsiiProxy_RestApi) AddRequestValidator(id *string, props *RequestValidatorOptions) RequestValidator {
	var returns RequestValidator

	_jsii_.Invoke(
		r,
		"addRequestValidator",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Adds a usage plan.
// Experimental.
func (r *jsiiProxy_RestApi) AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan {
	var returns UsagePlan

	_jsii_.Invoke(
		r,
		"addUsagePlan",
		[]interface{}{id, props},
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
func (r *jsiiProxy_RestApi) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Gets the "execute-api" ARN.
// Experimental.
func (r *jsiiProxy_RestApi) ArnForExecuteApi(method *string, path *string, stage *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"arnForExecuteApi",
		[]interface{}{method, path, stage},
		&returns,
	)

	return returns
}

// Deprecated: This method will be made internal. No replacement
func (r *jsiiProxy_RestApi) ConfigureCloudWatchRole(apiResource CfnRestApi) {
	_jsii_.InvokeVoid(
		r,
		"configureCloudWatchRole",
		[]interface{}{apiResource},
	)
}

// Deprecated: This method will be made internal. No replacement
func (r *jsiiProxy_RestApi) ConfigureDeployment(props *RestApiBaseProps) {
	_jsii_.InvokeVoid(
		r,
		"configureDeployment",
		[]interface{}{props},
	)
}

// Experimental.
func (r *jsiiProxy_RestApi) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_RestApi) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_RestApi) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns the given named metric for this API.
// Experimental.
func (r *jsiiProxy_RestApi) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the API cache in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApi) MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricCacheHitCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the backend in a given period, when API caching is enabled.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApi) MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricCacheMissCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of client-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApi) MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricClientError",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the total number API requests in a given period.
//
// Default: sample count over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApi) MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the time between when API Gateway relays a request to the backend and when it receives a response from the backend.
//
// Default: average over 5 minutes.
// Experimental.
func (r *jsiiProxy_RestApi) MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricIntegrationLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time between when API Gateway receives a request from a client and when it returns a response to the client.
//
// The latency includes the integration latency and other API Gateway overhead.
//
// Default: average over 5 minutes.
// Experimental.
func (r *jsiiProxy_RestApi) MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of server-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApi) MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricServerError",
		[]interface{}{props},
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
func (r *jsiiProxy_RestApi) OnPrepare() {
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
func (r *jsiiProxy_RestApi) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_RestApi) OnValidate() *[]*string {
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
func (r *jsiiProxy_RestApi) Prepare() {
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
func (r *jsiiProxy_RestApi) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_RestApi) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns the URL for an HTTP path.
//
// Fails if `deploymentStage` is not set either by `deploy` or explicitly.
// Experimental.
func (r *jsiiProxy_RestApi) UrlForPath(path *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"urlForPath",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Performs validation of the REST API.
// Experimental.
func (r *jsiiProxy_RestApi) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Attributes that can be specified when importing a RestApi.
// Experimental.
type RestApiAttributes struct {
	// The ID of the API Gateway RestApi.
	// Experimental.
	RestApiId *string `json:"restApiId"`
	// The resource ID of the root resource.
	// Experimental.
	RootResourceId *string `json:"rootResourceId"`
}

// Base implementation that are common to various implementations of IRestApi.
// Experimental.
type RestApiBase interface {
	awscdk.Resource
	IRestApi
	DeploymentStage() Stage
	SetDeploymentStage(val Stage)
	DomainName() DomainName
	Env() *awscdk.ResourceEnvironment
	LatestDeployment() Deployment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RestApiId() *string
	RestApiName() *string
	RestApiRootResourceId() *string
	Root() IResource
	Stack() awscdk.Stack
	AddApiKey(id *string, options *ApiKeyOptions) IApiKey
	AddDomainName(id *string, options *DomainNameOptions) DomainName
	AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse
	AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ArnForExecuteApi(method *string, path *string, stage *string) *string
	ConfigureCloudWatchRole(apiResource CfnRestApi)
	ConfigureDeployment(props *RestApiBaseProps)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	UrlForPath(path *string) *string
	Validate() *[]*string
}

// The jsii proxy struct for RestApiBase
type jsiiProxy_RestApiBase struct {
	internal.Type__awscdkResource
	jsiiProxy_IRestApi
}

func (j *jsiiProxy_RestApiBase) DeploymentStage() Stage {
	var returns Stage
	_jsii_.Get(
		j,
		"deploymentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) DomainName() DomainName {
	var returns DomainName
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) LatestDeployment() Deployment {
	var returns Deployment
	_jsii_.Get(
		j,
		"latestDeployment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) RestApiName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) RestApiRootResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiRootResourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) Root() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"root",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_RestApiBase) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewRestApiBase_Override(r RestApiBase, scope constructs.Construct, id *string, props *RestApiBaseProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.RestApiBase",
		[]interface{}{scope, id, props},
		r,
	)
}

func (j *jsiiProxy_RestApiBase) SetDeploymentStage(val Stage) {
	_jsii_.Set(
		j,
		"deploymentStage",
		val,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func RestApiBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RestApiBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func RestApiBase_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.RestApiBase",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add an ApiKey.
// Experimental.
func (r *jsiiProxy_RestApiBase) AddApiKey(id *string, options *ApiKeyOptions) IApiKey {
	var returns IApiKey

	_jsii_.Invoke(
		r,
		"addApiKey",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an API Gateway domain name and maps it to this API.
// Experimental.
func (r *jsiiProxy_RestApiBase) AddDomainName(id *string, options *DomainNameOptions) DomainName {
	var returns DomainName

	_jsii_.Invoke(
		r,
		"addDomainName",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a new gateway response.
// Experimental.
func (r *jsiiProxy_RestApiBase) AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse {
	var returns GatewayResponse

	_jsii_.Invoke(
		r,
		"addGatewayResponse",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a usage plan.
// Experimental.
func (r *jsiiProxy_RestApiBase) AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan {
	var returns UsagePlan

	_jsii_.Invoke(
		r,
		"addUsagePlan",
		[]interface{}{id, props},
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
func (r *jsiiProxy_RestApiBase) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		r,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Gets the "execute-api" ARN.
// Experimental.
func (r *jsiiProxy_RestApiBase) ArnForExecuteApi(method *string, path *string, stage *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"arnForExecuteApi",
		[]interface{}{method, path, stage},
		&returns,
	)

	return returns
}

// Deprecated: This method will be made internal. No replacement
func (r *jsiiProxy_RestApiBase) ConfigureCloudWatchRole(apiResource CfnRestApi) {
	_jsii_.InvokeVoid(
		r,
		"configureCloudWatchRole",
		[]interface{}{apiResource},
	)
}

// Deprecated: This method will be made internal. No replacement
func (r *jsiiProxy_RestApiBase) ConfigureDeployment(props *RestApiBaseProps) {
	_jsii_.InvokeVoid(
		r,
		"configureDeployment",
		[]interface{}{props},
	)
}

// Experimental.
func (r *jsiiProxy_RestApiBase) GeneratePhysicalName() *string {
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
func (r *jsiiProxy_RestApiBase) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (r *jsiiProxy_RestApiBase) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns the given named metric for this API.
// Experimental.
func (r *jsiiProxy_RestApiBase) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the API cache in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApiBase) MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricCacheHitCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the backend in a given period, when API caching is enabled.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApiBase) MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricCacheMissCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of client-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApiBase) MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricClientError",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the total number API requests in a given period.
//
// Default: sample count over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApiBase) MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the time between when API Gateway relays a request to the backend and when it receives a response from the backend.
//
// Default: average over 5 minutes.
// Experimental.
func (r *jsiiProxy_RestApiBase) MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricIntegrationLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time between when API Gateway receives a request from a client and when it returns a response to the client.
//
// The latency includes the integration latency and other API Gateway overhead.
//
// Default: average over 5 minutes.
// Experimental.
func (r *jsiiProxy_RestApiBase) MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of server-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (r *jsiiProxy_RestApiBase) MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		r,
		"metricServerError",
		[]interface{}{props},
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
func (r *jsiiProxy_RestApiBase) OnPrepare() {
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
func (r *jsiiProxy_RestApiBase) OnSynthesize(session constructs.ISynthesisSession) {
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
func (r *jsiiProxy_RestApiBase) OnValidate() *[]*string {
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
func (r *jsiiProxy_RestApiBase) Prepare() {
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
func (r *jsiiProxy_RestApiBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		r,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (r *jsiiProxy_RestApiBase) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns the URL for an HTTP path.
//
// Fails if `deploymentStage` is not set either by `deploy` or explicitly.
// Experimental.
func (r *jsiiProxy_RestApiBase) UrlForPath(path *string) *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"urlForPath",
		[]interface{}{path},
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
func (r *jsiiProxy_RestApiBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		r,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents the props that all Rest APIs share.
// Experimental.
type RestApiBaseProps struct {
	// Automatically configure an AWS CloudWatch role for API Gateway.
	// Experimental.
	CloudWatchRole *bool `json:"cloudWatchRole"`
	// Indicates if a Deployment should be automatically created for this API, and recreated when the API model (resources, methods) changes.
	//
	// Since API Gateway deployments are immutable, When this option is enabled
	// (by default), an AWS::ApiGateway::Deployment resource will automatically
	// created with a logical ID that hashes the API model (methods, resources
	// and options). This means that when the model changes, the logical ID of
	// this CloudFormation resource will change, and a new deployment will be
	// created.
	//
	// If this is set, `latestDeployment` will refer to the `Deployment` object
	// and `deploymentStage` will refer to a `Stage` that points to this
	// deployment. To customize the stage options, use the `deployOptions`
	// property.
	//
	// A CloudFormation Output will also be defined with the root URL endpoint
	// of this REST API.
	// Experimental.
	Deploy *bool `json:"deploy"`
	// Options for the API Gateway stage that will always point to the latest deployment when `deploy` is enabled.
	//
	// If `deploy` is disabled,
	// this value cannot be set.
	// Experimental.
	DeployOptions *StageOptions `json:"deployOptions"`
	// Specifies whether clients can invoke the API using the default execute-api endpoint.
	//
	// To require that clients use a custom domain name to invoke the
	// API, disable the default endpoint.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-restapi.html
	//
	// Experimental.
	DisableExecuteApiEndpoint *bool `json:"disableExecuteApiEndpoint"`
	// Configure a custom domain name and map it to this API.
	// Experimental.
	DomainName *DomainNameOptions `json:"domainName"`
	// Export name for the CfnOutput containing the API endpoint.
	// Experimental.
	EndpointExportName *string `json:"endpointExportName"`
	// A list of the endpoint types of the API.
	//
	// Use this property when creating
	// an API.
	// Experimental.
	EndpointTypes *[]EndpointType `json:"endpointTypes"`
	// Indicates whether to roll back the resource if a warning occurs while API Gateway is creating the RestApi resource.
	// Experimental.
	FailOnWarnings *bool `json:"failOnWarnings"`
	// Custom header parameters for the request.
	// See: https://docs.aws.amazon.com/cli/latest/reference/apigateway/import-rest-api.html
	//
	// Experimental.
	Parameters *map[string]*string `json:"parameters"`
	// A policy document that contains the permissions for this RestApi.
	// Experimental.
	Policy awsiam.PolicyDocument `json:"policy"`
	// A name for the API Gateway RestApi resource.
	// Experimental.
	RestApiName *string `json:"restApiName"`
	// Retains old deployment resources when the API changes.
	//
	// This allows
	// manually reverting stages to point to old deployments via the AWS
	// Console.
	// Experimental.
	RetainDeployments *bool `json:"retainDeployments"`
}

// Represents the props that all Rest APIs share.
// Deprecated: - superseded by `RestApiBaseProps`
type RestApiOptions struct {
	// Automatically configure an AWS CloudWatch role for API Gateway.
	// Deprecated: - superseded by `RestApiBaseProps`
	CloudWatchRole *bool `json:"cloudWatchRole"`
	// Indicates if a Deployment should be automatically created for this API, and recreated when the API model (resources, methods) changes.
	//
	// Since API Gateway deployments are immutable, When this option is enabled
	// (by default), an AWS::ApiGateway::Deployment resource will automatically
	// created with a logical ID that hashes the API model (methods, resources
	// and options). This means that when the model changes, the logical ID of
	// this CloudFormation resource will change, and a new deployment will be
	// created.
	//
	// If this is set, `latestDeployment` will refer to the `Deployment` object
	// and `deploymentStage` will refer to a `Stage` that points to this
	// deployment. To customize the stage options, use the `deployOptions`
	// property.
	//
	// A CloudFormation Output will also be defined with the root URL endpoint
	// of this REST API.
	// Deprecated: - superseded by `RestApiBaseProps`
	Deploy *bool `json:"deploy"`
	// Options for the API Gateway stage that will always point to the latest deployment when `deploy` is enabled.
	//
	// If `deploy` is disabled,
	// this value cannot be set.
	// Deprecated: - superseded by `RestApiBaseProps`
	DeployOptions *StageOptions `json:"deployOptions"`
	// Specifies whether clients can invoke the API using the default execute-api endpoint.
	//
	// To require that clients use a custom domain name to invoke the
	// API, disable the default endpoint.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-restapi.html
	//
	// Deprecated: - superseded by `RestApiBaseProps`
	DisableExecuteApiEndpoint *bool `json:"disableExecuteApiEndpoint"`
	// Configure a custom domain name and map it to this API.
	// Deprecated: - superseded by `RestApiBaseProps`
	DomainName *DomainNameOptions `json:"domainName"`
	// Export name for the CfnOutput containing the API endpoint.
	// Deprecated: - superseded by `RestApiBaseProps`
	EndpointExportName *string `json:"endpointExportName"`
	// A list of the endpoint types of the API.
	//
	// Use this property when creating
	// an API.
	// Deprecated: - superseded by `RestApiBaseProps`
	EndpointTypes *[]EndpointType `json:"endpointTypes"`
	// Indicates whether to roll back the resource if a warning occurs while API Gateway is creating the RestApi resource.
	// Deprecated: - superseded by `RestApiBaseProps`
	FailOnWarnings *bool `json:"failOnWarnings"`
	// Custom header parameters for the request.
	// See: https://docs.aws.amazon.com/cli/latest/reference/apigateway/import-rest-api.html
	//
	// Deprecated: - superseded by `RestApiBaseProps`
	Parameters *map[string]*string `json:"parameters"`
	// A policy document that contains the permissions for this RestApi.
	// Deprecated: - superseded by `RestApiBaseProps`
	Policy awsiam.PolicyDocument `json:"policy"`
	// A name for the API Gateway RestApi resource.
	// Deprecated: - superseded by `RestApiBaseProps`
	RestApiName *string `json:"restApiName"`
	// Retains old deployment resources when the API changes.
	//
	// This allows
	// manually reverting stages to point to old deployments via the AWS
	// Console.
	// Deprecated: - superseded by `RestApiBaseProps`
	RetainDeployments *bool `json:"retainDeployments"`
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Deprecated: - superseded by `RestApiBaseProps`
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Deprecated: - superseded by `RestApiBaseProps`
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Deprecated: - superseded by `RestApiBaseProps`
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
}

// Props to create a new instance of RestApi.
// Experimental.
type RestApiProps struct {
	// Automatically configure an AWS CloudWatch role for API Gateway.
	// Experimental.
	CloudWatchRole *bool `json:"cloudWatchRole"`
	// Indicates if a Deployment should be automatically created for this API, and recreated when the API model (resources, methods) changes.
	//
	// Since API Gateway deployments are immutable, When this option is enabled
	// (by default), an AWS::ApiGateway::Deployment resource will automatically
	// created with a logical ID that hashes the API model (methods, resources
	// and options). This means that when the model changes, the logical ID of
	// this CloudFormation resource will change, and a new deployment will be
	// created.
	//
	// If this is set, `latestDeployment` will refer to the `Deployment` object
	// and `deploymentStage` will refer to a `Stage` that points to this
	// deployment. To customize the stage options, use the `deployOptions`
	// property.
	//
	// A CloudFormation Output will also be defined with the root URL endpoint
	// of this REST API.
	// Experimental.
	Deploy *bool `json:"deploy"`
	// Options for the API Gateway stage that will always point to the latest deployment when `deploy` is enabled.
	//
	// If `deploy` is disabled,
	// this value cannot be set.
	// Experimental.
	DeployOptions *StageOptions `json:"deployOptions"`
	// Specifies whether clients can invoke the API using the default execute-api endpoint.
	//
	// To require that clients use a custom domain name to invoke the
	// API, disable the default endpoint.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-restapi.html
	//
	// Experimental.
	DisableExecuteApiEndpoint *bool `json:"disableExecuteApiEndpoint"`
	// Configure a custom domain name and map it to this API.
	// Experimental.
	DomainName *DomainNameOptions `json:"domainName"`
	// Export name for the CfnOutput containing the API endpoint.
	// Experimental.
	EndpointExportName *string `json:"endpointExportName"`
	// A list of the endpoint types of the API.
	//
	// Use this property when creating
	// an API.
	// Experimental.
	EndpointTypes *[]EndpointType `json:"endpointTypes"`
	// Indicates whether to roll back the resource if a warning occurs while API Gateway is creating the RestApi resource.
	// Experimental.
	FailOnWarnings *bool `json:"failOnWarnings"`
	// Custom header parameters for the request.
	// See: https://docs.aws.amazon.com/cli/latest/reference/apigateway/import-rest-api.html
	//
	// Experimental.
	Parameters *map[string]*string `json:"parameters"`
	// A policy document that contains the permissions for this RestApi.
	// Experimental.
	Policy awsiam.PolicyDocument `json:"policy"`
	// A name for the API Gateway RestApi resource.
	// Experimental.
	RestApiName *string `json:"restApiName"`
	// Retains old deployment resources when the API changes.
	//
	// This allows
	// manually reverting stages to point to old deployments via the AWS
	// Console.
	// Experimental.
	RetainDeployments *bool `json:"retainDeployments"`
	// Adds a CORS preflight OPTIONS method to this resource and all child resources.
	//
	// You can add CORS at the resource-level using `addCorsPreflight`.
	// Experimental.
	DefaultCorsPreflightOptions *CorsOptions `json:"defaultCorsPreflightOptions"`
	// An integration to use as a default for all methods created within this API unless an integration is specified.
	// Experimental.
	DefaultIntegration Integration `json:"defaultIntegration"`
	// Method options to use as a default for all methods created within this API unless custom options are specified.
	// Experimental.
	DefaultMethodOptions *MethodOptions `json:"defaultMethodOptions"`
	// The source of the API key for metering requests according to a usage plan.
	// Experimental.
	ApiKeySourceType ApiKeySourceType `json:"apiKeySourceType"`
	// The list of binary media mime-types that are supported by the RestApi resource, such as "image/png" or "application/octet-stream".
	// Experimental.
	BinaryMediaTypes *[]*string `json:"binaryMediaTypes"`
	// The ID of the API Gateway RestApi resource that you want to clone.
	// Experimental.
	CloneFrom IRestApi `json:"cloneFrom"`
	// A description of the purpose of this API Gateway RestApi resource.
	// Experimental.
	Description *string `json:"description"`
	// The EndpointConfiguration property type specifies the endpoint types of a REST API.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-endpointconfiguration.html
	//
	// Experimental.
	EndpointConfiguration *EndpointConfiguration `json:"endpointConfiguration"`
	// A nullable integer that is used to enable compression (with non-negative between 0 and 10485760 (10M) bytes, inclusive) or disable compression (when undefined) on an API.
	//
	// When compression is enabled, compression or
	// decompression is not applied on the payload if the payload size is
	// smaller than this value. Setting it to zero allows compression for any
	// payload size.
	// Experimental.
	MinimumCompressionSize *float64 `json:"minimumCompressionSize"`
}

// OpenAPI specification from an S3 archive.
// Experimental.
type S3ApiDefinition interface {
	ApiDefinition
	Bind(_scope awscdk.Construct) *ApiDefinitionConfig
}

// The jsii proxy struct for S3ApiDefinition
type jsiiProxy_S3ApiDefinition struct {
	jsiiProxy_ApiDefinition
}

// Experimental.
func NewS3ApiDefinition(bucket awss3.IBucket, key *string, objectVersion *string) S3ApiDefinition {
	_init_.Initialize()

	j := jsiiProxy_S3ApiDefinition{}

	_jsii_.Create(
		"monocdk.aws_apigateway.S3ApiDefinition",
		[]interface{}{bucket, key, objectVersion},
		&j,
	)

	return &j
}

// Experimental.
func NewS3ApiDefinition_Override(s S3ApiDefinition, bucket awss3.IBucket, key *string, objectVersion *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.S3ApiDefinition",
		[]interface{}{bucket, key, objectVersion},
		s,
	)
}

// Loads the API specification from a local disk asset.
// Experimental.
func S3ApiDefinition_FromAsset(file *string, options *awss3assets.AssetOptions) AssetApiDefinition {
	_init_.Initialize()

	var returns AssetApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.S3ApiDefinition",
		"fromAsset",
		[]interface{}{file, options},
		&returns,
	)

	return returns
}

// Creates an API definition from a specification file in an S3 bucket.
// Experimental.
func S3ApiDefinition_FromBucket(bucket awss3.IBucket, key *string, objectVersion *string) S3ApiDefinition {
	_init_.Initialize()

	var returns S3ApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.S3ApiDefinition",
		"fromBucket",
		[]interface{}{bucket, key, objectVersion},
		&returns,
	)

	return returns
}

// Create an API definition from an inline object.
//
// The inline object must follow the
// schema of OpenAPI 2.0 or OpenAPI 3.0
//
// TODO: EXAMPLE
//
// Experimental.
func S3ApiDefinition_FromInline(definition interface{}) InlineApiDefinition {
	_init_.Initialize()

	var returns InlineApiDefinition

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.S3ApiDefinition",
		"fromInline",
		[]interface{}{definition},
		&returns,
	)

	return returns
}

// Called when the specification is initialized to allow this object to bind to the stack, add resources and have fun.
// Experimental.
func (s *jsiiProxy_S3ApiDefinition) Bind(_scope awscdk.Construct) *ApiDefinitionConfig {
	var returns *ApiDefinitionConfig

	_jsii_.Invoke(
		s,
		"bind",
		[]interface{}{_scope},
		&returns,
	)

	return returns
}

// The minimum version of the SSL protocol that you want API Gateway to use for HTTPS connections.
// Experimental.
type SecurityPolicy string

const (
	SecurityPolicy_TLS_1_0 SecurityPolicy = "TLS_1_0"
	SecurityPolicy_TLS_1_2 SecurityPolicy = "TLS_1_2"
)

// Represents a REST API in Amazon API Gateway, created with an OpenAPI specification.
//
// Some properties normally accessible on @see {@link RestApi} - such as the description -
// must be declared in the specification. All Resources and Methods need to be defined as
// part of the OpenAPI specification file, and cannot be added via the CDK.
//
// By default, the API will automatically be deployed and accessible from a
// public endpoint.
// Experimental.
type SpecRestApi interface {
	RestApiBase
	DeploymentStage() Stage
	SetDeploymentStage(val Stage)
	DomainName() DomainName
	Env() *awscdk.ResourceEnvironment
	LatestDeployment() Deployment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RestApiId() *string
	RestApiName() *string
	RestApiRootResourceId() *string
	Root() IResource
	Stack() awscdk.Stack
	AddApiKey(id *string, options *ApiKeyOptions) IApiKey
	AddDomainName(id *string, options *DomainNameOptions) DomainName
	AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse
	AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	ArnForExecuteApi(method *string, path *string, stage *string) *string
	ConfigureCloudWatchRole(apiResource CfnRestApi)
	ConfigureDeployment(props *RestApiBaseProps)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	UrlForPath(path *string) *string
	Validate() *[]*string
}

// The jsii proxy struct for SpecRestApi
type jsiiProxy_SpecRestApi struct {
	jsiiProxy_RestApiBase
}

func (j *jsiiProxy_SpecRestApi) DeploymentStage() Stage {
	var returns Stage
	_jsii_.Get(
		j,
		"deploymentStage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) DomainName() DomainName {
	var returns DomainName
	_jsii_.Get(
		j,
		"domainName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) LatestDeployment() Deployment {
	var returns Deployment
	_jsii_.Get(
		j,
		"latestDeployment",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) RestApiName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) RestApiRootResourceId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiRootResourceId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) Root() IResource {
	var returns IResource
	_jsii_.Get(
		j,
		"root",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_SpecRestApi) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewSpecRestApi(scope constructs.Construct, id *string, props *SpecRestApiProps) SpecRestApi {
	_init_.Initialize()

	j := jsiiProxy_SpecRestApi{}

	_jsii_.Create(
		"monocdk.aws_apigateway.SpecRestApi",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewSpecRestApi_Override(s SpecRestApi, scope constructs.Construct, id *string, props *SpecRestApiProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.SpecRestApi",
		[]interface{}{scope, id, props},
		s,
	)
}

func (j *jsiiProxy_SpecRestApi) SetDeploymentStage(val Stage) {
	_jsii_.Set(
		j,
		"deploymentStage",
		val,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func SpecRestApi_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.SpecRestApi",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func SpecRestApi_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.SpecRestApi",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add an ApiKey.
// Experimental.
func (s *jsiiProxy_SpecRestApi) AddApiKey(id *string, options *ApiKeyOptions) IApiKey {
	var returns IApiKey

	_jsii_.Invoke(
		s,
		"addApiKey",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Defines an API Gateway domain name and maps it to this API.
// Experimental.
func (s *jsiiProxy_SpecRestApi) AddDomainName(id *string, options *DomainNameOptions) DomainName {
	var returns DomainName

	_jsii_.Invoke(
		s,
		"addDomainName",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a new gateway response.
// Experimental.
func (s *jsiiProxy_SpecRestApi) AddGatewayResponse(id *string, options *GatewayResponseOptions) GatewayResponse {
	var returns GatewayResponse

	_jsii_.Invoke(
		s,
		"addGatewayResponse",
		[]interface{}{id, options},
		&returns,
	)

	return returns
}

// Adds a usage plan.
// Experimental.
func (s *jsiiProxy_SpecRestApi) AddUsagePlan(id *string, props *UsagePlanProps) UsagePlan {
	var returns UsagePlan

	_jsii_.Invoke(
		s,
		"addUsagePlan",
		[]interface{}{id, props},
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
func (s *jsiiProxy_SpecRestApi) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Gets the "execute-api" ARN.
// Experimental.
func (s *jsiiProxy_SpecRestApi) ArnForExecuteApi(method *string, path *string, stage *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"arnForExecuteApi",
		[]interface{}{method, path, stage},
		&returns,
	)

	return returns
}

// Deprecated: This method will be made internal. No replacement
func (s *jsiiProxy_SpecRestApi) ConfigureCloudWatchRole(apiResource CfnRestApi) {
	_jsii_.InvokeVoid(
		s,
		"configureCloudWatchRole",
		[]interface{}{apiResource},
	)
}

// Deprecated: This method will be made internal. No replacement
func (s *jsiiProxy_SpecRestApi) ConfigureDeployment(props *RestApiBaseProps) {
	_jsii_.InvokeVoid(
		s,
		"configureDeployment",
		[]interface{}{props},
	)
}

// Experimental.
func (s *jsiiProxy_SpecRestApi) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_SpecRestApi) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_SpecRestApi) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns the given named metric for this API.
// Experimental.
func (s *jsiiProxy_SpecRestApi) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the API cache in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (s *jsiiProxy_SpecRestApi) MetricCacheHitCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricCacheHitCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of requests served from the backend in a given period, when API caching is enabled.
//
// Default: sum over 5 minutes
// Experimental.
func (s *jsiiProxy_SpecRestApi) MetricCacheMissCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricCacheMissCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of client-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (s *jsiiProxy_SpecRestApi) MetricClientError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricClientError",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the total number API requests in a given period.
//
// Default: sample count over 5 minutes
// Experimental.
func (s *jsiiProxy_SpecRestApi) MetricCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the time between when API Gateway relays a request to the backend and when it receives a response from the backend.
//
// Default: average over 5 minutes.
// Experimental.
func (s *jsiiProxy_SpecRestApi) MetricIntegrationLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricIntegrationLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time between when API Gateway receives a request from a client and when it returns a response to the client.
//
// The latency includes the integration latency and other API Gateway overhead.
//
// Default: average over 5 minutes.
// Experimental.
func (s *jsiiProxy_SpecRestApi) MetricLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Metric for the number of server-side errors captured in a given period.
//
// Default: sum over 5 minutes
// Experimental.
func (s *jsiiProxy_SpecRestApi) MetricServerError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		s,
		"metricServerError",
		[]interface{}{props},
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
func (s *jsiiProxy_SpecRestApi) OnPrepare() {
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
func (s *jsiiProxy_SpecRestApi) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_SpecRestApi) OnValidate() *[]*string {
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
func (s *jsiiProxy_SpecRestApi) Prepare() {
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
func (s *jsiiProxy_SpecRestApi) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_SpecRestApi) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns the URL for an HTTP path.
//
// Fails if `deploymentStage` is not set either by `deploy` or explicitly.
// Experimental.
func (s *jsiiProxy_SpecRestApi) UrlForPath(path *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"urlForPath",
		[]interface{}{path},
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
func (s *jsiiProxy_SpecRestApi) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Props to instantiate a new SpecRestApi.
// Experimental.
type SpecRestApiProps struct {
	// Automatically configure an AWS CloudWatch role for API Gateway.
	// Experimental.
	CloudWatchRole *bool `json:"cloudWatchRole"`
	// Indicates if a Deployment should be automatically created for this API, and recreated when the API model (resources, methods) changes.
	//
	// Since API Gateway deployments are immutable, When this option is enabled
	// (by default), an AWS::ApiGateway::Deployment resource will automatically
	// created with a logical ID that hashes the API model (methods, resources
	// and options). This means that when the model changes, the logical ID of
	// this CloudFormation resource will change, and a new deployment will be
	// created.
	//
	// If this is set, `latestDeployment` will refer to the `Deployment` object
	// and `deploymentStage` will refer to a `Stage` that points to this
	// deployment. To customize the stage options, use the `deployOptions`
	// property.
	//
	// A CloudFormation Output will also be defined with the root URL endpoint
	// of this REST API.
	// Experimental.
	Deploy *bool `json:"deploy"`
	// Options for the API Gateway stage that will always point to the latest deployment when `deploy` is enabled.
	//
	// If `deploy` is disabled,
	// this value cannot be set.
	// Experimental.
	DeployOptions *StageOptions `json:"deployOptions"`
	// Specifies whether clients can invoke the API using the default execute-api endpoint.
	//
	// To require that clients use a custom domain name to invoke the
	// API, disable the default endpoint.
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-restapi.html
	//
	// Experimental.
	DisableExecuteApiEndpoint *bool `json:"disableExecuteApiEndpoint"`
	// Configure a custom domain name and map it to this API.
	// Experimental.
	DomainName *DomainNameOptions `json:"domainName"`
	// Export name for the CfnOutput containing the API endpoint.
	// Experimental.
	EndpointExportName *string `json:"endpointExportName"`
	// A list of the endpoint types of the API.
	//
	// Use this property when creating
	// an API.
	// Experimental.
	EndpointTypes *[]EndpointType `json:"endpointTypes"`
	// Indicates whether to roll back the resource if a warning occurs while API Gateway is creating the RestApi resource.
	// Experimental.
	FailOnWarnings *bool `json:"failOnWarnings"`
	// Custom header parameters for the request.
	// See: https://docs.aws.amazon.com/cli/latest/reference/apigateway/import-rest-api.html
	//
	// Experimental.
	Parameters *map[string]*string `json:"parameters"`
	// A policy document that contains the permissions for this RestApi.
	// Experimental.
	Policy awsiam.PolicyDocument `json:"policy"`
	// A name for the API Gateway RestApi resource.
	// Experimental.
	RestApiName *string `json:"restApiName"`
	// Retains old deployment resources when the API changes.
	//
	// This allows
	// manually reverting stages to point to old deployments via the AWS
	// Console.
	// Experimental.
	RetainDeployments *bool `json:"retainDeployments"`
	// An OpenAPI definition compatible with API Gateway.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-import-api.html
	//
	// Experimental.
	ApiDefinition ApiDefinition `json:"apiDefinition"`
}

// Experimental.
type Stage interface {
	awscdk.Resource
	IStage
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RestApi() IRestApi
	Stack() awscdk.Stack
	StageName() *string
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
	UrlForPath(path *string) *string
	Validate() *[]*string
}

// The jsii proxy struct for Stage
type jsiiProxy_Stage struct {
	internal.Type__awscdkResource
	jsiiProxy_IStage
}

func (j *jsiiProxy_Stage) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) RestApi() IRestApi {
	var returns IRestApi
	_jsii_.Get(
		j,
		"restApi",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stage) StageName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stageName",
		&returns,
	)
	return returns
}


// Experimental.
func NewStage(scope constructs.Construct, id *string, props *StageProps) Stage {
	_init_.Initialize()

	j := jsiiProxy_Stage{}

	_jsii_.Create(
		"monocdk.aws_apigateway.Stage",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewStage_Override(s Stage, scope constructs.Construct, id *string, props *StageProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.Stage",
		[]interface{}{scope, id, props},
		s,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func Stage_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Stage",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func Stage_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.Stage",
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
func (s *jsiiProxy_Stage) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		s,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (s *jsiiProxy_Stage) GeneratePhysicalName() *string {
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
func (s *jsiiProxy_Stage) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (s *jsiiProxy_Stage) GetResourceNameAttribute(nameAttr *string) *string {
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
func (s *jsiiProxy_Stage) OnPrepare() {
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
func (s *jsiiProxy_Stage) OnSynthesize(session constructs.ISynthesisSession) {
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
func (s *jsiiProxy_Stage) OnValidate() *[]*string {
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
func (s *jsiiProxy_Stage) Prepare() {
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
func (s *jsiiProxy_Stage) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		s,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (s *jsiiProxy_Stage) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns the invoke URL for a certain path.
// Experimental.
func (s *jsiiProxy_Stage) UrlForPath(path *string) *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"urlForPath",
		[]interface{}{path},
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
func (s *jsiiProxy_Stage) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		s,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Experimental.
type StageOptions struct {
	// Indicates whether the cached responses are encrypted.
	// Experimental.
	CacheDataEncrypted *bool `json:"cacheDataEncrypted"`
	// Specifies the time to live (TTL), in seconds, for cached responses.
	//
	// The
	// higher the TTL, the longer the response will be cached.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-caching.html
	//
	// Experimental.
	CacheTtl awscdk.Duration `json:"cacheTtl"`
	// Specifies whether responses should be cached and returned for requests.
	//
	// A
	// cache cluster must be enabled on the stage for responses to be cached.
	// Experimental.
	CachingEnabled *bool `json:"cachingEnabled"`
	// Specifies whether data trace logging is enabled for this method, which effects the log entries pushed to Amazon CloudWatch Logs.
	// Experimental.
	DataTraceEnabled *bool `json:"dataTraceEnabled"`
	// Specifies the logging level for this method, which effects the log entries pushed to Amazon CloudWatch Logs.
	// Experimental.
	LoggingLevel MethodLoggingLevel `json:"loggingLevel"`
	// Specifies whether Amazon CloudWatch metrics are enabled for this method.
	// Experimental.
	MetricsEnabled *bool `json:"metricsEnabled"`
	// Specifies the throttling burst limit.
	//
	// The total rate of all requests in your AWS account is limited to 5,000 requests.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html
	//
	// Experimental.
	ThrottlingBurstLimit *float64 `json:"throttlingBurstLimit"`
	// Specifies the throttling rate limit.
	//
	// The total rate of all requests in your AWS account is limited to 10,000 requests per second (rps).
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html
	//
	// Experimental.
	ThrottlingRateLimit *float64 `json:"throttlingRateLimit"`
	// The CloudWatch Logs log group.
	// Experimental.
	AccessLogDestination IAccessLogDestination `json:"accessLogDestination"`
	// A single line format of access logs of data, as specified by selected $content variables.
	//
	// The format must include at least `AccessLogFormat.contextRequestId()`.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html#context-variable-reference
	//
	// Experimental.
	AccessLogFormat AccessLogFormat `json:"accessLogFormat"`
	// Indicates whether cache clustering is enabled for the stage.
	// Experimental.
	CacheClusterEnabled *bool `json:"cacheClusterEnabled"`
	// The stage's cache cluster size.
	// Experimental.
	CacheClusterSize *string `json:"cacheClusterSize"`
	// The identifier of the client certificate that API Gateway uses to call your integration endpoints in the stage.
	// Experimental.
	ClientCertificateId *string `json:"clientCertificateId"`
	// A description of the purpose of the stage.
	// Experimental.
	Description *string `json:"description"`
	// The version identifier of the API documentation snapshot.
	// Experimental.
	DocumentationVersion *string `json:"documentationVersion"`
	// Method deployment options for specific resources/methods.
	//
	// These will
	// override common options defined in `StageOptions#methodOptions`.
	// Experimental.
	MethodOptions *map[string]*MethodDeploymentOptions `json:"methodOptions"`
	// The name of the stage, which API Gateway uses as the first path segment in the invoked Uniform Resource Identifier (URI).
	// Experimental.
	StageName *string `json:"stageName"`
	// Specifies whether Amazon X-Ray tracing is enabled for this method.
	// Experimental.
	TracingEnabled *bool `json:"tracingEnabled"`
	// A map that defines the stage variables.
	//
	// Variable names must consist of
	// alphanumeric characters, and the values must match the following regular
	// expression: [A-Za-z0-9-._~:/?#&amp;=,]+.
	// Experimental.
	Variables *map[string]*string `json:"variables"`
}

// Experimental.
type StageProps struct {
	// Indicates whether the cached responses are encrypted.
	// Experimental.
	CacheDataEncrypted *bool `json:"cacheDataEncrypted"`
	// Specifies the time to live (TTL), in seconds, for cached responses.
	//
	// The
	// higher the TTL, the longer the response will be cached.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-caching.html
	//
	// Experimental.
	CacheTtl awscdk.Duration `json:"cacheTtl"`
	// Specifies whether responses should be cached and returned for requests.
	//
	// A
	// cache cluster must be enabled on the stage for responses to be cached.
	// Experimental.
	CachingEnabled *bool `json:"cachingEnabled"`
	// Specifies whether data trace logging is enabled for this method, which effects the log entries pushed to Amazon CloudWatch Logs.
	// Experimental.
	DataTraceEnabled *bool `json:"dataTraceEnabled"`
	// Specifies the logging level for this method, which effects the log entries pushed to Amazon CloudWatch Logs.
	// Experimental.
	LoggingLevel MethodLoggingLevel `json:"loggingLevel"`
	// Specifies whether Amazon CloudWatch metrics are enabled for this method.
	// Experimental.
	MetricsEnabled *bool `json:"metricsEnabled"`
	// Specifies the throttling burst limit.
	//
	// The total rate of all requests in your AWS account is limited to 5,000 requests.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html
	//
	// Experimental.
	ThrottlingBurstLimit *float64 `json:"throttlingBurstLimit"`
	// Specifies the throttling rate limit.
	//
	// The total rate of all requests in your AWS account is limited to 10,000 requests per second (rps).
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html
	//
	// Experimental.
	ThrottlingRateLimit *float64 `json:"throttlingRateLimit"`
	// The CloudWatch Logs log group.
	// Experimental.
	AccessLogDestination IAccessLogDestination `json:"accessLogDestination"`
	// A single line format of access logs of data, as specified by selected $content variables.
	//
	// The format must include at least `AccessLogFormat.contextRequestId()`.
	// See: https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html#context-variable-reference
	//
	// Experimental.
	AccessLogFormat AccessLogFormat `json:"accessLogFormat"`
	// Indicates whether cache clustering is enabled for the stage.
	// Experimental.
	CacheClusterEnabled *bool `json:"cacheClusterEnabled"`
	// The stage's cache cluster size.
	// Experimental.
	CacheClusterSize *string `json:"cacheClusterSize"`
	// The identifier of the client certificate that API Gateway uses to call your integration endpoints in the stage.
	// Experimental.
	ClientCertificateId *string `json:"clientCertificateId"`
	// A description of the purpose of the stage.
	// Experimental.
	Description *string `json:"description"`
	// The version identifier of the API documentation snapshot.
	// Experimental.
	DocumentationVersion *string `json:"documentationVersion"`
	// Method deployment options for specific resources/methods.
	//
	// These will
	// override common options defined in `StageOptions#methodOptions`.
	// Experimental.
	MethodOptions *map[string]*MethodDeploymentOptions `json:"methodOptions"`
	// The name of the stage, which API Gateway uses as the first path segment in the invoked Uniform Resource Identifier (URI).
	// Experimental.
	StageName *string `json:"stageName"`
	// Specifies whether Amazon X-Ray tracing is enabled for this method.
	// Experimental.
	TracingEnabled *bool `json:"tracingEnabled"`
	// A map that defines the stage variables.
	//
	// Variable names must consist of
	// alphanumeric characters, and the values must match the following regular
	// expression: [A-Za-z0-9-._~:/?#&amp;=,]+.
	// Experimental.
	Variables *map[string]*string `json:"variables"`
	// The deployment that this stage points to [disable-awslint:ref-via-interface].
	// Experimental.
	Deployment Deployment `json:"deployment"`
}

// Container for defining throttling parameters to API stages or methods.
// Experimental.
type ThrottleSettings struct {
	// The maximum API request rate limit over a time ranging from one to a few seconds.
	// Experimental.
	BurstLimit *float64 `json:"burstLimit"`
	// The API request steady-state rate limit (average requests per second over an extended period of time).
	// Experimental.
	RateLimit *float64 `json:"rateLimit"`
}

// Represents per-method throttling for a resource.
// Experimental.
type ThrottlingPerMethod struct {
	// [disable-awslint:ref-via-interface] The method for which you specify the throttling settings.
	// Experimental.
	Method Method `json:"method"`
	// Specifies the overall request rate (average requests per second) and burst capacity.
	// Experimental.
	Throttle *ThrottleSettings `json:"throttle"`
}

// Token based lambda authorizer that recognizes the caller's identity as a bearer token, such as a JSON Web Token (JWT) or an OAuth token.
//
// Based on the token, authorization is performed by a lambda function.
// Experimental.
type TokenAuthorizer interface {
	Authorizer
	IAuthorizer
	AuthorizationType() AuthorizationType
	AuthorizerArn() *string
	AuthorizerId() *string
	Env() *awscdk.ResourceEnvironment
	Handler() awslambda.IFunction
	Node() awscdk.ConstructNode
	PhysicalName() *string
	RestApiId() *string
	SetRestApiId(val *string)
	Role() awsiam.IRole
	Stack() awscdk.Stack
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LazyRestApiId() *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	SetupPermissions()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for TokenAuthorizer
type jsiiProxy_TokenAuthorizer struct {
	jsiiProxy_Authorizer
	jsiiProxy_IAuthorizer
}

func (j *jsiiProxy_TokenAuthorizer) AuthorizationType() AuthorizationType {
	var returns AuthorizationType
	_jsii_.Get(
		j,
		"authorizationType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) AuthorizerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) AuthorizerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"authorizerId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) Handler() awslambda.IFunction {
	var returns awslambda.IFunction
	_jsii_.Get(
		j,
		"handler",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) RestApiId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"restApiId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) Role() awsiam.IRole {
	var returns awsiam.IRole
	_jsii_.Get(
		j,
		"role",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TokenAuthorizer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewTokenAuthorizer(scope constructs.Construct, id *string, props *TokenAuthorizerProps) TokenAuthorizer {
	_init_.Initialize()

	j := jsiiProxy_TokenAuthorizer{}

	_jsii_.Create(
		"monocdk.aws_apigateway.TokenAuthorizer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewTokenAuthorizer_Override(t TokenAuthorizer, scope constructs.Construct, id *string, props *TokenAuthorizerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.TokenAuthorizer",
		[]interface{}{scope, id, props},
		t,
	)
}

func (j *jsiiProxy_TokenAuthorizer) SetRestApiId(val *string) {
	_jsii_.Set(
		j,
		"restApiId",
		val,
	)
}

// Return whether the given object is an Authorizer.
// Experimental.
func TokenAuthorizer_IsAuthorizer(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.TokenAuthorizer",
		"isAuthorizer",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func TokenAuthorizer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.TokenAuthorizer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func TokenAuthorizer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.TokenAuthorizer",
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
func (t *jsiiProxy_TokenAuthorizer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		t,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (t *jsiiProxy_TokenAuthorizer) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TokenAuthorizer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TokenAuthorizer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Returns a token that resolves to the Rest Api Id at the time of synthesis.
//
// Throws an error, during token resolution, if no RestApi is attached to this authorizer.
// Experimental.
func (t *jsiiProxy_TokenAuthorizer) LazyRestApiId() *string {
	var returns *string

	_jsii_.Invoke(
		t,
		"lazyRestApiId",
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
func (t *jsiiProxy_TokenAuthorizer) OnPrepare() {
	_jsii_.InvokeVoid(
		t,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (t *jsiiProxy_TokenAuthorizer) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
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
func (t *jsiiProxy_TokenAuthorizer) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TokenAuthorizer) Prepare() {
	_jsii_.InvokeVoid(
		t,
		"prepare",
		nil, // no parameters
	)
}

// Sets up the permissions necessary for the API Gateway service to invoke the Lambda function.
// Experimental.
func (t *jsiiProxy_TokenAuthorizer) SetupPermissions() {
	_jsii_.InvokeVoid(
		t,
		"setupPermissions",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (t *jsiiProxy_TokenAuthorizer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TokenAuthorizer) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		t,
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
func (t *jsiiProxy_TokenAuthorizer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for TokenAuthorizer.
// Experimental.
type TokenAuthorizerProps struct {
	// The handler for the authorizer lambda function.
	//
	// The handler must follow a very specific protocol on the input it receives and the output it needs to produce.
	// API Gateway has documented the handler's input specification
	// {@link https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-input.html | here} and output specification
	// {@link https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html | here}.
	// Experimental.
	Handler awslambda.IFunction `json:"handler"`
	// An optional IAM role for APIGateway to assume before calling the Lambda-based authorizer.
	//
	// The IAM role must be
	// assumable by 'apigateway.amazonaws.com'.
	// Experimental.
	AssumeRole awsiam.IRole `json:"assumeRole"`
	// An optional human friendly name for the authorizer.
	//
	// Note that, this is not the primary identifier of the authorizer.
	// Experimental.
	AuthorizerName *string `json:"authorizerName"`
	// How long APIGateway should cache the results.
	//
	// Max 1 hour.
	// Disable caching by setting this to 0.
	// Experimental.
	ResultsCacheTtl awscdk.Duration `json:"resultsCacheTtl"`
	// The request header mapping expression for the bearer token.
	//
	// This is typically passed as part of the header, in which case
	// this should be `method.request.header.Authorizer` where Authorizer is the header containing the bearer token.
	// See: https://docs.aws.amazon.com/apigateway/api-reference/link-relation/authorizer-create/#identitySource
	//
	// Experimental.
	IdentitySource *string `json:"identitySource"`
	// An optional regex to be matched against the authorization token.
	//
	// When matched the authorizer lambda is invoked,
	// otherwise a 401 Unauthorized is returned to the client.
	// Experimental.
	ValidationRegex *string `json:"validationRegex"`
}

// Experimental.
type UsagePlan interface {
	awscdk.Resource
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	UsagePlanId() *string
	AddApiKey(apiKey IApiKey, options *AddApiKeyOptions)
	AddApiStage(apiStage *UsagePlanPerApiStage)
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

// The jsii proxy struct for UsagePlan
type jsiiProxy_UsagePlan struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_UsagePlan) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UsagePlan) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UsagePlan) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UsagePlan) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_UsagePlan) UsagePlanId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"usagePlanId",
		&returns,
	)
	return returns
}


// Experimental.
func NewUsagePlan(scope constructs.Construct, id *string, props *UsagePlanProps) UsagePlan {
	_init_.Initialize()

	j := jsiiProxy_UsagePlan{}

	_jsii_.Create(
		"monocdk.aws_apigateway.UsagePlan",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewUsagePlan_Override(u UsagePlan, scope constructs.Construct, id *string, props *UsagePlanProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.UsagePlan",
		[]interface{}{scope, id, props},
		u,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func UsagePlan_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.UsagePlan",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func UsagePlan_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.UsagePlan",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Adds an ApiKey.
// Experimental.
func (u *jsiiProxy_UsagePlan) AddApiKey(apiKey IApiKey, options *AddApiKeyOptions) {
	_jsii_.InvokeVoid(
		u,
		"addApiKey",
		[]interface{}{apiKey, options},
	)
}

// Adds an apiStage.
// Experimental.
func (u *jsiiProxy_UsagePlan) AddApiStage(apiStage *UsagePlanPerApiStage) {
	_jsii_.InvokeVoid(
		u,
		"addApiStage",
		[]interface{}{apiStage},
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
func (u *jsiiProxy_UsagePlan) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		u,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (u *jsiiProxy_UsagePlan) GeneratePhysicalName() *string {
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
func (u *jsiiProxy_UsagePlan) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (u *jsiiProxy_UsagePlan) GetResourceNameAttribute(nameAttr *string) *string {
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
func (u *jsiiProxy_UsagePlan) OnPrepare() {
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
func (u *jsiiProxy_UsagePlan) OnSynthesize(session constructs.ISynthesisSession) {
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
func (u *jsiiProxy_UsagePlan) OnValidate() *[]*string {
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
func (u *jsiiProxy_UsagePlan) Prepare() {
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
func (u *jsiiProxy_UsagePlan) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		u,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (u *jsiiProxy_UsagePlan) ToString() *string {
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
func (u *jsiiProxy_UsagePlan) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		u,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Represents the API stages that a usage plan applies to.
// Experimental.
type UsagePlanPerApiStage struct {
	// Experimental.
	Api IRestApi `json:"api"`
	// [disable-awslint:ref-via-interface].
	// Experimental.
	Stage Stage `json:"stage"`
	// Experimental.
	Throttle *[]*ThrottlingPerMethod `json:"throttle"`
}

// Experimental.
type UsagePlanProps struct {
	// ApiKey to be associated with the usage plan.
	// Deprecated: use `addApiKey()`
	ApiKey IApiKey `json:"apiKey"`
	// API Stages to be associated with the usage plan.
	// Experimental.
	ApiStages *[]*UsagePlanPerApiStage `json:"apiStages"`
	// Represents usage plan purpose.
	// Experimental.
	Description *string `json:"description"`
	// Name for this usage plan.
	// Experimental.
	Name *string `json:"name"`
	// Number of requests clients can make in a given time period.
	// Experimental.
	Quota *QuotaSettings `json:"quota"`
	// Overall throttle settings for the API.
	// Experimental.
	Throttle *ThrottleSettings `json:"throttle"`
}

// Define a new VPC Link Specifies an API Gateway VPC link for a RestApi to access resources in an Amazon Virtual Private Cloud (VPC).
// Experimental.
type VpcLink interface {
	awscdk.Resource
	IVpcLink
	Env() *awscdk.ResourceEnvironment
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	VpcLinkId() *string
	AddTargets(targets ...awselasticloadbalancingv2.INetworkLoadBalancer)
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

// The jsii proxy struct for VpcLink
type jsiiProxy_VpcLink struct {
	internal.Type__awscdkResource
	jsiiProxy_IVpcLink
}

func (j *jsiiProxy_VpcLink) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_VpcLink) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_VpcLink) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_VpcLink) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_VpcLink) VpcLinkId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"vpcLinkId",
		&returns,
	)
	return returns
}


// Experimental.
func NewVpcLink(scope constructs.Construct, id *string, props *VpcLinkProps) VpcLink {
	_init_.Initialize()

	j := jsiiProxy_VpcLink{}

	_jsii_.Create(
		"monocdk.aws_apigateway.VpcLink",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewVpcLink_Override(v VpcLink, scope constructs.Construct, id *string, props *VpcLinkProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_apigateway.VpcLink",
		[]interface{}{scope, id, props},
		v,
	)
}

// Import a VPC Link by its Id.
// Experimental.
func VpcLink_FromVpcLinkId(scope constructs.Construct, id *string, vpcLinkId *string) IVpcLink {
	_init_.Initialize()

	var returns IVpcLink

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.VpcLink",
		"fromVpcLinkId",
		[]interface{}{scope, id, vpcLinkId},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func VpcLink_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.VpcLink",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func VpcLink_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_apigateway.VpcLink",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Experimental.
func (v *jsiiProxy_VpcLink) AddTargets(targets ...awselasticloadbalancingv2.INetworkLoadBalancer) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		v,
		"addTargets",
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
func (v *jsiiProxy_VpcLink) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		v,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (v *jsiiProxy_VpcLink) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_VpcLink) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_VpcLink) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_VpcLink) OnPrepare() {
	_jsii_.InvokeVoid(
		v,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (v *jsiiProxy_VpcLink) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		v,
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
func (v *jsiiProxy_VpcLink) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_VpcLink) Prepare() {
	_jsii_.InvokeVoid(
		v,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (v *jsiiProxy_VpcLink) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		v,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (v *jsiiProxy_VpcLink) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		v,
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
func (v *jsiiProxy_VpcLink) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		v,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a VpcLink.
// Experimental.
type VpcLinkProps struct {
	// The description of the VPC link.
	// Experimental.
	Description *string `json:"description"`
	// The network load balancers of the VPC targeted by the VPC link.
	//
	// The network load balancers must be owned by the same AWS account of the API owner.
	// Experimental.
	Targets *[]awselasticloadbalancingv2.INetworkLoadBalancer `json:"targets"`
	// The name used to label and identify the VPC link.
	// Experimental.
	VpcLinkName *string `json:"vpcLinkName"`
}


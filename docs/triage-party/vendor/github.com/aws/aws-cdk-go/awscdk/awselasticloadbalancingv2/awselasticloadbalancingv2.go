package awselasticloadbalancingv2

import (
	_init_ "github.com/aws/aws-cdk-go/awscdk/jsii"
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancingv2/internal"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/constructs-go/constructs/v3"
)

// Properties for adding a new action to a listener.
// Experimental.
type AddApplicationActionProps struct {
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Experimental.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Priority of this target group.
	//
	// The rule with the lowest priority will be used for every request.
	// If priority is not given, these target groups will be added as
	// defaults, and must not have conditions.
	//
	// Priorities must be unique.
	// Experimental.
	Priority *float64 `json:"priority"`
	// Action to perform.
	// Experimental.
	Action ListenerAction `json:"action"`
}

// Properties for adding a new target group to a listener.
// Experimental.
type AddApplicationTargetGroupsProps struct {
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Experimental.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Priority of this target group.
	//
	// The rule with the lowest priority will be used for every request.
	// If priority is not given, these target groups will be added as
	// defaults, and must not have conditions.
	//
	// Priorities must be unique.
	// Experimental.
	Priority *float64 `json:"priority"`
	// Target groups to forward requests to.
	// Experimental.
	TargetGroups *[]IApplicationTargetGroup `json:"targetGroups"`
}

// Properties for adding new targets to a listener.
// Experimental.
type AddApplicationTargetsProps struct {
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Experimental.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Priority of this target group.
	//
	// The rule with the lowest priority will be used for every request.
	// If priority is not given, these target groups will be added as
	// defaults, and must not have conditions.
	//
	// Priorities must be unique.
	// Experimental.
	Priority *float64 `json:"priority"`
	// The amount of time for Elastic Load Balancing to wait before deregistering a target.
	//
	// The range is 0-3600 seconds.
	// Experimental.
	DeregistrationDelay awscdk.Duration `json:"deregistrationDelay"`
	// Health check configuration.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// The protocol to use.
	// Experimental.
	Protocol ApplicationProtocol `json:"protocol"`
	// The protocol version to use.
	// Experimental.
	ProtocolVersion ApplicationProtocolVersion `json:"protocolVersion"`
	// The time period during which the load balancer sends a newly registered target a linearly increasing share of the traffic to the target group.
	//
	// The range is 30-900 seconds (15 minutes).
	// Experimental.
	SlowStart awscdk.Duration `json:"slowStart"`
	// The stickiness cookie expiration period.
	//
	// Setting this value enables load balancer stickiness.
	//
	// After this period, the cookie is considered stale. The minimum value is
	// 1 second and the maximum value is 7 days (604800 seconds).
	// Experimental.
	StickinessCookieDuration awscdk.Duration `json:"stickinessCookieDuration"`
	// The name of an application-based stickiness cookie.
	//
	// Names that start with the following prefixes are not allowed: AWSALB, AWSALBAPP,
	// and AWSALBTG; they're reserved for use by the load balancer.
	//
	// Note: `stickinessCookieName` parameter depends on the presence of `stickinessCookieDuration` parameter.
	// If `stickinessCookieDuration` is not set, `stickinessCookieName` will be omitted.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/sticky-sessions.html
	//
	// Experimental.
	StickinessCookieName *string `json:"stickinessCookieName"`
	// The name of the target group.
	//
	// This name must be unique per region per account, can have a maximum of
	// 32 characters, must contain only alphanumeric characters or hyphens, and
	// must not begin or end with a hyphen.
	// Experimental.
	TargetGroupName *string `json:"targetGroupName"`
	// The targets to add to this target group.
	//
	// Can be `Instance`, `IPAddress`, or any self-registering load balancing
	// target. All target must be of the same type.
	// Experimental.
	Targets *[]IApplicationLoadBalancerTarget `json:"targets"`
}

// Properties for adding a fixed response to a listener.
// Deprecated: Use `ApplicationListener.addAction` instead.
type AddFixedResponseProps struct {
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Priority of this target group.
	//
	// The rule with the lowest priority will be used for every request.
	// If priority is not given, these target groups will be added as
	// defaults, and must not have conditions.
	//
	// Priorities must be unique.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Priority *float64 `json:"priority"`
	// The HTTP response code (2XX, 4XX or 5XX).
	// Deprecated: Use `ApplicationListener.addAction` instead.
	StatusCode *string `json:"statusCode"`
	// The content type.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	ContentType ContentType `json:"contentType"`
	// The message.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	MessageBody *string `json:"messageBody"`
}

// Properties for adding a new action to a listener.
// Experimental.
type AddNetworkActionProps struct {
	// Action to perform.
	// Experimental.
	Action NetworkListenerAction `json:"action"`
}

// Properties for adding new network targets to a listener.
// Experimental.
type AddNetworkTargetsProps struct {
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// The amount of time for Elastic Load Balancing to wait before deregistering a target.
	//
	// The range is 0-3600 seconds.
	// Experimental.
	DeregistrationDelay awscdk.Duration `json:"deregistrationDelay"`
	// Health check configuration.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// Indicates whether client IP preservation is enabled.
	// Experimental.
	PreserveClientIp *bool `json:"preserveClientIp"`
	// Protocol for target group, expects TCP, TLS, UDP, or TCP_UDP.
	// Experimental.
	Protocol Protocol `json:"protocol"`
	// Indicates whether Proxy Protocol version 2 is enabled.
	// Experimental.
	ProxyProtocolV2 *bool `json:"proxyProtocolV2"`
	// The name of the target group.
	//
	// This name must be unique per region per account, can have a maximum of
	// 32 characters, must contain only alphanumeric characters or hyphens, and
	// must not begin or end with a hyphen.
	// Experimental.
	TargetGroupName *string `json:"targetGroupName"`
	// The targets to add to this target group.
	//
	// Can be `Instance`, `IPAddress`, or any self-registering load balancing
	// target. If you use either `Instance` or `IPAddress` as targets, all
	// target must be of the same type.
	// Experimental.
	Targets *[]INetworkLoadBalancerTarget `json:"targets"`
}

// Properties for adding a redirect response to a listener.
// Deprecated: Use `ApplicationListener.addAction` instead.
type AddRedirectResponseProps struct {
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Priority of this target group.
	//
	// The rule with the lowest priority will be used for every request.
	// If priority is not given, these target groups will be added as
	// defaults, and must not have conditions.
	//
	// Priorities must be unique.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Priority *float64 `json:"priority"`
	// The HTTP redirect code (HTTP_301 or HTTP_302).
	// Deprecated: Use `ApplicationListener.addAction` instead.
	StatusCode *string `json:"statusCode"`
	// The hostname.
	//
	// This component is not percent-encoded. The hostname can contain #{host}.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Host *string `json:"host"`
	// The absolute path, starting with the leading "/".
	//
	// This component is not percent-encoded.
	// The path can contain #{host}, #{path}, and #{port}.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Path *string `json:"path"`
	// The port.
	//
	// You can specify a value from 1 to 65535 or #{port}.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Port *string `json:"port"`
	// The protocol.
	//
	// You can specify HTTP, HTTPS, or #{protocol}. You can redirect HTTP to HTTP,
	// HTTP to HTTPS, and HTTPS to HTTPS. You cannot redirect HTTPS to HTTP.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Protocol *string `json:"protocol"`
	// The query parameters, URL-encoded when necessary, but not percent-encoded.
	//
	// Do not include the leading "?", as it is automatically added.
	// You can specify any of the reserved keywords.
	// Deprecated: Use `ApplicationListener.addAction` instead.
	Query *string `json:"query"`
}

// Properties for adding a conditional load balancing rule.
// Experimental.
type AddRuleProps struct {
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Experimental.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// May contain up to three '*' wildcards.
	//
	// Requires that priority is set.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Priority of this target group.
	//
	// The rule with the lowest priority will be used for every request.
	// If priority is not given, these target groups will be added as
	// defaults, and must not have conditions.
	//
	// Priorities must be unique.
	// Experimental.
	Priority *float64 `json:"priority"`
}

// Define an ApplicationListener.
// Experimental.
type ApplicationListener interface {
	BaseListener
	IApplicationListener
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	ListenerArn() *string
	LoadBalancer() IApplicationLoadBalancer
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddAction(id *string, props *AddApplicationActionProps)
	AddCertificateArns(id *string, arns *[]*string)
	AddCertificates(id *string, certificates *[]IListenerCertificate)
	AddFixedResponse(id *string, props *AddFixedResponseProps)
	AddRedirectResponse(id *string, props *AddRedirectResponseProps)
	AddTargetGroups(id *string, props *AddApplicationTargetGroupsProps)
	AddTargets(id *string, props *AddApplicationTargetsProps) ApplicationTargetGroup
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ApplicationListener
type jsiiProxy_ApplicationListener struct {
	jsiiProxy_BaseListener
	jsiiProxy_IApplicationListener
}

func (j *jsiiProxy_ApplicationListener) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationListener) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationListener) ListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationListener) LoadBalancer() IApplicationLoadBalancer {
	var returns IApplicationLoadBalancer
	_jsii_.Get(
		j,
		"loadBalancer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationListener) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationListener) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationListener) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewApplicationListener(scope constructs.Construct, id *string, props *ApplicationListenerProps) ApplicationListener {
	_init_.Initialize()

	j := jsiiProxy_ApplicationListener{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListener",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewApplicationListener_Override(a ApplicationListener, scope constructs.Construct, id *string, props *ApplicationListenerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListener",
		[]interface{}{scope, id, props},
		a,
	)
}

// Import an existing listener.
// Experimental.
func ApplicationListener_FromApplicationListenerAttributes(scope constructs.Construct, id *string, attrs *ApplicationListenerAttributes) IApplicationListener {
	_init_.Initialize()

	var returns IApplicationListener

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListener",
		"fromApplicationListenerAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Look up an ApplicationListener.
// Experimental.
func ApplicationListener_FromLookup(scope constructs.Construct, id *string, options *ApplicationListenerLookupOptions) IApplicationListener {
	_init_.Initialize()

	var returns IApplicationListener

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListener",
		"fromLookup",
		[]interface{}{scope, id, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ApplicationListener_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListener",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ApplicationListener_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListener",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Perform the given default action on incoming requests.
//
// This allows full control of the default action of the load balancer,
// including Action chaining, fixed responses and redirect responses. See
// the `ListenerAction` class for all options.
//
// It's possible to add routing conditions to the Action added in this way.
// At least one Action must be added without conditions (which becomes the
// default Action).
// Experimental.
func (a *jsiiProxy_ApplicationListener) AddAction(id *string, props *AddApplicationActionProps) {
	_jsii_.InvokeVoid(
		a,
		"addAction",
		[]interface{}{id, props},
	)
}

// Add one or more certificates to this listener.
//
// After the first certificate, this creates ApplicationListenerCertificates
// resources since cloudformation requires the certificates array on the
// listener resource to have a length of 1.
// Deprecated: Use `addCertificates` instead.
func (a *jsiiProxy_ApplicationListener) AddCertificateArns(id *string, arns *[]*string) {
	_jsii_.InvokeVoid(
		a,
		"addCertificateArns",
		[]interface{}{id, arns},
	)
}

// Add one or more certificates to this listener.
//
// After the first certificate, this creates ApplicationListenerCertificates
// resources since cloudformation requires the certificates array on the
// listener resource to have a length of 1.
// Experimental.
func (a *jsiiProxy_ApplicationListener) AddCertificates(id *string, certificates *[]IListenerCertificate) {
	_jsii_.InvokeVoid(
		a,
		"addCertificates",
		[]interface{}{id, certificates},
	)
}

// Add a fixed response.
// Deprecated: Use `addAction()` instead
func (a *jsiiProxy_ApplicationListener) AddFixedResponse(id *string, props *AddFixedResponseProps) {
	_jsii_.InvokeVoid(
		a,
		"addFixedResponse",
		[]interface{}{id, props},
	)
}

// Add a redirect response.
// Deprecated: Use `addAction()` instead
func (a *jsiiProxy_ApplicationListener) AddRedirectResponse(id *string, props *AddRedirectResponseProps) {
	_jsii_.InvokeVoid(
		a,
		"addRedirectResponse",
		[]interface{}{id, props},
	)
}

// Load balance incoming requests to the given target groups.
//
// All target groups will be load balanced to with equal weight and without
// stickiness. For a more complex configuration than that, use `addAction()`.
//
// It's possible to add routing conditions to the TargetGroups added in this
// way. At least one TargetGroup must be added without conditions (which will
// become the default Action for this listener).
// Experimental.
func (a *jsiiProxy_ApplicationListener) AddTargetGroups(id *string, props *AddApplicationTargetGroupsProps) {
	_jsii_.InvokeVoid(
		a,
		"addTargetGroups",
		[]interface{}{id, props},
	)
}

// Load balance incoming requests to the given load balancing targets.
//
// This method implicitly creates an ApplicationTargetGroup for the targets
// involved, and a 'forward' action to route traffic to the given TargetGroup.
//
// If you want more control over the precise setup, create the TargetGroup
// and use `addAction` yourself.
//
// It's possible to add conditions to the targets added in this way. At least
// one set of targets must be added without conditions.
//
// Returns: The newly created target group
// Experimental.
func (a *jsiiProxy_ApplicationListener) AddTargets(id *string, props *AddApplicationTargetsProps) ApplicationTargetGroup {
	var returns ApplicationTargetGroup

	_jsii_.Invoke(
		a,
		"addTargets",
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
func (a *jsiiProxy_ApplicationListener) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_ApplicationListener) GeneratePhysicalName() *string {
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
func (a *jsiiProxy_ApplicationListener) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (a *jsiiProxy_ApplicationListener) GetResourceNameAttribute(nameAttr *string) *string {
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
func (a *jsiiProxy_ApplicationListener) OnPrepare() {
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
func (a *jsiiProxy_ApplicationListener) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_ApplicationListener) OnValidate() *[]*string {
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
func (a *jsiiProxy_ApplicationListener) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Register that a connectable that has been added to this load balancer.
//
// Don't call this directly. It is called by ApplicationTargetGroup.
// Experimental.
func (a *jsiiProxy_ApplicationListener) RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port) {
	_jsii_.InvokeVoid(
		a,
		"registerConnectable",
		[]interface{}{connectable, portRange},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_ApplicationListener) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_ApplicationListener) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate this listener.
// Experimental.
func (a *jsiiProxy_ApplicationListener) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to reference an existing listener.
// Experimental.
type ApplicationListenerAttributes struct {
	// ARN of the listener.
	// Experimental.
	ListenerArn *string `json:"listenerArn"`
	// The default port on which this listener is listening.
	// Experimental.
	DefaultPort *float64 `json:"defaultPort"`
	// Security group of the load balancer this listener is associated with.
	// Experimental.
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
	// Whether the imported security group allows all outbound traffic or not when imported using `securityGroupId`.
	//
	// Unless set to `false`, no egress rules will be added to the security group.
	// Deprecated: use `securityGroup` instead
	SecurityGroupAllowsAllOutbound *bool `json:"securityGroupAllowsAllOutbound"`
	// Security group ID of the load balancer this listener is associated with.
	// Deprecated: use `securityGroup` instead
	SecurityGroupId *string `json:"securityGroupId"`
}

// Add certificates to a listener.
// Experimental.
type ApplicationListenerCertificate interface {
	awscdk.Construct
	Node() awscdk.ConstructNode
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ApplicationListenerCertificate
type jsiiProxy_ApplicationListenerCertificate struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_ApplicationListenerCertificate) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


// Experimental.
func NewApplicationListenerCertificate(scope constructs.Construct, id *string, props *ApplicationListenerCertificateProps) ApplicationListenerCertificate {
	_init_.Initialize()

	j := jsiiProxy_ApplicationListenerCertificate{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListenerCertificate",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewApplicationListenerCertificate_Override(a ApplicationListenerCertificate, scope constructs.Construct, id *string, props *ApplicationListenerCertificateProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListenerCertificate",
		[]interface{}{scope, id, props},
		a,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func ApplicationListenerCertificate_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListenerCertificate",
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
func (a *jsiiProxy_ApplicationListenerCertificate) OnPrepare() {
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
func (a *jsiiProxy_ApplicationListenerCertificate) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_ApplicationListenerCertificate) OnValidate() *[]*string {
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
func (a *jsiiProxy_ApplicationListenerCertificate) Prepare() {
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
func (a *jsiiProxy_ApplicationListenerCertificate) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_ApplicationListenerCertificate) ToString() *string {
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
func (a *jsiiProxy_ApplicationListenerCertificate) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for adding a set of certificates to a listener.
// Experimental.
type ApplicationListenerCertificateProps struct {
	// The listener to attach the rule to.
	// Experimental.
	Listener IApplicationListener `json:"listener"`
	// ARNs of certificates to attach.
	//
	// Duplicates are not allowed.
	// Deprecated: Use `certificates` instead.
	CertificateArns *[]*string `json:"certificateArns"`
	// Certificates to attach.
	//
	// Duplicates are not allowed.
	// Experimental.
	Certificates *[]IListenerCertificate `json:"certificates"`
}

// Options for ApplicationListener lookup.
// Experimental.
type ApplicationListenerLookupOptions struct {
	// Filter listeners by listener port.
	// Experimental.
	ListenerPort *float64 `json:"listenerPort"`
	// Filter listeners by associated load balancer arn.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// Filter listeners by associated load balancer tags.
	// Experimental.
	LoadBalancerTags *map[string]*string `json:"loadBalancerTags"`
	// ARN of the listener to look up.
	// Experimental.
	ListenerArn *string `json:"listenerArn"`
	// Filter listeners by listener protocol.
	// Experimental.
	ListenerProtocol ApplicationProtocol `json:"listenerProtocol"`
}

// Properties for defining a standalone ApplicationListener.
// Experimental.
type ApplicationListenerProps struct {
	// The certificates to use on this listener.
	// Deprecated: Use the `certificates` property instead
	CertificateArns *[]*string `json:"certificateArns"`
	// Certificate list of ACM cert ARNs.
	// Experimental.
	Certificates *[]IListenerCertificate `json:"certificates"`
	// Default action to take for requests to this listener.
	//
	// This allows full control of the default action of the load balancer,
	// including Action chaining, fixed responses and redirect responses.
	//
	// See the `ListenerAction` class for all options.
	//
	// Cannot be specified together with `defaultTargetGroups`.
	// Experimental.
	DefaultAction ListenerAction `json:"defaultAction"`
	// Default target groups to load balance to.
	//
	// All target groups will be load balanced to with equal weight and without
	// stickiness. For a more complex configuration than that, use
	// either `defaultAction` or `addAction()`.
	//
	// Cannot be specified together with `defaultAction`.
	// Experimental.
	DefaultTargetGroups *[]IApplicationTargetGroup `json:"defaultTargetGroups"`
	// Allow anyone to connect to this listener.
	//
	// If this is specified, the listener will be opened up to anyone who can reach it.
	// For internal load balancers this is anyone in the same VPC. For public load
	// balancers, this is anyone on the internet.
	//
	// If you want to be more selective about who can access this load
	// balancer, set this to `false` and use the listener's `connections`
	// object to selectively grant access to the listener.
	// Experimental.
	Open *bool `json:"open"`
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// The protocol to use.
	// Experimental.
	Protocol ApplicationProtocol `json:"protocol"`
	// The security policy that defines which ciphers and protocols are supported.
	// Experimental.
	SslPolicy SslPolicy `json:"sslPolicy"`
	// The load balancer to attach this listener to.
	// Experimental.
	LoadBalancer IApplicationLoadBalancer `json:"loadBalancer"`
}

// Define a new listener rule.
// Experimental.
type ApplicationListenerRule interface {
	awscdk.Construct
	ListenerRuleArn() *string
	Node() awscdk.ConstructNode
	AddCondition(condition ListenerCondition)
	AddFixedResponse(fixedResponse *FixedResponse)
	AddRedirectResponse(redirectResponse *RedirectResponse)
	AddTargetGroup(targetGroup IApplicationTargetGroup)
	ConfigureAction(action ListenerAction)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	SetCondition(field *string, values *[]*string)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ApplicationListenerRule
type jsiiProxy_ApplicationListenerRule struct {
	internal.Type__awscdkConstruct
}

func (j *jsiiProxy_ApplicationListenerRule) ListenerRuleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerRuleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationListenerRule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


// Experimental.
func NewApplicationListenerRule(scope constructs.Construct, id *string, props *ApplicationListenerRuleProps) ApplicationListenerRule {
	_init_.Initialize()

	j := jsiiProxy_ApplicationListenerRule{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListenerRule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewApplicationListenerRule_Override(a ApplicationListenerRule, scope constructs.Construct, id *string, props *ApplicationListenerRuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListenerRule",
		[]interface{}{scope, id, props},
		a,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func ApplicationListenerRule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationListenerRule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Add a non-standard condition to this rule.
// Experimental.
func (a *jsiiProxy_ApplicationListenerRule) AddCondition(condition ListenerCondition) {
	_jsii_.InvokeVoid(
		a,
		"addCondition",
		[]interface{}{condition},
	)
}

// Add a fixed response.
// Deprecated: Use configureAction instead
func (a *jsiiProxy_ApplicationListenerRule) AddFixedResponse(fixedResponse *FixedResponse) {
	_jsii_.InvokeVoid(
		a,
		"addFixedResponse",
		[]interface{}{fixedResponse},
	)
}

// Add a redirect response.
// Deprecated: Use configureAction instead
func (a *jsiiProxy_ApplicationListenerRule) AddRedirectResponse(redirectResponse *RedirectResponse) {
	_jsii_.InvokeVoid(
		a,
		"addRedirectResponse",
		[]interface{}{redirectResponse},
	)
}

// Add a TargetGroup to load balance to.
// Deprecated: Use configureAction instead
func (a *jsiiProxy_ApplicationListenerRule) AddTargetGroup(targetGroup IApplicationTargetGroup) {
	_jsii_.InvokeVoid(
		a,
		"addTargetGroup",
		[]interface{}{targetGroup},
	)
}

// Configure the action to perform for this rule.
// Experimental.
func (a *jsiiProxy_ApplicationListenerRule) ConfigureAction(action ListenerAction) {
	_jsii_.InvokeVoid(
		a,
		"configureAction",
		[]interface{}{action},
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
func (a *jsiiProxy_ApplicationListenerRule) OnPrepare() {
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
func (a *jsiiProxy_ApplicationListenerRule) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_ApplicationListenerRule) OnValidate() *[]*string {
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
func (a *jsiiProxy_ApplicationListenerRule) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Add a non-standard condition to this rule.
//
// If the condition conflicts with an already set condition, it will be overwritten by the one you specified.
// Deprecated: use `addCondition` instead.
func (a *jsiiProxy_ApplicationListenerRule) SetCondition(field *string, values *[]*string) {
	_jsii_.InvokeVoid(
		a,
		"setCondition",
		[]interface{}{field, values},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_ApplicationListenerRule) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_ApplicationListenerRule) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate the rule.
// Experimental.
func (a *jsiiProxy_ApplicationListenerRule) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining a listener rule.
// Experimental.
type ApplicationListenerRuleProps struct {
	// Priority of the rule.
	//
	// The rule with the lowest priority will be used for every request.
	//
	// Priorities must be unique.
	// Experimental.
	Priority *float64 `json:"priority"`
	// Action to perform when requests are received.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	// Experimental.
	Action ListenerAction `json:"action"`
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Experimental.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Fixed response to return.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	// Deprecated: Use `action` instead.
	FixedResponse *FixedResponse `json:"fixedResponse"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// Paths may contain up to three '*' wildcards.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Redirect response to return.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	// Deprecated: Use `action` instead.
	RedirectResponse *RedirectResponse `json:"redirectResponse"`
	// Target groups to forward requests to.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	//
	// Implies a `forward` action.
	// Experimental.
	TargetGroups *[]IApplicationTargetGroup `json:"targetGroups"`
	// The listener to attach the rule to.
	// Experimental.
	Listener IApplicationListener `json:"listener"`
}

// Define an Application Load Balancer.
// Experimental.
type ApplicationLoadBalancer interface {
	BaseLoadBalancer
	IApplicationLoadBalancer
	Connections() awsec2.Connections
	Env() *awscdk.ResourceEnvironment
	IpAddressType() IpAddressType
	LoadBalancerArn() *string
	LoadBalancerCanonicalHostedZoneId() *string
	LoadBalancerDnsName() *string
	LoadBalancerFullName() *string
	LoadBalancerName() *string
	LoadBalancerSecurityGroups() *[]*string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	Vpc() awsec2.IVpc
	AddListener(id *string, props *BaseApplicationListenerProps) ApplicationListener
	AddRedirect(props *ApplicationLoadBalancerRedirectConfig) ApplicationListener
	AddSecurityGroup(securityGroup awsec2.ISecurityGroup)
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LogAccessLogs(bucket awss3.IBucket, prefix *string)
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricActiveConnectionCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricClientTlsNegotiationErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricConsumedLCUs(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricElbAuthError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricElbAuthFailure(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricElbAuthLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricElbAuthSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHttpCodeElb(code HttpCodeElb, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHttpCodeTarget(code HttpCodeTarget, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHttpFixedResponseCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHttpRedirectCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHttpRedirectUrlLimitExceededCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIpv6ProcessedBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIpv6RequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNewConnectionCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricProcessedBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRejectedConnectionCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRuleEvaluations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTargetConnectionErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTargetResponseTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTargetTLSNegotiationErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RemoveAttribute(key *string)
	SetAttribute(key *string, value *string)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ApplicationLoadBalancer
type jsiiProxy_ApplicationLoadBalancer struct {
	jsiiProxy_BaseLoadBalancer
	jsiiProxy_IApplicationLoadBalancer
}

func (j *jsiiProxy_ApplicationLoadBalancer) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) IpAddressType() IpAddressType {
	var returns IpAddressType
	_jsii_.Get(
		j,
		"ipAddressType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) LoadBalancerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) LoadBalancerCanonicalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) LoadBalancerDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) LoadBalancerFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) LoadBalancerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) LoadBalancerSecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"loadBalancerSecurityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationLoadBalancer) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}


// Experimental.
func NewApplicationLoadBalancer(scope constructs.Construct, id *string, props *ApplicationLoadBalancerProps) ApplicationLoadBalancer {
	_init_.Initialize()

	j := jsiiProxy_ApplicationLoadBalancer{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationLoadBalancer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewApplicationLoadBalancer_Override(a ApplicationLoadBalancer, scope constructs.Construct, id *string, props *ApplicationLoadBalancerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationLoadBalancer",
		[]interface{}{scope, id, props},
		a,
	)
}

// Import an existing Application Load Balancer.
// Experimental.
func ApplicationLoadBalancer_FromApplicationLoadBalancerAttributes(scope constructs.Construct, id *string, attrs *ApplicationLoadBalancerAttributes) IApplicationLoadBalancer {
	_init_.Initialize()

	var returns IApplicationLoadBalancer

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationLoadBalancer",
		"fromApplicationLoadBalancerAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Look up an application load balancer.
// Experimental.
func ApplicationLoadBalancer_FromLookup(scope constructs.Construct, id *string, options *ApplicationLoadBalancerLookupOptions) IApplicationLoadBalancer {
	_init_.Initialize()

	var returns IApplicationLoadBalancer

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationLoadBalancer",
		"fromLookup",
		[]interface{}{scope, id, options},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ApplicationLoadBalancer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationLoadBalancer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func ApplicationLoadBalancer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationLoadBalancer",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a new listener to this load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) AddListener(id *string, props *BaseApplicationListenerProps) ApplicationListener {
	var returns ApplicationListener

	_jsii_.Invoke(
		a,
		"addListener",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

// Add a redirection listener to this load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) AddRedirect(props *ApplicationLoadBalancerRedirectConfig) ApplicationListener {
	var returns ApplicationListener

	_jsii_.Invoke(
		a,
		"addRedirect",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Add a security group to this load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) AddSecurityGroup(securityGroup awsec2.ISecurityGroup) {
	_jsii_.InvokeVoid(
		a,
		"addSecurityGroup",
		[]interface{}{securityGroup},
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
func (a *jsiiProxy_ApplicationLoadBalancer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		a,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) GeneratePhysicalName() *string {
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
func (a *jsiiProxy_ApplicationLoadBalancer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (a *jsiiProxy_ApplicationLoadBalancer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		a,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Enable access logging for this load balancer.
//
// A region must be specified on the stack containing the load balancer; you cannot enable logging on
// environment-agnostic stacks. See https://docs.aws.amazon.com/cdk/latest/guide/environments.html
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) LogAccessLogs(bucket awss3.IBucket, prefix *string) {
	_jsii_.InvokeVoid(
		a,
		"logAccessLogs",
		[]interface{}{bucket, prefix},
	)
}

// Return the given named metric for this Application Load Balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The total number of concurrent TCP connections active from clients to the load balancer and from the load balancer to targets.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricActiveConnectionCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricActiveConnectionCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of TLS connections initiated by the client that did not establish a session with the load balancer.
//
// Possible causes include a
// mismatch of ciphers or protocols.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricClientTlsNegotiationErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricClientTlsNegotiationErrorCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of load balancer capacity units (LCU) used by your load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricConsumedLCUs(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricConsumedLCUs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of user authentications that could not be completed.
//
// Because an authenticate action was misconfigured, the load balancer
// couldn't establish a connection with the IdP, or the load balancer
// couldn't complete the authentication flow due to an internal error.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricElbAuthError(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricElbAuthError",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of user authentications that could not be completed because the IdP denied access to the user or an authorization code was used more than once.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricElbAuthFailure(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricElbAuthFailure",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time elapsed, in milliseconds, to query the IdP for the ID token and user info.
//
// If one or more of these operations fail, this is the time to failure.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricElbAuthLatency(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricElbAuthLatency",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of authenticate actions that were successful.
//
// This metric is incremented at the end of the authentication workflow,
// after the load balancer has retrieved the user claims from the IdP.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricElbAuthSuccess(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricElbAuthSuccess",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of HTTP 3xx/4xx/5xx codes that originate from the load balancer.
//
// This does not include any response codes generated by the targets.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricHttpCodeElb(code HttpCodeElb, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHttpCodeElb",
		[]interface{}{code, props},
		&returns,
	)

	return returns
}

// The number of HTTP 2xx/3xx/4xx/5xx response codes generated by all targets in the load balancer.
//
// This does not include any response codes generated by the load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricHttpCodeTarget(code HttpCodeTarget, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHttpCodeTarget",
		[]interface{}{code, props},
		&returns,
	)

	return returns
}

// The number of fixed-response actions that were successful.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricHttpFixedResponseCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHttpFixedResponseCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of redirect actions that were successful.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricHttpRedirectCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHttpRedirectCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of redirect actions that couldn't be completed because the URL in the response location header is larger than 8K.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricHttpRedirectUrlLimitExceededCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHttpRedirectUrlLimitExceededCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of bytes processed by the load balancer over IPv6.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricIpv6ProcessedBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricIpv6ProcessedBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of IPv6 requests received by the load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricIpv6RequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricIpv6RequestCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of new TCP connections established from clients to the load balancer and from the load balancer to targets.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricNewConnectionCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricNewConnectionCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of bytes processed by the load balancer over IPv4 and IPv6.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricProcessedBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricProcessedBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of connections that were rejected because the load balancer had reached its maximum number of connections.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricRejectedConnectionCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricRejectedConnectionCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of requests processed over IPv4 and IPv6.
//
// This count includes only the requests with a response generated by a target of the load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricRequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricRequestCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of rules processed by the load balancer given a request rate averaged over an hour.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricRuleEvaluations(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricRuleEvaluations",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of connections that were not successfully established between the load balancer and target.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricTargetConnectionErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTargetConnectionErrorCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time elapsed, in seconds, after the request leaves the load balancer until a response from the target is received.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricTargetResponseTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTargetResponseTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of TLS connections initiated by the load balancer that did not establish a session with the target.
//
// Possible causes include a mismatch of ciphers or protocols.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) MetricTargetTLSNegotiationErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTargetTLSNegotiationErrorCount",
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
func (a *jsiiProxy_ApplicationLoadBalancer) OnPrepare() {
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
func (a *jsiiProxy_ApplicationLoadBalancer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_ApplicationLoadBalancer) OnValidate() *[]*string {
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
func (a *jsiiProxy_ApplicationLoadBalancer) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Remove an attribute from the load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) RemoveAttribute(key *string) {
	_jsii_.InvokeVoid(
		a,
		"removeAttribute",
		[]interface{}{key},
	)
}

// Set a non-standard attribute on the load balancer.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/application-load-balancers.html#load-balancer-attributes
//
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) SetAttribute(key *string, value *string) {
	_jsii_.InvokeVoid(
		a,
		"setAttribute",
		[]interface{}{key, value},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_ApplicationLoadBalancer) ToString() *string {
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
func (a *jsiiProxy_ApplicationLoadBalancer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to reference an existing load balancer.
// Experimental.
type ApplicationLoadBalancerAttributes struct {
	// ARN of the load balancer.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// ID of the load balancer's security group.
	// Experimental.
	SecurityGroupId *string `json:"securityGroupId"`
	// The canonical hosted zone ID of this load balancer.
	// Experimental.
	LoadBalancerCanonicalHostedZoneId *string `json:"loadBalancerCanonicalHostedZoneId"`
	// The DNS name of this load balancer.
	// Experimental.
	LoadBalancerDnsName *string `json:"loadBalancerDnsName"`
	// Whether the security group allows all outbound traffic or not.
	//
	// Unless set to `false`, no egress rules will be added to the security group.
	// Experimental.
	SecurityGroupAllowsAllOutbound *bool `json:"securityGroupAllowsAllOutbound"`
	// The VPC this load balancer has been created in, if available.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
}

// Options for looking up an ApplicationLoadBalancer.
// Experimental.
type ApplicationLoadBalancerLookupOptions struct {
	// Find by load balancer's ARN.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// Match load balancer tags.
	// Experimental.
	LoadBalancerTags *map[string]*string `json:"loadBalancerTags"`
}

// Properties for defining an Application Load Balancer.
// Experimental.
type ApplicationLoadBalancerProps struct {
	// The VPC network to place the load balancer in.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Indicates whether deletion protection is enabled.
	// Experimental.
	DeletionProtection *bool `json:"deletionProtection"`
	// Whether the load balancer has an internet-routable address.
	// Experimental.
	InternetFacing *bool `json:"internetFacing"`
	// Name of the load balancer.
	// Experimental.
	LoadBalancerName *string `json:"loadBalancerName"`
	// Which subnets place the load balancer in.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// Indicates whether HTTP/2 is enabled.
	// Experimental.
	Http2Enabled *bool `json:"http2Enabled"`
	// The load balancer idle timeout, in seconds.
	// Experimental.
	IdleTimeout awscdk.Duration `json:"idleTimeout"`
	// The type of IP addresses to use.
	//
	// Only applies to application load balancers.
	// Experimental.
	IpAddressType IpAddressType `json:"ipAddressType"`
	// Security group to associate with this load balancer.
	// Experimental.
	SecurityGroup awsec2.ISecurityGroup `json:"securityGroup"`
}

// Properties for a redirection config.
// Experimental.
type ApplicationLoadBalancerRedirectConfig struct {
	// Allow anyone to connect to this listener.
	//
	// If this is specified, the listener will be opened up to anyone who can reach it.
	// For internal load balancers this is anyone in the same VPC. For public load
	// balancers, this is anyone on the internet.
	//
	// If you want to be more selective about who can access this load
	// balancer, set this to `false` and use the listener's `connections`
	// object to selectively grant access to the listener.
	// Experimental.
	Open *bool `json:"open"`
	// The port number to listen to.
	// Experimental.
	SourcePort *float64 `json:"sourcePort"`
	// The protocol of the listener being created.
	// Experimental.
	SourceProtocol ApplicationProtocol `json:"sourceProtocol"`
	// The port number to redirect to.
	// Experimental.
	TargetPort *float64 `json:"targetPort"`
	// The protocol of the redirection target.
	// Experimental.
	TargetProtocol ApplicationProtocol `json:"targetProtocol"`
}

// Load balancing protocol for application load balancers.
// Experimental.
type ApplicationProtocol string

const (
	ApplicationProtocol_HTTP ApplicationProtocol = "HTTP"
	ApplicationProtocol_HTTPS ApplicationProtocol = "HTTPS"
)

// Load balancing protocol version for application load balancers.
// Experimental.
type ApplicationProtocolVersion string

const (
	ApplicationProtocolVersion_GRPC ApplicationProtocolVersion = "GRPC"
	ApplicationProtocolVersion_HTTP1 ApplicationProtocolVersion = "HTTP1"
	ApplicationProtocolVersion_HTTP2 ApplicationProtocolVersion = "HTTP2"
)

// Define an Application Target Group.
// Experimental.
type ApplicationTargetGroup interface {
	TargetGroupBase
	IApplicationTargetGroup
	DefaultPort() *float64
	FirstLoadBalancerFullName() *string
	HealthCheck() *HealthCheck
	SetHealthCheck(val *HealthCheck)
	LoadBalancerArns() *string
	LoadBalancerAttached() awscdk.IDependable
	LoadBalancerAttachedDependencies() awscdk.ConcreteDependable
	Node() awscdk.ConstructNode
	TargetGroupArn() *string
	TargetGroupFullName() *string
	TargetGroupLoadBalancerArns() *[]*string
	TargetGroupName() *string
	TargetType() TargetType
	SetTargetType(val TargetType)
	AddLoadBalancerTarget(props *LoadBalancerTargetProps)
	AddTarget(targets ...IApplicationLoadBalancerTarget)
	ConfigureHealthCheck(healthCheck *HealthCheck)
	EnableCookieStickiness(duration awscdk.Duration, cookieName *string)
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHttpCodeTarget(code HttpCodeTarget, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricIpv6RequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricRequestCountPerTarget(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTargetConnectionErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTargetResponseTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTargetTLSNegotiationErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricUnhealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port)
	RegisterListener(listener IApplicationListener, associatingConstruct constructs.IConstruct)
	SetAttribute(key *string, value *string)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for ApplicationTargetGroup
type jsiiProxy_ApplicationTargetGroup struct {
	jsiiProxy_TargetGroupBase
	jsiiProxy_IApplicationTargetGroup
}

func (j *jsiiProxy_ApplicationTargetGroup) DefaultPort() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"defaultPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) FirstLoadBalancerFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"firstLoadBalancerFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) HealthCheck() *HealthCheck {
	var returns *HealthCheck
	_jsii_.Get(
		j,
		"healthCheck",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) LoadBalancerArns() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) LoadBalancerAttached() awscdk.IDependable {
	var returns awscdk.IDependable
	_jsii_.Get(
		j,
		"loadBalancerAttached",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) LoadBalancerAttachedDependencies() awscdk.ConcreteDependable {
	var returns awscdk.ConcreteDependable
	_jsii_.Get(
		j,
		"loadBalancerAttachedDependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) TargetGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) TargetGroupFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) TargetGroupLoadBalancerArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"targetGroupLoadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) TargetGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ApplicationTargetGroup) TargetType() TargetType {
	var returns TargetType
	_jsii_.Get(
		j,
		"targetType",
		&returns,
	)
	return returns
}


// Experimental.
func NewApplicationTargetGroup(scope constructs.Construct, id *string, props *ApplicationTargetGroupProps) ApplicationTargetGroup {
	_init_.Initialize()

	j := jsiiProxy_ApplicationTargetGroup{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationTargetGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewApplicationTargetGroup_Override(a ApplicationTargetGroup, scope constructs.Construct, id *string, props *ApplicationTargetGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ApplicationTargetGroup",
		[]interface{}{scope, id, props},
		a,
	)
}

func (j *jsiiProxy_ApplicationTargetGroup) SetHealthCheck(val *HealthCheck) {
	_jsii_.Set(
		j,
		"healthCheck",
		val,
	)
}

func (j *jsiiProxy_ApplicationTargetGroup) SetTargetType(val TargetType) {
	_jsii_.Set(
		j,
		"targetType",
		val,
	)
}

// Import an existing target group.
// Experimental.
func ApplicationTargetGroup_FromTargetGroupAttributes(scope constructs.Construct, id *string, attrs *TargetGroupAttributes) IApplicationTargetGroup {
	_init_.Initialize()

	var returns IApplicationTargetGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationTargetGroup",
		"fromTargetGroupAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Import an existing target group.
// Deprecated: Use `fromTargetGroupAttributes` instead
func ApplicationTargetGroup_Import(scope constructs.Construct, id *string, props *TargetGroupImportProps) IApplicationTargetGroup {
	_init_.Initialize()

	var returns IApplicationTargetGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationTargetGroup",
		"import",
		[]interface{}{scope, id, props},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func ApplicationTargetGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ApplicationTargetGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Register the given load balancing target as part of this group.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) AddLoadBalancerTarget(props *LoadBalancerTargetProps) {
	_jsii_.InvokeVoid(
		a,
		"addLoadBalancerTarget",
		[]interface{}{props},
	)
}

// Add a load balancing target to this target group.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) AddTarget(targets ...IApplicationLoadBalancerTarget) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		a,
		"addTarget",
		args,
	)
}

// Set/replace the target group's health check.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) ConfigureHealthCheck(healthCheck *HealthCheck) {
	_jsii_.InvokeVoid(
		a,
		"configureHealthCheck",
		[]interface{}{healthCheck},
	)
}

// Enable sticky routing via a cookie to members of this target group.
//
// Note: If the `cookieName` parameter is set, application-based stickiness will be applied,
// otherwise it defaults to duration-based stickiness attributes (`lb_cookie`).
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/sticky-sessions.html
//
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) EnableCookieStickiness(duration awscdk.Duration, cookieName *string) {
	_jsii_.InvokeVoid(
		a,
		"enableCookieStickiness",
		[]interface{}{duration, cookieName},
	)
}

// Return the given named metric for this Application Load Balancer Target Group.
//
// Returns the metric for this target group from the point of view of the first
// load balancer load balancing to it. If you have multiple load balancers load
// sending traffic to the same target group, you will have to override the dimensions
// on this metric.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The number of healthy hosts in the target group.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHealthyHostCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of HTTP 2xx/3xx/4xx/5xx response codes generated by all targets in this target group.
//
// This does not include any response codes generated by the load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricHttpCodeTarget(code HttpCodeTarget, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricHttpCodeTarget",
		[]interface{}{code, props},
		&returns,
	)

	return returns
}

// The number of IPv6 requests received by the target group.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricIpv6RequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricIpv6RequestCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of requests processed over IPv4 and IPv6.
//
// This count includes only the requests with a response generated by a target of the load balancer.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricRequestCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricRequestCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The average number of requests received by each target in a target group.
//
// The only valid statistic is Sum. Note that this represents the average not the sum.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricRequestCountPerTarget(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricRequestCountPerTarget",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of connections that were not successfully established between the load balancer and target.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricTargetConnectionErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTargetConnectionErrorCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The time elapsed, in seconds, after the request leaves the load balancer until a response from the target is received.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricTargetResponseTime(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTargetResponseTime",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of TLS connections initiated by the load balancer that did not establish a session with the target.
//
// Possible causes include a mismatch of ciphers or protocols.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricTargetTLSNegotiationErrorCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricTargetTLSNegotiationErrorCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of unhealthy hosts in the target group.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) MetricUnhealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		a,
		"metricUnhealthyHostCount",
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
func (a *jsiiProxy_ApplicationTargetGroup) OnPrepare() {
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
func (a *jsiiProxy_ApplicationTargetGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (a *jsiiProxy_ApplicationTargetGroup) OnValidate() *[]*string {
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
func (a *jsiiProxy_ApplicationTargetGroup) Prepare() {
	_jsii_.InvokeVoid(
		a,
		"prepare",
		nil, // no parameters
	)
}

// Register a connectable as a member of this target group.
//
// Don't call this directly. It will be called by load balancing targets.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port) {
	_jsii_.InvokeVoid(
		a,
		"registerConnectable",
		[]interface{}{connectable, portRange},
	)
}

// Register a listener that is load balancing to this target group.
//
// Don't call this directly. It will be called by listeners.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) RegisterListener(listener IApplicationListener, associatingConstruct constructs.IConstruct) {
	_jsii_.InvokeVoid(
		a,
		"registerListener",
		[]interface{}{listener, associatingConstruct},
	)
}

// Set a non-standard attribute on the target group.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-target-groups.html#target-group-attributes
//
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) SetAttribute(key *string, value *string) {
	_jsii_.InvokeVoid(
		a,
		"setAttribute",
		[]interface{}{key, value},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		a,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) ToString() *string {
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
// Experimental.
func (a *jsiiProxy_ApplicationTargetGroup) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		a,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for defining an Application Target Group.
// Experimental.
type ApplicationTargetGroupProps struct {
	// The amount of time for Elastic Load Balancing to wait before deregistering a target.
	//
	// The range is 0-3600 seconds.
	// Experimental.
	DeregistrationDelay awscdk.Duration `json:"deregistrationDelay"`
	// Health check configuration.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The name of the target group.
	//
	// This name must be unique per region per account, can have a maximum of
	// 32 characters, must contain only alphanumeric characters or hyphens, and
	// must not begin or end with a hyphen.
	// Experimental.
	TargetGroupName *string `json:"targetGroupName"`
	// The type of targets registered to this TargetGroup, either IP or Instance.
	//
	// All targets registered into the group must be of this type. If you
	// register targets to the TargetGroup in the CDK app, the TargetType is
	// determined automatically.
	// Experimental.
	TargetType TargetType `json:"targetType"`
	// The virtual private cloud (VPC).
	//
	// only if `TargetType` is `Ip` or `InstanceId`
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// The protocol to use.
	// Experimental.
	Protocol ApplicationProtocol `json:"protocol"`
	// The protocol version to use.
	// Experimental.
	ProtocolVersion ApplicationProtocolVersion `json:"protocolVersion"`
	// The time period during which the load balancer sends a newly registered target a linearly increasing share of the traffic to the target group.
	//
	// The range is 30-900 seconds (15 minutes).
	// Experimental.
	SlowStart awscdk.Duration `json:"slowStart"`
	// The stickiness cookie expiration period.
	//
	// Setting this value enables load balancer stickiness.
	//
	// After this period, the cookie is considered stale. The minimum value is
	// 1 second and the maximum value is 7 days (604800 seconds).
	// Experimental.
	StickinessCookieDuration awscdk.Duration `json:"stickinessCookieDuration"`
	// The name of an application-based stickiness cookie.
	//
	// Names that start with the following prefixes are not allowed: AWSALB, AWSALBAPP,
	// and AWSALBTG; they're reserved for use by the load balancer.
	//
	// Note: `stickinessCookieName` parameter depends on the presence of `stickinessCookieDuration` parameter.
	// If `stickinessCookieDuration` is not set, `stickinessCookieName` will be omitted.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/sticky-sessions.html
	//
	// Experimental.
	StickinessCookieName *string `json:"stickinessCookieName"`
	// The targets to add to this target group.
	//
	// Can be `Instance`, `IPAddress`, or any self-registering load balancing
	// target. If you use either `Instance` or `IPAddress` as targets, all
	// target must be of the same type.
	// Experimental.
	Targets *[]IApplicationLoadBalancerTarget `json:"targets"`
}

// Options for `ListenerAction.authenciateOidc()`.
// Experimental.
type AuthenticateOidcOptions struct {
	// The authorization endpoint of the IdP.
	//
	// This must be a full URL, including the HTTPS protocol, the domain, and the path.
	// Experimental.
	AuthorizationEndpoint *string `json:"authorizationEndpoint"`
	// The OAuth 2.0 client identifier.
	// Experimental.
	ClientId *string `json:"clientId"`
	// The OAuth 2.0 client secret.
	// Experimental.
	ClientSecret awscdk.SecretValue `json:"clientSecret"`
	// The OIDC issuer identifier of the IdP.
	//
	// This must be a full URL, including the HTTPS protocol, the domain, and the path.
	// Experimental.
	Issuer *string `json:"issuer"`
	// What action to execute next.
	// Experimental.
	Next ListenerAction `json:"next"`
	// The token endpoint of the IdP.
	//
	// This must be a full URL, including the HTTPS protocol, the domain, and the path.
	// Experimental.
	TokenEndpoint *string `json:"tokenEndpoint"`
	// The user info endpoint of the IdP.
	//
	// This must be a full URL, including the HTTPS protocol, the domain, and the path.
	// Experimental.
	UserInfoEndpoint *string `json:"userInfoEndpoint"`
	// The query parameters (up to 10) to include in the redirect request to the authorization endpoint.
	// Experimental.
	AuthenticationRequestExtraParams *map[string]*string `json:"authenticationRequestExtraParams"`
	// The behavior if the user is not authenticated.
	// Experimental.
	OnUnauthenticatedRequest UnauthenticatedAction `json:"onUnauthenticatedRequest"`
	// The set of user claims to be requested from the IdP.
	//
	// To verify which scope values your IdP supports and how to separate multiple values, see the documentation for your IdP.
	// Experimental.
	Scope *string `json:"scope"`
	// The name of the cookie used to maintain session information.
	// Experimental.
	SessionCookieName *string `json:"sessionCookieName"`
	// The maximum duration of the authentication session.
	// Experimental.
	SessionTimeout awscdk.Duration `json:"sessionTimeout"`
}

// Basic properties for an ApplicationListener.
// Experimental.
type BaseApplicationListenerProps struct {
	// The certificates to use on this listener.
	// Deprecated: Use the `certificates` property instead
	CertificateArns *[]*string `json:"certificateArns"`
	// Certificate list of ACM cert ARNs.
	// Experimental.
	Certificates *[]IListenerCertificate `json:"certificates"`
	// Default action to take for requests to this listener.
	//
	// This allows full control of the default action of the load balancer,
	// including Action chaining, fixed responses and redirect responses.
	//
	// See the `ListenerAction` class for all options.
	//
	// Cannot be specified together with `defaultTargetGroups`.
	// Experimental.
	DefaultAction ListenerAction `json:"defaultAction"`
	// Default target groups to load balance to.
	//
	// All target groups will be load balanced to with equal weight and without
	// stickiness. For a more complex configuration than that, use
	// either `defaultAction` or `addAction()`.
	//
	// Cannot be specified together with `defaultAction`.
	// Experimental.
	DefaultTargetGroups *[]IApplicationTargetGroup `json:"defaultTargetGroups"`
	// Allow anyone to connect to this listener.
	//
	// If this is specified, the listener will be opened up to anyone who can reach it.
	// For internal load balancers this is anyone in the same VPC. For public load
	// balancers, this is anyone on the internet.
	//
	// If you want to be more selective about who can access this load
	// balancer, set this to `false` and use the listener's `connections`
	// object to selectively grant access to the listener.
	// Experimental.
	Open *bool `json:"open"`
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// The protocol to use.
	// Experimental.
	Protocol ApplicationProtocol `json:"protocol"`
	// The security policy that defines which ciphers and protocols are supported.
	// Experimental.
	SslPolicy SslPolicy `json:"sslPolicy"`
}

// Basic properties for defining a rule on a listener.
// Experimental.
type BaseApplicationListenerRuleProps struct {
	// Priority of the rule.
	//
	// The rule with the lowest priority will be used for every request.
	//
	// Priorities must be unique.
	// Experimental.
	Priority *float64 `json:"priority"`
	// Action to perform when requests are received.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	// Experimental.
	Action ListenerAction `json:"action"`
	// Rule applies if matches the conditions.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html
	//
	// Experimental.
	Conditions *[]ListenerCondition `json:"conditions"`
	// Fixed response to return.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	// Deprecated: Use `action` instead.
	FixedResponse *FixedResponse `json:"fixedResponse"`
	// Rule applies if the requested host matches the indicated host.
	//
	// May contain up to three '*' wildcards.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#host-conditions
	//
	// Deprecated: Use `conditions` instead.
	HostHeader *string `json:"hostHeader"`
	// Rule applies if the requested path matches the given path pattern.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPattern *string `json:"pathPattern"`
	// Rule applies if the requested path matches any of the given patterns.
	//
	// Paths may contain up to three '*' wildcards.
	// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#path-conditions
	//
	// Deprecated: Use `conditions` instead.
	PathPatterns *[]*string `json:"pathPatterns"`
	// Redirect response to return.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	// Deprecated: Use `action` instead.
	RedirectResponse *RedirectResponse `json:"redirectResponse"`
	// Target groups to forward requests to.
	//
	// Only one of `action`, `fixedResponse`, `redirectResponse` or `targetGroups` can be specified.
	//
	// Implies a `forward` action.
	// Experimental.
	TargetGroups *[]IApplicationTargetGroup `json:"targetGroups"`
}

// Base class for listeners.
// Experimental.
type BaseListener interface {
	awscdk.Resource
	Env() *awscdk.ResourceEnvironment
	ListenerArn() *string
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

// The jsii proxy struct for BaseListener
type jsiiProxy_BaseListener struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_BaseListener) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseListener) ListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseListener) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseListener) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseListener) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewBaseListener_Override(b BaseListener, scope constructs.Construct, id *string, additionalProps interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.BaseListener",
		[]interface{}{scope, id, additionalProps},
		b,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BaseListener_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.BaseListener",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func BaseListener_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.BaseListener",
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
func (b *jsiiProxy_BaseListener) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (b *jsiiProxy_BaseListener) GeneratePhysicalName() *string {
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
func (b *jsiiProxy_BaseListener) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (b *jsiiProxy_BaseListener) GetResourceNameAttribute(nameAttr *string) *string {
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
func (b *jsiiProxy_BaseListener) OnPrepare() {
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
func (b *jsiiProxy_BaseListener) OnSynthesize(session constructs.ISynthesisSession) {
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
func (b *jsiiProxy_BaseListener) OnValidate() *[]*string {
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
func (b *jsiiProxy_BaseListener) Prepare() {
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
func (b *jsiiProxy_BaseListener) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BaseListener) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate this listener.
// Experimental.
func (b *jsiiProxy_BaseListener) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options for listener lookup.
// Experimental.
type BaseListenerLookupOptions struct {
	// Filter listeners by listener port.
	// Experimental.
	ListenerPort *float64 `json:"listenerPort"`
	// Filter listeners by associated load balancer arn.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// Filter listeners by associated load balancer tags.
	// Experimental.
	LoadBalancerTags *map[string]*string `json:"loadBalancerTags"`
}

// Base class for both Application and Network Load Balancers.
// Experimental.
type BaseLoadBalancer interface {
	awscdk.Resource
	Env() *awscdk.ResourceEnvironment
	LoadBalancerArn() *string
	LoadBalancerCanonicalHostedZoneId() *string
	LoadBalancerDnsName() *string
	LoadBalancerFullName() *string
	LoadBalancerName() *string
	LoadBalancerSecurityGroups() *[]*string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	Vpc() awsec2.IVpc
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LogAccessLogs(bucket awss3.IBucket, prefix *string)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RemoveAttribute(key *string)
	SetAttribute(key *string, value *string)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for BaseLoadBalancer
type jsiiProxy_BaseLoadBalancer struct {
	internal.Type__awscdkResource
}

func (j *jsiiProxy_BaseLoadBalancer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) LoadBalancerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) LoadBalancerCanonicalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) LoadBalancerDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) LoadBalancerFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) LoadBalancerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) LoadBalancerSecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"loadBalancerSecurityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_BaseLoadBalancer) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}


// Experimental.
func NewBaseLoadBalancer_Override(b BaseLoadBalancer, scope constructs.Construct, id *string, baseProps *BaseLoadBalancerProps, additionalProps interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.BaseLoadBalancer",
		[]interface{}{scope, id, baseProps, additionalProps},
		b,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func BaseLoadBalancer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.BaseLoadBalancer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func BaseLoadBalancer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.BaseLoadBalancer",
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
func (b *jsiiProxy_BaseLoadBalancer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		b,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (b *jsiiProxy_BaseLoadBalancer) GeneratePhysicalName() *string {
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
func (b *jsiiProxy_BaseLoadBalancer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
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
func (b *jsiiProxy_BaseLoadBalancer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		b,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Enable access logging for this load balancer.
//
// A region must be specified on the stack containing the load balancer; you cannot enable logging on
// environment-agnostic stacks. See https://docs.aws.amazon.com/cdk/latest/guide/environments.html
// Experimental.
func (b *jsiiProxy_BaseLoadBalancer) LogAccessLogs(bucket awss3.IBucket, prefix *string) {
	_jsii_.InvokeVoid(
		b,
		"logAccessLogs",
		[]interface{}{bucket, prefix},
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
func (b *jsiiProxy_BaseLoadBalancer) OnPrepare() {
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
func (b *jsiiProxy_BaseLoadBalancer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (b *jsiiProxy_BaseLoadBalancer) OnValidate() *[]*string {
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
func (b *jsiiProxy_BaseLoadBalancer) Prepare() {
	_jsii_.InvokeVoid(
		b,
		"prepare",
		nil, // no parameters
	)
}

// Remove an attribute from the load balancer.
// Experimental.
func (b *jsiiProxy_BaseLoadBalancer) RemoveAttribute(key *string) {
	_jsii_.InvokeVoid(
		b,
		"removeAttribute",
		[]interface{}{key},
	)
}

// Set a non-standard attribute on the load balancer.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/application-load-balancers.html#load-balancer-attributes
//
// Experimental.
func (b *jsiiProxy_BaseLoadBalancer) SetAttribute(key *string, value *string) {
	_jsii_.InvokeVoid(
		b,
		"setAttribute",
		[]interface{}{key, value},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (b *jsiiProxy_BaseLoadBalancer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		b,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (b *jsiiProxy_BaseLoadBalancer) ToString() *string {
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
func (b *jsiiProxy_BaseLoadBalancer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		b,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Options for looking up load balancers.
// Experimental.
type BaseLoadBalancerLookupOptions struct {
	// Find by load balancer's ARN.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// Match load balancer tags.
	// Experimental.
	LoadBalancerTags *map[string]*string `json:"loadBalancerTags"`
}

// Shared properties of both Application and Network Load Balancers.
// Experimental.
type BaseLoadBalancerProps struct {
	// The VPC network to place the load balancer in.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Indicates whether deletion protection is enabled.
	// Experimental.
	DeletionProtection *bool `json:"deletionProtection"`
	// Whether the load balancer has an internet-routable address.
	// Experimental.
	InternetFacing *bool `json:"internetFacing"`
	// Name of the load balancer.
	// Experimental.
	LoadBalancerName *string `json:"loadBalancerName"`
	// Which subnets place the load balancer in.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
}

// Basic properties for a Network Listener.
// Experimental.
type BaseNetworkListenerProps struct {
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// Certificate list of ACM cert ARNs.
	// Experimental.
	Certificates *[]IListenerCertificate `json:"certificates"`
	// Default action to take for requests to this listener.
	//
	// This allows full control of the default Action of the load balancer,
	// including weighted forwarding. See the `NetworkListenerAction` class for
	// all options.
	//
	// Cannot be specified together with `defaultTargetGroups`.
	// Experimental.
	DefaultAction NetworkListenerAction `json:"defaultAction"`
	// Default target groups to load balance to.
	//
	// All target groups will be load balanced to with equal weight and without
	// stickiness. For a more complex configuration than that, use
	// either `defaultAction` or `addAction()`.
	//
	// Cannot be specified together with `defaultAction`.
	// Experimental.
	DefaultTargetGroups *[]INetworkTargetGroup `json:"defaultTargetGroups"`
	// Protocol for listener, expects TCP, TLS, UDP, or TCP_UDP.
	// Experimental.
	Protocol Protocol `json:"protocol"`
	// SSL Policy.
	// Experimental.
	SslPolicy SslPolicy `json:"sslPolicy"`
}

// Basic properties of both Application and Network Target Groups.
// Experimental.
type BaseTargetGroupProps struct {
	// The amount of time for Elastic Load Balancing to wait before deregistering a target.
	//
	// The range is 0-3600 seconds.
	// Experimental.
	DeregistrationDelay awscdk.Duration `json:"deregistrationDelay"`
	// Health check configuration.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The name of the target group.
	//
	// This name must be unique per region per account, can have a maximum of
	// 32 characters, must contain only alphanumeric characters or hyphens, and
	// must not begin or end with a hyphen.
	// Experimental.
	TargetGroupName *string `json:"targetGroupName"`
	// The type of targets registered to this TargetGroup, either IP or Instance.
	//
	// All targets registered into the group must be of this type. If you
	// register targets to the TargetGroup in the CDK app, the TargetType is
	// determined automatically.
	// Experimental.
	TargetType TargetType `json:"targetType"`
	// The virtual private cloud (VPC).
	//
	// only if `TargetType` is `Ip` or `InstanceId`
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
}

// A CloudFormation `AWS::ElasticLoadBalancingV2::Listener`.
type CfnListener interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AlpnPolicy() *[]*string
	SetAlpnPolicy(val *[]*string)
	AttrListenerArn() *string
	Certificates() interface{}
	SetCertificates(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	DefaultActions() interface{}
	SetDefaultActions(val interface{})
	LoadBalancerArn() *string
	SetLoadBalancerArn(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Port() *float64
	SetPort(val *float64)
	Protocol() *string
	SetProtocol(val *string)
	Ref() *string
	SslPolicy() *string
	SetSslPolicy(val *string)
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

// The jsii proxy struct for CfnListener
type jsiiProxy_CfnListener struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnListener) AlpnPolicy() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"alpnPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) AttrListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrListenerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) Certificates() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"certificates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) DefaultActions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"defaultActions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) LoadBalancerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) Port() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"port",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) Protocol() *string {
	var returns *string
	_jsii_.Get(
		j,
		"protocol",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) SslPolicy() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sslPolicy",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListener) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ElasticLoadBalancingV2::Listener`.
func NewCfnListener(scope awscdk.Construct, id *string, props *CfnListenerProps) CfnListener {
	_init_.Initialize()

	j := jsiiProxy_CfnListener{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnListener",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ElasticLoadBalancingV2::Listener`.
func NewCfnListener_Override(c CfnListener, scope awscdk.Construct, id *string, props *CfnListenerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnListener",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnListener) SetAlpnPolicy(val *[]*string) {
	_jsii_.Set(
		j,
		"alpnPolicy",
		val,
	)
}

func (j *jsiiProxy_CfnListener) SetCertificates(val interface{}) {
	_jsii_.Set(
		j,
		"certificates",
		val,
	)
}

func (j *jsiiProxy_CfnListener) SetDefaultActions(val interface{}) {
	_jsii_.Set(
		j,
		"defaultActions",
		val,
	)
}

func (j *jsiiProxy_CfnListener) SetLoadBalancerArn(val *string) {
	_jsii_.Set(
		j,
		"loadBalancerArn",
		val,
	)
}

func (j *jsiiProxy_CfnListener) SetPort(val *float64) {
	_jsii_.Set(
		j,
		"port",
		val,
	)
}

func (j *jsiiProxy_CfnListener) SetProtocol(val *string) {
	_jsii_.Set(
		j,
		"protocol",
		val,
	)
}

func (j *jsiiProxy_CfnListener) SetSslPolicy(val *string) {
	_jsii_.Set(
		j,
		"sslPolicy",
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
func CfnListener_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListener",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnListener_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListener",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnListener_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListener",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnListener_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_elasticloadbalancingv2.CfnListener",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnListener) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnListener) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnListener) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnListener) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnListener) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnListener) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnListener) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnListener) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnListener) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnListener) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnListener) OnPrepare() {
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
func (c *jsiiProxy_CfnListener) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnListener) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnListener) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnListener) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnListener) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnListener) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnListener) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnListener) ToString() *string {
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
func (c *jsiiProxy_CfnListener) Validate() *[]*string {
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
func (c *jsiiProxy_CfnListener) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnListener_ActionProperty struct {
	// `CfnListener.ActionProperty.Type`.
	Type *string `json:"type"`
	// `CfnListener.ActionProperty.AuthenticateCognitoConfig`.
	AuthenticateCognitoConfig interface{} `json:"authenticateCognitoConfig"`
	// `CfnListener.ActionProperty.AuthenticateOidcConfig`.
	AuthenticateOidcConfig interface{} `json:"authenticateOidcConfig"`
	// `CfnListener.ActionProperty.FixedResponseConfig`.
	FixedResponseConfig interface{} `json:"fixedResponseConfig"`
	// `CfnListener.ActionProperty.ForwardConfig`.
	ForwardConfig interface{} `json:"forwardConfig"`
	// `CfnListener.ActionProperty.Order`.
	Order *float64 `json:"order"`
	// `CfnListener.ActionProperty.RedirectConfig`.
	RedirectConfig interface{} `json:"redirectConfig"`
	// `CfnListener.ActionProperty.TargetGroupArn`.
	TargetGroupArn *string `json:"targetGroupArn"`
}

type CfnListener_AuthenticateCognitoConfigProperty struct {
	// `CfnListener.AuthenticateCognitoConfigProperty.UserPoolArn`.
	UserPoolArn *string `json:"userPoolArn"`
	// `CfnListener.AuthenticateCognitoConfigProperty.UserPoolClientId`.
	UserPoolClientId *string `json:"userPoolClientId"`
	// `CfnListener.AuthenticateCognitoConfigProperty.UserPoolDomain`.
	UserPoolDomain *string `json:"userPoolDomain"`
	// `CfnListener.AuthenticateCognitoConfigProperty.AuthenticationRequestExtraParams`.
	AuthenticationRequestExtraParams interface{} `json:"authenticationRequestExtraParams"`
	// `CfnListener.AuthenticateCognitoConfigProperty.OnUnauthenticatedRequest`.
	OnUnauthenticatedRequest *string `json:"onUnauthenticatedRequest"`
	// `CfnListener.AuthenticateCognitoConfigProperty.Scope`.
	Scope *string `json:"scope"`
	// `CfnListener.AuthenticateCognitoConfigProperty.SessionCookieName`.
	SessionCookieName *string `json:"sessionCookieName"`
	// `CfnListener.AuthenticateCognitoConfigProperty.SessionTimeout`.
	SessionTimeout *string `json:"sessionTimeout"`
}

type CfnListener_AuthenticateOidcConfigProperty struct {
	// `CfnListener.AuthenticateOidcConfigProperty.AuthorizationEndpoint`.
	AuthorizationEndpoint *string `json:"authorizationEndpoint"`
	// `CfnListener.AuthenticateOidcConfigProperty.ClientId`.
	ClientId *string `json:"clientId"`
	// `CfnListener.AuthenticateOidcConfigProperty.ClientSecret`.
	ClientSecret *string `json:"clientSecret"`
	// `CfnListener.AuthenticateOidcConfigProperty.Issuer`.
	Issuer *string `json:"issuer"`
	// `CfnListener.AuthenticateOidcConfigProperty.TokenEndpoint`.
	TokenEndpoint *string `json:"tokenEndpoint"`
	// `CfnListener.AuthenticateOidcConfigProperty.UserInfoEndpoint`.
	UserInfoEndpoint *string `json:"userInfoEndpoint"`
	// `CfnListener.AuthenticateOidcConfigProperty.AuthenticationRequestExtraParams`.
	AuthenticationRequestExtraParams interface{} `json:"authenticationRequestExtraParams"`
	// `CfnListener.AuthenticateOidcConfigProperty.OnUnauthenticatedRequest`.
	OnUnauthenticatedRequest *string `json:"onUnauthenticatedRequest"`
	// `CfnListener.AuthenticateOidcConfigProperty.Scope`.
	Scope *string `json:"scope"`
	// `CfnListener.AuthenticateOidcConfigProperty.SessionCookieName`.
	SessionCookieName *string `json:"sessionCookieName"`
	// `CfnListener.AuthenticateOidcConfigProperty.SessionTimeout`.
	SessionTimeout *string `json:"sessionTimeout"`
}

type CfnListener_CertificateProperty struct {
	// `CfnListener.CertificateProperty.CertificateArn`.
	CertificateArn *string `json:"certificateArn"`
}

type CfnListener_FixedResponseConfigProperty struct {
	// `CfnListener.FixedResponseConfigProperty.StatusCode`.
	StatusCode *string `json:"statusCode"`
	// `CfnListener.FixedResponseConfigProperty.ContentType`.
	ContentType *string `json:"contentType"`
	// `CfnListener.FixedResponseConfigProperty.MessageBody`.
	MessageBody *string `json:"messageBody"`
}

type CfnListener_ForwardConfigProperty struct {
	// `CfnListener.ForwardConfigProperty.TargetGroups`.
	TargetGroups interface{} `json:"targetGroups"`
	// `CfnListener.ForwardConfigProperty.TargetGroupStickinessConfig`.
	TargetGroupStickinessConfig interface{} `json:"targetGroupStickinessConfig"`
}

type CfnListener_RedirectConfigProperty struct {
	// `CfnListener.RedirectConfigProperty.StatusCode`.
	StatusCode *string `json:"statusCode"`
	// `CfnListener.RedirectConfigProperty.Host`.
	Host *string `json:"host"`
	// `CfnListener.RedirectConfigProperty.Path`.
	Path *string `json:"path"`
	// `CfnListener.RedirectConfigProperty.Port`.
	Port *string `json:"port"`
	// `CfnListener.RedirectConfigProperty.Protocol`.
	Protocol *string `json:"protocol"`
	// `CfnListener.RedirectConfigProperty.Query`.
	Query *string `json:"query"`
}

type CfnListener_TargetGroupStickinessConfigProperty struct {
	// `CfnListener.TargetGroupStickinessConfigProperty.DurationSeconds`.
	DurationSeconds *float64 `json:"durationSeconds"`
	// `CfnListener.TargetGroupStickinessConfigProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
}

type CfnListener_TargetGroupTupleProperty struct {
	// `CfnListener.TargetGroupTupleProperty.TargetGroupArn`.
	TargetGroupArn *string `json:"targetGroupArn"`
	// `CfnListener.TargetGroupTupleProperty.Weight`.
	Weight *float64 `json:"weight"`
}

// A CloudFormation `AWS::ElasticLoadBalancingV2::ListenerCertificate`.
type CfnListenerCertificate interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Certificates() interface{}
	SetCertificates(val interface{})
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	ListenerArn() *string
	SetListenerArn(val *string)
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

// The jsii proxy struct for CfnListenerCertificate
type jsiiProxy_CfnListenerCertificate struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnListenerCertificate) Certificates() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"certificates",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) ListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerCertificate) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ElasticLoadBalancingV2::ListenerCertificate`.
func NewCfnListenerCertificate(scope awscdk.Construct, id *string, props *CfnListenerCertificateProps) CfnListenerCertificate {
	_init_.Initialize()

	j := jsiiProxy_CfnListenerCertificate{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerCertificate",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ElasticLoadBalancingV2::ListenerCertificate`.
func NewCfnListenerCertificate_Override(c CfnListenerCertificate, scope awscdk.Construct, id *string, props *CfnListenerCertificateProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerCertificate",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnListenerCertificate) SetCertificates(val interface{}) {
	_jsii_.Set(
		j,
		"certificates",
		val,
	)
}

func (j *jsiiProxy_CfnListenerCertificate) SetListenerArn(val *string) {
	_jsii_.Set(
		j,
		"listenerArn",
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
func CfnListenerCertificate_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerCertificate",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnListenerCertificate_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerCertificate",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnListenerCertificate_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerCertificate",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnListenerCertificate_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerCertificate",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnListenerCertificate) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnListenerCertificate) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnListenerCertificate) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnListenerCertificate) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnListenerCertificate) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnListenerCertificate) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnListenerCertificate) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnListenerCertificate) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnListenerCertificate) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnListenerCertificate) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnListenerCertificate) OnPrepare() {
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
func (c *jsiiProxy_CfnListenerCertificate) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnListenerCertificate) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnListenerCertificate) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnListenerCertificate) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnListenerCertificate) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnListenerCertificate) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnListenerCertificate) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnListenerCertificate) ToString() *string {
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
func (c *jsiiProxy_CfnListenerCertificate) Validate() *[]*string {
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
func (c *jsiiProxy_CfnListenerCertificate) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnListenerCertificate_CertificateProperty struct {
	// `CfnListenerCertificate.CertificateProperty.CertificateArn`.
	CertificateArn *string `json:"certificateArn"`
}

// Properties for defining a `AWS::ElasticLoadBalancingV2::ListenerCertificate`.
type CfnListenerCertificateProps struct {
	// `AWS::ElasticLoadBalancingV2::ListenerCertificate.Certificates`.
	Certificates interface{} `json:"certificates"`
	// `AWS::ElasticLoadBalancingV2::ListenerCertificate.ListenerArn`.
	ListenerArn *string `json:"listenerArn"`
}

// Properties for defining a `AWS::ElasticLoadBalancingV2::Listener`.
type CfnListenerProps struct {
	// `AWS::ElasticLoadBalancingV2::Listener.DefaultActions`.
	DefaultActions interface{} `json:"defaultActions"`
	// `AWS::ElasticLoadBalancingV2::Listener.LoadBalancerArn`.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// `AWS::ElasticLoadBalancingV2::Listener.AlpnPolicy`.
	AlpnPolicy *[]*string `json:"alpnPolicy"`
	// `AWS::ElasticLoadBalancingV2::Listener.Certificates`.
	Certificates interface{} `json:"certificates"`
	// `AWS::ElasticLoadBalancingV2::Listener.Port`.
	Port *float64 `json:"port"`
	// `AWS::ElasticLoadBalancingV2::Listener.Protocol`.
	Protocol *string `json:"protocol"`
	// `AWS::ElasticLoadBalancingV2::Listener.SslPolicy`.
	SslPolicy *string `json:"sslPolicy"`
}

// A CloudFormation `AWS::ElasticLoadBalancingV2::ListenerRule`.
type CfnListenerRule interface {
	awscdk.CfnResource
	awscdk.IInspectable
	Actions() interface{}
	SetActions(val interface{})
	AttrIsDefault() awscdk.IResolvable
	AttrRuleArn() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	Conditions() interface{}
	SetConditions(val interface{})
	CreationStack() *[]*string
	ListenerArn() *string
	SetListenerArn(val *string)
	LogicalId() *string
	Node() awscdk.ConstructNode
	Priority() *float64
	SetPriority(val *float64)
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

// The jsii proxy struct for CfnListenerRule
type jsiiProxy_CfnListenerRule struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnListenerRule) Actions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"actions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) AttrIsDefault() awscdk.IResolvable {
	var returns awscdk.IResolvable
	_jsii_.Get(
		j,
		"attrIsDefault",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) AttrRuleArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrRuleArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) Conditions() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"conditions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) ListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) Priority() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"priority",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnListenerRule) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ElasticLoadBalancingV2::ListenerRule`.
func NewCfnListenerRule(scope awscdk.Construct, id *string, props *CfnListenerRuleProps) CfnListenerRule {
	_init_.Initialize()

	j := jsiiProxy_CfnListenerRule{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerRule",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ElasticLoadBalancingV2::ListenerRule`.
func NewCfnListenerRule_Override(c CfnListenerRule, scope awscdk.Construct, id *string, props *CfnListenerRuleProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerRule",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnListenerRule) SetActions(val interface{}) {
	_jsii_.Set(
		j,
		"actions",
		val,
	)
}

func (j *jsiiProxy_CfnListenerRule) SetConditions(val interface{}) {
	_jsii_.Set(
		j,
		"conditions",
		val,
	)
}

func (j *jsiiProxy_CfnListenerRule) SetListenerArn(val *string) {
	_jsii_.Set(
		j,
		"listenerArn",
		val,
	)
}

func (j *jsiiProxy_CfnListenerRule) SetPriority(val *float64) {
	_jsii_.Set(
		j,
		"priority",
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
func CfnListenerRule_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerRule",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnListenerRule_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerRule",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnListenerRule_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerRule",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnListenerRule_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_elasticloadbalancingv2.CfnListenerRule",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnListenerRule) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnListenerRule) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnListenerRule) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnListenerRule) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnListenerRule) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnListenerRule) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnListenerRule) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnListenerRule) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnListenerRule) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnListenerRule) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnListenerRule) OnPrepare() {
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
func (c *jsiiProxy_CfnListenerRule) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnListenerRule) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnListenerRule) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnListenerRule) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnListenerRule) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnListenerRule) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnListenerRule) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnListenerRule) ToString() *string {
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
func (c *jsiiProxy_CfnListenerRule) Validate() *[]*string {
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
func (c *jsiiProxy_CfnListenerRule) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnListenerRule_ActionProperty struct {
	// `CfnListenerRule.ActionProperty.Type`.
	Type *string `json:"type"`
	// `CfnListenerRule.ActionProperty.AuthenticateCognitoConfig`.
	AuthenticateCognitoConfig interface{} `json:"authenticateCognitoConfig"`
	// `CfnListenerRule.ActionProperty.AuthenticateOidcConfig`.
	AuthenticateOidcConfig interface{} `json:"authenticateOidcConfig"`
	// `CfnListenerRule.ActionProperty.FixedResponseConfig`.
	FixedResponseConfig interface{} `json:"fixedResponseConfig"`
	// `CfnListenerRule.ActionProperty.ForwardConfig`.
	ForwardConfig interface{} `json:"forwardConfig"`
	// `CfnListenerRule.ActionProperty.Order`.
	Order *float64 `json:"order"`
	// `CfnListenerRule.ActionProperty.RedirectConfig`.
	RedirectConfig interface{} `json:"redirectConfig"`
	// `CfnListenerRule.ActionProperty.TargetGroupArn`.
	TargetGroupArn *string `json:"targetGroupArn"`
}

type CfnListenerRule_AuthenticateCognitoConfigProperty struct {
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.UserPoolArn`.
	UserPoolArn *string `json:"userPoolArn"`
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.UserPoolClientId`.
	UserPoolClientId *string `json:"userPoolClientId"`
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.UserPoolDomain`.
	UserPoolDomain *string `json:"userPoolDomain"`
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.AuthenticationRequestExtraParams`.
	AuthenticationRequestExtraParams interface{} `json:"authenticationRequestExtraParams"`
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.OnUnauthenticatedRequest`.
	OnUnauthenticatedRequest *string `json:"onUnauthenticatedRequest"`
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.Scope`.
	Scope *string `json:"scope"`
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.SessionCookieName`.
	SessionCookieName *string `json:"sessionCookieName"`
	// `CfnListenerRule.AuthenticateCognitoConfigProperty.SessionTimeout`.
	SessionTimeout *float64 `json:"sessionTimeout"`
}

type CfnListenerRule_AuthenticateOidcConfigProperty struct {
	// `CfnListenerRule.AuthenticateOidcConfigProperty.AuthorizationEndpoint`.
	AuthorizationEndpoint *string `json:"authorizationEndpoint"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.ClientId`.
	ClientId *string `json:"clientId"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.ClientSecret`.
	ClientSecret *string `json:"clientSecret"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.Issuer`.
	Issuer *string `json:"issuer"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.TokenEndpoint`.
	TokenEndpoint *string `json:"tokenEndpoint"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.UserInfoEndpoint`.
	UserInfoEndpoint *string `json:"userInfoEndpoint"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.AuthenticationRequestExtraParams`.
	AuthenticationRequestExtraParams interface{} `json:"authenticationRequestExtraParams"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.OnUnauthenticatedRequest`.
	OnUnauthenticatedRequest *string `json:"onUnauthenticatedRequest"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.Scope`.
	Scope *string `json:"scope"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.SessionCookieName`.
	SessionCookieName *string `json:"sessionCookieName"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.SessionTimeout`.
	SessionTimeout *float64 `json:"sessionTimeout"`
	// `CfnListenerRule.AuthenticateOidcConfigProperty.UseExistingClientSecret`.
	UseExistingClientSecret interface{} `json:"useExistingClientSecret"`
}

type CfnListenerRule_FixedResponseConfigProperty struct {
	// `CfnListenerRule.FixedResponseConfigProperty.StatusCode`.
	StatusCode *string `json:"statusCode"`
	// `CfnListenerRule.FixedResponseConfigProperty.ContentType`.
	ContentType *string `json:"contentType"`
	// `CfnListenerRule.FixedResponseConfigProperty.MessageBody`.
	MessageBody *string `json:"messageBody"`
}

type CfnListenerRule_ForwardConfigProperty struct {
	// `CfnListenerRule.ForwardConfigProperty.TargetGroups`.
	TargetGroups interface{} `json:"targetGroups"`
	// `CfnListenerRule.ForwardConfigProperty.TargetGroupStickinessConfig`.
	TargetGroupStickinessConfig interface{} `json:"targetGroupStickinessConfig"`
}

type CfnListenerRule_HostHeaderConfigProperty struct {
	// `CfnListenerRule.HostHeaderConfigProperty.Values`.
	Values *[]*string `json:"values"`
}

type CfnListenerRule_HttpHeaderConfigProperty struct {
	// `CfnListenerRule.HttpHeaderConfigProperty.HttpHeaderName`.
	HttpHeaderName *string `json:"httpHeaderName"`
	// `CfnListenerRule.HttpHeaderConfigProperty.Values`.
	Values *[]*string `json:"values"`
}

type CfnListenerRule_HttpRequestMethodConfigProperty struct {
	// `CfnListenerRule.HttpRequestMethodConfigProperty.Values`.
	Values *[]*string `json:"values"`
}

type CfnListenerRule_PathPatternConfigProperty struct {
	// `CfnListenerRule.PathPatternConfigProperty.Values`.
	Values *[]*string `json:"values"`
}

type CfnListenerRule_QueryStringConfigProperty struct {
	// `CfnListenerRule.QueryStringConfigProperty.Values`.
	Values interface{} `json:"values"`
}

type CfnListenerRule_QueryStringKeyValueProperty struct {
	// `CfnListenerRule.QueryStringKeyValueProperty.Key`.
	Key *string `json:"key"`
	// `CfnListenerRule.QueryStringKeyValueProperty.Value`.
	Value *string `json:"value"`
}

type CfnListenerRule_RedirectConfigProperty struct {
	// `CfnListenerRule.RedirectConfigProperty.StatusCode`.
	StatusCode *string `json:"statusCode"`
	// `CfnListenerRule.RedirectConfigProperty.Host`.
	Host *string `json:"host"`
	// `CfnListenerRule.RedirectConfigProperty.Path`.
	Path *string `json:"path"`
	// `CfnListenerRule.RedirectConfigProperty.Port`.
	Port *string `json:"port"`
	// `CfnListenerRule.RedirectConfigProperty.Protocol`.
	Protocol *string `json:"protocol"`
	// `CfnListenerRule.RedirectConfigProperty.Query`.
	Query *string `json:"query"`
}

type CfnListenerRule_RuleConditionProperty struct {
	// `CfnListenerRule.RuleConditionProperty.Field`.
	Field *string `json:"field"`
	// `CfnListenerRule.RuleConditionProperty.HostHeaderConfig`.
	HostHeaderConfig interface{} `json:"hostHeaderConfig"`
	// `CfnListenerRule.RuleConditionProperty.HttpHeaderConfig`.
	HttpHeaderConfig interface{} `json:"httpHeaderConfig"`
	// `CfnListenerRule.RuleConditionProperty.HttpRequestMethodConfig`.
	HttpRequestMethodConfig interface{} `json:"httpRequestMethodConfig"`
	// `CfnListenerRule.RuleConditionProperty.PathPatternConfig`.
	PathPatternConfig interface{} `json:"pathPatternConfig"`
	// `CfnListenerRule.RuleConditionProperty.QueryStringConfig`.
	QueryStringConfig interface{} `json:"queryStringConfig"`
	// `CfnListenerRule.RuleConditionProperty.SourceIpConfig`.
	SourceIpConfig interface{} `json:"sourceIpConfig"`
	// `CfnListenerRule.RuleConditionProperty.Values`.
	Values *[]*string `json:"values"`
}

type CfnListenerRule_SourceIpConfigProperty struct {
	// `CfnListenerRule.SourceIpConfigProperty.Values`.
	Values *[]*string `json:"values"`
}

type CfnListenerRule_TargetGroupStickinessConfigProperty struct {
	// `CfnListenerRule.TargetGroupStickinessConfigProperty.DurationSeconds`.
	DurationSeconds *float64 `json:"durationSeconds"`
	// `CfnListenerRule.TargetGroupStickinessConfigProperty.Enabled`.
	Enabled interface{} `json:"enabled"`
}

type CfnListenerRule_TargetGroupTupleProperty struct {
	// `CfnListenerRule.TargetGroupTupleProperty.TargetGroupArn`.
	TargetGroupArn *string `json:"targetGroupArn"`
	// `CfnListenerRule.TargetGroupTupleProperty.Weight`.
	Weight *float64 `json:"weight"`
}

// Properties for defining a `AWS::ElasticLoadBalancingV2::ListenerRule`.
type CfnListenerRuleProps struct {
	// `AWS::ElasticLoadBalancingV2::ListenerRule.Actions`.
	Actions interface{} `json:"actions"`
	// `AWS::ElasticLoadBalancingV2::ListenerRule.Conditions`.
	Conditions interface{} `json:"conditions"`
	// `AWS::ElasticLoadBalancingV2::ListenerRule.ListenerArn`.
	ListenerArn *string `json:"listenerArn"`
	// `AWS::ElasticLoadBalancingV2::ListenerRule.Priority`.
	Priority *float64 `json:"priority"`
}

// A CloudFormation `AWS::ElasticLoadBalancingV2::LoadBalancer`.
type CfnLoadBalancer interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrCanonicalHostedZoneId() *string
	AttrDnsName() *string
	AttrLoadBalancerFullName() *string
	AttrLoadBalancerName() *string
	AttrSecurityGroups() *[]*string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	IpAddressType() *string
	SetIpAddressType(val *string)
	LoadBalancerAttributes() interface{}
	SetLoadBalancerAttributes(val interface{})
	LogicalId() *string
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Ref() *string
	Scheme() *string
	SetScheme(val *string)
	SecurityGroups() *[]*string
	SetSecurityGroups(val *[]*string)
	Stack() awscdk.Stack
	SubnetMappings() interface{}
	SetSubnetMappings(val interface{})
	Subnets() *[]*string
	SetSubnets(val *[]*string)
	Tags() awscdk.TagManager
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

// The jsii proxy struct for CfnLoadBalancer
type jsiiProxy_CfnLoadBalancer struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnLoadBalancer) AttrCanonicalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrCanonicalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrLoadBalancerFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrLoadBalancerFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrLoadBalancerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrLoadBalancerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) AttrSecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"attrSecurityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) IpAddressType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ipAddressType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) LoadBalancerAttributes() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"loadBalancerAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Scheme() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scheme",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) SecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) SubnetMappings() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"subnetMappings",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Subnets() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"subnets",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) Type() *string {
	var returns *string
	_jsii_.Get(
		j,
		"type",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnLoadBalancer) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}


// Create a new `AWS::ElasticLoadBalancingV2::LoadBalancer`.
func NewCfnLoadBalancer(scope awscdk.Construct, id *string, props *CfnLoadBalancerProps) CfnLoadBalancer {
	_init_.Initialize()

	j := jsiiProxy_CfnLoadBalancer{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnLoadBalancer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ElasticLoadBalancingV2::LoadBalancer`.
func NewCfnLoadBalancer_Override(c CfnLoadBalancer, scope awscdk.Construct, id *string, props *CfnLoadBalancerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnLoadBalancer",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetIpAddressType(val *string) {
	_jsii_.Set(
		j,
		"ipAddressType",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetLoadBalancerAttributes(val interface{}) {
	_jsii_.Set(
		j,
		"loadBalancerAttributes",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetScheme(val *string) {
	_jsii_.Set(
		j,
		"scheme",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetSecurityGroups(val *[]*string) {
	_jsii_.Set(
		j,
		"securityGroups",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetSubnetMappings(val interface{}) {
	_jsii_.Set(
		j,
		"subnetMappings",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetSubnets(val *[]*string) {
	_jsii_.Set(
		j,
		"subnets",
		val,
	)
}

func (j *jsiiProxy_CfnLoadBalancer) SetType(val *string) {
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
func CfnLoadBalancer_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnLoadBalancer",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnLoadBalancer_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnLoadBalancer",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnLoadBalancer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnLoadBalancer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnLoadBalancer_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_elasticloadbalancingv2.CfnLoadBalancer",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnLoadBalancer) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnLoadBalancer) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnLoadBalancer) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnLoadBalancer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnLoadBalancer) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnLoadBalancer) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnLoadBalancer) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnLoadBalancer) OnPrepare() {
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
func (c *jsiiProxy_CfnLoadBalancer) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLoadBalancer) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnLoadBalancer) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnLoadBalancer) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnLoadBalancer) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnLoadBalancer) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnLoadBalancer) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnLoadBalancer) ToString() *string {
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
func (c *jsiiProxy_CfnLoadBalancer) Validate() *[]*string {
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
func (c *jsiiProxy_CfnLoadBalancer) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnLoadBalancer_LoadBalancerAttributeProperty struct {
	// `CfnLoadBalancer.LoadBalancerAttributeProperty.Key`.
	Key *string `json:"key"`
	// `CfnLoadBalancer.LoadBalancerAttributeProperty.Value`.
	Value *string `json:"value"`
}

type CfnLoadBalancer_SubnetMappingProperty struct {
	// `CfnLoadBalancer.SubnetMappingProperty.SubnetId`.
	SubnetId *string `json:"subnetId"`
	// `CfnLoadBalancer.SubnetMappingProperty.AllocationId`.
	AllocationId *string `json:"allocationId"`
	// `CfnLoadBalancer.SubnetMappingProperty.IPv6Address`.
	IPv6Address *string `json:"iPv6Address"`
	// `CfnLoadBalancer.SubnetMappingProperty.PrivateIPv4Address`.
	PrivateIPv4Address *string `json:"privateIPv4Address"`
}

// Properties for defining a `AWS::ElasticLoadBalancingV2::LoadBalancer`.
type CfnLoadBalancerProps struct {
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.IpAddressType`.
	IpAddressType *string `json:"ipAddressType"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.LoadBalancerAttributes`.
	LoadBalancerAttributes interface{} `json:"loadBalancerAttributes"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.Name`.
	Name *string `json:"name"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.Scheme`.
	Scheme *string `json:"scheme"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.SecurityGroups`.
	SecurityGroups *[]*string `json:"securityGroups"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.SubnetMappings`.
	SubnetMappings interface{} `json:"subnetMappings"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.Subnets`.
	Subnets *[]*string `json:"subnets"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::ElasticLoadBalancingV2::LoadBalancer.Type`.
	Type *string `json:"type"`
}

// A CloudFormation `AWS::ElasticLoadBalancingV2::TargetGroup`.
type CfnTargetGroup interface {
	awscdk.CfnResource
	awscdk.IInspectable
	AttrLoadBalancerArns() *[]*string
	AttrTargetGroupFullName() *string
	AttrTargetGroupName() *string
	CfnOptions() awscdk.ICfnResourceOptions
	CfnProperties() *map[string]interface{}
	CfnResourceType() *string
	CreationStack() *[]*string
	HealthCheckEnabled() interface{}
	SetHealthCheckEnabled(val interface{})
	HealthCheckIntervalSeconds() *float64
	SetHealthCheckIntervalSeconds(val *float64)
	HealthCheckPath() *string
	SetHealthCheckPath(val *string)
	HealthCheckPort() *string
	SetHealthCheckPort(val *string)
	HealthCheckProtocol() *string
	SetHealthCheckProtocol(val *string)
	HealthCheckTimeoutSeconds() *float64
	SetHealthCheckTimeoutSeconds(val *float64)
	HealthyThresholdCount() *float64
	SetHealthyThresholdCount(val *float64)
	LogicalId() *string
	Matcher() interface{}
	SetMatcher(val interface{})
	Name() *string
	SetName(val *string)
	Node() awscdk.ConstructNode
	Port() *float64
	SetPort(val *float64)
	Protocol() *string
	SetProtocol(val *string)
	ProtocolVersion() *string
	SetProtocolVersion(val *string)
	Ref() *string
	Stack() awscdk.Stack
	Tags() awscdk.TagManager
	TargetGroupAttributes() interface{}
	SetTargetGroupAttributes(val interface{})
	Targets() interface{}
	SetTargets(val interface{})
	TargetType() *string
	SetTargetType(val *string)
	UnhealthyThresholdCount() *float64
	SetUnhealthyThresholdCount(val *float64)
	UpdatedProperites() *map[string]interface{}
	VpcId() *string
	SetVpcId(val *string)
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

// The jsii proxy struct for CfnTargetGroup
type jsiiProxy_CfnTargetGroup struct {
	internal.Type__awscdkCfnResource
	internal.Type__awscdkIInspectable
}

func (j *jsiiProxy_CfnTargetGroup) AttrLoadBalancerArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"attrLoadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) AttrTargetGroupFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrTargetGroupFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) AttrTargetGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"attrTargetGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) CfnOptions() awscdk.ICfnResourceOptions {
	var returns awscdk.ICfnResourceOptions
	_jsii_.Get(
		j,
		"cfnOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) CfnProperties() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"cfnProperties",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) CfnResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cfnResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) HealthCheckEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"healthCheckEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) HealthCheckIntervalSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"healthCheckIntervalSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) HealthCheckPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"healthCheckPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) HealthCheckPort() *string {
	var returns *string
	_jsii_.Get(
		j,
		"healthCheckPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) HealthCheckProtocol() *string {
	var returns *string
	_jsii_.Get(
		j,
		"healthCheckProtocol",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) HealthCheckTimeoutSeconds() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"healthCheckTimeoutSeconds",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) HealthyThresholdCount() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"healthyThresholdCount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) LogicalId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"logicalId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Matcher() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"matcher",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Port() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"port",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Protocol() *string {
	var returns *string
	_jsii_.Get(
		j,
		"protocol",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) ProtocolVersion() *string {
	var returns *string
	_jsii_.Get(
		j,
		"protocolVersion",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Ref() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ref",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Tags() awscdk.TagManager {
	var returns awscdk.TagManager
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) TargetGroupAttributes() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"targetGroupAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) Targets() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"targets",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) TargetType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) UnhealthyThresholdCount() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"unhealthyThresholdCount",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) UpdatedProperites() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"updatedProperites",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_CfnTargetGroup) VpcId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"vpcId",
		&returns,
	)
	return returns
}


// Create a new `AWS::ElasticLoadBalancingV2::TargetGroup`.
func NewCfnTargetGroup(scope awscdk.Construct, id *string, props *CfnTargetGroupProps) CfnTargetGroup {
	_init_.Initialize()

	j := jsiiProxy_CfnTargetGroup{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnTargetGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Create a new `AWS::ElasticLoadBalancingV2::TargetGroup`.
func NewCfnTargetGroup_Override(c CfnTargetGroup, scope awscdk.Construct, id *string, props *CfnTargetGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.CfnTargetGroup",
		[]interface{}{scope, id, props},
		c,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetHealthCheckEnabled(val interface{}) {
	_jsii_.Set(
		j,
		"healthCheckEnabled",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetHealthCheckIntervalSeconds(val *float64) {
	_jsii_.Set(
		j,
		"healthCheckIntervalSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetHealthCheckPath(val *string) {
	_jsii_.Set(
		j,
		"healthCheckPath",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetHealthCheckPort(val *string) {
	_jsii_.Set(
		j,
		"healthCheckPort",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetHealthCheckProtocol(val *string) {
	_jsii_.Set(
		j,
		"healthCheckProtocol",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetHealthCheckTimeoutSeconds(val *float64) {
	_jsii_.Set(
		j,
		"healthCheckTimeoutSeconds",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetHealthyThresholdCount(val *float64) {
	_jsii_.Set(
		j,
		"healthyThresholdCount",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetMatcher(val interface{}) {
	_jsii_.Set(
		j,
		"matcher",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetName(val *string) {
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetPort(val *float64) {
	_jsii_.Set(
		j,
		"port",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetProtocol(val *string) {
	_jsii_.Set(
		j,
		"protocol",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetProtocolVersion(val *string) {
	_jsii_.Set(
		j,
		"protocolVersion",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetTargetGroupAttributes(val interface{}) {
	_jsii_.Set(
		j,
		"targetGroupAttributes",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetTargets(val interface{}) {
	_jsii_.Set(
		j,
		"targets",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetTargetType(val *string) {
	_jsii_.Set(
		j,
		"targetType",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetUnhealthyThresholdCount(val *float64) {
	_jsii_.Set(
		j,
		"unhealthyThresholdCount",
		val,
	)
}

func (j *jsiiProxy_CfnTargetGroup) SetVpcId(val *string) {
	_jsii_.Set(
		j,
		"vpcId",
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
func CfnTargetGroup_IsCfnElement(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnTargetGroup",
		"isCfnElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a CfnResource.
// Experimental.
func CfnTargetGroup_IsCfnResource(construct constructs.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnTargetGroup",
		"isCfnResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func CfnTargetGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.CfnTargetGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func CfnTargetGroup_CFN_RESOURCE_TYPE_NAME() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"monocdk.aws_elasticloadbalancingv2.CfnTargetGroup",
		"CFN_RESOURCE_TYPE_NAME",
		&returns,
	)
	return returns
}

// Syntactic sugar for `addOverride(path, undefined)`.
// Experimental.
func (c *jsiiProxy_CfnTargetGroup) AddDeletionOverride(path *string) {
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
func (c *jsiiProxy_CfnTargetGroup) AddDependsOn(target awscdk.CfnResource) {
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
func (c *jsiiProxy_CfnTargetGroup) AddMetadata(key *string, value interface{}) {
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
func (c *jsiiProxy_CfnTargetGroup) AddOverride(path *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addOverride",
		[]interface{}{path, value},
	)
}

// Adds an override that deletes the value of a property from the resource definition.
// Experimental.
func (c *jsiiProxy_CfnTargetGroup) AddPropertyDeletionOverride(propertyPath *string) {
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
func (c *jsiiProxy_CfnTargetGroup) AddPropertyOverride(propertyPath *string, value interface{}) {
	_jsii_.InvokeVoid(
		c,
		"addPropertyOverride",
		[]interface{}{propertyPath, value},
	)
}

// Sets the deletion policy of the resource based on the removal policy specified.
// Experimental.
func (c *jsiiProxy_CfnTargetGroup) ApplyRemovalPolicy(policy awscdk.RemovalPolicy, options *awscdk.RemovalPolicyOptions) {
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
func (c *jsiiProxy_CfnTargetGroup) GetAtt(attributeName *string) awscdk.Reference {
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
func (c *jsiiProxy_CfnTargetGroup) GetMetadata(key *string) interface{} {
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
func (c *jsiiProxy_CfnTargetGroup) Inspect(inspector awscdk.TreeInspector) {
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
func (c *jsiiProxy_CfnTargetGroup) OnPrepare() {
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
func (c *jsiiProxy_CfnTargetGroup) OnSynthesize(session constructs.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTargetGroup) OnValidate() *[]*string {
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
func (c *jsiiProxy_CfnTargetGroup) OverrideLogicalId(newLogicalId *string) {
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
func (c *jsiiProxy_CfnTargetGroup) Prepare() {
	_jsii_.InvokeVoid(
		c,
		"prepare",
		nil, // no parameters
	)
}

func (c *jsiiProxy_CfnTargetGroup) RenderProperties(props *map[string]interface{}) *map[string]interface{} {
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
func (c *jsiiProxy_CfnTargetGroup) ShouldSynthesize() *bool {
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
func (c *jsiiProxy_CfnTargetGroup) Synthesize(session awscdk.ISynthesisSession) {
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
func (c *jsiiProxy_CfnTargetGroup) ToString() *string {
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
func (c *jsiiProxy_CfnTargetGroup) Validate() *[]*string {
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
func (c *jsiiProxy_CfnTargetGroup) ValidateProperties(_properties interface{}) {
	_jsii_.InvokeVoid(
		c,
		"validateProperties",
		[]interface{}{_properties},
	)
}

type CfnTargetGroup_MatcherProperty struct {
	// `CfnTargetGroup.MatcherProperty.GrpcCode`.
	GrpcCode *string `json:"grpcCode"`
	// `CfnTargetGroup.MatcherProperty.HttpCode`.
	HttpCode *string `json:"httpCode"`
}

type CfnTargetGroup_TargetDescriptionProperty struct {
	// `CfnTargetGroup.TargetDescriptionProperty.Id`.
	Id *string `json:"id"`
	// `CfnTargetGroup.TargetDescriptionProperty.AvailabilityZone`.
	AvailabilityZone *string `json:"availabilityZone"`
	// `CfnTargetGroup.TargetDescriptionProperty.Port`.
	Port *float64 `json:"port"`
}

type CfnTargetGroup_TargetGroupAttributeProperty struct {
	// `CfnTargetGroup.TargetGroupAttributeProperty.Key`.
	Key *string `json:"key"`
	// `CfnTargetGroup.TargetGroupAttributeProperty.Value`.
	Value *string `json:"value"`
}

// Properties for defining a `AWS::ElasticLoadBalancingV2::TargetGroup`.
type CfnTargetGroupProps struct {
	// `AWS::ElasticLoadBalancingV2::TargetGroup.HealthCheckEnabled`.
	HealthCheckEnabled interface{} `json:"healthCheckEnabled"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.HealthCheckIntervalSeconds`.
	HealthCheckIntervalSeconds *float64 `json:"healthCheckIntervalSeconds"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.HealthCheckPath`.
	HealthCheckPath *string `json:"healthCheckPath"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.HealthCheckPort`.
	HealthCheckPort *string `json:"healthCheckPort"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.HealthCheckProtocol`.
	HealthCheckProtocol *string `json:"healthCheckProtocol"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.HealthCheckTimeoutSeconds`.
	HealthCheckTimeoutSeconds *float64 `json:"healthCheckTimeoutSeconds"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.HealthyThresholdCount`.
	HealthyThresholdCount *float64 `json:"healthyThresholdCount"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.Matcher`.
	Matcher interface{} `json:"matcher"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.Name`.
	Name *string `json:"name"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.Port`.
	Port *float64 `json:"port"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.Protocol`.
	Protocol *string `json:"protocol"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.ProtocolVersion`.
	ProtocolVersion *string `json:"protocolVersion"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.Tags`.
	Tags *[]*awscdk.CfnTag `json:"tags"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.TargetGroupAttributes`.
	TargetGroupAttributes interface{} `json:"targetGroupAttributes"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.Targets`.
	Targets interface{} `json:"targets"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.TargetType`.
	TargetType *string `json:"targetType"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.UnhealthyThresholdCount`.
	UnhealthyThresholdCount *float64 `json:"unhealthyThresholdCount"`
	// `AWS::ElasticLoadBalancingV2::TargetGroup.VpcId`.
	VpcId *string `json:"vpcId"`
}

// The content type for a fixed response.
// Deprecated: superceded by `FixedResponseOptions`.
type ContentType string

const (
	ContentType_TEXT_PLAIN ContentType = "TEXT_PLAIN"
	ContentType_TEXT_CSS ContentType = "TEXT_CSS"
	ContentType_TEXT_HTML ContentType = "TEXT_HTML"
	ContentType_APPLICATION_JAVASCRIPT ContentType = "APPLICATION_JAVASCRIPT"
	ContentType_APPLICATION_JSON ContentType = "APPLICATION_JSON"
)

// A fixed response.
// Deprecated: superceded by `ListenerAction.fixedResponse()`.
type FixedResponse struct {
	// The HTTP response code (2XX, 4XX or 5XX).
	// Deprecated: superceded by `ListenerAction.fixedResponse()`.
	StatusCode *string `json:"statusCode"`
	// The content type.
	// Deprecated: superceded by `ListenerAction.fixedResponse()`.
	ContentType ContentType `json:"contentType"`
	// The message.
	// Deprecated: superceded by `ListenerAction.fixedResponse()`.
	MessageBody *string `json:"messageBody"`
}

// Options for `ListenerAction.fixedResponse()`.
// Experimental.
type FixedResponseOptions struct {
	// Content Type of the response.
	//
	// Valid Values: text/plain | text/css | text/html | application/javascript | application/json
	// Experimental.
	ContentType *string `json:"contentType"`
	// The response body.
	// Experimental.
	MessageBody *string `json:"messageBody"`
}

// Options for `ListenerAction.forward()`.
// Experimental.
type ForwardOptions struct {
	// For how long clients should be directed to the same target group.
	//
	// Range between 1 second and 7 days.
	// Experimental.
	StickinessDuration awscdk.Duration `json:"stickinessDuration"`
}

// Properties for configuring a health check.
// Experimental.
type HealthCheck struct {
	// Indicates whether health checks are enabled.
	//
	// If the target type is lambda,
	// health checks are disabled by default but can be enabled. If the target type
	// is instance or ip, health checks are always enabled and cannot be disabled.
	// Experimental.
	Enabled *bool `json:"enabled"`
	// GRPC code to use when checking for a successful response from a target.
	//
	// You can specify values between 0 and 99. You can specify multiple values
	// (for example, "0,1") or a range of values (for example, "0-5").
	// Experimental.
	HealthyGrpcCodes *string `json:"healthyGrpcCodes"`
	// HTTP code to use when checking for a successful response from a target.
	//
	// For Application Load Balancers, you can specify values between 200 and
	// 499, and the default value is 200. You can specify multiple values (for
	// example, "200,202") or a range of values (for example, "200-299").
	// Experimental.
	HealthyHttpCodes *string `json:"healthyHttpCodes"`
	// The number of consecutive health checks successes required before considering an unhealthy target healthy.
	//
	// For Application Load Balancers, the default is 5. For Network Load Balancers, the default is 3.
	// Experimental.
	HealthyThresholdCount *float64 `json:"healthyThresholdCount"`
	// The approximate number of seconds between health checks for an individual target.
	// Experimental.
	Interval awscdk.Duration `json:"interval"`
	// The ping path destination where Elastic Load Balancing sends health check requests.
	// Experimental.
	Path *string `json:"path"`
	// The port that the load balancer uses when performing health checks on the targets.
	// Experimental.
	Port *string `json:"port"`
	// The protocol the load balancer uses when performing health checks on targets.
	//
	// The TCP protocol is supported for health checks only if the protocol of the target group is TCP, TLS, UDP, or TCP_UDP.
	// The TLS, UDP, and TCP_UDP protocols are not supported for health checks.
	// Experimental.
	Protocol Protocol `json:"protocol"`
	// The amount of time, in seconds, during which no response from a target means a failed health check.
	//
	// For Application Load Balancers, the range is 2-60 seconds and the
	// default is 5 seconds. For Network Load Balancers, this is 10 seconds for
	// TCP and HTTPS health checks and 6 seconds for HTTP health checks.
	// Experimental.
	Timeout awscdk.Duration `json:"timeout"`
	// The number of consecutive health check failures required before considering a target unhealthy.
	//
	// For Application Load Balancers, the default is 2. For Network Load
	// Balancers, this value must be the same as the healthy threshold count.
	// Experimental.
	UnhealthyThresholdCount *float64 `json:"unhealthyThresholdCount"`
}

// Count of HTTP status originating from the load balancer.
//
// This count does not include any response codes generated by the targets.
// Experimental.
type HttpCodeElb string

const (
	HttpCodeElb_ELB_3XX_COUNT HttpCodeElb = "ELB_3XX_COUNT"
	HttpCodeElb_ELB_4XX_COUNT HttpCodeElb = "ELB_4XX_COUNT"
	HttpCodeElb_ELB_5XX_COUNT HttpCodeElb = "ELB_5XX_COUNT"
)

// Count of HTTP status originating from the targets.
// Experimental.
type HttpCodeTarget string

const (
	HttpCodeTarget_TARGET_2XX_COUNT HttpCodeTarget = "TARGET_2XX_COUNT"
	HttpCodeTarget_TARGET_3XX_COUNT HttpCodeTarget = "TARGET_3XX_COUNT"
	HttpCodeTarget_TARGET_4XX_COUNT HttpCodeTarget = "TARGET_4XX_COUNT"
	HttpCodeTarget_TARGET_5XX_COUNT HttpCodeTarget = "TARGET_5XX_COUNT"
)

// Properties to reference an existing listener.
// Experimental.
type IApplicationListener interface {
	awsec2.IConnectable
	awscdk.IResource
	// Add one or more certificates to this listener.
	// Deprecated: use `addCertificates()`
	AddCertificateArns(id *string, arns *[]*string)
	// Add one or more certificates to this listener.
	// Experimental.
	AddCertificates(id *string, certificates *[]IListenerCertificate)
	// Load balance incoming requests to the given target groups.
	//
	// It's possible to add conditions to the TargetGroups added in this way.
	// At least one TargetGroup must be added without conditions.
	// Experimental.
	AddTargetGroups(id *string, props *AddApplicationTargetGroupsProps)
	// Load balance incoming requests to the given load balancing targets.
	//
	// This method implicitly creates an ApplicationTargetGroup for the targets
	// involved.
	//
	// It's possible to add conditions to the targets added in this way. At least
	// one set of targets must be added without conditions.
	//
	// Returns: The newly created target group
	// Experimental.
	AddTargets(id *string, props *AddApplicationTargetsProps) ApplicationTargetGroup
	// Register that a connectable that has been added to this load balancer.
	//
	// Don't call this directly. It is called by ApplicationTargetGroup.
	// Experimental.
	RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port)
	// ARN of the listener.
	// Experimental.
	ListenerArn() *string
}

// The jsii proxy for IApplicationListener
type jsiiProxy_IApplicationListener struct {
	internal.Type__awsec2IConnectable
	internal.Type__awscdkIResource
}

func (i *jsiiProxy_IApplicationListener) AddCertificateArns(id *string, arns *[]*string) {
	_jsii_.InvokeVoid(
		i,
		"addCertificateArns",
		[]interface{}{id, arns},
	)
}

func (i *jsiiProxy_IApplicationListener) AddCertificates(id *string, certificates *[]IListenerCertificate) {
	_jsii_.InvokeVoid(
		i,
		"addCertificates",
		[]interface{}{id, certificates},
	)
}

func (i *jsiiProxy_IApplicationListener) AddTargetGroups(id *string, props *AddApplicationTargetGroupsProps) {
	_jsii_.InvokeVoid(
		i,
		"addTargetGroups",
		[]interface{}{id, props},
	)
}

func (i *jsiiProxy_IApplicationListener) AddTargets(id *string, props *AddApplicationTargetsProps) ApplicationTargetGroup {
	var returns ApplicationTargetGroup

	_jsii_.Invoke(
		i,
		"addTargets",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_IApplicationListener) RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port) {
	_jsii_.InvokeVoid(
		i,
		"registerConnectable",
		[]interface{}{connectable, portRange},
	)
}

func (j *jsiiProxy_IApplicationListener) ListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationListener) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationListener) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationListener) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationListener) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// An application load balancer.
// Experimental.
type IApplicationLoadBalancer interface {
	awsec2.IConnectable
	ILoadBalancerV2
	// Add a new listener to this load balancer.
	// Experimental.
	AddListener(id *string, props *BaseApplicationListenerProps) ApplicationListener
	// The IP Address Type for this load balancer.
	// Experimental.
	IpAddressType() IpAddressType
	// The ARN of this load balancer.
	// Experimental.
	LoadBalancerArn() *string
	// The VPC this load balancer has been created in (if available).
	//
	// If this interface is the result of an import call to fromApplicationLoadBalancerAttributes,
	// the vpc attribute will be undefined unless specified in the optional properties of that method.
	// Experimental.
	Vpc() awsec2.IVpc
}

// The jsii proxy for IApplicationLoadBalancer
type jsiiProxy_IApplicationLoadBalancer struct {
	internal.Type__awsec2IConnectable
	jsiiProxy_ILoadBalancerV2
}

func (i *jsiiProxy_IApplicationLoadBalancer) AddListener(id *string, props *BaseApplicationListenerProps) ApplicationListener {
	var returns ApplicationListener

	_jsii_.Invoke(
		i,
		"addListener",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) IpAddressType() IpAddressType {
	var returns IpAddressType
	_jsii_.Get(
		j,
		"ipAddressType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) LoadBalancerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) Connections() awsec2.Connections {
	var returns awsec2.Connections
	_jsii_.Get(
		j,
		"connections",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) LoadBalancerCanonicalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) LoadBalancerDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_IApplicationLoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// Interface for constructs that can be targets of an application load balancer.
// Experimental.
type IApplicationLoadBalancerTarget interface {
	// Attach load-balanced target to a TargetGroup.
	//
	// May return JSON to directly add to the [Targets] list, or return undefined
	// if the target will register itself with the load balancer.
	// Experimental.
	AttachToApplicationTargetGroup(targetGroup IApplicationTargetGroup) *LoadBalancerTargetProps
}

// The jsii proxy for IApplicationLoadBalancerTarget
type jsiiProxy_IApplicationLoadBalancerTarget struct {
	_ byte // padding
}

func (i *jsiiProxy_IApplicationLoadBalancerTarget) AttachToApplicationTargetGroup(targetGroup IApplicationTargetGroup) *LoadBalancerTargetProps {
	var returns *LoadBalancerTargetProps

	_jsii_.Invoke(
		i,
		"attachToApplicationTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// A Target Group for Application Load Balancers.
// Experimental.
type IApplicationTargetGroup interface {
	ITargetGroup
	// Add a load balancing target to this target group.
	// Experimental.
	AddTarget(targets ...IApplicationLoadBalancerTarget)
	// Register a connectable as a member of this target group.
	//
	// Don't call this directly. It will be called by load balancing targets.
	// Experimental.
	RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port)
	// Register a listener that is load balancing to this target group.
	//
	// Don't call this directly. It will be called by listeners.
	// Experimental.
	RegisterListener(listener IApplicationListener, associatingConstruct constructs.IConstruct)
}

// The jsii proxy for IApplicationTargetGroup
type jsiiProxy_IApplicationTargetGroup struct {
	jsiiProxy_ITargetGroup
}

func (i *jsiiProxy_IApplicationTargetGroup) AddTarget(targets ...IApplicationLoadBalancerTarget) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		i,
		"addTarget",
		args,
	)
}

func (i *jsiiProxy_IApplicationTargetGroup) RegisterConnectable(connectable awsec2.IConnectable, portRange awsec2.Port) {
	_jsii_.InvokeVoid(
		i,
		"registerConnectable",
		[]interface{}{connectable, portRange},
	)
}

func (i *jsiiProxy_IApplicationTargetGroup) RegisterListener(listener IApplicationListener, associatingConstruct constructs.IConstruct) {
	_jsii_.InvokeVoid(
		i,
		"registerListener",
		[]interface{}{listener, associatingConstruct},
	)
}

// Interface for listener actions.
// Experimental.
type IListenerAction interface {
	// Render the actions in this chain.
	// Experimental.
	RenderActions() *[]*CfnListener_ActionProperty
}

// The jsii proxy for IListenerAction
type jsiiProxy_IListenerAction struct {
	_ byte // padding
}

func (i *jsiiProxy_IListenerAction) RenderActions() *[]*CfnListener_ActionProperty {
	var returns *[]*CfnListener_ActionProperty

	_jsii_.Invoke(
		i,
		"renderActions",
		nil, // no parameters
		&returns,
	)

	return returns
}

// A certificate source for an ELBv2 listener.
// Experimental.
type IListenerCertificate interface {
	// The ARN of the certificate to use.
	// Experimental.
	CertificateArn() *string
}

// The jsii proxy for IListenerCertificate
type jsiiProxy_IListenerCertificate struct {
	_ byte // padding
}

func (j *jsiiProxy_IListenerCertificate) CertificateArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certificateArn",
		&returns,
	)
	return returns
}

// Experimental.
type ILoadBalancerV2 interface {
	awscdk.IResource
	// The canonical hosted zone ID of this load balancer.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	LoadBalancerCanonicalHostedZoneId() *string
	// The DNS name of this load balancer.
	//
	// TODO: EXAMPLE
	//
	// Experimental.
	LoadBalancerDnsName() *string
}

// The jsii proxy for ILoadBalancerV2
type jsiiProxy_ILoadBalancerV2 struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_ILoadBalancerV2) LoadBalancerCanonicalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ILoadBalancerV2) LoadBalancerDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerDnsName",
		&returns,
	)
	return returns
}

// Properties to reference an existing listener.
// Experimental.
type INetworkListener interface {
	awscdk.IResource
	// ARN of the listener.
	// Experimental.
	ListenerArn() *string
}

// The jsii proxy for INetworkListener
type jsiiProxy_INetworkListener struct {
	internal.Type__awscdkIResource
}

func (j *jsiiProxy_INetworkListener) ListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerArn",
		&returns,
	)
	return returns
}

// Properties for adding a certificate to a listener.
//
// This interface exists for backwards compatibility.
// Deprecated: Use IListenerCertificate instead
type INetworkListenerCertificateProps interface {
	IListenerCertificate
}

// The jsii proxy for INetworkListenerCertificateProps
type jsiiProxy_INetworkListenerCertificateProps struct {
	jsiiProxy_IListenerCertificate
}

// A network load balancer.
// Experimental.
type INetworkLoadBalancer interface {
	ILoadBalancerV2
	awsec2.IVpcEndpointServiceLoadBalancer
	// Add a listener to this load balancer.
	//
	// Returns: The newly created listener
	// Experimental.
	AddListener(id *string, props *BaseNetworkListenerProps) NetworkListener
	// The VPC this load balancer has been created in (if available).
	// Experimental.
	Vpc() awsec2.IVpc
}

// The jsii proxy for INetworkLoadBalancer
type jsiiProxy_INetworkLoadBalancer struct {
	jsiiProxy_ILoadBalancerV2
	internal.Type__awsec2IVpcEndpointServiceLoadBalancer
}

func (i *jsiiProxy_INetworkLoadBalancer) AddListener(id *string, props *BaseNetworkListenerProps) NetworkListener {
	var returns NetworkListener

	_jsii_.Invoke(
		i,
		"addListener",
		[]interface{}{id, props},
		&returns,
	)

	return returns
}

func (j *jsiiProxy_INetworkLoadBalancer) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_INetworkLoadBalancer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_INetworkLoadBalancer) LoadBalancerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_INetworkLoadBalancer) LoadBalancerCanonicalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_INetworkLoadBalancer) LoadBalancerDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_INetworkLoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_INetworkLoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

// Interface for constructs that can be targets of an network load balancer.
// Experimental.
type INetworkLoadBalancerTarget interface {
	// Attach load-balanced target to a TargetGroup.
	//
	// May return JSON to directly add to the [Targets] list, or return undefined
	// if the target will register itself with the load balancer.
	// Experimental.
	AttachToNetworkTargetGroup(targetGroup INetworkTargetGroup) *LoadBalancerTargetProps
}

// The jsii proxy for INetworkLoadBalancerTarget
type jsiiProxy_INetworkLoadBalancerTarget struct {
	_ byte // padding
}

func (i *jsiiProxy_INetworkLoadBalancerTarget) AttachToNetworkTargetGroup(targetGroup INetworkTargetGroup) *LoadBalancerTargetProps {
	var returns *LoadBalancerTargetProps

	_jsii_.Invoke(
		i,
		"attachToNetworkTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// A network target group.
// Experimental.
type INetworkTargetGroup interface {
	ITargetGroup
	// Add a load balancing target to this target group.
	// Experimental.
	AddTarget(targets ...INetworkLoadBalancerTarget)
	// Register a listener that is load balancing to this target group.
	//
	// Don't call this directly. It will be called by listeners.
	// Experimental.
	RegisterListener(listener INetworkListener)
}

// The jsii proxy for INetworkTargetGroup
type jsiiProxy_INetworkTargetGroup struct {
	jsiiProxy_ITargetGroup
}

func (i *jsiiProxy_INetworkTargetGroup) AddTarget(targets ...INetworkLoadBalancerTarget) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		i,
		"addTarget",
		args,
	)
}

func (i *jsiiProxy_INetworkTargetGroup) RegisterListener(listener INetworkListener) {
	_jsii_.InvokeVoid(
		i,
		"registerListener",
		[]interface{}{listener},
	)
}

// A target group.
// Experimental.
type ITargetGroup interface {
	awscdk.IConstruct
	// A token representing a list of ARNs of the load balancers that route traffic to this target group.
	// Experimental.
	LoadBalancerArns() *string
	// Return an object to depend on the listeners added to this target group.
	// Experimental.
	LoadBalancerAttached() awscdk.IDependable
	// ARN of the target group.
	// Experimental.
	TargetGroupArn() *string
}

// The jsii proxy for ITargetGroup
type jsiiProxy_ITargetGroup struct {
	internal.Type__awscdkIConstruct
}

func (j *jsiiProxy_ITargetGroup) LoadBalancerArns() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITargetGroup) LoadBalancerAttached() awscdk.IDependable {
	var returns awscdk.IDependable
	_jsii_.Get(
		j,
		"loadBalancerAttached",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ITargetGroup) TargetGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupArn",
		&returns,
	)
	return returns
}

// An EC2 instance that is the target for load balancing.
//
// If you register a target of this type, you are responsible for making
// sure the load balancer's security group can connect to the instance.
// Deprecated: Use IpTarget from the
type InstanceTarget interface {
	IApplicationLoadBalancerTarget
	INetworkLoadBalancerTarget
	AttachToApplicationTargetGroup(targetGroup IApplicationTargetGroup) *LoadBalancerTargetProps
	AttachToNetworkTargetGroup(targetGroup INetworkTargetGroup) *LoadBalancerTargetProps
}

// The jsii proxy struct for InstanceTarget
type jsiiProxy_InstanceTarget struct {
	jsiiProxy_IApplicationLoadBalancerTarget
	jsiiProxy_INetworkLoadBalancerTarget
}

// Create a new Instance target.
// Deprecated: Use IpTarget from the
func NewInstanceTarget(instanceId *string, port *float64) InstanceTarget {
	_init_.Initialize()

	j := jsiiProxy_InstanceTarget{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.InstanceTarget",
		[]interface{}{instanceId, port},
		&j,
	)

	return &j
}

// Create a new Instance target.
// Deprecated: Use IpTarget from the
func NewInstanceTarget_Override(i InstanceTarget, instanceId *string, port *float64) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.InstanceTarget",
		[]interface{}{instanceId, port},
		i,
	)
}

// Register this instance target with a load balancer.
//
// Don't call this, it is called automatically when you add the target to a
// load balancer.
// Deprecated: Use IpTarget from the
func (i *jsiiProxy_InstanceTarget) AttachToApplicationTargetGroup(targetGroup IApplicationTargetGroup) *LoadBalancerTargetProps {
	var returns *LoadBalancerTargetProps

	_jsii_.Invoke(
		i,
		"attachToApplicationTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// Register this instance target with a load balancer.
//
// Don't call this, it is called automatically when you add the target to a
// load balancer.
// Deprecated: Use IpTarget from the
func (i *jsiiProxy_InstanceTarget) AttachToNetworkTargetGroup(targetGroup INetworkTargetGroup) *LoadBalancerTargetProps {
	var returns *LoadBalancerTargetProps

	_jsii_.Invoke(
		i,
		"attachToNetworkTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// What kind of addresses to allocate to the load balancer.
// Experimental.
type IpAddressType string

const (
	IpAddressType_IPV4 IpAddressType = "IPV4"
	IpAddressType_DUAL_STACK IpAddressType = "DUAL_STACK"
)

// An IP address that is a target for load balancing.
//
// Specify IP addresses from the subnets of the virtual private cloud (VPC) for
// the target group, the RFC 1918 range (10.0.0.0/8, 172.16.0.0/12, and
// 192.168.0.0/16), and the RFC 6598 range (100.64.0.0/10). You can't specify
// publicly routable IP addresses.
//
// If you register a target of this type, you are responsible for making
// sure the load balancer's security group can send packets to the IP address.
// Deprecated: Use IpTarget from the
type IpTarget interface {
	IApplicationLoadBalancerTarget
	INetworkLoadBalancerTarget
	AttachToApplicationTargetGroup(targetGroup IApplicationTargetGroup) *LoadBalancerTargetProps
	AttachToNetworkTargetGroup(targetGroup INetworkTargetGroup) *LoadBalancerTargetProps
}

// The jsii proxy struct for IpTarget
type jsiiProxy_IpTarget struct {
	jsiiProxy_IApplicationLoadBalancerTarget
	jsiiProxy_INetworkLoadBalancerTarget
}

// Create a new IPAddress target.
//
// The availabilityZone parameter determines whether the target receives
// traffic from the load balancer nodes in the specified Availability Zone
// or from all enabled Availability Zones for the load balancer.
//
// This parameter is not supported if the target type of the target group
// is instance. If the IP address is in a subnet of the VPC for the target
// group, the Availability Zone is automatically detected and this
// parameter is optional. If the IP address is outside the VPC, this
// parameter is required.
//
// With an Application Load Balancer, if the IP address is outside the VPC
// for the target group, the only supported value is all.
//
// Default is automatic.
// Deprecated: Use IpTarget from the
func NewIpTarget(ipAddress *string, port *float64, availabilityZone *string) IpTarget {
	_init_.Initialize()

	j := jsiiProxy_IpTarget{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.IpTarget",
		[]interface{}{ipAddress, port, availabilityZone},
		&j,
	)

	return &j
}

// Create a new IPAddress target.
//
// The availabilityZone parameter determines whether the target receives
// traffic from the load balancer nodes in the specified Availability Zone
// or from all enabled Availability Zones for the load balancer.
//
// This parameter is not supported if the target type of the target group
// is instance. If the IP address is in a subnet of the VPC for the target
// group, the Availability Zone is automatically detected and this
// parameter is optional. If the IP address is outside the VPC, this
// parameter is required.
//
// With an Application Load Balancer, if the IP address is outside the VPC
// for the target group, the only supported value is all.
//
// Default is automatic.
// Deprecated: Use IpTarget from the
func NewIpTarget_Override(i IpTarget, ipAddress *string, port *float64, availabilityZone *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.IpTarget",
		[]interface{}{ipAddress, port, availabilityZone},
		i,
	)
}

// Register this instance target with a load balancer.
//
// Don't call this, it is called automatically when you add the target to a
// load balancer.
// Deprecated: Use IpTarget from the
func (i *jsiiProxy_IpTarget) AttachToApplicationTargetGroup(targetGroup IApplicationTargetGroup) *LoadBalancerTargetProps {
	var returns *LoadBalancerTargetProps

	_jsii_.Invoke(
		i,
		"attachToApplicationTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// Register this instance target with a load balancer.
//
// Don't call this, it is called automatically when you add the target to a
// load balancer.
// Deprecated: Use IpTarget from the
func (i *jsiiProxy_IpTarget) AttachToNetworkTargetGroup(targetGroup INetworkTargetGroup) *LoadBalancerTargetProps {
	var returns *LoadBalancerTargetProps

	_jsii_.Invoke(
		i,
		"attachToNetworkTargetGroup",
		[]interface{}{targetGroup},
		&returns,
	)

	return returns
}

// What to do when a client makes a request to a listener.
//
// Some actions can be combined with other ones (specifically,
// you can perform authentication before serving the request).
//
// Multiple actions form a linked chain; the chain must always terminate in a
// *(weighted)forward*, *fixedResponse* or *redirect* action.
//
// If an action supports chaining, the next action can be indicated
// by passing it in the `next` property.
//
// (Called `ListenerAction` instead of the more strictly correct
// `ListenerAction` because this is the class most users interact
// with, and we want to make it not too visually overwhelming).
// Experimental.
type ListenerAction interface {
	IListenerAction
	Next() ListenerAction
	Bind(scope awscdk.Construct, listener IApplicationListener, associatingConstruct awscdk.IConstruct)
	RenderActions() *[]*CfnListener_ActionProperty
	Renumber(actions *[]*CfnListener_ActionProperty) *[]*CfnListener_ActionProperty
}

// The jsii proxy struct for ListenerAction
type jsiiProxy_ListenerAction struct {
	jsiiProxy_IListenerAction
}

func (j *jsiiProxy_ListenerAction) Next() ListenerAction {
	var returns ListenerAction
	_jsii_.Get(
		j,
		"next",
		&returns,
	)
	return returns
}


// Create an instance of ListenerAction.
//
// The default class should be good enough for most cases and
// should be created by using one of the static factory functions,
// but allow overriding to make sure we allow flexibility for the future.
// Experimental.
func NewListenerAction(actionJson *CfnListener_ActionProperty, next ListenerAction) ListenerAction {
	_init_.Initialize()

	j := jsiiProxy_ListenerAction{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ListenerAction",
		[]interface{}{actionJson, next},
		&j,
	)

	return &j
}

// Create an instance of ListenerAction.
//
// The default class should be good enough for most cases and
// should be created by using one of the static factory functions,
// but allow overriding to make sure we allow flexibility for the future.
// Experimental.
func NewListenerAction_Override(l ListenerAction, actionJson *CfnListener_ActionProperty, next ListenerAction) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ListenerAction",
		[]interface{}{actionJson, next},
		l,
	)
}

// Authenticate using an identity provider (IdP) that is compliant with OpenID Connect (OIDC).
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/listener-authenticate-users.html#oidc-requirements
//
// Experimental.
func ListenerAction_AuthenticateOidc(options *AuthenticateOidcOptions) ListenerAction {
	_init_.Initialize()

	var returns ListenerAction

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerAction",
		"authenticateOidc",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Return a fixed response.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#fixed-response-actions
//
// Experimental.
func ListenerAction_FixedResponse(statusCode *float64, options *FixedResponseOptions) ListenerAction {
	_init_.Initialize()

	var returns ListenerAction

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerAction",
		"fixedResponse",
		[]interface{}{statusCode, options},
		&returns,
	)

	return returns
}

// Forward to one or more Target Groups.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#forward-actions
//
// Experimental.
func ListenerAction_Forward(targetGroups *[]IApplicationTargetGroup, options *ForwardOptions) ListenerAction {
	_init_.Initialize()

	var returns ListenerAction

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerAction",
		"forward",
		[]interface{}{targetGroups, options},
		&returns,
	)

	return returns
}

// Redirect to a different URI.
//
// A URI consists of the following components:
// protocol://hostname:port/path?query. You must modify at least one of the
// following components to avoid a redirect loop: protocol, hostname, port, or
// path. Any components that you do not modify retain their original values.
//
// You can reuse URI components using the following reserved keywords:
//
// - `#{protocol}`
// - `#{host}`
// - `#{port}`
// - `#{path}` (the leading "/" is removed)
// - `#{query}`
//
// For example, you can change the path to "/new/#{path}", the hostname to
// "example.#{host}", or the query to "#{query}&value=xyz".
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#redirect-actions
//
// Experimental.
func ListenerAction_Redirect(options *RedirectOptions) ListenerAction {
	_init_.Initialize()

	var returns ListenerAction

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerAction",
		"redirect",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Forward to one or more Target Groups which are weighted differently.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-listeners.html#forward-actions
//
// Experimental.
func ListenerAction_WeightedForward(targetGroups *[]*WeightedTargetGroup, options *ForwardOptions) ListenerAction {
	_init_.Initialize()

	var returns ListenerAction

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerAction",
		"weightedForward",
		[]interface{}{targetGroups, options},
		&returns,
	)

	return returns
}

// Called when the action is being used in a listener.
// Experimental.
func (l *jsiiProxy_ListenerAction) Bind(scope awscdk.Construct, listener IApplicationListener, associatingConstruct awscdk.IConstruct) {
	_jsii_.InvokeVoid(
		l,
		"bind",
		[]interface{}{scope, listener, associatingConstruct},
	)
}

// Render the actions in this chain.
// Experimental.
func (l *jsiiProxy_ListenerAction) RenderActions() *[]*CfnListener_ActionProperty {
	var returns *[]*CfnListener_ActionProperty

	_jsii_.Invoke(
		l,
		"renderActions",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Renumber the "order" fields in the actions array.
//
// We don't number for 0 or 1 elements, but otherwise number them 1...#actions
// so ELB knows about the right order.
//
// Do this in `ListenerAction` instead of in `Listener` so that we give
// users the opportunity to override by subclassing and overriding `renderActions`.
// Experimental.
func (l *jsiiProxy_ListenerAction) Renumber(actions *[]*CfnListener_ActionProperty) *[]*CfnListener_ActionProperty {
	var returns *[]*CfnListener_ActionProperty

	_jsii_.Invoke(
		l,
		"renumber",
		[]interface{}{actions},
		&returns,
	)

	return returns
}

// A certificate source for an ELBv2 listener.
// Experimental.
type ListenerCertificate interface {
	IListenerCertificate
	CertificateArn() *string
}

// The jsii proxy struct for ListenerCertificate
type jsiiProxy_ListenerCertificate struct {
	jsiiProxy_IListenerCertificate
}

func (j *jsiiProxy_ListenerCertificate) CertificateArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certificateArn",
		&returns,
	)
	return returns
}


// Experimental.
func NewListenerCertificate(certificateArn *string) ListenerCertificate {
	_init_.Initialize()

	j := jsiiProxy_ListenerCertificate{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ListenerCertificate",
		[]interface{}{certificateArn},
		&j,
	)

	return &j
}

// Experimental.
func NewListenerCertificate_Override(l ListenerCertificate, certificateArn *string) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ListenerCertificate",
		[]interface{}{certificateArn},
		l,
	)
}

// Use any certificate, identified by its ARN, as a listener certificate.
// Experimental.
func ListenerCertificate_FromArn(certificateArn *string) ListenerCertificate {
	_init_.Initialize()

	var returns ListenerCertificate

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCertificate",
		"fromArn",
		[]interface{}{certificateArn},
		&returns,
	)

	return returns
}

// Use an ACM certificate as a listener certificate.
// Experimental.
func ListenerCertificate_FromCertificateManager(acmCertificate awscertificatemanager.ICertificate) ListenerCertificate {
	_init_.Initialize()

	var returns ListenerCertificate

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCertificate",
		"fromCertificateManager",
		[]interface{}{acmCertificate},
		&returns,
	)

	return returns
}

// ListenerCondition providers definition.
// Experimental.
type ListenerCondition interface {
	RenderRawCondition() interface{}
}

// The jsii proxy struct for ListenerCondition
type jsiiProxy_ListenerCondition struct {
	_ byte // padding
}

// Experimental.
func NewListenerCondition_Override(l ListenerCondition) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.ListenerCondition",
		nil, // no parameters
		l,
	)
}

// Create a host-header listener rule condition.
// Experimental.
func ListenerCondition_HostHeaders(values *[]*string) ListenerCondition {
	_init_.Initialize()

	var returns ListenerCondition

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCondition",
		"hostHeaders",
		[]interface{}{values},
		&returns,
	)

	return returns
}

// Create a http-header listener rule condition.
// Experimental.
func ListenerCondition_HttpHeader(name *string, values *[]*string) ListenerCondition {
	_init_.Initialize()

	var returns ListenerCondition

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCondition",
		"httpHeader",
		[]interface{}{name, values},
		&returns,
	)

	return returns
}

// Create a http-request-method listener rule condition.
// Experimental.
func ListenerCondition_HttpRequestMethods(values *[]*string) ListenerCondition {
	_init_.Initialize()

	var returns ListenerCondition

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCondition",
		"httpRequestMethods",
		[]interface{}{values},
		&returns,
	)

	return returns
}

// Create a path-pattern listener rule condition.
// Experimental.
func ListenerCondition_PathPatterns(values *[]*string) ListenerCondition {
	_init_.Initialize()

	var returns ListenerCondition

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCondition",
		"pathPatterns",
		[]interface{}{values},
		&returns,
	)

	return returns
}

// Create a query-string listener rule condition.
// Experimental.
func ListenerCondition_QueryStrings(values *[]*QueryStringCondition) ListenerCondition {
	_init_.Initialize()

	var returns ListenerCondition

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCondition",
		"queryStrings",
		[]interface{}{values},
		&returns,
	)

	return returns
}

// Create a source-ip listener rule condition.
// Experimental.
func ListenerCondition_SourceIps(values *[]*string) ListenerCondition {
	_init_.Initialize()

	var returns ListenerCondition

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.ListenerCondition",
		"sourceIps",
		[]interface{}{values},
		&returns,
	)

	return returns
}

// Render the raw Cfn listener rule condition object.
// Experimental.
func (l *jsiiProxy_ListenerCondition) RenderRawCondition() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		l,
		"renderRawCondition",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Result of attaching a target to load balancer.
// Experimental.
type LoadBalancerTargetProps struct {
	// What kind of target this is.
	// Experimental.
	TargetType TargetType `json:"targetType"`
	// JSON representing the target's direct addition to the TargetGroup list.
	//
	// May be omitted if the target is going to register itself later.
	// Experimental.
	TargetJson interface{} `json:"targetJson"`
}

// Options for `NetworkListenerAction.forward()`.
// Experimental.
type NetworkForwardOptions struct {
	// For how long clients should be directed to the same target group.
	//
	// Range between 1 second and 7 days.
	// Experimental.
	StickinessDuration awscdk.Duration `json:"stickinessDuration"`
}

// Define a Network Listener.
// Experimental.
type NetworkListener interface {
	BaseListener
	INetworkListener
	Env() *awscdk.ResourceEnvironment
	ListenerArn() *string
	LoadBalancer() INetworkLoadBalancer
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	AddAction(_id *string, props *AddNetworkActionProps)
	AddTargetGroups(_id *string, targetGroups ...INetworkTargetGroup)
	AddTargets(id *string, props *AddNetworkTargetsProps) NetworkTargetGroup
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

// The jsii proxy struct for NetworkListener
type jsiiProxy_NetworkListener struct {
	jsiiProxy_BaseListener
	jsiiProxy_INetworkListener
}

func (j *jsiiProxy_NetworkListener) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkListener) ListenerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"listenerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkListener) LoadBalancer() INetworkLoadBalancer {
	var returns INetworkLoadBalancer
	_jsii_.Get(
		j,
		"loadBalancer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkListener) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkListener) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkListener) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}


// Experimental.
func NewNetworkListener(scope constructs.Construct, id *string, props *NetworkListenerProps) NetworkListener {
	_init_.Initialize()

	j := jsiiProxy_NetworkListener{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkListener",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewNetworkListener_Override(n NetworkListener, scope constructs.Construct, id *string, props *NetworkListenerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkListener",
		[]interface{}{scope, id, props},
		n,
	)
}

// Looks up a network listener.
// Experimental.
func NetworkListener_FromLookup(scope constructs.Construct, id *string, options *NetworkListenerLookupOptions) INetworkListener {
	_init_.Initialize()

	var returns INetworkListener

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkListener",
		"fromLookup",
		[]interface{}{scope, id, options},
		&returns,
	)

	return returns
}

// Import an existing listener.
// Experimental.
func NetworkListener_FromNetworkListenerArn(scope constructs.Construct, id *string, networkListenerArn *string) INetworkListener {
	_init_.Initialize()

	var returns INetworkListener

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkListener",
		"fromNetworkListenerArn",
		[]interface{}{scope, id, networkListenerArn},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func NetworkListener_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkListener",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func NetworkListener_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkListener",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Perform the given Action on incoming requests.
//
// This allows full control of the default Action of the load balancer,
// including weighted forwarding. See the `NetworkListenerAction` class for
// all options.
// Experimental.
func (n *jsiiProxy_NetworkListener) AddAction(_id *string, props *AddNetworkActionProps) {
	_jsii_.InvokeVoid(
		n,
		"addAction",
		[]interface{}{_id, props},
	)
}

// Load balance incoming requests to the given target groups.
//
// All target groups will be load balanced to with equal weight and without
// stickiness. For a more complex configuration than that, use `addAction()`.
// Experimental.
func (n *jsiiProxy_NetworkListener) AddTargetGroups(_id *string, targetGroups ...INetworkTargetGroup) {
	args := []interface{}{_id}
	for _, a := range targetGroups {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		n,
		"addTargetGroups",
		args,
	)
}

// Load balance incoming requests to the given load balancing targets.
//
// This method implicitly creates a NetworkTargetGroup for the targets
// involved, and a 'forward' action to route traffic to the given TargetGroup.
//
// If you want more control over the precise setup, create the TargetGroup
// and use `addAction` yourself.
//
// It's possible to add conditions to the targets added in this way. At least
// one set of targets must be added without conditions.
//
// Returns: The newly created target group
// Experimental.
func (n *jsiiProxy_NetworkListener) AddTargets(id *string, props *AddNetworkTargetsProps) NetworkTargetGroup {
	var returns NetworkTargetGroup

	_jsii_.Invoke(
		n,
		"addTargets",
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
func (n *jsiiProxy_NetworkListener) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		n,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (n *jsiiProxy_NetworkListener) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkListener) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkListener) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkListener) OnPrepare() {
	_jsii_.InvokeVoid(
		n,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NetworkListener) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
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
func (n *jsiiProxy_NetworkListener) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkListener) Prepare() {
	_jsii_.InvokeVoid(
		n,
		"prepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NetworkListener) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (n *jsiiProxy_NetworkListener) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Validate this listener.
// Experimental.
func (n *jsiiProxy_NetworkListener) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// What to do when a client makes a request to a listener.
//
// Some actions can be combined with other ones (specifically,
// you can perform authentication before serving the request).
//
// Multiple actions form a linked chain; the chain must always terminate in a
// *(weighted)forward*, *fixedResponse* or *redirect* action.
//
// If an action supports chaining, the next action can be indicated
// by passing it in the `next` property.
// Experimental.
type NetworkListenerAction interface {
	IListenerAction
	Next() NetworkListenerAction
	Bind(scope awscdk.Construct, listener INetworkListener)
	RenderActions() *[]*CfnListener_ActionProperty
	Renumber(actions *[]*CfnListener_ActionProperty) *[]*CfnListener_ActionProperty
}

// The jsii proxy struct for NetworkListenerAction
type jsiiProxy_NetworkListenerAction struct {
	jsiiProxy_IListenerAction
}

func (j *jsiiProxy_NetworkListenerAction) Next() NetworkListenerAction {
	var returns NetworkListenerAction
	_jsii_.Get(
		j,
		"next",
		&returns,
	)
	return returns
}


// Create an instance of NetworkListenerAction.
//
// The default class should be good enough for most cases and
// should be created by using one of the static factory functions,
// but allow overriding to make sure we allow flexibility for the future.
// Experimental.
func NewNetworkListenerAction(actionJson *CfnListener_ActionProperty, next NetworkListenerAction) NetworkListenerAction {
	_init_.Initialize()

	j := jsiiProxy_NetworkListenerAction{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkListenerAction",
		[]interface{}{actionJson, next},
		&j,
	)

	return &j
}

// Create an instance of NetworkListenerAction.
//
// The default class should be good enough for most cases and
// should be created by using one of the static factory functions,
// but allow overriding to make sure we allow flexibility for the future.
// Experimental.
func NewNetworkListenerAction_Override(n NetworkListenerAction, actionJson *CfnListener_ActionProperty, next NetworkListenerAction) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkListenerAction",
		[]interface{}{actionJson, next},
		n,
	)
}

// Forward to one or more Target Groups.
// Experimental.
func NetworkListenerAction_Forward(targetGroups *[]INetworkTargetGroup, options *NetworkForwardOptions) NetworkListenerAction {
	_init_.Initialize()

	var returns NetworkListenerAction

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkListenerAction",
		"forward",
		[]interface{}{targetGroups, options},
		&returns,
	)

	return returns
}

// Forward to one or more Target Groups which are weighted differently.
// Experimental.
func NetworkListenerAction_WeightedForward(targetGroups *[]*NetworkWeightedTargetGroup, options *NetworkForwardOptions) NetworkListenerAction {
	_init_.Initialize()

	var returns NetworkListenerAction

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkListenerAction",
		"weightedForward",
		[]interface{}{targetGroups, options},
		&returns,
	)

	return returns
}

// Called when the action is being used in a listener.
// Experimental.
func (n *jsiiProxy_NetworkListenerAction) Bind(scope awscdk.Construct, listener INetworkListener) {
	_jsii_.InvokeVoid(
		n,
		"bind",
		[]interface{}{scope, listener},
	)
}

// Render the actions in this chain.
// Experimental.
func (n *jsiiProxy_NetworkListenerAction) RenderActions() *[]*CfnListener_ActionProperty {
	var returns *[]*CfnListener_ActionProperty

	_jsii_.Invoke(
		n,
		"renderActions",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Renumber the "order" fields in the actions array.
//
// We don't number for 0 or 1 elements, but otherwise number them 1...#actions
// so ELB knows about the right order.
//
// Do this in `NetworkListenerAction` instead of in `Listener` so that we give
// users the opportunity to override by subclassing and overriding `renderActions`.
// Experimental.
func (n *jsiiProxy_NetworkListenerAction) Renumber(actions *[]*CfnListener_ActionProperty) *[]*CfnListener_ActionProperty {
	var returns *[]*CfnListener_ActionProperty

	_jsii_.Invoke(
		n,
		"renumber",
		[]interface{}{actions},
		&returns,
	)

	return returns
}

// Options for looking up a network listener.
// Experimental.
type NetworkListenerLookupOptions struct {
	// Filter listeners by listener port.
	// Experimental.
	ListenerPort *float64 `json:"listenerPort"`
	// Filter listeners by associated load balancer arn.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// Filter listeners by associated load balancer tags.
	// Experimental.
	LoadBalancerTags *map[string]*string `json:"loadBalancerTags"`
	// Protocol of the listener port.
	// Experimental.
	ListenerProtocol Protocol `json:"listenerProtocol"`
}

// Properties for a Network Listener attached to a Load Balancer.
// Experimental.
type NetworkListenerProps struct {
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// Certificate list of ACM cert ARNs.
	// Experimental.
	Certificates *[]IListenerCertificate `json:"certificates"`
	// Default action to take for requests to this listener.
	//
	// This allows full control of the default Action of the load balancer,
	// including weighted forwarding. See the `NetworkListenerAction` class for
	// all options.
	//
	// Cannot be specified together with `defaultTargetGroups`.
	// Experimental.
	DefaultAction NetworkListenerAction `json:"defaultAction"`
	// Default target groups to load balance to.
	//
	// All target groups will be load balanced to with equal weight and without
	// stickiness. For a more complex configuration than that, use
	// either `defaultAction` or `addAction()`.
	//
	// Cannot be specified together with `defaultAction`.
	// Experimental.
	DefaultTargetGroups *[]INetworkTargetGroup `json:"defaultTargetGroups"`
	// Protocol for listener, expects TCP, TLS, UDP, or TCP_UDP.
	// Experimental.
	Protocol Protocol `json:"protocol"`
	// SSL Policy.
	// Experimental.
	SslPolicy SslPolicy `json:"sslPolicy"`
	// The load balancer to attach this listener to.
	// Experimental.
	LoadBalancer INetworkLoadBalancer `json:"loadBalancer"`
}

// Define a new network load balancer.
// Experimental.
type NetworkLoadBalancer interface {
	BaseLoadBalancer
	INetworkLoadBalancer
	Env() *awscdk.ResourceEnvironment
	LoadBalancerArn() *string
	LoadBalancerCanonicalHostedZoneId() *string
	LoadBalancerDnsName() *string
	LoadBalancerFullName() *string
	LoadBalancerName() *string
	LoadBalancerSecurityGroups() *[]*string
	Node() awscdk.ConstructNode
	PhysicalName() *string
	Stack() awscdk.Stack
	Vpc() awsec2.IVpc
	AddListener(id *string, props *BaseNetworkListenerProps) NetworkListener
	ApplyRemovalPolicy(policy awscdk.RemovalPolicy)
	GeneratePhysicalName() *string
	GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string
	GetResourceNameAttribute(nameAttr *string) *string
	LogAccessLogs(bucket awss3.IBucket, prefix *string)
	Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricActiveFlowCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricConsumedLCUs(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricNewFlowCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricProcessedBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTcpClientResetCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTcpElbResetCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricTcpTargetResetCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricUnHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RemoveAttribute(key *string)
	SetAttribute(key *string, value *string)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for NetworkLoadBalancer
type jsiiProxy_NetworkLoadBalancer struct {
	jsiiProxy_BaseLoadBalancer
	jsiiProxy_INetworkLoadBalancer
}

func (j *jsiiProxy_NetworkLoadBalancer) Env() *awscdk.ResourceEnvironment {
	var returns *awscdk.ResourceEnvironment
	_jsii_.Get(
		j,
		"env",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) LoadBalancerArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) LoadBalancerCanonicalHostedZoneId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerCanonicalHostedZoneId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) LoadBalancerDnsName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerDnsName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) LoadBalancerFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) LoadBalancerName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) LoadBalancerSecurityGroups() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"loadBalancerSecurityGroups",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) PhysicalName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"physicalName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) Stack() awscdk.Stack {
	var returns awscdk.Stack
	_jsii_.Get(
		j,
		"stack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkLoadBalancer) Vpc() awsec2.IVpc {
	var returns awsec2.IVpc
	_jsii_.Get(
		j,
		"vpc",
		&returns,
	)
	return returns
}


// Experimental.
func NewNetworkLoadBalancer(scope constructs.Construct, id *string, props *NetworkLoadBalancerProps) NetworkLoadBalancer {
	_init_.Initialize()

	j := jsiiProxy_NetworkLoadBalancer{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkLoadBalancer",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewNetworkLoadBalancer_Override(n NetworkLoadBalancer, scope constructs.Construct, id *string, props *NetworkLoadBalancerProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkLoadBalancer",
		[]interface{}{scope, id, props},
		n,
	)
}

// Looks up the network load balancer.
// Experimental.
func NetworkLoadBalancer_FromLookup(scope constructs.Construct, id *string, options *NetworkLoadBalancerLookupOptions) INetworkLoadBalancer {
	_init_.Initialize()

	var returns INetworkLoadBalancer

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkLoadBalancer",
		"fromLookup",
		[]interface{}{scope, id, options},
		&returns,
	)

	return returns
}

// Experimental.
func NetworkLoadBalancer_FromNetworkLoadBalancerAttributes(scope constructs.Construct, id *string, attrs *NetworkLoadBalancerAttributes) INetworkLoadBalancer {
	_init_.Initialize()

	var returns INetworkLoadBalancer

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkLoadBalancer",
		"fromNetworkLoadBalancerAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func NetworkLoadBalancer_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkLoadBalancer",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Check whether the given construct is a Resource.
// Experimental.
func NetworkLoadBalancer_IsResource(construct awscdk.IConstruct) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkLoadBalancer",
		"isResource",
		[]interface{}{construct},
		&returns,
	)

	return returns
}

// Add a listener to this load balancer.
//
// Returns: The newly created listener
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) AddListener(id *string, props *BaseNetworkListenerProps) NetworkListener {
	var returns NetworkListener

	_jsii_.Invoke(
		n,
		"addListener",
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
func (n *jsiiProxy_NetworkLoadBalancer) ApplyRemovalPolicy(policy awscdk.RemovalPolicy) {
	_jsii_.InvokeVoid(
		n,
		"applyRemovalPolicy",
		[]interface{}{policy},
	)
}

// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) GeneratePhysicalName() *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkLoadBalancer) GetResourceArnAttribute(arnAttr *string, arnComponents *awscdk.ArnComponents) *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkLoadBalancer) GetResourceNameAttribute(nameAttr *string) *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"getResourceNameAttribute",
		[]interface{}{nameAttr},
		&returns,
	)

	return returns
}

// Enable access logging for this load balancer.
//
// A region must be specified on the stack containing the load balancer; you cannot enable logging on
// environment-agnostic stacks. See https://docs.aws.amazon.com/cdk/latest/guide/environments.html
//
// This is extending the BaseLoadBalancer.logAccessLogs method to match the bucket permissions described
// at https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-access-logs.html#access-logging-bucket-requirements
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) LogAccessLogs(bucket awss3.IBucket, prefix *string) {
	_jsii_.InvokeVoid(
		n,
		"logAccessLogs",
		[]interface{}{bucket, prefix},
	)
}

// Return the given named metric for this Network Load Balancer.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) Metric(metricName *string, props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metric",
		[]interface{}{metricName, props},
		&returns,
	)

	return returns
}

// The total number of concurrent TCP flows (or connections) from clients to targets.
//
// This metric includes connections in the SYN_SENT and ESTABLISHED states.
// TCP connections are not terminated at the load balancer, so a client
// opening a TCP connection to a target counts as a single flow.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) MetricActiveFlowCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricActiveFlowCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of load balancer capacity units (LCU) used by your load balancer.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) MetricConsumedLCUs(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricConsumedLCUs",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of targets that are considered healthy.
// Deprecated: use ``NetworkTargetGroup.metricHealthyHostCount`` instead
func (n *jsiiProxy_NetworkLoadBalancer) MetricHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricHealthyHostCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of new TCP flows (or connections) established from clients to targets in the time period.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) MetricNewFlowCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricNewFlowCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of bytes processed by the load balancer, including TCP/IP headers.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) MetricProcessedBytes(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricProcessedBytes",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of reset (RST) packets sent from a client to a target.
//
// These resets are generated by the client and forwarded by the load balancer.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) MetricTcpClientResetCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricTcpClientResetCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of reset (RST) packets generated by the load balancer.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) MetricTcpElbResetCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricTcpElbResetCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The total number of reset (RST) packets sent from a target to a client.
//
// These resets are generated by the target and forwarded by the load balancer.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) MetricTcpTargetResetCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricTcpTargetResetCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of targets that are considered unhealthy.
// Deprecated: use ``NetworkTargetGroup.metricUnHealthyHostCount`` instead
func (n *jsiiProxy_NetworkLoadBalancer) MetricUnHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricUnHealthyHostCount",
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
func (n *jsiiProxy_NetworkLoadBalancer) OnPrepare() {
	_jsii_.InvokeVoid(
		n,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
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
func (n *jsiiProxy_NetworkLoadBalancer) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkLoadBalancer) Prepare() {
	_jsii_.InvokeVoid(
		n,
		"prepare",
		nil, // no parameters
	)
}

// Remove an attribute from the load balancer.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) RemoveAttribute(key *string) {
	_jsii_.InvokeVoid(
		n,
		"removeAttribute",
		[]interface{}{key},
	)
}

// Set a non-standard attribute on the load balancer.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/application-load-balancers.html#load-balancer-attributes
//
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) SetAttribute(key *string, value *string) {
	_jsii_.InvokeVoid(
		n,
		"setAttribute",
		[]interface{}{key, value},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (n *jsiiProxy_NetworkLoadBalancer) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkLoadBalancer) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to reference an existing load balancer.
// Experimental.
type NetworkLoadBalancerAttributes struct {
	// ARN of the load balancer.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// The canonical hosted zone ID of this load balancer.
	// Experimental.
	LoadBalancerCanonicalHostedZoneId *string `json:"loadBalancerCanonicalHostedZoneId"`
	// The DNS name of this load balancer.
	// Experimental.
	LoadBalancerDnsName *string `json:"loadBalancerDnsName"`
	// The VPC to associate with the load balancer.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
}

// Options for looking up an NetworkLoadBalancer.
// Experimental.
type NetworkLoadBalancerLookupOptions struct {
	// Find by load balancer's ARN.
	// Experimental.
	LoadBalancerArn *string `json:"loadBalancerArn"`
	// Match load balancer tags.
	// Experimental.
	LoadBalancerTags *map[string]*string `json:"loadBalancerTags"`
}

// Properties for a network load balancer.
// Experimental.
type NetworkLoadBalancerProps struct {
	// The VPC network to place the load balancer in.
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// Indicates whether deletion protection is enabled.
	// Experimental.
	DeletionProtection *bool `json:"deletionProtection"`
	// Whether the load balancer has an internet-routable address.
	// Experimental.
	InternetFacing *bool `json:"internetFacing"`
	// Name of the load balancer.
	// Experimental.
	LoadBalancerName *string `json:"loadBalancerName"`
	// Which subnets place the load balancer in.
	// Experimental.
	VpcSubnets *awsec2.SubnetSelection `json:"vpcSubnets"`
	// Indicates whether cross-zone load balancing is enabled.
	// Experimental.
	CrossZoneEnabled *bool `json:"crossZoneEnabled"`
}

// Define a Network Target Group.
// Experimental.
type NetworkTargetGroup interface {
	TargetGroupBase
	INetworkTargetGroup
	DefaultPort() *float64
	FirstLoadBalancerFullName() *string
	HealthCheck() *HealthCheck
	SetHealthCheck(val *HealthCheck)
	LoadBalancerArns() *string
	LoadBalancerAttached() awscdk.IDependable
	LoadBalancerAttachedDependencies() awscdk.ConcreteDependable
	Node() awscdk.ConstructNode
	TargetGroupArn() *string
	TargetGroupFullName() *string
	TargetGroupLoadBalancerArns() *[]*string
	TargetGroupName() *string
	TargetType() TargetType
	SetTargetType(val TargetType)
	AddLoadBalancerTarget(props *LoadBalancerTargetProps)
	AddTarget(targets ...INetworkLoadBalancerTarget)
	ConfigureHealthCheck(healthCheck *HealthCheck)
	MetricHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	MetricUnHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	RegisterListener(listener INetworkListener)
	SetAttribute(key *string, value *string)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for NetworkTargetGroup
type jsiiProxy_NetworkTargetGroup struct {
	jsiiProxy_TargetGroupBase
	jsiiProxy_INetworkTargetGroup
}

func (j *jsiiProxy_NetworkTargetGroup) DefaultPort() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"defaultPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) FirstLoadBalancerFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"firstLoadBalancerFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) HealthCheck() *HealthCheck {
	var returns *HealthCheck
	_jsii_.Get(
		j,
		"healthCheck",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) LoadBalancerArns() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) LoadBalancerAttached() awscdk.IDependable {
	var returns awscdk.IDependable
	_jsii_.Get(
		j,
		"loadBalancerAttached",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) LoadBalancerAttachedDependencies() awscdk.ConcreteDependable {
	var returns awscdk.ConcreteDependable
	_jsii_.Get(
		j,
		"loadBalancerAttachedDependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) TargetGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) TargetGroupFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) TargetGroupLoadBalancerArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"targetGroupLoadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) TargetGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_NetworkTargetGroup) TargetType() TargetType {
	var returns TargetType
	_jsii_.Get(
		j,
		"targetType",
		&returns,
	)
	return returns
}


// Experimental.
func NewNetworkTargetGroup(scope constructs.Construct, id *string, props *NetworkTargetGroupProps) NetworkTargetGroup {
	_init_.Initialize()

	j := jsiiProxy_NetworkTargetGroup{}

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkTargetGroup",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

// Experimental.
func NewNetworkTargetGroup_Override(n NetworkTargetGroup, scope constructs.Construct, id *string, props *NetworkTargetGroupProps) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.NetworkTargetGroup",
		[]interface{}{scope, id, props},
		n,
	)
}

func (j *jsiiProxy_NetworkTargetGroup) SetHealthCheck(val *HealthCheck) {
	_jsii_.Set(
		j,
		"healthCheck",
		val,
	)
}

func (j *jsiiProxy_NetworkTargetGroup) SetTargetType(val TargetType) {
	_jsii_.Set(
		j,
		"targetType",
		val,
	)
}

// Import an existing target group.
// Experimental.
func NetworkTargetGroup_FromTargetGroupAttributes(scope constructs.Construct, id *string, attrs *TargetGroupAttributes) INetworkTargetGroup {
	_init_.Initialize()

	var returns INetworkTargetGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkTargetGroup",
		"fromTargetGroupAttributes",
		[]interface{}{scope, id, attrs},
		&returns,
	)

	return returns
}

// Import an existing listener.
// Deprecated: Use `fromTargetGroupAttributes` instead
func NetworkTargetGroup_Import(scope constructs.Construct, id *string, props *TargetGroupImportProps) INetworkTargetGroup {
	_init_.Initialize()

	var returns INetworkTargetGroup

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkTargetGroup",
		"import",
		[]interface{}{scope, id, props},
		&returns,
	)

	return returns
}

// Return whether the given object is a Construct.
// Experimental.
func NetworkTargetGroup_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.NetworkTargetGroup",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Register the given load balancing target as part of this group.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) AddLoadBalancerTarget(props *LoadBalancerTargetProps) {
	_jsii_.InvokeVoid(
		n,
		"addLoadBalancerTarget",
		[]interface{}{props},
	)
}

// Add a load balancing target to this target group.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) AddTarget(targets ...INetworkLoadBalancerTarget) {
	args := []interface{}{}
	for _, a := range targets {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		n,
		"addTarget",
		args,
	)
}

// Set/replace the target group's health check.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) ConfigureHealthCheck(healthCheck *HealthCheck) {
	_jsii_.InvokeVoid(
		n,
		"configureHealthCheck",
		[]interface{}{healthCheck},
	)
}

// The number of targets that are considered healthy.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) MetricHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricHealthyHostCount",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// The number of targets that are considered unhealthy.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) MetricUnHealthyHostCount(props *awscloudwatch.MetricOptions) awscloudwatch.Metric {
	var returns awscloudwatch.Metric

	_jsii_.Invoke(
		n,
		"metricUnHealthyHostCount",
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
func (n *jsiiProxy_NetworkTargetGroup) OnPrepare() {
	_jsii_.InvokeVoid(
		n,
		"onPrepare",
		nil, // no parameters
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) OnSynthesize(session constructs.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
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
func (n *jsiiProxy_NetworkTargetGroup) OnValidate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkTargetGroup) Prepare() {
	_jsii_.InvokeVoid(
		n,
		"prepare",
		nil, // no parameters
	)
}

// Register a listener that is load balancing to this target group.
//
// Don't call this directly. It will be called by listeners.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) RegisterListener(listener INetworkListener) {
	_jsii_.InvokeVoid(
		n,
		"registerListener",
		[]interface{}{listener},
	)
}

// Set a non-standard attribute on the target group.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-target-groups.html#target-group-attributes
//
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) SetAttribute(key *string, value *string) {
	_jsii_.InvokeVoid(
		n,
		"setAttribute",
		[]interface{}{key, value},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		n,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (n *jsiiProxy_NetworkTargetGroup) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		n,
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
func (n *jsiiProxy_NetworkTargetGroup) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		n,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties for a new Network Target Group.
// Experimental.
type NetworkTargetGroupProps struct {
	// The amount of time for Elastic Load Balancing to wait before deregistering a target.
	//
	// The range is 0-3600 seconds.
	// Experimental.
	DeregistrationDelay awscdk.Duration `json:"deregistrationDelay"`
	// Health check configuration.
	// Experimental.
	HealthCheck *HealthCheck `json:"healthCheck"`
	// The name of the target group.
	//
	// This name must be unique per region per account, can have a maximum of
	// 32 characters, must contain only alphanumeric characters or hyphens, and
	// must not begin or end with a hyphen.
	// Experimental.
	TargetGroupName *string `json:"targetGroupName"`
	// The type of targets registered to this TargetGroup, either IP or Instance.
	//
	// All targets registered into the group must be of this type. If you
	// register targets to the TargetGroup in the CDK app, the TargetType is
	// determined automatically.
	// Experimental.
	TargetType TargetType `json:"targetType"`
	// The virtual private cloud (VPC).
	//
	// only if `TargetType` is `Ip` or `InstanceId`
	// Experimental.
	Vpc awsec2.IVpc `json:"vpc"`
	// The port on which the listener listens for requests.
	// Experimental.
	Port *float64 `json:"port"`
	// Indicates whether client IP preservation is enabled.
	// Experimental.
	PreserveClientIp *bool `json:"preserveClientIp"`
	// Protocol for target group, expects TCP, TLS, UDP, or TCP_UDP.
	// Experimental.
	Protocol Protocol `json:"protocol"`
	// Indicates whether Proxy Protocol version 2 is enabled.
	// Experimental.
	ProxyProtocolV2 *bool `json:"proxyProtocolV2"`
	// The targets to add to this target group.
	//
	// Can be `Instance`, `IPAddress`, or any self-registering load balancing
	// target. If you use either `Instance` or `IPAddress` as targets, all
	// target must be of the same type.
	// Experimental.
	Targets *[]INetworkLoadBalancerTarget `json:"targets"`
}

// A Target Group and weight combination.
// Experimental.
type NetworkWeightedTargetGroup struct {
	// The target group.
	// Experimental.
	TargetGroup INetworkTargetGroup `json:"targetGroup"`
	// The target group's weight.
	//
	// Range is [0..1000).
	// Experimental.
	Weight *float64 `json:"weight"`
}

// Backend protocol for network load balancers and health checks.
// Experimental.
type Protocol string

const (
	Protocol_HTTP Protocol = "HTTP"
	Protocol_HTTPS Protocol = "HTTPS"
	Protocol_TCP Protocol = "TCP"
	Protocol_TLS Protocol = "TLS"
	Protocol_UDP Protocol = "UDP"
	Protocol_TCP_UDP Protocol = "TCP_UDP"
)

// Properties for the key/value pair of the query string.
// Experimental.
type QueryStringCondition struct {
	// The query string value for the condition.
	// Experimental.
	Value *string `json:"value"`
	// The query string key for the condition.
	// Experimental.
	Key *string `json:"key"`
}

// Options for `ListenerAction.redirect()`.
//
// A URI consists of the following components:
// protocol://hostname:port/path?query. You must modify at least one of the
// following components to avoid a redirect loop: protocol, hostname, port, or
// path. Any components that you do not modify retain their original values.
//
// You can reuse URI components using the following reserved keywords:
//
// - `#{protocol}`
// - `#{host}`
// - `#{port}`
// - `#{path}` (the leading "/" is removed)
// - `#{query}`
//
// For example, you can change the path to "/new/#{path}", the hostname to
// "example.#{host}", or the query to "#{query}&value=xyz".
// Experimental.
type RedirectOptions struct {
	// The hostname.
	//
	// This component is not percent-encoded. The hostname can contain #{host}.
	// Experimental.
	Host *string `json:"host"`
	// The absolute path, starting with the leading "/".
	//
	// This component is not percent-encoded. The path can contain #{host}, #{path}, and #{port}.
	// Experimental.
	Path *string `json:"path"`
	// The HTTP redirect code.
	//
	// The redirect is either permanent (HTTP 301) or temporary (HTTP 302).
	// Experimental.
	Permanent *bool `json:"permanent"`
	// The port.
	//
	// You can specify a value from 1 to 65535 or #{port}.
	// Experimental.
	Port *string `json:"port"`
	// The protocol.
	//
	// You can specify HTTP, HTTPS, or #{protocol}. You can redirect HTTP to HTTP, HTTP to HTTPS, and HTTPS to HTTPS. You cannot redirect HTTPS to HTTP.
	// Experimental.
	Protocol *string `json:"protocol"`
	// The query parameters, URL-encoded when necessary, but not percent-encoded.
	//
	// Do not include the leading "?", as it is automatically added. You can specify any of the reserved keywords.
	// Experimental.
	Query *string `json:"query"`
}

// A redirect response.
// Deprecated: superceded by `ListenerAction.redirect()`.
type RedirectResponse struct {
	// The HTTP redirect code (HTTP_301 or HTTP_302).
	// Deprecated: superceded by `ListenerAction.redirect()`.
	StatusCode *string `json:"statusCode"`
	// The hostname.
	//
	// This component is not percent-encoded. The hostname can contain #{host}.
	// Deprecated: superceded by `ListenerAction.redirect()`.
	Host *string `json:"host"`
	// The absolute path, starting with the leading "/".
	//
	// This component is not percent-encoded.
	// The path can contain #{host}, #{path}, and #{port}.
	// Deprecated: superceded by `ListenerAction.redirect()`.
	Path *string `json:"path"`
	// The port.
	//
	// You can specify a value from 1 to 65535 or #{port}.
	// Deprecated: superceded by `ListenerAction.redirect()`.
	Port *string `json:"port"`
	// The protocol.
	//
	// You can specify HTTP, HTTPS, or #{protocol}. You can redirect HTTP to HTTP,
	// HTTP to HTTPS, and HTTPS to HTTPS. You cannot redirect HTTPS to HTTP.
	// Deprecated: superceded by `ListenerAction.redirect()`.
	Protocol *string `json:"protocol"`
	// The query parameters, URL-encoded when necessary, but not percent-encoded.
	//
	// Do not include the leading "?", as it is automatically added.
	// You can specify any of the reserved keywords.
	// Deprecated: superceded by `ListenerAction.redirect()`.
	Query *string `json:"query"`
}

// Elastic Load Balancing provides the following security policies for Application Load Balancers.
//
// We recommend the Recommended policy for general use. You can
// use the ForwardSecrecy policy if you require Forward Secrecy
// (FS).
//
// You can use one of the TLS policies to meet compliance and security
// standards that require disabling certain TLS protocol versions, or to
// support legacy clients that require deprecated ciphers.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-https-listener.html
//
// Experimental.
type SslPolicy string

const (
	SslPolicy_RECOMMENDED SslPolicy = "RECOMMENDED"
	SslPolicy_FORWARD_SECRECY_TLS12_RES_GCM SslPolicy = "FORWARD_SECRECY_TLS12_RES_GCM"
	SslPolicy_FORWARD_SECRECY_TLS12_RES SslPolicy = "FORWARD_SECRECY_TLS12_RES"
	SslPolicy_FORWARD_SECRECY_TLS12 SslPolicy = "FORWARD_SECRECY_TLS12"
	SslPolicy_FORWARD_SECRECY_TLS11 SslPolicy = "FORWARD_SECRECY_TLS11"
	SslPolicy_FORWARD_SECRECY SslPolicy = "FORWARD_SECRECY"
	SslPolicy_TLS12 SslPolicy = "TLS12"
	SslPolicy_TLS12_EXT SslPolicy = "TLS12_EXT"
	SslPolicy_TLS11 SslPolicy = "TLS11"
	SslPolicy_LEGACY SslPolicy = "LEGACY"
)

// Properties to reference an existing target group.
// Experimental.
type TargetGroupAttributes struct {
	// ARN of the target group.
	// Experimental.
	TargetGroupArn *string `json:"targetGroupArn"`
	// Port target group is listening on.
	// Deprecated: - This property is unused and the wrong type. No need to use it.
	DefaultPort *string `json:"defaultPort"`
	// A Token representing the list of ARNs for the load balancer routing to this target group.
	// Experimental.
	LoadBalancerArns *string `json:"loadBalancerArns"`
}

// Define the target of a load balancer.
// Experimental.
type TargetGroupBase interface {
	awscdk.Construct
	ITargetGroup
	DefaultPort() *float64
	FirstLoadBalancerFullName() *string
	HealthCheck() *HealthCheck
	SetHealthCheck(val *HealthCheck)
	LoadBalancerArns() *string
	LoadBalancerAttached() awscdk.IDependable
	LoadBalancerAttachedDependencies() awscdk.ConcreteDependable
	Node() awscdk.ConstructNode
	TargetGroupArn() *string
	TargetGroupFullName() *string
	TargetGroupLoadBalancerArns() *[]*string
	TargetGroupName() *string
	TargetType() TargetType
	SetTargetType(val TargetType)
	AddLoadBalancerTarget(props *LoadBalancerTargetProps)
	ConfigureHealthCheck(healthCheck *HealthCheck)
	OnPrepare()
	OnSynthesize(session constructs.ISynthesisSession)
	OnValidate() *[]*string
	Prepare()
	SetAttribute(key *string, value *string)
	Synthesize(session awscdk.ISynthesisSession)
	ToString() *string
	Validate() *[]*string
}

// The jsii proxy struct for TargetGroupBase
type jsiiProxy_TargetGroupBase struct {
	internal.Type__awscdkConstruct
	jsiiProxy_ITargetGroup
}

func (j *jsiiProxy_TargetGroupBase) DefaultPort() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"defaultPort",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) FirstLoadBalancerFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"firstLoadBalancerFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) HealthCheck() *HealthCheck {
	var returns *HealthCheck
	_jsii_.Get(
		j,
		"healthCheck",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) LoadBalancerArns() *string {
	var returns *string
	_jsii_.Get(
		j,
		"loadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) LoadBalancerAttached() awscdk.IDependable {
	var returns awscdk.IDependable
	_jsii_.Get(
		j,
		"loadBalancerAttached",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) LoadBalancerAttachedDependencies() awscdk.ConcreteDependable {
	var returns awscdk.ConcreteDependable
	_jsii_.Get(
		j,
		"loadBalancerAttachedDependencies",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) Node() awscdk.ConstructNode {
	var returns awscdk.ConstructNode
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) TargetGroupArn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupArn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) TargetGroupFullName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupFullName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) TargetGroupLoadBalancerArns() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"targetGroupLoadBalancerArns",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) TargetGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_TargetGroupBase) TargetType() TargetType {
	var returns TargetType
	_jsii_.Get(
		j,
		"targetType",
		&returns,
	)
	return returns
}


// Experimental.
func NewTargetGroupBase_Override(t TargetGroupBase, scope constructs.Construct, id *string, baseProps *BaseTargetGroupProps, additionalProps interface{}) {
	_init_.Initialize()

	_jsii_.Create(
		"monocdk.aws_elasticloadbalancingv2.TargetGroupBase",
		[]interface{}{scope, id, baseProps, additionalProps},
		t,
	)
}

func (j *jsiiProxy_TargetGroupBase) SetHealthCheck(val *HealthCheck) {
	_jsii_.Set(
		j,
		"healthCheck",
		val,
	)
}

func (j *jsiiProxy_TargetGroupBase) SetTargetType(val TargetType) {
	_jsii_.Set(
		j,
		"targetType",
		val,
	)
}

// Return whether the given object is a Construct.
// Experimental.
func TargetGroupBase_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	var returns *bool

	_jsii_.StaticInvoke(
		"monocdk.aws_elasticloadbalancingv2.TargetGroupBase",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Register the given load balancing target as part of this group.
// Experimental.
func (t *jsiiProxy_TargetGroupBase) AddLoadBalancerTarget(props *LoadBalancerTargetProps) {
	_jsii_.InvokeVoid(
		t,
		"addLoadBalancerTarget",
		[]interface{}{props},
	)
}

// Set/replace the target group's health check.
// Experimental.
func (t *jsiiProxy_TargetGroupBase) ConfigureHealthCheck(healthCheck *HealthCheck) {
	_jsii_.InvokeVoid(
		t,
		"configureHealthCheck",
		[]interface{}{healthCheck},
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
func (t *jsiiProxy_TargetGroupBase) OnPrepare() {
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
func (t *jsiiProxy_TargetGroupBase) OnSynthesize(session constructs.ISynthesisSession) {
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
func (t *jsiiProxy_TargetGroupBase) OnValidate() *[]*string {
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
func (t *jsiiProxy_TargetGroupBase) Prepare() {
	_jsii_.InvokeVoid(
		t,
		"prepare",
		nil, // no parameters
	)
}

// Set a non-standard attribute on the target group.
// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-target-groups.html#target-group-attributes
//
// Experimental.
func (t *jsiiProxy_TargetGroupBase) SetAttribute(key *string, value *string) {
	_jsii_.InvokeVoid(
		t,
		"setAttribute",
		[]interface{}{key, value},
	)
}

// Allows this construct to emit artifacts into the cloud assembly during synthesis.
//
// This method is usually implemented by framework-level constructs such as `Stack` and `Asset`
// as they participate in synthesizing the cloud assembly.
// Experimental.
func (t *jsiiProxy_TargetGroupBase) Synthesize(session awscdk.ISynthesisSession) {
	_jsii_.InvokeVoid(
		t,
		"synthesize",
		[]interface{}{session},
	)
}

// Returns a string representation of this construct.
// Experimental.
func (t *jsiiProxy_TargetGroupBase) ToString() *string {
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
// Experimental.
func (t *jsiiProxy_TargetGroupBase) Validate() *[]*string {
	var returns *[]*string

	_jsii_.Invoke(
		t,
		"validate",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Properties to reference an existing target group.
// Deprecated: Use TargetGroupAttributes instead
type TargetGroupImportProps struct {
	// ARN of the target group.
	// Deprecated: Use TargetGroupAttributes instead
	TargetGroupArn *string `json:"targetGroupArn"`
	// Port target group is listening on.
	// Deprecated: - This property is unused and the wrong type. No need to use it.
	DefaultPort *string `json:"defaultPort"`
	// A Token representing the list of ARNs for the load balancer routing to this target group.
	// Deprecated: Use TargetGroupAttributes instead
	LoadBalancerArns *string `json:"loadBalancerArns"`
}

// How to interpret the load balancing target identifiers.
// Experimental.
type TargetType string

const (
	TargetType_INSTANCE TargetType = "INSTANCE"
	TargetType_IP TargetType = "IP"
	TargetType_LAMBDA TargetType = "LAMBDA"
)

// What to do with unauthenticated requests.
// Experimental.
type UnauthenticatedAction string

const (
	UnauthenticatedAction_DENY UnauthenticatedAction = "DENY"
	UnauthenticatedAction_ALLOW UnauthenticatedAction = "ALLOW"
	UnauthenticatedAction_AUTHENTICATE UnauthenticatedAction = "AUTHENTICATE"
)

// A Target Group and weight combination.
// Experimental.
type WeightedTargetGroup struct {
	// The target group.
	// Experimental.
	TargetGroup IApplicationTargetGroup `json:"targetGroup"`
	// The target group's weight.
	//
	// Range is [0..1000).
	// Experimental.
	Weight *float64 `json:"weight"`
}


/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

import (
	"time"
)

const (
	// DefaultAPIServerPort defines the API server port when defining a Load Balancer.
	DefaultAPIServerPort = 6443
	// DefaultAPIServerPortString defines the API server port as a string for convenience.
	DefaultAPIServerPortString = "6443"
	// DefaultAPIServerHealthCheckPath the API server health check path.
	DefaultAPIServerHealthCheckPath = "/readyz"
	// DefaultAPIServerHealthCheckIntervalSec the API server health check interval in seconds.
	DefaultAPIServerHealthCheckIntervalSec = 10
	// DefaultAPIServerHealthCheckTimeoutSec the API server health check timeout in seconds.
	DefaultAPIServerHealthCheckTimeoutSec = 5
	// DefaultAPIServerHealthThresholdCount the API server health check threshold count.
	DefaultAPIServerHealthThresholdCount = 5
	// DefaultAPIServerUnhealthThresholdCount the API server unhealthy check threshold count.
	DefaultAPIServerUnhealthThresholdCount = 3
)

// LoadBalancerType defines the type of load balancer to use.
type LoadBalancerType string

var (
	// LoadBalancerTypeClassic is the classic ELB type.
	LoadBalancerTypeClassic = LoadBalancerType("classic")
	// LoadBalancerTypeELB is the ELB type.
	LoadBalancerTypeELB = LoadBalancerType("elb")
	// LoadBalancerTypeALB is the ALB type.
	LoadBalancerTypeALB = LoadBalancerType("alb")
	// LoadBalancerTypeNLB is the NLB type.
	LoadBalancerTypeNLB = LoadBalancerType("nlb")
	// LoadBalancerTypeDisabled disables the load balancer.
	LoadBalancerTypeDisabled = LoadBalancerType("disabled")
)

// AWSLoadBalancerSpec defines the desired state of an AWS load balancer.
type AWSLoadBalancerSpec struct {
	// Name sets the name of the classic ELB load balancer. As per AWS, the name must be unique
	// within your set of load balancers for the region, must have a maximum of 32 characters, must
	// contain only alphanumeric characters or hyphens, and cannot begin or end with a hyphen. Once
	// set, the value cannot be changed.
	// +kubebuilder:validation:MaxLength:=32
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9]([A-Za-z0-9]{0,31}|[-A-Za-z0-9]{0,30}[A-Za-z0-9])$`
	// +optional
	Name *string `json:"name,omitempty"`

	// Scheme sets the scheme of the load balancer (defaults to internet-facing)
	// +kubebuilder:default=internet-facing
	// +kubebuilder:validation:Enum=internet-facing;internal
	// +optional
	Scheme *ELBScheme `json:"scheme,omitempty"`

	// CrossZoneLoadBalancing enables the classic ELB cross availability zone balancing.
	//
	// With cross-zone load balancing, each load balancer node for your Classic Load Balancer
	// distributes requests evenly across the registered instances in all enabled Availability Zones.
	// If cross-zone load balancing is disabled, each load balancer node distributes requests evenly across
	// the registered instances in its Availability Zone only.
	//
	// Defaults to false.
	// +optional
	CrossZoneLoadBalancing bool `json:"crossZoneLoadBalancing"`

	// Subnets sets the subnets that should be applied to the control plane load balancer (defaults to discovered subnets for managed VPCs or an empty set for unmanaged VPCs)
	// +optional
	Subnets []string `json:"subnets,omitempty"`

	// HealthCheckProtocol sets the protocol type for ELB health check target
	// default value is ELBProtocolSSL
	// +kubebuilder:validation:Enum=TCP;SSL;HTTP;HTTPS;TLS;UDP
	// +optional
	HealthCheckProtocol *ELBProtocol `json:"healthCheckProtocol,omitempty"`

	// HealthCheck sets custom health check configuration to the API target group.
	// +optional
	HealthCheck *TargetGroupHealthCheckAPISpec `json:"healthCheck,omitempty"`

	// AdditionalSecurityGroups sets the security groups used by the load balancer. Expected to be security group IDs
	// This is optional - if not provided new security groups will be created for the load balancer
	// +optional
	AdditionalSecurityGroups []string `json:"additionalSecurityGroups,omitempty"`

	// AdditionalListeners sets the additional listeners for the control plane load balancer.
	// This is only applicable to Network Load Balancer (NLB) types for the time being.
	// +listType=map
	// +listMapKey=port
	// +optional
	AdditionalListeners []AdditionalListenerSpec `json:"additionalListeners,omitempty"`

	// IngressRules sets the ingress rules for the control plane load balancer.
	// +optional
	IngressRules []IngressRule `json:"ingressRules,omitempty"`

	// LoadBalancerType sets the type for a load balancer. The default type is classic.
	// +kubebuilder:default=classic
	// +kubebuilder:validation:Enum:=classic;elb;alb;nlb;disabled
	LoadBalancerType LoadBalancerType `json:"loadBalancerType,omitempty"`

	// DisableHostsRewrite disabled the hair pinning issue solution that adds the NLB's address as 127.0.0.1 to the hosts
	// file of each instance. This is by default, false.
	DisableHostsRewrite bool `json:"disableHostsRewrite,omitempty"`

	// PreserveClientIP lets the user control if preservation of client ips must be retained or not.
	// If this is enabled 6443 will be opened to 0.0.0.0/0.
	PreserveClientIP bool `json:"preserveClientIP,omitempty"`
}

// AdditionalListenerSpec defines the desired state of an
// additional listener on an AWS load balancer.
type AdditionalListenerSpec struct {
	// Port sets the port for the additional listener.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int64 `json:"port"`

	// Protocol sets the protocol for the additional listener.
	// Currently only TCP is supported.
	// +kubebuilder:validation:Enum=TCP
	// +kubebuilder:default=TCP
	Protocol ELBProtocol `json:"protocol,omitempty"`

	// HealthCheck sets the optional custom health check configuration to the API target group.
	// +optional
	HealthCheck *TargetGroupHealthCheckAdditionalSpec `json:"healthCheck,omitempty"`
}

// ELBScheme defines the scheme of a load balancer.
type ELBScheme string

var (
	// ELBSchemeInternetFacing defines an internet-facing, publicly
	// accessible AWS ELB scheme.
	ELBSchemeInternetFacing = ELBScheme("internet-facing")

	// ELBSchemeInternal defines an internal-only facing
	// load balancer internal to an ELB.
	ELBSchemeInternal = ELBScheme("internal")
)

func (e ELBScheme) String() string {
	return string(e)
}

// Equals returns true if two ELBScheme are equal.
func (e ELBScheme) Equals(other *ELBScheme) bool {
	if other == nil {
		return false
	}

	return e == *other
}

// ELBProtocol defines listener protocols for a load balancer.
type ELBProtocol string

func (e ELBProtocol) String() string {
	return string(e)
}

var (
	// ELBProtocolTCP defines the ELB API string representing the TCP protocol.
	ELBProtocolTCP = ELBProtocol("TCP")
	// ELBProtocolSSL defines the ELB API string representing the TLS protocol.
	ELBProtocolSSL = ELBProtocol("SSL")
	// ELBProtocolHTTP defines the ELB API string representing the HTTP protocol at L7.
	ELBProtocolHTTP = ELBProtocol("HTTP")
	// ELBProtocolHTTPS defines the ELB API string representing the HTTP protocol at L7.
	ELBProtocolHTTPS = ELBProtocol("HTTPS")
	// ELBProtocolTLS defines the NLB API string representing the TLS protocol.
	ELBProtocolTLS = ELBProtocol("TLS")
	// ELBProtocolUDP defines the NLB API string representing the UDP protocol.
	ELBProtocolUDP = ELBProtocol("UDP")
)

// LoadBalancerAttribute defines a set of attributes for a V2 load balancer.
type LoadBalancerAttribute string

var (
	// LoadBalancerAttributeEnableLoadBalancingCrossZone defines the attribute key for enabling load balancing cross zone.
	LoadBalancerAttributeEnableLoadBalancingCrossZone = "load_balancing.cross_zone.enabled"
	// LoadBalancerAttributeIdleTimeTimeoutSeconds defines the attribute key for idle timeout.
	LoadBalancerAttributeIdleTimeTimeoutSeconds = "idle_timeout.timeout_seconds"
	// LoadBalancerAttributeIdleTimeDefaultTimeoutSecondsInSeconds defines the default idle timeout in seconds.
	LoadBalancerAttributeIdleTimeDefaultTimeoutSecondsInSeconds = "60"
)

// Listener defines an AWS network load balancer listener.
type Listener struct {
	Protocol    ELBProtocol     `json:"protocol"`
	Port        int64           `json:"port"`
	TargetGroup TargetGroupSpec `json:"targetGroup"`
}

// LoadBalancer defines an AWS load balancer.
type LoadBalancer struct {
	// ARN of the load balancer. Unlike the ClassicLB, ARN is used mostly
	// to define and get it.
	ARN string `json:"arn,omitempty"`
	// The name of the load balancer. It must be unique within the set of load balancers
	// defined in the region. It also serves as identifier.
	// +optional
	Name string `json:"name,omitempty"`

	// DNSName is the dns name of the load balancer.
	DNSName string `json:"dnsName,omitempty"`

	// Scheme is the load balancer scheme, either internet-facing or private.
	Scheme ELBScheme `json:"scheme,omitempty"`

	// AvailabilityZones is an array of availability zones in the VPC attached to the load balancer.
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// SubnetIDs is an array of subnets in the VPC attached to the load balancer.
	SubnetIDs []string `json:"subnetIds,omitempty"`

	// SecurityGroupIDs is an array of security groups assigned to the load balancer.
	SecurityGroupIDs []string `json:"securityGroupIds,omitempty"`

	// ClassicELBListeners is an array of classic elb listeners associated with the load balancer. There must be at least one.
	ClassicELBListeners []ClassicELBListener `json:"listeners,omitempty"`

	// HealthCheck is the classic elb health check associated with the load balancer.
	HealthCheck *ClassicELBHealthCheck `json:"healthChecks,omitempty"`

	// ClassicElbAttributes defines extra attributes associated with the load balancer.
	ClassicElbAttributes ClassicELBAttributes `json:"attributes,omitempty"`

	// Tags is a map of tags associated with the load balancer.
	Tags map[string]string `json:"tags,omitempty"`

	// ELBListeners is an array of listeners associated with the load balancer. There must be at least one.
	ELBListeners []Listener `json:"elbListeners,omitempty"`

	// ELBAttributes defines extra attributes associated with v2 load balancers.
	ELBAttributes map[string]*string `json:"elbAttributes,omitempty"`

	// LoadBalancerType sets the type for a load balancer. The default type is classic.
	// +kubebuilder:validation:Enum:=classic;elb;alb;nlb
	LoadBalancerType LoadBalancerType `json:"loadBalancerType,omitempty"`
}

// IsUnmanaged returns true if the Classic ELB is unmanaged.
func (b *LoadBalancer) IsUnmanaged(clusterName string) bool {
	return b.Name != "" && !Tags(b.Tags).HasOwned(clusterName)
}

// IsManaged returns true if Classic ELB is managed.
func (b *LoadBalancer) IsManaged(clusterName string) bool {
	return !b.IsUnmanaged(clusterName)
}

// ClassicELBAttributes defines extra attributes associated with a classic load balancer.
type ClassicELBAttributes struct {
	// IdleTimeout is time that the connection is allowed to be idle (no data
	// has been sent over the connection) before it is closed by the load balancer.
	IdleTimeout time.Duration `json:"idleTimeout,omitempty"`

	// CrossZoneLoadBalancing enables the classic load balancer load balancing.
	// +optional
	CrossZoneLoadBalancing bool `json:"crossZoneLoadBalancing,omitempty"`
}

// ClassicELBListener defines an AWS classic load balancer listener.
type ClassicELBListener struct {
	Protocol         ELBProtocol `json:"protocol"`
	Port             int64       `json:"port"`
	InstanceProtocol ELBProtocol `json:"instanceProtocol"`
	InstancePort     int64       `json:"instancePort"`
}

// ClassicELBHealthCheck defines an AWS classic load balancer health check.
type ClassicELBHealthCheck struct {
	Target             string        `json:"target"`
	Interval           time.Duration `json:"interval"`
	Timeout            time.Duration `json:"timeout"`
	HealthyThreshold   int64         `json:"healthyThreshold"`
	UnhealthyThreshold int64         `json:"unhealthyThreshold"`
}

// TargetGroupAttribute defines attribute key values for V2 Load Balancer Attributes.
type TargetGroupAttribute string

var (
	// TargetGroupAttributeEnablePreserveClientIP defines the attribute key for enabling preserve client IP.
	TargetGroupAttributeEnablePreserveClientIP = "preserve_client_ip.enabled"
)

// TargetGroupHealthCheck defines health check settings for the target group.
type TargetGroupHealthCheck struct {
	Protocol                *string `json:"protocol,omitempty"`
	Path                    *string `json:"path,omitempty"`
	Port                    *string `json:"port,omitempty"`
	IntervalSeconds         *int64  `json:"intervalSeconds,omitempty"`
	TimeoutSeconds          *int64  `json:"timeoutSeconds,omitempty"`
	ThresholdCount          *int64  `json:"thresholdCount,omitempty"`
	UnhealthyThresholdCount *int64  `json:"unhealthyThresholdCount,omitempty"`
}

// TargetGroupHealthCheckAPISpec defines the optional health check settings for the API target group.
type TargetGroupHealthCheckAPISpec struct {
	// The approximate amount of time, in seconds, between health checks of an individual
	// target.
	// +kubebuilder:validation:Minimum=5
	// +kubebuilder:validation:Maximum=300
	// +optional
	IntervalSeconds *int64 `json:"intervalSeconds,omitempty"`

	// The amount of time, in seconds, during which no response from a target means
	// a failed health check.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=120
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// The number of consecutive health check successes required before considering
	// a target healthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	ThresholdCount *int64 `json:"thresholdCount,omitempty"`

	// The number of consecutive health check failures required before considering
	// a target unhealthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	UnhealthyThresholdCount *int64 `json:"unhealthyThresholdCount,omitempty"`
}

// TargetGroupHealthCheckAdditionalSpec defines the optional health check settings for the additional target groups.
type TargetGroupHealthCheckAdditionalSpec struct {
	// The protocol to use to health check connect with the target. When not specified the Protocol
	// will be the same of the listener.
	// +kubebuilder:validation:Enum=TCP;HTTP;HTTPS
	// +optional
	Protocol *string `json:"protocol,omitempty"`

	// The port the load balancer uses when performing health checks for additional target groups. When
	// not specified this value will be set for the same of listener port.
	// +optional
	Port *string `json:"port,omitempty"`

	// The destination for health checks on the targets when using the protocol HTTP or HTTPS,
	// otherwise the path will be ignored.
	// +optional
	Path *string `json:"path,omitempty"`
	// The approximate amount of time, in seconds, between health checks of an individual
	// target.
	// +kubebuilder:validation:Minimum=5
	// +kubebuilder:validation:Maximum=300
	// +optional
	IntervalSeconds *int64 `json:"intervalSeconds,omitempty"`

	// The amount of time, in seconds, during which no response from a target means
	// a failed health check.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=120
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// The number of consecutive health check successes required before considering
	// a target healthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	ThresholdCount *int64 `json:"thresholdCount,omitempty"`

	// The number of consecutive health check failures required before considering
	// a target unhealthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	UnhealthyThresholdCount *int64 `json:"unhealthyThresholdCount,omitempty"`
}

// TargetGroupSpec specifies target group settings for a given listener.
// This is created first, and the ARN is then passed to the listener.
type TargetGroupSpec struct {
	// Name of the TargetGroup. Must be unique over the same group of listeners.
	// +kubebuilder:validation:MaxLength=32
	Name string `json:"name"`
	// Port is the exposed port
	Port int64 `json:"port"`
	// +kubebuilder:validation:Enum=tcp;tls;udp;TCP;TLS;UDP
	Protocol ELBProtocol `json:"protocol"`
	VpcID    string      `json:"vpcId"`
	// HealthCheck is the elb health check associated with the load balancer.
	HealthCheck *TargetGroupHealthCheck `json:"targetGroupHealthCheck,omitempty"`
}

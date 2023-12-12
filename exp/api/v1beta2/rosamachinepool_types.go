package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ROSAMachinePoolReadyCondition condition reports on the successful reconciliation of ROSAMachinePool.
	ROSAMachinePoolReadyCondition clusterv1.ConditionType = "ROSAMachinePoolReady"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosamachinepools,scope=Namespaced,categories=cluster-api,shortName=rosamp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="MachinePool ready status"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="Number of replicas"

// ROSAMachinePool is a representation of a node pool in a cluster.
type ROSAMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ROSAMachinePoolSpec   `json:"spec,omitempty"`
	Status ROSAMachinePoolStatus `json:"status,omitempty"`
}
type ROSAMachinePoolSpec struct {
	AWS         *ROSAAWS                    `json:"aws,omitempty"`
	Autoscaling *ROSAMachinePoolAutoscaling `json:"autoscaling,omitempty"`
	AutoRepair  bool                        `json:"autoRepair,omitempty"`

	// TODO(alberto): Enable and propagate this API input.
	// Labels           map[string]string            `json:"labels,omitempty"`
	// Taints           []*Taint                     `json:"taints,omitempty"`
	// TuningConfigs    []string                     `json:"tuningConfigs,omitempty"`
	// Version          *Version                     `json:"version,omitempty"`
}

type ROSAAWS struct {
	InstanceType     string            `json:"instanceType,omitempty"`
	AvailabilityZone string            `json:"availabilityZone,omitempty"`
	Subnet           string            `json:"subnet,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

// ROSAMachinePoolAutoscaling in a node pool.
type ROSAMachinePoolAutoscaling struct {
	MaxReplica int `json:"maxReplica,omitempty"`
	MinReplica int `json:"minReplica,omitempty"`
}

// ROSAMachinePoolStatus in a node pool.
type ROSAMachinePoolStatus struct {
	// Ready denotes that the AWSManagedMachinePool nodegroup has joined
	// the cluster
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`

	// Conditions defines current service state of the managed machine pool
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// ID is the ID given by ROSA.
	ID string `json:"id,omitempty"`
}

// +kubebuilder:object:root=true

// ROSAMachinePoolList contains a list of AWSManagedMachinePools.
type ROSAMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSAMachinePool `json:"items"`
}

// GetConditions returns the observations of the operational state of the AWSManagedMachinePool resource.
func (r *ROSAMachinePool) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the AWSManagedMachinePool to the predescribed clusterv1.Conditions.
func (r *ROSAMachinePool) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSAMachinePool{}, &ROSAMachinePoolList{})
}

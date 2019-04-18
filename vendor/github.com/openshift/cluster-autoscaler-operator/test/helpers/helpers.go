package helpers

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
)

// TestDeployment wraps the appsv1.Deployment type to add helper methods.
type TestDeployment struct {
	appsv1.Deployment
}

// NewTestDeployment returns a new TestDeployment wrapping the given
// appsv1.Deployment object.
func NewTestDeployment(dep *appsv1.Deployment) *TestDeployment {
	return &TestDeployment{Deployment: *dep}
}

// Copy returns a pointer to a new TestDeploymentdeep with a deep copy of the
// originally wrapped appsv1.Deployment object.
func (d *TestDeployment) Copy() *TestDeployment {
	newDeployment := &appsv1.Deployment{}
	d.Deployment.DeepCopyInto(newDeployment)

	return NewTestDeployment(newDeployment)
}

// WithAvailableReplicas returns a copy of the object with the AvailableReplicas
// field set to the given value.
func (d *TestDeployment) WithAvailableReplicas(n int32) *TestDeployment {
	newDeployment := d.Copy()
	newDeployment.Status.AvailableReplicas = n

	return newDeployment
}

// WithReleaseVersion returns a copy of the object with the release version
// annotation set to the given value.
func (d *TestDeployment) WithReleaseVersion(v string) *TestDeployment {
	newDeployment := d.Copy()
	annotations := newDeployment.GetAnnotations()

	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[util.ReleaseVersionAnnotation] = v
	newDeployment.SetAnnotations(annotations)

	return newDeployment
}

// WithAnnotations returns a copy of the object with the annotations set to the
// given value.
func (d *TestDeployment) WithAnnotations(a map[string]string) *TestDeployment {
	newDeployment := d.Copy()
	newDeployment.SetAnnotations(a)

	return newDeployment
}

// Object returns a copy of the wrapped appsv1.Deployment object.
func (d *TestDeployment) Object() *appsv1.Deployment {
	return d.Deployment.DeepCopy()
}

// TestClusterOperator wraps the ClusterOperator type to add helper methods.
type TestClusterOperator struct {
	configv1.ClusterOperator
}

// NewTestClusterOperator returns a new TestDeployment wrapping the given
// OpenShift ClusterOperator object.
func NewTestClusterOperator(co *configv1.ClusterOperator) *TestClusterOperator {
	return &TestClusterOperator{ClusterOperator: *co}
}

// Copy returns a deep copy of the wrapped object.
func (co *TestClusterOperator) Copy() *TestClusterOperator {
	newCO := &configv1.ClusterOperator{}
	co.ClusterOperator.DeepCopyInto(newCO)

	return NewTestClusterOperator(newCO)
}

// WithConditions returns a copy of the wrapped ClusterOperator object with the
// status conditions set to the given list.
func (co *TestClusterOperator) WithConditions(conds []configv1.ClusterOperatorStatusCondition) *TestClusterOperator {
	newCO := co.Copy()
	newCO.Status.Conditions = conds

	return newCO
}

// WithVersion returns a copy of the object with the "operator"
// OperandVersion's version field set to the given value.
func (co *TestClusterOperator) WithVersion(v string) *TestClusterOperator {
	newCO := co.Copy()

	if newCO.Status.Versions == nil {
		newCO.Status.Versions = []configv1.OperandVersion{{Name: "operator"}}
	}

	found := false

	for i := range newCO.Status.Versions {
		if newCO.Status.Versions[i].Name == "operator" {
			found = true
			newCO.Status.Versions[i].Version = v
		}
	}

	if !found {
		newCO.Status.Versions = append(newCO.Status.Versions, configv1.OperandVersion{
			Name:    "operator",
			Version: v,
		})
	}

	return newCO
}

// Object returns a copy of the wrapped configv1.ClusterOperator object.
func (co *TestClusterOperator) Object() *configv1.ClusterOperator {
	return co.ClusterOperator.DeepCopy()
}

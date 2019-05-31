package operator

import (
	"context"
	"fmt"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	osconfig "github.com/openshift/client-go/config/clientset/versioned"
	autoscalingv1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/util"
	cvorm "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Reason messages used in status conditions.
const (
	ReasonEmpty             = ""
	ReasonMissingDependency = "MissingDependency"
	ReasonSyncing           = "SyncingResources"
	ReasonCheckAutoscaler   = "UnableToCheckAutoscalers"
)

// StatusReporter reports the status of the operator to the OpenShift
// cluster-version-operator via ClusterOperator resource status.
type StatusReporter struct {
	client       client.Client
	configClient osconfig.Interface
	config       *StatusReporterConfig
}

// StatusReporterConfig represents the configuration of a given StatusReporter.
type StatusReporterConfig struct {
	ClusterAutoscalerName      string
	ClusterAutoscalerNamespace string
	ReleaseVersion             string
	RelatedObjects             []configv1.ObjectReference
}

// NewStatusReporter returns a new StatusReporter instance.
func NewStatusReporter(mgr manager.Manager, cfg *StatusReporterConfig) (*StatusReporter, error) {
	var err error

	reporter := &StatusReporter{
		client: mgr.GetClient(),
		config: cfg,
	}

	// Create a client for OpenShift config objects.
	reporter.configClient, err = osconfig.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	return reporter, nil
}

// SetReleaseVersion sets the configured release version.
func (r *StatusReporter) SetReleaseVersion(version string) {
	r.config.ReleaseVersion = version
}

// SetRelatedObjects sets the configured related objects.
func (r *StatusReporter) SetRelatedObjects(objs []configv1.ObjectReference) {
	r.config.RelatedObjects = objs
}

// AddRelatedObjects adds to the list of related objects.
func (r *StatusReporter) AddRelatedObjects(objs []configv1.ObjectReference) {
	for _, obj := range objs {
		r.config.RelatedObjects = append(r.config.RelatedObjects, obj)
	}
}

// GetClusterOperator fetches the the operator's ClusterOperator object.
func (r *StatusReporter) GetClusterOperator() (*configv1.ClusterOperator, error) {
	return r.configClient.ConfigV1().ClusterOperators().Get(OperatorName, metav1.GetOptions{})
}

// GetOrCreateClusterOperator gets, or if necessary, creates the
// operator's ClusterOperator object and returns it.
func (r *StatusReporter) GetOrCreateClusterOperator() (*configv1.ClusterOperator, error) {
	clusterOperator := &configv1.ClusterOperator{
		ObjectMeta: metav1.ObjectMeta{
			Name: OperatorName,
		},
	}

	existing, err := r.GetClusterOperator()

	if errors.IsNotFound(err) {
		return r.configClient.ConfigV1().ClusterOperators().Create(clusterOperator)
	}

	return existing, err
}

// ApplyStatus applies the given ClusterOperator status to the operator's
// ClusterOperator object if necessary.  The currently configured RelatedObjects
// are automatically set on the status.  If no ClusterOperator objects exists,
// one is created.
func (r *StatusReporter) ApplyStatus(status configv1.ClusterOperatorStatus) error {
	var modified bool

	co, err := r.GetOrCreateClusterOperator()
	if err != nil {
		return err
	}

	// Set the currently configured related objects.
	status.RelatedObjects = r.config.RelatedObjects

	// If no versions were set explicitly, continue reporting previous versions.
	if status.Versions == nil {
		status.Versions = co.Status.Versions
	}

	// Update LastTransitionTime for all conditions if necessary.
	for i := range status.Conditions {
		condType := status.Conditions[i].Type
		timestamp := metav1.NewTime(time.Now())

		c := cvorm.FindOperatorStatusCondition(co.Status.Conditions, condType)

		// If found, and status doesn't match, update.
		if c != nil && c.Status != status.Conditions[i].Status {
			status.Conditions[i].LastTransitionTime = timestamp
		}

		// If found, and status matches, copy previous.
		if c != nil && c.Status == status.Conditions[i].Status {
			status.Conditions[i].LastTransitionTime = c.LastTransitionTime
		}

		// If it's still nil, update it.
		if status.Conditions[i].LastTransitionTime.IsZero() {
			status.Conditions[i].LastTransitionTime = timestamp
		}
	}

	// If any versions have changed, we need to reset the transition time on the
	// Progressing condition whether or not we actually had any work to do.
	if !equality.Semantic.DeepEqual(status.Versions, co.Status.Versions) {
		util.ResetProgressingTime(&status.Conditions)
	}

	// Copy the current ClusterOperator and overwrite the status.
	requiredCO := &configv1.ClusterOperator{}
	co.DeepCopyInto(requiredCO)
	requiredCO.Status = status

	cvorm.EnsureClusterOperatorStatus(&modified, co, *requiredCO)

	if modified {
		_, err := r.configClient.ConfigV1().ClusterOperators().UpdateStatus(co)
		return err
	}

	return nil
}

// available reports the operator as available, not progressing, and
// not degraded -- optionally setting a reason and message.  This will
// update the reported operator version.  It should only be called if
// the operands are fully updated and available.
func (r *StatusReporter) available(reason, message string) error {
	status := configv1.ClusterOperatorStatus{
		Conditions: []configv1.ClusterOperatorStatusCondition{
			{
				Type:    configv1.OperatorAvailable,
				Status:  configv1.ConditionTrue,
				Reason:  reason,
				Message: message,
			},
			{
				Type:   configv1.OperatorProgressing,
				Status: configv1.ConditionFalse,
			},
			{
				Type:   configv1.OperatorDegraded,
				Status: configv1.ConditionFalse,
			},
		},
		Versions: []configv1.OperandVersion{
			{
				Name:    "operator",
				Version: r.config.ReleaseVersion,
			},
		},
	}

	klog.Infof("Operator status available: %s", message)

	return r.ApplyStatus(status)
}

// degraded reports the operator as degraded but available, and not
// progressing -- optionally setting a reason and message.
func (r *StatusReporter) degraded(reason, message string) error {
	status := configv1.ClusterOperatorStatus{
		Conditions: []configv1.ClusterOperatorStatusCondition{
			{
				Type:   configv1.OperatorAvailable,
				Status: configv1.ConditionTrue,
			},
			{
				Type:   configv1.OperatorProgressing,
				Status: configv1.ConditionFalse,
			},
			{
				Type:    configv1.OperatorDegraded,
				Status:  configv1.ConditionTrue,
				Reason:  reason,
				Message: message,
			},
		},
	}

	klog.Warningf("Operator status degraded: %s", message)

	return r.ApplyStatus(status)
}

// progressing reports the operator as progressing but available, and not
// degraded -- optionally setting a reason and message.
func (r *StatusReporter) progressing(reason, message string) error {
	status := configv1.ClusterOperatorStatus{
		Conditions: []configv1.ClusterOperatorStatusCondition{
			{
				Type:   configv1.OperatorAvailable,
				Status: configv1.ConditionTrue,
			},
			{
				Type:    configv1.OperatorProgressing,
				Status:  configv1.ConditionTrue,
				Reason:  reason,
				Message: message,
			},
			{
				Type:   configv1.OperatorDegraded,
				Status: configv1.ConditionFalse,
			},
		},
	}

	klog.Infof("Operator status progressing: %s", message)

	return r.ApplyStatus(status)
}

// Start checks the status of dependencies and reports the operator's status. It
// will poll until stopCh is closed or prerequisites are satisfied, in which
// case it will report the operator as available the configured version and wait
// for stopCh to close before returning.
func (r *StatusReporter) Start(stopCh <-chan struct{}) error {
	interval := 15 * time.Second

	// Poll the status of our prerequisites and set our status
	// accordingly.  Rather than return errors and stop polling, most
	// errors here should just be reported in the status message.
	pollFunc := func() (bool, error) {
		return r.ReportStatus()
	}

	err := wait.PollImmediateUntil(interval, pollFunc, stopCh)

	// Block until the stop channel is closed.
	<-stopCh

	return err
}

// ReportStatus checks the status of each dependency and operand and reports the
// appropriate status via the operator's ClusterOperator object.
func (r *StatusReporter) ReportStatus() (bool, error) {
	// Check that the machine-api-operator is reporting available.
	ok, err := r.CheckMachineAPI()
	if err != nil {
		msg := fmt.Sprintf("error checking machine-api status: %v", err)
		r.degraded(ReasonMissingDependency, msg)
		return false, nil
	}

	if !ok {
		r.degraded(ReasonMissingDependency, "machine-api not ready")
		return false, nil
	}

	// Check that any CluterAutoscaler deployments are updated and available.
	ok, err = r.CheckClusterAutoscaler()
	if err != nil {
		msg := fmt.Sprintf("error checking autoscaler status: %v", err)
		r.degraded(ReasonCheckAutoscaler, msg)
		return false, nil
	}

	if !ok {
		msg := fmt.Sprintf("updating to %s", r.config.ReleaseVersion)
		r.progressing(ReasonSyncing, msg)
		return false, nil
	}

	msg := fmt.Sprintf("at version %s", r.config.ReleaseVersion)
	r.available(ReasonEmpty, msg)

	return true, nil
}

// CheckMachineAPI checks the status of the machine-api-operator as
// reported to the CVO.  It returns true if the operator is available
// and not degraded.
func (r *StatusReporter) CheckMachineAPI() (bool, error) {
	mao, err := r.configClient.ConfigV1().ClusterOperators().
		Get("machine-api", metav1.GetOptions{})

	if err != nil {
		klog.Errorf("failed to get dependency machine-api status: %v", err)
		return false, err
	}

	conds := mao.Status.Conditions

	if cvorm.IsOperatorStatusConditionTrue(conds, configv1.OperatorAvailable) &&
		cvorm.IsOperatorStatusConditionFalse(conds, configv1.OperatorDegraded) {
		return true, nil
	}

	klog.Infof("machine-api-operator not ready yet")
	return false, nil
}

// CheckClusterAutoscaler checks the status of any cluster-autoscaler
// deployments. It returns a bool indicating whether the deployments are
// available and fully updated to the latest version and an error.
func (r *StatusReporter) CheckClusterAutoscaler() (bool, error) {
	ca := &autoscalingv1.ClusterAutoscaler{}
	caName := client.ObjectKey{Name: r.config.ClusterAutoscalerName}

	if err := r.client.Get(context.TODO(), caName, ca); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("No ClusterAutoscaler. Reporting available.")
			return true, nil
		}

		klog.Errorf("Error getting ClusterAutoscaler: %v", err)
		return false, err
	}

	deployment := &appsv1.Deployment{}
	deploymentName := client.ObjectKey{
		Name:      fmt.Sprintf("%s-%s", OperatorName, r.config.ClusterAutoscalerName),
		Namespace: r.config.ClusterAutoscalerNamespace,
	}

	if err := r.client.Get(context.TODO(), deploymentName, deployment); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("No ClusterAutoscaler deployment. Reporting unavailable.")
			return false, nil
		}

		klog.Errorf("Error getting ClusterAutoscaler deployment: %v", err)
		return false, err
	}

	if !util.ReleaseVersionMatches(deployment, r.config.ReleaseVersion) {
		klog.Info("ClusterAutoscaler deployment version not current.")
		return false, nil
	}

	if !util.DeploymentUpdated(deployment) {
		klog.Info("ClusterAutoscaler deployment updating.")
		return false, nil
	}

	klog.Info("ClusterAutoscaler deployment is available and updated.")

	return true, nil
}

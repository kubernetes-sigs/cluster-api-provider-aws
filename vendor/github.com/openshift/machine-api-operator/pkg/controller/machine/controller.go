/*
Copyright 2018 The Kubernetes Authors.

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

package machine

import (
	"context"
	"errors"
	"fmt"
	"time"

	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/openshift/machine-api-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"k8s.io/kubectl/pkg/drain"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	NodeNameEnvVar = "NODE_NAME"
	requeueAfter   = 30 * time.Second

	// ExcludeNodeDrainingAnnotation annotation explicitly skips node draining if set
	ExcludeNodeDrainingAnnotation = "machine.openshift.io/exclude-node-draining"

	// MachineRegionLabelName as annotation name for a machine region
	MachineRegionLabelName = "machine.openshift.io/region"

	// MachineAZLabelName as annotation name for a machine AZ
	MachineAZLabelName = "machine.openshift.io/zone"

	// MachineInstanceStateAnnotationName as annotation name for a machine instance state
	MachineInstanceStateAnnotationName = "machine.openshift.io/instance-state"

	// MachineInstanceTypeLabelName as annotation name for a machine instance type
	MachineInstanceTypeLabelName = "machine.openshift.io/instance-type"

	// MachineInterruptibleInstanceLabelName as annotaiton name for interruptible instances
	MachineInterruptibleInstanceLabelName = "machine.openshift.io/interruptible-instance"

	// https://github.com/openshift/enhancements/blob/master/enhancements/machine-instance-lifecycle.md
	// This is not a transient error, but
	// indicates a state that will likely need to be fixed before progress can be made
	// e.g Instance does NOT exist but Machine has providerID/address
	// e.g Cloud service returns a 4xx response
	phaseFailed = "Failed"

	// Instance does NOT exist
	// Machine has NOT been given providerID/address
	phaseProvisioning = "Provisioning"

	// Instance exists
	// Machine has been given providerID/address
	// Machine has NOT been given nodeRef
	phaseProvisioned = "Provisioned"

	// Instance exists
	// Machine has been given providerID/address
	// Machine has been given a nodeRef
	phaseRunning = "Running"

	// Machine has a deletion timestamp
	phaseDeleting = "Deleting"

	// Hardcoded instance state set on machine failure
	unknownInstanceState = "Unknown"

	skipWaitForDeleteTimeoutSeconds = 60 * 5
)

var DefaultActuator Actuator

func AddWithActuator(mgr manager.Manager, actuator Actuator) error {
	return add(mgr, newReconciler(mgr, actuator))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, actuator Actuator) reconcile.Reconciler {
	r := &ReconcileMachine{
		Client:        mgr.GetClient(),
		eventRecorder: mgr.GetEventRecorderFor("machine-controller"),
		config:        mgr.GetConfig(),
		scheme:        mgr.GetScheme(),
		actuator:      actuator,
	}
	return r
}

func stringPointerDeref(stringPointer *string) string {
	if stringPointer != nil {
		return *stringPointer
	}
	return ""
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("machine_controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Machine
	return c.Watch(
		&source.Kind{Type: &machinev1.Machine{}},
		&handler.EnqueueRequestForObject{},
	)
}

// ReconcileMachine reconciles a Machine object
type ReconcileMachine struct {
	client.Client
	config *rest.Config
	scheme *runtime.Scheme

	eventRecorder record.EventRecorder

	actuator Actuator
}

// Reconcile reads that state of the cluster for a Machine object and makes changes based on the state read
// and what is in the Machine.Spec
// +kubebuilder:rbac:groups=machine.openshift.io,resources=machines;machines/status,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileMachine) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// TODO(mvladev): Can context be passed from Kubebuilder?
	ctx := context.TODO()

	// Fetch the Machine instance
	m := &machinev1.Machine{}
	if err := r.Client.Get(ctx, request.NamespacedName, m); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Implement controller logic here
	machineName := m.GetName()
	klog.Infof("%v: reconciling Machine", machineName)

	if errList := m.Validate(); len(errList) > 0 {
		err := fmt.Errorf("%v: machine validation failed: %v", machineName, errList.ToAggregate().Error())
		klog.Error(err)
		r.eventRecorder.Eventf(m, corev1.EventTypeWarning, "FailedValidate", err.Error())
		return reconcile.Result{}, err
	}

	// If object hasn't been deleted and doesn't have a finalizer, add one
	// Add a finalizer to newly created objects.
	if m.ObjectMeta.DeletionTimestamp.IsZero() {
		finalizerCount := len(m.Finalizers)

		if !util.Contains(m.Finalizers, machinev1.MachineFinalizer) {
			m.Finalizers = append(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer)
		}

		if len(m.Finalizers) > finalizerCount {
			if err := r.Client.Update(ctx, m); err != nil {
				klog.Infof("%v: failed to add finalizers to machine: %v", machineName, err)
				return reconcile.Result{}, err
			}

			// Since adding the finalizer updates the object return to avoid later update issues
			return reconcile.Result{}, nil
		}
	}

	if !m.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := r.setPhase(m, phaseDeleting, ""); err != nil {
			return reconcile.Result{}, err
		}

		// no-op if finalizer has been removed.
		if !util.Contains(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer) {
			klog.Infof("%v: reconciling machine causes a no-op as there is no finalizer", machineName)
			return reconcile.Result{}, nil
		}

		klog.Infof("%v: reconciling machine triggers delete", machineName)
		// Drain node before deletion
		// If a machine is not linked to a node, just delete the machine. Since a node
		// can be unlinked from a machine when the node goes NotReady and is removed
		// by cloud controller manager. In that case some machines would never get
		// deleted without a manual intervention.
		if _, exists := m.ObjectMeta.Annotations[ExcludeNodeDrainingAnnotation]; !exists && m.Status.NodeRef != nil {
			if err := r.drainNode(m); err != nil {
				klog.Errorf("%v: failed to drain node for machine: %v", machineName, err)
				return delayIfRequeueAfterError(err)
			}
		}

		if err := r.actuator.Delete(ctx, m); err != nil {
			// isInvalidMachineConfiguration will take care of the case where the
			// configuration is invalid from the beginning. len(m.Status.Addresses) > 0
			// will handle the case when a machine configuration was invalidated
			// after an instance was created. So only a small window is left when
			// we can loose instances, e.g. right after request to create one
			// was sent and before a list of node addresses was set.
			if len(m.Status.Addresses) > 0 || !isInvalidMachineConfigurationError(err) {
				klog.Errorf("%v: failed to delete machine: %v", machineName, err)
				return delayIfRequeueAfterError(err)
			}
		}

		instanceExists, err := r.actuator.Exists(ctx, m)
		if err != nil {
			klog.Errorf("%v: failed to check if machine exists: %v", machineName, err)
			return reconcile.Result{}, err
		}

		if instanceExists {
			klog.V(3).Infof("%v: can't proceed deleting machine while cloud instance is being terminated, requeuing", machineName)
			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		}

		if m.Status.NodeRef != nil {
			klog.Infof("%v: deleting node %q for machine", machineName, m.Status.NodeRef.Name)
			if err := r.deleteNode(ctx, m.Status.NodeRef.Name); err != nil {
				klog.Errorf("%v: error deleting node for machine: %v", machineName, err)
				return reconcile.Result{}, err
			}
		}

		// Remove finalizer on successful deletion.
		m.ObjectMeta.Finalizers = util.Filter(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer)
		if err := r.Client.Update(context.Background(), m); err != nil {
			klog.Errorf("%v: failed to remove finalizer from machine: %v", machineName, err)
			return reconcile.Result{}, err
		}

		klog.Infof("%v: machine deletion successful", machineName)
		return reconcile.Result{}, nil
	}

	if machineIsFailed(m) {
		klog.Warningf("%v: machine has gone %q phase. It won't reconcile", machineName, phaseFailed)
		return reconcile.Result{}, nil
	}

	instanceExists, err := r.actuator.Exists(ctx, m)
	if err != nil {
		klog.Errorf("%v: failed to check if machine exists: %v", machineName, err)
		return reconcile.Result{}, err
	}

	if instanceExists {
		klog.Infof("%v: reconciling machine triggers idempotent update", machineName)
		if err := r.actuator.Update(ctx, m); err != nil {
			klog.Errorf("%v: error updating machine: %v", machineName, err)
			return delayIfRequeueAfterError(err)
		}

		if !machineIsProvisioned(m) {
			klog.Errorf("%v: instance exists but providerID or addresses has not been given to the machine yet, requeuing", machineName)
			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		}

		if !machineHasNode(m) {
			// Requeue until we reach running phase
			if err := r.setPhase(m, phaseProvisioned, ""); err != nil {
				return reconcile.Result{}, err
			}
			klog.Infof("%v: has no node yet, requeuing", machineName)
			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		}

		return reconcile.Result{}, r.setPhase(m, phaseRunning, "")
	}

	// Instance does not exist but the machine has been given a providerID/address.
	// This can only be reached if an instance was deleted outside the machine API
	if machineIsProvisioned(m) {
		if err := r.setPhase(m, phaseFailed, "Can't find created instance."); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// Machine resource created and instance does not exist yet.
	if err := r.setPhase(m, phaseProvisioning, ""); err != nil {
		return reconcile.Result{}, err
	}
	klog.Infof("%v: reconciling machine triggers idempotent create", machineName)
	if err := r.actuator.Create(ctx, m); err != nil {
		klog.Warningf("%v: failed to create machine: %v", machineName, err)
		if isInvalidMachineConfigurationError(err) {
			if err := r.setPhase(m, phaseFailed, err.Error()); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		return delayIfRequeueAfterError(err)
	}

	klog.Infof("%v: created instance, requeuing", machineName)
	return reconcile.Result{RequeueAfter: requeueAfter}, nil
}

func (r *ReconcileMachine) drainNode(machine *machinev1.Machine) error {
	kubeClient, err := kubernetes.NewForConfig(r.config)
	if err != nil {
		return fmt.Errorf("unable to build kube client: %v", err)
	}
	node, err := kubeClient.CoreV1().Nodes().Get(context.Background(), machine.Status.NodeRef.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If an admin deletes the node directly, we'll end up here.
			klog.Infof("Could not find node from noderef, it may have already been deleted: %v", machine.Status.NodeRef.Name)
			return nil
		}
		return fmt.Errorf("unable to get node %q: %v", machine.Status.NodeRef.Name, err)
	}

	drainer := &drain.Helper{
		Client:              kubeClient,
		Force:               true,
		IgnoreAllDaemonSets: true,
		DeleteLocalData:     true,
		GracePeriodSeconds:  -1,
		// If a pod is not evicted in 20 seconds, retry the eviction next time the
		// machine gets reconciled again (to allow other machines to be reconciled).
		Timeout: 20 * time.Second,
		OnPodDeletedOrEvicted: func(pod *corev1.Pod, usingEviction bool) {
			verbStr := "Deleted"
			if usingEviction {
				verbStr = "Evicted"
			}
			klog.Info(fmt.Sprintf("%s pod from Node", verbStr),
				"pod", fmt.Sprintf("%s/%s", pod.Name, pod.Namespace))
		},
		Out:    writer{klog.Info},
		ErrOut: writer{klog.Error},
	}

	if nodeIsUnreachable(node) {
		klog.Infof("%q: Node %q is unreachable, draining will wait %d seconds after pod is signalled for deletion and skip after it",
			machine.Name, node.Name, skipWaitForDeleteTimeoutSeconds)
		drainer.SkipWaitForDeleteTimeoutSeconds = skipWaitForDeleteTimeoutSeconds
	}

	if err := drain.RunCordonOrUncordon(drainer, node, true); err != nil {
		// Can't cordon a node
		klog.Warningf("cordon failed for node %q: %v", node.Name, err)
		return &RequeueAfterError{RequeueAfter: 20 * time.Second}
	}

	if err := drain.RunNodeDrain(drainer, node.Name); err != nil {
		// Machine still tries to terminate after drain failure
		klog.Warningf("drain failed for machine %q: %v", machine.Name, err)
		return &RequeueAfterError{RequeueAfter: 20 * time.Second}
	}

	klog.Infof("drain successful for machine %q", machine.Name)
	r.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Node %q drained", node.Name)

	return nil
}

func (r *ReconcileMachine) deleteNode(ctx context.Context, name string) error {
	var node corev1.Node
	if err := r.Client.Get(ctx, client.ObjectKey{Name: name}, &node); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Infof("Node %q not found", name)
			return nil
		}
		klog.Errorf("Failed to get node %q: %v", name, err)
		return err
	}
	return r.Client.Delete(ctx, &node)
}

func delayIfRequeueAfterError(err error) (reconcile.Result, error) {
	var requeueAfterError *RequeueAfterError
	if errors.As(err, &requeueAfterError) {
		klog.Infof("Actuator returned requeue-after error: %v", requeueAfterError)
		return reconcile.Result{Requeue: true, RequeueAfter: requeueAfterError.RequeueAfter}, nil
	}
	return reconcile.Result{}, err
}

func isInvalidMachineConfigurationError(err error) bool {
	var machineError *MachineError
	if errors.As(err, &machineError) {
		if machineError.Reason == machinev1.InvalidConfigurationMachineError {
			klog.Infof("Actuator returned invalid configuration error: %v", machineError)
			return true
		}
	}
	return false
}

func (r *ReconcileMachine) setPhase(machine *machinev1.Machine, phase string, errorMessage string) error {
	if stringPointerDeref(machine.Status.Phase) != phase {
		klog.V(3).Infof("%v: going into phase %q", machine.GetName(), phase)
		// A call to Patch will mutate our local copy of the machine to match what is stored in the API.
		// Before we make any changes to the status subresource on our local copy, we need to patch the object first,
		// otherwise our local changes to the status subresource will be lost.
		if phase == phaseFailed {
			err := r.patchFailedMachineInstanceAnnotation(machine)
			if err != nil {
				klog.Errorf("Failed to update machine %q: %v", machine.GetName(), err)
				return err
			}
		}

		// Since we may have mutated the local copy of the machine above, we need to calculate baseToPatch here.
		// Any updates to the status must be done after this point.
		baseToPatch := client.MergeFrom(machine.DeepCopy())

		if phase == phaseFailed {
			if err := r.overrideFailedMachineProviderStatusState(machine); err != nil {
				klog.Errorf("Failed to update machine provider status %q: %v", machine.GetName(), err)
				return err
			}
		}

		machine.Status.Phase = &phase
		machine.Status.ErrorMessage = nil
		now := metav1.Now()
		machine.Status.LastUpdated = &now
		if phase == phaseFailed && errorMessage != "" {
			machine.Status.ErrorMessage = &errorMessage
		}
		if err := r.Client.Status().Patch(context.Background(), machine, baseToPatch); err != nil {
			klog.Errorf("Failed to update machine status %q: %v", machine.GetName(), err)
			return err
		}
	}
	return nil
}

func (r *ReconcileMachine) patchFailedMachineInstanceAnnotation(machine *machinev1.Machine) error {
	baseToPatch := client.MergeFrom(machine.DeepCopy())
	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}
	machine.Annotations[MachineInstanceStateAnnotationName] = unknownInstanceState
	if err := r.Client.Patch(context.Background(), machine, baseToPatch); err != nil {
		return err
	}
	return nil
}

// overrideFailedMachineProviderStatusState patches the state of the VM in the provider status if it is set.
// Not all providers set a state, but AWS, Azure, GCP and vSphere do.
// If the machine has gone into the Failed phase, and the providerStatus has already been set,
// the VM is in an unknown state. This function overrides the state.
func (r *ReconcileMachine) overrideFailedMachineProviderStatusState(machine *machinev1.Machine) error {
	if machine.Status.ProviderStatus == nil {
		return nil
	}

	// instanceState is used by AWS, GCP and vSphere; vmState is used by Azure.
	const instanceStateField = "instanceState"
	const vmStateField = "vmState"

	providerStatus, err := runtime.DefaultUnstructuredConverter.ToUnstructured(machine.Status.ProviderStatus)
	if err != nil {
		return fmt.Errorf("could not covert provider status to unstructured: %v", err)
	}

	// if the instanceState is set already, update it to unknown
	if _, found, err := unstructured.NestedString(providerStatus, instanceStateField); err == nil && found {
		if err := unstructured.SetNestedField(providerStatus, unknownInstanceState, instanceStateField); err != nil {
			return fmt.Errorf("could not set %s: %v", instanceStateField, err)
		}
	}

	// if the vmState is set already, update it to unknown
	if _, found, err := unstructured.NestedString(providerStatus, vmStateField); err == nil && found {
		if err := unstructured.SetNestedField(providerStatus, unknownInstanceState, vmStateField); err != nil {
			return fmt.Errorf("could not set %s: %v", instanceStateField, err)
		}
	}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(providerStatus, machine.Status.ProviderStatus); err != nil {
		return fmt.Errorf("could not convert provider status from unstructured: %v", err)
	}

	return nil
}

func machineIsProvisioned(machine *machinev1.Machine) bool {
	return len(machine.Status.Addresses) > 0 || stringPointerDeref(machine.Spec.ProviderID) != ""
}

func machineHasNode(machine *machinev1.Machine) bool {
	return machine.Status.NodeRef != nil
}

func machineIsFailed(machine *machinev1.Machine) bool {
	return stringPointerDeref(machine.Status.Phase) == phaseFailed
}

func nodeIsUnreachable(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionUnknown {
			return true
		}
	}

	return false
}

// writer implements io.Writer interface as a pass-through for klog.
type writer struct {
	logFunc func(args ...interface{})
}

// Write passes string(p) into writer's logFunc and always returns len(p)
func (w writer) Write(p []byte) (n int, err error) {
	w.logFunc(string(p))
	return len(p), nil
}

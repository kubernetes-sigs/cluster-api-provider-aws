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
	"fmt"
	"time"

	"github.com/go-log/log/info"
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	commonerrors "github.com/openshift/cluster-api/pkg/apis/machine/common"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	controllerError "github.com/openshift/cluster-api/pkg/controller/error"
	kubedrain "github.com/openshift/cluster-api/pkg/drain"
	clusterapiError "github.com/openshift/cluster-api/pkg/errors"
	"github.com/openshift/cluster-api/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	NodeNameEnvVar = "NODE_NAME"

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
	name := m.Name
	klog.Infof("Reconciling Machine %q", name)

	if errList := m.Validate(); len(errList) > 0 {
		err := fmt.Errorf("%q machine validation failed: %v", m.Name, errList.ToAggregate().Error())
		klog.Error(err)
		r.eventRecorder.Eventf(m, corev1.EventTypeWarning, "FailedValidate", err.Error())
		return reconcile.Result{}, err
	}

	// Cluster might be nil as some providers might not require a cluster object
	// for machine management.
	cluster, err := r.getCluster(ctx, m)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set the ownerRef with foreground deletion if there is a linked cluster.
	if cluster != nil && len(m.OwnerReferences) == 0 {
		blockOwnerDeletion := true
		m.OwnerReferences = append(m.OwnerReferences, metav1.OwnerReference{
			APIVersion:         cluster.APIVersion,
			Kind:               cluster.Kind,
			Name:               cluster.Name,
			UID:                cluster.UID,
			BlockOwnerDeletion: &blockOwnerDeletion,
		})
	}

	// If object hasn't been deleted and doesn't have a finalizer, add one
	// Add a finalizer to newly created objects.
	if m.ObjectMeta.DeletionTimestamp.IsZero() {
		finalizerCount := len(m.Finalizers)

		if cluster != nil && !util.Contains(m.Finalizers, metav1.FinalizerDeleteDependents) {
			m.Finalizers = append(m.ObjectMeta.Finalizers, metav1.FinalizerDeleteDependents)
		}

		if !util.Contains(m.Finalizers, machinev1.MachineFinalizer) {
			m.Finalizers = append(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer)
		}

		if len(m.Finalizers) > finalizerCount {
			if err := r.Client.Update(ctx, m); err != nil {
				klog.Infof("Failed to add finalizers to machine %q: %v", name, err)
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
			klog.Infof("Reconciling machine %q causes a no-op as there is no finalizer", name)
			return reconcile.Result{}, nil
		}

		klog.Infof("Reconciling machine %q triggers delete", name)
		// Drain node before deletion
		// If a machine is not linked to a node, just delete the machine. Since a node
		// can be unlinked from a machine when the node goes NotReady and is removed
		// by cloud controller manager. In that case some machines would never get
		// deleted without a manual intervention.
		if _, exists := m.ObjectMeta.Annotations[ExcludeNodeDrainingAnnotation]; !exists && m.Status.NodeRef != nil {
			if err := r.drainNode(m); err != nil {
				klog.Errorf("Failed to drain node for machine %q: %v", name, err)
				return delayIfRequeueAfterError(err)
			}
		}

		if err := r.actuator.Delete(ctx, cluster, m); err != nil {
			// isInvalidMachineConfiguration will take care of the case where the
			// configuration is invalid from the beginning. len(m.Status.Addresses) > 0
			// will handle the case when a machine configuration was invalidated
			// after an instance was created. So only a small window is left when
			// we can loose instances, e.g. right after request to create one
			// was sent and before a list of node addresses was set.
			if len(m.Status.Addresses) > 0 || !isInvalidMachineConfigurationError(err) {
				klog.Errorf("Failed to delete machine %q: %v", name, err)
				return delayIfRequeueAfterError(err)
			}
		}

		if m.Status.NodeRef != nil {
			klog.Infof("Deleting node %q for machine %q", m.Status.NodeRef.Name, m.Name)
			if err := r.deleteNode(ctx, m.Status.NodeRef.Name); err != nil {
				klog.Errorf("Error deleting node %q for machine %q", name, err)
				return reconcile.Result{}, err
			}
		}

		// Remove finalizer on successful deletion.
		m.ObjectMeta.Finalizers = util.Filter(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer)
		if err := r.Client.Update(context.Background(), m); err != nil {
			klog.Errorf("Failed to remove finalizer from machine %q: %v", name, err)
			return reconcile.Result{}, err
		}

		klog.Infof("Machine %q deletion successful", name)
		return reconcile.Result{}, nil
	}

	if machineIsFailed(m) {
		klog.Warningf("Machine %q has gone %q phase. It won't reconcile", name, phaseFailed)
		return reconcile.Result{}, nil
	}

	instanceExists, err := r.actuator.Exists(ctx, cluster, m)
	if err != nil {
		klog.Errorf("Failed to check if machine %q exists: %v", name, err)
		return reconcile.Result{}, err
	}

	if instanceExists {
		klog.Infof("Reconciling machine %q triggers idempotent update", name)
		if err := r.actuator.Update(ctx, cluster, m); err != nil {
			klog.Errorf(`Error updating machine "%s/%s": %v`, m.Namespace, name, err)
			return delayIfRequeueAfterError(err)
		}

		if !machineIsProvisioned(m) {
			klog.Errorf(`Instance for Machine "%s/%s exists but providerID or addresses has not been given to the machine yet"`, m.Namespace, name)
			return reconcile.Result{}, err
		}
		if machineHasNode(m) {
			if err := r.setPhase(m, phaseRunning, ""); err != nil {
				return reconcile.Result{}, err
			}
		} else {
			if err := r.setPhase(m, phaseProvisioned, ""); err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
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
	klog.Infof("Reconciling machine object %v triggers idempotent create.", m.ObjectMeta.Name)
	if err := r.actuator.Create(ctx, cluster, m); err != nil {
		klog.Warningf("Failed to create machine %q: %v", name, err)
		if isInvalidMachineConfigurationError(err) {
			if err := r.setPhase(m, phaseFailed, err.Error()); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		return delayIfRequeueAfterError(err)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileMachine) drainNode(machine *machinev1.Machine) error {
	kubeClient, err := kubernetes.NewForConfig(r.config)
	if err != nil {
		return fmt.Errorf("unable to build kube client: %v", err)
	}
	node, err := kubeClient.CoreV1().Nodes().Get(machine.Status.NodeRef.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If an admin deletes the node directly, we'll end up here.
			klog.Infof("Could not find node from noderef, it may have already been deleted: %v", machine.Status.NodeRef.Name)
			return nil
		}
		return fmt.Errorf("unable to get node %q: %v", machine.Status.NodeRef.Name, err)
	}

	if err := kubedrain.Drain(
		kubeClient,
		[]*corev1.Node{node},
		&kubedrain.DrainOptions{
			Force:              true,
			IgnoreDaemonsets:   true,
			DeleteLocalData:    true,
			GracePeriodSeconds: -1,
			Logger:             info.New(klog.V(0)),
			// If a pod is not evicted in 20 second, retry the eviction next time the
			// machine gets reconciled again (to allow other machines to be reconciled)
			Timeout: 20 * time.Second,
		},
	); err != nil {
		// Machine still tries to terminate after drain failure
		klog.Warningf("drain failed for machine %q: %v", machine.Name, err)
		return &controllerError.RequeueAfterError{RequeueAfter: 20 * time.Second}
	}

	klog.Infof("drain successful for machine %q", machine.Name)
	r.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Node %q drained", node.Name)

	return nil
}

func (r *ReconcileMachine) getCluster(ctx context.Context, machine *machinev1.Machine) (*clusterv1.Cluster, error) {
	if machine.Labels[machinev1.MachineClusterLabelName] == "" {
		klog.Infof("Machine %q in namespace %q doesn't specify %q label, assuming nil cluster", machine.Name, machine.Namespace, machinev1.MachineClusterLabelName)
		return nil, nil
	}

	cluster := &clusterv1.Cluster{}
	key := client.ObjectKey{
		Namespace: machine.Namespace,
		Name:      machine.Labels[machinev1.MachineClusterLabelName],
	}

	if err := r.Client.Get(ctx, key, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
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
	switch t := err.(type) {
	case *controllerError.RequeueAfterError:
		klog.Infof("Actuator returned requeue-after error: %v", err)
		return reconcile.Result{Requeue: true, RequeueAfter: t.RequeueAfter}, nil
	}
	return reconcile.Result{}, err
}

func isInvalidMachineConfigurationError(err error) bool {
	switch t := err.(type) {
	case *clusterapiError.MachineError:
		if t.Reason == commonerrors.InvalidConfigurationMachineError {
			klog.Infof("Actuator returned invalid configuration error: %v", err)
			return true
		}
	}
	return false
}

func (r *ReconcileMachine) setPhase(machine *machinev1.Machine, phase string, errorMessage string) error {
	if stringPointerDeref(machine.Status.Phase) != phase {
		klog.V(3).Infof("Machine %q going into phase %q", machine.GetName(), phase)
		baseToPatch := client.MergeFrom(machine.DeepCopy())
		machine.Status.Phase = &phase
		machine.Status.ErrorMessage = nil
		now := metav1.Now()
		machine.Status.LastUpdated = &now
		if phase == phaseFailed && errorMessage != "" {
			machine.Status.ErrorMessage = &errorMessage
		}
		if err := r.Client.Status().Patch(context.Background(), machine, baseToPatch); err != nil {
			klog.Errorf("Failed to update machine %q: %v", machine.GetName(), err)
			return err
		}
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
	if stringPointerDeref(machine.Status.Phase) == phaseFailed {
		return true
	}
	return false
}

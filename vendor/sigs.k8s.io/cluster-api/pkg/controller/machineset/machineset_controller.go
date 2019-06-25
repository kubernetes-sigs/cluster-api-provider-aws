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

package machineset

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// controllerName is the name of this controller
const controllerName = "machineset-controller"

var (
	// controllerKind contains the schema.GroupVersionKind for this controller type.
	controllerKind = clusterv1alpha1.SchemeGroupVersion.WithKind("MachineSet")

	// stateConfirmationTimeout is the amount of time allowed to wait for desired state.
	stateConfirmationTimeout = 10 * time.Second

	// stateConfirmationInterval is the amount of time between polling for the desired state.
	// The polling is against a local memory cache.
	stateConfirmationInterval = 100 * time.Millisecond
)

// Add creates a new MachineSet Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	r := newReconciler(mgr)
	return add(mgr, r, r.MachineToMachineSets)
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) *ReconcileMachineSet {
	return &ReconcileMachineSet{Client: mgr.GetClient(), scheme: mgr.GetScheme(), recorder: mgr.GetRecorder(controllerName)}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler, mapFn handler.ToRequestsFunc) error {
	// Create a new controller.
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to MachineSet.
	err = c.Watch(
		&source.Kind{Type: &clusterv1alpha1.MachineSet{}},
		&handler.EnqueueRequestForObject{},
	)
	if err != nil {
		return err
	}

	// Watch for changes to Machines and reconcile the owner MachineSet.
	err = c.Watch(
		&source.Kind{Type: &clusterv1alpha1.Machine{}},
		&handler.EnqueueRequestForOwner{IsController: true, OwnerType: &clusterv1alpha1.MachineSet{}},
	)
	if err != nil {
		return err
	}

	// Watch for changes to Machines using a mapping function to MachineSets.
	// This watcher is required for use cases like adoption. In case a Machine doesn't have
	// a controller reference, it'll look for potential matching MachineSet to reconcile.
	return c.Watch(
		&source.Kind{Type: &clusterv1alpha1.Machine{}},
		&handler.EnqueueRequestsFromMapFunc{ToRequests: mapFn},
	)
}

// ReconcileMachineSet reconciles a MachineSet object
type ReconcileMachineSet struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a MachineSet object and makes changes based on the state read
// and what is in the MachineSet.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=cluster.k8s.io,resources=machinesets;machinesets/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.k8s.io,resources=machines,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileMachineSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the MachineSet instance
	ctx := context.TODO()
	machineSet := &clusterv1alpha1.MachineSet{}
	if err := r.Get(ctx, request.NamespacedName, machineSet); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Ignore deleted MachineSets, this can happen when foregroundDeletion
	// is enabled
	if machineSet.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	result, err := r.reconcile(ctx, machineSet)
	if err != nil {
		klog.Errorf("Failed to reconcile MachineSet %q: %v", request.NamespacedName, err)
		r.recorder.Eventf(machineSet, corev1.EventTypeWarning, "ReconcileError", "%v", err)
	}
	return result, err
}

func (r *ReconcileMachineSet) reconcile(ctx context.Context, machineSet *clusterv1alpha1.MachineSet) (reconcile.Result, error) {
	klog.V(4).Infof("Reconcile machineset %v", machineSet.Name)
	allMachines := &clusterv1alpha1.MachineList{}

	if err := r.Client.List(context.Background(), client.InNamespace(machineSet.Namespace), allMachines); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to list machines")
	}

	// Make sure that label selector can match template's labels.
	// TODO(vincepri): Move to a validation (admission) webhook when supported.
	selector, err := metav1.LabelSelectorAsSelector(&machineSet.Spec.Selector)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to parse MachineSet %q label selector", machineSet.Name)
	}

	if !selector.Matches(labels.Set(machineSet.Spec.Template.Labels)) {
		return reconcile.Result{}, errors.Errorf("failed validation on MachineSet %q label selector, cannot match any machines ", machineSet.Name)
	}

	// Cluster might be nil as some providers might not require a cluster object
	// for machine management.
	cluster, err := r.getCluster(machineSet)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set the ownerRef with foreground deletion if there is a linked cluster.
	if cluster != nil && len(machineSet.OwnerReferences) == 0 {
		blockOwnerDeletion := true
		machineSet.OwnerReferences = append(machineSet.OwnerReferences, metav1.OwnerReference{
			APIVersion:         cluster.APIVersion,
			Kind:               cluster.Kind,
			Name:               cluster.Name,
			UID:                cluster.UID,
			BlockOwnerDeletion: &blockOwnerDeletion,
		})
	}

	// Add foregroundDeletion finalizer if MachineSet isn't deleted and linked to a cluster.
	if cluster != nil &&
		machineSet.ObjectMeta.DeletionTimestamp.IsZero() &&
		!util.Contains(machineSet.Finalizers, metav1.FinalizerDeleteDependents) {

		machineSet.Finalizers = append(machineSet.ObjectMeta.Finalizers, metav1.FinalizerDeleteDependents)

		if err := r.Client.Update(context.Background(), machineSet); err != nil {
			klog.Infof("Failed to add finalizers to MachineSet %q: %v", machineSet.Name, err)
			return reconcile.Result{}, err
		}

		// Since adding the finalizer updates the object return to avoid later update issues
		return reconcile.Result{Requeue: true}, nil
	}

	// Filter out irrelevant machines (deleting/mismatch labels) and claim orphaned machines.
	filteredMachines := make([]*clusterv1alpha1.Machine, 0, len(allMachines.Items))
	for idx := range allMachines.Items {
		machine := &allMachines.Items[idx]
		if shouldExcludeMachine(machineSet, machine) {
			continue
		}

		// Attempt to adopt machine if it meets previous conditions and it has no controller references.
		if metav1.GetControllerOf(machine) == nil {
			if err := r.adoptOrphan(machineSet, machine); err != nil {
				klog.Warningf("Failed to adopt Machine %q into MachineSet %q: %v", machine.Name, machineSet.Name, err)
				r.recorder.Eventf(machineSet, corev1.EventTypeWarning, "FailedAdopt", "Failed to adopt Machine %q: %v", machine.Name, err)
				continue
			}
			klog.Infof("Adopted Machine %q into MachineSet %q", machine.Name, machineSet.Name)
			r.recorder.Eventf(machineSet, corev1.EventTypeNormal, "SuccessfulAdopt", "Adopted Machine %q", machine.Name)
		}

		filteredMachines = append(filteredMachines, machine)
	}

	syncErr := r.syncReplicas(machineSet, filteredMachines)

	ms := machineSet.DeepCopy()
	newStatus := r.calculateStatus(ms, filteredMachines)

	// Always updates status as machines come up or die.
	updatedMS, err := updateMachineSetStatus(r.Client, machineSet, newStatus)
	if err != nil {
		if syncErr != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to sync machines: %v. failed to update machine set status", syncErr)
		}
		return reconcile.Result{}, errors.Wrap(err, "failed to update machine set status")
	}

	if syncErr != nil {
		return reconcile.Result{}, errors.Wrapf(syncErr, "failed to sync Machineset replicas")
	}

	var replicas int32
	if updatedMS.Spec.Replicas != nil {
		replicas = *updatedMS.Spec.Replicas
	}

	// Resync the MachineSet after MinReadySeconds as a last line of defense to guard against clock-skew.
	// Clock-skew is an issue as it may impact whether an available replica is counted as a ready replica.
	// A replica is available if the amount of time since last transition exceeds MinReadySeconds.
	// If there was a clock skew, checking whether the amount of time since last transition to ready state
	// exceeds MinReadySeconds could be incorrect.
	// To avoid an available replica stuck in the ready state, we force a reconcile after MinReadySeconds,
	// at which point it should confirm any available replica to be available.
	if updatedMS.Spec.MinReadySeconds > 0 &&
		updatedMS.Status.ReadyReplicas == replicas &&
		updatedMS.Status.AvailableReplicas != replicas {

		return reconcile.Result{RequeueAfter: time.Duration(updatedMS.Spec.MinReadySeconds) * time.Second}, nil
	}

	return reconcile.Result{}, nil
}

// getCluster reuturns the Cluster associated with the MachineSet, if any.
func (r *ReconcileMachineSet) getCluster(ms *clusterv1alpha1.MachineSet) (*clusterv1alpha1.Cluster, error) {
	if ms.Spec.Template.Labels[clusterv1alpha1.MachineClusterLabelName] == "" {
		klog.Infof("MachineSet %q in namespace %q doesn't specify %q label, assuming nil cluster", ms.Name, ms.Namespace, clusterv1alpha1.MachineClusterLabelName)
		return nil, nil
	}

	cluster := &clusterv1alpha1.Cluster{}
	key := client.ObjectKey{
		Namespace: ms.Namespace,
		Name:      ms.Spec.Template.Labels[clusterv1alpha1.MachineClusterLabelName],
	}

	if err := r.Client.Get(context.Background(), key, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}

// syncReplicas scales Machine resources up or down.
func (r *ReconcileMachineSet) syncReplicas(ms *clusterv1alpha1.MachineSet, machines []*clusterv1alpha1.Machine) error {
	if ms.Spec.Replicas == nil {
		return errors.Errorf("the Replicas field in Spec for machineset %v is nil, this should not be allowed", ms.Name)
	}

	diff := len(machines) - int(*(ms.Spec.Replicas))

	if diff < 0 {
		diff *= -1
		klog.Infof("Too few replicas for %v %s/%s, need %d, creating %d",
			controllerKind, ms.Namespace, ms.Name, *(ms.Spec.Replicas), diff)

		var machineList []*clusterv1alpha1.Machine
		var errstrings []string
		for i := 0; i < diff; i++ {
			klog.Infof("Creating machine %d of %d, ( spec.replicas(%d) > currentMachineCount(%d) )",
				i+1, diff, *(ms.Spec.Replicas), len(machines))

			machine := r.createMachine(ms)
			if err := r.Client.Create(context.Background(), machine); err != nil {
				klog.Errorf("Unable to create Machine %q: %v", machine.Name, err)
				r.recorder.Eventf(ms, corev1.EventTypeWarning, "FailedCreate", "Failed to create machine %q: %v", machine.Name, err)
				errstrings = append(errstrings, err.Error())
				continue
			}
			klog.Infof("Created machine %d of %d with name %q", i+1, diff, machine.Name)
			r.recorder.Eventf(ms, corev1.EventTypeNormal, "SuccessfulCreate", "Created machine %q", machine.Name)

			machineList = append(machineList, machine)
		}

		if len(errstrings) > 0 {
			return errors.New(strings.Join(errstrings, "; "))
		}

		return r.waitForMachineCreation(machineList)
	} else if diff > 0 {
		klog.Infof("Too many replicas for %v %s/%s, need %d, deleting %d",
			controllerKind, ms.Namespace, ms.Name, *(ms.Spec.Replicas), diff)

		deletePriorityFunc, err := getDeletePriorityFunc(ms)
		if err != nil {
			return err
		}
		klog.Infof("Found %s delete policy", ms.Spec.DeletePolicy)
		// Choose which Machines to delete.
		machinesToDelete := getMachinesToDeletePrioritized(machines, diff, deletePriorityFunc)

		// TODO: Add cap to limit concurrent delete calls.
		errCh := make(chan error, diff)
		var wg sync.WaitGroup
		wg.Add(diff)
		for _, machine := range machinesToDelete {
			go func(targetMachine *clusterv1alpha1.Machine) {
				defer wg.Done()
				err := r.Client.Delete(context.Background(), targetMachine)
				if err != nil {
					klog.Errorf("Unable to delete Machine %s: %v", targetMachine.Name, err)
					r.recorder.Eventf(ms, corev1.EventTypeWarning, "FailedDelete", "Failed to delete machine %q: %v", targetMachine.Name, err)
					errCh <- err
				}
				klog.Infof("Deleted machine %q", targetMachine.Name)
				r.recorder.Eventf(ms, corev1.EventTypeNormal, "SuccessfulDelete", "Deleted machine %q", targetMachine.Name)
			}(machine)
		}
		wg.Wait()

		select {
		case err := <-errCh:
			// all errors have been reported before and they're likely to be the same, so we'll only return the first one we hit.
			if err != nil {
				return err
			}
		default:
		}

		return r.waitForMachineDeletion(machinesToDelete)
	}

	return nil
}

// createMachine creates a Machine resource. The name of the newly created resource is going
// to be created by the API server, we set the generateName field.
func (r *ReconcileMachineSet) createMachine(machineSet *clusterv1alpha1.MachineSet) *clusterv1alpha1.Machine {
	gv := clusterv1alpha1.SchemeGroupVersion
	machine := &clusterv1alpha1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       gv.WithKind("Machine").Kind,
			APIVersion: gv.String(),
		},
		ObjectMeta: machineSet.Spec.Template.ObjectMeta,
		Spec:       machineSet.Spec.Template.Spec,
	}
	machine.ObjectMeta.GenerateName = fmt.Sprintf("%s-", machineSet.Name)
	machine.ObjectMeta.OwnerReferences = []metav1.OwnerReference{*metav1.NewControllerRef(machineSet, controllerKind)}
	machine.Namespace = machineSet.Namespace
	return machine
}

// shouldExcludeMachine returns true if the machine should be filtered out, false otherwise.
func shouldExcludeMachine(machineSet *clusterv1alpha1.MachineSet, machine *clusterv1alpha1.Machine) bool {
	// Ignore inactive machines.
	if metav1.GetControllerOf(machine) != nil && !metav1.IsControlledBy(machine, machineSet) {
		klog.V(4).Infof("%s not controlled by %v", machine.Name, machineSet.Name)
		return true
	}

	if machine.ObjectMeta.DeletionTimestamp != nil {
		return true
	}

	if !hasMatchingLabels(machineSet, machine) {
		return true
	}

	return false
}

// adoptOrphan sets the MachineSet as a controller OwnerReference to the Machine.
func (r *ReconcileMachineSet) adoptOrphan(machineSet *clusterv1alpha1.MachineSet, machine *clusterv1alpha1.Machine) error {
	newRef := *metav1.NewControllerRef(machineSet, controllerKind)
	machine.OwnerReferences = append(machine.OwnerReferences, newRef)
	return r.Client.Update(context.Background(), machine)
}

func (r *ReconcileMachineSet) waitForMachineCreation(machineList []*clusterv1alpha1.Machine) error {
	for _, machine := range machineList {
		pollErr := util.PollImmediate(stateConfirmationInterval, stateConfirmationTimeout, func() (bool, error) {
			key := client.ObjectKey{Namespace: machine.Namespace, Name: machine.Name}

			if err := r.Client.Get(context.Background(), key, &clusterv1alpha1.Machine{}); err != nil {
				if apierrors.IsNotFound(err) {
					return false, nil
				}
				klog.Error(err)
				return false, err
			}

			return true, nil
		})

		if pollErr != nil {
			klog.Error(pollErr)
			return errors.Wrap(pollErr, "failed waiting for machine object to be created")
		}
	}

	return nil
}

func (r *ReconcileMachineSet) waitForMachineDeletion(machineList []*clusterv1alpha1.Machine) error {
	for _, machine := range machineList {
		pollErr := util.PollImmediate(stateConfirmationInterval, stateConfirmationTimeout, func() (bool, error) {
			m := &clusterv1alpha1.Machine{}
			key := client.ObjectKey{Namespace: machine.Namespace, Name: machine.Name}

			err := r.Client.Get(context.Background(), key, m)
			if apierrors.IsNotFound(err) || !m.DeletionTimestamp.IsZero() {
				return true, nil
			}

			return false, err
		})

		if pollErr != nil {
			klog.Error(pollErr)
			return errors.Wrap(pollErr, "failed waiting for machine object to be deleted")
		}
	}
	return nil
}

// MachineToMachineSets is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// for MachineSets that might adopt an orphaned Machine.
func (r *ReconcileMachineSet) MachineToMachineSets(o handler.MapObject) []reconcile.Request {
	result := []reconcile.Request{}

	m := &clusterv1alpha1.Machine{}
	key := client.ObjectKey{Namespace: o.Meta.GetNamespace(), Name: o.Meta.GetName()}
	if err := r.Client.Get(context.Background(), key, m); err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("Unable to retrieve Machine %q for possible MachineSet adoption: %v", key, err)
		}
		return nil
	}

	// Check if the controller reference is already set and
	// return an empty result when one is found.
	for _, ref := range m.ObjectMeta.OwnerReferences {
		if ref.Controller != nil && *ref.Controller {
			return result
		}
	}

	mss := r.getMachineSetsForMachine(m)
	if len(mss) == 0 {
		klog.V(4).Infof("Found no MachineSet for Machine %q", m.Name)
		return nil
	}

	for _, ms := range mss {
		name := client.ObjectKey{Namespace: ms.Namespace, Name: ms.Name}
		result = append(result, reconcile.Request{NamespacedName: name})
	}

	return result
}

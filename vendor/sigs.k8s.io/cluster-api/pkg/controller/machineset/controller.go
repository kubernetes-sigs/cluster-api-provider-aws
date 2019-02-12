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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
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

var controllerKind = clusterv1alpha1.SchemeGroupVersion.WithKind("MachineSet")

// stateConfirmationTimeout is the amount of time allowed to wait for desired state.
var stateConfirmationTimeout = 10 * time.Second

// stateConfirmationInterval is the amount of time between polling for the desired state.
// The polling is against a local memory cache.
var stateConfirmationInterval = 100 * time.Millisecond

// Add creates a new MachineSet Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	r := newReconciler(mgr)
	return add(mgr, r, r.MachineToMachineSets)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileMachineSet {
	return &ReconcileMachineSet{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler, mapFn handler.ToRequestsFunc) error {
	// Create a new controller
	c, err := controller.New("machineset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to MachineSet
	err = c.Watch(&source.Kind{Type: &clusterv1alpha1.MachineSet{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Map Machine changes to MachineSets using ControllerRef
	err = c.Watch(
		&source.Kind{Type: &clusterv1alpha1.Machine{}},
		&handler.EnqueueRequestForOwner{IsController: true, OwnerType: &clusterv1alpha1.MachineSet{}},
	)
	if err != nil {
		return err
	}

	// Map Machine changes to MachineSets by machining labels
	err = c.Watch(
		&source.Kind{Type: &clusterv1alpha1.Machine{}},
		&handler.EnqueueRequestsFromMapFunc{ToRequests: mapFn})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMachineSet{}

// ReconcileMachineSet reconciles a MachineSet object
type ReconcileMachineSet struct {
	client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileMachineSet) MachineToMachineSets(o handler.MapObject) []reconcile.Request {
	result := []reconcile.Request{}
	m := &clusterv1alpha1.Machine{}
	key := client.ObjectKey{Namespace: o.Meta.GetNamespace(), Name: o.Meta.GetName()}
	err := r.Client.Get(context.Background(), key, m)
	if err != nil {
		klog.Errorf("Unable to retrieve Machine %v from store: %v", key, err)
		return nil
	}

	for _, ref := range m.ObjectMeta.OwnerReferences {
		if ref.Controller != nil && *ref.Controller {
			return result
		}
	}

	mss := r.getMachineSetsForMachine(m)
	if len(mss) == 0 {
		klog.V(4).Infof("Found no machine set for machine: %v", m.Name)
		return nil
	}

	for _, ms := range mss {
		result = append(result, reconcile.Request{
			NamespacedName: client.ObjectKey{Namespace: ms.Namespace, Name: ms.Name}})
	}

	return result
}

// Reconcile reads that state of the cluster for a MachineSet object and makes changes based on the state read
// and what is in the MachineSet.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=cluster.k8s.io,resources=machinesets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.k8s.io,resources=machines,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileMachineSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the MachineSet instance
	machineSet := &clusterv1alpha1.MachineSet{}
	err := r.Get(context.TODO(), request.NamespacedName, machineSet)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	klog.V(4).Infof("Reconcile machineset %v", machineSet.Name)
	allMachines := &clusterv1alpha1.MachineList{}

	err = r.Client.List(context.Background(), client.InNamespace(machineSet.Namespace), allMachines)
	if err != nil {
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

	// Filter out irrelevant machines (deleting/mismatch labels) and claim orphaned machines.
	filteredMachines := make([]*clusterv1alpha1.Machine, 0, len(allMachines.Items))
	for idx := range allMachines.Items {
		machine := &allMachines.Items[idx]
		if shouldExcludeMachine(machineSet, machine) {
			continue
		}
		// Attempt to adopt machine if it meets previous conditions and it has no controller ref.
		if metav1.GetControllerOf(machine) == nil {
			if err := r.adoptOrphan(machineSet, machine); err != nil {
				klog.Warningf("failed to adopt machine %v into machineset %v. %v", machine.Name, machineSet.Name, err)
				continue
			}
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
	if syncErr == nil && updatedMS.Spec.MinReadySeconds > 0 &&
		updatedMS.Status.ReadyReplicas == replicas &&
		updatedMS.Status.AvailableReplicas != replicas {

		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

// syncReplicas essentially scales machine resources up and down.
func (r *ReconcileMachineSet) syncReplicas(ms *clusterv1alpha1.MachineSet, machines []*clusterv1alpha1.Machine) error {
	if ms.Spec.Replicas == nil {
		return errors.Errorf("the Replicas field in Spec for machineset %v is nil, this should not be allowed", ms.Name)
	}
	diff := len(machines) - int(*(ms.Spec.Replicas))

	if diff < 0 {
		diff *= -1
		klog.Infof("Too few replicas for %v %s/%s, need %d, creating %d", controllerKind, ms.Namespace, ms.Name, *(ms.Spec.Replicas), diff)

		var machineList []*clusterv1alpha1.Machine
		var errstrings []string
		for i := 0; i < diff; i++ {
			klog.Infof("creating machine %d of %d, ( spec.replicas(%d) > currentMachineCount(%d) )", i+1, diff, *(ms.Spec.Replicas), len(machines))
			machine := r.createMachine(ms)
			err := r.Client.Create(context.Background(), machine)
			if err != nil {
				klog.Errorf("unable to create a machine = %s, due to %v", machine.Name, err)
				errstrings = append(errstrings, err.Error())
				continue
			}
			machineList = append(machineList, machine)
		}

		if len(errstrings) > 0 {
			return errors.New(strings.Join(errstrings, "; "))
		}
		return r.waitForMachineCreation(machineList)
	} else if diff > 0 {
		klog.Infof("Too many replicas for %v %s/%s, need %d, deleting %d", controllerKind, ms.Namespace, ms.Name, *(ms.Spec.Replicas), diff)

		// Choose which Machines to delete.
		machinesToDelete := getMachinesToDeletePrioritized(machines, diff, simpleDeletePriority)

		// TODO: Add cap to limit concurrent delete calls.
		errCh := make(chan error, diff)
		var wg sync.WaitGroup
		wg.Add(diff)
		for _, machine := range machinesToDelete {
			go func(targetMachine *clusterv1alpha1.Machine) {
				defer wg.Done()
				err := r.Client.Delete(context.Background(), targetMachine)
				if err != nil {
					klog.Errorf("unable to delete a machine = %s, due to %v", targetMachine.Name, err)
					errCh <- err
				}
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

// createMachine creates a machine resource.
// the name of the newly created resource is going to be created by the API server, we set the generateName field
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

func (r *ReconcileMachineSet) adoptOrphan(machineSet *clusterv1alpha1.MachineSet, machine *clusterv1alpha1.Machine) error {
	// Add controller reference.
	ownerRefs := machine.ObjectMeta.GetOwnerReferences()
	if ownerRefs == nil {
		ownerRefs = []metav1.OwnerReference{}
	}

	newRef := *metav1.NewControllerRef(machineSet, controllerKind)
	ownerRefs = append(ownerRefs, newRef)
	machine.ObjectMeta.SetOwnerReferences(ownerRefs)
	if err := r.Client.Update(context.Background(), machine); err != nil {
		klog.Warningf("Failed to update machine owner reference. %v", err)
		return err
	}
	return nil
}

func (r *ReconcileMachineSet) waitForMachineCreation(machineList []*clusterv1alpha1.Machine) error {
	for _, machine := range machineList {
		pollErr := util.PollImmediate(stateConfirmationInterval, stateConfirmationTimeout, func() (bool, error) {
			err := r.Client.Get(context.Background(),
				client.ObjectKey{Namespace: machine.Namespace, Name: machine.Name},
				&clusterv1alpha1.Machine{})
			if err == nil {
				return true, nil
			}
			klog.Error(err)
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
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
			err := r.Client.Get(context.Background(),
				client.ObjectKey{Namespace: machine.Namespace, Name: machine.Name},
				m)
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

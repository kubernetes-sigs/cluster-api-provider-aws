/*
Copyright 2019 The Kubernetes Authors.

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

package awsmachine

import (
	"context"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new AWSMachine Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	r := newReconciler(mgr)
	return add(mgr, r, r.MachineToProviderMachines)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileAWSMachine {
	return &ReconcileAWSMachine{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler, mapFn handler.ToRequestsFunc) error {
	// Create a new controller
	c, err := controller.New("awsmachine-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to AWSMachine
	err = c.Watch(
		&source.Kind{Type: &infrav1.AWSMachine{}},
		&handler.EnqueueRequestForObject{},
	)
	if err != nil {
		return err
	}

	return c.Watch(
		&source.Kind{Type: &clusterv1.Machine{}},
		&handler.EnqueueRequestsFromMapFunc{ToRequests: mapFn},
	)
}

var _ reconcile.Reconciler = &ReconcileAWSMachine{}

// ReconcileAWSMachine reconciles a AWSMachine object
type ReconcileAWSMachine struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AWSMachine object and makes changes based on the state read
// and what is in the AWSMachine.Spec
func (r *ReconcileAWSMachine) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	// Fetch the AWSMachine instance.
	awsm := &infrav1.AWSMachine{}
	err := r.Get(ctx, request.NamespacedName, awsm)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Machine.
	m, err := r.getMachineOwner(ctx, awsm.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	} else if m == nil {
		klog.Infof("Waiting for Machine Controller to set OwnerRef on AWSMachine %q/%q", awsm.Namespace, awsm.Name)
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Make sure bootstrap data is available and populated.
	if m.Spec.Bootstrap.Data == nil || *m.Spec.Bootstrap.Data == "" {
		klog.Infof("Waiting for bootstrap data to be available on AWSMachine %q/%q", awsm.Namespace, awsm.Name)
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	return reconcile.Result{}, nil
}

// getMachineOwner returns the Machine object owning the current resource.
func (r *ReconcileAWSMachine) getMachineOwner(ctx context.Context, meta metav1.ObjectMeta) (*clusterv1.Machine, error) {
	for _, ref := range meta.OwnerReferences {
		if ref.Kind == "Machine" && ref.APIVersion == clusterv1.SchemeGroupVersion.String() {
			m := &clusterv1.Machine{}
			key := client.ObjectKey{Name: ref.Name, Namespace: meta.Namespace}
			if err := r.Get(ctx, key, m); err != nil {
				return nil, err
			}
			return m, nil
		}
	}
	return nil, nil
}

// MachineToProviderMachines is a handler.ToRequestsFunc to be used to enqeue reconciliation
// requests for the references infrastructure provider.
func (r *ReconcileAWSMachine) MachineToProviderMachines(o handler.MapObject) []reconcile.Request {
	m, ok := o.Object.(*clusterv1.Machine)
	if !ok {
		return nil
	}

	// Return early if the api group or kind don't match what we expect.
	gvk := m.Spec.InfrastructureRef.GroupVersionKind()
	if gvk.Group != infrav1.SchemeGroupVersion.Group || gvk.Kind != "AWSMachine" {
		return nil
	}

	return []reconcile.Request{
		{
			NamespacedName: client.ObjectKey{
				Namespace: m.Namespace,
				Name:      m.Spec.InfrastructureRef.Name,
			},
		},
	}
}

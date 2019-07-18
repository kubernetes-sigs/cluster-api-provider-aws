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
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/cluster-api/pkg/controller/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/pkg/errors"
	"sigs.k8s.io/cluster-api/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	waitForClusterInfrastructureReadyDuration   = 15 * time.Second //nolint
	waitForControlPlaneMachineExistenceDuration = 5 * time.Second  //nolint
	waitForControlPlaneReadyDuration            = 5 * time.Second  //nolint
)

// Add creates a new AWSMachine Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileAWSMachine {
	return &ReconcileAWSMachine{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
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
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: util.MachineToInfrastructureMapFunc(schema.GroupVersionKind{
				Group:   infrav1.SchemeGroupVersion.Group,
				Version: infrav1.SchemeGroupVersion.Version,
				Kind:    "AWSMachine",
			}),
		},
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

	// If the AWSMachine hasn't been deleted and doesn't have a finalizer, add one.
	if awsm.ObjectMeta.DeletionTimestamp.IsZero() {
		if !util.Contains(awsm.Finalizers, clusterv1.MachineFinalizer) {
			awsm.Finalizers = append(awsm.ObjectMeta.Finalizers, clusterv1.MachineFinalizer)
		}
	}

	// Create the scope.
	scope, err := scope.NewMachineScope(scope.MachineScopeParams{
		ProviderMachine: awsm,
		Client:          r.Client,
	})
	if err != nil {
		if requeueErr, ok := errors.Cause(err).(capierrors.HasRequeueAfterError); ok {
			return reconcile.Result{RequeueAfter: requeueErr.GetRequeueAfter()}, nil
		}
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}
	defer scope.Close()

	// Make sure bootstrap data is available and populated.
	if scope.Machine.Spec.Bootstrap.Data == nil || *scope.Machine.Spec.Bootstrap.Data == "" {
		scope.Info("Waiting for bootstrap data to be available")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Call the internal reconciler.
	if err := r.reconcile(ctx, scope); err != nil {
		if requeueErr, ok := errors.Cause(err).(capierrors.HasRequeueAfterError); ok {
			scope.Info("Reconciliation asked to requeue", "after", requeueErr.GetRequeueAfter())
			return reconcile.Result{Requeue: true, RequeueAfter: requeueErr.GetRequeueAfter()}, nil
		}
		scope.Error(err, "Failed to reconcile AWSMachine")
		return reconcile.Result{}, err
	}

	// TODO(vincepri): Handle deletion and finalizers here.

	return reconcile.Result{}, nil
}

func (r *ReconcileAWSMachine) reconcile(ctx context.Context, scope *scope.MachineScope) error {
	// If the AWSMachine is in an error state, return early.
	if scope.ProviderMachine.Status.ErrorReason != nil || scope.ProviderMachine.Status.ErrorMessage != nil {
		scope.Info("Error state detected, skipping reconciliation")
		return nil
	}

	if scope.Parent.Cluster.Annotations[infrav1.AnnotationClusterInfrastructureReady] != infrav1.ValueReady {
		scope.Info("Cluster infrastructure is not ready yet, requeuing machine")
		return &capierrors.RequeueAfterError{RequeueAfter: waitForClusterInfrastructureReadyDuration}
	}

	ec2svc := ec2.NewService(scope.Parent)

	// Get or create the instance.
	instance, err := r.getOrCreate(scope, ec2svc)
	if err != nil {
		return err
	}

	// Set an error message if we couldn't find the Machine.
	if instance == nil {
		scope.SetErrorReason(common.UpdateMachineError)
		scope.SetErrorMessage(errors.New("EC2 instance cannot be found"))
		return nil
	}

	// Make sure Spec.ProviderID is always set.
	scope.SetProviderID(fmt.Sprintf("aws:////%s", instance.ID))

	// Proceed to reconcile the AWSMachine state.
	scope.SetInstanceState(instance.State)

	// TODO(vincepri): Remove this annotation when clusterctl is no longer relevant.
	scope.SetAnnotation("cluster-api-provider-aws", "true")

	switch instance.State {
	case infrav1.InstanceStateRunning:
		scope.Info("Machine instance is running", "instance-id", *scope.GetInstanceID())
	case infrav1.InstanceStatePending:
		scope.Info("Machine instance is pending", "instance-id", *scope.GetInstanceID())
	default:
		scope.SetErrorReason(common.UpdateMachineError)
		scope.SetErrorMessage(errors.Errorf("EC2 instance state %q is unexpected", instance.State))
	}

	if err := r.reconcileLBAttachment(scope, instance); err != nil {
		return errors.Errorf("failed to reconcile LB attachment: %+v", err)
	}

	// Check if we can perform any in-place updates.
	return r.update(scope, instance)
}

// findInstance queries the EC2 apis and retrieves the instance if it exists, returns nil otherwise.
func (r *ReconcileAWSMachine) findInstance(scope *scope.MachineScope, ec2svc *ec2.Service) (*infrav1.Instance, error) {
	// Parse the ProviderID.
	pid, err := noderefutil.NewProviderID(scope.GetProviderID())
	if err != nil && err != noderefutil.ErrEmptyProviderID {
		return nil, errors.Wrapf(err, "failed to parse Spec.ProviderID")
	}

	// If the ProviderID is populated, describe the instance using the ID.
	if err == nil {
		instance, err := ec2svc.InstanceIfExists(pointer.StringPtr(pid.ID()))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to query AWSMachine instance")
		}
		return instance, nil
	}

	// If the ProviderID is empty, first try to query the instance using tags,
	// otherwise create a new one.
	instance, err := ec2svc.GetRunningInstanceByTags(scope)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query AWSMachine instance by tags")
	}
	return instance, nil
}

func (r *ReconcileAWSMachine) getOrCreate(scope *scope.MachineScope, ec2svc *ec2.Service) (*infrav1.Instance, error) {
	instance, err := r.findInstance(scope, ec2svc)
	if err != nil {
		return nil, err
	}

	if instance == nil {
		// Create a new AWSMachine instance if we couldn't find a running instance.
		instance, err = ec2svc.CreateInstance(scope)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create AWSMachine instance")
		}
	}

	return instance, nil
}

func (r *ReconcileAWSMachine) update(scope *scope.MachineScope, instance *infrav1.Instance) error {
	ec2svc := ec2.NewService(scope.Parent)

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	if errs := r.isMachineOutdated(&scope.ProviderMachine.Spec, instance); len(errs) > 0 {
		return errors.Errorf("found attempt to change immutable state for machine %q: %+q", scope.Name(), errs)
	}

	existingSecurityGroups, err := ec2svc.GetInstanceSecurityGroups(*scope.GetInstanceID())
	if err != nil {
		return err
	}

	// Ensure that the security groups are correct.
	_, err = r.ensureSecurityGroups(
		ec2svc,
		scope,
		scope.ProviderMachine.Spec.AdditionalSecurityGroups,
		existingSecurityGroups,
	)
	if err != nil {
		return errors.Errorf("failed to apply security groups: %+v", err)
	}

	// Ensure that the tags are correct.
	_, err = r.ensureTags(
		ec2svc,
		scope.ProviderMachine,
		scope.GetInstanceID(),
		scope.ProviderMachine.Spec.AdditionalTags,
	)
	if err != nil {
		return errors.Errorf("failed to ensure tags: %+v", err)
	}

	return nil
}

func (r *ReconcileAWSMachine) reconcileLBAttachment(scope *scope.MachineScope, i *infrav1.Instance) error {
	if !scope.IsControlPlane() {
		return nil
	}

	elbsvc := elb.NewService(scope.Parent)
	if err := elbsvc.RegisterInstanceWithAPIServerELB(i.ID); err != nil {
		return errors.Wrapf(err, "could not register control plane instance %q with load balancer", i.ID)
	}
	return nil
}

// isMachineOudated checks that no immutable fields have been updated in an
// Update request.
// Returns a slice of errors representing attempts to change immutable state
func (r *ReconcileAWSMachine) isMachineOutdated(spec *infrav1.AWSMachineSpec, i *infrav1.Instance) (errs []error) {
	// Instance Type
	if spec.InstanceType != i.Type {
		errs = append(errs, errors.Errorf("instance type cannot be mutated from %q to %q", i.Type, spec.InstanceType))
	}

	// IAM Profile
	if spec.IAMInstanceProfile != i.IAMProfile {
		errs = append(errs, errors.Errorf("instance IAM profile cannot be mutated from %q to %q", i.IAMProfile, spec.IAMInstanceProfile))
	}

	// SSH Key Name
	if spec.KeyName != aws.StringValue(i.KeyName) {
		errs = append(errs, errors.Errorf("SSH key name cannot be mutated from %q to %q", aws.StringValue(i.KeyName), spec.KeyName))
	}

	// Root Device Size
	if spec.RootDeviceSize > 0 && spec.RootDeviceSize != i.RootDeviceSize {
		errs = append(errs, errors.Errorf("Root volume size cannot be mutated from %v to %v", i.RootDeviceSize, spec.RootDeviceSize))
	}

	// Subnet ID
	// spec.Subnet is a *AWSResourceReference and could technically be
	// a *string, ARN or Filter. However, elsewhere in the code it is only used
	// as a *string, so do the same here.
	if spec.Subnet != nil {
		if aws.StringValue(spec.Subnet.ID) != i.SubnetID {
			errs = append(errs, errors.Errorf("machine subnet ID cannot be mutated from %q to %q",
				i.SubnetID, aws.StringValue(spec.Subnet.ID)))
		}
	}

	// PublicIP check is a little more complicated as the machineConfig is a
	// simple bool indicating if the instance should have a public IP or not,
	// while the instanceDescription contains the public IP assigned to the
	// instance.
	// Work out whether the instance already has a public IP or not based on
	// the length of the PublicIP string. Anything >0 is assumed to mean it does
	// have a public IP.
	instanceHasPublicIP := false
	if len(aws.StringValue(i.PublicIP)) > 0 {
		instanceHasPublicIP = true
	}

	if aws.BoolValue(spec.PublicIP) != instanceHasPublicIP {
		errs = append(errs, errors.Errorf(`public IP setting cannot be mutated from "%v" to "%v"`,
			instanceHasPublicIP, aws.BoolValue(spec.PublicIP)))
	}

	return errs
}

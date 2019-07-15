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
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/tokens"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/cluster-api/pkg/controller/remote"
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
	defaultTokenTTL                             = 10 * time.Minute
	waitForClusterInfrastructureReadyDuration   = 15 * time.Second
	waitForControlPlaneMachineExistenceDuration = 5 * time.Second
	waitForControlPlaneReadyDuration            = 5 * time.Second
)

// Add creates a new AWSMachine Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	r := newReconciler(mgr)
	return add(mgr, r, r.MachineToProviderMachines)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileAWSMachine {
	coreClient, err := corev1.NewForConfig(mgr.GetConfig())
	if err != nil {
		klog.Fatalf("Failed to create corev1 client from configuration: %v", err)
	}

	return &ReconcileAWSMachine{
		Client:                 mgr.GetClient(),
		scheme:                 mgr.GetScheme(),
		coreClient:             coreClient,
		controlPlaneInitLocker: newControlPlaneInitLocker(coreClient),
	}
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

	// schema.GroupVersionKind{
	// 	Group:   infrav1.SchemeGroupVersion.Group,
	// 	Version: infrav1.SchemeGroupVersion.Version,
	// 	Kind:    "AWSMachine",
	// }

	return c.Watch(
		&source.Kind{Type: &clusterv1.Machine{}},
		&handler.EnqueueRequestsFromMapFunc{ToRequests: mapFn},
	)
}

var _ reconcile.Reconciler = &ReconcileAWSMachine{}

// ReconcileAWSMachine reconciles a AWSMachine object
type ReconcileAWSMachine struct {
	client.Client
	scheme                 *runtime.Scheme
	coreClient             corev1.CoreV1Interface
	controlPlaneInitLocker ControlPlaneInitLocker
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

	// Store Machine early state to allow patching.
	patchFrom := client.MergeFrom(awsm.DeepCopy())

	// If the Machine hasn't been deleted and doesn't have a finalizer, add one.
	if awsm.ObjectMeta.DeletionTimestamp.IsZero() {
		if !util.Contains(awsm.Finalizers, clusterv1.MachineFinalizer) {
			awsm.Finalizers = append(awsm.ObjectMeta.Finalizers, clusterv1.MachineFinalizer)
			if err := r.Client.Patch(ctx, awsm, patchFrom); err != nil {
				return reconcile.Result{}, errors.Wrapf(err, "failed to add finalizer to AWSMachine %q/%q", awsm.Namespace, awsm.Name)
			}
			// Since adding the finalizer updates the object return to avoid later update issues
			return reconcile.Result{Requeue: true}, nil
		}
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

	scope, err := actuators.NewMachineScope(actuators.MachineScopeParams{
		ProviderMachine: awsm,
		Machine:         m,
		Client:          r.Client,
	})
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}
	defer scope.Close()

	return reconcile.Result{}, nil
}

func (r *ReconcileAWSMachine) reconcile(ctx context.Context, scope *actuators.MachineScope) error {
	exist, err := r.exists(scope)
	if err != nil {
		klog.Errorf("Failed to check if Machine %q infrastructure exists: %v", scope.Name(), err)
		return err
	}

	// Reconcile ProviderID.
	if scope.ProviderMachine.Spec.ProviderID == nil || *scope.ProviderMachine.Spec.ProviderID == "" {
		providerID := fmt.Sprintf("aws:////%s", *scope.ProviderMachine.Status.InstanceID)
		scope.ProviderMachine.Spec.ProviderID = &providerID
	}

	if exist {
		scope.Info("Reconciling Machine triggers idempotent update")
		return r.update(scope)
	}

	// Machine resource created. Machine does not yet exist.
	scope.Info("Reconciling Machine triggers idempotent create")
	return r.create(scope)
}

// create creates a machine and is invoked by the machine controller.
func (r *ReconcileAWSMachine) create(scope *actuators.MachineScope) error {
	if scope.Cluster.Cluster.Annotations[infrav1.AnnotationClusterInfrastructureReady] != infrav1.ValueReady {
		scope.Info("Cluster infrastructure is not ready yet - requeuing machine")
		return &capierrors.RequeueAfterError{RequeueAfter: waitForClusterInfrastructureReadyDuration}
	}

	ec2svc := ec2.NewService(scope.Cluster)

	scope.Info("Retrieving machines for cluster")
	machineList := &clusterv1.MachineList{}
	if err := scope.Client.List(context.Background(), machineList, actuators.ListOptionsForCluster(scope.Cluster.Name())); err != nil {
		return errors.Wrapf(err, "failed to retrieve machines in cluster %q", scope.Cluster.Name())
	}

	controlPlaneMachines := util.GetControlPlaneMachinesFromList(machineList)
	if len(controlPlaneMachines) == 0 {
		scope.Info("No control plane machines exist yet - requeuing")
		return &capierrors.RequeueAfterError{RequeueAfter: waitForControlPlaneMachineExistenceDuration}
	}

	join, err := r.isNodeJoin(scope)
	if err != nil {
		return err
	}

	var bootstrapToken string
	if join {
		remoteClient, err := remote.NewClusterClient(r.Client, scope.Cluster.Cluster)
		if err != nil {
			return errors.Wrapf(err, "unable to proceed until control plane is ready (error creating client) for cluster %q",
				path.Join(scope.Cluster.Namespace(), scope.Cluster.Name()))
		}

		coreClient, err := remoteClient.CoreV1()
		if err != nil {
			return errors.Wrapf(err, "unable to proceed until control plane is ready (error creating client) for cluster %q",
				path.Join(scope.Cluster.Namespace(), scope.Cluster.Name()))
		}

		scope.Info("Machine will join the cluster")
		bootstrapToken, err = tokens.NewBootstrap(coreClient, defaultTokenTTL)
		if err != nil {
			return errors.Wrapf(err, "failed to create new bootstrap token")
		}
	} else {
		scope.Info("Machine will init the cluster")
	}

	i, err := ec2svc.CreateOrGetMachine(scope, bootstrapToken)
	if err != nil {
		return errors.Errorf("failed to create or get machine: %+v", err)
	}

	scope.ProviderMachine.Status.InstanceID = &i.ID
	scope.ProviderMachine.Status.InstanceState = &i.State

	if scope.ProviderMachine.Annotations == nil {
		scope.ProviderMachine.Annotations = map[string]string{}
	}
	scope.ProviderMachine.Annotations["cluster-api-provider-aws"] = "true"

	if err := r.reconcileLBAttachment(scope, i); err != nil {
		return errors.Errorf("failed to reconcile LB attachment: %+v", err)
	}

	scope.Info("Create completed")
	return nil
}

func (r *ReconcileAWSMachine) exists(scope *actuators.MachineScope) (bool, error) {
	ec2svc := ec2.NewService(scope.Cluster)

	// TODO worry about pointers. instance if exists returns *any* instance
	if scope.ProviderMachine.Status.InstanceID == nil {
		return false, nil
	}

	instance, err := ec2svc.InstanceIfExists(scope.ProviderMachine.Status.InstanceID)
	if err != nil {
		return false, errors.Errorf("failed to retrieve instance: %+v", err)
	}

	if instance == nil {
		return false, nil
	}

	scope.Info("Found instance for machine", "instance", instance)

	switch instance.State {
	case infrav1.InstanceStateRunning:
		scope.Info("Machine instance is running", "instance-id", *scope.ProviderMachine.Status.InstanceID)
	case infrav1.InstanceStatePending:
		scope.Info("Machine instance is pending", "instance-id", *scope.ProviderMachine.Status.InstanceID)
	default:
		return false, nil
	}

	scope.ProviderMachine.Status.InstanceState = &instance.State

	if err := r.reconcileLBAttachment(scope, instance); err != nil {
		return true, err
	}

	return true, nil
}

func (r *ReconcileAWSMachine) update(scope *actuators.MachineScope) error {
	ec2svc := ec2.NewService(scope.Cluster)

	// Get the current instance description from AWS.
	instanceDescription, err := ec2svc.InstanceIfExists(scope.ProviderMachine.Status.InstanceID)
	if err != nil {
		return errors.Errorf("failed to get instance: %+v", err)
	}

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	if errs := r.isMachineOutdated(&scope.ProviderMachine.Spec, instanceDescription); len(errs) > 0 {
		return errors.Errorf("found attempt to change immutable state for machine %q: %+q", scope.Name(), errs)
	}

	existingSecurityGroups, err := ec2svc.GetInstanceSecurityGroups(*scope.ProviderMachine.Status.InstanceID)
	if err != nil {
		return err
	}

	// Ensure that the security groups are correct.
	_, err = r.ensureSecurityGroups(
		ec2svc,
		scope,
		*scope.ProviderMachine.Status.InstanceID,
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
		scope.ProviderMachine.Status.InstanceID,
		scope.ProviderMachine.Spec.AdditionalTags,
	)
	if err != nil {
		return errors.Errorf("failed to ensure tags: %+v", err)
	}

	return nil
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

func (r *ReconcileAWSMachine) reconcileLBAttachment(scope *actuators.MachineScope, i *infrav1.Instance) error {
	if !scope.IsControlPlane() {
		return nil
	}

	elbsvc := elb.NewService(scope.Cluster)
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

func (r *ReconcileAWSMachine) isNodeJoin(scope *actuators.MachineScope) (bool, error) {
	if scope.Cluster.Cluster.Annotations[v1alpha1.AnnotationControlPlaneReady] == v1alpha1.ValueReady {
		return true, nil
	}

	if !util.IsControlPlaneMachine(scope.Machine) {
		// This isn't a control plane machine - have to wait
		scope.Info("No control plane machines exist yet - requeuing")
		return true, &capierrors.RequeueAfterError{RequeueAfter: waitForControlPlaneMachineExistenceDuration}
	}

	if r.controlPlaneInitLocker.Acquire(scope.Cluster) {
		return false, nil
	}

	scope.Info("Unable to acquire control plane configmap lock - requeuing")
	return true, &capierrors.RequeueAfterError{RequeueAfter: waitForControlPlaneReadyDuration}
}

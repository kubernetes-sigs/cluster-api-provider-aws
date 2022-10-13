/*
Copyright 2020 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/controllers"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/asginstancestate"
	asg "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/autoscaling"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// AWSMachinePoolReconciler reconciles a AWSMachinePool object.
type AWSMachinePoolReconciler struct {
	client.Client
	Recorder          record.EventRecorder
	WatchFilterValue  string
	asgServiceFactory func(cloud.ClusterScoper) services.ASGInterface
	ec2ServiceFactory func(scope.EC2Scope) services.EC2Interface
}

func (r *AWSMachinePoolReconciler) getASGService(scope cloud.ClusterScoper) services.ASGInterface {
	if r.asgServiceFactory != nil {
		return r.asgServiceFactory(scope)
	}
	return asg.NewService(scope)
}

func (r *AWSMachinePoolReconciler) getEC2Service(scope scope.EC2Scope) services.EC2Interface {
	if r.ec2ServiceFactory != nil {
		return r.ec2ServiceFactory(scope)
	}

	return ec2.NewService(scope)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

// Reconcile is the reconciliation loop for AWSMachinePool.
func (r *AWSMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the AWSMachinePool .
	awsMachinePool := &expinfrav1.AWSMachinePool{}
	err := r.Get(ctx, req.NamespacedName, awsMachinePool)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the CAPI MachinePool
	machinePool, err := getOwnerMachinePool(ctx, r.Client, awsMachinePool.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if machinePool == nil {
		log.Info("MachinePool Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}
	log = log.WithValues("machinePool", klog.KObj(machinePool))

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		log.Info("MachinePool is missing cluster label or cluster does not exist")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	infraCluster, err := r.getInfraCluster(ctx, log, cluster, awsMachinePool)
	if err != nil {
		return ctrl.Result{}, errors.New("error getting infra provider cluster or control plane object")
	}
	if infraCluster == nil {
		log.Info("AWSCluster or AWSManagedControlPlane is not ready yet")
		return ctrl.Result{}, nil
	}

	// Create the machine pool scope
	machinePoolScope, err := scope.NewMachinePoolScope(scope.MachinePoolScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		MachinePool:    machinePool,
		InfraCluster:   infraCluster,
		AWSMachinePool: awsMachinePool,
	})
	if err != nil {
		log.Error(err, "failed to create scope")
		return ctrl.Result{}, err
	}

	// Always close the scope when exiting this function so we can persist any AWSMachine changes.
	defer func() {
		// set Ready condition before AWSMachinePool is patched
		conditions.SetSummary(machinePoolScope.AWSMachinePool,
			conditions.WithConditions(
				expinfrav1.ASGReadyCondition,
				expinfrav1.LaunchTemplateReadyCondition,
			),
			conditions.WithStepCounterIfOnly(
				expinfrav1.ASGReadyCondition,
				expinfrav1.LaunchTemplateReadyCondition,
			),
		)

		if err := machinePoolScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	switch infraScope := infraCluster.(type) {
	case *scope.ManagedControlPlaneScope:
		if !awsMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
			return r.reconcileDelete(machinePoolScope, infraScope, infraScope)
		}

		return r.reconcileNormal(ctx, machinePoolScope, infraScope, infraScope)
	case *scope.ClusterScope:
		if !awsMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
			return r.reconcileDelete(machinePoolScope, infraScope, infraScope)
		}

		return r.reconcileNormal(ctx, machinePoolScope, infraScope, infraScope)
	default:
		return ctrl.Result{}, errors.New("infraCluster has unknown type")
	}
}

func (r *AWSMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&expinfrav1.AWSMachinePool{}).
		Watches(
			&source.Kind{Type: &expclusterv1.MachinePool{}},
			handler.EnqueueRequestsFromMapFunc(machinePoolToInfrastructureMapFunc(expinfrav1.GroupVersion.WithKind("AWSMachinePool"))),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		Complete(r)
}

func (r *AWSMachinePoolReconciler) reconcileNormal(ctx context.Context, machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope) (ctrl.Result, error) {
	clusterScope.Info("Reconciling AWSMachinePool")

	// If the AWSMachine is in an error state, return early.
	if machinePoolScope.HasFailed() {
		machinePoolScope.Info("Error state detected, skipping reconciliation")

		// TODO: If we are in a failed state, delete the secret regardless of instance state

		return ctrl.Result{}, nil
	}

	// If the AWSMachinepool doesn't have our finalizer, add it
	controllerutil.AddFinalizer(machinePoolScope.AWSMachinePool, expinfrav1.MachinePoolFinalizer)

	// Register finalizer immediately to avoid orphaning AWS resources
	if err := machinePoolScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	if !machinePoolScope.Cluster.Status.InfrastructureReady {
		machinePoolScope.Info("Cluster infrastructure is not ready yet")
		conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated
	if machinePoolScope.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName == nil {
		machinePoolScope.Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	ec2Svc := r.getEC2Service(ec2Scope)
	asgsvc := r.getASGService(clusterScope)

	canUpdateLaunchTemplate := func() (bool, error) {
		// If there is a change: before changing the template, check if there exist an ongoing instance refresh,
		// because only 1 instance refresh can be "InProgress". If template is updated when refresh cannot be started,
		// that change will not trigger a refresh. Do not start an instance refresh if only userdata changed.
		return asgsvc.CanStartASGInstanceRefresh(machinePoolScope)
	}
	runPostLaunchTemplateUpdateOperation := func() error {
		// After creating a new version of launch template, instance refresh is required
		// to trigger a rolling replacement of all previously launched instances.
		// If ONLY the userdata changed, previously launched instances continue to use the old launch
		// template.
		//
		// FIXME(dlipovetsky,sedefsavas): If the controller terminates, or the StartASGInstanceRefresh returns an error,
		// this conditional will not evaluate to true the next reconcile. If any machines use an older
		// Launch Template version, and the difference between the older and current versions is _more_
		// than userdata, we should start an Instance Refresh.
		machinePoolScope.Info("starting instance refresh", "number of instances", machinePoolScope.MachinePool.Spec.Replicas)
		return asgsvc.StartASGInstanceRefresh(machinePoolScope)
	}
	if err := ec2Svc.ReconcileLaunchTemplate(machinePoolScope, canUpdateLaunchTemplate, runPostLaunchTemplateUpdateOperation); err != nil {
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedLaunchTemplateReconcile", "Failed to reconcile launch template: %v", err)
		machinePoolScope.Error(err, "failed to reconcile launch template")
		return ctrl.Result{}, err
	}

	// set the LaunchTemplateReady condition
	conditions.MarkTrue(machinePoolScope.AWSMachinePool, expinfrav1.LaunchTemplateReadyCondition)

	// Find existing ASG
	asg, err := r.findASG(machinePoolScope, asgsvc)
	if err != nil {
		conditions.MarkUnknown(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, expinfrav1.ASGNotFoundReason, err.Error())
		return ctrl.Result{}, err
	}

	if asg == nil {
		// Create new ASG
		if _, err := r.createPool(machinePoolScope, clusterScope); err != nil {
			conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, expinfrav1.ASGProvisionFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := r.updatePool(machinePoolScope, clusterScope, asg); err != nil {
		machinePoolScope.Error(err, "error updating AWSMachinePool")
		return ctrl.Result{}, err
	}

	launchTemplateID := machinePoolScope.GetLaunchTemplateIDStatus()
	asgName := machinePoolScope.Name()
	resourceServiceToUpdate := []scope.ResourceServiceToUpdate{
		{
			ResourceID:      &launchTemplateID,
			ResourceService: ec2Svc,
		},
		{
			ResourceID:      &asgName,
			ResourceService: asgsvc,
		},
	}
	err = ec2Svc.ReconcileTags(machinePoolScope, resourceServiceToUpdate)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "error updating tags")
	}

	// Make sure Spec.ProviderID is always set.
	machinePoolScope.AWSMachinePool.Spec.ProviderID = asg.ID
	providerIDList := make([]string, len(asg.Instances))

	instancestateSvc := asginstancestate.NewService(ec2Scope)
	if err := instancestateSvc.ReconcileASGEC2Events(); err != nil {
		// non fatal error, so we continue
		clusterScope.Error(err, "non-fatal: failed to set up EventBridge")
	}

	for i, ec2 := range asg.Instances {
		providerIDList[i] = fmt.Sprintf("aws:///%s/%s", ec2.AvailabilityZone, ec2.ID)
	}

	machinePoolScope.SetAnnotation("cluster-api-provider-aws", "true")
	r.checkForTerminatingInstancesInMachinePool(ctx, clusterScope, asgsvc, machinePoolScope.Name())
	machinePoolScope.AWSMachinePool.Spec.ProviderIDList = providerIDList
	machinePoolScope.AWSMachinePool.Status.Replicas = int32(len(providerIDList))
	machinePoolScope.AWSMachinePool.Status.Ready = true
	conditions.MarkTrue(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition)

	err = machinePoolScope.UpdateInstanceStatuses(ctx, asg.Instances)
	if err != nil {
		machinePoolScope.Info("Failed updating instances", "instances", asg.Instances)
	}

	return ctrl.Result{}, nil
}

func (r *AWSMachinePoolReconciler) reconcileDelete(machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope) (ctrl.Result, error) {
	clusterScope.Info("Handling deleted AWSMachinePool")

	ec2Svc := r.getEC2Service(ec2Scope)
	asgSvc := r.getASGService(clusterScope)

	asg, err := r.findASG(machinePoolScope, asgSvc)
	if err != nil {
		return ctrl.Result{}, err
	}

	if asg == nil {
		machinePoolScope.V(2).Info("Unable to locate ASG")
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeNormal, "NoASGFound", "Unable to find matching ASG")
	} else {
		machinePoolScope.SetASGStatus(asg.Status)
		switch asg.Status {
		case expinfrav1.ASGStatusDeleteInProgress:
			// ASG is already deleting
			machinePoolScope.SetNotReady()
			conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, expinfrav1.ASGDeletionInProgress, clusterv1.ConditionSeverityWarning, "")
			r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "DeletionInProgress", "ASG deletion in progress: %q", asg.Name)
			machinePoolScope.Info("ASG is already deleting", "name", asg.Name)
		default:
			machinePoolScope.Info("Deleting ASG", "id", asg.Name, "status", asg.Status)
			if err := asgSvc.DeleteASGAndWait(asg.Name); err != nil {
				r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedDelete", "Failed to delete ASG %q: %v", asg.Name, err)
				return ctrl.Result{}, errors.Wrap(err, "failed to delete ASG")
			}
		}
	}

	launchTemplateID := machinePoolScope.AWSMachinePool.Status.LaunchTemplateID
	launchTemplate, _, err := ec2Svc.GetLaunchTemplate(machinePoolScope.LaunchTemplateName())
	if err != nil {
		return ctrl.Result{}, err
	}

	if launchTemplate == nil {
		machinePoolScope.V(2).Info("Unable to locate launch template")
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeNormal, "NoASGFound", "Unable to find matching ASG")
		controllerutil.RemoveFinalizer(machinePoolScope.AWSMachinePool, expinfrav1.MachinePoolFinalizer)
		return ctrl.Result{}, nil
	}

	machinePoolScope.Info("deleting launch template", "name", launchTemplate.Name)
	if err := ec2Svc.DeleteLaunchTemplate(launchTemplateID); err != nil {
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedDelete", "Failed to delete launch template %q: %v", launchTemplate.Name, err)
		return ctrl.Result{}, errors.Wrap(err, "failed to delete ASG")
	}

	machinePoolScope.Info("successfully deleted AutoScalingGroup and Launch Template")

	// remove finalizer
	controllerutil.RemoveFinalizer(machinePoolScope.AWSMachinePool, expinfrav1.MachinePoolFinalizer)

	instancestateSvc := asginstancestate.NewService(ec2Scope)
	if err := instancestateSvc.DeleteASGEC2Events(); err != nil {
		clusterScope.Error(err, "non-fatal: failed to delete EventBridge notifications")
	}

	return ctrl.Result{}, nil
}

func (r *AWSMachinePoolReconciler) updatePool(machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper, existingASG *expinfrav1.AutoScalingGroup) error {
	asgSvc := r.getASGService(clusterScope)
	if asgNeedsUpdates(machinePoolScope, existingASG) {
		machinePoolScope.Info("updating AutoScalingGroup")

		if err := asgSvc.UpdateASG(machinePoolScope); err != nil {
			r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedUpdate", "Failed to update ASG: %v", err)
			return errors.Wrap(err, "unable to update ASG")
		}
	}

	suspendedProcessesSlice := machinePoolScope.AWSMachinePool.Spec.SuspendProcesses.ConvertSetValuesToStringSlice()
	if !cmp.Equal(existingASG.CurrentlySuspendProcesses, suspendedProcessesSlice) {
		var (
			toBeSuspended []string
			toBeResumed   []string

			currentlySuspended = make(map[string]struct{})
			desiredSuspended   = make(map[string]struct{})
		)

		// Convert the items to a map, so it's easy to create an effective diff from these two slices.
		for _, p := range existingASG.CurrentlySuspendProcesses {
			currentlySuspended[p] = struct{}{}
		}

		for _, p := range suspendedProcessesSlice {
			desiredSuspended[p] = struct{}{}
		}

		// Anything that remains in the desired items is not currently suspended so must be suspended.
		// Anything that remains in the currentlySuspended list must be resumed since they were not part of
		// desiredSuspended.
		for k := range desiredSuspended {
			if _, ok := currentlySuspended[k]; ok {
				delete(desiredSuspended, k)
			}
			delete(currentlySuspended, k)
		}

		// Convert them back into lists so
		for k := range desiredSuspended {
			toBeSuspended = append(toBeSuspended, k)
		}

		for k := range currentlySuspended {
			toBeResumed = append(toBeResumed, k)
		}

		if len(toBeSuspended) > 0 {
			if err := asgSvc.SuspendProcesses(existingASG.Name, toBeSuspended); err != nil {
				return errors.Wrapf(err, "failed to suspend processes while trying update pool")
			}
		}
		if len(toBeResumed) > 0 {
			if err := asgSvc.ResumeProcesses(existingASG.Name, toBeResumed); err != nil {
				return errors.Wrapf(err, "failed to resume processes while trying update pool")
			}
		}
	}
	return nil
}

func (r *AWSMachinePoolReconciler) createPool(machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper) (*expinfrav1.AutoScalingGroup, error) {
	clusterScope.Info("Initializing ASG client")

	asgsvc := r.getASGService(clusterScope)

	machinePoolScope.Info("Creating Autoscaling Group")
	asg, err := asgsvc.CreateASG(machinePoolScope)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create AWSMachinePool")
	}
	suspendedProcessesSlice := machinePoolScope.AWSMachinePool.Spec.SuspendProcesses.ConvertSetValuesToStringSlice()
	if err := asgsvc.SuspendProcesses(asg.Name, suspendedProcessesSlice); err != nil {
		return nil, errors.Wrapf(err, "failed to suspend processes while trying to create Pool")
	}
	return asg, nil
}

func (r *AWSMachinePoolReconciler) findASG(machinePoolScope *scope.MachinePoolScope, asgsvc services.ASGInterface) (*expinfrav1.AutoScalingGroup, error) {
	// Query the instance using tags.
	asg, err := asgsvc.GetASGByName(machinePoolScope)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query AWSMachinePool by name")
	}

	return asg, nil
}

// asgNeedsUpdates compares incoming AWSMachinePool and compares against existing ASG.
func asgNeedsUpdates(machinePoolScope *scope.MachinePoolScope, existingASG *expinfrav1.AutoScalingGroup) bool {
	if machinePoolScope.MachinePool.Spec.Replicas != nil {
		if existingASG.DesiredCapacity == nil || *machinePoolScope.MachinePool.Spec.Replicas != *existingASG.DesiredCapacity {
			return true
		}
	} else if existingASG.DesiredCapacity != nil {
		return true
	}

	if machinePoolScope.AWSMachinePool.Spec.MaxSize != existingASG.MaxSize {
		return true
	}

	if machinePoolScope.AWSMachinePool.Spec.MinSize != existingASG.MinSize {
		return true
	}

	if machinePoolScope.AWSMachinePool.Spec.CapacityRebalance != existingASG.CapacityRebalance {
		return true
	}

	if !cmp.Equal(machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy, existingASG.MixedInstancesPolicy) {
		machinePoolScope.Info("got a mixed diff here", "incoming", machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy, "existing", existingASG.MixedInstancesPolicy)
		return true
	}

	// todo subnet diff

	return false
}

// getOwnerMachinePool returns the MachinePool object owning the current resource.
func getOwnerMachinePool(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*expclusterv1.MachinePool, error) {
	for _, ref := range obj.OwnerReferences {
		if ref.Kind != "MachinePool" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if gv.Group == expclusterv1.GroupVersion.Group {
			return getMachinePoolByName(ctx, c, obj.Namespace, ref.Name)
		}
	}
	return nil, nil
}

// getMachinePoolByName finds and return a Machine object using the specified params.
func getMachinePoolByName(ctx context.Context, c client.Client, namespace, name string) (*expclusterv1.MachinePool, error) {
	m := &expclusterv1.MachinePool{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, m); err != nil {
		return nil, err
	}
	return m, nil
}

func machinePoolToInfrastructureMapFunc(gvk schema.GroupVersionKind) handler.MapFunc {
	return func(o client.Object) []reconcile.Request {
		m, ok := o.(*expclusterv1.MachinePool)
		if !ok {
			panic(fmt.Sprintf("Expected a MachinePool but got a %T", o))
		}

		gk := gvk.GroupKind()
		// Return early if the GroupKind doesn't match what we expect
		infraGK := m.Spec.Template.Spec.InfrastructureRef.GroupVersionKind().GroupKind()
		if gk != infraGK {
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: client.ObjectKey{
					Namespace: m.Namespace,
					Name:      m.Spec.Template.Spec.InfrastructureRef.Name,
				},
			},
		}
	}
}

func (r *AWSMachinePoolReconciler) getInfraCluster(ctx context.Context, log logr.Logger, cluster *clusterv1.Cluster, awsMachinePool *expinfrav1.AWSMachinePool) (scope.EC2Scope, error) {
	var clusterScope *scope.ClusterScope
	var managedControlPlaneScope *scope.ManagedControlPlaneScope
	var err error

	if cluster.Spec.ControlPlaneRef != nil && cluster.Spec.ControlPlaneRef.Kind == controllers.AWSManagedControlPlaneRefKind {
		controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
		controlPlaneName := client.ObjectKey{
			Namespace: awsMachinePool.Namespace,
			Name:      cluster.Spec.ControlPlaneRef.Name,
		}

		if err := r.Get(ctx, controlPlaneName, controlPlane); err != nil {
			// AWSManagedControlPlane is not ready
			return nil, nil //nolint:nilerr
		}

		managedControlPlaneScope, err = scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
			Client:         r.Client,
			Logger:         &log,
			Cluster:        cluster,
			ControlPlane:   controlPlane,
			ControllerName: "awsManagedControlPlane",
		})
		if err != nil {
			return nil, err
		}

		return managedControlPlaneScope, nil
	}

	awsCluster := &infrav1.AWSCluster{}

	infraClusterName := client.ObjectKey{
		Namespace: awsMachinePool.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}

	if err := r.Client.Get(ctx, infraClusterName, awsCluster); err != nil {
		// AWSCluster is not ready
		return nil, nil //nolint:nilerr
	}

	// Create the cluster scope
	clusterScope, err = scope.NewClusterScope(scope.ClusterScopeParams{
		Client:         r.Client,
		Logger:         &log,
		Cluster:        cluster,
		AWSCluster:     awsCluster,
		ControllerName: "awsmachine",
	})
	if err != nil {
		return nil, err
	}

	return clusterScope, nil
}

func (r *AWSMachinePoolReconciler) checkForTerminatingInstancesInMachinePool(ctx context.Context, clusterScope cloud.ClusterScoper, asgsvc services.ASGInterface, asgName string) {

	sqsSvs := scope.NewGlobalSQSClient(clusterScope, clusterScope)
	queueURL, _ := r.getQueueURL(clusterScope.Name(), sqsSvs)
	resp, err := sqsSvs.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: aws.String(queueURL)})
	if err != nil {
		clusterScope.Error(err, "non-fatal: failed to receive messages from instance state queue")
		return
	}
	for _, msg := range resp.Messages {

		m := message{}
		err := json.Unmarshal([]byte(*msg.Body), &m)
		if err != nil {
			clusterScope.Error(err, "non-fatal: failed to receive messages from instance state queue")
			return
		}
		if r.isInstanceTerminating(m) {
			clusterScope.Info("Instance shutting down", "instance-id", m.Detail.EC2InstanceId)
			err = r.cordonNode(ctx, clusterScope, m.Detail.EC2InstanceId)
			if err != nil {
				clusterScope.Error(err, "non-fatal: failed to cordon the node", "instance-id", m.Detail.EC2InstanceId)
			} else {
				clusterScope.Info("Node cordoned sucessfully", "instance-id", m.Detail.EC2InstanceId)
			}
			errs := r.evictPods(ctx, clusterScope, m.Detail.EC2InstanceId)
			if len(errs) > 0 {
				for i := 0; i < len(errs); i++ {
					clusterScope.Error(err, "non-fatal: failed to evict pod from the node", "instance-id", m.Detail.EC2InstanceId)
				}
			} else {
				clusterScope.Info("All pods evicted sucessfully from the node", "instance-id", m.Detail.EC2InstanceId)
			}
			err = asgsvc.CompleteLifeCycleEvent(asgName, m.Detail.EC2InstanceId)
			if err != nil {
				clusterScope.Error(err, "non-fatal: failed to remove pre deletion hook for the node", "instance-id", m.Detail.EC2InstanceId)
			} else {
				clusterScope.Info("Pre deletion hook removed sucessfully from the node", "instance-id", m.Detail.EC2InstanceId)
			}
			_, err = sqsSvs.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				clusterScope.Error(err, "error deleting message", "queueURL", queueURL, "messageReceiptHandle", msg.ReceiptHandle)
			}
		} else {
			_, err = sqsSvs.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})

			if err != nil {
				clusterScope.Error(err, "error deleting message", "queueURL", queueURL, "messageReceiptHandle", msg.ReceiptHandle)
			}
		}
	}
}

func (r *AWSMachinePoolReconciler) getQueueURL(clusterName string, sqsSvs sqsiface.SQSAPI) (string, error) {
	queueName := asginstancestate.GenerateQueueName(clusterName)
	resp, err := sqsSvs.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})
	if err != nil {
		return "", err
	}

	return *resp.QueueUrl, nil
}

func (r *AWSMachinePoolReconciler) isInstanceTerminating(msg message) bool {

	if msg.Source != "aws.autoscaling" || msg.Detail.LifecycleTransition == "" {
		return false
	}
	if msg.Detail.LifecycleTransition == LifecycleTransitionTerminating {
		return true
	}
	return false
}

func (r *AWSMachinePoolReconciler) cordonNode(ctx context.Context, clusterScope cloud.ClusterScoper, nodeInstanceID string) error {

	node, err := r.getNodeFromInstanceID(ctx, clusterScope, nodeInstanceID)
	if err != nil {
		return err
	}

	// disable schedulabling of pods on the node
	node.Spec.Unschedulable = true

	workloadClient, err := remote.NewClusterClient(ctx, "", r.Client, util.ObjectKey(clusterScope.ClusterObj()))
	if err != nil {
		return err
	}

	if err := workloadClient.Update(context.Background(), node); err != nil {
		return err
	}

	return nil
}

func (r *AWSMachinePoolReconciler) evictPods(ctx context.Context, clusterScope cloud.ClusterScoper, nodeInstanceID string) []error {
	pods, err := r.podsToEvict(ctx, clusterScope, nodeInstanceID)
	if err != nil {
		return []error{fmt.Errorf("failed to get the list of pods to evict: %v", err)}
	}
	errCh := make(chan error, len(pods))
	retErrs := []error{}

	var wg sync.WaitGroup
	var isDone bool
	defer func() { isDone = true }()

	wg.Add(len(pods))

	for _, pod := range pods {
		go func(p corev1.Pod) {
			defer wg.Done()
			for {
				if isDone {
					return
				}
				err := r.evictPod(ctx, clusterScope, &p)
				if err == nil {
					clusterScope.Info("Successfully evicted pod", "namespace", p.Namespace, "pod", p.Name, "instance-id", nodeInstanceID)
					return
				} else {
					errCh <- fmt.Errorf("error evicting pod %s/%s on node : %v", p.Namespace, p.Name, err)
					return
				}
			}
		}(pod)
	}

	finished := make(chan struct{})
	go func() { wg.Wait(); finished <- struct{}{} }()

	select {
	case <-finished:
		break
	case err := <-errCh:
		retErrs = append(retErrs, err)
	}

	return retErrs
}

func (r *AWSMachinePoolReconciler) podsToEvict(ctx context.Context, clusterScope cloud.ClusterScoper, nodeInstanceID string) ([]corev1.Pod, error) {

	var podsToEvict []corev1.Pod
	node, err := r.getNodeFromInstanceID(ctx, clusterScope, nodeInstanceID)
	if err != nil {
		return podsToEvict, nil
	}
	workloadClient, err := remote.NewClusterClient(ctx, "", r.Client, util.ObjectKey(clusterScope.ClusterObj()))
	if err != nil {
		return podsToEvict, nil
	}
	// filter to list all the pods on the node.
	listOptions := client.ListOptions{
		Namespace: metav1.NamespaceAll,
		Raw: &metav1.ListOptions{
			FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": node.Name}).String()},
	}

	listOptions.ApplyToList(&listOptions)

	podList := &corev1.PodList{}

	if err := workloadClient.List(context.Background(), podList, &listOptions); err != nil {
		return podsToEvict, nil
	}
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		if controllerRef := metav1.GetControllerOf(&pod); controllerRef != nil && controllerRef.Kind == "DaemonSet" {
			continue
		}
		if _, found := pod.ObjectMeta.Annotations[corev1.MirrorPodAnnotationKey]; found {
			continue
		}
		podsToEvict = append(podsToEvict, pod)
	}

	return podsToEvict, nil
}

func (r *AWSMachinePoolReconciler) evictPod(ctx context.Context, clusterScope cloud.ClusterScoper, pod *corev1.Pod) error {

	// get the config of the cluster where we need to evict pods
	clusterConfig, err := remote.RESTConfig(ctx, "", r.Client, util.ObjectKey(clusterScope.ClusterObj()))
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return err
	}

	eviction := &policy.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
	}

	return kubeClient.PolicyV1beta1().Evictions(eviction.Namespace).Evict(ctx, eviction)
}

func (r *AWSMachinePoolReconciler) getNodeFromInstanceID(ctx context.Context, clusterScope cloud.ClusterScoper, nodeInstanceID string) (*corev1.Node, error) {

	nodes := &corev1.NodeList{}
	workloadClient, err := remote.NewClusterClient(ctx, "", r.Client, util.ObjectKey(clusterScope.ClusterObj()))
	if err != nil {
		return &corev1.Node{}, nil
	}
	if err := workloadClient.List(context.Background(), nodes); err != nil {
		return &corev1.Node{}, nil
	}
	node := corev1.Node{}
	for _, node = range nodes.Items {
		// provider ID is of the format aws:///us-east-1c/i-016bbceabf8257d39
		if strings.Contains(node.Spec.ProviderID, nodeInstanceID) {
			break
		}
	}
	return &node, nil
}

type message struct {
	Source     string         `json:"source"`
	DetailType string         `json:"detail-type,omitempty"`
	Detail     *messageDetail `json:"detail,omitempty"`
}

type messageDetail struct {
	EC2InstanceId       string              `json:"EC2InstanceId,omitempty"`
	LifecycleTransition LifecycleTransition `json:"LifecycleTransition"`
}

// LifecycleTransition describes the state of an AWS instance lauched via. AWS auto scaling group
type LifecycleTransition string

var (
	// LifecycleTransitionTerminating is the string representing an instance in a terminating wait state
	LifecycleTransitionTerminating = LifecycleTransition("autoscaling:EC2_INSTANCE_TERMINATING")
)

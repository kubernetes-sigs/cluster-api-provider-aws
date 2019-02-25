package machineautoscaler

import (
	"context"
	"errors"

	"github.com/golang/glog"
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1alpha1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	// MachineTargetFinalizer is the finalizer added to MachineAutoscaler
	// instances to allow for cleanup of annotations on target resources.
	MachineTargetFinalizer = "machinetarget.autoscaling.openshift.io"

	// MachineTargetOwnerAnnotation is the annotation used to mark a
	// target resource's autoscaling as owned by a MachineAutoscaler.
	MachineTargetOwnerAnnotation = "autoscaling.openshift.io/machineautoscaler"

	minSizeAnnotation = "machine.openshift.io/cluster-api-autoscaler-node-group-min-size"
	maxSizeAnnotation = "machine.openshift.io/cluster-api-autoscaler-node-group-max-size"
)

var (
	// ErrUnsupportedTarget is the error returned when a target references an
	// object with an unsupported GroupVersionKind.
	ErrUnsupportedTarget = errors.New("unsupported MachineAutoscaler target")

	// ErrInvalidTarget is the error returned when a target reference is invalid
	// in some way, e.g. not having a name set.
	ErrInvalidTarget = errors.New("invalid MachineAutoscaler target")

	// ErrNoSupportedTargets is the error returned during initialization if none
	// of the supported MachineAutoscaler targets are registered with the API.
	ErrNoSupportedTargets = errors.New("no supported target types available")
)

// SupportedTargetGVKs is the list of GroupVersionKinds supported as targets for
// a MachineAutocaler instance.
var SupportedTargetGVKs = []schema.GroupVersionKind{
	{Group: "cluster.k8s.io", Version: "v1alpha1", Kind: "MachineDeployment"},
	{Group: "cluster.k8s.io", Version: "v1alpha1", Kind: "MachineSet"},
	{Group: "machine.openshift.io", Version: "v1beta1", Kind: "MachineDeployment"},
	{Group: "machine.openshift.io", Version: "v1beta1", Kind: "MachineSet"},
}

// Config represents the configuration for a reconciler instance.
type Config struct {
	// The namespace for MachineAutosclaers and their targets.
	Namespace string
}

// NewReconciler returns a new Reconciler.
func NewReconciler(mgr manager.Manager, cfg *Config) *Reconciler {
	return &Reconciler{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		config: cfg,
	}
}

// AddToManager adds a new Controller to mgr with r as the reconcile.Reconciler
func (r *Reconciler) AddToManager(mgr manager.Manager) error {
	// Create a new controller
	c, err := controller.New("machineautoscaler-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MachineAutoscaler
	err = c.Watch(&source.Kind{Type: &v1alpha1.MachineAutoscaler{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to each supported target resource type and enqueue
	// reconcile requests for their owning MachineAutoscaler resources.
	for _, gvk := range SupportedTargetGVKs {
		target := &unstructured.Unstructured{}
		target.SetGroupVersionKind(gvk)

		err := c.Watch(
			&source.Kind{Type: target},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(targetOwnerRequest),
			})

		// If we get an error indicating that no matching type is registered
		// with the API, then remove it from the list of supported target GVKs.
		// If the type is later registered, a restart of the operator will pick
		// it up and properly reconcile any MachineAutoscalers referencing it.
		if err != nil && meta.IsNoMatchError(err) {
			glog.Warningf("Removing support for unregistered target type: %s", gvk)
			SupportedTargetGVKs = removeSupportedGVK(gvk)
		} else if err != nil {
			return err
		}
	}

	// Fail if we didn't find any of the supported target types registered.
	if len(SupportedTargetGVKs) < 1 {
		return ErrNoSupportedTargets
	}

	return nil
}

var _ reconcile.Reconciler = &Reconciler{}

// Reconciler reconciles a MachineAutoscaler object
type Reconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	config *Config
}

// Reconcile reads that state of the cluster for a MachineAutoscaler object and
// makes changes based on the state read and what is in the
// MachineAutoscaler.Spec
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	glog.Infof("Reconciling MachineAutoscaler %s/%s\n", request.Namespace, request.Name)

	// Fetch the MachineAutoscaler instance
	ma := &v1alpha1.MachineAutoscaler{}
	err := r.client.Get(context.TODO(), request.NamespacedName, ma)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile
			// request.  Owned objects are automatically garbage collected. For
			// additional cleanup logic use finalizers.
			// Return and don't requeue.
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		glog.Errorf("Error reading MachineAutoscaler: %v", err)
		return reconcile.Result{}, err
	}

	// Handle MachineAutoscaler deletion. This should happen directly after the
	// MachineAutoscaler has been fetched, before any further reconciliation.
	if ma.GetDeletionTimestamp() != nil {
		return r.HandleDelete(ma)
	}

	targetRef := objectReference(ma.Spec.ScaleTargetRef)

	target, err := r.GetTarget(targetRef)
	if err != nil {
		glog.Errorf("Error getting target: %v", err)
		return reconcile.Result{}, err
	}

	// Set the MachineAutoscaler as the owner of the target.
	ownerModifed, err := target.SetOwner(ma)
	if err != nil {
		glog.Errorf("Error setting target owner: %v", err)
		return reconcile.Result{}, err
	}

	// If the owner is newly added, remove any existing limits.
	// This will force an update to bring things into sync.
	if ownerModifed {
		target.RemoveLimits()
	}

	// If there is a previously observed target referenced in the
	// status, and it has changed relative to the current target, the
	// previous target must be finalized, e.g. annotations removed.
	if ma.Status.LastTargetRef != nil && r.TargetChanged(ma) {
		glog.V(2).Infof("%s: Target changed", request.NamespacedName)

		lastTargetRef := objectReference(*ma.Status.LastTargetRef)

		lastTarget, err := r.GetTarget(lastTargetRef)
		if err != nil && !apierrors.IsNotFound(err) {
			// If there was a problem (other than a 404) fetching the
			// previous target, we should retry.  Otherwise, it may
			// retain autoscaling configuration.
			glog.Errorf("Error fetching previous target: %v", err)
			return reconcile.Result{}, err
		}

		// If the target changed, and we were able to fetch the
		// previous target successfully, finalize it.
		if lastTarget != nil {
			err := r.FinalizeTarget(lastTarget)

			// Ignore 404s, the resource has most likely been deleted.
			if err != nil && !apierrors.IsNotFound(err) {
				glog.Errorf("Error finalizing previous target: %v", err)
				return reconcile.Result{}, err
			}
		}

		// Set the previous target equal to the current target.
		if err := r.SetLastTarget(ma, targetRef); err != nil {
			glog.Errorf("Error setting previous target: %v", err)
			return reconcile.Result{}, err
		}
	}

	// Set the previous target if we don't have one.
	if ma.Status.LastTargetRef == nil {
		if err := r.SetLastTarget(ma, targetRef); err != nil {
			glog.Errorf("Error setting previous target: %v", err)
			return reconcile.Result{}, err
		}
	}

	// Ensure our finalizers have been added.
	if err := r.EnsureFinalizer(ma); err != nil {
		glog.Errorf("Error setting finalizer: %v", err)
		return reconcile.Result{}, err
	}

	min := int(ma.Spec.MinReplicas)
	max := int(ma.Spec.MaxReplicas)

	if err := r.UpdateTarget(target, min, max); err != nil {
		glog.Errorf("Error updating target: %v", err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// HandleDelete is called by Reconcile to handle MachineAutoscaler deletion,
// i.e. finalize the resource and remove finalizers.
func (r *Reconciler) HandleDelete(ma *v1alpha1.MachineAutoscaler) (reconcile.Result, error) {
	targetRef := objectReference(ma.Spec.ScaleTargetRef)

	target, err := r.GetTarget(targetRef)
	if err != nil && !apierrors.IsNotFound(err) {
		glog.Errorf("Error getting target for finalization: %v", err)
		return reconcile.Result{}, err
	}

	if target != nil {
		err := r.FinalizeTarget(target)

		// Ignore 404s, the resource has most likely been deleted.
		if err != nil && !apierrors.IsNotFound(err) {
			glog.Errorf("Error finalizing target: %v", err)
			return reconcile.Result{}, err
		}
	}

	if err := r.RemoveFinalizer(ma); err != nil {
		glog.Errorf("Error removing finalizer: %v", err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// GetTarget fetches the object targeted by the given reference.
func (r *Reconciler) GetTarget(ref *corev1.ObjectReference) (*MachineTarget, error) {
	obj := &unstructured.Unstructured{}
	gvk := ref.GroupVersionKind()

	if valid, err := ValidateReference(ref); !valid {
		return nil, err
	}

	obj.SetGroupVersionKind(gvk)

	err := r.client.Get(context.TODO(), client.ObjectKey{
		Namespace: r.config.Namespace,
		Name:      ref.Name,
	}, obj)

	if err != nil {
		return nil, err
	}

	target, err := MachineTargetFromObject(obj)
	if err != nil {
		glog.Errorf("Failed to convert object to MachineTarget: %v", err)
		return nil, err
	}

	return target, nil
}

// UpdateTarget updates the min and max annotations on the given target.
func (r *Reconciler) UpdateTarget(target *MachineTarget, min, max int) error {
	// Update the target object's annotations if necessary.
	if target.NeedsUpdate(min, max) {
		target.SetLimits(min, max)

		return r.client.Update(context.TODO(), target)
	}

	return nil
}

// FinalizeTarget handles finalizers for the given target.
func (r *Reconciler) FinalizeTarget(target *MachineTarget) error {
	modified := target.Finalize()

	if modified {
		return r.client.Update(context.TODO(), target)
	}

	return nil
}

// TargetChanged indicates whether a MachineAutoscaler's current target has
// changed relative to the last observed target noted in the status.
func (r *Reconciler) TargetChanged(ma *v1alpha1.MachineAutoscaler) bool {
	currentRef := ma.Spec.ScaleTargetRef
	lastRef := ma.Status.LastTargetRef

	if lastRef != nil && !equality.Semantic.DeepEqual(currentRef, *lastRef) {
		return true
	}

	return false
}

// SetLastTarget updates the give MachineAutoscaler's status with the given
// object as the last observed target.
func (r *Reconciler) SetLastTarget(ma *v1alpha1.MachineAutoscaler, ref *corev1.ObjectReference) error {
	ma.Status.LastTargetRef = &v1alpha1.CrossVersionObjectReference{
		APIVersion: ref.APIVersion,
		Kind:       ref.Kind,
		Name:       ref.Name,
	}

	return r.client.Status().Update(context.TODO(), ma)
}

// EnsureFinalizer adds finalizers to the given MachineAutoscaler if necessary.
func (r *Reconciler) EnsureFinalizer(ma *v1alpha1.MachineAutoscaler) error {
	for _, f := range ma.GetFinalizers() {
		// Bail early if we already have the finalizer.
		if f == MachineTargetFinalizer {
			return nil
		}
	}

	f := append(ma.GetFinalizers(), MachineTargetFinalizer)
	ma.SetFinalizers(f)

	return r.client.Update(context.TODO(), ma)
}

// RemoveFinalizer removes this packages's finalizers from the given
// MachineAutoscaler instance.
func (r *Reconciler) RemoveFinalizer(ma *v1alpha1.MachineAutoscaler) error {
	f, found := util.FilterString(ma.GetFinalizers(), MachineTargetFinalizer)

	if found == 0 {
		return nil
	}

	ma.SetFinalizers(f)

	return r.client.Update(context.TODO(), ma)
}

// SupportedTarget indicates whether a GVK is supported as a target.
func SupportedTarget(gvk schema.GroupVersionKind) bool {
	for _, supported := range SupportedTargetGVKs {
		if gvk == supported {
			return true
		}
	}

	return false
}

// ValidateReference validates that an object reference is valid, i.e. that it
// has a name and a supported GroupVersionKind.  If this method returns false,
// indicating that the reference is not valid, it MUST return a non-nil error.
func ValidateReference(obj *corev1.ObjectReference) (bool, error) {
	if obj == nil {
		return false, ErrInvalidTarget
	}

	if obj.Name == "" {
		return false, ErrInvalidTarget
	}

	if !SupportedTarget(obj.GroupVersionKind()) {
		return false, ErrUnsupportedTarget
	}

	return true, nil
}

// targetOwnerRequest is used with handler.EnqueueRequestsFromMapFunc to enqueue
// reconcile requests for the owning MachineAutoscaler of a watched target.
func targetOwnerRequest(a handler.MapObject) []reconcile.Request {
	target, err := MachineTargetFromObject(a.Object)
	if err != nil {
		glog.Errorf("Failed to convert object to MachineTarget: %v", err)
		return nil
	}

	owner, err := target.GetOwner()
	if err != nil {
		glog.V(2).Infof("Will not reconcile: %v", err)
		return nil
	}

	glog.V(2).Infof("Queuing reconcile for owner of %s/%s.",
		target.GetNamespace(), target.GetName())

	return []reconcile.Request{{NamespacedName: owner}}
}

// objectReference returns a new corev1.ObjectReference for the given
// CrossVersionObjectReference from a MachineAutoscaler target.
func objectReference(ref v1alpha1.CrossVersionObjectReference) *corev1.ObjectReference {
	obj := &corev1.ObjectReference{}
	gvk := schema.FromAPIVersionAndKind(ref.APIVersion, ref.Kind)

	obj.SetGroupVersionKind(gvk)
	obj.Name = ref.Name

	return obj
}

// removeSupportedGVK removes the given type from the list of supported GVKs for
// MachineAutoscaler targets.
func removeSupportedGVK(gvk schema.GroupVersionKind) []schema.GroupVersionKind {
	newSlice := SupportedTargetGVKs[:0] // Share the backing array.

	for _, x := range SupportedTargetGVKs {
		if x != gvk {
			newSlice = append(newSlice, x)
		}
	}

	return newSlice
}

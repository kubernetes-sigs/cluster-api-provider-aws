/*
Copyright 2024 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// AWSMachineTemplateReconciler reconciles AWSMachineTemplate objects.
//
// This controller automatically populates capacity information for AWSMachineTemplate resources
// to enable autoscaling from zero.
//
// See: https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
type AWSMachineTemplateReconciler struct {
	client.Client
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinetemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinetemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

// Reconcile populates capacity information for AWSMachineTemplate.
func (r *AWSMachineTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// Fetch the AWSMachineTemplate
	awsMachineTemplate := &infrav1.AWSMachineTemplate{}
	if err := r.Get(ctx, req.NamespacedName, awsMachineTemplate); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Skip if capacity is already set
	if len(awsMachineTemplate.Status.Capacity) > 0 {
		return ctrl.Result{}, nil
	}

	// Get instance type from spec
	instanceType := awsMachineTemplate.Spec.Template.Spec.InstanceType
	if instanceType == "" {
		return ctrl.Result{}, nil
	}

	// Find the region by checking ownerReferences
	region, err := r.getRegion(ctx, awsMachineTemplate)
	if err != nil {
		return ctrl.Result{}, err
	}
	if region == "" {
		return ctrl.Result{}, nil
	}

	// Create global scope for this region
	// Reference: exp/instancestate/awsinstancestate_controller.go:68-76
	globalScope, err := scope.NewGlobalScope(scope.GlobalScopeParams{
		ControllerName: "awsmachinetemplate",
		Region:         region,
	})
	if err != nil {
		record.Warnf(awsMachineTemplate, "AWSSessionFailed", "Failed to create AWS session for region %q: %v", region, err)
		return ctrl.Result{}, nil
	}

	// Query instance type capacity
	capacity, err := r.getInstanceTypeCapacity(ctx, globalScope, instanceType)
	if err != nil {
		record.Warnf(awsMachineTemplate, "CapacityQueryFailed", "Failed to query capacity for instance type %q: %v", instanceType, err)
		return ctrl.Result{}, nil
	}

	// Update status with capacity
	awsMachineTemplate.Status.Capacity = capacity

	if err := r.Status().Update(ctx, awsMachineTemplate); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to update AWSMachineTemplate status")
	}

	log.Info("Successfully populated capacity information", "instanceType", instanceType, "region", region, "capacity", capacity)
	return ctrl.Result{}, nil
}

// getRegion finds the region by checking the template's owner cluster reference.
func (r *AWSMachineTemplateReconciler) getRegion(ctx context.Context, template *infrav1.AWSMachineTemplate) (string, error) {
	// Check if the template has a Cluster owner reference (e.g., in Hypershift scenarios)
	for _, owner := range template.GetOwnerReferences() {
		if owner.Kind == "Cluster" && owner.APIVersion == clusterv1.GroupVersion.String() {
			cluster := &clusterv1.Cluster{}
			clusterKey := client.ObjectKey{
				Namespace: template.Namespace,
				Name:      owner.Name,
			}
			if err := r.Get(ctx, clusterKey, cluster); err != nil {
				if !apierrors.IsNotFound(err) {
					return "", errors.Wrapf(err, "failed to get Cluster %s/%s", template.Namespace, owner.Name)
				}
				continue
			}

			// Try to get region from AWSCluster (standard cluster)
			if cluster.Spec.InfrastructureRef != nil && cluster.Spec.InfrastructureRef.Kind == "AWSCluster" {
				awsCluster := &infrav1.AWSCluster{}
				if err := r.Get(ctx, client.ObjectKey{
					Namespace: cluster.Namespace,
					Name:      cluster.Spec.InfrastructureRef.Name,
				}, awsCluster); err != nil {
					if !apierrors.IsNotFound(err) {
						return "", errors.Wrapf(err, "failed to get AWSCluster %s/%s", cluster.Namespace, cluster.Spec.InfrastructureRef.Name)
					}
				} else if awsCluster.Spec.Region != "" {
					return awsCluster.Spec.Region, nil
				}
			}
		}
	}

	return "", nil
}

// getInstanceTypeCapacity queries AWS EC2 API for instance type capacity.
func (r *AWSMachineTemplateReconciler) getInstanceTypeCapacity(ctx context.Context, globalScope *scope.GlobalScope, instanceType string) (corev1.ResourceList, error) {
	// Create EC2 client from global scope
	ec2Client := ec2.NewFromConfig(globalScope.Session())

	// Query instance type information
	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []ec2types.InstanceType{ec2types.InstanceType(instanceType)},
	}

	result, err := ec2Client.DescribeInstanceTypes(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe instance type %q", instanceType)
	}

	if len(result.InstanceTypes) == 0 {
		return nil, errors.Errorf("no information found for instance type %q", instanceType)
	}

	// Extract capacity information
	info := result.InstanceTypes[0]
	resourceList := corev1.ResourceList{}

	// CPU
	if info.VCpuInfo != nil && info.VCpuInfo.DefaultVCpus != nil {
		resourceList[corev1.ResourceCPU] = *resource.NewQuantity(int64(*info.VCpuInfo.DefaultVCpus), resource.DecimalSI)
	}

	// Memory
	if info.MemoryInfo != nil && info.MemoryInfo.SizeInMiB != nil {
		memoryBytes := *info.MemoryInfo.SizeInMiB * 1024 * 1024
		resourceList[corev1.ResourceMemory] = *resource.NewQuantity(memoryBytes, resource.BinarySI)
	}
	return resourceList, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSMachineTemplateReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.AWSMachineTemplate{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Complete(r)
}

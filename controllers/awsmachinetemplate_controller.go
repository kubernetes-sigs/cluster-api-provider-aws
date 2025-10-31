/*
Copyright 2025 The Kubernetes Authors.

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
	"strings"

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
	ec2service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
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
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=kubeadmcontrolplanes,verbs=get;list;watch
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

	// Skip if capacity and nodeInfo are already set
	if len(awsMachineTemplate.Status.Capacity) > 0 && awsMachineTemplate.Status.NodeInfo != nil {
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

	// Create EC2 client from global scope
	ec2Client := ec2.NewFromConfig(globalScope.Session())

	// Query instance type capacity
	capacity, err := r.getInstanceTypeCapacity(ctx, ec2Client, instanceType)
	if err != nil {
		record.Warnf(awsMachineTemplate, "CapacityQueryFailed", "Failed to query capacity for instance type %q: %v", instanceType, err)
		return ctrl.Result{}, nil
	}

	// Query node info (architecture and OS)
	nodeInfo, err := r.getNodeInfo(ctx, ec2Client, awsMachineTemplate, instanceType)
	if err != nil {
		record.Warnf(awsMachineTemplate, "NodeInfoQueryFailed", "Failed to query node info for instance type %q: %v", instanceType, err)
		return ctrl.Result{}, nil
	}

	// Save original before modifying, then update all status fields at once
	original := awsMachineTemplate.DeepCopy()
	if len(capacity) > 0 {
		awsMachineTemplate.Status.Capacity = capacity
	}
	if nodeInfo != nil && (nodeInfo.Architecture != "" || nodeInfo.OperatingSystem != "") {
		awsMachineTemplate.Status.NodeInfo = nodeInfo
	}
	if err := r.Status().Patch(ctx, awsMachineTemplate, client.MergeFrom(original)); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to update AWSMachineTemplate status")
	}

	log.Info("Successfully populated capacity and nodeInfo", "instanceType", instanceType, "region", region, "capacity", capacity, "nodeInfo", nodeInfo)
	return ctrl.Result{}, nil
}

// getRegion finds the region by checking the template's owner cluster reference.
func (r *AWSMachineTemplateReconciler) getRegion(ctx context.Context, template *infrav1.AWSMachineTemplate) (string, error) {
	// Get the owner cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, template.ObjectMeta)
	if err != nil {
		return "", err
	}
	if cluster == nil {
		return "", errors.New("no owner cluster found")
	}

	// Get region from AWSCluster (standard EC2-based cluster)
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

	return "", nil
}

// getInstanceTypeCapacity queries AWS EC2 API for instance type capacity information.
// Returns the resource list (CPU, Memory).
func (r *AWSMachineTemplateReconciler) getInstanceTypeCapacity(ctx context.Context, ec2Client *ec2.Client, instanceType string) (corev1.ResourceList, error) {
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

// getNodeInfo queries node information (architecture and OS) for the AWSMachineTemplate.
// It uses AMI ID if specified, otherwise attempts AMI lookup or falls back to instance type info.
func (r *AWSMachineTemplateReconciler) getNodeInfo(ctx context.Context, ec2Client *ec2.Client, template *infrav1.AWSMachineTemplate, instanceType string) (*infrav1.NodeInfo, error) {
	nodeInfo := &infrav1.NodeInfo{}
	amiID := template.Spec.Template.Spec.AMI.ID
	if amiID != nil && *amiID != "" {
		// AMI ID is specified, query it directly
		arch, os, err := r.getNodeInfoFromAMI(ctx, ec2Client, *amiID)
		if err == nil {
			if arch != "" {
				nodeInfo.Architecture = arch
			}
			if os != "" {
				nodeInfo.OperatingSystem = os
			}
		}
	} else {
		// AMI ID is not specified, query instance type to get architecture
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

		instanceTypeInfo := result.InstanceTypes[0]

		// Infer architecture from instance type
		var architecture string
		if instanceTypeInfo.ProcessorInfo != nil && len(instanceTypeInfo.ProcessorInfo.SupportedArchitectures) == 1 {
			// Use the supported architecture
			switch instanceTypeInfo.ProcessorInfo.SupportedArchitectures[0] {
			case ec2types.ArchitectureTypeX8664:
				architecture = ec2service.Amd64ArchitectureTag
				nodeInfo.Architecture = infrav1.ArchitectureAmd64
			case ec2types.ArchitectureTypeArm64:
				architecture = ec2service.Arm64ArchitectureTag
				nodeInfo.Architecture = infrav1.ArchitectureArm64
			}
		} else {
			return nil, errors.Errorf("instance type must support exactly one architecture, got %d", len(instanceTypeInfo.ProcessorInfo.SupportedArchitectures))
		}

		// Attempt to get Kubernetes version from MachineDeployment
		kubernetesVersion, versionErr := r.getKubernetesVersion(ctx, template)
		if versionErr == nil && kubernetesVersion != "" {
			// Try to look up AMI using the version
			image, err := ec2service.DefaultAMILookup(
				ec2Client,
				template.Spec.Template.Spec.ImageLookupOrg,
				template.Spec.Template.Spec.ImageLookupBaseOS,
				kubernetesVersion,
				architecture,
				template.Spec.Template.Spec.ImageLookupFormat,
			)
			if err == nil && image != nil {
				// Successfully found AMI, extract accurate nodeInfo from it
				arch, os, _ := r.getNodeInfoFromAMI(ctx, ec2Client, *image.ImageId)
				if arch != "" {
					nodeInfo.Architecture = arch
				}
				if os != "" {
					nodeInfo.OperatingSystem = os
				}
				return nodeInfo, nil
			}
			// AMI lookup failed, fall through to defaults
		}
	}

	return nodeInfo, nil
}

// getNodeInfoFromAMI queries the AMI to determine architecture and operating system.
func (r *AWSMachineTemplateReconciler) getNodeInfoFromAMI(ctx context.Context, ec2Client *ec2.Client, amiID string) (infrav1.Architecture, string, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []string{amiID},
	}

	result, err := ec2Client.DescribeImages(ctx, input)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to describe AMI %q", amiID)
	}

	if len(result.Images) == 0 {
		return "", "", errors.Errorf("no information found for AMI %q", amiID)
	}

	image := result.Images[0]

	// Get architecture from AMI
	var arch infrav1.Architecture
	switch image.Architecture {
	case ec2types.ArchitectureValuesX8664:
		arch = infrav1.ArchitectureAmd64
	case ec2types.ArchitectureValuesArm64:
		arch = infrav1.ArchitectureArm64
	}

	// Determine OS - default to Linux, change to Windows if detected
	// Most AMIs are Linux-based, so we initialize with Linux as the default
	os := infrav1.OperatingSystemLinux

	// 1. Check Platform field (most reliable for Windows detection)
	if image.Platform == ec2types.PlatformValuesWindows {
		os = infrav1.OperatingSystemWindows
	}

	// 2. Check PlatformDetails field for Windows indication
	if os != infrav1.OperatingSystemWindows && image.PlatformDetails != nil {
		platformDetails := strings.ToLower(*image.PlatformDetails)
		if strings.Contains(platformDetails, infrav1.OperatingSystemWindows) {
			os = infrav1.OperatingSystemWindows
		}
	}

	return arch, os, nil
}

// getKubernetesVersion attempts to find the Kubernetes version by querying MachineDeployments
// or KubeadmControlPlanes that reference this AWSMachineTemplate.
func (r *AWSMachineTemplateReconciler) getKubernetesVersion(ctx context.Context, template *infrav1.AWSMachineTemplate) (string, error) {
	// Try to find version from MachineDeployment first
	machineDeploymentList := &clusterv1.MachineDeploymentList{}
	if err := r.List(ctx, machineDeploymentList, client.InNamespace(template.Namespace)); err != nil {
		return "", errors.Wrap(err, "failed to list MachineDeployments")
	}

	// Find MachineDeployments that reference this AWSMachineTemplate
	for _, md := range machineDeploymentList.Items {
		if md.Spec.Template.Spec.InfrastructureRef.Kind == "AWSMachineTemplate" &&
			md.Spec.Template.Spec.InfrastructureRef.Name == template.Name &&
			md.Spec.Template.Spec.Version != nil {
			return *md.Spec.Template.Spec.Version, nil
		}
	}

	// If not found in MachineDeployment, try KubeadmControlPlane
	kcpList := &controlplanev1.KubeadmControlPlaneList{}
	if err := r.List(ctx, kcpList, client.InNamespace(template.Namespace)); err != nil {
		return "", errors.Wrap(err, "failed to list KubeadmControlPlanes")
	}

	// Find KubeadmControlPlanes that reference this AWSMachineTemplate
	for _, kcp := range kcpList.Items {
		if kcp.Spec.MachineTemplate.InfrastructureRef.Kind == "AWSMachineTemplate" &&
			kcp.Spec.MachineTemplate.InfrastructureRef.Name == template.Name &&
			kcp.Spec.Version != "" {
			return kcp.Spec.Version, nil
		}
	}

	return "", errors.New("no MachineDeployment or KubeadmControlPlane found referencing this AWSMachineTemplate with a version")
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

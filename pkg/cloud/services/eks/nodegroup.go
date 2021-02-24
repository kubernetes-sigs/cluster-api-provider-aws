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

package eks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
)

func (s *NodegroupService) describeNodegroup() (*awseks.Nodegroup, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	nodegroupName := s.scope.NodegroupName()
	s.scope.V(2).Info("describing eks node group", "cluster", eksClusterName, "nodegroup", nodegroupName)
	input := &awseks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
	}

	out, err := s.EKSClient.DescribeNodegroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case awseks.ErrCodeResourceNotFoundException:
				return nil, nil
			default:
				return nil, errors.Wrap(err, "failed to describe nodegroup")
			}
		} else {
			return nil, errors.Wrap(err, "failed to describe nodegroup")
		}
	}

	return out.Nodegroup, nil
}

func (s *NodegroupService) scalingConfig() *awseks.NodegroupScalingConfig {
	var replicas int32 = 1
	if s.scope.MachinePool.Spec.Replicas != nil {
		replicas = *s.scope.MachinePool.Spec.Replicas
	}
	cfg := awseks.NodegroupScalingConfig{
		DesiredSize: aws.Int64(int64(replicas)),
	}
	scaling := s.scope.ManagedMachinePool.Spec.Scaling
	if scaling == nil {
		return &cfg
	}
	if scaling.MaxSize != nil {
		cfg.MaxSize = aws.Int64(int64(*scaling.MaxSize))
	}
	if scaling.MaxSize != nil {
		cfg.MinSize = aws.Int64(int64(*scaling.MinSize))
	}
	return &cfg
}

func (s *NodegroupService) subnets() []string {
	subnetIDs := s.scope.SubnetIDs()
	// If not specified, use all
	if len(subnetIDs) == 0 {
		subnetIDs := []string{}
		for _, subnet := range s.scope.ControlPlaneSubnets() {
			subnetIDs = append(subnetIDs, subnet.ID)
		}
		return subnetIDs
	}
	return subnetIDs
}

func (s *NodegroupService) roleArn() (*string, error) {
	var role *iam.Role
	if s.scope.RoleName() != "" {
		var err error
		role, err = s.GetIAMRole(s.scope.RoleName())
		if err != nil {
			return nil, errors.Wrapf(err, "error getting node group IAM role: %s", s.scope.RoleName())
		}
	}
	return role.Arn, nil
}

func ngTags(key string, additionalTags infrav1.Tags) map[string]string {
	tags := additionalTags.DeepCopy()
	tags[infrav1.ClusterAWSCloudProviderTagKey(key)] = string(infrav1.ResourceLifecycleOwned)
	return tags
}

func (s *NodegroupService) remoteAccess(data *ec2.RequestLaunchTemplateData) error {
	pool := s.scope.ManagedMachinePool.Spec
	if pool.RemoteAccess == nil {
		return nil
	}
	controlPlane := s.scope.ControlPlane
	sshKeyName := pool.RemoteAccess.SSHKeyName
	if sshKeyName == nil {
		sshKeyName = controlPlane.Spec.SSHKeyName
	}
	data.KeyName = sshKeyName
	if !pool.RemoteAccess.Public && pool.RemoteAccess.SourceSecurityGroups != nil {
		data.SecurityGroupIds = append(data.SecurityGroupIds, aws.StringSlice(pool.RemoteAccess.SourceSecurityGroups)...)
	}
	return nil
}

func (m *NodegroupService) CreateLaunchTemplateData(ec2Svc *ec2svc.Service, scope scope.EC2Scope, userData []byte) (*ec2.RequestLaunchTemplateData, error) {
	spec := m.scope.ManagedMachinePool.Spec

	var imageID string
	// if set imageID is set we need to make the bootstrap  TODO(felipeweb)
	//
	// if amiID := spec.AMIID; amiID != nil {
	// 	imageID = *amiID
	// } else {
	// 	paramName, err := ami.EKSAMIParameter(*m.scope.MachinePool.Spec.Template.Spec.Version, *m.scope.ManagedMachinePool.Spec.AMIType)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	lookupAMI, err := ec2Svc.SSMAMILookup(paramName)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	imageID = lookupAMI
	// }

	data := &ec2.RequestLaunchTemplateData{
		// ImageId:      aws.String(imageID),
		InstanceType: m.scope.ManagedMachinePool.Spec.InstanceType,
		// Userdata must have Content-Type: multipart/mixed; boundary="==BOUNDARY==" TODO
		// UserData:     pointer.StringPtr(base64.StdEncoding.EncodeToString(userData)),
	}

	sgIDs, err := m.scope.CoreSecurityGroups(scope)
	if err != nil {
		return nil, err
	}
	data.SecurityGroupIds = aws.StringSlice(sgIDs)
	if spec.DiskSize != nil {
		size := int64(*spec.DiskSize)
		deviceName, devices, err := ec2Svc.CheckRootVolume(&infrav1.Volume{Size: size}, imageID)
		if err != nil {
			return nil, err
		}

		for _, dev := range devices {
			if aws.StringValue(dev.DeviceName) != aws.StringValue(deviceName) {
				continue
			}
			data.BlockDeviceMappings = []*ec2.LaunchTemplateBlockDeviceMappingRequest{{
				DeviceName: deviceName,
				Ebs: &ec2.LaunchTemplateEbsBlockDeviceRequest{
					VolumeType:          dev.Ebs.VolumeType,
					DeleteOnTermination: dev.Ebs.DeleteOnTermination,
					Encrypted:           dev.Ebs.Encrypted,
					Iops:                dev.Ebs.Iops,
					KmsKeyId:            dev.Ebs.KmsKeyId,
					SnapshotId:          dev.Ebs.SnapshotId,
					Throughput:          dev.Ebs.Throughput,
					VolumeSize:          aws.Int64(int64(*spec.DiskSize)),
				},
				NoDevice:    dev.NoDevice,
				VirtualName: dev.VirtualName,
			}}
		}
	}

	if err := m.remoteAccess(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *NodegroupService) LaunchTemplateNeedsUpdate(incoming *ec2.RequestLaunchTemplateData, existing *ec2.LaunchTemplateVersion) (bool, error) {
	if aws.StringValue(incoming.InstanceType) != aws.StringValue(existing.LaunchTemplateData.InstanceType) || aws.StringValue(incoming.ImageId) != aws.StringValue(existing.LaunchTemplateData.ImageId) {
		return true, nil
	}

	return false, nil
}

func (s *NodegroupService) createNodegroup() (*awseks.Nodegroup, error) {
	managedPool := s.scope.ManagedMachinePool.Spec

	if managedPool.LaunchTemplate == nil || managedPool.LaunchTemplate.ID == "" {
		return nil, errors.New("unable to create nodegroup without launch template")
	}

	eksClusterName := s.scope.KubernetesClusterName()
	nodegroupName := s.scope.NodegroupName()
	additionalTags := s.scope.AdditionalTags()
	roleArn, err := s.roleArn()
	if err != nil {
		return nil, err
	}
	tags := ngTags(s.scope.ClusterName(), additionalTags)

	input := &awseks.CreateNodegroupInput{
		ScalingConfig: s.scalingConfig(),
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
		Subnets:       aws.StringSlice(s.subnets()),
		NodeRole:      roleArn,
		Labels:        aws.StringMap(managedPool.Labels),
		Tags:          aws.StringMap(tags),
		LaunchTemplate: &awseks.LaunchTemplateSpecification{
			Id:      aws.String(managedPool.LaunchTemplate.ID),
			Version: aws.String(managedPool.LaunchTemplate.Version),
		},
	}
	if managedPool.DiskSize != nil {
		// TODO !
		input.DiskSize = aws.Int64(int64(*managedPool.DiskSize))
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "created invalid CreateNodegroupInput")
	}

	out, err := s.EKSClient.CreateNodegroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// TODO
			case awseks.ErrCodeResourceNotFoundException:
				return nil, nil
			default:
				return nil, errors.Wrap(err, "failed to create nodegroup")
			}
		} else {
			return nil, errors.Wrap(err, "failed to create nodegroup")
		}
	}

	return out.Nodegroup, nil
}

func (s *NodegroupService) deleteNodegroupAndWait() (reterr error) {
	eksClusterName := s.scope.KubernetesClusterName()
	nodegroupName := s.scope.NodegroupName()
	if err := s.scope.NodegroupReadyFalse(clusterv1.DeletingReason, ""); err != nil {
		return err
	}
	defer func() {
		if reterr != nil {
			record.Warnf(
				s.scope.ManagedMachinePool, "FailedDeleteEKSNodegroup", "Failed to delete EKS nodegroup %s: %v", s.scope.NodegroupName(), reterr,
			)
			if err := s.scope.NodegroupReadyFalse("DeletingFailed", reterr.Error()); err != nil {
				reterr = err
			}
		} else if err := s.scope.NodegroupReadyFalse(clusterv1.DeletedReason, ""); err != nil {
			reterr = err
		}
	}()
	input := &awseks.DeleteNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
	}
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "created invalid DeleteNodegroupInput")
	}

	_, err := s.EKSClient.DeleteNodegroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// TODO
			case awseks.ErrCodeResourceNotFoundException:
				return nil
			default:
				return errors.Wrap(err, "failed to delete nodegroup")
			}
		} else {
			return errors.Wrap(err, "failed to delete nodegroup")
		}
	}

	waitInput := &awseks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
	}
	err = s.EKSClient.WaitUntilNodegroupDeleted(waitInput)
	if err != nil {
		return errors.Wrapf(err, "failed waiting for EKS nodegroup %s to delete", nodegroupName)
	}

	return nil
}

func (s *NodegroupService) reconcileNodegroupVersion(ng *awseks.Nodegroup) error {
	eksClusterName := s.scope.KubernetesClusterName()
	spec := s.scope.ManagedMachinePool.Spec
	if spec.LaunchTemplate == nil {
		return errors.New("cant' reconcile version without a launch template reference")
	}
	ltVersion := spec.LaunchTemplate.Version
	if aws.StringValue(ng.LaunchTemplate.Id) != spec.LaunchTemplate.ID ||
		aws.StringValue(ng.LaunchTemplate.Version) != ltVersion {

		input := &awseks.UpdateNodegroupVersionInput{
			ClusterName:   aws.String(eksClusterName),
			NodegroupName: aws.String(s.scope.NodegroupName()),
			LaunchTemplate: &awseks.LaunchTemplateSpecification{
				Id:      aws.String(spec.LaunchTemplate.ID),
				Version: aws.String(ltVersion),
			},
		}

		updateMsg := fmt.Sprintf("to launch template id/version: %s/%s", *input.LaunchTemplate.Id, *input.LaunchTemplate.Version)

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateNodegroupVersion(input); err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					return false, aerr
				}
				return false, err
			}
			record.Eventf(s.scope.ManagedMachinePool, "SuccessfulUpdateEKSNodegroup", "Updated EKS nodegroup %s %s", eksClusterName, updateMsg)
			return true, nil
		}); err != nil {
			record.Warnf(s.scope.ManagedMachinePool, "FailedUpdateEKSNodegroup", "failed to update the EKS nodegroup %s %s: %v", eksClusterName, updateMsg, err)
			return errors.Wrapf(err, "failed to update EKS nodegroup")
		}
	}
	return nil
}

func createLabelUpdate(specLabels map[string]string, ng *awseks.Nodegroup) *awseks.UpdateLabelsPayload {
	current := ng.Labels
	payload := awseks.UpdateLabelsPayload{}
	for k, v := range specLabels {
		if currentV, ok := current[k]; !ok || currentV == nil || v != *currentV {
			payload.AddOrUpdateLabels[k] = aws.String(v)
		}
	}
	for k := range current {
		if _, ok := specLabels[k]; !ok {
			payload.RemoveLabels = append(payload.RemoveLabels, aws.String(k))
		}
	}
	if len(payload.AddOrUpdateLabels) > 0 || len(payload.RemoveLabels) > 0 {
		return &payload
	}
	return nil
}

func (s *NodegroupService) reconcileNodegroupConfig(ng *awseks.Nodegroup) error {
	eksClusterName := s.scope.KubernetesClusterName()
	machinePool := s.scope.MachinePool.Spec
	managedPool := s.scope.ManagedMachinePool.Spec
	input := &awseks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(managedPool.EKSNodegroupName),
	}
	var needsUpdate bool
	if labelPayload := createLabelUpdate(managedPool.Labels, ng); labelPayload != nil {
		input.Labels = labelPayload
		needsUpdate = true
	}
	if machinePool.Replicas == nil {
		if ng.ScalingConfig.DesiredSize != nil && *ng.ScalingConfig.DesiredSize != 1 {
			input.ScalingConfig = s.scalingConfig()
			needsUpdate = true
		}
	} else if ng.ScalingConfig.DesiredSize == nil || int64(*machinePool.Replicas) != *ng.ScalingConfig.DesiredSize {
		input.ScalingConfig = s.scalingConfig()
		needsUpdate = true
	}
	if !needsUpdate {
		return nil
	}
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "created invalid UpdateNodegroupConfigInput")
	}

	_, err := s.EKSClient.UpdateNodegroupConfig(input)
	if err != nil {
		return errors.Wrap(err, "failed to update nodegroup config")
	}

	return nil
}

func (s *NodegroupService) reconcileNodegroup() error {
	eksClusterName := s.scope.KubernetesClusterName()
	eksNodegroupName := s.scope.NodegroupName()

	ng, err := s.describeNodegroup()
	if err != nil {
		return errors.Wrap(err, "failed to describe nodegroup")
	}

	if ng == nil {
		ng, err = s.createNodegroup()
		if err != nil {
			return errors.Wrap(err, "failed to create nodegroup")
		}
		s.scope.Info("Created EKS nodegroup in AWS", "cluster-name", eksClusterName, "nodegroup-name", eksNodegroupName)
	} else {
		tagKey := infrav1.ClusterAWSCloudProviderTagKey(s.scope.ClusterName())
		ownedTag := ng.Tags[tagKey]
		if ownedTag == nil {
			return errors.Wrapf(err, "owner of %s mismatch: %s", eksNodegroupName, s.scope.ClusterName())
		}
		s.scope.V(2).Info("Found owned EKS nodegroup in AWS", "cluster-name", eksClusterName, "nodegroup-name", eksNodegroupName)
	}

	if err := s.setStatus(ng); err != nil {
		return errors.Wrap(err, "failed to set status")
	}

	switch *ng.Status {
	case awseks.NodegroupStatusCreating, awseks.NodegroupStatusUpdating:
		ng, err = s.waitForNodegroupActive()
	default:
		break
	}

	if err != nil {
		return errors.Wrap(err, "failed to wait for nodegroup to be active")
	}

	if err := s.reconcileNodegroupVersion(ng); err != nil {
		return errors.Wrap(err, "failed to reconcile nodegroup version")
	}

	if err := s.reconcileNodegroupConfig(ng); err != nil {
		return errors.Wrap(err, "failed to reconcile nodegroup config")
	}

	if err := s.reconcileTags(ng); err != nil {
		return errors.Wrapf(err, "failed to reconcile nodegroup tags")
	}

	return nil
}

func (s *NodegroupService) setStatus(ng *awseks.Nodegroup) error {
	managedPool := s.scope.ManagedMachinePool
	switch *ng.Status {
	case awseks.NodegroupStatusDeleting:
		managedPool.Status.Ready = false
	case awseks.NodegroupStatusCreateFailed, awseks.NodegroupStatusDeleteFailed:
		managedPool.Status.Ready = false
		// TODO FailureReason
		failureMsg := fmt.Sprintf("EKS nodegroup in failed %s status", *ng.Status)
		managedPool.Status.FailureMessage = &failureMsg
	case awseks.NodegroupStatusActive:
		managedPool.Status.Ready = true
		managedPool.Status.FailureMessage = nil
		// TODO FailureReason
	case awseks.NodegroupStatusCreating:
		managedPool.Status.Ready = false
	case awseks.NodegroupStatusUpdating:
		managedPool.Status.Ready = true
	default:
		return errors.Errorf("unexpected EKS nodegroup status %s", *ng.Status)
	}
	if managedPool.Status.Ready && ng.Resources != nil && len(ng.Resources.AutoScalingGroups) > 0 {
		req := autoscaling.DescribeAutoScalingGroupsInput{}
		for _, asg := range ng.Resources.AutoScalingGroups {
			req.AutoScalingGroupNames = append(req.AutoScalingGroupNames, asg.Name)
		}
		groups, err := s.AutoscalingClient.DescribeAutoScalingGroups(&req)
		if err != nil {
			return errors.Wrap(err, "failed to describe AutoScalingGroup for nodegroup")
		}

		var replicas int32
		var providerIDList []string
		for _, group := range groups.AutoScalingGroups {
			replicas += int32(len(group.Instances))
			for _, instance := range group.Instances {
				id, err := noderefutil.NewProviderID(fmt.Sprintf("aws://%s/%s", *instance.AvailabilityZone, *instance.InstanceId))
				if err != nil {
					s.Error(err, "couldn't create provider ID for instance", "id", *instance.InstanceId)
					continue
				}
				providerIDList = append(providerIDList, id.String())
			}
		}
		managedPool.Spec.ProviderIDList = providerIDList
		managedPool.Status.Replicas = replicas
	}
	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to update nodegroup")
	}
	return nil
}

func (s *NodegroupService) waitForNodegroupActive() (*awseks.Nodegroup, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	eksNodegroupName := s.scope.NodegroupName()
	req := awseks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(eksNodegroupName),
	}
	if err := s.EKSClient.WaitUntilNodegroupActive(&req); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for EKS nodegroup %q", *req.NodegroupName)
	}

	s.scope.Info("EKS nodegroup is now available", "nodegroup-name", eksNodegroupName)

	ng, err := s.describeNodegroup()
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe EKS nodegroup")
	}
	if err := s.setStatus(ng); err != nil {
		return nil, errors.Wrap(err, "failed to set status")
	}

	return ng, nil
}

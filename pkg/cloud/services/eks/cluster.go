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
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/version"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

func (s *Service) reconcileCluster(ctx context.Context) error {
	s.scope.V(2).Info("Reconciling EKS cluster")

	cluster, err := s.describeEKSCluster()
	if err != nil {
		return errors.Wrap(err, "failed to describe eks clusters")
	}

	if cluster == nil {
		cluster, err = s.createCluster()
		if err != nil {
			return errors.Wrap(err, "failed to create cluster")
		}
		s.scope.Info("Created EKS control plane: %s", *cluster.Name)
	} else {
		s.scope.V(2).Info("Found EKS control plane: %s", *cluster.Name)
	}

	if err := s.setStatus(cluster); err != nil {
		return errors.Wrap(err, "failed to set status")
	}

	// Wait for our cluster to be ready if necessary
	switch *cluster.Status {
	case eks.ClusterStatusUpdating:
		cluster, err = s.waitForClusterUpdate()
	case eks.ClusterStatusCreating:
		cluster, err = s.waitForClusterActive()
	default:
		break
	}
	if err != nil {
		return errors.Wrap(err, "failed to wait for cluster to be active")
	}

	if !s.scope.ControlPlane.Status.Ready {
		return nil
	}

	s.scope.V(2).Info("EKS Control Plane active endpoint: %s", *cluster.Endpoint)

	s.scope.ControlPlane.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
		Host: *cluster.Endpoint,
		Port: 443,
	}

	if err := s.reconcileKubeconfig(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling kubeconfig")
	}

	if err := s.reconcileClusterVersion(cluster); err != nil {
		return errors.Wrap(err, "failed reconciling cluster version")
	}

	if err := s.reconcileClusterConfig(cluster); err != nil {
		return errors.Wrap(err, "failed reconciling cluster config")
	}

	return nil
}

func (s *Service) setStatus(cluster *eks.Cluster) error {
	switch *cluster.Status {
	case eks.ClusterStatusDeleting:
		s.scope.ControlPlane.Status.Initialized = false
		s.scope.ControlPlane.Status.Ready = false
	case eks.ClusterStatusFailed:
		s.scope.ControlPlane.Status.Initialized = true
		s.scope.ControlPlane.Status.Ready = false
		// TODO FailureReason
		failureMsg := fmt.Sprintf("EKS cluster in unexpected %s state", *cluster.Status)
		s.scope.ControlPlane.Status.FailureMessage = &failureMsg
	case eks.ClusterStatusActive:
		s.scope.ControlPlane.Status.Initialized = true
		s.scope.ControlPlane.Status.Ready = true
		s.scope.ControlPlane.Status.FailureMessage = nil
		// TODO FailureReason
	case eks.ClusterStatusCreating:
		s.scope.ControlPlane.Status.Initialized = true
		s.scope.ControlPlane.Status.Ready = false
	case eks.ClusterStatusUpdating:
		s.scope.ControlPlane.Status.Initialized = true
		s.scope.ControlPlane.Status.Ready = true
	default:
		return errors.Errorf("unexpected EKS cluster status %s", *cluster.Status)
	}
	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to update control plane")
	}
	return nil
}

// deleteCluster deletes an EKS cluster
func (s *Service) deleteCluster() error {
	cluster, err := s.describeEKSCluster()
	if err != nil {
		if awserrors.IsNotFound(err) {
			s.scope.V(4).Info("eks cluster does not exist")
			return nil
		}
		return errors.Wrap(err, "unable to describe eks cluster")
	}
	if cluster == nil {
		return nil
	}

	err = s.deleteClusterAndWait(cluster)
	if err != nil {
		record.Warnf(s.scope.ControlPlane, "FailedDeleteEKSCluster", "Failed to delete EKS cluster %s: %v", *cluster.Name, err)
		return errors.Wrap(err, "unable to delete EKS cluster")
	}
	record.Eventf(s.scope.ControlPlane, "SuccessfulDeleteEKSCluster", "Deleted EKS Cluster %s", *cluster.Name)

	return nil
}

func (s *Service) deleteClusterAndWait(cluster *eks.Cluster) error {
	s.scope.Info("Deleting EKS cluster", "eks-cluster", cluster.Name)

	input := &eks.DeleteClusterInput{
		Name: cluster.Name,
	}
	_, err := s.EKSClient.DeleteCluster(input)
	if err != nil {
		return errors.Wrapf(err, "failed to request delete of eks cluster %s", *cluster.Name)
	}

	waitInput := &eks.DescribeClusterInput{
		Name: cluster.Name,
	}

	err = s.EKSClient.WaitUntilClusterDeleted(waitInput)
	if err != nil {
		return errors.Wrapf(err, "failed waiting for eks cluster %s to delete", *cluster.Name)
	}

	return nil
}

func makeEksLogging(loggingSpec map[string]bool) *eks.Logging {
	var on = true
	var off = false
	enabled := eks.LogSetup{Enabled: &on}
	disabled := eks.LogSetup{Enabled: &off}
	for k, v := range loggingSpec {
		loggingKey := k
		if v {
			enabled.Types = append(enabled.Types, &loggingKey)
		} else {
			disabled.Types = append(disabled.Types, &loggingKey)
		}
	}
	var clusterLogging []*eks.LogSetup
	if len(enabled.Types) > 0 {
		clusterLogging = append(clusterLogging, &enabled)
	}
	if len(disabled.Types) > 0 {
		clusterLogging = append(clusterLogging, &disabled)
	}
	if len(clusterLogging) > 0 {
		return &eks.Logging{
			ClusterLogging: clusterLogging,
		}
	}
	return nil
}

func (s *Service) createCluster() (*eks.Cluster, error) {
	// TODO: Do we need to just add the private subnets?
	subnets := s.scope.Subnets()
	if len(subnets) < 2 {
		return nil, awserrors.NewFailedDependency(
			fmt.Sprintf("failed to create eks control plane %q, at least 2 subnets is required", s.scope.Name()),
		)
	}

	zones := subnets.GetUniqueZones()
	if len(zones) < 2 {
		return nil, awserrors.NewFailedDependency(
			fmt.Sprintf("failed to create eks control plane %q, subnets in at least 2 different az's are required", s.scope.Name()),
		)
	}

	subnetIds := make([]*string, 0)
	for _, subnet := range subnets {
		subnetIds = append(subnetIds, &subnet.ID)
	}

	// Make sure to use the MachineScope here to get the merger of AWSCluster and AWSMachine tags
	additionalTags := s.scope.AdditionalTags()

	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)
	tags := make(map[string]*string)
	for k, v := range additionalTags {
		tagValue := v
		tags[k] = &tagValue
	}

	version := strings.Replace(*s.scope.ControlPlane.Spec.Version, "v", "", -1)

	role, err := s.getIAMRole(*s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting control plane iam role: %s", *s.scope.ControlPlane.Spec.RoleName)
	}

	logging := makeEksLogging(s.scope.ControlPlane.Spec.Logging)
	input := &eks.CreateClusterInput{
		Name: &s.scope.Cluster.Name,
		//ClientRequestToken: aws.String(uuid.New().String()),
		Version: aws.String(version),
		Logging: logging,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds: subnetIds,
		},
		RoleArn: role.Arn,
		Tags:    tags,
	}

	var out *eks.CreateClusterOutput
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if out, err = s.EKSClient.CreateCluster(input); err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				return false, aerr
			}
			return false, err
		}
		return true, nil
	}, awserrors.ResourceNotFound); err != nil { //TODO: change the error that can be retried
		record.Warnf(s.scope.ControlPlane, "FaiedCreateEKSCluster", "Failed to create a new EKS cluster: %v", err)
		return nil, errors.Wrapf(err, "failed to create EKS cluster")
	}

	record.Eventf(s.scope.ControlPlane, "SuccessfulCreateEKSCluster", "Created a new EKS cluster %q", s.scope.Name())
	return out.Cluster, nil
}

func (s *Service) waitForClusterActive() (*eks.Cluster, error) {
	req := eks.DescribeClusterInput{
		Name: &s.scope.Cluster.Name,
	}
	if err := s.EKSClient.WaitUntilClusterActive(&req); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for eks control plane %q", *req.Name)
	}

	s.scope.Info("EKS control plane is now available", "cluster-name", s.scope.Cluster.Name)

	cluster, err := s.describeEKSCluster()
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe eks clusters")
	}
	if err := s.setStatus(cluster); err != nil {
		return nil, errors.Wrap(err, "failed to set status")
	}

	return cluster, nil
}

func (s *Service) reconcileClusterConfig(cluster *eks.Cluster) error {
	var needsUpdate bool
	input := eks.UpdateClusterConfigInput{Name: cluster.Name}

	if updateLogging := s.reconcileLogging(cluster.Logging); updateLogging != nil {
		needsUpdate = true
		input.Logging = updateLogging
	}

	if needsUpdate {
		if err := input.Validate(); err != nil {
			return errors.Wrap(err, "created invalid UpdateClusterConfigInput")
		}
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateClusterConfig(&input); err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					return false, aerr
				}
				return false, err
			}
			return true, nil
		}); err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "failed to update the EKS control plane: %v", err)
			return errors.Wrapf(err, "failed to update EKS cluster")
		}
	}
	return nil
}

func (s *Service) reconcileLogging(logging *eks.Logging) *eks.Logging {
	for _, logSetup := range logging.ClusterLogging {
		for _, l := range logSetup.Types {
			if enabled, ok := s.scope.ControlPlane.Spec.Logging[*l]; ok && enabled != *logSetup.Enabled {
				return makeEksLogging(s.scope.ControlPlane.Spec.Logging)
			}
		}
	}
	return nil
}

func (s *Service) reconcileClusterVersion(cluster *eks.Cluster) error {
	specVersion := version.MustParseGeneric(*s.scope.ControlPlane.Spec.Version)
	clusterVersion := version.MustParseGeneric(*cluster.Version)
	if clusterVersion.LessThan(specVersion) {
		nextVersion := clusterVersion.WithMinor(clusterVersion.Minor() + 1)
		nextVersionString := fmt.Sprintf("%d.%d", nextVersion.Major(), nextVersion.Minor())

		input := &eks.UpdateClusterVersionInput{
			Name:    cluster.Name,
			Version: &nextVersionString,
		}

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateClusterVersion(input); err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					return false, aerr
				}
				return false, err
			}
			return true, nil
		}); err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "failed to update the EKS control plane: %v", err)
			return errors.Wrapf(err, "failed to update EKS cluster")
		}
	}
	return nil
}

func (s *Service) waitForClusterUpdate() (*eks.Cluster, error) {
	cluster, err := s.waitForClusterActive()
	if err != nil {
		return nil, err
	}

	record.Eventf(s.scope.ControlPlane, "SuccessfulUpdateEKSControlPlane", "Upgraded control plane to %s", *cluster.Version)
	return cluster, nil
}

func (s *Service) describeEKSCluster() (*eks.Cluster, error) {
	input := &eks.DescribeClusterInput{
		Name: &s.scope.Cluster.Name,
	}

	out, err := s.EKSClient.DescribeCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				return nil, nil
			default:
				return nil, errors.Wrap(err, "failed to describe cluster")
			}
		} else {
			return nil, errors.Wrap(err, "failed to describe cluster")
		}
	}

	return out.Cluster, nil
}

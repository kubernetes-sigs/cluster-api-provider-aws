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
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/version"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/internal/tristate"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
)

func (s *Service) reconcileCluster(ctx context.Context) error {
	s.scope.V(2).Info("Reconciling EKS cluster")

	eksClusterName := s.scope.KubernetesClusterName()

	cluster, err := s.describeEKSCluster(eksClusterName)
	if err != nil {
		return errors.Wrap(err, "failed to describe eks clusters")
	}

	if cluster == nil {
		cluster, err = s.createCluster(eksClusterName)
		if err != nil {
			return errors.Wrap(err, "failed to create cluster")
		}
		s.scope.Info("Created EKS cluster in AWS", "cluster-name", eksClusterName)
	} else {
		tagKey := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())
		ownedTag := cluster.Tags[tagKey]
		if ownedTag == nil {
			return fmt.Errorf("checking owner of %s is %s: %w", s.scope.KubernetesClusterName(), s.scope.Name(), err)
		}

		s.scope.V(2).Info("Found owned EKS cluster in AWS", "cluster-name", eksClusterName)
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

	s.scope.V(2).Info("EKS Control Plane active", "endpoint", *cluster.Endpoint)

	s.scope.ControlPlane.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
		Host: *cluster.Endpoint,
		Port: 443,
	}

	if err := s.reconcileSecurityGroups(cluster); err != nil {
		return errors.Wrap(err, "failed reconciling security groups")
	}

	if err := s.reconcileKubeconfig(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling kubeconfig")
	}

	if err := s.reconcileAdditionalKubeconfigs(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling additional kubeconfigs")
	}

	if err := s.reconcileClusterVersion(cluster); err != nil {
		return errors.Wrap(err, "failed reconciling cluster version")
	}

	if err := s.reconcileClusterConfig(cluster); err != nil {
		return errors.Wrap(err, "failed reconciling cluster config")
	}

	if err := s.reconcileTags(cluster); err != nil {
		return errors.Wrap(err, "failed updating cluster tags")
	}

	if err := s.reconcileOIDCProvider(cluster); err != nil {
		return errors.Wrap(err, "failed reconciling OIDC provider for cluster")
	}

	return nil
}

func (s *Service) setStatus(cluster *eks.Cluster) error {
	switch *cluster.Status {
	case eks.ClusterStatusDeleting:
		s.scope.ControlPlane.Status.Ready = false
	case eks.ClusterStatusFailed:
		s.scope.ControlPlane.Status.Ready = false
		// TODO FailureReason
		failureMsg := fmt.Sprintf("EKS cluster in unexpected %s state", *cluster.Status)
		s.scope.ControlPlane.Status.FailureMessage = &failureMsg
	case eks.ClusterStatusActive:
		s.scope.ControlPlane.Status.Ready = true
		s.scope.ControlPlane.Status.FailureMessage = nil
		// TODO FailureReason
	case eks.ClusterStatusCreating:
		s.scope.ControlPlane.Status.Ready = false
	case eks.ClusterStatusUpdating:
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
	eksClusterName := s.scope.KubernetesClusterName()

	if eksClusterName == "" {
		s.scope.V(2).Info("no EKS cluster name, skipping EKS cluster deletion")
		return nil
	}

	cluster, err := s.describeEKSCluster(eksClusterName)
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
		record.Warnf(s.scope.ControlPlane, "FailedDeleteEKSCluster", "Failed to delete EKS cluster %s: %v", s.scope.KubernetesClusterName(), err)
		return errors.Wrap(err, "unable to delete EKS cluster")
	}
	record.Eventf(s.scope.ControlPlane, "SuccessfulDeleteEKSCluster", "Deleted EKS Cluster %s", s.scope.KubernetesClusterName())

	return nil
}

func (s *Service) deleteClusterAndWait(cluster *eks.Cluster) error {
	s.scope.Info("Deleting EKS cluster", "cluster-name", s.scope.KubernetesClusterName())

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

func makeEksEncryptionConfigs(encryptionConfig *ekscontrolplanev1.EncryptionConfig) []*eks.EncryptionConfig {
	if encryptionConfig == nil {
		return []*eks.EncryptionConfig{}
	}
	return []*eks.EncryptionConfig{{
		Provider: &eks.Provider{
			KeyArn: encryptionConfig.Provider,
		},
		Resources: encryptionConfig.Resources,
	}}
}

func makeVpcConfig(subnets infrav1.Subnets, endpointAccess ekscontrolplanev1.EndpointAccess) (*eks.VpcConfigRequest, error) {
	// TODO: Do we need to just add the private subnets?
	if len(subnets) < 2 {
		return nil, awserrors.NewFailedDependency("at least 2 subnets is required")
	}

	zones := subnets.GetUniqueZones()
	if len(zones) < 2 {
		return nil, awserrors.NewFailedDependency("subnets in at least 2 different az's are required")
	}

	subnetIds := make([]*string, 0)
	for _, subnet := range subnets {
		subnetIds = append(subnetIds, &subnet.ID)
	}

	cidrs := make([]*string, 0)
	for _, cidr := range endpointAccess.PublicCIDRs {
		_, ipNet, err := net.ParseCIDR(*cidr)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't parse PublicCIDRs")
		}
		parsedCIDR := ipNet.String()
		cidrs = append(cidrs, &parsedCIDR)
	}

	vpcConfig := &eks.VpcConfigRequest{
		EndpointPublicAccess:  endpointAccess.Public,
		EndpointPrivateAccess: endpointAccess.Private,
		SubnetIds:             subnetIds,
	}

	if len(cidrs) > 0 {
		vpcConfig.PublicAccessCidrs = cidrs
	}

	return vpcConfig, nil
}

func makeEksLogging(loggingSpec *ekscontrolplanev1.ControlPlaneLoggingSpec) *eks.Logging {
	if loggingSpec == nil {
		return nil
	}
	var on = true
	var off = false
	var enabledTypes []string
	var disabledTypes []string

	appendToTypes := func(logType string, field bool) {
		if field {
			enabledTypes = append(enabledTypes, logType)
		} else {
			disabledTypes = append(disabledTypes, logType)
		}
	}

	appendToTypes(eks.LogTypeApi, loggingSpec.APIServer)
	appendToTypes(eks.LogTypeAudit, loggingSpec.Audit)
	appendToTypes(eks.LogTypeAuthenticator, loggingSpec.Authenticator)
	appendToTypes(eks.LogTypeControllerManager, loggingSpec.ControllerManager)
	appendToTypes(eks.LogTypeScheduler, loggingSpec.Scheduler)

	var clusterLogging []*eks.LogSetup
	if len(enabledTypes) > 0 {
		enabled := eks.LogSetup{
			Enabled: &on,
			Types:   aws.StringSlice(enabledTypes),
		}
		clusterLogging = append(clusterLogging, &enabled)
	}
	if len(disabledTypes) > 0 {
		disabled := eks.LogSetup{
			Enabled: &off,
			Types:   aws.StringSlice(disabledTypes),
		}
		clusterLogging = append(clusterLogging, &disabled)
	}
	if len(clusterLogging) > 0 {
		return &eks.Logging{
			ClusterLogging: clusterLogging,
		}
	}
	return nil
}

func (s *Service) createCluster(eksClusterName string) (*eks.Cluster, error) {
	logging := makeEksLogging(s.scope.ControlPlane.Spec.Logging)
	encryptionConfigs := makeEksEncryptionConfigs(s.scope.ControlPlane.Spec.EncryptionConfig)
	vpcConfig, err := makeVpcConfig(s.scope.Subnets(), s.scope.ControlPlane.Spec.EndpointAccess)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create vpc config for cluster")
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

	role, err := s.GetIAMRole(*s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting control plane iam role: %s", *s.scope.ControlPlane.Spec.RoleName)
	}

	v := versionToEKS(parseEKSVersion(*s.scope.ControlPlane.Spec.Version))

	input := &eks.CreateClusterInput{
		Name: aws.String(eksClusterName),
		//ClientRequestToken: aws.String(uuid.New().String()),
		Version:            aws.String(v),
		Logging:            logging,
		EncryptionConfig:   encryptionConfigs,
		ResourcesVpcConfig: vpcConfig,
		RoleArn:            role.Arn,
		Tags:               tags,
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
	eksClusterName := s.scope.KubernetesClusterName()
	req := eks.DescribeClusterInput{
		Name: aws.String(eksClusterName),
	}
	if err := s.EKSClient.WaitUntilClusterActive(&req); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for eks control plane %q", *req.Name)
	}

	s.scope.Info("EKS control plane is now available", "cluster-name", eksClusterName)

	cluster, err := s.describeEKSCluster(eksClusterName)
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
	input := eks.UpdateClusterConfigInput{Name: aws.String(s.scope.KubernetesClusterName())}

	if updateLogging := s.reconcileLogging(cluster.Logging); updateLogging != nil {
		needsUpdate = true
		input.Logging = updateLogging
	}

	updateVpcConfig, err := s.reconcileVpcConfig(cluster.ResourcesVpcConfig)
	if err != nil {
		return errors.Wrap(err, "couldn't create vpc config for cluster")
	}
	if updateVpcConfig != nil {
		needsUpdate = true
		input.ResourcesVpcConfig = updateVpcConfig
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
			record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "failed to update the EKS control plane %s: %v", s.scope.KubernetesClusterName(), err)
			return errors.Wrapf(err, "failed to update EKS cluster")
		}
	}
	return nil
}

func (s *Service) reconcileLogging(logging *eks.Logging) *eks.Logging {
	for _, logSetup := range logging.ClusterLogging {
		for _, l := range logSetup.Types {
			enabled := s.scope.ControlPlane.Spec.Logging.IsLogEnabled(*l)
			if enabled != *logSetup.Enabled {
				return makeEksLogging(s.scope.ControlPlane.Spec.Logging)
			}
		}
	}
	return nil
}

func publicAccessCIDRsEqual(as []*string, bs []*string) bool {
	all := "0.0.0.0/0"
	if len(as) == 0 {
		as = []*string{&all}
	}
	if len(bs) == 0 {
		bs = []*string{&all}
	}
	return sets.NewString(aws.StringValueSlice(as)...).Equal(
		sets.NewString(aws.StringValueSlice(bs)...),
	)
}

func (s *Service) reconcileVpcConfig(vpcConfig *eks.VpcConfigResponse) (*eks.VpcConfigRequest, error) {
	endpointAccess := s.scope.ControlPlane.Spec.EndpointAccess
	updatedVpcConfig, err := makeVpcConfig(s.scope.Subnets(), endpointAccess)
	if err != nil {
		return nil, err
	}
	needsUpdate := !tristate.EqualWithDefault(false, vpcConfig.EndpointPrivateAccess, updatedVpcConfig.EndpointPrivateAccess) ||
		!tristate.EqualWithDefault(true, vpcConfig.EndpointPublicAccess, updatedVpcConfig.EndpointPublicAccess) ||
		!publicAccessCIDRsEqual(vpcConfig.PublicAccessCidrs, updatedVpcConfig.PublicAccessCidrs)
	if needsUpdate {
		return &eks.VpcConfigRequest{
			EndpointPublicAccess:  updatedVpcConfig.EndpointPublicAccess,
			EndpointPrivateAccess: updatedVpcConfig.EndpointPrivateAccess,
			PublicAccessCidrs:     updatedVpcConfig.PublicAccessCidrs,
		}, nil
	}
	return nil, nil
}

func parseEKSVersion(raw string) *version.Version {
	v := version.MustParseGeneric(raw)
	return version.MustParseGeneric(fmt.Sprintf("%d.%d", v.Major(), v.Minor()))
}

func versionToEKS(v *version.Version) string {
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

func (s *Service) reconcileClusterVersion(cluster *eks.Cluster) error {
	specVersion := parseEKSVersion(*s.scope.ControlPlane.Spec.Version)
	clusterVersion := version.MustParseGeneric(*cluster.Version)

	if clusterVersion.LessThan(specVersion) {
		// NOTE: you can only upgrade increments of minor versions. If you want to upgrade 1.14 to 1.16 we
		// need to go 1.14-> 1.15 and then 1.15 -> 1.16.
		nextVersionString := versionToEKS(clusterVersion.WithMinor(clusterVersion.Minor() + 1))

		input := &eks.UpdateClusterVersionInput{
			Name:    aws.String(s.scope.KubernetesClusterName()),
			Version: &nextVersionString,
		}

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateClusterVersion(input); err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					return false, aerr
				}
				return false, err
			}
			record.Eventf(s.scope.ControlPlane, "SuccessfulUpdateEKSControlPlane", "Updated EKS control plane %s to version %s", s.scope.KubernetesClusterName(), nextVersionString)
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

	record.Eventf(s.scope.ControlPlane, "SuccessfulUpdateEKSControlPlane", "Updated EKS control plane %s", s.scope.KubernetesClusterName())
	return cluster, nil
}

func (s *Service) describeEKSCluster(eksClusterName string) (*eks.Cluster, error) {
	input := &eks.DescribeClusterInput{
		Name: aws.String(eksClusterName),
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

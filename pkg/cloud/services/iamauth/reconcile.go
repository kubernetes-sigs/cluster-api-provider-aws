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

package iamauth

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// ReconcileIAMAuthenticator is used to create the aws-iam-authenticator in a cluster.
func (s *Service) ReconcileIAMAuthenticator(ctx context.Context) error {
	s.scope.Info(
		"Reconciling aws-iam-authenticator configuration",
		"cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()),
	)

	remoteClient, err := s.scope.RemoteClient()
	if err != nil {
		s.scope.Error(err, "getting client for remote cluster")
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

	authBackend, err := NewBackend(s.backend, remoteClient)
	if err != nil {
		return fmt.Errorf("getting aws-iam-authenticator backend: %w", err)
	}

	// Discover node roles from worker templates.
	nodeRoles, err := s.getRolesForWorkers(ctx)
	if err != nil {
		s.scope.Error(err, "getting roles for remote workers")
		return fmt.Errorf("getting roles for remote workers: %w", err)
	}

	iamCfg := s.scope.IAMAuthConfig()

	// Compose the desired role set: node roles ∪ iamCfg.RoleMappings.
	// Node roles MUST be included — the reconcile is a full replace of the
	// aws-auth backend, so omitting them would revoke kubelet auth on all
	// worker nodes.
	desiredRoles := make([]ekscontrolplanev1.RoleMapping, 0, len(nodeRoles)+len(iamCfg.RoleMappings))
	for roleName := range nodeRoles {
		roleARN, err := s.getARNForRole(ctx, roleName)
		if err != nil {
			return fmt.Errorf("failed to get ARN for role %s: %w", roleName, err)
		}
		desiredRoles = append(desiredRoles, ekscontrolplanev1.RoleMapping{
			RoleARN: roleARN,
			KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
				UserName: EC2NodeUserName,
				Groups:   NodeGroups,
			},
		})
	}
	desiredRoles = append(desiredRoles, iamCfg.RoleMappings...)

	// Dedup by ARN and sort deterministically. nodeRoles is a map (unordered)
	// and a user-configured RoleMapping may collide with a discovered node
	// role's ARN; without dedup+sort the CM backend would churn the aws-auth
	// ConfigMap on every reconcile and the CRD backend would create duplicate
	// IAMIdentityMapping CRs. User-configured mappings are appended after node
	// roles so they win on ARN collision (explicit intent overrides discovery).
	desiredRoles = dedupAndSortRoles(desiredRoles)
	desiredUsers := dedupAndSortUsers(iamCfg.UserMappings)

	s.scope.Debug(
		"Reconciling IAM authenticator mappings",
		"node-roles", len(nodeRoles),
		"user-roles", len(iamCfg.RoleMappings),
		"user-mappings", len(desiredUsers),
	)

	if err := authBackend.ReconcileMappings(desiredRoles, desiredUsers); err != nil {
		return fmt.Errorf("reconciling iam mappings: %w", err)
	}

	s.scope.Info("Reconciled aws-iam-authenticator configuration", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	return nil
}

func (s *Service) getARNForRole(ctx context.Context, role string) (string, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(role),
	}
	out, err := s.IAMClient.GetRole(ctx, input)
	if err != nil {
		return "", errors.Wrap(err, "unable to get role")
	}
	if out.Role == nil || out.Role.Arn == nil {
		return "", fmt.Errorf("role %s not found or ARN is missing", role)
	}
	return *out.Role.Arn, nil
}

func (s *Service) getRolesForWorkers(ctx context.Context) (map[string]struct{}, error) {
	allRoles := map[string]struct{}{}
	if err := s.getRolesForMachineDeployments(ctx, allRoles); err != nil {
		return nil, fmt.Errorf("failed to get roles from machine deployments %w", err)
	}
	if err := s.getRolesForMachinePools(ctx, allRoles); err != nil {
		return nil, fmt.Errorf("failed to get roles from machine pools %w", err)
	}
	return allRoles, nil
}

func (s *Service) getRolesForMachineDeployments(ctx context.Context, allRoles map[string]struct{}) error {
	deploymentList := &clusterv1.MachineDeploymentList{}
	selectors := []client.ListOption{
		client.InNamespace(s.scope.Namespace()),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: s.scope.Name(),
		},
	}
	err := s.client.List(ctx, deploymentList, selectors...)
	if err != nil {
		return fmt.Errorf("failed to list machine deployments for cluster %s/%s: %w", s.scope.Namespace(), s.scope.Name(), err)
	}

	for _, deployment := range deploymentList.Items {
		ref := deployment.Spec.Template.Spec.InfrastructureRef
		if ref.Kind != "AWSMachineTemplate" {
			continue
		}
		awsMachineTemplate := &infrav1.AWSMachineTemplate{}
		err := s.client.Get(ctx, client.ObjectKey{
			Name:      ref.Name,
			Namespace: s.scope.Namespace(),
		}, awsMachineTemplate)
		if err != nil {
			return fmt.Errorf("failed to get AWSMachine %s/%s: %w", s.scope.Namespace(), ref.Name, err)
		}
		instanceProfile := awsMachineTemplate.Spec.Template.Spec.IAMInstanceProfile
		if _, ok := allRoles[instanceProfile]; !ok && instanceProfile != "" {
			allRoles[instanceProfile] = struct{}{}
		}
	}
	return nil
}

func (s *Service) getRolesForMachinePools(ctx context.Context, allRoles map[string]struct{}) error {
	machinePoolList := &clusterv1.MachinePoolList{}
	selectors := []client.ListOption{
		client.InNamespace(s.scope.Namespace()),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: s.scope.Name(),
		},
	}
	err := s.client.List(ctx, machinePoolList, selectors...)
	if err != nil {
		return fmt.Errorf("failed to list machine pools for cluster %s/%s: %w", s.scope.Namespace(), s.scope.Name(), err)
	}
	for _, pool := range machinePoolList.Items {
		ref := pool.Spec.Template.Spec.InfrastructureRef
		switch ref.Kind {
		case "AWSMachinePool":
			if err := s.getRolesForAWSMachinePool(ctx, ref, allRoles); err != nil {
				return err
			}
		case "AWSManagedMachinePool":
			if err := s.getRolesForAWSManagedMachinePool(ctx, ref, allRoles); err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

func (s *Service) getRolesForAWSMachinePool(ctx context.Context, ref clusterv1.ContractVersionedObjectReference, allRoles map[string]struct{}) error {
	awsMachinePool := &expinfrav1.AWSMachinePool{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      ref.Name,
		Namespace: s.scope.Namespace(),
	}, awsMachinePool)
	if err != nil {
		return fmt.Errorf("failed to get AWSMachine %s/%s: %w", s.scope.Namespace(), ref.Name, err)
	}
	instanceProfile := awsMachinePool.Spec.AWSLaunchTemplate.IamInstanceProfile
	if _, ok := allRoles[instanceProfile]; !ok && instanceProfile != "" {
		allRoles[instanceProfile] = struct{}{}
	}
	return nil
}

func (s *Service) getRolesForAWSManagedMachinePool(ctx context.Context, ref clusterv1.ContractVersionedObjectReference, allRoles map[string]struct{}) error {
	awsManagedMachinePool := &expinfrav1.AWSManagedMachinePool{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      ref.Name,
		Namespace: s.scope.Namespace(),
	}, awsManagedMachinePool)
	if err != nil {
		return fmt.Errorf("failed to get AWSMachine %s/%s: %w", s.scope.Namespace(), ref.Name, err)
	}
	instanceProfile := awsManagedMachinePool.Spec.RoleName
	if _, ok := allRoles[instanceProfile]; !ok && instanceProfile != "" {
		allRoles[instanceProfile] = struct{}{}
	}
	return nil
}

// dedupAndSortRoles collapses RoleMapping entries with duplicate RoleARNs and
// returns the result sorted by RoleARN. Later entries win on collision, which
// makes user-configured mappings from iamCfg.RoleMappings override any
// same-ARN entry that came from node-role discovery.
func dedupAndSortRoles(in []ekscontrolplanev1.RoleMapping) []ekscontrolplanev1.RoleMapping {
	byARN := make(map[string]ekscontrolplanev1.RoleMapping, len(in))
	for _, m := range in {
		byARN[m.RoleARN] = m
	}
	out := make([]ekscontrolplanev1.RoleMapping, 0, len(byARN))
	for _, m := range byARN {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RoleARN < out[j].RoleARN })
	return out
}

// dedupAndSortUsers is the UserMapping analog of dedupAndSortRoles. It
// collapses duplicate UserARNs (later wins) and returns entries sorted by
// UserARN.
func dedupAndSortUsers(in []ekscontrolplanev1.UserMapping) []ekscontrolplanev1.UserMapping {
	byARN := make(map[string]ekscontrolplanev1.UserMapping, len(in))
	for _, m := range in {
		byARN[m.UserARN] = m
	}
	out := make([]ekscontrolplanev1.UserMapping, 0, len(byARN))
	for _, m := range byARN {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UserARN < out[j].UserARN })
	return out
}

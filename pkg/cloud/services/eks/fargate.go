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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func requeueProfileUpdating() reconcile.Result {
	return reconcile.Result{RequeueAfter: 10 * time.Second}
}

func requeueRoleUpdating() reconcile.Result {
	return reconcile.Result{RequeueAfter: 10 * time.Second}
}

// Reconcile is the entrypoint for FargateProfile reconciliation
func (s *FargateService) Reconcile() (reconcile.Result, error) {
	s.scope.V(2).Info("Reconciling EKS fargate profile")

	requeue, err := s.reconcileFargateIAMRole()
	if err != nil {
		conditions.MarkFalse(
			s.scope.FargateProfile,
			infrav1exp.IAMFargateRolesReadyCondition,
			infrav1exp.IAMFargateRolesReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			err.Error(),
		)
		return reconcile.Result{}, err
	}
	// When the role is updated, we requeue to let e.g. trust relationship
	// propagate
	if requeue {
		return requeueRoleUpdating(), nil
	}

	conditions.MarkTrue(s.scope.FargateProfile, infrav1exp.IAMFargateRolesReadyCondition)

	requeue, err = s.reconcileFargateProfile()
	if err != nil {
		conditions.MarkFalse(
			s.scope.FargateProfile,
			clusterv1.ReadyCondition,
			infrav1exp.EKSFargateReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			err.Error(),
		)
		return reconcile.Result{}, err
	}
	if requeue {
		return requeueProfileUpdating(), nil
	}

	return reconcile.Result{}, nil
}

func (s *FargateService) reconcileFargateProfile() (requeue bool, err error) {
	eksClusterName := s.scope.KubernetesClusterName()
	profileName := s.scope.FargateProfile.Spec.ProfileName

	profile, err := s.describeFargateProfile()
	if err != nil {
		return false, errors.Wrap(err, "failed to describe profile")
	}

	if profile == nil {
		profile, err = s.createFargateProfile()
		if err != nil {
			return false, errors.Wrap(err, "failed to create profile")
		}
		// Force status to creating
		profile.Status = aws.String(eks.FargateProfileStatusCreating)
		s.scope.Info("Created EKS fargate profile", "cluster-name", eksClusterName, "profile-name", profileName)
	} else {
		tagKey := infrav1.ClusterAWSCloudProviderTagKey(s.scope.ClusterName())
		ownedTag := profile.Tags[tagKey]
		if ownedTag == nil {
			return false, errors.Wrapf(err, "owned tag not found for this cluster")
		}
		s.scope.V(2).Info("Found owned EKS fargate profile", "cluster-name", eksClusterName, "profile-name", profileName)
	}

	return s.handleStatus(profile), nil
}

func (s *FargateService) handleStatus(profile *eks.FargateProfile) (requeue bool) {
	s.V(2).Info("fargate profile", "status", *profile.Status)
	switch *profile.Status {
	case eks.FargateProfileStatusCreating:
		s.scope.FargateProfile.Status.Ready = false
		if conditions.IsTrue(s.scope.FargateProfile, infrav1exp.EKSFargateDeletingCondition) {
			conditions.MarkFalse(s.scope.FargateProfile, infrav1exp.EKSFargateDeletingCondition, infrav1exp.EKSFargateCreatingReason, clusterv1.ConditionSeverityInfo, "")
		}
		if !conditions.IsTrue(s.scope.FargateProfile, infrav1exp.EKSFargateCreatingCondition) {
			record.Eventf(s.scope.FargateProfile, "InitiatedCreateEKSFargateProfile", "Started creating EKS fargate profile %s", s.scope.FargateProfile.Spec.ProfileName)
			conditions.MarkTrue(s.scope.FargateProfile, infrav1exp.EKSFargateCreatingCondition)
		}
		conditions.MarkFalse(s.scope.FargateProfile, infrav1exp.EKSFargateProfileReadyCondition, infrav1exp.EKSFargateCreatingReason, clusterv1.ConditionSeverityInfo, "")
	case eks.FargateProfileStatusCreateFailed, eks.FargateProfileStatusDeleteFailed:
		s.scope.FargateProfile.Status.Ready = false
		s.scope.FargateProfile.Status.FailureMessage = aws.String(fmt.Sprintf("unexpected profile status: %s", *profile.Status))
		reason := capierrors.MachineStatusError(infrav1exp.EKSFargateFailedReason)
		s.scope.FargateProfile.Status.FailureReason = &reason
		conditions.MarkFalse(s.scope.FargateProfile, infrav1exp.EKSFargateProfileReadyCondition, infrav1exp.EKSFargateFailedReason, clusterv1.ConditionSeverityError, "")
	case eks.FargateProfileStatusActive:
		s.scope.FargateProfile.Status.Ready = true
		if conditions.IsTrue(s.scope.FargateProfile, infrav1exp.EKSFargateCreatingCondition) {
			record.Eventf(s.scope.FargateProfile, "SuccessfulCreateEKSFargateProfile", "Created new EKS fargate profile %s", s.scope.FargateProfile.Spec.ProfileName)
			conditions.MarkFalse(s.scope.FargateProfile, infrav1exp.EKSFargateCreatingCondition, infrav1exp.EKSFargateCreatedReason, clusterv1.ConditionSeverityInfo, "")
		}
		conditions.MarkTrue(s.scope.FargateProfile, infrav1exp.EKSFargateProfileReadyCondition)
	case eks.FargateProfileStatusDeleting:
		s.scope.FargateProfile.Status.Ready = false
		if !conditions.IsTrue(s.scope.FargateProfile, infrav1exp.EKSFargateDeletingCondition) {
			record.Eventf(s.scope.FargateProfile, "InitiatedDeleteEKSFargateProfile", "Started deleting EKS fargate profile %s", s.scope.FargateProfile.Spec.ProfileName)
			conditions.MarkTrue(s.scope.FargateProfile, infrav1exp.EKSFargateDeletingCondition)
		}
		conditions.MarkFalse(s.scope.FargateProfile, infrav1exp.EKSFargateProfileReadyCondition, infrav1exp.EKSFargateDeletingReason, clusterv1.ConditionSeverityInfo, "")
	}
	switch *profile.Status {
	case eks.FargateProfileStatusCreating, eks.FargateProfileStatusDeleting:
		return true
	default:
		return false
	}
}

// ReconcileDelete is the entrypoint for FargateProfile reconciliation
func (s *FargateService) ReconcileDelete() (reconcile.Result, error) {
	s.scope.V(2).Info("Reconciling EKS fargate profile deletion")

	requeue, err := s.deleteFargateProfile()
	if err != nil {
		conditions.MarkFalse(
			s.scope.FargateProfile,
			clusterv1.ReadyCondition,
			infrav1exp.EKSFargateReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			err.Error(),
		)
		return reconcile.Result{}, err
	}

	if requeue {
		return requeueProfileUpdating(), nil
	}

	err = s.deleteFargateIAMRole()
	if err != nil {
		conditions.MarkFalse(
			s.scope.FargateProfile,
			infrav1exp.IAMFargateRolesReadyCondition,
			infrav1exp.IAMFargateRolesReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			err.Error(),
		)
	}
	return reconcile.Result{}, err
}

func (s *FargateService) describeFargateProfile() (*eks.FargateProfile, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	profileName := s.scope.FargateProfile.Spec.ProfileName
	input := &eks.DescribeFargateProfileInput{
		ClusterName:        aws.String(eksClusterName),
		FargateProfileName: aws.String(profileName),
	}

	out, err := s.EKSClient.DescribeFargateProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == eks.ErrCodeResourceNotFoundException {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to describe fargate profile")
	}

	return out.FargateProfile, nil
}

func (s *FargateService) createFargateProfile() (*eks.FargateProfile, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	profileName := s.scope.FargateProfile.Spec.ProfileName

	additionalTags := s.scope.AdditionalTags()

	roleArn, err := s.roleArn()
	if err != nil {
		return nil, err
	}

	tags := ngTags(s.scope.ClusterName(), additionalTags)

	subnets := s.scope.FargateProfile.Spec.SubnetIDs
	if len(subnets) == 0 {
		subnets = []string{}
		for _, s := range s.scope.ControlPlane.Spec.NetworkSpec.Subnets.FilterPrivate() {
			subnets = append(subnets, s.ID)
		}
	}

	selectors := []*eks.FargateProfileSelector{}
	for _, s := range s.scope.FargateProfile.Spec.Selectors {
		selectors = append(selectors, &eks.FargateProfileSelector{
			Labels:    aws.StringMap(s.Labels),
			Namespace: aws.String(s.Namespace),
		})
	}

	input := &eks.CreateFargateProfileInput{
		ClusterName:         aws.String(eksClusterName),
		FargateProfileName:  aws.String(profileName),
		PodExecutionRoleArn: roleArn,
		Subnets:             aws.StringSlice(subnets),
		Tags:                aws.StringMap(tags),
		Selectors:           selectors,
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "created invalid CreateFargateProfileInput")
	}

	out, err := s.EKSClient.CreateFargateProfile(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fargate profile")
	}

	return out.FargateProfile, nil
}

func (s *FargateService) deleteFargateProfile() (requeue bool, err error) {
	eksClusterName := s.scope.KubernetesClusterName()
	profileName := s.scope.FargateProfile.Spec.ProfileName

	profile, err := s.describeFargateProfile()
	if err != nil {
		return false, errors.Wrap(err, "failed to describe profile")
	}
	if profile == nil {
		if conditions.IsTrue(s.scope.FargateProfile, infrav1exp.EKSFargateDeletingCondition) {
			record.Eventf(s.scope.FargateProfile, "SuccessfulDeleteEKSFargateProfile", "Deleted EKS fargate profile %s", s.scope.FargateProfile.Spec.ProfileName)
			conditions.MarkFalse(s.scope.FargateProfile, infrav1exp.EKSFargateDeletingCondition, infrav1exp.EKSFargateDeletedReason, clusterv1.ConditionSeverityInfo, "")
		}
		conditions.MarkFalse(s.scope.FargateProfile, infrav1exp.EKSFargateProfileReadyCondition, infrav1exp.EKSFargateDeletedReason, clusterv1.ConditionSeverityInfo, "")
		return false, nil
	}

	switch aws.StringValue(profile.Status) {
	case eks.FargateProfileStatusCreating, eks.FargateProfileStatusDeleting, eks.FargateProfileStatusDeleteFailed:
		return s.handleStatus(profile), nil
	case eks.FargateProfileStatusActive, eks.FargateProfileStatusCreateFailed:
	}

	input := &eks.DeleteFargateProfileInput{
		ClusterName:        aws.String(eksClusterName),
		FargateProfileName: aws.String(profileName),
	}
	if err := input.Validate(); err != nil {
		return false, errors.Wrap(err, "created invalid DeleteFargateProfileInput")
	}

	out, err := s.EKSClient.DeleteFargateProfile(input)
	if err != nil {
		return false, errors.Wrap(err, "failed to delete fargate profile")
	}

	profile = out.FargateProfile
	profile.Status = aws.String(eks.FargateProfileStatusDeleting)

	return s.handleStatus(profile), nil
}

func (s *FargateService) roleArn() (*string, error) {
	var role *iam.Role
	if s.scope.RoleName() != "" {
		var err error
		role, err = s.GetIAMRole(s.scope.RoleName())
		if err != nil {
			return nil, errors.Wrapf(err, "error getting fargate profile IAM role: %s", s.scope.RoleName())
		}
	}
	return role.Arn, nil
}

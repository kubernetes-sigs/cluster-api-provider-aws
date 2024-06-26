package machinepool

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api/util/conditions"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

type Service struct {
	ASGInterface services.ASGInterface
	EC2Interface services.EC2Interface
}

// NewService returns a new service given the ec2 api client.
func NewService(EC2Interface services.EC2Interface, ASGInterface services.ASGInterface) *Service {
	return &Service{
		ASGInterface: ASGInterface,
		EC2Interface: EC2Interface,
	}
}

// ReconcileLifecycleHooks periodically reconciles a lifecycle hook for the ASG.
func (s *Service) ReconcileLifecycleHooks(scope scope.LifecycleHookScope) error {
	for _, hook := range scope.GetLifecycleHooks() {
		if err := s.reconcileLifecycleHook(scope, &hook); err != nil {
			return err
		}
	}

	// Get a list of lifecycle hooks that are registered with the ASG but not defined in the MachinePool and delete them.
	hooks, err := s.ASGInterface.GetLifecycleHooks(scope)
	if err != nil {
		return err
	}
	for _, hook := range hooks {
		found := false
		for _, definedHook := range scope.GetLifecycleHooks() {
			if hook.Name == definedHook.Name {
				found = true
				break
			}
		}
		if !found {
			scope.Info("Deleting lifecycle hook", "hook", hook.Name)
			if err := s.ASGInterface.DeleteLifecycleHook(scope, hook); err != nil {
				conditions.MarkFalse(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition, expinfrav1.LifecycleHookDeletionFailedReason, clusterv1.ConditionSeverityError, err.Error())
				return err
			}
		}
	}

	return nil
}

func (s *Service) reconcileLifecycleHook(scope scope.LifecycleHookScope, hook *expinfrav1.AWSLifecycleHook) error {
	scope.Info("Checking for existing lifecycle hook")
	existingHook, err := s.ASGInterface.GetLifecycleHook(scope, hook)
	if err != nil {
		conditions.MarkUnknown(scope.GetMachinePool(), expinfrav1.LifecycleHookReadyCondition, expinfrav1.LifecycleHookNotFoundReason, err.Error())
		return err
	}

	if existingHook == nil {
		scope.Info("Creating lifecycle hook")
		if err := s.ASGInterface.CreateLifecycleHook(scope, hook); err != nil {
			conditions.MarkFalse(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition, expinfrav1.LifecycleHookCreationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
		return nil
	}

	// If the lifecycle hook exists, we need to check if it's up to date
	needsUpdate := s.ASGInterface.LifecycleHookNeedsUpdate(scope, existingHook, hook)

	if needsUpdate {
		scope.Info("Updating lifecycle hook")
		if err := s.ASGInterface.UpdateLifecycleHook(scope, hook); err != nil {
			conditions.MarkFalse(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition, expinfrav1.LifecycleHookUpdateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
	}

	conditions.MarkTrue(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition)
	return nil
}

// ReconcileLaunchTemplate reconciles a launch template and triggers instance refresh conditionally, depending on
// changes.
//
//nolint:gocyclo
func (s *Service) ReconcileLaunchTemplate(
	scope scope.LaunchTemplateScope,
	ec2svc services.EC2Interface,
	canUpdateLaunchTemplate func() (bool, error),
	runPostLaunchTemplateUpdateOperation func() error,
) error {
	bootstrapData, bootstrapDataSecretKey, err := scope.GetRawBootstrapData()
	if err != nil {
		record.Eventf(scope.GetMachinePool(), corev1.EventTypeWarning, "FailedGetBootstrapData", err.Error())
		return err
	}
	bootstrapDataHash := userdata.ComputeHash(bootstrapData)

	scope.Info("checking for existing launch template")
	launchTemplate, launchTemplateUserDataHash, launchTemplateUserDataSecretKey, err := ec2svc.GetLaunchTemplate(scope.LaunchTemplateName())
	if err != nil {
		conditions.MarkUnknown(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, err.Error())
		return err
	}

	imageID, err := ec2svc.DiscoverLaunchTemplateAMI(scope)
	if err != nil {
		conditions.MarkFalse(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateCreateFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}

	if launchTemplate == nil {
		scope.Info("no existing launch template found, creating")
		launchTemplateID, err := ec2svc.CreateLaunchTemplate(scope, imageID, *bootstrapDataSecretKey, bootstrapData)
		if err != nil {
			conditions.MarkFalse(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateCreateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}

		scope.SetLaunchTemplateIDStatus(launchTemplateID)
		return scope.PatchObject()
	}

	// LaunchTemplateID is set during LaunchTemplate creation, but for a scenario such as `clusterctl move`, status fields become blank.
	// If launchTemplate already exists but LaunchTemplateID field in the status is empty, get the ID and update the status.
	if scope.GetLaunchTemplateIDStatus() == "" {
		launchTemplateID, err := ec2svc.GetLaunchTemplateID(scope.LaunchTemplateName())
		if err != nil {
			conditions.MarkUnknown(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, err.Error())
			return err
		}
		scope.SetLaunchTemplateIDStatus(launchTemplateID)
		return scope.PatchObject()
	}

	if scope.GetLaunchTemplateLatestVersionStatus() == "" {
		launchTemplateVersion, err := ec2svc.GetLaunchTemplateLatestVersion(scope.GetLaunchTemplateIDStatus())
		if err != nil {
			conditions.MarkUnknown(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, err.Error())
			return err
		}
		scope.SetLaunchTemplateLatestVersionStatus(launchTemplateVersion)
		return scope.PatchObject()
	}

	annotation, err := ec2.MachinePoolAnnotationJSON(scope, ec2.TagsLastAppliedAnnotation)
	if err != nil {
		return err
	}

	// Check if the instance tags were changed. If they were, create a new LaunchTemplate.
	tagsChanged, _, _, _ := ec2.TagsChanged(annotation, scope.AdditionalTags()) //nolint:dogsled

	needsUpdate, err := ec2svc.LaunchTemplateNeedsUpdate(scope, scope.GetLaunchTemplate(), launchTemplate)
	if err != nil {
		return err
	}

	amiChanged := *imageID != *launchTemplate.AMI.ID

	// `launchTemplateUserDataSecretKey` can be nil since it comes from a tag on the launch template
	// which may not exist in older launch templates created by older CAPA versions.
	// On change, we trigger instance refresh (rollout of new nodes). Therefore, do not consider it a change if the
	// launch template does not have the respective tag yet, as it could be surprising to users. Instead, ensure the
	// tag is stored on the newly-generated launch template version, without rolling out nodes.
	userDataSecretKeyChanged := launchTemplateUserDataSecretKey != nil && bootstrapDataSecretKey.String() != launchTemplateUserDataSecretKey.String()
	launchTemplateNeedsUserDataSecretKeyTag := launchTemplateUserDataSecretKey == nil

	if needsUpdate || tagsChanged || amiChanged || userDataSecretKeyChanged {
		canUpdate, err := canUpdateLaunchTemplate()
		if err != nil {
			return err
		}
		if !canUpdate {
			conditions.MarkFalse(scope.GetSetter(), expinfrav1.PreLaunchTemplateUpdateCheckCondition, expinfrav1.PreLaunchTemplateUpdateCheckFailedReason, clusterv1.ConditionSeverityWarning, "")
			return errors.New("cannot update the launch template, prerequisite not met")
		}
	}

	userDataHashChanged := launchTemplateUserDataHash != bootstrapDataHash

	// Create a new launch template version if there's a difference in configuration, tags,
	// userdata, OR we've discovered a new AMI ID.
	if needsUpdate || tagsChanged || amiChanged || userDataHashChanged || userDataSecretKeyChanged || launchTemplateNeedsUserDataSecretKeyTag {
		scope.Info("creating new version for launch template", "existing", launchTemplate, "incoming", scope.GetLaunchTemplate(), "needsUpdate", needsUpdate, "tagsChanged", tagsChanged, "amiChanged", amiChanged, "userDataHashChanged", userDataHashChanged, "userDataSecretKeyChanged", userDataSecretKeyChanged)
		// There is a limit to the number of Launch Template Versions.
		// We ensure that the number of versions does not grow without bound by following a simple rule: Before we create a new version, we delete one old version, if there is at least one old version that is not in use.
		if err := ec2svc.PruneLaunchTemplateVersions(scope.GetLaunchTemplateIDStatus()); err != nil {
			return err
		}
		if err := ec2svc.CreateLaunchTemplateVersion(scope.GetLaunchTemplateIDStatus(), scope, imageID, *bootstrapDataSecretKey, bootstrapData); err != nil {
			return err
		}
		version, err := ec2svc.GetLaunchTemplateLatestVersion(scope.GetLaunchTemplateIDStatus())
		if err != nil {
			return err
		}

		scope.SetLaunchTemplateLatestVersionStatus(version)
		if err := scope.PatchObject(); err != nil {
			return err
		}
	}

	if needsUpdate || tagsChanged || amiChanged || userDataSecretKeyChanged {
		if err := runPostLaunchTemplateUpdateOperation(); err != nil {
			conditions.MarkFalse(scope.GetSetter(), expinfrav1.PostLaunchTemplateUpdateOperationCondition, expinfrav1.PostLaunchTemplateUpdateOperationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
		conditions.MarkTrue(scope.GetSetter(), expinfrav1.PostLaunchTemplateUpdateOperationCondition)
	}

	return nil
}

// ReconcileTags reconciles the tags for the AWSMachinePool instances.
func (s *Service) ReconcileTags(scope scope.LaunchTemplateScope, resourceServicesToUpdate []scope.ResourceServiceToUpdate) error {
	additionalTags := scope.AdditionalTags()

	_, err := s.ensureTags(scope, resourceServicesToUpdate, additionalTags)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ensureTags(scope scope.LaunchTemplateScope, resourceServicesToUpdate []scope.ResourceServiceToUpdate, additionalTags map[string]string) (bool, error) {
	annotation, err := ec2.MachinePoolAnnotationJSON(scope, ec2.TagsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	// Check if the instance tags were changed. If they were, update them.
	// It would be possible here to only send new/updated tags, but for the
	// moment we send everything, even if only a single tag was created or
	// upated.
	changed, created, deleted, newAnnotation := ec2.TagsChanged(annotation, additionalTags)
	if changed {
		for _, resourceServiceToUpdate := range resourceServicesToUpdate {
			err := resourceServiceToUpdate.ResourceService.UpdateResourceTags(resourceServiceToUpdate.ResourceID, created, deleted)
			if err != nil {
				return false, err
			}
		}

		// We also need to update the annotation if anything changed.
		err = ec2.UpdateMachinePoolAnnotationJSON(scope, ec2.TagsLastAppliedAnnotation, newAnnotation)
		if err != nil {
			return false, err
		}
	}

	return changed, nil
}

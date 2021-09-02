/*
Copyright 2018 The Kubernetes Authors.

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

package securitygroup

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	errlist "k8s.io/apimachinery/pkg/util/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const (
	// IPProtocolTCP is how EC2 represents the TCP protocol in ingress rules
	IPProtocolTCP = "tcp"

	// IPProtocolUDP is how EC2 represents the UDP protocol in ingress rules
	IPProtocolUDP = "udp"

	// IPProtocolICMP is how EC2 represents the ICMP protocol in ingress rules
	IPProtocolICMP = "icmp"

	// IPProtocolICMPv6 is how EC2 represents the ICMPv6 protocol in ingress rules
	IPProtocolICMPv6 = "58"
)

var (
	defaultRoles = []infrav1.SecurityGroupRole{
		infrav1.SecurityGroupBastion,
		infrav1.SecurityGroupAPIServerLB,
		infrav1.SecurityGroupLB,
		infrav1.SecurityGroupControlPlane,
		infrav1.SecurityGroupNode,
	}
)

func (s *Service) ReconcileSecurityGroups() error {
	s.scope.V(2).Info("Reconciling security groups")

	if s.scope.Network().SecurityGroups == nil {
		s.scope.Network().SecurityGroups = make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
	}

	sgs, err := s.describeSecurityGroupsByName()
	if err != nil {
		return err
	}

	// First iteration makes sure that the security group are valid and fully created.
	for i := range s.roles {
		role := s.roles[i]
		sg := s.getDefaultSecurityGroup(role)
		existing, ok := sgs[*sg.GroupName]

		if !ok {
			if !s.scope.Bastion().Enabled && role == "bastion" {
				break
			} else if err := s.createSecurityGroup(role, sg); err != nil {
				return err
			}

			s.scope.SecurityGroups()[role] = infrav1.SecurityGroup{
				ID:   *sg.GroupId,
				Name: *sg.GroupName,
			}
			s.scope.V(2).Info("Created security group for role", "role", role, "security-group", s.scope.SecurityGroups()[role])
			continue
		}

		// TODO(vincepri): validate / update security group if necessary.
		s.scope.SecurityGroups()[role] = existing

		if s.isEKSOwned(existing) {
			s.scope.V(2).Info("Security group is EKS owned", "role", role, "security-group", s.scope.SecurityGroups()[role])
			continue
		}

		// Make sure tags are up to date.
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			buildParams := s.getSecurityGroupTagParams(existing.Name, existing.ID, role)
			tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
			if err := tagsBuilder.Ensure(existing.Tags); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.GroupNotFound); err != nil {
			return errors.Wrapf(err, "failed to ensure tags on security group %q", existing.ID)
		}
	}

	// Second iteration creates or updates all permissions on the security group to match
	// the specified ingress rules.
	for i := range s.scope.SecurityGroups() {
		sg := s.scope.SecurityGroups()[i]
		s.scope.V(2).Info("second pass security group reconciliation", "group-id", sg.ID, "name", sg.Name, "role", i)
		if sg.Tags.HasAWSCloudProviderOwned(s.scope.Name()) || s.isEKSOwned(sg) {
			// skip rule reconciliation, as we expect the in-cluster cloud integration to manage them
			continue
		}
		current := sg.IngressRules

		want, err := s.getSecurityGroupIngressRules(i)
		if err != nil {
			return err
		}

		toRevoke := current.Difference(want)
		if len(toRevoke) > 0 {
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				if err := s.revokeSecurityGroupIngressRules(sg.ID, toRevoke); err != nil {
					return false, err
				}
				return true, nil
			}, awserrors.GroupNotFound); err != nil {
				return errors.Wrapf(err, "failed to revoke security group ingress rules for %q", sg.ID)
			}

			s.scope.V(2).Info("Revoked ingress rules from security group", "revoked-ingress-rules", toRevoke, "security-group-id", sg.ID)
		}

		toAuthorize := want.Difference(current)
		if len(toAuthorize) > 0 {
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				if err := s.authorizeSecurityGroupIngressRules(sg.ID, toAuthorize); err != nil {
					return false, err
				}
				return true, nil
			}, awserrors.GroupNotFound); err != nil {
				return err
			}

			s.scope.V(2).Info("Authorized ingress rules in security group", "authorized-ingress-rules", toAuthorize, "security-group-id", sg.ID)
		}
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.ClusterSecurityGroupsReadyCondition)
	return nil
}

func (s *Service) DeleteSecurityGroups() error {
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.ClusterSecurityGroupsReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	for _, sg := range s.scope.SecurityGroups() {
		current := sg.IngressRules

		if s.isEKSOwned(sg) {
			continue
		}

		if err := s.revokeAllSecurityGroupIngressRules(sg.ID); awserrors.IsIgnorableSecurityGroupError(err) != nil {
			conditions.MarkFalse(s.scope.InfraCluster(), infrav1.ClusterSecurityGroupsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
			return err
		}

		s.scope.V(2).Info("Revoked ingress rules from security group", "revoked-ingress-rules", current, "security-group-id", sg.ID)
	}

	for i := range s.scope.SecurityGroups() {
		sg := s.scope.SecurityGroups()[i]

		if s.isEKSOwned(sg) {
			continue
		}

		if err := s.deleteSecurityGroup(&sg, "managed"); err != nil {
			conditions.MarkFalse(s.scope.InfraCluster(), infrav1.ClusterSecurityGroupsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
			return err
		}
	}

	clusterGroups, err := s.describeClusterOwnedSecurityGroups()
	if err != nil {
		return err
	}

	for i := range clusterGroups {
		sg := clusterGroups[i]
		if deleteErr := s.deleteSecurityGroup(&sg, "cluster managed"); deleteErr != nil {
			err = errlist.NewAggregate([]error{err, deleteErr})
		}
	}

	if err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.ClusterSecurityGroupsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.ClusterSecurityGroupsReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	return nil
}

func (s *Service) deleteSecurityGroup(sg *infrav1.SecurityGroup, typ string) error {
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sg.ID),
	}

	if _, err := s.EC2Client.DeleteSecurityGroup(input); awserrors.IsIgnorableSecurityGroupError(err) != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedDeleteSecurityGroup", "Failed to delete %s SecurityGroup %q: %v", typ, sg.ID, err)
		return errors.Wrapf(err, "failed to delete security group %q", sg.ID)
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteSecurityGroup", "Deleted %s SecurityGroup %q", typ, sg.ID)
	s.scope.V(2).Info("Deleted security group", "security-group-id", sg.ID, "kind", typ)

	return nil
}

func (s *Service) describeClusterOwnedSecurityGroups() ([]infrav1.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(s.scope.VPC().ID),
			filter.EC2.ProviderOwned(s.scope.Name()),
		},
	}

	groups := []infrav1.SecurityGroup{}

	err := s.EC2Client.DescribeSecurityGroupsPages(input, func(out *ec2.DescribeSecurityGroupsOutput, last bool) bool {
		for _, group := range out.SecurityGroups {
			if group != nil {
				groups = append(groups, makeInfraSecurityGroup(group))
			}
		}
		return true
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe cluster-owned security groups in vpc %q", s.scope.VPC().ID)
	}
	return groups, nil
}

func (s *Service) describeSecurityGroupsByName() (map[string]infrav1.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(s.scope.VPC().ID),
			filter.EC2.Cluster(s.scope.Name()),
		},
	}

	out, err := s.EC2Client.DescribeSecurityGroups(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe security groups in vpc %q", s.scope.VPC().ID)
	}

	res := make(map[string]infrav1.SecurityGroup, len(out.SecurityGroups))
	for _, ec2sg := range out.SecurityGroups {
		sg := makeInfraSecurityGroup(ec2sg)

		for _, ec2rule := range ec2sg.IpPermissions {
			sg.IngressRules = append(sg.IngressRules, ingressRuleFromSDKType(ec2rule))
		}

		res[sg.Name] = sg
	}

	return res, nil
}

func makeInfraSecurityGroup(ec2sg *ec2.SecurityGroup) infrav1.SecurityGroup {
	return infrav1.SecurityGroup{
		ID:   *ec2sg.GroupId,
		Name: *ec2sg.GroupName,
		Tags: converters.TagsToMap(ec2sg.Tags),
	}
}

func (s *Service) createSecurityGroup(role infrav1.SecurityGroupRole, input *ec2.SecurityGroup) error {
	sgTags := s.getSecurityGroupTagParams(aws.StringValue(input.GroupName), services.TemporaryResourceID, role)
	out, err := s.EC2Client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		VpcId:       input.VpcId,
		GroupName:   input.GroupName,
		Description: aws.String(fmt.Sprintf("Kubernetes cluster %s: %s", s.scope.Name(), role)),
		TagSpecifications: []*ec2.TagSpecification{
			tags.BuildParamsToTagSpecification(ec2.ResourceTypeSecurityGroup, sgTags),
		},
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateSecurityGroup", "Failed to create managed SecurityGroup for Role %q: %v", role, err)
		return errors.Wrapf(err, "failed to create security group %q in vpc %q", role, aws.StringValue(input.VpcId))
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateSecurityGroup", "Created managed SecurityGroup %q for Role %q", aws.StringValue(out.GroupId), role)

	// Set the group id.
	input.GroupId = out.GroupId

	return nil
}

func (s *Service) authorizeSecurityGroupIngressRules(id string, rules infrav1.IngressRules) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{GroupId: aws.String(id)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.EC2Client.AuthorizeSecurityGroupIngress(input); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAuthorizeSecurityGroupIngressRules", "Failed to authorize security group ingress rules %v for SecurityGroup %q: %v", rules, id, err)
		return errors.Wrapf(err, "failed to authorize security group %q ingress rules: %v", id, rules)
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulAuthorizeSecurityGroupIngressRules", "Authorized security group ingress rules %v for SecurityGroup %q", rules, id)
	return nil
}

func (s *Service) revokeSecurityGroupIngressRules(id string, rules infrav1.IngressRules) error {
	input := &ec2.RevokeSecurityGroupIngressInput{GroupId: aws.String(id)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.EC2Client.RevokeSecurityGroupIngress(input); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedRevokeSecurityGroupIngressRules", "Failed to revoke security group ingress rules %v for SecurityGroup %q: %v", rules, id, err)
		return errors.Wrapf(err, "failed to revoke security group %q ingress rules: %v", id, rules)
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulRevokeSecurityGroupIngressRules", "Revoked security group ingress rules %v for SecurityGroup %q", rules, id)
	return nil
}

func (s *Service) revokeAllSecurityGroupIngressRules(id string) error {
	describeInput := &ec2.DescribeSecurityGroupsInput{GroupIds: []*string{aws.String(id)}}

	securityGroups, err := s.EC2Client.DescribeSecurityGroups(describeInput)
	if err != nil {
		return errors.Wrapf(err, "failed to query security group %q", id)
	}

	for _, sg := range securityGroups.SecurityGroups {
		if len(sg.IpPermissions) > 0 {
			revokeInput := &ec2.RevokeSecurityGroupIngressInput{
				GroupId:       aws.String(id),
				IpPermissions: sg.IpPermissions,
			}
			if _, err := s.EC2Client.RevokeSecurityGroupIngress(revokeInput); err != nil {
				record.Warnf(s.scope.InfraCluster(), "FailedRevokeSecurityGroupIngressRules", "Failed to revoke all security group ingress rules for SecurityGroup %q: %v", *sg.GroupId, err)
				return errors.Wrapf(err, "failed to revoke security group %q ingress rules", id)
			}
			record.Eventf(s.scope.InfraCluster(), "SuccessfulRevokeSecurityGroupIngressRules", "Revoked all security group ingress rules for SecurityGroup %q", *sg.GroupId)
		}
	}

	return nil
}

func (s *Service) defaultSSHIngressRule(sourceSecurityGroupID string) *infrav1.IngressRule {
	return &infrav1.IngressRule{
		Description:            "SSH",
		Protocol:               infrav1.SecurityGroupProtocolTCP,
		FromPort:               22,
		ToPort:                 22,
		SourceSecurityGroupIDs: []string{sourceSecurityGroupID},
	}
}

func (s *Service) getSecurityGroupIngressRules(role infrav1.SecurityGroupRole) (infrav1.IngressRules, error) {
	// Set source of CNI ingress rules to be control plane and node security groups
	s.scope.V(2).Info("getting security group ingress rules", "role", role)

	cniRules := make(infrav1.IngressRules, len(s.scope.CNIIngressRules()))
	for i, r := range s.scope.CNIIngressRules() {
		cniRules[i] = &infrav1.IngressRule{
			Description: r.Description,
			Protocol:    r.Protocol,
			FromPort:    r.FromPort,
			ToPort:      r.ToPort,
			SourceSecurityGroupIDs: []string{
				s.scope.SecurityGroups()[infrav1.SecurityGroupControlPlane].ID,
				s.scope.SecurityGroups()[infrav1.SecurityGroupNode].ID,
			},
		}
	}

	switch role {
	case infrav1.SecurityGroupBastion:
		return infrav1.IngressRules{
			{
				Description: "SSH",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    22,
				ToPort:      22,
				CidrBlocks:  s.scope.Bastion().AllowedCIDRBlocks,
			},
		}, nil
	case infrav1.SecurityGroupControlPlane:
		rules := infrav1.IngressRules{
			s.defaultSSHIngressRule(s.scope.SecurityGroups()[infrav1.SecurityGroupBastion].ID),
			{
				Description: "Kubernetes API",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    6443,
				ToPort:      6443,
				SourceSecurityGroupIDs: []string{
					s.scope.SecurityGroups()[infrav1.SecurityGroupAPIServerLB].ID,
					s.scope.SecurityGroups()[infrav1.SecurityGroupControlPlane].ID,
					s.scope.SecurityGroups()[infrav1.SecurityGroupNode].ID,
				},
			},
			{
				Description:            "etcd",
				Protocol:               infrav1.SecurityGroupProtocolTCP,
				FromPort:               2379,
				ToPort:                 2379,
				SourceSecurityGroupIDs: []string{s.scope.SecurityGroups()[infrav1.SecurityGroupControlPlane].ID},
			},
			{
				Description:            "etcd peer",
				Protocol:               infrav1.SecurityGroupProtocolTCP,
				FromPort:               2380,
				ToPort:                 2380,
				SourceSecurityGroupIDs: []string{s.scope.SecurityGroups()[infrav1.SecurityGroupControlPlane].ID},
			},
		}
		return append(cniRules, rules...), nil

	case infrav1.SecurityGroupNode:
		rules := infrav1.IngressRules{
			s.defaultSSHIngressRule(s.scope.SecurityGroups()[infrav1.SecurityGroupBastion].ID),
			{
				Description: "Node Port Services",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    30000,
				ToPort:      32767,
				CidrBlocks:  []string{services.AnyIPv4CidrBlock},
			},
			{
				Description: "Kubelet API",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    10250,
				ToPort:      10250,
				SourceSecurityGroupIDs: []string{
					s.scope.SecurityGroups()[infrav1.SecurityGroupControlPlane].ID,
					// This is needed to support metrics-server deployments
					s.scope.SecurityGroups()[infrav1.SecurityGroupNode].ID,
				},
			},
		}
		return append(cniRules, rules...), nil
	case infrav1.SecurityGroupEKSNodeAdditional:
		if !s.scope.Bastion().Enabled {
			return infrav1.IngressRules{
				s.defaultSSHIngressRule(s.scope.SecurityGroups()[infrav1.SecurityGroupBastion].ID),
			}, nil
		}
	case infrav1.SecurityGroupAPIServerLB:
		return infrav1.IngressRules{
			{
				Description: "Kubernetes API",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    int64(s.scope.APIServerPort()),
				ToPort:      int64(s.scope.APIServerPort()),
				CidrBlocks:  []string{services.AnyIPv4CidrBlock},
			},
		}, nil
	case infrav1.SecurityGroupLB:
		// We hand this group off to the in-cluster cloud provider, so these rules aren't used
		return infrav1.IngressRules{}, nil
	}

	return nil, errors.Errorf("Cannot determine ingress rules for unknown security group role %q", role)
}

func (s *Service) getSecurityGroupName(clusterName string, role infrav1.SecurityGroupRole) string {
	groupPrefix := clusterName
	if strings.HasPrefix(clusterName, "sg-") {
		groupPrefix = "@" + clusterName
	}
	return fmt.Sprintf("%s-%v", groupPrefix, role)
}

func (s *Service) getDefaultSecurityGroup(role infrav1.SecurityGroupRole) *ec2.SecurityGroup {
	name := s.getSecurityGroupName(s.scope.Name(), role)

	return &ec2.SecurityGroup{
		GroupName: aws.String(name),
		VpcId:     aws.String(s.scope.VPC().ID),
		Tags:      converters.MapToTags(infrav1.Build(s.getSecurityGroupTagParams(name, "", role))),
	}
}

func (s *Service) getSecurityGroupTagParams(name, id string, role infrav1.SecurityGroupRole) infrav1.BuildParams {
	additional := s.scope.AdditionalTags()
	if role == infrav1.SecurityGroupLB {
		additional[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)
	}
	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		ResourceID:  id,
		Role:        aws.String(string(role)),
		Additional:  additional,
	}
}

func (s *Service) isEKSOwned(sg infrav1.SecurityGroup) bool {
	_, ok := sg.Tags["aws:eks:cluster-name"]
	return ok
}

func ingressRuleToSDKType(i *infrav1.IngressRule) (res *ec2.IpPermission) {
	// AWS seems to ignore the From/To port when set on protocols where it doesn't apply, but
	// we avoid serializing it out for clarity's sake.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
	switch i.Protocol {
	case infrav1.SecurityGroupProtocolTCP,
		infrav1.SecurityGroupProtocolUDP,
		infrav1.SecurityGroupProtocolICMP,
		infrav1.SecurityGroupProtocolICMPv6:
		res = &ec2.IpPermission{
			IpProtocol: aws.String(string(i.Protocol)),
			FromPort:   aws.Int64(i.FromPort),
			ToPort:     aws.Int64(i.ToPort),
		}
	case infrav1.SecurityGroupProtocolAll, infrav1.SecurityGroupProtocolIPinIP:
		res = &ec2.IpPermission{
			IpProtocol: aws.String(string(i.Protocol)),
		}
	}

	for _, cidr := range i.CidrBlocks {
		ipRange := &ec2.IpRange{
			CidrIp: aws.String(cidr),
		}

		if i.Description != "" {
			ipRange.Description = aws.String(i.Description)
		}

		res.IpRanges = append(res.IpRanges, ipRange)
	}

	for _, groupID := range i.SourceSecurityGroupIDs {
		userIDGroupPair := &ec2.UserIdGroupPair{
			GroupId: aws.String(groupID),
		}

		if i.Description != "" {
			userIDGroupPair.Description = aws.String(i.Description)
		}

		res.UserIdGroupPairs = append(res.UserIdGroupPairs, userIDGroupPair)
	}

	return res
}

func ingressRuleFromSDKType(v *ec2.IpPermission) (res *infrav1.IngressRule) {
	// Ports are only well-defined for TCP and UDP protocols, but EC2 overloads the port range
	// in the case of ICMP(v6) traffic to indicate which codes are allowed. For all other protocols,
	// including the custom "-1" All Traffic protcol, FromPort and ToPort are omitted from the response.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
	switch *v.IpProtocol {
	case IPProtocolTCP,
		IPProtocolUDP,
		IPProtocolICMP,
		IPProtocolICMPv6:
		res = &infrav1.IngressRule{
			Protocol: infrav1.SecurityGroupProtocol(*v.IpProtocol),
			FromPort: *v.FromPort,
			ToPort:   *v.ToPort,
		}
	default:
		res = &infrav1.IngressRule{
			Protocol: infrav1.SecurityGroupProtocol(*v.IpProtocol),
		}
	}

	for _, ec2range := range v.IpRanges {
		if ec2range.Description != nil && *ec2range.Description != "" {
			res.Description = *ec2range.Description
		}

		res.CidrBlocks = append(res.CidrBlocks, *ec2range.CidrIp)
	}

	for _, pair := range v.UserIdGroupPairs {
		if pair.GroupId == nil {
			continue
		}

		if pair.Description != nil && *pair.Description != "" {
			res.Description = *pair.Description
		}

		res.SourceSecurityGroupIDs = append(res.SourceSecurityGroupIDs, *pair.GroupId)
	}

	return res
}

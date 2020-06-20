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

package ec2

import (
	"fmt"

	errlist "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
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

// Declare all security group roles that the reconcile loop takes care of.
var roles = []infrav1.SecurityGroupRole{
	infrav1.SecurityGroupBastion,
	infrav1.SecurityGroupControlPlane,
	infrav1.SecurityGroupAPIServerLB,
	infrav1.SecurityGroupNode,
	infrav1.SecurityGroupLB,
}

func (s *Service) reconcileSecurityGroups() error {
	s.scope.V(2).Info("Reconciling security groups")

	if s.scope.Network().SecurityGroups == nil {
		s.scope.Network().SecurityGroups = make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
	}

	var sgs map[string]infrav1.SecurityGroup
	var err error

	// Security group overrides are mapped by Role rather than their security group name
	// They are copied into the main 'sgs' list by their group name later
	var securityGroupOverrides map[infrav1.SecurityGroupRole]*ec2.SecurityGroup
	securityGroupOverrides, err = s.describeSecurityGroupOverridesByID()
	if err != nil {
		return err
	}

	sgs, err = s.describeSecurityGroupsByName()
	if err != nil {
		return err
	}

	// Add security group overrides to known security group map
	for _, securityGroupOverride := range securityGroupOverrides {
		sg := s.ec2SecurityGroupToSecurityGroup(securityGroupOverride)
		sgs[sg.Name] = sg
	}

	// First iteration makes sure that the security group are valid and fully created.
	for i := range roles {
		role := roles[i]
		sg := s.getDefaultSecurityGroup(role)

		sgOverride, ok := securityGroupOverrides[role]
		if ok {
			s.scope.V(2).Info("Using security group override", "role", role, "security group", sgOverride.GroupName)
			sg = sgOverride
		}

		existing, ok := sgs[*sg.GroupName]

		if !ok {
			if err := s.createSecurityGroup(role, sg); err != nil {
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

		// Make sure tags are up to date.
		if !s.securityGroupIsOverridden(*sg.GroupId) {
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				if err := tags.Ensure(existing.Tags, &tags.ApplyParams{
					EC2Client:   s.scope.EC2,
					BuildParams: s.getSecurityGroupTagParams(existing.Name, existing.ID, role),
				}); err != nil {
					return false, err
				}
				return true, nil
			}, awserrors.GroupNotFound); err != nil {
				return errors.Wrapf(err, "failed to ensure tags on security group %q", existing.ID)
			}
		}
	}

	// Second iteration creates or updates all permissions on the security group to match
	// the specified ingress rules.
	for role := range s.scope.SecurityGroups() {
		sg := s.scope.SecurityGroups()[role]
		if s.securityGroupIsOverridden(sg.ID) {
			// skip rule/tag reconciliation on security groups that are overridden, assuming they're managed by another process
			continue
		}

		if sg.Tags.HasAWSCloudProviderOwned(s.scope.Name()) {
			// skip rule reconciliation, as we expect the in-cluster cloud integration to manage them
			continue
		}

		current := sg.IngressRules

		want, err := s.getSecurityGroupIngressRules(role)
		if err != nil {
			return err
		}

		toRevoke := current.Difference(want)
		toAuthorize := want.Difference(current)

		if securityGroupOverrides[role] != nil {
			s.scope.V(4).Info("Skipping security group rules reconciliation for overidden security group", "revocations", toRevoke, "authorizations", toAuthorize)
			continue
		}

		if len(toRevoke) > 0 {
			securityGroupID := sg.ID
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				if err := s.revokeSecurityGroupIngressRules(securityGroupID, toRevoke); err != nil {
					return false, err
				}
				return true, nil
			}, awserrors.GroupNotFound); err != nil {
				return errors.Wrapf(err, "failed to revoke security group ingress rules for %q", sg.ID)
			}

			s.scope.V(2).Info("Revoked ingress rules from security group", "revoked-ingress-rules", toRevoke, "security-group-id", sg.ID)
		}

		if len(toAuthorize) > 0 {
			securityGroupID := sg.ID
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				if err := s.authorizeSecurityGroupIngressRules(securityGroupID, toAuthorize); err != nil {
					return false, err
				}
				return true, nil
			}, awserrors.GroupNotFound); err != nil {
				return err
			}

			s.scope.V(2).Info("Authorized ingress rules in security group", "authorized-ingress-rules", toAuthorize, "security-group-id", sg.ID)
		}
	}

	return nil
}

func (s *Service) securityGroupIsOverridden(securityGroupID string) bool {
	for _, overrideID := range s.scope.SecurityGroupOverrides() {
		if overrideID == securityGroupID {
			return true
		}
	}
	return false
}

func (s *Service) describeSecurityGroupOverridesByID() (map[infrav1.SecurityGroupRole]*ec2.SecurityGroup, error) {
	securityGroupIds := map[infrav1.SecurityGroupRole]*string{}
	input := &ec2.DescribeSecurityGroupsInput{}

	overrides := s.scope.SecurityGroupOverrides()

	// return if no security group overrides have been provided
	if len(overrides) == 0 {
		return nil, nil
	}

	if len(overrides) > 0 {
		for _, role := range roles {
			securityGroupID, ok := s.scope.SecurityGroupOverrides()[role]
			if !ok {
				return nil, fmt.Errorf("security group overrides have been provided for some but not all roles - missing security group for role %s", role)
			}
			securityGroupIds[role] = aws.String(securityGroupID)
			input.GroupIds = append(input.GroupIds, aws.String(securityGroupID))
		}
	}

	out, err := s.scope.EC2.DescribeSecurityGroups(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe security groups in vpc %q", s.scope.VPC().ID)
	}

	res := make(map[infrav1.SecurityGroupRole]*ec2.SecurityGroup, len(out.SecurityGroups))
	for _, role := range roles {
		for _, ec2sg := range out.SecurityGroups {
			if *ec2sg.GroupId == *securityGroupIds[role] {
				s.scope.V(2).Info("found security group override", "role", role, "security group", *ec2sg.GroupName)

				res[role] = ec2sg
				break
			}
		}
	}

	return res, nil
}

func (s *Service) deleteSecurityGroups() error {
	for i := range s.scope.SecurityGroups() {
		sg := s.scope.SecurityGroups()[i]
		current := sg.IngressRules

		// do not attempt to delete security groups that are overrides
		if s.securityGroupIsOverridden(sg.ID) {
			continue
		}

		if err := s.revokeAllSecurityGroupIngressRules(sg.ID); awserrors.IsIgnorableSecurityGroupError(err) != nil {
			return err
		}

		s.scope.V(2).Info("Revoked ingress rules from security group", "revoked-ingress-rules", current, "security-group-id", sg.ID)
	}

	for i := range s.scope.SecurityGroups() {
		sg := s.scope.SecurityGroups()[i]
		// do not attempt to delete security groups that are overrides
		if s.securityGroupIsOverridden(sg.ID) {
			continue
		}

		if err := s.deleteSecurityGroup(&sg, "managed"); err != nil {
			return err
		}
	}

	clusterGroups, err := s.describeClusterOwnedSecurityGroups()
	if err != nil {
		return err
	}

	errs := []error{}
	for i := range clusterGroups {
		sg := clusterGroups[i]
		// do not attempt to delete security groups that are overrides
		if s.securityGroupIsOverridden(sg.ID) {
			continue
		}

		if err := s.deleteSecurityGroup(&sg, "cluster managed"); err != nil {
			errs = append(errs, err)
		}

		if len(errs) != 0 {
			return errlist.NewAggregate(errs)
		}
	}

	return nil
}

func (s *Service) deleteSecurityGroup(sg *infrav1.SecurityGroup, typ string) error {
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sg.ID),
	}

	if _, err := s.scope.EC2.DeleteSecurityGroup(input); awserrors.IsIgnorableSecurityGroupError(err) != nil {
		record.Warnf(s.scope.AWSCluster, "FailedDeleteSecurityGroup", "Failed to delete %s SecurityGroup %q: %v", typ, sg.ID, err)
		return errors.Wrapf(err, "failed to delete security group %q", sg.ID)
	}

	record.Eventf(s.scope.AWSCluster, "SuccessfulDeleteSecurityGroup", "Deleted %s SecurityGroup %q", typ, sg.ID)
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

	err := s.scope.EC2.DescribeSecurityGroupsPages(input, func(out *ec2.DescribeSecurityGroupsOutput, last bool) bool {
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

func (s *Service) ec2SecurityGroupToSecurityGroup(ec2SecurityGroup *ec2.SecurityGroup) infrav1.SecurityGroup {
	sg := infrav1.SecurityGroup{
		ID:   *ec2SecurityGroup.GroupId,
		Name: *ec2SecurityGroup.GroupName,
		Tags: converters.TagsToMap(ec2SecurityGroup.Tags),
	}

	for _, ec2rule := range ec2SecurityGroup.IpPermissions {
		sg.IngressRules = append(sg.IngressRules, ingressRuleFromSDKType(ec2rule))
	}
	return sg
}

func (s *Service) describeSecurityGroupsByName() (map[string]infrav1.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(s.scope.VPC().ID),
			filter.EC2.Cluster(s.scope.Name()),
		},
	}

	out, err := s.scope.EC2.DescribeSecurityGroups(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe security groups in vpc %q", s.scope.VPC().ID)
	}

	res := make(map[string]infrav1.SecurityGroup, len(out.SecurityGroups))
	for _, ec2sg := range out.SecurityGroups {
		sg := s.ec2SecurityGroupToSecurityGroup(ec2sg)
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
	out, err := s.scope.EC2.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		VpcId:       input.VpcId,
		GroupName:   input.GroupName,
		Description: aws.String(fmt.Sprintf("Kubernetes cluster %s: %s", s.scope.Name(), role)),
	})

	if err != nil {
		record.Warnf(s.scope.AWSCluster, "FailedCreateSecurityGroup", "Failed to create managed SecurityGroup for Role %q: %v", role, err)
		return errors.Wrapf(err, "failed to create security group %q in vpc %q", role, aws.StringValue(input.VpcId))
	}

	record.Eventf(s.scope.AWSCluster, "SuccessfulCreateSecurityGroup", "Created managed SecurityGroup %q for Role %q", aws.StringValue(out.GroupId), role)

	// Set the group id.
	input.GroupId = out.GroupId

	// Tag the security group.
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.scope.EC2.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{out.GroupId},
			Tags:      input.Tags,
		}); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.GroupNotFound); err != nil {
		record.Warnf(s.scope.AWSCluster, "FailedTagSecurityGroup", "Failed to tag managed SecurityGroup %q: %v", aws.StringValue(out.GroupId), err)
		return errors.Wrapf(err, "failed to tag security group %q in vpc %q", aws.StringValue(out.GroupId), aws.StringValue(input.VpcId))
	}

	record.Eventf(s.scope.AWSCluster, "SuccessfulTagSecurityGroup", "Tagged managed SecurityGroup %q", aws.StringValue(out.GroupId))
	return nil
}

func (s *Service) authorizeSecurityGroupIngressRules(id string, rules infrav1.IngressRules) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{GroupId: aws.String(id)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.scope.EC2.AuthorizeSecurityGroupIngress(input); err != nil {
		record.Warnf(s.scope.AWSCluster, "FailedAuthorizeSecurityGroupIngressRules", "Failed to authorize security group ingress rules %v for SecurityGroup %q: %v", rules, id, err)
		return errors.Wrapf(err, "failed to authorize security group %q ingress rules: %v", id, rules)
	}

	record.Eventf(s.scope.AWSCluster, "SuccessfulAuthorizeSecurityGroupIngressRules", "Authorized security group ingress rules %v for SecurityGroup %q", rules, id)
	return nil
}

func (s *Service) revokeSecurityGroupIngressRules(id string, rules infrav1.IngressRules) error {
	input := &ec2.RevokeSecurityGroupIngressInput{GroupId: aws.String(id)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.scope.EC2.RevokeSecurityGroupIngress(input); err != nil {
		record.Warnf(s.scope.AWSCluster, "FailedRevokeSecurityGroupIngressRules", "Failed to revoke security group ingress rules %v for SecurityGroup %q: %v", rules, id, err)
		return errors.Wrapf(err, "failed to revoke security group %q ingress rules: %v", id, rules)
	}

	record.Eventf(s.scope.AWSCluster, "SuccessfulRevokeSecurityGroupIngressRules", "Revoked security group ingress rules %v for SecurityGroup %q", rules, id)
	return nil
}

func (s *Service) revokeAllSecurityGroupIngressRules(id string) error {
	describeInput := &ec2.DescribeSecurityGroupsInput{GroupIds: []*string{aws.String(id)}}

	securityGroups, err := s.scope.EC2.DescribeSecurityGroups(describeInput)
	if err != nil {
		return errors.Wrapf(err, "failed to query security group %q", id)
	}

	for _, sg := range securityGroups.SecurityGroups {
		if len(sg.IpPermissions) > 0 {
			revokeInput := &ec2.RevokeSecurityGroupIngressInput{
				GroupId:       aws.String(id),
				IpPermissions: sg.IpPermissions,
			}
			if _, err := s.scope.EC2.RevokeSecurityGroupIngress(revokeInput); err != nil {
				record.Warnf(s.scope.AWSCluster, "FailedRevokeSecurityGroupIngressRules", "Failed to revoke all security group ingress rules for SecurityGroup %q: %v", *sg.GroupId, err)
				return errors.Wrapf(err, "failed to revoke security group %q ingress rules", id)
			}
			record.Eventf(s.scope.AWSCluster, "SuccessfulRevokeSecurityGroupIngressRules", "Revoked all security group ingress rules for SecurityGroup %q", *sg.GroupId)
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
				CidrBlocks:  []string{anyIPv4CidrBlock},
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
				CidrBlocks:  []string{anyIPv4CidrBlock},
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
	case infrav1.SecurityGroupAPIServerLB:
		return infrav1.IngressRules{
			{
				Description: "Kubernetes API",
				Protocol:    infrav1.SecurityGroupProtocolTCP,
				FromPort:    int64(s.scope.APIServerPort()),
				ToPort:      int64(s.scope.APIServerPort()),
				CidrBlocks:  []string{anyIPv4CidrBlock},
			},
		}, nil
	case infrav1.SecurityGroupLB:
		// We hand this group off to the in-cluster cloud provider, so these rules aren't used
		return infrav1.IngressRules{}, nil
	}

	return nil, errors.Errorf("Cannot determine ingress rules for unknown security group role %q", role)
}

func (s *Service) getSecurityGroupName(clusterName string, role infrav1.SecurityGroupRole) string {
	return fmt.Sprintf("%s-%v", clusterName, role)
}

func (s *Service) getDefaultSecurityGroup(role infrav1.SecurityGroupRole) *ec2.SecurityGroup {
	name := s.getSecurityGroupName(s.scope.Name(), role)

	return &ec2.SecurityGroup{
		GroupName: aws.String(name),
		VpcId:     aws.String(s.scope.VPC().ID),
		Tags:      converters.MapToTags(infrav1.Build(s.getSecurityGroupTagParams(name, "", role))),
	}
}

func (s *Service) getSecurityGroupTagParams(name string, id string, role infrav1.SecurityGroupRole) infrav1.BuildParams {
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
	default:
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

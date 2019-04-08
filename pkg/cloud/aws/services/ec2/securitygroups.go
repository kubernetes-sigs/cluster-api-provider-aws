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

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
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

func (s *Service) reconcileSecurityGroups() error {
	klog.V(2).Infof("Reconciling security groups")

	if s.scope.Network().SecurityGroups == nil {
		s.scope.Network().SecurityGroups = make(map[v1alpha1.SecurityGroupRole]*v1alpha1.SecurityGroup)
	}

	sgs, err := s.describeSecurityGroupsByName()
	if err != nil {
		return err
	}

	// Declare all security group roles that the reconcile loop takes care of.
	roles := []v1alpha1.SecurityGroupRole{
		v1alpha1.SecurityGroupBastion,
		v1alpha1.SecurityGroupControlPlane,
		v1alpha1.SecurityGroupNode,
	}

	// First iteration makes sure that the security group are valid and fully created.
	for _, role := range roles {
		sg := s.getDefaultSecurityGroup(role)
		existing, ok := sgs[*sg.GroupName]

		if !ok {
			if err := s.createSecurityGroup(role, sg); err != nil {
				return err
			}

			s.scope.SecurityGroups()[role] = &v1alpha1.SecurityGroup{
				ID:   *sg.GroupId,
				Name: *sg.GroupName,
			}
			klog.V(2).Infof("Security group for role %q: %v", role, s.scope.SecurityGroups()[role])
			continue
		}

		// TODO(vincepri): validate / update security group if necessary.
		s.scope.SecurityGroups()[role] = existing

		// Make sure tags are up to date.
		err := tags.Ensure(existing.Tags, &tags.ApplyParams{
			EC2Client:   s.scope.EC2,
			BuildParams: s.getSecurityGroupTagParams(existing.Name, role),
		})

		if err != nil {
			return errors.Wrapf(err, "failed to ensure tags on security group %q", existing.ID)
		}
	}

	// Second iteration creates or updates all permissions on the security group to match
	// the specified ingress rules.
	for role, sg := range s.scope.SecurityGroups() {
		current := sg.IngressRules

		want, err := s.getSecurityGroupIngressRules(role)
		if err != nil {
			return err
		}

		toRevoke := current.Difference(want)
		if len(toRevoke) > 0 {
			if err := s.revokeSecurityGroupIngressRules(sg.ID, toRevoke); err != nil {
				return errors.Wrapf(err, "failed to revoke security group ingress rules for %q", sg.ID)
			}

			klog.V(2).Infof("Revoked ingress rules %v from security group %q", toRevoke, sg)
		}

		toAuthorize := want.Difference(current)
		if len(toAuthorize) > 0 {
			if err := s.authorizeSecurityGroupIngressRules(sg.ID, toAuthorize); err != nil {
				return err
			}

			klog.V(2).Infof("Authorized ingress rules %v in security group %q", toAuthorize, sg)
		}
	}

	return nil
}

func (s *Service) deleteSecurityGroups() error {
	for _, sg := range s.scope.SecurityGroups() {
		current := sg.IngressRules

		if err := s.revokeSecurityGroupIngressRules(sg.ID, current); awserrors.IsIgnorableSecurityGroupError(err) != nil {
			return err
		}

		klog.V(2).Infof("Revoked ingress rules %v from security group %q", current, sg.ID)
	}

	for _, sg := range s.scope.SecurityGroups() {
		input := &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(sg.ID),
		}

		if _, err := s.scope.EC2.DeleteSecurityGroup(input); awserrors.IsIgnorableSecurityGroupError(err) != nil {
			return errors.Wrapf(err, "failed to delete security group %q", sg.ID)
		}

		klog.V(2).Infof("Deleted security group security group %q", sg.ID)
	}

	return nil
}

func (s *Service) describeSecurityGroupsByName() (map[string]*v1alpha1.SecurityGroup, error) {
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

	res := make(map[string]*v1alpha1.SecurityGroup, len(out.SecurityGroups))
	for _, ec2sg := range out.SecurityGroups {
		sg := &v1alpha1.SecurityGroup{
			ID:   *ec2sg.GroupId,
			Name: *ec2sg.GroupName,
			Tags: converters.TagsToMap(ec2sg.Tags),
		}

		for _, ec2rule := range ec2sg.IpPermissions {
			sg.IngressRules = append(sg.IngressRules, ingressRuleFromSDKType(ec2rule))
		}

		res[sg.Name] = sg
	}

	return res, nil
}

func (s *Service) createSecurityGroup(role v1alpha1.SecurityGroupRole, input *ec2.SecurityGroup) error {
	out, err := s.scope.EC2.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		VpcId:       input.VpcId,
		GroupName:   input.GroupName,
		Description: aws.String(fmt.Sprintf("Kubernetes cluster %s: %s", s.scope.Name(), role)),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to create security group %q in vpc %q", *input.GroupName, *input.VpcId)
	}

	// Set the group id.
	input.GroupId = out.GroupId

	// Tag the security group.
	if _, err := s.scope.EC2.CreateTags(&ec2.CreateTagsInput{Resources: []*string{out.GroupId}, Tags: input.Tags}); err != nil {
		return errors.Wrapf(err, "failed to tag security group %q in vpc %q", *input.GroupName, *input.VpcId)
	}

	return nil
}

func (s *Service) authorizeSecurityGroupIngressRules(id string, rules v1alpha1.IngressRules) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{GroupId: aws.String(id)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.scope.EC2.AuthorizeSecurityGroupIngress(input); err != nil {
		return errors.Wrapf(err, "failed to authorize security group %q ingress rules: %v", id, rules)
	}

	return nil
}

func (s *Service) revokeSecurityGroupIngressRules(id string, rules v1alpha1.IngressRules) error {
	input := &ec2.RevokeSecurityGroupIngressInput{GroupId: aws.String(id)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.scope.EC2.RevokeSecurityGroupIngress(input); err != nil {
		return errors.Wrapf(err, "failed to revoke security group %q ingress rules: %v", id, rules)
	}

	return nil
}

func (s *Service) defaultSSHIngressRule(sourceSecurityGroupID string) *v1alpha1.IngressRule {
	return &v1alpha1.IngressRule{
		Description:            "SSH",
		Protocol:               v1alpha1.SecurityGroupProtocolTCP,
		FromPort:               22,
		ToPort:                 22,
		SourceSecurityGroupIDs: []string{sourceSecurityGroupID},
	}
}

func (s *Service) getSecurityGroupIngressRules(role v1alpha1.SecurityGroupRole) (v1alpha1.IngressRules, error) {
	switch role {
	case v1alpha1.SecurityGroupBastion:
		return v1alpha1.IngressRules{
			{
				Description: "SSH",
				Protocol:    v1alpha1.SecurityGroupProtocolTCP,
				FromPort:    22,
				ToPort:      22,
				CidrBlocks:  []string{anyIPv4CidrBlock},
			},
		}, nil
	case v1alpha1.SecurityGroupControlPlane:
		return v1alpha1.IngressRules{
			s.defaultSSHIngressRule(s.scope.SecurityGroups()[v1alpha1.SecurityGroupBastion].ID),
			{
				Description: "Kubernetes API",
				Protocol:    v1alpha1.SecurityGroupProtocolTCP,
				FromPort:    6443,
				ToPort:      6443,
				CidrBlocks:  []string{anyIPv4CidrBlock},
			},
			{
				Description:            "etcd",
				Protocol:               v1alpha1.SecurityGroupProtocolTCP,
				FromPort:               2379,
				ToPort:                 2379,
				SourceSecurityGroupIDs: []string{s.scope.SecurityGroups()[v1alpha1.SecurityGroupControlPlane].ID},
			},
			{
				Description:            "etcd peer",
				Protocol:               v1alpha1.SecurityGroupProtocolTCP,
				FromPort:               2380,
				ToPort:                 2380,
				SourceSecurityGroupIDs: []string{s.scope.SecurityGroups()[v1alpha1.SecurityGroupControlPlane].ID},
			},
			{
				Description: "bgp (calico)",
				Protocol:    v1alpha1.SecurityGroupProtocolTCP,
				FromPort:    179,
				ToPort:      179,
				SourceSecurityGroupIDs: []string{
					s.scope.SecurityGroups()[v1alpha1.SecurityGroupControlPlane].ID,
					s.scope.SecurityGroups()[v1alpha1.SecurityGroupNode].ID,
				},
			},
		}, nil

	case v1alpha1.SecurityGroupNode:
		return v1alpha1.IngressRules{
			s.defaultSSHIngressRule(s.scope.SecurityGroups()[v1alpha1.SecurityGroupBastion].ID),
			{
				Description: "Node Port Services",
				Protocol:    v1alpha1.SecurityGroupProtocolTCP,
				FromPort:    30000,
				ToPort:      32767,
				CidrBlocks:  []string{anyIPv4CidrBlock},
			},
			{
				Description:            "Kubelet API",
				Protocol:               v1alpha1.SecurityGroupProtocolTCP,
				FromPort:               10250,
				ToPort:                 10250,
				SourceSecurityGroupIDs: []string{s.scope.SecurityGroups()[v1alpha1.SecurityGroupControlPlane].ID},
			},
			{
				Description: "bgp (calico)",
				Protocol:    v1alpha1.SecurityGroupProtocolTCP,
				FromPort:    179,
				ToPort:      179,
				SourceSecurityGroupIDs: []string{
					s.scope.SecurityGroups()[v1alpha1.SecurityGroupControlPlane].ID,
					s.scope.SecurityGroups()[v1alpha1.SecurityGroupNode].ID,
				},
			},
		}, nil
	}

	return nil, errors.Errorf("Cannot determine ingress rules for unknown security group role %q", role)
}

func (s *Service) getSecurityGroupName(clusterName string, role v1alpha1.SecurityGroupRole) string {
	return fmt.Sprintf("%s-%v", clusterName, role)
}

func (s *Service) getDefaultSecurityGroup(role v1alpha1.SecurityGroupRole) *ec2.SecurityGroup {
	name := s.getSecurityGroupName(s.scope.Name(), role)

	return &ec2.SecurityGroup{
		GroupName: aws.String(name),
		VpcId:     aws.String(s.scope.VPC().ID),
		Tags:      converters.MapToTags(tags.Build(s.getSecurityGroupTagParams(name, role))),
	}
}

func (s *Service) getSecurityGroupTagParams(name string, role v1alpha1.SecurityGroupRole) tags.BuildParams {
	return tags.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   tags.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(string(role)),
	}
}

func ingressRuleToSDKType(i *v1alpha1.IngressRule) (res *ec2.IpPermission) {
	// AWS seems to ignore the From/To port when set on protocols where it doesn't apply, but
	// we avoid serializing it out for clarity's sake.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
	switch i.Protocol {
	case v1alpha1.SecurityGroupProtocolTCP,
		v1alpha1.SecurityGroupProtocolUDP,
		v1alpha1.SecurityGroupProtocolICMP,
		v1alpha1.SecurityGroupProtocolICMPv6:
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

func ingressRuleFromSDKType(v *ec2.IpPermission) (res *v1alpha1.IngressRule) {
	// Ports are only well-defined for TCP and UDP protocols, but EC2 overloads the port range
	// in the case of ICMP(v6) traffic to indicate which codes are allowed. For all other protocols,
	// including the custom "-1" All Traffic protcol, FromPort and ToPort are omitted from the response.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
	switch *v.IpProtocol {
	case IPProtocolTCP,
		IPProtocolUDP,
		IPProtocolICMP,
		IPProtocolICMPv6:
		res = &v1alpha1.IngressRule{
			Protocol: v1alpha1.SecurityGroupProtocol(*v.IpProtocol),
			FromPort: *v.FromPort,
			ToPort:   *v.ToPort,
		}
	default:
		res = &v1alpha1.IngressRule{
			Protocol: v1alpha1.SecurityGroupProtocol(*v.IpProtocol),
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

// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

func (s *Service) reconcileSecurityGroups(clusterName string, network *v1alpha1.Network) error {
	glog.V(2).Infof("Reconciling security groups")

	if network.SecurityGroups == nil {
		network.SecurityGroups = make(map[v1alpha1.SecurityGroupRole]*v1alpha1.SecurityGroup)
	}

	sgs, err := s.describeSecurityGroupsByName(clusterName, network.VPC.ID)
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
		sg := s.getDefaultSecurityGroup(clusterName, network.VPC.ID, role)

		if existing, ok := sgs[*sg.GroupName]; !ok {
			if err := s.createSecurityGroup(clusterName, role, sg); err != nil {
				return err
			}

			network.SecurityGroups[role] = &v1alpha1.SecurityGroup{
				ID:   *sg.GroupId,
				Name: *sg.GroupName,
			}
		} else {
			// TODO(vincepri): validate / update security group if necessary.
			network.SecurityGroups[role] = existing
		}

		glog.V(2).Infof("Security group for role %q: %v", role, network.SecurityGroups[role])
	}

	// Second iteration creates or updates all permissions on the security group to match
	// the specified ingress rules.
	for role, sg := range network.SecurityGroups {
		current := sg.IngressRules
		want, err := s.getSecurityGroupIngressRules(role, network)
		if err != nil {
			return err
		}

		toAuthorize := want.Difference(current)
		if len(toAuthorize) > 0 {
			if err := s.authorizeSecurityGroupIngressRules(sg.ID, toAuthorize); err != nil {
				return err
			}

			glog.V(2).Infof("Authorized ingress rules %v in security group %q", toAuthorize, sg)
		}

		toRevoke := current.Difference(want)
		if len(toRevoke) > 0 {
			if err := s.revokeSecurityGroupIngressRules(sg.ID, toRevoke); err != nil {
				return err
			}

			glog.V(2).Infof("Revoked ingress rules %v from security group %q", toRevoke, sg)
		}
	}

	return nil
}

func (s *Service) describeSecurityGroupsByName(clusterName string, vpcID string) (map[string]*v1alpha1.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	input.Filters = s.addTagFilters(clusterName, input.Filters)

	out, err := s.EC2.DescribeSecurityGroups(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe security groups in vpc %q", vpcID)
	}

	res := make(map[string]*v1alpha1.SecurityGroup, len(out.SecurityGroups))
	for _, ec2sg := range out.SecurityGroups {
		sg := &v1alpha1.SecurityGroup{
			ID:   *ec2sg.GroupId,
			Name: *ec2sg.GroupName,
		}

		for _, ec2rule := range ec2sg.IpPermissions {
			sg.IngressRules = append(sg.IngressRules, ingressRuleFromSDKType(ec2rule))
		}

		res[sg.Name] = sg
	}

	return res, nil
}

func (s *Service) createSecurityGroup(clusterName string, role v1alpha1.SecurityGroupRole, input *ec2.SecurityGroup) error {
	out, err := s.EC2.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		VpcId:       input.VpcId,
		GroupName:   input.GroupName,
		Description: aws.String(fmt.Sprintf("Kubernetes cluster %s: %s", clusterName, role)),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to create security group %q in vpc %q", *input.GroupName, *input.VpcId)
	}

	// Set the group id.
	input.GroupId = out.GroupId

	// Tag the security group.
	if _, err := s.EC2.CreateTags(&ec2.CreateTagsInput{Resources: []*string{out.GroupId}, Tags: input.Tags}); err != nil {
		return errors.Wrapf(err, "failed to tag security group %q in vpc %q", *input.GroupName, *input.VpcId)
	}

	return nil
}

func (s *Service) authorizeSecurityGroupIngressRules(groupID string, rules v1alpha1.IngressRules) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{GroupId: aws.String(groupID)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.EC2.AuthorizeSecurityGroupIngress(input); err != nil {
		return errors.Wrapf(err, "failed to authorize security group ingress rules for %q", groupID)
	}

	return nil
}

func (s *Service) revokeSecurityGroupIngressRules(groupID string, rules v1alpha1.IngressRules) error {
	input := &ec2.RevokeSecurityGroupIngressInput{GroupId: aws.String(groupID)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, ingressRuleToSDKType(rule))
	}

	if _, err := s.EC2.RevokeSecurityGroupIngress(input); err != nil {
		return errors.Wrapf(err, "failed to revoke security group ingress rules for %q", groupID)
	}

	return nil
}

func (s *Service) getSecurityGroupIngressRules(role v1alpha1.SecurityGroupRole, network *v1alpha1.Network) (v1alpha1.IngressRules, error) {
	switch role {
	case v1alpha1.SecurityGroupBastion:
		return v1alpha1.IngressRules{
			{
				Description: "SSH",
				Protocol:    "tcp",
				FromPort:    22,
				ToPort:      22,
				CidrBlocks:  []string{"0.0.0.0/0"},
			},
		}, nil
	case v1alpha1.SecurityGroupControlPlane:
		return v1alpha1.IngressRules{
			{
				Description: "SSH",
				Protocol:    "tcp",
				FromPort:    22,
				ToPort:      22,
				CidrBlocks:  []string{"0.0.0.0/0"},
			},
			{
				Description: "Kubernetes API",
				Protocol:    "tcp",
				FromPort:    6443,
				ToPort:      6443,
				CidrBlocks:  []string{"0.0.0.0/0"},
			},
			{
				Description:           "etcd",
				Protocol:              "tcp",
				FromPort:              2379,
				ToPort:                2379,
				SourceSecurityGroupID: aws.String(network.SecurityGroups[v1alpha1.SecurityGroupControlPlane].ID),
			},
			{
				Description:           "etcd peer",
				Protocol:              "tcp",
				FromPort:              2380,
				ToPort:                2380,
				SourceSecurityGroupID: aws.String(network.SecurityGroups[v1alpha1.SecurityGroupControlPlane].ID),
			},
		}, nil

	case v1alpha1.SecurityGroupNode:
		return v1alpha1.IngressRules{
			{
				Description: "SSH",
				Protocol:    "tcp",
				FromPort:    22,
				ToPort:      22,
				CidrBlocks:  []string{"0.0.0.0/0"},
			},
			{
				Description: "Node Port Services",
				Protocol:    "tcp",
				FromPort:    30000,
				ToPort:      32767,
				CidrBlocks:  []string{"0.0.0.0/0"},
			},
			{
				Description:           "Kubelet API",
				Protocol:              "tcp",
				FromPort:              10250,
				ToPort:                10250,
				SourceSecurityGroupID: aws.String(network.SecurityGroups[v1alpha1.SecurityGroupControlPlane].ID),
			},
		}, nil
	}

	return nil, errors.Errorf("Cannot determine ingress rules for unknown security group role %q", role)
}

func (s *Service) getSecurityGroupName(clusterName string, role v1alpha1.SecurityGroupRole) string {
	return fmt.Sprintf("%s-%v", clusterName, role)
}

func (s *Service) getDefaultSecurityGroup(clusterName string, vpcID string, role v1alpha1.SecurityGroupRole) *ec2.SecurityGroup {
	name := s.getSecurityGroupName(clusterName, role)

	return &ec2.SecurityGroup{
		GroupName: aws.String(name),
		VpcId:     aws.String(vpcID),
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
			&ec2.Tag{
				Key:   aws.String(s.clusterTagKey(clusterName)),
				Value: aws.String(string(ResourceLifecycleOwned)),
			},
		},
	}
}

func ingressRuleToSDKType(i *v1alpha1.IngressRule) *ec2.IpPermission {
	res := &ec2.IpPermission{
		IpProtocol: aws.String(i.Protocol),
		FromPort:   aws.Int64(i.FromPort),
		ToPort:     aws.Int64(i.ToPort),
	}

	for _, cidr := range i.CidrBlocks {
		res.IpRanges = append(res.IpRanges, &ec2.IpRange{
			Description: aws.String(i.Description),
			CidrIp:      aws.String(cidr),
		})
	}

	if i.SourceSecurityGroupID != nil {
		res.UserIdGroupPairs = append(res.UserIdGroupPairs, &ec2.UserIdGroupPair{
			Description: aws.String(i.Description),
			GroupId:     aws.String(*i.SourceSecurityGroupID),
		})
	}

	return res
}

func ingressRuleFromSDKType(v *ec2.IpPermission) *v1alpha1.IngressRule {
	res := &v1alpha1.IngressRule{
		Protocol: *v.IpProtocol,
		FromPort: *v.FromPort,
		ToPort:   *v.ToPort,
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

		res.SourceSecurityGroupID = pair.GroupId
	}

	return res
}

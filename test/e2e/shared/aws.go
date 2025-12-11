//go:build e2e
// +build e2e

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

package shared

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awscredsv2 "github.com/aws/aws-sdk-go-v2/credentials"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailtypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	configservicetypes "github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	ecrpublictypes "github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	servicequotastypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cloudformation/service"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
)

type AWSInfrastructureSpec struct {
	ClusterName            string `json:"clustername"`
	VpcCidr                string `json:"vpcCidr"`
	PublicSubnetCidr       string `json:"publicSubnetCidr"`
	PrivateSubnetCidr      string `json:"privateSubnetCidr"`
	AvailabilityZone       string `json:"availabilityZone"`
	ExternalSecurityGroups bool   `json:"externalSecurityGroups"`
}

type AWSInfrastructureState struct {
	PrivateSubnetID     *string `json:"privateSubnetID"`
	PrivateSubnetState  *string `json:"privateSubnetState"`
	PublicSubnetID      *string `json:"publicSubnetID"`
	PublicSubnetState   *string `json:"publicSubnetState"`
	VpcState            *string `json:"vpcState"`
	NatGatewayState     *string `json:"natGatewayState"`
	PublicRouteTableID  *string `json:"publicRouteTableID"`
	PrivateRouteTableID *string `json:"privateRouteTableID"`
}
type AWSInfrastructure struct {
	Spec            AWSInfrastructureSpec `json:"spec"`
	Context         *E2EContext
	VPC             *ec2types.Vpc
	Subnets         []*ec2types.Subnet
	RouteTables     []*ec2types.RouteTable
	InternetGateway *ec2types.InternetGateway
	ElasticIP       *ec2types.Address
	NatGateway      *ec2types.NatGateway
	State           AWSInfrastructureState `json:"state"`
	Peering         *ec2types.VpcPeeringConnection
}

func (i *AWSInfrastructure) New(ais AWSInfrastructureSpec, e2eCtx *E2EContext) AWSInfrastructure {
	i.Spec = ais
	i.Context = e2eCtx
	return *i
}

func (i *AWSInfrastructure) CreateVPC() AWSInfrastructure {
	cv, err := CreateVPC(i.Context, i.Spec.ClusterName+"-vpc", i.Spec.VpcCidr)
	if err != nil {
		i.State.VpcState = ptr.To[string](fmt.Sprintf("failed: %v", err))
		return *i
	}

	i.VPC = cv
	i.State.VpcState = aws.String(string(cv.State))
	return *i
}

func (i *AWSInfrastructure) RefreshVPCState() AWSInfrastructure {
	if i.VPC == nil {
		return *i
	}
	vpc, err := GetVPC(i.Context, *i.VPC.VpcId)
	if err != nil {
		return *i
	}
	if vpc != nil {
		i.VPC = vpc
		i.State.VpcState = aws.String(string(vpc.State))
	}
	return *i
}

func (i *AWSInfrastructure) CreatePublicSubnet() AWSInfrastructure {
	subnet, err := CreateSubnet(i.Context, i.Spec.ClusterName, i.Spec.PublicSubnetCidr, i.Spec.AvailabilityZone, *i.VPC.VpcId, "public")
	if err != nil {
		i.State.PublicSubnetState = ptr.To[string]("failed")
		return *i
	}
	i.State.PublicSubnetID = subnet.SubnetId
	i.State.PublicSubnetState = aws.String(string(subnet.State))
	i.Subnets = append(i.Subnets, subnet)
	return *i
}

func (i *AWSInfrastructure) CreatePrivateSubnet() AWSInfrastructure {
	subnet, err := CreateSubnet(i.Context, i.Spec.ClusterName, i.Spec.PrivateSubnetCidr, i.Spec.AvailabilityZone, *i.VPC.VpcId, "private")
	if err != nil {
		i.State.PrivateSubnetState = ptr.To[string]("failed")
		return *i
	}
	i.State.PrivateSubnetID = subnet.SubnetId
	i.State.PrivateSubnetState = aws.String(string(subnet.State))
	i.Subnets = append(i.Subnets, subnet)
	return *i
}

func (i *AWSInfrastructure) CreateInternetGateway() AWSInfrastructure {
	igwC, err := CreateInternetGateway(i.Context, i.Spec.ClusterName+"-igw")
	if err != nil {
		return *i
	}
	_, aerr := AttachInternetGateway(i.Context, *igwC.InternetGatewayId, *i.VPC.VpcId)
	if aerr != nil {
		i.InternetGateway = igwC
		return *i
	}
	i.InternetGateway = igwC
	return *i
}

func (i *AWSInfrastructure) AllocateAddress() AWSInfrastructure {
	aa, err := AllocateAddress(i.Context, i.Spec.ClusterName+"-eip")
	if err != nil {
		return *i
	}

	var addr *ec2types.Address
	Eventually(func(gomega Gomega) {
		addr, _ = GetAddress(i.Context, *aa.AllocationId)
	}, 2*time.Minute, 5*time.Second).Should(Succeed())
	i.ElasticIP = addr
	return *i
}

func (i *AWSInfrastructure) CreateNatGateway(ct string) AWSInfrastructure {
	var s *ec2types.Subnet
	Eventually(func(gomega Gomega) {
		s, _ = GetSubnetByName(i.Context, i.Spec.ClusterName+"-subnet-"+ct)
	}, 2*time.Minute, 5*time.Second).Should(Succeed())
	if s == nil {
		return *i
	}
	ngwC, ngwce := CreateNatGateway(i.Context, i.Spec.ClusterName+"-nat", ct, *i.ElasticIP.AllocationId, *s.SubnetId)
	if ngwce != nil {
		return *i
	}
	if WaitForNatGatewayState(i.Context, *ngwC.NatGatewayId, "available") {
		ngw, _ := GetNatGateway(i.Context, *ngwC.NatGatewayId)
		i.NatGateway = ngw
		i.State.NatGatewayState = aws.String(string(ngw.State))
		return *i
	}
	i.NatGateway = ngwC
	return *i
}

func (i *AWSInfrastructure) CreateRouteTable(subnetType string) AWSInfrastructure {
	rt, err := CreateRouteTable(i.Context, i.Spec.ClusterName+"-rt-"+subnetType, *i.VPC.VpcId)
	if err != nil {
		return *i
	}
	switch subnetType {
	case "public":
		if a, _ := AssociateRouteTable(i.Context, *rt.RouteTableId, *i.State.PublicSubnetID); a != nil {
			i.State.PublicRouteTableID = rt.RouteTableId
		}
	case "private":
		if a, _ := AssociateRouteTable(i.Context, *rt.RouteTableId, *i.State.PrivateSubnetID); a != nil {
			i.State.PrivateRouteTableID = rt.RouteTableId
		}
	}
	return *i
}

func (i *AWSInfrastructure) GetRouteTable(rtID string) AWSInfrastructure {
	rt, err := GetRouteTable(i.Context, rtID)
	if err != nil {
		return *i
	}
	if rt != nil {
		i.RouteTables = append(i.RouteTables, rt)
	}
	return *i
}

// CreateInfrastructure creates a VPC, two subnets with appropriate tags based on type(private/public)
// an internet gateway, an elastic IP address, a NAT gateway, a route table for each subnet and
// routes to their respective gateway.
func (i *AWSInfrastructure) CreateInfrastructure() AWSInfrastructure {
	i.CreateVPC()
	Eventually(func() string {
		return *i.RefreshVPCState().State.VpcState
	}, 2*time.Minute, 5*time.Second).Should(Equal("available"), "Expected VPC state to eventually become 'available'")

	By(fmt.Sprintf("Created VPC - %s", *i.VPC.VpcId))
	if i.VPC != nil {
		i.CreatePublicSubnet()
		if i.State.PublicSubnetID != nil {
			By(fmt.Sprintf("Created Public Subnet - %s", *i.State.PublicSubnetID))
		}
		i.CreatePrivateSubnet()
		if i.State.PrivateSubnetID != nil {
			By(fmt.Sprintf("Created Private Subnet - %s", *i.State.PrivateSubnetID))
		}
		i.CreateInternetGateway()
		if i.InternetGateway != nil {
			By(fmt.Sprintf("Created Internet Gateway - %s", *i.InternetGateway.InternetGatewayId))
		}
	}
	i.AllocateAddress()
	if i.ElasticIP != nil && i.ElasticIP.AllocationId != nil {
		By(fmt.Sprintf("Created Elastic IP - %s", *i.ElasticIP.AllocationId))
		i.CreateNatGateway("public")
		if i.NatGateway != nil && i.NatGateway.NatGatewayId != nil {
			WaitForNatGatewayState(i.Context, *i.NatGateway.NatGatewayId, "available")
			By(fmt.Sprintf("Created NAT Gateway - %s", *i.NatGateway.NatGatewayId))
		}
	}
	if len(i.Subnets) == 2 {
		i.CreateRouteTable("public")
		if i.State.PublicRouteTableID != nil {
			By(fmt.Sprintf("Created public route table - %s", *i.State.PublicRouteTableID))
		}
		i.CreateRouteTable("private")
		if i.State.PrivateRouteTableID != nil {
			By(fmt.Sprintf("Created private route table - %s", *i.State.PrivateRouteTableID))
		}
		if i.InternetGateway != nil && i.InternetGateway.InternetGatewayId != nil {
			CreateRoute(i.Context, *i.State.PublicRouteTableID, "0.0.0.0/0", nil, i.InternetGateway.InternetGatewayId, nil)
		}
		if i.NatGateway != nil && i.NatGateway.NatGatewayId != nil {
			CreateRoute(i.Context, *i.State.PrivateRouteTableID, "0.0.0.0/0", i.NatGateway.NatGatewayId, nil, nil)
		}
		if i.State.PublicRouteTableID != nil {
			i.GetRouteTable(*i.State.PublicRouteTableID)
		}
		if i.State.PrivateRouteTableID != nil {
			i.GetRouteTable(*i.State.PrivateRouteTableID)
		}
	}
	return *i
}

// DeleteInfrastructure has calls added to discover and delete potential orphaned resources created
// by CAPA. In an attempt to avoid dependency violations it works in the following order
// Instances, Load Balancers, Route Tables, NAT gateway, Elastic IP, Internet Gateway,
// Security Group Rules, Security Groups, Subnets, VPC.
func (i *AWSInfrastructure) DeleteInfrastructure(ctx context.Context) {
	instances, _ := ListClusterEC2Instances(i.Context, i.Spec.ClusterName)
	for _, instance := range instances {
		if instance.State.Code != aws.Int32(48) {
			By(fmt.Sprintf("Deleting orphaned instance: %s - %v", *instance.InstanceId, TerminateInstance(i.Context, *instance.InstanceId)))
		}
	}
	WaitForInstanceState(i.Context, i.Spec.ClusterName, "terminated")

	loadbalancers, _ := ListLoadBalancers(ctx, i.Context, i.Spec.ClusterName)
	for _, lb := range loadbalancers {
		By(fmt.Sprintf("Deleting orphaned load balancer: %s - %v", *lb.LoadBalancerName, DeleteLoadBalancer(ctx, i.Context, *lb.LoadBalancerName)))
	}

	for _, rt := range i.RouteTables {
		for _, a := range rt.Associations {
			By(fmt.Sprintf("Disassociating route table - %s - %v", *a.RouteTableAssociationId, DisassociateRouteTable(i.Context, *a.RouteTableAssociationId)))
		}
		By(fmt.Sprintf("Deleting route table - %s - %v", *rt.RouteTableId, DeleteRouteTable(i.Context, *rt.RouteTableId)))
	}

	if i.NatGateway != nil {
		By(fmt.Sprintf("Deleting NAT Gateway - %s - %v", *i.NatGateway.NatGatewayId, DeleteNatGateway(i.Context, *i.NatGateway.NatGatewayId)))
		WaitForNatGatewayState(i.Context, *i.NatGateway.NatGatewayId, "deleted")
	}

	if i.ElasticIP != nil {
		By(fmt.Sprintf("Deleting Elastic IP - %s - %v", *i.ElasticIP.AllocationId, ReleaseAddress(i.Context, *i.ElasticIP.AllocationId)))
	}

	if i.InternetGateway != nil {
		By(fmt.Sprintf("Detaching Internet Gateway - %s - %v", *i.InternetGateway.InternetGatewayId, DetachInternetGateway(i.Context, *i.InternetGateway.InternetGatewayId, *i.VPC.VpcId)))
		By(fmt.Sprintf("Deleting Internet Gateway - %s - %v", *i.InternetGateway.InternetGatewayId, DeleteInternetGateway(i.Context, *i.InternetGateway.InternetGatewayId)))
	}

	sgGroups, _ := GetSecurityGroupsByVPC(i.Context, *i.VPC.VpcId)
	for _, sg := range sgGroups {
		if *sg.GroupName != "default" {
			sgRules, _ := ListSecurityGroupRules(i.Context, *sg.GroupId)
			for _, sgr := range sgRules {
				var d bool
				if *sgr.IsEgress {
					for d = DeleteSecurityGroupRule(i.Context, *sgr.GroupId, *sgr.SecurityGroupRuleId, "egress"); !d; {
						d = DeleteSecurityGroupRule(i.Context, *sgr.GroupId, *sgr.SecurityGroupRuleId, "egress")
					}
					By(fmt.Sprintf("Deleting Egress Security Group Rule - %s - %v", *sgr.SecurityGroupRuleId, d))
				} else {
					for d = DeleteSecurityGroupRule(i.Context, *sgr.GroupId, *sgr.SecurityGroupRuleId, "ingress"); !d; {
						d = DeleteSecurityGroupRule(i.Context, *sgr.GroupId, *sgr.SecurityGroupRuleId, "ingress")
					}
					By(fmt.Sprintf("Deleting Ingress Security Group Rule - %s - %v", *sgr.SecurityGroupRuleId, d))
				}
			}
		}
	}

	sgGroups, _ = GetSecurityGroupsByVPC(i.Context, *i.VPC.VpcId)
	for _, sg := range sgGroups {
		if *sg.GroupName != "default" {
			By(fmt.Sprintf("Deleting Security Group - %s - %v", *sg.GroupId, DeleteSecurityGroup(i.Context, *sg.GroupId)))
		}
	}

	for _, subnet := range i.Subnets {
		By(fmt.Sprintf("Deleting Subnet - %s - %v", *subnet.SubnetId, DeleteSubnet(i.Context, *subnet.SubnetId)))
	}

	if i.VPC != nil {
		By(fmt.Sprintf("Deleting VPC - %s - %v", *i.VPC.VpcId, DeleteVPC(i.Context, *i.VPC.VpcId)))
	}
}

func NewAWSSession() *aws.Config {
	By("Getting an AWS IAM session - from environment")
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	Expect(err).NotTo(HaveOccurred())
	_, err = cfg.Credentials.Retrieve(context.Background())
	Expect(err).NotTo(HaveOccurred())
	return &cfg
}

func NewAWSSessionRepoWithKey(accessKey *iamtypes.AccessKey) *aws.Config {
	By("Getting an AWS IAM session - from access key")
	region, err := credentials.ResolveRegion("us-east-1")
	Expect(err).NotTo(HaveOccurred())
	staticCredProvider := awscredsv2.NewStaticCredentialsProvider(aws.ToString(accessKey.AccessKeyId), aws.ToString(accessKey.SecretAccessKey), "")
	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(region),
		config.WithCredentialsProvider(staticCredProvider),
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	Expect(err).NotTo(HaveOccurred())
	_, err = cfg.Credentials.Retrieve(context.Background())
	Expect(err).NotTo(HaveOccurred())
	return &cfg
}

func NewAWSSessionWithKey(accessKey *iamtypes.AccessKey) *aws.Config {
	By("Getting an AWS IAM session - from access key")
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	staticCredProvider := awscredsv2.NewStaticCredentialsProvider(aws.ToString(accessKey.AccessKeyId), aws.ToString(accessKey.SecretAccessKey), "")
	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(region),
		config.WithCredentialsProvider(staticCredProvider),
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	Expect(err).NotTo(HaveOccurred())
	_, err = cfg.Credentials.Retrieve(context.Background())
	Expect(err).NotTo(HaveOccurred())
	return &cfg
}

// createCloudFormationStack ensures the cloudformation stack is up to date.
func createCloudFormationStack(ctx context.Context, cfg *aws.Config, t *cfn_bootstrap.Template, tags map[string]string) error {
	By(fmt.Sprintf("Creating AWS CloudFormation stack for AWS IAM resources: stack-name=%s", t.Spec.StackName))
	cfnClient := cfn.NewFromConfig(*cfg)
	// CloudFormation stack will clean up on a failure, we don't need an Eventually here.
	// The `create` already does a WaitUntilStackCreateComplete.
	cfnSvc := cloudformation.NewService(
		&cloudformation.CFNClient{
			Client: cfnClient,
		})

	if err := cfnSvc.ReconcileBootstrapNoUpdate(ctx, t.Spec.StackName, *renderCustomCloudFormation(t), tags); err != nil {
		By(fmt.Sprintf("Error reconciling Cloud formation stack %v", err))
		spewCloudFormationResources(ctx, cfnClient, t)

		// always clean up on a failure because we could leak these resources and the next cloud formation create would
		// fail with the same problem.
		deleteMultitenancyRoles(ctx, cfg)
		deleteResourcesInCloudFormation(ctx, cfg, t)
		return err
	}

	spewCloudFormationResources(ctx, cfnClient, t)
	return nil
}

func spewCloudFormationResources(ctx context.Context, cfnClient *cfn.Client, t *cfn_bootstrap.Template) {
	output, err := cfnClient.DescribeStackEvents(ctx, &cfn.DescribeStackEventsInput{StackName: aws.String(t.Spec.StackName), NextToken: aws.String("1")})
	if err != nil {
		By(fmt.Sprintf("Error describin Cloud formation stack events %v, skipping", err))
	} else {
		By("========= Stack Event Output Begin =========")
		for _, event := range output.StackEvents {
			By(fmt.Sprintf("Event details for %s : Resource: %s, Status: %s, Reason: %s", aws.ToString(event.LogicalResourceId), aws.ToString(event.ResourceType), event.ResourceStatus, aws.ToString(event.ResourceStatusReason)))
		}
		By("========= Stack Event Output End =========")
	}
	out, err := cfnClient.DescribeStackResources(ctx, &cfn.DescribeStackResourcesInput{
		StackName: aws.String(t.Spec.StackName),
	})
	if err != nil {
		By(fmt.Sprintf("Error describing Stack Resources %v, skipping", err))
	} else {
		By("========= Stack Resources Output Begin =========")
		By("Resource\tType\tStatus")

		for _, r := range out.StackResources {
			By(fmt.Sprintf("%s\t%s\t%s\t%s",
				aws.ToString(r.ResourceType),
				aws.ToString(r.PhysicalResourceId),
				r.ResourceStatus,
				aws.ToString(r.ResourceStatusReason)))
		}
		By("========= Stack Resources Output End =========")
	}
}

func SetMultitenancyEnvVars(ctx context.Context, cfg *aws.Config) error {
	for _, roles := range MultiTenancyRoles {
		if err := roles.SetEnvVars(ctx, cfg); err != nil {
			return err
		}
	}
	return nil
}

// Delete resources that already exists.
func deleteResourcesInCloudFormation(ctx context.Context, cfg *aws.Config, t *cfn_bootstrap.Template) {
	iamSvc := iam.NewFromConfig(*cfg)
	temp := *renderCustomCloudFormation(t)
	var (
		iamUsers         []*cfn_iam.User
		iamRoles         []*cfn_iam.Role
		instanceProfiles []*cfn_iam.InstanceProfile
		policies         []*cfn_iam.ManagedPolicy
		groups           []*cfn_iam.Group
	)
	// the deletion order of these resources is important. Policies need to be last,
	// so they don't have any attached resources which prevents their deletion.
	// temp.Resources is a map. Traversing that directly results in undetermined order.
	for _, val := range temp.Resources {
		switch configservicetypes.ResourceType(val.AWSCloudFormationType()) {
		case configservicetypes.ResourceTypeUser:
			user := val.(*cfn_iam.User)
			iamUsers = append(iamUsers, user)
		case configservicetypes.ResourceTypeRole:
			role := val.(*cfn_iam.Role)
			iamRoles = append(iamRoles, role)
		case configservicetypes.ResourceTypeIAMInstanceProfile:
			profile := val.(*cfn_iam.InstanceProfile)
			instanceProfiles = append(instanceProfiles, profile)
		case configservicetypes.ResourceType("AWS::IAM::ManagedPolicy"):
			policy := val.(*cfn_iam.ManagedPolicy)
			policies = append(policies, policy)
		case configservicetypes.ResourceTypeGroup:
			group := val.(*cfn_iam.Group)
			groups = append(groups, group)
		}
	}
	for _, user := range iamUsers {
		By(fmt.Sprintf("deleting the following user: %q", user.UserName))
		repeat := false
		Eventually(func(gomega Gomega) bool {
			err := DeleteUser(ctx, cfg, user.UserName)
			if err != nil && !repeat {
				By(fmt.Sprintf("failed to delete user '%q'; reason: %+v", user.UserName, err))
				repeat = true
			}
			var noSuchEntityErr *iamtypes.NoSuchEntityException
			return err == nil || errors.As(err, &noSuchEntityErr)
		}, 5*time.Minute, 5*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed deleting the user: %q", user.UserName))
	}
	for _, role := range iamRoles {
		By(fmt.Sprintf("deleting the following role: %s", role.RoleName))
		repeat := false
		Eventually(func(gomega Gomega) bool {
			err := DeleteRole(ctx, cfg, role.RoleName)
			if err != nil && !repeat {
				By(fmt.Sprintf("failed to delete role '%s'; reason: %+v", role.RoleName, err))
				repeat = true
			}
			var noSuchEntityErr *iamtypes.NoSuchEntityException
			return err == nil || errors.As(err, &noSuchEntityErr)
		}, 5*time.Minute, 5*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed deleting the following role: %q", role.RoleName))
	}
	for _, profile := range instanceProfiles {
		By(fmt.Sprintf("cleanup for profile with name '%s'", profile.InstanceProfileName))
		repeat := false
		Eventually(func(gomega Gomega) bool {
			_, err := iamSvc.DeleteInstanceProfile(ctx, &iam.DeleteInstanceProfileInput{
				InstanceProfileName: aws.String(profile.InstanceProfileName),
			})
			if err != nil && !repeat {
				By(fmt.Sprintf("failed to delete role '%s'; reason: %+v", profile.InstanceProfileName, err))
				repeat = true
			}
			var noSuchEntityErr *iamtypes.NoSuchEntityException
			return err == nil || errors.As(err, &noSuchEntityErr)
		}, 5*time.Minute, 5*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed cleaning up profile with name %q", profile.InstanceProfileName))
	}
	for _, group := range groups {
		repeat := false
		Eventually(func(gomega Gomega) bool {
			_, err := iamSvc.DeleteGroup(ctx, &iam.DeleteGroupInput{
				GroupName: aws.String(group.GroupName),
			})
			if err != nil && !repeat {
				By(fmt.Sprintf("failed to delete group '%s'; reason: %+v", group.GroupName, err))
				repeat = true
			}
			var noSuchEntityErr *iamtypes.NoSuchEntityException
			return err == nil || errors.As(err, &noSuchEntityErr)
		}, 5*time.Minute, 5*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed deleting group %q", group.GroupName))
	}
	for _, policy := range policies {
		listPoliciesOutput, err := iamSvc.ListPolicies(ctx, &iam.ListPoliciesInput{})
		Expect(err).NotTo(HaveOccurred())
		if len(listPoliciesOutput.Policies) > 0 {
			for _, p := range listPoliciesOutput.Policies {
				if aws.ToString(p.PolicyName) == policy.ManagedPolicyName {
					By(fmt.Sprintf("cleanup for policy '%s'", aws.ToString(p.PolicyName)))
					repeat := false
					Eventually(func(gomega Gomega) bool {
						_, err := iamSvc.DeletePolicy(ctx, &iam.DeletePolicyInput{
							PolicyArn: p.Arn,
						})
						if err != nil && !repeat {
							By(fmt.Sprintf("failed to delete policy '%s'; reason: %+v", policy.Description, err))
							repeat = true
						}
						var noSuchEntityErr *iamtypes.NoSuchEntityException
						return err == nil || errors.As(err, &noSuchEntityErr)
					}, 5*time.Minute, 5*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed to delete policy %q", aws.ToString(p.PolicyName)))
					// TODO: why is there a break here? Don't we want to clean up everything?
					break
				}
			}
		}
	}
}

// TODO: remove once test infra accounts are fixed.
func deleteMultitenancyRoles(ctx context.Context, cfg *aws.Config) {
	if err := DeleteRole(ctx, cfg, "multi-tenancy-role"); err != nil {
		By(fmt.Sprintf("failed to delete role multi-tenancy-role %s", err))
	}
	if err := DeleteRole(ctx, cfg, "multi-tenancy-nested-role"); err != nil {
		By(fmt.Sprintf("failed to delete role multi-tenancy-nested-role %s", err))
	}
}

// detachAllPoliciesForRole detaches all policies for role.
func detachAllPoliciesForRole(ctx context.Context, cfg *aws.Config, name string) error {
	iamSvc := iam.NewFromConfig(*cfg)

	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(name),
	}

	policies, err := iamSvc.ListAttachedRolePolicies(ctx, input)
	if err != nil {
		return errors.New("error fetching policies for role")
	}

	for _, p := range policies.AttachedPolicies {
		input := &iam.DetachRolePolicyInput{
			RoleName:  aws.String(name),
			PolicyArn: p.PolicyArn,
		}

		_, err := iamSvc.DetachRolePolicy(ctx, input)
		if err != nil {
			return errors.New("failed detaching policy from a role")
		}
	}
	return nil
}

// DeleteUser deletes an IAM user in a best effort manner.
func DeleteUser(ctx context.Context, cfg *aws.Config, name string) error {
	iamSvc := iam.NewFromConfig(*cfg)

	// if user does not exist, return.
	_, err := iamSvc.GetUser(ctx, &iam.GetUserInput{UserName: aws.String(name)})
	if err != nil {
		return err
	}

	_, err = iamSvc.DeleteUser(ctx, &iam.DeleteUserInput{UserName: aws.String(name)})
	if err != nil {
		return err
	}

	return nil
}

// DeleteRole deletes roles in a best effort manner.
func DeleteRole(ctx context.Context, cfg *aws.Config, name string) error {
	iamSvc := iam.NewFromConfig(*cfg)

	// if role does not exist, return.
	_, err := iamSvc.GetRole(ctx, &iam.GetRoleInput{RoleName: aws.String(name)})
	if err != nil {
		return err
	}

	if err := detachAllPoliciesForRole(ctx, cfg, name); err != nil {
		return err
	}

	_, err = iamSvc.DeleteRole(ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
	if err != nil {
		return err
	}

	return nil
}

func GetPolicyArn(ctx context.Context, cfg aws.Config, name string) string {
	iamSvc := iam.NewFromConfig(cfg)

	policyList, err := iamSvc.ListPolicies(ctx, &iam.ListPoliciesInput{
		Scope: iamtypes.PolicyScopeTypeLocal,
	})
	Expect(err).NotTo(HaveOccurred())

	for _, policy := range policyList.Policies {
		if aws.ToString(policy.PolicyName) == name {
			return aws.ToString(policy.Arn)
		}
	}
	return ""
}

func logAccountDetails(cfg *aws.Config) {
	By("Getting AWS account details")
	stsSvc := sts.NewFromConfig(*cfg)

	output, err := stsSvc.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Couldn't get sts caller identity: err=%s\n", err)
		return
	}

	fmt.Fprintf(GinkgoWriter, "Using AWS account: %s\n", *output.Account)
}

// deleteCloudFormationStack removes the provisioned clusterawsadm stack.
func deleteCloudFormationStack(cfg *aws.Config, t *cfn_bootstrap.Template) {
	By(fmt.Sprintf("Deleting %s CloudFormation stack", t.Spec.StackName))
	ctx := context.TODO()
	CFN := cfn.NewFromConfig(*cfg)
	cfnSvc := cloudformation.NewService(
		&cloudformation.CFNClient{
			Client: CFN,
		})
	err := cfnSvc.DeleteStack(context.TODO(), t.Spec.StackName, nil)
	if err != nil {
		var retainResources []string
		out, err := CFN.DescribeStackResources(ctx, &cfn.DescribeStackResourcesInput{StackName: aws.String(t.Spec.StackName)})
		Expect(err).NotTo(HaveOccurred())
		for _, v := range out.StackResources {
			if v.ResourceStatus == cfntypes.ResourceStatusDeleteFailed {
				retainResources = append(retainResources, aws.ToString(v.LogicalResourceId))
			}
		}
		err = cfnSvc.DeleteStack(ctx, t.Spec.StackName, retainResources)
		Expect(err).NotTo(HaveOccurred())
	}
	err = cfnSvc.CFN.WaitUntilStackDeleteComplete(ctx, &cfn.DescribeStacksInput{
		StackName: aws.String(t.Spec.StackName),
	}, cloudformation.MaxWaitCreateUpdateDelete)
	Expect(err).NotTo(HaveOccurred())
}

func ensureTestImageUploaded(ctx context.Context, e2eCtx *E2EContext) error {
	sessionForRepo := NewAWSSessionRepoWithKey(e2eCtx.Environment.BootstrapAccessKey)

	ecrSvc := ecrpublic.NewFromConfig(*sessionForRepo)
	repoName := ""
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		output, err := ecrSvc.CreateRepository(ctx, &ecrpublic.CreateRepositoryInput{
			RepositoryName: aws.String("capa/update"),
			CatalogData: &ecrpublictypes.RepositoryCatalogDataInput{
				AboutText: aws.String("Created by cluster-api-provider-aws/test/e2e/shared/aws.go for E2E tests"),
			},
		})

		if err != nil {
			if !awserrors.IsRepositoryExists(err) {
				return false, err
			}
			out, err := ecrSvc.DescribeRepositories(ctx, &ecrpublic.DescribeRepositoriesInput{RepositoryNames: []string{"capa/update"}})
			if err != nil || len(out.Repositories) == 0 {
				return false, err
			}
			repoName = aws.ToString(out.Repositories[0].RepositoryUri)
		} else {
			repoName = aws.ToString(output.Repository.RepositoryUri)
		}

		return true, nil
	}, awserrors.UnrecognizedClientException); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "docker", "inspect", "--format='{{index .Id}}'", "gcr.io/k8s-staging-cluster-api/capa-manager:e2e")
	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut
	err := cmd.Run()
	if err != nil {
		return err
	}

	imageSha := strings.ReplaceAll(strings.TrimSuffix(stdOut.String(), "\n"), "'", "")

	ecrImageName := repoName + ":e2e"
	cmd = exec.CommandContext(ctx, "docker", "tag", imageSha, ecrImageName) //nolint:gosec
	err = cmd.Run()
	if err != nil {
		return err
	}

	outToken, err := ecrSvc.GetAuthorizationToken(ctx, &ecrpublic.GetAuthorizationTokenInput{})
	if err != nil {
		return err
	}

	// Auth token is in username:password format. To login using it, we need to decode first and separate password and username
	decodedUsernamePassword, _ := b64.StdEncoding.DecodeString(aws.ToString(outToken.AuthorizationData.AuthorizationToken))

	strList := strings.Split(string(decodedUsernamePassword), ":")
	if len(strList) != 2 {
		return errors.New("failed to decode ECR authentication token")
	}

	cmd = exec.CommandContext(ctx, "docker", "login", "--username", strList[0], "--password", strList[1], "public.ecr.aws") //nolint:gosec
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.CommandContext(ctx, "docker", "push", ecrImageName)
	err = cmd.Run()
	if err != nil {
		return err
	}
	e2eCtx.E2EConfig.Variables["CAPI_IMAGES_REGISTRY"] = repoName
	e2eCtx.E2EConfig.Variables["E2E_IMAGE_TAG"] = "e2e"
	return nil
}

// ensureNoServiceLinkedRoles removes an auto-created IAM role, and tests
// the controller's IAM permissions to use ELB and Spot instances successfully.
func ensureNoServiceLinkedRoles(ctx context.Context, cfg *aws.Config) {
	iamSvc := iam.NewFromConfig(*cfg)

	By("Deleting AWS IAM Service Linked Role: role-name=AWSServiceRoleForElasticLoadBalancing")
	_, err := iamSvc.DeleteServiceLinkedRole(ctx, &iam.DeleteServiceLinkedRoleInput{
		RoleName: aws.String("AWSServiceRoleForElasticLoadBalancing"),
	})

	var noSuchEntityErr *iamtypes.NoSuchEntityException
	if !errors.As(err, &noSuchEntityErr) {
		Expect(err).NotTo(HaveOccurred())
	}

	By("Deleting AWS IAM Service Linked Role: role-name=AWSServiceRoleForEC2Spot")
	_, err = iamSvc.DeleteServiceLinkedRole(ctx, &iam.DeleteServiceLinkedRoleInput{
		RoleName: aws.String("AWSServiceRoleForEC2Spot"),
	})

	if !errors.As(err, &noSuchEntityErr) {
		Expect(err).NotTo(HaveOccurred())
	}
}

// ensureSSHKeyPair ensures A SSH key is present under the name.
func ensureSSHKeyPair(config aws.Config, keyPairName string) {
	By(fmt.Sprintf("Ensuring presence of SSH key in EC2: key-name=%s", keyPairName))
	ec2c := ec2.NewFromConfig(config)
	_, err := ec2c.CreateKeyPair(context.TODO(), &ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	var oe *smithy.OperationError
	if errors.As(err, &oe) && !strings.Contains(oe.Unwrap().Error(), "InvalidKeyPair.Duplicate") {
		Expect(err).NotTo(HaveOccurred())
	}
}

func ensureStackTags(cfg *aws.Config, stackName string, expectedTags map[string]string) {
	By(fmt.Sprintf("Ensuring AWS CloudFormation stack is created or updated with the specified tags: stack-name=%s", stackName))
	CFN := cfn.NewFromConfig(*cfg)
	r, err := CFN.DescribeStacks(context.TODO(), &cfn.DescribeStacksInput{StackName: &stackName})
	Expect(err).NotTo(HaveOccurred())
	stacks := r.Stacks
	Expect(len(stacks)).To(BeNumerically("==", 1))
	stackTags := stacks[0].Tags
	Expect(len(stackTags)).To(BeNumerically("==", len(expectedTags)))
	for _, tag := range stackTags {
		Expect(*tag.Value).To(BeIdenticalTo(expectedTags[*tag.Key]))
	}
}

// encodeCredentials leverages clusterawsadm to encode AWS credentials.
func encodeCredentials(accessKey *iamtypes.AccessKey, region string) string {
	creds := credentials.AWSCredentials{
		Region:          region,
		AccessKeyID:     *accessKey.AccessKeyId,
		SecretAccessKey: *accessKey.SecretAccessKey,
	}
	encCreds, err := creds.RenderBase64EncodedAWSDefaultProfile()
	Expect(err).NotTo(HaveOccurred())
	return encCreds
}

// newUserAccessKey generates a new AWS Access Key pair based off of the
// bootstrap user. This tests that the CloudFormation policy is correct.
func newUserAccessKey(ctx context.Context, cfg *aws.Config, userName string) *iamtypes.AccessKey {
	iamSvc := iam.NewFromConfig(*cfg)

	keyOuts, _ := iamSvc.ListAccessKeys(ctx, &iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	})
	for i := range keyOuts.AccessKeyMetadata {
		By(fmt.Sprintf("Deleting an existing access key: user-name=%s", userName))
		_, err := iamSvc.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
			UserName:    aws.String(userName),
			AccessKeyId: keyOuts.AccessKeyMetadata[i].AccessKeyId,
		})
		Expect(err).NotTo(HaveOccurred())
	}
	By(fmt.Sprintf("Creating an access key: user-name=%s", userName))
	out, err := iamSvc.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{UserName: aws.String(userName)})
	Expect(err).NotTo(HaveOccurred())
	Expect(out.AccessKey).ToNot(BeNil())
	Expect(out.AccessKey.AccessKeyId).ToNot(BeNil())
	Expect(out.AccessKey.SecretAccessKey).ToNot(BeNil())

	return &iamtypes.AccessKey{
		AccessKeyId:     out.AccessKey.AccessKeyId,
		SecretAccessKey: out.AccessKey.SecretAccessKey,
	}
}

func DumpCloudTrailEvents(e2eCtx *E2EContext) {
	if e2eCtx.BootstrapUserAWSSession == nil {
		Fail("Couldn't dump cloudtrail events: no AWS client was set up (please look at previous errors)")
		return
	}

	client := cloudtrail.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)
	events := []cloudtrailtypes.Event{}

	paginator := cloudtrail.NewLookupEventsPaginator(client, &cloudtrail.LookupEventsInput{
		StartTime: aws.Time(e2eCtx.StartOfSuite),
		EndTime:   aws.Time(time.Now()),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "Couldn't get AWS CloudTrail events: err=%v\n", err)
			break
		}
		events = append(events, page.Events...)
	}

	logPath := filepath.Join(e2eCtx.Settings.ArtifactFolder, "cloudtrail-events.yaml")
	dat, err := yaml.Marshal(events)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed to marshal AWS CloudTrail events: err=%v\n", err)
	}
	if err := os.WriteFile(logPath, dat, 0o600); err != nil {
		fmt.Fprintf(GinkgoWriter, "Couldn't write cloudtrail events to file: file=%q err=%s\n", logPath, err)
		return
	}
}

// conformanceImageID looks up a specific image for a given
// Kubernetes version in the e2econfig.
func conformanceImageID(e2eCtx *E2EContext) string {
	ver := e2eCtx.E2EConfig.MustGetVariable("CONFORMANCE_CI_ARTIFACTS_KUBERNETES_VERSION")
	amiName := AMIPrefix + ver + "*"

	By(fmt.Sprintf("Searching for AMI: name=%s", amiName))
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)
	filters := []ec2types.Filter{
		{
			Name:   aws.String("name"),
			Values: []string{amiName},
		},
	}
	filters = append(filters, ec2types.Filter{
		Name:   aws.String("owner-id"),
		Values: []string{DefaultImageLookupOrg},
	})
	resp, err := ec2Svc.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{
		Filters: filters,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(len(resp.Images)).To(Not(BeZero()))
	imageID := aws.ToString(resp.Images[0].ImageId)
	By(fmt.Sprintf("Using AMI: image-id=%s", imageID))
	return imageID
}

func GetAvailabilityZones(config aws.Config) []ec2types.AvailabilityZone {
	ec2Client := ec2.NewFromConfig(config)
	azs, err := ec2Client.DescribeAvailabilityZones(context.TODO(), nil)
	Expect(err).NotTo(HaveOccurred())
	return azs.AvailabilityZones
}

type ServiceQuota struct {
	ServiceCode         string
	QuotaName           string
	QuotaCode           string
	Value               int
	DesiredMinimumValue int
	RequestStatus       string
}

func EnsureServiceQuotas(sess *aws.Config) (map[string]*ServiceQuota, map[string]*servicequotastypes.ServiceQuota) {
	ctx := context.TODO()
	limitedResources := getLimitedResources()
	serviceQuotasClient := servicequotas.NewFromConfig(*sess)

	originalQuotas := map[string]*servicequotastypes.ServiceQuota{}

	for k, v := range limitedResources {
		out, err := serviceQuotasClient.GetServiceQuota(ctx, &servicequotas.GetServiceQuotaInput{
			QuotaCode:   aws.String(v.QuotaCode),
			ServiceCode: aws.String(v.ServiceCode),
		})
		Expect(err).NotTo(HaveOccurred())
		originalQuotas[k] = out.Quota
		v.Value = int(aws.ToFloat64(out.Quota.Value))
		limitedResources[k] = v
		if v.Value < v.DesiredMinimumValue {
			v.attemptRaiseServiceQuotaRequest(ctx, serviceQuotasClient)
		}
	}

	return limitedResources, originalQuotas
}

func (s *ServiceQuota) attemptRaiseServiceQuotaRequest(ctx context.Context, serviceQuotasClient *servicequotas.Client) {
	s.updateServiceQuotaRequestStatus(ctx, serviceQuotasClient)
	if s.RequestStatus == "" {
		s.raiseServiceRequest(ctx, serviceQuotasClient)
	}
}

func (s *ServiceQuota) raiseServiceRequest(ctx context.Context, serviceQuotasClient *servicequotas.Client) {
	fmt.Printf("Requesting service quota increase for %s/%s to %d\n", s.ServiceCode, s.QuotaName, s.DesiredMinimumValue)
	out, err := serviceQuotasClient.RequestServiceQuotaIncrease(
		ctx,
		&servicequotas.RequestServiceQuotaIncreaseInput{
			DesiredValue: aws.Float64(float64(s.DesiredMinimumValue)),
			ServiceCode:  aws.String(s.ServiceCode),
			QuotaCode:    aws.String(s.QuotaCode),
		},
	)
	if err != nil {
		fmt.Printf("Unable to raise quota for %s/%s: %s\n", s.ServiceCode, s.QuotaName, err)
	} else {
		s.RequestStatus = string(out.RequestedQuota.Status)
	}
}

func (s *ServiceQuota) updateServiceQuotaRequestStatus(ctx context.Context, serviceQuotasClient *servicequotas.Client) {
	params := &servicequotas.ListRequestedServiceQuotaChangeHistoryInput{
		ServiceCode: aws.String(s.ServiceCode),
	}
	latestRequest := &servicequotastypes.RequestedServiceQuotaChange{}

	paginator := servicequotas.NewListRequestedServiceQuotaChangeHistoryPaginator(serviceQuotasClient, params)
	for paginator.HasMorePages() {
		page, _ := paginator.NextPage(ctx)
		for _, v := range page.RequestedQuotas {
			if int(aws.ToFloat64(v.DesiredValue)) >= s.DesiredMinimumValue && aws.ToString(v.QuotaCode) == s.QuotaCode && aws.ToTime(v.Created).After(aws.ToTime(latestRequest.Created)) {
				latestRequest = &v
			}
		}
	}

	if latestRequest.Status != "" {
		s.RequestStatus = string(latestRequest.Status)
	}
}

// DumpEKSClusters dumps the EKS clusters in the environment.
func DumpEKSClusters(ctx context.Context, e2eCtx *E2EContext) {
	name := "no-bootstrap-cluster"
	if e2eCtx.Environment.BootstrapClusterProxy != nil {
		name = e2eCtx.Environment.BootstrapClusterProxy.GetName()
	}
	logPath := filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", name, "aws-resources")
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		fmt.Fprintf(GinkgoWriter, "Couldn't create directory: path=%q, err=%s\n", logPath, err)
	}
	fmt.Fprintf(GinkgoWriter, "Folder created for eks clusters: %q\n", logPath)

	input := &eks.ListClustersInput{}
	var eksClient *eks.Client
	if e2eCtx.BootstrapUserAWSSession == nil && e2eCtx.AWSSession != nil {
		eksClient = eks.NewFromConfig(*e2eCtx.AWSSession)
	} else if e2eCtx.BootstrapUserAWSSession != nil {
		eksClient = eks.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)
	} else {
		Fail("Couldn't list EKS clusters: no AWS client was set up (please look at previous errors)")
		return
	}

	output, err := eksClient.ListClusters(ctx, input)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Couldn't list EKS clusters: err=%s\n", err)
		return
	}

	for _, clusterName := range output.Clusters {
		describeInput := &eks.DescribeClusterInput{
			Name: aws.String(clusterName),
		}
		describeOutput, err := eksClient.DescribeCluster(ctx, describeInput)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "Couldn't describe EKS clusters: name=%q err=%s\n", clusterName, err)
			continue
		}
		dumpEKSCluster(describeOutput.Cluster, logPath)
	}
}

func dumpEKSCluster(cluster *ekstypes.Cluster, logPath string) {
	clusterYAML, err := yaml.Marshal(cluster)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Couldn't marshal cluster to yaml: name=%q err=%s\n", *cluster.Name, err)
		return
	}

	fileName := fmt.Sprintf("%s.yaml", *cluster.Name)
	clusterLog := path.Join(logPath, fileName)
	f, err := os.OpenFile(clusterLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm) //nolint:gosec
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Couldn't open log file: name=%q err=%s\n", clusterLog, err)
		return
	}
	defer f.Close()

	if err := os.WriteFile(f.Name(), clusterYAML, 0o600); err != nil {
		fmt.Fprintf(GinkgoWriter, "Couldn't write cluster yaml to file: name=%q file=%q err=%s\n", *cluster.Name, f.Name(), err)
		return
	}
}

// To calculate how much resources a test consumes, these helper functions below can be used.
// ListVpcInternetGateways, ListNATGateways, ListRunningEC2, ListVPC.

func ListVpcInternetGateways(e2eCtx *E2EContext) ([]ec2types.InternetGateway, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	out, err := ec2Svc.DescribeInternetGateways(context.TODO(), &ec2.DescribeInternetGatewaysInput{})
	if err != nil {
		return nil, err
	}

	return out.InternetGateways, nil
}

func ListNATGateways(e2eCtx *E2EContext) (map[string]ec2types.NatGateway, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	describeNatGatewayInput := &ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			filter.EC2.NATGatewayStates(ec2types.NatGatewayStateAvailable),
		},
	}

	gateways := make(map[string]ec2types.NatGateway)

	paginator := ec2.NewDescribeNatGatewaysPaginator(ec2Svc, describeNatGatewayInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to describe NAT gateways: %w", err)
		}
		for _, r := range page.NatGateways {
			gateways[aws.ToString(r.SubnetId)] = r
		}
	}

	return gateways, nil
}

// listRunningEC2 returns a list of running EC2 instances.
func listRunningEC2(e2eCtx *E2EContext) ([]instance, error) { //nolint:unused
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	resp, err := ec2Svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{filter.EC2.InstanceStates(ec2types.InstanceStateNameRunning)},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("no machines found")
	}
	instances := []instance{}
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			tags := i.Tags
			name := ""
			for _, t := range tags {
				if aws.ToString(t.Key) == "Name" {
					name = aws.ToString(t.Value)
				}
			}
			if name == "" {
				continue
			}
			instances = append(instances,
				instance{
					name:       name,
					instanceID: aws.ToString(i.InstanceId),
				},
			)
		}
	}
	return instances, nil
}

func ListClusterEC2Instances(e2eCtx *E2EContext, clusterName string) ([]ec2types.Instance, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)
	filter := ec2types.Filter{
		Name:   aws.String("tag-key"),
		Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/" + clusterName},
	}
	input := &ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	instances := []ec2types.Instance{}
	for _, r := range result.Reservations {
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}

func WaitForInstanceState(e2eCtx *E2EContext, clusterName string, state string) bool {
	Eventually(func(gomega Gomega) bool {
		st := map[string]int{
			"pending":       0,
			"running":       0,
			"shutting-down": 0,
			"terminated":    0,
		}
		instances, _ := ListClusterEC2Instances(e2eCtx, clusterName)
		for _, i := range instances {
			iState := i.State.Name
			st[string(iState)]++
		}
		if st[state] == len(instances) || len(instances) == 0 {
			return true
		}
		return false
	}, 5*time.Minute, 5*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed waiting for all cluster's EC2 instance to be in %q state", state))

	return false
}

func TerminateInstance(e2eCtx *E2EContext, instanceID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	if _, err := ec2Svc.TerminateInstances(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func ListVPC(e2eCtx *E2EContext) int {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			filter.EC2.VPCStates(ec2types.VpcStateAvailable),
		},
	}

	out, err := ec2Svc.DescribeVpcs(context.TODO(), input)
	if err != nil {
		return 0
	}

	return len(out.Vpcs)
}

func GetVPC(e2eCtx *E2EContext, vpcID string) (*ec2types.Vpc, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("vpc-id"),
		Values: []string{vpcID},
	}

	input := &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeVpcs(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.Vpcs == nil {
		return nil, nil
	}
	return &result.Vpcs[0], nil
}

func GetVPCByName(e2eCtx *E2EContext, vpcName string) (*ec2types.Vpc, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("tag:Name"),
		Values: []string{vpcName},
	}

	input := &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeVpcs(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if len(result.Vpcs) == 0 {
		return nil, awserrors.NewNotFound("Vpc not found")
	}
	return &result.Vpcs[0], nil
}

func GetVPCEndpointsByID(e2eCtx *E2EContext, vpcID string) ([]ec2types.VpcEndpoint, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DescribeVpcEndpointsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	}

	res := []ec2types.VpcEndpoint{}
	paginator := ec2.NewDescribeVpcEndpointsPaginator(ec2Svc, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to describe VPC endpoints: %w", err)
		}
		res = append(res, page.VpcEndpoints...)
	}

	return res, nil
}

func CreateVPC(e2eCtx *E2EContext, vpcName string, cidrBlock string) (*ec2types.Vpc, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidrBlock),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeVpc,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(vpcName),
					},
				},
			},
		},
	}
	result, err := ec2Svc.CreateVpc(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.Vpc, nil
}

func DisassociateVpcCidrBlock(e2eCtx *E2EContext, assocID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DisassociateVpcCidrBlockInput{
		AssociationId: aws.String(assocID),
	}

	if _, err := ec2Svc.DisassociateVpcCidrBlock(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func DeleteVPC(e2eCtx *E2EContext, vpcID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcID),
	}
	if _, err := ec2Svc.DeleteVpc(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func ListVpcSubnets(e2eCtx *E2EContext, vpcID string) ([]ec2types.Subnet, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("vpc-id"),
		Values: []string{vpcID},
	}

	input := &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeSubnets(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.Subnets, nil
}

func GetSubnet(e2eCtx *E2EContext, subnetID string) (*ec2types.Subnet, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("subnet-id"),
		Values: []string{subnetID},
	}

	input := &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeSubnets(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.Subnets == nil {
		return nil, nil
	}
	return &result.Subnets[0], nil
}

func GetSubnetByName(e2eCtx *E2EContext, name string) (*ec2types.Subnet, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("tag:Name"),
		Values: []string{name},
	}

	input := &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeSubnets(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.Subnets == nil {
		return nil, nil
	}
	return &result.Subnets[0], nil
}

func CreateSubnet(e2eCtx *E2EContext, clusterName string, cidrBlock string, az string, vpcID string, st string) (*ec2types.Subnet, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(cidrBlock),
		VpcId:     aws.String(vpcID),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSubnet,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(clusterName + "-subnet-" + st),
					},
				},
			},
		},
	}

	if az != "" {
		input.AvailabilityZone = aws.String(az)
	}

	result, err := ec2Svc.CreateSubnet(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.Subnet, nil
}

func DeleteSubnet(e2eCtx *E2EContext, subnetID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetID),
	}

	if _, err := ec2Svc.DeleteSubnet(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func GetAddress(e2eCtx *E2EContext, allocationID string) (*ec2types.Address, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("allocation-id"),
		Values: []string{allocationID},
	}

	input := &ec2.DescribeAddressesInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeAddresses(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.Addresses == nil {
		return nil, nil
	}
	return &result.Addresses[0], nil
}

func AllocateAddress(e2eCtx *E2EContext, eipName string) (*ec2.AllocateAddressOutput, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.AllocateAddressInput{
		Domain: ec2types.DomainTypeVpc,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeElasticIp,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(eipName),
					},
				},
			},
		},
	}

	result, err := ec2Svc.AllocateAddress(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DisassociateAddress(e2eCtx *E2EContext, assocID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DisassociateAddressInput{
		AssociationId: aws.String(assocID),
	}

	if _, err := ec2Svc.DisassociateAddress(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func ReleaseAddress(e2eCtx *E2EContext, allocationID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.ReleaseAddressInput{
		AllocationId: aws.String(allocationID),
	}

	if _, err := ec2Svc.ReleaseAddress(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func CreateNatGateway(e2eCtx *E2EContext, gatewayName string, connectType string, allocationID string, subnetID string) (*ec2types.NatGateway, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateNatGatewayInput{
		SubnetId: aws.String(subnetID),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeNatgateway,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(gatewayName),
					},
				},
			},
		},
	}

	if connectType != "" {
		input.ConnectivityType = ec2types.ConnectivityType(connectType)
	}

	if allocationID != "" {
		input.AllocationId = aws.String(allocationID)
	}

	result, err := ec2Svc.CreateNatGateway(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.NatGateway, nil
}

func GetNatGateway(e2eCtx *E2EContext, gatewayID string) (*ec2types.NatGateway, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("nat-gateway-id"),
		Values: []string{gatewayID},
	}

	input := &ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeNatGateways(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.NatGateways == nil {
		return nil, nil
	}
	return &result.NatGateways[0], nil
}

func DeleteNatGateway(e2eCtx *E2EContext, gatewayID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteNatGatewayInput{
		NatGatewayId: aws.String(gatewayID),
	}

	if _, err := ec2Svc.DeleteNatGateway(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func WaitForNatGatewayState(e2eCtx *E2EContext, gatewayID string, state string) bool {
	Eventually(func(gomega Gomega) bool {
		gw, _ := GetNatGateway(e2eCtx, gatewayID)
		gwState := gw.State
		return gwState == ec2types.NatGatewayState(state)
	}, 3*time.Minute, 5*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed waiting for NAT Gateway to be in %q state", state))
	return false
}

func CreateInternetGateway(e2eCtx *E2EContext, gatewayName string) (*ec2types.InternetGateway, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateInternetGatewayInput{
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeInternetGateway,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(gatewayName),
					},
				},
			},
		},
	}

	result, err := ec2Svc.CreateInternetGateway(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.InternetGateway, nil
}

func GetInternetGateway(e2eCtx *E2EContext, gatewayID string) (*ec2types.InternetGateway, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("internet-gateway-id"),
		Values: []string{gatewayID},
	}

	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeInternetGateways(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.InternetGateways == nil {
		return nil, nil
	}
	return &result.InternetGateways[0], nil
}

func DeleteInternetGateway(e2eCtx *E2EContext, gatewayID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(gatewayID),
	}

	if _, err := ec2Svc.DeleteInternetGateway(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func AttachInternetGateway(e2eCtx *E2EContext, gatewayID string, vpcID string) (bool, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(gatewayID),
		VpcId:             aws.String(vpcID),
	}

	if _, err := ec2Svc.AttachInternetGateway(context.TODO(), input); err != nil {
		return false, err
	}
	return true, nil
}

func DetachInternetGateway(e2eCtx *E2EContext, gatewayID string, vpcID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(gatewayID),
		VpcId:             aws.String(vpcID),
	}

	if _, err := ec2Svc.DetachInternetGateway(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func CreatePeering(e2eCtx *E2EContext, peerName string, vpcID string, peerVpcID string) (*ec2types.VpcPeeringConnection, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateVpcPeeringConnectionInput{
		VpcId:     aws.String(vpcID),
		PeerVpcId: aws.String(peerVpcID),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeVpcPeeringConnection,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(peerName),
					},
				},
			},
		},
	}

	result, err := ec2Svc.CreateVpcPeeringConnection(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.VpcPeeringConnection, nil
}

func GetPeering(e2eCtx *E2EContext, peeringID string) (*ec2types.VpcPeeringConnection, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("vpc-peering-connection-id"),
		Values: []string{peeringID},
	}

	input := &ec2.DescribeVpcPeeringConnectionsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeVpcPeeringConnections(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.VpcPeeringConnections == nil {
		return nil, nil
	}
	return &result.VpcPeeringConnections[0], nil
}

func DeletePeering(e2eCtx *E2EContext, peeringID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(peeringID),
	}

	if _, err := ec2Svc.DeleteVpcPeeringConnection(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func AcceptPeering(e2eCtx *E2EContext, peeringID string) (*ec2types.VpcPeeringConnection, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.AcceptVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(peeringID),
	}

	result, err := ec2Svc.AcceptVpcPeeringConnection(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.VpcPeeringConnection, nil
}

func CreateRouteTable(e2eCtx *E2EContext, rtName string, vpcID string) (*ec2types.RouteTable, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcID),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeRouteTable,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(rtName),
					},
				},
			},
		},
	}

	result, err := ec2Svc.CreateRouteTable(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.RouteTable, nil
}

func ListVpcRouteTables(e2eCtx *E2EContext, vpcID string) ([]ec2types.RouteTable, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("vpc-id"),
		Values: []string{vpcID},
	}

	input := &ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeRouteTables(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.RouteTables, nil
}

func ListSubnetRouteTables(e2eCtx *E2EContext, subnetID string) ([]ec2types.RouteTable, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("association.subnet-id"),
		Values: []string{subnetID},
	}

	input := &ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeRouteTables(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.RouteTables, nil
}

func GetRouteTable(e2eCtx *E2EContext, rtID string) (*ec2types.RouteTable, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("route-table-id"),
		Values: []string{rtID},
	}

	input := &ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeRouteTables(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.RouteTables == nil {
		return nil, nil
	}
	return &result.RouteTables[0], nil
}

func DeleteRouteTable(e2eCtx *E2EContext, rtID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteRouteTableInput{
		RouteTableId: aws.String(rtID),
	}

	if _, err := ec2Svc.DeleteRouteTable(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func CreateRoute(e2eCtx *E2EContext, rtID string, destinationCidr string, natID *string, igwID *string, pcxID *string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateRouteInput{
		RouteTableId:         &rtID,
		DestinationCidrBlock: aws.String(destinationCidr),
	}

	if natID != nil {
		input.NatGatewayId = natID
	}

	if igwID != nil {
		input.GatewayId = igwID
	}

	if pcxID != nil {
		input.VpcPeeringConnectionId = pcxID
	}

	_, err := ec2Svc.CreateRoute(context.TODO(), input)
	return err == nil
}

func DeleteRoute(e2eCtx *E2EContext, rtID string, destinationCidr string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteRouteInput{
		RouteTableId:         aws.String(rtID),
		DestinationCidrBlock: aws.String(destinationCidr),
	}

	if _, err := ec2Svc.DeleteRoute(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func AssociateRouteTable(e2eCtx *E2EContext, rtID string, subnetID string) (*ec2.AssociateRouteTableOutput, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rtID),
		SubnetId:     aws.String(subnetID),
	}

	result, err := ec2Svc.AssociateRouteTable(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DisassociateRouteTable(e2eCtx *E2EContext, assocID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DisassociateRouteTableInput{
		AssociationId: aws.String(assocID),
	}

	if _, err := ec2Svc.DisassociateRouteTable(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func CreateSecurityGroup(e2eCtx *E2EContext, sgName string, sgDescription string, vpcID string) (*ec2.CreateSecurityGroupOutput, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(vpcID),
		GroupName:   aws.String(sgName),
		Description: aws.String(sgDescription),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(sgName),
					},
				},
			},
		},
	}

	result, err := ec2Svc.CreateSecurityGroup(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetSecurityGroupByFilters(e2eCtx *E2EContext, filters []ec2types.Filter) ([]ec2types.SecurityGroup, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: filters,
	}
	result, err := ec2Svc.DescribeSecurityGroups(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.SecurityGroups, nil
}

func GetSecurityGroup(e2eCtx *E2EContext, sgID string) (*ec2types.SecurityGroup, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("group-id"),
		Values: []string{sgID},
	}

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeSecurityGroups(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.SecurityGroups == nil {
		return nil, nil
	}
	return &result.SecurityGroups[0], nil
}

func GetSecurityGroupsByVPC(e2eCtx *E2EContext, vpcID string) ([]ec2types.SecurityGroup, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name: aws.String("vpc-id"),
		Values: []string{
			vpcID,
		},
	}

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeSecurityGroups(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.SecurityGroups == nil {
		return nil, nil
	}
	return result.SecurityGroups, nil
}

func DeleteSecurityGroup(e2eCtx *E2EContext, sgID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sgID),
	}

	if _, err := ec2Svc.DeleteSecurityGroup(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func ListSecurityGroupRules(e2eCtx *E2EContext, sgID string) ([]ec2types.SecurityGroupRule, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("group-id"),
		Values: []string{sgID},
	}

	input := &ec2.DescribeSecurityGroupRulesInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeSecurityGroupRules(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.SecurityGroupRules, nil
}

func GetSecurityGroupRule(e2eCtx *E2EContext, sgrID string) (*ec2types.SecurityGroupRule, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	filter := ec2types.Filter{
		Name:   aws.String("security-group-rule-id"),
		Values: []string{sgrID},
	}

	input := &ec2.DescribeSecurityGroupRulesInput{
		Filters: []ec2types.Filter{
			filter,
		},
	}

	result, err := ec2Svc.DescribeSecurityGroupRules(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if result.SecurityGroupRules == nil {
		return nil, nil
	}
	return &result.SecurityGroupRules[0], nil
}

func CreateSecurityGroupIngressRule(e2eCtx *E2EContext, sgID string, sgrDescription string, cidr string, protocol string, fromPort int32, toPort int32) (bool, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	ipPerm := ec2types.IpPermission{
		FromPort:   aws.Int32(fromPort),
		ToPort:     aws.Int32(toPort),
		IpProtocol: aws.String(protocol),
		IpRanges: []ec2types.IpRange{
			{
				CidrIp:      aws.String(cidr),
				Description: aws.String(sgrDescription),
			},
		},
	}

	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(sgID),
		IpPermissions: []ec2types.IpPermission{
			ipPerm,
		},
	}

	result, err := ec2Svc.AuthorizeSecurityGroupIngress(context.TODO(), input)
	if err != nil {
		return false, err
	}
	return *result.Return, nil
}

func CreateSecurityGroupEgressRule(e2eCtx *E2EContext, sgID string, sgrDescription string, cidr string, protocol string, fromPort int32, toPort int32) (bool, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	ipPerm := ec2types.IpPermission{
		FromPort:   aws.Int32(fromPort),
		ToPort:     aws.Int32(toPort),
		IpProtocol: aws.String(protocol),
		IpRanges: []ec2types.IpRange{
			{
				CidrIp:      aws.String(cidr),
				Description: aws.String(sgrDescription),
			},
		},
	}

	input := &ec2.AuthorizeSecurityGroupEgressInput{
		GroupId: aws.String(sgID),
		IpPermissions: []ec2types.IpPermission{
			ipPerm,
		},
	}
	result, err := ec2Svc.AuthorizeSecurityGroupEgress(context.TODO(), input)
	if err != nil {
		return false, err
	}
	return *result.Return, nil
}

func CreateSecurityGroupRule(e2eCtx *E2EContext, sgID string, sgrDescription string, cidr string, protocol string, fromPort int32, toPort int32, rt string) (bool, error) {
	switch rt {
	case "ingress":
		return CreateSecurityGroupIngressRule(e2eCtx, sgID, sgrDescription, cidr, protocol, fromPort, toPort)
	case "egress":
		return CreateSecurityGroupEgressRule(e2eCtx, sgID, sgrDescription, cidr, protocol, fromPort, toPort)
	}
	return false, nil
}

func CreateSecurityGroupIngressRuleWithSourceSG(e2eCtx *E2EContext, sgID string, protocol string, toPort int32, sourceSecurityGroupID string) (bool, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	ipPerm := ec2types.IpPermission{
		FromPort:   aws.Int32(toPort),
		ToPort:     aws.Int32(toPort),
		IpProtocol: aws.String(protocol),
		UserIdGroupPairs: []ec2types.UserIdGroupPair{
			{
				GroupId: aws.String(sourceSecurityGroupID),
			},
		},
	}
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(sgID),
		IpPermissions: []ec2types.IpPermission{
			ipPerm,
		},
	}

	result, err := ec2Svc.AuthorizeSecurityGroupIngress(context.TODO(), input)
	if err != nil {
		return false, err
	}
	return *result.Return, nil
}

func DeleteSecurityGroupIngressRule(e2eCtx *E2EContext, sgID, sgrID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.RevokeSecurityGroupIngressInput{
		SecurityGroupRuleIds: []string{sgrID},
		GroupId:              aws.String(sgID),
	}

	if _, err := ec2Svc.RevokeSecurityGroupIngress(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func DeleteSecurityGroupEgressRule(e2eCtx *E2EContext, sgID, sgrID string) bool {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.RevokeSecurityGroupEgressInput{
		SecurityGroupRuleIds: []string{sgrID},
		GroupId:              aws.String(sgID),
	}

	if _, err := ec2Svc.RevokeSecurityGroupEgress(context.TODO(), input); err != nil {
		return false
	}
	return true
}

func DeleteSecurityGroupRule(e2eCtx *E2EContext, sgID, sgrID, rt string) bool {
	switch rt {
	case "ingress":
		return DeleteSecurityGroupIngressRule(e2eCtx, sgID, sgrID)
	case "egress":
		return DeleteSecurityGroupEgressRule(e2eCtx, sgID, sgrID)
	}
	return false
}

func ListLoadBalancers(ctx context.Context, e2eCtx *E2EContext, clusterName string) ([]elbtypes.LoadBalancerDescription, error) {
	elbSvc := elb.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []string{clusterName + "-apiserver"},
	}

	result, err := elbSvc.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.LoadBalancerDescriptions == nil {
		return nil, nil
	}
	return result.LoadBalancerDescriptions, nil
}

func DeleteLoadBalancer(ctx context.Context, e2eCtx *E2EContext, loadbalancerName string) bool {
	elbSvc := elb.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(loadbalancerName),
	}

	if _, err := elbSvc.DeleteLoadBalancer(ctx, input); err != nil {
		return false
	}
	return true
}

func CreateEFS(e2eCtx *E2EContext, creationToken string) (*efs.CreateFileSystemOutput, error) {
	efsSvc := efs.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	input := &efs.CreateFileSystemInput{
		CreationToken: aws.String(creationToken),
		Encrypted:     aws.Bool(true),
	}
	efsOutput, err := efsSvc.CreateFileSystem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return efsOutput, nil
}

func DescribeEFS(ctx context.Context, e2eCtx *E2EContext, efsID string) (*efstypes.FileSystemDescription, error) {
	efsSvc := efs.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	input := &efs.DescribeFileSystemsInput{
		FileSystemId: aws.String(efsID),
	}
	efsOutput, err := efsSvc.DescribeFileSystems(ctx, input)
	if err != nil {
		return nil, err
	}
	if efsOutput == nil || len(efsOutput.FileSystems) == 0 {
		return nil, &efstypes.FileSystemNotFound{}
	}
	return &efsOutput.FileSystems[0], nil
}

func DeleteEFS(ctx context.Context, e2eCtx *E2EContext, efsID string) (*efs.DeleteFileSystemOutput, error) {
	efsSvc := efs.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	input := &efs.DeleteFileSystemInput{
		FileSystemId: aws.String(efsID),
	}
	result, err := efsSvc.DeleteFileSystem(ctx, input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetEFSState(ctx context.Context, e2eCtx *E2EContext, efsID string) (*string, error) {
	efs, err := DescribeEFS(ctx, e2eCtx, efsID)
	if err != nil {
		return nil, err
	}
	state := string(efs.LifeCycleState)
	return &state, nil
}

func CreateMountTargetOnEFS(ctx context.Context, e2eCtx *E2EContext, efsID string, vpcID string, sg string) (*efs.CreateMountTargetOutput, error) {
	efsSvc := efs.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	subnets, err := ListVpcSubnets(e2eCtx, vpcID)
	if err != nil {
		return nil, err
	}
	input := &efs.CreateMountTargetInput{
		FileSystemId:   aws.String(efsID),
		SecurityGroups: []string{sg},
		SubnetId:       subnets[0].SubnetId,
	}
	result, err := efsSvc.CreateMountTarget(ctx, input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteMountTarget(ctx context.Context, e2eCtx *E2EContext, mountTargetID string) (*efs.DeleteMountTargetOutput, error) {
	efsSvc := efs.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	input := &efs.DeleteMountTargetInput{
		MountTargetId: aws.String(mountTargetID),
	}
	result, err := efsSvc.DeleteMountTarget(ctx, input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetMountTarget(ctx context.Context, e2eCtx *E2EContext, mountTargetID string) (*efstypes.MountTargetDescription, error) {
	efsSvc := efs.NewFromConfig(*e2eCtx.BootstrapUserAWSSession)

	input := &efs.DescribeMountTargetsInput{
		MountTargetId: aws.String(mountTargetID),
	}
	result, err := efsSvc.DescribeMountTargets(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(result.MountTargets) == 0 {
		return nil, &efstypes.MountTargetNotFound{}
	}
	return &result.MountTargets[0], nil
}

func GetMountTargetState(ctx context.Context, e2eCtx *E2EContext, mountTargetID string) (*string, error) {
	result, err := GetMountTarget(ctx, e2eCtx, mountTargetID)
	if err != nil {
		return nil, err
	}

	state := string(result.LifeCycleState)
	return &state, nil
}

func getAvailabilityZone(e2eCtx *E2EContext) string {
	az := e2eCtx.E2EConfig.MustGetVariable(AwsAvailabilityZone1)
	return az
}

func getInstanceFamily(e2eCtx *E2EContext) string {
	machineType := e2eCtx.E2EConfig.MustGetVariable(AwsNodeMachineType)
	// from instance type get instace family behind the dot
	// for example: t3a.medium -> t3
	machineTypeSplit := strings.Split(machineType, ".")
	if len(machineTypeSplit) > 0 {
		return machineTypeSplit[0]
	}
	return "t3"
}

func AllocateHost(ctx context.Context, e2eCtx *E2EContext) (string, error) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)
	input := &ec2.AllocateHostsInput{
		AvailabilityZone: aws.String(getAvailabilityZone(e2eCtx)),
		InstanceFamily:   aws.String(getInstanceFamily(e2eCtx)),
		Quantity:         aws.Int32(1),
	}
	output, err := ec2Svc.AllocateHosts(ctx, input)
	Expect(err).ToNot(HaveOccurred(), "Failed to allocate  host")
	Expect(len(output.HostIds)).To(BeNumerically(">", 0), "No dedicated host ID returned")
	fmt.Println("Allocated Host ID: ", output.HostIds[0])
	hostID := output.HostIds[0]
	return hostID, nil
}

func ReleaseHost(ctx context.Context, e2eCtx *E2EContext, hostID string) {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.ReleaseHostsInput{
		HostIds: []string{hostID},
	}

	_, err := ec2Svc.ReleaseHosts(ctx, input)
	Expect(err).ToNot(HaveOccurred(), "Failed to release host %s", hostID)
	fmt.Println("Released Host ID: ", hostID)
}

func GetHostID(ctx context.Context, e2eCtx *E2EContext, instanceID string) string {
	ec2Svc := ec2.NewFromConfig(*e2eCtx.AWSSession)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := ec2Svc.DescribeInstances(ctx, input)
	Expect(err).ToNot(HaveOccurred(), "Failed to get host ID for instance %s", instanceID)
	Expect(len(result.Reservations)).To(BeNumerically(">", 0), "No reservation returned")
	Expect(len(result.Reservations[0].Instances)).To(BeNumerically(">", 0), "No instance returned")
	placement := *result.Reservations[0].Instances[0].Placement
	hostID := *placement.HostId
	fmt.Println("Host ID: ", hostID)
	return hostID
}

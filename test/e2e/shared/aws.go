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
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/service"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/yaml"
)

func NewAWSSession() client.ConfigProvider {
	By("Getting an AWS IAM session - from environment")
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true).WithRegion(region)
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	})
	Expect(err).NotTo(HaveOccurred())
	_, err = sess.Config.Credentials.Get()
	Expect(err).NotTo(HaveOccurred())
	return sess
}

func NewAWSSessionWithKey(accessKey *iam.AccessKey) client.ConfigProvider {
	By("Getting an AWS IAM session - from access key")
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true).WithRegion(region)
	config.Credentials = awscreds.NewStaticCredentials(*accessKey.AccessKeyId, *accessKey.SecretAccessKey, "")

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: *config,
	})
	Expect(err).NotTo(HaveOccurred())
	_, err = sess.Config.Credentials.Get()
	Expect(err).NotTo(HaveOccurred())
	return sess
}

// createCloudFormationStack ensures the cloudformation stack is up to date
func createCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) error {
	Byf("Creating AWS CloudFormation stack for AWS IAM resources: stack-name=%s", t.Spec.StackName)
	CFN := cfn.New(prov)
	cfnSvc := cloudformation.NewService(CFN)

	err := cfnSvc.ReconcileBootstrapStack(t.Spec.StackName, *renderCustomCloudFormation(t))
	if err != nil {
		stack, err := CFN.DescribeStacks(&cfn.DescribeStacksInput{StackName: aws.String(t.Spec.StackName)})
		if err == nil && len(stack.Stacks) > 0 {
			deleteMultitenancyRoles(prov)
			if aws.StringValue(stack.Stacks[0].StackStatus) == cfn.StackStatusRollbackFailed ||
				aws.StringValue(stack.Stacks[0].StackStatus) == cfn.StackStatusRollbackComplete ||
				aws.StringValue(stack.Stacks[0].StackStatus) == cfn.StackStatusRollbackInProgress {
				// If cloudformation stack creation fails due to resources that already exist, stack stays in rollback status and must be manually deleted.
				// Delete resources that failed because they already exists.
				deleteResourcesInCloudFormation(prov, t)
			}
		}
	}
	return err
}

func SetMultitenancyEnvVars(prov client.ConfigProvider) error {
	for _, roles := range MultiTenancyRoles {
		if err := roles.SetEnvVars(prov); err != nil {
			return err
		}
	}
	return nil
}

// Delete resources that already exists.
func deleteResourcesInCloudFormation(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	iamSvc := iam.New(prov)
	temp := *renderCustomCloudFormation(t)
	for _, val := range temp.Resources {
		tayp := val.AWSCloudFormationType()
		if tayp == configservice.ResourceTypeAwsIamRole {
			role := val.(*cfn_iam.Role)
			iamSvc.DeleteRole(&iam.DeleteRoleInput{RoleName: aws.String(role.RoleName)})
		}
		if val.AWSCloudFormationType() == "AWS::IAM::InstanceProfile" {
			profile := val.(*cfn_iam.InstanceProfile)
			iamSvc.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{InstanceProfileName: aws.String(profile.InstanceProfileName)})
		}
		if val.AWSCloudFormationType() == "AWS::IAM::ManagedPolicy" {
			policy := val.(*cfn_iam.ManagedPolicy)
			policies, err := iamSvc.ListPolicies(&iam.ListPoliciesInput{})
			Expect(err).NotTo(HaveOccurred())
			if len(policies.Policies) > 0 {
				for _, p := range policies.Policies {
					if aws.StringValue(p.PolicyName) == policy.ManagedPolicyName {
						iamSvc.DeletePolicy(&iam.DeletePolicyInput{PolicyArn: p.Arn})
						break
					}
				}
			}
		}
		if val.AWSCloudFormationType() == configservice.ResourceTypeAwsIamGroup {
			group := val.(*cfn_iam.Group)
			iamSvc.DeleteGroup(&iam.DeleteGroupInput{GroupName: aws.String(group.GroupName)})
		}
	}
}

// TODO: remove once test infra accounts are fixed.
func deleteMultitenancyRoles(prov client.ConfigProvider) {
	DeleteRole(prov, "multi-tenancy-role")
	DeleteRole(prov, "multi-tenancy-nested-role")
}

// detachAllPoliciesForRole detaches all policies for role
func detachAllPoliciesForRole(prov client.ConfigProvider, name string) error {
	iamSvc := iam.New(prov)

	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &name,
	}
	policies, err := iamSvc.ListAttachedRolePolicies(input)
	if err != nil {
		return errors.New("error fetching policies for role")
	}
	for _, p := range policies.AttachedPolicies {
		input := &iam.DetachRolePolicyInput{
			RoleName:  aws.String(name),
			PolicyArn: p.PolicyArn,
		}

		_, err := iamSvc.DetachRolePolicy(input)
		if err != nil {
			return errors.New("failed detaching policy from a role")
		}
	}
	return nil
}

// Best effort deletes roles.
func DeleteRole(prov client.ConfigProvider, name string) {
	iamSvc := iam.New(prov)

	// if role does not exist, return
	_, err := iamSvc.GetRole(&iam.GetRoleInput{RoleName: aws.String(name)})
	if err != nil {
		return
	}

	if err := detachAllPoliciesForRole(prov, name); err != nil {
		return
	}

	iamSvc.DeleteRole(&iam.DeleteRoleInput{RoleName: aws.String(name)})
}

func GetPolicyArn(prov client.ConfigProvider, name string) string {
	iamSvc := iam.New(prov)
	policyList, err := iamSvc.ListPolicies(&iam.ListPoliciesInput{
		Scope: aws.String(iam.PolicyScopeTypeLocal),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, policy := range policyList.Policies {
		if aws.StringValue(policy.PolicyName) == name {
			return aws.StringValue(policy.Arn)
		}
	}
	return ""
}

// deleteCloudFormationStack removes the provisioned clusterawsadm stack
func deleteCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	Byf("Deleting %s CloudFormation stack", t.Spec.StackName)
	CFN := cfn.New(prov)
	cfnSvc := cloudformation.NewService(CFN)
	err := cfnSvc.DeleteStack(t.Spec.StackName, nil)
	if err != nil {
		var retainResources []*string
		out, err := CFN.DescribeStackResources(&cfn.DescribeStackResourcesInput{StackName: aws.String(t.Spec.StackName)})
		Expect(err).NotTo(HaveOccurred())
		for _, v := range out.StackResources {
			if aws.StringValue(v.ResourceStatus) == cfn.ResourceStatusDeleteFailed {
				retainResources = append(retainResources, v.LogicalResourceId)
			}
		}
		err = cfnSvc.DeleteStack(t.Spec.StackName, retainResources)
		Expect(err).NotTo(HaveOccurred())
	}
	CFN.WaitUntilStackDeleteComplete(&cfn.DescribeStacksInput{
		StackName: aws.String(t.Spec.StackName),
	})
}

// ensureNoServiceLinkedRoles removes an auto-created IAM role, and tests
// the controller's IAM permissions to use ELB and Spot instances successfully
func ensureNoServiceLinkedRoles(prov client.ConfigProvider) {
	Byf("Deleting AWS IAM Service Linked Role: role-name=AWSServiceRoleForElasticLoadBalancing")
	iamSvc := iam.New(prov)
	_, err := iamSvc.DeleteServiceLinkedRole(&iam.DeleteServiceLinkedRoleInput{
		RoleName: aws.String("AWSServiceRoleForElasticLoadBalancing"),
	})
	if code, _ := awserrors.Code(err); code != iam.ErrCodeNoSuchEntityException {
		Expect(err).NotTo(HaveOccurred())
	}

	Byf("Deleting AWS IAM Service Linked Role: role-name=AWSServiceRoleForEC2Spot")
	_, err = iamSvc.DeleteServiceLinkedRole(&iam.DeleteServiceLinkedRoleInput{
		RoleName: aws.String("AWSServiceRoleForEC2Spot"),
	})
	if code, _ := awserrors.Code(err); code != iam.ErrCodeNoSuchEntityException {
		Expect(err).NotTo(HaveOccurred())
	}
}

// ensureSSHKeyPair ensures A SSH key is present under the name
func ensureSSHKeyPair(prov client.ConfigProvider, keyPairName string) {
	Byf("Ensuring presence of SSH key in EC2: key-name=%s", keyPairName)
	ec2c := ec2.New(prov)
	_, err := ec2c.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	if code, _ := awserrors.Code(err); code != "InvalidKeyPair.Duplicate" {
		Expect(err).NotTo(HaveOccurred())
	}
}

// encodeCredentials leverages clusterawsadm to encode AWS credentials
func encodeCredentials(accessKey *iam.AccessKey, region string) string {
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
func newUserAccessKey(prov client.ConfigProvider, userName string) *iam.AccessKey {
	iamSvc := iam.New(prov)
	keyOuts, _ := iamSvc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	})
	for i := range keyOuts.AccessKeyMetadata {
		Byf("Deleting an existing access key: user-name=%s", userName)
		_, err := iamSvc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			UserName:    aws.String(userName),
			AccessKeyId: keyOuts.AccessKeyMetadata[i].AccessKeyId,
		})
		Expect(err).NotTo(HaveOccurred())
	}
	Byf("Creating an access key: user-name=%s", userName)
	out, err := iamSvc.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: aws.String(userName)})
	Expect(err).NotTo(HaveOccurred())
	Expect(out.AccessKey).ToNot(BeNil())

	return &iam.AccessKey{
		AccessKeyId:     out.AccessKey.AccessKeyId,
		SecretAccessKey: out.AccessKey.SecretAccessKey,
	}
}

// conformanceImageID looks up a specific image for a given
// Kubernetes version in the e2econfig
func conformanceImageID(e2eCtx *E2EContext) string {
	ver := e2eCtx.E2EConfig.GetVariable("CONFORMANCE_CI_ARTIFACTS_KUBERNETES_VERSION")
	strippedVer := strings.Replace(ver, "v", "", 1)
	amiName := AMIPrefix + strippedVer + "*"

	Byf("Searching for AMI: name=%s", amiName)
	ec2Svc := ec2.New(e2eCtx.AWSSession)
	filters := []*ec2.Filter{
		{
			Name:   aws.String("name"),
			Values: []*string{aws.String(amiName)},
		},
	}
	filters = append(filters, &ec2.Filter{
		Name:   aws.String("owner-id"),
		Values: []*string{aws.String(DefaultImageLookupOrg)},
	})
	resp, err := ec2Svc.DescribeImages(&ec2.DescribeImagesInput{
		Filters: filters,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(len(resp.Images)).To(Not(BeZero()))
	imageID := aws.StringValue(resp.Images[0].ImageId)
	Byf("Using AMI: image-id=%s", imageID)
	return imageID
}

func GetAvailabilityZones(sess client.ConfigProvider) []*ec2.AvailabilityZone {
	ec2Client := ec2.New(sess)
	azs, err := ec2Client.DescribeAvailabilityZones(nil)
	Expect(err).NotTo(HaveOccurred())
	return azs.AvailabilityZones
}

func DumpEKSClusters(ctx context.Context, e2eCtx *E2EContext) {
	logPath := filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName(), "aws-resources")
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		fmt.Fprintf(GinkgoWriter, "couldn't create directory: path=%s, err=%s", logPath, err)
	}
	fmt.Fprintf(GinkgoWriter, "folder created for eks clusters: %s\n", logPath)

	input := &eks.ListClustersInput{}
	eksClient := eks.New(e2eCtx.BootstratpUserAWSSession)
	output, err := eksClient.ListClusters(input)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "couldn't list EKS clusters: err=%s", err)
		return
	}

	for _, clusterName := range output.Clusters {
		describeInput := &eks.DescribeClusterInput{
			Name: clusterName,
		}
		describeOutput, err := eksClient.DescribeCluster(describeInput)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "couldn't describe EKS clusters: name=%s err=%s", *clusterName, err)
			continue
		}
		dumpEKSCluster(describeOutput.Cluster, logPath)
	}

}

func dumpEKSCluster(cluster *eks.Cluster, logPath string) {
	clusterYAML, err := yaml.Marshal(cluster)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "couldn't marshal cluster to yaml: name=%s err=%s", *cluster.Name, err)
		return
	}

	fileName := fmt.Sprintf("%s.yaml", *cluster.Name)
	clusterLog := path.Join(logPath, fileName)
	f, err := os.OpenFile(clusterLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "couldn't open log file: name=%s err=%s", clusterLog, err)
		return
	}
	defer f.Close()

	if err := ioutil.WriteFile(f.Name(), clusterYAML, 0600); err != nil {
		fmt.Fprintf(GinkgoWriter, "couldn't write cluster yaml to file: name=%s file=%s err=%s", *cluster.Name, f.Name(), err)
		return
	}
}

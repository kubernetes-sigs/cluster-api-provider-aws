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
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"

	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/service"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
)

func NewAWSSession() client.ConfigProvider {
	By("Getting an AWS IAM session")
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

// createCloudFormationStack ensures the cloudformation stack is up to date
func createCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	Byf("Creating AWS CloudFormation stack for AWS IAM resources: stack-name=%s", t.Spec.StackName)
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	cfnTemplate := t.RenderCloudFormation()
	Expect(
		cfnSvc.ReconcileBootstrapStack(t.Spec.StackName, *cfnTemplate),
	).To(Succeed())
}

// deleteCloudFormationStack removes the provisioned clusterawsadm stack
func deleteCloudFormationStack(prov client.ConfigProvider, t *cfn_bootstrap.Template) {
	Byf("Deleting %s CloudFormation stack", t.Spec.StackName)
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.DeleteStack(t.Spec.StackName),
	).To(Succeed())
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

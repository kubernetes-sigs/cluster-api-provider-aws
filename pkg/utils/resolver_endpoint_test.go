package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/onsi/gomega"
	"os"
	"strings"
	"testing"
)

func TestResolverEndpointAWSGov(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	err := os.Setenv("AWS_USE_FIPS_ENDPOINT", "true")
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	// Test us-gov and fips enabled regions
	err = os.Setenv("AWS_REGION", "us-west-2")
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	stsSvc := sts.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(stsSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(stsSvc.Endpoint, "sts-fips.us-west-2.amazonaws.com")).To(gomega.BeTrue())

	cfnSvc := cloudformation.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(cfnSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(cfnSvc.Endpoint, "cloudformation-fips.us-west-2.amazonaws.com")).To(gomega.BeTrue())

	err = os.Setenv("AWS_REGION", "us-east-2")
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}
	sess, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	stsSvc = sts.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(stsSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(stsSvc.Endpoint, "sts-fips.us-east-2.amazonaws.com")).To(gomega.BeTrue())

	cfnSvc = cloudformation.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(cfnSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(cfnSvc.Endpoint, "cloudformation-fips.us-east-2.amazonaws.com")).To(gomega.BeTrue())

	err = os.Setenv("AWS_REGION", endpoints.UsGovWest1RegionID)
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	sess, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	cfnSvc = cloudformation.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(cfnSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(cfnSvc.Endpoint, "cloudformation.us-gov-west-1.amazonaws.com")).To(gomega.BeTrue())

	stsSvc = sts.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(stsSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(stsSvc.Endpoint, "sts.us-gov-west-1.amazonaws.com")).To(gomega.BeTrue())

	// Test non fips endpoint regions
	err = os.Unsetenv("AWS_USE_FIPS_ENDPOINT")
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}
	err = os.Setenv("AWS_REGION", endpoints.EuNorth1RegionID)
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}
	sess, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	cfnSvc = cloudformation.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(cfnSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(cfnSvc.Endpoint, "cloudformation.eu-north-1.amazonaws.com")).To(gomega.BeTrue())

	stsSvc = sts.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(stsSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(stsSvc.Endpoint, "sts.amazonaws.com")).To(gomega.BeTrue())

	err = os.Setenv("AWS_REGION", endpoints.ApSoutheast1RegionID)
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	sess, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	cfnSvc = cloudformation.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(cfnSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(cfnSvc.Endpoint, "cloudformation.ap-southeast-1.amazonaws.com")).To(gomega.BeTrue())

	stsSvc = sts.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(stsSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(strings.Contains(stsSvc.Endpoint, "sts.amazonaws.com")).To(gomega.BeTrue())
}

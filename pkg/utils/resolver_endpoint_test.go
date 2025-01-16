package utils

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/onsi/gomega"
)

func TestResolverEndpointAWSSCS(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	err := os.Setenv("AWS_USE_FIPS_ENDPOINT", "true")
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	// Test us-gov and fips enabled regions
	err = os.Setenv("AWS_REGION", "us-iso-east-1")
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	stsSvc := cloudformation.New(sess, aws.NewConfig().WithEndpointResolver(CustomEndpointResolverForAWS()))
	g.Expect(stsSvc.ServiceID).ToNot(gomega.BeEmpty())
	g.Expect(stsSvc.Endpoint).To(gomega.Equal("https://cloudformation.us-iso-east-1.sc2s.sgov.gov"))
}

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

func TestDisableFipsEndpointForSC2S(t *testing.T) {
	tests := []struct {
		name               string
		region             string
		initialFipsState   endpoints.FIPSEndpointState
		expectedFipsState  endpoints.FIPSEndpointState
		expectModification bool
	}{
		{
			name:               "Disable FIPS for us-iso-east-1",
			region:             endpoints.UsIsoEast1RegionID,
			initialFipsState:   endpoints.FIPSEndpointStateEnabled,
			expectedFipsState:  endpoints.FIPSEndpointStateDisabled,
			expectModification: true,
		},
		{
			name:               "Disable FIPS for us-iso-west-1",
			region:             endpoints.UsIsoWest1RegionID,
			initialFipsState:   endpoints.FIPSEndpointStateEnabled,
			expectedFipsState:  endpoints.FIPSEndpointStateDisabled,
			expectModification: true,
		},
		{
			name:               "Disable FIPS for us-isob-east-1",
			region:             endpoints.UsIsobEast1RegionID,
			initialFipsState:   endpoints.FIPSEndpointStateEnabled,
			expectedFipsState:  endpoints.FIPSEndpointStateDisabled,
			expectModification: true,
		},
		{
			name:               "Do not modify FIPS for unknown region",
			region:             "unknown-region",
			initialFipsState:   endpoints.FIPSEndpointStateEnabled,
			expectedFipsState:  endpoints.FIPSEndpointStateEnabled,
			expectModification: false,
		}, {
			name:               "Do not modify FIPS for us-east-1",
			region:             endpoints.UsEast1RegionID,
			initialFipsState:   endpoints.FIPSEndpointStateEnabled,
			expectedFipsState:  endpoints.FIPSEndpointStateEnabled,
			expectModification: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			options := &endpoints.Options{
				UseFIPSEndpoint: tt.initialFipsState,
			}

			optFn := disableFipsEndpointForSC2S(tt.region)
			optFn(options)

			if options.UseFIPSEndpoint != tt.expectedFipsState {
				t.Errorf("Unexpected FIPS state for region %q: got %v, want %v", tt.region, options.UseFIPSEndpoint, tt.expectedFipsState)
			}

			modified := options.UseFIPSEndpoint != tt.initialFipsState
			if modified != tt.expectModification {
				t.Errorf("Unexpected modification for region %q: got %v, want %v", tt.region, modified, tt.expectModification)
			}
		})
	}
}

func TestPartitionForRegion(t *testing.T) {
	tests := []struct {
		name          string
		region        string
		wantPartition string
		wantError     bool
	}{
		{
			name:          "Standard AWS region",
			region:        "us-east-1",
			wantPartition: "aws",
			wantError:     false,
		},
		{
			name:          "China region",
			region:        "cn-north-1",
			wantPartition: "aws-cn",
			wantError:     false,
		},
		{
			name:          "GovCloud region",
			region:        "us-gov-west-1",
			wantPartition: "aws-us-gov",
			wantError:     false,
		},
		{
			name:          "ISO partition region",
			region:        "us-iso-east-1",
			wantPartition: "aws-iso",
			wantError:     false,
		},
		{
			name:          "ISOB partition region",
			region:        "us-isob-east-1",
			wantPartition: "aws-iso-b",
			wantError:     false,
		},
		{
			name:          "Unknown region partition aws",
			region:        "unknown-region",
			wantPartition: "aws",
			wantError:     false,
		},
		{
			name:          "Empty region",
			region:        "",
			wantPartition: "",
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			partition, err := PartitionForRegion(tt.region)

			if partition != tt.wantPartition {
				t.Errorf("Test %q failed: expected partition %q, got %q", tt.name, tt.wantPartition, partition)
			}

			if err != nil && !tt.wantError {
				t.Errorf("Test %q failed: expected no error, got %q", tt.name, err.Error())
			} else if err == nil && tt.wantError {
				t.Errorf("Test %q failed: expected error, got nil", tt.name)
			}
		})
	}
}

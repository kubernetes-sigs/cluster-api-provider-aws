package utils

import (
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"k8s.io/klog/v2/klogr"
	"os"
)

var isFIPSEndpointEnabled bool

func CustomEndpointResolverForAWS() endpoints.ResolverFunc {

	log := klogr.New()
	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {

		resolve, err := endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
		if err != nil {
			return resolve, err
		}

		log.V(1).Info("CustomEndpointResolverForAWS", " region: ", region, " service: ", service, " optFns: ", optFns)
		// Handle only for US-GOV regions exceptions
		switch region {
		case endpoints.UsGovEast1RegionID:
			switch service {
			case "cloudformation":
				resolve.URL = "https://cloudformation.us-gov-east-1.amazonaws.com"
			case "monitoring":
				resolve.URL = "https://monitoring.us-gov-east-1.amazonaws.com"
			case "events":
				resolve.URL = "https://events.us-gov-east-1.amazonaws.com"
			case "logs":
				resolve.URL = "https://logs.us-gov-east-1.amazonaws.com"
			case "elasticloadbalancing":
				resolve.URL = "https://elasticloadbalancing.us-gov-east-1.amazonaws.com"
			case "ssm":
				resolve.URL = "https://ssm.us-gov-east-1.amazonaws.com"
			case "support":
				resolve.URL = "https://support.us-gov-west-1.amazonaws.com"
			case "states":
				resolve.URL = "https://states-fips.us-gov-east-1.amazonaws.com"
			case "serverlessrepo":
				resolve.URL = "https://serverlessrepo.us-gov-east-1.amazonaws.com"
			case "sts":
				resolve.URL = "https://sts.us-gov-east-1.amazonaws.com"
			case "iam":
				resolve.URL = "https://iam.us-gov.amazonaws.com"
			case "cloudtrail":
				resolve.URL = "https://cloudtrail.us-gov-east-1.amazonaws.com"
			case "autoscaling-plans":
				resolve.URL = "https://autoscaling-plans.us-gov-east-1.amazonaws.com"
			case "autoscaling":
				resolve.URL = "https://autoscaling.us-gov-east-1.amazonaws.com"
			}

		case endpoints.UsGovWest1RegionID:
			switch service {
			case "cloudformation":
				resolve.URL = "https://cloudformation.us-gov-west-1.amazonaws.com"
			case "monitoring":
				resolve.URL = "https://monitoring.us-gov-west-1.amazonaws.com"
			case "events":
				resolve.URL = "https://events.us-gov-west-1.amazonaws.com"
			case "logs":
				resolve.URL = "https://logs.us-gov-west-1.amazonaws.com"
			case "elasticloadbalancing":
				resolve.URL = "https://elasticloadbalancing.us-gov-west-1.amazonaws.com"
			case "ssm":
				resolve.URL = "https://ssm.us-gov-west-1.amazonaws.com"
			case "support":
				resolve.URL = "https://support.us-gov-west-1.amazonaws.com"
			case "states":
				resolve.URL = "https://states.us-gov-west-1.amazonaws.com"
			case "serverlessrepo":
				resolve.URL = "https://serverlessrepo.us-gov-west-1.amazonaws.com"
			case "sts":
				resolve.URL = "https://sts.us-gov-west-1.amazonaws.com"
			case "iam":
				resolve.URL = "https://iam.us-gov.amazonaws.com"
			case "cloudtrail":
				resolve.URL = "https://cloudtrail.us-gov-west-1.amazonaws.com"
			case "autoscaling-plans":
				resolve.URL = "https://autoscaling-plans.us-gov-west-1.amazonaws.com"
			case "autoscaling":
				resolve.URL = "https://autoscaling.us-gov-west-1.amazonaws.com"
			}

		case endpoints.UsEast1RegionID:
			switch service {
			case "autoscaling":
				resolve.URL = "https://autoscaling.us-east-1.amazonaws.com"
			}

		case endpoints.UsEast2RegionID:
			switch service {
			case "autoscaling":
				resolve.URL = "https://autoscaling.us-east-2.amazonaws.com"
			}

		case endpoints.UsWest1RegionID:
			switch service {
			case "autoscaling":
				resolve.URL = "https://autoscaling.us-west-1.amazonaws.com"
			}

		case endpoints.UsWest2RegionID:
			switch service {
			case "autoscaling":
				resolve.URL = "https://autoscaling.us-west-2.amazonaws.com"
			}
		}

		log.V(1).Info("CustomEndpointResolverForAWS", "resolve: ", resolve)
		return resolve, nil
	}

	return resolver
}

func ResetFipsEndpointEnv(region string) error {
	log := klogr.New()
	if isFIPSEndpointEnabled && shouldResetFIPSEndpointEnv(region) {
		log.V(1).Info("ResetFipsEndpointEnv required for non fips regions")
		err := os.Unsetenv("AWS_USE_FIPS_ENDPOINT")
		if err != nil {
			log.Error(err, "Failed to unset env AWS_USE_FIPS_ENDPOINT")
			return err
		}
		isFIPSEndpointEnabled = false
	}

	return nil
}

func shouldResetFIPSEndpointEnv(region string) bool {
	switch region {
	case endpoints.UsEast1RegionID, endpoints.UsEast2RegionID, endpoints.UsWest1RegionID, endpoints.UsWest2RegionID, endpoints.UsGovEast1RegionID, endpoints.UsGovWest1RegionID:
	default:
		return true
	}
	return false
}

func init() {
	isEnabled := os.Getenv("AWS_USE_FIPS_ENDPOINT")
	if isEnabled == "true" {
		isFIPSEndpointEnabled = true
	}
}

func CustomEndpointResolverForAWSIRSA(s3region string) endpoints.ResolverFunc {

	log := klogr.New()
	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {

		resolve, err := endpoints.DefaultResolver().EndpointFor(service, s3region, optFns...)
		if err != nil {
			return resolve, err
		}

		log.V(1).Info("CustomEndpointResolverForAWSIRSA", "resolve: ", resolve)
		return resolve, nil
	}

	return resolver
}

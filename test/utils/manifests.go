package utils

import (
	"fmt"
	"os"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// GenerateAwsCredentialsSecretFromEnv generates secret with AWS credentials
func GenerateAwsCredentialsSecretFromEnv(secretName, namespace string) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			awsclient.AwsCredsSecretIDKey:     []byte(os.Getenv("AWS_ACCESS_KEY_ID")),
			awsclient.AwsCredsSecretAccessKey: []byte(os.Getenv("AWS_SECRET_ACCESS_KEY")),
		},
	}
}

func TestingMachineProviderConfig(awsCredentialsSecretName string, clusterID string) (clusterv1alpha1.ProviderConfig, error) {
	publicIP := true
	machinePc := &providerconfigv1.AWSMachineProviderConfig{
		AMI: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{
				{
					Name:   "tag:image_stage",
					Values: []string{"base"},
				},
				{
					Name:   "tag:operating_system",
					Values: []string{"rhel"},
				},
				{
					Name:   "tag:ready",
					Values: []string{"yes"},
				},
			},
		},
		CredentialsSecret: &apiv1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: providerconfigv1.Placement{
			Region:           "us-east-1",
			AvailabilityZone: "us-east-1a",
		},
		Subnet: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-worker-*", clusterID)},
				},
			},
		},
		Tags: []providerconfigv1.TagSpecification{
			{
				Name:  "openshift-node-group-config",
				Value: "node-config-master",
			},
			{
				Name:  "host-type",
				Value: "master",
			},
			{
				Name:  "sub-host-type",
				Value: "default",
			},
		},
		SecurityGroups: []providerconfigv1.AWSResourceReference{
			{
				Filters: []providerconfigv1.Filter{
					{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-*", clusterID)},
					},
				},
			},
		},
		PublicIP: &publicIP,
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return clusterv1alpha1.ProviderConfig{}, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeProviderConfig(machinePc)
	if err != nil {
		return clusterv1alpha1.ProviderConfig{}, fmt.Errorf("EncodeToProviderConfig failed: %v", err)
	}
	return *config, nil
}

func MasterMachineProviderConfig(awsCredentialsSecretName, masterUserDataSecretName, clusterID string) (clusterv1alpha1.ProviderConfig, error) {
	publicIP := true
	machinePc := &providerconfigv1.AWSMachineProviderConfig{
		AMI: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{
				{
					Name:   "tag:image_stage",
					Values: []string{"base"},
				},
				{
					Name:   "tag:operating_system",
					Values: []string{"rhel"},
				},
				{
					Name:   "tag:ready",
					Values: []string{"yes"},
				},
			},
		},
		CredentialsSecret: &apiv1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: providerconfigv1.Placement{
			Region:           "us-east-1",
			AvailabilityZone: "us-east-1a",
		},
		Subnet: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-worker-*", clusterID)},
				},
			},
		},
		SecurityGroups: []providerconfigv1.AWSResourceReference{
			{
				Filters: []providerconfigv1.Filter{
					{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-*", clusterID)},
					},
				},
			},
		},
		PublicIP: &publicIP,
		UserDataSecret: &apiv1.LocalObjectReference{
			Name: masterUserDataSecretName,
		},
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return clusterv1alpha1.ProviderConfig{}, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeProviderConfig(machinePc)
	if err != nil {
		return clusterv1alpha1.ProviderConfig{}, fmt.Errorf("EncodeToProviderConfig failed: %v", err)
	}
	return *config, nil
}

func WorkerMachineSetProviderConfig(awsCredentialsSecretName, workerUserDataSecretName, clusterID string) (clusterv1alpha1.ProviderConfig, error) {
	publicIP := true
	machinePc := &providerconfigv1.AWSMachineProviderConfig{
		AMI: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{
				{
					Name:   "tag:image_stage",
					Values: []string{"base"},
				},
				{
					Name:   "tag:operating_system",
					Values: []string{"rhel"},
				},
				{
					Name:   "tag:ready",
					Values: []string{"yes"},
				},
			},
		},
		CredentialsSecret: &apiv1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: providerconfigv1.Placement{
			Region:           "us-east-1",
			AvailabilityZone: "us-east-1a",
		},
		Subnet: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-worker-*", clusterID)},
				},
			},
		},
		SecurityGroups: []providerconfigv1.AWSResourceReference{
			{
				Filters: []providerconfigv1.Filter{
					{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-*", clusterID)},
					},
				},
			},
		},
		PublicIP: &publicIP,
		UserDataSecret: &apiv1.LocalObjectReference{
			Name: workerUserDataSecretName,
		},
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return clusterv1alpha1.ProviderConfig{}, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeProviderConfig(machinePc)
	if err != nil {
		return clusterv1alpha1.ProviderConfig{}, fmt.Errorf("EncodeToProviderConfig failed: %v", err)
	}
	return *config, nil
}

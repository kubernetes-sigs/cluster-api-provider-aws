package v1beta1

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEncodeAndDecodeProviderStatus(t *testing.T) {

	codec, err := NewCodec()
	if err != nil {
		t.Fatal(err)
	}
	time := metav1.Time{
		Time: time.Date(2018, 6, 3, 0, 0, 0, 0, time.Local),
	}

	instanceState := "running"
	instanceID := "id"
	providerStatus := &AWSMachineProviderStatus{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSMachineProviderStatus",
			APIVersion: "awsproviderconfig.openshift.io/v1beta1",
		},
		InstanceState: &instanceState,
		InstanceID:    &instanceID,
		Conditions: []AWSMachineProviderCondition{
			{
				Type:               "MachineCreation",
				Status:             "True",
				Reason:             "MachineCreationSucceeded",
				Message:            "machine successfully created",
				LastTransitionTime: time,
				LastProbeTime:      time,
			},
		},
	}
	providerStatusEncoded, err := codec.EncodeProviderStatus(providerStatus)
	if err != nil {
		t.Error(err)
	}

	// without deep copy
	{
		providerStatusDecoded := &AWSMachineProviderStatus{}
		codec.DecodeProviderStatus(providerStatusEncoded, providerStatusDecoded)
		if !reflect.DeepEqual(providerStatus, providerStatusDecoded) {
			t.Errorf("failed EncodeProviderStatus/DecodeProviderStatus. Expected: %+v, got: %+v", providerStatus, providerStatusDecoded)
		}
	}

	// with deep copy
	{
		providerStatusDecoded := &AWSMachineProviderStatus{}
		codec.DecodeProviderStatus(providerStatusEncoded.DeepCopy(), providerStatusDecoded)
		if !reflect.DeepEqual(providerStatus, providerStatusDecoded) {
			t.Errorf("failed EncodeProviderStatus/DecodeProviderStatus. Expected: %+v, got: %+v", providerStatus, providerStatusDecoded)
		}
	}
}

func TestEncodeAndDecodeProviderSpec(t *testing.T) {
	codec, err := NewCodec()
	if err != nil {
		t.Fatal(err)
	}

	publicIP := true
	amiID := "id"
	awsCredentialsSecretName := "test"
	clusterID := "test"

	providerConfig := &AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSMachineProviderConfig",
			APIVersion: "awsproviderconfig.openshift.io/v1beta1",
		},
		AMI: AWSResourceReference{
			Filters: []Filter{
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
			ID: &amiID,
		},
		CredentialsSecret: &apiv1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: Placement{
			Region:           "us-east-1",
			AvailabilityZone: "us-east-1a",
		},
		Subnet: AWSResourceReference{
			Filters: []Filter{
				{
					Name:   "tag:Name",
					Values: []string{clusterID},
				},
			},
		},
		Tags: []TagSpecification{
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
		SecurityGroups: []AWSResourceReference{
			{
				Filters: []Filter{
					{
						Name:   "tag:Name",
						Values: []string{clusterID},
					},
				},
			},
		},
		PublicIP: &publicIP,
		IAMInstanceProfile: &AWSResourceReference{
			ID: aws.String("IAMInstanceProfile"),
		},
		UserDataSecret: &corev1.LocalObjectReference{
			Name: "userSecret",
		},
		KeyName: aws.String("sshkeyname"),
		LoadBalancers: []LoadBalancerReference{
			{
				Name: "balancer",
				Type: ClassicLoadBalancerType,
			},
		},
	}

	providerConfigEncoded, err := codec.EncodeProviderSpec(providerConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Without deep copy
	{
		providerConfigDecoded := &AWSMachineProviderConfig{}
		codec.DecodeProviderSpec(providerConfigEncoded, providerConfigDecoded)

		if !reflect.DeepEqual(providerConfig, providerConfigDecoded) {
			t.Errorf("failed EncodeProviderSpec/DecodeProviderSpec. Expected: %+v, got: %+v", providerConfig, providerConfigDecoded)
		}
	}

	// With deep copy
	{
		providerConfigDecoded := &AWSMachineProviderConfig{}
		codec.DecodeProviderSpec(providerConfigEncoded, providerConfigDecoded)

		if !reflect.DeepEqual(providerConfig.DeepCopy(), providerConfigDecoded) {
			t.Errorf("failed EncodeProviderSpec/DecodeProviderSpec. Expected: %+v, got: %+v", providerConfig, providerConfigDecoded)
		}
	}
}

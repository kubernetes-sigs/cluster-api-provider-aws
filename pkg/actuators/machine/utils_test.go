package machine

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
}

func TestProviderConfigFromMachine(t *testing.T) {

	providerConfig := &providerconfigv1.AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "awsproviderconfig.openshift.io/v1beta1",
			Kind:       "AWSMachineProviderConfig",
		},
		InstanceType: "testInstance",
		AMI:          providerconfigv1.AWSResourceReference{ID: nil},
		Tags: []providerconfigv1.TagSpecification{
			{Name: "", Value: ""},
		},
		IAMInstanceProfile: &providerconfigv1.AWSResourceReference{ID: nil},
		UserDataSecret:     &corev1.LocalObjectReference{Name: ""},
		Subnet: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{{
				Name:   "tag:Name",
				Values: []string{""},
			}},
		},
		Placement: providerconfigv1.Placement{Region: "", AvailabilityZone: ""},
		SecurityGroups: []providerconfigv1.AWSResourceReference{{
			Filters: []providerconfigv1.Filter{{
				Name:   "tag:Name",
				Values: []string{""},
			}},
		}},
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Error(err)
	}
	encodedProviderSpec, err := codec.EncodeProviderSpec(providerConfig)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		machine *machinev1.Machine
	}{
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "configFromSpecProviderConfigValue",
					Namespace: "",
					Labels: map[string]string{
						"foo": "a",
					},
				},
				TypeMeta: metav1.TypeMeta{
					Kind: "Machine",
				},
				Spec: machinev1.MachineSpec{
					ProviderSpec: *encodedProviderSpec,
				},
			},
		},
	}

	for _, tc := range testCases {
		decodedProviderConfig, err := providerConfigFromMachine(tc.machine, codec)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(decodedProviderConfig, providerConfig) {
			t.Errorf("Test case %s. Expected: %v, got: %v", tc.machine.Name, providerConfig, decodedProviderConfig)
		}
	}
}

func TestExtractNodeAddresses(t *testing.T) {
	testCases := []struct {
		testcase          string
		instance          *ec2.Instance
		expectedAddresses []corev1.NodeAddress
	}{
		{
			testcase: "one-public",
			instance: &ec2.Instance{
				PublicIpAddress: aws.String("1.1.1.1"),
				PublicDnsName:   aws.String("ec2.example.net"),
			},
			expectedAddresses: []corev1.NodeAddress{
				{Type: corev1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: corev1.NodeExternalDNS, Address: "ec2.example.net"},
			},
		},
		{
			testcase: "one-private",
			instance: &ec2.Instance{
				PrivateDnsName: aws.String("ec2.example.net"),
				NetworkInterfaces: []*ec2.InstanceNetworkInterface{
					{
						Status: aws.String(ec2.NetworkInterfaceStatusInUse),
						PrivateIpAddresses: []*ec2.InstancePrivateIpAddress{
							{
								Primary:          aws.Bool(true),
								PrivateIpAddress: aws.String("10.0.0.5"),
							},
						},
					},
				},
			},
			expectedAddresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "10.0.0.5"},
				{Type: corev1.NodeInternalDNS, Address: "ec2.example.net"},
				{Type: corev1.NodeHostName, Address: "ec2.example.net"},
			},
		},
		{
			testcase: "multiple-private",
			instance: &ec2.Instance{
				PrivateDnsName: aws.String("ec2.example.net"),
				NetworkInterfaces: []*ec2.InstanceNetworkInterface{
					{
						Status: aws.String(ec2.NetworkInterfaceStatusInUse),
						PrivateIpAddresses: []*ec2.InstancePrivateIpAddress{
							{
								Primary:          aws.Bool(true),
								PrivateIpAddress: aws.String("10.0.0.5"),
							},
						},
					},
					{
						Status: aws.String(ec2.NetworkInterfaceStatusInUse),
						PrivateIpAddresses: []*ec2.InstancePrivateIpAddress{
							{
								Primary:          aws.Bool(false),
								PrivateIpAddress: aws.String("10.0.0.6"),
							},
						},
					},
				},
			},
			expectedAddresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "10.0.0.5"},
				{Type: corev1.NodeInternalIP, Address: "10.0.0.6"},
				{Type: corev1.NodeInternalDNS, Address: "ec2.example.net"},
				{Type: corev1.NodeHostName, Address: "ec2.example.net"},
			},
		},
		{
			testcase: "ipv6-private",
			instance: &ec2.Instance{
				PrivateDnsName: aws.String("ec2.example.net"),
				NetworkInterfaces: []*ec2.InstanceNetworkInterface{
					{
						Status: aws.String(ec2.NetworkInterfaceStatusInUse),
						Ipv6Addresses: []*ec2.InstanceIpv6Address{
							{
								Ipv6Address: aws.String("2600:1f18:4254:5100:ef8a:7b65:7782:9248"),
							},
						},
						PrivateIpAddresses: []*ec2.InstancePrivateIpAddress{
							{
								Primary:          aws.Bool(true),
								PrivateIpAddress: aws.String("10.0.0.5"),
							},
						},
					},
				},
			},
			expectedAddresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "2600:1f18:4254:5100:ef8a:7b65:7782:9248"},
				{Type: corev1.NodeInternalIP, Address: "10.0.0.5"},
				{Type: corev1.NodeInternalDNS, Address: "ec2.example.net"},
				{Type: corev1.NodeHostName, Address: "ec2.example.net"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcase, func(t *testing.T) {
			addresses, err := extractNodeAddresses(tc.instance)
			if err != nil {
				t.Errorf("Unexpected extractNodeAddresses error: %v", err)
			}

			if !equality.Semantic.DeepEqual(addresses, tc.expectedAddresses) {
				t.Errorf("expected: %v, got: %v", tc.expectedAddresses, addresses)
			}
		})
	}
}

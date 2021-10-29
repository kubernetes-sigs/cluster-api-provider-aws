package machine

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
}

func TestExtractNodeAddresses(t *testing.T) {
	testCases := []struct {
		testcase          string
		instance          *ec2.Instance
		expectedAddresses []corev1.NodeAddress
		domainNames       []string
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
			domainNames: nil,
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
			domainNames: nil,
		},
		{
			testcase: "custom-domain",
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
				{Type: corev1.NodeInternalDNS, Address: "ec2.openshift.com"},
				{Type: corev1.NodeInternalDNS, Address: "ec2.openshift.io"},
			},
			domainNames: []string{"openshift.com", "openshift.io"},
		},
		{
			testcase: "custom-domain no duplicates",
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
			domainNames: []string{"example.net", "example.net"},
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
			domainNames: nil,
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
			domainNames: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcase, func(t *testing.T) {
			addresses, err := extractNodeAddresses(tc.instance, tc.domainNames)
			if err != nil {
				t.Errorf("Unexpected extractNodeAddresses error: %v", err)
			}

			if !equality.Semantic.DeepEqual(addresses, tc.expectedAddresses) {
				t.Errorf("expected: %v, got: %v", tc.expectedAddresses, addresses)
			}
		})
	}
}

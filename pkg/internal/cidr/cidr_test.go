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

package cidr

import (
	"net"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSplitIntoSubnetsIPv4(t *testing.T) {
	RegisterTestingT(t)
	tests := []struct {
		name        string
		cidrblock   string
		subnetcount int
		expected    []*net.IPNet
	}{
		{
			// https://aws.amazon.com/about-aws/whats-new/2018/10/amazon-eks-now-supports-additional-vpc-cidr-blocks/
			name:        "default secondary cidr block configuration with primary cidr",
			cidrblock:   "100.64.0.0/16",
			subnetcount: 3,
			expected: []*net.IPNet{
				{
					IP:   net.IPv4(100, 64, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
				{
					IP:   net.IPv4(100, 64, 64, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
				{
					IP:   net.IPv4(100, 64, 128, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
			},
		},
		{
			// https://aws.amazon.com/about-aws/whats-new/2018/10/amazon-eks-now-supports-additional-vpc-cidr-blocks/
			name:        "default secondary cidr block configuration with alternative cidr",
			cidrblock:   "198.19.0.0/16",
			subnetcount: 3,
			expected: []*net.IPNet{
				{
					IP:   net.IPv4(198, 19, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
				{
					IP:   net.IPv4(198, 19, 64, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
				{
					IP:   net.IPv4(198, 19, 128, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
			},
		},
		{
			name:        "slash 16 cidr with one subnet",
			cidrblock:   "1.1.0.0/16",
			subnetcount: 1,
			expected: []*net.IPNet{
				{
					IP:   net.IPv4(1, 1, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 0, 0),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := SplitIntoSubnetsIPv4(tc.cidrblock, tc.subnetcount)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(ConsistOf(tc.expected))
		})
	}
}

func TestParseIPv4CIDR(t *testing.T) {
	RegisterTestingT(t)

	input := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
		"2001:db8::/32",
		"193.168.3.20/7",
	}

	output, err := GetIPv4Cidrs(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(HaveLen(1))
}

func TestParseIPv6CIDR(t *testing.T) {
	RegisterTestingT(t)

	input := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
		"2001:db8::/32",
		"193.168.3.20/7",
	}

	output, err := GetIPv6Cidrs(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(HaveLen(2))
}

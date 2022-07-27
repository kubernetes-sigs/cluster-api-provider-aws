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

package cidr_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/internal/cidr"
)

func TestParseIPv4CIDR(t *testing.T) {
	RegisterTestingT(t)

	input := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
		"2001:db8::/32",
		"193.168.3.20/7",
	}

	output, err := cidr.GetIPv4Cidrs(input)
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

	output, err := cidr.GetIPv6Cidrs(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(HaveLen(2))
}

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

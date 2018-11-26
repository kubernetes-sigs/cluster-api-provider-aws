package actuators

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
)

// AWSClients contains all the aws clients used by the scopes.
type AWSClients struct {
	EC2 ec2iface.EC2API
	ELB elbiface.ELBAPI
}

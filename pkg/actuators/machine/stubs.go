package machine

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
)

const (
	defaultNamespace         = "default"
	defaultAvailabilityZone  = "us-east-1a"
	region                   = "us-east-1"
	awsCredentialsSecretName = "aws-credentials-secret"
	userDataSecretName       = "aws-actuator-user-data-secret"

	keyName   = "aws-actuator-key-name"
	clusterID = "aws-actuator-cluster"
)

const userDataBlob = `#cloud-config
write_files:
- path: /root/node_bootstrap/node_settings.yaml
  owner: 'root:root'
  permissions: '0640'
  content: |
    node_config_name: node-config-master
runcmd:
- [ cat, /root/node_bootstrap/node_settings.yaml]
`

func stubProviderConfig() *providerconfigv1.AWSMachineProviderConfig {
	return &providerconfigv1.AWSMachineProviderConfig{
		AMI: providerconfigv1.AWSResourceReference{
			ID: aws.String("ami-a9acbbd6"),
		},
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: providerconfigv1.Placement{
			Region:           region,
			AvailabilityZone: defaultAvailabilityZone,
		},
		Subnet: providerconfigv1.AWSResourceReference{
			ID: aws.String("subnet-0e56b13a64ff8a941"),
		},
		IAMInstanceProfile: &providerconfigv1.AWSResourceReference{
			ID: aws.String("openshift_master_launch_instances"),
		},
		KeyName: aws.String(keyName),
		UserDataSecret: &corev1.LocalObjectReference{
			Name: userDataSecretName,
		},
		Tags: []providerconfigv1.TagSpecification{
			{Name: "openshift-node-group-config", Value: "node-config-master"},
			{Name: "host-type", Value: "master"},
			{Name: "sub-host-type", Value: "default"},
		},
		SecurityGroups: []providerconfigv1.AWSResourceReference{
			{ID: aws.String("sg-00868b02fbe29de17")},
			{ID: aws.String("sg-0a4658991dc5eb40a")},
			{ID: aws.String("sg-009a70e28fa4ba84e")},
			{ID: aws.String("sg-07323d56fb932c84c")},
			{ID: aws.String("sg-08b1ffd32874d59a2")},
		},
		PublicIP: aws.Bool(true),
		LoadBalancers: []providerconfigv1.LoadBalancerReference{
			{
				Name: "cluster-con",
				Type: providerconfigv1.ClassicLoadBalancerType,
			},
			{
				Name: "cluster-ext",
				Type: providerconfigv1.ClassicLoadBalancerType,
			},
			{
				Name: "cluster-int",
				Type: providerconfigv1.ClassicLoadBalancerType,
			},
			{
				Name: "cluster-net-lb",
				Type: providerconfigv1.NetworkLoadBalancerType,
			},
		},
	}
}

func stubMachine() (*machinev1.Machine, error) {
	machinePc := stubProviderConfig()

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return nil, fmt.Errorf("failed creating codec: %v", err)
	}
	providerSpec, err := codec.EncodeProviderSpec(machinePc)
	if err != nil {
		return nil, fmt.Errorf("codec.EncodeProviderSpec failed: %v", err)
	}

	machine := &machinev1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-actuator-testing-machine",
			Namespace: defaultNamespace,
			Labels: map[string]string{
				providerconfigv1.ClusterIDLabel: clusterID,
			},
			Annotations: map[string]string{
				// skip node draining since it's not mocked
				machinecontroller.ExcludeNodeDrainingAnnotation: "",
			},
		},

		Spec: machinev1.MachineSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"node-role.kubernetes.io/master": "",
					"node-role.kubernetes.io/infra":  "",
				},
			},
			ProviderSpec: *providerSpec,
		},
	}

	return machine, nil
}

func stubUserDataSecret() *corev1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userDataSecretName,
			Namespace: defaultNamespace,
		},
		Data: map[string][]byte{
			userDataSecretKey: []byte(userDataBlob),
		},
	}
}

func stubAwsCredentialsSecret() *corev1.Secret {
	return GenerateAwsCredentialsSecretFromEnv(awsCredentialsSecretName, defaultNamespace)
}

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

func stubInstance(imageID, instanceID string) *ec2.Instance {
	return &ec2.Instance{
		ImageId:    aws.String(imageID),
		InstanceId: aws.String(instanceID),
		State: &ec2.InstanceState{
			Name: aws.String(ec2.InstanceStateNameRunning),
			Code: aws.Int64(16),
		},
		LaunchTime:       aws.Time(time.Now()),
		PublicDnsName:    aws.String("publicDNS"),
		PrivateDnsName:   aws.String("privateDNS"),
		PublicIpAddress:  aws.String("1.1.1.1"),
		PrivateIpAddress: aws.String("1.1.1.1"),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("key"),
				Value: aws.String("value"),
			},
		},
		IamInstanceProfile: &ec2.IamInstanceProfile{
			Id: aws.String("profile"),
		},
		SubnetId: aws.String("subnetID"),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String("us-east-1a"),
		},
		SecurityGroups: []*ec2.GroupIdentifier{
			{
				GroupName: aws.String("groupName"),
			},
		},
	}
}

func stubPCSecurityGroups(groups []providerconfigv1.AWSResourceReference) *providerconfigv1.AWSMachineProviderConfig {
	pc := stubProviderConfig()
	pc.SecurityGroups = groups
	return pc
}

func stubPCSubnet(subnet providerconfigv1.AWSResourceReference) *providerconfigv1.AWSMachineProviderConfig {
	pc := stubProviderConfig()
	pc.Subnet = subnet
	return pc
}

func stubPCAMI(ami providerconfigv1.AWSResourceReference) *providerconfigv1.AWSMachineProviderConfig {
	pc := stubProviderConfig()
	pc.AMI = ami
	return pc
}

func stubDescribeLoadBalancersOutput() *elbv2.DescribeLoadBalancersOutput {
	return &elbv2.DescribeLoadBalancersOutput{
		LoadBalancers: []*elbv2.LoadBalancer{
			{
				LoadBalancerName: aws.String("lbname"),
				LoadBalancerArn:  aws.String("lbarn"),
			},
		},
	}
}

func stubDescribeTargetGroupsOutput() *elbv2.DescribeTargetGroupsOutput {
	return &elbv2.DescribeTargetGroupsOutput{
		TargetGroups: []*elbv2.TargetGroup{
			{
				TargetType:     aws.String(elbv2.TargetTypeEnumInstance),
				TargetGroupArn: aws.String("arn1"),
			},
			{
				TargetType:     aws.String(elbv2.TargetTypeEnumIp),
				TargetGroupArn: aws.String("arn2"),
			},
		},
	}
}

func stubReservation(imageID, instanceID string) *ec2.Reservation {
	az := defaultAvailabilityZone
	return &ec2.Reservation{
		Instances: []*ec2.Instance{
			{
				ImageId:    aws.String(imageID),
				InstanceId: aws.String(instanceID),
				State: &ec2.InstanceState{
					Name: aws.String(ec2.InstanceStateNameRunning),
					Code: aws.Int64(16),
				},
				LaunchTime: aws.Time(time.Now()),
				Placement: &ec2.Placement{
					AvailabilityZone: &az,
				},
			},
		},
	}
}

func stubDescribeInstancesOutput(imageID, instanceID string, state string) *ec2.DescribeInstancesOutput {
	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						ImageId:    aws.String(imageID),
						InstanceId: aws.String(instanceID),
						State: &ec2.InstanceState{
							Name: aws.String(state),
							Code: aws.Int64(16),
						},
						LaunchTime: aws.Time(time.Now()),
					},
				},
			},
		},
	}
}

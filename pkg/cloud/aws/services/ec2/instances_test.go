/*
Copyright 2018 The Kubernetes Authors.

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

package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
)

func TestInstanceIfExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check      func(instance *v1alpha1.Instance, err error)
	}{
		{
			name:       "does not exist",
			instanceID: "hello",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("hello")},
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String("test-vpc")},
						},
						{
							Name:   aws.String("instance-state-name"),
							Values: []*string{aws.String("pending"), aws.String("running")},
						},
					},
				})).
					Return(nil, awserrors.NewNotFound(errors.New("not found")))
			},
			check: func(instance *v1alpha1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance != nil {
					t.Fatalf("Did not expect anything but got something: %+v", instance)
				}
			},
		},
		{
			name:       "instance exists",
			instanceID: "id-1",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("id-1")},
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String("test-vpc")},
						},
						{
							Name:   aws.String("instance-state-name"),
							Values: []*string{aws.String("pending"), aws.String("running")},
						},
					},
				})).
					Return(&ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{
								Instances: []*ec2.Instance{
									{
										InstanceId:   aws.String("id-1"),
										InstanceType: aws.String("m5.large"),
										SubnetId:     aws.String("subnet-1"),
										ImageId:      aws.String("ami-1"),
										IamInstanceProfile: &ec2.IamInstanceProfile{
											Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
										},
										State: &ec2.InstanceState{
											Code: aws.Int64(16),
											Name: aws.String(ec2.StateAvailable),
										},
									},
								},
							},
						},
					}, nil)
			},
			check: func(instance *v1alpha1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance == nil {
					t.Fatalf("expected instance but got nothing")
				}

				if instance.ID != "id-1" {
					t.Fatalf("expected id-1 but got: %v", instance.ID)
				}
			},
		},
		{
			name:       "error describing instances",
			instanceID: "one",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInstances(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("one")},
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String("test-vpc")},
						},
						{
							Name:   aws.String("instance-state-name"),
							Values: []*string{aws.String("pending"), aws.String("running")},
						},
					},
				}).
					Return(nil, errors.New("some unknown error"))
			},
			check: func(i *v1alpha1.Instance, err error) {
				if err == nil {
					t.Fatalf("expected an error but got none.")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := actuators.NewScope(actuators.ScopeParams{
				Cluster: &clusterv1.Cluster{},
				AWSClients: actuators.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
			})

			scope.ClusterConfig = &v1alpha1.AWSClusterProviderSpec{
				NetworkSpec: v1alpha1.NetworkSpec{
					VPC: v1alpha1.VPCSpec{
						ID: "test-vpc",
					},
				},
			}

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			instance, err := s.InstanceIfExists(&tc.instanceID)
			tc.check(instance, err)
		})
	}
}

func TestTerminateInstance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	instanceNotFoundError := errors.New("instance not found")

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check      func(err error)
	}{
		{
			name:       "instance exists",
			instanceID: "i-exist",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: []*string{aws.String("i-exist")},
				})).
					Return(&ec2.TerminateInstancesOutput{}, nil)
			},
			check: func(err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name:       "instance does not exist",
			instanceID: "i-donotexist",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: []*string{aws.String("i-donotexist")},
				})).
					Return(&ec2.TerminateInstancesOutput{}, instanceNotFoundError)
			},
			check: func(err error) {
				if err == nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := actuators.NewScope(actuators.ScopeParams{
				Cluster: &clusterv1.Cluster{},
				AWSClients: actuators.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
			})

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			err = s.TerminateInstance(tc.instanceID)
			tc.check(err)
		})
	}
}

func TestCreateInstance(t *testing.T) {
	testCaCert := []byte(`
-----BEGIN CERTIFICATE-----
MIID6jCCAtICCQCa6H6nD76FxzANBgkqhkiG9w0BAQsFADCBtjELMAkGA1UEBhMC
VVMxCzAJBgNVBAgMAldBMRMwEQYDVQQHDAprdWJlcm5ldGVzMRQwEgYDVQQKDAtj
bHVzdGVyLWFwaTEhMB8GA1UECwwYY2x1c3Rlci1hcGktcHJvdmlkZXItYXdzMTAw
LgYDVQQDDCdzaWdzLms4cy5pby5jbHVzdGVyLWFwaS1wcm92aWRlci1hd3MuYWYx
GjAYBgkqhkiG9w0BCQEWC2Zvb0BiYXIuY29tMB4XDTE5MDExMTA5MTgxNVoXDTIx
MTAwNzA5MTgxNVowgbYxCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJXQTETMBEGA1UE
BwwKa3ViZXJuZXRlczEUMBIGA1UECgwLY2x1c3Rlci1hcGkxITAfBgNVBAsMGGNs
dXN0ZXItYXBpLXByb3ZpZGVyLWF3czEwMC4GA1UEAwwnc2lncy5rOHMuaW8uY2x1
c3Rlci1hcGktcHJvdmlkZXItYXdzLmFmMRowGAYJKoZIhvcNAQkBFgtmb29AYmFy
LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKHVoYcq6NiS2lch
ai62dDU+wStXJzFkF3URQ7auYDmL3Xz+01yxdARdafO3fweXSsfuxcGZ/DDBzRBB
ROXeJI1zxV6xk+OlI0puOabo6m5ji4RdTTFqt94afnK43qcDMDOnh0u6F5UZXZlr
T7XNnO++6e7elZ+9jJJ/NKPXDGKo9+M7kmypTLcI5b5pH4qn1coe8a5Li+FQONEM
j+Ttomqr0s84DyFBSNZYKvRVL1AdH/6r213pco5Qm9RDkw9HZr83Y1PjyQ77C7FQ
IPquny5XkjZjq65Bz8I8s+MoPQgBOr8JvVfc3Jt8u10qD4JOeRFnhZOygaApgswg
9XZhdMsCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAVfucXOzEy88NQ+fz5FV1D1PO
No6uqi2Q9fqGU9Lfnj3PhXr0sb0tAXGnZEg8i1317xMXqzA9J9umqg3ADsOsR3sL
NR41dkjP2ROfTW1wkEGBaRzp/TOagMy1IeeS9MPd4gRH3cZqgUvrQJCrX8878gxk
jor3R8gPhjvV74KrZD4lIF7IHHv4cCBaejm+3GwOIbTNoHXa4PadVwbcjWp6P8UB
dTga1FiyISsMchVaVKD5aX7hkxMP1/C98KdVzWQ4k12TBOhZDYUS67M4ibBtw/og
vuO9LYxDXLVY9F7W4ccyCqe27Cj1xyAvdZxwhITrib8Wg5CMqoRpqTw5V3+TpA==
-----END CERTIFICATE-----
	`)
	testcases := []struct {
		name          string
		machine       clusterv1.Machine
		machineConfig *v1alpha1.AWSMachineProviderSpec
		clusterStatus *v1alpha1.AWSClusterProviderStatus
		clusterConfig *v1alpha1.AWSClusterProviderSpec
		cluster       clusterv1.Cluster
		expect        func(m *mock_ec2iface.MockEC2APIMockRecorder)
		check         func(instance *v1alpha1.Instance, err error)
	}{
		{
			name: "simple",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
			},
			machineConfig: &v1alpha1.AWSMachineProviderSpec{
				AMI: v1alpha1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
			},
			clusterStatus: &v1alpha1.AWSClusterProviderStatus{
				Network: v1alpha1.Network{
					SecurityGroups: map[v1alpha1.SecurityGroupRole]*v1alpha1.SecurityGroup{
						v1alpha1.SecurityGroupControlPlane: {
							ID: "1",
						},
						v1alpha1.SecurityGroupNode: {
							ID: "2",
						},
					},
					APIServerELB: v1alpha1.ClassicELB{
						DNSName: "test-apiserver.us-east-1.aws",
					},
				},
			},
			clusterConfig: &v1alpha1.AWSClusterProviderSpec{
				NetworkSpec: v1alpha1.NetworkSpec{
					Subnets: v1alpha1.Subnets{
						&v1alpha1.SubnetSpec{
							ID:       "subnet-1",
							IsPublic: false,
						},
						&v1alpha1.SubnetSpec{
							IsPublic: false,
						},
					},
				},
				CAKeyPair: v1alpha1.KeyPair{
					Cert: testCaCert,
					Key:  []byte("y"),
				},
			},
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test1",
				},
				Spec: clusterv1.ClusterSpec{
					ClusterNetwork: clusterv1.ClusterNetworkingConfig{
						ServiceDomain: "cluster.local",
						Services: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
						Pods: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name: aws.String("ami-1"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								IamInstanceProfile: &ec2.IamInstanceProfile{
									Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
								},
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-1"),
								ImageId:      aws.String("ami-1"),
							},
						},
					}, nil)
				m.WaitUntilInstanceRunningWithContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			check: func(instance *v1alpha1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			// defer mockCtrl.Finish()
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := actuators.NewMachineScope(actuators.MachineScopeParams{
				Cluster: &tc.cluster,
				Machine: &clusterv1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"set": "node",
						},
					},
				},
				AWSClients: actuators.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
			})

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			scope.Scope.ClusterConfig = tc.clusterConfig
			scope.Scope.ClusterStatus = tc.clusterStatus

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope.Scope)
			instance, err := s.createInstance(scope, "token", "kubeconfig")
			tc.check(instance, err)
		})
	}
}

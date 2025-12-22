/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilfeature "k8s.io/component-base/featuregate/testing"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	utildefaulting "sigs.k8s.io/cluster-api-provider-aws/v2/util/defaulting"
)

func TestMachineDefault(t *testing.T) {
	machine := &AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	t.Run("for AWSMachine", utildefaulting.DefaultValidateTest(context.Background(), machine, &awsMachineWebhook{}))
	g := NewWithT(t)
	err := (&awsMachineWebhook{}).Default(context.Background(), machine)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(machine.Spec.CloudInit.SecureSecretsBackend).To(Equal(SecretBackendSecretsManager))
}

func TestAWSMachineCreate(t *testing.T) {
	tests := []struct {
		name    string
		machine *AWSMachine
		wantErr bool
	}{
		{
			name: "ensure IOPS exists if type equal to io1",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						Type: "io1",
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io2",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						Type: "io2",
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure root volume throughput is nonnegative",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						Throughput: aws.Int64(-125),
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure root volume with device name works (for clusterctl move)",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						DeviceName: "name",
						Type:       "gp2",
						Size:       *aws.Int64(8),
					},
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "ensure non root volume have device names",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io1 for non root volumes",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{
							DeviceName: "name",
							Type:       "io1",
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io2 for non root volumes",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{
							DeviceName: "name",
							Type:       "io2",
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure non root volume throughput is nonnegative",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{
							Throughput: aws.Int64(-125),
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "additional security groups may have id",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: aws.String("id"),
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "additional security groups may have filters",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							Filters: []Filter{
								{
									Name:   "example-name",
									Values: []string{"example-value"},
								},
							},
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "additional security groups can't have both id and filters",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: aws.String("id"),
							Filters: []Filter{
								{
									Name:   "example-name",
									Values: []string{"example-value"},
								},
							},
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "valid additional tags are accepted",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalTags: Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid case, MarketType set to MarketTypeCapacityBlock and spotMarketOptions are specified",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					MarketType:        MarketTypeCapacityBlock,
					SpotMarketOptions: &SpotMarketOptions{},
					InstanceType:      "test",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid case, MarketType set to MarketTypeOnDemand and spotMarketOptions are specified",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					MarketType:        MarketTypeOnDemand,
					SpotMarketOptions: &SpotMarketOptions{},
					InstanceType:      "test",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid MarketType set to MarketTypeCapacityBlock is specified and CapacityReservationId is not provided",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					MarketType:   MarketTypeCapacityBlock,
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "valid MarketType set to MarketTypeCapacityBlock and CapacityReservationId are specified",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					MarketType:            MarketTypeCapacityBlock,
					CapacityReservationID: aws.String("cr-12345678901234567"),
					InstanceType:          "test",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid case, CapacityReservationId is set and CapacityReservationPreference is not `capacity-reservation-only`",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: CapacityReservationPreferenceNone,
				},
			},
			wantErr: true,
		},
		{
			name: "valid CapacityReservationId is set and CapacityReservationPreference is not specified",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType:          "test",
					CapacityReservationID: aws.String("cr-12345678901234567"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid CapacityReservationId is set and CapacityReservationPreference is `capacity-reservation-only`",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: CapacityReservationPreferenceOnly,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid CapacityReservationPreference is `CapacityReservationsOnly` and MarketType is `Spot`",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: CapacityReservationPreferenceOnly,
					MarketType:                    MarketTypeSpot,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CapacityReservationPreference is `CapacityReservationsOnly` and SpotMarketOptions is non-nil",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: CapacityReservationPreferenceOnly,
					SpotMarketOptions:             &SpotMarketOptions{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty instance type not allowed",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "",
				},
			},
			wantErr: true,
		},
		{
			name: "instance type minimum length is 2",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "t",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tags return error",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalTags: Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ignition proxy and TLS can be from version 3.1",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					Ignition: &Ignition{
						Version: "3.1",
						Proxy: &IgnitionProxy{
							HTTPProxy: ptr.To("http://proxy.example.com:3128"),
						},
						TLS: &IgnitionTLS{
							CASources: []IgnitionCASource{"s3://example.com/ca.pem"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ignition tls with invalid CASources URL",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					Ignition: &Ignition{
						Version: "3.1",
						TLS: &IgnitionTLS{
							CASources: []IgnitionCASource{"data;;"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ignition proxy with valid URLs, and noproxy",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					Ignition: &Ignition{
						Version: "3.1",
						Proxy: &IgnitionProxy{
							HTTPProxy:  ptr.To("http://proxy.example.com:3128"),
							HTTPSProxy: ptr.To("https://proxy.example.com:3128"),
							NoProxy: []IgnitionNoProxy{
								"10.0.0.1",         // single ip
								"example.com",      // domain
								".example.com",     // all subdomains
								"example.com:3128", // domain with port
								"10.0.0.1:3128",    // ip with port
								"10.0.0.0/8",       // cidr block
								"*",                // no proxy wildcard
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ignition proxy with invalid HTTPProxy URL",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					Ignition: &Ignition{
						Version: "3.1",
						Proxy: &IgnitionProxy{
							HTTPProxy: ptr.To("*:80"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ignition proxy with invalid HTTPSProxy URL",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					Ignition: &Ignition{
						Version: "3.1",
						Proxy: &IgnitionProxy{
							HTTPSProxy: ptr.To("*:80"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ignition proxy with invalid noproxy URL",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					Ignition: &Ignition{
						Version: "3.1",
						Proxy: &IgnitionProxy{
							NoProxy: []IgnitionNoProxy{"&"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "cannot use ignition proxy with version 2.3",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					Ignition: &Ignition{
						Version: "2.3.0",
						Proxy: &IgnitionProxy{
							HTTPProxy: ptr.To("http://proxy.example.com:3128"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "hostAffinity=invalid is invalid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("invalid"),
				},
			},
			wantErr: true,
		},
		{
			name: "hostAffinity=host does not require hostID or dynamicHostAllocation",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("host"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=host with hostID is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("host"),
					HostID:       ptr.To("h-09dcf61cb388b0149"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=host with dynamicHostAllocation is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("host"),
					DynamicHostAllocation: &DynamicHostAllocationSpec{
						Tags: map[string]string{"env": "test"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=default without hostID and dynamicHostAllocation is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("default"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=default with hostID is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("default"),
					HostID:       ptr.To("h-09dcf61cb388b0149"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=default with dynamicHostAllocation is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("default"),
					DynamicHostAllocation: &DynamicHostAllocationSpec{
						Tags: map[string]string{"env": "test"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity omitted (=default) without hostID and dynamicHostAllocation is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity omitted (=default) with hostID is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostID:       ptr.To("h-09dcf61cb388b0149"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity omitted (=default) with dynamicHostAllocation is valid",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					DynamicHostAllocation: &DynamicHostAllocationSpec{
						Tags: map[string]string{"env": "test"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=host with both hostID and dynamicHostAllocation is not valid (mutually exclusive)",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("host"),
					HostID:       aws.String("h-1234567890abcdef0"),
					DynamicHostAllocation: &DynamicHostAllocationSpec{
						Tags: map[string]string{
							"Environment": "test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "hostAffinity=default with both hostID and dynamicHostAllocation is not valid (mutually exclusive)",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("default"),
					HostID:       aws.String("h-1234567890abcdef0"),
					DynamicHostAllocation: &DynamicHostAllocationSpec{
						Tags: map[string]string{
							"Environment": "test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "create with valid BYOIPv4",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &ElasticIPPool{
						PublicIpv4Pool:              aws.String("ipv4pool-ec2-0123456789abcdef0"),
						PublicIpv4PoolFallBackOrder: ptr.To(PublicIpv4PoolFallbackOrderAmazonPool),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error when BYOIPv4 without fallback",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &ElasticIPPool{
						PublicIpv4Pool: aws.String("ipv4pool-ec2-0123456789abcdef0"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error when BYOIPv4 without public ipv4 pool",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &ElasticIPPool{
						PublicIpv4PoolFallBackOrder: ptr.To(PublicIpv4PoolFallbackOrderAmazonPool),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error when BYOIPv4 with non-public IP set",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(false),
					ElasticIPPool: &ElasticIPPool{
						PublicIpv4Pool:              aws.String("ipv4pool-ec2-0123456789abcdef0"),
						PublicIpv4PoolFallBackOrder: ptr.To(PublicIpv4PoolFallbackOrderAmazonPool),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error when BYOIPv4 with invalid pool name",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &ElasticIPPool{
						PublicIpv4Pool:              aws.String("ipv4poolx-ec2-0123456789abcdef"),
						PublicIpv4PoolFallBackOrder: ptr.To(PublicIpv4PoolFallbackOrderAmazonPool),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utilfeature.SetFeatureGateDuringTest(t, feature.Gates, feature.BootstrapFormatIgnition, true)

			machine := tt.machine.DeepCopy()
			machine.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "machine-",
				Namespace:    "default",
			}
			ctx := context.TODO()
			if err := testEnv.Create(ctx, machine); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
			testEnv.Delete(ctx, machine)
		})
	}
}

func TestAWSMachineUpdate(t *testing.T) {
	tests := []struct {
		name       string
		oldMachine *AWSMachine
		newMachine *AWSMachine
		wantErr    bool
	}{
		{
			name: "change in providerid, cloudinit, tags, securitygroups",
			oldMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID:               nil,
					AdditionalTags:           nil,
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID:   ptr.To[string]("ID"),
					InstanceType: "test",
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: ptr.To[string]("ID"),
						},
					},
					CloudInit: CloudInit{
						SecretPrefix: "test",
						SecretCount:  5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "change in fields other than providerid, tags and securitygroups",
			oldMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID:               nil,
					AdditionalTags:           nil,
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ImageLookupOrg: "test",
					InstanceType:   "test",
					ProviderID:     ptr.To[string]("ID"),
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: ptr.To[string]("ID"),
						},
					},
					PrivateDNSName: &PrivateDNSName{
						EnableResourceNameDNSAAAARecord: aws.Bool(true),
						EnableResourceNameDNSARecord:    aws.Bool(true),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "change in tags adding invalid ones",
			oldMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID: nil,
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID: nil,
					AdditionalTags: Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			machine := tt.oldMachine.DeepCopy()
			machine.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "machine-",
				Namespace:    "default",
			}
			if err := testEnv.Create(ctx, machine); err != nil {
				t.Errorf("failed to create machine: %v", err)
			}
			machine.Spec = tt.newMachine.Spec
			if err := testEnv.Update(ctx, machine); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAWSMachineSecretsBackend(t *testing.T) {
	baseMachine := &AWSMachine{
		Spec: AWSMachineSpec{
			ProviderID:               nil,
			AdditionalTags:           nil,
			AdditionalSecurityGroups: nil,
			InstanceType:             "test",
		},
	}

	tests := []struct {
		name                   string
		cloudInit              CloudInit
		expectedSecretsBackend string
	}{
		{
			name:                   "with insecure skip secrets manager unset",
			cloudInit:              CloudInit{InsecureSkipSecretsManager: false},
			expectedSecretsBackend: "secrets-manager",
		},
		{
			name:                   "with insecure skip secrets manager unset and secrets backend set",
			cloudInit:              CloudInit{InsecureSkipSecretsManager: false, SecureSecretsBackend: "ssm-parameter-store"},
			expectedSecretsBackend: "ssm-parameter-store",
		},
		{
			name:                   "with insecure skip secrets manager set",
			cloudInit:              CloudInit{InsecureSkipSecretsManager: true},
			expectedSecretsBackend: "",
		},
	}

	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			machine := baseMachine.DeepCopy()
			machine.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "machine-",
				Namespace:    "default",
			}
			machine.Spec.CloudInit = tt.cloudInit
			if err := testEnv.Create(ctx, machine); err != nil {
				t.Errorf("failed to create machine: %v", err)
			}
			g := NewWithT(t)
			g.Expect(machine.Spec.CloudInit.SecureSecretsBackend).To(Equal(SecretBackend(tt.expectedSecretsBackend)))
		})
	}
}

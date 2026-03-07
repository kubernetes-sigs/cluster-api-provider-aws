/*
Copyright 2026 The Kubernetes Authors.

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

package webhooks

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilfeature "k8s.io/component-base/featuregate/testing"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	utildefaulting "sigs.k8s.io/cluster-api-provider-aws/v2/util/defaulting"
)

func TestMachineDefault(t *testing.T) {
	machine := &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	t.Run("for AWSMachine", utildefaulting.DefaultValidateTest(context.Background(), machine, &AWSMachine{}))
	g := NewWithT(t)
	err := (&AWSMachine{}).Default(context.Background(), machine)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(machine.Spec.CloudInit.SecureSecretsBackend).To(Equal(infrav1.SecretBackendSecretsManager))
}

func TestAWSMachineCreate(t *testing.T) {
	tests := []struct {
		name    string
		machine *infrav1.AWSMachine
		wantErr bool
	}{
		{
			name: "ensure IOPS exists if type equal to io1",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					RootVolume: &infrav1.Volume{
						Type: "io1",
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io2",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					RootVolume: &infrav1.Volume{
						Type: "io2",
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure root volume throughput is within range",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					RootVolume: &infrav1.Volume{
						Throughput: aws.Int64(-125),
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure root volume with device name works (for clusterctl move)",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					RootVolume: &infrav1.Volume{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					NonRootVolumes: []infrav1.Volume{
						{},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io1 for non root volumes",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					NonRootVolumes: []infrav1.Volume{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					NonRootVolumes: []infrav1.Volume{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					NonRootVolumes: []infrav1.Volume{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					AdditionalSecurityGroups: []infrav1.AWSResourceReference{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					AdditionalSecurityGroups: []infrav1.AWSResourceReference{
						{
							Filters: []infrav1.Filter{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					AdditionalSecurityGroups: []infrav1.AWSResourceReference{
						{
							ID: aws.String("id"),
							Filters: []infrav1.Filter{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					AdditionalTags: infrav1.Tags{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					MarketType:        infrav1.MarketTypeCapacityBlock,
					SpotMarketOptions: &infrav1.SpotMarketOptions{},
					InstanceType:      "test",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid case, MarketType set to MarketTypeOnDemand and spotMarketOptions are specified",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					MarketType:        infrav1.MarketTypeOnDemand,
					SpotMarketOptions: &infrav1.SpotMarketOptions{},
					InstanceType:      "test",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid MarketType set to MarketTypeCapacityBlock is specified and CapacityReservationId is not provided",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					MarketType:   infrav1.MarketTypeCapacityBlock,
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "valid MarketType set to MarketTypeCapacityBlock and CapacityReservationId are specified",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					MarketType:            infrav1.MarketTypeCapacityBlock,
					CapacityReservationID: aws.String("cr-12345678901234567"),
					InstanceType:          "test",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid case, CapacityReservationId is set and CapacityReservationPreference is not `capacity-reservation-only`",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: infrav1.CapacityReservationPreferenceNone,
				},
			},
			wantErr: true,
		},
		{
			name: "valid CapacityReservationId is set and CapacityReservationPreference is not specified",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:          "test",
					CapacityReservationID: aws.String("cr-12345678901234567"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid CapacityReservationId is set and CapacityReservationPreference is `capacity-reservation-only`",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: infrav1.CapacityReservationPreferenceOnly,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid CapacityReservationPreference is `CapacityReservationsOnly` and MarketType is `Spot`",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: infrav1.CapacityReservationPreferenceOnly,
					MarketType:                    infrav1.MarketTypeSpot,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CapacityReservationPreference is `CapacityReservationsOnly` and SpotMarketOptions is non-nil",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:                  "test",
					CapacityReservationID:         aws.String("cr-12345678901234567"),
					CapacityReservationPreference: infrav1.CapacityReservationPreferenceOnly,
					SpotMarketOptions:             &infrav1.SpotMarketOptions{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty instance type not allowed",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "",
				},
			},
			wantErr: true,
		},
		{
			name: "instance type minimum length is 2",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "t",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tags return error",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					AdditionalTags: infrav1.Tags{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Ignition: &infrav1.Ignition{
						Version: "3.1",
						Proxy: &infrav1.IgnitionProxy{
							HTTPProxy: ptr.To("http://proxy.example.com:3128"),
						},
						TLS: &infrav1.IgnitionTLS{
							CASources: []infrav1.IgnitionCASource{"s3://example.com/ca.pem"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ignition tls with invalid CASources URL",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Ignition: &infrav1.Ignition{
						Version: "3.1",
						TLS: &infrav1.IgnitionTLS{
							CASources: []infrav1.IgnitionCASource{"data;;"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ignition proxy with valid URLs, and noproxy",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Ignition: &infrav1.Ignition{
						Version: "3.1",
						Proxy: &infrav1.IgnitionProxy{
							HTTPProxy:  ptr.To("http://proxy.example.com:3128"),
							HTTPSProxy: ptr.To("https://proxy.example.com:3128"),
							NoProxy: []infrav1.IgnitionNoProxy{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Ignition: &infrav1.Ignition{
						Version: "3.1",
						Proxy: &infrav1.IgnitionProxy{
							HTTPProxy: ptr.To("*:80"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ignition proxy with invalid HTTPSProxy URL",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Ignition: &infrav1.Ignition{
						Version: "3.1",
						Proxy: &infrav1.IgnitionProxy{
							HTTPSProxy: ptr.To("*:80"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ignition proxy with invalid noproxy URL",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Ignition: &infrav1.Ignition{
						Version: "3.1",
						Proxy: &infrav1.IgnitionProxy{
							NoProxy: []infrav1.IgnitionNoProxy{"&"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "cannot use ignition proxy with version 2.3",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Ignition: &infrav1.Ignition{
						Version: "2.3.0",
						Proxy: &infrav1.IgnitionProxy{
							HTTPProxy: ptr.To("http://proxy.example.com:3128"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "hostAffinity=invalid is invalid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("invalid"),
				},
			},
			wantErr: true,
		},
		{
			name: "hostAffinity=host does not require hostID or dynamicHostAllocation",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostAffinity: ptr.To("host"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=host with hostID is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostAffinity: ptr.To("host"),
					HostID:       ptr.To("h-09dcf61cb388b0149"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=host with dynamicHostAllocation is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostAffinity: ptr.To("host"),
					DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
						Tags: map[string]string{"env": "test"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=default without hostID and dynamicHostAllocation is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					HostAffinity: ptr.To("default"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=default with hostID is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostAffinity: ptr.To("default"),
					HostID:       ptr.To("h-09dcf61cb388b0149"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=default with dynamicHostAllocation is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostAffinity: ptr.To("default"),
					DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
						Tags: map[string]string{"env": "test"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity omitted (=default) without hostID and dynamicHostAllocation is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity omitted (=default) with hostID is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostID:       ptr.To("h-09dcf61cb388b0149"),
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity omitted (=default) with dynamicHostAllocation is valid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
						Tags: map[string]string{"env": "test"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hostAffinity=host with both hostID and dynamicHostAllocation is not valid (mutually exclusive)",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostAffinity: ptr.To("host"),
					HostID:       aws.String("h-1234567890abcdef0"),
					DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
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
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "host",
					HostAffinity: ptr.To("default"),
					HostID:       aws.String("h-1234567890abcdef0"),
					DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
						Tags: map[string]string{
							"Environment": "test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "hostID without tenancy=host is invalid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "default",
					HostID:       ptr.To("h-09dcf61cb388b0149"),
				},
			},
			wantErr: true,
		},
		{
			name: "hostAffinity=host without tenancy=host is invalid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "default",
					HostAffinity: ptr.To("host"),
				},
			},
			wantErr: true,
		},
		{
			name: "dynamicHostAllocation without tenancy=host is invalid",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "test",
					Tenancy:      "dedicated",
					DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
						Tags: map[string]string{"env": "test"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "create with valid BYOIPv4",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &infrav1.ElasticIPPool{
						PublicIpv4Pool:              aws.String("ipv4pool-ec2-0123456789abcdef0"),
						PublicIpv4PoolFallBackOrder: ptr.To(infrav1.PublicIpv4PoolFallbackOrderAmazonPool),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error when BYOIPv4 without fallback",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &infrav1.ElasticIPPool{
						PublicIpv4Pool: aws.String("ipv4pool-ec2-0123456789abcdef0"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error when BYOIPv4 without public ipv4 pool",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &infrav1.ElasticIPPool{
						PublicIpv4PoolFallBackOrder: ptr.To(infrav1.PublicIpv4PoolFallbackOrderAmazonPool),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error when BYOIPv4 with non-public IP set",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(false),
					ElasticIPPool: &infrav1.ElasticIPPool{
						PublicIpv4Pool:              aws.String("ipv4pool-ec2-0123456789abcdef0"),
						PublicIpv4PoolFallBackOrder: ptr.To(infrav1.PublicIpv4PoolFallbackOrderAmazonPool),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error when BYOIPv4 with invalid pool name",
			machine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					InstanceType: "type",
					PublicIP:     aws.Bool(true),
					ElasticIPPool: &infrav1.ElasticIPPool{
						PublicIpv4Pool:              aws.String("ipv4poolx-ec2-0123456789abcdef"),
						PublicIpv4PoolFallBackOrder: ptr.To(infrav1.PublicIpv4PoolFallbackOrderAmazonPool),
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
		oldMachine *infrav1.AWSMachine
		newMachine *infrav1.AWSMachine
		wantErr    bool
	}{
		{
			name: "change in providerid, cloudinit, tags, securitygroups",
			oldMachine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					ProviderID:               nil,
					AdditionalTags:           nil,
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			newMachine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					ProviderID:   ptr.To[string]("ID"),
					InstanceType: "test",
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []infrav1.AWSResourceReference{
						{
							ID: ptr.To[string]("ID"),
						},
					},
					CloudInit: infrav1.CloudInit{
						SecretPrefix: "test",
						SecretCount:  5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "change in fields other than providerid, tags and securitygroups",
			oldMachine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					ProviderID:               nil,
					AdditionalTags:           nil,
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			newMachine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					ImageLookupOrg: "test",
					InstanceType:   "test",
					ProviderID:     ptr.To[string]("ID"),
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []infrav1.AWSResourceReference{
						{
							ID: ptr.To[string]("ID"),
						},
					},
					PrivateDNSName: &infrav1.PrivateDNSName{
						EnableResourceNameDNSAAAARecord: aws.Bool(true),
						EnableResourceNameDNSARecord:    aws.Bool(true),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "change in tags adding invalid ones",
			oldMachine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					ProviderID: nil,
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			newMachine: &infrav1.AWSMachine{
				Spec: infrav1.AWSMachineSpec{
					ProviderID: nil,
					AdditionalTags: infrav1.Tags{
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
	baseMachine := &infrav1.AWSMachine{
		Spec: infrav1.AWSMachineSpec{
			ProviderID:               nil,
			AdditionalTags:           nil,
			AdditionalSecurityGroups: nil,
			InstanceType:             "test",
		},
	}

	tests := []struct {
		name                   string
		cloudInit              infrav1.CloudInit
		expectedSecretsBackend string
	}{
		{
			name:                   "with insecure skip secrets manager unset",
			cloudInit:              infrav1.CloudInit{InsecureSkipSecretsManager: false},
			expectedSecretsBackend: "secrets-manager",
		},
		{
			name:                   "with insecure skip secrets manager unset and secrets backend set",
			cloudInit:              infrav1.CloudInit{InsecureSkipSecretsManager: false, SecureSecretsBackend: "ssm-parameter-store"},
			expectedSecretsBackend: "ssm-parameter-store",
		},
		{
			name:                   "with insecure skip secrets manager set",
			cloudInit:              infrav1.CloudInit{InsecureSkipSecretsManager: true},
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
			g.Expect(machine.Spec.CloudInit.SecureSecretsBackend).To(Equal(infrav1.SecretBackend(tt.expectedSecretsBackend)))
		})
	}
}

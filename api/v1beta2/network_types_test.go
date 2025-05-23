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
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"
)

func TestSGDifference(t *testing.T) {
	tests := []struct {
		name     string
		self     IngressRules
		input    IngressRules
		expected IngressRules
	}{
		{
			name:     "self and input are nil",
			self:     nil,
			input:    nil,
			expected: nil,
		},
		{
			name: "input is nil",
			self: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
			},
			input: nil,
			expected: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
			},
		},
		{
			name: "self has more rules",
			self: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
				{
					Description: "MY-SSH",
					Protocol:    SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
			},
			input: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
			},
			expected: IngressRules{
				{
					Description: "MY-SSH",
					Protocol:    SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			out := tc.self.Difference(tc.input)

			g.Expect(out).To(Equal(tc.expected))
		})
	}
}

var (
	stubNetworkTypeSubnetsAvailabilityZone = []*SubnetSpec{
		{
			ID:               "subnet-id-us-east-1a-private",
			AvailabilityZone: "us-east-1a",
			IsPublic:         false,
			ZoneType:         ptr.To(ZoneTypeAvailabilityZone),
		},
		{
			ID:               "subnet-id-us-east-1a-public",
			AvailabilityZone: "us-east-1a",
			IsPublic:         true,
			ZoneType:         ptr.To(ZoneTypeAvailabilityZone),
		},
	}
	stubNetworkTypeSubnetsLocalZone = []*SubnetSpec{
		{
			ID:               "subnet-id-us-east-1-nyc-1-private",
			AvailabilityZone: "us-east-1-nyc-1a",
			IsPublic:         false,
			ZoneType:         ptr.To(ZoneTypeLocalZone),
		},
		{
			ID:               "subnet-id-us-east-1-nyc-1-public",
			AvailabilityZone: "us-east-1-nyc-1a",
			IsPublic:         true,
			ZoneType:         ptr.To(ZoneTypeLocalZone),
		},
	}
	stubNetworkTypeSubnetsWavelengthZone = []*SubnetSpec{
		{
			ID:               "subnet-id-us-east-1-wl1-nyc-wlz-1-private",
			AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
			IsPublic:         false,
			ZoneType:         ptr.To(ZoneTypeWavelengthZone),
		},
		{
			ID:               "subnet-id-us-east-1-wl1-nyc-wlz-1-public",
			AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
			IsPublic:         true,
			ZoneType:         ptr.To(ZoneTypeWavelengthZone),
		},
	}

	subnetsAllZones = Subnets{
		{
			ResourceID:       "subnet-az-1a",
			AvailabilityZone: "us-east-1a",
		},
		{
			ResourceID:       "subnet-az-1b",
			IsPublic:         true,
			AvailabilityZone: "us-east-1a",
		},
		{
			ResourceID:       "subnet-az-2a",
			IsPublic:         false,
			AvailabilityZone: "us-east-1b",
		},
		{
			ResourceID:       "subnet-az-2b",
			IsPublic:         true,
			AvailabilityZone: "us-east-1b",
		},
		{
			ResourceID:       "subnet-az-3a",
			ZoneType:         ptr.To(ZoneTypeAvailabilityZone),
			IsPublic:         false,
			AvailabilityZone: "us-east-1c",
		},
		{
			ResourceID:       "subnet-az-3b",
			ZoneType:         ptr.To(ZoneTypeAvailabilityZone),
			IsPublic:         true,
			AvailabilityZone: "us-east-1c",
		},
		{
			ResourceID:       "subnet-lz-1a",
			ZoneType:         ptr.To(ZoneTypeLocalZone),
			IsPublic:         false,
			AvailabilityZone: "us-east-1-nyc-1a",
		},
		{
			ResourceID:       "subnet-lz-2b",
			ZoneType:         ptr.To(ZoneTypeLocalZone),
			IsPublic:         true,
			AvailabilityZone: "us-east-1-nyc-1a",
		},
		{
			ResourceID:       "subnet-wl-1a",
			ZoneType:         ptr.To(ZoneTypeWavelengthZone),
			IsPublic:         false,
			AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
		},
		{
			ResourceID:       "subnet-wl-1b",
			ZoneType:         ptr.To(ZoneTypeWavelengthZone),
			IsPublic:         true,
			AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
		},
	}
)

type testStubNetworkTypes struct{}

func (ts *testStubNetworkTypes) deepCopyToSubnets(stub []*SubnetSpec) (subnets Subnets) {
	for _, sn := range stub {
		subnets = append(subnets, *sn.DeepCopy())
	}
	return subnets
}

func (ts *testStubNetworkTypes) deepCopySubnets(stub []*SubnetSpec) (subnets []*SubnetSpec) {
	for _, s := range stub {
		subnets = append(subnets, s.DeepCopy())
	}
	return subnets
}

func (ts *testStubNetworkTypes) getSubnetsAvailabilityZones() (subnets []*SubnetSpec) {
	return ts.deepCopySubnets(stubNetworkTypeSubnetsAvailabilityZone)
}

func (ts *testStubNetworkTypes) getSubnetsLocalZones() (subnets []*SubnetSpec) {
	return ts.deepCopySubnets(stubNetworkTypeSubnetsLocalZone)
}

func (ts *testStubNetworkTypes) getSubnetsWavelengthZones() (subnets []*SubnetSpec) {
	return ts.deepCopySubnets(stubNetworkTypeSubnetsWavelengthZone)
}

func (ts *testStubNetworkTypes) getSubnets() (sns Subnets) {
	subnets := []*SubnetSpec{}
	subnets = append(subnets, ts.getSubnetsAvailabilityZones()...)
	subnets = append(subnets, ts.getSubnetsLocalZones()...)
	subnets = append(subnets, ts.getSubnetsWavelengthZones()...)
	sns = ts.deepCopyToSubnets(subnets)
	return sns
}

func TestSubnetSpec_IsEdge(t *testing.T) {
	stub := testStubNetworkTypes{}
	tests := []struct {
		name string
		spec *SubnetSpec
		want bool
	}{
		{
			name: "az without type is not edge",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsAvailabilityZones()[0]
				s.ZoneType = nil
				return s
			}(),
			want: false,
		},
		{
			name: "az is not edge",
			spec: stub.getSubnetsAvailabilityZones()[0],
			want: false,
		},
		{
			name: "localzone is edge",
			spec: stub.getSubnetsLocalZones()[0],
			want: true,
		},
		{
			name: "wavelength is edge",
			spec: stub.getSubnetsWavelengthZones()[0],
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.spec
			if got := s.IsEdge(); got != tt.want {
				t.Errorf("SubnetSpec.IsEdge() returned unexpected value = got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestSubnetSpec_IsEdgeWavelength(t *testing.T) {
	stub := testStubNetworkTypes{}
	tests := []struct {
		name string
		spec *SubnetSpec
		want bool
	}{
		{
			name: "az without type is not edge wavelength",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsAvailabilityZones()[0]
				s.ZoneType = nil
				return s
			}(),
			want: false,
		},
		{
			name: "az is not edge wavelength",
			spec: stub.getSubnetsAvailabilityZones()[0],
			want: false,
		},
		{
			name: "localzone is not edge wavelength",
			spec: stub.getSubnetsLocalZones()[0],
			want: false,
		},
		{
			name: "wavelength is edge wavelength",
			spec: stub.getSubnetsWavelengthZones()[0],
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.spec
			if got := s.IsEdgeWavelength(); got != tt.want {
				t.Errorf("SubnetSpec.IsEdgeWavelength() returned unexpected value = got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestSubnetSpec_SetZoneInfo(t *testing.T) {
	stub := testStubNetworkTypes{}
	tests := []struct {
		name    string
		spec    *SubnetSpec
		zones   []types.AvailabilityZone
		want    *SubnetSpec
		wantErr string
	}{
		{
			name: "set zone information to availability zone subnet",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsAvailabilityZones()[0]
				s.ZoneType = nil
				s.ParentZoneName = nil
				return s
			}(),
			zones: []types.AvailabilityZone{
				{
					ZoneName: ptr.To[string]("us-east-1a"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
			},
			want: stub.getSubnetsAvailabilityZones()[0],
		},
		{
			name: "set zone information to availability zone subnet with many zones",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsAvailabilityZones()[0]
				s.ZoneType = nil
				s.ParentZoneName = nil
				return s
			}(),
			zones: []types.AvailabilityZone{
				{
					ZoneName: ptr.To[string]("us-east-1b"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
				{
					ZoneName: ptr.To[string]("us-east-1a"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
			},
			want: stub.getSubnetsAvailabilityZones()[0],
		},
		{
			name: "want error when zone metadata is not provided",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsAvailabilityZones()[0]
				s.ZoneType = nil
				s.ParentZoneName = nil
				return s
			}(),
			zones:   []types.AvailabilityZone{},
			wantErr: `unable to update zone information for subnet 'subnet-id-us-east-1a-private' and zone 'us-east-1a'`,
		},
		{
			name: "want error when subnet's available zone is not set",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsAvailabilityZones()[0]
				s.AvailabilityZone = ""
				return s
			}(),
			zones: []types.AvailabilityZone{
				{
					ZoneName: ptr.To[string]("us-east-1a"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
			},
			wantErr: `unable to update zone information for subnet 'subnet-id-us-east-1a-private'`,
		},
		{
			name: "set zone information to local zone subnet",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsLocalZones()[0]
				s.ZoneType = nil
				s.ParentZoneName = nil
				return s
			}(),
			zones: []types.AvailabilityZone{
				{
					ZoneName: ptr.To[string]("us-east-1b"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
				{
					ZoneName: ptr.To[string]("us-east-1a"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
				{
					ZoneName: ptr.To[string]("us-east-1-nyc-1a"),
					ZoneType: ptr.To[string]("local-zone"),
				},
			},
			want: stub.getSubnetsLocalZones()[0],
		},
		{
			name: "set zone information to wavelength zone subnet",
			spec: func() *SubnetSpec {
				s := stub.getSubnetsWavelengthZones()[0]
				s.ZoneType = nil
				s.ParentZoneName = nil
				return s
			}(),
			zones: []types.AvailabilityZone{
				{
					ZoneName: ptr.To[string]("us-east-1b"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
				{
					ZoneName: ptr.To[string]("us-east-1a"),
					ZoneType: ptr.To[string]("availability-zone"),
				},
				{
					ZoneName: ptr.To[string]("us-east-1-wl1-nyc-wlz-1"),
					ZoneType: ptr.To[string]("wavelength-zone"),
				},
				{
					ZoneName: ptr.To[string]("us-east-1-nyc-1a"),
					ZoneType: ptr.To[string]("local-zone"),
				},
			},
			want: stub.getSubnetsWavelengthZones()[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.spec
			err := s.SetZoneInfo(tt.zones)
			if err != nil {
				if len(tt.wantErr) == 0 {
					t.Fatalf("SubnetSpec.SetZoneInfo() got unexpected error: %v", err)
				}
				if len(tt.wantErr) > 0 && err.Error() != tt.wantErr {
					t.Fatalf("SubnetSpec.SetZoneInfo() got unexpected error message:\n got: %v,\nwant: %v", err, tt.wantErr)
				} else {
					return
				}
			}
			if !cmp.Equal(s, tt.want) {
				t.Errorf("SubnetSpec.SetZoneInfo() got unwanted value:\n %v", cmp.Diff(s, tt.want))
			}
		})
	}
}

func TestSubnets_IDs(t *testing.T) {
	tests := []struct {
		name    string
		subnets Subnets
		want    []string
	}{
		{
			name:    "no valid subnet IDs",
			subnets: Subnets{},
			want:    []string{},
		},
		{
			name: "no valid subnet IDs",
			subnets: Subnets{
				{
					ResourceID: "subnet-lz-1",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
				},
				{
					ResourceID: "subnet-wl-1",
					ZoneType:   ptr.To(ZoneTypeWavelengthZone),
				},
			},
			want: []string{},
		},
		{
			name: "should have only subnet IDs from availability zone",
			subnets: Subnets{
				{
					ResourceID: "subnet-az-1",
				},
				{
					ResourceID: "subnet-az-2",
					ZoneType:   ptr.To(ZoneTypeAvailabilityZone),
				},
				{
					ResourceID: "subnet-lz-1",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
				},
			},
			want: []string{"subnet-az-1", "subnet-az-2"},
		},
		{
			name: "should have only subnet IDs from availability zone",
			subnets: Subnets{
				{
					ResourceID: "subnet-az-1",
				},
				{
					ResourceID: "subnet-az-2",
					ZoneType:   ptr.To(ZoneTypeAvailabilityZone),
				},
				{
					ResourceID: "subnet-lz-1",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
				},
				{
					ResourceID: "subnet-wl-1",
					ZoneType:   ptr.To(ZoneTypeWavelengthZone),
				},
			},
			want: []string{"subnet-az-1", "subnet-az-2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subnets.IDs(); !cmp.Equal(got, tt.want) {
				t.Errorf("Subnets.IDs() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSubnets_IDsWithEdge(t *testing.T) {
	tests := []struct {
		name    string
		subnets Subnets
		want    []string
	}{
		{
			name:    "invalid subnet IDs",
			subnets: nil,
			want:    []string{},
		},
		{
			name:    "invalid subnet IDs",
			subnets: Subnets{},
			want:    []string{},
		},
		{
			name: "subnet IDs for all zones",
			subnets: Subnets{
				{
					ResourceID: "subnet-az-1",
				},
				{
					ResourceID: "subnet-az-2",
					ZoneType:   ptr.To(ZoneTypeAvailabilityZone),
				},
				{
					ResourceID: "subnet-lz-1",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
				},
			},
			want: []string{"subnet-az-1", "subnet-az-2", "subnet-lz-1"},
		},
		{
			name: "subnet IDs for all zones",
			subnets: Subnets{
				{
					ResourceID: "subnet-az-1",
				},
				{
					ResourceID: "subnet-az-2",
					ZoneType:   ptr.To(ZoneTypeAvailabilityZone),
				},
				{
					ResourceID: "subnet-lz-1",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
				},
				{
					ResourceID: "subnet-wl-1",
					ZoneType:   ptr.To(ZoneTypeWavelengthZone),
				},
			},
			want: []string{"subnet-az-1", "subnet-az-2", "subnet-lz-1", "subnet-wl-1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subnets.IDsWithEdge(); !cmp.Equal(got, tt.want) {
				t.Errorf("Subnets.IDsWithEdge() got unwanted value:\n %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSubnets_FilterPrivate(t *testing.T) {
	tests := []struct {
		name    string
		subnets Subnets
		want    Subnets
	}{
		{
			name:    "no private subnets",
			subnets: nil,
			want:    nil,
		},
		{
			name:    "no private subnets",
			subnets: Subnets{},
			want:    nil,
		},
		{
			name: "no private subnets",
			subnets: Subnets{
				{
					ResourceID: "subnet-az-1b",
					IsPublic:   true,
				},
				{
					ResourceID: "subnet-az-2b",
					IsPublic:   true,
				},
				{
					ResourceID: "subnet-az-3b",
					ZoneType:   ptr.To(ZoneTypeAvailabilityZone),
					IsPublic:   true,
				},
				{
					ResourceID: "subnet-lz-1a",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
					IsPublic:   false,
				},
				{
					ResourceID: "subnet-lz-2b",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
					IsPublic:   true,
				},
			},
			want: nil,
		},
		{
			name:    "private subnets",
			subnets: subnetsAllZones,
			want: Subnets{
				{
					ResourceID:       "subnet-az-1a",
					AvailabilityZone: "us-east-1a",
				},
				{
					ResourceID:       "subnet-az-2a",
					IsPublic:         false,
					AvailabilityZone: "us-east-1b",
				},
				{
					ResourceID:       "subnet-az-3a",
					ZoneType:         ptr.To(ZoneTypeAvailabilityZone),
					IsPublic:         false,
					AvailabilityZone: "us-east-1c",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subnets.FilterPrivate(); !cmp.Equal(got, tt.want) {
				t.Errorf("Subnets.FilterPrivate() got unwanted value:\n %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSubnets_FilterPublic(t *testing.T) {
	tests := []struct {
		name    string
		subnets Subnets
		want    Subnets
	}{
		{
			name:    "empty subnets",
			subnets: nil,
			want:    nil,
		},
		{
			name:    "empty subnets",
			subnets: Subnets{},
			want:    nil,
		},
		{
			name: "no public subnets",
			subnets: Subnets{
				{
					ResourceID: "subnet-az-1a",
					IsPublic:   false,
				},
				{
					ResourceID: "subnet-az-2a",
					IsPublic:   false,
				},
				{
					ResourceID: "subnet-az-3a",
					ZoneType:   ptr.To(ZoneTypeAvailabilityZone),
					IsPublic:   false,
				},
				{
					ResourceID: "subnet-lz-1a",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
					IsPublic:   false,
				},
				{
					ResourceID: "subnet-lz-2b",
					ZoneType:   ptr.To(ZoneTypeLocalZone),
					IsPublic:   true,
				},
			},
			want: nil,
		},
		{
			name:    "public subnets",
			subnets: subnetsAllZones,
			want: Subnets{
				{
					ResourceID:       "subnet-az-1b",
					IsPublic:         true,
					AvailabilityZone: "us-east-1a",
				},
				{
					ResourceID:       "subnet-az-2b",
					IsPublic:         true,
					AvailabilityZone: "us-east-1b",
				},
				{
					ResourceID:       "subnet-az-3b",
					ZoneType:         ptr.To(ZoneTypeAvailabilityZone),
					IsPublic:         true,
					AvailabilityZone: "us-east-1c",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subnets.FilterPublic(); !cmp.Equal(got, tt.want) {
				t.Errorf("Subnets.FilterPublic() got unwanted value:\n %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSubnets_GetUniqueZones(t *testing.T) {
	tests := []struct {
		name    string
		subnets Subnets
		want    []string
	}{
		{
			name:    "no subnets",
			subnets: Subnets{},
			want:    []string{},
		},
		{
			name:    "all subnets and zones",
			subnets: subnetsAllZones,
			want: []string{
				"us-east-1a",
				"us-east-1b",
				"us-east-1c",
				"us-east-1-nyc-1a",
				"us-east-1-wl1-nyc-wlz-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subnets.GetUniqueZones(); !cmp.Equal(got, tt.want) {
				t.Errorf("Subnets.GetUniqueZones() got unwanted value:\n %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSubnets_HasPublicSubnetWavelength(t *testing.T) {
	stub := testStubNetworkTypes{}
	tests := []struct {
		name    string
		subnets Subnets
		want    bool
	}{
		{
			name:    "no subnets",
			subnets: Subnets{},
			want:    false,
		},
		{
			name:    "no wavelength",
			subnets: stub.deepCopyToSubnets(stub.getSubnetsAvailabilityZones()),
			want:    false,
		},
		{
			name:    "no wavelength",
			subnets: stub.deepCopyToSubnets(stub.getSubnetsLocalZones()),
			want:    false,
		},
		{
			name: "has only private subnets in wavelength zones",
			subnets: Subnets{
				{
					ID:               "subnet-id-us-east-1-wl1-nyc-wlz-1-private",
					AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
					IsPublic:         false,
					ZoneType:         ptr.To(ZoneTypeWavelengthZone),
				},
			},
			want: false,
		},
		{
			name:    "has public subnets in wavelength zones",
			subnets: stub.getSubnets(),
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subnets.HasPublicSubnetWavelength(); got != tt.want {
				t.Errorf("Subnets.HasPublicSubnetWavelength() got unwanted value:\n %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

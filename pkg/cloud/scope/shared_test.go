/*
Copyright 2021 The Kubernetes Authors.

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

package scope

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
)

func TestSubnetPlacement(t *testing.T) {
	testCases := []struct {
		name                string
		specSubnetIDs       []string
		specAZs             []string
		parentAZs           []string
		subnetPlacementType *expinfrav1.AZSubnetType
		controlPlaneSubnets infrav1.Subnets
		logger              *logger.Logger
		expectedSubnetIDs   []string
		expectError         bool
	}{
		{
			name:                "spec subnets expected",
			specSubnetIDs:       []string{"subnet-az1"},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: nil,
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az1"},
			expectError:       false,
		},
		{
			name:                "spec azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: nil,
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az2", "subnet-az3"},
			expectError:       false,
		},
		{
			name:                "parent azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: nil,
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az4", "subnet-az5"},
			expectError:       false,
		},
		{
			name:                "spec private azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePrivate),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az6",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
					Tags: infrav1.Tags{
						infrav1.NameAWSSubnetAssociation: infrav1.SecondarySubnetTagValue,
					},
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az3"},
			expectError:       false,
		},
		{
			name:                "spec public azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePublic),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az2"},
			expectError:       false,
		},
		{
			name:                "spec all azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypeAll),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az2", "subnet-az3"},
			expectError:       false,
		},
		{
			name:                "spec public no azs found",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1a"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePublic),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{},
			expectError:       true,
		},
		{
			name:                "parent private azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePrivate),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az5"},
			expectError:       false,
		},
		{
			name:                "parent public azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePublic),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az4"},
			expectError:       false,
		},
		{
			name:                "parent all azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypeAll),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az4", "subnet-az5"},
			expectError:       false,
		},
		{
			name:          "use control plane subnets",
			specSubnetIDs: []string{},
			specAZs:       []string{},
			parentAZs:     []string{},
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az1", "subnet-az3", "subnet-az5"},
			expectError:       false,
		},
		{
			name:                "no placement",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{},
			controlPlaneSubnets: infrav1.Subnets{},
			logger:              logger.NewLogger(klog.Background()),
			expectedSubnetIDs:   []string{},
			expectError:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			strategy, err := newDefaultSubnetPlacementStrategy(tc.logger)
			g.Expect(err).NotTo(HaveOccurred())

			actualSubnetIDs, err := strategy.Place(&placementInput{
				SpecSubnetIDs:           tc.specSubnetIDs,
				SpecAvailabilityZones:   tc.specAZs,
				ParentAvailabilityZones: tc.parentAZs,
				ControlplaneSubnets:     tc.controlPlaneSubnets,
				SubnetPlacementType:     tc.subnetPlacementType,
			})

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			}

			g.Expect(err).To(BeNil())
			g.Expect(actualSubnetIDs).To(Equal(tc.expectedSubnetIDs))
		})
	}
}

func TestClusterScopeStandardReady(t *testing.T) {
	for _, tt := range []struct {
		name                   string
		managed, bastion, ipv6 bool
		noPrivateSubnets       bool
		missing                string
		want                   metav1.ConditionStatus
	}{
		{name: "unmanaged", want: metav1.ConditionTrue},
		{name: "managed", managed: true, want: metav1.ConditionTrue},
		{name: "managed with bastion and IPv6", managed: true, bastion: true, ipv6: true, want: metav1.ConditionTrue},
		{name: "missing NAT gateway", managed: true, missing: infrav1.AWSClusterNatGatewaysReadyV1Beta2Condition, want: metav1.ConditionUnknown},
		{name: "missing optional bastion", managed: true, bastion: true, missing: infrav1.AWSClusterBastionHostReadyV1Beta2Condition, want: metav1.ConditionUnknown},
		{name: "bastion skipped without private subnets", managed: true, bastion: true, noPrivateSubnets: true, missing: infrav1.AWSClusterBastionHostReadyV1Beta2Condition, want: metav1.ConditionTrue},
		{name: "missing IPv6 gateway", managed: true, ipv6: true, missing: infrav1.AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition, want: metav1.ConditionUnknown},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			ctx := context.Background()
			scheme := runtime.NewScheme()
			g.Expect(infrav1.AddToScheme(scheme)).To(Succeed())
			obj := &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster", Namespace: "default", Generation: 2}}
			obj.Spec.Bastion.Enabled = tt.bastion
			obj.Spec.NetworkSpec.Subnets = infrav1.Subnets{{ID: "subnet", IsPublic: tt.noPrivateSubnets}}
			if !tt.managed {
				obj.Spec.NetworkSpec.VPC.ID = "vpc-external"
			}
			if tt.ipv6 {
				obj.Spec.NetworkSpec.VPC.IPv6 = &infrav1.IPv6{}
			}
			c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(obj).WithObjects(obj).Build()
			ph, err := v1beta1patch.NewHelper(obj, c)
			g.Expect(err).ToNot(HaveOccurred())
			s := &ClusterScope{AWSCluster: obj, Cluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster"}}, patchHelper: ph}
			types := []string{infrav1.AWSClusterVpcReadyV1Beta2Condition, infrav1.AWSClusterSubnetsReadyV1Beta2Condition, infrav1.AWSClusterSecurityGroupsReadyV1Beta2Condition, infrav1.AWSClusterLoadBalancerReadyV1Beta2Condition}
			if tt.managed {
				types = append(types, infrav1.AWSClusterInternetGatewayReadyV1Beta2Condition, infrav1.AWSClusterNatGatewaysReadyV1Beta2Condition, infrav1.AWSClusterRouteTablesReadyV1Beta2Condition, infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition)
			}
			if tt.bastion {
				types = append(types, infrav1.AWSClusterBastionHostReadyV1Beta2Condition)
			}
			if tt.ipv6 {
				types = append(types, infrav1.AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition)
			}
			for _, typ := range types {
				if tt.noPrivateSubnets && typ == infrav1.AWSClusterBastionHostReadyV1Beta2Condition {
					continue
				}
				legacyType := clusterv1beta1.ConditionType(typ)
				if typ == infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition {
					legacyType = infrav1.VpcEndpointsReadyCondition
				}
				v1beta1conditions.MarkTrue(obj, legacyType)
				if typ != tt.missing {
					conditions.Set(obj, metav1.Condition{Type: typ, Status: metav1.ConditionTrue, Reason: clusterv1.ReadyReason})
				}
			}
			// An unrelated condition written after the helper's snapshot must survive the patch.
			latest := &infrav1.AWSCluster{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), latest)).To(Succeed())
			conditions.Set(latest, metav1.Condition{Type: "OtherController", Status: metav1.ConditionFalse, Reason: "Observed"})
			g.Expect(c.Status().Update(ctx, latest)).To(Succeed())
			g.Expect(s.PatchObject()).To(Succeed())
			stored := &infrav1.AWSCluster{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
			ready := conditions.Get(stored, infrav1.AWSClusterReadyV1Beta2Condition)
			g.Expect(ready).ToNot(BeNil())
			g.Expect(ready.Status).To(Equal(tt.want))
			g.Expect(ready.ObservedGeneration).To(Equal(int64(2)))
			g.Expect(v1beta1conditions.IsTrue(stored, clusterv1beta1.ReadyCondition)).To(BeTrue())
			g.Expect(conditions.Get(stored, "OtherController").Status).To(Equal(metav1.ConditionFalse))
		})
	}
}

func TestPoolScopeStandardConditions(t *testing.T) {
	t.Run("waiting for ROSA control plane", func(t *testing.T) {
		g := NewWithT(t)
		ctx := context.Background()
		scheme := runtime.NewScheme()
		g.Expect(expinfrav1.AddToScheme(scheme)).To(Succeed())
		obj := &expinfrav1.ROSAMachinePool{ObjectMeta: metav1.ObjectMeta{Name: "pool", Namespace: "default"}}
		c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(obj).WithObjects(obj).Build()
		ph, err := v1beta1patch.NewHelper(obj, c)
		g.Expect(err).ToNot(HaveOccurred())
		s := &RosaMachinePoolScope{RosaMachinePool: obj, patchHelper: ph}
		g.Expect(s.RosaMachinePoolReadyFalse(expinfrav1.WaitingForRosaControlPlaneReason, "")).To(Succeed())
		stored := &expinfrav1.ROSAMachinePool{}
		g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
		g.Expect(v1beta1conditions.IsFalse(stored, expinfrav1.RosaMachinePoolReadyCondition)).To(BeTrue())
		standard := conditions.Get(stored, expinfrav1.ROSAMachinePoolReadyV1Beta2Condition)
		g.Expect(standard).ToNot(BeNil())
		g.Expect(standard.Status).To(Equal(metav1.ConditionFalse))
		g.Expect(standard.Reason).To(Equal(expinfrav1.ROSAMachinePoolWaitingForROSAControlPlaneV1Beta2Reason))
	})
	t.Run("deleting EKS nodegroup", func(t *testing.T) {
		g := NewWithT(t)
		ctx := context.Background()
		scheme := runtime.NewScheme()
		g.Expect(expinfrav1.AddToScheme(scheme)).To(Succeed())
		obj := &expinfrav1.AWSManagedMachinePool{ObjectMeta: metav1.ObjectMeta{Name: "pool", Namespace: "default"}}
		c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(obj).WithObjects(obj).Build()
		ph, err := v1beta1patch.NewHelper(obj, c)
		g.Expect(err).ToNot(HaveOccurred())
		s := &ManagedMachinePoolScope{ManagedMachinePool: obj, patchHelper: ph}
		g.Expect(s.NodegroupReadyFalse(clusterv1beta1.DeletingReason, "")).To(Succeed())
		g.Expect(s.IAMReadyFalse("DeletingFailed", "role in use")).To(Succeed())
		stored := &expinfrav1.AWSManagedMachinePool{}
		g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
		g.Expect(v1beta1conditions.IsFalse(stored, expinfrav1.EKSNodegroupReadyCondition)).To(BeTrue())
		nodegroup := conditions.Get(stored, expinfrav1.AWSManagedMachinePoolEKSNodegroupReadyV1Beta2Condition)
		g.Expect(nodegroup.Reason).To(Equal(clusterv1.DeletingReason))
		role := conditions.Get(stored, expinfrav1.AWSManagedMachinePoolIAMNodegroupRolesReadyV1Beta2Condition)
		g.Expect(role.Status).To(Equal(metav1.ConditionFalse))
		g.Expect(role.Reason).To(Equal("DeletingFailed"))
		g.Expect(role.Message).To(Equal("role in use"))
	})
}

func TestClusterScopeUnownedStandardConditionConflicts(t *testing.T) {
	for _, conditionType := range []string{
		infrav1.AWSClusterCarrierGatewayReadyV1Beta2Condition,
		infrav1.AWSClusterSecondaryCidrsReadyV1Beta2Condition,
		infrav1.AWSClusterS3BucketReadyV1Beta2Condition,
	} {
		t.Run(conditionType, func(t *testing.T) {
			g := NewWithT(t)
			ctx := context.Background()
			scheme := runtime.NewScheme()
			g.Expect(infrav1.AddToScheme(scheme)).To(Succeed())
			obj := &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster", Namespace: "default"}}
			c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(obj).WithObjects(obj).Build()
			ph, err := v1beta1patch.NewHelper(obj, c)
			g.Expect(err).ToNot(HaveOccurred())
			s := &ClusterScope{AWSCluster: obj, Cluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster"}}, patchHelper: ph}
			conditions.Set(obj, metav1.Condition{Type: conditionType, Status: metav1.ConditionTrue, Reason: clusterv1.ReadyReason})

			latest := &infrav1.AWSCluster{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), latest)).To(Succeed())
			conditions.Set(latest, metav1.Condition{Type: conditionType, Status: metav1.ConditionFalse, Reason: "ConcurrentFailure"})
			g.Expect(c.Status().Update(ctx, latest)).To(Succeed())

			g.Expect(s.PatchObject()).To(MatchError(ContainSubstring("conflict")))
			stored := &infrav1.AWSCluster{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
			g.Expect(conditions.Get(stored, conditionType).Reason).To(Equal("ConcurrentFailure"))
		})
	}
}

func TestBastionRequired(t *testing.T) {
	t.Run("nil bastion", func(t *testing.T) {
		g := NewWithT(t)
		g.Expect(BastionRequired(nil, infrav1.Subnets{{ID: "private"}})).To(BeFalse())
	})

	for _, tt := range []struct {
		name    string
		enabled bool
		subnets infrav1.Subnets
		want    bool
	}{
		{name: "disabled with private subnets", subnets: infrav1.Subnets{{ID: "private"}}},
		{name: "enabled without subnets", enabled: true},
		{name: "enabled with public subnets only", enabled: true, subnets: infrav1.Subnets{{ID: "public", IsPublic: true}}},
		{name: "enabled with private subnets", enabled: true, subnets: infrav1.Subnets{{ID: "private"}}, want: true},
		{name: "enabled with mixed subnets", enabled: true, subnets: infrav1.Subnets{{ID: "public", IsPublic: true}, {ID: "private"}}, want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(BastionRequired(&infrav1.Bastion{Enabled: tt.enabled}, tt.subnets)).To(Equal(tt.want))
		})
	}
}

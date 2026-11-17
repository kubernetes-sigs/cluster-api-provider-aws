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

package ec2

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/mock_services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
)

const pruneGateLTID = "lt-0123456789abcdef0"

// realPruneEC2 mocks every EC2Interface call except PruneLaunchTemplateVersions, which runs the real implementation.
type realPruneEC2 struct {
	*mock_services.MockEC2Interface
	real *Service
}

func (r realPruneEC2) PruneLaunchTemplateVersions(id string, inUseVersions []string) (*ec2types.LaunchTemplateVersion, error) {
	return r.real.PruneLaunchTemplateVersions(id, inUseVersions)
}

// ltVersions returns DescribeLaunchTemplateVersions output for the given version numbers (first is the default).
func ltVersions(versions ...int64) *ec2.DescribeLaunchTemplateVersionsOutput {
	out := &ec2.DescribeLaunchTemplateVersionsOutput{}
	for i, v := range versions {
		out.LaunchTemplateVersions = append(out.LaunchTemplateVersions, ec2types.LaunchTemplateVersion{
			LaunchTemplateId: aws.String(pruneGateLTID), VersionNumber: aws.Int64(v), DefaultVersion: aws.Bool(i == 0),
		})
	}
	return out
}

func newPruneGateMachinePoolScope(g *WithT, statusVersion string, bootstrapData []byte) (*scope.MachinePoolScope, *scope.ClusterScope, apimachinerytypes.NamespacedName) {
	amp := newAWSMachinePool()
	amp.Status.LaunchTemplateID = pruneGateLTID
	amp.Status.LaunchTemplateVersion = ptr.To(statusVersion)
	mp := newMachinePool()
	mp.Spec.Template.Spec.Bootstrap.DataSecretName = ptr.To("bootstrap-data")
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "bootstrap-data", Namespace: amp.Namespace},
		Data:       map[string][]byte{"value": bootstrapData},
	}

	scheme, err := setupScheme()
	g.Expect(err).NotTo(HaveOccurred())
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(amp, secret).WithStatusSubresource(amp).Build()
	cs, err := setupClusterScope(cl)
	g.Expect(err).NotTo(HaveOccurred())
	ms, err := scope.NewMachinePoolScope(scope.MachinePoolScopeParams{
		Client:         cl,
		InfraCluster:   cs,
		Cluster:        newCluster(),
		MachinePool:    mp,
		AWSMachinePool: amp,
	})
	g.Expect(err).NotTo(HaveOccurred())
	return ms, cs, apimachinerytypes.NamespacedName{Namespace: amp.Namespace, Name: "bootstrap-data"}
}

// expectNewAMIVersion mocks an existing LT whose only drift is a newly discovered AMI, so a new version is required.
func expectNewAMIVersion(ec2Svc *mock_services.MockEC2Interface, name string, bootstrapData []byte, secretKey apimachinerytypes.NamespacedName) {
	ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(
		&expinfrav1.AWSLaunchTemplate{Name: name, AMI: infrav1.AMIReference{ID: ptr.To("ami-old")}},
		userdata.ComputeHash(bootstrapData), &secretKey, nil, nil)
	ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any(), gomock.Any()).Return(ptr.To("ami-new"), nil)
	ec2Svc.EXPECT().LaunchTemplateNeedsUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, services.LaunchTemplateNeedsUpdateReasonNone, nil)
}

func reconcileWithPruneGate(s *Service, ms *scope.MachinePoolScope, cs *scope.ClusterScope, ec2Svc services.EC2Interface, gate func() (bool, []string, error)) error {
	_, err := s.ReconcileLaunchTemplate(context.TODO(), ms, ms, cs, ec2Svc, nil,
		func() (bool, *autoscalingtypes.InstanceRefreshStatus, error) { return true, nil, nil },
		func() error { return nil },
		func() error { return nil },
		gate,
	)
	return err
}

// TestReconcileLaunchTemplatePruneGate: versions [1, 6, 7], nodegroup stuck on v6, new AMI needs v8; v8 is always created.
func TestReconcileLaunchTemplatePruneGate(t *testing.T) {
	tests := []struct {
		name          string
		canPrune      bool
		inUse         []string
		gateErr       error
		expectPrune   bool
		expectInUseTo []string
	}{
		{name: "gate closed (consumer not on latest version): skip prune, still create new version", canPrune: false, inUse: []string{"6"}},
		{name: "gate open (consumer converged): prune with in-use versions, then create new version", canPrune: true, inUse: []string{"7"}, expectPrune: true, expectInUseTo: []string{"7"}},
		{name: "gate error: fail closed, skip prune, still create new version", gateErr: errors.New("describe nodegroup failed")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			bootstrapData := []byte("shell-script")
			ms, cs, secretKey := newPruneGateMachinePoolScope(g, "7", bootstrapData)

			ec2Svc := mock_services.NewMockEC2Interface(mockCtrl)
			expectNewAMIVersion(ec2Svc, ms.Name(), bootstrapData, secretKey)
			pruneTimes := 0
			if tc.expectPrune {
				pruneTimes = 1
			}
			ec2Svc.EXPECT().PruneLaunchTemplateVersions(pruneGateLTID, tc.expectInUseTo).Return(nil, nil).Times(pruneTimes)
			ec2Svc.EXPECT().CreateLaunchTemplateVersion(pruneGateLTID, gomock.Any(), gomock.Eq(ptr.To("ami-new")), secretKey, gomock.Any(), gomock.Any()).Return(nil)
			ec2Svc.EXPECT().GetLaunchTemplateLatestVersion(pruneGateLTID).Return("8", nil)

			err := reconcileWithPruneGate(NewService(cs), ms, cs, ec2Svc, func() (bool, []string, error) { return tc.canPrune, tc.inUse, tc.gateErr })
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(ms.GetLaunchTemplateLatestVersionStatus()).To(Equal("8"))
		})
	}
}

// TestReconcileLaunchTemplateStaleStatusKeepsInUseVersion: status stale at v6 lets the gate pass, but in-use v6 is not deleted.
func TestReconcileLaunchTemplateStaleStatusKeepsInUseVersion(t *testing.T) {
	g := NewWithT(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	bootstrapData := []byte("shell-script")
	ms, cs, secretKey := newPruneGateMachinePoolScope(g, "6", bootstrapData)

	ec2API := mocks.NewMockEC2API(mockCtrl)
	ec2API.EXPECT().DescribeLaunchTemplateVersions(gomock.Any(), gomock.Any()).Return(ltVersions(1, 6, 7), nil)
	ec2API.EXPECT().DeleteLaunchTemplateVersions(gomock.Any(), gomock.Any()).Times(0)
	s := NewService(cs)
	s.EC2Client = ec2API

	mockSvc := mock_services.NewMockEC2Interface(mockCtrl)
	expectNewAMIVersion(mockSvc, ms.Name(), bootstrapData, secretKey)
	mockSvc.EXPECT().CreateLaunchTemplateVersion(pruneGateLTID, gomock.Any(), gomock.Eq(ptr.To("ami-new")), secretKey, gomock.Any(), gomock.Any()).Return(nil)
	mockSvc.EXPECT().GetLaunchTemplateLatestVersion(pruneGateLTID).Return("8", nil)

	err := reconcileWithPruneGate(s, ms, cs, realPruneEC2{MockEC2Interface: mockSvc, real: s}, func() (bool, []string, error) { return true, []string{"6"}, nil })
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ms.GetLaunchTemplateLatestVersionStatus()).To(Equal("8"))
}

// TestPruneLaunchTemplateVersionsInUse: the prune candidate is never deleted while it is in use.
func TestPruneLaunchTemplateVersionsInUse(t *testing.T) {
	tests := []struct {
		name          string
		versions      *ec2.DescribeLaunchTemplateVersionsOutput
		inUse         []string
		expectDeleted string
	}{
		{name: "no in-use versions: oldest non-default version is pruned", versions: ltVersions(1, 6, 7), expectDeleted: "6"},
		{name: "candidate in use: nothing pruned", versions: ltVersions(1, 6, 7), inUse: []string{"6"}},
		{name: "only latest in use: candidate pruned", versions: ltVersions(1, 6, 7), inUse: []string{"7"}, expectDeleted: "6"},
		{name: "fewer than three versions: nothing pruned", versions: ltVersions(1, 6)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			cs, err := setupClusterScope(fake.NewClientBuilder().WithScheme(scheme).Build())
			g.Expect(err).NotTo(HaveOccurred())

			ec2API := mocks.NewMockEC2API(mockCtrl)
			ec2API.EXPECT().DescribeLaunchTemplateVersions(gomock.Any(), gomock.Any()).Return(tc.versions, nil)
			if tc.expectDeleted != "" {
				ec2API.EXPECT().DeleteLaunchTemplateVersions(gomock.Any(), gomock.Eq(&ec2.DeleteLaunchTemplateVersionsInput{
					LaunchTemplateId: aws.String(pruneGateLTID), Versions: []string{tc.expectDeleted},
				})).Return(&ec2.DeleteLaunchTemplateVersionsOutput{}, nil)
			}
			s := NewService(cs)
			s.EC2Client = ec2API

			pruned, err := s.PruneLaunchTemplateVersions(pruneGateLTID, tc.inUse)
			g.Expect(err).NotTo(HaveOccurred())
			if tc.expectDeleted == "" {
				g.Expect(pruned).To(BeNil())
				return
			}
			g.Expect(pruned).NotTo(BeNil())
			g.Expect(aws.ToInt64(pruned.VersionNumber)).To(BeNumerically("==", 6))
		})
	}
}

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

package eks

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	eksiam "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/iam"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth/mock_iamauth"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

const (
	testEKSClusterName = "test-cluster"
	testNodegroupName  = "test-ng"
)

type poolDeleteMocks struct {
	eks *mock_eksiface.MockEKSAPI
	iam *mock_iamauth.MockIAMAPI
}

func newPoolDeleteService(t *testing.T, enableIAM bool, roleName string) (*NodegroupService, *poolDeleteMocks, *expinfrav1.AWSManagedMachinePool) {
	t.Helper()
	g := NewWithT(t)

	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	m := &poolDeleteMocks{
		eks: mock_eksiface.NewMockEKSAPI(mockCtrl),
		iam: mock_iamauth.NewMockIAMAPI(mockCtrl),
	}

	scheme := runtime.NewScheme()
	g.Expect(clusterv1.AddToScheme(scheme)).To(Succeed())
	g.Expect(expinfrav1.AddToScheme(scheme)).To(Succeed())
	g.Expect(ekscontrolplanev1.AddToScheme(scheme)).To(Succeed())

	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{
		ObjectMeta: metav1.ObjectMeta{Name: "cp", Namespace: "default"},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
			EKSClusterName: testEKSClusterName,
			Region:         "us-west-2",
		},
	}
	machinePool := &clusterv1.MachinePool{ObjectMeta: metav1.ObjectMeta{Name: "mp", Namespace: "default"}}
	pool := &expinfrav1.AWSManagedMachinePool{
		ObjectMeta: metav1.ObjectMeta{Name: "mp", Namespace: "default"},
		Spec: expinfrav1.AWSManagedMachinePoolSpec{
			EKSNodegroupName: testNodegroupName,
			RoleName:         roleName,
		},
	}
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(controlPlane, machinePool, pool).WithStatusSubresource(pool, machinePool).Build()

	log := logger.NewLogger(logr.Discard())
	machinePoolScope, err := scope.NewManagedMachinePoolScope(scope.ManagedMachinePoolScopeParams{
		Client:             cl,
		Logger:             log,
		Cluster:            &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "c", Namespace: "default"}},
		ControlPlane:       controlPlane,
		MachinePool:        machinePool,
		ManagedMachinePool: pool,
		EnableIAM:          enableIAM,
	})
	g.Expect(err).NotTo(HaveOccurred())

	return &NodegroupService{
		scope:      machinePoolScope,
		EKSClient:  m.eks,
		IAMService: eksiam.IAMService{Wrapper: log, IAMClient: m.iam},
	}, m, pool
}

func nodegroupWithStatus(status ekstypes.NodegroupStatus) *eks.DescribeNodegroupOutput {
	return &eks.DescribeNodegroupOutput{Nodegroup: &ekstypes.Nodegroup{
		NodegroupName: aws.String(testNodegroupName),
		Status:        status,
	}}
}

func ownedRole(name string) *iam.GetRoleOutput {
	return &iam.GetRoleOutput{Role: &iamtypes.Role{
		RoleName: aws.String(name),
		Tags: []iamtypes.Tag{{
			Key:   aws.String(infrav1.ClusterAWSCloudProviderTagKey(testEKSClusterName)),
			Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
		}},
	}}
}

func unmanagedRole(name string) *iam.GetRoleOutput {
	return &iam.GetRoleOutput{Role: &iamtypes.Role{RoleName: aws.String(name)}}
}

var (
	errNoSuchEntity = &smithy.GenericAPIError{Code: "NoSuchEntity", Message: "role not found"}
	errAccessDenied = &smithy.GenericAPIError{Code: "AccessDenied", Message: "denied"}
)

func expectDescribe(m *poolDeleteMocks, out *eks.DescribeNodegroupOutput, err error) *gomock.Call {
	return m.eks.EXPECT().DescribeNodegroup(gomock.Any(), gomock.Any()).Return(out, err)
}

func expectDescribeNotFound(m *poolDeleteMocks) *gomock.Call {
	return expectDescribe(m, nil, &ekstypes.ResourceNotFoundException{})
}

func expectDeleteNodegroup(m *poolDeleteMocks, err error) *gomock.Call {
	return m.eks.EXPECT().DeleteNodegroup(gomock.Any(), gomock.Any()).Return(&eks.DeleteNodegroupOutput{}, err)
}

func expectWait(m *poolDeleteMocks, err error) *gomock.Call {
	return m.eks.EXPECT().WaitUntilNodegroupDeleted(gomock.Any(), gomock.Any(), gomock.Any()).Return(err)
}

func expectGetRole(m *poolDeleteMocks, wantName string, out *iam.GetRoleOutput, err error) *gomock.Call {
	return m.iam.EXPECT().GetRole(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, in *iam.GetRoleInput, _ ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
			if aws.ToString(in.RoleName) != wantName {
				return nil, errors.Errorf("GetRole called with %q, want %q", aws.ToString(in.RoleName), wantName)
			}
			return out, err
		})
}

// expectRoleDeleted expects the detach-all-policies + DeleteRole sequence.
func expectRoleDeleted(m *poolDeleteMocks, name string) []*gomock.Call {
	return []*gomock.Call{
		m.iam.EXPECT().ListAttachedRolePolicies(gomock.Any(), gomock.Any()).Return(&iam.ListAttachedRolePoliciesOutput{}, nil),
		m.iam.EXPECT().DeleteRole(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, in *iam.DeleteRoleInput, _ ...func(*iam.Options)) (*iam.DeleteRoleOutput, error) {
				if aws.ToString(in.RoleName) != name {
					return nil, errors.Errorf("DeleteRole called with %q, want %q", aws.ToString(in.RoleName), name)
				}
				return &iam.DeleteRoleOutput{}, nil
			}),
	}
}

func TestReconcilePoolDelete(t *testing.T) {
	const explicitRole = "my-nodegroup-role"

	tests := []struct {
		name      string
		enableIAM bool
		roleName  string
		expect    func(m *poolDeleteMocks)
		wantErr   string
	}{
		{
			name:      "nodegroup already gone, owned IAM role is still deleted",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribeNotFound(m)
				gomock.InOrder(append([]*gomock.Call{expectGetRole(m, explicitRole, ownedRole(explicitRole), nil)}, expectRoleDeleted(m, explicitRole)...)...)
			},
		},
		{
			name:      "describe returns an empty nodegroup, owned IAM role is still deleted",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribe(m, &eks.DescribeNodegroupOutput{}, nil)
				gomock.InOrder(append([]*gomock.Call{expectGetRole(m, explicitRole, ownedRole(explicitRole), nil)}, expectRoleDeleted(m, explicitRole)...)...)
			},
		},
		{
			name:      "nodegroup gone and RoleName empty: no IAM calls, no error",
			enableIAM: true,
			roleName:  "",
			expect: func(m *poolDeleteMocks) {
				expectDescribeNotFound(m)
			},
		},
		{
			name:      "nodegroup gone, IAM role already deleted",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribeNotFound(m)
				expectGetRole(m, explicitRole, nil, errNoSuchEntity)
			},
		},
		{
			name:      "nodegroup gone, unmanaged (BYO / bootstrap-stack) role is left alone",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribeNotFound(m)
				expectGetRole(m, explicitRole, unmanagedRole(explicitRole), nil)
			},
		},
		{
			name:      "nodegroup gone, IAM management disabled: no IAM calls",
			enableIAM: false,
			roleName:  "",
			expect: func(m *poolDeleteMocks) {
				expectDescribeNotFound(m)
			},
		},
		{
			name:      "nodegroup active: delete, wait, then delete role in that order",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				calls := []*gomock.Call{
					expectDescribe(m, nodegroupWithStatus(ekstypes.NodegroupStatusActive), nil),
					expectDeleteNodegroup(m, nil),
					expectWait(m, nil),
					expectGetRole(m, explicitRole, ownedRole(explicitRole), nil),
				}
				gomock.InOrder(append(calls, expectRoleDeleted(m, explicitRole)...)...)
			},
		},
		{
			name:      "nodegroup already DELETING: only wait, then delete role",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				calls := []*gomock.Call{
					expectDescribe(m, nodegroupWithStatus(ekstypes.NodegroupStatusDeleting), nil),
					expectWait(m, nil),
					expectGetRole(m, explicitRole, ownedRole(explicitRole), nil),
				}
				gomock.InOrder(append(calls, expectRoleDeleted(m, explicitRole)...)...)
			},
		},
		{
			name:      "nodegroup in DELETE_FAILED: delete is retried",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				calls := []*gomock.Call{
					expectDescribe(m, nodegroupWithStatus(ekstypes.NodegroupStatusDeleteFailed), nil),
					expectDeleteNodegroup(m, nil),
					expectWait(m, nil),
					expectGetRole(m, explicitRole, ownedRole(explicitRole), nil),
				}
				gomock.InOrder(append(calls, expectRoleDeleted(m, explicitRole)...)...)
			},
		},
		{
			name:      "DeleteNodegroup reports NotFound (raced with deletion): role is still deleted, no wait",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				calls := []*gomock.Call{
					expectDescribe(m, nodegroupWithStatus(ekstypes.NodegroupStatusActive), nil),
					expectDeleteNodegroup(m, &ekstypes.ResourceNotFoundException{}),
					expectGetRole(m, explicitRole, ownedRole(explicitRole), nil),
				}
				gomock.InOrder(append(calls, expectRoleDeleted(m, explicitRole)...)...)
			},
		},
		{
			name:      "DeleteNodegroup fails: error returned, role untouched",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribe(m, nodegroupWithStatus(ekstypes.NodegroupStatusActive), nil)
				expectDeleteNodegroup(m, errors.New("boom"))
			},
			wantErr: "failed to delete nodegroup",
		},
		{
			name:      "wait for nodegroup deletion fails: error returned, role untouched so a retry can finish",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribe(m, nodegroupWithStatus(ekstypes.NodegroupStatusActive), nil)
				expectDeleteNodegroup(m, nil)
				expectWait(m, errors.New("timed out"))
			},
			wantErr: "failed waiting for EKS nodegroup",
		},
		{
			name:      "describe fails with a non-NotFound error: nothing is deleted",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribe(m, nil, errAccessDenied)
			},
			wantErr: "failed to describe",
		},
		{
			name:      "GetRole fails with a non-NotFound error: error returned",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribeNotFound(m)
				expectGetRole(m, explicitRole, nil, errAccessDenied)
			},
			wantErr: "failed to delete nodegroup IAM role",
		},
		{
			name:      "DeleteRole fails: error returned so the finalizer is kept and the reconcile retries",
			enableIAM: true,
			roleName:  explicitRole,
			expect: func(m *poolDeleteMocks) {
				expectDescribeNotFound(m)
				expectGetRole(m, explicitRole, ownedRole(explicitRole), nil)
				m.iam.EXPECT().ListAttachedRolePolicies(gomock.Any(), gomock.Any()).Return(&iam.ListAttachedRolePoliciesOutput{}, nil)
				m.iam.EXPECT().DeleteRole(gomock.Any(), gomock.Any()).Return(nil, errAccessDenied)
			},
			wantErr: "failed to delete nodegroup IAM role",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			s, m, pool := newPoolDeleteService(t, tc.enableIAM, tc.roleName)
			tc.expect(m)

			err := s.ReconcilePoolDelete(context.TODO())
			if tc.wantErr != "" {
				g.Expect(err).To(MatchError(ContainSubstring(tc.wantErr)))
				return
			}
			g.Expect(err).NotTo(HaveOccurred())

			cond := v1beta1conditions.Get(pool, expinfrav1.EKSNodegroupReadyCondition)
			g.Expect(cond).NotTo(BeNil())
			g.Expect(cond.Status).To(Equal(corev1.ConditionFalse))
			g.Expect(cond.Reason).To(Equal(clusterv1beta1.DeletedReason))
		})
	}
}

// TestReconcilePoolDeleteResumesAfterInterruption: the first reconcile is interrupted
// before the IAM step, EKS finishes the deletion, and the next reconcile must remove the role.
func TestReconcilePoolDeleteResumesAfterInterruption(t *testing.T) {
	g := NewWithT(t)
	const role = "my-nodegroup-role"
	s, m, _ := newPoolDeleteService(t, true, role)

	gomock.InOrder(
		expectDescribe(m, nodegroupWithStatus(ekstypes.NodegroupStatusActive), nil),
		expectDeleteNodegroup(m, nil),
		expectWait(m, errors.New("timed out")),
	)
	g.Expect(s.ReconcilePoolDelete(context.TODO())).NotTo(Succeed())

	calls := []*gomock.Call{
		expectDescribeNotFound(m),
		expectGetRole(m, role, ownedRole(role), nil),
	}
	gomock.InOrder(append(calls, expectRoleDeleted(m, role)...)...)
	g.Expect(s.ReconcilePoolDelete(context.TODO())).To(Succeed())
}

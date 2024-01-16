package eks

import (
	"github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"testing"
)

func TestSetStatus(t *testing.T) {
	g := NewWithT(t)
	degraded := eks.NodegroupStatusDegraded
	code := eks.NodegroupIssueCodeAsgInstanceLaunchFailures
	message := "VcpuLimitExceeded"
	resourceId := "my-worker-nodes"

	s := &NodegroupService{
		scope: &scope.ManagedMachinePoolScope{
			ManagedMachinePool: &v1beta1.AWSManagedMachinePool{
				Status: v1beta1.AWSManagedMachinePoolStatus{
					Ready: false,
				},
			},
		},
	}

	issue := &eks.Issue{
		Code:        &code,
		Message:     &message,
		ResourceIds: []*string{&resourceId},
	}
	ng := &eks.Nodegroup{
		Status: &degraded,
		Health: &eks.NodegroupHealth{
			Issues: []*eks.Issue{issue},
		},
	}

	err := s.setStatus(ng)
	g.Expect(err).ToNot(BeNil())
	// ensure machine pool status values are set as expected
	g.Expect(*s.scope.ManagedMachinePool.Status.FailureMessage).To(ContainSubstring(issue.GoString()))
	g.Expect(s.scope.ManagedMachinePool.Status.Ready).To(Equal(false))
	g.Expect(*s.scope.ManagedMachinePool.Status.FailureReason).To(Equal(capierrors.InsufficientResourcesMachineError))
}

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

package rosa

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
)

func buildOCMAutoscaler(g *WithT, maxNodesTotal, maxPodGracePeriod int, maxNodeProvisionTime string, podPriorityThreshold int) *cmv1.ClusterAutoscaler {
	autoscaler, err := cmv1.NewClusterAutoscaler().
		ResourceLimits(cmv1.NewAutoscalerResourceLimits().MaxNodesTotal(maxNodesTotal)).
		MaxPodGracePeriod(maxPodGracePeriod).
		MaxNodeProvisionTime(maxNodeProvisionTime).
		PodPriorityThreshold(podPriorityThreshold).
		Build()
	g.Expect(err).ToNot(HaveOccurred())
	return autoscaler
}

func TestBuildClusterAutoscalerPatch(t *testing.T) {
	t.Run("returns nil when OCM already has the defaults and spec is empty", func(t *testing.T) {
		g := NewWithT(t)
		current := buildOCMAutoscaler(g, 0, 600, "15m", -10)

		patch, err := BuildClusterAutoscalerPatch(&rosacontrolplanev1.ClusterAutoscaler{}, current)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(patch).To(BeNil())
	})

	t.Run("sends only the field that differs", func(t *testing.T) {
		g := NewWithT(t)
		current := buildOCMAutoscaler(g, 0, 600, "15m", -10)
		spec := &rosacontrolplanev1.ClusterAutoscaler{
			ResourceLimits: &rosacontrolplanev1.ClusterAutoscalerResourceLimits{MaxNodesTotal: ptr.To[int32](50)},
		}

		patch, err := BuildClusterAutoscalerPatch(spec, current)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(patch).ToNot(BeNil())

		maxNodesTotal, ok := patch.ResourceLimits().GetMaxNodesTotal()
		g.Expect(ok).To(BeTrue())
		g.Expect(maxNodesTotal).To(Equal(50))
		_, ok = patch.GetMaxPodGracePeriod()
		g.Expect(ok).To(BeFalse())
		_, ok = patch.GetMaxNodeProvisionTime()
		g.Expect(ok).To(BeFalse())
		_, ok = patch.GetPodPriorityThreshold()
		g.Expect(ok).To(BeFalse())
	})

	t.Run("resets values set outside of CAPA to the defaults when spec is empty", func(t *testing.T) {
		g := NewWithT(t)
		current := buildOCMAutoscaler(g, 50, 900, "30m", 10)

		patch, err := BuildClusterAutoscalerPatch(&rosacontrolplanev1.ClusterAutoscaler{}, current)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(patch).ToNot(BeNil())
		g.Expect(patch.ResourceLimits().MaxNodesTotal()).To(Equal(0))
		g.Expect(patch.MaxPodGracePeriod()).To(Equal(600))
		g.Expect(patch.MaxNodeProvisionTime()).To(Equal("15m"))
		g.Expect(patch.PodPriorityThreshold()).To(Equal(-10))
	})

	t.Run("compares durations by value, not by string", func(t *testing.T) {
		g := NewWithT(t)
		current := buildOCMAutoscaler(g, 0, 600, "15m0s", -10)

		patch, err := BuildClusterAutoscalerPatch(&rosacontrolplanev1.ClusterAutoscaler{
			MaxNodeProvisionTime: &metav1.Duration{Duration: 15 * time.Minute},
		}, current)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(patch).To(BeNil())
	})

	t.Run("sends durations in compact form", func(t *testing.T) {
		g := NewWithT(t)
		current := buildOCMAutoscaler(g, 0, 600, "15m", -10)

		patch, err := BuildClusterAutoscalerPatch(&rosacontrolplanev1.ClusterAutoscaler{
			MaxNodeProvisionTime: &metav1.Duration{Duration: 30 * time.Minute},
		}, current)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(patch.MaxNodeProvisionTime()).To(Equal("30m"))
	})

	t.Run("treats fields missing in OCM as different", func(t *testing.T) {
		g := NewWithT(t)
		current, err := cmv1.NewClusterAutoscaler().Build()
		g.Expect(err).ToNot(HaveOccurred())

		patch, err := BuildClusterAutoscalerPatch(&rosacontrolplanev1.ClusterAutoscaler{}, current)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(patch).ToNot(BeNil())
		_, ok := patch.ResourceLimits().GetMaxNodesTotal()
		g.Expect(ok).To(BeTrue())
		_, ok = patch.GetMaxPodGracePeriod()
		g.Expect(ok).To(BeTrue())
		_, ok = patch.GetMaxNodeProvisionTime()
		g.Expect(ok).To(BeTrue())
		_, ok = patch.GetPodPriorityThreshold()
		g.Expect(ok).To(BeTrue())
	})

	t.Run("handles a nil current autoscaler", func(t *testing.T) {
		g := NewWithT(t)
		patch, err := BuildClusterAutoscalerPatch(&rosacontrolplanev1.ClusterAutoscaler{
			PodPriorityThreshold: ptr.To[int32](5),
		}, nil)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(patch).ToNot(BeNil())
		g.Expect(patch.PodPriorityThreshold()).To(Equal(5))
	})
}

func TestClusterAutoscalerStatusFromOCM(t *testing.T) {
	t.Run("returns nil for a nil autoscaler", func(t *testing.T) {
		g := NewWithT(t)
		g.Expect(ClusterAutoscalerStatusFromOCM(nil)).To(BeNil())
	})

	t.Run("converts all fields", func(t *testing.T) {
		g := NewWithT(t)
		status := ClusterAutoscalerStatusFromOCM(buildOCMAutoscaler(g, 50, 600, "15m", -10))
		g.Expect(status).To(Equal(&rosacontrolplanev1.ClusterAutoscalerStatus{
			ResourceLimits:       &rosacontrolplanev1.ClusterAutoscalerResourceLimitsStatus{MaxNodesTotal: ptr.To[int32](50)},
			MaxPodGracePeriod:    ptr.To[int32](600),
			MaxNodeProvisionTime: &metav1.Duration{Duration: 15 * time.Minute},
			PodPriorityThreshold: ptr.To[int32](-10),
		}))
	})

	t.Run("leaves missing or unparsable fields unset", func(t *testing.T) {
		g := NewWithT(t)
		current, err := cmv1.NewClusterAutoscaler().MaxNodeProvisionTime("not-a-duration").Build()
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ClusterAutoscalerStatusFromOCM(current)).To(Equal(&rosacontrolplanev1.ClusterAutoscalerStatus{}))
	})
}

func TestClusterAutoscalerPatchFields(t *testing.T) {
	g := NewWithT(t)

	empty, err := cmv1.NewClusterAutoscaler().Build()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(ClusterAutoscalerPatchFields(empty)).To(BeEmpty())
	g.Expect(ClusterAutoscalerPatchFields(nil)).To(BeEmpty())

	patch, err := BuildClusterAutoscalerPatch(&rosacontrolplanev1.ClusterAutoscaler{
		ResourceLimits:       &rosacontrolplanev1.ClusterAutoscalerResourceLimits{MaxNodesTotal: ptr.To[int32](50)},
		PodPriorityThreshold: ptr.To[int32](5),
	}, buildOCMAutoscaler(g, 0, 600, "15m", -10))
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(ClusterAutoscalerPatchFields(patch)).To(Equal([]string{"resourceLimits.maxNodesTotal", "podPriorityThreshold"}))
}

func TestFormatDuration(t *testing.T) {
	g := NewWithT(t)
	g.Expect(formatDuration(15 * time.Minute)).To(Equal("15m"))
	g.Expect(formatDuration(time.Hour)).To(Equal("1h"))
	g.Expect(formatDuration(90 * time.Minute)).To(Equal("1h30m"))
	g.Expect(formatDuration(90 * time.Second)).To(Equal("1m30s"))
	g.Expect(formatDuration(45 * time.Second)).To(Equal("45s"))
}

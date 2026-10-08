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
	"testing"

	. "github.com/onsi/gomega"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestROSAMachinePoolEc2MetadataHTTPTokensImmutability(t *testing.T) {
	w := &ROSAMachinePool{}

	t.Run("changing ec2MetadataHttpTokens is rejected", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensOptional,
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("ec2MetadataHttpTokens"))
	})

	t.Run("adding ec2MetadataHttpTokens after creation is rejected", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-pool",
				InstanceType: "m5.large",
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("ec2MetadataHttpTokens"))
	})

	t.Run("same ec2MetadataHttpTokens value is allowed", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("both empty ec2MetadataHttpTokens is allowed", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-pool",
				InstanceType: "m5.large",
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-pool",
				InstanceType: "m5.large",
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).ToNot(HaveOccurred())
	})
}

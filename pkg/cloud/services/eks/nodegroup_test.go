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
	"testing"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestNodegroupAMIType(t *testing.T) {
	al2 := expinfrav1.Al2x86_64
	al2023 := expinfrav1.Al2023x86_64
	custom := expinfrav1.Custom
	lookupAL2023 := infrav1.AmazonLinux2023

	tests := []struct {
		name   string
		pool   expinfrav1.AWSManagedMachinePoolSpec
		expect ekstypes.AMITypes
	}{
		{
			name:   "no launch template, no AMI type",
			pool:   expinfrav1.AWSManagedMachinePoolSpec{},
			expect: "",
		},
		{
			name:   "no launch template uses the AMI type",
			pool:   expinfrav1.AWSManagedMachinePoolSpec{AMIType: &al2023},
			expect: ekstypes.AMITypesAl2023X8664Standard,
		},
		{
			name: "launch template with looked-up AMI omits the defaulted AMI type",
			pool: expinfrav1.AWSManagedMachinePoolSpec{
				AMIType: &al2,
				AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
					AMI: infrav1.AMIReference{EKSOptimizedLookupType: &lookupAL2023},
				},
			},
			expect: "",
		},
		{
			name: "launch template with explicit AMI ID omits the AMI type",
			pool: expinfrav1.AWSManagedMachinePoolSpec{
				AMIType: &al2,
				AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
					AMI: infrav1.AMIReference{ID: ptr.To("ami-123")},
				},
			},
			expect: "",
		},
		{
			name: "launch template with AMI type CUSTOM omits it",
			pool: expinfrav1.AWSManagedMachinePoolSpec{
				AMIType:           &custom,
				AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{},
			},
			expect: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(nodegroupAMIType(&tt.pool)).To(Equal(tt.expect))
		})
	}
}

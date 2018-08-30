// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster_test

import (
	"testing"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type ec2 struct{}

func (e *ec2) ReconcileVPC(id string) (*ec2svc.VPC, error) {
	return &ec2svc.VPC{
		ID: id,
	}, nil
}

func TestReconcile(t *testing.T) {
	c, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create codec: %v", err)
	}
	ap := cluster.ActuatorParams{
		Codec:      c,
		EC2Service: &ec2{},
	}

	a, err := cluster.NewActuator(ap)
	if err != nil {
		t.Fatalf("could not create an actuator: %v", err)
	}

	if err := a.Reconcile(&clusterv1.Cluster{}); err != nil {
		t.Fatalf("failed to reconcile cluster: %v", err)
	}
}

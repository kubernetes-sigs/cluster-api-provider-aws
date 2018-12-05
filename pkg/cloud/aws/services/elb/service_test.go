/*
Copyright 2018 The Kubernetes Authors.

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

package elb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	//nolint
)

func TestNewService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
	elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

	scope, err := actuators.NewScope(actuators.ScopeParams{
		Cluster: &clusterv1.Cluster{},
		AWSClients: actuators.AWSClients{
			EC2: ec2Mock,
			ELB: elbMock,
		},
	})

	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}

	s := NewService(scope)
	if s == nil {
		t.Fatalf("Service shouldn't be nil")
	}
}

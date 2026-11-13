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

package identityprovider

import (
	"context"
	"testing"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"k8s.io/klog/v2"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

func TestMultiProviderPlan(t *testing.T) {
	clusterName := "default.cluster"
	g := NewWithT(t)

	mockControl := gomock.NewController(t)
	defer mockControl.Finish()

	eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
	eksMock.EXPECT().AssociateIdentityProviderConfig(gomock.Eq(context.TODO()), gomock.Any()).Times(1)

	currentByName := map[string]*OidcIdentityProviderConfig{
		"existing": createCurrentIdentityProvider("existing", "arn:existing", string(ekstypes.ConfigStatusActive), createTags()),
	}
	desiredByName := map[string]*OidcIdentityProviderConfig{
		"existing": createDesiredIdentityProvider("existing", createTags()),
		"new":      createDesiredIdentityProvider("new", createTags()),
	}

	planner := NewMultiProviderPlan(clusterName, currentByName, desiredByName, eksMock, logger.NewLogger(klog.Background()))
	procedures, err := planner.Create(context.TODO())
	g.Expect(err).NotTo(HaveOccurred())

	associateCount := 0
	for _, proc := range procedures {
		if proc.Name() == procedureNameAssociateIdentityProvider {
			associateCount++
		}
	}
	g.Expect(associateCount).To(Equal(1))

	for _, proc := range procedures {
		if proc.Name() == procedureNameAssociateIdentityProvider {
			g.Expect(proc.Do(context.TODO())).To(Succeed())
		}
	}
}

func TestUnionProviderNames(t *testing.T) {
	g := NewWithT(t)

	current := map[string]*OidcIdentityProviderConfig{"a": {}, "b": {}}
	desired := map[string]*OidcIdentityProviderConfig{"b": {}, "c": {}}
	names := unionProviderNames(current, desired)
	g.Expect(names).To(ConsistOf("a", "b", "c"))
}

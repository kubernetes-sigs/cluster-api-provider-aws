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

package integration_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"

	appsv1 "k8s.io/api/apps/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
)

const (
	controllerNamespace = "aws-provider-system"
	controllerName      = "aws-provider-controller-manager"
)

var _ = Describe("Metacluster", func() {
	Describe("manager container", func() {
		It("Should be healthy", func() {
			Eventually(
				func() (*appsv1.StatefulSet, error) {
					statefulSet := &appsv1.StatefulSet{}
					if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: controllerNamespace, Name: controllerName}, statefulSet); err != nil {
						return nil, err
					}
					return statefulSet, nil
				},
				2*time.Minute, 5*time.Second,
			).Should(haveReplicas(1))
		})
	})
})

// haveReplicas matches a stateful set with i replicas
func haveReplicas(i int32) types.GomegaMatcher {
	return PointTo(
		MatchFields(IgnoreExtras, Fields{
			"Status": MatchFields(IgnoreExtras, Fields{
				"Replicas":      Equal(i),
				"ReadyReplicas": Equal(i),
			}),
		}),
	)
}

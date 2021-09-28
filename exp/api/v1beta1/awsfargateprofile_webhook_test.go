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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/eks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utildefaulting "sigs.k8s.io/cluster-api/util/defaulting"
)

func TestAWSFargateProfileDefault(t *testing.T) {
	fargate := &AWSFargateProfile{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
		Spec: FargateProfileSpec{
			ClusterName: "clustername",
		},
	}
	t.Run("for AWSFargateProfile", utildefaulting.DefaultValidateTest(fargate))
	fargate.Default()
	g := NewWithT(t)
	g.Expect(fargate.GetLabels()[clusterv1.ClusterLabelName]).To(BeEquivalentTo(fargate.Spec.ClusterName))
	name, err := eks.GenerateEKSName(fargate.Name, fargate.Namespace, maxProfileNameLength)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fargate.Spec.ProfileName).To(BeEquivalentTo(name))
}

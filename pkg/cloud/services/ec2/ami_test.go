/*
Copyright 2019 The Kubernetes Authors.

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

package ec2

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAMIs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
	}{
		{
			name: "simple test",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("pretty new"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("pretty old"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock.EXPECT())

			clusterScope, err := setupCluster("test-cluster")
			g.Expect(err).To(Not(HaveOccurred()))

			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			id, err := s.defaultAMIIDLookup("", "", "base os-baseos version", "1.11.1")
			if err != nil {
				t.Fatalf("did not expect error calling a mock: %v", err)
			}
			if id != "pretty new" {
				t.Fatalf("returned %q expected 'pretty new'", id)
			}
		})
	}
}

func TestAMIsWithInvalidCreationDate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
	}{
		{
			name: "simple test",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("pretty new"),
								CreationDate: aws.String("invalid creation date"),
							},
							{
								ImageId:      aws.String("pretty old"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			g := NewWithT(t)

			clusterScope, err := setupCluster("test-cluster")
			g.Expect(err).To(Not(HaveOccurred()))

			tc.expect(ec2Mock.EXPECT())

			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			_, err = s.defaultAMIIDLookup("", "", "base os-baseos version", "1.11.1")
			if err == nil {
				t.Fatalf("expected an error but did not get one")
			}
		})
	}
}

func setupCluster(clusterName string) (*scope.ClusterScope, error) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec:       infrav1.AWSClusterSpec{},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).Build()
	return scope.NewClusterScope(scope.ClusterScopeParams{
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: clusterName},
		},
		AWSCluster: awsCluster,
		Client:     client,
	})
}

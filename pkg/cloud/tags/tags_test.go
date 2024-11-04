/*
Copyright 2020 The Kubernetes Authors.

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

package tags

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
)

var (
	expectedTags = []*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test"),
		},
		{
			Key:   aws.String("k1"),
			Value: aws.String("v1"),
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/testcluster"),
			Value: aws.String("owned"),
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
			Value: aws.String("testrole"),
		},
	}
)

func TestTagsComputeDiff(t *testing.T) {
	pName := "test"
	pRole := "testrole"
	bp := infrav1.BuildParams{
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		ClusterName: "testcluster",
		Name:        &pName,
		Role:        &pRole,
		Additional:  map[string]string{"k1": "v1"},
	}

	tests := []struct {
		name     string
		input    infrav1.Tags
		expected infrav1.Tags
	}{
		{
			name:  "input is nil",
			input: nil,
			expected: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v1",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         pRole,
			},
		},
		{
			name: "same input",
			input: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v1",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         pRole,
			},
			expected: infrav1.Tags{},
		},
		{
			name: "input with external tags",
			input: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v1",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         pRole,
				"k2":                                  "v2",
			},
			expected: infrav1.Tags{},
		},
		{
			name: "input with modified values",
			input: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v2",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         "testrole2",
				"k2":                                  "v2",
			},
			expected: infrav1.Tags{
				"k1":                          "v1",
				infrav1.NameAWSClusterAPIRole: pRole,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := computeDiff(tc.input, bp)
			if e, a := tc.expected, out; !cmp.Equal(e, a) {
				t.Errorf("expected %#v, got %#v", e, a)
			}
		})
	}
}

func TestTagsEnsureWithEC2(t *testing.T) {
	tests := []struct {
		name    string
		builder Builder
		expect  func(m *mocks.MockEC2APIMockRecorder)
	}{
		{
			name: "Should return error when create tag fails",
			builder: Builder{params: &infrav1.BuildParams{
				Lifecycle:   infrav1.ResourceLifecycleOwned,
				ClusterName: "testcluster",
				Name:        aws.String("test"),
				Role:        aws.String("testrole"),
				Additional:  map[string]string{"k1": "v1"},
			}},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.CreateTagsWithContext(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{""}),
					Tags:      expectedTags,
				})).Return(nil, errors.New("failed to create tag"))
			},
		},
		{
			name: "Should return error when optional configuration for builder is nil",
			builder: Builder{params: &infrav1.BuildParams{
				Lifecycle:   infrav1.ResourceLifecycleOwned,
				ClusterName: "testcluster",
				Name:        aws.String("test"),
				Role:        aws.String("testrole"),
				Additional:  map[string]string{"k1": "v1"},
			}, applyFunc: nil},
		},
		{
			name:    "Should return error when build params is nil",
			builder: Builder{params: nil},
		},
		{
			name: "Should ensure tags successfully",
			builder: Builder{params: &infrav1.BuildParams{
				Lifecycle:   infrav1.ResourceLifecycleOwned,
				ClusterName: "testcluster",
				Name:        aws.String("test"),
				Role:        aws.String("testrole"),
				Additional:  map[string]string{"k1": "v1"},
			}},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.CreateTagsWithContext(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{""}),
					Tags:      expectedTags,
				})).Return(nil, nil)
			},
		},
		{
			name: "Should filter internal aws tags",
			builder: Builder{params: &infrav1.BuildParams{
				Lifecycle:   infrav1.ResourceLifecycleOwned,
				ClusterName: "testcluster",
				Name:        aws.String("test"),
				Role:        aws.String("testrole"),
				Additional:  map[string]string{"k1": "v1", "aws:cloudformation:stack-name": "cloudformation-stack-name"},
			}},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.CreateTagsWithContext(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{""}),
					Tags:      expectedTags,
				})).Return(nil, nil)
			},
		},
	}

	g := NewWithT(t)
	mockCtrl := gomock.NewController(t)
	ec2Mock := mocks.NewMockEC2API(mockCtrl)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var builder *Builder
			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
				builder = New(tc.builder.params, WithEC2(ec2Mock))
			} else {
				builder = New(tc.builder.params, func(builder *Builder) {})
			}
			err := builder.Ensure(nil)
			if err != nil {
				g.Expect(err).To(Not(BeNil()))
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestTagsEnsureWithEKS(t *testing.T) {
	tests := []struct {
		name    string
		builder Builder
		expect  func(m *mock_eksiface.MockEKSAPIMockRecorder)
	}{
		{
			name: "Should return error when tag resources fails",
			builder: Builder{params: &infrav1.BuildParams{
				Lifecycle:   infrav1.ResourceLifecycleOwned,
				ClusterName: "testcluster",
				Name:        aws.String("test"),
				Role:        aws.String("testrole"),
				Additional:  map[string]string{"k1": "v1"},
			}},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.TagResource(gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String(""),
					Tags:        map[string]*string{"Name": aws.String("test"), "k1": aws.String("v1"), "sigs.k8s.io/cluster-api-provider-aws/cluster/testcluster": aws.String("owned"), "sigs.k8s.io/cluster-api-provider-aws/role": aws.String("testrole")},
				})).Return(nil, errors.New("failed to tag resource"))
			},
		},
		{
			name: "Should ensure tags successfully",
			builder: Builder{params: &infrav1.BuildParams{
				Lifecycle:   infrav1.ResourceLifecycleOwned,
				ClusterName: "testcluster",
				Name:        aws.String("test"),
				Role:        aws.String("testrole"),
				Additional:  map[string]string{"k1": "v1"},
			}},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.TagResource(gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String(""),
					Tags:        map[string]*string{"Name": aws.String("test"), "k1": aws.String("v1"), "sigs.k8s.io/cluster-api-provider-aws/cluster/testcluster": aws.String("owned"), "sigs.k8s.io/cluster-api-provider-aws/role": aws.String("testrole")},
				})).Return(nil, nil)
			},
		},
	}

	g := NewWithT(t)
	mockCtrl := gomock.NewController(t)
	eksMock := mock_eksiface.NewMockEKSAPI(mockCtrl)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var builder *Builder
			if tc.expect != nil {
				tc.expect(eksMock.EXPECT())
				builder = New(tc.builder.params, WithEKS(eksMock))
			} else {
				builder = New(tc.builder.params, func(builder *Builder) {})
			}
			err := builder.Ensure(nil)
			if err != nil {
				g.Expect(err).To(Not(BeNil()))
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestTagsBuildParamsToTagSpecification(t *testing.T) {
	g := NewWithT(t)
	tagSpec := BuildParamsToTagSpecification("test-resource", infrav1.BuildParams{
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		ClusterName: "testcluster",
		Name:        aws.String("test"),
		Role:        aws.String("testrole"),
		Additional:  map[string]string{"k1": "v1"},
	})
	expectedTagSpec := &ec2.TagSpecification{
		ResourceType: aws.String("test-resource"),
		Tags:         expectedTags,
	}
	g.Expect(expectedTagSpec).To(Equal(tagSpec))
}

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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
)

var (
	expectedTags = []ec2types.Tag{
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
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{""},
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
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{""},
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
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{""},
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
				m.TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String(""),
					Tags:        map[string]string{"Name": "test", "k1": "v1", "sigs.k8s.io/cluster-api-provider-aws/cluster/testcluster": "owned", "sigs.k8s.io/cluster-api-provider-aws/role": "testrole"},
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
				m.TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String(""),
					Tags:        map[string]string{"Name": "test", "k1": "v1", "sigs.k8s.io/cluster-api-provider-aws/cluster/testcluster": "owned", "sigs.k8s.io/cluster-api-provider-aws/role": "testrole"},
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
				builder = New(tc.builder.params, WithEKS(context.TODO(), eksMock))
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

func TestEnsureOnlyAppliesDiffTagsEC2(t *testing.T) {
	pName := "test"
	pRole := "testrole"
	params := &infrav1.BuildParams{
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		ClusterName: "testcluster",
		ResourceID:  "res-123",
		Name:        &pName,
		Role:        &pRole,
		Additional:  map[string]string{"k1": "v1"},
	}

	tests := []struct {
		name         string
		current      infrav1.Tags
		expectedTags []ec2types.Tag
		expect       func(m *mocks.MockEC2APIMockRecorder, expectedTags []ec2types.Tag)
	}{
		{
			name:    "no current tags, all tags are new",
			current: nil,
			expectedTags: []ec2types.Tag{
				{Key: aws.String("Name"), Value: aws.String(pName)},
				{Key: aws.String("k1"), Value: aws.String("v1")},
				{Key: aws.String(infrav1.ClusterTagKey("testcluster")), Value: aws.String(string(infrav1.ResourceLifecycleOwned))},
				{Key: aws.String(infrav1.NameAWSClusterAPIRole), Value: aws.String(pRole)},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder, expectedTags []ec2types.Tag) {
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"res-123"},
					Tags:      expectedTags,
				})).Return(nil, nil)
			},
		},
		{
			name: "all tags match, no apply call",
			current: infrav1.Tags{
				"Name":                               pName,
				"k1":                                 "v1",
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
			},
		},
		{
			name: "only changed tags are passed to apply",
			current: infrav1.Tags{
				"Name":                               pName,
				"k1":                                 "old-value",
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
			},
			expectedTags: []ec2types.Tag{
				{Key: aws.String("k1"), Value: aws.String("v1")},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder, expectedTags []ec2types.Tag) {
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"res-123"},
					Tags:      expectedTags,
				})).Return(nil, nil)
			},
		},
		{
			name: "only missing tags are passed to apply",
			current: infrav1.Tags{
				"Name":                               pName,
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
			},
			expectedTags: []ec2types.Tag{
				{Key: aws.String("k1"), Value: aws.String("v1")},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder, expectedTags []ec2types.Tag) {
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"res-123"},
					Tags:      expectedTags,
				})).Return(nil, nil)
			},
		},
		{
			name: "external tags do not cause apply",
			current: infrav1.Tags{
				"Name":                               pName,
				"k1":                                 "v1",
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
				"external-tag":                       "external-value",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			// The `diff` parameter should equal the test case's `expectedTags`.
			// This is to show better error messages than with mocking.
			var actualDiff infrav1.Tags
			builder := New(params, func(b *Builder) {
				b.applyFunc = func(_ *infrav1.BuildParams, diff infrav1.Tags) error {
					actualDiff = diff
					return nil
				}
			})
			err := builder.Ensure(tc.current)
			g.Expect(err).To(BeNil())
			var expectedTagsAsInfraType infrav1.Tags
			for _, tag := range tc.expectedTags {
				if expectedTagsAsInfraType == nil {
					expectedTagsAsInfraType = infrav1.Tags{}
				}
				expectedTagsAsInfraType[*tag.Key] = *tag.Value
			}
			g.Expect(actualDiff).To(Equal(expectedTagsAsInfraType), "Diff calculation is wrong, or the test case has a wrong expectedTags")

			// Now check with mocking since the actual AWS API calls are most important.
			// This ensures that the `diff` is used, and not all tags (incl. unchanged ones)
			// get sent in the AWS request.
			mockCtrl := gomock.NewController(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT(), tc.expectedTags)
			}

			builder = New(params, WithEC2(ec2Mock))

			err = builder.Ensure(tc.current)
			g.Expect(err).To(BeNil())
		})
	}
}

func TestEnsureOnlyAppliesDiffTagsEKS(t *testing.T) {
	pName := "test"
	pRole := "testrole"
	params := &infrav1.BuildParams{
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		ClusterName: "testcluster",
		ResourceID:  "res-123",
		Name:        &pName,
		Role:        &pRole,
		Additional:  map[string]string{"k1": "v1"},
	}

	tests := []struct {
		name         string
		current      infrav1.Tags
		expectedTags map[string]string
		expect       func(m *mock_eksiface.MockEKSAPIMockRecorder, expectedTags map[string]string)
	}{
		{
			name:    "no current tags, all tags are new",
			current: nil,
			expectedTags: map[string]string{
				"Name":                               pName,
				"k1":                                 "v1",
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder, expectedTags map[string]string) {
				m.TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String("res-123"),
					Tags:        expectedTags,
				})).Return(nil, nil)
			},
		},
		{
			name: "all tags match, no apply call",
			current: infrav1.Tags{
				"Name":                               pName,
				"k1":                                 "v1",
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
			},
		},
		{
			name: "only changed tags are passed to apply",
			current: infrav1.Tags{
				"Name":                               pName,
				"k1":                                 "old-value",
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
			},
			expectedTags: map[string]string{
				"k1": "v1",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder, expectedTags map[string]string) {
				m.TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String("res-123"),
					Tags:        expectedTags,
				})).Return(nil, nil)
			},
		},
		{
			name: "only missing tags are passed to apply",
			current: infrav1.Tags{
				"Name":                               pName,
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
			},
			expectedTags: map[string]string{
				"k1": "v1",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder, expectedTags map[string]string) {
				m.TagResource(gomock.Eq(context.TODO()), gomock.Eq(&eks.TagResourceInput{
					ResourceArn: aws.String("res-123"),
					Tags:        expectedTags,
				})).Return(nil, nil)
			},
		},
		{
			name: "external tags do not cause apply",
			current: infrav1.Tags{
				"Name":                               pName,
				"k1":                                 "v1",
				infrav1.ClusterTagKey("testcluster"): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:        pRole,
				"external-tag":                       "external-value",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			// The `diff` parameter should equal the test case's `expectedTags`.
			// This is to show better error messages than with mocking.
			var actualDiff infrav1.Tags
			builder := New(params, func(b *Builder) {
				b.applyFunc = func(_ *infrav1.BuildParams, diff infrav1.Tags) error {
					actualDiff = diff
					return nil
				}
			})
			err := builder.Ensure(tc.current)
			g.Expect(err).To(BeNil())
			var expectedTagsAsInfraType infrav1.Tags
			for k, v := range tc.expectedTags {
				if expectedTagsAsInfraType == nil {
					expectedTagsAsInfraType = infrav1.Tags{}
				}
				expectedTagsAsInfraType[k] = v
			}
			g.Expect(actualDiff).To(Equal(expectedTagsAsInfraType), "Diff calculation is wrong, or the test case has a wrong expectedTags")

			// Now check with mocking since the actual AWS API calls are most important.
			// This ensures that the `diff` is used, and not all tags (incl. unchanged ones)
			// get sent in the AWS request.
			mockCtrl := gomock.NewController(t)
			eksMock := mock_eksiface.NewMockEKSAPI(mockCtrl)

			if tc.expect != nil {
				tc.expect(eksMock.EXPECT(), tc.expectedTags)
			}

			builder = New(params, WithEKS(context.TODO(), eksMock))

			err = builder.Ensure(tc.current)
			g.Expect(err).To(BeNil())
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
	expectedTagSpec := ec2types.TagSpecification{
		ResourceType: "test-resource",
		Tags:         expectedTags,
	}
	g.Expect(expectedTagSpec).To(Equal(tagSpec))
}

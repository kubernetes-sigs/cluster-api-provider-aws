/*
Copyright The Kubernetes Authors.

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

package iamservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	go_cfn "github.com/awslabs/goformation/v4/cloudformation"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	"github.com/golang/mock/gomock"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/iamservice"
	mock_iamservice "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/iamservice/mock_iamservice"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

func testPolicyDocument() *iamv1.PolicyDocument {
	return &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: iamv1.Statements{
			{
				Effect: iamv1.EffectAllow,
				Principal: iamv1.Principals{
					iamv1.PrincipalService: iamv1.PrincipalID{"ec2.amazonaws.com"},
				},
				Action: iamv1.Actions{"sts:AssumeRole"},
			},
		},
	}
}

func rolesTemplate() go_cfn.Template {
	return go_cfn.Template{
		Resources: go_cfn.Resources{
			"AWSIAMRoleControllers": &cfn_iam.Role{
				RoleName:                 fmt.Sprintf("controllers%s", iamv1.DefaultNameSuffix),
				Description:              "For the Kubernetes Cluster API Provider AWS Controllers",
				AssumeRolePolicyDocument: testPolicyDocument(),
			},
		},
	}
}

func instanceProfilesTemplate() go_cfn.Template {
	return go_cfn.Template{
		Resources: go_cfn.Resources{
			"AWSIAMInstanceProfileControlPlane": &cfn_iam.InstanceProfile{
				InstanceProfileName: fmt.Sprintf("controllers%s", iamv1.DefaultNameSuffix),
			},
		},
	}
}

func policiesTemplate() go_cfn.Template {
	return go_cfn.Template{
		Resources: go_cfn.Resources{
			"ControllersPolicy": &cfn_iam.ManagedPolicy{
				ManagedPolicyName: fmt.Sprintf("controllers%s", iamv1.DefaultNameSuffix),
				Description:       "For the Kubernetes Cluster API Provider AWS Controllers",
				PolicyDocument:    testPolicyDocument(),
			},
		},
	}
}

func mergeTemplates(templates ...go_cfn.Template) go_cfn.Template {
	merged := go_cfn.Template{
		Resources: go_cfn.Resources{},
	}
	for _, t := range templates {
		for k, v := range t.Resources {
			merged.Resources[k] = v
		}
	}
	return merged
}

var testTags = map[string]string{"sigs.k8s.io/cluster-api-provider-aws/managed": "true"}

func TestCreateResources(t *testing.T) {
	// TODO: Add test cases to include conditions that invoke AddRoleToInstanceProfile() and ListPoliciesI)
	tests := []struct {
		name      string
		template  go_cfn.Template
		setupMock func(*mock_iamservice.MockIAMAPI)
		expectErr bool
	}{
		{
			name:     "create roles",
			template: rolesTemplate(),
			setupMock: func(m *mock_iamservice.MockIAMAPI) {
				m.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&iam.CreateRoleOutput{}, nil).Times(1)
			},
			expectErr: false,
		},
		{
			name:     "create instance profiles",
			template: instanceProfilesTemplate(),
			setupMock: func(m *mock_iamservice.MockIAMAPI) {
				m.EXPECT().CreateInstanceProfile(gomock.Any(), gomock.Any()).Return(&iam.CreateInstanceProfileOutput{}, nil).Times(1)
			},
			expectErr: false,
		},
		{
			name:     "create policies",
			template: policiesTemplate(),
			setupMock: func(m *mock_iamservice.MockIAMAPI) {
				m.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iam.CreatePolicyOutput{}, nil).Times(1)
			},
			expectErr: false,
		},
		{
			name:     "create roles, instance profiles and policies",
			template: mergeTemplates(rolesTemplate(), instanceProfilesTemplate(), policiesTemplate()),
			setupMock: func(m *mock_iamservice.MockIAMAPI) {
				m.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&iam.CreateRoleOutput{}, nil).Times(1)
				m.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any()).Return(&iam.AttachRolePolicyOutput{}, nil).AnyTimes()
				m.EXPECT().CreateInstanceProfile(gomock.Any(), gomock.Any()).Return(&iam.CreateInstanceProfileOutput{}, nil).Times(1)
				m.EXPECT().AddRoleToInstanceProfile(gomock.Any(), gomock.Any()).Return(&iam.AddRoleToInstanceProfileOutput{}, nil).AnyTimes()
				m.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iam.CreatePolicyOutput{}, nil).Times(1)
				m.EXPECT().ListPolicies(gomock.Any(), gomock.Any()).Return(&iam.ListPoliciesOutput{}, nil).AnyTimes()
			},
			expectErr: false,
		},
		{
			name:     "create roles and policies",
			template: mergeTemplates(rolesTemplate(), policiesTemplate()),
			setupMock: func(m *mock_iamservice.MockIAMAPI) {
				m.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&iam.CreateRoleOutput{}, nil).Times(1)
				m.EXPECT().AttachRolePolicy(gomock.Any(), gomock.Any()).Return(&iam.AttachRolePolicyOutput{}, nil).AnyTimes()
				m.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iam.CreatePolicyOutput{}, nil).Times(1)
				m.EXPECT().ListPolicies(gomock.Any(), gomock.Any()).Return(&iam.ListPoliciesOutput{}, nil).AnyTimes()
			},
			expectErr: false,
		},
		{
			name:     "create instance profiles and roles",
			template: mergeTemplates(instanceProfilesTemplate(), rolesTemplate()),
			setupMock: func(m *mock_iamservice.MockIAMAPI) {
				m.EXPECT().CreateInstanceProfile(gomock.Any(), gomock.Any()).Return(&iam.CreateInstanceProfileOutput{}, nil).Times(1)
				m.EXPECT().AddRoleToInstanceProfile(gomock.Any(), gomock.Any()).Return(&iam.AddRoleToInstanceProfileOutput{}, nil).AnyTimes()
				m.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(&iam.CreateRoleOutput{}, nil).Times(1)
			},
			expectErr: false,
		},
		{
			name:     "create policies and instance profiles",
			template: mergeTemplates(policiesTemplate(), instanceProfilesTemplate()),
			setupMock: func(m *mock_iamservice.MockIAMAPI) {
				m.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).Return(&iam.CreatePolicyOutput{}, nil).Times(1)
				m.EXPECT().CreateInstanceProfile(gomock.Any(), gomock.Any()).Return(&iam.CreateInstanceProfileOutput{}, nil).Times(1)
			},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_iamservice.NewMockIAMAPI(ctrl)
			tt.setupMock(m)
			svc := iamservice.New(m)
			err := svc.CreateResources(context.TODO(), tt.template, testTags)
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

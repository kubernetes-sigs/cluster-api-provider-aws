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

package mocks

import (
	"github.com/aws/aws-sdk-go/aws/session"
	gomock "github.com/golang/mock/gomock"
	v1alpha1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services"
)

// SDKGetter is a mocked getter.
type SDKGetter struct {
	EC2Mock *MockEC2Interface
	ELBMock *MockELBInterface
}

// Session returns a nil session.
func (*SDKGetter) Session(clusterConfig *v1alpha1.AWSClusterProviderConfig) *session.Session {
	return nil
}

// EC2 returns a mocked EC2 service client.
func (t *SDKGetter) EC2(session *session.Session) services.EC2Interface {
	return t.EC2Mock
}

// ELB returns a mocked EC2 service client.
func (t *SDKGetter) ELB(session *session.Session) services.ELBInterface {
	return t.ELBMock
}

// NewSDKGetter returns a new SDKGetter.
func NewSDKGetter(ctrl *gomock.Controller) *SDKGetter {
	return &SDKGetter{
		EC2Mock: NewMockEC2Interface(ctrl),
		ELBMock: NewMockELBInterface(ctrl),
	}
}

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

// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks (interfaces: EKSAPI)

// Package mock_eksiface is a generated GoMock package.
package mock_eksiface

import (
	context "context"
	reflect "reflect"
	time "time"

	eks "github.com/aws/aws-sdk-go-v2/service/eks"
	gomock "github.com/golang/mock/gomock"
)

// MockEKSAPI is a mock of EKSAPI interface.
type MockEKSAPI struct {
	ctrl     *gomock.Controller
	recorder *MockEKSAPIMockRecorder
}

// MockEKSAPIMockRecorder is the mock recorder for MockEKSAPI.
type MockEKSAPIMockRecorder struct {
	mock *MockEKSAPI
}

// NewMockEKSAPI creates a new mock instance.
func NewMockEKSAPI(ctrl *gomock.Controller) *MockEKSAPI {
	mock := &MockEKSAPI{ctrl: ctrl}
	mock.recorder = &MockEKSAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEKSAPI) EXPECT() *MockEKSAPIMockRecorder {
	return m.recorder
}

// AssociateEncryptionConfig mocks base method.
func (m *MockEKSAPI) AssociateEncryptionConfig(arg0 context.Context, arg1 *eks.AssociateEncryptionConfigInput, arg2 ...func(*eks.Options)) (*eks.AssociateEncryptionConfigOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AssociateEncryptionConfig", varargs...)
	ret0, _ := ret[0].(*eks.AssociateEncryptionConfigOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssociateEncryptionConfig indicates an expected call of AssociateEncryptionConfig.
func (mr *MockEKSAPIMockRecorder) AssociateEncryptionConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssociateEncryptionConfig", reflect.TypeOf((*MockEKSAPI)(nil).AssociateEncryptionConfig), varargs...)
}

// AssociateIdentityProviderConfig mocks base method.
func (m *MockEKSAPI) AssociateIdentityProviderConfig(arg0 context.Context, arg1 *eks.AssociateIdentityProviderConfigInput, arg2 ...func(*eks.Options)) (*eks.AssociateIdentityProviderConfigOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AssociateIdentityProviderConfig", varargs...)
	ret0, _ := ret[0].(*eks.AssociateIdentityProviderConfigOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssociateIdentityProviderConfig indicates an expected call of AssociateIdentityProviderConfig.
func (mr *MockEKSAPIMockRecorder) AssociateIdentityProviderConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssociateIdentityProviderConfig", reflect.TypeOf((*MockEKSAPI)(nil).AssociateIdentityProviderConfig), varargs...)
}

// CreateAddon mocks base method.
func (m *MockEKSAPI) CreateAddon(arg0 context.Context, arg1 *eks.CreateAddonInput, arg2 ...func(*eks.Options)) (*eks.CreateAddonOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAddon", varargs...)
	ret0, _ := ret[0].(*eks.CreateAddonOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAddon indicates an expected call of CreateAddon.
func (mr *MockEKSAPIMockRecorder) CreateAddon(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAddon", reflect.TypeOf((*MockEKSAPI)(nil).CreateAddon), varargs...)
}

// CreateCluster mocks base method.
func (m *MockEKSAPI) CreateCluster(arg0 context.Context, arg1 *eks.CreateClusterInput, arg2 ...func(*eks.Options)) (*eks.CreateClusterOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateCluster", varargs...)
	ret0, _ := ret[0].(*eks.CreateClusterOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCluster indicates an expected call of CreateCluster.
func (mr *MockEKSAPIMockRecorder) CreateCluster(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockEKSAPI)(nil).CreateCluster), varargs...)
}

// CreateFargateProfile mocks base method.
func (m *MockEKSAPI) CreateFargateProfile(arg0 context.Context, arg1 *eks.CreateFargateProfileInput, arg2 ...func(*eks.Options)) (*eks.CreateFargateProfileOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateFargateProfile", varargs...)
	ret0, _ := ret[0].(*eks.CreateFargateProfileOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFargateProfile indicates an expected call of CreateFargateProfile.
func (mr *MockEKSAPIMockRecorder) CreateFargateProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFargateProfile", reflect.TypeOf((*MockEKSAPI)(nil).CreateFargateProfile), varargs...)
}

// CreateNodegroup mocks base method.
func (m *MockEKSAPI) CreateNodegroup(arg0 context.Context, arg1 *eks.CreateNodegroupInput, arg2 ...func(*eks.Options)) (*eks.CreateNodegroupOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateNodegroup", varargs...)
	ret0, _ := ret[0].(*eks.CreateNodegroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNodegroup indicates an expected call of CreateNodegroup.
func (mr *MockEKSAPIMockRecorder) CreateNodegroup(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNodegroup", reflect.TypeOf((*MockEKSAPI)(nil).CreateNodegroup), varargs...)
}

// DeleteAddon mocks base method.
func (m *MockEKSAPI) DeleteAddon(arg0 context.Context, arg1 *eks.DeleteAddonInput, arg2 ...func(*eks.Options)) (*eks.DeleteAddonOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAddon", varargs...)
	ret0, _ := ret[0].(*eks.DeleteAddonOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAddon indicates an expected call of DeleteAddon.
func (mr *MockEKSAPIMockRecorder) DeleteAddon(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAddon", reflect.TypeOf((*MockEKSAPI)(nil).DeleteAddon), varargs...)
}

// DeleteCluster mocks base method.
func (m *MockEKSAPI) DeleteCluster(arg0 context.Context, arg1 *eks.DeleteClusterInput, arg2 ...func(*eks.Options)) (*eks.DeleteClusterOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteCluster", varargs...)
	ret0, _ := ret[0].(*eks.DeleteClusterOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCluster indicates an expected call of DeleteCluster.
func (mr *MockEKSAPIMockRecorder) DeleteCluster(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCluster", reflect.TypeOf((*MockEKSAPI)(nil).DeleteCluster), varargs...)
}

// DeleteFargateProfile mocks base method.
func (m *MockEKSAPI) DeleteFargateProfile(arg0 context.Context, arg1 *eks.DeleteFargateProfileInput, arg2 ...func(*eks.Options)) (*eks.DeleteFargateProfileOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteFargateProfile", varargs...)
	ret0, _ := ret[0].(*eks.DeleteFargateProfileOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFargateProfile indicates an expected call of DeleteFargateProfile.
func (mr *MockEKSAPIMockRecorder) DeleteFargateProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFargateProfile", reflect.TypeOf((*MockEKSAPI)(nil).DeleteFargateProfile), varargs...)
}

// DeleteNodegroup mocks base method.
func (m *MockEKSAPI) DeleteNodegroup(arg0 context.Context, arg1 *eks.DeleteNodegroupInput, arg2 ...func(*eks.Options)) (*eks.DeleteNodegroupOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteNodegroup", varargs...)
	ret0, _ := ret[0].(*eks.DeleteNodegroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNodegroup indicates an expected call of DeleteNodegroup.
func (mr *MockEKSAPIMockRecorder) DeleteNodegroup(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNodegroup", reflect.TypeOf((*MockEKSAPI)(nil).DeleteNodegroup), varargs...)
}

// DescribeAddon mocks base method.
func (m *MockEKSAPI) DescribeAddon(arg0 context.Context, arg1 *eks.DescribeAddonInput, arg2 ...func(*eks.Options)) (*eks.DescribeAddonOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeAddon", varargs...)
	ret0, _ := ret[0].(*eks.DescribeAddonOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeAddon indicates an expected call of DescribeAddon.
func (mr *MockEKSAPIMockRecorder) DescribeAddon(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeAddon", reflect.TypeOf((*MockEKSAPI)(nil).DescribeAddon), varargs...)
}

// DescribeAddonConfiguration mocks base method.
func (m *MockEKSAPI) DescribeAddonConfiguration(arg0 context.Context, arg1 *eks.DescribeAddonConfigurationInput, arg2 ...func(*eks.Options)) (*eks.DescribeAddonConfigurationOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeAddonConfiguration", varargs...)
	ret0, _ := ret[0].(*eks.DescribeAddonConfigurationOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeAddonConfiguration indicates an expected call of DescribeAddonConfiguration.
func (mr *MockEKSAPIMockRecorder) DescribeAddonConfiguration(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeAddonConfiguration", reflect.TypeOf((*MockEKSAPI)(nil).DescribeAddonConfiguration), varargs...)
}

// DescribeAddonVersions mocks base method.
func (m *MockEKSAPI) DescribeAddonVersions(arg0 context.Context, arg1 *eks.DescribeAddonVersionsInput, arg2 ...func(*eks.Options)) (*eks.DescribeAddonVersionsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeAddonVersions", varargs...)
	ret0, _ := ret[0].(*eks.DescribeAddonVersionsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeAddonVersions indicates an expected call of DescribeAddonVersions.
func (mr *MockEKSAPIMockRecorder) DescribeAddonVersions(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeAddonVersions", reflect.TypeOf((*MockEKSAPI)(nil).DescribeAddonVersions), varargs...)
}

// DescribeCluster mocks base method.
func (m *MockEKSAPI) DescribeCluster(arg0 context.Context, arg1 *eks.DescribeClusterInput, arg2 ...func(*eks.Options)) (*eks.DescribeClusterOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeCluster", varargs...)
	ret0, _ := ret[0].(*eks.DescribeClusterOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeCluster indicates an expected call of DescribeCluster.
func (mr *MockEKSAPIMockRecorder) DescribeCluster(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeCluster", reflect.TypeOf((*MockEKSAPI)(nil).DescribeCluster), varargs...)
}

// DescribeFargateProfile mocks base method.
func (m *MockEKSAPI) DescribeFargateProfile(arg0 context.Context, arg1 *eks.DescribeFargateProfileInput, arg2 ...func(*eks.Options)) (*eks.DescribeFargateProfileOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeFargateProfile", varargs...)
	ret0, _ := ret[0].(*eks.DescribeFargateProfileOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeFargateProfile indicates an expected call of DescribeFargateProfile.
func (mr *MockEKSAPIMockRecorder) DescribeFargateProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeFargateProfile", reflect.TypeOf((*MockEKSAPI)(nil).DescribeFargateProfile), varargs...)
}

// DescribeIdentityProviderConfig mocks base method.
func (m *MockEKSAPI) DescribeIdentityProviderConfig(arg0 context.Context, arg1 *eks.DescribeIdentityProviderConfigInput, arg2 ...func(*eks.Options)) (*eks.DescribeIdentityProviderConfigOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeIdentityProviderConfig", varargs...)
	ret0, _ := ret[0].(*eks.DescribeIdentityProviderConfigOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeIdentityProviderConfig indicates an expected call of DescribeIdentityProviderConfig.
func (mr *MockEKSAPIMockRecorder) DescribeIdentityProviderConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeIdentityProviderConfig", reflect.TypeOf((*MockEKSAPI)(nil).DescribeIdentityProviderConfig), varargs...)
}

// DescribeNodegroup mocks base method.
func (m *MockEKSAPI) DescribeNodegroup(arg0 context.Context, arg1 *eks.DescribeNodegroupInput, arg2 ...func(*eks.Options)) (*eks.DescribeNodegroupOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeNodegroup", varargs...)
	ret0, _ := ret[0].(*eks.DescribeNodegroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeNodegroup indicates an expected call of DescribeNodegroup.
func (mr *MockEKSAPIMockRecorder) DescribeNodegroup(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeNodegroup", reflect.TypeOf((*MockEKSAPI)(nil).DescribeNodegroup), varargs...)
}

// DescribeUpdate mocks base method.
func (m *MockEKSAPI) DescribeUpdate(arg0 context.Context, arg1 *eks.DescribeUpdateInput, arg2 ...func(*eks.Options)) (*eks.DescribeUpdateOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeUpdate", varargs...)
	ret0, _ := ret[0].(*eks.DescribeUpdateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeUpdate indicates an expected call of DescribeUpdate.
func (mr *MockEKSAPIMockRecorder) DescribeUpdate(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeUpdate", reflect.TypeOf((*MockEKSAPI)(nil).DescribeUpdate), varargs...)
}

// DisassociateIdentityProviderConfig mocks base method.
func (m *MockEKSAPI) DisassociateIdentityProviderConfig(arg0 context.Context, arg1 *eks.DisassociateIdentityProviderConfigInput, arg2 ...func(*eks.Options)) (*eks.DisassociateIdentityProviderConfigOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DisassociateIdentityProviderConfig", varargs...)
	ret0, _ := ret[0].(*eks.DisassociateIdentityProviderConfigOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisassociateIdentityProviderConfig indicates an expected call of DisassociateIdentityProviderConfig.
func (mr *MockEKSAPIMockRecorder) DisassociateIdentityProviderConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisassociateIdentityProviderConfig", reflect.TypeOf((*MockEKSAPI)(nil).DisassociateIdentityProviderConfig), varargs...)
}

// ListAddons mocks base method.
func (m *MockEKSAPI) ListAddons(arg0 context.Context, arg1 *eks.ListAddonsInput, arg2 ...func(*eks.Options)) (*eks.ListAddonsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListAddons", varargs...)
	ret0, _ := ret[0].(*eks.ListAddonsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAddons indicates an expected call of ListAddons.
func (mr *MockEKSAPIMockRecorder) ListAddons(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAddons", reflect.TypeOf((*MockEKSAPI)(nil).ListAddons), varargs...)
}

// ListClusters mocks base method.
func (m *MockEKSAPI) ListClusters(arg0 context.Context, arg1 *eks.ListClustersInput, arg2 ...func(*eks.Options)) (*eks.ListClustersOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListClusters", varargs...)
	ret0, _ := ret[0].(*eks.ListClustersOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClusters indicates an expected call of ListClusters.
func (mr *MockEKSAPIMockRecorder) ListClusters(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusters", reflect.TypeOf((*MockEKSAPI)(nil).ListClusters), varargs...)
}

// ListIdentityProviderConfigs mocks base method.
func (m *MockEKSAPI) ListIdentityProviderConfigs(arg0 context.Context, arg1 *eks.ListIdentityProviderConfigsInput, arg2 ...func(*eks.Options)) (*eks.ListIdentityProviderConfigsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListIdentityProviderConfigs", varargs...)
	ret0, _ := ret[0].(*eks.ListIdentityProviderConfigsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIdentityProviderConfigs indicates an expected call of ListIdentityProviderConfigs.
func (mr *MockEKSAPIMockRecorder) ListIdentityProviderConfigs(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIdentityProviderConfigs", reflect.TypeOf((*MockEKSAPI)(nil).ListIdentityProviderConfigs), varargs...)
}

// TagResource mocks base method.
func (m *MockEKSAPI) TagResource(arg0 context.Context, arg1 *eks.TagResourceInput, arg2 ...func(*eks.Options)) (*eks.TagResourceOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TagResource", varargs...)
	ret0, _ := ret[0].(*eks.TagResourceOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TagResource indicates an expected call of TagResource.
func (mr *MockEKSAPIMockRecorder) TagResource(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TagResource", reflect.TypeOf((*MockEKSAPI)(nil).TagResource), varargs...)
}

// UntagResource mocks base method.
func (m *MockEKSAPI) UntagResource(arg0 context.Context, arg1 *eks.UntagResourceInput, arg2 ...func(*eks.Options)) (*eks.UntagResourceOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UntagResource", varargs...)
	ret0, _ := ret[0].(*eks.UntagResourceOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UntagResource indicates an expected call of UntagResource.
func (mr *MockEKSAPIMockRecorder) UntagResource(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UntagResource", reflect.TypeOf((*MockEKSAPI)(nil).UntagResource), varargs...)
}

// UpdateAddon mocks base method.
func (m *MockEKSAPI) UpdateAddon(arg0 context.Context, arg1 *eks.UpdateAddonInput, arg2 ...func(*eks.Options)) (*eks.UpdateAddonOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateAddon", varargs...)
	ret0, _ := ret[0].(*eks.UpdateAddonOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAddon indicates an expected call of UpdateAddon.
func (mr *MockEKSAPIMockRecorder) UpdateAddon(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAddon", reflect.TypeOf((*MockEKSAPI)(nil).UpdateAddon), varargs...)
}

// UpdateClusterConfig mocks base method.
func (m *MockEKSAPI) UpdateClusterConfig(arg0 context.Context, arg1 *eks.UpdateClusterConfigInput, arg2 ...func(*eks.Options)) (*eks.UpdateClusterConfigOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateClusterConfig", varargs...)
	ret0, _ := ret[0].(*eks.UpdateClusterConfigOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateClusterConfig indicates an expected call of UpdateClusterConfig.
func (mr *MockEKSAPIMockRecorder) UpdateClusterConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterConfig", reflect.TypeOf((*MockEKSAPI)(nil).UpdateClusterConfig), varargs...)
}

// UpdateClusterVersion mocks base method.
func (m *MockEKSAPI) UpdateClusterVersion(arg0 context.Context, arg1 *eks.UpdateClusterVersionInput, arg2 ...func(*eks.Options)) (*eks.UpdateClusterVersionOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateClusterVersion", varargs...)
	ret0, _ := ret[0].(*eks.UpdateClusterVersionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateClusterVersion indicates an expected call of UpdateClusterVersion.
func (mr *MockEKSAPIMockRecorder) UpdateClusterVersion(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterVersion", reflect.TypeOf((*MockEKSAPI)(nil).UpdateClusterVersion), varargs...)
}

// UpdateNodegroupConfig mocks base method.
func (m *MockEKSAPI) UpdateNodegroupConfig(arg0 context.Context, arg1 *eks.UpdateNodegroupConfigInput, arg2 ...func(*eks.Options)) (*eks.UpdateNodegroupConfigOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateNodegroupConfig", varargs...)
	ret0, _ := ret[0].(*eks.UpdateNodegroupConfigOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNodegroupConfig indicates an expected call of UpdateNodegroupConfig.
func (mr *MockEKSAPIMockRecorder) UpdateNodegroupConfig(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodegroupConfig", reflect.TypeOf((*MockEKSAPI)(nil).UpdateNodegroupConfig), varargs...)
}

// UpdateNodegroupVersion mocks base method.
func (m *MockEKSAPI) UpdateNodegroupVersion(arg0 context.Context, arg1 *eks.UpdateNodegroupVersionInput, arg2 ...func(*eks.Options)) (*eks.UpdateNodegroupVersionOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateNodegroupVersion", varargs...)
	ret0, _ := ret[0].(*eks.UpdateNodegroupVersionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNodegroupVersion indicates an expected call of UpdateNodegroupVersion.
func (mr *MockEKSAPIMockRecorder) UpdateNodegroupVersion(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodegroupVersion", reflect.TypeOf((*MockEKSAPI)(nil).UpdateNodegroupVersion), varargs...)
}

// WaitUntilAddonDeleted mocks base method.
func (m *MockEKSAPI) WaitUntilAddonDeleted(arg0 context.Context, arg1 *eks.DescribeAddonInput, arg2 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilAddonDeleted", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilAddonDeleted indicates an expected call of WaitUntilAddonDeleted.
func (mr *MockEKSAPIMockRecorder) WaitUntilAddonDeleted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilAddonDeleted", reflect.TypeOf((*MockEKSAPI)(nil).WaitUntilAddonDeleted), arg0, arg1, arg2)
}

// WaitUntilClusterActive mocks base method.
func (m *MockEKSAPI) WaitUntilClusterActive(arg0 context.Context, arg1 *eks.DescribeClusterInput, arg2 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilClusterActive", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilClusterActive indicates an expected call of WaitUntilClusterActive.
func (mr *MockEKSAPIMockRecorder) WaitUntilClusterActive(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilClusterActive", reflect.TypeOf((*MockEKSAPI)(nil).WaitUntilClusterActive), arg0, arg1, arg2)
}

// WaitUntilClusterDeleted mocks base method.
func (m *MockEKSAPI) WaitUntilClusterDeleted(arg0 context.Context, arg1 *eks.DescribeClusterInput, arg2 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilClusterDeleted", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilClusterDeleted indicates an expected call of WaitUntilClusterDeleted.
func (mr *MockEKSAPIMockRecorder) WaitUntilClusterDeleted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilClusterDeleted", reflect.TypeOf((*MockEKSAPI)(nil).WaitUntilClusterDeleted), arg0, arg1, arg2)
}

// WaitUntilClusterUpdating mocks base method.
func (m *MockEKSAPI) WaitUntilClusterUpdating(arg0 context.Context, arg1 *eks.DescribeClusterInput, arg2 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilClusterUpdating", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilClusterUpdating indicates an expected call of WaitUntilClusterUpdating.
func (mr *MockEKSAPIMockRecorder) WaitUntilClusterUpdating(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilClusterUpdating", reflect.TypeOf((*MockEKSAPI)(nil).WaitUntilClusterUpdating), arg0, arg1, arg2)
}

// WaitUntilNodegroupActive mocks base method.
func (m *MockEKSAPI) WaitUntilNodegroupActive(arg0 context.Context, arg1 *eks.DescribeNodegroupInput, arg2 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilNodegroupActive", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilNodegroupActive indicates an expected call of WaitUntilNodegroupActive.
func (mr *MockEKSAPIMockRecorder) WaitUntilNodegroupActive(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilNodegroupActive", reflect.TypeOf((*MockEKSAPI)(nil).WaitUntilNodegroupActive), arg0, arg1, arg2)
}

// WaitUntilNodegroupDeleted mocks base method.
func (m *MockEKSAPI) WaitUntilNodegroupDeleted(arg0 context.Context, arg1 *eks.DescribeNodegroupInput, arg2 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilNodegroupDeleted", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilNodegroupDeleted indicates an expected call of WaitUntilNodegroupDeleted.
func (mr *MockEKSAPIMockRecorder) WaitUntilNodegroupDeleted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilNodegroupDeleted", reflect.TypeOf((*MockEKSAPI)(nil).WaitUntilNodegroupDeleted), arg0, arg1, arg2)
}

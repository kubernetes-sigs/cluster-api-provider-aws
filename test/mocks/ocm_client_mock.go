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
// Source: sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa (interfaces: OCMClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	aws "github.com/openshift/rosa/pkg/aws"
	ocm "github.com/openshift/rosa/pkg/ocm"
)

// MockOCMClient is a mock of OCMClient interface.
type MockOCMClient struct {
	ctrl     *gomock.Controller
	recorder *MockOCMClientMockRecorder
}

// MockOCMClientMockRecorder is the mock recorder for MockOCMClient.
type MockOCMClientMockRecorder struct {
	mock *MockOCMClient
}

// NewMockOCMClient creates a new mock instance.
func NewMockOCMClient(ctrl *gomock.Controller) *MockOCMClient {
	mock := &MockOCMClient{ctrl: ctrl}
	mock.recorder = &MockOCMClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOCMClient) EXPECT() *MockOCMClientMockRecorder {
	return m.recorder
}

// AckVersionGate mocks base method.
func (m *MockOCMClient) AckVersionGate(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AckVersionGate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AckVersionGate indicates an expected call of AckVersionGate.
func (mr *MockOCMClientMockRecorder) AckVersionGate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AckVersionGate", reflect.TypeOf((*MockOCMClient)(nil).AckVersionGate), arg0, arg1)
}

// AddHTPasswdUser mocks base method.
func (m *MockOCMClient) AddHTPasswdUser(arg0, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddHTPasswdUser", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddHTPasswdUser indicates an expected call of AddHTPasswdUser.
func (mr *MockOCMClientMockRecorder) AddHTPasswdUser(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHTPasswdUser", reflect.TypeOf((*MockOCMClient)(nil).AddHTPasswdUser), arg0, arg1, arg2, arg3)
}

// CreateCluster mocks base method.
func (m *MockOCMClient) CreateCluster(arg0 ocm.Spec) (*v1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCluster", arg0)
	ret0, _ := ret[0].(*v1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCluster indicates an expected call of CreateCluster.
func (mr *MockOCMClientMockRecorder) CreateCluster(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockOCMClient)(nil).CreateCluster), arg0)
}

// CreateIdentityProvider mocks base method.
func (m *MockOCMClient) CreateIdentityProvider(arg0 string, arg1 *v1.IdentityProvider) (*v1.IdentityProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIdentityProvider", arg0, arg1)
	ret0, _ := ret[0].(*v1.IdentityProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIdentityProvider indicates an expected call of CreateIdentityProvider.
func (mr *MockOCMClientMockRecorder) CreateIdentityProvider(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIdentityProvider", reflect.TypeOf((*MockOCMClient)(nil).CreateIdentityProvider), arg0, arg1)
}

// CreateNodePool mocks base method.
func (m *MockOCMClient) CreateNodePool(arg0 string, arg1 *v1.NodePool) (*v1.NodePool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNodePool", arg0, arg1)
	ret0, _ := ret[0].(*v1.NodePool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNodePool indicates an expected call of CreateNodePool.
func (mr *MockOCMClientMockRecorder) CreateNodePool(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNodePool", reflect.TypeOf((*MockOCMClient)(nil).CreateNodePool), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockOCMClient) CreateUser(arg0, arg1 string, arg2 *v1.User) (*v1.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockOCMClientMockRecorder) CreateUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockOCMClient)(nil).CreateUser), arg0, arg1, arg2)
}

// DeleteCluster mocks base method.
func (m *MockOCMClient) DeleteCluster(arg0 string, arg1 bool, arg2 *aws.Creator) (*v1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCluster", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCluster indicates an expected call of DeleteCluster.
func (mr *MockOCMClientMockRecorder) DeleteCluster(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCluster", reflect.TypeOf((*MockOCMClient)(nil).DeleteCluster), arg0, arg1, arg2)
}

// DeleteNodePool mocks base method.
func (m *MockOCMClient) DeleteNodePool(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNodePool", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNodePool indicates an expected call of DeleteNodePool.
func (mr *MockOCMClientMockRecorder) DeleteNodePool(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNodePool", reflect.TypeOf((*MockOCMClient)(nil).DeleteNodePool), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockOCMClient) DeleteUser(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockOCMClientMockRecorder) DeleteUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockOCMClient)(nil).DeleteUser), arg0, arg1, arg2)
}

// GetCluster mocks base method.
func (m *MockOCMClient) GetCluster(arg0 string, arg1 *aws.Creator) (*v1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", arg0, arg1)
	ret0, _ := ret[0].(*v1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCluster indicates an expected call of GetCluster.
func (mr *MockOCMClientMockRecorder) GetCluster(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockOCMClient)(nil).GetCluster), arg0, arg1)
}

// GetControlPlaneUpgradePolicies mocks base method.
func (m *MockOCMClient) GetControlPlaneUpgradePolicies(arg0 string) ([]*v1.ControlPlaneUpgradePolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetControlPlaneUpgradePolicies", arg0)
	ret0, _ := ret[0].([]*v1.ControlPlaneUpgradePolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetControlPlaneUpgradePolicies indicates an expected call of GetControlPlaneUpgradePolicies.
func (mr *MockOCMClientMockRecorder) GetControlPlaneUpgradePolicies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetControlPlaneUpgradePolicies", reflect.TypeOf((*MockOCMClient)(nil).GetControlPlaneUpgradePolicies), arg0)
}

// GetHTPasswdUserList mocks base method.
func (m *MockOCMClient) GetHTPasswdUserList(arg0, arg1 string) (*v1.HTPasswdUserList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHTPasswdUserList", arg0, arg1)
	ret0, _ := ret[0].(*v1.HTPasswdUserList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHTPasswdUserList indicates an expected call of GetHTPasswdUserList.
func (mr *MockOCMClientMockRecorder) GetHTPasswdUserList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHTPasswdUserList", reflect.TypeOf((*MockOCMClient)(nil).GetHTPasswdUserList), arg0, arg1)
}

// GetHypershiftNodePoolUpgrade mocks base method.
func (m *MockOCMClient) GetHypershiftNodePoolUpgrade(arg0, arg1, arg2 string) (*v1.NodePool, *v1.NodePoolUpgradePolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHypershiftNodePoolUpgrade", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.NodePool)
	ret1, _ := ret[1].(*v1.NodePoolUpgradePolicy)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetHypershiftNodePoolUpgrade indicates an expected call of GetHypershiftNodePoolUpgrade.
func (mr *MockOCMClientMockRecorder) GetHypershiftNodePoolUpgrade(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHypershiftNodePoolUpgrade", reflect.TypeOf((*MockOCMClient)(nil).GetHypershiftNodePoolUpgrade), arg0, arg1, arg2)
}

// GetIdentityProviders mocks base method.
func (m *MockOCMClient) GetIdentityProviders(arg0 string) ([]*v1.IdentityProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIdentityProviders", arg0)
	ret0, _ := ret[0].([]*v1.IdentityProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIdentityProviders indicates an expected call of GetIdentityProviders.
func (mr *MockOCMClientMockRecorder) GetIdentityProviders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIdentityProviders", reflect.TypeOf((*MockOCMClient)(nil).GetIdentityProviders), arg0)
}

// GetMissingGateAgreementsHypershift mocks base method.
func (m *MockOCMClient) GetMissingGateAgreementsHypershift(arg0 string, arg1 *v1.ControlPlaneUpgradePolicy) ([]*v1.VersionGate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMissingGateAgreementsHypershift", arg0, arg1)
	ret0, _ := ret[0].([]*v1.VersionGate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMissingGateAgreementsHypershift indicates an expected call of GetMissingGateAgreementsHypershift.
func (mr *MockOCMClientMockRecorder) GetMissingGateAgreementsHypershift(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMissingGateAgreementsHypershift", reflect.TypeOf((*MockOCMClient)(nil).GetMissingGateAgreementsHypershift), arg0, arg1)
}

// GetNodePool mocks base method.
func (m *MockOCMClient) GetNodePool(arg0, arg1 string) (*v1.NodePool, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodePool", arg0, arg1)
	ret0, _ := ret[0].(*v1.NodePool)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNodePool indicates an expected call of GetNodePool.
func (mr *MockOCMClientMockRecorder) GetNodePool(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodePool", reflect.TypeOf((*MockOCMClient)(nil).GetNodePool), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockOCMClient) GetUser(arg0, arg1, arg2 string) (*v1.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockOCMClientMockRecorder) GetUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockOCMClient)(nil).GetUser), arg0, arg1, arg2)
}

// ScheduleHypershiftControlPlaneUpgrade mocks base method.
func (m *MockOCMClient) ScheduleHypershiftControlPlaneUpgrade(arg0 string, arg1 *v1.ControlPlaneUpgradePolicy) (*v1.ControlPlaneUpgradePolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleHypershiftControlPlaneUpgrade", arg0, arg1)
	ret0, _ := ret[0].(*v1.ControlPlaneUpgradePolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScheduleHypershiftControlPlaneUpgrade indicates an expected call of ScheduleHypershiftControlPlaneUpgrade.
func (mr *MockOCMClientMockRecorder) ScheduleHypershiftControlPlaneUpgrade(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleHypershiftControlPlaneUpgrade", reflect.TypeOf((*MockOCMClient)(nil).ScheduleHypershiftControlPlaneUpgrade), arg0, arg1)
}

// ScheduleNodePoolUpgrade mocks base method.
func (m *MockOCMClient) ScheduleNodePoolUpgrade(arg0, arg1 string, arg2 *v1.NodePoolUpgradePolicy) (*v1.NodePoolUpgradePolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleNodePoolUpgrade", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.NodePoolUpgradePolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScheduleNodePoolUpgrade indicates an expected call of ScheduleNodePoolUpgrade.
func (mr *MockOCMClientMockRecorder) ScheduleNodePoolUpgrade(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleNodePoolUpgrade", reflect.TypeOf((*MockOCMClient)(nil).ScheduleNodePoolUpgrade), arg0, arg1, arg2)
}

// UpdateCluster mocks base method.
func (m *MockOCMClient) UpdateCluster(arg0 string, arg1 *aws.Creator, arg2 ocm.Spec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCluster", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCluster indicates an expected call of UpdateCluster.
func (mr *MockOCMClientMockRecorder) UpdateCluster(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCluster", reflect.TypeOf((*MockOCMClient)(nil).UpdateCluster), arg0, arg1, arg2)
}

// UpdateNodePool mocks base method.
func (m *MockOCMClient) UpdateNodePool(arg0 string, arg1 *v1.NodePool) (*v1.NodePool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNodePool", arg0, arg1)
	ret0, _ := ret[0].(*v1.NodePool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNodePool indicates an expected call of UpdateNodePool.
func (mr *MockOCMClientMockRecorder) UpdateNodePool(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodePool", reflect.TypeOf((*MockOCMClient)(nil).UpdateNodePool), arg0, arg1)
}

// ValidateHypershiftVersion mocks base method.
func (m *MockOCMClient) ValidateHypershiftVersion(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateHypershiftVersion", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateHypershiftVersion indicates an expected call of ValidateHypershiftVersion.
func (mr *MockOCMClientMockRecorder) ValidateHypershiftVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateHypershiftVersion", reflect.TypeOf((*MockOCMClient)(nil).ValidateHypershiftVersion), arg0, arg1)
}

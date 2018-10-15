// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services (interfaces: EC2Interface)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	v1alpha1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	v1alpha10 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// MockEC2Interface is a mock of EC2Interface interface
type MockEC2Interface struct {
	ctrl     *gomock.Controller
	recorder *MockEC2InterfaceMockRecorder
}

// MockEC2InterfaceMockRecorder is the mock recorder for MockEC2Interface
type MockEC2InterfaceMockRecorder struct {
	mock *MockEC2Interface
}

// NewMockEC2Interface creates a new mock instance
func NewMockEC2Interface(ctrl *gomock.Controller) *MockEC2Interface {
	mock := &MockEC2Interface{ctrl: ctrl}
	mock.recorder = &MockEC2InterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEC2Interface) EXPECT() *MockEC2InterfaceMockRecorder {
	return m.recorder
}

// CreateInstance mocks base method
func (m *MockEC2Interface) CreateInstance(arg0 *v1alpha10.Machine, arg1 *v1alpha1.AWSMachineProviderConfig, arg2 *v1alpha10.Cluster, arg3 *v1alpha1.AWSClusterProviderConfig) (*v1alpha1.Instance, error) {
	ret := m.ctrl.Call(m, "CreateInstance", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*v1alpha1.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateInstance indicates an expected call of CreateInstance
func (mr *MockEC2InterfaceMockRecorder) CreateInstance(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInstance", reflect.TypeOf((*MockEC2Interface)(nil).CreateInstance), arg0, arg1, arg2, arg3)
}

// CreateOrGetMachine mocks base method
func (m *MockEC2Interface) CreateOrGetMachine(arg0 *v1alpha10.Machine, arg1 *v1alpha1.AWSMachineProviderConfig, arg2 *v1alpha10.Cluster, arg3 *v1alpha1.AWSClusterProviderConfig) (*v1alpha1.Instance, error) {
	ret := m.ctrl.Call(m, "CreateOrGetMachine", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*v1alpha1.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrGetMachine indicates an expected call of CreateOrGetMachine
func (mr *MockEC2InterfaceMockRecorder) CreateOrGetMachine(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrGetMachine", reflect.TypeOf((*MockEC2Interface)(nil).CreateOrGetMachine), arg0, arg1, arg2, arg3)
}

// DeleteBastion mocks base method
func (m *MockEC2Interface) DeleteBastion(arg0 string, arg1 *v1alpha1.AWSClusterProviderConfigStatus) error {
	ret := m.ctrl.Call(m, "DeleteBastion", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBastion indicates an expected call of DeleteBastion
func (mr *MockEC2InterfaceMockRecorder) DeleteBastion(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBastion", reflect.TypeOf((*MockEC2Interface)(nil).DeleteBastion), arg0, arg1)
}

// DeleteNetwork mocks base method
func (m *MockEC2Interface) DeleteNetwork(arg0 string, arg1 *v1alpha1.Network) error {
	ret := m.ctrl.Call(m, "DeleteNetwork", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNetwork indicates an expected call of DeleteNetwork
func (mr *MockEC2InterfaceMockRecorder) DeleteNetwork(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNetwork", reflect.TypeOf((*MockEC2Interface)(nil).DeleteNetwork), arg0, arg1)
}

// InstanceIfExists mocks base method
func (m *MockEC2Interface) InstanceIfExists(arg0 *string) (*v1alpha1.Instance, error) {
	ret := m.ctrl.Call(m, "InstanceIfExists", arg0)
	ret0, _ := ret[0].(*v1alpha1.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceIfExists indicates an expected call of InstanceIfExists
func (mr *MockEC2InterfaceMockRecorder) InstanceIfExists(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceIfExists", reflect.TypeOf((*MockEC2Interface)(nil).InstanceIfExists), arg0)
}

// ReconcileBastion mocks base method
func (m *MockEC2Interface) ReconcileBastion(arg0, arg1 string, arg2 *v1alpha1.AWSClusterProviderConfigStatus) error {
	ret := m.ctrl.Call(m, "ReconcileBastion", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileBastion indicates an expected call of ReconcileBastion
func (mr *MockEC2InterfaceMockRecorder) ReconcileBastion(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileBastion", reflect.TypeOf((*MockEC2Interface)(nil).ReconcileBastion), arg0, arg1, arg2)
}

// ReconcileNetwork mocks base method
func (m *MockEC2Interface) ReconcileNetwork(arg0 string, arg1 *v1alpha1.Network) error {
	ret := m.ctrl.Call(m, "ReconcileNetwork", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileNetwork indicates an expected call of ReconcileNetwork
func (mr *MockEC2InterfaceMockRecorder) ReconcileNetwork(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileNetwork", reflect.TypeOf((*MockEC2Interface)(nil).ReconcileNetwork), arg0, arg1)
}

// TerminateInstance mocks base method
func (m *MockEC2Interface) TerminateInstance(arg0 string) error {
	ret := m.ctrl.Call(m, "TerminateInstance", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TerminateInstance indicates an expected call of TerminateInstance
func (mr *MockEC2InterfaceMockRecorder) TerminateInstance(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TerminateInstance", reflect.TypeOf((*MockEC2Interface)(nil).TerminateInstance), arg0)
}

// UpdateInstanceSecurityGroups mocks base method
func (m *MockEC2Interface) UpdateInstanceSecurityGroups(arg0 string, arg1 []string) error {
	ret := m.ctrl.Call(m, "UpdateInstanceSecurityGroups", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInstanceSecurityGroups indicates an expected call of UpdateInstanceSecurityGroups
func (mr *MockEC2InterfaceMockRecorder) UpdateInstanceSecurityGroups(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInstanceSecurityGroups", reflect.TypeOf((*MockEC2Interface)(nil).UpdateInstanceSecurityGroups), arg0, arg1)
}

// UpdateResourceTags mocks base method
func (m *MockEC2Interface) UpdateResourceTags(arg0 *string, arg1, arg2 map[string]string) error {
	ret := m.ctrl.Call(m, "UpdateResourceTags", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateResourceTags indicates an expected call of UpdateResourceTags
func (mr *MockEC2InterfaceMockRecorder) UpdateResourceTags(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateResourceTags", reflect.TypeOf((*MockEC2Interface)(nil).UpdateResourceTags), arg0, arg1, arg2)
}

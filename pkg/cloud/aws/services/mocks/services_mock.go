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

// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services (interfaces: EC2Interface,ELBInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	v1alpha1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	actuators "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
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

// CreateOrGetMachine mocks base method
func (m *MockEC2Interface) CreateOrGetMachine(arg0 *actuators.MachineScope, arg1 string) (*v1alpha1.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrGetMachine", arg0, arg1)
	ret0, _ := ret[0].(*v1alpha1.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrGetMachine indicates an expected call of CreateOrGetMachine
func (mr *MockEC2InterfaceMockRecorder) CreateOrGetMachine(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrGetMachine", reflect.TypeOf((*MockEC2Interface)(nil).CreateOrGetMachine), arg0, arg1)
}

// DeleteBastion mocks base method
func (m *MockEC2Interface) DeleteBastion() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBastion")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBastion indicates an expected call of DeleteBastion
func (mr *MockEC2InterfaceMockRecorder) DeleteBastion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBastion", reflect.TypeOf((*MockEC2Interface)(nil).DeleteBastion))
}

// DeleteNetwork mocks base method
func (m *MockEC2Interface) DeleteNetwork() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNetwork")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNetwork indicates an expected call of DeleteNetwork
func (mr *MockEC2InterfaceMockRecorder) DeleteNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNetwork", reflect.TypeOf((*MockEC2Interface)(nil).DeleteNetwork))
}

// GetCoreSecurityGroups mocks base method
func (m *MockEC2Interface) GetCoreSecurityGroups(arg0 *actuators.MachineScope) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoreSecurityGroups", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoreSecurityGroups indicates an expected call of GetCoreSecurityGroups
func (mr *MockEC2InterfaceMockRecorder) GetCoreSecurityGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoreSecurityGroups", reflect.TypeOf((*MockEC2Interface)(nil).GetCoreSecurityGroups), arg0)
}

// GetInstanceSecurityGroups mocks base method
func (m *MockEC2Interface) GetInstanceSecurityGroups(arg0 string) (map[string][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstanceSecurityGroups", arg0)
	ret0, _ := ret[0].(map[string][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstanceSecurityGroups indicates an expected call of GetInstanceSecurityGroups
func (mr *MockEC2InterfaceMockRecorder) GetInstanceSecurityGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstanceSecurityGroups", reflect.TypeOf((*MockEC2Interface)(nil).GetInstanceSecurityGroups), arg0)
}

// InstanceIfExists mocks base method
func (m *MockEC2Interface) InstanceIfExists(arg0 *string) (*v1alpha1.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceIfExists", arg0)
	ret0, _ := ret[0].(*v1alpha1.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceIfExists indicates an expected call of InstanceIfExists
func (mr *MockEC2InterfaceMockRecorder) InstanceIfExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceIfExists", reflect.TypeOf((*MockEC2Interface)(nil).InstanceIfExists), arg0)
}

// ReconcileBastion mocks base method
func (m *MockEC2Interface) ReconcileBastion() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileBastion")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileBastion indicates an expected call of ReconcileBastion
func (mr *MockEC2InterfaceMockRecorder) ReconcileBastion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileBastion", reflect.TypeOf((*MockEC2Interface)(nil).ReconcileBastion))
}

// ReconcileNetwork mocks base method
func (m *MockEC2Interface) ReconcileNetwork() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileNetwork")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileNetwork indicates an expected call of ReconcileNetwork
func (mr *MockEC2InterfaceMockRecorder) ReconcileNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileNetwork", reflect.TypeOf((*MockEC2Interface)(nil).ReconcileNetwork))
}

// TerminateInstance mocks base method
func (m *MockEC2Interface) TerminateInstance(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TerminateInstance", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TerminateInstance indicates an expected call of TerminateInstance
func (mr *MockEC2InterfaceMockRecorder) TerminateInstance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TerminateInstance", reflect.TypeOf((*MockEC2Interface)(nil).TerminateInstance), arg0)
}

// UpdateInstanceSecurityGroups mocks base method
func (m *MockEC2Interface) UpdateInstanceSecurityGroups(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInstanceSecurityGroups", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInstanceSecurityGroups indicates an expected call of UpdateInstanceSecurityGroups
func (mr *MockEC2InterfaceMockRecorder) UpdateInstanceSecurityGroups(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInstanceSecurityGroups", reflect.TypeOf((*MockEC2Interface)(nil).UpdateInstanceSecurityGroups), arg0, arg1)
}

// UpdateResourceTags mocks base method
func (m *MockEC2Interface) UpdateResourceTags(arg0 *string, arg1, arg2 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateResourceTags", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateResourceTags indicates an expected call of UpdateResourceTags
func (mr *MockEC2InterfaceMockRecorder) UpdateResourceTags(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateResourceTags", reflect.TypeOf((*MockEC2Interface)(nil).UpdateResourceTags), arg0, arg1, arg2)
}

// MockELBInterface is a mock of ELBInterface interface
type MockELBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockELBInterfaceMockRecorder
}

// MockELBInterfaceMockRecorder is the mock recorder for MockELBInterface
type MockELBInterfaceMockRecorder struct {
	mock *MockELBInterface
}

// NewMockELBInterface creates a new mock instance
func NewMockELBInterface(ctrl *gomock.Controller) *MockELBInterface {
	mock := &MockELBInterface{ctrl: ctrl}
	mock.recorder = &MockELBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockELBInterface) EXPECT() *MockELBInterfaceMockRecorder {
	return m.recorder
}

// DeleteLoadbalancers mocks base method
func (m *MockELBInterface) DeleteLoadbalancers() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLoadbalancers")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLoadbalancers indicates an expected call of DeleteLoadbalancers
func (mr *MockELBInterfaceMockRecorder) DeleteLoadbalancers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLoadbalancers", reflect.TypeOf((*MockELBInterface)(nil).DeleteLoadbalancers))
}

// GetAPIServerDNSName mocks base method
func (m *MockELBInterface) GetAPIServerDNSName() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPIServerDNSName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAPIServerDNSName indicates an expected call of GetAPIServerDNSName
func (mr *MockELBInterfaceMockRecorder) GetAPIServerDNSName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIServerDNSName", reflect.TypeOf((*MockELBInterface)(nil).GetAPIServerDNSName))
}

// ReconcileLoadbalancers mocks base method
func (m *MockELBInterface) ReconcileLoadbalancers() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileLoadbalancers")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileLoadbalancers indicates an expected call of ReconcileLoadbalancers
func (mr *MockELBInterfaceMockRecorder) ReconcileLoadbalancers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileLoadbalancers", reflect.TypeOf((*MockELBInterface)(nil).ReconcileLoadbalancers))
}

// RegisterInstanceWithAPIServerELB mocks base method
func (m *MockELBInterface) RegisterInstanceWithAPIServerELB(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterInstanceWithAPIServerELB", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterInstanceWithAPIServerELB indicates an expected call of RegisterInstanceWithAPIServerELB
func (mr *MockELBInterfaceMockRecorder) RegisterInstanceWithAPIServerELB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterInstanceWithAPIServerELB", reflect.TypeOf((*MockELBInterface)(nil).RegisterInstanceWithAPIServerELB), arg0)
}

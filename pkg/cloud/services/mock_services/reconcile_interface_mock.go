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
// Source: sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services (interfaces: MachinePoolReconcileInterface)

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	scope "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	services "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
)

// MockMachinePoolReconcileInterface is a mock of MachinePoolReconcileInterface interface.
type MockMachinePoolReconcileInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMachinePoolReconcileInterfaceMockRecorder
}

// MockMachinePoolReconcileInterfaceMockRecorder is the mock recorder for MockMachinePoolReconcileInterface.
type MockMachinePoolReconcileInterfaceMockRecorder struct {
	mock *MockMachinePoolReconcileInterface
}

// NewMockMachinePoolReconcileInterface creates a new mock instance.
func NewMockMachinePoolReconcileInterface(ctrl *gomock.Controller) *MockMachinePoolReconcileInterface {
	mock := &MockMachinePoolReconcileInterface{ctrl: ctrl}
	mock.recorder = &MockMachinePoolReconcileInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachinePoolReconcileInterface) EXPECT() *MockMachinePoolReconcileInterfaceMockRecorder {
	return m.recorder
}

// ReconcileLaunchTemplate mocks base method.
func (m *MockMachinePoolReconcileInterface) ReconcileLaunchTemplate(arg0 scope.IgnitionScope, arg1 scope.LaunchTemplateScope, arg2 scope.S3Scope, arg3 services.EC2Interface, arg4 services.ObjectStoreInterface, arg5 func() (bool, error), arg6 func() error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileLaunchTemplate", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileLaunchTemplate indicates an expected call of ReconcileLaunchTemplate.
func (mr *MockMachinePoolReconcileInterfaceMockRecorder) ReconcileLaunchTemplate(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileLaunchTemplate", reflect.TypeOf((*MockMachinePoolReconcileInterface)(nil).ReconcileLaunchTemplate), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// ReconcileTags mocks base method.
func (m *MockMachinePoolReconcileInterface) ReconcileTags(arg0 scope.LaunchTemplateScope, arg1 []scope.ResourceServiceToUpdate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileTags", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileTags indicates an expected call of ReconcileTags.
func (mr *MockMachinePoolReconcileInterfaceMockRecorder) ReconcileTags(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileTags", reflect.TypeOf((*MockMachinePoolReconcileInterface)(nil).ReconcileTags), arg0, arg1)
}

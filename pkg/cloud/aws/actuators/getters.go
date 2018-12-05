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

package actuators

var (
	DefaultScopeGetter        ScopeGetter        = ScopeGetterFunc(NewScope)
	DefaultMachineScopeGetter MachineScopeGetter = MachineScopeGetterFunc(NewMachineScope)
)

type ScopeGetter interface {
	GetScope(params ScopeParams) (*Scope, error)
}

type ScopeGetterFunc func(params ScopeParams) (*Scope, error)

func (f ScopeGetterFunc) GetScope(params ScopeParams) (*Scope, error) {
	return f(params)
}

type MachineScopeGetter interface {
	GetMachineScope(params MachineScopeParams) (*MachineScope, error)
}

type MachineScopeGetterFunc func(params MachineScopeParams) (*MachineScope, error)

func (f MachineScopeGetterFunc) GetMachineScope(params MachineScopeParams) (*MachineScope, error) {
	return f(params)
}

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

package certificates

import (
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
)

// Service groups certificate related operations together and allows
// certificate updates to be applied to the actuator scope.
type Service struct {
	scope *actuators.Scope
}

// NewService returns a new certificates service for the given actuators scope.
func NewService(scope *actuators.Scope) *Service {
	return &Service{
		scope: scope,
	}
}

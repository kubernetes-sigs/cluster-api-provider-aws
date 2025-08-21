/*
Copyright 2022 The Kubernetes Authors.

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

// Package mock_ssmiface provides a mock implementation for the SSMAPI interface.
// Run go generate to regenerate this mock.
//
//go:generate ../../../../../hack/tools/bin/mockgen -destination mock_ssmiface.go -package mock_ssmiface sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm SSMAPI
//go:generate /usr/bin/env bash -c "cat ../../../../../hack/boilerplate/boilerplate.generatego.txt mock_ssmiface.go > _mock_ssmiface.go && mv _mock_ssmiface.go mock_ssmiface.go"
package mock_ssmiface //nolint:stylecheck

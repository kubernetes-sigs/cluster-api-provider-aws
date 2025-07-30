/*
Copyright 2020 The Kubernetes Authors.

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

// Package mock_stsiface provides a mock implementation for the STSClient interface.
// Run go generate to regenerate this mock.
//
//go:generate ../../../../../hack/tools/bin/mockgen -destination stsiface_mock_v2.go -package mock_stsiface sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts STSClient
//go:generate /usr/bin/env bash -c "cat ../../../../../hack/boilerplate/boilerplate.generatego.txt stsiface_mock_v2.go > _stsiface_mock_v2.go && mv _stsiface_mock_v2.go stsiface_mock_v2.go"
//go:generate ../../../../../hack/tools/bin/mockgen -destination stsiface_mock_v1.go -package mock_stsiface github.com/aws/aws-sdk-go/service/sts/stsiface STSAPI
//go:generate /usr/bin/env bash -c "cat ../../../../../hack/boilerplate/boilerplate.generatego.txt stsiface_mock_v1.go > _stsiface_mock_v1.go && mv _stsiface_mock_v1.go stsiface_mock_v1.go"
package mock_stsiface //nolint:stylecheck

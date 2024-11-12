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

// Package mocks provides a way to generate mock objects for OCM and ClusterScoper services.
//
//go:generate ../../hack/tools/bin/mockgen -destination capa_clusterscoper_mock.go -package mocks sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud ClusterScoper
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt capa_clusterscoper_mock.go > _capa_clusterscoper_mock.go && mv _capa_clusterscoper_mock.go capa_clusterscoper_mock.go"
//go:generate ../../hack/tools/bin/mockgen -destination ocm_client_mock.go -package mocks sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa OCMClient
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt ocm_client_mock.go > _ocm_client_mock.go && mv _ocm_client_mock.go ocm_client_mock.go"
package mocks

// +build integration

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

package integration_test

import (
	. "github.com/onsi/ginkgo"
	"sigs.k8s.io/cluster-api/util"

	"sigs.k8s.io/cluster-api/test/helpers/kind"
)

const (
	kindTimeout = 5 * 60
	// controllerNamespace = "aws-provider-system"
	// controllerName      = "aws-provider-controller-manager"
)

var _ = Describe("Metacluster", func() {
	var kindCluster kind.Cluster
	BeforeEach(func() {
		kindCluster = kind.Cluster{
			Name: "capa-test-" + util.RandomString(6),
		}
		kindCluster.Setup()
		// client = kindCluster.KubeClient()
	}, kindTimeout)

	AfterEach(func() {
		kindCluster.Teardown()
	})

	// TODO: validate that the controller-manager is deployed and the
	// types are available
})

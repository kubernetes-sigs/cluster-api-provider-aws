// +build e2e

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

package unmanaged

import (
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/cluster-api/test/framework"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

var (
	e2eCtx *shared.E2EContext
)

func init() {
	e2eCtx = shared.NewE2EContext()
	shared.CreateDefaultFlags(e2eCtx)
	SetDefaultEventuallyTimeout(20 * time.Minute)
	SetDefaultEventuallyPollingInterval(10 * time.Second)
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "capa-e2e", []ginkgo.Reporter{framework.CreateJUnitReporterForProw(e2eCtx.Settings.ArtifactFolder)})
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	return shared.Node1BeforeSuite(e2eCtx)
}, func(data []byte) {
	shared.AllNodesBeforeSuite(e2eCtx, data)
})

var _ = ginkgo.SynchronizedAfterSuite(
	func() {
		shared.AllNodesAfterSuite(e2eCtx)
	},
	func() {
		shared.Node1AfterSuite(e2eCtx)
	},
)

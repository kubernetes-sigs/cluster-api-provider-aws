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

package conformance

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/cluster-api/test/framework"

	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

var (
	e2eCtx *shared.E2EContext
)

func init() {
	e2eCtx = shared.NewE2EContextWithFlags(initScheme)
}

func TestE2EConformance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "capa-e2e-conformance", []Reporter{framework.CreateJUnitReporterForProw(e2eCtx.ArtifactFolder)})
}

var _ = SynchronizedBeforeSuite(func() []byte {
	return shared.Node1BeforeSuite(e2eCtx)
}, func(data []byte) {
	shared.AllNodesBeforeSuite(e2eCtx, data)
})

var _ = SynchronizedAfterSuite(func() {
	shared.Node1AfterSuite(e2eCtx)
}, func() {
	shared.AllNodesAfterSuite(e2eCtx)
})

// initScheme creates a new GVK scheme
func initScheme() *runtime.Scheme {
	sc := runtime.NewScheme()
	framework.TryAddDefaultSchemes(sc)
	_ = v1alpha3.AddToScheme(sc)
	_ = clientgoscheme.AddToScheme(sc)
	return sc
}

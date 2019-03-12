package operators

import (
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
)

var _ = g.Describe("[Feature:Operators] Cluster autoscaler operator deployment should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(isDeploymentAvailable(client, "cluster-autoscaler-operator")).To(o.BeTrue())
	})
})

var _ = g.Describe("[Feature:Operators] Cluster autoscaler cluster operator status should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(isStatusAvailable(client, "cluster-autoscaler")).To(o.BeTrue())
	})
})

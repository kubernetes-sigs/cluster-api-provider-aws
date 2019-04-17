package operators

import (
	"fmt"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
)

var (
	deploymentDeprecatedName = "clusterapi-manager-controllers"
)

var _ = g.Describe("[Feature:Operators] Machine API operator deployment should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(isDeploymentAvailable(client, "machine-api-operator")).To(o.BeTrue())
	})

	g.It("reconcile controllers deployment", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		deploymentName := "machine-api-controllers"
		initialDeployment, err := getDeployment(client, deploymentName)
		if err != nil {
			initialDeployment, err = getDeployment(client, deploymentDeprecatedName)
			o.Expect(err).NotTo(o.HaveOccurred())
			deploymentName = deploymentDeprecatedName
		}

		g.By(fmt.Sprintf("checking deployment %q is available", deploymentName))
		o.Expect(isDeploymentAvailable(client, deploymentName)).To(o.BeTrue())

		g.By(fmt.Sprintf("deleting deployment %q", deploymentName))
		err = deleteDeployment(client, initialDeployment)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By(fmt.Sprintf("checking deployment %q is available again", deploymentName))
		o.Expect(isDeploymentAvailable(client, deploymentName)).To(o.BeTrue())
	})
})

var _ = g.Describe("[Feature:Operators] Machine API cluster operator status should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(isStatusAvailable(client, "machine-api")).To(o.BeTrue())
	})
})

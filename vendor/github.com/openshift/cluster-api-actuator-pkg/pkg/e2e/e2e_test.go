package e2e

import (
	"testing"

	"github.com/golang/glog"
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	osconfigv1 "github.com/openshift/api/config/v1"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	caov1alpha1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis"
	healthcheckingv1alpha1 "github.com/openshift/machine-api-operator/pkg/apis/healthchecking/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"

	_ "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/autoscaler"
	_ "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/infra"
	_ "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/machinehealthcheck"
	_ "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/operators"
)

func init() {
	if err := mapiv1beta1.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}

	if err := caov1alpha1.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}

	if err := osconfigv1.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}

	if err := healthcheckingv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}
}

func TestE2E(t *testing.T) {
	o.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Machine Suite")
}

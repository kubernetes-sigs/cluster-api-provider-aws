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

package e2e_test

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"

	"fmt"

	capa "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/util/kind"
	capi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clientset "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
)

const (
	kindTimeout      = 5 * 60
	namespace        = "capa-test"
	clusterName      = "capa-test-cluster"
	controlPlaneName = "capa-test-control-plane"

	stackName   = "cluster-api-provider-aws-sigs-k8s-io"
	keyPairName = "cluster-api-provider-aws-sigs-k8s-io"
)

var (
	region      string
	credFile    = flag.String("credFile", "", "path to an AWS credentials file")
	clusterYAML = flag.String("clusterYAML", "", "path to the YAML for the cluster we're creating")
	machineYAML = flag.String("machineYAML", "", "path to the YAML describing the control plane we're creating")
)

func init() {
	if region = os.Getenv("AWS_DEFAULT_REGION"); region == "" {
		if region = os.Getenv("AWS_REGION"); region == "" {
			region = "us-east-1"
		}
	}
}

var _ = Describe("AWS", func() {
	var (
		cluster kind.Cluster
		client  *clientset.Clientset
	)
	BeforeEach(func() {
		fmt.Fprintf(GinkgoWriter, "running in AWS region: %s\n", region)
		cluster.Setup()
		cfg := cluster.RestConfig()
		var err error
		client, err = clientset.NewForConfig(cfg)
		Expect(err).To(BeNil())
	}, kindTimeout)

	AfterEach(func() {
		cluster.Teardown()
	})

	Describe("control plane node", func() {
		It("should be running", func() {

			sess := getSession()
			accountID := getAccountID(sess)

			createKeyPair(sess)
			createIAMRoles(sess, accountID)

			createNamespace(cluster.KubeClient())

			By("Creating a cluster")
			clusterapi := client.ClusterV1alpha1().Clusters(namespace)
			_, err := clusterapi.Create(makeCluster())
			Expect(err).To(BeNil())

			By("Creating a machine")
			machineapi := client.ClusterV1alpha1().Machines(namespace)
			_, err = machineapi.Create(makeMachine())
			Expect(err).To(BeNil())

			Eventually(
				func() (*capa.AWSMachineProviderStatus, error) {
					machine, err := machineapi.Get(controlPlaneName, metav1.GetOptions{})
					if err != nil {
						return nil, err
					}
					return capa.MachineStatusFromProviderStatus(machine.Status.ProviderStatus)
				},
				10*time.Minute, 15*time.Second,
			).Should(beHealthy())
		})
	})
})

func beHealthy() types.GomegaMatcher {
	return PointTo(
		MatchFields(IgnoreExtras, Fields{
			"InstanceState": PointTo(Equal(capa.InstanceStateRunning)),
		}),
	)
}

func makeCluster() *capi.Cluster {
	yaml, err := ioutil.ReadFile(*clusterYAML)
	Expect(err).To(BeNil())

	deserializer := serializer.NewCodecFactory(getScheme()).UniversalDeserializer()
	cluster := &capi.Cluster{}
	obj, _, err := deserializer.Decode(yaml, nil, cluster)
	Expect(err).To(BeNil())
	cluster, ok := obj.(*capi.Cluster)
	Expect(ok).To(BeTrue(), "Wanted cluster, got %T", obj)

	cluster.ObjectMeta.Name = clusterName

	awsSpec, err := capa.ClusterConfigFromProviderSpec(cluster.Spec.ProviderSpec)
	Expect(err).To(BeNil())
	awsSpec.SSHKeyName = keyPairName
	awsSpec.Region = region
	cluster.Spec.ProviderSpec.Value, err = capa.EncodeClusterSpec(awsSpec)
	Expect(err).To(BeNil())

	return cluster
}

func makeMachine() *capi.Machine {
	yaml, err := ioutil.ReadFile(*machineYAML)
	Expect(err).To(BeNil())

	deserializer := serializer.NewCodecFactory(getScheme()).UniversalDeserializer()
	obj, _, err := deserializer.Decode(yaml, nil, &capi.MachineList{})
	Expect(err).To(BeNil())
	machineList, ok := obj.(*capi.MachineList)
	Expect(ok).To(BeTrue(), "Wanted machine, got %T", obj)

	machines := machine.GetControlPlaneMachines(machineList)
	Expect(machines).NotTo(BeEmpty())

	machine := machines[0]
	machine.ObjectMeta.Name = controlPlaneName
	machine.ObjectMeta.Labels[capi.MachineClusterLabelName] = clusterName

	awsSpec, err := actuators.MachineConfigFromProviderSpec(nil, machine.Spec.ProviderSpec)
	Expect(err).To(BeNil())
	awsSpec.KeyName = keyPairName
	machine.Spec.ProviderSpec.Value, err = capa.EncodeMachineSpec(awsSpec)
	Expect(err).To(BeNil())

	return machine
}

func getScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	capi.SchemeBuilder.AddToScheme(s)
	capa.SchemeBuilder.AddToScheme(s)
	return s
}

func createNamespace(client kubernetes.Interface) {
	_, err := client.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	Expect(err).To(BeNil())
}

func getAccountID(prov client.ConfigProvider) string {
	stsSvc := sts.NewService(awssts.New(prov))
	accountID, err := stsSvc.AccountID()
	Expect(err).To(BeNil())
	return accountID
}

func createIAMRoles(prov client.ConfigProvider, accountID string) {
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.ReconcileBootstrapStack(stackName, accountID),
	).To(Succeed())
}

func createKeyPair(prov client.ConfigProvider) {
	ec2c := ec2.New(prov)
	_, err := ec2c.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	if code, _ := awserrors.Code(err); code != "InvalidKeyPair.Duplicate" {
		Expect(err).To(BeNil())
	}
}

func getSession() client.ConfigProvider {
	creds := credentials.NewCredentials(&credentials.SharedCredentialsProvider{
		Filename: *credFile,
	})
	sess, err := session.NewSession(aws.NewConfig().WithCredentials(creds).WithRegion(region))
	Expect(err).To(BeNil())
	return sess
}

// +build integration

/*
Copyright 2019 The Kubernetes Authors.

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

package commons

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
	bootstrapv1 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/sts"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiFlag "sigs.k8s.io/cluster-api/test/helpers/flag"
	"sigs.k8s.io/cluster-api/test/helpers/kind"
	"sigs.k8s.io/cluster-api/test/helpers/scheme"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	CapiNamespace       = "capi-system"
	CapiDeploymentName  = "capi-controller-manager"
	CabpkNamespace      = "cabpk-system"
	CabpkDeploymentName = "cabpk-controller-manager"
	CapaNamespace       = "capa-system"
	CapaDeploymentName  = "capa-controller-manager"
)

const AWSCredentialsTemplate = `[default]
aws_access_key_id = {{ .AccessKeyID }}
aws_secret_access_key = {{ .SecretAccessKey }}
region = {{ .Region }}
`

var (
	kustomizeBinary = capiFlag.DefineOrLookupStringFlag("kustomizeBinary", "kustomize", "path to the kustomize binary")
	capaComponents  = capiFlag.DefineOrLookupStringFlag("capaComponents", "", "capa components to load")
	Region          *string
)

type awsCredential struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

func DeployCAPAComponents(kindCluster kind.Cluster, suiteTmpDir string, accessKey *iam.AccessKey) {
	if capaComponents != nil && *capaComponents != "" {
		ApplyManifests(kindCluster, capaComponents)
		return
	}

	fmt.Fprintf(GinkgoWriter, "Generating CAPA manifests\n")

	// Build the manifests using kustomize
	capaManifests, err := exec.Command(*kustomizeBinary, "build", "../../config/default").Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(GinkgoWriter, "Error: %s\n", string(exitError.Stderr))
		}
	}
	Expect(err).NotTo(HaveOccurred())

	// envsubst the credentials
	Expect(err).NotTo(HaveOccurred())
	b64credentials := generateB64Credentials(accessKey)
	os.Setenv("AWS_B64ENCODED_CREDENTIALS", b64credentials)
	manifestsContent := os.ExpandEnv(string(capaManifests))

	// write out the manifests
	manifestFile := path.Join(suiteTmpDir, "infrastructure-components.yaml")
	Expect(ioutil.WriteFile(manifestFile, []byte(manifestsContent), 0644)).To(Succeed())

	// apply generated manifests
	ApplyManifests(kindCluster, &manifestFile)
}

func GetSession() client.ConfigProvider {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	Expect(err).NotTo(HaveOccurred())
	return sess
}

func GetAccountID(prov client.ConfigProvider) string {
	stsSvc := sts.NewService(awssts.New(prov))
	accountID, err := stsSvc.AccountID()
	Expect(err).NotTo(HaveOccurred())
	return accountID
}

func CreateIAMRoles(prov client.ConfigProvider, accountID string, stackName string) {
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.ReconcileBootstrapStack(stackName, accountID, "aws"),
	).To(Succeed())
}

func DeleteIAMRoles(prov client.ConfigProvider, stackName string) {
	cfnSvc := cloudformation.NewService(cfn.New(prov))
	Expect(
		cfnSvc.DeleteStack(stackName),
	).To(Succeed())
}

func CreateKeyPair(prov client.ConfigProvider, keyPairName string) {
	ec2c := ec2.New(prov)
	_, err := ec2c.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(keyPairName)})
	if code, _ := awserrors.Code(err); code != "InvalidKeyPair.Duplicate" {
		Expect(err).NotTo(HaveOccurred())
	}
}

func LoadManagerImage(kindCluster kind.Cluster, managerImage *string) {
	if managerImage != nil && *managerImage != "" {
		kindCluster.LoadImage(*managerImage)
	}
}

func ApplyManifests(kindCluster kind.Cluster, manifests *string) {
	Expect(manifests).ToNot(BeNil())
	fmt.Fprintf(GinkgoWriter, "Applying manifests for %s\n", *manifests)
	Expect(*manifests).ToNot(BeEmpty())
	kindCluster.ApplyYAML(*manifests)
}

func generateB64Credentials(accessKey *iam.AccessKey) string {
	creds := awsCredential{
		Region:          *Region,
		AccessKeyID:     *accessKey.AccessKeyId,
		SecretAccessKey: *accessKey.SecretAccessKey,
	}

	tmpl, err := template.New("AWS Credentials").Parse(AWSCredentialsTemplate)
	Expect(err).NotTo(HaveOccurred())

	var profile bytes.Buffer
	Expect(tmpl.Execute(&profile, creds)).To(Succeed())

	encCreds := base64.StdEncoding.EncodeToString(profile.Bytes())
	return encCreds
}

func SetupScheme() *runtime.Scheme {
	s := scheme.SetupScheme()
	Expect(bootstrapv1.AddToScheme(s)).To(Succeed())
	Expect(infrav1.AddToScheme(s)).To(Succeed())
	return s
}

func RetrieveAllLogs(kindCluster kind.Cluster, kindClient crclient.Client) {
	capiLogs := retrieveLogs(CapiNamespace, CapiDeploymentName, kindCluster, kindClient)
	cabpkLogs := retrieveLogs(CabpkNamespace, CabpkDeploymentName, kindCluster, kindClient)
	capaLogs := retrieveLogs(CapaNamespace, CapaDeploymentName, kindCluster, kindClient)

	// If running in prow, output the logs to the artifacts path
	artifactPath, exists := os.LookupEnv("ARTIFACTS")
	if exists {
		ioutil.WriteFile(path.Join(artifactPath, "capi.log"), []byte(capiLogs), 0644)
		ioutil.WriteFile(path.Join(artifactPath, "cabpk.log"), []byte(cabpkLogs), 0644)
		ioutil.WriteFile(path.Join(artifactPath, "capa.log"), []byte(capaLogs), 0644)
		return
	}

	fmt.Fprintf(GinkgoWriter, "CAPI Logs:\n%s\n", capiLogs)
	fmt.Fprintf(GinkgoWriter, "CABPK Logs:\n%s\n", cabpkLogs)
	fmt.Fprintf(GinkgoWriter, "CAPA Logs:\n%s\n", capaLogs)
}

func retrieveLogs(namespace, deploymentName string, kindCluster kind.Cluster, kindClient crclient.Client) string {
	deployment := &appsv1.Deployment{}
	Expect(kindClient.Get(context.TODO(), crclient.ObjectKey{Namespace: namespace, Name: deploymentName}, deployment)).To(Succeed())

	pods := &corev1.PodList{}

	selector, err := metav1.LabelSelectorAsMap(deployment.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())

	Expect(kindClient.List(context.TODO(), pods, crclient.InNamespace(namespace), crclient.MatchingLabels(selector))).To(Succeed())
	Expect(pods.Items).NotTo(BeEmpty())

	clientset, err := kubernetes.NewForConfig(kindCluster.RestConfig())
	Expect(err).NotTo(HaveOccurred())

	podLogs, err := clientset.CoreV1().Pods(namespace).GetLogs(pods.Items[0].Name, &corev1.PodLogOptions{Container: "manager"}).Stream()
	Expect(err).NotTo(HaveOccurred())
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	Expect(err).NotTo(HaveOccurred())

	return buf.String()
}

func CreateNamespace(namespace string, kindClient crclient.Client) {
	fmt.Fprintf(GinkgoWriter, "creating namespace %q\n", namespace)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	Expect(kindClient.Create(context.TODO(), ns)).To(Succeed())
}

func MakeCluster(namespace, name, awsClusterName string, kindClient crclient.Client) error {
	fmt.Fprintf(GinkgoWriter, "Creating Cluster %s/%s\n", namespace, name)
	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: clusterv1.ClusterSpec{
			ClusterNetwork: &clusterv1.ClusterNetwork{
				Pods: &clusterv1.NetworkRanges{
					CIDRBlocks: []string{"192.168.0.0/16"},
				},
			},
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSCluster",
				APIVersion: infrav1.GroupVersion.String(),
				Name:       awsClusterName,
				Namespace:  namespace,
			},
		},
	}
	return kindClient.Create(context.TODO(), cluster)
}

func MakeAWSCluster(namespace, name, keyPairName string, kindClient crclient.Client) error {
	fmt.Fprintf(GinkgoWriter, "Creating AWSCluster %s/%s\n", namespace, name)
	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSClusterSpec{
			Region:     *Region,
			SSHKeyName: keyPairName,
		},
	}
	return kindClient.Create(context.TODO(), awsCluster)
}

func IsClusterProvisioningCompleted(clusterName, namespace string, kindClient crclient.Client, waitMins time.Duration) bool {
	fmt.Fprintf(GinkgoWriter, "Ensuring Cluster: %s is provisioned under Namepsace: %s\n", clusterName, namespace)
	clusterProvisioned := false
	endTime := time.Now().Add(waitMins * time.Minute)
	for time.Now().Before(endTime) {
		cluster := &clusterv1.Cluster{}
		if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: clusterName}, cluster); err == nil {
			if cluster.Status.Phase == string(clusterv1.ClusterPhaseProvisioned) {
				clusterProvisioned = true
				break
			}
		}
		time.Sleep(15 * time.Second)
	}
	return clusterProvisioned
}

func DeleteCluster(namespace, clusterName string, kindClient crclient.Client) {
	fmt.Fprintf(GinkgoWriter, "Deleting Cluster %s/%s\n", namespace, clusterName)
	Expect(kindClient.Delete(context.TODO(), &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
	})).To(Succeed())

	Eventually(
		func() bool {
			cluster := &clusterv1.Cluster{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: clusterName}, cluster); err != nil {
				if apierrors.IsNotFound(err) {
					return true
				}
				return false
			}
			return false
		},
		20*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

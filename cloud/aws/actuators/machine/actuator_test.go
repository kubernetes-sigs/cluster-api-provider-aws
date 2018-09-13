package machine

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/fake"

	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go/aws"
	fakeawsclient "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/client/fake"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/test/utils"
)

const (
	controllerLogName = "awsMachine"
	defaultLogLevel   = "info"

	defaultNamespace         = "default"
	defaultAvailabilityZone  = "us-east-1a"
	region                   = "us-east-1"
	awsCredentialsSecretName = "aws-credentials-secret"
	userDataSecretName       = "aws-actuator-user-data-secret"

	keyName     = "aws-actuator-key-name"
	clusterName = "aws-actuator-cluster"
)

const userDataBlob = `#cloud-config
write_files:
- path: /root/node_bootstrap/node_settings.yaml
  owner: 'root:root'
  permissions: '0640'
  content: |
    node_config_name: node-config-master
runcmd:
- [ cat, /root/node_bootstrap/node_settings.yaml]
`

func testMachineAPIResources(clusterID string) (*clusterv1.Machine, *clusterv1.Cluster, *apiv1.Secret, *apiv1.Secret, error) {
	awsCredentialsSecret := utils.GenerateAwsCredentialsSecretFromEnv(awsCredentialsSecretName, defaultNamespace)

	userDataSecret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userDataSecretName,
			Namespace: defaultNamespace,
		},
		Data: map[string][]byte{
			userDataSecretKey: []byte(userDataBlob),
		},
	}

	machinePc := &providerconfigv1.AWSMachineProviderConfig{
		AMI: providerconfigv1.AWSResourceReference{
			ID: aws.String("ami-a9acbbd6"),
		},
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: providerconfigv1.Placement{
			Region:           region,
			AvailabilityZone: defaultAvailabilityZone,
		},
		Subnet: providerconfigv1.AWSResourceReference{
			ID: aws.String("subnet-0e56b13a64ff8a941"),
		},
		IAMInstanceProfile: &providerconfigv1.AWSResourceReference{
			ID: aws.String("openshift_master_launch_instances"),
		},
		KeyName: aws.String(keyName),
		UserDataSecret: &corev1.LocalObjectReference{
			Name: userDataSecretName,
		},
		Tags: []providerconfigv1.TagSpecification{
			{Name: "openshift-node-group-config", Value: "node-config-master"},
			{Name: "host-type", Value: "master"},
			{Name: "sub-host-type", Value: "default"},
		},
		SecurityGroups: []providerconfigv1.AWSResourceReference{
			{ID: aws.String("sg-00868b02fbe29de17")}, // aws-actuator
			{ID: aws.String("sg-0a4658991dc5eb40a")}, // aws-actuator_master
			{ID: aws.String("sg-009a70e28fa4ba84e")}, // aws-actuator_master_k8s
			{ID: aws.String("sg-07323d56fb932c84c")}, // aws-actuator_infra
			{ID: aws.String("sg-08b1ffd32874d59a2")}, // aws-actuator_infra_k8s
		},
		PublicIP: aws.Bool(true),
	}

	var buf bytes.Buffer
	if err := providerconfigv1.Encoder.Encode(machinePc, &buf); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("encoding failed: %v", err)
	}

	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-actuator-testing-machine",
			Namespace: defaultNamespace,
			Labels: map[string]string{
				providerconfigv1.ClusterNameLabel: clusterID,
				providerconfigv1.MachineRoleLabel: "infra",
				providerconfigv1.MachineTypeLabel: "master",
			},
		},

		Spec: clusterv1.MachineSpec{
			ProviderConfig: clusterv1.ProviderConfig{
				Value: &runtime.RawExtension{Raw: buf.Bytes()},
			},
		},
	}

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID,
			Namespace: defaultNamespace,
		},
	}

	return machine, cluster, awsCredentialsSecret, userDataSecret, nil
}

func TestCreateAndDeleteMachine(t *testing.T) {

	// kube client is needed to fetch aws credentials:
	// - kubeClient.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	// cluster client for updating machine statues
	// - clusterClient.ClusterV1alpha1().Machines(machineCopy.Namespace).UpdateStatus(machineCopy)

	machine, cluster, awsCredentialsSecret, userDataSecret, err := testMachineAPIResources(clusterName)
	if err != nil {
		t.Fatal(err)
	}

	fakeKubeClient := kubernetesfake.NewSimpleClientset(awsCredentialsSecret, userDataSecret)
	fakeClient := fake.NewSimpleClientset(machine)
	logger := log.WithField("controller", controllerLogName)

	params := ActuatorParams{
		ClusterClient:    fakeClient,
		KubeClient:       fakeKubeClient,
		AwsClientBuilder: fakeawsclient.NewClient,
		Logger:           logger,
	}

	actuator, err := NewActuator(params)
	if err != nil {
		t.Fatalf("Could not create Openstack machine actuator: %v", err)
	}

	// Create a machine
	if err := actuator.Create(cluster, machine); err != nil {
		t.Errorf("Unable to create instance for machine: %v", err)
	}

	// Get the machine
	if exists, err := actuator.Exists(cluster, machine); err != nil || !exists {
		t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
	} else {
		t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
	}

	// TODO(jchaloup): Wait until the machine is ready

	// Update a machine
	if err := actuator.Update(cluster, machine); err != nil {
		t.Errorf("Unable to create instance for machine: %v", err)
	}

	// Get the machine
	if exists, err := actuator.Exists(cluster, machine); err != nil || !exists {
		t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
	} else {
		t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
	}

	// Delete a machine
	if err := actuator.Delete(cluster, machine); err != nil {
		t.Errorf("Unable to delete instance for machine: %v", err)
	}
}

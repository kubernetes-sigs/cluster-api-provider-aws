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

package main

// Tests individual AWS actuator actions. This is meant to be executed
// in a machine that has access to AWS either as an instance with the right role
// or creds in ~/.aws/credentials

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/ghodss/yaml"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"

	"text/template"

	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	instanceIDAnnotation = "cluster-operator.openshift.io/aws-instance-id"
	ami                  = "ami-03f6257a"
	region               = "eu-west-1"
	size                 = "t1.micro"

	pollInterval           = 5 * time.Second
	timeoutPoolAWSInterval = 10 * time.Minute
)

func usage() {
	fmt.Printf("Usage: %s\n\n", os.Args[0])
}

var rootCmd = &cobra.Command{
	Use:   "aws-actuator-test",
	Short: "Test for Cluster API AWS actuator",
}

func createCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create machine instance for specified cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkFlags(cmd); err != nil {
				return err
			}
			cluster, machine, awsCredentials, userData, err := readClusterResources(
				&manifestParams{
					ClusterID: cmd.Flag("environment-id").Value.String(),
				},
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("aws-credentials").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := createActuator(machine, awsCredentials, userData)
			result, err := actuator.CreateMachine(cluster, machine)
			if err != nil {
				return err
			}
			fmt.Printf("Machine creation was successful! InstanceID: %s\n", *result.InstanceId)
			return nil
		},
	}
}

func deleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete machine instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkFlags(cmd); err != nil {
				return err
			}
			cluster, machine, awsCredentials, userData, err := readClusterResources(
				&manifestParams{
					ClusterID: cmd.Flag("environment-id").Value.String(),
				},
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("aws-credentials").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := createActuator(machine, awsCredentials, userData)
			err = actuator.DeleteMachine(cluster, machine)
			if err != nil {
				return err
			}
			fmt.Printf("Machine delete operation was successful.\n")
			return nil
		},
	}
}

func existsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "exists",
		Short: "Determine if underlying machine instance exists",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkFlags(cmd); err != nil {
				return err
			}
			cluster, machine, awsCredentials, userData, err := readClusterResources(
				&manifestParams{
					ClusterID: cmd.Flag("environment-id").Value.String(),
				},
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("aws-credentials").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := createActuator(machine, awsCredentials, userData)
			exists, err := actuator.Exists(cluster, machine)
			if err != nil {
				return err
			}
			if exists {
				fmt.Printf("Underlying machine's instance exists.\n")
			} else {
				fmt.Printf("Underlying machine's instance not found.\n")
			}
			return nil
		},
	}
}

func readMachineSetManifest(manifestParams *manifestParams, manifestLoc string) (*clusterv1.MachineSet, error) {
	machineset := &clusterv1.MachineSet{}
	manifestBytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", manifestLoc, err)
	}

	t, err := template.New("machinesetuserdata").Parse(string(manifestBytes))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, *manifestParams)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buf.Bytes(), &machineset); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v: %v", manifestLoc, err)
	}

	return machineset, nil
}

func readMachineManifest(manifestParams *manifestParams, manifestLoc string) (*clusterv1.Machine, error) {
	machine := &clusterv1.Machine{}
	manifestBytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", manifestLoc, err)
	}

	t, err := template.New("machineuserdata").Parse(string(manifestBytes))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, *manifestParams)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buf.Bytes(), &machine); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v: %v", manifestLoc, err)
	}

	return machine, nil
}

func readClusterManifest(manifestLoc string) (*clusterv1.Cluster, error) {
	cluster := &clusterv1.Cluster{}
	bytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", manifestLoc, err)
	}

	if err = yaml.Unmarshal(bytes, &cluster); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v: %v", manifestLoc, err)
	}

	return cluster, nil
}

func readSecretManifest(manifestLoc string) (*apiv1.Secret, error) {
	secret := &apiv1.Secret{}
	bytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", manifestLoc, err)
	}
	if err = yaml.Unmarshal(bytes, &secret); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v: %v", manifestLoc, err)
	}
	return secret, nil
}

const workerUserDataBlob = `#!/bin/bash

cat <<HEREDOC > /root/user-data.sh
#!/bin/bash

cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kube*
EOF
setenforce 0
yum install -y kubelet-1.11.3 kubeadm-1.11.3 --disableexcludes=kubernetes

cat <<EOF > /etc/default/kubelet
KUBELET_KUBEADM_EXTRA_ARGS=--cgroup-driver=systemd
EOF

kubeadm join {{ .MasterIP }}:8443 --token 2iqzqm.85bs0x6miyx1nm7l --discovery-token-unsafe-skip-ca-verification

HEREDOC

bash /root/user-data.sh > /root/user-data.logs
`

type userDataParams struct {
	MasterIP string
}

func generateWorkerUserData(masterIP string, workerUserDataSecret *apiv1.Secret) (*apiv1.Secret, error) {
	params := userDataParams{
		MasterIP: masterIP,
	}
	t, err := template.New("workeruserdata").Parse(workerUserDataBlob)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, params)
	if err != nil {
		return nil, err
	}

	secret := workerUserDataSecret.DeepCopy()
	secret.Data["userData"] = []byte(buf.String())

	return secret, nil
}

// TestConfig stores clients for managing various resources
type TestConfig struct {
	KubeClient *kubernetes.Clientset
	CAPIClient *clientset.Clientset
}

func createCluster(testConfig *TestConfig, cluster *clusterv1.Cluster) error {
	glog.Infof("Creating %q cluster...", strings.Join([]string{cluster.Namespace, cluster.Name}, "/"))
	if _, err := testConfig.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{}); err != nil {
		if _, err := testConfig.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Create(cluster); err != nil {
			return fmt.Errorf("unable to create cluster: %v", err)
		}
	}

	return nil
}

func createMachineSet(testConfig *TestConfig, machineset *clusterv1.MachineSet) error {
	glog.Infof("Creating %q machineset...", strings.Join([]string{machineset.Namespace, machineset.Name}, "/"))
	if _, err := testConfig.CAPIClient.ClusterV1alpha1().MachineSets(machineset.Namespace).Get(machineset.Name, metav1.GetOptions{}); err != nil {
		if _, err := testConfig.CAPIClient.ClusterV1alpha1().MachineSets(machineset.Namespace).Create(machineset); err != nil {
			return fmt.Errorf("unable to create machineset: %v", err)
		}
	}

	return nil
}

func createSecret(testConfig *TestConfig, secret *apiv1.Secret) error {
	glog.Infof("Creating %q secret...", strings.Join([]string{secret.Namespace, secret.Name}, "/"))
	if _, err := testConfig.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{}); err != nil {
		if _, err := testConfig.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret); err != nil {
			return fmt.Errorf("unable to create secret: %v", err)
		}
	}

	return nil
}

func createNamespace(testConfig *TestConfig, namespace string) error {
	nsObj := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	glog.Infof("Creating %q namespace...", nsObj.Name)
	if _, err := testConfig.KubeClient.CoreV1().Namespaces().Get(nsObj.Name, metav1.GetOptions{}); err != nil {
		if _, err := testConfig.KubeClient.CoreV1().Namespaces().Create(nsObj); err != nil {
			return fmt.Errorf("unable to create namespace: %v", err)
		}
	}

	return nil
}

func bootstrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap kubernetes cluster with kubeadm",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestsDir := cmd.Flag("manifests").Value.String()
			if manifestsDir == "" {
				return fmt.Errorf("--manifests needs to be set")
			}

			machinePrefix := cmd.Flag("environment-id").Value.String()

			if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
				return fmt.Errorf("AWS_ACCESS_KEY_ID env needs to be set")
			}
			if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
				return fmt.Errorf("AWS_SECRET_ACCESS_KEY env needs to be set")
			}

			glog.Infof("Reading cluster manifest from %v", path.Join(manifestsDir, "cluster.yaml"))
			cluster, err := readClusterManifest(path.Join(manifestsDir, "cluster.yaml"))
			if err != nil {
				return err
			}

			glog.Infof("Reading master machine manifest from %v", path.Join(manifestsDir, "master-machine.yaml"))
			masterMachine, err := readMachineManifest(
				&manifestParams{
					ClusterID: machinePrefix,
				},
				path.Join(manifestsDir, "master-machine.yaml"),
			)
			if err != nil {
				return err
			}

			glog.Infof("Reading master user data manifest from %v", path.Join(manifestsDir, "master-userdata.yaml"))
			masterUserDataSecret, err := readSecretManifest(path.Join(manifestsDir, "master-userdata.yaml"))
			if err != nil {
				return err
			}

			if machinePrefix != "" {
				masterMachine.Name = machinePrefix + "-" + masterMachine.Name
			}

			var awsCredentialsSecret *apiv1.Secret
			if cmd.Flag("aws-credentials").Value.String() != "" {
				glog.Infof("Reading aws credentials manifest from %v", cmd.Flag("aws-credentials").Value.String())
				awsCredentialsSecret, err = readSecretManifest(cmd.Flag("aws-credentials").Value.String())
				if err != nil {
					return err
				}
			}

			glog.Infof("Creating master machine")
			actuator := createActuator(masterMachine, awsCredentialsSecret, masterUserDataSecret)
			result, err := actuator.CreateMachine(cluster, masterMachine)
			if err != nil {
				return err
			}

			glog.Infof("Master machine created with ipv4: %v, InstanceId: %v", *result.PrivateIpAddress, *result.InstanceId)

			masterMachinePublicDNS := ""
			masterMachinePrivateIP := ""
			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Info("Waiting for master machine PublicDNS")
				result, err := actuator.Describe(cluster, masterMachine)
				if err != nil {
					glog.Info(err)
					return false, nil
				}

				glog.Infof("PublicDnsName: %v\n", *result.PublicDnsName)
				if *result.PublicDnsName == "" {
					return false, nil
				}

				masterMachinePublicDNS = *result.PublicDnsName
				masterMachinePrivateIP = *result.PrivateIpAddress
				return true, nil
			})
			if err != nil {
				glog.Errorf("Unable to get DNS name: %v", err)
			}

			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Infof("Pulling kubeconfig from %v:8443", masterMachinePublicDNS)
				output, err := cmdRun("ssh", fmt.Sprintf("ec2-user@%v", masterMachinePublicDNS), "sudo cat /etc/kubernetes/admin.conf")
				if err != nil {
					glog.Infof("Unable to pull kubeconfig: %v, %v", err, string(output))
					return false, nil
				}

				f, err := os.Create("kubeconfig")
				if err != nil {
					return false, err
				}

				if _, err = f.Write(output); err != nil {
					f.Close()
					return false, err
				}
				f.Close()

				return true, nil
			})
			if err != nil {
				glog.Errorf("Unable to create kubeconfig: %v", err)
			}

			glog.Infof("Running kubectl --kubeconfig=kubeconfig config set-cluster kubernetes --server=https://%v:8443", masterMachinePublicDNS)
			if _, err := cmdRun("kubectl", "--kubeconfig=kubeconfig", "config", "set-cluster", "kubernetes", fmt.Sprintf("--server=https://%v:8443", masterMachinePublicDNS)); err != nil {
				return err
			}

			// Wait until the cluster comes up
			config, err := controller.GetConfig("kubeconfig")
			if err != nil {
				return fmt.Errorf("unable to create config: %v", err)
			}

			kubeClient, err := kubernetes.NewForConfig(config)
			if err != nil {
				glog.Fatalf("Could not create kubernetes client to talk to the apiserver: %v", err)
			}

			capiclient, err := clientset.NewForConfig(config)
			if err != nil {
				glog.Fatalf("Could not create client for talking to the apiserver: %v", err)
			}

			tc := &TestConfig{
				KubeClient: kubeClient,
				CAPIClient: capiclient,
			}

			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Info("Waiting for all nodes to come up")
				nodesList, err := kubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
				if err != nil {
					return false, nil
				}

				nodesReady := true
				for _, node := range nodesList.Items {
					ready := false
					for _, c := range node.Status.Conditions {
						if c.Type != apiv1.NodeReady {
							continue
						}
						ready = true
					}
					glog.Infof("Is node %q ready?: %v\n", node.Name, ready)
					if !ready {
						nodesReady = false
					}
				}

				return nodesReady, nil
			})
			if err != nil {
				glog.Errorf("Failed to get all nodes up and running: %v", err)
			}

			glog.Info("Deploying cluster-api stack")
			glog.Info("Deploying aws credentials")

			if err := createNamespace(tc, "test"); err != nil {
				return err
			}

			awsCredentials := &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-credentials-secret",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"awsAccessKeyId":     []byte(os.Getenv("AWS_ACCESS_KEY_ID")),
					"awsSecretAccessKey": []byte(os.Getenv("AWS_SECRET_ACCESS_KEY")),
				},
			}

			if err := createSecret(tc, awsCredentials); err != nil {
				return err
			}

			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Info("Deploying cluster-api server")
				if output, err := cmdRun("kubectl", "--kubeconfig=kubeconfig", "apply", fmt.Sprintf("-f=%v", path.Join(manifestsDir, "cluster-api-server.yaml")), "--validate=false"); err != nil {
					glog.Infof("Unable to apply %v manifest: %v\n%v", path.Join(manifestsDir, "cluster-api-server.yaml"), err, string(output))
					return false, nil
				}

				return true, nil
			})
			if err != nil {
				glog.Errorf("Unable to deploy cluster-api server %v", err)
			}

			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Info("Deploying cluster-api controllers")
				if output, err := cmdRun("kubectl", "--kubeconfig=kubeconfig", "apply", fmt.Sprintf("-f=%v", path.Join(manifestsDir, "provider-components.yml"))); err != nil {
					glog.Infof("Unable to apply %v manifest: %v\n%v", path.Join(manifestsDir, "provider-components.yml"), err, string(output))
					return false, nil
				}
				return true, nil
			})
			if err != nil {
				glog.Errorf("Unable to deploy cluster-api controllers: %v", err)
			}

			testCluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tb-asg-35",
					Namespace: "test",
				},
				Spec: clusterv1.ClusterSpec{
					ClusterNetwork: clusterv1.ClusterNetworkingConfig{
						Services: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"10.0.0.1/24"},
						},
						Pods: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"10.0.0.1/24"},
						},
						ServiceDomain: "example.com",
					},
				},
			}

			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Infof("Deploying cluster resource")

				if err := createCluster(tc, testCluster); err != nil {
					glog.Infof("Unable to deploy cluster manifest: %v", err)
					return false, nil
				}

				return true, nil
			})
			if err != nil {
				glog.Errorf("Unable to deploy cluster resource %v", err)
			}

			glog.Infof("Reading worker user data manifest from %v", path.Join(manifestsDir, "worker-userdata.yaml"))
			workerUserDataSecret, err := readSecretManifest(path.Join(manifestsDir, "worker-userdata.yaml"))
			if err != nil {
				return err
			}

			glog.Infof("Generating worker machine set user data for master listening at %v", masterMachinePrivateIP)
			workerUserDataSecret, err = generateWorkerUserData(masterMachinePrivateIP, workerUserDataSecret)
			if err != nil {
				return fmt.Errorf("unable to generate worker user data: %v", err)
			}

			if err := createSecret(tc, workerUserDataSecret); err != nil {
				return err
			}

			glog.Infof("Reading worker machine manifest from %v", path.Join(manifestsDir, "worker-machineset.yaml"))
			workerMachineSet, err := readMachineSetManifest(
				&manifestParams{
					ClusterID: machinePrefix,
				},
				path.Join(manifestsDir, "worker-machineset.yaml"),
			)
			if err != nil {
				return err
			}

			if machinePrefix != "" {
				workerMachineSet.Name = machinePrefix + "-" + workerMachineSet.Name
			}

			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Info("Deploying worker machineset")
				if err := createMachineSet(tc, workerMachineSet); err != nil {
					glog.Infof("unable to create machineset: %v", err)
					return false, nil
				}

				return true, nil
			})
			if err != nil {
				glog.Errorf("Unable to deploy worker machineset: %v", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringP("manifests", "", "", "Directory with bootstrapping manifests")
	return cmd
}

func cmdRun(binaryPath string, args ...string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	return cmd.CombinedOutput()
}

func init() {
	rootCmd.PersistentFlags().StringP("machine", "m", "", "Machine manifest")
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "Cluster manifest")
	rootCmd.PersistentFlags().StringP("aws-credentials", "a", "", "Secret manifest with aws credentials")
	rootCmd.PersistentFlags().StringP("userdata", "u", "", "User data manifest")
	cUser, err := user.Current()
	if err != nil {
		rootCmd.PersistentFlags().StringP("environment-id", "p", "", "Directory with bootstrapping manifests")
	} else {
		rootCmd.PersistentFlags().StringP("environment-id", "p", cUser.Username, "Machine prefix, by default set to the current user")
	}

	rootCmd.AddCommand(createCommand())

	rootCmd.AddCommand(deleteCommand())

	rootCmd.AddCommand(existsCommand())

	rootCmd.AddCommand(bootstrapCommand())
}

type manifestParams struct {
	ClusterID string
}

func readClusterResources(manifestParams *manifestParams, clusterLoc, machineLoc, awsCredentialSecretLoc, userDataLoc string) (*clusterv1.Cluster, *clusterv1.Machine, *apiv1.Secret, *apiv1.Secret, error) {
	var err error
	machine, err := readMachineManifest(manifestParams, machineLoc)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cluster := &clusterv1.Cluster{}
	{
		bytes, err := ioutil.ReadFile(clusterLoc)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("cluster manifest %q: %v", clusterLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &cluster); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("cluster manifest %q: %v", clusterLoc, err)
		}
	}

	var awsCredentialsSecret *apiv1.Secret
	if awsCredentialSecretLoc != "" {
		awsCredentialsSecret = &apiv1.Secret{}
		bytes, err := ioutil.ReadFile(awsCredentialSecretLoc)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("aws credentials manifest %q: %v", awsCredentialSecretLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &awsCredentialsSecret); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("aws credentials manifest %q: %v", awsCredentialSecretLoc, err)
		}
	}

	var userDataSecret *apiv1.Secret
	if userDataLoc != "" {
		userDataSecret = &apiv1.Secret{}
		bytes, err := ioutil.ReadFile(userDataLoc)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("user data manifest %q: %v", userDataLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &userDataSecret); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("user data manifest %q: %v", userDataLoc, err)
		}
	}

	return cluster, machine, awsCredentialsSecret, userDataSecret, nil
}

func createActuator(machine *clusterv1.Machine, awsCredentials *apiv1.Secret, userData *apiv1.Secret) *machineactuator.Actuator {
	objList := []runtime.Object{}
	if awsCredentials != nil {
		objList = append(objList, awsCredentials)
	}
	if userData != nil {
		objList = append(objList, userData)
	}
	fakeKubeClient := kubernetesfake.NewSimpleClientset(objList...)
	fakeClient := fake.NewFakeClient(machine)

	params := machineactuator.ActuatorParams{
		Client:           fakeClient,
		KubeClient:       fakeKubeClient,
		AwsClientBuilder: awsclient.NewClient,
	}

	actuator, _ := machineactuator.NewActuator(params)
	return actuator
}

func checkFlags(cmd *cobra.Command) error {
	if cmd.Flag("cluster").Value.String() == "" {
		return fmt.Errorf("--%v/-%v flag is required", cmd.Flag("cluster").Name, cmd.Flag("cluster").Shorthand)
	}
	if cmd.Flag("machine").Value.String() == "" {
		return fmt.Errorf("--%v/-%v flag is required", cmd.Flag("machine").Name, cmd.Flag("machine").Shorthand)
	}
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %v\n", err)
		os.Exit(1)
	}
}

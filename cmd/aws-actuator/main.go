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
	"os/user"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/client"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/machine"

	"text/template"

	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/fake"
)

const (
	instanceIDAnnotation = "cluster-operator.openshift.io/aws-instance-id"
	ami                  = "ami-03f6257a"
	region               = "eu-west-1"
	size                 = "t1.micro"
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
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("aws-credentials").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := createActuator(machine, awsCredentials, userData, log.WithField("example", "create-machine"))
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
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("aws-credentials").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := createActuator(machine, awsCredentials, userData, log.WithField("example", "create-machine"))
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
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("aws-credentials").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := createActuator(machine, awsCredentials, userData, log.WithField("example", "create-machine"))
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

func readMachineManifest(manifestLoc string) (*clusterv1.Machine, error) {
	machine := &clusterv1.Machine{}
	bytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %v: %v", manifestLoc, err)
	}

	if err = yaml.Unmarshal(bytes, &machine); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal %v: %v", manifestLoc, err)
	}

	return machine, nil
}

func readClusterManifest(manifestLoc string) (*clusterv1.Cluster, error) {
	cluster := &clusterv1.Cluster{}
	bytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %v: %v", manifestLoc, err)
	}

	if err = yaml.Unmarshal(bytes, &cluster); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal %v: %v", manifestLoc, err)
	}

	return cluster, nil
}

func readSecretManifest(manifestLoc string) (*apiv1.Secret, error) {
	secret := &apiv1.Secret{}
	bytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %v: %v", manifestLoc, err)
	}
	if err = yaml.Unmarshal(bytes, &secret); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal %v: %v", manifestLoc, err)
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
yum install -y kubelet kubeadm --disableexcludes=kubernetes

cat <<EOF > /etc/default/kubelet
KUBELET_KUBEADM_EXTRA_ARGS=--cgroup-driver=systemd
EOF

kubeadm join {{ .MasterIp }}:8443 --token 2iqzqm.85bs0x6miyx1nm7l --discovery-token-unsafe-skip-ca-verification

HEREDOC

bash /root/user-data.sh > /root/user-data.logs
`

type userDataParams struct {
	MasterIp string
}

func generateWorkerUserData(masterIp string, workerUserDataSecret *apiv1.Secret) (*apiv1.Secret, error) {
	params := userDataParams{
		MasterIp: masterIp,
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

func bootstrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap kubernetes cluster with kubeadm",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestsDir := cmd.Flag("manifests").Value.String()
			if manifestsDir == "" {
				return fmt.Errorf("--manifests needs to be set")
			}

			log.Infof("Reading cluster manifest from %v", path.Join(manifestsDir, "cluster.yaml"))
			cluster, err := readClusterManifest(path.Join(manifestsDir, "cluster.yaml"))
			if err != nil {
				return err
			}

			log.Infof("Reading master machine manifest from %v", path.Join(manifestsDir, "master-machine.yaml"))
			masterMachine, err := readMachineManifest(path.Join(manifestsDir, "master-machine.yaml"))
			if err != nil {
				return err
			}

			log.Infof("Reading master user data manifest from %v", path.Join(manifestsDir, "master-userdata.yaml"))
			masterUserDataSecret, err := readSecretManifest(path.Join(manifestsDir, "master-userdata.yaml"))
			if err != nil {
				return err
			}

			log.Infof("Reading worker machine manifest from %v", path.Join(manifestsDir, "worker-machine.yaml"))
			workerMachine, err := readMachineManifest(path.Join(manifestsDir, "worker-machine.yaml"))
			if err != nil {
				return err
			}

			log.Infof("Reading worker user data manifest from %v", path.Join(manifestsDir, "worker-userdata.yaml"))
			workerUserDataSecret, err := readSecretManifest(path.Join(manifestsDir, "worker-userdata.yaml"))
			if err != nil {
				return err
			}

			var awsCredentialsSecret *apiv1.Secret
			if cmd.Flag("aws-credentials").Value.String() != "" {
				log.Infof("Reading aws credentials manifest from %v", cmd.Flag("aws-credentials").Value.String())
				awsCredentialsSecret, err = readSecretManifest(cmd.Flag("aws-credentials").Value.String())
				if err != nil {
					return err
				}
			}

			log.Infof("Creating master machine")
			actuator := createActuator(masterMachine, awsCredentialsSecret, masterUserDataSecret, log.WithField("bootstrap", "create-master-machine"))
			result, err := actuator.CreateMachine(cluster, masterMachine)
			if err != nil {
				return err
			}

			log.Infof("Master machine created with ipv4: %v, InstanceId: %v", *result.PrivateIpAddress, *result.InstanceId)

			log.Infof("Generating worker user data for master listening at %v", *result.PrivateIpAddress)
			workerUserDataSecret, err = generateWorkerUserData(*result.PrivateIpAddress, workerUserDataSecret)
			if err != nil {
				return fmt.Errorf("Unable to generate worker user data: %v\n", err)
			}

			log.Infof("Creating worker machine")
			actuator = createActuator(workerMachine, awsCredentialsSecret, workerUserDataSecret, log.WithField("bootstrap", "create-worker-machine"))
			result, err = actuator.CreateMachine(cluster, workerMachine)
			if err != nil {
				return err
			}

			log.Infof("Worker machine created with InstanceId: %v", *result.InstanceId)

			return nil
		},
	}

	cmd.PersistentFlags().StringP("manifests", "", "", "Directory with bootstrapping manifests")
	cUser, err := user.Current()
	if err != nil {
		cmd.PersistentFlags().StringP("machine-prefix", "p", "", "Directory with bootstrapping manifests")
	} else {
		cmd.PersistentFlags().StringP("machine-prefix", "p", cUser.Username, "Machine prefix, by default set to the current user")
	}

	return cmd
}

func init() {
	rootCmd.PersistentFlags().StringP("machine", "m", "", "Machine manifest")
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "Cluster manifest")
	rootCmd.PersistentFlags().StringP("aws-credentials", "a", "", "Secret manifest with aws credentials")
	rootCmd.PersistentFlags().StringP("userdata", "u", "", "User data manifest")

	rootCmd.AddCommand(createCommand())

	rootCmd.AddCommand(deleteCommand())

	rootCmd.AddCommand(existsCommand())

	rootCmd.AddCommand(bootstrapCommand())
}

func readClusterResources(clusterLoc, machineLoc, awsCredentialSecretLoc, userDataLoc string) (*clusterv1.Cluster, *clusterv1.Machine, *apiv1.Secret, *apiv1.Secret, error) {
	machine := &clusterv1.Machine{}
	{
		bytes, err := ioutil.ReadFile(machineLoc)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("machine manifest %q: %v", machineLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &machine); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("machine manifest %q: %v", machineLoc, err)
		}
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

func createActuator(machine *clusterv1.Machine, awsCredentials *apiv1.Secret, userData *apiv1.Secret, logger *log.Entry) *machineactuator.Actuator {
	objList := []runtime.Object{}
	if awsCredentials != nil {
		objList = append(objList, awsCredentials)
	}
	if userData != nil {
		objList = append(objList, userData)
	}
	fakeKubeClient := kubernetesfake.NewSimpleClientset(objList...)
	fakeClient := fake.NewSimpleClientset(machine)
	actuator, _ := machineactuator.NewActuator(fakeKubeClient, fakeClient, logger, awsclient.NewClient)
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
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %v\n", err)
		os.Exit(1)
	}
}

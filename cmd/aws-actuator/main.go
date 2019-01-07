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
	"context"
	goflag "flag"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	"github.com/openshift/cluster-api-actuator-pkg/pkg/manifests"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/actuators/machine"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	testutils "sigs.k8s.io/cluster-api-provider-aws/test/utils"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	region                   = "us-east-1"
	awsCredentialsSecretName = "aws-credentials-secret"
	pollInterval             = 5 * time.Second
	timeoutPoolAWSInterval   = 10 * time.Minute
)

func init() {
	// Add types to scheme
	clusterv1.AddToScheme(scheme.Scheme)

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

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
}

func usage() {
	fmt.Printf("Usage: %s\n\n", os.Args[0])
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
				return fmt.Errorf("unable to create read resources: %v", err)
			}

			actuator, err := createActuator(machine, awsCredentials, userData)
			if err != nil {
				return fmt.Errorf("unable to create actuator: %v", err)
			}
			result, err := actuator.CreateMachine(cluster, machine)
			if err != nil {
				return fmt.Errorf("unable to create machine: %v", err)
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

			if err != nil {
				return fmt.Errorf("unable to create read resources: %v", err)
			}

			actuator, err := createActuator(machine, awsCredentials, userData)
			if err != nil {
				return fmt.Errorf("unable to create actuator: %v", err)
			}
			if err = actuator.DeleteMachine(cluster, machine); err != nil {
				return fmt.Errorf("unable to delete machine: %v", err)
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
				return fmt.Errorf("unable to create read resources: %v", err)
			}

			actuator, err := createActuator(machine, awsCredentials, userData)
			if err != nil {
				return fmt.Errorf("unable to create actuator: %v", err)
			}
			exists, err := actuator.Exists(context.TODO(), cluster, machine)
			if err != nil {
				return fmt.Errorf("unable to check if machine exists: %v", err)
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

func bootstrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap kubernetes cluster with kubeadm",
		RunE: func(cmd *cobra.Command, args []string) error {
			machinePrefix := cmd.Flag("environment-id").Value.String()

			mastermachinepk := cmd.Flag("master-machine-private-key").Value.String()
			if mastermachinepk == "" {
				return fmt.Errorf("--master-machine-private-key needs to be set")
			}

			if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
				return fmt.Errorf("AWS_ACCESS_KEY_ID env needs to be set")
			}
			if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
				return fmt.Errorf("AWS_SECRET_ACCESS_KEY env needs to be set")
			}

			testNamespace := &apiv1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			}

			testCluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      machinePrefix,
					Namespace: testNamespace.Name,
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

			awsCredentialsSecret := testutils.GenerateAwsCredentialsSecretFromEnv(awsCredentialsSecretName, testNamespace.Name)

			// Create master machine and verify the master node is ready
			masterUserDataSecret, err := manifests.MasterMachineUserDataSecret(
				"masteruserdatasecret",
				testNamespace.Name,
				[]string{"\\$(curl -s http://169.254.169.254/latest/meta-data/public-hostname)", "\\$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)"},
			)
			if err != nil {
				return err
			}

			masterMachineProviderConfig, err := testutils.MasterMachineProviderConfig(awsCredentialsSecret.Name, masterUserDataSecret.Name, testCluster.Name)
			if err != nil {
				return err
			}

			masterMachine := manifests.MasterMachine(testCluster.Name, testCluster.Namespace, masterMachineProviderConfig)

			glog.Infof("Creating master machine")

			actuator, err := createActuator(masterMachine, awsCredentialsSecret, masterUserDataSecret)
			if err != nil {
				return fmt.Errorf("unable to create actuator: %v", err)
			}
			result, err := actuator.CreateMachine(testCluster, masterMachine)
			if err != nil {
				glog.Errorf("unable to create machine: %v", err)
				return fmt.Errorf("unable to create machine: %v", err)
			}

			glog.Infof("Master machine created with ipv4: %v, InstanceId: %v", *result.PrivateIpAddress, *result.InstanceId)

			masterMachinePrivateIP := ""
			err = wait.Poll(pollInterval, timeoutPoolAWSInterval, func() (bool, error) {
				glog.Info("Waiting for master machine PublicDNS")
				result, err := actuator.Describe(testCluster, masterMachine)
				if err != nil {
					glog.Info(err)
					return false, nil
				}

				glog.Infof("PublicDnsName: %v\n", *result.PublicDnsName)
				if *result.PublicDnsName == "" {
					return false, nil
				}

				masterMachinePrivateIP = *result.PrivateIpAddress
				return true, nil
			})
			if err != nil {
				glog.Errorf("Unable to get DNS name: %v", err)
				return err
			}

			f := framework.Framework{
				SSH: &framework.SSHConfig{
					Key:  mastermachinepk,
					User: "ec2-user",
				},
			}

			objList := []runtime.Object{awsCredentialsSecret}
			fakeClient := fake.NewFakeClient(objList...)
			awsClient, err := awsclient.NewClient(fakeClient, awsCredentialsSecret.Name, awsCredentialsSecret.Namespace, region)
			if err != nil {
				glog.Errorf("Unable to create aws client: %v", err)
				return err
			}

			acw := machineactuator.NewAwsClientWrapper(awsClient)
			glog.Infof("Collecting master kubeconfig")
			restConfig, err := f.GetMasterMachineRestConfig(masterMachine, acw)
			if err != nil {
				glog.Errorf("Unable to pull kubeconfig: %v", err)
				return err
			}

			clusterFramework, err := framework.NewFrameworkFromConfig(
				restConfig,
				&framework.SSHConfig{
					Key:  mastermachinepk,
					User: "ec2-user",
				},
			)
			if err != nil {
				return err
			}

			clusterFramework.ErrNotExpected = func(err error) {
				if err != nil {
					glog.Fatal(err)
				}
			}

			clusterFramework.By = func(msg string) {
				glog.Info(msg)
			}

			clusterFramework.MachineControllerImage = "openshift/origin-aws-machine-controllers:v4.0.0"
			clusterFramework.MachineManagerImage = "openshift/origin-aws-machine-controllers:v4.0.0"
			clusterFramework.NodelinkControllerImage = "registry.svc.ci.openshift.org/openshift/origin-v4.0-2019-01-03-031244@sha256:152c0a4ea7cda1731e45af87e33909421dcde7a8fcf4e973cd098a8bae892c50"

			glog.Info("Waiting for all nodes to come up")
			err = clusterFramework.WaitForNodesToGetReady(1)
			if err != nil {
				return err
			}

			glog.Infof("Creating %q namespace", testNamespace.Name)
			if _, err := clusterFramework.KubeClient.CoreV1().Namespaces().Create(testNamespace); err != nil {
				return err
			}

			clusterFramework.DeployClusterAPIStack(testNamespace.Name, "")
			clusterFramework.CreateClusterAndWait(testCluster)
			createSecretAndWait(clusterFramework, awsCredentialsSecret)

			workerUserDataSecret, err := manifests.WorkerMachineUserDataSecret("workeruserdatasecret", testNamespace.Name, masterMachinePrivateIP)
			if err != nil {
				return err
			}

			createSecretAndWait(clusterFramework, workerUserDataSecret)
			workerMachineSetProviderConfig, err := testutils.WorkerMachineSetProviderConfig(awsCredentialsSecret.Name, workerUserDataSecret.Name, testCluster.Name)
			if err != nil {
				return err
			}
			workerMachineSet := manifests.WorkerMachineSet(testCluster.Name, testCluster.Namespace, workerMachineSetProviderConfig)
			clusterFramework.CreateMachineSetAndWait(workerMachineSet, acw)

			return nil
		},
	}

	cmd.PersistentFlags().StringP("manifests", "", "", "Directory with bootstrapping manifests")
	cmd.PersistentFlags().StringP("master-machine-private-key", "", "", "Private key file of the master machine to pull kubeconfig")
	return cmd
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %v\n", err)
		os.Exit(1)
	}
}

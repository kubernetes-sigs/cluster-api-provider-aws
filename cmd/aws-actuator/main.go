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
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/client"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/machine"

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

func init() {
	rootCmd.PersistentFlags().StringP("machine", "m", "", "Machine manifest")
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "Cluster manifest")
	rootCmd.PersistentFlags().StringP("aws-credentials", "a", "", "Secret manifest with aws credentials")
	rootCmd.PersistentFlags().StringP("userdata", "u", "", "User data manifest")

	rootCmd.AddCommand(&cobra.Command{
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
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete INSTANCE-ID",
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
	})

	rootCmd.AddCommand(&cobra.Command{
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
	})
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

/*
Copyright 2021 The Kubernetes Authors.

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
package copy

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/flags"
	ec2service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
)

var (
	ownerID           string
	kubernetesVersion string
	opSystem          string
)

func addSourceRegion(c *cobra.Command) {
	c.Flags().String("source-region", "", "Set if wanting to copy an AMI from a different region")
}

func GetSourceRegion(c *cobra.Command) (string, error) {
	explicitRegion := c.Flags().Lookup("source-region").Value.String()
	if explicitRegion != "" {
		return explicitRegion, nil
	}
	region, err := flags.GetRegion(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not resolve source AWS region, define it with --source-region flag or as an environment variable.")
		return "", err
	}
	return region, nil
}

func addOsFlag(c *cobra.Command) {
	c.Flags().StringVar(&opSystem, "os", "", "Operating system of the AMI to be copied")
	if err := c.MarkFlagRequired("os"); err != nil {
		panic(errors.Wrap(err, "error marking required --os flag"))
	}
}

func addKubernetesVersionFlag(c *cobra.Command) {
	c.Flags().StringVar(&kubernetesVersion, "kubernetes-version", "", "Kubernetes version of the AMI to be copied")
	if err := c.MarkFlagRequired("kubernetes-version"); err != nil {
		panic(errors.Wrap(err, "error marking required --kubernetes-version flag"))
	}
}

func addOwnerIDFlag(c *cobra.Command) {
	c.Flags().StringVar(&ownerID, "owner-id", ec2service.DefaultMachineAMIOwnerID, "The source AWS owner ID, where the AMI will be copied from")
}

func addDryRunFlag(c *cobra.Command) {
	c.Flags().Bool("dry-run", false, "Check if AMI exists and can be copied")
}

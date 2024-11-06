/*
Copyright 2020 The Kubernetes Authors.

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

package iam

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
	"sigs.k8s.io/yaml"

	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cloudformation/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/configreader"
)

func printConfigCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "print-config",
		Short: "Print configuration",
		Long:  templates.LongDesc("Print configuration"),
		Example: templates.Examples(`
		# Print the default configuration.
		clusterawsadm bootstrap iam print-config

		# Apply defaults to a configuration file and print the result
		clusterawsadm bootstrap iam print-config --config bootstrap_config.yaml
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}

			conf := bootstrapv1.NewAWSIAMConfiguration()
			conf.Spec = *t.Spec

			b, err := yaml.Marshal(conf)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	}
	addConfigFlag(newCmd)

	return newCmd
}

func getBootstrapTemplate(cmd *cobra.Command) (*bootstrap.Template, error) {
	flag := cmd.Flags().Lookup("config")
	if flag == nil || flag.Value.String() == "" {
		t := bootstrap.NewTemplate()
		return &t, nil
	}

	val := flag.Value.String()

	iamConfiguration, err := configreader.LoadConfigFile(val)
	if err != nil {
		return nil, err
	}

	return &bootstrap.Template{
		Spec: &iamConfiguration.Spec,
	}, nil
}

func addConfigFlag(c *cobra.Command) {
	c.Flags().String("config", "", templates.LongDesc(`
		clusterawsadm will load a bootstrap configuration from this file. The path may be
		absolute or relative; relative paths start at the current working directory.

		The configuration file is a Kubernetes YAML using the
		bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1/AWSIAMConfiguration
		kind.

		Documentation for this kind can be found at:
		https://pkg.go.dev/sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1beta1

		To see the default configuration, run 'clusterawsadm bootstrap iam print-config'.
	`))
}

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

// Package credentials provides a way to encode credentials for use with Kubernetes Cluster API Provider AWS.
package credentials

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/flags"
	creds "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/credentials"
)

const (
	rawSharedConfig    = "rawSharedConfig"
	base64SharedConfig = "base64SharedConfig"
	// BackupAWSRegion is a region used purely for resolving credentials if it is
	// not available in the default lookup chain.
	backupAWSRegion = "us-east-1"

	// CredentialHelp provides an explanation as to how credentials are resolved by
	// clusterawsadm.
	//nolint:gosec
	CredentialHelp = `
	The utility will attempt to find credentials in the following order:

	  1. Check for the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables.
	  2. Read the default credentials from the shared configuration files ~/.aws/credentials
	     or the default profile in ~/.aws/config.
	  3. Check for the presence of an EC2 IAM instance profile if it's running on AWS.
	  4. Check for ECS credentials.

	IAM role assumption can be performed by using any valid configuration for the AWS CLI at:
	https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html. For role assumption
	to be used, a region is required for the utility to use the AWS Security Token Service (STS). The utility
	resolves the region in the following order:

	  1. Check for the --region flag.
	  2. Check for the AWS_REGION environment variable.
	  3. Check for the DEFAULT_AWS_REGION environment variable.
	  4. Check that a region is specified in the shared configuration file.
	`
	// EncodingHelp provides an explanation for how clusterawsadm will generate ini-files
	// for the resolved credentials.
	EncodingHelp = `
	The utility will then generate an ini-file with a default profile corresponding to
	the resolved credentials.

	If a region cannot be found, for the purposes of using AWS Security Token Service, this
	utility will fall back to us-east-1. This does not affect the region in which clusters
	will be created.

	In the case of an instance profile or role assumption, note that encoded credentials are
	time-limited.
	`

	examples = `
		# Encode credentials from the environment for use with clusterctl
		export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
		clusterctl init --infrastructure aws
	`
)

var errInvalidOutputFlag = errors.New("invalid output flag. Expected rawSharedConfig or base64SharedConfig")

// RootCmd is the root of the `alpha bootstrap command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "credentials",
		Short: `Encode credentials to use with Kubernetes Cluster API Provider AWS`,
		Long: templates.LongDesc(`
			Encode credentials to use with Kubernetes Cluster API Provider AWS.
			` + CredentialHelp + EncodingHelp),
		Example: templates.Examples(examples),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	newCmd.AddCommand(generateAWSDefaultProfileWithChain())

	return newCmd
}

func getOutputFlag(cmd *cobra.Command) (string, error) {
	val := cmd.Flags().Lookup("output").Value.String()
	switch val {
	case rawSharedConfig, base64SharedConfig:
		return val, nil
	default:
		return "", errInvalidOutputFlag
	}
}

func generateAWSDefaultProfileWithChain() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "encode-as-profile",
		Short: "Generate an AWS profile from the current environment",
		Long: templates.LongDesc(`
		Generate an AWS profile from the current environment for the ephemeral bootstrap cluster.
		` + CredentialHelp + EncodingHelp),
		Example: templates.Examples(examples),
		RunE: func(c *cobra.Command, args []string) error {
			flags.CredentialWarning(c)

			output, err := getOutputFlag(c)
			if err != nil {
				return err
			}

			region, err := flags.GetRegion(c)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not resolve AWS region, defaulting to %s.\n\n", backupAWSRegion)
				region = backupAWSRegion
			}

			awsCreds, err := creds.NewAWSCredentialFromDefaultChain(region)
			if err != nil {
				return flags.ResolveAWSError(err)
			}

			out := ""
			switch output {
			case rawSharedConfig:
				out, err = awsCreds.RenderAWSDefaultProfile()
			case base64SharedConfig:
				out, err = awsCreds.RenderBase64EncodedAWSDefaultProfile()
			}

			if err != nil {
				return err
			}

			fmt.Println(out)

			return nil
		},
	}

	newCmd.Flags().String("output", string(base64SharedConfig), "Output for credential configuration (rawSharedConfig, base64SharedConfig)")
	flags.AddRegionFlag(newCmd)
	return newCmd
}

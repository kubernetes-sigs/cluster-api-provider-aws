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

package flags

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
)

func ResolveAWSError(err error) error {
	code, _ := awserrors.Code(err)
	if code == awserrors.NoCredentialProviders {
		return errors.New("could not resolve default credentials. Please see https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html for how to provide credentials")
	}
	return err
}

func AddRegionFlag(c *cobra.Command) {
	c.Flags().String("region", "", "The AWS region in which to provision")
}

func GetRegion(c *cobra.Command) (string, error) {
	explicitRegion := c.Flags().Lookup("region").Value.String()
	return credentials.ResolveRegion(explicitRegion)
}

func GetRegionWithError(c *cobra.Command) (string, error) {
	region, err := GetRegion(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not resolve AWS region, define it with --region flag or as an environment variable.")
		return "", err
	}
	return region, nil
}

func MarkAlphaDeprecated(c *cobra.Command) {
	c.Deprecated = "and will be removed in 0.6.0"
}

func CredentialWarning(c *cobra.Command) {
	fmt.Fprintf(os.Stderr, "\nWARNING: `%s` should only be used for bootstrapping.\n\n", c.Name())
}

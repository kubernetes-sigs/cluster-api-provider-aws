/*
Copyright 2024 The Kubernetes Authors.

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

package credentials

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

// UseEKSPodIdentityCmd is a CLI command that will enable using EKS pod identity for CAPA.
func UseEKSPodIdentityCmd() *cobra.Command {
	clusterName := ""
	region := ""
	namespace := ""
	serviceAccount := ""
	roleName := ""

	newCmd := &cobra.Command{
		Use:   "use-pod-identity",
		Short: "Enable EKS pod identity with CAPA",
		Long: templates.LongDesc(`
			Updates CAPA running in an EKS cluster to use EKS pod identity
		`),
		Example: templates.Examples(`
		clusterawsadm controller use-pod-identity --cluster-name cluster1
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return usePodIdentity(region, clusterName, namespace, serviceAccount, roleName)
		},
	}

	newCmd.Flags().StringVarP(&region, "region", "r", "", "The AWS region containing the EKS cluster")
	newCmd.Flags().StringVarP(&clusterName, "cluster-name", "n", "", "The name of the EKS management cluster")
	newCmd.Flags().StringVar(&namespace, "namespace", "capa-system", "The namespace of CAPA controller")
	newCmd.Flags().StringVar(&serviceAccount, "service-account", "capa-controller-manager", "The service account for the CAPA controller")
	newCmd.Flags().StringVar(&roleName, "role-name", "controllers.cluster-api-provider-aws.sigs.k8s.io", "The name of the CAPA controller role. If you have used a prefix or suffix this will need to be changed.")

	newCmd.MarkFlagRequired("cluster-name")

	return newCmd
}

func usePodIdentity(region, clusterName, namespace, serviceAccount, roleName string) error {
	cfg := aws.Config{}
	if region != "" {
		cfg.Region = aws.String(region)
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            cfg,
	})
	if err != nil {
		return fmt.Errorf("failed creating aws session: %w", err)
	}

	roleArn, err := getRoleArn(sess, roleName)
	if err != nil {
		return err
	}

	eksClient := eks.New(sess)

	listInput := &eks.ListPodIdentityAssociationsInput{
		ClusterName: aws.String(clusterName),
		Namespace:   aws.String(namespace),
	}

	listOutput, err := eksClient.ListPodIdentityAssociations(listInput)
	if err != nil {
		return fmt.Errorf("listing existing pod identity associations for cluster %s in namespace %s: %w", clusterName, namespace, err)
	}

	for _, association := range listOutput.Associations {
		if *association.ServiceAccount == serviceAccount {
			needsUpdate, err := podIdentityNeedsUpdate(eksClient, association, roleName)
			if err != nil {
				return err
			}
			if !needsUpdate {
				fmt.Printf("EKS pod association for service account %s already exists, no action taken\n", serviceAccount)
			}

			return updatePodIdentity(eksClient, association, roleName)
		}
	}

	fmt.Printf("Creating pod association for service account %s.....\n", serviceAccount)

	createInput := &eks.CreatePodIdentityAssociationInput{
		ClusterName:    &clusterName,
		Namespace:      &namespace,
		RoleArn:        &roleArn,
		ServiceAccount: &serviceAccount,
	}

	output, err := eksClient.CreatePodIdentityAssociation(createInput)
	if err != nil {
		return fmt.Errorf("failed to create pod identity association: %w", err)
	}

	fmt.Printf("Created pod identity association (%s)\n", *output.Association.AssociationId)

	return nil
}

func podIdentityNeedsUpdate(client *eks.EKS, association *eks.PodIdentityAssociationSummary, roleArn string) (bool, error) {
	input := &eks.DescribePodIdentityAssociationInput{
		AssociationId: association.AssociationId,
		ClusterName:   association.ClusterName,
	}

	output, err := client.DescribePodIdentityAssociation(input)
	if err != nil {
		return false, fmt.Errorf("failed describing pod identity association: %w", err)
	}

	return *output.Association.RoleArn != roleArn, nil
}

func updatePodIdentity(client *eks.EKS, association *eks.PodIdentityAssociationSummary, roleArn string) error {
	input := &eks.UpdatePodIdentityAssociationInput{
		AssociationId: association.AssociationId,
		ClusterName:   association.ClusterName,
		RoleArn:       &roleArn,
	}

	_, err := client.UpdatePodIdentityAssociation(input)
	if err != nil {
		return fmt.Errorf("failed updating pod identity association: %w", err)
	}

	fmt.Printf("Updated pod identity to use role %s\n", roleArn)

	return nil
}

func getRoleArn(sess *session.Session, roleName string) (string, error) {
	client := iam.New(sess)

	input := &iam.GetRoleInput{
		RoleName: &roleName,
	}

	output, err := client.GetRole(input)
	if err != nil {
		return "", fmt.Errorf("failed looking up role %s: %w", roleName, err)
	}

	return *output.Role.Arn, nil
}

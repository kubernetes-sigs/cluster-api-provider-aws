package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/bootstrap/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/flags"
	iamservice "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/iamservice"
)

func createResources() *cobra.Command {
	newCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "[alpha] Create CAPA managed AWS IAM entities",
		Args:    cobra.NoArgs,
		Long: templates.LongDesc(`
		Create CAPA managed AWS IAM entities like roles, instances profiles and policies
		from the bootstrap configuration file (uses default configuration if not provided)
		` + credentials.CredentialHelp),
		Example: templates.Examples(`
		clusterawsadm bootstrap iam create
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			t, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}
			if err := resolveTemplateRegion(t, cmd); err != nil {
				fmt.Println("AWS_REGION env not set and --region flag not provided, default configuration will be used")
			}
			fmt.Printf("Attempting to create CAPA managed AWS IAM entities %s\n", t.Spec.StackName)
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(t.Spec.Region))
			if err != nil {
				return err
			}
			iamSvc := iam.NewFromConfig(cfg)
			svc := iamservice.New(iamSvc)
			err = svc.CreateResources(ctx, *t.RenderCloudFormation(), t.Spec.StackTags)
			if err != nil {
				return err
			}
			return nil
		},
	}
	addConfigFlag(newCmd)
	flags.AddRegionFlag(newCmd)
	return newCmd
}

func deleteResources() *cobra.Command {
	newCmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "[alpha] Delete CAPA managed AWS IAM entities",
		Args:    cobra.NoArgs,
		Long: templates.LongDesc(`
		Delete CAPA managed AWS IAM entities like roles, instances profiles and policies
		from the bootstrap configuration file (uses default configuration if not provided)
		` + credentials.CredentialHelp),
		Example: templates.Examples(`
		clusterawsadm bootstrap iam delete
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			t, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}
			if err := resolveTemplateRegion(t, cmd); err != nil {
				fmt.Println("AWS_REGION env not set and --region flag not provided, default configuration will be used")
			}
			fmt.Printf("Attempting to delete CAPA managed AWS IAM entities %s\n", t.Spec.StackName)
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(t.Spec.Region))
			if err != nil {
				return err
			}
			iamsvc := iam.NewFromConfig(cfg)
			svc := iamservice.New(iamsvc)
			err = svc.DeleteResources(ctx, *t.RenderCloudFormation(), t.Spec.StackTags)
			if err != nil {
				return err
			}
			return nil
		},
	}
	addConfigFlag(newCmd)
	flags.AddRegionFlag(newCmd)
	return newCmd
}

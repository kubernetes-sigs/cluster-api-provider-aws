package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/bootstrap/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/flags"
	iamservice "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/iamservice"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

func createServices() *cobra.Command {
	newCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create AWS IAM resources",
		Args:    cobra.NoArgs,
		Long: cmd.LongDesc(`
		` + credentials.CredentialHelp),
		Example: cmd.Examples(`
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := getBootstrapTemplate(cmd)
			if err != nil {
				return err
			}
			if err := resolveTemplateRegion(t, cmd); err != nil {
				fmt.Println("AWS_REGION env not set and --region flag not provided, default configuration will be used")
			}
			fmt.Printf("Attempting to create AWS Resources %s\n", t.Spec.StackName)
			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Region: aws.String(t.Spec.Region)},
			})
			if err != nil {
				return err
			}
			iamSvc := iam.New(sess)
			svc := iamservice.New(iamSvc)
			err = svc.CreateServices(*t.RenderCloudFormation(), t.Spec.StackTags)
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

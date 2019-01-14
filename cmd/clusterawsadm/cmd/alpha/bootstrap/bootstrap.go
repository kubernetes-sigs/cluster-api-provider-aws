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

package bootstrap

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	awssts "github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/sts"
)

// KubernetesAWSSecret is the template to generate an encoded version of the
// users' AWS credentials
// nolint
const KubernetesAWSSecret = `apiVersion: v1
kind: Secret
metadata:
  name: credentials.cluster-api-provider-aws.sigs.k8s.io
type: Opaque
data:
  credentials: {{ .CredentialsFile }}
`

// AWSCredentialsTemplate generates an AWS credentials file that can
// be loaded by the various SDKs.
const AWSCredentialsTemplate = `[default]
aws_access_key_id = {{ .AccessKeyID }}
aws_secret_access_key = {{ .SecretAccessKey }}
region = {{ .Region }}
{{if .SessionToken }}
aws_session_token = {{ .SessionToken }}
{{end}}
`

// RootCmd is the root of the `alpha bootstrap command`
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap cloudformation",
		Long:  `Create and apply bootstrap AWS CloudFormation template to create IAM permissions for the Cluster API`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	newCmd.AddCommand(generateCmd())
	newCmd.AddCommand(createStackCmd())
	newCmd.AddCommand(generateIAMPolicyDocJSON())
	newCmd.AddCommand(encodeAWSSecret())
	newCmd.AddCommand(generateAWSDefaultProfile())
	return newCmd
}

func generateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-cloudformation [AWS Account ID]",
		Short: "Generate bootstrap AWS CloudFormation template",
		Long: `Generate bootstrap AWS CloudFormation template with initial IAM policies.
You must enter an AWS account ID to generate the CloudFormation template.

Instructions for obtaining the AWS account ID can be found on https://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Printf("Error: requires AWS Account ID as an argument\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(200)
			}
			if !sts.ValidateAccountID(args[0]) {
				fmt.Printf("Error: provided AWS Account ID is invalid\n\n")
				cmd.Help()
				os.Exit(201)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			template := cloudformation.BootstrapTemplate(args[0])
			j, err := template.YAML()
			if err != nil {
				return err
			}

			fmt.Print(string(j))
			return nil
		},
	}
	return newCmd
}

func createStackCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "create-stack",
		Short: "Create a new AWS CloudFormation stack using the bootstrap template",
		Long:  "Create a new AWS CloudFormation stack using the bootstrap template",
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := "cluster-api-provider-aws-sigs-k8s-io"
			fmt.Printf("Attempting to create CloudFormation stack %s", stackName)
			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			})
			if err != nil {
				return err
			}

			stsSvc := sts.NewService(awssts.New(sess))
			accountID, stsErr := stsSvc.AccountID()
			if stsErr != nil {
				return stsErr
			}

			cfnSvc := cloudformation.NewService(cfn.New(sess))
			err = cfnSvc.ReconcileBootstrapStack(stackName, accountID)
			if err != nil {
				return err
			}

			return cfnSvc.ShowStackResources(stackName)
		},
	}

	return newCmd
}

func generateIAMPolicyDocJSON() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-iam-policy-docs [AWS Account ID] [Directory for JSON]",
		Short: "Generate PolicyDocument JSON for all ManagedIAMPolicies",
		Long:  `Generate PolicyDocument JSON for all ManagedIAMPolicies`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				fmt.Printf("Error: requires, as arguments, an AWS Account ID and a directory for the exported JSON\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(300)
			}
			accountID := args[0]
			policyDocDir := args[1]

			var err error
			if !sts.ValidateAccountID(accountID) {
				fmt.Printf("Error: provided AWS Account ID is invalid\n\n")
				cmd.Help()
				os.Exit(301)
			}

			if _, err = os.Stat(policyDocDir); os.IsNotExist(err) {
				err = os.Mkdir(policyDocDir, 0755)
				if err != nil {
					fmt.Printf("Error: failed to make directory %q, %v", policyDocDir, err)
					cmd.Help()
					os.Exit(302)
				}
			}
			if err != nil {
				fmt.Printf("Error: failed to stat directory %q, %v", policyDocDir, err)
				cmd.Help()
				os.Exit(303)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID := args[0]
			policyDocDir := args[1]
			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			})
			if err != nil {
				return fmt.Errorf("Error: failed to create a session: %v", err)
			}

			cfnSvc := cloudformation.NewService(cfn.New(sess))
			err = cfnSvc.GenerateManagedIAMPolicyDocuments(policyDocDir, accountID)

			if err != nil {
				return fmt.Errorf("Error: failed to generate PolicyDocument for all ManagedIAMPolicies: %v", err)
			}

			fmt.Printf("PolicyDocument for all ManagedIAMPolicies successfully generated in JSON at %q\n", policyDocDir)
			return nil
		},
	}
	return newCmd
}

func encodeAWSSecret() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "encode-aws-credentials",
		Short: "Encode AWS credentials as a base64 encoded Kubernetes secret",
		Long:  "Encode AWS credentials as a base64 encoded Kubernetes secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := getCredentialsFromEnvironment()

			if err != nil {
				return err
			}

			err = generateAWSKubernetesSecret(*creds)

			if err != nil {
				return err
			}

			return nil
		},
	}

	return newCmd
}

func generateAWSDefaultProfile() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-aws-default-profile",
		Short: "Generate an AWS profile from the current environment",
		Long:  "Generate an AWS profile from the current environment to be saved into minikube",
		RunE: func(cmd *cobra.Command, args []string) error {

			creds, err := getCredentialsFromEnvironment()

			if err != nil {
				return err
			}

			profile, err := renderAWSDefaultProfile(*creds)

			if err != nil {
				return err
			}

			fmt.Println(profile.String())

			return nil
		},
	}

	return newCmd
}

func getCredentialsFromEnvironment() (*awsCredential, error) {
	creds := awsCredential{}

	region, err := getEnv("AWS_REGION")
	if err != nil {
		return nil, err
	}
	creds.Region = region

	accessKeyID, err := getEnv("AWS_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}
	creds.AccessKeyID = accessKeyID

	secretAccessKey, err := getEnv("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}
	creds.SecretAccessKey = secretAccessKey

	sessionToken, err := getEnv("AWS_SESSION_TOKEN")
	if err != nil {
		creds.SessionToken = ""
	} else {
		creds.SessionToken = sessionToken
	}

	return &creds, nil
}

type awsCredential struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

type awsCredentialsFile struct {
	CredentialsFile string
}

func getEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("Environment variable %q not found", key)
	}
	return val, nil
}

func renderAWSDefaultProfile(creds awsCredential) (*bytes.Buffer, error) {
	tmpl, err := template.New("AWS Credentials").Parse(AWSCredentialsTemplate)
	if err != nil {
		return nil, err
	}

	var credsFileStr bytes.Buffer
	err = tmpl.Execute(&credsFileStr, creds)
	if err != nil {
		return nil, err
	}

	return &credsFileStr, nil
}

func generateAWSKubernetesSecret(creds awsCredential) error {

	profile, err := renderAWSDefaultProfile(creds)

	if err != nil {
		return err
	}

	encCreds := base64.StdEncoding.EncodeToString(profile.Bytes())

	credsFile := awsCredentialsFile{
		CredentialsFile: encCreds,
	}

	secretTmpl, err := template.New("AWS Credentials Secret").Parse(KubernetesAWSSecret)
	if err != nil {
		return err
	}
	var out bytes.Buffer

	err = secretTmpl.Execute(&out, credsFile)

	if err != nil {
		return err
	}

	fmt.Println(out.String())

	return nil
}

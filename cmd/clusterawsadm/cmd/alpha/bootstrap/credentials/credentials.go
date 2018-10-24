// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package credentials

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
)

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

// Cmd is the root of the `alpha bootstrap aws command`
func Cmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "credentials",
		Short: "credentials",
		Long:  `credentials commands`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	newCmd.AddCommand(generateAWSDefaultProfile())
	newCmd.AddCommand(generateCredentials())
	return newCmd
}

func generateCredentials() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-credentials [user name]",
		Short: "Generate an AWS profile for a given IAM user",
		Long:  "Generate an AWS profile for a given IAM user to be saved as a kubernetes secret",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Printf("Error: requires user name\n\n")
				cmd.Help()
				os.Exit(200)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			sess, err := session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			})
			if err != nil {
				return err
			}

			svc := iam.New(sess)
			out, err := svc.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: aws.String(username)})
			if err != nil {
				return err
			}

			creds := awsCredential{
				Region:          aws.StringValue(sess.Config.Region),
				AccessKeyID:     aws.StringValue(out.AccessKey.AccessKeyId),
				SecretAccessKey: aws.StringValue(out.AccessKey.SecretAccessKey),
			}

			profile, err := renderAWSDefaultProfile(creds)
			if err != nil {
				return err
			}

			fmt.Println(profile.String())

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

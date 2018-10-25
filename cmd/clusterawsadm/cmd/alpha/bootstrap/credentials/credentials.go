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
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
)

// KubernetesAWSSecret is the template to generate an encoded version of the
// users' AWS credentials
const KubernetesAWSSecret = `apiVersion: v1
kind: Secret
metadata:
  name: credentials.cluster-api-provider-aws.sigs.k8s.io
  namespace: aws-provider-system
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
	newCmd.AddCommand(generateCredentials())
	return newCmd
}

func generateCredentials() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "generate-credentials [user name]",
		Short: "Generate a Kubernetes secret containing an AWS profile for a given IAM user",
		Long:  "Generate a Kubernetes secret containing an AWS profile for a given IAM user to be saved as a kubernetes secret",
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

			return generateAWSKubernetesSecret(creds)
		},
	}

	return newCmd
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
		CredentialsFile: string(encCreds),
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

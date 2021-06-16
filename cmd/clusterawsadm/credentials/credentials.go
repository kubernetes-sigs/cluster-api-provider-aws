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

package credentials

import (
	"bytes"
	"encoding/base64"
	"errors"
	"text/template"

	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

// AWSDefaultRegion is the default AWS region.
const AWSDefaultRegion = "us-east-1"

// ErrNoAWSRegionConfigured is an error singleton for when no AWS region is configured.
var ErrNoAWSRegionConfigured = errors.New("no AWS region configured. Use --region or set AWS_REGION or DEFAULT_AWS_REGION environment variable")

// AWSCredentials defines the specs for AWS credentials.
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

// NewAWSCredentialFromDefaultChain will create a new credential provider chain from the
// default chain.
func NewAWSCredentialFromDefaultChain(region string) (*AWSCredentials, error) {
	creds := AWSCredentials{}
	conf := aws.NewConfig()
	conf.CredentialsChainVerboseErrors = aws.Bool(true)
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *conf,
	})
	if err != nil {
		return nil, err
	}
	chainCreds, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	creds.Region = region
	creds.AccessKeyID = chainCreds.AccessKeyID
	creds.SecretAccessKey = chainCreds.SecretAccessKey
	creds.SessionToken = chainCreds.SessionToken

	return &creds, nil
}

// ResolveRegion will attempt to resolve an AWS region based on the customer's configuration.
func ResolveRegion(explicitRegion string) (string, error) {
	if explicitRegion != "" {
		return explicitRegion, nil
	}
	region, err := util.GetEnv("AWS_REGION")
	if err == nil {
		return region, nil
	}
	region, err = util.GetEnv("DEFAULT_AWS_REGION")
	if err == nil {
		return region, nil
	}
	return "", ErrNoAWSRegionConfigured
}

// RenderAWSDefaultProfile will render the AWS default profile.
func (c AWSCredentials) RenderAWSDefaultProfile() (string, error) {
	tmpl, err := template.New("AWS Credentials").Parse(AWSCredentialsTemplate)
	if err != nil {
		return "", err
	}

	var credsFileStr bytes.Buffer
	err = tmpl.Execute(&credsFileStr, c)
	if err != nil {
		return "", err
	}

	return credsFileStr.String(), nil
}

// RenderBase64EncodedAWSDefaultProfile will render the AWS default profile, encoded in base 64.
func (c AWSCredentials) RenderBase64EncodedAWSDefaultProfile() (string, error) {
	profile, err := c.RenderAWSDefaultProfile()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(profile)), nil
}

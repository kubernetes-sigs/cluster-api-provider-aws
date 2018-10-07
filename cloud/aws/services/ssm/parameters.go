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

package ssm

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/awserrors"
)

const (
	// Prefix is the parameter prefix used by Cluster API Provider AWS
	Prefix = "/sigs.k8s.io/cluster-api-provider-aws"
)

func (s *Service) ReconcileParameter(cluster string, path string, value string) error {

	err := s.putParameter(cluster, path, value, false)

	if code, ok := awserrors.Code(err); ok && code == "ParameterAlreadyExists" {
		return nil
	}

	return err
}

func (s *Service) putParameter(cluster string, path string, value string, overwrite bool) error {

	input := &ssm.PutParameterInput{
		Name:      aws.String(ResolvePath(cluster, path)),
		Value:     aws.String(value),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(overwrite),
	}

	_, err := s.SSM.PutParameter(input)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetParameter(cluster string, path string) (string, error) {

	input := &ssm.GetParameterInput{
		Name: aws.String(ResolvePath(cluster, path)),
	}

	out, err := s.SSM.GetParameter(input)

	if err != nil {
		return "", err
	}

	return aws.StringValue(out.Parameter.Value), nil
}

func (s *Service) DeleteParameter(cluster string, path string) error {

	input := &ssm.DeleteParameterInput{
		Name: aws.String(ResolvePath(cluster, path)),
	}

	_, err := s.SSM.DeleteParameter(input)

	if code, _ := awserrors.Code(err); code == "ParameterNotFound" {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

// ResolvePath provides a qualified SSM path
func ResolvePath(cluster string, path string) string {
	return fmt.Sprintf("%s/%s/%s", Prefix, path, cluster)
}

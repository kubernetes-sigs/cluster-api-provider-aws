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

package ec2

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"go.opencensus.io/trace"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/instrumentation"
)

// GetConsole returns the latest console output of an instance
func (s *Service) GetConsole(ctx context.Context, instanceID string) (string, error) {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "ec2", "GetConsole"),
	)
	defer span.End()

	input := &ec2.GetConsoleOutputInput{
		InstanceId: aws.String(instanceID),
		Latest:     aws.Bool(true),
	}

	out, err := s.EC2.GetConsoleOutput(input)

	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(aws.StringValue(out.Output))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

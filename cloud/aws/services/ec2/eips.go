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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

func (s *Service) allocateAddress() (string, error) {
	out, err := s.EC2.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to create Elastic IP address")
	}

	return *out.AllocationId, nil
}

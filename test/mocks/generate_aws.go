/*
Copyright 2022 The Kubernetes Authors.

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

//go:generate ../../hack/tools/bin/mockgen -destination aws_elbv2_mock.go -package mocks github.com/aws/aws-sdk-go/service/elbv2/elbv2iface ELBV2API
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt aws_elbv2_mock.go > _aws_elbv2_mock.go && mv _aws_elbv2_mock.go aws_elbv2_mock.go"

//go:generate ../../hack/tools/bin/mockgen -destination aws_elb_mock.go -package mocks github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt aws_elb_mock.go > _aws_elb_mock.go && mv _aws_elb_mock.go aws_elb_mock.go"

//go:generate ../../hack/tools/bin/mockgen -destination aws_rgtagging_mock.go -package mocks github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface ResourceGroupsTaggingAPIAPI
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt aws_rgtagging_mock.go > _aws_rgtagging_mock.go && mv _aws_rgtagging_mock.go aws_rgtagging_mock.go"

//go:generate ../../hack/tools/bin/mockgen -destination aws_ec2api_mock.go -package mocks github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt aws_ec2api_mock.go > _aws_ec2api_mock.go && mv _aws_ec2api_mock.go aws_ec2api_mock.go"

package mocks

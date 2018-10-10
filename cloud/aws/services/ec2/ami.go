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

// defaultAMILookup returns the default AMI based on region
func (s *Service) defaultAMILookup(region string) string {
	//TODO(chuckha) Replace this function with a method to filter public images.
	switch region {
	case "ap-northeast-1":
		return "ami-01b9b95604fddade1"
	case "ap-northeast-2":
		return "ami-0429d18e76a0b3705"
	case "ap-south-1":
		return "ami-09930fca077b07e18"
	case "ap-southeast-1":
		return "ami-0b2e5665eda719758"
	case "ap-southeast-2":
		return "ami-01c2df5ce9bc2f573"
	case "ca-central-1":
		return "ami-0f5fbb71a0c65e000"
	case "eu-central-1":
		return "ami-0fad7824ed21125b1"
	case "eu-west-1":
		return "ami-0da760e590e7de0e8"
	case "eu-west-2":
		return "ami-04137690b33cb5a8e"
	case "eu-west-3":
		return "ami-028272d8c8e9ff369"
	case "sa-east-1":
		return "ami-041a08e5511d25535"
	case "us-east-1":
		return "ami-0de61b6929e8f091c"
	case "us-east-2":
		return "ami-0a2463ac1e1f46b95"
	case "us-west-1":
		return "ami-05dc1567db5bd869a"
	case "us-west-2":
		return "ami-0f33a1d90f189e0a1"
	default:
		return "unknown region"
	}
}

func (s *Service) defaultBastionAMILookup(region string) string {
	switch region {
	case "ap-northeast-1":
		return "ami-d39a02b5"
	case "ap-northeast-2":
		return "ami-67973709"
	case "ap-south-1":
		return "ami-5d055232"
	case "ap-southeast-1":
		return "ami-325d2e4e"
	case "ap-southeast-2":
		return "ami-37df2255"
	case "ca-central-1":
		return "ami-f0870294"
	case "eu-central-1":
		return "ami-af79ebc0"
	case "eu-west-1":
		return "ami-4d46d534"
	case "eu-west-2":
		return "ami-d7aab2b3"
	case "eu-west-3":
		return "ami-5e0eb923"
	case "sa-east-1":
		return "ami-1157157d"
	case "us-east-1":
		return "ami-41e0b93b"
	case "us-east-2":
		return "ami-2581aa40"
	case "us-west-1":
		return "ami-79aeae19"
	case "us-west-2":
		return "ami-1ee65166"
	default:
		return "unknown region"
	}
}

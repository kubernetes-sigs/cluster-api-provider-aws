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

package ec2

// defaultAMILookup returns the default AMI based on region
func (s *Service) defaultAMILookup(region string) string {
	//TODO(chuckha) Replace this function with a method to filter public images.
	switch region {
	case "ap-northeast-1":
		return "ami-06525773f0332c0ca"
	case "ap-northeast-2":
		return "ami-0914a816a80ab8779"
	case "ap-south-1":
		return "ami-04b17f6875b4d9c29"
	case "ap-southeast-1":
		return "ami-0ccaa193408259c82"
	case "ap-southeast-2":
		return "ami-0385d030f1b7df12c"
	case "ca-central-1":
		return "ami-031ca32942ffbe25f"
	case "eu-central-1":
		return "ami-078f9fe8bdaa81aa8"
	case "eu-west-1":
		return "ami-09547cd6f9856a79b"
	case "eu-west-2":
		return "ami-08efe1a8b5f18f381"
	case "eu-west-3":
		return "ami-042c578ede4ff3f33"
	case "sa-east-1":
		return "ami-0dd7e2894db509204"
	case "us-east-1":
		return "ami-026e9f3a713727945"
	case "us-east-2":
		return "ami-0ec6d3241fb7776fe"
	case "us-west-1":
		return "ami-06ec1c533176de131"
	case "us-west-2":
		return "ami-0cfa2d1fa5cc93615"
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

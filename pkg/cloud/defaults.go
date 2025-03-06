/*
Copyright 2024 The Kubernetes Authors.

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

package cloud

const (
	// TemporaryResourceID is the name used temporarily when creating AWS resources.
	TemporaryResourceID = "temporary-resource-id"
	// AnyIPv4CidrBlock is the CIDR block to match all IPv4 addresses.
	AnyIPv4CidrBlock = "0.0.0.0/0"
	// AnyIPv6CidrBlock is the CIDR block to match all IPv6 addresses.
	AnyIPv6CidrBlock = "::/0"
)

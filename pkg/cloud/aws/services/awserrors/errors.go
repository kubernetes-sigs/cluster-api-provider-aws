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

package awserrors

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// Code returns the AWS error code as a string
func Code(err error) (string, bool) {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Code(), true
	}
	return "", false
}

// Message returns the AWS error message as a string
func Message(err error) string {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Message()
	}
	return ""
}

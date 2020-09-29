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

package ssm

import (
	"bytes"
	"net/mail"
	"testing"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

func TestUserData(t *testing.T) {
	service := Service{}
	endpoints := []scope.ServiceEndpoint{}
	doc, _ := service.UserData("secretARN", 1, "eu-west-1", endpoints)

	_, err := mail.ReadMessage(bytes.NewBuffer(doc))
	if err != nil {
		t.Fatalf("Cannot parse MIME doc: %+v\n%s", err, string(doc))
	}
}

func TestUserDataEndpoints(t *testing.T) {
	service := Service{}
	endpoints := []scope.ServiceEndpoint{
		scope.ServiceEndpoint{
			URL:           "localhost",
			SigningRegion: "localhost",
			ServiceID:     "ssm",
		},
	}
	doc, _ := service.UserData("secretARN", 1, "eu-west-1", endpoints)

	_, err := mail.ReadMessage(bytes.NewBuffer(doc))
	if err != nil {
		t.Fatalf("Cannot parse MIME doc: %+v\n%s", err, string(doc))
	}
}

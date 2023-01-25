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

package v1beta2

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestSSHKeyName(t *testing.T) {
	tests := []struct {
		name       string
		sshKeyName *string
		wantErr    bool
	}{
		{
			name:       "SSH key name is nil is valid",
			sshKeyName: nil,
			wantErr:    false,
		},
		{
			name:       "SSH key name is an empty string is valid",
			sshKeyName: aws.String(""),
			wantErr:    false,
		},
		{
			name:       "SSH key name with alphanumeric characters is valid",
			sshKeyName: aws.String("test123"),
			wantErr:    false,
		},
		{
			name:       "SSH key name with underscore is valid",
			sshKeyName: aws.String("test_key"),
			wantErr:    false,
		},
		{
			name:       "SSH key name with dash is valid",
			sshKeyName: aws.String(`test-key`),
			wantErr:    false,
		},
		{
			name:       "SSH key name with tab is not valid",
			sshKeyName: aws.String("test-capi\t"),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cluster := &AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "cluster-",
					Namespace:    "default",
				},
				Spec: AWSClusterSpec{
					SSHKeyName: tt.sshKeyName,
				},
			}
			machine := &AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "machine-",
					Namespace:    "default",
				},
				Spec: AWSMachineSpec{
					SSHKeyName:   tt.sshKeyName,
					InstanceType: "test",
				},
			}
			for _, obj := range []client.Object{cluster, machine} {
				ctx := context.TODO()
				if err := testEnv.Create(ctx, obj); (err != nil) != tt.wantErr {
					t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

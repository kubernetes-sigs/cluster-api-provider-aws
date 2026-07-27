/*
Copyright 2026 The Kubernetes Authors.

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

package list

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	amiv1 "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/ami/v1beta1"
)

func newAMIList(items ...amiv1.AWSAMI) *amiv1.AWSAMIList {
	return &amiv1.AWSAMIList{
		TypeMeta: metav1.TypeMeta{
			Kind: amiv1.AWSAMIListKind,
		},
		Items: items,
	}
}

func newAMI(kubernetesVersion, os, region, imageID string) amiv1.AWSAMI {
	return amiv1.AWSAMI{
		TypeMeta: metav1.TypeMeta{
			Kind: amiv1.AWSAMIKind,
		},
		Spec: amiv1.AWSAMISpec{
			KubernetesVersion: kubernetesVersion,
			OS:                os,
			Region:            region,
			ImageID:           imageID,
		},
	}
}

func TestPrintAMIList(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		list      *amiv1.AWSAMIList
		wantErr   bool
		wantJSON  bool
		wantItems int
		wantText  string
	}{
		{
			name:      "empty list with json output is valid JSON",
			format:    "json",
			list:      newAMIList(),
			wantJSON:  true,
			wantItems: 0,
		},
		{
			name:     "empty list with table output prints message",
			format:   "table",
			list:     newAMIList(),
			wantText: "No AMIs found",
		},
		{
			name:      "non-empty list with json output round-trips items",
			format:    "json",
			list:      newAMIList(newAMI("v1.32.0", "ubuntu-24.04", "us-east-1", "ami-123")),
			wantJSON:  true,
			wantItems: 1,
		},
		{
			name:     "non-empty list with table output prints rows",
			format:   "table",
			list:     newAMIList(newAMI("v1.32.0", "ubuntu-24.04", "us-east-1", "ami-123")),
			wantText: "ami-123",
		},
		{
			name:     "empty list with yaml output is not the table message",
			format:   "yaml",
			list:     newAMIList(),
			wantText: "items:",
		},
		{
			name:    "unknown format returns error",
			format:  "xml",
			list:    newAMIList(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := printAMIList(&out, tt.format, tt.list)
			if tt.wantErr {
				if err == nil {
					t.Fatal("printAMIList() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("printAMIList() unexpected error: %v", err)
			}

			if tt.wantJSON {
				var got struct {
					Items []struct {
						Spec amiv1.AWSAMISpec `json:"spec"`
					} `json:"items"`
				}
				if err := json.Unmarshal(out.Bytes(), &got); err != nil {
					t.Fatalf("output is not valid JSON: %v\noutput: %s", err, out.String())
				}
				if len(got.Items) != tt.wantItems {
					t.Errorf("got %d items, want %d", len(got.Items), tt.wantItems)
				}
			}

			if tt.wantText != "" && !strings.Contains(out.String(), tt.wantText) {
				t.Errorf("output %q does not contain %q", out.String(), tt.wantText)
			}
		})
	}
}

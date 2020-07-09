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

package userdata

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

func TestNewNode(t *testing.T) {
	format.TruncatedDiff = false
	g := NewWithT(t)

	type args struct {
		input *NodeInput
	}
	tests := []struct {
		name          string
		args          args
		expectedBytes []byte
		expectErr     bool
	}{
		{
			name: "success case",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster
`),
			expectErr: false,
		},
		{
			name: "with extra args",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					KubeletExtraArgs: map[string]string{
						"foo":           "bar",
						"pizza-topping": "pepperoni",
					},
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --kubelet-extra-args '--foo=bar --pizza-topping=pepperoni'
`),
			expectErr: false,
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			bytes, err := NewNode(testcase.args.input)
			if testcase.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(string(bytes)).To(Equal(string(testcase.expectedBytes)))
		})
	}
}

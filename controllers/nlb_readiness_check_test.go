/*
Copyright 2019 The Kubernetes Authors.

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

package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestInjectNLBReadinessCheck(t *testing.T) {
	tests := []struct {
		name         string
		userData     []byte
		host         string
		port         int32
		expectErr    bool
		expectInject string
		expectOrder  []string
	}{
		{
			name: "standard worker cloud-init with kubeadm join",
			userData: []byte(`## template: jinja
#cloud-config

write_files:
-   path: /run/kubeadm/kubeadm-join-config.yaml
    owner: root:root
    permissions: '0640'
    content: |
      ---
      apiVersion: kubeadm.k8s.io/v1beta4
      kind: JoinConfiguration
runcmd:
  - kubeadm join --config /run/kubeadm/kubeadm-join-config.yaml  && echo success > /run/cluster-api/bootstrap-success.complete
`),
			host:      "my-nlb.elb.amazonaws.com",
			port:      6443,
			expectErr: false,
			expectOrder: []string{
				"Waiting for API server NLB",
				"kubeadm join --config",
			},
		},
		{
			name: "cloud-init with preKubeadmCommands before kubeadm join",
			userData: []byte(`## template: jinja
#cloud-config

runcmd:
  - "echo preKubeadm"
  - "apt-get update"
  - kubeadm join --config /run/kubeadm/kubeadm-join-config.yaml  && echo success > /run/cluster-api/bootstrap-success.complete
  - "echo postKubeadm"
`),
			host:      "test-nlb.us-east-1.elb.amazonaws.com",
			port:      6443,
			expectErr: false,
			expectOrder: []string{
				"echo preKubeadm",
				"apt-get update",
				"Waiting for API server NLB",
				"kubeadm join --config",
				"echo postKubeadm",
			},
		},
		{
			name: "custom API server port",
			userData: []byte(`runcmd:
  - kubeadm join --config /run/kubeadm/kubeadm-join-config.yaml  && echo success > /run/cluster-api/bootstrap-success.complete
`),
			host:         "my-nlb.elb.amazonaws.com",
			port:         8443,
			expectErr:    false,
			expectInject: "https://my-nlb.elb.amazonaws.com:8443/readyz",
		},
		{
			name:      "no kubeadm join line",
			userData:  []byte("runcmd:\n  - echo hello\n"),
			host:      "my-nlb.elb.amazonaws.com",
			port:      6443,
			expectErr: true,
		},
		{
			name:      "empty userdata",
			userData:  []byte{},
			host:      "my-nlb.elb.amazonaws.com",
			port:      6443,
			expectErr: true,
		},
		{
			name:      "nil userdata",
			userData:  nil,
			host:      "my-nlb.elb.amazonaws.com",
			port:      6443,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			result, err := injectNLBReadinessCheck(tt.userData, tt.host, tt.port)
			if tt.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())

			resultStr := string(result)

			g.Expect(resultStr).To(ContainSubstring("curl not found, cannot perform API server NLB readiness check, continuing anyway"))
			g.Expect(resultStr).To(ContainSubstring("Waiting for API server NLB"))
			g.Expect(resultStr).To(ContainSubstring(fmt.Sprintf("https://%s:%d/readyz", tt.host, tt.port)))

			g.Expect(resultStr).To(ContainSubstring("kubeadm join --config"))

			if tt.expectInject != "" {
				g.Expect(resultStr).To(ContainSubstring(tt.expectInject))
			}

			if len(tt.expectOrder) > 0 {
				prevIdx := 0
				for _, s := range tt.expectOrder {
					idx := strings.Index(resultStr[prevIdx:], s)
					g.Expect(idx).To(BeNumerically(">=", 0), "expected %q to appear after previous marker", s)
					prevIdx += idx + len(s)
				}
			}
		})
	}
}

func makeIgnitionJSON(t *testing.T, script string) []byte {
	t.Helper()
	ign := map[string]interface{}{
		"ignition": map[string]interface{}{"version": "2.3.0"},
		"storage": map[string]interface{}{
			"files": []interface{}{
				map[string]interface{}{
					"filesystem": "root",
					"path":       "/etc/kubeadm.sh",
					"contents": map[string]interface{}{
						"source": "data:," + url.PathEscape(script),
					},
					"mode": 448,
				},
				map[string]interface{}{
					"filesystem": "root",
					"path":       "/etc/kubeadm.yml",
					"contents": map[string]interface{}{
						"source": "data:,---%0Afoo%0A",
					},
					"mode": 384,
				},
			},
		},
	}
	data, err := json.Marshal(ign)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func decodeIgnitionKubeadmSh(t *testing.T, ignJSON []byte) string {
	t.Helper()
	var ign map[string]json.RawMessage
	if err := json.Unmarshal(ignJSON, &ign); err != nil {
		t.Fatal(err)
	}
	var storage map[string]json.RawMessage
	if err := json.Unmarshal(ign["storage"], &storage); err != nil {
		t.Fatal(err)
	}
	var files []map[string]json.RawMessage
	if err := json.Unmarshal(storage["files"], &files); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		var path string
		if err := json.Unmarshal(f["path"], &path); err != nil {
			continue
		}
		if path != "/etc/kubeadm.sh" {
			continue
		}
		var contents map[string]json.RawMessage
		if err := json.Unmarshal(f["contents"], &contents); err != nil {
			t.Fatal(err)
		}
		var source string
		if err := json.Unmarshal(contents["source"], &source); err != nil {
			t.Fatal(err)
		}
		decoded, err := url.PathUnescape(strings.TrimPrefix(source, "data:,"))
		if err != nil {
			t.Fatal(err)
		}
		return decoded
	}
	t.Fatal("could not find /etc/kubeadm.sh in Ignition JSON")
	return ""
}

func TestInjectNLBReadinessCheckIgnition(t *testing.T) {
	tests := []struct {
		name        string
		script      string
		host        string
		port        int32
		expectErr   bool
		expectOrder []string
	}{
		{
			name: "standard Ignition with kubeadm join",
			script: `#!/bin/bash
set -e

pre-command
another-pre-command

kubeadm join --config /etc/kubeadm.yml
mkdir -p /run/cluster-api && echo success > /run/cluster-api/bootstrap-success.complete
mv /etc/kubeadm.yml /tmp/

post-kubeadm-command
`,
			host:      "my-nlb.elb.amazonaws.com",
			port:      6443,
			expectErr: false,
			expectOrder: []string{
				"pre-command",
				"another-pre-command",
				"Waiting for API server NLB",
				"kubeadm join",
				"post-kubeadm-command",
			},
		},
		{
			name: "custom port",
			script: `#!/bin/bash
set -e
kubeadm join --config /etc/kubeadm.yml
mv /etc/kubeadm.yml /tmp/
`,
			host:      "my-nlb.elb.amazonaws.com",
			port:      8443,
			expectErr: false,
			expectOrder: []string{
				"https://my-nlb.elb.amazonaws.com:8443/readyz",
				"kubeadm join",
			},
		},
		{
			name: "no kubeadm join in script",
			script: `#!/bin/bash
set -e
echo hello
`,
			host:      "my-nlb.elb.amazonaws.com",
			port:      6443,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			var userData []byte
			if tt.script != "" {
				userData = makeIgnitionJSON(t, tt.script)
			}

			result, err := injectNLBReadinessCheckIgnition(userData, tt.host, tt.port)
			if tt.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())

			var parsed map[string]json.RawMessage
			g.Expect(json.Unmarshal(result, &parsed)).To(Succeed())

			script := decodeIgnitionKubeadmSh(t, result)

			g.Expect(script).To(ContainSubstring("curl not found, cannot perform API server NLB readiness check, continuing anyway"))
			g.Expect(script).To(ContainSubstring("Waiting for API server NLB"))
			g.Expect(script).To(ContainSubstring(fmt.Sprintf("https://%s:%d/readyz", tt.host, tt.port)))
			g.Expect(script).To(ContainSubstring("kubeadm join"))

			if len(tt.expectOrder) > 0 {
				prevIdx := 0
				for _, s := range tt.expectOrder {
					idx := strings.Index(script[prevIdx:], s)
					g.Expect(idx).To(BeNumerically(">=", 0), "expected %q to appear after previous marker", s)
					prevIdx += idx + len(s)
				}
			}
		})
	}
}

func TestInjectNLBReadinessCheckIgnitionEdgeCases(t *testing.T) {
	g := NewWithT(t)

	_, err := injectNLBReadinessCheckIgnition(nil, "host", 6443)
	g.Expect(err).To(HaveOccurred())

	_, err = injectNLBReadinessCheckIgnition([]byte("not json"), "host", 6443)
	g.Expect(err).To(HaveOccurred())

	_, err = injectNLBReadinessCheckIgnition([]byte(`{"ignition":{"version":"2.3.0"}}`), "host", 6443)
	g.Expect(err).To(HaveOccurred())

	noKubeadm := []byte(`{"ignition":{"version":"2.3.0"},"storage":{"files":[{"path":"/etc/other","contents":{"source":"data:,foo"}}]}}`)
	_, err = injectNLBReadinessCheckIgnition(noKubeadm, "host", 6443)
	g.Expect(err).To(HaveOccurred())
}

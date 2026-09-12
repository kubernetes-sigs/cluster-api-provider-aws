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

package userdata

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
)

func TestNewCustomHybridUserdata(t *testing.T) {
	format.TruncatedDiff = false
	g := NewWithT(t)

	tests := []struct {
		name         string
		template     string
		input        *NodeadmInput
		expectErr    bool
		errContains  string
		verifyOutput func(output string) bool
	}{
		{
			name: "basic template with all required variables",
			template: `#!/bin/bash
CLUSTER={{.ClusterName}}
REGION={{.Region}}
ACTIVATION_ID={{.ActivationID}}
ACTIVATION_CODE={{.ActivationCode}}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "CLUSTER=test-cluster") &&
					strings.Contains(output, "REGION=us-west-2") &&
					strings.Contains(output, "ACTIVATION_ID=act-123456") &&
					strings.Contains(output, "ACTIVATION_CODE=code-abcdef")
			},
		},
		{
			name: "template with kubernetes version",
			template: `#!/bin/bash
CLUSTER={{.ClusterName}}
REGION={{.Region}}
K8S_VERSION={{.KubernetesVersion}}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				KubernetesVersion: "1.29",
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "CLUSTER=test-cluster") &&
					strings.Contains(output, "REGION=us-west-2") &&
					strings.Contains(output, "K8S_VERSION=1.29")
			},
		},
		{
			name: "template with kubelet flags",
			template: `#!/bin/bash
{{- if .KubeletFlags }}
KUBELET_FLAGS="{{join .KubeletFlags " "}}"
{{- end }}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				KubeletFlags: []string{
					"--node-labels=env=prod",
					"--max-pods=110",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, `KUBELET_FLAGS="--node-labels=env=prod --max-pods=110"`)
			},
		},
		{
			name: "template with kubelet config",
			template: `#!/bin/bash
cat <<EOF > /etc/kubelet-config.yaml
{{.KubeletConfig}}
EOF
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				KubeletConfig: &runtime.RawExtension{
					Raw: []byte(`{"maxPods":110,"evictionHard":{"memory.available":"500Mi"}}`),
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "maxPods: 110") &&
					strings.Contains(output, "evictionHard:") &&
					strings.Contains(output, "memory.available: 500Mi")
			},
		},
		{
			name: "template with containerd config",
			template: `#!/bin/bash
cat <<EOF > /etc/containerd/config.toml
{{.ContainerdConfig}}
EOF
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				ContainerdConfig:  `[plugins."io.containerd.grpc.v1.cri"]`,
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, `[plugins."io.containerd.grpc.v1.cri"]`)
			},
		},
		{
			name: "template with Indent function",
			template: `config:
{{ Indent 2 .KubeletConfig }}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				KubeletConfig: &runtime.RawExtension{
					Raw: []byte(`{"maxPods":110,"evictionHard":{"memory.available":"500Mi"}}`),
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "  maxPods: 110") &&
					strings.Contains(output, "  evictionHard:")
			},
		},
		{
			name: "template with base64Encode function",
			template: `#!/bin/bash
ENCODED={{base64Encode .ActivationCode}}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "secret-code",
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				// base64 of "secret-code" is "c2VjcmV0LWNvZGU="
				return strings.Contains(output, "ENCODED=c2VjcmV0LWNvZGU=")
			},
		},
		{
			name: "template with default function",
			template: `#!/bin/bash
CONFIG={{default "default-config" .ContainerdConfig}}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				ContainerdConfig:  "", // empty, should use default
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "CONFIG=default-config")
			},
		},
		{
			name:     "template with trimSpace function",
			template: `TRIMMED={{trimSpace "  hello world  "}}`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "TRIMMED=hello world")
			},
		},
		{
			name: "template with string manipulation functions",
			template: `#!/bin/bash
LOWER={{lower .ClusterName}}
UPPER={{upper .Region}}
REPLACED={{replace .ClusterName "-" "_"}}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "Test-Cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "LOWER=test-cluster") &&
					strings.Contains(output, "UPPER=US-WEST-2") &&
					strings.Contains(output, "REPLACED=Test_Cluster")
			},
		},
		{
			name: "template with conditional logic",
			template: `#!/bin/bash
{{- if .KubeletFlags }}
HAS_FLAGS=true
{{- else }}
HAS_FLAGS=false
{{- end }}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				KubeletFlags:      []string{"--flag1"},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "HAS_FLAGS=true")
			},
		},
		{
			name: "template with range over kubelet flags",
			template: `#!/bin/bash
{{- range .KubeletFlags }}
echo "Flag: {{.}}"
{{- end }}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				KubeletFlags: []string{
					"--flag1=value1",
					"--flag2=value2",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, `echo "Flag: --flag1=value1"`) &&
					strings.Contains(output, `echo "Flag: --flag2=value2"`)
			},
		},
		{
			name: "realistic nodeadm template",
			template: `#!/bin/bash
set -euo pipefail

# Write nodeadm configuration
mkdir -p /etc/nodeadm
cat <<'NODECONFIG' > /etc/nodeadm/config.yaml
apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    name: {{.ClusterName}}
    region: {{.Region}}
  hybrid:
    ssm:
      activationId: {{.ActivationID}}
      activationCode: {{.ActivationCode}}
{{- if .KubeletFlags }}
  kubelet:
    flags:
{{- range .KubeletFlags }}
      - "{{.}}"
{{- end }}
{{- end }}
NODECONFIG

# Run nodeadm
/usr/local/bin/nodeadm init -c file:///etc/nodeadm/config.yaml
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "prod-cluster",
				Region:            "eu-west-1",
				KubernetesVersion: "1.29",
				ActivationID:      "act-prod-123",
				ActivationCode:    "code-prod-abc",
				KubeletFlags: []string{
					"--node-labels=node-type=hybrid",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "#!/bin/bash") &&
					strings.Contains(output, "set -euo pipefail") &&
					strings.Contains(output, "name: prod-cluster") &&
					strings.Contains(output, "region: eu-west-1") &&
					strings.Contains(output, "activationId: act-prod-123") &&
					strings.Contains(output, "activationCode: code-prod-abc") &&
					strings.Contains(output, `- "--node-labels=node-type=hybrid"`) &&
					strings.Contains(output, "/usr/local/bin/nodeadm init")
			},
		},
		{
			name: "template with contains/hasPrefix/hasSuffix",
			template: `#!/bin/bash
{{- if contains .ClusterName "prod" }}
IS_PROD=true
{{- end }}
{{- if hasPrefix .Region "us-" }}
IS_US=true
{{- end }}
{{- if hasSuffix .Region "-2" }}
IS_AZ2=true
{{- end }}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "my-prod-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "IS_PROD=true") &&
					strings.Contains(output, "IS_US=true") &&
					strings.Contains(output, "IS_AZ2=true")
			},
		},
		{
			name: "template with all NodeadmConfigSpec fields",
			template: `#!/bin/bash
{{- range .PreNodeadmCommands }}
PRE={{.}}
{{- end }}
{{- range .Files }}
FILE={{.Path}}:{{.Content}}
{{- end }}
{{- range .Users }}
USER={{.Name}}
{{- end }}
{{- if .NTP }}
NTP={{join .NTP.Servers ","}}
{{- end }}
{{- if .DiskSetup }}
{{- range .DiskSetup.Filesystems }}
FS={{.Device}}:{{.Filesystem}}
{{- end }}
{{- end }}
{{- range .Mounts }}
MOUNT={{index . 0}}:{{index . 1}}
{{- end }}
{{- range $gate, $enabled := .FeatureGates }}
FEATURE={{$gate}}:{{$enabled}}
{{- end }}
{{- if .ContainerdBaseRuntimeSpec }}
BASE_RUNTIME={{.ContainerdBaseRuntimeSpec}}
{{- end }}
`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
				FeatureGates: map[eksbootstrapv1.Feature]bool{
					eksbootstrapv1.FeatureFastImagePull: true,
				},
				PreNodeadmCommands: []string{"echo pre-nodeadm"},
				Files: []eksbootstrapv1.File{
					{Path: "/etc/example", Content: "file-content"},
				},
				Users: []eksbootstrapv1.User{
					{Name: "test-user"},
				},
				NTP: &eksbootstrapv1.NTP{
					Enabled: ptr.To(true),
					Servers: []string{"time.example.com"},
				},
				DiskSetup: &eksbootstrapv1.DiskSetup{
					Filesystems: []eksbootstrapv1.Filesystem{
						{Device: "/dev/xvdb", Filesystem: "ext4", Label: "data"},
					},
				},
				Mounts: []eksbootstrapv1.MountPoints{
					{"/dev/xvdb", "/data"},
				},
				ContainerdBaseRuntimeSpec: &runtime.RawExtension{
					Raw: []byte(`{"process":{"terminal":false}}`),
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "PRE=echo pre-nodeadm") &&
					strings.Contains(output, "FILE=/etc/example:file-content") &&
					strings.Contains(output, "USER=test-user") &&
					strings.Contains(output, "NTP=time.example.com") &&
					strings.Contains(output, "FS=/dev/xvdb:ext4") &&
					strings.Contains(output, "MOUNT=/dev/xvdb:/data") &&
					strings.Contains(output, "FEATURE=FastImagePull:true") &&
					strings.Contains(output, "terminal: false")
			},
		},
		// Error cases
		{
			name:     "missing cluster name",
			template: `{{.ClusterName}}`,
			input: &NodeadmInput{
				Hybrid:            true,
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr:   true,
			errContains: "cluster name is required",
		},
		{
			name:     "missing region",
			template: `{{.Region}}`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr:   true,
			errContains: "region is required",
		},
		{
			name:     "missing kubernetes version",
			template: `{{.KubernetesVersion}}`,
			input: &NodeadmInput{
				Hybrid:         true,
				ClusterName:    "test-cluster",
				Region:         "us-west-2",
				ActivationID:   "act-123456",
				ActivationCode: "code-abcdef",
			},
			expectErr:   true,
			errContains: "kubernetes version is required",
		},
		{
			name:     "missing activation ID",
			template: `{{.ActivationID}}`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationCode:    "code-abcdef",
			},
			expectErr:   true,
			errContains: "SSM activation ID is required",
		},
		{
			name:     "missing activation code",
			template: `{{.ActivationCode}}`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
			},
			expectErr:   true,
			errContains: "SSM activation code is required",
		},
		{
			name:        "nil input",
			template:    `{{.ClusterName}}`,
			input:       nil,
			expectErr:   true,
			errContains: "custom hybrid input is required",
		},
		{
			name:     "empty template",
			template: "",
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr:   true,
			errContains: "custom userdata template is required",
		},
		{
			name:     "invalid template syntax",
			template: `{{.ClusterName`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr:   true,
			errContains: "failed to parse custom userdata template",
		},
		{
			name:     "unknown field in template",
			template: `{{.UnknownField}}`,
			input: &NodeadmInput{
				Hybrid:            true,
				ClusterName:       "test-cluster",
				Region:            "us-west-2",
				KubernetesVersion: "1.29",
				ActivationID:      "act-123456",
				ActivationCode:    "code-abcdef",
			},
			expectErr:   true,
			errContains: "failed to execute custom userdata template",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := NewCustomHybridUserdata(tc.template, tc.input)

			if tc.expectErr {
				g.Expect(err).To(HaveOccurred())
				if tc.errContains != "" {
					g.Expect(err.Error()).To(ContainSubstring(tc.errContains))
				}
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			if tc.verifyOutput != nil {
				g.Expect(tc.verifyOutput(string(output))).To(BeTrue(),
					"Output verification failed for: %s\n\nActual output:\n%s", tc.name, string(output))
			}
		})
	}
}

func TestCustomHybridUserdataWithEmptyOptionalFields(t *testing.T) {
	g := NewWithT(t)

	template := `#!/bin/bash
CLUSTER={{.ClusterName}}
{{- if .KubeletFlags }}
FLAGS={{join .KubeletFlags ","}}
{{- end }}
{{- if .KubeletConfig }}
CONFIG={{.KubeletConfig}}
{{- end }}
{{- if .ContainerdConfig }}
CONTAINERD={{.ContainerdConfig}}
{{- end }}
`
	input := &NodeadmInput{
		Hybrid:            true,
		ClusterName:       "test-cluster",
		Region:            "us-west-2",
		KubernetesVersion: "1.29",
		ActivationID:      "act-123456",
		ActivationCode:    "code-abcdef",
		// All optional fields are empty/nil
		KubeletFlags:     nil,
		KubeletConfig:    nil,
		ContainerdConfig: "",
	}

	output, err := NewCustomHybridUserdata(template, input)
	g.Expect(err).NotTo(HaveOccurred())

	outputStr := string(output)
	g.Expect(outputStr).To(ContainSubstring("CLUSTER=test-cluster"))
	g.Expect(outputStr).NotTo(ContainSubstring("FLAGS="))
	g.Expect(outputStr).NotTo(ContainSubstring("CONFIG="))
	g.Expect(outputStr).NotTo(ContainSubstring("CONTAINERD="))
}

func TestCustomHybridUserdataComplexTemplate(t *testing.T) {
	g := NewWithT(t)

	// This tests a complex real-world-like template
	template := `#!/bin/bash
set -euo pipefail

# bootstrap script for {{.ClusterName}} in {{.Region}} (k8s: {{.KubernetesVersion}})

echo "Starting hybrid node bootstrap..."

# Create directories
mkdir -p /etc/nodeadm
mkdir -p /var/log/nodeadm

# Write nodeadm configuration
cat <<'EOF' > /etc/nodeadm/config.yaml
apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    name: {{.ClusterName}}
    region: {{.Region}}
  hybrid:
    ssm:
      activationId: {{.ActivationID}}
      activationCode: {{.ActivationCode}}
{{- if or .KubeletFlags .KubeletConfig }}
  kubelet:
{{- if .KubeletConfig }}
    config:
{{ Indent 6 .KubeletConfig }}
{{- end }}
{{- if .KubeletFlags }}
    flags:
{{- range .KubeletFlags }}
      - "{{.}}"
{{- end }}
{{- end }}
{{- end }}
{{- if .ContainerdConfig }}
  containerd:
    config: |
{{ Indent 6 .ContainerdConfig }}
{{- end }}
EOF

# Run pre-hooks
if [ -x /opt/scripts/pre-bootstrap.sh ]; then
    /opt/scripts/pre-bootstrap.sh
fi

# Initialize nodeadm
echo "Initializing nodeadm for k8s {{.KubernetesVersion}}..."
/usr/local/bin/nodeadm init -c file:///etc/nodeadm/config.yaml 2>&1 | tee /var/log/nodeadm/init.log

# Run post-hooks
if [ -x /opt/scripts/post-bootstrap.sh ]; then
    /opt/scripts/post-bootstrap.sh
fi

echo "Bootstrap complete for {{.ClusterName}}"
`

	input := &NodeadmInput{
		Hybrid:            true,
		ClusterName:       "test-clusterprod",
		Region:            "eu-central-1",
		KubernetesVersion: "1.29",
		ActivationID:      "act-xxx-123",
		ActivationCode:    "code-xxx-secret",
		KubeletFlags: []string{
			"--node-labels=environment=production,team=platform",
			"--register-with-taints=dedicated=isolated:NoSchedule",
		},
		KubeletConfig: &runtime.RawExtension{
			Raw: []byte(`{"maxPods":250,"evictionHard":{"memory.available":"1Gi","nodefs.available":"10%"}}`),
		},
		ContainerdConfig: `[plugins."io.containerd.grpc.v1.cri".containerd]
discard_unpacked_layers = false`,
	}

	output, err := NewCustomHybridUserdata(template, input)
	g.Expect(err).NotTo(HaveOccurred())

	outputStr := string(output)

	// Verify structure
	g.Expect(outputStr).To(ContainSubstring("#!/bin/bash"))
	g.Expect(outputStr).To(ContainSubstring("set -euo pipefail"))

	// Verify cluster info
	g.Expect(outputStr).To(ContainSubstring("name: test-clusterprod"))
	g.Expect(outputStr).To(ContainSubstring("region: eu-central-1"))

	// Verify Kubernetes version
	g.Expect(outputStr).To(ContainSubstring("k8s: 1.29"))
	g.Expect(outputStr).To(ContainSubstring("Initializing nodeadm for k8s 1.29"))

	// Verify SSM credentials
	g.Expect(outputStr).To(ContainSubstring("activationId: act-xxx-123"))
	g.Expect(outputStr).To(ContainSubstring("activationCode: code-xxx-secret"))

	// Verify kubelet config
	g.Expect(outputStr).To(ContainSubstring("maxPods: 250"))
	g.Expect(outputStr).To(ContainSubstring("memory.available:"))

	// Verify kubelet flags
	g.Expect(outputStr).To(ContainSubstring(`"--node-labels=environment=production,team=platform"`))
	g.Expect(outputStr).To(ContainSubstring(`"--register-with-taints=dedicated=isolated:NoSchedule"`))

	// Verify containerd config
	g.Expect(outputStr).To(ContainSubstring("discard_unpacked_layers = false"))

	// Verify comments and echo statements preserved
	g.Expect(outputStr).To(ContainSubstring("bootstrap script"))
	g.Expect(outputStr).To(ContainSubstring("Bootstrap complete for test-clusterprod"))
}

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
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/alessio/shellescape"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

const (
	defaultBootstrapCommand = "/etc/eks/bootstrap.sh"
	boundary                = "//"

	// AMIFamilyAL2 is the Amazon Linux 2 AMI family.
	AMIFamilyAL2 = "AmazonLinux2"
	// AMIFamilyAL2023 is the Amazon Linux 2023 AMI family.
	AMIFamilyAL2023 = "AmazonLinux2023"

	nodeUserData = `#cloud-config
{{template "files" .Files}}
runcmd:
{{- template "commands" .PreBootstrapCommands }}
  - {{ .BootstrapCommand }} {{.ClusterName}} {{- template "args" . }}
{{- template "commands" .PostBootstrapCommands }}
{{- template "ntp" .NTP }}
{{- template "users" .Users }}
{{- template "disk_setup" .DiskSetup}}
{{- template "fs_setup" .DiskSetup}}
{{- template "mounts" .Mounts}}
`

	// Shell script part template for AL2023.
	shellScriptPartTemplate = `--{{.Boundary}}
Content-Type: text/x-shellscript; charset="us-ascii"

#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset
{{- if or .PreBootstrapCommands .PostBootstrapCommands }}

{{- range .PreBootstrapCommands}}
{{.}}
{{- end}}
{{- range .PostBootstrapCommands}}
{{.}}
{{- end}}
{{- end}}`

	// Node config part template for AL2023.
	nodeConfigPartTemplate = `
--{{.Boundary}}
Content-Type: application/node.eks.aws

---
apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    name: {{.ClusterName}}
    apiServerEndpoint: {{.APIServerEndpoint}}
    certificateAuthority: {{.CACert}}
    cidr: {{if .ServiceCIDR}}{{.ServiceCIDR}}{{else}}172.20.0.0/16{{end}}
  kubelet:
    config:
      maxPods: {{.MaxPods}}
      {{- with .DNSClusterIP }}
      clusterDNS:
      - {{.}}
      {{- end }}
    flags:
    - "--node-labels={{.NodeLabels}}"

--{{.Boundary}}--`
)

// NodeInput contains all the information required to generate user data for a node.
type NodeInput struct {
	ClusterName           string
	KubeletExtraArgs      map[string]string
	ContainerRuntime      *string
	DNSClusterIP          *string
	DockerConfigJSON      *string
	APIRetryAttempts      *int
	PauseContainerAccount *string
	PauseContainerVersion *string
	UseMaxPods            *bool
	// NOTE: currently the IPFamily/ServiceIPV6Cidr isn't exposed to the user.
	// TODO (richardcase): remove the above comment when IPV6 / dual stack is implemented.
	IPFamily                 *string
	ServiceIPV6Cidr          *string
	PreBootstrapCommands     []string
	PostBootstrapCommands    []string
	BootstrapCommandOverride *string
	Files                    []eksbootstrapv1.File
	DiskSetup                *eksbootstrapv1.DiskSetup
	Mounts                   []eksbootstrapv1.MountPoints
	Users                    []eksbootstrapv1.User
	NTP                      *eksbootstrapv1.NTP

	// AMI Family Type to determine userdata format
	AMIFamilyType string

	// AL2023 specific fields
	AMIImageID        string
	APIServerEndpoint string
	Boundary          string
	CACert            string
	CapacityType      *v1beta2.ManagedMachinePoolCapacityType
	ServiceCIDR       string // Service CIDR range for the cluster
	ClusterDNS        string
	MaxPods           *int32
	NodeGroupName     string
	NodeLabels        string // Not exposed in CRD, computed from user input
}

// PauseContainerInfo holds pause container information for templates.
type PauseContainerInfo struct {
	AccountNumber *string
	Version       *string
}

// DockerConfigJSONEscaped returns the DockerConfigJSON escaped for use in cloud-init.
func (ni *NodeInput) DockerConfigJSONEscaped() string {
	if ni.DockerConfigJSON == nil || len(*ni.DockerConfigJSON) == 0 {
		return "''"
	}

	return shellescape.Quote(*ni.DockerConfigJSON)
}

// BootstrapCommand returns the bootstrap command to be used on a node instance.
func (ni *NodeInput) BootstrapCommand() string {
	if ni.BootstrapCommandOverride != nil && *ni.BootstrapCommandOverride != "" {
		return *ni.BootstrapCommandOverride
	}

	return defaultBootstrapCommand
}

// NewNode returns the user data string to be used on a node instance.
func NewNode(input *NodeInput) ([]byte, error) {
	// For AL2023, use the multipart MIME format
	if input.AMIFamilyType == AMIFamilyAL2023 {
		return generateAL2023UserData(input)
	}

	// For AL2 and other types, use the standard cloud-config format
	return generateStandardUserData(input)
}

// generateStandardUserData generates userdata for AL2 and other standard node types.
func generateStandardUserData(input *NodeInput) ([]byte, error) {
	tm := template.New("Node").Funcs(defaultTemplateFuncMap)

	if _, err := tm.Parse(filesTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse args template: %w", err)
	}

	if _, err := tm.Parse(argsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse args template: %w", err)
	}

	if _, err := tm.Parse(kubeletArgsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse kubeletExtraArgs template: %w", err)
	}

	if _, err := tm.Parse(commandsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse commandsTemplate template: %w", err)
	}

	if _, err := tm.Parse(ntpTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse ntp template: %w", err)
	}

	if _, err := tm.Parse(usersTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse users template: %w", err)
	}

	if _, err := tm.Parse(diskSetupTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse disk setup template: %w", err)
	}

	if _, err := tm.Parse(fsSetupTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse fs setup template: %w", err)
	}

	if _, err := tm.Parse(mountsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse mounts template: %w", err)
	}

	t, err := tm.Parse(nodeUserData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Node template: %w", err)
	}

	var out bytes.Buffer
	if err := t.Execute(&out, input); err != nil {
		return nil, fmt.Errorf("failed to generate Node template: %w", err)
	}

	return out.Bytes(), nil
}

// generateAL2023UserData generates userdata for Amazon Linux 2023 nodes.
func generateAL2023UserData(input *NodeInput) ([]byte, error) {
	if err := validateAL2023Input(input); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	// Write MIME header
	if _, err := buf.WriteString(fmt.Sprintf("MIME-Version: 1.0\nContent-Type: multipart/mixed; boundary=%q\n\n", input.Boundary)); err != nil {
		return nil, fmt.Errorf("failed to write MIME header: %v", err)
	}

	// Write shell script part if needed
	if len(input.PreBootstrapCommands) > 0 || len(input.PostBootstrapCommands) > 0 {
		shellScriptTemplate := template.Must(template.New("shell").Parse(shellScriptPartTemplate))
		if err := shellScriptTemplate.Execute(&buf, input); err != nil {
			return nil, fmt.Errorf("failed to execute shell script template: %v", err)
		}
		if _, err := buf.WriteString("\n"); err != nil {
			return nil, fmt.Errorf("failed to write newline: %v", err)
		}
	}

	// Write node config part
	nodeConfigTemplate := template.Must(template.New("node").Parse(nodeConfigPartTemplate))
	if err := nodeConfigTemplate.Execute(&buf, input); err != nil {
		return nil, fmt.Errorf("failed to execute node config template: %v", err)
	}

	return buf.Bytes(), nil
}

// getNodeLabels returns the string representation of node-labels flags for nodeadm.
func (ni *NodeInput) getNodeLabels() string {
	if ni.KubeletExtraArgs != nil {
		if _, ok := ni.KubeletExtraArgs["node-labels"]; ok {
			return ni.KubeletExtraArgs["node-labels"]
		}
	}
	nodeLabels := make([]string, 0, 3)
	if ni.AMIImageID != "" {
		nodeLabels = append(nodeLabels, fmt.Sprintf("eks.amazonaws.com/nodegroup-image=%s", ni.AMIImageID))
	}
	if ni.NodeGroupName != "" {
		nodeLabels = append(nodeLabels, fmt.Sprintf("eks.amazonaws.com/nodegroup=%s", ni.NodeGroupName))
	}
	nodeLabels = append(nodeLabels, fmt.Sprintf("eks.amazonaws.com/capacityType=%s", ni.getCapacityTypeString()))
	return strings.Join(nodeLabels, ",")
}

// getCapacityTypeString returns the string representation of the capacity type.
func (ni *NodeInput) getCapacityTypeString() string {
	if ni.CapacityType == nil {
		return "ON_DEMAND"
	}
	switch *ni.CapacityType {
	case v1beta2.ManagedMachinePoolCapacityTypeSpot:
		return "SPOT"
	case v1beta2.ManagedMachinePoolCapacityTypeOnDemand:
		return "ON_DEMAND"
	default:
		return strings.ToUpper(string(*ni.CapacityType))
	}
}

// validateAL2023Input validates the input for AL2023 user data generation.
func validateAL2023Input(input *NodeInput) error {
	if input.APIServerEndpoint == "" {
		return fmt.Errorf("API server endpoint is required for AL2023")
	}
	if input.CACert == "" {
		return fmt.Errorf("CA certificate is required for AL2023")
	}
	if input.ClusterName == "" {
		return fmt.Errorf("cluster name is required for AL2023")
	}
	if input.NodeGroupName == "" {
		return fmt.Errorf("node group name is required for AL2023")
	}

	if input.MaxPods == nil {
		if input.UseMaxPods != nil && *input.UseMaxPods {
			input.MaxPods = ptr.To[int32](110)
		} else {
			input.MaxPods = ptr.To[int32](58)
		}
	}
	if input.DNSClusterIP != nil {
		input.ClusterDNS = *input.DNSClusterIP
	}

	if input.Boundary == "" {
		input.Boundary = boundary
	}
	input.NodeLabels = input.getNodeLabels()

	klog.V(2).Infof("AL2023 Userdata Generation - maxPods: %d, node-labels: %s",
		*input.MaxPods, input.NodeLabels)

	return nil
}

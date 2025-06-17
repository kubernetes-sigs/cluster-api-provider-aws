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

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
)

const (
	defaultBootstrapCommand = "/etc/eks/bootstrap.sh"

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
)

// NodeInput defines the context to generate a node user data.
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

// AL2023UserDataInput defines the required input for generating AL2023 userdata
type AL2023UserDataInput struct {
	ClusterName       string
	APIServerEndpoint string
	CACert            string
	NodeGroupName     string
	MaxPods           int
	ClusterDNS        string
	AMIImageID        string
	CapacityType      string
}

// ValidateAL2023UserDataInput validates the input for AL2023 userdata generation
func ValidateAL2023UserDataInput(input *AL2023UserDataInput) error {
	if input.ClusterName == "" {
		return fmt.Errorf("cluster name is required")
	}
	if input.APIServerEndpoint == "" {
		return fmt.Errorf("API server endpoint is required")
	}
	if !strings.HasPrefix(input.APIServerEndpoint, "https://") {
		return fmt.Errorf("API server endpoint must start with https://")
	}
	if input.CACert == "" {
		return fmt.Errorf("CA certificate is required")
	}
	if input.NodeGroupName == "" {
		return fmt.Errorf("node group name is required")
	}
	if input.MaxPods <= 0 {
		return fmt.Errorf("max pods must be greater than 0")
	}
	if input.ClusterDNS == "" {
		return fmt.Errorf("cluster DNS is required")
	}
	if input.AMIImageID == "" {
		return fmt.Errorf("AMI image ID is required")
	}
	if input.CapacityType == "" {
		return fmt.Errorf("capacity type is required")
	}
	return nil
}

// GenerateAL2023UserData generates userdata for Amazon Linux 2023 nodes with validation and retry
func GenerateAL2023UserData(input *AL2023UserDataInput) ([]byte, error) {
	// Validate input
	if err := ValidateAL2023UserDataInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Generate userdata with validated input
	userData := fmt.Sprintf(`MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="//"

--//
Content-Type: application/node.eks.aws

---
apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    apiServerEndpoint: %s
    certificateAuthority: %s
    cidr: 10.96.0.0/12
    name: %s
  kubelet:
    config:
      maxPods: %d
      clusterDNS:
      - %s
    flags:
    - "--node-labels=eks.amazonaws.com/nodegroup-image=%s,eks.amazonaws.com/capacityType=%s,eks.amazonaws.com/nodegroup=%s"

--//--`,
		input.APIServerEndpoint,
		input.CACert,
		input.ClusterName,
		input.MaxPods,
		input.ClusterDNS,
		input.AMIImageID,
		input.CapacityType,
		input.NodeGroupName)

	return []byte(userData), nil
}

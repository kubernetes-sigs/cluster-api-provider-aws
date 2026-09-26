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
	"bytes"
	"fmt"
	texttemplate "text/template"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	bootstraptemplate "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/internal/template"
)

// validateCustomHybridNodeadmInput validates fields required for custom hybrid userdata generation.
func validateCustomHybridNodeadmInput(input *NodeadmInput) error {
	if input == nil {
		return fmt.Errorf("custom hybrid input is required")
	}
	if !input.Hybrid {
		return fmt.Errorf("hybrid mode is required for custom hybrid userdata")
	}
	if err := validateNodeadmInput(input); err != nil {
		return err
	}
	if input.KubernetesVersion == "" {
		return fmt.Errorf("kubernetes version is required for custom hybrid userdata")
	}
	return nil
}

// NewCustomHybridUserdata generates a generic userdata from a user-provided template
// with runtime variable interpolation. The template uses Go text/template syntax.
func NewCustomHybridUserdata(templateStr string, input *NodeadmInput) ([]byte, error) {
	if err := validateCustomHybridNodeadmInput(input); err != nil {
		return nil, err
	}

	if templateStr == "" {
		return nil, fmt.Errorf("custom userdata template is required")
	}

	// Convert KubeletConfig RawExtension to YAML string for template use
	kubeletConfigStr := ""
	if input.KubeletConfig != nil {
		var err error
		kubeletConfigStr, err = bootstraptemplate.ToYAML(input.KubeletConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to convert kubelet config to YAML: %w", err)
		}
	}
	containerdBaseRuntimeSpecStr := ""
	if input.ContainerdBaseRuntimeSpec != nil {
		var err error
		containerdBaseRuntimeSpecStr, err = bootstraptemplate.ToYAML(input.ContainerdBaseRuntimeSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to convert containerd base runtime spec to YAML: %w", err)
		}
	}

	// Build template data with KubeletConfig as a rendered string
	templateData := struct {
		ClusterName               string
		Region                    string
		KubernetesVersion         string
		ActivationID              string
		ActivationCode            string
		KubeletFlags              []string
		KubeletConfig             string
		ContainerdConfig          string
		ContainerdBaseRuntimeSpec string
		FeatureGates              map[eksbootstrapv1.Feature]bool
		PreNodeadmCommands        []string
		Files                     []eksbootstrapv1.File
		Users                     []eksbootstrapv1.User
		NTP                       *eksbootstrapv1.NTP
		DiskSetup                 *eksbootstrapv1.DiskSetup
		Mounts                    []eksbootstrapv1.MountPoints
	}{
		ClusterName:               input.ClusterName,
		Region:                    input.Region,
		KubernetesVersion:         input.KubernetesVersion,
		ActivationID:              input.ActivationID,
		ActivationCode:            input.ActivationCode,
		KubeletFlags:              input.KubeletFlags,
		KubeletConfig:             kubeletConfigStr,
		ContainerdConfig:          input.ContainerdConfig,
		ContainerdBaseRuntimeSpec: containerdBaseRuntimeSpecStr,
		FeatureGates:              input.FeatureGates,
		PreNodeadmCommands:        input.PreNodeadmCommands,
		Files:                     input.Files,
		Users:                     input.Users,
		NTP:                       input.NTP,
		DiskSetup:                 input.DiskSetup,
		Mounts:                    input.Mounts,
	}

	// Parse the user-provided template
	tmpl, err := texttemplate.New("customHybridUserdata").
		Funcs(bootstraptemplate.FuncMap()).
		Option("missingkey=error"). // Fail on missing keys for better error messages
		Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse custom userdata template: %w", err)
	}

	// Execute the template with the provided input
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return nil, fmt.Errorf("failed to execute custom userdata template: %w", err)
	}

	return buf.Bytes(), nil
}

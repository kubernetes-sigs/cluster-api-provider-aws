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
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"

	"k8s.io/apimachinery/pkg/runtime"
)

// CustomHybridInput contains all runtime variables available for custom template interpolation.
// These variables are populated at runtime and made available to user-provided templates.
type CustomHybridInput struct {
	// ClusterName is the EKS cluster name (required).
	ClusterName string

	// Region is the AWS region where the cluster is located (required).
	Region string

	// KubernetesVersion is the Kubernetes version of the cluster (required).
	// Examples: "1.29", "v1.29.0"
	KubernetesVersion string

	// ActivationID is the SSM activation ID for hybrid node registration (required).
	ActivationID string

	// ActivationCode is the SSM activation code for hybrid node registration (required).
	ActivationCode string

	// KubeletFlags contains additional kubelet command-line flags (optional).
	KubeletFlags []string

	// KubeletConfig contains the kubelet configuration as a RawExtension (optional).
	// It will be converted to a YAML string before template rendering.
	KubeletConfig *runtime.RawExtension

	// ContainerdConfig contains the containerd configuration (optional).
	ContainerdConfig string
}

// customTemplateFuncMap provides template functions for custom userdata templates.
// These functions enable common operations within user-provided templates.
var customTemplateFuncMap = template.FuncMap{
	// Indent adds indentation to each line of content
	"Indent": templateYAMLIndent,
	// join concatenates slice elements with a separator
	"join": strings.Join,
	// base64Encode encodes a string to base64
	"base64Encode": func(s string) string {
		return base64.StdEncoding.EncodeToString([]byte(s))
	},
	// default returns the default value if the given value is empty
	"default": func(def, val string) string {
		if val == "" {
			return def
		}
		return val
	},
	// trimSpace removes leading and trailing whitespace
	"trimSpace": strings.TrimSpace,
	// contains checks if a string contains a substring
	"contains": strings.Contains,
	// hasPrefix checks if a string has a prefix
	"hasPrefix": strings.HasPrefix,
	// hasSuffix checks if a string has a suffix
	"hasSuffix": strings.HasSuffix,
	// replace replaces all occurrences of old with new in s
	"replace": strings.ReplaceAll,
	// lower converts a string to lowercase
	"lower": strings.ToLower,
	// upper converts a string to uppercase
	"upper": strings.ToUpper,
}

// validateCustomHybridInput validates the required fields for custom hybrid userdata generation.
func validateCustomHybridInput(input *CustomHybridInput) error {
	if input == nil {
		return fmt.Errorf("custom hybrid input is required")
	}
	if input.ClusterName == "" {
		return fmt.Errorf("cluster name is required for custom hybrid userdata")
	}
	if input.Region == "" {
		return fmt.Errorf("region is required for custom hybrid userdata")
	}
	if input.KubernetesVersion == "" {
		return fmt.Errorf("kubernetes version is required for custom hybrid userdata")
	}
	if input.ActivationID == "" {
		return fmt.Errorf("SSM activation ID is required for custom hybrid userdata")
	}
	if input.ActivationCode == "" {
		return fmt.Errorf("SSM activation code is required for custom hybrid userdata")
	}
	return nil
}

// NewCustomHybridUserdata generates a generic userdata from a user-provided template
// with runtime variable interpolation. The template uses Go text/template syntax.
func NewCustomHybridUserdata(templateStr string, input *CustomHybridInput) ([]byte, error) {
	if err := validateCustomHybridInput(input); err != nil {
		return nil, err
	}

	if templateStr == "" {
		return nil, fmt.Errorf("custom userdata template is required")
	}

	// Convert KubeletConfig RawExtension to YAML string for template use
	kubeletConfigStr := ""
	if input.KubeletConfig != nil {
		var err error
		kubeletConfigStr, err = templateToYAML(input.KubeletConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to convert kubelet config to YAML: %w", err)
		}
	}

	// Build template data with KubeletConfig as a rendered string
	templateData := struct {
		ClusterName       string
		Region            string
		KubernetesVersion string
		ActivationID      string
		ActivationCode    string
		KubeletFlags      []string
		KubeletConfig     string
		ContainerdConfig  string
	}{
		ClusterName:       input.ClusterName,
		Region:            input.Region,
		KubernetesVersion: input.KubernetesVersion,
		ActivationID:      input.ActivationID,
		ActivationCode:    input.ActivationCode,
		KubeletFlags:      input.KubeletFlags,
		KubeletConfig:     kubeletConfigStr,
		ContainerdConfig:  input.ContainerdConfig,
	}

	// Parse the user-provided template
	tmpl, err := template.New("customHybridUserdata").
		Funcs(customTemplateFuncMap).
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

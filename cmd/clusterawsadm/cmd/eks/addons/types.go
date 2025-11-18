/*
Copyright 2021 The Kubernetes Authors.

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

package addons

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type availableAddon struct {
	// +required
	Name            string          `json:"name"`
	// +required
	Type            string          `json:"type"`
	// +required
	Version         string          `json:"version"`
	// +required
	Architecture    []string        `json:"architecture"`
	// +required
	Compatibilities []compatibility `json:"compatibilities"`
}

type compatibility struct {
	// +required
	ClusterVersion   string   `json:"clusterVersion"`
	// +required
	DefaultVersion   bool     `json:"defaultVersion"`
	// +required
	PlatformVersions []string `json:"platformVersions"`
}

func (c compatibility) String() string {
	return fmt.Sprintf("Cluster version: %s, default: %t, platform versions: %q", c.ClusterVersion, c.DefaultVersion, c.PlatformVersions)
}

type availableAddonsList struct {
	// +required
	Cluster string           `json:"cluster"`
	// +required
	Addons  []availableAddon `json:"addons"`
}

func (a *availableAddonsList) ToTable() *metav1.Table {
	table := &metav1.Table{
		TypeMeta: metav1.TypeMeta{
			APIVersion: metav1.SchemeGroupVersion.String(),
			Kind:       "Table",
		},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Type",
				Type: "string",
			},
			{
				Name: "Version",
				Type: "string",
			},
			{
				Name: "Architectures",
				Type: "string",
			},
			{
				Name: "Compatibilities",
				Type: "string",
			},
		},
	}

	for _, addon := range a.Addons {
		row := metav1.TableRow{
			Cells: []interface{}{addon.Name, addon.Type, addon.Version, addon.Architecture, addon.Compatibilities},
		}
		table.Rows = append(table.Rows, row)
	}

	return table
}

type installedAddon struct {
	// +optional
	Name    string
	// +optional
	Version string

	// +optional
	AddonARN string
	// +optional
	RoleARN  *string

	// +optional
	Status string
	// +optional
	Tags   map[string]string

	// +optional
	HealthIssues []issue

	// +optional
	CreatedAt  time.Time
	// +optional
	ModifiedAt time.Time
}

type issue struct {
	// +optional
	Code        string
	// +optional
	Message     string
	// +optional
	ResourceIDs []string
}

type installedAddonsList struct {
	// +required
	Cluster string           `json:"cluster"`
	// +required
	Addons  []installedAddon `json:"addons"`
}

func (a *installedAddonsList) ToTable() *metav1.Table {
	//NOTE: the ARN hasn't been included as it makes the table too wide
	table := &metav1.Table{
		TypeMeta: metav1.TypeMeta{
			APIVersion: metav1.SchemeGroupVersion.String(),
			Kind:       "Table",
		},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Version",
				Type: "string",
			},
			{
				Name: "Status",
				Type: "string",
			},
			{
				Name: "Created",
				Type: "string",
			},
			{
				Name: "Modified",
				Type: "string",
			},
			{
				Name: "SA ARN",
				Type: "string",
			},
			{
				Name: "Issues",
				Type: "string",
			},
		},
	}

	for _, addon := range a.Addons {
		if addon.RoleARN == nil {
			addon.RoleARN = aws.String("")
		}

		row := metav1.TableRow{
			Cells: []interface{}{addon.Name, addon.Version, addon.Status, addon.CreatedAt, addon.ModifiedAt, *addon.RoleARN, len(addon.HealthIssues)},
		}
		table.Rows = append(table.Rows, row)
	}
	return table
}

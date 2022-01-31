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

	"github.com/aws/aws-sdk-go/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type availableAddon struct {
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	Version         string          `json:"version"`
	Architecture    []string        `json:"architecture"`
	Compatibilities []compatibility `json:"compatibilities"`
}

type compatibility struct {
	ClusterVersion   string   `json:"clusterVersion"`
	DefaultVersion   bool     `json:"defaultVersion"`
	PlatformVersions []string `json:"platformVersions"`
}

func (c compatibility) String() string {
	return fmt.Sprintf("Cluster version: %s, default: %t, platform versions: %q", c.ClusterVersion, c.DefaultVersion, c.PlatformVersions)
}

type availableAddonsList struct {
	Cluster string           `json:"cluster"`
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
	Name    string
	Version string

	AddonARN string
	RoleARN  *string

	Status string
	Tags   map[string]*string

	HealthIssues []issue

	CreatedAt  time.Time
	ModifiedAt time.Time
}

type issue struct {
	Code        string
	Message     string
	ResourceIds []string
}

type installedAddonsList struct {
	Cluster string           `json:"cluster"`
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

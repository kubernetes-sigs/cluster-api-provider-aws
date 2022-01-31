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

package resource

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// AWSResource defines an AWS resource.
type AWSResource struct {
	Partition string `json:"partition"`
	Service   string `json:"service"`
	Region    string `json:"region"`
	AccountID string `json:"account_id"`
	Resource  string `json:"resource"`
	ARN       string `json:"arn"`
}

// AWSResourceList defines list of AWSResources.
type AWSResourceList struct {
	ClusterName  string        `json:"cluster_name"`
	AWSResources []AWSResource `json:"aws_resources"`
}

// ToTable converts AWSResourceList to Table.
func (a *AWSResourceList) ToTable() *metav1.Table {
	table := &metav1.Table{
		TypeMeta: metav1.TypeMeta{
			APIVersion: metav1.SchemeGroupVersion.String(),
			Kind:       "Table",
		},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{
				Name: "Partition",
				Type: "string",
			},
			{
				Name: "Service",
				Type: "string",
			},
			{
				Name: "Region",
				Type: "string",
			},
			{
				Name: "AccountID",
				Type: "string",
			},
			{
				Name: "Resource",
				Type: "string",
			},
			{
				Name: "ARN",
				Type: "string",
			},
		},
	}

	for _, resource := range a.AWSResources {
		row := metav1.TableRow{
			Cells: []interface{}{resource.Partition, resource.Service, resource.Region, resource.AccountID, resource.Resource, resource.ARN},
		}
		table.Rows = append(table.Rows, row)
	}
	return table
}

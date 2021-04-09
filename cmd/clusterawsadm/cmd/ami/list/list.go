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

package copy

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/flags"
	cmdout "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/printers"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

var (
	kubernetesVersion string
	opSystem          string
	outputPrinter     string
)

type amiInfo struct {
	OS                string `json:"os"`
	Region            string `json:"region"`
	ID                string `json:"amiID"`
	CreationDate      string `json:"creationDate"`
	KubernetesVersion string `json:"kubernetesVersion"`
	Name              string `json:"name"`
}

func ListAMICmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "list",
		Short: "List AMIs from the default AWS account where AMIs are stored",
		Long: cmd.LongDesc(`
			List AMIs based on Kubernetes version, OS, region. If no arguments are provided,
			it will print all AMIs in all regions, OS types for the supported Kubernetes versions.
            Supported Kubernetes versions start from the latest stable version and goes 2 release back:
			if the latest stable release is v1.20.4- v1.19.x and v1.18.x are supported.
			Note: First release of each version will be skipped, e.g., v1.21.0
			To list AMIs of unsupported Kubernetes versions, --kubernetes-version flag needs to be provided.
		`),
		Example: cmd.Examples(`
		# List AMIs from the default AWS account where AMIs are stored.
		# Available os options: centos-7, ubuntu-18.04, ubuntu-20.04, amazon-2
		clusterawsadm ami list --kubernetes-version=v1.18.12 --os=ubuntu-20.04  --region=us-west-2
		# To list all supported AMIs in all supported Kubernetes versions, regions, and linux distributions:
		clusterawsadm ami list
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			supportedOsList := []string{}
			if opSystem == "" {
				supportedOsList = getSupportedOsList()
			} else {
				supportedOsList = append(supportedOsList, opSystem)
			}
			imageRegionList := []string{}
			region := cmd.Flags().Lookup("region").Value.String()
			if region == "" {
				imageRegionList = getimageRegionList()
			} else {
				imageRegionList = append(imageRegionList, region)
			}

			supportedVersions := []string{}
			if kubernetesVersion == "" {
				var err error
				supportedVersions, err = getSupportedKubernetesVersions()
				if err != nil {
					fmt.Println("Failed to calculate supported Kubernetes versions")
					return err
				}
			} else {
				supportedVersions = append(supportedVersions, kubernetesVersion)
			}

			var sessionCache sync.Map
			imageMap := make(map[string][]*ec2.Image)
			for _, region := range imageRegionList {
				var sess *session.Session
				var err error
				if s, ok := sessionCache.Load(region); ok {
					sess = s.(*session.Session)
				} else {
					sess, err = session.NewSessionWithOptions(session.Options{
						SharedConfigState: session.SharedConfigEnable,
						Config:            aws.Config{Region: aws.String(region)},
					})
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						return err
					}
					sessionCache.Store(region, sess)
				}

				ec2Client := ec2.New(sess)
				imagesForRegion, err := getAllImages(ec2Client, "")
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return err
				}

				for key, image := range imagesForRegion {
					images, ok := imageMap[key]
					if !ok {
						images = make([]*ec2.Image, 0)
					}
					imageMap[key] = append(images, image...)
				}
			}

			listByVersion := amiList{
				AmiList: []amiInfo{},
			}
			for _, version := range supportedVersions {
				for _, region := range imageRegionList {
					for _, os := range supportedOsList {
						image, err := findAMI(imageMap, os, version)
						if err != nil {
							return err
						}
						listByVersion.AmiList = append(listByVersion.AmiList, amiInfo{
							OS:                os,
							Region:            region,
							ID:                *image.ImageId,
							CreationDate:      *image.CreationDate,
							KubernetesVersion: version,
							Name:              *image.Name,
						})
					}
				}
			}
			printer, err := cmdout.New(outputPrinter, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed creating output printer: %w", err)
			}

			if outputPrinter == string(cmdout.PrinterTypeTable) {
				table := listByVersion.ToTable()
				printer.Print(table)
			} else {
				printer.Print(listByVersion)
			}

			return nil
		},
	}

	flags.AddRegionFlag(newCmd)
	addOsFlag(newCmd)
	addKubernetesVersionFlag(newCmd)
	addOutputFlag(newCmd)
	return newCmd
}

func addOsFlag(c *cobra.Command) {
	c.Flags().StringVar(&opSystem, "os", "", "Operating system of the AMI to be copied")
}

func addKubernetesVersionFlag(c *cobra.Command) {
	c.Flags().StringVar(&kubernetesVersion, "kubernetes-version", "", "Kubernetes version of the AMI to be copied")
}

func addOutputFlag(c *cobra.Command) {
	c.Flags().StringVarP(&outputPrinter, "output", "o", "table", "The output format of the results. Possible values: table,json,yaml")
}

type amiList struct {
	AmiList []amiInfo `json:"AMIs"`
}

func (a *amiList) ToTable() *metav1.Table {
	table := &metav1.Table{
		TypeMeta: metav1.TypeMeta{
			APIVersion: metav1.SchemeGroupVersion.String(),
			Kind:       "Table",
		},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{
				Name: "Kubernetes-version",
				Type: "string",
			},
			{
				Name: "Region",
				Type: "string",
			},
			{
				Name: "OS",
				Type: "string",
			},
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Ami-id",
				Type: "string",
			},
		},
	}

	for _, ami := range a.AmiList {

		row := metav1.TableRow{
			Cells: []interface{}{ami.KubernetesVersion, ami.Region, ami.OS, ami.Name, ami.ID},
		}
		table.Rows = append(table.Rows, row)

	}
	return table
}

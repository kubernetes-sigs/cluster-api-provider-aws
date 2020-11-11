// +build e2e

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

package shared

import "flag"

func CreateDefaultFlags(ctx *E2EContext) {
	flag.StringVar(&ctx.Settings.ConfigPath, "config-path", "", "path to the e2e config file")
	flag.StringVar(&ctx.Settings.ArtifactFolder, "artifacts-folder", "", "folder where e2e test artifact should be stored")
	flag.BoolVar(&ctx.Settings.UseCIArtifacts, "kubetest.use-ci-artifacts", false, "use the latest build from the main branch of the Kubernetes repository")
	flag.StringVar(&ctx.Settings.KubetestConfigFilePath, "kubetest.config-file", "", "path to the kubetest configuration file")
	flag.IntVar(&ctx.Settings.GinkgoNodes, "kubetest.ginkgo-nodes", 1, "number of ginkgo nodes to use")
	flag.IntVar(&ctx.Settings.GinkgoSlowSpecThreshold, "kubetest.ginkgo-slowSpecThreshold", 120, "time in s before spec is marked as slow")
	flag.BoolVar(&ctx.Settings.UseExistingCluster, "use-existing-cluster", false, "if true, the test uses the current cluster instead of creating a new one (default discovery rules apply)")
	flag.BoolVar(&ctx.Settings.SkipCleanup, "skip-cleanup", false, "if true, the resource cleanup after tests will be skipped")
	flag.BoolVar(&ctx.Settings.SkipCloudFormationDeletion, "skip-cloudformation-deletion", false, "if true, an AWS CloudFormation stack will not be deleted")
	flag.BoolVar(&ctx.Settings.SkipCloudFormationCreation, "skip-cloudformation-creation", false, "if true, an AWS CloudFormation stack will not be created")
	//flag.StringVar(&ctx.DataFolder, "data-folder", "", "path to the data folder")
}

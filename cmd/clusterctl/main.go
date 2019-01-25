/*
Copyright 2018 The Kubernetes Authors.

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

package main

import (
	"sigs.k8s.io/cluster-api-provider-aws/cmd/versioninfo"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
)

func registerCustomCommands() {
	cmd.RootCmd.AddCommand(versioninfo.VersionCmd())
}

func main() {
	clusterActuator := cluster.NewActuator(cluster.ActuatorParams{})
	common.RegisterClusterProvisioner("aws", clusterActuator)
	registerCustomCommands()
	cmd.Execute()
}

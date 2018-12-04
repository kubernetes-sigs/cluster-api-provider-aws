// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"

	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/versioninfo"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
)

// initLogs is a temporary hack to enable proper logging until upstream dependencies
// are migrated to fully utilize klog instead of glog.
func initLogs() {
	flag.Set("logtostderr", "true")
	flags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(flags)
	flags.Set("alsologtostderr", "true")
	flags.Set("v", "4")
	flag.Parse()
}

func registerCustomCommands() {
	cmd.RootCmd.AddCommand(versioninfo.VersionCmd())
}

func main() {
	initLogs()
	clusterActuator := cluster.NewActuator(cluster.ActuatorParams{})
	common.RegisterClusterProvisioner("aws", clusterActuator)
	registerCustomCommands()
	cmd.Execute()
}

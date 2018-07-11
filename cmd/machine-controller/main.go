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
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/util/logs"
	"sigs.k8s.io/cluster-api/pkg/controller/config"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/controllers/machine"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/controllers/machine/options"
)

func init() {
	config.ControllerConfig.AddFlags(pflag.CommandLine)
}

func main() {
	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
	pflag.Parse()

	logs.InitLogs()
	defer logs.FlushLogs()

	machineServer := options.NewServer()
	if err := machine.Run(machineServer); err != nil {
		glog.Errorf("Failed to start the machine controller. Err: %v", err)
	}
}

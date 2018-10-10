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

package cluster

import (
	"github.com/golang/glog"

	clusteractuator "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clustercommon "sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
)

// ProviderName is the name of the cloud provider
const ProviderName = "aws"

func init() {
	codec, err := v1alpha1.NewCodec()
	if err != nil {
		glog.Fatalf("Could not create codec: %v", err)
	}

	params := clusteractuator.ActuatorParams{
		Codec: codec,
	}

	actuator, err := clusteractuator.NewActuator(params)
	if err != nil {
		glog.Fatalf("Could not create aws cluster actuator: %v", err)
	}

	clustercommon.RegisterClusterProvisioner(ProviderName, actuator)

}

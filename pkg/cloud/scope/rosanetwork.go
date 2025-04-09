 /*
 Copyright The Kubernetes Authors.
 
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

package scope

import (
	// "context"

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	rosanetworkv1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"k8s.io/klog/v2"
)

type RosaNetworkScopeParams struct  {
	Client         client.Client
	Logger         *logger.Logger
	ControllerName string
}

type RosaNetworkScope struct {
	Logger		    logger.Logger
	Client		    client.Client
	session		    awsclient.ConfigProvider
	ControllerName  string
	RosaNetwork     *rosanetworkv1.RosaNetwork
}

// NewROSAControlPlaneScope creates a new ROSAControlPlaneScope from the supplied parameters.
func NewRosaNetworkScope(params RosaNetworkScopeParams) (*RosaNetworkScope, error) {
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	rosaNetworkScope := &RosaNetworkScope{
	    Logger:         *params.Logger,
	    Client:         params.Client,
	    ControllerName: params.ControllerName,
	}

	return rosaNetworkScope, nil
}

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

package controller

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api-provider-aws/version"
)

// BootstrapCredsSecret defines the tag for capa manager bootstrap credentials.
const BootstrapCredsSecret = "capa-manager-bootstrap-credentials"

// GetClient creates the config for a kubernetes client and returns a client-go client for the cluster.
func GetClient(kubeconfigPath string, kubeconfigContext string) (*kubernetes.Clientset, error) {
	// If a kubeconfig file isn't provided, find one in the standard locations.
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfigPath != "" {
		rules.ExplicitPath = kubeconfigPath
	}

	config, err := rules.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load Kubeconfig")
	}

	configOverrides := &clientcmd.ConfigOverrides{}
	if kubeconfigContext == "" {
		configOverrides.CurrentContext = kubeconfigContext
	}

	restConfig, err := clientcmd.NewDefaultClientConfig(*config, configOverrides).ClientConfig()
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid configuration:") {
			return nil, errors.New(strings.Replace(err.Error(), "invalid configuration:", "invalid kubeconfig file; clusterawsadm requires a valid kubeconfig file to connect to the management cluster:", 1))
		}
		return nil, err
	}
	restConfig.UserAgent = fmt.Sprintf("clusterawsadm/%s (%s)", version.Get().GitVersion, version.Get().Platform)

	// Get a client-go client for the cluster
	cs, err := kubernetes.NewForConfig(restConfig)
	return cs, err
}

// PrintBootstrapCredentials will print the bootstrap credentials.
func PrintBootstrapCredentials(secret *corev1.Secret) {
	if creds, ok := secret.Data["credentials"]; ok {
		if base64.StdEncoding.EncodeToString(creds) == "Cg==" {
			fmt.Println("Credentials are zeroed")
		} else {
			fmt.Println(string(creds))
		}
	}
}

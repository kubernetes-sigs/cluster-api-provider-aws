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

package credentials

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/controller"
)

// UpdateCredentialsInput defines the specs for update credentials input.
type UpdateCredentialsInput struct {
	KubeconfigPath    string
	KubeconfigContext string
	Credentials       string
	Namespace         string
}

// UpdateCredentials updates the CAPA controller bootstrap secret
// RolloutControllers() must be called after any change to the controller bootstrap secret to take effect.
func UpdateCredentials(input UpdateCredentialsInput) error {
	client, err := controller.GetClient(input.KubeconfigPath, input.KubeconfigContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get client-go client for the cluster: %s\n", err.Error())
		return err
	}

	creds := input.Credentials
	if creds == "" {
		creds = "Cg=="
	}

	patch := fmt.Sprintf("{\"data\":{\"credentials\": \"%s\"}}", creds)
	_, err = client.CoreV1().Secrets(input.Namespace).Patch(
		context.TODO(),
		controller.BootstrapCredsSecret,
		types.MergePatchType,
		[]byte(patch),
		metav1.PatchOptions{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to patch bootstrap credentials secret: %s\n", err.Error())
		return err
	}

	secret, err := client.CoreV1().Secrets(input.Namespace).Get(context.TODO(), controller.BootstrapCredsSecret, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get bootstrap credentials secret: %s\n", err.Error())
		return err
	}
	controller.PrintBootstrapCredentials(secret)
	return nil
}

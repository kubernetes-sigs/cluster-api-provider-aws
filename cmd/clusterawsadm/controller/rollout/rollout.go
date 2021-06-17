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

package rollout

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/controller"
)

// ControllerDeploymentName is a tag for capa controller manager.
const ControllerDeploymentName = "capa-controller-manager"

// RolloutControllersInput defines the specs for rollout controllers input.
type RolloutControllersInput struct {
	KubeconfigPath    string
	KubeconfigContext string
	Namespace         string
}

// RolloutControllers initiates rollout restrart on the CAPA controller deployment.
// Must be called after any change to the controller bootstrap secret.
func RolloutControllers(input RolloutControllersInput) error {
	client, err := controller.GetClient(input.KubeconfigPath, input.KubeconfigContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get client-go client for the cluster: %s\n", err.Error())
		return err
	}

	dp, err := client.AppsV1().Deployments(input.Namespace).Get(context.TODO(), ControllerDeploymentName, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get %s deployment: %s\n", ControllerDeploymentName, err.Error())
		return err
	}

	initialDeployMarshalled, err := json.Marshal(dp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal controller deployment: %s\n", err.Error())
		return err
	}

	updatedDeployment := dp.DeepCopy()
	updatedDeployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	updatedDeployMarshalled, err := json.Marshal(updatedDeployment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal updated controller deployment: %s\n", err.Error())
		return err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(initialDeployMarshalled, updatedDeployMarshalled, appsv1.Deployment{})
	if err != nil {
		panic(err)
	}

	_, err = client.AppsV1().Deployments(input.Namespace).Patch(
		context.TODO(),
		ControllerDeploymentName,
		types.StrategicMergePatchType,
		patch,
		metav1.PatchOptions{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initiate rollout restart on CAPA controller manager deployment: %s\n", err.Error())
		return err
	}
	return nil
}

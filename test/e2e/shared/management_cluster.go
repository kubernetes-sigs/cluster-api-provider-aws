/*
Copyright 2026 The Kubernetes Authors.

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

import (
	"errors"
	"os"
	"path/filepath"
)

const selfHostedManagementClusterKubeconfigFile = "self-hosted-management-cluster.kubeconfig"

func validateManagementClusterFlags(useExistingCluster, provision, teardown bool) error {
	if useExistingCluster && provision {
		return errors.New("use-existing-cluster and provision-self-hosted-management-cluster cannot be enabled together")
	}
	if provision && teardown {
		return errors.New("provision-self-hosted-management-cluster and teardown-self-hosted-management-cluster cannot be enabled together")
	}
	return nil
}

func managementClusterLifecycleOnly(provision, teardown bool) bool {
	return provision || teardown
}

func removeFileIfExists(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func persistSelfHostedManagementClusterKubeconfig(sourcePath, artifactFolder string) (string, error) {
	managementKubeconfig, err := os.ReadFile(sourcePath) //nolint:gosec // The source path is created by the e2e framework.
	if err != nil {
		return "", err
	}
	managementKubeconfigPath := filepath.Join(artifactFolder, selfHostedManagementClusterKubeconfigFile)
	if err := os.WriteFile(managementKubeconfigPath, managementKubeconfig, 0o600); err != nil { //nolint:gosec // The artifact directory is supplied by the e2e runner.
		return "", err
	}
	return managementKubeconfigPath, nil
}

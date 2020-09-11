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

package eks

import (
	"fmt"
	"strings"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/hash"
)

const (
	// maxCharsName maximum number of characters for the name
	maxCharsName = 100

	clusterPrefix = "capa_"
)

// GenerateEKSName generates a name of an EKS cluster or nodegroup
func GenerateEKSName(clusterName, namespace string) (string, error) {
	escapedName := strings.Replace(clusterName, ".", "_", -1)
	eksName := fmt.Sprintf("%s_%s", namespace, escapedName)

	if len(eksName) < maxCharsName {
		return eksName, nil
	}

	hashLength := 32 - len(clusterPrefix)
	hashedName, err := hash.Base36TruncatedHash(eksName, hashLength)
	if err != nil {
		return "", fmt.Errorf("creating hash from cluster name: %w", err)
	}

	return fmt.Sprintf("%s%s", clusterPrefix, hashedName), nil
}

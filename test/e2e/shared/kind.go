//go:build e2e
// +build e2e

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
	"context"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	. "github.com/onsi/ginkgo/v2"
	"sigs.k8s.io/kind/pkg/build/nodeimage"
	"sigs.k8s.io/kind/pkg/cmd"

	"sigs.k8s.io/cluster-api/test/infrastructure/container"
	"sigs.k8s.io/cluster-api/test/infrastructure/kind"
)

// kindNodeImageForVersion returns the kind node image the Cluster API test framework uses when creating a kind
// cluster for the given Kubernetes version, and whether that image has to be built locally.
// Images pinned by digest in the framework's kind mapping are published by kind and are pulled as usual;
// any other version resolves to a plain kindest/node:<version> tag which may not have been published (yet).
func kindNodeImageForVersion(kubernetesVersion string) (image string, needsBuild bool, err error) {
	semVer, err := semver.ParseTolerant(kubernetesVersion)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse Kubernetes version %q: %w", kubernetesVersion, err)
	}
	image = kind.GetMapping(semVer, "").Image
	return image, !strings.Contains(image, "@sha256:"), nil
}

// ensureKindNodeImage returns the kind node image to use for the given Kubernetes version, making sure it is
// available locally by building it from the Kubernetes release artifacts (as `kind build node-image <version>`
// does) when kind has not published it.
func ensureKindNodeImage(ctx context.Context, kubernetesVersion string) (string, error) {
	image, needsBuild, err := kindNodeImageForVersion(kubernetesVersion)
	if err != nil {
		return "", err
	}
	if !needsBuild {
		By(fmt.Sprintf("Using published kind node image %s for Kubernetes %s", image, kubernetesVersion))
		return image, nil
	}

	containerRuntime, err := container.NewDockerClient()
	if err != nil {
		return "", fmt.Errorf("failed to get Docker runtime client: %w", err)
	}
	exists, err := containerRuntime.ImageExistsLocally(ctx, image)
	if err != nil {
		return "", fmt.Errorf("failed to check if image %s exists locally: %w", image, err)
	}
	if exists {
		By(fmt.Sprintf("Kind node image %s already exists locally", image))
		return image, nil
	}

	By(fmt.Sprintf("Building kind node image %s for Kubernetes %s", image, kubernetesVersion))
	if err := nodeimage.Build(
		nodeimage.WithImage(image),
		nodeimage.WithKubeParam(kubernetesVersion),
		nodeimage.WithLogger(cmd.NewLogger()),
	); err != nil {
		return "", fmt.Errorf("failed to build kind node image %s: %w", image, err)
	}
	return image, nil
}

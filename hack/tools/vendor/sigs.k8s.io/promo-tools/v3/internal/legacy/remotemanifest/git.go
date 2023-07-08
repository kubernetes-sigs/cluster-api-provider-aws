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

package remotemanifest

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry/schema"
)

const (
	gitCloneDepth = 1
)

// Git stores the Git-based information in order to fetch the Git contents (and
// to parse them into promoter manifests).
type Git struct {
	repoURL             *url.URL
	repoBranch          string
	thinManifestDirPath string
}

// Fetch gets the remote Git contents and parses it into promoter manifests. It
// could be the case that the repository is defined simply as a local path on
// disk (in the case of e2e tests where we do not have a full-fledged online
// repository for the manifests we want to audit) --- in such cases, we have to
// use the local path instead of freshly cloning a remote repo.
func (remote *Git) Fetch() ([]schema.Manifest, error) {
	// There is no remote; use the local path directly.
	if remote.repoURL.String() == "" {
		manifests, err := schema.ParseThinManifestsFromDir(
			remote.thinManifestDirPath, false)
		if err != nil {
			return nil, err
		}

		return manifests, nil
	}

	var repoPath string
	repoPath, err := cloneToTempDir(remote.repoURL, remote.repoBranch)
	if err != nil {
		return nil, err
	}

	manifests, err := schema.ParseThinManifestsFromDir(
		filepath.Join(repoPath, remote.thinManifestDirPath), false)
	if err != nil {
		return nil, err
	}

	// Garbage-collect freshly-cloned repo (we don't need it any more).
	err = os.RemoveAll(repoPath)
	if err != nil {
		// We don't really care too much about failures about removing the
		// (temporary) repoPath directory, because we'll clone a fresh one
		// anyway in the future. So don't return an error even if this fails.
		logrus.Errorf("Could not remove temporary Git repo %v: %v", repoPath, err)
	}

	return manifests, nil
}

func cloneToTempDir(
	repoURL fmt.Stringer,
	branch string,
) (string, error) {
	tdir, err := os.MkdirTemp("", "k8s.io-")
	if err != nil {
		return "", err
	}

	r, err := gogit.PlainClone(tdir, false, &gogit.CloneOptions{
		URL:           repoURL.String(),
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
		Depth:         gitCloneDepth,
	})
	if err != nil {
		return "", err
	}

	sha, err := getHeadSha(r)
	if err == nil {
		logrus.Infof("cloned %v at revision %v", tdir, sha)
	}

	return tdir, nil
}

func getHeadSha(repo *gogit.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

// NewGit creates a new Git implementation for defining a remote set of promoter
// schema.
func NewGit(
	repoURLStr string,
	repoBranch string,
	thinManifestDirPath string,
) (*Git, error) {
	remote := Git{}

	repoURL, err := url.Parse(repoURLStr)
	if err != nil {
		return nil, err
	}

	remote.repoURL = repoURL
	remote.repoBranch = repoBranch
	remote.thinManifestDirPath = thinManifestDirPath

	return &remote, nil
}

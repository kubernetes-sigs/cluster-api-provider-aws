/*
Copyright 2019 The Kubernetes Authors.

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

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// nlbReadinessScript returns a bash snippet that polls the NLB /readyz
// endpoint until it responds, with a 5-minute timeout. It is used by both
// the cloud-init and Ignition injection paths.
func nlbReadinessScript(host string, port int32) string {
	return fmt.Sprintf(
		`if ! command -v curl > /dev/null 2>&1; then echo "WARNING: curl not found, cannot perform API server NLB readiness check, continuing anyway"; else echo "Waiting for API server NLB to become ready at %s:%d..."; for i in $(seq 1 60); do if curl -sk --max-time 5 https://%s:%d/readyz > /dev/null 2>&1; then echo "API server NLB is ready"; break; fi; echo "Attempt $i/60: API server NLB not ready, retrying in 5s..."; sleep 5; done; if [ "$i" = 60 ]; then echo "WARNING: API server NLB readiness check timed out after 300s, proceeding anyway"; fi; fi`,
		host, port, host, port,
	)
}

// injectNLBReadinessCheck inserts a health check command into cloud-init
// userdata that polls the NLB endpoint before kubeadm join runs. This
// prevents worker join failures caused by NLB propagation delays.
// The check is non-fatal: it proceeds after a timeout so that the original
// kubeadm join behavior is preserved as a fallback.
func injectNLBReadinessCheck(userData []byte, host string, port int32) ([]byte, error) {
	if len(userData) == 0 {
		return nil, errors.New("empty userdata")
	}

	const kubeadmJoinMarker = "kubeadm join --config"

	idx := bytes.Index(userData, []byte(kubeadmJoinMarker))
	if idx == -1 {
		return nil, errors.New("could not find kubeadm join command in cloud-init userdata")
	}

	lineStart := bytes.LastIndex(userData[:idx], []byte("\n"))
	if lineStart == -1 {
		lineStart = 0
	} else {
		lineStart++
	}

	healthCheckCmd := fmt.Sprintf(
		"  - '/bin/bash -c ''%s'''\n",
		nlbReadinessScript(host, port),
	)

	var buf bytes.Buffer
	buf.Grow(len(userData) + len(healthCheckCmd))
	buf.Write(userData[:lineStart])
	buf.WriteString(healthCheckCmd)
	buf.Write(userData[lineStart:])

	return buf.Bytes(), nil
}

// injectNLBReadinessCheckIgnition inserts a health check into Ignition
// userdata by parsing the JSON, finding the /etc/kubeadm.sh file entry,
// decoding its data: URI, injecting the check before kubeadm join, and
// re-serializing.
func injectNLBReadinessCheckIgnition(userData []byte, host string, port int32) ([]byte, error) {
	if len(userData) == 0 {
		return nil, errors.New("empty userdata")
	}

	// Use a minimal struct for the parts we need to modify; json.Unmarshal
	// into a map preserves all other fields via re-marshal.
	var ignConfig map[string]json.RawMessage
	if err := json.Unmarshal(userData, &ignConfig); err != nil {
		return nil, errors.Wrap(err, "failed to parse Ignition JSON")
	}

	storageRaw, ok := ignConfig["storage"]
	if !ok {
		return nil, errors.New("Ignition config has no storage section")
	}

	var storage map[string]json.RawMessage
	if err := json.Unmarshal(storageRaw, &storage); err != nil {
		return nil, errors.Wrap(err, "failed to parse Ignition storage section")
	}

	filesRaw, ok := storage["files"]
	if !ok {
		return nil, errors.New("Ignition config has no files in storage")
	}

	var files []map[string]json.RawMessage
	if err := json.Unmarshal(filesRaw, &files); err != nil {
		return nil, errors.Wrap(err, "failed to parse Ignition files array")
	}

	const kubeadmShPath = "/etc/kubeadm.sh"
	const kubeadmJoinMarker = "kubeadm join"

	found := false
	for i, file := range files {
		var path string
		if pathRaw, ok := file["path"]; ok {
			if err := json.Unmarshal(pathRaw, &path); err != nil {
				continue
			}
		}
		if path != kubeadmShPath {
			continue
		}

		contentsRaw, ok := file["contents"]
		if !ok {
			return nil, errors.Errorf("Ignition file %s has no contents", kubeadmShPath)
		}

		var contents map[string]json.RawMessage
		if err := json.Unmarshal(contentsRaw, &contents); err != nil {
			return nil, errors.Wrapf(err, "failed to parse contents of %s", kubeadmShPath)
		}

		sourceRaw, ok := contents["source"]
		if !ok {
			return nil, errors.Errorf("Ignition file %s has no contents.source", kubeadmShPath)
		}

		var source string
		if err := json.Unmarshal(sourceRaw, &source); err != nil {
			return nil, errors.Wrapf(err, "failed to parse source of %s", kubeadmShPath)
		}

		// Decode the data: URI. Format is "data:,<url-encoded-content>"
		if !strings.HasPrefix(source, "data:,") {
			return nil, errors.Errorf("unexpected source format in %s: %q", kubeadmShPath, source[:min(len(source), 30)])
		}
		script, err := url.PathUnescape(source[len("data:,"):])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to URL-decode source of %s", kubeadmShPath)
		}

		joinIdx := strings.Index(script, kubeadmJoinMarker)
		if joinIdx == -1 {
			return nil, errors.Errorf("could not find %q in %s", kubeadmJoinMarker, kubeadmShPath)
		}

		lineStart := strings.LastIndex(script[:joinIdx], "\n")
		if lineStart == -1 {
			lineStart = 0
		} else {
			lineStart++
		}

		healthCheck := nlbReadinessScript(host, port) + "\n"
		script = script[:lineStart] + healthCheck + script[lineStart:]

		// Re-encode as a data: URI and write back.
		newSource := "data:," + url.PathEscape(script)
		newSourceJSON, err := json.Marshal(newSource)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal new source")
		}
		contents["source"] = newSourceJSON

		newContents, err := json.Marshal(contents)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal new contents")
		}
		files[i]["contents"] = newContents

		found = true
		break
	}

	if !found {
		return nil, errors.Errorf("could not find %s in Ignition config", kubeadmShPath)
	}

	newFilesJSON, err := json.Marshal(files)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal files array")
	}
	storage["files"] = newFilesJSON

	newStorageJSON, err := json.Marshal(storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal storage section")
	}
	ignConfig["storage"] = newStorageJSON

	return json.Marshal(ignConfig)
}

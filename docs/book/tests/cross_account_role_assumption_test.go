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

package tests

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	goruntime "runtime"
	"strings"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

const (
	crossAccountRoleAssumptionDoc = "../src/topics/eks/cross-account-role-assumption.md"
	managementRoleARN             = "arn:aws:iam::111122223333:role/controllers.cluster-api-provider-aws.sigs.k8s.io"
	workloadRoleARN               = "arn:aws:iam::444455556666:role/controllers.cluster-api-provider-aws.sigs.k8s.io"
	oidcProvider                  = "oidc.eks.us-east-1.amazonaws.com/id/EXAMPLE"
)

type iamStatement struct {
	Effect    string                    `json:"Effect"`
	Principal map[string]string         `json:"Principal"`
	Action    string                    `json:"Action"`
	Resource  string                    `json:"Resource"`
	Condition map[string]map[string]any `json:"Condition"`
}

type iamPolicy struct {
	Version   string         `json:"Version"`
	Statement []iamStatement `json:"Statement"`
}

func TestCrossAccountRoleAssumptionExamples(t *testing.T) {
	blocks := readFencedCodeBlocks(t)
	assertBlockCount(t, blocks, "json", 3)
	assertBlockCount(t, blocks, "yaml", 2)
	assertBlockCount(t, blocks, "bash", 3)

	t.Run("JSON examples", func(t *testing.T) {
		for i, block := range blocks["json"] {
			var document any
			if err := json.Unmarshal([]byte(block), &document); err != nil {
				t.Errorf("JSON block %d is invalid: %v", i+1, err)
			}
		}
	})

	t.Run("YAML examples", func(t *testing.T) {
		for i, block := range blocks["yaml"] {
			decoder := utilyaml.NewYAMLOrJSONDecoder(strings.NewReader(block), 4096)
			for document := 1; ; document++ {
				var value any
				if err := decoder.Decode(&value); err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					t.Errorf("YAML block %d, document %d is invalid: %v", i+1, document, err)
					break
				}
			}
		}
	})

	t.Run("shell instructions", func(t *testing.T) {
		runShellInstructions(t, blocks["bash"])
	})

	t.Run("IAM role chain", func(t *testing.T) {
		validateIAMRoleChain(t, blocks["json"])
	})

	t.Run("CAPA API contract", func(t *testing.T) {
		validateCAPAAPIContract(t, blocks["yaml"])
	})

	t.Run("provider manifest contract", func(t *testing.T) {
		manifest := readFileRelativeToTest(t, "../../../config/rbac/serviceaccount.yaml")
		expectedSubstitution := "${AWS_CONTROLLER_IAM_ROLE/#arn/eks.amazonaws.com/role-arn: arn}"
		if !bytes.Contains(manifest, []byte(expectedSubstitution)) {
			t.Fatalf("provider service account manifest does not consume AWS_CONTROLLER_IAM_ROLE using %q", expectedSubstitution)
		}
	})
}

func validateIAMRoleChain(t *testing.T, blocks []string) {
	t.Helper()

	var webIdentityTrust iamStatement
	if err := json.Unmarshal([]byte(blocks[0]), &webIdentityTrust); err != nil {
		t.Fatalf("failed to decode management role trust statement: %v", err)
	}
	assertEqual(t, "management trust action", webIdentityTrust.Action, "sts:AssumeRoleWithWebIdentity")
	assertEqual(t, "management trust federated principal", webIdentityTrust.Principal["Federated"],
		"arn:aws:iam::111122223333:oidc-provider/"+oidcProvider)
	stringEquals := webIdentityTrust.Condition["StringEquals"]
	assertEqual(t, "OIDC audience", stringEquals[oidcProvider+":aud"], "sts.amazonaws.com")
	assertEqual(t, "OIDC subjects", stringEquals[oidcProvider+":sub"], []any{
		"system:serviceaccount:capa-system:capa-controller-manager",
		"system:serviceaccount:capa-eks-control-plane-system:capa-eks-control-plane-controller-manager",
	})

	managementPolicy := decodeIAMPolicy(t, blocks[1], "management role permissions")
	assertEqual(t, "management permission action", managementPolicy.Statement[0].Action, "sts:AssumeRole")
	assertEqual(t, "management permission resource", managementPolicy.Statement[0].Resource, workloadRoleARN)

	workloadTrust := decodeIAMPolicy(t, blocks[2], "workload role trust")
	assertEqual(t, "workload trust action", workloadTrust.Statement[0].Action, "sts:AssumeRole")
	assertEqual(t, "workload trust principal", workloadTrust.Statement[0].Principal["AWS"], managementRoleARN)
}

func decodeIAMPolicy(t *testing.T, block, description string) iamPolicy {
	t.Helper()

	var policy iamPolicy
	if err := json.Unmarshal([]byte(block), &policy); err != nil {
		t.Fatalf("failed to decode %s: %v", description, err)
	}
	if policy.Version != "2012-10-17" {
		t.Fatalf("%s uses unexpected policy version %q", description, policy.Version)
	}
	if len(policy.Statement) != 1 {
		t.Fatalf("%s must contain exactly one statement, got %d", description, len(policy.Statement))
	}
	return policy
}

func validateCAPAAPIContract(t *testing.T, blocks []string) {
	t.Helper()

	scheme := runtime.NewScheme()
	if err := infrav1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to register infrastructure API: %v", err)
	}
	if err := ekscontrolplanev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to register EKS control plane API: %v", err)
	}
	decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()

	var controllerIdentity *infrav1.AWSClusterControllerIdentity
	var roleIdentity *infrav1.AWSClusterRoleIdentity
	var controlPlane *ekscontrolplanev1.AWSManagedControlPlane
	for blockIndex, block := range blocks {
		for documentIndex, document := range strings.Split(block, "\n---\n") {
			jsonDocument, err := utilyaml.ToJSON([]byte(document))
			if err != nil {
				t.Fatalf("failed to convert YAML block %d, document %d to JSON: %v", blockIndex+1, documentIndex+1, err)
			}
			object, _, err := decoder.Decode(jsonDocument, nil, nil)
			if err != nil {
				t.Fatalf("YAML block %d, document %d does not decode as a CAPA API object: %v", blockIndex+1, documentIndex+1, err)
			}
			switch typed := object.(type) {
			case *infrav1.AWSClusterControllerIdentity:
				controllerIdentity = typed
			case *infrav1.AWSClusterRoleIdentity:
				roleIdentity = typed
			case *ekscontrolplanev1.AWSManagedControlPlane:
				controlPlane = typed
			default:
				t.Fatalf("documentation contains unexpected CAPA object type %T", object)
			}
		}
	}

	if controllerIdentity == nil || roleIdentity == nil || controlPlane == nil {
		t.Fatalf("expected controller identity, role identity, and managed control plane examples")
	}
	assertEqual(t, "controller identity name", controllerIdentity.Name, "default")
	assertEqual(t, "controller allowed namespaces", controllerIdentity.Spec.AllowedNamespaces.NamespaceList, []string{"default"})
	assertEqual(t, "role identity ARN", roleIdentity.Spec.RoleArn, workloadRoleARN)
	assertIdentityReference(t, "role source identity", roleIdentity.Spec.SourceIdentityRef,
		infrav1.ControllerIdentityKind, controllerIdentity.Name)
	assertIdentityReference(t, "managed control plane identity", controlPlane.Spec.IdentityRef,
		infrav1.ClusterRoleIdentityKind, roleIdentity.Name)
}

func assertIdentityReference(t *testing.T, description string, reference *infrav1.AWSIdentityReference, kind infrav1.AWSIdentityKind, name string) {
	t.Helper()

	if reference == nil {
		t.Fatalf("%s is missing", description)
	}
	assertEqual(t, description+" kind", reference.Kind, kind)
	assertEqual(t, description+" name", reference.Name, name)
}

func assertEqual(t *testing.T, description string, actual, expected any) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("%s: got %#v, want %#v", description, actual, expected)
	}
}

func assertBlockCount(t *testing.T, blocks map[string][]string, language string, expected int) {
	t.Helper()

	if actual := len(blocks[language]); actual != expected {
		t.Fatalf("expected %d %s blocks, got %d", expected, language, actual)
	}
}

func readFencedCodeBlocks(t *testing.T) map[string][]string {
	t.Helper()

	content := readFileRelativeToTest(t, crossAccountRoleAssumptionDoc)

	blocks := map[string][]string{}
	var language string
	var block strings.Builder
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if language == "" {
			if candidate, found := strings.CutPrefix(line, "```"); found {
				language = strings.TrimSpace(candidate)
				block.Reset()
			}
			continue
		}

		if line == "```" {
			blocks[language] = append(blocks[language], block.String())
			language = ""
			continue
		}
		block.WriteString(line)
		block.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("failed to scan documentation: %v", err)
	}
	if language != "" {
		t.Fatalf("unterminated %q code block", language)
	}

	return blocks
}

func readFileRelativeToTest(t *testing.T, relativePath string) []byte {
	t.Helper()

	_, sourceFile, _, ok := goruntime.Caller(0)
	if !ok {
		t.Fatal("failed to determine test source path")
	}

	content, err := os.ReadFile(filepath.Join(filepath.Dir(sourceFile), relativePath)) //nolint:gosec // Paths are repository-owned test inputs.
	if err != nil {
		t.Fatalf("failed to read %s: %v", relativePath, err)
	}
	return content
}

func runShellInstructions(t *testing.T, blocks []string) {
	t.Helper()

	binDir := t.TempDir()
	writeExecutable(t, binDir, "clusterctl", `#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
[[ "$*" == "init --infrastructure aws" ]]
[[ "${AWS_CONTROLLER_IAM_ROLE}" == "arn:aws:iam::111122223333:role/controllers.cluster-api-provider-aws.sigs.k8s.io" ]]
`)
	writeExecutable(t, binDir, "kubectl", `#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
[[ "$*" == "-n capa-system get serviceaccount capa-controller-manager -o jsonpath={.metadata.annotations.eks\\.amazonaws\\.com/role-arn}{\"\\n\"}" ]]
printf '%s\n' "${AWS_CONTROLLER_IAM_ROLE}"
`)

	script := strings.Join(blocks, "\n")
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	// The script comes from the repository-owned documentation under test.
	command := exec.CommandContext(ctx, "bash", "-euo", "pipefail", "-c", script) //nolint:gosec
	command.Env = append(os.Environ(), fmt.Sprintf("PATH=%s%c%s", binDir, os.PathListSeparator, os.Getenv("PATH")))
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("documented shell instructions failed: %v\n%s", err, output)
	}
}

func writeExecutable(t *testing.T, directory, name, content string) {
	t.Helper()

	// The temporary stub must be executable and is only accessible to the current user.
	if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o700); err != nil { //nolint:gosec
		t.Fatalf("failed to create %s stub: %v", name, err)
	}
}

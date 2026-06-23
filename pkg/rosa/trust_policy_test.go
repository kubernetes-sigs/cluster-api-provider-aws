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

package rosa

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestInjectExternalIDIntoTrustPolicy(t *testing.T) {
	t.Run("injects external ID into policy with no condition", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "arn:aws:iam::123456789012:root"},
				"Action": "sts:AssumeRole"
			}]
		}`

		result, changed, err := InjectExternalIDIntoTrustPolicy(policy, "my-external-id")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeTrue())

		var doc map[string]interface{}
		g.Expect(json.Unmarshal([]byte(result), &doc)).To(Succeed())

		statements := doc["Statement"].([]interface{})
		stmt := statements[0].(map[string]interface{})
		condition := stmt["Condition"].(map[string]interface{})
		stringEquals := condition["StringEquals"].(map[string]interface{})
		g.Expect(stringEquals["sts:ExternalId"]).To(Equal("my-external-id"))
	})

	t.Run("injects external ID into policy with existing condition but no StringEquals", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "arn:aws:iam::123456789012:root"},
				"Action": "sts:AssumeRole",
				"Condition": {
					"Bool": {"aws:MultiFactorAuthPresent": "true"}
				}
			}]
		}`

		result, changed, err := InjectExternalIDIntoTrustPolicy(policy, "ext-123")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeTrue())

		var doc map[string]interface{}
		g.Expect(json.Unmarshal([]byte(result), &doc)).To(Succeed())

		statements := doc["Statement"].([]interface{})
		stmt := statements[0].(map[string]interface{})
		condition := stmt["Condition"].(map[string]interface{})
		g.Expect(condition["Bool"]).ToNot(BeNil())
		stringEquals := condition["StringEquals"].(map[string]interface{})
		g.Expect(stringEquals["sts:ExternalId"]).To(Equal("ext-123"))
	})

	t.Run("injects external ID into policy with existing StringEquals condition", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "arn:aws:iam::123456789012:root"},
				"Action": "sts:AssumeRole",
				"Condition": {
					"StringEquals": {"aws:PrincipalTag/Environment": "production"}
				}
			}]
		}`

		result, changed, err := InjectExternalIDIntoTrustPolicy(policy, "ext-456")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeTrue())

		var doc map[string]interface{}
		g.Expect(json.Unmarshal([]byte(result), &doc)).To(Succeed())

		statements := doc["Statement"].([]interface{})
		stmt := statements[0].(map[string]interface{})
		condition := stmt["Condition"].(map[string]interface{})
		stringEquals := condition["StringEquals"].(map[string]interface{})
		g.Expect(stringEquals["aws:PrincipalTag/Environment"]).To(Equal("production"))
		g.Expect(stringEquals["sts:ExternalId"]).To(Equal("ext-456"))
	})

	t.Run("skips injection when correct external ID already present", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "arn:aws:iam::123456789012:root"},
				"Action": "sts:AssumeRole",
				"Condition": {
					"StringEquals": {"sts:ExternalId": "my-id"}
				}
			}]
		}`

		_, changed, err := InjectExternalIDIntoTrustPolicy(policy, "my-id")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeFalse())
	})

	t.Run("updates external ID when different value present", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "arn:aws:iam::123456789012:root"},
				"Action": "sts:AssumeRole",
				"Condition": {
					"StringEquals": {"sts:ExternalId": "old-id"}
				}
			}]
		}`

		result, changed, err := InjectExternalIDIntoTrustPolicy(policy, "new-id")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeTrue())

		var doc map[string]interface{}
		g.Expect(json.Unmarshal([]byte(result), &doc)).To(Succeed())

		statements := doc["Statement"].([]interface{})
		stmt := statements[0].(map[string]interface{})
		condition := stmt["Condition"].(map[string]interface{})
		stringEquals := condition["StringEquals"].(map[string]interface{})
		g.Expect(stringEquals["sts:ExternalId"]).To(Equal("new-id"))
	})

	t.Run("does not modify non-AssumeRole statements", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": "arn:aws:iam::123456789012:root"},
					"Action": "sts:AssumeRole"
				},
				{
					"Effect": "Allow",
					"Principal": {"Service": "lambda.amazonaws.com"},
					"Action": "sts:AssumeRoleWithWebIdentity"
				}
			]
		}`

		result, changed, err := InjectExternalIDIntoTrustPolicy(policy, "ext-id")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeTrue())

		var doc map[string]interface{}
		g.Expect(json.Unmarshal([]byte(result), &doc)).To(Succeed())

		statements := doc["Statement"].([]interface{})
		stmt0 := statements[0].(map[string]interface{})
		g.Expect(stmt0["Condition"]).ToNot(BeNil())
		stmt1 := statements[1].(map[string]interface{})
		g.Expect(stmt1["Condition"]).To(BeNil())
	})

	t.Run("handles Action as array containing sts:AssumeRole", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "arn:aws:iam::123456789012:root"},
				"Action": ["sts:AssumeRole", "sts:TagSession"]
			}]
		}`

		result, changed, err := InjectExternalIDIntoTrustPolicy(policy, "array-ext-id")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeTrue())

		var doc map[string]interface{}
		g.Expect(json.Unmarshal([]byte(result), &doc)).To(Succeed())

		statements := doc["Statement"].([]interface{})
		stmt := statements[0].(map[string]interface{})
		condition := stmt["Condition"].(map[string]interface{})
		stringEquals := condition["StringEquals"].(map[string]interface{})
		g.Expect(stringEquals["sts:ExternalId"]).To(Equal("array-ext-id"))
	})

	t.Run("no change when no AssumeRole statements exist", func(t *testing.T) {
		g := NewWithT(t)

		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"Service": "lambda.amazonaws.com"},
				"Action": "sts:AssumeRoleWithWebIdentity"
			}]
		}`

		_, changed, err := InjectExternalIDIntoTrustPolicy(policy, "ext-id")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(changed).To(BeFalse())
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		g := NewWithT(t)

		_, _, err := InjectExternalIDIntoTrustPolicy("not json", "ext-id")
		g.Expect(err).To(HaveOccurred())
	})
}

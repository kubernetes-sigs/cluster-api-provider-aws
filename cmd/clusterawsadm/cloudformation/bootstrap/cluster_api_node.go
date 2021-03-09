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

package bootstrap

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/iam/v1alpha1"
)

func (t Template) secretPolicy(secureSecretsBackend infrav1.SecretBackend) iamv1.StatementEntry {
	switch secureSecretsBackend {
	case infrav1.SecretBackendSecretsManager:
		return iamv1.StatementEntry{
			Effect: iamv1.EffectAllow,
			Resource: iamv1.Resources{
				"arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*",
			},
			Action: iamv1.Actions{
				"secretsmanager:DeleteSecret",
				"secretsmanager:GetSecretValue",
			},
		}
	case infrav1.SecretBackendSSMParameterStore:
		return iamv1.StatementEntry{
			Effect: iamv1.EffectAllow,
			Resource: iamv1.Resources{
				"arn:*:ssm:*:*:parameter/cluster.x-k8s.io/*",
			},
			Action: iamv1.Actions{
				"ssm:DeleteParameter",
				"ssm:GetParameter",
			},
		}
	}
	return iamv1.StatementEntry{}
}

func (t Template) sessionManagerPolicy() iamv1.StatementEntry {
	return iamv1.StatementEntry{
		Effect:   iamv1.EffectAllow,
		Resource: iamv1.Resources{iamv1.Any},
		Action: iamv1.Actions{
			"ssm:UpdateInstanceInformation",
			"ssmmessages:CreateControlChannel",
			"ssmmessages:CreateDataChannel",
			"ssmmessages:OpenControlChannel",
			"ssmmessages:OpenDataChannel",
			"s3:GetEncryptionConfiguration",
		},
	}
}

func (t Template) nodeManagedPolicies() []string {
	policies := t.Spec.Nodes.ExtraPolicyAttachments

	if t.Spec.EKS.Enable {
		policies = append(policies,
			t.generateAWSManagedPolicyARN("AmazonEKSWorkerNodePolicy"),
			t.generateAWSManagedPolicyARN("AmazonEKS_CNI_Policy"),
		)
	}

	if t.Spec.Nodes.EC2ContainerRegistryReadOnly {
		policies = append(policies, t.generateAWSManagedPolicyARN("AmazonEC2ContainerRegistryReadOnly"))
	}

	return policies
}

func (t Template) nodePolicy() *iamv1.PolicyDocument {
	policyDocument := t.cloudProviderNodeAwsPolicy()
	for _, secureSecretsBackend := range t.Spec.SecureSecretsBackends {
		policyDocument.Statement = append(
			policyDocument.Statement,
			t.secretPolicy(secureSecretsBackend),
		)
	}
	policyDocument.Statement = append(
		policyDocument.Statement,
		t.sessionManagerPolicy(),
	)

	return policyDocument
}

func (t Template) generateAWSManagedPolicyARN(name string) string {
	return "arn:" + t.Spec.Partition + ":iam::aws:policy/" + name
}

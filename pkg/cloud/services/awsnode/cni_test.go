/*
Copyright 2022 The Kubernetes Authors.

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

package awsnode

import (
	"context"
	"testing"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/apis/crd/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

func TestReconcileCniVpcCniValues(t *testing.T) {
	tests := []struct {
		name       string
		cniValues  ekscontrolplanev1.VpcCni
		daemonSet  *v1.DaemonSet
		consistsOf []corev1.EnvVar
	}{
		{
			name: "users can set environment values",
			cniValues: ekscontrolplanev1.VpcCni{
				Env: []corev1.EnvVar{
					{
						Name:  "NAME1",
						Value: "VALUE1",
					},
				},
			},
			daemonSet: &v1.DaemonSet{
				TypeMeta: metav1.TypeMeta{
					Kind: "DaemonSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      awsNodeName,
					Namespace: awsNodeNamespace,
				},
				Spec: v1.DaemonSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name: awsNodeName,
									Env:  []corev1.EnvVar{},
								},
							},
						},
					},
				},
			},
			consistsOf: []corev1.EnvVar{
				{
					Name:  "NAME1",
					Value: "VALUE1",
				},
			},
		},
		{
			name: "users can set environment values without duplications",
			cniValues: ekscontrolplanev1.VpcCni{
				Env: []corev1.EnvVar{
					{
						Name:  "NAME1",
						Value: "VALUE1",
					},
					{
						Name:  "NAME1",
						Value: "VALUE2",
					},
				},
			},
			daemonSet: &v1.DaemonSet{
				TypeMeta: metav1.TypeMeta{
					Kind: "DaemonSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      awsNodeName,
					Namespace: awsNodeNamespace,
				},
				Spec: v1.DaemonSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name: awsNodeName,
									Env:  []corev1.EnvVar{},
								},
							},
						},
					},
				},
			},
			consistsOf: []corev1.EnvVar{
				{
					Name:  "NAME1",
					Value: "VALUE2",
				},
			},
		},
		{
			name: "users can set environment values overwriting existing values",
			cniValues: ekscontrolplanev1.VpcCni{
				Env: []corev1.EnvVar{
					{
						Name:  "NAME1",
						Value: "VALUE1",
					},
					{
						Name:  "NAME2",
						Value: "VALUE2",
					},
				},
			},
			daemonSet: &v1.DaemonSet{
				TypeMeta: metav1.TypeMeta{
					Kind: "DaemonSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      awsNodeName,
					Namespace: awsNodeNamespace,
				},
				Spec: v1.DaemonSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name: awsNodeName,
									Env: []corev1.EnvVar{
										{
											Name:  "NAME1",
											Value: "OVERWRITE",
										},
										{
											Name:  "NAME3",
											Value: "VALUE3",
										},
									},
								},
							},
						},
					},
				},
			},
			consistsOf: []corev1.EnvVar{
				{
					Name:  "NAME1",
					Value: "VALUE1",
				},
				{
					Name:  "NAME2",
					Value: "VALUE2",
				},
				{
					Name:  "NAME3",
					Value: "VALUE3",
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name+" without secondary cidr", func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			g := NewWithT(t)
			mockClient := &cachingClient{
				getValue: tc.daemonSet,
			}
			m := &mockScope{
				client: mockClient,
				cni:    tc.cniValues,
			}
			s := NewService(m)

			err := s.ReconcileCNI(context.Background())
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(mockClient.updateChain).NotTo(BeEmpty())
			ds, ok := mockClient.updateChain[0].(*v1.DaemonSet)
			g.Expect(ok).To(BeTrue())
			g.Expect(ds.Spec.Template.Spec.Containers).NotTo(BeEmpty())
			g.Expect(ds.Spec.Template.Spec.Containers[0].Env).To(ConsistOf(tc.consistsOf))
		})
	}

	for _, tc := range tests {
		t.Run(tc.name+" with secondary cidr", func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			g := NewWithT(t)
			mockClient := &cachingClient{
				getValue: tc.daemonSet,
			}
			m := &mockScope{
				client:             mockClient,
				cni:                tc.cniValues,
				secondaryCidrBlock: aws.String("100.0.0.1/20"),
				securityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					"node": {
						ID:   "subnet-1234",
						Name: "node",
					},
				},
				subnets: infrav1.Subnets{
					{
						// we aren't testing reconcileSubnets where this extra conf is added so putting it in by hand
						ID:        "subnet-1234",
						CidrBlock: "100.0.0.1/20",
						Tags: infrav1.Tags{
							infrav1.NameAWSSubnetAssociation: infrav1.SecondarySubnetTagValue,
						},
					},
				},
			}
			s := NewService(m)

			err := s.ReconcileCNI(context.Background())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(mockClient.updateChain).NotTo(BeEmpty()) // 0: eniconfig 1: daemonset
			eniconf, ok := mockClient.updateChain[0].(*v1alpha1.ENIConfig)
			g.Expect(ok).To(BeTrue())
			g.Expect(len(eniconf.Spec.SecurityGroups)).To(Equal(1))
			g.Expect(eniconf.Spec.SecurityGroups[0]).To(Equal(m.securityGroups["node"].ID))
			g.Expect(eniconf.Spec.Subnet).To(Equal(m.subnets[0].ID))

			ds, ok := mockClient.updateChain[1].(*v1.DaemonSet)
			g.Expect(ok).To(BeTrue())
			g.Expect(ds.Spec.Template.Spec.Containers).NotTo(BeEmpty())
			g.Expect(ds.Spec.Template.Spec.Containers[0].Env).To(ConsistOf(tc.consistsOf))
		})
	}
}

type cachingClient struct {
	client.Client
	getValue    client.Object
	updateChain []client.Object
}

func (c *cachingClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if _, ok := obj.(*v1.DaemonSet); ok {
		daemonset, _ := obj.(*v1.DaemonSet)
		*daemonset = *c.getValue.(*v1.DaemonSet)
	}
	return nil
}

func (c *cachingClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	c.updateChain = append(c.updateChain, obj)
	return nil
}

func (c *cachingClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return nil
}

type mockScope struct {
	scope.AWSNodeScope
	client             client.Client
	cni                ekscontrolplanev1.VpcCni
	secondaryCidrBlock *string
	securityGroups     map[infrav1.SecurityGroupRole]infrav1.SecurityGroup
	subnets            infrav1.Subnets
}

func (s *mockScope) RemoteClient() (client.Client, error) {
	return s.client, nil
}

func (s *mockScope) VpcCni() ekscontrolplanev1.VpcCni {
	return s.cni
}

func (s *mockScope) Info(msg string, keysAndValues ...interface{}) {

}

func (s *mockScope) Name() string {
	return "mock-name"
}

func (s *mockScope) Namespace() string {
	return "mock-namespace"
}

func (s *mockScope) DisableVPCCNI() bool {
	return false
}

func (s *mockScope) SecondaryCidrBlock() *string {
	return s.secondaryCidrBlock
}

func (s *mockScope) SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup {
	return s.securityGroups
}

func (s *mockScope) Subnets() infrav1.Subnets {
	return s.subnets
}

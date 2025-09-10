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

// Package gc provides a way to handle AWS garbage collection on deletion.
package gc

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec" // import all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // import all oidc plugins
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/annotations"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/cluster-api/util/patch"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clusterv1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
}

// CmdProcessor handles the garbage collection commands.
type CmdProcessor struct {
	client client.Client

	clusterName string
	namespace   string
}

// GCInput holds the configuration for the command processor.
type GCInput struct {
	ClusterName    string
	Namespace      string
	KubeconfigPath string
}

// CmdProcessorOption is a function type to supply options when creating the command processor.
type CmdProcessorOption func(proc *CmdProcessor) error

// WithClient is an option that enable you to explicitly supply a client.
func WithClient(client client.Client) CmdProcessorOption {
	return func(proc *CmdProcessor) error {
		proc.client = client

		return nil
	}
}

// New creates a new instance of the command processor.
func New(input GCInput, opts ...CmdProcessorOption) (*CmdProcessor, error) {
	cmd := &CmdProcessor{
		clusterName: input.ClusterName,
		namespace:   input.Namespace,
	}

	for _, opt := range opts {
		if err := opt(cmd); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	if cmd.client == nil {
		config, err := clientcmd.BuildConfigFromFlags("", input.KubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("building client config: %w", err)
		}

		cl, err := client.New(config, client.Options{Scheme: scheme})
		if err != nil {
			return nil, fmt.Errorf("creating new client: %w", err)
		}

		cmd.client = cl
	}

	return cmd, nil
}

// Enable is used to enable external resource garbage collection for a cluster.
func (c *CmdProcessor) Enable(ctx context.Context) error {
	if err := c.setAnnotationAndPatch(ctx, infrav1.ExternalResourceGCAnnotation, "true"); err != nil {
		return fmt.Errorf("setting gc annotation to true: %w", err)
	}

	return nil
}

// Disable is used to disable external resource garbage collection for a cluster.
func (c *CmdProcessor) Disable(ctx context.Context) error {
	if err := c.setAnnotationAndPatch(ctx, infrav1.ExternalResourceGCAnnotation, "false"); err != nil {
		return fmt.Errorf("setting gc annotation to false: %w", err)
	}

	return nil
}

// Configure is used to configure external resource garbage collection for a cluster.
func (c *CmdProcessor) Configure(ctx context.Context, gcTasks []string) error {
	supportedGCTasks := []infrav1.GCTask{infrav1.GCTaskLoadBalancer, infrav1.GCTaskTargetGroup, infrav1.GCTaskSecurityGroup}

	for _, gcTask := range gcTasks {
		found := false

		for _, supportedGCTask := range supportedGCTasks {
			if gcTask == string(supportedGCTask) {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unsupported gc task: %s", gcTask)
		}
	}

	annotationValue := strings.Join(gcTasks, ",")

	if err := c.setAnnotationAndPatch(ctx, infrav1.ExternalResourceGCTasksAnnotation, annotationValue); err != nil {
		return fmt.Errorf("setting gc tasks annotation to %s: %w", annotationValue, err)
	}

	return nil
}

func (c *CmdProcessor) setAnnotationAndPatch(ctx context.Context, annotationName, annotationValue string) error {
	infraObj, err := c.getInfraCluster(ctx)
	if err != nil {
		return err
	}

	patchHelper, err := patch.NewHelper(infraObj, c.client)
	if err != nil {
		return fmt.Errorf("creating patch helper: %w", err)
	}

	if annotationValue != "" {
		annotations.Set(infraObj, annotationName, annotationValue)
	} else {
		annotations.Delete(infraObj, annotationName)
	}

	if err := patchHelper.Patch(ctx, infraObj); err != nil {
		return fmt.Errorf("patching infra cluster with gc annotation: %w", err)
	}

	return nil
}

func (c *CmdProcessor) getInfraCluster(ctx context.Context) (*unstructured.Unstructured, error) {
	cluster := &clusterv1.Cluster{}

	key := client.ObjectKey{
		Name:      c.clusterName,
		Namespace: c.namespace,
	}

	if err := c.client.Get(ctx, key, cluster); err != nil {
		return nil, fmt.Errorf("getting capi cluster %s/%s: %w", c.namespace, c.clusterName, err)
	}

	ref := cluster.Spec.InfrastructureRef
	obj, err := external.Get(ctx, c.client, ref)
	if err != nil {
		return nil, fmt.Errorf("getting infra cluster %s/%s: %w", ref.Namespace, ref.Name, err)
	}

	return obj, nil
}

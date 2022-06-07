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
package controllers

import (
	"context"

  awsmanagedcontrolplanecontrollers "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/internal/controllers"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

type AWSManagedControlPlaneReconciler struct {
	client.Client
	Recorder  record.EventRecorder
	Endpoints []scope.ServiceEndpoint

	EnableIAM            bool
	AllowAdditionalRoles bool
	WatchFilterValue     string
}

func (r *AWSManagedControlPlaneReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return (&awsmanagedcontrolplanecontrollers.AWSManagedControlPlaneReconciler{
    Client:               r.Client,
		Recorder:             r.Recorder,
		Endpoints:            r.Endpoints,
		EnableIAM:            r.EnableIAM,
		AllowAdditionalRoles: r.AllowAdditionalRoles,
		WatchFilterValue:     r.WatchFilterValue,
	}).SetupWithManager(ctx, mgr, options)
}

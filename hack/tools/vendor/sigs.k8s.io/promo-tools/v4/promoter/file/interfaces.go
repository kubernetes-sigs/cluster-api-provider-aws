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

package file

import (
	"context"

	api "sigs.k8s.io/promo-tools/v4/api/files"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// SyncFileOp defines a synchronization operation.
//
//counterfeiter:generate . SyncFileOp
type SyncFileOp interface {
	Run(ctx context.Context) error
}

// Provider defines a file provider, able to work with GCS or S3.
type Provider interface {
	// Scheme returns the URI scheme we handle.
	Scheme() string

	// OpenFilestore opens a handle to the specified filestore.
	OpenFilestore(ctx context.Context, filestore *api.Filestore, useServiceAccount, confirm bool) (syncFilestore, error)
}

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

package audit

import (
	reg "sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry/registry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/logclient"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/remotemanifest"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/report"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/stream"
)

// GcrReadingFacility holds functions used to create streams for reading the
// repository and manifest list.
type GcrReadingFacility struct {
	ReadRepo         func(*reg.SyncContext, registry.Context) stream.Producer
	ReadManifestList func(*reg.SyncContext, *reg.GCRManifestListContext) stream.Producer
}

// ServerContext holds all of the initialization data for the server to start
// up.
type ServerContext struct {
	ID                     string
	RemoteManifestFacility remotemanifest.Facility
	ErrorReportingFacility report.ReportingFacility
	LoggingFacility        logclient.LoggingFacility
	GcrReadingFacility     GcrReadingFacility
}

// PubSubMessageInner is the inner struct that holds the actual Pub/Sub
// information.
type PubSubMessageInner struct {
	Data []byte `json:"data,omitempty"`
	ID   string `json:"id"`
}

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Message      PubSubMessageInner `json:"message"`
	Subscription string             `json:"subscription"`
}

const (
	// LogName is the auditing log name to use. This is the name that comes up
	// for "gcloud logging logs list".
	LogName = "cip-audit-log"
)

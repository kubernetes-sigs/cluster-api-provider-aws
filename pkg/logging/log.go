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

// Package logging defines the logger to use within the cluster-api-provider-aws
// package. Swap the implementation here.
package logging

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/logging/klog"
)

// Log is the log used in this project.
type Log struct {
	klog.Logger
}

var _ logr.Logger = &Log{}
var _ logr.InfoLogger = &Log{}

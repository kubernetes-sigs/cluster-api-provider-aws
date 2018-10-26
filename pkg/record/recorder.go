/*
Copyright 2018 The Kubernetes authors.

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

package record

import (
	"strings"
	"sync"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	clusterapiclientsetscheme "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/scheme"
)

const (
	defaultRecorderSource = "cluster-api-aws-provider"
)

var (
	initOnce        sync.Once
	defaultRecorder record.EventRecorder
)

func init() {
	defaultRecorder = new(record.FakeRecorder)
}

// Init initializes the global default recorder. It can only be called once. Subsequent calls are considered noops.
func Init(kubeClient *clientset.Clientset) {
	initOnce.Do(func() {
		scheme := runtime.NewScheme()
		if err := corev1.AddToScheme(scheme); err != nil {
			glog.Fatal(err)
		}
		clusterapiclientsetscheme.AddToScheme(scheme)

		broadcaster := record.NewBroadcaster()
		broadcaster.StartLogging(glog.Infof)
		broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.RESTClient()).Events("")})
		defaultRecorder = broadcaster.NewRecorder(scheme, corev1.EventSource{Component: defaultRecorderSource})
	})
}

// Event constructs an event from the given information and puts it in the queue for sending.
func Event(object runtime.Object, reason, message string) {
	defaultRecorder.Event(object, corev1.EventTypeNormal, strings.Title(reason), message)
}

// Eventf is just like Event, but with Sprintf for the message field.
func Eventf(object runtime.Object, reason, message string, args ...interface{}) {
	defaultRecorder.Eventf(object, corev1.EventTypeNormal, strings.Title(reason), message, args)
}

// Event constructs a warning event from the given information and puts it in the queue for sending.
func Warn(object runtime.Object, reason, message string) {
	defaultRecorder.Event(object, corev1.EventTypeWarning, strings.Title(reason), message)
}

// Eventf is just like Event, but with Sprintf for the message field.
func Warnf(object runtime.Object, reason, message string, args ...interface{}) {
	defaultRecorder.Eventf(object, corev1.EventTypeWarning, strings.Title(reason), message, args)
}

package events

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ref "k8s.io/client-go/tools/reference"
	clusterapiclientsetscheme "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/scheme"
)

// ObjectRec holds a contained event recorder
type ObjectRec struct {
	rec *record.EventRecorder
	obj runtime.Object
}

type key int

const (
	contextKey key = 0

	// Normal represents a normal event
	Normal = v1.EventTypeNormal

	// Warning event
	Warning = v1.EventTypeWarning

	// Failure event
	Failure = "Failure"

	// Pending event
	Pending = "Pending"
)

var (
	globalEventRec *record.EventRecorder
)

type logEvent struct {
	Object  *string `json:"object,omitempty"`
	Type    string  `json:"type"`
	Reason  string  `json:"reason"`
	Message string  `json:"message"`
}

// Eventf sends a formatted event to the event recorder sink
func (o *ObjectRec) Eventf(eventtype, reason, messageFmt string, args ...interface{}) {
	rec := *o.rec
	if o.obj == nil {
		o.Infof(eventtype, reason, messageFmt, args...)
	}
	rec.Eventf(o.obj, eventtype, reason, messageFmt, args...)
}

// Event sends an event to the event recorder sink
func (o *ObjectRec) Event(eventtype, reason, message string) {
	rec := *o.rec
	if o.obj == nil {
		o.Info(eventtype, reason, message)
	}
	rec.Event(o.obj, eventtype, reason, message)
}

func objectDescription(obj runtime.Object) string {
	ref, err := ref.GetReference(clusterapiclientsetscheme.Scheme, obj)
	if err != nil {
		glog.Error("Could not construct object reference")
	}

	return fmt.Sprintf("%v", ref.Name)
	//	gvk := obj.GetObjectKind().GroupVersionKind().GroupKind()
	//return fmt.Sprintf("%s/%s", gvk.Group, gvk.Kind)
}

func jsonLog(obj runtime.Object, eventtype, reason, message string) string {
	le := &logEvent{
		Type:    eventtype,
		Reason:  reason,
		Message: message,
	}
	if obj != nil {
		desc := objectDescription(obj)
		le.Object = &desc
	}
	str, _ := json.Marshal(le)
	return string(str)
}

// Info puts a structured log event on level 2
func (o *ObjectRec) Info(eventtype, reason, message string) {
	glog.InfoDepth(1, jsonLog(o.obj, eventtype, reason, message))
}

// Infof puts a structured log event on level 2
func (o *ObjectRec) Infof(eventtype, reason, messageFmt string, args ...interface{}) {
	glog.InfoDepth(1, jsonLog(o.obj, eventtype, reason, fmt.Sprintf(messageFmt, args...)))
}

// Errorf puts a structured log event on level 2
func (o *ObjectRec) Errorf(eventtype, reason, messageFmt string, args ...interface{}) {
	glog.ErrorDepth(1, jsonLog(o.obj, eventtype, reason, fmt.Sprintf(messageFmt, args...)))
}

// Error puts a structured log event on level 2
func (o *ObjectRec) Error(eventtype, reason, message string) {
	glog.ErrorDepth(1, jsonLog(o.obj, eventtype, reason, message))
}

// NewRecorderContext attaches an event recorder to a context
func NewStdObjRecorder(obj runtime.Object) (*ObjectRec, error) {
	rec, err := GetStdEventRecorder()
	if err != nil {
		return nil, err
	}
	o := &ObjectRec{
		rec: rec,
		obj: obj,
	}
	return o, nil
}

// GetStdEventRecorder returns the globally registered event recorder if available
func GetStdEventRecorder() (*record.EventRecorder, error) {
	if globalEventRec != nil {
		return globalEventRec, nil
	}
	return nil, errors.New("No object recorder found")
}

// SetStdEventRecorder sets a global, standard event recorder
func SetStdEventRecorder(rec *record.EventRecorder) {
	globalEventRec = rec
}

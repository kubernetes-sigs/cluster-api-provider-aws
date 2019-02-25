package resourcebuilder

import (
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"

	"github.com/openshift/cluster-version-operator/lib"
)

var (
	// Mapper is default ResourceMapper.
	Mapper = NewResourceMapper()
)

// ResourceMapper maps {Group, Version} to a function that returns Interface and an error.
type ResourceMapper struct {
	l *sync.Mutex

	gvkToNew map[schema.GroupVersionKind]NewInteraceFunc
}

// AddToMap adds all keys from caller to input.
// Locks the input ResourceMapper before adding the keys from caller.
func (rm *ResourceMapper) AddToMap(irm *ResourceMapper) {
	irm.l.Lock()
	defer irm.l.Unlock()
	for k, v := range rm.gvkToNew {
		irm.gvkToNew[k] = v
	}
}

// Exist returns true when gvk is known.
func (rm *ResourceMapper) Exists(gvk schema.GroupVersionKind) bool {
	_, ok := rm.gvkToNew[gvk]
	return ok
}

// RegisterGVK adds GVK to NewInteraceFunc mapping.
// It does not lock before adding the mapping.
func (rm *ResourceMapper) RegisterGVK(gvk schema.GroupVersionKind, f NewInteraceFunc) {
	rm.gvkToNew[gvk] = f
}

// NewResourceMapper returns a new map.
// This is required a we cannot push to uninitialized map.
func NewResourceMapper() *ResourceMapper {
	m := map[schema.GroupVersionKind]NewInteraceFunc{}
	return &ResourceMapper{
		l:        &sync.Mutex{},
		gvkToNew: m,
	}
}

type MetaV1ObjectModifierFunc func(metav1.Object)

// NewInteraceFunc returns an Interface.
// It requires rest Config that can be used to create a client
// and the Manifest.
type NewInteraceFunc func(rest *rest.Config, m lib.Manifest) Interface

type Interface interface {
	WithModifier(MetaV1ObjectModifierFunc) Interface
	Do() error
}

// New returns Interface using the mapping stored in mapper for m Manifest.
func New(mapper *ResourceMapper, rest *rest.Config, m lib.Manifest) (Interface, error) {
	f, ok := mapper.gvkToNew[m.GVK]
	if !ok {
		return nil, fmt.Errorf("No mapping found for gvk: %v", m.GVK)
	}
	return f(rest, m), nil
}

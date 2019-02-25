package dynamicclient

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/api/meta"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type resourceClientFactory struct {
	dynamicClient dynamic.Interface
	restMapper    *restmapper.DeferredDiscoveryRESTMapper
}

var (
	// this stores the singleton in a package local
	singletonFactory *resourceClientFactory
	once             sync.Once
)

// Private constructor for once.Do
func newSingletonFactory(config *rest.Config) func() {
	return func() {
		cachedDiscoveryClient := cached.NewMemCacheClient(kubernetes.NewForConfigOrDie(config).Discovery())
		restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
		restMapper.Reset()

		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		singletonFactory = &resourceClientFactory{
			dynamicClient: dynamicClient,
			restMapper:    restMapper,
		}
		singletonFactory.runBackgroundCacheReset(1 * time.Minute)
	}
}

// New returns the resource client using a singleton factory
func New(config *rest.Config, gvk schema.GroupVersionKind, namespace string) (dynamic.ResourceInterface, error) {
	once.Do(newSingletonFactory(config))
	return singletonFactory.getResourceClient(gvk, namespace)
}

// getResourceClient returns the dynamic client for the resource specified by the gvk.
func (c *resourceClientFactory) getResourceClient(gvk schema.GroupVersionKind, namespace string) (dynamic.ResourceInterface, error) {
	var (
		gvr        *schema.GroupVersionResource
		namespaced bool
		err        error
	)
	gvr, namespaced, err = gvkToGVR(gvk, c.restMapper)
	if meta.IsNoMatchError(err) {
		// refresh the restMapperCache and try once more.
		c.restMapper.Reset()
		gvr, namespaced, err = gvkToGVR(gvk, c.restMapper)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get resource type")
	}

	// sometimes manifests of non-namespaced resources
	// might have namespace set.
	// preventing such cases.
	ns := namespace
	if !namespaced {
		ns = ""
	}
	return c.dynamicClient.Resource(*gvr).Namespace(ns), nil
}

func gvkToGVR(gvk schema.GroupVersionKind, restMapper *restmapper.DeferredDiscoveryRESTMapper) (*schema.GroupVersionResource, bool, error) {
	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if meta.IsNoMatchError(err) {
		return nil, false, err
	}
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to get the resource REST mapping for GroupVersionKind(%s)", gvk.String())
	}

	return &mapping.Resource, mapping.Scope.Name() == meta.RESTScopeNameNamespace, nil
}

// runBackgroundCacheReset - Starts the rest mapper cache reseting
// at a duration given.
func (c *resourceClientFactory) runBackgroundCacheReset(duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			c.restMapper.Reset()
		}
	}()
}

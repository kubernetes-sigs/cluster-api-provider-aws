/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
	"flag"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/pflag"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	kubeclientset "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"sigs.k8s.io/controller-runtime/pkg/client/config"

	mapiv1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	mapiclient "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	mapiinformersfactory "github.com/openshift/cluster-api/pkg/client/informers_generated/externalversions"
	mapiinformers "github.com/openshift/cluster-api/pkg/client/informers_generated/externalversions/machine/v1beta1"
	lister "github.com/openshift/cluster-api/pkg/client/listers_generated/machine/v1beta1"

	kubeinformers "k8s.io/client-go/informers"
)

const (
	// maxRetries is the number of times a service will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the
	// sequence of delays between successive queuings of a service.
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries     = 15
	controllerName = "nodelink"

	// machineAnnotationKey is the annotation storing a link between a node and it's machine. Should match upstream cluster-api machine controller. (node.go)
	machineAnnotationKey = "machine.openshift.io/machine"
)

// NewController returns a new *Controller.
func NewController(
	nodeInformer coreinformers.NodeInformer,
	machineInformer mapiinformers.MachineInformer,
	kubeClient kubeclientset.Interface,
	capiClient mapiclient.Interface,
) *Controller {

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	// TODO: remove the wrapper when every clients have moved to use the clientset.
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.CoreV1().RESTClient()).Events("")})

	c := &Controller{
		capiClient: capiClient,
		kubeClient: kubeClient,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "nodelink"),
	}

	nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addNode,
		UpdateFunc: c.updateNode,
		DeleteFunc: c.deleteNode,
	})

	machineInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addMachine,
		UpdateFunc: c.updateMachine,
		DeleteFunc: c.deleteMachine,
	})

	c.nodeLister = nodeInformer.Lister()
	c.nodesSynced = nodeInformer.Informer().HasSynced

	c.machinesLister = machineInformer.Lister()
	c.machinesSynced = machineInformer.Informer().HasSynced

	c.syncHandler = c.syncNode
	c.enqueueNode = c.enqueue

	c.machineAddress = make(map[string]*mapiv1.Machine)

	return c
}

// Controller monitors nodes and links them to their machines when possible, as well as applies desired labels and taints.
type Controller struct {
	capiClient mapiclient.Interface
	kubeClient kubeclientset.Interface

	// To allow injection for testing.
	syncHandler func(hKey string) error

	// used for unit testing
	enqueueNode func(node *corev1.Node)

	nodeLister  corelister.NodeLister
	nodesSynced cache.InformerSynced

	machinesLister lister.MachineLister
	machinesSynced cache.InformerSynced

	// Machines that need to be synced
	queue workqueue.RateLimitingInterface

	// machines cache map[internalIP]machine
	machineAddress    map[string]*mapiv1.Machine
	machineAddressMux sync.Mutex
}

func (c *Controller) addNode(obj interface{}) {
	node := obj.(*corev1.Node)
	glog.V(3).Infof("Adding node: %q", node.Name)
	c.enqueueNode(node)
}

func (c *Controller) updateNode(old, cur interface{}) {
	//oldNode := old.(*corev1.Node)
	curNode := cur.(*corev1.Node)

	glog.V(3).Infof("Updating node: %q", curNode.Name)
	c.enqueueNode(curNode)
}

func (c *Controller) deleteNode(obj interface{}) {

	node, ok := obj.(*corev1.Node)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		node, ok = tombstone.Obj.(*corev1.Node)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Node %#v", obj))
			return
		}
	}

	glog.V(3).Infof("Deleting node")
	c.enqueueNode(node)
}

func (c *Controller) addMachine(obj interface{}) {
	machine := obj.(*mapiv1.Machine)

	c.machineAddressMux.Lock()
	defer c.machineAddressMux.Unlock()

	for _, a := range machine.Status.Addresses {
		// Use the internal IP to look for matches:
		if a.Type == corev1.NodeInternalIP {
			glog.V(3).Infof("Adding machine %q into machineAddress list for %q", machine.Name, a.Address)
			c.machineAddress[a.Address] = machine
			break
		}
	}
}

func (c *Controller) updateMachine(old, cur interface{}) {
	machine := cur.(*mapiv1.Machine)

	c.machineAddressMux.Lock()
	defer c.machineAddressMux.Unlock()

	for _, a := range machine.Status.Addresses {
		// Use the internal IP to look for matches:
		if a.Type == corev1.NodeInternalIP {
			c.machineAddress[a.Address] = machine
			glog.V(3).Infof("Updating machine addresses list. Machine: %q, address: %q", machine.Name, a.Address)
			break
		}
	}
}

func (c *Controller) deleteMachine(obj interface{}) {
	machine := obj.(*mapiv1.Machine)

	c.machineAddressMux.Lock()
	defer c.machineAddressMux.Unlock()

	for _, a := range machine.Status.Addresses {
		// Use the internal IP to look for matches:
		if a.Type == corev1.NodeInternalIP {
			delete(c.machineAddress, a.Address)
			break
		}
	}
	glog.V(3).Infof("Delete obsolete machines from machine addresses list")
}

// WaitForCacheSync is a wrapper around cache.WaitForCacheSync that generates log messages
// indicating that the controller identified by controllerName is waiting for syncs, followed by
// either a successful or failed sync.
func WaitForCacheSync(controllerName string, stopCh <-chan struct{}, cacheSyncs ...cache.InformerSynced) bool {
	glog.Infof("Waiting for caches to sync for %s controller", controllerName)

	if !cache.WaitForCacheSync(stopCh, cacheSyncs...) {
		utilruntime.HandleError(fmt.Errorf("unable to sync caches for %s controller", controllerName))
		return false
	}

	glog.Infof("Caches are synced for %s controller", controllerName)
	return true
}

// Run runs c; will not return until stopCh is closed. workers determines how
// many machines will be handled in parallel.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	glog.Infof("Starting nodelink controller")
	defer glog.Infof("Shutting down nodelink controller")

	if !WaitForCacheSync("machine", stopCh, c.machinesSynced) {
		return
	}

	if !WaitForCacheSync("node", stopCh, c.nodesSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

func (c *Controller) enqueue(node *corev1.Node) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(node)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", node, err))
		return
	}

	c.queue.Add(key)
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (c *Controller) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncHandler(key.(string))
	c.handleErr(err, key)

	return true
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < maxRetries {
		glog.Infof("Error syncing node %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	glog.Infof("Dropping node %q out of the queue: %v", key, err)
	c.queue.Forget(key)
}

// syncNode will sync the node with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (c *Controller) syncNode(key string) error {
	startTime := time.Now()
	glog.V(3).Infof("Syncing node")
	defer func() {
		glog.V(3).Infof("Finished syncing node, duration: %s", time.Now().Sub(startTime))
	}()

	_, _, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	node, err := c.nodeLister.Get(key)
	if errors.IsNotFound(err) {
		glog.Infof("Error syncing, Node %s has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}

	return c.processNode(node)
}
func (c *Controller) processNode(node *corev1.Node) error {
	machineKey, ok := node.Annotations[machineAnnotationKey]
	// No machine annotation, this is likely the first time we've seen the node,
	// need to load all machines and search for one with matching IP.
	var matchingMachine *mapiv1.Machine
	if ok {
		var err error
		namespace, machineName, err := cache.SplitMetaNamespaceKey(machineKey)
		if err != nil {
			glog.Infof("Error processing node %q. Machine annotation format is incorrect %q: %v", node.Name, machineKey, err)
			return err
		}
		matchingMachine, err = c.machinesLister.Machines(namespace).Get(machineName)
		// If machine matching annotation is not found, we'll still try to find one via IP matching:
		if err != nil {
			if errors.IsNotFound(err) {
				glog.Warningf("Machine %q associated to node %q has been deleted, will attempt to find new machine by IP", machineKey, node.Name)
			} else {
				return err
			}
		}
	}

	if matchingMachine == nil {
		// Find this nodes internal IP so we can search for a matching machine:
		var nodeInternalIP string
		for _, a := range node.Status.Addresses {
			if a.Type == corev1.NodeInternalIP {
				nodeInternalIP = a.Address
				break
			}
		}
		if nodeInternalIP == "" {
			glog.Warningf("Unable to find InternalIP for node %q", node.Name)
			return fmt.Errorf("unable to find InternalIP for node: %q", node.Name)
		}

		glog.V(4).Infof("Searching machine cache for IP match for node %q", node.Name)
		c.machineAddressMux.Lock()
		machine, found := c.machineAddress[nodeInternalIP]
		c.machineAddressMux.Unlock()

		if found {
			matchingMachine = machine
		}
	}

	if matchingMachine == nil {
		return fmt.Errorf("no machine was found for node: %q", node.Name)
	}

	glog.V(3).Infof("Found machine %s for node %s", machineKey, node.Name)
	modNode := node.DeepCopy()
	if modNode.Annotations == nil {
		modNode.Annotations = map[string]string{}
	}
	modNode.Annotations[machineAnnotationKey] = fmt.Sprintf("%s/%s", matchingMachine.Namespace, matchingMachine.Name)

	if modNode.Labels == nil {
		modNode.Labels = map[string]string{}
	}

	for k, v := range matchingMachine.Spec.Labels {
		glog.V(3).Infof("Copying label %s = %s", k, v)
		modNode.Labels[k] = v
	}

	addTaintsToNode(modNode, matchingMachine)

	if !reflect.DeepEqual(node, modNode) {
		glog.V(3).Infof("Node %q has changed, updating", modNode.Name)
		_, err := c.kubeClient.CoreV1().Nodes().Update(modNode)
		if err != nil {
			glog.Errorf("Error updating node: %v", err)
			return err
		}
	}
	return nil
}

// addTaintsToNode adds taints from machine object to the node object
// Taints are to be an authoritative list on the machine spec per cluster-api comments.
// However, we believe many components can directly taint a node and there is no direct source of truth that should enforce a single writer of taints
func addTaintsToNode(node *corev1.Node, machine *mapiv1.Machine) {
	for _, mTaint := range machine.Spec.Taints {
		glog.V(3).Infof("Adding taint %v from machine %q to node %q", mTaint, machine.Name, node.Name)
		alreadyPresent := false
		for _, nTaint := range node.Spec.Taints {
			if nTaint.Key == mTaint.Key && nTaint.Effect == mTaint.Effect {
				glog.V(3).Infof("Skipping to add machine taint, %v, to the node. Node already has a taint with same key and effect", mTaint)
				alreadyPresent = true
				break
			}
		}
		if !alreadyPresent {
			node.Spec.Taints = append(node.Spec.Taints, mTaint)
		}
	}
}

var (
	logLevel string
)

func main() {

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	config, err := config.GetConfig()
	if err != nil {
		glog.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create client for talking to the apiserver: %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create kubernetes client to talk to the apiserver: %v", err)
	}

	// TODO(jchaloup): set the resync period reasonably
	kubeSharedInformers := kubeinformers.NewSharedInformerFactory(kubeClient, 5*time.Second)
	mapiInformers := mapiinformersfactory.NewSharedInformerFactory(client, 5*time.Second)

	ctrl := NewController(
		kubeSharedInformers.Core().V1().Nodes(),
		mapiInformers.Machine().V1beta1().Machines(),
		kubeClient,
		client,
	)

	go ctrl.Run(1, wait.NeverStop)

	mapiInformers.Start(wait.NeverStop)
	kubeSharedInformers.Start(wait.NeverStop)

	select {}
}

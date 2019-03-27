package operator

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	osclientset "github.com/openshift/client-go/config/clientset/versioned"

	osconfigv1 "github.com/openshift/api/config/v1"
	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformersv1 "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	coreclientsetv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisterv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/scheme"
)

const (
	// maxRetries is the number of times a key will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// a machineconfig pool is going to be requeued:
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries        = 15
	ownedManifestsDir = "owned-manifests"
)

// Operator defines machine api operator.
type Operator struct {
	namespace, name string

	imagesFile string
	config     string

	kubeClient    kubernetes.Interface
	osClient      osclientset.Interface
	eventRecorder record.EventRecorder

	syncHandler func(ic string) error

	deployLister       appslisterv1.DeploymentLister
	deployListerSynced cache.InformerSynced

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue           workqueue.RateLimitingInterface
	operandVersions []osconfigv1.OperandVersion
}

// New returns a new machine config operator.
func New(
	namespace, name string,
	imagesFile string,

	config string,

	deployInformer appsinformersv1.DeploymentInformer,

	kubeClient kubernetes.Interface,
	osClient osclientset.Interface,
) *Operator {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&coreclientsetv1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	// we must report the version from the release payload when we report available at that level
	operandVersions := []osconfigv1.OperandVersion{}
	if releaseVersion := os.Getenv("RELEASE_VERSION"); len(releaseVersion) > 0 {
		operandVersions = append(operandVersions, osconfigv1.OperandVersion{Name: "operator", Version: releaseVersion})
	}

	optr := &Operator{
		namespace:       namespace,
		name:            name,
		imagesFile:      imagesFile,
		kubeClient:      kubeClient,
		osClient:        osClient,
		eventRecorder:   eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "machineapioperator"}),
		queue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "machineapioperator"),
		operandVersions: operandVersions,
	}

	deployInformer.Informer().AddEventHandler(optr.eventHandler())

	optr.config = config
	optr.syncHandler = optr.sync

	optr.deployLister = deployInformer.Lister()
	optr.deployListerSynced = deployInformer.Informer().HasSynced

	return optr
}

// Run runs the machine config operator.
func (optr *Operator) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer optr.queue.ShutDown()

	glog.Info("Starting Machine API Operator")
	defer glog.Info("Shutting down Machine API Operator")

	if !cache.WaitForCacheSync(stopCh,
		optr.deployListerSynced) {
		glog.Error("Failed to sync caches")
		return
	}
	glog.Info("Synced up caches")
	for i := 0; i < workers; i++ {
		go wait.Until(optr.worker, time.Second, stopCh)
	}

	<-stopCh
}

func (optr *Operator) eventHandler() cache.ResourceEventHandler {
	workQueueKey := fmt.Sprintf("%s/%s", optr.namespace, optr.name)
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { optr.queue.Add(workQueueKey) },
		UpdateFunc: func(old, new interface{}) { optr.queue.Add(workQueueKey) },
		DeleteFunc: func(obj interface{}) { optr.queue.Add(workQueueKey) },
	}
}

func (optr *Operator) worker() {
	for optr.processNextWorkItem() {
	}
}

func (optr *Operator) processNextWorkItem() bool {
	key, quit := optr.queue.Get()
	if quit {
		return false
	}
	defer optr.queue.Done(key)

	glog.V(4).Infof("Processing key %s", key)
	err := optr.syncHandler(key.(string))
	optr.handleErr(err, key)

	return true
}

func (optr *Operator) handleErr(err error, key interface{}) {
	if err == nil {
		optr.queue.Forget(key)
		return
	}

	if optr.queue.NumRequeues(key) < maxRetries {
		glog.V(1).Infof("Error syncing operator %v: %v", key, err)
		optr.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	glog.V(1).Infof("Dropping operator %q out of the queue: %v", key, err)
	optr.queue.Forget(key)
}

func (optr *Operator) sync(key string) error {
	startTime := time.Now()
	glog.V(4).Infof("Started syncing operator %q (%v)", key, startTime)
	defer func() {
		glog.V(4).Infof("Finished syncing operator %q (%v)", key, time.Since(startTime))
	}()

	operatorConfig, err := optr.maoConfigFromInstallConfig()
	if err != nil {
		glog.Errorf("Failed getting operator config: %v", err)
		return err
	}

	return optr.syncAll(*operatorConfig)
}

func (optr *Operator) maoConfigFromInstallConfig() (*OperatorConfig, error) {
	installConfig, err := getInstallConfig(optr.kubeClient)
	if err != nil {
		return nil, err
	}

	provider, err := getProviderFromInstallConfig(installConfig)
	if err != nil {
		return nil, err
	}

	images, err := getImagesFromJSONFile(optr.imagesFile)
	if err != nil {
		return nil, err
	}

	providerControllerImage, err := getProviderControllerFromImages(provider, *images)
	if err != nil {
		return nil, err
	}

	machineAPIOperatorImage, err := getMachineAPIOperatorFromImages(*images)
	if err != nil {
		return nil, err
	}

	return &OperatorConfig{
		optr.namespace,
		Controllers{
			providerControllerImage,
			machineAPIOperatorImage,
			machineAPIOperatorImage,
		},
	}, nil
}

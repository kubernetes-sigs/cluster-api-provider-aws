package cvo

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/golang/glog"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	coreclientsetv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	configv1 "github.com/openshift/api/config/v1"
	clientset "github.com/openshift/client-go/config/clientset/versioned"
	configinformersv1 "github.com/openshift/client-go/config/informers/externalversions/config/v1"
	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"

	"github.com/openshift/cluster-version-operator/lib"
	"github.com/openshift/cluster-version-operator/lib/resourceapply"
	"github.com/openshift/cluster-version-operator/lib/resourcebuilder"
	"github.com/openshift/cluster-version-operator/lib/resourcemerge"
	"github.com/openshift/cluster-version-operator/lib/validation"
	"github.com/openshift/cluster-version-operator/pkg/cvo/internal"
	"github.com/openshift/cluster-version-operator/pkg/cvo/internal/dynamicclient"
	"github.com/openshift/cluster-version-operator/pkg/payload"
)

const (
	// maxRetries is the number of times a machineconfig pool will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// a machineconfig pool is going to be requeued:
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15
)

// Operator defines cluster version operator. The CVO attempts to reconcile the appropriate image
// onto the cluster, writing status to the ClusterVersion object as it goes. A background loop
// periodically checks for new updates from a server described by spec.upstream and spec.channel.
// The main CVO sync loop is the single writer of ClusterVersion status.
//
// The CVO updates multiple conditions, but synthesizes them into a summary message on the
// Progressing condition to answer the question of "what version is available on the cluster".
// When errors occur, the Failing condition of the status is updated with a detailed message and
// reason, and then the reason is used to summarize the error onto the Progressing condition's
// message for a simple overview.
//
// The CVO periodically syncs the whole image to the cluster even if no version transition is
// detected in order to undo accidental actions.
//
// A release image is expected to contain a CVO binary, the manifests necessary to update the
// CVO, and the manifests of the other operators on the cluster. During an update the operator
// attempts to copy the contents of the image manifests into a temporary directory using a
// batch job and a shared host-path, then applies the CVO manifests using the image image
// for the CVO deployment. The deployment is then expected to launch the new process, and the
// new operator picks up the lease and applies the rest of the image.
type Operator struct {
	// nodename allows CVO to sync fetchPayload to same node as itself.
	nodename string
	// namespace and name are used to find the ClusterVersion, OperatorStatus.
	namespace, name string

	// releaseImage is the image the current operator points to and allows
	// templating of the CVO deployment manifest.
	releaseImage string
	// releaseVersion is a string identifier for the current version, read
	// from the image of the operator. It may be empty if no version exists, in
	// which case no available updates will be returned.
	releaseVersion string
	// releaseCreated, if set, is the timestamp of the current update.
	releaseCreated time.Time

	// restConfig is used to create resourcebuilder.
	restConfig *rest.Config

	client        clientset.Interface
	kubeClient    kubernetes.Interface
	eventRecorder record.EventRecorder

	// minimumUpdateCheckInterval is the minimum duration to check for updates from
	// the upstream.
	minimumUpdateCheckInterval time.Duration
	// payloadDir is intended for testing. If unset it will default to '/'
	payloadDir string
	// defaultUpstreamServer is intended for testing.
	defaultUpstreamServer string
	// syncBackoff allows the tests to use a quicker backoff
	syncBackoff wait.Backoff

	cvLister    configlistersv1.ClusterVersionLister
	coLister    configlistersv1.ClusterOperatorLister
	cacheSynced []cache.InformerSynced

	// queue tracks applying updates to a cluster.
	queue workqueue.RateLimitingInterface
	// availableUpdatesQueue tracks checking for updates from the update server.
	availableUpdatesQueue workqueue.RateLimitingInterface

	// statusLock guards access to modifying available updates
	statusLock       sync.Mutex
	availableUpdates *availableUpdates

	configSync ConfigSyncWorker
	// statusInterval is how often the configSync worker is allowed to retrigger
	// the main sync status loop.
	statusInterval time.Duration

	// lastAtLock guards access to controller memory about the sync loop
	lastAtLock          sync.Mutex
	lastResourceVersion int64
}

// New returns a new cluster version operator.
func New(
	nodename string,
	namespace, name string,
	releaseImage string,
	overridePayloadDir string,
	minimumInterval time.Duration,
	cvInformer configinformersv1.ClusterVersionInformer,
	coInformer configinformersv1.ClusterOperatorInformer,
	restConfig *rest.Config,
	client clientset.Interface,
	kubeClient kubernetes.Interface,
	enableMetrics bool,
) *Operator {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&coreclientsetv1.EventSinkImpl{Interface: kubeClient.CoreV1().Events(namespace)})

	optr := &Operator{
		nodename:     nodename,
		namespace:    namespace,
		name:         name,
		releaseImage: releaseImage,

		statusInterval:             15 * time.Second,
		minimumUpdateCheckInterval: minimumInterval,
		payloadDir:                 overridePayloadDir,
		defaultUpstreamServer:      "https://api.openshift.com/api/upgrades_info/v1/graph",

		restConfig:    restConfig,
		client:        client,
		kubeClient:    kubeClient,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: namespace}),

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "clusterversion"),
		availableUpdatesQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "availableupdates"),
	}

	optr.configSync = NewSyncWorker(
		optr.defaultPayloadRetriever(),
		NewResourceBuilder(optr.restConfig),
		minimumInterval,
		wait.Backoff{
			Duration: time.Second * 10,
			Factor:   1.3,
			Steps:    3,
		},
	)

	cvInformer.Informer().AddEventHandler(optr.eventHandler())

	optr.coLister = coInformer.Lister()
	optr.cacheSynced = append(optr.cacheSynced, coInformer.Informer().HasSynced)

	optr.cvLister = cvInformer.Lister()
	optr.cacheSynced = append(optr.cacheSynced, cvInformer.Informer().HasSynced)

	if enableMetrics {
		if err := optr.registerMetrics(coInformer.Informer()); err != nil {
			panic(err)
		}
	}

	if update, err := payload.LoadUpdate(optr.defaultPayloadDir(), releaseImage); err != nil {
		glog.Warningf("The local release contents are invalid - no current version can be determined from disk: %v", err)
	} else {
		optr.releaseCreated = update.ImageRef.CreationTimestamp.Time
		// XXX: set this to the cincinnati version in preference
		if _, err := semver.Parse(update.ImageRef.Name); err != nil {
			glog.Warningf("The local release contents name %q is not a valid semantic version - no current version will be reported: %v", update.ImageRef.Name, err)
		} else {
			optr.releaseVersion = update.ImageRef.Name
		}
	}

	return optr
}

// Run runs the cluster version operator until stopCh is completed. Workers is ignored for now.
func (optr *Operator) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer optr.queue.ShutDown()

	glog.Infof("Starting ClusterVersionOperator with minimum reconcile period %s", optr.minimumUpdateCheckInterval)
	defer glog.Info("Shutting down ClusterVersionOperator")

	if !cache.WaitForCacheSync(stopCh, optr.cacheSynced...) {
		glog.Info("Caches never synchronized")
		return
	}

	// trigger the first cluster version reconcile always
	optr.queue.Add(optr.queueKey())

	// start the config sync loop, and have it notify the queue when new status is detected
	go runThrottledStatusNotifier(stopCh, optr.statusInterval, 2, optr.configSync.StatusCh(), func() { optr.queue.Add(optr.queueKey()) })
	go optr.configSync.Start(8, stopCh)

	go wait.Until(func() { optr.worker(optr.queue, optr.sync) }, time.Second, stopCh)
	go wait.Until(func() { optr.worker(optr.availableUpdatesQueue, optr.availableUpdatesSync) }, time.Second, stopCh)

	<-stopCh
}

func (optr *Operator) queueKey() string {
	return fmt.Sprintf("%s/%s", optr.namespace, optr.name)
}

// eventHandler queues an update for the cluster version on any change to the given object.
// Callers should use this with a scoped informer.
func (optr *Operator) eventHandler() cache.ResourceEventHandler {
	workQueueKey := optr.queueKey()
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			optr.queue.Add(workQueueKey)
			optr.availableUpdatesQueue.Add(workQueueKey)
		},
		UpdateFunc: func(old, new interface{}) {
			optr.queue.Add(workQueueKey)
			optr.availableUpdatesQueue.Add(workQueueKey)
		},
		DeleteFunc: func(obj interface{}) {
			optr.queue.Add(workQueueKey)
		},
	}
}

func (optr *Operator) worker(queue workqueue.RateLimitingInterface, syncHandler func(string) error) {
	for processNextWorkItem(queue, syncHandler, optr.syncFailingStatus) {
	}
}

type syncFailingStatusFunc func(config *configv1.ClusterVersion, err error) error

func processNextWorkItem(queue workqueue.RateLimitingInterface, syncHandler func(string) error, syncFailingStatus syncFailingStatusFunc) bool {
	key, quit := queue.Get()
	if quit {
		return false
	}
	defer queue.Done(key)

	err := syncHandler(key.(string))
	handleErr(queue, err, key, syncFailingStatus)
	return true
}

func handleErr(queue workqueue.RateLimitingInterface, err error, key interface{}, syncFailingStatus syncFailingStatusFunc) {
	if err == nil {
		queue.Forget(key)
		return
	}

	if queue.NumRequeues(key) < maxRetries {
		glog.V(2).Infof("Error syncing operator %v: %v", key, err)
		queue.AddRateLimited(key)
		return
	}

	err = syncFailingStatus(nil, err)
	utilruntime.HandleError(err)
	glog.V(2).Infof("Dropping operator %q out of the queue %v: %v", key, queue, err)
	queue.Forget(key)
}

// sync ensures:
//
// 1. A ClusterVersion object exists
// 2. The ClusterVersion object has the appropriate status for the state of the cluster
// 3. The configSync object is kept up to date maintaining the user's desired version
//
// It returns an error if it could not update the cluster version object.
func (optr *Operator) sync(key string) error {
	startTime := time.Now()
	glog.V(4).Infof("Started syncing cluster version %q (%v)", key, startTime)
	defer func() {
		glog.V(4).Infof("Finished syncing cluster version %q (%v)", key, time.Since(startTime))
	}()

	// ensure the cluster version exists, that the object is valid, and that
	// all initial conditions are set.
	original, changed, err := optr.getOrCreateClusterVersion()
	if err != nil {
		return err
	}
	if changed {
		glog.V(4).Infof("Cluster version changed, waiting for newer event")
		return nil
	}

	// ensure that the object we do have is valid
	errs := validation.ValidateClusterVersion(original)
	// for fields that have meaning that are incomplete, clear them
	// prevents us from loading clearly malformed payloads
	config := validation.ClearInvalidFields(original, errs)

	// identify the desired next version
	desired, ok := findUpdateFromConfig(config)
	if ok {
		glog.V(4).Infof("Desired version from spec is %#v", desired)
	} else {
		desired = optr.currentVersion()
		glog.V(4).Infof("Desired version from operator is %#v", desired)
	}

	// handle the case of a misconfigured CVO by doing nothing
	if len(desired.Image) == 0 {
		return optr.syncStatus(original, config, &SyncWorkerStatus{
			Failure: fmt.Errorf("No configured operator version, unable to update cluster"),
		}, errs)
	}

	// inform the config sync loop about our desired state
	reconciling := resourcemerge.IsOperatorStatusConditionTrue(config.Status.Conditions, configv1.OperatorAvailable) &&
		resourcemerge.IsOperatorStatusConditionFalse(config.Status.Conditions, configv1.OperatorProgressing)
	status := optr.configSync.Update(config.Generation, desired, config.Spec.Overrides, reconciling)

	// write cluster version status
	return optr.syncStatus(original, config, status, errs)
}

// availableUpdatesSync is triggered on cluster version change (and periodic requeues) to
// sync available updates. It only modifies cluster version.
func (optr *Operator) availableUpdatesSync(key string) error {
	startTime := time.Now()
	glog.V(4).Infof("Started syncing available updates %q (%v)", key, startTime)
	defer func() {
		glog.V(4).Infof("Finished syncing available updates %q (%v)", key, time.Since(startTime))
	}()

	config, err := optr.cvLister.Get(optr.name)
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if errs := validation.ValidateClusterVersion(config); len(errs) > 0 {
		return nil
	}

	return optr.syncAvailableUpdates(config)
}

// isOlderThanLastUpdate returns true if the cluster version is older than
// the last update we saw.
func (optr *Operator) isOlderThanLastUpdate(config *configv1.ClusterVersion) bool {
	i, err := strconv.ParseInt(config.ResourceVersion, 10, 64)
	if err != nil {
		return false
	}
	optr.lastAtLock.Lock()
	defer optr.lastAtLock.Unlock()
	return i < optr.lastResourceVersion
}

// rememberLastUpdate records the most recent resource version we
// have seen from the server for cluster versions.
func (optr *Operator) rememberLastUpdate(config *configv1.ClusterVersion) {
	if config == nil {
		return
	}
	i, err := strconv.ParseInt(config.ResourceVersion, 10, 64)
	if err != nil {
		return
	}
	optr.lastAtLock.Lock()
	defer optr.lastAtLock.Unlock()
	optr.lastResourceVersion = i
}

func (optr *Operator) getOrCreateClusterVersion() (*configv1.ClusterVersion, bool, error) {
	obj, err := optr.cvLister.Get(optr.name)
	if err == nil {
		// if we are waiting to see a newer cached version, just exit
		if optr.isOlderThanLastUpdate(obj) {
			return nil, true, nil
		}
		return obj, false, nil
	}

	if !apierrors.IsNotFound(err) {
		return nil, false, err
	}

	var upstream configv1.URL
	if len(optr.defaultUpstreamServer) > 0 {
		u := configv1.URL(optr.defaultUpstreamServer)
		upstream = u
	}
	id, _ := uuid.NewRandom()

	// XXX: generate ClusterVersion from options calculated above.
	config := &configv1.ClusterVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: optr.name,
		},
		Spec: configv1.ClusterVersionSpec{
			Upstream:  upstream,
			Channel:   "fast",
			ClusterID: configv1.ClusterID(id.String()),
		},
	}

	actual, _, err := resourceapply.ApplyClusterVersionFromCache(optr.cvLister, optr.client.ConfigV1(), config)
	if apierrors.IsAlreadyExists(err) {
		return nil, true, nil
	}
	return actual, true, err
}

// versionString returns a string describing the desired version.
func versionString(update configv1.Update) string {
	if len(update.Version) > 0 {
		return update.Version
	}
	if len(update.Image) > 0 {
		return update.Image
	}
	return "<unknown>"
}

// currentVersion returns an update object describing the current known cluster version.
func (optr *Operator) currentVersion() configv1.Update {
	return configv1.Update{
		Version: optr.releaseVersion,
		Image:   optr.releaseImage,
	}
}

// SetSyncWorkerForTesting updates the sync worker for whitebox testing.
func (optr *Operator) SetSyncWorkerForTesting(worker ConfigSyncWorker) {
	optr.configSync = worker
}

// resourceBuilder provides the default builder implementation for the operator.
// It is abstracted for testing.
type resourceBuilder struct {
	config   *rest.Config
	modifier resourcebuilder.MetaV1ObjectModifierFunc
}

// NewResourceBuilder creates the default resource builder implementation.
func NewResourceBuilder(config *rest.Config) ResourceBuilder {
	return &resourceBuilder{config: config}
}

func (b *resourceBuilder) BuilderFor(m *lib.Manifest) (resourcebuilder.Interface, error) {
	if resourcebuilder.Mapper.Exists(m.GVK) {
		return resourcebuilder.New(resourcebuilder.Mapper, b.config, *m)
	}
	client, err := dynamicclient.New(b.config, m.GVK, m.Object().GetNamespace())
	if err != nil {
		return nil, err
	}
	return internal.NewGenericBuilder(client, *m)
}

func (b *resourceBuilder) Apply(m *lib.Manifest) error {
	builder, err := b.BuilderFor(m)
	if err != nil {
		return err
	}
	if b.modifier != nil {
		builder = builder.WithModifier(b.modifier)
	}
	return builder.Do()
}

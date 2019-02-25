package cvo

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	configv1 "github.com/openshift/api/config/v1"

	"github.com/openshift/cluster-version-operator/lib"
	"github.com/openshift/cluster-version-operator/lib/resourcebuilder"
	"github.com/openshift/cluster-version-operator/pkg/payload"
)

// ConfigSyncWorker abstracts how the image is synchronized to the server. Introduced for testing.
type ConfigSyncWorker interface {
	Start(maxWorkers int, stopCh <-chan struct{})
	Update(generation int64, desired configv1.Update, overrides []configv1.ComponentOverride, reconciling bool) *SyncWorkerStatus
	StatusCh() <-chan SyncWorkerStatus
}

// PayloadRetriever abstracts how a desired version is extracted to disk. Introduced for testing.
type PayloadRetriever interface {
	RetrievePayload(ctx context.Context, desired configv1.Update) (string, error)
}

// ResourceBuilder abstracts how a manifest is created on the server. Introduced for testing.
type ResourceBuilder interface {
	Apply(*lib.Manifest) error
}

// StatusReporter abstracts how status is reported by the worker run method. Introduced for testing.
type StatusReporter interface {
	Report(status SyncWorkerStatus)
}

// SyncWork represents the work that should be done in a sync iteration.
type SyncWork struct {
	Generation  int64
	Desired     configv1.Update
	Overrides   []configv1.ComponentOverride
	Reconciling bool
	Completed   int
}

// Empty returns true if the image is empty for this work.
func (w SyncWork) Empty() bool {
	return len(w.Desired.Image) == 0
}

// SyncWorkerStatus is the status of the sync worker at a given time.
type SyncWorkerStatus struct {
	Generation int64

	Step    string
	Failure error

	Fraction float32

	Completed   int
	Reconciling bool
	VersionHash string

	Actual configv1.Update
}

// DeepCopy copies the worker status.
func (w SyncWorkerStatus) DeepCopy() *SyncWorkerStatus {
	return &w
}

// SyncWorker retrieves and applies the desired image, tracking the status for the parent to
// monitor. The worker accepts a desired state via Update() and works to keep that state in
// sync. Once a particular image version is synced, it will be updated no more often than
// minimumReconcileInterval.
//
// State transitions:
//
//   Initial: wait for valid Update(), report empty status
//     Update() -> Sync
//   Sync: attempt to invoke the syncOnce() method
//     syncOnce() returns an error -> Error
//     syncOnce() returns nil -> Reconciling
//   Reconciling: invoke syncOnce() no more often than reconcileInterval
//     Update() with different values -> Sync
//     syncOnce() returns an error -> Error
//     syncOnce() returns nil -> Reconciling
//   Error: backoff until we are attempting every reconcileInterval
//     syncOnce() returns an error -> Error
//     syncOnce() returns nil -> Reconciling
//
type SyncWorker struct {
	backoff     wait.Backoff
	retriever   PayloadRetriever
	builder     ResourceBuilder
	reconciling bool

	// minimumReconcileInterval is the minimum time between reconcile attempts, and is
	// used to define the maximum backoff interval when syncOnce() returns an error.
	minimumReconcileInterval time.Duration

	// coordination between the sync loop and external callers
	notify chan struct{}
	report chan SyncWorkerStatus

	// lock guards changes to these fields
	lock     sync.Mutex
	work     *SyncWork
	cancelFn func()
	status   SyncWorkerStatus

	// updated by the run method only
	payload *payload.Update
}

// NewSyncWorker initializes a ConfigSyncWorker that will retrieve payloads to disk, apply them via builder
// to a server, and obey limits about how often to reconcile or retry on errors.
func NewSyncWorker(retriever PayloadRetriever, builder ResourceBuilder, reconcileInterval time.Duration, backoff wait.Backoff) ConfigSyncWorker {
	return &SyncWorker{
		retriever: retriever,
		builder:   builder,
		backoff:   backoff,

		minimumReconcileInterval: reconcileInterval,

		notify: make(chan struct{}, 1),
		// report is a large buffered channel to improve local testing - most consumers should invoke
		// Status() or use the result of calling Update() instead because the channel can be out of date
		// if the reader is not fast enough.
		report: make(chan SyncWorkerStatus, 500),
	}
}

// StatusCh returns a channel that reports status from the worker. The channel is buffered and events
// can be lost, so this is best used as a trigger to read the latest status.
func (w *SyncWorker) StatusCh() <-chan SyncWorkerStatus {
	return w.report
}

// Update instructs the sync worker to start synchronizing the desired update. The reconciling boolean is
// ignored unless this is the first time that Update has been called. The returned status represents either
// the initial state or whatever the last recorded status was.
// TODO: in the future it may be desirable for changes that alter desired to wait briefly before returning,
//   giving the sync loop the opportunity to observe our change and begin working towards it.
func (w *SyncWorker) Update(generation int64, desired configv1.Update, overrides []configv1.ComponentOverride, reconciling bool) *SyncWorkerStatus {
	w.lock.Lock()
	defer w.lock.Unlock()

	work := &SyncWork{
		Generation: generation,
		Desired:    desired,
		Overrides:  overrides,
	}

	if work.Empty() || equalSyncWork(w.work, work) {
		return w.status.DeepCopy()
	}

	// initialize the reconciliation flag and the status the first time
	// update is invoked
	if w.work == nil {
		if reconciling {
			work.Reconciling = true
		}
		w.status = SyncWorkerStatus{
			Generation:  generation,
			Reconciling: work.Reconciling,
			Actual:      work.Desired,
		}
	}

	// notify the sync loop that we changed config
	w.work = work
	if w.cancelFn != nil {
		w.cancelFn()
		w.cancelFn = nil
	}
	select {
	case w.notify <- struct{}{}:
	default:
	}

	return w.status.DeepCopy()
}

// Start periodically invokes run, detecting whether content has changed.
// It is edge-triggered when Update() is invoked and level-driven after the
// syncOnce() has succeeded for a given input (we are said to be "reconciling").
func (w *SyncWorker) Start(maxWorkers int, stopCh <-chan struct{}) {
	glog.V(5).Infof("Starting sync worker")

	work := &SyncWork{}

	wait.Until(func() {
		consecutiveErrors := 0
		errorInterval := w.minimumReconcileInterval / 16

		var next <-chan time.Time
		for {
			waitingToReconcile := work.Reconciling
			select {
			case <-stopCh:
				glog.V(5).Infof("Stopped worker")
				return
			case <-next:
				waitingToReconcile = false
				glog.V(5).Infof("Wait finished")
			case <-w.notify:
				glog.V(5).Infof("Work updated")
			}

			// determine whether we need to do work
			changed := w.calculateNext(work)
			if !changed && waitingToReconcile {
				glog.V(5).Infof("No change, waiting")
				continue
			}

			// until Update() has been called at least once, we do nothing
			if work.Empty() {
				next = time.After(w.minimumReconcileInterval)
				glog.V(5).Infof("No work, waiting")
				continue
			}

			// actually apply the image, allowing for calls to be cancelled
			err := func() error {
				ctx, cancelFn := context.WithCancel(context.Background())
				w.lock.Lock()
				w.cancelFn = cancelFn
				w.lock.Unlock()
				defer cancelFn()

				// reporter hides status updates that occur earlier than the previous failure,
				// so that we don't fail, then immediately start reporting an earlier status
				reporter := &statusWrapper{w: w, previousStatus: w.Status()}
				glog.V(5).Infof("Previous sync status: %#v", reporter.previousStatus)
				return w.syncOnce(ctx, work, maxWorkers, reporter)
			}()
			if err != nil {
				// backoff wait
				// TODO: replace with wait.Backoff when 1.13 client-go is available
				consecutiveErrors++
				interval := w.minimumReconcileInterval
				if consecutiveErrors < 4 {
					interval = errorInterval
					for i := 0; i < consecutiveErrors; i++ {
						interval *= 2
					}
				}
				next = time.After(wait.Jitter(interval, 0.2))

				utilruntime.HandleError(fmt.Errorf("unable to synchronize image (waiting %s): %v", interval, err))
				continue
			}
			glog.V(5).Infof("Sync succeeded, reconciling")

			work.Completed++
			work.Reconciling = true
			next = time.After(w.minimumReconcileInterval)
		}
	}, 10*time.Millisecond, stopCh)

	glog.V(5).Infof("Worker shut down")
}

// statusWrapper prevents a newer status update from overwriting a previous
// failure from later in the sync process.
type statusWrapper struct {
	w              *SyncWorker
	previousStatus *SyncWorkerStatus
}

func (w *statusWrapper) Report(status SyncWorkerStatus) {
	p := w.previousStatus
	if p.Failure != nil && status.Failure == nil {
		if p.Actual == status.Actual {
			if status.Fraction < p.Fraction {
				glog.V(5).Infof("Dropping status report from earlier in sync loop")
				return
			}
		}
	}
	w.w.updateStatus(status)
}

// calculateNext updates the passed work object with the desired next state and
// returns true if any changes were made. The reconciling flag is set the first
// time work transitions from empty to not empty (as a result of someone invoking
// Update).
func (w *SyncWorker) calculateNext(work *SyncWork) bool {
	w.lock.Lock()
	defer w.lock.Unlock()

	changed := !equalSyncWork(w.work, work)

	// if this is the first time through the loop, initialize reconciling to
	// the state Update() calculated (to allow us to start in reconciling)
	if work.Empty() {
		work.Reconciling = w.work.Reconciling
	} else {
		if changed {
			work.Reconciling = false
		}
	}
	// always clear the completed variable if we are not reconciling
	if !work.Reconciling {
		work.Completed = 0
	}

	if w.work != nil {
		work.Desired = w.work.Desired
		work.Overrides = w.work.Overrides
	}

	work.Generation = w.work.Generation

	return changed
}

// equalUpdate returns true if two updates have the same image.
func equalUpdate(a, b configv1.Update) bool {
	return a.Image == b.Image
}

// equalSyncWork returns true if a and b are equal.
func equalSyncWork(a, b *SyncWork) bool {
	if a == b {
		return true
	}
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	return equalUpdate(a.Desired, b.Desired) && reflect.DeepEqual(a.Overrides, b.Overrides)
}

// updateStatus records the current status of the sync action for observation
// by others. It sends a copy of the update to the report channel for improved
// testability.
func (w *SyncWorker) updateStatus(update SyncWorkerStatus) {
	w.lock.Lock()
	defer w.lock.Unlock()

	glog.V(5).Infof("Status change %#v", update)
	w.status = update
	select {
	case w.report <- update:
	default:
		if glog.V(5) {
			glog.Infof("Status report channel was full %#v", update)
		}
	}
}

// Desired returns the state the SyncWorker is trying to achieve.
func (w *SyncWorker) Desired() configv1.Update {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.work == nil {
		return configv1.Update{}
	}
	return w.work.Desired
}

// Status returns a copy of the current worker status.
func (w *SyncWorker) Status() *SyncWorkerStatus {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.status.DeepCopy()
}

// sync retrieves the image and applies it to the server, returning an error if
// the update could not be completely applied. The status is updated as we progress.
// Cancelling the context will abort the execution of the sync.
func (w *SyncWorker) syncOnce(ctx context.Context, work *SyncWork, maxWorkers int, reporter StatusReporter) error {
	glog.V(4).Infof("Running sync %s on generation %d", versionString(work.Desired), work.Generation)
	update := work.Desired

	// cache the payload until the release image changes
	validPayload := w.payload
	if validPayload == nil || !equalUpdate(configv1.Update{Image: validPayload.ReleaseImage}, update) {
		glog.V(4).Infof("Loading payload")
		reporter.Report(SyncWorkerStatus{Step: "RetrievePayload", Reconciling: work.Reconciling, Actual: update})
		payloadDir, err := w.retriever.RetrievePayload(ctx, update)
		if err != nil {
			reporter.Report(SyncWorkerStatus{Failure: err, Step: "RetrievePayload", Reconciling: work.Reconciling, Actual: update})
			return err
		}
		payloadUpdate, err := payload.LoadUpdate(payloadDir, update.Image)
		if err != nil {
			reporter.Report(SyncWorkerStatus{Failure: err, Step: "VerifyPayload", Reconciling: work.Reconciling, Actual: update})
			return err
		}
		w.payload = payloadUpdate
		glog.V(4).Infof("Payload loaded from %s with hash %s", payloadUpdate.ReleaseImage, payloadUpdate.ManifestHash)
	}

	return w.apply(ctx, w.payload, work, maxWorkers, reporter)
}

// apply updates the server with the contents of the provided image or returns an error.
// Cancelling the context will abort the execution of the sync. Will be executed in parallel if
// maxWorkers is set greater than 1.
func (w *SyncWorker) apply(ctx context.Context, payloadUpdate *payload.Update, work *SyncWork, maxWorkers int, reporter StatusReporter) error {
	update := configv1.Update{
		Version: payloadUpdate.ReleaseVersion,
		Image:   payloadUpdate.ReleaseImage,
	}

	// encapsulate status reporting in a threadsafe updater
	version := payloadUpdate.ReleaseVersion
	total := len(payloadUpdate.Manifests)
	cr := &consistentReporter{
		status: SyncWorkerStatus{
			Generation:  work.Generation,
			Reconciling: work.Reconciling,
			VersionHash: payloadUpdate.ManifestHash,
			Actual:      update,
		},
		completed: work.Completed,
		version:   version,
		total:     total,
		reporter:  reporter,
	}

	var tasks []*payload.Task
	for i := range payloadUpdate.Manifests {
		tasks = append(tasks, &payload.Task{
			Index:    i + 1,
			Total:    total,
			Manifest: &payloadUpdate.Manifests[i],
			Backoff:  w.backoff,
		})
	}
	graph := payload.NewTaskGraph(tasks)
	graph.Split(payload.SplitOnJobs)
	graph.Parallelize(payload.ByNumberAndComponent)

	// update each object
	err := payload.RunGraph(ctx, graph, maxWorkers, func(ctx context.Context, tasks []*payload.Task) error {
		for _, task := range tasks {
			if contextIsCancelled(ctx) {
				return cr.CancelError()
			}
			cr.Update()

			glog.V(4).Infof("Running sync for %s", task)
			glog.V(5).Infof("Manifest: %s", string(task.Manifest.Raw))

			ov, ok := getOverrideForManifest(work.Overrides, task.Manifest)
			if ok && ov.Unmanaged {
				glog.V(4).Infof("Skipping %s as unmanaged", task)
				continue
			}

			if err := task.Run(version, w.builder); err != nil {
				return err
			}
			cr.Inc()
			glog.V(4).Infof("Done syncing for %s", task)
		}
		return nil
	})
	if err != nil {
		cr.Error(err)
		return err
	}

	// update the
	cr.Complete()
	return nil
}

var (
	metricPayload = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cluster_version_payload",
		Help: "Report the number of entries in the payload.",
	}, []string{"version", "type"})
)

func init() {
	prometheus.MustRegister(
		metricPayload,
	)
}

// consistentReporter hides the details of calculating the status based on the progress
// of the graph runner.
type consistentReporter struct {
	lock      sync.Mutex
	status    SyncWorkerStatus
	version   string
	completed int
	total     int
	done      int
	reporter  StatusReporter
}

func (r *consistentReporter) Inc() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.done++
}

func (r *consistentReporter) Update() {
	r.lock.Lock()
	defer r.lock.Unlock()
	metricPayload.WithLabelValues(r.version, "pending").Set(float64(r.total - r.done))
	metricPayload.WithLabelValues(r.version, "applied").Set(float64(r.done))
	copied := r.status
	copied.Step = "ApplyResources"
	copied.Fraction = float32(r.done) / float32(r.total)
	r.reporter.Report(copied)
}

func (r *consistentReporter) Error(err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	copied := r.status
	copied.Step = "ApplyResources"
	copied.Fraction = float32(r.done) / float32(r.total)
	copied.Failure = err
	r.reporter.Report(copied)
}

func (r *consistentReporter) CancelError() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return fmt.Errorf("update was cancelled at %d/%d", r.done, r.total)
}

func (r *consistentReporter) Complete() {
	r.lock.Lock()
	defer r.lock.Unlock()
	metricPayload.WithLabelValues(r.version, "pending").Set(float64(r.total))
	metricPayload.WithLabelValues(r.version, "applied").Set(float64(r.total))
	copied := r.status
	copied.Completed = r.completed + 1
	copied.Reconciling = true
	copied.Fraction = 1
	r.reporter.Report(copied)
}

// getOverrideForManifest returns the override and true when override exists for manifest.
func getOverrideForManifest(overrides []configv1.ComponentOverride, manifest *lib.Manifest) (configv1.ComponentOverride, bool) {
	for idx, ov := range overrides {
		kind, namespace, name := manifest.GVK.Kind, manifest.Object().GetNamespace(), manifest.Object().GetName()
		if ov.Kind == kind &&
			(namespace == "" || ov.Namespace == namespace) && // cluster-scoped objects don't have namespace.
			ov.Name == name {
			return overrides[idx], true
		}
	}
	return configv1.ComponentOverride{}, false
}

// ownerKind contains the schema.GroupVersionKind for type that owns objects managed by CVO.
var ownerKind = configv1.SchemeGroupVersion.WithKind("ClusterVersion")

func ownerRefModifier(config *configv1.ClusterVersion) resourcebuilder.MetaV1ObjectModifierFunc {
	oref := metav1.NewControllerRef(config, ownerKind)
	return func(obj metav1.Object) {
		obj.SetOwnerReferences([]metav1.OwnerReference{*oref})
	}
}

// contextIsCancelled returns true if the provided context is cancelled.
func contextIsCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// runThrottledStatusNotifier invokes fn every time ch is updated, but no more often than once
// every interval. If bucket is non-zero then the channel is throttled like a rate limiter bucket.
func runThrottledStatusNotifier(stopCh <-chan struct{}, interval time.Duration, bucket int, ch <-chan SyncWorkerStatus, fn func()) {
	// notify the status change function fairly infrequently to avoid updating
	// the caller status more frequently than is needed
	throttle := rate.NewLimiter(rate.Every(interval), bucket)
	wait.Until(func() {
		ctx := context.Background()
		var last SyncWorkerStatus
		for {
			select {
			case <-stopCh:
				return
			case next := <-ch:
				// only throttle if we aren't on an edge
				if next.Generation == last.Generation && next.Actual == last.Actual && next.Reconciling == last.Reconciling && (next.Failure != nil) == (last.Failure != nil) {
					if err := throttle.Wait(ctx); err != nil {
						utilruntime.HandleError(fmt.Errorf("unable to throttle status notification: %v", err))
					}
				}
				last = next

				fn()
			}
		}
	}, 1*time.Second, stopCh)
}

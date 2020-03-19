package termination

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/go-logr/logr"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	awsTerminationEndpointURL = "http://169.254.169.254/latest/meta-data/spot/termination-time"
)

// Handler represents a handler that will run to check the termination
// notice endpoint and delete Machine's if the instance termination notice is fulfilled.
type Handler interface {
	Run(stop <-chan struct{}) error
}

// NewHandler constructs a new Handler
func NewHandler(logger logr.Logger, cfg *rest.Config, pollInterval time.Duration, nodeName string) (Handler, error) {
	c, err := client.New(cfg, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	pollURL, err := url.Parse(awsTerminationEndpointURL)
	if err != nil {
		// This should never happen
		panic(err)
	}

	logger = logger.WithValues("node", nodeName)

	return &handler{
		client:       c,
		pollURL:      pollURL,
		pollInterval: pollInterval,
		nodeName:     nodeName,
		log:          logger,
	}, nil
}

// handler implements the logic to check the termination endpoint and delete the
// machine associated with the node
type handler struct {
	client       client.Client
	pollURL      *url.URL
	pollInterval time.Duration
	nodeName     string
	log          logr.Logger
}

// Run starts the handler and runs the termination logic
func (h *handler) Run(stop <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())

	errs := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		errs <- h.run(ctx, wg)
	}()

	select {
	case <-stop:
		cancel()
		// Wait for run to stop
		wg.Wait()
		return nil
	case err := <-errs:
		cancel()
		return err
	}
}

func (h *handler) run(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	machine, err := h.getMachineForNode(ctx)
	if err != nil {
		return fmt.Errorf("error fetching machine for node (%q): %v", h.nodeName, err)
	}

	logger := h.log.WithValues("namespace", machine.Namespace, "machine", machine.Name)
	logger.V(1).Info("Monitoring node for machine")

	return fmt.Errorf("not implemented")
}

// getMachineForNodeName finds the Machine associated with the Node name given
func (h *handler) getMachineForNode(ctx context.Context) (*machinev1.Machine, error) {
	return nil, fmt.Errorf("machine not found for node %q", h.nodeName)
}

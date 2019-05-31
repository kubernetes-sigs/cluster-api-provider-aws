package operator

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/appscode/jsonpatch"
	"k8s.io/apimachinery/pkg/types"
	admissionregistrationv1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// WebhookConfigurationName is the name of the webhook configurations to be
// updated with the current CA certificate.
const WebhookConfigurationName = "autoscaling.openshift.io"

// WebhookCAUpdater updates webhook configuratons to inject the CA certficiate
// bundle read from disk at the configured location.
type WebhookCAUpdater struct {
	caPath string
	client *admissionregistrationv1beta1.AdmissionregistrationV1beta1Client
}

// NewWebhookCAUpdater returns a new WebhookCAUpdater.
func NewWebhookCAUpdater(mgr manager.Manager, caPath string) (*WebhookCAUpdater, error) {
	var err error

	w := &WebhookCAUpdater{caPath: caPath}

	w.client, err = admissionregistrationv1beta1.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Start fetches the current CA bundle from disk and sets it on the webhook
// configurations, then simply waits for the stop channel to close.
func (w *WebhookCAUpdater) Start(stopCh <-chan struct{}) error {
	ca, err := w.GetEncodedCA()
	if err != nil {
		return err
	}

	// TODO: This should probably replace the caBundle in all webhook client
	// configurations in the object, but unfortuntaely that's not easy to do
	// with a JSON patch.  For now this only modifies the first entry.
	patch := []jsonpatch.Operation{
		{
			Operation: "replace",
			Path:      "/webhooks/0/clientConfig/caBundle",
			Value:     ca,
		},
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	_, err = w.client.ValidatingWebhookConfigurations().
		Patch(WebhookConfigurationName, types.JSONPatchType, patchBytes)

	if err != nil {
		return err
	}

	klog.Info("Updated webhook configuration CA certificates.")

	// Block until the stop channel is closed.
	<-stopCh

	return nil
}

// GetEncodedCA returns the base64 encoded CA certificate used for securing
// admission webhook server connections.
func (w *WebhookCAUpdater) GetEncodedCA() (string, error) {
	ca, err := ioutil.ReadFile(w.caPath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ca), nil
}

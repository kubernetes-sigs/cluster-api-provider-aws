package clusterautoscaler

import (
	"context"
	"net/http"

	autoscalingv1 "github.com/openshift/cluster-autoscaler-operator/pkg/apis/autoscaling/v1"
	"k8s.io/klog"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator validates ClusterAutoscaler resources.
type Validator struct {
	client  client.Client
	decoder *admission.Decoder
}

// Handle handles HTTP requests for admission webhook servers.
func (v *Validator) Handle(ctx context.Context, req admission.Request) admission.Response {
	ca := &autoscalingv1.ClusterAutoscaler{}

	if err := v.decoder.Decode(req, ca); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	klog.Infof("Validation webhook called for ClustAutoscaler: %s", ca.GetName())

	// TODO: Implement validations.
	return admission.Allowed("ALLOW ALL")
}

// InjectClient injects the client.
func (v *Validator) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

// InjectDecoder injects the decoder.
func (v *Validator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

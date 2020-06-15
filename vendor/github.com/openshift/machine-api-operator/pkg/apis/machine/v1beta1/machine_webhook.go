package v1beta1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	osconfigv1 "github.com/openshift/api/config/v1"
	osclientset "github.com/openshift/client-go/config/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog"
	"k8s.io/utils/pointer"
	aws "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	yaml "sigs.k8s.io/yaml"
)

var (
	defaultAWSIAMInstanceProfile = func(clusterID string) *string {
		return pointer.StringPtr(fmt.Sprintf("%s-worker-profile", clusterID))
	}
	defaultAWSSecurityGroup = func(clusterID string) string {
		return fmt.Sprintf("%s-worker-sg", clusterID)
	}
	defaultAWSSubnet = func(clusterID, az string) string {
		return fmt.Sprintf("%s-private-%s", clusterID, az)
	}
)

const (
	defaultAWSUserDataSecret    = "worker-user-data"
	defaultAWSCredentialsSecret = "aws-cloud-credentials"
	defaultAWSInstanceType      = "m4.large"
)

func getInfra() (*osconfigv1.Infrastructure, error) {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := osclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	infra, err := client.ConfigV1().Infrastructures().Get(context.Background(), "cluster", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return infra, nil
}

type handlerValidationFn func(h *validatorHandler, m *Machine) (bool, utilerrors.Aggregate)
type handlerMutationFn func(h *defaulterHandler, m *Machine) (bool, utilerrors.Aggregate)

// validatorHandler validates Machine API resources.
// implements type Handler interface.
// https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/webhook/admission#Handler
type validatorHandler struct {
	clusterID         string
	webhookOperations handlerValidationFn
	decoder           *admission.Decoder
}

// defaulterHandler defaults Machine API resources.
// implements type Handler interface.
// https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/webhook/admission#Handler
type defaulterHandler struct {
	clusterID         string
	webhookOperations handlerMutationFn
	decoder           *admission.Decoder
}

// NewValidator returns a new validatorHandler.
func NewMachineValidator() (*validatorHandler, error) {
	infra, err := getInfra()
	if err != nil {
		return nil, err
	}

	return createMachineValidator(infra.Status.PlatformStatus.Type, infra.Status.InfrastructureName), nil
}

func createMachineValidator(platform osconfigv1.PlatformType, clusterID string) *validatorHandler {
	h := &validatorHandler{
		clusterID: clusterID,
	}

	switch platform {
	case osconfigv1.AWSPlatformType:
		h.webhookOperations = validateAWS
	default:
		// just no-op
		h.webhookOperations = func(h *validatorHandler, m *Machine) (bool, utilerrors.Aggregate) {
			return true, nil
		}
	}
	return h
}

// NewDefaulter returns a new defaulterHandler.
func NewMachineDefaulter() (*defaulterHandler, error) {
	infra, err := getInfra()
	if err != nil {
		return nil, err
	}

	return createMachineDefaulter(infra.Status.PlatformStatus.Type, infra.Status.InfrastructureName), nil
}

func createMachineDefaulter(platform osconfigv1.PlatformType, clusterID string) *defaulterHandler {
	h := &defaulterHandler{
		clusterID: clusterID,
	}

	switch platform {
	case osconfigv1.AWSPlatformType:
		h.webhookOperations = defaultAWS
	default:
		// just no-op
		h.webhookOperations = func(h *defaulterHandler, m *Machine) (bool, utilerrors.Aggregate) {
			return true, nil
		}
	}
	return h
}

// InjectDecoder injects the decoder.
func (v *validatorHandler) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

// InjectDecoder injects the decoder.
func (v *defaulterHandler) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

// Handle handles HTTP requests for admission webhook servers.
func (h *validatorHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	m := &Machine{}

	if err := h.decoder.Decode(req, m); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	klog.V(3).Infof("Validate webhook called for Machine: %s", m.GetName())

	if ok, err := h.webhookOperations(h, m); !ok {
		return admission.Denied(err.Error())
	}

	return admission.Allowed("Machine valid")
}

// Handle handles HTTP requests for admission webhook servers.
func (h *defaulterHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	m := &Machine{}

	if err := h.decoder.Decode(req, m); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	klog.V(3).Infof("Mutate webhook called for Machine: %s", m.GetName())

	if ok, err := h.webhookOperations(h, m); !ok {
		return admission.Denied(err.Error())
	}

	marshaledMachine, err := json.Marshal(m)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledMachine)
}

func defaultAWS(h *defaulterHandler, m *Machine) (bool, utilerrors.Aggregate) {
	klog.V(3).Infof("Defaulting AWS providerSpec")

	var errs []error
	providerSpec := new(aws.AWSMachineProviderConfig)
	if err := yaml.Unmarshal(m.Spec.ProviderSpec.Value.Raw, &providerSpec); err != nil {
		errs = append(
			errs,
			field.Invalid(
				field.NewPath("providerSpec", "value"),
				providerSpec,
				err.Error(),
			),
		)
		return false, utilerrors.NewAggregate(errs)
	}

	if providerSpec.InstanceType == "" {
		providerSpec.InstanceType = defaultAWSInstanceType
	}
	if providerSpec.IAMInstanceProfile == nil {
		providerSpec.IAMInstanceProfile = &aws.AWSResourceReference{ID: defaultAWSIAMInstanceProfile(h.clusterID)}
	}
	if providerSpec.UserDataSecret == nil {
		providerSpec.UserDataSecret = &corev1.LocalObjectReference{Name: defaultAWSUserDataSecret}
	}

	if providerSpec.CredentialsSecret == nil {
		providerSpec.CredentialsSecret = &corev1.LocalObjectReference{Name: defaultAWSCredentialsSecret}
	}

	if providerSpec.SecurityGroups == nil {
		providerSpec.SecurityGroups = []aws.AWSResourceReference{
			{
				Filters: []aws.Filter{
					{
						Name:   "tag:Name",
						Values: []string{defaultAWSSecurityGroup(h.clusterID)},
					},
				},
			},
		}
	}

	if providerSpec.Subnet.ARN == nil && providerSpec.Subnet.ID == nil && providerSpec.Subnet.Filters == nil {
		providerSpec.Subnet.Filters = []aws.Filter{
			{
				Name:   "tag:Name",
				Values: []string{defaultAWSSubnet(h.clusterID, providerSpec.Placement.AvailabilityZone)},
			},
		}
	}

	rawBytes, err := json.Marshal(providerSpec)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return false, utilerrors.NewAggregate(errs)
	}

	m.Spec.ProviderSpec.Value = &runtime.RawExtension{Raw: rawBytes}
	return true, nil
}

func validateAWS(h *validatorHandler, m *Machine) (bool, utilerrors.Aggregate) {
	klog.V(3).Infof("Validating AWS providerSpec")

	var errs []error
	providerSpec := new(aws.AWSMachineProviderConfig)
	if err := yaml.Unmarshal(m.Spec.ProviderSpec.Value.Raw, &providerSpec); err != nil {
		errs = append(
			errs,
			field.Invalid(
				field.NewPath("providerSpec", "value"),
				providerSpec,
				err.Error(),
			),
		)
		return false, utilerrors.NewAggregate(errs)
	}

	if providerSpec.AMI.ARN == nil && providerSpec.AMI.Filters == nil && providerSpec.AMI.ID == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "ami"),
				"expected either providerSpec.ami.arn or providerSpec.ami.filters or providerSpec.ami.id to be populated",
			),
		)
	}

	if providerSpec.InstanceType == "" {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "instanceType"),
				"expected providerSpec.instanceType to be populated",
			),
		)
	}

	if providerSpec.IAMInstanceProfile == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "iamInstanceProfile"),
				"expected providerSpec.iamInstanceProfile to be populated",
			),
		)
	}

	if providerSpec.UserDataSecret == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "userDataSecret"),
				"expected providerSpec.userDataSecret to be populated",
			),
		)
	}

	if providerSpec.CredentialsSecret == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "credentialsSecret"),
				"expected providerSpec.credentialsSecret to be populated",
			),
		)
	}

	if providerSpec.SecurityGroups == nil {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "securityGroups"),
				"expected providerSpec.securityGroups to be populated",
			),
		)
	}

	if providerSpec.Subnet.ARN == nil && providerSpec.Subnet.ID == nil && providerSpec.Subnet.Filters == nil && providerSpec.Placement.AvailabilityZone == "" {
		errs = append(
			errs,
			field.Required(
				field.NewPath("providerSpec", "subnet"),
				"expected either providerSpec.subnet.arn or providerSpec.subnet.id or providerSpec.subnet.filters or providerSpec.placement.availabilityZone to be populated",
			),
		)
	}
	// TODO(alberto): Validate providerSpec.BlockDevices.
	// https://github.com/openshift/cluster-api-provider-aws/pull/299#discussion_r433920532

	if len(errs) > 0 {
		return false, utilerrors.NewAggregate(errs)
	}

	return true, nil
}

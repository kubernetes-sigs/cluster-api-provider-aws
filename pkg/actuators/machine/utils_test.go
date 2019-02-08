package machine

import (
	"reflect"
	"testing"

	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
}

func TestProviderConfigFromMachine(t *testing.T) {

	providerConfig := &providerconfigv1.AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "awsproviderconfig.openshift.io/v1beta1",
			Kind:       "AWSMachineProviderConfig",
		},
		InstanceType: "testInstance",
		AMI:          providerconfigv1.AWSResourceReference{ID: nil},
		Tags: []providerconfigv1.TagSpecification{
			{Name: "", Value: ""},
		},
		IAMInstanceProfile: &providerconfigv1.AWSResourceReference{ID: nil},
		UserDataSecret:     &corev1.LocalObjectReference{Name: ""},
		Subnet: providerconfigv1.AWSResourceReference{
			Filters: []providerconfigv1.Filter{{
				Name:   "tag:Name",
				Values: []string{""},
			}},
		},
		Placement: providerconfigv1.Placement{Region: "", AvailabilityZone: ""},
		SecurityGroups: []providerconfigv1.AWSResourceReference{{
			Filters: []providerconfigv1.Filter{{
				Name:   "tag:Name",
				Values: []string{""},
			}},
		}},
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Error(err)
	}
	encodedProviderSpec, err := codec.EncodeProviderSpec(providerConfig)
	if err != nil {
		t.Error(err)
	}

	machineClass := &machinev1.MachineClass{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-cluster-api",
			Name:      "testClass",
		},
		ProviderSpec: *encodedProviderSpec.Value,
	}

	testCases := []struct {
		machine *machinev1.Machine
	}{
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "configFromSpecProviderConfigValue",
					Namespace: "",
					Labels: map[string]string{
						"foo": "a",
					},
				},
				TypeMeta: metav1.TypeMeta{
					Kind: "Machine",
				},
				Spec: machinev1.MachineSpec{
					ProviderSpec: *encodedProviderSpec,
				},
			},
		},
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "configFromClass",
					Namespace: "",
					Labels: map[string]string{
						"foo": "a",
					},
				},
				TypeMeta: metav1.TypeMeta{
					Kind: "Machine",
				},
				Spec: machinev1.MachineSpec{
					ProviderSpec: machinev1.ProviderSpec{
						ValueFrom: &machinev1.ProviderSpecSource{
							MachineClass: &machinev1.MachineClassRef{
								ObjectReference: &corev1.ObjectReference{
									Kind:      "MachineClass",
									Name:      "testClass",
									Namespace: "openshift-cluster-api",
								},
							},
						},
					},
				},
			},
		},
	}

	client := fake.NewFakeClient(machineClass)
	for _, tc := range testCases {
		decodedProviderConfig, err := providerConfigFromMachine(client, tc.machine, codec)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(decodedProviderConfig, providerConfig) {
			t.Errorf("Test case %s. Expected: %v, got: %v", tc.machine.Name, providerConfig, decodedProviderConfig)
		}
	}
}

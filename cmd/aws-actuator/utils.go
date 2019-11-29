package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os/exec"

	"github.com/ghodss/yaml"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
)

type manifestParams struct {
	ClusterID string
}

func readMachineManifest(manifestParams *manifestParams, manifestLoc string) (*machinev1.Machine, error) {
	machine := &machinev1.Machine{}
	manifestBytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", manifestLoc, err)
	}

	t, err := template.New("machineuserdata").Parse(string(manifestBytes))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, *manifestParams)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buf.Bytes(), &machine); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v: %v", manifestLoc, err)
	}

	return machine, nil
}

func readClusterResources(manifestParams *manifestParams, machineLoc, awsCredentialSecretLoc, userDataLoc string) (*machinev1.Machine, *apiv1.Secret, *apiv1.Secret, error) {
	machine, err := readMachineManifest(manifestParams, machineLoc)
	if err != nil {
		return nil, nil, nil, err
	}

	var awsCredentialsSecret *apiv1.Secret
	if awsCredentialSecretLoc != "" {
		awsCredentialsSecret = &apiv1.Secret{}
		bytes, err := ioutil.ReadFile(awsCredentialSecretLoc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("aws credentials manifest %q: %v", awsCredentialSecretLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &awsCredentialsSecret); err != nil {
			return nil, nil, nil, fmt.Errorf("aws credentials manifest %q: %v", awsCredentialSecretLoc, err)
		}
	}

	var userDataSecret *apiv1.Secret
	if userDataLoc != "" {
		userDataSecret = &apiv1.Secret{}
		bytes, err := ioutil.ReadFile(userDataLoc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("user data manifest %q: %v", userDataLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &userDataSecret); err != nil {
			return nil, nil, nil, fmt.Errorf("user data manifest %q: %v", userDataLoc, err)
		}
	}

	return machine, awsCredentialsSecret, userDataSecret, nil
}

// CreateActuator creates actuator with fake clientsets
func createActuator(machine *machinev1.Machine, awsCredentials, userData *apiv1.Secret) (*machineactuator.Actuator, error) {
	objList := []runtime.Object{machine}
	if awsCredentials != nil {
		objList = append(objList, awsCredentials)
	}
	if userData != nil {
		objList = append(objList, userData)
	}
	fakeClient := fake.NewFakeClient(objList...)

	codec, err := v1beta1.NewCodec()
	if err != nil {
		return nil, err
	}

	params := machineactuator.ActuatorParams{
		Client:           fakeClient,
		AwsClientBuilder: awsclient.NewClient,
		Codec:            codec,
		// use empty recorder dropping any event recorded
		EventRecorder: &record.FakeRecorder{},
	}

	actuator, err := machineactuator.NewActuator(params)
	if err != nil {
		return nil, err
	}
	return actuator, nil
}

func cmdRun(binaryPath string, args ...string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	return cmd.CombinedOutput()
}

package utils

import (
	"github.com/golang/glog"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// CreateActuator creates actuator with fake clientsets
func CreateActuator(machine *clusterv1.Machine, awsCredentials *apiv1.Secret, userData *apiv1.Secret) *machineactuator.Actuator {
	objList := []runtime.Object{machine}
	if awsCredentials != nil {
		objList = append(objList, awsCredentials)
	}
	if userData != nil {
		objList = append(objList, userData)
	}
	fakeClient := fake.NewFakeClient(objList...)

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		glog.Fatal(err)
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
		glog.Error(err)
	}
	return actuator
}

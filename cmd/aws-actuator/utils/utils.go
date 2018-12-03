package utils

import (
	"github.com/golang/glog"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// CreateActuator creates actuator with fake clientsets
func CreateActuator(machine *clusterv1.Machine, awsCredentials *apiv1.Secret, userData *apiv1.Secret) *machineactuator.Actuator {
	objList := []runtime.Object{}
	if awsCredentials != nil {
		objList = append(objList, awsCredentials)
	}
	if userData != nil {
		objList = append(objList, userData)
	}
	fakeClient := fake.NewFakeClient(machine)
	fakeKubeClient := kubernetesfake.NewSimpleClientset(objList...)

	params := machineactuator.ActuatorParams{
		Client:           fakeClient,
		KubeClient:       fakeKubeClient,
		AwsClientBuilder: awsclient.NewClient,
	}

	actuator, err := machineactuator.NewActuator(params)
	if err != nil {
		glog.Error(err)
	}
	return actuator
}

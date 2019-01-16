/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package machine

// should not need to import the ec2 sdk here
import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/deployer"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/tokens"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	controllerError "sigs.k8s.io/cluster-api/pkg/controller/error"
)

// Actuator is responsible for performing machine reconciliation.
type Actuator struct {
	*deployer.Deployer

	client           client.ClusterV1alpha1Interface
	KubeClientConfig *rest.Config
}

// ActuatorParams holds parameter information for Actuator.
type ActuatorParams struct {
	Client           client.ClusterV1alpha1Interface
	KubeClientConfig *rest.Config
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		Deployer:         deployer.New(deployer.Params{ScopeGetter: actuators.DefaultScopeGetter}),
		client:           params.Client,
		KubeClientConfig: params.KubeClientConfig,
	}
}

// Create creates a machine and is invoked by the machine controller.
func (a *Actuator) Create(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	klog.Infof("Creating machine %v for cluster %v", machine.Name, cluster)
	var scope *actuators.MachineScope
	var err error
	if cluster != nil {
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client})
		if err != nil {
			return errors.Errorf("failed to create scope: %+v", err)
		}
	} else {
		// create scope without cluster
		if a.KubeClientConfig == nil {
			return errors.New("failed to initialize new corev1 client. Either cluster object or KubeClientConfig must be non-nil. Both are found to be nil")
		}
		kubeClient, err := kubernetes.NewForConfig(a.KubeClientConfig)
		if err != nil {
			return errors.Errorf("failed to initialize new k8s client: %+v", err)
		}
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client, KubeClient: kubeClient})
		if err != nil {
			return errors.Errorf("failed to create scope: %+v", err)
		}
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope)

	var bootstrapToken string
	var clientConfig *rest.Config
	if machine.ObjectMeta.Labels["set"] == "node" {
		if cluster != nil {
			kubeConfig, err := a.GetKubeConfig(cluster, nil)
			if err != nil {
				return errors.Errorf("failed to retrieve kubeconfig during machine creation: %+v", err)
			}
			controlPlaneURL, err := a.GetIP(cluster, nil)
			if err != nil {
				return errors.Errorf("failed to retrieve controlplane url during machine creation: %+v", err)
			}

			clientConfig, err = clientcmd.BuildConfigFromKubeconfigGetter(controlPlaneURL, func() (*clientcmdapi.Config, error) {
				return clientcmd.Load([]byte(kubeConfig))
			})

			if err != nil {
				return errors.Errorf("failed to retrieve kubeconfig during machine creation: %+v", err)
			}

		} else {
			// use kubeclient
			if a.KubeClientConfig == nil {
				return errors.New("failed to initialize new corev1 client. Either cluster object or KubeClientConfig must be non-nil. Both are found to be nil")
			}
			clientConfig = a.KubeClientConfig
		}
		coreClient, err := corev1.NewForConfig(clientConfig)
		if err != nil {
			return errors.Errorf("failed to initialize new corev1 client: %+v", err)
		}

		bootstrapToken, err = tokens.NewBootstrap(coreClient, 10*time.Minute)
		if err != nil {
			return errors.Errorf("failed to create new bootstrap token: %+v", err)
		}
	}

	i, err := ec2svc.CreateOrGetMachine(scope, bootstrapToken)
	if err != nil {
		if awserrors.IsFailedDependency(errors.Cause(err)) {
			klog.Errorf("network not ready to launch instances yet: %+v", err)
			return &controllerError.RequeueAfterError{
				RequeueAfter: time.Minute,
			}
		}

		return errors.Errorf("failed to create or get machine: %+v", err)
	}

	scope.MachineStatus.InstanceID = &i.ID
	scope.MachineStatus.InstanceState = aws.String(string(i.State))

	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}

	machine.Annotations["cluster-api-provider-aws"] = "true"

	if err := a.reconcileLBAttachment(scope, machine, i); err != nil {
		return errors.Errorf("failed to reconcile LB attachment: %+v", err)
	}

	return nil
}

func (a *Actuator) reconcileLBAttachment(scope *actuators.MachineScope, m *clusterv1.Machine, i *v1alpha1.Instance) error {
	elbsvc := elb.NewService(scope.Scope)
	if m.ObjectMeta.Labels["set"] == "controlplane" {
		if err := elbsvc.RegisterInstanceWithAPIServerELB(i.ID); err != nil {
			return errors.Wrapf(err, "could not register control plane instance %q with load balancer", i.ID)
		}
	}

	return nil
}

// Delete deletes a machine and is invoked by the Machine Controller
func (a *Actuator) Delete(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	klog.Infof("Deleting machine %+v for cluster %v.", machine, cluster)
	var scope *actuators.MachineScope
	var err error
	if cluster != nil {
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client})
		if err != nil {
			return errors.Errorf("failed to create scope: %+v", err)
		}
	} else {
		// create scope without cluster
		if a.KubeClientConfig == nil {
			return errors.New("failed to initialize new corev1 client. Either cluster object or KubeClientConfig must be non-nil. Both are found to be nil")
		}
		kubeClient, err := kubernetes.NewForConfig(a.KubeClientConfig)
		if err != nil {
			return errors.Errorf("failed to initialize new k8s client: %+v", err)
		}
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client, KubeClient: kubeClient})
		if err != nil {
			return errors.Errorf("failed to delete machine %v: %+v", machine.Name, err)
		}
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope)
	klog.Infof("machinestatus %+v", scope.MachineStatus)

	if scope.MachineStatus.InstanceID == nil {
		klog.Info("Instance is nil and therefore does not exist")
		return nil
	}

	instance, err := ec2svc.InstanceIfExists(*scope.MachineStatus.InstanceID)
	if err != nil {
		return errors.Errorf("failed to get instance: %+v", err)
	}

	if instance == nil {
		// The machine hasn't been created yet
		klog.Info("Instance is nil and therefore does not exist")
		return nil
	}

	// Check the instance state. If it's already shutting down or terminated,
	// do nothing. Otherwise attempt to delete it.
	// This decision is based on the ec2-instance-lifecycle graph at
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	switch instance.State {
	case v1alpha1.InstanceStateShuttingDown, v1alpha1.InstanceStateTerminated:
		klog.Infof("instance %q is shutting down or already terminated", machine.Name)
		return nil
	default:
		if err := ec2svc.TerminateInstance(aws.StringValue(scope.MachineStatus.InstanceID)); err != nil {
			return errors.Errorf("failed to terminate instance: %+v", err)
		}
	}

	klog.Info("shutdown signal was sent. Shutting down machine.")
	return nil
}

// Update updates a machine and is invoked by the Machine Controller.
// If the Update attempts to mutate any immutable state, the method will error
// and no updates will be performed.
func (a *Actuator) Update(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	klog.Infof("Updating machine %v for cluster %v.", machine.Name, cluster)
	var scope *actuators.MachineScope
	var err error
	if cluster != nil {
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client})
		if err != nil {
			return errors.Errorf("failed to create scope: %+v", err)
		}
	} else {
		// create scope without cluster
		if a.KubeClientConfig == nil {
			return errors.New("failed to initialize new corev1 client. Either cluster object or KubeClientConfig must be non-nil. Both are found to be nil")
		}
		kubeClient, err := kubernetes.NewForConfig(a.KubeClientConfig)
		if err != nil {
			return errors.Errorf("failed to initialize new k8s client: %+v", err)
		}
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client, KubeClient: kubeClient})
		if err != nil {
			return errors.Errorf("failed to update machine %v: %+v", machine.Name, err)
		}
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope)

	// Get the current instance description from AWS.
	instance, err := ec2svc.InstanceIfExists(*scope.MachineStatus.InstanceID)
	if err != nil {
		return errors.Errorf("failed to get instance: %+v", err)
	}
	klog.Infof("Found instance for machine %q: %v", machine.Name, instance)

	switch instance.State {
	case v1alpha1.InstanceStateRunning:
		klog.Infof("Machine %s, %v is running", machine.Name, scope.MachineStatus.InstanceID)
	case v1alpha1.InstanceStatePending:
		klog.Infof("Machine %s, %v is pending", machine.Name, scope.MachineStatus.InstanceID)
	default:
		return nil
	}
	scope.MachineStatus.InstanceState = aws.String(string(instance.State))

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	// TODO: Implement immutable state check.

	// Ensure that the security groups are correct.
	_, err = a.ensureSecurityGroups(
		ec2svc,
		machine,
		*scope.MachineStatus.InstanceID,
		scope.MachineConfig.AdditionalSecurityGroups,
		instance.SecurityGroupIDs,
	)
	if err != nil {
		return errors.Errorf("failed to apply security groups: %+v", err)
	}

	// Ensure that the tags are correct.
	_, err = a.ensureTags(ec2svc, machine, scope.MachineStatus.InstanceID, scope.MachineConfig.AdditionalTags)
	if err != nil {
		return errors.Errorf("failed to ensure tags: %+v", err)
	}

	return nil
}

// Exists test for the existence of a machine and is invoked by the Machine Controller
func (a *Actuator) Exists(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	klog.Infof("Checking if machine %v for cluster %v exists", machine.Name, cluster)

	var scope *actuators.MachineScope
	var err error
	if cluster != nil {
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client})
		if err != nil {
			return false, errors.Errorf("failed to create scope: %+v", err)
		}
	} else {
		// create scope without cluster
		if a.KubeClientConfig == nil {
			return false, errors.New("failed to initialize new corev1 client. Either cluster object or KubeClientConfig must be non-nil. Both are found to be nil")
		}
		kubeClient, err := kubernetes.NewForConfig(a.KubeClientConfig)
		if err != nil {
			return false, errors.Errorf("failed to initialize new k8s client: %+v", err)
		}
		scope, err = actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.client, KubeClient: kubeClient})
		if err != nil {
			return false, errors.Errorf("failed to verify existence of machine %v: %+v", machine.Name, err)
		}
	}

	defer scope.Close()

	ec2svc  := ec2.NewService(scope)

	// TODO worry about pointers. instance if exists returns *any* instance
	if scope.MachineStatus.InstanceID == nil {
		return false, nil
	}

	instance, err := ec2svc.InstanceIfExists(*scope.MachineStatus.InstanceID)
	if err != nil {
		return false, errors.Errorf("failed to retrieve instance: %+v", err)
	}

	if instance == nil {
		return false, nil
	}

	klog.Infof("Found instance for machine %q: %v", machine.Name, instance)

	switch instance.State {
	case v1alpha1.InstanceStateRunning:
		klog.Infof("Machine %s, %v is running", machine.Name, scope.MachineStatus.InstanceID)
	case v1alpha1.InstanceStatePending:
		klog.Infof("Machine %s, %v is pending", machine.Name, scope.MachineStatus.InstanceID)
	default:
		return false, nil
	}
	scope.MachineStatus.InstanceState = aws.String(string(instance.State))

	if err := a.reconcileLBAttachment(scope, machine, instance); err != nil {
		return true, err
	}

	return true, nil
}

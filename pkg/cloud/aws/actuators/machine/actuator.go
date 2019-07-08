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
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/klogr"
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

const (
	defaultTokenTTL = 10 * time.Minute
)

//+kubebuilder:rbac:groups=awsprovider.k8s.io,resources=awsmachineproviderconfigs;awsmachineproviderstatuses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.k8s.io,resources=machines;machines/status;machinedeployments;machinedeployments/status;machinesets;machinesets/status;machineclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=nodes;events,verbs=get;list;watch;create;update;patch;delete

// Actuator is responsible for performing machine reconciliation.
type Actuator struct {
	*deployer.Deployer

	coreClient    corev1.CoreV1Interface
	clusterClient client.ClusterV1alpha1Interface
	log           logr.Logger
}

// ActuatorParams holds parameter information for Actuator.
type ActuatorParams struct {
	CoreClient     corev1.CoreV1Interface
	ClusterClient  client.ClusterV1alpha1Interface
	LoggingContext string
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		Deployer:      deployer.New(deployer.Params{ScopeGetter: actuators.DefaultScopeGetter}),
		coreClient:    params.CoreClient,
		clusterClient: params.ClusterClient,
		log:           klogr.New().WithName(params.LoggingContext),
	}
}

// GetControlPlaneMachines retrieves all non-deleted control plane nodes from a MachineList
func GetControlPlaneMachines(machineList *clusterv1.MachineList) []*clusterv1.Machine {
	var cpm []*clusterv1.Machine
	for _, m := range machineList.Items {
		if m.DeletionTimestamp.IsZero() && m.Spec.Versions.ControlPlane != "" {
			cpm = append(cpm, m.DeepCopy())
		}
	}
	return cpm
}

// defining equality as name and namespace are equivalent and not checking any other fields.
func machinesEqual(m1 *clusterv1.Machine, m2 *clusterv1.Machine) bool {
	return m1.Name == m2.Name && m1.Namespace == m2.Namespace
}

// isNodeJoin determines if a machine, in scope, should join of the cluster.
// TODO: Make this thread safe kubernetes-sigs/cluster-api#925
// https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/745#discussion_r280506890
func (a *Actuator) isNodeJoin(scope *actuators.MachineScope, controlPlaneMachines []*clusterv1.Machine) (bool, error) {
	switch set := scope.Machine.ObjectMeta.Labels["set"]; set {
	case "node":
		// Worker machines, not part of the controlplane, will always join the cluster.
		return true, nil
	case "controlplane":
		// Controlplane machines will join the cluster if the cluster has an existing control plane.
		controlplaneExists := false
		var err error

		var sharedScope *actuators.MachineScope

		for _, cm := range controlPlaneMachines {
			if sharedScope == nil {
				sharedScope, err = actuators.NewMachineScope(actuators.MachineScopeParams{
					Machine:    cm,
					Cluster:    scope.Cluster,
					Client:     a.clusterClient,
					Logger:     a.log,
					AWSClients: scope.AWSClients,
				})
				if err != nil {
					return false, errors.Wrapf(err, "failed to create machine scope for %s/%s", cm.Namespace, cm.Name)
				}
			}
			sharedScope.Machine = cm
			ec2svc := ec2.NewService(sharedScope.Scope)

			controlplaneExists, err = ec2svc.MachineExists(sharedScope)
			if err != nil {
				return false, errors.Wrapf(err, "failed to verify existence of machine %s/%s", sharedScope.Machine.Namespace, sharedScope.Machine.Name)
			}

			if controlplaneExists {
				return controlplaneExists, err
			}
		}

		a.log.V(2).Info("Will machine join the controlplane", "machine-name", scope.Machine.Name, "machine-namespace", scope.Machine.Name, "should-join-control-plane", controlplaneExists)
		return controlplaneExists, err

	default:
		return false, errors.Errorf("Unknown value %q for label `set` on machine %q", set, scope.Machine.Name)
	}
}

// Create creates a machine and is invoked by the machine controller.
func (a *Actuator) Create(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	if cluster == nil {
		return errors.Errorf("missing cluster for machine %s/%s", machine.Namespace, machine.Name)
	}
	a.log.Info("Creating machine in cluster", "machine-name", machine.Name, "machine-namespace", machine.Namespace, "cluster-name", cluster.Name)

	scope, err := actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.clusterClient, Logger: a.log})
	if err != nil {
		return errors.Errorf("failed to create scope: %+v", err)
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope.Scope)

	clusterMachines, err := scope.MachineClient.List(v1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve machines in cluster %q", cluster.Name)
	}

	controlPlaneMachines := GetControlPlaneMachines(clusterMachines)

	isNodeJoin, err := a.isNodeJoin(scope, controlPlaneMachines)
	if err != nil {
		return errors.Wrapf(err, "failed to determine whether machine %q should join cluster %q", machine.Name, cluster.Name)
	}

	var bootstrapToken string
	if isNodeJoin {
		a.log.V(2).Info("Machine will join the cluster", "cluster", cluster.Name, "machine-name", machine.Name, "machine-namespace", machine.Namespace)
		coreClient, err := a.coreV1Client(cluster)
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve corev1 client for cluster %q", cluster.Name)
		}

		bootstrapToken, err = tokens.NewBootstrap(coreClient, defaultTokenTTL)
		if err != nil {
			return errors.Wrapf(err, "failed to create new bootstrap token")
		}
	}

	i, err := ec2svc.CreateOrGetMachine(scope, bootstrapToken)
	if err != nil {
		if awserrors.IsFailedDependency(errors.Cause(err)) {
			a.log.Error(err, "network not ready to launch instances yet")
			return &controllerError.RequeueAfterError{
				RequeueAfter: time.Minute,
			}
		}

		return errors.Errorf("failed to create or get machine: %+v", err)
	}

	scope.MachineStatus.InstanceID = &i.ID
	scope.MachineStatus.InstanceState = &i.State

	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}

	machine.Annotations["cluster-api-provider-aws"] = "true"

	if err := a.reconcileLBAttachment(scope, machine, i); err != nil {
		return errors.Errorf("failed to reconcile LB attachment: %+v", err)
	}
	a.log.Info("Create completed", "machine-name", machine.Name)
	return nil
}

func (a *Actuator) coreV1Client(cluster *clusterv1.Cluster) (corev1.CoreV1Interface, error) {
	controlPlaneDNSName, err := a.GetIP(cluster, nil)
	if err != nil {
		return nil, errors.Errorf("failed to retrieve controlplane (GetIP): %+v", err)
	}

	controlPlaneURL := fmt.Sprintf("https://%s:6443", controlPlaneDNSName)

	kubeConfig, err := a.GetKubeConfig(cluster, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve kubeconfig for cluster %q.", cluster.Name)
	}

	clientConfig, err := clientcmd.BuildConfigFromKubeconfigGetter(controlPlaneURL, func() (*clientcmdapi.Config, error) {
		return clientcmd.Load([]byte(kubeConfig))
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get client config for cluster at %q", controlPlaneURL)
	}

	return corev1.NewForConfig(clientConfig)
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
	if cluster == nil {
		return errors.Errorf("missing cluster for machine %s/%s", machine.Namespace, machine.Name)
	}
	a.log.Info("Deleting machine in cluster", "machine-name", machine.Name, "machine-namespace", machine.Namespace, "cluster-name", cluster.Name)

	scope, err := actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.clusterClient, Logger: a.log})
	if err != nil {
		return errors.Errorf("failed to create scope: %+v", err)
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope.Scope)

	instance, err := ec2svc.InstanceIfExists(scope.MachineStatus.InstanceID)
	if err != nil {
		return errors.Errorf("failed to get instance: %+v", err)
	}

	if instance == nil {
		instance, err = ec2svc.InstanceByTags(scope)
		if err != nil {
			return errors.Errorf("failed to query instance by tags: %+v", err)
		} else if instance == nil {
			// The machine hasn't been created yet
			a.log.V(3).Info("Instance is nil and therefore does not exist")
			return nil
		}
	}

	// Check the instance state. If it's already shutting down or terminated,
	// do nothing. Otherwise attempt to delete it.
	// This decision is based on the ec2-instance-lifecycle graph at
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	switch instance.State {
	case v1alpha1.InstanceStateShuttingDown, v1alpha1.InstanceStateTerminated:
		a.log.Info("Machine instance is shutting down or already terminated")
		return nil
	default:
		a.log.Info("Terminating machine")
		if err := ec2svc.TerminateInstance(instance.ID); err != nil {
			return errors.Errorf("failed to terminate instance: %+v", err)
		}
	}

	return nil
}

// isMachineOudated checks that no immutable fields have been updated in an
// Update request.
// Returns a slice of errors representing attempts to change immutable state
func (a *Actuator) isMachineOutdated(machineSpec *v1alpha1.AWSMachineProviderSpec, instance *v1alpha1.Instance) (errs []error) {
	// Instance Type
	if machineSpec.InstanceType != instance.Type {
		errs = append(errs, errors.Errorf("instance type cannot be mutated from %q to %q", instance.Type, machineSpec.InstanceType))
	}

	// IAM Profile
	if machineSpec.IAMInstanceProfile != instance.IAMProfile {
		errs = append(errs, errors.Errorf("instance IAM profile cannot be mutated from %q to %q", instance.IAMProfile, machineSpec.IAMInstanceProfile))
	}

	// SSH Key Name
	if machineSpec.KeyName != aws.StringValue(instance.KeyName) {
		errs = append(errs, errors.Errorf("SSH key name cannot be mutated from %q to %q", aws.StringValue(instance.KeyName), machineSpec.KeyName))
	}

	// Root Device Size
	if machineSpec.RootDeviceSize > 0 && machineSpec.RootDeviceSize != instance.RootDeviceSize {
		errs = append(errs, errors.Errorf("Root volume size cannot be mutated from %v to %v", instance.RootDeviceSize, machineSpec.RootDeviceSize))
	}

	// Subnet ID
	// machineSpec.Subnet is a *AWSResourceReference and could technically be
	// a *string, ARN or Filter. However, elsewhere in the code it is only used
	// as a *string, so do the same here.
	if machineSpec.Subnet != nil {
		if aws.StringValue(machineSpec.Subnet.ID) != instance.SubnetID {
			errs = append(errs, errors.Errorf("machine subnet ID cannot be mutated from %q to %q", instance.SubnetID, aws.StringValue(machineSpec.Subnet.ID)))
		}
	}

	// PublicIP check is a little more complicated as the machineConfig is a
	// simple bool indicating if the instance should have a public IP or not,
	// while the instanceDescription contains the public IP assigned to the
	// instance.
	// Work out whether the instance already has a public IP or not based on
	// the length of the PublicIP string. Anything >0 is assumed to mean it does
	// have a public IP.
	instanceHasPublicIP := false
	if len(aws.StringValue(instance.PublicIP)) > 0 {
		instanceHasPublicIP = true
	}

	if aws.BoolValue(machineSpec.PublicIP) != instanceHasPublicIP {
		errs = append(errs, errors.Errorf(`public IP setting cannot be mutated from "%v" to "%v"`, instanceHasPublicIP, aws.BoolValue(machineSpec.PublicIP)))
	}

	return errs
}

// Update updates a machine and is invoked by the Machine Controller.
// If the Update attempts to mutate any immutable state, the method will error
// and no updates will be performed.
func (a *Actuator) Update(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	if cluster == nil {
		return errors.Errorf("missing cluster for machine %s/%s", machine.Namespace, machine.Name)
	}

	a.log.Info("Updating machine in cluster", "machine-name", machine.Name, "machine-namespace", machine.Namespace, "cluster-name", cluster.Name)

	scope, err := actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.clusterClient, Logger: a.log})
	if err != nil {
		return errors.Errorf("failed to create scope: %+v", err)
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope.Scope)

	// Get the current instance description from AWS.
	instanceDescription, err := ec2svc.InstanceIfExists(scope.MachineStatus.InstanceID)
	if err != nil {
		return errors.Errorf("failed to get instance: %+v", err)
	}

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	if errs := a.isMachineOutdated(scope.MachineConfig, instanceDescription); len(errs) > 0 {
		return errors.Errorf("found attempt to change immutable state for machine %q: %+q", machine.Name, errs)
	}

	existingSecurityGroups, err := ec2svc.GetInstanceSecurityGroups(*scope.MachineStatus.InstanceID)
	if err != nil {
		return err
	}

	// Ensure that the security groups are correct.
	_, err = a.ensureSecurityGroups(
		ec2svc,
		scope,
		*scope.MachineStatus.InstanceID,
		scope.MachineConfig.AdditionalSecurityGroups,
		existingSecurityGroups,
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
	if cluster == nil {
		return false, errors.Errorf("missing cluster for machine %s/%s", machine.Namespace, machine.Name)
	}

	a.log.Info("Checking if machine exists in cluster", "machine-name", machine.Name, "machine-namespace", machine.Namespace, "cluster-name", cluster.Name)

	scope, err := actuators.NewMachineScope(actuators.MachineScopeParams{Machine: machine, Cluster: cluster, Client: a.clusterClient, Logger: a.log})
	if err != nil {
		return false, errors.Errorf("failed to create scope: %+v", err)
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope.Scope)

	// TODO worry about pointers. instance if exists returns *any* instance
	if scope.MachineStatus.InstanceID == nil {
		return false, nil
	}

	instance, err := ec2svc.InstanceIfExists(scope.MachineStatus.InstanceID)
	if err != nil {
		return false, errors.Errorf("failed to retrieve instance: %+v", err)
	}

	if instance == nil {
		return false, nil
	}

	a.log.Info("Found instance for machine", "machine-name", machine.Name, "machine-namespace", machine.Namespace, "instance", instance)

	switch instance.State {
	case v1alpha1.InstanceStateRunning:
		a.log.Info("Machine instance is running", "instance-id", *scope.MachineStatus.InstanceID)
	case v1alpha1.InstanceStatePending:
		a.log.Info("Machine instance is pending", "instance-id", *scope.MachineStatus.InstanceID)
	default:
		return false, nil
	}

	scope.MachineStatus.InstanceState = &instance.State

	if err := a.reconcileLBAttachment(scope, machine, instance); err != nil {
		return true, err
	}

	if machine.Spec.ProviderID == nil || *machine.Spec.ProviderID == "" {
		providerID := fmt.Sprintf("aws:////%s", *scope.MachineStatus.InstanceID)
		scope.Machine.Spec.ProviderID = &providerID
	}

	return true, nil
}

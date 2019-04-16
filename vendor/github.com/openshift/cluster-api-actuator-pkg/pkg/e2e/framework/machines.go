package framework

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"

	"github.com/openshift/cluster-api-actuator-pkg/pkg/types"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	controllernode "github.com/openshift/cluster-api/pkg/controller/node"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (f *Framework) DeleteMachineAndWait(machine *machinev1beta1.Machine, client types.CloudProviderClient) {
	f.By(fmt.Sprintf("Deleting %q machine", machine.Name))
	err := f.CAPIClient.MachineV1beta1().Machines(machine.Namespace).Delete(machine.Name, &metav1.DeleteOptions{})
	f.IgnoreNotFoundErr(err)

	// Verify the testing machine has been destroyed
	f.By("Verify instance is terminated")
	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.MachineV1beta1().Machines(machine.Namespace).Get(machine.Name, metav1.GetOptions{})
		if err == nil {
			glog.V(2).Info("Waiting for machine to be deleted")
			return false, nil
		}
		if strings.Contains(err.Error(), "not found") {
			return true, nil
		}
		return false, nil
	})
	f.IgnoreNotFoundErr(err)

	if client != nil {
		err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
			glog.V(2).Info("Waiting for instance to be terminated")
			runningInstances, err := client.GetRunningInstances(machine)
			if err != nil {
				return false, fmt.Errorf("unable to get running instances: %v", err)
			}
			runningInstancesLen := len(runningInstances)
			if runningInstancesLen > 0 {
				glog.V(2).Info("Machine is running")
				return false, nil
			}
			return true, nil
		})
		f.ErrNotExpected(err)
	}
}

func (f *Framework) waitForMachineToRun(machine *machinev1beta1.Machine, client types.CloudProviderClient) {
	f.By(fmt.Sprintf("Waiting for %q machine", machine.Name))
	// Verify machine has been deployed
	err := wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		if _, err := f.CAPIClient.MachineV1beta1().Machines(machine.Namespace).Get(machine.Name, metav1.GetOptions{}); err != nil {
			glog.V(2).Infof("Waiting for '%v/%v' machine to be created", machine.Namespace, machine.Name)
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)

	f.By("Verify machine's underlying instance is running")
	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		glog.V(2).Info("Waiting for instance to come up")
		runningInstances, err := client.GetRunningInstances(machine)
		if err != nil {
			glog.V(2).Infof("unable to get running instances: %v", err)
			return false, nil
		}
		runningInstancesLen := len(runningInstances)
		if runningInstancesLen == 1 {
			glog.V(2).Info("Machine is running")
			return true, nil
		}
		if runningInstancesLen > 1 {
			return false, fmt.Errorf("Found %q instances instead of one", runningInstancesLen)
		}
		return false, nil
	})
	f.ErrNotExpected(err)
}

func (f *Framework) waitForMachineToTerminate(machine *machinev1beta1.Machine, client types.CloudProviderClient) error {
	f.By("Verify machine's underlying instance is not running")
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		glog.V(2).Info("Waiting for instance to terminate")
		runningInstances, err := client.GetRunningInstances(machine)
		if err != nil {
			glog.V(2).Infof("unable to get running instances from cloud provider: %v", err)
			return false, nil
		}
		runningInstancesLen := len(runningInstances)
		if runningInstancesLen > 1 {
			return false, fmt.Errorf("Found %q running instances for %q", runningInstancesLen, machine.Name)
		}
		return true, nil
	})
	// We need to allow to follow
	if err != nil {
		glog.V(2).Infof("unable to wait for instance(s) to terminate: %v", err)
		return err
	}

	f.By(fmt.Sprintf("Waiting for %q machine object to be deleted", machine.Name))
	// Verify machine has been deployed
	err = wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		if _, err := f.CAPIClient.MachineV1beta1().Machines(machine.Namespace).Get(machine.Name, metav1.GetOptions{}); err != nil {
			glog.V(2).Infof("err: %v", err)
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		glog.V(2).Infof("unable to wait for machine to get deleted: %v", err)
		return err
	}
	return nil
}

func (f *Framework) CreateMachineAndWait(machine *machinev1beta1.Machine, client types.CloudProviderClient) {
	f.By(fmt.Sprintf("Creating %q machine", machine.Name))
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.MachineV1beta1().Machines(machine.Namespace).Create(machine)
		if err != nil {
			glog.V(2).Infof("can't create machine: %v", err)
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)

	f.waitForMachineToRun(machine, client)
}

func (f *Framework) CreateMachineSetAndWait(machineset *machinev1beta1.MachineSet, client types.CloudProviderClient) {
	f.By(fmt.Sprintf("Creating %q machineset", machineset.Name))
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.MachineV1beta1().MachineSets(machineset.Namespace).Create(machineset)
		if err != nil {
			glog.V(2).Infof("can't create machineset: %v", err)
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)

	// Verify machineset has been deployed
	err = wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		if _, err := f.CAPIClient.MachineV1beta1().MachineSets(machineset.Namespace).Get(machineset.Name, metav1.GetOptions{}); err != nil {
			glog.V(2).Infof("Waiting for machineset to be created: %v", err)
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)

	f.By("Verify machineset's underlying instances is running")
	machines, err := f.CAPIClient.MachineV1beta1().Machines(machineset.Namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(machineset.Spec.Selector.MatchLabels).String(),
	})
	f.ErrNotExpected(err)

	for _, machine := range machines.Items {
		f.waitForMachineToRun(&machine, client)
	}
}

func (f *Framework) DeleteMachineSetAndWait(machineset *machinev1beta1.MachineSet, client types.CloudProviderClient) error {
	f.By(fmt.Sprintf("Get all %q machineset's machines", machineset.Name))
	machines, err := f.CAPIClient.MachineV1beta1().Machines(machineset.Namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(machineset.Spec.Selector.MatchLabels).String(),
	})
	if err != nil {
		return err
	}

	f.By(fmt.Sprintf("Deleting %q machineset", machineset.Name))
	err = f.CAPIClient.MachineV1beta1().MachineSets(machineset.Namespace).Delete(machineset.Name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	f.By("Waiting for all machines to be deleted")
	for _, machine := range machines.Items {
		f.waitForMachineToTerminate(&machine, client)
		// TODO(jchaloup): run it one more time depending on the error returned
	}

	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.MachineV1beta1().MachineSets(machineset.Namespace).Get(machineset.Name, metav1.GetOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return true, nil
			}
			return false, nil
		}
		return true, nil
	})
	return err
}

func (f *Framework) WaitForNodesToGetReady(count int) error {
	return wait.Poll(PollNodeInterval, PoolNodesReadyTimeout, func() (bool, error) {
		items, err := f.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			glog.V(2).Infof("unable to list nodes: %v", err)
			return false, nil
		}
		if len(items.Items) < count {
			glog.V(2).Infof("Waiting for %v nodes to come up, have %v", count, len(items.Items))
			return false, nil
		}
		allNodesReady := true
		for _, node := range items.Items {
			for _, condition := range node.Status.Conditions {
				if condition.Type != apiv1.NodeReady {
					continue
				}
				if condition.Status != apiv1.ConditionTrue {
					glog.V(2).Infof("Node %q not ready", node.Name)
					allNodesReady = false
				} else {
					glog.V(2).Infof("Node %q is ready", node.Name)
				}
				break
			}
		}

		if !allNodesReady {
			return false, nil
		}

		return true, nil
	})
}

// GetMachineFromNode returns the machine referenced by the "controllernode.MachineAnnotationKey" annotation in the given node
func GetMachineFromNode(client runtimeclient.Client, node *corev1.Node) (*mapiv1beta1.Machine, error) {
	machineNamespaceKey, ok := node.Annotations[controllernode.MachineAnnotationKey]
	if !ok {
		return nil, fmt.Errorf("node %q does not have a MachineAnnotationKey %q", node.Name, controllernode.MachineAnnotationKey)
	}
	namespace, machineName, err := cache.SplitMetaNamespaceKey(machineNamespaceKey)
	if err != nil {
		return nil, fmt.Errorf("machine annotation format is incorrect %v: %v", machineNamespaceKey, err)
	}

	if namespace != TestContext.MachineApiNamespace {
		return nil, fmt.Errorf("Machine %q is forbidden to live outside of default %v namespace", machineNamespaceKey, TestContext.MachineApiNamespace)
	}

	machine, err := GetMachine(context.TODO(), client, machineName)
	if err != nil {
		return nil, fmt.Errorf("error querying api for machine object: %v", err)
	}

	return machine, nil
}

func ReadKubeconfigFromServer(sshConfig *SSHConfig) (string, error) {
	client, err := createSSHClient(sshConfig)
	if err != nil {
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	var se bytes.Buffer
	session.Stdout = &b
	session.Stderr = &se
	if err := session.Run("sudo cat /root/.kube/config"); err != nil {
		return "", fmt.Errorf("failed to collect kubeconfig: %v, %v", err.Error(), se.String())
	}
	return b.String(), nil
}

func (f *Framework) GetMasterMachineRestConfig(masterMachine *machinev1beta1.Machine, client types.CloudProviderClient) (*rest.Config, error) {
	var masterPublicDNSName string
	err := wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		var err error
		masterPublicDNSName, err = client.GetPublicDNSName(masterMachine)
		if err != nil {
			glog.V(2).Infof("Unable to collect master public DNS name: %v", err)
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	var masterKubeconfig string
	err = wait.Poll(PollInterval, PoolKubeConfigTimeout, func() (bool, error) {
		glog.V(2).Infof("Pulling kubeconfig from %v:8443", masterPublicDNSName)
		var err error
		masterKubeconfig, err = ReadKubeconfigFromServer(&SSHConfig{
			User: f.SSH.User,
			Host: masterPublicDNSName,
			Key:  f.SSH.Key,
		})
		if err != nil {
			glog.V(2).Infof("Unable to pull kubeconfig: %v", err)
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("Master running on https://" + masterPublicDNSName + ":8443")

	config, err := clientcmd.Load([]byte(masterKubeconfig))
	if err != nil {
		return nil, err
	}

	config.Clusters["kubernetes"].Server = "https://" + masterPublicDNSName + ":8443"
	return clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
}

type machineToDelete struct {
	machine   *machinev1beta1.Machine
	framework *Framework
	client    types.CloudProviderClient
}

type machinesetToDelete struct {
	machineset *machinev1beta1.MachineSet
	framework  *Framework
	client     types.CloudProviderClient
}

type MachinesToDelete struct {
	machines    []machineToDelete
	machinesets []machinesetToDelete
}

func InitMachinesToDelete() *MachinesToDelete {
	return &MachinesToDelete{
		machines:    make([]machineToDelete, 0),
		machinesets: make([]machinesetToDelete, 0),
	}
}

func (m *MachinesToDelete) AddMachine(machine *machinev1beta1.Machine, framework *Framework, client types.CloudProviderClient) {
	m.machines = append([]machineToDelete{{machine: machine, framework: framework, client: client}}, m.machines...)
}

func (m *MachinesToDelete) AddMachineSet(machineset *machinev1beta1.MachineSet, framework *Framework, client types.CloudProviderClient) {
	m.machinesets = append([]machinesetToDelete{{machineset: machineset, framework: framework, client: client}}, m.machinesets...)
}

func (m *MachinesToDelete) Delete() {
	for _, item := range m.machinesets {
		item.framework.DeleteMachineSetAndWait(item.machineset, item.client)
	}

	m.machinesets = make([]machinesetToDelete, 0)

	for _, item := range m.machines {
		item.framework.DeleteMachineAndWait(item.machine, item.client)
	}

	m.machines = make([]machineToDelete, 0)
}

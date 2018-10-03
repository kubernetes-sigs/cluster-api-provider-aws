package framework

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/prometheus/common/log"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CloudProviderClient interface {
	// Get running instances (of a given cloud provider) managed by the machine object
	GetRunningInstances(machine *clusterv1alpha1.Machine) ([]interface{}, error)
	// Get running instance public DNS name
	GetPublicDNSName(machine *clusterv1alpha1.Machine) (string, error)
	// Get private IP
	GetPrivateIP(machine *clusterv1alpha1.Machine) (string, error)
}

func (f *Framework) DeleteMachineAndWait(machine *clusterv1alpha1.Machine, client CloudProviderClient) {
	By(fmt.Sprintf("Deleting %q machine", machine.Name))
	err := f.CAPIClient.ClusterV1alpha1().Machines(machine.Namespace).Delete(machine.Name, &metav1.DeleteOptions{})
	IgnoreNotFoundErr(err)

	// Verify the testing machine has been destroyed
	By("Verify instance is terminated")
	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().Machines(machine.Namespace).Get(machine.Name, metav1.GetOptions{})
		if err == nil {
			log.Info("Waiting for machine to be deleted")
			return false, nil
		}
		if strings.Contains(err.Error(), "not found") {
			return true, nil
		}
		return false, nil
	})
	IgnoreNotFoundErr(err)

	if client != nil {
		err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
			log.Info("Waiting for instance to be terminated")
			runningInstances, err := client.GetRunningInstances(machine)
			if err != nil {
				return false, fmt.Errorf("unable to get running instances from aws: %v", err)
			}
			runningInstancesLen := len(runningInstances)
			if runningInstancesLen > 0 {
				log.Info("Machine is running")
				return false, nil
			}
			return true, nil
		})
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Framework) waitForMachineToRun(machine *clusterv1alpha1.Machine, client CloudProviderClient) {
	By(fmt.Sprintf("Waiting for %q machine", machine.Name))
	// Verify machine has been deployed
	err := wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		if _, err := f.CAPIClient.ClusterV1alpha1().Machines(machine.Namespace).Get(machine.Name, metav1.GetOptions{}); err != nil {
			log.Infof("Waiting for '%v/%v' machine to be created", machine.Namespace, machine.Name)
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())

	By("Verify machine's underlying instance is running")
	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		log.Info("Waiting for instance to come up")
		runningInstances, err := client.GetRunningInstances(machine)
		if err != nil {
			return false, fmt.Errorf("unable to get running instances from aws: %v", err)
		}
		runningInstancesLen := len(runningInstances)
		if runningInstancesLen == 1 {
			log.Info("Machine is running")
			return true, nil
		}
		if runningInstancesLen > 1 {
			return false, fmt.Errorf("Found %q instances instead of one", runningInstancesLen)
		}
		return false, nil
	})
	Expect(err).NotTo(HaveOccurred())
}

func (f *Framework) waitForMachineToTerminate(machine *clusterv1alpha1.Machine, client CloudProviderClient) error {
	By("Verify machine's underlying instance is not running")
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		log.Info("Waiting for instance to terminate")
		runningInstances, err := client.GetRunningInstances(machine)
		if err != nil {
			return false, fmt.Errorf("unable to get running instances from cloud provider: %v", err)
		}
		runningInstancesLen := len(runningInstances)
		if runningInstancesLen > 1 {
			return false, fmt.Errorf("Found %q running instances for %q", runningInstancesLen, machine.Name)
		}
		return true, nil
	})
	// We need to allow to follow
	if err != nil {
		log.Info("unable to wait for instance(s) to terminate: %v", err)
		return err
	}

	By(fmt.Sprintf("Waiting for %q machine object to be deleted", machine.Name))
	// Verify machine has been deployed
	err = wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		if _, err := f.CAPIClient.ClusterV1alpha1().Machines(machine.Namespace).Get(machine.Name, metav1.GetOptions{}); err != nil {
			log.Infof("err: %v", err)
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		log.Info("unable to wait for machine to get deleted: %v", err)
		return err
	}
	return nil
}

func (f *Framework) CreateMachineAndWait(machine *clusterv1alpha1.Machine, client CloudProviderClient) {
	By(fmt.Sprintf("Creating %q machine", machine.Name))
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().Machines(machine.Namespace).Create(machine)
		if err != nil {
			log.Infof("can't create machine: %v", err)
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())

	f.waitForMachineToRun(machine, client)
}

func (f *Framework) CreateMachineSetAndWait(machineset *clusterv1alpha1.MachineSet, client CloudProviderClient) {
	By(fmt.Sprintf("Creating %q machineset", machineset.Name))
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().MachineSets(machineset.Namespace).Create(machineset)
		if err != nil {
			log.Infof("can't create machineset: %v", err)
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())

	// Verify machineset has been deployed
	err = wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		if _, err := f.CAPIClient.ClusterV1alpha1().MachineSets(machineset.Namespace).Get(machineset.Name, metav1.GetOptions{}); err != nil {
			log.Infof("Waiting for machineset to be created: %v", err)
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())

	By("Verify machineset's underlying instances is running")
	machines, err := f.CAPIClient.ClusterV1alpha1().Machines(machineset.Namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(machineset.Spec.Selector.MatchLabels).String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, machine := range machines.Items {
		f.waitForMachineToRun(&machine, client)
	}
}

func (f *Framework) DeleteMachineSetAndWait(machineset *clusterv1alpha1.MachineSet, client CloudProviderClient) error {
	By("Get all machineset's machines")
	machines, err := f.CAPIClient.ClusterV1alpha1().Machines(machineset.Namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(machineset.Spec.Selector.MatchLabels).String(),
	})
	if err != nil {
		return err
	}

	By(fmt.Sprintf("Deleting %q machineset", machineset.Name))
	err = f.CAPIClient.ClusterV1alpha1().MachineSets(machineset.Namespace).Delete(machineset.Name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	By("Waiting for all machines to be deleted")
	for _, machine := range machines.Items {
		f.waitForMachineToTerminate(&machine, client)
		// TODO(jchaloup): run it one more time depending on the error returned
	}

	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().MachineSets(machineset.Namespace).Get(machineset.Name, metav1.GetOptions{})
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

func (f *Framework) WaitForNodesToGetReady() {
	err := wait.Poll(PollNodeInterval, PoolNodesReadyTimeout, func() (bool, error) {
		items, err := f.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			return false, fmt.Errorf("unable to list nodes: %v", err)
		}
		if len(items.Items) < 2 {
			log.Infof("Waiting for both nodes to come up, have %v", len(items.Items))
			return false, nil
		}
		allNodesReady := true
		for _, node := range items.Items {
			for _, condition := range node.Status.Conditions {
				if condition.Type != apiv1.NodeReady {
					continue
				}
				if condition.Status != apiv1.ConditionTrue {
					log.Infof("Node %q not ready", node.Name)
					allNodesReady = false
				} else {
					log.Infof("Node %q is ready", node.Name)
				}
				break
			}
		}

		if !allNodesReady {
			return false, nil
		}

		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())
}

func ReadKubeconfigFromServer(sshuser, sshhost, prkey string) (string, error) {
	fp, err := os.Open(prkey)
	if err != nil {
		return "", fmt.Errorf("unable to open id_rsa file: %v", err)
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return "", fmt.Errorf("unable to read id_rsa file: %v", err)
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: sshuser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", sshhost+":22", config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: " + err.Error())
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

func (f *Framework) GetMasterMachineRestConfig(masterMachine *clusterv1alpha1.Machine, client CloudProviderClient) (*rest.Config, error) {
	var masterPublicDNSName string
	err := wait.Poll(PollInterval, TimeoutPoolMachineRunningInterval, func() (bool, error) {
		var err error
		masterPublicDNSName, err = client.GetPublicDNSName(masterMachine)
		if err != nil {
			log.Infof("Unable to collect master public DNS name: %v", err)
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	var masterKubeconfig string
	err = wait.Poll(PollInterval, PoolKubeConfigTimeout, func() (bool, error) {
		log.Infof("Pulling kubeconfig from %v:8443", masterPublicDNSName)
		var err error
		masterKubeconfig, err = ReadKubeconfigFromServer("ec2-user", masterPublicDNSName, f.SSHKey)
		if err != nil {
			log.Infof("Unable to pull kubeconfig: %v", err)
			return false, nil
		}

		return true, nil
	})

	log.Infof("Master running on https://" + masterPublicDNSName + ":8443")

	config, err := clientcmd.Load([]byte(masterKubeconfig))
	if err != nil {
		return nil, err
	}
	config.Clusters["kubernetes"].Server = "https://" + masterPublicDNSName + ":8443"
	return clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
}

type machineToDelete struct {
	machine   *clusterv1alpha1.Machine
	framework *Framework
	client    CloudProviderClient
}

type machinesetToDelete struct {
	machineset *clusterv1alpha1.MachineSet
	framework  *Framework
	client     CloudProviderClient
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

func (m *MachinesToDelete) AddMachine(machine *clusterv1alpha1.Machine, framework *Framework, client CloudProviderClient) {
	m.machines = append(m.machines, machineToDelete{machine: machine, framework: framework, client: client})
}

func (m *MachinesToDelete) AddMachineSet(machineset *clusterv1alpha1.MachineSet, framework *Framework, client CloudProviderClient) {
	m.machinesets = append(m.machinesets, machinesetToDelete{machineset: machineset, framework: framework, client: client})
}

func (m *MachinesToDelete) Delete() {
	for _, item := range m.machines {
		item.framework.DeleteMachineAndWait(item.machine, item.client)
	}

	for _, item := range m.machinesets {
		item.framework.DeleteMachineSetAndWait(item.machineset, item.client)
	}
}

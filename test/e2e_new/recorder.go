// +build e2e

/*
Copyright 2020 The Kubernetes Authors.

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

package e2e_new

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"encoding/csv"
	"encoding/json"

	tablecsv "github.com/frictionlessdata/tableschema-go/csv"
	"github.com/frictionlessdata/tableschema-go/schema"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	corev1 "k8s.io/api/core/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	kcpv1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1alpha3"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	TestSchema = schema.Schema{
		Fields: []schema.Field{
			{
				Name:   "Time",
				Type:   schema.DateTimeType,
				Format: "%Y-%m-%dT%H:%M:%S.%f+%z",
			},
			{
				Name: "ClusterName",
				Type: schema.StringType,
			},
			{
				Name: "ClusterControlPlaneReady",
				Type: schema.BooleanType,
			},
			{
				Name: "ClusterControlPlaneInitialized",
				Type: schema.BooleanType,
			},
			{
				Name: "ClusterInfrastructureReady",
				Type: schema.BooleanType,
			},
			{
				Name: "ClusterReady",
				Type: schema.BooleanType,
			},
			{
				Name: "ClusterWaitingForControlPlane",
				Type: schema.BooleanType,
			},
			{
				Name: "ClusterWaitingForInfrastructure",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPMachinesReady",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPCertificatesAvailable",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPCertificateGenerationFailed",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPAvailable",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPWaitingForKubeadmInit",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPMachinesSpecUpToDate",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPRollingUpdateInProgress",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPReplicas",
				Type: schema.IntegerType,
			},
			{
				Name: "KCPUpdatedReplicas",
				Type: schema.IntegerType,
			},
			{
				Name: "KCPReadyReplicas",
				Type: schema.IntegerType,
			},
			{
				Name: "KCPUnavailableReplicas",
				Type: schema.IntegerType,
			},
			{
				Name: "KCPInitialized",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPReady",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPResized",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPScalingUp",
				Type: schema.BooleanType,
			},
			{
				Name: "KCPScalingDown",
				Type: schema.BooleanType,
			},
			{
				Name: "LoadBalancerFailed",
				Type: schema.BooleanType,
			},
			{
				Name: "LoadBalancerReady",
				Type: schema.BooleanType,
			},
			{
				Name: "TotalAWSMachines",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalAWSMachinesWaitingForBootstrapData",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalAWSMachinesWaitingForClusterInfrastructure",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalControlPlaneAWSMachines",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesELBAttached",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesELBAttachFailed",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesELBDetachFailed",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesNotFound",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesNotReady",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesProvisionFailed",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesProvisionStartedButFailed",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesReady",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesStopped",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalInstancesTerminated",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalMachines",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalMachinesBootstrapReady",
				Type: schema.IntegerType,
			},
			{
				Name: "TotalMachinesWaitingForDataSecret",
				Type: schema.IntegerType,
			},
			{
				Name: "WaitingForDNSName",
				Type: schema.BooleanType,
			},
			{
				Name: "WaitingForDNSResolution",
				Type: schema.BooleanType,
			},
		},
	}
)

type TestRecord struct {
	ClusterControlPlaneReady                        bool
	ClusterControlPlaneInitialized                  bool
	ClusterInfrastructureReady                      bool
	ClusterName                                     string
	ClusterReady                                    bool
	ClusterWaitingForControlPlane                   bool
	ClusterWaitingForInfrastructure                 bool
	KCPMachinesReady                                bool
	KCPCertificatesAvailable                        bool
	KCPCertificateGenerationFailed                  bool
	KCPAvailable                                    bool
	KCPWaitingForKubeadmInit                        bool
	KCPMachinesSpecUpToDate                         bool
	KCPRollingUpdateInProgress                      bool
	KCPReplicas                                     int
	KCPUpdatedReplicas                              int
	KCPReadyReplicas                                int
	KCPUnavailableReplicas                          int
	KCPInitialized                                  bool
	KCPReady                                        bool
	KCPResized                                      bool
	KCPScalingUp                                    bool
	KCPScalingDown                                  bool
	LoadBalancerFailed                              bool
	LoadBalancerReady                               bool
	Time                                            time.Time
	TotalAWSMachines                                int
	TotalAWSMachinesWaitingForBootstrapData         int
	TotalAWSMachinesWaitingForClusterInfrastructure int
	TotalControlPlaneAWSMachines                    int
	TotalInstancesELBAttached                       int
	TotalInstancesELBAttachFailed                   int
	TotalInstancesELBDetachFailed                   int
	TotalInstancesNotFound                          int
	TotalInstancesNotReady                          int
	TotalInstancesProvisionFailed                   int
	TotalInstancesProvisionStartedButFailed         int
	TotalInstancesReady                             int
	TotalInstancesStopped                           int
	TotalInstancesTerminated                        int
	TotalMachines                                   int
	TotalMachinesBootstrapReady                     int
	TotalMachinesWaitingForDataSecret               int
	WaitingForDNSName                               bool
	WaitingForDNSResolution                         bool
}

func newCsvFile() (*csv.Writer, *os.File, *os.File, error) {
	statisticsDir := path.Join(artifactFolder, "statistics")
	if err := os.MkdirAll(statisticsDir, 0o750); err != nil {
		return nil, nil, nil, err
	}
	tsLogFileName := filepath.Join(statisticsDir, fmt.Sprintf("time-series.%d.log", config.GinkgoConfig.ParallelNode))
	tsLogFile, err := os.OpenFile(tsLogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, nil, err
	}
	csvFileName := fmt.Sprintf("time-series.%d.csv", config.GinkgoConfig.ParallelNode)
	csvFile, err := os.Create(path.Join(statisticsDir, csvFileName))
	if err != nil {
		return nil, nil, nil, err
	}
	csvWriter := tablecsv.NewWriter(csvFile)

	descriptor := map[string]interface{}{
		"resources": []interface{}{
			map[string]interface{}{
				"name":    "time-series",
				"path":    csvFileName,
				"format":  "csv",
				"profile": "tabular-data-resource",
				"schema":  TestSchema,
			},
		},
	}
	descriptorJson, err := json.MarshalIndent(descriptor, "", "  ")
	if err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Unable to create time series descriptor: %v\n", err)
	}
	if err := ioutil.WriteFile(path.Join(statisticsDir, fmt.Sprintf("time-series-datatable.%d.json", config.GinkgoConfig.ParallelNode)), descriptorJson, 0o640); err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Unable to create time series descriptor: %v\n", err)
	}
	fieldNames := make([]string, len(TestSchema.Fields))
	for i := range TestSchema.Fields {
		fieldNames[i] = TestSchema.Fields[i].Name
	}
	if err := csvWriter.Write(fieldNames); err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Unable to write fields: %v\n", err)
	}
	return csvWriter, csvFile, tsLogFile, nil
}

func recordAllClusters(ctx context.Context, logF io.Writer, w *csv.Writer, clusterProxy framework.ClusterProxy, artifactFolder string, namespace *corev1.Namespace) {
	clusterList := &clusterv1.ClusterList{}
	timestamp := time.Now().UTC()
	mgmtClient := clusterProxy.GetClient()
	if err := mgmtClient.List(ctx, clusterList, client.InNamespace(namespace.Name)); err != nil {
		fmt.Fprintf(logF, "Couldn't list clusters: %v\n", err)
	}
	for i := range clusterList.Items {
		cluster := clusterList.Items[i]
		if err := recordCluster(ctx, logF, w, clusterProxy, cluster.Namespace, cluster.Name, timestamp); err != nil {
			fmt.Fprintf(logF, "Couldn't record cluster: %v\n", err)
		}
	}
}

func recordMachine(machine *clusterv1.Machine, record *TestRecord) {
	if instanceReadyCond := conditions.Get(machine, infrav1.InstanceReadyCondition); instanceReadyCond != nil {
		if instanceReadyCond.Status == corev1.ConditionTrue {
			record.TotalInstancesReady++
		} else {
			record.TotalInstancesNotReady++
			switch instanceReadyCond.Reason {
			case infrav1.InstanceStoppedReason:
				record.TotalInstancesStopped++
			case infrav1.InstanceTerminatedReason:
				record.TotalInstancesTerminated++
			}
		}
	}
	if instanceELBCond := conditions.Get(machine, infrav1.ELBAttachedCondition); instanceELBCond != nil {
		if instanceELBCond.Status == corev1.ConditionTrue {
			record.TotalInstancesELBAttached++
		} else {
			switch instanceELBCond.Reason {
			case infrav1.ELBAttachFailedReason:
				record.TotalInstancesELBAttachFailed++
			case infrav1.ELBDetachFailedReason:
				record.TotalInstancesELBDetachFailed++
			}
		}
	}
}

func saveRecord(f io.Writer, w *csv.Writer, record *TestRecord) {
	row, err := TestSchema.UncastRow(record)
	if err != nil {
		fmt.Fprintf(f, "Unable to uncast row: %v\n", err)
	}
	row[0] = record.Time.Format("2006-01-02T15:04:05.999999+00:00")
	if err := w.Write(row); err != nil {
		fmt.Fprintf(f, "Unable to write row: %v\n", err)
	}
	w.Flush()
}

func recordCluster(ctx context.Context, f io.Writer, w *csv.Writer, clusterProxy framework.ClusterProxy, namespace, name string, timestamp time.Time) error {
	record := TestRecord{
		ClusterName: name,
		Time:        timestamp,
	}

	defer saveRecord(f, w, &record)
	cluster := &clusterv1.Cluster{}
	if err := clusterProxy.GetClient().Get(ctx, apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, cluster); err != nil {
		return err
	}
	recordClusterResource(cluster, &record)
	infraClusterRef := cluster.Spec.InfrastructureRef
	if infraClusterRef != nil {
		infraCluster := &infrav1.AWSCluster{}
		if err := clusterProxy.GetClient().Get(ctx, apimachinerytypes.NamespacedName{Namespace: infraClusterRef.Namespace, Name: infraClusterRef.Name}, infraCluster); err != nil {
			return err
		}
		recordInfraCluster(infraCluster, &record)
	}
	kcpRef := cluster.Spec.ControlPlaneRef
	if kcpRef != nil {
		kcp := &kcpv1.KubeadmControlPlane{}
		if err := clusterProxy.GetClient().Get(ctx, apimachinerytypes.NamespacedName{Namespace: kcpRef.Namespace, Name: kcpRef.Name}, kcp); err != nil {
			return err
		}
		recordKCP(kcp, &record)
	}
	machines, err := machinesForCluster(ctx, clusterProxy, namespace, name)
	if err != nil {
		return err
	}
	for i := range machines {
		record.TotalMachines++
		machine := machines[i]
		recordMachine(&machine, &record)
		infraRef := machine.Spec.InfrastructureRef
		if infraRef.GroupVersionKind().Kind != "AWSMachine" {
			continue
		}
		record.TotalAWSMachines++
		infraMachine := &infrav1.AWSMachine{}
		if err := clusterProxy.GetClient().Get(ctx, apimachinerytypes.NamespacedName{Namespace: infraRef.Namespace, Name: infraRef.Name}, infraMachine); err != nil {
			return err
		}
		recordInfraMachine(infraMachine, &record)
	}

	return nil
}

func recordInfraMachine(infraMachine *infrav1.AWSMachine, record *TestRecord) {
	if readyCond := conditions.Get(infraMachine, infrav1.InstanceReadyCondition); readyCond != nil {
		if readyCond.Status == corev1.ConditionTrue {
			record.TotalInstancesReady++
		} else {
			switch readyCond.Reason {
			case infrav1.InstanceNotFoundReason:
				record.TotalInstancesNotFound++
			case infrav1.InstanceNotReadyReason:
				record.TotalInstancesNotReady++
			case infrav1.InstanceProvisionFailedReason:
				record.TotalInstancesProvisionFailed++
			case infrav1.InstanceStoppedReason:
				record.TotalInstancesStopped++
			case infrav1.InstanceTerminatedReason:
				record.TotalInstancesTerminated++
			}
		}
	}
	if elbAttachedCond := conditions.Get(infraMachine, infrav1.ELBAttachedCondition); elbAttachedCond != nil {
		if elbAttachedCond.Status == corev1.ConditionTrue {
			record.TotalInstancesELBAttached++
		} else {
			switch elbAttachedCond.Reason {
			case infrav1.ELBAttachFailedReason:
				record.TotalInstancesELBAttachFailed++
			case infrav1.ELBDetachFailedReason:
				record.TotalInstancesELBDetachFailed++
			}
		}
	}
}

func machinesForCluster(ctx context.Context, clusterProxy framework.ClusterProxy, namespace, name string) ([]clusterv1.Machine, error) {
	labels := map[string]string{clusterv1.ClusterLabelName: name}
	machineList := &clusterv1.MachineList{}
	mgmtClient := clusterProxy.GetClient()
	if err := mgmtClient.List(ctx, machineList, client.InNamespace(namespace), client.MatchingLabels(labels)); err != nil {
		return nil, err
	}
	return machineList.Items, nil
}

func recordKCP(kcp *kcpv1.KubeadmControlPlane, record *TestRecord) {
	status := kcp.Status
	record.KCPReadyReplicas = int(status.ReadyReplicas)
	record.KCPReplicas = int(status.Replicas)
	record.KCPUnavailableReplicas = int(status.UnavailableReplicas)
	record.KCPUpdatedReplicas = int(status.UpdatedReplicas)
	if machinesReadyCondition := conditions.Get(kcp, kcpv1.MachinesReadyCondition); machinesReadyCondition != nil {
		if machinesReadyCondition.Status == corev1.ConditionTrue {
			record.KCPMachinesReady = true
		}
	}
	if certificatesAvailableCond := conditions.Get(kcp, kcpv1.CertificatesAvailableCondition); certificatesAvailableCond != nil {
		if certificatesAvailableCond.Status == corev1.ConditionTrue {
			record.KCPCertificatesAvailable = true
		} else if certificatesAvailableCond.Reason == kcpv1.CertificatesGenerationFailedReason {
			record.KCPCertificateGenerationFailed = true
		}
	}
	if availableCond := conditions.Get(kcp, kcpv1.AvailableCondition); availableCond != nil {
		if availableCond.Status == corev1.ConditionTrue {
			record.KCPAvailable = true
		} else if availableCond.Reason == kcpv1.WaitingForKubeadmInitReason {
			record.KCPWaitingForKubeadmInit = true
		}
	}
	if machineSpecCondition := conditions.Get(kcp, kcpv1.MachinesSpecUpToDateCondition); machineSpecCondition != nil {
		if machineSpecCondition.Status == corev1.ConditionTrue {
			record.KCPMachinesSpecUpToDate = true
		} else if machineSpecCondition.Reason == kcpv1.RollingUpdateInProgressReason {
			record.KCPRollingUpdateInProgress = true
		}
	}
	if resizedCondition := conditions.Get(kcp, kcpv1.ResizedCondition); resizedCondition != nil {
		if resizedCondition.Status == corev1.ConditionTrue {
			record.KCPResized = true
		} else {
			switch resizedCondition.Reason {
			case kcpv1.ScalingUpReason:
				record.KCPScalingUp = true
			case kcpv1.ScalingDownReason:
				record.KCPScalingDown = true
			}
		}
	}
}

func recordInfraCluster(infraCluster *infrav1.AWSCluster, record *TestRecord) {
	if loadBalancerReady := conditions.Get(infraCluster, infrav1.LoadBalancerReadyCondition); loadBalancerReady != nil {
		if loadBalancerReady.Status == corev1.ConditionTrue {
			record.LoadBalancerReady = true
		} else {
			switch loadBalancerReady.Reason {
			case infrav1.WaitForDNSNameReason:
				record.WaitingForDNSName = true
			case infrav1.WaitForDNSNameResolveReason:
				record.WaitingForDNSResolution = true
			case infrav1.LoadBalancerFailedReason:
				record.LoadBalancerFailed = true
			}
		}
	}
}

func recordClusterResource(cluster *clusterv1.Cluster, record *TestRecord) {
	record.ClusterControlPlaneInitialized = cluster.Status.ControlPlaneInitialized
	if clusterInfraReady := conditions.Get(cluster, clusterv1.InfrastructureReadyCondition); clusterInfraReady != nil {
		if clusterInfraReady.Status == corev1.ConditionTrue {
			record.ClusterInfrastructureReady = true
		} else if clusterInfraReady.Reason == clusterv1.WaitingForInfrastructureFallbackReason {
			record.ClusterWaitingForInfrastructure = true
		}
	}
	if clusterControlPlaneReady := conditions.Get(cluster, clusterv1.ControlPlaneReadyCondition); clusterControlPlaneReady != nil {
		if clusterControlPlaneReady.Status == corev1.ConditionTrue {
			record.ClusterControlPlaneReady = true
		}
	}
	if clusterReady := conditions.Get(cluster, clusterv1.ReadyCondition); clusterReady != nil {
		if clusterReady.Status == corev1.ConditionTrue {
			record.ClusterReady = true
		}
	}
}

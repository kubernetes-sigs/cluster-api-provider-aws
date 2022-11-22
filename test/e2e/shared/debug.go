package shared

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	"sigs.k8s.io/cluster-api/test/framework"
)

type podDumper func(ctx context.Context, artifactDir string, clusterClient *kubernetes.Clientset, pods ...corev1.Pod)

var defaultPodDumpers = []podDumper{dumpPodState, dumpPodEvents}

func DumpWorkloadClusterResources(ctx context.Context, e2eCtx *E2EContext) {
	By("Getting all namespaces in bootstrap cluster")
	namespaces := corev1.NamespaceList{}
	clusterClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
	err := clusterClient.List(ctx, &namespaces)
	Expect(err).NotTo(HaveOccurred())

	for _, ns := range namespaces.Items {
		clusters := framework.GetAllClustersByNamespace(ctx, framework.GetAllClustersByNamespaceInput{
			Lister:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: ns.Name,
		})

		for _, cluster := range clusters {
			clusterClient := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, cluster.Namespace, cluster.Name).GetClientSet()
			dumpClusterWorkloads(ctx, e2eCtx, cluster.Name, clusterClient)
		}
	}
}

func dumpClusterWorkloads(ctx context.Context, e2eCtx *E2EContext, name string, clusterClient *kubernetes.Clientset) {
	Byf("Dumping workloads for cluster %s", name)

	logPath := filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", name)
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		fmt.Fprintf(GinkgoWriter, "couldn't create directory: path=%s, err=%s", logPath, err)
	}
	fmt.Fprintf(GinkgoWriter, "folder created for cluster: %s\n", logPath)

	//TODO: do we need the ability to filter these in the e2e config?
	namespacesToDump, err := clusterClient.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())

	for _, ns := range namespacesToDump.Items {
		dumpPodsForNamespace(ctx, &ns, logPath, clusterClient)
		//TODO: add any additional resource dumping here
	}
}

func dumpPodsForNamespace(ctx context.Context, ns *corev1.Namespace, artifactFolder string, clusterClient *kubernetes.Clientset) {
	logPath := filepath.Join(artifactFolder, ns.Name, "pods")
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		fmt.Fprintf(GinkgoWriter, "ERROR: couldn't create directory: path=%s, err=%s", logPath, err)
		return
	}
	fmt.Fprintf(GinkgoWriter, "folder created for cluster pods: %s\n", logPath)

	pods, err := clusterClient.CoreV1().Pods(ns.Name).List(ctx, v1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())

	if len(pods.Items) == 0 {
		fmt.Fprintf(GinkgoWriter, "no pods in namespace: %s\n", ns.Name)
		return
	}

	wg := sync.WaitGroup{}
	for _, dumpFunc := range defaultPodDumpers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dumpFunc(ctx, logPath, clusterClient, pods.Items...)
		}()
	}
	wg.Wait()
}

func dumpPodState(ctx context.Context, artifactFolder string, clusterClient *kubernetes.Clientset, pods ...corev1.Pod) {
	for _, pod := range pods {
		logPath := fmt.Sprintf("%s/%s_state.yaml", artifactFolder, pod.Name)

		data, err := yaml.Marshal(pod)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "ERROR: couldn't yaml marshal pod: pod=%s, ns=%s, err=%s", pod.Name, pod.Namespace, err)
			return
		}
		err = os.WriteFile(logPath, data, os.ModePerm)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "ERROR: couldn't save yaml: path=%s, err=%s", logPath, err)
			return
		}
	}
}

func dumpPodEvents(ctx context.Context, artifactFolder string, clusterClient *kubernetes.Clientset, pods ...corev1.Pod) {
	for _, pod := range pods {
		logPath := fmt.Sprintf("%s/%s_events.yaml", artifactFolder, pod.Name)

		events, err := clusterClient.CoreV1().Events(pod.Namespace).List(ctx, v1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%s", pod.Name),
		})
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "ERROR: list events for pod: pod=%s, err=%s", pod.Name, err)
			continue
		}
		for i := range events.Items {
			e := events.Items[i]
			e.ManagedFields = nil
			events.Items[i] = e
		}

		data, err := yaml.Marshal(events)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "ERROR: couldn't yaml marshal pod events: pod=%s, ns=%s, err=%s", pod.Name, pod.Namespace, err)
			continue
		}
		err = os.WriteFile(logPath, data, os.ModePerm)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "ERROR: couldn't save yaml: path=%s, err=%s", logPath, err)
			continue
		}
	}
}

func dumpPodLogs(ctx context.Context, artifactFolder string, clusterClient *kubernetes.Clientset, pods ...corev1.Pod) {
	for _, pod := range pods {
		containers := append(pod.Spec.Containers, pod.Spec.InitContainers...)
		for _, container := range containers {
			opts := &corev1.PodLogOptions{
				Container: container.Name,
				Previous:  false,
			}

			res, err := clusterClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts).Stream(ctx)
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "ERROR: couldn't get pod logs: pod=%s, err=%s", pod.Name, err)
				continue
			}
			defer res.Close()

			builder := &strings.Builder{}
			if _, err = io.Copy(builder, res); err != nil {
				fmt.Fprintf(GinkgoWriter, "ERROR: couldn't stream pod logs: pod=%s, err=%s", pod.Name, err)
				continue
			}

			logPath := fmt.Sprintf("%s/%s_%s.log", artifactFolder, pod.Name, &container.Name)
			err = os.WriteFile(logPath, []byte(builder.String()), os.ModePerm)
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "ERROR: couldn't save pod logs: path=%s, err=%s", logPath, err)
				continue
			}
		}
	}
}

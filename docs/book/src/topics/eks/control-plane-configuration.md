# Control Plane Component Configuration

EKS lets you set some parameters of the Kubernetes scheduler, API server and controller manager.
See [Advanced Kubernetes control plane configuration](https://docs.aws.amazon.com/eks/latest/userguide/control-plane-configuration.html)
for the supported parameters, their defaults and their supported values.

You set them with `kubeSchedulerConfig`, `kubeAPIServerConfig` and `kubeControllerManagerConfig` in the
`AWSManagedControlPlane`:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta2
metadata:
  name: "capi-managed-test-control-plane"
spec:
  ...
  kubeSchedulerConfig:
    nodeResourcesFit:
      scoringStrategy:
        type: MostAllocated
        resources:
          - name: cpu
            weight: 1
          - name: memory
            weight: 1
  kubeAPIServerConfig:
    eventTTL: 15m
    serviceNodePortRange:
      minPort: 30000
      maxPort: 32767
  kubeControllerManagerConfig:
    horizontalPodAutoscalerControllerConfig:
      horizontalPodAutoscalerSyncPeriod: 10s
    podGCControllerConfig:
      terminatedPodGCThreshold: 10000
```

CAPA sets these parameters when it creates the cluster, and updates them when they differ from the cluster.

- Only the parameters that you set are compared and updated. Parameters that you do not set keep their current value.
- Removing a parameter from the spec does not reset it on the cluster. To return a parameter to its default,
  set it to the default value.
- `horizontalPodAutoscalerSyncPeriod` and `terminatedPodGCThreshold` require a provisioned control plane tier,
  set with `controlPlaneScalingConfig`.

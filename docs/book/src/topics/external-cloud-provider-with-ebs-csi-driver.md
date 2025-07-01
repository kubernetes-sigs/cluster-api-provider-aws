# External AWS Cloud Provider and AWS CSI Driver

## Overview
The support for in-tree cloud providers and the CSI drivers is coming to an end and CAPA supports various upgrade paths
to use [external cloud provider (Cloud Controller Manager - CCM) ](https://github.com/kubernetes/cloud-provider-aws) and external CSI drivers.
This document explains how to create a CAPA cluster with external CSI/CCM plugins and how to upgrade existing clusters that rely on in-tree providers.


## Creating clusters with external CSI/CCM and validating
For clusters that will use external CCM, `cloud-provider: external` flag needs to be set in KubeadmConfig resources in both `KubeadmControlPlane` and `MachineDeployment` resources.

    clusterConfiguration:
      apiServer:
        extraArgs:
          cloud-provider: external
      controllerManager:
        extraArgs:
          cloud-provider: external
    initConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: external
    joinConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: external


External CCM and EBS CSI driver can be installed manually or using ClusterResourceSets (CRS) onto the CAPA workload cluster.
To install them with CRS, create a CRS resource on the management cluster with labels, for example `csi: external` and `ccm: external` labels.
Then, when creating `Cluster` objects for workload clusters that should have this CSR applied, create them with matching labels `csi: external` and `ccm: external` for CSI and CCM, respectively.

Manifests for installing the AWS CCM and the AWS EBS CSI driver are available from their respective
GitHub repositories (see [here for the AWS CCM](https://github.com/kubernetes/cloud-provider-aws) and
[here for the AWS EBS CSI driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver)).

An example of a workload cluster manifest with labels assigned for matching to a CRS can be found
[here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/tree/main/templates/cluster-template.yaml).

### Verifying dynamically provisioned volumes with CSI driver
Once you have the cluster with external CCM and CSI controller running successfully, you can test the CSI driver functioning with following steps after switching to workload cluster:
1. Create a service (say,`nginx`)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-svc
  namespace: default
spec:
  clusterIP: None
  ports:
    - name: nginx-web
      port: 80
  selector:
    app: nginx
```
2. Create a storageclass and statefulset for the service created above with the persistent volume assigned to the storageclass:
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: aws-ebs-volumes
provisioner: ebs.csi.aws.com
volumeBindingMode: WaitForFirstConsumer
parameters:
  csi.storage.k8s.io/fstype: xfs
  type: io1
  iopsPerGB: "100"
allowedTopologies:
  - matchLabelExpressions:
      - key: topology.ebs.csi.aws.com/zone
        values:
          - us-east-1a
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: nginx-statefulset
spec:
  serviceName: "nginx-svc"
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: registry.k8s.io/nginx-slim:0.8
          ports:
            - name: nginx-web
              containerPort: 80
          volumeMounts:
            - name: nginx-volumes
              mountPath: /usr/share/nginx/html
      volumes:
        - name: nginx-volumes
          persistentVolumeClaim:
            claimName: nginx-volumes
  volumeClaimTemplates:
    - metadata:
        name: nginx-volumes
      spec:
        storageClassName: "aws-ebs-volumes"
        accessModes: [ "ReadWriteOnce" ]
        resources:
          requests:
            storage: 4Gi
 ```
3. Once you apply the above manifest, the EBS volumes will be created and attached to the worker nodes.

>**IMPORTANT WARNING:** The CRDs from the AWS EBS CSI driver and AWS external cloud provider gives issue while installing the respective controllers on the AWS Cluster, it doesn't allow statefulsets to create the volume on existing EC2 instance.
> We need the CSI controller deployment and CCM pinned to the control plane which has right permissions to create, attach
> and mount the volumes to EC2 instances. To achieve this, you should add the node affinity rules to the CSI driver controller deployment and CCM DaemonSet manifests.
> ```yaml
> tolerations:
> - key: node-role.kubernetes.io/master
>   effect: NoSchedule
> - effect: NoSchedule
>   key: node-role.kubernetes.io/control-plane
> affinity:
>   nodeAffinity:
>   requiredDuringSchedulingIgnoredDuringExecution:
>     nodeSelectorTerms:
>       - matchExpressions:
>           - key: node-role.kubernetes.io/control-plane
>             operator: Exists
>       - matchExpressions:
>           - key: node-role.kubernetes.io/master
>             operator: Exists
>```


## Validated upgrade paths for existing clusters

From Kubernetes 1.23 onwards, `CSIMigrationAWS` flag is enabled by default, which requires the installation of [external CSI driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver), unless `CSIMigrationAWS` is disabled by the user.
For installing external CSI/CCM in the upgraded cluster, CRS can be used, see the section above for details.

CCM and CSI do not need to be migrated to use external plugins at the same time,
external CSI drivers works with in-tree CCM (Warning: using in-tree CSI with external CCM does not work).

**Following 3 upgrade paths are validated:**
- Scenario 1: During upgrade to v1.23.x, disabling `CSIMigrationAWS` flag and keep using in-tree CCM and CSI.
- Scenario 2: During upgrade to v1.23.x, enabling `CSIMigrationAWS` flag and using in-tree CCM with external CSI.
- Scenario 3: During upgrade to v1.23.x, enabling `CSIMigrationAWS` flag and using external CCM and CSI.


|                             | CSI      | CCM      | feature-gate CSIMigrationAWS | external-cloud-volume-plugin |
|-----------------------------|----------|----------|------------------------------|------------------------------|
| **Scenario 1**              |          |          |                              |                              |
| From Kubernetes  < v1.23    | in-tree  | in-tree  | off                          | NA                           |
| To Kubernetes    >= v1.23   | in-tree  | in-tree  | off                          | NA                           |
| **Scenario 2**              |          |          |                              |                              |
| From Kubernetes  < v1.23    | in-tree  | in-tree  | off                          | NA                           |
| To Kubernetes    >= v1.23   | external | in-tree  | on                           | NA                           |
| **Scenario 3**              |          |          |                              |                              |
| From Kubernetes  < v1.23    | in-tree  | in-tree  | off                          | NA                           |
| To Kubernetes    >= v1.23   | external | external | on                           | aws                          |


**KubeadmConfig in the upgraded cluster for scenario 1:**

    clusterConfiguration:
      apiServer:
        extraArgs:
          cloud-provider: aws
      controllerManager:
        extraArgs:
          cloud-provider: aws
          feature-gates: CSIMigrationAWS=false
    initConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: aws
          feature-gates: CSIMigrationAWS=false
        name: '{{ ds.meta_data.local_hostname }}'
    joinConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: aws
          feature-gates: CSIMigrationAWS=false

**KubeadmConfig in the upgraded cluster for scenario 2:**

When `CSIMigrationAWS=true`, installed external CSI driver will be used while relying on in-tree CCM.

    clusterConfiguration:
      apiServer:
        extraArgs:
          cloud-provider: aws
          feature-gates: CSIMigrationAWS=true   // Set only if Kubernetes version < 1.23.x, otherwise this flag is enabled by default.
      controllerManager:
        extraArgs:
          cloud-provider: aws
          feature-gates: CSIMigrationAWS=true   // Set only if Kubernetes version < 1.23.x, otherwise this flag is enabled by default.
    initConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: aws
          feature-gates: CSIMigrationAWS=true   // Set only if Kubernetes version < 1.23.x, otherwise this flag is enabled by default.
    joinConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: aws
          feature-gates: CSIMigrationAWS=true   // Set only if Kubernetes version < 1.23.x, otherwise this flag is enabled by default.

**KubeadmConfig in the upgraded cluster for scenario 3:**

`external-cloud-volume-plugin` flag needs to be set for old Kubelets to keep talking to in-tree CCM and upgrade fails without this is set.


    clusterConfiguration:
      apiServer:
        extraArgs:
          cloud-provider: external
      controllerManager:
        extraArgs:
          cloud-provider: external
          external-cloud-volume-plugin: aws
    initConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: external
    joinConfiguration:
      nodeRegistration:
        kubeletExtraArgs:
          cloud-provider: external

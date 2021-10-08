# External AWS cloud provider with AWS EBS CSI driver

## Overview
From Kubernetes 1.21 onwards, the support for its in-tree AWS cloud provider and the EBS CSI driver is removed, hence there is a need to use [out-of-tree cloud provider (Cloud Controller Manager - CCM) ](https://github.com/kubernetes/cloud-provider-aws) and a CSI driver in CAPA.
For details, see [Status of project and documentation of Cloud provider AWS](https://github.com/kubernetes/cloud-provider-aws/issues/42)

## Using external cloud provider and EBS CSI driver in AWS workloads
Once Management cluster is ready, install external CCM and EBS CSI driver onto the CAPA workload cluster either manually or using ClusterResourceSets (CRS).
To install them with CRS, create a CRS resource on the management cluster with labels, for example `csi: external` and `ccm: external` labels. 
Then, when creating `Cluster` objects for workload clusters that should have this CSR applied, create them with matching labels `csi: external` and `ccm: external` for CSI and CCM, respectively.

Example manifests for installing the AWS CCM and the AWS EBS CSI driver, and for creating corresponding CRS resources,
can be found [here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/tree/main/test/e2e/data/infrastructure-aws/kustomize_sources/external-cloud-provider).

An example of a workload cluster manifest with labels assigned for matching to a CRS can be found 
[here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/tree/main/templates/cluster-template-external-cloud-provider.yaml).

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
          image: k8s.gcr.io/nginx-slim:0.8
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

>**IMPORTANT WARNING:** The CRDs from the AWS EBS CSI driver and AWS out-of-tree cloud provider gives issue while installing the respective controllers on the AWS Cluster, it doesn't allow statefulsets to create the volume on existing EC2 instance.
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


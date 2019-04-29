# Creating clusters using cross account role assumption using KIAM

This document outlines the list of steps to create the target cluster via cross account role assumption using [KIAM](https://github.com/uswitch/kiam). 
KIAM lets the controller pod(s) to assume an AWS role that enables them create AWS resources necessary to create an
operational cluster. This way we wouldn't have to mount any AWS credentials or load environment variables to 
supply AWS credentials to the CAPA controller. This is automatically taken care by the KIAM components.
Note: If you dont want to use KIAM and rather want to mount the credentials as secrets, you may still achieve cross 
account role assumption by using multiple profiles. 

### Glossary

* Management cluster - The cluster that runs in AWS and is used to create target clusters in different AWS accounts
* Target account - The account where the target cluster is created
* Source account - The AWS account where the CAPA controllers for management cluster runs. 

## Goals
1. The CAPA controllers are running in an AWS account and you want to create the target cluster in another AWS account.
2. This assumes that you start with no existing clusters. 

## High level steps

1. Creating a management cluster in AWS - This can be done by running the phases in clusterctl
    * Uses the existing provider components yaml
2. Setting up cross account IAM roles
3. Deploying the KIAM server/agent
4. Create the target cluster (through KIAM)
    * Uses different provider components with no secrets and annotation to indicate the IAM Role to assume. 

## 1. Creating the management cluster in AWS
Using clusterctl command we can create a new cluster in AWS which in turn will act as the 
management cluster to create the target cluster(in a different AWS account. This can be achieved by using the phases in 
clusterctl to perform all the steps except the pivoting. This will provide us with a bare-bones functioning cluster that 
we can use as a management cluster. 
To begin with follow the steps in this [getting started guide](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/master/docs/getting-started.md) to setup the environment
except for creating the actual cluster. Instead follow the steps below to create the cluster.

create a new cluster using kind for bootstrapping purpose by running:

```$sh
kind create cluster --name <CLUSTER_NAME>
```

and get its kube config path by running

```$sh
export KIND_KUBECONFIG=`kind get kubeconfig-path`
```

Use the following commands to create new management cluster in AWS. 
```$sh
clusterctl alpha phases apply-cluster-api-components --provider-components cmd/clusterctl/examples/aws/out/provider-components.yaml \
--kubeconfig $KIND_KUBECONFIG

clusterctl alpha phases apply-cluster --cluster cmd/clusterctl/examples/aws/out/cluster.yaml --kubeconfig $KIND_KUBECONFIG
```

We only need to create the control plane on the cluster running in AWS source account. Since the example includes definition for a worker node, you may delete it. 

```$sh
clusterctl alpha phases apply-machines --machines cmd/clusterctl/examples/aws/out/machines.yaml --kubeconfig $KIND_KUBECONFIG

clusterctl alpha phases get-kubeconfig --provider aws --cluster-name <CLUSTER_NAME> --kubeconfig $KIND_KUBECONFIG

export AWS_KUBECONFIG=$PWD/kubeconfig

kubectl apply -f cmd/clusterctl/examples/aws/out/addons.yaml --kubeconfig $AWS_KUBECONFIG
```

Verify that all the pods in the kube-system namespace are running smoothly. Also you may remove the additional node in 
the machines example yaml since we are only interested in running the controllers that runs in control plane node
(although its not required to make any changes there). You can destroy your local kind cluster by running 

```$sh
kind delete cluster --name <CLUSTER_NAME>
```

## 2. Setting up cross account roles

In this step we will create new roles/policy in across 2 different AWS accounts.
Let us start by creating the roles in the account where the AWS controller runs. Following the directions 
posted in [KIAM repo](https://github.com/uswitch/kiam/blob/master/docs/IAM.md) create a "kiam_server" role
in AWS that only has a single managed policy with a single permission "sts:AssumeRole". Also add a trust policy on the 
 "kiam_server" role to include the role attached to the Control plane instance as a trusted entity. This looks something
 like this:
 
 In "kiam_server" role (Source AWS account):
 
 ```$json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::<AWS_ACCOUNT_NUMBER>:role/control-plane.cluster-api-provider-aws.sigs.k8s.io"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

Next we must establish a link between this "kiam_server" role on source AWS account and the role on target AWS account that has the permissions to create new cluster.
Begin by running the clusterawsadm cli to create a new stack on the target account where the target cluster is created. Make sure you use the credentials for target AWS account before creating the stack.

```clusterawsadm alpha bootstrap create-stack```

Then sign-in to the target AWS account to establish the link as mentioned above. Create a new Role with the permission policy set to "controllers.cluster-api-provider-aws.sigs.k8s.io". Lets name this role "cluster-api" for future reference. Add a new trust relationship to include the "kiam_server" role from the source account as trusted entity. This is shown below:

In "controllers.cluster-api-provider-aws.sigs.k8s.io" role(target AWS account)

```$json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::<SOURCE_AWS_ACCOUNT_NUMBER>:role/kserver"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

## 3. Deploying the KIAM server & agent
By now, your target cluster must be up & running. Make sure your KUBECONFIG pointing to the cluster in the target account. 

Create new secrets using the steps outlined [here](https://github.com/uswitch/kiam/blob/master/docs/TLS.md) 
Apply the manifest shown below: 
Make sure you update the argument to include your source AWS account "--assume-role-arn=arn:aws:iam::<SOURCE_AWS_ACCOUNT>:role/kiam_server"
server.yaml

```
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  namespace: kube-system
  name: kiam-server
spec:
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9620"
      labels:
        app: kiam
        role: server
    spec:
      serviceAccountName: kiam-server
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
      - operator: "Exists"
      volumes:
        - name: ssl-certs
          hostPath:
            # for AWS linux or RHEL distros
            # path: /etc/pki/ca-trust/extracted/pem/
            path: /etc/ssl/certs/
        - name: tls
          secret:
            secretName: kiam-server-tls
      hostNetwork: true
      containers:
        - name: kiam
          image: quay.io/uswitch/kiam:v3.2
          imagePullPolicy: Always
          command:
            - /kiam
          args:
            - server
            - --json-log
            - --level=warn
            - --bind=0.0.0.0:443
            - --cert=/etc/kiam/tls/server.pem
            - --key=/etc/kiam/tls/server-key.pem
            - --ca=/etc/kiam/tls/ca.pem
            - --role-base-arn-autodetect
            - --assume-role-arn=arn:aws:iam::<SOURCE_AWS_ACCOUNT>:role/kiam_server
            - --sync=1m
            - --prometheus-listen-addr=0.0.0.0:9620
            - --prometheus-sync-interval=5s
          volumeMounts:
            - mountPath: /etc/ssl/certs
              name: ssl-certs
            - mountPath: /etc/kiam/tls
              name: tls
          livenessProbe:
            exec:
              command:
              - /kiam
              - health
              - --cert=/etc/kiam/tls/server.pem
              - --key=/etc/kiam/tls/server-key.pem
              - --ca=/etc/kiam/tls/ca.pem
              - --server-address=127.0.0.1:443
              - --gateway-timeout-creation=1s
              - --timeout=5s
            initialDelaySeconds: 10
            periodSeconds: 10
            timeoutSeconds: 10
          readinessProbe:
            exec:
              command:
              - /kiam
              - health
              - --cert=/etc/kiam/tls/server.pem
              - --key=/etc/kiam/tls/server-key.pem
              - --ca=/etc/kiam/tls/ca.pem
              - --server-address=127.0.0.1:443
              - --gateway-timeout-creation=1s
              - --timeout=5s
            initialDelaySeconds: 3
            periodSeconds: 10
            timeoutSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: kiam-server
  namespace: kube-system
spec:
  clusterIP: None
  selector:
    app: kiam
    role: server
  ports:
  - name: grpclb
    port: 443
    targetPort: 443
    protocol: TCP
```

agent.yaml 

```
apiVersion: apps/v1
kind: DaemonSet
metadata:
  namespace: kube-system
  name: kiam-agent
spec:
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9620"
      labels:
        app: kiam
        role: agent
    spec:
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      tolerations:
        - operator: "Exists"
      volumes:
        - name: ssl-certs
          hostPath:
            # for AWS linux or RHEL distros
            #path: /etc/pki/ca-trust/extracted/pem/
            path: /etc/ssl/certs/
        - name: tls
          secret:
            secretName: kiam-agent-tls
        - name: xtables
          hostPath:
            path: /run/xtables.lock
            type: FileOrCreate
      containers:
        - name: kiam
          securityContext:
            capabilities:
              add: ["NET_ADMIN"]
          image: quay.io/uswitch/kiam:v3.2
          imagePullPolicy: Always
          command:
            - /kiam
          args:
            - agent
            - --iptables
            - --host-interface=cali+
            - --json-log
            - --port=8181
            - --cert=/etc/kiam/tls/agent.pem
            - --key=/etc/kiam/tls/agent-key.pem
            - --ca=/etc/kiam/tls/ca.pem
            - --server-address=kiam-server:443
            - --prometheus-listen-addr=0.0.0.0:9620
            - --prometheus-sync-interval=5s
            - --gateway-timeout-creation=1s
          env:
            - name: HOST_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - mountPath: /etc/ssl/certs
              name: ssl-certs
            - mountPath: /etc/kiam/tls
              name: tls
            - mountPath: /var/run/xtables.lock
              name: xtables
          livenessProbe:
            httpGet:
              path: /ping
              port: 8181
            initialDelaySeconds: 3
            periodSeconds: 3
```

server-rbac.yaml

```
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: kiam-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kiam-read
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  - pods
  verbs:
  - watch
  - get
  - list
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kiam-read
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kiam-read
subjects:
- kind: ServiceAccount
  name: kiam-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kiam-write
rules:
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kiam-write
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kiam-write
subjects:
- kind: ServiceAccount
  name: kiam-server
  namespace: kube-system
```

After deploying the above components make sure that the kiam_server and kiam_agent pods are up and running. 

## 4. Create the target cluster
Make sure you create copy of the "aws/out" directory called "out2". To create the target cluster we must update the provider_components.yaml generated in the out2 directory as shown below (to be run from the repository root directory):

```
cp cmd/clusterctl/examples/aws/out cmd/clusterctl/examples/aws/out2
vi cmd/clusterctl/examples/aws/out2/provider-components.yaml
```

1. Remove the credentials secret added at the bottom of the provider_components.yaml and do not mount the secret 
2. Add the following annotation to the template of aws-provider-controller-manager stateful set to specify the new role that was created in target account.
```
      annotations:
        iam.amazonaws.com/role: arn:aws:iam::<TARGET_AWS_ACCOUNT>:role/cluster-api
```
3. Also add this below annotation to the "aws-provider-system" namespace 
```
  annotations:
    iam.amazonaws.com/permitted: ".*"
```

Create a new cluster using the steps similar to the one used to create the source cluster. They are as follows:
```$sh
export SOURCE_KUBECONFIG=<PATH_TO_SOURCE_CLUSTER_KUBECONFIG>

clusterctl alpha phases apply-cluster-api-components --provider-components cmd/clusterctl/examples/aws/out2/provider-components.yaml \
--kubeconfig $SOURCE_KUBECONFIG

kubectl -f apply cmd/clusterctl/examples/aws/out2/cluster.yaml --kubeconfig $SOURCE_KUBECONFIG

kubectl apply -f cmd/clusterctl/examples/aws/out2/machines.yaml --kubeconfig $SOURCE_KUBECONFIG

clusterctl alpha phases get-kubeconfig --provider aws --cluster-name <TARGET_CLUSTER_NAME> --kubeconfig $SOURCE_KUBECONFIG
export KUBECONFIG=$PWD/kubeconfig

kubectl apply -f cmd/clusterctl/examples/aws/out2/addons.yaml

```
This creates the new cluster in the target AWS account.

**Creating clusters using cross account role assumption using kiam**

This document outlines the list of steps to create the target cluster via cross account role assumption using KIAM. 
KIAM lets the controller pod(s) to assume an AWS role that enables them create AWS resources necessary to create an
operational cluster. This way we wouldn't have to mount any AWS credentials or load environment variables to 
supply AWS credentials to the CAPA controller. This is automatically taken care by the KIAM components.
Note: If you dont want to use KIAM and rather want to mount the credentials as secrets, you may still achieve cross 
account role assumption by using multiple profiles. (TODO add this section at the bottom) 

**Glossory**

* Trusting Account - The account where the cluster is created
* Trusted Account - The AWS account where the CAPA controller runs. The "Trusting" account trusts the "Trusted" account
to create a new cluster in "Trusting" account. 

**Assumptions:**
1. The CAPA controllers are running in 1 AWS account and you want to create the target cluster in another AWS account.
(We could also use this doc to create a cluster in the same account)
2. This assumes that you start with no existing clusters. 

**High level steps**

1. Creating a bootstrap/management cluster in AWS - This can be done by running the phases in clusterctl
    * Uses the existing provider components yaml
2. Setting up cross account roles
3. Deploying the Kiam server/agent
4. Create the target cluster (through kiam)
    * Uses different provider components with no secrets and annotation to indicate the IAM Role to assume. 

**Creating the bootstrap cluster in AWS**
Using clusterctl command we can create a new cluster in AWS which in turn will act as the 
bootstrap cluster to create the target cluster(in a different AWS account. This can be achieved by using the phases in 
clusterctl to perform all the steps except the pivoting. This will provide us with a bare-bones functioning cluster that 
we can use as a bootstrap cluster. 
To begin with follow the steps in this getting started guide
(https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/master/docs/getting-started.md) to setup the environment
except for creating the actual cluster. Instead follow the steps below to create the cluster.

create a new cluster using kind for bootstrapping purpose by running:
```$xslt
kind create cluster
```
and get its kube config path by running
```
export KIND_KUBECONFIG=`kind get kubeconfig-path`
```

Use the following commands to create new (bootstrap) cluster in AWS
```
kubectl alpha phases apply-cluster-api-components --provider-components cmd/clusterctl/examples/aws/out/provider-components.yaml \
--kubeconfig $KIND_KUBECONFIG

kubectl alpha phases apply-cluster --cluster cmd/clusterctl/examples/aws/out/cluster.yaml --kubeconfig $KIND_KUBECONFIG

kubectl alpha phases apply-machines --machines cmd/clusterctl/examples/aws/out/machines.yaml 
--kubeconfig $KIND_KUBECONFIG

kubectl alpha phases get-kubeconfig --provider aws --cluster-name <CLUSTER_NAME> --kubeconfig $KIND_KUBECONFIG
export AWS_KUBECONFIG=`pwd`/kubeconfig

kubectl alpha phases apply-addons -a cmd/clusterctl/examples/aws/out/addons.yaml 
--kubeconfig $AWS_KUBECONFIG

```

Verify that all the pods in the kube-system namespace are running smoothly. Also you may remove the additional node in 
the machines example yaml since we are only interested in running the controllers that runs in control plane node
(although its not required to make any changes there). You can destroy your local kind cluster by running 
```$xslt
make kind-reset
```

**Setting up cross account roles:**

In this step we will new roles/policy in total across 2 different AWS accounts.
First lets start by creating the roles in the account where the AWS controller runs. Lets call this the "trusted" 
account since this account is trusted by the "trusting" account where the cluster is created. Following the directions 
posted here:https://github.com/uswitch/kiam/blob/master/docs/IAM.md create a "kiam_server" role
in AWS that only has a single managed policy with a single permission "sts:AssumeRole". Also add a trust policy on the 
 "kiam_server" role to include the role attached to the Control plane instance as a trusted entity. This looks something
 like this:
 ```$xslt
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

Next we must establish a link between this "kiam_server" and the role on trusting account that includes the permissions to 
create new cluster. This is done by using the similar steps as shown above. 
Sign in to the trusting account.
This requires running the clusterawsadm cli to create a new stack on the trusting account where the target cluster is 
created

```clusterawsadm alpha bootstrap create-stack```

The last role to be created is the one that is attached to the control plane that runs the CAPA controllers.
This role must have a minimal set of permissions 


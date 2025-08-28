# Using EKS Pod Identity for CAPA Controller

You can use [EKS Pod Identity](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html) to supply the credentials for the CAPA controller when the management is in EKS. This is an alternative to using the static boostrap credentials or IRSA.

## Pre-requisites

- Management cluster must be an EKS cluster
- AWS environment variables set for your account

## Steps

1. Install the **Amazon EKS Pod Identity Agent** EKS addon into the cluster. This can be done using the AWS console or using the AWS cli. 

> NOTE: If your management cluster is a "self-managed" CAPI cluster then its possible to install the addon via the **EKSManagedControlPlane**.

2. Create an EKS pod identity association for CAPA by running the following (replacing **<clustername>** with the name of your EKS cluster):

```bash
clusterawsadm controller use-pod-identity --cluster-name <clustername>
```

3. Ensure any credentials set for the controller are removed (a.k.a zeroed out):

```bash
clusterawsadm controller zero-credentials --namespace=capa-system
```

4. Force CAPA to restart so that the AWS credentials are injected:

```bash
clusterawsadm controller rollout-controller --kubeconfig=kubeconfig --namespace=capa-system
```
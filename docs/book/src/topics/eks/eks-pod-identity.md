# Using EKS Pod Identity for CAPA Controller

You can use [EKS Pod Identity](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html) to supply the credentials for the CAPA controller when the management cluster is in EKS. This is an alternative to using static bootstrap credentials or IRSA.

## Pre-requisites

- Management cluster must be an EKS cluster
- If using aws-cli, AWS [environment variables set](https://docs.aws.amazon.com/cli/v1/userguide/cli-configure-envvars.html) to target your account

## Steps

1. Install the **Amazon EKS Pod Identity Agent** EKS addon into the cluster. This can be done using the AWS console or using the AWS cli.

> NOTE: If your management cluster is managed by CAPI, it's possible to install the addon via the **EKSManagedControlPlane**.

2. Create an EKS pod identity association for CAPA by running the following (replacing `mycluster` with the name of your EKS cluster):

```bash
clusterawsadm controller use-pod-identity --cluster-name mycluster
```

3. Ensure any credentials set for the controller are removed (zeroed out):

```bash
clusterawsadm controller zero-credentials --namespace=capa-system
```

4. Force CAPA to restart so that the AWS credentials are injected:

```bash
clusterawsadm controller rollout-controller --kubeconfig=kubeconfig --namespace=capa-system
```

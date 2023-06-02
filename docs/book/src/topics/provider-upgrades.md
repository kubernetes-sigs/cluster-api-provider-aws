# Upgrading the Cluster API AWS Provider

The `clusterctl` tool ships with an upgrade path designed to check for available updates and apply a contract to upgrade a CAPA management cluster to the latest version. The upgrade path will download manifests for the latest versions of Cluster API and the AWS provider, apply those to your management cluster and scale down the `Deployments` running older versions.

During the upgrade, clusters provisioned by CAPA will continue to function normally. Still, the ability to use the management cluster to provision new clusters/manage or reconcile existing clusters will be restricted. For production workloads, it's best to upgrade your CAPA management cluster when you know you won't need to create or modify existing clusters.

## Prerequisites

1. If you've set any feature gates for Cluster API and/or the AWS provider, ensure you set them again using environment variables. If you omit this step, your management cluster will run into errors while trying to reconcile with existing clusters that depend on any feature gates. You would have to manually edit the `capi-controller-manager` and `capa-controller-manager` `Deployments` to open them again. Note that this also holds true for any other arguments you may be passing to `capa-controller-manager` - such as `awscluster-concurrency`, `sync-period` and so on.

2. If you're managing clusters in other accounts, the upgrade is also liable to reset the static management credentials used by CAPA to assume roles in those accounts - it's advisable to keep those credentials handy should you need to feed those into the cluster again.

To that end, some sample environment variables that you should set before running the upgrade are listed here:

```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=<your-static-access-key>
export AWS_SECRET_ACCESS_KEY=<your-static-secret-access-key>
export AWS_CONTROLLER_IAM_ROLE=<your-admin-iam-role>
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile) # This is used to help encode your environment variables

export EXP_EKS=true
export EXP_EKS_IAM=true
export EXP_MACHINE_POOL=true
export CAPA_EKS_ADD_ROLES=true
```

## Upgrade CAPI & CAPA

Assuming `$KUBECONFIG` is pointed at the kubeconfig for your management cluster:

1. Run `clusterctl upgrade plan`. This checks your management cluster's components against the latest releases to see if a new version is available.

If a new version is available - run `clusterctl upgrade apply --contract v1beta1`. This will power down existing Cluster API and AWS Provider controllers and deploy new ones.

Once the upgrade is complete, CAPA will pick up where it left off, managing and reconciling existing clusters.

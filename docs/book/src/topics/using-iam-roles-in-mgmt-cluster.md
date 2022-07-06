# Using IAM roles in management cluster instead of AWS credentials

## Overview

Sometimes users might want to use IAM roles to deploy management clusters. If the user already has a management cluster which was created using the AWS credentials, CAPA provides a way to use IAM roles instead of using these credentials.

## Pre-requisites
User has a bootstrap cluster created with AWS credentials. These credentials can be temporary as well.
To create temporary credentials, please follow [this doc](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_request.html).

We can verify whether this bootstrap cluster is using AWS credentials by checking the `capa-manager-bootstrap-credentials` secret created in `capa-system` namespace:
```bash
kubectl get secret -n capa-system capa-manager-bootstrap-credentials -o=jsonpath='{.data.credentials}' | { base64 -d 2>/dev/null || base64 -D; }
```
which will give output similar to below:
```bash
[default]
aws_access_key_id = <your-access-key>
aws_secret_access_key = <your-secret-access-key>
region = us-east-1

aws_session_token = <session-token>
```

## Goal
Create a management cluster which uses instance profiles (IAM roles) attached to EC2 instance.

## Steps for CAPA-managed clusters
1. Create a workload cluster on existing bootstrap cluster. Refer [quick start guide](https://cluster-api.sigs.k8s.io/user/quick-start.html) for more details.
   Since only control-plane nodes have the required IAM roles attached, CAPA deployment should have the necessary tolerations for master (control-plane) node and node selector for master.
> **Note:** A cluster with a single control plane node won’t be sufficient here due to the `NoSchedule` taint.

2. Get the kubeconfig for the new target management cluster(created in previous step) once it is up and running.

3. Zero the credentials CAPA controller started with, such that target management cluster uses empty credentials and not the previous credentials used to create bootstrap cluster using:
   ```bash
   clusterawsadm controller zero-credentials --namespace=capa-system
   ```
   For more details, please refer [zero-credentials doc](https://cluster-api-aws.sigs.k8s.io/clusterawsadm/clusterawsadm_controller_zero-credentials.html).

4. Rollout and restart on capa-controller-manager deployment using:
   ```bash
   clusterawsadm controller rollout-controller --kubeconfig=kubeconfig --namespace=capa-system
   ```
   For more details, please refer [rollout-controller doc](https://cluster-api-aws.sigs.k8s.io/clusterawsadm/clusterawsadm_controller_rollout-controller.html).

5. Use `clusterctl init` with the new cluster’s kubeconfig to install the provider components. For more details on preparing for init, please refer [clusterctl init doc](https://cluster-api.sigs.k8s.io/clusterctl/commands/init.html).

6. Use `clusterctl move` to move the Cluster API resources from the bootstrap cluster to the target management cluster. For more details on preparing for move, please refer [clusterctl move doc](https://cluster-api.sigs.k8s.io/clusterctl/commands/move.html).

7. Once the resources are moved to target management cluster successfully, `capa-manager-bootstrap-credentials` will be created as nil, and hence CAPA controllers will fall back to use the attached instance profiles.

8. Delete the bootstrap cluster with the AWS credentials.

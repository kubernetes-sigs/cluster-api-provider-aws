# Troubleshooting

## Cluster API with Docker - common issues with docker - 

When provisioning workload clusters using Cluster API with the Docker infrastructure provider,
provisioning might be stuck:

1. if there are stopped containers on your machine from previous runs. Clean unused containers with [docker rm -f ](https://docs.docker.com/engine/reference/commandline/rm/).

2. if the Docker space on your disk is being exhausted
    * Run [docker system df](https://docs.docker.com/engine/reference/commandline/system_df/) to inspect the disk space consumed by Docker resources.
    * Run [docker system prune --volumes](https://docs.docker.com/engine/reference/commandline/system_prune/) to prune dangling images, containers, volumes and networks.


## Cluster API with Docker  - "too many open files"
When creating many nodes using Cluster API and Docker infrastructure, either by creating large Clusters or a number of small Clusters, the OS may run into inotify limits which prevent new nodes from being provisioned.
If the error  `Failed to create inotify object: Too many open files` is present in the logs of the Docker Infrastructure provider this limit is being hit.

On Linux this issue can be resolved by increasing the inotify watch limits with:

```bash
sysctl fs.inotify.max_user_watches=1048576
sysctl fs.inotify.max_user_instances=8192
```

Newly created clusters should be able to take advantage of the increased limits.

### MacOS and Docker Desktop -  "too many open files"
This error was also observed in Docker Desktop 4.3 and 4.4 on MacOS. It can be resolved by updating to Docker Desktop for Mac 4.5 or using a version lower than 4.3.

[The upstream issue for this error is closed as of the release of Docker 4.5.0](https://github.com/docker/for-mac/issues/6071)

Note: The below workaround is not recommended unless upgrade or downgrade cannot be performed.

If using a version of Docker Desktop for Mac 4.3 or 4.4, the following workaround can be used:

Increase the maximum inotify file watch settings in the Docker Desktop VM:

1) Enter the Docker Desktop VM
```bash
nc -U ~/Library/Containers/com.docker.docker/Data/debug-shell.sock
```
2) Increase the inotify limits using sysctl
```bash
sysctl fs.inotify.max_user_watches=1048576
sysctl fs.inotify.max_user_instances=8192
```
3) Exit the Docker Desktop VM
```bash
exit
```

## Resources aren't being created

TODO

## Target cluster's control plane machine is up but target cluster's apiserver not working as expected

If `aws-provider-controller-manager-0` logs did not help, you might want to look into cloud-init logs, `/var/log/cloud-init-output.log`, on the controller host.
Verifying kubelet status and logs may also provide hints:
```bash
journalctl -u kubelet.service
systemctl status kubelet
```
For reaching controller host from your local machine:
```bash
 ssh -i <private-key> -o "ProxyCommand ssh -W %h:%p -i <private-key> ubuntu@<bastion-IP>" ubuntu@<controller-host-IP>
 ```

`private-key` is the private key from the key-pair discussed in the `ssh key pair` section above.

## kubelet on the control plane host failing with error: NoCredentialProviders
```bash
failed to run Kubelet: could not init cloud provider "aws": error finding instance i-0c276f2a1f1c617b2: "error listing AWS instances: \"NoCredentialProviders: no valid providers in chain. Deprecated.\\n\\tFor verbose messaging see aws.Config.CredentialsChainVerboseErrors\""
```
This error can occur if `CloudFormation` stack is not created properly and IAM instance profile is missing appropriate roles. Run following command to inspect IAM instance profile:
```bash
$ aws iam get-instance-profile --instance-profile-name control-plane.cluster-api-provider-aws.sigs.k8s.io --output json
{
    "InstanceProfile": {
        "InstanceProfileId": "AIPAJQABLZS4A3QDU576Q",
        "Roles": [
            {
                "AssumeRolePolicyDocument": {
                    "Version": "2012-10-17",
                    "Statement": [
                        {
                            "Action": "sts:AssumeRole",
                            "Effect": "Allow",
                            "Principal": {
                                "Service": "ec2.amazonaws.com"
                            }
                        }
                    ]
                },
                "RoleId": "AROAJQABLZS4A3QDU576Q",
                "CreateDate": "2019-05-13T16:45:12Z",
                "RoleName": "control-plane.cluster-api-provider-aws.sigs.k8s.io",
                "Path": "/",
                "Arn": "arn:aws:iam::123456789012:role/control-plane.cluster-api-provider-aws.sigs.k8s.io"
            }
        ],
        "CreateDate": "2019-05-13T16:45:28Z",
        "InstanceProfileName": "control-plane.cluster-api-provider-aws.sigs.k8s.io",
        "Path": "/",
        "Arn": "arn:aws:iam::123456789012:instance-profile/control-plane.cluster-api-provider-aws.sigs.k8s.io"
    }
}

```
If instance profile does not look as expected, you may try recreating the CloudFormation stack using `clusterawsadm` as explained in the above sections.


## Recover a management cluster after losing the api server load balancer

These steps outline the process for recovering a management cluster after losing the load balancer for the api server. These steps are needed because AWS load balancers have dynamically generated DNS names. This means that when a load balancer is deleted CAPA will recreate the load balancer but it will have a different DNS name that does not match the original, so we need to update some resources as well as the certs to match the new name to make the cluster healthy again. There are a few different scenarios which this could happen.

* The load balancer gets deleted by some external process or user.
* If a cluster is created with the same name as the management cluster in a different namespace and then deleted it will delete the existing load balancer. This is due to ownership of AWS resources being managed by tags. See this [issue](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/969#issuecomment-519121056) for reference.

### **Access the api server locally**

1. ssh to a control plane node and modify the `/etc/kubernetes/admin.conf`

    * Replace the `server` with `server: https://localhost:6443`

    * Add `insecure-skip-tls-verify: true`

    * Comment out `certificate-authority-data:`

2. Export the kubeconfig and ensure you can connect

    ```bash
    export KUBECONFIG=/etc/kubernetes/admin.conf
    kubectl get nodes
    ```


### **Get rid of the lingering duplicate cluster**

**This step is only needed in the scenario that duplicate cluster was created and deleted which caused the API server load balancer to be deleted.**

1. since there is a duplicate cluster that is trying to be deleted and can't due to some resources being unable to cleanup since they are in use we need to stop the conflicting reconciliation process. Edit the duplicate aws cluster object and remove the `finalizers`

    ```bash
    kubectl edit awscluster <clustername>
    ```
2. next run `kubectl describe awscluster <clustername>` to validate that the finalizers have been removed

3. `kubectl get clusters` to verify the cluster is gone


### **Make at least one node `Ready`**

1. Right now all endpoints are down due to nodes not being ready. this is problematic for coredns adn cni pods in particular. let's get one control plane node back healthy. on the control plane node we logged into edit the `/etc/kubernetes/kubelet.conf`

    * Replace the `server` with `server: https://localhost:6443`

    * Add `insecure-skip-tls-verify: true`

    * Comment out `certificate-authority-data:`

    * Restart the kubelet `systemctl restart kubelet`

2. `kubectl get nodes` and validate that the node is in a  ready state.
3. After a few minutes most things should start scheduling themselves on the new node. The pods that did not restart on their own that were causing issues were core-dns,kube-proxy, and cni pods.Those should be restart manually.
4. (optional) tail the capa logs to see the load balancer start to reconcile

    ```bash
    kubectl logs -f -n capa-system deployments.apps/capa-controller-manager`
    ```

### **Update the control plane nodes with new LB settings**

1. To be safe we will do this on all CP nodes rather than having them recreate to avoid potential data loss issues. Follow the following steps for **each** CP node.

2. Regenrate the certs for the api server using the new name. Make sure to update your service cidr and endpoint in the below command.

    ```bash
    rm /etc/kubernetes/pki/apiserver.crt
    rm /etc/kubernetes/pki/apiserver.key

    kubeadm init phase certs apiserver --control-plane-endpoint="mynewendpoint.com" --service-cidr=100.64.0.0/13 -v10
    ```

3. Update settings in `/etc/kubernetes/admin.conf`

    * Replace the `server` with `server: https://<your-new-lb.com>:6443`
    
    * Remove `insecure-skip-tls-verify: true`

    * Uncomment `certificate-authority-data:`

    * Export the kubeconfig and ensure you can connect 

        ```bash
        export KUBECONFIG=/etc/kubernetes/admin.conf
        kubectl get nodes
        ```

4. Update the settings in `/etc/kubernetes/kubelet.conf`

    * Replace the `server` with `server: https://your-new-lb.com:6443`

    * Remove `insecure-skip-tls-verify: true`

    * Uncomment `certificate-authority-data:`

    * restart the kubelet `systemctl restart kubelet`

5. Just as we did before we need new pods to pick up api server cache changes so  you will want to force restart pods like cni pods, kube-proxy, core-dns , etc.

### Update capi settings for new LB DNS name

1. Update the control plane endpoint on the `awscluster` and `cluster` objects. To do this we need to disable the validatingwebhooks. We will back them up and then delete so we can apply later.

    ```bash
    kubectl get validatingwebhookconfigurations capa-validating-webhook-configuration -o yaml > capa-webhook && kubectl delete validatingwebhookconfigurations capa-validating-webhook-configuration

    kubectl get validatingwebhookconfigurations capi-validating-webhook-configuration -o yaml > capi-webhook && kubectl delete validatingwebhookconfigurations capi-validating-webhook-configuration
    ```

2. Edit the `spec.controlPlaneEndpoint.host` field on both `awscluster` and `cluster` to have the new endpoint

3. Re-apply your webhooks

    ```bash
    kubectl apply -f capi-webhook
    kubectl apply -f capa-webhook
    ```


4. Update the following config maps and replace the old control plane name with the new one.

    ```bash
    kubectl edit cm -n kube-system kubeadm-config
    kubectl edit cm -n kube-system kube-proxy
    kubectl edit cm -n kube-public cluster-info
    ```

5. Edit the cluster kubeconfig secret that capi uses to talk to the management cluster. You will need to decode teh secret, replace the endpoint and re-encode and save.

    ```bash
    kubectl edit secret -n <namespace> <cluster-name>-kubeconfig`
    ```
6. At this point things should start to reconcile on their own, but we can use the commands in the next step to force it. 


### Roll all of the nodes to make sure everything is fresh


1. 
   ```bash
   kubectl patch kcp <clusternamekcp> -n namespace --type merge -p "{\"spec\":{\"rolloutAfter\":\"`date +'%Y-%m-%dT%TZ'`\"}}"
   ```
   
2. ```bash
    kubectl patch machinedeployment CLUSTER_NAME-md-0 -n namespace --type merge -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"`date +'%s'`\"}}}}}"
    ```

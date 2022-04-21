# Troubleshooting

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

### **Access the api server locally**

1. ssh to a control plane node and modify the `/etc/kubernetes/admin.conf`

    * Replace the `server` with `server: https://localhost:6443`

    * Add `insecure-skip-tls-verify: true`

    * Comment out `certificate-authority-data:`

2. Export the kubeconfig and ensure you can connect

    ```bash
    export KUBECONFIG=/etc/kubernetes/admin.conf`
    kubectl get nodes
    ```


### **Get rid of the lingering duplicate cluster**

1. since there is a duplicate cluster that is trying to be deleted and can't due to some resources being unable to cleanup since they are in use we need to stop the conflicting reconciliation process. Edit the duplicate aws cluster object and remove the `finalizers`

    ```bash
    kubectl edit awscluster <clustername>
    ```
2. `kubectl get clusters` to verify it's gone


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

1. `kubectl patch kcp <clusternamekcp> -n namespace --type merge -p "{\"spec\":{\"rolloutAfter\":\"`date +'%Y-%m-%dT%TZ'`\"}}"`
   
2. `kubectl patch machinedeployment CLUSTER_NAME-md-0 -n namespace --type merge -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"`date +'%s'`\"}}}}}"`
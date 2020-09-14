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

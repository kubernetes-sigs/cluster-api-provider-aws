# IAM Permissions

## Required to use clusterawsadm to provision IAM roles via CloudFormation

If using `clusterawsadm` to automate deployment of IAM roles via CloudFormation,
you must have IAM administrative access as `clusterawsadm` will provision IAM
roles and policies.

## Required by Cluster API Provider AWS controllers

The Cluster API Provider AWS controller requires permissions to use EC2, ELB
Autoscaling and optionally EKS. If provisioning IAM roles using `clusterawsadm`,
these will be set up as the `controllers.cluster-api-provider-aws.sigs.k8s.io`
IAM Policy, and attached to the `controllers.cluster-api-provider-aws.sigs.k8s.io`
and `control-plane.cluster-api-provider-aws.sigs.k8s.io` IAM roles.

### EC2 Provisioned Kubernetes Clusters

``` json
{
  "Type": "AWS::IAM::ManagedPolicy",
  "Properties": {
    "Description": "For the Kubernetes Cluster API Provider AWS Controllers",
    "ManagedPolicyName": "controllers.cluster-api-provider-aws.sigs.k8s.io",
    "PolicyDocument": {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "ec2:DescribeIpamPools",
            "ec2:AllocateIpamPoolCidr",
            "ec2:AttachNetworkInterface",
            "ec2:DetachNetworkInterface",
            "ec2:AllocateAddress",
            "ec2:AssignIpv6Addresses",
            "ec2:AssignPrivateIpAddresses",
            "ec2:UnassignPrivateIpAddresses",
            "ec2:AssociateRouteTable",
            "ec2:AssociateVpcCidrBlock",
            "ec2:AttachInternetGateway",
            "ec2:AuthorizeSecurityGroupIngress",
            "ec2:CreateCarrierGateway",
            "ec2:CreateInternetGateway",
            "ec2:CreateEgressOnlyInternetGateway",
            "ec2:CreateNatGateway",
            "ec2:CreateNetworkInterface",
            "ec2:CreateRoute",
            "ec2:CreateRouteTable",
            "ec2:CreateSecurityGroup",
            "ec2:CreateSubnet",
            "ec2:CreateTags",
            "ec2:CreateVpc",
            "ec2:CreateVpcEndpoint",
            "ec2:DisassociateVpcCidrBlock",
            "ec2:ModifyVpcAttribute",
            "ec2:ModifyVpcEndpoint",
            "ec2:DeleteCarrierGateway",
            "ec2:DeleteInternetGateway",
            "ec2:DeleteEgressOnlyInternetGateway",
            "ec2:DeleteNatGateway",
            "ec2:DeleteRouteTable",
            "ec2:ReplaceRoute",
            "ec2:DeleteSecurityGroup",
            "ec2:DeleteSubnet",
            "ec2:DeleteTags",
            "ec2:DeleteVpc",
            "ec2:DeleteVpcEndpoints",
            "ec2:DescribeAccountAttributes",
            "ec2:DescribeAddresses",
            "ec2:DescribeAvailabilityZones",
            "ec2:DescribeCarrierGateways",
            "ec2:DescribeInstances",
            "ec2:DescribeInstanceTypes",
            "ec2:DescribeInternetGateways",
            "ec2:DescribeEgressOnlyInternetGateways",
            "ec2:DescribeInstanceTypes",
            "ec2:DescribeImages",
            "ec2:DescribeNatGateways",
            "ec2:DescribeNetworkInterfaces",
            "ec2:DescribeNetworkInterfaceAttribute",
            "ec2:DescribeRouteTables",
            "ec2:DescribeSecurityGroups",
            "ec2:DescribeSubnets",
            "ec2:DescribeVpcs",
            "ec2:DescribeDhcpOptions",
            "ec2:DescribeVpcAttribute",
            "ec2:DescribeVpcEndpoints",
            "ec2:DescribeVolumes",
            "ec2:DescribeTags",
            "ec2:DetachInternetGateway",
            "ec2:DisassociateRouteTable",
            "ec2:DisassociateAddress",
            "ec2:ModifyInstanceAttribute",
            "ec2:ModifyNetworkInterfaceAttribute",
            "ec2:ModifySubnetAttribute",
            "ec2:ReleaseAddress",
            "ec2:RevokeSecurityGroupEgress",
            "ec2:RevokeSecurityGroupIngress",
            "ec2:RunInstances",
            "ec2:TerminateInstances",
            "ec2:GetSecurityGroupsForVpc",
            "tag:GetResources",
            "elasticloadbalancing:AddTags",
            "elasticloadbalancing:CreateLoadBalancer",
            "elasticloadbalancing:ConfigureHealthCheck",
            "elasticloadbalancing:DeleteLoadBalancer",
            "elasticloadbalancing:DeleteTargetGroup",
            "elasticloadbalancing:DescribeLoadBalancers",
            "elasticloadbalancing:DescribeLoadBalancerAttributes",
            "elasticloadbalancing:DescribeTargetGroups",
            "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
            "elasticloadbalancing:SetSecurityGroups",
            "elasticloadbalancing:DescribeTags",
            "elasticloadbalancing:ModifyLoadBalancerAttributes",
            "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
            "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
            "elasticloadbalancing:RemoveTags",
            "elasticloadbalancing:SetSubnets",
            "elasticloadbalancing:ModifyTargetGroupAttributes",
            "elasticloadbalancing:CreateTargetGroup",
            "elasticloadbalancing:DescribeListeners",
            "elasticloadbalancing:CreateListener",
            "elasticloadbalancing:DescribeTargetHealth",
            "elasticloadbalancing:RegisterTargets",
            "elasticloadbalancing:DeregisterTargets",
            "elasticloadbalancing:DeleteListener",
            "autoscaling:DescribeAutoScalingGroups",
            "autoscaling:DescribeInstanceRefreshes",
            "autoscaling:DeleteLifecycleHook",
            "autoscaling:DescribeLifecycleHooks",
            "autoscaling:PutLifecycleHook",
            "ec2:CreateLaunchTemplate",
            "ec2:CreateLaunchTemplateVersion",
            "ec2:DescribeLaunchTemplates",
            "ec2:DescribeLaunchTemplateVersions",
            "ec2:DeleteLaunchTemplate",
            "ec2:DeleteLaunchTemplateVersions",
            "ec2:DescribeKeyPairs",
            "ec2:ModifyInstanceMetadataOptions",
            "eks:CreateAccessEntry",
            "eks:DeleteAccessEntry",
            "eks:DescribeAccessEntry",
            "eks:UpdateAccessEntry",
            "eks:ListAccessEntries",
            "eks:AssociateAccessPolicy",
            "eks:DisassociateAccessPolicy",
            "eks:ListAssociatedAccessPolicies"
          ],
          "Resource": [
            "*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "autoscaling:CancelInstanceRefresh",
            "autoscaling:CreateAutoScalingGroup",
            "autoscaling:UpdateAutoScalingGroup",
            "autoscaling:CreateOrUpdateTags",
            "autoscaling:StartInstanceRefresh",
            "autoscaling:DeleteAutoScalingGroup",
            "autoscaling:DeleteTags"
          ],
          "Resource": [
            "arn:*:autoscaling:*:*:autoScalingGroup:*:autoScalingGroupName/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "autoscaling.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/spot.amazonaws.com/AWSServiceRoleForEC2Spot"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "spot.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:PassRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/*.cluster-api-provider-aws.sigs.k8s.io"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "secretsmanager:CreateSecret",
            "secretsmanager:DeleteSecret",
            "secretsmanager:TagResource"
          ],
          "Resource": [
            "arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*"
          ]
        }
      ]
    },
    "Roles": [
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sbGVycyIgfQ==",
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sUGxhbmUiIH0="
    ]
  }
}
```

### With EKS Support

``` json
{
  "Type": "AWS::IAM::ManagedPolicy",
  "Properties": {
    "Description": "For the Kubernetes Cluster API Provider AWS Controllers",
    "ManagedPolicyName": "controllers.cluster-api-provider-aws.sigs.k8s.io",
    "PolicyDocument": {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "ec2:DescribeIpamPools",
            "ec2:AllocateIpamPoolCidr",
            "ec2:AttachNetworkInterface",
            "ec2:DetachNetworkInterface",
            "ec2:AllocateAddress",
            "ec2:AssignIpv6Addresses",
            "ec2:AssignPrivateIpAddresses",
            "ec2:UnassignPrivateIpAddresses",
            "ec2:AssociateRouteTable",
            "ec2:AssociateVpcCidrBlock",
            "ec2:AttachInternetGateway",
            "ec2:AuthorizeSecurityGroupIngress",
            "ec2:CreateCarrierGateway",
            "ec2:CreateInternetGateway",
            "ec2:CreateEgressOnlyInternetGateway",
            "ec2:CreateNatGateway",
            "ec2:CreateNetworkInterface",
            "ec2:CreateRoute",
            "ec2:CreateRouteTable",
            "ec2:CreateSecurityGroup",
            "ec2:CreateSubnet",
            "ec2:CreateTags",
            "ec2:CreateVpc",
            "ec2:CreateVpcEndpoint",
            "ec2:DisassociateVpcCidrBlock",
            "ec2:ModifyVpcAttribute",
            "ec2:ModifyVpcEndpoint",
            "ec2:DeleteCarrierGateway",
            "ec2:DeleteInternetGateway",
            "ec2:DeleteEgressOnlyInternetGateway",
            "ec2:DeleteNatGateway",
            "ec2:DeleteRouteTable",
            "ec2:ReplaceRoute",
            "ec2:DeleteSecurityGroup",
            "ec2:DeleteSubnet",
            "ec2:DeleteTags",
            "ec2:DeleteVpc",
            "ec2:DeleteVpcEndpoints",
            "ec2:DescribeAccountAttributes",
            "ec2:DescribeAddresses",
            "ec2:DescribeAvailabilityZones",
            "ec2:DescribeCarrierGateways",
            "ec2:DescribeInstances",
            "ec2:DescribeInstanceTypes",
            "ec2:DescribeInternetGateways",
            "ec2:DescribeEgressOnlyInternetGateways",
            "ec2:DescribeInstanceTypes",
            "ec2:DescribeImages",
            "ec2:DescribeNatGateways",
            "ec2:DescribeNetworkInterfaces",
            "ec2:DescribeNetworkInterfaceAttribute",
            "ec2:DescribeRouteTables",
            "ec2:DescribeSecurityGroups",
            "ec2:DescribeSubnets",
            "ec2:DescribeVpcs",
            "ec2:DescribeDhcpOptions",
            "ec2:DescribeVpcAttribute",
            "ec2:DescribeVpcEndpoints",
            "ec2:DescribeVolumes",
            "ec2:DescribeTags",
            "ec2:DetachInternetGateway",
            "ec2:DisassociateRouteTable",
            "ec2:DisassociateAddress",
            "ec2:ModifyInstanceAttribute",
            "ec2:ModifyNetworkInterfaceAttribute",
            "ec2:ModifySubnetAttribute",
            "ec2:ReleaseAddress",
            "ec2:RevokeSecurityGroupEgress",
            "ec2:RevokeSecurityGroupIngress",
            "ec2:RunInstances",
            "ec2:TerminateInstances",
            "ec2:GetSecurityGroupsForVpc",
            "tag:GetResources",
            "elasticloadbalancing:AddTags",
            "elasticloadbalancing:CreateLoadBalancer",
            "elasticloadbalancing:ConfigureHealthCheck",
            "elasticloadbalancing:DeleteLoadBalancer",
            "elasticloadbalancing:DeleteTargetGroup",
            "elasticloadbalancing:DescribeLoadBalancers",
            "elasticloadbalancing:DescribeLoadBalancerAttributes",
            "elasticloadbalancing:DescribeTargetGroups",
            "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
            "elasticloadbalancing:SetSecurityGroups",
            "elasticloadbalancing:DescribeTags",
            "elasticloadbalancing:ModifyLoadBalancerAttributes",
            "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
            "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
            "elasticloadbalancing:RemoveTags",
            "elasticloadbalancing:SetSubnets",
            "elasticloadbalancing:ModifyTargetGroupAttributes",
            "elasticloadbalancing:CreateTargetGroup",
            "elasticloadbalancing:DescribeListeners",
            "elasticloadbalancing:CreateListener",
            "elasticloadbalancing:DescribeTargetHealth",
            "elasticloadbalancing:RegisterTargets",
            "elasticloadbalancing:DeregisterTargets",
            "elasticloadbalancing:DeleteListener",
            "autoscaling:DescribeAutoScalingGroups",
            "autoscaling:DescribeInstanceRefreshes",
            "autoscaling:DeleteLifecycleHook",
            "autoscaling:DescribeLifecycleHooks",
            "autoscaling:PutLifecycleHook",
            "ec2:CreateLaunchTemplate",
            "ec2:CreateLaunchTemplateVersion",
            "ec2:DescribeLaunchTemplates",
            "ec2:DescribeLaunchTemplateVersions",
            "ec2:DeleteLaunchTemplate",
            "ec2:DeleteLaunchTemplateVersions",
            "ec2:DescribeKeyPairs",
            "ec2:ModifyInstanceMetadataOptions",
            "eks:CreateAccessEntry",
            "eks:DeleteAccessEntry",
            "eks:DescribeAccessEntry",
            "eks:UpdateAccessEntry",
            "eks:ListAccessEntries",
            "eks:AssociateAccessPolicy",
            "eks:DisassociateAccessPolicy",
            "eks:ListAssociatedAccessPolicies"
          ],
          "Resource": [
            "*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "autoscaling:CancelInstanceRefresh",
            "autoscaling:CreateAutoScalingGroup",
            "autoscaling:UpdateAutoScalingGroup",
            "autoscaling:CreateOrUpdateTags",
            "autoscaling:StartInstanceRefresh",
            "autoscaling:DeleteAutoScalingGroup",
            "autoscaling:DeleteTags"
          ],
          "Resource": [
            "arn:*:autoscaling:*:*:autoScalingGroup:*:autoScalingGroupName/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "autoscaling.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/spot.amazonaws.com/AWSServiceRoleForEC2Spot"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "spot.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:PassRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/*.cluster-api-provider-aws.sigs.k8s.io"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "secretsmanager:CreateSecret",
            "secretsmanager:DeleteSecret",
            "secretsmanager:TagResource"
          ],
          "Resource": [
            "arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*"
          ]
        }
      ]
    },
    "Roles": [
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sbGVycyIgfQ==",
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sUGxhbmUiIH0="
    ]
  }
}
```

### With S3 Support
``` json
{
  "Type": "AWS::IAM::ManagedPolicy",
  "Properties": {
    "Description": "For the Kubernetes Cluster API Provider AWS Controllers",
    "ManagedPolicyName": "controllers.cluster-api-provider-aws.sigs.k8s.io",
    "PolicyDocument": {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "ec2:DescribeIpamPools",
            "ec2:AllocateIpamPoolCidr",
            "ec2:AttachNetworkInterface",
            "ec2:DetachNetworkInterface",
            "ec2:AllocateAddress",
            "ec2:AssignIpv6Addresses",
            "ec2:AssignPrivateIpAddresses",
            "ec2:UnassignPrivateIpAddresses",
            "ec2:AssociateRouteTable",
            "ec2:AssociateVpcCidrBlock",
            "ec2:AttachInternetGateway",
            "ec2:AuthorizeSecurityGroupIngress",
            "ec2:CreateCarrierGateway",
            "ec2:CreateInternetGateway",
            "ec2:CreateEgressOnlyInternetGateway",
            "ec2:CreateNatGateway",
            "ec2:CreateNetworkInterface",
            "ec2:CreateRoute",
            "ec2:CreateRouteTable",
            "ec2:CreateSecurityGroup",
            "ec2:CreateSubnet",
            "ec2:CreateTags",
            "ec2:CreateVpc",
            "ec2:CreateVpcEndpoint",
            "ec2:DisassociateVpcCidrBlock",
            "ec2:ModifyVpcAttribute",
            "ec2:ModifyVpcEndpoint",
            "ec2:DeleteCarrierGateway",
            "ec2:DeleteInternetGateway",
            "ec2:DeleteEgressOnlyInternetGateway",
            "ec2:DeleteNatGateway",
            "ec2:DeleteRouteTable",
            "ec2:ReplaceRoute",
            "ec2:DeleteSecurityGroup",
            "ec2:DeleteSubnet",
            "ec2:DeleteTags",
            "ec2:DeleteVpc",
            "ec2:DeleteVpcEndpoints",
            "ec2:DescribeAccountAttributes",
            "ec2:DescribeAddresses",
            "ec2:DescribeAvailabilityZones",
            "ec2:DescribeCarrierGateways",
            "ec2:DescribeInstances",
            "ec2:DescribeInstanceTypes",
            "ec2:DescribeInternetGateways",
            "ec2:DescribeEgressOnlyInternetGateways",
            "ec2:DescribeInstanceTypes",
            "ec2:DescribeImages",
            "ec2:DescribeNatGateways",
            "ec2:DescribeNetworkInterfaces",
            "ec2:DescribeNetworkInterfaceAttribute",
            "ec2:DescribeRouteTables",
            "ec2:DescribeSecurityGroups",
            "ec2:DescribeSubnets",
            "ec2:DescribeVpcs",
            "ec2:DescribeDhcpOptions",
            "ec2:DescribeVpcAttribute",
            "ec2:DescribeVpcEndpoints",
            "ec2:DescribeVolumes",
            "ec2:DescribeTags",
            "ec2:DetachInternetGateway",
            "ec2:DisassociateRouteTable",
            "ec2:DisassociateAddress",
            "ec2:ModifyInstanceAttribute",
            "ec2:ModifyNetworkInterfaceAttribute",
            "ec2:ModifySubnetAttribute",
            "ec2:ReleaseAddress",
            "ec2:RevokeSecurityGroupEgress",
            "ec2:RevokeSecurityGroupIngress",
            "ec2:RunInstances",
            "ec2:TerminateInstances",
            "ec2:GetSecurityGroupsForVpc",
            "tag:GetResources",
            "elasticloadbalancing:AddTags",
            "elasticloadbalancing:CreateLoadBalancer",
            "elasticloadbalancing:ConfigureHealthCheck",
            "elasticloadbalancing:DeleteLoadBalancer",
            "elasticloadbalancing:DeleteTargetGroup",
            "elasticloadbalancing:DescribeLoadBalancers",
            "elasticloadbalancing:DescribeLoadBalancerAttributes",
            "elasticloadbalancing:DescribeTargetGroups",
            "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
            "elasticloadbalancing:SetSecurityGroups",
            "elasticloadbalancing:DescribeTags",
            "elasticloadbalancing:ModifyLoadBalancerAttributes",
            "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
            "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
            "elasticloadbalancing:RemoveTags",
            "elasticloadbalancing:SetSubnets",
            "elasticloadbalancing:ModifyTargetGroupAttributes",
            "elasticloadbalancing:CreateTargetGroup",
            "elasticloadbalancing:DescribeListeners",
            "elasticloadbalancing:CreateListener",
            "elasticloadbalancing:DescribeTargetHealth",
            "elasticloadbalancing:RegisterTargets",
            "elasticloadbalancing:DeregisterTargets",
            "elasticloadbalancing:DeleteListener",
            "autoscaling:DescribeAutoScalingGroups",
            "autoscaling:DescribeInstanceRefreshes",
            "autoscaling:DeleteLifecycleHook",
            "autoscaling:DescribeLifecycleHooks",
            "autoscaling:PutLifecycleHook",
            "ec2:CreateLaunchTemplate",
            "ec2:CreateLaunchTemplateVersion",
            "ec2:DescribeLaunchTemplates",
            "ec2:DescribeLaunchTemplateVersions",
            "ec2:DeleteLaunchTemplate",
            "ec2:DeleteLaunchTemplateVersions",
            "ec2:DescribeKeyPairs",
            "ec2:ModifyInstanceMetadataOptions",
            "eks:CreateAccessEntry",
            "eks:DeleteAccessEntry",
            "eks:DescribeAccessEntry",
            "eks:UpdateAccessEntry",
            "eks:ListAccessEntries",
            "eks:AssociateAccessPolicy",
            "eks:DisassociateAccessPolicy",
            "eks:ListAssociatedAccessPolicies"
          ],
          "Resource": [
            "*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "autoscaling:CancelInstanceRefresh",
            "autoscaling:CreateAutoScalingGroup",
            "autoscaling:UpdateAutoScalingGroup",
            "autoscaling:CreateOrUpdateTags",
            "autoscaling:StartInstanceRefresh",
            "autoscaling:DeleteAutoScalingGroup",
            "autoscaling:DeleteTags"
          ],
          "Resource": [
            "arn:*:autoscaling:*:*:autoScalingGroup:*:autoScalingGroupName/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "autoscaling.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:CreateServiceLinkedRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/aws-service-role/spot.amazonaws.com/AWSServiceRoleForEC2Spot"
          ],
          "Condition": {
            "StringLike": {
              "iam:AWSServiceName": "spot.amazonaws.com"
            }
          }
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:PassRole"
          ],
          "Resource": [
            "arn:*:iam::*:role/*.cluster-api-provider-aws.sigs.k8s.io"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "secretsmanager:CreateSecret",
            "secretsmanager:DeleteSecret",
            "secretsmanager:TagResource"
          ],
          "Resource": [
            "arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "s3:CreateBucket",
            "s3:DeleteBucket",
            "s3:DeleteObject",
            "s3:GetObject",
            "s3:ListBucket",
            "s3:PutBucketPolicy",
            "s3:PutBucketTagging",
            "s3:PutLifecycleConfiguration",
            "s3:PutObject"
          ],
          "Resource": [
            "arn:*:s3:::cluster-api-provider-aws-*"
          ]
        }
      ]
    },
    "Roles": [
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sbGVycyIgfQ==",
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sUGxhbmUiIH0="
    ]
  }
}
```

## Required by the Kubernetes AWS Cloud Provider

These permissions are used by the Kubernetes AWS Cloud Provider. If you are
running with the in-tree cloud provider, this will typically be used by the
`controller-manager` pod in the `kube-system` namespace.

If provisioning IAM roles using `clusterawsadm`,
these will be set up as the `control-plane.cluster-api-provider-aws.sigs.k8s.io`
IAM Policy, and attached to the `control-plane.cluster-api-provider-aws.sigs.k8s.io`
IAM role.

``` json
{
  "Type": "AWS::IAM::ManagedPolicy",
  "Properties": {
    "Description": "For the Kubernetes Cloud Provider AWS Control Plane",
    "ManagedPolicyName": "control-plane.cluster-api-provider-aws.sigs.k8s.io",
    "PolicyDocument": {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "autoscaling:DescribeAutoScalingGroups",
            "autoscaling:DescribeLaunchConfigurations",
            "autoscaling:DescribeTags",
            "ec2:AssignIpv6Addresses",
            "ec2:DescribeInstances",
            "ec2:DescribeImages",
            "ec2:DescribeRegions",
            "ec2:DescribeRouteTables",
            "ec2:DescribeSecurityGroups",
            "ec2:DescribeSubnets",
            "ec2:DescribeVolumes",
            "ec2:CreateSecurityGroup",
            "ec2:CreateTags",
            "ec2:CreateVolume",
            "ec2:ModifyInstanceAttribute",
            "ec2:ModifyVolume",
            "ec2:AttachVolume",
            "ec2:AuthorizeSecurityGroupIngress",
            "ec2:CreateRoute",
            "ec2:DeleteRoute",
            "ec2:DeleteSecurityGroup",
            "ec2:DeleteVolume",
            "ec2:DetachVolume",
            "ec2:RevokeSecurityGroupIngress",
            "ec2:DescribeVpcs",
            "elasticloadbalancing:AddTags",
            "elasticloadbalancing:AttachLoadBalancerToSubnets",
            "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
            "elasticloadbalancing:SetSecurityGroups",
            "elasticloadbalancing:CreateLoadBalancer",
            "elasticloadbalancing:CreateLoadBalancerPolicy",
            "elasticloadbalancing:CreateLoadBalancerListeners",
            "elasticloadbalancing:ConfigureHealthCheck",
            "elasticloadbalancing:DeleteLoadBalancer",
            "elasticloadbalancing:DeleteLoadBalancerListeners",
            "elasticloadbalancing:DescribeLoadBalancers",
            "elasticloadbalancing:DescribeLoadBalancerAttributes",
            "elasticloadbalancing:DetachLoadBalancerFromSubnets",
            "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
            "elasticloadbalancing:ModifyLoadBalancerAttributes",
            "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
            "elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
            "elasticloadbalancing:CreateListener",
            "elasticloadbalancing:CreateTargetGroup",
            "elasticloadbalancing:DeleteListener",
            "elasticloadbalancing:DeleteTargetGroup",
            "elasticloadbalancing:DeregisterTargets",
            "elasticloadbalancing:DescribeListeners",
            "elasticloadbalancing:DescribeLoadBalancerPolicies",
            "elasticloadbalancing:DescribeTargetGroups",
            "elasticloadbalancing:DescribeTargetHealth",
            "elasticloadbalancing:ModifyListener",
            "elasticloadbalancing:ModifyTargetGroup",
            "elasticloadbalancing:RegisterTargets",
            "elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
            "iam:CreateServiceLinkedRole",
            "kms:DescribeKey"
          ],
          "Resource": [
            "*"
          ]
        }
      ]
    },
    "Roles": [
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sUGxhbmUiIH0="
    ]
  }
}
```
## Required by all nodes

All nodes require these permissions in order to run, and are used by the AWS
cloud provider run by kubelet.

If provisioning IAM roles using `clusterawsadm`,
these will be set up as the `nodes.cluster-api-provider-aws.sigs.k8s.io`
IAM Policy, and attached to the `nodes.cluster-api-provider-aws.sigs.k8s.io`
IAM role.


``` json
{
  "Type": "AWS::IAM::ManagedPolicy",
  "Properties": {
    "Description": "For the Kubernetes Cloud Provider AWS nodes",
    "ManagedPolicyName": "nodes.cluster-api-provider-aws.sigs.k8s.io",
    "PolicyDocument": {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "ec2:AssignIpv6Addresses",
            "ec2:DescribeInstances",
            "ec2:DescribeRegions",
            "ec2:CreateTags",
            "ec2:DescribeTags",
            "ec2:DescribeNetworkInterfaces",
            "ec2:DescribeInstanceTypes",
            "ecr:GetAuthorizationToken",
            "ecr:BatchCheckLayerAvailability",
            "ecr:GetDownloadUrlForLayer",
            "ecr:GetRepositoryPolicy",
            "ecr:DescribeRepositories",
            "ecr:ListImages",
            "ecr:BatchGetImage"
          ],
          "Resource": [
            "*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "secretsmanager:DeleteSecret",
            "secretsmanager:GetSecretValue"
          ],
          "Resource": [
            "arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "ssm:UpdateInstanceInformation",
            "ssmmessages:CreateControlChannel",
            "ssmmessages:CreateDataChannel",
            "ssmmessages:OpenControlChannel",
            "ssmmessages:OpenDataChannel",
            "s3:GetEncryptionConfiguration"
          ],
          "Resource": [
            "*"
          ]
        }
      ]
    },
    "Roles": [
      "eyAiUmVmIjogIkFXU0lBTVJvbGVDb250cm9sUGxhbmUiIH0=",
      "eyAiUmVmIjogIkFXU0lBTVJvbGVOb2RlcyIgfQ=="
    ]
  }
}
```

When using EKS, the `AmazonEKSWorkerNodePolicy` and `AmazonEKS_CNI_Policy`
AWS managed policies will also be attached to
`nodes.cluster-api-provider-aws.sigs.k8s.io` IAM role.

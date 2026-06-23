# Control Plane Egress Routing

## Overview

The default behaviour of EKS is to manage the routing of controlplane traffic to resources in your VPC, for example, when webhooks are called. This is called "AWS Managed".

There may be situations where you want full control over how the traffic is routed. This is especially useful for highly regulated environments. This is called "customer routed".

The table below summarizes the 2 modes:

| Mode | Description | 
| --- | --- | 
|  `aws-managed`  | Default behavior. Amazon EKS manages the egress path from the control plane ENIs. You don’t need to configure NAT gateways or other routing infrastructure for control plane traffic. | 
|  `customer-routed`  | You manage the egress path from the control plane in your VPC subnets. You are responsible for ensuring that the control plane can reach required endpoints (such as webhook servers, OIDC providers, and other resources). You provide an egress path, such as a NAT gateway, NAT instance, transit gateway, or firewall appliance. You also configure the route table, network ACL, and security group rules that allow this traffic. | 

**Important**  
In `customer-routed` mode, you’re responsible for ensuring proper network connectivity from the control plane. Misconfigurations in your VPC networking can cause control plane operations to fail. These misconfigurations include a missing egress path, restrictive network ACLs, or incorrect security groups. Affected operations include admission webhook calls and OIDC authentication.

## How to use

You can set the control plane egress mode using the `controlPlaneEgressMode` field of the `AWSManagedControlPlane`:

```yaml
kind: AWSManagedControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "capi-managed-test-control-plane"
spec:
  ...
  controlPlaneEgressMode: aws-managed
```

## Further information

See the AWS documentation [here](https://docs.aws.amazon.com/eks/latest/userguide/control-plane-egress.html).


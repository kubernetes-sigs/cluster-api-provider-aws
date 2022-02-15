# CAPA Version Support

## Supported Versions

The Cluster API Provider AWS team maintains branches for **v1.x (v1beta1)**, **v0.7 (v1alpha4)**, and **v0.6 (v1alpha3)**.
These branches end of support follows their compatible Cluster API branch support.

Releases include these components:

- Cluster API Provider for AWS
- clusterawsadm client

## Compatibility with Cluster API and Kubernetes Versions

This provider's versions are compatible with the following versions of Cluster API
and support all Kubernetes versions that is supported by its compatible Cluster API version:


|                              | v1alpha3 (v0.3) | v1alpha4 (v0.4) | v1beta1 (v1.x) |
| ---------------------------- | --------------- | --------------- | -------------- |
| AWS Provider v1alpha3 (v0.6) | ✓               |                 |                |
| AWS Provider v1alpha4 (v0.7) |                 | ✓               |                |
| AWS Provider v1beta1 (v1.x)  |                 |                 | ✓              |


(See [Kubernetes support matrix][cluster-api-supported-v] of Cluster API versions).

[cluster-api-supported-v]: https://cluster-api.sigs.k8s.io/reference/versions.html

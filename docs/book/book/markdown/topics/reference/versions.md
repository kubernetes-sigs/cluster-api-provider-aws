# CAPA Version Support

## Release Versioning

CAPA follows the [semantic versionining][semver] specification:

MAJOR version release for incompatible API changes, 
MINOR version release for backwards compatible feature additions, 
and PATCH version release for only bug fixes.

**Example versions:**

- Minor release: `v0.1.0`
- Patch release: `v0.1.1`
- Major release: `v1.0.0`


## Compatibility with Cluster API Versions

CAPA's versions are compatible with the following versions of Cluster API

|     API Version              | Cluster API v1alpha3 (v0.3) | Cluster API v1alpha4 (v0.4) |
| ---------------------------- | --------------------------- | --------------------------- |
| AWS Provider v1alpha3 (v0.6) | ✓                           |                             |
| AWS Provider v1alpha4 (v0.7) |                             | ✓                           |


CAPA v1beta1 versions are not released in lock-step with Cluster API releases.
Multiple CAPA minor releases can use the same Cluster API minor release.

For compatibility, check the release notes [here](https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/) to see which v1beta1 Cluster API version each CAPA version is compatible with.

For example:
- CAPA v1.0.x, v1.1.x, v1.2.x is compatible with Cluster API v1.0.x
- CAPA v1.3.x is compatible with Cluster API v1.1.x

## End-of-Life Timeline

CAPA team maintains branches for **v1.x (v1beta1)**, **v0.7 (v1alpha4)**, and **v0.6 (v1alpha3)**.

CAPA branches follow their compatible Cluster API branch EOL date.

| API Version   | Branch      | Supported Until |
| ------------- |-------------|-----------------|
| **v1alpha4**  | release-0.7 | 2022-04-06      |
| **v1alpha3**  | release-0.6 | 2022-02-23      |

## Compatibility with Kubernetes Versions

 CAPA API versions support all Kubernetes versions that is supported by its compatible Cluster API version:

|     API Versions             | CAPI v1alpha3 (v0.3) | CAPI v1alpha4 (v0.4) | CAPI v1beta1 (v1.x) |
| ---------------------------- | --------------- | --------------- | -------------- |
| CAPA v1alpha3 (v0.6)         | ✓               |                 |                |
| CAPA v1alpha4 (v0.7)         |                 | ✓               |                |
| CAPA v1beta1 (v1.x)          |                 |                 | ✓              |


(See [Kubernetes support matrix][cluster-api-supported-v] of Cluster API versions).

[cluster-api-supported-v]: https://cluster-api.sigs.k8s.io/reference/versions.html
[semver]: https://semver.org/#semantic-versioning-200

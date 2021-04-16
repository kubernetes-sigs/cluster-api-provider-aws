# Multi-tenancy

Starting from v0.6.5, single controller multi-tenancy is supported that allows using a different AWS Identity for each workload cluster.
For details, see the [multi-tenancy proposal](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/docs/proposal/20200506-single-controller-multitenancy.md).


For multi-tenancy support, a reference field (`identityRef`) is added to `AWSCluster`, which describes the identity to be used when reconciling the cluster.

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: <IdentityType>
    name: <IdentityName>
```

Identity resources are used to describe IAM identities that will be used during reconciliation.
There are three identity types: AWSClusterControllerIdentity, AWSClusterStaticIdentity, and AWSClusterRoleIdentity.
Once an IAM identity is created in AWS, the corresponding values should be used to create a identity resource.

## AWSClusterControllerIdentity

Before multi-tenancy support, all AWSClusters were being reconciled using the credentials that are used by Cluster API Provider AWS Controllers.
`AWSClusterControllerIdentity` is used to restrict the usage of controller credentials only to AWSClusters that are in `allowedNamespaces`.
Since CAPA controllers use a single set of credentials, `AWSClusterControllerIdentity` is a singleton, and can only be created with `name: default`.

For backward compatibility, `AutoControllerIdentityCreator` experimental feature is added, which is responsible to create the `AWSClusterControllerIdentity` singleton if it does not exist.
- **Feature status:** Experimental
- **Feature gate:** AutoControllerIdentityCreator=true
`AutoControllerIdentityCreator` creates `AWSClusterControllerIdentity` singleton with empty `allowedNamespaces` (allowedNamespaces: {}) to grant access to the `AWSClusterControllerIdentity` from all namespaces.

Example:
```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSClusterControllerIdentity
    name: default
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterControllerIdentity
metadata:
  name: "default"
spec:
  allowedNamespaces:{}  # matches all namespaces
```
`AWSClusterControllerIdentity` is immutable to avoid any unwanted overrides to the allowed namespaces, especially during upgrading clusters.

## AWSClusterIdentityIdentity
`AWSClusterIdentityIdentity` represents static AWS credentials, which are stored in a `Secret`.

Example: Below, an `AWSClusterIdentityIdentity` is created that allows access to the `AWSClusters` that are in "test" namespace.
The identity credentials that will be used by "test" AWSCluster are stored in "test-account-creds" secret.


```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSClusterIdentityIdentity
    name: test-account
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterIdentityIdentity
metadata:
  name: "test-account"
spec:
  secretRef:
    name: test-account-creds
    namespace: capa-system
  allowedNamespaces:
    selector:
      matchLabels:
        ns: "testlabel"
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    cluster.x-k8s.io/ns: "testlabel"
  name: "test"
---
apiVersion: v1
kind: Secret
metadata:
  name: "test-account-creds"
  namespace: capa-system
stringData:
 accessKeyID: AKIAIOSFODNN7EXAMPLE
 secretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

## AWSClusterRoleIdentity
`AWSClusterRoleIdentity` allows CAPA to assume a role either in the same or another AWS account, using the STS::AssumeRole API.
The assumed role could be used by the AWSClusters that is in the `allowedNamespaces`.

Example:
Below, an `AWSClusterRoleIdentity` instance, which will be used by AWSCluster "test", is created.
This role will be assumed by the source identity at runtime. Source identity can be of any identity type.
Role is assumed in the beginning once and after, whenever the assumed role's credentials are expired.

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSCluster
metadata:
  name: "test"
  namespace: "test"
spec:
  region: "eu-west-1"
  identityRef:
    kind: AWSClusterRoleIdentity
    name: test-account-role
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterRoleIdentity
metadata:
  name: "test-account-role"
spec:
  allowedNamespaces:
    list: # allows only "test" namespace to use this identity
      "test"
  roleARN: "arn:aws:iam::123456789:role/CAPARole"
  sourceIdentityRef:
    kind: AWSClusterIdentityIdentity
    name: test-account-creds
```

Nested role assumption is also supported.
Example: Below, "multi-tenancy-nested-role" will be assumed by "multi-tenancy-role", which will be assumed by the "default" `AWSClusterControllerIdentity`

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterRoleIdentity
metadata:
  name: multi-tenancy-role
spec:
  allowedNamespaces:
    list: []
  durationSeconds: 900 # default and min value is 900 seconds
  roleARN: arn:aws:iam::11122233344:role/multi-tenancy-role
  sessionName: multi-tenancy-role-session
  sourceidentityRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: AWSClusterControllerIdentity
    name: default
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterRoleIdentity
metadata:
  name: multi-tenancy-nested-role
spec:
  allowedNamespaces:
    list: []
  roleARN: arn:aws:iam::11122233355:role/multi-tenancy-nested-role
  sessionName: multi-tenancy-nested-role-session
  sourceidentityRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: AWSClusterRoleIdentity
    name: multi-tenancy-role
```

## Secure Access to Identitys
`allowedNamespaces` field is used to grant access to the namespaces to use Identitys.
Only AWSClusters that are created in one of the Identity's allowed namespaces can use that Identity.
`allowedNamespaces` are defined by providing either a list of namespaces or label selector to select namespaces.

### Examples

An empty `allowedNamespaces` indicates that the Identity can be used by all namespaces.

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterControllerIdentity
spec:
  allowedNamespaces:{}  # matches all namespaces
```

Having a nil `list` and a nil `selector` is the same with having an empty `allowedNamespaces` (Identity can be used by all namespaces).

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterControllerIdentity
spec:
  allowedNamespaces:
    list: nil
    selector: nil
```

A nil `allowedNamespaces` indicates that the Identity cannot be used from any namespace.

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterControllerIdentity
spec:
  allowedNamespaces:  # this is same with not providing the field at all or allowedNamespaces: null
```

The union of namespaces that are matched by `selector` and the namespaces that are in the `list` is granted access to the identity.
The namespaces that are not in the list and not matching the selector will not have access.

Nil or empty `list` matches no namespaces. Nil or empty `selector` matches no namespaces.
If `list` is nil and `selector` is empty OR `list` is empty and `selector` is nil, Identity cannot be used from any namespace.
Because in this case, `allowedNamespaces` is not empty or nil, and neither `list` nor `selector` allows any namespaces, so the union is empty.

```yaml
# Matches no namespaces
allowedNamespaces:
  list: []
```
```yaml
# Matches no namespaces
allowedNamespaces:
  selector: {}
```
```yaml
# Matches no namespaces
allowedNamespaces:
  list: null
  selector: {}
```
```yaml
# Matches no namespaces
allowedNamespaces:
  list: []
  selector: {}
```

**Important** The default behaviour of an empty label selector is to match all objects, however here we do not follow that behavior to avoid unintended access to the identitys.
This is consistent with core cluster API selectors, e.g., Machine and ClusterResourceSet selectors. The result of matchLabels and matchExpressions are ANDed.


In Kubernetes selectors, `matchLabels` and `matchExpressions` are ANDed.
In the example below, list is empty/nil, so does not allow any namespaces and selector matches with only `default` namespace.
Since `list` and `selector` results are ORed, `default` namespace can use this identity.

```yaml
kind: namespace
metadata:
  name: default
  labels:
    environment: dev
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: AWSClusterControllerIdentity
spec:
  allowedNamespaces:
    list: null # or []
    selector:
      matchLabels:
        namespace: default
      matchExpressions:
        - {key: environment, operator: In, values: [dev]}
```

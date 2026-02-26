# Restricting Cluster API to certain namespaces <!-- omit in toc -->

Cluster-api-provider-aws controllers by default, reconcile cluster-api objects
across all namespaces in the cluster. However, it is possible to restrict
reconciliation to a single namespace and this document tells you how.

## Contents <!-- omit in toc -->

- [Use cases](#use-cases)
- [Configuring `cluster-api-provider-aws` controllers](#configuring-cluster-api-provider-aws-controllers)

## Use cases

- Grouping clusters into a namespace based on the AWS account will allow
  managing clusters across multiple AWS accounts. This will require each
  `cluster-api-provider-aws` controller to have credentials to their respective
  AWS accounts. These credentials can be created as kubernetes secret and be
  mounted in the pod at `/home/.aws` or as environment variables.
- Grouping clusters into a namespace based on their environment, (test,
  qualification, canary, production) will allow a phased rolling out of
  `cluster-api-provider-aws` releases.
- Grouping clusters into a namespace based on the infrastructure provider will
  allow running multiple cluster-api provider implementations side-by-side and
  manage clusters across infrastructure providers.

## Configuring `cluster-api-provider-aws` controllers

- Create the namespace that `cluster-api-provider-aws` controller will watch for
  cluster-api objects

```(bash)
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: my-pet-clusters #edit if necessary
EOF
```

- Deploy/edit `aws-provider-controller-manager` controller statefulset

Specifically, edit the container spec for `cluster-api-aws-controller`, in the
`aws-provider-controller-manager` statefulset, to pass a value to the `namespace`
CLI flag.

```(bash)
        - -namespace=my-pet-clusters # edit this if necessary
```

Once the `aws-provider-controller-manager-0` pod restarts,
`cluster-api-provider-aws` controllers will only reconcile the cluster-api
objects in the `my-pet-clusters` namespace.

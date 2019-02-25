# CVO Metrics

The Cluster Version Operator reports the following metrics:

The cluster version is reported as seconds since the epoch with labels for `version` and
`image`. The `type` label reports which value is reported:

* `current` - the version the operator is applying right now (the running CVO version) and the age of the payload
* `cluster` - the same as current, but the value is the creation timestamp of the cluster version (cluster age)
* `failure` - if the failure condition is set, reports the last transition time for both desired and current versions.
* `desired` - reported if different from current as the most recent timestamp on the cluster version
* `completed` - the time the most recent version was completely applied, or is zero.

```
# HELP cluster_version Reports the version of the cluster.
# TYPE cluster_version gauge
cluster_version{image="test/image:1",type="current",version="4.0.2"} 130000000
cluster_version{image="test/image:1",type="failure",version="4.0.2"} 132000400
cluster_version{image="test/image:2",type="desired",version="4.0.3"} 132000400
cluster_version{image="test/image:1",type="completed",version="4.0.2"} 132000100
cluster_version{image="test/image:1",type="cluster",version="4.0.2"} 131000000
# HELP cluster_version_available_updates Report the count of available versions for an upstream and channel.
# TYPE cluster_version_available_updates gauge
cluster_version_available_updates{channel="fast",upstream="https://api.openshift.com/api/upgrades_info/v1/graph"} 0
```

Metrics about cluster operators:

```
# HELP cluster_operator_conditions Report the conditions for active cluster operators. 0 is False and 1 is True.
# TYPE cluster_operator_conditions gauge
cluster_operator_conditions{condition="Available",name="version",namespace="openshift-cluster-version"} 1
cluster_operator_conditions{condition="Failing",name="version",namespace="openshift-cluster-version"} 0
cluster_operator_conditions{condition="Progressing",name="version",namespace="openshift-cluster-version"} 0
cluster_operator_conditions{condition="RetrievedUpdates",name="version",namespace="openshift-cluster-version"} 0
# HELP cluster_operator_up Reports key highlights of the active cluster operators.
# TYPE cluster_operator_up gauge
cluster_operator_up{name="version",namespace="openshift-cluster-version",version="4.0.1"} 1
```

Metrics reported while applying the image:

```
# HELP cluster_version_payload Report the number of entries in the image.
# TYPE cluster_version_payload gauge
cluster_version_payload{type="applied",version="4.0.3"} 0
cluster_version_payload{type="pending",version="4.0.3"} 1
# HELP cluster_operator_payload_errors Report the number of errors encountered applying the image.
# TYPE cluster_operator_payload_errors gauge
cluster_operator_payload_errors{version="4.0.3"} 10
```
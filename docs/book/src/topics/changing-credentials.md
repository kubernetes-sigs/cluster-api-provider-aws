# Changing Credentials

Although `clusterctl init` embeds a set of AWS credentials into the configuration of the Cluster API Provider for AWS (CAPA) when initializing the management cluster, it may be necessary after deployment to change or update those AWS credentials.

CAPA stores the credentials it uses as a Kubernetes Secret in the "capa-system" namespace. You can use `kubectl -n capa-system get secrets` and you'll see the "capa-manager-bootstrap-credentials" Secret. The credentials themselves are stored as a key named `credentials`; you can use this command to retrieve the credentials and decode them (if you're using macOS, change the `-d` to `-D`):

    kubectl -n capa-system get secret capa-manager-bootstrap-credentials \
    -o jsonpath="{.data.credentials}" | base64 -d

The command will return something like this (but with valid access key ID, secret access key, and region values, obviously):

    [default]
    aws_access_key_id = <access-key-id-value-here>
    aws_secret_access_key = <secret-access-key-value-here>
    region = <aws-region-here>

There are two ways to change the AWS credentials used by CAPA:

1. As part of upgrading Cluster API components using `clusterctl upgrade`
2. Manually

## When Upgrading Cluster API Components

When upgrading the Cluster API components using `clusterctl upgrade` in a management cluster that has CAPA installed, `clusterctl` requires that the variable "AWS_B64ENCODED_CREDENTIALS" is defined with a value. This value is used by `clusterctl` to update the Kubernetes Secret described above. Thus, if the management cluster needs to be upgraded, users can change/update the AWS credentials used by CAPA as part of the upgrade process.

The process would look like this:

1. Use `clusterawsadm` to define the "AWS_B64ENCODED_CREDENTIALS" variable. With `clusterawsadm` 0.5.4 and earlier, the command is `clusterawsadm alpha bootstrap encode-aws-credentials`; with version 0.5.5 or later, the command is `clusterawsadm bootstrap credentials encode-as-profile`. Assign the output of this command to the "AWS_B64ENCODED_CREDENTIALS" variable. Users should ensure that `clusterawsadm` picks up the _new_ or _updated_ credentials when running this command.
2. Use `clusterctl` to upgrade the management cluster, first by running `clusterctl upgrade plan` and then by running `clusterctl upgrade apply` (using the command provided in the output of `clusterctl upgrade plan`). As part of the upgrade process, the credentials stored in the "AWS_B64ENCODED_CREDENTIALS" variable will be placed into the Kubernetes Secret used by CAPA.

If the management cluster does not need to be upgraded or can't be upgraded at the current time, then users will need to manually change the credentials.

## Manually Changing Credentials

To manually change the credentials used by CAPA, follow these steps:

1. Use `clusterawsadm` to create the material needed for the Kubernetes Secret. With `clusterawsadm` 0.5.5 and earlier, the command is `clusterawsadm alpha bootstrap encode-aws-credentials`; with version 0.6.0 or later, the command is `clusterawsadm bootstrap credentials encode-as-profile`. Record the output of this command, as it is needed in a later step.
2. Use `kubectl -n capa-system edit secret capa-manager-bootstrap-credentials` to edit the Secret. Replace the existing value of the `data.credentials` field with the new value created in step 1 using `clusterawsadm`. Save your changes.
3. For the CAPA controller manager to pick up the new credentials in the Secret, restart it with `kubectl -n capa-system rollout restart deployment capa-controller-manager`.

Upon restart, CAPA will now use the updated credentials in the Secret to communicate with AWS.

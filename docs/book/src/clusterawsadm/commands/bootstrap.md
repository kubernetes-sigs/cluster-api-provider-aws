# clusterawsadm bootstrap

To use Cluster API Provider AWS, an AWS account needs to
be prepared with AWS Identity and Access Management (IAM) roles that will be used by
clusters as well as provide Cluster API Provider AWS with credentials to use to provision infrastructure.

The `clusterawsadm bootstrap` command provides access management related utilities. 

## clusterawsadm bootstrap iam
The `clusterawsadm bootstrap iam` command provides CRUD operations on IAM roles necessary for controllers.
The clusterawsadm utility retrieve the local credentials and uses them to create a CloudFormation stack in your AWS account with the correct IAM resources.

### Create/update an AWS CloudFormation stack that has the required IAM roles 
CloudFormation stack can be customized using a config file that has `AWSIAMConfiguration`,

**Create bootstrap file:**
This is the config file that can be used to customize the CloudFormation stack, this step can be bypassed if no customization is necessary.
Documentation for this kind can be found at [here](https://pkg.go.dev/sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/bootstrap/v1alpha1
).

```bash
$ cat config-bootstrap.yaml
apiVersion: bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1
kind: AWSIAMConfiguration
spec:
  bootstrapUser:
    enable: true
  eventBridge:
    enable: true
  eks:
    enable: true
```

**Create/Update CloudFormation stack:**

When running this command for the first time, it will attempt to create a new CloudFormation stack that have the necessary IAM roles.
If the CloudFormation stack already exists but there are some changes, it updates the stack.
```bash
$ clusterawsadm bootstrap iam create-cloudformation-stack --config=config-bootstrap.yaml
Attempting to create AWS CloudFormation stack cluster-api-provider-aws-sigs-k8s-io
```

**NOTE:** During update, the same region that is used during creation must be used with this command. Although IAM roles are global, CloudFormation stacks are not,
hence update request must be sent to the same region; otherwise update fails because it tries to create a new CloudFormation stack in the new region,
and since the same roles have already been created before and available globally.
```bash
$ clusterawsadm bootstrap iam create-cloudformation-stack     
Attempting to create AWS CloudFormation stack cluster-api-provider-aws-sigs-k8s-io
Error: failed to create AWS CloudFormation stack: ResourceNotReady: failed waiting for successful resource state
```

### View AWS IAM policies and CloudFormation stack

To view the CloudFormation Template:
```bash
$ clusterawsadm bootstrap iam print-cloudformation-template
```

To view the bootstrap config that is used to customize the CloudFormation stack:
```bash
$ clusterawsadm bootstrap iam print-config
```

To view the IAM policies such as `AWSIAMManagedPolicyControllers`:
```bash
$ clusterawsadm bootstrap iam print-policy --document <policy-name>
```

## clusterawsadm bootstrap credentials
This command is used to encode credentials to use with Kubernetes Cluster API Provider AWS.
It retrieves AWS account credentials from your environment and encode it with base64.

`clusterctl` requires that the variable "AWS_B64ENCODED_CREDENTIALS" is defined with a value.
This command is needed both during first installation of Cluster API provider AWS (`clusterctl init`) and during upgrade (`clusterctl upgrade`)
A Kubernetes Secret is created with this data during `clusterctl init`, which is used by CAPA controllers. See [upgrade](./../changing-credentials.md) for more information about how credentials that are used by CAPA can be upgraded using this Secret.
```bash
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)
```
<p>Packages:</p>
<ul>
<li>
<a href="#ami.aws.infrastructure.cluster.x-k8s.io%2fv1beta1">ami.aws.infrastructure.cluster.x-k8s.io/v1beta1</a>
</li>
<li>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io%2fv1alpha1">bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1</a>
</li>
<li>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io%2fv1beta1">bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1</a>
</li>
<li>
<a href="#bootstrap.cluster.x-k8s.io%2fv1beta1">bootstrap.cluster.x-k8s.io/v1beta1</a>
</li>
<li>
<a href="#bootstrap.cluster.x-k8s.io%2fv1beta2">bootstrap.cluster.x-k8s.io/v1beta2</a>
</li>
<li>
<a href="#controlplane.cluster.x-k8s.io%2fv1beta1">controlplane.cluster.x-k8s.io/v1beta1</a>
</li>
<li>
<a href="#controlplane.cluster.x-k8s.io%2fv1beta2">controlplane.cluster.x-k8s.io/v1beta2</a>
</li>
<li>
<a href="#infrastructure.cluster.x-k8s.io%2fv1beta1">infrastructure.cluster.x-k8s.io/v1beta1</a>
</li>
<li>
<a href="#infrastructure.cluster.x-k8s.io%2fv1beta2">infrastructure.cluster.x-k8s.io/v1beta2</a>
</li>
</ul>
<h2 id="ami.aws.infrastructure.cluster.x-k8s.io/v1beta1">ami.aws.infrastructure.cluster.x-k8s.io/v1beta1</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the AMI v1beta1 API group</p>
</p>
Resource Types:
<ul></ul>
<h3 id="ami.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSAMI">AWSAMI
</h3>
<p>
<p>AWSAMI defines an AMI.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#ami.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSAMISpec">
AWSAMISpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>os</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imageID</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>kubernetesVersion</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="ami.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSAMISpec">AWSAMISpec
</h3>
<p>
(<em>Appears on:</em><a href="#ami.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSAMI">AWSAMI</a>)
</p>
<p>
<p>AWSAMISpec defines an AMI.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>os</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imageID</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>kubernetesVersion</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1">bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains API Schema definitions for the bootstrap v1alpha1 API group</p>
</p>
Resource Types:
<ul></ul>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfiguration">AWSIAMConfiguration
</h3>
<p>
<p>AWSIAMConfiguration controls the creation of AWS Identity and Access Management (IAM) resources for use
by Kubernetes clusters and Kubernetes Cluster API Provider AWS.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">
AWSIAMConfigurationSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>namePrefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NamePrefix will be prepended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to &ldquo;&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>nameSuffix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NameSuffix will be appended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to
&ldquo;.cluster-api-provider-aws.sigs.k8s.io&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlane</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ControlPlane">
ControlPlane
</a>
</em>
</td>
<td>
<p>ControlPlane controls the configuration of the AWS IAM role for a Kubernetes cluster&rsquo;s control plane nodes.</p>
</td>
</tr>
<tr>
<td>
<code>clusterAPIControllers</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ClusterAPIControllers">
ClusterAPIControllers
</a>
</em>
</td>
<td>
<p>ClusterAPIControllers controls the configuration of an IAM role and policy specifically for Kubernetes Cluster API Provider AWS.</p>
</td>
</tr>
<tr>
<td>
<code>nodes</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.Nodes">
Nodes
</a>
</em>
</td>
<td>
<p>Nodes controls the configuration of the AWS IAM role for all nodes in a Kubernetes cluster.</p>
</td>
</tr>
<tr>
<td>
<code>bootstrapUser</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.BootstrapUser">
BootstrapUser
</a>
</em>
</td>
<td>
<p>BootstrapUser contains a list of elements that is specific
to the configuration and enablement of an IAM user.</p>
</td>
</tr>
<tr>
<td>
<code>stackName</code><br/>
<em>
string
</em>
</td>
<td>
<p>StackName defines the name of the AWS CloudFormation stack.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>Region controls which region the control-plane is created in if not specified on the command line or
via environment variables.</p>
</td>
</tr>
<tr>
<td>
<code>eks</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.EKSConfig">
EKSConfig
</a>
</em>
</td>
<td>
<p>EKS controls the configuration related to EKS. Settings in here affect the control plane
and nodes roles</p>
</td>
</tr>
<tr>
<td>
<code>eventBridge</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.EventBridgeConfig">
EventBridgeConfig
</a>
</em>
</td>
<td>
<p>EventBridge controls configuration for consuming EventBridge events</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>secureSecretBackends</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecretBackend">
[]SecretBackend
</a>
</em>
</td>
<td>
<p>SecureSecretsBackend, when set to parameter-store will create AWS Systems Manager
Parameter Storage policies. By default or with the value of secrets-manager,
will generate AWS Secrets Manager policies instead.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfiguration">AWSIAMConfiguration</a>)
</p>
<p>
<p>AWSIAMConfigurationSpec defines the specification of the AWSIAMConfiguration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>namePrefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NamePrefix will be prepended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to &ldquo;&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>nameSuffix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NameSuffix will be appended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to
&ldquo;.cluster-api-provider-aws.sigs.k8s.io&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlane</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ControlPlane">
ControlPlane
</a>
</em>
</td>
<td>
<p>ControlPlane controls the configuration of the AWS IAM role for a Kubernetes cluster&rsquo;s control plane nodes.</p>
</td>
</tr>
<tr>
<td>
<code>clusterAPIControllers</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ClusterAPIControllers">
ClusterAPIControllers
</a>
</em>
</td>
<td>
<p>ClusterAPIControllers controls the configuration of an IAM role and policy specifically for Kubernetes Cluster API Provider AWS.</p>
</td>
</tr>
<tr>
<td>
<code>nodes</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.Nodes">
Nodes
</a>
</em>
</td>
<td>
<p>Nodes controls the configuration of the AWS IAM role for all nodes in a Kubernetes cluster.</p>
</td>
</tr>
<tr>
<td>
<code>bootstrapUser</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.BootstrapUser">
BootstrapUser
</a>
</em>
</td>
<td>
<p>BootstrapUser contains a list of elements that is specific
to the configuration and enablement of an IAM user.</p>
</td>
</tr>
<tr>
<td>
<code>stackName</code><br/>
<em>
string
</em>
</td>
<td>
<p>StackName defines the name of the AWS CloudFormation stack.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>Region controls which region the control-plane is created in if not specified on the command line or
via environment variables.</p>
</td>
</tr>
<tr>
<td>
<code>eks</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.EKSConfig">
EKSConfig
</a>
</em>
</td>
<td>
<p>EKS controls the configuration related to EKS. Settings in here affect the control plane
and nodes roles</p>
</td>
</tr>
<tr>
<td>
<code>eventBridge</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.EventBridgeConfig">
EventBridgeConfig
</a>
</em>
</td>
<td>
<p>EventBridge controls configuration for consuming EventBridge events</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>secureSecretBackends</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecretBackend">
[]SecretBackend
</a>
</em>
</td>
<td>
<p>SecureSecretsBackend, when set to parameter-store will create AWS Systems Manager
Parameter Storage policies. By default or with the value of secrets-manager,
will generate AWS Secrets Manager policies instead.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">AWSIAMRoleSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ClusterAPIControllers">ClusterAPIControllers</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ControlPlane">ControlPlane</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.EKSConfig">EKSConfig</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.Nodes">Nodes</a>)
</p>
<p>
<p>AWSIAMRoleSpec defines common configuration for AWS IAM roles created by
Kubernetes Cluster API Provider AWS.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable if set to true will not create the AWS IAM role. Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>extraPolicyAttachments</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ExtraPolicyAttachments is a list of additional policies to be attached to the IAM role.</p>
</td>
</tr>
<tr>
<td>
<code>extraStatements</code><br/>
<em>
[]Cluster API AWS iam/api/v1beta1.StatementEntry
</em>
</td>
<td>
<p>ExtraStatements are additional IAM statements to be included inline for the role.</p>
</td>
</tr>
<tr>
<td>
<code>trustStatements</code><br/>
<em>
[]Cluster API AWS iam/api/v1beta1.StatementEntry
</em>
</td>
<td>
<p>TrustStatements is an IAM PolicyDocument defining what identities are allowed to assume this role.
See &ldquo;sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/iam/v1beta1&rdquo; for more documentation.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a map of tags to be applied to the AWS IAM role.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.BootstrapUser">BootstrapUser
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>BootstrapUser contains a list of elements that is specific
to the configuration and enablement of an IAM user.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enable controls whether or not a bootstrap AWS IAM user will be created.
This can be used to scope down the initial credentials used to bootstrap the
cluster.
Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>userName</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserName controls the username of the bootstrap user. Defaults to
&ldquo;bootstrapper.cluster-api-provider-aws.sigs.k8s.io&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>groupName</code><br/>
<em>
string
</em>
</td>
<td>
<p>GroupName controls the group the user will belong to. Defaults to
&ldquo;bootstrapper.cluster-api-provider-aws.sigs.k8s.io&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>extraPolicyAttachments</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ExtraPolicyAttachments is a list of additional policies to be attached to the IAM user.</p>
</td>
</tr>
<tr>
<td>
<code>extraGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ExtraGroups is a list of groups to add this user to.</p>
</td>
</tr>
<tr>
<td>
<code>extraStatements</code><br/>
<em>
[]Cluster API AWS iam/api/v1beta1.StatementEntry
</em>
</td>
<td>
<p>ExtraStatements are additional AWS IAM policy document statements to be included inline for the user.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a map of tags to be applied to the AWS IAM user.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ClusterAPIControllers">ClusterAPIControllers
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>ClusterAPIControllers controls the configuration of the AWS IAM role for
the Kubernetes Cluster API Provider AWS controller.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSIAMRoleSpec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSIAMRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>allowedEC2InstanceProfiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AllowedEC2InstanceProfiles controls which EC2 roles are allowed to be
consumed by Cluster API when creating an ec2 instance. Defaults to
*.<suffix>, where suffix is defaulted to .cluster-api-provider-aws.sigs.k8s.io</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.ControlPlane">ControlPlane
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>ControlPlane controls the configuration of the AWS IAM role for
the control plane of provisioned Kubernetes clusters.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSIAMRoleSpec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSIAMRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>disableClusterAPIControllerPolicyAttachment</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableClusterAPIControllerPolicyAttachment, if set to true, will not attach the AWS IAM policy for Cluster
API Provider AWS to the control plane role. Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>disableCloudProviderPolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableCloudProviderPolicy if set to true, will not generate and attach the AWS IAM policy for the AWS Cloud Provider.</p>
</td>
</tr>
<tr>
<td>
<code>enableCSIPolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableCSIPolicy if set to true, will generate and attach the AWS IAM policy for the EBS CSI Driver.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.EKSConfig">EKSConfig
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>EKSConfig represents the EKS related configuration config.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable controls whether EKS-related permissions are granted</p>
</td>
</tr>
<tr>
<td>
<code>iamRoleCreation</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AllowIAMRoleCreation controls whether the EKS controllers have permissions for creating IAM
roles per cluster</p>
</td>
</tr>
<tr>
<td>
<code>enableUserEKSConsolePolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableUserEKSConsolePolicy controls the creation of the policy to view EKS nodes and workloads.</p>
</td>
</tr>
<tr>
<td>
<code>defaultControlPlaneRole</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>DefaultControlPlaneRole controls the configuration of the AWS IAM role for
the EKS control plane. This is the default role that will be used if
no role is included in the spec and automatic creation of the role
isn&rsquo;t enabled</p>
</td>
</tr>
<tr>
<td>
<code>managedMachinePool</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>ManagedMachinePool controls the configuration of the AWS IAM role for
used by EKS managed machine pools.</p>
</td>
</tr>
<tr>
<td>
<code>fargate</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>Fargate controls the configuration of the AWS IAM role for
used by EKS managed machine pools.</p>
</td>
</tr>
<tr>
<td>
<code>kmsAliasPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>KMSAliasPrefix is prefix to use to restrict permission to KMS keys to only those that have an alias
name that is prefixed by this.
Defaults to cluster-api-provider-aws-*</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.EventBridgeConfig">EventBridgeConfig
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>EventBridgeConfig represents configuration for enabling experimental feature to consume
EventBridge EC2 events.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enable controls whether permissions are granted to consume EC2 events</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.Nodes">Nodes
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>Nodes controls the configuration of the AWS IAM role for worker nodes
in a cluster created by Kubernetes Cluster API Provider AWS.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSIAMRoleSpec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSIAMRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>disableCloudProviderPolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableCloudProviderPolicy if set to true, will not generate and attach the policy for the AWS Cloud Provider.
Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>ec2ContainerRegistryReadOnly</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EC2ContainerRegistryReadOnly controls whether the node has read-only access to the
EC2 container registry</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1">bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the bootstrap v1beta1 API group</p>
</p>
Resource Types:
<ul></ul>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfiguration">AWSIAMConfiguration
</h3>
<p>
<p>AWSIAMConfiguration controls the creation of AWS Identity and Access Management (IAM) resources for use
by Kubernetes clusters and Kubernetes Cluster API Provider AWS.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">
AWSIAMConfigurationSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>namePrefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NamePrefix will be prepended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to &ldquo;&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>nameSuffix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NameSuffix will be appended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to
&ldquo;.cluster-api-provider-aws.sigs.k8s.io&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlane</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ControlPlane">
ControlPlane
</a>
</em>
</td>
<td>
<p>ControlPlane controls the configuration of the AWS IAM role for a Kubernetes cluster&rsquo;s control plane nodes.</p>
</td>
</tr>
<tr>
<td>
<code>clusterAPIControllers</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ClusterAPIControllers">
ClusterAPIControllers
</a>
</em>
</td>
<td>
<p>ClusterAPIControllers controls the configuration of an IAM role and policy specifically for Kubernetes Cluster API Provider AWS.</p>
</td>
</tr>
<tr>
<td>
<code>nodes</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.Nodes">
Nodes
</a>
</em>
</td>
<td>
<p>Nodes controls the configuration of the AWS IAM role for all nodes in a Kubernetes cluster.</p>
</td>
</tr>
<tr>
<td>
<code>bootstrapUser</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.BootstrapUser">
BootstrapUser
</a>
</em>
</td>
<td>
<p>BootstrapUser contains a list of elements that is specific
to the configuration and enablement of an IAM user.</p>
</td>
</tr>
<tr>
<td>
<code>stackName</code><br/>
<em>
string
</em>
</td>
<td>
<p>StackName defines the name of the AWS CloudFormation stack.</p>
</td>
</tr>
<tr>
<td>
<code>stackTags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>StackTags defines the tags of the AWS CloudFormation stack.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>Region controls which region the control-plane is created in if not specified on the command line or
via environment variables.</p>
</td>
</tr>
<tr>
<td>
<code>eks</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.EKSConfig">
EKSConfig
</a>
</em>
</td>
<td>
<p>EKS controls the configuration related to EKS. Settings in here affect the control plane
and nodes roles</p>
</td>
</tr>
<tr>
<td>
<code>eventBridge</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.EventBridgeConfig">
EventBridgeConfig
</a>
</em>
</td>
<td>
<p>EventBridge controls configuration for consuming EventBridge events</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>secureSecretBackends</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecretBackend">
[]SecretBackend
</a>
</em>
</td>
<td>
<p>SecureSecretsBackend, when set to parameter-store will create AWS Systems Manager
Parameter Storage policies. By default or with the value of secrets-manager,
will generate AWS Secrets Manager policies instead.</p>
</td>
</tr>
<tr>
<td>
<code>s3Buckets</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.S3Buckets">
S3Buckets
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Buckets, when enabled, will add controller nodes permissions to
create S3 Buckets for workload clusters.
TODO: This field could be a pointer, but it seems it breaks setting default values?</p>
</td>
</tr>
<tr>
<td>
<code>allowAssumeRole</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AllowAssumeRole enables the sts:AssumeRole permission within the CAPA policies</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfiguration">AWSIAMConfiguration</a>)
</p>
<p>
<p>AWSIAMConfigurationSpec defines the specification of the AWSIAMConfiguration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>namePrefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NamePrefix will be prepended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to &ldquo;&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>nameSuffix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NameSuffix will be appended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to
&ldquo;.cluster-api-provider-aws.sigs.k8s.io&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlane</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ControlPlane">
ControlPlane
</a>
</em>
</td>
<td>
<p>ControlPlane controls the configuration of the AWS IAM role for a Kubernetes cluster&rsquo;s control plane nodes.</p>
</td>
</tr>
<tr>
<td>
<code>clusterAPIControllers</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ClusterAPIControllers">
ClusterAPIControllers
</a>
</em>
</td>
<td>
<p>ClusterAPIControllers controls the configuration of an IAM role and policy specifically for Kubernetes Cluster API Provider AWS.</p>
</td>
</tr>
<tr>
<td>
<code>nodes</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.Nodes">
Nodes
</a>
</em>
</td>
<td>
<p>Nodes controls the configuration of the AWS IAM role for all nodes in a Kubernetes cluster.</p>
</td>
</tr>
<tr>
<td>
<code>bootstrapUser</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.BootstrapUser">
BootstrapUser
</a>
</em>
</td>
<td>
<p>BootstrapUser contains a list of elements that is specific
to the configuration and enablement of an IAM user.</p>
</td>
</tr>
<tr>
<td>
<code>stackName</code><br/>
<em>
string
</em>
</td>
<td>
<p>StackName defines the name of the AWS CloudFormation stack.</p>
</td>
</tr>
<tr>
<td>
<code>stackTags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>StackTags defines the tags of the AWS CloudFormation stack.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>Region controls which region the control-plane is created in if not specified on the command line or
via environment variables.</p>
</td>
</tr>
<tr>
<td>
<code>eks</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.EKSConfig">
EKSConfig
</a>
</em>
</td>
<td>
<p>EKS controls the configuration related to EKS. Settings in here affect the control plane
and nodes roles</p>
</td>
</tr>
<tr>
<td>
<code>eventBridge</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.EventBridgeConfig">
EventBridgeConfig
</a>
</em>
</td>
<td>
<p>EventBridge controls configuration for consuming EventBridge events</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>secureSecretBackends</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecretBackend">
[]SecretBackend
</a>
</em>
</td>
<td>
<p>SecureSecretsBackend, when set to parameter-store will create AWS Systems Manager
Parameter Storage policies. By default or with the value of secrets-manager,
will generate AWS Secrets Manager policies instead.</p>
</td>
</tr>
<tr>
<td>
<code>s3Buckets</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.S3Buckets">
S3Buckets
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Buckets, when enabled, will add controller nodes permissions to
create S3 Buckets for workload clusters.
TODO: This field could be a pointer, but it seems it breaks setting default values?</p>
</td>
</tr>
<tr>
<td>
<code>allowAssumeRole</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AllowAssumeRole enables the sts:AssumeRole permission within the CAPA policies</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">AWSIAMRoleSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ClusterAPIControllers">ClusterAPIControllers</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ControlPlane">ControlPlane</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.EKSConfig">EKSConfig</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.Nodes">Nodes</a>)
</p>
<p>
<p>AWSIAMRoleSpec defines common configuration for AWS IAM roles created by
Kubernetes Cluster API Provider AWS.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable if set to true will not create the AWS IAM role. Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>extraPolicyAttachments</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ExtraPolicyAttachments is a list of additional policies to be attached to the IAM role.</p>
</td>
</tr>
<tr>
<td>
<code>extraStatements</code><br/>
<em>
[]Cluster API AWS iam/api/v1beta1.StatementEntry
</em>
</td>
<td>
<p>ExtraStatements are additional IAM statements to be included inline for the role.</p>
</td>
</tr>
<tr>
<td>
<code>trustStatements</code><br/>
<em>
[]Cluster API AWS iam/api/v1beta1.StatementEntry
</em>
</td>
<td>
<p>TrustStatements is an IAM PolicyDocument defining what identities are allowed to assume this role.
See &ldquo;sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/iam/v1beta1&rdquo; for more documentation.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a map of tags to be applied to the AWS IAM role.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.BootstrapUser">BootstrapUser
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>BootstrapUser contains a list of elements that is specific
to the configuration and enablement of an IAM user.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enable controls whether or not a bootstrap AWS IAM user will be created.
This can be used to scope down the initial credentials used to bootstrap the
cluster.
Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>userName</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserName controls the username of the bootstrap user. Defaults to
&ldquo;bootstrapper.cluster-api-provider-aws.sigs.k8s.io&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>groupName</code><br/>
<em>
string
</em>
</td>
<td>
<p>GroupName controls the group the user will belong to. Defaults to
&ldquo;bootstrapper.cluster-api-provider-aws.sigs.k8s.io&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>extraPolicyAttachments</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ExtraPolicyAttachments is a list of additional policies to be attached to the IAM user.</p>
</td>
</tr>
<tr>
<td>
<code>extraGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ExtraGroups is a list of groups to add this user to.</p>
</td>
</tr>
<tr>
<td>
<code>extraStatements</code><br/>
<em>
[]Cluster API AWS iam/api/v1beta1.StatementEntry
</em>
</td>
<td>
<p>ExtraStatements are additional AWS IAM policy document statements to be included inline for the user.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a map of tags to be applied to the AWS IAM user.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ClusterAPIControllers">ClusterAPIControllers
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>ClusterAPIControllers controls the configuration of the AWS IAM role for
the Kubernetes Cluster API Provider AWS controller.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSIAMRoleSpec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSIAMRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>allowedEC2InstanceProfiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AllowedEC2InstanceProfiles controls which EC2 roles are allowed to be
consumed by Cluster API when creating an ec2 instance. Defaults to
*.<suffix>, where suffix is defaulted to .cluster-api-provider-aws.sigs.k8s.io</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.ControlPlane">ControlPlane
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>ControlPlane controls the configuration of the AWS IAM role for
the control plane of provisioned Kubernetes clusters.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSIAMRoleSpec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSIAMRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>disableClusterAPIControllerPolicyAttachment</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableClusterAPIControllerPolicyAttachment, if set to true, will not attach the AWS IAM policy for Cluster
API Provider AWS to the control plane role. Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>disableCloudProviderPolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableCloudProviderPolicy if set to true, will not generate and attach the AWS IAM policy for the AWS Cloud Provider.</p>
</td>
</tr>
<tr>
<td>
<code>enableCSIPolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableCSIPolicy if set to true, will generate and attach the AWS IAM policy for the EBS CSI Driver.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.EKSConfig">EKSConfig
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>EKSConfig represents the EKS related configuration config.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable controls whether EKS-related permissions are granted</p>
</td>
</tr>
<tr>
<td>
<code>iamRoleCreation</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AllowIAMRoleCreation controls whether the EKS controllers have permissions for creating IAM
roles per cluster</p>
</td>
</tr>
<tr>
<td>
<code>enableUserEKSConsolePolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableUserEKSConsolePolicy controls the creation of the policy to view EKS nodes and workloads.</p>
</td>
</tr>
<tr>
<td>
<code>defaultControlPlaneRole</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>DefaultControlPlaneRole controls the configuration of the AWS IAM role for
the EKS control plane. This is the default role that will be used if
no role is included in the spec and automatic creation of the role
isn&rsquo;t enabled</p>
</td>
</tr>
<tr>
<td>
<code>managedMachinePool</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>ManagedMachinePool controls the configuration of the AWS IAM role for
used by EKS managed machine pools.</p>
</td>
</tr>
<tr>
<td>
<code>fargate</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>Fargate controls the configuration of the AWS IAM role for
used by EKS managed machine pools.</p>
</td>
</tr>
<tr>
<td>
<code>kmsAliasPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>KMSAliasPrefix is prefix to use to restrict permission to KMS keys to only those that have an alias
name that is prefixed by this.
Defaults to cluster-api-provider-aws-*</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.EventBridgeConfig">EventBridgeConfig
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>EventBridgeConfig represents configuration for enabling experimental feature to consume
EventBridge EC2 events.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enable controls whether permissions are granted to consume EC2 events</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.Nodes">Nodes
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>Nodes controls the configuration of the AWS IAM role for worker nodes
in a cluster created by Kubernetes Cluster API Provider AWS.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSIAMRoleSpec</code><br/>
<em>
<a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">
AWSIAMRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSIAMRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>disableCloudProviderPolicy</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableCloudProviderPolicy if set to true, will not generate and attach the policy for the AWS Cloud Provider.
Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>ec2ContainerRegistryReadOnly</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EC2ContainerRegistryReadOnly controls whether the node has read-only access to the
EC2 container registry</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.S3Buckets">S3Buckets
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>S3Buckets controls the configuration of the AWS IAM role for S3 buckets
which can be created for storing bootstrap data for nodes requiring it.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enable controls whether permissions are granted to manage S3 buckets.</p>
</td>
</tr>
<tr>
<td>
<code>namePrefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>NamePrefix will be prepended to every AWS IAM role bucket name. Defaults to &ldquo;cluster-api-provider-aws-&rdquo;.
AWSCluster S3 Bucket name must be prefixed with the same prefix.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="bootstrap.cluster.x-k8s.io/v1beta1">bootstrap.cluster.x-k8s.io/v1beta1</h2>
Resource Types:
<ul></ul>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta1.EKSConfig">EKSConfig
</h3>
<p>
<p>EKSConfig is the schema for the Amazon EKS Machine Bootstrap Configuration API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigSpec">
EKSConfigSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>kubeletExtraArgs</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script</p>
</td>
</tr>
<tr>
<td>
<code>containerRuntime</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ContainerRuntime specify the container runtime to use when bootstrapping EKS.</p>
</td>
</tr>
<tr>
<td>
<code>dnsClusterIP</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DNSClusterIP overrides the IP address to use for DNS queries within the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>dockerConfigJson</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
This is expected to be a json string.</p>
</td>
</tr>
<tr>
<td>
<code>apiRetryAttempts</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIRetryAttempts is the number of retry attempts for AWS API call.</p>
</td>
</tr>
<tr>
<td>
<code>pauseContainer</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.PauseContainer">
PauseContainer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PauseContainer allows customization of the pause container to use.</p>
</td>
</tr>
<tr>
<td>
<code>useMaxPods</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UseMaxPods  sets &ndash;max-pods for the kubelet when true.</p>
</td>
</tr>
<tr>
<td>
<code>serviceIPV6Cidr</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
the ip family will be set to ipv6.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigStatus">
EKSConfigStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigSpec">EKSConfigSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfig">EKSConfig</a>, <a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplateResource">EKSConfigTemplateResource</a>)
</p>
<p>
<p>EKSConfigSpec defines the desired state of Amazon EKS Bootstrap Configuration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>kubeletExtraArgs</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script</p>
</td>
</tr>
<tr>
<td>
<code>containerRuntime</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ContainerRuntime specify the container runtime to use when bootstrapping EKS.</p>
</td>
</tr>
<tr>
<td>
<code>dnsClusterIP</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DNSClusterIP overrides the IP address to use for DNS queries within the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>dockerConfigJson</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
This is expected to be a json string.</p>
</td>
</tr>
<tr>
<td>
<code>apiRetryAttempts</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIRetryAttempts is the number of retry attempts for AWS API call.</p>
</td>
</tr>
<tr>
<td>
<code>pauseContainer</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.PauseContainer">
PauseContainer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PauseContainer allows customization of the pause container to use.</p>
</td>
</tr>
<tr>
<td>
<code>useMaxPods</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UseMaxPods  sets &ndash;max-pods for the kubelet when true.</p>
</td>
</tr>
<tr>
<td>
<code>serviceIPV6Cidr</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
the ip family will be set to ipv6.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigStatus">EKSConfigStatus
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfig">EKSConfig</a>)
</p>
<p>
<p>EKSConfigStatus defines the observed state of the Amazon EKS Bootstrap Configuration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready indicates the BootstrapData secret is ready to be consumed</p>
</td>
</tr>
<tr>
<td>
<code>dataSecretName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DataSecretName is the name of the secret that stores the bootstrap data script.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set on non-retryable errors</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set on non-retryable errors</p>
</td>
</tr>
<tr>
<td>
<code>observedGeneration</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>ObservedGeneration is the latest generation observed by the controller.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the EKSConfig.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplate">EKSConfigTemplate
</h3>
<p>
<p>EKSConfigTemplate is the Amazon EKS Bootstrap Configuration Template API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplateSpec">
EKSConfigTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplateResource">
EKSConfigTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplateResource">EKSConfigTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplateSpec">EKSConfigTemplateSpec</a>)
</p>
<p>
<p>EKSConfigTemplateResource defines the Template structure.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigSpec">
EKSConfigSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>kubeletExtraArgs</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script</p>
</td>
</tr>
<tr>
<td>
<code>containerRuntime</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ContainerRuntime specify the container runtime to use when bootstrapping EKS.</p>
</td>
</tr>
<tr>
<td>
<code>dnsClusterIP</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DNSClusterIP overrides the IP address to use for DNS queries within the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>dockerConfigJson</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
This is expected to be a json string.</p>
</td>
</tr>
<tr>
<td>
<code>apiRetryAttempts</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIRetryAttempts is the number of retry attempts for AWS API call.</p>
</td>
</tr>
<tr>
<td>
<code>pauseContainer</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.PauseContainer">
PauseContainer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PauseContainer allows customization of the pause container to use.</p>
</td>
</tr>
<tr>
<td>
<code>useMaxPods</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UseMaxPods  sets &ndash;max-pods for the kubelet when true.</p>
</td>
</tr>
<tr>
<td>
<code>serviceIPV6Cidr</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
the ip family will be set to ipv6.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplateSpec">EKSConfigTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplate">EKSConfigTemplate</a>)
</p>
<p>
<p>EKSConfigTemplateSpec defines the desired state of templated EKSConfig Amazon EKS Bootstrap Configuration resources.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigTemplateResource">
EKSConfigTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta1.PauseContainer">PauseContainer
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta1.EKSConfigSpec">EKSConfigSpec</a>)
</p>
<p>
<p>PauseContainer contains details of pause container.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>accountNumber</code><br/>
<em>
string
</em>
</td>
<td>
<p>AccountNumber is the AWS account number to pull the pause container from.</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version is the tag of the pause container to use.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="bootstrap.cluster.x-k8s.io/v1beta2">bootstrap.cluster.x-k8s.io/v1beta2</h2>
<p>
<p>Package v1beta2 contains API Schema definitions for the Amazon EKS Bootstrap v1beta2 API group.</p>
</p>
Resource Types:
<ul></ul>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.DiskSetup">DiskSetup
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">EKSConfigSpec</a>)
</p>
<p>
<p>DiskSetup defines input for generated disk_setup and fs_setup in cloud-init.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>partitions</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.Partition">
[]Partition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Partitions specifies the list of the partitions to setup.</p>
</td>
</tr>
<tr>
<td>
<code>filesystems</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.Filesystem">
[]Filesystem
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Filesystems specifies the list of file systems to setup.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.EKSConfig">EKSConfig
</h3>
<p>
<p>EKSConfig is the schema for the Amazon EKS Machine Bootstrap Configuration API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">
EKSConfigSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>kubeletExtraArgs</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script</p>
</td>
</tr>
<tr>
<td>
<code>containerRuntime</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ContainerRuntime specify the container runtime to use when bootstrapping EKS.</p>
</td>
</tr>
<tr>
<td>
<code>dnsClusterIP</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DNSClusterIP overrides the IP address to use for DNS queries within the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>dockerConfigJson</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
This is expected to be a json string.</p>
</td>
</tr>
<tr>
<td>
<code>apiRetryAttempts</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIRetryAttempts is the number of retry attempts for AWS API call.</p>
</td>
</tr>
<tr>
<td>
<code>pauseContainer</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.PauseContainer">
PauseContainer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PauseContainer allows customization of the pause container to use.</p>
</td>
</tr>
<tr>
<td>
<code>useMaxPods</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UseMaxPods  sets &ndash;max-pods for the kubelet when true.</p>
</td>
</tr>
<tr>
<td>
<code>serviceIPV6Cidr</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
the ip family will be set to ipv6.</p>
</td>
</tr>
<tr>
<td>
<code>preBootstrapCommands</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PreBootstrapCommands specifies extra commands to run before bootstrapping nodes to the cluster</p>
</td>
</tr>
<tr>
<td>
<code>postBootstrapCommands</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PostBootstrapCommands specifies extra commands to run after bootstrapping nodes to the cluster</p>
</td>
</tr>
<tr>
<td>
<code>boostrapCommandOverride</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>BootstrapCommandOverride allows you to override the bootstrap command to use for EKS nodes.</p>
</td>
</tr>
<tr>
<td>
<code>files</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.File">
[]File
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Files specifies extra files to be passed to user_data upon creation.</p>
</td>
</tr>
<tr>
<td>
<code>diskSetup</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.DiskSetup">
DiskSetup
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>DiskSetup specifies options for the creation of partition tables and file systems on devices.</p>
</td>
</tr>
<tr>
<td>
<code>mounts</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.MountPoints">
[]MountPoints
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Mounts specifies a list of mount points to be setup.</p>
</td>
</tr>
<tr>
<td>
<code>users</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.User">
[]User
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Users specifies extra users to add</p>
</td>
</tr>
<tr>
<td>
<code>ntp</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.NTP">
NTP
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NTP specifies NTP configuration</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigStatus">
EKSConfigStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">EKSConfigSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfig">EKSConfig</a>, <a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplateResource">EKSConfigTemplateResource</a>)
</p>
<p>
<p>EKSConfigSpec defines the desired state of Amazon EKS Bootstrap Configuration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>kubeletExtraArgs</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script</p>
</td>
</tr>
<tr>
<td>
<code>containerRuntime</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ContainerRuntime specify the container runtime to use when bootstrapping EKS.</p>
</td>
</tr>
<tr>
<td>
<code>dnsClusterIP</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DNSClusterIP overrides the IP address to use for DNS queries within the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>dockerConfigJson</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
This is expected to be a json string.</p>
</td>
</tr>
<tr>
<td>
<code>apiRetryAttempts</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIRetryAttempts is the number of retry attempts for AWS API call.</p>
</td>
</tr>
<tr>
<td>
<code>pauseContainer</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.PauseContainer">
PauseContainer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PauseContainer allows customization of the pause container to use.</p>
</td>
</tr>
<tr>
<td>
<code>useMaxPods</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UseMaxPods  sets &ndash;max-pods for the kubelet when true.</p>
</td>
</tr>
<tr>
<td>
<code>serviceIPV6Cidr</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
the ip family will be set to ipv6.</p>
</td>
</tr>
<tr>
<td>
<code>preBootstrapCommands</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PreBootstrapCommands specifies extra commands to run before bootstrapping nodes to the cluster</p>
</td>
</tr>
<tr>
<td>
<code>postBootstrapCommands</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PostBootstrapCommands specifies extra commands to run after bootstrapping nodes to the cluster</p>
</td>
</tr>
<tr>
<td>
<code>boostrapCommandOverride</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>BootstrapCommandOverride allows you to override the bootstrap command to use for EKS nodes.</p>
</td>
</tr>
<tr>
<td>
<code>files</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.File">
[]File
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Files specifies extra files to be passed to user_data upon creation.</p>
</td>
</tr>
<tr>
<td>
<code>diskSetup</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.DiskSetup">
DiskSetup
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>DiskSetup specifies options for the creation of partition tables and file systems on devices.</p>
</td>
</tr>
<tr>
<td>
<code>mounts</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.MountPoints">
[]MountPoints
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Mounts specifies a list of mount points to be setup.</p>
</td>
</tr>
<tr>
<td>
<code>users</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.User">
[]User
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Users specifies extra users to add</p>
</td>
</tr>
<tr>
<td>
<code>ntp</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.NTP">
NTP
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NTP specifies NTP configuration</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigStatus">EKSConfigStatus
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfig">EKSConfig</a>)
</p>
<p>
<p>EKSConfigStatus defines the observed state of the Amazon EKS Bootstrap Configuration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready indicates the BootstrapData secret is ready to be consumed</p>
</td>
</tr>
<tr>
<td>
<code>dataSecretName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DataSecretName is the name of the secret that stores the bootstrap data script.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set on non-retryable errors</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set on non-retryable errors</p>
</td>
</tr>
<tr>
<td>
<code>observedGeneration</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>ObservedGeneration is the latest generation observed by the controller.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the EKSConfig.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplate">EKSConfigTemplate
</h3>
<p>
<p>EKSConfigTemplate is the Amazon EKS Bootstrap Configuration Template API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplateSpec">
EKSConfigTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplateResource">
EKSConfigTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplateResource">EKSConfigTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplateSpec">EKSConfigTemplateSpec</a>)
</p>
<p>
<p>EKSConfigTemplateResource defines the Template structure.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">
EKSConfigSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>kubeletExtraArgs</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script</p>
</td>
</tr>
<tr>
<td>
<code>containerRuntime</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ContainerRuntime specify the container runtime to use when bootstrapping EKS.</p>
</td>
</tr>
<tr>
<td>
<code>dnsClusterIP</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DNSClusterIP overrides the IP address to use for DNS queries within the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>dockerConfigJson</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
This is expected to be a json string.</p>
</td>
</tr>
<tr>
<td>
<code>apiRetryAttempts</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>APIRetryAttempts is the number of retry attempts for AWS API call.</p>
</td>
</tr>
<tr>
<td>
<code>pauseContainer</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.PauseContainer">
PauseContainer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PauseContainer allows customization of the pause container to use.</p>
</td>
</tr>
<tr>
<td>
<code>useMaxPods</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UseMaxPods  sets &ndash;max-pods for the kubelet when true.</p>
</td>
</tr>
<tr>
<td>
<code>serviceIPV6Cidr</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
the ip family will be set to ipv6.</p>
</td>
</tr>
<tr>
<td>
<code>preBootstrapCommands</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PreBootstrapCommands specifies extra commands to run before bootstrapping nodes to the cluster</p>
</td>
</tr>
<tr>
<td>
<code>postBootstrapCommands</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PostBootstrapCommands specifies extra commands to run after bootstrapping nodes to the cluster</p>
</td>
</tr>
<tr>
<td>
<code>boostrapCommandOverride</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>BootstrapCommandOverride allows you to override the bootstrap command to use for EKS nodes.</p>
</td>
</tr>
<tr>
<td>
<code>files</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.File">
[]File
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Files specifies extra files to be passed to user_data upon creation.</p>
</td>
</tr>
<tr>
<td>
<code>diskSetup</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.DiskSetup">
DiskSetup
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>DiskSetup specifies options for the creation of partition tables and file systems on devices.</p>
</td>
</tr>
<tr>
<td>
<code>mounts</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.MountPoints">
[]MountPoints
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Mounts specifies a list of mount points to be setup.</p>
</td>
</tr>
<tr>
<td>
<code>users</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.User">
[]User
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Users specifies extra users to add</p>
</td>
</tr>
<tr>
<td>
<code>ntp</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.NTP">
NTP
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NTP specifies NTP configuration</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplateSpec">EKSConfigTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplate">EKSConfigTemplate</a>)
</p>
<p>
<p>EKSConfigTemplateSpec defines the desired state of templated EKSConfig Amazon EKS Bootstrap Configuration resources.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigTemplateResource">
EKSConfigTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.Encoding">Encoding
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.File">File</a>)
</p>
<p>
<p>Encoding specifies the cloud-init file encoding.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;base64&#34;</p></td>
<td><p>Base64 implies the contents of the file are encoded as base64.</p>
</td>
</tr><tr><td><p>&#34;gzip&#34;</p></td>
<td><p>Gzip implies the contents of the file are encoded with gzip.</p>
</td>
</tr><tr><td><p>&#34;gzip&#43;base64&#34;</p></td>
<td><p>GzipBase64 implies the contents of the file are first base64 encoded and then gzip encoded.</p>
</td>
</tr></tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.File">File
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">EKSConfigSpec</a>)
</p>
<p>
<p>File defines the input for generating write_files in cloud-init.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>path</code><br/>
<em>
string
</em>
</td>
<td>
<p>Path specifies the full path on disk where to store the file.</p>
</td>
</tr>
<tr>
<td>
<code>owner</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Owner specifies the ownership of the file, e.g. &ldquo;root:root&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>permissions</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Permissions specifies the permissions to assign to the file, e.g. &ldquo;0640&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>encoding</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.Encoding">
Encoding
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Encoding specifies the encoding of the file contents.</p>
</td>
</tr>
<tr>
<td>
<code>append</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Append specifies whether to append Content to existing file if Path exists.</p>
</td>
</tr>
<tr>
<td>
<code>content</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Content is the actual content of the file.</p>
</td>
</tr>
<tr>
<td>
<code>contentFrom</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.FileSource">
FileSource
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ContentFrom is a referenced source of content to populate the file.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.FileSource">FileSource
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.File">File</a>)
</p>
<p>
<p>FileSource is a union of all possible external source types for file data.
Only one field may be populated in any given instance. Developers adding new
sources of data for target systems should add them here.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>secret</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.SecretFileSource">
SecretFileSource
</a>
</em>
</td>
<td>
<p>Secret represents a secret that should populate this file.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.Filesystem">Filesystem
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.DiskSetup">DiskSetup</a>)
</p>
<p>
<p>Filesystem defines the file systems to be created.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>device</code><br/>
<em>
string
</em>
</td>
<td>
<p>Device specifies the device name</p>
</td>
</tr>
<tr>
<td>
<code>filesystem</code><br/>
<em>
string
</em>
</td>
<td>
<p>Filesystem specifies the file system type.</p>
</td>
</tr>
<tr>
<td>
<code>label</code><br/>
<em>
string
</em>
</td>
<td>
<p>Label specifies the file system label to be used. If set to None, no label is used.</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Partition specifies the partition to use. The valid options are: &ldquo;auto|any&rdquo;, &ldquo;auto&rdquo;, &ldquo;any&rdquo;, &ldquo;none&rdquo;, and <NUM>, where NUM is the actual partition number.</p>
</td>
</tr>
<tr>
<td>
<code>overwrite</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Overwrite defines whether or not to overwrite any existing filesystem.
If true, any pre-existing file system will be destroyed. Use with Caution.</p>
</td>
</tr>
<tr>
<td>
<code>extraOpts</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ExtraOpts defined extra options to add to the command for creating the file system.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.MountPoints">MountPoints
(<code>[]string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">EKSConfigSpec</a>)
</p>
<p>
<p>MountPoints defines input for generated mounts in cloud-init.</p>
</p>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.NTP">NTP
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">EKSConfigSpec</a>)
</p>
<p>
<p>NTP defines input for generated ntp in cloud-init.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>servers</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Servers specifies which NTP servers to use</p>
</td>
</tr>
<tr>
<td>
<code>enabled</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Enabled specifies whether NTP should be enabled</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.Partition">Partition
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.DiskSetup">DiskSetup</a>)
</p>
<p>
<p>Partition defines how to create and layout a partition.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>device</code><br/>
<em>
string
</em>
</td>
<td>
<p>Device is the name of the device.</p>
</td>
</tr>
<tr>
<td>
<code>layout</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Layout specifies the device layout.
If it is true, a single partition will be created for the entire device.
When layout is false, it means don&rsquo;t partition or ignore existing partitioning.</p>
</td>
</tr>
<tr>
<td>
<code>overwrite</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Overwrite describes whether to skip checks and create the partition if a partition or filesystem is found on the device.
Use with caution. Default is &lsquo;false&rsquo;.</p>
</td>
</tr>
<tr>
<td>
<code>tableType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>TableType specifies the tupe of partition table. The following are supported:
&lsquo;mbr&rsquo;: default and setups a MS-DOS partition table
&lsquo;gpt&rsquo;: setups a GPT partition table</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.PasswdSource">PasswdSource
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.User">User</a>)
</p>
<p>
<p>PasswdSource is a union of all possible external source types for passwd data.
Only one field may be populated in any given instance. Developers adding new
sources of data for target systems should add them here.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>secret</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.SecretPasswdSource">
SecretPasswdSource
</a>
</em>
</td>
<td>
<p>Secret represents a secret that should populate this password.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.PauseContainer">PauseContainer
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">EKSConfigSpec</a>)
</p>
<p>
<p>PauseContainer contains details of pause container.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>accountNumber</code><br/>
<em>
string
</em>
</td>
<td>
<p>AccountNumber is the AWS account number to pull the pause container from.</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version is the tag of the pause container to use.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.SecretFileSource">SecretFileSource
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.FileSource">FileSource</a>)
</p>
<p>
<p>SecretFileSource adapts a Secret into a FileSource.</p>
<p>The contents of the target Secret&rsquo;s Data field will be presented
as files using the keys in the Data field as the file names.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the secret in the KubeadmBootstrapConfig&rsquo;s namespace to use.</p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>Key is the key in the secret&rsquo;s data map for this value.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.SecretPasswdSource">SecretPasswdSource
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.PasswdSource">PasswdSource</a>)
</p>
<p>
<p>SecretPasswdSource adapts a Secret into a PasswdSource.</p>
<p>The contents of the target Secret&rsquo;s Data field will be presented
as passwd using the keys in the Data field as the file names.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the secret in the KubeadmBootstrapConfig&rsquo;s namespace to use.</p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>Key is the key in the secret&rsquo;s data map for this value.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="bootstrap.cluster.x-k8s.io/v1beta2.User">User
</h3>
<p>
(<em>Appears on:</em><a href="#bootstrap.cluster.x-k8s.io/v1beta2.EKSConfigSpec">EKSConfigSpec</a>)
</p>
<p>
<p>User defines the input for a generated user in cloud-init.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name specifies the username</p>
</td>
</tr>
<tr>
<td>
<code>gecos</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Gecos specifies the gecos to use for the user</p>
</td>
</tr>
<tr>
<td>
<code>groups</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Groups specifies the additional groups for the user</p>
</td>
</tr>
<tr>
<td>
<code>homeDir</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>HomeDir specifies the home directory to use for the user</p>
</td>
</tr>
<tr>
<td>
<code>inactive</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Inactive specifies whether to mark the user as inactive</p>
</td>
</tr>
<tr>
<td>
<code>shell</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Shell specifies the user&rsquo;s shell</p>
</td>
</tr>
<tr>
<td>
<code>passwd</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Passwd specifies a hashed password for the user</p>
</td>
</tr>
<tr>
<td>
<code>passwdFrom</code><br/>
<em>
<a href="#bootstrap.cluster.x-k8s.io/v1beta2.PasswdSource">
PasswdSource
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PasswdFrom is a referenced source of passwd to populate the passwd.</p>
</td>
</tr>
<tr>
<td>
<code>primaryGroup</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrimaryGroup specifies the primary group for the user</p>
</td>
</tr>
<tr>
<td>
<code>lockPassword</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>LockPassword specifies if password login should be disabled</p>
</td>
</tr>
<tr>
<td>
<code>sudo</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Sudo specifies a sudo role for the user</p>
</td>
</tr>
<tr>
<td>
<code>sshAuthorizedKeys</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHAuthorizedKeys specifies a list of ssh authorized keys for the user</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="controlplane.cluster.x-k8s.io/v1beta1">controlplane.cluster.x-k8s.io/v1beta1</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the controlplane v1beta1 API group</p>
</p>
Resource Types:
<ul></ul>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlane">AWSManagedControlPlane
</h3>
<p>
<p>AWSManagedControlPlane is the schema for the Amazon EKS Managed Control Plane API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">
AWSManagedControlPlaneSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>eksClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSClusterName allows you to specify the name of the EKS cluster in
AWS. If you don&rsquo;t specify a name then a default name will be created
based on the namespace and name of the managed control plane.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryCidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines the desired Kubernetes version. If no version number
is supplied then the latest version of Kubernetes that EKS supports
will be used.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role that gives EKS
permission to make API calls. If the role is pre-existing
we will treat it as unmanaged and not delete it on
deletion. If the EKSEnableIAM feature flag is true
and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the control plane role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>logging</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.ControlPlaneLoggingSpec">
ControlPlaneLoggingSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Logging specifies which EKS Cluster logs should be enabled. Entries for
each of the enabled logs will be sent to CloudWatch</p>
</td>
</tr>
<tr>
<td>
<code>encryptionConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.EncryptionConfig">
EncryptionConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptionConfig specifies the encryption configuration for the cluster</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>iamAuthenticatorConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.IAMAuthenticatorConfig">
IAMAuthenticatorConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMAuthenticatorConfig allows the specification of any additional user or role mappings
for use when generating the aws-iam-authenticator configuration. If this is nil the
default configuration is still generated for the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>endpointAccess</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.EndpointAccess">
EndpointAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Endpoints specifies access to this cluster&rsquo;s control plane endpoints</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>tokenMethod</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.EKSTokenMethod">
EKSTokenMethod
</a>
</em>
</td>
<td>
<p>TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
iam-authenticator - obtains a client token using iam-authentictor
aws-cli - obtains a client token using the AWS CLI
Defaults to iam-authenticator</p>
</td>
</tr>
<tr>
<td>
<code>associateOIDCProvider</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AssociateOIDCProvider can be enabled to automatically create an identity
provider for the controller for use with IAM roles for service accounts</p>
</td>
</tr>
<tr>
<td>
<code>addons</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta1.Addon">
[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta1.Addon
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Addons defines the EKS addons to enable with the EKS cluster.</p>
</td>
</tr>
<tr>
<td>
<code>oidcIdentityProviderConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.OIDCIdentityProviderConfig">
OIDCIdentityProviderConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityProviderconfig is used to specify the oidc provider config
to be attached with this eks cluster</p>
</td>
</tr>
<tr>
<td>
<code>disableVPCCNI</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableVPCCNI indicates that the Amazon VPC CNI should be disabled. With EKS clusters the
Amazon VPC CNI is automatically installed into the cluster. For clusters where you want
to use an alternate CNI this option provides a way to specify that the Amazon VPC CNI
should be deleted. You cannot set this to true if you are using the
Amazon VPC CNI addon.</p>
</td>
</tr>
<tr>
<td>
<code>vpcCni</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.VpcCni">
VpcCni
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>VpcCni is used to set configuration options for the VPC CNI plugin</p>
</td>
</tr>
<tr>
<td>
<code>kubeProxy</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.KubeProxy">
KubeProxy
</a>
</em>
</td>
<td>
<p>KubeProxy defines managed attributes of the kube-proxy daemonset</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneStatus">
AWSManagedControlPlaneStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlane">AWSManagedControlPlane</a>)
</p>
<p>
<p>AWSManagedControlPlaneSpec defines the desired state of an Amazon EKS Cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>eksClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSClusterName allows you to specify the name of the EKS cluster in
AWS. If you don&rsquo;t specify a name then a default name will be created
based on the namespace and name of the managed control plane.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryCidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines the desired Kubernetes version. If no version number
is supplied then the latest version of Kubernetes that EKS supports
will be used.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role that gives EKS
permission to make API calls. If the role is pre-existing
we will treat it as unmanaged and not delete it on
deletion. If the EKSEnableIAM feature flag is true
and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the control plane role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>logging</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.ControlPlaneLoggingSpec">
ControlPlaneLoggingSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Logging specifies which EKS Cluster logs should be enabled. Entries for
each of the enabled logs will be sent to CloudWatch</p>
</td>
</tr>
<tr>
<td>
<code>encryptionConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.EncryptionConfig">
EncryptionConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptionConfig specifies the encryption configuration for the cluster</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>iamAuthenticatorConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.IAMAuthenticatorConfig">
IAMAuthenticatorConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMAuthenticatorConfig allows the specification of any additional user or role mappings
for use when generating the aws-iam-authenticator configuration. If this is nil the
default configuration is still generated for the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>endpointAccess</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.EndpointAccess">
EndpointAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Endpoints specifies access to this cluster&rsquo;s control plane endpoints</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>tokenMethod</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.EKSTokenMethod">
EKSTokenMethod
</a>
</em>
</td>
<td>
<p>TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
iam-authenticator - obtains a client token using iam-authentictor
aws-cli - obtains a client token using the AWS CLI
Defaults to iam-authenticator</p>
</td>
</tr>
<tr>
<td>
<code>associateOIDCProvider</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AssociateOIDCProvider can be enabled to automatically create an identity
provider for the controller for use with IAM roles for service accounts</p>
</td>
</tr>
<tr>
<td>
<code>addons</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta1.Addon">
[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta1.Addon
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Addons defines the EKS addons to enable with the EKS cluster.</p>
</td>
</tr>
<tr>
<td>
<code>oidcIdentityProviderConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.OIDCIdentityProviderConfig">
OIDCIdentityProviderConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityProviderconfig is used to specify the oidc provider config
to be attached with this eks cluster</p>
</td>
</tr>
<tr>
<td>
<code>disableVPCCNI</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableVPCCNI indicates that the Amazon VPC CNI should be disabled. With EKS clusters the
Amazon VPC CNI is automatically installed into the cluster. For clusters where you want
to use an alternate CNI this option provides a way to specify that the Amazon VPC CNI
should be deleted. You cannot set this to true if you are using the
Amazon VPC CNI addon.</p>
</td>
</tr>
<tr>
<td>
<code>vpcCni</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.VpcCni">
VpcCni
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>VpcCni is used to set configuration options for the VPC CNI plugin</p>
</td>
</tr>
<tr>
<td>
<code>kubeProxy</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.KubeProxy">
KubeProxy
</a>
</em>
</td>
<td>
<p>KubeProxy defines managed attributes of the kube-proxy daemonset</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlane">AWSManagedControlPlane</a>)
</p>
<p>
<p>AWSManagedControlPlaneStatus defines the observed state of an Amazon EKS Cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>networkStatus</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkStatus">
NetworkStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Networks holds details about the AWS networking resources used by the control plane</p>
</td>
</tr>
<tr>
<td>
<code>failureDomains</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.FailureDomains
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureDomains specifies a list fo available availability zones that can be used</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">
Instance
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion holds details of the instance that is used as a bastion jump box</p>
</td>
</tr>
<tr>
<td>
<code>oidcProvider</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.OIDCProviderStatus">
OIDCProviderStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OIDCProvider holds the status of the identity provider for this cluster</p>
</td>
</tr>
<tr>
<td>
<code>externalManagedControlPlane</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ExternalManagedControlPlane indicates to cluster-api that the control plane
is managed by an external service such as AKS, EKS, GKE, etc.</p>
</td>
</tr>
<tr>
<td>
<code>initialized</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Initialized denotes whether or not the control plane has the
uploaded kubernetes config-map.</p>
</td>
</tr>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the AWSManagedControlPlane API Server is ready to
receive requests and that the VPC infra is ready.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ErrorMessage indicates that there is a terminal problem reconciling the
state, and will be set to a descriptive error message.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<p>Conditions specifies the cpnditions for the managed control plane</p>
</td>
</tr>
<tr>
<td>
<code>addons</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.AddonState">
[]AddonState
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Addons holds the current status of the EKS addons</p>
</td>
</tr>
<tr>
<td>
<code>identityProviderStatus</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.IdentityProviderStatus">
IdentityProviderStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityProviderStatus holds the status for
associated identity provider</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.Addon">Addon
</h3>
<p>
<p>Addon represents a EKS addon.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the addon</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version is the version of the addon to use</p>
</td>
</tr>
<tr>
<td>
<code>configuration</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration of the EKS addon</p>
</td>
</tr>
<tr>
<td>
<code>conflictResolution</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.AddonResolution">
AddonResolution
</a>
</em>
</td>
<td>
<p>ConflictResolution is used to declare what should happen if there
are parameter conflicts. Defaults to none</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountRoleArn is the ARN of an IAM role to bind to the addons service account</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.AddonIssue">AddonIssue
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AddonState">AddonState</a>)
</p>
<p>
<p>AddonIssue represents an issue with an addon.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>code</code><br/>
<em>
string
</em>
</td>
<td>
<p>Code is the issue code</p>
</td>
</tr>
<tr>
<td>
<code>message</code><br/>
<em>
string
</em>
</td>
<td>
<p>Message is the textual description of the issue</p>
</td>
</tr>
<tr>
<td>
<code>resourceIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ResourceIDs is a list of resource ids for the issue</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.AddonResolution">AddonResolution
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.Addon">Addon</a>)
</p>
<p>
<p>AddonResolution defines the method for resolving parameter conflicts.</p>
</p>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.AddonState">AddonState
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>)
</p>
<p>
<p>AddonState represents the state of an addon.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the addon</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version is the version of the addon to use</p>
</td>
</tr>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<p>ARN is the AWS ARN of the addon</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>ServiceAccountRoleArn is the ARN of the IAM role used for the service account</p>
</td>
</tr>
<tr>
<td>
<code>createdAt</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>CreatedAt is the date and time the addon was created at</p>
</td>
</tr>
<tr>
<td>
<code>modifiedAt</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>ModifiedAt is the date and time the addon was last modified</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
string
</em>
</td>
<td>
<p>Status is the status of the addon</p>
</td>
</tr>
<tr>
<td>
<code>issues</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.AddonIssue">
[]AddonIssue
</a>
</em>
</td>
<td>
<p>Issues is a list of issue associated with the addon</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.AddonStatus">AddonStatus
(<code>string</code> alias)</p></h3>
<p>
<p>AddonStatus defines the status for an addon.</p>
</p>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.ControlPlaneLoggingSpec">ControlPlaneLoggingSpec
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>ControlPlaneLoggingSpec defines what EKS control plane logs that should be enabled.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiServer</code><br/>
<em>
bool
</em>
</td>
<td>
<p>APIServer indicates if the Kubernetes API Server log (kube-apiserver) shoulkd be enabled</p>
</td>
</tr>
<tr>
<td>
<code>audit</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Audit indicates if the Kubernetes API audit log should be enabled</p>
</td>
</tr>
<tr>
<td>
<code>authenticator</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Authenticator indicates if the iam authenticator log should be enabled</p>
</td>
</tr>
<tr>
<td>
<code>controllerManager</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ControllerManager indicates if the controller manager (kube-controller-manager) log should be enabled</p>
</td>
</tr>
<tr>
<td>
<code>scheduler</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Scheduler indicates if the Kubernetes scheduler (kube-scheduler) log should be enabled</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.EKSTokenMethod">EKSTokenMethod
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>EKSTokenMethod defines the method for obtaining a client token to use when connecting to EKS.</p>
</p>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.EncryptionConfig">EncryptionConfig
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>EncryptionConfig specifies the encryption configuration for the EKS clsuter.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>provider</code><br/>
<em>
string
</em>
</td>
<td>
<p>Provider specifies the ARN or alias of the CMK (in AWS KMS)</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
[]*string
</em>
</td>
<td>
<p>Resources specifies the resources to be encrypted</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.EndpointAccess">EndpointAccess
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>EndpointAccess specifies how control plane endpoints are accessible.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>public</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Public controls whether control plane endpoints are publicly accessible</p>
</td>
</tr>
<tr>
<td>
<code>publicCIDRs</code><br/>
<em>
[]*string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicCIDRs specifies which blocks can access the public endpoint</p>
</td>
</tr>
<tr>
<td>
<code>private</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Private points VPC-internal control plane access to the private endpoint</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.IAMAuthenticatorConfig">IAMAuthenticatorConfig
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>IAMAuthenticatorConfig represents an aws-iam-authenticator configuration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mapRoles</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.RoleMapping">
[]RoleMapping
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleMappings is a list of role mappings</p>
</td>
</tr>
<tr>
<td>
<code>mapUsers</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.UserMapping">
[]UserMapping
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>UserMappings is a list of user mappings</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.IdentityProviderStatus">IdentityProviderStatus
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>)
</p>
<p>
<p>IdentityProviderStatus holds the status for associated identity provider</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<p>ARN holds the ARN of associated identity provider</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
string
</em>
</td>
<td>
<p>Status holds current status of associated identity provider</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.KubeProxy">KubeProxy
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>KubeProxy specifies how the kube-proxy daemonset is managed.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable set to true indicates that kube-proxy should be disabled. With EKS clusters
kube-proxy is automatically installed into the cluster. For clusters where you want
to use kube-proxy functionality that is provided with an alternate CNI, this option
provides a way to specify that the kube-proxy daemonset should be deleted. You cannot
set this to true if you are using the Amazon kube-proxy addon.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.KubernetesMapping">KubernetesMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.RoleMapping">RoleMapping</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.UserMapping">UserMapping</a>)
</p>
<p>
<p>KubernetesMapping represents the kubernetes RBAC mapping.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>username</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserName is a kubernetes RBAC user subject</p>
</td>
</tr>
<tr>
<td>
<code>groups</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Groups is a list of kubernetes RBAC groups</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.OIDCIdentityProviderConfig">OIDCIdentityProviderConfig
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>OIDCIdentityProviderConfig defines the configuration for an OIDC identity provider.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>clientId</code><br/>
<em>
string
</em>
</td>
<td>
<p>This is also known as audience. The ID for the client application that makes
authentication requests to the OpenID identity provider.</p>
</td>
</tr>
<tr>
<td>
<code>groupsClaim</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The JWT claim that the provider uses to return your groups.</p>
</td>
</tr>
<tr>
<td>
<code>groupsPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The prefix that is prepended to group claims to prevent clashes with existing
names (such as system: groups). For example, the valueoidc: will create group
names like oidc:engineering and oidc:infra.</p>
</td>
</tr>
<tr>
<td>
<code>identityProviderConfigName</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the OIDC provider configuration.</p>
<p>IdentityProviderConfigName is a required field</p>
</td>
</tr>
<tr>
<td>
<code>issuerUrl</code><br/>
<em>
string
</em>
</td>
<td>
<p>The URL of the OpenID identity provider that allows the API server to discover
public signing keys for verifying tokens. The URL must begin with https://
and should correspond to the iss claim in the provider&rsquo;s OIDC ID tokens.
Per the OIDC standard, path components are allowed but query parameters are
not. Typically the URL consists of only a hostname, like <a href="https://server.example.org">https://server.example.org</a>
or <a href="https://example.com">https://example.com</a>. This URL should point to the level below .well-known/openid-configuration
and must be publicly accessible over the internet.</p>
</td>
</tr>
<tr>
<td>
<code>requiredClaims</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The key value pairs that describe required claims in the identity token.
If set, each claim is verified to be present in the token with a matching
value. For the maximum number of claims that you can require, see Amazon
EKS service quotas (<a href="https://docs.aws.amazon.com/eks/latest/userguide/service-quotas.html">https://docs.aws.amazon.com/eks/latest/userguide/service-quotas.html</a>)
in the Amazon EKS User Guide.</p>
</td>
</tr>
<tr>
<td>
<code>usernameClaim</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The JSON Web Token (JWT) claim to use as the username. The default is sub,
which is expected to be a unique identifier of the end user. You can choose
other claims, such as email or name, depending on the OpenID identity provider.
Claims other than email are prefixed with the issuer URL to prevent naming
clashes with other plug-ins.</p>
</td>
</tr>
<tr>
<td>
<code>usernamePrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The prefix that is prepended to username claims to prevent clashes with existing
names. If you do not provide this field, and username is a value other than
email, the prefix defaults to issuerurl#. You can use the value - to disable
all prefixing.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>tags to apply to oidc identity provider association</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.OIDCProviderStatus">OIDCProviderStatus
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>)
</p>
<p>
<p>OIDCProviderStatus holds the status of the AWS OIDC identity provider.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<p>ARN holds the ARN of the provider</p>
</td>
</tr>
<tr>
<td>
<code>trustPolicy</code><br/>
<em>
string
</em>
</td>
<td>
<p>TrustPolicy contains the boilerplate IAM trust policy to use for IRSA</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.RoleMapping">RoleMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.IAMAuthenticatorConfig">IAMAuthenticatorConfig</a>)
</p>
<p>
<p>RoleMapping represents a mapping from a IAM role to Kubernetes users and groups.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>rolearn</code><br/>
<em>
string
</em>
</td>
<td>
<p>RoleARN is the AWS ARN for the role to map</p>
</td>
</tr>
<tr>
<td>
<code>KubernetesMapping</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.KubernetesMapping">
KubernetesMapping
</a>
</em>
</td>
<td>
<p>
(Members of <code>KubernetesMapping</code> are embedded into this type.)
</p>
<p>KubernetesMapping holds the RBAC details for the mapping</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.UserMapping">UserMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.IAMAuthenticatorConfig">IAMAuthenticatorConfig</a>)
</p>
<p>
<p>UserMapping represents a mapping from an IAM user to Kubernetes users and groups.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>userarn</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserARN is the AWS ARN for the user to map</p>
</td>
</tr>
<tr>
<td>
<code>KubernetesMapping</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta1.KubernetesMapping">
KubernetesMapping
</a>
</em>
</td>
<td>
<p>
(Members of <code>KubernetesMapping</code> are embedded into this type.)
</p>
<p>KubernetesMapping holds the RBAC details for the mapping</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta1.VpcCni">VpcCni
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>VpcCni specifies configuration related to the VPC CNI.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>env</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Env defines a list of environment variables to apply to the <code>aws-node</code> DaemonSet</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="controlplane.cluster.x-k8s.io/v1beta2">controlplane.cluster.x-k8s.io/v1beta2</h2>
<p>
<p>Package v1beta2 contains API Schema definitions for the controlplane v1beta2 API group</p>
</p>
Resource Types:
<ul></ul>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlane">AWSManagedControlPlane
</h3>
<p>
<p>AWSManagedControlPlane is the schema for the Amazon EKS Managed Control Plane API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">
AWSManagedControlPlaneSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>eksClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSClusterName allows you to specify the name of the EKS cluster in
AWS. If you don&rsquo;t specify a name then a default name will be created
based on the namespace and name of the managed control plane.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryCidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines the desired Kubernetes version. If no version number
is supplied then the latest version of Kubernetes that EKS supports
will be used.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role that gives EKS
permission to make API calls. If the role is pre-existing
we will treat it as unmanaged and not delete it on
deletion. If the EKSEnableIAM feature flag is true
and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the control plane role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>logging</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.ControlPlaneLoggingSpec">
ControlPlaneLoggingSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Logging specifies which EKS Cluster logs should be enabled. Entries for
each of the enabled logs will be sent to CloudWatch</p>
</td>
</tr>
<tr>
<td>
<code>encryptionConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.EncryptionConfig">
EncryptionConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptionConfig specifies the encryption configuration for the cluster</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>iamAuthenticatorConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.IAMAuthenticatorConfig">
IAMAuthenticatorConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMAuthenticatorConfig allows the specification of any additional user or role mappings
for use when generating the aws-iam-authenticator configuration. If this is nil the
default configuration is still generated for the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>endpointAccess</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.EndpointAccess">
EndpointAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Endpoints specifies access to this cluster&rsquo;s control plane endpoints</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>tokenMethod</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.EKSTokenMethod">
EKSTokenMethod
</a>
</em>
</td>
<td>
<p>TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
iam-authenticator - obtains a client token using iam-authentictor
aws-cli - obtains a client token using the AWS CLI
Defaults to iam-authenticator</p>
</td>
</tr>
<tr>
<td>
<code>associateOIDCProvider</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AssociateOIDCProvider can be enabled to automatically create an identity
provider for the controller for use with IAM roles for service accounts</p>
</td>
</tr>
<tr>
<td>
<code>addons</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2.Addon">
[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2.Addon
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Addons defines the EKS addons to enable with the EKS cluster.</p>
</td>
</tr>
<tr>
<td>
<code>oidcIdentityProviderConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.OIDCIdentityProviderConfig">
OIDCIdentityProviderConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityProviderconfig is used to specify the oidc provider config
to be attached with this eks cluster</p>
</td>
</tr>
<tr>
<td>
<code>vpcCni</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.VpcCni">
VpcCni
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>VpcCni is used to set configuration options for the VPC CNI plugin</p>
</td>
</tr>
<tr>
<td>
<code>kubeProxy</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.KubeProxy">
KubeProxy
</a>
</em>
</td>
<td>
<p>KubeProxy defines managed attributes of the kube-proxy daemonset</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneStatus">
AWSManagedControlPlaneStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlane">AWSManagedControlPlane</a>)
</p>
<p>
<p>AWSManagedControlPlaneSpec defines the desired state of an Amazon EKS Cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>eksClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSClusterName allows you to specify the name of the EKS cluster in
AWS. If you don&rsquo;t specify a name then a default name will be created
based on the namespace and name of the managed control plane.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryCidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines the desired Kubernetes version. If no version number
is supplied then the latest version of Kubernetes that EKS supports
will be used.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role that gives EKS
permission to make API calls. If the role is pre-existing
we will treat it as unmanaged and not delete it on
deletion. If the EKSEnableIAM feature flag is true
and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the control plane role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>logging</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.ControlPlaneLoggingSpec">
ControlPlaneLoggingSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Logging specifies which EKS Cluster logs should be enabled. Entries for
each of the enabled logs will be sent to CloudWatch</p>
</td>
</tr>
<tr>
<td>
<code>encryptionConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.EncryptionConfig">
EncryptionConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptionConfig specifies the encryption configuration for the cluster</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>iamAuthenticatorConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.IAMAuthenticatorConfig">
IAMAuthenticatorConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMAuthenticatorConfig allows the specification of any additional user or role mappings
for use when generating the aws-iam-authenticator configuration. If this is nil the
default configuration is still generated for the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>endpointAccess</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.EndpointAccess">
EndpointAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Endpoints specifies access to this cluster&rsquo;s control plane endpoints</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>tokenMethod</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.EKSTokenMethod">
EKSTokenMethod
</a>
</em>
</td>
<td>
<p>TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
iam-authenticator - obtains a client token using iam-authentictor
aws-cli - obtains a client token using the AWS CLI
Defaults to iam-authenticator</p>
</td>
</tr>
<tr>
<td>
<code>associateOIDCProvider</code><br/>
<em>
bool
</em>
</td>
<td>
<p>AssociateOIDCProvider can be enabled to automatically create an identity
provider for the controller for use with IAM roles for service accounts</p>
</td>
</tr>
<tr>
<td>
<code>addons</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2.Addon">
[]sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2.Addon
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Addons defines the EKS addons to enable with the EKS cluster.</p>
</td>
</tr>
<tr>
<td>
<code>oidcIdentityProviderConfig</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.OIDCIdentityProviderConfig">
OIDCIdentityProviderConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityProviderconfig is used to specify the oidc provider config
to be attached with this eks cluster</p>
</td>
</tr>
<tr>
<td>
<code>vpcCni</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.VpcCni">
VpcCni
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>VpcCni is used to set configuration options for the VPC CNI plugin</p>
</td>
</tr>
<tr>
<td>
<code>kubeProxy</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.KubeProxy">
KubeProxy
</a>
</em>
</td>
<td>
<p>KubeProxy defines managed attributes of the kube-proxy daemonset</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlane">AWSManagedControlPlane</a>)
</p>
<p>
<p>AWSManagedControlPlaneStatus defines the observed state of an Amazon EKS Cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>networkStatus</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkStatus">
NetworkStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Networks holds details about the AWS networking resources used by the control plane</p>
</td>
</tr>
<tr>
<td>
<code>failureDomains</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.FailureDomains
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureDomains specifies a list fo available availability zones that can be used</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">
Instance
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion holds details of the instance that is used as a bastion jump box</p>
</td>
</tr>
<tr>
<td>
<code>oidcProvider</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.OIDCProviderStatus">
OIDCProviderStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OIDCProvider holds the status of the identity provider for this cluster</p>
</td>
</tr>
<tr>
<td>
<code>externalManagedControlPlane</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ExternalManagedControlPlane indicates to cluster-api that the control plane
is managed by an external service such as AKS, EKS, GKE, etc.</p>
</td>
</tr>
<tr>
<td>
<code>initialized</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Initialized denotes whether or not the control plane has the
uploaded kubernetes config-map.</p>
</td>
</tr>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the AWSManagedControlPlane API Server is ready to
receive requests and that the VPC infra is ready.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ErrorMessage indicates that there is a terminal problem reconciling the
state, and will be set to a descriptive error message.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<p>Conditions specifies the cpnditions for the managed control plane</p>
</td>
</tr>
<tr>
<td>
<code>addons</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.AddonState">
[]AddonState
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Addons holds the current status of the EKS addons</p>
</td>
</tr>
<tr>
<td>
<code>identityProviderStatus</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.IdentityProviderStatus">
IdentityProviderStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityProviderStatus holds the status for
associated identity provider</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.Addon">Addon
</h3>
<p>
<p>Addon represents a EKS addon.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the addon</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version is the version of the addon to use</p>
</td>
</tr>
<tr>
<td>
<code>configuration</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration of the EKS addon</p>
</td>
</tr>
<tr>
<td>
<code>conflictResolution</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.AddonResolution">
AddonResolution
</a>
</em>
</td>
<td>
<p>ConflictResolution is used to declare what should happen if there
are parameter conflicts. Defaults to none</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountRoleArn is the ARN of an IAM role to bind to the addons service account</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AddonIssue">AddonIssue
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AddonState">AddonState</a>)
</p>
<p>
<p>AddonIssue represents an issue with an addon.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>code</code><br/>
<em>
string
</em>
</td>
<td>
<p>Code is the issue code</p>
</td>
</tr>
<tr>
<td>
<code>message</code><br/>
<em>
string
</em>
</td>
<td>
<p>Message is the textual description of the issue</p>
</td>
</tr>
<tr>
<td>
<code>resourceIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ResourceIDs is a list of resource ids for the issue</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AddonResolution">AddonResolution
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.Addon">Addon</a>)
</p>
<p>
<p>AddonResolution defines the method for resolving parameter conflicts.</p>
</p>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AddonState">AddonState
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>)
</p>
<p>
<p>AddonState represents the state of an addon.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the addon</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version is the version of the addon to use</p>
</td>
</tr>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<p>ARN is the AWS ARN of the addon</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>ServiceAccountRoleArn is the ARN of the IAM role used for the service account</p>
</td>
</tr>
<tr>
<td>
<code>createdAt</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>CreatedAt is the date and time the addon was created at</p>
</td>
</tr>
<tr>
<td>
<code>modifiedAt</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>ModifiedAt is the date and time the addon was last modified</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
string
</em>
</td>
<td>
<p>Status is the status of the addon</p>
</td>
</tr>
<tr>
<td>
<code>issues</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.AddonIssue">
[]AddonIssue
</a>
</em>
</td>
<td>
<p>Issues is a list of issue associated with the addon</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AddonStatus">AddonStatus
(<code>string</code> alias)</p></h3>
<p>
<p>AddonStatus defines the status for an addon.</p>
</p>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.ControlPlaneLoggingSpec">ControlPlaneLoggingSpec
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>ControlPlaneLoggingSpec defines what EKS control plane logs that should be enabled.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiServer</code><br/>
<em>
bool
</em>
</td>
<td>
<p>APIServer indicates if the Kubernetes API Server log (kube-apiserver) shoulkd be enabled</p>
</td>
</tr>
<tr>
<td>
<code>audit</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Audit indicates if the Kubernetes API audit log should be enabled</p>
</td>
</tr>
<tr>
<td>
<code>authenticator</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Authenticator indicates if the iam authenticator log should be enabled</p>
</td>
</tr>
<tr>
<td>
<code>controllerManager</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ControllerManager indicates if the controller manager (kube-controller-manager) log should be enabled</p>
</td>
</tr>
<tr>
<td>
<code>scheduler</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Scheduler indicates if the Kubernetes scheduler (kube-scheduler) log should be enabled</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.EKSTokenMethod">EKSTokenMethod
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>EKSTokenMethod defines the method for obtaining a client token to use when connecting to EKS.</p>
</p>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.EncryptionConfig">EncryptionConfig
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>EncryptionConfig specifies the encryption configuration for the EKS clsuter.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>provider</code><br/>
<em>
string
</em>
</td>
<td>
<p>Provider specifies the ARN or alias of the CMK (in AWS KMS)</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
[]*string
</em>
</td>
<td>
<p>Resources specifies the resources to be encrypted</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.EndpointAccess">EndpointAccess
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>EndpointAccess specifies how control plane endpoints are accessible.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>public</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Public controls whether control plane endpoints are publicly accessible</p>
</td>
</tr>
<tr>
<td>
<code>publicCIDRs</code><br/>
<em>
[]*string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicCIDRs specifies which blocks can access the public endpoint</p>
</td>
</tr>
<tr>
<td>
<code>private</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Private points VPC-internal control plane access to the private endpoint</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.IAMAuthenticatorConfig">IAMAuthenticatorConfig
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>IAMAuthenticatorConfig represents an aws-iam-authenticator configuration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mapRoles</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.RoleMapping">
[]RoleMapping
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleMappings is a list of role mappings</p>
</td>
</tr>
<tr>
<td>
<code>mapUsers</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.UserMapping">
[]UserMapping
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>UserMappings is a list of user mappings</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.IdentityProviderStatus">IdentityProviderStatus
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>)
</p>
<p>
<p>IdentityProviderStatus holds the status for associated identity provider.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<p>ARN holds the ARN of associated identity provider</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
string
</em>
</td>
<td>
<p>Status holds current status of associated identity provider</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.KubeProxy">KubeProxy
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>KubeProxy specifies how the kube-proxy daemonset is managed.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable set to true indicates that kube-proxy should be disabled. With EKS clusters
kube-proxy is automatically installed into the cluster. For clusters where you want
to use kube-proxy functionality that is provided with an alternate CNI, this option
provides a way to specify that the kube-proxy daemonset should be deleted. You cannot
set this to true if you are using the Amazon kube-proxy addon.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.KubernetesMapping">KubernetesMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.RoleMapping">RoleMapping</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.UserMapping">UserMapping</a>)
</p>
<p>
<p>KubernetesMapping represents the kubernetes RBAC mapping.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>username</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserName is a kubernetes RBAC user subject</p>
</td>
</tr>
<tr>
<td>
<code>groups</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Groups is a list of kubernetes RBAC groups</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.OIDCIdentityProviderConfig">OIDCIdentityProviderConfig
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>OIDCIdentityProviderConfig represents the configuration for an OIDC identity provider.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>clientId</code><br/>
<em>
string
</em>
</td>
<td>
<p>This is also known as audience. The ID for the client application that makes
authentication requests to the OpenID identity provider.</p>
</td>
</tr>
<tr>
<td>
<code>groupsClaim</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The JWT claim that the provider uses to return your groups.</p>
</td>
</tr>
<tr>
<td>
<code>groupsPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The prefix that is prepended to group claims to prevent clashes with existing
names (such as system: groups). For example, the valueoidc: will create group
names like oidc:engineering and oidc:infra.</p>
</td>
</tr>
<tr>
<td>
<code>identityProviderConfigName</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the OIDC provider configuration.</p>
<p>IdentityProviderConfigName is a required field</p>
</td>
</tr>
<tr>
<td>
<code>issuerUrl</code><br/>
<em>
string
</em>
</td>
<td>
<p>The URL of the OpenID identity provider that allows the API server to discover
public signing keys for verifying tokens. The URL must begin with https://
and should correspond to the iss claim in the provider&rsquo;s OIDC ID tokens.
Per the OIDC standard, path components are allowed but query parameters are
not. Typically the URL consists of only a hostname, like <a href="https://server.example.org">https://server.example.org</a>
or <a href="https://example.com">https://example.com</a>. This URL should point to the level below .well-known/openid-configuration
and must be publicly accessible over the internet.</p>
</td>
</tr>
<tr>
<td>
<code>requiredClaims</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The key value pairs that describe required claims in the identity token.
If set, each claim is verified to be present in the token with a matching
value. For the maximum number of claims that you can require, see Amazon
EKS service quotas (<a href="https://docs.aws.amazon.com/eks/latest/userguide/service-quotas.html">https://docs.aws.amazon.com/eks/latest/userguide/service-quotas.html</a>)
in the Amazon EKS User Guide.</p>
</td>
</tr>
<tr>
<td>
<code>usernameClaim</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The JSON Web Token (JWT) claim to use as the username. The default is sub,
which is expected to be a unique identifier of the end user. You can choose
other claims, such as email or name, depending on the OpenID identity provider.
Claims other than email are prefixed with the issuer URL to prevent naming
clashes with other plug-ins.</p>
</td>
</tr>
<tr>
<td>
<code>usernamePrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The prefix that is prepended to username claims to prevent clashes with existing
names. If you do not provide this field, and username is a value other than
email, the prefix defaults to issuerurl#. You can use the value - to disable
all prefixing.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>tags to apply to oidc identity provider association</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.OIDCProviderStatus">OIDCProviderStatus
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>)
</p>
<p>
<p>OIDCProviderStatus holds the status of the AWS OIDC identity provider.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<p>ARN holds the ARN of the provider</p>
</td>
</tr>
<tr>
<td>
<code>trustPolicy</code><br/>
<em>
string
</em>
</td>
<td>
<p>TrustPolicy contains the boilerplate IAM trust policy to use for IRSA</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.RoleMapping">RoleMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.IAMAuthenticatorConfig">IAMAuthenticatorConfig</a>)
</p>
<p>
<p>RoleMapping represents a mapping from a IAM role to Kubernetes users and groups.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>rolearn</code><br/>
<em>
string
</em>
</td>
<td>
<p>RoleARN is the AWS ARN for the role to map</p>
</td>
</tr>
<tr>
<td>
<code>KubernetesMapping</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.KubernetesMapping">
KubernetesMapping
</a>
</em>
</td>
<td>
<p>
(Members of <code>KubernetesMapping</code> are embedded into this type.)
</p>
<p>KubernetesMapping holds the RBAC details for the mapping</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.UserMapping">UserMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.IAMAuthenticatorConfig">IAMAuthenticatorConfig</a>)
</p>
<p>
<p>UserMapping represents a mapping from an IAM user to Kubernetes users and groups.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>userarn</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserARN is the AWS ARN for the user to map</p>
</td>
</tr>
<tr>
<td>
<code>KubernetesMapping</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.KubernetesMapping">
KubernetesMapping
</a>
</em>
</td>
<td>
<p>
(Members of <code>KubernetesMapping</code> are embedded into this type.)
</p>
<p>KubernetesMapping holds the RBAC details for the mapping</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.VpcCni">VpcCni
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>VpcCni specifies configuration related to the VPC CNI.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable indicates that the Amazon VPC CNI should be disabled. With EKS clusters the
Amazon VPC CNI is automatically installed into the cluster. For clusters where you want
to use an alternate CNI this option provides a way to specify that the Amazon VPC CNI
should be deleted. You cannot set this to true if you are using the
Amazon VPC CNI addon.</p>
</td>
</tr>
<tr>
<td>
<code>env</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Env defines a list of environment variables to apply to the <code>aws-node</code> DaemonSet</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.AWSRolesRef">AWSRolesRef
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec</a>)
</p>
<p>
<p>AWSRolesRef contains references to various AWS IAM roles required for operators to make calls against the AWS API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ingressARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>The referenced role must have a trust relationship that allows it to be assumed via web identity.
<a href="https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_oidc.html">https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_oidc.html</a>.
Example:
{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Principal&rdquo;: {
&ldquo;Federated&rdquo;: &ldquo;{{ .ProviderARN }}&rdquo;
},
&ldquo;Action&rdquo;: &ldquo;sts:AssumeRoleWithWebIdentity&rdquo;,
&ldquo;Condition&rdquo;: {
&ldquo;StringEquals&rdquo;: {
&ldquo;{{ .ProviderName }}:sub&rdquo;: {{ .ServiceAccounts }}
}
}
}
]
}</p>
<p>IngressARN is an ARN value referencing a role appropriate for the Ingress Operator.</p>
<p>The following is an example of a valid policy document:</p>
<p>{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;elasticloadbalancing:DescribeLoadBalancers&rdquo;,
&ldquo;tag:GetResources&rdquo;,
&ldquo;route53:ListHostedZones&rdquo;
],
&ldquo;Resource&rdquo;: &ldquo;*&rdquo;
},
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;route53:ChangeResourceRecordSets&rdquo;
],
&ldquo;Resource&rdquo;: [
&ldquo;arn:aws:route53:::PUBLIC_ZONE_ID&rdquo;,
&ldquo;arn:aws:route53:::PRIVATE_ZONE_ID&rdquo;
]
}
]
}</p>
</td>
</tr>
<tr>
<td>
<code>imageRegistryARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageRegistryARN is an ARN value referencing a role appropriate for the Image Registry Operator.</p>
<p>The following is an example of a valid policy document:</p>
<p>{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;s3:CreateBucket&rdquo;,
&ldquo;s3:DeleteBucket&rdquo;,
&ldquo;s3:PutBucketTagging&rdquo;,
&ldquo;s3:GetBucketTagging&rdquo;,
&ldquo;s3:PutBucketPublicAccessBlock&rdquo;,
&ldquo;s3:GetBucketPublicAccessBlock&rdquo;,
&ldquo;s3:PutEncryptionConfiguration&rdquo;,
&ldquo;s3:GetEncryptionConfiguration&rdquo;,
&ldquo;s3:PutLifecycleConfiguration&rdquo;,
&ldquo;s3:GetLifecycleConfiguration&rdquo;,
&ldquo;s3:GetBucketLocation&rdquo;,
&ldquo;s3:ListBucket&rdquo;,
&ldquo;s3:GetObject&rdquo;,
&ldquo;s3:PutObject&rdquo;,
&ldquo;s3:DeleteObject&rdquo;,
&ldquo;s3:ListBucketMultipartUploads&rdquo;,
&ldquo;s3:AbortMultipartUpload&rdquo;,
&ldquo;s3:ListMultipartUploadParts&rdquo;
],
&ldquo;Resource&rdquo;: &ldquo;*&rdquo;
}
]
}</p>
</td>
</tr>
<tr>
<td>
<code>storageARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>StorageARN is an ARN value referencing a role appropriate for the Storage Operator.</p>
<p>The following is an example of a valid policy document:</p>
<p>{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;ec2:AttachVolume&rdquo;,
&ldquo;ec2:CreateSnapshot&rdquo;,
&ldquo;ec2:CreateTags&rdquo;,
&ldquo;ec2:CreateVolume&rdquo;,
&ldquo;ec2:DeleteSnapshot&rdquo;,
&ldquo;ec2:DeleteTags&rdquo;,
&ldquo;ec2:DeleteVolume&rdquo;,
&ldquo;ec2:DescribeInstances&rdquo;,
&ldquo;ec2:DescribeSnapshots&rdquo;,
&ldquo;ec2:DescribeTags&rdquo;,
&ldquo;ec2:DescribeVolumes&rdquo;,
&ldquo;ec2:DescribeVolumesModifications&rdquo;,
&ldquo;ec2:DetachVolume&rdquo;,
&ldquo;ec2:ModifyVolume&rdquo;
],
&ldquo;Resource&rdquo;: &ldquo;*&rdquo;
}
]
}</p>
</td>
</tr>
<tr>
<td>
<code>networkARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>NetworkARN is an ARN value referencing a role appropriate for the Network Operator.</p>
<p>The following is an example of a valid policy document:</p>
<p>{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;ec2:DescribeInstances&rdquo;,
&ldquo;ec2:DescribeInstanceStatus&rdquo;,
&ldquo;ec2:DescribeInstanceTypes&rdquo;,
&ldquo;ec2:UnassignPrivateIpAddresses&rdquo;,
&ldquo;ec2:AssignPrivateIpAddresses&rdquo;,
&ldquo;ec2:UnassignIpv6Addresses&rdquo;,
&ldquo;ec2:AssignIpv6Addresses&rdquo;,
&ldquo;ec2:DescribeSubnets&rdquo;,
&ldquo;ec2:DescribeNetworkInterfaces&rdquo;
],
&ldquo;Resource&rdquo;: &ldquo;*&rdquo;
}
]
}</p>
</td>
</tr>
<tr>
<td>
<code>kubeCloudControllerARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>KubeCloudControllerARN is an ARN value referencing a role appropriate for the KCM/KCC.
Source: <a href="https://cloud-provider-aws.sigs.k8s.io/prerequisites/#iam-policies">https://cloud-provider-aws.sigs.k8s.io/prerequisites/#iam-policies</a></p>
<p>The following is an example of a valid policy document:</p>
<p>{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Action&rdquo;: [
&ldquo;autoscaling:DescribeAutoScalingGroups&rdquo;,
&ldquo;autoscaling:DescribeLaunchConfigurations&rdquo;,
&ldquo;autoscaling:DescribeTags&rdquo;,
&ldquo;ec2:DescribeAvailabilityZones&rdquo;,
&ldquo;ec2:DescribeInstances&rdquo;,
&ldquo;ec2:DescribeImages&rdquo;,
&ldquo;ec2:DescribeRegions&rdquo;,
&ldquo;ec2:DescribeRouteTables&rdquo;,
&ldquo;ec2:DescribeSecurityGroups&rdquo;,
&ldquo;ec2:DescribeSubnets&rdquo;,
&ldquo;ec2:DescribeVolumes&rdquo;,
&ldquo;ec2:CreateSecurityGroup&rdquo;,
&ldquo;ec2:CreateTags&rdquo;,
&ldquo;ec2:CreateVolume&rdquo;,
&ldquo;ec2:ModifyInstanceAttribute&rdquo;,
&ldquo;ec2:ModifyVolume&rdquo;,
&ldquo;ec2:AttachVolume&rdquo;,
&ldquo;ec2:AuthorizeSecurityGroupIngress&rdquo;,
&ldquo;ec2:CreateRoute&rdquo;,
&ldquo;ec2:DeleteRoute&rdquo;,
&ldquo;ec2:DeleteSecurityGroup&rdquo;,
&ldquo;ec2:DeleteVolume&rdquo;,
&ldquo;ec2:DetachVolume&rdquo;,
&ldquo;ec2:RevokeSecurityGroupIngress&rdquo;,
&ldquo;ec2:DescribeVpcs&rdquo;,
&ldquo;elasticloadbalancing:AddTags&rdquo;,
&ldquo;elasticloadbalancing:AttachLoadBalancerToSubnets&rdquo;,
&ldquo;elasticloadbalancing:ApplySecurityGroupsToLoadBalancer&rdquo;,
&ldquo;elasticloadbalancing:CreateLoadBalancer&rdquo;,
&ldquo;elasticloadbalancing:CreateLoadBalancerPolicy&rdquo;,
&ldquo;elasticloadbalancing:CreateLoadBalancerListeners&rdquo;,
&ldquo;elasticloadbalancing:ConfigureHealthCheck&rdquo;,
&ldquo;elasticloadbalancing:DeleteLoadBalancer&rdquo;,
&ldquo;elasticloadbalancing:DeleteLoadBalancerListeners&rdquo;,
&ldquo;elasticloadbalancing:DescribeLoadBalancers&rdquo;,
&ldquo;elasticloadbalancing:DescribeLoadBalancerAttributes&rdquo;,
&ldquo;elasticloadbalancing:DetachLoadBalancerFromSubnets&rdquo;,
&ldquo;elasticloadbalancing:DeregisterInstancesFromLoadBalancer&rdquo;,
&ldquo;elasticloadbalancing:ModifyLoadBalancerAttributes&rdquo;,
&ldquo;elasticloadbalancing:RegisterInstancesWithLoadBalancer&rdquo;,
&ldquo;elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer&rdquo;,
&ldquo;elasticloadbalancing:AddTags&rdquo;,
&ldquo;elasticloadbalancing:CreateListener&rdquo;,
&ldquo;elasticloadbalancing:CreateTargetGroup&rdquo;,
&ldquo;elasticloadbalancing:DeleteListener&rdquo;,
&ldquo;elasticloadbalancing:DeleteTargetGroup&rdquo;,
&ldquo;elasticloadbalancing:DeregisterTargets&rdquo;,
&ldquo;elasticloadbalancing:DescribeListeners&rdquo;,
&ldquo;elasticloadbalancing:DescribeLoadBalancerPolicies&rdquo;,
&ldquo;elasticloadbalancing:DescribeTargetGroups&rdquo;,
&ldquo;elasticloadbalancing:DescribeTargetHealth&rdquo;,
&ldquo;elasticloadbalancing:ModifyListener&rdquo;,
&ldquo;elasticloadbalancing:ModifyTargetGroup&rdquo;,
&ldquo;elasticloadbalancing:RegisterTargets&rdquo;,
&ldquo;elasticloadbalancing:SetLoadBalancerPoliciesOfListener&rdquo;,
&ldquo;iam:CreateServiceLinkedRole&rdquo;,
&ldquo;kms:DescribeKey&rdquo;
],
&ldquo;Resource&rdquo;: [
&ldquo;*&rdquo;
],
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;
}
]
}</p>
</td>
</tr>
<tr>
<td>
<code>nodePoolManagementARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>NodePoolManagementARN is an ARN value referencing a role appropriate for the CAPI Controller.</p>
<p>The following is an example of a valid policy document:</p>
<p>{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Action&rdquo;: [
&ldquo;ec2:AssociateRouteTable&rdquo;,
&ldquo;ec2:AttachInternetGateway&rdquo;,
&ldquo;ec2:AuthorizeSecurityGroupIngress&rdquo;,
&ldquo;ec2:CreateInternetGateway&rdquo;,
&ldquo;ec2:CreateNatGateway&rdquo;,
&ldquo;ec2:CreateRoute&rdquo;,
&ldquo;ec2:CreateRouteTable&rdquo;,
&ldquo;ec2:CreateSecurityGroup&rdquo;,
&ldquo;ec2:CreateSubnet&rdquo;,
&ldquo;ec2:CreateTags&rdquo;,
&ldquo;ec2:DeleteInternetGateway&rdquo;,
&ldquo;ec2:DeleteNatGateway&rdquo;,
&ldquo;ec2:DeleteRouteTable&rdquo;,
&ldquo;ec2:DeleteSecurityGroup&rdquo;,
&ldquo;ec2:DeleteSubnet&rdquo;,
&ldquo;ec2:DeleteTags&rdquo;,
&ldquo;ec2:DescribeAccountAttributes&rdquo;,
&ldquo;ec2:DescribeAddresses&rdquo;,
&ldquo;ec2:DescribeAvailabilityZones&rdquo;,
&ldquo;ec2:DescribeImages&rdquo;,
&ldquo;ec2:DescribeInstances&rdquo;,
&ldquo;ec2:DescribeInternetGateways&rdquo;,
&ldquo;ec2:DescribeNatGateways&rdquo;,
&ldquo;ec2:DescribeNetworkInterfaces&rdquo;,
&ldquo;ec2:DescribeNetworkInterfaceAttribute&rdquo;,
&ldquo;ec2:DescribeRouteTables&rdquo;,
&ldquo;ec2:DescribeSecurityGroups&rdquo;,
&ldquo;ec2:DescribeSubnets&rdquo;,
&ldquo;ec2:DescribeVpcs&rdquo;,
&ldquo;ec2:DescribeVpcAttribute&rdquo;,
&ldquo;ec2:DescribeVolumes&rdquo;,
&ldquo;ec2:DetachInternetGateway&rdquo;,
&ldquo;ec2:DisassociateRouteTable&rdquo;,
&ldquo;ec2:DisassociateAddress&rdquo;,
&ldquo;ec2:ModifyInstanceAttribute&rdquo;,
&ldquo;ec2:ModifyNetworkInterfaceAttribute&rdquo;,
&ldquo;ec2:ModifySubnetAttribute&rdquo;,
&ldquo;ec2:RevokeSecurityGroupIngress&rdquo;,
&ldquo;ec2:RunInstances&rdquo;,
&ldquo;ec2:TerminateInstances&rdquo;,
&ldquo;tag:GetResources&rdquo;,
&ldquo;ec2:CreateLaunchTemplate&rdquo;,
&ldquo;ec2:CreateLaunchTemplateVersion&rdquo;,
&ldquo;ec2:DescribeLaunchTemplates&rdquo;,
&ldquo;ec2:DescribeLaunchTemplateVersions&rdquo;,
&ldquo;ec2:DeleteLaunchTemplate&rdquo;,
&ldquo;ec2:DeleteLaunchTemplateVersions&rdquo;
],
&ldquo;Resource&rdquo;: [
&ldquo;<em>&rdquo;
],
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;
},
{
&ldquo;Condition&rdquo;: {
&ldquo;StringLike&rdquo;: {
&ldquo;iam:AWSServiceName&rdquo;: &ldquo;elasticloadbalancing.amazonaws.com&rdquo;
}
},
&ldquo;Action&rdquo;: [
&ldquo;iam:CreateServiceLinkedRole&rdquo;
],
&ldquo;Resource&rdquo;: [
&ldquo;arn:</em>:iam::<em>:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing&rdquo;
],
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;
},
{
&ldquo;Action&rdquo;: [
&ldquo;iam:PassRole&rdquo;
],
&ldquo;Resource&rdquo;: [
&ldquo;arn:</em>:iam::<em>:role/</em>-worker-role&rdquo;
],
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;
},
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;kms:Decrypt&rdquo;,
&ldquo;kms:ReEncrypt&rdquo;,
&ldquo;kms:GenerateDataKeyWithoutPlainText&rdquo;,
&ldquo;kms:DescribeKey&rdquo;
],
&ldquo;Resource&rdquo;: &ldquo;<em>&rdquo;
},
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;kms:CreateGrant&rdquo;
],
&ldquo;Resource&rdquo;: &ldquo;</em>&rdquo;,
&ldquo;Condition&rdquo;: {
&ldquo;Bool&rdquo;: {
&ldquo;kms:GrantIsForAWSResource&rdquo;: true
}
}
}
]
}</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneOperatorARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>ControlPlaneOperatorARN  is an ARN value referencing a role appropriate for the Control Plane Operator.</p>
<p>The following is an example of a valid policy document:</p>
<p>{
&ldquo;Version&rdquo;: &ldquo;2012-10-17&rdquo;,
&ldquo;Statement&rdquo;: [
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;ec2:CreateVpcEndpoint&rdquo;,
&ldquo;ec2:DescribeVpcEndpoints&rdquo;,
&ldquo;ec2:ModifyVpcEndpoint&rdquo;,
&ldquo;ec2:DeleteVpcEndpoints&rdquo;,
&ldquo;ec2:CreateTags&rdquo;,
&ldquo;route53:ListHostedZones&rdquo;,
&ldquo;ec2:CreateSecurityGroup&rdquo;,
&ldquo;ec2:AuthorizeSecurityGroupIngress&rdquo;,
&ldquo;ec2:AuthorizeSecurityGroupEgress&rdquo;,
&ldquo;ec2:DeleteSecurityGroup&rdquo;,
&ldquo;ec2:RevokeSecurityGroupIngress&rdquo;,
&ldquo;ec2:RevokeSecurityGroupEgress&rdquo;,
&ldquo;ec2:DescribeSecurityGroups&rdquo;,
&ldquo;ec2:DescribeVpcs&rdquo;,
],
&ldquo;Resource&rdquo;: &ldquo;*&rdquo;
},
{
&ldquo;Effect&rdquo;: &ldquo;Allow&rdquo;,
&ldquo;Action&rdquo;: [
&ldquo;route53:ChangeResourceRecordSets&rdquo;,
&ldquo;route53:ListResourceRecordSets&rdquo;
],
&ldquo;Resource&rdquo;: &ldquo;arn:aws:route53:::%s&rdquo;
}
]
}</p>
</td>
</tr>
<tr>
<td>
<code>kmsProviderARN</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.DefaultMachinePoolSpec">DefaultMachinePoolSpec
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec</a>)
</p>
<p>
<p>DefaultMachinePoolSpec defines the configuration for the required worker nodes provisioned as part of the cluster creation.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The instance type to use, for example <code>r5.xlarge</code>. Instance type ref; <a href="https://aws.amazon.com/ec2/instance-types/">https://aws.amazon.com/ec2/instance-types/</a></p>
</td>
</tr>
<tr>
<td>
<code>autoscaling</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolAutoScaling">
RosaMachinePoolAutoScaling
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Autoscaling specifies auto scaling behaviour for the default MachinePool. Autoscaling min/max value
must be equal or multiple of the availability zones count.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider">ExternalAuthProvider
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec</a>)
</p>
<p>
<p>ExternalAuthProvider is an external OIDC identity provider that can issue tokens for this cluster</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the OIDC provider</p>
</td>
</tr>
<tr>
<td>
<code>issuer</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenIssuer">
TokenIssuer
</a>
</em>
</td>
<td>
<p>Issuer describes attributes of the OIDC token issuer</p>
</td>
</tr>
<tr>
<td>
<code>oidcClients</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.OIDCClientConfig">
[]OIDCClientConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>OIDCClients contains configuration for the platform&rsquo;s clients that
need to request tokens from the issuer</p>
</td>
</tr>
<tr>
<td>
<code>claimMappings</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenClaimMappings">
TokenClaimMappings
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ClaimMappings describes rules on how to transform information from an
ID token into a cluster identity</p>
</td>
</tr>
<tr>
<td>
<code>claimValidationRules</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenClaimValidationRule">
[]TokenClaimValidationRule
</a>
</em>
</td>
<td>
<p>ClaimValidationRules are rules that are applied to validate token claims to authenticate users.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.NetworkSpec">NetworkSpec
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec</a>)
</p>
<p>
<p>NetworkSpec for ROSA-HCP.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>machineCIDR</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IP addresses block used by OpenShift while installing the cluster, for example &ldquo;10.0.0.0/16&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>podCIDR</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IP address block from which to assign pod IP addresses, for example <code>10.128.0.0/14</code>.</p>
</td>
</tr>
<tr>
<td>
<code>serviceCIDR</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IP address block from which to assign service IP addresses, for example <code>172.30.0.0/16</code>.</p>
</td>
</tr>
<tr>
<td>
<code>hostPrefix</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>Network host prefix which is defaulted to <code>23</code> if not specified.</p>
</td>
</tr>
<tr>
<td>
<code>networkType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The CNI network type default is OVNKubernetes.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.OIDCClientConfig">OIDCClientConfig
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider">ExternalAuthProvider</a>)
</p>
<p>
<p>OIDCClientConfig contains configuration for the platform&rsquo;s client that
need to request tokens from the issuer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>componentName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ComponentName is the name of the component that is supposed to consume this
client configuration</p>
</td>
</tr>
<tr>
<td>
<code>componentNamespace</code><br/>
<em>
string
</em>
</td>
<td>
<p>ComponentNamespace is the namespace of the component that is supposed to consume this
client configuration</p>
</td>
</tr>
<tr>
<td>
<code>clientID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ClientID is the identifier of the OIDC client from the OIDC provider</p>
</td>
</tr>
<tr>
<td>
<code>clientSecret</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#secretreference-v1-core">
Kubernetes core/v1.SecretReference
</a>
</em>
</td>
<td>
<p>ClientSecret refers to a secret that
contains the client secret in the <code>clientSecret</code> key of the <code>.data</code> field</p>
</td>
</tr>
<tr>
<td>
<code>extraScopes</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ExtraScopes is an optional set of scopes to request tokens with.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.PrefixedClaimMapping">PrefixedClaimMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenClaimMappings">TokenClaimMappings</a>)
</p>
<p>
<p>PrefixedClaimMapping defines claims with a prefix.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>claim</code><br/>
<em>
string
</em>
</td>
<td>
<p>Claim is a JWT token claim to be used in the mapping</p>
</td>
</tr>
<tr>
<td>
<code>prefix</code><br/>
<em>
string
</em>
</td>
<td>
<p>Prefix is a string to prefix the value from the token in the result of the
claim mapping.</p>
<p>By default, no prefixing occurs.</p>
<p>Example: if <code>prefix</code> is set to &ldquo;myoidc:&rdquo;&rdquo; and the <code>claim</code> in JWT contains
an array of strings &ldquo;a&rdquo;, &ldquo;b&rdquo; and  &ldquo;c&rdquo;, the mapping will result in an
array of string &ldquo;myoidc:a&rdquo;, &ldquo;myoidc:b&rdquo; and &ldquo;myoidc:c&rdquo;.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.ROSAControlPlane">ROSAControlPlane
</h3>
<p>
<p>ROSAControlPlane is the Schema for the ROSAControlPlanes API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">
RosaControlPlaneSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>rosaClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>Cluster name must be valid DNS-1035 label, so it must consist of lower case alphanumeric
characters or &lsquo;-&rsquo;, start with an alphabetic character, end with an alphanumeric character
and have a max length of 54 characters.</p>
</td>
</tr>
<tr>
<td>
<code>domainPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DomainPrefix is an optional prefix added to the cluster&rsquo;s domain name. It will be used
when generating a sub-domain for the cluster on openshiftapps domain. It must be valid DNS-1035 label
consisting of lower case alphanumeric characters or &lsquo;-&rsquo;, start with an alphabetic character
end with an alphanumeric character and have a max length of 15 characters.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The Subnet IDs to use when installing the cluster.
SubnetIDs should come in pairs; two per availability zone, one private and one public.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones describe AWS AvailabilityZones of the worker nodes.
should match the AvailabilityZones of the provided Subnets.
a machinepool will be created for each availabilityZone.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>OpenShift semantic version, for example &ldquo;4.14.5&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>rolesRef</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSRolesRef">
AWSRolesRef
</a>
</em>
</td>
<td>
<p>AWS IAM roles used to perform credential requests by the openshift operators.</p>
</td>
</tr>
<tr>
<td>
<code>oidcID</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the internal OpenID Connect Provider.</p>
</td>
</tr>
<tr>
<td>
<code>enableExternalAuthProviders</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>EnableExternalAuthProviders enables external authentication configuration for the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>externalAuthProviders</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider">
[]ExternalAuthProvider
</a>
</em>
</td>
<td>
<p>ExternalAuthProviders are external OIDC identity providers that can issue tokens for this cluster.
Can only be set if &ldquo;enableExternalAuthProviders&rdquo; is set to &ldquo;True&rdquo;.</p>
<p>At most one provider can be configured.</p>
</td>
</tr>
<tr>
<td>
<code>installerRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstallerRoleARN is an AWS IAM role that OpenShift Cluster Manager will assume to create the cluster..</p>
</td>
</tr>
<tr>
<td>
<code>supportRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>SupportRoleARN is an AWS IAM role used by Red Hat SREs to enable
access to the cluster account in order to provide support.</p>
</td>
</tr>
<tr>
<td>
<code>workerRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>WorkerRoleARN is an AWS IAM role that will be attached to worker instances.</p>
</td>
</tr>
<tr>
<td>
<code>billingAccount</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>BillingAccount is an optional AWS account to use for billing the subscription fees for ROSA clusters.
The cost of running each ROSA cluster will be billed to the infrastructure account in which the cluster
is running.</p>
</td>
</tr>
<tr>
<td>
<code>defaultMachinePoolSpec</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.DefaultMachinePoolSpec">
DefaultMachinePoolSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>DefaultMachinePoolSpec defines the configuration for the default machinepool(s) provisioned as part of the cluster creation.
One MachinePool will be created with this configuration per AvailabilityZone. Those default machinepools are required for openshift cluster operators
to work properly.
As these machinepool not created using ROSAMachinePool CR, they will not be visible/managed by ROSA CAPI provider.
<code>rosa list machinepools -c &lt;rosaClusterName&gt;</code> can be used to view those machinepools.</p>
<p>This field will be removed in the future once the current limitation is resolved.</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Network config for the ROSA HCP cluster.</p>
</td>
</tr>
<tr>
<td>
<code>endpointAccess</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaEndpointAccessType">
RosaEndpointAccessType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EndpointAccess specifies the publishing scope of cluster endpoints. The
default is Public.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags are user-defined tags to be added on the AWS resources associated with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>etcdEncryptionKMSARN</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EtcdEncryptionKMSARN is the ARN of the KMS key used to encrypt etcd. The key itself needs to be
created out-of-band by the user and tagged with <code>red-hat:true</code>.</p>
</td>
</tr>
<tr>
<td>
<code>auditLogRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AuditLogRoleARN defines the role that is used to forward audit logs to AWS CloudWatch.
If not set, audit log forwarding is disabled.</p>
</td>
</tr>
<tr>
<td>
<code>provisionShardID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProvisionShardID defines the shard where rosa control plane components will be hosted.</p>
</td>
</tr>
<tr>
<td>
<code>credentialsSecretRef</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CredentialsSecretRef references a secret with necessary credentials to connect to the OCM API.
The secret should contain the following data keys:
- ocmToken: eyJhbGciOiJIUzI1NiIsI&hellip;.
- ocmApiUrl: Optional, defaults to &lsquo;<a href="https://api.openshift.com'">https://api.openshift.com&rsquo;</a></p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneStatus">
RosaControlPlaneStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.ROSAControlPlane">ROSAControlPlane</a>)
</p>
<p>
<p>RosaControlPlaneSpec defines the desired state of ROSAControlPlane.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>rosaClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>Cluster name must be valid DNS-1035 label, so it must consist of lower case alphanumeric
characters or &lsquo;-&rsquo;, start with an alphabetic character, end with an alphanumeric character
and have a max length of 54 characters.</p>
</td>
</tr>
<tr>
<td>
<code>domainPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>DomainPrefix is an optional prefix added to the cluster&rsquo;s domain name. It will be used
when generating a sub-domain for the cluster on openshiftapps domain. It must be valid DNS-1035 label
consisting of lower case alphanumeric characters or &lsquo;-&rsquo;, start with an alphabetic character
end with an alphanumeric character and have a max length of 15 characters.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The Subnet IDs to use when installing the cluster.
SubnetIDs should come in pairs; two per availability zone, one private and one public.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones describe AWS AvailabilityZones of the worker nodes.
should match the AvailabilityZones of the provided Subnets.
a machinepool will be created for each availabilityZone.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>OpenShift semantic version, for example &ldquo;4.14.5&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>rolesRef</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSRolesRef">
AWSRolesRef
</a>
</em>
</td>
<td>
<p>AWS IAM roles used to perform credential requests by the openshift operators.</p>
</td>
</tr>
<tr>
<td>
<code>oidcID</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the internal OpenID Connect Provider.</p>
</td>
</tr>
<tr>
<td>
<code>enableExternalAuthProviders</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>EnableExternalAuthProviders enables external authentication configuration for the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>externalAuthProviders</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider">
[]ExternalAuthProvider
</a>
</em>
</td>
<td>
<p>ExternalAuthProviders are external OIDC identity providers that can issue tokens for this cluster.
Can only be set if &ldquo;enableExternalAuthProviders&rdquo; is set to &ldquo;True&rdquo;.</p>
<p>At most one provider can be configured.</p>
</td>
</tr>
<tr>
<td>
<code>installerRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstallerRoleARN is an AWS IAM role that OpenShift Cluster Manager will assume to create the cluster..</p>
</td>
</tr>
<tr>
<td>
<code>supportRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>SupportRoleARN is an AWS IAM role used by Red Hat SREs to enable
access to the cluster account in order to provide support.</p>
</td>
</tr>
<tr>
<td>
<code>workerRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>WorkerRoleARN is an AWS IAM role that will be attached to worker instances.</p>
</td>
</tr>
<tr>
<td>
<code>billingAccount</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>BillingAccount is an optional AWS account to use for billing the subscription fees for ROSA clusters.
The cost of running each ROSA cluster will be billed to the infrastructure account in which the cluster
is running.</p>
</td>
</tr>
<tr>
<td>
<code>defaultMachinePoolSpec</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.DefaultMachinePoolSpec">
DefaultMachinePoolSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>DefaultMachinePoolSpec defines the configuration for the default machinepool(s) provisioned as part of the cluster creation.
One MachinePool will be created with this configuration per AvailabilityZone. Those default machinepools are required for openshift cluster operators
to work properly.
As these machinepool not created using ROSAMachinePool CR, they will not be visible/managed by ROSA CAPI provider.
<code>rosa list machinepools -c &lt;rosaClusterName&gt;</code> can be used to view those machinepools.</p>
<p>This field will be removed in the future once the current limitation is resolved.</p>
</td>
</tr>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Network config for the ROSA HCP cluster.</p>
</td>
</tr>
<tr>
<td>
<code>endpointAccess</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaEndpointAccessType">
RosaEndpointAccessType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EndpointAccess specifies the publishing scope of cluster endpoints. The
default is Public.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags are user-defined tags to be added on the AWS resources associated with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>etcdEncryptionKMSARN</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EtcdEncryptionKMSARN is the ARN of the KMS key used to encrypt etcd. The key itself needs to be
created out-of-band by the user and tagged with <code>red-hat:true</code>.</p>
</td>
</tr>
<tr>
<td>
<code>auditLogRoleARN</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AuditLogRoleARN defines the role that is used to forward audit logs to AWS CloudWatch.
If not set, audit log forwarding is disabled.</p>
</td>
</tr>
<tr>
<td>
<code>provisionShardID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProvisionShardID defines the shard where rosa control plane components will be hosted.</p>
</td>
</tr>
<tr>
<td>
<code>credentialsSecretRef</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CredentialsSecretRef references a secret with necessary credentials to connect to the OCM API.
The secret should contain the following data keys:
- ocmToken: eyJhbGciOiJIUzI1NiIsI&hellip;.
- ocmApiUrl: Optional, defaults to &lsquo;<a href="https://api.openshift.com'">https://api.openshift.com&rsquo;</a></p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneStatus">RosaControlPlaneStatus
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.ROSAControlPlane">ROSAControlPlane</a>)
</p>
<p>
<p>RosaControlPlaneStatus defines the observed state of ROSAControlPlane.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>externalManagedControlPlane</code><br/>
<em>
bool
</em>
</td>
<td>
<p>ExternalManagedControlPlane indicates to cluster-api that the control plane
is managed by an external service such as AKS, EKS, GKE, etc.</p>
</td>
</tr>
<tr>
<td>
<code>initialized</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Initialized denotes whether or not the control plane has the
uploaded kubernetes config-map.</p>
</td>
</tr>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the ROSAControlPlane API Server is ready to receive requests.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the state and will be set to a descriptive error message.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the spec or the configuration of
the controller, and that manual intervention is required.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<p>Conditions specifies the conditions for the managed control plane</p>
</td>
</tr>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is the cluster ID given by ROSA.</p>
</td>
</tr>
<tr>
<td>
<code>consoleURL</code><br/>
<em>
string
</em>
</td>
<td>
<p>ConsoleURL is the url for the openshift console.</p>
</td>
</tr>
<tr>
<td>
<code>oidcEndpointURL</code><br/>
<em>
string
</em>
</td>
<td>
<p>OIDCEndpointURL is the endpoint url for the managed OIDC provider.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.RosaEndpointAccessType">RosaEndpointAccessType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec</a>)
</p>
<p>
<p>RosaEndpointAccessType specifies the publishing scope of cluster endpoints.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Private&#34;</p></td>
<td><p>Private endpoint access allows only private API server access and private
node communication with the control plane.</p>
</td>
</tr><tr><td><p>&#34;Public&#34;</p></td>
<td><p>Public endpoint access allows public API server access and
private node communication with the control plane.</p>
</td>
</tr></tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.TokenAudience">TokenAudience
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenIssuer">TokenIssuer</a>)
</p>
<p>
<pre><code>TokenAudience is the audience that the token was issued for.
</code></pre>
</p>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.TokenClaimMappings">TokenClaimMappings
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider">ExternalAuthProvider</a>)
</p>
<p>
<p>TokenClaimMappings describes rules on how to transform information from an
ID token into a cluster identity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>username</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.UsernameClaimMapping">
UsernameClaimMapping
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Username is a name of the claim that should be used to construct
usernames for the cluster identity.</p>
<p>Default value: &ldquo;sub&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>groups</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.PrefixedClaimMapping">
PrefixedClaimMapping
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Groups is a name of the claim that should be used to construct
groups for the cluster identity.
The referenced claim must use array of strings values.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.TokenClaimValidationRule">TokenClaimValidationRule
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider">ExternalAuthProvider</a>)
</p>
<p>
<p>TokenClaimValidationRule validates token claims to authenticate users.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenValidationRuleType">
TokenValidationRuleType
</a>
</em>
</td>
<td>
<p>Type sets the type of the validation rule</p>
</td>
</tr>
<tr>
<td>
<code>requiredClaim</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenRequiredClaim">
TokenRequiredClaim
</a>
</em>
</td>
<td>
<p>RequiredClaim allows configuring a required claim name and its expected value</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.TokenIssuer">TokenIssuer
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.ExternalAuthProvider">ExternalAuthProvider</a>)
</p>
<p>
<p>TokenIssuer describes attributes of the OIDC token issuer</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>issuerURL</code><br/>
<em>
string
</em>
</td>
<td>
<p>URL is the serving URL of the token issuer.
Must use the https:// scheme.</p>
</td>
</tr>
<tr>
<td>
<code>audiences</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenAudience">
[]TokenAudience
</a>
</em>
</td>
<td>
<p>Audiences is an array of audiences that the token was issued for.
Valid tokens must include at least one of these values in their
&ldquo;aud&rdquo; claim.
Must be set to exactly one value.</p>
</td>
</tr>
<tr>
<td>
<code>issuerCertificateAuthority</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>CertificateAuthority is a reference to a config map in the
configuration namespace. The .data of the configMap must contain
the &ldquo;ca-bundle.crt&rdquo; key.
If unset, system trust is used instead.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.TokenRequiredClaim">TokenRequiredClaim
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenClaimValidationRule">TokenClaimValidationRule</a>)
</p>
<p>
<p>TokenRequiredClaim allows configuring a required claim name and its expected value.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>claim</code><br/>
<em>
string
</em>
</td>
<td>
<p>Claim is a name of a required claim. Only claims with string values are
supported.</p>
</td>
</tr>
<tr>
<td>
<code>requiredValue</code><br/>
<em>
string
</em>
</td>
<td>
<p>RequiredValue is the required value for the claim.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.TokenValidationRuleType">TokenValidationRuleType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenClaimValidationRule">TokenClaimValidationRule</a>)
</p>
<p>
<p>TokenValidationRuleType defines the type of the validation rule.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;RequiredClaim&#34;</p></td>
<td><p>TokenValidationRuleTypeRequiredClaim defines the type for RequiredClaim.</p>
</td>
</tr></tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.UsernameClaimMapping">UsernameClaimMapping
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.TokenClaimMappings">TokenClaimMappings</a>)
</p>
<p>
<p>UsernameClaimMapping defines the claim that should be used to construct usernames for the cluster identity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>claim</code><br/>
<em>
string
</em>
</td>
<td>
<p>Claim is a JWT token claim to be used in the mapping</p>
</td>
</tr>
<tr>
<td>
<code>prefixPolicy</code><br/>
<em>
<a href="#controlplane.cluster.x-k8s.io/v1beta2.UsernamePrefixPolicy">
UsernamePrefixPolicy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrefixPolicy specifies how a prefix should apply.</p>
<p>By default, claims other than <code>email</code> will be prefixed with the issuer URL to
prevent naming clashes with other plugins.</p>
<p>Set to &ldquo;NoPrefix&rdquo; to disable prefixing.</p>
<p>Example:
(1) <code>prefix</code> is set to &ldquo;myoidc:&rdquo; and <code>claim</code> is set to &ldquo;username&rdquo;.
If the JWT claim <code>username</code> contains value <code>userA</code>, the resulting
mapped value will be &ldquo;myoidc:userA&rdquo;.
(2) <code>prefix</code> is set to &ldquo;myoidc:&rdquo; and <code>claim</code> is set to &ldquo;email&rdquo;. If the
JWT <code>email</code> claim contains value &ldquo;userA@myoidc.tld&rdquo;, the resulting
mapped value will be &ldquo;myoidc:userA@myoidc.tld&rdquo;.
(3) <code>prefix</code> is unset, <code>issuerURL</code> is set to <code>https://myoidc.tld</code>,
the JWT claims include &ldquo;username&rdquo;:&ldquo;userA&rdquo; and &ldquo;email&rdquo;:&ldquo;userA@myoidc.tld&rdquo;,
and <code>claim</code> is set to:
(a) &ldquo;username&rdquo;: the mapped value will be &ldquo;<a href="https://myoidc.tld#userA&quot;">https://myoidc.tld#userA&rdquo;</a>
(b) &ldquo;email&rdquo;: the mapped value will be &ldquo;userA@myoidc.tld&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>prefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Prefix is prepended to claim to prevent clashes with existing names.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="controlplane.cluster.x-k8s.io/v1beta2.UsernamePrefixPolicy">UsernamePrefixPolicy
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.UsernameClaimMapping">UsernameClaimMapping</a>)
</p>
<p>
<p>UsernamePrefixPolicy specifies how a prefix should apply.</p>
</p>
<hr/>
<h2 id="infrastructure.cluster.x-k8s.io/v1beta1">infrastructure.cluster.x-k8s.io/v1beta1</h2>
<p>
<p>Package v1beta1 contains the v1beta1 API implementation.</p>
</p>
Resource Types:
<ul></ul>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AMIReference">AMIReference
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec</a>)
</p>
<p>
<p>AMIReference is a reference to a specific AWS resource by ID, ARN, or filters.
Only one of ID, ARN or Filters may be specified. Specifying more than one will result in
a validation error.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ID of resource</p>
</td>
</tr>
<tr>
<td>
<code>eksLookupType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.EKSAMILookupType">
EKSAMILookupType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSOptimizedLookupType If specified, will look up an EKS Optimized image in SSM Parameter store</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSCluster">AWSCluster
</h3>
<p>
<p>AWSCluster is the schema for Amazon EC2 based Kubernetes Cluster API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">
AWSClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.S3Bucket">
S3Bucket
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Bucket contains options to configure a supporting S3 bucket for this
cluster - currently used for nodes requiring Ignition
(<a href="https://coreos.github.io/ignition/">https://coreos.github.io/ignition/</a>) for bootstrapping (requires
BootstrapFormatIgnition feature flag to be enabled).</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStatus">
AWSClusterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterControllerIdentity">AWSClusterControllerIdentity
</h3>
<p>
<p>AWSClusterControllerIdentity is the Schema for the awsclustercontrolleridentities API
It is used to grant access to use Cluster API Provider AWS Controller credentials.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterControllerIdentitySpec">
AWSClusterControllerIdentitySpec
</a>
</em>
</td>
<td>
<p>Spec for this AWSClusterControllerIdentity.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterControllerIdentitySpec">AWSClusterControllerIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterControllerIdentity">AWSClusterControllerIdentity</a>)
</p>
<p>
<p>AWSClusterControllerIdentitySpec defines the specifications for AWSClusterControllerIdentity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">AWSClusterIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterControllerIdentitySpec">AWSClusterControllerIdentitySpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStaticIdentitySpec">AWSClusterStaticIdentitySpec</a>)
</p>
<p>
<p>AWSClusterIdentitySpec defines the Spec struct for AWSClusterIdentity types.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>allowedNamespaces</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AllowedNamespaces">
AllowedNamespaces
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AllowedNamespaces is used to identify which namespaces are allowed to use the identity from.
Namespaces can be selected either using an array of namespaces or with label selector.
An empty allowedNamespaces object indicates that AWSClusters can use this identity from any namespace.
If this object is nil, no namespaces will be allowed (default behaviour, if this field is not provided)
A namespace should be either in the NamespaceList or match with Selector to use the identity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterRoleIdentity">AWSClusterRoleIdentity
</h3>
<p>
<p>AWSClusterRoleIdentity is the Schema for the awsclusterroleidentities API
It is used to assume a role using the provided sourceRef.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterRoleIdentitySpec">
AWSClusterRoleIdentitySpec
</a>
</em>
</td>
<td>
<p>Spec for this AWSClusterRoleIdentity.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>AWSRoleSpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSRoleSpec">
AWSRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>externalID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>A unique identifier that might be required when you assume a role in another account.
If the administrator of the account to which the role belongs provided you with an
external ID, then provide that value in the ExternalId parameter. This value can be
any string, such as a passphrase or account number. A cross-account role is usually
set up to trust everyone in an account. Therefore, the administrator of the trusting
account might send an external ID to the administrator of the trusted account. That
way, only someone with the ID can assume the role, rather than everyone in the
account. For more information about the external ID, see How to Use an External ID
When Granting Access to Your AWS Resources to a Third Party in the IAM User Guide.</p>
</td>
</tr>
<tr>
<td>
<code>sourceIdentityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>SourceIdentityRef is a reference to another identity which will be chained to do
role assumption. All identity types are accepted.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterRoleIdentity">AWSClusterRoleIdentity</a>)
</p>
<p>
<p>AWSClusterRoleIdentitySpec defines the specifications for AWSClusterRoleIdentity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>AWSRoleSpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSRoleSpec">
AWSRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>externalID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>A unique identifier that might be required when you assume a role in another account.
If the administrator of the account to which the role belongs provided you with an
external ID, then provide that value in the ExternalId parameter. This value can be
any string, such as a passphrase or account number. A cross-account role is usually
set up to trust everyone in an account. Therefore, the administrator of the trusting
account might send an external ID to the administrator of the trusted account. That
way, only someone with the ID can assume the role, rather than everyone in the
account. For more information about the external ID, see How to Use an External ID
When Granting Access to Your AWS Resources to a Third Party in the IAM User Guide.</p>
</td>
</tr>
<tr>
<td>
<code>sourceIdentityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>SourceIdentityRef is a reference to another identity which will be chained to do
role assumption. All identity types are accepted.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">AWSClusterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSCluster">AWSCluster</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplateResource">AWSClusterTemplateResource</a>)
</p>
<p>
<p>AWSClusterSpec defines the desired state of an EC2-based Kubernetes cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.S3Bucket">
S3Bucket
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Bucket contains options to configure a supporting S3 bucket for this
cluster - currently used for nodes requiring Ignition
(<a href="https://coreos.github.io/ignition/">https://coreos.github.io/ignition/</a>) for bootstrapping (requires
BootstrapFormatIgnition feature flag to be enabled).</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStaticIdentity">AWSClusterStaticIdentity
</h3>
<p>
<p>AWSClusterStaticIdentity is the Schema for the awsclusterstaticidentities API
It represents a reference to an AWS access key ID and secret access key, stored in a secret.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStaticIdentitySpec">
AWSClusterStaticIdentitySpec
</a>
</em>
</td>
<td>
<p>Spec for this AWSClusterStaticIdentity</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
string
</em>
</td>
<td>
<p>Reference to a secret containing the credentials. The secret should
contain the following data keys:
AccessKeyID: AKIAIOSFODNN7EXAMPLE
SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
SessionToken: Optional</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStaticIdentitySpec">AWSClusterStaticIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStaticIdentity">AWSClusterStaticIdentity</a>)
</p>
<p>
<p>AWSClusterStaticIdentitySpec defines the specifications for AWSClusterStaticIdentity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
string
</em>
</td>
<td>
<p>Reference to a secret containing the credentials. The secret should
contain the following data keys:
AccessKeyID: AKIAIOSFODNN7EXAMPLE
SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
SessionToken: Optional</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStatus">AWSClusterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSCluster">AWSCluster</a>)
</p>
<p>
<p>AWSClusterStatus defines the observed state of AWSCluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>networkStatus</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkStatus">
NetworkStatus
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>failureDomains</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.FailureDomains
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Instance">
Instance
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplate">AWSClusterTemplate
</h3>
<p>
<p>AWSClusterTemplate is the schema for Amazon EC2 based Kubernetes Cluster Templates.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplateSpec">
AWSClusterTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplateResource">
AWSClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplateResource">AWSClusterTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplateSpec">AWSClusterTemplateSpec</a>)
</p>
<p>
<p>AWSClusterTemplateResource defines the desired state of AWSClusterTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.
More info: <a href="https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata">https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata</a></p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">
AWSClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.S3Bucket">
S3Bucket
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Bucket contains options to configure a supporting S3 bucket for this
cluster - currently used for nodes requiring Ignition
(<a href="https://coreos.github.io/ignition/">https://coreos.github.io/ignition/</a>) for bootstrapping (requires
BootstrapFormatIgnition feature flag to be enabled).</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplateSpec">AWSClusterTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplate">AWSClusterTemplate</a>)
</p>
<p>
<p>AWSClusterTemplateSpec defines the desired state of AWSClusterTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterTemplateResource">
AWSClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityKind">AWSIdentityKind
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityReference">AWSIdentityReference</a>)
</p>
<p>
<p>AWSIdentityKind defines allowed AWS identity types.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityReference">AWSIdentityReference
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">AWSClusterSpec</a>)
</p>
<p>
<p>AWSIdentityReference specifies a identity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the identity.</p>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSIdentityKind">
AWSIdentityKind
</a>
</em>
</td>
<td>
<p>Kind of the identity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSLoadBalancerSpec">AWSLoadBalancerSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">AWSClusterSpec</a>)
</p>
<p>
<p>AWSLoadBalancerSpec defines the desired state of an AWS load balancer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Name sets the name of the classic ELB load balancer. As per AWS, the name must be unique
within your set of load balancers for the region, must have a maximum of 32 characters, must
contain only alphanumeric characters or hyphens, and cannot begin or end with a hyphen. Once
set, the value cannot be changed.</p>
</td>
</tr>
<tr>
<td>
<code>scheme</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBScheme">
ClassicELBScheme
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Scheme sets the scheme of the load balancer (defaults to internet-facing)</p>
</td>
</tr>
<tr>
<td>
<code>crossZoneLoadBalancing</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>CrossZoneLoadBalancing enables the classic ELB cross availability zone balancing.</p>
<p>With cross-zone load balancing, each load balancer node for your Classic Load Balancer
distributes requests evenly across the registered instances in all enabled Availability Zones.
If cross-zone load balancing is disabled, each load balancer node distributes requests evenly across
the registered instances in its Availability Zone only.</p>
<p>Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets sets the subnets that should be applied to the control plane load balancer (defaults to discovered subnets for managed VPCs or an empty set for unmanaged VPCs)</p>
</td>
</tr>
<tr>
<td>
<code>healthCheckProtocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBProtocol">
ClassicELBProtocol
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>HealthCheckProtocol sets the protocol type for classic ELB health check target
default value is ClassicELBProtocolSSL</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups sets the security groups used by the load balancer. Expected to be security group IDs
This is optional - if not provided new security groups will be created for the load balancer</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachine">AWSMachine
</h3>
<p>
<p>AWSMachine is the schema for Amazon EC2 machines.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">
AWSMachineSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceID is the EC2 instance ID for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
AWSMachine&rsquo;s value takes precedence.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMInstanceProfile is a name of an IAM instance profile to assign to the instance</p>
</td>
</tr>
<tr>
<td>
<code>publicIP</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicIP specifies whether the instance should get a public IP.
Precedence for this setting is as follows:
1. This field if set
2. Cluster/flavor setting
3. Subnet default</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instance. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator. It is possible to specify either IDs of Filters. Using Filters
will cause additional requests to AWS API and if tags change the attached security groups might change too.</p>
</td>
</tr>
<tr>
<td>
<code>failureDomain</code><br/>
<em>
string
</em>
</td>
<td>
<p>FailureDomain is the failure domain unique identifier this Machine should be attached to, as defined in Cluster API.
For this infrastructure provider, the ID is equivalent to an AWS Availability Zone.
If multiple subnets are matched for the availability zone, the first one returned is picked.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">
AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnet is a reference to the subnet to use for this instance. If not specified,
the cluster subnet will be used.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces is a list of ENIs to associate with the instance.
A maximum of 2 may be specified.</p>
</td>
</tr>
<tr>
<td>
<code>uncompressedUserData</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UncompressedUserData specify whether the user data is gzip-compressed before it is sent to ec2 instance.
cloud-init has built-in support for gzip-compressed user data
user data stored in aws secret manager is always gzip-compressed.</p>
</td>
</tr>
<tr>
<td>
<code>cloudInit</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.CloudInit">
CloudInit
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Ignition">
Ignition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ignition defined options related to the bootstrapping systems where Ignition is used.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineStatus">
AWSMachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineProviderConditionType">AWSMachineProviderConditionType
(<code>string</code> alias)</p></h3>
<p>
<p>AWSMachineProviderConditionType is a valid value for AWSMachineProviderCondition.Type.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachine">AWSMachine</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateResource">AWSMachineTemplateResource</a>)
</p>
<p>
<p>AWSMachineSpec defines the desired state of an Amazon EC2 instance.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceID is the EC2 instance ID for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
AWSMachine&rsquo;s value takes precedence.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMInstanceProfile is a name of an IAM instance profile to assign to the instance</p>
</td>
</tr>
<tr>
<td>
<code>publicIP</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicIP specifies whether the instance should get a public IP.
Precedence for this setting is as follows:
1. This field if set
2. Cluster/flavor setting
3. Subnet default</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instance. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator. It is possible to specify either IDs of Filters. Using Filters
will cause additional requests to AWS API and if tags change the attached security groups might change too.</p>
</td>
</tr>
<tr>
<td>
<code>failureDomain</code><br/>
<em>
string
</em>
</td>
<td>
<p>FailureDomain is the failure domain unique identifier this Machine should be attached to, as defined in Cluster API.
For this infrastructure provider, the ID is equivalent to an AWS Availability Zone.
If multiple subnets are matched for the availability zone, the first one returned is picked.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">
AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnet is a reference to the subnet to use for this instance. If not specified,
the cluster subnet will be used.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces is a list of ENIs to associate with the instance.
A maximum of 2 may be specified.</p>
</td>
</tr>
<tr>
<td>
<code>uncompressedUserData</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UncompressedUserData specify whether the user data is gzip-compressed before it is sent to ec2 instance.
cloud-init has built-in support for gzip-compressed user data
user data stored in aws secret manager is always gzip-compressed.</p>
</td>
</tr>
<tr>
<td>
<code>cloudInit</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.CloudInit">
CloudInit
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Ignition">
Ignition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ignition defined options related to the bootstrapping systems where Ignition is used.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineStatus">AWSMachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachine">AWSMachine</a>)
</p>
<p>
<p>AWSMachineStatus defines the observed state of AWSMachine.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready is true when the provider resource is ready.</p>
</td>
</tr>
<tr>
<td>
<code>interruptible</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interruptible reports that this machine is using spot instances and can therefore be interrupted by CAPI when it receives a notice that the spot instance is to be terminated by AWS.
This will be set to true when SpotMarketOptions is not nil (i.e. this machine is using a spot instance).</p>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
[]Cluster API api/v1beta1.MachineAddress
</a>
</em>
</td>
<td>
<p>Addresses contains the AWS instance associated addresses.</p>
</td>
</tr>
<tr>
<td>
<code>instanceState</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.InstanceState">
InstanceState
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceState is the state of the AWS instance for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the Machine and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the Machine and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the AWSMachine.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplate">AWSMachineTemplate
</h3>
<p>
<p>AWSMachineTemplate is the schema for the Amazon EC2 Machine Templates API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateSpec">
AWSMachineTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateResource">
AWSMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateStatus">
AWSMachineTemplateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateResource">AWSMachineTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateSpec">AWSMachineTemplateSpec</a>)
</p>
<p>
<p>AWSMachineTemplateResource describes the data needed to create am AWSMachine from a template.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.
More info: <a href="https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata">https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata</a></p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">
AWSMachineSpec
</a>
</em>
</td>
<td>
<p>Spec is the specification of the desired behavior of the machine.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceID is the EC2 instance ID for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
AWSMachine&rsquo;s value takes precedence.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMInstanceProfile is a name of an IAM instance profile to assign to the instance</p>
</td>
</tr>
<tr>
<td>
<code>publicIP</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicIP specifies whether the instance should get a public IP.
Precedence for this setting is as follows:
1. This field if set
2. Cluster/flavor setting
3. Subnet default</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instance. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator. It is possible to specify either IDs of Filters. Using Filters
will cause additional requests to AWS API and if tags change the attached security groups might change too.</p>
</td>
</tr>
<tr>
<td>
<code>failureDomain</code><br/>
<em>
string
</em>
</td>
<td>
<p>FailureDomain is the failure domain unique identifier this Machine should be attached to, as defined in Cluster API.
For this infrastructure provider, the ID is equivalent to an AWS Availability Zone.
If multiple subnets are matched for the availability zone, the first one returned is picked.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">
AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnet is a reference to the subnet to use for this instance. If not specified,
the cluster subnet will be used.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces is a list of ENIs to associate with the instance.
A maximum of 2 may be specified.</p>
</td>
</tr>
<tr>
<td>
<code>uncompressedUserData</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UncompressedUserData specify whether the user data is gzip-compressed before it is sent to ec2 instance.
cloud-init has built-in support for gzip-compressed user data
user data stored in aws secret manager is always gzip-compressed.</p>
</td>
</tr>
<tr>
<td>
<code>cloudInit</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.CloudInit">
CloudInit
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Ignition">
Ignition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ignition defined options related to the bootstrapping systems where Ignition is used.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateSpec">AWSMachineTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplate">AWSMachineTemplate</a>)
</p>
<p>
<p>AWSMachineTemplateSpec defines the desired state of AWSMachineTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateResource">
AWSMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplateStatus">AWSMachineTemplateStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineTemplate">AWSMachineTemplate</a>)
</p>
<p>
<p>AWSMachineTemplateStatus defines a status for an AWSMachineTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>capacity</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcelist-v1-core">
Kubernetes core/v1.ResourceList
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Capacity defines the resource capacity for this machine.
This value is used for autoscaling from zero operations as defined in:
<a href="https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md">https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">AWSResourceReference
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec</a>)
</p>
<p>
<p>AWSResourceReference is a reference to a specific AWS resource by ID or filters.
Only one of ID or Filters may be specified. Specifying more than one will result in
a validation error.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ID of resource</p>
</td>
</tr>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ARN of resource.
Deprecated: This field has no function and is going to be removed in the next release.</p>
</td>
</tr>
<tr>
<td>
<code>filters</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Filter">
[]Filter
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Filters is a set of key/value pairs used to identify a resource
They are applied according to the rules defined by the AWS API:
<a href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html">https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSRoleSpec">AWSRoleSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec</a>)
</p>
<p>
<p>AWSRoleSpec defines the specifications for all identities based around AWS roles.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>roleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>The Amazon Resource Name (ARN) of the role to assume.</p>
</td>
</tr>
<tr>
<td>
<code>sessionName</code><br/>
<em>
string
</em>
</td>
<td>
<p>An identifier for the assumed role session</p>
</td>
</tr>
<tr>
<td>
<code>durationSeconds</code><br/>
<em>
int32
</em>
</td>
<td>
<p>The duration, in seconds, of the role session before it is renewed.</p>
</td>
</tr>
<tr>
<td>
<code>inlinePolicy</code><br/>
<em>
string
</em>
</td>
<td>
<p>An IAM policy as a JSON-encoded string that you want to use as an inline session policy.</p>
</td>
</tr>
<tr>
<td>
<code>policyARNs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The Amazon Resource Names (ARNs) of the IAM managed policies that you want
to use as managed session policies.
The policies must exist in the same account as the role.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AZSelectionScheme">AZSelectionScheme
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.VPCSpec">VPCSpec</a>)
</p>
<p>
<p>AZSelectionScheme defines the scheme of selecting AZs.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AllowedNamespaces">AllowedNamespaces
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterIdentitySpec">AWSClusterIdentitySpec</a>)
</p>
<p>
<p>AllowedNamespaces is a selector of namespaces that AWSClusters can
use this ClusterPrincipal from. This is a standard Kubernetes LabelSelector,
a label query over a set of resources. The result of matchLabels and
matchExpressions are ANDed.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>list</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>An nil or empty list indicates that AWSClusters cannot use the identity from any namespace.</p>
</td>
</tr>
<tr>
<td>
<code>selector</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>An empty selector indicates that AWSClusters cannot use this
AWSClusterIdentity from any namespace.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Bastion">Bastion
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">AWSClusterSpec</a>)
</p>
<p>
<p>Bastion defines a bastion host.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Enabled allows this provider to create a bastion host instance
with a public ip to access the VPC private network.</p>
</td>
</tr>
<tr>
<td>
<code>disableIngressRules</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>DisableIngressRules will ensure there are no Ingress rules in the bastion host&rsquo;s security group.
Requires AllowedCIDRBlocks to be empty.</p>
</td>
</tr>
<tr>
<td>
<code>allowedCIDRBlocks</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AllowedCIDRBlocks is a list of CIDR blocks allowed to access the bastion host.
They are set as ingress rules for the Bastion host&rsquo;s Security Group (defaults to 0.0.0.0/0).</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType will use the specified instance type for the bastion. If not specified,
Cluster API Provider AWS will use t3.micro for all regions except us-east-1, where t2.micro
will be the default.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMI will use the specified AMI to boot the bastion. If not specified,
the AMI will default to one picked out in public space.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.BuildParams">BuildParams
</h3>
<p>
<p>BuildParams is used to build tags around an aws resource.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Lifecycle</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ResourceLifecycle">
ResourceLifecycle
</a>
</em>
</td>
<td>
<p>Lifecycle determines the resource lifecycle.</p>
</td>
</tr>
<tr>
<td>
<code>ClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ClusterName is the cluster associated with the resource.</p>
</td>
</tr>
<tr>
<td>
<code>ResourceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ResourceID is the unique identifier of the resource to be tagged.</p>
</td>
</tr>
<tr>
<td>
<code>Name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Name is the name of the resource, it&rsquo;s applied as the tag &ldquo;Name&rdquo; on AWS.</p>
</td>
</tr>
<tr>
<td>
<code>Role</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Role is the role associated to the resource.</p>
</td>
</tr>
<tr>
<td>
<code>Additional</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Any additional tags to be added to the resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.CNIIngressRule">CNIIngressRule
</h3>
<p>
<p>CNIIngressRule defines an AWS ingress rule for CNI requirements.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>description</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroupProtocol">
SecurityGroupProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>fromPort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>toPort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.CNIIngressRules">CNIIngressRules
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1.CNIIngressRule</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.CNISpec">CNISpec</a>)
</p>
<p>
<p>CNIIngressRules is a slice of CNIIngressRule.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.CNISpec">CNISpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkSpec">NetworkSpec</a>)
</p>
<p>
<p>CNISpec defines configuration for CNI.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cniIngressRules</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.CNIIngressRules">
CNIIngressRules
</a>
</em>
</td>
<td>
<p>CNIIngressRules specify rules to apply to control plane and worker node security groups.
The source for the rule will be set to control plane and worker security group IDs.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ClassicELB">ClassicELB
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkStatus">NetworkStatus</a>)
</p>
<p>
<p>ClassicELB defines an AWS classic load balancer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The name of the load balancer. It must be unique within the set of load balancers
defined in the region. It also serves as identifier.</p>
</td>
</tr>
<tr>
<td>
<code>dnsName</code><br/>
<em>
string
</em>
</td>
<td>
<p>DNSName is the dns name of the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>scheme</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBScheme">
ClassicELBScheme
</a>
</em>
</td>
<td>
<p>Scheme is the load balancer scheme, either internet-facing or private.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones in the VPC attached to the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SubnetIDs is an array of subnets in the VPC attached to the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SecurityGroupIDs is an array of security groups assigned to the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>listeners</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBListener">
[]ClassicELBListener
</a>
</em>
</td>
<td>
<p>Listeners is an array of classic elb listeners associated with the load balancer. There must be at least one.</p>
</td>
</tr>
<tr>
<td>
<code>healthChecks</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBHealthCheck">
ClassicELBHealthCheck
</a>
</em>
</td>
<td>
<p>HealthCheck is the classic elb health check associated with the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>attributes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBAttributes">
ClassicELBAttributes
</a>
</em>
</td>
<td>
<p>Attributes defines extra attributes associated with the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Tags is a map of tags associated with the load balancer.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBAttributes">ClassicELBAttributes
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELB">ClassicELB</a>)
</p>
<p>
<p>ClassicELBAttributes defines extra attributes associated with a classic load balancer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>idleTimeout</code><br/>
<em>
<a href="https://golang.org/pkg/time/#Duration">
time.Duration
</a>
</em>
</td>
<td>
<p>IdleTimeout is time that the connection is allowed to be idle (no data
has been sent over the connection) before it is closed by the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>crossZoneLoadBalancing</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>CrossZoneLoadBalancing enables the classic load balancer load balancing.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBHealthCheck">ClassicELBHealthCheck
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELB">ClassicELB</a>)
</p>
<p>
<p>ClassicELBHealthCheck defines an AWS classic load balancer health check.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>target</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>interval</code><br/>
<em>
<a href="https://golang.org/pkg/time/#Duration">
time.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>timeout</code><br/>
<em>
<a href="https://golang.org/pkg/time/#Duration">
time.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>healthyThreshold</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>unhealthyThreshold</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBListener">ClassicELBListener
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELB">ClassicELB</a>)
</p>
<p>
<p>ClassicELBListener defines an AWS classic load balancer listener.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBProtocol">
ClassicELBProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>port</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instanceProtocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBProtocol">
ClassicELBProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instancePort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBProtocol">ClassicELBProtocol
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLoadBalancerSpec">AWSLoadBalancerSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBListener">ClassicELBListener</a>)
</p>
<p>
<p>ClassicELBProtocol defines listener protocols for a classic load balancer.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ClassicELBScheme">ClassicELBScheme
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLoadBalancerSpec">AWSLoadBalancerSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELB">ClassicELB</a>)
</p>
<p>
<p>ClassicELBScheme defines the scheme of a classic load balancer.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.CloudInit">CloudInit
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec</a>)
</p>
<p>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>insecureSkipSecretsManager</code><br/>
<em>
bool
</em>
</td>
<td>
<p>InsecureSkipSecretsManager, when set to true will not use AWS Secrets Manager
or AWS Systems Manager Parameter Store to ensure privacy of userdata.
By default, a cloud-init boothook shell script is prepended to download
the userdata from Secrets Manager and additionally delete the secret.</p>
</td>
</tr>
<tr>
<td>
<code>secretCount</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretCount is the number of secrets used to form the complete secret</p>
</td>
</tr>
<tr>
<td>
<code>secretPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretPrefix is the prefix for the secret name. This is stored
temporarily, and deleted when the machine registers as a node against
the workload cluster.</p>
</td>
</tr>
<tr>
<td>
<code>secureSecretsBackend</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SecretBackend">
SecretBackend
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecureSecretsBackend, when set to parameter-store will utilize the AWS Systems Manager
Parameter Storage to distribute secrets. By default or with the value of secrets-manager,
will use AWS Secrets Manager instead.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.EKSAMILookupType">EKSAMILookupType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AMIReference">AMIReference</a>)
</p>
<p>
<p>EKSAMILookupType specifies which AWS AMI to use for a AWSMachine and AWSMachinePool.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Filter">Filter
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSResourceReference">AWSResourceReference</a>)
</p>
<p>
<p>Filter is a filter used to identify an AWS resource.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the filter. Filter names are case-sensitive.</p>
</td>
</tr>
<tr>
<td>
<code>values</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Values includes one or more filter values. Filter values are case-sensitive.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.IPv6">IPv6
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.VPCSpec">VPCSpec</a>)
</p>
<p>
<p>IPv6 contains ipv6 specific settings for the network.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>CidrBlock is the CIDR block provided by Amazon when VPC has enabled IPv6.</p>
</td>
</tr>
<tr>
<td>
<code>poolId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PoolID is the IP pool which must be defined in case of BYO IP is defined.</p>
</td>
</tr>
<tr>
<td>
<code>egressOnlyInternetGatewayId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EgressOnlyInternetGatewayID is the id of the egress only internet gateway associated with an IPv6 enabled VPC.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Ignition">Ignition
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec</a>)
</p>
<p>
<p>Ignition defines options related to the bootstrapping systems where Ignition is used.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines which version of Ignition will be used to generate bootstrap data.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.IngressRule">IngressRule
</h3>
<p>
<p>IngressRule defines an AWS ingress rule for security groups.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>description</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroupProtocol">
SecurityGroupProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>fromPort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>toPort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cidrBlocks</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>List of CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.</p>
</td>
</tr>
<tr>
<td>
<code>ipv6CidrBlocks</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>List of IPv6 CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.</p>
</td>
</tr>
<tr>
<td>
<code>sourceSecurityGroupIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The security group id to allow access from. Cannot be specified with CidrBlocks.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.IngressRules">IngressRules
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1.IngressRule</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroup">SecurityGroup</a>)
</p>
<p>
<p>IngressRules is a slice of AWS ingress rules for security groups.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Instance">Instance
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStatus">AWSClusterStatus</a>)
</p>
<p>
<p>Instance describes an AWS instance.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instanceState</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.InstanceState">
InstanceState
</a>
</em>
</td>
<td>
<p>The current state of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
<p>The instance type.</p>
</td>
</tr>
<tr>
<td>
<code>subnetId</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the subnet of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageId</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the AMI used to launch the instance.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the SSH key pair.</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SecurityGroupIDs are one or more security group IDs this instance belongs to.</p>
</td>
</tr>
<tr>
<td>
<code>userData</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserData is the raw data script passed to the instance which is run upon bootstrap.
This field must not be base64 encoded and should only be used when running a new instance.</p>
</td>
</tr>
<tr>
<td>
<code>iamProfile</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the IAM instance profile associated with the instance, if applicable.</p>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
[]Cluster API api/v1beta1.MachineAddress
</a>
</em>
</td>
<td>
<p>Addresses contains the AWS instance associated addresses.</p>
</td>
</tr>
<tr>
<td>
<code>privateIp</code><br/>
<em>
string
</em>
</td>
<td>
<p>The private IPv4 address assigned to the instance.</p>
</td>
</tr>
<tr>
<td>
<code>publicIp</code><br/>
<em>
string
</em>
</td>
<td>
<p>The public IPv4 address assigned to the instance, if applicable.</p>
</td>
</tr>
<tr>
<td>
<code>enaSupport</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Specifies whether enhanced networking with ENA is enabled.</p>
</td>
</tr>
<tr>
<td>
<code>ebsOptimized</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Indicates whether the instance is optimized for Amazon EBS I/O.</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the root storage volume.</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Specifies ENIs attached to instance</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>The tags associated with the instance.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZone</code><br/>
<em>
string
</em>
</td>
<td>
<p>Availability zone of instance</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<p>SpotMarketOptions option for configuring instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
<tr>
<td>
<code>volumeIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IDs of the instance&rsquo;s volumes</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.InstanceState">InstanceState
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineStatus">AWSMachineStatus</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.Instance">Instance</a>)
</p>
<p>
<p>InstanceState describes the state of an AWS instance.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.NetworkSpec">NetworkSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">AWSClusterSpec</a>)
</p>
<p>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>vpc</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.VPCSpec">
VPCSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>VPC configuration.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Subnets">
Subnets
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets configuration.</p>
</td>
</tr>
<tr>
<td>
<code>cni</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.CNISpec">
CNISpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CNI configuration</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupOverrides</code><br/>
<em>
map[sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1.SecurityGroupRole]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecurityGroupOverrides is an optional set of security groups to use for cluster instances
This is optional - if not provided new security groups will be created for the cluster</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.NetworkStatus">NetworkStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterStatus">AWSClusterStatus</a>)
</p>
<p>
<p>NetworkStatus encapsulates AWS networking resources.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>securityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroup">
map[sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1.SecurityGroupRole]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1.SecurityGroup
</a>
</em>
</td>
<td>
<p>SecurityGroups is a map from the role/kind of the security group to its unique name, if any.</p>
</td>
</tr>
<tr>
<td>
<code>apiServerElb</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ClassicELB">
ClassicELB
</a>
</em>
</td>
<td>
<p>APIServerELB is the Kubernetes api server classic load balancer.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ResourceLifecycle">ResourceLifecycle
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.BuildParams">BuildParams</a>)
</p>
<p>
<p>ResourceLifecycle configures the lifecycle of a resource.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.RouteTable">RouteTable
</h3>
<p>
<p>RouteTable defines an AWS routing table.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.S3Bucket">S3Bucket
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">AWSClusterSpec</a>)
</p>
<p>
<p>S3Bucket defines a supporting S3 bucket for the cluster, currently can be optionally used for Ignition.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>controlPlaneIAMInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<p>ControlPlaneIAMInstanceProfile is a name of the IAMInstanceProfile, which will be allowed
to read control-plane node bootstrap data from S3 Bucket.</p>
</td>
</tr>
<tr>
<td>
<code>nodesIAMInstanceProfiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>NodesIAMInstanceProfiles is a list of IAM instance profiles, which will be allowed to read
worker nodes bootstrap data from S3 Bucket.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name defines name of S3 Bucket to be created.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.SecretBackend">SecretBackend
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.CloudInit">CloudInit</a>)
</p>
<p>
<p>SecretBackend defines variants for backend secret storage.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroup">SecurityGroup
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkStatus">NetworkStatus</a>)
</p>
<p>
<p>SecurityGroup defines an AWS security group.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is a unique identifier.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the security group name.</p>
</td>
</tr>
<tr>
<td>
<code>ingressRule</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.IngressRules">
IngressRules
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IngressRules is the inbound rules associated with the security group.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a map of tags associated with the security group.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroupProtocol">SecurityGroupProtocol
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.CNIIngressRule">CNIIngressRule</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.IngressRule">IngressRule</a>)
</p>
<p>
<p>SecurityGroupProtocol defines the protocol type for a security group rule.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroupRole">SecurityGroupRole
(<code>string</code> alias)</p></h3>
<p>
<p>SecurityGroupRole defines the unique role of a security group.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.SpotMarketOptions">SpotMarketOptions
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.Instance">Instance</a>)
</p>
<p>
<p>SpotMarketOptions defines the options available to a user when configuring
Machines to run on Spot instances.
Most users should provide an empty struct.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxPrice</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxPrice defines the maximum price the user is willing to pay for Spot VM instances</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.SubnetSpec">SubnetSpec
</h3>
<p>
<p>SubnetSpec configures an AWS Subnet.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID defines a unique identifier to reference this resource.</p>
</td>
</tr>
<tr>
<td>
<code>cidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<p>CidrBlock is the CIDR block to be used when the provider creates a managed VPC.</p>
</td>
</tr>
<tr>
<td>
<code>ipv6CidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPv6CidrBlock is the IPv6 CIDR block to be used when the provider creates a managed VPC.
A subnet can have an IPv4 and an IPv6 address.
IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZone</code><br/>
<em>
string
</em>
</td>
<td>
<p>AvailabilityZone defines the availability zone to use for this subnet in the cluster&rsquo;s region.</p>
</td>
</tr>
<tr>
<td>
<code>isPublic</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>IsPublic defines the subnet as a public subnet. A subnet is public when it is associated with a route table that has a route to an internet gateway.</p>
</td>
</tr>
<tr>
<td>
<code>isIpv6</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>IsIPv6 defines the subnet as an IPv6 subnet. A subnet is IPv6 when it is associated with a VPC that has IPv6 enabled.
IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.</p>
</td>
</tr>
<tr>
<td>
<code>routeTableId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RouteTableID is the routing table id associated with the subnet.</p>
</td>
</tr>
<tr>
<td>
<code>natGatewayId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NatGatewayID is the NAT gateway id associated with the subnet.
Ignored unless the subnet is managed by the provider, in which case this is set on the public subnet where the NAT gateway resides. It is then used to determine routes for private subnets in the same AZ as the public subnet.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a collection of tags describing the resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Subnets">Subnets
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1.SubnetSpec</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkSpec">NetworkSpec</a>)
</p>
<p>
<p>Subnets is a slice of Subnet.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Tags">Tags
(<code>map[string]string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSClusterSpec">AWSClusterSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.BuildParams">BuildParams</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.SecurityGroup">SecurityGroup</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.SubnetSpec">SubnetSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.VPCSpec">VPCSpec</a>)
</p>
<p>
<p>Tags defines a map of tags.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.VPCSpec">VPCSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.NetworkSpec">NetworkSpec</a>)
</p>
<p>
<p>VPCSpec configures an AWS VPC.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is the vpc-id of the VPC this provider should use to create resources.</p>
</td>
</tr>
<tr>
<td>
<code>cidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<p>CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
Defaults to 10.0.0.0/16.</p>
</td>
</tr>
<tr>
<td>
<code>ipv6</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.IPv6">
IPv6
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPv6 contains ipv6 specific settings for the network. Supported only in managed clusters.
This field cannot be set on AWSCluster object.</p>
</td>
</tr>
<tr>
<td>
<code>internetGatewayId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InternetGatewayID is the id of the internet gateway associated with the VPC.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a collection of tags describing the resource.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneUsageLimit</code><br/>
<em>
int
</em>
</td>
<td>
<p>AvailabilityZoneUsageLimit specifies the maximum number of availability zones (AZ) that
should be used in a region when automatically creating subnets. If a region has more
than this number of AZs then this number of AZs will be picked randomly when creating
default subnets. Defaults to 3</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneSelection</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AZSelectionScheme">
AZSelectionScheme
</a>
</em>
</td>
<td>
<p>AvailabilityZoneSelection specifies how AZs should be selected if there are more AZs
in a region than specified by AvailabilityZoneUsageLimit. There are 2 selection schemes:
Ordered - selects based on alphabetical order
Random - selects AZs randomly in a region
Defaults to Ordered</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Volume">Volume
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.Instance">Instance</a>)
</p>
<p>
<p>Volume encapsulates the configuration options for the storage device.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>deviceName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Device name</p>
</td>
</tr>
<tr>
<td>
<code>size</code><br/>
<em>
int64
</em>
</td>
<td>
<p>Size specifies size (in Gi) of the storage device.
Must be greater than the image snapshot size or 8 (whichever is greater).</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.VolumeType">
VolumeType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Type is the type of the volume (e.g. gp2, io1, etc&hellip;).</p>
</td>
</tr>
<tr>
<td>
<code>iops</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>IOPS is the number of IOPS requested for the disk. Not applicable to all types.</p>
</td>
</tr>
<tr>
<td>
<code>throughput</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>Throughput to provision in MiB/s supported for the volume type. Not applicable to all types.</p>
</td>
</tr>
<tr>
<td>
<code>encrypted</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Encrypted is whether the volume should be encrypted or not.</p>
</td>
</tr>
<tr>
<td>
<code>encryptionKey</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptionKey is the KMS key to use to encrypt the volume. Can be either a KMS key ID or ARN.
If Encrypted is set and this is omitted, the default AWS key will be used.
The key must already exist and be accessible by the controller.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.VolumeType">VolumeType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.Volume">Volume</a>)
</p>
<p>
<p>VolumeType describes the EBS volume type.
See: <a href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html">https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html</a></p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ASGStatus">ASGStatus
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolStatus">AWSMachinePoolStatus</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AutoScalingGroup">AutoScalingGroup</a>)
</p>
<p>
<p>ASGStatus is a status string returned by the autoscaling API.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSFargateProfile">AWSFargateProfile
</h3>
<p>
<p>AWSFargateProfile is the Schema for the awsfargateprofiles API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.FargateProfileSpec">
FargateProfileSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>clusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ClusterName is the name of the Cluster this object belongs to.</p>
</td>
</tr>
<tr>
<td>
<code>profileName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProfileName specifies the profile name.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for this fargate pool
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>selectors</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.FargateSelector">
[]FargateSelector
</a>
</em>
</td>
<td>
<p>Selectors specify fargate pod selectors.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.FargateProfileStatus">
FargateProfileStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">AWSLaunchTemplate
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>AWSLaunchTemplate defines the desired state of AWSLaunchTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the launch template.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name or the Amazon Resource Name (ARN) of the instance profile associated
with the IAM role for the instance. The instance profile contains the IAM
role.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string
(do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>versionNumber</code><br/>
<em>
int64
</em>
</td>
<td>
<p>VersionNumber is the version of the launch template that is applied.
Typically a new version is created when at least one of the following happens:
1) A new launch template spec is applied.
2) One or more parameters in an existing template is changed.
3) A new AMI is discovered.</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instances. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<p>SpotMarketOptions are options for configuring AWSMachinePool instances to be run using AWS Spot instances.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePool">AWSMachinePool
</h3>
<p>
<p>AWSMachinePool is the Schema for the awsmachinepools API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolSpec">
AWSMachinePoolSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the ARN of the associated ASG</p>
</td>
</tr>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MinSize defines the minimum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MaxSize defines the maximum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets is an array of subnet configurations</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<p>AWSLaunchTemplate specifies the launch template and version to use when an instance is launched.</p>
</td>
</tr>
<tr>
<td>
<code>mixedInstancesPolicy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.MixedInstancesPolicy">
MixedInstancesPolicy
</a>
</em>
</td>
<td>
<p>MixedInstancesPolicy describes how multiple instance types will be used by the ASG.</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the identification IDs of machine instances provided by the provider.
This field must match the provider IDs as seen on the node objects corresponding to a machine pool&rsquo;s machine instances.</p>
</td>
</tr>
<tr>
<td>
<code>defaultCoolDown</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of time, in seconds, after a scaling activity completes before another scaling activity can start.
If no value is supplied by user a default value of 300 seconds is set</p>
</td>
</tr>
<tr>
<td>
<code>refreshPreferences</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.RefreshPreferences">
RefreshPreferences
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RefreshPreferences describes set of preferences associated with the instance refresh request.</p>
</td>
</tr>
<tr>
<td>
<code>capacityRebalance</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Enable or disable the capacity rebalance autoscaling group feature</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolStatus">
AWSMachinePoolStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolInstanceStatus">AWSMachinePoolInstanceStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolStatus">AWSMachinePoolStatus</a>)
</p>
<p>
<p>AWSMachinePoolInstanceStatus defines the status of the AWSMachinePoolInstance.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceID is the identification of the Machine Instance within ASG</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines the Kubernetes version for the Machine Instance</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolSpec">AWSMachinePoolSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePool">AWSMachinePool</a>)
</p>
<p>
<p>AWSMachinePoolSpec defines the desired state of AWSMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the ARN of the associated ASG</p>
</td>
</tr>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MinSize defines the minimum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MaxSize defines the maximum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets is an array of subnet configurations</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<p>AWSLaunchTemplate specifies the launch template and version to use when an instance is launched.</p>
</td>
</tr>
<tr>
<td>
<code>mixedInstancesPolicy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.MixedInstancesPolicy">
MixedInstancesPolicy
</a>
</em>
</td>
<td>
<p>MixedInstancesPolicy describes how multiple instance types will be used by the ASG.</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the identification IDs of machine instances provided by the provider.
This field must match the provider IDs as seen on the node objects corresponding to a machine pool&rsquo;s machine instances.</p>
</td>
</tr>
<tr>
<td>
<code>defaultCoolDown</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of time, in seconds, after a scaling activity completes before another scaling activity can start.
If no value is supplied by user a default value of 300 seconds is set</p>
</td>
</tr>
<tr>
<td>
<code>refreshPreferences</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.RefreshPreferences">
RefreshPreferences
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RefreshPreferences describes set of preferences associated with the instance refresh request.</p>
</td>
</tr>
<tr>
<td>
<code>capacityRebalance</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Enable or disable the capacity rebalance autoscaling group feature</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolStatus">AWSMachinePoolStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePool">AWSMachinePool</a>)
</p>
<p>
<p>AWSMachinePoolStatus defines the observed state of AWSMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready is true when the provider resource is ready.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Replicas is the most recently observed number of replicas</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the AWSMachinePool.</p>
</td>
</tr>
<tr>
<td>
<code>instances</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolInstanceStatus">
[]AWSMachinePoolInstanceStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Instances contains the status for each instance in the pool</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateID</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The version of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the Machine and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the Machine and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>asgStatus</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ASGStatus">
ASGStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePool">AWSManagedMachinePool
</h3>
<p>
<p>AWSManagedMachinePool is the Schema for the awsmanagedmachinepools API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">
AWSManagedMachinePoolSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>eksNodegroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSNodegroupName specifies the name of the nodegroup in AWS
corresponding to this MachinePool. If you don&rsquo;t specify a name
then a default name will be created based on the namespace and
name of the managed machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the node group role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for the node group.
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>amiVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIVersion defines the desired AMI release version. If no version number
is supplied then the latest version for the Kubernetes version
will be used</p>
</td>
</tr>
<tr>
<td>
<code>amiType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachineAMIType">
ManagedMachineAMIType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIType defines the AMI type</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels specifies labels for the Kubernetes node objects</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Taints">
Taints
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints specifies the taints to apply to the nodes of the machine pool</p>
</td>
</tr>
<tr>
<td>
<code>diskSize</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>DiskSize specifies the root disk size</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceType specifies the AWS instance type</p>
</td>
</tr>
<tr>
<td>
<code>scaling</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachinePoolScaling">
ManagedMachinePoolScaling
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Scaling specifies scaling for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>remoteAccess</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedRemoteAccess">
ManagedRemoteAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RemoteAccess specifies how machines can be accessed remotely</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the provider IDs of instances in the
autoscaling group corresponding to the nodegroup represented by this
machine pool</p>
</td>
</tr>
<tr>
<td>
<code>capacityType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachinePoolCapacityType">
ManagedMachinePoolCapacityType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CapacityType specifies the capacity type for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>updateConfig</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.UpdateConfig">
UpdateConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>UpdateConfig holds the optional config to control the behaviour of the update
to the nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AWSLaunchTemplate specifies the launch template to use to create the managed node group.
If AWSLaunchTemplate is specified, certain node group configuraions outside of launch template
are prohibited (<a href="https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html">https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html</a>).</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolStatus">
AWSManagedMachinePoolStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePool">AWSManagedMachinePool</a>)
</p>
<p>
<p>AWSManagedMachinePoolSpec defines the desired state of AWSManagedMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>eksNodegroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSNodegroupName specifies the name of the nodegroup in AWS
corresponding to this MachinePool. If you don&rsquo;t specify a name
then a default name will be created based on the namespace and
name of the managed machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the node group role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for the node group.
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>amiVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIVersion defines the desired AMI release version. If no version number
is supplied then the latest version for the Kubernetes version
will be used</p>
</td>
</tr>
<tr>
<td>
<code>amiType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachineAMIType">
ManagedMachineAMIType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIType defines the AMI type</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels specifies labels for the Kubernetes node objects</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Taints">
Taints
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints specifies the taints to apply to the nodes of the machine pool</p>
</td>
</tr>
<tr>
<td>
<code>diskSize</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>DiskSize specifies the root disk size</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceType specifies the AWS instance type</p>
</td>
</tr>
<tr>
<td>
<code>scaling</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachinePoolScaling">
ManagedMachinePoolScaling
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Scaling specifies scaling for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>remoteAccess</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedRemoteAccess">
ManagedRemoteAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RemoteAccess specifies how machines can be accessed remotely</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the provider IDs of instances in the
autoscaling group corresponding to the nodegroup represented by this
machine pool</p>
</td>
</tr>
<tr>
<td>
<code>capacityType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachinePoolCapacityType">
ManagedMachinePoolCapacityType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CapacityType specifies the capacity type for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>updateConfig</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.UpdateConfig">
UpdateConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>UpdateConfig holds the optional config to control the behaviour of the update
to the nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AWSLaunchTemplate specifies the launch template to use to create the managed node group.
If AWSLaunchTemplate is specified, certain node group configuraions outside of launch template
are prohibited (<a href="https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html">https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html</a>).</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolStatus">AWSManagedMachinePoolStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePool">AWSManagedMachinePool</a>)
</p>
<p>
<p>AWSManagedMachinePoolStatus defines the observed state of AWSManagedMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the AWSManagedMachinePool nodegroup has joined
the cluster</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Replicas is the most recently observed number of replicas.</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The ID of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The version of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the MachinePool and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of MachinePools
can be added as events to the MachinePool object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the MachinePool and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the MachinePool&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of MachinePools
can be added as events to the MachinePool object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the managed machine pool</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.AutoScalingGroup">AutoScalingGroup
</h3>
<p>
<p>AutoScalingGroup describes an AWS autoscaling group.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>The tags associated with the instance.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>desiredCapacity</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>placementGroup</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>defaultCoolDown</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>capacityRebalance</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>mixedInstancesPolicy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.MixedInstancesPolicy">
MixedInstancesPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.ASGStatus">
ASGStatus
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instances</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">
[]Instance
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.BlockDeviceMapping">BlockDeviceMapping
</h3>
<p>
<p>BlockDeviceMapping specifies the block devices for the instance.
You can specify virtual devices and EBS volumes.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>deviceName</code><br/>
<em>
string
</em>
</td>
<td>
<p>The device name exposed to the EC2 instance (for example, /dev/sdh or xvdh).</p>
</td>
</tr>
<tr>
<td>
<code>ebs</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.EBS">
EBS
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>You can specify either VirtualName or Ebs, but not both.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.EBS">EBS
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.BlockDeviceMapping">BlockDeviceMapping</a>)
</p>
<p>
<p>EBS can be used to automatically set up EBS volumes when an instance is launched.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>encrypted</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Encrypted is whether the volume should be encrypted or not.</p>
</td>
</tr>
<tr>
<td>
<code>volumeSize</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>The size of the volume, in GiB.
This can be a number from 1-1,024 for standard, 4-16,384 for io1, 1-16,384
for gp2, and 500-16,384 for st1 and sc1. If you specify a snapshot, the volume
size must be equal to or larger than the snapshot size.</p>
</td>
</tr>
<tr>
<td>
<code>volumeType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The volume type
For more information, see Amazon EBS Volume Types (<a href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html">https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html</a>)</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.FargateProfileSpec">FargateProfileSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSFargateProfile">AWSFargateProfile</a>)
</p>
<p>
<p>FargateProfileSpec defines the desired state of FargateProfile.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>clusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ClusterName is the name of the Cluster this object belongs to.</p>
</td>
</tr>
<tr>
<td>
<code>profileName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProfileName specifies the profile name.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for this fargate pool
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>selectors</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.FargateSelector">
[]FargateSelector
</a>
</em>
</td>
<td>
<p>Selectors specify fargate pod selectors.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.FargateProfileStatus">FargateProfileStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSFargateProfile">AWSFargateProfile</a>)
</p>
<p>
<p>FargateProfileStatus defines the observed state of FargateProfile.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the FargateProfile is available.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the FargateProfile and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the FargateProfile&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of
FargateProfiles can be added as events to the FargateProfile object
and/or logged in the controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the FargateProfile and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the FargateProfile&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of
FargateProfiles can be added as events to the FargateProfile
object and/or logged in the controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current state of the Fargate profile.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.FargateSelector">FargateSelector
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.FargateProfileSpec">FargateProfileSpec</a>)
</p>
<p>
<p>FargateSelector specifies a selector for pods that should run on this fargate pool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Labels specifies which pod labels this selector should match.</p>
</td>
</tr>
<tr>
<td>
<code>namespace</code><br/>
<em>
string
</em>
</td>
<td>
<p>Namespace specifies which namespace this selector should match.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.InstancesDistribution">InstancesDistribution
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.MixedInstancesPolicy">MixedInstancesPolicy</a>)
</p>
<p>
<p>InstancesDistribution to configure distribution of On-Demand Instances and Spot Instances.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>onDemandAllocationStrategy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.OnDemandAllocationStrategy">
OnDemandAllocationStrategy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>spotAllocationStrategy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.SpotAllocationStrategy">
SpotAllocationStrategy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>onDemandBaseCapacity</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>onDemandPercentageAboveBaseCapacity</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachineAMIType">ManagedMachineAMIType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedMachineAMIType specifies which AWS AMI to use for a managed MachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;AL2_ARM_64&#34;</p></td>
<td><p>Al2Arm64 is the Arm AMI type.</p>
</td>
</tr><tr><td><p>&#34;AL2_x86_64&#34;</p></td>
<td><p>Al2x86_64 is the default AMI type.</p>
</td>
</tr><tr><td><p>&#34;AL2_x86_64_GPU&#34;</p></td>
<td><p>Al2x86_64GPU is the x86-64 GPU AMI type.</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachinePoolCapacityType">ManagedMachinePoolCapacityType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedMachinePoolCapacityType specifies the capacity type to be used for the managed MachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;onDemand&#34;</p></td>
<td><p>ManagedMachinePoolCapacityTypeOnDemand is the default capacity type, to launch on-demand instances.</p>
</td>
</tr><tr><td><p>&#34;spot&#34;</p></td>
<td><p>ManagedMachinePoolCapacityTypeSpot is the spot instance capacity type to launch spot instances.</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ManagedMachinePoolScaling">ManagedMachinePoolScaling
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedMachinePoolScaling specifies scaling options.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.ManagedRemoteAccess">ManagedRemoteAccess
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedRemoteAccess specifies remote access settings for EC2 instances.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<p>SSHKeyName specifies which EC2 SSH key can be used to access machines.
If left empty, the key from the control plane is used.</p>
</td>
</tr>
<tr>
<td>
<code>sourceSecurityGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SourceSecurityGroups specifies which security groups are allowed access</p>
</td>
</tr>
<tr>
<td>
<code>public</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Public specifies whether to open port 22 to the public internet</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.MixedInstancesPolicy">MixedInstancesPolicy
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AutoScalingGroup">AutoScalingGroup</a>)
</p>
<p>
<p>MixedInstancesPolicy for an Auto Scaling group.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instancesDistribution</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.InstancesDistribution">
InstancesDistribution
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>overrides</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.Overrides">
[]Overrides
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.OnDemandAllocationStrategy">OnDemandAllocationStrategy
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.InstancesDistribution">InstancesDistribution</a>)
</p>
<p>
<p>OnDemandAllocationStrategy indicates how to allocate instance types to fulfill On-Demand capacity.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Overrides">Overrides
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.MixedInstancesPolicy">MixedInstancesPolicy</a>)
</p>
<p>
<p>Overrides are used to override the instance type specified by the launch template with multiple
instance types that can be used to launch On-Demand Instances and Spot Instances.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.RefreshPreferences">RefreshPreferences
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolSpec">AWSMachinePoolSpec</a>)
</p>
<p>
<p>RefreshPreferences defines the specs for instance refreshing.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>strategy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The strategy to use for the instance refresh. The only valid value is Rolling.
A rolling update is an update that is applied to all instances in an Auto
Scaling group until all instances have been updated.</p>
</td>
</tr>
<tr>
<td>
<code>instanceWarmup</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>The number of seconds until a newly launched instance is configured and ready
to use. During this time, the next replacement will not be initiated.
The default is to use the value for the health check grace period defined for the group.</p>
</td>
</tr>
<tr>
<td>
<code>minHealthyPercentage</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of capacity as a percentage in ASG that must remain healthy
during an instance refresh. The default is 90.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.SpotAllocationStrategy">SpotAllocationStrategy
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.InstancesDistribution">InstancesDistribution</a>)
</p>
<p>
<p>SpotAllocationStrategy indicates how to allocate instances across Spot Instance pools.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Tags">Tags
(<code>map[string]string</code> alias)</p></h3>
<p>
<p>Tags is a mapping for tags.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Taint">Taint
</h3>
<p>
<p>Taint defines the specs for a Kubernetes taint.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>effect</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta1.TaintEffect">
TaintEffect
</a>
</em>
</td>
<td>
<p>Effect specifies the effect for the taint</p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>Key is the key of the taint</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
string
</em>
</td>
<td>
<p>Value is the value of the taint</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.TaintEffect">TaintEffect
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.Taint">Taint</a>)
</p>
<p>
<p>TaintEffect is the effect for a Kubernetes taint.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.Taints">Taints
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta1.Taint</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>Taints is an array of Taints.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta1.UpdateConfig">UpdateConfig
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>UpdateConfig is the configuration options for updating a nodegroup. Only one of MaxUnavailable
and MaxUnavailablePercentage should be specified.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxUnavailable</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxUnavailable is the maximum number of nodes unavailable at once during a version update.
Nodes will be updated in parallel. The maximum number is 100.</p>
</td>
</tr>
<tr>
<td>
<code>maxUnavailablePrecentage</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxUnavailablePercentage is the maximum percentage of nodes unavailable during a version update. This
percentage of nodes will be updated in parallel, up to 100 nodes at once.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="infrastructure.cluster.x-k8s.io/v1beta2">infrastructure.cluster.x-k8s.io/v1beta2</h2>
<p>
<p>Package v1beta2 contains the v1beta2 API implementation.</p>
</p>
Resource Types:
<ul></ul>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AMIReference">AMIReference
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">AWSLaunchTemplate</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">AWSLaunchTemplate</a>)
</p>
<p>
<p>AMIReference is a reference to a specific AWS resource by ID, ARN, or filters.
Only one of ID, ARN or Filters may be specified. Specifying more than one will result in
a validation error.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ID of resource</p>
</td>
</tr>
<tr>
<td>
<code>eksLookupType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.EKSAMILookupType">
EKSAMILookupType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSOptimizedLookupType If specified, will look up an EKS Optimized image in SSM Parameter store</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSCluster">AWSCluster
</h3>
<p>
<p>AWSCluster is the schema for Amazon EC2 based Kubernetes Cluster API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">
AWSClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryControlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecondaryControlPlaneLoadBalancer is an additional load balancer that can be used for the control plane.</p>
<p>An example use case is to have a separate internal load balancer for internal traffic,
and a separate external load balancer for external traffic.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.S3Bucket">
S3Bucket
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Bucket contains options to configure a supporting S3 bucket for this
cluster - currently used for nodes requiring Ignition
(<a href="https://coreos.github.io/ignition/">https://coreos.github.io/ignition/</a>) for bootstrapping (requires
BootstrapFormatIgnition feature flag to be enabled).</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStatus">
AWSClusterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterControllerIdentity">AWSClusterControllerIdentity
</h3>
<p>
<p>AWSClusterControllerIdentity is the Schema for the awsclustercontrolleridentities API
It is used to grant access to use Cluster API Provider AWS Controller credentials.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterControllerIdentitySpec">
AWSClusterControllerIdentitySpec
</a>
</em>
</td>
<td>
<p>Spec for this AWSClusterControllerIdentity.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterControllerIdentitySpec">AWSClusterControllerIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterControllerIdentity">AWSClusterControllerIdentity</a>)
</p>
<p>
<p>AWSClusterControllerIdentitySpec defines the specifications for AWSClusterControllerIdentity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">AWSClusterIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterControllerIdentitySpec">AWSClusterControllerIdentitySpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStaticIdentitySpec">AWSClusterStaticIdentitySpec</a>)
</p>
<p>
<p>AWSClusterIdentitySpec defines the Spec struct for AWSClusterIdentity types.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>allowedNamespaces</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AllowedNamespaces">
AllowedNamespaces
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AllowedNamespaces is used to identify which namespaces are allowed to use the identity from.
Namespaces can be selected either using an array of namespaces or with label selector.
An empty allowedNamespaces object indicates that AWSClusters can use this identity from any namespace.
If this object is nil, no namespaces will be allowed (default behaviour, if this field is not provided)
A namespace should be either in the NamespaceList or match with Selector to use the identity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterRoleIdentity">AWSClusterRoleIdentity
</h3>
<p>
<p>AWSClusterRoleIdentity is the Schema for the awsclusterroleidentities API
It is used to assume a role using the provided sourceRef.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterRoleIdentitySpec">
AWSClusterRoleIdentitySpec
</a>
</em>
</td>
<td>
<p>Spec for this AWSClusterRoleIdentity.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>AWSRoleSpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSRoleSpec">
AWSRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>externalID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>A unique identifier that might be required when you assume a role in another account.
If the administrator of the account to which the role belongs provided you with an
external ID, then provide that value in the ExternalId parameter. This value can be
any string, such as a passphrase or account number. A cross-account role is usually
set up to trust everyone in an account. Therefore, the administrator of the trusting
account might send an external ID to the administrator of the trusted account. That
way, only someone with the ID can assume the role, rather than everyone in the
account. For more information about the external ID, see How to Use an External ID
When Granting Access to Your AWS Resources to a Third Party in the IAM User Guide.</p>
</td>
</tr>
<tr>
<td>
<code>sourceIdentityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>SourceIdentityRef is a reference to another identity which will be chained to do
role assumption. All identity types are accepted.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterRoleIdentity">AWSClusterRoleIdentity</a>)
</p>
<p>
<p>AWSClusterRoleIdentitySpec defines the specifications for AWSClusterRoleIdentity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>AWSRoleSpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSRoleSpec">
AWSRoleSpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSRoleSpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>externalID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>A unique identifier that might be required when you assume a role in another account.
If the administrator of the account to which the role belongs provided you with an
external ID, then provide that value in the ExternalId parameter. This value can be
any string, such as a passphrase or account number. A cross-account role is usually
set up to trust everyone in an account. Therefore, the administrator of the trusting
account might send an external ID to the administrator of the trusted account. That
way, only someone with the ID can assume the role, rather than everyone in the
account. For more information about the external ID, see How to Use an External ID
When Granting Access to Your AWS Resources to a Third Party in the IAM User Guide.</p>
</td>
</tr>
<tr>
<td>
<code>sourceIdentityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>SourceIdentityRef is a reference to another identity which will be chained to do
role assumption. All identity types are accepted.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">AWSClusterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSCluster">AWSCluster</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplateResource">AWSClusterTemplateResource</a>)
</p>
<p>
<p>AWSClusterSpec defines the desired state of an EC2-based Kubernetes cluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryControlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecondaryControlPlaneLoadBalancer is an additional load balancer that can be used for the control plane.</p>
<p>An example use case is to have a separate internal load balancer for internal traffic,
and a separate external load balancer for external traffic.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.S3Bucket">
S3Bucket
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Bucket contains options to configure a supporting S3 bucket for this
cluster - currently used for nodes requiring Ignition
(<a href="https://coreos.github.io/ignition/">https://coreos.github.io/ignition/</a>) for bootstrapping (requires
BootstrapFormatIgnition feature flag to be enabled).</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStaticIdentity">AWSClusterStaticIdentity
</h3>
<p>
<p>AWSClusterStaticIdentity is the Schema for the awsclusterstaticidentities API
It represents a reference to an AWS access key ID and secret access key, stored in a secret.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStaticIdentitySpec">
AWSClusterStaticIdentitySpec
</a>
</em>
</td>
<td>
<p>Spec for this AWSClusterStaticIdentity</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
string
</em>
</td>
<td>
<p>Reference to a secret containing the credentials. The secret should
contain the following data keys:
AccessKeyID: AKIAIOSFODNN7EXAMPLE
SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
SessionToken: Optional</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStaticIdentitySpec">AWSClusterStaticIdentitySpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStaticIdentity">AWSClusterStaticIdentity</a>)
</p>
<p>
<p>AWSClusterStaticIdentitySpec defines the specifications for AWSClusterStaticIdentity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>AWSClusterIdentitySpec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">
AWSClusterIdentitySpec
</a>
</em>
</td>
<td>
<p>
(Members of <code>AWSClusterIdentitySpec</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br/>
<em>
string
</em>
</td>
<td>
<p>Reference to a secret containing the credentials. The secret should
contain the following data keys:
AccessKeyID: AKIAIOSFODNN7EXAMPLE
SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
SessionToken: Optional</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStatus">AWSClusterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSCluster">AWSCluster</a>)
</p>
<p>
<p>AWSClusterStatus defines the observed state of AWSCluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>networkStatus</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkStatus">
NetworkStatus
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>failureDomains</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.FailureDomains
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">
Instance
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplate">AWSClusterTemplate
</h3>
<p>
<p>AWSClusterTemplate is the schema for Amazon EC2 based Kubernetes Cluster Templates.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplateSpec">
AWSClusterTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplateResource">
AWSClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplateResource">AWSClusterTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplateSpec">AWSClusterTemplateSpec</a>)
</p>
<p>
<p>AWSClusterTemplateResource defines the desired state of AWSClusterTemplateResource.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.
More info: <a href="https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata">https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata</a></p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">
AWSClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>network</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">
NetworkSpec
</a>
</em>
</td>
<td>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
<p>The AWS Region the cluster lives in.</p>
</td>
</tr>
<tr>
<td>
<code>partition</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Partition is the AWS security partition being used. Defaults to &ldquo;aws&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryControlPlaneLoadBalancer</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">
AWSLoadBalancerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecondaryControlPlaneLoadBalancer is an additional load balancer that can be used for the control plane.</p>
<p>An example use case is to have a separate internal load balancer for internal traffic,
and a separate external load balancer for external traffic.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up machine images when
a machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.
Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
OS and kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupOrg is the AWS Organization ID to look up machine images when a
machine does not specify an AMI. When set, this will be used for all
cluster machines unless a machine specifies a different ImageLookupOrg.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system used to look
up machine images when a machine does not specify an AMI. When set, this
will be used for all cluster machines unless a machine specifies a
different ImageLookupBaseOS.</p>
</td>
</tr>
<tr>
<td>
<code>bastion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Bastion">
Bastion
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Bastion contains options to configure the bastion host.</p>
</td>
</tr>
<tr>
<td>
<code>identityRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">
AWSIdentityReference
</a>
</em>
</td>
<td>
<p>IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
If no identity is specified, the default identity for this controller will be used.</p>
</td>
</tr>
<tr>
<td>
<code>s3Bucket</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.S3Bucket">
S3Bucket
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>S3Bucket contains options to configure a supporting S3 bucket for this
cluster - currently used for nodes requiring Ignition
(<a href="https://coreos.github.io/ignition/">https://coreos.github.io/ignition/</a>) for bootstrapping (requires
BootstrapFormatIgnition feature flag to be enabled).</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplateSpec">AWSClusterTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplate">AWSClusterTemplate</a>)
</p>
<p>
<p>AWSClusterTemplateSpec defines the desired state of AWSClusterTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterTemplateResource">
AWSClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityKind">AWSIdentityKind
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">AWSIdentityReference</a>)
</p>
<p>
<p>AWSIdentityKind defines allowed AWS identity types.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityReference">AWSIdentityReference
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">AWSClusterSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec</a>)
</p>
<p>
<p>AWSIdentityReference specifies a identity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the identity.</p>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSIdentityKind">
AWSIdentityKind
</a>
</em>
</td>
<td>
<p>Kind of the identity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">AWSLoadBalancerSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">AWSClusterSpec</a>)
</p>
<p>
<p>AWSLoadBalancerSpec defines the desired state of an AWS load balancer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Name sets the name of the classic ELB load balancer. As per AWS, the name must be unique
within your set of load balancers for the region, must have a maximum of 32 characters, must
contain only alphanumeric characters or hyphens, and cannot begin or end with a hyphen. Once
set, the value cannot be changed.</p>
</td>
</tr>
<tr>
<td>
<code>scheme</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBScheme">
ELBScheme
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Scheme sets the scheme of the load balancer (defaults to internet-facing)</p>
</td>
</tr>
<tr>
<td>
<code>crossZoneLoadBalancing</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>CrossZoneLoadBalancing enables the classic ELB cross availability zone balancing.</p>
<p>With cross-zone load balancing, each load balancer node for your Classic Load Balancer
distributes requests evenly across the registered instances in all enabled Availability Zones.
If cross-zone load balancing is disabled, each load balancer node distributes requests evenly across
the registered instances in its Availability Zone only.</p>
<p>Defaults to false.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets sets the subnets that should be applied to the control plane load balancer (defaults to discovered subnets for managed VPCs or an empty set for unmanaged VPCs)</p>
</td>
</tr>
<tr>
<td>
<code>healthCheckProtocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBProtocol">
ELBProtocol
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>HealthCheckProtocol sets the protocol type for ELB health check target
default value is ELBProtocolSSL</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups sets the security groups used by the load balancer. Expected to be security group IDs
This is optional - if not provided new security groups will be created for the load balancer</p>
</td>
</tr>
<tr>
<td>
<code>additionalListeners</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AdditionalListenerSpec">
[]AdditionalListenerSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalListeners sets the additional listeners for the control plane load balancer.
This is only applicable to Network Load Balancer (NLB) types for the time being.</p>
</td>
</tr>
<tr>
<td>
<code>ingressRules</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IngressRule">
[]IngressRule
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IngressRules sets the ingress rules for the control plane load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>loadBalancerType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancerType">
LoadBalancerType
</a>
</em>
</td>
<td>
<p>LoadBalancerType sets the type for a load balancer. The default type is classic.</p>
</td>
</tr>
<tr>
<td>
<code>disableHostsRewrite</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableHostsRewrite disabled the hair pinning issue solution that adds the NLB&rsquo;s address as 127.0.0.1 to the hosts
file of each instance. This is by default, false.</p>
</td>
</tr>
<tr>
<td>
<code>preserveClientIP</code><br/>
<em>
bool
</em>
</td>
<td>
<p>PreserveClientIP lets the user control if preservation of client ips must be retained or not.
If this is enabled 6443 will be opened to 0.0.0.0/0.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachine">AWSMachine
</h3>
<p>
<p>AWSMachine is the schema for Amazon EC2 machines.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">
AWSMachineSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceID is the EC2 instance ID for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>instanceMetadataOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">
InstanceMetadataOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceMetadataOptions is the metadata options for the EC2 instance.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
AWSMachine&rsquo;s value takes precedence.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMInstanceProfile is a name of an IAM instance profile to assign to the instance</p>
</td>
</tr>
<tr>
<td>
<code>publicIP</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicIP specifies whether the instance should get a public IP.
Precedence for this setting is as follows:
1. This field if set
2. Cluster/flavor setting
3. Subnet default</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instance. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator. It is possible to specify either IDs of Filters. Using Filters
will cause additional requests to AWS API and if tags change the attached security groups might change too.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnet is a reference to the subnet to use for this instance. If not specified,
the cluster subnet will be used.</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupOverrides</code><br/>
<em>
map[sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.SecurityGroupRole]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecurityGroupOverrides is an optional set of security groups to use for the node.
This is optional - if not provided security groups from the cluster will be used.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces is a list of ENIs to associate with the instance.
A maximum of 2 may be specified.</p>
</td>
</tr>
<tr>
<td>
<code>uncompressedUserData</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UncompressedUserData specify whether the user data is gzip-compressed before it is sent to ec2 instance.
cloud-init has built-in support for gzip-compressed user data
user data stored in aws secret manager is always gzip-compressed.</p>
</td>
</tr>
<tr>
<td>
<code>cloudInit</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.CloudInit">
CloudInit
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Ignition">
Ignition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ignition defined options related to the bootstrapping systems where Ignition is used.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupName specifies the name of the placement group in which to launch the instance.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupPartition</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupPartition is the partition number within the placement group in which to launch the instance.
This value is only valid if the placement group, referred in <code>PlacementGroupName</code>, was created with
strategy set to partition.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
<tr>
<td>
<code>privateDnsName</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.PrivateDNSName">
PrivateDNSName
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrivateDNSName is the options for the instance hostname.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineStatus">
AWSMachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineProviderConditionType">AWSMachineProviderConditionType
(<code>string</code> alias)</p></h3>
<p>
<p>AWSMachineProviderConditionType is a valid value for AWSMachineProviderCondition.Type.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachine">AWSMachine</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateResource">AWSMachineTemplateResource</a>)
</p>
<p>
<p>AWSMachineSpec defines the desired state of an Amazon EC2 instance.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceID is the EC2 instance ID for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>instanceMetadataOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">
InstanceMetadataOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceMetadataOptions is the metadata options for the EC2 instance.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
AWSMachine&rsquo;s value takes precedence.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMInstanceProfile is a name of an IAM instance profile to assign to the instance</p>
</td>
</tr>
<tr>
<td>
<code>publicIP</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicIP specifies whether the instance should get a public IP.
Precedence for this setting is as follows:
1. This field if set
2. Cluster/flavor setting
3. Subnet default</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instance. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator. It is possible to specify either IDs of Filters. Using Filters
will cause additional requests to AWS API and if tags change the attached security groups might change too.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnet is a reference to the subnet to use for this instance. If not specified,
the cluster subnet will be used.</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupOverrides</code><br/>
<em>
map[sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.SecurityGroupRole]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecurityGroupOverrides is an optional set of security groups to use for the node.
This is optional - if not provided security groups from the cluster will be used.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces is a list of ENIs to associate with the instance.
A maximum of 2 may be specified.</p>
</td>
</tr>
<tr>
<td>
<code>uncompressedUserData</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UncompressedUserData specify whether the user data is gzip-compressed before it is sent to ec2 instance.
cloud-init has built-in support for gzip-compressed user data
user data stored in aws secret manager is always gzip-compressed.</p>
</td>
</tr>
<tr>
<td>
<code>cloudInit</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.CloudInit">
CloudInit
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Ignition">
Ignition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ignition defined options related to the bootstrapping systems where Ignition is used.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupName specifies the name of the placement group in which to launch the instance.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupPartition</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupPartition is the partition number within the placement group in which to launch the instance.
This value is only valid if the placement group, referred in <code>PlacementGroupName</code>, was created with
strategy set to partition.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
<tr>
<td>
<code>privateDnsName</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.PrivateDNSName">
PrivateDNSName
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrivateDNSName is the options for the instance hostname.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineStatus">AWSMachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachine">AWSMachine</a>)
</p>
<p>
<p>AWSMachineStatus defines the observed state of AWSMachine.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready is true when the provider resource is ready.</p>
</td>
</tr>
<tr>
<td>
<code>interruptible</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interruptible reports that this machine is using spot instances and can therefore be interrupted by CAPI when it receives a notice that the spot instance is to be terminated by AWS.
This will be set to true when SpotMarketOptions is not nil (i.e. this machine is using a spot instance).</p>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
[]Cluster API api/v1beta1.MachineAddress
</a>
</em>
</td>
<td>
<p>Addresses contains the AWS instance associated addresses.</p>
</td>
</tr>
<tr>
<td>
<code>instanceState</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceState">
InstanceState
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceState is the state of the AWS instance for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the Machine and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the Machine and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the AWSMachine.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplate">AWSMachineTemplate
</h3>
<p>
<p>AWSMachineTemplate is the schema for the Amazon EC2 Machine Templates API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateSpec">
AWSMachineTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateResource">
AWSMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateStatus">
AWSMachineTemplateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateResource">AWSMachineTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateSpec">AWSMachineTemplateSpec</a>)
</p>
<p>
<p>AWSMachineTemplateResource describes the data needed to create am AWSMachine from a template.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.
More info: <a href="https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata">https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata</a></p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">
AWSMachineSpec
</a>
</em>
</td>
<td>
<p>Spec is the specification of the desired behavior of the machine.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceID is the EC2 instance ID for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>instanceMetadataOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">
InstanceMetadataOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceMetadataOptions is the metadata options for the EC2 instance.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
AWSMachine&rsquo;s value takes precedence.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IAMInstanceProfile is a name of an IAM instance profile to assign to the instance</p>
</td>
</tr>
<tr>
<td>
<code>publicIP</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicIP specifies whether the instance should get a public IP.
Precedence for this setting is as follows:
1. This field if set
2. Cluster/flavor setting
3. Subnet default</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instance. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator. It is possible to specify either IDs of Filters. Using Filters
will cause additional requests to AWS API and if tags change the attached security groups might change too.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnet is a reference to the subnet to use for this instance. If not specified,
the cluster subnet will be used.</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupOverrides</code><br/>
<em>
map[sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.SecurityGroupRole]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecurityGroupOverrides is an optional set of security groups to use for the node.
This is optional - if not provided security groups from the cluster will be used.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NetworkInterfaces is a list of ENIs to associate with the instance.
A maximum of 2 may be specified.</p>
</td>
</tr>
<tr>
<td>
<code>uncompressedUserData</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>UncompressedUserData specify whether the user data is gzip-compressed before it is sent to ec2 instance.
cloud-init has built-in support for gzip-compressed user data
user data stored in aws secret manager is always gzip-compressed.</p>
</td>
</tr>
<tr>
<td>
<code>cloudInit</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.CloudInit">
CloudInit
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</td>
</tr>
<tr>
<td>
<code>ignition</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Ignition">
Ignition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ignition defined options related to the bootstrapping systems where Ignition is used.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupName specifies the name of the placement group in which to launch the instance.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupPartition</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupPartition is the partition number within the placement group in which to launch the instance.
This value is only valid if the placement group, referred in <code>PlacementGroupName</code>, was created with
strategy set to partition.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
<tr>
<td>
<code>privateDnsName</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.PrivateDNSName">
PrivateDNSName
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrivateDNSName is the options for the instance hostname.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateSpec">AWSMachineTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplate">AWSMachineTemplate</a>)
</p>
<p>
<p>AWSMachineTemplateSpec defines the desired state of AWSMachineTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateResource">
AWSMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateStatus">AWSMachineTemplateStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplate">AWSMachineTemplate</a>)
</p>
<p>
<p>AWSMachineTemplateStatus defines a status for an AWSMachineTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>capacity</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcelist-v1-core">
Kubernetes core/v1.ResourceList
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Capacity defines the resource capacity for this machine.
This value is used for autoscaling from zero operations as defined in:
<a href="https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md">https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineTemplateWebhook">AWSMachineTemplateWebhook
</h3>
<p>
<p>AWSMachineTemplateWebhook implements a custom validation webhook for AWSMachineTemplate.
Note: we use a custom validator to access the request context for SSA of AWSMachineTemplate.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedCluster">AWSManagedCluster
</h3>
<p>
<p>AWSManagedCluster is the Schema for the awsmanagedclusters API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedClusterSpec">
AWSManagedClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedClusterStatus">
AWSManagedClusterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedClusterSpec">AWSManagedClusterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedCluster">AWSManagedCluster</a>)
</p>
<p>
<p>AWSManagedClusterSpec defines the desired state of AWSManagedCluster</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedClusterStatus">AWSManagedClusterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedCluster">AWSManagedCluster</a>)
</p>
<p>
<p>AWSManagedClusterStatus defines the observed state of AWSManagedCluster</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready is when the AWSManagedControlPlane has a API server URL.</p>
</td>
</tr>
<tr>
<td>
<code>failureDomains</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.FailureDomains
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureDomains specifies a list fo available availability zones that can be used</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">AWSResourceReference
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">AWSLaunchTemplate</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">AWSLaunchTemplate</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec</a>)
</p>
<p>
<p>AWSResourceReference is a reference to a specific AWS resource by ID or filters.
Only one of ID or Filters may be specified. Specifying more than one will result in
a validation error.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ID of resource</p>
</td>
</tr>
<tr>
<td>
<code>filters</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Filter">
[]Filter
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Filters is a set of key/value pairs used to identify a resource
They are applied according to the rules defined by the AWS API:
<a href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html">https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSRoleSpec">AWSRoleSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterRoleIdentitySpec">AWSClusterRoleIdentitySpec</a>)
</p>
<p>
<p>AWSRoleSpec defines the specifications for all identities based around AWS roles.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>roleARN</code><br/>
<em>
string
</em>
</td>
<td>
<p>The Amazon Resource Name (ARN) of the role to assume.</p>
</td>
</tr>
<tr>
<td>
<code>sessionName</code><br/>
<em>
string
</em>
</td>
<td>
<p>An identifier for the assumed role session</p>
</td>
</tr>
<tr>
<td>
<code>durationSeconds</code><br/>
<em>
int32
</em>
</td>
<td>
<p>The duration, in seconds, of the role session before it is renewed.</p>
</td>
</tr>
<tr>
<td>
<code>inlinePolicy</code><br/>
<em>
string
</em>
</td>
<td>
<p>An IAM policy as a JSON-encoded string that you want to use as an inline session policy.</p>
</td>
</tr>
<tr>
<td>
<code>policyARNs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The Amazon Resource Names (ARNs) of the IAM managed policies that you want
to use as managed session policies.
The policies must exist in the same account as the role.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AZSelectionScheme">AZSelectionScheme
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.VPCSpec">VPCSpec</a>)
</p>
<p>
<p>AZSelectionScheme defines the scheme of selecting AZs.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AdditionalListenerSpec">AdditionalListenerSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">AWSLoadBalancerSpec</a>)
</p>
<p>
<p>AdditionalListenerSpec defines the desired state of an
additional listener on an AWS load balancer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>port</code><br/>
<em>
int64
</em>
</td>
<td>
<p>Port sets the port for the additional listener.</p>
</td>
</tr>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBProtocol">
ELBProtocol
</a>
</em>
</td>
<td>
<p>Protocol sets the protocol for the additional listener.
Currently only TCP is supported.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AllowedNamespaces">AllowedNamespaces
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterIdentitySpec">AWSClusterIdentitySpec</a>)
</p>
<p>
<p>AllowedNamespaces is a selector of namespaces that AWSClusters can
use this ClusterPrincipal from. This is a standard Kubernetes LabelSelector,
a label query over a set of resources. The result of matchLabels and
matchExpressions are ANDed.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>list</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>An nil or empty list indicates that AWSClusters cannot use the identity from any namespace.</p>
</td>
</tr>
<tr>
<td>
<code>selector</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>An empty selector indicates that AWSClusters cannot use this
AWSClusterIdentity from any namespace.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Bastion">Bastion
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">AWSClusterSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>Bastion defines a bastion host.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Enabled allows this provider to create a bastion host instance
with a public ip to access the VPC private network.</p>
</td>
</tr>
<tr>
<td>
<code>disableIngressRules</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>DisableIngressRules will ensure there are no Ingress rules in the bastion host&rsquo;s security group.
Requires AllowedCIDRBlocks to be empty.</p>
</td>
</tr>
<tr>
<td>
<code>allowedCIDRBlocks</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AllowedCIDRBlocks is a list of CIDR blocks allowed to access the bastion host.
They are set as ingress rules for the Bastion host&rsquo;s Security Group (defaults to 0.0.0.0/0).</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType will use the specified instance type for the bastion. If not specified,
Cluster API Provider AWS will use t3.micro for all regions except us-east-1, where t2.micro
will be the default.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMI will use the specified AMI to boot the bastion. If not specified,
the AMI will default to one picked out in public space.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.BuildParams">BuildParams
</h3>
<p>
<p>BuildParams is used to build tags around an aws resource.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Lifecycle</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ResourceLifecycle">
ResourceLifecycle
</a>
</em>
</td>
<td>
<p>Lifecycle determines the resource lifecycle.</p>
</td>
</tr>
<tr>
<td>
<code>ClusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ClusterName is the cluster associated with the resource.</p>
</td>
</tr>
<tr>
<td>
<code>ResourceID</code><br/>
<em>
string
</em>
</td>
<td>
<p>ResourceID is the unique identifier of the resource to be tagged.</p>
</td>
</tr>
<tr>
<td>
<code>Name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Name is the name of the resource, it&rsquo;s applied as the tag &ldquo;Name&rdquo; on AWS.</p>
</td>
</tr>
<tr>
<td>
<code>Role</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Role is the role associated to the resource.</p>
</td>
</tr>
<tr>
<td>
<code>Additional</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Any additional tags to be added to the resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.CNIIngressRule">CNIIngressRule
</h3>
<p>
<p>CNIIngressRule defines an AWS ingress rule for CNI requirements.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>description</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroupProtocol">
SecurityGroupProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>fromPort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>toPort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.CNIIngressRules">CNIIngressRules
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.CNIIngressRule</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.CNISpec">CNISpec</a>)
</p>
<p>
<p>CNIIngressRules is a slice of CNIIngressRule.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.CNISpec">CNISpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">NetworkSpec</a>)
</p>
<p>
<p>CNISpec defines configuration for CNI.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cniIngressRules</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.CNIIngressRules">
CNIIngressRules
</a>
</em>
</td>
<td>
<p>CNIIngressRules specify rules to apply to control plane and worker node security groups.
The source for the rule will be set to control plane and worker security group IDs.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ClassicELBAttributes">ClassicELBAttributes
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">LoadBalancer</a>)
</p>
<p>
<p>ClassicELBAttributes defines extra attributes associated with a classic load balancer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>idleTimeout</code><br/>
<em>
<a href="https://golang.org/pkg/time/#Duration">
time.Duration
</a>
</em>
</td>
<td>
<p>IdleTimeout is time that the connection is allowed to be idle (no data
has been sent over the connection) before it is closed by the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>crossZoneLoadBalancing</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>CrossZoneLoadBalancing enables the classic load balancer load balancing.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ClassicELBHealthCheck">ClassicELBHealthCheck
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">LoadBalancer</a>)
</p>
<p>
<p>ClassicELBHealthCheck defines an AWS classic load balancer health check.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>target</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>interval</code><br/>
<em>
<a href="https://golang.org/pkg/time/#Duration">
time.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>timeout</code><br/>
<em>
<a href="https://golang.org/pkg/time/#Duration">
time.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>healthyThreshold</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>unhealthyThreshold</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ClassicELBListener">ClassicELBListener
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">LoadBalancer</a>)
</p>
<p>
<p>ClassicELBListener defines an AWS classic load balancer listener.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBProtocol">
ELBProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>port</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instanceProtocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBProtocol">
ELBProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instancePort</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.CloudInit">CloudInit
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>)
</p>
<p>
<p>CloudInit defines options related to the bootstrapping systems where
CloudInit is used.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>insecureSkipSecretsManager</code><br/>
<em>
bool
</em>
</td>
<td>
<p>InsecureSkipSecretsManager, when set to true will not use AWS Secrets Manager
or AWS Systems Manager Parameter Store to ensure privacy of userdata.
By default, a cloud-init boothook shell script is prepended to download
the userdata from Secrets Manager and additionally delete the secret.</p>
</td>
</tr>
<tr>
<td>
<code>secretCount</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretCount is the number of secrets used to form the complete secret</p>
</td>
</tr>
<tr>
<td>
<code>secretPrefix</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretPrefix is the prefix for the secret name. This is stored
temporarily, and deleted when the machine registers as a node against
the workload cluster.</p>
</td>
</tr>
<tr>
<td>
<code>secureSecretsBackend</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecretBackend">
SecretBackend
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecureSecretsBackend, when set to parameter-store will utilize the AWS Systems Manager
Parameter Storage to distribute secrets. By default or with the value of secrets-manager,
will use AWS Secrets Manager instead.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.EKSAMILookupType">EKSAMILookupType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AMIReference">AMIReference</a>)
</p>
<p>
<p>EKSAMILookupType specifies which AWS AMI to use for a AWSMachine and AWSMachinePool.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ELBProtocol">ELBProtocol
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">AWSLoadBalancerSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AdditionalListenerSpec">AdditionalListenerSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.ClassicELBListener">ClassicELBListener</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.Listener">Listener</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.TargetGroupSpec">TargetGroupSpec</a>)
</p>
<p>
<p>ELBProtocol defines listener protocols for a load balancer.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ELBScheme">ELBScheme
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">AWSLoadBalancerSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">LoadBalancer</a>)
</p>
<p>
<p>ELBScheme defines the scheme of a load balancer.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Filter">Filter
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">AWSResourceReference</a>)
</p>
<p>
<p>Filter is a filter used to identify an AWS resource.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the filter. Filter names are case-sensitive.</p>
</td>
</tr>
<tr>
<td>
<code>values</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Values includes one or more filter values. Filter values are case-sensitive.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.GCTask">GCTask
(<code>string</code> alias)</p></h3>
<p>
<p>GCTask defines a task to be executed by the garbage collector.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.HTTPTokensState">HTTPTokensState
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">InstanceMetadataOptions</a>)
</p>
<p>
<p>HTTPTokensState describes the state of InstanceMetadataOptions.HTTPTokensState</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IPAMPool">IPAMPool
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.IPv6">IPv6</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.VPCSpec">VPCSpec</a>)
</p>
<p>
<p>IPAMPool defines the IPAM pool to be used for VPC.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is the ID of the IPAM pool this provider should use to create VPC.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the IPAM pool this provider should use to create VPC.</p>
</td>
</tr>
<tr>
<td>
<code>netmaskLength</code><br/>
<em>
int64
</em>
</td>
<td>
<p>The netmask length of the IPv4 CIDR you want to allocate to VPC from
an Amazon VPC IP Address Manager (IPAM) pool.
Defaults to /16 for IPv4 if not specified.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IPv6">IPv6
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.VPCSpec">VPCSpec</a>)
</p>
<p>
<p>IPv6 contains ipv6 specific settings for the network.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>CidrBlock is the CIDR block provided by Amazon when VPC has enabled IPv6.
Mutually exclusive with IPAMPool.</p>
</td>
</tr>
<tr>
<td>
<code>poolId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PoolID is the IP pool which must be defined in case of BYO IP is defined.
Must be specified if CidrBlock is set.
Mutually exclusive with IPAMPool.</p>
</td>
</tr>
<tr>
<td>
<code>egressOnlyInternetGatewayId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EgressOnlyInternetGatewayID is the id of the egress only internet gateway associated with an IPv6 enabled VPC.</p>
</td>
</tr>
<tr>
<td>
<code>ipamPool</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IPAMPool">
IPAMPool
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPAMPool defines the IPAMv6 pool to be used for VPC.
Mutually exclusive with CidrBlock.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Ignition">Ignition
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>)
</p>
<p>
<p>Ignition defines options related to the bootstrapping systems where Ignition is used.
For more information on Ignition configuration, see <a href="https://coreos.github.io/butane/specs/">https://coreos.github.io/butane/specs/</a></p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines which version of Ignition will be used to generate bootstrap data.</p>
</td>
</tr>
<tr>
<td>
<code>storageType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IgnitionStorageTypeOption">
IgnitionStorageTypeOption
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>StorageType defines how to store the boostrap user data for Ignition.
This can be used to instruct Ignition from where to fetch the user data to bootstrap an instance.</p>
<p>When omitted, the storage option will default to ClusterObjectStore.</p>
<p>When set to &ldquo;ClusterObjectStore&rdquo;, if the capability is available and a Cluster ObjectStore configuration
is correctly provided in the Cluster object (under .spec.s3Bucket),
an object store will be used to store bootstrap user data.</p>
<p>When set to &ldquo;UnencryptedUserData&rdquo;, EC2 Instance User Data will be used to store the machine bootstrap user data, unencrypted.
This option is considered less secure than others as user data may contain sensitive informations (keys, certificates, etc.)
and users with ec2:DescribeInstances permission or users running pods
that can access the ec2 metadata service have access to this sensitive information.
So this is only to be used at ones own risk, and only when other more secure options are not viable.</p>
</td>
</tr>
<tr>
<td>
<code>proxy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IgnitionProxy">
IgnitionProxy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Proxy defines proxy settings for Ignition.
Only valid for Ignition versions 3.1 and above.</p>
</td>
</tr>
<tr>
<td>
<code>tls</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IgnitionTLS">
IgnitionTLS
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>TLS defines TLS settings for Ignition.
Only valid for Ignition versions 3.1 and above.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IgnitionCASource">IgnitionCASource
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.IgnitionTLS">IgnitionTLS</a>)
</p>
<p>
<p>IgnitionCASource defines the source of the certificate authority to use for Ignition.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IgnitionNoProxy">IgnitionNoProxy
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.IgnitionProxy">IgnitionProxy</a>)
</p>
<p>
<p>IgnitionNoProxy defines the list of domains to not proxy for Ignition.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IgnitionProxy">IgnitionProxy
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.Ignition">Ignition</a>)
</p>
<p>
<p>IgnitionProxy defines proxy settings for Ignition.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>httpProxy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>HTTPProxy is the HTTP proxy to use for Ignition.
A single URL that specifies the proxy server to use for HTTP and HTTPS requests,
unless overridden by the HTTPSProxy or NoProxy options.</p>
</td>
</tr>
<tr>
<td>
<code>httpsProxy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>HTTPSProxy is the HTTPS proxy to use for Ignition.
A single URL that specifies the proxy server to use for HTTPS requests,
unless overridden by the NoProxy option.</p>
</td>
</tr>
<tr>
<td>
<code>noProxy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IgnitionNoProxy">
[]IgnitionNoProxy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NoProxy is the list of domains to not proxy for Ignition.
Specifies a list of strings to hosts that should be excluded from proxying.</p>
<p>Each value is represented by:
- An IP address prefix (1.2.3.4)
- An IP address prefix in CIDR notation (1.2.3.<sup>4</sup>&frasl;<sub>8</sub>)
- A domain name
- A domain name matches that name and all subdomains
- A domain name with a leading . matches subdomains only
- A special DNS label (*), indicates that no proxying should be done</p>
<p>An IP address prefix and domain name can also include a literal port number (1.2.3.4:80).</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IgnitionStorageTypeOption">IgnitionStorageTypeOption
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.Ignition">Ignition</a>)
</p>
<p>
<p>IgnitionStorageTypeOption defines the different storage types for Ignition.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IgnitionTLS">IgnitionTLS
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.Ignition">Ignition</a>)
</p>
<p>
<p>IgnitionTLS defines TLS settings for Ignition.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>certificateAuthorities</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IgnitionCASource">
[]IgnitionCASource
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CASources defines the list of certificate authorities to use for Ignition.
The value is the certificate bundle (in PEM format). The bundle can contain multiple concatenated certificates.
Supported schemes are http, https, tftp, s3, arn, gs, and <code>data</code> (RFC 2397) URL scheme.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IngressRule">IngressRule
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">AWSLoadBalancerSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">NetworkSpec</a>)
</p>
<p>
<p>IngressRule defines an AWS ingress rule for security groups.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>description</code><br/>
<em>
string
</em>
</td>
<td>
<p>Description provides extended information about the ingress rule.</p>
</td>
</tr>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroupProtocol">
SecurityGroupProtocol
</a>
</em>
</td>
<td>
<p>Protocol is the protocol for the ingress rule. Accepted values are &ldquo;-1&rdquo; (all), &ldquo;4&rdquo; (IP in IP),&ldquo;tcp&rdquo;, &ldquo;udp&rdquo;, &ldquo;icmp&rdquo;, and &ldquo;58&rdquo; (ICMPv6), &ldquo;50&rdquo; (ESP).</p>
</td>
</tr>
<tr>
<td>
<code>fromPort</code><br/>
<em>
int64
</em>
</td>
<td>
<p>FromPort is the start of port range.</p>
</td>
</tr>
<tr>
<td>
<code>toPort</code><br/>
<em>
int64
</em>
</td>
<td>
<p>ToPort is the end of port range.</p>
</td>
</tr>
<tr>
<td>
<code>cidrBlocks</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>List of CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.</p>
</td>
</tr>
<tr>
<td>
<code>ipv6CidrBlocks</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>List of IPv6 CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.</p>
</td>
</tr>
<tr>
<td>
<code>sourceSecurityGroupIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The security group id to allow access from. Cannot be specified with CidrBlocks.</p>
</td>
</tr>
<tr>
<td>
<code>sourceSecurityGroupRoles</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroupRole">
[]SecurityGroupRole
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>The security group role to allow access from. Cannot be specified with CidrBlocks.
The field will be combined with source security group IDs if specified.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.IngressRules">IngressRules
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.IngressRule</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroup">SecurityGroup</a>)
</p>
<p>
<p>IngressRules is a slice of AWS ingress rules for security groups.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Instance">Instance
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStatus">AWSClusterStatus</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AutoScalingGroup">AutoScalingGroup</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AutoScalingGroup">AutoScalingGroup</a>)
</p>
<p>
<p>Instance describes an AWS instance.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instanceState</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceState">
InstanceState
</a>
</em>
</td>
<td>
<p>The current state of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
<p>The instance type.</p>
</td>
</tr>
<tr>
<td>
<code>subnetId</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the subnet of the instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageId</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the AMI used to launch the instance.</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the SSH key pair.</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SecurityGroupIDs are one or more security group IDs this instance belongs to.</p>
</td>
</tr>
<tr>
<td>
<code>userData</code><br/>
<em>
string
</em>
</td>
<td>
<p>UserData is the raw data script passed to the instance which is run upon bootstrap.
This field must not be base64 encoded and should only be used when running a new instance.</p>
</td>
</tr>
<tr>
<td>
<code>iamProfile</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the IAM instance profile associated with the instance, if applicable.</p>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
[]Cluster API api/v1beta1.MachineAddress
</a>
</em>
</td>
<td>
<p>Addresses contains the AWS instance associated addresses.</p>
</td>
</tr>
<tr>
<td>
<code>privateIp</code><br/>
<em>
string
</em>
</td>
<td>
<p>The private IPv4 address assigned to the instance.</p>
</td>
</tr>
<tr>
<td>
<code>publicIp</code><br/>
<em>
string
</em>
</td>
<td>
<p>The public IPv4 address assigned to the instance, if applicable.</p>
</td>
</tr>
<tr>
<td>
<code>enaSupport</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Specifies whether enhanced networking with ENA is enabled.</p>
</td>
</tr>
<tr>
<td>
<code>ebsOptimized</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Indicates whether the instance is optimized for Amazon EBS I/O.</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the root storage volume.</p>
</td>
</tr>
<tr>
<td>
<code>nonRootVolumes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
[]Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Configuration options for the non root storage volumes.</p>
</td>
</tr>
<tr>
<td>
<code>networkInterfaces</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Specifies ENIs attached to instance</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>The tags associated with the instance.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZone</code><br/>
<em>
string
</em>
</td>
<td>
<p>Availability zone of instance</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<p>SpotMarketOptions option for configuring instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupName specifies the name of the placement group in which to launch the instance.</p>
</td>
</tr>
<tr>
<td>
<code>placementGroupPartition</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>PlacementGroupPartition is the partition number within the placement group in which to launch the instance.
This value is only valid if the placement group, referred in <code>PlacementGroupName</code>, was created with
strategy set to partition.</p>
</td>
</tr>
<tr>
<td>
<code>tenancy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tenancy indicates if instance should run on shared or single-tenant hardware.</p>
</td>
</tr>
<tr>
<td>
<code>volumeIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IDs of the instance&rsquo;s volumes</p>
</td>
</tr>
<tr>
<td>
<code>instanceMetadataOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">
InstanceMetadataOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceMetadataOptions is the metadata options for the EC2 instance.</p>
</td>
</tr>
<tr>
<td>
<code>privateDnsName</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.PrivateDNSName">
PrivateDNSName
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrivateDNSName is the options for the instance hostname.</p>
</td>
</tr>
<tr>
<td>
<code>publicIPOnLaunch</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>PublicIPOnLaunch is the option to associate a public IP on instance launch</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">InstanceMetadataOptions
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">Instance</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">AWSLaunchTemplate</a>)
</p>
<p>
<p>InstanceMetadataOptions describes metadata options for the EC2 instance.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>httpEndpoint</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataState">
InstanceMetadataState
</a>
</em>
</td>
<td>
<p>Enables or disables the HTTP metadata endpoint on your instances.</p>
<p>If you specify a value of disabled, you cannot access your instance metadata.</p>
<p>Default: enabled</p>
</td>
</tr>
<tr>
<td>
<code>httpPutResponseHopLimit</code><br/>
<em>
int64
</em>
</td>
<td>
<p>The desired HTTP PUT response hop limit for instance metadata requests. The
larger the number, the further instance metadata requests can travel.</p>
<p>Default: 1</p>
</td>
</tr>
<tr>
<td>
<code>httpTokens</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.HTTPTokensState">
HTTPTokensState
</a>
</em>
</td>
<td>
<p>The state of token usage for your instance metadata requests.</p>
<p>If the state is optional, you can choose to retrieve instance metadata with
or without a session token on your request. If you retrieve the IAM role
credentials without a token, the version 1.0 role credentials are returned.
If you retrieve the IAM role credentials using a valid session token, the
version 2.0 role credentials are returned.</p>
<p>If the state is required, you must send a session token with any instance
metadata retrieval requests. In this state, retrieving the IAM role credentials
always returns the version 2.0 credentials; the version 1.0 credentials are
not available.</p>
<p>Default: optional</p>
</td>
</tr>
<tr>
<td>
<code>instanceMetadataTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataState">
InstanceMetadataState
</a>
</em>
</td>
<td>
<p>Set to enabled to allow access to instance tags from the instance metadata.
Set to disabled to turn off access to instance tags from the instance metadata.
For more information, see Work with instance tags using the instance metadata
(<a href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#work-with-tags-in-IMDS">https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#work-with-tags-in-IMDS</a>).</p>
<p>Default: disabled</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataState">InstanceMetadataState
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">InstanceMetadataOptions</a>)
</p>
<p>
<p>InstanceMetadataState describes the state of InstanceMetadataOptions.HttpEndpoint and InstanceMetadataOptions.InstanceMetadataTags</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.InstanceState">InstanceState
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineStatus">AWSMachineStatus</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">Instance</a>)
</p>
<p>
<p>InstanceState describes the state of an AWS instance.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Listener">Listener
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">LoadBalancer</a>)
</p>
<p>
<p>Listener defines an AWS network load balancer listener.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBProtocol">
ELBProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>port</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetGroup</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.TargetGroupSpec">
TargetGroupSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">LoadBalancer
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkStatus">NetworkStatus</a>)
</p>
<p>
<p>LoadBalancer defines an AWS load balancer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code><br/>
<em>
string
</em>
</td>
<td>
<p>ARN of the load balancer. Unlike the ClassicLB, ARN is used mostly
to define and get it.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The name of the load balancer. It must be unique within the set of load balancers
defined in the region. It also serves as identifier.</p>
</td>
</tr>
<tr>
<td>
<code>dnsName</code><br/>
<em>
string
</em>
</td>
<td>
<p>DNSName is the dns name of the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>scheme</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBScheme">
ELBScheme
</a>
</em>
</td>
<td>
<p>Scheme is the load balancer scheme, either internet-facing or private.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones in the VPC attached to the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SubnetIDs is an array of subnets in the VPC attached to the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupIds</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SecurityGroupIDs is an array of security groups assigned to the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>listeners</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ClassicELBListener">
[]ClassicELBListener
</a>
</em>
</td>
<td>
<p>ClassicELBListeners is an array of classic elb listeners associated with the load balancer. There must be at least one.</p>
</td>
</tr>
<tr>
<td>
<code>healthChecks</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ClassicELBHealthCheck">
ClassicELBHealthCheck
</a>
</em>
</td>
<td>
<p>HealthCheck is the classic elb health check associated with the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>attributes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ClassicELBAttributes">
ClassicELBAttributes
</a>
</em>
</td>
<td>
<p>ClassicElbAttributes defines extra attributes associated with the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Tags is a map of tags associated with the load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>elbListeners</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Listener">
[]Listener
</a>
</em>
</td>
<td>
<p>ELBListeners is an array of listeners associated with the load balancer. There must be at least one.</p>
</td>
</tr>
<tr>
<td>
<code>elbAttributes</code><br/>
<em>
map[string]*string
</em>
</td>
<td>
<p>ELBAttributes defines extra attributes associated with v2 load balancers.</p>
</td>
</tr>
<tr>
<td>
<code>loadBalancerType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancerType">
LoadBalancerType
</a>
</em>
</td>
<td>
<p>LoadBalancerType sets the type for a load balancer. The default type is classic.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancerAttribute">LoadBalancerAttribute
(<code>string</code> alias)</p></h3>
<p>
<p>LoadBalancerAttribute defines a set of attributes for a V2 load balancer.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancerType">LoadBalancerType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLoadBalancerSpec">AWSLoadBalancerSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">LoadBalancer</a>)
</p>
<p>
<p>LoadBalancerType defines the type of load balancer to use.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">NetworkSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">AWSClusterSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>)
</p>
<p>
<p>NetworkSpec encapsulates all things related to AWS network.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>vpc</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.VPCSpec">
VPCSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>VPC configuration.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Subnets">
Subnets
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets configuration.</p>
</td>
</tr>
<tr>
<td>
<code>cni</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.CNISpec">
CNISpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CNI configuration</p>
</td>
</tr>
<tr>
<td>
<code>securityGroupOverrides</code><br/>
<em>
map[sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.SecurityGroupRole]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecurityGroupOverrides is an optional set of security groups to use for cluster instances
This is optional - if not provided new security groups will be created for the cluster</p>
</td>
</tr>
<tr>
<td>
<code>additionalControlPlaneIngressRules</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IngressRule">
[]IngressRule
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalControlPlaneIngressRules is an optional set of ingress rules to add to the control plane</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.NetworkStatus">NetworkStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterStatus">AWSClusterStatus</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneStatus">AWSManagedControlPlaneStatus</a>)
</p>
<p>
<p>NetworkStatus encapsulates AWS networking resources.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>securityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroup">
map[sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.SecurityGroupRole]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.SecurityGroup
</a>
</em>
</td>
<td>
<p>SecurityGroups is a map from the role/kind of the security group to its unique name, if any.</p>
</td>
</tr>
<tr>
<td>
<code>apiServerElb</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">
LoadBalancer
</a>
</em>
</td>
<td>
<p>APIServerELB is the Kubernetes api server load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>secondaryAPIServerELB</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.LoadBalancer">
LoadBalancer
</a>
</em>
</td>
<td>
<p>SecondaryAPIServerELB is the secondary Kubernetes api server load balancer.</p>
</td>
</tr>
<tr>
<td>
<code>natGatewaysIPs</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>NatGatewaysIPs contains the public IPs of the NAT Gateways</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.PrivateDNSName">PrivateDNSName
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">Instance</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">AWSLaunchTemplate</a>)
</p>
<p>
<p>PrivateDNSName is the options for the instance hostname.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enableResourceNameDnsAAAARecord</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>EnableResourceNameDNSAAAARecord indicates whether to respond to DNS queries for instance hostnames with DNS AAAA records.</p>
</td>
</tr>
<tr>
<td>
<code>enableResourceNameDnsARecord</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>EnableResourceNameDNSARecord indicates whether to respond to DNS queries for instance hostnames with DNS A records.</p>
</td>
</tr>
<tr>
<td>
<code>hostnameType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The type of hostname to assign to an instance.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ResourceLifecycle">ResourceLifecycle
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.BuildParams">BuildParams</a>)
</p>
<p>
<p>ResourceLifecycle configures the lifecycle of a resource.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.RouteTable">RouteTable
</h3>
<p>
<p>RouteTable defines an AWS routing table.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.S3Bucket">S3Bucket
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">AWSClusterSpec</a>)
</p>
<p>
<p>S3Bucket defines a supporting S3 bucket for the cluster, currently can be optionally used for Ignition.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>controlPlaneIAMInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneIAMInstanceProfile is a name of the IAMInstanceProfile, which will be allowed
to read control-plane node bootstrap data from S3 Bucket.</p>
</td>
</tr>
<tr>
<td>
<code>nodesIAMInstanceProfiles</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NodesIAMInstanceProfiles is a list of IAM instance profiles, which will be allowed to read
worker nodes bootstrap data from S3 Bucket.</p>
</td>
</tr>
<tr>
<td>
<code>presignedURLDuration</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PresignedURLDuration defines the duration for which presigned URLs are valid.</p>
<p>This is used to generate presigned URLs for S3 Bucket objects, which are used by
control-plane and worker nodes to fetch bootstrap data.</p>
<p>When enabled, the IAM instance profiles specified are not used.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name defines name of S3 Bucket to be created.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SecretBackend">SecretBackend
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.CloudInit">CloudInit</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMConfigurationSpec">AWSIAMConfigurationSpec</a>)
</p>
<p>
<p>SecretBackend defines variants for backend secret storage.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroup">SecurityGroup
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkStatus">NetworkStatus</a>)
</p>
<p>
<p>SecurityGroup defines an AWS security group.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is a unique identifier.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name is the security group name.</p>
</td>
</tr>
<tr>
<td>
<code>ingressRule</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IngressRules">
IngressRules
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IngressRules is the inbound rules associated with the security group.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a map of tags associated with the security group.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroupProtocol">SecurityGroupProtocol
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.CNIIngressRule">CNIIngressRule</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.IngressRule">IngressRule</a>)
</p>
<p>
<p>SecurityGroupProtocol defines the protocol type for a security group rule.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroupRole">SecurityGroupRole
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.IngressRule">IngressRule</a>)
</p>
<p>
<p>SecurityGroupRole defines the unique role of a security group.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SpotMarketOptions">SpotMarketOptions
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">Instance</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">AWSLaunchTemplate</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">AWSLaunchTemplate</a>)
</p>
<p>
<p>SpotMarketOptions defines the options available to a user when configuring
Machines to run on Spot instances.
Most users should provide an empty struct.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxPrice</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxPrice defines the maximum price the user is willing to pay for Spot VM instances</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SubnetSpec">SubnetSpec
</h3>
<p>
<p>SubnetSpec configures an AWS Subnet.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID defines a unique identifier to reference this resource.
If you&rsquo;re bringing your subnet, set the AWS subnet-id here, it must start with <code>subnet-</code>.</p>
<p>When the VPC is managed by CAPA, and you&rsquo;d like the provider to create a subnet for you,
the id can be set to any placeholder value that does not start with <code>subnet-</code>;
upon creation, the subnet AWS identifier will be populated in the <code>ResourceID</code> field and
the <code>id</code> field is going to be used as the subnet name. If you specify a tag
called <code>Name</code>, it takes precedence.</p>
</td>
</tr>
<tr>
<td>
<code>resourceID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ResourceID is the subnet identifier from AWS, READ ONLY.
This field is populated when the provider manages the subnet.</p>
</td>
</tr>
<tr>
<td>
<code>cidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<p>CidrBlock is the CIDR block to be used when the provider creates a managed VPC.</p>
</td>
</tr>
<tr>
<td>
<code>ipv6CidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPv6CidrBlock is the IPv6 CIDR block to be used when the provider creates a managed VPC.
A subnet can have an IPv4 and an IPv6 address.
IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZone</code><br/>
<em>
string
</em>
</td>
<td>
<p>AvailabilityZone defines the availability zone to use for this subnet in the cluster&rsquo;s region.</p>
</td>
</tr>
<tr>
<td>
<code>isPublic</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>IsPublic defines the subnet as a public subnet. A subnet is public when it is associated with a route table that has a route to an internet gateway.</p>
</td>
</tr>
<tr>
<td>
<code>isIpv6</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>IsIPv6 defines the subnet as an IPv6 subnet. A subnet is IPv6 when it is associated with a VPC that has IPv6 enabled.
IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.</p>
</td>
</tr>
<tr>
<td>
<code>routeTableId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RouteTableID is the routing table id associated with the subnet.</p>
</td>
</tr>
<tr>
<td>
<code>natGatewayId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NatGatewayID is the NAT gateway id associated with the subnet.
Ignored unless the subnet is managed by the provider, in which case this is set on the public subnet where the NAT gateway resides. It is then used to determine routes for private subnets in the same AZ as the public subnet.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a collection of tags describing the resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Subnets">Subnets
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2.SubnetSpec</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">NetworkSpec</a>)
</p>
<p>
<p>Subnets is a slice of Subnet.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Tags">Tags
(<code>map[string]string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSClusterSpec">AWSClusterSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.BuildParams">BuildParams</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.SecurityGroup">SecurityGroup</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.SubnetSpec">SubnetSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.VPCSpec">VPCSpec</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.AWSIAMRoleSpec">AWSIAMRoleSpec</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1alpha1.BootstrapUser">BootstrapUser</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.AWSIAMRoleSpec">AWSIAMRoleSpec</a>, <a href="#bootstrap.aws.infrastructure.cluster.x-k8s.io/v1beta1.BootstrapUser">BootstrapUser</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta1.OIDCIdentityProviderConfig">OIDCIdentityProviderConfig</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.AWSManagedControlPlaneSpec">AWSManagedControlPlaneSpec</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.OIDCIdentityProviderConfig">OIDCIdentityProviderConfig</a>, <a href="#controlplane.cluster.x-k8s.io/v1beta2.RosaControlPlaneSpec">RosaControlPlaneSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AutoScalingGroup">AutoScalingGroup</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.FargateProfileSpec">FargateProfileSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AutoScalingGroup">AutoScalingGroup</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.FargateProfileSpec">FargateProfileSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolSpec">RosaMachinePoolSpec</a>)
</p>
<p>
<p>Tags defines a map of tags.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.TargetGroupAttribute">TargetGroupAttribute
(<code>string</code> alias)</p></h3>
<p>
<p>TargetGroupAttribute defines attribute key values for V2 Load Balancer Attributes.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.TargetGroupHealthCheck">TargetGroupHealthCheck
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.TargetGroupSpec">TargetGroupSpec</a>)
</p>
<p>
<p>TargetGroupHealthCheck defines health check settings for the target group.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>protocol</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>path</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>port</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>intervalSeconds</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>timeoutSeconds</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>thresholdCount</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.TargetGroupSpec">TargetGroupSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.Listener">Listener</a>)
</p>
<p>
<p>TargetGroupSpec specifies target group settings for a given listener.
This is created first, and the ARN is then passed to the listener.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of the TargetGroup. Must be unique over the same group of listeners.</p>
</td>
</tr>
<tr>
<td>
<code>port</code><br/>
<em>
int64
</em>
</td>
<td>
<p>Port is the exposed port</p>
</td>
</tr>
<tr>
<td>
<code>protocol</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ELBProtocol">
ELBProtocol
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>vpcId</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>targetGroupHealthCheck</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.TargetGroupHealthCheck">
TargetGroupHealthCheck
</a>
</em>
</td>
<td>
<p>HealthCheck is the elb health check associated with the load balancer.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.VPCSpec">VPCSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.NetworkSpec">NetworkSpec</a>)
</p>
<p>
<p>VPCSpec configures an AWS VPC.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is the vpc-id of the VPC this provider should use to create resources.</p>
</td>
</tr>
<tr>
<td>
<code>cidrBlock</code><br/>
<em>
string
</em>
</td>
<td>
<p>CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
Defaults to 10.0.0.0/16.
Mutually exclusive with IPAMPool.</p>
</td>
</tr>
<tr>
<td>
<code>ipamPool</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IPAMPool">
IPAMPool
</a>
</em>
</td>
<td>
<p>IPAMPool defines the IPAMv4 pool to be used for VPC.
Mutually exclusive with CidrBlock.</p>
</td>
</tr>
<tr>
<td>
<code>ipv6</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.IPv6">
IPv6
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPv6 contains ipv6 specific settings for the network. Supported only in managed clusters.
This field cannot be set on AWSCluster object.</p>
</td>
</tr>
<tr>
<td>
<code>internetGatewayId</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InternetGatewayID is the id of the internet gateway associated with the VPC.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<p>Tags is a collection of tags describing the resource.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneUsageLimit</code><br/>
<em>
int
</em>
</td>
<td>
<p>AvailabilityZoneUsageLimit specifies the maximum number of availability zones (AZ) that
should be used in a region when automatically creating subnets. If a region has more
than this number of AZs then this number of AZs will be picked randomly when creating
default subnets. Defaults to 3</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneSelection</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AZSelectionScheme">
AZSelectionScheme
</a>
</em>
</td>
<td>
<p>AvailabilityZoneSelection specifies how AZs should be selected if there are more AZs
in a region than specified by AvailabilityZoneUsageLimit. There are 2 selection schemes:
Ordered - selects based on alphabetical order
Random - selects AZs randomly in a region
Defaults to Ordered</p>
</td>
</tr>
<tr>
<td>
<code>emptyRoutesDefaultVPCSecurityGroup</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>EmptyRoutesDefaultVPCSecurityGroup specifies whether the default VPC security group ingress
and egress rules should be removed.</p>
<p>By default, when creating a VPC, AWS creates a security group called <code>default</code> with ingress and egress
rules that allow traffic from anywhere. The group could be used as a potential surface attack and
it&rsquo;s generally suggested that the group rules are removed or modified appropriately.</p>
<p>NOTE: This only applies when the VPC is managed by the Cluster API AWS controller.</p>
</td>
</tr>
<tr>
<td>
<code>privateDnsHostnameTypeOnLaunch</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrivateDNSHostnameTypeOnLaunch is the type of hostname to assign to instances in the subnet at launch.
For IPv4-only and dual-stack (IPv4 and IPv6) subnets, an instance DNS name can be based on the instance IPv4 address (ip-name)
or the instance ID (resource-name). For IPv6 only subnets, an instance DNS name must be based on the instance ID (resource-name).</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Volume">Volume
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachineSpec">AWSMachineSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">Instance</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta1.AWSLaunchTemplate">AWSLaunchTemplate</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">AWSLaunchTemplate</a>)
</p>
<p>
<p>Volume encapsulates the configuration options for the storage device.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>deviceName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Device name</p>
</td>
</tr>
<tr>
<td>
<code>size</code><br/>
<em>
int64
</em>
</td>
<td>
<p>Size specifies size (in Gi) of the storage device.
Must be greater than the image snapshot size or 8 (whichever is greater).</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.VolumeType">
VolumeType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Type is the type of the volume (e.g. gp2, io1, etc&hellip;).</p>
</td>
</tr>
<tr>
<td>
<code>iops</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>IOPS is the number of IOPS requested for the disk. Not applicable to all types.</p>
</td>
</tr>
<tr>
<td>
<code>throughput</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>Throughput to provision in MiB/s supported for the volume type. Not applicable to all types.</p>
</td>
</tr>
<tr>
<td>
<code>encrypted</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Encrypted is whether the volume should be encrypted or not.</p>
</td>
</tr>
<tr>
<td>
<code>encryptionKey</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EncryptionKey is the KMS key to use to encrypt the volume. Can be either a KMS key ID or ARN.
If Encrypted is set and this is omitted, the default AWS key will be used.
The key must already exist and be accessible by the controller.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.VolumeType">VolumeType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">Volume</a>)
</p>
<p>
<p>VolumeType describes the EBS volume type.
See: <a href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html">https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html</a></p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ASGStatus">ASGStatus
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolStatus">AWSMachinePoolStatus</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AutoScalingGroup">AutoScalingGroup</a>)
</p>
<p>
<p>ASGStatus is a status string returned by the autoscaling API.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSFargateProfile">AWSFargateProfile
</h3>
<p>
<p>AWSFargateProfile is the Schema for the awsfargateprofiles API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.FargateProfileSpec">
FargateProfileSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>clusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ClusterName is the name of the Cluster this object belongs to.</p>
</td>
</tr>
<tr>
<td>
<code>profileName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProfileName specifies the profile name.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for this fargate pool
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>selectors</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.FargateSelector">
[]FargateSelector
</a>
</em>
</td>
<td>
<p>Selectors specify fargate pod selectors.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.FargateProfileStatus">
FargateProfileStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">AWSLaunchTemplate
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>AWSLaunchTemplate defines the desired state of AWSLaunchTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name of the launch template.</p>
</td>
</tr>
<tr>
<td>
<code>iamInstanceProfile</code><br/>
<em>
string
</em>
</td>
<td>
<p>The name or the Amazon Resource Name (ARN) of the instance profile associated
with the IAM role for the instance. The instance profile contains the IAM
role.</p>
</td>
</tr>
<tr>
<td>
<code>ami</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AMIReference">
AMIReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMI is the reference to the AMI from which to create the machine instance.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupFormat</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImageLookupFormat is the AMI naming format to look up the image for this
machine It will be ignored if an explicit AMI is set. Supports
substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
kubernetes version, respectively. The BaseOS will be the value in
ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
defined by the packages produced by kubernetes/release without v as a
prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
also: <a href="https://golang.org/pkg/text/template/">https://golang.org/pkg/text/template/</a></p>
</td>
</tr>
<tr>
<td>
<code>imageLookupOrg</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>imageLookupBaseOS</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageLookupBaseOS is the name of the base operating system to use for
image lookup the AMI is not set.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType is the type of instance to create. Example: m4.xlarge</p>
</td>
</tr>
<tr>
<td>
<code>rootVolume</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Volume">
Volume
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RootVolume encapsulates the configuration options for the root volume</p>
</td>
</tr>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string
(do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)</p>
</td>
</tr>
<tr>
<td>
<code>versionNumber</code><br/>
<em>
int64
</em>
</td>
<td>
<p>VersionNumber is the version of the launch template that is applied.
Typically a new version is created when at least one of the following happens:
1) A new launch template spec is applied.
2) One or more parameters in an existing template is changed.
3) A new AMI is discovered.</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an array of references to security groups that should be applied to the
instances. These security groups would be set in addition to any security groups defined
at the cluster level or in the actuator.</p>
</td>
</tr>
<tr>
<td>
<code>spotMarketOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SpotMarketOptions">
SpotMarketOptions
</a>
</em>
</td>
<td>
<p>SpotMarketOptions are options for configuring AWSMachinePool instances to be run using AWS Spot instances.</p>
</td>
</tr>
<tr>
<td>
<code>instanceMetadataOptions</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstanceMetadataOptions">
InstanceMetadataOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceMetadataOptions defines the behavior for applying metadata to instances.</p>
</td>
</tr>
<tr>
<td>
<code>privateDnsName</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.PrivateDNSName">
PrivateDNSName
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>PrivateDNSName is the options for the instance hostname.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePool">AWSMachinePool
</h3>
<p>
<p>AWSMachinePool is the Schema for the awsmachinepools API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">
AWSMachinePoolSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the ARN of the associated ASG</p>
</td>
</tr>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MinSize defines the minimum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MaxSize defines the maximum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneSubnetType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AZSubnetType">
AZSubnetType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AvailabilityZoneSubnetType specifies which type of subnets to use when an availability zone is specified.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets is an array of subnet configurations</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<p>AWSLaunchTemplate specifies the launch template and version to use when an instance is launched.</p>
</td>
</tr>
<tr>
<td>
<code>mixedInstancesPolicy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.MixedInstancesPolicy">
MixedInstancesPolicy
</a>
</em>
</td>
<td>
<p>MixedInstancesPolicy describes how multiple instance types will be used by the ASG.</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the identification IDs of machine instances provided by the provider.
This field must match the provider IDs as seen on the node objects corresponding to a machine pool&rsquo;s machine instances.</p>
</td>
</tr>
<tr>
<td>
<code>defaultCoolDown</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of time, in seconds, after a scaling activity completes before another scaling activity can start.
If no value is supplied by user a default value of 300 seconds is set</p>
</td>
</tr>
<tr>
<td>
<code>defaultInstanceWarmup</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of time, in seconds, until a new instance is considered to
have finished initializing and resource consumption to become stable
after it enters the InService state.
If no value is supplied by user a default value of 300 seconds is set</p>
</td>
</tr>
<tr>
<td>
<code>refreshPreferences</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RefreshPreferences">
RefreshPreferences
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RefreshPreferences describes set of preferences associated with the instance refresh request.</p>
</td>
</tr>
<tr>
<td>
<code>capacityRebalance</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Enable or disable the capacity rebalance autoscaling group feature</p>
</td>
</tr>
<tr>
<td>
<code>suspendProcesses</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SuspendProcessesTypes">
SuspendProcessesTypes
</a>
</em>
</td>
<td>
<p>SuspendProcesses defines a list of processes to suspend for the given ASG. This is constantly reconciled.
If a process is removed from this list it will automatically be resumed.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolStatus">
AWSMachinePoolStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolInstanceStatus">AWSMachinePoolInstanceStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolStatus">AWSMachinePoolStatus</a>)
</p>
<p>
<p>AWSMachinePoolInstanceStatus defines the status of the AWSMachinePoolInstance.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceID is the identification of the Machine Instance within ASG</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version defines the Kubernetes version for the Machine Instance</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePool">AWSMachinePool</a>)
</p>
<p>
<p>AWSMachinePoolSpec defines the desired state of AWSMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the ARN of the associated ASG</p>
</td>
</tr>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MinSize defines the minimum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
<p>MaxSize defines the maximum size of the group.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneSubnetType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AZSubnetType">
AZSubnetType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AvailabilityZoneSubnetType specifies which type of subnets to use when an availability zone is specified.</p>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSResourceReference">
[]AWSResourceReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Subnets is an array of subnet configurations</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
AWS provider.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<p>AWSLaunchTemplate specifies the launch template and version to use when an instance is launched.</p>
</td>
</tr>
<tr>
<td>
<code>mixedInstancesPolicy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.MixedInstancesPolicy">
MixedInstancesPolicy
</a>
</em>
</td>
<td>
<p>MixedInstancesPolicy describes how multiple instance types will be used by the ASG.</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the identification IDs of machine instances provided by the provider.
This field must match the provider IDs as seen on the node objects corresponding to a machine pool&rsquo;s machine instances.</p>
</td>
</tr>
<tr>
<td>
<code>defaultCoolDown</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of time, in seconds, after a scaling activity completes before another scaling activity can start.
If no value is supplied by user a default value of 300 seconds is set</p>
</td>
</tr>
<tr>
<td>
<code>defaultInstanceWarmup</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of time, in seconds, until a new instance is considered to
have finished initializing and resource consumption to become stable
after it enters the InService state.
If no value is supplied by user a default value of 300 seconds is set</p>
</td>
</tr>
<tr>
<td>
<code>refreshPreferences</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RefreshPreferences">
RefreshPreferences
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RefreshPreferences describes set of preferences associated with the instance refresh request.</p>
</td>
</tr>
<tr>
<td>
<code>capacityRebalance</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Enable or disable the capacity rebalance autoscaling group feature</p>
</td>
</tr>
<tr>
<td>
<code>suspendProcesses</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SuspendProcessesTypes">
SuspendProcessesTypes
</a>
</em>
</td>
<td>
<p>SuspendProcesses defines a list of processes to suspend for the given ASG. This is constantly reconciled.
If a process is removed from this list it will automatically be resumed.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolStatus">AWSMachinePoolStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePool">AWSMachinePool</a>)
</p>
<p>
<p>AWSMachinePoolStatus defines the observed state of AWSMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready is true when the provider resource is ready.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Replicas is the most recently observed number of replicas</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the AWSMachinePool.</p>
</td>
</tr>
<tr>
<td>
<code>instances</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolInstanceStatus">
[]AWSMachinePoolInstanceStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Instances contains the status for each instance in the pool</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateID</code><br/>
<em>
string
</em>
</td>
<td>
<p>The ID of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The version of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the Machine and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the Machine and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>asgStatus</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ASGStatus">
ASGStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePool">AWSManagedMachinePool
</h3>
<p>
<p>AWSManagedMachinePool is the Schema for the awsmanagedmachinepools API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">
AWSManagedMachinePoolSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>eksNodegroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSNodegroupName specifies the name of the nodegroup in AWS
corresponding to this MachinePool. If you don&rsquo;t specify a name
then a default name will be created based on the namespace and
name of the managed machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneSubnetType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AZSubnetType">
AZSubnetType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AvailabilityZoneSubnetType specifies which type of subnets to use when an availability zone is specified.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the node group role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for the node group.
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>amiVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIVersion defines the desired AMI release version. If no version number
is supplied then the latest version for the Kubernetes version
will be used</p>
</td>
</tr>
<tr>
<td>
<code>amiType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachineAMIType">
ManagedMachineAMIType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIType defines the AMI type</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels specifies labels for the Kubernetes node objects</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Taints">
Taints
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints specifies the taints to apply to the nodes of the machine pool</p>
</td>
</tr>
<tr>
<td>
<code>diskSize</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>DiskSize specifies the root disk size</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceType specifies the AWS instance type</p>
</td>
</tr>
<tr>
<td>
<code>scaling</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachinePoolScaling">
ManagedMachinePoolScaling
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Scaling specifies scaling for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>remoteAccess</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedRemoteAccess">
ManagedRemoteAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RemoteAccess specifies how machines can be accessed remotely</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the provider IDs of instances in the
autoscaling group corresponding to the nodegroup represented by this
machine pool</p>
</td>
</tr>
<tr>
<td>
<code>capacityType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachinePoolCapacityType">
ManagedMachinePoolCapacityType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CapacityType specifies the capacity type for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>updateConfig</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.UpdateConfig">
UpdateConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>UpdateConfig holds the optional config to control the behaviour of the update
to the nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AWSLaunchTemplate specifies the launch template to use to create the managed node group.
If AWSLaunchTemplate is specified, certain node group configuraions outside of launch template
are prohibited (<a href="https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html">https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html</a>).</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolStatus">
AWSManagedMachinePoolStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePool">AWSManagedMachinePool</a>)
</p>
<p>
<p>AWSManagedMachinePoolSpec defines the desired state of AWSManagedMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>eksNodegroupName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EKSNodegroupName specifies the name of the nodegroup in AWS
corresponding to this MachinePool. If you don&rsquo;t specify a name
then a default name will be created based on the namespace and
name of the managed machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZones</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AvailabilityZones is an array of availability zones instances can run in</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZoneSubnetType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AZSubnetType">
AZSubnetType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AvailabilityZoneSubnetType specifies which type of subnets to use when an availability zone is specified.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleAdditionalPolicies</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleAdditionalPolicies allows you to attach additional polices to
the node group role. You must enable the EKSAllowAddRoles
feature flag to incorporate these into the created role.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for the node group.
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>amiVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIVersion defines the desired AMI release version. If no version number
is supplied then the latest version for the Kubernetes version
will be used</p>
</td>
</tr>
<tr>
<td>
<code>amiType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachineAMIType">
ManagedMachineAMIType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AMIType defines the AMI type</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels specifies labels for the Kubernetes node objects</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Taints">
Taints
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints specifies the taints to apply to the nodes of the machine pool</p>
</td>
</tr>
<tr>
<td>
<code>diskSize</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>DiskSize specifies the root disk size</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InstanceType specifies the AWS instance type</p>
</td>
</tr>
<tr>
<td>
<code>scaling</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachinePoolScaling">
ManagedMachinePoolScaling
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Scaling specifies scaling for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>remoteAccess</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedRemoteAccess">
ManagedRemoteAccess
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>RemoteAccess specifies how machines can be accessed remotely</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList are the provider IDs of instances in the
autoscaling group corresponding to the nodegroup represented by this
machine pool</p>
</td>
</tr>
<tr>
<td>
<code>capacityType</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachinePoolCapacityType">
ManagedMachinePoolCapacityType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CapacityType specifies the capacity type for the ASG behind this pool</p>
</td>
</tr>
<tr>
<td>
<code>updateConfig</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.UpdateConfig">
UpdateConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>UpdateConfig holds the optional config to control the behaviour of the update
to the nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>awsLaunchTemplate</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSLaunchTemplate">
AWSLaunchTemplate
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AWSLaunchTemplate specifies the launch template to use to create the managed node group.
If AWSLaunchTemplate is specified, certain node group configuraions outside of launch template
are prohibited (<a href="https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html">https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html</a>).</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolStatus">AWSManagedMachinePoolStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePool">AWSManagedMachinePool</a>)
</p>
<p>
<p>AWSManagedMachinePoolStatus defines the observed state of AWSManagedMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the AWSManagedMachinePool nodegroup has joined
the cluster</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Replicas is the most recently observed number of replicas.</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The ID of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>launchTemplateVersion</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The version of the launch template</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the MachinePool and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of MachinePools
can be added as events to the MachinePool object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the MachinePool and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the MachinePool&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of MachinePools
can be added as events to the MachinePool object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the managed machine pool</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AZSubnetType">AZSubnetType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>AZSubnetType is the type of subnet to use when an availability zone is specified.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;all&#34;</p></td>
<td><p>AZSubnetTypeAll is all subnets in an availability zone.</p>
</td>
</tr><tr><td><p>&#34;private&#34;</p></td>
<td><p>AZSubnetTypePrivate is a private subnet.</p>
</td>
</tr><tr><td><p>&#34;public&#34;</p></td>
<td><p>AZSubnetTypePublic is a public subnet.</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.AutoScalingGroup">AutoScalingGroup
</h3>
<p>
<p>AutoScalingGroup describes an AWS autoscaling group.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>The tags associated with the instance.</p>
</td>
</tr>
<tr>
<td>
<code>tags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>desiredCapacity</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>placementGroup</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>subnets</code><br/>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>defaultCoolDown</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>defaultInstanceWarmup</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>capacityRebalance</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>mixedInstancesPolicy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.MixedInstancesPolicy">
MixedInstancesPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>Status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ASGStatus">
ASGStatus
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instances</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Instance">
[]Instance
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>currentlySuspendProcesses</code><br/>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.BlockDeviceMapping">BlockDeviceMapping
</h3>
<p>
<p>BlockDeviceMapping specifies the block devices for the instance.
You can specify virtual devices and EBS volumes.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>deviceName</code><br/>
<em>
string
</em>
</td>
<td>
<p>The device name exposed to the EC2 instance (for example, /dev/sdh or xvdh).</p>
</td>
</tr>
<tr>
<td>
<code>ebs</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.EBS">
EBS
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>You can specify either VirtualName or Ebs, but not both.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.EBS">EBS
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.BlockDeviceMapping">BlockDeviceMapping</a>)
</p>
<p>
<p>EBS can be used to automatically set up EBS volumes when an instance is launched.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>encrypted</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Encrypted is whether the volume should be encrypted or not.</p>
</td>
</tr>
<tr>
<td>
<code>volumeSize</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>The size of the volume, in GiB.
This can be a number from 1-1,024 for standard, 4-16,384 for io1, 1-16,384
for gp2, and 500-16,384 for st1 and sc1. If you specify a snapshot, the volume
size must be equal to or larger than the snapshot size.</p>
</td>
</tr>
<tr>
<td>
<code>volumeType</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The volume type
For more information, see Amazon EBS Volume Types (<a href="https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html">https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html</a>)</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.FargateProfileSpec">FargateProfileSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSFargateProfile">AWSFargateProfile</a>)
</p>
<p>
<p>FargateProfileSpec defines the desired state of FargateProfile.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>clusterName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ClusterName is the name of the Cluster this object belongs to.</p>
</td>
</tr>
<tr>
<td>
<code>profileName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ProfileName specifies the profile name.</p>
</td>
</tr>
<tr>
<td>
<code>subnetIDs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SubnetIDs specifies which subnets are used for the
auto scaling group of this nodegroup.</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
ones added by default.</p>
</td>
</tr>
<tr>
<td>
<code>roleName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>RoleName specifies the name of IAM role for this fargate pool
If the role is pre-existing we will treat it as unmanaged
and not delete it on deletion. If the EKSEnableIAM feature
flag is true and no name is supplied then a role is created.</p>
</td>
</tr>
<tr>
<td>
<code>selectors</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.FargateSelector">
[]FargateSelector
</a>
</em>
</td>
<td>
<p>Selectors specify fargate pod selectors.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.FargateProfileStatus">FargateProfileStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSFargateProfile">AWSFargateProfile</a>)
</p>
<p>
<p>FargateProfileStatus defines the observed state of FargateProfile.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the FargateProfile is available.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the FargateProfile and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the FargateProfile&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of
FargateProfiles can be added as events to the FargateProfile object
and/or logged in the controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the FargateProfile and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the FargateProfile&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of
FargateProfiles can be added as events to the FargateProfile
object and/or logged in the controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current state of the Fargate profile.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.FargateSelector">FargateSelector
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.FargateProfileSpec">FargateProfileSpec</a>)
</p>
<p>
<p>FargateSelector specifies a selector for pods that should run on this fargate pool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Labels specifies which pod labels this selector should match.</p>
</td>
</tr>
<tr>
<td>
<code>namespace</code><br/>
<em>
string
</em>
</td>
<td>
<p>Namespace specifies which namespace this selector should match.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.InstancesDistribution">InstancesDistribution
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.MixedInstancesPolicy">MixedInstancesPolicy</a>)
</p>
<p>
<p>InstancesDistribution to configure distribution of On-Demand Instances and Spot Instances.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>onDemandAllocationStrategy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.OnDemandAllocationStrategy">
OnDemandAllocationStrategy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>spotAllocationStrategy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.SpotAllocationStrategy">
SpotAllocationStrategy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>onDemandBaseCapacity</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>onDemandPercentageAboveBaseCapacity</code><br/>
<em>
int64
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachineAMIType">ManagedMachineAMIType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedMachineAMIType specifies which AWS AMI to use for a managed MachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;AL2_ARM_64&#34;</p></td>
<td><p>Al2Arm64 is the Arm AMI type.</p>
</td>
</tr><tr><td><p>&#34;AL2_x86_64&#34;</p></td>
<td><p>Al2x86_64 is the default AMI type.</p>
</td>
</tr><tr><td><p>&#34;AL2_x86_64_GPU&#34;</p></td>
<td><p>Al2x86_64GPU is the x86-64 GPU AMI type.</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachinePoolCapacityType">ManagedMachinePoolCapacityType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedMachinePoolCapacityType specifies the capacity type to be used for the managed MachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;onDemand&#34;</p></td>
<td><p>ManagedMachinePoolCapacityTypeOnDemand is the default capacity type, to launch on-demand instances.</p>
</td>
</tr><tr><td><p>&#34;spot&#34;</p></td>
<td><p>ManagedMachinePoolCapacityTypeSpot is the spot instance capacity type to launch spot instances.</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ManagedMachinePoolScaling">ManagedMachinePoolScaling
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedMachinePoolScaling specifies scaling options.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>minSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxSize</code><br/>
<em>
int32
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ManagedRemoteAccess">ManagedRemoteAccess
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>ManagedRemoteAccess specifies remote access settings for EC2 instances.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>sshKeyName</code><br/>
<em>
string
</em>
</td>
<td>
<p>SSHKeyName specifies which EC2 SSH key can be used to access machines.
If left empty, the key from the control plane is used.</p>
</td>
</tr>
<tr>
<td>
<code>sourceSecurityGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>SourceSecurityGroups specifies which security groups are allowed access</p>
</td>
</tr>
<tr>
<td>
<code>public</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Public specifies whether to open port 22 to the public internet</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.MixedInstancesPolicy">MixedInstancesPolicy
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.AutoScalingGroup">AutoScalingGroup</a>)
</p>
<p>
<p>MixedInstancesPolicy for an Auto Scaling group.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instancesDistribution</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstancesDistribution">
InstancesDistribution
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>overrides</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Overrides">
[]Overrides
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.OnDemandAllocationStrategy">OnDemandAllocationStrategy
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstancesDistribution">InstancesDistribution</a>)
</p>
<p>
<p>OnDemandAllocationStrategy indicates how to allocate instance types to fulfill On-Demand capacity.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Overrides">Overrides
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.MixedInstancesPolicy">MixedInstancesPolicy</a>)
</p>
<p>
<p>Overrides are used to override the instance type specified by the launch template with multiple
instance types that can be used to launch On-Demand Instances and Spot Instances.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Processes">Processes
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.SuspendProcessesTypes">SuspendProcessesTypes</a>)
</p>
<p>
<p>Processes defines the processes which can be enabled or disabled individually.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>launch</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>terminate</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>addToLoadBalancer</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>alarmNotification</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>azRebalance</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>healthCheck</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>instanceRefresh</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replaceUnhealthy</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>scheduledActions</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ROSACluster">ROSACluster
</h3>
<p>
<p>ROSACluster is the Schema for the ROSAClusters API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ROSAClusterSpec">
ROSAClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.ROSAClusterStatus">
ROSAClusterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ROSAClusterSpec">ROSAClusterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.ROSACluster">ROSACluster</a>)
</p>
<p>
<p>ROSAClusterSpec defines the desired state of ROSACluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ROSAClusterStatus">ROSAClusterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.ROSACluster">ROSACluster</a>)
</p>
<p>
<p>ROSAClusterStatus defines the observed state of ROSACluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready is when the ROSAControlPlane has a API server URL.</p>
</td>
</tr>
<tr>
<td>
<code>failureDomains</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.FailureDomains
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureDomains specifies a list fo available availability zones that can be used</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.ROSAMachinePool">ROSAMachinePool
</h3>
<p>
<p>ROSAMachinePool is the Schema for the rosamachinepools API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolSpec">
RosaMachinePoolSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>nodePoolName</code><br/>
<em>
string
</em>
</td>
<td>
<p>NodePoolName specifies the name of the nodepool in Rosa
must be a valid DNS-1035 label, so it must consist of lower case alphanumeric and have a max length of 15 characters.</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version specifies the OpenShift version of the nodes associated with this machinepool.
ROSAControlPlane version is used if not set.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZone</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AvailabilityZone is an optinal field specifying the availability zone where instances of this machine pool should run
For Multi-AZ clusters, you can create a machine pool in a Single-AZ of your choice.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels specifies labels for the Kubernetes node objects</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaTaint">
[]RosaTaint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints specifies the taints to apply to the nodes of the machine pool</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags are user-defined tags to be added on the underlying EC2 instances associated with this machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>autoRepair</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>AutoRepair specifies whether health checks should be enabled for machines
in the NodePool. The default is false.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType specifies the AWS instance type</p>
</td>
</tr>
<tr>
<td>
<code>autoscaling</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolAutoScaling">
RosaMachinePoolAutoScaling
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Autoscaling specifies auto scaling behaviour for this MachinePool.
required if Replicas is not configured</p>
</td>
</tr>
<tr>
<td>
<code>tuningConfigs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>TuningConfigs specifies the names of the tuning configs to be applied to this MachinePool.
Tuning configs must already exist.</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an optional set of security groups to associate
with all node instances of the machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList contain a ProviderID for each machine instance that&rsquo;s currently managed by this machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>nodeDrainGracePeriod</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NodeDrainGracePeriod is grace period for how long Pod Disruption Budget-protected workloads will be
respected during upgrades. After this grace period, any workloads protected by Pod Disruption
Budgets that have not been successfully drained from a node will be forcibly evicted.</p>
<p>Valid values are from 0 to 1 week(10080m|168h) .
0 or empty value means that the MachinePool can be drained without any time limitation.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolStatus">
RosaMachinePoolStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.RefreshPreferences">RefreshPreferences
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec</a>)
</p>
<p>
<p>RefreshPreferences defines the specs for instance refreshing.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>disable</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Disable, if true, disables instance refresh from triggering when new launch templates are detected.
This is useful in scenarios where ASG nodes are externally managed.</p>
</td>
</tr>
<tr>
<td>
<code>strategy</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The strategy to use for the instance refresh. The only valid value is Rolling.
A rolling update is an update that is applied to all instances in an Auto
Scaling group until all instances have been updated.</p>
</td>
</tr>
<tr>
<td>
<code>instanceWarmup</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>The number of seconds until a newly launched instance is configured and ready
to use. During this time, the next replacement will not be initiated.
The default is to use the value for the health check grace period defined for the group.</p>
</td>
</tr>
<tr>
<td>
<code>minHealthyPercentage</code><br/>
<em>
int64
</em>
</td>
<td>
<em>(Optional)</em>
<p>The amount of capacity as a percentage in ASG that must remain healthy
during an instance refresh. The default is 90.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolAutoScaling">RosaMachinePoolAutoScaling
</h3>
<p>
(<em>Appears on:</em><a href="#controlplane.cluster.x-k8s.io/v1beta2.DefaultMachinePoolSpec">DefaultMachinePoolSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolSpec">RosaMachinePoolSpec</a>)
</p>
<p>
<p>RosaMachinePoolAutoScaling specifies scaling options.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>minReplicas</code><br/>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxReplicas</code><br/>
<em>
int
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolSpec">RosaMachinePoolSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.ROSAMachinePool">ROSAMachinePool</a>)
</p>
<p>
<p>RosaMachinePoolSpec defines the desired state of RosaMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>nodePoolName</code><br/>
<em>
string
</em>
</td>
<td>
<p>NodePoolName specifies the name of the nodepool in Rosa
must be a valid DNS-1035 label, so it must consist of lower case alphanumeric and have a max length of 15 characters.</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version specifies the OpenShift version of the nodes associated with this machinepool.
ROSAControlPlane version is used if not set.</p>
</td>
</tr>
<tr>
<td>
<code>availabilityZone</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AvailabilityZone is an optinal field specifying the availability zone where instances of this machine pool should run
For Multi-AZ clusters, you can create a machine pool in a Single-AZ of your choice.</p>
</td>
</tr>
<tr>
<td>
<code>subnet</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels specifies labels for the Kubernetes node objects</p>
</td>
</tr>
<tr>
<td>
<code>taints</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaTaint">
[]RosaTaint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Taints specifies the taints to apply to the nodes of the machine pool</p>
</td>
</tr>
<tr>
<td>
<code>additionalTags</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Tags">
Tags
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalTags are user-defined tags to be added on the underlying EC2 instances associated with this machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>autoRepair</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>AutoRepair specifies whether health checks should be enabled for machines
in the NodePool. The default is false.</p>
</td>
</tr>
<tr>
<td>
<code>instanceType</code><br/>
<em>
string
</em>
</td>
<td>
<p>InstanceType specifies the AWS instance type</p>
</td>
</tr>
<tr>
<td>
<code>autoscaling</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolAutoScaling">
RosaMachinePoolAutoScaling
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Autoscaling specifies auto scaling behaviour for this MachinePool.
required if Replicas is not configured</p>
</td>
</tr>
<tr>
<td>
<code>tuningConfigs</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>TuningConfigs specifies the names of the tuning configs to be applied to this MachinePool.
Tuning configs must already exist.</p>
</td>
</tr>
<tr>
<td>
<code>additionalSecurityGroups</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>AdditionalSecurityGroups is an optional set of security groups to associate
with all node instances of the machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>providerIDList</code><br/>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderIDList contain a ProviderID for each machine instance that&rsquo;s currently managed by this machine pool.</p>
</td>
</tr>
<tr>
<td>
<code>nodeDrainGracePeriod</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NodeDrainGracePeriod is grace period for how long Pod Disruption Budget-protected workloads will be
respected during upgrades. After this grace period, any workloads protected by Pod Disruption
Budgets that have not been successfully drained from a node will be forcibly evicted.</p>
<p>Valid values are from 0 to 1 week(10080m|168h) .
0 or empty value means that the MachinePool can be drained without any time limitation.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolStatus">RosaMachinePoolStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.ROSAMachinePool">ROSAMachinePool</a>)
</p>
<p>
<p>RosaMachinePoolStatus defines the observed state of RosaMachinePool.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Ready denotes that the RosaMachinePool nodepool has joined
the cluster</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<em>(Optional)</em>
<p>Replicas is the most recently observed number of replicas.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the managed machine pool</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the state and will be set to a descriptive error message.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the spec or the configuration of
the controller, and that manual intervention is required.</p>
</td>
</tr>
<tr>
<td>
<code>id</code><br/>
<em>
string
</em>
</td>
<td>
<p>ID is the ID given by ROSA.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.RosaTaint">RosaTaint
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.RosaMachinePoolSpec">RosaMachinePoolSpec</a>)
</p>
<p>
<p>RosaTaint represents a taint to be applied to a node.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>The taint key to be applied to a node.</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>The taint value corresponding to the taint key.</p>
</td>
</tr>
<tr>
<td>
<code>effect</code><br/>
<em>
<a href="https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#tainteffect-v1-core">
Kubernetes core/v1.TaintEffect
</a>
</em>
</td>
<td>
<p>The effect of the taint on pods that do not tolerate the taint.
Valid effects are NoSchedule, PreferNoSchedule and NoExecute.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SpotAllocationStrategy">SpotAllocationStrategy
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.InstancesDistribution">InstancesDistribution</a>)
</p>
<p>
<p>SpotAllocationStrategy indicates how to allocate instances across Spot Instance pools.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.SuspendProcessesTypes">SuspendProcessesTypes
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSMachinePoolSpec">AWSMachinePoolSpec</a>)
</p>
<p>
<p>SuspendProcessesTypes contains user friendly auto-completable values for suspended process names.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>all</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>processes</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.Processes">
Processes
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Tags">Tags
(<code>map[string]string</code> alias)</p></h3>
<p>
<p>Tags is a mapping for tags.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Taint">Taint
</h3>
<p>
<p>Taint defines the specs for a Kubernetes taint.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>effect</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1beta2.TaintEffect">
TaintEffect
</a>
</em>
</td>
<td>
<p>Effect specifies the effect for the taint</p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<p>Key is the key of the taint</p>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
string
</em>
</td>
<td>
<p>Value is the value of the taint</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.TaintEffect">TaintEffect
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.Taint">Taint</a>)
</p>
<p>
<p>TaintEffect is the effect for a Kubernetes taint.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.Taints">Taints
(<code>[]sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2.Taint</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>Taints is an array of Taints.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1beta2.UpdateConfig">UpdateConfig
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1beta2.AWSManagedMachinePoolSpec">AWSManagedMachinePoolSpec</a>)
</p>
<p>
<p>UpdateConfig is the configuration options for updating a nodegroup. Only one of MaxUnavailable
and MaxUnavailablePercentage should be specified.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxUnavailable</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxUnavailable is the maximum number of nodes unavailable at once during a version update.
Nodes will be updated in parallel. The maximum number is 100.</p>
</td>
</tr>
<tr>
<td>
<code>maxUnavailablePercentage</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>MaxUnavailablePercentage is the maximum percentage of nodes unavailable during a version update. This
percentage of nodes will be updated in parallel, up to 100 nodes at once.</p>
</td>
</tr>
</tbody>
</table>
<hr/>

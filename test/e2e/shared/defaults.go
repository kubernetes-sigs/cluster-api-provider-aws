package shared

const (
	DefaultSSHKeyPairName        = "cluster-api-provider-aws-sigs-k8s-io"
	AMIPrefix                    = "capa-ami-ubuntu-18.04-"
	DefaultImageLookupOrg        = "258751437250"
	KubernetesVersion            = "KUBERNETES_VERSION"
	CNIPath                      = "CNI"
	CNIResources                 = "CNI_RESOURCES"
	AwsNodeMachineType           = "AWS_NODE_MACHINE_TYPE"
	AwsAvailabilityZone1         = "AWS_AVAILABILITY_ZONE_1"
	AwsAvailabilityZone2         = "AWS_AVAILABILITY_ZONE_2"
	MultiAzFlavor                = "multi-az"
	LimitAzFlavor                = "limit-az"
	SpotInstancesFlavor          = "spot-instances"
	SSMFlavor                    = "ssm"
	StorageClassFailureZoneLabel = "failure-domain.beta.kubernetes.io/zone"
)

package filter

const (
	// TagNameKubernetesClusterPrefix is the tag name we use to differentiate multiple
	// logically independent clusters running in the same AZ.
	// The tag key = TagNameKubernetesClusterPrefix + clusterID
	// The tag value is an ownership value
	TagNameKubernetesClusterPrefix = "kubernetes.io/cluster/"

	// TagNameAWSProviderManaged is the tag name we use to differentiate
	// cluster-api-provider-aws owned components from other tooling that
	// uses TagNameKubernetesClusterPrefix
	TagNameAWSProviderManaged = "sigs.k8s.io/cluster-api-provider-aws/managed"

	// TagNameAWSClusterAPIRole is the tag name we use to mark roles for resources
	// dedicated to this cluster api provider implementation.
	TagNameAWSClusterAPIRole = "sigs.k8s.io/cluster-api-provider-aws/role"

	// TagValueAPIServerRole describes the value for the apiserver role
	TagValueAPIServerRole = "apiserver"

	// TagValueBastionRole describes the value for the bastion role
	TagValueBastionRole = "bastion"

	// TagValueCommonRole describes the value for the common role
	TagValueCommonRole = "common"
)

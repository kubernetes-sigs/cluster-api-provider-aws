package ami

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
)

const (
	// EKS AMI ID SSM Parameter name
	eksAmiSSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2%s/recommended/image_id"
)

func formatVersionForEKS(version string) (string, error) {
	parsed, err := semver.ParseTolerant(version)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse version")
	}

	return fmt.Sprintf("%d.%d", parsed.Major, parsed.Minor), nil
}

func EKSAMIParameter(kubernetesVersion string, imageType infrav1exp.ManagedMachineAMIType) (string, error) {
	// format ssm parameter path properly
	formattedVersion, err := formatVersionForEKS(kubernetesVersion)
	if err != nil {
		return "", err
	}
	var imageTypeSuffix string
	switch imageType {
	case infrav1exp.Al2x86_64, "":
		imageTypeSuffix = ""
	case infrav1exp.Al2x86_64GPU:
		imageTypeSuffix = "-gpu"
	case infrav1exp.Al2Arm64:
		imageTypeSuffix = "-arm64"
	}

	return fmt.Sprintf(eksAmiSSMParameterFormat, formattedVersion, imageTypeSuffix), nil
}

package deployer_test

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/deployer"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

func statusToRawExtension(status *v1alpha1.AWSClusterProviderStatus, t *testing.T) *runtime.RawExtension {
	codec, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create a codec: %v", err)
	}
	out, err := codec.EncodeProviderStatus(status)
	if err != nil {
		t.Fatalf("failed to convert status to extension: %v", err)
	}
	return out
}

func TestGetIP(t *testing.T) {
	testcases := []struct {
		name                  string
		clusterProviderStatus v1alpha1.AWSClusterProviderStatus
	}{
		{
			name: "simple test",
			clusterProviderStatus: v1alpha1.AWSClusterProviderStatus{
				Network: v1alpha1.Network{
					APIServerELB: v1alpha1.ClassicELB{
						DNSName: "example.com",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			deployer := deployer.AWSDeployer{}
			cluster := &clusterv1.Cluster{
				Status: clusterv1.ClusterStatus{
					ProviderStatus: statusToRawExtension(&tc.clusterProviderStatus, t),
				},
			}
			machine := &clusterv1.Machine{}

			ip, err := deployer.GetIP(cluster, machine)
			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}
			if ip != tc.clusterProviderStatus.Network.APIServerELB.DNSName {
				t.Fatalf("expected %v got %v", tc.clusterProviderStatus.Network.APIServerELB.DNSName, ip)
			}
		})
	}
}

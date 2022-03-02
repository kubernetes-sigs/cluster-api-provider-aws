package ami

import (
	"reflect"
	"testing"
)

func TestAllPatchVersions(t *testing.T) {
	tests := []struct {
		name          string
		latestVersion string
		want          []string
	}{
		{
			name:          "Should return 3 different patch versions for v1.23.3",
			latestVersion: "v1.23.3",
			want:          []string{"v1.23.2", "v1.23.1", "v1.23.0"},
		},
		{
			name:          "Should return empty patch versions for v1.23.0",
			latestVersion: "v1.23.0",
			want:          []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allPatchesForVersion(tt.latestVersion)
			if err != nil {
				t.Fatalf("error while fetching patch versions %+v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("allPatchesForVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package custom

type AMIBuildConfig struct {
	K8sReleases map[string]string `json:"k8s_releases"`
}

type AMIBuildConfigDefaults struct {
	Amazon2    map[string]string `json:"amazon-2"`
	Centos7    map[string]string `json:"centos-7"`
	Flatcar    map[string]string `json:"flatcar"`
	Ubuntu1804 map[string]string `json:"ubuntu-1804"`
	Ubuntu2004 map[string]string `json:"ubuntu-2004"`
	Default    map[string]string `json:"default"`
}

type ReleaseVersion struct {
	Major int
	Minor int
	Patch int
}

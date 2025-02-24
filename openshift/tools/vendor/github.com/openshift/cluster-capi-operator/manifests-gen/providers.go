package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/klauspost/compress/zstd"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	configclient "sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/repository"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/yamlprocessor"
	utilyaml "sigs.k8s.io/cluster-api/util/yaml"
	"sigs.k8s.io/yaml"
)

const (
	powerVSProvider                = "powervs"
	ibmCloudProvider               = "ibmcloud"
	coreCAPIProvider               = "cluster-api"
	metadataFilePath               = "./metadata.yaml"
	defaultKustomizeComponentsPath = "./config/default"
	// customizedComponentsFilename is a name for file containing customized infrastructure components.
	// This file helps with code review as it is always uncompressed unlike the components configMap.
	customizedComponentsFilename = "infrastructure-components-openshift.yaml"
)

type provider struct {
	Name       string                    `json:"name"`
	Type       clusterctlv1.ProviderType `json:"type"`
	Version    string                    `json:"version"`
	components repository.Components
	metadata   []byte
}

// loadComponents loads components from the given provider.
func (p *provider) loadComponents() error {
	// Create new clusterctl config client
	configClient, err := configclient.New(context.Background(), "")
	if err != nil {
		return fmt.Errorf("error creating clusterctl config client: %w", err)
	}

	// Create new clusterctl provider client
	providerConfig, err := configClient.Providers().Get(p.Name, p.Type)
	if err != nil {
		return fmt.Errorf("error creating clusterctl provider client: %w", err)
	}

	// Set options for yaml processor
	options := repository.ComponentsOptions{
		TargetNamespace:     targetNamespace,
		SkipTemplateProcess: true,
	}

	// Compile assets using kustomize.
	kustomizeComponentsPath := path.Join(projDir, *kustomizeDir)
	fmt.Printf("> Generating OpenShift manifests based on kustomize.yaml from %q\n", kustomizeComponentsPath)
	rawComponents, err := fetchAndCompileComponents(kustomizeComponentsPath)
	if err != nil {
		return fmt.Errorf("error fetching and compiling assets using kustomize: %w", err)
	}

	// Ininitialize new clusterctl repository components, this should some yaml processing
	p.components, err = repository.NewComponents(repository.ComponentsInput{
		Provider:     providerConfig,
		ConfigClient: configClient,
		Processor:    yamlprocessor.NewSimpleProcessor(),
		RawYaml:      rawComponents,
		Options:      options,
	})
	if err != nil {
		return fmt.Errorf("error initializing new clusterctl components repository: %w", err)
	}

	content, err := os.ReadFile(path.Join(projDir, metadataFilePath))
	if err != nil {
		return fmt.Errorf("error while reading metadata file: %w", err)
	}

	p.metadata = content

	return nil
}

func (p *provider) providerTypeName() string {
	return strings.ReplaceAll(strings.ToLower(string(p.Type)), "provider", "")
}

// writeProviderComponentsToManifest allows to write provider components directly to a manifest.
// This differs from writeProviderComponentsConfigmap as it won't store the components in a ConfigMap
// but directly on the YAML manifests as YAML representation of an unstructured objects.
func (p *provider) writeProviderComponentsToManifest(fileName string, objs []unstructured.Unstructured) error {
	if len(objs) == 0 {
		return nil
	}

	combined, err := utilyaml.FromUnstructured(objs)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(*manifestsPath, fileName), ensureNewLine(combined), 0600)
}

// writeProviderCustomizedComponents writes the customized infrastructure components to allow for code review
func (p *provider) writeProviderCustomizedComponents(resourceMap map[resourceKey][]unstructured.Unstructured) error {
	crds, err := utilyaml.FromUnstructured(resourceMap[crdKey])
	if err != nil {
		return fmt.Errorf("error converting unstructured object to YAML: %w", err)
	}

	other, err := utilyaml.FromUnstructured(resourceMap[otherKey])
	if err != nil {
		return fmt.Errorf("error converting unstructured object to YAML: %w", err)
	}

	combined := utilyaml.JoinYaml(crds, other)

	return os.WriteFile(path.Join(*basePath, "openshift", customizedComponentsFilename), ensureNewLine(combined), 0600)
}

// writeProviderComponentsConfigmap allows to write provider components to the provider (transport) ConfigMap.
func (p *provider) writeProviderComponentsConfigmap(fileName string, objs []unstructured.Unstructured) error {
	annotations := openshiftAnnotations
	annotations[featureSetAnnotationKey] = featureSetAnnotationValue

	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: targetNamespace,
			Labels: map[string]string{
				"provider.cluster.x-k8s.io/name":    p.Name,
				"provider.cluster.x-k8s.io/type":    p.providerTypeName(),
				"provider.cluster.x-k8s.io/version": p.Version,
			},
			Annotations: annotations,
		},
		Data: map[string]string{
			"metadata": string(p.metadata),
		},
	}

	cm, err := addComponentsToConfigMap(cm, objs)
	if err != nil {
		return fmt.Errorf("failed to inject components into configmap: %w", err)
	}

	cmYaml, err := yaml.Marshal(cm)
	if err != nil {
		return fmt.Errorf("error marshaling provider ConfigMap to YAML: %w", err)
	}

	return os.WriteFile(path.Join(*manifestsPath, fileName), ensureNewLine(cmYaml), 0600)
}

// addComponentsToConfigMap adds the components to configMap. The components are compressed if their size would be over 950KB.
func addComponentsToConfigMap(cm corev1.ConfigMap, objs []unstructured.Unstructured) (corev1.ConfigMap, error) {
	combined, err := utilyaml.FromUnstructured(objs)
	if err != nil {
		return corev1.ConfigMap{}, fmt.Errorf("error converting unstructured object to YAML: %w", err)
	}

	compressionRequired := len(combined) > 950_000 // 98KB under the configMap 1MiB limit to leave space for rest of the configMap
	if compressionRequired {
		compressed, err := compressToZstd(combined)
		if err != nil {
			return corev1.ConfigMap{}, fmt.Errorf("failed to compress components: %w", err)
		}

		cm.BinaryData = map[string][]byte{
			"components-zstd": compressed.Bytes(),
		}
	} else {
		cm.Data["components"] = string(combined)
	}

	return cm, nil
}

// compressToZstd compressses bytes using zstd
func compressToZstd(data []byte) (bytes.Buffer, error) {
	var compressed bytes.Buffer

	writer, err := zstd.NewWriter(&compressed, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to initialize zstd writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write(data)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to encode using zstd: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to close zstd writer: %w", err)
	}

	return compressed, nil
}

func importProvider(p provider) error {
	fmt.Printf("Processing provider %s\n", p.Name)

	// Load manifests from github for specific provider

	// for Power VS the upstream cluster api provider name is ibmcloud
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/cmd/clusterctl/client/config/providers_client.go#L210-L214
	var initialProviderName string
	if p.Name == powerVSProvider {
		initialProviderName = powerVSProvider
		p.Name = ibmCloudProvider
	}

	if err := p.loadComponents(); err != nil {
		return fmt.Errorf("error loading provider components: %w", err)
	}
	if p.Name == powerVSProvider {
		p.Name = ibmCloudProvider
	}

	// Perform all manifest transformations

	// We need to perform Power VS specific customization which may not needed for ibmcloud
	if initialProviderName == powerVSProvider {
		p.Name = powerVSProvider
	}
	resourceMap := processObjects(p.components.Objs(), p.Name)

	// Write RBAC components to manifests, they will be managed by CVO
	if p.Name == powerVSProvider {
		p.Name = ibmCloudProvider
	}

	// Write modified infrastructure components to allow for code review.
	if err := p.writeProviderCustomizedComponents(resourceMap); err != nil {
		return fmt.Errorf("error writing %v file: %w", customizedComponentsFilename, err)
	}

	// Store provider components in the provider ConfigMap.
	cmFileName := fmt.Sprintf("%s04_cm.%s-%s.yaml", manifestPrefix, strings.ToLower(p.providerTypeName()), p.Name)
	if err := p.writeProviderComponentsConfigmap(cmFileName, resourceMap[otherKey]); err != nil {
		return fmt.Errorf("error writing provider ConfigMap: %w", err)
	}

	// Optionally write a separate CRD manifest file,
	// to apply CRDs directly via CVO rather than through the cluster-capi-operator,
	// useful in cases where the platform is not supported but some CRDs are needed
	// by other OCP operators other than the cluster-capi-operator.
	if len(resourceMap[crdKey]) > 0 {
		cmFileName := fmt.Sprintf("%s04_crd.%s-%s.yaml", manifestPrefix, strings.ToLower(p.providerTypeName()), p.Name)
		if err := p.writeProviderComponentsToManifest(cmFileName, resourceMap[crdKey]); err != nil {
			return fmt.Errorf("error writing provider CRDs: %w", err)
		}
	}

	return nil
}

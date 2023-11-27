package main

import (
	"fmt"
	"os"
	"path"
	"strings"

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
	powerVSProvider         = "powervs"
	ibmCloudProvider        = "ibmcloud"
	metadataFilePath        = "./metadata.yaml"
	kustomizeComponentsPath = "./config/default"
)

type provider struct {
	Name       string                    `json:"name"`
	PType      clusterctlv1.ProviderType `json:"type"`
	Version    string                    `json:"version"`
	components repository.Components
	metadata   []byte
}

// loadComponents loads components from the given provider.
func (p *provider) loadComponents() error {
	// Create new clusterctl config client
	configClient, err := configclient.New("")
	if err != nil {
		return err
	}

	// Create new clusterctl provider client
	providerConfig, err := configClient.Providers().Get(p.Name, p.PType)
	if err != nil {
		return err
	}

	// Set options for yaml processor
	options := repository.ComponentsOptions{
		TargetNamespace:     targetNamespace,
		SkipTemplateProcess: true,
	}

	// Compile assets using kustomize.
	rawComponents, err := fetchAndCompileComponents(path.Join(projDir, kustomizeComponentsPath))
	if err != nil {
		return err
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
		return err
	}

	content, err := os.ReadFile(path.Join(projDir, metadataFilePath))
	if err != nil {
		return fmt.Errorf("error while reading metadata file: %w", err)
	}

	p.metadata = content

	return err
}

func (p *provider) providerTypeName() string {
	return strings.ReplaceAll(strings.ToLower(string(p.PType)), "provider", "")
}

func (p *provider) writeProviderComponentsConfigmap(fileName string, objs []unstructured.Unstructured) error {
	combined, err := utilyaml.FromUnstructured(objs)
	if err != nil {
		return err
	}

	annotations := openshifAnnotations
	annotations[techPreviewAnnotation] = techPreviewAnnotationValue

	cm := &corev1.ConfigMap{
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
			"metadata":   string(p.metadata),
			"components": string(combined),
		},
	}

	cmYaml, err := yaml.Marshal(cm)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(*manifestsPath, fileName), ensureNewLine(cmYaml), 0600)
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
		return err
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

	// rbacFileName := fmt.Sprintf("%s03_rbac-roles.%s-%s.yaml", manifestPrefix, p.providerTypeName(), p.Name)
	// if err := writeComponentsToManifests(rbacFileName, resourceMap[rbacKey]); err != nil {
	// 	return err
	// }

	// // Write CRD components to manifests, they will be managed by CVO
	// crdFileName := fmt.Sprintf("%s02_crd.%s-%s.yaml", manifestPrefix, p.providerTypeName(), p.Name)
	// if err := writeComponentsToManifests(crdFileName, resourceMap[crdKey]); err != nil {
	// 	return err
	// }

	// Write all other components(deployments, services, secret, etc) to a ConfigMap,
	// they will be managed by CAPI operator
	// fName := strings.ToLower(p.providerTypeName() + "-" + p.Name + ".yaml")
	cmFileName := fmt.Sprintf("%s04_cm.%s-%s.yaml", manifestPrefix, strings.ToLower(p.providerTypeName()), p.Name)
	if err := p.writeProviderComponentsConfigmap(cmFileName, resourceMap[otherKey]); err != nil {
		return err
	}

	return nil
}

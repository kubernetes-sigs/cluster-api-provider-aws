package main

import (
	"bytes"
	"strings"

	certmangerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	admissionregistration "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type resourceKey string

const (
	crdKey                                                            resourceKey = "crds"
	otherKey                                                          resourceKey = "other"
	managedByAnnotationValueClusterCAPIOperatorInfraClusterController             = "cluster-capi-operator-infracluster-controller"
)

var (
	openshiftAnnotations = map[string]string{
		"exclude.release.openshift.io/internal-openshift-hosted":      "true",
		"include.release.openshift.io/self-managed-high-availability": "true",
		"include.release.openshift.io/single-node-developer":          "true",
	}

	// Workload annotations are used by the workload admission webhook to modify pod
	// resources and correctly schedule them while also pinning them to specific CPUSets.
	// See for more info:
	// https://github.com/openshift/enhancements/blob/master/enhancements/workload-partitioning/wide-availability-workload-partitioning.md
	openshiftWorkloadAnnotation = map[string]string{
		"target.workload.openshift.io/management": `{"effect": "PreferredDuringScheduling"}`,
	}

	// featureSetAnnotationValue is a multiple-feature-sets annotation value
	// adhering to the: %s,%s,... notation defined in the openshift/library-go/pkg/manifest parser.
	featureSetAnnotationValue = "CustomNoUpgrade,TechPreviewNoUpgrade"
	featureSetAnnotationKey   = "release.openshift.io/feature-set"
)

func processObjects(objs []unstructured.Unstructured, providerName string) map[resourceKey][]unstructured.Unstructured {
	resourceMap := map[resourceKey][]unstructured.Unstructured{}
	providerConfigMapObjs := []unstructured.Unstructured{}
	crdObjs := []unstructured.Unstructured{}

	objs = addInfraClusterProtectionPolicy(objs, providerName)

	serviceSecretNames := findWebhookServiceSecretName(objs)

	for _, obj := range objs {
		providerCustomizations(&obj, providerName)
		switch obj.GetKind() {
		case "ClusterRole", "Role", "ClusterRoleBinding", "RoleBinding", "ServiceAccount":
			setOpenShiftAnnotations(obj, false)
			setNoUpgradeAnnotations(obj)
			providerConfigMapObjs = append(providerConfigMapObjs, obj)
		case "MutatingWebhookConfiguration":
			// Explicitly remove defaulting webhooks for the cluster-api provider.
			// We don't need CAPI to set any default to the cluster object because
			// we have a custom controller for reconciling it.
			// For more information: https://issues.redhat.com/browse/OCPCLOUD-1506
			removeClusterDefaultingWebhooks(&obj)
			replaceCertManagerAnnotations(&obj)
			providerConfigMapObjs = append(providerConfigMapObjs, obj)
		case "ValidatingWebhookConfiguration":
			removeClusterValidatingWebhooks(&obj)
			replaceCertManagerAnnotations(&obj)
			providerConfigMapObjs = append(providerConfigMapObjs, obj)
		case "CustomResourceDefinition":
			replaceCertManagerAnnotations(&obj)
			removeConversionWebhook(&obj)
			setOpenShiftAnnotations(obj, true)
			// Apply NoUpgrade annotations unless IPAM CRDs,
			// as those are in General Availability.
			if !isCRDGroup(&obj, "ipam.cluster.x-k8s.io") {
				setNoUpgradeAnnotations(obj)
			}
			// Store Core CAPI CRDs in their own manifest to get them applied by CVO directly.
			// We want these to be installed independently from whether the cluster-capi-operator is enabled,
			// as other Openshift components rely on them.
			if providerName == coreCAPIProvider {
				crdObjs = append(crdObjs, obj)
			} else {
				providerConfigMapObjs = append(providerConfigMapObjs, obj)
			}
		case "Service":
			replaceCertMangerServiceSecret(&obj, serviceSecretNames)
			setOpenShiftAnnotations(obj, true)
			setNoUpgradeAnnotations(obj)
			providerConfigMapObjs = append(providerConfigMapObjs, obj)
		case "Deployment":
			customizeDeployments(&obj)
			if providerName == "operator" {
				setOpenShiftAnnotations(obj, false)
				setNoUpgradeAnnotations(obj)
			}
			providerConfigMapObjs = append(providerConfigMapObjs, obj)
		case "ValidatingAdmissionPolicy":
			providerConfigMapObjs = append(providerConfigMapObjs, obj)
		case "ValidatingAdmissionPolicyBinding":
			providerConfigMapObjs = append(providerConfigMapObjs, obj)
		case "Certificate", "Issuer", "Namespace", "Secret": // skip
		}
	}

	resourceMap[crdKey] = crdObjs
	resourceMap[otherKey] = providerConfigMapObjs

	return resourceMap
}

func setOpenShiftAnnotations(obj unstructured.Unstructured, merge bool) {
	if !merge || len(obj.GetAnnotations()) == 0 {
		obj.SetAnnotations(openshiftAnnotations)
	}

	anno := obj.GetAnnotations()
	if anno == nil {
		anno = map[string]string{}
	}

	for k, v := range openshiftAnnotations {
		anno[k] = v
	}
	obj.SetAnnotations(anno)
}

func setNoUpgradeAnnotations(obj unstructured.Unstructured) {
	anno := obj.GetAnnotations()
	if anno == nil {
		anno = map[string]string{}
	}

	anno[featureSetAnnotationKey] = featureSetAnnotationValue
	obj.SetAnnotations(anno)
}

func findWebhookServiceSecretName(objs []unstructured.Unstructured) map[string]string {
	serviceSecretNames := map[string]string{}
	certSecretNames := map[string]string{}

	secretFromCertNN := func(certNN string) (string, bool) {
		if len(certNN) == 0 {
			return "", false
		}
		certName := strings.Split(certNN, "/")[1]
		secretName, ok := certSecretNames[certName]
		if !ok || secretName == "" {
			return "", false
		}
		return secretName, true
	}
	// find service, then cert, then secret
	// return map[certName] = secretName
	for i, obj := range objs {
		switch obj.GetKind() {
		case "Certificate":
			cert := &certmangerv1.Certificate{}
			if err := scheme.Convert(&objs[i], cert, nil); err != nil {
				panic(err)
			}
			certSecretNames[cert.Name] = cert.Spec.SecretName
		}
	}
	for _, obj := range objs {
		switch obj.GetKind() {
		case "CustomResourceDefinition":
			crd := &apiextensionsv1.CustomResourceDefinition{}
			if err := scheme.Convert(&obj, crd, nil); err != nil {
				panic(err)
			}
			if certNN, ok := crd.Annotations["cert-manager.io/inject-ca-from"]; ok {
				secretName, ok := secretFromCertNN(certNN)
				if !ok {
					panic("can't find secret from cert: " + certNN)
				}
				if crd.Spec.Conversion != nil {
					serviceSecretNames[crd.Spec.Conversion.Webhook.ClientConfig.Service.Name] = secretName
				}
			}

		case "MutatingWebhookConfiguration":
			mwc := &admissionregistration.MutatingWebhookConfiguration{}
			if err := scheme.Convert(&obj, mwc, nil); err != nil {
				panic(err)
			}
			if certNN, ok := mwc.Annotations["cert-manager.io/inject-ca-from"]; ok {
				secretName, ok := secretFromCertNN(certNN)
				if !ok {
					panic("can't find secret from cert: " + certNN)
				}
				serviceSecretNames[mwc.Webhooks[0].ClientConfig.Service.Name] = secretName
			}

		case "ValidatingWebhookConfiguration":
			vwc := &admissionregistration.ValidatingWebhookConfiguration{}
			if err := scheme.Convert(&obj, vwc, nil); err != nil {
				panic(err)
			}
			if certNN, ok := vwc.Annotations["cert-manager.io/inject-ca-from"]; ok {
				secretName, ok := secretFromCertNN(certNN)
				if !ok {
					panic("can't find secret from cert:CustomResourceDefinition " + certNN)
				}
				serviceSecretNames[vwc.Webhooks[0].ClientConfig.Service.Name] = secretName
			}
		}
	}
	return serviceSecretNames
}

func customizeDeployments(obj *unstructured.Unstructured) {
	deployment := &appsv1.Deployment{}
	if err := scheme.Convert(obj, deployment, nil); err != nil {
		panic(err)
	}
	deployment.Spec.Template.Spec.PriorityClassName = "system-cluster-critical"

	deployment.Spec.Template.Annotations = mergeMaps(deployment.Spec.Template.Annotations, openshiftWorkloadAnnotation)

	for i := range deployment.Spec.Template.Spec.Containers {
		container := &deployment.Spec.Template.Spec.Containers[i]
		// Add resource requests
		container.Resources.Requests = corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("10m"),
			corev1.ResourceMemory: resource.MustParse("50Mi"),
		}
		// Remove any existing resource limits. See: https://github.com/openshift/enhancements/blob/master/CONVENTIONS.md#resources-and-limits
		container.Resources.Limits = corev1.ResourceList{}
		// Remove all image references if they are external, they will be substituted operator later
		if !strings.HasPrefix(container.Image, "registry.ci.openshift.org") {
			container.Image = "to.be/replaced:v99"
		}
		if container.Name == "kube-rbac-proxy" {
			container.Image = "registry.ci.openshift.org/openshift:kube-rbac-proxy"
		}

		// This helps with debugging and is enforced in OCP, see https://issues.redhat.com/browse/OCPBUGS-33170.
		container.TerminationMessagePolicy = corev1.TerminationMessageFallbackToLogsOnError
	}

	if err := scheme.Convert(deployment, obj, nil); err != nil {
		panic(err)
	}
}

func replaceCertManagerAnnotations(obj *unstructured.Unstructured) {
	anns := obj.GetAnnotations()
	if anns == nil {
		anns = map[string]string{}
	}
	if _, ok := anns["cert-manager.io/inject-ca-from"]; ok {
		anns["service.beta.openshift.io/inject-cabundle"] = "true"
		delete(anns, "cert-manager.io/inject-ca-from")
		obj.SetAnnotations(anns)
	}
}

func replaceCertMangerServiceSecret(obj *unstructured.Unstructured, serviceSecretNames map[string]string) {
	anns := obj.GetAnnotations()
	if anns == nil {
		anns = map[string]string{}
	}
	if name, ok := serviceSecretNames[obj.GetName()]; ok {
		anns["service.beta.openshift.io/serving-cert-secret-name"] = name
		obj.SetAnnotations(anns)
	}
}

func removeConversionWebhook(obj *unstructured.Unstructured) {
	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := scheme.Convert(obj, crd, nil); err != nil {
		panic(err)
	}
	crd.Spec.Conversion = nil
	if err := scheme.Convert(crd, obj, nil); err != nil {
		panic(err)
	}
}

// isCRDGroup checks whether the object provided is a CRD for the specified API group.
func isCRDGroup(obj *unstructured.Unstructured, group string) bool {
	switch obj.GetKind() {
	case "CustomResourceDefinition":
		crd := &apiextensions.CustomResourceDefinition{}
		if err := scheme.Convert(obj, crd, nil); err != nil {
			panic(err)
		}

		if crd.Spec.Group == group {
			return true
		}

		return false
	default:
		return false
	}
}

// ensureNewLine makes sure that there is one new line at the end of the file for git
func ensureNewLine(b []byte) []byte {
	return append(bytes.TrimRight(b, "\n"), []byte("\n")...)
}

func fetchAndCompileComponents(url string) ([]byte, error) {
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

	fSys := filesys.MakeFsOnDisk()

	m, err := k.Run(fSys, url)
	if err != nil {
		return nil, err
	}

	return m.AsYaml()
}

// Variadic function to merge maps of like kind.
// Note: keys of next map will override keys in previous map if previous map contains same key.
func mergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := map[K]V{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// addInfraClusterProtectionPolicy adds a Validating Admission Policy and Binding for protecting
// InfraClusters created by the cluster-capi-operator from deletion and editing.
func addInfraClusterProtectionPolicy(objs []unstructured.Unstructured, providerName string) []unstructured.Unstructured {
	policy := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "admissionregistration.k8s.io/v1beta1",
			"kind":       "ValidatingAdmissionPolicy",
			"metadata": map[string]interface{}{
				"name": "openshift-cluster-api-protect-" + providerName + "cluster",
			},
			"spec": map[string]interface{}{
				"failurePolicy": "Fail",
				"paramKind": map[string]interface{}{
					"apiVersion": "config.openshift.io/v1",
					"kind":       "Infrastructure",
				},
				"matchConstraints": map[string]interface{}{
					"resourceRules": []interface{}{
						map[string]interface{}{
							"apiGroups":   []interface{}{"infrastructure.cluster.x-k8s.io"},
							"apiVersions": []interface{}{"*"},
							"operations":  []interface{}{"DELETE"},
							"resources":   []interface{}{providerName + "clusters"},
						},
					},
				},
				"validations": []interface{}{
					map[string]interface{}{
						"expression": "!(oldObject.metadata.name == params.status.infrastructureName)",
						"message":    "InfraCluster resources with metadata.name corresponding to the cluster infrastructureName cannot be deleted.",
					},
				},
			},
		},
	}

	binding := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "admissionregistration.k8s.io/v1beta1",
			"kind":       "ValidatingAdmissionPolicyBinding",
			"metadata": map[string]interface{}{
				"name": "openshift-cluster-api-protect-" + providerName + "cluster",
			},
			"spec": map[string]interface{}{
				"paramRef": map[string]interface{}{
					"name":                    "cluster",
					"parameterNotFoundAction": "Deny",
				},
				"policyName":        "openshift-cluster-api-protect-" + providerName + "cluster",
				"validationActions": []interface{}{"Deny"},
				"matchResources": map[string]interface{}{
					"namespaceSelector": map[string]interface{}{
						"matchLabels": map[string]interface{}{
							"kubernetes.io/metadata.name": "openshift-cluster-api",
						},
					},
				},
			},
		},
	}

	return append(objs, *policy, *binding)
}

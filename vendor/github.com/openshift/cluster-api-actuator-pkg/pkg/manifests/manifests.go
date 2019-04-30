package manifests

import (
	"bytes"
	"fmt"
	"text/template"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"

	// apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func ClusterCRDManifest() *v1beta1.CustomResourceDefinition {
	return &v1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "clusters.cluster.k8s.io",
			Labels: map[string]string{
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "cluster.k8s.io",
			Version: "v1alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "clusters",
				Kind:   "Cluster",
			},
			Scope: "Namespaced",
			Subresources: &v1beta1.CustomResourceSubresources{
				Status: &v1beta1.CustomResourceSubresourceStatus{},
			},
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"spec": {
							Type: "object",
							Required: []string{
								"clusterNetwork",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"clusterNetwork": {
									Type: "object",
									Required: []string{
										"services",
										"pods",
										"serviceDomain",
									},
									Properties: map[string]v1beta1.JSONSchemaProps{
										"pods": {
											Type: "object",
											Required: []string{
												"cidrBlocks",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"cidrBlocks": {
													Type: "array",
													Items: &v1beta1.JSONSchemaPropsOrArray{
														Schema: &v1beta1.JSONSchemaProps{
															Type: "string",
														},
													},
												},
											},
										},
										"serviceDomain": {
											Type: "string",
										},
										"services": {
											Type: "object",
											Required: []string{
												"cidrBlocks",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"cidrBlocks": {
													Type: "array",
													Items: &v1beta1.JSONSchemaPropsOrArray{
														Schema: &v1beta1.JSONSchemaProps{
															Type: "string",
														},
													},
												},
											},
										},
									},
								},
								"providerSpec": {
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"value": {
											Type: "object",
										},
										"valueFrom": {
											Type: "object",
										},
									},
								},
							},
						},
						"status": {
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"providerStatus": {
									Type: "object",
								},
								"apiEndpoints": {
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
											Required: []string{
												"host",
												"port",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"host": {
													Type: "string",
												},
												"port": {
													Type:   "integer",
													Format: "int64",
												},
											},
										},
									},
								},
								"errorMessage": {
									Type: "string",
								},
								"errorReason": {
									Type: "string",
								},
							},
						},
						"apiVersion": {
							Type: "string",
						},
						"kind": {
							Type: "string",
						},
						"metadata": {
							Type: "object",
						},
					},
				},
			},
		},
	}
}

func MachineCRDManifest() *v1beta1.CustomResourceDefinition {
	return &v1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "machines.machine.openshift.io",
			Labels: map[string]string{
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "machine.openshift.io",
			Version: "v1beta1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "machines",
				Kind:   "Machine",
			},
			Scope: "Namespaced",
			Subresources: &v1beta1.CustomResourceSubresources{
				Status: &v1beta1.CustomResourceSubresourceStatus{},
			},
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": {
							Type: "string",
						},
						"kind": {
							Type: "string",
						},
						"metadata": {
							Type: "object",
						},
						"spec": {
							Type: "object",
							Required: []string{
								"providerSpec",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"versions": {
									Type: "object",
									Required: []string{
										"kubelet",
									},
									Properties: map[string]v1beta1.JSONSchemaProps{
										"controlPlane": {
											Type: "string",
										},
										"kubelet": {
											Type: "string",
										},
									},
								},
								"configSource": {
									Type: "object",
								},
								"metadata": {
									Type: "object",
								},
								"providerSpec": {
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"value": {
											Type: "object",
										},
										"valueFrom": {
											Type: "object",
										},
									},
								},
								"taints": {
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
							},
						},
						"status": {
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"providerStatus": {
									Type: "object",
								},
								"versions": {
									Type: "object",
									Required: []string{
										"kubelet",
									},
									Properties: map[string]v1beta1.JSONSchemaProps{
										"controlPlane": {
											Type: "string",
										},
										"kubelet": {
											Type: "string",
										},
									},
								},
								"addresses": {
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
								"conditions": {
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
								"errorMessage": {
									Type: "string",
								},
								"errorReason": {
									Type: "string",
								},
								"lastUpdated": {
									Type:   "string",
									Format: "date-time",
								},
								"nodeRef": {
									Type: "object",
								},
							},
						},
					},
				},
			},
		},
	}
}

func MachineSetCRDManifest() *v1beta1.CustomResourceDefinition {
	return &v1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "machinesets.machine.openshift.io",
			Labels: map[string]string{
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "machine.openshift.io",
			Version: "v1beta1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "machinesets",
				Kind:   "MachineSet",
			},
			Scope: v1beta1.ResourceScope("Namespaced"),
			Subresources: &v1beta1.CustomResourceSubresources{
				Status: &v1beta1.CustomResourceSubresourceStatus{},
			},
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"status": {
							Type: "object",
							Required: []string{
								"replicas",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"errorMessage": {
									Type: "string",
								},
								"errorReason": {
									Type: "string",
								},
								"fullyLabeledReplicas": {
									Type:   "integer",
									Format: "int32",
								},
								"observedGeneration": {
									Type:   "integer",
									Format: "int64",
								},
								"readyReplicas": {
									Type:   "integer",
									Format: "int32",
								},
								"replicas": {
									Type:   "integer",
									Format: "int32",
								},
								"availableReplicas": {
									Type:   "integer",
									Format: "int32",
								},
							},
						},
						"apiVersion": {
							Type: "string",
						},
						"kind": {
							Type: "string",
						},
						"metadata": {
							Type: "object",
						},
						"spec": {
							Type: "object",
							Required: []string{
								"selector",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"minReadySeconds": {
									Type:   "integer",
									Format: "int32",
								},
								"replicas": {
									Type:   "integer",
									Format: "int32",
								},
								"selector": {
									Type: "object",
								},
								"template": {
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"metadata": {
											Type: "object",
										},
										"spec": {
											Type: "object",
											Required: []string{
												"providerSpec",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"configSource": {
													Type: "object",
												},
												"metadata": {
													Type: "object",
												},
												"providerSpec": {
													Type: "object",
													Properties: map[string]v1beta1.JSONSchemaProps{
														"value": {
															Type: "object",
														},
														"valueFrom": {
															Type: "object",
														},
													},
												},
												"taints": {
													Type: "array",
													Items: &v1beta1.JSONSchemaPropsOrArray{
														Schema: &v1beta1.JSONSchemaProps{
															Type: "object",
														},
													},
												},
												"versions": {
													Type: "object",
													Required: []string{
														"kubelet",
													},
													Properties: map[string]v1beta1.JSONSchemaProps{
														"controlPlane": {
															Type: "string",
														},
														"kubelet": {
															Type: "string",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func MachineDeploymentCRDManifest() *v1beta1.CustomResourceDefinition {
	return &v1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "machinedeployments.machine.openshift.io",
			Labels: map[string]string{
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "machine.openshift.io",
			Version: "v1beta1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "machinedeployments",
				Kind:   "MachineDeployment",
			},
			Scope: "Namespaced",
			Subresources: &v1beta1.CustomResourceSubresources{
				Status: &v1beta1.CustomResourceSubresourceStatus{},
			},
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": {
							Type: "string",
						},
						"kind": {
							Type: "string",
						},
						"metadata": {
							Type: "object",
						},
						"spec": {
							Type: "object",
							Required: []string{
								"selector",
								"template",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"paused": {
									Type: "boolean",
								},
								"progressDeadlineSeconds": {
									Type:   "integer",
									Format: "int32",
								},
								"replicas": {
									Type:   "integer",
									Format: "int32",
								},
								"revisionHistoryLimit": {
									Type:   "integer",
									Format: "int32",
								},
								"selector": {
									Type: "object",
								},
								"strategy": {
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"rollingUpdate": {
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"maxSurge": {
													Type: "object",
												},
												"maxUnavailable": {
													Type: "object",
												},
											},
										},
										"type": {
											Type: "string",
										},
									},
								},
								"template": {
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"metadata": {
											Type: "object",
										},
										"spec": {
											Type: "object",
											Required: []string{
												"providerSpec",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"versions": {
													Type: "object",
													Required: []string{
														"kubelet",
													},
													Properties: map[string]v1beta1.JSONSchemaProps{
														"controlPlane": {
															Type: "string",
														},
														"kubelet": {
															Type: "string",
														},
													},
												},
												"configSource": {
													Type: "object",
												},
												"metadata": {
													Type: "object",
												},
												"providerSpec": {
													Type: "object",
													Properties: map[string]v1beta1.JSONSchemaProps{
														"valueFrom": {
															Type: "object",
														},
														"value": {
															Type: "object",
														},
													},
												},
												"taints": {
													Type: "array",
													Items: &v1beta1.JSONSchemaPropsOrArray{
														Schema: &v1beta1.JSONSchemaProps{
															Type: "object",
														},
													},
												},
											},
										},
									},
								},
								"minReadySeconds": {
									Type:   "integer",
									Format: "int32",
								},
							},
						},
						"status": {
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"readyReplicas": {
									Type:   "integer",
									Format: "int32",
								},
								"replicas": {
									Type:   "integer",
									Format: "int32",
								},
								"unavailableReplicas": {
									Type:   "integer",
									Format: "int32",
								},
								"updatedReplicas": {
									Type:   "integer",
									Format: "int32",
								},
								"availableReplicas": {
									Type:   "integer",
									Format: "int32",
								},
								"observedGeneration": {
									Type:   "integer",
									Format: "int64",
								},
							},
						},
					},
				},
			},
		},
	}
}

func ClusterRoleManifest() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "manager-role",
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"update",
					"patch",
					"delete",
				},
				APIGroups: []string{
					"machine.openshift.io",
				},
				Resources: []string{
					"clusters",
					"clusters/status",
					"machines",
					"machines/status",
					"machinesets",
					"machinesets/status",
					"machinedeployments",
					"machinedeployments/status",
				},
			},
			{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"update",
					"patch",
					"delete",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
				},
			},
			{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"secrets",
				},
			},
			{
				Verbs: []string{
					"create",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods/eviction",
				},
			},
			{
				Verbs: []string{
					"list",
					"watch",
					"get",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
				},
			},
			{
				Verbs: []string{
					"list",
					"watch",
					"get",
				},
				APIGroups: []string{
					"extensions",
				},
				Resources: []string{
					"daemonsets",
				},
			},
		},
	}
}

func ClusterRoleBinding(clusterAPINamespace string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "manager-rolebinding",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: clusterAPINamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "manager-role",
		},
	}
}

func ClusterAPIControllersDeployment(clusterAPINamespace, machineControllerImage, machineManagerImage, nodelinkControllerImage, ActuatorPrivateKey string) *appsv1.Deployment {
	var replicas int32 = 1
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi-controllers",
			Namespace: clusterAPINamespace,
			Labels: map[string]string{
				"api": "clusterapi",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"api": "clusterapi",
				},
			},
			Replicas: &replicas,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"api": "clusterapi",
					},
				},
				Spec: apiv1.PodSpec{
					NodeSelector: map[string]string{
						"node-role.kubernetes.io/master": "",
					},
					Tolerations: []apiv1.Toleration{
						{
							Effect: apiv1.TaintEffectNoSchedule,
							Key:    "node-role.kubernetes.io/master",
						},
						{
							Key:      "CriticalAddonsOnly",
							Operator: "Exists",
						},
						{
							Effect:   apiv1.TaintEffectNoExecute,
							Key:      "node.alpha.kubernetes.io/notReady",
							Operator: "Exists",
						},
						{
							Effect:   apiv1.TaintEffectNoExecute,
							Key:      "node.alpha.kubernetes.io/unreachable",
							Operator: "Exists",
						},
					},
					Containers: []apiv1.Container{
						{
							Name:  fmt.Sprintf("machine-controller"),
							Image: machineControllerImage,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/kubernetes",
								},
								{
									Name:      "certs",
									MountPath: "/etc/ssl/certs",
								},
								{
									Name:      "kubeadm",
									MountPath: "/usr/bin/kubeadm",
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
							Command: []string{"./machine-controller-manager"},
							Args: []string{
								"--logtostderr=true",
								"--v=3",
							},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("20Mi"),
								},
								Limits: apiv1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("30Mi"),
								},
							},
						},
						{
							Name: "nodelink-controller",
							// TODO(jchaloup): use other than the latest tag
							Image: nodelinkControllerImage,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/kubernetes",
								},
								{
									Name:      "certs",
									MountPath: "/etc/ssl/certs",
								},
							},
							Command: []string{"./nodelink-controller"},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("20Mi"),
								},
								Limits: apiv1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("30Mi"),
								},
							},
						},
						{
							Name:  "manager",
							Image: machineManagerImage,
							Command: []string{
								"/manager",
							},
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									"memory": resource.MustParse("30Mi"),
									"cpu":    resource.MustParse("100m"),
								},
								Requests: apiv1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("30Mi"),
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "config",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/etc/kubernetes",
								},
							},
						},
						{
							Name: "certs",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/etc/ssl/certs",
								},
							},
						},
						{
							Name: "kubeadm",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/usr/bin/kubeadm",
								},
							},
						},
					},
				},
			},
		},
	}

	if ActuatorPrivateKey != "" {
		var defaultMode int32 = 384
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, apiv1.Volume{
			Name: ActuatorPrivateKey,
			VolumeSource: apiv1.VolumeSource{
				Secret: &apiv1.SecretVolumeSource{
					SecretName:  ActuatorPrivateKey,
					DefaultMode: &defaultMode,
				},
			},
		})
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, apiv1.VolumeMount{
			Name:      ActuatorPrivateKey,
			MountPath: "/root/.ssh/actuator.pem",
			ReadOnly:  true,
		})
	}

	return deployment
}

func TestingMachine(clusterID string, namespace string, providerSpec machinev1beta1.ProviderSpec) *machinev1beta1.Machine {
	randomUUID := string(uuid.NewUUID())
	machine := &machinev1beta1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:         clusterID + "-machine-" + randomUUID[:6],
			Namespace:    namespace,
			GenerateName: "vs-master-",
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: machinev1beta1.MachineSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"node-role.kubernetes.io/compute": "",
				},
			},
			ProviderSpec: providerSpec,
		},
	}

	return machine
}

func MasterMachine(clusterID, namespace string, providerSpec machinev1beta1.ProviderSpec) *machinev1beta1.Machine {
	randomUUID := string(uuid.NewUUID())
	machine := &machinev1beta1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:         clusterID + "-master-machine-" + randomUUID[:6],
			Namespace:    namespace,
			GenerateName: "vs-master-",
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: machinev1beta1.MachineSpec{
			ProviderSpec: providerSpec,
		},
	}

	return machine
}

func MasterMachineUserDataSecret(secretName, namespace string, apiserverCertExtraSans []string) (*apiv1.Secret, error) {
	params := userDataParams{
		ApiserverCertExtraSans: apiserverCertExtraSans,
	}
	t, err := template.New("masteruserdata").Parse(masterUserDataBlob)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, params)
	if err != nil {
		return nil, err
	}

	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"userData": []byte(buf.String()),
		},
	}, nil
}

func WorkerMachineUserDataSecret(secretName, namespace, masterIP string) (*apiv1.Secret, error) {
	params := userDataParams{
		MasterIP: masterIP,
	}
	t, err := template.New("workeruserdata").Parse(workerUserDataBlob)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, params)
	if err != nil {
		return nil, err
	}

	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"userData": []byte(buf.String()),
		},
	}, nil
}

func WorkerMachineSet(clusterID, namespace string, providerSpec machinev1beta1.ProviderSpec) *machinev1beta1.MachineSet {
	var replicas int32 = 1
	randomUUID := string(uuid.NewUUID())
	return &machinev1beta1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:         clusterID + "-worker-machineset-" + randomUUID[:6],
			Namespace:    namespace,
			GenerateName: clusterID + "-worker-machine-" + randomUUID[:6] + "-",
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: machinev1beta1.MachineSetSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"machine.openshift.io/cluster-api-machineset": clusterID + "-worker-machineset-" + randomUUID[:6],
					"machine.openshift.io/cluster-api-cluster":    clusterID,
				},
			},
			Replicas: &replicas,
			Template: machinev1beta1.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: clusterID + "-worker-machine-" + randomUUID[:6] + "-",
					Labels: map[string]string{
						"machine.openshift.io/cluster-api-machineset": clusterID + "-worker-machineset-" + randomUUID[:6],
						"machine.openshift.io/cluster-api-cluster":    clusterID,
					},
				},
				Spec: machinev1beta1.MachineSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"node-role.kubernetes.io/compute": "",
						},
					},
					ProviderSpec: providerSpec,
				},
			},
		},
	}
}

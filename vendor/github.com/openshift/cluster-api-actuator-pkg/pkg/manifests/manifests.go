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
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
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
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Required: []string{
								"clusterNetwork",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"clusterNetwork": v1beta1.JSONSchemaProps{
									Type: "object",
									Required: []string{
										"services",
										"pods",
										"serviceDomain",
									},
									Properties: map[string]v1beta1.JSONSchemaProps{
										"pods": v1beta1.JSONSchemaProps{
											Type: "object",
											Required: []string{
												"cidrBlocks",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"cidrBlocks": v1beta1.JSONSchemaProps{
													Type: "array",
													Items: &v1beta1.JSONSchemaPropsOrArray{
														Schema: &v1beta1.JSONSchemaProps{
															Type: "string",
														},
													},
												},
											},
										},
										"serviceDomain": v1beta1.JSONSchemaProps{
											Type: "string",
										},
										"services": v1beta1.JSONSchemaProps{
											Type: "object",
											Required: []string{
												"cidrBlocks",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"cidrBlocks": v1beta1.JSONSchemaProps{
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
								"providerConfig": v1beta1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"value": v1beta1.JSONSchemaProps{
											Type: "object",
										},
										"valueFrom": v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
							},
						},
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"providerStatus": v1beta1.JSONSchemaProps{
									Type: "object",
								},
								"apiEndpoints": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
											Required: []string{
												"host",
												"port",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"host": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"port": v1beta1.JSONSchemaProps{
													Type:   "integer",
													Format: "int64",
												},
											},
										},
									},
								},
								"errorMessage": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"errorReason": v1beta1.JSONSchemaProps{
									Type: "string",
								},
							},
						},
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
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
			Name: "machines.cluster.k8s.io",
			Labels: map[string]string{
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "cluster.k8s.io",
			Version: "v1alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "machines",
				Kind:   "Machine",
			},
			Scope: "Namespaced",
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
							Type: "object",
						},
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Required: []string{
								"providerConfig",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"versions": v1beta1.JSONSchemaProps{
									Type: "object",
									Required: []string{
										"kubelet",
									},
									Properties: map[string]v1beta1.JSONSchemaProps{
										"controlPlane": v1beta1.JSONSchemaProps{
											Type: "string",
										},
										"kubelet": v1beta1.JSONSchemaProps{
											Type: "string",
										},
									},
								},
								"configSource": v1beta1.JSONSchemaProps{
									Type: "object",
								},
								"metadata": v1beta1.JSONSchemaProps{
									Type: "object",
								},
								"providerConfig": v1beta1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"value": v1beta1.JSONSchemaProps{
											Type: "object",
										},
										"valueFrom": v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
								"taints": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
							},
						},
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"providerStatus": v1beta1.JSONSchemaProps{
									Type: "object",
								},
								"versions": v1beta1.JSONSchemaProps{
									Type: "object",
									Required: []string{
										"kubelet",
									},
									Properties: map[string]v1beta1.JSONSchemaProps{
										"controlPlane": v1beta1.JSONSchemaProps{
											Type: "string",
										},
										"kubelet": v1beta1.JSONSchemaProps{
											Type: "string",
										},
									},
								},
								"addresses": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
								"conditions": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
										},
									},
								},
								"errorMessage": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"errorReason": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"lastUpdated": v1beta1.JSONSchemaProps{
									Type:   "string",
									Format: "date-time",
								},
								"nodeRef": v1beta1.JSONSchemaProps{
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
			Name: "machinesets.cluster.k8s.io",
			Labels: map[string]string{
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "cluster.k8s.io",
			Version: "v1alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "machinesets",
				Kind:   "MachineSet",
			},
			Scope: v1beta1.ResourceScope("Namespaced"),
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Required: []string{
								"replicas",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"errorMessage": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"errorReason": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"fullyLabeledReplicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"observedGeneration": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int64",
								},
								"readyReplicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"replicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"availableReplicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
							},
						},
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
							Type: "object",
						},
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Required: []string{
								"selector",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"minReadySeconds": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"replicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"selector": v1beta1.JSONSchemaProps{
									Type: "object",
								},
								"template": v1beta1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"metadata": v1beta1.JSONSchemaProps{
											Type: "object",
										},
										"spec": v1beta1.JSONSchemaProps{
											Type: "object",
											Required: []string{
												"providerConfig",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"configSource": v1beta1.JSONSchemaProps{
													Type: "object",
												},
												"metadata": v1beta1.JSONSchemaProps{
													Type: "object",
												},
												"providerConfig": v1beta1.JSONSchemaProps{
													Type: "object",
													Properties: map[string]v1beta1.JSONSchemaProps{
														"value": v1beta1.JSONSchemaProps{
															Type: "object",
														},
														"valueFrom": v1beta1.JSONSchemaProps{
															Type: "object",
														},
													},
												},
												"taints": v1beta1.JSONSchemaProps{
													Type: "array",
													Items: &v1beta1.JSONSchemaPropsOrArray{
														Schema: &v1beta1.JSONSchemaProps{
															Type: "object",
														},
													},
												},
												"versions": v1beta1.JSONSchemaProps{
													Type: "object",
													Required: []string{
														"kubelet",
													},
													Properties: map[string]v1beta1.JSONSchemaProps{
														"controlPlane": v1beta1.JSONSchemaProps{
															Type: "string",
														},
														"kubelet": v1beta1.JSONSchemaProps{
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
			Name: "machinedeployments.cluster.k8s.io",
			Labels: map[string]string{
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "cluster.k8s.io",
			Version: "v1alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "machinedeployments",
				Kind:   "MachineDeployment",
			},
			Scope: "Namespaced",
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
							Type: "object",
						},
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Required: []string{
								"selector",
								"template",
							},
							Properties: map[string]v1beta1.JSONSchemaProps{
								"paused": v1beta1.JSONSchemaProps{
									Type: "boolean",
								},
								"progressDeadlineSeconds": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"replicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"revisionHistoryLimit": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"selector": v1beta1.JSONSchemaProps{
									Type: "object",
								},
								"strategy": v1beta1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"rollingUpdate": v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"maxSurge": v1beta1.JSONSchemaProps{
													Type: "object",
												},
												"maxUnavailable": v1beta1.JSONSchemaProps{
													Type: "object",
												},
											},
										},
										"type": v1beta1.JSONSchemaProps{
											Type: "string",
										},
									},
								},
								"template": v1beta1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"metadata": v1beta1.JSONSchemaProps{
											Type: "object",
										},
										"spec": v1beta1.JSONSchemaProps{
											Type: "object",
											Required: []string{
												"providerConfig",
											},
											Properties: map[string]v1beta1.JSONSchemaProps{
												"versions": v1beta1.JSONSchemaProps{
													Type: "object",
													Required: []string{
														"kubelet",
													},
													Properties: map[string]v1beta1.JSONSchemaProps{
														"controlPlane": v1beta1.JSONSchemaProps{
															Type: "string",
														},
														"kubelet": v1beta1.JSONSchemaProps{
															Type: "string",
														},
													},
												},
												"configSource": v1beta1.JSONSchemaProps{
													Type: "object",
												},
												"metadata": v1beta1.JSONSchemaProps{
													Type: "object",
												},
												"providerConfig": v1beta1.JSONSchemaProps{
													Type: "object",
													Properties: map[string]v1beta1.JSONSchemaProps{
														"valueFrom": v1beta1.JSONSchemaProps{
															Type: "object",
														},
														"value": v1beta1.JSONSchemaProps{
															Type: "object",
														},
													},
												},
												"taints": v1beta1.JSONSchemaProps{
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
								"minReadySeconds": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
							},
						},
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"readyReplicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"replicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"unavailableReplicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"updatedReplicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"availableReplicas": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"observedGeneration": v1beta1.JSONSchemaProps{
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
			rbacv1.PolicyRule{
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
					"cluster.k8s.io",
				},
				Resources: []string{
					"clusters",
				},
			},
			rbacv1.PolicyRule{
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
					"cluster.k8s.io",
				},
				Resources: []string{
					"machines",
				},
			},
			rbacv1.PolicyRule{
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
					"cluster.k8s.io",
				},
				Resources: []string{
					"machinedeployments",
				},
			},
			rbacv1.PolicyRule{
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
					"cluster.k8s.io",
				},
				Resources: []string{
					"machinesets",
				},
			},
			rbacv1.PolicyRule{
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
					"cluster.k8s.io",
				},
				Resources: []string{
					"machines",
				},
			},
			rbacv1.PolicyRule{
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
			rbacv1.PolicyRule{
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
					"cluster.k8s.io",
				},
				Resources: []string{
					"machines",
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
			rbacv1.Subject{
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

func ManagerManifest(clusterAPINamespace, managerImage string) *appsv1.StatefulSet {
	var terminationGracePeriodSeconds int64 = 10
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controller-manager",
			Namespace: clusterAPINamespace,
			Labels: map[string]string{
				"control-plane":           "controller-manager",
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"control-plane":           "controller-manager",
					"controller-tools.k8s.io": "1.0",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"control-plane":           "controller-manager",
						"controller-tools.k8s.io": "1.0",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:  "manager",
							Image: managerImage,
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
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Tolerations: []apiv1.Toleration{
						apiv1.Toleration{
							Key:    "node-role.kubernetes.io/master",
							Effect: "NoSchedule",
						},
						apiv1.Toleration{
							Key:      "CriticalAddonsOnly",
							Operator: "Exists",
						},
						apiv1.Toleration{
							Key:      "node.alpha.kubernetes.io/notReady",
							Operator: "Exists",
							Effect:   "NoExecute",
						},
						apiv1.Toleration{
							Key:      "node.alpha.kubernetes.io/unreachable",
							Operator: "Exists",
							Effect:   "NoExecute",
						},
					},
				},
			},
			ServiceName: "controller-manager-service",
		},
	}
}

func ManagerService(clusterAPINamespace string) *apiv1.Service {
	return &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controller-manager-service",
			Namespace: clusterAPINamespace,
			Labels: map[string]string{
				"control-plane":           "controller-manager",
				"controller-tools.k8s.io": "1.0",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port: 443,
				},
			},
			Selector: map[string]string{
				"control-plane":           "controller-manager",
				"controller-tools.k8s.io": "1.0",
			},
		},
	}
}

func ClusterAPIControllersDeployment(clusterAPINamespace, actuatorImage string, ActuatorPrivateKey string) *appsv1.Deployment {
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
							Image: actuatorImage,
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
								"--log-level=debug",
								"--kubeconfig=/etc/kubernetes/admin.conf",
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
							Image: "openshift/origin-machine-api-operator:latest",
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
							Args:    []string{"--kubeconfig=/etc/kubernetes/admin.conf"},
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
		deployment.Spec.Template.Spec.Containers[1].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[1].VolumeMounts, apiv1.VolumeMount{
			Name:      ActuatorPrivateKey,
			MountPath: "/root/.ssh/actuator.pem",
			ReadOnly:  true,
		})
	}

	return deployment
}

func TestingMachine(clusterID string, namespace string, providerConfig clusterv1alpha1.ProviderConfig) *clusterv1alpha1.Machine {
	randomUUID := string(uuid.NewUUID())
	machine := &clusterv1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:         clusterID + "-machine-" + randomUUID[:6],
			Namespace:    namespace,
			GenerateName: "vs-master-",
			Labels: map[string]string{
				"sigs.k8s.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: clusterv1alpha1.MachineSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"node-role.kubernetes.io/compute": "",
				},
			},
			ProviderConfig: providerConfig,
			Versions: clusterv1alpha1.MachineVersionInfo{
				Kubelet:      "1.10.1",
				ControlPlane: "1.10.1",
			},
		},
	}

	return machine
}

func MasterMachine(clusterID, namespace string, providerConfig clusterv1alpha1.ProviderConfig) *clusterv1alpha1.Machine {
	randomUUID := string(uuid.NewUUID())
	machine := &clusterv1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:         clusterID + "-master-machine-" + randomUUID[:6],
			Namespace:    namespace,
			GenerateName: "vs-master-",
			Labels: map[string]string{
				"sigs.k8s.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: clusterv1alpha1.MachineSpec{
			ProviderConfig: providerConfig,
			Versions: clusterv1alpha1.MachineVersionInfo{
				Kubelet:      "1.10.1",
				ControlPlane: "1.10.1",
			},
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

func WorkerMachineSet(clusterID, namespace string, providerConfig clusterv1alpha1.ProviderConfig) *clusterv1alpha1.MachineSet {
	var replicas int32 = 1
	randomUUID := string(uuid.NewUUID())
	return &clusterv1alpha1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:         clusterID + "-worker-machineset-" + randomUUID[:6],
			Namespace:    namespace,
			GenerateName: clusterID + "-worker-machine-" + randomUUID[:6] + "-",
			Labels: map[string]string{
				"sigs.k8s.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: clusterv1alpha1.MachineSetSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"sigs.k8s.io/cluster-api-machineset": clusterID + "-worker-machineset-" + randomUUID[:6],
					"sigs.k8s.io/cluster-api-cluster":    clusterID,
				},
			},
			Replicas: &replicas,
			Template: clusterv1alpha1.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: clusterID + "-worker-machine-" + randomUUID[:6] + "-",
					Labels: map[string]string{
						"sigs.k8s.io/cluster-api-machineset": clusterID + "-worker-machineset-" + randomUUID[:6],
						"sigs.k8s.io/cluster-api-cluster":    clusterID,
					},
				},
				Spec: clusterv1alpha1.MachineSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"node-role.kubernetes.io/compute": "",
						},
					},
					ProviderConfig: providerConfig,
					Versions: clusterv1alpha1.MachineVersionInfo{
						Kubelet:      "1.10.1",
						ControlPlane: "1.10.1",
					},
				},
			},
		},
	}
}

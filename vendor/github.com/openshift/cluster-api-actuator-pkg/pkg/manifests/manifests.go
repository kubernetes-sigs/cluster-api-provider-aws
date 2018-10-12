package manifests

import (
	"bytes"
	"fmt"
	"text/template"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/cert/triple"
	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func ClusterAPIServerAPIServiceObjects(clusterAPINamespace string) (*apiv1.Secret, *apiregistrationv1beta1.APIService, error) {
	// Copied from the https://github.com/kubernetes-sigs/cluster-api/blob/master/pkg/deployer/clusterapiserver.go#L46 (getApiServerCertsForNamespace)
	name := "clusterapi"

	caKeyPair, err := triple.NewCA(fmt.Sprintf("%s-certificate-authority", name))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create root-ca: %v", err)
	}

	apiServerKeyPair, err := triple.NewServerKeyPair(
		caKeyPair,
		fmt.Sprintf("%s.%s.svc", name, clusterAPINamespace),
		name,
		clusterAPINamespace,
		"cluster.local",
		[]string{},
		[]string{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create apiserver key pair: %v", err)
	}

	return &apiv1.Secret{
			Type: "kubernetes.io/tls",
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cluster-apiserver-certs",
				Namespace: clusterAPINamespace,
				Labels: map[string]string{
					"api":       "clusterapi",
					"apiserver": "true",
				},
			},
			Data: map[string][]byte{
				"tls.crt": cert.EncodeCertPEM(apiServerKeyPair.Cert),
				"tls.key": cert.EncodePrivateKeyPEM(apiServerKeyPair.Key),
			},
		}, &apiregistrationv1beta1.APIService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "v1alpha1.cluster.k8s.io",
				Namespace: clusterAPINamespace,
				Labels: map[string]string{
					"api":       "clusterapi",
					"apiserver": "true",
				},
			},
			Spec: apiregistrationv1beta1.APIServiceSpec{
				Version:              "v1alpha1",
				Group:                "cluster.k8s.io",
				GroupPriorityMinimum: 2000,
				Service: &apiregistrationv1beta1.ServiceReference{
					Name:      "clusterapi",
					Namespace: clusterAPINamespace,
				},
				VersionPriority: 10,
				CABundle:        cert.EncodeCertPEM(caKeyPair.Cert),
			},
		}, nil
}

func ClusterAPIService(clusterAPINamespace string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi",
			Namespace: clusterAPINamespace,
			Labels: map[string]string{
				"api":       "clusterapi",
				"apiserver": "true",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:       443,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt(443),
				},
			},
			Selector: map[string]string{
				"api":       "clusterapi",
				"apiserver": "true",
			},
		},
	}
}

func ClusterAPIDeployment(clusterAPINamespace string) *appsv1beta2.Deployment {
	var replicas int32 = 1
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi-apiserver",
			Namespace: clusterAPINamespace,
			Labels: map[string]string{
				"api":       "clusterapi",
				"apiserver": "true",
			},
		},
		Spec: appsv1beta2.DeploymentSpec{
			// https://github.com/kubernetes/kubernetes/issues/51133
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"api":       "clusterapi",
					"apiserver": "true",
				},
			},
			Replicas: &replicas,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"api":       "clusterapi",
						"apiserver": "true",
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
							Name:  "apiserver",
							Image: "gcr.io/k8s-cluster-api/cluster-apiserver:0.0.6",
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "cluster-apiserver-certs",
									MountPath: "/apiserver.local.config/certificates",
									ReadOnly:  true,
								},
								{
									Name:      "config",
									MountPath: "/etc/kubernetes",
								},
								{
									Name:      "certs",
									MountPath: "/etc/ssl/certs",
								},
							},
							Command: []string{"./apiserver"},
							Args: []string{
								"--etcd-servers=http://etcd-clusterapi-svc:2379",
								"--tls-cert-file=/apiserver.local.config/certificates/tls.crt",
								"--tls-private-key-file=/apiserver.local.config/certificates/tls.key",
								"--audit-log-path=-",
								"--audit-log-maxage=0",
								"--audit-log-maxbackup=0",
								"--authorization-kubeconfig=/etc/kubernetes/admin.conf",
								"--kubeconfig=/etc/kubernetes/admin.conf",
							},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("50Mi"),
								},
								Limits: apiv1.ResourceList{
									"cpu":    resource.MustParse("300m"),
									"memory": resource.MustParse("200Mi"),
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "cluster-apiserver-certs",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "cluster-apiserver-certs",
								},
							},
						},
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
					},
				},
			},
		},
	}
}

func ClusterAPIControllersDeployment(clusterAPINamespace, actuatorImage string, ActuatorPrivateKey string) *appsv1beta2.Deployment {
	var replicas int32 = 1
	deployment := &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi-controllers",
			Namespace: clusterAPINamespace,
			Labels: map[string]string{
				"api": "clusterapi",
			},
		},
		Spec: appsv1beta2.DeploymentSpec{
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
							Name:  "controller-manager",
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
							},
							Command: []string{"./controller-manager"},
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
							Command: []string{"./machine-controller"},
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

func ClusterAPIRoleBinding(clusterAPINamespace string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi",
			Namespace: "kube-system",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "extension-apiserver-authentication-reader",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: clusterAPINamespace,
			},
		},
	}
}

func ClusterAPIEtcdCluster(clusterAPINamespace string) *appsv1beta2.StatefulSet {
	var terminationGracePeriodSeconds int64 = 10
	var replicas int32 = 1
	hostPathDirectoryOrCreate := apiv1.HostPathDirectoryOrCreate
	return &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "etcd-clusterapi",
			Namespace: clusterAPINamespace,
		},
		Spec: appsv1beta2.StatefulSetSpec{
			ServiceName: "etcd",
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "etcd",
				},
			},
			Replicas: &replicas,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "etcd",
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
							Name:  "etcd",
							Image: "k8s.gcr.io/etcd:3.1.12",
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "etcd-data-dir",
									MountPath: "/etcd-data-dir",
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "ETCD_DATA_DIR",
									Value: "/etcd-data-dir",
								},
							},
							Command: []string{
								"/usr/local/bin/etcd",
								"--listen-client-urls",
								"http://0.0.0.0:2379",
								"--advertise-client-urls",
								"http://localhost:2379",
							},
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 2379,
								},
							},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									"cpu":    resource.MustParse("100m"),
									"memory": resource.MustParse("50Mi"),
								},
								Limits: apiv1.ResourceList{
									"cpu":    resource.MustParse("200m"),
									"memory": resource.MustParse("300Mi"),
								},
							},
							ReadinessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Port: intstr.FromInt(2379),
										Path: "/health",
									},
								},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      2,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    1,
							},
							LivenessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Port: intstr.FromInt(2379),
										Path: "/health",
									},
								},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      2,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "etcd-data-dir",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/etc/kubernetes",
									Type: &hostPathDirectoryOrCreate,
								},
							},
						},
					},
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
				},
			},
		},
	}
}

func ClusterAPIEtcdService(clusterAPINamespace string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "etcd-clusterapi-svc",
			Namespace: clusterAPINamespace,
			Labels: map[string]string{
				"app": "etcd",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:       2379,
					TargetPort: intstr.FromInt(2379),
					Name:       "etcd",
				},
			},
			Selector: map[string]string{
				"app": "etcd",
			},
		},
	}
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

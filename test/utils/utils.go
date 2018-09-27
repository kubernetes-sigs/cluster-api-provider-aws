package utils

import (
	"bytes"
	"fmt"
	"os"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/client"
	clusterapiaproviderawsv1alpha1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// GenerateAwsCredentialsSecretFromEnv generates secret with AWS credentials
func GenerateAwsCredentialsSecretFromEnv(secretName, namespace string) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			awsclient.AwsCredsSecretIDKey:     []byte(os.Getenv("AWS_ACCESS_KEY_ID")),
			awsclient.AwsCredsSecretAccessKey: []byte(os.Getenv("AWS_SECRET_ACCESS_KEY")),
		},
	}
}

const (
	tlscrt = `-----BEGIN CERTIFICATE-----
MIIDYDCCAkigAwIBAgIIOtBgHbOzgecwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UE
AxMgY2x1c3RlcmFwaS1jZXJ0aWZpY2F0ZS1hdXRob3JpdHkwHhcNMTgwODA5MTQ0
MTU1WhcNMTkwODA5MTQ0MTU1WjAhMR8wHQYDVQQDExZjbHVzdGVyYXBpLmRlZmF1
bHQuc3ZjMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0C7AVOjCcnLp
Fx+JYzCrcaBGdzE34TylVM1XQha+du82fTv7qD0Y/wTaNv0eZmaH2IcuJ9ynGwIR
Zm3V428uls+MiAN9d8jUWoq3B/DiUZ1dDPCRYwvqTT02O70GF6Fe0jFMqJU5I9f2
OM99iyP9tiXYC423kBPL7m7d9TLXHKkhvFdFlUKZZI2kI47/oIbkm6UVwnFOJo47
zsa8T+jw6QHPmAvpigX1PpBLjyjkXbgWdrgI5rejiSVs9LIu0klFD/CsLR0o/F2x
9qkofc8f1HRXxLNlq9JuLB0BgYgxpIjBoZ5oLrbgxk2zuVOsH/6Kms3HW9PlWoxh
PiVoBa5EJQIDAQABo4GRMIGOMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggr
BgEFBQcDATBnBgNVHREEYDBeggpjbHVzdGVyYXBpghJjbHVzdGVyYXBpLmRlZmF1
bHSCFmNsdXN0ZXJhcGkuZGVmYXVsdC5zdmOCJGNsdXN0ZXJhcGkuZGVmYXVsdC5z
dmMuY2x1c3Rlci5sb2NhbDANBgkqhkiG9w0BAQsFAAOCAQEArZlNH6T0CWrUr8Vm
xvBkejNTbODsYuuhdqGgs9/KmoqLXSrdjWa++7NciqGsKoDZMIru2wzD84nRHDKR
5gyByRHF5W583EXQNjwgXxJpWnmlovRikproz99e9XBDG3x3mS9b+JjUbMxOv9oI
wWfmyH3wHsas3zZCCSis49n3xsARkRY6EswzmGKLno9UYIdnG2DzOe2Jv4N8HeN7
OLU9tLAiWFDabsnWqvbyOzL3GdjaWfUDjmFhX8EK8+5PMeREwu+GZiO1o1QEV688
5kzMP+8fwQVTV03vJi7QnhucpOhW0VZfoiOKt0kthPg+2R01p8hwG0qdx+U5kHzT
T3gOkA==
-----END CERTIFICATE-----`
	tlskey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0C7AVOjCcnLpFx+JYzCrcaBGdzE34TylVM1XQha+du82fTv7
qD0Y/wTaNv0eZmaH2IcuJ9ynGwIRZm3V428uls+MiAN9d8jUWoq3B/DiUZ1dDPCR
YwvqTT02O70GF6Fe0jFMqJU5I9f2OM99iyP9tiXYC423kBPL7m7d9TLXHKkhvFdF
lUKZZI2kI47/oIbkm6UVwnFOJo47zsa8T+jw6QHPmAvpigX1PpBLjyjkXbgWdrgI
5rejiSVs9LIu0klFD/CsLR0o/F2x9qkofc8f1HRXxLNlq9JuLB0BgYgxpIjBoZ5o
Lrbgxk2zuVOsH/6Kms3HW9PlWoxhPiVoBa5EJQIDAQABAoIBAQCw4khA3NP6cnBi
WUVepgfFr6yvsX4NPn4ro500ZibG31Go7sJQnDkU1YajmkWuNAfQjmtFK1JAvG0U
XtaRO/KV6Rs6pdyBXn4vwBTsBlwFhHN/fxfI1GLr5cqiz2TRxybN6V19D+1Q6zol
4waEprv3fAgpKOyC2o83s7Oblur3SaRlaXV33UYU67IhY7YPNUqvmphpyQJUTbEK
9DsLPQsAKwqEFOGNq5fDCfN4eRc6U5AV5JkKq561OdkC/tbyp/nX0nburkuPtI9f
F+o3sKZbqY2PPRgCmhXrlSk1mSbFCvu9U/lLyTy6PcSxkf83fMPdP1VwJwwWzjBa
08DAV9OBAoGBANhuELPuhY9tbCJAVjZH0GwV94yIadtMBnNxdCYdg4Z+bo4UsFO3
C3LweDbJ6vAWrbKI0tgzvnNpL4Y8E1QI2J6TvQ12z8auCZAR8OJcA+T1EISUgq6j
FU/VELIRa6gm8CX5WFbIjgfTzpaMn55WZMHSw/TyPuIZEP+UORzVCoZ9AoGBAPY+
rD/itlrG7INj+f0sTh7CTVcv6IkP5YasoW5/FHRzjuwK9u6dfyVQ04N7f2KnY36k
ezdu2w9G/7YA6QdklCPGDWmPrcDDeBaMeNL/TOc3qO1Q6U/jkF/fWuv032ijDCVD
rz5mM/02tcW8VMrvVgWkLrxUOsr1Ag0yObJtmRzJAoGAWimQH8VQMq4dDC/NOpO0
SjLki9EQeGE1lsY+4toMvuzQ1bPcuSNaS6nOCtUXYKmx9tx1Kch0oNPDDqLcUnfU
9ksJySAj8trx9OjkdwhqPumw1eqgfmxGJpnWeLg1JzoBdXBo0s5+DNi6CZHPtUC8
fNp29AYvGDXlFPQEzvQZjGkCgYBRhtp8pFD/qRCxR66C1eJfaLE2hpQUnQC/H/Sq
osRg8cmF+PNceSSZdDMzOvYn8YeNbGOnLLq2SilrVs3QNsqdNXtHUdyTD6R4wrVW
FlSd0N3LBJjabFtmgoqVyJMXD7R7ufcRT8EyuqRf/USNk8QFRiB7FeAJRikRuWlE
2+hvkQKBgQCE3B4vZKrbUwXuyXnuAiuVFQBSrjVJlpqY38+VucB0UCzUg4BegYCq
uLVzPuKQRUTVK66PB9KYFsSqGeUsEGpD0pBk3OoqN6Cdio0517IbeSp5tw/JaGnR
ORj2OTRfjZoPsqrOlUY8A3TVLZc408ECv15PedfW+2TY/2u36XkBFA==
-----END RSA PRIVATE KEY-----`
)

// ClusterAPIServerCertsSecret builds secret object with tls crt and key
func ClusterAPIServerCertsSecret() *apiv1.Secret {
	return &apiv1.Secret{
		Type: "kubernetes.io/tls",
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-apiserver-certs",
			Namespace: "default",
			Labels: map[string]string{
				"api":       "clusterapi",
				"apiserver": "true",
			},
		},
		Data: map[string][]byte{
			"tls.crt": []byte(tlscrt),
			"tls.key": []byte(tlskey),
		},
	}
}

const (
	cabundle = `-----BEGIN CERTIFICATE-----
MIIC9DCCAdygAwIBAgIBADANBgkqhkiG9w0BAQsFADArMSkwJwYDVQQDEyBjbHVz
dGVyYXBpLWNlcnRpZmljYXRlLWF1dGhvcml0eTAeFw0xODA4MDkxNDQxNTVaFw0y
ODA4MDYxNDQxNTVaMCsxKTAnBgNVBAMTIGNsdXN0ZXJhcGktY2VydGlmaWNhdGUt
YXV0aG9yaXR5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5u7BGB85
+3OOCHWCBB6+XSRNCqzXSk0BPZ9b0lnGAdeK4wqshvIHagCvXtFdLwCms1PWo3G0
T3a1chzMPRPSpWMmrWCz/Hnms8n9FEcf/ADbWrl8Wb9jBelseXiPKgraV8vzz7LG
rot6LWd4F5TiJ8Canhiwu1jOTd+ab4tqsFCSNBd4sMw4cihz3cW6h6vfNRmgeCiG
5zxKdQNPMOO+gMkB6IQvOpnLvYD1lCNDLLv4ssJ6w1mLeTvF0OAOp95lnOy41KIJ
V+TZTiLdvyhWQ2bDzL/loo6+GIBy5bLxgAWqTpd7nU94IsXJnCSPGm3S5N//zN2R
ysM55ryR6ULZZQIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUw
AwEB/zANBgkqhkiG9w0BAQsFAAOCAQEA2Br70cLOPW3X7lzJupe8ZAMC8RN+1KpN
1/1v6+R5Mq0KgJQsNqgP/YAdTtaoqTlhea2qLzpOH0f4yMekarHOweqf5K/Rp/7h
Mo4ocOVTfmw9+Vffx6OQTxqM6uvK3IwzfPIkna1alKzANiqVC9Q844Ls0n6D2Ck5
+n9Ro6MGywX2nEoP7vlRGvpwz11WFcqcOMjwY5uiIiuIR8hS6jNJbz9H/0Ng50wz
NHRNsymZGoLKX00An2rUVy0NwNYr400GEQBuppKl4rZ6l5PTx7fh7/3wAN3JfXE0
P66QocO5ZrXAqwezokmymjx31khIgFvNl8Tc1WhEf7MOmdSwWnQCPw==
-----END CERTIFICATE-----`
)

func ClusterAPIAPIService() *apiregistrationv1beta1.APIService {
	return &apiregistrationv1beta1.APIService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "v1alpha1.cluster.k8s.io",
			Namespace: "default",
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
				Namespace: "default",
			},
			VersionPriority: 10,
			CABundle:        []byte(cabundle),
		},
	}
}

func ClusterAPIService() *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi",
			Namespace: "default",
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

func ClusterAPIDeployment() *appsv1beta2.Deployment {
	var replicas int32 = 1
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi-apiserver",
			Namespace: "default",
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

func ClusterAPIControllersDeployment() *appsv1beta2.Deployment {
	var replicas int32 = 1
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clusterapi-controllers",
			Namespace: "default",
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
							Image: "gcr.io/k8s-cluster-api/aws-machine-controller:0.0.1",
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
							Name:  "aws-machine-controller",
							Image: "gcr.io/k8s-cluster-api/aws-machine-controller:0.0.1",
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
}

func ClusterAPIRoleBinding() *rbacv1.RoleBinding {
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
				Namespace: "default",
			},
		},
	}
}

func ClusterAPIEtcdCluster() *appsv1beta2.StatefulSet {
	var terminationGracePeriodSeconds int64 = 10
	var replicas int32 = 1
	hostPathDirectoryOrCreate := apiv1.HostPathDirectoryOrCreate
	return &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "etcd-clusterapi",
			Namespace: "default",
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

func ClusterAPIEtcdService() *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "etcd-clusterapi-svc",
			Namespace: "default",
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

func TestingMachine(awsCredentialsSecretName string, clusterID string, namespace string) *clusterv1alpha1.Machine {
	publicIP := true
	machinePc := &clusterapiaproviderawsv1alpha1.AWSMachineProviderConfig{
		AMI: clusterapiaproviderawsv1alpha1.AWSResourceReference{
			Filters: []clusterapiaproviderawsv1alpha1.Filter{
				{
					Name:   "tag:image_stage",
					Values: []string{"base"},
				},
				{
					Name:   "tag:operating_system",
					Values: []string{"rhel"},
				},
				{
					Name:   "tag:ready",
					Values: []string{"yes"},
				},
			},
		},
		CredentialsSecret: &apiv1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: clusterapiaproviderawsv1alpha1.Placement{
			Region:           "us-east-1",
			AvailabilityZone: "us-east-1a",
		},
		Subnet: clusterapiaproviderawsv1alpha1.AWSResourceReference{
			Filters: []clusterapiaproviderawsv1alpha1.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-worker-*", clusterID)},
				},
			},
		},
		Tags: []clusterapiaproviderawsv1alpha1.TagSpecification{
			{
				Name:  "openshift-node-group-config",
				Value: "node-config-master",
			},
			{
				Name:  "host-type",
				Value: "master",
			},
			{
				Name:  "sub-host-type",
				Value: "default",
			},
		},
		SecurityGroups: []clusterapiaproviderawsv1alpha1.AWSResourceReference{
			{
				Filters: []clusterapiaproviderawsv1alpha1.Filter{
					{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-*", clusterID)},
					},
				},
			},
		},
		PublicIP: &publicIP,
	}

	var buf bytes.Buffer
	if err := providerconfigv1.Encoder.Encode(machinePc, &buf); err != nil {
		panic(fmt.Errorf("AWSMachineProviderConfig encoding failed: %v", err))
	}

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
			ProviderConfig: clusterv1alpha1.ProviderConfig{
				Value: &runtime.RawExtension{Raw: buf.Bytes()},
			},
			Versions: clusterv1alpha1.MachineVersionInfo{
				Kubelet:      "1.10.1",
				ControlPlane: "1.10.1",
			},
		},
	}

	return machine
}

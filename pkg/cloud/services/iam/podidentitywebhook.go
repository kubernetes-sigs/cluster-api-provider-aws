package iam

import (
	"context"
	"errors"

	v14 "k8s.io/api/admissionregistration/v1"
	v13 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	podIdentityWebhookName  = "pod-identity-webhook"
	podIdentityWebhookImage = "amazon/amazon-eks-pod-identity-webhook:v0.4.0"

	labelNodeRoleMaster       = "node-role.kubernetes.io/master"
	labelNodeRoleControlPlane = "node-role.kubernetes.io/control-plane"
)

func reconcileServiceAccount(ctx context.Context, ns string, remoteClient client.Client) error {
	check := &corev1.ServiceAccount{}
	if err := remoteClient.Get(ctx, types.NamespacedName{
		Name:      podIdentityWebhookName,
		Namespace: ns,
	}, check); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if check.UID != "" {
		return nil
	}

	sa := &corev1.ServiceAccount{
		ObjectMeta: objectMeta(podIdentityWebhookName, ns),
	}

	return remoteClient.Create(ctx, sa)
}

func reconcileClusterRole(ctx context.Context, ns string, remoteClient client.Client) error {
	check := &v12.ClusterRole{}
	if err := remoteClient.Get(ctx, types.NamespacedName{
		Name:      podIdentityWebhookName,
		Namespace: ns,
	}, check); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if check.UID != "" {
		return nil
	}

	cr := &v12.ClusterRole{
		ObjectMeta: objectMeta(podIdentityWebhookName, ns),
		Rules: []v12.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups:     []string{""},
				Resources:     []string{"secrets"},
				Verbs:         []string{"get", "update", "patch"},
				ResourceNames: []string{podIdentityWebhookName},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{"certificates.k8s.io"},
				Resources: []string{"certificatesigningrequests"},
				Verbs:     []string{"create", "get", "list", "watch"},
			},
		},
	}

	return remoteClient.Create(ctx, cr)
}

func reconcileClusterRoleBinding(ctx context.Context, ns string, remoteClient client.Client) error {
	check := &v12.ClusterRoleBinding{}
	if err := remoteClient.Get(ctx, types.NamespacedName{
		Name: podIdentityWebhookName,
	}, check); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if check.UID != "" {
		return nil
	}

	crb := &v12.ClusterRoleBinding{
		ObjectMeta: objectMeta(podIdentityWebhookName, ""),
		RoleRef: v12.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     podIdentityWebhookName,
		},
		Subjects: []v12.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      podIdentityWebhookName,
				Namespace: ns,
			},
		},
	}

	return remoteClient.Create(ctx, crb)
}

func reconcileService(ctx context.Context, ns string, remoteClient client.Client) error {
	check := &corev1.Service{}
	if err := remoteClient.Get(ctx, types.NamespacedName{
		Name:      podIdentityWebhookName,
		Namespace: ns,
	}, check); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if check.UID != "" {
		return nil
	}
	service := &corev1.Service{
		ObjectMeta: objectMeta(podIdentityWebhookName, ns),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       443,
					TargetPort: intstr.FromInt(443),
				},
			},
			Selector: map[string]string{
				"app": podIdentityWebhookName,
			},
		},
	}

	return remoteClient.Create(ctx, service)
}

func reconcileDeployment(ctx context.Context, ns string, secret *corev1.Secret, remoteClient client.Client) error {
	check := &v13.Deployment{}
	if err := remoteClient.Get(ctx, types.NamespacedName{
		Name:      podIdentityWebhookName,
		Namespace: ns,
	}, check); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	nodeAffinity := &corev1.NodeAffinity{
		PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
			{
				Weight: 10,
				Preference: corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      labelNodeRoleMaster,
							Operator: corev1.NodeSelectorOpExists,
						}, {
							Key:      labelNodeRoleControlPlane,
							Operator: corev1.NodeSelectorOpExists,
						},
					},
				},
			}, {
				Weight: 10,
				Preference: corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      labelNodeRoleControlPlane,
							Operator: corev1.NodeSelectorOpExists,
						},
					},
				},
			},
		},
	}

	tolerations := []corev1.Toleration{
		{
			Key:    labelNodeRoleControlPlane,
			Effect: corev1.TaintEffectNoSchedule,
		}, {
			Key:    labelNodeRoleMaster,
			Effect: corev1.TaintEffectNoSchedule,
		},
	}

	if check.UID == "" {
		replicas := int32(1)

		deployment := &v13.Deployment{
			ObjectMeta: objectMeta(podIdentityWebhookName, ns),
			Spec: v13.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": podIdentityWebhookName,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": podIdentityWebhookName,
						},
					},
					Spec: corev1.PodSpec{
						Affinity:           &corev1.Affinity{NodeAffinity: nodeAffinity},
						Tolerations:        tolerations,
						ServiceAccountName: podIdentityWebhookName,
						Containers: []corev1.Container{
							{
								Name:            podIdentityWebhookName,
								Image:           podIdentityWebhookImage,
								ImagePullPolicy: corev1.PullIfNotPresent,
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "webhook-certs",
										MountPath: "/etc/webhook/certs/",
										ReadOnly:  false,
									},
								},
								Command: []string{
									"/webhook",
									"--in-cluster=false",
									"--namespace=" + ns,
									"--service-name=" + podIdentityWebhookName,
									"--annotation-prefix=eks.amazonaws.com",
									"--token-audience=sts.amazonaws.com",
									"--logtostderr",
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "webhook-certs",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: secret.Name,
									},
								},
							},
						},
					},
				},
			},
		}

		return remoteClient.Create(ctx, deployment)
	}

	needsUpdate := false
	if check.Spec.Template.Spec.Affinity == nil {
		check.Spec.Template.Spec.Affinity = &corev1.Affinity{}
	}
	if check.Spec.Template.Spec.Affinity.NodeAffinity == nil {
		check.Spec.Template.Spec.Affinity.NodeAffinity = &corev1.NodeAffinity{}
	}

	for _, aff := range nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
		found := false
		for _, a := range check.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
			if len(a.Preference.MatchExpressions) == len(aff.Preference.MatchExpressions) {
				for _, e := range a.Preference.MatchExpressions {
					for _, e2 := range aff.Preference.MatchExpressions {
						if e.Key == e2.Key && e.Operator == e2.Operator {
							found = true
						}
					}
				}
			}
		}
		if !found {
			check.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(check.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution, aff)
			needsUpdate = true
		}
	}

	for _, tol := range tolerations {
		found := false
		for _, t := range check.Spec.Template.Spec.Tolerations {
			if t.Key == tol.Key && t.Effect == tol.Effect {
				found = true
			}
		}
		if !found {
			check.Spec.Template.Spec.Tolerations = append(check.Spec.Template.Spec.Tolerations, tol)
			needsUpdate = true
		}
	}

	if needsUpdate {
		return remoteClient.Update(ctx, check)
	}

	return nil
}

func reconcileMutatingWebHook(ctx context.Context, ns string, secret *corev1.Secret, remoteClient client.Client) error {
	check := &v14.MutatingWebhookConfiguration{}
	if err := remoteClient.Get(ctx, types.NamespacedName{
		Name:      podIdentityWebhookName,
		Namespace: ns,
	}, check); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if check.UID != "" {
		return nil
	}

	caBundle, ok := secret.Data["ca.crt"]
	if !ok {
		return errors.New("no CA certificate for the pod identity webhook certificate")
	}

	mwhMeta := objectMeta(podIdentityWebhookName, ns)
	fail := v14.Ignore
	none := v14.SideEffectClassNone
	mutate := "/mutate"
	mwh := &v14.MutatingWebhookConfiguration{
		ObjectMeta: mwhMeta,
		Webhooks: []v14.MutatingWebhook{
			{
				Name:          podIdentityWebhookName + ".amazonaws.com",
				FailurePolicy: &fail,
				ClientConfig: v14.WebhookClientConfig{
					Service: &v14.ServiceReference{
						Name:      podIdentityWebhookName,
						Namespace: ns,
						Path:      &mutate,
					},
					CABundle: caBundle,
				},
				Rules: []v14.RuleWithOperations{
					{
						Operations: []v14.OperationType{v14.Create},
						Rule: v14.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{"v1"},
							Resources:   []string{"pods"},
						},
					},
				},
				SideEffects:             &none,
				AdmissionReviewVersions: []string{"v1beta1"},
			},
		},
	}

	return remoteClient.Create(ctx, mwh)
}

// reconcilePodIdentityWebhookComponents will create sa, cr, crb, service, deployment and a mutating webhook in kube-system. The
// only difference between this and upstream is we are using cert-manager in the management cluster to create a cert
// instead installing it in the work load cluster.
// https://github.com/aws/amazon-eks-pod-identity-webhook/tree/master/deploy
func reconcilePodIdentityWebhookComponents(ctx context.Context, ns string, secret *corev1.Secret, remoteClient client.Client) error {
	// TODO only creates the object if they don't exist, could create some comparison logic for updates
	if err := reconcileServiceAccount(ctx, ns, remoteClient); err != nil {
		return err
	}

	if err := reconcileClusterRole(ctx, ns, remoteClient); err != nil {
		return err
	}

	if err := reconcileClusterRoleBinding(ctx, ns, remoteClient); err != nil {
		return err
	}

	if err := reconcileService(ctx, ns, remoteClient); err != nil {
		return err
	}

	if err := reconcileDeployment(ctx, ns, secret, remoteClient); err != nil {
		return err
	}

	if err := reconcileMutatingWebHook(ctx, ns, secret, remoteClient); err != nil {
		return err
	}

	return nil
}

func objectMeta(name, namespace string) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Labels: map[string]string{
			v1beta1.ProviderLabelName: "infrastructure-aws",
		},
		Name: name,
	}
	if namespace != "" {
		meta.Namespace = namespace
	}
	return meta
}

// reconcileCertificateSecret takes a secret and moves it to the workload cluster.
func reconcileCertificateSecret(ctx context.Context, cert *corev1.Secret, remoteClient client.Client) error {
	// check if the secret was created by cert-manager
	certCheck := &corev1.Secret{}
	if err := remoteClient.Get(ctx, types.NamespacedName{
		Name:      cert.Name,
		Namespace: cert.Namespace,
	}, certCheck); err != nil && !apierrors.IsNotFound(err) {
		// will return not found if waiting for cert-manager and will reconcile again later due to error
		return err
	}

	if certCheck.UID == "" {
		cert.ResourceVersion = ""
		return remoteClient.Create(ctx, cert)
	}

	return nil
}

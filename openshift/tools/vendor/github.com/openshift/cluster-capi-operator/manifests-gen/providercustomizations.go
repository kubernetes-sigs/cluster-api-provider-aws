package main

import (
	"regexp"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func providerCustomizations(obj *unstructured.Unstructured, providerName string) {
	switch providerName {
	case "azure":
		azureCustomizations(obj)
	case "gcp":
		gcpCustomizations(obj)
	case powerVSProvider:
		powerVSCustomizations(obj)
	}
}

func azureCustomizations(obj *unstructured.Unstructured) {
	switch obj.GetKind() {
	case "Deployment":
		deployment := &appsv1.Deployment{}
		if err := scheme.Convert(obj, deployment, nil); err != nil {
			panic(err)
		}

		// Modify bootstrap secret keys as they don't match with what is created by CCO.
		for i := range deployment.Spec.Template.Spec.Containers {
			container := &deployment.Spec.Template.Spec.Containers[i]
			if container.Name == "manager" {
				for j := range container.Env {
					env := &container.Env[j]
					switch env.Name {
					case "AZURE_SUBSCRIPTION_ID":
						env.ValueFrom.SecretKeyRef.Key = "azure_subscription_id"
					case "AZURE_TENANT_ID":
						env.ValueFrom.SecretKeyRef.Key = "azure_tenant_id"
					case "AZURE_CLIENT_ID":
						env.ValueFrom.SecretKeyRef.Key = "azure_client_id"
					case "AZURE_CLIENT_SECRET":
						env.ValueFrom.SecretKeyRef.Key = "azure_client_secret"
					}
				}
			}
		}

		if err := scheme.Convert(deployment, obj, nil); err != nil {
			panic(err)
		}
	}
}

func gcpCustomizations(obj *unstructured.Unstructured) {
	switch obj.GetKind() {
	case "Deployment":
		deployment := &appsv1.Deployment{}
		if err := scheme.Convert(obj, deployment, nil); err != nil {
			panic(err)
		}

		// Modify bootstrap secret keys as they don't match with what is created by CCO.
		for i := range deployment.Spec.Template.Spec.Containers {
			container := &deployment.Spec.Template.Spec.Containers[i]
			if container.Name == "manager" {
				for j := range container.Env {
					env := &container.Env[j]
					switch env.Name {
					case "GOOGLE_APPLICATION_CREDENTIALS":
						env.Value = "/home/.gcp/service_account.json"
					}
				}
			}
		}

		if err := scheme.Convert(deployment, obj, nil); err != nil {
			panic(err)
		}
	}
}

func powerVSCustomizations(obj *unstructured.Unstructured) {
	switch obj.GetKind() {
	case "Deployment":
		deployment := &appsv1.Deployment{}
		if err := scheme.Convert(obj, deployment, nil); err != nil {
			panic(err)
		}

		for i := range deployment.Spec.Template.Spec.Containers {
			container := &deployment.Spec.Template.Spec.Containers[i]
			if container.Name == "manager" {
				for j := range container.Args {
					arg := &container.Args[j]
					if *arg == "--provider-id-fmt=${PROVIDER_ID_FORMAT:=v1}" {
						container.Args[j] = "--provider-id-fmt=v2"
					}
				}
			}
		}
		if err := scheme.Convert(deployment, obj, nil); err != nil {
			panic(err)
		}
	}
}

func removeClusterDefaultingWebhooks(obj *unstructured.Unstructured) {
	var providerWebhooks = regexp.MustCompile(`^default\.[a-z]+cluster[a-z]*\.infrastructure\.cluster\.x-k8s\.io$`)

	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{}
	if err := scheme.Convert(obj, mutatingWebhookConfiguration, nil); err != nil {
		panic(err)
	}

	webhooks := []admissionregistrationv1.MutatingWebhook{}
	for _, webhook := range mutatingWebhookConfiguration.Webhooks {
		// We don't need these specific webhooks.
		if webhook.Name == "default.cluster.cluster.x-k8s.io" || webhook.Name == "default.clusterclass.cluster.x-k8s.io" || webhook.Name == "default.clusterresourceset.addons.cluster.x-k8s.io" {
			continue
		}

		// We also don't need provider webhooks for clusters.
		if providerWebhooks.MatchString(webhook.Name) {
			continue
		}

		webhooks = append(webhooks, webhook)
	}

	mutatingWebhookConfiguration.Webhooks = webhooks

	if err := scheme.Convert(mutatingWebhookConfiguration, obj, nil); err != nil {
		panic(err)
	}
}

func removeClusterValidatingWebhooks(obj *unstructured.Unstructured) {
	var providerWebhooks = regexp.MustCompile(`^validation\.[a-z]+cluster[a-z]*\.infrastructure\.cluster\.x-k8s\.io$`)

	validatingWebhookConfiguration := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	if err := scheme.Convert(obj, validatingWebhookConfiguration, nil); err != nil {
		panic(err)
	}

	webhooks := []admissionregistrationv1.ValidatingWebhook{}
	for _, webhook := range validatingWebhookConfiguration.Webhooks {
		// We don't need these specific webhooks.
		if webhook.Name == "validation.cluster.cluster.x-k8s.io" || webhook.Name == "validation.clusterclass.cluster.x-k8s.io" || webhook.Name == "default.clusterresourceset.addons.cluster.x-k8s.io" {
			continue
		}

		// We also don't need provider webhooks for clusters.
		if providerWebhooks.MatchString(webhook.Name) {
			continue
		}

		webhooks = append(webhooks, webhook)
	}

	validatingWebhookConfiguration.Webhooks = webhooks

	if err := scheme.Convert(validatingWebhookConfiguration, obj, nil); err != nil {
		panic(err)
	}
}

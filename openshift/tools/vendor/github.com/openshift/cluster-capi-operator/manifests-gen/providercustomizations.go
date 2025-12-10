package main

import (
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

package util

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// ServiceAccountNamespaceFile contains path to the file that contains namespace
const ServiceAccountNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// GetNamespace returns the namespace of the pod where the code is running
func GetNamespace(namespaceFile string) (string, error) {
	data, err := ioutil.ReadFile(namespaceFile)
	if err != nil {
		return "", fmt.Errorf("failed to determine namespace from %s: %v", namespaceFile, err)
	}

	return strings.TrimSpace(string(data)), nil
}

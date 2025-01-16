package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

const (
	permissionsBoundaryFile = "/home/.aws/permissionsBoundary"
)

var (
	cachedPermissionsBoundary string
	cacheLock                 sync.Mutex
)

func isSecretPartition(partitionID string) bool {
	return partitionID == endpoints.AwsIsoPartitionID || partitionID == endpoints.AwsIsoBPartitionID
}

func GetPermissionsBoundary(partitionID string) (string, error) {

	if !isSecretPartition(partitionID) {
		return "", nil
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()

	if cachedPermissionsBoundary != "" {
		return cachedPermissionsBoundary, nil
	}

	permissionsBoundary, err := readPermissionsBoundaryFromFile(permissionsBoundaryFile)
	if err != nil {
		return "", err
	}

	cachedPermissionsBoundary = permissionsBoundary
	return cachedPermissionsBoundary, nil
}

func readPermissionsBoundaryFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to open permissions boundary file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read permissions boundary file: %w", err)
	}

	return "", nil
}

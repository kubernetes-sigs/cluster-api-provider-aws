package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestIsSecretPartition(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		partitionID    string
		expectedResult bool
	}{
		{
			name:           "Valid AWS ISO Partition",
			partitionID:    "aws-iso",
			expectedResult: true,
		},
		{
			name:           "Valid AWS ISO-B Partition",
			partitionID:    "aws-iso-b",
			expectedResult: true,
		},
		{
			name:           "Invalid Partition",
			partitionID:    "aws",
			expectedResult: false,
		},
		{
			name:           "Empty Partition ID",
			partitionID:    "",
			expectedResult: false,
		},
		{
			name:           "Random Partition ID",
			partitionID:    "random-partition",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSecretPartition(tt.partitionID)
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v for partitionID %s", tt.expectedResult, result, tt.partitionID)
			}
		})
	}
}

func TestReadPermissionsBoundaryFromFile(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		fileContent   string
		expectedARN   string
		expectedError bool
	}{
		{
			name:          "Valid file with ARN without new line",
			fileContent:   "arn:aws:iam::123456789012:policy/MyPermissionsBoundary",
			expectedARN:   "arn:aws:iam::123456789012:policy/MyPermissionsBoundary",
			expectedError: false,
		},
		{
			name:          "Valid file with ARN",
			fileContent:   "arn:aws:iam::123456789012:policy/MyPermissionsBoundary\n",
			expectedARN:   "arn:aws:iam::123456789012:policy/MyPermissionsBoundary",
			expectedError: false,
		},
		{
			name:          "Empty file",
			fileContent:   "",
			expectedARN:   "",
			expectedError: false,
		},
		{
			name:          "File with only whitespace",
			fileContent:   "\n \n",
			expectedARN:   "",
			expectedError: false,
		},
		{
			name:          "Multiple lines, pick first non-empty line",
			fileContent:   "\n \narn:aws:iam::123456789012:policy/MyPermissionsBoundary\nanother-line\n",
			expectedARN:   "arn:aws:iam::123456789012:policy/MyPermissionsBoundary",
			expectedError: false,
		},
		{
			name:          "File does not exist",
			fileContent:   "",
			expectedARN:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir := t.TempDir()
			filePath := fmt.Sprintf("%s/permissionsBoundary", tempDir)

			if tt.name != "File does not exist" {
				// Write test content to the file
				file, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
				defer file.Close()

				_, err = file.WriteString(tt.fileContent)
				if err != nil {
					t.Fatalf("failed to write to test file: %v", err)
				}
			}

			// Call the function
			arn, err := readPermissionsBoundaryFromFile(filePath)

			// Verify the output
			if arn != tt.expectedARN {
				t.Errorf("expected ARN: %s, got: %s", tt.expectedARN, arn)
			}

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err != nil)
			}
		})
	}
}
